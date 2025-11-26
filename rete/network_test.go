// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

func TestReteNetwork_Creation(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	if network.RootNode == nil {
		t.Error("RootNode should not be nil")
	}
	if network.TypeNodes == nil {
		t.Error("TypeNodes should not be nil")
	}
}

func TestReteNetwork_SubmitFact(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Ajouter un TypeNode
	typeDef := TypeDefinition{
		Name:   "Person",
		Fields: []Field{{Name: "name", Type: "string"}},
	}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode
	network.RootNode.AddChild(typeNode)

	// Soumettre un fait
	fact := &Fact{
		ID:     "p1",
		Type:   "Person",
		Fields: map[string]interface{}{"name": "Alice"},
	}

	err := network.SubmitFact(fact)
	if err != nil {
		t.Errorf("SubmitFact failed: %v", err)
	}

	// Vérifier que le fait est dans le RootNode
	rootMemory := network.RootNode.GetMemory()
	if len(rootMemory.Facts) != 1 {
		t.Errorf("Expected 1 fact in root, got %d", len(rootMemory.Facts))
	}
}

func TestReteNetwork_RetractFact(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Ajouter un TypeNode
	typeDef := TypeDefinition{
		Name:   "Person",
		Fields: []Field{{Name: "name", Type: "string"}},
	}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode
	network.RootNode.AddChild(typeNode)

	// Soumettre un fait
	fact := &Fact{
		ID:     "p1",
		Type:   "Person",
		Fields: map[string]interface{}{"name": "Alice"},
	}
	network.SubmitFact(fact)

	// Rétracter le fait
	err := network.RetractFact("Person_p1")
	if err != nil {
		t.Errorf("RetractFact failed: %v", err)
	}

	// Vérifier que le fait a été supprimé
	rootMemory := network.RootNode.GetMemory()
	if len(rootMemory.Facts) != 0 {
		t.Errorf("Expected 0 facts in root after retract, got %d", len(rootMemory.Facts))
	}

	typeMemory := typeNode.GetMemory()
	if len(typeMemory.Facts) != 0 {
		t.Errorf("Expected 0 facts in type node after retract, got %d", len(typeMemory.Facts))
	}
}

func TestReteNetwork_RetractFact_NotFound(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Essayer de rétracter un fait qui n'existe pas
	err := network.RetractFact("Person_non_existent")
	if err == nil {
		t.Error("RetractFact should error when fact not found")
	}
}

func TestScenario_AddAndRetractMultipleFacts(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Setup réseau
	typeDef := TypeDefinition{
		Name:   "Person",
		Fields: []Field{{Name: "name", Type: "string"}},
	}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode
	network.RootNode.AddChild(typeNode)

	// Ajouter plusieurs faits
	facts := []*Fact{
		{ID: "p1", Type: "Person", Fields: map[string]interface{}{"name": "Alice"}},
		{ID: "p2", Type: "Person", Fields: map[string]interface{}{"name": "Bob"}},
		{ID: "p3", Type: "Person", Fields: map[string]interface{}{"name": "Charlie"}},
	}

	for _, fact := range facts {
		if err := network.SubmitFact(fact); err != nil {
			t.Errorf("SubmitFact failed for %s: %v", fact.ID, err)
		}
	}

	// Vérifier que tous sont présents
	rootMemory := network.RootNode.GetMemory()
	if len(rootMemory.Facts) != 3 {
		t.Errorf("Expected 3 facts in root, got %d", len(rootMemory.Facts))
	}

	// Rétracter le fait du milieu
	if err := network.RetractFact("Person_p2"); err != nil {
		t.Errorf("RetractFact failed: %v", err)
	}

	// Vérifier qu'il reste 2 faits
	rootMemory = network.RootNode.GetMemory()
	if len(rootMemory.Facts) != 2 {
		t.Errorf("Expected 2 facts in root after retract, got %d", len(rootMemory.Facts))
	}

	// Vérifier que p1 et p3 sont toujours là
	if _, exists := rootMemory.GetFact("Person_p1"); !exists {
		t.Error("p1 should still exist")
	}
	if _, exists := rootMemory.GetFact("Person_p3"); !exists {
		t.Error("p3 should still exist")
	}
	if _, exists := rootMemory.GetFact("Person_p2"); exists {
		t.Error("p2 should be removed")
	}
}
