// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProgramStateIterativeParsing(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "program_state_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files

	// 1. Type definitions file
	typesFile := filepath.Join(tempDir, "types.tsd")
	typesContent := `// Type definitions
type Person(id: string, name: string, age:number)
type Company(id: string, name: string, employees:number)`

	err = os.WriteFile(typesFile, []byte(typesContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write types file: %v", err)
	}

	// 2. Rules file
	rulesFile := filepath.Join(tempDir, "rules.tsd")
	rulesContent := `// Business rules
rule r1 : {p: Person, c: Company} / p.age > 25 AND c.employees > 100 ==> AddToWorkforce(p.id, c.id)`

	err = os.WriteFile(rulesFile, []byte(rulesContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write rules file: %v", err)
	}

	// 3. Facts file
	factsFile := filepath.Join(tempDir, "facts.tsd")
	factsContent := `// Initial facts
Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 22)
Company(id: "C001", name: "TechCorp", employees: 250)`

	err = os.WriteFile(factsFile, []byte(factsContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write facts file: %v", err)
	}

	// Test iterative parsing
	ps := NewProgramState()

	// Parse types first
	err = ps.ParseAndMerge(typesFile)
	if err != nil {
		t.Fatalf("Failed to parse types file: %v", err)
	}

	// Verify types were parsed
	if ps.GetTypesCount() != 2 {
		t.Errorf("Expected 2 types, got %d", ps.GetTypesCount())
	}

	if _, exists := ps.GetTypes()["Person"]; !exists {
		t.Error("Person type not found")
	}

	if _, exists := ps.GetTypes()["Company"]; !exists {
		t.Error("Company type not found")
	}

	// Parse rules (should validate against existing types)
	err = ps.ParseAndMerge(rulesFile)
	if err != nil {
		t.Fatalf("Failed to parse rules file: %v", err)
	}

	// Verify rules were parsed
	if ps.GetRulesCount() == 0 {
		t.Error("No rules were parsed")
	}

	// Parse facts (should validate against existing types)
	err = ps.ParseAndMerge(factsFile)
	if err != nil {
		t.Fatalf("Failed to parse facts file: %v", err)
	}

	// Verify facts were parsed
	if ps.GetFactsCount() != 3 {
		t.Errorf("Expected 3 facts, got %d", ps.GetFactsCount())
	}

	// Verify all files were recorded
	if len(ps.GetFilesParsed()) != 3 {
		t.Errorf("Expected 3 parsed files, got %d", len(ps.GetFilesParsed()))
	}

	// Convert back to Program structure
	program := ps.ToProgram()
	if program == nil {
		t.Fatal("Failed to convert to Program")
	}

	if len(program.Types) != 2 {
		t.Errorf("Program types: expected 2, got %d", len(program.Types))
	}

	if len(program.Facts) != 3 {
		t.Errorf("Program facts: expected 3, got %d", len(program.Facts))
	}

	t.Logf("Successfully parsed %d files with %d types, %d rules, %d facts",
		len(ps.GetFilesParsed()), ps.GetTypesCount(), ps.GetRulesCount(), ps.GetFactsCount())
}

func TestProgramStateTypeValidation(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "program_state_validation_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create types file
	typesFile := filepath.Join(tempDir, "types.tsd")
	typesContent := `type Person(id: string, name:string)`
	err = os.WriteFile(typesFile, []byte(typesContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write types file: %v", err)
	}

	// Create facts file with invalid type
	invalidFactsFile := filepath.Join(tempDir, "invalid_facts.tsd")
	invalidFactsContent := `UnknownType(id: "U001", name: "Test")`
	err = os.WriteFile(invalidFactsFile, []byte(invalidFactsContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid facts file: %v", err)
	}

	// Test validation
	ps := NewProgramState()

	// Parse types
	err = ps.ParseAndMerge(typesFile)
	if err != nil {
		t.Fatalf("Failed to parse types: %v", err)
	}

	// Try to parse facts with undefined type - should not fail (non-blocking)
	// but should record validation errors
	err = ps.ParseAndMerge(invalidFactsFile)
	if err != nil {
		t.Fatalf("ParseAndMerge should not fail with validation errors: %v", err)
	}

	// Check that validation errors were recorded
	if !ps.HasErrors() {
		t.Error("Expected validation errors to be recorded for undefined type")
	} else {
		t.Logf("Correctly detected %d validation error(s)", ps.GetErrorCount())
	}
}

func TestProgramState_ParseAndMergeContent(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		filename    string
		expectError bool
		checkFunc   func(*testing.T, *ProgramState)
	}{
		{
			name:     "parse type from content",
			content:  `type Person(id: string, name: string, age: number)`,
			filename: "test.tsd",
			checkFunc: func(t *testing.T, ps *ProgramState) {
				if ps.GetTypesCount() != 1 {
					t.Errorf("Expected 1 type, got %d", ps.GetTypesCount())
				}
				if _, exists := ps.GetTypes()["Person"]; !exists {
					t.Error("Person type not found")
				}
			},
		},
		{
			name: "parse rule from content",
			content: `type Person(id: string, age: number)
rule r1 : {p: Person} / p.age > 18 ==> Approve(p.id)`,
			filename: "test.tsd",
			checkFunc: func(t *testing.T, ps *ProgramState) {
				if ps.GetRulesCount() != 1 {
					t.Errorf("Expected 1 rule, got %d", ps.GetRulesCount())
				}
			},
		},
		{
			name: "parse fact from content",
			content: `type Person(id: string)
Person(id: "P001")`,
			filename: "test.tsd",
			checkFunc: func(t *testing.T, ps *ProgramState) {
				if ps.GetFactsCount() != 1 {
					t.Errorf("Expected 1 fact, got %d", ps.GetFactsCount())
				}
			},
		},
		{
			name:        "invalid syntax",
			content:     `this is not valid TSD syntax {{{}}}`,
			filename:    "invalid.tsd",
			expectError: true,
		},
		{
			name:        "empty content",
			content:     ``,
			filename:    "empty.tsd",
			expectError: true,
		},
		{
			name: "multiple types in content",
			content: `type Person(id: string)
type Company(id: string, name: string)`,
			filename: "multi.tsd",
			checkFunc: func(t *testing.T, ps *ProgramState) {
				if ps.GetTypesCount() != 2 {
					t.Errorf("Expected 2 types, got %d", ps.GetTypesCount())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := NewProgramState()
			err := ps.ParseAndMergeContent(tt.content, tt.filename)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.checkFunc != nil {
					tt.checkFunc(t, ps)
				}
			}
		})
	}
}

func TestProgramState_ErrorHandling(t *testing.T) {
	ps := NewProgramState()

	// Test initial state
	if ps.HasErrors() {
		t.Error("New ProgramState should have no errors")
	}
	if ps.GetErrorCount() != 0 {
		t.Errorf("Expected 0 errors, got %d", ps.GetErrorCount())
	}

	// Add an error
	ps.AddError(ValidationError{Message: "test error 1", File: "test.tsd", Type: "test"})
	if !ps.HasErrors() {
		t.Error("ProgramState should have errors after AddError")
	}
	if ps.GetErrorCount() != 1 {
		t.Errorf("Expected 1 error, got %d", ps.GetErrorCount())
	}

	// Add another error
	ps.AddError(ValidationError{Message: "test error 2", File: "test.tsd", Type: "test"})
	if ps.GetErrorCount() != 2 {
		t.Errorf("Expected 2 errors, got %d", ps.GetErrorCount())
	}

	// Get errors
	errors := ps.GetErrors()
	if len(errors) != 2 {
		t.Errorf("Expected 2 errors in list, got %d", len(errors))
	}

	// Clear errors
	ps.ClearErrors()
	if ps.HasErrors() {
		t.Error("ProgramState should have no errors after ClearErrors")
	}
	if ps.GetErrorCount() != 0 {
		t.Errorf("Expected 0 errors after clear, got %d", ps.GetErrorCount())
	}
}

func TestProgramState_Reset(t *testing.T) {
	ps := NewProgramState()

	// Add some data
	ps.ParseAndMergeContent(`type Person(id: string)`, "test.tsd")
	ps.ParseAndMergeContent(`Person(id: "P001")`, "test.tsd")
	ps.AddError(ValidationError{Message: "test error", File: "test.tsd", Type: "test"})

	// Verify data exists
	if ps.GetTypesCount() == 0 || ps.GetFactsCount() == 0 || !ps.HasErrors() {
		t.Fatal("Setup failed - expected some data in ProgramState")
	}

	// Reset
	ps.Reset()

	// Verify everything is cleared
	if ps.GetTypesCount() != 0 {
		t.Errorf("Expected 0 types after reset, got %d", ps.GetTypesCount())
	}
	if ps.GetRulesCount() != 0 {
		t.Errorf("Expected 0 rules after reset, got %d", ps.GetRulesCount())
	}
	if ps.GetFactsCount() != 0 {
		t.Errorf("Expected 0 facts after reset, got %d", ps.GetFactsCount())
	}
	if ps.HasErrors() {
		t.Error("Expected no errors after reset")
	}
}

func TestProgramState_MergeConflicts(t *testing.T) {
	ps := NewProgramState()

	// Parse first type definition
	content1 := `type Person(id: string, name: string)`
	err := ps.ParseAndMergeContent(content1, "file1.tsd")
	if err != nil {
		t.Fatalf("Failed to parse first content: %v", err)
	}

	// Try to parse conflicting type definition with same name
	content2 := `type Person(id: string, age: number)`
	err = ps.ParseAndMergeContent(content2, "file2.tsd")
	// This might succeed or fail depending on merge logic
	// Just verify no crash occurs
	if err != nil {
		t.Logf("Merge conflict detected (expected): %v", err)
	}

	// Verify at least one Person type exists
	if _, exists := ps.GetTypes()["Person"]; !exists {
		t.Error("Person type should exist after merge")
	}
}

func TestProgramState_ToProgram(t *testing.T) {
	ps := NewProgramState()

	// Add types and rules
	content := `type Person(id: string, age: number)
rule r1 : {p: Person} / p.age > 18 ==> Approve(p.id)
Person(id: "P001", age: 25)`

	err := ps.ParseAndMergeContent(content, "test.tsd")
	if err != nil {
		t.Fatalf("Failed to parse content: %v", err)
	}

	// Convert to Program
	program := ps.ToProgram()

	// Verify program structure
	if len(program.Types) != 1 {
		t.Errorf("Expected 1 type in program, got %d", len(program.Types))
	}
	if len(program.Expressions) != 1 {
		t.Errorf("Expected 1 expression in program, got %d", len(program.Expressions))
	}
	if len(program.Facts) != 1 {
		t.Errorf("Expected 1 fact in program, got %d", len(program.Facts))
	}
}

func TestProgramState_ParseNonExistentFile(t *testing.T) {
	ps := NewProgramState()

	// Try to parse a file that doesn't exist
	err := ps.ParseAndMerge("/nonexistent/path/to/file.tsd")
	if err == nil {
		t.Error("Expected error when parsing non-existent file")
	}
}
