// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"log"
	"os"
	"testing"
)

// TestActionExecutor_BasicExecution teste l'exécution basique d'une action
func TestActionExecutor_BasicExecution(t *testing.T) {
	t.Log("🧪 TEST EXÉCUTION BASIQUE D'ACTION")
	t.Log("===================================")

	// Créer un réseau RETE avec types
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

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
	network.Types = append(network.Types, personType)

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
	executor := network.ActionExecutor
	err := executor.ExecuteAction(action, token)

	if err != nil {
		t.Fatalf("❌ Erreur exécution action: %v", err)
	}

	t.Log("✅ Action exécutée avec succès")
}

// TestActionExecutor_VariableArgument teste l'utilisation d'une variable comme argument
func TestActionExecutor_VariableArgument(t *testing.T) {
	t.Log("🧪 TEST ARGUMENT VARIABLE")
	t.Log("=========================")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
		},
	}
	network.Types = append(network.Types, personType)

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

	err := network.ActionExecutor.ExecuteAction(action, token)

	if err != nil {
		t.Fatalf("❌ Erreur exécution action: %v", err)
	}

	t.Log("✅ Action avec variable exécutée avec succès")
}

// TestActionExecutor_FieldAccessArgument teste l'accès à un attribut
func TestActionExecutor_FieldAccessArgument(t *testing.T) {
	t.Log("🧪 TEST ARGUMENT ACCÈS ATTRIBUT")
	t.Log("================================")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
			{Name: "age", Type: "number"},
		},
	}
	network.Types = append(network.Types, personType)

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

	err := network.ActionExecutor.ExecuteAction(action, token)

	if err != nil {
		t.Fatalf("❌ Erreur exécution action: %v", err)
	}

	t.Log("✅ Action avec accès attribut exécutée avec succès")
}

// TestActionExecutor_MultipleArguments teste plusieurs arguments de types différents
func TestActionExecutor_MultipleArguments(t *testing.T) {
	t.Log("🧪 TEST ARGUMENTS MULTIPLES")
	t.Log("============================")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
			{Name: "age", Type: "number"},
		},
	}
	network.Types = append(network.Types, personType)

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

	err := network.ActionExecutor.ExecuteAction(action, token)

	if err != nil {
		t.Fatalf("❌ Erreur exécution action: %v", err)
	}

	t.Log("✅ Action avec arguments multiples exécutée avec succès")
}

// TestActionExecutor_ArithmeticExpression teste les expressions arithmétiques
func TestActionExecutor_ArithmeticExpression(t *testing.T) {
	t.Log("🧪 TEST EXPRESSION ARITHMÉTIQUE")
	t.Log("================================")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "salary", Type: "number"},
		},
	}
	network.Types = append(network.Types, personType)

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

	err := network.ActionExecutor.ExecuteAction(action, token)

	if err != nil {
		t.Fatalf("❌ Erreur exécution action: %v", err)
	}

	t.Log("✅ Expression arithmétique évaluée avec succès")
}

// TestActionExecutor_MultipleJobs teste l'exécution de plusieurs jobs
func TestActionExecutor_MultipleJobs(t *testing.T) {
	t.Log("🧪 TEST JOBS MULTIPLES")
	t.Log("======================")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
		},
	}
	network.Types = append(network.Types, personType)

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

	err := network.ActionExecutor.ExecuteAction(action, token)

	if err != nil {
		t.Fatalf("❌ Erreur exécution action: %v", err)
	}

	t.Log("✅ Trois jobs exécutés en séquence avec succès")
}

// TestActionExecutor_ValidationErrors teste les erreurs de validation
func TestActionExecutor_ValidationErrors(t *testing.T) {
	t.Log("🧪 TEST ERREURS DE VALIDATION")
	t.Log("==============================")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
		},
	}
	network.Types = append(network.Types, personType)

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

	testCases := []struct {
		name        string
		action      *Action
		shouldError bool
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
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := network.ActionExecutor.ExecuteAction(tc.action, token)

			if tc.shouldError && err == nil {
				t.Errorf("❌ Attendait une erreur pour '%s', mais aucune erreur", tc.name)
			} else if !tc.shouldError && err != nil {
				t.Errorf("❌ N'attendait pas d'erreur pour '%s', reçu: %v", tc.name, err)
			} else if tc.shouldError && err != nil {
				t.Logf("✅ Erreur correctement détectée: %v", err)
			}
		})
	}
}

// TestActionExecutor_Logging teste le logging des actions
func TestActionExecutor_Logging(t *testing.T) {
	t.Log("🧪 TEST LOGGING DES ACTIONS")
	t.Log("============================")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
		},
	}
	network.Types = append(network.Types, personType)

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
	network.ActionExecutor.SetLogging(true)
	err := network.ActionExecutor.ExecuteAction(action, token)
	if err != nil {
		t.Fatalf("❌ Erreur avec logging activé: %v", err)
	}

	// Tester avec logging désactivé
	network.ActionExecutor.SetLogging(false)
	err = network.ActionExecutor.ExecuteAction(action, token)
	if err != nil {
		t.Fatalf("❌ Erreur avec logging désactivé: %v", err)
	}

	t.Log("✅ Logging fonctionne correctement")
}

// TestActionExecutor_CustomLogger teste l'utilisation d'un logger personnalisé
func TestActionExecutor_CustomLogger(t *testing.T) {
	t.Log("🧪 TEST LOGGER PERSONNALISÉ")
	t.Log("============================")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Créer un logger personnalisé
	customLogger := log.New(os.Stdout, "[CUSTOM] ", log.LstdFlags)
	network.ActionExecutor = NewActionExecutor(network, customLogger)

	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
		},
	}
	network.Types = append(network.Types, personType)

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

	err := network.ActionExecutor.ExecuteAction(action, token)
	if err != nil {
		t.Fatalf("❌ Erreur avec logger personnalisé: %v", err)
	}

	t.Log("✅ Logger personnalisé fonctionne correctement")
}
