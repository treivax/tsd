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
	metrics *ChainBuildMetrics
	config  *ChainPerformanceConfig
	mutex   sync.RWMutex

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

// NewPrometheusExporter crée un nouveau exporteur Prometheus
func NewPrometheusExporter(metrics *ChainBuildMetrics, config *ChainPerformanceConfig) *PrometheusExporter {
	if config == nil {
		config = DefaultChainPerformanceConfig()
	}

	return &PrometheusExporter{
		metrics:  metrics,
		config:   config,
		registry: make(map[string]*prometheusMetric),
	}
}

// RegisterMetrics enregistre toutes les métriques Prometheus
func (pe *PrometheusExporter) RegisterMetrics() {
	prefix := pe.config.PrometheusPrefix

	// Métriques de chaînes
	pe.registerMetric(fmt.Sprintf("%s_chains_built_total", prefix),
		"Total number of alpha chains built",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_chains_length_avg", prefix),
		"Average length of alpha chains",
		"gauge")

	// Métriques de nœuds
	pe.registerMetric(fmt.Sprintf("%s_nodes_created_total", prefix),
		"Total number of alpha nodes created",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_nodes_reused_total", prefix),
		"Total number of alpha nodes reused",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_nodes_sharing_ratio", prefix),
		"Ratio of node sharing (0.0 to 1.0)",
		"gauge")

	// Métriques de cache de hash
	pe.registerMetric(fmt.Sprintf("%s_hash_cache_hits_total", prefix),
		"Total number of hash cache hits",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_hash_cache_misses_total", prefix),
		"Total number of hash cache misses",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_hash_cache_size", prefix),
		"Current size of hash cache",
		"gauge")

	pe.registerMetric(fmt.Sprintf("%s_hash_cache_efficiency", prefix),
		"Hash cache efficiency (0.0 to 1.0)",
		"gauge")

	// Métriques de cache de connexion
	pe.registerMetric(fmt.Sprintf("%s_connection_cache_hits_total", prefix),
		"Total number of connection cache hits",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_connection_cache_misses_total", prefix),
		"Total number of connection cache misses",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_connection_cache_efficiency", prefix),
		"Connection cache efficiency (0.0 to 1.0)",
		"gauge")

	// Métriques de temps
	pe.registerMetric(fmt.Sprintf("%s_build_time_seconds_total", prefix),
		"Total time spent building chains in seconds",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_build_time_seconds_avg", prefix),
		"Average time spent building a chain in seconds",
		"gauge")

	pe.registerMetric(fmt.Sprintf("%s_hash_compute_time_seconds_total", prefix),
		"Total time spent computing hashes in seconds",
		"counter")
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
	if pe.metrics == nil {
		return
	}

	snapshot := pe.metrics.GetSnapshot()
	prefix := pe.config.PrometheusPrefix

	pe.mutex.Lock()
	defer pe.mutex.Unlock()

	// Chaînes
	pe.updateValue(fmt.Sprintf("%s_chains_built_total", prefix), float64(snapshot.TotalChainsBuilt))
	pe.updateValue(fmt.Sprintf("%s_chains_length_avg", prefix), snapshot.AverageChainLength)

	// Nœuds
	pe.updateValue(fmt.Sprintf("%s_nodes_created_total", prefix), float64(snapshot.TotalNodesCreated))
	pe.updateValue(fmt.Sprintf("%s_nodes_reused_total", prefix), float64(snapshot.TotalNodesReused))
	pe.updateValue(fmt.Sprintf("%s_nodes_sharing_ratio", prefix), snapshot.SharingRatio)

	// Cache de hash
	pe.updateValue(fmt.Sprintf("%s_hash_cache_hits_total", prefix), float64(snapshot.HashCacheHits))
	pe.updateValue(fmt.Sprintf("%s_hash_cache_misses_total", prefix), float64(snapshot.HashCacheMisses))
	pe.updateValue(fmt.Sprintf("%s_hash_cache_size", prefix), float64(snapshot.HashCacheSize))
	pe.updateValue(fmt.Sprintf("%s_hash_cache_efficiency", prefix), pe.metrics.GetHashCacheEfficiency())

	// Cache de connexion
	pe.updateValue(fmt.Sprintf("%s_connection_cache_hits_total", prefix), float64(snapshot.ConnectionCacheHits))
	pe.updateValue(fmt.Sprintf("%s_connection_cache_misses_total", prefix), float64(snapshot.ConnectionCacheMisses))
	pe.updateValue(fmt.Sprintf("%s_connection_cache_efficiency", prefix), pe.metrics.GetConnectionCacheEfficiency())

	// Temps
	pe.updateValue(fmt.Sprintf("%s_build_time_seconds_total", prefix), snapshot.TotalBuildTime.Seconds())
	pe.updateValue(fmt.Sprintf("%s_build_time_seconds_avg", prefix), snapshot.AverageBuildTime.Seconds())
	pe.updateValue(fmt.Sprintf("%s_hash_compute_time_seconds_total", prefix), snapshot.TotalHashComputeTime.Seconds())
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
