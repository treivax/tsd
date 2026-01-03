// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
	"testing"
	"time"
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

// TestOptimizedCache_Clear teste la méthode Clear (0% couverture).
func TestOptimizedCache_Clear(t *testing.T) {
	t.Log("🧪 TEST: Clear - Vidage complet du cache")

	cache := NewOptimizedCache(10, DefaultCacheTTL)

	// Ajouter plusieurs entrées
	for i := 0; i < 5; i++ {
		delta := NewFactDelta("Test~1", "Test")
		cache.Put(fmt.Sprintf("key%d", i), delta)
	}

	// Vérifier que le cache contient des entrées
	stats := cache.GetStats()
	if stats.Size != 5 {
		t.Fatalf("❌ Cache devrait contenir 5 entrées, got %d", stats.Size)
	}

	// Clear le cache
	cache.Clear()

	// Vérifier que le cache est vide
	stats = cache.GetStats()
	if stats.Size != 0 {
		t.Errorf("❌ Cache devrait être vide après Clear, got size=%d", stats.Size)
	}

	// Vérifier qu'aucune clé n'est accessible
	for i := 0; i < 5; i++ {
		_, found := cache.Get(fmt.Sprintf("key%d", i))
		if found {
			t.Errorf("❌ key%d ne devrait plus être dans le cache après Clear", i)
		}
	}

	// Vérifier qu'on peut ajouter de nouvelles entrées après Clear
	delta := NewFactDelta("Test~new", "Test")
	cache.Put("newkey", delta)

	result, found := cache.Get("newkey")
	if !found {
		t.Error("❌ Devrait pouvoir ajouter des entrées après Clear")
	}
	if result.FactID != "Test~new" {
		t.Errorf("❌ Expected FactID 'Test~new', got '%s'", result.FactID)
	}

	t.Log("✅ Clear vide correctement le cache")
}

// TestOptimizedCache_RemoveExpired teste la méthode removeExpired (0% couverture).
func TestOptimizedCache_RemoveExpired(t *testing.T) {
	t.Log("🧪 TEST: removeExpired - Suppression entrées expirées")

	// Créer cache avec TTL très court (1ms)
	cache := NewOptimizedCache(10, 1*time.Millisecond)

	delta := NewFactDelta("Test~1", "Test")
	cache.Put("key1", delta)

	// Vérifier que l'entrée existe
	_, found := cache.Get("key1")
	if !found {
		t.Fatal("❌ key1 devrait exister juste après insertion")
	}

	// Attendre expiration
	time.Sleep(5 * time.Millisecond)

	// Appeler removeExpired manuellement (normalement appelé en interne)
	cache.removeExpired("key1")

	// Vérifier que l'entrée a été supprimée
	stats := cache.GetStats()
	if stats.Size != 0 {
		t.Errorf("❌ Cache devrait être vide après removeExpired, got size=%d", stats.Size)
	}

	if stats.Evictions != 1 {
		t.Errorf("❌ Expected 1 eviction, got %d", stats.Evictions)
	}

	t.Log("✅ removeExpired supprime correctement les entrées expirées")
}

// TestOptimizedCache_RemoveNode teste indirectement removeNode via éviction.
func TestOptimizedCache_RemoveNode(t *testing.T) {
	t.Log("🧪 TEST: removeNode - Suppression nœud LRU (via éviction)")

	cache := NewOptimizedCache(2, DefaultCacheTTL)

	delta1 := NewFactDelta("Test~1", "Test")
	delta2 := NewFactDelta("Test~2", "Test")
	delta3 := NewFactDelta("Test~3", "Test")

	// Remplir le cache à capacité max
	cache.Put("key1", delta1)
	cache.Put("key2", delta2)

	stats := cache.GetStats()
	if stats.Size != 2 {
		t.Fatalf("❌ Cache devrait contenir 2 entrées, got %d", stats.Size)
	}

	// Forcer une éviction (removeNode sera appelé en interne)
	cache.Put("key3", delta3)

	// Vérifier qu'une éviction a eu lieu
	stats = cache.GetStats()
	if stats.Size != 2 {
		t.Errorf("❌ Cache devrait toujours contenir 2 entrées, got %d", stats.Size)
	}

	if stats.Evictions != 1 {
		t.Errorf("❌ Expected 1 eviction (removeNode called), got %d", stats.Evictions)
	}

	// key1 (LRU) devrait être évincée
	_, found := cache.Get("key1")
	if found {
		t.Error("❌ key1 devrait avoir été évincée (removeNode)")
	}

	// key3 devrait exister
	_, found = cache.Get("key3")
	if !found {
		t.Error("❌ key3 devrait être dans le cache")
	}

	t.Log("✅ removeNode (via éviction LRU) fonctionne correctement")
}

// TestOptimizedCache_TTLExpiration teste l'expiration TTL complète.
func TestOptimizedCache_TTLExpiration(t *testing.T) {
	t.Log("🧪 TEST: TTL Expiration - Vérification expiration automatique")

	// Cache avec TTL très court
	cache := NewOptimizedCache(10, 2*time.Millisecond)

	delta := NewFactDelta("Test~1", "Test")
	cache.Put("expiring_key", delta)

	// Vérifier existence immédiate
	_, found := cache.Get("expiring_key")
	if !found {
		t.Fatal("❌ Clé devrait exister immédiatement après insertion")
	}

	// Attendre expiration
	time.Sleep(5 * time.Millisecond)

	// Get devrait retourner not found (car TTL expiré)
	_, found = cache.Get("expiring_key")
	if found {
		t.Error("❌ Clé devrait avoir expiré après TTL")
	}

	t.Log("✅ TTL expiration fonctionne correctement")
}
