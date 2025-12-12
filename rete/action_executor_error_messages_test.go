// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestActionExecutor_ErrorMessages_VariableList teste que les messages d'erreur
// listent bien les variables disponibles pour aider au debugging.
func TestActionExecutor_ErrorMessages_VariableList(t *testing.T) {
	t.Parallel()
	t.Log("🧪 TEST - Messages d'erreur avec liste des variables disponibles")

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Créer des faits
	userFact := &Fact{
		ID:   "u1",
		Type: "User",
		Fields: map[string]interface{}{
			"id":   "u1",
			"name": "Alice",
		},
	}
	orderFact := &Fact{
		ID:   "o1",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":      "o1",
			"user_id": "u1",
		},
	}

	// Token avec 2 bindings
	token := &Token{
		ID:    "t1",
		Facts: []*Fact{userFact, orderFact},
		Bindings: NewBindingChain().
			Add("user", userFact).
			Add("order", orderFact),
	}

	tests := []struct {
		name              string
		action            *Action
		expectedInError   []string
		unexpectedInError []string
	}{
		{
			name: "Variable inexistante dans variable reference",
			action: &Action{
				Type: "action",
				Jobs: []JobCall{
					{
						Type: "jobCall",
						Name: "process",
						Args: []interface{}{
							map[string]interface{}{
								"type": "variable",
								"name": "product", // N'existe pas
							},
						},
					},
				},
			},
			expectedInError: []string{
				"product",
				"non trouvée",
				"Variables disponibles",
				"user",
				"order",
			},
			unexpectedInError: []string{},
		},
		{
			name: "Variable inexistante dans fieldAccess",
			action: &Action{
				Type: "action",
				Jobs: []JobCall{
					{
						Type: "jobCall",
						Name: "process",
						Args: []interface{}{
							map[string]interface{}{
								"type":   "fieldAccess",
								"object": "task", // N'existe pas
								"field":  "status",
							},
						},
					},
				},
			},
			expectedInError: []string{
				"task",
				"non trouvée",
				"Variables disponibles",
				"user",
				"order",
			},
			unexpectedInError: []string{},
		},
		{
			name: "Champ inexistant - ne liste PAS les variables (erreur différente)",
			action: &Action{
				Type: "action",
				Jobs: []JobCall{
					{
						Type: "jobCall",
						Name: "process",
						Args: []interface{}{
							map[string]interface{}{
								"type":   "fieldAccess",
								"object": "user",  // Existe
								"field":  "email", // N'existe pas
							},
						},
					},
				},
			},
			expectedInError: []string{
				"email",
				"non trouvé",
			},
			unexpectedInError: []string{
				"Variables disponibles", // Ne devrait pas lister les variables pour un champ manquant
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := env.Network.ActionExecutor.ExecuteAction(tt.action, token)
			require.Error(t, err, "Should return error")

			errorMsg := err.Error()
			t.Logf("Message d'erreur reçu:\n%s", errorMsg)

			// Vérifier que tous les strings attendus sont présents
			for _, expected := range tt.expectedInError {
				if !strings.Contains(errorMsg, expected) {
					t.Errorf("❌ Message d'erreur devrait contenir '%s'\nMessage: %s", expected, errorMsg)
				}
			}

			// Vérifier qu'aucun string non désiré n'est présent
			for _, unexpected := range tt.unexpectedInError {
				if strings.Contains(errorMsg, unexpected) {
					t.Errorf("❌ Message d'erreur ne devrait PAS contenir '%s'\nMessage: %s", unexpected, errorMsg)
				}
			}
		})
	}

	t.Log("✅ Messages d'erreur affichent correctement les variables disponibles")
}

// TestExecutionContext_ResolveVariable_WithBindingChain teste la résolution
// de variables via BindingChain dans ExecutionContext.
func TestExecutionContext_ResolveVariable_WithBindingChain(t *testing.T) {
	t.Parallel()
	t.Log("🧪 TEST - ExecutionContext résout les variables via BindingChain")

	// Setup
	userFact := &Fact{
		ID:   "u1",
		Type: "User",
		Fields: map[string]interface{}{
			"id":   "u1",
			"name": "Alice",
		},
	}
	orderFact := &Fact{
		ID:   "o1",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":      "o1",
			"user_id": "u1",
		},
	}
	productFact := &Fact{
		ID:   "p1",
		Type: "Product",
		Fields: map[string]interface{}{
			"id":   "p1",
			"name": "Laptop",
		},
	}

	// Token avec 3 bindings
	token := &Token{
		ID:    "t1",
		Facts: []*Fact{userFact, orderFact, productFact},
		Bindings: NewBindingChain().
			Add("user", userFact).
			Add("order", orderFact).
			Add("product", productFact),
	}

	ctx := NewExecutionContext(token, nil)

	// Test 1: Variables existantes
	tests := []struct {
		name     string
		varName  string
		expected *Fact
	}{
		{"user", "user", userFact},
		{"order", "order", orderFact},
		{"product", "product", productFact},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ctx.GetVariable(tt.varName)
			require.NotNil(t, result, "Variable %s devrait exister", tt.varName)
			require.Equal(t, tt.expected.ID, result.ID, "Fact ID devrait correspondre")
			require.Equal(t, tt.expected.Type, result.Type, "Fact Type devrait correspondre")
		})
	}

	// Test 2: Variable inexistante
	t.Run("Variable inexistante", func(t *testing.T) {
		result := ctx.GetVariable("nonexistent")
		require.Nil(t, result, "Variable inexistante devrait retourner nil")
	})

	t.Log("✅ ExecutionContext résout correctement les variables via BindingChain")
}

// TestTerminalNode_ExecuteAction_AllVariablesAvailable teste l'exécution
// d'une action avec toutes les variables disponibles via BindingChain.
func TestTerminalNode_ExecuteAction_AllVariablesAvailable(t *testing.T) {
	t.Parallel()
	t.Log("🧪 TEST - TerminalNode exécute action avec toutes variables disponibles")

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Définir les types
	env.Network.Types = []TypeDefinition{
		{
			Type: "typeDefinition",
			Name: "User",
			Fields: []Field{
				{Name: "id", Type: "string"},
				{Name: "name", Type: "string"},
			},
		},
		{
			Type: "typeDefinition",
			Name: "Order",
			Fields: []Field{
				{Name: "id", Type: "string"},
				{Name: "user_id", Type: "string"},
				{Name: "total", Type: "number"},
			},
		},
		{
			Type: "typeDefinition",
			Name: "Product",
			Fields: []Field{
				{Name: "id", Type: "string"},
				{Name: "name", Type: "string"},
				{Name: "price", Type: "number"},
			},
		},
	}

	// Créer des faits
	userFact := &Fact{
		ID:   "u1",
		Type: "User",
		Fields: map[string]interface{}{
			"id":   "u1",
			"name": "Alice",
		},
	}
	orderFact := &Fact{
		ID:   "o1",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":      "o1",
			"user_id": "u1",
			"total":   100.0,
		},
	}
	productFact := &Fact{
		ID:   "p1",
		Type: "Product",
		Fields: map[string]interface{}{
			"id":    "p1",
			"name":  "Laptop",
			"price": 999.99,
		},
	}

	// Token avec 3 bindings
	token := &Token{
		ID:    "t_final",
		Facts: []*Fact{userFact, orderFact, productFact},
		Bindings: NewBindingChain().
			Add("user", userFact).
			Add("order", orderFact).
			Add("product", productFact),
	}

	// Action utilisant les 3 variables
	action := &Action{
		Type: "action",
		Jobs: []JobCall{
			{
				Type: "jobCall",
				Name: "log_purchase",
				Args: []interface{}{
					// Argument 1: user.name
					map[string]interface{}{
						"type":   "fieldAccess",
						"object": "user",
						"field":  "name",
					},
					// Argument 2: product.name
					map[string]interface{}{
						"type":   "fieldAccess",
						"object": "product",
						"field":  "name",
					},
					// Argument 3: order.total
					map[string]interface{}{
						"type":   "fieldAccess",
						"object": "order",
						"field":  "total",
					},
				},
			},
		},
	}

	// Exécuter l'action
	err := env.Network.ActionExecutor.ExecuteAction(action, token)
	require.NoError(t, err, "Action avec 3 variables devrait s'exécuter sans erreur")

	t.Log("✅ Action exécutée avec succès avec toutes les variables disponibles")
}
