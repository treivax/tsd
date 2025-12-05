// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

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

// getHealthStatus returns a formatted health status string
func (pm *StrongModePerformanceMetrics) getHealthStatus() string {
	if pm.IsHealthy {
		return "✅ Healthy"
	}
	return "⚠️  Needs Attention"
}
