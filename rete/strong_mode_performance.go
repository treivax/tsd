// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sync"
	"time"
)

// StrongModePerformanceMetrics tracks real-world performance of Strong mode operations
type StrongModePerformanceMetrics struct {
	// Transaction-level metrics
	TransactionCount       int           `json:"transaction_count"`
	SuccessfulTransactions int           `json:"successful_transactions"`
	FailedTransactions     int           `json:"failed_transactions"`
	TotalTransactionTime   time.Duration `json:"total_transaction_time"`
	AvgTransactionTime     time.Duration `json:"avg_transaction_time"`
	MinTransactionTime     time.Duration `json:"min_transaction_time"`
	MaxTransactionTime     time.Duration `json:"max_transaction_time"`

	// Fact-level metrics
	TotalFactsProcessed    int           `json:"total_facts_processed"`
	TotalFactsPersisted    int           `json:"total_facts_persisted"`
	TotalFactsFailed       int           `json:"total_facts_failed"`
	AvgFactsPerTransaction float64       `json:"avg_facts_per_transaction"`
	AvgTimePerFact         time.Duration `json:"avg_time_per_fact"`

	// Verification metrics
	TotalVerifications    int           `json:"total_verifications"`
	SuccessfulVerifies    int           `json:"successful_verifies"`
	FailedVerifies        int           `json:"failed_verifies"`
	TotalRetries          int           `json:"total_retries"`
	AvgRetriesPerFact     float64       `json:"avg_retries_per_fact"`
	MaxRetriesPerFact     int           `json:"max_retries_per_fact"`
	TotalVerificationTime time.Duration `json:"total_verification_time"`
	AvgVerificationTime   time.Duration `json:"avg_verification_time"`

	// Timeout metrics
	TotalTimeouts    int           `json:"total_timeouts"`
	TimeoutRate      float64       `json:"timeout_rate"`
	AvgTimeToTimeout time.Duration `json:"avg_time_to_timeout"`

	// Commit metrics
	TotalCommits      int           `json:"total_commits"`
	SuccessfulCommits int           `json:"successful_commits"`
	FailedCommits     int           `json:"failed_commits"`
	AvgCommitTime     time.Duration `json:"avg_commit_time"`

	// Rollback metrics
	TotalRollbacks  int            `json:"total_rollbacks"`
	RollbackReasons map[string]int `json:"rollback_reasons"`

	// Configuration effectiveness
	CurrentConfig       *TransactionOptions `json:"current_config"`
	ConfigChangeHistory []ConfigChange      `json:"config_change_history"`

	// Health indicators
	IsHealthy        bool     `json:"is_healthy"`
	HealthScore      float64  `json:"health_score"`      // 0-100
	PerformanceGrade string   `json:"performance_grade"` // A, B, C, D, F
	Recommendations  []string `json:"recommendations"`

	// Time tracking
	StartTime   time.Time `json:"start_time"`
	LastUpdated time.Time `json:"last_updated"`

	mutex sync.RWMutex
}

// ConfigChange records a configuration parameter change
type ConfigChange struct {
	Timestamp         time.Time   `json:"timestamp"`
	Parameter         string      `json:"parameter"`
	OldValue          interface{} `json:"old_value"`
	NewValue          interface{} `json:"new_value"`
	Reason            string      `json:"reason"`
	ImpactObserved    bool        `json:"impact_observed"`
	ImpactDescription string      `json:"impact_description,omitempty"`
}

// NewStrongModePerformanceMetrics creates a new performance metrics collector
func NewStrongModePerformanceMetrics() *StrongModePerformanceMetrics {
	return &StrongModePerformanceMetrics{
		RollbackReasons:     make(map[string]int),
		ConfigChangeHistory: make([]ConfigChange, 0),
		CurrentConfig:       DefaultTransactionOptions(),
		StartTime:           time.Now(),
		LastUpdated:         time.Now(),
		MinTransactionTime:  time.Duration(0),
		IsHealthy:           true,
		HealthScore:         100.0,
		PerformanceGrade:    "A",
	}
}

// RecordTransaction records metrics from a completed transaction
func (pm *StrongModePerformanceMetrics) RecordTransaction(
	duration time.Duration,
	factCount int,
	success bool,
	coherenceMetrics *CoherenceMetrics,
) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.TransactionCount++
	pm.TotalTransactionTime += duration
	pm.LastUpdated = time.Now()

	if success {
		pm.SuccessfulTransactions++
	} else {
		pm.FailedTransactions++
	}

	// Update transaction time stats
	if pm.MinTransactionTime == 0 || duration < pm.MinTransactionTime {
		pm.MinTransactionTime = duration
	}
	if duration > pm.MaxTransactionTime {
		pm.MaxTransactionTime = duration
	}
	pm.AvgTransactionTime = pm.TotalTransactionTime / time.Duration(pm.TransactionCount)

	// Update fact-level metrics
	pm.TotalFactsProcessed += factCount
	if coherenceMetrics != nil {
		pm.TotalFactsPersisted += coherenceMetrics.FactsPersisted
		pm.TotalFactsFailed += coherenceMetrics.FactsFailed
		pm.TotalVerifications += coherenceMetrics.TotalVerifyAttempts
		pm.TotalRetries += coherenceMetrics.TotalRetries
		pm.TotalTimeouts += coherenceMetrics.TotalTimeouts
		pm.TotalVerificationTime += coherenceMetrics.TotalWaitTime

		if coherenceMetrics.MaxRetriesForSingleFact > pm.MaxRetriesPerFact {
			pm.MaxRetriesPerFact = coherenceMetrics.MaxRetriesForSingleFact
		}

		if coherenceMetrics.WasRolledBack {
			pm.TotalRollbacks++
			if coherenceMetrics.RollbackReason != "" {
				pm.RollbackReasons[coherenceMetrics.RollbackReason]++
			}
		}
	}

	// Recalculate averages
	if pm.TransactionCount > 0 {
		pm.AvgFactsPerTransaction = float64(pm.TotalFactsProcessed) / float64(pm.TransactionCount)
	}
	if pm.TotalFactsPersisted > 0 {
		pm.AvgRetriesPerFact = float64(pm.TotalRetries) / float64(pm.TotalFactsPersisted)
		pm.AvgTimePerFact = pm.TotalTransactionTime / time.Duration(pm.TotalFactsPersisted)
	}
	if pm.TotalVerifications > 0 {
		pm.AvgVerificationTime = pm.TotalVerificationTime / time.Duration(pm.TotalVerifications)
		pm.SuccessfulVerifies = pm.TotalFactsPersisted
		pm.FailedVerifies = pm.TotalVerifications - pm.SuccessfulVerifies
	}
	if pm.TotalFactsProcessed > 0 {
		pm.TimeoutRate = float64(pm.TotalTimeouts) / float64(pm.TotalFactsProcessed)
	}

	// Update health indicators
	pm.updateHealthIndicators()
}

// RecordCommit records a commit operation
func (pm *StrongModePerformanceMetrics) RecordCommit(duration time.Duration, success bool) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.TotalCommits++
	if success {
		pm.SuccessfulCommits++
	} else {
		pm.FailedCommits++
	}

	// Update average commit time
	if pm.TotalCommits > 0 {
		pm.AvgCommitTime = (pm.AvgCommitTime*time.Duration(pm.TotalCommits-1) + duration) / time.Duration(pm.TotalCommits)
	}
}

// RecordConfigChange records a configuration parameter change
func (pm *StrongModePerformanceMetrics) RecordConfigChange(
	parameter string,
	oldValue, newValue interface{},
	reason string,
) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	change := ConfigChange{
		Timestamp:      time.Now(),
		Parameter:      parameter,
		OldValue:       oldValue,
		NewValue:       newValue,
		Reason:         reason,
		ImpactObserved: false,
	}

	pm.ConfigChangeHistory = append(pm.ConfigChangeHistory, change)
}

// updateHealthIndicators calculates health score and recommendations
func (pm *StrongModePerformanceMetrics) updateHealthIndicators() {
	// Calculate health score (0-100)
	score := 100.0

	// Deduct points for high failure rate
	if pm.TransactionCount > 0 {
		failureRate := float64(pm.FailedTransactions) / float64(pm.TransactionCount)
		if failureRate > 0.05 {
			score -= (failureRate - 0.05) * 100 // Deduct up to 10 points
		}
	}

	// Deduct points for high timeout rate
	if pm.TimeoutRate > 0.01 {
		score -= (pm.TimeoutRate - 0.01) * 500 // Deduct up to 20 points
	}

	// Deduct points for high retry rate
	if pm.AvgRetriesPerFact > 0.5 {
		score -= (pm.AvgRetriesPerFact - 0.5) * 20 // Deduct up to 20 points
	}

	// Deduct points for many rollbacks
	if pm.TransactionCount > 0 {
		rollbackRate := float64(pm.TotalRollbacks) / float64(pm.TransactionCount)
		if rollbackRate > 0.02 {
			score -= (rollbackRate - 0.02) * 300 // Deduct up to 15 points
		}
	}

	// Ensure score is in valid range
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	pm.HealthScore = score
	pm.IsHealthy = score >= 80.0

	// Assign performance grade
	switch {
	case score >= 90:
		pm.PerformanceGrade = "A"
	case score >= 80:
		pm.PerformanceGrade = "B"
	case score >= 70:
		pm.PerformanceGrade = "C"
	case score >= 60:
		pm.PerformanceGrade = "D"
	default:
		pm.PerformanceGrade = "F"
	}

	// Generate recommendations
	pm.Recommendations = pm.generateRecommendations()
}

// generateRecommendations generates tuning recommendations based on metrics
func (pm *StrongModePerformanceMetrics) generateRecommendations() []string {
	recommendations := make([]string, 0)

	// High timeout rate
	if pm.TimeoutRate > 0.05 {
		recommendations = append(recommendations,
			fmt.Sprintf("⚠️  High timeout rate (%.2f%%). Consider increasing SubmissionTimeout (current: %v)",
				pm.TimeoutRate*100, pm.CurrentConfig.SubmissionTimeout))
	}

	// High retry rate
	if pm.AvgRetriesPerFact > 1.0 {
		recommendations = append(recommendations,
			fmt.Sprintf("⚠️  High average retries per fact (%.2f). Consider increasing VerifyRetryDelay (current: %v) or MaxVerifyRetries (current: %d)",
				pm.AvgRetriesPerFact, pm.CurrentConfig.VerifyRetryDelay, pm.CurrentConfig.MaxVerifyRetries))
	}

	// Low retry rate - could optimize
	if pm.AvgRetriesPerFact < 0.1 && pm.CurrentConfig.MaxVerifyRetries > 5 {
		recommendations = append(recommendations,
			fmt.Sprintf("✅ Low retry rate (%.2f). You could reduce MaxVerifyRetries (current: %d) to improve performance",
				pm.AvgRetriesPerFact, pm.CurrentConfig.MaxVerifyRetries))
	}

	// Fast verification - could optimize delay
	if pm.AvgVerificationTime > 0 && pm.AvgVerificationTime < 10*time.Millisecond && pm.CurrentConfig.VerifyRetryDelay > 20*time.Millisecond {
		recommendations = append(recommendations,
			fmt.Sprintf("✅ Fast verification (avg: %v). You could reduce VerifyRetryDelay (current: %v) to improve performance",
				pm.AvgVerificationTime, pm.CurrentConfig.VerifyRetryDelay))
	}

	// Slow verification
	if pm.AvgVerificationTime > 100*time.Millisecond {
		recommendations = append(recommendations,
			fmt.Sprintf("⚠️  Slow verification (avg: %v). Consider investigating storage performance or increasing VerifyRetryDelay",
				pm.AvgVerificationTime))
	}

	// High rollback rate
	if pm.TransactionCount > 0 {
		rollbackRate := float64(pm.TotalRollbacks) / float64(pm.TransactionCount)
		if rollbackRate > 0.05 {
			recommendations = append(recommendations,
				fmt.Sprintf("⚠️  High rollback rate (%.2f%%). Top reasons: %v",
					rollbackRate*100, pm.getTopRollbackReasons(3)))
		}
	}

	// Good performance
	if pm.HealthScore >= 95 && len(recommendations) == 0 {
		recommendations = append(recommendations,
			"✅ Excellent performance! Current configuration is well-tuned for your workload.")
	}

	return recommendations
}

// getTopRollbackReasons returns the top N rollback reasons
func (pm *StrongModePerformanceMetrics) getTopRollbackReasons(n int) []string {
	type reasonCount struct {
		reason string
		count  int
	}

	counts := make([]reasonCount, 0, len(pm.RollbackReasons))
	for reason, count := range pm.RollbackReasons {
		counts = append(counts, reasonCount{reason, count})
	}

	// Simple bubble sort for top N
	for i := 0; i < len(counts) && i < n; i++ {
		for j := i + 1; j < len(counts); j++ {
			if counts[j].count > counts[i].count {
				counts[i], counts[j] = counts[j], counts[i]
			}
		}
	}

	result := make([]string, 0, n)
	for i := 0; i < len(counts) && i < n; i++ {
		result = append(result, fmt.Sprintf("%s (%d)", counts[i].reason, counts[i].count))
	}

	return result
}

// GetReport generates a comprehensive performance report
func (pm *StrongModePerformanceMetrics) GetReport() string {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	uptime := time.Since(pm.StartTime)

	return fmt.Sprintf(`
╔════════════════════════════════════════════════════════════════╗
║          Strong Mode Performance Report                        ║
╠════════════════════════════════════════════════════════════════╣
║ 📊 Overall Health                                              ║
║   Score:                  %.1f/100 (%s)                     ║
║   Status:                 %s                                   ║
║   Uptime:                 %v                                   ║
║   Last Updated:           %v ago                               ║
╠════════════════════════════════════════════════════════════════╣
║ 🔄 Transactions                                                ║
║   Total:                  %d                                   ║
║   Successful:             %d (%.1f%%)                          ║
║   Failed:                 %d (%.1f%%)                          ║
║   Avg Duration:           %v                                   ║
║   Min/Max Duration:       %v / %v                              ║
╠════════════════════════════════════════════════════════════════╣
║ 📦 Facts                                                       ║
║   Total Processed:        %d                                   ║
║   Persisted:              %d (%.1f%%)                          ║
║   Failed:                 %d (%.1f%%)                          ║
║   Avg per Transaction:    %.2f                                 ║
║   Avg Time per Fact:      %v                                   ║
╠════════════════════════════════════════════════════════════════╣
║ ✅ Verifications                                               ║
║   Total Attempts:         %d                                   ║
║   Successful:             %d (%.1f%%)                          ║
║   Total Retries:          %d                                   ║
║   Avg Retries per Fact:   %.2f                                 ║
║   Max Retries (1 fact):   %d                                   ║
║   Avg Verification Time:  %v                                   ║
╠════════════════════════════════════════════════════════════════╣
║ ⏱️  Timeouts                                                   ║
║   Total:                  %d                                   ║
║   Rate:                   %.2f%%                               ║
╠════════════════════════════════════════════════════════════════╣
║ 💾 Commits                                                     ║
║   Total:                  %d                                   ║
║   Successful:             %d (%.1f%%)                          ║
║   Avg Commit Time:        %v                                   ║
╠════════════════════════════════════════════════════════════════╣
║ 🔙 Rollbacks                                                   ║
║   Total:                  %d                                   ║
║   Top Reasons:            %v                                   ║
╠════════════════════════════════════════════════════════════════╣
║ ⚙️  Current Configuration                                      ║
║   SubmissionTimeout:      %v                                   ║
║   VerifyRetryDelay:       %v                                   ║
║   MaxVerifyRetries:       %d                                   ║
║   VerifyOnCommit:         %v                                   ║
║   Config Changes:         %d                                   ║
╠════════════════════════════════════════════════════════════════╣
║ 💡 Recommendations                                             ║
%s
╚════════════════════════════════════════════════════════════════╝
`,
		pm.HealthScore, pm.PerformanceGrade,
		pm.getHealthStatus(),
		uptime,
		time.Since(pm.LastUpdated),
		pm.TransactionCount,
		pm.SuccessfulTransactions, pm.getSuccessRate(),
		pm.FailedTransactions, pm.getFailureRate(),
		pm.AvgTransactionTime,
		pm.MinTransactionTime, pm.MaxTransactionTime,
		pm.TotalFactsProcessed,
		pm.TotalFactsPersisted, pm.getFactPersistRate(),
		pm.TotalFactsFailed, pm.getFactFailureRate(),
		pm.AvgFactsPerTransaction,
		pm.AvgTimePerFact,
		pm.TotalVerifications,
		pm.SuccessfulVerifies, pm.getVerifySuccessRate(),
		pm.TotalRetries,
		pm.AvgRetriesPerFact,
		pm.MaxRetriesPerFact,
		pm.AvgVerificationTime,
		pm.TotalTimeouts,
		pm.TimeoutRate*100,
		pm.TotalCommits,
		pm.SuccessfulCommits, pm.getCommitSuccessRate(),
		pm.AvgCommitTime,
		pm.TotalRollbacks,
		pm.getTopRollbackReasons(3),
		pm.CurrentConfig.SubmissionTimeout,
		pm.CurrentConfig.VerifyRetryDelay,
		pm.CurrentConfig.MaxVerifyRetries,
		pm.CurrentConfig.VerifyOnCommit,
		len(pm.ConfigChangeHistory),
		pm.formatRecommendations(),
	)
}

// Helper functions for percentage calculations
func (pm *StrongModePerformanceMetrics) getSuccessRate() float64 {
	if pm.TransactionCount == 0 {
		return 0
	}
	return float64(pm.SuccessfulTransactions) / float64(pm.TransactionCount) * 100
}

func (pm *StrongModePerformanceMetrics) getFailureRate() float64 {
	if pm.TransactionCount == 0 {
		return 0
	}
	return float64(pm.FailedTransactions) / float64(pm.TransactionCount) * 100
}

func (pm *StrongModePerformanceMetrics) getFactPersistRate() float64 {
	if pm.TotalFactsProcessed == 0 {
		return 0
	}
	return float64(pm.TotalFactsPersisted) / float64(pm.TotalFactsProcessed) * 100
}

func (pm *StrongModePerformanceMetrics) getFactFailureRate() float64 {
	if pm.TotalFactsProcessed == 0 {
		return 0
	}
	return float64(pm.TotalFactsFailed) / float64(pm.TotalFactsProcessed) * 100
}

func (pm *StrongModePerformanceMetrics) getVerifySuccessRate() float64 {
	if pm.TotalVerifications == 0 {
		return 0
	}
	return float64(pm.SuccessfulVerifies) / float64(pm.TotalVerifications) * 100
}

func (pm *StrongModePerformanceMetrics) getCommitSuccessRate() float64 {
	if pm.TotalCommits == 0 {
		return 0
	}
	return float64(pm.SuccessfulCommits) / float64(pm.TotalCommits) * 100
}

func (pm *StrongModePerformanceMetrics) getHealthStatus() string {
	if pm.IsHealthy {
		return "✅ Healthy"
	}
	return "⚠️  Needs Attention"
}

func (pm *StrongModePerformanceMetrics) formatRecommendations() string {
	if len(pm.Recommendations) == 0 {
		return "║   None - system performing optimally                          ║\n"
	}

	result := ""
	for _, rec := range pm.Recommendations {
		// Wrap long recommendations
		if len(rec) > 60 {
			rec = rec[:57] + "..."
		}
		result += fmt.Sprintf("║   %-60s ║\n", rec)
	}
	return result
}

// GetSummary returns a short summary suitable for logging
func (pm *StrongModePerformanceMetrics) GetSummary() string {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	return fmt.Sprintf(
		"Strong Mode: %d txns (%.1f%% success) | %d facts (%.2f retries/fact) | Health: %.0f%% (%s)",
		pm.TransactionCount,
		pm.getSuccessRate(),
		pm.TotalFactsPersisted,
		pm.AvgRetriesPerFact,
		pm.HealthScore,
		pm.PerformanceGrade,
	)
}

// Clone creates a deep copy for safe reading without locks
func (pm *StrongModePerformanceMetrics) Clone() *StrongModePerformanceMetrics {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	clone := &StrongModePerformanceMetrics{
		TransactionCount:       pm.TransactionCount,
		SuccessfulTransactions: pm.SuccessfulTransactions,
		FailedTransactions:     pm.FailedTransactions,
		TotalTransactionTime:   pm.TotalTransactionTime,
		AvgTransactionTime:     pm.AvgTransactionTime,
		MinTransactionTime:     pm.MinTransactionTime,
		MaxTransactionTime:     pm.MaxTransactionTime,
		TotalFactsProcessed:    pm.TotalFactsProcessed,
		TotalFactsPersisted:    pm.TotalFactsPersisted,
		TotalFactsFailed:       pm.TotalFactsFailed,
		AvgFactsPerTransaction: pm.AvgFactsPerTransaction,
		AvgTimePerFact:         pm.AvgTimePerFact,
		TotalVerifications:     pm.TotalVerifications,
		SuccessfulVerifies:     pm.SuccessfulVerifies,
		FailedVerifies:         pm.FailedVerifies,
		TotalRetries:           pm.TotalRetries,
		AvgRetriesPerFact:      pm.AvgRetriesPerFact,
		MaxRetriesPerFact:      pm.MaxRetriesPerFact,
		TotalVerificationTime:  pm.TotalVerificationTime,
		AvgVerificationTime:    pm.AvgVerificationTime,
		TotalTimeouts:          pm.TotalTimeouts,
		TimeoutRate:            pm.TimeoutRate,
		AvgTimeToTimeout:       pm.AvgTimeToTimeout,
		TotalCommits:           pm.TotalCommits,
		SuccessfulCommits:      pm.SuccessfulCommits,
		FailedCommits:          pm.FailedCommits,
		AvgCommitTime:          pm.AvgCommitTime,
		TotalRollbacks:         pm.TotalRollbacks,
		IsHealthy:              pm.IsHealthy,
		HealthScore:            pm.HealthScore,
		PerformanceGrade:       pm.PerformanceGrade,
		StartTime:              pm.StartTime,
		LastUpdated:            pm.LastUpdated,
	}

	// Deep copy maps and slices
	clone.RollbackReasons = make(map[string]int)
	for k, v := range pm.RollbackReasons {
		clone.RollbackReasons[k] = v
	}

	clone.ConfigChangeHistory = make([]ConfigChange, len(pm.ConfigChangeHistory))
	copy(clone.ConfigChangeHistory, pm.ConfigChangeHistory)

	clone.Recommendations = make([]string, len(pm.Recommendations))
	copy(clone.Recommendations, pm.Recommendations)

	if pm.CurrentConfig != nil {
		clone.CurrentConfig = pm.CurrentConfig.Clone()
	}

	return clone
}
