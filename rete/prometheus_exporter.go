// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// PrometheusExporter exporte les métriques RETE vers Prometheus
type PrometheusExporter struct {
	alphaMetrics *ChainBuildMetrics
	betaMetrics  *BetaChainMetrics
	config       *ChainPerformanceConfig
	mutex        sync.RWMutex

	// Registres des métriques
	registry map[string]*prometheusMetric
}

// prometheusMetric représente une métrique Prometheus
type prometheusMetric struct {
	name       string
	help       string
	metricType string // "counter", "gauge", "histogram"
	value      float64
	labels     map[string]string
}

// NewPrometheusExporter crée un nouveau exporteur Prometheus pour les métriques alpha
func NewPrometheusExporter(metrics *ChainBuildMetrics, config *ChainPerformanceConfig) *PrometheusExporter {
	if config == nil {
		config = DefaultChainPerformanceConfig()
	}

	return &PrometheusExporter{
		alphaMetrics: metrics,
		config:       config,
		registry:     make(map[string]*prometheusMetric),
	}
}

// NewPrometheusExporterWithBeta crée un exporteur avec métriques alpha et beta
func NewPrometheusExporterWithBeta(alphaMetrics *ChainBuildMetrics, betaMetrics *BetaChainMetrics, config *ChainPerformanceConfig) *PrometheusExporter {
	if config == nil {
		config = DefaultChainPerformanceConfig()
	}

	return &PrometheusExporter{
		alphaMetrics: alphaMetrics,
		betaMetrics:  betaMetrics,
		config:       config,
		registry:     make(map[string]*prometheusMetric),
	}
}

// RegisterMetrics enregistre toutes les métriques Prometheus
func (pe *PrometheusExporter) RegisterMetrics() {
	prefix := pe.config.PrometheusPrefix

	// Métriques de chaînes alpha
	pe.registerMetric(fmt.Sprintf("%s_alpha_chains_built_total", prefix),
		"Total number of alpha chains built",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_alpha_chains_length_avg", prefix),
		"Average length of alpha chains",
		"gauge")

	// Métriques de nœuds alpha
	pe.registerMetric(fmt.Sprintf("%s_alpha_nodes_created_total", prefix),
		"Total number of alpha nodes created",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_alpha_nodes_reused_total", prefix),
		"Total number of alpha nodes reused",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_alpha_nodes_sharing_ratio", prefix),
		"Ratio of alpha node sharing (0.0 to 1.0)",
		"gauge")

	// Métriques de cache de hash alpha
	pe.registerMetric(fmt.Sprintf("%s_alpha_hash_cache_hits_total", prefix),
		"Total number of alpha hash cache hits",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_alpha_hash_cache_misses_total", prefix),
		"Total number of alpha hash cache misses",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_alpha_hash_cache_size", prefix),
		"Current size of alpha hash cache",
		"gauge")

	pe.registerMetric(fmt.Sprintf("%s_alpha_hash_cache_efficiency", prefix),
		"Alpha hash cache efficiency (0.0 to 1.0)",
		"gauge")

	// Métriques de cache de connexion alpha
	pe.registerMetric(fmt.Sprintf("%s_alpha_connection_cache_hits_total", prefix),
		"Total number of alpha connection cache hits",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_alpha_connection_cache_misses_total", prefix),
		"Total number of alpha connection cache misses",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_alpha_connection_cache_efficiency", prefix),
		"Alpha connection cache efficiency (0.0 to 1.0)",
		"gauge")

	// Métriques de temps alpha
	pe.registerMetric(fmt.Sprintf("%s_alpha_build_time_seconds_total", prefix),
		"Total time spent building alpha chains in seconds",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_alpha_build_time_seconds_avg", prefix),
		"Average time spent building an alpha chain in seconds",
		"gauge")

	pe.registerMetric(fmt.Sprintf("%s_alpha_hash_compute_time_seconds_total", prefix),
		"Total time spent computing alpha hashes in seconds",
		"counter")

	// Métriques de chaînes beta (si disponibles)
	if pe.betaMetrics != nil {
		pe.registerMetric(fmt.Sprintf("%s_beta_chains_built_total", prefix),
			"Total number of beta chains built",
			"counter")

		pe.registerMetric(fmt.Sprintf("%s_beta_chains_length_avg", prefix),
			"Average length of beta chains",
			"gauge")

		// Métriques de nœuds beta
		pe.registerMetric(fmt.Sprintf("%s_beta_nodes_created_total", prefix),
			"Total number of beta nodes created",
			"counter")

		pe.registerMetric(fmt.Sprintf("%s_beta_nodes_reused_total", prefix),
			"Total number of beta nodes reused",
			"counter")

		pe.registerMetric(fmt.Sprintf("%s_beta_nodes_sharing_ratio", prefix),
			"Ratio of beta node sharing (0.0 to 1.0)",
			"gauge")

		// Métriques de jointures
		pe.registerMetric(fmt.Sprintf("%s_beta_joins_executed_total", prefix),
			"Total number of beta joins executed",
			"counter")

		pe.registerMetric(fmt.Sprintf("%s_beta_joins_time_seconds_avg", prefix),
			"Average time per beta join in seconds",
			"gauge")

		pe.registerMetric(fmt.Sprintf("%s_beta_joins_selectivity_avg", prefix),
			"Average beta join selectivity (0.0 to 1.0)",
			"gauge")

		pe.registerMetric(fmt.Sprintf("%s_beta_joins_result_size_avg", prefix),
			"Average beta join result size",
			"gauge")

		// Métriques de cache de hash beta
		pe.registerMetric(fmt.Sprintf("%s_beta_hash_cache_hits_total", prefix),
			"Total number of beta hash cache hits",
			"counter")

		pe.registerMetric(fmt.Sprintf("%s_beta_hash_cache_misses_total", prefix),
			"Total number of beta hash cache misses",
			"counter")

		pe.registerMetric(fmt.Sprintf("%s_beta_hash_cache_size", prefix),
			"Current size of beta hash cache",
			"gauge")

		pe.registerMetric(fmt.Sprintf("%s_beta_hash_cache_efficiency", prefix),
			"Beta hash cache efficiency (0.0 to 1.0)",
			"gauge")

		// Métriques de cache de jointure
		pe.registerMetric(fmt.Sprintf("%s_beta_join_cache_hits_total", prefix),
			"Total number of beta join cache hits",
			"counter")

		pe.registerMetric(fmt.Sprintf("%s_beta_join_cache_misses_total", prefix),
			"Total number of beta join cache misses",
			"counter")

		pe.registerMetric(fmt.Sprintf("%s_beta_join_cache_size", prefix),
			"Current size of beta join cache",
			"gauge")

		pe.registerMetric(fmt.Sprintf("%s_beta_join_cache_evictions_total", prefix),
			"Total number of beta join cache evictions",
			"counter")

		pe.registerMetric(fmt.Sprintf("%s_beta_join_cache_efficiency", prefix),
			"Beta join cache efficiency (0.0 to 1.0)",
			"gauge")

		// Métriques de cache de connexion beta
		pe.registerMetric(fmt.Sprintf("%s_beta_connection_cache_hits_total", prefix),
			"Total number of beta connection cache hits",
			"counter")

		pe.registerMetric(fmt.Sprintf("%s_beta_connection_cache_misses_total", prefix),
			"Total number of beta connection cache misses",
			"counter")

		pe.registerMetric(fmt.Sprintf("%s_beta_connection_cache_efficiency", prefix),
			"Beta connection cache efficiency (0.0 to 1.0)",
			"gauge")

		// Métriques de cache de préfixe beta
		pe.registerMetric(fmt.Sprintf("%s_beta_prefix_cache_hits_total", prefix),
			"Total number of beta prefix cache hits",
			"counter")

		pe.registerMetric(fmt.Sprintf("%s_beta_prefix_cache_misses_total", prefix),
			"Total number of beta prefix cache misses",
			"counter")

		pe.registerMetric(fmt.Sprintf("%s_beta_prefix_cache_size", prefix),
			"Current size of beta prefix cache",
			"gauge")

		pe.registerMetric(fmt.Sprintf("%s_beta_prefix_cache_efficiency", prefix),
			"Beta prefix cache efficiency (0.0 to 1.0)",
			"gauge")

		// Métriques de temps beta
		pe.registerMetric(fmt.Sprintf("%s_beta_build_time_seconds_total", prefix),
			"Total time spent building beta chains in seconds",
			"counter")

		pe.registerMetric(fmt.Sprintf("%s_beta_build_time_seconds_avg", prefix),
			"Average time spent building a beta chain in seconds",
			"gauge")

		pe.registerMetric(fmt.Sprintf("%s_beta_hash_compute_time_seconds_total", prefix),
			"Total time spent computing beta hashes in seconds",
			"counter")
	}
}

// registerMetric enregistre une métrique
func (pe *PrometheusExporter) registerMetric(name, help, metricType string) {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()

	pe.registry[name] = &prometheusMetric{
		name:       name,
		help:       help,
		metricType: metricType,
		labels:     make(map[string]string),
	}
}

// UpdateMetrics met à jour toutes les métriques avec les valeurs actuelles
func (pe *PrometheusExporter) UpdateMetrics() {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()

	prefix := pe.config.PrometheusPrefix

	// Métriques alpha
	if pe.alphaMetrics != nil {
		snapshot := pe.alphaMetrics.GetSnapshot()

		// Chaînes alpha
		pe.updateValue(fmt.Sprintf("%s_alpha_chains_built_total", prefix), float64(snapshot.TotalChainsBuilt))
		pe.updateValue(fmt.Sprintf("%s_alpha_chains_length_avg", prefix), snapshot.AverageChainLength)

		// Nœuds alpha
		pe.updateValue(fmt.Sprintf("%s_alpha_nodes_created_total", prefix), float64(snapshot.TotalNodesCreated))
		pe.updateValue(fmt.Sprintf("%s_alpha_nodes_reused_total", prefix), float64(snapshot.TotalNodesReused))
		pe.updateValue(fmt.Sprintf("%s_alpha_nodes_sharing_ratio", prefix), snapshot.SharingRatio)

		// Cache de hash alpha
		pe.updateValue(fmt.Sprintf("%s_alpha_hash_cache_hits_total", prefix), float64(snapshot.HashCacheHits))
		pe.updateValue(fmt.Sprintf("%s_alpha_hash_cache_misses_total", prefix), float64(snapshot.HashCacheMisses))
		pe.updateValue(fmt.Sprintf("%s_alpha_hash_cache_size", prefix), float64(snapshot.HashCacheSize))
		pe.updateValue(fmt.Sprintf("%s_alpha_hash_cache_efficiency", prefix), pe.alphaMetrics.GetHashCacheEfficiency())

		// Cache de connexion alpha
		pe.updateValue(fmt.Sprintf("%s_alpha_connection_cache_hits_total", prefix), float64(snapshot.ConnectionCacheHits))
		pe.updateValue(fmt.Sprintf("%s_alpha_connection_cache_misses_total", prefix), float64(snapshot.ConnectionCacheMisses))
		pe.updateValue(fmt.Sprintf("%s_alpha_connection_cache_efficiency", prefix), pe.alphaMetrics.GetConnectionCacheEfficiency())

		// Temps alpha
		pe.updateValue(fmt.Sprintf("%s_alpha_build_time_seconds_total", prefix), snapshot.TotalBuildTime.Seconds())
		pe.updateValue(fmt.Sprintf("%s_alpha_build_time_seconds_avg", prefix), snapshot.AverageBuildTime.Seconds())
		pe.updateValue(fmt.Sprintf("%s_alpha_hash_compute_time_seconds_total", prefix), snapshot.TotalHashComputeTime.Seconds())
	}

	// Métriques beta
	if pe.betaMetrics != nil {
		snapshot := pe.betaMetrics.GetSnapshot()

		// Chaînes beta
		pe.updateValue(fmt.Sprintf("%s_beta_chains_built_total", prefix), float64(snapshot.TotalChainsBuilt))
		pe.updateValue(fmt.Sprintf("%s_beta_chains_length_avg", prefix), snapshot.AverageChainLength)

		// Nœuds beta
		pe.updateValue(fmt.Sprintf("%s_beta_nodes_created_total", prefix), float64(snapshot.TotalNodesCreated))
		pe.updateValue(fmt.Sprintf("%s_beta_nodes_reused_total", prefix), float64(snapshot.TotalNodesReused))
		pe.updateValue(fmt.Sprintf("%s_beta_nodes_sharing_ratio", prefix), snapshot.SharingRatio)

		// Jointures
		pe.updateValue(fmt.Sprintf("%s_beta_joins_executed_total", prefix), float64(snapshot.TotalJoinsExecuted))
		pe.updateValue(fmt.Sprintf("%s_beta_joins_time_seconds_avg", prefix), snapshot.AverageJoinTime.Seconds())
		pe.updateValue(fmt.Sprintf("%s_beta_joins_selectivity_avg", prefix), snapshot.AverageJoinSelectivity)
		pe.updateValue(fmt.Sprintf("%s_beta_joins_result_size_avg", prefix), snapshot.AverageResultSize)

		// Cache de hash beta
		pe.updateValue(fmt.Sprintf("%s_beta_hash_cache_hits_total", prefix), float64(snapshot.HashCacheHits))
		pe.updateValue(fmt.Sprintf("%s_beta_hash_cache_misses_total", prefix), float64(snapshot.HashCacheMisses))
		pe.updateValue(fmt.Sprintf("%s_beta_hash_cache_size", prefix), float64(snapshot.HashCacheSize))
		pe.updateValue(fmt.Sprintf("%s_beta_hash_cache_efficiency", prefix), pe.betaMetrics.GetHashCacheEfficiency())

		// Cache de jointure
		pe.updateValue(fmt.Sprintf("%s_beta_join_cache_hits_total", prefix), float64(snapshot.JoinCacheHits))
		pe.updateValue(fmt.Sprintf("%s_beta_join_cache_misses_total", prefix), float64(snapshot.JoinCacheMisses))
		pe.updateValue(fmt.Sprintf("%s_beta_join_cache_size", prefix), float64(snapshot.JoinCacheSize))
		pe.updateValue(fmt.Sprintf("%s_beta_join_cache_evictions_total", prefix), float64(snapshot.JoinCacheEvictions))
		pe.updateValue(fmt.Sprintf("%s_beta_join_cache_efficiency", prefix), pe.betaMetrics.GetJoinCacheEfficiency())

		// Cache de connexion beta
		pe.updateValue(fmt.Sprintf("%s_beta_connection_cache_hits_total", prefix), float64(snapshot.ConnectionCacheHits))
		pe.updateValue(fmt.Sprintf("%s_beta_connection_cache_misses_total", prefix), float64(snapshot.ConnectionCacheMisses))
		pe.updateValue(fmt.Sprintf("%s_beta_connection_cache_efficiency", prefix), pe.betaMetrics.GetConnectionCacheEfficiency())

		// Cache de préfixe beta
		pe.updateValue(fmt.Sprintf("%s_beta_prefix_cache_hits_total", prefix), float64(snapshot.PrefixCacheHits))
		pe.updateValue(fmt.Sprintf("%s_beta_prefix_cache_misses_total", prefix), float64(snapshot.PrefixCacheMisses))
		pe.updateValue(fmt.Sprintf("%s_beta_prefix_cache_size", prefix), float64(snapshot.PrefixCacheSize))
		pe.updateValue(fmt.Sprintf("%s_beta_prefix_cache_efficiency", prefix), pe.betaMetrics.GetPrefixCacheEfficiency())

		// Temps beta
		pe.updateValue(fmt.Sprintf("%s_beta_build_time_seconds_total", prefix), snapshot.TotalBuildTime.Seconds())
		pe.updateValue(fmt.Sprintf("%s_beta_build_time_seconds_avg", prefix), snapshot.AverageBuildTime.Seconds())
		pe.updateValue(fmt.Sprintf("%s_beta_hash_compute_time_seconds_total", prefix), snapshot.TotalHashComputeTime.Seconds())
	}
}

// updateValue met à jour la valeur d'une métrique
func (pe *PrometheusExporter) updateValue(name string, value float64) {
	if metric, exists := pe.registry[name]; exists {
		metric.value = value
	}
}

// Handler retourne un http.Handler pour le endpoint /metrics
func (pe *PrometheusExporter) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mettre à jour les métriques avant de les exporter
		pe.UpdateMetrics()

		// Écrire les métriques au format Prometheus
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")

		pe.mutex.RLock()
		defer pe.mutex.RUnlock()

		for _, metric := range pe.registry {
			// HELP
			fmt.Fprintf(w, "# HELP %s %s\n", metric.name, metric.help)
			// TYPE
			fmt.Fprintf(w, "# TYPE %s %s\n", metric.name, metric.metricType)
			// VALUE
			if len(metric.labels) > 0 {
				fmt.Fprintf(w, "%s{%s} %v\n", metric.name, pe.formatLabels(metric.labels), metric.value)
			} else {
				fmt.Fprintf(w, "%s %v\n", metric.name, metric.value)
			}
		}
	})
}

// formatLabels formate les labels au format Prometheus
func (pe *PrometheusExporter) formatLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return ""
	}

	result := ""
	first := true
	for key, value := range labels {
		if !first {
			result += ","
		}
		result += fmt.Sprintf("%s=\"%s\"", key, value)
		first = false
	}
	return result
}

// ServeHTTP démarre un serveur HTTP pour exposer les métriques
func (pe *PrometheusExporter) ServeHTTP(addr string) error {
	http.Handle("/metrics", pe.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	return http.ListenAndServe(addr, nil)
}

// StartAutoUpdate démarre une goroutine qui met à jour automatiquement les métriques
func (pe *PrometheusExporter) StartAutoUpdate(interval time.Duration) chan struct{} {
	stop := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				pe.UpdateMetrics()
			case <-stop:
				return
			}
		}
	}()

	return stop
}

// GetMetricsText retourne les métriques au format texte Prometheus
func (pe *PrometheusExporter) GetMetricsText() string {
	pe.UpdateMetrics()

	pe.mutex.RLock()
	defer pe.mutex.RUnlock()

	result := ""
	for _, metric := range pe.registry {
		result += fmt.Sprintf("# HELP %s %s\n", metric.name, metric.help)
		result += fmt.Sprintf("# TYPE %s %s\n", metric.name, metric.metricType)
		if len(metric.labels) > 0 {
			result += fmt.Sprintf("%s{%s} %v\n", metric.name, pe.formatLabels(metric.labels), metric.value)
		} else {
			result += fmt.Sprintf("%s %v\n", metric.name, metric.value)
		}
	}
	return result
}

// PrometheusMetricsSnapshot retourne un snapshot des métriques pour Prometheus
type PrometheusMetricsSnapshot struct {
	Timestamp time.Time              `json:"timestamp"`
	Metrics   map[string]interface{} `json:"metrics"`
}

// GetSnapshot retourne un snapshot JSON des métriques
func (pe *PrometheusExporter) GetSnapshot() PrometheusMetricsSnapshot {
	pe.UpdateMetrics()

	pe.mutex.RLock()
	defer pe.mutex.RUnlock()

	metrics := make(map[string]interface{})
	for name, metric := range pe.registry {
		metrics[name] = map[string]interface{}{
			"value":  metric.value,
			"type":   metric.metricType,
			"help":   metric.help,
			"labels": metric.labels,
		}
	}

	return PrometheusMetricsSnapshot{
		Timestamp: time.Now(),
		Metrics:   metrics,
	}
}
