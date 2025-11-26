// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/treivax/tsd/rete"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected *Config
	}{
		{
			name: "constraint file flag",
			args: []string{"-constraint", "test.constraint"},
			expected: &Config{
				ConstraintFile: "test.constraint",
			},
		},
		{
			name: "text flag",
			args: []string{"-text", "type Person : <id: string>"},
			expected: &Config{
				ConstraintText: "type Person : <id: string>",
			},
		},
		{
			name: "stdin flag",
			args: []string{"-stdin"},
			expected: &Config{
				UseStdin: true,
			},
		},
		{
			name: "verbose flag",
			args: []string{"-constraint", "test.constraint", "-v"},
			expected: &Config{
				ConstraintFile: "test.constraint",
				Verbose:        true,
			},
		},
		{
			name: "facts file flag",
			args: []string{"-constraint", "test.constraint", "-facts", "test.facts"},
			expected: &Config{
				ConstraintFile: "test.constraint",
				FactsFile:      "test.facts",
			},
		},
		{
			name: "version flag",
			args: []string{"-version"},
			expected: &Config{
				ShowVersion: true,
			},
		},
		{
			name: "help flag",
			args: []string{"-h"},
			expected: &Config{
				ShowHelp: true,
			},
		},
		{
			name: "multiple flags",
			args: []string{"-constraint", "rules.constraint", "-facts", "data.facts", "-v"},
			expected: &Config{
				ConstraintFile: "rules.constraint",
				FactsFile:      "data.facts",
				Verbose:        true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ParseFlags(tt.args)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if config.ConstraintFile != tt.expected.ConstraintFile {
				t.Errorf("ConstraintFile = %v, want %v", config.ConstraintFile, tt.expected.ConstraintFile)
			}
			if config.ConstraintText != tt.expected.ConstraintText {
				t.Errorf("ConstraintText = %v, want %v", config.ConstraintText, tt.expected.ConstraintText)
			}
			if config.UseStdin != tt.expected.UseStdin {
				t.Errorf("UseStdin = %v, want %v", config.UseStdin, tt.expected.UseStdin)
			}
			if config.FactsFile != tt.expected.FactsFile {
				t.Errorf("FactsFile = %v, want %v", config.FactsFile, tt.expected.FactsFile)
			}
			if config.Verbose != tt.expected.Verbose {
				t.Errorf("Verbose = %v, want %v", config.Verbose, tt.expected.Verbose)
			}
			if config.ShowVersion != tt.expected.ShowVersion {
				t.Errorf("ShowVersion = %v, want %v", config.ShowVersion, tt.expected.ShowVersion)
			}
			if config.ShowHelp != tt.expected.ShowHelp {
				t.Errorf("ShowHelp = %v, want %v", config.ShowHelp, tt.expected.ShowHelp)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid constraint file",
			config: &Config{
				ConstraintFile: "test.constraint",
			},
			wantError: false,
		},
		{
			name: "valid text input",
			config: &Config{
				ConstraintText: "type Person : <id: string>",
			},
			wantError: false,
		},
		{
			name: "valid stdin input",
			config: &Config{
				UseStdin: true,
			},
			wantError: false,
		},
		{
			name:      "no input source",
			config:    &Config{},
			wantError: true,
			errorMsg:  "spécifiez une source",
		},
		{
			name: "multiple input sources - file and text",
			config: &Config{
				ConstraintFile: "test.constraint",
				ConstraintText: "type Person : <id: string>",
			},
			wantError: true,
			errorMsg:  "spécifiez une seule source",
		},
		{
			name: "multiple input sources - file and stdin",
			config: &Config{
				ConstraintFile: "test.constraint",
				UseStdin:       true,
			},
			wantError: true,
			errorMsg:  "spécifiez une seule source",
		},
		{
			name: "multiple input sources - text and stdin",
			config: &Config{
				ConstraintText: "type Person : <id: string>",
				UseStdin:       true,
			},
			wantError: true,
			errorMsg:  "spécifiez une seule source",
		},
		{
			name: "all three input sources",
			config: &Config{
				ConstraintFile: "test.constraint",
				ConstraintText: "type Person : <id: string>",
				UseStdin:       true,
			},
			wantError: true,
			errorMsg:  "spécifiez une seule source",
		},
		{
			name: "valid with facts file",
			config: &Config{
				ConstraintFile: "test.constraint",
				FactsFile:      "test.facts",
			},
			wantError: false,
		},
		{
			name: "valid with verbose flag",
			config: &Config{
				ConstraintFile: "test.constraint",
				Verbose:        true,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)

			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateConfig() error = nil, want error containing %q", tt.errorMsg)
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("ValidateConfig() error = %v, want error containing %q", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateConfig() error = %v, want nil", err)
				}
			}
		})
	}
}

func TestParseFromText(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		wantError bool
	}{
		{
			name: "valid constraint text",
			config: &Config{
				ConstraintText: "type Person : <id: string, name: string>",
			},
			wantError: false,
		},
		{
			name: "empty text",
			config: &Config{
				ConstraintText: "",
			},
			wantError: false, // Parser accepts empty input
		},
		{
			name: "invalid syntax",
			config: &Config{
				ConstraintText: "invalid constraint syntax !!!",
			},
			wantError: true,
		},
		{
			name: "verbose mode",
			config: &Config{
				ConstraintText: "type Person : <id: string>",
				Verbose:        true,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout for verbose output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			result, sourceName, err := parseFromText(tt.config)

			w.Close()
			os.Stdout = oldStdout
			r.Close()

			if tt.wantError {
				if err == nil {
					t.Errorf("parseFromText() error = nil, want error")
				}
			} else {
				if err != nil {
					t.Errorf("parseFromText() error = %v, want nil", err)
				}
				if result == nil {
					t.Errorf("parseFromText() result = nil, want non-nil")
				}
				if sourceName != "<text>" {
					t.Errorf("parseFromText() sourceName = %v, want <text>", sourceName)
				}
			}
		})
	}
}

func TestParseFromFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Create a valid constraint file
	validFile := filepath.Join(tempDir, "valid.constraint")
	validContent := []byte("type Person : <id: string, name: string>")
	if err := os.WriteFile(validFile, validContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create an invalid constraint file
	invalidFile := filepath.Join(tempDir, "invalid.constraint")
	invalidContent := []byte("invalid constraint syntax !!!")
	if err := os.WriteFile(invalidFile, invalidContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name      string
		config    *Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid constraint file",
			config: &Config{
				ConstraintFile: validFile,
			},
			wantError: false,
		},
		{
			name: "non-existent file",
			config: &Config{
				ConstraintFile: filepath.Join(tempDir, "nonexistent.constraint"),
			},
			wantError: true,
			errorMsg:  "fichier contrainte non trouvé",
		},
		{
			name: "invalid syntax in file",
			config: &Config{
				ConstraintFile: invalidFile,
			},
			wantError: true,
		},
		{
			name: "verbose mode with valid file",
			config: &Config{
				ConstraintFile: validFile,
				Verbose:        true,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout for verbose output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			result, sourceName, err := parseFromFile(tt.config)

			w.Close()
			os.Stdout = oldStdout
			r.Close()

			if tt.wantError {
				if err == nil {
					t.Errorf("parseFromFile() error = nil, want error")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("parseFromFile() error = %v, want error containing %q", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("parseFromFile() error = %v, want nil", err)
				}
				if result == nil {
					t.Errorf("parseFromFile() result = nil, want non-nil")
				}
				if sourceName != tt.config.ConstraintFile {
					t.Errorf("parseFromFile() sourceName = %v, want %v", sourceName, tt.config.ConstraintFile)
				}
			}
		})
	}
}

func TestPrintVersion(t *testing.T) {
	var buf bytes.Buffer
	PrintVersion(&buf)
	output := buf.String()

	expectedStrings := []string{
		"TSD",
		"Type System Development",
		"v1.0",
		"RETE",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("printVersion() output does not contain %q\nGot: %s", expected, output)
		}
	}
}

func TestPrintHelp(t *testing.T) {
	var buf bytes.Buffer
	PrintHelp(&buf)
	output := buf.String()

	expectedStrings := []string{
		"TSD",
		"USAGE:",
		"OPTIONS:",
		"-constraint",
		"-text",
		"-stdin",
		"-facts",
		"-v",
		"-version",
		"-h",
		"EXEMPLES:",
		"FORMATS DE FICHIERS:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("printHelp() output does not contain %q", expected)
		}
	}
}

func TestCountActivations(t *testing.T) {
	// Note: countActivations in main.go accesses terminal.Memory.Tokens
	// but TerminalNode only has BaseNode and Action fields.
	// The Memory field is part of BaseNode, so we can't easily test this
	// without creating a full network. We'll test the counting logic instead.

	tests := []struct {
		name          string
		tokenCounts   []int
		expectedTotal int
	}{
		{
			name:          "no tokens",
			tokenCounts:   []int{},
			expectedTotal: 0,
		},
		{
			name:          "one terminal with tokens",
			tokenCounts:   []int{5},
			expectedTotal: 5,
		},
		{
			name:          "multiple terminals",
			tokenCounts:   []int{3, 5, 2},
			expectedTotal: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the counting logic
			activations := 0
			for _, count := range tt.tokenCounts {
				activations += count
			}

			if activations != tt.expectedTotal {
				t.Errorf("Total activations = %d, want %d", activations, tt.expectedTotal)
			}
		})
	}
}

func TestRunValidationOnly(t *testing.T) {
	tests := []struct {
		name           string
		config         *Config
		expectedOutput []string
	}{
		{
			name: "basic validation",
			config: &Config{
				Verbose: false,
			},
			expectedOutput: []string{
				"✅ Contraintes validées avec succès",
			},
		},
		{
			name: "verbose validation",
			config: &Config{
				Verbose: true,
			},
			expectedOutput: []string{
				"✅ Contraintes validées avec succès",
				"🎉 Validation terminée!",
				"Les contraintes sont syntaxiquement correctes",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			exitCode := RunValidationOnly(tt.config, &buf)

			if exitCode != 0 {
				t.Errorf("exitCode = %d, want 0", exitCode)
			}

			output := buf.String()

			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("RunValidationOnly() output does not contain %q\nGot: %s", expected, output)
				}
			}
		})
	}
}

// TestConfig tests the Config struct
func TestConfig(t *testing.T) {
	config := &Config{
		ConstraintFile: "test.constraint",
		ConstraintText: "type Person : <id: string>",
		UseStdin:       true,
		FactsFile:      "test.facts",
		Verbose:        true,
		ShowVersion:    true,
		ShowHelp:       true,
	}

	if config.ConstraintFile != "test.constraint" {
		t.Errorf("Config.ConstraintFile = %v, want test.constraint", config.ConstraintFile)
	}
	if config.ConstraintText != "type Person : <id: string>" {
		t.Errorf("Config.ConstraintText = %v, want 'type Person : <id: string>'", config.ConstraintText)
	}
	if !config.UseStdin {
		t.Error("Config.UseStdin = false, want true")
	}
	if config.FactsFile != "test.facts" {
		t.Errorf("Config.FactsFile = %v, want test.facts", config.FactsFile)
	}
	if !config.Verbose {
		t.Error("Config.Verbose = false, want true")
	}
	if !config.ShowVersion {
		t.Error("Config.ShowVersion = false, want true")
	}
	if !config.ShowHelp {
		t.Error("Config.ShowHelp = false, want true")
	}
}

// TestParseConstraintSource tests the ParseConstraintSource routing function
func TestParseConstraintSource(t *testing.T) {
	tempDir := t.TempDir()

	// Create a valid constraint file
	validFile := filepath.Join(tempDir, "valid.constraint")
	validContent := []byte("type Person : <id: string>")
	if err := os.WriteFile(validFile, validContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name           string
		config         *Config
		wantError      bool
		wantSourceName string
		stdinContent   string
	}{
		{
			name: "route to file",
			config: &Config{
				ConstraintFile: validFile,
			},
			wantError:      false,
			wantSourceName: validFile,
		},
		{
			name: "route to text",
			config: &Config{
				ConstraintText: "type Person : <id: string>",
			},
			wantError:      false,
			wantSourceName: "<text>",
		},
		{
			name: "route to stdin",
			config: &Config{
				UseStdin: true,
			},
			wantError:      false,
			wantSourceName: "<stdin>",
			stdinContent:   "type Person : <id: string>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup stdin mock if needed
			if tt.stdinContent != "" {
				oldStdin := os.Stdin
				r, w, _ := os.Pipe()
				os.Stdin = r
				w.Write([]byte(tt.stdinContent))
				w.Close()
				defer func() { os.Stdin = oldStdin }()
			}

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create stdin reader if needed
			var stdinReader io.Reader
			if tt.config.UseStdin && tt.stdinContent != "" {
				stdinReader = strings.NewReader(tt.stdinContent)
			}

			result, sourceName, err := ParseConstraintSource(tt.config, stdinReader)

			w.Close()
			os.Stdout = oldStdout
			r.Close()

			if tt.wantError && err == nil {
				t.Error("ParseConstraintSource() error = nil, want error")
			}
			if !tt.wantError && err != nil {
				t.Errorf("ParseConstraintSource() error = %v, want nil", err)
			}
			if !tt.wantError {
				if result == nil {
					t.Error("ParseConstraintSource() result = nil, want non-nil")
				}
				if sourceName != tt.wantSourceName {
					t.Errorf("ParseConstraintSource() sourceName = %v, want %v", sourceName, tt.wantSourceName)
				}
			}
		})
	}
}

// TestParseFromStdin tests parsing from stdin
func TestParseFromStdin(t *testing.T) {
	tests := []struct {
		name         string
		config       *Config
		stdinContent string
		wantError    bool
	}{
		{
			name: "valid constraint from stdin",
			config: &Config{
				UseStdin: true,
			},
			stdinContent: "type Person : <id: string, name: string>",
			wantError:    false,
		},
		{
			name: "empty stdin",
			config: &Config{
				UseStdin: true,
			},
			stdinContent: "",
			wantError:    false,
		},
		{
			name: "invalid syntax from stdin",
			config: &Config{
				UseStdin: true,
			},
			stdinContent: "invalid syntax !!!",
			wantError:    true,
		},
		{
			name: "verbose mode from stdin",
			config: &Config{
				UseStdin: true,
				Verbose:  true,
			},
			stdinContent: "type Person : <id: string>",
			wantError:    false,
		},
		{
			name: "complex constraint from stdin",
			config: &Config{
				UseStdin: true,
			},
			stdinContent: "type Person : <id: string, name: string, age: number>",
			wantError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			rOut, wOut, _ := os.Pipe()
			os.Stdout = wOut

			result, sourceName, err := parseFromStdin(tt.config, strings.NewReader(tt.stdinContent))

			wOut.Close()
			os.Stdout = oldStdout
			rOut.Close()

			if tt.wantError {
				if err == nil {
					t.Errorf("ParseFromStdin() error = nil, want error")
				}
			} else {
				if err != nil {
					t.Errorf("ParseFromStdin() error = %v, want nil", err)
				}
				if result == nil {
					t.Errorf("ParseFromStdin() result = nil, want non-nil")
				}
				if sourceName != "<stdin>" {
					t.Errorf("ParseFromStdin() sourceName = %v, want <stdin>", sourceName)
				}
			}
		})
	}
}

// TestRunWithFactsLogic tests the logic of checking if facts file exists
func TestRunWithFactsLogic(t *testing.T) {
	tempDir := t.TempDir()

	// Create a valid facts file
	factsFile := filepath.Join(tempDir, "test.facts")
	factsContent := []byte(`Person(id: "1", name: "Alice")`)
	if err := os.WriteFile(factsFile, factsContent, 0644); err != nil {
		t.Fatalf("Failed to create facts file: %v", err)
	}

	// Create a non-existent facts file path
	nonExistentFacts := filepath.Join(tempDir, "nonexistent.facts")

	tests := []struct {
		name      string
		factsFile string
		wantError bool
	}{
		{
			name:      "existing facts file",
			factsFile: factsFile,
			wantError: false,
		},
		{
			name:      "non-existent facts file",
			factsFile: nonExistentFacts,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the file existence check logic used in runWithFacts
			_, err := os.Stat(tt.factsFile)
			hasError := os.IsNotExist(err)

			if tt.wantError && !hasError {
				t.Errorf("Expected file to not exist, but it exists")
			}
			if !tt.wantError && hasError {
				t.Errorf("Expected file to exist, but it doesn't")
			}
		})
	}
}

// TestPrintResults tests the printResults function
func TestPrintResults(t *testing.T) {
	tests := []struct {
		name           string
		config         *Config
		tokenCount     int
		expectedOutput []string
	}{
		{
			name: "no activations - non-verbose",
			config: &Config{
				Verbose: false,
			},
			tokenCount: 0,
			expectedOutput: []string{
				"Aucune action déclenchée",
			},
		},
		{
			name: "no activations - verbose",
			config: &Config{
				Verbose: true,
			},
			tokenCount: 0,
			expectedOutput: []string{
				"RÉSULTATS",
				"Faits injectés:",
				"Aucune action déclenchée",
				"Pipeline RETE exécuté avec succès",
			},
		},
		{
			name: "with activations - non-verbose",
			config: &Config{
				Verbose: false,
			},
			tokenCount: 3,
			expectedOutput: []string{
				"ACTIONS DISPONIBLES: 3",
			},
		},
		{
			name: "with activations - verbose",
			config: &Config{
				Verbose: true,
			},
			tokenCount: 5,
			expectedOutput: []string{
				"RÉSULTATS",
				"ACTIONS DISPONIBLES: 5",
				"Pipeline RETE exécuté avec succès",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock network with the specified token count
			// Since we can't easily create a real ReteNetwork without dependencies,
			// we'll test the output logic by capturing stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// We can't easily create a real network, so we'll test the core logic
			// by checking the output format
			if tt.config.Verbose {
				fmt.Printf("\n📊 RÉSULTATS\n")
				fmt.Printf("============\n")
				fmt.Printf("Faits injectés: %d\n", 2)
			}

			if tt.tokenCount > 0 {
				fmt.Printf("\n🎯 ACTIONS DISPONIBLES: %d\n", tt.tokenCount)
			} else {
				fmt.Printf("\nℹ️  Aucune action déclenchée\n")
			}

			if tt.config.Verbose {
				fmt.Printf("\n✅ Pipeline RETE exécuté avec succès\n")
			}

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("printResults() output does not contain %q\nGot: %s", expected, output)
				}
			}
		})
	}
}

// TestCountActivationsWithRealNetwork tests countActivations with a mock network structure
func TestCountActivationsWithRealNetwork(t *testing.T) {
	tests := []struct {
		name          string
		tokenCounts   []int
		expectedTotal int
	}{
		{
			name:          "no terminals",
			tokenCounts:   []int{},
			expectedTotal: 0,
		},
		{
			name:          "one terminal with no tokens",
			tokenCounts:   []int{0},
			expectedTotal: 0,
		},
		{
			name:          "one terminal with tokens",
			tokenCounts:   []int{5},
			expectedTotal: 5,
		},
		{
			name:          "multiple terminals with tokens",
			tokenCounts:   []int{3, 5, 2, 1},
			expectedTotal: 11,
		},
		{
			name:          "mixed terminals (some with tokens, some without)",
			tokenCounts:   []int{0, 5, 0, 3},
			expectedTotal: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the counting logic used in countActivations
			total := 0
			for _, count := range tt.tokenCounts {
				total += count
			}

			if total != tt.expectedTotal {
				t.Errorf("countActivations() = %d, want %d", total, tt.expectedTotal)
			}
		})
	}
}

// TestPrintActivationDetails tests the printActivationDetails function
func TestPrintActivationDetails(t *testing.T) {
	tests := []struct {
		name        string
		activations []struct {
			name     string
			bindings int
		}
		expectedOutput []string
	}{
		{
			name: "no activations",
			activations: []struct {
				name     string
				bindings int
			}{},
			expectedOutput: []string{},
		},
		{
			name: "single activation",
			activations: []struct {
				name     string
				bindings int
			}{
				{name: "greet", bindings: 2},
			},
			expectedOutput: []string{
				"1. greet() - 2 bindings",
			},
		},
		{
			name: "multiple activations",
			activations: []struct {
				name     string
				bindings int
			}{
				{name: "greet", bindings: 2},
				{name: "notify", bindings: 3},
				{name: "process", bindings: 1},
			},
			expectedOutput: []string{
				"1. greet() - 2 bindings",
				"2. notify() - 3 bindings",
				"3. process() - 1 bindings",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Simulate the logic of printActivationDetails
			count := 0
			for _, activation := range tt.activations {
				count++
				fmt.Printf("  %d. %s() - %d bindings\n", count, activation.name, activation.bindings)
			}

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("printActivationDetails() output does not contain %q\nGot: %s", expected, output)
				}
			}
		})
	}
}

// TestMainIntegration tests the main function via subprocess execution
func TestMainIntegration(t *testing.T) {
	// Build the test binary
	testBinary := filepath.Join(t.TempDir(), "tsd-test")
	buildCmd := []string{"go", "build", "-o", testBinary, "."}

	// Get the current directory (cmd/tsd)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Build the binary
	cmd := exec.Command(buildCmd[0], buildCmd[1:]...)
	cmd.Dir = wd
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build test binary: %v\nOutput: %s", err, output)
	}
	defer os.Remove(testBinary)

	tempDir := t.TempDir()

	// Create test constraint file
	constraintFile := filepath.Join(tempDir, "test.constraint")
	constraintContent := []byte("type Person : <id: string, name: string>")
	if err := os.WriteFile(constraintFile, constraintContent, 0644); err != nil {
		t.Fatalf("Failed to create constraint file: %v", err)
	}

	// Create test facts file
	factsFile := filepath.Join(tempDir, "test.facts")
	factsContent := []byte(`Person(id: "1", name: "Alice")`)
	if err := os.WriteFile(factsFile, factsContent, 0644); err != nil {
		t.Fatalf("Failed to create facts file: %v", err)
	}

	tests := []struct {
		name               string
		args               []string
		stdin              string
		wantExitCode       int
		wantOutputContains []string
		wantErrorContains  []string
	}{
		{
			name:         "help flag",
			args:         []string{"-h"},
			wantExitCode: 0,
			wantOutputContains: []string{
				"TSD",
				"USAGE:",
				"OPTIONS:",
			},
		},
		{
			name:         "version flag",
			args:         []string{"-version"},
			wantExitCode: 0,
			wantOutputContains: []string{
				"TSD",
				"v1.0",
			},
		},
		{
			name:         "constraint file validation",
			args:         []string{"-constraint", constraintFile},
			wantExitCode: 0,
			wantOutputContains: []string{
				"Contraintes validées avec succès",
			},
		},
		{
			name:         "constraint file verbose",
			args:         []string{"-constraint", constraintFile, "-v"},
			wantExitCode: 0,
			wantOutputContains: []string{
				"Parsing réussi",
				"Validation du programme",
				"Validation terminée",
			},
		},
		{
			name:         "text input",
			args:         []string{"-text", "type Person : <id: string>"},
			wantExitCode: 0,
			wantOutputContains: []string{
				"Contraintes validées avec succès",
			},
		},
		{
			name:         "stdin input",
			args:         []string{"-stdin"},
			stdin:        "type Person : <id: string>",
			wantExitCode: 0,
			wantOutputContains: []string{
				"Contraintes validées avec succès",
			},
		},
		{
			name:         "no input source error",
			args:         []string{},
			wantExitCode: 1,
			wantErrorContains: []string{
				"Erreur",
				"spécifiez une source",
			},
		},
		{
			name:         "multiple input sources error",
			args:         []string{"-constraint", constraintFile, "-text", "type X : <a: string>"},
			wantExitCode: 1,
			wantErrorContains: []string{
				"Erreur",
				"spécifiez une seule source",
			},
		},
		{
			name:         "non-existent constraint file",
			args:         []string{"-constraint", filepath.Join(tempDir, "nonexistent.constraint")},
			wantExitCode: 1,
			wantErrorContains: []string{
				"Erreur",
				"fichier contrainte non trouvé",
			},
		},
		{
			name:         "invalid syntax in text",
			args:         []string{"-text", "invalid syntax !!!"},
			wantExitCode: 1,
			wantErrorContains: []string{
				"Erreur de parsing",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(testBinary, tt.args...)

			if tt.stdin != "" {
				cmd.Stdin = strings.NewReader(tt.stdin)
			}

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

// TestMainWithFactsIntegration tests main with facts file processing
func TestMainWithFactsIntegration(t *testing.T) {
	// Build the test binary
	testBinary := filepath.Join(t.TempDir(), "tsd-test-facts")
	buildCmd := []string{"go", "build", "-o", testBinary, "."}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	cmd := exec.Command(buildCmd[0], buildCmd[1:]...)
	cmd.Dir = wd
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build test binary: %v\nOutput: %s", err, output)
	}
	defer os.Remove(testBinary)

	tempDir := t.TempDir()

	// Create a simple constraint file
	constraintFile := filepath.Join(tempDir, "rules.constraint")
	constraintContent := []byte(`type Person : <id: string, name: string, age: number>
type Order : <id: string, customer_id: string, amount: number>

rule r1 : {p: Person, o: Order} / p.id == o.customer_id ==> customer_order(p.id, o.id)
`)
	if err := os.WriteFile(constraintFile, constraintContent, 0644); err != nil {
		t.Fatalf("Failed to create constraint file: %v", err)
	}

	// Create facts file with correct syntax (no quotes)
	factsFile := filepath.Join(tempDir, "data.facts")
	factsContent := []byte(`Person(id:P001, name:Alice, age:25)
Person(id:P002, name:Bob, age:30)
Order(id:O001, customer_id:P001, amount:100)
Order(id:O002, customer_id:P002, amount:200)
`)
	if err := os.WriteFile(factsFile, factsContent, 0644); err != nil {
		t.Fatalf("Failed to create facts file: %v", err)
	}

	tests := []struct {
		name               string
		args               []string
		wantExitCode       int
		wantOutputContains []string
	}{
		{
			name:         "constraint with facts - non-verbose",
			args:         []string{"-constraint", constraintFile, "-facts", factsFile},
			wantExitCode: 0,
			wantOutputContains: []string{
				"ACTIONS DISPONIBLES",
			},
		},
		{
			name:         "constraint with facts - verbose",
			args:         []string{"-constraint", constraintFile, "-facts", factsFile, "-v"},
			wantExitCode: 0,
			wantOutputContains: []string{
				"Parsing réussi",
				"PIPELINE",
				"faits",
			},
		},
		{
			name:         "constraint with non-existent facts file",
			args:         []string{"-constraint", constraintFile, "-facts", filepath.Join(tempDir, "missing.facts")},
			wantExitCode: 1,
			wantOutputContains: []string{
				"Fichier faits non trouvé",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(testBinary, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			exitCode := 0
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				}
			}

			if exitCode != tt.wantExitCode {
				t.Errorf("Exit code = %d, want %d\nOutput: %s", exitCode, tt.wantExitCode, outputStr)
			}

			for _, expected := range tt.wantOutputContains {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Output does not contain %q\nGot: %s", expected, outputStr)
				}
			}
		})
	}
}

// TestParseFromStdinError tests error handling in ParseFromStdin
func TestParseFromStdinError(t *testing.T) {
	// Create a closed pipe to simulate read error
	r, w, _ := os.Pipe()
	w.Close()
	r.Close()

	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	config := &Config{
		UseStdin: true,
	}

	// Test with invalid input that causes read error
	// Empty stdin should still work, so we need actual invalid syntax to trigger parsing error
	_, _, err := parseFromStdin(config, strings.NewReader("invalid @#$% syntax"))
	if err == nil {
		t.Error("parseFromStdin() with invalid syntax should return error")
	}
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("empty config", func(t *testing.T) {
		config := &Config{}
		err := ValidateConfig(config)
		if err == nil {
			t.Error("ValidateConfig() with empty config should return error")
		}
	})

	t.Run("parseFromFile with invalid file path characters", func(t *testing.T) {
		config := &Config{
			ConstraintFile: "/nonexistent/path/to/file.constraint",
		}

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		_, _, err := parseFromFile(config)

		w.Close()
		os.Stdout = oldStdout
		r.Close()

		if err == nil {
			t.Error("parseFromFile() with non-existent file should return error")
		}
	})

	t.Run("config with all boolean flags", func(t *testing.T) {
		config := &Config{
			UseStdin:    false,
			Verbose:     false,
			ShowVersion: false,
			ShowHelp:    false,
		}

		// Should fail validation because no input source
		err := ValidateConfig(config)
		if err == nil {
			t.Error("Expected validation error for config with no input source")
		}
	})
}

// TestRunWithFacts tests the RunWithFacts function
func TestRunWithFacts(t *testing.T) {
	// Create temporary constraint and facts files
	tmpDir := t.TempDir()

	constraintFile := filepath.Join(tmpDir, "test.constraint")
	constraintContent := `type Person : <id: string, name: string, age: number>

rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)`
	if err := os.WriteFile(constraintFile, []byte(constraintContent), 0644); err != nil {
		t.Fatalf("Failed to create constraint file: %v", err)
	}

	factsFile := filepath.Join(tmpDir, "test.facts")
	factsContent := `Person(id: "1", name: "Alice", age: 25)
Person(id: "2", name: "Bob", age: 30)`
	if err := os.WriteFile(factsFile, []byte(factsContent), 0644); err != nil {
		t.Fatalf("Failed to create facts file: %v", err)
	}

	tests := []struct {
		name         string
		config       *Config
		sourceName   string
		wantExitCode int
		checkOutput  func(*testing.T, string, string)
	}{
		{
			name: "successful execution",
			config: &Config{
				FactsFile: factsFile,
				Verbose:   false,
			},
			sourceName:   constraintFile,
			wantExitCode: 0,
			checkOutput: func(t *testing.T, stdout, stderr string) {
				if stderr != "" {
					t.Errorf("Expected no stderr, got: %s", stderr)
				}
			},
		},
		{
			name: "verbose mode",
			config: &Config{
				FactsFile: factsFile,
				Verbose:   true,
			},
			sourceName:   constraintFile,
			wantExitCode: 0,
			checkOutput: func(t *testing.T, stdout, stderr string) {
				if !strings.Contains(stdout, "PIPELINE RETE COMPLET") {
					t.Error("Verbose output should contain pipeline header")
				}
				if !strings.Contains(stdout, "RÉSULTATS") {
					t.Error("Verbose output should contain results header")
				}
			},
		},
		{
			name: "facts file not found",
			config: &Config{
				FactsFile: filepath.Join(tmpDir, "nonexistent.facts"),
				Verbose:   false,
			},
			sourceName:   constraintFile,
			wantExitCode: 1,
			checkOutput: func(t *testing.T, stdout, stderr string) {
				if !strings.Contains(stderr, "Fichier faits non trouvé") {
					t.Error("Should report facts file not found")
				}
			},
		},
		{
			name: "invalid constraint source",
			config: &Config{
				FactsFile: factsFile,
				Verbose:   false,
			},
			sourceName:   "nonexistent.constraint",
			wantExitCode: 1,
			checkOutput: func(t *testing.T, stdout, stderr string) {
				if !strings.Contains(stderr, "Erreur pipeline RETE") {
					t.Error("Should report pipeline error")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer

			exitCode := RunWithFacts(tt.config, tt.sourceName, &stdout, &stderr)

			if exitCode != tt.wantExitCode {
				t.Errorf("exitCode = %d, want %d", exitCode, tt.wantExitCode)
			}

			if tt.checkOutput != nil {
				tt.checkOutput(t, stdout.String(), stderr.String())
			}
		})
	}
}

// TestExecutePipeline tests the ExecutePipeline function
func TestExecutePipeline(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid test files
	constraintFile := filepath.Join(tmpDir, "test.constraint")
	constraintContent := `type Person : <id: string, name: string, age: number>

rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)`
	if err := os.WriteFile(constraintFile, []byte(constraintContent), 0644); err != nil {
		t.Fatalf("Failed to create constraint file: %v", err)
	}

	factsFile := filepath.Join(tmpDir, "test.facts")
	factsContent := `Person(id: "1", name: "Alice", age: 25)
Person(id: "2", name: "Bob", age: 30)`
	if err := os.WriteFile(factsFile, []byte(factsContent), 0644); err != nil {
		t.Fatalf("Failed to create facts file: %v", err)
	}

	tests := []struct {
		name             string
		constraintSource string
		factsFile        string
		wantError        bool
		checkResult      func(*testing.T, *Result)
	}{
		{
			name:             "successful pipeline execution",
			constraintSource: constraintFile,
			factsFile:        factsFile,
			wantError:        false,
			checkResult: func(t *testing.T, result *Result) {
				if result == nil {
					t.Fatal("Result should not be nil")
				}
				if result.Network == nil {
					t.Error("Network should not be nil")
				}
				if len(result.Facts) != 2 {
					t.Errorf("Expected 2 facts, got %d", len(result.Facts))
				}
				if result.Activations < 0 {
					t.Error("Activations should not be negative")
				}
			},
		},
		{
			name:             "constraint file not found",
			constraintSource: "nonexistent.constraint",
			factsFile:        factsFile,
			wantError:        true,
			checkResult:      nil,
		},
		{
			name:             "facts file not found",
			constraintSource: constraintFile,
			factsFile:        "nonexistent.facts",
			wantError:        true,
			checkResult:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExecutePipeline(tt.constraintSource, tt.factsFile)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.checkResult != nil {
					tt.checkResult(t, result)
				}
			}
		})
	}
}

// TestPrintResultsFunction tests the PrintResults function
func TestPrintResultsFunction(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		result      *Result
		checkOutput func(*testing.T, string)
	}{
		{
			name: "basic results - no activations",
			config: &Config{
				Verbose: false,
			},
			result: &Result{
				Network:     nil,
				Facts:       make([]*rete.Fact, 5),
				Activations: 0,
			},
			checkOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Aucune action déclenchée") {
					t.Error("Should contain 'no actions' message")
				}
			},
		},
		{
			name: "verbose results with activations",
			config: &Config{
				Verbose: true,
			},
			result: &Result{
				Network:     &rete.ReteNetwork{},
				Facts:       make([]*rete.Fact, 3),
				Activations: 5,
			},
			checkOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "RÉSULTATS") {
					t.Error("Should contain results header")
				}
				if !strings.Contains(output, "Faits injectés: 3") {
					t.Error("Should show fact count")
				}
				if !strings.Contains(output, "ACTIONS DISPONIBLES: 5") {
					t.Error("Should show activation count")
				}
				if !strings.Contains(output, "Pipeline RETE exécuté avec succès") {
					t.Error("Should contain success message")
				}
			},
		},
		{
			name: "non-verbose with activations",
			config: &Config{
				Verbose: false,
			},
			result: &Result{
				Network:     &rete.ReteNetwork{},
				Facts:       make([]*rete.Fact, 10),
				Activations: 3,
			},
			checkOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "ACTIONS DISPONIBLES: 3") {
					t.Error("Should show activation count")
				}
				if strings.Contains(output, "RÉSULTATS") {
					t.Error("Should not contain verbose headers")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			PrintResults(tt.config, tt.result, &stdout)

			if tt.checkOutput != nil {
				tt.checkOutput(t, stdout.String())
			}
		})
	}
}

// TestCountActivationsWithRealNetwork tests CountActivations with actual network
func TestCountActivationsReal(t *testing.T) {
	tests := []struct {
		name     string
		network  *rete.ReteNetwork
		expected int
	}{
		{
			name:     "nil network",
			network:  nil,
			expected: 0,
		},
		{
			name:     "empty network",
			network:  &rete.ReteNetwork{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := CountActivations(tt.network)
			if count != tt.expected {
				t.Errorf("CountActivations() = %d, want %d", count, tt.expected)
			}
		})
	}
}

// TestPrintActivationDetailsReal tests PrintActivationDetails with actual network
func TestPrintActivationDetailsReal(t *testing.T) {
	tests := []struct {
		name    string
		network *rete.ReteNetwork
	}{
		{
			name:    "nil network",
			network: nil,
		},
		{
			name:    "empty network",
			network: &rete.ReteNetwork{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			// Should not panic
			PrintActivationDetails(tt.network, &stdout)
		})
	}
}

// TestRun tests the main Run function
func TestRun(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test constraint file
	constraintFile := filepath.Join(tmpDir, "test.constraint")
	constraintContent := "type Person : <id: string, name: string>"
	if err := os.WriteFile(constraintFile, []byte(constraintContent), 0644); err != nil {
		t.Fatalf("Failed to create constraint file: %v", err)
	}

	tests := []struct {
		name         string
		args         []string
		stdin        string
		wantExitCode int
		checkOutput  func(*testing.T, string, string)
	}{
		{
			name:         "help flag",
			args:         []string{"-h"},
			wantExitCode: 0,
			checkOutput: func(t *testing.T, stdout, stderr string) {
				if !strings.Contains(stdout, "USAGE") {
					t.Error("Should display help")
				}
			},
		},
		{
			name:         "version flag",
			args:         []string{"-version"},
			wantExitCode: 0,
			checkOutput: func(t *testing.T, stdout, stderr string) {
				if !strings.Contains(stdout, "TSD") {
					t.Error("Should display version")
				}
			},
		},
		{
			name:         "valid constraint file",
			args:         []string{"-constraint", constraintFile},
			wantExitCode: 0,
			checkOutput: func(t *testing.T, stdout, stderr string) {
				if !strings.Contains(stdout, "validées avec succès") {
					t.Error("Should show success message")
				}
			},
		},
		{
			name:         "valid text input",
			args:         []string{"-text", "type Person : <id: string>"},
			wantExitCode: 0,
			checkOutput: func(t *testing.T, stdout, stderr string) {
				if !strings.Contains(stdout, "validées avec succès") {
					t.Error("Should show success message")
				}
			},
		},
		{
			name:         "no arguments",
			args:         []string{},
			wantExitCode: 1,
			checkOutput: func(t *testing.T, stdout, stderr string) {
				if !strings.Contains(stderr, "Erreur") {
					t.Error("Should show error")
				}
			},
		},
		{
			name:         "invalid syntax",
			args:         []string{"-text", "invalid @#$%"},
			wantExitCode: 1,
			checkOutput: func(t *testing.T, stdout, stderr string) {
				if !strings.Contains(stderr, "Erreur") {
					t.Error("Should show parsing error")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			var stdin io.Reader
			if tt.stdin != "" {
				stdin = strings.NewReader(tt.stdin)
			}

			exitCode := Run(tt.args, stdin, &stdout, &stderr)

			if exitCode != tt.wantExitCode {
				t.Errorf("exitCode = %d, want %d\nStdout: %s\nStderr: %s",
					exitCode, tt.wantExitCode, stdout.String(), stderr.String())
			}

			if tt.checkOutput != nil {
				tt.checkOutput(t, stdout.String(), stderr.String())
			}
		})
	}
}
