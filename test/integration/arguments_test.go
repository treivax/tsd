package main

import (
	"testing"

	"github.com/treivax/tsd/rete"
)

// TestVariableArguments tests the ability to pass complete variables as action arguments
func TestVariableArguments(t *testing.T) {
	helper := NewTestHelper()

	constraintFile := "../../constraint/test/integration/variable_action_test.constraint"
	factsFile := "../../constraint/test/integration/variable_action_test.facts"

	// Build network from constraint and facts
	network, facts, _ := helper.BuildNetworkFromConstraintFileWithMassiveFacts(t, constraintFile, factsFile)

	// Validate basic structure
	if len(facts) == 0 {
		t.Fatal("No facts were parsed")
	}

	if len(network.TerminalNodes) == 0 {
		t.Fatal("No terminal nodes created")
	}

	// Verify that actions have variable arguments
	foundVariableArg := false
	for _, terminal := range network.TerminalNodes {
		if terminal.Action != nil {
			for _, arg := range terminal.Action.Job.Args {
				if argMap, ok := arg.(map[string]interface{}); ok {
					if argType, hasType := argMap["type"]; hasType && argType == "variable" {
						foundVariableArg = true
						break
					}
				}
			}
		}
	}

	if !foundVariableArg {
		t.Error("Expected to find at least one variable argument in actions")
	}

	t.Logf("✅ Variable arguments test passed: %d facts, %d terminal nodes", len(facts), len(network.TerminalNodes))
}

// TestComprehensiveMixedArguments tests all combinations of arguments: variables, fields, and mixed
func TestComprehensiveMixedArguments(t *testing.T) {
	helper := NewTestHelper()

	constraintFile := "../../constraint/test/integration/comprehensive_args_test.constraint"
	factsFile := "../../constraint/test/integration/comprehensive_args_test.facts"

	network, facts, _ := helper.BuildNetworkFromConstraintFileWithMassiveFacts(t, constraintFile, factsFile)

	// Track different argument types found
	foundVariableArg := false
	foundFieldArg := false

	for _, terminal := range network.TerminalNodes {
		if terminal.Action != nil {
			for _, arg := range terminal.Action.Job.Args {
				if argMap, ok := arg.(map[string]interface{}); ok {
					if argType, hasType := argMap["type"]; hasType {
						switch argType {
						case "variable":
							foundVariableArg = true
						case "fieldAccess":
							foundFieldArg = true
						}
					}
				}
			}
		}
	}

	// Verify we found both types
	if !foundVariableArg {
		t.Error("Expected to find variable arguments")
	}
	if !foundFieldArg {
		t.Error("Expected to find field access arguments")
	}

	t.Logf("✅ Mixed arguments test passed: %d facts, %d terminal nodes, found variables=%v, fields=%v",
		len(facts), len(network.TerminalNodes), foundVariableArg, foundFieldArg)
}

// TestErrorDetectionInArguments tests that invalid arguments are properly detected
func TestErrorDetectionInArguments(t *testing.T) {
	helper := NewTestHelper()

	constraintFile := "../../constraint/test/integration/error_args_test.constraint"
	factsFile := "../../constraint/test/integration/error_args_test.facts"

	// This should fail due to validation errors in the constraint file
	// Call the pipeline directly instead of using the helper to check for errors
	storage := rete.NewMemoryStorage()
	network, facts, err := helper.pipeline.BuildNetworkFromConstraintFileWithFacts(constraintFile, factsFile, storage)

	// We EXPECT an error to be returned (the file contains intentional errors)
	if err == nil {
		t.Errorf("❌ Expected error to be detected in constraint file, but network was built successfully with %d facts, %d terminal nodes",
			len(facts), len(network.TerminalNodes))
		return
	}

	// Verify the error message indicates a validation problem
	errMsg := err.Error()
	if !containsAny(errMsg, []string{"validation", "incompatibilité", "unknown", "invalid"}) {
		t.Errorf("❌ Error was detected but message doesn't indicate a validation error: %v", err)
		return
	}

	t.Logf("✅ Error detection test passed: validation error correctly detected: %v", err)
}

// containsAny checks if a string contains any of the provided substrings
func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

// TestBasicNetworkIntegrity validates the basic integrity of a simple RETE network
func TestBasicNetworkIntegrity(t *testing.T) {
	helper := NewTestHelper()

	// Use an existing simple constraint file for testing
	constraintFile := "../../constraint/test/integration/variable_action_test.constraint"
	factsFile := "../../constraint/test/integration/variable_action_test.facts"

	network, facts, storage := helper.BuildNetworkFromConstraintFileWithMassiveFacts(t, constraintFile, factsFile)

	// Verify network structure
	if len(network.TypeNodes) == 0 {
		t.Error("No type nodes created")
	}

	if len(network.TerminalNodes) == 0 {
		t.Error("No terminal nodes created")
	}

	if len(facts) == 0 {
		t.Error("No facts were created")
	}

	// Verify storage is not nil
	if storage == nil {
		t.Error("Storage should not be nil")
	}

	t.Logf("✅ Basic network integrity test passed: %d type nodes, %d terminal nodes, %d facts",
		len(network.TypeNodes), len(network.TerminalNodes), len(facts))
}
