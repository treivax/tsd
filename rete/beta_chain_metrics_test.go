// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
	"time"
)

func TestNewBetaChainMetrics(t *testing.T) {
	metrics := NewBetaChainMetrics()

	if metrics == nil {
		t.Fatal("NewBetaChainMetrics returned nil")
	}

	if metrics.ChainDetails == nil {
		t.Error("ChainDetails should be initialized")
	}

	if len(metrics.ChainDetails) != 0 {
		t.Error("ChainDetails should be empty initially")
	}
}

func TestRecordChainBuild(t *testing.T) {
	metrics := NewBetaChainMetrics()

	detail := BetaChainMetricDetail{
		RuleID:       "rule1",
		ChainLength:  5,
		NodesCreated: 3,
		NodesReused:  2,
		BuildTime:    100 * time.Millisecond,
		Timestamp:    time.Now(),
	}

	metrics.RecordChainBuild(detail)

	if metrics.TotalChainsBuilt != 1 {
		t.Errorf("Expected TotalChainsBuilt=1, got %d", metrics.TotalChainsBuilt)
	}

	if metrics.TotalNodesCreated != 3 {
		t.Errorf("Expected TotalNodesCreated=3, got %d", metrics.TotalNodesCreated)
	}

	if metrics.TotalNodesReused != 2 {
		t.Errorf("Expected TotalNodesReused=2, got %d", metrics.TotalNodesReused)
	}

	expectedAvgLength := 5.0
	if metrics.AverageChainLength != expectedAvgLength {
		t.Errorf("Expected AverageChainLength=%.2f, got %.2f", expectedAvgLength, metrics.AverageChainLength)
	}

	expectedSharingRatio := 2.0 / 5.0 // 2 reused out of 5 total
	if metrics.SharingRatio != expectedSharingRatio {
		t.Errorf("Expected SharingRatio=%.2f, got %.2f", expectedSharingRatio, metrics.SharingRatio)
	}
}

func TestRecordChainBuild_MultipleChains(t *testing.T) {
	metrics := NewBetaChainMetrics()

	// First chain: 3 created, 2 reused = 5 total
	metrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:       "rule1",
		ChainLength:  5,
		NodesCreated: 3,
		NodesReused:  2,
		BuildTime:    100 * time.Millisecond,
		Timestamp:    time.Now(),
	})

	// Second chain: 1 created, 4 reused = 5 total
	metrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:       "rule2",
		ChainLength:  5,
		NodesCreated: 1,
		NodesReused:  4,
		BuildTime:    50 * time.Millisecond,
		Timestamp:    time.Now(),
	})

	if metrics.TotalChainsBuilt != 2 {
		t.Errorf("Expected TotalChainsBuilt=2, got %d", metrics.TotalChainsBuilt)
	}

	if metrics.TotalNodesCreated != 4 {
		t.Errorf("Expected TotalNodesCreated=4, got %d", metrics.TotalNodesCreated)
	}

	if metrics.TotalNodesReused != 6 {
		t.Errorf("Expected TotalNodesReused=6, got %d", metrics.TotalNodesReused)
	}

	expectedAvgLength := 10.0 / 2.0 // 10 total nodes / 2 chains
	if metrics.AverageChainLength != expectedAvgLength {
		t.Errorf("Expected AverageChainLength=%.2f, got %.2f", expectedAvgLength, metrics.AverageChainLength)
	}

	expectedSharingRatio := 6.0 / 10.0 // 6 reused out of 10 total
	if metrics.SharingRatio != expectedSharingRatio {
		t.Errorf("Expected SharingRatio=%.2f, got %.2f", expectedSharingRatio, metrics.SharingRatio)
	}

	expectedAvgBuildTime := 75 * time.Millisecond // (100 + 50) / 2
	if metrics.AverageBuildTime != expectedAvgBuildTime {
		t.Errorf("Expected AverageBuildTime=%v, got %v", expectedAvgBuildTime, metrics.AverageBuildTime)
	}
}

func TestRecordJoinExecution(t *testing.T) {
	metrics := NewBetaChainMetrics()

	// Record a join: 10 left × 20 right = 200 possible, 50 actual results
	metrics.RecordJoinExecution(10, 20, 50, 10*time.Millisecond)

	if metrics.TotalJoinsExecuted != 1 {
		t.Errorf("Expected TotalJoinsExecuted=1, got %d", metrics.TotalJoinsExecuted)
	}

	expectedSelectivity := 50.0 / 200.0 // 0.25
	if metrics.AverageJoinSelectivity != expectedSelectivity {
		t.Errorf("Expected AverageJoinSelectivity=%.4f, got %.4f", expectedSelectivity, metrics.AverageJoinSelectivity)
	}

	if metrics.AverageResultSize != 50.0 {
		t.Errorf("Expected AverageResultSize=50.0, got %.2f", metrics.AverageResultSize)
	}

	if metrics.AverageJoinTime != 10*time.Millisecond {
		t.Errorf("Expected AverageJoinTime=10ms, got %v", metrics.AverageJoinTime)
	}
}

func TestRecordJoinExecution_Multiple(t *testing.T) {
	metrics := NewBetaChainMetrics()

	// First join: selectivity = 50/200 = 0.25, result size = 50
	metrics.RecordJoinExecution(10, 20, 50, 10*time.Millisecond)

	// Second join: selectivity = 100/200 = 0.5, result size = 100
	metrics.RecordJoinExecution(10, 20, 100, 20*time.Millisecond)

	if metrics.TotalJoinsExecuted != 2 {
		t.Errorf("Expected TotalJoinsExecuted=2, got %d", metrics.TotalJoinsExecuted)
	}

	// Average selectivity = (0.25 + 0.5) / 2 = 0.375
	expectedSelectivity := 0.375
	if metrics.AverageJoinSelectivity != expectedSelectivity {
		t.Errorf("Expected AverageJoinSelectivity=%.4f, got %.4f", expectedSelectivity, metrics.AverageJoinSelectivity)
	}

	// Average result size = (50 + 100) / 2 = 75
	if metrics.AverageResultSize != 75.0 {
		t.Errorf("Expected AverageResultSize=75.0, got %.2f", metrics.AverageResultSize)
	}

	// Average join time = (10ms + 20ms) / 2 = 15ms
	expectedAvgTime := 15 * time.Millisecond
	if metrics.AverageJoinTime != expectedAvgTime {
		t.Errorf("Expected AverageJoinTime=%v, got %v", expectedAvgTime, metrics.AverageJoinTime)
	}
}

func TestHashCacheMetrics(t *testing.T) {
	metrics := NewBetaChainMetrics()

	metrics.RecordHashCacheHit()
	metrics.RecordHashCacheHit()
	metrics.RecordHashCacheMiss()

	if metrics.HashCacheHits != 2 {
		t.Errorf("Expected HashCacheHits=2, got %d", metrics.HashCacheHits)
	}

	if metrics.HashCacheMisses != 1 {
		t.Errorf("Expected HashCacheMisses=1, got %d", metrics.HashCacheMisses)
	}

	expectedEfficiency := 2.0 / 3.0 // 2 hits out of 3 total
	efficiency := metrics.GetHashCacheEfficiency()
	if efficiency != expectedEfficiency {
		t.Errorf("Expected hash cache efficiency=%.4f, got %.4f", expectedEfficiency, efficiency)
	}

	metrics.UpdateHashCacheSize(100)
	if metrics.HashCacheSize != 100 {
		t.Errorf("Expected HashCacheSize=100, got %d", metrics.HashCacheSize)
	}
}

func TestJoinCacheMetrics(t *testing.T) {
	metrics := NewBetaChainMetrics()

	metrics.RecordJoinCacheHit()
	metrics.RecordJoinCacheHit()
	metrics.RecordJoinCacheHit()
	metrics.RecordJoinCacheMiss()
	metrics.RecordJoinCacheEviction()
	metrics.RecordJoinCacheEviction()

	if metrics.JoinCacheHits != 3 {
		t.Errorf("Expected JoinCacheHits=3, got %d", metrics.JoinCacheHits)
	}

	if metrics.JoinCacheMisses != 1 {
		t.Errorf("Expected JoinCacheMisses=1, got %d", metrics.JoinCacheMisses)
	}

	if metrics.JoinCacheEvictions != 2 {
		t.Errorf("Expected JoinCacheEvictions=2, got %d", metrics.JoinCacheEvictions)
	}

	expectedEfficiency := 3.0 / 4.0 // 3 hits out of 4 total
	efficiency := metrics.GetJoinCacheEfficiency()
	if efficiency != expectedEfficiency {
		t.Errorf("Expected join cache efficiency=%.4f, got %.4f", expectedEfficiency, efficiency)
	}

	metrics.UpdateJoinCacheSize(500)
	if metrics.JoinCacheSize != 500 {
		t.Errorf("Expected JoinCacheSize=500, got %d", metrics.JoinCacheSize)
	}
}

func TestConnectionCacheMetrics(t *testing.T) {
	metrics := NewBetaChainMetrics()

	metrics.RecordConnectionCacheHit()
	metrics.RecordConnectionCacheHit()
	metrics.RecordConnectionCacheMiss()
	metrics.RecordConnectionCacheMiss()
	metrics.RecordConnectionCacheMiss()

	if metrics.ConnectionCacheHits != 2 {
		t.Errorf("Expected ConnectionCacheHits=2, got %d", metrics.ConnectionCacheHits)
	}

	if metrics.ConnectionCacheMisses != 3 {
		t.Errorf("Expected ConnectionCacheMisses=3, got %d", metrics.ConnectionCacheMisses)
	}

	expectedEfficiency := 2.0 / 5.0 // 2 hits out of 5 total
	efficiency := metrics.GetConnectionCacheEfficiency()
	if efficiency != expectedEfficiency {
		t.Errorf("Expected connection cache efficiency=%.4f, got %.4f", expectedEfficiency, efficiency)
	}
}

func TestPrefixCacheMetrics(t *testing.T) {
	metrics := NewBetaChainMetrics()

	metrics.RecordPrefixCacheHit()
	metrics.RecordPrefixCacheMiss()
	metrics.RecordPrefixCacheMiss()

	if metrics.PrefixCacheHits != 1 {
		t.Errorf("Expected PrefixCacheHits=1, got %d", metrics.PrefixCacheHits)
	}

	if metrics.PrefixCacheMisses != 2 {
		t.Errorf("Expected PrefixCacheMisses=2, got %d", metrics.PrefixCacheMisses)
	}

	expectedEfficiency := 1.0 / 3.0 // 1 hit out of 3 total
	efficiency := metrics.GetPrefixCacheEfficiency()
	if efficiency != expectedEfficiency {
		t.Errorf("Expected prefix cache efficiency=%.4f, got %.4f", expectedEfficiency, efficiency)
	}

	metrics.UpdatePrefixCacheSize(200)
	if metrics.PrefixCacheSize != 200 {
		t.Errorf("Expected PrefixCacheSize=200, got %d", metrics.PrefixCacheSize)
	}
}

func TestCacheEfficiency_ZeroAccess(t *testing.T) {
	metrics := NewBetaChainMetrics()

	// All cache efficiencies should be 0.0 with no accesses
	if eff := metrics.GetHashCacheEfficiency(); eff != 0.0 {
		t.Errorf("Expected hash cache efficiency=0.0 with no accesses, got %.4f", eff)
	}

	if eff := metrics.GetJoinCacheEfficiency(); eff != 0.0 {
		t.Errorf("Expected join cache efficiency=0.0 with no accesses, got %.4f", eff)
	}

	if eff := metrics.GetConnectionCacheEfficiency(); eff != 0.0 {
		t.Errorf("Expected connection cache efficiency=0.0 with no accesses, got %.4f", eff)
	}

	if eff := metrics.GetPrefixCacheEfficiency(); eff != 0.0 {
		t.Errorf("Expected prefix cache efficiency=0.0 with no accesses, got %.4f", eff)
	}
}

func TestGetSnapshot(t *testing.T) {
	metrics := NewBetaChainMetrics()

	// Populate with data
	metrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:       "rule1",
		ChainLength:  5,
		NodesCreated: 3,
		NodesReused:  2,
		BuildTime:    100 * time.Millisecond,
	})

	metrics.RecordJoinExecution(10, 20, 50, 10*time.Millisecond)
	metrics.RecordHashCacheHit()
	metrics.RecordHashCacheMiss()

	// Get snapshot
	snapshot := metrics.GetSnapshot()

	// Verify snapshot data
	if snapshot.TotalChainsBuilt != 1 {
		t.Errorf("Snapshot: Expected TotalChainsBuilt=1, got %d", snapshot.TotalChainsBuilt)
	}

	if snapshot.TotalNodesCreated != 3 {
		t.Errorf("Snapshot: Expected TotalNodesCreated=3, got %d", snapshot.TotalNodesCreated)
	}

	if snapshot.HashCacheHits != 1 {
		t.Errorf("Snapshot: Expected HashCacheHits=1, got %d", snapshot.HashCacheHits)
	}

	// Modify original metrics
	metrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:       "rule2",
		ChainLength:  3,
		NodesCreated: 1,
		NodesReused:  2,
		BuildTime:    50 * time.Millisecond,
	})

	// Verify snapshot is unchanged
	if snapshot.TotalChainsBuilt != 1 {
		t.Error("Snapshot should not be affected by changes to original metrics")
	}

	if metrics.TotalChainsBuilt != 2 {
		t.Error("Original metrics should be updated")
	}
}

func TestReset(t *testing.T) {
	metrics := NewBetaChainMetrics()

	// Populate with data
	metrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:       "rule1",
		ChainLength:  5,
		NodesCreated: 3,
		NodesReused:  2,
		BuildTime:    100 * time.Millisecond,
	})

	metrics.RecordJoinExecution(10, 20, 50, 10*time.Millisecond)
	metrics.RecordHashCacheHit()
	metrics.RecordJoinCacheHit()
	metrics.UpdateJoinCacheSize(100)

	// Reset
	metrics.Reset()

	// Verify all fields are reset
	if metrics.TotalChainsBuilt != 0 {
		t.Error("TotalChainsBuilt should be reset to 0")
	}

	if metrics.TotalNodesCreated != 0 {
		t.Error("TotalNodesCreated should be reset to 0")
	}

	if metrics.TotalJoinsExecuted != 0 {
		t.Error("TotalJoinsExecuted should be reset to 0")
	}

	if metrics.HashCacheHits != 0 {
		t.Error("HashCacheHits should be reset to 0")
	}

	if metrics.JoinCacheSize != 0 {
		t.Error("JoinCacheSize should be reset to 0")
	}

	if len(metrics.ChainDetails) != 0 {
		t.Error("ChainDetails should be empty after reset")
	}
}

func TestGetSummary(t *testing.T) {
	metrics := NewBetaChainMetrics()

	metrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:        "rule1",
		ChainLength:   5,
		NodesCreated:  3,
		NodesReused:   2,
		BuildTime:     100 * time.Millisecond,
		JoinsExecuted: 2,
		TotalJoinTime: 20 * time.Millisecond,
	})

	metrics.RecordHashCacheHit()
	metrics.RecordHashCacheMiss()
	metrics.RecordJoinCacheHit()
	metrics.RecordJoinCacheMiss()
	metrics.RecordJoinCacheMiss()

	summary := metrics.GetSummary()

	if summary == nil {
		t.Fatal("GetSummary returned nil")
	}

	// Check chains section
	chains, ok := summary["chains"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected chains section in summary")
	}

	if chains["total_built"].(int) != 1 {
		t.Error("Summary chains total_built incorrect")
	}

	// Check joins section
	joins, ok := summary["joins"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected joins section in summary")
	}

	if joins["total_executed"].(int) != 2 {
		t.Error("Summary joins total_executed incorrect")
	}

	// Check cache sections
	hashCache, ok := summary["hash_cache"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected hash_cache section in summary")
	}

	if hashCache["hits"].(int) != 1 {
		t.Error("Summary hash_cache hits incorrect")
	}

	joinCache, ok := summary["join_cache"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected join_cache section in summary")
	}

	if joinCache["misses"].(int) != 2 {
		t.Error("Summary join_cache misses incorrect")
	}
}

func TestGetTopChainsByBuildTime(t *testing.T) {
	metrics := NewBetaChainMetrics()

	// Add chains with different build times
	metrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:    "rule1",
		BuildTime: 100 * time.Millisecond,
	})

	metrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:    "rule2",
		BuildTime: 300 * time.Millisecond,
	})

	metrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:    "rule3",
		BuildTime: 200 * time.Millisecond,
	})

	// Get top 2
	top := metrics.GetTopChainsByBuildTime(2)

	if len(top) != 2 {
		t.Fatalf("Expected 2 chains, got %d", len(top))
	}

	// Should be sorted by build time descending
	if top[0].RuleID != "rule2" {
		t.Errorf("Expected rule2 first, got %s", top[0].RuleID)
	}

	if top[1].RuleID != "rule3" {
		t.Errorf("Expected rule3 second, got %s", top[1].RuleID)
	}
}

func TestGetTopChainsByLength(t *testing.T) {
	metrics := NewBetaChainMetrics()

	// Add chains with different lengths
	metrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:      "rule1",
		ChainLength: 5,
	})

	metrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:      "rule2",
		ChainLength: 10,
	})

	metrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:      "rule3",
		ChainLength: 7,
	})

	// Get top 2
	top := metrics.GetTopChainsByLength(2)

	if len(top) != 2 {
		t.Fatalf("Expected 2 chains, got %d", len(top))
	}

	// Should be sorted by length descending
	if top[0].RuleID != "rule2" || top[0].ChainLength != 10 {
		t.Errorf("Expected rule2 with length 10 first, got %s with length %d", top[0].RuleID, top[0].ChainLength)
	}

	if top[1].RuleID != "rule3" || top[1].ChainLength != 7 {
		t.Errorf("Expected rule3 with length 7 second, got %s with length %d", top[1].RuleID, top[1].ChainLength)
	}
}

func TestGetTopChainsByJoinTime(t *testing.T) {
	metrics := NewBetaChainMetrics()

	// Add chains with different join times
	metrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:        "rule1",
		TotalJoinTime: 50 * time.Millisecond,
	})

	metrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:        "rule2",
		TotalJoinTime: 150 * time.Millisecond,
	})

	metrics.RecordChainBuild(BetaChainMetricDetail{
		RuleID:        "rule3",
		TotalJoinTime: 100 * time.Millisecond,
	})

	// Get top 2
	top := metrics.GetTopChainsByJoinTime(2)

	if len(top) != 2 {
		t.Fatalf("Expected 2 chains, got %d", len(top))
	}

	// Should be sorted by join time descending
	if top[0].RuleID != "rule2" {
		t.Errorf("Expected rule2 first, got %s", top[0].RuleID)
	}

	if top[1].RuleID != "rule3" {
		t.Errorf("Expected rule3 second, got %s", top[1].RuleID)
	}
}

func TestGetTopChains_EmptyMetrics(t *testing.T) {
	metrics := NewBetaChainMetrics()

	topByTime := metrics.GetTopChainsByBuildTime(5)
	if len(topByTime) != 0 {
		t.Error("Expected empty result for empty metrics (by time)")
	}

	topByLength := metrics.GetTopChainsByLength(5)
	if len(topByLength) != 0 {
		t.Error("Expected empty result for empty metrics (by length)")
	}

	topByJoinTime := metrics.GetTopChainsByJoinTime(5)
	if len(topByJoinTime) != 0 {
		t.Error("Expected empty result for empty metrics (by join time)")
	}
}

func TestGetTopChains_RequestMoreThanAvailable(t *testing.T) {
	metrics := NewBetaChainMetrics()

	// Add 2 chains
	metrics.RecordChainBuild(BetaChainMetricDetail{RuleID: "rule1"})
	metrics.RecordChainBuild(BetaChainMetricDetail{RuleID: "rule2"})

	// Request 5, should get 2
	top := metrics.GetTopChainsByBuildTime(5)
	if len(top) != 2 {
		t.Errorf("Expected 2 chains (all available), got %d", len(top))
	}
}

func TestGetJoinPerformanceStats(t *testing.T) {
	metrics := NewBetaChainMetrics()

	// Record some joins
	metrics.RecordJoinExecution(10, 20, 50, 10*time.Millisecond)
	metrics.RecordJoinExecution(5, 10, 25, 5*time.Millisecond)

	stats := metrics.GetJoinPerformanceStats()

	if stats == nil {
		t.Fatal("GetJoinPerformanceStats returned nil")
	}

	if stats["total_joins"].(int) != 2 {
		t.Error("Expected total_joins=2")
	}

	// Check that computed fields exist
	if _, ok := stats["average_time_per_join_ms"]; !ok {
		t.Error("Expected average_time_per_join_ms in stats")
	}

	if _, ok := stats["throughput_joins_per_sec"]; !ok {
		t.Error("Expected throughput_joins_per_sec in stats")
	}
}

func TestGetCacheStats(t *testing.T) {
	metrics := NewBetaChainMetrics()

	metrics.RecordHashCacheHit()
	metrics.RecordHashCacheMiss()
	metrics.RecordJoinCacheHit()
	metrics.RecordConnectionCacheHit()
	metrics.RecordPrefixCacheHit()

	stats := metrics.GetCacheStats()

	if stats == nil {
		t.Fatal("GetCacheStats returned nil")
	}

	// Check all cache types are present
	cacheTypes := []string{"hash_cache", "join_cache", "connection_cache", "prefix_cache"}
	for _, cacheType := range cacheTypes {
		if _, ok := stats[cacheType]; !ok {
			t.Errorf("Expected %s in cache stats", cacheType)
		}
	}

	// Check hash cache details
	hashCache := stats["hash_cache"].(map[string]interface{})
	if hashCache["hits"].(int) != 1 {
		t.Error("Expected hash_cache hits=1")
	}
}

func TestAddHashComputeTime(t *testing.T) {
	metrics := NewBetaChainMetrics()

	metrics.AddHashComputeTime(100 * time.Millisecond)
	metrics.AddHashComputeTime(200 * time.Millisecond)

	expected := 300 * time.Millisecond
	if metrics.TotalHashComputeTime != expected {
		t.Errorf("Expected TotalHashComputeTime=%v, got %v", expected, metrics.TotalHashComputeTime)
	}
}

func TestThreadSafety(t *testing.T) {
	metrics := NewBetaChainMetrics()

	// Run concurrent operations
	done := make(chan bool, 4)

	// Goroutine 1: Record chain builds
	go func() {
		for i := 0; i < 100; i++ {
			metrics.RecordChainBuild(BetaChainMetricDetail{
				RuleID:       "rule",
				ChainLength:  5,
				NodesCreated: 3,
				NodesReused:  2,
			})
		}
		done <- true
	}()

	// Goroutine 2: Record joins
	go func() {
		for i := 0; i < 100; i++ {
			metrics.RecordJoinExecution(10, 20, 50, time.Millisecond)
		}
		done <- true
	}()

	// Goroutine 3: Record cache hits
	go func() {
		for i := 0; i < 100; i++ {
			metrics.RecordHashCacheHit()
			metrics.RecordJoinCacheHit()
		}
		done <- true
	}()

	// Goroutine 4: Get snapshots
	go func() {
		for i := 0; i < 100; i++ {
			_ = metrics.GetSnapshot()
			_ = metrics.GetSummary()
		}
		done <- true
	}()

	// Wait for all goroutines
	for i := 0; i < 4; i++ {
		<-done
	}

	// Verify final counts
	if metrics.TotalChainsBuilt != 100 {
		t.Errorf("Expected 100 chains built, got %d", metrics.TotalChainsBuilt)
	}

	if metrics.TotalJoinsExecuted != 100 {
		t.Errorf("Expected 100 joins executed, got %d", metrics.TotalJoinsExecuted)
	}

	if metrics.HashCacheHits != 100 {
		t.Errorf("Expected 100 hash cache hits, got %d", metrics.HashCacheHits)
	}
}
