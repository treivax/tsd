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

// TestProgramState_MergeTypes_EdgeCases tests edge cases for type merging
func TestProgramState_MergeTypes_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		types1      string
		types2      string
		wantErr     bool
		errContains string
		checkCount  int
	}{
		{
			name:       "Identical types from different files",
			types1:     "type Person(id: string, name: string)",
			types2:     "type Person(id: string, name: string)",
			wantErr:    false,
			checkCount: 1,
		},
		{
			name:       "Compatible types - one with more fields",
			types1:     "type Person(id: string, name: string)",
			types2:     "type Person(id: string, name: string, age: number)",
			wantErr:    false,
			checkCount: 1,
		},
		{
			name:        "Incompatible types - different field types",
			types1:      "type Person(id: string, age: string)",
			types2:      "type Person(id: string, age: number)",
			wantErr:     true,
			errContains: "redefined incompatibly",
		},
		{
			name:       "Multiple distinct types",
			types1:     "type Person(id: string)\ntype Company(id: string)",
			types2:     "type Order(id: string)\ntype Product(id: string)",
			wantErr:    false,
			checkCount: 4,
		},
		{
			name:       "Empty types list",
			types1:     "type Person(id: string)",
			types2:     "",
			wantErr:    false,
			checkCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			// Create first file
			file1 := filepath.Join(tempDir, "types1.tsd")
			if err := os.WriteFile(file1, []byte(tt.types1), 0644); err != nil {
				t.Fatal(err)
			}

			// Create second file
			file2 := filepath.Join(tempDir, "types2.tsd")
			if err := os.WriteFile(file2, []byte(tt.types2), 0644); err != nil {
				t.Fatal(err)
			}

			ps := NewProgramState()

			// Parse first file
			err := ps.ParseAndMerge(file1)
			if err != nil {
				t.Fatalf("Failed to parse first file: %v", err)
			}

			// Parse second file
			err = ps.ParseAndMerge(file2)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error containing %q, got nil", tt.errContains)
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Error %q does not contain %q", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.checkCount > 0 && len(ps.Types) != tt.checkCount {
					t.Errorf("Expected %d types, got %d", tt.checkCount, len(ps.Types))
				}
			}
		})
	}
}

// TestProgramState_MergeRules_EdgeCases tests edge cases for rule merging
func TestProgramState_MergeRules_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		types      string
		rules1     string
		rules2     string
		wantRules  int
		wantErrors int
	}{
		{
			name:       "Duplicate rule ID - should skip second",
			types:      "type T(x: number)",
			rules1:     `rule r1: {t: T} / t.x > 0 ==> DoA()`,
			rules2:     `rule r1: {t: T} / t.x < 100 ==> DoB()`,
			wantRules:  1,
			wantErrors: 1,
		},
		{
			name:       "Rule with undefined type - should skip",
			types:      "type T(x: number)",
			rules1:     `rule r1: {u: Unknown} / u.x > 0 ==> DoA()`,
			rules2:     `rule r2: {t: T} / t.x > 0 ==> DoB()`,
			wantRules:  1,
			wantErrors: 1,
		},
		{
			name:       "Multiple valid rules",
			types:      "type T(x: number)",
			rules1:     `rule r1: {t: T} / t.x > 0 ==> DoA()`,
			rules2:     `rule r2: {t: T} / t.x < 100 ==> DoB()`,
			wantRules:  2,
			wantErrors: 0,
		},
		{
			name:       "Multiple rules with different IDs",
			types:      "type T(x: number)",
			rules1:     `rule r1: {t: T} / t.x > 0 ==> DoA()`,
			rules2:     `rule r3: {t: T} / t.x < 100 ==> DoB()`,
			wantRules:  2,
			wantErrors: 0,
		},
		{
			name:       "Rule with undefined field",
			types:      "type T(x: number)",
			rules1:     `rule r1: {t: T} / t.y > 0 ==> DoA()`,
			rules2:     `rule r2: {t: T} / t.x > 0 ==> DoB()`,
			wantRules:  1,
			wantErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			// Create types file
			typesFile := filepath.Join(tempDir, "types.tsd")
			if err := os.WriteFile(typesFile, []byte(tt.types), 0644); err != nil {
				t.Fatal(err)
			}

			// Create rules files
			rules1File := filepath.Join(tempDir, "rules1.tsd")
			if err := os.WriteFile(rules1File, []byte(tt.rules1), 0644); err != nil {
				t.Fatal(err)
			}

			rules2File := filepath.Join(tempDir, "rules2.tsd")
			if err := os.WriteFile(rules2File, []byte(tt.rules2), 0644); err != nil {
				t.Fatal(err)
			}

			ps := NewProgramState()

			// Parse types
			if err := ps.ParseAndMerge(typesFile); err != nil {
				t.Fatalf("Failed to parse types: %v", err)
			}

			// Parse rules (errors should be non-blocking)
			_ = ps.ParseAndMerge(rules1File)
			_ = ps.ParseAndMerge(rules2File)

			if len(ps.Rules) != tt.wantRules {
				t.Errorf("Expected %d rules, got %d", tt.wantRules, len(ps.Rules))
			}

			if len(ps.Errors) != tt.wantErrors {
				t.Errorf("Expected %d errors, got %d", tt.wantErrors, len(ps.Errors))
			}
		})
	}
}

// TestProgramState_MergeFacts_EdgeCases tests edge cases for fact merging
func TestProgramState_MergeFacts_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		types      string
		facts      string
		wantFacts  int
		wantErrors int
	}{
		{
			name:  "Valid facts",
			types: "type Person(id: string, age: number)",
			facts: `Person(id: "p1", age: 30)
Person(id: "p2", age: 25)`,
			wantFacts:  2,
			wantErrors: 0,
		},
		{
			name:       "Fact with undefined type",
			types:      "type Person(id: string, age: number)",
			facts:      `Unknown(id: "u1", value: 10)`,
			wantFacts:  0,
			wantErrors: 1,
		},
		{
			name:       "Fact with undefined field",
			types:      "type Person(id: string, age: number)",
			facts:      `Person(id: "p1", name: "Alice")`,
			wantFacts:  0,
			wantErrors: 1,
		},
		{
			name:       "Fact with wrong field type",
			types:      "type Person(id: string, age: number)",
			facts:      `Person(id: "p1", age: "thirty")`,
			wantFacts:  0,
			wantErrors: 1,
		},
		{
			name:  "Mixed valid and invalid facts",
			types: "type Person(id: string, age: number)",
			facts: `Person(id: "p1", age: 30)
Person(id: "p2", name: "Invalid")
Person(id: "p3", age: 25)`,
			wantFacts:  2,
			wantErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			// Create types file
			typesFile := filepath.Join(tempDir, "types.tsd")
			if err := os.WriteFile(typesFile, []byte(tt.types), 0644); err != nil {
				t.Fatal(err)
			}

			// Create facts file
			factsFile := filepath.Join(tempDir, "facts.tsd")
			if err := os.WriteFile(factsFile, []byte(tt.facts), 0644); err != nil {
				t.Fatal(err)
			}

			ps := NewProgramState()

			// Parse types
			if err := ps.ParseAndMerge(typesFile); err != nil {
				t.Fatalf("Failed to parse types: %v", err)
			}

			// Parse facts (errors should be non-blocking)
			_ = ps.ParseAndMerge(factsFile)

			if len(ps.Facts) != tt.wantFacts {
				t.Errorf("Expected %d facts, got %d", tt.wantFacts, len(ps.Facts))
			}

			if len(ps.Errors) != tt.wantErrors {
				t.Errorf("Expected %d errors, got %d", tt.wantErrors, len(ps.Errors))
			}
		})
	}
}

// TestProgramState_ParseAndMergeContent_EdgeCases tests edge cases for content parsing
func TestProgramState_ParseAndMergeContent_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		filename    string
		setupState  func(*ProgramState)
		wantErr     bool
		errContains string
	}{
		{
			name:        "Empty content",
			content:     "",
			filename:    "test.tsd",
			wantErr:     true,
			errContains: "content cannot be empty",
		},
		{
			name:        "Empty filename",
			content:     "type T(x: number)",
			filename:    "",
			wantErr:     true,
			errContains: "filename cannot be empty",
		},
		{
			name:     "Nil state",
			content:  "type T(x: number)",
			filename: "test.tsd",
			setupState: func(ps *ProgramState) {
				// Will test with nil state
			},
			wantErr:     true,
			errContains: "ProgramState is nil",
		},
		{
			name:     "Valid content",
			content:  "type Person(id: string, name: string)",
			filename: "test.tsd",
			wantErr:  false,
		},
		{
			name:     "Content with comments",
			content:  "// This is a comment\ntype Person(id: string, name: string)\n// Another comment",
			filename: "test.tsd",
			wantErr:  false,
		},
		{
			name:     "Content with reset",
			content:  "reset\ntype Person(id: string)",
			filename: "test.tsd",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ps *ProgramState
			if tt.name == "Nil state" {
				ps = nil
			} else {
				ps = NewProgramState()
				if tt.setupState != nil {
					tt.setupState(ps)
				}
			}

			err := ps.ParseAndMergeContent(tt.content, tt.filename)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error containing %q, got nil", tt.errContains)
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Error %q does not contain %q", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestProgramState_Reset_EdgeCases tests edge cases for reset functionality
func TestProgramState_Reset_EdgeCases(t *testing.T) {
	tempDir := t.TempDir()

	// Create initial files
	typesFile := filepath.Join(tempDir, "types.tsd")
	typesContent := "type Person(id: string, name: string)"
	if err := os.WriteFile(typesFile, []byte(typesContent), 0644); err != nil {
		t.Fatal(err)
	}

	rulesFile := filepath.Join(tempDir, "rules.tsd")
	rulesContent := "rule r1: {p: Person} / p.id != \"\" ==> Notify()"
	if err := os.WriteFile(rulesFile, []byte(rulesContent), 0644); err != nil {
		t.Fatal(err)
	}

	factsFile := filepath.Join(tempDir, "facts.tsd")
	factsContent := `Person(id: "p1", name: "Alice")`
	if err := os.WriteFile(factsFile, []byte(factsContent), 0644); err != nil {
		t.Fatal(err)
	}

	resetFile := filepath.Join(tempDir, "reset.tsd")
	resetContent := "reset\ntype NewType(x: number)"
	if err := os.WriteFile(resetFile, []byte(resetContent), 0644); err != nil {
		t.Fatal(err)
	}

	ps := NewProgramState()

	// Parse initial files
	if err := ps.ParseAndMerge(typesFile); err != nil {
		t.Fatalf("Failed to parse types: %v", err)
	}
	if err := ps.ParseAndMerge(rulesFile); err != nil {
		t.Fatalf("Failed to parse rules: %v", err)
	}
	if err := ps.ParseAndMerge(factsFile); err != nil {
		t.Fatalf("Failed to parse facts: %v", err)
	}

	// Verify initial state
	if len(ps.Types) != 1 {
		t.Errorf("Expected 1 type before reset, got %d", len(ps.Types))
	}
	if len(ps.Rules) != 1 {
		t.Errorf("Expected 1 rule before reset, got %d", len(ps.Rules))
	}
	if len(ps.Facts) != 1 {
		t.Errorf("Expected 1 fact before reset, got %d", len(ps.Facts))
	}
	if len(ps.FilesParsed) != 3 {
		t.Errorf("Expected 3 files parsed before reset, got %d", len(ps.FilesParsed))
	}

	// Mark a rule ID as used
	ps.RuleIDs["r1"] = true

	// Parse reset file
	if err := ps.ParseAndMerge(resetFile); err != nil {
		t.Fatalf("Failed to parse reset file: %v", err)
	}

	// Verify state after reset
	if len(ps.Types) != 1 {
		t.Errorf("Expected 1 type after reset, got %d", len(ps.Types))
	}
	if _, exists := ps.Types["Person"]; exists {
		t.Error("Old type 'Person' should not exist after reset")
	}
	if _, exists := ps.Types["NewType"]; !exists {
		t.Error("New type 'NewType' should exist after reset")
	}
	if len(ps.Rules) != 0 {
		t.Errorf("Expected 0 rules after reset, got %d", len(ps.Rules))
	}
	if len(ps.Facts) != 0 {
		t.Errorf("Expected 0 facts after reset, got %d", len(ps.Facts))
	}
	if len(ps.RuleIDs) != 0 {
		t.Errorf("Expected 0 rule IDs after reset, got %d", len(ps.RuleIDs))
	}
	if len(ps.FilesParsed) != 1 {
		t.Errorf("Expected 1 file parsed after reset, got %d", len(ps.FilesParsed))
	}
}

// TestProgramState_ValidateFieldAccesses_EdgeCases tests edge cases for field access validation
func TestProgramState_ValidateFieldAccesses_EdgeCases(t *testing.T) {
	ps := NewProgramState()

	// Add a type
	ps.Types["Person"] = &TypeDefinition{
		Type: "type",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "age", Type: "number"},
		},
	}

	variables := map[string]string{
		"p": "Person",
	}

	tests := []struct {
		name        string
		data        interface{}
		wantErr     bool
		errContains string
	}{
		{
			name: "Valid field access",
			data: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "id",
			},
			wantErr: false,
		},
		{
			name: "Invalid field access - nonexistent field",
			data: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "name",
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "Nested valid field access",
			data: map[string]interface{}{
				"type": "comparison",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "age",
				},
				"operator": ">",
				"right":    18,
			},
			wantErr: false,
		},
		{
			name: "Array with field accesses",
			data: []interface{}{
				map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "id",
				},
				map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "age",
				},
			},
			wantErr: false,
		},
		{
			name:    "Non-fieldAccess map",
			data:    map[string]interface{}{"type": "literal", "value": 42},
			wantErr: false,
		},
		{
			name:    "Simple value",
			data:    "simple string",
			wantErr: false,
		},
		{
			name:    "Nil data",
			data:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ps.validateFieldAccesses(tt.data, variables)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error containing %q, got nil", tt.errContains)
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Error %q does not contain %q", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestProgramState_ToProgram_EdgeCases tests edge cases for ToProgram conversion
func TestProgramState_ToProgram_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		setupState func(*ProgramState)
		checkFunc  func(*testing.T, *Program)
	}{
		{
			name: "Empty state",
			setupState: func(ps *ProgramState) {
				// Empty state
			},
			checkFunc: func(t *testing.T, p *Program) {
				if len(p.Types) != 0 {
					t.Errorf("Expected 0 types, got %d", len(p.Types))
				}
				if len(p.Expressions) != 0 {
					t.Errorf("Expected 0 expressions, got %d", len(p.Expressions))
				}
				if len(p.Facts) != 0 {
					t.Errorf("Expected 0 facts, got %d", len(p.Facts))
				}
			},
		},
		{
			name: "State with only types",
			setupState: func(ps *ProgramState) {
				ps.Types["Person"] = &TypeDefinition{
					Type: "type",
					Name: "Person",
					Fields: []Field{
						{Name: "id", Type: "string"},
					},
				}
			},
			checkFunc: func(t *testing.T, p *Program) {
				if len(p.Types) != 1 {
					t.Errorf("Expected 1 type, got %d", len(p.Types))
				}
			},
		},
		{
			name: "State with multiple of each",
			setupState: func(ps *ProgramState) {
				ps.Types["T1"] = &TypeDefinition{Type: "type", Name: "T1"}
				ps.Types["T2"] = &TypeDefinition{Type: "type", Name: "T2"}
				ps.Rules = append(ps.Rules, &Expression{RuleId: "r1"})
				ps.Rules = append(ps.Rules, &Expression{RuleId: "r2"})
				ps.Facts = append(ps.Facts, &Fact{TypeName: "T1"})
				ps.Facts = append(ps.Facts, &Fact{TypeName: "T2"})
			},
			checkFunc: func(t *testing.T, p *Program) {
				if len(p.Types) != 2 {
					t.Errorf("Expected 2 types, got %d", len(p.Types))
				}
				if len(p.Expressions) != 2 {
					t.Errorf("Expected 2 expressions, got %d", len(p.Expressions))
				}
				if len(p.Facts) != 2 {
					t.Errorf("Expected 2 facts, got %d", len(p.Facts))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := NewProgramState()
			tt.setupState(ps)

			program := ps.ToProgram()

			if program == nil {
				t.Fatal("ToProgram() returned nil")
			}

			tt.checkFunc(t, program)
		})
	}
}
