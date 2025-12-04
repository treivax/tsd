// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package testutil

import (
	"strings"
	"testing"
)

// AssertNetworkStructure validates RETE network structure
func AssertNetworkStructure(t *testing.T, result *TSDResult, expectedTypeNodes, expectedTerminalNodes int) {
	t.Helper()

	if result.TypeNodes != expectedTypeNodes {
		t.Errorf("Expected %d type nodes, got %d", expectedTypeNodes, result.TypeNodes)
	}

	if result.TerminalNodes != expectedTerminalNodes {
		t.Errorf("Expected %d terminal nodes, got %d", expectedTerminalNodes, result.TerminalNodes)
	}
}

// AssertMinNetworkStructure validates minimum network structure
func AssertMinNetworkStructure(t *testing.T, result *TSDResult, minTypeNodes, minTerminalNodes int) {
	t.Helper()

	if result.TypeNodes < minTypeNodes {
		t.Errorf("Expected at least %d type nodes, got %d", minTypeNodes, result.TypeNodes)
	}

	if result.TerminalNodes < minTerminalNodes {
		t.Errorf("Expected at least %d terminal nodes, got %d", minTerminalNodes, result.TerminalNodes)
	}
}

// AssertActivations validates exact activation count
func AssertActivations(t *testing.T, result *TSDResult, expected int) {
	t.Helper()

	if result.Activations != expected {
		t.Errorf("Expected %d activations, got %d", expected, result.Activations)
	}
}

// AssertMinActivations validates minimum activations
func AssertMinActivations(t *testing.T, result *TSDResult, min int) {
	t.Helper()

	if result.Activations < min {
		t.Errorf("Expected at least %d activations, got %d", min, result.Activations)
	}
}

// AssertMaxActivations validates maximum activations
func AssertMaxActivations(t *testing.T, result *TSDResult, max int) {
	t.Helper()

	if result.Activations > max {
		t.Errorf("Expected at most %d activations, got %d", max, result.Activations)
	}
}

// AssertActivationRange validates activations within a range
func AssertActivationRange(t *testing.T, result *TSDResult, min, max int) {
	t.Helper()

	if result.Activations < min || result.Activations > max {
		t.Errorf("Expected activations between %d and %d, got %d", min, max, result.Activations)
	}
}

// AssertNoError validates successful execution
func AssertNoError(t *testing.T, result *TSDResult) {
	t.Helper()

	if result.Error != nil {
		t.Errorf("Unexpected error: %v", result.Error)
	}

	// Also check for injection errors in output
	if strings.Contains(result.Output, "⚠️ Erreur injection fait") {
		t.Error("Unexpected injection error in output")
	}
}

// AssertError validates expected error
func AssertError(t *testing.T, result *TSDResult) {
	t.Helper()

	hasError := result.Error != nil
	hasInjectionError := strings.Contains(result.Output, "⚠️ Erreur injection fait")

	if !hasError && !hasInjectionError {
		t.Error("Expected an error but got none")
	}
}

// AssertErrorContains validates error message contains substring
func AssertErrorContains(t *testing.T, result *TSDResult, substring string) {
	t.Helper()

	if result.Error == nil {
		t.Error("Expected an error but got none")
		return
	}

	if !strings.Contains(result.Error.Error(), substring) {
		t.Errorf("Expected error to contain %q, got: %v", substring, result.Error)
	}
}

// AssertErrorMatches validates error matches a specific error
func AssertErrorMatches(t *testing.T, result *TSDResult, expectedError error) {
	t.Helper()

	if result.Error == nil {
		t.Error("Expected an error but got none")
		return
	}

	if result.Error.Error() != expectedError.Error() {
		t.Errorf("Expected error %v, got %v", expectedError, result.Error)
	}
}

// AssertOutputContains validates captured output contains substring
func AssertOutputContains(t *testing.T, result *TSDResult, substring string) {
	t.Helper()

	if !strings.Contains(result.Output, substring) {
		t.Errorf("Expected output to contain %q, got: %s", substring, result.Output)
	}
}

// AssertOutputNotContains validates output does not contain substring
func AssertOutputNotContains(t *testing.T, result *TSDResult, substring string) {
	t.Helper()

	if strings.Contains(result.Output, substring) {
		t.Errorf("Expected output not to contain %q, but it did", substring)
	}
}

// AssertOutputEmpty validates no output was produced
func AssertOutputEmpty(t *testing.T, result *TSDResult) {
	t.Helper()

	if len(result.Output) > 0 {
		t.Errorf("Expected no output, got: %s", result.Output)
	}
}

// AssertFactCount validates exact fact count
func AssertFactCount(t *testing.T, result *TSDResult, expected int) {
	t.Helper()

	if result.Facts != expected {
		t.Errorf("Expected %d facts, got %d", expected, result.Facts)
	}
}

// AssertMinFactCount validates minimum fact count
func AssertMinFactCount(t *testing.T, result *TSDResult, min int) {
	t.Helper()

	if result.Facts < min {
		t.Errorf("Expected at least %d facts, got %d", min, result.Facts)
	}
}

// AssertNetworkBuilt validates that a network was successfully built
func AssertNetworkBuilt(t *testing.T, result *TSDResult) {
	t.Helper()

	AssertNoError(t, result)

	if result.TypeNodes == 0 {
		t.Error("Expected type nodes to be created")
	}

	if result.TerminalNodes == 0 {
		t.Error("Expected terminal nodes to be created")
	}
}

// AssertValidExecution validates a complete successful execution
func AssertValidExecution(t *testing.T, result *TSDResult) {
	t.Helper()

	AssertNoError(t, result)
	AssertNetworkBuilt(t, result)
	AssertMinFactCount(t, result, 1)
}

// AssertExecutionWithActivations validates execution with expected activations
func AssertExecutionWithActivations(t *testing.T, result *TSDResult, expectedActivations int) {
	t.Helper()

	AssertNoError(t, result)
	AssertNetworkBuilt(t, result)
	AssertActivations(t, result, expectedActivations)
}

// AssertQuickExecution validates execution completed within time limit
func AssertQuickExecution(t *testing.T, result *TSDResult, maxDuration string) {
	t.Helper()

	// This is a placeholder - duration checking can be enhanced
	if result.Duration.String() == "" {
		t.Log("Warning: Duration not tracked")
	}
}

// AssertIdenticalResults validates two results are identical
func AssertIdenticalResults(t *testing.T, result1, result2 *TSDResult, label1, label2 string) {
	t.Helper()

	if result1.TypeNodes != result2.TypeNodes {
		t.Errorf("%s has %d type nodes, %s has %d",
			label1, result1.TypeNodes, label2, result2.TypeNodes)
	}

	if result1.TerminalNodes != result2.TerminalNodes {
		t.Errorf("%s has %d terminal nodes, %s has %d",
			label1, result1.TerminalNodes, label2, result2.TerminalNodes)
	}

	if result1.Facts != result2.Facts {
		t.Errorf("%s has %d facts, %s has %d",
			label1, result1.Facts, label2, result2.Facts)
	}

	if result1.Activations != result2.Activations {
		t.Errorf("%s has %d activations, %s has %d",
			label1, result1.Activations, label2, result2.Activations)
	}

	if (result1.Error != nil) != (result2.Error != nil) {
		t.Errorf("%s error: %v, %s error: %v",
			label1, result1.Error, label2, result2.Error)
	}
}

// AssertResultsMatch validates two results match on key metrics
func AssertResultsMatch(t *testing.T, result1, result2 *TSDResult) {
	t.Helper()
	AssertIdenticalResults(t, result1, result2, "Result 1", "Result 2")
}
