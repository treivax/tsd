// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
	"time"
)

// TestDefaultAdvancedPipelineConfig tests the default configuration function
func TestDefaultAdvancedPipelineConfig(t *testing.T) {
	t.Run("returns non-nil config", func(t *testing.T) {
		config := DefaultAdvancedPipelineConfig()
		if config == nil {
			t.Fatal("DefaultAdvancedPipelineConfig returned nil")
		}
	})

	t.Run("has reasonable default values", func(t *testing.T) {
		config := DefaultAdvancedPipelineConfig()

		// Check transaction timeout is set
		if config.TransactionTimeout <= 0 {
			t.Errorf("TransactionTimeout should be positive, got %v", config.TransactionTimeout)
		}

		// Check max transaction size is reasonable
		if config.MaxTransactionSize <= 0 {
			t.Errorf("MaxTransactionSize should be positive, got %d", config.MaxTransactionSize)
		}

		// Verify default values
		expectedTimeout := 30 * time.Second
		if config.TransactionTimeout != expectedTimeout {
			t.Errorf("Expected TransactionTimeout %v, got %v", expectedTimeout, config.TransactionTimeout)
		}

		expectedMaxSize := int64(100 * 1024 * 1024) // 100 MB
		if config.MaxTransactionSize != expectedMaxSize {
			t.Errorf("Expected MaxTransactionSize %d, got %d", expectedMaxSize, config.MaxTransactionSize)
		}

		// Check auto-rollback is enabled by default
		if !config.AutoRollbackOnError {
			t.Error("AutoRollbackOnError should be true by default")
		}

		// Check auto-commit is disabled by default
		if config.AutoCommit {
			t.Error("AutoCommit should be false by default")
		}
	})

	t.Run("multiple calls return independent configs", func(t *testing.T) {
		config1 := DefaultAdvancedPipelineConfig()
		config2 := DefaultAdvancedPipelineConfig()

		// Should be different instances
		if config1 == config2 {
			t.Error("Multiple calls should return different config instances")
		}

		// Modify one and verify the other is unchanged
		config1.TransactionTimeout = 1 * time.Minute
		if config2.TransactionTimeout == 1*time.Minute {
			t.Error("Modifying one config should not affect another")
		}
	})
}

// TestAdvancedPipelineConfig_Modification tests config modification
func TestAdvancedPipelineConfig_Modification(t *testing.T) {
	t.Run("can modify transaction timeout", func(t *testing.T) {
		config := DefaultAdvancedPipelineConfig()
		newTimeout := 1 * time.Minute
		config.TransactionTimeout = newTimeout

		if config.TransactionTimeout != newTimeout {
			t.Errorf("Expected timeout %v, got %v", newTimeout, config.TransactionTimeout)
		}
	})

	t.Run("can modify max transaction size", func(t *testing.T) {
		config := DefaultAdvancedPipelineConfig()
		newSize := int64(500 * 1024 * 1024) // 500 MB
		config.MaxTransactionSize = newSize

		if config.MaxTransactionSize != newSize {
			t.Errorf("Expected size %d, got %d", newSize, config.MaxTransactionSize)
		}
	})

	t.Run("can toggle auto-commit", func(t *testing.T) {
		config := DefaultAdvancedPipelineConfig()
		config.AutoCommit = true

		if !config.AutoCommit {
			t.Error("AutoCommit should be true after setting")
		}

		config.AutoCommit = false
		if config.AutoCommit {
			t.Error("AutoCommit should be false after toggling")
		}
	})

	t.Run("can toggle auto-rollback", func(t *testing.T) {
		config := DefaultAdvancedPipelineConfig()
		config.AutoRollbackOnError = false

		if config.AutoRollbackOnError {
			t.Error("AutoRollbackOnError should be false after setting")
		}

		config.AutoRollbackOnError = true
		if !config.AutoRollbackOnError {
			t.Error("AutoRollbackOnError should be true after toggling")
		}
	})
}

// TestAdvancedMetrics_Structure tests the AdvancedMetrics struct
func TestAdvancedMetrics_Structure(t *testing.T) {
	t.Run("can create and populate metrics", func(t *testing.T) {
		metrics := &AdvancedMetrics{
			ValidationWithContextDuration: 100 * time.Millisecond,
			TypesFoundInContext:           5,
			ValidationErrors:              []string{"error1", "error2"},
			GCDuration:                    50 * time.Millisecond,
			NodesCollected:                10,
			MemoryFreed:                   1024,
			GCPerformed:                   true,
			TransactionID:                 "tx_123",
			TransactionFootprint:          2048,
			ChangesTracked:                3,
			RollbackPerformed:             false,
			RollbackDuration:              0,
			TransactionDuration:           200 * time.Millisecond,
		}

		// Verify all fields are accessible
		if metrics.ValidationWithContextDuration != 100*time.Millisecond {
			t.Error("ValidationWithContextDuration not set correctly")
		}
		if metrics.TypesFoundInContext != 5 {
			t.Error("TypesFoundInContext not set correctly")
		}
		if len(metrics.ValidationErrors) != 2 {
			t.Error("ValidationErrors not set correctly")
		}
		if !metrics.GCPerformed {
			t.Error("GCPerformed not set correctly")
		}
		if metrics.TransactionID != "tx_123" {
			t.Error("TransactionID not set correctly")
		}
	})

	t.Run("metrics with zero values", func(t *testing.T) {
		metrics := &AdvancedMetrics{}

		if metrics.ValidationWithContextDuration != 0 {
			t.Error("Zero-value duration should be 0")
		}
		if metrics.TypesFoundInContext != 0 {
			t.Error("Zero-value int should be 0")
		}
		if metrics.ValidationErrors != nil {
			t.Error("Zero-value slice should be nil")
		}
		if metrics.GCPerformed {
			t.Error("Zero-value bool should be false")
		}
	})
}

// TestAdvancedPipelineConfig_EdgeCases tests edge case values
func TestAdvancedPipelineConfig_EdgeCases(t *testing.T) {
	t.Run("zero timeout", func(t *testing.T) {
		config := DefaultAdvancedPipelineConfig()
		config.TransactionTimeout = 0

		if config.TransactionTimeout != 0 {
			t.Error("Should allow zero timeout")
		}
	})

	t.Run("very large transaction size", func(t *testing.T) {
		config := DefaultAdvancedPipelineConfig()
		config.MaxTransactionSize = 1024 * 1024 * 1024 * 10 // 10 GB

		if config.MaxTransactionSize != 1024*1024*1024*10 {
			t.Error("Should allow large transaction size")
		}
	})

	t.Run("negative timeout", func(t *testing.T) {
		config := DefaultAdvancedPipelineConfig()
		config.TransactionTimeout = -1 * time.Second

		// Should be able to set it (validation would happen at usage time)
		if config.TransactionTimeout != -1*time.Second {
			t.Error("Should allow setting negative timeout")
		}
	})
}

// TestAdvancedMetrics_Calculations tests metric calculations
func TestAdvancedMetrics_Calculations(t *testing.T) {
	t.Run("total operation time", func(t *testing.T) {
		metrics := &AdvancedMetrics{
			ValidationWithContextDuration: 100 * time.Millisecond,
			GCDuration:                    50 * time.Millisecond,
			TransactionDuration:           200 * time.Millisecond,
		}

		totalTime := metrics.ValidationWithContextDuration +
			metrics.GCDuration +
			metrics.TransactionDuration

		expectedTotal := 350 * time.Millisecond
		if totalTime != expectedTotal {
			t.Errorf("Expected total time %v, got %v", expectedTotal, totalTime)
		}
	})

	t.Run("rollback overhead", func(t *testing.T) {
		metrics := &AdvancedMetrics{
			TransactionDuration: 200 * time.Millisecond,
			RollbackDuration:    50 * time.Millisecond,
			RollbackPerformed:   true,
		}

		if !metrics.RollbackPerformed {
			t.Error("RollbackPerformed should be true")
		}

		if metrics.RollbackDuration > metrics.TransactionDuration {
			t.Error("Rollback duration should typically be less than transaction duration")
		}
	})
}
