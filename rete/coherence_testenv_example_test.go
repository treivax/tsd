// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Example conversion of coherence tests to use TestEnvironment
// This demonstrates the pattern for converting existing tests
func TestCoherence_StorageSync_WithTestEnv(t *testing.T) {
	t.Parallel() // Now safe with isolated environment!
	env := NewTestEnvironment(t, WithLogLevel(LogLevelDebug))
	defer env.Cleanup()
	// Set up a simple type
	content := `type TestType(value: number)`
	env.RequireIngestFileContent(content)
	// Add a fact
	fact := Fact{
		ID:   "F1",
		Type: "TestType",
		Fields: map[string]interface{}{
			"value": float64(42),
		},
	}
	env.RequireSubmitFact(fact)
	// Call Sync
	err := env.Storage.Sync()
	assert.NoError(t, err, "Storage.Sync() should not fail")
	// Verify fact is still present
	retrievedFact := env.Storage.GetFact("TestType_F1")
	assert.NotNil(t, retrievedFact, "Fact should still be present after Sync()")
	assert.Equal(t, "F1", retrievedFact.ID)
	assert.Equal(t, "TestType", retrievedFact.Type)
	// Verify we can query by count
	assert.Equal(t, 1, env.GetFactCount(), "Should have exactly 1 fact")
}
func TestCoherence_InternalIDCorrectness_WithTestEnv(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	// Set up type
	content := `type MyType(value: number)`
	env.RequireIngestFileContent(content)
	// Create a fact
	fact := Fact{
		ID:   "TEST123",
		Type: "MyType",
		Fields: map[string]interface{}{
			"value": float64(42),
		},
	}
	// Verify internal ID
	internalID := fact.GetInternalID()
	assert.Equal(t, "MyType_TEST123", internalID, "Internal ID should be Type_ID")
	// Submit fact
	env.RequireSubmitFact(fact)
	// Retrieve with internal ID
	retrievedFact := env.Storage.GetFact(internalID)
	require.NotNil(t, retrievedFact, "Fact should be retrievable by internal ID")
	assert.Equal(t, "TEST123", retrievedFact.ID)
	assert.Equal(t, "MyType", retrievedFact.Type)
}
func TestCoherence_MultipleFactSubmission_WithTestEnv(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t, WithTimestamps(false))
	defer env.Cleanup()
	// Set up type with rule
	content := `type Product(id: number, name: string, price: number)
action print(message: string)
rule ExpensiveProducts : {p: Product} / p.price > 1000 ==> print(p.name)`
	env.RequireIngestFileContent(content)
	// Submit multiple facts
	products := []Fact{
		{ID: "P1", Type: "Product", Fields: map[string]interface{}{"id": float64(1), "name": "Laptop", "price": float64(1500)}},
		{ID: "P2", Type: "Product", Fields: map[string]interface{}{"id": float64(2), "name": "Mouse", "price": float64(50)}},
		{ID: "P3", Type: "Product", Fields: map[string]interface{}{"id": float64(3), "name": "Monitor", "price": float64(2000)}},
	}
	for _, p := range products {
		env.RequireSubmitFact(p)
	}
	// Verify all facts stored
	assert.Equal(t, 3, env.GetFactCount(), "Should have 3 facts")
	// Verify can filter by type
	productFacts := env.GetFactsByType("Product")
	assert.Len(t, productFacts, 3, "Should retrieve all Product facts")
	// Verify logs show operations
	logs := env.GetLogs()
	assert.Contains(t, logs, "INGESTION INCRÉMENTALE TERMINÉE", "Should log ingestion")
	env.AssertNoErrors(t)
}
func TestCoherence_TransactionPattern_WithTestEnv(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	// Set up type
	content := `type Order(id: number, total: number)`
	env.RequireIngestFileContent(content)
	// Create transaction explicitly
	tx := env.Network.BeginTransaction()
	env.Network.SetTransaction(tx)
	// Add facts via transaction
	fact1 := Fact{
		ID:   "O1",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":    float64(1),
			"total": float64(100),
		},
	}
	cmd1 := NewAddFactCommand(env.Storage, &fact1)
	err := tx.RecordAndExecute(cmd1)
	require.NoError(t, err, "First fact should be added successfully")
	// Verify fact is accessible before commit
	assert.NotNil(t, env.Storage.GetFact("Order_O1"), "Fact should be accessible before commit")
	// Commit transaction
	err = tx.Commit()
	require.NoError(t, err, "Transaction commit should succeed")
	// Verify fact persists after commit
	assert.NotNil(t, env.Storage.GetFact("Order_O1"), "Fact should persist after commit")
	assert.Equal(t, 1, env.GetFactCount())
}
func TestCoherence_SubEnvironmentSharing_WithTestEnv(t *testing.T) {
	t.Parallel()
	// Main environment
	mainEnv := NewTestEnvironment(t, WithLogPrefix("[MAIN]"), WithTimestamps(false))
	defer mainEnv.Cleanup()
	// Set up shared type and fact
	content := `type SharedData(id: number, value: string)`
	mainEnv.RequireIngestFileContent(content)
	sharedFact := Fact{
		ID:   "S1",
		Type: "SharedData",
		Fields: map[string]interface{}{
			"id":    float64(1),
			"value": "shared",
		},
	}
	mainEnv.RequireSubmitFact(sharedFact)
	// Create sub-environment with shared storage
	subEnv := mainEnv.NewSubEnvironment(WithLogPrefix("[SUB]"))
	// Sub-environment can see main environment's facts
	assert.Equal(t, 1, subEnv.GetFactCount(), "Sub-env should see main env facts")
	assert.NotNil(t, subEnv.Storage.GetFact("SharedData_S1"), "Sub-env can access shared fact")
	// Add fact in sub-environment
	subFact := Fact{
		ID:   "S2",
		Type: "SharedData",
		Fields: map[string]interface{}{
			"id":    float64(2),
			"value": "from_sub",
		},
	}
	subEnv.RequireSubmitFact(subFact)
	// Main environment can see sub-environment's facts (shared storage)
	assert.Equal(t, 2, mainEnv.GetFactCount(), "Main env should see sub env facts")
	// But logs are isolated
	mainEnv.Logger.Info("main message")
	subEnv.Logger.Info("sub message")
	mainLogs := mainEnv.GetLogs()
	subLogs := subEnv.GetLogs()
	assert.Contains(t, mainLogs, "[MAIN]", "Main logs have MAIN prefix")
	assert.Contains(t, mainLogs, "main message")
	assert.NotContains(t, mainLogs, "sub message", "Main logs don't include sub logs")
	assert.Contains(t, subLogs, "[SUB]", "Sub logs have SUB prefix")
	assert.Contains(t, subLogs, "sub message")
	assert.NotContains(t, subLogs, "main message", "Sub logs don't include main logs")
}
func TestCoherence_ConcurrentAccess_WithTestEnv(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t, WithLogLevel(LogLevelInfo), WithTimestamps(false))
	defer env.Cleanup()
	// Set up type
	content := `type Counter(id: number, value: number)`
	env.RequireIngestFileContent(content)
	// Submit fact via helper (which handles transactions)
	fact := Fact{
		ID:   "C1",
		Type: "Counter",
		Fields: map[string]interface{}{
			"id":    float64(1),
			"value": float64(0),
		},
	}
	err := env.SubmitFact(fact)
	require.NoError(t, err, "Fact submission should succeed")
	// Verify fact is stored
	assert.Equal(t, 1, env.GetFactCount(), "Should have 1 fact")
	// Verify we can retrieve it
	retrieved := env.Storage.GetFact("Counter_C1")
	require.NotNil(t, retrieved, "Fact should be retrievable")
	assert.Equal(t, "C1", retrieved.ID)
}
