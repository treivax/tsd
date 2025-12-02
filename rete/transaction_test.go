// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

// Helper: Build a test network with N facts
func buildTestNetworkWithFacts(size int) *ReteNetwork {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Add a simple type
	typeDef := TypeDefinition{
		Name: "TestType",
		Fields: []Field{
			{Name: "value", Type: "int"},
			{Name: "name", Type: "string"},
		},
	}
	network.Types = append(network.Types, typeDef)

	// Add facts
	for i := 0; i < size; i++ {
		fact := &Fact{
			ID:   fmt.Sprintf("fact_%d", i),
			Type: "TestType",
			Fields: map[string]interface{}{
				"value": i,
				"name":  fmt.Sprintf("test_%d", i),
			},
			Timestamp: time.Now(),
		}
		storage.AddFact(fact)
	}

	return network
}

// Helper: Create a unique test fact
func createUniqueTestFact() *Fact {
	return &Fact{
		ID:   fmt.Sprintf("fact_%d", time.Now().UnixNano()),
		Type: "TestType",
		Fields: map[string]interface{}{
			"value": 42,
			"name":  "test",
		},
		Timestamp: time.Now(),
	}
}

// TestTransaction_CommitAppliesChanges tests that commit applies all changes
func TestTransaction_CommitAppliesChanges(t *testing.T) {
	t.Log("🔍 VALIDATION : Commit doit appliquer les changements")

	network := buildTestNetworkWithFacts(10)
	initialFactCount := len(network.Storage.GetAllFacts())

	tx := network.BeginTransaction()
	network.SetTransaction(tx)

	// Ajouter plusieurs faits
	facts := []*Fact{
		createUniqueTestFact(),
		createUniqueTestFact(),
		createUniqueTestFact(),
	}

	for _, fact := range facts {
		err := network.SubmitFact(fact)
		if err != nil {
			t.Fatalf("Failed to submit fact: %v", err)
		}
	}

	// Commit
	err := tx.Commit()
	if err != nil {
		t.Fatalf("Commit failed: %v", err)
	}

	network.SetTransaction(nil)

	// ✅ EXTRACTION DEPUIS RÉSEAU RÉEL
	actualCount := len(network.Storage.GetAllFacts())
	expectedCount := initialFactCount + len(facts)

	if actualCount != expectedCount {
		t.Errorf("Expected %d facts after commit, got %d", expectedCount, actualCount)
	}

	t.Logf("✅ Commit successful: %d facts added, total now %d", len(facts), actualCount)
}

// TestTransaction_RollbackRevertsAllChanges tests that rollback reverts all changes
func TestTransaction_RollbackRevertsAllChanges(t *testing.T) {
	t.Log("🔍 VALIDATION : Rollback doit annuler TOUTES les modifications")

	network := buildTestNetworkWithFacts(10)

	// Capturer l'état initial
	beforeFacts := network.Storage.GetAllFacts()
	beforeCount := len(beforeFacts)
	beforeIDs := make(map[string]bool)
	for _, fact := range beforeFacts {
		beforeIDs[fact.GetInternalID()] = true
	}

	tx := network.BeginTransaction()
	network.SetTransaction(tx)

	// Multiple operations
	addedFacts := make([]*Fact, 0, 10)
	for i := 0; i < 10; i++ {
		fact := createUniqueTestFact()
		addedFacts = append(addedFacts, fact)
		err := network.SubmitFact(fact)
		if err != nil {
			t.Fatalf("Failed to submit fact %d: %v", i, err)
		}
	}

	// Vérifier que les faits ont été ajoutés
	duringCount := len(network.Storage.GetAllFacts())
	if duringCount != beforeCount+10 {
		t.Errorf("Expected %d facts during transaction, got %d", beforeCount+10, duringCount)
	}

	// Rollback
	err := tx.Rollback()
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	network.SetTransaction(nil)

	// ✅ VÉRIFICATION STRICTE
	afterFacts := network.Storage.GetAllFacts()
	afterCount := len(afterFacts)

	if afterCount != beforeCount {
		t.Errorf("Rollback failed: expected %d facts, got %d", beforeCount, afterCount)
	}

	// Vérifier que les faits sont identiques (pas juste le count)
	for _, afterFact := range afterFacts {
		if !beforeIDs[afterFact.GetInternalID()] {
			t.Errorf("Unexpected fact %s after rollback", afterFact.GetInternalID())
		}
	}

	for _, beforeFact := range beforeFacts {
		found := false
		for _, afterFact := range afterFacts {
			if beforeFact.GetInternalID() == afterFact.GetInternalID() {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Fact %s lost after rollback", beforeFact.GetInternalID())
		}
	}

	// Vérifier que les faits ajoutés ont bien été supprimés
	for _, addedFact := range addedFacts {
		if network.Storage.GetFact(addedFact.GetInternalID()) != nil {
			t.Errorf("Added fact %s still present after rollback", addedFact.GetInternalID())
		}
	}

	t.Logf("✅ Rollback successful: state restored to %d facts", beforeCount)
}

// TestTransaction_MultipleOperations tests transactions with many operations
func TestTransaction_MultipleOperations(t *testing.T) {
	t.Log("🔍 TEST : Transaction avec multiples opérations")

	network := buildTestNetworkWithFacts(5)
	initialCount := len(network.Storage.GetAllFacts())

	tx := network.BeginTransaction()
	network.SetTransaction(tx)

	// Ajouter 50 faits
	for i := 0; i < 50; i++ {
		fact := createUniqueTestFact()
		if err := network.SubmitFact(fact); err != nil {
			t.Fatalf("Failed to submit fact %d: %v", i, err)
		}
	}

	// Vérifier le nombre de commandes
	commandCount := tx.GetCommandCount()
	if commandCount != 50 {
		t.Errorf("Expected 50 commands, got %d", commandCount)
	}

	// Commit
	if err := tx.Commit(); err != nil {
		t.Fatalf("Commit failed: %v", err)
	}

	network.SetTransaction(nil)

	// Vérifier le résultat
	finalCount := len(network.Storage.GetAllFacts())
	if finalCount != initialCount+50 {
		t.Errorf("Expected %d facts, got %d", initialCount+50, finalCount)
	}

	t.Logf("✅ Transaction with %d operations successful", commandCount)
}

// TestTransaction_CannotCommitTwice tests that committing twice fails
func TestTransaction_CannotCommitTwice(t *testing.T) {
	t.Log("🔍 TEST : Impossible de commit deux fois")

	network := buildTestNetworkWithFacts(5)
	tx := network.BeginTransaction()
	network.SetTransaction(tx)

	network.SubmitFact(createUniqueTestFact())

	// Premier commit
	if err := tx.Commit(); err != nil {
		t.Fatalf("First commit failed: %v", err)
	}

	// Deuxième commit (doit échouer)
	err := tx.Commit()
	if err == nil {
		t.Error("Expected error on second commit, got nil")
	}

	t.Log("✅ Cannot commit twice - correct behavior")
}

// TestTransaction_CannotRollbackAfterCommit tests that rollback after commit fails
func TestTransaction_CannotRollbackAfterCommit(t *testing.T) {
	t.Log("🔍 TEST : Impossible de rollback après commit")

	network := buildTestNetworkWithFacts(5)
	tx := network.BeginTransaction()
	network.SetTransaction(tx)

	network.SubmitFact(createUniqueTestFact())

	// Commit
	if err := tx.Commit(); err != nil {
		t.Fatalf("Commit failed: %v", err)
	}

	// Rollback (doit échouer)
	err := tx.Rollback()
	if err == nil {
		t.Error("Expected error on rollback after commit, got nil")
	}

	t.Log("✅ Cannot rollback after commit - correct behavior")
}

// TestTransaction_CannotRollbackTwice tests that rolling back twice fails
func TestTransaction_CannotRollbackTwice(t *testing.T) {
	t.Log("🔍 TEST : Impossible de rollback deux fois")

	network := buildTestNetworkWithFacts(5)
	tx := network.BeginTransaction()
	network.SetTransaction(tx)

	network.SubmitFact(createUniqueTestFact())

	// Premier rollback
	if err := tx.Rollback(); err != nil {
		t.Fatalf("First rollback failed: %v", err)
	}

	// Deuxième rollback (doit échouer)
	err := tx.Rollback()
	if err == nil {
		t.Error("Expected error on second rollback, got nil")
	}

	t.Log("✅ Cannot rollback twice - correct behavior")
}

// TestTransaction_EmptyTransaction tests an empty transaction
func TestTransaction_EmptyTransaction(t *testing.T) {
	t.Log("🔍 TEST : Transaction vide")

	network := buildTestNetworkWithFacts(5)
	initialCount := len(network.Storage.GetAllFacts())

	tx := network.BeginTransaction()
	network.SetTransaction(tx)

	// Pas d'opérations

	// Commit
	if err := tx.Commit(); err != nil {
		t.Fatalf("Commit failed: %v", err)
	}

	network.SetTransaction(nil)

	// Vérifier que rien n'a changé
	finalCount := len(network.Storage.GetAllFacts())
	if finalCount != initialCount {
		t.Errorf("Expected %d facts, got %d", initialCount, finalCount)
	}

	t.Log("✅ Empty transaction handled correctly")
}

// TestTransaction_GetCommandCount tests command counting
func TestTransaction_GetCommandCount(t *testing.T) {
	t.Log("🔍 TEST : Comptage des commandes")

	network := buildTestNetworkWithFacts(5)
	tx := network.BeginTransaction()
	network.SetTransaction(tx)

	if tx.GetCommandCount() != 0 {
		t.Errorf("Expected 0 commands initially, got %d", tx.GetCommandCount())
	}

	// Ajouter 3 faits
	for i := 0; i < 3; i++ {
		network.SubmitFact(createUniqueTestFact())
	}

	if tx.GetCommandCount() != 3 {
		t.Errorf("Expected 3 commands, got %d", tx.GetCommandCount())
	}

	tx.Commit()

	// Après commit, les commandes sont libérées
	if tx.GetCommandCount() != 0 {
		t.Errorf("Expected 0 commands after commit, got %d", tx.GetCommandCount())
	}

	t.Log("✅ Command counting works correctly")
}

// TestTransaction_GetDuration tests duration tracking
func TestTransaction_GetDuration(t *testing.T) {
	t.Log("🔍 TEST : Durée de transaction")

	network := buildTestNetworkWithFacts(5)
	tx := network.BeginTransaction()

	// Attendre un peu
	time.Sleep(10 * time.Millisecond)

	duration := tx.GetDuration()
	if duration < 10*time.Millisecond {
		t.Errorf("Expected duration >= 10ms, got %v", duration)
	}

	t.Logf("✅ Transaction duration: %v", duration)
}

// TestTransaction_String tests string representation
func TestTransaction_String(t *testing.T) {
	t.Log("🔍 TEST : Représentation string")

	network := buildTestNetworkWithFacts(5)
	tx := network.BeginTransaction()
	network.SetTransaction(tx)

	network.SubmitFact(createUniqueTestFact())

	str := tx.String()
	if str == "" {
		t.Error("Transaction.String() returned empty string")
	}

	t.Logf("Transaction: %s", str)

	tx.Commit()
	str = tx.String()
	t.Logf("After commit: %s", str)

	t.Log("✅ String representation works")
}

// TestTransaction_MemoryFootprint tests memory footprint estimation
func TestTransaction_MemoryFootprint(t *testing.T) {
	t.Log("🔍 TEST : Empreinte mémoire")

	network := buildTestNetworkWithFacts(100)

	// Mesurer mémoire avant transaction
	runtime.GC()
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	tx := network.BeginTransaction()
	network.SetTransaction(tx)

	// Ajouter quelques faits
	for i := 0; i < 10; i++ {
		network.SubmitFact(createUniqueTestFact())
	}

	// Mesurer mémoire après
	runtime.GC()
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	actualOverhead := int64(m2.Alloc - m1.Alloc)
	estimatedFootprint := tx.GetMemoryFootprint()

	t.Logf("Actual memory overhead: %d bytes", actualOverhead)
	t.Logf("Estimated footprint: %d bytes", estimatedFootprint)

	// L'estimation devrait être dans le bon ordre de grandeur
	// (pas nécessairement exacte, mais raisonnable)
	if estimatedFootprint == 0 {
		t.Error("Memory footprint should not be zero")
	}

	tx.Commit()
	network.SetTransaction(nil)

	t.Log("✅ Memory footprint estimation works")
}

// TestTransaction_WithoutNetwork tests transaction without setting on network
func TestTransaction_WithoutNetwork(t *testing.T) {
	t.Log("🔍 TEST : Transaction sans réseau")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	tx := network.BeginTransaction()
	// Ne PAS appeler network.SetTransaction(tx)

	initialCount := len(storage.GetAllFacts())

	// Soumettre un fait (devrait aller en mode normal, pas transactionnel)
	fact := createUniqueTestFact()
	network.SubmitFact(fact)

	// Le fait devrait être ajouté directement (pas via transaction)
	afterCount := len(storage.GetAllFacts())
	if afterCount != initialCount+1 {
		t.Errorf("Expected fact to be added directly, got %d facts", afterCount)
	}

	// La transaction ne devrait avoir aucune commande
	if tx.GetCommandCount() != 0 {
		t.Errorf("Expected 0 commands in unused transaction, got %d", tx.GetCommandCount())
	}

	t.Log("✅ Transaction without network works correctly")
}

// TestTransaction_ConcurrentAccess tests basic thread safety
func TestTransaction_ConcurrentAccess(t *testing.T) {
	t.Log("🔍 TEST : Accès concurrent (basique)")

	network := buildTestNetworkWithFacts(10)
	tx := network.BeginTransaction()
	network.SetTransaction(tx)

	done := make(chan bool, 2)

	// Goroutine 1: ajouter des faits
	go func() {
		for i := 0; i < 5; i++ {
			network.SubmitFact(createUniqueTestFact())
			time.Sleep(1 * time.Millisecond)
		}
		done <- true
	}()

	// Goroutine 2: lire le nombre de commandes
	go func() {
		for i := 0; i < 5; i++ {
			_ = tx.GetCommandCount()
			time.Sleep(1 * time.Millisecond)
		}
		done <- true
	}()

	// Attendre la fin
	<-done
	<-done

	tx.Commit()
	network.SetTransaction(nil)

	t.Log("✅ Basic concurrent access works (no panic)")
}

// TestTransaction_BeginTransactionIsO1 tests that BeginTransaction is O(1)
func TestTransaction_BeginTransactionIsO1(t *testing.T) {
	t.Log("🔍 TEST : BeginTransaction est O(1)")

	sizes := []int{100, 1000, 10000, 100000}
	times := make([]time.Duration, len(sizes))

	for i, size := range sizes {
		network := buildTestNetworkWithFacts(size)

		start := time.Now()
		tx := network.BeginTransaction()
		times[i] = time.Since(start)

		tx.Commit()

		t.Logf("Size %d: BeginTransaction took %v", size, times[i])
	}

	// Vérifier que le temps ne croît pas linéairement
	// Le temps pour 100k faits ne devrait pas être 1000x le temps pour 100 faits
	ratio := float64(times[len(times)-1]) / float64(times[0])
	expectedMaxRatio := 100.0 // Tolérance généreuse pour variations système

	if ratio > expectedMaxRatio {
		t.Errorf("❌ BeginTransaction n'est pas O(1): ratio=%.2f (max=%.2f)", ratio, expectedMaxRatio)
	} else {
		t.Logf("✅ BeginTransaction est O(1): ratio=%.2f", ratio)
	}
}

// TestTransaction_GetCommands tests getting commands list
func TestTransaction_GetCommands(t *testing.T) {
	t.Log("🔍 TEST : Récupération de la liste des commandes")

	network := buildTestNetworkWithFacts(5)
	tx := network.BeginTransaction()
	network.SetTransaction(tx)

	// Ajouter 3 faits
	for i := 0; i < 3; i++ {
		network.SubmitFact(createUniqueTestFact())
	}

	commands := tx.GetCommands()
	if len(commands) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(commands))
	}

	// Vérifier que c'est une copie (pas la référence interne)
	commands[0] = nil
	if tx.GetCommandCount() != 3 {
		t.Error("Modifying returned commands affected internal state")
	}

	tx.Commit()

	// Après commit, GetCommands devrait retourner nil
	commands = tx.GetCommands()
	if commands != nil {
		t.Error("Expected nil commands after commit")
	}

	t.Log("✅ GetCommands works correctly")
}
