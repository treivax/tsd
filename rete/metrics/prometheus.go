// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusExporter expose les métriques RETE au format Prometheus
type PrometheusExporter struct {
	// Ingestion metrics
	ingestionDuration  *prometheus.HistogramVec
	ingestionTotal     *prometheus.CounterVec
	typesAdded         prometheus.Counter
	rulesAdded         prometheus.Counter
	factsSubmitted     prometheus.Counter
	factsCollected     prometheus.Counter
	factsPropagated    prometheus.Counter
	terminalsAdded     prometheus.Counter
	propagationTargets prometheus.Gauge

	// Network state metrics
	typeNodes     prometheus.Gauge
	terminalNodes prometheus.Gauge
	alphaNodes    prometheus.Gauge
	betaNodes     prometheus.Gauge

	// Phase duration metrics
	parsingDuration        prometheus.Histogram
	validationDuration     prometheus.Histogram
	typeCreationDuration   prometheus.Histogram
	ruleCreationDuration   prometheus.Histogram
	factCollectionDuration prometheus.Histogram
	propagationDuration    prometheus.Histogram
	factSubmissionDuration prometheus.Histogram

	// Efficiency metrics
	efficiencyScore prometheus.Gauge
	bottleneckPhase *prometheus.GaugeVec

	// Transaction metrics
	transactionTotal    *prometheus.CounterVec
	transactionDuration *prometheus.HistogramVec
	commandsExecuted    prometheus.Counter
	rollbacksTotal      prometheus.Counter

	// Storage metrics
	storageFacts      prometheus.Gauge
	storageOperations *prometheus.CounterVec
	storageErrors     *prometheus.CounterVec

	// Coherence metrics
	coherenceViolations *prometheus.CounterVec
	coherenceChecks     prometheus.Counter
	coherenceMode       *prometheus.GaugeVec

	// Performance metrics
	ruleEvaluations    prometheus.Counter
	ruleMatches        prometheus.Counter
	ruleEvaluationTime prometheus.Histogram

	// Arithmetic decomposition metrics (if enabled)
	arithmeticDecompositions prometheus.Counter
	arithmeticOptimizations  prometheus.Counter

	registry *prometheus.Registry
	mutex    sync.RWMutex
}

// NewPrometheusExporter crée un nouveau exporteur Prometheus
func NewPrometheusExporter(namespace string) *PrometheusExporter {
	if namespace == "" {
		namespace = "tsd"
	}

	registry := prometheus.NewRegistry()
	factory := promauto.With(registry)

	exporter := &PrometheusExporter{
		registry: registry,

		// Ingestion metrics
		ingestionDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "ingestion_duration_seconds",
				Help:      "Duration of RETE ingestion operations by phase",
				Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"phase", "mode"},
		),

		ingestionTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "ingestion_total",
				Help:      "Total number of RETE ingestion operations",
			},
			[]string{"mode", "status"},
		),

		typesAdded: factory.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "types_added_total",
				Help:      "Total number of types added to RETE network",
			},
		),

		rulesAdded: factory.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "rules_added_total",
				Help:      "Total number of rules added to RETE network",
			},
		),

		factsSubmitted: factory.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "facts_submitted_total",
				Help:      "Total number of facts submitted to RETE network",
			},
		),

		factsCollected: factory.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "facts_collected_total",
				Help:      "Total number of existing facts collected",
			},
		),

		factsPropagated: factory.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "facts_propagated_total",
				Help:      "Total number of facts propagated through network",
			},
		),

		terminalsAdded: factory.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "terminals_added_total",
				Help:      "Total number of terminal nodes added",
			},
		),

		propagationTargets: factory.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "propagation_targets",
				Help:      "Current number of propagation targets",
			},
		),

		// Network state metrics
		typeNodes: factory.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "type_nodes",
				Help:      "Current number of type nodes in RETE network",
			},
		),

		terminalNodes: factory.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "terminal_nodes",
				Help:      "Current number of terminal nodes in RETE network",
			},
		),

		alphaNodes: factory.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "alpha_nodes",
				Help:      "Current number of alpha nodes in RETE network",
			},
		),

		betaNodes: factory.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "beta_nodes",
				Help:      "Current number of beta nodes in RETE network",
			},
		),

		// Phase duration metrics
		parsingDuration: factory.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "parsing_duration_seconds",
				Help:      "Duration of parsing phase",
				Buckets:   prometheus.DefBuckets,
			},
		),

		validationDuration: factory.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "validation_duration_seconds",
				Help:      "Duration of validation phase",
				Buckets:   prometheus.DefBuckets,
			},
		),

		typeCreationDuration: factory.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "type_creation_duration_seconds",
				Help:      "Duration of type creation phase",
				Buckets:   prometheus.DefBuckets,
			},
		),

		ruleCreationDuration: factory.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "rule_creation_duration_seconds",
				Help:      "Duration of rule creation phase",
				Buckets:   prometheus.DefBuckets,
			},
		),

		factCollectionDuration: factory.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "fact_collection_duration_seconds",
				Help:      "Duration of fact collection phase",
				Buckets:   prometheus.DefBuckets,
			},
		),

		propagationDuration: factory.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "propagation_duration_seconds",
				Help:      "Duration of propagation phase",
				Buckets:   prometheus.DefBuckets,
			},
		),

		factSubmissionDuration: factory.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "fact_submission_duration_seconds",
				Help:      "Duration of fact submission phase",
				Buckets:   prometheus.DefBuckets,
			},
		),

		// Efficiency metrics
		efficiencyScore: factory.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "efficiency_score",
				Help:      "RETE network efficiency score (0-1)",
			},
		),

		bottleneckPhase: factory.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "bottleneck_phase_percentage",
				Help:      "Percentage of time spent in each phase",
			},
			[]string{"phase"},
		),

		// Transaction metrics
		transactionTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "transaction_total",
				Help:      "Total number of transactions",
			},
			[]string{"operation", "status"},
		),

		transactionDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "transaction_duration_seconds",
				Help:      "Duration of transaction operations",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"operation"},
		),

		commandsExecuted: factory.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "commands_executed_total",
				Help:      "Total number of commands executed in transactions",
			},
		),

		rollbacksTotal: factory.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "rete",
				Name:      "rollbacks_total",
				Help:      "Total number of transaction rollbacks",
			},
		),

		// Storage metrics
		storageFacts: factory.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "storage",
				Name:      "facts_total",
				Help:      "Current number of facts in storage",
			},
		),

		storageOperations: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "storage",
				Name:      "operations_total",
				Help:      "Total number of storage operations",
			},
			[]string{"operation", "type"},
		),

		storageErrors: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "storage",
				Name:      "errors_total",
				Help:      "Total number of storage errors",
			},
			[]string{"operation", "error_type"},
		),

		// Coherence metrics
		coherenceViolations: factory.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "coherence",
				Name:      "violations_total",
				Help:      "Total number of coherence violations detected",
			},
			[]string{"type"},
		),

		coherenceChecks: factory.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "coherence",
				Name:      "checks_total",
				Help:      "Total number of coherence checks performed",
			},
		),

		coherenceMode: factory.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "coherence",
				Name:      "mode",
				Help:      "Current coherence mode (1=strong, 2=relaxed, 3=eventual)",
			},
			[]string{"mode"},
		),

		// Performance metrics
		ruleEvaluations: factory.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "performance",
				Name:      "rule_evaluations_total",
				Help:      "Total number of rule evaluations",
			},
		),

		ruleMatches: factory.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "performance",
				Name:      "rule_matches_total",
				Help:      "Total number of successful rule matches",
			},
		),

		ruleEvaluationTime: factory.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "performance",
				Name:      "rule_evaluation_duration_seconds",
				Help:      "Duration of rule evaluation",
				Buckets:   []float64{.00001, .00005, .0001, .0005, .001, .005, .01, .05, .1},
			},
		),

		// Arithmetic decomposition metrics
		arithmeticDecompositions: factory.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "arithmetic",
				Name:      "decompositions_total",
				Help:      "Total number of arithmetic decompositions",
			},
		),

		arithmeticOptimizations: factory.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "arithmetic",
				Name:      "optimizations_total",
				Help:      "Total number of arithmetic optimizations applied",
			},
		),
	}

	return exporter
}

// RecordIngestionMetrics enregistre les métriques d'une ingestion
func (e *PrometheusExporter) RecordIngestionMetrics(metrics *IngestionMetrics) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	mode := "full"
	if metrics.WasIncremental {
		mode = "incremental"
	}

	status := "success"

	// Record total duration
	e.ingestionDuration.WithLabelValues("total", mode).Observe(metrics.TotalDuration.Seconds())
	e.ingestionTotal.WithLabelValues(mode, status).Inc()

	// Record phase durations
	e.ingestionDuration.WithLabelValues("parsing", mode).Observe(metrics.ParsingDuration.Seconds())
	e.ingestionDuration.WithLabelValues("validation", mode).Observe(metrics.ValidationDuration.Seconds())
	e.ingestionDuration.WithLabelValues("type_creation", mode).Observe(metrics.TypeCreationDuration.Seconds())
	e.ingestionDuration.WithLabelValues("rule_creation", mode).Observe(metrics.RuleCreationDuration.Seconds())
	e.ingestionDuration.WithLabelValues("fact_collection", mode).Observe(metrics.FactCollectionDuration.Seconds())
	e.ingestionDuration.WithLabelValues("propagation", mode).Observe(metrics.PropagationDuration.Seconds())
	e.ingestionDuration.WithLabelValues("fact_submission", mode).Observe(metrics.FactSubmissionDuration.Seconds())

	// Record individual phase histograms
	e.parsingDuration.Observe(metrics.ParsingDuration.Seconds())
	e.validationDuration.Observe(metrics.ValidationDuration.Seconds())
	e.typeCreationDuration.Observe(metrics.TypeCreationDuration.Seconds())
	e.ruleCreationDuration.Observe(metrics.RuleCreationDuration.Seconds())
	e.factCollectionDuration.Observe(metrics.FactCollectionDuration.Seconds())
	e.propagationDuration.Observe(metrics.PropagationDuration.Seconds())
	e.factSubmissionDuration.Observe(metrics.FactSubmissionDuration.Seconds())

	// Record counters
	e.typesAdded.Add(float64(metrics.TypesAdded))
	e.rulesAdded.Add(float64(metrics.RulesAdded))
	e.factsSubmitted.Add(float64(metrics.FactsSubmitted))
	e.factsCollected.Add(float64(metrics.ExistingFactsCollected))
	e.factsPropagated.Add(float64(metrics.FactsPropagated))
	e.terminalsAdded.Add(float64(metrics.NewTerminalsAdded))
	e.propagationTargets.Set(float64(metrics.PropagationTargets))

	// Record network state
	e.typeNodes.Set(float64(metrics.TotalTypeNodes))
	e.terminalNodes.Set(float64(metrics.TotalTerminalNodes))
	e.alphaNodes.Set(float64(metrics.TotalAlphaNodes))
	e.betaNodes.Set(float64(metrics.TotalBetaNodes))

	// Calculate and record efficiency
	efficiencyScore := 1.0
	if !metrics.IsEfficient() {
		efficiencyScore = 0.5
	}
	e.efficiencyScore.Set(efficiencyScore)

	// Record bottleneck percentages
	if metrics.TotalDuration > 0 {
		phases := map[string]time.Duration{
			"parsing":         metrics.ParsingDuration,
			"validation":      metrics.ValidationDuration,
			"type_creation":   metrics.TypeCreationDuration,
			"rule_creation":   metrics.RuleCreationDuration,
			"fact_collection": metrics.FactCollectionDuration,
			"propagation":     metrics.PropagationDuration,
			"fact_submission": metrics.FactSubmissionDuration,
		}

		for phase, duration := range phases {
			percentage := float64(duration) / float64(metrics.TotalDuration) * 100
			e.bottleneckPhase.WithLabelValues(phase).Set(percentage)
		}
	}
}

// RecordTransaction enregistre une opération de transaction
func (e *PrometheusExporter) RecordTransaction(operation string, duration time.Duration, success bool) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	status := "success"
	if !success {
		status = "error"
	}

	e.transactionTotal.WithLabelValues(operation, status).Inc()
	e.transactionDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordCommand enregistre l'exécution d'une commande
func (e *PrometheusExporter) RecordCommand() {
	e.commandsExecuted.Inc()
}

// RecordRollback enregistre un rollback
func (e *PrometheusExporter) RecordRollback() {
	e.rollbacksTotal.Inc()
}

// RecordStorageOperation enregistre une opération de storage
func (e *PrometheusExporter) RecordStorageOperation(operation, opType string) {
	e.storageOperations.WithLabelValues(operation, opType).Inc()
}

// RecordStorageError enregistre une erreur de storage
func (e *PrometheusExporter) RecordStorageError(operation, errorType string) {
	e.storageErrors.WithLabelValues(operation, errorType).Inc()
}

// SetStorageFacts met à jour le nombre de faits en storage
func (e *PrometheusExporter) SetStorageFacts(count int) {
	e.storageFacts.Set(float64(count))
}

// RecordCoherenceViolation enregistre une violation de cohérence
func (e *PrometheusExporter) RecordCoherenceViolation(violationType string) {
	e.coherenceViolations.WithLabelValues(violationType).Inc()
}

// RecordCoherenceCheck enregistre une vérification de cohérence
func (e *PrometheusExporter) RecordCoherenceCheck() {
	e.coherenceChecks.Inc()
}

// SetCoherenceMode définit le mode de cohérence actuel
func (e *PrometheusExporter) SetCoherenceMode(mode string) {
	// Reset all modes
	e.coherenceMode.WithLabelValues("strong").Set(0)
	e.coherenceMode.WithLabelValues("relaxed").Set(0)
	e.coherenceMode.WithLabelValues("eventual").Set(0)

	// Set current mode
	switch mode {
	case "strong":
		e.coherenceMode.WithLabelValues("strong").Set(1)
	case "relaxed":
		e.coherenceMode.WithLabelValues("relaxed").Set(2)
	case "eventual":
		e.coherenceMode.WithLabelValues("eventual").Set(3)
	}
}

// RecordRuleEvaluation enregistre une évaluation de règle
func (e *PrometheusExporter) RecordRuleEvaluation(duration time.Duration, matched bool) {
	e.ruleEvaluations.Inc()
	if matched {
		e.ruleMatches.Inc()
	}
	e.ruleEvaluationTime.Observe(duration.Seconds())
}

// RecordArithmeticDecomposition enregistre une décomposition arithmétique
func (e *PrometheusExporter) RecordArithmeticDecomposition() {
	e.arithmeticDecompositions.Inc()
}

// RecordArithmeticOptimization enregistre une optimisation arithmétique
func (e *PrometheusExporter) RecordArithmeticOptimization() {
	e.arithmeticOptimizations.Inc()
}

// Registry retourne le registre Prometheus
func (e *PrometheusExporter) Registry() *prometheus.Registry {
	return e.registry
}

// IngestionMetrics représente les métriques d'ingestion (copie de rete package pour éviter dépendance circulaire)
type IngestionMetrics struct {
	ParsingDuration        time.Duration
	ValidationDuration     time.Duration
	TypeCreationDuration   time.Duration
	RuleCreationDuration   time.Duration
	FactCollectionDuration time.Duration
	PropagationDuration    time.Duration
	FactSubmissionDuration time.Duration
	TotalDuration          time.Duration

	TypesAdded             int
	RulesAdded             int
	FactsSubmitted         int
	ExistingFactsCollected int
	FactsPropagated        int
	NewTerminalsAdded      int
	PropagationTargets     int

	WasReset          bool
	WasIncremental    bool
	ValidationSkipped bool
	StartTime         time.Time
	EndTime           time.Time

	TotalTypeNodes     int
	TotalTerminalNodes int
	TotalAlphaNodes    int
	TotalBetaNodes     int
}

// IsEfficient vérifie si l'ingestion est efficace
func (m *IngestionMetrics) IsEfficient() bool {
	if m.NewTerminalsAdded > 0 && m.ExistingFactsCollected > 0 {
		maxPropagations := m.ExistingFactsCollected * m.NewTerminalsAdded
		if m.FactsPropagated > maxPropagations {
			return false
		}

		if m.TotalDuration > 0 && m.TotalDuration > 1*time.Millisecond {
			propagationRatio := float64(m.PropagationDuration) / float64(m.TotalDuration)
			if propagationRatio > 0.3 {
				return false
			}
		}
	}

	return true
}
