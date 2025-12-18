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

func TestReteNetwork_InsertFact(t *testing.T) {
	t.Log("🧪 TEST InsertFact - Insertion dynamique de faits")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Test 1: Insertion simple
	fact1 := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Alice",
			"age":  30,
		},
	}

	err := network.InsertFact(fact1)
	if err != nil {
		t.Errorf("❌ InsertFact failed: %v", err)
	}

	// Vérifier que le fait est dans le storage
	storedFact := storage.GetFact("Person_p1")
	if storedFact == nil {
		t.Error("❌ Fact not found in storage")
	}

	// Test 2: Insertion avec ID déjà existant (doit échouer)
	fact2 := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Bob",
		},
	}

	err = network.InsertFact(fact2)
	if err == nil {
		t.Error("❌ Expected error when inserting duplicate fact")
	}

	// Test 3: Insertion avec fait nil
	err = network.InsertFact(nil)
	if err == nil {
		t.Error("❌ Expected error when inserting nil fact")
	}

	// Test 4: Insertion sans type
	fact3 := &Fact{
		ID: "p2",
		Fields: map[string]interface{}{
			"name": "Charlie",
		},
	}

	err = network.InsertFact(fact3)
	if err == nil {
		t.Error("❌ Expected error when inserting fact without type")
	}

	// Test 5: Insertion sans ID
	fact4 := &Fact{
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "David",
		},
	}

	err = network.InsertFact(fact4)
	if err == nil {
		t.Error("❌ Expected error when inserting fact without ID")
	}

	t.Log("✅ InsertFact tests passed")
}

func TestReteNetwork_UpdateFact(t *testing.T) {
	t.Log("🧪 TEST UpdateFact - Mise à jour dynamique de faits")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Setup: Ajouter un fait initial
	initialFact := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Alice",
			"age":  30,
		},
	}

	err := storage.AddFact(initialFact)
	if err != nil {
		t.Fatalf("❌ Setup failed: %v", err)
	}

	// Test 1: Mise à jour simple
	updatedFact := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Alice Smith",
			"age":  31,
		},
	}

	err = network.UpdateFact(updatedFact)
	if err != nil {
		t.Errorf("❌ UpdateFact failed: %v", err)
	}

	// Vérifier la mise à jour
	storedFact := storage.GetFact("Person_p1")
	if storedFact == nil {
		t.Fatal("❌ Fact not found after update")
	}

	if storedFact.Fields["age"] != 31 {
		t.Errorf("❌ Expected age 31, got %v", storedFact.Fields["age"])
	}

	if storedFact.Fields["name"] != "Alice Smith" {
		t.Errorf("❌ Expected name 'Alice Smith', got %v", storedFact.Fields["name"])
	}

	// Test 2: Mise à jour fait inexistant
	nonExistentFact := &Fact{
		ID:   "p999",
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Ghost",
		},
	}

	err = network.UpdateFact(nonExistentFact)
	if err == nil {
		t.Error("❌ Expected error when updating non-existent fact")
	}

	// Test 3: Mise à jour avec fait nil
	err = network.UpdateFact(nil)
	if err == nil {
		t.Error("❌ Expected error when updating nil fact")
	}

	// Test 4: Mise à jour sans type
	invalidFact := &Fact{
		ID: "p1",
		Fields: map[string]interface{}{
			"name": "Invalid",
		},
	}

	err = network.UpdateFact(invalidFact)
	if err == nil {
		t.Error("❌ Expected error when updating fact without type")
	}

	// Test 5: Mise à jour sans ID
	invalidFact2 := &Fact{
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Invalid",
		},
	}

	err = network.UpdateFact(invalidFact2)
	if err == nil {
		t.Error("❌ Expected error when updating fact without ID")
	}

	t.Log("✅ UpdateFact tests passed")
}

func TestReteNetwork_FactOperationsIntegration(t *testing.T) {
	t.Log("🧪 TEST Integration - Insert, Update, Retract")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Scenario: Insert -> Update -> Retract

	// 1. Insert
	fact1 := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Alice",
			"age":  30,
		},
	}

	err := network.InsertFact(fact1)
	if err != nil {
		t.Fatalf("❌ InsertFact failed: %v", err)
	}
	t.Log("✅ Step 1: Fact inserted")

	// Vérifier insertion
	storedFact := storage.GetFact("Person_p1")
	if storedFact == nil {
		t.Fatal("❌ Fact not found after insert")
	}
	if storedFact.Fields["age"] != 30 {
		t.Errorf("❌ Expected age 30, got %v", storedFact.Fields["age"])
	}

	// 2. Update
	fact2 := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Alice Smith",
			"age":  31,
		},
	}

	err = network.UpdateFact(fact2)
	if err != nil {
		t.Fatalf("❌ UpdateFact failed: %v", err)
	}
	t.Log("✅ Step 2: Fact updated")

	// Vérifier mise à jour
	storedFact = storage.GetFact("Person_p1")
	if storedFact == nil {
		t.Fatal("❌ Fact not found after update")
	}
	if storedFact.Fields["age"] != 31 {
		t.Errorf("❌ Expected age 31, got %v", storedFact.Fields["age"])
	}
	if storedFact.Fields["name"] != "Alice Smith" {
		t.Errorf("❌ Expected name 'Alice Smith', got %v", storedFact.Fields["name"])
	}

	// 3. Retract
	err = network.RetractFact("Person_p1")
	if err != nil {
		t.Fatalf("❌ RetractFact failed: %v", err)
	}
	t.Log("✅ Step 3: Fact retracted")

	// Vérifier suppression
	storedFact = storage.GetFact("Person_p1")
	if storedFact != nil {
		t.Error("❌ Fact should have been removed")
	}

	allFacts := storage.GetAllFacts()
	if len(allFacts) != 0 {
		t.Errorf("❌ Expected 0 facts, got %d", len(allFacts))
	}

	t.Log("✅ Integration test passed")
}
