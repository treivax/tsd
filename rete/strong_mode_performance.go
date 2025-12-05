// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"time"
)

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
