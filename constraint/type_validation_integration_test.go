// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

// TestRuleValidation_UndefinedType tests that a rule referencing an undefined type
// produces a non-blocking error and is rejected
func TestRuleValidation_UndefinedType(t *testing.T) {
	ps := NewProgramState()

	content := `
type Person(id: string, name: string,age:number)

rule r1 : {u: UnknownType} / u.name == "test" ==> log(u.id)
`

	err := ps.ParseAndMergeContent(content, "test.tsd")
	if err != nil {
		t.Fatalf("ParseAndMergeContent should not return blocking error: %v", err)
	}

	// Check that the error was recorded
	if !ps.HasErrors() {
		t.Error("Expected non-blocking error to be recorded")
	}

	errors := ps.GetErrors()
	if len(errors) == 0 {
		t.Error("Expected at least one error")
	}

	found := false
	for _, err := range errors {
		if err.Type == "rule" && strings.Contains(err.Message, "UnknownType") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected error about undefined type UnknownType")
	}

	// Check that the rule was rejected (not added)
	if ps.GetRulesCount() != 0 {
		t.Errorf("Expected 0 rules (invalid rule should be rejected), got %d", ps.GetRulesCount())
	}

	// Check that valid type was still added
	if ps.GetTypesCount() != 1 {
		t.Errorf("Expected 1 type definition, got %d", ps.GetTypesCount())
	}
}

// TestRuleValidation_InvalidFieldAccess tests that a rule accessing a non-existent field
// produces a non-blocking error and is rejected
func TestRuleValidation_InvalidFieldAccess(t *testing.T) {
	ps := NewProgramState()

	content := `
type Person(id: string, name: string,age:number)

rule r2 : {p: Person} / p.salary > 1000 ==> high_earner(p.id)
`

	err := ps.ParseAndMergeContent(content, "test.tsd")
	if err != nil {
		t.Fatalf("ParseAndMergeContent should not return blocking error: %v", err)
	}

	// Check that the error was recorded
	if !ps.HasErrors() {
		t.Error("Expected non-blocking error to be recorded")
	}

	errors := ps.GetErrors()
	found := false
	for _, err := range errors {
		if err.Type == "rule" && strings.Contains(err.Message, "salary") {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected error about invalid field 'salary', got: %v", errors)
	}

	// Check that the rule was rejected
	if ps.GetRulesCount() != 0 {
		t.Errorf("Expected 0 rules (invalid rule should be rejected), got %d", ps.GetRulesCount())
	}
}

// TestRuleValidation_ValidRule tests that a valid rule is accepted
func TestRuleValidation_ValidRule(t *testing.T) {
	ps := NewProgramState()

	content := `
type Person(id: string, name: string,age:number)

rule r3 : {p: Person} / p.age > 18 ==> adult(p.id, p.name)
`

	err := ps.ParseAndMergeContent(content, "test.tsd")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if ps.HasErrors() {
		t.Errorf("Expected no errors, got: %v", ps.GetErrors())
	}

	if ps.GetRulesCount() != 1 {
		t.Errorf("Expected 1 rule, got %d", ps.GetRulesCount())
	}
}

// TestFactValidation_UndefinedType tests that a fact referencing an undefined type
// produces a non-blocking error and is rejected
func TestFactValidation_UndefinedType(t *testing.T) {
	ps := NewProgramState()

	content := `
type Person(id: string, name: string, age: number)

UnknownType(id: "X001", data: "test")
`

	err := ps.ParseAndMergeContent(content, "test.tsd")
	if err != nil {
		t.Fatalf("ParseAndMergeContent should not return blocking error: %v", err)
	}

	// Check that the error was recorded
	if !ps.HasErrors() {
		t.Error("Expected non-blocking error to be recorded")
	}

	errors := ps.GetErrors()
	found := false
	for _, err := range errors {
		if err.Type == "fact" && strings.Contains(err.Message, "UnknownType") {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected error about undefined type UnknownType, got: %v", errors)
	}

	// Check that the fact was rejected
	if ps.GetFactsCount() != 0 {
		t.Errorf("Expected 0 facts (invalid fact should be rejected), got %d", ps.GetFactsCount())
	}

	// Check that valid type was still added
	if ps.GetTypesCount() != 1 {
		t.Errorf("Expected 1 type definition, got %d", ps.GetTypesCount())
	}
}

// TestFactValidation_InvalidField tests that a fact with an undefined field
// produces a non-blocking error and is rejected
func TestFactValidation_InvalidField(t *testing.T) {
	ps := NewProgramState()

	content := `
type Person(id: string, name: string, age: number)

Person(id: "P001", name: "Alice", salary: 50000)
`

	err := ps.ParseAndMergeContent(content, "test.tsd")
	if err != nil {
		t.Fatalf("ParseAndMergeContent should not return blocking error: %v", err)
	}

	// Check that the error was recorded
	if !ps.HasErrors() {
		t.Error("Expected non-blocking error to be recorded")
	}

	errors := ps.GetErrors()
	found := false
	for _, err := range errors {
		if err.Type == "fact" && strings.Contains(err.Message, "salary") {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected error about undefined field 'salary', got: %v", errors)
	}

	// Check that the fact was rejected
	if ps.GetFactsCount() != 0 {
		t.Errorf("Expected 0 facts (invalid fact should be rejected), got %d", ps.GetFactsCount())
	}
}

// TestFactValidation_WrongFieldType tests that a fact with wrong field type
// produces a non-blocking error and is rejected
func TestFactValidation_WrongFieldType(t *testing.T) {
	ps := NewProgramState()

	content := `
type Person(id: string, name: string, age: number)

Person(id: "P001", name: "Alice", age: "twenty-five")
`

	err := ps.ParseAndMergeContent(content, "test.tsd")
	if err != nil {
		t.Fatalf("ParseAndMergeContent should not return blocking error: %v", err)
	}

	// Check that the error was recorded
	if !ps.HasErrors() {
		t.Error("Expected non-blocking error to be recorded")
	}

	errors := ps.GetErrors()
	found := false
	for _, err := range errors {
		if err.Type == "fact" && (strings.Contains(err.Message, "expected number") || strings.Contains(err.Message, "age")) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected error about wrong type for 'age', got: %v", errors)
	}

	// Check that the fact was rejected
	if ps.GetFactsCount() != 0 {
		t.Errorf("Expected 0 facts (invalid fact should be rejected), got %d", ps.GetFactsCount())
	}
}

// TestFactValidation_ValidFact tests that a valid fact is accepted
func TestFactValidation_ValidFact(t *testing.T) {
	ps := NewProgramState()

	content := `
type Person(id: string, name: string, age: number)

Person(id: "P001", name: "Alice", age: 25)
`

	err := ps.ParseAndMergeContent(content, "test.tsd")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if ps.HasErrors() {
		t.Errorf("Expected no errors, got: %v", ps.GetErrors())
	}

	if ps.GetFactsCount() != 1 {
		t.Errorf("Expected 1 fact, got %d", ps.GetFactsCount())
	}
}

// TestMixedValidation_PartialSuccess tests that valid items are accepted
// while invalid items are rejected with non-blocking errors
func TestMixedValidation_PartialSuccess(t *testing.T) {
	ps := NewProgramState()

	content := `
type Person(id: string, name: string,age:number)
type Product(id: string, name: string,price:number)

Person(id: "P001", name: "Alice", age: 25)
Person(id: "P002", name: "Bob", invalidField: "test")
Product(id: "PR001", name: "Widget", price: 19.99)

rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
rule r2 : {x: UnknownType} / x.value > 0 ==> process(x.id)
rule r3 : {pr: Product} / pr.price < 100 ==> affordable(pr.id)
`

	err := ps.ParseAndMergeContent(content, "test.tsd")
	if err != nil {
		t.Fatalf("ParseAndMergeContent should not return blocking error: %v", err)
	}

	// Check that errors were recorded for invalid items
	if !ps.HasErrors() {
		t.Error("Expected non-blocking errors to be recorded")
	}

	errorCount := ps.GetErrorCount()
	if errorCount != 2 {
		t.Errorf("Expected 2 errors (1 invalid fact, 1 invalid rule), got %d", errorCount)
	}

	// Check that valid items were still added
	if ps.GetTypesCount() != 2 {
		t.Errorf("Expected 2 types, got %d", ps.GetTypesCount())
	}

	if ps.GetFactsCount() != 2 {
		t.Errorf("Expected 2 valid facts (out of 3), got %d", ps.GetFactsCount())
	}

	if ps.GetRulesCount() != 2 {
		t.Errorf("Expected 2 valid rules (out of 3), got %d", ps.GetRulesCount())
	}
}

// TestMultipleFiles_ValidationAcrossFiles tests that validation works correctly
// across multiple file parsing operations
func TestMultipleFiles_ValidationAcrossFiles(t *testing.T) {
	ps := NewProgramState()

	// File 1: Define types
	content1 := `
type Person(id: string, name: string, age: number)
`
	err := ps.ParseAndMergeContent(content1, "types.tsd")
	if err != nil {
		t.Fatalf("Unexpected error in file 1: %v", err)
	}

	// File 2: Add valid and invalid facts
	content2 := `
Person(id: "P001", name: "Alice", age: 25)
Person(id: "P002", name: "Bob", salary: 50000)
`
	err = ps.ParseAndMergeContent(content2, "facts.tsd")
	if err != nil {
		t.Fatalf("ParseAndMergeContent should not return blocking error: %v", err)
	}

	// File 3: Add valid and invalid rules
	content3 := `
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
rule r2 : {p: Person} / p.salary > 0 ==> employed(p.id)
`
	err = ps.ParseAndMergeContent(content3, "rules.tsd")
	if err != nil {
		t.Fatalf("ParseAndMergeContent should not return blocking error: %v", err)
	}

	// Check results
	if !ps.HasErrors() {
		t.Error("Expected errors to be recorded")
	}

	if ps.GetErrorCount() != 2 {
		t.Errorf("Expected 2 errors, got %d: %v", ps.GetErrorCount(), ps.GetErrors())
	}

	if ps.GetTypesCount() != 1 {
		t.Errorf("Expected 1 type, got %d", ps.GetTypesCount())
	}

	if ps.GetFactsCount() != 1 {
		t.Errorf("Expected 1 valid fact, got %d", ps.GetFactsCount())
	}

	if ps.GetRulesCount() != 1 {
		t.Errorf("Expected 1 valid rule, got %d", ps.GetRulesCount())
	}

	// Verify error details
	errors := ps.GetErrors()
	factErrorFound := false
	ruleErrorFound := false

	for _, err := range errors {
		if err.Type == "fact" && err.File == "facts.tsd" {
			factErrorFound = true
		}
		if err.Type == "rule" && err.File == "rules.tsd" {
			ruleErrorFound = true
		}
	}

	if !factErrorFound {
		t.Error("Expected fact error from facts.tsd")
	}
	if !ruleErrorFound {
		t.Error("Expected rule error from rules.tsd")
	}
}

// TestValidation_ComplexFieldAccess tests validation of nested field accesses
func TestValidation_ComplexFieldAccess(t *testing.T) {
	ps := NewProgramState()

	content := `
type Person(id: string, name: string,age:number)

rule r1 : {p: Person} / p.age > 18 AND p.name == "Alice" ==> process(p.id)
rule r2 : {p: Person} / p.age > 18 AND p.salary > 1000 ==> high_earner(p.id)
`

	err := ps.ParseAndMergeContent(content, "test.tsd")
	if err != nil {
		t.Fatalf("ParseAndMergeContent should not return blocking error: %v", err)
	}

	// First rule should be valid
	// Second rule should be invalid (salary field doesn't exist)
	if !ps.HasErrors() {
		t.Error("Expected error for invalid field access")
	}

	if ps.GetRulesCount() != 1 {
		t.Errorf("Expected 1 valid rule, got %d", ps.GetRulesCount())
	}

	// Check that the valid rule has ID "r1"
	if ps.GetRules()[0].RuleId != "r1" {
		t.Errorf("Expected valid rule to be r1, got %s", ps.GetRules()[0].RuleId)
	}
}

// TestValidation_ActionFieldAccess tests that field accesses in actions are validated
func TestValidation_ActionFieldAccess(t *testing.T) {
	ps := NewProgramState()

	content := `
type Person(id: string, name: string,age:number)

rule r1 : {p: Person} / p.age > 18 ==> process(p.id, p.name)
rule r2 : {p: Person} / p.age > 18 ==> process(p.id, p.address)
`

	err := ps.ParseAndMergeContent(content, "test.tsd")
	if err != nil {
		t.Fatalf("ParseAndMergeContent should not return blocking error: %v", err)
	}

	// First rule should be valid
	// Second rule should be invalid (address field doesn't exist in action)
	if !ps.HasErrors() {
		t.Error("Expected error for invalid field access in action")
	}

	if ps.GetRulesCount() != 1 {
		t.Errorf("Expected 1 valid rule, got %d", ps.GetRulesCount())
	}

	// Check that the valid rule has ID "r1"
	if ps.GetRules()[0].RuleId != "r1" {
		t.Errorf("Expected valid rule to be r1, got %s", ps.GetRules()[0].RuleId)
	}
}
