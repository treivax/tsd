// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
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
