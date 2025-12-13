// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
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
	cascadeLevel int,
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

	// Create signature with cascade level to prevent incorrect sharing between cascade levels
	sig := &JoinNodeSignature{
		Condition:    condition,
		LeftVars:     leftVars,
		RightVars:    rightVars,
		AllVars:      allVars,
		VarTypes:     varTypes,
		CascadeLevel: cascadeLevel,
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

// Note: Statistics and introspection functions have been extracted to beta_sharing_stats.go

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

// Note: Hash computation, normalization, and helper functions have been
// extracted to beta_sharing_hash.go and beta_sharing_helpers.go
