// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"sort"
)

// ============================================================================
// Query and Retrieval Functions
// ============================================================================

// GetRuleMetrics retourne les métriques pour une règle spécifique
func (adm *ArithmeticDecompositionMetrics) GetRuleMetrics(ruleID string) *RuleArithmeticMetrics {
	adm.mutex.RLock()
	defer adm.mutex.RUnlock()

	if rule, exists := adm.ruleMetrics[ruleID]; exists {
		// Retourner une copie pour éviter les modifications externes
		return adm.copyRuleMetrics(rule)
	}
	return nil
}

// GetGlobalMetrics retourne les métriques globales
func (adm *ArithmeticDecompositionMetrics) GetGlobalMetrics() GlobalArithmeticMetrics {
	adm.mutex.RLock()
	defer adm.mutex.RUnlock()

	// Calculer les percentiles si activés
	if adm.config.CollectPercentiles {
		adm.calculatePercentiles()
	}

	return adm.global
}

// GetAllRuleMetrics retourne les métriques de toutes les règles
func (adm *ArithmeticDecompositionMetrics) GetAllRuleMetrics() map[string]*RuleArithmeticMetrics {
	adm.mutex.RLock()
	defer adm.mutex.RUnlock()

	result := make(map[string]*RuleArithmeticMetrics)
	for ruleID, metrics := range adm.ruleMetrics {
		result[ruleID] = adm.copyRuleMetrics(metrics)
	}
	return result
}

// GetTopRulesByEvaluations retourne les N règles avec le plus d'évaluations
func (adm *ArithmeticDecompositionMetrics) GetTopRulesByEvaluations(n int) []*RuleArithmeticMetrics {
	adm.mutex.RLock()
	defer adm.mutex.RUnlock()

	rules := make([]*RuleArithmeticMetrics, 0, len(adm.ruleMetrics))
	for _, rule := range adm.ruleMetrics {
		rules = append(rules, adm.copyRuleMetrics(rule))
	}

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].TotalEvaluations > rules[j].TotalEvaluations
	})

	if n > len(rules) {
		n = len(rules)
	}
	return rules[:n]
}

// GetTopRulesByDuration retourne les N règles avec le temps total le plus élevé
func (adm *ArithmeticDecompositionMetrics) GetTopRulesByDuration(n int) []*RuleArithmeticMetrics {
	adm.mutex.RLock()
	defer adm.mutex.RUnlock()

	rules := make([]*RuleArithmeticMetrics, 0, len(adm.ruleMetrics))
	for _, rule := range adm.ruleMetrics {
		rules = append(rules, adm.copyRuleMetrics(rule))
	}

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].TotalEvaluationTime > rules[j].TotalEvaluationTime
	})

	if n > len(rules) {
		n = len(rules)
	}
	return rules[:n]
}

// GetSlowestRules retourne les N règles avec le temps moyen le plus élevé
func (adm *ArithmeticDecompositionMetrics) GetSlowestRules(n int) []*RuleArithmeticMetrics {
	adm.mutex.RLock()
	defer adm.mutex.RUnlock()

	rules := make([]*RuleArithmeticMetrics, 0, len(adm.ruleMetrics))
	for _, rule := range adm.ruleMetrics {
		if rule.TotalEvaluations > 0 {
			rules = append(rules, adm.copyRuleMetrics(rule))
		}
	}

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].AvgEvaluationTime > rules[j].AvgEvaluationTime
	})

	if n > len(rules) {
		n = len(rules)
	}
	return rules[:n]
}

// GetSummary retourne un résumé formaté des métriques
func (adm *ArithmeticDecompositionMetrics) GetSummary() map[string]interface{} {
	adm.mutex.RLock()
	defer adm.mutex.RUnlock()

	return map[string]interface{}{
		"rules": map[string]interface{}{
			"total_with_arithmetic":   adm.global.TotalRulesWithArithmetic,
			"total_decomposed_chains": adm.global.TotalDecomposedChains,
			"tracked_rules":           len(adm.ruleMetrics),
		},
		"nodes": map[string]interface{}{
			"total_atomic":         adm.global.TotalAtomicNodes,
			"total_comparison":     adm.global.TotalComparisonNodes,
			"average_chain_length": adm.global.AverageChainLength,
		},
		"evaluations": map[string]interface{}{
			"total":        adm.global.TotalEvaluations,
			"total_time":   adm.global.TotalEvaluationTime.String(),
			"average_time": adm.global.AverageEvaluationTime.String(),
			"min_time":     adm.global.MinEvaluationTime.String(),
			"max_time":     adm.global.MaxEvaluationTime.String(),
		},
		"cache": map[string]interface{}{
			"hits":               adm.global.TotalCacheHits,
			"misses":             adm.global.TotalCacheMisses,
			"hit_rate":           adm.global.CacheGlobalHitRate,
			"size":               adm.global.CacheSize,
			"evictions":          adm.global.CacheEvictions,
			"memory_usage_bytes": adm.global.CacheMemoryUsage,
		},
		"validation": map[string]interface{}{
			"total_validations": adm.global.GraphValidations,
			"cycles_detected":   adm.global.CyclesDetected,
			"max_graph_depth":   adm.global.MaxGraphDepth,
		},
	}
}
