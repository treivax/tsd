package rete

import (
	"testing"
	"time"
)

// TestMemoryStorageOperations teste toutes les opérations de stockage en mémoire
func TestMemoryStorageOperations(t *testing.T) {
	storage := NewMemoryStorage()

	// Test ListNodes avec stockage vide
	nodes, err := storage.ListNodes()
	if err != nil {
		t.Fatalf("Erreur lors du listage des nœuds vides: %v", err)
	}
	if len(nodes) != 0 {
		t.Fatalf("Attendu 0 nœuds, obtenu %d", len(nodes))
	}

	// Créer des mémoires de test
	fact1 := &Fact{
		ID:        "fact1",
		Type:      "TestType",
		Fields:    map[string]interface{}{"name": "test1", "value": 42},
		Timestamp: time.Now(),
	}
	fact2 := &Fact{
		ID:        "fact2",
		Type:      "TestType",
		Fields:    map[string]interface{}{"name": "test2", "value": 84},
		Timestamp: time.Now(),
	}

	memory1 := &WorkingMemory{
		NodeID: "node1",
		Facts:  map[string]*Fact{"fact1": fact1},
		Tokens: make(map[string]*Token),
	}

	memory2 := &WorkingMemory{
		NodeID: "node2",
		Facts:  map[string]*Fact{"fact2": fact2},
		Tokens: make(map[string]*Token),
	}

	// Test SaveMemory
	err = storage.SaveMemory("node1", memory1)
	if err != nil {
		t.Fatalf("Erreur lors de la sauvegarde de la mémoire: %v", err)
	}

	err = storage.SaveMemory("node2", memory2)
	if err != nil {
		t.Fatalf("Erreur lors de la sauvegarde de la deuxième mémoire: %v", err)
	}

	// Test ListNodes avec stockage non vide
	nodes, err = storage.ListNodes()
	if err != nil {
		t.Fatalf("Erreur lors du listage des nœuds: %v", err)
	}
	if len(nodes) != 2 {
		t.Fatalf("Attendu 2 nœuds, obtenu %d", len(nodes))
	}

	// Vérifier que les deux nœuds sont présents
	nodeMap := make(map[string]bool)
	for _, node := range nodes {
		nodeMap[node] = true
	}
	if !nodeMap["node1"] || !nodeMap["node2"] {
		t.Fatalf("Les nœuds attendus ne sont pas présents: %v", nodes)
	}

	// Test LoadMemory
	loadedMemory, err := storage.LoadMemory("node1")
	if err != nil {
		t.Fatalf("Erreur lors du chargement de la mémoire: %v", err)
	}

	if len(loadedMemory.Facts) != 1 {
		t.Fatalf("Attendu 1 fait, obtenu %d", len(loadedMemory.Facts))
	}

	// Vérifier le fait chargé par son ID
	loadedFact, exists := loadedMemory.Facts["fact1"]
	if !exists {
		t.Fatal("Le fait fact1 n'a pas été trouvé dans la mémoire chargée")
	}

	if loadedFact.ID != "fact1" {
		t.Fatalf("Attendu ID fact1, obtenu %s", loadedFact.ID)
	}

	// Test LoadMemory avec nœud inexistant
	_, err = storage.LoadMemory("nonexistent")
	if err == nil {
		t.Fatal("Attendu une erreur pour un nœud inexistant")
	}

	// Test DeleteMemory
	err = storage.DeleteMemory("node1")
	if err != nil {
		t.Fatalf("Erreur lors de la suppression: %v", err)
	}

	// Vérifier que le nœud a été supprimé
	nodes, err = storage.ListNodes()
	if err != nil {
		t.Fatalf("Erreur lors du listage après suppression: %v", err)
	}
	if len(nodes) != 1 {
		t.Fatalf("Attendu 1 nœud après suppression, obtenu %d", len(nodes))
	}

	// Vérifier que le chargement du nœud supprimé échoue
	_, err = storage.LoadMemory("node1")
	if err == nil {
		t.Fatal("Attendu une erreur pour un nœud supprimé")
	}
}

// TestMemoryStorageConcurrency teste les opérations concurrentes
func TestMemoryStorageConcurrency(t *testing.T) {
	storage := NewMemoryStorage()

	// Test de modifications concurrentes pour vérifier la sécurité des mutex
	concurrentFact := &Fact{
		ID:        "concurrent_fact",
		Type:      "TestType",
		Fields:    map[string]interface{}{"test": true},
		Timestamp: time.Now(),
	}

	memory := &WorkingMemory{
		NodeID: "concurrent_node",
		Facts:  map[string]*Fact{"concurrent_fact": concurrentFact},
		Tokens: make(map[string]*Token),
	}

	// Sauvegarde initiale
	err := storage.SaveMemory("concurrent_node", memory)
	if err != nil {
		t.Fatalf("Erreur lors de la sauvegarde initiale: %v", err)
	}

	// Test de lecture concurrente
	go func() {
		_, _ = storage.LoadMemory("concurrent_node")
	}()

	// Test de listage concurrent
	go func() {
		_, _ = storage.ListNodes()
	}()

	// Test de suppression
	err = storage.DeleteMemory("concurrent_node")
	if err != nil {
		t.Fatalf("Erreur lors de la suppression concurrente: %v", err)
	}
}

// TestWorkingMemoryIntegration teste l'intégration avec WorkingMemory
func TestWorkingMemoryIntegration(t *testing.T) {
	storage := NewMemoryStorage()

	// Créer une mémoire de travail complexe
	complexFact := &Fact{
		ID:   "complex_fact",
		Type: "ComplexType",
		Fields: map[string]interface{}{
			"name":   "ComplexObject",
			"value":  3.14159,
			"active": true,
			"tags":   []string{"important", "test"},
			"metadata": map[string]interface{}{
				"created": "2023-01-01",
				"author":  "test",
			},
		},
		Timestamp: time.Now(),
	}

	workingMemory := &WorkingMemory{
		NodeID: "complex_node",
		Facts:  map[string]*Fact{"complex_fact": complexFact},
		Tokens: make(map[string]*Token),
	}

	// Sauvegarder
	err := storage.SaveMemory("complex_node", workingMemory)
	if err != nil {
		t.Fatalf("Erreur lors de la sauvegarde de la mémoire complexe: %v", err)
	}

	// Charger et vérifier
	loaded, err := storage.LoadMemory("complex_node")
	if err != nil {
		t.Fatalf("Erreur lors du chargement de la mémoire complexe: %v", err)
	}

	// Vérifier que les données sont identiques
	if len(loaded.Facts) != 1 {
		t.Fatalf("Attendu 1 fait complexe, obtenu %d", len(loaded.Facts))
	}

	// Récupérer le fait par son ID
	fact, exists := loaded.Facts["complex_fact"]
	if !exists {
		t.Fatal("Le fait complex_fact n'a pas été trouvé")
	}

	if fact.ID != "complex_fact" {
		t.Fatalf("ID du fait incorrect: attendu 'complex_fact', obtenu '%s'", fact.ID)
	}

	// Vérifier quelques champs spécifiques
	if fact.Fields["name"] != "ComplexObject" {
		t.Fatalf("Nom incorrect: attendu 'ComplexObject', obtenu '%v'", fact.Fields["name"])
	}

	if fact.Fields["value"] != 3.14159 {
		t.Fatalf("Valeur incorrecte: attendu 3.14159, obtenu '%v'", fact.Fields["value"])
	}

	if fact.Fields["active"] != true {
		t.Fatalf("Statut actif incorrect: attendu true, obtenu '%v'", fact.Fields["active"])
	}
}
