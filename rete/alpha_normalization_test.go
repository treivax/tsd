package rete

import (
	"github.com/stretchr/testify/require"
	"github.com/treivax/tsd/constraint"
	"testing"
)

// TestExtractFromLogicalExpressionMap_ConstraintLogicalOperation tests the []constraint.LogicalOperation path
func TestExtractFromLogicalExpressionMap_ConstraintLogicalOperation(t *testing.T) {
	tests := []struct {
		name          string
		expr          map[string]interface{}
		expectError   bool
		errorContains string
		minConditions int
		expectOpType  string
	}{
		{
			name: "operations as []constraint.LogicalOperation",
			expr: map[string]interface{}{
				"type": "logicalExpr",
				"left": map[string]interface{}{
					"type":     "comparison",
					"left":     "a",
					"operator": ">",
					"right":    "b",
				},
				"operations": []constraint.LogicalOperation{
					{
						Op: "AND",
						Right: map[string]interface{}{
							"type":     "comparison",
							"left":     "c",
							"operator": "<",
							"right":    "d",
						},
					},
				},
			},
			minConditions: 2,
			expectOpType:  "AND",
		},
		{
			name: "operations as []constraint.LogicalOperation with multiple ops",
			expr: map[string]interface{}{
				"type": "logicalExpr",
				"left": map[string]interface{}{
					"type":     "comparison",
					"left":     "a",
					"operator": ">",
					"right":    "b",
				},
				"operations": []constraint.LogicalOperation{
					{
						Op: "OR",
						Right: map[string]interface{}{
							"type":     "comparison",
							"left":     "c",
							"operator": "<",
							"right":    "d",
						},
					},
					{
						Op: "OR",
						Right: map[string]interface{}{
							"type":     "comparison",
							"left":     "e",
							"operator": "==",
							"right":    "f",
						},
					},
				},
			},
			minConditions: 3,
			expectOpType:  "OR",
		},
		{
			name: "operations as []constraint.LogicalOperation - mixed operators",
			expr: map[string]interface{}{
				"type": "logicalExpr",
				"left": map[string]interface{}{
					"type":     "comparison",
					"left":     "a",
					"operator": ">",
					"right":    "b",
				},
				"operations": []constraint.LogicalOperation{
					{
						Op: "AND",
						Right: map[string]interface{}{
							"type":     "comparison",
							"left":     "c",
							"operator": "<",
							"right":    "d",
						},
					},
					{
						Op: "OR",
						Right: map[string]interface{}{
							"type":     "comparison",
							"left":     "e",
							"operator": ">=",
							"right":    "f",
						},
					},
				},
			},
			minConditions: 3,
			expectOpType:  "MIXED",
		},
		{
			name: "operations as invalid type (string)",
			expr: map[string]interface{}{
				"type": "logicalExpr",
				"left": map[string]interface{}{
					"type":     "comparison",
					"left":     "a",
					"operator": ">",
					"right":    "b",
				},
				"operations": "invalid_string",
			},
			expectError:   true,
			errorContains: "operations doit être un tableau",
		},
		{
			name: "operations as []interface{} with non-map element",
			expr: map[string]interface{}{
				"type": "logicalExpr",
				"left": map[string]interface{}{
					"type":     "comparison",
					"left":     "a",
					"operator": ">",
					"right":    "b",
				},
				"operations": []interface{}{
					"not_a_map",
				},
			},
			expectError:   true,
			errorContains: "operation doit être une map",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conditions, opType, err := extractFromLogicalExpressionMap(tt.expr)
			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.GreaterOrEqual(t, len(conditions), tt.minConditions, "not enough conditions")
				require.Equal(t, tt.expectOpType, opType, "operator type mismatch")
			}
		})
	}
}

// TestNormalizeORExpression_Coverage tests OR expression normalization
func TestNormalizeORExpression_Coverage(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		expectError   bool
		errorContains string
	}{
		{
			name:          "nil expression",
			input:         nil,
			expectError:   true,
			errorContains: "expression nil",
		},
		{
			name: "non-OR expression - should error",
			input: constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: map[string]interface{}{
					"type":     "comparison",
					"left":     "a",
					"operator": ">",
					"right":    "b",
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "AND",
						Right: map[string]interface{}{
							"type":     "comparison",
							"left":     "c",
							"operator": "<",
							"right":    "d",
						},
					},
				},
			},
			expectError:   true,
			errorContains: "expression n'est pas de type OR",
		},
		{
			name: "valid OR LogicalExpression",
			input: constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: map[string]interface{}{
					"type":     "comparison",
					"left":     "z",
					"operator": ">",
					"right":    "a",
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "OR",
						Right: map[string]interface{}{
							"type":     "comparison",
							"left":     "a",
							"operator": "<",
							"right":    "b",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "valid OR map expression",
			input: map[string]interface{}{
				"type": "logicalExpr",
				"left": map[string]interface{}{
					"type":     "comparison",
					"left":     "x",
					"operator": ">",
					"right":    "y",
				},
				"operations": []map[string]interface{}{
					{
						"op": "OR",
						"right": map[string]interface{}{
							"type":     "comparison",
							"left":     "a",
							"operator": "<",
							"right":    "b",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name:          "unsupported expression type",
			input:         "string_expression",
			expectError:   true,
			errorContains: "type d'expression non supporté",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizeORExpression(tt.input)
			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}

// TestCanonicalMap_SpecialTypes tests additional special type handling in canonicalMap
func TestCanonicalMap_SpecialTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		contains string // substring that should be in output
	}{
		{
			name: "logicalExpression type with operations",
			input: map[string]interface{}{
				"type": "logicalExpr",
				"left": "leftValue",
				"operations": []interface{}{
					map[string]interface{}{
						"op":    "AND",
						"right": "rightValue",
					},
				},
			},
			contains: "logical",
		},
		{
			name: "logicalExpression without operations",
			input: map[string]interface{}{
				"type": "logicalExpr",
				"left": "value",
			},
			contains: "logical",
		},
		{
			name: "binary_op alternative type name",
			input: map[string]interface{}{
				"type":     "binary_op",
				"left":     "a",
				"operator": "*",
				"right":    "b",
			},
			contains: "binaryOp",
		},
		{
			name: "comparison type",
			input: map[string]interface{}{
				"type":     "comparison",
				"left":     "x",
				"operator": "<=",
				"right":    "y",
			},
			contains: "binaryOp",
		},
		{
			name: "binary operation with op field instead of operator",
			input: map[string]interface{}{
				"type":  "binaryOperation",
				"left":  "m",
				"op":    "/",
				"right": "n",
			},
			contains: "binaryOp",
		},
		{
			name: "stringLiteral type",
			input: map[string]interface{}{
				"type":  "stringLiteral",
				"value": "hello world",
			},
			contains: "literal",
		},
		{
			name: "booleanLiteral type",
			input: map[string]interface{}{
				"type":  "booleanLiteral",
				"value": false,
			},
			contains: "literal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := canonicalMap(tt.input)
			require.Contains(t, result, tt.contains)
		})
	}
}

// TestCanonicalValue_AdditionalTypes tests additional type handling in canonicalValue
func TestCanonicalValue_AdditionalTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		contains string
	}{
		{
			name: "constraint.FieldAccess",
			input: constraint.FieldAccess{
				Type:   "fieldAccess",
				Object: "obj",
				Field:  "field",
			},
			contains: "fieldAccess(obj,field)",
		},
		{
			name: "constraint.NumberLiteral",
			input: constraint.NumberLiteral{
				Type:  "numberLiteral",
				Value: 123.45,
			},
			contains: "literal(123.45)",
		},
		{
			name: "constraint.StringLiteral",
			input: constraint.StringLiteral{
				Type:  "stringLiteral",
				Value: "test",
			},
			contains: "literal(test)",
		},
		{
			name: "constraint.BooleanLiteral true",
			input: constraint.BooleanLiteral{
				Type:  "booleanLiteral",
				Value: true,
			},
			contains: "literal(true)",
		},
		{
			name: "constraint.BooleanLiteral false",
			input: constraint.BooleanLiteral{
				Type:  "booleanLiteral",
				Value: false,
			},
			contains: "literal(false)",
		},
		{
			name: "constraint.BinaryOperation",
			input: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Left:     "a",
				Operator: "-",
				Right:    "b",
			},
			contains: "binaryOp",
		},
		{
			name: "constraint.LogicalExpression",
			input: constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: "left",
				Operations: []constraint.LogicalOperation{
					{Op: "AND", Right: "right"},
				},
			},
			contains: "logical",
		},
		{
			name:     "int32",
			input:    int32(42),
			contains: "int(42)",
		},
		{
			name:     "float32",
			input:    float32(3.14),
			contains: "float(3.14)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := canonicalValue(tt.input)
			require.Contains(t, result, tt.contains)
		})
	}
}

// TestAnalyzeExpression_Coverage tests expression type analysis
func TestAnalyzeExpression_Coverage(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		expectType    ExpressionType
		expectError   bool
		errorContains string
	}{
		{
			name:          "nil expression",
			input:         nil,
			expectError:   true,
			errorContains: "expression nil",
		},
		{
			name: "AND LogicalExpression",
			input: constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: map[string]interface{}{
					"type":     "comparison",
					"left":     "a",
					"operator": ">",
					"right":    "b",
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "AND",
						Right: map[string]interface{}{
							"type":     "comparison",
							"left":     "c",
							"operator": "<",
							"right":    "d",
						},
					},
				},
			},
			expectType: ExprTypeAND,
		},
		{
			name: "OR LogicalExpression",
			input: constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: map[string]interface{}{
					"type":     "comparison",
					"left":     "x",
					"operator": ">",
					"right":    "y",
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "OR",
						Right: map[string]interface{}{
							"type":     "comparison",
							"left":     "m",
							"operator": "==",
							"right":    "n",
						},
					},
				},
			},
			expectType: ExprTypeOR,
		},
		{
			name: "Mixed LogicalExpression",
			input: constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: map[string]interface{}{
					"type":     "comparison",
					"left":     "a",
					"operator": ">",
					"right":    "b",
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "AND",
						Right: map[string]interface{}{
							"type":     "comparison",
							"left":     "c",
							"operator": "<",
							"right":    "d",
						},
					},
					{
						Op: "OR",
						Right: map[string]interface{}{
							"type":     "comparison",
							"left":     "e",
							"operator": ">=",
							"right":    "f",
						},
					},
				},
			},
			expectType: ExprTypeMixed,
		},
		{
			name: "Simple comparison",
			input: map[string]interface{}{
				"type":     "comparison",
				"left":     "a",
				"operator": ">",
				"right":    "b",
			},
			expectType: ExprTypeSimple,
		},
		{
			name: "AND map expression",
			input: map[string]interface{}{
				"type": "logicalExpr",
				"left": map[string]interface{}{
					"type":     "comparison",
					"left":     "p",
					"operator": ">",
					"right":    "q",
				},
				"operations": []map[string]interface{}{
					{
						"op": "AND",
						"right": map[string]interface{}{
							"type":     "comparison",
							"left":     "r",
							"operator": "<",
							"right":    "s",
						},
					},
				},
			},
			expectType: ExprTypeAND,
		},
		{
			name: "OR map expression",
			input: map[string]interface{}{
				"type": "logicalExpr",
				"left": map[string]interface{}{
					"type":     "comparison",
					"left":     "x",
					"operator": "==",
					"right":    "y",
				},
				"operations": []interface{}{
					map[string]interface{}{
						"op": "OR",
						"right": map[string]interface{}{
							"type":     "comparison",
							"left":     "z",
							"operator": "!=",
							"right":    "w",
						},
					},
				},
			},
			expectType: ExprTypeOR,
		},
		{
			name:          "unsupported type",
			input:         "string",
			expectError:   true,
			errorContains: "type d'expression non supporté",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exprType, err := AnalyzeExpression(tt.input)
			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectType, exprType)
			}
		})
	}
}
