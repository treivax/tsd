// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"fmt"
	"testing"
	"time"
)

// TestNewBetaJoinCache tests cache creation
func TestNewBetaJoinCache(t *testing.T) {
	config := DefaultChainPerformanceConfig()
	cache := NewBetaJoinCache(config)
	if cache == nil {
		t.Fatal("Expected non-nil cache")
	}
	if cache.lruCache == nil {
		t.Error("Expected LRU cache to be initialized")
	}
	if cache.config != config {
		t.Error("Expected config to match")
	}
}

// TestNewBetaJoinCache_NilConfig tests creation with nil config
func TestNewBetaJoinCache_NilConfig(t *testing.T) {
	cache := NewBetaJoinCache(nil)
	if cache == nil {
		t.Fatal("Expected non-nil cache")
	}
	if cache.config == nil {
		t.Error("Expected default config to be set")
	}
}

// TestGetSetJoinResult tests basic cache get/set operations
func TestGetSetJoinResult(t *testing.T) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	config.BetaJoinResultCacheMaxSize = 100
	cache := NewBetaJoinCache(config)
	// Create test data
	leftToken := &Token{
		ID:     "token1",
		Facts:  []*Fact{},
		NodeID: "node1",
	}
	rightFact := &Fact{
		ID:     "fact1",
		Type:   "TestType",
		Fields: map[string]interface{}{"field": "value"},
	}
	joinNode := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	// Initially, should not be in cache
	if result, found := cache.GetJoinResult(leftToken, rightFact, joinNode); found {
		t.Errorf("Expected cache miss, got hit with result: %v", result)
	}
	// Set a result
	expectedResult := &JoinResult{
		Matched:  true,
		Token:    leftToken,
		JoinType: "binary",
	}
	cache.SetJoinResult(leftToken, rightFact, joinNode, expectedResult)
	// Now should be in cache
	result, found := cache.GetJoinResult(leftToken, rightFact, joinNode)
	if !found {
		t.Fatal("Expected cache hit after set")
	}
	if !result.Matched {
		t.Error("Expected matched=true")
	}
	if result.Token != leftToken {
		t.Error("Expected token to match")
	}
	if result.JoinType != "binary" {
		t.Errorf("Expected JoinType 'binary', got '%s'", result.JoinType)
	}
}

// TestGetSetJoinResult_Disabled tests cache when disabled
func TestGetSetJoinResult_Disabled(t *testing.T) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = false
	cache := NewBetaJoinCache(config)
	leftToken := &Token{ID: "token1"}
	rightFact := &Fact{ID: "fact1", Type: "Test"}
	joinNode := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	result := &JoinResult{Matched: true}
	cache.SetJoinResult(leftToken, rightFact, joinNode, result)
	// Should not find anything when disabled
	if _, found := cache.GetJoinResult(leftToken, rightFact, joinNode); found {
		t.Error("Expected no cache hit when cache is disabled")
	}
}

// TestCacheHitMiss tests hit/miss tracking
func TestCacheHitMiss(t *testing.T) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	cache := NewBetaJoinCache(config)
	leftToken := &Token{ID: "token1"}
	rightFact := &Fact{ID: "fact1", Type: "Test"}
	joinNode := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	// First access - miss
	cache.GetJoinResult(leftToken, rightFact, joinNode)
	stats := cache.GetStats()
	if stats["hits"].(int64) != 0 {
		t.Errorf("Expected 0 hits, got %d", stats["hits"])
	}
	if stats["misses"].(int64) != 1 {
		t.Errorf("Expected 1 miss, got %d", stats["misses"])
	}
	// Set and get - hit
	result := &JoinResult{Matched: true}
	cache.SetJoinResult(leftToken, rightFact, joinNode, result)
	cache.GetJoinResult(leftToken, rightFact, joinNode)
	stats = cache.GetStats()
	if stats["hits"].(int64) != 1 {
		t.Errorf("Expected 1 hit, got %d", stats["hits"])
	}
	if stats["misses"].(int64) != 1 {
		t.Errorf("Expected 1 miss, got %d", stats["misses"])
	}
}

// TestBetaJoinCache_GetHitRate tests hit rate calculation
func TestBetaJoinCache_GetHitRate(t *testing.T) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	cache := NewBetaJoinCache(config)
	// Initially 0
	hitRate := cache.GetHitRate()
	if hitRate != 0.0 {
		t.Errorf("Expected initial hit rate 0.0, got %.2f", hitRate)
	}
	leftToken := &Token{ID: "token1"}
	rightFact := &Fact{ID: "fact1", Type: "Test"}
	joinNode := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	// 1 miss
	cache.GetJoinResult(leftToken, rightFact, joinNode)
	hitRate = cache.GetHitRate()
	if hitRate != 0.0 {
		t.Errorf("Expected hit rate 0.0 after miss, got %.2f", hitRate)
	}
	// Set and get (1 hit, 1 miss = 50%)
	result := &JoinResult{Matched: true}
	cache.SetJoinResult(leftToken, rightFact, joinNode, result)
	cache.GetJoinResult(leftToken, rightFact, joinNode)
	hitRate = cache.GetHitRate()
	if hitRate != 0.5 {
		t.Errorf("Expected hit rate 0.5, got %.2f", hitRate)
	}
	// Another hit (2 hits, 1 miss = 66.67%)
	cache.GetJoinResult(leftToken, rightFact, joinNode)
	hitRate = cache.GetHitRate()
	expectedRate := 2.0 / 3.0
	if hitRate < expectedRate-0.01 || hitRate > expectedRate+0.01 {
		t.Errorf("Expected hit rate ~%.2f, got %.2f", expectedRate, hitRate)
	}
}

// TestCacheEviction tests LRU eviction behavior
func TestCacheEviction(t *testing.T) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	config.BetaJoinResultCacheMaxSize = 2 // Small cache for testing
	cache := NewBetaJoinCache(config)
	joinNode := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	// Add 3 entries to a cache of size 2
	for i := 1; i <= 3; i++ {
		token := &Token{ID: fmt.Sprintf("token%d", i)}
		fact := &Fact{ID: fmt.Sprintf("fact%d", i), Type: "Test"}
		result := &JoinResult{Matched: true}
		cache.SetJoinResult(token, fact, joinNode, result)
	}
	// First entry should have been evicted
	token1 := &Token{ID: "token1"}
	fact1 := &Fact{ID: "fact1", Type: "Test"}
	if _, found := cache.GetJoinResult(token1, fact1, joinNode); found {
		t.Error("Expected first entry to be evicted")
	}
	// Last two should still be there
	token2 := &Token{ID: "token2"}
	fact2 := &Fact{ID: "fact2", Type: "Test"}
	if _, found := cache.GetJoinResult(token2, fact2, joinNode); !found {
		t.Error("Expected second entry to be in cache")
	}
	token3 := &Token{ID: "token3"}
	fact3 := &Fact{ID: "fact3", Type: "Test"}
	if _, found := cache.GetJoinResult(token3, fact3, joinNode); !found {
		t.Error("Expected third entry to be in cache")
	}
}

// TestCacheTTL tests TTL expiration
func TestCacheTTL(t *testing.T) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	config.BetaJoinResultCacheTTL = 50 * time.Millisecond
	cache := NewBetaJoinCache(config)
	leftToken := &Token{ID: "token1"}
	rightFact := &Fact{ID: "fact1", Type: "Test"}
	joinNode := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	// Set a result
	result := &JoinResult{Matched: true}
	cache.SetJoinResult(leftToken, rightFact, joinNode, result)
	// Should be in cache immediately
	if _, found := cache.GetJoinResult(leftToken, rightFact, joinNode); !found {
		t.Error("Expected cache hit immediately after set")
	}
	// Wait for TTL to expire
	time.Sleep(100 * time.Millisecond)
	// Should not be in cache after expiration
	if _, found := cache.GetJoinResult(leftToken, rightFact, joinNode); found {
		t.Error("Expected cache miss after TTL expiration")
	}
}

// TestInvalidateForFact tests fact invalidation
func TestInvalidateForFact(t *testing.T) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	cache := NewBetaJoinCache(config)
	leftToken := &Token{ID: "token1"}
	rightFact := &Fact{ID: "fact1", Type: "Test"}
	joinNode := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	// Set a result
	result := &JoinResult{Matched: true}
	cache.SetJoinResult(leftToken, rightFact, joinNode, result)
	// Verify it's in cache
	if _, found := cache.GetJoinResult(leftToken, rightFact, joinNode); !found {
		t.Fatal("Expected cache hit before invalidation")
	}
	// Invalidate
	invalidated := cache.InvalidateForFact(rightFact.GetInternalID())
	if invalidated == 0 {
		t.Error("Expected at least 1 invalidation")
	}
	// Should not be in cache after invalidation
	if _, found := cache.GetJoinResult(leftToken, rightFact, joinNode); found {
		t.Error("Expected cache miss after invalidation")
	}
}

// TestInvalidateForToken tests token invalidation
func TestInvalidateForToken(t *testing.T) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	cache := NewBetaJoinCache(config)
	leftToken := &Token{ID: "token1"}
	rightFact := &Fact{ID: "fact1", Type: "Test"}
	joinNode := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	// Set a result
	result := &JoinResult{Matched: true}
	cache.SetJoinResult(leftToken, rightFact, joinNode, result)
	// Verify it's in cache
	if _, found := cache.GetJoinResult(leftToken, rightFact, joinNode); !found {
		t.Fatal("Expected cache hit before invalidation")
	}
	// Invalidate
	invalidated := cache.InvalidateForToken(leftToken.ID)
	if invalidated == 0 {
		t.Error("Expected at least 1 invalidation")
	}
	// Should not be in cache after invalidation
	if _, found := cache.GetJoinResult(leftToken, rightFact, joinNode); found {
		t.Error("Expected cache miss after invalidation")
	}
}

// TestBetaJoinCache_Clear tests cache clearing
func TestBetaJoinCache_Clear(t *testing.T) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	cache := NewBetaJoinCache(config)
	joinNode := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	// Add multiple entries
	for i := 1; i <= 5; i++ {
		token := &Token{ID: fmt.Sprintf("token%d", i)}
		fact := &Fact{ID: fmt.Sprintf("fact%d", i), Type: "Test"}
		result := &JoinResult{Matched: true}
		cache.SetJoinResult(token, fact, joinNode, result)
	}
	// Verify size
	size := cache.GetSize()
	if size != 5 {
		t.Errorf("Expected cache size 5, got %d", size)
	}
	// Clear
	cache.Clear()
	// Verify empty
	size = cache.GetSize()
	if size != 0 {
		t.Errorf("Expected cache size 0 after clear, got %d", size)
	}
	// Verify entries not found
	token1 := &Token{ID: "token1"}
	fact1 := &Fact{ID: "fact1", Type: "Test"}
	if _, found := cache.GetJoinResult(token1, fact1, joinNode); found {
		t.Error("Expected cache miss after clear")
	}
}

// TestGetStats tests statistics retrieval
func TestGetStats(t *testing.T) {
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	config.BetaJoinResultCacheMaxSize = 100
	config.BetaJoinResultCacheTTL = time.Minute
	cache := NewBetaJoinCache(config)
	stats := cache.GetStats()
	if !stats["enabled"].(bool) {
		t.Error("Expected enabled=true")
	}
	if stats["capacity"].(int) != 100 {
		t.Errorf("Expected capacity 100, got %d", stats["capacity"])
	}
	if stats["ttl_seconds"].(float64) != 60.0 {
		t.Errorf("Expected TTL 60 seconds, got %.1f", stats["ttl_seconds"])
	}
	if stats["size"].(int) != 0 {
		t.Errorf("Expected initial size 0, got %d", stats["size"])
	}
}

// TestGetStats_Disabled tests stats when cache is disabled
func TestGetStats_Disabled(t *testing.T) {
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = false
	cache := NewBetaJoinCache(config)
	stats := cache.GetStats()
	if stats["enabled"].(bool) {
		t.Error("Expected enabled=false")
	}
}

// TestGetSize tests size retrieval
func TestGetSize(t *testing.T) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	cache := NewBetaJoinCache(config)
	if cache.GetSize() != 0 {
		t.Error("Expected initial size 0")
	}
	joinNode := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	// Add entries
	for i := 1; i <= 3; i++ {
		token := &Token{ID: fmt.Sprintf("token%d", i)}
		fact := &Fact{ID: fmt.Sprintf("fact%d", i), Type: "Test"}
		result := &JoinResult{Matched: true}
		cache.SetJoinResult(token, fact, joinNode, result)
	}
	if cache.GetSize() != 3 {
		t.Errorf("Expected size 3, got %d", cache.GetSize())
	}
}

// TestCleanExpired tests cleaning of expired entries
func TestCleanExpired(t *testing.T) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	config.BetaJoinResultCacheTTL = 50 * time.Millisecond
	cache := NewBetaJoinCache(config)
	joinNode := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	// Add entries
	for i := 1; i <= 3; i++ {
		token := &Token{ID: fmt.Sprintf("token%d", i)}
		fact := &Fact{ID: fmt.Sprintf("fact%d", i), Type: "Test"}
		result := &JoinResult{Matched: true}
		cache.SetJoinResult(token, fact, joinNode, result)
	}
	// Wait for expiration
	time.Sleep(100 * time.Millisecond)
	// Clean expired
	cleaned := cache.CleanExpired()
	if cleaned == 0 {
		t.Error("Expected at least 1 expired entry to be cleaned")
	}
	// Verify cache is empty or smaller
	size := cache.GetSize()
	if size >= 3 {
		t.Errorf("Expected size < 3 after cleaning, got %d", size)
	}
}

// TestResetStats tests statistics reset
func TestResetStats(t *testing.T) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	cache := NewBetaJoinCache(config)
	leftToken := &Token{ID: "token1"}
	rightFact := &Fact{ID: "fact1", Type: "Test"}
	joinNode := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	// Generate some hits and misses
	cache.GetJoinResult(leftToken, rightFact, joinNode) // miss
	result := &JoinResult{Matched: true}
	cache.SetJoinResult(leftToken, rightFact, joinNode, result)
	cache.GetJoinResult(leftToken, rightFact, joinNode) // hit
	stats := cache.GetStats()
	if stats["hits"].(int64) == 0 {
		t.Error("Expected non-zero hits before reset")
	}
	// Reset
	cache.ResetStats()
	stats = cache.GetStats()
	if stats["hits"].(int64) != 0 {
		t.Errorf("Expected 0 hits after reset, got %d", stats["hits"])
	}
	if stats["misses"].(int64) != 0 {
		t.Errorf("Expected 0 misses after reset, got %d", stats["misses"])
	}
}

// TestDifferentJoinNodes tests caching with different join nodes
func TestDifferentJoinNodes(t *testing.T) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	cache := NewBetaJoinCache(config)
	leftToken := &Token{ID: "token1"}
	rightFact := &Fact{ID: "fact1", Type: "Test"}
	joinNode1 := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	joinNode2 := NewJoinNode("join2", nil, []string{"a"}, []string{"b"}, nil, storage)
	// Set result for joinNode1
	result1 := &JoinResult{Matched: true, JoinType: "type1"}
	cache.SetJoinResult(leftToken, rightFact, joinNode1, result1)
	// Set different result for joinNode2
	result2 := &JoinResult{Matched: false, JoinType: "type2"}
	cache.SetJoinResult(leftToken, rightFact, joinNode2, result2)
	// Verify correct results are retrieved
	retrieved1, found1 := cache.GetJoinResult(leftToken, rightFact, joinNode1)
	if !found1 || !retrieved1.Matched || retrieved1.JoinType != "type1" {
		t.Error("Expected to retrieve correct result for joinNode1")
	}
	retrieved2, found2 := cache.GetJoinResult(leftToken, rightFact, joinNode2)
	if !found2 || retrieved2.Matched || retrieved2.JoinType != "type2" {
		t.Error("Expected to retrieve correct result for joinNode2")
	}
}

// TestBetaJoinCache_ConcurrentAccess tests thread-safety
func TestBetaJoinCache_ConcurrentAccess(t *testing.T) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	config.BetaJoinResultCacheMaxSize = 1000
	cache := NewBetaJoinCache(config)
	joinNode := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	done := make(chan bool)
	errors := make(chan error, 100)
	// Concurrent writers
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				token := &Token{ID: fmt.Sprintf("token_%d_%d", id, j)}
				fact := &Fact{ID: fmt.Sprintf("fact_%d_%d", id, j), Type: "Test"}
				result := &JoinResult{Matched: true}
				cache.SetJoinResult(token, fact, joinNode, result)
			}
			done <- true
		}(i)
	}
	// Concurrent readers
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				token := &Token{ID: fmt.Sprintf("token_%d_%d", id, j)}
				fact := &Fact{ID: fmt.Sprintf("fact_%d_%d", id, j), Type: "Test"}
				cache.GetJoinResult(token, fact, joinNode)
			}
			done <- true
		}(i)
	}
	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}
	close(errors)
	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent access error: %v", err)
	}
	// Verify cache is still functional
	stats := cache.GetStats()
	if stats["enabled"].(bool) != true {
		t.Error("Cache should still be enabled after concurrent access")
	}
}

// BenchmarkCacheGetHit benchmarks cache get with hit
func BenchmarkCacheGetHit(b *testing.B) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	cache := NewBetaJoinCache(config)
	leftToken := &Token{ID: "token1"}
	rightFact := &Fact{ID: "fact1", Type: "Test"}
	joinNode := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	result := &JoinResult{Matched: true}
	cache.SetJoinResult(leftToken, rightFact, joinNode, result)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetJoinResult(leftToken, rightFact, joinNode)
	}
}

// BenchmarkCacheGetMiss benchmarks cache get with miss
func BenchmarkCacheGetMiss(b *testing.B) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	cache := NewBetaJoinCache(config)
	joinNode := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		token := &Token{ID: fmt.Sprintf("token%d", i)}
		fact := &Fact{ID: fmt.Sprintf("fact%d", i), Type: "Test"}
		cache.GetJoinResult(token, fact, joinNode)
	}
}

// BenchmarkCacheSet benchmarks cache set
func BenchmarkCacheSet(b *testing.B) {
	storage := NewMemoryStorage()
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	config.BetaJoinResultCacheMaxSize = 100000
	cache := NewBetaJoinCache(config)
	joinNode := NewJoinNode("join1", nil, []string{"a"}, []string{"b"}, nil, storage)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		token := &Token{ID: fmt.Sprintf("token%d", i)}
		fact := &Fact{ID: fmt.Sprintf("fact%d", i), Type: "Test"}
		result := &JoinResult{Matched: true}
		cache.SetJoinResult(token, fact, joinNode, result)
	}
}
