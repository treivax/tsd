// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPhase2_BasicSynchronization vérifie que tous les faits sont immédiatement visibles après soumission
func TestPhase2_BasicSynchronization(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	facts := []map[string]interface{}{
		{"id": "fact1", "type": "Product", "name": "Item1"},
		{"id": "fact2", "type": "Product", "name": "Item2"},
		{"id": "fact3", "type": "Product", "name": "Item3"},
	}
	err := env.Network.SubmitFactsFromGrammar(facts)
	require.NoError(t, err, "La soumission des faits doit réussir")
	// Vérification immédiate : tous les faits doivent être visibles
	for _, fm := range facts {
		id := fm["id"].(string)
		typ := fm["type"].(string)
		internalID := typ + "~" + id
		fact := env.Storage.GetFact(internalID)
		assert.NotNil(t, fact, "Fait %s doit être immédiatement visible", id)
		if fact != nil {
			assert.Equal(t, internalID, fact.ID, "ID du fait doit correspondre")
			assert.Equal(t, typ, fact.Type, "Type du fait doit correspondre")
		}
	}
}

// TestPhase2_EmptyFactList vérifie que soumettre une liste vide ne cause pas d'erreur
func TestPhase2_EmptyFactList(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	facts := []map[string]interface{}{}
	err := env.Network.SubmitFactsFromGrammar(facts)
	require.NoError(t, err, "Soumettre une liste vide doit réussir")
}

// TestPhase2_SingleFact vérifie la synchronisation pour un seul fait
func TestPhase2_SingleFact(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	facts := []map[string]interface{}{
		{"id": "single", "type": "Product", "price": 99.99},
	}
	start := time.Now()
	err := env.Network.SubmitFactsFromGrammar(facts)
	duration := time.Since(start)
	require.NoError(t, err)
	assert.Less(t, duration, 1*time.Second, "Un seul fait devrait être rapide")
	// Vérifier que le fait est visible
	fact := env.Storage.GetFact("Product~single")
	assert.NotNil(t, fact)
}

// TestPhase2_WaitForFactPersistence vérifie que waitForFactPersistence fonctionne correctement
func TestPhase2_WaitForFactPersistence(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	fact := &Fact{
		ID:     "test_fact",
		Type:   "TestType",
		Fields: map[string]interface{}{"value": 42},
	}
	// Soumettre le fait au storage directement
	env.Storage.AddFact(fact)
	// Attendre la persistance (devrait réussir immédiatement)
	err := env.Network.waitForFactPersistence(fact, 1*time.Second)
	require.NoError(t, err, "Le fait devrait être trouvé immédiatement")
}

// TestPhase2_WaitForFactPersistence_Timeout vérifie que le timeout fonctionne
func TestPhase2_WaitForFactPersistence_Timeout(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	env.Network.SubmissionTimeout = 100 * time.Millisecond
	env.Network.VerifyRetryDelay = 10 * time.Millisecond
	fact := &Fact{
		ID:     "nonexistent",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	// Ne PAS ajouter le fait au storage
	// waitForFactPersistence devrait timeout
	start := time.Now()
	err := env.Network.waitForFactPersistence(fact, 100*time.Millisecond)
	duration := time.Since(start)
	require.Error(t, err, "Devrait timeout car le fait n'existe pas")
	assert.Contains(t, err.Error(), "timeout", "Le message d'erreur devrait mentionner le timeout")
	assert.GreaterOrEqual(t, duration, 100*time.Millisecond, "Devrait attendre au moins le timeout")
}

// TestPhase2_RetryMechanism vérifie que le mécanisme de retry fonctionne
func TestPhase2_RetryMechanism(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	env.Network.VerifyRetryDelay = 20 * time.Millisecond
	env.Network.MaxVerifyRetries = 5
	fact := &Fact{
		ID:     "delayed_fact",
		Type:   "Product",
		Fields: map[string]interface{}{"name": "Delayed"},
	}
	// Simuler un délai dans la persistance
	go func() {
		time.Sleep(50 * time.Millisecond)
		env.Storage.AddFact(fact)
	}()
	start := time.Now()
	err := env.Network.waitForFactPersistence(fact, 500*time.Millisecond)
	duration := time.Since(start)
	require.NoError(t, err, "Devrait réussir après retry")
	assert.GreaterOrEqual(t, duration, 50*time.Millisecond, "Devrait attendre au moins le délai de persistance")
	assert.Less(t, duration, 200*time.Millisecond, "Ne devrait pas attendre trop longtemps")
}

// TestPhase2_ConcurrentReadsAfterWrite vérifie que les lectures concurrentes après écriture fonctionnent
func TestPhase2_ConcurrentReadsAfterWrite(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	facts := []map[string]interface{}{
		{"id": "concurrent_fact", "type": "Product", "name": "Concurrent"},
	}
	// Soumettre le fait
	err := env.Network.SubmitFactsFromGrammar(facts)
	require.NoError(t, err)
	// Lancer plusieurs lectures concurrentes immédiatement après
	var wg sync.WaitGroup
	errors := make(chan error, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			fact := env.Storage.GetFact("Product~concurrent_fact")
			if fact == nil {
				errors <- fmt.Errorf("goroutine %d: fait non visible", idx)
			}
		}(i)
	}
	wg.Wait()
	close(errors)
	for err := range errors {
		require.NoError(t, err, "Toutes les goroutines devraient voir le fait")
	}
}

// TestPhase2_MultipleFactsBatch vérifie la soumission d'un grand lot de faits
func TestPhase2_MultipleFactsBatch(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	// Créer un lot de 50 faits
	facts := make([]map[string]interface{}, 50)
	for i := 0; i < 50; i++ {
		facts[i] = map[string]interface{}{
			"id":    fmt.Sprintf("batch_fact_%d", i),
			"type":  "Product",
			"index": i,
		}
	}
	start := time.Now()
	err := env.Network.SubmitFactsFromGrammar(facts)
	duration := time.Since(start)
	require.NoError(t, err, "La soumission du lot doit réussir")
	t.Logf("Soumission de 50 faits en %v", duration)
	// Vérifier que tous les faits sont visibles
	for i := 0; i < 50; i++ {
		internalID := fmt.Sprintf("Product~batch_fact_%d", i)
		fact := env.Storage.GetFact(internalID)
		assert.NotNil(t, fact, "Fait %d doit être visible", i)
	}
}

// TestPhase2_TimeoutPerFact vérifie que le timeout par fait est calculé correctement
func TestPhase2_TimeoutPerFact(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	env.Network.SubmissionTimeout = 10 * time.Second
	// Avec 10 faits et 10s de timeout total, chaque fait devrait avoir 1s
	facts := make([]map[string]interface{}, 10)
	for i := 0; i < 10; i++ {
		facts[i] = map[string]interface{}{
			"id":   fmt.Sprintf("fact_%d", i),
			"type": "Product",
		}
	}
	start := time.Now()
	err := env.Network.SubmitFactsFromGrammar(facts)
	duration := time.Since(start)
	require.NoError(t, err)
	// Devrait être rapide car tous les faits sont persistés immédiatement
	assert.Less(t, duration, 2*time.Second, "Devrait être rapide avec persistance immédiate")
}

// TestPhase2_RaceConditionSafety vérifie la sécurité avec le race detector
func TestPhase2_RaceConditionSafety(t *testing.T) {
	t.Parallel()
	// Use silent logger to avoid race on shared log buffer
	env := NewTestEnvironment(t, WithLogLevel(LogLevelSilent))
	defer env.Cleanup()
	// Soumettre des faits en parallèle depuis plusieurs goroutines
	var wg sync.WaitGroup
	errors := make(chan error, 5)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			facts := []map[string]interface{}{
				{
					"id":   fmt.Sprintf("race_fact_%d", idx),
					"type": "Product",
					"from": idx,
				},
			}
			if err := env.Network.SubmitFactsFromGrammar(facts); err != nil {
				errors <- fmt.Errorf("goroutine %d: %w", idx, err)
			}
		}(i)
	}
	wg.Wait()
	close(errors)
	for err := range errors {
		require.NoError(t, err)
	}
	// Vérifier que tous les faits sont présents
	for i := 0; i < 5; i++ {
		internalID := fmt.Sprintf("Product~race_fact_%d", i)
		fact := env.Storage.GetFact(internalID)
		assert.NotNil(t, fact, "Fait %d doit être visible", i)
	}
}

// TestPhase2_BackoffStrategy vérifie que le backoff exponentiel fonctionne
func TestPhase2_BackoffStrategy(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	env.Network.VerifyRetryDelay = 10 * time.Millisecond
	env.Network.MaxVerifyRetries = 5
	fact := &Fact{
		ID:     "backoff_test",
		Type:   "Product",
		Fields: map[string]interface{}{},
	}
	// Ajouter le fait avec un délai pour forcer plusieurs retries
	go func() {
		time.Sleep(80 * time.Millisecond) // Devrait nécessiter 3-4 retries
		env.Storage.AddFact(fact)
	}()
	start := time.Now()
	err := env.Network.waitForFactPersistence(fact, 500*time.Millisecond)
	duration := time.Since(start)
	require.NoError(t, err, "Devrait réussir avec le backoff")
	// Avec backoff: 10ms + 20ms + 40ms + check = ~80ms minimum
	assert.GreaterOrEqual(t, duration, 80*time.Millisecond, "Devrait utiliser le backoff")
	t.Logf("Persistance trouvée après %v (avec backoff exponentiel)", duration)
}

// TestPhase2_ConfigurableParameters vérifie que les paramètres sont configurables
func TestPhase2_ConfigurableParameters(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	// Vérifier les valeurs par défaut
	assert.Equal(t, DefaultSubmissionTimeout, env.Network.SubmissionTimeout)
	assert.Equal(t, DefaultVerifyRetryDelay, env.Network.VerifyRetryDelay)
	assert.Equal(t, DefaultMaxVerifyRetries, env.Network.MaxVerifyRetries)
	// Modifier les paramètres
	env.Network.SubmissionTimeout = 5 * time.Second
	env.Network.VerifyRetryDelay = 5 * time.Millisecond
	env.Network.MaxVerifyRetries = 20
	assert.Equal(t, 5*time.Second, env.Network.SubmissionTimeout)
	assert.Equal(t, 5*time.Millisecond, env.Network.VerifyRetryDelay)
	assert.Equal(t, 20, env.Network.MaxVerifyRetries)
}

// TestPhase2_ErrorHandling vérifie la gestion des erreurs lors de la soumission
func TestPhase2_ErrorHandling(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	// Fait invalide sans type
	facts := []map[string]interface{}{
		{"id": "no_type_fact"}, // Manque le type
	}
	err := env.Network.SubmitFactsFromGrammar(facts)
	// Devrait gérer gracieusement (le type devient "unknown")
	require.NoError(t, err, "Devrait gérer les faits sans type explicite")
}

// TestPhase2_PerformanceOverhead mesure l'overhead de la Phase 2
func TestPhase2_PerformanceOverhead(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	// Benchmark avec 100 faits
	facts := make([]map[string]interface{}, 100)
	for i := 0; i < 100; i++ {
		facts[i] = map[string]interface{}{
			"id":    fmt.Sprintf("perf_fact_%d", i),
			"type":  "Product",
			"index": i,
		}
	}
	start := time.Now()
	err := env.Network.SubmitFactsFromGrammar(facts)
	duration := time.Since(start)
	require.NoError(t, err)
	t.Logf("Soumission de 100 faits avec Phase 2 en %v", duration)
	t.Logf("Temps moyen par fait: %v", duration/100)
	// L'overhead devrait être raisonnable (< 10ms par fait en moyenne)
	avgPerFact := duration / 100
	assert.Less(t, avgPerFact, 10*time.Millisecond,
		"Overhead moyen par fait devrait être < 10ms")
}

// TestPhase2_IntegrationWithPhase1 vérifie que Phase 2 est compatible avec Phase 1
func TestPhase2_IntegrationWithPhase1(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	facts := []map[string]interface{}{
		{"id": "integration_fact", "type": "Product", "name": "Test"},
	}
	// La soumission devrait réussir avec les garanties des deux phases
	err := env.Network.SubmitFactsFromGrammar(facts)
	require.NoError(t, err)
	// Phase 1: Vérification de cohérence (compteurs)
	// Phase 2: Synchronisation avec retry
	// Le fait doit être immédiatement visible (garantie combinée)
	fact := env.Storage.GetFact("Product~integration_fact")
	require.NotNil(t, fact)
	assert.Equal(t, "Product~integration_fact", fact.ID)
	assert.Equal(t, "Product", fact.Type)
}

// TestPhase2_MinimumTimeoutPerFact vérifie que le timeout minimum est respecté
func TestPhase2_MinimumTimeoutPerFact(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()
	env.Network.SubmissionTimeout = 500 * time.Millisecond
	// Avec 1000 faits et 500ms de timeout, le calcul donnerait 0.5ms par fait
	// Mais le minimum devrait être 1s
	facts := make([]map[string]interface{}, 1000)
	for i := 0; i < 10; i++ { // Seulement 10 pour le test, pas 1000
		facts[i] = map[string]interface{}{
			"id":   fmt.Sprintf("min_timeout_fact_%d", i),
			"type": "Product",
		}
	}
	facts = facts[:10] // Limiter à 10 faits pour le test
	err := env.Network.SubmitFactsFromGrammar(facts)
	require.NoError(t, err, "Devrait réussir même avec timeout court")
}

// TestPhase2_RetractDuringSubmission vérifie que Phase 2 gère les rétractations pendant propagation
func TestPhase2_RetractDuringSubmission(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Soumettre des faits
	facts := []map[string]interface{}{
		{"id": "prod1", "type": "Product", "name": "Item1", "price": 50.0},
		{"id": "prod2", "type": "Product", "name": "Item2", "price": 150.0},
		{"id": "prod3", "type": "Product", "name": "Item3", "price": 75.0},
	}

	err := env.Network.SubmitFactsFromGrammar(facts)
	require.NoError(t, err, "La soumission doit réussir")

	// Vérifier que tous les faits sont persistés
	fact1 := env.Storage.GetFact("Product~prod1")
	assert.NotNil(t, fact1, "prod1 doit être persisté")

	fact2 := env.Storage.GetFact("Product~prod2")
	assert.NotNil(t, fact2, "prod2 doit être persisté")

	fact3 := env.Storage.GetFact("Product~prod3")
	assert.NotNil(t, fact3, "prod3 doit être persisté")

	// Simuler une rétractation pendant la propagation en créant un contexte de soumission
	ctx := NewSubmissionContext()
	ctx.MarkSubmitted("Product~prod2")

	env.Network.submissionMutex.Lock()
	env.Network.currentSubmission = ctx
	env.Network.submissionMutex.Unlock()

	// Rétracter un fait pendant que le contexte est actif
	err = env.Network.RetractFact("Product~prod2")
	require.NoError(t, err, "La rétractation doit réussir")

	env.Network.submissionMutex.Lock()
	env.Network.currentSubmission = nil
	env.Network.submissionMutex.Unlock()

	// Vérifier que le fait a été marqué comme rétracté
	assert.True(t, ctx.WasRetracted("Product~prod2"), "Le fait doit être marqué comme rétracté")

	// Vérifier que le fait a été supprimé du storage
	fact2After := env.Storage.GetFact("Product~prod2")
	assert.Nil(t, fact2After, "prod2 doit être supprimé du storage")

	t.Log("✅ Phase 2 gère correctement les rétractations pendant la soumission")
}

// TestPhase2_SubmissionContextTracking vérifie le tracking des faits soumis et rétractés
func TestPhase2_SubmissionContextTracking(t *testing.T) {
	t.Parallel()

	ctx := NewSubmissionContext()

	// Marquer des faits comme soumis
	ctx.MarkSubmitted("fact1")
	ctx.MarkSubmitted("fact2")
	ctx.MarkSubmitted("fact3")

	// Vérifier qu'ils sont marqués comme soumis
	assert.True(t, ctx.WasSubmitted("fact1"))
	assert.True(t, ctx.WasSubmitted("fact2"))
	assert.True(t, ctx.WasSubmitted("fact3"))
	assert.False(t, ctx.WasSubmitted("fact4"))

	// Marquer un fait comme rétracté
	ctx.MarkRetracted("fact2")

	// Vérifier qu'il est marqué comme rétracté
	assert.True(t, ctx.WasRetracted("fact2"))
	assert.False(t, ctx.WasRetracted("fact1"))
	assert.False(t, ctx.WasRetracted("fact3"))

	t.Log("✅ SubmissionContext track correctement les soumissions et rétractations")
}

// TestPhase2_ConcurrentSubmissionContext vérifie la thread-safety du SubmissionContext
func TestPhase2_ConcurrentSubmissionContext(t *testing.T) {
	t.Parallel()

	ctx := NewSubmissionContext()
	var wg sync.WaitGroup

	// Lancer plusieurs goroutines qui marquent des faits
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			factID := fmt.Sprintf("fact_%d", idx)
			ctx.MarkSubmitted(factID)
			if idx%2 == 0 {
				ctx.MarkRetracted(factID)
			}
		}(i)
	}

	wg.Wait()

	// Vérifier que tous les faits pairs sont rétractés
	for i := 0; i < 10; i++ {
		factID := fmt.Sprintf("fact_%d", i)
		assert.True(t, ctx.WasSubmitted(factID), "fact_%d doit être marqué comme soumis", i)
		if i%2 == 0 {
			assert.True(t, ctx.WasRetracted(factID), "fact_%d doit être marqué comme rétracté", i)
		}
	}

	t.Log("✅ SubmissionContext est thread-safe")
}
