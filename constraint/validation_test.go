// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"testing"
)

// TestValidateFactWithInvalidType tests validateFact with an undefined type
func TestValidateFactWithInvalidType(t *testing.T) {
	ps := NewProgramState()

	// Add a type definition
	ps.Types["Person"] = &TypeDefinition{
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "identifier"},
			{Name: "name", Type: "string"},
		},
	}

	// Create a fact with an undefined type
	fact := &Fact{
		TypeName: "UnknownType",
		Fields: []FactField{
			{Name: "id", Value: FactValue{Type: "identifier", Value: "U001"}},
		},
	}

	err := ps.validateFact(fact, "test.constraint")
	if err == nil {
		t.Error("Expected error for undefined type, got nil")
	}
	if err != nil && err.Error() != "fact references undefined type UnknownType in test.constraint" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// TestValidateFactWithInvalidField tests validateFact with an undefined field
func TestValidateFactWithInvalidField(t *testing.T) {
	ps := NewProgramState()

	// Add a type definition
	ps.Types["Person"] = &TypeDefinition{
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "identifier"},
			{Name: "name", Type: "string"},
		},
	}

	// Create a fact with an undefined field
	fact := &Fact{
		TypeName: "Person",
		Fields: []FactField{
			{Name: "id", Value: FactValue{Type: "identifier", Value: "U001"}},
			{Name: "invalidField", Value: FactValue{Type: "string", Value: "test"}},
		},
	}

	err := ps.validateFact(fact, "test.constraint")
	if err == nil {
		t.Error("Expected error for undefined field, got nil")
	}
}

// TestValidateFactSuccess tests validateFact with a valid fact
func TestValidateFactSuccess(t *testing.T) {
	ps := NewProgramState()

	// Add a type definition
	ps.Types["Person"] = &TypeDefinition{
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "identifier"},
			{Name: "name", Type: "string"},
		},
	}

	// Create a valid fact
	fact := &Fact{
		TypeName: "Person",
		Fields: []FactField{
			{Name: "id", Value: FactValue{Type: "identifier", Value: "U001"}},
			{Name: "name", Value: FactValue{Type: "string", Value: "John"}},
		},
	}

	err := ps.validateFact(fact, "test.constraint")
	if err != nil {
		t.Errorf("Unexpected error for valid fact: %v", err)
	}
}

// TestMergeTypesNilProgramState tests mergeTypes with nil ProgramState
func TestMergeTypesNilProgramState(t *testing.T) {
	var ps *ProgramState = nil

	types := []TypeDefinition{
		{Name: "Person", Fields: []Field{{Name: "id", Type: "identifier"}}},
	}

	err := ps.mergeTypes(types, "test.constraint")
	if err == nil {
		t.Error("Expected error for nil ProgramState, got nil")
	}
}

// TestParseAndMergeContentEmptyContent tests ParseAndMergeContent with empty content
func TestParseAndMergeContentEmptyContent(t *testing.T) {
	ps := NewProgramState()

	err := ps.ParseAndMergeContent("", "test.constraint")
	if err == nil {
		t.Error("Expected error for empty content, got nil")
	}
	if err != nil && err.Error() != "content cannot be empty" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// TestParseAndMergeContentEmptyFilename tests ParseAndMergeContent with empty filename
func TestParseAndMergeContentEmptyFilename(t *testing.T) {
	ps := NewProgramState()

	err := ps.ParseAndMergeContent("type Person { id: identifier }", "")
	if err == nil {
		t.Error("Expected error for empty filename, got nil")
	}
	if err != nil && err.Error() != "filename cannot be empty" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// TestAreTypesCompatible tests type compatibility checking
func TestAreTypesCompatible(t *testing.T) {
	ps := NewProgramState()

	// Test 1: Identical types should be compatible
	type1 := &TypeDefinition{
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "identifier"},
			{Name: "name", Type: "string"},
		},
	}

	type2 := &TypeDefinition{
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "identifier"},
			{Name: "name", Type: "string"},
		},
	}

	if !ps.areTypesCompatible(type1, type2) {
		t.Error("Expected identical types to be compatible")
	}

	// Test 2: Types with different field types should be incompatible
	type3 := &TypeDefinition{
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"}, // Different type for id
			{Name: "name", Type: "string"},
		},
	}

	if ps.areTypesCompatible(type1, type3) {
		t.Error("Expected types with different field types to be incompatible")
	}

	// Test 3: Different names should be incompatible
	type4 := &TypeDefinition{
		Name: "User",
		Fields: []Field{
			{Name: "id", Type: "identifier"},
			{Name: "name", Type: "string"},
		},
	}

	if ps.areTypesCompatible(type1, type4) {
		t.Error("Expected types with different names to be incompatible")
	}
}

// TestErrorMethods tests ValidationError methods
func TestErrorMethods(t *testing.T) {
	ps := NewProgramState()

	err := ValidationError{
		File:    "test.constraint",
		Type:    ErrorTypeFact,
		Message: "test error",
		Line:    10,
	}

	ps.AddError(err)

	if !ps.HasErrors() {
		t.Error("Expected HasErrors to return true")
	}

	if ps.GetErrorCount() != 1 {
		t.Errorf("Expected 1 error, got %d", ps.GetErrorCount())
	}

	errors := ps.GetErrors()
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}

	// Test that GetErrors returns a copy
	errors[0].Message = "modified"
	if ps.Errors[0].Message == "modified" {
		t.Error("GetErrors should return a copy, not the original slice")
	}

	ps.ClearErrors()
	if ps.HasErrors() {
		t.Error("Expected HasErrors to return false after ClearErrors")
	}
}
