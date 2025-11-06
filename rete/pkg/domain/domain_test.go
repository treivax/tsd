package domain

import (
	"errors"
	"testing"
)

func TestFact(t *testing.T) {
	t.Run("NewFact", func(t *testing.T) {
		fields := map[string]interface{}{
			"value": 25.5,
		}
		fact := NewFact("temperature", "sensor", fields)
		
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
}

func TestToken(t *testing.T) {
	t.Run("NewToken", func(t *testing.T) {
		fields := map[string]interface{}{
			"value": "test",
		}
		fact := NewFact("test", "type", fields)
		facts := []*Fact{fact}
		token := NewToken("token1", "node1", facts)
		
		if len(token.Facts) != 1 {
			t.Error("Le token devrait contenir 1 fait")
		}
		if token.Facts[0] != fact {
			t.Error("Le token devrait contenir le fait correct")
		}
		if token.NodeID != "node1" {
			t.Errorf("Attendu NodeID 'node1', reçu '%s'", token.NodeID)
		}
	})
}

func TestWorkingMemory(t *testing.T) {
	t.Run("NewWorkingMemory", func(t *testing.T) {
		wm := NewWorkingMemory("test-node")
		if wm == nil {
			t.Fatal("WorkingMemory ne doit pas être nil")
		}
		if wm.NodeID != "test-node" {
			t.Errorf("Attendu NodeID 'test-node', reçu '%s'", wm.NodeID)
		}
	})

	t.Run("AddFact", func(t *testing.T) {
		wm := NewWorkingMemory("test-node")
		fields := map[string]interface{}{
			"value": "test",
		}
		fact := NewFact("test", "type", fields)
		wm.AddFact(fact)

		if len(wm.Facts) != 1 {
			t.Errorf("Attendu 1 fait, reçu %d", len(wm.Facts))
		}
	})
}

func TestErrors(t *testing.T) {
	t.Run("ValidationError", func(t *testing.T) {
		err := NewValidationError("field", "value", "test error")
		if err == nil {
			t.Fatal("L'erreur ne doit pas être nil")
		}
		if err.Field != "field" {
			t.Errorf("Attendu field 'field', reçu '%s'", err.Field)
		}
	})

	t.Run("NodeError", func(t *testing.T) {
		cause := errors.New("cause error")
		err := NewNodeError("node-1", "TestNode", cause)
		if err == nil {
			t.Fatal("L'erreur ne doit pas être nil")
		}
		if err.NodeID != "node-1" {
			t.Errorf("Attendu NodeID 'node-1', reçu '%s'", err.NodeID)
		}
		if err.Unwrap() != cause {
			t.Error("L'erreur devrait wrapper la cause")
		}
	})
}