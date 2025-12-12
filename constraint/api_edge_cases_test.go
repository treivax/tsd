// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestReadFileContent_EdgeCases tests edge cases for ReadFileContent
func TestReadFileContent_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		wantErr     bool
		errContains string
	}{
		{
			name: "Empty file",
			setupFile: func(t *testing.T) string {
				f, err := os.CreateTemp("", "empty_*.tsd")
				if err != nil {
					t.Fatal(err)
				}
				defer f.Close()
				return f.Name()
			},
			wantErr: false,
		},
		{
			name: "Nonexistent file",
			setupFile: func(t *testing.T) string {
				return "/nonexistent/path/to/file.tsd"
			},
			wantErr:     true,
			errContains: "failed to read file",
		},
		{
			name: "Directory instead of file",
			setupFile: func(t *testing.T) string {
				dir, err := os.MkdirTemp("", "testdir_*")
				if err != nil {
					t.Fatal(err)
				}
				return dir
			},
			wantErr:     true,
			errContains: "failed to read file",
		},
		{
			name: "File with special characters",
			setupFile: func(t *testing.T) string {
				f, err := os.CreateTemp("", "special_*.tsd")
				if err != nil {
					t.Fatal(err)
				}
				defer f.Close()
				content := "type Test(field: string) // Special: éàü 中文 🎉\n"
				if _, err := f.WriteString(content); err != nil {
					t.Fatal(err)
				}
				return f.Name()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := tt.setupFile(t)
			defer func() {
				// Clean up if file exists
				if _, err := os.Stat(filename); err == nil {
					os.RemoveAll(filename)
				}
			}()

			content, err := ReadFileContent(filename)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ReadFileContent() expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ReadFileContent() error = %v, should contain %q", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ReadFileContent() unexpected error = %v", err)
				}
				if tt.name == "Empty file" && content != "" {
					t.Errorf("ReadFileContent() expected empty string, got %q", content)
				}
			}
		})
	}
}

// TestParseConstraint_EdgeCases tests edge cases for ParseConstraint
func TestParseConstraint_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		input       []byte
		wantErr     bool
		errContains string
	}{
		{
			name:     "Empty input",
			filename: "empty.tsd",
			input:    []byte(""),
			wantErr:  false, // Empty programs are valid
		},
		{
			name:     "Only whitespace",
			filename: "whitespace.tsd",
			input:    []byte("   \n\t\n   "),
			wantErr:  false,
		},
		{
			name:     "Only comments",
			filename: "comments.tsd",
			input:    []byte("// Comment 1\n// Comment 2\n/* Block comment */"),
			wantErr:  false,
		},
		{
			name:        "Invalid syntax",
			filename:    "invalid.tsd",
			input:       []byte("this is not valid syntax !!!"),
			wantErr:     true,
			errContains: "invalid.tsd",
		},
		{
			name:     "Multiple types",
			filename: "multi.tsd",
			input:    []byte("type A(x: number)\ntype B(y: string)\ntype C(z: bool)"),
			wantErr:  false,
		},
		{
			name:        "Unterminated string",
			filename:    "bad_string.tsd",
			input:       []byte(`type Test(name: "unterminated)`),
			wantErr:     true,
			errContains: "bad_string.tsd",
		},
		{
			name:     "Unicode content",
			filename: "unicode.tsd",
			input:    []byte("type 用户(名字: string, 年龄: number)"),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseConstraint(tt.filename, tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseConstraint() expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ParseConstraint() error = %v, should contain %q", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ParseConstraint() unexpected error = %v", err)
				}
				if result == nil {
					t.Error("ParseConstraint() returned nil result without error")
				}
			}
		})
	}
}

// TestValidateConstraintProgram_EdgeCases tests edge cases for ValidateConstraintProgram
func TestValidateConstraintProgram_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		errContains string
	}{
		{
			name: "Program with action using undefined type",
			input: `
type Person(id: string, name: string)
action Notify(personId: string)
rule r1: {p: UnknownType} / p.id == "test" ==> Notify(p.id)
`,
			wantErr:     true,
			errContains: "UnknownType",
		},
		{
			name: "Action with wrong parameter type",
			input: `
type Person(id: string, age: number)
action SetAge(personId: string, newAge: string)
rule r1: {p: Person} / p.age > 18 ==> SetAge(p.id, p.age)
`,
			wantErr:     true,
			errContains: "type mismatch",
		},
		{
			name: "Program with no types",
			input: `
action DoSomething()
`,
			wantErr: false,
		},
		{
			name: "Action with no parameters",
			input: `
type Event(id: string)
action Trigger()
rule r1: {e: Event} / e.id == "test" ==> Trigger()
`,
			wantErr: false,
		},
		{
			name: "Multiple actions in rule",
			input: `
type Order(id: string, status: string)
action UpdateStatus(orderId: string, status: string)
action SendNotification(orderId: string)
rule r1: {o: Order} / o.status == "pending" ==> UpdateStatus(o.id, "processed")
`,
			wantErr: false,
		},
		{
			name: "Rule with field access to undefined field",
			input: `
type Product(id: string, name: string)
action Alert(msg: string)
rule r1: {p: Product} / p.price > 100 ==> Alert("Expensive")
`,
			wantErr:     true,
			errContains: "price",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseConstraint("test.tsd", []byte(tt.input))
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			err = ValidateConstraintProgram(result)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateConstraintProgram() expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateConstraintProgram() error = %v, should contain %q", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateConstraintProgram() unexpected error = %v", err)
				}
			}
		})
	}
}

// TestExtractFactsFromProgram_EdgeCases tests edge cases for ExtractFactsFromProgram
func TestExtractFactsFromProgram_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantFacts   int
		wantErr     bool
		errContains string
	}{
		{
			name: "Program with no facts",
			input: `
type Person(id: string, name: string)
rule r1: {p: Person} / p.id == "test" ==> Alert()
`,
			wantFacts: 0,
			wantErr:   false,
		},
		{
			name: "Program with multiple facts",
			input: `
type Person(id: string, name: string, age: number)
Person(id: "p1", name: "Alice", age: 30)
Person(id: "p2", name: "Bob", age: 25)
Person(id: "p3", name: "Charlie", age: 35)
`,
			wantFacts: 3,
			wantErr:   false,
		},
		{
			name: "Facts with bool values",
			input: `
type Config(key: string, enabled: bool)
Config(key: "feature1", enabled: true)
Config(key: "feature2", enabled: false)
`,
			wantFacts: 2,
			wantErr:   false,
		},
		{
			name: "Facts with number values",
			input: `
type Metric(name: string, value: number)
Metric(name: "cpu", value: 75.5)
Metric(name: "memory", value: 80)
Metric(name: "disk", value: 45.3)
`,
			wantFacts: 3,
			wantErr:   false,
		},
		{
			name: "Mixed types and facts",
			input: `
type A(x: number)
type B(y: string)
A(x: 10)
B(y: "test")
A(x: 20)
`,
			wantFacts: 3,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseConstraint("test.tsd", []byte(tt.input))
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			facts, err := ExtractFactsFromProgram(result)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ExtractFactsFromProgram() expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ExtractFactsFromProgram() error = %v, should contain %q", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ExtractFactsFromProgram() unexpected error = %v", err)
				}
				if len(facts) != tt.wantFacts {
					t.Errorf("ExtractFactsFromProgram() got %d facts, want %d", len(facts), tt.wantFacts)
				}
			}
		})
	}
}

// TestConvertResultToProgram_EdgeCases tests edge cases for ConvertResultToProgram
func TestConvertResultToProgram_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		checkTypes  bool
		checkRules  bool
		checkFacts  bool
		wantErr     bool
		errContains string
	}{
		{
			name:       "Empty program",
			input:      "",
			checkTypes: false,
			checkRules: false,
			checkFacts: false,
			wantErr:    false,
		},
		{
			name:       "Only types",
			input:      "type Person(id: string)\ntype Company(name: string)",
			checkTypes: true,
			checkRules: false,
			checkFacts: false,
			wantErr:    false,
		},
		{
			name:       "Only facts",
			input:      "type T(x: number)\nT(x: 1)\nT(x: 2)",
			checkTypes: true,
			checkRules: false,
			checkFacts: true,
			wantErr:    false,
		},
		{
			name: "Complex program",
			input: `
type Order(id: string, amount: number, status: string)
action ProcessOrder(orderId: string)
Order(id: "O1", amount: 100, status: "pending")
Order(id: "O2", amount: 200, status: "pending")
rule r1: {o: Order} / o.status == "pending" AND o.amount > 50 ==> ProcessOrder(o.id)
`,
			checkTypes: true,
			checkRules: true,
			checkFacts: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseConstraint("test.tsd", []byte(tt.input))
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			program, err := ConvertResultToProgram(result)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ConvertResultToProgram() expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ConvertResultToProgram() error = %v, should contain %q", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ConvertResultToProgram() unexpected error = %v", err)
				}
				if program == nil {
					t.Error("ConvertResultToProgram() returned nil program")
					return
				}

				if tt.checkTypes && len(program.Types) == 0 {
					t.Error("Expected types but got none")
				}
				if tt.checkRules && len(program.Expressions) == 0 {
					t.Error("Expected rules but got none")
				}
				if tt.checkFacts && len(program.Facts) == 0 {
					t.Error("Expected facts but got none")
				}
			}
		})
	}
}

// TestConvertToReteProgram_EdgeCases tests edge cases for ConvertToReteProgram
func TestConvertToReteProgram_EdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		checkTypes    bool
		checkActions  bool
		checkRules    bool
		checkRemovals bool
	}{
		{
			name:       "Minimal program",
			input:      "type T(x: number)",
			checkTypes: true,
		},
		{
			name: "Program with actions",
			input: `
type Event(id: string)
action Log(msg: string)
rule r1: {e: Event} / e.id != "" ==> Log("Event detected")
`,
			checkTypes:   true,
			checkActions: true,
			checkRules:   true,
		},
		{
			name: "Program with multiple types and rules",
			input: `
type A(x: number)
type B(y: string)
type C(z: bool)
rule r1: {a: A} / a.x > 0 ==> DoA()
rule r2: {b: B} / b.y == "test" ==> DoB()
`,
			checkTypes: true,
			checkRules: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseConstraint("test.tsd", []byte(tt.input))
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			program, err := ConvertResultToProgram(result)
			if err != nil {
				t.Fatalf("Failed to convert to program: %v", err)
			}

			reteProgram, err := ConvertToReteProgram(program)
			if err != nil {
				t.Fatalf("ConvertToReteProgram() failed: %v", err)
			}
			if reteProgram == nil {
				t.Fatal("ConvertToReteProgram() returned nil")
			}

			reteProgramMap, ok := reteProgram.(map[string]interface{})
			if !ok {
				t.Fatal("ConvertToReteProgram() did not return a map")
			}

			if tt.checkTypes {
				types, hasTypes := reteProgramMap["types"]
				if !hasTypes {
					t.Error("RETE program missing 'types' field")
				}
				if typesSlice, ok := types.([]interface{}); ok {
					if len(typesSlice) == 0 {
						t.Error("Expected types in RETE program but got none")
					}
				}
			}

			if tt.checkActions {
				actions, hasActions := reteProgramMap["actions"]
				if !hasActions {
					t.Error("RETE program missing 'actions' field")
				}
				if actionsSlice, ok := actions.([]interface{}); ok {
					if len(actionsSlice) == 0 {
						t.Error("Expected actions in RETE program but got none")
					}
				}
			}

			if tt.checkRules {
				expressions, hasExpressions := reteProgramMap["expressions"]
				if !hasExpressions {
					t.Error("RETE program missing 'expressions' field")
				}
				if expSlice, ok := expressions.([]interface{}); ok {
					if len(expSlice) == 0 {
						t.Error("Expected expressions in RETE program but got none")
					}
				}
			}
		})
	}
}

// TestParseConstraintFile_EdgeCases tests edge cases for ParseConstraintFile
func TestParseConstraintFile_EdgeCases(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		setupFile   func() string
		wantErr     bool
		errContains string
	}{
		{
			name: "Very large file",
			setupFile: func() string {
				filename := filepath.Join(tempDir, "large.tsd")
				content := "type T(x: number)\n"
				// Create a file with many repeated type definitions
				for i := 0; i < 1000; i++ {
					content += "T(x: " + string(rune('0'+i%10)) + ")\n"
				}
				if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
				return filename
			},
			wantErr: false,
		},
		{
			name: "File with BOM",
			setupFile: func() string {
				filename := filepath.Join(tempDir, "bom.tsd")
				content := "\xEF\xBB\xBFtype Test(x: number)"
				if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
				return filename
			},
			wantErr:     true,
			errContains: "no match found",
		},
		{
			name: "File with mixed line endings",
			setupFile: func() string {
				filename := filepath.Join(tempDir, "mixed_newlines.tsd")
				content := "type A(x: number)\r\ntype B(y: string)\ntype C(z: bool)\r\n"
				if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
				return filename
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := tt.setupFile()

			result, err := ParseConstraintFile(filename)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseConstraintFile() expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ParseConstraintFile() error = %v, should contain %q", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ParseConstraintFile() unexpected error = %v", err)
				}
				if result == nil {
					t.Error("ParseConstraintFile() returned nil result")
				}
			}
		})
	}
}

// TestIterativeParser_ErrorRecovery tests error recovery in iterative parsing
func TestIterativeParser_ErrorRecovery(t *testing.T) {
	parser := NewIterativeParser()

	// Parse valid types
	validTypes := "type Person(id: string, name: string)"
	err := parser.ParseContent(validTypes, "types.tsd")
	if err != nil {
		t.Fatalf("Failed to parse valid types: %v", err)
	}

	// Try to parse invalid content
	invalidContent := "this is not valid syntax"
	err = parser.ParseContent(invalidContent, "invalid.tsd")
	if err == nil {
		t.Error("Expected error for invalid content, got nil")
	}

	// Verify that valid types are still present
	program := parser.GetProgram()
	if len(program.Types) != 1 {
		t.Errorf("Expected 1 type after error, got %d", len(program.Types))
	}

	// Parse more valid content
	validFacts := "Person(id: \"p1\", name: \"Alice\")"
	err = parser.ParseContent(validFacts, "facts.tsd")
	if err != nil {
		t.Fatalf("Failed to parse valid facts after error: %v", err)
	}

	program = parser.GetProgram()
	if len(program.Facts) != 1 {
		t.Errorf("Expected 1 fact, got %d", len(program.Facts))
	}
}

// TestIterativeParser_ConcurrentAccess tests concurrent access patterns
func TestIterativeParser_ConcurrentAccess(t *testing.T) {
	parser := NewIterativeParser()

	// Parse initial content
	content := "type T(x: number)\nT(x: 1)"
	err := parser.ParseContent(content, "test.tsd")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Get program multiple times (simulating concurrent reads)
	for i := 0; i < 10; i++ {
		program := parser.GetProgram()
		if program == nil {
			t.Error("GetProgram() returned nil")
			continue
		}
		if len(program.Types) != 1 {
			t.Errorf("Iteration %d: expected 1 type, got %d", i, len(program.Types))
		}
	}

	// Get state multiple times
	for i := 0; i < 10; i++ {
		state := parser.GetState()
		if state == nil {
			t.Error("GetState() returned nil")
			continue
		}
		if state.GetTypesCount() != 1 {
			t.Errorf("Iteration %d: expected 1 type in state, got %d", i, state.GetTypesCount())
		}
	}

	// Get statistics multiple times
	for i := 0; i < 10; i++ {
		stats := parser.GetParsingStatistics()
		if stats.TypesCount != 1 {
			t.Errorf("Iteration %d: expected 1 type in stats, got %d", i, stats.TypesCount)
		}
	}
}
