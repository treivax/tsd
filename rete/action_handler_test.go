// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"bytes"
	"strings"
	"testing"
)
// TestActionRegistry_Basic teste les fonctionnalités de base du registry
func TestActionRegistry_Basic(t *testing.T) {
	t.Log("🧪 TEST ACTION REGISTRY - FONCTIONNALITÉS DE BASE")
	t.Log("===================================================")
	registry := NewActionRegistry()
	// Test: Registry vide au départ
	if registry.Count() != 0 {
		t.Errorf("❌ Registry devrait être vide, contient %d handlers", registry.Count())
	}
	// Test: Enregistrer une action print
	printAction := NewPrintAction(nil)
	err := registry.Register(printAction)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'enregistrement: %v", err)
	}
	// Test: Vérifier que l'action est enregistrée
	if registry.Count() != 1 {
		t.Errorf("❌ Registry devrait contenir 1 handler, contient %d", registry.Count())
	}
	if !registry.Has(ActionNamePrint) {
		t.Error("❌ L'action print devrait être enregistrée")
	}
	// Test: Récupérer l'action
	handler := registry.Get(ActionNamePrint)
	if handler == nil {
		t.Fatal("❌ L'action print devrait être récupérable")
	}
	if handler.GetName() != ActionNamePrint {
		t.Errorf("❌ Nom de l'action incorrect: attendu '%s', reçu '%s'", ActionNamePrint, handler.GetName())
	}
	t.Log("✅ Tests de base réussis")
}
// TestActionRegistry_Unregister teste la désinscription d'actions
func TestActionRegistry_Unregister(t *testing.T) {
	t.Log("🧪 TEST ACTION REGISTRY - DÉSINSCRIPTION")
	t.Log("=========================================")
	registry := NewActionRegistry()
	printAction := NewPrintAction(nil)
	registry.Register(printAction)
	// Vérifier que l'action est présente
	if !registry.Has(ActionNamePrint) {
		t.Fatal("❌ L'action devrait être enregistrée")
	}
	// Désinscrire l'action
	registry.Unregister(ActionNamePrint)
	// Vérifier que l'action n'est plus présente
	if registry.Has(ActionNamePrint) {
		t.Error("❌ L'action ne devrait plus être enregistrée")
	}
	if registry.Count() != 0 {
		t.Errorf("❌ Registry devrait être vide, contient %d handlers", registry.Count())
	}
	t.Log("✅ Test de désinscription réussi")
}
// TestActionRegistry_Multiple teste l'enregistrement multiple
func TestActionRegistry_Multiple(t *testing.T) {
	t.Log("🧪 TEST ACTION REGISTRY - ENREGISTREMENT MULTIPLE")
	t.Log("=================================================")
	registry := NewActionRegistry()
	// Créer plusieurs actions de test
	printAction := NewPrintAction(nil)
	mockAction := &MockActionHandler{name: "mock_action"}
	handlers := []ActionHandler{printAction, mockAction}
	// Enregistrer toutes les actions
	err := registry.RegisterMultiple(handlers)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'enregistrement multiple: %v", err)
	}
	// Vérifier que toutes les actions sont enregistrées
	if registry.Count() != 2 {
		t.Errorf("❌ Registry devrait contenir 2 handlers, contient %d", registry.Count())
	}
	names := registry.GetRegisteredNames()
	if len(names) != 2 {
		t.Errorf("❌ Devrait avoir 2 noms, reçu %d", len(names))
	}
	t.Log("✅ Test d'enregistrement multiple réussi")
}
// TestActionRegistry_Clear teste le nettoyage du registry
func TestActionRegistry_Clear(t *testing.T) {
	t.Log("🧪 TEST ACTION REGISTRY - NETTOYAGE")
	t.Log("====================================")
	registry := NewActionRegistry()
	printAction := NewPrintAction(nil)
	registry.Register(printAction)
	// Vérifier que l'action est présente
	if registry.Count() != 1 {
		t.Fatal("❌ Registry devrait contenir 1 handler")
	}
	// Nettoyer le registry
	registry.Clear()
	// Vérifier que le registry est vide
	if registry.Count() != 0 {
		t.Errorf("❌ Registry devrait être vide après Clear(), contient %d handlers", registry.Count())
	}
	t.Log("✅ Test de nettoyage réussi")
}
// TestPrintAction_StringArgument teste l'action print avec une chaîne
func TestPrintAction_StringArgument(t *testing.T) {
	t.Log("🧪 TEST PRINT ACTION - ARGUMENT STRING")
	t.Log("======================================")
	// Créer un buffer pour capturer la sortie
	var output bytes.Buffer
	printAction := NewPrintAction(&output)
	// Créer un contexte d'exécution minimal
	ctx := NewExecutionContext(nil, nil)
	// Exécuter l'action avec une chaîne
	testString := "Hello, World!"
	args := []interface{}{testString}
	err := printAction.Execute(args, ctx)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'exécution: %v", err)
	}
	// Vérifier la sortie
	result := strings.TrimSpace(output.String())
	if result != testString {
		t.Errorf("❌ Sortie incorrecte: attendu '%s', reçu '%s'", testString, result)
	}
	t.Log("✅ Test avec argument string réussi")
}
// TestPrintAction_NumberArgument teste l'action print avec un nombre
func TestPrintAction_NumberArgument(t *testing.T) {
	t.Log("🧪 TEST PRINT ACTION - ARGUMENT NUMBER")
	t.Log("======================================")
	var output bytes.Buffer
	printAction := NewPrintAction(&output)
	ctx := NewExecutionContext(nil, nil)
	// Test avec un float64
	args := []interface{}{42.5}
	err := printAction.Execute(args, ctx)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'exécution: %v", err)
	}
	result := strings.TrimSpace(output.String())
	if result != "42.5" {
		t.Errorf("❌ Sortie incorrecte: attendu '42.5', reçu '%s'", result)
	}
	t.Log("✅ Test avec argument number réussi")
}
// TestPrintAction_BooleanArgument teste l'action print avec un booléen
func TestPrintAction_BooleanArgument(t *testing.T) {
	t.Log("🧪 TEST PRINT ACTION - ARGUMENT BOOLEAN")
	t.Log("=======================================")
	var output bytes.Buffer
	printAction := NewPrintAction(&output)
	ctx := NewExecutionContext(nil, nil)
	// Test avec un booléen
	args := []interface{}{true}
	err := printAction.Execute(args, ctx)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'exécution: %v", err)
	}
	result := strings.TrimSpace(output.String())
	if result != "true" {
		t.Errorf("❌ Sortie incorrecte: attendu 'true', reçu '%s'", result)
	}
	t.Log("✅ Test avec argument boolean réussi")
}
// TestPrintAction_FactArgument teste l'action print avec un fait
func TestPrintAction_FactArgument(t *testing.T) {
	t.Log("🧪 TEST PRINT ACTION - ARGUMENT FACT")
	t.Log("====================================")
	var output bytes.Buffer
	printAction := NewPrintAction(&output)
	ctx := NewExecutionContext(nil, nil)
	// Créer un fait de test
	fact := &Fact{
		ID:   "person_1",
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Alice",
			"age":  30.0,
		},
	}
	args := []interface{}{fact}
	err := printAction.Execute(args, ctx)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'exécution: %v", err)
	}
	result := strings.TrimSpace(output.String())
	if !strings.Contains(result, "Person") || !strings.Contains(result, "person_1") {
		t.Errorf("❌ La sortie devrait contenir le type et l'ID du fait: %s", result)
	}
	t.Log("✅ Test avec argument fact réussi")
}
// TestPrintAction_NoArguments teste l'action print sans arguments
func TestPrintAction_NoArguments(t *testing.T) {
	t.Log("🧪 TEST PRINT ACTION - SANS ARGUMENTS")
	t.Log("=====================================")
	var output bytes.Buffer
	printAction := NewPrintAction(&output)
	ctx := NewExecutionContext(nil, nil)
	// Exécuter sans arguments
	args := []interface{}{}
	err := printAction.Execute(args, ctx)
	// Devrait retourner une erreur
	if err == nil {
		t.Error("❌ Devrait retourner une erreur quand il n'y a pas d'arguments")
	}
	t.Log("✅ Test sans arguments réussi (erreur attendue)")
}
// TestPrintAction_Validate teste la validation de l'action print
func TestPrintAction_Validate(t *testing.T) {
	t.Log("🧪 TEST PRINT ACTION - VALIDATION")
	t.Log("=================================")
	printAction := NewPrintAction(nil)
	// Test: validation réussie avec un argument valide
	args := []interface{}{"test"}
	err := printAction.Validate(args)
	if err != nil {
		t.Errorf("❌ La validation devrait réussir: %v", err)
	}
	// Test: validation échouée sans arguments
	emptyArgs := []interface{}{}
	err = printAction.Validate(emptyArgs)
	if err == nil {
		t.Error("❌ La validation devrait échouer sans arguments")
	}
	t.Log("✅ Test de validation réussi")
}
// TestActionExecutor_WithRegistry teste l'intégration du registry dans l'executor
func TestActionExecutor_WithRegistry(t *testing.T) {
	t.Log("🧪 TEST ACTION EXECUTOR - AVEC REGISTRY")
	t.Log("=======================================")
	// Créer un réseau et un executor
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	executor := NewActionExecutor(network, nil)
	// Vérifier que le registry est initialisé
	registry := executor.GetRegistry()
	if registry == nil {
		t.Fatal("❌ Le registry devrait être initialisé")
	}
	// Vérifier que l'action print est enregistrée par défaut
	if !registry.Has(ActionNamePrint) {
		t.Error("❌ L'action print devrait être enregistrée par défaut")
	}
	t.Log("✅ Test d'intégration réussi")
}
// TestActionExecutor_CustomAction teste l'enregistrement d'une action personnalisée
func TestActionExecutor_CustomAction(t *testing.T) {
	t.Log("🧪 TEST ACTION EXECUTOR - ACTION PERSONNALISÉE")
	t.Log("==============================================")
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	executor := NewActionExecutor(network, nil)
	// Créer et enregistrer une action personnalisée
	customAction := &MockActionHandler{name: "custom"}
	err := executor.RegisterAction(customAction)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'enregistrement: %v", err)
	}
	// Vérifier que l'action est enregistrée
	if !executor.GetRegistry().Has("custom") {
		t.Error("❌ L'action personnalisée devrait être enregistrée")
	}
	t.Log("✅ Test d'action personnalisée réussi")
}
// TestActionExecutor_UndefinedAction teste le comportement avec une action non définie
func TestActionExecutor_UndefinedAction(t *testing.T) {
	t.Log("🧪 TEST ACTION EXECUTOR - ACTION NON DÉFINIE")
	t.Log("============================================")
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	executor := NewActionExecutor(network, nil)
	// Créer un job avec une action non définie
	job := JobCall{
		Type: "jobCall",
		Name: "undefined_action",
		Args: []interface{}{"test"},
	}
	ctx := NewExecutionContext(nil, network)
	// Exécuter le job
	err := executor.executeJob(job, ctx, 0)
	// Ne devrait pas retourner d'erreur (juste logger)
	if err != nil {
		t.Errorf("❌ Une action non définie ne devrait pas causer d'erreur: %v", err)
	}
	t.Log("✅ Test d'action non définie réussi")
}
// MockActionHandler est un handler de test
type MockActionHandler struct {
	name           string
	executeCalled  bool
	validateCalled bool
	lastArgs       []interface{}
}
func (m *MockActionHandler) Execute(args []interface{}, ctx *ExecutionContext) error {
	m.executeCalled = true
	m.lastArgs = args
	return nil
}
func (m *MockActionHandler) GetName() string {
	return m.name
}
func (m *MockActionHandler) Validate(args []interface{}) error {
	m.validateCalled = true
	return nil
}
// TestActionRegistry_NilHandler teste l'enregistrement d'un handler nil
func TestActionRegistry_NilHandler(t *testing.T) {
	t.Log("🧪 TEST ACTION REGISTRY - HANDLER NIL")
	t.Log("=====================================")
	registry := NewActionRegistry()
	err := registry.Register(nil)
	if err == nil {
		t.Error("❌ L'enregistrement d'un handler nil devrait retourner une erreur")
	}
	t.Log("✅ Test handler nil réussi")
}
// TestActionRegistry_EmptyName teste l'enregistrement d'un handler avec nom vide
func TestActionRegistry_EmptyName(t *testing.T) {
	t.Log("🧪 TEST ACTION REGISTRY - NOM VIDE")
	t.Log("==================================")
	registry := NewActionRegistry()
	emptyNameHandler := &MockActionHandler{name: ""}
	err := registry.Register(emptyNameHandler)
	if err == nil {
		t.Error("❌ L'enregistrement d'un handler avec nom vide devrait retourner une erreur")
	}
	t.Log("✅ Test nom vide réussi")
}
// TestPrintAction_SetOutput teste le changement de sortie
func TestPrintAction_SetOutput(t *testing.T) {
	t.Log("🧪 TEST PRINT ACTION - CHANGEMENT DE SORTIE")
	t.Log("===========================================")
	// Créer avec une première sortie
	var output1 bytes.Buffer
	printAction := NewPrintAction(&output1)
	ctx := NewExecutionContext(nil, nil)
	// Exécuter
	args := []interface{}{"first"}
	printAction.Execute(args, ctx)
	// Changer la sortie
	var output2 bytes.Buffer
	printAction.SetOutput(&output2)
	// Exécuter à nouveau
	args2 := []interface{}{"second"}
	printAction.Execute(args2, ctx)
	// Vérifier les sorties
	result1 := strings.TrimSpace(output1.String())
	result2 := strings.TrimSpace(output2.String())
	if result1 != "first" {
		t.Errorf("❌ Première sortie incorrecte: '%s'", result1)
	}
	if result2 != "second" {
		t.Errorf("❌ Deuxième sortie incorrecte: '%s'", result2)
	}
	t.Log("✅ Test changement de sortie réussi")
}
// TestPrintAction_IntegerTypes teste l'action print avec différents types d'entiers
func TestPrintAction_IntegerTypes(t *testing.T) {
	t.Log("🧪 TEST PRINT ACTION - TYPES D'ENTIERS")
	t.Log("======================================")
	var output bytes.Buffer
	printAction := NewPrintAction(&output)
	ctx := NewExecutionContext(nil, nil)
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"int", int(42), "42"},
		{"int64", int64(100), "100"},
		{"float64", float64(3.14), "3.14"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.Reset()
			args := []interface{}{tt.value}
			err := printAction.Execute(args, ctx)
			if err != nil {
				t.Errorf("❌ Erreur lors de l'exécution: %v", err)
			}
			result := strings.TrimSpace(output.String())
			if result != tt.expected {
				t.Errorf("❌ Sortie incorrecte pour %s: attendu '%s', reçu '%s'",
					tt.name, tt.expected, result)
			}
		})
	}
	t.Log("✅ Test types d'entiers réussi")
}
// TestActionRegistry_GetAll teste la récupération de tous les handlers
func TestActionRegistry_GetAll(t *testing.T) {
	t.Log("🧪 TEST ACTION REGISTRY - GET ALL")
	t.Log("=================================")
	registry := NewActionRegistry()
	// Enregistrer plusieurs handlers
	printAction := NewPrintAction(nil)
	mockAction := &MockActionHandler{name: "mock"}
	registry.Register(printAction)
	registry.Register(mockAction)
	// Récupérer tous les handlers
	allHandlers := registry.GetAll()
	if len(allHandlers) != 2 {
		t.Errorf("❌ Devrait avoir 2 handlers, reçu %d", len(allHandlers))
	}
	if _, exists := allHandlers[ActionNamePrint]; !exists {
		t.Error("❌ L'action print devrait être dans la liste")
	}
	if _, exists := allHandlers["mock"]; !exists {
		t.Error("❌ L'action mock devrait être dans la liste")
	}
	t.Log("✅ Test GetAll réussi")
}
// TestPrintAction_NilFact teste l'action print avec un fait nil
func TestPrintAction_NilFact(t *testing.T) {
	t.Log("🧪 TEST PRINT ACTION - FAIT NIL")
	t.Log("================================")
	var output bytes.Buffer
	printAction := NewPrintAction(&output)
	ctx := NewExecutionContext(nil, nil)
	// Créer un fait nil
	var nilFact *Fact = nil
	args := []interface{}{nilFact}
	err := printAction.Execute(args, ctx)
	if err == nil {
		t.Error("❌ Devrait retourner une erreur avec un fait nil")
	}
	t.Log("✅ Test fait nil réussi (erreur attendue)")
}