// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestFunctionRegistry(t *testing.T) {
	t.Run("NewFunctionRegistry creates registry with defaults", func(t *testing.T) {
		fr := NewFunctionRegistry()
		if fr == nil {
			t.Fatal("NewFunctionRegistry() returned nil")
		}

		// Test default functions
		if returnType := fr.GetReturnType("LENGTH", "default"); returnType != ValueTypeNumber {
			t.Errorf("Expected LENGTH to return %s, got %s", ValueTypeNumber, returnType)
		}
		if returnType := fr.GetReturnType("UPPER", "default"); returnType != ValueTypeString {
			t.Errorf("Expected UPPER to return %s, got %s", ValueTypeString, returnType)
		}
	})

	t.Run("RegisterFunction adds new function", func(t *testing.T) {
		fr := NewFunctionRegistry()
		fr.RegisterFunction("CUSTOM_FUNC", ValueTypeBool, []string{ValueTypeString})

		returnType := fr.GetReturnType("CUSTOM_FUNC", "default")
		if returnType != ValueTypeBool {
			t.Errorf("Expected CUSTOM_FUNC to return %s, got %s", ValueTypeBool, returnType)
		}
	})

	t.Run("HasFunction checks function existence", func(t *testing.T) {
		fr := NewFunctionRegistry()

		if !fr.HasFunction("LENGTH") {
			t.Error("Expected LENGTH to exist in registry")
		}
		if fr.HasFunction("NONEXISTENT") {
			t.Error("Expected NONEXISTENT to not exist in registry")
		}
	})

	t.Run("GetReturnType is case insensitive", func(t *testing.T) {
		fr := NewFunctionRegistry()

		tests := []string{"LENGTH", "length", "LeNgTh"}
		for _, funcName := range tests {
			if returnType := fr.GetReturnType(funcName, "default"); returnType != ValueTypeNumber {
				t.Errorf("Expected %s to return %s, got %s", funcName, ValueTypeNumber, returnType)
			}
		}
	})

	t.Run("GetReturnType returns default for unknown function", func(t *testing.T) {
		fr := NewFunctionRegistry()

		returnType := fr.GetReturnType("UNKNOWN_FUNCTION", "mydefault")
		if returnType != "mydefault" {
			t.Errorf("Expected default value 'mydefault', got %s", returnType)
		}
	})

	t.Run("GetSignature returns full signature", func(t *testing.T) {
		fr := NewFunctionRegistry()

		sig, exists := fr.GetSignature("SUBSTRING")
		if !exists {
			t.Error("Expected SUBSTRING signature to exist")
		}
		if sig.ReturnType != ValueTypeString {
			t.Errorf("Expected SUBSTRING return type %s, got %s", ValueTypeString, sig.ReturnType)
		}
	})
}

func TestValidationHelpers(t *testing.T) {
	t.Run("findVariableType finds variable in single pattern", func(t *testing.T) {
		expr := Expression{
			Set: Set{
				Variables: []TypedVariable{
					{Name: "p", DataType: "Person"},
					{Name: "o", DataType: "Order"},
				},
			},
		}

		varType, err := findVariableType(expr, "p")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if varType != "Person" {
			t.Errorf("Expected 'Person', got %s", varType)
		}
	})

	t.Run("findVariableType finds variable in multi-pattern", func(t *testing.T) {
		expr := Expression{
			Patterns: []Set{
				{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
					},
				},
				{
					Variables: []TypedVariable{
						{Name: "o", DataType: "Order"},
					},
				},
			},
		}

		varType, err := findVariableType(expr, "o")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if varType != "Order" {
			t.Errorf("Expected 'Order', got %s", varType)
		}
	})

	t.Run("findVariableType returns error for empty name", func(t *testing.T) {
		expr := Expression{}

		_, err := findVariableType(expr, "")
		if err == nil {
			t.Error("Expected error for empty variable name")
		}
	})

	t.Run("findVariableType returns error for not found", func(t *testing.T) {
		expr := Expression{
			Set: Set{
				Variables: []TypedVariable{
					{Name: "p", DataType: "Person"},
				},
			},
		}

		_, err := findVariableType(expr, "notfound")
		if err == nil {
			t.Error("Expected error for variable not found")
		}
	})

	t.Run("sanitizeForLog removes control characters", func(t *testing.T) {
		input := "Hello\x00World\x01Test"
		output := sanitizeForLog(input, 100)

		// Should not contain null bytes
		if strings.Contains(output, "\x00") {
			t.Error("Output should not contain null bytes")
		}
	})

	t.Run("sanitizeForLog truncates long strings", func(t *testing.T) {
		input := strings.Repeat("a", 1000)
		output := sanitizeForLog(input, 100)

		if len(output) > 110 { // 100 + "..."
			t.Errorf("Output should be truncated to ~100 chars, got %d", len(output))
		}
	})

	t.Run("validateInputNotNil checks nil values", func(t *testing.T) {
		err := validateInputNotNil(map[string]interface{}{
			"field1": "value",
			"field2": nil,
		})

		if err == nil {
			t.Error("Expected error for nil value")
		}
		if !strings.Contains(err.Error(), "field2") {
			t.Errorf("Error should mention field2, got: %v", err)
		}
	})

	t.Run("validateInputNotNil accepts non-nil values", func(t *testing.T) {
		err := validateInputNotNil(map[string]interface{}{
			"field1": "value",
			"field2": 123,
		})

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("isArithmeticOperator recognizes arithmetic ops", func(t *testing.T) {
		tests := []struct {
			op       string
			expected bool
		}{
			{"+", true},
			{"-", true},
			{"*", true},
			{"/", true},
			{"%", true},
			{"==", false},
			{"!=", false},
			{"<", false},
		}

		for _, tt := range tests {
			result := isArithmeticOperator(tt.op)
			if result != tt.expected {
				t.Errorf("isArithmeticOperator(%q) = %v, want %v", tt.op, result, tt.expected)
			}
		}
	})

	t.Run("isComparisonOperator recognizes comparison ops", func(t *testing.T) {
		tests := []struct {
			op       string
			expected bool
		}{
			{"==", true},
			{"!=", true},
			{"<", true},
			{">", true},
			{"<=", true},
			{">=", true},
			{"+", false},
			{"-", false},
		}

		for _, tt := range tests {
			result := isComparisonOperator(tt.op)
			if result != tt.expected {
				t.Errorf("isComparisonOperator(%q) = %v, want %v", tt.op, result, tt.expected)
			}
		}
	})
}

func TestSafeBase64Decode(t *testing.T) {
	t.Run("decodes valid base64", func(t *testing.T) {
		encoded := "SGVsbG8gV29ybGQ=" // "Hello World"
		decoded, err := safeBase64Decode(encoded)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if decoded != "Hello World" {
			t.Errorf("Expected 'Hello World', got %s", decoded)
		}
	})

	t.Run("rejects invalid base64", func(t *testing.T) {
		encoded := "not valid base64!!!"
		_, err := safeBase64Decode(encoded)

		if err == nil {
			t.Error("Expected error for invalid base64")
		}
	})

	t.Run("rejects non-UTF8 output", func(t *testing.T) {
		// Binary data that's not valid UTF-8
		binary := []byte{0xFF, 0xFE, 0xFD}
		encoded := base64.StdEncoding.EncodeToString(binary)

		_, err := safeBase64Decode(encoded)
		if err == nil {
			t.Error("Expected error for non-UTF8 data")
		}
		if !strings.Contains(err.Error(), "UTF-8") {
			t.Errorf("Error should mention UTF-8, got: %v", err)
		}
	})

	t.Run("rejects oversized payload", func(t *testing.T) {
		// Create a payload larger than MaxBase64DecodeSize
		largeData := make([]byte, MaxBase64DecodeSize+1000)
		encoded := base64.StdEncoding.EncodeToString(largeData)

		_, err := safeBase64Decode(encoded)
		if err == nil {
			t.Error("Expected error for oversized payload")
		}
	})
}

func TestDepthLimits(t *testing.T) {
	t.Run("ValidateConstraintFieldAccess respects depth limit", func(t *testing.T) {
		program := Program{
			Types: []TypeDefinition{
				{
					Name: "Person",
					Fields: []Field{
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

		// Create deeply nested constraint (depth > MaxValidationDepth)
		constraint := map[string]interface{}{"type": "comparison"}
		current := constraint

		// Create nesting deeper than MaxValidationDepth
		for i := 0; i < MaxValidationDepth+10; i++ {
			nested := map[string]interface{}{"type": "comparison"}
			current["left"] = nested
			current = nested
		}

		err := ValidateConstraintFieldAccess(program, constraint, 0)
		if err == nil {
			t.Error("Expected error for excessive depth")
		}
		if !strings.Contains(err.Error(), "depth exceeded") {
			t.Errorf("Error should mention depth, got: %v", err)
		}
	})

	t.Run("ValidateTypeCompatibility respects depth limit", func(t *testing.T) {
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

		// Create deeply nested constraint with valid structure
		var createNested func(depth int) map[string]interface{}
		createNested = func(depth int) map[string]interface{} {
			if depth <= 0 {
				return map[string]interface{}{
					"type": "comparison",
					"left": map[string]interface{}{
						"type": "number",
					},
					"right": map[string]interface{}{
						"type": "number",
					},
				}
			}
			return map[string]interface{}{
				"type": "comparison",
				"left": createNested(depth - 1),
				"right": map[string]interface{}{
					"type": "number",
				},
			}
		}

		constraint := createNested(MaxValidationDepth + 10)

		err := ValidateTypeCompatibility(program, constraint, 0)
		if err == nil {
			t.Error("Expected error for excessive depth")
		}
		if !strings.Contains(err.Error(), "depth exceeded") {
			t.Errorf("Error should mention depth, got: %v", err)
		}
	})

	t.Run("inferArgumentType respects depth limit", func(t *testing.T) {
		av := NewActionValidator(nil, nil)

		// This test would need a very deeply nested arg structure
		// For now, just verify the depth parameter exists
		arg := map[string]interface{}{
			"type": "string",
		}

		_, err := av.inferArgumentType(arg, nil, MaxValidationDepth+1)
		if err == nil {
			t.Error("Expected error for excessive depth")
		}
		if !strings.Contains(err.Error(), "depth exceeded") {
			t.Errorf("Error should mention depth, got: %v", err)
		}
	})
}
