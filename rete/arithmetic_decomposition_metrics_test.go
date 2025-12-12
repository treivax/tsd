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
// TestNewArithmeticDecompositionMetrics teste la création de métriques
func TestNewArithmeticDecompositionMetrics(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	if metrics == nil {
		t.Fatal("Expected non-nil metrics")
	}
	if metrics.ruleMetrics == nil {
		t.Error("Expected initialized ruleMetrics map")
	}
	if !metrics.config.Enabled {
		t.Error("Expected metrics to be enabled by default")
	}
	if len(metrics.config.HistogramBuckets) == 0 {
		t.Error("Expected histogram buckets to be configured")
	}
}
// TestRecordActivation teste l'enregistrement d'activations
func TestRecordActivation(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	ruleID := "rule1"
	// Enregistrer des activations réussies
	metrics.RecordActivation(ruleID, true, 100*time.Microsecond)
	metrics.RecordActivation(ruleID, true, 150*time.Microsecond)
	metrics.RecordActivation(ruleID, false, 200*time.Microsecond)
	rule := metrics.GetRuleMetrics(ruleID)
	if rule == nil {
		t.Fatal("Expected rule metrics to be created")
	}
	if rule.TotalActivations != 3 {
		t.Errorf("Expected 3 total activations, got %d", rule.TotalActivations)
	}
	if rule.SuccessfulActivations != 2 {
		t.Errorf("Expected 2 successful activations, got %d", rule.SuccessfulActivations)
	}
	if rule.FailedActivations != 1 {
		t.Errorf("Expected 1 failed activation, got %d", rule.FailedActivations)
	}
	global := metrics.GetGlobalMetrics()
	if global.TotalActivations != 3 {
		t.Errorf("Expected 3 global activations, got %d", global.TotalActivations)
	}
}
// TestRecordEvaluation teste l'enregistrement d'évaluations
func TestRecordEvaluation(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	ruleID := "rule1"
	// Enregistrer des évaluations avec différents temps
	durations := []time.Duration{
		100 * time.Microsecond,
		200 * time.Microsecond,
		150 * time.Microsecond,
		300 * time.Microsecond,
	}
	for _, d := range durations {
		metrics.RecordEvaluation(ruleID, true, d)
	}
	rule := metrics.GetRuleMetrics(ruleID)
	if rule == nil {
		t.Fatal("Expected rule metrics to be created")
	}
	if rule.TotalEvaluations != 4 {
		t.Errorf("Expected 4 total evaluations, got %d", rule.TotalEvaluations)
	}
	if rule.SuccessfulEvaluations != 4 {
		t.Errorf("Expected 4 successful evaluations, got %d", rule.SuccessfulEvaluations)
	}
	if rule.MinEvaluationTime != 100*time.Microsecond {
		t.Errorf("Expected min time 100µs, got %v", rule.MinEvaluationTime)
	}
	if rule.MaxEvaluationTime != 300*time.Microsecond {
		t.Errorf("Expected max time 300µs, got %v", rule.MaxEvaluationTime)
	}
	expectedTotal := time.Duration(100+200+150+300) * time.Microsecond
	if rule.TotalEvaluationTime != expectedTotal {
		t.Errorf("Expected total time %v, got %v", expectedTotal, rule.TotalEvaluationTime)
	}
	expectedAvg := expectedTotal / 4
	if rule.AvgEvaluationTime != expectedAvg {
		t.Errorf("Expected avg time %v, got %v", expectedAvg, rule.AvgEvaluationTime)
	}
}
// TestRecordEvaluationHistogram teste l'histogramme des temps
func TestRecordEvaluationHistogram(t *testing.T) {
	config := DefaultMetricsConfig()
	config.CollectHistograms = true
	metrics := NewArithmeticDecompositionMetrics(config)
	ruleID := "rule1"
	// Enregistrer des évaluations qui tombent dans différents buckets
	metrics.RecordEvaluation(ruleID, true, 3*time.Microsecond)   // bucket 5
	metrics.RecordEvaluation(ruleID, true, 8*time.Microsecond)   // bucket 10
	metrics.RecordEvaluation(ruleID, true, 8*time.Microsecond)   // bucket 10
	metrics.RecordEvaluation(ruleID, true, 30*time.Microsecond)  // bucket 50
	metrics.RecordEvaluation(ruleID, true, 200*time.Microsecond) // bucket 250
	rule := metrics.GetRuleMetrics(ruleID)
	if rule == nil {
		t.Fatal("Expected rule metrics to be created")
	}
	if len(rule.EvaluationTimeHistogram) == 0 {
		t.Error("Expected histogram to have entries")
	}
	// Vérifier les buckets attendus
	if rule.EvaluationTimeHistogram[5] != 1 {
		t.Errorf("Expected bucket 5 to have 1 entry, got %d", rule.EvaluationTimeHistogram[5])
	}
	if rule.EvaluationTimeHistogram[10] != 2 {
		t.Errorf("Expected bucket 10 to have 2 entries, got %d", rule.EvaluationTimeHistogram[10])
	}
	if rule.EvaluationTimeHistogram[50] != 1 {
		t.Errorf("Expected bucket 50 to have 1 entry, got %d", rule.EvaluationTimeHistogram[50])
	}
	if rule.EvaluationTimeHistogram[250] != 1 {
		t.Errorf("Expected bucket 250 to have 1 entry, got %d", rule.EvaluationTimeHistogram[250])
	}
}
// TestRecordCacheHitMiss teste l'enregistrement des cache hits/misses
func TestRecordCacheHitMiss(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	ruleID := "rule1"
	// Enregistrer des hits et misses
	metrics.RecordCacheHit(ruleID)
	metrics.RecordCacheHit(ruleID)
	metrics.RecordCacheHit(ruleID)
	metrics.RecordCacheMiss(ruleID)
	rule := metrics.GetRuleMetrics(ruleID)
	if rule == nil {
		t.Fatal("Expected rule metrics to be created")
	}
	if rule.CacheHits != 3 {
		t.Errorf("Expected 3 cache hits, got %d", rule.CacheHits)
	}
	if rule.CacheMisses != 1 {
		t.Errorf("Expected 1 cache miss, got %d", rule.CacheMisses)
	}
	expectedHitRate := 3.0 / 4.0
	if rule.CacheHitRate != expectedHitRate {
		t.Errorf("Expected cache hit rate %.2f, got %.2f", expectedHitRate, rule.CacheHitRate)
	}
	if !rule.CacheEnabled {
		t.Error("Expected cache to be marked as enabled")
	}
	global := metrics.GetGlobalMetrics()
	if global.TotalCacheHits != 3 {
		t.Errorf("Expected 3 global cache hits, got %d", global.TotalCacheHits)
	}
	if global.TotalCacheMisses != 1 {
		t.Errorf("Expected 1 global cache miss, got %d", global.TotalCacheMisses)
	}
	if global.CacheGlobalHitRate != expectedHitRate {
		t.Errorf("Expected global hit rate %.2f, got %.2f", expectedHitRate, global.CacheGlobalHitRate)
	}
}
// TestRecordChainStructure teste l'enregistrement de la structure de chaîne
func TestRecordChainStructure(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	ruleID := "rule1"
	intermediateResults := []string{"temp1", "temp2", "temp3"}
	dependencies := map[string][]string{
		"temp1": {},
		"temp2": {"temp1"},
		"temp3": {"temp1", "temp2"},
	}
	metrics.RecordChainStructure(ruleID, 5, 3, 2, intermediateResults, dependencies)
	rule := metrics.GetRuleMetrics(ruleID)
	if rule == nil {
		t.Fatal("Expected rule metrics to be created")
	}
	if rule.ChainLength != 5 {
		t.Errorf("Expected chain length 5, got %d", rule.ChainLength)
	}
	if rule.AtomicStepsCount != 3 {
		t.Errorf("Expected 3 atomic steps, got %d", rule.AtomicStepsCount)
	}
	if rule.ComparisonStepsCount != 2 {
		t.Errorf("Expected 2 comparison steps, got %d", rule.ComparisonStepsCount)
	}
	if len(rule.IntermediateResults) != 3 {
		t.Errorf("Expected 3 intermediate results, got %d", len(rule.IntermediateResults))
	}
	if len(rule.Dependencies) != 3 {
		t.Errorf("Expected 3 dependencies entries, got %d", len(rule.Dependencies))
	}
	// La profondeur max devrait être 2 (temp3 dépend de temp2 qui dépend de temp1)
	if rule.MaxDependencyDepth != 2 {
		t.Errorf("Expected max depth 2, got %d", rule.MaxDependencyDepth)
	}
	global := metrics.GetGlobalMetrics()
	if global.TotalDecomposedChains != 1 {
		t.Errorf("Expected 1 decomposed chain, got %d", global.TotalDecomposedChains)
	}
	if global.TotalAtomicNodes != 3 {
		t.Errorf("Expected 3 atomic nodes, got %d", global.TotalAtomicNodes)
	}
	if global.TotalComparisonNodes != 2 {
		t.Errorf("Expected 2 comparison nodes, got %d", global.TotalComparisonNodes)
	}
}
// TestRecordCircularDependency teste l'enregistrement de dépendances circulaires
func TestRecordCircularDependency(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	ruleID := "rule1"
	cyclePath := []string{"A", "B", "C", "A"}
	metrics.RecordCircularDependency(ruleID, cyclePath)
	rule := metrics.GetRuleMetrics(ruleID)
	if rule == nil {
		t.Fatal("Expected rule metrics to be created")
	}
	if !rule.HasCircularDeps {
		t.Error("Expected HasCircularDeps to be true")
	}
	if rule.Metadata["cycle_path"] == nil {
		t.Error("Expected cycle_path in metadata")
	}
	global := metrics.GetGlobalMetrics()
	if global.TotalCircularDepsDetected != 1 {
		t.Errorf("Expected 1 circular dependency detected, got %d", global.TotalCircularDepsDetected)
	}
	if global.CyclesDetected != 1 {
		t.Errorf("Expected 1 cycle detected, got %d", global.CyclesDetected)
	}
}
// TestRecordGraphValidation teste l'enregistrement de validations de graphe
func TestRecordGraphValidation(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	metrics.RecordGraphValidation(5, false)
	metrics.RecordGraphValidation(8, true)
	metrics.RecordGraphValidation(3, false)
	global := metrics.GetGlobalMetrics()
	if global.GraphValidations != 3 {
		t.Errorf("Expected 3 graph validations, got %d", global.GraphValidations)
	}
	if global.CyclesDetected != 1 {
		t.Errorf("Expected 1 cycle detected, got %d", global.CyclesDetected)
	}
	if global.MaxGraphDepth != 8 {
		t.Errorf("Expected max graph depth 8, got %d", global.MaxGraphDepth)
	}
}
// TestUpdateCacheStatistics teste la mise à jour des stats de cache
func TestUpdateCacheStatistics(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	metrics.UpdateCacheStatistics(100, 5, 1024*1024)
	global := metrics.GetGlobalMetrics()
	if global.CacheSize != 100 {
		t.Errorf("Expected cache size 100, got %d", global.CacheSize)
	}
	if global.CacheEvictions != 5 {
		t.Errorf("Expected 5 cache evictions, got %d", global.CacheEvictions)
	}
	if global.CacheMemoryUsage != 1024*1024 {
		t.Errorf("Expected cache memory usage 1MB, got %d", global.CacheMemoryUsage)
	}
}
// TestGetAllRuleMetrics teste la récupération de toutes les métriques
func TestGetAllRuleMetrics(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	// Créer plusieurs règles
	ruleIDs := []string{"rule1", "rule2", "rule3"}
	for _, ruleID := range ruleIDs {
		metrics.RecordActivation(ruleID, true, 100*time.Microsecond)
		metrics.RecordEvaluation(ruleID, true, 200*time.Microsecond)
	}
	allMetrics := metrics.GetAllRuleMetrics()
	if len(allMetrics) != 3 {
		t.Errorf("Expected 3 rules, got %d", len(allMetrics))
	}
	for _, ruleID := range ruleIDs {
		if _, exists := allMetrics[ruleID]; !exists {
			t.Errorf("Expected metrics for %s", ruleID)
		}
	}
}
// TestGetTopRulesByEvaluations teste le classement par nombre d'évaluations
func TestGetTopRulesByEvaluations(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	// Créer des règles avec différents nombres d'évaluations
	metrics.RecordEvaluation("rule1", true, 100*time.Microsecond)
	metrics.RecordEvaluation("rule2", true, 100*time.Microsecond)
	metrics.RecordEvaluation("rule2", true, 100*time.Microsecond)
	metrics.RecordEvaluation("rule2", true, 100*time.Microsecond)
	metrics.RecordEvaluation("rule3", true, 100*time.Microsecond)
	metrics.RecordEvaluation("rule3", true, 100*time.Microsecond)
	topRules := metrics.GetTopRulesByEvaluations(2)
	if len(topRules) != 2 {
		t.Errorf("Expected 2 top rules, got %d", len(topRules))
	}
	if topRules[0].RuleID != "rule2" {
		t.Errorf("Expected rule2 to be first, got %s", topRules[0].RuleID)
	}
	if topRules[1].RuleID != "rule3" {
		t.Errorf("Expected rule3 to be second, got %s", topRules[1].RuleID)
	}
}
// TestGetTopRulesByDuration teste le classement par durée totale
func TestGetTopRulesByDuration(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	// Créer des règles avec différentes durées
	metrics.RecordEvaluation("rule1", true, 1000*time.Microsecond)
	metrics.RecordEvaluation("rule2", true, 500*time.Microsecond)
	metrics.RecordEvaluation("rule2", true, 1000*time.Microsecond)
	metrics.RecordEvaluation("rule3", true, 100*time.Microsecond)
	topRules := metrics.GetTopRulesByDuration(2)
	if len(topRules) != 2 {
		t.Errorf("Expected 2 top rules, got %d", len(topRules))
	}
	if topRules[0].RuleID != "rule2" {
		t.Errorf("Expected rule2 to be first (1500µs total), got %s", topRules[0].RuleID)
	}
	if topRules[1].RuleID != "rule1" {
		t.Errorf("Expected rule1 to be second (1000µs total), got %s", topRules[1].RuleID)
	}
}
// TestGetSlowestRules teste le classement par temps moyen
func TestGetSlowestRules(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	// rule1: 1 évaluation de 1000µs = 1000µs moyen
	metrics.RecordEvaluation("rule1", true, 1000*time.Microsecond)
	// rule2: 2 évaluations de 500µs = 500µs moyen
	metrics.RecordEvaluation("rule2", true, 500*time.Microsecond)
	metrics.RecordEvaluation("rule2", true, 500*time.Microsecond)
	// rule3: 4 évaluations de 100µs = 100µs moyen
	metrics.RecordEvaluation("rule3", true, 100*time.Microsecond)
	metrics.RecordEvaluation("rule3", true, 100*time.Microsecond)
	metrics.RecordEvaluation("rule3", true, 100*time.Microsecond)
	metrics.RecordEvaluation("rule3", true, 100*time.Microsecond)
	slowestRules := metrics.GetSlowestRules(2)
	if len(slowestRules) != 2 {
		t.Errorf("Expected 2 slowest rules, got %d", len(slowestRules))
	}
	if slowestRules[0].RuleID != "rule1" {
		t.Errorf("Expected rule1 to be slowest, got %s", slowestRules[0].RuleID)
	}
	if slowestRules[1].RuleID != "rule2" {
		t.Errorf("Expected rule2 to be second slowest, got %s", slowestRules[1].RuleID)
	}
}
// TestArithmeticMetricsGetSummary teste le résumé des métriques
func TestArithmeticMetricsGetSummary(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	// Ajouter quelques métriques
	metrics.RecordActivation("rule1", true, 100*time.Microsecond)
	metrics.RecordEvaluation("rule1", true, 200*time.Microsecond)
	metrics.RecordChainStructure("rule1", 5, 3, 2, []string{"temp1"}, map[string][]string{})
	metrics.RecordCacheHit("rule1")
	metrics.RecordGraphValidation(5, false)
	summary := metrics.GetSummary()
	if summary == nil {
		t.Fatal("Expected non-nil summary")
	}
	rules := summary["rules"].(map[string]interface{})
	if rules["total_decomposed_chains"].(int) != 1 {
		t.Error("Expected 1 decomposed chain in summary")
	}
	cache := summary["cache"].(map[string]interface{})
	if cache["hits"].(int64) != 1 {
		t.Error("Expected 1 cache hit in summary")
	}
	validation := summary["validation"].(map[string]interface{})
	if validation["total_validations"].(int64) != 1 {
		t.Error("Expected 1 validation in summary")
	}
}
// TestArithmeticMetricsReset teste la réinitialisation des métriques
func TestArithmeticMetricsReset(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	// Ajouter des métriques
	metrics.RecordActivation("rule1", true, 100*time.Microsecond)
	metrics.RecordEvaluation("rule1", true, 200*time.Microsecond)
	metrics.RecordCacheHit("rule1")
	// Vérifier qu'elles existent
	if metrics.GetRuleMetrics("rule1") == nil {
		t.Fatal("Expected rule metrics before reset")
	}
	global := metrics.GetGlobalMetrics()
	if global.TotalActivations == 0 {
		t.Fatal("Expected global metrics before reset")
	}
	// Réinitialiser
	metrics.Reset()
	// Vérifier que tout est réinitialisé
	if metrics.GetRuleMetrics("rule1") != nil {
		t.Error("Expected no rule metrics after reset")
	}
	global = metrics.GetGlobalMetrics()
	if global.TotalActivations != 0 {
		t.Error("Expected zero activations after reset")
	}
	if global.TotalEvaluations != 0 {
		t.Error("Expected zero evaluations after reset")
	}
	if global.TotalCacheHits != 0 {
		t.Error("Expected zero cache hits after reset")
	}
}
// TestMaxRulesToTrack teste la limite du nombre de règles suivies
func TestMaxRulesToTrack(t *testing.T) {
	config := DefaultMetricsConfig()
	config.MaxRulesToTrack = 3
	metrics := NewArithmeticDecompositionMetrics(config)
	// Ajouter 5 règles (dépasse la limite de 3)
	for i := 1; i <= 5; i++ {
		ruleID := fmt.Sprintf("rule%d", i)
		metrics.RecordActivation(ruleID, true, 100*time.Microsecond)
		time.Sleep(1 * time.Millisecond) // Pour assurer des timestamps différents
	}
	allMetrics := metrics.GetAllRuleMetrics()
	// Devrait avoir seulement 3 règles (les plus récentes)
	if len(allMetrics) > config.MaxRulesToTrack {
		t.Errorf("Expected max %d rules, got %d", config.MaxRulesToTrack, len(allMetrics))
	}
}
// TestMetricsDisabled teste que les métriques peuvent être désactivées
func TestMetricsDisabled(t *testing.T) {
	config := DefaultMetricsConfig()
	config.Enabled = false
	metrics := NewArithmeticDecompositionMetrics(config)
	// Tenter d'enregistrer des métriques
	metrics.RecordActivation("rule1", true, 100*time.Microsecond)
	metrics.RecordEvaluation("rule1", true, 200*time.Microsecond)
	metrics.RecordCacheHit("rule1")
	// Aucune métrique ne devrait être enregistrée
	if metrics.GetRuleMetrics("rule1") != nil {
		t.Error("Expected no metrics when disabled")
	}
	global := metrics.GetGlobalMetrics()
	if global.TotalActivations != 0 {
		t.Error("Expected zero activations when disabled")
	}
}
// TestConcurrentMetrics teste la sécurité concurrentielle
func TestConcurrentMetrics(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100
	// Lancer plusieurs goroutines qui enregistrent des métriques
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ruleID := fmt.Sprintf("rule%d", id)
			for j := 0; j < numOperations; j++ {
				metrics.RecordActivation(ruleID, true, 100*time.Microsecond)
				metrics.RecordEvaluation(ruleID, true, 200*time.Microsecond)
				metrics.RecordCacheHit(ruleID)
				metrics.RecordCacheMiss(ruleID)
			}
		}(i)
	}
	// Attendre que toutes les goroutines se terminent
	wg.Wait()
	// Vérifier que les métriques sont cohérentes
	global := metrics.GetGlobalMetrics()
	expectedActivations := int64(numGoroutines * numOperations)
	if global.TotalActivations != expectedActivations {
		t.Errorf("Expected %d activations, got %d", expectedActivations, global.TotalActivations)
	}
	if global.TotalEvaluations != expectedActivations {
		t.Errorf("Expected %d evaluations, got %d", expectedActivations, global.TotalEvaluations)
	}
	expectedCacheOps := int64(numGoroutines * numOperations)
	if global.TotalCacheHits != expectedCacheOps {
		t.Errorf("Expected %d cache hits, got %d", expectedCacheOps, global.TotalCacheHits)
	}
	if global.TotalCacheMisses != expectedCacheOps {
		t.Errorf("Expected %d cache misses, got %d", expectedCacheOps, global.TotalCacheMisses)
	}
}
// TestCalculateMaxDepth teste le calcul de profondeur maximale
func TestCalculateMaxDepth(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	testCases := []struct {
		name         string
		dependencies map[string][]string
		expectedMax  int
	}{
		{
			name:         "empty dependencies",
			dependencies: map[string][]string{},
			expectedMax:  0,
		},
		{
			name: "single level",
			dependencies: map[string][]string{
				"A": {},
				"B": {},
			},
			expectedMax: 0,
		},
		{
			name: "two levels",
			dependencies: map[string][]string{
				"A": {},
				"B": {"A"},
			},
			expectedMax: 1,
		},
		{
			name: "three levels linear",
			dependencies: map[string][]string{
				"A": {},
				"B": {"A"},
				"C": {"B"},
			},
			expectedMax: 2,
		},
		{
			name: "branching dependencies",
			dependencies: map[string][]string{
				"A": {},
				"B": {},
				"C": {"A", "B"},
				"D": {"C"},
			},
			expectedMax: 2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ruleID := "test_rule"
			metrics.RecordChainStructure(ruleID, 1, 1, 0, []string{}, tc.dependencies)
			rule := metrics.GetRuleMetrics(ruleID)
			if rule.MaxDependencyDepth != tc.expectedMax {
				t.Errorf("Expected max depth %d, got %d", tc.expectedMax, rule.MaxDependencyDepth)
			}
			// Réinitialiser pour le prochain test
			metrics.Reset()
		})
	}
}
// TestGlobalAverages teste le calcul des moyennes globales
func TestGlobalAverages(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	// Ajouter plusieurs règles avec différentes structures
	metrics.RecordChainStructure("rule1", 5, 3, 2, []string{}, map[string][]string{
		"A": {},
		"B": {"A"},
	})
	metrics.RecordChainStructure("rule2", 7, 4, 3, []string{}, map[string][]string{
		"A": {},
		"B": {"A"},
		"C": {"B"},
	})
	metrics.RecordChainStructure("rule3", 3, 2, 1, []string{}, map[string][]string{})
	global := metrics.GetGlobalMetrics()
	expectedAvgChainLength := (5.0 + 7.0 + 3.0) / 3.0
	if global.AverageChainLength != expectedAvgChainLength {
		t.Errorf("Expected avg chain length %.2f, got %.2f", expectedAvgChainLength, global.AverageChainLength)
	}
	expectedAvgAtomicSteps := (3.0 + 4.0 + 2.0) / 3.0
	if global.AverageAtomicStepsPerChain != expectedAvgAtomicSteps {
		t.Errorf("Expected avg atomic steps %.2f, got %.2f", expectedAvgAtomicSteps, global.AverageAtomicStepsPerChain)
	}
	expectedAvgDepth := (1.0 + 2.0 + 0.0) / 3.0
	if global.AverageDependencyDepth != expectedAvgDepth {
		t.Errorf("Expected avg dependency depth %.2f, got %.2f", expectedAvgDepth, global.AverageDependencyDepth)
	}
}
// TestCopyRuleMetrics teste que les copies sont indépendantes
func TestCopyRuleMetrics(t *testing.T) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	ruleID := "rule1"
	metrics.RecordActivation(ruleID, true, 100*time.Microsecond)
	metrics.RecordChainStructure(ruleID, 5, 3, 2, []string{"temp1"}, map[string][]string{
		"temp1": {},
	})
	// Obtenir une copie
	copy1 := metrics.GetRuleMetrics(ruleID)
	if copy1 == nil {
		t.Fatal("Expected non-nil copy")
	}
	// Modifier la copie
	copy1.TotalActivations = 999
	copy1.IntermediateResults = append(copy1.IntermediateResults, "temp2")
	copy1.Dependencies["temp2"] = []string{"temp1"}
	// Obtenir une nouvelle copie et vérifier qu'elle n'est pas affectée
	copy2 := metrics.GetRuleMetrics(ruleID)
	if copy2 == nil {
		t.Fatal("Expected non-nil second copy")
	}
	if copy2.TotalActivations == 999 {
		t.Error("Copy should be independent - activations should not be 999")
	}
	if len(copy2.IntermediateResults) > 1 {
		t.Error("Copy should be independent - intermediate results should not be modified")
	}
	if _, exists := copy2.Dependencies["temp2"]; exists {
		t.Error("Copy should be independent - dependencies should not be modified")
	}
}
// BenchmarkRecordEvaluation benchmark l'enregistrement d'évaluations
func BenchmarkRecordEvaluation(b *testing.B) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.RecordEvaluation("rule1", true, 100*time.Microsecond)
	}
}
// BenchmarkGetRuleMetrics benchmark la récupération de métriques
func BenchmarkGetRuleMetrics(b *testing.B) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	// Pré-remplir
	for i := 0; i < 100; i++ {
		ruleID := "rule" + string(rune('0'+i%10))
		metrics.RecordEvaluation(ruleID, true, 100*time.Microsecond)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.GetRuleMetrics("rule1")
	}
}
// BenchmarkConcurrentRecording benchmark l'enregistrement concurrent
func BenchmarkConcurrentRecording(b *testing.B) {
	config := DefaultMetricsConfig()
	metrics := NewArithmeticDecompositionMetrics(config)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ruleID := "rule" + string(rune('0'+i%10))
			metrics.RecordEvaluation(ruleID, true, 100*time.Microsecond)
			i++
		}
	})
}