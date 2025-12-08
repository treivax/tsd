// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTestEnvironment_Basic(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	assert.NotNil(t, env.Network)
	assert.NotNil(t, env.Storage)
	assert.NotNil(t, env.Pipeline)
	assert.NotNil(t, env.Logger)
	assert.NotNil(t, env.LogBuffer)
	assert.NotEmpty(t, env.TempDir)
	assert.Equal(t, LogLevelInfo, env.Logger.GetLevel())
}

func TestNewTestEnvironment_WithOptions(t *testing.T) {
	env := NewTestEnvironment(t,
		WithLogLevel(LogLevelDebug),
		WithTimestamps(true),
		WithLogPrefix("[TEST]"),
	)
	defer env.Cleanup()

	assert.Equal(t, LogLevelDebug, env.Logger.GetLevel())
	env.Logger.Info("test message")
	logs := env.GetLogs()
	assert.Contains(t, logs, "[TEST]")
}

func TestTestEnvironment_GetLogs(t *testing.T) {
	env := NewTestEnvironment(t, WithTimestamps(false))
	defer env.Cleanup()

	assert.Empty(t, env.GetLogs())
	env.Logger.Info("test message")
	logs := env.GetLogs()
	assert.Contains(t, logs, "test message")
	assert.Contains(t, logs, "[INFO]")
}

func TestTestEnvironment_ClearLogs(t *testing.T) {
	env := NewTestEnvironment(t, WithTimestamps(false))
	defer env.Cleanup()

	env.Logger.Info("first message")
	assert.NotEmpty(t, env.GetLogs())
	env.ClearLogs()
	assert.Empty(t, env.GetLogs())
	env.Logger.Info("second message")
	logs := env.GetLogs()
	assert.Contains(t, logs, "second message")
	assert.NotContains(t, logs, "first message")
}

func TestTestEnvironment_CreateConstraintFile(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	content := `type Person(name: string, age: number)
action print(message: string)
rule Adults : {p: Person} / p.age >= 18 ==> print(p.name)`

	filename := env.CreateConstraintFile("test.constraint", content)
	assert.FileExists(t, filename)
	assert.Contains(t, filename, env.TempDir)
}

func TestTestEnvironment_IngestFileContent(t *testing.T) {
	env := NewTestEnvironment(t, WithTimestamps(false))
	defer env.Cleanup()

	content := `type Person(name: string, age: number)
action print(message: string)
rule Adults : {p: Person} / p.age >= 18 ==> print(p.name)`

	network, _, err := env.IngestFileContent(content)
	require.NoError(t, err)
	assert.NotNil(t, network)
	logs := env.GetLogs()
	assert.Contains(t, logs, "INGESTION INCRÉMENTALE TERMINÉE")
}

func TestTestEnvironment_RequireIngestFileContent(t *testing.T) {
	env := NewTestEnvironment(t, WithTimestamps(false))
	defer env.Cleanup()

	content := `type Employee(id: number, name: string, salary: number)
action print(message: string)
rule HighEarners : {e: Employee} / e.salary > 100000 ==> print(e.name)`

	network := env.RequireIngestFileContent(content)
	assert.NotNil(t, network)
}

func TestTestEnvironment_SubmitFact(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	content := `type Person(name: string, age: number)
action print(message: string)
rule Adults : {p: Person} / p.age >= 18 ==> print(p.name)`
	env.RequireIngestFileContent(content)

	fact := Fact{
		ID:   "alice1",
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Alice",
			"age":  float64(25),
		},
	}

	err := env.SubmitFact(fact)
	require.NoError(t, err)
	assert.Equal(t, 1, env.GetFactCount())
}

func TestTestEnvironment_GetFactCount(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	assert.Equal(t, 0, env.GetFactCount())
	content := `type Item(id: number)`
	env.RequireIngestFileContent(content)

	for i := 1; i <= 3; i++ {
		fact := Fact{
			ID:   fmt.Sprintf("item_%d", i),
			Type: "Item",
			Fields: map[string]interface{}{
				"id": float64(i),
			},
		}
		env.RequireSubmitFact(fact)
	}

	assert.Equal(t, 3, env.GetFactCount())
}

func TestTestEnvironment_GetFactsByType(t *testing.T) {
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	content := `type Person(name: string)
type Company(name: string)`
	env.RequireIngestFileContent(content)

	person1 := Fact{ID: "person1", Type: "Person", Fields: map[string]interface{}{"name": "Alice"}}
	person2 := Fact{ID: "person2", Type: "Person", Fields: map[string]interface{}{"name": "Bob"}}
	company := Fact{ID: "company1", Type: "Company", Fields: map[string]interface{}{"name": "Acme"}}

	env.RequireSubmitFact(person1)
	env.RequireSubmitFact(person2)
	env.RequireSubmitFact(company)

	persons := env.GetFactsByType("Person")
	assert.Len(t, persons, 2)

	companies := env.GetFactsByType("Company")
	assert.Len(t, companies, 1)

	nonExistent := env.GetFactsByType("Product")
	assert.Len(t, nonExistent, 0)
}

func TestTestEnvironment_SetLogLevel(t *testing.T) {
	env := NewTestEnvironment(t, WithLogLevel(LogLevelInfo), WithTimestamps(false))
	defer env.Cleanup()

	env.Logger.Debug("debug message")
	assert.Empty(t, env.GetLogs())

	env.SetLogLevel(LogLevelDebug)
	env.Logger.Debug("debug message 2")

	logs := env.GetLogs()
	assert.Contains(t, logs, "debug message 2")
}

func TestTestEnvironment_NewSubEnvironment(t *testing.T) {
	env := NewTestEnvironment(t, WithTimestamps(false))
	defer env.Cleanup()

	content := `type Person(name: string)`
	env.RequireIngestFileContent(content)

	fact := Fact{ID: "alice2", Type: "Person", Fields: map[string]interface{}{"name": "Alice"}}
	env.RequireSubmitFact(fact)

	subEnv := env.NewSubEnvironment(WithLogLevel(LogLevelDebug))

	assert.NotEqual(t, env.Network, subEnv.Network)
	assert.NotEqual(t, env.Logger, subEnv.Logger)
	assert.Equal(t, env.Storage, subEnv.Storage)
	assert.Equal(t, 1, subEnv.GetFactCount())

	env.Logger.Info("main log")
	subEnv.Logger.Info("sub log")

	mainLogs := env.GetLogs()
	subLogs := subEnv.GetLogs()

	assert.Contains(t, mainLogs, "main log")
	assert.NotContains(t, mainLogs, "sub log")
	assert.Contains(t, subLogs, "sub log")
	assert.NotContains(t, subLogs, "main log")
}

func TestTestEnvironment_Cleanup(t *testing.T) {
	env := NewTestEnvironment(t)

	cleanupCalled := false
	env.AddCleanup(func() {
		cleanupCalled = true
	})

	env.Cleanup()
	assert.True(t, cleanupCalled)
	env.Cleanup() // Should not panic
}

func TestTestEnvironment_MultipleCleanups(t *testing.T) {
	env := NewTestEnvironment(t)

	order := []int{}
	env.AddCleanup(func() { order = append(order, 1) })
	env.AddCleanup(func() { order = append(order, 2) })
	env.AddCleanup(func() { order = append(order, 3) })

	env.Cleanup()
	assert.Equal(t, []int{3, 2, 1}, order)
}

func TestTestEnvironment_WithCustomStorage(t *testing.T) {
	customStorage := NewMemoryStorage()
	env := NewTestEnvironment(t, WithCustomStorage(customStorage))
	defer env.Cleanup()

	assert.Equal(t, customStorage, env.Storage)
}

func TestTestEnvironment_ParallelSafety(t *testing.T) {
	t.Run("env1", func(t *testing.T) {
		t.Parallel()
		env := NewTestEnvironment(t, WithLogPrefix("[ENV1]"), WithTimestamps(false))
		defer env.Cleanup()

		env.Logger.Info("message from env1")
		logs := env.GetLogs()
		assert.Contains(t, logs, "[ENV1]")
		assert.NotContains(t, logs, "[ENV2]")
	})

	t.Run("env2", func(t *testing.T) {
		t.Parallel()
		env := NewTestEnvironment(t, WithLogPrefix("[ENV2]"), WithTimestamps(false))
		defer env.Cleanup()

		env.Logger.Info("message from env2")
		logs := env.GetLogs()
		assert.Contains(t, logs, "[ENV2]")
		assert.NotContains(t, logs, "[ENV1]")
	})
}
