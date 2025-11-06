package nodes

import (
	"testing"
	"time"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// MockLogger implémente l'interface Logger pour les tests
type MockLogger struct{}

func (ml *MockLogger) Debug(msg string, fields map[string]interface{})            {}
func (ml *MockLogger) Info(msg string, fields map[string]interface{})             {}
func (ml *MockLogger) Warn(msg string, fields map[string]interface{})             {}
func (ml *MockLogger) Error(msg string, err error, fields map[string]interface{}) {}

func TestNotNode_ProcessNegation(t *testing.T) {
	logger := &MockLogger{}
	notNode := NewNotNode("not1", logger)

	// Définir une condition de négation
	notNode.SetNegationCondition("amount > 50")

	// Créer un token et un fait de test
	token := &domain.Token{
		ID: "token1",
		Facts: []*domain.Fact{
			{
				ID:        "fact1",
				Type:      "Transaction",
				Fields:    map[string]interface{}{"amount": 100.0},
				Timestamp: time.Now(),
			},
		},
	}

	fact := &domain.Fact{
		ID:        "fact2",
		Type:      "Account",
		Fields:    map[string]interface{}{"balance": 500.0},
		Timestamp: time.Now(),
	}

	// Test de la négation (avec l'implémentation basique qui retourne true,
	// la négation devrait retourner false)
	result := notNode.ProcessNegation(token, fact)
	if result {
		t.Errorf("Expected ProcessNegation to return false, got %v", result)
	}
}

func TestNotNode_ProcessLeftToken(t *testing.T) {
	logger := &MockLogger{}
	notNode := NewNotNode("not1", logger)

	// Créer un token de test
	token := &domain.Token{
		ID: "token1",
		Facts: []*domain.Fact{
			{
				ID:        "fact1",
				Type:      "Transaction",
				Fields:    map[string]interface{}{"amount": 100.0},
				Timestamp: time.Now(),
			},
		},
	}

	// Traiter le token (pas de faits de droite, donc devrait propager)
	err := notNode.ProcessLeftToken(token)
	if err != nil {
		t.Errorf("ProcessLeftToken failed: %v", err)
	}

	// Vérifier que le token est stocké dans la mémoire
	tokens := notNode.betaMemory.GetTokens()
	if len(tokens) != 1 {
		t.Errorf("Expected 1 token in memory, got %d", len(tokens))
	}

	if tokens[0].ID != "token1" {
		t.Errorf("Expected token ID 'token1', got '%s'", tokens[0].ID)
	}
}

func TestExistsNode_CheckExistence(t *testing.T) {
	logger := &MockLogger{}
	existsNode := NewExistsNode("exists1", logger)

	// Définir une variable d'existence
	variable := domain.TypedVariable{
		Name:     "account",
		DataType: "Account",
	}
	existsNode.SetExistenceCondition(variable, "balance > 0")

	// Créer un token de test
	token := &domain.Token{
		ID: "token1",
		Facts: []*domain.Fact{
			{
				ID:        "fact1",
				Type:      "Transaction",
				Fields:    map[string]interface{}{"amount": 100.0},
				Timestamp: time.Now(),
			},
		},
	}

	// Ajouter un fait de droite correspondant au type
	fact := &domain.Fact{
		ID:        "fact2",
		Type:      "Account",
		Fields:    map[string]interface{}{"balance": 500.0},
		Timestamp: time.Now(),
	}
	existsNode.betaMemory.StoreFact(fact)

	// Test de l'existence (devrait retourner true avec l'implémentation basique)
	result := existsNode.CheckExistence(token)
	if !result {
		t.Errorf("Expected CheckExistence to return true, got %v", result)
	}
}

func TestExistsNode_ProcessLeftToken(t *testing.T) {
	logger := &MockLogger{}
	existsNode := NewExistsNode("exists1", logger)

	// Définir une variable d'existence
	variable := domain.TypedVariable{
		Name:     "account",
		DataType: "Account",
	}
	existsNode.SetExistenceCondition(variable, "balance > 0")

	// Ajouter un fait de droite correspondant
	fact := &domain.Fact{
		ID:        "fact2",
		Type:      "Account",
		Fields:    map[string]interface{}{"balance": 500.0},
		Timestamp: time.Now(),
	}
	existsNode.betaMemory.StoreFact(fact)

	// Créer un token de test
	token := &domain.Token{
		ID: "token1",
		Facts: []*domain.Fact{
			{
				ID:        "fact1",
				Type:      "Transaction",
				Fields:    map[string]interface{}{"amount": 100.0},
				Timestamp: time.Now(),
			},
		},
	}

	// Traiter le token
	err := existsNode.ProcessLeftToken(token)
	if err != nil {
		t.Errorf("ProcessLeftToken failed: %v", err)
	}

	// Vérifier que le token est stocké dans la mémoire
	tokens := existsNode.betaMemory.GetTokens()
	if len(tokens) != 1 {
		t.Errorf("Expected 1 token in memory, got %d", len(tokens))
	}
}

func TestAccumulateNode_ComputeSum(t *testing.T) {
	logger := &MockLogger{}
	accumulator := domain.AccumulateFunction{
		FunctionType: "SUM",
		Field:        "amount",
	}
	accNode := NewAccumulateNode("acc1", accumulator, logger)

	// Créer des faits avec des montants
	facts := []*domain.Fact{
		{
			ID:        "fact1",
			Type:      "Transaction",
			Fields:    map[string]interface{}{"amount": 100.0},
			Timestamp: time.Now(),
		},
		{
			ID:        "fact2",
			Type:      "Transaction",
			Fields:    map[string]interface{}{"amount": 200.0},
			Timestamp: time.Now(),
		},
		{
			ID:        "fact3",
			Type:      "Transaction",
			Fields:    map[string]interface{}{"amount": 150.0},
			Timestamp: time.Now(),
		},
	}

	// Créer un token de test
	token := &domain.Token{
		ID:    "token1",
		Facts: []*domain.Fact{facts[0]},
	}

	// Calculer la somme
	result, err := accNode.ComputeAggregate(token, facts)
	if err != nil {
		t.Errorf("ComputeAggregate failed: %v", err)
	}

	expectedSum := 450.0
	if result != expectedSum {
		t.Errorf("Expected sum %v, got %v", expectedSum, result)
	}
}

func TestAccumulateNode_ComputeCount(t *testing.T) {
	logger := &MockLogger{}
	accumulator := domain.AccumulateFunction{
		FunctionType: "COUNT",
		Field:        "",
	}
	accNode := NewAccumulateNode("acc1", accumulator, logger)

	// Créer des faits
	facts := []*domain.Fact{
		{
			ID:        "fact1",
			Type:      "Transaction",
			Fields:    map[string]interface{}{"amount": 100.0},
			Timestamp: time.Now(),
		},
		{
			ID:        "fact2",
			Type:      "Transaction",
			Fields:    map[string]interface{}{"amount": 200.0},
			Timestamp: time.Now(),
		},
		{
			ID:        "fact3",
			Type:      "Transaction",
			Fields:    map[string]interface{}{"amount": 150.0},
			Timestamp: time.Now(),
		},
	}

	// Créer un token de test
	token := &domain.Token{
		ID:    "token1",
		Facts: []*domain.Fact{facts[0]},
	}

	// Calculer le compte
	result, err := accNode.ComputeAggregate(token, facts)
	if err != nil {
		t.Errorf("ComputeAggregate failed: %v", err)
	}

	expectedCount := 3
	if result != expectedCount {
		t.Errorf("Expected count %v, got %v", expectedCount, result)
	}
}

func TestAccumulateNode_ComputeAverage(t *testing.T) {
	logger := &MockLogger{}
	accumulator := domain.AccumulateFunction{
		FunctionType: "AVG",
		Field:        "amount",
	}
	accNode := NewAccumulateNode("acc1", accumulator, logger)

	// Créer des faits avec des montants
	facts := []*domain.Fact{
		{
			ID:        "fact1",
			Type:      "Transaction",
			Fields:    map[string]interface{}{"amount": 100.0},
			Timestamp: time.Now(),
		},
		{
			ID:        "fact2",
			Type:      "Transaction",
			Fields:    map[string]interface{}{"amount": 200.0},
			Timestamp: time.Now(),
		},
		{
			ID:        "fact3",
			Type:      "Transaction",
			Fields:    map[string]interface{}{"amount": 150.0},
			Timestamp: time.Now(),
		},
	}

	// Créer un token de test
	token := &domain.Token{
		ID:    "token1",
		Facts: []*domain.Fact{facts[0]},
	}

	// Calculer la moyenne
	result, err := accNode.ComputeAggregate(token, facts)
	if err != nil {
		t.Errorf("ComputeAggregate failed: %v", err)
	}

	expectedAvg := 150.0
	if result != expectedAvg {
		t.Errorf("Expected average %v, got %v", expectedAvg, result)
	}
}

func TestAccumulateNode_ProcessLeftToken(t *testing.T) {
	logger := &MockLogger{}
	accumulator := domain.AccumulateFunction{
		FunctionType: "SUM",
		Field:        "amount",
	}
	accNode := NewAccumulateNode("acc1", accumulator, logger)

	// Ajouter des faits de droite
	facts := []*domain.Fact{
		{
			ID:        "fact1",
			Type:      "Transaction",
			Fields:    map[string]interface{}{"amount": 100.0},
			Timestamp: time.Now(),
		},
		{
			ID:        "fact2",
			Type:      "Transaction",
			Fields:    map[string]interface{}{"amount": 200.0},
			Timestamp: time.Now(),
		},
	}

	for _, fact := range facts {
		accNode.betaMemory.StoreFact(fact)
	}

	// Créer un token de test
	token := &domain.Token{
		ID: "token1",
		Facts: []*domain.Fact{
			{
				ID:        "fact0",
				Type:      "Account",
				Fields:    map[string]interface{}{"id": "acc1"},
				Timestamp: time.Now(),
			},
		},
	}

	// Traiter le token
	err := accNode.ProcessLeftToken(token)
	if err != nil {
		t.Errorf("ProcessLeftToken failed: %v", err)
	}

	// Vérifier que le token est stocké dans la mémoire
	tokens := accNode.betaMemory.GetTokens()
	if len(tokens) != 1 {
		t.Errorf("Expected 1 token in memory, got %d", len(tokens))
	}

	// Vérifier que la valeur agrégée est stockée
	accNode.mu.RLock()
	result, exists := accNode.accumulatedValues[token.ID]
	accNode.mu.RUnlock()

	if !exists {
		t.Errorf("Expected accumulated value to be stored for token %s", token.ID)
	}

	expectedSum := 300.0
	if result != expectedSum {
		t.Errorf("Expected accumulated sum %v, got %v", expectedSum, result)
	}
}

func TestAccumulateNode_ProcessRightFact(t *testing.T) {
	logger := &MockLogger{}
	accumulator := domain.AccumulateFunction{
		FunctionType: "SUM",
		Field:        "amount",
	}
	accNode := NewAccumulateNode("acc1", accumulator, logger)

	// Ajouter un token de gauche
	token := &domain.Token{
		ID: "token1",
		Facts: []*domain.Fact{
			{
				ID:        "fact0",
				Type:      "Account",
				Fields:    map[string]interface{}{"id": "acc1"},
				Timestamp: time.Now(),
			},
		},
	}
	accNode.betaMemory.StoreToken(token)

	// Ajouter un premier fait de droite
	fact1 := &domain.Fact{
		ID:        "fact1",
		Type:      "Transaction",
		Fields:    map[string]interface{}{"amount": 100.0},
		Timestamp: time.Now(),
	}

	err := accNode.ProcessRightFact(fact1)
	if err != nil {
		t.Errorf("ProcessRightFact failed: %v", err)
	}

	// Vérifier que le fait est stocké
	facts := accNode.betaMemory.GetFacts()
	if len(facts) != 1 {
		t.Errorf("Expected 1 fact in memory, got %d", len(facts))
	}

	// Ajouter un deuxième fait de droite
	fact2 := &domain.Fact{
		ID:        "fact2",
		Type:      "Transaction",
		Fields:    map[string]interface{}{"amount": 200.0},
		Timestamp: time.Now(),
	}

	err = accNode.ProcessRightFact(fact2)
	if err != nil {
		t.Errorf("ProcessRightFact failed: %v", err)
	}

	// Vérifier que la somme a été mise à jour
	accNode.mu.RLock()
	result, exists := accNode.accumulatedValues[token.ID]
	accNode.mu.RUnlock()

	if !exists {
		t.Errorf("Expected accumulated value to be stored for token %s", token.ID)
	}

	expectedSum := 300.0
	if result != expectedSum {
		t.Errorf("Expected updated sum %v, got %v", expectedSum, result)
	}
}
