// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureActionHandler is a simple action handler for testing
type captureActionHandler struct {
	name         string
	capturedArgs *[]interface{}
}

func (c *captureActionHandler) GetName() string {
	return c.name
}
func (c *captureActionHandler) Validate(args []interface{}) error {
	return nil
}
func (c *captureActionHandler) Execute(args []interface{}, ctx *ExecutionContext) error {
	*c.capturedArgs = args
	return nil
}

// TestCastInActions teste l'utilisation des opérateurs de cast dans les actions
func TestCastInActions(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	// Créer un fait Product avec des valeurs à caster
	product := &Fact{
		ID:   "p1",
		Type: "Product",
		Fields: map[string]interface{}{
			"name":     "Widget",
			"price":    99.99,
			"quantity": 10.0,
			"active":   "true",
		},
	}
	// Variable pour capturer les arguments
	var capturedArgs []interface{}
	// Créer un handler personnalisé pour capturer les résultats
	captureHandler := &captureActionHandler{
		name:         "testCast",
		capturedArgs: &capturedArgs,
	}
	err := env.Network.ActionExecutor.RegisterAction(captureHandler)
	require.NoError(t, err)
	t.Run("cast number to string in action", func(t *testing.T) {
		capturedArgs = nil
		action := &Action{
			Jobs: []JobCall{
				{
					Name: "testCast",
					Args: []interface{}{
						map[string]interface{}{
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
		}
		token := &Token{
			Facts:    []*Fact{product},
			Bindings: NewBindingChainWith("p", product),
		}
		err := env.Network.ActionExecutor.ExecuteAction(action, token)
		require.NoError(t, err, "Action with cast should execute successfully")
		require.Len(t, capturedArgs, 1)
		result, ok := capturedArgs[0].(string)
		require.True(t, ok, "Expected string result from cast, got %T", capturedArgs[0])
		assert.Equal(t, "99.99", result, "Cast should convert number to string")
	})
	t.Run("cast number to string (integer)", func(t *testing.T) {
		capturedArgs = nil
		action := &Action{
			Jobs: []JobCall{
				{
					Name: "testCast",
					Args: []interface{}{
						map[string]interface{}{
							"type":     "cast",
							"castType": "string",
							"expression": map[string]interface{}{
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
		require.NoError(t, err, "Action with cast should execute successfully")
		require.Len(t, capturedArgs, 1)
		result, ok := capturedArgs[0].(string)
		require.True(t, ok, "Expected string result from cast")
		assert.Equal(t, "10", result, "Cast should convert number to string")
	})
	t.Run("cast string to bool in action", func(t *testing.T) {
		capturedArgs = nil
		action := &Action{
			Jobs: []JobCall{
				{
					Name: "testCast",
					Args: []interface{}{
						map[string]interface{}{
							"type":     "cast",
							"castType": "bool",
							"expression": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "p",
								"field":  "active",
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
		require.NoError(t, err, "Action with cast should execute successfully")
		require.Len(t, capturedArgs, 1)
		result, ok := capturedArgs[0].(bool)
		require.True(t, ok, "Expected bool result from cast")
		assert.Equal(t, true, result, "Cast should convert string to bool")
	})
	t.Run("cast literal number to string", func(t *testing.T) {
		capturedArgs = nil
		action := &Action{
			Jobs: []JobCall{
				{
					Name: "testCast",
					Args: []interface{}{
						map[string]interface{}{
							"type":     "cast",
							"castType": "string",
							"expression": map[string]interface{}{
								"type":  "number",
								"value": 123.0,
							},
						},
					},
				},
			},
		}
		token := &Token{
			Facts: []*Fact{product},
		}
		err := env.Network.ActionExecutor.ExecuteAction(action, token)
		require.NoError(t, err, "Action with cast of literal should execute successfully")
		require.Len(t, capturedArgs, 1)
		result, ok := capturedArgs[0].(string)
		require.True(t, ok, "Expected string result from cast")
		assert.Equal(t, "123", result, "Cast should convert literal number to string")
	})
	t.Run("cast in arithmetic expression", func(t *testing.T) {
		capturedArgs = nil
		action := &Action{
			Jobs: []JobCall{
				{
					Name: "testCast",
					Args: []interface{}{
						map[string]interface{}{
							"type":     "cast",
							"castType": "string",
							"expression": map[string]interface{}{
								"type":     "binaryOperation",
								"operator": "*",
								"left": map[string]interface{}{
									"type":   "fieldAccess",
									"object": "p",
									"field":  "quantity",
								},
								"right": map[string]interface{}{
									"type":  "number",
									"value": 2.0,
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
		require.NoError(t, err, "Action with cast of expression should execute successfully")
		require.Len(t, capturedArgs, 1)
		result, ok := capturedArgs[0].(string)
		require.True(t, ok, "Expected string result from cast")
		assert.Equal(t, "20", result, "Cast should convert arithmetic result to string")
	})
}

// TestCastErrorHandlingInActions teste la gestion des erreurs de cast dans les actions
func TestCastErrorHandlingInActions(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	product := &Fact{
		ID:   "p1",
		Type: "Product",
		Fields: map[string]interface{}{
			"invalidNumber": "not-a-number",
		},
	}
	t.Run("invalid string to number cast", func(t *testing.T) {
		action := &Action{
			Jobs: []JobCall{
				{
					Name: "print",
					Args: []interface{}{
						map[string]interface{}{
							"type":     "cast",
							"castType": "number",
							"expression": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "p",
								"field":  "invalidNumber",
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
		assert.Error(t, err, "Invalid cast should produce an error")
		assert.Contains(t, err.Error(), "cannot cast", "Error should mention cast failure")
	})
	t.Run("invalid cast type", func(t *testing.T) {
		action := &Action{
			Jobs: []JobCall{
				{
					Name: "print",
					Args: []interface{}{
						map[string]interface{}{
							"type":     "cast",
							"castType": "integer", // Type non supporté
							"expression": map[string]interface{}{
								"type":  "number",
								"value": 123.0,
							},
						},
					},
				},
			},
		}
		token := &Token{
			Facts: []*Fact{product},
		}
		err := env.Network.ActionExecutor.ExecuteAction(action, token)
		assert.Error(t, err, "Invalid cast type should produce an error")
		assert.Contains(t, err.Error(), "type de cast non supporté", "Error should mention unsupported cast type")
	})
	t.Run("missing cast type", func(t *testing.T) {
		action := &Action{
			Jobs: []JobCall{
				{
					Name: "print",
					Args: []interface{}{
						map[string]interface{}{
							"type": "cast",
							// Pas de castType
							"expression": map[string]interface{}{
								"type":  "number",
								"value": 123.0,
							},
						},
					},
				},
			},
		}
		token := &Token{
			Facts: []*Fact{product},
		}
		err := env.Network.ActionExecutor.ExecuteAction(action, token)
		assert.Error(t, err, "Missing cast type should produce an error")
		assert.Contains(t, err.Error(), "type de cast manquant", "Error should mention missing cast type")
	})
	t.Run("missing expression", func(t *testing.T) {
		action := &Action{
			Jobs: []JobCall{
				{
					Name: "print",
					Args: []interface{}{
						map[string]interface{}{
							"type":     "cast",
							"castType": "string",
							// Pas d'expression
						},
					},
				},
			},
		}
		token := &Token{
			Facts: []*Fact{product},
		}
		err := env.Network.ActionExecutor.ExecuteAction(action, token)
		assert.Error(t, err, "Missing expression should produce an error")
		assert.Contains(t, err.Error(), "expression à caster manquante", "Error should mention missing expression")
	})
}

// TestCastMultipleArgumentsInAction teste le cast avec plusieurs arguments
func TestCastMultipleArgumentsInAction(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	order := &Fact{
		ID:   "o1",
		Type: "Order",
		Fields: map[string]interface{}{
			"orderId":  "ORD-123",
			"totalStr": "250.50",
			"quantity": 5.0,
		},
	}
	var capturedArgs []interface{}
	processHandler := &captureActionHandler{
		name:         "processOrder",
		capturedArgs: &capturedArgs,
	}
	err := env.Network.ActionExecutor.RegisterAction(processHandler)
	require.NoError(t, err)
	action := &Action{
		Jobs: []JobCall{
			{
				Name: "processOrder",
				Args: []interface{}{
					// Argument 1: orderId (pas de cast)
					map[string]interface{}{
						"type":   "fieldAccess",
						"object": "o",
						"field":  "orderId",
					},
					// Argument 2: total converti en number
					map[string]interface{}{
						"type":     "cast",
						"castType": "number",
						"expression": map[string]interface{}{
							"type":   "fieldAccess",
							"object": "o",
							"field":  "totalStr",
						},
					},
					// Argument 3: quantity converti en string
					map[string]interface{}{
						"type":     "cast",
						"castType": "string",
						"expression": map[string]interface{}{
							"type":   "fieldAccess",
							"object": "o",
							"field":  "quantity",
						},
					},
				},
			},
		},
	}
	token := &Token{
		Facts:    []*Fact{order},
		Bindings: NewBindingChainWith("o", order),
	}
	err = env.Network.ActionExecutor.ExecuteAction(action, token)
	require.NoError(t, err, "Action with multiple cast arguments should execute successfully")
	require.Len(t, capturedArgs, 3, "Should have 3 arguments")
	// Vérifier l'argument 1 (string, pas de cast)
	arg1, ok := capturedArgs[0].(string)
	require.True(t, ok, "First argument should be string")
	assert.Equal(t, "ORD-123", arg1)
	// Vérifier l'argument 2 (cast vers number)
	arg2, ok := capturedArgs[1].(float64)
	require.True(t, ok, "Second argument should be float64")
	assert.Equal(t, 250.50, arg2)
	// Vérifier l'argument 3 (cast vers string)
	arg3, ok := capturedArgs[2].(string)
	require.True(t, ok, "Third argument should be string")
	assert.Equal(t, "5", arg3)
}
