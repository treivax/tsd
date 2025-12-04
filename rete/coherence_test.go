// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCoherence_TransactionRollback teste que le rollback fonctionne correctement en cas d'incohérence
func TestCoherence_TransactionRollback(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Créer une transaction
	tx := env.Network.BeginTransaction()
	env.Network.SetTransaction(tx)

	// Ajouter quelques faits
	fact1 := &Fact{
		ID:     "F1",
		Type:   "TestType",
		Fields: map[string]interface{}{"value": 42},
	}
	fact2 := &Fact{
		ID:     "F2",
		Type:   "TestType",
		Fields: map[string]interface{}{"value": 100},
	}

	err := env.Network.SubmitFact(fact1)
	require.NoError(t, err)

	err = env.Network.SubmitFact(fact2)
	require.NoError(t, err)

	// Vérifier que les faits sont présents (avec ID interne)
	assert.NotNil(t, env.Storage.GetFact("TestType_F1"))
	assert.NotNil(t, env.Storage.GetFact("TestType_F2"))

	// Rollback
	err = tx.Rollback()
	require.NoError(t, err)

	// Vérifier que les faits ont été supprimés après rollback
	assert.Nil(t, env.Storage.GetFact("TestType_F1"), "Le fait F1 doit être supprimé après rollback")
	assert.Nil(t, env.Storage.GetFact("TestType_F2"), "Le fait F2 doit être supprimé après rollback")
}

// TestCoherence_StorageSync teste que Storage.Sync() fonctionne sans erreur
func TestCoherence_StorageSync(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Ajouter quelques faits
	fact := &Fact{
		ID:     "F1",
		Type:   "TestType",
		Fields: map[string]interface{}{"value": 42},
	}

	err := env.Storage.AddFact(fact)
	require.NoError(t, err)

	// Appeler Sync
	err = env.Storage.Sync()
	assert.NoError(t, err, "Storage.Sync() ne doit pas échouer")

	// Vérifier que le fait est toujours là (avec ID interne)
	retrievedFact := env.Storage.GetFact("TestType_F1")
	assert.NotNil(t, retrievedFact, "Le fait doit toujours être présent après Sync()")
	assert.Equal(t, "F1", retrievedFact.ID)
	assert.Equal(t, "TestType", retrievedFact.Type)
}

// TestCoherence_InternalIDCorrectness vérifie que les IDs internes sont correctement utilisés
func TestCoherence_InternalIDCorrectness(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Créer un fait
	fact := &Fact{
		ID:     "TEST123",
		Type:   "MyType",
		Fields: map[string]interface{}{"value": 42},
	}

	// Vérifier l'ID interne
	internalID := fact.GetInternalID()
	assert.Equal(t, "MyType_TEST123", internalID, "L'ID interne doit être Type_ID")

	// Vérifier que le storage utilise bien l'ID interne
	err := env.Storage.AddFact(fact)
	require.NoError(t, err)

	// Chercher avec l'ID interne (doit réussir)
	retrievedFact := env.Storage.GetFact(internalID)
	assert.NotNil(t, retrievedFact, "Le fait doit être trouvable avec l'ID interne")

	// Chercher avec l'ID simple (doit échouer)
	notFound := env.Storage.GetFact("TEST123")
	assert.Nil(t, notFound, "Le fait ne doit PAS être trouvable avec l'ID simple")
}

// TestCoherence_FactSubmissionConsistency vérifie que SubmitFactsFromGrammar valide la cohérence
func TestCoherence_FactSubmissionConsistency(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Démarrer une transaction
	tx := env.Network.BeginTransaction()
	env.Network.SetTransaction(tx)

	// Créer des faits à soumettre
	factsToSubmit := []map[string]interface{}{
		{
			"id":   "F1",
			"type": "TestType",
			"val":  100,
		},
		{
			"id":   "F2",
			"type": "TestType",
			"val":  200,
		},
		{
			"id":   "F3",
			"type": "TestType",
			"val":  300,
		},
	}

	// Soumettre les faits
	err := env.Network.SubmitFactsFromGrammar(factsToSubmit)
	require.NoError(t, err, "La soumission ne doit pas échouer si tous les faits sont persistés")

	// Vérifier que tous les faits sont présents avec leur ID interne
	assert.NotNil(t, env.Storage.GetFact("TestType_F1"))
	assert.NotNil(t, env.Storage.GetFact("TestType_F2"))
	assert.NotNil(t, env.Storage.GetFact("TestType_F3"))

	// Commit
	err = tx.Commit()
	require.NoError(t, err)

	// Vérifier que les faits sont toujours là après commit
	assert.NotNil(t, env.Storage.GetFact("TestType_F1"))
	assert.NotNil(t, env.Storage.GetFact("TestType_F2"))
	assert.NotNil(t, env.Storage.GetFact("TestType_F3"))
}

// TestCoherence_ConcurrentFactAddition teste l'ajout concurrent de faits
// Each goroutine uses its own isolated environment for thread safety
func TestCoherence_ConcurrentFactAddition(t *testing.T) {
	t.Parallel()

	numGoroutines := 10
	factsPerGoroutine := 5

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)
	factCounts := make(chan int, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Each goroutine gets its own isolated environment (network + storage)
			env := NewTestEnvironment(t, WithLogLevel(LogLevelSilent))
			defer env.Cleanup()

			// Create transaction for this goroutine's environment
			tx := env.Network.BeginTransaction()
			env.Network.SetTransaction(tx)

			for j := 0; j < factsPerGoroutine; j++ {
				fact := &Fact{
					ID:     fmt.Sprintf("G%d_F%d", id, j),
					Type:   "ConcurrentTest",
					Fields: map[string]interface{}{"goroutine": id, "index": j},
				}

				if err := env.Network.SubmitFact(fact); err != nil {
					errors <- fmt.Errorf("goroutine %d, fact %d: %w", id, j, err)
					return
				}
			}

			// Commit
			if err := tx.Commit(); err != nil {
				errors <- fmt.Errorf("goroutine %d commit: %w", id, err)
				return
			}

			// Report fact count from this environment
			factCounts <- len(env.Storage.GetAllFacts())
		}(i)
	}

	wg.Wait()
	close(errors)
	close(factCounts)

	// Vérifier qu'il n'y a pas eu d'erreurs
	for err := range errors {
		require.NoError(t, err)
	}

	// Vérifier que chaque goroutine a bien persisté ses faits
	totalFacts := 0
	for count := range factCounts {
		assert.Equal(t, factsPerGoroutine, count, "Chaque environment doit avoir exactement %d faits", factsPerGoroutine)
		totalFacts += count
	}

	expectedTotal := numGoroutines * factsPerGoroutine
	assert.Equal(t, expectedTotal, totalFacts, "Le nombre total de faits doit être %d", expectedTotal)
}

// TestCoherence_SyncAfterMultipleAdditions teste la synchronisation après plusieurs ajouts
func TestCoherence_SyncAfterMultipleAdditions(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Ajouter plusieurs faits
	for i := 0; i < 100; i++ {
		fact := &Fact{
			ID:     fmt.Sprintf("FACT_%d", i),
			Type:   "BulkTest",
			Fields: map[string]interface{}{"index": i},
		}
		err := env.Storage.AddFact(fact)
		require.NoError(t, err)
	}

	// Synchroniser
	err := env.Storage.Sync()
	assert.NoError(t, err, "Sync() doit réussir après ajouts multiples")

	// Vérifier que tous les faits sont toujours là
	allFacts := env.Storage.GetAllFacts()
	assert.Equal(t, 100, len(allFacts), "Tous les 100 faits doivent être présents après Sync()")

	// Vérifier quelques faits spécifiques
	fact0 := env.Storage.GetFact("BulkTest_FACT_0")
	fact50 := env.Storage.GetFact("BulkTest_FACT_50")
	fact99 := env.Storage.GetFact("BulkTest_FACT_99")

	assert.NotNil(t, fact0)
	assert.NotNil(t, fact50)
	assert.NotNil(t, fact99)
}

// TestCoherence_ReadAfterWriteGuarantee teste la garantie read-after-write
func TestCoherence_ReadAfterWriteGuarantee(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	tx := env.Network.BeginTransaction()
	env.Network.SetTransaction(tx)

	// Écrire un fait
	fact := &Fact{
		ID:     "RW_TEST",
		Type:   "ReadWriteTest",
		Fields: map[string]interface{}{"data": "test_value"},
	}

	err := env.Network.SubmitFact(fact)
	require.NoError(t, err)

	// Lire immédiatement (read-after-write)
	internalID := "ReadWriteTest_RW_TEST"
	retrievedFact := env.Storage.GetFact(internalID)

	// Le fait DOIT être visible immédiatement
	assert.NotNil(t, retrievedFact, "Le fait doit être visible immédiatement après écriture (read-after-write)")
	assert.Equal(t, "RW_TEST", retrievedFact.ID)
	assert.Equal(t, "ReadWriteTest", retrievedFact.Type)
	assert.Equal(t, "test_value", retrievedFact.Fields["data"])

	// Commit
	err = tx.Commit()
	require.NoError(t, err)

	// Vérifier que le fait est toujours là après commit
	retrievedFactAfterCommit := env.Storage.GetFact(internalID)
	assert.NotNil(t, retrievedFactAfterCommit, "Le fait doit rester visible après commit")
}
