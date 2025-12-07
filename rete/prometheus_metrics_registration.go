// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import "fmt"

// prometheus_metrics_registration.go contient des fonctions helper pour l'enregistrement
// des métriques Prometheus. Ces fonctions ont été extraites de RegisterMetrics() pour
// améliorer la lisibilité et la maintenabilité.

// registerAlphaChainMetrics enregistre les métriques liées aux chaînes alpha
func (pe *PrometheusExporter) registerAlphaChainMetrics(prefix string) {
	pe.registerMetric(fmt.Sprintf("%s_alpha_chains_built_total", prefix),
		"Total number of alpha chains built",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_alpha_chains_length_avg", prefix),
		"Average length of alpha chains",
		"gauge")
}

// registerAlphaNodeMetrics enregistre les métriques liées aux nœuds alpha
func (pe *PrometheusExporter) registerAlphaNodeMetrics(prefix string) {
	pe.registerMetric(fmt.Sprintf("%s_alpha_nodes_created_total", prefix),
		"Total number of alpha nodes created",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_alpha_nodes_reused_total", prefix),
		"Total number of alpha nodes reused",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_alpha_nodes_sharing_ratio", prefix),
		"Ratio of alpha node sharing (0.0 to 1.0)",
		"gauge")
}

// registerAlphaHashCacheMetrics enregistre les métriques du cache de hash alpha
func (pe *PrometheusExporter) registerAlphaHashCacheMetrics(prefix string) {
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
}

// registerAlphaConnectionCacheMetrics enregistre les métriques du cache de connexion alpha
func (pe *PrometheusExporter) registerAlphaConnectionCacheMetrics(prefix string) {
	pe.registerMetric(fmt.Sprintf("%s_alpha_connection_cache_hits_total", prefix),
		"Total number of alpha connection cache hits",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_alpha_connection_cache_misses_total", prefix),
		"Total number of alpha connection cache misses",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_alpha_connection_cache_efficiency", prefix),
		"Alpha connection cache efficiency (0.0 to 1.0)",
		"gauge")
}

// registerAlphaTimeMetrics enregistre les métriques de temps alpha
func (pe *PrometheusExporter) registerAlphaTimeMetrics(prefix string) {
	pe.registerMetric(fmt.Sprintf("%s_alpha_build_time_seconds_total", prefix),
		"Total time spent building alpha chains in seconds",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_alpha_build_time_seconds_avg", prefix),
		"Average time spent building an alpha chain in seconds",
		"gauge")

	pe.registerMetric(fmt.Sprintf("%s_alpha_hash_compute_time_seconds_total", prefix),
		"Total time spent computing alpha hashes in seconds",
		"counter")
}

// registerBetaChainMetrics enregistre les métriques liées aux chaînes beta
func (pe *PrometheusExporter) registerBetaChainMetrics(prefix string) {
	pe.registerMetric(fmt.Sprintf("%s_beta_chains_built_total", prefix),
		"Total number of beta chains built",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_beta_chains_length_avg", prefix),
		"Average length of beta chains",
		"gauge")
}

// registerBetaNodeMetrics enregistre les métriques liées aux nœuds beta
func (pe *PrometheusExporter) registerBetaNodeMetrics(prefix string) {
	pe.registerMetric(fmt.Sprintf("%s_beta_nodes_created_total", prefix),
		"Total number of beta nodes created",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_beta_nodes_reused_total", prefix),
		"Total number of beta nodes reused",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_beta_nodes_sharing_ratio", prefix),
		"Ratio of beta node sharing (0.0 to 1.0)",
		"gauge")
}

// registerBetaJoinMetrics enregistre les métriques de jointures beta
func (pe *PrometheusExporter) registerBetaJoinMetrics(prefix string) {
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
}

// registerBetaHashCacheMetrics enregistre les métriques du cache de hash beta
func (pe *PrometheusExporter) registerBetaHashCacheMetrics(prefix string) {
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
}

// registerBetaJoinCacheMetrics enregistre les métriques du cache de jointure beta
func (pe *PrometheusExporter) registerBetaJoinCacheMetrics(prefix string) {
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
}

// registerBetaConnectionCacheMetrics enregistre les métriques du cache de connexion beta
func (pe *PrometheusExporter) registerBetaConnectionCacheMetrics(prefix string) {
	pe.registerMetric(fmt.Sprintf("%s_beta_connection_cache_hits_total", prefix),
		"Total number of beta connection cache hits",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_beta_connection_cache_misses_total", prefix),
		"Total number of beta connection cache misses",
		"counter")

	pe.registerMetric(fmt.Sprintf("%s_beta_connection_cache_efficiency", prefix),
		"Beta connection cache efficiency (0.0 to 1.0)",
		"gauge")
}

// registerBetaPrefixCacheMetrics enregistre les métriques du cache de préfixe beta
func (pe *PrometheusExporter) registerBetaPrefixCacheMetrics(prefix string) {
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
}

// registerBetaTimeMetrics enregistre les métriques de temps beta
func (pe *PrometheusExporter) registerBetaTimeMetrics(prefix string) {
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

// registerAlphaMetrics enregistre toutes les métriques alpha
func (pe *PrometheusExporter) registerAlphaMetrics(prefix string) {
	pe.registerAlphaChainMetrics(prefix)
	pe.registerAlphaNodeMetrics(prefix)
	pe.registerAlphaHashCacheMetrics(prefix)
	pe.registerAlphaConnectionCacheMetrics(prefix)
	pe.registerAlphaTimeMetrics(prefix)
}

// registerBetaMetrics enregistre toutes les métriques beta
func (pe *PrometheusExporter) registerBetaMetrics(prefix string) {
	pe.registerBetaChainMetrics(prefix)
	pe.registerBetaNodeMetrics(prefix)
	pe.registerBetaJoinMetrics(prefix)
	pe.registerBetaHashCacheMetrics(prefix)
	pe.registerBetaJoinCacheMetrics(prefix)
	pe.registerBetaConnectionCacheMetrics(prefix)
	pe.registerBetaPrefixCacheMetrics(prefix)
	pe.registerBetaTimeMetrics(prefix)
}
