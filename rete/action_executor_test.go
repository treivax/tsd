// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestActionExecutor_BasicExecution teste l'exécution basique d'une action
func TestActionExecutor_BasicExecution(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t, WithLogLevel(LogLevelInfo))
	defer env.Cleanup()

	// Définir un type
	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
			{Name: "age", Type: "number"},
		},
	}
	env.Network.Types = append(env.Network.Types, personType)

	// Créer un fait
	fact := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   "p1",
			"name": "Alice",
			"age":  25.0,
		},
	}

	// Créer un token avec bindings
	token := &Token{
		ID:    "token1",
		Facts: []*Fact{fact},
		Bindings: map[string]*Fact{
			"p": fact,
		},
	}

	// Créer une action simple avec valeur littérale
	action := &Action{
		Type: "action",
		Jobs: []JobCall{
			{
				Type: "jobCall",
				Name: "log",
				Args: []interface{}{
					map[string]interface{}{
						"type":  "string",
						"value": "Test message",
					},
				},
			},
		},
	}

	// Exécuter l'action
	executor := env.Network.ActionExecutor
	err := executor.ExecuteAction(action, token)

	require.NoError(t, err, "Action execution should succeed")
	env.AssertNoErrors(t)
}

// TestActionExecutor_VariableArgument teste l'utilisation d'une variable comme argument
func TestActionExecutor_VariableArgument(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t, WithLogLevel(LogLevelInfo))
	defer env.Cleanup()

	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
		},
	}
	env.Network.Types = append(env.Network.Types, personType)

	fact := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   "p1",
			"name": "Bob",
		},
	}

	token := &Token{
		ID:    "token1",
		Facts: []*Fact{fact},
		Bindings: map[string]*Fact{
			"p": fact,
		},
	}

	// Action avec variable complète
	action := &Action{
		Type: "action",
		Jobs: []JobCall{
			{
				Type: "jobCall",
				Name: "process",
				Args: []interface{}{
					map[string]interface{}{
						"type": "variable",
						"name": "p",
					},
				},
			},
		},
	}

	err := env.Network.ActionExecutor.ExecuteAction(action, token)

	require.NoError(t, err, "Action with variable should execute successfully")
	env.AssertNoErrors(t)
}

// TestActionExecutor_FieldAccessArgument teste l'accès à un attribut
func TestActionExecutor_FieldAccessArgument(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t, WithLogLevel(LogLevelInfo))
	defer env.Cleanup()

	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
			{Name: "age", Type: "number"},
		},
	}
	env.Network.Types = append(env.Network.Types, personType)

	fact := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   "p1",
			"name": "Charlie",
			"age":  30.0,
		},
	}

	token := &Token{
		ID:    "token1",
		Facts: []*Fact{fact},
		Bindings: map[string]*Fact{
			"p": fact,
		},
	}

	// Action avec accès à un champ
	action := &Action{
		Type: "action",
		Jobs: []JobCall{
			{
				Type: "jobCall",
				Name: "notify",
				Args: []interface{}{
					map[string]interface{}{
						"type":   "fieldAccess",
						"object": "p",
						"field":  "name",
					},
				},
			},
		},
	}

	err := env.Network.ActionExecutor.ExecuteAction(action, token)

	require.NoError(t, err, "Action with field access should execute successfully")
	env.AssertNoErrors(t)
}

// TestActionExecutor_MultipleArguments teste plusieurs arguments de types différents
func TestActionExecutor_MultipleArguments(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t, WithLogLevel(LogLevelInfo))
	defer env.Cleanup()

	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
			{Name: "age", Type: "number"},
		},
	}
	env.Network.Types = append(env.Network.Types, personType)

	fact := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   "p1",
			"name": "Diana",
			"age":  28.0,
		},
	}

	token := &Token{
		ID:    "token1",
		Facts: []*Fact{fact},
		Bindings: map[string]*Fact{
			"p": fact,
		},
	}

	// Action avec plusieurs types d'arguments
	action := &Action{
		Type: "action",
		Jobs: []JobCall{
			{
				Type: "jobCall",
				Name: "process",
				Args: []interface{}{
					// Argument 1: variable complète
					map[string]interface{}{
						"type": "variable",
						"name": "p",
					},
					// Argument 2: accès attribut
					map[string]interface{}{
						"type":   "fieldAccess",
						"object": "p",
						"field":  "name",
					},
					// Argument 3: valeur littérale
					map[string]interface{}{
						"type":  "string",
						"value": "processed",
					},
				},
			},
		},
	}

	err := env.Network.ActionExecutor.ExecuteAction(action, token)

	require.NoError(t, err, "Action with multiple arguments should execute successfully")
	env.AssertNoErrors(t)
}

// TestActionExecutor_ArithmeticExpression teste les expressions arithmétiques
func TestActionExecutor_ArithmeticExpression(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t, WithLogLevel(LogLevelInfo))
	defer env.Cleanup()

	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "salary", Type: "number"},
		},
	}
	env.Network.Types = append(env.Network.Types, personType)

	fact := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":     "p1",
			"salary": 50000.0,
		},
	}

	token := &Token{
		ID:    "token1",
		Facts: []*Fact{fact},
		Bindings: map[string]*Fact{
			"p": fact,
		},
	}

	// Action avec expression arithmétique
	action := &Action{
		Type: "action",
		Jobs: []JobCall{
			{
				Type: "jobCall",
				Name: "calculate_bonus",
				Args: []interface{}{
					map[string]interface{}{
						"type":     "arithmetic",
						"operator": "*",
						"left": map[string]interface{}{
							"type":   "fieldAccess",
							"object": "p",
							"field":  "salary",
						},
						"right": map[string]interface{}{
							"type":  "number",
							"value": 1.1,
						},
					},
				},
			},
		},
	}

	err := env.Network.ActionExecutor.ExecuteAction(action, token)

	require.NoError(t, err, "Arithmetic expression should be evaluated successfully")
	env.AssertNoErrors(t)
}

// TestActionExecutor_MultipleJobs teste l'exécution de plusieurs jobs
func TestActionExecutor_MultipleJobs(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t, WithLogLevel(LogLevelInfo))
	defer env.Cleanup()

	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
		},
	}
	env.Network.Types = append(env.Network.Types, personType)

	fact := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   "p1",
			"name": "Eve",
		},
	}

	token := &Token{
		ID:    "token1",
		Facts: []*Fact{fact},
		Bindings: map[string]*Fact{
			"p": fact,
		},
	}

	// Action avec trois jobs
	action := &Action{
		Type: "action",
		Jobs: []JobCall{
			{
				Type: "jobCall",
				Name: "log",
				Args: []interface{}{
					map[string]interface{}{"type": "string", "value": "Job 1"},
				},
			},
			{
				Type: "jobCall",
				Name: "notify",
				Args: []interface{}{
					map[string]interface{}{
						"type":   "fieldAccess",
						"object": "p",
						"field":  "name",
					},
				},
			},
			{
				Type: "jobCall",
				Name: "update",
				Args: []interface{}{
					map[string]interface{}{"type": "variable", "name": "p"},
				},
			},
		},
	}

	err := env.Network.ActionExecutor.ExecuteAction(action, token)

	require.NoError(t, err, "Multiple jobs should execute in sequence successfully")
	env.AssertNoErrors(t)
}

// TestActionExecutor_ValidationErrors teste les erreurs de validation
func TestActionExecutor_ValidationErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		action      *Action
		shouldError bool
		errorMsg    string
	}{
		{
			name: "Variable inexistante",
			action: &Action{
				Type: "action",
				Jobs: []JobCall{
					{
						Type: "jobCall",
						Name: "process",
						Args: []interface{}{
							map[string]interface{}{
								"type": "variable",
								"name": "unknown",
							},
						},
					},
				},
			},
			shouldError: true,
			errorMsg:    "variable 'unknown' non trouvée",
		},
		{
			name: "Champ inexistant",
			action: &Action{
				Type: "action",
				Jobs: []JobCall{
					{
						Type: "jobCall",
						Name: "process",
						Args: []interface{}{
							map[string]interface{}{
								"type":   "fieldAccess",
								"object": "p",
								"field":  "nonexistent",
							},
						},
					},
				},
			},
			shouldError: true,
			errorMsg:    "champ 'nonexistent' non trouvé",
		},
		{
			name: "Division par zéro",
			action: &Action{
				Type: "action",
				Jobs: []JobCall{
					{
						Type: "jobCall",
						Name: "calculate",
						Args: []interface{}{
							map[string]interface{}{
								"type":     "arithmetic",
								"operator": "/",
								"left": map[string]interface{}{
									"type":  "number",
									"value": 100.0,
								},
								"right": map[string]interface{}{
									"type":  "number",
									"value": 0.0,
								},
							},
						},
					},
				},
			},
			shouldError: true,
			errorMsg:    "division par zéro",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			env := NewTestEnvironment(t, WithLogLevel(LogLevelWarn))
			defer env.Cleanup()

			personType := TypeDefinition{
				Type: "typeDefinition",
				Name: "Person",
				Fields: []Field{
					{Name: "id", Type: "string"},
					{Name: "name", Type: "string"},
				},
			}
			env.Network.Types = append(env.Network.Types, personType)

			fact := &Fact{
				ID:   "p1",
				Type: "Person",
				Fields: map[string]interface{}{
					"id":   "p1",
					"name": "Frank",
				},
			}

			token := &Token{
				ID:    "token1",
				Facts: []*Fact{fact},
				Bindings: map[string]*Fact{
					"p": fact,
				},
			}

			err := env.Network.ActionExecutor.ExecuteAction(tc.action, token)

			if tc.shouldError {
				assert.Error(t, err, "Expected error for %s", tc.name)
				if err != nil {
					assert.Contains(t, strings.ToLower(err.Error()), strings.ToLower(tc.errorMsg),
						"Error message should contain expected text")
				}
			} else {
				assert.NoError(t, err, "Did not expect error for %s", tc.name)
			}
		})
	}
}

// TestActionExecutor_Logging teste le logging des actions
func TestActionExecutor_Logging(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t, WithLogLevel(LogLevelDebug))
	defer env.Cleanup()

	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
		},
	}
	env.Network.Types = append(env.Network.Types, personType)

	fact := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id": "p1",
		},
	}

	token := &Token{
		ID:    "token1",
		Facts: []*Fact{fact},
		Bindings: map[string]*Fact{
			"p": fact,
		},
	}

	action := &Action{
		Type: "action",
		Jobs: []JobCall{
			{
				Type: "jobCall",
				Name: "test_action",
				Args: []interface{}{
					map[string]interface{}{"type": "variable", "name": "p"},
				},
			},
		},
	}

	// Tester avec logging activé
	env.Network.ActionExecutor.SetLogging(true)
	err := env.Network.ActionExecutor.ExecuteAction(action, token)
	require.NoError(t, err, "Action with logging enabled should succeed")

	// Tester avec logging désactivé
	env.Network.ActionExecutor.SetLogging(false)
	err = env.Network.ActionExecutor.ExecuteAction(action, token)
	require.NoError(t, err, "Action with logging disabled should succeed")

	env.AssertNoErrors(t)
}

// TestActionExecutor_CustomLogger teste l'utilisation d'un logger personnalisé
func TestActionExecutor_CustomLogger(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t, WithLogLevel(LogLevelInfo))
	defer env.Cleanup()

	// Le ActionExecutor utilise déjà le logger de l'environment via le Network
	// Testons simplement que l'exécution fonctionne avec le logger configuré

	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
		},
	}
	env.Network.Types = append(env.Network.Types, personType)

	fact := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id": "p1",
		},
	}

	token := &Token{
		ID:    "token1",
		Facts: []*Fact{fact},
		Bindings: map[string]*Fact{
			"p": fact,
		},
	}

	action := &Action{
		Type: "action",
		Jobs: []JobCall{
			{
				Type: "jobCall",
				Name: "custom_test",
				Args: []interface{}{},
			},
		},
	}

	// Activer le logging pour capturer l'exécution
	env.Network.ActionExecutor.SetLogging(true)

	err := env.Network.ActionExecutor.ExecuteAction(action, token)
	require.NoError(t, err, "Action with custom logger should succeed")

	// Note: ActionExecutor logs to stdout via standard log package,
	// not to TestEnvironment buffer. The important thing is no errors occurred.
	env.AssertNoErrors(t)
}
