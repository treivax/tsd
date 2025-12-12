package rete

import (
	"testing"
)

// TestDecomposeRecursive tests the decomposeRecursive function coverage
func TestDecomposeRecursive(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()
	tests := []struct {
		name      string
		expr      interface{}
		wantSteps int
		wantErr   bool
	}{
		{
			name: "simple number literal",
			expr: map[string]interface{}{
				"type":  "number",
				"value": 42,
			},
			wantSteps: 0,    // No decomposition needed for literals
			wantErr:   true, // Will error because no steps extracted
		},
		{
			name: "simple addition",
			expr: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "+",
				"left":     map[string]interface{}{"type": "number", "value": 1},
				"right":    map[string]interface{}{"type": "number", "value": 2},
			},
			wantSteps: 1,
			wantErr:   false,
		},
		{
			name: "nested arithmetic",
			expr: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "+",
				"left": map[string]interface{}{
					"type":     "binaryOp",
					"operator": "*",
					"left":     map[string]interface{}{"type": "number", "value": 2},
					"right":    map[string]interface{}{"type": "number", "value": 3},
				},
				"right": map[string]interface{}{"type": "number", "value": 4},
			},
			wantSteps: 2, // Two operations
			wantErr:   false,
		},
		{
			name: "field access",
			expr: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "x",
				"field":  "value",
			},
			wantSteps: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := decomposer.DecomposeExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecomposeExpression() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && len(result.Steps) != tt.wantSteps {
				t.Errorf("DecomposeExpression() got %d steps, want %d", len(result.Steps), tt.wantSteps)
			}
		})
	}
}

// TestSimplifySteps tests the SimplifySteps function (currently 0% coverage)
func TestSimplifySteps(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()
	// Test with simple expression
	expr := map[string]interface{}{
		"type":     "binaryOp",
		"operator": "+",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "*",
			"left":     map[string]interface{}{"type": "number", "value": 2},
			"right":    map[string]interface{}{"type": "number", "value": 3},
		},
		"right": map[string]interface{}{"type": "number", "value": 4},
	}
	result, err := decomposer.DecomposeExpression(expr)
	if err != nil {
		t.Fatalf("DecomposeExpression failed: %v", err)
	}
	// Call SimplifySteps
	simplified := decomposer.SimplifySteps(result.Steps)
	// Verify simplified steps exist
	if simplified == nil {
		t.Error("SimplifySteps returned nil")
	}
	if len(simplified) != len(result.Steps) {
		t.Errorf("SimplifySteps changed number of steps: got %d, want %d", len(simplified), len(result.Steps))
	}
}

// TestValidateDecomposition tests the ValidateDecomposition function (currently 0% coverage)
func TestValidateDecomposition(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()
	tests := []struct {
		name    string
		expr    interface{}
		wantErr bool
	}{
		{
			name: "valid simple expression",
			expr: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "+",
				"left":     map[string]interface{}{"type": "number", "value": 1},
				"right":    map[string]interface{}{"type": "number", "value": 2},
			},
			wantErr: false,
		},
		{
			name: "valid nested expression",
			expr: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "*",
				"left": map[string]interface{}{
					"type":     "binaryOp",
					"operator": "+",
					"left":     map[string]interface{}{"type": "number", "value": 1},
					"right":    map[string]interface{}{"type": "number", "value": 2},
				},
				"right": map[string]interface{}{"type": "number", "value": 3},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := decomposer.DecomposeExpression(tt.expr)
			if err != nil {
				t.Fatalf("DecomposeExpression failed: %v", err)
			}
			err = decomposer.ValidateDecomposition(tt.expr, result)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDecomposition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestFormatSteps tests the FormatSteps function (currently 0% coverage)
func TestFormatSteps(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()
	expr := map[string]interface{}{
		"type":     "binaryOp",
		"operator": "+",
		"left":     map[string]interface{}{"type": "number", "value": 1},
		"right":    map[string]interface{}{"type": "number", "value": 2},
	}
	result, err := decomposer.DecomposeExpression(expr)
	if err != nil {
		t.Fatalf("DecomposeExpression failed: %v", err)
	}
	formatted := decomposer.FormatSteps(result.Steps)
	if formatted == "" {
		t.Error("FormatSteps returned empty string")
	}
	if len(formatted) == 0 {
		t.Error("FormatSteps should return non-empty string for decomposed expression")
	}
}

// TestGetDecompositionStats tests the GetDecompositionStats function (currently 0% coverage)
func TestGetDecompositionStats(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()
	// Perform some decompositions
	expressions := []interface{}{
		map[string]interface{}{
			"type":     "binaryOp",
			"operator": "+",
			"left":     map[string]interface{}{"type": "number", "value": 1},
			"right":    map[string]interface{}{"type": "number", "value": 2},
		},
		map[string]interface{}{
			"type":     "binaryOp",
			"operator": "*",
			"left":     map[string]interface{}{"type": "number", "value": 3},
			"right":    map[string]interface{}{"type": "number", "value": 4},
		},
	}
	for _, expr := range expressions {
		result, err := decomposer.DecomposeExpression(expr)
		if err != nil {
			t.Fatalf("DecomposeExpression failed: %v", err)
		}
		stats := decomposer.GetDecompositionStats(result)
		if stats == nil {
			t.Error("GetDecompositionStats returned nil")
		}
		valid, ok := stats["valid"].(bool)
		if !ok || !valid {
			t.Error("GetDecompositionStats should return valid=true for successful decomposition")
		}
	}
}

// TestDecomposeExpressionEdgeCases tests edge cases for expression decomposition
func TestDecomposeExpressionEdgeCases(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()
	tests := []struct {
		name    string
		expr    interface{}
		wantErr bool
	}{
		{
			name:    "nil expression",
			expr:    nil,
			wantErr: true,
		},
		{
			name:    "empty map",
			expr:    map[string]interface{}{},
			wantErr: true,
		},
		{
			name: "missing operator",
			expr: map[string]interface{}{
				"type":  "binaryOp",
				"left":  map[string]interface{}{"type": "number", "value": 1},
				"right": map[string]interface{}{"type": "number", "value": 2},
			},
			wantErr: true,
		},
		{
			name: "all arithmetic operators",
			expr: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "+",
				"left": map[string]interface{}{
					"type":     "binaryOp",
					"operator": "-",
					"left": map[string]interface{}{
						"type":     "binaryOp",
						"operator": "*",
						"left": map[string]interface{}{
							"type":     "binaryOp",
							"operator": "/",
							"left":     map[string]interface{}{"type": "number", "value": 8},
							"right":    map[string]interface{}{"type": "number", "value": 2},
						},
						"right": map[string]interface{}{"type": "number", "value": 3},
					},
					"right": map[string]interface{}{"type": "number", "value": 5},
				},
				"right": map[string]interface{}{"type": "number", "value": 1},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decomposer.DecomposeExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecomposeExpression() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestShouldDecomposeComplexity tests ShouldDecompose with various complexity levels
func TestShouldDecomposeComplexity(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()
	tests := []struct {
		name       string
		expr       interface{}
		wantResult bool
	}{
		{
			name: "simple literal - should not decompose",
			expr: map[string]interface{}{
				"type":  "number",
				"value": 42,
			},
			wantResult: false,
		},
		{
			name: "single operation - should not decompose",
			expr: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "+",
				"left":     map[string]interface{}{"type": "number", "value": 1},
				"right":    map[string]interface{}{"type": "number", "value": 2},
			},
			wantResult: false,
		},
		{
			name: "multiple operations - should decompose",
			expr: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "+",
				"left": map[string]interface{}{
					"type":     "binaryOp",
					"operator": "*",
					"left":     map[string]interface{}{"type": "number", "value": 2},
					"right":    map[string]interface{}{"type": "number", "value": 3},
				},
				"right": map[string]interface{}{
					"type":     "binaryOp",
					"operator": "-",
					"left":     map[string]interface{}{"type": "number", "value": 5},
					"right":    map[string]interface{}{"type": "number", "value": 1},
				},
			},
			wantResult: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := decomposer.ShouldDecompose(tt.expr)
			if result != tt.wantResult {
				t.Errorf("ShouldDecompose() = %v, want %v", result, tt.wantResult)
			}
		})
	}
}

// TestGetComplexity tests the GetComplexity function
func TestGetComplexity(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()
	tests := []struct {
		name          string
		expr          interface{}
		minComplexity int
	}{
		{
			name: "literal has low complexity",
			expr: map[string]interface{}{
				"type":  "number",
				"value": 42,
			},
			minComplexity: 0,
		},
		{
			name: "single operation has complexity >= 1",
			expr: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "+",
				"left":     map[string]interface{}{"type": "number", "value": 1},
				"right":    map[string]interface{}{"type": "number", "value": 2},
			},
			minComplexity: 1,
		},
		{
			name: "nested operations have higher complexity",
			expr: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "+",
				"left": map[string]interface{}{
					"type":     "binaryOp",
					"operator": "*",
					"left":     map[string]interface{}{"type": "number", "value": 2},
					"right":    map[string]interface{}{"type": "number", "value": 3},
				},
				"right": map[string]interface{}{"type": "number", "value": 4},
			},
			minComplexity: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			complexity := decomposer.GetComplexity(tt.expr)
			if complexity < tt.minComplexity {
				t.Errorf("GetComplexity() = %d, want at least %d", complexity, tt.minComplexity)
			}
		})
	}
}

// TestValidateDecompositionNilCases tests ValidateDecomposition with nil inputs
func TestValidateDecompositionNilCases(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()
	tests := []struct {
		name       string
		original   interface{}
		decomposed *DecomposedExpression
		wantErr    bool
	}{
		{
			name:       "nil decomposed expression",
			original:   map[string]interface{}{"type": "number", "value": 1},
			decomposed: nil,
			wantErr:    true,
		},
		{
			name:     "empty steps",
			original: map[string]interface{}{"type": "number", "value": 1},
			decomposed: &DecomposedExpression{
				Steps: []SimpleCondition{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := decomposer.ValidateDecomposition(tt.original, tt.decomposed)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDecomposition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetDecompositionStatsNil tests GetDecompositionStats with nil input
func TestGetDecompositionStatsNil(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()
	stats := decomposer.GetDecompositionStats(nil)
	if stats == nil {
		t.Error("GetDecompositionStats should not return nil")
	}
	valid, ok := stats["valid"].(bool)
	if !ok {
		t.Error("GetDecompositionStats should return a 'valid' field")
	}
	if valid {
		t.Error("GetDecompositionStats should return valid=false for nil input")
	}
}
