// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// =============================================================================
// Benchmark: Chain Construction Performance
// =============================================================================
// BenchmarkBetaChainBuild_WithSharing benchmarks chain building with sharing enabled
func BenchmarkBetaChainBuild_WithSharing(b *testing.B) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	patterns := createSimilarPatterns(10)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ruleID := fmt.Sprintf("rule_%d", i%10) // Simulate 10 different rules
		_, err := builder.BuildChain(patterns, ruleID)
		if err != nil {
			b.Fatalf("BuildChain failed: %v", err)
		}
	}
	b.StopTimer()
	// Report metrics
	metrics := builder.GetMetrics()
	sharingRatio := float64(metrics.TotalNodesReused) / float64(metrics.TotalNodesCreated+metrics.TotalNodesReused) * 100
	b.ReportMetric(sharingRatio, "sharing_%")
	b.ReportMetric(float64(metrics.TotalNodesReused), "nodes_reused")
	b.ReportMetric(float64(metrics.TotalNodesCreated), "nodes_created")
}

// BenchmarkBetaChainBuild_WithoutSharing benchmarks chain building with sharing disabled
func BenchmarkBetaChainBuild_WithoutSharing(b *testing.B) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       false,
		HashCacheSize: 0,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	patterns := createSimilarPatterns(10)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ruleID := fmt.Sprintf("rule_%d", i%10)
		_, err := builder.BuildChain(patterns, ruleID)
		if err != nil {
			b.Fatalf("BuildChain failed: %v", err)
		}
	}
	b.StopTimer()
	// Report metrics
	metrics := builder.GetMetrics()
	b.ReportMetric(float64(metrics.TotalNodesCreated), "nodes_created")
}

// BenchmarkBetaChainBuild_SimilarPatterns_10Rules benchmarks with 10 similar rules
func BenchmarkBetaChainBuild_SimilarPatterns_10Rules(b *testing.B) {
	benchmarkWithRuleCount(b, 10, createSimilarPatterns, true)
}

// BenchmarkBetaChainBuild_SimilarPatterns_100Rules benchmarks with 100 similar rules
func BenchmarkBetaChainBuild_SimilarPatterns_100Rules(b *testing.B) {
	benchmarkWithRuleCount(b, 100, createSimilarPatterns, true)
}

// BenchmarkBetaChainBuild_MixedPatterns_10Rules benchmarks with 10 mixed rules
func BenchmarkBetaChainBuild_MixedPatterns_10Rules(b *testing.B) {
	benchmarkWithRuleCount(b, 10, createMixedPatterns, true)
}

// BenchmarkBetaChainBuild_MixedPatterns_100Rules benchmarks with 100 mixed rules
func BenchmarkBetaChainBuild_MixedPatterns_100Rules(b *testing.B) {
	benchmarkWithRuleCount(b, 100, createMixedPatterns, true)
}

// BenchmarkBetaChainBuild_ComplexRules benchmarks complex rules with 5+ joins
func BenchmarkBetaChainBuild_ComplexRules(b *testing.B) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 2000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	// Create complex patterns with 7 joins
	patterns := createComplexPatterns(7)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ruleID := fmt.Sprintf("complex_rule_%d", i%20)
		_, err := builder.BuildChain(patterns, ruleID)
		if err != nil {
			b.Fatalf("BuildChain failed: %v", err)
		}
	}
	b.StopTimer()
	reportBenchmarkMetrics(b, builder)
}

// =============================================================================
// Benchmark: Join Cache Performance
// =============================================================================
// BenchmarkJoinCache_Hits benchmarks join cache hit performance
func BenchmarkJoinCache_Hits(b *testing.B) {
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	config.BetaJoinResultCacheMaxSize = 1000
	cache := NewBetaJoinCache(config)
	// Pre-populate cache
	for i := 0; i < 100; i++ {
		fact := &Fact{
			ID:     fmt.Sprintf("fact_%d", i),
			Type:   "TestFact",
			Fields: map[string]interface{}{"id": i},
		}
		token := &Token{
			ID:    fmt.Sprintf("token_%d", i),
			Facts: []*Fact{fact},
		}
		joinNode := &JoinNode{
			BaseNode: BaseNode{ID: fmt.Sprintf("join_%d", i)},
		}
		result := &JoinResult{
			Matched: true,
			Token:   token,
		}
		cache.SetJoinResult(token, fact, joinNode, result)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fact := &Fact{
			ID:     fmt.Sprintf("fact_%d", i%100),
			Type:   "TestFact",
			Fields: map[string]interface{}{"id": i % 100},
		}
		token := &Token{
			ID:    fmt.Sprintf("token_%d", i%100),
			Facts: []*Fact{fact},
		}
		joinNode := &JoinNode{
			BaseNode: BaseNode{ID: fmt.Sprintf("join_%d", i%100)},
		}
		_, found := cache.GetJoinResult(token, fact, joinNode)
		if !found {
			b.Fatalf("Expected cache hit for token %d", i%100)
		}
	}
	b.StopTimer()
	hitRate := cache.GetHitRate() * 100
	b.ReportMetric(hitRate, "hit_rate_%")
}

// BenchmarkJoinCache_Misses benchmarks join cache miss performance
func BenchmarkJoinCache_Misses(b *testing.B) {
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	config.BetaJoinResultCacheMaxSize = 1000
	cache := NewBetaJoinCache(config)
	// Pre-populate cache with different keys
	for i := 0; i < 100; i++ {
		fact := &Fact{
			ID:     fmt.Sprintf("fact_%d", i),
			Type:   "TestFact",
			Fields: map[string]interface{}{"id": i},
		}
		token := &Token{
			ID:    fmt.Sprintf("token_%d", i),
			Facts: []*Fact{fact},
		}
		joinNode := &JoinNode{
			BaseNode: BaseNode{ID: fmt.Sprintf("join_%d", i)},
		}
		result := &JoinResult{
			Matched: true,
			Token:   token,
		}
		cache.SetJoinResult(token, fact, joinNode, result)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fact := &Fact{
			ID:     fmt.Sprintf("miss_fact_%d", i%100),
			Type:   "TestFact",
			Fields: map[string]interface{}{"id": i % 100},
		}
		token := &Token{
			ID:    fmt.Sprintf("miss_token_%d", i%100),
			Facts: []*Fact{fact},
		}
		joinNode := &JoinNode{
			BaseNode: BaseNode{ID: "miss_join"},
		}
		_, found := cache.GetJoinResult(token, fact, joinNode)
		if found {
			b.Fatal("Expected cache miss")
		}
	}
	b.StopTimer()
	stats := cache.GetStats()
	missRate := (1.0 - cache.GetHitRate()) * 100
	b.ReportMetric(missRate, "miss_rate_%")
	b.ReportMetric(float64(stats["misses"].(int64)), "misses")
}

// BenchmarkJoinCache_Evictions benchmarks cache behavior under eviction pressure
func BenchmarkJoinCache_Evictions(b *testing.B) {
	cacheSize := 100
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	config.BetaJoinResultCacheMaxSize = cacheSize
	cache := NewBetaJoinCache(config)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fact := &Fact{
			ID:     fmt.Sprintf("evict_fact_%d", i),
			Type:   "TestFact",
			Fields: map[string]interface{}{"id": i},
		}
		token := &Token{
			ID:    fmt.Sprintf("evict_token_%d", i),
			Facts: []*Fact{fact},
		}
		joinNode := &JoinNode{
			BaseNode: BaseNode{ID: fmt.Sprintf("evict_join_%d", i)},
		}
		result := &JoinResult{
			Matched: true,
			Token:   token,
		}
		cache.SetJoinResult(token, fact, joinNode, result)
	}
	b.StopTimer()
	stats := cache.GetStats()
	b.ReportMetric(float64(stats["evictions"].(int64)), "evictions")
	b.ReportMetric(float64(cache.GetSize()), "final_size")
}

// BenchmarkJoinCache_MixedWorkload benchmarks realistic mixed read/write workload
func BenchmarkJoinCache_MixedWorkload(b *testing.B) {
	config := DefaultChainPerformanceConfig()
	config.BetaCacheEnabled = true
	config.BetaJoinResultCacheEnabled = true
	config.BetaJoinResultCacheMaxSize = 500
	cache := NewBetaJoinCache(config)
	// Pre-populate with some data
	for i := 0; i < 100; i++ {
		fact := &Fact{
			ID:     fmt.Sprintf("fact_%d", i),
			Type:   "TestFact",
			Fields: map[string]interface{}{"id": i},
		}
		token := &Token{
			ID:    fmt.Sprintf("token_%d", i),
			Facts: []*Fact{fact},
		}
		joinNode := &JoinNode{
			BaseNode: BaseNode{ID: fmt.Sprintf("join_%d", i)},
		}
		result := &JoinResult{
			Matched: true,
			Token:   token,
		}
		cache.SetJoinResult(token, fact, joinNode, result)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// 70% reads, 30% writes
		if i%10 < 7 {
			fact := &Fact{
				ID:     fmt.Sprintf("fact_%d", i%150),
				Type:   "TestFact",
				Fields: map[string]interface{}{"id": i % 150},
			}
			token := &Token{
				ID:    fmt.Sprintf("token_%d", i%150),
				Facts: []*Fact{fact},
			}
			joinNode := &JoinNode{
				BaseNode: BaseNode{ID: fmt.Sprintf("join_%d", i%150)},
			}
			cache.GetJoinResult(token, fact, joinNode)
		} else {
			fact := &Fact{
				ID:     fmt.Sprintf("fact_%d", i),
				Type:   "TestFact",
				Fields: map[string]interface{}{"id": i},
			}
			token := &Token{
				ID:    fmt.Sprintf("token_%d", i),
				Facts: []*Fact{fact},
			}
			joinNode := &JoinNode{
				BaseNode: BaseNode{ID: fmt.Sprintf("join_%d", i)},
			}
			result := &JoinResult{
				Matched: true,
				Token:   token,
			}
			cache.SetJoinResult(token, fact, joinNode, result)
		}
	}
	b.StopTimer()
	stats := cache.GetStats()
	hitRate := cache.GetHitRate() * 100
	b.ReportMetric(hitRate, "hit_rate_%")
	b.ReportMetric(float64(stats["evictions"].(int64)), "evictions")
}

// =============================================================================
// Benchmark: Join Order Optimization
// =============================================================================
// BenchmarkJoinOrder_Optimal benchmarks pre-optimized join order
func BenchmarkJoinOrder_Optimal(b *testing.B) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	// Patterns already in optimal order (low to high selectivity)
	patterns := []JoinPattern{
		createPatternWithSelectivity("p", "o", 0.1),
		createPatternWithSelectivity("o", "pay", 0.3),
		createPatternWithSelectivity("pay", "ship", 0.5),
		createPatternWithSelectivity("ship", "inv", 0.7),
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ruleID := fmt.Sprintf("optimal_rule_%d", i)
		_, err := builder.BuildChain(patterns, ruleID)
		if err != nil {
			b.Fatalf("BuildChain failed: %v", err)
		}
	}
	b.StopTimer()
	reportBenchmarkMetrics(b, builder)
}

// BenchmarkJoinOrder_Suboptimal benchmarks suboptimal join order that needs reordering
func BenchmarkJoinOrder_Suboptimal(b *testing.B) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	// Patterns in suboptimal order (high to low selectivity)
	patterns := []JoinPattern{
		createPatternWithSelectivity("p", "o", 0.9),
		createPatternWithSelectivity("o", "pay", 0.7),
		createPatternWithSelectivity("pay", "ship", 0.3),
		createPatternWithSelectivity("ship", "inv", 0.1),
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ruleID := fmt.Sprintf("suboptimal_rule_%d", i)
		_, err := builder.BuildChain(patterns, ruleID)
		if err != nil {
			b.Fatalf("BuildChain failed: %v", err)
		}
	}
	b.StopTimer()
	reportBenchmarkMetrics(b, builder)
}

// =============================================================================
// Benchmark: Hash Computation
// =============================================================================
// BenchmarkHashCompute_Simple benchmarks simple condition hashing via GetOrCreateJoinNode
func BenchmarkHashCompute_Simple(b *testing.B) {
	storage := NewMemoryStorage()
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 0, // Disable caching
	}
	registry := NewBetaSharingRegistry(config, nil)
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
		"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "id"},
		"right":    map[string]interface{}{"type": "fieldAccess", "object": "o", "field": "customer_id"},
	}
	leftVars := []string{"p"}
	rightVars := []string{"o"}
	allVars := []string{"p", "o"}
	varTypes := map[string]string{"p": "Person", "o": "Order"}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, _, err := registry.GetOrCreateJoinNode(condition, leftVars, rightVars, allVars, varTypes, storage)
		if err != nil {
			b.Fatalf("GetOrCreateJoinNode failed: %v", err)
		}
	}
}

// BenchmarkHashCompute_Complex benchmarks complex condition hashing via GetOrCreateJoinNode
func BenchmarkHashCompute_Complex(b *testing.B) {
	storage := NewMemoryStorage()
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 0, // Disable caching
	}
	registry := NewBetaSharingRegistry(config, nil)
	condition := map[string]interface{}{
		"type":     "and",
		"operator": "&&",
		"left": map[string]interface{}{
			"type":     "comparison",
			"operator": "==",
			"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "id"},
			"right":    map[string]interface{}{"type": "fieldAccess", "object": "o", "field": "customer_id"},
		},
		"right": map[string]interface{}{
			"type":     "comparison",
			"operator": ">",
			"left":     map[string]interface{}{"type": "fieldAccess", "object": "o", "field": "amount"},
			"right":    map[string]interface{}{"type": "literal", "value": 100},
		},
	}
	leftVars := []string{"p"}
	rightVars := []string{"o"}
	allVars := []string{"p", "o"}
	varTypes := map[string]string{"p": "Person", "o": "Order"}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, _, err := registry.GetOrCreateJoinNode(condition, leftVars, rightVars, allVars, varTypes, storage)
		if err != nil {
			b.Fatalf("GetOrCreateJoinNode failed: %v", err)
		}
	}
}

// BenchmarkHashCompute_WithCache benchmarks hashing with cache enabled via GetOrCreateJoinNode
func BenchmarkHashCompute_WithCache(b *testing.B) {
	storage := NewMemoryStorage()
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
	}
	registry := NewBetaSharingRegistry(config, nil)
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
		"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "id"},
		"right":    map[string]interface{}{"type": "fieldAccess", "object": "o", "field": "customer_id"},
	}
	leftVars := []string{"p"}
	rightVars := []string{"o"}
	allVars := []string{"p", "o"}
	varTypes := map[string]string{"p": "Person", "o": "Order"}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, _, err := registry.GetOrCreateJoinNode(condition, leftVars, rightVars, allVars, varTypes, storage)
		if err != nil {
			b.Fatalf("GetOrCreateJoinNode failed: %v", err)
		}
	}
	b.StopTimer()
	stats := registry.GetSharingStats()
	if stats.HashCacheHitRate > 0 {
		hitRate := stats.HashCacheHitRate * 100
		b.ReportMetric(hitRate, "cache_hit_%")
	}
}

// =============================================================================
// Benchmark: High Load Scenarios
// =============================================================================
// BenchmarkBetaChainBuild_HighLoad_ManyFacts benchmarks with many facts
func BenchmarkBetaChainBuild_HighLoad_ManyFacts(b *testing.B) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Note: Facts are added to network via separate mechanism
	// This benchmark focuses on chain building performance
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 2000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	patterns := createSimilarPatterns(5)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ruleID := fmt.Sprintf("rule_%d", i%50)
		_, err := builder.BuildChain(patterns, ruleID)
		if err != nil {
			b.Fatalf("BuildChain failed: %v", err)
		}
	}
	b.StopTimer()
	reportBenchmarkMetrics(b, builder)
}

// BenchmarkBetaChainBuild_HighLoad_ManyRules benchmarks with many rules
func BenchmarkBetaChainBuild_HighLoad_ManyRules(b *testing.B) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 5000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	// Build 1000 rules upfront
	patterns := createSimilarPatterns(5)
	for i := 0; i < 1000; i++ {
		ruleID := fmt.Sprintf("prebuilt_rule_%d", i)
		_, err := builder.BuildChain(patterns, ruleID)
		if err != nil {
			b.Fatalf("Failed to prebuild rule %d: %v", i, err)
		}
	}
	// Reset metrics after prebuild
	builder.ResetMetrics()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ruleID := fmt.Sprintf("new_rule_%d", i)
		_, err := builder.BuildChain(patterns, ruleID)
		if err != nil {
			b.Fatalf("BuildChain failed: %v", err)
		}
	}
	b.StopTimer()
	reportBenchmarkMetrics(b, builder)
}

// =============================================================================
// Benchmark: Memory Usage
// =============================================================================
// BenchmarkMemory_WithSharing benchmarks memory usage with sharing
func BenchmarkMemory_WithSharing(b *testing.B) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 2000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	patterns := createSimilarPatterns(10)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ruleID := fmt.Sprintf("rule_%d", i%100)
		_, err := builder.BuildChain(patterns, ruleID)
		if err != nil {
			b.Fatalf("BuildChain failed: %v", err)
		}
	}
	b.StopTimer()
	reportBenchmarkMetrics(b, builder)
}

// BenchmarkMemory_WithoutSharing benchmarks memory usage without sharing
func BenchmarkMemory_WithoutSharing(b *testing.B) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       false,
		HashCacheSize: 0,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	patterns := createSimilarPatterns(10)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ruleID := fmt.Sprintf("rule_%d", i%100)
		_, err := builder.BuildChain(patterns, ruleID)
		if err != nil {
			b.Fatalf("BuildChain failed: %v", err)
		}
	}
	b.StopTimer()
	reportBenchmarkMetrics(b, builder)
}

// =============================================================================
// Benchmark: Prefix Sharing Performance
// =============================================================================
// BenchmarkPrefixSharing_Enabled benchmarks with prefix sharing enabled
func BenchmarkPrefixSharing_Enabled(b *testing.B) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	builder.SetPrefixSharingEnabled(true)
	// Create patterns with common prefixes
	patterns := createPatternsWithCommonPrefix(8)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ruleID := fmt.Sprintf("prefix_rule_%d", i%20)
		_, err := builder.BuildChain(patterns, ruleID)
		if err != nil {
			b.Fatalf("BuildChain failed: %v", err)
		}
	}
	b.StopTimer()
	reportBenchmarkMetrics(b, builder)
	b.ReportMetric(float64(builder.GetPrefixCacheSize()), "prefix_cache_size")
}

// BenchmarkPrefixSharing_Disabled benchmarks with prefix sharing disabled
func BenchmarkPrefixSharing_Disabled(b *testing.B) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	builder.SetPrefixSharingEnabled(false)
	patterns := createPatternsWithCommonPrefix(8)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ruleID := fmt.Sprintf("prefix_rule_%d", i%20)
		_, err := builder.BuildChain(patterns, ruleID)
		if err != nil {
			b.Fatalf("BuildChain failed: %v", err)
		}
	}
	b.StopTimer()
	reportBenchmarkMetrics(b, builder)
}

// =============================================================================
// Helper Functions
// =============================================================================
// benchmarkWithRuleCount is a helper for benchmarking with varying rule counts
func benchmarkWithRuleCount(b *testing.B, ruleCount int, patternCreator func(int) []JoinPattern, sharingEnabled bool) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       sharingEnabled,
		HashCacheSize: 2000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	patterns := patternCreator(5)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ruleID := fmt.Sprintf("rule_%d", i%ruleCount)
		_, err := builder.BuildChain(patterns, ruleID)
		if err != nil {
			b.Fatalf("BuildChain failed: %v", err)
		}
	}
	b.StopTimer()
	reportBenchmarkMetrics(b, builder)
}

// reportBenchmarkMetrics reports standard metrics for benchmarks
func reportBenchmarkMetrics(b *testing.B, builder *BetaChainBuilder) {
	metrics := builder.GetMetrics()
	totalNodes := metrics.TotalNodesCreated + metrics.TotalNodesReused
	if totalNodes > 0 {
		sharingRatio := float64(metrics.TotalNodesReused) / float64(totalNodes) * 100
		b.ReportMetric(sharingRatio, "sharing_%")
	}
	b.ReportMetric(float64(metrics.TotalNodesCreated), "nodes_created")
	b.ReportMetric(float64(metrics.TotalNodesReused), "nodes_reused")
	b.ReportMetric(float64(metrics.TotalChainsBuilt), "chains_built")
	b.ReportMetric(metrics.AverageChainLength, "avg_chain_len")
	if metrics.HashCacheHits+metrics.HashCacheMisses > 0 {
		hashHitRate := float64(metrics.HashCacheHits) / float64(metrics.HashCacheHits+metrics.HashCacheMisses) * 100
		b.ReportMetric(hashHitRate, "hash_hit_%")
	}
	if metrics.PrefixCacheHits+metrics.PrefixCacheMisses > 0 {
		prefixHitRate := float64(metrics.PrefixCacheHits) / float64(metrics.PrefixCacheHits+metrics.PrefixCacheMisses) * 100
		b.ReportMetric(prefixHitRate, "prefix_hit_%")
	}
}

// createSimilarPatterns creates patterns with similar join conditions
func createSimilarPatterns(count int) []JoinPattern {
	patterns := make([]JoinPattern, count)
	vars := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}
	for i := 0; i < count; i++ {
		leftVar := vars[i]
		rightVar := vars[i+1]
		patterns[i] = JoinPattern{
			LeftVars:  []string{leftVar},
			RightVars: []string{rightVar},
			AllVars:   []string{leftVar, rightVar},
			VarTypes: map[string]string{
				leftVar:  "Entity",
				rightVar: "Entity",
			},
			Condition: map[string]interface{}{
				"type":     "comparison",
				"operator": "==",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": leftVar,
					"field":  "id",
				},
				"right": map[string]interface{}{
					"type":   "fieldAccess",
					"object": rightVar,
					"field":  "parent_id",
				},
			},
			Selectivity:   0.3,
			EstimatedCost: float64(i) * 0.5,
		}
	}
	return patterns
}

// createMixedPatterns creates patterns with varying complexity
func createMixedPatterns(count int) []JoinPattern {
	patterns := make([]JoinPattern, count)
	vars := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}
	operators := []string{"==", "!=", ">", "<", ">=", "<="}
	fields := []string{"id", "value", "amount", "status", "type", "category"}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < count; i++ {
		leftVar := vars[i]
		rightVar := vars[i+1]
		operator := operators[rng.Intn(len(operators))]
		field := fields[rng.Intn(len(fields))]
		patterns[i] = JoinPattern{
			LeftVars:  []string{leftVar},
			RightVars: []string{rightVar},
			AllVars:   []string{leftVar, rightVar},
			VarTypes: map[string]string{
				leftVar:  fmt.Sprintf("Type%d", rand.Intn(5)),
				rightVar: fmt.Sprintf("Type%d", rand.Intn(5)),
			},
			Condition: map[string]interface{}{
				"type":     "comparison",
				"operator": operator,
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": leftVar,
					"field":  field,
				},
				"right": map[string]interface{}{
					"type":   "fieldAccess",
					"object": rightVar,
					"field":  field,
				},
			},
			Selectivity:   rand.Float64(),
			EstimatedCost: rand.Float64() * float64(count),
		}
	}
	return patterns
}

// createComplexPatterns creates complex patterns with multiple conditions
func createComplexPatterns(count int) []JoinPattern {
	patterns := make([]JoinPattern, count)
	vars := []string{"p", "o", "pay", "ship", "inv", "cust", "prod", "cat", "tag", "loc", "addr", "contact"}
	for i := 0; i < count; i++ {
		leftVar := vars[i]
		rightVar := vars[i+1]
		patterns[i] = JoinPattern{
			LeftVars:  []string{leftVar},
			RightVars: []string{rightVar},
			AllVars:   []string{leftVar, rightVar},
			VarTypes: map[string]string{
				leftVar:  fmt.Sprintf("ComplexType%d", i),
				rightVar: fmt.Sprintf("ComplexType%d", i+1),
			},
			Condition: map[string]interface{}{
				"type":     "and",
				"operator": "&&",
				"left": map[string]interface{}{
					"type":     "comparison",
					"operator": "==",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": leftVar,
						"field":  "id",
					},
					"right": map[string]interface{}{
						"type":   "fieldAccess",
						"object": rightVar,
						"field":  "parent_id",
					},
				},
				"right": map[string]interface{}{
					"type":     "comparison",
					"operator": ">",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": rightVar,
						"field":  "value",
					},
					"right": map[string]interface{}{
						"type":  "literal",
						"value": 100,
					},
				},
			},
			Selectivity:   0.1 + (float64(i) * 0.1),
			EstimatedCost: float64(i) * 1.5,
		}
	}
	return patterns
}

// createPatternWithSelectivity creates a pattern with specific selectivity
func createPatternWithSelectivity(leftVar, rightVar string, selectivity float64) JoinPattern {
	return JoinPattern{
		LeftVars:  []string{leftVar},
		RightVars: []string{rightVar},
		AllVars:   []string{leftVar, rightVar},
		VarTypes: map[string]string{
			leftVar:  "Entity",
			rightVar: "Entity",
		},
		Condition: map[string]interface{}{
			"type":     "comparison",
			"operator": "==",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": leftVar,
				"field":  "id",
			},
			"right": map[string]interface{}{
				"type":   "fieldAccess",
				"object": rightVar,
				"field":  "ref_id",
			},
		},
		Selectivity:   selectivity,
		EstimatedCost: selectivity * 10,
	}
}

// createPatternsWithCommonPrefix creates patterns that share common prefixes
func createPatternsWithCommonPrefix(count int) []JoinPattern {
	patterns := make([]JoinPattern, count)
	vars := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	// First few patterns are identical (common prefix)
	commonPrefixLength := count / 2
	for i := 0; i < commonPrefixLength; i++ {
		leftVar := vars[i]
		rightVar := vars[i+1]
		patterns[i] = JoinPattern{
			LeftVars:  []string{leftVar},
			RightVars: []string{rightVar},
			AllVars:   []string{leftVar, rightVar},
			VarTypes: map[string]string{
				leftVar:  "PrefixEntity",
				rightVar: "PrefixEntity",
			},
			Condition: map[string]interface{}{
				"type":     "comparison",
				"operator": "==",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": leftVar,
					"field":  "common_field",
				},
				"right": map[string]interface{}{
					"type":   "fieldAccess",
					"object": rightVar,
					"field":  "common_field",
				},
			},
			Selectivity:   0.2,
			EstimatedCost: float64(i) * 0.3,
		}
	}
	// Remaining patterns vary
	for i := commonPrefixLength; i < count; i++ {
		leftVar := vars[i]
		rightVar := vars[i+1]
		patterns[i] = JoinPattern{
			LeftVars:  []string{leftVar},
			RightVars: []string{rightVar},
			AllVars:   []string{leftVar, rightVar},
			VarTypes: map[string]string{
				leftVar:  "VariedEntity",
				rightVar: "VariedEntity",
			},
			Condition: map[string]interface{}{
				"type":     "comparison",
				"operator": "!=",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": leftVar,
					"field":  fmt.Sprintf("field_%d", i),
				},
				"right": map[string]interface{}{
					"type":   "fieldAccess",
					"object": rightVar,
					"field":  fmt.Sprintf("field_%d", i),
				},
			},
			Selectivity:   0.5 + (float64(i) * 0.05),
			EstimatedCost: float64(i) * 0.8,
		}
	}
	return patterns
}
