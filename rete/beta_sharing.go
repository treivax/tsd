// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"sync/atomic"
	"time"
)

// ============================================================================
// BetaSharingRegistryImpl Implementation
// ============================================================================

// GetOrCreateJoinNode returns a shared JoinNode for the given signature.
func (bsr *BetaSharingRegistryImpl) GetOrCreateJoinNode(
	condition interface{},
	leftVars []string,
	rightVars []string,
	allVars []string,
	varTypes map[string]string,
	storage Storage,
) (*JoinNode, string, bool, error) {
	if !bsr.config.Enabled {
		// Beta sharing disabled - create unique node
		nodeID := fmt.Sprintf("join_%d", time.Now().UnixNano())
		var conditionMap map[string]interface{}
		if condition != nil {
			if cm, ok := condition.(map[string]interface{}); ok {
				conditionMap = cm
			}
		}
		node := NewJoinNode(nodeID, conditionMap, leftVars, rightVars, varTypes, storage)
		return node, nodeID, false, nil
	}

	// Create signature
	sig := &JoinNodeSignature{
		Condition: condition,
		LeftVars:  leftVars,
		RightVars: rightVars,
		AllVars:   allVars,
		VarTypes:  varTypes,
	}

	// Compute hash (with caching if hasher supports it)
	var hash string
	var err error

	if bsr.hasher != nil {
		hash, err = bsr.hasher.ComputeHashCached(sig)
	} else {
		hash, err = bsr.computeHashDirect(sig)
	}

	if err != nil {
		return nil, "", false, fmt.Errorf("failed to compute hash: %w", err)
	}

	// Check if node already exists
	bsr.mutex.RLock()
	existingNode, exists := bsr.sharedJoinNodes[hash]
	bsr.mutex.RUnlock()

	if exists {
		// Node found - record metrics
		if bsr.config.EnableMetrics {
			RecordSharedReuse(bsr.metrics)
		}

		return existingNode, hash, true, nil
	}

	// Node not found - create new one
	bsr.mutex.Lock()
	defer bsr.mutex.Unlock()

	// Double-check after acquiring write lock
	if existingNode, exists := bsr.sharedJoinNodes[hash]; exists {
		if bsr.config.EnableMetrics {
			RecordSharedReuse(bsr.metrics)
		}
		return existingNode, hash, true, nil
	}

	// Create new JoinNode (convert condition to map[string]interface{})
	var conditionMap map[string]interface{}
	if condition != nil {
		if cm, ok := condition.(map[string]interface{}); ok {
			conditionMap = cm
		} else {
			conditionMap = nil
		}
	}
	node := NewJoinNode(hash, conditionMap, leftVars, rightVars, varTypes, storage)
	bsr.sharedJoinNodes[hash] = node
	bsr.hashToNodeID[hash] = hash // Track hash to node ID mapping

	// Record metrics
	if bsr.config.EnableMetrics {
		RecordUniqueCreation(bsr.metrics)
	}

	return node, hash, false, nil
}

// RegisterJoinNode explicitly registers an existing JoinNode.
func (bsr *BetaSharingRegistryImpl) RegisterJoinNode(node *JoinNode, hash string) error {
	if !bsr.config.Enabled {
		return fmt.Errorf("beta sharing is disabled")
	}

	bsr.mutex.Lock()
	defer bsr.mutex.Unlock()

	// Check if hash is already registered
	if existingNode, exists := bsr.sharedJoinNodes[hash]; exists {
		if existingNode != node {
			return fmt.Errorf("hash %s already registered with different node", hash)
		}
		// Same node already exists
		return nil
	}

	// Register new node
	bsr.sharedJoinNodes[hash] = node

	return nil
}

// AddRuleToJoinNode associates a rule with a join node.
func (bsr *BetaSharingRegistryImpl) AddRuleToJoinNode(nodeID, ruleID string) error {
	if !bsr.config.Enabled {
		return nil
	}

	bsr.mutex.Lock()
	defer bsr.mutex.Unlock()

	if _, exists := bsr.joinNodeRules[nodeID]; !exists {
		bsr.joinNodeRules[nodeID] = make(map[string]bool)
	}

	bsr.joinNodeRules[nodeID][ruleID] = true

	// Also register with lifecycle manager if available
	if bsr.lifecycleManager != nil {
		bsr.lifecycleManager.AddRuleToNode(nodeID, ruleID, ruleID)
	}

	return nil
}

// RemoveRuleFromJoinNode removes a rule's reference from a join node.
// Returns true if the node has no more rules and can be deleted.
func (bsr *BetaSharingRegistryImpl) RemoveRuleFromJoinNode(nodeID, ruleID string) (bool, error) {
	if !bsr.config.Enabled {
		return false, nil
	}

	bsr.mutex.Lock()
	defer bsr.mutex.Unlock()

	rules, exists := bsr.joinNodeRules[nodeID]
	if !exists {
		return false, fmt.Errorf("join node %s not found in rule tracking", nodeID)
	}

	delete(rules, ruleID)

	// Also update lifecycle manager if available
	if bsr.lifecycleManager != nil {
		bsr.lifecycleManager.RemoveRuleFromNode(nodeID, ruleID)
	}

	// If no more rules reference this node, it can be deleted
	canDelete := len(rules) == 0
	if canDelete {
		delete(bsr.joinNodeRules, nodeID)
	}

	return canDelete, nil
}

// RegisterRuleForJoinNode registers a rule as using a specific join node.
// This ensures proper tracking for lifecycle management and reference counting.
func (bsr *BetaSharingRegistryImpl) RegisterRuleForJoinNode(nodeID, ruleID string) error {
	if !bsr.config.Enabled {
		return nil // No-op when sharing is disabled
	}

	bsr.mutex.Lock()
	defer bsr.mutex.Unlock()

	// Initialize rule map for this node if needed
	if _, exists := bsr.joinNodeRules[nodeID]; !exists {
		bsr.joinNodeRules[nodeID] = make(map[string]bool)
	}

	// Add the rule reference
	bsr.joinNodeRules[nodeID][ruleID] = true

	// Sync with lifecycle manager if available
	if bsr.lifecycleManager != nil {
		// Lifecycle manager already tracks this via AddRuleToNode in beta_chain_builder
		// This is just for consistency check
	}

	return nil
}

// UnregisterJoinNode completely removes a join node from the registry.
// Should only be called when the node is being deleted from the network.
func (bsr *BetaSharingRegistryImpl) UnregisterJoinNode(nodeID string) error {
	if !bsr.config.Enabled {
		return nil
	}

	bsr.mutex.Lock()
	defer bsr.mutex.Unlock()

	// Remove from join node rules tracking
	delete(bsr.joinNodeRules, nodeID)

	// Remove from shared nodes if present
	delete(bsr.sharedJoinNodes, nodeID)

	// Find and remove from hash mapping
	for hash, id := range bsr.hashToNodeID {
		if id == nodeID {
			delete(bsr.hashToNodeID, hash)
			break
		}
	}

	return nil
}

// GetJoinNodeRules returns all rules using a specific join node.
func (bsr *BetaSharingRegistryImpl) GetJoinNodeRules(nodeID string) []string {
	bsr.mutex.RLock()
	defer bsr.mutex.RUnlock()

	rules, exists := bsr.joinNodeRules[nodeID]
	if !exists {
		return []string{}
	}

	ruleList := make([]string, 0, len(rules))
	for ruleID := range rules {
		ruleList = append(ruleList, ruleID)
	}
	return ruleList
}

// GetJoinNodeRefCount returns the number of rules referencing a join node.
func (bsr *BetaSharingRegistryImpl) GetJoinNodeRefCount(nodeID string) int {
	bsr.mutex.RLock()
	defer bsr.mutex.RUnlock()

	rules, exists := bsr.joinNodeRules[nodeID]
	if !exists {
		return 0
	}
	return len(rules)
}

// ReleaseJoinNode decrements refcount and removes node if unused.
func (bsr *BetaSharingRegistryImpl) ReleaseJoinNode(hash string) error {
	if !bsr.config.Enabled {
		return nil // No-op when disabled
	}

	bsr.mutex.Lock()
	defer bsr.mutex.Unlock()

	nodeID, exists := bsr.hashToNodeID[hash]
	if !exists {
		return fmt.Errorf("join node with hash %s not found", hash)
	}

	// Check if node still has rules referencing it
	if rules, exists := bsr.joinNodeRules[nodeID]; exists && len(rules) > 0 {
		return fmt.Errorf("cannot release join node %s: still referenced by %d rule(s)", nodeID, len(rules))
	}

	// Remove node from all tracking structures
	delete(bsr.sharedJoinNodes, hash)
	delete(bsr.hashToNodeID, hash)
	delete(bsr.joinNodeRules, nodeID)

	return nil
}

// ReleaseJoinNodeByID removes a join node by its node ID.
// Returns true if the node was found and removed.
func (bsr *BetaSharingRegistryImpl) ReleaseJoinNodeByID(nodeID string) (bool, error) {
	if !bsr.config.Enabled {
		return false, nil
	}

	bsr.mutex.Lock()
	defer bsr.mutex.Unlock()

	// Check if node still has rules referencing it
	if rules, exists := bsr.joinNodeRules[nodeID]; exists && len(rules) > 0 {
		return false, fmt.Errorf("cannot release join node %s: still referenced by %d rule(s)", nodeID, len(rules))
	}

	// Find and remove the hash mapping
	var foundHash string
	for hash, nid := range bsr.hashToNodeID {
		if nid == nodeID {
			foundHash = hash
			break
		}
	}

	if foundHash != "" {
		delete(bsr.sharedJoinNodes, foundHash)
		delete(bsr.hashToNodeID, foundHash)
	}

	delete(bsr.joinNodeRules, nodeID)

	return foundHash != "", nil
}

// GetSharingStats returns current sharing metrics.
func (bsr *BetaSharingRegistryImpl) GetSharingStats() *BetaSharingStats {
	bsr.mutex.RLock()
	totalShared := len(bsr.sharedJoinNodes)
	bsr.mutex.RUnlock()

	stats := &BetaSharingStats{
		TotalSharedNodes:  totalShared,
		TotalRequests:     atomic.LoadInt64(&bsr.metrics.TotalJoinNodesRequested),
		SharedReuses:      atomic.LoadInt64(&bsr.metrics.SharedJoinNodesReused),
		UniqueCreations:   atomic.LoadInt64(&bsr.metrics.UniqueJoinNodesCreated),
		SharingRatio:      CalculateSharingRatio(bsr.metrics),
		HashCacheHitRate:  CalculateCacheHitRate(bsr.metrics),
		AverageHashTimeMs: float64(bsr.metrics.AverageHashTimeNs()) / 1_000_000.0,
		Timestamp:         time.Now(),
	}

	return stats
}

// ListSharedJoinNodes returns all shared join node hashes.
func (bsr *BetaSharingRegistryImpl) ListSharedJoinNodes() []string {
	bsr.mutex.RLock()
	defer bsr.mutex.RUnlock()

	hashes := make([]string, 0, len(bsr.sharedJoinNodes))
	for hash := range bsr.sharedJoinNodes {
		hashes = append(hashes, hash)
	}

	sort.Strings(hashes)
	return hashes
}

// GetSharedJoinNodeDetails returns detailed info about a shared node.
func (bsr *BetaSharingRegistryImpl) GetSharedJoinNodeDetails(hash string) (*JoinNodeDetails, error) {
	bsr.mutex.RLock()
	node, exists := bsr.sharedJoinNodes[hash]
	bsr.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("join node with hash %s not found", hash)
	}

	details := &JoinNodeDetails{
		Hash:             hash,
		NodeID:           node.GetID(),
		ReferenceCount:   1, // Simplified - no refcounting in this version
		LeftVars:         node.LeftVariables,
		RightVars:        node.RightVariables,
		AllVars:          node.AllVariables,
		VarTypes:         node.VariableTypes,
		LeftMemorySize:   len(node.LeftMemory.Tokens),
		RightMemorySize:  len(node.RightMemory.Tokens),
		ResultMemorySize: len(node.ResultMemory.Tokens),
		CreatedAt:        time.Time{}, // TODO: Track creation time
		LastAccessedAt:   time.Now(),
		ActivationCount:  0, // TODO: Track activation count
	}

	return details, nil
}

// ClearCache clears the hash cache.
func (bsr *BetaSharingRegistryImpl) ClearCache() {
	if bsr.hashCache != nil {
		bsr.hashCache.Clear()
	}
}

// Shutdown performs cleanup and releases all resources.
func (bsr *BetaSharingRegistryImpl) Shutdown() error {
	bsr.mutex.Lock()
	defer bsr.mutex.Unlock()

	// Clear all nodes
	bsr.sharedJoinNodes = make(map[string]*JoinNode)

	// Clear cache
	if bsr.hashCache != nil {
		bsr.hashCache.Clear()
	}

	return nil
}

// computeHashDirect computes hash without a hasher (fallback).
func (bsr *BetaSharingRegistryImpl) computeHashDirect(sig *JoinNodeSignature) (string, error) {
	startTime := time.Now()

	// Normalize signature
	canonical, err := bsr.normalizeSignatureFallback(sig)
	if err != nil {
		return "", err
	}

	// Compute hash
	hash, err := ComputeJoinNodeHash(canonical)
	if err != nil {
		return "", err
	}

	// Record metrics
	if bsr.config.EnableMetrics {
		duration := time.Since(startTime).Nanoseconds()
		RecordHashComputation(bsr.metrics, duration)
	}

	return hash, nil
}

// normalizeSignatureFallback provides basic normalization when no normalizer is configured.
func (bsr *BetaSharingRegistryImpl) normalizeSignatureFallback(sig *JoinNodeSignature) (*CanonicalJoinSignature, error) {
	canonical := &CanonicalJoinSignature{
		Version:   "1.0",
		Condition: sig.Condition,
	}

	// Sort variables if normalization is enabled
	if bsr.config.NormalizeOrder {
		canonical.LeftVars = sortStrings(sig.LeftVars)
		canonical.RightVars = sortStrings(sig.RightVars)
		canonical.AllVars = sortStrings(sig.AllVars)
	} else {
		canonical.LeftVars = sig.LeftVars
		canonical.RightVars = sig.RightVars
		canonical.AllVars = sig.AllVars
	}

	// Convert and sort variable types
	canonical.VarTypes = sortVarTypes(sig.VarTypes)

	return canonical, nil
}

// ============================================================================
// DefaultJoinNodeNormalizer Implementation
// ============================================================================

// defaultJoinNodeNormalizer is the default implementation of JoinNodeNormalizer.
type defaultJoinNodeNormalizer struct {
	config BetaSharingConfig
}

// NewDefaultJoinNodeNormalizer creates a default normalizer.
func NewDefaultJoinNodeNormalizer(config BetaSharingConfig) JoinNodeNormalizer {
	return &defaultJoinNodeNormalizer{config: config}
}

// NormalizeSignature converts a join signature to canonical form.
func (n *defaultJoinNodeNormalizer) NormalizeSignature(sig *JoinNodeSignature) (*CanonicalJoinSignature, error) {
	canonical := &CanonicalJoinSignature{
		Version: "1.0",
	}

	// Sort variables if normalization is enabled
	if n.config.NormalizeOrder {
		canonical.LeftVars = sortStrings(sig.LeftVars)
		canonical.RightVars = sortStrings(sig.RightVars)
		canonical.AllVars = sortStrings(sig.AllVars)
	} else {
		canonical.LeftVars = sig.LeftVars
		canonical.RightVars = sig.RightVars
		canonical.AllVars = sig.AllVars
	}

	// Convert and sort variable types
	canonical.VarTypes = sortVarTypes(sig.VarTypes)

	// Normalize condition
	normalizedCondition, err := n.NormalizeCondition(sig.Condition)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize condition: %w", err)
	}
	canonical.Condition = normalizedCondition

	return canonical, nil
}

// NormalizeCondition converts a condition AST to canonical form.
func (n *defaultJoinNodeNormalizer) NormalizeCondition(condition interface{}) (interface{}, error) {
	if condition == nil {
		return nil, nil
	}

	// Handle map[string]interface{} condition
	if condMap, ok := condition.(map[string]interface{}); ok {
		return n.normalizeConditionMap(condMap)
	}

	// Return as-is for other types
	return condition, nil
}

// normalizeConditionMap normalizes a condition represented as a map.
func (n *defaultJoinNodeNormalizer) normalizeConditionMap(condition map[string]interface{}) (map[string]interface{}, error) {
	normalized := make(map[string]interface{})

	// Copy all fields
	for key, value := range condition {
		normalized[key] = value
	}

	// Normalize operator synonyms
	if opType, ok := condition["type"].(string); ok {
		if opType == "comparison" {
			normalized["type"] = "binaryOperation"
		}
	}

	// Normalize commutative operators (== and !=)
	if n.config.EnableAdvancedNormalization {
		if operator, ok := condition["operator"].(string); ok {
			if operator == "==" || operator == "!=" {
				// Sort operands for commutative operators
				if left, hasLeft := condition["left"]; hasLeft {
					if right, hasRight := condition["right"]; hasRight {
						leftJSON, _ := json.Marshal(left)
						rightJSON, _ := json.Marshal(right)

						// Swap if left > right lexicographically
						if string(leftJSON) > string(rightJSON) {
							normalized["left"] = right
							normalized["right"] = left
						}
					}
				}
			}
		}
	}

	return normalized, nil
}

// ============================================================================
// DefaultJoinNodeHasher Implementation
// ============================================================================

// defaultJoinNodeHasher is the default implementation of JoinNodeHasher.
type defaultJoinNodeHasher struct {
	config     BetaSharingConfig
	normalizer JoinNodeNormalizer
	cache      *LRUCache
}

// NewDefaultJoinNodeHasher creates a default hasher with LRU caching.
func NewDefaultJoinNodeHasher(config BetaSharingConfig) JoinNodeHasher {
	return &defaultJoinNodeHasher{
		config:     config,
		normalizer: NewDefaultJoinNodeNormalizer(config),
		cache:      NewLRUCache(config.HashCacheSize, 0), // 0 = no TTL
	}
}

// ComputeHash computes a hash for a canonical join signature.
func (h *defaultJoinNodeHasher) ComputeHash(canonical *CanonicalJoinSignature) (string, error) {
	// Serialize to JSON
	jsonBytes, err := json.Marshal(canonical)
	if err != nil {
		return "", fmt.Errorf("failed to marshal canonical signature: %w", err)
	}

	// Compute SHA-256 hash
	hashBytes := sha256.Sum256(jsonBytes)
	hashHex := hex.EncodeToString(hashBytes[:8])

	return "join_" + hashHex, nil
}

// ComputeHashCached computes a hash with LRU caching.
func (h *defaultJoinNodeHasher) ComputeHashCached(sig *JoinNodeSignature) (string, error) {
	// Create cache key from raw signature
	cacheKeyBytes, err := json.Marshal(sig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal signature for cache key: %w", err)
	}
	cacheKey := string(cacheKeyBytes)

	// Check cache
	if cachedValue, found := h.cache.Get(cacheKey); found {
		if cachedHash, ok := cachedValue.(string); ok {
			return cachedHash, nil
		}
	}

	// Cache miss - normalize and compute hash
	canonical, err := h.normalizer.NormalizeSignature(sig)
	if err != nil {
		return "", err
	}

	hash, err := h.ComputeHash(canonical)
	if err != nil {
		return "", err
	}

	// Store in cache
	h.cache.Set(cacheKey, hash)

	return hash, nil
}

// ============================================================================
// Helper Functions
// ============================================================================

// NormalizeJoinCondition normalizes a join condition for consistent hashing.
// This is a standalone helper function for backward compatibility.
func NormalizeJoinCondition(condition map[string]interface{}) (map[string]interface{}, error) {
	if condition == nil {
		return map[string]interface{}{"type": "simple"}, nil
	}

	normalized := make(map[string]interface{})

	// Normalize type field
	if condType, hasType := condition["type"]; hasType {
		if typeStr, ok := condType.(string); ok {
			switch typeStr {
			case "comparison", "binaryOperation":
				normalized["type"] = "binaryOperation"
			default:
				normalized["type"] = typeStr
			}
		} else {
			normalized["type"] = condType
		}
	}

	// Handle commutative operators
	if operator, hasOp := condition["operator"].(string); hasOp {
		normalized["operator"] = operator

		if operator == "==" || operator == "!=" {
			// For commutative operators, sort left/right by lexicographic order
			left, hasLeft := condition["left"]
			right, hasRight := condition["right"]

			if hasLeft && hasRight {
				leftJSON, _ := json.Marshal(left)
				rightJSON, _ := json.Marshal(right)

				if string(leftJSON) > string(rightJSON) {
					normalized["left"] = right
					normalized["right"] = left
				} else {
					normalized["left"] = left
					normalized["right"] = right
				}
			}
		} else {
			// Non-commutative operators: preserve order
			if left, hasLeft := condition["left"]; hasLeft {
				normalized["left"] = left
			}
			if right, hasRight := condition["right"]; hasRight {
				normalized["right"] = right
			}
		}
	}

	// Copy remaining fields
	for key, value := range condition {
		if key != "type" && key != "operator" && key != "left" && key != "right" {
			normalized[key] = value
		}
	}

	return normalized, nil
}

// ComputeJoinHash computes a hash for a join specification.
// This is a standalone helper function for backward compatibility.
func ComputeJoinHash(condition map[string]interface{}, leftVars, rightVars []string, varTypes map[string]string) (string, error) {
	// Create canonical signature
	canonical := &CanonicalJoinSignature{
		Version:   "1.0",
		LeftVars:  sortStrings(leftVars),
		RightVars: sortStrings(rightVars),
		AllVars:   sortStrings(append(leftVars, rightVars...)),
		VarTypes:  sortVarTypes(varTypes),
		Condition: condition,
	}

	// Normalize condition
	if condition != nil {
		normalizedCond, err := NormalizeJoinCondition(condition)
		if err != nil {
			return "", err
		}
		canonical.Condition = normalizedCond
	}

	return ComputeJoinNodeHash(canonical)
}
