// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"testing"
)

// TestResetInstruction_StoragePreservation verifies that the Storage reference
// is preserved across resets (important for performance)
func TestResetInstruction_StoragePreservation(t *testing.T) {
	helper := NewTestHelper()

	constraintFile := "../../constraint/test/integration/reset_integration_test.tsd"
	factsFile := "../../constraint/test/integration/reset_integration_test.tsd"

	network, _, storage := helper.BuildNetworkFromConstraintFileWithMassiveFacts(t, constraintFile, factsFile)

	// Verify storage is the same instance (preserved)
	if storage == nil {
		t.Fatal("Storage should not be nil")
	}

	// Verify network has a valid storage reference
	if network.Storage == nil {
		t.Error("Network storage should not be nil after reset")
	}

	// Both should point to the same storage
	if network.Storage != storage {
		t.Error("Network storage and returned storage should be the same instance")
	}

	t.Logf("✅ Storage preservation test passed: storage correctly preserved across reset")
}
