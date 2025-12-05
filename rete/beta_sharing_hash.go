// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// ============================================================================
// Hash Computation and Normalization (BetaSharingRegistryImpl methods)
// ============================================================================

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
