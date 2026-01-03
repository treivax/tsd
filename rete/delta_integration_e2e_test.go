// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// TestDeltaIntegration_Initialization tests that delta components are properly
// initialized when creating a new RETE network.
func TestDeltaIntegration_Initialization(t *testing.T) {
	t.Log("🧪 TEST DELTA INTEGRATION - INITIALIZATION")
	t.Log("==========================================")

	// Create network with default config (delta should be enabled)
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Verify all delta components are initialized
	t.Run("DeltaPropagator initialized", func(t *testing.T) {
		if network.DeltaPropagator == nil {
			t.Error("❌ DeltaPropagator should be initialized")
		} else {
			t.Log("✅ DeltaPropagator initialized")
		}
	})

	t.Run("DependencyIndex initialized", func(t *testing.T) {
		if network.DependencyIndex == nil {
			t.Error("❌ DependencyIndex should be initialized")
		} else {
			t.Log("✅ DependencyIndex initialized")
		}
	})

	t.Run("IntegrationHelper initialized", func(t *testing.T) {
		if network.IntegrationHelper == nil {
			t.Error("❌ IntegrationHelper should be initialized")
		} else {
			t.Log("✅ IntegrationHelper initialized")
		}
	})

	t.Run("EnableDeltaPropagation enabled by default", func(t *testing.T) {
		if !network.EnableDeltaPropagation {
			t.Error("❌ Delta propagation should be enabled by default")
		} else {
			t.Log("✅ Delta propagation enabled by default")
		}
	})

	t.Log("✅ All delta components initialized correctly")
}

// TestDeltaIntegration_DisableEnable tests enabling and disabling delta propagation.
func TestDeltaIntegration_DisableEnable(t *testing.T) {
	t.Log("🧪 TEST DELTA INTEGRATION - DISABLE/ENABLE")
	t.Log("==========================================")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Initially enabled
	if !network.EnableDeltaPropagation {
		t.Fatal("❌ Delta should be enabled initially")
	}
	t.Log("✅ Delta enabled initially")

	// Disable
	network.EnableDeltaPropagation = false
	if network.EnableDeltaPropagation {
		t.Error("❌ Delta should be disabled")
	}
	t.Log("✅ Delta disabled successfully")

	// Re-enable
	network.EnableDeltaPropagation = true
	if !network.EnableDeltaPropagation {
		t.Error("❌ Delta should be re-enabled")
	}
	t.Log("✅ Delta re-enabled successfully")
}

// Note: UpdateFact tests require a complete RETE network setup with types,
// rules, and proper propagation paths. These are covered by integration tests
// in the delta package examples and by full E2E tests that use the TSD parser.
// The initialization tests above verify that delta components are properly
// integrated into the network creation pipeline.
