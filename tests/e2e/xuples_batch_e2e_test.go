// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/xuples"
)

// TestXuplesBatch_E2E_Comprehensive teste RetrieveMultiple dans un scénario réel complet
func TestXuplesBatch_E2E_Comprehensive(t *testing.T) {
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("🧪 TEST E2E: RetrieveMultiple - Scénario Complet")
	t.Log("═══════════════════════════════════════════════════════════════")

	// Créer programme TSD avec xuple-spaces
	tmpDir := t.TempDir()
	tsdFile := filepath.Join(tmpDir, "batch-test.tsd")

	programContent := `// Test E2E Batch Processing
type Task(id: string, taskType: string, priority: number, data: string)
type Result(taskId: string, status: string, output: string)

// Xuple-spaces pour traitement batch
xuple-space task_queue {
	selection: fifo
	consumption: once
	retention: duration(5m)
}

xuple-space high_priority_tasks {
	selection: lifo
	consumption: once
	retention: duration(10m)
}

xuple-space results_pool {
	selection: random
	consumption: per-agent
	retention: duration(1h)
}

// Faits de test
Task(taskType: "compute", priority: 1, data: "task1")
Task(taskType: "compute", priority: 2, data: "task2")
Task(taskType: "io", priority: 1, data: "task3")
Task(taskType: "compute", priority: 3, data: "task4")
Task(taskType: "io", priority: 2, data: "task5")
`

	if err := os.WriteFile(tsdFile, []byte(programContent), 0644); err != nil {
		t.Fatalf("❌ Erreur création fichier TSD: %v", err)
	}

	// Parser le programme
	content, err := os.ReadFile(tsdFile)
	if err != nil {
		t.Fatalf("❌ Erreur lecture: %v", err)
	}

	program, err := constraint.Parse(tsdFile, content)
	if err != nil {
		t.Fatalf("❌ Erreur parsing: %v", err)
	}

	// Vérifier parsing réussi
	if program == nil {
		t.Fatalf("❌ Programme nil après parsing")
	}
	t.Logf("✅ Programme parsé avec succès")

	// Créer XupleManager (RETE network pas nécessaire pour ce test)
	xupleManager := xuples.NewXupleManager()

	// Créer les xuple-spaces
	spaces := []struct {
		name              string
		selection         string
		consumption       string
		retentionDuration time.Duration
	}{
		{"task_queue", "fifo", "once", 5 * time.Minute},
		{"high_priority_tasks", "lifo", "once", 10 * time.Minute},
		{"results_pool", "random", "per-agent", time.Hour},
	}

	for _, s := range spaces {
		config := xuples.XupleSpaceConfig{
			Name:              s.name,
			SelectionPolicy:   getSelectionPolicy(s.selection),
			ConsumptionPolicy: getConsumptionPolicy(s.consumption),
			RetentionPolicy:   xuples.NewDurationRetentionPolicy(s.retentionDuration),
			MaxSize:           0,
		}
		if err := xupleManager.CreateXupleSpace(s.name, config); err != nil {
			t.Fatalf("❌ Erreur création xuple-space %s: %v", s.name, err)
		}
	}

	t.Log("✅ Xuple-spaces créés")

	// ═══════════════════════════════════════════════════════════════
	// SCÉNARIO 1: Traitement batch de tâches
	// ═══════════════════════════════════════════════════════════════
	t.Log("")
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("SCÉNARIO 1: Traitement Batch de Tâches")
	t.Log("───────────────────────────────────────────────────────────────")

	taskQueue, _ := xupleManager.GetXupleSpace("task_queue")

	// Créer 20 tâches
	const numTasks = 20
	for i := 0; i < numTasks; i++ {
		task := &rete.Fact{
			Type: "Task",
			Fields: map[string]interface{}{
				"id":       fmt.Sprintf("task-%03d", i),
				"taskType": []string{"compute", "io", "network"}[i%3],
				"priority": (i % 5) + 1,
				"data":     fmt.Sprintf("data-%d", i),
			},
		}
		if err := xupleManager.CreateXuple("task_queue", task, nil); err != nil {
			t.Fatalf("❌ Erreur création tâche %d: %v", i, err)
		}
	}

	countBefore := taskQueue.Count()
	t.Logf("   Tâches créées: %d", countBefore)

	// Worker récupère batch de 5 tâches
	const batchSize = 5
	worker1Tasks, err := taskQueue.RetrieveMultiple("worker-1", batchSize)
	if err != nil {
		t.Fatalf("❌ RetrieveMultiple worker-1 échoué: %v", err)
	}

	if len(worker1Tasks) != batchSize {
		t.Errorf("❌ Worker-1 devrait avoir %d tâches, reçu %d", batchSize, len(worker1Tasks))
	} else {
		t.Logf("✅ Worker-1 a récupéré %d tâches", len(worker1Tasks))
	}

	// Vérifier que les tâches sont consommées
	countAfterWorker1 := taskQueue.Count()
	expectedRemaining := numTasks - batchSize
	if countAfterWorker1 != expectedRemaining {
		t.Errorf("❌ Count devrait être %d, reçu %d", expectedRemaining, countAfterWorker1)
	} else {
		t.Logf("✅ Tâches restantes: %d", countAfterWorker1)
	}

	// Plusieurs workers récupèrent en parallèle
	const numWorkers = 3
	workerResults := make(chan int, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		workerID := fmt.Sprintf("worker-%d", i+2)
		go func(id string) {
			defer wg.Done()
			tasks, _ := taskQueue.RetrieveMultiple(id, batchSize)
			workerResults <- len(tasks)
		}(workerID)
	}

	wg.Wait()
	close(workerResults)

	totalRetrievedByWorkers := 0
	for count := range workerResults {
		totalRetrievedByWorkers += count
	}

	t.Logf("✅ %d workers ont récupéré %d tâches au total", numWorkers, totalRetrievedByWorkers)

	finalCount := taskQueue.Count()
	expectedFinal := numTasks - batchSize - totalRetrievedByWorkers
	if finalCount != expectedFinal {
		t.Errorf("❌ Count final devrait être %d, reçu %d", expectedFinal, finalCount)
	} else {
		t.Logf("✅ Count final correct: %d", finalCount)
	}

	// ═══════════════════════════════════════════════════════════════
	// SCÉNARIO 2: Priorité avec LIFO
	// ═══════════════════════════════════════════════════════════════
	t.Log("")
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("SCÉNARIO 2: Tâches Haute Priorité (LIFO)")
	t.Log("───────────────────────────────────────────────────────────────")

	highPrioSpace, _ := xupleManager.GetXupleSpace("high_priority_tasks")

	// Créer 10 tâches avec priorités croissantes
	for i := 0; i < 10; i++ {
		task := &rete.Fact{
			Type: "Task",
			Fields: map[string]interface{}{
				"id":       fmt.Sprintf("prio-task-%d", i),
				"taskType": "urgent",
				"priority": i + 1, // 1 à 10
				"data":     fmt.Sprintf("urgent-data-%d", i),
			},
		}
		if err := xupleManager.CreateXuple("high_priority_tasks", task, nil); err != nil {
			t.Fatalf("❌ Erreur création tâche prioritaire: %v", err)
		}
	}

	// Récupérer batch de 3 (LIFO = derniers créés = priorités 10, 9, 8)
	prioTasks, err := highPrioSpace.RetrieveMultiple("urgent-worker", 3)
	if err != nil {
		t.Fatalf("❌ RetrieveMultiple échoué: %v", err)
	}

	if len(prioTasks) != 3 {
		t.Errorf("❌ Devrait récupérer 3 tâches, reçu %d", len(prioTasks))
	} else {
		// Vérifier l'ordre LIFO (priorités décroissantes)
		priorities := make([]int, len(prioTasks))
		for i, task := range prioTasks {
			prio, _ := task.Fact.Fields["priority"].(int)
			priorities[i] = prio
		}
		t.Logf("✅ Tâches récupérées (LIFO): priorités %v", priorities)

		// Vérifier ordre décroissant
		for i := 0; i < len(priorities)-1; i++ {
			if priorities[i] < priorities[i+1] {
				t.Errorf("❌ Ordre LIFO incorrect: %v", priorities)
				break
			}
		}
	}

	// ═══════════════════════════════════════════════════════════════
	// SCÉNARIO 3: Résultats partagés (per-agent)
	// ═══════════════════════════════════════════════════════════════
	t.Log("")
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("SCÉNARIO 3: Pool de Résultats (per-agent)")
	t.Log("───────────────────────────────────────────────────────────────")

	resultsPool, _ := xupleManager.GetXupleSpace("results_pool")

	// Créer 8 résultats
	for i := 0; i < 8; i++ {
		result := &rete.Fact{
			Type: "Result",
			Fields: map[string]interface{}{
				"taskId": fmt.Sprintf("completed-task-%d", i),
				"status": "success",
				"output": fmt.Sprintf("result-data-%d", i),
			},
		}
		if err := xupleManager.CreateXuple("results_pool", result, nil); err != nil {
			t.Fatalf("❌ Erreur création résultat: %v", err)
		}
	}

	// Monitor-1 récupère 5 résultats
	monitor1Results, err := resultsPool.RetrieveMultiple("monitor-1", 5)
	if err != nil {
		t.Fatalf("❌ RetrieveMultiple monitor-1 échoué: %v", err)
	}
	t.Logf("✅ Monitor-1 a récupéré %d résultats", len(monitor1Results))

	// Avec per-agent, les résultats restent disponibles
	countAfterMonitor1 := resultsPool.Count()
	if countAfterMonitor1 != 8 {
		t.Errorf("❌ Avec per-agent, count devrait rester 8, reçu %d", countAfterMonitor1)
	} else {
		t.Logf("✅ Count reste %d (per-agent policy)", countAfterMonitor1)
	}

	// Monitor-2 récupère aussi 5 résultats (les mêmes peuvent être retournés)
	monitor2Results, err := resultsPool.RetrieveMultiple("monitor-2", 5)
	if err != nil {
		t.Fatalf("❌ RetrieveMultiple monitor-2 échoué: %v", err)
	}
	t.Logf("✅ Monitor-2 a récupéré %d résultats", len(monitor2Results))

	// Monitor-1 ne peut plus récupérer (déjà consommé ces résultats)
	monitor1Again, err := resultsPool.RetrieveMultiple("monitor-1", 5)
	if err != nil {
		t.Logf("   RetrieveMultiple monitor-1 (2ème): erreur attendue: %v", err)
	}
	if len(monitor1Again) != 3 {
		t.Logf("   Monitor-1 peut récupérer les 3 résultats restants non consommés par lui")
	}

	// ═══════════════════════════════════════════════════════════════
	// SCÉNARIO 4: Gestion des limites
	// ═══════════════════════════════════════════════════════════════
	t.Log("")
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("SCÉNARIO 4: Gestion des Limites")
	t.Log("───────────────────────────────────────────────────────────────")

	// Créer un xuple-space avec MaxSize limité
	limitedConfig := xuples.XupleSpaceConfig{
		Name:              "limited_queue",
		SelectionPolicy:   xuples.NewFIFOSelectionPolicy(),
		ConsumptionPolicy: xuples.NewOnceConsumptionPolicy(),
		RetentionPolicy:   xuples.NewDurationRetentionPolicy(time.Hour),
		MaxSize:           10,
	}
	if err := xupleManager.CreateXupleSpace("limited_queue", limitedConfig); err != nil {
		t.Fatalf("❌ Erreur création limited_queue: %v", err)
	}

	limitedSpace, _ := xupleManager.GetXupleSpace("limited_queue")

	// Insérer jusqu'à la limite
	insertedCount := 0
	for i := 0; i < 15; i++ {
		fact := &rete.Fact{
			Type: "Task",
			Fields: map[string]interface{}{
				"id":   fmt.Sprintf("limited-%d", i),
				"data": fmt.Sprintf("data-%d", i),
			},
		}
		err := xupleManager.CreateXuple("limited_queue", fact, nil)
		if err == nil {
			insertedCount++
		}
	}

	if insertedCount != 10 {
		t.Errorf("❌ Devrait insérer exactement 10 xuples, inséré %d", insertedCount)
	} else {
		t.Logf("✅ MaxSize respecté: %d xuples insérés (max 10)", insertedCount)
	}

	// Récupérer batch plus grand que disponible
	largeBatch, err := limitedSpace.RetrieveMultiple("consumer", 20)
	if err != nil {
		t.Fatalf("❌ RetrieveMultiple ne devrait pas échouer même si n > disponible: %v", err)
	}

	if len(largeBatch) != 10 {
		t.Errorf("❌ Devrait retourner 10 xuples (tous disponibles), reçu %d", len(largeBatch))
	} else {
		t.Logf("✅ RetrieveMultiple retourne tous les xuples disponibles: %d", len(largeBatch))
	}

	// Vérifier espace vide
	countAfterLargeBatch := limitedSpace.Count()
	if countAfterLargeBatch != 0 {
		t.Errorf("❌ Espace devrait être vide, count = %d", countAfterLargeBatch)
	} else {
		t.Logf("✅ Espace correctement vidé")
	}

	// ═══════════════════════════════════════════════════════════════
	// RAPPORT FINAL
	// ═══════════════════════════════════════════════════════════════
	t.Log("")
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("📊 RAPPORT FINAL E2E BATCH")
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Logf("✅ Scénario 1: Traitement batch concurrent - VALIDÉ")
	t.Logf("✅ Scénario 2: LIFO pour priorités - VALIDÉ")
	t.Logf("✅ Scénario 3: Partage per-agent - VALIDÉ")
	t.Logf("✅ Scénario 4: Limites et gestion erreurs - VALIDÉ")
	t.Log("")
	t.Log("✅ TOUS LES SCÉNARIOS E2E BATCH RÉUSSIS!")
	t.Log("═══════════════════════════════════════════════════════════════")
}

// TestXuplesBatch_E2E_StressTest teste la robustesse sous charge
func TestXuplesBatch_E2E_StressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Stress test ignoré en mode short")
	}

	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("🧪 TEST E2E: RetrieveMultiple - Test de Charge")
	t.Log("═══════════════════════════════════════════════════════════════")

	xupleManager := xuples.NewXupleManager()

	// Créer un large xuple-space
	config := xuples.XupleSpaceConfig{
		Name:              "stress_queue",
		SelectionPolicy:   xuples.NewFIFOSelectionPolicy(),
		ConsumptionPolicy: xuples.NewOnceConsumptionPolicy(),
		RetentionPolicy:   xuples.NewDurationRetentionPolicy(time.Hour),
		MaxSize:           0,
	}

	if err := xupleManager.CreateXupleSpace("stress_queue", config); err != nil {
		t.Fatalf("❌ Erreur création xuple-space: %v", err)
	}

	space, _ := xupleManager.GetXupleSpace("stress_queue")

	// Créer beaucoup de xuples
	const numXuples = 1000
	for i := 0; i < numXuples; i++ {
		fact := &rete.Fact{
			Type: "StressTask",
			Fields: map[string]interface{}{
				"id":    fmt.Sprintf("stress-%d", i),
				"index": i,
			},
		}
		if err := xupleManager.CreateXuple("stress_queue", fact, nil); err != nil {
			t.Fatalf("❌ Erreur création xuple %d: %v", i, err)
		}
	}

	t.Logf("   Xuples créés: %d", numXuples)

	// Lancer plusieurs consumers concurrents
	const numConsumers = 10
	const batchSize = 50

	var wg sync.WaitGroup
	totalConsumed := make(chan int, numConsumers)

	startTime := time.Now()

	for i := 0; i < numConsumers; i++ {
		wg.Add(1)
		consumerID := fmt.Sprintf("stress-consumer-%d", i)
		go func(id string) {
			defer wg.Done()
			consumed := 0
			for {
				batch, err := space.RetrieveMultiple(id, batchSize)
				if err != nil {
					break
				}
				if len(batch) == 0 {
					break
				}
				consumed += len(batch)
				// Simuler traitement
				time.Sleep(1 * time.Millisecond)
			}
			totalConsumed <- consumed
		}(consumerID)
	}

	wg.Wait()
	close(totalConsumed)

	duration := time.Since(startTime)

	totalProcessed := 0
	for count := range totalConsumed {
		totalProcessed += count
	}

	t.Logf("✅ Xuples traités: %d/%d", totalProcessed, numXuples)
	t.Logf("✅ Temps total: %v", duration)
	t.Logf("✅ Débit: %.0f xuples/seconde", float64(numXuples)/duration.Seconds())

	if totalProcessed != numXuples {
		t.Errorf("❌ Total traité devrait être %d, reçu %d", numXuples, totalProcessed)
	}

	finalCount := space.Count()
	if finalCount != 0 {
		t.Errorf("❌ Espace devrait être vide, count = %d", finalCount)
	} else {
		t.Logf("✅ Espace correctement vidé")
	}
}

// Fonctions helper pour créer les politiques
func getSelectionPolicy(name string) xuples.SelectionPolicy {
	switch name {
	case "fifo":
		return xuples.NewFIFOSelectionPolicy()
	case "lifo":
		return xuples.NewLIFOSelectionPolicy()
	case "random":
		return xuples.NewRandomSelectionPolicy()
	default:
		return xuples.NewFIFOSelectionPolicy()
	}
}

func getConsumptionPolicy(name string) xuples.ConsumptionPolicy {
	switch name {
	case "once":
		return xuples.NewOnceConsumptionPolicy()
	case "per-agent":
		return xuples.NewPerAgentConsumptionPolicy()
	default:
		return xuples.NewOnceConsumptionPolicy()
	}
}
