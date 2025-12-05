// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"time"
)

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

// formatRecommendations formats recommendations for report display
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
