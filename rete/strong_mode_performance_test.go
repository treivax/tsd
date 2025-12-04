// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"strings"
	"testing"
	"time"
)

func TestNewStrongModePerformanceMetrics(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	if pm == nil {
		t.Fatal("Expected non-nil performance metrics")
	}

	if pm.TransactionCount != 0 {
		t.Errorf("Expected initial transaction count to be 0, got %d", pm.TransactionCount)
	}

	if pm.HealthScore != 100.0 {
		t.Errorf("Expected initial health score to be 100, got %.1f", pm.HealthScore)
	}

	if pm.PerformanceGrade != "A" {
		t.Errorf("Expected initial grade to be A, got %s", pm.PerformanceGrade)
	}

	if !pm.IsHealthy {
		t.Error("Expected initial state to be healthy")
	}

	if pm.CurrentConfig == nil {
		t.Error("Expected CurrentConfig to be initialized")
	}

	if pm.RollbackReasons == nil {
		t.Error("Expected RollbackReasons map to be initialized")
	}

	if pm.ConfigChangeHistory == nil {
		t.Error("Expected ConfigChangeHistory to be initialized")
	}
}

func TestRecordTransaction_Basic(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	// Create mock coherence metrics
	cm := &CoherenceMetrics{
		FactsSubmitted:          10,
		FactsPersisted:          10,
		FactsFailed:             0,
		TotalVerifyAttempts:     10,
		TotalRetries:            2,
		TotalTimeouts:           0,
		TotalWaitTime:           100 * time.Millisecond,
		MaxRetriesForSingleFact: 1,
	}

	duration := 500 * time.Millisecond
	pm.RecordTransaction(duration, 10, true, cm)

	if pm.TransactionCount != 1 {
		t.Errorf("Expected transaction count 1, got %d", pm.TransactionCount)
	}

	if pm.SuccessfulTransactions != 1 {
		t.Errorf("Expected 1 successful transaction, got %d", pm.SuccessfulTransactions)
	}

	if pm.FailedTransactions != 0 {
		t.Errorf("Expected 0 failed transactions, got %d", pm.FailedTransactions)
	}

	if pm.TotalFactsProcessed != 10 {
		t.Errorf("Expected 10 facts processed, got %d", pm.TotalFactsProcessed)
	}

	if pm.TotalFactsPersisted != 10 {
		t.Errorf("Expected 10 facts persisted, got %d", pm.TotalFactsPersisted)
	}

	if pm.TotalRetries != 2 {
		t.Errorf("Expected 2 retries, got %d", pm.TotalRetries)
	}

	if pm.AvgTransactionTime != duration {
		t.Errorf("Expected avg transaction time %v, got %v", duration, pm.AvgTransactionTime)
	}
}

func TestRecordTransaction_Failed(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	cm := &CoherenceMetrics{
		FactsSubmitted: 10,
		FactsPersisted: 5,
		FactsFailed:    5,
		WasRolledBack:  true,
		RollbackReason: "Storage failure",
	}

	pm.RecordTransaction(1*time.Second, 10, false, cm)

	if pm.FailedTransactions != 1 {
		t.Errorf("Expected 1 failed transaction, got %d", pm.FailedTransactions)
	}

	if pm.TotalRollbacks != 1 {
		t.Errorf("Expected 1 rollback, got %d", pm.TotalRollbacks)
	}

	if pm.RollbackReasons["Storage failure"] != 1 {
		t.Errorf("Expected rollback reason to be recorded")
	}

	if pm.TotalFactsFailed != 5 {
		t.Errorf("Expected 5 failed facts, got %d", pm.TotalFactsFailed)
	}
}

func TestRecordTransaction_MultipleTransactions(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	// Record 10 successful transactions
	for i := 0; i < 10; i++ {
		cm := &CoherenceMetrics{
			FactsSubmitted:      5,
			FactsPersisted:      5,
			TotalVerifyAttempts: 5,
			TotalRetries:        1,
		}
		pm.RecordTransaction(100*time.Millisecond, 5, true, cm)
	}

	if pm.TransactionCount != 10 {
		t.Errorf("Expected 10 transactions, got %d", pm.TransactionCount)
	}

	if pm.SuccessfulTransactions != 10 {
		t.Errorf("Expected 10 successful transactions, got %d", pm.SuccessfulTransactions)
	}

	if pm.TotalFactsProcessed != 50 {
		t.Errorf("Expected 50 facts processed, got %d", pm.TotalFactsProcessed)
	}

	expectedAvgFacts := 5.0
	if pm.AvgFactsPerTransaction != expectedAvgFacts {
		t.Errorf("Expected avg facts per tx %.1f, got %.1f", expectedAvgFacts, pm.AvgFactsPerTransaction)
	}

	if pm.TotalRetries != 10 {
		t.Errorf("Expected 10 total retries, got %d", pm.TotalRetries)
	}
}

func TestRecordTransaction_TransactionTimingStats(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	durations := []time.Duration{
		50 * time.Millisecond,
		100 * time.Millisecond,
		200 * time.Millisecond,
		150 * time.Millisecond,
		75 * time.Millisecond,
	}

	for _, d := range durations {
		cm := &CoherenceMetrics{
			FactsSubmitted: 5,
			FactsPersisted: 5,
		}
		pm.RecordTransaction(d, 5, true, cm)
	}

	if pm.MinTransactionTime != 50*time.Millisecond {
		t.Errorf("Expected min time 50ms, got %v", pm.MinTransactionTime)
	}

	if pm.MaxTransactionTime != 200*time.Millisecond {
		t.Errorf("Expected max time 200ms, got %v", pm.MaxTransactionTime)
	}

	expectedAvg := 115 * time.Millisecond
	if pm.AvgTransactionTime != expectedAvg {
		t.Errorf("Expected avg time %v, got %v", expectedAvg, pm.AvgTransactionTime)
	}
}

func TestRecordCommit(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	pm.RecordCommit(50*time.Millisecond, true)
	pm.RecordCommit(60*time.Millisecond, true)
	pm.RecordCommit(70*time.Millisecond, false)

	if pm.TotalCommits != 3 {
		t.Errorf("Expected 3 commits, got %d", pm.TotalCommits)
	}

	if pm.SuccessfulCommits != 2 {
		t.Errorf("Expected 2 successful commits, got %d", pm.SuccessfulCommits)
	}

	if pm.FailedCommits != 1 {
		t.Errorf("Expected 1 failed commit, got %d", pm.FailedCommits)
	}

	expectedAvg := 60 * time.Millisecond
	if pm.AvgCommitTime != expectedAvg {
		t.Errorf("Expected avg commit time %v, got %v", expectedAvg, pm.AvgCommitTime)
	}
}

func TestRecordConfigChange(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	pm.RecordConfigChange(
		"VerifyRetryDelay",
		50*time.Millisecond,
		20*time.Millisecond,
		"Optimizing for fast storage",
	)

	if len(pm.ConfigChangeHistory) != 1 {
		t.Fatalf("Expected 1 config change, got %d", len(pm.ConfigChangeHistory))
	}

	change := pm.ConfigChangeHistory[0]
	if change.Parameter != "VerifyRetryDelay" {
		t.Errorf("Expected parameter VerifyRetryDelay, got %s", change.Parameter)
	}

	if change.Reason != "Optimizing for fast storage" {
		t.Errorf("Unexpected reason: %s", change.Reason)
	}

	if change.ImpactObserved {
		t.Error("Expected ImpactObserved to be false initially")
	}
}

func TestHealthIndicators_HighFailureRate(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	// Record 10 transactions, 4 failed (40% failure rate)
	for i := 0; i < 6; i++ {
		cm := &CoherenceMetrics{FactsSubmitted: 10, FactsPersisted: 10}
		pm.RecordTransaction(100*time.Millisecond, 10, true, cm)
	}
	for i := 0; i < 4; i++ {
		cm := &CoherenceMetrics{FactsSubmitted: 10, FactsPersisted: 0, FactsFailed: 10}
		pm.RecordTransaction(100*time.Millisecond, 10, false, cm)
	}

	if pm.IsHealthy {
		t.Error("Expected system to be unhealthy with 40% failure rate")
	}

	if pm.HealthScore > 80.0 {
		t.Errorf("Expected health score < 80, got %.1f", pm.HealthScore)
	}

	if pm.PerformanceGrade == "A" {
		t.Errorf("Expected grade worse than A, got %s", pm.PerformanceGrade)
	}
}

func TestHealthIndicators_HighTimeoutRate(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	// High timeout rate
	cm := &CoherenceMetrics{
		FactsSubmitted:      100,
		FactsPersisted:      90,
		TotalTimeouts:       10, // 10% timeout rate
		TotalVerifyAttempts: 100,
	}

	pm.RecordTransaction(1*time.Second, 100, true, cm)

	if pm.TimeoutRate != 0.1 {
		t.Errorf("Expected timeout rate 0.1, got %.2f", pm.TimeoutRate)
	}

	if pm.HealthScore >= 90.0 {
		t.Errorf("Expected health score < 90 with high timeouts, got %.1f", pm.HealthScore)
	}
}

func TestHealthIndicators_HighRetryRate(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	// High retry rate (3 retries per fact)
	cm := &CoherenceMetrics{
		FactsSubmitted:      10,
		FactsPersisted:      10,
		TotalRetries:        30,
		TotalVerifyAttempts: 40,
	}

	pm.RecordTransaction(500*time.Millisecond, 10, true, cm)

	if pm.AvgRetriesPerFact != 3.0 {
		t.Errorf("Expected avg retries 3.0, got %.1f", pm.AvgRetriesPerFact)
	}

	if pm.HealthScore >= 90.0 {
		t.Errorf("Expected health score < 90 with high retries, got %.1f", pm.HealthScore)
	}
}

func TestHealthIndicators_ExcellentPerformance(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	// Perfect performance
	for i := 0; i < 100; i++ {
		cm := &CoherenceMetrics{
			FactsSubmitted:      10,
			FactsPersisted:      10,
			TotalRetries:        1,
			TotalTimeouts:       0,
			TotalVerifyAttempts: 11,
		}
		pm.RecordTransaction(50*time.Millisecond, 10, true, cm)
	}

	if !pm.IsHealthy {
		t.Error("Expected system to be healthy")
	}

	if pm.HealthScore < 95.0 {
		t.Errorf("Expected health score >= 95, got %.1f", pm.HealthScore)
	}

	if pm.PerformanceGrade != "A" {
		t.Errorf("Expected grade A, got %s", pm.PerformanceGrade)
	}
}

func TestPerformanceGrades(t *testing.T) {
	testCases := []struct {
		name          string
		failureRate   float64
		expectedGrade string
	}{
		{"Excellent", 0.01, "A"},
		{"Good", 0.08, "B"},
		{"Fair", 0.15, "C"},
		{"Poor", 0.25, "D"},
		{"Failing", 0.40, "F"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pm := NewStrongModePerformanceMetrics()

			successCount := int(100 * (1 - tc.failureRate))
			failureCount := int(100 * tc.failureRate)

			for i := 0; i < successCount; i++ {
				cm := &CoherenceMetrics{FactsSubmitted: 10, FactsPersisted: 10}
				pm.RecordTransaction(100*time.Millisecond, 10, true, cm)
			}
			for i := 0; i < failureCount; i++ {
				cm := &CoherenceMetrics{FactsSubmitted: 10, FactsFailed: 10}
				pm.RecordTransaction(100*time.Millisecond, 10, false, cm)
			}

			if pm.PerformanceGrade != tc.expectedGrade {
				t.Errorf("Expected grade %s, got %s (health score: %.1f)",
					tc.expectedGrade, pm.PerformanceGrade, pm.HealthScore)
			}
		})
	}
}

func TestRecommendations_HighTimeout(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	// Create high timeout scenario
	cm := &CoherenceMetrics{
		FactsSubmitted: 100,
		FactsPersisted: 90,
		TotalTimeouts:  10, // 10% timeout rate
	}

	pm.RecordTransaction(1*time.Second, 100, true, cm)

	hasTimeoutRecommendation := false
	for _, rec := range pm.Recommendations {
		if strings.Contains(rec, "timeout rate") {
			hasTimeoutRecommendation = true
			break
		}
	}

	if !hasTimeoutRecommendation {
		t.Error("Expected recommendation about high timeout rate")
	}
}

func TestRecommendations_HighRetries(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	cm := &CoherenceMetrics{
		FactsSubmitted: 10,
		FactsPersisted: 10,
		TotalRetries:   15, // 1.5 retries per fact
	}

	pm.RecordTransaction(500*time.Millisecond, 10, true, cm)

	hasRetryRecommendation := false
	for _, rec := range pm.Recommendations {
		if strings.Contains(rec, "retries per fact") {
			hasRetryRecommendation = true
			break
		}
	}

	if !hasRetryRecommendation {
		t.Error("Expected recommendation about high retry rate")
	}
}

func TestRecommendations_ExcellentPerformance(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	// Perfect metrics
	for i := 0; i < 100; i++ {
		cm := &CoherenceMetrics{
			FactsSubmitted:      10,
			FactsPersisted:      10,
			TotalRetries:        1,
			TotalVerifyAttempts: 11,
		}
		pm.RecordTransaction(50*time.Millisecond, 10, true, cm)
	}

	hasExcellentMessage := false
	for _, rec := range pm.Recommendations {
		if strings.Contains(rec, "Excellent performance") {
			hasExcellentMessage = true
			break
		}
	}

	if !hasExcellentMessage {
		t.Error("Expected recommendation about excellent performance")
	}
}

func TestGetReport_Format(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	// Add some data
	cm := &CoherenceMetrics{
		FactsSubmitted:      10,
		FactsPersisted:      10,
		TotalRetries:        2,
		TotalVerifyAttempts: 12,
	}

	pm.RecordTransaction(100*time.Millisecond, 10, true, cm)

	report := pm.GetReport()

	if report == "" {
		t.Error("Expected non-empty report")
	}

	// Check for key sections
	expectedSections := []string{
		"Strong Mode Performance Report",
		"Overall Health",
		"Transactions",
		"Facts",
		"Verifications",
		"Current Configuration",
		"Recommendations",
	}

	for _, section := range expectedSections {
		if !strings.Contains(report, section) {
			t.Errorf("Report missing section: %s", section)
		}
	}
}

func TestStrongModePerformance_GetSummary(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	cm := &CoherenceMetrics{
		FactsSubmitted:      10,
		FactsPersisted:      10,
		TotalRetries:        2,
		TotalVerifyAttempts: 12,
	}

	pm.RecordTransaction(100*time.Millisecond, 10, true, cm)

	summary := pm.GetSummary()

	if summary == "" {
		t.Error("Expected non-empty summary")
	}

	// Summary should contain key metrics
	expectedTerms := []string{"Strong Mode", "txns", "facts", "retries", "Health"}
	for _, term := range expectedTerms {
		if !strings.Contains(summary, term) {
			t.Errorf("Summary missing term: %s", term)
		}
	}
}

func TestClone(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	// Populate with data
	cm := &CoherenceMetrics{
		FactsSubmitted: 10,
		FactsPersisted: 10,
		TotalRetries:   2,
		WasRolledBack:  true,
		RollbackReason: "Test rollback",
	}

	pm.RecordTransaction(100*time.Millisecond, 10, true, cm)
	pm.RecordConfigChange("test", "old", "new", "testing")

	// Clone
	clone := pm.Clone()

	if clone == nil {
		t.Fatal("Expected non-nil clone")
	}

	// Verify key fields
	if clone.TransactionCount != pm.TransactionCount {
		t.Errorf("Clone transaction count mismatch: %d != %d", clone.TransactionCount, pm.TransactionCount)
	}

	if clone.HealthScore != pm.HealthScore {
		t.Errorf("Clone health score mismatch: %.1f != %.1f", clone.HealthScore, pm.HealthScore)
	}

	if len(clone.RollbackReasons) != len(pm.RollbackReasons) {
		t.Error("Clone rollback reasons length mismatch")
	}

	if len(clone.ConfigChangeHistory) != len(pm.ConfigChangeHistory) {
		t.Error("Clone config change history length mismatch")
	}

	// Verify deep copy (modifying clone shouldn't affect original)
	clone.TransactionCount = 999
	if pm.TransactionCount == 999 {
		t.Error("Clone is not independent of original")
	}
}

func TestStrongModePerformance_ConcurrentAccess(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	// Simulate concurrent access
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				cm := &CoherenceMetrics{
					FactsSubmitted: 10,
					FactsPersisted: 10,
					TotalRetries:   1,
				}
				pm.RecordTransaction(100*time.Millisecond, 10, true, cm)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should have 1000 transactions total
	if pm.TransactionCount != 1000 {
		t.Errorf("Expected 1000 transactions, got %d", pm.TransactionCount)
	}

	// Should still be able to get report without panicking
	report := pm.GetReport()
	if report == "" {
		t.Error("Expected non-empty report after concurrent access")
	}
}

func TestTopRollbackReasons(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	// Add various rollback reasons
	reasons := map[string]int{
		"Storage failure":   10,
		"Timeout":           5,
		"Network error":     3,
		"Permission denied": 1,
	}

	for reason, count := range reasons {
		for i := 0; i < count; i++ {
			cm := &CoherenceMetrics{
				WasRolledBack:  true,
				RollbackReason: reason,
			}
			pm.RecordTransaction(100*time.Millisecond, 10, false, cm)
		}
	}

	top := pm.getTopRollbackReasons(3)

	if len(top) > 3 {
		t.Errorf("Expected at most 3 reasons, got %d", len(top))
	}

	// First should be "Storage failure" (10 occurrences)
	if !strings.Contains(top[0], "Storage failure") {
		t.Errorf("Expected first reason to be Storage failure, got %s", top[0])
	}

	// Should contain count
	if !strings.Contains(top[0], "10") {
		t.Error("Expected count to be included in reason string")
	}
}

func TestMetricsAccumulation(t *testing.T) {
	pm := NewStrongModePerformanceMetrics()

	// Record multiple transactions with different characteristics
	transactions := []struct {
		duration time.Duration
		facts    int
		retries  int
		timeouts int
		success  bool
	}{
		{100 * time.Millisecond, 10, 1, 0, true},
		{200 * time.Millisecond, 20, 3, 1, true},
		{150 * time.Millisecond, 15, 2, 0, true},
		{300 * time.Millisecond, 25, 5, 2, false},
	}

	totalRetries := 0
	totalTimeouts := 0
	totalFacts := 0

	for _, tx := range transactions {
		cm := &CoherenceMetrics{
			FactsSubmitted:      tx.facts,
			FactsPersisted:      tx.facts,
			TotalRetries:        tx.retries,
			TotalTimeouts:       tx.timeouts,
			TotalVerifyAttempts: tx.facts + tx.retries,
		}

		if !tx.success {
			cm.FactsFailed = tx.facts
			cm.FactsPersisted = 0
		}

		pm.RecordTransaction(tx.duration, tx.facts, tx.success, cm)

		totalRetries += tx.retries
		totalTimeouts += tx.timeouts
		totalFacts += tx.facts
	}

	if pm.TotalRetries != totalRetries {
		t.Errorf("Expected %d total retries, got %d", totalRetries, pm.TotalRetries)
	}

	if pm.TotalTimeouts != totalTimeouts {
		t.Errorf("Expected %d total timeouts, got %d", totalTimeouts, pm.TotalTimeouts)
	}

	if pm.TotalFactsProcessed != totalFacts {
		t.Errorf("Expected %d total facts, got %d", totalFacts, pm.TotalFactsProcessed)
	}

	if pm.SuccessfulTransactions != 3 {
		t.Errorf("Expected 3 successful transactions, got %d", pm.SuccessfulTransactions)
	}

	if pm.FailedTransactions != 1 {
		t.Errorf("Expected 1 failed transaction, got %d", pm.FailedTransactions)
	}
}
