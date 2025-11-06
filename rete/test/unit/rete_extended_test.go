package rete

import (
	"testing"
	"time"
)

// TestWorkingMemoryOperations teste les opérations sur la mémoire de travail
func TestWorkingMemoryOperations(t *testing.T) {
	// Test RemoveFact
	memory := &WorkingMemory{
		NodeID: "test_memory",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	// Créer un fait de test
	fact := &Fact{
		ID:        "test_fact",
		Type:      "TestType",
		Fields:    map[string]interface{}{"value": 42},
		Timestamp: time.Now(),
	}

	// Ajouter le fait
	memory.AddFact(fact)
	if len(memory.Facts) != 1 {
		t.Fatalf("Attendu 1 fait après ajout, obtenu %d", len(memory.Facts))
	}

	// Test RemoveFact
	memory.RemoveFact("test_fact")
	if len(memory.Facts) != 0 {
		t.Fatalf("Attendu 0 fait après suppression, obtenu %d", len(memory.Facts))
	}

	// Test RemoveFact avec ID inexistant (ne doit pas causer d'erreur)
	memory.RemoveFact("nonexistent")

	// Test GetFacts
	fact1 := &Fact{ID: "fact1", Type: "Type1", Fields: map[string]interface{}{"a": 1}, Timestamp: time.Now()}
	fact2 := &Fact{ID: "fact2", Type: "Type2", Fields: map[string]interface{}{"b": 2}, Timestamp: time.Now()}

	memory.AddFact(fact1)
	memory.AddFact(fact2)

	facts := memory.GetFacts()
	if len(facts) != 2 {
		t.Fatalf("Attendu 2 faits, obtenu %d", len(facts))
	}

	// Vérifier que les faits sont bien présents
	foundFact1, foundFact2 := false, false
	for _, f := range facts {
		if f.ID == "fact1" {
			foundFact1 = true
		}
		if f.ID == "fact2" {
			foundFact2 = true
		}
	}

	if !foundFact1 || !foundFact2 {
		t.Fatal("Tous les faits n'ont pas été récupérés par GetFacts")
	}
}

// TestNodeTypes teste les fonctions GetType
func TestNodeTypes(t *testing.T) {
	storage := NewMemoryStorage()

	// Test RootNode GetType
	rootNode := NewRootNode(storage)
	if rootNode.GetType() != "root" {
		t.Fatalf("Type du RootNode incorrect: attendu 'root', obtenu '%s'", rootNode.GetType())
	}

	// Test TypeNode GetType
	typeDef := TypeDefinition{
		Type:   "typeDefinition",
		Name:   "TestType",
		Fields: []Field{{Name: "test", Type: "string"}},
	}
	typeNode := NewTypeNode("TestType", typeDef, storage)
	if typeNode.GetType() != "type" {
		t.Fatalf("Type du TypeNode incorrect: attendu 'type', obtenu '%s'", typeNode.GetType())
	}

	// Test AlphaNode GetType
	builder := NewAlphaConditionBuilder()
	condition := builder.FieldEquals("test", "value", "value")
	alphaNode := NewAlphaNode("alpha_test", condition, "x", storage)
	if alphaNode.GetType() != "alpha" {
		t.Fatalf("Type de l'AlphaNode incorrect: attendu 'alpha', obtenu '%s'", alphaNode.GetType())
	}

	// Test TerminalNode GetType
	action := &Action{
		Type: "print",
		Job:  JobCall{Type: "job", Name: "print", Args: []string{"message"}},
	}
	terminalNode := NewTerminalNode("test_rule", action, storage)
	if terminalNode.GetType() != "terminal" {
		t.Fatalf("Type du TerminalNode incorrect: attendu 'terminal', obtenu '%s'", terminalNode.GetType())
	}
}

// TestActivateLeft teste les fonctions ActivateLeft non couvertes
func TestActivateLeft(t *testing.T) {
	storage := NewMemoryStorage()

	// Test RootNode ActivateLeft - ne devrait rien faire
	rootNode := NewRootNode(storage)
	token := &Token{
		ID:     "test_token",
		Facts:  []*Fact{},
		NodeID: "root",
	}

	// Cette fonction ne devrait pas causer d'erreur même si elle ne fait rien
	rootNode.ActivateLeft(token)

	// Test TypeNode ActivateLeft - ne devrait rien faire
	typeDef := TypeDefinition{
		Type:   "typeDefinition",
		Name:   "TestType",
		Fields: []Field{{Name: "test", Type: "string"}},
	}
	typeNode := NewTypeNode("TestType", typeDef, storage)
	typeNode.ActivateLeft(token)

	// Test AlphaNode ActivateLeft - ne devrait rien faire
	builder := NewAlphaConditionBuilder()
	condition := builder.FieldEquals("test", "value", "value")
	alphaNode := NewAlphaNode("alpha_test", condition, "x", storage)
	alphaNode.ActivateLeft(token)
}

// TestTerminalNodeActivateRight teste ActivateRight du TerminalNode
func TestTerminalNodeActivateRight(t *testing.T) {
	// Créer un TerminalNode avec une action simple
	action := &Action{
		Type: "print",
		Job:  JobCall{Type: "job", Name: "print", Args: []string{"message"}},
	}

	terminalNode := NewTerminalNode("test_rule", action, NewMemoryStorage())

	// Créer un fait de test
	fact := &Fact{
		ID:        "terminal_test_fact",
		Type:      "TestType",
		Fields:    map[string]interface{}{"value": 123},
		Timestamp: time.Now(),
	}

	// Activer le nœud terminal - devrait retourner une erreur car les nœuds terminaux ne reçoivent pas directement de faits
	err := terminalNode.ActivateRight(fact)
	if err == nil {
		t.Fatal("ActivateRight devrait retourner une erreur pour un nœud terminal")
	}
	// Vérifier que l'erreur est celle attendue
	expectedError := "les nœuds terminaux ne reçoivent pas de faits directement"
	if err.Error() != expectedError {
		t.Fatalf("Erreur inattendue: attendu '%s', obtenu '%s'", expectedError, err.Error())
	}
}

// TestSaveMemory teste la fonction SaveMemory des nœuds
func TestSaveMemory(t *testing.T) {
	// Test avec un AlphaNode
	builder := NewAlphaConditionBuilder()
	condition := builder.FieldEquals("test", "value", "value")

	storage := NewMemoryStorage()
	alphaNode := NewAlphaNode("save_test_node", condition, "x", storage)

	// Ajouter un fait à la mémoire du nœud
	fact := &Fact{
		ID:        "save_test_fact",
		Type:      "TestType",
		Fields:    map[string]interface{}{"test": "value"},
		Timestamp: time.Now(),
	}

	alphaNode.Memory.AddFact(fact)

	// Sauvegarder la mémoire
	err := alphaNode.SaveMemory()
	if err != nil {
		t.Fatalf("SaveMemory a échoué: %v", err)
	}

	// Vérifier que la mémoire a été sauvegardée
	nodes, err := storage.ListNodes()
	if err != nil {
		t.Fatalf("Erreur lors du listage des nœuds: %v", err)
	}

	found := false
	for _, nodeID := range nodes {
		if nodeID == alphaNode.ID {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("La mémoire du nœud n'a pas été sauvegardée")
	}
}

// TestNodeMemoryOperations teste les opérations sur la mémoire des nœuds
func TestNodeMemoryOperations(t *testing.T) {
	// Créer différents types de nœuds et tester leurs mémoires
	builder := NewAlphaConditionBuilder()
	// Utilisons une condition plus simple avec une constante
	condition := builder.True()

	alphaNode := NewAlphaNode("memory_test_node", condition, "x", NewMemoryStorage())

	// Test avec n'importe quel fait (la condition True devrait toujours passer)
	matchingFact := &Fact{
		ID:        "matching_fact",
		Type:      "Document",
		Fields:    map[string]interface{}{"category": "important", "title": "Test Doc"},
		Timestamp: time.Now(),
	}

	err := alphaNode.ActivateRight(matchingFact)
	if err != nil {
		t.Fatalf("ActivateRight a échoué: %v", err)
	}

	// Vérifier que le fait a été ajouté à la mémoire
	if len(alphaNode.Memory.Facts) != 1 {
		t.Fatalf("Attendu 1 fait dans la mémoire, obtenu %d", len(alphaNode.Memory.Facts))
	}

	// Test avec un autre fait
	anotherFact := &Fact{
		ID:        "another_fact",
		Type:      "Document",
		Fields:    map[string]interface{}{"category": "normal", "title": "Other Doc"},
		Timestamp: time.Now(),
	}

	err = alphaNode.ActivateRight(anotherFact)
	if err != nil {
		t.Fatalf("ActivateRight ne devrait pas échouer: %v", err)
	}

	// Vérifier que le deuxième fait a également été ajouté (condition True)
	if len(alphaNode.Memory.Facts) != 2 {
		t.Fatalf("Attendu 2 faits dans la mémoire, obtenu %d", len(alphaNode.Memory.Facts))
	}
}
