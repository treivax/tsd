// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// CoherenceMetrics contient les métriques détaillées pour la cohérence et la synchronisation
type CoherenceMetrics struct {
	// Métriques de soumission des faits
	FactsSubmitted      int `json:"facts_submitted"`
	FactsPersisted      int `json:"facts_persisted"`
	FactsRetried        int `json:"facts_retried"`
	FactsFailed         int `json:"facts_failed"`
	FactsPropagated     int `json:"facts_propagated"`
	TerminalActivations int `json:"terminal_activations"`

	// Métriques de synchronisation
	TotalVerifyAttempts     int `json:"total_verify_attempts"`
	TotalRetries            int `json:"total_retries"`
	TotalTimeouts           int `json:"total_timeouts"`
	MaxRetriesForSingleFact int `json:"max_retries_for_single_fact"`

	// Métriques de temps
	TotalWaitTime time.Duration `json:"total_wait_time"`
	MaxWaitTime   time.Duration `json:"max_wait_time"`
	MinWaitTime   time.Duration `json:"min_wait_time"`
	AvgWaitTime   time.Duration `json:"avg_wait_time"`

	TotalSyncTime       time.Duration `json:"total_sync_time"`
	TotalSubmissionTime time.Duration `json:"total_submission_time"`

	// Métriques de cohérence pré-commit
	PreCommitChecks    int `json:"pre_commit_checks"`
	PreCommitSuccesses int `json:"pre_commit_successes"`
	PreCommitFailures  int `json:"pre_commit_failures"`

	// Détails par phase
	PhaseMetrics map[string]*PhaseMetrics `json:"phase_metrics"`

	// Timestamps
	StartTime     time.Time     `json:"start_time"`
	EndTime       time.Time     `json:"end_time"`
	TotalDuration time.Duration `json:"total_duration"`

	// État
	TransactionID  string `json:"transaction_id,omitempty"`
	WasRolledBack  bool   `json:"was_rolled_back"`
	RollbackReason string `json:"rollback_reason,omitempty"`
}

// PhaseMetrics contient les métriques pour une phase spécifique d'ingestion
type PhaseMetrics struct {
	PhaseName      string        `json:"phase_name"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
	Duration       time.Duration `json:"duration"`
	ItemsProcessed int           `json:"items_processed"`
	Errors         int           `json:"errors"`
	Succeeded      bool          `json:"succeeded"`
}

// CoherenceMetricsCollector collecte les métriques de cohérence pendant l'ingestion
type CoherenceMetricsCollector struct {
	metrics        *CoherenceMetrics
	mutex          sync.RWMutex
	activePhases   map[string]*PhaseMetrics
	minWaitTimeSet bool
}

// NewCoherenceMetricsCollector crée un nouveau collecteur de métriques de cohérence
func NewCoherenceMetricsCollector() *CoherenceMetricsCollector {
	return &CoherenceMetricsCollector{
		metrics: &CoherenceMetrics{
			StartTime:    time.Now(),
			PhaseMetrics: make(map[string]*PhaseMetrics),
			MinWaitTime:  time.Duration(0),
		},
		activePhases:   make(map[string]*PhaseMetrics),
		minWaitTimeSet: false,
	}
}

// StartPhase démarre le chronomètre pour une phase
func (cmc *CoherenceMetricsCollector) StartPhase(phaseName string) {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()

	phase := &PhaseMetrics{
		PhaseName: phaseName,
		StartTime: time.Now(),
	}
	cmc.activePhases[phaseName] = phase
}

// EndPhase termine le chronomètre pour une phase
func (cmc *CoherenceMetricsCollector) EndPhase(phaseName string, itemsProcessed int, succeeded bool) {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()

	phase, exists := cmc.activePhases[phaseName]
	if !exists {
		// Phase non démarrée, créer une entrée par défaut
		phase = &PhaseMetrics{
			PhaseName: phaseName,
			StartTime: time.Now(),
		}
	}

	phase.EndTime = time.Now()
	phase.Duration = phase.EndTime.Sub(phase.StartTime)
	phase.ItemsProcessed = itemsProcessed
	phase.Succeeded = succeeded

	cmc.metrics.PhaseMetrics[phaseName] = phase
	delete(cmc.activePhases, phaseName)
}

// RecordFactSubmitted enregistre un fait soumis
func (cmc *CoherenceMetricsCollector) RecordFactSubmitted() {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()
	cmc.metrics.FactsSubmitted++
}

// RecordFactPersisted enregistre un fait persisté
func (cmc *CoherenceMetricsCollector) RecordFactPersisted() {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()
	cmc.metrics.FactsPersisted++
}

// RecordFactRetried enregistre un fait qui a nécessité un retry
func (cmc *CoherenceMetricsCollector) RecordFactRetried() {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()
	cmc.metrics.FactsRetried++
}

// RecordFactFailed enregistre un fait qui a échoué
func (cmc *CoherenceMetricsCollector) RecordFactFailed() {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()
	cmc.metrics.FactsFailed++
}

// RecordFactPropagated enregistre un fait propagé
func (cmc *CoherenceMetricsCollector) RecordFactPropagated() {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()
	cmc.metrics.FactsPropagated++
}

// RecordTerminalActivation enregistre une activation de terminal
func (cmc *CoherenceMetricsCollector) RecordTerminalActivation() {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()
	cmc.metrics.TerminalActivations++
}

// RecordVerifyAttempt enregistre une tentative de vérification
func (cmc *CoherenceMetricsCollector) RecordVerifyAttempt() {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()
	cmc.metrics.TotalVerifyAttempts++
}

// RecordRetry enregistre un retry avec le nombre de tentatives pour ce fait
func (cmc *CoherenceMetricsCollector) RecordRetry(attemptCount int) {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()
	cmc.metrics.TotalRetries++
	if attemptCount > cmc.metrics.MaxRetriesForSingleFact {
		cmc.metrics.MaxRetriesForSingleFact = attemptCount
	}
}

// RecordTimeout enregistre un timeout
func (cmc *CoherenceMetricsCollector) RecordTimeout() {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()
	cmc.metrics.TotalTimeouts++
}

// RecordWaitTime enregistre le temps d'attente pour la persistance d'un fait
func (cmc *CoherenceMetricsCollector) RecordWaitTime(waitTime time.Duration) {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()

	cmc.metrics.TotalWaitTime += waitTime

	if waitTime > cmc.metrics.MaxWaitTime {
		cmc.metrics.MaxWaitTime = waitTime
	}

	if !cmc.minWaitTimeSet || waitTime < cmc.metrics.MinWaitTime {
		cmc.metrics.MinWaitTime = waitTime
		cmc.minWaitTimeSet = true
	}
}

// RecordSyncTime enregistre le temps de synchronisation du storage
func (cmc *CoherenceMetricsCollector) RecordSyncTime(syncTime time.Duration) {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()
	cmc.metrics.TotalSyncTime += syncTime
}

// RecordSubmissionTime enregistre le temps total de soumission
func (cmc *CoherenceMetricsCollector) RecordSubmissionTime(submissionTime time.Duration) {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()
	cmc.metrics.TotalSubmissionTime = submissionTime
}

// RecordPreCommitCheck enregistre une vérification pré-commit
func (cmc *CoherenceMetricsCollector) RecordPreCommitCheck(success bool) {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()
	cmc.metrics.PreCommitChecks++
	if success {
		cmc.metrics.PreCommitSuccesses++
	} else {
		cmc.metrics.PreCommitFailures++
	}
}

// SetTransactionID définit l'ID de transaction
func (cmc *CoherenceMetricsCollector) SetTransactionID(txID string) {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()
	cmc.metrics.TransactionID = txID
}

// RecordRollback enregistre un rollback avec sa raison
func (cmc *CoherenceMetricsCollector) RecordRollback(reason string) {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()
	cmc.metrics.WasRolledBack = true
	cmc.metrics.RollbackReason = reason
}

// Finalize finalise les métriques en calculant les statistiques finales
func (cmc *CoherenceMetricsCollector) Finalize() *CoherenceMetrics {
	cmc.mutex.Lock()
	defer cmc.mutex.Unlock()

	cmc.metrics.EndTime = time.Now()
	cmc.metrics.TotalDuration = cmc.metrics.EndTime.Sub(cmc.metrics.StartTime)

	// Calculer le temps d'attente moyen
	if cmc.metrics.FactsPersisted > 0 {
		cmc.metrics.AvgWaitTime = cmc.metrics.TotalWaitTime / time.Duration(cmc.metrics.FactsPersisted)
	}

	return cmc.metrics
}

// GetMetrics retourne une copie des métriques actuelles
func (cmc *CoherenceMetricsCollector) GetMetrics() *CoherenceMetrics {
	cmc.mutex.RLock()
	defer cmc.mutex.RUnlock()

	// Créer une copie profonde
	metricsCopy := *cmc.metrics
	metricsCopy.PhaseMetrics = make(map[string]*PhaseMetrics)
	for k, v := range cmc.metrics.PhaseMetrics {
		phaseCopy := *v
		metricsCopy.PhaseMetrics[k] = &phaseCopy
	}

	return &metricsCopy
}

// ToJSON exporte les métriques en JSON
func (cm *CoherenceMetrics) ToJSON() (string, error) {
	data, err := json.MarshalIndent(cm, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erreur lors de l'export JSON: %w", err)
	}
	return string(data), nil
}

// String retourne une représentation formatée des métriques de cohérence
func (cm *CoherenceMetrics) String() string {
	return fmt.Sprintf(`
📊 Métriques de Cohérence RETE
════════════════════════════════════════
📦 Faits:
   - Soumis:               %d
   - Persistés:            %d
   - Avec retry:           %d
   - Échoués:              %d
   - Propagés:             %d
   - Activations term.:    %d

🔄 Synchronisation:
   - Tentatives vérif.:    %d
   - Total retries:        %d
   - Max retries (1 fait): %d
   - Timeouts:             %d

⏱️  Temps d'attente:
   - Total:                %v
   - Moyen:                %v
   - Max:                  %v
   - Min:                  %v

⏰ Temps système:
   - Sync storage:         %v
   - Soumission totale:    %v
   - Durée totale:         %v

✅ Cohérence pré-commit:
   - Vérifications:        %d
   - Succès:               %d
   - Échecs:               %d

🔄 Transaction:
   - ID:                   %s
   - Rollback:             %v
   - Raison:               %s

📋 Phases:
%s
════════════════════════════════════════`,
		cm.FactsSubmitted,
		cm.FactsPersisted,
		cm.FactsRetried,
		cm.FactsFailed,
		cm.FactsPropagated,
		cm.TerminalActivations,
		cm.TotalVerifyAttempts,
		cm.TotalRetries,
		cm.MaxRetriesForSingleFact,
		cm.TotalTimeouts,
		cm.TotalWaitTime,
		cm.AvgWaitTime,
		cm.MaxWaitTime,
		cm.MinWaitTime,
		cm.TotalSyncTime,
		cm.TotalSubmissionTime,
		cm.TotalDuration,
		cm.PreCommitChecks,
		cm.PreCommitSuccesses,
		cm.PreCommitFailures,
		cm.TransactionID,
		cm.WasRolledBack,
		cm.RollbackReason,
		cm.formatPhases(),
	)
}

// formatPhases formate les métriques des phases
func (cm *CoherenceMetrics) formatPhases() string {
	if len(cm.PhaseMetrics) == 0 {
		return "   (aucune phase enregistrée)"
	}

	result := ""
	for _, phase := range cm.PhaseMetrics {
		status := "✅"
		if !phase.Succeeded {
			status = "❌"
		}
		result += fmt.Sprintf("   %s %s: %v (%d items, %d erreurs)\n",
			status,
			phase.PhaseName,
			phase.Duration,
			phase.ItemsProcessed,
			phase.Errors,
		)
	}
	return result
}

// Summary retourne un résumé court des métriques de cohérence
func (cm *CoherenceMetrics) Summary() string {
	successRate := 0.0
	if cm.FactsSubmitted > 0 {
		successRate = float64(cm.FactsPersisted) / float64(cm.FactsSubmitted) * 100
	}

	return fmt.Sprintf(
		"Cohérence: %d/%d faits persistés (%.1f%%) | %d retries | %d timeouts | wait moyen: %v",
		cm.FactsPersisted,
		cm.FactsSubmitted,
		successRate,
		cm.TotalRetries,
		cm.TotalTimeouts,
		cm.AvgWaitTime,
	)
}

// IsHealthy retourne true si les métriques indiquent un système en bonne santé
func (cm *CoherenceMetrics) IsHealthy() bool {
	// Critères de santé:
	// 1. Taux de succès >= 95%
	// 2. Pas trop de timeouts (< 5% des faits)
	// 3. Nombre moyen de retries acceptable (< 2 par fait persisté)
	// 4. Pas de rollback (sauf si explicitement voulu)

	if cm.FactsSubmitted == 0 {
		return true // Pas de faits = pas de problème
	}

	successRate := float64(cm.FactsPersisted) / float64(cm.FactsSubmitted)
	if successRate < 0.95 {
		return false
	}

	timeoutRate := float64(cm.TotalTimeouts) / float64(cm.FactsSubmitted)
	if timeoutRate >= 0.05 {
		return false
	}

	if cm.FactsPersisted > 0 {
		avgRetries := float64(cm.TotalRetries) / float64(cm.FactsPersisted)
		if avgRetries >= 2.0 {
			return false
		}
	}

	// Un rollback n'est pas forcément un problème (peut être intentionnel)
	// mais on le signale quand même
	if cm.WasRolledBack && cm.RollbackReason != "" {
		return false
	}

	return true
}

// GetHealthReport génère un rapport de santé détaillé
func (cm *CoherenceMetrics) GetHealthReport() string {
	if cm.IsHealthy() {
		return "✅ Système en bonne santé"
	}

	issues := []string{}

	if cm.FactsSubmitted > 0 {
		successRate := float64(cm.FactsPersisted) / float64(cm.FactsSubmitted)
		if successRate < 0.95 {
			issues = append(issues, fmt.Sprintf("❌ Taux de succès bas: %.1f%%", successRate*100))
		}

		timeoutRate := float64(cm.TotalTimeouts) / float64(cm.FactsSubmitted)
		if timeoutRate >= 0.05 {
			issues = append(issues, fmt.Sprintf("⚠️  Trop de timeouts: %.1f%%", timeoutRate*100))
		}

		if cm.FactsPersisted > 0 {
			avgRetries := float64(cm.TotalRetries) / float64(cm.FactsPersisted)
			if avgRetries >= 2.0 {
				issues = append(issues, fmt.Sprintf("⚠️  Trop de retries: %.2f en moyenne", avgRetries))
			}
		}
	}

	if cm.WasRolledBack && cm.RollbackReason != "" {
		issues = append(issues, fmt.Sprintf("🔙 Rollback: %s", cm.RollbackReason))
	}

	if len(issues) == 0 {
		return "✅ Système en bonne santé"
	}

	result := "⚠️  Problèmes détectés:\n"
	for _, issue := range issues {
		result += "   " + issue + "\n"
	}

	return result
}
