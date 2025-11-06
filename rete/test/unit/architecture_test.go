package rete
package rete

import (
	"testing"
	
	"github.com/treivax/tsd/rete/pkg/domain"
	"github.com/treivax/tsd/rete/pkg/nodes"
)

// TestNewArchitecture teste la nouvelle architecture du module RETE
func TestNewArchitecture(t *testing.T) {
	t.Log("Test de la nouvelle architecture RETE")

	// Test des types de base du domaine
	t.Run("TestDomainFact", func(t *testing.T) {
		fact := domain.NewFact("temperature", 25.5)
		if fact.Type != "temperature" {
			t.Errorf("Attendu type 'temperature', reçu '%s'", fact.Type)
		}
		if fact.Value != 25.5 {
			t.Errorf("Attendu valeur 25.5, reçu %v", fact.Value)
		}
	})

	t.Run("TestWorkingMemory", func(t *testing.T) {
		wm := domain.NewWorkingMemory()
		if wm == nil {
			t.Fatal("WorkingMemory ne doit pas être nil")
		}

		fact := domain.NewFact("test", "value")
		wm.AddFact(fact)

		facts := wm.GetFacts()
		if len(facts) != 1 {
			t.Errorf("Attendu 1 fait, reçu %d", len(facts))
		}
	})

	t.Run("TestToken", func(t *testing.T) {
		fact := domain.NewFact("test", "value")
		token := domain.NewToken(fact)
		
		if token.Fact != fact {
			t.Error("Le token devrait contenir le fait")
		}
	})

	// Test du BaseNode
	t.Run("TestBaseNode", func(t *testing.T) {
		node := nodes.NewBaseNode("test-node", "TestNode")
		if node.GetID() != "test-node" {
			t.Errorf("Attendu ID 'test-node', reçu '%s'", node.GetID())
		}
		if node.GetType() != "TestNode" {
			t.Errorf("Attendu type 'TestNode', reçu '%s'", node.GetType())
		}
	})

	// Test des erreurs structurées
	t.Run("TestDomainErrors", func(t *testing.T) {
		err := domain.NewValidationError("test error", "field", "value")
		if err == nil {
			t.Fatal("L'erreur ne doit pas être nil")
		}

		if !domain.IsValidationError(err) {
			t.Error("Devrait être une ValidationError")
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