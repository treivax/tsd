// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// TestArithmeticDecomposer_SimpleExpression tests simple expression decomposition
func TestArithmeticDecomposer_SimpleExpression(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()

	// Expression: c.qte * 23
	expr := map[string]interface{}{
		"type":     "binaryOp",
		"operator": "*",
		"left": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "qte",
		},
		"right": map[string]interface{}{
			"type":  "number",
			"value": 23,
		},
	}

	steps, err := decomposer.DecomposeToDecomposedConditions(expr)
	if err != nil {
		t.Fatalf("Error decomposing: %v", err)
	}

	// Should produce 1 step
	if len(steps) != 1 {
		t.Errorf("Expected 1 step, got %d", len(steps))
	}

	// Verify step 1
	step1 := steps[0]
	if step1.ResultName != "temp_1" {
		t.Errorf("Expected ResultName 'temp_1', got '%s'", step1.ResultName)
	}

	if !step1.IsAtomic {
		t.Error("Step should be atomic")
	}

	if step1.Operator != "*" {
		t.Errorf("Expected operator '*', got '%s'", step1.Operator)
	}

	if len(step1.Dependencies) != 0 {
		t.Errorf("Step 1 should have no dependencies, got %v", step1.Dependencies)
	}
}

// TestArithmeticDecomposer_ComplexExpression tests complex nested expression
func TestArithmeticDecomposer_ComplexExpression(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()

	// Expression: (c.qte * 23 - 10 + c.remise * 43) > 0
	// Should decompose into:
	// Step 1: c.qte * 23 → temp_1
	// Step 2: temp_1 - 10 → temp_2
	// Step 3: c.remise * 43 → temp_3
	// Step 4: temp_2 + temp_3 → temp_4
	// Step 5: temp_4 > 0 → temp_5

	expr := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "+",
			"left": map[string]interface{}{
				"type":     "binaryOp",
				"operator": "-",
				"left": map[string]interface{}{
					"type":     "binaryOp",
					"operator": "*",
					"left": map[string]interface{}{
						"type":  "fieldAccess",
						"field": "qte",
					},
					"right": map[string]interface{}{
						"type":  "number",
						"value": 23,
					},
				},
				"right": map[string]interface{}{
					"type":  "number",
					"value": 10,
				},
			},
			"right": map[string]interface{}{
				"type":     "binaryOp",
				"operator": "*",
				"left": map[string]interface{}{
					"type":  "fieldAccess",
					"field": "remise",
				},
				"right": map[string]interface{}{
					"type":  "number",
					"value": 43,
				},
			},
		},
		"right": map[string]interface{}{
			"type":  "number",
			"value": 0,
		},
	}

	steps, err := decomposer.DecomposeToDecomposedConditions(expr)
	if err != nil {
		t.Fatalf("Error decomposing: %v", err)
	}

	// Should produce 5 steps
	if len(steps) != 5 {
		t.Errorf("Expected 5 steps, got %d", len(steps))
		for i, step := range steps {
			t.Logf("Step %d: %s (deps: %v)", i+1, step.ResultName, step.Dependencies)
		}
	}

	// Verify step 1: c.qte * 23 → temp_1
	if steps[0].ResultName != "temp_1" {
		t.Errorf("Step 1: Expected temp_1, got %s", steps[0].ResultName)
	}
	if steps[0].Operator != "*" {
		t.Errorf("Step 1: Expected operator *, got %s", steps[0].Operator)
	}
	if len(steps[0].Dependencies) != 0 {
		t.Errorf("Step 1: Expected no dependencies, got %v", steps[0].Dependencies)
	}

	// Verify step 2: temp_1 - 10 → temp_2
	if steps[1].ResultName != "temp_2" {
		t.Errorf("Step 2: Expected temp_2, got %s", steps[1].ResultName)
	}
	if steps[1].Operator != "-" {
		t.Errorf("Step 2: Expected operator -, got %s", steps[1].Operator)
	}
	if len(steps[1].Dependencies) != 1 || steps[1].Dependencies[0] != "temp_1" {
		t.Errorf("Step 2: Expected dependency on temp_1, got %v", steps[1].Dependencies)
	}

	// Verify step 2 left operand is tempResult reference
	leftMap, ok := steps[1].Left.(map[string]interface{})
	if !ok {
		t.Fatal("Step 2 left should be a map")
	}
	if leftMap["type"] != "tempResult" {
		t.Errorf("Step 2 left should be tempResult, got %s", leftMap["type"])
	}
	if leftMap["step_name"] != "temp_1" {
		t.Errorf("Step 2 left should reference temp_1, got %s", leftMap["step_name"])
	}

	// Verify step 3: c.remise * 43 → temp_3
	if steps[2].ResultName != "temp_3" {
		t.Errorf("Step 3: Expected temp_3, got %s", steps[2].ResultName)
	}
	if steps[2].Operator != "*" {
		t.Errorf("Step 3: Expected operator *, got %s", steps[2].Operator)
	}

	// Verify step 4: temp_2 + temp_3 → temp_4
	if steps[3].ResultName != "temp_4" {
		t.Errorf("Step 4: Expected temp_4, got %s", steps[3].ResultName)
	}
	if steps[3].Operator != "+" {
		t.Errorf("Step 4: Expected operator +, got %s", steps[3].Operator)
	}
	if len(steps[3].Dependencies) != 2 {
		t.Errorf("Step 4: Expected 2 dependencies, got %v", steps[3].Dependencies)
	}

	// Verify step 5: temp_4 > 0 → temp_5 (comparison)
	if steps[4].ResultName != "temp_5" {
		t.Errorf("Step 5: Expected temp_5, got %s", steps[4].ResultName)
	}
	if steps[4].Type != "comparison" {
		t.Errorf("Step 5: Expected type comparison, got %s", steps[4].Type)
	}
	if steps[4].Operator != ">" {
		t.Errorf("Step 5: Expected operator >, got %s", steps[4].Operator)
	}
}

// TestArithmeticDecomposer_ShouldDecompose tests complexity threshold
func TestArithmeticDecomposer_ShouldDecompose(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()

	// Simple expression: c.qte * 23
	simpleExpr := map[string]interface{}{
		"type":     "binaryOp",
		"operator": "*",
		"left":     map[string]interface{}{"type": "fieldAccess", "field": "qte"},
		"right":    map[string]interface{}{"type": "number", "value": 23},
	}

	// Complexity = 1 (should NOT decompose, needs > 1)
	if decomposer.ShouldDecompose(simpleExpr) {
		t.Error("Simple expression should not decompose (only 1 operation, needs > 1)")
	}

	// Complex expression: (a * b) + (c * d)
	complexExpr := map[string]interface{}{
		"type":     "binaryOp",
		"operator": "+",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "*",
			"left":     map[string]interface{}{"type": "fieldAccess", "field": "a"},
			"right":    map[string]interface{}{"type": "fieldAccess", "field": "b"},
		},
		"right": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "*",
			"left":     map[string]interface{}{"type": "fieldAccess", "field": "c"},
			"right":    map[string]interface{}{"type": "fieldAccess", "field": "d"},
		},
	}

	// Complexity = 3 (two multiplications + one addition)
	complexity := decomposer.GetComplexity(complexExpr)
	if complexity != 3 {
		t.Errorf("Expected complexity 3, got %d", complexity)
	}

	// Should decompose
	if !decomposer.ShouldDecompose(complexExpr) {
		t.Error("Complex expression should decompose")
	}
}

// TestArithmeticDecomposer_Dependencies tests dependency extraction
func TestArithmeticDecomposer_Dependencies(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()

	// Expression: (a + b) * (c - d)
	// Should create:
	// Step 1: a + b → temp_1
	// Step 2: c - d → temp_2
	// Step 3: temp_1 * temp_2 → temp_3 (depends on temp_1 and temp_2)

	expr := map[string]interface{}{
		"type":     "binaryOp",
		"operator": "*",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "+",
			"left":     map[string]interface{}{"type": "fieldAccess", "field": "a"},
			"right":    map[string]interface{}{"type": "fieldAccess", "field": "b"},
		},
		"right": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "-",
			"left":     map[string]interface{}{"type": "fieldAccess", "field": "c"},
			"right":    map[string]interface{}{"type": "fieldAccess", "field": "d"},
		},
	}

	steps, err := decomposer.DecomposeToDecomposedConditions(expr)
	if err != nil {
		t.Fatalf("Error decomposing: %v", err)
	}

	if len(steps) != 3 {
		t.Fatalf("Expected 3 steps, got %d", len(steps))
	}

	// Step 3 should depend on temp_1 and temp_2
	step3 := steps[2]
	if len(step3.Dependencies) != 2 {
		t.Errorf("Step 3 should have 2 dependencies, got %v", step3.Dependencies)
	}

	expectedDeps := map[string]bool{"temp_1": true, "temp_2": true}
	for _, dep := range step3.Dependencies {
		if !expectedDeps[dep] {
			t.Errorf("Unexpected dependency: %s", dep)
		}
	}
}

// TestArithmeticDecomposer_DecomposeExpression tests legacy API
func TestArithmeticDecomposer_DecomposeExpression(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()

	expr := map[string]interface{}{
		"type":     "binaryOp",
		"operator": "*",
		"left":     map[string]interface{}{"type": "fieldAccess", "field": "qte"},
		"right":    map[string]interface{}{"type": "number", "value": 23},
	}

	decomposed, err := decomposer.DecomposeExpression(expr)
	if err != nil {
		t.Fatalf("Error decomposing: %v", err)
	}

	if decomposed == nil {
		t.Fatal("Expected decomposed expression, got nil")
	}

	if len(decomposed.Steps) != 1 {
		t.Errorf("Expected 1 step, got %d", len(decomposed.Steps))
	}
}

// BenchmarkArithmeticDecomposer_SimpleExpression benchmarks simple decomposition
func BenchmarkArithmeticDecomposer_SimpleExpression(b *testing.B) {
	decomposer := NewArithmeticExpressionDecomposer()

	expr := map[string]interface{}{
		"type":     "binaryOp",
		"operator": "*",
		"left":     map[string]interface{}{"type": "fieldAccess", "field": "qte"},
		"right":    map[string]interface{}{"type": "number", "value": 23},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = decomposer.DecomposeToDecomposedConditions(expr)
	}
}

// BenchmarkArithmeticDecomposer_ComplexExpression benchmarks complex decomposition
func BenchmarkArithmeticDecomposer_ComplexExpression(b *testing.B) {
	decomposer := NewArithmeticExpressionDecomposer()

	// (a * b - c) + (d * e)
	expr := map[string]interface{}{
		"type":     "binaryOp",
		"operator": "+",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "-",
			"left": map[string]interface{}{
				"type":     "binaryOp",
				"operator": "*",
				"left":     map[string]interface{}{"type": "fieldAccess", "field": "a"},
				"right":    map[string]interface{}{"type": "fieldAccess", "field": "b"},
			},
			"right": map[string]interface{}{"type": "fieldAccess", "field": "c"},
		},
		"right": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "*",
			"left":     map[string]interface{}{"type": "fieldAccess", "field": "d"},
			"right":    map[string]interface{}{"type": "fieldAccess", "field": "e"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = decomposer.DecomposeToDecomposedConditions(expr)
	}
}
