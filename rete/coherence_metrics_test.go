// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"
)
// TestNewCoherenceMetricsCollector vérifie la création d'un nouveau collecteur
func TestNewCoherenceMetricsCollector(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	if collector == nil {
		t.Fatal("Le collecteur ne devrait pas être nil")
	}
	if collector.metrics == nil {
		t.Fatal("Les métriques ne devraient pas être nil")
	}
	if collector.metrics.PhaseMetrics == nil {
		t.Fatal("PhaseMetrics ne devrait pas être nil")
	}
	if collector.activePhases == nil {
		t.Fatal("activePhases ne devrait pas être nil")
	}
	if collector.metrics.StartTime.IsZero() {
		t.Error("StartTime devrait être initialisé")
	}
}
// TestRecordFactOperations teste l'enregistrement des opérations sur les faits
func TestRecordFactOperations(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	// Enregistrer des opérations
	collector.RecordFactSubmitted()
	collector.RecordFactSubmitted()
	collector.RecordFactSubmitted()
	collector.RecordFactPersisted()
	collector.RecordFactPersisted()
	collector.RecordFactRetried()
	collector.RecordFactFailed()
	collector.RecordFactPropagated()
	collector.RecordFactPropagated()
	collector.RecordFactPropagated()
	collector.RecordFactPropagated()
	collector.RecordTerminalActivation()
	collector.RecordTerminalActivation()
	// Vérifier les compteurs
	metrics := collector.GetMetrics()
	if metrics.FactsSubmitted != 3 {
		t.Errorf("FactsSubmitted: attendu 3, obtenu %d", metrics.FactsSubmitted)
	}
	if metrics.FactsPersisted != 2 {
		t.Errorf("FactsPersisted: attendu 2, obtenu %d", metrics.FactsPersisted)
	}
	if metrics.FactsRetried != 1 {
		t.Errorf("FactsRetried: attendu 1, obtenu %d", metrics.FactsRetried)
	}
	if metrics.FactsFailed != 1 {
		t.Errorf("FactsFailed: attendu 1, obtenu %d", metrics.FactsFailed)
	}
	if metrics.FactsPropagated != 4 {
		t.Errorf("FactsPropagated: attendu 4, obtenu %d", metrics.FactsPropagated)
	}
	if metrics.TerminalActivations != 2 {
		t.Errorf("TerminalActivations: attendu 2, obtenu %d", metrics.TerminalActivations)
	}
}
// TestRecordSynchronizationMetrics teste l'enregistrement des métriques de synchronisation
func TestRecordSynchronizationMetrics(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	// Enregistrer des tentatives et retries
	collector.RecordVerifyAttempt()
	collector.RecordVerifyAttempt()
	collector.RecordVerifyAttempt()
	collector.RecordVerifyAttempt()
	collector.RecordVerifyAttempt()
	collector.RecordRetry(1)
	collector.RecordRetry(2)
	collector.RecordRetry(3)
	collector.RecordRetry(1)
	collector.RecordTimeout()
	collector.RecordTimeout()
	// Vérifier les compteurs
	metrics := collector.GetMetrics()
	if metrics.TotalVerifyAttempts != 5 {
		t.Errorf("TotalVerifyAttempts: attendu 5, obtenu %d", metrics.TotalVerifyAttempts)
	}
	if metrics.TotalRetries != 4 {
		t.Errorf("TotalRetries: attendu 4, obtenu %d", metrics.TotalRetries)
	}
	if metrics.MaxRetriesForSingleFact != 3 {
		t.Errorf("MaxRetriesForSingleFact: attendu 3, obtenu %d", metrics.MaxRetriesForSingleFact)
	}
	if metrics.TotalTimeouts != 2 {
		t.Errorf("TotalTimeouts: attendu 2, obtenu %d", metrics.TotalTimeouts)
	}
}
// TestRecordWaitTime teste l'enregistrement des temps d'attente
func TestRecordWaitTime(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	// Enregistrer des temps d'attente
	collector.RecordWaitTime(10 * time.Millisecond)
	collector.RecordWaitTime(20 * time.Millisecond)
	collector.RecordWaitTime(5 * time.Millisecond)
	collector.RecordWaitTime(30 * time.Millisecond)
	metrics := collector.GetMetrics()
	expectedTotal := 65 * time.Millisecond
	if metrics.TotalWaitTime != expectedTotal {
		t.Errorf("TotalWaitTime: attendu %v, obtenu %v", expectedTotal, metrics.TotalWaitTime)
	}
	if metrics.MaxWaitTime != 30*time.Millisecond {
		t.Errorf("MaxWaitTime: attendu 30ms, obtenu %v", metrics.MaxWaitTime)
	}
	if metrics.MinWaitTime != 5*time.Millisecond {
		t.Errorf("MinWaitTime: attendu 5ms, obtenu %v", metrics.MinWaitTime)
	}
}
// TestRecordTimings teste l'enregistrement des temps système
func TestRecordTimings(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	syncTime := 100 * time.Millisecond
	submissionTime := 500 * time.Millisecond
	collector.RecordSyncTime(syncTime)
	collector.RecordSubmissionTime(submissionTime)
	metrics := collector.GetMetrics()
	if metrics.TotalSyncTime != syncTime {
		t.Errorf("TotalSyncTime: attendu %v, obtenu %v", syncTime, metrics.TotalSyncTime)
	}
	if metrics.TotalSubmissionTime != submissionTime {
		t.Errorf("TotalSubmissionTime: attendu %v, obtenu %v", submissionTime, metrics.TotalSubmissionTime)
	}
}
// TestRecordPreCommitCheck teste l'enregistrement des vérifications pré-commit
func TestRecordPreCommitCheck(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	collector.RecordPreCommitCheck(true)
	collector.RecordPreCommitCheck(true)
	collector.RecordPreCommitCheck(false)
	collector.RecordPreCommitCheck(true)
	metrics := collector.GetMetrics()
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
// TestPhaseMetrics teste la gestion des métriques par phase
func TestPhaseMetrics(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	// Phase 1: Succès
	collector.StartPhase("parsing")
	time.Sleep(10 * time.Millisecond)
	collector.EndPhase("parsing", 5, true)
	// Phase 2: Échec
	collector.StartPhase("validation")
	time.Sleep(5 * time.Millisecond)
	collector.EndPhase("validation", 3, false)
	metrics := collector.GetMetrics()
	if len(metrics.PhaseMetrics) != 2 {
		t.Fatalf("PhaseMetrics: attendu 2 phases, obtenu %d", len(metrics.PhaseMetrics))
	}
	// Vérifier phase parsing
	parsingPhase := metrics.PhaseMetrics["parsing"]
	if parsingPhase == nil {
		t.Fatal("Phase 'parsing' non trouvée")
	}
	if parsingPhase.PhaseName != "parsing" {
		t.Errorf("PhaseName: attendu 'parsing', obtenu '%s'", parsingPhase.PhaseName)
	}
	if parsingPhase.ItemsProcessed != 5 {
		t.Errorf("ItemsProcessed: attendu 5, obtenu %d", parsingPhase.ItemsProcessed)
	}
	if !parsingPhase.Succeeded {
		t.Error("La phase parsing devrait avoir réussi")
	}
	if parsingPhase.Duration < 10*time.Millisecond {
		t.Errorf("Duration devrait être >= 10ms, obtenu %v", parsingPhase.Duration)
	}
	// Vérifier phase validation
	validationPhase := metrics.PhaseMetrics["validation"]
	if validationPhase == nil {
		t.Fatal("Phase 'validation' non trouvée")
	}
	if validationPhase.Succeeded {
		t.Error("La phase validation ne devrait pas avoir réussi")
	}
}
// TestTransactionTracking teste le suivi des transactions
func TestTransactionTracking(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	txID := "tx-12345"
	rollbackReason := "test rollback"
	collector.SetTransactionID(txID)
	collector.RecordRollback(rollbackReason)
	metrics := collector.GetMetrics()
	if metrics.TransactionID != txID {
		t.Errorf("TransactionID: attendu '%s', obtenu '%s'", txID, metrics.TransactionID)
	}
	if !metrics.WasRolledBack {
		t.Error("WasRolledBack devrait être true")
	}
	if metrics.RollbackReason != rollbackReason {
		t.Errorf("RollbackReason: attendu '%s', obtenu '%s'", rollbackReason, metrics.RollbackReason)
	}
}
// TestFinalize teste la finalisation des métriques
func TestFinalize(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	// Ajouter des données
	collector.RecordFactSubmitted()
	collector.RecordFactSubmitted()
	collector.RecordFactPersisted()
	collector.RecordWaitTime(10 * time.Millisecond)
	collector.RecordWaitTime(20 * time.Millisecond)
	time.Sleep(50 * time.Millisecond)
	// Finaliser
	metrics := collector.Finalize()
	if metrics.EndTime.IsZero() {
		t.Error("EndTime devrait être défini")
	}
	if metrics.EndTime.Before(metrics.StartTime) {
		t.Error("EndTime devrait être après StartTime")
	}
	if metrics.TotalDuration < 50*time.Millisecond {
		t.Errorf("TotalDuration devrait être >= 50ms, obtenu %v", metrics.TotalDuration)
	}
	// Vérifier le temps d'attente moyen (1 fact persisté)
	expectedAvg := 30 * time.Millisecond // (10 + 20) / 1
	if metrics.AvgWaitTime != expectedAvg {
		t.Errorf("AvgWaitTime: attendu %v, obtenu %v", expectedAvg, metrics.AvgWaitTime)
	}
}
// TestCoherenceMetricsConcurrentAccess teste l'accès concurrent aux métriques
func TestCoherenceMetricsConcurrentAccess(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	var wg sync.WaitGroup
	operations := 100
	// Lancer plusieurs goroutines pour enregistrer des métriques
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				collector.RecordFactSubmitted()
				collector.RecordFactPersisted()
				collector.RecordVerifyAttempt()
				collector.RecordWaitTime(time.Millisecond)
			}
		}()
	}
	wg.Wait()
	metrics := collector.GetMetrics()
	expectedCount := 10 * operations
	if metrics.FactsSubmitted != expectedCount {
		t.Errorf("FactsSubmitted: attendu %d, obtenu %d", expectedCount, metrics.FactsSubmitted)
	}
	if metrics.FactsPersisted != expectedCount {
		t.Errorf("FactsPersisted: attendu %d, obtenu %d", expectedCount, metrics.FactsPersisted)
	}
	if metrics.TotalVerifyAttempts != expectedCount {
		t.Errorf("TotalVerifyAttempts: attendu %d, obtenu %d", expectedCount, metrics.TotalVerifyAttempts)
	}
}
// TestJSONExport teste l'export JSON des métriques
func TestJSONExport(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	collector.RecordFactSubmitted()
	collector.RecordFactPersisted()
	collector.RecordRetry(2)
	collector.SetTransactionID("tx-test")
	metrics := collector.Finalize()
	jsonStr, err := metrics.ToJSON()
	if err != nil {
		t.Fatalf("Erreur lors de l'export JSON: %v", err)
	}
	// Vérifier que c'est du JSON valide
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Fatalf("JSON invalide: %v", err)
	}
	// Vérifier quelques champs
	if parsed["facts_submitted"] != float64(1) {
		t.Errorf("facts_submitted incorrect dans JSON")
	}
	if parsed["transaction_id"] != "tx-test" {
		t.Errorf("transaction_id incorrect dans JSON")
	}
}
// TestSummary teste la génération du résumé
func TestSummary(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	collector.RecordFactSubmitted()
	collector.RecordFactSubmitted()
	collector.RecordFactPersisted()
	collector.RecordFactPersisted()
	collector.RecordRetry(1)
	collector.RecordTimeout()
	collector.RecordWaitTime(15 * time.Millisecond)
	metrics := collector.Finalize()
	summary := metrics.Summary()
	// Vérifier que le résumé contient les informations clés
	if !strings.Contains(summary, "2/2") {
		t.Error("Le résumé devrait contenir '2/2 faits persistés'")
	}
	if !strings.Contains(summary, "100.0%") {
		t.Error("Le résumé devrait contenir '100.0%'")
	}
	if !strings.Contains(summary, "1 retries") {
		t.Error("Le résumé devrait mentionner les retries")
	}
}
// TestIsHealthy teste la détection de santé du système
func TestIsHealthy(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*CoherenceMetricsCollector)
		expected bool
	}{
		{
			name: "Système sain - 100% succès",
			setup: func(c *CoherenceMetricsCollector) {
				for i := 0; i < 10; i++ {
					c.RecordFactSubmitted()
					c.RecordFactPersisted()
				}
			},
			expected: true,
		},
		{
			name: "Système sain - 95% succès",
			setup: func(c *CoherenceMetricsCollector) {
				for i := 0; i < 20; i++ {
					c.RecordFactSubmitted()
				}
				for i := 0; i < 19; i++ {
					c.RecordFactPersisted()
				}
			},
			expected: true,
		},
		{
			name: "Système malsain - faible taux de succès",
			setup: func(c *CoherenceMetricsCollector) {
				for i := 0; i < 10; i++ {
					c.RecordFactSubmitted()
				}
				for i := 0; i < 5; i++ {
					c.RecordFactPersisted()
				}
			},
			expected: false,
		},
		{
			name: "Système malsain - trop de timeouts",
			setup: func(c *CoherenceMetricsCollector) {
				for i := 0; i < 10; i++ {
					c.RecordFactSubmitted()
					c.RecordFactPersisted()
				}
				for i := 0; i < 1; i++ {
					c.RecordTimeout()
				}
			},
			expected: false,
		},
		{
			name: "Système malsain - trop de retries",
			setup: func(c *CoherenceMetricsCollector) {
				for i := 0; i < 10; i++ {
					c.RecordFactSubmitted()
					c.RecordFactPersisted()
					c.RecordRetry(2)
					c.RecordRetry(2)
				}
			},
			expected: false,
		},
		{
			name: "Système malsain - rollback avec raison",
			setup: func(c *CoherenceMetricsCollector) {
				for i := 0; i < 10; i++ {
					c.RecordFactSubmitted()
					c.RecordFactPersisted()
				}
				c.RecordRollback("coherence failure")
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := NewCoherenceMetricsCollector()
			tt.setup(collector)
			metrics := collector.Finalize()
			result := metrics.IsHealthy()
			if result != tt.expected {
				t.Errorf("IsHealthy(): attendu %v, obtenu %v", tt.expected, result)
			}
		})
	}
}
// TestGetHealthReport teste la génération du rapport de santé
func TestGetHealthReport(t *testing.T) {
	// Système sain
	collector := NewCoherenceMetricsCollector()
	for i := 0; i < 10; i++ {
		collector.RecordFactSubmitted()
		collector.RecordFactPersisted()
	}
	metrics := collector.Finalize()
	report := metrics.GetHealthReport()
	if !strings.Contains(report, "bonne santé") {
		t.Error("Le rapport devrait indiquer une bonne santé")
	}
	// Système avec problèmes
	collector2 := NewCoherenceMetricsCollector()
	for i := 0; i < 10; i++ {
		collector2.RecordFactSubmitted()
	}
	for i := 0; i < 5; i++ {
		collector2.RecordFactPersisted()
	}
	collector2.RecordTimeout()
	collector2.RecordRollback("test failure")
	metrics2 := collector2.Finalize()
	report2 := metrics2.GetHealthReport()
	if !strings.Contains(report2, "Problèmes détectés") {
		t.Error("Le rapport devrait mentionner des problèmes")
	}
	if !strings.Contains(report2, "Taux de succès bas") {
		t.Error("Le rapport devrait mentionner le taux de succès bas")
	}
	if !strings.Contains(report2, "Rollback") {
		t.Error("Le rapport devrait mentionner le rollback")
	}
}
// TestStringFormatting teste le formatage de la chaîne complète
func TestStringFormatting(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	collector.RecordFactSubmitted()
	collector.RecordFactPersisted()
	collector.RecordRetry(1)
	collector.RecordWaitTime(10 * time.Millisecond)
	collector.SetTransactionID("tx-format-test")
	collector.StartPhase("test_phase")
	time.Sleep(5 * time.Millisecond)
	collector.EndPhase("test_phase", 1, true)
	metrics := collector.Finalize()
	str := metrics.String()
	// Vérifier que la sortie contient les sections principales
	requiredSections := []string{
		"Métriques de Cohérence RETE",
		"Faits:",
		"Synchronisation:",
		"Temps d'attente:",
		"Temps système:",
		"Cohérence pré-commit:",
		"Transaction:",
		"Phases:",
	}
	for _, section := range requiredSections {
		if !strings.Contains(str, section) {
			t.Errorf("La sortie devrait contenir la section '%s'", section)
		}
	}
	// Vérifier que les valeurs sont présentes
	if !strings.Contains(str, "tx-format-test") {
		t.Error("La sortie devrait contenir l'ID de transaction")
	}
	if !strings.Contains(str, "test_phase") {
		t.Error("La sortie devrait contenir le nom de la phase")
	}
}
// TestAvgWaitTimeCalculation teste le calcul du temps d'attente moyen
func TestAvgWaitTimeCalculation(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	// Cas 1: Aucun fait persisté
	metrics1 := collector.Finalize()
	if metrics1.AvgWaitTime != 0 {
		t.Error("AvgWaitTime devrait être 0 quand aucun fait n'est persisté")
	}
	// Cas 2: Plusieurs faits avec différents temps d'attente
	collector2 := NewCoherenceMetricsCollector()
	collector2.RecordFactPersisted()
	collector2.RecordFactPersisted()
	collector2.RecordFactPersisted()
	collector2.RecordWaitTime(10 * time.Millisecond)
	collector2.RecordWaitTime(20 * time.Millisecond)
	collector2.RecordWaitTime(30 * time.Millisecond)
	metrics2 := collector2.Finalize()
	expectedAvg := 20 * time.Millisecond // (10 + 20 + 30) / 3
	if metrics2.AvgWaitTime != expectedAvg {
		t.Errorf("AvgWaitTime: attendu %v, obtenu %v", expectedAvg, metrics2.AvgWaitTime)
	}
}
// TestEndPhaseWithoutStart teste EndPhase appelé sans StartPhase
func TestEndPhaseWithoutStart(t *testing.T) {
	collector := NewCoherenceMetricsCollector()
	// Appeler EndPhase sans StartPhase ne devrait pas paniquer
	collector.EndPhase("orphan_phase", 10, true)
	metrics := collector.GetMetrics()
	phase := metrics.PhaseMetrics["orphan_phase"]
	if phase == nil {
		t.Fatal("La phase orpheline devrait quand même être enregistrée")
	}
	if phase.ItemsProcessed != 10 {
		t.Errorf("ItemsProcessed: attendu 10, obtenu %d", phase.ItemsProcessed)
	}
	if !phase.Succeeded {
		t.Error("La phase devrait être marquée comme réussie")
	}
}