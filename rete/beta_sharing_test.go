// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// TestBetaSharingRegistry_CreateWithDefaultConfig tests creating a registry with default config
func TestBetaSharingRegistry_CreateWithDefaultConfig(t *testing.T) {
	config := DefaultBetaSharingConfig()
	registry := NewBetaSharingRegistry(config, nil)

	if registry == nil {
		t.Fatal("Registry should not be nil")
	}

	if registry.config.Enabled {
		t.Error("Beta sharing should be disabled by default")
	}

	if registry.config.HashCacheSize != 1000 {
		t.Errorf("Expected HashCacheSize 1000, got %d", registry.config.HashCacheSize)
	}
}

// TestBetaSharingRegistry_GetOrCreateJoinNode_Disabled tests behavior when sharing is disabled
func TestBetaSharingRegistry_GetOrCreateJoinNode_Disabled(t *testing.T) {
	config := DefaultBetaSharingConfig()
	config.Enabled = false
	registry := NewBetaSharingRegistry(config, nil)
	storage := NewMemoryStorage()

	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "p",
			"field":  "id",
		},
		"right": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "o",
			"field":  "customer_id",
		},
	}

	leftVars := []string{"p"}
	rightVars := []string{"o"}
	allVars := []string{"p", "o"}
	varTypes := map[string]string{"p": "Person", "o": "Order"}

	// First call
	node1, hash1, reused1, err := registry.GetOrCreateJoinNode(
		condition, leftVars, rightVars, allVars, varTypes, storage,
	)
	if err != nil {
		t.Fatalf("Error creating node1: %v", err)
	}

	if reused1 {
		t.Error("Node should not be marked as reused when sharing is disabled")
	}

	// Second call - should create a new node (not shared)
	node2, hash2, reused2, err := registry.GetOrCreateJoinNode(
		condition, leftVars, rightVars, allVars, varTypes, storage,
	)
	if err != nil {
		t.Fatalf("Error creating node2: %v", err)
	}

	if reused2 {
		t.Error("Node should not be reused when sharing is disabled")
	}

	if node1 == node2 {
		t.Error("Nodes should be different when sharing is disabled")
	}

	if hash1 == hash2 {
		t.Log("Note: Hashes may be the same (timestamp-based) but nodes are different")
	}
}

// TestBetaSharingRegistry_GetOrCreateJoinNode_SameCondition tests sharing with identical conditions
func TestBetaSharingRegistry_GetOrCreateJoinNode_SameCondition(t *testing.T) {
	config := DefaultBetaSharingConfig()
	config.Enabled = true
	registry := NewBetaSharingRegistry(config, nil)
	storage := NewMemoryStorage()

	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "p",
			"field":  "id",
		},
		"right": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "o",
			"field":  "customer_id",
		},
	}

	leftVars := []string{"p"}
	rightVars := []string{"o"}
	allVars := []string{"p", "o"}
	varTypes := map[string]string{"p": "Person", "o": "Order"}

	// First call
	node1, hash1, reused1, err := registry.GetOrCreateJoinNode(
		condition, leftVars, rightVars, allVars, varTypes, storage,
	)
	if err != nil {
		t.Fatalf("Error creating node1: %v", err)
	}

	if reused1 {
		t.Error("First node should not be marked as reused")
	}

	// Second call - should reuse the same node
	node2, hash2, reused2, err := registry.GetOrCreateJoinNode(
		condition, leftVars, rightVars, allVars, varTypes, storage,
	)
	if err != nil {
		t.Fatalf("Error creating node2: %v", err)
	}

	if !reused2 {
		t.Error("Second node should be marked as reused")
	}

	if hash1 != hash2 {
		t.Errorf("Hashes should be identical: %s vs %s", hash1, hash2)
	}

	if node1 != node2 {
		t.Error("Nodes should be the same object (shared)")
	}

	t.Logf("✅ JoinNode shared correctly (hash: %s)", hash1)
}

// TestBetaSharingRegistry_GetOrCreateJoinNode_DifferentConditions tests different conditions create different nodes
func TestBetaSharingRegistry_GetOrCreateJoinNode_DifferentConditions(t *testing.T) {
	config := DefaultBetaSharingConfig()
	config.Enabled = true
	registry := NewBetaSharingRegistry(config, nil)
	storage := NewMemoryStorage()

	condition1 := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "p",
			"field":  "id",
		},
		"right": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "o",
			"field":  "customer_id",
		},
	}

	condition2 := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "p",
			"field":  "age",
		},
		"right": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "o",
			"field":  "min_age",
		},
	}

	leftVars := []string{"p"}
	rightVars := []string{"o"}
	allVars := []string{"p", "o"}
	varTypes := map[string]string{"p": "Person", "o": "Order"}

	node1, hash1, _, err := registry.GetOrCreateJoinNode(
		condition1, leftVars, rightVars, allVars, varTypes, storage,
	)
	if err != nil {
		t.Fatalf("Error creating node1: %v", err)
	}

	node2, hash2, _, err := registry.GetOrCreateJoinNode(
		condition2, leftVars, rightVars, allVars, varTypes, storage,
	)
	if err != nil {
		t.Fatalf("Error creating node2: %v", err)
	}

	if hash1 == hash2 {
		t.Errorf("Different conditions should have different hashes: %s", hash1)
	}

	if node1 == node2 {
		t.Error("Different conditions should create different nodes")
	}

	t.Logf("✅ Different conditions create separate nodes (hash1: %s, hash2: %s)", hash1, hash2)
}

// TestBetaSharingRegistry_RegisterJoinNode tests explicit node registration
func TestBetaSharingRegistry_RegisterJoinNode(t *testing.T) {
	config := DefaultBetaSharingConfig()
	config.Enabled = true
	registry := NewBetaSharingRegistry(config, nil)
	storage := NewMemoryStorage()

	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
	}

	node := NewJoinNode("test_node", condition, []string{"p"}, []string{"o"}, map[string]string{"p": "Person"}, storage)

	err := registry.RegisterJoinNode(node, "test_hash")
	if err != nil {
		t.Fatalf("Error registering node: %v", err)
	}

	// Try to retrieve the node
	nodes := registry.ListSharedJoinNodes()
	if len(nodes) != 1 {
		t.Errorf("Expected 1 shared node, got %d", len(nodes))
	}

	if nodes[0] != "test_hash" {
		t.Errorf("Expected hash 'test_hash', got '%s'", nodes[0])
	}

	t.Logf("✅ Node registered successfully")
}

// TestBetaSharingRegistry_ReleaseJoinNode tests node release and cleanup
func TestBetaSharingRegistry_ReleaseJoinNode(t *testing.T) {
	config := DefaultBetaSharingConfig()
	config.Enabled = true
	registry := NewBetaSharingRegistry(config, nil)
	storage := NewMemoryStorage()

	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
	}

	leftVars := []string{"p"}
	rightVars := []string{"o"}
	allVars := []string{"p", "o"}
	varTypes := map[string]string{"p": "Person"}

	// Create node
	_, hash, _, err := registry.GetOrCreateJoinNode(
		condition, leftVars, rightVars, allVars, varTypes, storage,
	)
	if err != nil {
		t.Fatalf("Error creating node: %v", err)
	}

	// Verify node exists
	nodes := registry.ListSharedJoinNodes()
	if len(nodes) != 1 {
		t.Errorf("Expected 1 node before release, got %d", len(nodes))
	}

	// Release node
	err = registry.ReleaseJoinNode(hash)
	if err != nil {
		t.Fatalf("Error releasing node: %v", err)
	}

	// Verify node is removed
	nodes = registry.ListSharedJoinNodes()
	if len(nodes) != 0 {
		t.Errorf("Expected 0 nodes after release, got %d", len(nodes))
	}

	t.Logf("✅ Node released and removed successfully")
}

// TestBetaSharingRegistry_GetSharingStats tests statistics collection
func TestBetaSharingRegistry_GetSharingStats(t *testing.T) {
	config := DefaultBetaSharingConfig()
	config.Enabled = true
	config.EnableMetrics = true
	registry := NewBetaSharingRegistry(config, nil)
	storage := NewMemoryStorage()

	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
	}

	leftVars := []string{"p"}
	rightVars := []string{"o"}
	allVars := []string{"p", "o"}
	varTypes := map[string]string{"p": "Person"}

	// Create first node
	registry.GetOrCreateJoinNode(condition, leftVars, rightVars, allVars, varTypes, storage)

	// Reuse same node
	registry.GetOrCreateJoinNode(condition, leftVars, rightVars, allVars, varTypes, storage)

	// Get stats
	stats := registry.GetSharingStats()

	if stats.TotalRequests != 2 {
		t.Errorf("Expected 2 total requests, got %d", stats.TotalRequests)
	}

	if stats.SharedReuses != 1 {
		t.Errorf("Expected 1 shared reuse, got %d", stats.SharedReuses)
	}

	if stats.UniqueCreations != 1 {
		t.Errorf("Expected 1 unique creation, got %d", stats.UniqueCreations)
	}

	if stats.SharingRatio != 0.5 {
		t.Errorf("Expected sharing ratio 0.5, got %f", stats.SharingRatio)
	}

	t.Logf("✅ Stats: %d requests, %d reuses, %d unique, ratio: %.2f",
		stats.TotalRequests, stats.SharedReuses, stats.UniqueCreations, stats.SharingRatio)
}

// TestBetaSharingRegistry_ListSharedJoinNodes tests listing shared nodes
func TestBetaSharingRegistry_ListSharedJoinNodes(t *testing.T) {
	config := DefaultBetaSharingConfig()
	config.Enabled = true
	registry := NewBetaSharingRegistry(config, nil)
	storage := NewMemoryStorage()

	condition1 := map[string]interface{}{"type": "comparison", "operator": "=="}
	condition2 := map[string]interface{}{"type": "comparison", "operator": ">"}

	leftVars := []string{"p"}
	rightVars := []string{"o"}
	allVars := []string{"p", "o"}
	varTypes := map[string]string{"p": "Person"}

	registry.GetOrCreateJoinNode(condition1, leftVars, rightVars, allVars, varTypes, storage)
	registry.GetOrCreateJoinNode(condition2, leftVars, rightVars, allVars, varTypes, storage)

	nodes := registry.ListSharedJoinNodes()

	if len(nodes) != 2 {
		t.Errorf("Expected 2 shared nodes, got %d", len(nodes))
	}

	// Verify list is sorted
	if len(nodes) == 2 && nodes[0] > nodes[1] {
		t.Error("List should be sorted alphabetically")
	}

	t.Logf("✅ Listed %d shared nodes: %v", len(nodes), nodes)
}

// TestBetaSharingRegistry_GetSharedJoinNodeDetails tests retrieving node details
func TestBetaSharingRegistry_GetSharedJoinNodeDetails(t *testing.T) {
	config := DefaultBetaSharingConfig()
	config.Enabled = true
	registry := NewBetaSharingRegistry(config, nil)
	storage := NewMemoryStorage()

	condition := map[string]interface{}{"type": "comparison", "operator": "=="}
	leftVars := []string{"p"}
	rightVars := []string{"o"}
	allVars := []string{"p", "o"}
	varTypes := map[string]string{"p": "Person", "o": "Order"}

	_, hash, _, err := registry.GetOrCreateJoinNode(
		condition, leftVars, rightVars, allVars, varTypes, storage,
	)
	if err != nil {
		t.Fatalf("Error creating node: %v", err)
	}

	details, err := registry.GetSharedJoinNodeDetails(hash)
	if err != nil {
		t.Fatalf("Error getting details: %v", err)
	}

	if details.Hash != hash {
		t.Errorf("Expected hash %s, got %s", hash, details.Hash)
	}

	if len(details.LeftVars) != 1 || details.LeftVars[0] != "p" {
		t.Errorf("Expected LeftVars ['p'], got %v", details.LeftVars)
	}

	if len(details.RightVars) != 1 || details.RightVars[0] != "o" {
		t.Errorf("Expected RightVars ['o'], got %v", details.RightVars)
	}

	t.Logf("✅ Node details retrieved: %+v", details)
}

// TestBetaSharingRegistry_ClearCache tests cache clearing
func TestBetaSharingRegistry_ClearCache(t *testing.T) {
	config := DefaultBetaSharingConfig()
	config.Enabled = true
	registry := NewBetaSharingRegistry(config, nil)

	registry.ClearCache()

	// Verify cache is empty (size should be 0)
	if registry.hashCache != nil && registry.hashCache.Len() != 0 {
		t.Errorf("Cache should be empty after clear, got size %d", registry.hashCache.Len())
	}

	t.Logf("✅ Cache cleared successfully")
}

// TestBetaSharingRegistry_Shutdown tests shutdown cleanup
func TestBetaSharingRegistry_Shutdown(t *testing.T) {
	config := DefaultBetaSharingConfig()
	config.Enabled = true
	registry := NewBetaSharingRegistry(config, nil)
	storage := NewMemoryStorage()

	condition := map[string]interface{}{"type": "comparison"}
	leftVars := []string{"p"}
	rightVars := []string{"o"}
	allVars := []string{"p", "o"}
	varTypes := map[string]string{"p": "Person"}

	registry.GetOrCreateJoinNode(condition, leftVars, rightVars, allVars, varTypes, storage)

	err := registry.Shutdown()
	if err != nil {
		t.Fatalf("Error during shutdown: %v", err)
	}

	// Verify all nodes are cleared
	nodes := registry.ListSharedJoinNodes()
	if len(nodes) != 0 {
		t.Errorf("Expected 0 nodes after shutdown, got %d", len(nodes))
	}

	t.Logf("✅ Shutdown completed successfully")
}

// TestNormalizeJoinCondition tests condition normalization
func TestNormalizeJoinCondition(t *testing.T) {
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
		"left":     map[string]interface{}{"field": "id"},
		"right":    map[string]interface{}{"field": "customer_id"},
	}

	normalized, err := NormalizeJoinCondition(condition)
	if err != nil {
		t.Fatalf("Error normalizing condition: %v", err)
	}

	if normalized["type"] != "binaryOperation" {
		t.Errorf("Expected type 'binaryOperation', got %v", normalized["type"])
	}

	t.Logf("✅ Condition normalized: %+v", normalized)
}

// TestComputeJoinHash tests hash computation
func TestComputeJoinHash(t *testing.T) {
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
	}

	leftVars := []string{"p"}
	rightVars := []string{"o"}
	varTypes := map[string]string{"p": "Person", "o": "Order"}

	hash1, err := ComputeJoinHash(condition, leftVars, rightVars, varTypes)
	if err != nil {
		t.Fatalf("Error computing hash: %v", err)
	}

	// Compute again - should be identical
	hash2, err := ComputeJoinHash(condition, leftVars, rightVars, varTypes)
	if err != nil {
		t.Fatalf("Error computing hash: %v", err)
	}

	if hash1 != hash2 {
		t.Errorf("Hash should be deterministic: %s vs %s", hash1, hash2)
	}

	if hash1[:5] != "join_" {
		t.Errorf("Hash should start with 'join_', got %s", hash1)
	}

	t.Logf("✅ Hash computed consistently: %s", hash1)
}

// TestDefaultJoinNodeNormalizer tests the default normalizer
func TestDefaultJoinNodeNormalizer(t *testing.T) {
	config := DefaultBetaSharingConfig()
	config.NormalizeOrder = true
	normalizer := NewDefaultJoinNodeNormalizer(config)

	sig := &JoinNodeSignature{
		Condition: map[string]interface{}{"type": "comparison"},
		LeftVars:  []string{"b", "a"},
		RightVars: []string{"d", "c"},
		AllVars:   []string{"b", "a", "d", "c"},
		VarTypes:  map[string]string{"a": "TypeA", "b": "TypeB"},
	}

	canonical, err := normalizer.NormalizeSignature(sig)
	if err != nil {
		t.Fatalf("Error normalizing signature: %v", err)
	}

	// Variables should be sorted
	if canonical.LeftVars[0] != "a" || canonical.LeftVars[1] != "b" {
		t.Errorf("LeftVars should be sorted: %v", canonical.LeftVars)
	}

	if canonical.RightVars[0] != "c" || canonical.RightVars[1] != "d" {
		t.Errorf("RightVars should be sorted: %v", canonical.RightVars)
	}

	t.Logf("✅ Signature normalized with sorted variables")
}

// TestDefaultJoinNodeHasher tests the default hasher
func TestDefaultJoinNodeHasher(t *testing.T) {
	config := DefaultBetaSharingConfig()
	hasher := NewDefaultJoinNodeHasher(config)

	sig := &JoinNodeSignature{
		Condition: map[string]interface{}{"type": "comparison"},
		LeftVars:  []string{"p"},
		RightVars: []string{"o"},
		AllVars:   []string{"p", "o"},
		VarTypes:  map[string]string{"p": "Person"},
	}

	// First call (cache miss)
	hash1, err := hasher.ComputeHashCached(sig)
	if err != nil {
		t.Fatalf("Error computing hash: %v", err)
	}

	// Second call (cache hit)
	hash2, err := hasher.ComputeHashCached(sig)
	if err != nil {
		t.Fatalf("Error computing hash: %v", err)
	}

	if hash1 != hash2 {
		t.Errorf("Cached hash should be identical: %s vs %s", hash1, hash2)
	}

	t.Logf("✅ Hasher with cache working correctly: %s", hash1)
}

// TestBetaBuildMetrics_AverageHashTimeNs tests metrics calculation
func TestBetaBuildMetrics_AverageHashTimeNs(t *testing.T) {
	metrics := &BetaBuildMetrics{}

	// No computations yet
	if metrics.AverageHashTimeNs() != 0 {
		t.Error("Average should be 0 when no hashes computed")
	}

	// Add some hash computations
	RecordHashComputation(metrics, 1000)
	RecordHashComputation(metrics, 2000)
	RecordHashComputation(metrics, 3000)

	avg := metrics.AverageHashTimeNs()
	expectedAvg := int64(2000)

	if avg != expectedAvg {
		t.Errorf("Expected average %d, got %d", expectedAvg, avg)
	}

	t.Logf("✅ Average hash time: %d ns", avg)
}

// TestCanShareJoinNodes tests compatibility checking
func TestCanShareJoinNodes(t *testing.T) {
	sig1 := &JoinNodeSignature{
		LeftVars:  []string{"p"},
		RightVars: []string{"o"},
		VarTypes:  map[string]string{"p": "Person", "o": "Order"},
	}

	sig2 := &JoinNodeSignature{
		LeftVars:  []string{"p"},
		RightVars: []string{"o"},
		VarTypes:  map[string]string{"p": "Person", "o": "Order"},
	}

	sig3 := &JoinNodeSignature{
		LeftVars:  []string{"p"},
		RightVars: []string{"q"},
		VarTypes:  map[string]string{"p": "Person", "q": "Quote"},
	}

	if !CanShareJoinNodes(sig1, sig2) {
		t.Error("Identical signatures should be shareable")
	}

	if CanShareJoinNodes(sig1, sig3) {
		t.Error("Different right vars should not be shareable")
	}

	t.Logf("✅ Compatibility checking works correctly")
}

// TestLRUCache tests the LRU cache implementation
func TestLRUCache(t *testing.T) {
	cache := NewLRUCache(2, 0) // Capacity of 2, no TTL

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	// Retrieve key1
	if val, found := cache.Get("key1"); !found || val != "value1" {
		t.Errorf("Expected 'value1', got %v (found: %v)", val, found)
	}

	// Add key3 - should evict key2 (least recently used)
	cache.Set("key3", "value3")

	// key2 should be evicted
	if _, found := cache.Get("key2"); found {
		t.Error("key2 should have been evicted")
	}

	// key1 and key3 should still exist
	if _, found := cache.Get("key1"); !found {
		t.Error("key1 should still exist")
	}

	if _, found := cache.Get("key3"); !found {
		t.Error("key3 should still exist")
	}

	// Test size
	if cache.Len() != 2 {
		t.Errorf("Expected size 2, got %d", cache.Len())
	}

	// Test clear
	cache.Clear()
	if cache.Len() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", cache.Len())
	}

	t.Logf("✅ LRU cache works correctly")
}

// BenchmarkBetaSharingRegistry_GetOrCreateJoinNode benchmarks node creation
func BenchmarkBetaSharingRegistry_GetOrCreateJoinNode(b *testing.B) {
	config := DefaultBetaSharingConfig()
	config.Enabled = true
	registry := NewBetaSharingRegistry(config, nil)
	storage := NewMemoryStorage()

	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
	}

	leftVars := []string{"p"}
	rightVars := []string{"o"}
	allVars := []string{"p", "o"}
	varTypes := map[string]string{"p": "Person"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registry.GetOrCreateJoinNode(condition, leftVars, rightVars, allVars, varTypes, storage)
	}
}

// BenchmarkComputeJoinHash benchmarks hash computation
func BenchmarkComputeJoinHash(b *testing.B) {
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
	}

	leftVars := []string{"p"}
	rightVars := []string{"o"}
	varTypes := map[string]string{"p": "Person"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ComputeJoinHash(condition, leftVars, rightVars, varTypes)
	}
}
