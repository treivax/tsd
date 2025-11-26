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
	typesFile := filepath.Join(tempDir, "types.constraint")
	typesContent := `// Type definitions
type Person : <id: string, name: string, age: number>
type Company : <id: string, name: string, employees: number>`

	err = os.WriteFile(typesFile, []byte(typesContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write types file: %v", err)
	}

	// 2. Rules file
	rulesFile := filepath.Join(tempDir, "rules.constraint")
	rulesContent := `// Business rules
rule r1 : {p: Person, c: Company} / p.age > 25 AND c.employees > 100 ==> AddToWorkforce(p.id, c.id)`

	err = os.WriteFile(rulesFile, []byte(rulesContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write rules file: %v", err)
	}

	// 3. Facts file
	factsFile := filepath.Join(tempDir, "facts.constraint")
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
	if len(ps.Types) != 2 {
		t.Errorf("Expected 2 types, got %d", len(ps.Types))
	}

	if _, exists := ps.Types["Person"]; !exists {
		t.Error("Person type not found")
	}

	if _, exists := ps.Types["Company"]; !exists {
		t.Error("Company type not found")
	}

	// Parse rules (should validate against existing types)
	err = ps.ParseAndMerge(rulesFile)
	if err != nil {
		t.Fatalf("Failed to parse rules file: %v", err)
	}

	// Verify rules were parsed
	if len(ps.Rules) == 0 {
		t.Error("No rules were parsed")
	}

	// Parse facts (should validate against existing types)
	err = ps.ParseAndMerge(factsFile)
	if err != nil {
		t.Fatalf("Failed to parse facts file: %v", err)
	}

	// Verify facts were parsed
	if len(ps.Facts) != 3 {
		t.Errorf("Expected 3 facts, got %d", len(ps.Facts))
	}

	// Verify all files were recorded
	if len(ps.FilesParsed) != 3 {
		t.Errorf("Expected 3 parsed files, got %d", len(ps.FilesParsed))
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
		len(ps.FilesParsed), len(ps.Types), len(ps.Rules), len(ps.Facts))
}

func TestProgramStateTypeValidation(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "program_state_validation_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create types file
	typesFile := filepath.Join(tempDir, "types.constraint")
	typesContent := `type Person : <id: string, name: string>`
	err = os.WriteFile(typesFile, []byte(typesContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write types file: %v", err)
	}

	// Create facts file with invalid type
	invalidFactsFile := filepath.Join(tempDir, "invalid_facts.constraint")
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
