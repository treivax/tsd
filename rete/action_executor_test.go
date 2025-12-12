// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestActionExecutor_evaluateComparison(t *testing.T) {
	executor := NewActionExecutor(nil, nil)
	tests := []struct {
		name      string
		left      interface{}
		operator  string
		right     interface{}
		expected  interface{}
		wantError bool
	}{
		// Equality tests
		{
			name:     "equal integers",
			left:     5,
			operator: "==",
			right:    5,
			expected: true,
		},
		{
			name:     "not equal integers",
			left:     5,
			operator: "==",
			right:    10,
			expected: false,
		},
		{
			name:     "equal strings",
			left:     "hello",
			operator: "==",
			right:    "hello",
			expected: true,
		},
		{
			name:     "not equal strings",
			left:     "hello",
			operator: "==",
			right:    "world",
			expected: false,
		},
		{
			name:     "equal float and int",
			left:     5.0,
			operator: "==",
			right:    5,
			expected: true,
		},
		{
			name:     "equal int64 and int",
			left:     int64(5),
			operator: "==",
			right:    5,
			expected: true,
		},
		// Inequality tests
		{
			name:     "not equal operator - true",
			left:     5,
			operator: "!=",
			right:    10,
			expected: true,
		},
		{
			name:     "not equal operator - false",
			left:     5,
			operator: "!=",
			right:    5,
			expected: false,
		},
		// Less than tests
		{
			name:     "less than - true",
			left:     5,
			operator: "<",
			right:    10,
			expected: true,
		},
		{
			name:     "less than - false",
			left:     10,
			operator: "<",
			right:    5,
			expected: false,
		},
		{
			name:     "less than float",
			left:     5.5,
			operator: "<",
			right:    10.2,
			expected: true,
		},
		// Less than or equal tests
		{
			name:     "less than or equal - true (less)",
			left:     5,
			operator: "<=",
			right:    10,
			expected: true,
		},
		{
			name:     "less than or equal - true (equal)",
			left:     5,
			operator: "<=",
			right:    5,
			expected: true,
		},
		{
			name:     "less than or equal - false",
			left:     10,
			operator: "<=",
			right:    5,
			expected: false,
		},
		// Greater than tests
		{
			name:     "greater than - true",
			left:     10,
			operator: ">",
			right:    5,
			expected: true,
		},
		{
			name:     "greater than - false",
			left:     5,
			operator: ">",
			right:    10,
			expected: false,
		},
		{
			name:     "greater than float",
			left:     10.5,
			operator: ">",
			right:    5.2,
			expected: true,
		},
		// Greater than or equal tests
		{
			name:     "greater than or equal - true (greater)",
			left:     10,
			operator: ">=",
			right:    5,
			expected: true,
		},
		{
			name:     "greater than or equal - true (equal)",
			left:     5,
			operator: ">=",
			right:    5,
			expected: true,
		},
		{
			name:     "greater than or equal - false",
			left:     5,
			operator: ">=",
			right:    10,
			expected: false,
		},
		// Error cases
		{
			name:      "numeric comparison with string left",
			left:      "hello",
			operator:  "<",
			right:     5,
			wantError: true,
		},
		{
			name:      "numeric comparison with string right",
			left:      5,
			operator:  ">",
			right:     "hello",
			wantError: true,
		},
		{
			name:      "unknown operator",
			left:      5,
			operator:  "???",
			right:     10,
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := executor.evaluateComparison(tt.left, tt.operator, tt.right)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
func TestActionExecutor_areEqual(t *testing.T) {
	executor := NewActionExecutor(nil, nil)
	tests := []struct {
		name     string
		left     interface{}
		right    interface{}
		expected bool
	}{
		// Numeric equality
		{
			name:     "equal integers",
			left:     5,
			right:    5,
			expected: true,
		},
		{
			name:     "not equal integers",
			left:     5,
			right:    10,
			expected: false,
		},
		{
			name:     "equal floats",
			left:     5.5,
			right:    5.5,
			expected: true,
		},
		{
			name:     "equal int and float",
			left:     5,
			right:    5.0,
			expected: true,
		},
		{
			name:     "equal float and int",
			left:     5.0,
			right:    5,
			expected: true,
		},
		{
			name:     "equal int64 and int",
			left:     int64(5),
			right:    5,
			expected: true,
		},
		{
			name:     "equal int and int64",
			left:     5,
			right:    int64(5),
			expected: true,
		},
		{
			name:     "equal int32 and int",
			left:     int32(5),
			right:    5,
			expected: true,
		},
		{
			name:     "not equal different numeric values",
			left:     5.5,
			right:    10.2,
			expected: false,
		},
		// String equality
		{
			name:     "equal strings",
			left:     "hello",
			right:    "hello",
			expected: true,
		},
		{
			name:     "not equal strings",
			left:     "hello",
			right:    "world",
			expected: false,
		},
		// Boolean equality
		{
			name:     "equal booleans true",
			left:     true,
			right:    true,
			expected: true,
		},
		{
			name:     "equal booleans false",
			left:     false,
			right:    false,
			expected: true,
		},
		{
			name:     "not equal booleans",
			left:     true,
			right:    false,
			expected: false,
		},
		// Mixed type comparisons (non-numeric)
		{
			name:     "string vs int",
			left:     "5",
			right:    5,
			expected: false,
		},
		{
			name:     "bool vs int",
			left:     true,
			right:    1,
			expected: false,
		},
		// Nil comparisons
		{
			name:     "nil equals nil",
			left:     nil,
			right:    nil,
			expected: true,
		},
		{
			name:     "nil not equals int",
			left:     nil,
			right:    5,
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.areEqual(tt.left, tt.right)
			assert.Equal(t, tt.expected, result)
		})
	}
}

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
		ID:       "token1",
		Facts:    []*Fact{fact},
		Bindings: NewBindingChainWith("p", fact),
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
		ID:       "token1",
		Facts:    []*Fact{fact},
		Bindings: NewBindingChainWith("p", fact),
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
		ID:       "token1",
		Facts:    []*Fact{fact},
		Bindings: NewBindingChainWith("p", fact),
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
		ID:       "token1",
		Facts:    []*Fact{fact},
		Bindings: NewBindingChainWith("p", fact),
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
		ID:       "token1",
		Facts:    []*Fact{fact},
		Bindings: NewBindingChainWith("p", fact),
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
		ID:       "token1",
		Facts:    []*Fact{fact},
		Bindings: NewBindingChainWith("p", fact),
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
				ID:       "token1",
				Facts:    []*Fact{fact},
				Bindings: NewBindingChainWith("p", fact),
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
		ID:       "token1",
		Facts:    []*Fact{fact},
		Bindings: NewBindingChainWith("p", fact),
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
		ID:       "token1",
		Facts:    []*Fact{fact},
		Bindings: NewBindingChainWith("p", fact),
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
func TestActionExecutor_RegisterDefaultActions(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	// Get registry before registering defaults
	registry := env.Network.ActionExecutor.GetRegistry()
	require.NotNil(t, registry, "Registry should not be nil")
	// Register default actions
	env.Network.ActionExecutor.RegisterDefaultActions()
	// Verify print action is registered
	handler := registry.Get("print")
	require.NotNil(t, handler, "Print action should be registered")
	// Test that we can execute the print action
	fact := &Fact{
		ID:   "f1",
		Type: "TestFact",
		Fields: map[string]interface{}{
			"value": "test",
		},
	}
	token := &Token{
		ID:       "t1",
		Facts:    []*Fact{fact},
		Bindings: NewBindingChainWith("f", fact),
	}
	action := &Action{
		Type: "action",
		Job: &JobCall{
			Type: "jobCall",
			Name: "print",
			Args: []interface{}{
				map[string]interface{}{
					"type":  "string",
					"value": "test message",
				},
			},
		},
	}
	err := env.Network.ActionExecutor.ExecuteAction(action, token)
	require.NoError(t, err, "Print action should execute successfully")
}
func TestActionExecutor_EvaluateBinaryOperation_Coverage(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	tests := []struct {
		name      string
		argMap    map[string]interface{}
		expectErr bool
	}{
		{
			name: "addition with numbers",
			argMap: map[string]interface{}{
				"operator": "+",
				"left":     float64(5),
				"right":    float64(3),
			},
		},
		{
			name: "subtraction",
			argMap: map[string]interface{}{
				"operator": "-",
				"left":     float64(10),
				"right":    float64(4),
			},
		},
		{
			name: "multiplication",
			argMap: map[string]interface{}{
				"operator": "*",
				"left":     float64(6),
				"right":    float64(7),
			},
		},
		{
			name: "division",
			argMap: map[string]interface{}{
				"operator": "/",
				"left":     float64(15),
				"right":    float64(3),
			},
		},
		{
			name: "modulo",
			argMap: map[string]interface{}{
				"operator": "%",
				"left":     float64(17),
				"right":    float64(5),
			},
		},
		{
			name: "comparison equals",
			argMap: map[string]interface{}{
				"operator": "==",
				"left":     float64(5),
				"right":    float64(5),
			},
		},
		{
			name: "comparison less than",
			argMap: map[string]interface{}{
				"operator": "<",
				"left":     float64(3),
				"right":    float64(5),
			},
		},
		{
			name: "invalid left operand type",
			argMap: map[string]interface{}{
				"operator": "+",
				"left":     "not a number",
				"right":    float64(5),
			},
			expectErr: true,
		},
		{
			name: "invalid right operand type",
			argMap: map[string]interface{}{
				"operator": "+",
				"left":     float64(5),
				"right":    "not a number",
			},
			expectErr: true,
		},
		{
			name: "unknown operator",
			argMap: map[string]interface{}{
				"operator": "??",
				"left":     float64(5),
				"right":    float64(3),
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewExecutionContext(&Token{}, env.Network)
			result, err := env.Network.ActionExecutor.evaluateBinaryOperation(tt.argMap, ctx)
			if tt.expectErr {
				require.Error(t, err, "Expected error for %s", tt.name)
			} else {
				require.NoError(t, err, "Should not error for %s", tt.name)
				require.NotNil(t, result, "Result should not be nil")
			}
		})
	}
}
func TestActionExecutor_ValidateFieldType_Coverage(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "name", Type: "string"},
			{Name: "age", Type: "number"},
			{Name: "active", Type: "bool"},
		},
	}
	env.Network.Types = append(env.Network.Types, personType)
	tests := []struct {
		name        string
		value       interface{}
		expectedTyp string
		expectErr   bool
	}{
		{
			name:        "string value",
			value:       "test",
			expectedTyp: "string",
			expectErr:   false,
		},
		{
			name:        "int value as number",
			value:       42,
			expectedTyp: "number",
			expectErr:   false,
		},
		{
			name:        "int64 value as number",
			value:       int64(100),
			expectedTyp: "number",
			expectErr:   false,
		},
		{
			name:        "float64 value as number",
			value:       3.14,
			expectedTyp: "number",
			expectErr:   false,
		},
		{
			name:        "bool value",
			value:       true,
			expectedTyp: "bool",
			expectErr:   false,
		},
		{
			name:        "string value expecting number - type mismatch",
			value:       "not a number",
			expectedTyp: "number",
			expectErr:   true,
		},
		{
			name:        "number value expecting string - type mismatch",
			value:       42,
			expectedTyp: "string",
			expectErr:   true,
		},
		{
			name:        "bool value expecting string - type mismatch",
			value:       true,
			expectedTyp: "string",
			expectErr:   true,
		},
		{
			name:        "unsupported type",
			value:       []interface{}{1, 2, 3},
			expectedTyp: "string",
			expectErr:   true,
		},
		{
			name:        "nil value",
			value:       nil,
			expectedTyp: "string",
			expectErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := env.Network.ActionExecutor.validateFieldType(tt.expectedTyp, tt.value)
			if tt.expectErr {
				require.Error(t, err, "Expected error for %s", tt.name)
			} else {
				require.NoError(t, err, "Should not error for %s", tt.name)
			}
		})
	}
}
func TestActionExecutor_EvaluateArgument_EdgeCases(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "name", Type: "string"},
			{Name: "age", Type: "number"},
		},
	}
	env.Network.Types = append(env.Network.Types, personType)
	fact := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Alice",
			"age":  float64(30),
		},
	}
	token := &Token{
		ID:       "t1",
		Facts:    []*Fact{fact},
		Bindings: NewBindingChainWith("p", fact),
	}
	ctx := NewExecutionContext(token, env.Network)
	tests := []struct {
		name      string
		arg       interface{}
		expectErr bool
	}{
		{
			name: "string literal",
			arg: map[string]interface{}{
				"type":  "string",
				"value": "hello",
			},
			expectErr: false,
		},
		{
			name: "number literal",
			arg: map[string]interface{}{
				"type":  "number",
				"value": 42,
			},
			expectErr: false,
		},
		{
			name: "boolean literal",
			arg: map[string]interface{}{
				"type":  "bool",
				"value": true,
			},
			expectErr: false,
		},
		{
			name: "variable reference",
			arg: map[string]interface{}{
				"type": "variable",
				"name": "p",
			},
			expectErr: false,
		},
		{
			name: "field access",
			arg: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "name",
			},
			expectErr: false,
		},
		{
			name: "binary operation",
			arg: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "+",
				"left": map[string]interface{}{
					"type":  "number",
					"value": 10,
				},
				"right": map[string]interface{}{
					"type":  "number",
					"value": 5,
				},
			},
			expectErr: false,
		},
		{
			name: "invalid variable",
			arg: map[string]interface{}{
				"type": "variable",
				"name": "unknown",
			},
			expectErr: true,
		},
		{
			name: "invalid field access",
			arg: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "unknownField",
			},
			expectErr: true,
		},
		{
			name: "unknown type - returns as-is",
			arg: map[string]interface{}{
				"type": "unknownType",
			},
			expectErr: false,
		},
		{
			name:      "plain string (not a map)",
			arg:       "plain value",
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := env.Network.ActionExecutor.evaluateArgument(tt.arg, ctx)
			if tt.expectErr {
				require.Error(t, err, "Expected error for %s", tt.name)
			} else {
				require.NoError(t, err, "Should not error for %s", tt.name)
				require.NotNil(t, result, "Result should not be nil")
			}
		})
	}
}
func TestActionExecutor_RegisterAction(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	// Create a custom action handler
	customHandler := &mockActionHandler{
		name: "customAction",
	}
	// Register the custom action
	err := env.Network.ActionExecutor.RegisterAction(customHandler)
	require.NoError(t, err, "Should register custom action")
	// Verify it's registered
	handler := env.Network.ActionExecutor.GetRegistry().Get("customAction")
	require.NotNil(t, handler, "Custom action should be in registry")
	require.Equal(t, "customAction", handler.GetName(), "Handler name should match")
}

// mockActionHandler is a simple mock for testing
type mockActionHandler struct {
	name string
}

func (m *mockActionHandler) GetName() string {
	return m.name
}
func (m *mockActionHandler) Execute(args []interface{}, ctx *ExecutionContext) error {
	return nil
}
func (m *mockActionHandler) Validate(args []interface{}) error {
	return nil
}
func TestActionExecutor_EvaluateFactCreation_Coverage(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
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
	token := &Token{
		ID:    "t1",
		Facts: []*Fact{},
	}
	ctx := NewExecutionContext(token, env.Network)
	tests := []struct {
		name        string
		argMap      map[string]interface{}
		expectError bool
	}{
		{
			name: "valid fact creation with all fields",
			argMap: map[string]interface{}{
				"typeName": "Person",
				"fields": map[string]interface{}{
					"id": map[string]interface{}{
						"type":  "string",
						"value": "P001",
					},
					"name": map[string]interface{}{
						"type":  "string",
						"value": "Alice",
					},
					"age": map[string]interface{}{
						"type":  "number",
						"value": 30,
					},
				},
			},
			expectError: false,
		},
		{
			name: "missing typeName",
			argMap: map[string]interface{}{
				"fields": map[string]interface{}{
					"id": "test",
				},
			},
			expectError: true,
		},
		{
			name: "undefined type",
			argMap: map[string]interface{}{
				"typeName": "UnknownType",
				"fields": map[string]interface{}{
					"id": "test",
				},
			},
			expectError: true,
		},
		{
			name: "missing fields",
			argMap: map[string]interface{}{
				"typeName": "Person",
			},
			expectError: true,
		},
		{
			name: "error in field evaluation",
			argMap: map[string]interface{}{
				"typeName": "Person",
				"fields": map[string]interface{}{
					"id": map[string]interface{}{
						"type": "variable",
						"name": "unknownVar",
					},
				},
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := env.Network.ActionExecutor.evaluateFactCreation(tt.argMap, ctx)
			if tt.expectError {
				require.Error(t, err, "Expected error for %s", tt.name)
			} else {
				require.NoError(t, err, "Should not error for %s", tt.name)
				require.NotNil(t, result, "Result should not be nil")
				fact, ok := result.(*Fact)
				require.True(t, ok, "Result should be a Fact")
				require.Equal(t, "Person", fact.Type)
			}
		})
	}
}
func TestActionExecutor_EvaluateFactModification_Coverage(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
			{Name: "age", Type: "number"},
			{Name: "active", Type: "bool"},
		},
	}
	env.Network.Types = append(env.Network.Types, personType)
	fact := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":     "P001",
			"name":   "Alice",
			"age":    float64(30),
			"active": true,
		},
	}
	token := &Token{
		ID:       "t1",
		Facts:    []*Fact{fact},
		Bindings: NewBindingChainWith("p", fact),
	}
	ctx := NewExecutionContext(token, env.Network)
	tests := []struct {
		name        string
		argMap      map[string]interface{}
		expectError bool
	}{
		{
			name: "valid field modification",
			argMap: map[string]interface{}{
				"variable": "p",
				"field":    "name",
				"value": map[string]interface{}{
					"type":  "string",
					"value": "Bob",
				},
			},
			expectError: false,
		},
		{
			name: "modify number field",
			argMap: map[string]interface{}{
				"variable": "p",
				"field":    "age",
				"value": map[string]interface{}{
					"type":  "number",
					"value": 31,
				},
			},
			expectError: false,
		},
		{
			name: "modify bool field",
			argMap: map[string]interface{}{
				"variable": "p",
				"field":    "active",
				"value": map[string]interface{}{
					"type":  "bool",
					"value": false,
				},
			},
			expectError: false,
		},
		{
			name: "missing variable",
			argMap: map[string]interface{}{
				"field": "name",
				"value": "test",
			},
			expectError: true,
		},
		{
			name: "missing field",
			argMap: map[string]interface{}{
				"variable": "p",
				"value":    "test",
			},
			expectError: true,
		},
		{
			name: "missing value",
			argMap: map[string]interface{}{
				"variable": "p",
				"field":    "name",
			},
			expectError: true,
		},
		{
			name: "variable not found",
			argMap: map[string]interface{}{
				"variable": "unknown",
				"field":    "name",
				"value":    "test",
			},
			expectError: true,
		},
		{
			name: "error evaluating value",
			argMap: map[string]interface{}{
				"variable": "p",
				"field":    "name",
				"value": map[string]interface{}{
					"type": "variable",
					"name": "unknownVar",
				},
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := env.Network.ActionExecutor.evaluateFactModification(tt.argMap, ctx)
			if tt.expectError {
				require.Error(t, err, "Expected error for %s", tt.name)
			} else {
				require.NoError(t, err, "Should not error for %s", tt.name)
				require.NotNil(t, result, "Result should not be nil")
				modifiedFact, ok := result.(*Fact)
				require.True(t, ok, "Result should be a Fact")
				require.Equal(t, "Person", modifiedFact.Type)
				require.Equal(t, "p1", modifiedFact.ID, "Should preserve fact ID")
			}
		})
	}
}
func TestActionExecutor_EvaluateArithmetic_ErrorPaths(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	emptyToken := &Token{
		ID:       "empty",
		Bindings: NewBindingChain(),
	}
	ctx := NewExecutionContext(emptyToken, env.Network)

	tests := []struct {
		name          string
		argMap        map[string]interface{}
		expectError   bool
		errorContains string
	}{
		{
			name: "missing operator",
			argMap: map[string]interface{}{
				"left":  10.0,
				"right": 5.0,
			},
			expectError:   true,
			errorContains: "opérateur manquant",
		},
		{
			name: "invalid left operand",
			argMap: map[string]interface{}{
				"operator": "+",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "nonexistent",
					"field":  "value",
				},
				"right": 5.0,
			},
			expectError:   true,
			errorContains: "erreur évaluation left",
		},
		{
			name: "invalid right operand",
			argMap: map[string]interface{}{
				"operator": "+",
				"left":     10.0,
				"right": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "nonexistent",
					"field":  "value",
				},
			},
			expectError:   true,
			errorContains: "erreur évaluation right",
		},
		{
			name: "valid addition",
			argMap: map[string]interface{}{
				"operator": "+",
				"left":     10.0,
				"right":    5.0,
			},
			expectError: false,
		},
		{
			name: "valid subtraction",
			argMap: map[string]interface{}{
				"operator": "-",
				"left":     10.0,
				"right":    3.0,
			},
			expectError: false,
		},
		{
			name: "valid multiplication",
			argMap: map[string]interface{}{
				"operator": "*",
				"left":     4.0,
				"right":    5.0,
			},
			expectError: false,
		},
		{
			name: "valid division",
			argMap: map[string]interface{}{
				"operator": "/",
				"left":     20.0,
				"right":    4.0,
			},
			expectError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := env.Network.ActionExecutor.evaluateArithmetic(tt.argMap, ctx)
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
func TestActionExecutor_ValidateFactFields_Coverage(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
			{Name: "age", Type: "number"},
			{Name: "active", Type: "bool"},
		},
	}
	tests := []struct {
		name        string
		fields      map[string]interface{}
		expectError bool
	}{
		{
			name: "all valid fields",
			fields: map[string]interface{}{
				"id":     "P001",
				"name":   "Alice",
				"age":    float64(30),
				"active": true,
			},
			expectError: false,
		},
		{
			name: "missing required field - error",
			fields: map[string]interface{}{
				"id":   "P001",
				"name": "Alice",
				"age":  float64(30),
			},
			expectError: true,
		},
		{
			name: "extra field not in type - error",
			fields: map[string]interface{}{
				"id":     "P001",
				"name":   "Alice",
				"age":    float64(30),
				"active": true,
				"extra":  "ignored",
			},
			expectError: true,
		},
		{
			name: "invalid field type - string expected",
			fields: map[string]interface{}{
				"id":   123,
				"name": "Alice",
			},
			expectError: true,
		},
		{
			name: "invalid field type - number expected",
			fields: map[string]interface{}{
				"id":   "P001",
				"name": "Alice",
				"age":  "not a number",
			},
			expectError: true,
		},
		{
			name: "invalid field type - bool expected",
			fields: map[string]interface{}{
				"id":     "P001",
				"name":   "Alice",
				"age":    float64(30),
				"active": "not a bool",
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := env.Network.ActionExecutor.validateFactFields(&personType, tt.fields)
			if tt.expectError {
				require.Error(t, err, "Expected error for %s", tt.name)
			} else {
				require.NoError(t, err, "Should not error for %s", tt.name)
			}
		})
	}
}
