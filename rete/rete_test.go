package rete

import (
	"testing"
	"time"
)

// ========== TESTS DE BASE ==========

func TestFact_Creation(t *testing.T) {
	fact := &Fact{
		ID:   "test_001",
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Alice",
			"age":  30,
		},
		Timestamp: time.Now(),
	}

	if fact.ID != "test_001" {
		t.Errorf("Expected ID 'test_001', got '%s'", fact.ID)
	}
	if fact.Type != "Person" {
		t.Errorf("Expected Type 'Person', got '%s'", fact.Type)
	}
}

func TestFact_GetField(t *testing.T) {
	fact := &Fact{
		ID:   "test_001",
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Alice",
			"age":  30,
		},
	}

	// Test champ existant
	name, exists := fact.GetField("name")
	if !exists {
		t.Error("Field 'name' should exist")
	}
	if name != "Alice" {
		t.Errorf("Expected name 'Alice', got '%v'", name)
	}

	// Test champ inexistant
	_, exists = fact.GetField("city")
	if exists {
		t.Error("Field 'city' should not exist")
	}
}

func TestWorkingMemory_AddFact(t *testing.T) {
	wm := &WorkingMemory{
		NodeID: "test_node",
		Facts:  make(map[string]*Fact),
	}

	fact := &Fact{ID: "f1", Type: "Person"}
	wm.AddFact(fact)

	if len(wm.Facts) != 1 {
		t.Errorf("Expected 1 fact, got %d", len(wm.Facts))
	}

	retrieved, exists := wm.GetFact("f1")
	if !exists {
		t.Error("Fact should exist in working memory")
	}
	if retrieved.ID != "f1" {
		t.Errorf("Expected fact ID 'f1', got '%s'", retrieved.ID)
	}
}

func TestWorkingMemory_RemoveFact(t *testing.T) {
	wm := &WorkingMemory{
		NodeID: "test_node",
		Facts:  make(map[string]*Fact),
	}

	fact := &Fact{ID: "f1", Type: "Person"}
	wm.AddFact(fact)
	wm.RemoveFact("f1")

	if len(wm.Facts) != 0 {
		t.Errorf("Expected 0 facts after removal, got %d", len(wm.Facts))
	}

	_, exists := wm.GetFact("f1")
	if exists {
		t.Error("Fact should not exist after removal")
	}
}

func TestRootNode_ActivateRetract(t *testing.T) {
	storage := NewMemoryStorage()
	root := NewRootNode(storage)

	fact := &Fact{ID: "f1", Type: "Person"}
	root.ActivateRight(fact)

	// Rétracter le fait
	err := root.ActivateRetract("f1")
	if err != nil {
		t.Errorf("ActivateRetract failed: %v", err)
	}

	// Vérifier que le fait a été supprimé
	memory := root.GetMemory()
	if len(memory.Facts) != 0 {
		t.Errorf("Expected 0 facts after retract, got %d", len(memory.Facts))
	}
}

func TestTypeNode_ActivateRetract(t *testing.T) {
	storage := NewMemoryStorage()
	typeDef := TypeDefinition{
		Name:   "Person",
		Fields: []Field{{Name: "name", Type: "string"}},
	}

	typeNode := NewTypeNode("Person", typeDef, storage)

	fact := &Fact{
		ID:     "p1",
		Type:   "Person",
		Fields: map[string]interface{}{"name": "Alice"},
	}

	typeNode.ActivateRight(fact)
	typeNode.ActivateRetract("p1")

	memory := typeNode.GetMemory()
	if len(memory.Facts) != 0 {
		t.Errorf("Expected 0 facts after retract, got %d", len(memory.Facts))
	}
}

func TestAlphaNode_ActivateRetract(t *testing.T) {
	storage := NewMemoryStorage()
	alphaNode := NewAlphaNode("alpha_1", nil, "p", storage)

	fact := &Fact{ID: "f1", Type: "Person"}
	alphaNode.Memory.AddFact(fact)

	err := alphaNode.ActivateRetract("f1")
	if err != nil {
		t.Errorf("ActivateRetract failed: %v", err)
	}

	memory := alphaNode.GetMemory()
	if len(memory.Facts) != 0 {
		t.Errorf("Expected 0 facts after retract, got %d", len(memory.Facts))
	}
}

func TestTerminalNode_ActivateRetract(t *testing.T) {
	storage := NewMemoryStorage()
	action := &Action{
		Job: JobCall{Name: "alert", Args: []interface{}{}},
	}

	terminal := NewTerminalNode("term_1", action, storage)

	fact := &Fact{ID: "f1", Type: "Person"}
	token := &Token{
		ID:    "tok_1",
		Facts: []*Fact{fact},
	}

	terminal.ActivateLeft(token)

	// Rétracter le fait
	err := terminal.ActivateRetract("f1")
	if err != nil {
		t.Errorf("ActivateRetract failed: %v", err)
	}

	memory := terminal.GetMemory()
	if len(memory.Tokens) != 0 {
		t.Errorf("Expected 0 tokens after retract, got %d", len(memory.Tokens))
	}
}

func TestJoinNode_ActivateRetract(t *testing.T) {
	storage := NewMemoryStorage()
	joinNode := NewJoinNode("join_1", nil, []string{"p"}, []string{"o"}, map[string]string{}, storage)

	// Ajouter des tokens dans les mémoires
	fact1 := &Fact{ID: "p1", Type: "Person"}
	token1 := &Token{
		ID:       "tok_p1",
		Facts:    []*Fact{fact1},
		Bindings: map[string]*Fact{"p": fact1},
	}
	joinNode.LeftMemory.AddToken(token1)

	fact2 := &Fact{ID: "o1", Type: "Order"}
	token2 := &Token{
		ID:       "tok_o1",
		Facts:    []*Fact{fact2},
		Bindings: map[string]*Fact{"o": fact2},
	}
	joinNode.RightMemory.AddToken(token2)

	// Rétracter p1
	err := joinNode.ActivateRetract("p1")
	if err != nil {
		t.Errorf("ActivateRetract failed: %v", err)
	}

	// Vérifier que le token de gauche a été supprimé
	leftTokens := joinNode.LeftMemory.GetTokens()
	if len(leftTokens) != 0 {
		t.Errorf("Expected 0 tokens in left memory after retract, got %d", len(leftTokens))
	}
}

func TestExistsNode_ActivateRetract(t *testing.T) {
	storage := NewMemoryStorage()

	existsConditions := map[string]interface{}{}
	existsNode := NewExistsNode("exists_1", existsConditions, "p", "o", map[string]string{}, storage)

	// Ajouter un fait dans la mémoire d'existence
	fact := &Fact{ID: "o1", Type: "Order"}
	existsNode.ExistsMemory.AddFact(fact)

	// Rétracter le fait d'existence
	err := existsNode.ActivateRetract("o1")
	if err != nil {
		t.Errorf("ActivateRetract failed: %v", err)
	}

	// Vérifier que le fait a été supprimé
	existsFacts := existsNode.ExistsMemory.GetFacts()
	if len(existsFacts) != 0 {
		t.Errorf("Expected 0 facts in exists memory after retract, got %d", len(existsFacts))
	}
}
