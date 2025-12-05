// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
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

// Note: Query and retrieval functions have been extracted to arithmetic_decomposition_metrics_query.go

// Reset réinitialise toutes les métriques
func (adm *ArithmeticDecompositionMetrics) Reset() {
	adm.mutex.Lock()
	defer adm.mutex.Unlock()

	adm.ruleMetrics = make(map[string]*RuleArithmeticMetrics)
	adm.global = GlobalArithmeticMetrics{
		MinEvaluationTime: time.Duration(1<<63 - 1),
	}
}

// Note: Private helper functions have been extracted to arithmetic_decomposition_metrics_helpers.go
