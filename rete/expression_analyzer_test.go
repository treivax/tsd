// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"github.com/treivax/tsd/constraint"
	"strings"
	"testing"
)

// TestAnalyzeExpression_Simple teste l'analyse d'expressions simples
func TestAnalyzeExpression_Simple(t *testing.T) {
	tests := []struct {
		name     string
		expr     interface{}
		expected ExpressionType
		wantErr  bool
	}{
		{
			name: "simple binary operation - comparison",
			expr: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
				Operator: ">",
				Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
			},
			expected: ExprTypeSimple,
			wantErr:  false,
		},
		{
			name: "simple binary operation - equality",
			expr: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
				Operator: "==",
				Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "active"},
			},
			expected: ExprTypeSimple,
			wantErr:  false,
		},
		{
			name: "simple condition map format",
			expr: map[string]interface{}{
				"type":     "binaryOperation",
				"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
				"operator": ">=",
				"right":    map[string]interface{}{"type": "numberLiteral", "value": 21},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AnalyzeExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("AnalyzeExpression() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestApplyDeMorganTransformation teste la transformation de De Morgan
func TestApplyDeMorganTransformation(t *testing.T) {
	tests := []struct {
		name              string
		expr              interface{}
		expectTransformed bool
		validateResult    func(t *testing.T, result interface{})
	}{
		{
			name: "NOT(A OR B) -> (NOT A) AND (NOT B)",
			expr: constraint.NotConstraint{
				Expression: constraint.LogicalExpression{
					Left: constraint.BinaryOperation{
						Left:     constraint.FieldAccess{Object: "p", Field: "age"},
						Operator: ">",
						Right:    constraint.NumberLiteral{Value: 18},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Left:     constraint.FieldAccess{Object: "p", Field: "salary"},
								Operator: "<",
								Right:    constraint.NumberLiteral{Value: 50000},
							},
						},
					},
				},
			},
			expectTransformed: true,
			validateResult: func(t *testing.T, result interface{}) {
				logicalExpr, ok := result.(constraint.LogicalExpression)
				if !ok {
					t.Errorf("Expected LogicalExpression, got %T", result)
					return
				}
				// Vérifier que le left est enveloppé dans NOT
				if _, ok := logicalExpr.Left.(constraint.NotConstraint); !ok {
					t.Errorf("Expected left to be NotConstraint, got %T", logicalExpr.Left)
				}
				// Vérifier que l'opérateur est AND
				if len(logicalExpr.Operations) == 0 {
					t.Error("Expected at least one operation")
					return
				}
				op := logicalExpr.Operations[0].Op
				if op != "AND" && op != "&&" {
					t.Errorf("Expected AND operator, got %s", op)
				}
				// Vérifier que right est enveloppé dans NOT
				if _, ok := logicalExpr.Operations[0].Right.(constraint.NotConstraint); !ok {
					t.Errorf("Expected right to be NotConstraint, got %T", logicalExpr.Operations[0].Right)
				}
			},
		},
		{
			name: "NOT(A AND B) -> (NOT A) OR (NOT B)",
			expr: constraint.NotConstraint{
				Expression: constraint.LogicalExpression{
					Left: constraint.BinaryOperation{
						Left:     constraint.FieldAccess{Object: "p", Field: "age"},
						Operator: ">",
						Right:    constraint.NumberLiteral{Value: 18},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "AND",
							Right: constraint.BinaryOperation{
								Left:     constraint.FieldAccess{Object: "p", Field: "salary"},
								Operator: ">=",
								Right:    constraint.NumberLiteral{Value: 50000},
							},
						},
					},
				},
			},
			expectTransformed: true,
			validateResult: func(t *testing.T, result interface{}) {
				logicalExpr, ok := result.(constraint.LogicalExpression)
				if !ok {
					t.Errorf("Expected LogicalExpression, got %T", result)
					return
				}
				// Vérifier que l'opérateur est OR
				if len(logicalExpr.Operations) == 0 {
					t.Error("Expected at least one operation")
					return
				}
				op := logicalExpr.Operations[0].Op
				if op != "OR" && op != "||" {
					t.Errorf("Expected OR operator, got %s", op)
				}
			},
		},
		{
			name: "NOT(simple) should not transform",
			expr: constraint.NotConstraint{
				Expression: constraint.BinaryOperation{
					Left:     constraint.FieldAccess{Object: "p", Field: "active"},
					Operator: "==",
					Right:    constraint.BooleanLiteral{Value: true},
				},
			},
			expectTransformed: false,
			validateResult: func(t *testing.T, result interface{}) {
				_, ok := result.(constraint.NotConstraint)
				if !ok {
					t.Errorf("Expected NotConstraint to remain unchanged, got %T", result)
				}
			},
		},
		{
			name: "NOT(A OR B OR C) -> (NOT A) AND (NOT B) AND (NOT C)",
			expr: constraint.NotConstraint{
				Expression: constraint.LogicalExpression{
					Left: constraint.BinaryOperation{
						Left:     constraint.FieldAccess{Object: "p", Field: "status"},
						Operator: "==",
						Right:    constraint.StringLiteral{Value: "active"},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Left:     constraint.FieldAccess{Object: "p", Field: "status"},
								Operator: "==",
								Right:    constraint.StringLiteral{Value: "pending"},
							},
						},
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Left:     constraint.FieldAccess{Object: "p", Field: "status"},
								Operator: "==",
								Right:    constraint.StringLiteral{Value: "review"},
							},
						},
					},
				},
			},
			expectTransformed: true,
			validateResult: func(t *testing.T, result interface{}) {
				logicalExpr, ok := result.(constraint.LogicalExpression)
				if !ok {
					t.Errorf("Expected LogicalExpression, got %T", result)
					return
				}
				// Vérifier qu'il y a 2 opérations
				if len(logicalExpr.Operations) != 2 {
					t.Errorf("Expected 2 operations, got %d", len(logicalExpr.Operations))
				}
			},
		},
		{
			name: "NOT(A OR B) map format",
			expr: map[string]interface{}{
				"type": "notConstraint",
				"expression": map[string]interface{}{
					"type": "logicalExpression",
					"left": map[string]interface{}{
						"type":     "binaryOperation",
						"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
						"operator": ">",
						"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
					},
					"operations": []interface{}{
						map[string]interface{}{
							"op": "OR",
							"right": map[string]interface{}{
								"type":     "binaryOperation",
								"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "salary"},
								"operator": "<",
								"right":    map[string]interface{}{"type": "numberLiteral", "value": 50000},
							},
						},
					},
				},
			},
			expectTransformed: true,
			validateResult: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Errorf("Expected map[string]interface{}, got %T", result)
					return
				}
				if resultMap["type"] != "logicalExpression" {
					t.Errorf("Expected type logicalExpression, got %v", resultMap["type"])
				}
			},
		},
		{
			name: "Simple expression should not transform",
			expr: constraint.BinaryOperation{
				Left:     constraint.FieldAccess{Object: "p", Field: "age"},
				Operator: ">",
				Right:    constraint.NumberLiteral{Value: 18},
			},
			expectTransformed: false,
			validateResult: func(t *testing.T, result interface{}) {
				_, ok := result.(constraint.BinaryOperation)
				if !ok {
					t.Errorf("Expected BinaryOperation to remain unchanged, got %T", result)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, transformed := ApplyDeMorganTransformation(tt.expr)
			if transformed != tt.expectTransformed {
				t.Errorf("ApplyDeMorganTransformation() transformed = %v, want %v", transformed, tt.expectTransformed)
			}
			if tt.validateResult != nil {
				tt.validateResult(t, result)
			}
		})
	}
}

// TestShouldApplyDeMorgan teste la logique de décision pour appliquer De Morgan
func TestShouldApplyDeMorgan(t *testing.T) {
	tests := []struct {
		name     string
		expr     interface{}
		expected bool
	}{
		{
			name: "NOT(A OR B) should apply",
			expr: constraint.NotConstraint{
				Expression: constraint.LogicalExpression{
					Left: constraint.BinaryOperation{
						Left:     constraint.FieldAccess{Object: "p", Field: "age"},
						Operator: ">",
						Right:    constraint.NumberLiteral{Value: 18},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Left:     constraint.FieldAccess{Object: "p", Field: "salary"},
								Operator: "<",
								Right:    constraint.NumberLiteral{Value: 50000},
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "NOT(A AND B) simple should apply",
			expr: constraint.NotConstraint{
				Expression: constraint.LogicalExpression{
					Left: constraint.BinaryOperation{
						Left:     constraint.FieldAccess{Object: "p", Field: "age"},
						Operator: ">",
						Right:    constraint.NumberLiteral{Value: 18},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "AND",
							Right: constraint.BinaryOperation{
								Left:     constraint.FieldAccess{Object: "p", Field: "salary"},
								Operator: ">=",
								Right:    constraint.NumberLiteral{Value: 50000},
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "NOT(simple) should not apply",
			expr: constraint.NotConstraint{
				Expression: constraint.BinaryOperation{
					Left:     constraint.FieldAccess{Object: "p", Field: "active"},
					Operator: "==",
					Right:    constraint.BooleanLiteral{Value: true},
				},
			},
			expected: false,
		},
		{
			name: "Simple expression should not apply",
			expr: constraint.BinaryOperation{
				Left:     constraint.FieldAccess{Object: "p", Field: "age"},
				Operator: ">",
				Right:    constraint.NumberLiteral{Value: 18},
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShouldApplyDeMorgan(tt.expr)
			if result != tt.expected {
				t.Errorf("ShouldApplyDeMorgan() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestOptimizationHints teste la génération de hints d'optimisation
func TestOptimizationHints(t *testing.T) {
	tests := []struct {
		name          string
		expr          interface{}
		expectedHints []string
	}{
		{
			name: "NOT(A OR B) should suggest apply_demorgan_not_or",
			expr: constraint.NotConstraint{
				Expression: constraint.LogicalExpression{
					Left: constraint.BinaryOperation{
						Left:     constraint.FieldAccess{Object: "p", Field: "age"},
						Operator: ">",
						Right:    constraint.NumberLiteral{Value: 18},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Left:     constraint.FieldAccess{Object: "p", Field: "salary"},
								Operator: "<",
								Right:    constraint.NumberLiteral{Value: 50000},
							},
						},
					},
				},
			},
			expectedHints: []string{"apply_demorgan_not_or"},
		},
		{
			name: "NOT(A AND B) should suggest apply_demorgan_not_and",
			expr: constraint.NotConstraint{
				Expression: constraint.LogicalExpression{
					Left: constraint.BinaryOperation{
						Left:     constraint.FieldAccess{Object: "p", Field: "age"},
						Operator: ">",
						Right:    constraint.NumberLiteral{Value: 18},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "AND",
							Right: constraint.BinaryOperation{
								Left:     constraint.FieldAccess{Object: "p", Field: "salary"},
								Operator: ">=",
								Right:    constraint.NumberLiteral{Value: 50000},
							},
						},
					},
				},
			},
			expectedHints: []string{"apply_demorgan_not_and"},
		},
		{
			name: "Mixed expression should suggest normalize_to_dnf",
			expr: constraint.LogicalExpression{
				Left: constraint.BinaryOperation{
					Left:     constraint.FieldAccess{Object: "p", Field: "age"},
					Operator: ">",
					Right:    constraint.NumberLiteral{Value: 18},
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "AND",
						Right: constraint.BinaryOperation{
							Left:     constraint.FieldAccess{Object: "p", Field: "salary"},
							Operator: ">=",
							Right:    constraint.NumberLiteral{Value: 50000},
						},
					},
					{
						Op: "OR",
						Right: constraint.BinaryOperation{
							Left:     constraint.FieldAccess{Object: "p", Field: "vip"},
							Operator: "==",
							Right:    constraint.BooleanLiteral{Value: true},
						},
					},
				},
			},
			expectedHints: []string{"normalize_to_dnf", "requires_beta_node"},
		},
		{
			name: "OR expression should suggest consider_dnf_expansion",
			expr: constraint.LogicalExpression{
				Left: constraint.BinaryOperation{
					Left:     constraint.FieldAccess{Object: "p", Field: "status"},
					Operator: "==",
					Right:    constraint.StringLiteral{Value: "active"},
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "OR",
						Right: constraint.BinaryOperation{
							Left:     constraint.FieldAccess{Object: "p", Field: "status"},
							Operator: "==",
							Right:    constraint.StringLiteral{Value: "pending"},
						},
					},
				},
			},
			expectedHints: []string{"consider_dnf_expansion", "requires_beta_node"},
		},
		{
			name: "Complex AND should suggest alpha_sharing and reordering",
			expr: constraint.LogicalExpression{
				Left: constraint.BinaryOperation{
					Left:     constraint.FieldAccess{Object: "p", Field: "age"},
					Operator: ">",
					Right:    constraint.NumberLiteral{Value: 18},
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "AND",
						Right: constraint.BinaryOperation{
							Left:     constraint.FieldAccess{Object: "p", Field: "salary"},
							Operator: ">=",
							Right:    constraint.NumberLiteral{Value: 50000},
						},
					},
					{
						Op: "AND",
						Right: constraint.BinaryOperation{
							Left:     constraint.FieldAccess{Object: "p", Field: "active"},
							Operator: "==",
							Right:    constraint.BooleanLiteral{Value: true},
						},
					},
				},
			},
			expectedHints: []string{"alpha_sharing_opportunity", "consider_reordering"},
		},
		{
			name: "Arithmetic should suggest simplification",
			expr: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Left:     constraint.FieldAccess{Object: "p", Field: "price"},
				Operator: "*",
				Right:    constraint.NumberLiteral{Value: 1.2},
			},
			expectedHints: []string{"consider_arithmetic_simplification"},
		},
		{
			name: "Simple expression should have minimal hints",
			expr: constraint.BinaryOperation{
				Left:     constraint.FieldAccess{Object: "p", Field: "age"},
				Operator: ">",
				Right:    constraint.NumberLiteral{Value: 18},
			},
			expectedHints: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := GetExpressionInfo(tt.expr)
			if err != nil {
				t.Fatalf("GetExpressionInfo() error = %v", err)
			}
			// Vérifier que tous les hints attendus sont présents
			for _, expectedHint := range tt.expectedHints {
				found := false
				for _, hint := range info.OptimizationHints {
					if hint == expectedHint {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected hint '%s' not found in %v", expectedHint, info.OptimizationHints)
				}
			}
		})
	}
}

// TestGetExpressionInfo_WithOptimizationHints teste GetExpressionInfo avec hints
func TestGetExpressionInfo_WithOptimizationHints(t *testing.T) {
	tests := []struct {
		name               string
		expr               interface{}
		checkHints         func(t *testing.T, hints []string)
		expectedComplexity int
	}{
		{
			name: "NOT(A OR B) with hints",
			expr: constraint.NotConstraint{
				Expression: constraint.LogicalExpression{
					Left: constraint.BinaryOperation{
						Left:     constraint.FieldAccess{Object: "p", Field: "age"},
						Operator: ">",
						Right:    constraint.NumberLiteral{Value: 18},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Left:     constraint.FieldAccess{Object: "p", Field: "salary"},
								Operator: "<",
								Right:    constraint.NumberLiteral{Value: 50000},
							},
						},
					},
				},
			},
			checkHints: func(t *testing.T, hints []string) {
				if len(hints) == 0 {
					t.Error("Expected at least one hint")
				}
				hasDemorgan := false
				for _, hint := range hints {
					if strings.Contains(hint, "demorgan") {
						hasDemorgan = true
						break
					}
				}
				if !hasDemorgan {
					t.Error("Expected demorgan hint")
				}
			},
			expectedComplexity: 4, // 2 (NOT) + 2 (OR with 2 terms)
		},
		{
			name: "Complex mixed expression",
			expr: constraint.LogicalExpression{
				Left: constraint.BinaryOperation{
					Left:     constraint.FieldAccess{Object: "p", Field: "age"},
					Operator: ">",
					Right:    constraint.NumberLiteral{Value: 18},
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "AND",
						Right: constraint.BinaryOperation{
							Left:     constraint.FieldAccess{Object: "p", Field: "salary"},
							Operator: ">=",
							Right:    constraint.NumberLiteral{Value: 50000},
						},
					},
					{
						Op: "OR",
						Right: constraint.BinaryOperation{
							Left:     constraint.FieldAccess{Object: "p", Field: "vip"},
							Operator: "==",
							Right:    constraint.BooleanLiteral{Value: true},
						},
					},
				},
			},
			checkHints: func(t *testing.T, hints []string) {
				if len(hints) == 0 {
					t.Error("Expected at least one hint")
				}
				hasNormalize := false
				for _, hint := range hints {
					if strings.Contains(hint, "normalize") {
						hasNormalize = true
					}
				}
				if !hasNormalize {
					t.Error("Expected normalize hint for mixed expression")
				}
				// Note: complexity hint only appears for complexity >= 4
				// This mixed expression has complexity 3 (1 left + 2 ops)
			},
			expectedComplexity: 3, // 1 left + 2 operations
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := GetExpressionInfo(tt.expr)
			if err != nil {
				t.Fatalf("GetExpressionInfo() error = %v", err)
			}
			if info.Complexity != tt.expectedComplexity {
				t.Errorf("Complexity = %d, want %d", info.Complexity, tt.expectedComplexity)
			}
			if info.OptimizationHints == nil {
				t.Error("OptimizationHints should not be nil")
				return
			}
			if tt.checkHints != nil {
				tt.checkHints(t, info.OptimizationHints)
			}
		})
	}
}

// TestDeMorganTransformationRoundtrip teste que les transformations sont correctes
func TestDeMorganTransformationRoundtrip(t *testing.T) {
	tests := []struct {
		name string
		expr interface{}
	}{
		{
			name: "NOT(A OR B) transformation preserves structure",
			expr: constraint.NotConstraint{
				Expression: constraint.LogicalExpression{
					Left: constraint.BinaryOperation{
						Left:     constraint.FieldAccess{Object: "p", Field: "age"},
						Operator: ">",
						Right:    constraint.NumberLiteral{Value: 18},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Left:     constraint.FieldAccess{Object: "p", Field: "salary"},
								Operator: "<",
								Right:    constraint.NumberLiteral{Value: 50000},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Appliquer la transformation
			transformed, applied := ApplyDeMorganTransformation(tt.expr)
			if !applied {
				t.Error("Expected transformation to be applied")
				return
			}
			// Analyser l'expression transformée
			resultType, err := AnalyzeExpression(transformed)
			if err != nil {
				t.Errorf("Failed to analyze transformed expression: %v", err)
				return
			}
			// L'expression transformée devrait être de type AND (pour NOT(OR))
			if resultType != ExprTypeAND {
				t.Errorf("Expected transformed expression to be AND, got %v", resultType)
			}
		})
	}
}

// TestOptimizationHintsIntegration teste l'intégration complète
func TestOptimizationHintsIntegration(t *testing.T) {
	// Expression complexe: NOT((A OR B) AND C)
	complexExpr := constraint.NotConstraint{
		Expression: constraint.LogicalExpression{
			Left: constraint.LogicalExpression{
				Left: constraint.BinaryOperation{
					Left:     constraint.FieldAccess{Object: "p", Field: "age"},
					Operator: ">",
					Right:    constraint.NumberLiteral{Value: 18},
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "OR",
						Right: constraint.BinaryOperation{
							Left:     constraint.FieldAccess{Object: "p", Field: "age"},
							Operator: "<",
							Right:    constraint.NumberLiteral{Value: 65},
						},
					},
				},
			},
			Operations: []constraint.LogicalOperation{
				{
					Op: "AND",
					Right: constraint.BinaryOperation{
						Left:     constraint.FieldAccess{Object: "p", Field: "active"},
						Operator: "==",
						Right:    constraint.BooleanLiteral{Value: true},
					},
				},
			},
		},
	}
	info, err := GetExpressionInfo(complexExpr)
	if err != nil {
		t.Fatalf("GetExpressionInfo() error = %v", err)
	}
	// Vérifier que des hints sont générés
	if len(info.OptimizationHints) == 0 {
		t.Error("Expected optimization hints for complex expression")
	}
	// Vérifier que InnerInfo existe
	if info.InnerInfo == nil {
		t.Error("Expected InnerInfo for NOT expression")
	}
	// Vérifier la complexité
	if info.Complexity < 4 {
		t.Errorf("Expected high complexity for nested expression, got %d", info.Complexity)
	}
	t.Logf("Generated hints: %v", info.OptimizationHints)
	t.Logf("Complexity: %d", info.Complexity)
	if info.InnerInfo != nil {
		t.Logf("Inner type: %v", info.InnerInfo.Type)
		t.Logf("Inner complexity: %d", info.InnerInfo.Complexity)
	}
}

// TestAnalyzeExpression_AND teste l'analyse d'expressions AND
func TestAnalyzeExpression_AND(t *testing.T) {
	tests := []struct {
		name     string
		expr     interface{}
		expected ExpressionType
		wantErr  bool
	}{
		{
			name: "simple AND expression",
			expr: constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
					Operator: ">",
					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "AND",
						Right: constraint.BinaryOperation{
							Type:     "binaryOperation",
							Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
							Operator: ">=",
							Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
						},
					},
				},
			},
			expected: ExprTypeAND,
			wantErr:  false,
		},
		{
			name: "multiple AND operations",
			expr: constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
					Operator: ">",
					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "AND",
						Right: constraint.BinaryOperation{
							Type:     "binaryOperation",
							Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
							Operator: ">=",
							Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
						},
					},
					{
						Op: "AND",
						Right: constraint.BinaryOperation{
							Type:     "binaryOperation",
							Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "active"},
							Operator: "==",
							Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
						},
					},
				},
			},
			expected: ExprTypeAND,
			wantErr:  false,
		},
		{
			name: "AND expression with && operator",
			expr: constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "x"},
					Operator: ">",
					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 0},
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "&&",
						Right: constraint.BinaryOperation{
							Type:     "binaryOperation",
							Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "y"},
							Operator: ">",
							Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 0},
						},
					},
				},
			},
			expected: ExprTypeAND,
			wantErr:  false,
		},
		{
			name: "AND expression map format",
			expr: map[string]interface{}{
				"type": "logicalExpression",
				"left": map[string]interface{}{
					"type":     "binaryOperation",
					"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
					"operator": ">",
					"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
				},
				"operations": []interface{}{
					map[string]interface{}{
						"op": "AND",
						"right": map[string]interface{}{
							"type":     "binaryOperation",
							"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "salary"},
							"operator": ">=",
							"right":    map[string]interface{}{"type": "numberLiteral", "value": 50000},
						},
					},
				},
			},
			expected: ExprTypeAND,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AnalyzeExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("AnalyzeExpression() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestAnalyzeExpression_OR teste l'analyse d'expressions OR
func TestAnalyzeExpression_OR(t *testing.T) {
	tests := []struct {
		name     string
		expr     interface{}
		expected ExpressionType
		wantErr  bool
	}{
		{
			name: "simple OR expression",
			expr: constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
					Operator: "==",
					Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "active"},
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "OR",
						Right: constraint.BinaryOperation{
							Type:     "binaryOperation",
							Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
							Operator: "==",
							Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "pending"},
						},
					},
				},
			},
			expected: ExprTypeOR,
			wantErr:  false,
		},
		{
			name: "multiple OR operations",
			expr: constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "type"},
					Operator: "==",
					Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "A"},
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "OR",
						Right: constraint.BinaryOperation{
							Type:     "binaryOperation",
							Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "type"},
							Operator: "==",
							Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "B"},
						},
					},
					{
						Op: "OR",
						Right: constraint.BinaryOperation{
							Type:     "binaryOperation",
							Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "type"},
							Operator: "==",
							Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "C"},
						},
					},
				},
			},
			expected: ExprTypeOR,
			wantErr:  false,
		},
		{
			name: "OR expression with || operator",
			expr: constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "vip"},
					Operator: "==",
					Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "||",
						Right: constraint.BinaryOperation{
							Type:     "binaryOperation",
							Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "premium"},
							Operator: "==",
							Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
						},
					},
				},
			},
			expected: ExprTypeOR,
			wantErr:  false,
		},
		{
			name: "OR expression map format",
			expr: map[string]interface{}{
				"type": "logicalExpression",
				"left": map[string]interface{}{
					"type":     "binaryOperation",
					"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "status"},
					"operator": "==",
					"right":    map[string]interface{}{"type": "stringLiteral", "value": "active"},
				},
				"operations": []interface{}{
					map[string]interface{}{
						"op": "OR",
						"right": map[string]interface{}{
							"type":     "binaryOperation",
							"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "status"},
							"operator": "==",
							"right":    map[string]interface{}{"type": "stringLiteral", "value": "pending"},
						},
					},
				},
			},
			expected: ExprTypeOR,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AnalyzeExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("AnalyzeExpression() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestAnalyzeExpression_Mixed_AND_OR teste l'analyse d'expressions mixtes
func TestAnalyzeExpression_Mixed_AND_OR(t *testing.T) {
	tests := []struct {
		name     string
		expr     interface{}
		expected ExpressionType
		wantErr  bool
	}{
		{
			name: "mixed AND and OR - (A AND B) OR C",
			expr: constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
					Operator: ">",
					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "AND",
						Right: constraint.BinaryOperation{
							Type:     "binaryOperation",
							Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
							Operator: ">=",
							Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
						},
					},
					{
						Op: "OR",
						Right: constraint.BinaryOperation{
							Type:     "binaryOperation",
							Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "vip"},
							Operator: "==",
							Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
						},
					},
				},
			},
			expected: ExprTypeMixed,
			wantErr:  false,
		},
		{
			name: "mixed OR and AND - A OR (B AND C)",
			expr: constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "admin"},
					Operator: "==",
					Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "OR",
						Right: constraint.BinaryOperation{
							Type:     "binaryOperation",
							Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
							Operator: ">",
							Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 21},
						},
					},
					{
						Op: "AND",
						Right: constraint.BinaryOperation{
							Type:     "binaryOperation",
							Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "verified"},
							Operator: "==",
							Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
						},
					},
				},
			},
			expected: ExprTypeMixed,
			wantErr:  false,
		},
		{
			name: "complex mixed expression",
			expr: constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "x"},
					Operator: ">",
					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 0},
				},
				Operations: []constraint.LogicalOperation{
					{
						Op: "AND",
						Right: constraint.BinaryOperation{
							Type:     "binaryOperation",
							Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "y"},
							Operator: ">",
							Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 0},
						},
					},
					{
						Op: "OR",
						Right: constraint.BinaryOperation{
							Type:     "binaryOperation",
							Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "z"},
							Operator: "<",
							Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 100},
						},
					},
					{
						Op: "AND",
						Right: constraint.BinaryOperation{
							Type:     "binaryOperation",
							Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "w"},
							Operator: "!=",
							Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 0},
						},
					},
				},
			},
			expected: ExprTypeMixed,
			wantErr:  false,
		},
		{
			name: "mixed expression map format",
			expr: map[string]interface{}{
				"type": "logicalExpression",
				"left": map[string]interface{}{
					"type":     "binaryOperation",
					"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
					"operator": ">",
					"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
				},
				"operations": []interface{}{
					map[string]interface{}{
						"op": "AND",
						"right": map[string]interface{}{
							"type":     "binaryOperation",
							"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "salary"},
							"operator": ">=",
							"right":    map[string]interface{}{"type": "numberLiteral", "value": 50000},
						},
					},
					map[string]interface{}{
						"op": "OR",
						"right": map[string]interface{}{
							"type":     "binaryOperation",
							"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "vip"},
							"operator": "==",
							"right":    map[string]interface{}{"type": "booleanLiteral", "value": true},
						},
					},
				},
			},
			expected: ExprTypeMixed,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AnalyzeExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("AnalyzeExpression() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestAnalyzeExpression_Arithmetic teste l'analyse d'expressions arithmétiques
func TestAnalyzeExpression_Arithmetic(t *testing.T) {
	tests := []struct {
		name     string
		expr     interface{}
		expected ExpressionType
		wantErr  bool
	}{
		{
			name: "arithmetic addition",
			expr: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "price"},
				Operator: "+",
				Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 10},
			},
			expected: ExprTypeArithmetic,
			wantErr:  false,
		},
		{
			name: "arithmetic multiplication",
			expr: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "price"},
				Operator: "*",
				Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 1.2},
			},
			expected: ExprTypeArithmetic,
			wantErr:  false,
		},
		{
			name: "arithmetic division",
			expr: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "total"},
				Operator: "/",
				Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 2},
			},
			expected: ExprTypeArithmetic,
			wantErr:  false,
		},
		{
			name: "arithmetic modulo",
			expr: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "count"},
				Operator: "%",
				Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 5},
			},
			expected: ExprTypeArithmetic,
			wantErr:  false,
		},
		{
			name: "arithmetic map format",
			expr: map[string]interface{}{
				"type":     "binaryOperation",
				"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "price"},
				"operator": "*",
				"right":    map[string]interface{}{"type": "numberLiteral", "value": 1.5},
			},
			expected: ExprTypeArithmetic,
			wantErr:  false,
		},
		{
			name: "arithmetic operation type",
			expr: map[string]interface{}{
				"type":     "arithmeticOperation",
				"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "x"},
				"operator": "+",
				"right":    map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "y"},
			},
			expected: ExprTypeArithmetic,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AnalyzeExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("AnalyzeExpression() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCanDecompose_AllTypes teste la décomposabilité de tous les types
func TestCanDecompose_AllTypes(t *testing.T) {
	tests := []struct {
		name     string
		exprType ExpressionType
		expected bool
	}{
		{
			name:     "SimpleCondition can decompose",
			exprType: ExprTypeSimple,
			expected: true,
		},
		{
			name:     "ANDExpression can decompose",
			exprType: ExprTypeAND,
			expected: true,
		},
		{
			name:     "ArithmeticChain can decompose",
			exprType: ExprTypeArithmetic,
			expected: true,
		},
		{
			name:     "NOTExpression can decompose",
			exprType: ExprTypeNOT,
			expected: true,
		},
		{
			name:     "ORExpression cannot decompose",
			exprType: ExprTypeOR,
			expected: false,
		},
		{
			name:     "MixedExpression cannot decompose",
			exprType: ExprTypeMixed,
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanDecompose(tt.exprType)
			if result != tt.expected {
				t.Errorf("CanDecompose(%v) = %v, want %v", tt.exprType, result, tt.expected)
			}
		})
	}
}

// TestShouldNormalize_AllTypes teste la nécessité de normalisation pour tous les types
func TestShouldNormalize_AllTypes(t *testing.T) {
	tests := []struct {
		name     string
		exprType ExpressionType
		expected bool
	}{
		{
			name:     "SimpleCondition does not need normalization",
			exprType: ExprTypeSimple,
			expected: false,
		},
		{
			name:     "ANDExpression does not need normalization",
			exprType: ExprTypeAND,
			expected: false,
		},
		{
			name:     "ArithmeticChain does not need normalization",
			exprType: ExprTypeArithmetic,
			expected: false,
		},
		{
			name:     "NOTExpression does not need normalization",
			exprType: ExprTypeNOT,
			expected: false,
		},
		{
			name:     "ORExpression should be normalized",
			exprType: ExprTypeOR,
			expected: true,
		},
		{
			name:     "MixedExpression should be normalized",
			exprType: ExprTypeMixed,
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShouldNormalize(tt.exprType)
			if result != tt.expected {
				t.Errorf("ShouldNormalize(%v) = %v, want %v", tt.exprType, result, tt.expected)
			}
		})
	}
}

// TestExpressionType_String teste la représentation textuelle des types
func TestExpressionType_String(t *testing.T) {
	tests := []struct {
		exprType ExpressionType
		expected string
	}{
		{ExprTypeSimple, "ExprTypeSimple"},
		{ExprTypeAND, "ExprTypeAND"},
		{ExprTypeOR, "ExprTypeOR"},
		{ExprTypeMixed, "ExprTypeMixed"},
		{ExprTypeArithmetic, "ExprTypeArithmetic"},
		{ExprTypeNOT, "ExprTypeNOT"},
		{ExpressionType(999), "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.exprType.String()
			if result != tt.expected {
				t.Errorf("ExpressionType.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestGetExpressionComplexity teste le calcul de complexité
func TestGetExpressionComplexity(t *testing.T) {
	tests := []struct {
		name     string
		exprType ExpressionType
		expected int
	}{
		{
			name:     "SimpleCondition complexity",
			exprType: ExprTypeSimple,
			expected: 1,
		},
		{
			name:     "ANDExpression complexity",
			exprType: ExprTypeAND,
			expected: 2,
		},
		{
			name:     "ArithmeticChain complexity",
			exprType: ExprTypeArithmetic,
			expected: 2,
		},
		{
			name:     "NOTExpression complexity",
			exprType: ExprTypeNOT,
			expected: 2,
		},
		{
			name:     "ORExpression complexity",
			exprType: ExprTypeOR,
			expected: 3,
		},
		{
			name:     "MixedExpression complexity",
			exprType: ExprTypeMixed,
			expected: 4,
		},
		{
			name:     "Unknown type complexity",
			exprType: ExpressionType(999),
			expected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetExpressionComplexity(tt.exprType)
			if result != tt.expected {
				t.Errorf("GetExpressionComplexity(%v) = %v, want %v", tt.exprType, result, tt.expected)
			}
		})
	}
}

// TestRequiresBetaNode teste si un type nécessite des nœuds beta
func TestRequiresBetaNode(t *testing.T) {
	tests := []struct {
		name     string
		exprType ExpressionType
		expected bool
	}{
		{
			name:     "SimpleCondition does not require beta",
			exprType: ExprTypeSimple,
			expected: false,
		},
		{
			name:     "ANDExpression does not require beta",
			exprType: ExprTypeAND,
			expected: false,
		},
		{
			name:     "ArithmeticChain does not require beta",
			exprType: ExprTypeArithmetic,
			expected: false,
		},
		{
			name:     "NOTExpression does not require beta",
			exprType: ExprTypeNOT,
			expected: false,
		},
		{
			name:     "ORExpression requires beta",
			exprType: ExprTypeOR,
			expected: true,
		},
		{
			name:     "MixedExpression requires beta",
			exprType: ExprTypeMixed,
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RequiresBetaNode(tt.exprType)
			if result != tt.expected {
				t.Errorf("RequiresBetaNode(%v) = %v, want %v", tt.exprType, result, tt.expected)
			}
		})
	}
}

// TestGetExpressionInfo teste la récupération d'informations complètes
func TestGetExpressionInfo(t *testing.T) {
	// Test avec une condition simple
	expr := constraint.BinaryOperation{
		Type:     "binaryOperation",
		Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		Operator: ">",
		Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	}
	info, err := GetExpressionInfo(expr)
	if err != nil {
		t.Fatalf("GetExpressionInfo() error = %v", err)
	}
	if info.Type != ExprTypeSimple {
		t.Errorf("info.Type = %v, want %v", info.Type, ExprTypeSimple)
	}
	if !info.CanDecompose {
		t.Errorf("info.CanDecompose = false, want true")
	}
	if info.ShouldNormalize {
		t.Errorf("info.ShouldNormalize = true, want false")
	}
	if info.Complexity != 1 {
		t.Errorf("info.Complexity = %v, want 1", info.Complexity)
	}
	if info.RequiresBeta {
		t.Errorf("info.RequiresBeta = true, want false")
	}
	// Test avec une expression mixte
	mixedExpr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
			Operator: ">",
			Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
					Operator: ">=",
					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
				},
			},
			{
				Op: "OR",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "vip"},
					Operator: "==",
					Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
				},
			},
		},
	}
	info2, err := GetExpressionInfo(mixedExpr)
	if err != nil {
		t.Fatalf("GetExpressionInfo() error = %v", err)
	}
	if info2.Type != ExprTypeMixed {
		t.Errorf("info2.Type = %v, want %v", info2.Type, ExprTypeMixed)
	}
	if info2.CanDecompose {
		t.Errorf("info2.CanDecompose = true, want false")
	}
	if !info2.ShouldNormalize {
		t.Errorf("info2.ShouldNormalize = false, want true")
	}
	if info2.Complexity != 3 {
		t.Errorf("info2.Complexity = %v, want 3", info2.Complexity)
	}
	if !info2.RequiresBeta {
		t.Errorf("info2.RequiresBeta = false, want true")
	}
	// Test avec nil
	_, err = GetExpressionInfo(nil)
	if err == nil {
		t.Errorf("GetExpressionInfo(nil) should return error")
	}
}

// TestAnalyzeExpression_EdgeCases teste les cas limites
func TestAnalyzeExpression_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		expr     interface{}
		expected ExpressionType
		wantErr  bool
	}{
		{
			name: "logical expression with no operations",
			expr: constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
					Operator: ">",
					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
				},
				Operations: []constraint.LogicalOperation{},
			},
			expected: ExprTypeSimple,
			wantErr:  false,
		},
		{
			name: "map without type",
			expr: map[string]interface{}{
				"left":     "something",
				"operator": ">",
				"right":    18,
			},
			expected: ExprTypeSimple,
			wantErr:  true,
		},
		{
			name: "map with unsupported type",
			expr: map[string]interface{}{
				"type": "unknownType",
			},
			expected: ExprTypeSimple,
			wantErr:  true,
		},
		{
			name: "logical expression map with no operations",
			expr: map[string]interface{}{
				"type": "logicalExpression",
				"left": map[string]interface{}{
					"type":     "binaryOperation",
					"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
					"operator": ">",
					"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
				},
			},
			expected: ExprTypeSimple,
			wantErr:  false,
		},
		{
			name: "logical expression map with empty operations",
			expr: map[string]interface{}{
				"type": "logicalExpression",
				"left": map[string]interface{}{
					"type":     "binaryOperation",
					"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
					"operator": ">",
					"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
				},
				"operations": []interface{}{},
			},
			expected: ExprTypeSimple,
			wantErr:  false,
		},
		{
			name: "constraint with arithmetic operator",
			expr: constraint.Constraint{
				Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "x"},
				Operator: "+",
				Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 5},
			},
			expected: ExprTypeArithmetic,
			wantErr:  false,
		},
		{
			name:     "empty constraint",
			expr:     constraint.Constraint{},
			expected: ExprTypeSimple,
			wantErr:  false,
		},
		{
			name: "NOT constraint",
			expr: constraint.NotConstraint{
				Type: "notConstraint",
				Expression: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "active"},
					Operator: "==",
					Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
				},
			},
			expected: ExprTypeNOT,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AnalyzeExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("AnalyzeExpression() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestAnalyzeExpression_NOT teste l'analyse d'expressions NOT
func TestAnalyzeExpression_NOT(t *testing.T) {
	tests := []struct {
		name     string
		expr     interface{}
		expected ExpressionType
		wantErr  bool
	}{
		{
			name: "simple NOT constraint",
			expr: constraint.NotConstraint{
				Type: "notConstraint",
				Expression: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "active"},
					Operator: "==",
					Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
				},
			},
			expected: ExprTypeNOT,
			wantErr:  false,
		},
		{
			name: "NOT with complex expression",
			expr: constraint.NotConstraint{
				Type: "notConstraint",
				Expression: constraint.LogicalExpression{
					Type: "logicalExpr",
					Left: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
						Operator: ">",
						Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "AND",
							Right: constraint.BinaryOperation{
								Type:     "binaryOperation",
								Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
								Operator: "<",
								Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
							},
						},
					},
				},
			},
			expected: ExprTypeNOT,
			wantErr:  false,
		},
		{
			name: "NOT map format",
			expr: map[string]interface{}{
				"type": "notConstraint",
				"expression": map[string]interface{}{
					"type":     "binaryOperation",
					"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "verified"},
					"operator": "==",
					"right":    map[string]interface{}{"type": "booleanLiteral", "value": true},
				},
			},
			expected: ExprTypeNOT,
			wantErr:  false,
		},
		{
			name: "NOT with 'not' type",
			expr: map[string]interface{}{
				"type": "not",
				"expression": map[string]interface{}{
					"type":     "binaryOperation",
					"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "status"},
					"operator": "==",
					"right":    map[string]interface{}{"type": "stringLiteral", "value": "inactive"},
				},
			},
			expected: ExprTypeNOT,
			wantErr:  false,
		},
		{
			name: "NOT with 'negation' type",
			expr: map[string]interface{}{
				"type": "negation",
				"expression": map[string]interface{}{
					"type":     "binaryOperation",
					"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "deleted"},
					"operator": "==",
					"right":    map[string]interface{}{"type": "booleanLiteral", "value": false},
				},
			},
			expected: ExprTypeNOT,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AnalyzeExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("AnalyzeExpression() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestAnalyzeExpression_NOT_Nested teste les expressions NOT imbriquées
func TestAnalyzeExpression_NOT_Nested(t *testing.T) {
	tests := []struct {
		name     string
		expr     interface{}
		expected ExpressionType
		wantErr  bool
	}{
		{
			name: "double NOT (NOT NOT expression)",
			expr: constraint.NotConstraint{
				Type: "notConstraint",
				Expression: constraint.NotConstraint{
					Type: "notConstraint",
					Expression: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "active"},
						Operator: "==",
						Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
					},
				},
			},
			expected: ExprTypeNOT,
			wantErr:  false,
		},
		{
			name: "NOT with OR expression inside",
			expr: constraint.NotConstraint{
				Type: "notConstraint",
				Expression: constraint.LogicalExpression{
					Type: "logicalExpr",
					Left: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
						Operator: "==",
						Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "inactive"},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Type:     "binaryOperation",
								Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
								Operator: "==",
								Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "deleted"},
							},
						},
					},
				},
			},
			expected: ExprTypeNOT,
			wantErr:  false,
		},
		{
			name: "NOT with Mixed expression inside",
			expr: constraint.NotConstraint{
				Type: "notConstraint",
				Expression: constraint.LogicalExpression{
					Type: "logicalExpr",
					Left: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
						Operator: ">",
						Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "AND",
							Right: constraint.BinaryOperation{
								Type:     "binaryOperation",
								Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
								Operator: ">=",
								Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
							},
						},
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Type:     "binaryOperation",
								Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "vip"},
								Operator: "==",
								Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
							},
						},
					},
				},
			},
			expected: ExprTypeNOT,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AnalyzeExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("AnalyzeExpression() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestGetExpressionInfo_NOT teste GetExpressionInfo avec expressions NOT
func TestGetExpressionInfo_NOT(t *testing.T) {
	// Test avec une expression NOT simple
	notExpr := constraint.NotConstraint{
		Type: "notConstraint",
		Expression: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "active"},
			Operator: "==",
			Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
		},
	}
	info, err := GetExpressionInfo(notExpr)
	if err != nil {
		t.Fatalf("GetExpressionInfo() error = %v", err)
	}
	if info.Type != ExprTypeNOT {
		t.Errorf("info.Type = %v, want %v", info.Type, ExprTypeNOT)
	}
	if !info.CanDecompose {
		t.Errorf("info.CanDecompose = false, want true")
	}
	if info.ShouldNormalize {
		t.Errorf("info.ShouldNormalize = true, want false")
	}
	// La complexité est maintenant 2 + complexité de l'expression interne (1 pour Simple)
	expectedComplexity := 3 // 2 (NOT) + 1 (Simple)
	if info.Complexity != expectedComplexity {
		t.Errorf("info.Complexity = %v, want %v", info.Complexity, expectedComplexity)
	}
	// Vérifier que l'expression interne a été analysée
	if info.InnerInfo == nil {
		t.Errorf("info.InnerInfo is nil, expected inner info")
	} else if info.InnerInfo.Type != ExprTypeSimple {
		t.Errorf("info.InnerInfo.Type = %v, want %v", info.InnerInfo.Type, ExprTypeSimple)
	}
	if info.RequiresBeta {
		t.Errorf("info.RequiresBeta = true, want false")
	}
	// Test avec une expression NOT complexe (NOT avec Mixed à l'intérieur)
	notMixedExpr := constraint.NotConstraint{
		Type: "notConstraint",
		Expression: constraint.LogicalExpression{
			Type: "logicalExpr",
			Left: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
				Operator: ">",
				Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
			},
			Operations: []constraint.LogicalOperation{
				{
					Op: "AND",
					Right: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
						Operator: ">=",
						Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
					},
				},
				{
					Op: "OR",
					Right: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "vip"},
						Operator: "==",
						Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
					},
				},
			},
		},
	}
	info2, err := GetExpressionInfo(notMixedExpr)
	if err != nil {
		t.Fatalf("GetExpressionInfo() error = %v", err)
	}
	// La négation d'une expression mixte reste une expression NOT
	if info2.Type != ExprTypeNOT {
		t.Errorf("info2.Type = %v, want %v", info2.Type, ExprTypeNOT)
	}
	if !info2.CanDecompose {
		t.Errorf("info2.CanDecompose = false, want true (NOT can be decomposed)")
	}
	if info2.ShouldNormalize {
		t.Errorf("info2.ShouldNormalize = true, want false")
	}
	// La complexité est maintenant 2 + complexité de l'expression interne (3 pour Mixed avec 2 ops)
	expectedComplexity2 := 5 // 2 (NOT) + 3 (Mixed: 1 left + 2 ops)
	if info2.Complexity != expectedComplexity2 {
		t.Errorf("info2.Complexity = %v, want %v", info2.Complexity, expectedComplexity2)
	}
	// Vérifier que l'expression interne a été analysée
	if info2.InnerInfo == nil {
		t.Errorf("info2.InnerInfo is nil, expected inner info")
	} else if info2.InnerInfo.Type != ExprTypeMixed {
		t.Errorf("info2.InnerInfo.Type = %v, want %v", info2.InnerInfo.Type, ExprTypeMixed)
	}
	if info2.RequiresBeta {
		t.Errorf("info2.RequiresBeta = true, want false")
	}
}

// TestAnalyzeExpression_Parenthesized teste l'analyse d'expressions parenthésées
func TestAnalyzeExpression_Parenthesized(t *testing.T) {
	tests := []struct {
		name     string
		expr     interface{}
		expected ExpressionType
		wantErr  bool
	}{
		{
			name: "Parenthesized simple expression",
			expr: map[string]interface{}{
				"type": "parenthesized",
				"expression": map[string]interface{}{
					"type":     "binaryOperation",
					"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
					"operator": ">",
					"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
				},
			},
			expected: ExprTypeSimple,
			wantErr:  false,
		},
		{
			name: "Parenthesized AND expression",
			expr: map[string]interface{}{
				"type": "parenthesizedExpression",
				"expr": map[string]interface{}{
					"type": "logicalExpr",
					"left": map[string]interface{}{
						"type":     "binaryOperation",
						"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
						"operator": ">",
						"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
					},
					"operations": []interface{}{
						map[string]interface{}{
							"op": "AND",
							"right": map[string]interface{}{
								"type":     "binaryOperation",
								"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "active"},
								"operator": "==",
								"right":    map[string]interface{}{"type": "booleanLiteral", "value": true},
							},
						},
					},
				},
			},
			expected: ExprTypeAND,
			wantErr:  false,
		},
		{
			name: "Parenthesized OR expression",
			expr: map[string]interface{}{
				"type": "group",
				"inner": map[string]interface{}{
					"type": "logicalExpr",
					"left": map[string]interface{}{
						"type":     "binaryOperation",
						"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "status"},
						"operator": "==",
						"right":    map[string]interface{}{"type": "stringLiteral", "value": "active"},
					},
					"operations": []interface{}{
						map[string]interface{}{
							"op": "OR",
							"right": map[string]interface{}{
								"type":     "binaryOperation",
								"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "status"},
								"operator": "==",
								"right":    map[string]interface{}{"type": "stringLiteral", "value": "pending"},
							},
						},
					},
				},
			},
			expected: ExprTypeOR,
			wantErr:  false,
		},
		{
			name: "Parenthesized Mixed expression",
			expr: map[string]interface{}{
				"type": "parenthesized",
				"expression": map[string]interface{}{
					"type": "logicalExpr",
					"left": map[string]interface{}{
						"type":     "binaryOperation",
						"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
						"operator": ">",
						"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
					},
					"operations": []interface{}{
						map[string]interface{}{
							"op": "AND",
							"right": map[string]interface{}{
								"type":     "binaryOperation",
								"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "salary"},
								"operator": ">=",
								"right":    map[string]interface{}{"type": "numberLiteral", "value": 50000},
							},
						},
						map[string]interface{}{
							"op": "OR",
							"right": map[string]interface{}{
								"type":     "binaryOperation",
								"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "vip"},
								"operator": "==",
								"right":    map[string]interface{}{"type": "booleanLiteral", "value": true},
							},
						},
					},
				},
			},
			expected: ExprTypeMixed,
			wantErr:  false,
		},
		{
			name: "Nested parenthesized expressions",
			expr: map[string]interface{}{
				"type": "parenthesized",
				"expression": map[string]interface{}{
					"type": "parenthesized",
					"expression": map[string]interface{}{
						"type":     "binaryOperation",
						"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
						"operator": ">",
						"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
					},
				},
			},
			expected: ExprTypeSimple,
			wantErr:  false,
		},
		{
			name: "Parenthesized expression without inner expression",
			expr: map[string]interface{}{
				"type": "parenthesized",
			},
			expected: ExprTypeSimple,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AnalyzeExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("AnalyzeExpression() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestAnalyzeInnerExpression teste l'analyse des expressions internes
func TestAnalyzeInnerExpression(t *testing.T) {
	tests := []struct {
		name     string
		expr     interface{}
		expected ExpressionType
		wantErr  bool
	}{
		{
			name: "NOT with simple inner expression",
			expr: map[string]interface{}{
				"type": "notConstraint",
				"constraint": map[string]interface{}{
					"type":     "binaryOperation",
					"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "active"},
					"operator": "==",
					"right":    map[string]interface{}{"type": "booleanLiteral", "value": true},
				},
			},
			expected: ExprTypeSimple,
			wantErr:  false,
		},
		{
			name: "NOT with AND inner expression",
			expr: map[string]interface{}{
				"type": "not",
				"expr": map[string]interface{}{
					"type": "logicalExpr",
					"left": map[string]interface{}{
						"type":     "binaryOperation",
						"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
						"operator": ">",
						"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
					},
					"operations": []interface{}{
						map[string]interface{}{
							"op": "AND",
							"right": map[string]interface{}{
								"type":     "binaryOperation",
								"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "active"},
								"operator": "==",
								"right":    map[string]interface{}{"type": "booleanLiteral", "value": true},
							},
						},
					},
				},
			},
			expected: ExprTypeAND,
			wantErr:  false,
		},
		{
			name: "NOT with OR inner expression",
			expr: map[string]interface{}{
				"type": "negation",
				"expression": map[string]interface{}{
					"type": "logicalExpr",
					"left": map[string]interface{}{
						"type":     "binaryOperation",
						"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "status"},
						"operator": "==",
						"right":    map[string]interface{}{"type": "stringLiteral", "value": "active"},
					},
					"operations": []interface{}{
						map[string]interface{}{
							"op": "OR",
							"right": map[string]interface{}{
								"type":     "binaryOperation",
								"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "status"},
								"operator": "==",
								"right":    map[string]interface{}{"type": "stringLiteral", "value": "pending"},
							},
						},
					},
				},
			},
			expected: ExprTypeOR,
			wantErr:  false,
		},
		{
			name: "NOT with Mixed inner expression",
			expr: map[string]interface{}{
				"type": "not",
				"expr": map[string]interface{}{
					"type": "logicalExpr",
					"left": map[string]interface{}{
						"type":     "binaryOperation",
						"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
						"operator": ">",
						"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
					},
					"operations": []interface{}{
						map[string]interface{}{
							"op": "AND",
							"right": map[string]interface{}{
								"type":     "binaryOperation",
								"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "salary"},
								"operator": ">=",
								"right":    map[string]interface{}{"type": "numberLiteral", "value": 50000},
							},
						},
						map[string]interface{}{
							"op": "OR",
							"right": map[string]interface{}{
								"type":     "binaryOperation",
								"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "vip"},
								"operator": "==",
								"right":    map[string]interface{}{"type": "booleanLiteral", "value": true},
							},
						},
					},
				},
			},
			expected: ExprTypeMixed,
			wantErr:  false,
		},
		{
			name: "Parenthesized with simple inner expression",
			expr: map[string]interface{}{
				"type": "parenthesized",
				"expression": map[string]interface{}{
					"type":     "binaryOperation",
					"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
					"operator": ">",
					"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
				},
			},
			expected: ExprTypeSimple,
			wantErr:  false,
		},
		{
			name:     "Expression without inner expression",
			expr:     map[string]interface{}{"type": "binaryOperation"},
			expected: ExprTypeSimple,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AnalyzeInnerExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeInnerExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("AnalyzeInnerExpression() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestGetExpressionInfo_WithInnerInfo teste GetExpressionInfo avec analyse des expressions internes
func TestGetExpressionInfo_WithInnerInfo(t *testing.T) {
	tests := []struct {
		name              string
		expr              interface{}
		expectedType      ExpressionType
		expectedInnerType ExpressionType
		hasInnerInfo      bool
	}{
		{
			name: "NOT with simple inner expression",
			expr: map[string]interface{}{
				"type": "notConstraint",
				"constraint": map[string]interface{}{
					"type":     "binaryOperation",
					"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "active"},
					"operator": "==",
					"right":    map[string]interface{}{"type": "booleanLiteral", "value": true},
				},
			},
			expectedType:      ExprTypeNOT,
			expectedInnerType: ExprTypeSimple,
			hasInnerInfo:      true,
		},
		{
			name: "NOT with AND inner expression",
			expr: map[string]interface{}{
				"type": "not",
				"expr": map[string]interface{}{
					"type": "logicalExpr",
					"left": map[string]interface{}{
						"type":     "binaryOperation",
						"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
						"operator": ">",
						"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
					},
					"operations": []interface{}{
						map[string]interface{}{
							"op": "AND",
							"right": map[string]interface{}{
								"type":     "binaryOperation",
								"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "salary"},
								"operator": ">=",
								"right":    map[string]interface{}{"type": "numberLiteral", "value": 50000},
							},
						},
					},
				},
			},
			expectedType:      ExprTypeNOT,
			expectedInnerType: ExprTypeAND,
			hasInnerInfo:      true,
		},
		{
			name: "NOT with OR inner expression",
			expr: map[string]interface{}{
				"type": "negation",
				"expression": map[string]interface{}{
					"type": "logicalExpr",
					"left": map[string]interface{}{
						"type":     "binaryOperation",
						"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "status"},
						"operator": "==",
						"right":    map[string]interface{}{"type": "stringLiteral", "value": "active"},
					},
					"operations": []interface{}{
						map[string]interface{}{
							"op": "OR",
							"right": map[string]interface{}{
								"type":     "binaryOperation",
								"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "status"},
								"operator": "==",
								"right":    map[string]interface{}{"type": "stringLiteral", "value": "pending"},
							},
						},
					},
				},
			},
			expectedType:      ExprTypeNOT,
			expectedInnerType: ExprTypeOR,
			hasInnerInfo:      true,
		},
		{
			name: "NOT with Mixed inner expression",
			expr: map[string]interface{}{
				"type": "not",
				"expr": map[string]interface{}{
					"type": "logicalExpr",
					"left": map[string]interface{}{
						"type":     "binaryOperation",
						"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
						"operator": ">",
						"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
					},
					"operations": []interface{}{
						map[string]interface{}{
							"op": "AND",
							"right": map[string]interface{}{
								"type":     "binaryOperation",
								"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "salary"},
								"operator": ">=",
								"right":    map[string]interface{}{"type": "numberLiteral", "value": 50000},
							},
						},
						map[string]interface{}{
							"op": "OR",
							"right": map[string]interface{}{
								"type":     "binaryOperation",
								"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "vip"},
								"operator": "==",
								"right":    map[string]interface{}{"type": "booleanLiteral", "value": true},
							},
						},
					},
				},
			},
			expectedType:      ExprTypeNOT,
			expectedInnerType: ExprTypeMixed,
			hasInnerInfo:      true,
		},
		{
			name: "Simple expression (no inner info)",
			expr: map[string]interface{}{
				"type":     "binaryOperation",
				"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
				"operator": ">",
				"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
			},
			expectedType: ExprTypeSimple,
			hasInnerInfo: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := GetExpressionInfo(tt.expr)
			if err != nil {
				t.Fatalf("GetExpressionInfo() error = %v", err)
			}
			if info.Type != tt.expectedType {
				t.Errorf("info.Type = %v, want %v", info.Type, tt.expectedType)
			}
			if tt.hasInnerInfo {
				if info.InnerInfo == nil {
					t.Errorf("info.InnerInfo is nil, expected inner info")
				} else {
					if info.InnerInfo.Type != tt.expectedInnerType {
						t.Errorf("info.InnerInfo.Type = %v, want %v", info.InnerInfo.Type, tt.expectedInnerType)
					}
					// Vérifier que la complexité a été ajustée
					expectedComplexity := 2 + info.InnerInfo.Complexity
					if info.Complexity != expectedComplexity {
						t.Errorf("info.Complexity = %v, want %v (2 + inner complexity)", info.Complexity, expectedComplexity)
					}
				}
			} else {
				if info.InnerInfo != nil {
					t.Errorf("info.InnerInfo is not nil, expected no inner info")
				}
			}
		})
	}
}

// TestNestedParenthesizedAndNOT teste les expressions avec parenthèses et NOT imbriqués
func TestNestedParenthesizedAndNOT(t *testing.T) {
	tests := []struct {
		name     string
		expr     interface{}
		expected ExpressionType
	}{
		{
			name: "NOT with parenthesized expression",
			expr: map[string]interface{}{
				"type": "not",
				"expr": map[string]interface{}{
					"type": "parenthesized",
					"expression": map[string]interface{}{
						"type": "logicalExpr",
						"left": map[string]interface{}{
							"type":     "binaryOperation",
							"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
							"operator": ">",
							"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
						},
						"operations": []interface{}{
							map[string]interface{}{
								"op": "AND",
								"right": map[string]interface{}{
									"type":     "binaryOperation",
									"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "active"},
									"operator": "==",
									"right":    map[string]interface{}{"type": "booleanLiteral", "value": true},
								},
							},
						},
					},
				},
			},
			expected: ExprTypeNOT,
		},
		{
			name: "Parenthesized NOT expression",
			expr: map[string]interface{}{
				"type": "parenthesized",
				"expression": map[string]interface{}{
					"type": "not",
					"expr": map[string]interface{}{
						"type":     "binaryOperation",
						"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "active"},
						"operator": "==",
						"right":    map[string]interface{}{"type": "booleanLiteral", "value": true},
					},
				},
			},
			expected: ExprTypeNOT,
		},
		{
			name: "Multiple levels of parentheses with NOT",
			expr: map[string]interface{}{
				"type": "parenthesized",
				"expression": map[string]interface{}{
					"type": "not",
					"expr": map[string]interface{}{
						"type": "parenthesized",
						"expression": map[string]interface{}{
							"type": "logicalExpr",
							"left": map[string]interface{}{
								"type":     "binaryOperation",
								"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
								"operator": ">",
								"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
							},
							"operations": []interface{}{
								map[string]interface{}{
									"op": "OR",
									"right": map[string]interface{}{
										"type":     "binaryOperation",
										"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "vip"},
										"operator": "==",
										"right":    map[string]interface{}{"type": "booleanLiteral", "value": true},
									},
								},
							},
						},
					},
				},
			},
			expected: ExprTypeNOT,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AnalyzeExpression(tt.expr)
			if err != nil {
				t.Fatalf("AnalyzeExpression() error = %v", err)
			}
			if got != tt.expected {
				t.Errorf("AnalyzeExpression() = %v, want %v", got, tt.expected)
			}
		})
	}
}
