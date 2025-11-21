package constraint

import (
	"testing"
)

// TestValidateRuleWithUndefinedType tests validateRule with undefined variable type
func TestValidateRuleWithUndefinedType(t *testing.T) {
	ps := NewProgramState()

	ps.Types["Person"] = &TypeDefinition{
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "identifier"},
			{Name: "name", Type: "string"},
		},
	}

	rule := &Expression{
		Set: Set{
			Variables: []TypedVariable{
				{Name: "p", DataType: "UnknownType"},
			},
		},
		Constraints: map[string]interface{}{},
	}

	err := ps.validateRule(rule, "test.constraint")
	if err == nil {
		t.Error("Expected error for undefined type in rule, got nil")
	}
}

// TestValidateRuleWithValidType tests validateRule with valid types
func TestValidateRuleWithValidType(t *testing.T) {
	ps := NewProgramState()

	ps.Types["Person"] = &TypeDefinition{
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "identifier"},
			{Name: "name", Type: "string"},
		},
	}

	rule := &Expression{
		Set: Set{
			Variables: []TypedVariable{
				{Name: "p", DataType: "Person"},
			},
		},
		Constraints: map[string]interface{}{},
	}

	err := ps.validateRule(rule, "test.constraint")
	if err != nil {
		t.Errorf("Unexpected error for valid rule: %v", err)
	}
}

// TestScanForFieldAccessWithInvalidField tests field access validation
func TestScanForFieldAccessWithInvalidField(t *testing.T) {
	ps := NewProgramState()

	ps.Types["Person"] = &TypeDefinition{
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "identifier"},
			{Name: "name", Type: "string"},
		},
	}

	variables := map[string]string{
		"p": "Person",
	}

	data := map[string]interface{}{
		"type":   "fieldAccess",
		"object": "p",
		"field":  "invalidField",
	}

	err := ps.scanForFieldAccess(data, variables)
	if err == nil {
		t.Error("Expected error for invalid field access, got nil")
	}
}

// TestScanForFieldAccessWithValidField tests valid field access
func TestScanForFieldAccessWithValidField(t *testing.T) {
	ps := NewProgramState()

	ps.Types["Person"] = &TypeDefinition{
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "identifier"},
			{Name: "name", Type: "string"},
		},
	}

	variables := map[string]string{
		"p": "Person",
	}

	data := map[string]interface{}{
		"type":   "fieldAccess",
		"object": "p",
		"field":  "name",
	}

	err := ps.scanForFieldAccess(data, variables)
	if err != nil {
		t.Errorf("Unexpected error for valid field access: %v", err)
	}
}

// TestPrintErrors tests PrintErrors method
func TestPrintErrors(t *testing.T) {
	ps := NewProgramState()

	ps.AddError(ValidationError{
		File:    "test.constraint",
		Type:    ErrorTypeFact,
		Message: "test error",
		Line:    10,
	})

	ps.PrintErrors()
}

// TestMergeTypesCompatible tests merging compatible types
func TestMergeTypesCompatible(t *testing.T) {
	ps := NewProgramState()

	types1 := []TypeDefinition{
		{
			Name: "Person",
			Fields: []Field{
				{Name: "id", Type: "identifier"},
			},
		},
	}

	err := ps.mergeTypes(types1, "file1.constraint")
	if err != nil {
		t.Errorf("Unexpected error merging types: %v", err)
	}

	types2 := []TypeDefinition{
		{
			Name: "Person",
			Fields: []Field{
				{Name: "id", Type: "identifier"},
				{Name: "name", Type: "string"},
			},
		},
	}

	err = ps.mergeTypes(types2, "file2.constraint")
	if err != nil {
		t.Errorf("Unexpected error merging compatible types: %v", err)
	}

	if len(ps.Types["Person"].Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(ps.Types["Person"].Fields))
	}
}

// TestMergeTypesIncompatible tests merging incompatible types
func TestMergeTypesIncompatible(t *testing.T) {
	ps := NewProgramState()

	types1 := []TypeDefinition{
		{
			Name: "Person",
			Fields: []Field{
				{Name: "id", Type: "identifier"},
			},
		},
	}

	err := ps.mergeTypes(types1, "file1.constraint")
	if err != nil {
		t.Errorf("Unexpected error merging types: %v", err)
	}

	types2 := []TypeDefinition{
		{
			Name: "Person",
			Fields: []Field{
				{Name: "id", Type: "string"},
			},
		},
	}

	err = ps.mergeTypes(types2, "file2.constraint")
	if err == nil {
		t.Error("Expected error for incompatible types, got nil")
	}
}

// TestMergeRulesNilProgramState tests mergeRules with nil ProgramState
func TestMergeRulesNilProgramState(t *testing.T) {
	var ps *ProgramState = nil

	rules := []Expression{
		{
			Set: Set{
				Variables: []TypedVariable{
					{Name: "p", DataType: "Person"},
				},
			},
		},
	}

	err := ps.mergeRules(rules, "test.constraint")
	if err == nil {
		t.Error("Expected error for nil ProgramState, got nil")
	}
}

// TestMergeFactsNilProgramState tests mergeFacts with nil ProgramState
func TestMergeFactsNilProgramState(t *testing.T) {
	var ps *ProgramState = nil

	facts := []Fact{
		{
			TypeName: "Person",
			Fields: []FactField{
				{Name: "id", Value: FactValue{Type: "identifier", Value: "P001"}},
			},
		},
	}

	err := ps.mergeFacts(facts, "test.constraint")
	if err == nil {
		t.Error("Expected error for nil ProgramState, got nil")
	}
}

// TestToProgram tests ToProgram conversion
func TestToProgram(t *testing.T) {
	ps := NewProgramState()

	ps.Types["Person"] = &TypeDefinition{
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "identifier"},
		},
	}

	ps.Rules = append(ps.Rules, &Expression{
		Set: Set{
			Variables: []TypedVariable{
				{Name: "p", DataType: "Person"},
			},
		},
	})

	ps.Facts = append(ps.Facts, &Fact{
		TypeName: "Person",
		Fields: []FactField{
			{Name: "id", Value: FactValue{Type: "identifier", Value: "P001"}},
		},
	})

	program := ps.ToProgram()

	if program == nil {
		t.Fatal("Expected program, got nil")
	}

	if len(program.Types) != 1 {
		t.Errorf("Expected 1 type, got %d", len(program.Types))
	}

	if len(program.Expressions) != 1 {
		t.Errorf("Expected 1 expression, got %d", len(program.Expressions))
	}

	if len(program.Facts) != 1 {
		t.Errorf("Expected 1 fact, got %d", len(program.Facts))
	}
}

// TestValidateFactValueTypeErrors tests validateFactValue with type mismatches
func TestValidateFactValueTypeErrors(t *testing.T) {
	ps := NewProgramState()

	// Test string type mismatch
	err := ps.validateFactValue(FactValue{Type: "number", Value: 42}, "string")
	if err == nil {
		t.Error("Expected error for number value with string type, got nil")
	}

	// Test number type mismatch
	err = ps.validateFactValue(FactValue{Type: "string", Value: "test"}, "number")
	if err == nil {
		t.Error("Expected error for string value with number type, got nil")
	}

	// Test bool type mismatch
	err = ps.validateFactValue(FactValue{Type: "string", Value: "true"}, "bool")
	if err == nil {
		t.Error("Expected error for string value with bool type, got nil")
	}
}

// TestValidateFactValueIdentifierAsString tests identifier accepted as string
func TestValidateFactValueIdentifierAsString(t *testing.T) {
	ps := NewProgramState()

	// Identifier should be accepted as string
	err := ps.validateFactValue(FactValue{Type: "identifier", Value: "ID001"}, "string")
	if err != nil {
		t.Errorf("Unexpected error for identifier as string: %v", err)
	}
}

// TestGetStringValueValidCases tests getStringValue with various inputs
func TestGetStringValueValidCases(t *testing.T) {
	// Test with valid string
	m := map[string]interface{}{
		"name": "Alice",
		"age":  25,
	}

	result := getStringValue(m, "name")
	if result != "Alice" {
		t.Errorf("Expected 'Alice', got '%s'", result)
	}

	// Test with non-string value
	result = getStringValue(m, "age")
	if result != "" {
		t.Errorf("Expected empty string for non-string value, got '%s'", result)
	}

	// Test with non-existent key
	result = getStringValue(m, "nonexistent")
	if result != "" {
		t.Errorf("Expected empty string for non-existent key, got '%s'", result)
	}
}

// TestScanForFieldAccessArrayNested tests nested array scanning
func TestScanForFieldAccessArrayNested(t *testing.T) {
	ps := NewProgramState()

	ps.Types["Person"] = &TypeDefinition{
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "identifier"},
			{Name: "name", Type: "string"},
		},
	}

	variables := map[string]string{
		"p": "Person",
	}

	// Nested array with maps
	data := []interface{}{
		[]interface{}{
			map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "name",
			},
		},
	}

	err := ps.scanForFieldAccess(data, variables)
	if err != nil {
		t.Errorf("Unexpected error for nested array: %v", err)
	}
}

// TestScanForFieldAccessUnknownVariable tests field access with unknown variable
func TestScanForFieldAccessUnknownVariable(t *testing.T) {
	ps := NewProgramState()

	ps.Types["Person"] = &TypeDefinition{
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "identifier"},
		},
	}

	variables := map[string]string{
		"p": "Person",
	}

	// Field access with unknown variable
	data := map[string]interface{}{
		"type":   "fieldAccess",
		"object": "unknown",
		"field":  "id",
	}

	// Should not error for unknown variable (just skip validation)
	err := ps.scanForFieldAccess(data, variables)
	if err != nil {
		t.Errorf("Unexpected error for unknown variable: %v", err)
	}
}

// TestMergeTypesEmptySlice tests mergeTypes with empty slice
func TestMergeTypesEmptySlice(t *testing.T) {
	ps := NewProgramState()

	err := ps.mergeTypes([]TypeDefinition{}, "test.constraint")
	if err != nil {
		t.Errorf("Unexpected error for empty types: %v", err)
	}

	if len(ps.Types) != 0 {
		t.Errorf("Expected 0 types, got %d", len(ps.Types))
	}
}

// TestMergeRulesWithInvalidRule tests mergeRules with rule referencing undefined type
func TestMergeRulesWithInvalidRule(t *testing.T) {
	ps := NewProgramState()

	rules := []Expression{
		{
			Set: Set{
				Variables: []TypedVariable{
					{Name: "p", DataType: "UndefinedType"},
				},
			},
			Constraints: map[string]interface{}{},
		},
	}

	err := ps.mergeRules(rules, "test.constraint")
	// Should not return error (non-blocking), but should record it
	if err != nil {
		t.Errorf("Unexpected error (should be non-blocking): %v", err)
	}

	if !ps.HasErrors() {
		t.Error("Expected validation error to be recorded")
	}

	if len(ps.Rules) != 0 {
		t.Errorf("Expected 0 valid rules, got %d", len(ps.Rules))
	}
}

// TestValidateRuleWithFieldAccessError tests validateRule with invalid field access
func TestValidateRuleWithFieldAccessError(t *testing.T) {
	ps := NewProgramState()

	ps.Types["Person"] = &TypeDefinition{
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "identifier"},
			{Name: "name", Type: "string"},
		},
	}

	rule := &Expression{
		Set: Set{
			Variables: []TypedVariable{
				{Name: "p", DataType: "Person"},
			},
		},
		Constraints: map[string]interface{}{
			"type":   "fieldAccess",
			"object": "p",
			"field":  "invalidField",
		},
	}

	err := ps.validateRule(rule, "test.constraint")
	if err == nil {
		t.Error("Expected error for invalid field access in constraints, got nil")
	}
}

// TestValidateFactWithWrongFieldType tests validateFact with wrong field type
func TestValidateFactWithWrongFieldType(t *testing.T) {
	ps := NewProgramState()

	ps.Types["Person"] = &TypeDefinition{
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "identifier"},
			{Name: "age", Type: "number"},
		},
	}

	// Fact with string value for number field
	fact := &Fact{
		TypeName: "Person",
		Fields: []FactField{
			{Name: "id", Value: FactValue{Type: "identifier", Value: "P001"}},
			{Name: "age", Value: FactValue{Type: "string", Value: "25"}},
		},
	}

	err := ps.validateFact(fact, "test.constraint")
	if err == nil {
		t.Error("Expected error for wrong field type, got nil")
	}
}

// TestValidateTypes tests ValidateTypes function
func TestValidateTypes(t *testing.T) {
	t.Run("all types defined", func(t *testing.T) {
		program := Program{
			Types: []TypeDefinition{
				{Name: "Person", Fields: []Field{{Name: "name", Type: "string"}}},
			},
			Expressions: []Expression{
				{
					Set: Set{
						Variables: []TypedVariable{{Name: "p", DataType: "Person"}},
					},
				},
			},
		}
		err := ValidateTypes(program)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("undefined type", func(t *testing.T) {
		program := Program{
			Types: []TypeDefinition{},
			Expressions: []Expression{
				{
					Set: Set{
						Variables: []TypedVariable{{Name: "p", DataType: "UnknownType"}},
					},
				},
			},
		}
		err := ValidateTypes(program)
		if err == nil {
			t.Error("Expected error for undefined type, got nil")
		}
	})
}

// TestGetTypeFields tests GetTypeFields function
func TestGetTypeFields(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "id", Type: "identifier"},
					{Name: "name", Type: "string"},
					{Name: "age", Type: "number"},
				},
			},
		},
	}

	t.Run("existing type", func(t *testing.T) {
		fields, err := GetTypeFields(program, "Person")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(fields) != 3 {
			t.Errorf("Expected 3 fields, got %d", len(fields))
		}
	})

	t.Run("non-existing type", func(t *testing.T) {
		_, err := GetTypeFields(program, "Unknown")
		if err == nil {
			t.Error("Expected error for non-existing type, got nil")
		}
	})
}

// TestValidateFieldAccess tests ValidateFieldAccess function
func TestValidateFieldAccess(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "name", Type: "string"},
					{Name: "age", Type: "number"},
				},
			},
		},
		Expressions: []Expression{
			{
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
					},
				},
			},
		},
	}

	t.Run("valid field access", func(t *testing.T) {
		fa := FieldAccess{Object: "p", Field: "name"}
		err := ValidateFieldAccess(program, fa, 0)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("invalid expression index", func(t *testing.T) {
		fa := FieldAccess{Object: "p", Field: "name"}
		err := ValidateFieldAccess(program, fa, 99)
		if err == nil {
			t.Error("Expected error for invalid expression index, got nil")
		}
	})

	t.Run("unknown variable", func(t *testing.T) {
		fa := FieldAccess{Object: "unknown", Field: "name"}
		err := ValidateFieldAccess(program, fa, 0)
		if err == nil {
			t.Error("Expected error for unknown variable, got nil")
		}
	})

	t.Run("unknown field", func(t *testing.T) {
		fa := FieldAccess{Object: "p", Field: "unknown"}
		err := ValidateFieldAccess(program, fa, 0)
		if err == nil {
			t.Error("Expected error for unknown field, got nil")
		}
	})
}

// TestGetFieldType tests GetFieldType function
func TestGetFieldType(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "name", Type: "string"},
					{Name: "age", Type: "number"},
				},
			},
		},
		Expressions: []Expression{
			{
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
					},
				},
			},
		},
	}

	t.Run("valid field type", func(t *testing.T) {
		fieldType, err := GetFieldType(program, "p", "name", 0)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if fieldType != "string" {
			t.Errorf("Expected 'string', got '%s'", fieldType)
		}
	})

	t.Run("invalid expression index", func(t *testing.T) {
		_, err := GetFieldType(program, "p", "name", 99)
		if err == nil {
			t.Error("Expected error for invalid expression index, got nil")
		}
	})

	t.Run("unknown variable", func(t *testing.T) {
		_, err := GetFieldType(program, "unknown", "name", 0)
		if err == nil {
			t.Error("Expected error for unknown variable, got nil")
		}
	})

	t.Run("unknown field", func(t *testing.T) {
		_, err := GetFieldType(program, "p", "unknown", 0)
		if err == nil {
			t.Error("Expected error for unknown field, got nil")
		}
	})
}

// TestGetValueType tests GetValueType function
func TestGetValueType(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "number type",
			value:    map[string]interface{}{"type": "number", "value": 42},
			expected: "number",
		},
		{
			name:     "string type",
			value:    map[string]interface{}{"type": "string", "value": "hello"},
			expected: "string",
		},
		{
			name:     "boolean type",
			value:    map[string]interface{}{"type": "boolean", "value": true},
			expected: "bool",
		},
		{
			name:     "true variable",
			value:    map[string]interface{}{"type": "variable", "name": "true"},
			expected: "bool",
		},
		{
			name:     "false variable",
			value:    map[string]interface{}{"type": "variable", "name": "false"},
			expected: "bool",
		},
		{
			name:     "other variable",
			value:    map[string]interface{}{"type": "variable", "name": "x"},
			expected: "variable",
		},
		{
			name:     "unknown type",
			value:    map[string]interface{}{"type": "unknown"},
			expected: "unknown",
		},
		{
			name:     "non-map value",
			value:    "plain string",
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetValueType(tt.value)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestExtractVariablesFromArg tests extractVariablesFromArg function
func TestExtractVariablesFromArg(t *testing.T) {
	t.Run("simple string", func(t *testing.T) {
		vars := extractVariablesFromArg("varName")
		if len(vars) != 1 || vars[0] != "varName" {
			t.Errorf("Expected [varName], got %v", vars)
		}
	})

	t.Run("fieldAccess map", func(t *testing.T) {
		arg := map[string]interface{}{
			"type":   "fieldAccess",
			"object": "person",
			"field":  "name",
		}
		vars := extractVariablesFromArg(arg)
		if len(vars) != 1 || vars[0] != "person" {
			t.Errorf("Expected [person], got %v", vars)
		}
	})

	t.Run("string literal", func(t *testing.T) {
		arg := map[string]interface{}{
			"type":  "string",
			"value": "hello",
		}
		vars := extractVariablesFromArg(arg)
		if len(vars) != 0 {
			t.Errorf("Expected no variables for string literal, got %v", vars)
		}
	})

	t.Run("number literal", func(t *testing.T) {
		arg := map[string]interface{}{
			"type":  "number",
			"value": 42,
		}
		vars := extractVariablesFromArg(arg)
		if len(vars) != 0 {
			t.Errorf("Expected no variables for number literal, got %v", vars)
		}
	})
}

// TestParseAndMergeContentValidInput tests ParseAndMergeContent with valid inputs
func TestParseAndMergeContentValidInput(t *testing.T) {
	ps := NewProgramState()

	content := "type Person : <id: string, name: string>"
	err := ps.ParseAndMergeContent(content, "test.constraint")
	if err != nil {
		t.Errorf("Expected no error for valid content, got %v", err)
	}

	if len(ps.FilesParsed) != 1 {
		t.Errorf("Expected 1 file parsed, got %d", len(ps.FilesParsed))
	}
}

// TestValidateFacts tests ValidateFacts function
func TestValidateFacts(t *testing.T) {
	t.Run("valid facts", func(t *testing.T) {
		program := Program{
			Types: []TypeDefinition{
				{
					Name: "Person",
					Fields: []Field{
						{Name: "id", Type: "string"},
						{Name: "age", Type: "number"},
					},
				},
			},
			Facts: []Fact{
				{
					TypeName: "Person",
					Fields: []FactField{
						{Name: "id", Value: FactValue{Type: "string", Value: "p1"}},
						{Name: "age", Value: FactValue{Type: "number", Value: 25}},
					},
				},
			},
		}
		err := ValidateFacts(program)
		if err != nil {
			t.Errorf("Expected no error for valid facts, got %v", err)
		}
	})

	t.Run("undefined type", func(t *testing.T) {
		program := Program{
			Types: []TypeDefinition{},
			Facts: []Fact{
				{TypeName: "Unknown", Fields: []FactField{}},
			},
		}
		err := ValidateFacts(program)
		if err == nil {
			t.Error("Expected error for undefined type, got nil")
		}
	})

	t.Run("undefined field", func(t *testing.T) {
		program := Program{
			Types: []TypeDefinition{
				{Name: "Person", Fields: []Field{{Name: "id", Type: "string"}}},
			},
			Facts: []Fact{
				{
					TypeName: "Person",
					Fields: []FactField{
						{Name: "unknown", Value: FactValue{Type: "string", Value: "x"}},
					},
				},
			},
		}
		err := ValidateFacts(program)
		if err == nil {
			t.Error("Expected error for undefined field, got nil")
		}
	})

	t.Run("wrong field type", func(t *testing.T) {
		program := Program{
			Types: []TypeDefinition{
				{Name: "Person", Fields: []Field{{Name: "age", Type: "number"}}},
			},
			Facts: []Fact{
				{
					TypeName: "Person",
					Fields: []FactField{
						{Name: "age", Value: FactValue{Type: "string", Value: "not a number"}},
					},
				},
			},
		}
		err := ValidateFacts(program)
		if err == nil {
			t.Error("Expected error for wrong field type, got nil")
		}
	})
}

// TestValidateFactFieldType tests ValidateFactFieldType function
func TestValidateFactFieldType(t *testing.T) {
	tests := []struct {
		name         string
		value        FactValue
		expectedType string
		shouldError  bool
	}{
		{"string type valid", FactValue{Type: "string", Value: "hello"}, "string", false},
		{"identifier as string", FactValue{Type: "identifier", Value: "id1"}, "string", false},
		{"number type valid", FactValue{Type: "number", Value: 42}, "number", false},
		{"boolean type valid", FactValue{Type: "boolean", Value: true}, "bool", false},
		{"boolean type valid (boolean)", FactValue{Type: "boolean", Value: false}, "boolean", false},
		{"string type mismatch", FactValue{Type: "number", Value: 42}, "string", true},
		{"number type mismatch", FactValue{Type: "string", Value: "x"}, "number", true},
		{"boolean type mismatch", FactValue{Type: "string", Value: "x"}, "bool", true},
		{"unknown type", FactValue{Type: "custom", Value: "x"}, "customType", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFactFieldType(tt.value, tt.expectedType, "TestType", "testField")
			if tt.shouldError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got %v", err)
			}
		})
	}
}

// TestValidateProgram tests ValidateProgram function
func TestValidateProgram(t *testing.T) {
	t.Run("valid program with all components", func(t *testing.T) {
		// Create a simple valid program structure
		programData := map[string]interface{}{
			"types": []map[string]interface{}{
				{
					"name": "Person",
					"fields": []map[string]interface{}{
						{"name": "id", "type": "string"},
						{"name": "name", "type": "string"},
					},
				},
			},
			"expressions": []map[string]interface{}{
				{
					"set": map[string]interface{}{
						"variables": []map[string]interface{}{
							{"name": "p", "dataType": "Person"},
						},
					},
					"constraints": map[string]interface{}{
						"type": "comparison",
						"left": map[string]interface{}{
							"type":   "fieldAccess",
							"object": "p",
							"field":  "id",
						},
						"operator": "==",
						"right": map[string]interface{}{
							"type":  "string",
							"value": "p1",
						},
					},
					"action": map[string]interface{}{
						"jobCall": map[string]interface{}{
							"name": "match",
							"args": []interface{}{"p"},
						},
					},
				},
			},
			"facts": []map[string]interface{}{},
		}

		err := ValidateProgram(programData)
		if err != nil {
			t.Errorf("Expected no error for valid program, got %v", err)
		}
	})

	t.Run("program with undefined type in expression", func(t *testing.T) {
		programData := map[string]interface{}{
			"types": []map[string]interface{}{},
			"expressions": []map[string]interface{}{
				{
					"set": map[string]interface{}{
						"variables": []map[string]interface{}{
							{"name": "p", "dataType": "UndefinedType"},
						},
					},
					"constraints": nil,
					"action": map[string]interface{}{
						"jobCall": map[string]interface{}{
							"name": "test",
							"args": []interface{}{},
						},
					},
				},
			},
			"facts": []map[string]interface{}{},
		}

		err := ValidateProgram(programData)
		if err == nil {
			t.Error("Expected error for undefined type, got nil")
		}
	})

	t.Run("program with missing action", func(t *testing.T) {
		programData := map[string]interface{}{
			"types": []map[string]interface{}{
				{
					"name":   "Person",
					"fields": []map[string]interface{}{},
				},
			},
			"expressions": []map[string]interface{}{
				{
					"set": map[string]interface{}{
						"variables": []map[string]interface{}{
							{"name": "p", "dataType": "Person"},
						},
					},
					"constraints": nil,
					"action":      nil,
				},
			},
			"facts": []map[string]interface{}{},
		}

		err := ValidateProgram(programData)
		if err == nil {
			t.Error("Expected error for missing action, got nil")
		}
	})
}

// TestConvertFactsToReteFormat tests ConvertFactsToReteFormat function
func TestConvertFactsToReteFormat(t *testing.T) {
	program := Program{
		Facts: []Fact{
			{
				TypeName: "Person",
				Fields: []FactField{
					{Name: "id", Value: FactValue{Type: "identifier", Value: "p1"}},
					{Name: "name", Value: FactValue{Type: "string", Value: map[string]interface{}{"value": "John"}}},
					{Name: "age", Value: FactValue{Type: "number", Value: map[string]interface{}{"value": 30}}},
					{Name: "active", Value: FactValue{Type: "boolean", Value: map[string]interface{}{"value": true}}},
				},
			},
		},
	}

	reteFacts := ConvertFactsToReteFormat(program)

	if len(reteFacts) != 1 {
		t.Fatalf("Expected 1 rete fact, got %d", len(reteFacts))
	}

	reteFact := reteFacts[0]

	if reteFact["reteType"] != "Person" {
		t.Errorf("Expected reteType 'Person', got %v", reteFact["reteType"])
	}

	if reteFact["id"] != "p1" {
		t.Errorf("Expected id 'p1', got %v", reteFact["id"])
	}

	if reteFact["name"] != "John" {
		t.Errorf("Expected name 'John', got %v", reteFact["name"])
	}

	if reteFact["age"] != 30 {
		t.Errorf("Expected age 30, got %v", reteFact["age"])
	}

	if reteFact["active"] != true {
		t.Errorf("Expected active true, got %v", reteFact["active"])
	}
}

// TestValidateAction tests ValidateAction function
func TestValidateAction(t *testing.T) {
	program := Program{
		Expressions: []Expression{
			{
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
						{Name: "o", DataType: "Order"},
					},
				},
			},
		},
	}

	t.Run("valid action with string args", func(t *testing.T) {
		action := Action{
			Job: JobCall{
				Name: "process",
				Args: []interface{}{"p", "o"},
			},
		}
		err := ValidateAction(program, action, 0)
		if err != nil {
			t.Errorf("Expected no error for valid action, got %v", err)
		}
	})

	t.Run("invalid expression index", func(t *testing.T) {
		action := Action{
			Job: JobCall{
				Name: "process",
				Args: []interface{}{"p"},
			},
		}
		err := ValidateAction(program, action, 99)
		if err == nil {
			t.Error("Expected error for invalid expression index, got nil")
		}
	})

	t.Run("action with undefined variable", func(t *testing.T) {
		action := Action{
			Job: JobCall{
				Name: "process",
				Args: []interface{}{"unknown"},
			},
		}
		err := ValidateAction(program, action, 0)
		if err == nil {
			t.Error("Expected error for undefined variable in action, got nil")
		}
	})

	t.Run("action with field access", func(t *testing.T) {
		action := Action{
			Job: JobCall{
				Name: "process",
				Args: []interface{}{
					map[string]interface{}{
						"type":   "fieldAccess",
						"object": "p",
						"field":  "id",
					},
				},
			},
		}
		err := ValidateAction(program, action, 0)
		if err != nil {
			t.Errorf("Expected no error for action with field access, got %v", err)
		}
	})
}

// TestValidateLogicalExpressionConstraint tests validateLogicalExpressionConstraint function
func TestValidateLogicalExpressionConstraint(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "age", Type: "number"},
					{Name: "name", Type: "string"},
				},
			},
		},
		Expressions: []Expression{
			{
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
					},
				},
			},
		},
	}

	t.Run("valid logical expression", func(t *testing.T) {
		constraint := map[string]interface{}{
			"type": "logicalExpr",
			"left": map[string]interface{}{
				"type": "comparison",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "age",
				},
				"operator": ">",
				"right": map[string]interface{}{
					"type":  "number",
					"value": 18,
				},
			},
			"operations": []interface{}{
				map[string]interface{}{
					"operator": "&&",
					"right": map[string]interface{}{
						"type": "comparison",
						"left": map[string]interface{}{
							"type":   "fieldAccess",
							"object": "p",
							"field":  "name",
						},
						"operator": "!=",
						"right": map[string]interface{}{
							"type":  "string",
							"value": "",
						},
					},
				},
			},
		}
		err := validateLogicalExpressionConstraint(program, constraint, 0)
		if err != nil {
			t.Errorf("Expected no error for valid logical expression, got %v", err)
		}
	})

	t.Run("logical expression without operations", func(t *testing.T) {
		constraint := map[string]interface{}{
			"type": "logicalExpr",
			"left": map[string]interface{}{
				"type": "comparison",
			},
		}
		err := validateLogicalExpressionConstraint(program, constraint, 0)
		if err != nil {
			t.Errorf("Expected no error when operations is missing, got %v", err)
		}
	})
}

// TestValidateBinaryOpConstraint tests validateBinaryOpConstraint function
func TestValidateBinaryOpConstraint(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "age", Type: "number"},
				},
			},
		},
		Expressions: []Expression{
			{
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
					},
				},
			},
		},
	}

	t.Run("valid binary operation", func(t *testing.T) {
		constraint := map[string]interface{}{
			"type": "binaryOp",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "age",
			},
			"operator": "+",
			"right": map[string]interface{}{
				"type":  "number",
				"value": 5,
			},
		}
		err := validateBinaryOpConstraint(program, constraint, 0)
		if err != nil {
			t.Errorf("Expected no error for valid binary operation, got %v", err)
		}
	})

	t.Run("binary operation without right operand", func(t *testing.T) {
		constraint := map[string]interface{}{
			"type": "binaryOp",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "age",
			},
			"operator": "+",
		}
		err := validateBinaryOpConstraint(program, constraint, 0)
		if err != nil {
			t.Errorf("Expected no error when right is missing, got %v", err)
		}
	})
}

// TestValidateConstraintFieldAccessRecursive tests ValidateConstraintFieldAccess with nested constraints
func TestValidateConstraintFieldAccessRecursive(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "age", Type: "number"},
					{Name: "name", Type: "string"},
				},
			},
		},
		Expressions: []Expression{
			{
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
					},
				},
			},
		},
	}

	t.Run("nested comparison constraint", func(t *testing.T) {
		constraint := map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "age",
			},
			"operator": ">=",
			"right": map[string]interface{}{
				"type":  "number",
				"value": 18,
			},
		}
		err := ValidateConstraintFieldAccess(program, constraint, 0)
		if err != nil {
			t.Errorf("Expected no error for valid nested constraint, got %v", err)
		}
	})

	t.Run("nested logical expression", func(t *testing.T) {
		constraint := map[string]interface{}{
			"type": "logicalExpr",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "age",
			},
			"operations": []interface{}{
				map[string]interface{}{
					"operator": "&&",
					"right": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "p",
						"field":  "name",
					},
				},
			},
		}
		err := ValidateConstraintFieldAccess(program, constraint, 0)
		if err != nil {
			t.Errorf("Expected no error for logical expression, got %v", err)
		}
	})

	t.Run("invalid field in nested constraint", func(t *testing.T) {
		constraint := map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "invalidField",
			},
			"operator": "==",
			"right": map[string]interface{}{
				"type":  "number",
				"value": 10,
			},
		}
		err := ValidateConstraintFieldAccess(program, constraint, 0)
		if err == nil {
			t.Error("Expected error for invalid field, got nil")
		}
	})

	t.Run("constraint without type", func(t *testing.T) {
		constraint := map[string]interface{}{
			"operator": "==",
		}
		err := ValidateConstraintFieldAccess(program, constraint, 0)
		if err != nil {
			t.Errorf("Expected no error for constraint without type, got %v", err)
		}
	})

	t.Run("non-map constraint", func(t *testing.T) {
		constraint := "not a map"
		err := ValidateConstraintFieldAccess(program, constraint, 0)
		if err != nil {
			t.Errorf("Expected no error for non-map constraint, got %v", err)
		}
	})
}
