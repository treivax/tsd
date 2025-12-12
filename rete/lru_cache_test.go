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
// TestNewLRUCache teste la création d'un cache LRU
func TestNewLRUCache(t *testing.T) {
	cache := NewLRUCache(100, 0)
	if cache == nil {
		t.Fatal("NewLRUCache retourne nil")
	}
	if cache.Capacity() != 100 {
		t.Errorf("Capacité devrait être 100, obtenu %d", cache.Capacity())
	}
	if cache.Len() != 0 {
		t.Errorf("Longueur initiale devrait être 0, obtenu %d", cache.Len())
	}
}
// TestNewLRUCache_ZeroCapacity teste la création avec capacité 0
func TestNewLRUCache_ZeroCapacity(t *testing.T) {
	cache := NewLRUCache(0, 0)
	if cache.Capacity() <= 0 {
		t.Error("Capacité devrait être > 0 même avec input 0")
	}
}
// TestLRUCache_SetGet teste Set et Get
func TestLRUCache_SetGet(t *testing.T) {
	cache := NewLRUCache(10, 0)
	// Set
	cache.Set("key1", "value1")
	cache.Set("key2", 42)
	cache.Set("key3", []string{"a", "b"})
	if cache.Len() != 3 {
		t.Errorf("Longueur devrait être 3, obtenu %d", cache.Len())
	}
	// Get existant
	val, ok := cache.Get("key1")
	if !ok {
		t.Error("key1 devrait exister")
	}
	if val != "value1" {
		t.Errorf("Valeur incorrecte pour key1: %v", val)
	}
	// Get non-existant
	_, ok = cache.Get("nonexistent")
	if ok {
		t.Error("nonexistent ne devrait pas exister")
	}
}
// TestLRUCache_Update teste la mise à jour d'une valeur existante
func TestLRUCache_Update(t *testing.T) {
	cache := NewLRUCache(10, 0)
	cache.Set("key1", "value1")
	cache.Set("key1", "value2") // Mise à jour
	if cache.Len() != 1 {
		t.Errorf("Longueur devrait être 1 après mise à jour, obtenu %d", cache.Len())
	}
	val, ok := cache.Get("key1")
	if !ok || val != "value2" {
		t.Error("Valeur devrait être mise à jour")
	}
}
// TestLRUCache_Eviction teste l'éviction LRU
func TestLRUCache_Eviction(t *testing.T) {
	cache := NewLRUCache(3, 0) // Capacité de 3
	// Remplir le cache
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	if cache.Len() != 3 {
		t.Errorf("Cache devrait être plein (3), obtenu %d", cache.Len())
	}
	// Ajouter un 4ème élément devrait évincer le plus ancien (key1)
	cache.Set("key4", "value4")
	if cache.Len() != 3 {
		t.Errorf("Cache devrait toujours avoir 3 éléments, obtenu %d", cache.Len())
	}
	// key1 devrait être évincé
	_, ok := cache.Get("key1")
	if ok {
		t.Error("key1 devrait avoir été évincé")
	}
	// Les autres devraient être présents
	if _, ok := cache.Get("key2"); !ok {
		t.Error("key2 devrait être présent")
	}
	if _, ok := cache.Get("key3"); !ok {
		t.Error("key3 devrait être présent")
	}
	if _, ok := cache.Get("key4"); !ok {
		t.Error("key4 devrait être présent")
	}
}
// TestLRUCache_LRUOrder teste l'ordre LRU
func TestLRUCache_LRUOrder(t *testing.T) {
	cache := NewLRUCache(3, 0)
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	// Accéder à key1 pour le déplacer vers le front
	cache.Get("key1")
	// Ajouter key4 devrait évincer key2 (le plus ancien non accédé)
	cache.Set("key4", "value4")
	// key2 devrait être évincé
	if _, ok := cache.Get("key2"); ok {
		t.Error("key2 devrait avoir été évincé")
	}
	// key1 devrait toujours être présent (récemment accédé)
	if _, ok := cache.Get("key1"); !ok {
		t.Error("key1 devrait être présent")
	}
}
// TestLRUCache_Delete teste la suppression
func TestLRUCache_Delete(t *testing.T) {
	cache := NewLRUCache(10, 0)
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	// Supprimer key1
	deleted := cache.Delete("key1")
	if !deleted {
		t.Error("Delete devrait retourner true pour key1")
	}
	if cache.Len() != 1 {
		t.Errorf("Longueur devrait être 1 après suppression, obtenu %d", cache.Len())
	}
	// Vérifier que key1 n'existe plus
	if _, ok := cache.Get("key1"); ok {
		t.Error("key1 ne devrait plus exister")
	}
	// Supprimer une clé inexistante
	deleted = cache.Delete("nonexistent")
	if deleted {
		t.Error("Delete devrait retourner false pour clé inexistante")
	}
}
// TestLRUCache_Clear teste le vidage du cache
func TestLRUCache_Clear(t *testing.T) {
	cache := NewLRUCache(10, 0)
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	cache.Clear()
	if cache.Len() != 0 {
		t.Errorf("Cache devrait être vide après Clear, obtenu %d", cache.Len())
	}
	// Vérifier qu'aucune clé n'existe
	for _, key := range []string{"key1", "key2", "key3"} {
		if _, ok := cache.Get(key); ok {
			t.Errorf("%s ne devrait plus exister après Clear", key)
		}
	}
}
// TestLRUCache_Stats teste les statistiques
func TestLRUCache_Stats(t *testing.T) {
	cache := NewLRUCache(3, 0)
	// Opérations
	cache.Set("key1", "value1") // Set 1
	cache.Set("key2", "value2") // Set 2
	cache.Get("key1")           // Hit 1
	cache.Get("nonexistent")    // Miss 1
	cache.Set("key3", "value3") // Set 3
	cache.Set("key4", "value4") // Set 4 (éviction)
	stats := cache.GetStats()
	if stats.Hits != 1 {
		t.Errorf("Hits devrait être 1, obtenu %d", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("Misses devrait être 1, obtenu %d", stats.Misses)
	}
	if stats.Sets != 4 {
		t.Errorf("Sets devrait être 4, obtenu %d", stats.Sets)
	}
	if stats.Evictions != 1 {
		t.Errorf("Evictions devrait être 1, obtenu %d", stats.Evictions)
	}
	if stats.Size != 3 {
		t.Errorf("Size devrait être 3, obtenu %d", stats.Size)
	}
	if stats.Capacity != 3 {
		t.Errorf("Capacity devrait être 3, obtenu %d", stats.Capacity)
	}
}
// TestLRUCache_HitRate teste le calcul du taux de hits
func TestLRUCache_HitRate(t *testing.T) {
	cache := NewLRUCache(10, 0)
	// Taux initial (aucun accès)
	if rate := cache.GetHitRate(); rate != 0.0 {
		t.Errorf("HitRate initial devrait être 0.0, obtenu %f", rate)
	}
	cache.Set("key1", "value1")
	cache.Get("key1") // Hit
	cache.Get("key2") // Miss
	cache.Get("key1") // Hit
	cache.Get("key3") // Miss
	// 2 hits / 4 accès = 0.5
	rate := cache.GetHitRate()
	expected := 0.5
	if rate != expected {
		t.Errorf("HitRate devrait être %f, obtenu %f", expected, rate)
	}
}
// TestLRUCache_TTL teste l'expiration avec TTL
func TestLRUCache_TTL(t *testing.T) {
	ttl := 100 * time.Millisecond
	cache := NewLRUCache(10, ttl)
	cache.Set("key1", "value1")
	// Immédiatement après Set, devrait être présent
	if _, ok := cache.Get("key1"); !ok {
		t.Error("key1 devrait être présent immédiatement")
	}
	// Attendre l'expiration
	time.Sleep(150 * time.Millisecond)
	// Après expiration, ne devrait plus être présent
	if _, ok := cache.Get("key1"); ok {
		t.Error("key1 devrait avoir expiré")
	}
}
// TestLRUCache_CleanExpired teste le nettoyage des éléments expirés
func TestLRUCache_CleanExpired(t *testing.T) {
	ttl := 100 * time.Millisecond
	cache := NewLRUCache(10, ttl)
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	// Attendre l'expiration
	time.Sleep(150 * time.Millisecond)
	// Nettoyer les expirés
	cleaned := cache.CleanExpired()
	if cleaned != 3 {
		t.Errorf("CleanExpired devrait retourner 3, obtenu %d", cleaned)
	}
	if cache.Len() != 0 {
		t.Errorf("Cache devrait être vide après nettoyage, obtenu %d", cache.Len())
	}
}
// TestLRUCache_CleanExpired_NoTTL teste CleanExpired sans TTL
func TestLRUCache_CleanExpired_NoTTL(t *testing.T) {
	cache := NewLRUCache(10, 0) // Pas de TTL
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cleaned := cache.CleanExpired()
	if cleaned != 0 {
		t.Errorf("CleanExpired devrait retourner 0 sans TTL, obtenu %d", cleaned)
	}
	if cache.Len() != 2 {
		t.Error("Les éléments ne devraient pas être supprimés sans TTL")
	}
}
// TestLRUCache_Keys teste Keys()
func TestLRUCache_Keys(t *testing.T) {
	cache := NewLRUCache(10, 0)
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	keys := cache.Keys()
	if len(keys) != 3 {
		t.Errorf("Keys devrait retourner 3 clés, obtenu %d", len(keys))
	}
	// Vérifier que toutes les clés sont présentes
	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}
	for _, expected := range []string{"key1", "key2", "key3"} {
		if !keyMap[expected] {
			t.Errorf("Clé %s manquante dans Keys()", expected)
		}
	}
}
// TestLRUCache_OldestNewest teste Oldest() et Newest()
func TestLRUCache_OldestNewest(t *testing.T) {
	cache := NewLRUCache(10, 0)
	// Cache vide
	if _, ok := cache.Oldest(); ok {
		t.Error("Oldest devrait retourner false pour cache vide")
	}
	if _, ok := cache.Newest(); ok {
		t.Error("Newest devrait retourner false pour cache vide")
	}
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	// Le plus récent devrait être key3
	newest, ok := cache.Newest()
	if !ok || newest != "key3" {
		t.Errorf("Newest devrait être key3, obtenu %s", newest)
	}
	// Le plus ancien devrait être key1
	oldest, ok := cache.Oldest()
	if !ok || oldest != "key1" {
		t.Errorf("Oldest devrait être key1, obtenu %s", oldest)
	}
	// Accéder à key1 pour le déplacer
	cache.Get("key1")
	// Maintenant le plus ancien devrait être key2
	oldest, ok = cache.Oldest()
	if !ok || oldest != "key2" {
		t.Errorf("Oldest devrait maintenant être key2, obtenu %s", oldest)
	}
}
// TestLRUCache_Contains teste Contains()
func TestLRUCache_Contains(t *testing.T) {
	cache := NewLRUCache(10, 0)
	cache.Set("key1", "value1")
	if !cache.Contains("key1") {
		t.Error("Contains devrait retourner true pour key1")
	}
	if cache.Contains("nonexistent") {
		t.Error("Contains devrait retourner false pour clé inexistante")
	}
	// Contains ne devrait pas affecter les statistiques
	statsBefore := cache.GetStats()
	cache.Contains("key1")
	statsAfter := cache.GetStats()
	if statsBefore.Hits != statsAfter.Hits {
		t.Error("Contains ne devrait pas incrémenter les hits")
	}
}
// TestLRUCache_ThreadSafety teste la sécurité des threads
func TestLRUCache_ThreadSafety(t *testing.T) {
	cache := NewLRUCache(100, 0)
	var wg sync.WaitGroup
	// Lancer plusieurs goroutines concurrentes
	numGoroutines := 10
	numOpsPerGoroutine := 100
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOpsPerGoroutine; j++ {
				key := fmt.Sprintf("key_%d_%d", id, j)
				cache.Set(key, j)
				cache.Get(key)
				if j%10 == 0 {
					cache.Delete(key)
				}
			}
		}(i)
	}
	wg.Wait()
	// Vérifier que le cache fonctionne toujours
	cache.Set("test", "value")
	if val, ok := cache.Get("test"); !ok || val != "value" {
		t.Error("Cache devrait fonctionner après accès concurrents")
	}
}
// TestLRUCacheStats_Methods teste les méthodes de LRUCacheStats
func TestLRUCacheStats_Methods(t *testing.T) {
	stats := LRUCacheStats{
		Hits:      80,
		Misses:    20,
		Evictions: 5,
		Sets:      100,
		Size:      50,
		Capacity:  100,
	}
	// HitRate
	expectedHitRate := 80.0 / 100.0
	if rate := stats.HitRate(); rate != expectedHitRate {
		t.Errorf("HitRate devrait être %f, obtenu %f", expectedHitRate, rate)
	}
	// EvictionRate
	expectedEvictionRate := 5.0 / 100.0
	if rate := stats.EvictionRate(); rate != expectedEvictionRate {
		t.Errorf("EvictionRate devrait être %f, obtenu %f", expectedEvictionRate, rate)
	}
	// FillRate
	expectedFillRate := 50.0 / 100.0
	if rate := stats.FillRate(); rate != expectedFillRate {
		t.Errorf("FillRate devrait être %f, obtenu %f", expectedFillRate, rate)
	}
}
// TestLRUCacheStats_ZeroValues teste les méthodes avec valeurs nulles
func TestLRUCacheStats_ZeroValues(t *testing.T) {
	stats := LRUCacheStats{}
	if rate := stats.HitRate(); rate != 0.0 {
		t.Errorf("HitRate devrait être 0.0 pour stats vides, obtenu %f", rate)
	}
	if rate := stats.EvictionRate(); rate != 0.0 {
		t.Errorf("EvictionRate devrait être 0.0 pour stats vides, obtenu %f", rate)
	}
	if rate := stats.FillRate(); rate != 0.0 {
		t.Errorf("FillRate devrait être 0.0 pour stats vides, obtenu %f", rate)
	}
}
// TestLRUCache_ResetStats teste la réinitialisation des stats
func TestLRUCache_ResetStats(t *testing.T) {
	cache := NewLRUCache(10, 0)
	cache.Set("key1", "value1")
	cache.Get("key1")
	cache.Get("key2")
	// Stats avant reset
	statsBefore := cache.GetStats()
	if statsBefore.Hits == 0 && statsBefore.Misses == 0 {
		t.Error("Stats devraient avoir des valeurs avant reset")
	}
	// Reset
	cache.ResetStats()
	// Stats après reset
	statsAfter := cache.GetStats()
	if statsAfter.Hits != 0 {
		t.Errorf("Hits devrait être 0 après reset, obtenu %d", statsAfter.Hits)
	}
	if statsAfter.Misses != 0 {
		t.Errorf("Misses devrait être 0 après reset, obtenu %d", statsAfter.Misses)
	}
	if statsAfter.Sets != 0 {
		t.Errorf("Sets devrait être 0 après reset, obtenu %d", statsAfter.Sets)
	}
	if statsAfter.Evictions != 0 {
		t.Errorf("Evictions devrait être 0 après reset, obtenu %d", statsAfter.Evictions)
	}
	// Le cache devrait toujours contenir les éléments
	if statsAfter.Size == 0 {
		t.Error("Size ne devrait pas être 0 (reset n'efface pas le cache)")
	}
}
// BenchmarkLRUCache_Set benchmark les opérations Set
func BenchmarkLRUCache_Set(b *testing.B) {
	cache := NewLRUCache(1000, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i%1000)
		cache.Set(key, i)
	}
}
// BenchmarkLRUCache_Get benchmark les opérations Get
func BenchmarkLRUCache_Get(b *testing.B) {
	cache := NewLRUCache(1000, 0)
	// Pré-remplir le cache
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key_%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i%1000)
		cache.Get(key)
	}
}
// BenchmarkLRUCache_SetGet benchmark les opérations mixtes
func BenchmarkLRUCache_SetGet(b *testing.B) {
	cache := NewLRUCache(1000, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i%1000)
		if i%2 == 0 {
			cache.Set(key, i)
		} else {
			cache.Get(key)
		}
	}
}