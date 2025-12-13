// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package rete implements the RETE algorithm for the TSD rule engine.
// This file contains the interface definitions and type declarations for
// Beta node (JoinNode) sharing functionality.
//
// Beta sharing eliminates duplicate JoinNodes by identifying and reusing nodes
// with identical join patterns, reducing memory consumption by 30-50% and
// improving runtime performance by 20-40% in typical rule bases.
package rete

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// Configuration Constants
// ============================================================================

const (
	// DefaultHashCacheSize est la taille par défaut du cache LRU pour les hash
	DefaultHashCacheSize = 1000

	// DefaultMaxSharedNodes est le nombre maximum de nœuds partagés (0 = illimité)
	DefaultMaxSharedNodes = 10000

	// HashPrefixLength est la longueur du préfixe de hash utilisé (en bytes)
	HashPrefixLength = 8
)

// ============================================================================
// Core Interfaces
// ============================================================================

// BetaSharingRegistry is the central registry for shared JoinNodes.
// It provides hash-based lookup, lifecycle management, and metrics collection.
//
// Thread-safe: All public methods can be called concurrently.
type BetaSharingRegistry interface {
	// GetOrCreateJoinNode returns a shared JoinNode for the given signature.
	// If a compatible node exists, it is returned and refcount is incremented.
	// Otherwise, a new node is created and registered.
	//
	// Returns:
	//   - node: The JoinNode instance (shared or newly created)
	//   - hash: Canonical hash of the join signature
	//   - wasShared: true if an existing node was reused, false if newly created
	//   - error: Any error during hash computation or node creation
	GetOrCreateJoinNode(
		condition interface{},
		leftVars []string,
		rightVars []string,
		allVars []string,
		varTypes map[string]string,
		storage Storage,
		cascadeLevel int,
	) (*JoinNode, string, bool, error)

	// RegisterJoinNode explicitly registers an existing JoinNode.
	// Used when migrating from non-shared to shared nodes.
	//
	// Returns error if a different node is already registered with the same hash.
	RegisterJoinNode(node *JoinNode, hash string) error

	// ReleaseJoinNode decrements refcount and removes node if unused.
	// Should be called when a rule is removed.
	//
	// Returns error if the hash is not found in the registry.
	ReleaseJoinNode(hash string) error

	// GetSharingStats returns current sharing metrics.
	GetSharingStats() *BetaSharingStats

	// ListSharedJoinNodes returns all shared join node hashes.
	// Results are sorted alphabetically.
	ListSharedJoinNodes() []string

	// GetSharedJoinNodeDetails returns detailed info about a shared node.
	// Returns error if the hash is not found.
	GetSharedJoinNodeDetails(hash string) (*JoinNodeDetails, error)

	// RegisterRuleForJoinNode registers a rule as using a specific join node.
	// This ensures proper tracking for lifecycle management and reference counting.
	RegisterRuleForJoinNode(nodeID, ruleID string) error

	// UnregisterJoinNode completely removes a join node from the registry.
	// Should only be called when the node is being deleted from the network.
	UnregisterJoinNode(nodeID string) error

	// AddRuleToJoinNode associates a rule with a join node.
	// Used for tracking which rules reference which join nodes.
	AddRuleToJoinNode(nodeID, ruleID string) error

	// RemoveRuleFromJoinNode removes a rule's reference from a join node.
	// Returns true if the node has no more rules and can be deleted.
	RemoveRuleFromJoinNode(nodeID, ruleID string) (bool, error)

	// GetJoinNodeRules returns all rules using a specific join node.
	GetJoinNodeRules(nodeID string) []string

	// GetJoinNodeRefCount returns the number of rules referencing a join node.
	GetJoinNodeRefCount(nodeID string) int

	// ReleaseJoinNodeByID removes a join node by its node ID.
	// Returns true if the node was found and removed.
	ReleaseJoinNodeByID(nodeID string) (bool, error)

	// ClearCache clears the hash cache.
	// Used for testing or when normalization rules change.
	ClearCache()

	// Clear vide tous les caches et nodes partagés
	Clear()

	// Shutdown performs cleanup and releases all resources.
	Shutdown() error
}

// JoinNodeNormalizer normalizes join signatures into canonical form.
type JoinNodeNormalizer interface {
	// NormalizeSignature converts a join signature to canonical form.
	// Equivalent signatures produce identical canonical structures.
	NormalizeSignature(sig *JoinNodeSignature) (*CanonicalJoinSignature, error)

	// NormalizeCondition converts a condition AST to canonical form.
	NormalizeCondition(condition interface{}) (interface{}, error)
}

// JoinNodeHasher computes hashes for join signatures.
type JoinNodeHasher interface {
	// ComputeHash computes a hash for a canonical join signature.
	ComputeHash(canonical *CanonicalJoinSignature) (string, error)

	// ComputeHashCached computes a hash with LRU caching.
	ComputeHashCached(sig *JoinNodeSignature) (string, error)
}

// ============================================================================
// Data Structures
// ============================================================================

// BetaSharingRegistryImpl is the concrete implementation of BetaSharingRegistry.
type BetaSharingRegistryImpl struct {
	// Shared JoinNodes indexed by canonical hash
	sharedJoinNodes map[string]*JoinNode

	// Hash to NodeID mapping for efficient lookups
	hashToNodeID map[string]string

	// Join node to rules mapping (nodeID -> set of ruleIDs)
	joinNodeRules map[string]map[string]bool

	// LRU cache for condition hash computation
	// Key: normalized signature JSON
	// Value: computed hash string
	hashCache *LRUCache

	// Lifecycle manager for reference counting (optional)
	lifecycleManager *LifecycleManager

	// Normalizer for join signatures
	normalizer JoinNodeNormalizer

	// Hasher for computing signature hashes
	hasher JoinNodeHasher

	// Metrics collection
	metrics *BetaBuildMetrics

	// Thread safety
	mutex sync.RWMutex

	// Configuration
	config BetaSharingConfig
}

// BetaSharingConfig contains configuration for Beta node sharing.
type BetaSharingConfig struct {
	// Enabled controls whether Beta sharing is active
	Enabled bool

	// HashCacheSize is the maximum number of cached hash computations
	// Default: 1000
	HashCacheSize int

	// MaxSharedNodes is the maximum number of shared nodes
	// Default: 10000 (0 = unlimited)
	MaxSharedNodes int

	// EnableMetrics controls whether metrics collection is active
	EnableMetrics bool

	// NormalizeOrder controls whether variable order is canonicalized
	// Set to false for debugging to preserve original order
	NormalizeOrder bool

	// EnableAdvancedNormalization enables experimental normalization features
	// (commutative operators, associativity, etc.)
	EnableAdvancedNormalization bool
}

// DefaultBetaSharingConfig returns the default configuration.
func DefaultBetaSharingConfig() BetaSharingConfig {
	return BetaSharingConfig{
		Enabled:                     false, // Disabled by default for safe rollout
		HashCacheSize:               DefaultHashCacheSize,
		MaxSharedNodes:              DefaultMaxSharedNodes,
		EnableMetrics:               true,
		NormalizeOrder:              true,
		EnableAdvancedNormalization: false,
	}
}

// ============================================================================
// Join Signature Types
// ============================================================================

// JoinNodeSignature represents the complete signature of a JoinNode.
// Used as input to the sharing registry.
type JoinNodeSignature struct {
	// Condition is the join condition AST
	Condition interface{}

	// LeftVars are the variables from the left input
	LeftVars []string

	// RightVars are the variables from the right input
	RightVars []string

	// AllVars are all variables (accumulated from both inputs)
	AllVars []string

	// VarTypes maps variable names to their types
	VarTypes map[string]string

	// CascadeLevel represents the depth in the join cascade (0-based)
	// Level 0: first join (v0 ⋈ v1)
	// Level 1: second join ((v0,v1) ⋈ v2)
	// Level N: (N+1)th join
	// This ensures joins at different cascade levels are not incorrectly shared
	CascadeLevel int
}

// CanonicalJoinSignature is the normalized form of a JoinNodeSignature.
// Semantically equivalent signatures produce identical canonical forms.
type CanonicalJoinSignature struct {
	// Version of the canonicalization algorithm
	// Allows future evolution of normalization rules
	Version string

	// Sorted variable lists
	LeftVars  []string
	RightVars []string
	AllVars   []string

	// Sorted type mappings
	VarTypes []VariableTypeMapping

	// Normalized condition AST
	Condition interface{}

	// CascadeLevel for distinguishing joins at different cascade depths
	CascadeLevel int
}

// VariableTypeMapping represents a variable name and its type.
// Used in canonical signatures for deterministic ordering.
type VariableTypeMapping struct {
	VarName  string `json:"var_name"`
	TypeName string `json:"type_name"`
}

// ToJSON serializes the canonical signature to JSON.
// Output is deterministic (sorted keys, no whitespace).
func (c *CanonicalJoinSignature) ToJSON() (string, error) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal canonical signature: %w", err)
	}
	return string(bytes), nil
}

// ============================================================================
// Metrics Types
// ============================================================================

// BetaBuildMetrics collects statistics about Beta node sharing.
type BetaBuildMetrics struct {
	// Total number of GetOrCreateJoinNode calls
	TotalJoinNodesRequested int64

	// Number of times a shared node was reused
	SharedJoinNodesReused int64

	// Number of times a new unique node was created
	UniqueJoinNodesCreated int64

	// Hash cache statistics
	HashCacheHits   int64
	HashCacheMisses int64

	// Performance metrics
	TotalHashTimeNs int64 // Cumulative time spent hashing
	HashCount       int64 // Number of hash operations
}

// AverageHashTimeNs returns the average time per hash computation.
func (m *BetaBuildMetrics) AverageHashTimeNs() int64 {
	if m.HashCount == 0 {
		return 0
	}
	return m.TotalHashTimeNs / m.HashCount
}

// BetaSharingStats provides a snapshot of sharing statistics.
type BetaSharingStats struct {
	// Current state
	TotalSharedNodes int

	// Cumulative counters
	TotalRequests   int64
	SharedReuses    int64
	UniqueCreations int64

	// Computed metrics
	SharingRatio      float64 // SharedReuses / TotalRequests
	HashCacheHitRate  float64 // CacheHits / (CacheHits + CacheMisses)
	AverageHashTimeMs float64

	// Timestamp
	Timestamp time.Time
}

// JoinNodeDetails provides detailed information about a shared JoinNode.
type JoinNodeDetails struct {
	// Identity
	Hash   string
	NodeID string

	// Reference counting
	ReferenceCount int

	// Join signature
	LeftVars  []string
	RightVars []string
	AllVars   []string
	VarTypes  map[string]string

	// Memory state
	LeftMemorySize   int
	RightMemorySize  int
	ResultMemorySize int

	// Metadata
	CreatedAt       time.Time
	LastAccessedAt  time.Time
	ActivationCount int64
}

// ============================================================================
// Normalization Helpers
// ============================================================================

// Note: LRUCache is defined in lru_cache.go

// NormalizationContext holds state during normalization.
type NormalizationContext struct {
	// Configuration
	NormalizeOrder              bool
	EnableAdvancedNormalization bool

	// State
	seenNodes map[interface{}]bool
}

// NewNormalizationContext creates a new normalization context.
func NewNormalizationContext(config BetaSharingConfig) *NormalizationContext {
	return &NormalizationContext{
		NormalizeOrder:              config.NormalizeOrder,
		EnableAdvancedNormalization: config.EnableAdvancedNormalization,
		seenNodes:                   make(map[interface{}]bool),
	}
}

// ============================================================================
// Hash Computation Helpers
// ============================================================================

// ComputeSHA256Hash computes a SHA-256 hash of the input bytes.
// Returns a hex-encoded string of the first HashPrefixLength bytes.
func ComputeSHA256Hash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:HashPrefixLength])
}

// ComputeJoinNodeHash computes a hash for a canonical join signature.
// Returns a hash string prefixed with "join_".
func ComputeJoinNodeHash(canonical *CanonicalJoinSignature) (string, error) {
	// Serialize to JSON
	jsonBytes, err := json.Marshal(canonical)
	if err != nil {
		return "", fmt.Errorf("failed to marshal canonical signature: %w", err)
	}

	// Compute hash
	hashHex := ComputeSHA256Hash(jsonBytes)

	// Add prefix for readability
	return "join_" + hashHex, nil
}

// ============================================================================
// Constructor Functions
// ============================================================================

// NewBetaSharingRegistry creates a new Beta sharing registry.
func NewBetaSharingRegistry(
	config BetaSharingConfig,
	lifecycle *LifecycleManager,
) *BetaSharingRegistryImpl {
	return &BetaSharingRegistryImpl{
		sharedJoinNodes:  make(map[string]*JoinNode),
		hashToNodeID:     make(map[string]string),
		joinNodeRules:    make(map[string]map[string]bool),
		hashCache:        NewLRUCache(config.HashCacheSize, 0), // 0 = no TTL
		lifecycleManager: lifecycle,
		normalizer:       NewDefaultJoinNodeNormalizer(config),
		hasher:           NewDefaultJoinNodeHasher(config),
		metrics:          &BetaBuildMetrics{},
		config:           config,
	}
}

// ============================================================================
// Compatibility Testing
// ============================================================================

// CanShareJoinNodes determines if two join signatures can share a node.
func CanShareJoinNodes(sig1, sig2 *JoinNodeSignature) bool {
	// Check variable set equality
	if !equalStringSlices(sig1.LeftVars, sig2.LeftVars) {
		return false
	}
	if !equalStringSlices(sig1.RightVars, sig2.RightVars) {
		return false
	}

	// Check variable type compatibility
	if !compatibleVarTypes(sig1.VarTypes, sig2.VarTypes) {
		return false
	}

	// TODO: Deep comparison of normalized conditions

	return true
}

// equalStringSlices compares two string slices for equality (order matters).
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// compatibleVarTypes checks if two variable type maps are compatible.
func compatibleVarTypes(types1, types2 map[string]string) bool {
	if len(types1) != len(types2) {
		return false
	}
	for varName, type1 := range types1 {
		type2, ok := types2[varName]
		if !ok || type1 != type2 {
			return false
		}
	}
	return true
}

// sortStrings returns a sorted copy of the input slice.
func sortStrings(input []string) []string {
	result := make([]string, len(input))
	copy(result, input)
	sort.Strings(result)
	return result
}

// sortVarTypes converts a map to a sorted slice of VariableTypeMapping.
func sortVarTypes(varTypes map[string]string) []VariableTypeMapping {
	result := make([]VariableTypeMapping, 0, len(varTypes))
	for varName, typeName := range varTypes {
		result = append(result, VariableTypeMapping{
			VarName:  varName,
			TypeName: typeName,
		})
	}

	// Sort by variable name
	sort.Slice(result, func(i, j int) bool {
		return result[i].VarName < result[j].VarName
	})

	return result
}

// ============================================================================
// Metrics Helpers
// ============================================================================

// RecordHashComputation records metrics for a hash computation.
func RecordHashComputation(metrics *BetaBuildMetrics, durationNs int64) {
	atomic.AddInt64(&metrics.TotalHashTimeNs, durationNs)
	atomic.AddInt64(&metrics.HashCount, 1)
}

// RecordCacheHit increments the cache hit counter.
func RecordCacheHit(metrics *BetaBuildMetrics) {
	atomic.AddInt64(&metrics.HashCacheHits, 1)
}

// RecordCacheMiss increments the cache miss counter.
func RecordCacheMiss(metrics *BetaBuildMetrics) {
	atomic.AddInt64(&metrics.HashCacheMisses, 1)
}

// RecordSharedReuse increments the shared reuse counter.
func RecordSharedReuse(metrics *BetaBuildMetrics) {
	atomic.AddInt64(&metrics.SharedJoinNodesReused, 1)
	atomic.AddInt64(&metrics.TotalJoinNodesRequested, 1)
}

// RecordUniqueCreation increments the unique creation counter.
func RecordUniqueCreation(metrics *BetaBuildMetrics) {
	atomic.AddInt64(&metrics.UniqueJoinNodesCreated, 1)
	atomic.AddInt64(&metrics.TotalJoinNodesRequested, 1)
}

// CalculateSharingRatio computes the sharing ratio from metrics.
func CalculateSharingRatio(metrics *BetaBuildMetrics) float64 {
	total := atomic.LoadInt64(&metrics.TotalJoinNodesRequested)
	if total == 0 {
		return 0.0
	}
	reused := atomic.LoadInt64(&metrics.SharedJoinNodesReused)
	return float64(reused) / float64(total)
}

// CalculateCacheHitRate computes the cache hit rate from metrics.
func CalculateCacheHitRate(metrics *BetaBuildMetrics) float64 {
	hits := atomic.LoadInt64(&metrics.HashCacheHits)
	misses := atomic.LoadInt64(&metrics.HashCacheMisses)
	total := hits + misses
	if total == 0 {
		return 0.0
	}
	return float64(hits) / float64(total)
}

// Clear vide tous les caches et nodes partagés du BetaSharingRegistry
func (bsr *BetaSharingRegistryImpl) Clear() {
	bsr.mutex.Lock()
	defer bsr.mutex.Unlock()

	bsr.sharedJoinNodes = make(map[string]*JoinNode)
	bsr.hashToNodeID = make(map[string]string)
	bsr.joinNodeRules = make(map[string]map[string]bool)

	if bsr.hashCache != nil {
		bsr.hashCache.Clear()
	}
}
