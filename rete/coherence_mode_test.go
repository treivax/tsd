// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCoherenceMode_String tests the String() method
func TestCoherenceMode_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		mode     CoherenceMode
		expected string
	}{
		{
			name:     "Strong mode",
			mode:     CoherenceModeStrong,
			expected: "Strong",
		},
		{
			name:     "Unknown mode",
			mode:     CoherenceMode(99),
			expected: "Unknown(99)",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.mode.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCoherenceMode_IsValid tests mode validation
func TestCoherenceMode_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		mode     CoherenceMode
		expected bool
	}{
		{
			name:     "Strong is valid",
			mode:     CoherenceModeStrong,
			expected: true,
		},
		{
			name:     "Invalid mode 1",
			mode:     CoherenceMode(1),
			expected: false,
		},
		{
			name:     "Invalid mode -1",
			mode:     CoherenceMode(-1),
			expected: false,
		},
		{
			name:     "Invalid mode 99",
			mode:     CoherenceMode(99),
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.mode.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestDefaultTransactionOptions tests default options
func TestDefaultTransactionOptions(t *testing.T) {
	t.Parallel()

	opts := DefaultTransactionOptions()

	require.NotNil(t, opts)
	assert.Equal(t, 30*time.Second, opts.SubmissionTimeout)
	assert.Equal(t, 50*time.Millisecond, opts.VerifyRetryDelay)
	assert.Equal(t, 10, opts.MaxVerifyRetries)
	assert.True(t, opts.VerifyOnCommit)
}

// TestTransactionOptions_Validate tests option validation
func TestTransactionOptions_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		opts      *TransactionOptions
		wantError bool
		errorMsg  string
	}{
		{
			name: "Valid options",
			opts: &TransactionOptions{
				SubmissionTimeout: 30 * time.Second,
				VerifyRetryDelay:  50 * time.Millisecond,
				MaxVerifyRetries:  10,
				VerifyOnCommit:    true,
			},
			wantError: false,
		},
		{
			name: "Negative SubmissionTimeout",
			opts: &TransactionOptions{
				SubmissionTimeout: -1 * time.Second,
				VerifyRetryDelay:  50 * time.Millisecond,
				MaxVerifyRetries:  10,
			},
			wantError: true,
			errorMsg:  "SubmissionTimeout cannot be negative",
		},
		{
			name: "Negative VerifyRetryDelay",
			opts: &TransactionOptions{
				SubmissionTimeout: 30 * time.Second,
				VerifyRetryDelay:  -50 * time.Millisecond,
				MaxVerifyRetries:  10,
			},
			wantError: true,
			errorMsg:  "VerifyRetryDelay cannot be negative",
		},
		{
			name: "Negative MaxVerifyRetries",
			opts: &TransactionOptions{
				SubmissionTimeout: 30 * time.Second,
				VerifyRetryDelay:  50 * time.Millisecond,
				MaxVerifyRetries:  -5,
			},
			wantError: true,
			errorMsg:  "MaxVerifyRetries cannot be negative",
		},
		{
			name: "Zero values are valid",
			opts: &TransactionOptions{
				SubmissionTimeout: 0,
				VerifyRetryDelay:  0,
				MaxVerifyRetries:  0,
				VerifyOnCommit:    false,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.opts.Validate()

			if tt.wantError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestTransactionOptions_Clone tests option cloning
func TestTransactionOptions_Clone(t *testing.T) {
	t.Parallel()

	t.Run("Clone nil", func(t *testing.T) {
		t.Parallel()

		var opts *TransactionOptions
		clone := opts.Clone()
		assert.Nil(t, clone)
	})

	t.Run("Clone valid options", func(t *testing.T) {
		t.Parallel()

		opts := &TransactionOptions{
			SubmissionTimeout: 45 * time.Second,
			VerifyRetryDelay:  75 * time.Millisecond,
			MaxVerifyRetries:  15,
			VerifyOnCommit:    false,
		}

		clone := opts.Clone()

		require.NotNil(t, clone)
		assert.Equal(t, opts.SubmissionTimeout, clone.SubmissionTimeout)
		assert.Equal(t, opts.VerifyRetryDelay, clone.VerifyRetryDelay)
		assert.Equal(t, opts.MaxVerifyRetries, clone.MaxVerifyRetries)
		assert.Equal(t, opts.VerifyOnCommit, clone.VerifyOnCommit)

		// Verify it's a deep copy
		clone.MaxVerifyRetries = 20
		assert.Equal(t, 15, opts.MaxVerifyRetries, "Original should not be modified")
	})
}

// TestDefaultNetworkCoherenceConfig tests default network config
func TestDefaultNetworkCoherenceConfig(t *testing.T) {
	t.Parallel()

	config := DefaultNetworkCoherenceConfig()

	require.NotNil(t, config)
	assert.Equal(t, 30*time.Second, config.DefaultOptions.SubmissionTimeout)
	assert.Equal(t, 50*time.Millisecond, config.DefaultOptions.VerifyRetryDelay)
	assert.Equal(t, 10, config.DefaultOptions.MaxVerifyRetries)
	assert.True(t, config.DefaultOptions.VerifyOnCommit)
	assert.True(t, config.EnableMetrics)
}

// TestCoherenceMetrics_Basic tests basic metrics tracking
func TestCoherenceMetrics_Basic(t *testing.T) {
	t.Parallel()

	collector := NewCoherenceMetricsCollector()
	metrics := collector.GetMetrics()

	require.NotNil(t, metrics)
	assert.Equal(t, 0, metrics.FactsSubmitted)
	assert.Equal(t, 0, metrics.TotalVerifyAttempts)
	assert.Equal(t, 0, metrics.TotalRetries)
}

// TestTransactionWithOptions tests transaction creation with options
func TestTransactionWithOptions(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t, WithLogLevel(LogLevelSilent))
	defer env.Cleanup()

	t.Run("Default options", func(t *testing.T) {
		tx := env.Network.BeginTransaction()

		require.NotNil(t, tx)
		require.NotNil(t, tx.Options)
		assert.Equal(t, 30*time.Second, tx.Options.SubmissionTimeout)
		assert.Equal(t, 50*time.Millisecond, tx.Options.VerifyRetryDelay)
		assert.Equal(t, 10, tx.Options.MaxVerifyRetries)
		assert.True(t, tx.Options.VerifyOnCommit)
	})

	t.Run("Custom options", func(t *testing.T) {
		opts := &TransactionOptions{
			SubmissionTimeout: 60 * time.Second,
			VerifyRetryDelay:  100 * time.Millisecond,
			MaxVerifyRetries:  20,
			VerifyOnCommit:    false,
		}

		tx := env.Network.BeginTransactionWithOptions(opts)

		require.NotNil(t, tx)
		require.NotNil(t, tx.Options)
		assert.Equal(t, 60*time.Second, tx.Options.SubmissionTimeout)
		assert.Equal(t, 100*time.Millisecond, tx.Options.VerifyRetryDelay)
		assert.Equal(t, 20, tx.Options.MaxVerifyRetries)
		assert.False(t, tx.Options.VerifyOnCommit)
	})

	t.Run("Nil options uses defaults", func(t *testing.T) {
		tx := env.Network.BeginTransactionWithOptions(nil)

		require.NotNil(t, tx)
		require.NotNil(t, tx.Options)
		assert.Equal(t, 30*time.Second, tx.Options.SubmissionTimeout)
	})
}

// TestCoherenceMetrics_ConcurrentAccess tests thread-safe metrics
func TestCoherenceMetrics_ConcurrentAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent test in short mode")
	}

	t.Parallel()

	collector := NewCoherenceMetricsCollector()
	done := make(chan bool)

	// Spawn multiple goroutines writing to metrics
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				collector.RecordFactSubmitted()
				collector.RecordFactPersisted()
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify counts (should be thread-safe)
	metrics := collector.GetMetrics()
	assert.Equal(t, 1000, metrics.FactsSubmitted)
	assert.Equal(t, 1000, metrics.FactsPersisted)
}

// TestTransactionOptions_EdgeCases tests edge cases
func TestTransactionOptions_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("Very large timeout", func(t *testing.T) {
		t.Parallel()

		opts := &TransactionOptions{
			SubmissionTimeout: 24 * time.Hour,
			VerifyRetryDelay:  1 * time.Second,
			MaxVerifyRetries:  1000,
			VerifyOnCommit:    true,
		}

		err := opts.Validate()
		assert.NoError(t, err)
	})

	t.Run("Zero timeout valid", func(t *testing.T) {
		t.Parallel()

		opts := &TransactionOptions{
			SubmissionTimeout: 0,
			VerifyRetryDelay:  0,
			MaxVerifyRetries:  0,
			VerifyOnCommit:    false,
		}

		err := opts.Validate()
		assert.NoError(t, err)
	})

	t.Run("Very small delays", func(t *testing.T) {
		t.Parallel()

		opts := &TransactionOptions{
			SubmissionTimeout: 1 * time.Nanosecond,
			VerifyRetryDelay:  1 * time.Nanosecond,
			MaxVerifyRetries:  1,
			VerifyOnCommit:    true,
		}

		err := opts.Validate()
		assert.NoError(t, err)
	})
}
