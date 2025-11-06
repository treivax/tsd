package domain
package domain

import (
	"testing"
)

func TestFact(t *testing.T) {
	fact := NewFact("temperature", 25.5)
	
	if fact.Type != "temperature" {
		t.Errorf("Attendu type 'temperature', reçu '%s'", fact.Type)
	}
	
	if fact.Value != 25.5 {
		t.Errorf("Attendu valeur 25.5, reçu %v", fact.Value)
	}
	
	if fact.ID == "" {
		t.Error("L'ID ne doit pas être vide")
	}
}

func TestToken(t *testing.T) {
	fact := NewFact("test", "value")
	token := NewToken(fact)
	
	if token.Fact != fact {
		t.Error("Le token devrait contenir le fait")
	}
	
	if token.ID == "" {
		t.Error("L'ID du token ne doit pas être vide")
	}
}

func TestWorkingMemory(t *testing.T) {
	wm := NewWorkingMemory()
	
	if wm == nil {
		t.Fatal("WorkingMemory ne doit pas être nil")
	}
	
	// Test ajout de fait
	fact := NewFact("test", "value")
	wm.AddFact(fact)
	
	facts := wm.GetFacts()
	if len(facts) != 1 {
		t.Errorf("Attendu 1 fait, reçu %d", len(facts))
	}
	
	// Test suppression de fait
	wm.RemoveFact(fact.ID)
	facts = wm.GetFacts()
	if len(facts) != 0 {
		t.Errorf("Attendu 0 fait après suppression, reçu %d", len(facts))
	}
}

func TestValidationError(t *testing.T) {
	err := NewValidationError("test error", "field", "value")
	
	if err == nil {
		t.Fatal("L'erreur ne doit pas être nil")
	}
	
	if !IsValidationError(err) {
		t.Error("Devrait être une ValidationError")
	}
	
	// Test que ce n'est pas une NodeError
	if IsNodeError(err) {
		t.Error("Ne devrait pas être une NodeError")
	}
}

func TestNodeError(t *testing.T) {
	err := NewNodeError("node-1", "test error")
	
	if err == nil {
		t.Fatal("L'erreur ne doit pas être nil")
	}
	
	if !IsNodeError(err) {
		t.Error("Devrait être une NodeError")
	}
	
	// Test que ce n'est pas une ValidationError
	if IsValidationError(err) {
		t.Error("Ne devrait pas être une ValidationError")
	}
}