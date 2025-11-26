// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRuleIdUniqueness tests that duplicate rule IDs are detected and rejected
func TestRuleIdUniqueness(t *testing.T) {
	t.Log("🧪 TEST: Rule ID Uniqueness Validation")
	t.Log("========================================")

	tempDir, err := os.MkdirTemp("", "rule_id_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create first file with a rule
	file1 := filepath.Join(tempDir, "rules1.constraint")
	content1 := `type Person : <id: string, age: number>

rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
rule r2 : {p: Person} / p.age < 18 ==> minor(p.id)`

	err = os.WriteFile(file1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("Failed to write file1: %v", err)
	}

	// Create second file with duplicate rule ID
	file2 := filepath.Join(tempDir, "rules2.constraint")
	content2 := `// This file has a duplicate rule ID
rule r1 : {p: Person} / p.age > 65 ==> senior(p.id)
rule r3 : {p: Person} / p.age == 18 ==> exactly_eighteen(p.id)`

	err = os.WriteFile(file2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("Failed to write file2: %v", err)
	}

	// Parse with ProgramState
	ps := NewProgramState()

	// Parse first file - should succeed
	err = ps.ParseAndMerge(file1)
	if err != nil {
		t.Fatalf("Failed to parse first file: %v", err)
	}

	// Verify we have 2 rules
	if len(ps.Rules) != 2 {
		t.Errorf("Expected 2 rules after first file, got %d", len(ps.Rules))
	}

	// Verify rule IDs are tracked
	if !ps.RuleIDs["r1"] {
		t.Error("Rule ID 'r1' should be tracked")
	}
	if !ps.RuleIDs["r2"] {
		t.Error("Rule ID 'r2' should be tracked")
	}

	t.Logf("✅ First file parsed: 2 rules (r1, r2)")

	// Parse second file - should detect duplicate r1
	err = ps.ParseAndMerge(file2)
	if err != nil {
		t.Fatalf("Failed to parse second file: %v", err)
	}

	// Should have only 3 rules total (r1 and r2 from file1, r3 from file2)
	// The duplicate r1 from file2 should be ignored
	if len(ps.Rules) != 3 {
		t.Errorf("Expected 3 rules total (duplicate r1 ignored), got %d", len(ps.Rules))
	}

	// Verify r3 was added
	if !ps.RuleIDs["r3"] {
		t.Error("Rule ID 'r3' should be tracked")
	}

	// Verify we have an error recorded for the duplicate
	foundDuplicateError := false
	for _, errRecord := range ps.Errors {
		if errRecord.Type == "rule" && errRecord.File == file2 {
			t.Logf("✅ Duplicate rule error recorded: %s", errRecord.Message)
			foundDuplicateError = true
			break
		}
	}

	if !foundDuplicateError {
		t.Error("❌ Expected an error to be recorded for duplicate rule ID")
	}

	t.Log("\n🎊 TEST PASSED: Duplicate rule IDs are detected and rejected")
}

// TestRuleIdUniquenessWithReset tests that rule IDs can be reused after a reset
func TestRuleIdUniquenessWithReset(t *testing.T) {
	t.Log("🧪 TEST: Rule ID Reuse After Reset")
	t.Log("===================================")

	tempDir, err := os.MkdirTemp("", "rule_id_reset_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create first file with rules
	file1 := filepath.Join(tempDir, "rules1.constraint")
	content1 := `type Person : <id: string, age: number>

rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
rule r2 : {p: Person} / p.age < 18 ==> minor(p.id)`

	err = os.WriteFile(file1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("Failed to write file1: %v", err)
	}

	// Create second file with reset and reused rule IDs
	file2 := filepath.Join(tempDir, "rules2.constraint")
	content2 := `reset

type Product : <id: string, price: number>

rule r1 : {prod: Product} / prod.price > 100 ==> expensive(prod.id)
rule r2 : {prod: Product} / prod.price < 50 ==> cheap(prod.id)`

	err = os.WriteFile(file2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("Failed to write file2: %v", err)
	}

	// Parse with ProgramState
	ps := NewProgramState()

	// Parse first file
	err = ps.ParseAndMerge(file1)
	if err != nil {
		t.Fatalf("Failed to parse first file: %v", err)
	}

	// Verify initial state
	if len(ps.Rules) != 2 {
		t.Errorf("Expected 2 rules after first file, got %d", len(ps.Rules))
	}
	if len(ps.Types) != 1 {
		t.Errorf("Expected 1 type after first file, got %d", len(ps.Types))
	}

	t.Logf("✅ Before reset: 2 rules (r1, r2), 1 type (Person)")

	// Parse second file with reset
	err = ps.ParseAndMerge(file2)
	if err != nil {
		t.Fatalf("Failed to parse second file with reset: %v", err)
	}

	// After reset, should have new rules with same IDs (allowed after reset)
	if len(ps.Rules) != 2 {
		t.Errorf("Expected 2 rules after reset, got %d", len(ps.Rules))
	}

	// Should have Product type, not Person (reset cleared Person)
	if len(ps.Types) != 1 {
		t.Errorf("Expected 1 type after reset, got %d", len(ps.Types))
	}

	if _, exists := ps.Types["Product"]; !exists {
		t.Error("Expected Product type after reset")
	}

	if _, exists := ps.Types["Person"]; exists {
		t.Error("Person type should be cleared after reset")
	}

	// Verify rule IDs are the new ones
	if !ps.RuleIDs["r1"] {
		t.Error("Rule ID 'r1' should be tracked after reset")
	}
	if !ps.RuleIDs["r2"] {
		t.Error("Rule ID 'r2' should be tracked after reset")
	}

	// Should have no errors (reusing IDs after reset is allowed)
	if len(ps.Errors) > 0 {
		t.Errorf("Expected no errors after reset, got %d: %v", len(ps.Errors), ps.Errors)
	}

	t.Logf("✅ After reset: 2 rules (r1, r2 reused), 1 type (Product)")
	t.Log("\n🎊 TEST PASSED: Rule IDs can be reused after reset")
}

// TestRuleIdUniquenessInSameFile tests that duplicate IDs in the same file are detected
func TestRuleIdUniquenessInSameFile(t *testing.T) {
	t.Log("🧪 TEST: Duplicate Rule IDs in Same File")
	t.Log("=========================================")

	tempDir, err := os.MkdirTemp("", "rule_id_same_file_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create file with duplicate rule IDs
	file := filepath.Join(tempDir, "duplicate_rules.constraint")
	content := `type Person : <id: string, age: number>

rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
rule r2 : {p: Person} / p.age < 18 ==> minor(p.id)
rule r1 : {p: Person} / p.age == 18 ==> exactly_eighteen(p.id)`

	err = os.WriteFile(file, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Parse with ProgramState
	ps := NewProgramState()

	err = ps.ParseAndMerge(file)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	// Should have only 2 rules (first r1 and r2)
	// The second r1 should be ignored
	if len(ps.Rules) != 2 {
		t.Errorf("Expected 2 rules (duplicate ignored), got %d", len(ps.Rules))
	}

	// Verify we have an error recorded
	if len(ps.Errors) != 1 {
		t.Errorf("Expected 1 error for duplicate rule ID, got %d", len(ps.Errors))
	}

	if len(ps.Errors) > 0 {
		t.Logf("✅ Duplicate error recorded: %s", ps.Errors[0].Message)
	}

	t.Log("\n🎊 TEST PASSED: Duplicate rule IDs in same file are detected")
}

// TestRuleIdEmptyAllowed tests that empty rule IDs are allowed (for backward compatibility testing)
func TestRuleIdEmptyAllowed(t *testing.T) {
	t.Log("🧪 TEST: Empty Rule IDs Handling")
	t.Log("=================================")

	ps := NewProgramState()

	// Manually create rules with empty IDs (simulating legacy or edge case)
	ps.Types["Person"] = &TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "age", Type: "number"},
		},
	}

	// Create rules with empty IDs
	rule1 := Expression{
		Type:   "expression",
		RuleId: "", // Empty ID
		Set: Set{
			Type: "set",
			Variables: []TypedVariable{
				{Type: "typedVariable", Name: "p", DataType: "Person"},
			},
		},
	}

	rule2 := Expression{
		Type:   "expression",
		RuleId: "", // Empty ID (duplicate empty)
		Set: Set{
			Type: "set",
			Variables: []TypedVariable{
				{Type: "typedVariable", Name: "p", DataType: "Person"},
			},
		},
	}

	// Merge rules
	err := ps.mergeRules([]Expression{rule1, rule2}, "test.constraint")
	if err != nil {
		t.Fatalf("Failed to merge rules: %v", err)
	}

	// Both rules should be added (empty IDs are not tracked)
	if len(ps.Rules) != 2 {
		t.Errorf("Expected 2 rules with empty IDs, got %d", len(ps.Rules))
	}

	// Should have no errors (empty IDs don't trigger uniqueness check)
	if len(ps.Errors) > 0 {
		t.Errorf("Expected no errors for empty IDs, got %d", len(ps.Errors))
	}

	t.Logf("✅ Both rules with empty IDs were accepted")
	t.Log("\n🎊 TEST PASSED: Empty rule IDs are handled correctly")
}

// TestRuleIdMultipleFiles tests uniqueness across multiple files
func TestRuleIdMultipleFiles(t *testing.T) {
	t.Log("🧪 TEST: Rule ID Uniqueness Across Multiple Files")
	t.Log("===================================================")

	tempDir, err := os.MkdirTemp("", "rule_id_multi_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create three files
	file1 := filepath.Join(tempDir, "file1.constraint")
	content1 := `type Person : <id: string, age: number>
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)`

	file2 := filepath.Join(tempDir, "file2.constraint")
	content2 := `rule r2 : {p: Person} / p.age < 18 ==> minor(p.id)`

	file3 := filepath.Join(tempDir, "file3.constraint")
	content3 := `rule r1 : {p: Person} / p.age == 18 ==> exactly_eighteen(p.id)
rule r3 : {p: Person} / p.age > 65 ==> senior(p.id)`

	for _, fileData := range []struct {
		path    string
		content string
	}{
		{file1, content1},
		{file2, content2},
		{file3, content3},
	} {
		err = os.WriteFile(fileData.path, []byte(fileData.content), 0644)
		if err != nil {
			t.Fatalf("Failed to write %s: %v", fileData.path, err)
		}
	}

	ps := NewProgramState()

	// Parse file1
	err = ps.ParseAndMerge(file1)
	if err != nil {
		t.Fatalf("Failed to parse file1: %v", err)
	}
	t.Logf("✅ File1: 1 rule (r1)")

	// Parse file2
	err = ps.ParseAndMerge(file2)
	if err != nil {
		t.Fatalf("Failed to parse file2: %v", err)
	}
	t.Logf("✅ File2: 1 rule (r2)")

	// Parse file3 - should detect duplicate r1
	err = ps.ParseAndMerge(file3)
	if err != nil {
		t.Fatalf("Failed to parse file3: %v", err)
	}
	t.Logf("✅ File3: 1 rule accepted (r3), 1 duplicate rejected (r1)")

	// Should have 3 rules total: r1 (file1), r2 (file2), r3 (file3)
	if len(ps.Rules) != 3 {
		t.Errorf("Expected 3 rules total, got %d", len(ps.Rules))
	}

	// Should have 1 error for duplicate r1 in file3
	duplicateErrors := 0
	for _, errRecord := range ps.Errors {
		if errRecord.Type == "rule" && errRecord.File == file3 {
			duplicateErrors++
			t.Logf("✅ Error recorded: %s", errRecord.Message)
		}
	}

	if duplicateErrors != 1 {
		t.Errorf("Expected 1 duplicate error, got %d", duplicateErrors)
	}

	// Verify tracked IDs
	expectedIDs := []string{"r1", "r2", "r3"}
	for _, id := range expectedIDs {
		if !ps.RuleIDs[id] {
			t.Errorf("Rule ID '%s' should be tracked", id)
		}
	}

	t.Log("\n🎊 TEST PASSED: Rule ID uniqueness enforced across multiple files")
}
