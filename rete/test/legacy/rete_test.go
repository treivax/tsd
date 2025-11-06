package rete

import (
	"fmt"
	"testing"
	"time"
)

func TestFactCreation(t *testing.T) {
	fact := &Fact{
		ID:   "test_1",
		Type: "TestType",
		Fields: map[string]interface{}{
			"name": "test",
			"age":  25,
		},
		Timestamp: time.Now(),
	}

	if fact.ID != "test_1" {
		t.Errorf("Expected ID 'test_1', got '%s'", fact.ID)
	}

	name, exists := fact.GetField("name")
	if !exists || name != "test" {
		t.Errorf("Expected field 'name' with value 'test'")
	}
}

func TestMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()

	// Test sauvegarde
	memory := &WorkingMemory{
		NodeID: "test_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	err := storage.SaveMemory("test_node", memory)
	if err != nil {
		t.Fatalf("Erreur sauvegarde: %v", err)
	}

	// Test chargement
	loadedMemory, err := storage.LoadMemory("test_node")
	if err != nil {
		t.Fatalf("Erreur chargement: %v", err)
	}

	if loadedMemory.NodeID != "test_node" {
		t.Errorf("Expected NodeID 'test_node', got '%s'", loadedMemory.NodeID)
	}

	// Test suppression
	err = storage.DeleteMemory("test_node")
	if err != nil {
		t.Fatalf("Erreur suppression: %v", err)
	}

	// Vérifier que c'est supprimé
	_, err = storage.LoadMemory("test_node")
	if err == nil {
		t.Errorf("Expected error when loading deleted memory")
	}
}

func TestReteNetworkCreation(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	if network.RootNode == nil {
		t.Errorf("RootNode should not be nil")
	}

	if network.RootNode.GetID() != "root" {
		t.Errorf("Expected root ID 'root', got '%s'", network.RootNode.GetID())
	}

	if len(network.TypeNodes) != 0 {
		t.Errorf("Expected 0 TypeNodes initially, got %d", len(network.TypeNodes))
	}
}

func TestTypeNodeValidation(t *testing.T) {
	typeDef := TypeDefinition{
		Type: "typeDefinition",
		Name: "TestType",
		Fields: []Field{
			{Name: "name", Type: "string"},
			{Name: "age", Type: "number"},
			{Name: "active", Type: "bool"},
		},
	}

	storage := NewMemoryStorage()
	typeNode := NewTypeNode("TestType", typeDef, storage)

	// Fait valide
	validFact := &Fact{
		ID:   "valid_1",
		Type: "TestType",
		Fields: map[string]interface{}{
			"name":   "Alice",
			"age":    25.0,
			"active": true,
		},
		Timestamp: time.Now(),
	}

	err := typeNode.ActivateRight(validFact)
	if err != nil {
		t.Errorf("Valid fact should not produce error: %v", err)
	}

	// Vérifier que le fait est en mémoire
	memory := typeNode.GetMemory()
	if len(memory.Facts) != 1 {
		t.Errorf("Expected 1 fact in memory, got %d", len(memory.Facts))
	}

	// Fait invalide (champ manquant)
	invalidFact := &Fact{
		ID:   "invalid_1",
		Type: "TestType",
		Fields: map[string]interface{}{
			"name": "Bob",
			// age manquant
			"active": true,
		},
		Timestamp: time.Now(),
	}

	err = typeNode.ActivateRight(invalidFact)
	if err == nil {
		t.Errorf("Invalid fact should produce error")
	}
}

func TestProgramLoading(t *testing.T) {
	program := &Program{
		Types: []TypeDefinition{
			{
				Type: "typeDefinition",
				Name: "Person",
				Fields: []Field{
					{Name: "name", Type: "string"},
					{Name: "age", Type: "number"},
				},
			},
		},
		Expressions: []Expression{
			{
				Type: "expression",
				Set: Set{
					Type: "set",
					Variables: []TypedVariable{
						{
							Type:     "typedVariable",
							Name:     "p",
							DataType: "Person",
						},
					},
				},
				Constraints: map[string]interface{}{
					"type": "simple",
				},
				Action: &Action{
					Type: "action",
					Job: JobCall{
						Type: "jobCall",
						Name: "notify",
						Args: []string{"p"},
					},
				},
			},
		},
	}

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	err := network.LoadFromAST(program)
	if err != nil {
		t.Fatalf("Error loading program: %v", err)
	}

	// Vérifier que les nœuds ont été créés
	if len(network.TypeNodes) != 1 {
		t.Errorf("Expected 1 TypeNode, got %d", len(network.TypeNodes))
	}

	if len(network.AlphaNodes) != 1 {
		t.Errorf("Expected 1 AlphaNode, got %d", len(network.AlphaNodes))
	}

	if len(network.TerminalNodes) != 1 {
		t.Errorf("Expected 1 TerminalNode, got %d", len(network.TerminalNodes))
	}
}

func TestFactSubmission(t *testing.T) {
	// Créer un réseau simple
	program := &Program{
		Types: []TypeDefinition{
			{
				Type: "typeDefinition",
				Name: "TestType",
				Fields: []Field{
					{Name: "value", Type: "number"},
				},
			},
		},
		Expressions: []Expression{
			{
				Type: "expression",
				Set: Set{
					Type: "set",
					Variables: []TypedVariable{
						{
							Type:     "typedVariable",
							Name:     "t",
							DataType: "TestType",
						},
					},
				},
				Constraints: map[string]interface{}{
					"type":     "binaryOperation",
					"operator": ">",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "t",
						"field":  "value",
					},
					"right": map[string]interface{}{
						"type":  "numberLiteral",
						"value": 0.0,
					},
				},
				Action: &Action{
					Type: "action",
					Job: JobCall{
						Type: "jobCall",
						Name: "process",
						Args: []string{"t"},
					},
				},
			},
		},
	}

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	err := network.LoadFromAST(program)
	if err != nil {
		t.Fatalf("Error loading program: %v", err)
	}

	// Soumettre un fait
	fact := &Fact{
		ID:   "fact_1",
		Type: "TestType",
		Fields: map[string]interface{}{
			"value": 42.0,
		},
		Timestamp: time.Now(),
	}

	err = network.SubmitFact(fact)
	if err != nil {
		t.Errorf("Error submitting fact: %v", err)
	}

	// Vérifier l'état du réseau
	state, err := network.GetNetworkState()
	if err != nil {
		t.Fatalf("Error getting network state: %v", err)
	}

	// Le fait devrait être dans le nœud racine
	rootMemory := state["root"]
	if rootMemory == nil || len(rootMemory.Facts) != 1 {
		t.Errorf("Expected 1 fact in root memory")
	}

	// Le fait devrait être dans le TypeNode
	typeMemory := state["type_TestType"]
	if typeMemory == nil || len(typeMemory.Facts) != 1 {
		t.Errorf("Expected 1 fact in type memory")
	}
}

// Benchmark pour tester les performances
func BenchmarkFactSubmission(b *testing.B) {
	program := &Program{
		Types: []TypeDefinition{
			{
				Type: "typeDefinition",
				Name: "BenchType",
				Fields: []Field{
					{Name: "id", Type: "number"},
				},
			},
		},
		Expressions: []Expression{
			{
				Type: "expression",
				Set: Set{
					Type: "set",
					Variables: []TypedVariable{
						{
							Type:     "typedVariable",
							Name:     "item",
							DataType: "BenchType",
						},
					},
				},
				Constraints: map[string]interface{}{"simple": true},
				Action: &Action{
					Type: "action",
					Job:  JobCall{Type: "jobCall", Name: "bench", Args: []string{"item"}},
				},
			},
		},
	}

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.LoadFromAST(program)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fact := &Fact{
			ID:   fmt.Sprintf("fact_%d", i),
			Type: "BenchType",
			Fields: map[string]interface{}{
				"id": float64(i),
			},
			Timestamp: time.Now(),
		}

		err := network.SubmitFact(fact)
		if err != nil {
			b.Fatalf("Error submitting fact: %v", err)
		}
	}
}
