// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectExitCode int
		checkOutput    func(*testing.T, string, string)
	}{
		{
			name:           "no arguments",
			args:           []string{},
			expectExitCode: 1,
			checkOutput: func(t *testing.T, stdout, stderr string) {
				if !strings.Contains(stderr, "Usage") {
					t.Error("Should print usage when no arguments")
				}
			},
		},
		{
			name:           "nonexistent file",
			args:           []string{"nonexistent.tsd"},
			expectExitCode: 1,
			checkOutput: func(t *testing.T, stdout, stderr string) {
				if !strings.Contains(stderr, "Erreur") {
					t.Error("Should print error for nonexistent file")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer

			exitCode := Run(tt.args, &stdout, &stderr)

			if exitCode != tt.expectExitCode {
				t.Errorf("exitCode = %d, want %d", exitCode, tt.expectExitCode)
			}

			if tt.checkOutput != nil {
				tt.checkOutput(t, stdout.String(), stderr.String())
			}
		})
	}
}

func TestRunWithValidFile(t *testing.T) {
	// Create temporary constraint file
	tmpfile, err := os.CreateTemp("", "test*.tsd")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	content := "type Person(id: string, name:string)"
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{tmpfile.Name()}, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("exitCode = %d, want 0, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Person") {
		t.Error("Output should contain parsed type Person")
	}
	if !strings.Contains(output, "id") {
		t.Error("Output should contain field id")
	}
	if !strings.Contains(output, "name") {
		t.Error("Output should contain field name")
	}
}

func TestRunWithInvalidSyntax(t *testing.T) {
	// Create temporary constraint file with invalid syntax
	tmpfile, err := os.CreateTemp("", "test*.tsd")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	content := "invalid syntax here @#$%"
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{tmpfile.Name()}, &stdout, &stderr)

	if exitCode != 1 {
		t.Errorf("exitCode = %d, want 1", exitCode)
	}

	output := stderr.String()
	if !strings.Contains(output, "Erreur") {
		t.Error("Should print error for invalid syntax")
	}
}

func TestParseFile(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name:        "valid type definition",
			content:     "type Person(id: string, name:string)",
			expectError: false,
		},
		{
			name:        "valid type with multiple fields",
			content:     "type Employee(id: string, name: string, age: number, salary:number)",
			expectError: false,
		},
		{
			name:        "invalid syntax",
			content:     "invalid @#$%",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpfile, err := os.CreateTemp("", "test*.tsd")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.content)); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpfile.Close()

			result, err := ParseFile(tmpfile.Name())

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectError && result == nil {
				t.Error("Result should not be nil for valid input")
			}
		})
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, err := ParseFile("nonexistent_file_xyz.tsd")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
	if !strings.Contains(err.Error(), "lecture fichier") {
		t.Errorf("Error should mention file reading, got: %v", err)
	}
}

func TestOutputJSON(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		errMsg string
	}{
		{
			name: "simple map",
			input: map[string]string{
				"key": "value",
			},
		},
		{
			name: "nested structure",
			input: map[string]interface{}{
				"name": "Person",
				"fields": []string{
					"id",
					"name",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := OutputJSON(tt.input, &buf)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			output := buf.String()
			if output == "" {
				t.Error("Output should not be empty")
			}

			// Check if output is valid-looking JSON
			if !strings.Contains(output, "{") && !strings.Contains(output, "[") {
				t.Error("Output should contain JSON brackets")
			}
		})
	}
}

func TestOutputJSONWithInvalidInput(t *testing.T) {
	var buf bytes.Buffer

	// Channels cannot be marshaled to JSON
	invalidInput := make(chan int)

	err := OutputJSON(invalidInput, &buf)

	if err == nil {
		t.Error("Expected error for invalid JSON input")
	}
}

func TestPrintHelp(t *testing.T) {
	var buf bytes.Buffer

	PrintHelp(&buf)

	output := buf.String()

	if !strings.Contains(output, "Usage") {
		t.Error("Help should contain Usage")
	}
	if !strings.Contains(output, "Exemple") {
		t.Error("Help should contain Exemple")
	}
	if !strings.Contains(output, "constraint-parser") {
		t.Error("Help should mention constraint-parser")
	}
}

func TestRunIntegrationWithRealConstraintFile(t *testing.T) {
	// Try to find an existing constraint file in the test directories
	testFiles := []string{
		"../../test/coverage/alpha/simple_alpha_test.tsd",
		"../../constraint/test/integration/basic_types.tsd",
		"../test/integration/basic_types.tsd",
	}

	var validTestFile string
	for _, file := range testFiles {
		if _, err := os.Stat(file); err == nil {
			validTestFile = file
			break
		}
	}

	if validTestFile == "" {
		t.Skip("No test constraint files found, skipping integration test")
	}

	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{validTestFile}, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("exitCode = %d, want 0 for valid constraint file, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if output == "" {
		t.Error("Output should not be empty for valid constraint file")
	}
}

func TestRunWithComplexConstraintFile(t *testing.T) {
	// Create a more complex constraint file
	tmpfile, err := os.CreateTemp("", "test*.tsd")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	content := `type Person(id: string, name: string, age:number)
type Company(id: string, name: string, employees:number)`
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{tmpfile.Name()}, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("exitCode = %d, want 0, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Person") {
		t.Error("Output should contain type Person")
	}
	if !strings.Contains(output, "Company") {
		t.Error("Output should contain type Company")
	}
}

func TestParseFileWithRelativePath(t *testing.T) {
	// Create a temp directory and file
	tmpDir, err := os.MkdirTemp("", "test_dir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpfile := filepath.Join(tmpDir, "test.tsd")
	content := "type Person(id:string)"
	if err := os.WriteFile(tmpfile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	result, err := ParseFile(tmpfile)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Error("Result should not be nil")
	}
}

func TestRunWithEmptyFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.tsd")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{tmpfile.Name()}, &stdout, &stderr)

	// Empty file is actually valid (0 types), so should succeed
	if exitCode != 0 {
		t.Errorf("exitCode = %d, want 0 for empty file", exitCode)
	}
}

func TestRunWithWhitespaceOnlyFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.tsd")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	content := "   \n\n   \t\t\n   "
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{tmpfile.Name()}, &stdout, &stderr)

	// Whitespace-only file is valid (0 types), so should succeed
	if exitCode != 0 {
		t.Errorf("exitCode = %d, want 0 for whitespace-only file", exitCode)
	}
}
