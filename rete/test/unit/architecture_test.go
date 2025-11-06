package rete

import (
	"testing"
	
	"github.com/treivax/tsd/rete/pkg/domain"
	"github.com/treivax/tsd/rete/pkg/nodes"
)

// TestNewArchitecture teste la nouvelle architecture du module RETE
func TestNewArchitecture(t *testing.T) {
	t.Log("Test de la nouvelle architecture RETE")

		// Test NewFact
	t.Run("NewFact", func(t *testing.T) {
		fields := map[string]interface{}{
			"value": 25.5,
		}
		fact := domain.NewFact("temperature", "sensor", fields)
		
		if fact.ID != "temperature" {
			t.Errorf("Attendu ID 'temperature', reçu '%s'", fact.ID)
		}
		if fact.Type != "sensor" {
			t.Errorf("Attendu type 'sensor', reçu '%s'", fact.Type)
		}
		if val, exists := fact.GetField("value"); !exists || val != 25.5 {
			t.Errorf("Attendu valeur 25.5, reçu %v", val)
		}
	})

	// Test NewWorkingMemory
	t.Run("NewWorkingMemory", func(t *testing.T) {
		wm := domain.NewWorkingMemory("test-node")
		
		if wm == nil {
			t.Error("WorkingMemory ne devrait pas être nil")
		}
		if wm.NodeID != "test-node" {
			t.Errorf("Attendu NodeID 'test-node', reçu '%s'", wm.NodeID)
		}
	})

	// Test AddFact
	t.Run("AddFact", func(t *testing.T) {
		fields := map[string]interface{}{
			"value": "test",
		}
		fact := domain.NewFact("test", "type", fields)
		wm := domain.NewWorkingMemory("test-node")
		wm.AddFact(fact)
		
		// Vérifier que le fait a été ajouté
		if len(wm.Facts) != 1 {
			t.Errorf("Attendu 1 fait, reçu %d", len(wm.Facts))
		}
	})

	// Test Token
	t.Run("NewToken", func(t *testing.T) {
		fields := map[string]interface{}{
			"value": "test",
		}
		fact := domain.NewFact("test", "type", fields)
		facts := []*domain.Fact{fact}
		token := domain.NewToken("token1", "node1", facts)
		
		if len(token.Facts) != 1 {
			t.Error("Le token devrait contenir 1 fait")
		}
		if token.Facts[0] != fact {
			t.Error("Le token devrait contenir le fait correct")
		}
	})

	t.Run("TestWorkingMemory", func(t *testing.T) {
		wm := domain.NewWorkingMemory("test-node")
		if wm == nil {
			t.Fatal("WorkingMemory ne doit pas être nil")
		}

		fields := map[string]interface{}{
			"value": "test",
		}
		fact := domain.NewFact("test", "type", fields)
		wm.AddFact(fact)

		if len(wm.Facts) != 1 {
			t.Errorf("Attendu 1 fait, reçu %d", len(wm.Facts))
		}
	})

	// Test du BaseNode avec un logger mock
	t.Run("TestBaseNode", func(t *testing.T) {
		// Logger mock simple
		logger := &mockLogger{}
		node := nodes.NewBaseNode("test-node", "TestNode", logger)
		if node.ID() != "test-node" {
			t.Errorf("Attendu ID 'test-node', reçu '%s'", node.ID())
		}
		if node.Type() != "TestNode" {
			t.Errorf("Attendu type 'TestNode', reçu '%s'", node.Type())
		}
	})

	// Test des erreurs structurées
	t.Run("TestDomainErrors", func(t *testing.T) {
		err := domain.NewValidationError("field", "value", "test error")
		if err == nil {
			t.Fatal("L'erreur ne doit pas être nil")
		}

		if err.Field != "field" {
			t.Errorf("Attendu field 'field', reçu '%s'", err.Field)
		}
	})
}

// TestCompatibility teste la compatibilité avec l'ancienne API
func TestCompatibility(t *testing.T) {
	t.Log("Test de compatibilité avec l'ancienne API")
	
	// Vérifier que les fonctions principales existent toujours
	t.Run("TestExistingFunctions", func(t *testing.T) {
		// Test que le module principal compile toujours
		// Ces tests peuvent être étendus au fur et à mesure de la migration
		t.Log("Les fonctions de base doivent être disponibles")
	})
}

// mockLogger implémente l'interface Logger pour les tests
type mockLogger struct{}

func (m *mockLogger) Debug(msg string, fields map[string]interface{}) {}
func (m *mockLogger) Info(msg string, fields map[string]interface{})  {}
func (m *mockLogger) Warn(msg string, fields map[string]interface{})  {}
func (m *mockLogger) Error(msg string, err error, fields map[string]interface{}) {}