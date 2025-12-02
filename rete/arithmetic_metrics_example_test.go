// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"time"
)

// ExampleArithmeticDecompositionMetrics_basicUsage démontre l'utilisation de base des métriques
func ExampleArithmeticDecompositionMetrics_basicUsage() {
	// Créer les métriques avec configuration par défaut
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)

	ruleID := "price_calculation"

	// Enregistrer une évaluation
	metrics.RecordEvaluation(ruleID, true, 150*time.Microsecond)
	metrics.RecordActivation(ruleID, true, 150*time.Microsecond)

	// Obtenir les métriques
	rule := metrics.GetRuleMetrics(ruleID)
	fmt.Printf("Total evaluations: %d\n", rule.TotalEvaluations)
	fmt.Printf("Average time: %v\n", rule.AvgEvaluationTime)

	// Output:
	// Total evaluations: 1
	// Average time: 150µs
}

// ExampleArithmeticDecompositionMetrics_withCache démontre l'utilisation avec le cache
func ExampleArithmeticDecompositionMetrics_withCache() {
	// Créer le cache
	cache := NewArithmeticResultCache(CacheConfig{
		MaxSize: 100,
		TTL:     5 * time.Minute,
		Enabled: true,
	})

	// Créer les métriques
	metrics := NewArithmeticDecompositionMetrics(DefaultMetricsConfig())

	ruleID := "discount_calculation"

	// Première évaluation - cache miss
	_, exists := cache.Get("discount_amount")
	if !exists {
		metrics.RecordCacheMiss(ruleID)
		// Calculer et stocker
		cache.Set("discount_amount", 15.0)
	}
	metrics.RecordEvaluation(ruleID, true, 200*time.Microsecond)

	// Deuxième évaluation - cache hit
	_, exists = cache.Get("discount_amount")
	if exists {
		metrics.RecordCacheHit(ruleID)
	}
	metrics.RecordEvaluation(ruleID, true, 10*time.Microsecond)

	// Afficher les stats
	rule := metrics.GetRuleMetrics(ruleID)
	fmt.Printf("Cache hit rate: %.0f%%\n", rule.CacheHitRate*100)

	// Output:
	// Cache hit rate: 50%
}

// ExampleArithmeticDecompositionMetrics_chainStructure démontre l'enregistrement de structure
func ExampleArithmeticDecompositionMetrics_chainStructure() {
	metrics := NewArithmeticDecompositionMetrics(DefaultMetricsConfig())

	ruleID := "complex_formula"

	// Structure: (price * quantity) - discount
	intermediateResults := []string{"subtotal", "final_amount"}
	dependencies := map[string][]string{
		"subtotal":     {},
		"final_amount": {"subtotal"},
	}

	metrics.RecordChainStructure(ruleID, 2, 2, 0, intermediateResults, dependencies)

	rule := metrics.GetRuleMetrics(ruleID)
	fmt.Printf("Chain length: %d\n", rule.ChainLength)
	fmt.Printf("Atomic steps: %d\n", rule.AtomicStepsCount)
	fmt.Printf("Max dependency depth: %d\n", rule.MaxDependencyDepth)

	// Output:
	// Chain length: 2
	// Atomic steps: 2
	// Max dependency depth: 1
}

// ExampleArithmeticDecompositionMetrics_topRules démontre le classement des règles
func ExampleArithmeticDecompositionMetrics_topRules() {
	metrics := NewArithmeticDecompositionMetrics(DefaultMetricsConfig())

	// Simuler plusieurs règles avec différentes performances
	rules := map[string]int{
		"fast_rule":   10,
		"medium_rule": 50,
		"slow_rule":   100,
	}

	for ruleID, duration := range rules {
		for i := 0; i < 5; i++ {
			metrics.RecordEvaluation(ruleID, true, time.Duration(duration)*time.Microsecond)
		}
	}

	// Obtenir les règles les plus lentes
	slowest := metrics.GetSlowestRules(2)
	fmt.Printf("Slowest rule: %s (avg: %v)\n", slowest[0].RuleID, slowest[0].AvgEvaluationTime)
	fmt.Printf("Second slowest: %s (avg: %v)\n", slowest[1].RuleID, slowest[1].AvgEvaluationTime)

	// Output:
	// Slowest rule: slow_rule (avg: 100µs)
	// Second slowest: medium_rule (avg: 50µs)
}

// ExampleArithmeticDecompositionMetrics_globalMetrics démontre les métriques globales
func ExampleArithmeticDecompositionMetrics_globalMetrics() {
	metrics := NewArithmeticDecompositionMetrics(DefaultMetricsConfig())

	// Simuler plusieurs règles
	for i := 1; i <= 3; i++ {
		ruleID := fmt.Sprintf("rule_%d", i)
		metrics.RecordChainStructure(ruleID, i+1, i, 1, []string{}, map[string][]string{})
		metrics.RecordEvaluation(ruleID, true, time.Duration(i*50)*time.Microsecond)
	}

	global := metrics.GetGlobalMetrics()
	fmt.Printf("Total rules: %d\n", global.TotalRulesWithArithmetic)
	fmt.Printf("Total chains: %d\n", global.TotalDecomposedChains)
	fmt.Printf("Total evaluations: %d\n", global.TotalEvaluations)
	fmt.Printf("Average chain length: %.1f\n", global.AverageChainLength)

	// Output:
	// Total rules: 3
	// Total chains: 3
	// Total evaluations: 3
	// Average chain length: 3.0
}

// ExampleArithmeticDecompositionMetrics_summary démontre le résumé complet
func ExampleArithmeticDecompositionMetrics_summary() {
	metrics := NewArithmeticDecompositionMetrics(DefaultMetricsConfig())

	// Ajouter quelques métriques
	ruleID := "example_rule"
	metrics.RecordChainStructure(ruleID, 3, 2, 1, []string{"temp1", "temp2"}, map[string][]string{})
	metrics.RecordEvaluation(ruleID, true, 100*time.Microsecond)
	metrics.RecordCacheHit(ruleID)
	metrics.UpdateCacheStatistics(10, 0, 1024)

	summary := metrics.GetSummary()

	rules := summary["rules"].(map[string]interface{})
	cache := summary["cache"].(map[string]interface{})

	fmt.Printf("Tracked rules: %d\n", rules["tracked_rules"])
	fmt.Printf("Cache hits: %d\n", cache["hits"])
	fmt.Printf("Cache size: %d\n", cache["size"])

	// Output:
	// Tracked rules: 1
	// Cache hits: 1
	// Cache size: 10
}

// ExampleArithmeticDecompositionMetrics_withCircularDetection démontre la détection de cycles
func ExampleArithmeticDecompositionMetrics_withCircularDetection() {
	metrics := NewArithmeticDecompositionMetrics(DefaultMetricsConfig())
	detector := NewCircularDependencyDetector()

	// Créer un graphe avec cycle
	detector.AddNode("A", []string{"B"})
	detector.AddNode("B", []string{"C"})
	detector.AddNode("C", []string{"A"}) // Cycle!

	result := detector.Validate()

	if result.HasCircularDeps {
		metrics.RecordCircularDependency("problematic_rule", result.CyclePath)
		metrics.RecordGraphValidation(result.MaxDepth, true)
	}

	global := metrics.GetGlobalMetrics()
	fmt.Printf("Cycles detected: %d\n", global.CyclesDetected)
	fmt.Printf("Total validations: %d\n", global.GraphValidations)

	// Output:
	// Cycles detected: 1
	// Total validations: 1
}

// ExampleArithmeticDecompositionMetrics_histogram démontre les histogrammes
func ExampleArithmeticDecompositionMetrics_histogram() {
	config := DefaultMetricsConfig()
	config.CollectHistograms = true
	metrics := NewArithmeticDecompositionMetrics(config)

	ruleID := "timed_rule"

	// Enregistrer des évaluations avec différents temps
	timings := []int{5, 15, 45, 150, 350}
	for _, t := range timings {
		metrics.RecordEvaluation(ruleID, true, time.Duration(t)*time.Microsecond)
	}

	rule := metrics.GetRuleMetrics(ruleID)
	fmt.Printf("Total evaluations: %d\n", rule.TotalEvaluations)
	fmt.Printf("Histogram buckets: %d\n", len(rule.EvaluationTimeHistogram))
	fmt.Printf("Min time: %v\n", rule.MinEvaluationTime)
	fmt.Printf("Max time: %v\n", rule.MaxEvaluationTime)

	// Output:
	// Total evaluations: 5
	// Histogram buckets: 5
	// Min time: 5µs
	// Max time: 350µs
}

// ExampleArithmeticDecompositionMetrics_fullIntegration démontre une intégration complète
func ExampleArithmeticDecompositionMetrics_fullIntegration() {
	// Setup complet avec tous les composants
	cache := NewArithmeticResultCache(CacheConfig{
		MaxSize: 100,
		TTL:     5 * time.Minute,
		Enabled: true,
	})

	detector := NewCircularDependencyDetector()
	metrics := NewArithmeticDecompositionMetrics(DefaultMetricsConfig())

	ruleID := "integrated_rule"

	// 1. Valider les dépendances
	dependencies := map[string][]string{
		"step1": {},
		"step2": {"step1"},
		"final": {"step2"},
	}

	for node, deps := range dependencies {
		detector.AddNode(node, deps)
	}

	validationResult := detector.Validate()
	metrics.RecordGraphValidation(validationResult.MaxDepth, validationResult.HasCircularDeps)

	// 2. Enregistrer la structure
	metrics.RecordChainStructure(ruleID, 3, 3, 0, []string{"step1", "step2", "final"}, dependencies)

	// 3. Simuler les évaluations avec cache
	// step1 - cache miss
	_, exists := cache.Get("step1")
	if !exists {
		metrics.RecordCacheMiss(ruleID)
		cache.Set("step1", 10.0)
		metrics.RecordEvaluation(ruleID, true, 100*time.Microsecond)
	}

	// step2 - cache miss, mais step1 en cache
	_, exists = cache.Get("step2")
	if !exists {
		metrics.RecordCacheMiss(ruleID)
		_, exists = cache.Get("step1")
		if exists {
			metrics.RecordCacheHit(ruleID)
		}
		cache.Set("step2", 20.0)
		metrics.RecordEvaluation(ruleID, true, 80*time.Microsecond)
	}

	// final - tout en cache
	_, exists = cache.Get("step1")
	if exists {
		metrics.RecordCacheHit(ruleID)
	}
	_, exists = cache.Get("step2")
	if exists {
		metrics.RecordCacheHit(ruleID)
	}
	metrics.RecordEvaluation(ruleID, true, 50*time.Microsecond)

	// 4. Enregistrer l'activation
	metrics.RecordActivation(ruleID, true, 230*time.Microsecond)

	// 5. Mettre à jour les stats du cache
	stats := cache.GetStatistics()
	// Estimation de mémoire: ~100 bytes par entrée
	estimatedMemory := int64(stats.CurrentSize * 100)
	metrics.UpdateCacheStatistics(stats.CurrentSize, stats.Evictions, estimatedMemory)

	// 6. Afficher le résumé
	rule := metrics.GetRuleMetrics(ruleID)
	fmt.Printf("Chain length: %d\n", rule.ChainLength)
	fmt.Printf("Total evaluations: %d\n", rule.TotalEvaluations)
	fmt.Printf("Cache hit rate: %.0f%%\n", rule.CacheHitRate*100)
	fmt.Printf("Max dependency depth: %d\n", rule.MaxDependencyDepth)

	global := metrics.GetGlobalMetrics()
	fmt.Printf("Graph validations: %d\n", global.GraphValidations)
	fmt.Printf("Cache size: %d\n", global.CacheSize)

	// Output:
	// Chain length: 3
	// Total evaluations: 3
	// Cache hit rate: 60%
	// Max dependency depth: 2
	// Graph validations: 1
	// Cache size: 3
}

// ExampleArithmeticDecompositionMetrics_performance démontre le suivi des performances
func ExampleArithmeticDecompositionMetrics_performance() {
	metrics := NewArithmeticDecompositionMetrics(DefaultMetricsConfig())

	// Simuler plusieurs règles avec des performances variées
	performanceData := map[string]struct {
		evaluations int
		avgDuration time.Duration
	}{
		"optimized_rule":    {evaluations: 1000, avgDuration: 10 * time.Microsecond},
		"standard_rule":     {evaluations: 500, avgDuration: 50 * time.Microsecond},
		"complex_rule":      {evaluations: 100, avgDuration: 200 * time.Microsecond},
		"experimental_rule": {evaluations: 50, avgDuration: 500 * time.Microsecond},
	}

	for ruleID, data := range performanceData {
		for i := 0; i < data.evaluations; i++ {
			metrics.RecordEvaluation(ruleID, true, data.avgDuration)
		}
	}

	// Analyser les performances
	fmt.Println("Top rules by evaluation count:")
	topByCount := metrics.GetTopRulesByEvaluations(3)
	for i, rule := range topByCount {
		fmt.Printf("%d. %s: %d evaluations\n", i+1, rule.RuleID, rule.TotalEvaluations)
	}

	fmt.Println("\nSlowest rules by average time:")
	slowest := metrics.GetSlowestRules(3)
	for i, rule := range slowest {
		fmt.Printf("%d. %s: %v avg\n", i+1, rule.RuleID, rule.AvgEvaluationTime)
	}

	// Output:
	// Top rules by evaluation count:
	// 1. optimized_rule: 1000 evaluations
	// 2. standard_rule: 500 evaluations
	// 3. complex_rule: 100 evaluations
	//
	// Slowest rules by average time:
	// 1. experimental_rule: 500µs avg
	// 2. complex_rule: 200µs avg
	// 3. standard_rule: 50µs avg
}
