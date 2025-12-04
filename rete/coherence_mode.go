// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"time"
)

// CoherenceMode defines the consistency guarantee level for transactions.
// Currently only Strong mode is implemented.
type CoherenceMode int

const (
	// CoherenceModeStrong provides the strictest consistency guarantees.
	// All reads reflect the most recent writes. Maximum correctness.
	// - Synchronous verification of all facts
	// - Retry mechanism with exponential backoff
	// - Read-after-write consistency enforced
	// - Suitable for: all production use cases
	// - Performance: ~100-1,000 facts/sec
	CoherenceModeStrong CoherenceMode = iota
)

// String returns the human-readable name of the coherence mode.
func (m CoherenceMode) String() string {
	switch m {
	case CoherenceModeStrong:
		return "Strong"
	default:
		return fmt.Sprintf("Unknown(%d)", m)
	}
}

// IsValid returns true if the coherence mode is a valid value.
func (m CoherenceMode) IsValid() bool {
	return m == CoherenceModeStrong
}

// TransactionOptions configures the behavior of a transaction.
type TransactionOptions struct {
	// SubmissionTimeout is the maximum time to wait for fact batch submission.
	// Default: 30 seconds
	SubmissionTimeout time.Duration

	// VerifyRetryDelay is the delay between verification retry attempts.
	// Default: 50ms
	VerifyRetryDelay time.Duration

	// MaxVerifyRetries is the maximum number of verification retry attempts.
	// Default: 10
	MaxVerifyRetries int

	// VerifyOnCommit controls whether all facts are verified on transaction commit.
	// Default: true
	VerifyOnCommit bool
}

// DefaultTransactionOptions returns the default transaction options.
func DefaultTransactionOptions() *TransactionOptions {
	return &TransactionOptions{
		SubmissionTimeout: 30 * time.Second,
		VerifyRetryDelay:  50 * time.Millisecond,
		MaxVerifyRetries:  10,
		VerifyOnCommit:    true,
	}
}

// Validate checks if the transaction options are valid.
func (opts *TransactionOptions) Validate() error {
	if opts.SubmissionTimeout < 0 {
		return fmt.Errorf("SubmissionTimeout cannot be negative: %v", opts.SubmissionTimeout)
	}
	if opts.VerifyRetryDelay < 0 {
		return fmt.Errorf("VerifyRetryDelay cannot be negative: %v", opts.VerifyRetryDelay)
	}
	if opts.MaxVerifyRetries < 0 {
		return fmt.Errorf("MaxVerifyRetries cannot be negative: %d", opts.MaxVerifyRetries)
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

// NetworkCoherenceConfig contains global coherence configuration for a RETE network.
type NetworkCoherenceConfig struct {
	// DefaultOptions are the default transaction options.
	DefaultOptions TransactionOptions

	// EnableMetrics controls whether coherence metrics are collected.
	// Default: true
	EnableMetrics bool
}

// DefaultNetworkCoherenceConfig returns the default global coherence configuration.
func DefaultNetworkCoherenceConfig() NetworkCoherenceConfig {
	return NetworkCoherenceConfig{
		DefaultOptions: *DefaultTransactionOptions(),
		EnableMetrics:  true,
	}
}

// Note: CoherenceMetrics is defined in coherence_metrics.go
// Strong mode will use the existing comprehensive metrics structure
