// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// getSuccessRate calculates the transaction success rate as a percentage
func (pm *StrongModePerformanceMetrics) getSuccessRate() float64 {
	if pm.TransactionCount == 0 {
		return 0
	}
	return float64(pm.SuccessfulTransactions) / float64(pm.TransactionCount) * 100
}

// getFailureRate calculates the transaction failure rate as a percentage
func (pm *StrongModePerformanceMetrics) getFailureRate() float64 {
	if pm.TransactionCount == 0 {
		return 0
	}
	return float64(pm.FailedTransactions) / float64(pm.TransactionCount) * 100
}

// getFactPersistRate calculates the fact persistence rate as a percentage
func (pm *StrongModePerformanceMetrics) getFactPersistRate() float64 {
	if pm.TotalFactsProcessed == 0 {
		return 0
	}
	return float64(pm.TotalFactsPersisted) / float64(pm.TotalFactsProcessed) * 100
}

// getFactFailureRate calculates the fact failure rate as a percentage
func (pm *StrongModePerformanceMetrics) getFactFailureRate() float64 {
	if pm.TotalFactsProcessed == 0 {
		return 0
	}
	return float64(pm.TotalFactsFailed) / float64(pm.TotalFactsProcessed) * 100
}

// getVerifySuccessRate calculates the verification success rate as a percentage
func (pm *StrongModePerformanceMetrics) getVerifySuccessRate() float64 {
	if pm.TotalVerifications == 0 {
		return 0
	}
	return float64(pm.SuccessfulVerifies) / float64(pm.TotalVerifications) * 100
}

// getCommitSuccessRate calculates the commit success rate as a percentage
func (pm *StrongModePerformanceMetrics) getCommitSuccessRate() float64 {
	if pm.TotalCommits == 0 {
		return 0
	}
	return float64(pm.SuccessfulCommits) / float64(pm.TotalCommits) * 100
}
