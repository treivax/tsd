// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"bytes"
	"strings"
	"testing"
)

// TestPrintActionIntegration_SimpleRule teste l'action print dans une règle simple
func TestPrintActionIntegration_SimpleRule(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION PRINT ACTION - RÈGLE SIMPLE")
	t.Log("===============================================")
	// Créer un réseau RETE
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Définir un type Person
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
	// Capturer la sortie de l'action print
	var output bytes.Buffer
	printAction := NewPrintAction(&output)
	// Remplacer l'action print par défaut pour capturer la sortie
	executor := network.ActionExecutor
	if executor == nil {
		t.Fatal("❌ ActionExecutor n'est pas initialisé")
	}
	executor.GetRegistry().Register(printAction)
	// Créer un fait
	fact := &Fact{
		ID:   "person_1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   "1",
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
	// Créer une action print
	action := &Action{
		Type: "action",
		Jobs: []JobCall{
			{
				Type: "jobCall",
				Name: "print",
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
	// Exécuter l'action
	err := executor.ExecuteAction(action, token)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'exécution de l'action: %v", err)
	}
	// Vérifier que l'action print a été exécutée
	result := strings.TrimSpace(output.String())
	if result != "Alice" {
		t.Errorf("❌ Sortie incorrecte: attendu 'Alice', reçu '%s'", result)
	}
	t.Log("✅ Test d'intégration règle simple réussi")
}

// TestPrintActionIntegration_MultipleJobs teste l'action print avec plusieurs jobs
func TestPrintActionIntegration_MultipleJobs(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION PRINT ACTION - PLUSIEURS JOBS")
	t.Log("=================================================")
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Définir un type Person
	personType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
			{Name: "status", Type: "string"},
		},
	}
	network.Types = append(network.Types, personType)
	// Capturer la sortie
	var output bytes.Buffer
	printAction := NewPrintAction(&output)
	network.ActionExecutor.GetRegistry().Register(printAction)
	// Créer un fait
	fact := &Fact{
		ID:   "person_1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":     "1",
			"name":   "Bob",
			"status": "active",
		},
	}
	token := &Token{
		ID:       "token1",
		Facts:    []*Fact{fact},
		Bindings: NewBindingChainWith("p", fact),
	}
	// Action avec plusieurs jobs print
	action := &Action{
		Type: "action",
		Jobs: []JobCall{
			{
				Type: "jobCall",
				Name: "print",
				Args: []interface{}{
					map[string]interface{}{
						"type":  "string",
						"value": "Active user:",
					},
				},
			},
			{
				Type: "jobCall",
				Name: "print",
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
		t.Fatalf("❌ Erreur lors de l'exécution: %v", err)
	}
	// Vérifier la sortie
	result := output.String()
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) != 2 {
		t.Errorf("❌ Devrait avoir 2 lignes de sortie, reçu %d", len(lines))
	}
	if !strings.Contains(result, "Active user:") {
		t.Error("❌ La sortie devrait contenir 'Active user:'")
	}
	if !strings.Contains(result, "Bob") {
		t.Error("❌ La sortie devrait contenir 'Bob'")
	}
	t.Log("✅ Test d'intégration plusieurs jobs réussi")
}

// TestPrintActionIntegration_WithNumbers teste l'action print avec des nombres
func TestPrintActionIntegration_WithNumbers(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION PRINT ACTION - AVEC NOMBRES")
	t.Log("================================================")
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	productType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Product",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
			{Name: "price", Type: "number"},
		},
	}
	network.Types = append(network.Types, productType)
	var output bytes.Buffer
	printAction := NewPrintAction(&output)
	network.ActionExecutor.GetRegistry().Register(printAction)
	// Créer un produit
	fact := &Fact{
		ID:   "product_1",
		Type: "Product",
		Fields: map[string]interface{}{
			"id":    "1",
			"name":  "Laptop",
			"price": 999.99,
		},
	}
	token := &Token{
		ID:       "token1",
		Facts:    []*Fact{fact},
		Bindings: NewBindingChainWith("prod", fact),
	}
	// Règle qui affiche le prix
	action := &Action{
		Type: "action",
		Jobs: []JobCall{
			{
				Type: "jobCall",
				Name: "print",
				Args: []interface{}{
					map[string]interface{}{
						"type":   "fieldAccess",
						"object": "prod",
						"field":  "price",
					},
				},
			},
		},
	}
	err := network.ActionExecutor.ExecuteAction(action, token)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'exécution: %v", err)
	}
	// Vérifier la sortie
	result := strings.TrimSpace(output.String())
	if result != "999.99" {
		t.Errorf("❌ Sortie incorrecte: attendu '999.99', reçu '%s'", result)
	}
	t.Log("✅ Test d'intégration avec nombres réussi")
}

// TestPrintActionIntegration_UndefinedAction teste le comportement avec une action non définie
func TestPrintActionIntegration_UndefinedAction(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION - ACTION NON DÉFINIE")
	t.Log("=========================================")
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
		ID:   "person_1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   "1",
			"name": "Charlie",
		},
	}
	token := &Token{
		ID:       "token1",
		Facts:    []*Fact{fact},
		Bindings: NewBindingChainWith("p", fact),
	}
	// Action non définie
	action := &Action{
		Type: "action",
		Jobs: []JobCall{
			{
				Type: "jobCall",
				Name: "undefined_custom_action",
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
	// Ne devrait pas causer d'erreur (l'action est juste loguée)
	err := network.ActionExecutor.ExecuteAction(action, token)
	if err != nil {
		t.Errorf("❌ Une action non définie ne devrait pas causer d'erreur: %v", err)
	}
	t.Log("✅ Test action non définie réussi (loguée uniquement)")
}

// TestPrintActionIntegration_MixedActions teste un mélange d'actions définies et non définies
func TestPrintActionIntegration_MixedActions(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION - ACTIONS MIXTES")
	t.Log("====================================")
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	eventType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Event",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "message", Type: "string"},
		},
	}
	network.Types = append(network.Types, eventType)
	var output bytes.Buffer
	printAction := NewPrintAction(&output)
	network.ActionExecutor.GetRegistry().Register(printAction)
	fact := &Fact{
		ID:   "event_1",
		Type: "Event",
		Fields: map[string]interface{}{
			"id":      "1",
			"message": "Test message",
		},
	}
	token := &Token{
		ID:       "token1",
		Facts:    []*Fact{fact},
		Bindings: NewBindingChainWith("e", fact),
	}
	// Action avec print (définie) et custom_action (non définie)
	action := &Action{
		Type: "action",
		Jobs: []JobCall{
			{
				Type: "jobCall",
				Name: "print",
				Args: []interface{}{
					map[string]interface{}{
						"type":   "fieldAccess",
						"object": "e",
						"field":  "message",
					},
				},
			},
			{
				Type: "jobCall",
				Name: "custom_undefined_action",
				Args: []interface{}{
					map[string]interface{}{
						"type":  "string",
						"value": "test",
					},
				},
			},
		},
	}
	err := network.ActionExecutor.ExecuteAction(action, token)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'exécution: %v", err)
	}
	// Vérifier que print a été exécuté
	result := strings.TrimSpace(output.String())
	if result != "Test message" {
		t.Errorf("❌ L'action print devrait avoir affiché 'Test message', reçu '%s'", result)
	}
	t.Log("✅ Test actions mixtes réussi")
}

// TestPrintActionIntegration_WithFact teste l'action print avec un fait complet
func TestPrintActionIntegration_WithFact(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION PRINT ACTION - AVEC FAIT COMPLET")
	t.Log("====================================================")
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
	var output bytes.Buffer
	printAction := NewPrintAction(&output)
	network.ActionExecutor.GetRegistry().Register(printAction)
	fact := &Fact{
		ID:   "person_1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   "1",
			"name": "David",
			"age":  35.0,
		},
	}
	token := &Token{
		ID:       "token1",
		Facts:    []*Fact{fact},
		Bindings: NewBindingChainWith("p", fact),
	}
	// Imprimer le fait complet (variable)
	action := &Action{
		Type: "action",
		Jobs: []JobCall{
			{
				Type: "jobCall",
				Name: "print",
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
		t.Fatalf("❌ Erreur lors de l'exécution: %v", err)
	}
	// Vérifier que la sortie contient des informations sur le fait
	result := strings.TrimSpace(output.String())
	if !strings.Contains(result, "Person") || !strings.Contains(result, "person_1") {
		t.Errorf("❌ La sortie devrait contenir le type et l'ID du fait: %s", result)
	}
	t.Log("✅ Test avec fait complet réussi")
}
