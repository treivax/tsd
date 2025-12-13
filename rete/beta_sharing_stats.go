// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sort"
	"sync/atomic"
	"time"
)

// ============================================================================
// Statistics and Introspection (BetaSharingRegistryImpl methods)
// ============================================================================

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
		CreatedAt:        node.CreatedAt(),
		LastAccessedAt:   node.GetLastActivatedAt(),
		ActivationCount:  node.GetActivationCount(),
	}

	return details, nil
}
