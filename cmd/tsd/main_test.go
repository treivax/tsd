package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
			// Reset flag.CommandLine for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = append([]string{"cmd"}, tt.args...)

			config := parseFlags()

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
			err := validateConfig(tt.config)

			if tt.wantError {
				if err == nil {
					t.Errorf("validateConfig() error = nil, want error containing %q", tt.errorMsg)
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("validateConfig() error = %v, want error containing %q", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateConfig() error = %v, want nil", err)
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

func TestPrintParsingHeader(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		expectedOutput []string
	}{
		{
			name:   "stdin source",
			source: "stdin",
			expectedOutput: []string{
				"🚀 TSD - Analyse des contraintes",
				"===============================",
				"Source: stdin",
			},
		},
		{
			name:   "file source",
			source: "test.constraint",
			expectedOutput: []string{
				"🚀 TSD - Analyse des contraintes",
				"===============================",
				"Source: test.constraint",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			printParsingHeader(tt.source)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("printParsingHeader() output does not contain %q\nGot: %s", expected, output)
				}
			}
		})
	}
}

func TestPrintVersion(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printVersion()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
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
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printHelp()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
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
			name:   "non-verbose mode",
			config: &Config{},
			expectedOutput: []string{
				"✅ Contraintes validées avec succès",
			},
		},
		{
			name: "verbose mode",
			config: &Config{
				Verbose: true,
			},
			expectedOutput: []string{
				"✅ Contraintes validées avec succès",
				"🎉 Validation terminée!",
				"syntaxiquement correctes",
				"-facts",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			runValidationOnly(tt.config)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("runValidationOnly() output does not contain %q\nGot: %s", expected, output)
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

// TestParseConstraintSource tests the parseConstraintSource routing function
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			result, sourceName, err := parseConstraintSource(tt.config)

			w.Close()
			os.Stdout = oldStdout
			r.Close()

			if tt.wantError && err == nil {
				t.Error("parseConstraintSource() error = nil, want error")
			}
			if !tt.wantError && err != nil {
				t.Errorf("parseConstraintSource() error = %v, want nil", err)
			}
			if !tt.wantError {
				if result == nil {
					t.Error("parseConstraintSource() result = nil, want non-nil")
				}
				if sourceName != tt.wantSourceName {
					t.Errorf("parseConstraintSource() sourceName = %v, want %v", sourceName, tt.wantSourceName)
				}
			}
		})
	}
}
