// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"
)

// TestConditionEvaluator_BinaryOp tests arithmetic operations
func TestConditionEvaluator_BinaryOp(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{"qte": 10, "remise": 5},
	}
	ctx := NewEvaluationContext(fact)
	evaluator := NewConditionEvaluator(nil)
	tests := []struct {
		name     string
		op       string
		left     interface{}
		right    interface{}
		expected float64
	}{
		{"multiplication", "*", 10.0, 23.0, 230.0},
		{"addition", "+", 100.0, 50.0, 150.0},
		{"subtraction", "-", 100.0, 10.0, 90.0},
		{"division", "/", 100.0, 4.0, 25.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := map[string]interface{}{
				"type":     "binaryOp",
				"operator": tt.op,
				"left":     map[string]interface{}{"type": "number", "value": tt.left},
				"right":    map[string]interface{}{"type": "number", "value": tt.right},
			}
			result, err := evaluator.EvaluateWithContext(condition, fact, ctx)
			if err != nil {
				t.Fatalf("Error evaluating: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestConditionEvaluator_Comparison tests comparison operations
func TestConditionEvaluator_Comparison(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	evaluator := NewConditionEvaluator(nil)
	tests := []struct {
		name     string
		op       string
		left     float64
		right    float64
		expected bool
	}{
		{"greater than true", ">", 100.0, 50.0, true},
		{"greater than false", ">", 50.0, 100.0, false},
		{"less than true", "<", 50.0, 100.0, true},
		{"less than false", "<", 100.0, 50.0, false},
		{"equal true", "==", 100.0, 100.0, true},
		{"equal false", "==", 100.0, 50.0, false},
		{"not equal true", "!=", 100.0, 50.0, true},
		{"not equal false", "!=", 100.0, 100.0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := map[string]interface{}{
				"type":     "comparison",
				"operator": tt.op,
				"left":     map[string]interface{}{"type": "number", "value": tt.left},
				"right":    map[string]interface{}{"type": "number", "value": tt.right},
			}
			result, err := evaluator.EvaluateWithContext(condition, fact, ctx)
			if err != nil {
				t.Fatalf("Error evaluating: %v", err)
			}
			boolResult, ok := result.(bool)
			if !ok {
				t.Fatalf("Result should be bool, got %T", result)
			}
			if boolResult != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, boolResult)
			}
		})
	}
}

// TestConditionEvaluator_TempResult tests intermediate result resolution
func TestConditionEvaluator_TempResult(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	evaluator := NewConditionEvaluator(nil)
	// Store intermediate results
	ctx.SetIntermediateResult("temp_1", 42.0)
	ctx.SetIntermediateResult("temp_2", 84.0)
	// Test resolving temp_1
	condition := map[string]interface{}{
		"type":      "tempResult",
		"step_name": "temp_1",
	}
	result, err := evaluator.EvaluateWithContext(condition, fact, ctx)
	if err != nil {
		t.Fatalf("Error evaluating: %v", err)
	}
	if result != 42.0 {
		t.Errorf("Expected 42.0, got %v", result)
	}
	// Test resolving temp_2
	condition["step_name"] = "temp_2"
	result, err = evaluator.EvaluateWithContext(condition, fact, ctx)
	if err != nil {
		t.Fatalf("Error evaluating: %v", err)
	}
	if result != 84.0 {
		t.Errorf("Expected 84.0, got %v", result)
	}
}

// TestConditionEvaluator_MissingDependency tests error handling for missing dependencies
func TestConditionEvaluator_MissingDependency(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	evaluator := NewConditionEvaluator(nil)
	// Try to resolve a non-existent temp result
	condition := map[string]interface{}{
		"type":      "tempResult",
		"step_name": "nonexistent",
	}
	_, err := evaluator.EvaluateWithContext(condition, fact, ctx)
	if err == nil {
		t.Fatal("Expected error for missing dependency, got nil")
	}
}

// TestConditionEvaluator_NestedExpression tests complex nested expressions
func TestConditionEvaluator_NestedExpression(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{"qte": 5, "remise": 10},
	}
	ctx := NewEvaluationContext(fact)
	evaluator := NewConditionEvaluator(nil)
	// Expression: (10 * 23) - 10 = 230 - 10 = 220
	condition := map[string]interface{}{
		"type":     "binaryOp",
		"operator": "-",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "*",
			"left":     map[string]interface{}{"type": "number", "value": 10.0},
			"right":    map[string]interface{}{"type": "number", "value": 23.0},
		},
		"right": map[string]interface{}{"type": "number", "value": 10.0},
	}
	result, err := evaluator.EvaluateWithContext(condition, fact, ctx)
	if err != nil {
		t.Fatalf("Error evaluating: %v", err)
	}
	expected := 220.0
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestConditionEvaluator_FieldAccess tests field value extraction
func TestConditionEvaluator_FieldAccess(t *testing.T) {
	fact := &Fact{
		ID:   "test1",
		Type: "TestType",
		Fields: map[string]interface{}{
			"qte":    10,
			"price":  42.5,
			"active": true,
		},
	}
	ctx := NewEvaluationContext(fact)
	evaluator := NewConditionEvaluator(nil)
	tests := []struct {
		field    string
		expected interface{}
	}{
		{"qte", 10},
		{"price", 42.5},
		{"active", true},
	}
	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			condition := map[string]interface{}{
				"type":  "fieldAccess",
				"field": tt.field,
			}
			result, err := evaluator.EvaluateWithContext(condition, fact, ctx)
			if err != nil {
				t.Fatalf("Error evaluating: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestConditionEvaluator_WithTempResults tests operations using temp results
func TestConditionEvaluator_WithTempResults(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	evaluator := NewConditionEvaluator(nil)
	// Simulate a chain: temp_1 = 100, temp_2 = temp_1 + 50
	ctx.SetIntermediateResult("temp_1", 100.0)
	condition := map[string]interface{}{
		"type":     "binaryOp",
		"operator": "+",
		"left":     map[string]interface{}{"type": "tempResult", "step_name": "temp_1"},
		"right":    map[string]interface{}{"type": "number", "value": 50.0},
	}
	result, err := evaluator.EvaluateWithContext(condition, fact, ctx)
	if err != nil {
		t.Fatalf("Error evaluating: %v", err)
	}
	expected := 150.0
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
	// Store as temp_2
	ctx.SetIntermediateResult("temp_2", result)
	// Now use temp_2 in another operation: temp_3 = temp_2 * 2
	condition = map[string]interface{}{
		"type":     "binaryOp",
		"operator": "*",
		"left":     map[string]interface{}{"type": "tempResult", "step_name": "temp_2"},
		"right":    map[string]interface{}{"type": "number", "value": 2.0},
	}
	result, err = evaluator.EvaluateWithContext(condition, fact, ctx)
	if err != nil {
		t.Fatalf("Error evaluating: %v", err)
	}
	expected = 300.0
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestConditionEvaluator_DivisionByZero tests division by zero error handling
func TestConditionEvaluator_DivisionByZero(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	evaluator := NewConditionEvaluator(nil)
	condition := map[string]interface{}{
		"type":     "binaryOp",
		"operator": "/",
		"left":     map[string]interface{}{"type": "number", "value": 100.0},
		"right":    map[string]interface{}{"type": "number", "value": 0.0},
	}
	_, err := evaluator.EvaluateWithContext(condition, fact, ctx)
	if err == nil {
		t.Fatal("Expected error for division by zero, got nil")
	}
}
