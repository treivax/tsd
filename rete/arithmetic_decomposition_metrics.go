// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"sort"
	"sync"
	"time"
)

// ArithmeticDecompositionMetrics collecte des métriques détaillées sur la décomposition
// arithmétique et l'évaluation des expressions. Thread-safe pour usage concurrent.
type ArithmeticDecompositionMetrics struct {
	// Métriques par règle
	ruleMetrics map[string]*RuleArithmeticMetrics

	// Métriques globales
	global GlobalArithmeticMetrics

	// Configuration
	config MetricsConfig

	// Protection concurrence
	mutex sync.RWMutex
}

// RuleArithmeticMetrics contient les métriques pour une règle spécifique
type RuleArithmeticMetrics struct {
	RuleID   string
	RuleName string

	// Compteurs d'activations
	TotalActivations      int64
	SuccessfulActivations int64
	FailedActivations     int64

	// Compteurs d'évaluations
	TotalEvaluations      int64
	SuccessfulEvaluations int64
	FailedEvaluations     int64

	// Structure de la chaîne
	ChainLength          int
	AtomicStepsCount     int
	ComparisonStepsCount int
	IntermediateResults  []string
	Dependencies         map[string][]string

	// Temps d'évaluation
	TotalEvaluationTime time.Duration
	MinEvaluationTime   time.Duration
	MaxEvaluationTime   time.Duration
	AvgEvaluationTime   time.Duration

	// Histogramme des temps (buckets en microsecondes)
	EvaluationTimeHistogram map[int]int64 // bucket -> count

	// Statistiques de cache
	CacheHits    int64
	CacheMisses  int64
	CacheHitRate float64
	CacheEnabled bool

	// Statistiques de dépendances
	MaxDependencyDepth int
	HasCircularDeps    bool
	IsolatedNodes      []string

	// Timestamps
	FirstSeen time.Time
	LastSeen  time.Time

	// Métadonnées
	Metadata map[string]interface{}
}

// GlobalArithmeticMetrics contient les métriques agrégées globales
type GlobalArithmeticMetrics struct {
	// Statistiques globales
	TotalRulesWithArithmetic int
	TotalDecomposedChains    int
	TotalAtomicNodes         int
	TotalComparisonNodes     int

	// Moyennes
	AverageChainLength         float64
	AverageAtomicStepsPerChain float64
	AverageDependencyDepth     float64

	// Ratios
	SharedNodesRatio   float64
	CacheGlobalHitRate float64

	// Compteurs globaux
	TotalActivations          int64
	TotalEvaluations          int64
	TotalCacheHits            int64
	TotalCacheMisses          int64
	TotalCircularDepsDetected int64

	// Temps globaux
	TotalEvaluationTime   time.Duration
	AverageEvaluationTime time.Duration
	MinEvaluationTime     time.Duration
	MaxEvaluationTime     time.Duration

	// Distribution des temps (percentiles)
	EvaluationTimeP50 time.Duration
	EvaluationTimeP95 time.Duration
	EvaluationTimeP99 time.Duration

	// Statistiques de cache global
	CacheSize        int
	CacheEvictions   int64
	CacheMemoryUsage int64

	// Statistiques de détection de cycles
	GraphValidations int64
	CyclesDetected   int64
	MaxGraphDepth    int
}

// MetricsConfig configure le comportement de la collecte de métriques
type MetricsConfig struct {
	Enabled             bool
	CollectHistograms   bool
	CollectPercentiles  bool
	HistogramBuckets    []int // En microsecondes
	MaxRulesToTrack     int
	RetentionDuration   time.Duration
	AggregationInterval time.Duration
}

// DefaultMetricsConfig retourne une configuration par défaut
func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		Enabled:            true,
		CollectHistograms:  true,
		CollectPercentiles: true,
		HistogramBuckets: []int{
			1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000,
		},
		MaxRulesToTrack:     1000,
		RetentionDuration:   24 * time.Hour,
		AggregationInterval: 1 * time.Minute,
	}
}

// NewArithmeticDecompositionMetrics crée une nouvelle instance de métriques
func NewArithmeticDecompositionMetrics(config MetricsConfig) *ArithmeticDecompositionMetrics {
	return &ArithmeticDecompositionMetrics{
		ruleMetrics: make(map[string]*RuleArithmeticMetrics),
		config:      config,
		global: GlobalArithmeticMetrics{
			MinEvaluationTime: time.Duration(1<<63 - 1), // Max duration
		},
	}
}

// RecordActivation enregistre une activation de règle
func (adm *ArithmeticDecompositionMetrics) RecordActivation(ruleID string, success bool, duration time.Duration) {
	if !adm.config.Enabled {
		return
	}

	adm.mutex.Lock()
	defer adm.mutex.Unlock()

	rule := adm.getOrCreateRuleMetrics(ruleID)
	rule.TotalActivations++
	if success {
		rule.SuccessfulActivations++
	} else {
		rule.FailedActivations++
	}
	rule.LastSeen = time.Now()

	adm.global.TotalActivations++
}

// RecordEvaluation enregistre une évaluation d'expression arithmétique
func (adm *ArithmeticDecompositionMetrics) RecordEvaluation(ruleID string, success bool, duration time.Duration) {
	if !adm.config.Enabled {
		return
	}

	adm.mutex.Lock()
	defer adm.mutex.Unlock()

	rule := adm.getOrCreateRuleMetrics(ruleID)
	rule.TotalEvaluations++
	if success {
		rule.SuccessfulEvaluations++
	} else {
		rule.FailedEvaluations++
	}

	// Mettre à jour les temps
	rule.TotalEvaluationTime += duration
	if rule.MinEvaluationTime == 0 || duration < rule.MinEvaluationTime {
		rule.MinEvaluationTime = duration
	}
	if duration > rule.MaxEvaluationTime {
		rule.MaxEvaluationTime = duration
	}
	if rule.TotalEvaluations > 0 {
		rule.AvgEvaluationTime = rule.TotalEvaluationTime / time.Duration(rule.TotalEvaluations)
	}

	// Histogramme
	if adm.config.CollectHistograms {
		bucket := adm.getHistogramBucket(duration)
		rule.EvaluationTimeHistogram[bucket]++
	}

	// Mettre à jour les globales
	adm.global.TotalEvaluations++
	adm.global.TotalEvaluationTime += duration
	if duration < adm.global.MinEvaluationTime {
		adm.global.MinEvaluationTime = duration
	}
	if duration > adm.global.MaxEvaluationTime {
		adm.global.MaxEvaluationTime = duration
	}
	if adm.global.TotalEvaluations > 0 {
		adm.global.AverageEvaluationTime = adm.global.TotalEvaluationTime / time.Duration(adm.global.TotalEvaluations)
	}
}

// RecordCacheHit enregistre un cache hit
func (adm *ArithmeticDecompositionMetrics) RecordCacheHit(ruleID string) {
	if !adm.config.Enabled {
		return
	}

	adm.mutex.Lock()
	defer adm.mutex.Unlock()

	rule := adm.getOrCreateRuleMetrics(ruleID)
	rule.CacheHits++
	rule.CacheEnabled = true
	adm.updateCacheHitRate(rule)

	adm.global.TotalCacheHits++
	adm.updateGlobalCacheHitRate()
}

// RecordCacheMiss enregistre un cache miss
func (adm *ArithmeticDecompositionMetrics) RecordCacheMiss(ruleID string) {
	if !adm.config.Enabled {
		return
	}

	adm.mutex.Lock()
	defer adm.mutex.Unlock()

	rule := adm.getOrCreateRuleMetrics(ruleID)
	rule.CacheMisses++
	rule.CacheEnabled = true
	adm.updateCacheHitRate(rule)

	adm.global.TotalCacheMisses++
	adm.updateGlobalCacheHitRate()
}

// RecordChainStructure enregistre la structure d'une chaîne décomposée
func (adm *ArithmeticDecompositionMetrics) RecordChainStructure(ruleID string, chainLength, atomicSteps, comparisonSteps int, intermediateResults []string, dependencies map[string][]string) {
	if !adm.config.Enabled {
		return
	}

	adm.mutex.Lock()
	defer adm.mutex.Unlock()

	rule := adm.getOrCreateRuleMetrics(ruleID)
	rule.ChainLength = chainLength
	rule.AtomicStepsCount = atomicSteps
	rule.ComparisonStepsCount = comparisonSteps
	rule.IntermediateResults = intermediateResults
	rule.Dependencies = dependencies

	// Calculer la profondeur max de dépendances
	rule.MaxDependencyDepth = adm.calculateMaxDepth(dependencies)

	// Mettre à jour les globales
	adm.global.TotalDecomposedChains++
	adm.global.TotalAtomicNodes += atomicSteps
	adm.global.TotalComparisonNodes += comparisonSteps
	adm.recalculateGlobalAverages()
}

// RecordCircularDependency enregistre la détection d'une dépendance circulaire
func (adm *ArithmeticDecompositionMetrics) RecordCircularDependency(ruleID string, cyclePath []string) {
	if !adm.config.Enabled {
		return
	}

	adm.mutex.Lock()
	defer adm.mutex.Unlock()

	rule := adm.getOrCreateRuleMetrics(ruleID)
	rule.HasCircularDeps = true
	if rule.Metadata == nil {
		rule.Metadata = make(map[string]interface{})
	}
	rule.Metadata["cycle_path"] = cyclePath

	adm.global.TotalCircularDepsDetected++
	adm.global.CyclesDetected++
}

// RecordGraphValidation enregistre une validation de graphe
func (adm *ArithmeticDecompositionMetrics) RecordGraphValidation(maxDepth int, hasCycles bool) {
	if !adm.config.Enabled {
		return
	}

	adm.mutex.Lock()
	defer adm.mutex.Unlock()

	adm.global.GraphValidations++
	if hasCycles {
		adm.global.CyclesDetected++
	}
	if maxDepth > adm.global.MaxGraphDepth {
		adm.global.MaxGraphDepth = maxDepth
	}
}

// UpdateCacheStatistics met à jour les statistiques du cache global
func (adm *ArithmeticDecompositionMetrics) UpdateCacheStatistics(size int, evictions int64, memoryUsage int64) {
	if !adm.config.Enabled {
		return
	}

	adm.mutex.Lock()
	defer adm.mutex.Unlock()

	adm.global.CacheSize = size
	adm.global.CacheEvictions = evictions
	adm.global.CacheMemoryUsage = memoryUsage
}

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

// Reset réinitialise toutes les métriques
func (adm *ArithmeticDecompositionMetrics) Reset() {
	adm.mutex.Lock()
	defer adm.mutex.Unlock()

	adm.ruleMetrics = make(map[string]*RuleArithmeticMetrics)
	adm.global = GlobalArithmeticMetrics{
		MinEvaluationTime: time.Duration(1<<63 - 1),
	}
}

// --- Méthodes privées ---

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
