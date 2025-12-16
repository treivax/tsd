// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"
)

func TestNewMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()
	if storage == nil {
		t.Fatal("NewMemoryStorage returned nil")
	}
	if storage.memories == nil {
		t.Error("memories map not initialized")
	}
	if len(storage.memories) != 0 {
		t.Errorf("Expected empty memories map, got %d entries", len(storage.memories))
	}
}
func TestSaveMemory(t *testing.T) {
	storage := NewMemoryStorage()
	memory := &WorkingMemory{
		NodeID: "node1",
		Facts:  make(map[string]*Fact),
		Tokens: map[string]*Token{
			"token1": {
				Facts: []*Fact{
					{
						ID:     "fact1",
						Type:   "Person",
						Fields: map[string]interface{}{"name": "Alice"},
					},
				},
			},
		},
	}
	err := storage.SaveMemory("node1", memory)
	if err != nil {
		t.Fatalf("SaveMemory failed: %v", err)
	}
	// Verify memory was saved
	if len(storage.memories) != 1 {
		t.Errorf("Expected 1 memory entry, got %d", len(storage.memories))
	}
	saved, exists := storage.memories["node1"]
	if !exists {
		t.Fatal("Memory not found for node1")
	}
	if len(saved.Tokens) != 1 {
		t.Errorf("Expected 1 token in saved memory, got %d", len(saved.Tokens))
	}
}
func TestSaveMemoryCreatesDeepCopy(t *testing.T) {
	storage := NewMemoryStorage()
	memory := &WorkingMemory{
		NodeID: "node1",
		Facts:  make(map[string]*Fact),
		Tokens: map[string]*Token{
			"token1": {
				Facts: []*Fact{
					{
						ID:     "fact1",
						Type:   "Person",
						Fields: map[string]interface{}{"name": "Alice"},
					},
				},
			},
		},
	}
	err := storage.SaveMemory("node1", memory)
	if err != nil {
		t.Fatalf("SaveMemory failed: %v", err)
	}
	// Modify original memory
	memory.Tokens["token2"] = &Token{
		Facts: []*Fact{
			{
				ID:     "fact2",
				Type:   "Person",
				Fields: map[string]interface{}{"name": "Bob"},
			},
		},
	}
	// Verify saved memory wasn't modified
	saved := storage.memories["node1"]
	if len(saved.Tokens) != 1 {
		t.Errorf("Expected saved memory to have 1 token (not affected by original modification), got %d", len(saved.Tokens))
	}
}
func TestLoadMemory(t *testing.T) {
	storage := NewMemoryStorage()
	memory := &WorkingMemory{
		NodeID: "node1",
		Facts:  make(map[string]*Fact),
		Tokens: map[string]*Token{
			"token1": {
				Facts: []*Fact{
					{
						ID:     "fact1",
						Type:   "Person",
						Fields: map[string]interface{}{"name": "Alice", "age": 30},
					},
				},
			},
		},
	}
	storage.SaveMemory("node1", memory)
	// Load memory
	loaded, err := storage.LoadMemory("node1")
	if err != nil {
		t.Fatalf("LoadMemory failed: %v", err)
	}
	if loaded == nil {
		t.Fatal("LoadMemory returned nil")
	}
	if len(loaded.Tokens) != 1 {
		t.Errorf("Expected 1 token in loaded memory, got %d", len(loaded.Tokens))
	}
	token, exists := loaded.Tokens["token1"]
	if !exists {
		t.Fatal("Expected token1 in loaded memory")
	}
	if len(token.Facts) != 1 {
		t.Errorf("Expected 1 fact in token, got %d", len(token.Facts))
	}
	fact := token.Facts[0]
	if fact.ID != "fact1" {
		t.Errorf("Expected fact ID 'fact1', got '%s'", fact.ID)
	}
	if fact.Type != "Person" {
		t.Errorf("Expected fact type 'Person', got '%s'", fact.Type)
	}
	if fact.Fields["name"] != "Alice" {
		t.Errorf("Expected name 'Alice', got '%v'", fact.Fields["name"])
	}
}
func TestLoadMemoryNonExistent(t *testing.T) {
	storage := NewMemoryStorage()
	loaded, err := storage.LoadMemory("nonexistent")
	if err == nil {
		t.Error("Expected error when loading non-existent memory")
	}
	if loaded != nil {
		t.Error("Expected nil memory for non-existent node")
	}
	expectedErrMsg := "mémoire non trouvée pour le nœud nonexistent"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}
func TestLoadMemoryReturnsDeepCopy(t *testing.T) {
	storage := NewMemoryStorage()
	memory := &WorkingMemory{
		NodeID: "node1",
		Facts:  make(map[string]*Fact),
		Tokens: map[string]*Token{
			"token1": {
				Facts: []*Fact{
					{
						ID:     "fact1",
						Type:   "Person",
						Fields: map[string]interface{}{"name": "Alice"},
					},
				},
			},
		},
	}
	storage.SaveMemory("node1", memory)
	// Load memory twice
	loaded1, _ := storage.LoadMemory("node1")
	loaded2, _ := storage.LoadMemory("node1")
	// Modify first loaded memory
	loaded1.Tokens["token2"] = &Token{
		Facts: []*Fact{
			{
				ID:     "fact2",
				Type:   "Car",
				Fields: map[string]interface{}{"model": "Tesla"},
			},
		},
	}
	// Verify second loaded memory is not affected
	if len(loaded2.Tokens) != 1 {
		t.Errorf("Expected loaded2 to have 1 token (independent copy), got %d", len(loaded2.Tokens))
	}
	// Verify stored memory is not affected
	stored := storage.memories["node1"]
	if len(stored.Tokens) != 1 {
		t.Errorf("Expected stored memory to have 1 token (not affected), got %d", len(stored.Tokens))
	}
}
func TestDeleteMemory(t *testing.T) {
	storage := NewMemoryStorage()
	memory := &WorkingMemory{
		NodeID: "test",
		Facts:  make(map[string]*Fact),
		Tokens: map[string]*Token{"t1": {Facts: []*Fact{{ID: "fact1", Type: "Test"}}}},
	}
	storage.SaveMemory("node1", memory)
	storage.SaveMemory("node2", memory)
	if len(storage.memories) != 2 {
		t.Errorf("Expected 2 memory entries before delete, got %d", len(storage.memories))
	}
	// Delete one memory
	err := storage.DeleteMemory("node1")
	if err != nil {
		t.Fatalf("DeleteMemory failed: %v", err)
	}
	// Verify memory was deleted
	if len(storage.memories) != 1 {
		t.Errorf("Expected 1 memory entry after delete, got %d", len(storage.memories))
	}
	_, exists := storage.memories["node1"]
	if exists {
		t.Error("node1 memory should have been deleted")
	}
	_, exists = storage.memories["node2"]
	if !exists {
		t.Error("node2 memory should still exist")
	}
}
func TestDeleteMemoryNonExistent(t *testing.T) {
	storage := NewMemoryStorage()
	// Deleting non-existent memory should not error
	err := storage.DeleteMemory("nonexistent")
	if err != nil {
		t.Errorf("DeleteMemory should not error on non-existent node, got: %v", err)
	}
}
func TestListNodes(t *testing.T) {
	storage := NewMemoryStorage()
	// Test empty storage
	nodes, err := storage.ListNodes()
	if err != nil {
		t.Fatalf("ListNodes failed: %v", err)
	}
	if len(nodes) != 0 {
		t.Errorf("Expected 0 nodes in empty storage, got %d", len(nodes))
	}
	// Add some memories
	memory := &WorkingMemory{
		NodeID: "test",
		Facts:  make(map[string]*Fact),
		Tokens: map[string]*Token{"t1": {Facts: []*Fact{{ID: "fact1", Type: "Test"}}}},
	}
	storage.SaveMemory("node1", memory)
	storage.SaveMemory("node2", memory)
	storage.SaveMemory("node3", memory)
	// List nodes
	nodes, err = storage.ListNodes()
	if err != nil {
		t.Fatalf("ListNodes failed: %v", err)
	}
	if len(nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(nodes))
	}
	// Verify all node IDs are present
	nodeMap := make(map[string]bool)
	for _, node := range nodes {
		nodeMap[node] = true
	}
	expectedNodes := []string{"node1", "node2", "node3"}
	for _, expected := range expectedNodes {
		if !nodeMap[expected] {
			t.Errorf("Expected node '%s' in list, but it was not found", expected)
		}
	}
}
func TestListNodesAfterDelete(t *testing.T) {
	storage := NewMemoryStorage()
	memory := &WorkingMemory{
		NodeID: "test",
		Facts:  make(map[string]*Fact),
		Tokens: map[string]*Token{"t1": {Facts: []*Fact{{ID: "fact1", Type: "Test"}}}},
	}
	storage.SaveMemory("node1", memory)
	storage.SaveMemory("node2", memory)
	storage.SaveMemory("node3", memory)
	// Delete one node
	storage.DeleteMemory("node2")
	// List nodes
	nodes, err := storage.ListNodes()
	if err != nil {
		t.Fatalf("ListNodes failed: %v", err)
	}
	if len(nodes) != 2 {
		t.Errorf("Expected 2 nodes after deletion, got %d", len(nodes))
	}
	// Verify node2 is not in the list
	for _, node := range nodes {
		if node == "node2" {
			t.Error("node2 should not be in the list after deletion")
		}
	}
}
func TestConcurrentMemoryAccess(t *testing.T) {
	storage := NewMemoryStorage()
	memory := &WorkingMemory{
		NodeID: "test",
		Facts:  make(map[string]*Fact),
		Tokens: map[string]*Token{"t1": {Facts: []*Fact{{ID: "fact1", Type: "Test"}}}},
	}
	done := make(chan bool)
	// Concurrent saves
	for i := 0; i < 10; i++ {
		go func(id int) {
			nodeID := string(rune('a' + id))
			storage.SaveMemory(nodeID, memory)
			done <- true
		}(i)
	}
	// Wait for all saves
	for i := 0; i < 10; i++ {
		<-done
	}
	// Verify all memories were saved
	nodes, _ := storage.ListNodes()
	if len(nodes) != 10 {
		t.Errorf("Expected 10 nodes after concurrent saves, got %d", len(nodes))
	}
	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func(id int) {
			nodeID := string(rune('a' + id))
			storage.LoadMemory(nodeID)
			done <- true
		}(i)
	}
	// Wait for all reads
	for i := 0; i < 10; i++ {
		<-done
	}
}
func TestSaveAndLoadComplexMemory(t *testing.T) {
	storage := NewMemoryStorage()
	// Create complex memory with multiple tokens and facts
	memory := &WorkingMemory{
		NodeID: "complex",
		Facts:  make(map[string]*Fact),
		Tokens: map[string]*Token{
			"token1": {
				Facts: []*Fact{
					{
						ID:   "fact1",
						Type: "Person",
						Fields: map[string]interface{}{
							"name":    "Alice",
							"age":     30,
							"city":    "Paris",
							"hobbies": []interface{}{"reading", "coding"},
						},
					},
					{
						ID:   "fact2",
						Type: "Person",
						Fields: map[string]interface{}{
							"name": "Bob",
							"age":  25,
						},
					},
				},
			},
			"token2": {
				Facts: []*Fact{
					{
						ID:   "fact3",
						Type: "Car",
						Fields: map[string]interface{}{
							"model": "Tesla",
							"year":  2023,
						},
					},
				},
			},
		},
	}
	err := storage.SaveMemory("complex", memory)
	if err != nil {
		t.Fatalf("SaveMemory failed for complex memory: %v", err)
	}
	loaded, err := storage.LoadMemory("complex")
	if err != nil {
		t.Fatalf("LoadMemory failed: %v", err)
	}
	if len(loaded.Tokens) != 2 {
		t.Errorf("Expected 2 tokens, got %d", len(loaded.Tokens))
	}
	token1, exists := loaded.Tokens["token1"]
	if !exists {
		t.Fatal("Expected token1 in loaded memory")
	}
	if len(token1.Facts) != 2 {
		t.Errorf("Expected 2 facts in token1, got %d", len(token1.Facts))
	}
	token2, exists := loaded.Tokens["token2"]
	if !exists {
		t.Fatal("Expected token2 in loaded memory")
	}
	if len(token2.Facts) != 1 {
		t.Errorf("Expected 1 fact in token2, got %d", len(token2.Facts))
	}
}
func TestOverwriteMemory(t *testing.T) {
	storage := NewMemoryStorage()
	memory1 := &WorkingMemory{
		NodeID: "node1",
		Facts:  make(map[string]*Fact),
		Tokens: map[string]*Token{"t1": {Facts: []*Fact{{ID: "fact1", Type: "Test1"}}}},
	}
	memory2 := &WorkingMemory{
		NodeID: "node1",
		Facts:  make(map[string]*Fact),
		Tokens: map[string]*Token{
			"t2": {Facts: []*Fact{{ID: "fact2", Type: "Test2"}}},
			"t3": {Facts: []*Fact{{ID: "fact3", Type: "Test3"}}},
		},
	}
	// Save first memory
	storage.SaveMemory("node1", memory1)
	loaded, _ := storage.LoadMemory("node1")
	if len(loaded.Tokens) != 1 {
		t.Errorf("Expected 1 token after first save, got %d", len(loaded.Tokens))
	}
	// Overwrite with second memory
	storage.SaveMemory("node1", memory2)
	loaded, _ = storage.LoadMemory("node1")
	if len(loaded.Tokens) != 2 {
		t.Errorf("Expected 2 tokens after overwrite, got %d", len(loaded.Tokens))
	}
	token, exists := loaded.Tokens["t2"]
	if !exists {
		t.Fatal("Expected t2 token after overwrite")
	}
	if token.Facts[0].ID != "fact2" {
		t.Errorf("Expected fact2 after overwrite, got %s", token.Facts[0].ID)
	}
}

func TestSync(t *testing.T) {
	t.Log("🧪 TEST Sync")
	t.Log("============")

	t.Run("sync avec storage vide", func(t *testing.T) {
		storage := NewMemoryStorage()

		err := storage.Sync()

		if err != nil {
			t.Errorf("❌ Sync() sur storage vide ne devrait pas échouer: %v", err)
			return
		}

		t.Log("✅ Sync() sur storage vide réussi")
	})

	t.Run("sync avec mémoires valides", func(t *testing.T) {
		storage := NewMemoryStorage()

		// Ajouter plusieurs mémoires valides
		for i := 1; i <= 3; i++ {
			memory := &WorkingMemory{
				NodeID: "node" + string(rune('0'+i)),
				Facts: map[string]*Fact{
					"fact1": {ID: "fact1", Type: "Test"},
				},
				Tokens: map[string]*Token{
					"token1": {Facts: []*Fact{{ID: "fact1", Type: "Test"}}},
				},
			}
			storage.SaveMemory("node"+string(rune('0'+i)), memory)
		}

		err := storage.Sync()

		if err != nil {
			t.Errorf("❌ Sync() ne devrait pas échouer avec mémoires valides: %v", err)
			return
		}

		t.Log("✅ Sync() avec mémoires valides réussi")
	})

	t.Run("sync initialise Facts nil", func(t *testing.T) {
		storage := NewMemoryStorage()

		// Créer une mémoire avec Facts nil
		memory := &WorkingMemory{
			NodeID: "node1",
			Facts:  nil, // Facts est nil
			Tokens: map[string]*Token{
				"token1": {Facts: []*Fact{{ID: "fact1", Type: "Test"}}},
			},
		}

		// Sauvegarder directement dans le storage (bypass SaveMemory qui fait une copie)
		storage.mutex.Lock()
		storage.memories["node1"] = memory
		storage.mutex.Unlock()

		err := storage.Sync()

		if err != nil {
			t.Errorf("❌ Sync() ne devrait pas échouer: %v", err)
			return
		}

		// Vérifier que Facts a été initialisé
		storage.mutex.RLock()
		mem := storage.memories["node1"]
		storage.mutex.RUnlock()

		if mem.Facts == nil {
			t.Error("❌ Sync() devrait initialiser Facts à une map vide")
			return
		}

		if len(mem.Facts) != 0 {
			t.Errorf("❌ Facts devrait être vide, longueur: %d", len(mem.Facts))
			return
		}

		t.Log("✅ Sync() a correctement initialisé Facts nil")
	})

	t.Run("sync initialise Tokens nil", func(t *testing.T) {
		storage := NewMemoryStorage()

		// Créer une mémoire avec Tokens nil
		memory := &WorkingMemory{
			NodeID: "node1",
			Facts: map[string]*Fact{
				"fact1": {ID: "fact1", Type: "Test"},
			},
			Tokens: nil, // Tokens est nil
		}

		// Sauvegarder directement dans le storage
		storage.mutex.Lock()
		storage.memories["node1"] = memory
		storage.mutex.Unlock()

		err := storage.Sync()

		if err != nil {
			t.Errorf("❌ Sync() ne devrait pas échouer: %v", err)
			return
		}

		// Vérifier que Tokens a été initialisé
		storage.mutex.RLock()
		mem := storage.memories["node1"]
		storage.mutex.RUnlock()

		if mem.Tokens == nil {
			t.Error("❌ Sync() devrait initialiser Tokens à une map vide")
			return
		}

		if len(mem.Tokens) != 0 {
			t.Errorf("❌ Tokens devrait être vide, longueur: %d", len(mem.Tokens))
			return
		}

		t.Log("✅ Sync() a correctement initialisé Tokens nil")
	})

	t.Run("sync détecte mémoire nulle", func(t *testing.T) {
		storage := NewMemoryStorage()

		// Insérer une mémoire nulle directement
		storage.mutex.Lock()
		storage.memories["node1"] = nil
		storage.mutex.Unlock()

		err := storage.Sync()

		if err == nil {
			t.Error("❌ Sync() devrait échouer avec une mémoire nulle")
			return
		}

		expectedErrMsg := "mémoire nulle pour le nœud node1"
		if err.Error() != expectedErrMsg {
			t.Errorf("❌ Message d'erreur incorrect: attendu %q, reçu %q", expectedErrMsg, err.Error())
			return
		}

		t.Log("✅ Sync() a correctement détecté la mémoire nulle")
	})

	t.Run("sync initialise Facts et Tokens ensemble", func(t *testing.T) {
		storage := NewMemoryStorage()

		// Créer une mémoire avec Facts et Tokens nil
		memory := &WorkingMemory{
			NodeID: "node1",
			Facts:  nil,
			Tokens: nil,
		}

		storage.mutex.Lock()
		storage.memories["node1"] = memory
		storage.mutex.Unlock()

		err := storage.Sync()

		if err != nil {
			t.Errorf("❌ Sync() ne devrait pas échouer: %v", err)
			return
		}

		// Vérifier que les deux ont été initialisés
		storage.mutex.RLock()
		mem := storage.memories["node1"]
		storage.mutex.RUnlock()

		if mem.Facts == nil {
			t.Error("❌ Facts devrait être initialisé")
		}
		if mem.Tokens == nil {
			t.Error("❌ Tokens devrait être initialisé")
		}

		t.Log("✅ Sync() a initialisé Facts et Tokens ensemble")
	})

	t.Run("sync préserve données existantes", func(t *testing.T) {
		storage := NewMemoryStorage()

		// Créer une mémoire avec données existantes
		memory := &WorkingMemory{
			NodeID: "node1",
			Facts: map[string]*Fact{
				"fact1": {ID: "fact1", Type: "Person", Fields: map[string]interface{}{"name": "Alice"}},
				"fact2": {ID: "fact2", Type: "Car", Fields: map[string]interface{}{"model": "Tesla"}},
			},
			Tokens: map[string]*Token{
				"token1": {Facts: []*Fact{{ID: "fact1", Type: "Person"}}},
				"token2": {Facts: []*Fact{{ID: "fact2", Type: "Car"}}},
			},
		}

		storage.mutex.Lock()
		storage.memories["node1"] = memory
		storage.mutex.Unlock()

		err := storage.Sync()

		if err != nil {
			t.Errorf("❌ Sync() ne devrait pas échouer: %v", err)
			return
		}

		// Vérifier que les données sont préservées
		storage.mutex.RLock()
		mem := storage.memories["node1"]
		storage.mutex.RUnlock()

		if len(mem.Facts) != 2 {
			t.Errorf("❌ Facts devrait contenir 2 éléments, reçu %d", len(mem.Facts))
		}
		if len(mem.Tokens) != 2 {
			t.Errorf("❌ Tokens devrait contenir 2 éléments, reçu %d", len(mem.Tokens))
		}

		t.Log("✅ Sync() a préservé les données existantes")
	})
}

func TestSyncConcurrency(t *testing.T) {
	t.Log("🧪 TEST Sync - Concurrence")
	t.Log("==========================")

	storage := NewMemoryStorage()

	// Ajouter plusieurs mémoires
	for i := 0; i < 5; i++ {
		memory := &WorkingMemory{
			NodeID: "node" + string(rune('0'+i)),
			Facts: map[string]*Fact{
				"fact1": {ID: "fact1", Type: "Test"},
			},
			Tokens: map[string]*Token{
				"token1": {Facts: []*Fact{{ID: "fact1", Type: "Test"}}},
			},
		}
		storage.SaveMemory("node"+string(rune('0'+i)), memory)
	}

	// Appeler Sync() en parallèle plusieurs fois
	done := make(chan error, 10)
	for i := 0; i < 10; i++ {
		go func() {
			done <- storage.Sync()
		}()
	}

	// Attendre tous les appels
	for i := 0; i < 10; i++ {
		err := <-done
		if err != nil {
			t.Errorf("❌ Sync() concurrent a échoué: %v", err)
		}
	}

	t.Log("✅ Sync() concurrent réussi sans race condition")
}

func TestClearMemoryStorage(t *testing.T) {
	t.Log("🧪 TEST Clear - MemoryStorage")
	t.Log("=============================")

	storage := NewMemoryStorage()

	// Ajouter des mémoires
	for i := 0; i < 5; i++ {
		memory := &WorkingMemory{
			NodeID: "node" + string(rune('0'+i)),
			Facts:  map[string]*Fact{"fact1": {ID: "fact1", Type: "Test"}},
			Tokens: map[string]*Token{"token1": {Facts: []*Fact{{ID: "fact1"}}}},
		}
		storage.SaveMemory("node"+string(rune('0'+i)), memory)
	}

	if len(storage.memories) != 5 {
		t.Errorf("❌ Attendu 5 mémoires avant Clear, reçu %d", len(storage.memories))
	}

	// Clear
	err := storage.Clear()
	if err != nil {
		t.Errorf("❌ Clear() ne devrait pas échouer: %v", err)
		return
	}

	// Vérifier que tout est vide
	if len(storage.memories) != 0 {
		t.Errorf("❌ Attendu 0 mémoires après Clear, reçu %d", len(storage.memories))
		return
	}

	nodes, _ := storage.ListNodes()
	if len(nodes) != 0 {
		t.Errorf("❌ Attendu 0 nodes après Clear, reçu %d", len(nodes))
	}

	t.Log("✅ Clear() a correctement vidé le storage")
}
