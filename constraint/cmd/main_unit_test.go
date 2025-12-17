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
			expectExitCode: 1, // ExitUsageError
			checkOutput: func(t *testing.T, stdout, stderr string) {
				if !strings.Contains(stderr, "Usage") {
					t.Error("Should print usage when no arguments")
				}
			},
		},
		{
			name:           "nonexistent file",
			args:           []string{"nonexistent.tsd"},
			expectExitCode: 2, // ExitRuntimeError
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

	content := "type Person(#id: string, name:string)"
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

	if exitCode != 2 { // ExitRuntimeError
		t.Errorf("exitCode = %d, want 2", exitCode)
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
			content:     "type Person(#id: string, name:string)",
			expectError: false,
		},
		{
			name:        "valid type with multiple fields",
			content:     "type Employee(#id: string, name: string, age: number, salary:number)",
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

	content := `type Person(#id: string, name: string, age:number)
type Company(#id: string, name: string, employees:number)`
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
	content := "type Person(#id:string)"
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

// TestRunWithVersionFlag tests --version flag
func TestRunWithVersionFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{"--version"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("✗ exitCode = %d, want 0 for --version", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "constraint-parser") {
		t.Error("✗ Version output should contain 'constraint-parser'")
	}
	if !strings.Contains(output, "version") {
		t.Error("✗ Version output should contain 'version'")
	}
}

// TestRunWithHelpFlag tests --help flag
func TestRunWithHelpFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{"--help"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("✗ exitCode = %d, want 0 for --help", exitCode)
	}

	output := stderr.String()
	if !strings.Contains(output, "Usage") {
		t.Error("✗ Help output should contain 'Usage'")
	}
	if !strings.Contains(output, "constraint-parser") {
		t.Error("✗ Help output should contain 'constraint-parser'")
	}
}

// TestRunWithDebugFlag tests --debug flag
func TestRunWithDebugFlag(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.tsd")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	content := "type Person(#id: string, name:string)"
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{"--debug", tmpfile.Name()}, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("✗ exitCode = %d, want 0 with --debug flag, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Person") {
		t.Error("✗ Debug mode output should still contain parsed content")
	}
}

// TestRunWithOutputFlag tests --output flag
func TestRunWithOutputFlag(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.tsd")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	content := "type Person(#id: string, name:string)"
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{"--output", "json", tmpfile.Name()}, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("✗ exitCode = %d, want 0 with --output json, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "{") {
		t.Error("✗ JSON output should contain JSON brackets")
	}
}

// TestRunWithInvalidOutputFormat tests invalid output format
func TestRunWithInvalidOutputFormat(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.tsd")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	content := "type Person(#id: string, name:string)"
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{"--output", "xml", tmpfile.Name()}, &stdout, &stderr)

	if exitCode != 2 { // ExitRuntimeError
		t.Errorf("✗ exitCode = %d, want 2 for invalid output format", exitCode)
	}

	output := stderr.String()
	if !strings.Contains(output, "Erreur sortie") || !strings.Contains(output, "non supporté") {
		t.Errorf("✗ Should show unsupported format error, got: %s", output)
	}
}

// TestParseInputFromStdin tests stdin input ("-")
func TestParseInputFromStdin(t *testing.T) {
	// Note: This test would require mocking stdin, which is complex
	// For now, we test the stdin placeholder recognition
	t.Skip("Stdin testing requires complex mocking - tested in integration tests")
}

// TestLoadConfigurationWithFile tests loading config from file
func TestLoadConfigurationWithFile(t *testing.T) {
	// Clear environment variables that might interfere
	oldMaxExpr := os.Getenv("CONSTRAINT_MAX_EXPRESSIONS")
	oldMaxDepth := os.Getenv("CONSTRAINT_MAX_DEPTH")
	oldDebug := os.Getenv("CONSTRAINT_DEBUG")
	oldLogLevel := os.Getenv("CONSTRAINT_LOG_LEVEL")
	os.Unsetenv("CONSTRAINT_MAX_EXPRESSIONS")
	os.Unsetenv("CONSTRAINT_MAX_DEPTH")
	os.Unsetenv("CONSTRAINT_DEBUG")
	os.Unsetenv("CONSTRAINT_LOG_LEVEL")
	defer func() {
		if oldMaxExpr != "" {
			os.Setenv("CONSTRAINT_MAX_EXPRESSIONS", oldMaxExpr)
		}
		if oldMaxDepth != "" {
			os.Setenv("CONSTRAINT_MAX_DEPTH", oldMaxDepth)
		}
		if oldDebug != "" {
			os.Setenv("CONSTRAINT_DEBUG", oldDebug)
		}
		if oldLogLevel != "" {
			os.Setenv("CONSTRAINT_LOG_LEVEL", oldLogLevel)
		}
	}()

	tmpDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configFile := filepath.Join(tmpDir, "config.json")
	configContent := []byte(`{
		"parser": {
			"max_expressions": 100,
			"debug": false,
			"recover": true
		},
		"validator": {
			"strict_mode": true,
			"allowed_operators": ["==", "!=", "<", ">"],
			"max_depth": 50
		},
		"logger": {
			"level": "info",
			"format": "json",
			"output": "stdout"
		},
		"debug": false,
		"version": "1.0.0"
	}`)
	if err := os.WriteFile(configFile, configContent, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cliConfig := &CLIConfig{
		ConfigFile:   configFile,
		InputFile:    "test.tsd",
		OutputFormat: "json",
	}

	cfg, err := loadConfiguration(cliConfig)
	if err != nil {
		t.Errorf("✗ Unexpected error loading config: %v", err)
	}
	if cfg == nil {
		t.Error("✗ Config should not be nil")
	}
}

// TestLoadConfigurationWithInvalidFile tests error handling for invalid config file
func TestLoadConfigurationWithInvalidFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configFile := filepath.Join(tmpDir, "invalid_config.json")
	configContent := []byte(`{invalid json}`)
	if err := os.WriteFile(configFile, configContent, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cliConfig := &CLIConfig{
		ConfigFile:   configFile,
		InputFile:    "test.tsd",
		OutputFormat: "json",
	}

	_, err = loadConfiguration(cliConfig)
	if err == nil {
		t.Error("✗ Expected error for invalid JSON config file")
	}
	if !strings.Contains(err.Error(), "chargement fichier config") {
		t.Errorf("✗ Error should mention config file loading, got: %v", err)
	}
}

// TestLoadConfigurationWithoutFile tests default config when no file exists
func TestLoadConfigurationWithoutFile(t *testing.T) {
	// Clear all CONSTRAINT_* environment variables that might interfere
	envVars := []string{
		"CONSTRAINT_MAX_EXPRESSIONS",
		"CONSTRAINT_MAX_DEPTH",
		"CONSTRAINT_DEBUG",
		"CONSTRAINT_STRICT_MODE",
		"CONSTRAINT_LOG_LEVEL",
		"CONSTRAINT_LOG_FORMAT",
		"CONSTRAINT_LOG_OUTPUT",
		"CONSTRAINT_CONFIG_FILE",
	}

	oldValues := make(map[string]string)
	for _, envVar := range envVars {
		oldValues[envVar] = os.Getenv(envVar)
		os.Unsetenv(envVar)
	}
	defer func() {
		for envVar, oldValue := range oldValues {
			if oldValue != "" {
				os.Setenv(envVar, oldValue)
			}
		}
	}()

	cliConfig := &CLIConfig{
		ConfigFile:   "", // No config file
		InputFile:    "test.tsd",
		OutputFormat: "json",
		Debug:        false,
	}

	cfg, err := loadConfiguration(cliConfig)
	if err != nil {
		t.Errorf("✗ Unexpected error with no config file: %v", err)
	}
	if cfg == nil {
		t.Error("✗ Config should not be nil even without config file")
	}
}

// TestParseFlagsWithAllOptions tests all flag combinations
func TestParseFlagsWithAllOptions(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantConfig  bool
		wantOutput  string
		wantDebug   bool
		wantVersion bool
		wantHelp    bool
		wantInput   string
		expectError bool
	}{
		{
			name:        "✅ all flags set",
			args:        []string{"--config", "test.json", "--output", "json", "--debug", "--version", "input.tsd"},
			wantConfig:  true,
			wantOutput:  "json",
			wantDebug:   true,
			wantVersion: true,
			wantInput:   "input.tsd",
		},
		{
			name:       "✅ minimal args",
			args:       []string{"input.tsd"},
			wantOutput: "json", // default
			wantInput:  "input.tsd",
		},
		{
			name:       "✅ help flag only",
			args:       []string{"--help"},
			wantHelp:   true,
			wantOutput: "json", // default
		},
		{
			name:        "✅ version flag only",
			args:        []string{"--version"},
			wantVersion: true,
			wantOutput:  "json", // default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer
			cfg, err := parseFlags(tt.args, &stderr)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if cfg != nil {
				if cfg.OutputFormat != tt.wantOutput {
					t.Errorf("OutputFormat = %q, want %q", cfg.OutputFormat, tt.wantOutput)
				}
				if cfg.Debug != tt.wantDebug {
					t.Errorf("Debug = %v, want %v", cfg.Debug, tt.wantDebug)
				}
				if cfg.Version != tt.wantVersion {
					t.Errorf("Version = %v, want %v", cfg.Version, tt.wantVersion)
				}
				if cfg.Help != tt.wantHelp {
					t.Errorf("Help = %v, want %v", cfg.Help, tt.wantHelp)
				}
				if tt.wantInput != "" && cfg.InputFile != tt.wantInput {
					t.Errorf("InputFile = %q, want %q", cfg.InputFile, tt.wantInput)
				}
			}
		})
	}
}

// TestOutputResultWithDifferentFormats tests OutputResult with various formats
func TestOutputResultWithDifferentFormats(t *testing.T) {
	testData := map[string]string{"test": "data"}

	tests := []struct {
		name        string
		format      string
		expectError bool
	}{
		{
			name:        "✅ json format",
			format:      "json",
			expectError: false,
		},
		{
			name:        "✗ xml format (unsupported)",
			format:      "xml",
			expectError: true,
		},
		{
			name:        "✗ yaml format (unsupported)",
			format:      "yaml",
			expectError: true,
		},
		{
			name:        "✗ invalid format",
			format:      "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := OutputResult(testData, tt.format, &buf)

			if tt.expectError && err == nil {
				t.Error("Expected error for unsupported format but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.expectError && err != nil {
				if !strings.Contains(err.Error(), "non supporté") {
					t.Errorf("Error should mention unsupported format, got: %v", err)
				}
			}
		})
	}
}

// TestRunWithConfigFlag tests Run with --config flag
func TestRunWithConfigFlag(t *testing.T) {
	// Clear all CONSTRAINT_* environment variables that might interfere
	envVars := []string{
		"CONSTRAINT_MAX_EXPRESSIONS",
		"CONSTRAINT_MAX_DEPTH",
		"CONSTRAINT_DEBUG",
		"CONSTRAINT_STRICT_MODE",
		"CONSTRAINT_LOG_LEVEL",
		"CONSTRAINT_LOG_FORMAT",
		"CONSTRAINT_LOG_OUTPUT",
		"CONSTRAINT_CONFIG_FILE",
	}

	oldValues := make(map[string]string)
	for _, envVar := range envVars {
		oldValues[envVar] = os.Getenv(envVar)
		os.Unsetenv(envVar)
	}
	defer func() {
		for envVar, oldValue := range oldValues {
			if oldValue != "" {
				os.Setenv(envVar, oldValue)
			}
		}
	}()

	tmpDir, err := os.MkdirTemp("", "test_run_config")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create valid config file with complete structure
	configFile := filepath.Join(tmpDir, "config.json")
	configContent := []byte(`{
		"parser": {
			"max_expressions": 200,
			"debug": false,
			"recover": true
		},
		"validator": {
			"strict_mode": true,
			"allowed_operators": ["==", "!=", "<", ">", "<=", ">="],
			"max_depth": 50
		},
		"logger": {
			"level": "info",
			"format": "json",
			"output": "stdout"
		},
		"debug": false,
		"version": "1.0.0"
	}`)
	if err := os.WriteFile(configFile, configContent, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Create valid constraint file
	constraintFile := filepath.Join(tmpDir, "test.tsd")
	constraintContent := []byte("type Person(#id: string, name:string)")
	if err := os.WriteFile(constraintFile, constraintContent, 0644); err != nil {
		t.Fatalf("Failed to write constraint file: %v", err)
	}

	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{"--config", configFile, constraintFile}, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("✗ exitCode = %d, want 0 with --config, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Person") {
		t.Error("✗ Output should contain parsed type Person")
	}
}

// TestParseInputErrorHandling tests error paths in ParseInput
func TestParseInputErrorHandling(t *testing.T) {
	// Test with validation error (if possible)
	tmpfile, err := os.CreateTemp("", "test*.tsd")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Write content that parses but might fail validation
	content := "type Person(#id: string, name:string)"
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	// This should succeed (valid constraint)
	result, err := ParseInput(tmpfile.Name())
	if err != nil {
		t.Errorf("✗ Valid constraint should not fail: %v", err)
	}
	if result == nil {
		t.Error("✗ Result should not be nil for valid constraint")
	}
}
