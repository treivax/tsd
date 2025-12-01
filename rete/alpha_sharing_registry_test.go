// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"sync"
	"testing"
)

// TestAlphaSharingRegistry_NewRegistry tests the creation of a new registry
func TestAlphaSharingRegistry_NewRegistry(t *testing.T) {
	t.Run("default constructor", func(t *testing.T) {
		registry := NewAlphaSharingRegistry()
		if registry == nil {
			t.Fatal("NewAlphaSharingRegistry returned nil")
		}
		if registry.sharedAlphaNodes == nil {
			t.Error("sharedAlphaNodes map not initialized")
		}
		if registry.config == nil {
			t.Error("config not initialized")
		}
		if registry.metrics == nil {
			t.Error("metrics not initialized")
		}
	})

	t.Run("with custom metrics", func(t *testing.T) {
		metrics := NewChainBuildMetrics()
		registry := NewAlphaSharingRegistryWithMetrics(metrics)
		if registry == nil {
			t.Fatal("NewAlphaSharingRegistryWithMetrics returned nil")
		}
		if registry.metrics != metrics {
			t.Error("metrics not set correctly")
		}
	})

	t.Run("with custom config", func(t *testing.T) {
		config := DefaultChainPerformanceConfig()
		config.HashCacheEnabled = false
		metrics := NewChainBuildMetrics()
		registry := NewAlphaSharingRegistryWithConfig(config, metrics)
		if registry == nil {
			t.Fatal("NewAlphaSharingRegistryWithConfig returned nil")
		}
		if registry.config != config {
			t.Error("config not set correctly")
		}
	})
}

// TestAlphaSharingRegistry_GetOrCreateAlphaNode_Sharing tests node sharing behavior
func TestAlphaSharingRegistry_GetOrCreateAlphaNode_Sharing(t *testing.T) {
	storage := NewMemoryStorage()
	registry := NewAlphaSharingRegistry()

	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left":     map[string]interface{}{"type": "fieldAccess", "field": "age"},
		"right":    map[string]interface{}{"type": "literal", "value": 18},
	}

	t.Run("first creation", func(t *testing.T) {
		node1, hash1, wasShared1, err := registry.GetOrCreateAlphaNode(condition, "p", storage)
		if err != nil {
			t.Fatalf("GetOrCreateAlphaNode failed: %v", err)
		}
		if node1 == nil {
			t.Fatal("node1 is nil")
		}
		if hash1 == "" {
			t.Error("hash1 is empty")
		}
		if wasShared1 {
			t.Error("first node should not be marked as shared")
		}
		if node1.ID != hash1 {
			t.Errorf("node ID (%s) should equal hash (%s)", node1.ID, hash1)
		}
	})

	t.Run("second creation with same condition - should reuse", func(t *testing.T) {
		node2, _, wasShared2, err := registry.GetOrCreateAlphaNode(condition, "p", storage)
		if err != nil {
			t.Fatalf("GetOrCreateAlphaNode failed: %v", err)
		}
		if node2 == nil {
			t.Fatal("node2 is nil")
		}
		if !wasShared2 {
			t.Error("second node should be marked as shared")
		}

		// Verify it's the same node
		node1, _, _, _ := registry.GetOrCreateAlphaNode(condition, "p", storage)
		if node1.ID != node2.ID {
			t.Errorf("nodes should be identical: %s vs %s", node1.ID, node2.ID)
		}
	})

	t.Run("different condition - should create new node", func(t *testing.T) {
		differentCondition := map[string]interface{}{
			"type":     "comparison",
			"operator": "<",
			"left":     map[string]interface{}{"type": "fieldAccess", "field": "age"},
			"right":    map[string]interface{}{"type": "literal", "value": 65},
		}
		node3, hash3, wasShared3, err := registry.GetOrCreateAlphaNode(differentCondition, "p", storage)
		if err != nil {
			t.Fatalf("GetOrCreateAlphaNode failed: %v", err)
		}
		if wasShared3 {
			t.Error("different condition should create new node")
		}

		node1, hash1, _, _ := registry.GetOrCreateAlphaNode(condition, "p", storage)
		if hash3 == hash1 {
			t.Error("different conditions should have different hashes")
		}
		if node3.ID == node1.ID {
			t.Error("different conditions should create different nodes")
		}
	})

	t.Run("different variable name - should create new node", func(t *testing.T) {
		_, hash4, wasShared4, err := registry.GetOrCreateAlphaNode(condition, "q", storage)
		if err != nil {
			t.Fatalf("GetOrCreateAlphaNode failed: %v", err)
		}
		if wasShared4 {
			t.Error("different variable name should create new node")
		}

		_, hash1, _, _ := registry.GetOrCreateAlphaNode(condition, "p", storage)
		if hash4 == hash1 {
			t.Error("different variable names should have different hashes")
		}
	})
}

// TestAlphaSharingRegistry_GetStats tests statistics collection
func TestAlphaSharingRegistry_GetStats(t *testing.T) {
	storage := NewMemoryStorage()
	registry := NewAlphaSharingRegistry()

	t.Run("empty registry", func(t *testing.T) {
		stats := registry.GetStats()
		if stats == nil {
			t.Fatal("GetStats returned nil")
		}

		totalNodes := stats["total_shared_alpha_nodes"].(int)
		if totalNodes != 0 {
			t.Errorf("total_shared_alpha_nodes = %d, want 0", totalNodes)
		}

		totalRefs := stats["total_rule_references"].(int)
		if totalRefs != 0 {
			t.Errorf("total_rule_references = %d, want 0", totalRefs)
		}

		avgRatio := stats["average_sharing_ratio"].(float64)
		if avgRatio != 0.0 {
			t.Errorf("average_sharing_ratio = %f, want 0.0", avgRatio)
		}
	})

	t.Run("single node with one child", func(t *testing.T) {
		registry.Reset()
		cond := map[string]interface{}{
			"type":     "comparison",
			"operator": ">",
			"left":     "age",
			"right":    18,
		}
		node, _, _, _ := registry.GetOrCreateAlphaNode(cond, "p", storage)

		// Add one child (simulate one rule using this node)
		child := NewAlphaNode("child1", nil, "p", storage)
		node.AddChild(child)

		stats := registry.GetStats()
		totalNodes := stats["total_shared_alpha_nodes"].(int)
		if totalNodes != 1 {
			t.Errorf("total_shared_alpha_nodes = %d, want 1", totalNodes)
		}

		totalRefs := stats["total_rule_references"].(int)
		if totalRefs != 1 {
			t.Errorf("total_rule_references = %d, want 1", totalRefs)
		}

		avgRatio := stats["average_sharing_ratio"].(float64)
		if avgRatio != 1.0 {
			t.Errorf("average_sharing_ratio = %f, want 1.0", avgRatio)
		}
	})

	t.Run("single node with multiple children", func(t *testing.T) {
		registry.Reset()
		cond := map[string]interface{}{
			"type":     "comparison",
			"operator": ">",
			"left":     "age",
			"right":    18,
		}
		node, _, _, _ := registry.GetOrCreateAlphaNode(cond, "p", storage)

		// Add three children (simulate three rules sharing this node)
		for i := 0; i < 3; i++ {
			child := NewAlphaNode("child"+string(rune(i)), nil, "p", storage)
			node.AddChild(child)
		}

		stats := registry.GetStats()
		totalNodes := stats["total_shared_alpha_nodes"].(int)
		if totalNodes != 1 {
			t.Errorf("total_shared_alpha_nodes = %d, want 1", totalNodes)
		}

		totalRefs := stats["total_rule_references"].(int)
		if totalRefs != 3 {
			t.Errorf("total_rule_references = %d, want 3", totalRefs)
		}

		avgRatio := stats["average_sharing_ratio"].(float64)
		if avgRatio != 3.0 {
			t.Errorf("average_sharing_ratio = %f, want 3.0", avgRatio)
		}
	})

	t.Run("multiple nodes with different child counts", func(t *testing.T) {
		registry.Reset()

		// Node 1: 2 children
		cond1 := map[string]interface{}{"type": "comparison", "operator": ">", "left": "age", "right": 18}
		node1, _, _, _ := registry.GetOrCreateAlphaNode(cond1, "p", storage)
		for i := 0; i < 2; i++ {
			node1.AddChild(NewAlphaNode("node1_child"+string(rune(i)), nil, "p", storage))
		}

		// Node 2: 3 children
		cond2 := map[string]interface{}{"type": "comparison", "operator": "<", "left": "age", "right": 65}
		node2, _, _, _ := registry.GetOrCreateAlphaNode(cond2, "p", storage)
		for i := 0; i < 3; i++ {
			node2.AddChild(NewAlphaNode("node2_child"+string(rune(i)), nil, "p", storage))
		}

		// Node 3: 1 child
		cond3 := map[string]interface{}{"type": "comparison", "operator": "==", "left": "name", "right": "Alice"}
		node3, _, _, _ := registry.GetOrCreateAlphaNode(cond3, "p", storage)
		node3.AddChild(NewAlphaNode("node3_child", nil, "p", storage))

		stats := registry.GetStats()
		totalNodes := stats["total_shared_alpha_nodes"].(int)
		if totalNodes != 3 {
			t.Errorf("total_shared_alpha_nodes = %d, want 3", totalNodes)
		}

		totalRefs := stats["total_rule_references"].(int)
		expected := 2 + 3 + 1
		if totalRefs != expected {
			t.Errorf("total_rule_references = %d, want %d", totalRefs, expected)
		}

		avgRatio := stats["average_sharing_ratio"].(float64)
		expectedAvg := float64(expected) / 3.0
		if avgRatio != expectedAvg {
			t.Errorf("average_sharing_ratio = %f, want %f", avgRatio, expectedAvg)
		}
	})
}

// TestAlphaSharingRegistry_RemoveNode tests node removal
func TestAlphaSharingRegistry_RemoveNode(t *testing.T) {
	storage := NewMemoryStorage()
	registry := NewAlphaSharingRegistry()

	cond := map[string]interface{}{"type": "comparison", "operator": ">", "left": "age", "right": 18}
	node, hash, _, _ := registry.GetOrCreateAlphaNode(cond, "p", storage)

	t.Run("remove existing node", func(t *testing.T) {
		err := registry.RemoveAlphaNode(hash)
		if err != nil {
			t.Errorf("RemoveAlphaNode failed: %v", err)
		}

		// Verify node is gone
		_, exists := registry.GetAlphaNode(hash)
		if exists {
			t.Error("node should not exist after removal")
		}

		stats := registry.GetStats()
		totalNodes := stats["total_shared_alpha_nodes"].(int)
		if totalNodes != 0 {
			t.Errorf("total_shared_alpha_nodes = %d, want 0 after removal", totalNodes)
		}
	})

	t.Run("remove non-existent node", func(t *testing.T) {
		err := registry.RemoveAlphaNode("nonexistent_hash")
		if err == nil {
			t.Error("RemoveAlphaNode should return error for non-existent node")
		}
	})

	t.Run("remove then recreate", func(t *testing.T) {
		node2, _, wasShared, _ := registry.GetOrCreateAlphaNode(cond, "p", storage)
		if wasShared {
			t.Error("recreated node should not be marked as shared")
		}
		if node2.ID != node.ID {
			t.Error("recreated node should have same hash-based ID")
		}
	})
}

// TestAlphaSharingRegistry_ListSharedAlphaNodes tests listing functionality
func TestAlphaSharingRegistry_ListSharedAlphaNodes(t *testing.T) {
	storage := NewMemoryStorage()
	registry := NewAlphaSharingRegistry()

	t.Run("empty registry", func(t *testing.T) {
		list := registry.ListSharedAlphaNodes()
		if len(list) != 0 {
			t.Errorf("ListSharedAlphaNodes should return empty list, got %d nodes", len(list))
		}
	})

	t.Run("multiple nodes", func(t *testing.T) {
		cond1 := map[string]interface{}{"type": "comparison", "operator": ">", "left": "age", "right": 18}
		cond2 := map[string]interface{}{"type": "comparison", "operator": "<", "left": "age", "right": 65}
		cond3 := map[string]interface{}{"type": "comparison", "operator": "==", "left": "name", "right": "Alice"}

		_, hash1, _, _ := registry.GetOrCreateAlphaNode(cond1, "p", storage)
		_, hash2, _, _ := registry.GetOrCreateAlphaNode(cond2, "p", storage)
		_, hash3, _, _ := registry.GetOrCreateAlphaNode(cond3, "p", storage)

		list := registry.ListSharedAlphaNodes()
		if len(list) != 3 {
			t.Errorf("ListSharedAlphaNodes should return 3 nodes, got %d", len(list))
		}

		// Verify all hashes are in the list
		hashSet := make(map[string]bool)
		for _, h := range list {
			hashSet[h] = true
		}
		if !hashSet[hash1] || !hashSet[hash2] || !hashSet[hash3] {
			t.Error("not all hashes found in list")
		}

		// Verify list is sorted
		for i := 1; i < len(list); i++ {
			if list[i-1] > list[i] {
				t.Error("list should be sorted")
				break
			}
		}
	})
}

// TestAlphaSharingRegistry_ResetRegistry tests registry reset
func TestAlphaSharingRegistry_ResetRegistry(t *testing.T) {
	storage := NewMemoryStorage()
	registry := NewAlphaSharingRegistry()

	// Create some nodes
	cond1 := map[string]interface{}{"type": "comparison", "operator": ">", "left": "age", "right": 18}
	cond2 := map[string]interface{}{"type": "comparison", "operator": "<", "left": "age", "right": 65}
	registry.GetOrCreateAlphaNode(cond1, "p", storage)
	registry.GetOrCreateAlphaNode(cond2, "p", storage)

	stats := registry.GetStats()
	totalBefore := stats["total_shared_alpha_nodes"].(int)
	if totalBefore != 2 {
		t.Errorf("should have 2 nodes before reset, got %d", totalBefore)
	}

	// Reset
	registry.Reset()

	// Verify everything is cleared
	stats = registry.GetStats()
	totalAfter := stats["total_shared_alpha_nodes"].(int)
	if totalAfter != 0 {
		t.Errorf("should have 0 nodes after reset, got %d", totalAfter)
	}

	list := registry.ListSharedAlphaNodes()
	if len(list) != 0 {
		t.Errorf("ListSharedAlphaNodes should return empty list after reset, got %d", len(list))
	}

	cacheSize := registry.GetHashCacheSize()
	if cacheSize != 0 {
		t.Errorf("hash cache should be empty after reset, got size %d", cacheSize)
	}
}

// TestAlphaSharingRegistry_GetSharedAlphaNodeDetails tests detailed info retrieval
func TestAlphaSharingRegistry_GetSharedAlphaNodeDetails(t *testing.T) {
	storage := NewMemoryStorage()
	registry := NewAlphaSharingRegistry()

	cond := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left":     "age",
		"right":    18,
	}

	node, hash, _, _ := registry.GetOrCreateAlphaNode(cond, "p", storage)

	// Add children
	child1 := NewTerminalNode("term1", nil, storage)
	child2 := NewTerminalNode("term2", nil, storage)
	node.AddChild(child1)
	node.AddChild(child2)

	t.Run("existing node", func(t *testing.T) {
		details := registry.GetSharedAlphaNodeDetails(hash)
		if details == nil {
			t.Fatal("GetSharedAlphaNodeDetails returned nil")
		}

		if details["hash"] != hash {
			t.Errorf("hash = %v, want %v", details["hash"], hash)
		}

		if details["node_id"] != node.ID {
			t.Errorf("node_id = %v, want %v", details["node_id"], node.ID)
		}

		if details["variable_name"] != "p" {
			t.Errorf("variable_name = %v, want 'p'", details["variable_name"])
		}

		childCount := details["child_count"].(int)
		if childCount != 2 {
			t.Errorf("child_count = %d, want 2", childCount)
		}

		childIDs := details["child_ids"].([]string)
		if len(childIDs) != 2 {
			t.Errorf("child_ids length = %d, want 2", len(childIDs))
		}
	})

	t.Run("non-existent node", func(t *testing.T) {
		details := registry.GetSharedAlphaNodeDetails("nonexistent")
		if details != nil {
			t.Error("GetSharedAlphaNodeDetails should return nil for non-existent node")
		}
	})
}

// TestAlphaSharingRegistry_ThreadSafety tests thread safety
func TestAlphaSharingRegistry_ThreadSafety(t *testing.T) {
	storage := NewMemoryStorage()
	registry := NewAlphaSharingRegistry()

	cond := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left":     "age",
		"right":    18,
	}

	// Run concurrent GetOrCreateAlphaNode calls
	numGoroutines := 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			node, hash, _, err := registry.GetOrCreateAlphaNode(cond, "p", storage)
			if err != nil {
				t.Errorf("GetOrCreateAlphaNode failed: %v", err)
			}
			if node == nil {
				t.Error("node is nil")
			}
			if hash == "" {
				t.Error("hash is empty")
			}
		}()
	}

	wg.Wait()

	// Should only have 1 shared node (all goroutines shared the same node)
	stats := registry.GetStats()
	totalNodes := stats["total_shared_alpha_nodes"].(int)
	if totalNodes != 1 {
		t.Errorf("should have only 1 shared node after concurrent access, got %d", totalNodes)
	}
}

// TestAlphaSharingRegistry_HashCache tests hash caching functionality
func TestAlphaSharingRegistry_HashCache(t *testing.T) {
	t.Run("with cache enabled", func(t *testing.T) {
		config := DefaultChainPerformanceConfig()
		config.HashCacheEnabled = true
		config.HashCacheMaxSize = 100
		metrics := NewChainBuildMetrics()
		registry := NewAlphaSharingRegistryWithConfig(config, metrics)

		cond := map[string]interface{}{"type": "comparison", "operator": ">", "left": "age", "right": 18}

		// First call - cache miss
		hash1, err := registry.ConditionHashCached(cond, "p")
		if err != nil {
			t.Fatalf("ConditionHashCached failed: %v", err)
		}

		// Second call - cache hit
		hash2, err := registry.ConditionHashCached(cond, "p")
		if err != nil {
			t.Fatalf("ConditionHashCached failed: %v", err)
		}

		if hash1 != hash2 {
			t.Error("hashes should be identical")
		}

		cacheSize := registry.GetHashCacheSize()
		if cacheSize == 0 {
			t.Error("cache should not be empty")
		}
	})

	t.Run("with cache disabled", func(t *testing.T) {
		config := DefaultChainPerformanceConfig()
		config.HashCacheEnabled = false
		metrics := NewChainBuildMetrics()
		registry := NewAlphaSharingRegistryWithConfig(config, metrics)

		cond := map[string]interface{}{"type": "comparison", "operator": ">", "left": "age", "right": 18}

		hash1, err := registry.ConditionHashCached(cond, "p")
		if err != nil {
			t.Fatalf("ConditionHashCached failed: %v", err)
		}

		hash2, err := registry.ConditionHashCached(cond, "p")
		if err != nil {
			t.Fatalf("ConditionHashCached failed: %v", err)
		}

		if hash1 != hash2 {
			t.Error("hashes should still be identical without cache")
		}

		// Cache should remain at 0 since caching is disabled
		cacheSize := registry.GetHashCacheSize()
		if cacheSize != 0 {
			t.Error("cache should be empty when disabled")
		}
	})

	t.Run("clear cache", func(t *testing.T) {
		config := DefaultChainPerformanceConfig()
		config.HashCacheEnabled = true
		metrics := NewChainBuildMetrics()
		registry := NewAlphaSharingRegistryWithConfig(config, metrics)

		cond := map[string]interface{}{"type": "comparison", "operator": ">", "left": "age", "right": 18}
		registry.ConditionHashCached(cond, "p")

		sizeBefore := registry.GetHashCacheSize()
		if sizeBefore == 0 {
			t.Error("cache should not be empty before clear")
		}

		registry.ClearHashCache()

		sizeAfter := registry.GetHashCacheSize()
		if sizeAfter != 0 {
			t.Errorf("cache should be empty after clear, got size %d", sizeAfter)
		}
	})
}

// TestAlphaSharingRegistry_ConditionNormalization tests condition normalization for sharing
func TestAlphaSharingRegistry_ConditionNormalization(t *testing.T) {
	storage := NewMemoryStorage()
	registry := NewAlphaSharingRegistry()

	t.Run("wrapped conditions should be normalized", func(t *testing.T) {
		// Direct condition
		directCond := map[string]interface{}{
			"type":     "comparison",
			"operator": ">",
			"left":     "age",
			"right":    18,
		}

		// Wrapped condition (as used in some parts of the system)
		wrappedCond := map[string]interface{}{
			"type": "constraint",
			"constraint": map[string]interface{}{
				"type":     "comparison",
				"operator": ">",
				"left":     "age",
				"right":    18,
			},
		}

		node1, hash1, _, _ := registry.GetOrCreateAlphaNode(directCond, "p", storage)
		node2, hash2, wasShared, _ := registry.GetOrCreateAlphaNode(wrappedCond, "p", storage)

		// Should share the same node after normalization
		if !wasShared {
			t.Error("wrapped condition should reuse direct condition's node after normalization")
		}
		if hash1 != hash2 {
			t.Errorf("hashes should be equal after normalization: %s vs %s", hash1, hash2)
		}
		if node1.ID != node2.ID {
			t.Error("should return same node after normalization")
		}
	})

	t.Run("comparison should normalize to binaryOperation", func(t *testing.T) {
		registry.Reset()

		cond1 := map[string]interface{}{
			"type":     "comparison",
			"operator": ">",
			"left":     "age",
			"right":    18,
		}

		cond2 := map[string]interface{}{
			"type":     "binaryOperation",
			"operator": ">",
			"left":     "age",
			"right":    18,
		}

		node1, hash1, _, _ := registry.GetOrCreateAlphaNode(cond1, "p", storage)
		node2, hash2, wasShared, _ := registry.GetOrCreateAlphaNode(cond2, "p", storage)

		// Should share after type normalization
		if !wasShared {
			t.Error("binaryOperation should reuse comparison node after normalization")
		}
		if hash1 != hash2 {
			t.Errorf("hashes should be equal after type normalization: %s vs %s", hash1, hash2)
		}
		if node1.ID != node2.ID {
			t.Error("should return same node after type normalization")
		}
	})
}
