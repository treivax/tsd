// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"time"
)

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
