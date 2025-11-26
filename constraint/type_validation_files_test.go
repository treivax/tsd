// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

// TestTypeValidation_ValidFile tests that a file with only valid items passes validation
func TestTypeValidation_ValidFile(t *testing.T) {
	ps := NewProgramState()

	err := ps.ParseAndMerge("test/validation/valid_items.tsd")
	if err != nil {
		t.Fatalf("Unexpected error parsing valid file: %v", err)
	}

	if ps.HasErrors() {
		t.Errorf("Valid file should not produce errors, got: %v", ps.GetErrors())
	}

	// Verify all types were loaded
	if len(ps.Types) != 3 {
		t.Errorf("Expected 3 types, got %d", len(ps.Types))
	}

	// Verify all facts were loaded
	if len(ps.Facts) != 8 {
		t.Errorf("Expected 8 facts, got %d", len(ps.Facts))
	}

	// Verify all rules were loaded
	if len(ps.Rules) != 5 {
		t.Errorf("Expected 5 rules, got %d", len(ps.Rules))
	}
}

// TestTypeValidation_InvalidFile tests that a file with invalid items produces appropriate errors
func TestTypeValidation_InvalidFile(t *testing.T) {
	ps := NewProgramState()

	err := ps.ParseAndMerge("test/validation/invalid_items.tsd")
	if err != nil {
		t.Fatalf("Should not return blocking error: %v", err)
	}

	// Should have validation errors
	if !ps.HasErrors() {
		t.Fatal("Expected validation errors for invalid items")
	}

	errorCount := ps.GetErrorCount()
	expectedErrors := 9 // 4 invalid facts + 5 invalid rules

	if errorCount != expectedErrors {
		t.Errorf("Expected %d errors, got %d", expectedErrors, errorCount)
		for i, err := range ps.GetErrors() {
			t.Logf("Error %d: [%s] %s", i+1, err.Type, err.Message)
		}
	}

	// Verify valid items were still loaded
	if len(ps.Types) != 2 {
		t.Errorf("Expected 2 types (valid ones), got %d", len(ps.Types))
	}

	if len(ps.Facts) != 2 {
		t.Errorf("Expected 2 valid facts, got %d", len(ps.Facts))
	}

	if len(ps.Rules) != 2 {
		t.Errorf("Expected 2 valid rules, got %d", len(ps.Rules))
	}

	// Verify error types
	factErrors := 0
	ruleErrors := 0

	for _, err := range ps.GetErrors() {
		if err.Type == "fact" {
			factErrors++
		} else if err.Type == "rule" {
			ruleErrors++
		}
	}

	if factErrors != 4 {
		t.Errorf("Expected 4 fact errors, got %d", factErrors)
	}

	if ruleErrors != 5 {
		t.Errorf("Expected 5 rule errors, got %d", ruleErrors)
	}
}

// TestTypeValidation_ErrorMessages tests that error messages are descriptive
func TestTypeValidation_ErrorMessages(t *testing.T) {
	ps := NewProgramState()

	err := ps.ParseAndMerge("test/validation/invalid_items.tsd")
	if err != nil {
		t.Fatalf("Should not return blocking error: %v", err)
	}

	errors := ps.GetErrors()

	// Check for specific error patterns
	errorPatterns := map[string]bool{
		"undefined type":    false,
		"undefined field":   false,
		"not found in type": false,
		"expected number":   false,
	}

	for _, err := range errors {
		for pattern := range errorPatterns {
			if strings.Contains(err.Message, pattern) {
				errorPatterns[pattern] = true
			}
		}
	}

	for pattern, found := range errorPatterns {
		if !found {
			t.Errorf("Expected to find error pattern '%s' in error messages", pattern)
		}
	}
}

// TestTypeValidation_FileTracking tests that errors correctly track source files
func TestTypeValidation_FileTracking(t *testing.T) {
	ps := NewProgramState()

	err := ps.ParseAndMerge("test/validation/invalid_items.tsd")
	if err != nil {
		t.Fatalf("Should not return blocking error: %v", err)
	}

	errors := ps.GetErrors()

	for _, err := range errors {
		if err.File != "test/validation/invalid_items.tsd" {
			t.Errorf("Expected error file to be 'test/validation/invalid_items.tsd', got '%s'", err.File)
		}
	}
}

// TestTypeValidation_MultipleFiles tests validation across multiple files
func TestTypeValidation_MultipleFiles(t *testing.T) {
	ps := NewProgramState()

	// Parse valid file first
	err := ps.ParseAndMerge("test/validation/valid_items.tsd")
	if err != nil {
		t.Fatalf("Unexpected error parsing valid file: %v", err)
	}

	if ps.HasErrors() {
		t.Errorf("Valid file should not produce errors")
	}

	initialFactsCount := len(ps.Facts)
	initialRulesCount := len(ps.Rules)

	// Parse invalid file second
	err = ps.ParseAndMerge("test/validation/invalid_items.tsd")
	if err != nil {
		t.Fatalf("Should not return blocking error: %v", err)
	}

	// Should have errors from invalid file
	if !ps.HasErrors() {
		t.Error("Expected validation errors from invalid file")
	}

	// Facts should increase by 2 (valid facts from invalid_items.tsd)
	expectedFacts := initialFactsCount + 2
	if len(ps.Facts) != expectedFacts {
		t.Errorf("Expected %d facts total, got %d", expectedFacts, len(ps.Facts))
	}

	// At least some rules should have been added or rejected
	if len(ps.Rules) < initialRulesCount {
		t.Errorf("Rule count should not decrease, got %d (was %d)", len(ps.Rules), initialRulesCount)
	}

	// All errors should be from the invalid file
	for _, err := range ps.GetErrors() {
		if err.File != "test/validation/invalid_items.tsd" {
			t.Errorf("Expected all errors from invalid_items.tsd, got error from %s", err.File)
		}
	}

	// Should have multiple errors in the second file
	if ps.GetErrorCount() < 5 {
		t.Errorf("Expected at least 5 errors, got %d", ps.GetErrorCount())
	}
}
