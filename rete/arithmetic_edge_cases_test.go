// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"math"
	"testing"
)
// TestArithmeticEdgeCases tests edge cases in arithmetic expression evaluation
func TestArithmeticEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		left      interface{}
		operator  string
		right     interface{}
		expected  interface{}
		shouldErr bool
		errMsg    string
	}{
		// Division by zero cases
		{
			name:      "Integer division by zero",
			left:      10,
			operator:  "/",
			right:     0,
			expected:  nil,
			shouldErr: true,
			errMsg:    "division by zero",
		},
		{
			name:      "Float division by zero",
			left:      10.5,
			operator:  "/",
			right:     0.0,
			expected:  nil,
			shouldErr: true,
			errMsg:    "division by zero",
		},
		{
			name:      "Modulo by zero",
			left:      10,
			operator:  "%",
			right:     0,
			expected:  nil,
			shouldErr: true,
			errMsg:    "modulo by zero",
		},
		// Type conversion cases
		{
			name:      "Int + Float",
			left:      5,
			operator:  "+",
			right:     2.5,
			expected:  7.5,
			shouldErr: false,
		},
		{
			name:      "Float + Int",
			left:      2.5,
			operator:  "+",
			right:     5,
			expected:  7.5,
			shouldErr: false,
		},
		{
			name:      "Int * Float",
			left:      3,
			operator:  "*",
			right:     2.5,
			expected:  7.5,
			shouldErr: false,
		},
		{
			name:      "Float / Int",
			left:      10.0,
			operator:  "/",
			right:     4,
			expected:  2.5,
			shouldErr: false,
		},
		// Large numbers
		{
			name:      "Large int addition",
			left:      int64(1000000000),
			operator:  "+",
			right:     int64(2000000000),
			expected:  float64(3000000000),
			shouldErr: false,
		},
		{
			name:      "Large float multiplication",
			left:      1e10,
			operator:  "*",
			right:     1e10,
			expected:  1e20,
			shouldErr: false,
		},
		// Negative numbers
		{
			name:      "Negative addition",
			left:      -5,
			operator:  "+",
			right:     -3,
			expected:  float64(-8),
			shouldErr: false,
		},
		{
			name:      "Negative multiplication",
			left:      -5,
			operator:  "*",
			right:     3,
			expected:  float64(-15),
			shouldErr: false,
		},
		{
			name:      "Negative division",
			left:      -10,
			operator:  "/",
			right:     2,
			expected:  float64(-5),
			shouldErr: false,
		},
		{
			name:      "Double negative multiplication",
			left:      -5,
			operator:  "*",
			right:     -3,
			expected:  float64(15),
			shouldErr: false,
		},
		// Zero cases
		{
			name:      "Zero addition",
			left:      0,
			operator:  "+",
			right:     5,
			expected:  float64(5),
			shouldErr: false,
		},
		{
			name:      "Zero multiplication",
			left:      5,
			operator:  "*",
			right:     0,
			expected:  float64(0),
			shouldErr: false,
		},
		{
			name:      "Zero division (0/x)",
			left:      0,
			operator:  "/",
			right:     5,
			expected:  float64(0),
			shouldErr: false,
		},
		{
			name:      "Zero modulo (0%x)",
			left:      0,
			operator:  "%",
			right:     5,
			expected:  float64(0),
			shouldErr: false,
		},
		// Modulo edge cases
		{
			name:      "Modulo with float (converts to int)",
			left:      10.5,
			operator:  "%",
			right:     3.0,
			expected:  float64(1), // Implementation converts to int64: 10 % 3 = 1
			shouldErr: false,
		},
		{
			name:      "Negative modulo",
			left:      -10,
			operator:  "%",
			right:     3,
			expected:  float64(-1),
			shouldErr: false,
		},
		// Subtraction edge cases
		{
			name:      "Subtraction resulting in zero",
			left:      5,
			operator:  "-",
			right:     5,
			expected:  float64(0),
			shouldErr: false,
		},
		{
			name:      "Subtraction resulting in negative",
			left:      3,
			operator:  "-",
			right:     5,
			expected:  float64(-2),
			shouldErr: false,
		},
		// Very small numbers
		{
			name:      "Small float addition",
			left:      0.0000001,
			operator:  "+",
			right:     0.0000002,
			expected:  0.0000003,
			shouldErr: false,
		},
		{
			name:      "Small float multiplication",
			left:      0.1,
			operator:  "*",
			right:     0.1,
			expected:  0.01,
			shouldErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewAlphaConditionEvaluator()
			result, err := evaluator.evaluateArithmeticOperation(tt.left, tt.operator, tt.right)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				// For floating point comparisons, use approximate equality
				expectedFloat, okExp := toFloat64(tt.expected)
				resultFloat, okRes := toFloat64(result)
				if !okExp || !okRes {
					t.Errorf("Could not convert to float64: expected=%v (%T), result=%v (%T)",
						tt.expected, tt.expected, result, result)
					return
				}
				if !floatEquals(expectedFloat, resultFloat, 1e-9) {
					t.Errorf("Expected %v, got %v (diff: %e)",
						expectedFloat, resultFloat, math.Abs(expectedFloat-resultFloat))
				}
			}
		})
	}
}
// TestArithmeticWithInvalidTypes tests arithmetic operations with invalid types
func TestArithmeticWithInvalidTypes(t *testing.T) {
	tests := []struct {
		name     string
		left     interface{}
		operator string
		right    interface{}
	}{
		{
			name:     "String + Number",
			left:     "hello",
			operator: "+",
			right:    5,
		},
		{
			name:     "Number + String",
			left:     5,
			operator: "+",
			right:    "hello",
		},
		{
			name:     "Boolean + Number",
			left:     true,
			operator: "+",
			right:    5,
		},
		{
			name:     "Nil + Number",
			left:     nil,
			operator: "+",
			right:    5,
		},
		{
			name:     "Number + Nil",
			left:     5,
			operator: "+",
			right:    nil,
		},
		{
			name:     "String * Number",
			left:     "test",
			operator: "*",
			right:    3,
		},
		{
			name:     "Array + Number",
			left:     []interface{}{1, 2, 3},
			operator: "+",
			right:    5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewAlphaConditionEvaluator()
			_, err := evaluator.evaluateArithmeticOperation(tt.left, tt.operator, tt.right)
			if err == nil {
				t.Errorf("Expected error for invalid type combination, got nil")
			}
		})
	}
}
// TestComplexNestedArithmetic tests nested arithmetic expressions
func TestComplexNestedArithmetic(t *testing.T) {
	// Create a simple evaluator with bindings
	evaluator := NewAlphaConditionEvaluator()
	// Create a fact with numeric fields
	fact := &Fact{
		ID:   "f1",
		Type: "TestType",
		Fields: map[string]interface{}{
			"a": 10.0,
			"b": 5.0,
			"c": 2.0,
		},
	}
	// Bind the fact by setting variableBindings directly
	evaluator.variableBindings = map[string]*Fact{
		"x": fact,
	}
	tests := []struct {
		name     string
		expr     map[string]interface{}
		expected float64
	}{
		{
			name: "Nested: (a + b) * c",
			expr: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "*",
				"left": map[string]interface{}{
					"type":     "binaryOp",
					"operator": "+",
					"left": map[string]interface{}{
						"type":     "fieldAccess",
						"variable": "x",
						"field":    "a",
					},
					"right": map[string]interface{}{
						"type":     "fieldAccess",
						"variable": "x",
						"field":    "b",
					},
				},
				"right": map[string]interface{}{
					"type":     "fieldAccess",
					"variable": "x",
					"field":    "c",
				},
			},
			expected: 30.0, // (10 + 5) * 2 = 30
		},
		{
			name: "Nested: a * (b + c)",
			expr: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "*",
				"left": map[string]interface{}{
					"type":     "fieldAccess",
					"variable": "x",
					"field":    "a",
				},
				"right": map[string]interface{}{
					"type":     "binaryOp",
					"operator": "+",
					"left": map[string]interface{}{
						"type":     "fieldAccess",
						"variable": "x",
						"field":    "b",
					},
					"right": map[string]interface{}{
						"type":     "fieldAccess",
						"variable": "x",
						"field":    "c",
					},
				},
			},
			expected: 70.0, // 10 * (5 + 2) = 70
		},
		{
			name: "Triple nested: (a + b) * (c + 1)",
			expr: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "*",
				"left": map[string]interface{}{
					"type":     "binaryOp",
					"operator": "+",
					"left": map[string]interface{}{
						"type":     "fieldAccess",
						"variable": "x",
						"field":    "a",
					},
					"right": map[string]interface{}{
						"type":     "fieldAccess",
						"variable": "x",
						"field":    "b",
					},
				},
				"right": map[string]interface{}{
					"type":     "binaryOp",
					"operator": "+",
					"left": map[string]interface{}{
						"type":     "fieldAccess",
						"variable": "x",
						"field":    "c",
					},
					"right": map[string]interface{}{
						"type":  "number",
						"value": 1.0,
					},
				},
			},
			expected: 45.0, // (10 + 5) * (2 + 1) = 45
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateValue(tt.expr)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			resultFloat, ok := toFloat64(result)
			if !ok {
				t.Fatalf("Could not convert result to float64: %v (%T)", result, result)
			}
			if !floatEquals(resultFloat, tt.expected, 1e-9) {
				t.Errorf("Expected %v, got %v", tt.expected, resultFloat)
			}
		})
	}
}
// TestArithmeticOperatorPrecedence tests that operator precedence is respected
func TestArithmeticOperatorPrecedence(t *testing.T) {
	// Note: Precedence should be handled by the parser, not the evaluator
	// The evaluator evaluates the parsed tree structure
	// This test verifies that properly structured trees evaluate correctly
	evaluator := NewAlphaConditionEvaluator()
	// Test: 2 + 3 * 4 should be 2 + (3 * 4) = 14, not (2 + 3) * 4 = 20
	// The tree structure should reflect this precedence
	expr := map[string]interface{}{
		"type":     "binaryOp",
		"operator": "+",
		"left": map[string]interface{}{
			"type":  "number",
			"value": 2.0,
		},
		"right": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "*",
			"left": map[string]interface{}{
				"type":  "number",
				"value": 3.0,
			},
			"right": map[string]interface{}{
				"type":  "number",
				"value": 4.0,
			},
		},
	}
	result, err := evaluator.evaluateValue(expr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := 14.0
	resultFloat, ok := toFloat64(result)
	if !ok {
		t.Fatalf("Could not convert result to float64: %v (%T)", result, result)
	}
	if !floatEquals(resultFloat, expected, 1e-9) {
		t.Errorf("Expected %v, got %v", expected, resultFloat)
	}
}
// Helper functions
func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	default:
		return 0, false
	}
}
func floatEquals(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}