// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"sync"
	"testing"
	"time"
	"github.com/treivax/tsd/constraint"
)
// TestNewNormalizationCache teste la création d'un nouveau cache
func TestNewNormalizationCache(t *testing.T) {
	cache := NewNormalizationCache(100)
	if cache == nil {
		t.Fatal("Expected cache to be created, got nil")
	}
	if !cache.IsEnabled() {
		t.Error("Expected cache to be enabled by default")
	}
	if cache.maxSize != 100 {
		t.Errorf("Expected maxSize 100, got %d", cache.maxSize)
	}
	if cache.Size() != 0 {
		t.Errorf("Expected empty cache, got size %d", cache.Size())
	}
	if cache.eviction != "lru" {
		t.Errorf("Expected eviction strategy 'lru', got '%s'", cache.eviction)
	}
}
// TestCacheEnableDisable teste l'activation/désactivation du cache
func TestCacheEnableDisable(t *testing.T) {
	cache := NewNormalizationCache(10)
	// Vérifier activé par défaut
	if !cache.IsEnabled() {
		t.Error("Cache should be enabled by default")
	}
	// Désactiver
	cache.Disable()
	if cache.IsEnabled() {
		t.Error("Cache should be disabled after Disable()")
	}
	// Activer
	cache.Enable()
	if !cache.IsEnabled() {
		t.Error("Cache should be enabled after Enable()")
	}
}
// TestCacheGetSet teste les opérations de base Get/Set
func TestCacheGetSet(t *testing.T) {
	cache := NewNormalizationCache(10)
	key := "test_key_1"
	value := "test_value_1"
	// Le cache devrait être vide
	if _, found := cache.Get(key); found {
		t.Error("Expected cache miss for new key")
	}
	// Ajouter une valeur
	cache.Set(key, value)
	// Récupérer la valeur
	retrieved, found := cache.Get(key)
	if !found {
		t.Error("Expected cache hit after Set")
	}
	if retrieved != value {
		t.Errorf("Expected value '%v', got '%v'", value, retrieved)
	}
}
// TestCacheStats teste les statistiques du cache
func TestCacheStats(t *testing.T) {
	cache := NewNormalizationCache(10)
	// Stats initiales
	stats := cache.GetStats()
	if stats.Hits != 0 || stats.Misses != 0 {
		t.Errorf("Expected 0 hits and misses, got %d hits, %d misses", stats.Hits, stats.Misses)
	}
	// Ajouter des entrées
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	// Cache miss
	cache.Get("key3")
	// Cache hits
	cache.Get("key1")
	cache.Get("key2")
	cache.Get("key1")
	// Vérifier les stats
	stats = cache.GetStats()
	if stats.Hits != 3 {
		t.Errorf("Expected 3 hits, got %d", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}
	if stats.Size != 2 {
		t.Errorf("Expected size 2, got %d", stats.Size)
	}
	expectedHitRate := 0.75 // 3/(3+1)
	if stats.HitRate != expectedHitRate {
		t.Errorf("Expected hit rate %.2f, got %.2f", expectedHitRate, stats.HitRate)
	}
}
// TestCacheClear teste le vidage du cache
func TestCacheClear(t *testing.T) {
	cache := NewNormalizationCache(10)
	// Ajouter des entrées
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	if cache.Size() != 3 {
		t.Errorf("Expected size 3, got %d", cache.Size())
	}
	// Vider le cache
	cache.Clear()
	if cache.Size() != 0 {
		t.Errorf("Expected size 0 after Clear, got %d", cache.Size())
	}
	// Vérifier que les entrées sont bien supprimées
	if _, found := cache.Get("key1"); found {
		t.Error("Expected cache miss after Clear")
	}
}
// TestCacheResetStats teste la réinitialisation des statistiques
func TestCacheResetStats(t *testing.T) {
	cache := NewNormalizationCache(10)
	// Générer des stats
	cache.Set("key1", "value1")
	cache.Get("key1")
	cache.Get("key2")
	stats := cache.GetStats()
	if stats.Hits == 0 && stats.Misses == 0 {
		t.Error("Expected some stats before reset")
	}
	// Réinitialiser
	cache.ResetStats()
	stats = cache.GetStats()
	if stats.Hits != 0 || stats.Misses != 0 {
		t.Errorf("Expected 0 hits and misses after reset, got %d hits, %d misses", stats.Hits, stats.Misses)
	}
	// La taille du cache devrait rester inchangée
	if stats.Size != 1 {
		t.Errorf("Expected size 1 after ResetStats, got %d", stats.Size)
	}
}
// TestCacheEvictionLRU teste l'éviction LRU
func TestCacheEvictionLRU(t *testing.T) {
	cache := NewNormalizationCache(3)
	// Remplir le cache
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	if cache.Size() != 3 {
		t.Errorf("Expected size 3, got %d", cache.Size())
	}
	// Accéder à key1 pour la marquer comme récente
	cache.Get("key1")
	// Ajouter key4 devrait évincer key2 (la moins récemment utilisée)
	cache.Set("key4", "value4")
	if cache.Size() != 3 {
		t.Errorf("Expected size 3 after eviction, got %d", cache.Size())
	}
	// key2 devrait être évincée
	if _, found := cache.Get("key2"); found {
		t.Error("Expected key2 to be evicted")
	}
	// key1, key3, key4 devraient être présentes
	if _, found := cache.Get("key1"); !found {
		t.Error("Expected key1 to be present")
	}
	if _, found := cache.Get("key3"); !found {
		t.Error("Expected key3 to be present")
	}
	if _, found := cache.Get("key4"); !found {
		t.Error("Expected key4 to be present")
	}
}
// TestCacheDisabledGetSet teste que Get/Set ne font rien quand le cache est désactivé
func TestCacheDisabledGetSet(t *testing.T) {
	cache := NewNormalizationCache(10)
	cache.Disable()
	cache.Set("key1", "value1")
	if cache.Size() != 0 {
		t.Errorf("Expected size 0 when disabled, got %d", cache.Size())
	}
	if _, found := cache.Get("key1"); found {
		t.Error("Expected cache miss when disabled")
	}
}
// TestComputeCacheKey teste le calcul de clés de cache
func TestComputeCacheKey(t *testing.T) {
	// Créer deux expressions identiques
	expr1 := constraint.BinaryOperation{
		Type:     "binaryOperation",
		Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		Operator: ">",
		Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	}
	expr2 := constraint.BinaryOperation{
		Type:     "binaryOperation",
		Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		Operator: ">",
		Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	}
	key1 := computeCacheKey(expr1)
	key2 := computeCacheKey(expr2)
	// Les clés devraient être identiques pour des expressions identiques
	if key1 != key2 {
		t.Errorf("Expected same keys for identical expressions, got %s and %s", key1, key2)
	}
	// Expression différente
	expr3 := constraint.BinaryOperation{
		Type:     "binaryOperation",
		Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		Operator: ">=",
		Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	}
	key3 := computeCacheKey(expr3)
	// Les clés devraient être différentes
	if key1 == key3 {
		t.Error("Expected different keys for different expressions")
	}
}
// TestNormalizeExpressionWithCache teste la normalisation avec cache
func TestNormalizeExpressionWithCache(t *testing.T) {
	cache := NewNormalizationCache(10)
	// Expression à normaliser
	expr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
			Operator: ">=",
			Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
					Operator: ">",
					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
				},
			},
		},
	}
	// Première normalisation (cache miss)
	normalized1, err := NormalizeExpressionWithCache(expr, cache)
	if err != nil {
		t.Fatalf("NormalizeExpressionWithCache failed: %v", err)
	}
	stats := cache.GetStats()
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}
	// Deuxième normalisation (cache hit)
	normalized2, err := NormalizeExpressionWithCache(expr, cache)
	if err != nil {
		t.Fatalf("NormalizeExpressionWithCache failed: %v", err)
	}
	stats = cache.GetStats()
	if stats.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", stats.Hits)
	}
	// Les résultats devraient être identiques
	conds1, _, _ := ExtractConditions(normalized1)
	conds2, _, _ := ExtractConditions(normalized2)
	if len(conds1) != len(conds2) {
		t.Fatalf("Different number of conditions: %d vs %d", len(conds1), len(conds2))
	}
	for i := range conds1 {
		if !CompareConditions(conds1[i], conds2[i]) {
			t.Errorf("Condition %d differs", i)
		}
	}
}
// TestNormalizeExpressionWithCacheDisabled teste avec cache désactivé
func TestNormalizeExpressionWithCacheDisabled(t *testing.T) {
	cache := NewNormalizationCache(10)
	cache.Disable()
	expr := constraint.BinaryOperation{
		Type:     "binaryOperation",
		Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		Operator: ">",
		Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	}
	// Normaliser avec cache désactivé
	_, err := NormalizeExpressionWithCache(expr, cache)
	if err != nil {
		t.Fatalf("NormalizeExpressionWithCache failed: %v", err)
	}
	// Aucune stat ne devrait être enregistrée
	stats := cache.GetStats()
	if stats.Hits != 0 || stats.Misses != 0 {
		t.Errorf("Expected 0 hits and misses when disabled, got %d hits, %d misses", stats.Hits, stats.Misses)
	}
}
// TestCacheConcurrency teste l'accès concurrent au cache
func TestCacheConcurrency(t *testing.T) {
	cache := NewNormalizationCache(100)
	const numGoroutines = 10
	const numOperations = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	// Lancer plusieurs goroutines qui accèdent au cache
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := "key" + string(rune(j%10))
				value := "value" + string(rune(id))
				// Set
				cache.Set(key, value)
				// Get
				cache.Get(key)
				// Stats
				_ = cache.GetStats()
			}
		}(i)
	}
	wg.Wait()
	// Vérifier que le cache n'est pas corrompu
	stats := cache.GetStats()
	if stats.Size > cache.maxSize {
		t.Errorf("Cache size exceeds max: %d > %d", stats.Size, cache.maxSize)
	}
}
// TestGlobalCache teste le cache global
func TestGlobalCache(t *testing.T) {
	// Créer et définir le cache global
	cache := NewNormalizationCache(10)
	SetGlobalCache(cache)
	// Récupérer le cache global
	retrieved := GetGlobalCache()
	if retrieved != cache {
		t.Error("Expected to retrieve the same cache instance")
	}
	// Tester NormalizeExpressionCached
	expr := constraint.BinaryOperation{
		Type:     "binaryOperation",
		Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		Operator: ">",
		Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	}
	_, err := NormalizeExpressionCached(expr)
	if err != nil {
		t.Fatalf("NormalizeExpressionCached failed: %v", err)
	}
	stats := cache.GetStats()
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}
	// Nettoyer
	SetGlobalCache(nil)
}
// TestSetCacheMaxSize teste le changement de taille max
func TestSetCacheMaxSize(t *testing.T) {
	cache := NewNormalizationCache(10)
	// Remplir le cache
	for i := 0; i < 10; i++ {
		cache.Set("key"+string(rune(i)), "value"+string(rune(i)))
	}
	if cache.Size() != 10 {
		t.Errorf("Expected size 10, got %d", cache.Size())
	}
	// Réduire la taille max
	cache.SetCacheMaxSize(5)
	if cache.maxSize != 5 {
		t.Errorf("Expected maxSize 5, got %d", cache.maxSize)
	}
	if cache.Size() > 5 {
		t.Errorf("Expected size <= 5 after SetCacheMaxSize, got %d", cache.Size())
	}
}
// TestSetEvictionStrategy teste le changement de stratégie d'éviction
func TestSetEvictionStrategy(t *testing.T) {
	cache := NewNormalizationCache(10)
	if cache.eviction != "lru" {
		t.Errorf("Expected initial eviction 'lru', got '%s'", cache.eviction)
	}
	// Changer la stratégie
	cache.SetEvictionStrategy("fifo")
	if cache.eviction != "fifo" {
		t.Errorf("Expected eviction 'fifo', got '%s'", cache.eviction)
	}
}
// TestCacheStatsString teste la méthode String de CacheStats
func TestCacheStatsString(t *testing.T) {
	stats := CacheStats{
		Hits:     100,
		Misses:   25,
		Size:     50,
		MaxSize:  100,
		HitRate:  0.8,
		Enabled:  true,
		Eviction: "lru",
	}
	str := stats.String()
	if str == "" {
		t.Error("Expected non-empty string")
	}
	// Vérifier que la string contient les informations importantes
	if !containsSubstring(str, "100") || !containsSubstring(str, "25") || !containsSubstring(str, "50") {
		t.Error("String should contain stats values")
	}
}
// TestCachePerformance teste les performances du cache
func TestCachePerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}
	cache := NewNormalizationCache(1000)
	// Expression à normaliser
	expr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
			Operator: ">",
			Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
					Operator: ">=",
					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
				},
			},
		},
	}
	const iterations = 10000
	// Sans cache
	startNoCache := time.Now()
	for i := 0; i < iterations; i++ {
		_, _ = NormalizeExpression(expr)
	}
	durationNoCache := time.Since(startNoCache)
	// Avec cache
	startWithCache := time.Now()
	for i := 0; i < iterations; i++ {
		_, _ = NormalizeExpressionWithCache(expr, cache)
	}
	durationWithCache := time.Since(startWithCache)
	t.Logf("Without cache: %v", durationNoCache)
	t.Logf("With cache: %v", durationWithCache)
	t.Logf("Speedup: %.2fx", float64(durationNoCache)/float64(durationWithCache))
	stats := cache.GetStats()
	t.Logf("Cache stats: %s", stats.String())
	// Le cache devrait être plus rapide
	if durationWithCache > durationNoCache {
		t.Logf("Warning: Cache is slower than no cache (might be due to system load)")
	}
	// Vérifier un bon taux de succès
	if stats.HitRate < 0.99 {
		t.Errorf("Expected hit rate > 0.99, got %.2f", stats.HitRate)
	}
}
// containsSubstring vérifie si une string contient une sous-string
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstringAt(s, substr))
}
func containsSubstringAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
// TestNewNormalizationCacheWithEviction teste la création avec stratégie d'éviction
func TestNewNormalizationCacheWithEviction(t *testing.T) {
	tests := []struct {
		name     string
		eviction string
		wantLRU  bool
	}{
		{"LRU", "lru", true},
		{"FIFO", "fifo", false},
		{"None", "none", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewNormalizationCacheWithEviction(10, tt.eviction)
			if cache.eviction != tt.eviction {
				t.Errorf("Expected eviction '%s', got '%s'", tt.eviction, cache.eviction)
			}
			hasLRU := cache.lru != nil
			if hasLRU != tt.wantLRU {
				t.Errorf("Expected LRU tracker: %v, got: %v", tt.wantLRU, hasLRU)
			}
		})
	}
}
// TestGetHitRate teste le calcul du taux de succès
func TestGetHitRate(t *testing.T) {
	cache := NewNormalizationCache(10)
	// Initialement 0
	if rate := cache.GetHitRate(); rate != 0.0 {
		t.Errorf("Expected hit rate 0.0, got %.2f", rate)
	}
	// Ajouter une entrée et générer des hits/misses
	cache.Set("key1", "value1")
	cache.Get("key1") // hit
	cache.Get("key2") // miss
	cache.Get("key1") // hit
	expectedRate := 2.0 / 3.0 // 2 hits, 1 miss
	if rate := cache.GetHitRate(); rate != expectedRate {
		t.Errorf("Expected hit rate %.2f, got %.2f", expectedRate, rate)
	}
}