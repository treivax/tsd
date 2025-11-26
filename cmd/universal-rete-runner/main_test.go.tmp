// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintHeader(t *testing.T) {
	var buf bytes.Buffer
	PrintHeader(&buf)

	output := buf.String()

	if !strings.Contains(output, "RUNNER UNIVERSEL") {
		t.Error("Header should contain 'RUNNER UNIVERSEL'")
	}
	if !strings.Contains(output, "TESTS COMPLETS") {
		t.Error("Header should contain 'TESTS COMPLETS'")
	}
	if !strings.Contains(output, "Date:") {
		t.Error("Header should contain date")
	}
}

func TestGetErrorTests(t *testing.T) {
	errorTests := GetErrorTests()

	if len(errorTests) == 0 {
		t.Error("GetErrorTests should return at least one error test")
	}

	if !errorTests["error_args_test"] {
		t.Error("error_args_test should be marked as an error test")
	}
}

func TestDiscoverTests(t *testing.T) {
	tests := DiscoverTests()

	// We should find at least some tests if the directories exist
	// Note: This test might fail if run outside the project directory
	if len(tests) > 0 {
		// Validate structure of discovered tests
		for _, test := range tests {
			if test.Name == "" {
				t.Error("Test name should not be empty")
			}
			if test.Category == "" {
				t.Error("Test category should not be empty")
			}
			if test.Constraint == "" {
				t.Error("Test constraint file should not be empty")
			}
			if test.Facts == "" {
				t.Error("Test facts file should not be empty")
			}
		}
	}
}

func TestExecuteTest(t *testing.T) {
	tests := []struct {
		name        string
		testFile    TestFile
		expectError bool
		checkResult func(*testing.T, TestResult)
	}{
		{
			name: "nonexistent files",
			testFile: TestFile{
				Name:       "nonexistent_test",
				Category:   "test",
				Constraint: "nonexistent.constraint",
				Facts:      "nonexistent.facts",
			},
			expectError: true,
			checkResult: func(t *testing.T, result TestResult) {
				if result.Error == nil {
					t.Error("Should have error for nonexistent files")
				}
				// When expecting error, test should pass if error occurred
				if !result.Passed && result.Error != nil {
					// This is correct behavior for error tests
				} else if result.Passed && result.Error == nil {
					t.Error("Test should not pass without error when error expected")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExecuteTest(tt.testFile, tt.expectError)

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}

			if result.Name != tt.testFile.Name {
				t.Errorf("Result name = %q, want %q", result.Name, tt.testFile.Name)
			}
			if result.Category != tt.testFile.Category {
				t.Errorf("Result category = %q, want %q", result.Category, tt.testFile.Category)
			}
		})
	}
}

func TestExecuteTestWithErrorExpected(t *testing.T) {
	testFile := TestFile{
		Name:       "error_test",
		Category:   "error",
		Constraint: "nonexistent.constraint",
		Facts:      "nonexistent.facts",
	}

	result := ExecuteTest(testFile, true)

	if result.Error != nil && !result.Passed {
		t.Error("Test expecting error should pass when error occurs")
	}
}

func TestRunTests(t *testing.T) {
	var stdout bytes.Buffer

	result := RunTests(&stdout)

	if result == nil {
		t.Fatal("RunTests should return a result")
	}

	if result.Total < 0 {
		t.Error("Total should not be negative")
	}

	if result.Passed < 0 {
		t.Error("Passed should not be negative")
	}

	if result.Failed < 0 {
		t.Error("Failed should not be negative")
	}

	if result.Passed+result.Failed != result.Total {
		t.Errorf("Passed (%d) + Failed (%d) should equal Total (%d)",
			result.Passed, result.Failed, result.Total)
	}

	output := stdout.String()
	if !strings.Contains(output, "RUNNER UNIVERSEL") {
		t.Error("Output should contain header")
	}
}

func TestRun(t *testing.T) {
	var stdout, stderr bytes.Buffer

	exitCode := Run(&stdout, &stderr)

	// Exit code should be 0 if all tests pass, 1 if any fail
	if exitCode != 0 && exitCode != 1 {
		t.Errorf("exitCode = %d, want 0 or 1", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "Résumé") {
		t.Error("Output should contain summary")
	}
	if !strings.Contains(output, "tests") {
		t.Error("Output should mention tests")
	}
}

func TestRunWithAllTestsPassing(t *testing.T) {
	// This test simulates the scenario where all tests pass
	// We can't fully test this without mock data, but we can test the structure

	var stdout, stderr bytes.Buffer
	exitCode := Run(&stdout, &stderr)

	output := stdout.String()

	// Should contain summary information
	if !strings.Contains(output, "réussis") || !strings.Contains(output, "échoués") {
		t.Error("Output should contain pass/fail counts")
	}

	// Exit code should be valid
	if exitCode < 0 || exitCode > 1 {
		t.Errorf("Invalid exit code: %d", exitCode)
	}
}

func TestTestResultStructure(t *testing.T) {
	result := TestResult{
		Name:          "test1",
		Category:      "alpha",
		Passed:        true,
		TypeNodes:     5,
		TerminalNodes: 2,
		Facts:         10,
		Activations:   3,
		Error:         nil,
		Output:        "test output",
	}

	if result.Name != "test1" {
		t.Error("TestResult should store name")
	}
	if result.Category != "alpha" {
		t.Error("TestResult should store category")
	}
	if !result.Passed {
		t.Error("TestResult should store passed status")
	}
	if result.TypeNodes != 5 {
		t.Error("TestResult should store type nodes count")
	}
}

func TestRunResultStructure(t *testing.T) {
	result := RunResult{
		Total:  10,
		Passed: 8,
		Failed: 2,
		Results: []TestResult{
			{Name: "test1", Passed: true},
			{Name: "test2", Passed: false},
		},
	}

	if result.Total != 10 {
		t.Error("RunResult should store total")
	}
	if result.Passed != 8 {
		t.Error("RunResult should store passed count")
	}
	if result.Failed != 2 {
		t.Error("RunResult should store failed count")
	}
	if len(result.Results) != 2 {
		t.Error("RunResult should store individual results")
	}
}

func TestTestFileStructure(t *testing.T) {
	testFile := TestFile{
		Name:       "test1",
		Category:   "alpha",
		Constraint: "test.constraint",
		Facts:      "test.facts",
	}

	if testFile.Name != "test1" {
		t.Error("TestFile should store name")
	}
	if testFile.Category != "alpha" {
		t.Error("TestFile should store category")
	}
	if testFile.Constraint != "test.constraint" {
		t.Error("TestFile should store constraint file")
	}
	if testFile.Facts != "test.facts" {
		t.Error("TestFile should store facts file")
	}
}

func TestDiscoverTestsReturnsValidStructure(t *testing.T) {
	tests := DiscoverTests()

	for _, test := range tests {
		if test.Name == "" && test.Constraint != "" {
			t.Error("If constraint is set, name should not be empty")
		}
		if test.Constraint != "" && test.Facts == "" {
			t.Error("If constraint is set, facts should also be set")
		}
	}
}

func TestExecuteTestHandlesNilGracefully(t *testing.T) {
	// Test with minimal TestFile
	testFile := TestFile{
		Name:       "minimal",
		Category:   "test",
		Constraint: "",
		Facts:      "",
	}

	// Should not panic
	result := ExecuteTest(testFile, false)

	if result.Name != "minimal" {
		t.Error("Should preserve test name")
	}
}

func TestGetErrorTestsReturnsMap(t *testing.T) {
	errorTests := GetErrorTests()

	if errorTests == nil {
		t.Fatal("GetErrorTests should not return nil")
	}

	// Should be a valid map
	_, ok := errorTests["error_args_test"]
	if !ok {
		t.Error("Map should contain error_args_test key")
	}
}

func TestPrintHeaderFormat(t *testing.T) {
	var buf bytes.Buffer
	PrintHeader(&buf)

	output := buf.String()
	lines := strings.Split(output, "\n")

	// Should have multiple lines
	if len(lines) < 4 {
		t.Error("Header should have multiple lines")
	}

	// First line should be a separator
	if !strings.Contains(lines[0], "═") {
		t.Error("First line should contain separator")
	}
}

func TestRunTestsProducesOutput(t *testing.T) {
	var stdout bytes.Buffer

	result := RunTests(&stdout)

	output := stdout.String()

	// Should produce some output
	if len(output) == 0 {
		t.Error("RunTests should produce output")
	}

	// Should have results
	if result.Results == nil {
		t.Error("Results should not be nil")
	}
}

func TestRunExitCodeLogic(t *testing.T) {
	var stdout, stderr bytes.Buffer

	exitCode := Run(&stdout, &stderr)

	// Verify exit code is within expected range
	if exitCode != 0 && exitCode != 1 {
		t.Errorf("Exit code should be 0 (success) or 1 (failure), got %d", exitCode)
	}

	output := stdout.String()

	// If exit code is 0, should say all tests passed
	if exitCode == 0 && !strings.Contains(output, "TOUS LES TESTS SONT PASSÉS") {
		t.Error("When exit code is 0, should indicate all tests passed")
	}
}
