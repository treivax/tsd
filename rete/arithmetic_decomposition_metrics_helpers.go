// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"sort"
	"time"
)

// ============================================================================
// Private Helper Functions
// ============================================================================

// getOrCreateRuleMetrics obtient ou crée les métriques pour une règle
func (adm *ArithmeticDecompositionMetrics) getOrCreateRuleMetrics(ruleID string) *RuleArithmeticMetrics {
	if rule, exists := adm.ruleMetrics[ruleID]; exists {
		return rule
	}

	// Vérifier la limite
	if len(adm.ruleMetrics) >= adm.config.MaxRulesToTrack {
		// Supprimer la règle la plus ancienne
		adm.evictOldestRule()
	}

	rule := &RuleArithmeticMetrics{
		RuleID:                  ruleID,
		EvaluationTimeHistogram: make(map[int]int64),
		Dependencies:            make(map[string][]string),
		Metadata:                make(map[string]interface{}),
		FirstSeen:               time.Now(),
		LastSeen:                time.Now(),
		MinEvaluationTime:       time.Duration(1<<63 - 1),
	}

	adm.ruleMetrics[ruleID] = rule
	adm.global.TotalRulesWithArithmetic++

	return rule
}

// getHistogramBucket trouve le bucket approprié pour une durée
func (adm *ArithmeticDecompositionMetrics) getHistogramBucket(duration time.Duration) int {
	micros := int(duration.Microseconds())

	// Trouver le bucket le plus proche
	for _, bucket := range adm.config.HistogramBuckets {
		if micros <= bucket {
			return bucket
		}
	}

	// Si plus grand que tous les buckets, utiliser le dernier
	if len(adm.config.HistogramBuckets) > 0 {
		return adm.config.HistogramBuckets[len(adm.config.HistogramBuckets)-1]
	}

	return micros
}

// updateCacheHitRate met à jour le taux de cache hit pour une règle
func (adm *ArithmeticDecompositionMetrics) updateCacheHitRate(rule *RuleArithmeticMetrics) {
	total := rule.CacheHits + rule.CacheMisses
	if total > 0 {
		rule.CacheHitRate = float64(rule.CacheHits) / float64(total)
	}
}

// updateGlobalCacheHitRate met à jour le taux global de cache hit
func (adm *ArithmeticDecompositionMetrics) updateGlobalCacheHitRate() {
	total := adm.global.TotalCacheHits + adm.global.TotalCacheMisses
	if total > 0 {
		adm.global.CacheGlobalHitRate = float64(adm.global.TotalCacheHits) / float64(total)
	}
}

// calculateMaxDepth calcule la profondeur maximale d'un graphe de dépendances
func (adm *ArithmeticDecompositionMetrics) calculateMaxDepth(dependencies map[string][]string) int {
	if len(dependencies) == 0 {
		return 0
	}

	depths := make(map[string]int)
	var calculateDepth func(string) int

	calculateDepth = func(node string) int {
		if depth, exists := depths[node]; exists {
			return depth
		}

		maxChildDepth := 0
		for _, dep := range dependencies[node] {
			childDepth := calculateDepth(dep)
			if childDepth+1 > maxChildDepth {
				maxChildDepth = childDepth + 1
			}
		}

		depths[node] = maxChildDepth
		return maxChildDepth
	}

	maxDepth := 0
	for node := range dependencies {
		depth := calculateDepth(node)
		if depth > maxDepth {
			maxDepth = depth
		}
	}

	return maxDepth
}

// recalculateGlobalAverages recalcule les moyennes globales
func (adm *ArithmeticDecompositionMetrics) recalculateGlobalAverages() {
	if adm.global.TotalDecomposedChains > 0 {
		totalChainLength := 0
		totalAtomicSteps := 0
		totalDepth := 0

		for _, rule := range adm.ruleMetrics {
			totalChainLength += rule.ChainLength
			totalAtomicSteps += rule.AtomicStepsCount
			totalDepth += rule.MaxDependencyDepth
		}

		count := float64(len(adm.ruleMetrics))
		if count > 0 {
			adm.global.AverageChainLength = float64(totalChainLength) / count
			adm.global.AverageAtomicStepsPerChain = float64(totalAtomicSteps) / count
			adm.global.AverageDependencyDepth = float64(totalDepth) / count
		}
	}
}

// calculatePercentiles calcule les percentiles des temps d'évaluation
func (adm *ArithmeticDecompositionMetrics) calculatePercentiles() {
	// Collecter tous les temps d'évaluation
	var times []time.Duration
	for _, rule := range adm.ruleMetrics {
		for bucket, count := range rule.EvaluationTimeHistogram {
			duration := time.Duration(bucket) * time.Microsecond
			for i := int64(0); i < count; i++ {
				times = append(times, duration)
			}
		}
	}

	if len(times) == 0 {
		return
	}

	sort.Slice(times, func(i, j int) bool {
		return times[i] < times[j]
	})

	// Calculer P50, P95, P99
	adm.global.EvaluationTimeP50 = times[len(times)*50/100]
	adm.global.EvaluationTimeP95 = times[len(times)*95/100]
	adm.global.EvaluationTimeP99 = times[len(times)*99/100]
}

// evictOldestRule supprime la règle la plus ancienne
func (adm *ArithmeticDecompositionMetrics) evictOldestRule() {
	var oldestRuleID string
	var oldestTime time.Time

	for ruleID, rule := range adm.ruleMetrics {
		if oldestTime.IsZero() || rule.LastSeen.Before(oldestTime) {
			oldestTime = rule.LastSeen
			oldestRuleID = ruleID
		}
	}

	if oldestRuleID != "" {
		delete(adm.ruleMetrics, oldestRuleID)
	}
}

// copyRuleMetrics crée une copie d'une métrique de règle
func (adm *ArithmeticDecompositionMetrics) copyRuleMetrics(rule *RuleArithmeticMetrics) *RuleArithmeticMetrics {
	copy := &RuleArithmeticMetrics{
		RuleID:                  rule.RuleID,
		RuleName:                rule.RuleName,
		TotalActivations:        rule.TotalActivations,
		SuccessfulActivations:   rule.SuccessfulActivations,
		FailedActivations:       rule.FailedActivations,
		TotalEvaluations:        rule.TotalEvaluations,
		SuccessfulEvaluations:   rule.SuccessfulEvaluations,
		FailedEvaluations:       rule.FailedEvaluations,
		ChainLength:             rule.ChainLength,
		AtomicStepsCount:        rule.AtomicStepsCount,
		ComparisonStepsCount:    rule.ComparisonStepsCount,
		TotalEvaluationTime:     rule.TotalEvaluationTime,
		MinEvaluationTime:       rule.MinEvaluationTime,
		MaxEvaluationTime:       rule.MaxEvaluationTime,
		AvgEvaluationTime:       rule.AvgEvaluationTime,
		EvaluationTimeHistogram: make(map[int]int64),
		CacheHits:               rule.CacheHits,
		CacheMisses:             rule.CacheMisses,
		CacheHitRate:            rule.CacheHitRate,
		CacheEnabled:            rule.CacheEnabled,
		MaxDependencyDepth:      rule.MaxDependencyDepth,
		HasCircularDeps:         rule.HasCircularDeps,
		FirstSeen:               rule.FirstSeen,
		LastSeen:                rule.LastSeen,
		IntermediateResults:     make([]string, len(rule.IntermediateResults)),
		IsolatedNodes:           make([]string, len(rule.IsolatedNodes)),
		Dependencies:            make(map[string][]string),
		Metadata:                make(map[string]interface{}),
	}

	// Copier les slices et maps
	copy.IntermediateResults = append([]string{}, rule.IntermediateResults...)
	copy.IsolatedNodes = append([]string{}, rule.IsolatedNodes...)

	for k, v := range rule.EvaluationTimeHistogram {
		copy.EvaluationTimeHistogram[k] = v
	}

	for k, v := range rule.Dependencies {
		copy.Dependencies[k] = append([]string{}, v...)
	}

	for k, v := range rule.Metadata {
		copy.Metadata[k] = v
	}

	return copy
}
