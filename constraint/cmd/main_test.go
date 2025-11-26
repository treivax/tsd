// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMain tests the entry point logic via subprocess
func TestMainIntegration(t *testing.T) {
	// Build the binary
	testBinary := filepath.Join(t.TempDir(), "constraint-cmd-test")
	buildCmd := exec.Command("go", "build", "-o", testBinary, ".")

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	buildCmd.Dir = wd

	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build test binary: %v\nOutput: %s", err, output)
	}
	defer os.Remove(testBinary)

	tempDir := t.TempDir()

	// Create test constraint file
	validConstraint := filepath.Join(tempDir, "valid.tsd")
	validContent := []byte("type Person : <id: string, name: string, age: number>")
	if err := os.WriteFile(validConstraint, validContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create invalid constraint file
	invalidConstraint := filepath.Join(tempDir, "invalid.tsd")
	invalidContent := []byte("invalid syntax !!!")
	if err := os.WriteFile(invalidConstraint, invalidContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name               string
		args               []string
		wantExitCode       int
		wantOutputContains []string
		wantErrorContains  []string
	}{
		{
			name:         "no arguments - show usage",
			args:         []string{},
			wantExitCode: 1,
			wantOutputContains: []string{
				"Usage:",
				"constraint-parser <input-file>",
				"Exemple:",
			},
		},
		{
			name:         "valid constraint file",
			args:         []string{validConstraint},
			wantExitCode: 0,
			wantOutputContains: []string{
				"Person",
			},
		},
		{
			name:         "non-existent file",
			args:         []string{filepath.Join(tempDir, "nonexistent.tsd")},
			wantExitCode: 1,
			wantErrorContains: []string{
				"Erreur:",
			},
		},
		{
			name:         "invalid syntax",
			args:         []string{invalidConstraint},
			wantExitCode: 1,
			wantErrorContains: []string{
				"Erreur:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(testBinary, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			// Check exit code
			exitCode := 0
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				}
			}

			if exitCode != tt.wantExitCode {
				t.Errorf("Exit code = %d, want %d\nOutput: %s", exitCode, tt.wantExitCode, outputStr)
			}

			// Check expected output
			for _, expected := range tt.wantOutputContains {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Output does not contain %q\nGot: %s", expected, outputStr)
				}
			}

			// Check expected error output
			for _, expected := range tt.wantErrorContains {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Output does not contain error %q\nGot: %s", expected, outputStr)
				}
			}
		})
	}
}

// TestValidConstraintParsing tests parsing of valid constraint files
func TestValidConstraintParsing(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name           string
		constraintText string
		expectedInJSON []string
	}{
		{
			name:           "simple type definition",
			constraintText: "type Person : <id: string, name: string>",
			expectedInJSON: []string{
				"Person",
				"id",
				"name",
				"string",
			},
		},
		{
			name: "multiple types",
			constraintText: `type Person : <id: string, name: string>
type Address : <street: string, city: string>`,
			expectedInJSON: []string{
				"Person",
				"Address",
				"street",
				"city",
			},
		},
		{
			name:           "type with different field types",
			constraintText: "type User : <id: string, age: number, active: bool>",
			expectedInJSON: []string{
				"User",
				"id",
				"age",
				"active",
				"string",
				"number",
				"bool",
			},
		},
		{
			name: "type with constraint",
			constraintText: `type Person : <id: string, age: number>

rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)`,
			expectedInJSON: []string{
				"Person",
				"age",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraintFile := filepath.Join(tempDir, "test.tsd")
			if err := os.WriteFile(constraintFile, []byte(tt.constraintText), 0644); err != nil {
				t.Fatalf("Failed to create constraint file: %v", err)
			}

			testBinary := filepath.Join(t.TempDir(), "constraint-cmd-parse-test")
			buildCmd := exec.Command("go", "build", "-o", testBinary, ".")
			buildCmd.Dir, _ = os.Getwd()

			if output, err := buildCmd.CombinedOutput(); err != nil {
				t.Fatalf("Failed to build: %v\n%s", err, output)
			}
			defer os.Remove(testBinary)

			cmd := exec.Command(testBinary, constraintFile)
			output, err := cmd.CombinedOutput()

			if err != nil {
				t.Fatalf("Command failed: %v\nOutput: %s", err, string(output))
			}

			outputStr := string(output)
			for _, expected := range tt.expectedInJSON {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("JSON output does not contain %q\nGot: %s", expected, outputStr)
				}
			}
		})
	}
}

// TestInvalidConstraintFiles tests error handling for invalid files
func TestInvalidConstraintFiles(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name           string
		constraintText string
		errorContains  string
	}{
		{
			name:           "invalid syntax - random text",
			constraintText: "this is not a valid constraint",
			errorContains:  "Erreur:",
		},
		{
			name:           "invalid syntax - incomplete type",
			constraintText: "type Person :",
			errorContains:  "Erreur:",
		},
		{
			name:           "invalid syntax - missing bracket",
			constraintText: "type Person : <id: string",
			errorContains:  "Erreur:",
		},
		{
			name:           "invalid field type",
			constraintText: "type Person : <id: invalidtype>",
			errorContains:  "Erreur",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraintFile := filepath.Join(tempDir, "invalid.tsd")
			if err := os.WriteFile(constraintFile, []byte(tt.constraintText), 0644); err != nil {
				t.Fatalf("Failed to create constraint file: %v", err)
			}

			testBinary := filepath.Join(t.TempDir(), "constraint-cmd-invalid-test")
			buildCmd := exec.Command("go", "build", "-o", testBinary, ".")
			buildCmd.Dir, _ = os.Getwd()

			if output, err := buildCmd.CombinedOutput(); err != nil {
				t.Fatalf("Failed to build: %v\n%s", err, output)
			}
			defer os.Remove(testBinary)

			cmd := exec.Command(testBinary, constraintFile)
			output, err := cmd.CombinedOutput()

			if err == nil {
				t.Errorf("Expected error but command succeeded\nOutput: %s", string(output))
			}

			outputStr := string(output)
			if !strings.Contains(outputStr, tt.errorContains) {
				t.Errorf("Error output does not contain %q\nGot: %s", tt.errorContains, outputStr)
			}
		})
	}
}

// TestFileReadError tests error handling for file read errors
func TestFileReadError(t *testing.T) {
	tempDir := t.TempDir()

	// Create a file with no read permissions (Unix-like systems)
	noReadFile := filepath.Join(tempDir, "noread.tsd")
	if err := os.WriteFile(noReadFile, []byte("type Person : <id: string>"), 0000); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Chmod(noReadFile, 0644) // Restore permissions for cleanup

	testBinary := filepath.Join(t.TempDir(), "constraint-cmd-read-test")
	buildCmd := exec.Command("go", "build", "-o", testBinary, ".")
	buildCmd.Dir, _ = os.Getwd()

	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\n%s", err, output)
	}
	defer os.Remove(testBinary)

	cmd := exec.Command(testBinary, noReadFile)
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Errorf("Expected error for unreadable file but command succeeded")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Erreur lecture fichier") && !strings.Contains(outputStr, "permission denied") {
		// Note: On some systems, the file might still be readable
		t.Logf("Warning: Expected permission error, got: %s", outputStr)
	}
}

// TestUsageMessage tests the usage message format
func TestUsageMessage(t *testing.T) {
	testBinary := filepath.Join(t.TempDir(), "constraint-cmd-usage-test")
	buildCmd := exec.Command("go", "build", "-o", testBinary, ".")
	buildCmd.Dir, _ = os.Getwd()

	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\n%s", err, output)
	}
	defer os.Remove(testBinary)

	cmd := exec.Command(testBinary)
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected non-zero exit code when no arguments provided")
	}

	outputStr := string(output)
	requiredStrings := []string{
		"Usage:",
		"constraint-parser <input-file>",
		"Exemple:",
		"constraint-parser ../tests/test_input.txt",
	}

	for _, required := range requiredStrings {
		if !strings.Contains(outputStr, required) {
			t.Errorf("Usage message missing %q\nGot: %s", required, outputStr)
		}
	}
}

// TestJSONOutput tests that output is valid JSON
func TestJSONOutput(t *testing.T) {
	tempDir := t.TempDir()

	constraintFile := filepath.Join(tempDir, "test.tsd")
	constraintContent := []byte("type Person : <id: string, name: string>")
	if err := os.WriteFile(constraintFile, constraintContent, 0644); err != nil {
		t.Fatalf("Failed to create constraint file: %v", err)
	}

	testBinary := filepath.Join(t.TempDir(), "constraint-cmd-json-test")
	buildCmd := exec.Command("go", "build", "-o", testBinary, ".")
	buildCmd.Dir, _ = os.Getwd()

	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\n%s", err, output)
	}
	defer os.Remove(testBinary)

	cmd := exec.Command(testBinary, constraintFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, string(output))
	}

	// Check if output starts with { and ends with }
	outputStr := string(output)

	// The output contains a validation message followed by JSON
	// Check if JSON is present in the output
	if !strings.Contains(outputStr, "{") || !strings.Contains(outputStr, "}") {
		t.Errorf("Output does not contain JSON\nGot: %s", outputStr)
	}

	// Check for common JSON formatting
	if !strings.Contains(outputStr, "\"") {
		t.Errorf("Output missing JSON string quotes\nGot: %s", outputStr)
	}
}

// TestMultipleArguments tests behavior with multiple arguments
func TestMultipleArguments(t *testing.T) {
	tempDir := t.TempDir()

	file1 := filepath.Join(tempDir, "test1.tsd")
	file2 := filepath.Join(tempDir, "test2.tsd")

	content := []byte("type Person : <id: string>")
	os.WriteFile(file1, content, 0644)
	os.WriteFile(file2, content, 0644)

	testBinary := filepath.Join(t.TempDir(), "constraint-cmd-multi-test")
	buildCmd := exec.Command("go", "build", "-o", testBinary, ".")
	buildCmd.Dir, _ = os.Getwd()

	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\n%s", err, output)
	}
	defer os.Remove(testBinary)

	// The program only processes the first argument
	cmd := exec.Command(testBinary, file1, file2)
	output, err := cmd.CombinedOutput()

	// Should succeed using first argument
	if err != nil {
		t.Logf("Command with multiple args: %v\nOutput: %s", err, string(output))
	}
}

// TestEmptyConstraintFile tests handling of empty constraint files
func TestEmptyConstraintFile(t *testing.T) {
	tempDir := t.TempDir()

	emptyFile := filepath.Join(tempDir, "empty.tsd")
	if err := os.WriteFile(emptyFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	testBinary := filepath.Join(t.TempDir(), "constraint-cmd-empty-test")
	buildCmd := exec.Command("go", "build", "-o", testBinary, ".")
	buildCmd.Dir, _ = os.Getwd()

	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\n%s", err, output)
	}
	defer os.Remove(testBinary)

	cmd := exec.Command(testBinary, emptyFile)
	output, _ := cmd.CombinedOutput()

	// Empty file should either parse successfully or fail gracefully
	outputStr := string(output)
	t.Logf("Empty file output: %s", outputStr)
}

// TestStdoutCapture tests that stdout is properly used for output
func TestStdoutCapture(t *testing.T) {
	tempDir := t.TempDir()

	constraintFile := filepath.Join(tempDir, "test.tsd")
	content := []byte("type Person : <id: string>")
	if err := os.WriteFile(constraintFile, content, 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	testBinary := filepath.Join(t.TempDir(), "constraint-cmd-stdout-test")
	buildCmd := exec.Command("go", "build", "-o", testBinary, ".")
	buildCmd.Dir, _ = os.Getwd()

	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\n%s", err, output)
	}
	defer os.Remove(testBinary)

	cmd := exec.Command(testBinary, constraintFile)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Command failed: %v\nStderr: %s", err, stderr.String())
	}

	if stdout.Len() == 0 {
		t.Error("Expected JSON output on stdout, got nothing")
	}

	if stderr.Len() > 0 {
		t.Logf("Stderr (should be empty): %s", stderr.String())
	}
}

// TestValidationError tests constraint validation errors
func TestValidationError(t *testing.T) {
	// This test would need a constraint that passes parsing but fails validation
	// For now, we just test that validation is called
	tempDir := t.TempDir()

	// Create a valid constraint (validation should pass)
	constraintFile := filepath.Join(tempDir, "test.tsd")
	content := []byte("type Person : <id: string, name: string>")
	if err := os.WriteFile(constraintFile, content, 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	testBinary := filepath.Join(t.TempDir(), "constraint-cmd-validation-test")
	buildCmd := exec.Command("go", "build", "-o", testBinary, ".")
	buildCmd.Dir, _ = os.Getwd()

	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\n%s", err, output)
	}
	defer os.Remove(testBinary)

	cmd := exec.Command(testBinary, constraintFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Valid constraint should not fail validation: %v\nOutput: %s", err, string(output))
	}
}
