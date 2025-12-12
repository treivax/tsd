// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"strings"
	"testing"
	"time"
)

func TestNewPrometheusExporterWithBeta(t *testing.T) {
	alphaMetrics := NewChainBuildMetrics()
	betaMetrics := NewBetaChainMetrics()
	config := DefaultChainPerformanceConfig()
	exporter := NewPrometheusExporterWithBeta(alphaMetrics, betaMetrics, config)
	if exporter == nil {
		t.Fatal("Expected non-nil exporter")
	}
	if exporter.alphaMetrics != alphaMetrics {
		t.Error("Expected alphaMetrics to match")
	}
	if exporter.betaMetrics != betaMetrics {
		t.Error("Expected betaMetrics to match")
	}
}
func TestPrometheusExporter_RegisterBetaMetrics(t *testing.T) {
	alphaMetrics := NewChainBuildMetrics()
	betaMetrics := NewBetaChainMetrics()
	config := DefaultChainPerformanceConfig()
	config.PrometheusPrefix = "test"
	exporter := NewPrometheusExporterWithBeta(alphaMetrics, betaMetrics, config)
	exporter.RegisterMetrics()
	// Verify beta metrics are registered
	expectedBetaMetrics := []string{
		"test_beta_chains_built_total",
		"test_beta_chains_length_avg",
		"test_beta_nodes_created_total",
		"test_beta_nodes_reused_total",
		"test_beta_nodes_sharing_ratio",
		"test_beta_joins_executed_total",
		"test_beta_joins_time_seconds_avg",
		"test_beta_joins_selectivity_avg",
		"test_beta_joins_result_size_avg",
		"test_beta_hash_cache_hits_total",
		"test_beta_hash_cache_misses_total",
		"test_beta_hash_cache_size",
		"test_beta_hash_cache_efficiency",
		"test_beta_join_cache_hits_total",
		"test_beta_join_cache_misses_total",
		"test_beta_join_cache_size",
		"test_beta_join_cache_evictions_total",
		"test_beta_join_cache_efficiency",
		"test_beta_connection_cache_hits_total",
		"test_beta_connection_cache_misses_total",
		"test_beta_connection_cache_efficiency",
		"test_beta_prefix_cache_hits_total",
		"test_beta_prefix_cache_misses_total",
		"test_beta_prefix_cache_size",
		"test_beta_prefix_cache_efficiency",
		"test_beta_build_time_seconds_total",
		"test_beta_build_time_seconds_avg",
		"test_beta_hash_compute_time_seconds_total",
	}
	exporter.mutex.RLock()
	defer exporter.mutex.RUnlock()
	for _, metricName := range expectedBetaMetrics {
		if _, exists := exporter.registry[metricName]; !exists {
			t.Errorf("Expected beta metric %s to be registered", metricName)
		}
	}
}
func TestPrometheusExporter_UpdateBetaMetrics(t *testing.T) {
	alphaMetrics := NewChainBuildMetrics()
	betaMetrics := NewBetaChainMetrics()
	config := DefaultChainPerformanceConfig()
	config.PrometheusPrefix = "test"
	// Populate beta metrics with data
	betaMetrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:       "rule1",
		ChainLength:  3,
		NodesCreated: 2,
		NodesReused:  1,
		BuildTime:    100 * time.Millisecond,
	})
	betaMetrics.RecordJoinExecution(10, 20, 50, 5*time.Millisecond)
	betaMetrics.RecordHashCacheHit()
	betaMetrics.RecordHashCacheMiss()
	betaMetrics.RecordJoinCacheHit()
	betaMetrics.UpdateJoinCacheSize(100)
	exporter := NewPrometheusExporterWithBeta(alphaMetrics, betaMetrics, config)
	exporter.RegisterMetrics()
	exporter.UpdateMetrics()
	exporter.mutex.RLock()
	defer exporter.mutex.RUnlock()
	// Verify beta metrics are updated
	if metric, exists := exporter.registry["test_beta_chains_built_total"]; exists {
		if metric.value != 1.0 {
			t.Errorf("Expected beta_chains_built_total=1.0, got %f", metric.value)
		}
	} else {
		t.Error("Expected beta_chains_built_total to exist")
	}
	if metric, exists := exporter.registry["test_beta_nodes_created_total"]; exists {
		if metric.value != 2.0 {
			t.Errorf("Expected beta_nodes_created_total=2.0, got %f", metric.value)
		}
	} else {
		t.Error("Expected beta_nodes_created_total to exist")
	}
	if metric, exists := exporter.registry["test_beta_joins_executed_total"]; exists {
		if metric.value != 1.0 {
			t.Errorf("Expected beta_joins_executed_total=1.0, got %f", metric.value)
		}
	} else {
		t.Error("Expected beta_joins_executed_total to exist")
	}
	if metric, exists := exporter.registry["test_beta_join_cache_size"]; exists {
		if metric.value != 100.0 {
			t.Errorf("Expected beta_join_cache_size=100.0, got %f", metric.value)
		}
	} else {
		t.Error("Expected beta_join_cache_size to exist")
	}
}
func TestPrometheusExporter_GetMetricsTextWithBeta(t *testing.T) {
	alphaMetrics := NewChainBuildMetrics()
	betaMetrics := NewBetaChainMetrics()
	config := DefaultChainPerformanceConfig()
	config.PrometheusPrefix = "rete"
	// Populate with data
	betaMetrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:       "rule1",
		ChainLength:  2,
		NodesCreated: 1,
		NodesReused:  1,
		BuildTime:    50 * time.Millisecond,
	})
	exporter := NewPrometheusExporterWithBeta(alphaMetrics, betaMetrics, config)
	exporter.RegisterMetrics()
	metricsText := exporter.GetMetricsText()
	// Verify beta metrics are present in the output
	expectedStrings := []string{
		"# HELP rete_beta_chains_built_total",
		"# TYPE rete_beta_chains_built_total counter",
		"rete_beta_chains_built_total 1",
		"# HELP rete_beta_nodes_created_total",
		"rete_beta_nodes_created_total 1",
		"# HELP rete_beta_nodes_reused_total",
		"rete_beta_nodes_reused_total 1",
		"# HELP rete_beta_nodes_sharing_ratio",
		"rete_beta_nodes_sharing_ratio 0.5",
	}
	for _, expected := range expectedStrings {
		if !strings.Contains(metricsText, expected) {
			t.Errorf("Expected metrics text to contain: %s", expected)
		}
	}
}
func TestPrometheusExporter_BetaCacheEfficiencyMetrics(t *testing.T) {
	alphaMetrics := NewChainBuildMetrics()
	betaMetrics := NewBetaChainMetrics()
	config := DefaultChainPerformanceConfig()
	config.PrometheusPrefix = "test"
	// Record cache activity
	betaMetrics.RecordHashCacheHit()
	betaMetrics.RecordHashCacheHit()
	betaMetrics.RecordHashCacheMiss()
	betaMetrics.RecordJoinCacheHit()
	betaMetrics.RecordJoinCacheHit()
	betaMetrics.RecordJoinCacheHit()
	betaMetrics.RecordJoinCacheMiss()
	betaMetrics.RecordConnectionCacheHit()
	betaMetrics.RecordConnectionCacheMiss()
	betaMetrics.RecordConnectionCacheMiss()
	betaMetrics.RecordPrefixCacheHit()
	betaMetrics.RecordPrefixCacheMiss()
	exporter := NewPrometheusExporterWithBeta(alphaMetrics, betaMetrics, config)
	exporter.RegisterMetrics()
	exporter.UpdateMetrics()
	exporter.mutex.RLock()
	defer exporter.mutex.RUnlock()
	// Hash cache efficiency: 2/3 = 0.666...
	if metric, exists := exporter.registry["test_beta_hash_cache_efficiency"]; exists {
		expected := 2.0 / 3.0
		if metric.value != expected {
			t.Errorf("Expected hash cache efficiency=%.4f, got %.4f", expected, metric.value)
		}
	}
	// Join cache efficiency: 3/4 = 0.75
	if metric, exists := exporter.registry["test_beta_join_cache_efficiency"]; exists {
		expected := 0.75
		if metric.value != expected {
			t.Errorf("Expected join cache efficiency=%.4f, got %.4f", expected, metric.value)
		}
	}
	// Connection cache efficiency: 1/3 = 0.333...
	if metric, exists := exporter.registry["test_beta_connection_cache_efficiency"]; exists {
		expected := 1.0 / 3.0
		if metric.value != expected {
			t.Errorf("Expected connection cache efficiency=%.4f, got %.4f", expected, metric.value)
		}
	}
	// Prefix cache efficiency: 1/2 = 0.5
	if metric, exists := exporter.registry["test_beta_prefix_cache_efficiency"]; exists {
		expected := 0.5
		if metric.value != expected {
			t.Errorf("Expected prefix cache efficiency=%.4f, got %.4f", expected, metric.value)
		}
	}
}
func TestPrometheusExporter_BetaJoinMetrics(t *testing.T) {
	alphaMetrics := NewChainBuildMetrics()
	betaMetrics := NewBetaChainMetrics()
	config := DefaultChainPerformanceConfig()
	config.PrometheusPrefix = "test"
	// Record join executions
	betaMetrics.RecordJoinExecution(10, 20, 50, 10*time.Millisecond)
	betaMetrics.RecordJoinExecution(5, 10, 25, 5*time.Millisecond)
	exporter := NewPrometheusExporterWithBeta(alphaMetrics, betaMetrics, config)
	exporter.RegisterMetrics()
	exporter.UpdateMetrics()
	exporter.mutex.RLock()
	defer exporter.mutex.RUnlock()
	// Verify join metrics
	if metric, exists := exporter.registry["test_beta_joins_executed_total"]; exists {
		if metric.value != 2.0 {
			t.Errorf("Expected joins_executed_total=2.0, got %f", metric.value)
		}
	}
	// Average join time: (10ms + 5ms) / 2 = 7.5ms = 0.0075s
	if metric, exists := exporter.registry["test_beta_joins_time_seconds_avg"]; exists {
		expected := 0.0075
		if metric.value != expected {
			t.Errorf("Expected avg join time=%.6f, got %.6f", expected, metric.value)
		}
	}
	// Average selectivity: (50/200 + 25/50) / 2 = (0.25 + 0.5) / 2 = 0.375
	if metric, exists := exporter.registry["test_beta_joins_selectivity_avg"]; exists {
		expected := 0.375
		if metric.value != expected {
			t.Errorf("Expected avg selectivity=%.4f, got %.4f", expected, metric.value)
		}
	}
	// Average result size: (50 + 25) / 2 = 37.5
	if metric, exists := exporter.registry["test_beta_joins_result_size_avg"]; exists {
		expected := 37.5
		if metric.value != expected {
			t.Errorf("Expected avg result size=%.2f, got %.2f", expected, metric.value)
		}
	}
}
func TestPrometheusExporter_AlphaAndBetaTogether(t *testing.T) {
	alphaMetrics := NewChainBuildMetrics()
	betaMetrics := NewBetaChainMetrics()
	config := DefaultChainPerformanceConfig()
	config.PrometheusPrefix = "rete"
	// Populate alpha metrics
	alphaMetrics.RecordChainBuild(ChainMetricDetail{
		RuleID:       "alpha_rule",
		ChainLength:  3,
		NodesCreated: 2,
		NodesReused:  1,
		BuildTime:    100 * time.Millisecond,
	})
	// Populate beta metrics
	betaMetrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:       "beta_rule",
		ChainLength:  2,
		NodesCreated: 1,
		NodesReused:  1,
		BuildTime:    50 * time.Millisecond,
	})
	exporter := NewPrometheusExporterWithBeta(alphaMetrics, betaMetrics, config)
	exporter.RegisterMetrics()
	metricsText := exporter.GetMetricsText()
	// Verify both alpha and beta metrics are present
	alphaStrings := []string{
		"rete_alpha_chains_built_total 1",
		"rete_alpha_nodes_created_total 2",
	}
	betaStrings := []string{
		"rete_beta_chains_built_total 1",
		"rete_beta_nodes_created_total 1",
	}
	for _, expected := range alphaStrings {
		if !strings.Contains(metricsText, expected) {
			t.Errorf("Expected alpha metric: %s", expected)
		}
	}
	for _, expected := range betaStrings {
		if !strings.Contains(metricsText, expected) {
			t.Errorf("Expected beta metric: %s", expected)
		}
	}
}
func TestPrometheusExporter_BetaMetricsWithoutAlpha(t *testing.T) {
	// Test exporter with only beta metrics (no alpha)
	betaMetrics := NewBetaChainMetrics()
	config := DefaultChainPerformanceConfig()
	config.PrometheusPrefix = "test"
	betaMetrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:       "rule1",
		ChainLength:  2,
		NodesCreated: 1,
		NodesReused:  1,
	})
	exporter := NewPrometheusExporterWithBeta(nil, betaMetrics, config)
	exporter.RegisterMetrics()
	exporter.UpdateMetrics()
	// Should not panic and should still export beta metrics
	metricsText := exporter.GetMetricsText()
	if !strings.Contains(metricsText, "test_beta_chains_built_total 1") {
		t.Error("Expected beta metrics to be exported even without alpha metrics")
	}
}
func TestPrometheusExporter_AlphaMetricsWithoutBeta(t *testing.T) {
	// Test exporter with only alpha metrics (no beta) - legacy mode
	alphaMetrics := NewChainBuildMetrics()
	config := DefaultChainPerformanceConfig()
	config.PrometheusPrefix = "test"
	alphaMetrics.RecordChainBuild(ChainMetricDetail{
		RuleID:       "rule1",
		ChainLength:  2,
		NodesCreated: 1,
		NodesReused:  1,
	})
	exporter := NewPrometheusExporter(alphaMetrics, config)
	exporter.RegisterMetrics()
	exporter.UpdateMetrics()
	metricsText := exporter.GetMetricsText()
	// Should have alpha metrics but not beta
	if !strings.Contains(metricsText, "test_alpha_chains_built_total 1") {
		t.Error("Expected alpha metrics to be exported")
	}
	if strings.Contains(metricsText, "beta") {
		t.Error("Should not contain beta metrics when exporter created without beta")
	}
}
