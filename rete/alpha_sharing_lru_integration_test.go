// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"fmt"
	"testing"
	"time"
)

// TestAlphaSharingLRUIntegration_DefaultConfig teste l'intégration avec la configuration par défaut
func TestAlphaSharingLRUIntegration_DefaultConfig(t *testing.T) {
	config := DefaultChainPerformanceConfig()
	metrics := NewChainBuildMetrics()
	registry := NewAlphaSharingRegistryWithConfig(config, metrics)
	// Vérifier que le LRU cache est initialisé
	if registry.lruHashCache == nil {
		t.Fatal("LRU cache devrait être initialisé avec la configuration par défaut")
	}
	// Vérifier la capacité
	if registry.lruHashCache.Capacity() != config.HashCacheMaxSize {
		t.Errorf("Capacité LRU = %d, attendu %d", registry.lruHashCache.Capacity(), config.HashCacheMaxSize)
	}
	// Tester le caching
	condition := map[string]interface{}{
		"type":     "binaryOperation",
		"operator": "==",
		"left":     "field1",
		"right":    100,
	}
	// Premier appel - cache miss
	hash1, err := registry.ConditionHashCached(condition, "p")
	if err != nil {
		t.Fatalf("Erreur calcul hash: %v", err)
	}
	// Vérifier les métriques
	if metrics.HashCacheMisses != 1 {
		t.Errorf("Cache misses = %d, attendu 1", metrics.HashCacheMisses)
	}
	// Deuxième appel - cache hit
	hash2, err := registry.ConditionHashCached(condition, "p")
	if err != nil {
		t.Fatalf("Erreur calcul hash: %v", err)
	}
	if hash1 != hash2 {
		t.Errorf("Les hash devraient être identiques: %s != %s", hash1, hash2)
	}
	if metrics.HashCacheHits != 1 {
		t.Errorf("Cache hits = %d, attendu 1", metrics.HashCacheHits)
	}
	// Vérifier la taille du cache
	if registry.GetHashCacheSize() != 1 {
		t.Errorf("Taille cache = %d, attendu 1", registry.GetHashCacheSize())
	}
	// Vérifier les statistiques détaillées
	stats := registry.GetHashCacheStats()
	if stats["type"] != "lru" {
		t.Errorf("Type cache = %v, attendu 'lru'", stats["type"])
	}
	if stats["size"] != 1 {
		t.Errorf("Taille cache stats = %v, attendu 1", stats["size"])
	}
}

// TestAlphaSharingLRUIntegration_HighPerformance teste avec la config haute performance
func TestAlphaSharingLRUIntegration_HighPerformance(t *testing.T) {
	config := HighPerformanceConfig()
	metrics := NewChainBuildMetrics()
	registry := NewAlphaSharingRegistryWithConfig(config, metrics)
	if registry.lruHashCache == nil {
		t.Fatal("LRU cache devrait être initialisé")
	}
	expectedCapacity := 100000
	if registry.lruHashCache.Capacity() != expectedCapacity {
		t.Errorf("Capacité = %d, attendu %d", registry.lruHashCache.Capacity(), expectedCapacity)
	}
	// Tester avec beaucoup de conditions
	numConditions := 1000
	hashes := make([]string, numConditions)
	for i := 0; i < numConditions; i++ {
		condition := map[string]interface{}{
			"type":     "binaryOperation",
			"operator": "==",
			"left":     fmt.Sprintf("field%d", i),
			"right":    i,
		}
		hash, err := registry.ConditionHashCached(condition, "p")
		if err != nil {
			t.Fatalf("Erreur calcul hash %d: %v", i, err)
		}
		hashes[i] = hash
	}
	// Vérifier que tous les hash sont en cache
	if registry.GetHashCacheSize() != numConditions {
		t.Errorf("Taille cache = %d, attendu %d", registry.GetHashCacheSize(), numConditions)
	}
	// Re-calculer les hash - tous devraient être des hits
	initialHits := metrics.HashCacheHits
	for i := 0; i < numConditions; i++ {
		condition := map[string]interface{}{
			"type":     "binaryOperation",
			"operator": "==",
			"left":     fmt.Sprintf("field%d", i),
			"right":    i,
		}
		hash, err := registry.ConditionHashCached(condition, "p")
		if err != nil {
			t.Fatalf("Erreur recalcul hash %d: %v", i, err)
		}
		if hash != hashes[i] {
			t.Errorf("Hash %d différent: %s != %s", i, hash, hashes[i])
		}
	}
	expectedHits := initialHits + int(numConditions)
	if int(metrics.HashCacheHits) != expectedHits {
		t.Errorf("Cache hits = %d, attendu %d", metrics.HashCacheHits, expectedHits)
	}
	// Vérifier le hit rate
	stats := registry.GetHashCacheStats()
	hitRate := stats["hit_rate"].(float64)
	if hitRate < 0.5 {
		t.Errorf("Hit rate trop faible: %.2f", hitRate)
	}
}

// TestAlphaSharingLRUIntegration_LRUEviction teste l'éviction LRU
func TestAlphaSharingLRUIntegration_LRUEviction(t *testing.T) {
	// Configuration avec petite capacité pour tester l'éviction
	config := DefaultChainPerformanceConfig()
	config.HashCacheMaxSize = 100
	metrics := NewChainBuildMetrics()
	registry := NewAlphaSharingRegistryWithConfig(config, metrics)
	// Ajouter 150 conditions (plus que la capacité)
	numConditions := 150
	for i := 0; i < numConditions; i++ {
		condition := map[string]interface{}{
			"type":     "binaryOperation",
			"operator": "==",
			"left":     fmt.Sprintf("field%d", i),
			"right":    i,
		}
		_, err := registry.ConditionHashCached(condition, "p")
		if err != nil {
			t.Fatalf("Erreur calcul hash %d: %v", i, err)
		}
	}
	// La taille du cache devrait être limitée à 100
	cacheSize := registry.GetHashCacheSize()
	if cacheSize != 100 {
		t.Errorf("Taille cache = %d, attendu 100 (capacité max)", cacheSize)
	}
	// Vérifier qu'il y a eu des évictions
	stats := registry.GetHashCacheStats()
	evictions := stats["evictions"].(int64)
	if evictions != 50 {
		t.Errorf("Évictions = %d, attendu 50", evictions)
	}
	// Les premières conditions devraient avoir été évincées
	// Re-calculer les 10 premières - devraient être des misses
	initialMisses := metrics.HashCacheMisses
	for i := 0; i < 10; i++ {
		condition := map[string]interface{}{
			"type":     "binaryOperation",
			"operator": "==",
			"left":     fmt.Sprintf("field%d", i),
			"right":    i,
		}
		_, err := registry.ConditionHashCached(condition, "p")
		if err != nil {
			t.Fatalf("Erreur recalcul hash %d: %v", i, err)
		}
	}
	// Devrait avoir au moins quelques misses (conditions évincées)
	if metrics.HashCacheMisses <= initialMisses {
		t.Errorf("Pas assez de misses détectés après éviction: %d", metrics.HashCacheMisses-initialMisses)
	}
}

// TestAlphaSharingLRUIntegration_TTLExpiration teste l'expiration TTL
func TestAlphaSharingLRUIntegration_TTLExpiration(t *testing.T) {
	// Configuration avec TTL court
	config := DefaultChainPerformanceConfig()
	config.HashCacheTTL = 100 * time.Millisecond
	metrics := NewChainBuildMetrics()
	registry := NewAlphaSharingRegistryWithConfig(config, metrics)
	condition := map[string]interface{}{
		"type":     "binaryOperation",
		"operator": "==",
		"left":     "field1",
		"right":    100,
	}
	// Premier calcul
	hash1, err := registry.ConditionHashCached(condition, "p")
	if err != nil {
		t.Fatalf("Erreur calcul hash: %v", err)
	}
	// Immédiatement après - devrait être un hit
	hash2, err := registry.ConditionHashCached(condition, "p")
	if err != nil {
		t.Fatalf("Erreur recalcul hash: %v", err)
	}
	if hash1 != hash2 {
		t.Errorf("Hash différents: %s != %s", hash1, hash2)
	}
	if metrics.HashCacheHits != 1 {
		t.Errorf("Cache hits = %d, attendu 1", metrics.HashCacheHits)
	}
	// Attendre l'expiration
	time.Sleep(150 * time.Millisecond)
	// Après expiration - devrait être un miss
	initialMisses := metrics.HashCacheMisses
	hash3, err := registry.ConditionHashCached(condition, "p")
	if err != nil {
		t.Fatalf("Erreur calcul hash après expiration: %v", err)
	}
	if hash3 != hash1 {
		t.Errorf("Hash après expiration différent: %s != %s", hash3, hash1)
	}
	if metrics.HashCacheMisses != initialMisses+1 {
		t.Errorf("Cache misses = %d, attendu %d", metrics.HashCacheMisses, initialMisses+1)
	}
}

// TestAlphaSharingLRUIntegration_CleanExpired teste le nettoyage des entrées expirées
func TestAlphaSharingLRUIntegration_CleanExpired(t *testing.T) {
	config := DefaultChainPerformanceConfig()
	config.HashCacheTTL = 50 * time.Millisecond
	metrics := NewChainBuildMetrics()
	registry := NewAlphaSharingRegistryWithConfig(config, metrics)
	// Ajouter plusieurs conditions
	numConditions := 10
	for i := 0; i < numConditions; i++ {
		condition := map[string]interface{}{
			"type":     "binaryOperation",
			"operator": "==",
			"left":     fmt.Sprintf("field%d", i),
			"right":    i,
		}
		_, err := registry.ConditionHashCached(condition, "p")
		if err != nil {
			t.Fatalf("Erreur calcul hash %d: %v", i, err)
		}
	}
	// Vérifier la taille
	if registry.GetHashCacheSize() != numConditions {
		t.Errorf("Taille cache = %d, attendu %d", registry.GetHashCacheSize(), numConditions)
	}
	// Attendre l'expiration
	time.Sleep(100 * time.Millisecond)
	// Nettoyer les entrées expirées
	cleaned := registry.CleanExpiredHashCache()
	if cleaned != numConditions {
		t.Errorf("Entrées nettoyées = %d, attendu %d", cleaned, numConditions)
	}
	// Le cache devrait être vide
	if registry.GetHashCacheSize() != 0 {
		t.Errorf("Taille cache après nettoyage = %d, attendu 0", registry.GetHashCacheSize())
	}
}

// TestAlphaSharingLRUIntegration_DisabledCache teste avec cache désactivé
func TestAlphaSharingLRUIntegration_DisabledCache(t *testing.T) {
	config := DisabledCachesConfig()
	metrics := NewChainBuildMetrics()
	registry := NewAlphaSharingRegistryWithConfig(config, metrics)
	// Le cache ne devrait pas être initialisé
	if registry.lruHashCache != nil {
		t.Error("LRU cache devrait être nil avec cache désactivé")
	}
	condition := map[string]interface{}{
		"type":     "binaryOperation",
		"operator": "==",
		"left":     "field1",
		"right":    100,
	}
	// Calculer plusieurs fois
	hash1, err := registry.ConditionHashCached(condition, "p")
	if err != nil {
		t.Fatalf("Erreur calcul hash: %v", err)
	}
	hash2, err := registry.ConditionHashCached(condition, "p")
	if err != nil {
		t.Fatalf("Erreur recalcul hash: %v", err)
	}
	// Les hash devraient être identiques
	if hash1 != hash2 {
		t.Errorf("Hash différents: %s != %s", hash1, hash2)
	}
	// Pas de hits/misses car pas de cache
	if metrics.HashCacheHits != 0 {
		t.Errorf("Cache hits devrait être 0, obtenu %d", metrics.HashCacheHits)
	}
	if metrics.HashCacheMisses != 0 {
		t.Errorf("Cache misses devrait être 0, obtenu %d", metrics.HashCacheMisses)
	}
}

// TestAlphaSharingLRUIntegration_ClearCache teste le vidage du cache
func TestAlphaSharingLRUIntegration_ClearCache(t *testing.T) {
	config := DefaultChainPerformanceConfig()
	metrics := NewChainBuildMetrics()
	registry := NewAlphaSharingRegistryWithConfig(config, metrics)
	// Ajouter des conditions
	for i := 0; i < 10; i++ {
		condition := map[string]interface{}{
			"type":     "binaryOperation",
			"operator": "==",
			"left":     fmt.Sprintf("field%d", i),
			"right":    i,
		}
		_, err := registry.ConditionHashCached(condition, "p")
		if err != nil {
			t.Fatalf("Erreur calcul hash %d: %v", i, err)
		}
	}
	// Vérifier que le cache contient des éléments
	if registry.GetHashCacheSize() == 0 {
		t.Error("Cache devrait contenir des éléments")
	}
	// Vider le cache
	registry.ClearHashCache()
	// Le cache devrait être vide
	if registry.GetHashCacheSize() != 0 {
		t.Errorf("Taille cache après Clear = %d, attendu 0", registry.GetHashCacheSize())
	}
	// Les statistiques LRU devraient refléter le vidage
	stats := registry.GetHashCacheStats()
	if stats["size"] != 0 {
		t.Errorf("Stats size = %v, attendu 0", stats["size"])
	}
}

// TestAlphaSharingLRUIntegration_ReteNetwork teste l'intégration avec ReteNetwork
func TestAlphaSharingLRUIntegration_ReteNetwork(t *testing.T) {
	storage := NewMemoryStorage()
	config := HighPerformanceConfig()
	// Créer un réseau avec configuration
	network := NewReteNetworkWithConfig(storage, config)
	// Vérifier que la configuration est appliquée
	if network.Config == nil {
		t.Fatal("Configuration devrait être définie")
	}
	// Vérifier que l'AlphaSharingManager utilise le LRU
	if network.AlphaSharingManager.lruHashCache == nil {
		t.Error("AlphaSharingManager devrait utiliser le cache LRU")
	}
	// Vérifier la capacité
	expectedCapacity := config.HashCacheMaxSize
	actualCapacity := network.AlphaSharingManager.lruHashCache.Capacity()
	if actualCapacity != expectedCapacity {
		t.Errorf("Capacité cache = %d, attendu %d", actualCapacity, expectedCapacity)
	}
	// Tester l'utilisation via le builder de chaînes
	builder := NewAlphaChainBuilder(network, storage)
	// Créer plusieurs conditions identiques
	typeDef := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "age", Type: "int"},
		},
	}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode
	for i := 0; i < 5; i++ {
		condition := SimpleCondition{
			Type:     "binaryOperation",
			Operator: "==",
			Left:     "age",
			Right:    25,
		}
		chain, err := builder.BuildChain([]SimpleCondition{condition}, "p", typeNode, fmt.Sprintf("rule_%d", i))
		if err != nil {
			t.Fatalf("Erreur construction chaîne %d: %v", i, err)
		}
		if len(chain.Nodes) == 0 {
			t.Errorf("Chaîne %d vide", i)
		}
	}
	// Vérifier qu'il y a eu du partage (cache hits)
	metrics := network.GetChainMetrics()
	if metrics.HashCacheHits == 0 {
		t.Error("Devrait avoir des cache hits avec conditions identiques")
	}
	// Vérifier les stats du cache
	stats := network.AlphaSharingManager.GetHashCacheStats()
	if stats["size"].(int) == 0 {
		t.Error("Cache devrait contenir des entrées")
	}
	hitRate := stats["hit_rate"].(float64)
	if hitRate == 0 {
		t.Error("Hit rate devrait être > 0")
	}
}

// TestAlphaSharingLRUIntegration_LowMemoryConfig teste la config basse mémoire
func TestAlphaSharingLRUIntegration_LowMemoryConfig(t *testing.T) {
	config := LowMemoryConfig()
	metrics := NewChainBuildMetrics()
	registry := NewAlphaSharingRegistryWithConfig(config, metrics)
	// Vérifier la capacité réduite
	expectedCapacity := 1000
	if registry.lruHashCache.Capacity() != expectedCapacity {
		t.Errorf("Capacité = %d, attendu %d", registry.lruHashCache.Capacity(), expectedCapacity)
	}
	// Ajouter plus que la capacité
	numConditions := 1500
	for i := 0; i < numConditions; i++ {
		condition := map[string]interface{}{
			"type":     "binaryOperation",
			"operator": "==",
			"left":     fmt.Sprintf("field%d", i),
			"right":    i,
		}
		_, err := registry.ConditionHashCached(condition, "p")
		if err != nil {
			t.Fatalf("Erreur calcul hash %d: %v", i, err)
		}
	}
	// La taille devrait être limitée
	cacheSize := registry.GetHashCacheSize()
	if cacheSize > expectedCapacity {
		t.Errorf("Taille cache = %d, devrait être <= %d", cacheSize, expectedCapacity)
	}
	// Devrait y avoir beaucoup d'évictions
	stats := registry.GetHashCacheStats()
	evictions := stats["evictions"].(int64)
	minExpectedEvictions := int64(numConditions - expectedCapacity)
	if evictions < minExpectedEvictions {
		t.Errorf("Évictions = %d, attendu au moins %d", evictions, minExpectedEvictions)
	}
}

// TestAlphaSharingLRUIntegration_ConcurrentAccess teste l'accès concurrent
func TestAlphaSharingLRUIntegration_ConcurrentAccess(t *testing.T) {
	config := DefaultChainPerformanceConfig()
	metrics := NewChainBuildMetrics()
	registry := NewAlphaSharingRegistryWithConfig(config, metrics)
	// Nombre de goroutines concurrentes
	numGoroutines := 10
	numConditionsPerGoroutine := 100
	// Canal pour synchroniser la fin
	done := make(chan bool, numGoroutines)
	// Lancer plusieurs goroutines qui calculent des hash
	for g := 0; g < numGoroutines; g++ {
		go func(goroutineID int) {
			for i := 0; i < numConditionsPerGoroutine; i++ {
				condition := map[string]interface{}{
					"type":     "binaryOperation",
					"operator": "==",
					"left":     fmt.Sprintf("field%d_%d", goroutineID, i),
					"right":    i,
				}
				_, err := registry.ConditionHashCached(condition, "p")
				if err != nil {
					t.Errorf("Erreur calcul hash goroutine %d, condition %d: %v", goroutineID, i, err)
				}
				// Re-calculer immédiatement pour tester les hits
				_, err = registry.ConditionHashCached(condition, "p")
				if err != nil {
					t.Errorf("Erreur recalcul hash goroutine %d, condition %d: %v", goroutineID, i, err)
				}
			}
			done <- true
		}(g)
	}
	// Attendre que toutes les goroutines terminent
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	// Vérifier qu'il y a eu des hits (pas de race conditions)
	if metrics.HashCacheHits == 0 {
		t.Error("Devrait avoir des cache hits après accès concurrent")
	}
	// Vérifier la cohérence du cache
	cacheSize := registry.GetHashCacheSize()
	if cacheSize == 0 {
		t.Error("Cache ne devrait pas être vide")
	}
	stats := registry.GetHashCacheStats()
	totalAccess := stats["hits"].(int64) + stats["misses"].(int64)
	expectedAccess := int64(numGoroutines * numConditionsPerGoroutine * 2) // *2 car on calcule deux fois
	if totalAccess != expectedAccess {
		t.Errorf("Total accès = %d, attendu %d", totalAccess, expectedAccess)
	}
}
