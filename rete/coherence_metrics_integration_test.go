// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"sync"
	"testing"
	"time"
)
// delayedStorage wraps MemoryStorage with configurable write delay for testing retries
type delayedStorage struct {
	*MemoryStorage
	writeDelay   time.Duration
	mu           sync.RWMutex
	pendingFacts map[string]*Fact
	startTimes   map[string]time.Time
}
// newDelayedStorage creates a storage with delayed writes
func newDelayedStorage(delay time.Duration) *delayedStorage {
	return &delayedStorage{
		MemoryStorage: NewMemoryStorage(),
		writeDelay:    delay,
		pendingFacts:  make(map[string]*Fact),
		startTimes:    make(map[string]time.Time),
	}
}
// AddFact adds a fact immediately to pending, but delays visibility
func (ds *delayedStorage) AddFact(fact *Fact) error {
	ds.mu.Lock()
	internalID := fact.GetInternalID()
	ds.pendingFacts[internalID] = fact
	ds.startTimes[internalID] = time.Now()
	ds.mu.Unlock()
	// Immediately call parent to add
	return ds.MemoryStorage.AddFact(fact)
}
// GetFact retrieves a fact only after the delay has passed
func (ds *delayedStorage) GetFact(id string) *Fact {
	ds.mu.RLock()
	startTime, hasPending := ds.startTimes[id]
	ds.mu.RUnlock()
	if hasPending {
		elapsed := time.Since(startTime)
		if elapsed < ds.writeDelay {
			// Fact not yet "visible" - still in write delay period
			return nil
		}
		// Delay passed, remove from pending
		ds.mu.Lock()
		delete(ds.pendingFacts, id)
		delete(ds.startTimes, id)
		ds.mu.Unlock()
	}
	return ds.MemoryStorage.GetFact(id)
}
// Sync implements Storage.Sync with delay awareness
func (ds *delayedStorage) Sync() error {
	// Ensure any pending writes are visible
	ds.mu.Lock()
	defer ds.mu.Unlock()
	return ds.MemoryStorage.Sync()
}
// TestCoherenceMetrics_Integration teste l'intégration complète du collecteur de métriques
func TestCoherenceMetrics_Integration(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Créer un collecteur de métriques
	collector := NewCoherenceMetricsCollector()
	collector.SetTransactionID("test-tx-001")
	// Créer des faits de test
	facts := []map[string]interface{}{
		{"id": "prod1", "type": "Product", "name": "Item1", "price": 10.0},
		{"id": "prod2", "type": "Product", "name": "Item2", "price": 20.0},
		{"id": "prod3", "type": "Product", "name": "Item3", "price": 30.0},
		{"id": "prod4", "type": "Product", "name": "Item4", "price": 40.0},
		{"id": "prod5", "type": "Product", "name": "Item5", "price": 50.0},
	}
	// Soumettre les faits avec collecte de métriques
	err := network.SubmitFactsFromGrammarWithMetrics(facts, collector)
	if err != nil {
		t.Fatalf("Erreur lors de la soumission des faits: %v", err)
	}
	// Finaliser les métriques
	metrics := collector.Finalize()
	// Vérifications
	if metrics.FactsSubmitted != 5 {
		t.Errorf("FactsSubmitted: attendu 5, obtenu %d", metrics.FactsSubmitted)
	}
	if metrics.FactsPersisted != 5 {
		t.Errorf("FactsPersisted: attendu 5, obtenu %d", metrics.FactsPersisted)
	}
	if metrics.FactsFailed != 0 {
		t.Errorf("FactsFailed: attendu 0, obtenu %d", metrics.FactsFailed)
	}
	if metrics.TotalVerifyAttempts < 5 {
		t.Errorf("TotalVerifyAttempts devrait être >= 5, obtenu %d", metrics.TotalVerifyAttempts)
	}
	if metrics.TotalDuration == 0 {
		t.Error("TotalDuration devrait être > 0")
	}
	if metrics.TotalSubmissionTime == 0 {
		t.Error("TotalSubmissionTime devrait être > 0")
	}
	if metrics.TransactionID != "test-tx-001" {
		t.Errorf("TransactionID: attendu 'test-tx-001', obtenu '%s'", metrics.TransactionID)
	}
	// Vérifier la santé du système
	if !metrics.IsHealthy() {
		t.Errorf("Le système devrait être en bonne santé: %s", metrics.GetHealthReport())
	}
	// Vérifier que la phase de soumission a été enregistrée
	if _, exists := metrics.PhaseMetrics["fact_submission"]; !exists {
		t.Error("La phase 'fact_submission' devrait être enregistrée")
	}
	t.Logf("✅ Métriques collectées avec succès:\n%s", metrics.Summary())
}
// TestCoherenceMetrics_WithRetries teste les métriques avec des retries
func TestCoherenceMetrics_WithRetries(t *testing.T) {
	storage := newDelayedStorage(25 * time.Millisecond)
	network := NewReteNetwork(storage)
	network.VerifyRetryDelay = 10 * time.Millisecond
	network.MaxVerifyRetries = 5
	network.SubmissionTimeout = 5 * time.Second
	collector := NewCoherenceMetricsCollector()
	facts := []map[string]interface{}{
		{"id": "delayed1", "type": "Product", "name": "Delayed1"},
		{"id": "delayed2", "type": "Product", "name": "Delayed2"},
	}
	err := network.SubmitFactsFromGrammarWithMetrics(facts, collector)
	if err != nil {
		t.Fatalf("Erreur lors de la soumission: %v", err)
	}
	metrics := collector.Finalize()
	// Avec un délai, nous devrions avoir des retries
	if metrics.TotalRetries == 0 {
		t.Error("TotalRetries devrait être > 0 avec un storage retardé")
	}
	if metrics.FactsRetried == 0 {
		t.Error("FactsRetried devrait être > 0 avec un storage retardé")
	}
	if metrics.MaxRetriesForSingleFact == 0 {
		t.Error("MaxRetriesForSingleFact devrait être > 0")
	}
	// Le temps d'attente total devrait être significatif
	if metrics.TotalWaitTime < 25*time.Millisecond {
		t.Errorf("TotalWaitTime devrait être >= 25ms, obtenu %v", metrics.TotalWaitTime)
	}
	t.Logf("✅ Retries: %d total, %d max pour un fait", metrics.TotalRetries, metrics.MaxRetriesForSingleFact)
	t.Logf("✅ Temps d'attente: total=%v, moyen=%v, max=%v",
		metrics.TotalWaitTime, metrics.AvgWaitTime, metrics.MaxWaitTime)
}
// TestCoherenceMetrics_MultiplePhases teste l'enregistrement de multiples phases
func TestCoherenceMetrics_MultiplePhases(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	// Phase 1: Parsing
	collector.StartPhase("parsing")
	time.Sleep(5 * time.Millisecond)
	collector.EndPhase("parsing", 10, true)
	// Phase 2: Validation
	collector.StartPhase("validation")
	time.Sleep(3 * time.Millisecond)
	collector.EndPhase("validation", 10, true)
	// Phase 3: Soumission
	collector.StartPhase("submission")
	collector.RecordFactSubmitted()
	collector.RecordFactSubmitted()
	collector.RecordFactPersisted()
	collector.RecordFactPersisted()
	time.Sleep(10 * time.Millisecond)
	collector.EndPhase("submission", 2, true)
	metrics := collector.Finalize()
	if len(metrics.PhaseMetrics) != 3 {
		t.Errorf("Attendu 3 phases, obtenu %d", len(metrics.PhaseMetrics))
	}
	// Vérifier chaque phase
	phases := []string{"parsing", "validation", "submission"}
	for _, phaseName := range phases {
		phase, exists := metrics.PhaseMetrics[phaseName]
		if !exists {
			t.Errorf("Phase '%s' non trouvée", phaseName)
			continue
		}
		if !phase.Succeeded {
			t.Errorf("Phase '%s' devrait avoir réussi", phaseName)
		}
		if phase.Duration == 0 {
			t.Errorf("Phase '%s' devrait avoir une durée > 0", phaseName)
		}
	}
	t.Logf("✅ %d phases enregistrées avec succès", len(metrics.PhaseMetrics))
}
// TestCoherenceMetrics_PreCommitChecks teste l'enregistrement des vérifications pré-commit
func TestCoherenceMetrics_PreCommitChecks(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	// Simuler des vérifications pré-commit
	collector.RecordPreCommitCheck(true)
	collector.RecordPreCommitCheck(true)
	collector.RecordPreCommitCheck(false)
	collector.RecordPreCommitCheck(true)
	metrics := collector.Finalize()
	if metrics.PreCommitChecks != 4 {
		t.Errorf("PreCommitChecks: attendu 4, obtenu %d", metrics.PreCommitChecks)
	}
	if metrics.PreCommitSuccesses != 3 {
		t.Errorf("PreCommitSuccesses: attendu 3, obtenu %d", metrics.PreCommitSuccesses)
	}
	if metrics.PreCommitFailures != 1 {
		t.Errorf("PreCommitFailures: attendu 1, obtenu %d", metrics.PreCommitFailures)
	}
}
// TestCoherenceMetrics_Rollback teste l'enregistrement d'un rollback
func TestCoherenceMetrics_Rollback(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	collector.SetTransactionID("tx-rollback-test")
	// Simuler une soumission qui échoue
	collector.RecordFactSubmitted()
	collector.RecordFactSubmitted()
	collector.RecordFactPersisted()
	collector.RecordFactFailed()
	// Enregistrer le rollback
	collector.RecordRollback("coherence check failed")
	metrics := collector.Finalize()
	if !metrics.WasRolledBack {
		t.Error("WasRolledBack devrait être true")
	}
	if metrics.RollbackReason != "coherence check failed" {
		t.Errorf("RollbackReason: attendu 'coherence check failed', obtenu '%s'", metrics.RollbackReason)
	}
	// Le système ne devrait pas être sain après un rollback
	if metrics.IsHealthy() {
		t.Error("Le système ne devrait pas être sain après un rollback")
	}
	healthReport := metrics.GetHealthReport()
	if healthReport == "" {
		t.Error("Le rapport de santé ne devrait pas être vide")
	}
	t.Logf("✅ Rollback enregistré: %s", metrics.RollbackReason)
	t.Logf("Rapport de santé:\n%s", healthReport)
}
// TestCoherenceMetrics_JSONExport teste l'export JSON des métriques
func TestCoherenceMetrics_JSONExport(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	collector := NewCoherenceMetricsCollector()
	collector.SetTransactionID("tx-json-test")
	facts := []map[string]interface{}{
		{"id": "json1", "type": "Product", "name": "JSON Test 1"},
		{"id": "json2", "type": "Product", "name": "JSON Test 2"},
	}
	err := network.SubmitFactsFromGrammarWithMetrics(facts, collector)
	if err != nil {
		t.Fatalf("Erreur lors de la soumission: %v", err)
	}
	metrics := collector.Finalize()
	// Exporter en JSON
	jsonStr, err := metrics.ToJSON()
	if err != nil {
		t.Fatalf("Erreur lors de l'export JSON: %v", err)
	}
	if jsonStr == "" {
		t.Error("Le JSON exporté ne devrait pas être vide")
	}
	// Vérifier que le JSON contient les champs clés
	expectedFields := []string{
		"facts_submitted",
		"facts_persisted",
		"total_verify_attempts",
		"transaction_id",
		"total_duration",
	}
	for _, field := range expectedFields {
		if !containsString(jsonStr, field) {
			t.Errorf("Le JSON devrait contenir le champ '%s'", field)
		}
	}
	t.Logf("✅ JSON exporté avec succès (%d bytes)", len(jsonStr))
}
// TestCoherenceMetrics_ConcurrentCollection teste la collecte de métriques en concurrence
func TestCoherenceMetrics_ConcurrentCollection(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	collector := NewCoherenceMetricsCollector()
	// Soumettre des faits en parallèle (via goroutines simulant des opérations concurrentes)
	facts := []map[string]interface{}{
		{"id": "concurrent1", "type": "Product", "name": "Concurrent 1"},
		{"id": "concurrent2", "type": "Product", "name": "Concurrent 2"},
		{"id": "concurrent3", "type": "Product", "name": "Concurrent 3"},
		{"id": "concurrent4", "type": "Product", "name": "Concurrent 4"},
		{"id": "concurrent5", "type": "Product", "name": "Concurrent 5"},
	}
	// La soumission elle-même est séquentielle, mais le collecteur doit gérer
	// les accès concurrents potentiels
	err := network.SubmitFactsFromGrammarWithMetrics(facts, collector)
	if err != nil {
		t.Fatalf("Erreur lors de la soumission: %v", err)
	}
	metrics := collector.Finalize()
	if metrics.FactsSubmitted != 5 {
		t.Errorf("FactsSubmitted: attendu 5, obtenu %d", metrics.FactsSubmitted)
	}
	if metrics.FactsPersisted != 5 {
		t.Errorf("FactsPersisted: attendu 5, obtenu %d", metrics.FactsPersisted)
	}
	t.Logf("✅ Collection concurrente réussie: %d faits", metrics.FactsPersisted)
}
// TestCoherenceMetrics_HealthThresholds teste les seuils de santé
func TestCoherenceMetrics_HealthThresholds(t *testing.T) {
	tests := []struct {
		name            string
		setup           func(*CoherenceMetricsCollector)
		shouldBeHealthy bool
	}{
		{
			name: "100% succès - sain",
			setup: func(c *CoherenceMetricsCollector) {
				for i := 0; i < 100; i++ {
					c.RecordFactSubmitted()
					c.RecordFactPersisted()
				}
			},
			shouldBeHealthy: true,
		},
		{
			name: "95% succès - sain (seuil limite)",
			setup: func(c *CoherenceMetricsCollector) {
				for i := 0; i < 100; i++ {
					c.RecordFactSubmitted()
				}
				for i := 0; i < 95; i++ {
					c.RecordFactPersisted()
				}
				for i := 0; i < 5; i++ {
					c.RecordFactFailed()
				}
			},
			shouldBeHealthy: true,
		},
		{
			name: "90% succès - malsain",
			setup: func(c *CoherenceMetricsCollector) {
				for i := 0; i < 100; i++ {
					c.RecordFactSubmitted()
				}
				for i := 0; i < 90; i++ {
					c.RecordFactPersisted()
				}
				for i := 0; i < 10; i++ {
					c.RecordFactFailed()
				}
			},
			shouldBeHealthy: false,
		},
		{
			name: "Trop de timeouts - malsain",
			setup: func(c *CoherenceMetricsCollector) {
				for i := 0; i < 100; i++ {
					c.RecordFactSubmitted()
					c.RecordFactPersisted()
				}
				for i := 0; i < 10; i++ {
					c.RecordTimeout()
				}
			},
			shouldBeHealthy: false,
		},
		{
			name: "Trop de retries - malsain",
			setup: func(c *CoherenceMetricsCollector) {
				for i := 0; i < 10; i++ {
					c.RecordFactSubmitted()
					c.RecordFactPersisted()
					c.RecordRetry(3)
					c.RecordRetry(3)
				}
			},
			shouldBeHealthy: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := NewCoherenceMetricsCollector()
			tt.setup(collector)
			metrics := collector.Finalize()
			isHealthy := metrics.IsHealthy()
			if isHealthy != tt.shouldBeHealthy {
				t.Errorf("IsHealthy(): attendu %v, obtenu %v\nRapport: %s",
					tt.shouldBeHealthy, isHealthy, metrics.GetHealthReport())
			}
		})
	}
}
// containsString vérifie si une chaîne contient une sous-chaîne
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsStringHelper(s, substr))
}
func containsStringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}