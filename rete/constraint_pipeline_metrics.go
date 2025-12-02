// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sync"
	"time"
)

// IngestionMetrics contient les métriques de performance pour l'ingestion incrémentale
type IngestionMetrics struct {
	// Durées
	ParsingDuration        time.Duration `json:"parsing_duration"`
	ValidationDuration     time.Duration `json:"validation_duration"`
	TypeCreationDuration   time.Duration `json:"type_creation_duration"`
	RuleCreationDuration   time.Duration `json:"rule_creation_duration"`
	FactCollectionDuration time.Duration `json:"fact_collection_duration"`
	PropagationDuration    time.Duration `json:"propagation_duration"`
	FactSubmissionDuration time.Duration `json:"fact_submission_duration"`
	TotalDuration          time.Duration `json:"total_duration"`

	// Compteurs
	TypesAdded             int `json:"types_added"`
	RulesAdded             int `json:"rules_added"`
	FactsSubmitted         int `json:"facts_submitted"`
	ExistingFactsCollected int `json:"existing_facts_collected"`
	FactsPropagated        int `json:"facts_propagated"`
	NewTerminalsAdded      int `json:"new_terminals_added"`
	PropagationTargets     int `json:"propagation_targets"` // Nombre de terminaux ciblés

	// États
	WasReset          bool      `json:"was_reset"`
	WasIncremental    bool      `json:"was_incremental"`
	ValidationSkipped bool      `json:"validation_skipped"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`

	// Réseau
	TotalTypeNodes     int `json:"total_type_nodes"`
	TotalTerminalNodes int `json:"total_terminal_nodes"`
	TotalAlphaNodes    int `json:"total_alpha_nodes"`
	TotalBetaNodes     int `json:"total_beta_nodes"`
}

// MetricsCollector collecte les métriques pendant l'ingestion
type MetricsCollector struct {
	metrics *IngestionMetrics
	mutex   sync.RWMutex
}

// NewMetricsCollector crée un nouveau collecteur de métriques
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: &IngestionMetrics{
			StartTime: time.Now(),
		},
	}
}

// RecordParsingDuration enregistre la durée du parsing
func (mc *MetricsCollector) RecordParsingDuration(duration time.Duration) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.ParsingDuration = duration
}

// RecordValidationDuration enregistre la durée de la validation
func (mc *MetricsCollector) RecordValidationDuration(duration time.Duration) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.ValidationDuration = duration
}

// RecordTypeCreationDuration enregistre la durée de création des types
func (mc *MetricsCollector) RecordTypeCreationDuration(duration time.Duration) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.TypeCreationDuration = duration
}

// RecordRuleCreationDuration enregistre la durée de création des règles
func (mc *MetricsCollector) RecordRuleCreationDuration(duration time.Duration) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.RuleCreationDuration = duration
}

// RecordFactCollectionDuration enregistre la durée de collection des faits
func (mc *MetricsCollector) RecordFactCollectionDuration(duration time.Duration) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.FactCollectionDuration = duration
}

// RecordPropagationDuration enregistre la durée de propagation
func (mc *MetricsCollector) RecordPropagationDuration(duration time.Duration) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.PropagationDuration = duration
}

// RecordFactSubmissionDuration enregistre la durée de soumission des faits
func (mc *MetricsCollector) RecordFactSubmissionDuration(duration time.Duration) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.FactSubmissionDuration = duration
}

// SetTypesAdded enregistre le nombre de types ajoutés
func (mc *MetricsCollector) SetTypesAdded(count int) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.TypesAdded = count
}

// SetRulesAdded enregistre le nombre de règles ajoutées
func (mc *MetricsCollector) SetRulesAdded(count int) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.RulesAdded = count
}

// SetFactsSubmitted enregistre le nombre de faits soumis
func (mc *MetricsCollector) SetFactsSubmitted(count int) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.FactsSubmitted = count
}

// SetExistingFactsCollected enregistre le nombre de faits existants collectés
func (mc *MetricsCollector) SetExistingFactsCollected(count int) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.ExistingFactsCollected = count
}

// SetFactsPropagated enregistre le nombre de faits propagés
func (mc *MetricsCollector) SetFactsPropagated(count int) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.FactsPropagated = count
}

// SetNewTerminalsAdded enregistre le nombre de nouveaux terminaux ajoutés
func (mc *MetricsCollector) SetNewTerminalsAdded(count int) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.NewTerminalsAdded = count
}

// SetPropagationTargets enregistre le nombre de cibles de propagation
func (mc *MetricsCollector) SetPropagationTargets(count int) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.PropagationTargets = count
}

// SetWasReset marque que le réseau a été réinitialisé
func (mc *MetricsCollector) SetWasReset(wasReset bool) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.WasReset = wasReset
}

// SetWasIncremental marque que c'était une ingestion incrémentale
func (mc *MetricsCollector) SetWasIncremental(wasIncremental bool) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.WasIncremental = wasIncremental
}

// SetValidationSkipped marque que la validation a été ignorée
func (mc *MetricsCollector) SetValidationSkipped(skipped bool) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.metrics.ValidationSkipped = skipped
}

// RecordNetworkState enregistre l'état final du réseau
func (mc *MetricsCollector) RecordNetworkState(network *ReteNetwork) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.metrics.TotalTypeNodes = len(network.TypeNodes)
	mc.metrics.TotalTerminalNodes = len(network.TerminalNodes)
	mc.metrics.TotalAlphaNodes = len(network.AlphaNodes)
	mc.metrics.TotalBetaNodes = len(network.BetaNodes)
}

// Finalize finalise les métriques en calculant la durée totale
func (mc *MetricsCollector) Finalize() *IngestionMetrics {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.metrics.EndTime = time.Now()
	mc.metrics.TotalDuration = mc.metrics.EndTime.Sub(mc.metrics.StartTime)

	return mc.metrics
}

// GetMetrics retourne une copie des métriques actuelles
func (mc *MetricsCollector) GetMetrics() *IngestionMetrics {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	// Créer une copie pour éviter les modifications concurrentes
	metricsCopy := *mc.metrics
	return &metricsCopy
}

// String retourne une représentation formatée des métriques
func (m *IngestionMetrics) String() string {
	return fmt.Sprintf(`
📊 Métriques d'Ingestion RETE
════════════════════════════════════════
⏱️  Durées:
   - Parsing:              %v
   - Validation:           %v
   - Création types:       %v
   - Création règles:      %v
   - Collection faits:     %v
   - Propagation:          %v
   - Soumission faits:     %v
   - TOTAL:                %v

📈 Compteurs:
   - Types ajoutés:        %d
   - Règles ajoutées:      %d
   - Faits soumis:         %d
   - Faits existants:      %d
   - Faits propagés:       %d
   - Nouveaux terminaux:   %d
   - Cibles propagation:   %d

🏗️  État du réseau:
   - TypeNodes:            %d
   - TerminalNodes:        %d
   - AlphaNodes:           %d
   - BetaNodes:            %d

🔄 Mode:
   - Reset:                %v
   - Incrémental:          %v
   - Validation ignorée:   %v
════════════════════════════════════════`,
		m.ParsingDuration,
		m.ValidationDuration,
		m.TypeCreationDuration,
		m.RuleCreationDuration,
		m.FactCollectionDuration,
		m.PropagationDuration,
		m.FactSubmissionDuration,
		m.TotalDuration,
		m.TypesAdded,
		m.RulesAdded,
		m.FactsSubmitted,
		m.ExistingFactsCollected,
		m.FactsPropagated,
		m.NewTerminalsAdded,
		m.PropagationTargets,
		m.TotalTypeNodes,
		m.TotalTerminalNodes,
		m.TotalAlphaNodes,
		m.TotalBetaNodes,
		m.WasReset,
		m.WasIncremental,
		m.ValidationSkipped,
	)
}

// Summary retourne un résumé court des métriques
func (m *IngestionMetrics) Summary() string {
	return fmt.Sprintf(
		"Ingestion: %v total | %d types, %d règles, %d faits | %d propagés vers %d nouveaux terminaux",
		m.TotalDuration,
		m.TypesAdded,
		m.RulesAdded,
		m.FactsSubmitted,
		m.FactsPropagated,
		m.NewTerminalsAdded,
	)
}

// IsEfficient retourne true si l'ingestion a été efficace
func (m *IngestionMetrics) IsEfficient() bool {
	// Considéré efficace si:
	// 1. Propagation ciblée (moins ou égal de faits propagés que d'existants × nouveaux terminaux)
	// 2. Durée de propagation raisonnable (< 30% du temps total si propagation nécessaire)
	//    Note: Pour de petites ingestions, le ratio peut être plus élevé à cause de l'overhead fixe

	if m.NewTerminalsAdded > 0 && m.ExistingFactsCollected > 0 {
		maxPropagations := m.ExistingFactsCollected * m.NewTerminalsAdded
		if m.FactsPropagated > maxPropagations {
			return false // Trop de propagations (pas de ciblage)
		}

		// Si propagation ciblée est correcte, c'est déjà efficace
		// Le ratio de temps est moins important que la propagation ciblée
		if m.TotalDuration > 0 && m.TotalDuration > 1*time.Millisecond {
			propagationRatio := float64(m.PropagationDuration) / float64(m.TotalDuration)
			// Pour de grandes ingestions (> 1ms), vérifier le ratio
			if propagationRatio > 0.3 {
				return false // Propagation prend trop de temps
			}
		}
	}

	return true
}

// GetBottleneck identifie le goulot d'étranglement principal
func (m *IngestionMetrics) GetBottleneck() string {
	durations := map[string]time.Duration{
		"Parsing":          m.ParsingDuration,
		"Validation":       m.ValidationDuration,
		"Création types":   m.TypeCreationDuration,
		"Création règles":  m.RuleCreationDuration,
		"Collection faits": m.FactCollectionDuration,
		"Propagation":      m.PropagationDuration,
		"Soumission faits": m.FactSubmissionDuration,
	}

	var maxDuration time.Duration
	var bottleneck string

	for name, duration := range durations {
		if duration > maxDuration {
			maxDuration = duration
			bottleneck = name
		}
	}

	if bottleneck == "" {
		return "Aucun goulot d'étranglement identifié"
	}

	percentage := float64(maxDuration) / float64(m.TotalDuration) * 100
	return fmt.Sprintf("%s (%.1f%% du temps total)", bottleneck, percentage)
}
