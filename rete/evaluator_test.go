// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/treivax/tsd/constraint"
)

func TestAlphaConditionEvaluator_evaluateConstraint(t *testing.T) {
	fact := &Fact{
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Alice",
			"age":  30,
		},
	}
	evaluator := NewAlphaConditionEvaluator()
	evaluator.variableBindings["p"] = fact
	tests := []struct {
		name        string
		constraint  constraint.Constraint
		expected    bool
		shouldError bool
	}{
		{
			name: "simple equality constraint",
			constraint: constraint.Constraint{
				Type:     "constraint",
				Operator: "==",
				Left: map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "name",
				},
				Right: map[string]interface{}{
					"type":  "stringLiteral",
					"value": "Alice",
				},
			},
			expected: true,
		},
		{
			name: "numeric comparison constraint",
			constraint: constraint.Constraint{
				Type:     "constraint",
				Operator: ">",
				Left: map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "age",
				},
				Right: map[string]interface{}{
					"type":  "numberLiteral",
					"value": 25.0,
				},
			},
			expected: true,
		},
		{
			name: "inequality constraint",
			constraint: constraint.Constraint{
				Type:     "constraint",
				Operator: "!=",
				Left: map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "name",
				},
				Right: map[string]interface{}{
					"type":  "stringLiteral",
					"value": "Bob",
				},
			},
			expected: true,
		},
		{
			name: "less than or equal constraint",
			constraint: constraint.Constraint{
				Type:     "constraint",
				Operator: "<=",
				Left: map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "age",
				},
				Right: map[string]interface{}{
					"type":  "numberLiteral",
					"value": 30.0,
				},
			},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateConstraint(tt.constraint)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
func TestAlphaConditionEvaluator_evaluateExpression(t *testing.T) {
	fact := &Fact{
		Type: "Person",
		Fields: map[string]interface{}{
			"name":   "Alice",
			"age":    30,
			"active": true,
		},
	}
	evaluator := NewAlphaConditionEvaluator()
	evaluator.variableBindings["p"] = fact
	tests := []struct {
		name        string
		expr        interface{}
		expected    bool
		shouldError bool
	}{
		{
			name: "boolean literal true",
			expr: constraint.BooleanLiteral{
				Type:  "booleanLiteral",
				Value: true,
			},
			expected: true,
		},
		{
			name: "boolean literal false",
			expr: constraint.BooleanLiteral{
				Type:  "booleanLiteral",
				Value: false,
			},
			expected: false,
		},
		{
			name: "map expression with boolean literal",
			expr: map[string]interface{}{
				"type":  "booleanLiteral",
				"value": true,
			},
			expected: true,
		},
		{
			name: "binary operation map",
			expr: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "==",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "name",
				},
				"right": map[string]interface{}{
					"type":  "stringLiteral",
					"value": "Alice",
				},
			},
			expected: true,
		},
		{
			name: "comparison type",
			expr: map[string]interface{}{
				"type":     "comparison",
				"operator": ">",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "age",
				},
				"right": map[string]interface{}{
					"type":  "numberLiteral",
					"value": 25.0,
				},
			},
			expected: true,
		},
		{
			name: "simple type always true",
			expr: map[string]interface{}{
				"type": "simple",
			},
			expected: true,
		},
		{
			name: "logical expression AND",
			expr: map[string]interface{}{
				"type": "logicalExpression",
				"left": map[string]interface{}{
					"type":  "booleanLiteral",
					"value": true,
				},
				"operations": []interface{}{
					map[string]interface{}{
						"op": "AND",
						"right": map[string]interface{}{
							"type":  "booleanLiteral",
							"value": true,
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "logical expression OR",
			expr: map[string]interface{}{
				"type": "logicalExpr",
				"left": map[string]interface{}{
					"type":  "booleanLiteral",
					"value": false,
				},
				"operations": []interface{}{
					map[string]interface{}{
						"op": "OR",
						"right": map[string]interface{}{
							"type":  "booleanLiteral",
							"value": true,
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "negation constraint",
			expr: map[string]interface{}{
				"type": "negation",
				"condition": map[string]interface{}{
					"type":  "booleanLiteral",
					"value": false,
				},
			},
			expected: true,
		},
		{
			name: "notConstraint",
			expr: map[string]interface{}{
				"type": "notConstraint",
				"expression": map[string]interface{}{
					"type":  "booleanLiteral",
					"value": false,
				},
			},
			expected: true,
		},
		{
			name: "existsConstraint",
			expr: map[string]interface{}{
				"type":     "existsConstraint",
				"variable": "x",
			},
			expected: true, // Usually true based on hash
		},
		{
			name: "unsupported expression type",
			expr: map[string]interface{}{
				"type": "unknownType",
			},
			shouldError: true,
		},
		{
			name:        "unsupported Go type",
			expr:        123, // int is not supported
			shouldError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateExpression(tt.expr)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
func TestAlphaConditionEvaluator_evaluateValue(t *testing.T) {
	fact := &Fact{
		Type: "Person",
		Fields: map[string]interface{}{
			"name":     "Alice",
			"age":      30,
			"active":   true,
			"salary":   50000.5,
			"nickname": nil,
		},
	}
	evaluator := NewAlphaConditionEvaluator()
	evaluator.variableBindings["p"] = fact
	tests := []struct {
		name        string
		value       interface{}
		expected    interface{}
		shouldError bool
	}{
		// Literal values
		{
			name: "string literal",
			value: map[string]interface{}{
				"type":  "stringLiteral",
				"value": "test",
			},
			expected: "test",
		},
		{
			name: "number literal",
			value: map[string]interface{}{
				"type":  "numberLiteral",
				"value": 42.0,
			},
			expected: 42.0,
		},
		{
			name: "boolean literal true",
			value: map[string]interface{}{
				"type":  "booleanLiteral",
				"value": true,
			},
			expected: true,
		},
		{
			name: "boolean literal false",
			value: map[string]interface{}{
				"type":  "booleanLiteral",
				"value": false,
			},
			expected: false,
		},
		// Field access
		{
			name: "field access - string",
			value: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "name",
			},
			expected: "Alice",
		},
		{
			name: "field access - number",
			value: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "age",
			},
			expected: 30,
		},
		{
			name: "field access - boolean",
			value: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "active",
			},
			expected: true,
		},
		{
			name: "field access - float",
			value: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "salary",
			},
			expected: 50000.5,
		},
		{
			name: "field access - nil value",
			value: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "nickname",
			},
			expected: nil,
		},
		// Direct values
		{
			name:     "direct string",
			value:    "hello",
			expected: "hello",
		},
		{
			name:     "direct number",
			value:    42,
			expected: 42,
		},
		{
			name:     "direct float",
			value:    3.14,
			expected: 3.14,
		},
		{
			name:     "direct bool",
			value:    true,
			expected: true,
		},
		// Constraint types
		{
			name: "string type",
			value: constraint.StringLiteral{
				Type:  "stringLiteral",
				Value: "constraint string",
			},
			expected: "constraint string",
		},
		{
			name: "number type",
			value: constraint.NumberLiteral{
				Type:  "numberLiteral",
				Value: 99.9,
			},
			expected: 99.9,
		},
		{
			name: "boolean type",
			value: constraint.BooleanLiteral{
				Type:  "booleanLiteral",
				Value: true,
			},
			expected: true,
		},
		{
			name: "field access type",
			value: constraint.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "name",
			},
			expected: "Alice",
		},
		// Error cases
		{
			name: "field access - non-existent field",
			value: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "nonexistent",
			},
			shouldError: true,
		},
		{
			name: "invalid field access - missing field key",
			value: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
			},
			shouldError: true,
		},
		{
			name: "map without type",
			value: map[string]interface{}{
				"value": "test",
			},
			shouldError: true,
		},
		{
			name: "unknown type in map",
			value: map[string]interface{}{
				"type":  "unknownType",
				"value": "test",
			},
			shouldError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateValue(tt.value)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
func TestAlphaConditionEvaluator_evaluateLogicalExpression(t *testing.T) {
	fact := &Fact{
		Type: "Test",
		Fields: map[string]interface{}{
			"value": 10,
		},
	}
	evaluator := NewAlphaConditionEvaluator()
	evaluator.variableBindings["p"] = fact
	tests := []struct {
		name        string
		expr        constraint.LogicalExpression
		expected    bool
		shouldError bool
	}{
		{
			name: "AND with two true values",
			expr: constraint.LogicalExpression{
				Type: "logicalExpression",
				Left: constraint.BooleanLiteral{Value: true},
				Operations: []constraint.LogicalOperation{
					{
						Op:    "AND",
						Right: constraint.BooleanLiteral{Value: true},
					},
				},
			},
			expected: true,
		},
		{
			name: "AND with false left",
			expr: constraint.LogicalExpression{
				Type: "logicalExpression",
				Left: constraint.BooleanLiteral{Value: false},
				Operations: []constraint.LogicalOperation{
					{
						Op:    "AND",
						Right: constraint.BooleanLiteral{Value: true},
					},
				},
			},
			expected: false,
		},
		{
			name: "OR with false and true",
			expr: constraint.LogicalExpression{
				Type: "logicalExpression",
				Left: constraint.BooleanLiteral{Value: false},
				Operations: []constraint.LogicalOperation{
					{
						Op:    "OR",
						Right: constraint.BooleanLiteral{Value: true},
					},
				},
			},
			expected: true,
		},
		{
			name: "multiple AND operations",
			expr: constraint.LogicalExpression{
				Type: "logicalExpression",
				Left: constraint.BooleanLiteral{Value: true},
				Operations: []constraint.LogicalOperation{
					{
						Op:    "AND",
						Right: constraint.BooleanLiteral{Value: true},
					},
					{
						Op:    "AND",
						Right: constraint.BooleanLiteral{Value: true},
					},
				},
			},
			expected: true,
		},
		{
			name: "mixed AND and OR operations",
			expr: constraint.LogicalExpression{
				Type: "logicalExpression",
				Left: constraint.BooleanLiteral{Value: false},
				Operations: []constraint.LogicalOperation{
					{
						Op:    "OR",
						Right: constraint.BooleanLiteral{Value: true},
					},
					{
						Op:    "AND",
						Right: constraint.BooleanLiteral{Value: true},
					},
				},
			},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateLogicalExpression(tt.expr)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
func TestAlphaConditionEvaluator_evaluateConstraintMap(t *testing.T) {
	fact := &Fact{
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Alice",
			"age":  30,
		},
	}
	evaluator := NewAlphaConditionEvaluator()
	evaluator.variableBindings["p"] = fact
	tests := []struct {
		name        string
		expr        map[string]interface{}
		expected    bool
		shouldError bool
	}{
		{
			name: "constraint with indirection",
			expr: map[string]interface{}{
				"constraint": map[string]interface{}{
					"type":     "binaryOperation",
					"operator": "==",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "p",
						"field":  "name",
					},
					"right": map[string]interface{}{
						"type":  "stringLiteral",
						"value": "Alice",
					},
				},
			},
			expected: true,
		},
		{
			name: "simple type constraint",
			expr: map[string]interface{}{
				"type": "simple",
			},
			expected: true,
		},
		{
			name: "passthrough type constraint",
			expr: map[string]interface{}{
				"type": "passthrough",
			},
			expected: true,
		},
		{
			name: "exists type constraint",
			expr: map[string]interface{}{
				"type": "exists",
			},
			expected: true,
		},
		{
			name: "direct constraint without indirection",
			expr: map[string]interface{}{
				"operator": "==",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "age",
				},
				"right": map[string]interface{}{
					"type":  "numberLiteral",
					"value": 30.0,
				},
			},
			expected: true,
		},
		{
			name: "missing operator",
			expr: map[string]interface{}{
				"left": map[string]interface{}{
					"type":  "stringLiteral",
					"value": "test",
				},
				"right": map[string]interface{}{
					"type":  "stringLiteral",
					"value": "test",
				},
			},
			shouldError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateConstraintMap(tt.expr)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
func TestAlphaConditionEvaluator_evaluateBinaryOperation(t *testing.T) {
	fact := &Fact{
		Type: "Test",
		Fields: map[string]interface{}{
			"value": 42,
		},
	}
	evaluator := NewAlphaConditionEvaluator()
	evaluator.variableBindings["p"] = fact
	tests := []struct {
		name        string
		op          constraint.BinaryOperation
		expected    bool
		shouldError bool
	}{
		{
			name: "equality comparison",
			op: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Operator: "==",
				Left:     constraint.NumberLiteral{Value: 42},
				Right:    constraint.NumberLiteral{Value: 42},
			},
			expected: true,
		},
		{
			name: "inequality comparison",
			op: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Operator: "!=",
				Left:     constraint.NumberLiteral{Value: 42},
				Right:    constraint.NumberLiteral{Value: 10},
			},
			expected: true,
		},
		{
			name: "greater than comparison",
			op: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Operator: ">",
				Left:     constraint.NumberLiteral{Value: 50},
				Right:    constraint.NumberLiteral{Value: 30},
			},
			expected: true,
		},
		{
			name: "less than comparison",
			op: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Operator: "<",
				Left:     constraint.NumberLiteral{Value: 10},
				Right:    constraint.NumberLiteral{Value: 20},
			},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateBinaryOperation(tt.op)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestEvaluator_AccessIDField teste l'accès au champ spécial 'id' dans les expressions
func TestEvaluator_AccessIDField(t *testing.T) {
	t.Log("🧪 TEST: Evaluator Access ID Field - Accès au champ 'id'")
	t.Log("=========================================================")

	fact := &Fact{
		ID:   "Person~Alice",
		Type: "Person",
		Fields: map[string]interface{}{
			"nom": "Alice",
			"age": 30,
		},
	}

	evaluator := NewAlphaConditionEvaluator()
	evaluator.variableBindings["p"] = fact

	t.Run("accès au champ id", func(t *testing.T) {
		idValue, err := evaluator.evaluateFieldAccessByName("p", "id")
		if err != nil {
			t.Fatalf("❌ evaluateFieldAccessByName(p, 'id') erreur = %v", err)
		}

		if idValue != "Person~Alice" {
			t.Errorf("❌ evaluateFieldAccessByName(p, 'id') = %v, attendu %v", idValue, "Person~Alice")
		} else {
			t.Logf("✅ Accès au champ 'id' réussi: %v", idValue)
		}
	})

	t.Run("accès au champ nom normal", func(t *testing.T) {
		nomValue, err := evaluator.evaluateFieldAccessByName("p", "nom")
		if err != nil {
			t.Fatalf("❌ evaluateFieldAccessByName(p, 'nom') erreur = %v", err)
		}

		if nomValue != "Alice" {
			t.Errorf("❌ evaluateFieldAccessByName(p, 'nom') = %v, attendu %v", nomValue, "Alice")
		} else {
			t.Logf("✅ Accès au champ 'nom' réussi: %v", nomValue)
		}
	})

	t.Run("expression complète avec accès à id", func(t *testing.T) {
		// Test avec une contrainte complète utilisant le champ id
		constraint := map[string]interface{}{
			"type":     "constraint",
			"operator": "==",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "id",
			},
			"right": map[string]interface{}{
				"type":  "stringLiteral",
				"value": "Person~Alice",
			},
		}

		result, err := evaluator.evaluateConstraintMap(constraint)
		if err != nil {
			t.Fatalf("❌ Évaluation de la contrainte id erreur = %v", err)
		}

		if !result {
			t.Errorf("❌ Expression 'p.id == \"Person~Alice\"' devrait être vraie")
		} else {
			t.Logf("✅ Expression 'p.id == \"Person~Alice\"' évaluée correctement")
		}
	})

	t.Run("expression avec id dans jointure", func(t *testing.T) {
		fact2 := &Fact{
			ID:   "Order~O123",
			Type: "Order",
			Fields: map[string]interface{}{
				"userId": "Person~Alice",
				"amount": 100.0,
			},
		}

		evaluator.variableBindings["o"] = fact2

		// Test jointure: o.userId == p.id
		constraint := map[string]interface{}{
			"type":     "constraint",
			"operator": "==",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "o",
				"field":  "userId",
			},
			"right": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "id",
			},
		}

		result, err := evaluator.evaluateConstraintMap(constraint)
		if err != nil {
			t.Fatalf("❌ Évaluation de la jointure erreur = %v", err)
		}

		if !result {
			t.Errorf("❌ Expression 'o.userId == p.id' devrait être vraie")
		} else {
			t.Logf("✅ Jointure 'o.userId == p.id' évaluée correctement")
		}
	})

	t.Run("accès à _id_ interdit", func(t *testing.T) {
		// Vérifier que la constante FieldNameID est bien mise à jour
		if FieldNameID != "_id_" {
			t.Errorf("❌ FieldNameID = %v, attendu '_id_'", FieldNameID)
		} else {
			t.Logf("✅ Constante FieldNameID = %v", FieldNameID)
		}

		// Tenter d'utiliser la constante dans l'évaluation - devrait échouer
		_, err := evaluator.evaluateFieldAccessByName("p", FieldNameID)
		if err == nil {
			t.Errorf("❌ Attendu une erreur pour accès à _id_, reçu nil")
		} else {
			t.Logf("✅ Accès à _id_ correctement interdit: %v", err)
		}
	})

	t.Log("✅ Test complet: Evaluator Access ID Field")
}
