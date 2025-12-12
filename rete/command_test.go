// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"
)

// TestAddFactCommand_ExecuteUndo tests the AddFact command execution and undo
func TestAddFactCommand_ExecuteUndo(t *testing.T) {
	t.Log("🔍 TEST : AddFactCommand Execute + Undo")
	storage := NewMemoryStorage()
	fact := &Fact{
		ID:   "test1",
		Type: "TestType",
		Fields: map[string]interface{}{
			"value": 42,
		},
	}
	cmd := NewAddFactCommand(storage, fact)
	// Test Execute
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	// Vérifier que le fait a été ajouté
	factID := fact.GetInternalID()
	retrieved := storage.GetFact(factID)
	if retrieved == nil {
		t.Error("Fact not found after Execute")
	}
	if retrieved != nil && retrieved.ID != "test1" {
		t.Errorf("Expected fact ID 'test1', got '%s'", retrieved.ID)
	}
	// Test Undo
	err = cmd.Undo()
	if err != nil {
		t.Fatalf("Undo failed: %v", err)
	}
	// Vérifier que le fait a été supprimé
	retrieved = storage.GetFact(factID)
	if retrieved != nil {
		t.Error("Fact still present after Undo")
	}
	t.Log("✅ AddFactCommand Execute + Undo successful")
}

// TestAddFactCommand_Idempotence tests that Execute+Undo+Execute produces same result
func TestAddFactCommand_Idempotence(t *testing.T) {
	t.Log("🔍 TEST : AddFactCommand idempotence")
	storage := NewMemoryStorage()
	fact := &Fact{
		ID:   "test2",
		Type: "TestType",
		Fields: map[string]interface{}{
			"value": 100,
		},
	}
	cmd := NewAddFactCommand(storage, fact)
	factID := fact.GetInternalID()
	// Execute -> Undo -> Execute
	cmd.Execute()
	cmd.Undo()
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Second Execute failed: %v", err)
	}
	// Vérifier état final
	retrieved := storage.GetFact(factID)
	if retrieved == nil {
		t.Error("Fact not found after second Execute")
	}
	t.Log("✅ AddFactCommand is idempotent")
}

// TestRemoveFactCommand_ExecuteUndo tests the RemoveFact command
func TestRemoveFactCommand_ExecuteUndo(t *testing.T) {
	t.Log("🔍 TEST : RemoveFactCommand Execute + Undo")
	storage := NewMemoryStorage()
	fact := &Fact{
		ID:   "test3",
		Type: "TestType",
		Fields: map[string]interface{}{
			"value": 200,
		},
	}
	// Ajouter le fait d'abord
	storage.AddFact(fact)
	factID := fact.GetInternalID()
	cmd := NewRemoveFactCommand(storage, factID)
	// Test Execute
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	// Vérifier que le fait a été supprimé
	retrieved := storage.GetFact(factID)
	if retrieved != nil {
		t.Error("Fact still present after Execute")
	}
	// Test Undo (doit restaurer le fait)
	err = cmd.Undo()
	if err != nil {
		t.Fatalf("Undo failed: %v", err)
	}
	// Vérifier que le fait a été restauré
	retrieved = storage.GetFact(factID)
	if retrieved == nil {
		t.Error("Fact not restored after Undo")
	}
	if retrieved != nil && retrieved.ID != "test3" {
		t.Errorf("Expected fact ID 'test3', got '%s'", retrieved.ID)
	}
	if retrieved != nil && retrieved.Fields["value"] != 200 {
		t.Errorf("Expected value 200, got %v", retrieved.Fields["value"])
	}
	t.Log("✅ RemoveFactCommand Execute + Undo successful")
}

// TestRemoveFactCommand_NonExistentFact tests removing a non-existent fact
func TestRemoveFactCommand_NonExistentFact(t *testing.T) {
	t.Log("🔍 TEST : RemoveFactCommand sur fait inexistant")
	storage := NewMemoryStorage()
	cmd := NewRemoveFactCommand(storage, "TestType_nonexistent")
	// Execute devrait échouer
	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error when removing non-existent fact, got nil")
	}
	t.Log("✅ RemoveFactCommand correctly fails on non-existent fact")
}

// TestCommandString tests the String() method of commands
func TestCommandString(t *testing.T) {
	t.Log("🔍 TEST : Command String() representation")
	storage := NewMemoryStorage()
	fact := &Fact{
		ID:   "test4",
		Type: "TestType",
		Fields: map[string]interface{}{
			"value": 42,
		},
	}
	addCmd := NewAddFactCommand(storage, fact)
	addStr := addCmd.String()
	if addStr == "" {
		t.Error("AddFactCommand.String() returned empty string")
	}
	t.Logf("AddFactCommand: %s", addStr)
	storage.AddFact(fact)
	removeCmd := NewRemoveFactCommand(storage, fact.GetInternalID())
	removeStr := removeCmd.String()
	if removeStr == "" {
		t.Error("RemoveFactCommand.String() returned empty string")
	}
	t.Logf("RemoveFactCommand: %s", removeStr)
	t.Log("✅ Command String() works correctly")
}

// TestMultipleCommands tests executing multiple commands in sequence
func TestMultipleCommands(t *testing.T) {
	t.Log("🔍 TEST : Multiple commands in sequence")
	storage := NewMemoryStorage()
	// Créer et exécuter plusieurs commandes
	facts := []*Fact{
		{ID: "f1", Type: "T1", Fields: map[string]interface{}{"v": 1}},
		{ID: "f2", Type: "T1", Fields: map[string]interface{}{"v": 2}},
		{ID: "f3", Type: "T1", Fields: map[string]interface{}{"v": 3}},
	}
	commands := make([]Command, len(facts))
	for i, fact := range facts {
		commands[i] = NewAddFactCommand(storage, fact)
		if err := commands[i].Execute(); err != nil {
			t.Fatalf("Command %d Execute failed: %v", i, err)
		}
	}
	// Vérifier que tous les faits sont présents
	allFacts := storage.GetAllFacts()
	if len(allFacts) != 3 {
		t.Errorf("Expected 3 facts, got %d", len(allFacts))
	}
	// Undo en ordre inverse
	for i := len(commands) - 1; i >= 0; i-- {
		if err := commands[i].Undo(); err != nil {
			t.Fatalf("Command %d Undo failed: %v", i, err)
		}
	}
	// Vérifier que tous les faits ont été supprimés
	allFacts = storage.GetAllFacts()
	if len(allFacts) != 0 {
		t.Errorf("Expected 0 facts after undo, got %d", len(allFacts))
	}
	t.Log("✅ Multiple commands executed and undone successfully")
}

// TestCommandError tests command error handling
func TestCommandError(t *testing.T) {
	t.Log("🔍 TEST : Command error handling")
	err := NewCommandError("TestCommand", "Execute", nil)
	if err.CommandName != "TestCommand" {
		t.Errorf("Expected CommandName 'TestCommand', got '%s'", err.CommandName)
	}
	if err.Operation != "Execute" {
		t.Errorf("Expected Operation 'Execute', got '%s'", err.Operation)
	}
	errStr := err.Error()
	if errStr == "" {
		t.Error("CommandError.Error() returned empty string")
	}
	t.Logf("CommandError: %s", errStr)
	t.Log("✅ Command error handling works correctly")
}

// TestCommandError_Unwrap tests the Unwrap method
func TestCommandError_Unwrap(t *testing.T) {
	t.Log("🔍 TEST : CommandError Unwrap")
	t.Run("unwrap with underlying error", func(t *testing.T) {
		underlyingErr := NewCommandError("InnerCommand", "Undo", nil)
		outerErr := NewCommandError("OuterCommand", "Execute", underlyingErr)
		unwrapped := outerErr.Unwrap()
		if unwrapped == nil {
			t.Fatal("Unwrap returned nil")
		}
		if unwrapped != underlyingErr {
			t.Error("Unwrap did not return the correct underlying error")
		}
	})
	t.Run("unwrap with nil error", func(t *testing.T) {
		err := NewCommandError("TestCommand", "Execute", nil)
		unwrapped := err.Unwrap()
		if unwrapped != nil {
			t.Errorf("Expected nil from Unwrap, got %v", unwrapped)
		}
	})
	t.Log("✅ CommandError Unwrap works correctly")
}

// BenchmarkAddFactCommand benchmarks AddFact command execution
func BenchmarkAddFactCommand(b *testing.B) {
	storage := NewMemoryStorage()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fact := &Fact{
			ID:     "bench_fact",
			Type:   "BenchType",
			Fields: map[string]interface{}{"value": i},
		}
		cmd := NewAddFactCommand(storage, fact)
		cmd.Execute()
		// Cleanup for next iteration
		storage.RemoveFact(fact.GetInternalID())
	}
}

// BenchmarkRemoveFactCommand benchmarks RemoveFact command execution
func BenchmarkRemoveFactCommand(b *testing.B) {
	storage := NewMemoryStorage()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		fact := &Fact{
			ID:     "bench_fact",
			Type:   "BenchType",
			Fields: map[string]interface{}{"value": i},
		}
		storage.AddFact(fact)
		b.StartTimer()
		cmd := NewRemoveFactCommand(storage, fact.GetInternalID())
		cmd.Execute()
	}
}
