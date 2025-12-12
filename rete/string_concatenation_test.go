// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestStringConcatenation teste la concaténation de chaînes avec l'opérateur +
func TestStringConcatenation(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	// Créer un fait pour les tests
	product := &Fact{
		ID:   "p1",
		Type: "Product",
		Fields: map[string]interface{}{
			"name":     "Widget",
			"price":    99.99,
			"category": "Electronics",
		},
	}
	var capturedArgs []interface{}
	testHandler := &captureActionHandler{
		name:         "testConcat",
		capturedArgs: &capturedArgs,
	}
	err := env.Network.ActionExecutor.RegisterAction(testHandler)
	require.NoError(t, err)
	t.Run("concatenate two string literals", func(t *testing.T) {
		capturedArgs = nil
		action := &Action{
			Jobs: []JobCall{
				{
					Name: "testConcat",
					Args: []interface{}{
						map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "+",
							"left": map[string]interface{}{
								"type":  "string",
								"value": "Hello ",
							},
							"right": map[string]interface{}{
								"type":  "string",
								"value": "World",
							},
						},
					},
				},
			},
		}
		token := &Token{
			Facts:    []*Fact{product},
			Bindings: NewBindingChainWith("p", product),
		}
		err := env.Network.ActionExecutor.ExecuteAction(action, token)
		require.NoError(t, err, "String concatenation should work")
		require.Len(t, capturedArgs, 1)
		result, ok := capturedArgs[0].(string)
		require.True(t, ok, "Expected string result, got %T", capturedArgs[0])
		assert.Equal(t, "Hello World", result, "Should concatenate strings")
	})
	t.Run("concatenate string literal and field access", func(t *testing.T) {
		capturedArgs = nil
		action := &Action{
			Jobs: []JobCall{
				{
					Name: "testConcat",
					Args: []interface{}{
						map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "+",
							"left": map[string]interface{}{
								"type":  "string",
								"value": "Product: ",
							},
							"right": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "p",
								"field":  "name",
							},
						},
					},
				},
			},
		}
		token := &Token{
			Facts:    []*Fact{product},
			Bindings: NewBindingChainWith("p", product),
		}
		err := env.Network.ActionExecutor.ExecuteAction(action, token)
		require.NoError(t, err, "String concatenation with field access should work")
		require.Len(t, capturedArgs, 1)
		result, ok := capturedArgs[0].(string)
		require.True(t, ok, "Expected string result, got %T", capturedArgs[0])
		assert.Equal(t, "Product: Widget", result, "Should concatenate string and field")
	})
	t.Run("concatenate string with cast number", func(t *testing.T) {
		capturedArgs = nil
		action := &Action{
			Jobs: []JobCall{
				{
					Name: "testConcat",
					Args: []interface{}{
						map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "+",
							"left": map[string]interface{}{
								"type":  "string",
								"value": "Price: $",
							},
							"right": map[string]interface{}{
								"type":     "cast",
								"castType": "string",
								"expression": map[string]interface{}{
									"type":   "fieldAccess",
									"object": "p",
									"field":  "price",
								},
							},
						},
					},
				},
			},
		}
		token := &Token{
			Facts:    []*Fact{product},
			Bindings: NewBindingChainWith("p", product),
		}
		err := env.Network.ActionExecutor.ExecuteAction(action, token)
		require.NoError(t, err, "String concatenation with cast should work")
		require.Len(t, capturedArgs, 1)
		result, ok := capturedArgs[0].(string)
		require.True(t, ok, "Expected string result, got %T", capturedArgs[0])
		assert.Equal(t, "Price: $99.99", result, "Should concatenate string and cast number")
	})
	t.Run("concatenate multiple strings", func(t *testing.T) {
		capturedArgs = nil
		// "Category: " + p.category + " - " + p.name
		action := &Action{
			Jobs: []JobCall{
				{
					Name: "testConcat",
					Args: []interface{}{
						map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "+",
							"left": map[string]interface{}{
								"type":     "binaryOperation",
								"operator": "+",
								"left": map[string]interface{}{
									"type":     "binaryOperation",
									"operator": "+",
									"left": map[string]interface{}{
										"type":  "string",
										"value": "Category: ",
									},
									"right": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "p",
										"field":  "category",
									},
								},
								"right": map[string]interface{}{
									"type":  "string",
									"value": " - ",
								},
							},
							"right": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "p",
								"field":  "name",
							},
						},
					},
				},
			},
		}
		token := &Token{
			Facts:    []*Fact{product},
			Bindings: NewBindingChainWith("p", product),
		}
		err := env.Network.ActionExecutor.ExecuteAction(action, token)
		require.NoError(t, err, "Multiple string concatenation should work")
		require.Len(t, capturedArgs, 1)
		result, ok := capturedArgs[0].(string)
		require.True(t, ok, "Expected string result, got %T", capturedArgs[0])
		assert.Equal(t, "Category: Electronics - Widget", result, "Should concatenate multiple strings")
	})
	t.Run("empty string concatenation", func(t *testing.T) {
		capturedArgs = nil
		action := &Action{
			Jobs: []JobCall{
				{
					Name: "testConcat",
					Args: []interface{}{
						map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "+",
							"left": map[string]interface{}{
								"type":  "string",
								"value": "",
							},
							"right": map[string]interface{}{
								"type":  "string",
								"value": "test",
							},
						},
					},
				},
			},
		}
		token := &Token{
			Facts:    []*Fact{product},
			Bindings: NewBindingChainWith("p", product),
		}
		err := env.Network.ActionExecutor.ExecuteAction(action, token)
		require.NoError(t, err, "Empty string concatenation should work")
		require.Len(t, capturedArgs, 1)
		result, ok := capturedArgs[0].(string)
		require.True(t, ok, "Expected string result, got %T", capturedArgs[0])
		assert.Equal(t, "test", result, "Should handle empty strings")
	})
}

// TestStringConcatenationInConditions teste la concaténation dans les conditions
func TestStringConcatenationInConditions(t *testing.T) {
	t.Parallel()
	evaluator := NewAlphaConditionEvaluator()
	product := &Fact{
		ID:   "p1",
		Type: "Product",
		Fields: map[string]interface{}{
			"name":     "Widget",
			"category": "Electronics",
		},
	}
	evaluator.variableBindings["p"] = product
	t.Run("concatenate and compare", func(t *testing.T) {
		// Build expression: ("Category: " + p.category) == "Category: Electronics"
		condition := map[string]interface{}{
			"type":     "comparison",
			"operator": "==",
			"left": map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "+",
				"left": map[string]interface{}{
					"type":  "string",
					"value": "Category: ",
				},
				"right": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "category",
				},
			},
			"right": map[string]interface{}{
				"type":  "string",
				"value": "Category: Electronics",
			},
		}
		// Pour les conditions, on utilise evaluateComparison directement
		leftVal, err := evaluator.evaluateValue(condition["left"])
		require.NoError(t, err, "Should evaluate left side")
		rightVal, err := evaluator.evaluateValue(condition["right"])
		require.NoError(t, err, "Should evaluate right side")
		assert.Equal(t, rightVal, leftVal, "Concatenated string should match expected value")
	})
}

// TestMixedArithmeticAndStringOps teste que les opérations numériques et strings ne se mélangent pas
func TestMixedArithmeticAndStringOps(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	product := &Fact{
		ID:   "p1",
		Type: "Product",
		Fields: map[string]interface{}{
			"price":    99.99,
			"quantity": 5.0,
		},
	}
	var capturedArgs []interface{}
	testHandler := &captureActionHandler{
		name:         "testOp",
		capturedArgs: &capturedArgs,
	}
	err := env.Network.ActionExecutor.RegisterAction(testHandler)
	require.NoError(t, err)
	t.Run("number addition should return number", func(t *testing.T) {
		capturedArgs = nil
		action := &Action{
			Jobs: []JobCall{
				{
					Name: "testOp",
					Args: []interface{}{
						map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "+",
							"left": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "p",
								"field":  "price",
							},
							"right": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "p",
								"field":  "quantity",
							},
						},
					},
				},
			},
		}
		token := &Token{
			Facts:    []*Fact{product},
			Bindings: NewBindingChainWith("p", product),
		}
		err := env.Network.ActionExecutor.ExecuteAction(action, token)
		require.NoError(t, err, "Number addition should work")
		require.Len(t, capturedArgs, 1)
		result, ok := capturedArgs[0].(float64)
		require.True(t, ok, "Expected float64 result, got %T", capturedArgs[0])
		assert.Equal(t, 104.99, result, "Should add numbers")
	})
	t.Run("string concatenation should return string", func(t *testing.T) {
		capturedArgs = nil
		action := &Action{
			Jobs: []JobCall{
				{
					Name: "testOp",
					Args: []interface{}{
						map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "+",
							"left": map[string]interface{}{
								"type":  "string",
								"value": "Hello ",
							},
							"right": map[string]interface{}{
								"type":  "string",
								"value": "World",
							},
						},
					},
				},
			},
		}
		token := &Token{
			Facts:    []*Fact{product},
			Bindings: NewBindingChainWith("p", product),
		}
		err := env.Network.ActionExecutor.ExecuteAction(action, token)
		require.NoError(t, err, "String concatenation should work")
		require.Len(t, capturedArgs, 1)
		result, ok := capturedArgs[0].(string)
		require.True(t, ok, "Expected string result, got %T", capturedArgs[0])
		assert.Equal(t, "Hello World", result, "Should concatenate strings")
	})
}
