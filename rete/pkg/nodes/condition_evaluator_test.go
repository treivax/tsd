// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package nodes

import (
	"testing"

	"github.com/treivax/tsd/rete/pkg/domain"
)

func TestConditionEvaluator_EvaluateCondition_InvalidCondition(t *testing.T) {
	evaluator := NewConditionEvaluator()

	fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": 25})
	token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

	// Test with invalid condition type
	condition := "invalid_string"
	_, err := evaluator.EvaluateCondition(condition, token, fact)
	if err == nil {
		t.Error("EvaluateCondition should return error for invalid condition")
	}
}

func TestConditionEvaluator_EvaluateBinaryCondition_MissingOperator(t *testing.T) {
	evaluator := NewConditionEvaluator()

	fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": 25})
	token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

	// Condition without operator
	condition := map[string]interface{}{
		"type": "binaryOperation",
		"left": map[string]interface{}{
			"field": "age",
		},
		"right": map[string]interface{}{
			"value": 25,
		},
	}

	_, err := evaluator.EvaluateBinaryCondition(condition, token, fact)
	if err == nil {
		t.Error("EvaluateBinaryCondition should return error when operator is missing")
	}
}

func TestConditionEvaluator_EvaluateBinaryCondition_AllOperators(t *testing.T) {
	evaluator := NewConditionEvaluator()

	tests := []struct {
		name     string
		operator string
		leftVal  int
		rightVal int
		want     bool
	}{
		{"equal true", "==", 25, 25, true},
		{"equal false", "==", 25, 30, false},
		{"not equal true", "!=", 25, 30, true},
		{"not equal false", "!=", 25, 25, false},
		{"less than true", "<", 20, 25, true},
		{"less than false", "<", 30, 25, false},
		{"greater than true", ">", 30, 25, true},
		{"greater than false", ">", 20, 25, false},
		{"less or equal true", "<=", 25, 25, true},
		{"less or equal false", "<=", 30, 25, false},
		{"greater or equal true", ">=", 25, 25, true},
		{"greater or equal false", ">=", 20, 25, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": tt.leftVal})
			token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

			condition := map[string]interface{}{
				"type": "binaryOperation",
				"left": map[string]interface{}{
					"field": "age",
				},
				"operator": tt.operator,
				"right": map[string]interface{}{
					"value": tt.rightVal,
				},
			}

			result, err := evaluator.EvaluateBinaryCondition(condition, token, fact)
			if err != nil {
				t.Errorf("EvaluateBinaryCondition should not return error: %v", err)
			}
			if result != tt.want {
				t.Errorf("EvaluateBinaryCondition() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestConditionEvaluator_ExtractFieldValue_Missing(t *testing.T) {
	evaluator := NewConditionEvaluator()

	fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": 25})
	token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

	leftExpr := map[string]interface{}{
		"field": "nonexistent",
	}

	_, err := evaluator.ExtractFieldValue(leftExpr, token, fact)
	if err == nil {
		t.Error("ExtractFieldValue should return error for missing field")
	}
}

func TestConditionEvaluator_ExtractFieldValue_EmptyToken(t *testing.T) {
	evaluator := NewConditionEvaluator()

	fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": 25})
	token := domain.NewToken("t1", "node1", []*domain.Fact{}) // Empty facts

	leftExpr := map[string]interface{}{
		"field": "age",
	}

	_, err := evaluator.ExtractFieldValue(leftExpr, token, fact)
	if err == nil {
		t.Error("ExtractFieldValue should return error for empty token")
	}
}

func TestConditionEvaluator_ExtractFieldValue_Success(t *testing.T) {
	evaluator := NewConditionEvaluator()

	fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": 25, "name": "John"})
	token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

	leftExpr := map[string]interface{}{
		"field": "age",
	}

	value, err := evaluator.ExtractFieldValue(leftExpr, token, fact)
	if err != nil {
		t.Errorf("ExtractFieldValue should not return error: %v", err)
	}
	if value != 25 {
		t.Errorf("ExtractFieldValue() = %v, want 25", value)
	}
}

func TestConditionEvaluator_ExtractConstantValue(t *testing.T) {
	evaluator := NewConditionEvaluator()

	tests := []struct {
		name      string
		rightExpr interface{}
		want      interface{}
		wantErr   bool
	}{
		{
			name:      "map with value",
			rightExpr: map[string]interface{}{"value": 42},
			want:      42,
			wantErr:   false,
		},
		{
			name:      "direct value",
			rightExpr: 42,
			want:      42,
			wantErr:   false,
		},
		{
			name:      "string value",
			rightExpr: map[string]interface{}{"value": "test"},
			want:      "test",
			wantErr:   false,
		},
		{
			name:      "missing value - returns map directly",
			rightExpr: map[string]interface{}{"other": "field"},
			want:      nil,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := evaluator.ExtractConstantValue(tt.rightExpr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractConstantValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Special handling for map case
				if tt.name == "missing value - returns map directly" {
					if _, ok := value.(map[string]interface{}); !ok {
						t.Errorf("ExtractConstantValue() should return a map, got %T", value)
					}
				} else if value != tt.want {
					t.Errorf("ExtractConstantValue() = %v, want %v", value, tt.want)
				}
			}
		})
	}
}

func TestConditionEvaluator_CompareValues_DifferentTypes(t *testing.T) {
	evaluator := NewConditionEvaluator()

	// Compare int with string using == (converts to string)
	result, err := evaluator.CompareValues(42, "==", "42")
	if err != nil {
		t.Errorf("CompareValues should not return error for == operator: %v", err)
	}
	if !result {
		t.Error("CompareValues with == should return true for 42 and '42'")
	}

	// Compare int with string using < (requires numeric conversion)
	_, err2 := evaluator.CompareValues(42, "<", "not_a_number")
	if err2 == nil {
		t.Error("CompareValues should return error for < with non-numeric string")
	}
}

func TestConditionEvaluator_CompareValues_StringComparison(t *testing.T) {
	evaluator := NewConditionEvaluator()

	tests := []struct {
		name     string
		left     interface{}
		operator string
		right    interface{}
		want     bool
		wantErr  bool
	}{
		{"string equal true", "abc", "==", "abc", true, false},
		{"string equal false", "abc", "==", "xyz", false, false},
		{"string not equal", "abc", "!=", "xyz", true, false},
		{"int to string", 42, "==", "42", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.CompareValues(tt.left, tt.operator, tt.right)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.want {
				t.Errorf("CompareValues() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestConditionEvaluator_CompareValues_FloatComparison(t *testing.T) {
	evaluator := NewConditionEvaluator()

	result, err := evaluator.CompareValues(3.14, "<", 3.15)
	if err != nil {
		t.Errorf("CompareValues should not return error: %v", err)
	}
	if !result {
		t.Error("3.14 < 3.15 should be true")
	}
}

func TestConditionEvaluator_CompareValues_InvalidOperator(t *testing.T) {
	evaluator := NewConditionEvaluator()

	_, err := evaluator.CompareValues(42, "INVALID", 42)
	if err == nil {
		t.Error("CompareValues should return error for invalid operator")
	}
}

func TestConditionEvaluator_NumericCompare(t *testing.T) {
	evaluator := NewConditionEvaluator()

	tests := []struct {
		name    string
		left    interface{}
		right   interface{}
		cmpFunc func(float64, float64) bool
		want    bool
		wantErr bool
	}{
		{"int comparison", 10, 20, func(l, r float64) bool { return l < r }, true, false},
		{"float comparison", 3.14, 2.71, func(l, r float64) bool { return l > r }, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.NumericCompare(tt.left, tt.right, tt.cmpFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("NumericCompare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.want {
				t.Errorf("NumericCompare() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestConditionEvaluator_ToFloat64(t *testing.T) {
	evaluator := NewConditionEvaluator()

	tests := []struct {
		name    string
		value   interface{}
		want    float64
		wantErr bool
	}{
		{"int", 42, 42.0, false},
		{"int64", int64(42), 42.0, false},
		{"bool", true, 0, true},
		{"nil", nil, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.ToFloat64(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.want {
				t.Errorf("ToFloat64() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestConditionEvaluator_EvaluateCondition_BinaryOperation(t *testing.T) {
	evaluator := NewConditionEvaluator()

	fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": 30})
	token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

	condition := map[string]interface{}{
		"type": "binaryOperation",
		"left": map[string]interface{}{
			"field": "age",
		},
		"operator": "==",
		"right": map[string]interface{}{
			"value": 30,
		},
	}

	result, err := evaluator.EvaluateCondition(condition, token, fact)
	if err != nil {
		t.Errorf("EvaluateCondition should not return error: %v", err)
	}
	if !result {
		t.Error("EvaluateCondition should return true for matching condition")
	}
}

func TestConditionEvaluator_EvaluateCondition_SimpleType(t *testing.T) {
	evaluator := NewConditionEvaluator()

	fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": 30})
	token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

	condition := map[string]interface{}{
		"type": "simple",
		"left": map[string]interface{}{
			"field": "age",
		},
		"operator": ">",
		"right": map[string]interface{}{
			"value": 25,
		},
	}

	result, err := evaluator.EvaluateCondition(condition, token, fact)
	if err != nil {
		t.Errorf("EvaluateCondition should not return error: %v", err)
	}
	if !result {
		t.Error("EvaluateCondition should return true for simple condition")
	}
}
