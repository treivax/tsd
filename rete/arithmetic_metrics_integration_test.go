// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
	"time"
)

// TestMetricsIntegrationWithCache teste l'intégration des métriques avec le cache
func TestMetricsIntegrationWithCache(t *testing.T) {
	// Créer le cache avec métriques
	cacheConfig := CacheConfig{
		MaxSize: 100,
		TTL:     5 * time.Minute,
		Enabled: true,
	}
	cache := NewArithmeticResultCache(cacheConfig)

	// Créer les métriques
	metricsConfig := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(metricsConfig)

	ruleID := "test_rule"

	// Simuler des évaluations avec cache miss puis hit
	result1 := 42.0

	// Premier accès - cache miss
	_, exists := cache.Get("temp1")
	if exists {
		t.Error("Expected cache miss on first access")
	}
	metrics.RecordCacheMiss(ruleID)
	metrics.RecordEvaluation(ruleID, true, 100*time.Microsecond)

	// Stocker dans le cache
	cache.Set("temp1", result1)

	// Deuxième accès - cache hit
	_, exists = cache.Get("temp1")
	if !exists {
		t.Error("Expected cache hit on second access")
	}
	metrics.RecordCacheHit(ruleID)
	metrics.RecordEvaluation(ruleID, true, 10*time.Microsecond) // Plus rapide grâce au cache

	// Vérifier les métriques
	rule := metrics.GetRuleMetrics(ruleID)
	if rule == nil {
		t.Fatal("Expected rule metrics")
	}

	if rule.CacheHits != 1 {
		t.Errorf("Expected 1 cache hit, got %d", rule.CacheHits)
	}

	if rule.CacheMisses != 1 {
		t.Errorf("Expected 1 cache miss, got %d", rule.CacheMisses)
	}

	if rule.CacheHitRate != 0.5 {
		t.Errorf("Expected cache hit rate 0.5, got %.2f", rule.CacheHitRate)
	}

	// Vérifier que les stats du cache correspondent
	stats := cache.GetStatistics()
	if stats.Hits != 1 {
		t.Errorf("Expected 1 cache hit in stats, got %d", stats.Hits)
	}

	if stats.Misses != 1 {
		t.Errorf("Expected 1 cache miss in stats, got %d", stats.Misses)
	}

	// Mettre à jour les métriques avec les stats du cache
	// Estimation de mémoire: ~100 bytes par entrée
	estimatedMemory := int64(stats.CurrentSize * 100)
	metrics.UpdateCacheStatistics(stats.CurrentSize, stats.Evictions, estimatedMemory)

	global := metrics.GetGlobalMetrics()
	if global.CacheSize != stats.CurrentSize {
		t.Errorf("Expected cache size %d, got %d", stats.CurrentSize, global.CacheSize)
	}
}

// TestMetricsIntegrationWithCircularDetector teste l'intégration avec le détecteur
func TestMetricsIntegrationWithCircularDetector(t *testing.T) {
	// Créer le détecteur
	detector := NewCircularDependencyDetector()

	// Créer les métriques
	metricsConfig := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(metricsConfig)

	ruleID := "test_rule"

	// Cas 1: Graphe sans cycle
	detector.AddNode("A", []string{})
	detector.AddNode("B", []string{"A"})
	detector.AddNode("C", []string{"B"})

	result := detector.Validate()
	if !result.Valid {
		t.Error("Expected valid graph without cycles")
	}

	metrics.RecordGraphValidation(result.MaxDepth, result.HasCircularDeps)

	global := metrics.GetGlobalMetrics()
	if global.GraphValidations != 1 {
		t.Errorf("Expected 1 validation, got %d", global.GraphValidations)
	}

	if global.CyclesDetected != 0 {
		t.Errorf("Expected 0 cycles, got %d", global.CyclesDetected)
	}

	if global.MaxGraphDepth != result.MaxDepth {
		t.Errorf("Expected max depth %d, got %d", result.MaxDepth, global.MaxGraphDepth)
	}

	// Cas 2: Graphe avec cycle
	detector = NewCircularDependencyDetector()
	detector.AddNode("A", []string{})
	detector.AddNode("B", []string{"A"})
	detector.AddNode("C", []string{"B"})
	detector.AddNode("A", []string{"C"}) // Cycle! A dépend maintenant de C

	result = detector.Validate()
	if result.Valid {
		t.Error("Expected invalid graph with cycle")
	}

	if !result.HasCircularDeps {
		t.Error("Expected cycle detection")
	}

	metrics.RecordGraphValidation(result.MaxDepth, result.HasCircularDeps)
	metrics.RecordCircularDependency(ruleID, result.CyclePath)

	global = metrics.GetGlobalMetrics()
	if global.GraphValidations != 2 {
		t.Errorf("Expected 2 validations, got %d", global.GraphValidations)
	}

	// Le nombre de cycles détectés peut être >= 1 car RecordGraphValidation et RecordCircularDependency
	// incrémentent tous les deux le compteur
	if global.CyclesDetected < 1 {
		t.Errorf("Expected at least 1 cycle, got %d", global.CyclesDetected)
	}

	rule := metrics.GetRuleMetrics(ruleID)
	if !rule.HasCircularDeps {
		t.Error("Expected HasCircularDeps to be true")
	}
}

// TestMetricsWithDecomposedChain teste les métriques avec une chaîne décomposée complète
func TestMetricsWithDecomposedChain(t *testing.T) {
	// Créer le cache
	cacheConfig := CacheConfig{
		MaxSize: 100,
		TTL:     5 * time.Minute,
		Enabled: true,
	}
	cache := NewArithmeticResultCache(cacheConfig)

	// Créer les métriques
	metricsConfig := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(metricsConfig)

	// Créer le détecteur
	detector := NewCircularDependencyDetector()

	ruleID := "complex_rule"

	// Simuler une décomposition: (a + b) * (c - d)
	// Étapes:
	// 1. temp1 = a + b
	// 2. temp2 = c - d
	// 3. result = temp1 * temp2

	intermediateResults := []string{"temp1", "temp2", "result"}
	dependencies := map[string][]string{
		"temp1":  {},
		"temp2":  {},
		"result": {"temp1", "temp2"},
	}

	// Valider les dépendances avec le détecteur
	for node, deps := range dependencies {
		detector.AddNode(node, deps)
	}

	validationResult := detector.Validate()
	if !validationResult.Valid {
		t.Fatal("Expected valid dependency graph")
	}

	metrics.RecordGraphValidation(validationResult.MaxDepth, validationResult.HasCircularDeps)

	// Enregistrer la structure de la chaîne
	chainLength := 3
	atomicSteps := 3 // Les 3 opérations arithmétiques
	comparisonSteps := 0

	metrics.RecordChainStructure(ruleID, chainLength, atomicSteps, comparisonSteps, intermediateResults, dependencies)

	// Simuler l'évaluation des étapes
	// Étape 1: temp1 = a + b (cache miss)
	start := time.Now()
	time.Sleep(50 * time.Microsecond)
	duration1 := time.Since(start)
	metrics.RecordCacheMiss(ruleID)
	metrics.RecordEvaluation(ruleID, true, duration1)

	// Stocker temp1 dans le cache
	cache.Set("temp1", 5.0)

	// Étape 2: temp2 = c - d (cache miss)
	start = time.Now()
	time.Sleep(50 * time.Microsecond)
	duration2 := time.Since(start)
	metrics.RecordCacheMiss(ruleID)
	metrics.RecordEvaluation(ruleID, true, duration2)

	// Stocker temp2 dans le cache
	cache.Set("temp2", 3.0)

	// Étape 3: result = temp1 * temp2 (cache hits pour temp1 et temp2)
	_, exists1 := cache.Get("temp1")
	if exists1 {
		metrics.RecordCacheHit(ruleID)
	}

	_, exists2 := cache.Get("temp2")
	if exists2 {
		metrics.RecordCacheHit(ruleID)
	}

	start = time.Now()
	time.Sleep(30 * time.Microsecond) // Plus rapide grâce au cache
	duration3 := time.Since(start)
	metrics.RecordEvaluation(ruleID, true, duration3)

	// Stocker le résultat final
	cache.Set("result", 15.0)

	// Enregistrer une activation de règle réussie
	metrics.RecordActivation(ruleID, true, duration1+duration2+duration3)

	// Vérifier les métriques de la règle
	rule := metrics.GetRuleMetrics(ruleID)
	if rule == nil {
		t.Fatal("Expected rule metrics")
	}

	if rule.ChainLength != 3 {
		t.Errorf("Expected chain length 3, got %d", rule.ChainLength)
	}

	if rule.AtomicStepsCount != 3 {
		t.Errorf("Expected 3 atomic steps, got %d", rule.AtomicStepsCount)
	}

	if rule.TotalEvaluations != 3 {
		t.Errorf("Expected 3 evaluations, got %d", rule.TotalEvaluations)
	}

	if rule.SuccessfulEvaluations != 3 {
		t.Errorf("Expected 3 successful evaluations, got %d", rule.SuccessfulEvaluations)
	}

	if rule.CacheHits != 2 {
		t.Errorf("Expected 2 cache hits, got %d", rule.CacheHits)
	}

	if rule.CacheMisses != 2 {
		t.Errorf("Expected 2 cache misses, got %d", rule.CacheMisses)
	}

	expectedHitRate := 0.5 // 2 hits / 4 total
	if rule.CacheHitRate != expectedHitRate {
		t.Errorf("Expected cache hit rate %.2f, got %.2f", expectedHitRate, rule.CacheHitRate)
	}

	if rule.MaxDependencyDepth != 1 {
		t.Errorf("Expected max dependency depth 1, got %d", rule.MaxDependencyDepth)
	}

	if rule.TotalActivations != 1 {
		t.Errorf("Expected 1 activation, got %d", rule.TotalActivations)
	}

	if rule.SuccessfulActivations != 1 {
		t.Errorf("Expected 1 successful activation, got %d", rule.SuccessfulActivations)
	}

	// Vérifier les métriques globales
	global := metrics.GetGlobalMetrics()

	if global.TotalDecomposedChains != 1 {
		t.Errorf("Expected 1 decomposed chain, got %d", global.TotalDecomposedChains)
	}

	if global.TotalAtomicNodes != 3 {
		t.Errorf("Expected 3 atomic nodes, got %d", global.TotalAtomicNodes)
	}

	if global.TotalEvaluations != 3 {
		t.Errorf("Expected 3 global evaluations, got %d", global.TotalEvaluations)
	}

	if global.TotalCacheHits != 2 {
		t.Errorf("Expected 2 global cache hits, got %d", global.TotalCacheHits)
	}

	if global.TotalCacheMisses != 2 {
		t.Errorf("Expected 2 global cache misses, got %d", global.TotalCacheMisses)
	}

	if global.GraphValidations != 1 {
		t.Errorf("Expected 1 graph validation, got %d", global.GraphValidations)
	}

	// Mettre à jour avec les stats du cache
	stats := cache.GetStatistics()
	estimatedMemory := int64(stats.CurrentSize * 100)
	metrics.UpdateCacheStatistics(stats.CurrentSize, stats.Evictions, estimatedMemory)

	global = metrics.GetGlobalMetrics()
	if global.CacheSize != 3 {
		t.Errorf("Expected cache size 3, got %d", global.CacheSize)
	}
}

// TestMetricsWithMultipleRulesAndCache teste plusieurs règles partageant le cache
func TestMetricsWithMultipleRulesAndCache(t *testing.T) {
	// Créer le cache partagé
	cacheConfig := CacheConfig{
		MaxSize: 100,
		TTL:     5 * time.Minute,
		Enabled: true,
	}
	cache := NewArithmeticResultCache(cacheConfig)

	// Créer les métriques
	metricsConfig := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(metricsConfig)

	// Règle 1: calcule temp1 = a + b
	rule1ID := "rule1"
	metrics.RecordChainStructure(rule1ID, 1, 1, 0, []string{"temp1"}, map[string][]string{"temp1": {}})

	_, exists := cache.Get("temp1")
	if exists {
		metrics.RecordCacheHit(rule1ID)
	} else {
		metrics.RecordCacheMiss(rule1ID)
		cache.Set("temp1", 10.0)
	}
	metrics.RecordEvaluation(rule1ID, true, 100*time.Microsecond)
	metrics.RecordActivation(rule1ID, true, 100*time.Microsecond)

	// Règle 2: utilise temp1 et calcule temp2 = temp1 * c
	rule2ID := "rule2"
	dependencies := map[string][]string{
		"temp1": {},
		"temp2": {"temp1"},
	}
	metrics.RecordChainStructure(rule2ID, 2, 2, 0, []string{"temp1", "temp2"}, dependencies)

	// temp1 devrait être dans le cache (hit)
	_, exists = cache.Get("temp1")
	if exists {
		metrics.RecordCacheHit(rule2ID)
	} else {
		t.Error("Expected temp1 to be in cache")
		metrics.RecordCacheMiss(rule2ID)
	}

	// temp2 n'est pas encore dans le cache (miss)
	_, exists = cache.Get("temp2")
	if exists {
		metrics.RecordCacheHit(rule2ID)
	} else {
		metrics.RecordCacheMiss(rule2ID)
		cache.Set("temp2", 20.0)
	}

	metrics.RecordEvaluation(rule2ID, true, 80*time.Microsecond) // Plus rapide grâce au cache
	metrics.RecordActivation(rule2ID, true, 80*time.Microsecond)

	// Règle 3: utilise temp1 et temp2 (tous deux en cache)
	rule3ID := "rule3"
	dependencies = map[string][]string{
		"temp1":  {},
		"temp2":  {},
		"result": {"temp1", "temp2"},
	}
	metrics.RecordChainStructure(rule3ID, 3, 3, 0, []string{"temp1", "temp2", "result"}, dependencies)

	// Les deux devraient être des hits
	_, exists = cache.Get("temp1")
	if exists {
		metrics.RecordCacheHit(rule3ID)
	}
	_, exists = cache.Get("temp2")
	if exists {
		metrics.RecordCacheHit(rule3ID)
	}

	metrics.RecordEvaluation(rule3ID, true, 50*time.Microsecond) // Très rapide
	metrics.RecordActivation(rule3ID, true, 50*time.Microsecond)

	// Vérifier les métriques de chaque règle
	rule1 := metrics.GetRuleMetrics(rule1ID)
	if rule1.CacheHits != 0 {
		t.Errorf("Rule1: Expected 0 cache hits, got %d", rule1.CacheHits)
	}
	if rule1.CacheMisses != 1 {
		t.Errorf("Rule1: Expected 1 cache miss, got %d", rule1.CacheMisses)
	}

	rule2 := metrics.GetRuleMetrics(rule2ID)
	if rule2.CacheHits != 1 {
		t.Errorf("Rule2: Expected 1 cache hit, got %d", rule2.CacheHits)
	}
	if rule2.CacheMisses != 1 {
		t.Errorf("Rule2: Expected 1 cache miss, got %d", rule2.CacheMisses)
	}

	rule3 := metrics.GetRuleMetrics(rule3ID)
	if rule3.CacheHits != 2 {
		t.Errorf("Rule3: Expected 2 cache hits, got %d", rule3.CacheHits)
	}
	if rule3.CacheMisses != 0 {
		t.Errorf("Rule3: Expected 0 cache misses, got %d", rule3.CacheMisses)
	}

	// Vérifier les métriques globales
	global := metrics.GetGlobalMetrics()

	if global.TotalRulesWithArithmetic != 3 {
		t.Errorf("Expected 3 rules with arithmetic, got %d", global.TotalRulesWithArithmetic)
	}

	if global.TotalCacheHits != 3 {
		t.Errorf("Expected 3 total cache hits, got %d", global.TotalCacheHits)
	}

	if global.TotalCacheMisses != 2 {
		t.Errorf("Expected 2 total cache misses, got %d", global.TotalCacheMisses)
	}

	expectedGlobalHitRate := 3.0 / 5.0 // 3 hits / 5 total
	if global.CacheGlobalHitRate != expectedGlobalHitRate {
		t.Errorf("Expected global hit rate %.2f, got %.2f", expectedGlobalHitRate, global.CacheGlobalHitRate)
	}

	// Vérifier que rule3 est la plus rapide (grâce au cache)
	slowestRules := metrics.GetSlowestRules(3)
	if len(slowestRules) != 3 {
		t.Fatalf("Expected 3 rules, got %d", len(slowestRules))
	}

	// rule1 devrait être la plus lente (pas de cache)
	if slowestRules[0].RuleID != rule1ID {
		t.Errorf("Expected %s to be slowest, got %s", rule1ID, slowestRules[0].RuleID)
	}

	// rule3 devrait être la plus rapide (tout en cache)
	if slowestRules[2].RuleID != rule3ID {
		t.Errorf("Expected %s to be fastest, got %s", rule3ID, slowestRules[2].RuleID)
	}
}

// TestMetricsWithCacheEviction teste les métriques lors d'évictions de cache
func TestMetricsWithCacheEviction(t *testing.T) {
	// Créer un petit cache qui forcera des évictions
	cacheConfig := CacheConfig{
		MaxSize: 2, // Seulement 2 entrées
		TTL:     5 * time.Minute,
		Enabled: true,
	}
	cache := NewArithmeticResultCache(cacheConfig)

	// Créer les métriques
	metricsConfig := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(metricsConfig)

	ruleID := "test_rule"

	// Ajouter 3 résultats (causera une éviction)
	cache.Set("temp1", 1.0)
	cache.Set("temp2", 2.0)
	cache.Set("temp3", 3.0) // Éviction de temp1

	// Mettre à jour les métriques avec les stats
	stats := cache.GetStatistics()
	estimatedMemory := int64(stats.CurrentSize * 100)
	metrics.UpdateCacheStatistics(stats.CurrentSize, stats.Evictions, estimatedMemory)

	global := metrics.GetGlobalMetrics()

	if global.CacheSize != 2 {
		t.Errorf("Expected cache size 2, got %d", global.CacheSize)
	}

	if global.CacheEvictions != 1 {
		t.Errorf("Expected 1 eviction, got %d", global.CacheEvictions)
	}

	// Tenter d'accéder à temp1 (devrait être un miss)
	_, exists := cache.Get("temp1")
	if !exists {
		metrics.RecordCacheMiss(ruleID)
	} else {
		t.Error("Expected temp1 to be evicted")
	}

	// Accéder à temp2 et temp3 (devraient être des hits)
	_, exists = cache.Get("temp2")
	if exists {
		metrics.RecordCacheHit(ruleID)
	}
	_, exists = cache.Get("temp3")
	if exists {
		metrics.RecordCacheHit(ruleID)
	}

	rule := metrics.GetRuleMetrics(ruleID)
	if rule.CacheHits != 2 {
		t.Errorf("Expected 2 cache hits, got %d", rule.CacheHits)
	}

	if rule.CacheMisses != 1 {
		t.Errorf("Expected 1 cache miss (evicted entry), got %d", rule.CacheMisses)
	}
}

// TestMetricsSummaryIntegration teste le résumé avec tous les composants
func TestMetricsSummaryIntegration(t *testing.T) {
	// Créer tous les composants
	cache := NewArithmeticResultCache(CacheConfig{
		MaxSize: 100,
		TTL:     5 * time.Minute,
		Enabled: true,
	})

	detector := NewCircularDependencyDetector()
	metrics := NewArithmeticDecompositionMetrics(DefaultMetricsConfig())

	// Simuler une session complète avec plusieurs règles
	for i := 0; i < 5; i++ {
		ruleID := "rule_" + string(rune('0'+i))

		// Enregistrer structure
		chainLength := i + 2
		atomicSteps := i + 1
		metrics.RecordChainStructure(ruleID, chainLength, atomicSteps, 1, []string{}, map[string][]string{})

		// Simuler des évaluations
		for j := 0; j < 10; j++ {
			duration := time.Duration(100+i*20) * time.Microsecond
			metrics.RecordEvaluation(ruleID, true, duration)

			// Mix de cache hits/misses
			if j%2 == 0 {
				metrics.RecordCacheHit(ruleID)
			} else {
				metrics.RecordCacheMiss(ruleID)
			}
		}

		metrics.RecordActivation(ruleID, true, time.Duration(100+i*20)*time.Microsecond)
	}

	// Valider un graphe
	detector.AddNode("A", []string{})
	detector.AddNode("B", []string{"A"})
	detector.AddNode("C", []string{"B"})
	result := detector.Validate()
	metrics.RecordGraphValidation(result.MaxDepth, result.HasCircularDeps)

	// Mettre à jour les stats du cache
	stats := cache.GetStatistics()
	estimatedMemory := int64(stats.CurrentSize * 100)
	metrics.UpdateCacheStatistics(stats.CurrentSize, stats.Evictions, estimatedMemory)

	// Obtenir le résumé
	summary := metrics.GetSummary()

	// Vérifier que le résumé contient toutes les sections
	if summary["rules"] == nil {
		t.Error("Expected rules section in summary")
	}
	if summary["nodes"] == nil {
		t.Error("Expected nodes section in summary")
	}
	if summary["evaluations"] == nil {
		t.Error("Expected evaluations section in summary")
	}
	if summary["cache"] == nil {
		t.Error("Expected cache section in summary")
	}
	if summary["validation"] == nil {
		t.Error("Expected validation section in summary")
	}

	// Vérifier les valeurs du résumé
	rules := summary["rules"].(map[string]interface{})
	if rules["tracked_rules"].(int) != 5 {
		t.Errorf("Expected 5 tracked rules, got %v", rules["tracked_rules"])
	}

	evaluations := summary["evaluations"].(map[string]interface{})
	if evaluations["total"].(int64) != 50 { // 5 rules * 10 evaluations
		t.Errorf("Expected 50 total evaluations, got %v", evaluations["total"])
	}

	cache_summary := summary["cache"].(map[string]interface{})
	if cache_summary["hits"].(int64) != 25 { // 5 rules * 5 hits each
		t.Errorf("Expected 25 cache hits, got %v", cache_summary["hits"])
	}
	if cache_summary["misses"].(int64) != 25 { // 5 rules * 5 misses each
		t.Errorf("Expected 25 cache misses, got %v", cache_summary["misses"])
	}

	validation := summary["validation"].(map[string]interface{})
	if validation["total_validations"].(int64) != 1 {
		t.Errorf("Expected 1 validation, got %v", validation["total_validations"])
	}
}
