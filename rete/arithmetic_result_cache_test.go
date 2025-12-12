// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"fmt"
	"sync"
	"testing"
	"time"
)
// TestNewArithmeticResultCache teste la création du cache
func TestNewArithmeticResultCache(t *testing.T) {
	config := DefaultCacheConfig()
	cache := NewArithmeticResultCache(config)
	if cache == nil {
		t.Fatal("Expected cache to be created, got nil")
	}
	if cache.maxSize != config.MaxSize {
		t.Errorf("Expected maxSize=%d, got %d", config.MaxSize, cache.maxSize)
	}
	if cache.ttl != config.TTL {
		t.Errorf("Expected ttl=%v, got %v", config.TTL, cache.ttl)
	}
	if !cache.enabled {
		t.Error("Expected cache to be enabled by default")
	}
	if cache.GetSize() != 0 {
		t.Errorf("Expected empty cache, got size=%d", cache.GetSize())
	}
}
// TestArithmeticCache_BasicSetAndGet teste les opérations de base
func TestArithmeticCache_BasicSetAndGet(t *testing.T) {
	config := DefaultCacheConfig()
	cache := NewArithmeticResultCache(config)
	key := "test_key"
	value := 42
	// Set
	cache.Set(key, value)
	// Get
	retrieved, found := cache.Get(key)
	if !found {
		t.Fatal("Expected to find key in cache")
	}
	if retrieved != value {
		t.Errorf("Expected value=%d, got %v", value, retrieved)
	}
	// Vérifier statistiques
	stats := cache.GetStatistics()
	if stats.Sets != 1 {
		t.Errorf("Expected Sets=1, got %d", stats.Sets)
	}
	if stats.Hits != 1 {
		t.Errorf("Expected Hits=1, got %d", stats.Hits)
	}
	if stats.Misses != 0 {
		t.Errorf("Expected Misses=0, got %d", stats.Misses)
	}
}
// TestArithmeticCache_CacheMiss teste les cache misses
func TestArithmeticCache_CacheMiss(t *testing.T) {
	config := DefaultCacheConfig()
	cache := NewArithmeticResultCache(config)
	_, found := cache.Get("nonexistent")
	if found {
		t.Error("Expected cache miss for nonexistent key")
	}
	stats := cache.GetStatistics()
	if stats.Misses != 1 {
		t.Errorf("Expected Misses=1, got %d", stats.Misses)
	}
}
// TestArithmeticCache_GenerateCacheKey teste la génération de clés
func TestArithmeticCache_GenerateCacheKey(t *testing.T) {
	tests := []struct {
		name         string
		resultName   string
		dependencies map[string]interface{}
		expectUnique bool
	}{
		{
			name:       "same_params_same_key",
			resultName: "temp_1",
			dependencies: map[string]interface{}{
				"field_a": 10,
				"field_b": 20,
			},
			expectUnique: false,
		},
		{
			name:       "different_values_different_key",
			resultName: "temp_1",
			dependencies: map[string]interface{}{
				"field_a": 10,
				"field_b": 30,
			},
			expectUnique: true,
		},
		{
			name:       "different_result_name",
			resultName: "temp_2",
			dependencies: map[string]interface{}{
				"field_a": 10,
				"field_b": 20,
			},
			expectUnique: true,
		},
	}
	baseKey := GenerateCacheKey("temp_1", map[string]interface{}{
		"field_a": 10,
		"field_b": 20,
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := GenerateCacheKey(tt.resultName, tt.dependencies)
			if key == "" {
				t.Error("Expected non-empty key")
			}
			if tt.expectUnique {
				if key == baseKey {
					t.Error("Expected different key, got same key")
				}
			} else {
				if key != baseKey {
					t.Error("Expected same key, got different key")
				}
			}
		})
	}
}
// TestArithmeticCache_WithDependencies teste les méthodes avec dépendances
func TestArithmeticCache_WithDependencies(t *testing.T) {
	config := DefaultCacheConfig()
	cache := NewArithmeticResultCache(config)
	resultName := "temp_multiply"
	dependencies := map[string]interface{}{
		"qty":   10,
		"price": 23,
	}
	value := 230
	// Set avec dépendances
	cache.SetWithDependencies(resultName, dependencies, value)
	// Get avec dépendances
	retrieved, found := cache.GetWithDependencies(resultName, dependencies)
	if !found {
		t.Fatal("Expected to find result in cache")
	}
	if retrieved != value {
		t.Errorf("Expected value=%d, got %v", value, retrieved)
	}
	// Get avec dépendances différentes (doit échouer)
	differentDeps := map[string]interface{}{
		"qty":   10,
		"price": 24, // Différent
	}
	_, found = cache.GetWithDependencies(resultName, differentDeps)
	if found {
		t.Error("Expected cache miss for different dependencies")
	}
}
// TestArithmeticCache_LRUEviction teste l'éviction LRU
func TestArithmeticCache_LRUEviction(t *testing.T) {
	config := CacheConfig{
		MaxSize: 3,
		TTL:     time.Hour,
		Enabled: true,
	}
	cache := NewArithmeticResultCache(config)
	// Ajouter 3 entrées
	cache.Set("key1", 1)
	cache.Set("key2", 2)
	cache.Set("key3", 3)
	if cache.GetSize() != 3 {
		t.Errorf("Expected size=3, got %d", cache.GetSize())
	}
	// Accéder à key1 pour la rendre plus récente
	cache.Get("key1")
	// Ajouter une 4ème entrée, doit évincer key2 (la plus ancienne non accédée)
	cache.Set("key4", 4)
	if cache.GetSize() != 3 {
		t.Errorf("Expected size=3 after eviction, got %d", cache.GetSize())
	}
	// key2 doit avoir été évincée
	_, found := cache.Get("key2")
	if found {
		t.Error("Expected key2 to be evicted")
	}
	// key1, key3, key4 doivent être présentes
	if _, found := cache.Get("key1"); !found {
		t.Error("Expected key1 to be present")
	}
	if _, found := cache.Get("key3"); !found {
		t.Error("Expected key3 to be present")
	}
	if _, found := cache.Get("key4"); !found {
		t.Error("Expected key4 to be present")
	}
	// Vérifier statistiques
	stats := cache.GetStatistics()
	if stats.Evictions != 1 {
		t.Errorf("Expected Evictions=1, got %d", stats.Evictions)
	}
}
// TestArithmeticCache_TTLExpiration teste l'expiration par TTL
func TestArithmeticCache_TTLExpiration(t *testing.T) {
	config := CacheConfig{
		MaxSize: 10,
		TTL:     100 * time.Millisecond,
		Enabled: true,
	}
	cache := NewArithmeticResultCache(config)
	// Ajouter une entrée
	cache.Set("key1", 42)
	// Vérifier qu'elle est présente
	value, found := cache.Get("key1")
	if !found || value != 42 {
		t.Fatal("Expected to find key1 immediately after set")
	}
	// Attendre expiration
	time.Sleep(150 * time.Millisecond)
	// Vérifier qu'elle est expirée
	_, found = cache.Get("key1")
	if found {
		t.Error("Expected key1 to be expired")
	}
}
// TestArithmeticCache_EnableDisable teste l'activation/désactivation dynamique
func TestArithmeticCache_EnableDisable(t *testing.T) {
	config := DefaultCacheConfig()
	cache := NewArithmeticResultCache(config)
	// Cache activé par défaut
	cache.Set("key1", 1)
	_, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1 when cache enabled")
	}
	// Désactiver le cache
	cache.SetEnabled(false)
	if cache.IsEnabled() {
		t.Error("Expected cache to be disabled")
	}
	// Les nouvelles écritures ne doivent pas être enregistrées
	cache.Set("key2", 2)
	_, found = cache.Get("key2")
	if found {
		t.Error("Expected not to find key2 when cache disabled")
	}
	// Les anciennes entrées ne doivent pas être accessibles
	_, found = cache.Get("key1")
	if found {
		t.Error("Expected not to find key1 when cache disabled")
	}
	// Réactiver
	cache.SetEnabled(true)
	// key1 doit toujours être dans le cache
	_, found = cache.Get("key1")
	if !found {
		t.Error("Expected to find key1 after re-enabling cache")
	}
}
// TestArithmeticCache_Clear teste le vidage du cache
func TestArithmeticCache_Clear(t *testing.T) {
	config := DefaultCacheConfig()
	cache := NewArithmeticResultCache(config)
	// Ajouter plusieurs entrées
	for i := 0; i < 10; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i)
	}
	if cache.GetSize() != 10 {
		t.Errorf("Expected size=10, got %d", cache.GetSize())
	}
	// Vider le cache
	cache.Clear()
	if cache.GetSize() != 0 {
		t.Errorf("Expected size=0 after clear, got %d", cache.GetSize())
	}
	// Vérifier qu'aucune entrée n'est accessible
	for i := 0; i < 10; i++ {
		_, found := cache.Get(fmt.Sprintf("key%d", i))
		if found {
			t.Errorf("Expected key%d to be cleared", i)
		}
	}
}
// TestArithmeticCache_ConcurrentAccess teste la sécurité thread-safe
func TestArithmeticCache_ConcurrentAccess(t *testing.T) {
	config := DefaultCacheConfig()
	cache := NewArithmeticResultCache(config)
	const numGoroutines = 100
	const opsPerGoroutine = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	// Lancer plusieurs goroutines qui lisent/écrivent en parallèle
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("key_%d_%d", id, j)
				value := id*1000 + j
				// Set
				cache.Set(key, value)
				// Get
				retrieved, found := cache.Get(key)
				if found && retrieved != value {
					t.Errorf("Concurrent access error: expected %d, got %v", value, retrieved)
				}
			}
		}(i)
	}
	wg.Wait()
	// Vérifier que le cache est dans un état cohérent
	stats := cache.GetStatistics()
	if stats.Sets == 0 {
		t.Error("Expected some Sets to be recorded")
	}
}
// TestArithmeticCache_GetHitRate teste le calcul du taux de succès
func TestArithmeticCache_GetHitRate(t *testing.T) {
	config := DefaultCacheConfig()
	cache := NewArithmeticResultCache(config)
	// Cache vide, hit rate = 0
	if cache.GetHitRate() != 0.0 {
		t.Errorf("Expected hit rate 0.0 for empty cache, got %f", cache.GetHitRate())
	}
	// Ajouter une entrée
	cache.Set("key1", 1)
	// 3 hits
	cache.Get("key1")
	cache.Get("key1")
	cache.Get("key1")
	// 2 misses
	cache.Get("key2")
	cache.Get("key3")
	// Hit rate = 3 / (3 + 2) = 0.6
	hitRate := cache.GetHitRate()
	expected := 0.6
	if hitRate < expected-0.01 || hitRate > expected+0.01 {
		t.Errorf("Expected hit rate ~%f, got %f", expected, hitRate)
	}
}
// TestArithmeticCache_GetTopEntries teste la récupération des entrées les plus utilisées
func TestArithmeticCache_GetTopEntries(t *testing.T) {
	config := DefaultCacheConfig()
	cache := NewArithmeticResultCache(config)
	// Ajouter des entrées avec différents nombres d'accès
	cache.Set("key1", 1)
	cache.Set("key2", 2)
	cache.Set("key3", 3)
	// Accéder avec différentes fréquences
	for i := 0; i < 10; i++ {
		cache.Get("key1")
	}
	for i := 0; i < 5; i++ {
		cache.Get("key2")
	}
	cache.Get("key3")
	// Récupérer top 2
	topEntries := cache.GetTopEntries(2)
	if len(topEntries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(topEntries))
	}
	// key1 doit être en premier (10 hits)
	if topEntries[0].HitCount != 10 {
		t.Errorf("Expected top entry to have 10 hits, got %d", topEntries[0].HitCount)
	}
	// key2 doit être en second (5 hits)
	if topEntries[1].HitCount != 5 {
		t.Errorf("Expected second entry to have 5 hits, got %d", topEntries[1].HitCount)
	}
}
// TestArithmeticCache_ResetStatistics teste la réinitialisation des stats
func TestArithmeticCache_ResetStatistics(t *testing.T) {
	config := DefaultCacheConfig()
	cache := NewArithmeticResultCache(config)
	// Générer des statistiques
	cache.Set("key1", 1)
	cache.Get("key1")
	cache.Get("key2") // miss
	stats := cache.GetStatistics()
	if stats.Hits == 0 || stats.Misses == 0 {
		t.Fatal("Expected some statistics to be recorded")
	}
	// Réinitialiser
	cache.ResetStatistics()
	stats = cache.GetStatistics()
	if stats.Hits != 0 {
		t.Errorf("Expected Hits=0 after reset, got %d", stats.Hits)
	}
	if stats.Misses != 0 {
		t.Errorf("Expected Misses=0 after reset, got %d", stats.Misses)
	}
	// Les entrées doivent toujours être présentes
	if cache.GetSize() != 1 {
		t.Errorf("Expected size=1 after stats reset, got %d", cache.GetSize())
	}
}
// TestArithmeticCache_Purge teste le nettoyage des entrées expirées
func TestArithmeticCache_Purge(t *testing.T) {
	config := CacheConfig{
		MaxSize: 10,
		TTL:     100 * time.Millisecond,
		Enabled: true,
	}
	cache := NewArithmeticResultCache(config)
	// Ajouter des entrées
	cache.Set("key1", 1)
	time.Sleep(50 * time.Millisecond)
	cache.Set("key2", 2)
	cache.Set("key3", 3)
	// Attendre que key1 expire
	time.Sleep(60 * time.Millisecond)
	// Purger
	expiredCount := cache.Purge()
	if expiredCount != 1 {
		t.Errorf("Expected 1 expired entry, got %d", expiredCount)
	}
	if cache.GetSize() != 2 {
		t.Errorf("Expected size=2 after purge, got %d", cache.GetSize())
	}
	// key1 doit être absente
	_, found := cache.Get("key1")
	if found {
		t.Error("Expected key1 to be purged")
	}
	// key2 et key3 doivent être présentes
	if _, found := cache.Get("key2"); !found {
		t.Error("Expected key2 to be present")
	}
	if _, found := cache.Get("key3"); !found {
		t.Error("Expected key3 to be present")
	}
}
// TestArithmeticCache_AutoPurge teste le nettoyage automatique
func TestArithmeticCache_AutoPurge(t *testing.T) {
	config := CacheConfig{
		MaxSize: 10,
		TTL:     100 * time.Millisecond,
		Enabled: true,
	}
	cache := NewArithmeticResultCache(config)
	// Démarrer auto-purge
	stopChan := cache.StartAutoPurge(50 * time.Millisecond)
	defer close(stopChan)
	// Ajouter une entrée
	cache.Set("key1", 1)
	// Attendre expiration et auto-purge
	time.Sleep(200 * time.Millisecond)
	// L'entrée doit avoir été purgée automatiquement
	if cache.GetSize() != 0 {
		t.Errorf("Expected size=0 after auto-purge, got %d", cache.GetSize())
	}
}
// TestArithmeticCache_EvictionCallback teste le callback d'éviction
func TestArithmeticCache_EvictionCallback(t *testing.T) {
	evictedKeys := make([]string, 0)
	var mu sync.Mutex
	config := CacheConfig{
		MaxSize: 2,
		TTL:     time.Hour,
		Enabled: true,
		OnEviction: func(key string, value interface{}) {
			mu.Lock()
			defer mu.Unlock()
			evictedKeys = append(evictedKeys, key)
		},
	}
	cache := NewArithmeticResultCache(config)
	// Ajouter 3 entrées (la première sera évincée)
	cache.Set("key1", 1)
	cache.Set("key2", 2)
	cache.Set("key3", 3)
	mu.Lock()
	if len(evictedKeys) != 1 {
		t.Errorf("Expected 1 evicted key, got %d", len(evictedKeys))
	}
	if evictedKeys[0] != "key1" {
		t.Errorf("Expected key1 to be evicted, got %s", evictedKeys[0])
	}
	mu.Unlock()
}
// TestArithmeticCache_EstimateMemoryUsage teste l'estimation mémoire
func TestArithmeticCache_EstimateMemoryUsage(t *testing.T) {
	config := DefaultCacheConfig()
	cache := NewArithmeticResultCache(config)
	// Cache vide
	if cache.EstimateMemoryUsage() != 0 {
		t.Errorf("Expected 0 bytes for empty cache, got %d", cache.EstimateMemoryUsage())
	}
	// Ajouter des entrées
	for i := 0; i < 10; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i)
	}
	usage := cache.EstimateMemoryUsage()
	if usage == 0 {
		t.Error("Expected non-zero memory usage")
	}
	// L'usage doit être proportionnel au nombre d'entrées
	const minBytesPerEntry = 100 // Au minimum
	if usage < int64(cache.GetSize()*minBytesPerEntry) {
		t.Errorf("Expected at least %d bytes, got %d", cache.GetSize()*minBytesPerEntry, usage)
	}
}
// TestArithmeticCache_GetSummary teste le résumé formaté
func TestArithmeticCache_GetSummary(t *testing.T) {
	config := DefaultCacheConfig()
	cache := NewArithmeticResultCache(config)
	// Ajouter quelques opérations
	cache.Set("key1", 1)
	cache.Get("key1")
	cache.Get("key2") // miss
	summary := cache.GetSummary()
	// Vérifier que les clés attendues sont présentes
	expectedKeys := []string{
		"enabled", "size", "max_size", "ttl",
		"hits", "misses", "hit_rate", "evictions", "sets",
		"avg_hit_time", "avg_miss_time",
	}
	for _, key := range expectedKeys {
		if _, exists := summary[key]; !exists {
			t.Errorf("Expected summary to contain key '%s'", key)
		}
	}
	// Vérifier quelques valeurs
	if summary["enabled"] != true {
		t.Error("Expected cache to be enabled in summary")
	}
	if summary["size"] != 1 {
		t.Errorf("Expected size=1 in summary, got %v", summary["size"])
	}
}
// BenchmarkArithmeticCache_Set benchmark les écritures
func BenchmarkArithmeticCache_Set(b *testing.B) {
	config := DefaultCacheConfig()
	cache := NewArithmeticResultCache(config)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%1000)
		cache.Set(key, i)
	}
}
// BenchmarkArithmeticCache_Get_Hit benchmark les lectures avec hits
func BenchmarkArithmeticCache_Get_Hit(b *testing.B) {
	config := DefaultCacheConfig()
	cache := NewArithmeticResultCache(config)
	// Pré-remplir le cache
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%1000)
		cache.Get(key)
	}
}
// BenchmarkArithmeticCache_Get_Miss benchmark les lectures avec misses
func BenchmarkArithmeticCache_Get_Miss(b *testing.B) {
	config := DefaultCacheConfig()
	cache := NewArithmeticResultCache(config)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("nonexistent%d", i)
		cache.Get(key)
	}
}
// BenchmarkArithmeticCache_WithDependencies benchmark avec dépendances
func BenchmarkArithmeticCache_WithDependencies(b *testing.B) {
	config := DefaultCacheConfig()
	cache := NewArithmeticResultCache(config)
	dependencies := map[string]interface{}{
		"field_a": 10,
		"field_b": 20,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultName := fmt.Sprintf("temp_%d", i%100)
		if i%2 == 0 {
			cache.SetWithDependencies(resultName, dependencies, i)
		} else {
			cache.GetWithDependencies(resultName, dependencies)
		}
	}
}