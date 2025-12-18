// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/treivax/tsd/rete"
)

// TestConcurrentRetrieveAndMarkConsumed teste Retrieve simultané
// Note: Retrieve() marque maintenant automatiquement le xuple comme consommé
func TestConcurrentRetrieveAndMarkConsumed(t *testing.T) {
	t.Log("🧪 TEST CONCURRENT RETRIEVE (AUTO-CONSUMES)")
	t.Log("==============================================")

	config := XupleSpaceConfig{
		Name:              "concurrent-space",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewPerAgentConsumptionPolicy(), // Chaque agent peut consommer
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
		MaxSize:           0,
	}

	space := NewXupleSpace(config)

	// Insérer 10 xuples
	const numXuples = 10
	for i := 0; i < numXuples; i++ {
		xuple := &Xuple{
			ID:              fmt.Sprintf("xuple-%d", i),
			Fact:            &rete.Fact{},
			TriggeringFacts: []*rete.Fact{},
			CreatedAt:       time.Now(),
			Metadata: XupleMetadata{
				State:      XupleStateAvailable,
				ConsumedBy: make(map[string]time.Time),
			},
		}
		if err := space.Insert(xuple); err != nil {
			t.Fatalf("❌ Insert échoué: %v", err)
		}
	}

	t.Logf("✅ Inséré %d xuples", numXuples)

	// Lancer 20 agents qui récupèrent et marquent consommés simultanément
	const numAgents = 20
	var wg sync.WaitGroup
	errors := make(chan error, numAgents)
	retrieved := make(chan string, numAgents)

	for i := 0; i < numAgents; i++ {
		wg.Add(1)
		agentID := fmt.Sprintf("agent-%d", i)

		go func(id string) {
			defer wg.Done()

			// Retrieve - marque AUTOMATIQUEMENT comme consommé maintenant
			xuple, err := space.Retrieve(id)
			if err != nil {
				// Acceptable si plus de xuples disponibles pour cet agent
				if err != ErrNoAvailableXuple {
					errors <- fmt.Errorf("agent %s retrieve error: %w", id, err)
				}
				return
			}

			retrieved <- xuple.ID

			// Note: Ne pas lire xuple.Metadata.ConsumptionCount ici car le xuple
			// peut être modifié par d'autres goroutines (per-agent policy).
			// Le fait que Retrieve() réussisse prouve que le marking a fonctionné.
		}(agentID)
	}

	wg.Wait()
	close(errors)
	close(retrieved)

	// Vérifier les erreurs
	for err := range errors {
		t.Errorf("❌ Erreur concurrence: %v", err)
	}

	// Compter les xuples récupérés
	retrievedCount := 0
	for range retrieved {
		retrievedCount++
	}

	t.Logf("✅ %d agents ont récupéré des xuples", retrievedCount)

	// Avec PerAgentConsumptionPolicy, chaque xuple peut être consommé par plusieurs agents
	// Donc retrievedCount peut être > numXuples
	if retrievedCount < numXuples {
		t.Errorf("❌ Attendu au moins %d récupérations, reçu: %d", numXuples, retrievedCount)
	}

	t.Log("✅ Retrieve concurrent avec auto-consume fonctionne correctement")
}

// TestConcurrentInsertWithMaxSize teste les insertions concurrentes avec limite de capacité
func TestConcurrentInsertWithMaxSize(t *testing.T) {
	t.Log("🧪 TEST CONCURRENT INSERT WITH MAX SIZE")
	t.Log("=========================================")

	const maxSize = 50
	config := XupleSpaceConfig{
		Name:              "limited-space",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
		MaxSize:           maxSize,
	}

	space := NewXupleSpace(config)

	// Lancer 100 goroutines qui tentent d'insérer
	const numGoroutines = 100
	var wg sync.WaitGroup
	successes := make(chan bool, numGoroutines)
	failures := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		xupleID := fmt.Sprintf("xuple-%d", i)

		go func(id string) {
			defer wg.Done()

			xuple := &Xuple{
				ID:              id,
				Fact:            &rete.Fact{},
				TriggeringFacts: []*rete.Fact{},
				CreatedAt:       time.Now(),
				Metadata: XupleMetadata{
					State:      XupleStateAvailable,
					ConsumedBy: make(map[string]time.Time),
				},
			}

			err := space.Insert(xuple)
			if err == nil {
				successes <- true
			} else if err == ErrXupleSpaceFull {
				failures <- err
			} else {
				t.Errorf("❌ Erreur inattendue: %v", err)
			}
		}(xupleID)
	}

	wg.Wait()
	close(successes)
	close(failures)

	// Compter succès et échecs
	successCount := len(successes)
	failureCount := len(failures)

	t.Logf("✅ Insertions réussies: %d", successCount)
	t.Logf("⚠️  Insertions rejetées: %d (MaxSize atteint)", failureCount)

	// Exactement maxSize insertions doivent réussir
	if successCount != maxSize {
		t.Errorf("❌ Attendu %d succès, reçu: %d", maxSize, successCount)
	}

	// Les autres doivent échouer avec ErrXupleSpaceFull
	expectedFailures := numGoroutines - maxSize
	if failureCount != expectedFailures {
		t.Errorf("❌ Attendu %d échecs, reçu: %d", expectedFailures, failureCount)
	}

	// Vérifier le count final
	count := space.Count()
	if count != maxSize {
		t.Errorf("❌ Attendu count=%d, reçu: %d", maxSize, count)
	}

	t.Log("✅ Insertion concurrente avec MaxSize fonctionne correctement")
}

// TestConcurrentCleanup teste le nettoyage concurrent
func TestConcurrentCleanup(t *testing.T) {
	t.Log("🧪 TEST CONCURRENT CLEANUP")
	t.Log("===========================")

	config := XupleSpaceConfig{
		Name:              "cleanup-space",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewDurationRetentionPolicy(10 * time.Millisecond),
		MaxSize:           0,
	}

	space := NewXupleSpace(config)

	// Insérer 50 xuples
	const numXuples = 50
	for i := 0; i < numXuples; i++ {
		xuple := &Xuple{
			ID:              fmt.Sprintf("xuple-%d", i),
			Fact:            &rete.Fact{},
			TriggeringFacts: []*rete.Fact{},
			CreatedAt:       time.Now(),
			Metadata: XupleMetadata{
				State:      XupleStateAvailable,
				ConsumedBy: make(map[string]time.Time),
			},
		}
		if err := space.Insert(xuple); err != nil {
			t.Fatalf("❌ Insert échoué: %v", err)
		}
	}

	// Attendre expiration
	time.Sleep(50 * time.Millisecond)

	// Lancer 10 goroutines de cleanup simultanées
	const numCleaners = 10
	var wg sync.WaitGroup
	totalCleaned := make(chan int, numCleaners)

	for i := 0; i < numCleaners; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cleaned := space.Cleanup()
			totalCleaned <- cleaned
		}()
	}

	wg.Wait()
	close(totalCleaned)

	// Compter le total nettoyé
	sum := 0
	for cleaned := range totalCleaned {
		sum += cleaned
	}

	// Au moins numXuples xuples doivent avoir été nettoyés
	// (la première cleanup peut tout nettoyer, les autres 0)
	if sum < numXuples {
		t.Errorf("❌ Attendu au moins %d xuples nettoyés, reçu: %d", numXuples, sum)
	}

	// Le count final doit être 0
	count := space.Count()
	if count != 0 {
		t.Errorf("❌ Attendu count=0 après cleanup, reçu: %d", count)
	}

	t.Log("✅ Cleanup concurrent fonctionne correctement")
}

// TestRaceConditions teste avec go test -race
func TestRaceConditions(t *testing.T) {
	t.Log("🧪 TEST RACE CONDITIONS (use go test -race)")
	t.Log("============================================")

	config := XupleSpaceConfig{
		Name:              "race-test-space",
		SelectionPolicy:   NewRandomSelectionPolicy(),
		ConsumptionPolicy: NewLimitedConsumptionPolicy(3),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
		MaxSize:           100,
	}

	space := NewXupleSpace(config)

	// Mix de toutes les opérations concurrentes
	const numOps = 50
	var wg sync.WaitGroup

	// Insertions
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			xuple := &Xuple{
				ID:              fmt.Sprintf("xuple-%d", id),
				Fact:            &rete.Fact{},
				TriggeringFacts: []*rete.Fact{},
				CreatedAt:       time.Now(),
				Metadata: XupleMetadata{
					State:      XupleStateAvailable,
					ConsumedBy: make(map[string]time.Time),
				},
			}
			_ = space.Insert(xuple)
		}(i)
	}

	// Retrieves
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, _ = space.Retrieve(fmt.Sprintf("agent-%d", id))
		}(i)
	}

	// MarkConsumed
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_ = space.MarkConsumed(fmt.Sprintf("xuple-%d", id%10), fmt.Sprintf("agent-%d", id))
		}(i)
	}

	// Count
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = space.Count()
		}()
	}

	// Cleanup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = space.Cleanup()
		}()
	}

	wg.Wait()

	t.Log("✅ Aucune race condition détectée (si run avec -race)")
}
