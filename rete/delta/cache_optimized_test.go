// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
)

func TestOptimizedCache(t *testing.T) {
	t.Log("🧪 TEST OPTIMIZED CACHE")
	t.Log("=======================")

	cache := NewOptimizedCache(3, DefaultCacheTTL)

	// Test Put et Get
	delta1 := NewFactDelta("Test~1", "Test")
	delta1.AddFieldChange("field1", "old", "new")

	cache.Put("key1", delta1)

	result, found := cache.Get("key1")
	if !found {
		t.Error("❌ Key1 should be found")
	}
	if result.FactID != "Test~1" {
		t.Errorf("❌ Expected FactID 'Test~1', got '%s'", result.FactID)
	}

	// Test miss
	_, found = cache.Get("nonexistent")
	if found {
		t.Error("❌ Nonexistent key should not be found")
	}

	// Test éviction LRU
	delta2 := NewFactDelta("Test~2", "Test")
	delta3 := NewFactDelta("Test~3", "Test")
	delta4 := NewFactDelta("Test~4", "Test")

	cache.Put("key2", delta2)
	cache.Put("key3", delta3)
	cache.Put("key4", delta4) // Devrait évincer key1 (LRU)

	_, found = cache.Get("key1")
	if found {
		t.Error("❌ Key1 should have been evicted")
	}

	_, found = cache.Get("key2")
	if !found {
		t.Error("❌ Key2 should still be in cache")
	}

	// Test stats
	stats := cache.GetStats()
	if stats.Size != 3 {
		t.Errorf("❌ Expected size 3, got %d", stats.Size)
	}

	if stats.Evictions != 1 {
		t.Errorf("❌ Expected 1 eviction, got %d", stats.Evictions)
	}

	t.Log("✅ OptimizedCache works correctly")
}

func TestOptimizedCache_LRU(t *testing.T) {
	t.Log("🧪 TEST LRU EVICTION")
	t.Log("====================")

	cache := NewOptimizedCache(2, DefaultCacheTTL)

	delta1 := NewFactDelta("Test~1", "Test")
	delta2 := NewFactDelta("Test~2", "Test")
	delta3 := NewFactDelta("Test~3", "Test")

	// Ajouter 2 entrées
	cache.Put("key1", delta1)
	cache.Put("key2", delta2)

	// Accéder à key1 (devient MRU)
	cache.Get("key1")

	// Ajouter key3 - devrait évincer key2 (LRU)
	cache.Put("key3", delta3)

	// key1 devrait toujours exister
	_, found := cache.Get("key1")
	if !found {
		t.Error("❌ Key1 should still be in cache (was accessed)")
	}

	// key2 devrait être évincée
	_, found = cache.Get("key2")
	if found {
		t.Error("❌ Key2 should have been evicted (LRU)")
	}

	// key3 devrait exister
	_, found = cache.Get("key3")
	if !found {
		t.Error("❌ Key3 should be in cache")
	}

	t.Log("✅ LRU eviction works correctly")
}

func TestOptimizedCache_Stats(t *testing.T) {
	t.Log("🧪 TEST CACHE STATS")
	t.Log("===================")

	cache := NewOptimizedCache(10, DefaultCacheTTL)

	delta := NewFactDelta("Test~1", "Test")
	cache.Put("key1", delta)

	// Hit
	cache.Get("key1")
	cache.Get("key1")

	// Miss
	cache.Get("nonexistent")

	stats := cache.GetStats()

	if stats.Hits != 2 {
		t.Errorf("❌ Expected 2 hits, got %d", stats.Hits)
	}

	if stats.Misses != 1 {
		t.Errorf("❌ Expected 1 miss, got %d", stats.Misses)
	}

	expectedRate := 2.0 / 3.0
	if stats.HitRate < expectedRate-0.01 || stats.HitRate > expectedRate+0.01 {
		t.Errorf("❌ Expected hit rate ~%.2f, got %.2f", expectedRate, stats.HitRate)
	}

	t.Log("✅ Cache stats work correctly")
}

func BenchmarkOptimizedCache_Get(b *testing.B) {
	cache := NewOptimizedCache(1000, DefaultCacheTTL)

	// Préremplir
	for i := 0; i < 100; i++ {
		delta := NewFactDelta("Test~1", "Test")
		cache.Put("key", delta)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}

func BenchmarkOptimizedCache_Put(b *testing.B) {
	cache := NewOptimizedCache(1000, DefaultCacheTTL)
	delta := NewFactDelta("Test~1", "Test")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		cache.Put("key", delta)
	}
}

func BenchmarkOptimizedCache_PutEvict(b *testing.B) {
	cache := NewOptimizedCache(10, DefaultCacheTTL)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		delta := NewFactDelta("Test~1", "Test")
		cache.Put("key", delta)
	}
}
