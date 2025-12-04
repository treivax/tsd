// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"time"
)

// CoherenceMode defines the consistency guarantee level for transactions.
// Different modes trade consistency for performance based on application requirements.
type CoherenceMode int

const (
	// CoherenceModeStrong provides the strictest consistency guarantees.
	// All reads reflect the most recent writes. Maximum correctness.
	// - Synchronous verification of all facts
	// - Retry mechanism with exponential backoff
	// - Read-after-write consistency enforced
	// - Suitable for: financial transactions, critical business rules
	// - Performance: ~100-1,000 facts/sec
	CoherenceModeStrong CoherenceMode = iota

	// CoherenceModeRelaxed provides balanced consistency.
	// Reads may lag behind writes by a bounded time (configurable).
	// - Quick verification with limited retries
	// - Bounded staleness (default: 100ms)
	// - Still maintains eventual verification
	// - Suitable for: real-time analytics, recommendation engines
	// - Performance: ~1,000-10,000 facts/sec
	CoherenceModeRelaxed

	// CoherenceModeEventual provides eventual consistency.
	// Reads may return stale data temporarily.
	// - No synchronous verification
	// - Best-effort consistency
	// - Trust that storage will eventually reflect writes
	// - Suitable for: logging, telemetry, high-volume sensor data
	// - Performance: >10,000 facts/sec
	CoherenceModeEventual
)

// String returns the human-readable name of the coherence mode.
func (m CoherenceMode) String() string {
	switch m {
	case CoherenceModeStrong:
		return "Strong"
	case CoherenceModeRelaxed:
		return "Relaxed"
	case CoherenceModeEventual:
		return "Eventual"
	default:
		return fmt.Sprintf("Unknown(%d)", m)
	}
}

// IsValid returns true if the coherence mode is a valid value.
func (m CoherenceMode) IsValid() bool {
	return m >= CoherenceModeStrong && m <= CoherenceModeEventual
}

// TransactionOptions configures the behavior of a transaction based on its coherence mode.
type TransactionOptions struct {
	// CoherenceMode specifies the consistency guarantee level.
	CoherenceMode CoherenceMode

	// MaxStaleness defines the maximum acceptable staleness for reads.
	// Only applicable in Relaxed mode.
	// Zero means no staleness allowed (equivalent to Strong mode).
	MaxStaleness time.Duration

	// SkipVerification disables fact verification entirely.
	// Only applicable in Eventual mode.
	// Default: false for Strong/Relaxed, true for Eventual.
	SkipVerification bool

	// VerifyOnCommit controls whether all facts are verified on transaction commit.
	// Default: true for Strong/Relaxed, false for Eventual.
	VerifyOnCommit bool

	// QuickVerifyRetries specifies how many quick verification attempts to make.
	// Only applicable in Relaxed mode.
	// Default: 2
	QuickVerifyRetries int

	// QuickVerifyDelay is the delay between quick verification attempts.
	// Only applicable in Relaxed mode.
	// Default: 20ms
	QuickVerifyDelay time.Duration
}

// DefaultOptionsForMode returns sensible default options for each coherence mode.
func DefaultOptionsForMode(mode CoherenceMode) *TransactionOptions {
	switch mode {
	case CoherenceModeStrong:
		return &TransactionOptions{
			CoherenceMode:      CoherenceModeStrong,
			MaxStaleness:       0, // No staleness allowed
			SkipVerification:   false,
			VerifyOnCommit:     true,
			QuickVerifyRetries: 0, // Not used in Strong mode
			QuickVerifyDelay:   0, // Not used in Strong mode
		}

	case CoherenceModeRelaxed:
		return &TransactionOptions{
			CoherenceMode:      CoherenceModeRelaxed,
			MaxStaleness:       100 * time.Millisecond, // 100ms staleness acceptable
			SkipVerification:   false,
			VerifyOnCommit:     true,
			QuickVerifyRetries: 2,                     // 2 quick verification attempts
			QuickVerifyDelay:   20 * time.Millisecond, // 20ms between attempts
		}

	case CoherenceModeEventual:
		return &TransactionOptions{
			CoherenceMode:      CoherenceModeEventual,
			MaxStaleness:       time.Hour, // Effectively unlimited
			SkipVerification:   true,      // Skip verification entirely
			VerifyOnCommit:     false,     // Trust that it will eventually be consistent
			QuickVerifyRetries: 0,         // Not used in Eventual mode
			QuickVerifyDelay:   0,         // Not used in Eventual mode
		}

	default:
		// Fallback to Strong mode for unknown modes
		return DefaultOptionsForMode(CoherenceModeStrong)
	}
}

// Validate checks if the transaction options are valid and consistent.
func (opts *TransactionOptions) Validate() error {
	if !opts.CoherenceMode.IsValid() {
		return fmt.Errorf("invalid coherence mode: %v", opts.CoherenceMode)
	}

	if opts.CoherenceMode == CoherenceModeRelaxed {
		if opts.MaxStaleness < 0 {
			return fmt.Errorf("MaxStaleness cannot be negative: %v", opts.MaxStaleness)
		}
		if opts.QuickVerifyRetries < 0 {
			return fmt.Errorf("QuickVerifyRetries cannot be negative: %d", opts.QuickVerifyRetries)
		}
		if opts.QuickVerifyDelay < 0 {
			return fmt.Errorf("QuickVerifyDelay cannot be negative: %v", opts.QuickVerifyDelay)
		}
	}

	return nil
}

// Clone creates a deep copy of the transaction options.
func (opts *TransactionOptions) Clone() *TransactionOptions {
	if opts == nil {
		return nil
	}

	clone := *opts
	return &clone
}

// StrongModeConfig contains configuration for Strong consistency mode.
type StrongModeConfig struct {
	// SubmissionTimeout is the maximum time to wait for fact batch submission.
	// Default: 30 seconds
	SubmissionTimeout time.Duration

	// VerifyRetryDelay is the delay between verification retry attempts.
	// Default: 50ms
	VerifyRetryDelay time.Duration

	// MaxVerifyRetries is the maximum number of verification retry attempts.
	// Default: 10
	MaxVerifyRetries int
}

// DefaultStrongModeConfig returns the default configuration for Strong mode.
func DefaultStrongModeConfig() StrongModeConfig {
	return StrongModeConfig{
		SubmissionTimeout: 30 * time.Second,
		VerifyRetryDelay:  50 * time.Millisecond,
		MaxVerifyRetries:  10,
	}
}

// RelaxedModeConfig contains configuration for Relaxed consistency mode.
type RelaxedModeConfig struct {
	// MaxStaleness is the maximum acceptable staleness for reads.
	// Default: 100ms
	MaxStaleness time.Duration

	// QuickVerifyRetries is the number of quick verification attempts.
	// Default: 2
	QuickVerifyRetries int

	// QuickVerifyDelay is the delay between quick verification attempts.
	// Default: 20ms
	QuickVerifyDelay time.Duration

	// AsyncVerification enables background verification (future feature).
	// Default: false
	AsyncVerification bool
}

// DefaultRelaxedModeConfig returns the default configuration for Relaxed mode.
func DefaultRelaxedModeConfig() RelaxedModeConfig {
	return RelaxedModeConfig{
		MaxStaleness:       100 * time.Millisecond,
		QuickVerifyRetries: 2,
		QuickVerifyDelay:   20 * time.Millisecond,
		AsyncVerification:  false, // Not yet implemented
	}
}

// EventualModeConfig contains configuration for Eventual consistency mode.
type EventualModeConfig struct {
	// SkipVerification disables all verification.
	// Default: true
	SkipVerification bool

	// AsyncVerification enables background verification (future feature).
	// Default: false
	AsyncVerification bool

	// LogFailures controls whether storage failures are logged.
	// Default: true (always log failures even in best-effort mode)
	LogFailures bool
}

// DefaultEventualModeConfig returns the default configuration for Eventual mode.
func DefaultEventualModeConfig() EventualModeConfig {
	return EventualModeConfig{
		SkipVerification:  true,
		AsyncVerification: false, // Not yet implemented
		LogFailures:       true,  // Always log failures
	}
}

// NetworkCoherenceConfig contains global coherence configuration for a RETE network.
type NetworkCoherenceConfig struct {
	// DefaultCoherenceMode is the mode used when not explicitly specified.
	// Default: CoherenceModeStrong (for backward compatibility and safety)
	DefaultCoherenceMode CoherenceMode

	// StrongMode contains configuration specific to Strong consistency mode.
	StrongMode StrongModeConfig

	// RelaxedMode contains configuration specific to Relaxed consistency mode.
	RelaxedMode RelaxedModeConfig

	// EventualMode contains configuration specific to Eventual consistency mode.
	EventualMode EventualModeConfig

	// EnableMetrics controls whether coherence metrics are collected.
	// Default: true
	EnableMetrics bool
}

// DefaultNetworkCoherenceConfig returns the default global coherence configuration.
func DefaultNetworkCoherenceConfig() NetworkCoherenceConfig {
	return NetworkCoherenceConfig{
		DefaultCoherenceMode: CoherenceModeStrong, // Safe default
		StrongMode:           DefaultStrongModeConfig(),
		RelaxedMode:          DefaultRelaxedModeConfig(),
		EventualMode:         DefaultEventualModeConfig(),
		EnableMetrics:        true,
	}
}

// CoherenceMetrics tracks metrics for coherence operations.
type CoherenceMetrics struct {
	// Transaction counters by mode
	StrongTransactions   int64
	RelaxedTransactions  int64
	EventualTransactions int64

	// Verification metrics
	VerificationAttempts    int64
	VerificationSuccesses   int64
	VerificationFailures    int64
	QuickVerificationHits   int64
	QuickVerificationMisses int64

	// Staleness tracking (Relaxed mode)
	TotalStalenessObserved time.Duration
	MaxStalenessObserved   time.Duration
	StalenessObservations  int64

	// Performance metrics
	StrongModeTotalDuration   time.Duration
	RelaxedModeTotalDuration  time.Duration
	EventualModeTotalDuration time.Duration
}

// IncrementModeCounter increments the transaction counter for the given mode.
func (cm *CoherenceMetrics) IncrementModeCounter(mode CoherenceMode) {
	switch mode {
	case CoherenceModeStrong:
		cm.StrongTransactions++
	case CoherenceModeRelaxed:
		cm.RelaxedTransactions++
	case CoherenceModeEventual:
		cm.EventualTransactions++
	}
}

// RecordStaleness records an observed staleness duration for Relaxed mode.
func (cm *CoherenceMetrics) RecordStaleness(staleness time.Duration) {
	cm.TotalStalenessObserved += staleness
	cm.StalenessObservations++

	if staleness > cm.MaxStalenessObserved {
		cm.MaxStalenessObserved = staleness
	}
}

// AverageStaleness returns the average observed staleness.
func (cm *CoherenceMetrics) AverageStaleness() time.Duration {
	if cm.StalenessObservations == 0 {
		return 0
	}
	return cm.TotalStalenessObserved / time.Duration(cm.StalenessObservations)
}

// RecordTransactionDuration records the duration of a transaction for the given mode.
func (cm *CoherenceMetrics) RecordTransactionDuration(mode CoherenceMode, duration time.Duration) {
	switch mode {
	case CoherenceModeStrong:
		cm.StrongModeTotalDuration += duration
	case CoherenceModeRelaxed:
		cm.RelaxedModeTotalDuration += duration
	case CoherenceModeEventual:
		cm.EventualModeTotalDuration += duration
	}
}

// AverageTransactionDuration returns the average transaction duration for the given mode.
func (cm *CoherenceMetrics) AverageTransactionDuration(mode CoherenceMode) time.Duration {
	var totalDuration time.Duration
	var count int64

	switch mode {
	case CoherenceModeStrong:
		totalDuration = cm.StrongModeTotalDuration
		count = cm.StrongTransactions
	case CoherenceModeRelaxed:
		totalDuration = cm.RelaxedModeTotalDuration
		count = cm.RelaxedTransactions
	case CoherenceModeEventual:
		totalDuration = cm.EventualModeTotalDuration
		count = cm.EventualTransactions
	}

	if count == 0 {
		return 0
	}
	return totalDuration / time.Duration(count)
}
