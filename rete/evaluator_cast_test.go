// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
// TestCastToNumber teste la conversion vers number
func TestCastToNumber(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    float64
		expectError bool
	}{
		// Number -> Number (déjà un nombre)
		{name: "float64 to number", input: 123.45, expected: 123.45, expectError: false},
		{name: "int to number", input: 42, expected: 42.0, expectError: false},
		{name: "int64 to number", input: int64(100), expected: 100.0, expectError: false},
		{name: "zero to number", input: 0.0, expected: 0.0, expectError: false},
		{name: "negative to number", input: -45.5, expected: -45.5, expectError: false},
		// String -> Number
		{name: "string integer to number", input: "123", expected: 123.0, expectError: false},
		{name: "string decimal to number", input: "12.5", expected: 12.5, expectError: false},
		{name: "string negative to number", input: "-45", expected: -45.0, expectError: false},
		{name: "string with spaces to number", input: "  123  ", expected: 123.0, expectError: false},
		{name: "string scientific notation", input: "1e3", expected: 1000.0, expectError: false},
		// String -> Number (erreurs)
		{name: "invalid string to number", input: "abc", expected: 0, expectError: true},
		{name: "mixed string to number", input: "12abc", expected: 0, expectError: true},
		{name: "empty string to number", input: "", expected: 0, expectError: true},
		// Bool -> Number
		{name: "true to number", input: true, expected: 1.0, expectError: false},
		{name: "false to number", input: false, expected: 0.0, expectError: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CastToNumber(tt.input)
			if tt.expectError {
				assert.Error(t, err, "Expected error for input: %v", tt.input)
			} else {
				require.NoError(t, err, "Unexpected error for input: %v", tt.input)
				assert.Equal(t, tt.expected, result, "Wrong conversion result")
			}
		})
	}
}
// TestCastToString teste la conversion vers string
func TestCastToString(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    string
		expectError bool
	}{
		// String -> String (déjà une chaîne)
		{name: "string to string", input: "hello", expected: "hello", expectError: false},
		{name: "empty string to string", input: "", expected: "", expectError: false},
		// Number -> String
		{name: "integer to string", input: 123.0, expected: "123", expectError: false},
		{name: "decimal to string", input: 12.5, expected: "12.5", expectError: false},
		{name: "negative to string", input: -45.0, expected: "-45", expectError: false},
		{name: "zero to string", input: 0.0, expected: "0", expectError: false},
		{name: "int type to string", input: 42, expected: "42", expectError: false},
		{name: "int64 type to string", input: int64(100), expected: "100", expectError: false},
		// Bool -> String
		{name: "true to string", input: true, expected: "true", expectError: false},
		{name: "false to string", input: false, expected: "false", expectError: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CastToString(tt.input)
			if tt.expectError {
				assert.Error(t, err, "Expected error for input: %v", tt.input)
			} else {
				require.NoError(t, err, "Unexpected error for input: %v", tt.input)
				assert.Equal(t, tt.expected, result, "Wrong conversion result")
			}
		})
	}
}
// TestCastToBool teste la conversion vers bool
func TestCastToBool(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    bool
		expectError bool
	}{
		// Bool -> Bool (déjà un booléen)
		{name: "true to bool", input: true, expected: true, expectError: false},
		{name: "false to bool", input: false, expected: false, expectError: false},
		// String -> Bool (valeurs vraies)
		{name: "string 'true' to bool", input: "true", expected: true, expectError: false},
		{name: "string 'TRUE' to bool", input: "TRUE", expected: true, expectError: false},
		{name: "string 'True' to bool", input: "True", expected: true, expectError: false},
		{name: "string '1' to bool", input: "1", expected: true, expectError: false},
		// String -> Bool (valeurs fausses)
		{name: "string 'false' to bool", input: "false", expected: false, expectError: false},
		{name: "string 'FALSE' to bool", input: "FALSE", expected: false, expectError: false},
		{name: "string 'False' to bool", input: "False", expected: false, expectError: false},
		{name: "string '0' to bool", input: "0", expected: false, expectError: false},
		{name: "empty string to bool", input: "", expected: false, expectError: false},
		// String -> Bool (avec espaces)
		{name: "string with spaces to bool", input: "  true  ", expected: true, expectError: false},
		{name: "string with spaces false", input: "  false  ", expected: false, expectError: false},
		// String -> Bool (comportement permissif)
		{name: "invalid string to bool (permissive)", input: "maybe", expected: false, expectError: false},
		{name: "other string to bool", input: "yes", expected: false, expectError: false},
		// Number -> Bool
		{name: "zero to bool", input: 0.0, expected: false, expectError: false},
		{name: "non-zero positive to bool", input: 1.0, expected: true, expectError: false},
		{name: "non-zero negative to bool", input: -5.0, expected: true, expectError: false},
		{name: "large number to bool", input: 999.0, expected: true, expectError: false},
		{name: "int zero to bool", input: 0, expected: false, expectError: false},
		{name: "int non-zero to bool", input: 42, expected: true, expectError: false},
		{name: "int64 zero to bool", input: int64(0), expected: false, expectError: false},
		{name: "int64 non-zero to bool", input: int64(100), expected: true, expectError: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CastToBool(tt.input)
			if tt.expectError {
				assert.Error(t, err, "Expected error for input: %v", tt.input)
			} else {
				require.NoError(t, err, "Unexpected error for input: %v", tt.input)
				assert.Equal(t, tt.expected, result, "Wrong conversion result")
			}
		})
	}
}
// TestEvaluateCast teste la fonction générique de cast
func TestEvaluateCast(t *testing.T) {
	tests := []struct {
		name        string
		castType    string
		input       interface{}
		expected    interface{}
		expectError bool
	}{
		// Cast vers number
		{name: "cast to number from string", castType: "number", input: "123", expected: 123.0, expectError: false},
		{name: "cast to number from bool", castType: "number", input: true, expected: 1.0, expectError: false},
		// Cast vers string
		{name: "cast to string from number", castType: "string", input: 123.0, expected: "123", expectError: false},
		{name: "cast to string from bool", castType: "string", input: true, expected: "true", expectError: false},
		// Cast vers bool
		{name: "cast to bool from string", castType: "bool", input: "true", expected: true, expectError: false},
		{name: "cast to bool from number", castType: "bool", input: 0.0, expected: false, expectError: false},
		// Type de cast invalide
		{name: "invalid cast type", castType: "integer", input: "123", expected: nil, expectError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := EvaluateCast(tt.castType, tt.input)
			if tt.expectError {
				assert.Error(t, err, "Expected error for cast type: %s", tt.castType)
			} else {
				require.NoError(t, err, "Unexpected error for cast type: %s", tt.castType)
				assert.Equal(t, tt.expected, result, "Wrong cast result")
			}
		})
	}
}
// TestCastInExpressions teste le cast dans des expressions plus complexes
func TestCastInExpressions(t *testing.T) {
	// Créer un évaluateur avec des données de test
	eval := &AlphaConditionEvaluator{
		Bindings: NewBindingChain().Add("p", {
				ID:   "p1").Add("price", "99.99").Add("quantity", "5").Add("active", "true").Add("count", 10.0).Add("flag", true),
			},
		},
	}
	tests := []struct {
		name        string
		castExpr    map[string]interface{}
		expected    interface{}
		expectError bool
	}{
		{
			name: "cast string field to number",
			castExpr: map[string]interface{}{
				"type":     "cast",
				"castType": "number",
				"expression": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "price",
				},
			},
			expected:    99.99,
			expectError: false,
		},
		{
			name: "cast string field to bool",
			castExpr: map[string]interface{}{
				"type":     "cast",
				"castType": "bool",
				"expression": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "active",
				},
			},
			expected:    true,
			expectError: false,
		},
		{
			name: "cast number field to string",
			castExpr: map[string]interface{}{
				"type":     "cast",
				"castType": "string",
				"expression": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "count",
				},
			},
			expected:    "10",
			expectError: false,
		},
		{
			name: "cast bool field to string",
			castExpr: map[string]interface{}{
				"type":     "cast",
				"castType": "string",
				"expression": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "flag",
				},
			},
			expected:    "true",
			expectError: false,
		},
		{
			name: "cast literal number to string",
			castExpr: map[string]interface{}{
				"type":     "cast",
				"castType": "string",
				"expression": map[string]interface{}{
					"type":  "number",
					"value": 123.0,
				},
			},
			expected:    "123",
			expectError: false,
		},
		{
			name: "cast literal string to number",
			castExpr: map[string]interface{}{
				"type":     "cast",
				"castType": "number",
				"expression": map[string]interface{}{
					"type":  "string",
					"value": "456",
				},
			},
			expected:    456.0,
			expectError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := eval.evaluateCastExpression(tt.castExpr)
			if tt.expectError {
				assert.Error(t, err, "Expected error for expression")
			} else {
				require.NoError(t, err, "Unexpected error for expression")
				assert.Equal(t, tt.expected, result, "Wrong expression result")
			}
		})
	}
}
// TestCastEdgeCases teste les cas limites et erreurs
func TestCastEdgeCases(t *testing.T) {
	t.Run("cast very large number to string", func(t *testing.T) {
		result, err := CastToString(1e10)
		require.NoError(t, err)
		assert.Equal(t, "10000000000", result)
	})
	t.Run("cast very small number to string", func(t *testing.T) {
		result, err := CastToString(0.0001)
		require.NoError(t, err)
		assert.Contains(t, result, "0.0001")
	})
	t.Run("cast zero to bool", func(t *testing.T) {
		result, err := CastToBool(0.0)
		require.NoError(t, err)
		assert.Equal(t, false, result)
	})
	t.Run("cast unicode string containing number", func(t *testing.T) {
		_, err := CastToNumber("１２３") // Full-width digits
		assert.Error(t, err)          // Should fail - only ASCII digits supported
	})
	t.Run("cast string with leading zeros", func(t *testing.T) {
		result, err := CastToNumber("0123")
		require.NoError(t, err)
		assert.Equal(t, 123.0, result)
	})
}