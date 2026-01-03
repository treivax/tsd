// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"sync"
	"time"
)

// PropagationMetrics collecte des statistiques sur les propagations delta.
//
// Ces métriques permettent de monitorer la performance et l'efficacité
// du système de propagation sélective.
type PropagationMetrics struct {
	TotalPropagations       int64
	DeltaPropagations       int64
	ClassicPropagations     int64
	FailedPropagations      int64
	TotalPropagationTime    time.Duration
	AvgPropagationTime      time.Duration
	MinPropagationTime      time.Duration
	MaxPropagationTime      time.Duration
	TotalNodesEvaluated     int64
	NodesSkippedByDelta     int64
	AvgNodesPerPropagation  float64
	TotalFieldsChanged      int64
	AvgFieldsPerPropagation float64
	FallbacksDueToFields    int64
	FallbacksDueToRatio     int64
	FallbacksDueToNodes     int64
	FallbacksDueToPK        int64
	FallbacksDueToError     int64
	FirstPropagation        time.Time
	LastPropagation         time.Time
	mutex                   sync.RWMutex
}

// NewPropagationMetrics crée une nouvelle instance de métriques.
func NewPropagationMetrics() *PropagationMetrics {
	return &PropagationMetrics{
		MinPropagationTime: time.Duration(1<<63 - 1),
	}
}

// RecordDeltaPropagation enregistre une propagation delta.
//
// Paramètres :
//   - duration : temps pris par la propagation
//   - nodesAffected : nombre de nœuds affectés
//   - fieldsChanged : nombre de champs modifiés
func (pm *PropagationMetrics) RecordDeltaPropagation(
	duration time.Duration,
	nodesAffected int,
	fieldsChanged int,
) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.TotalPropagations++
	pm.DeltaPropagations++
	pm.TotalNodesEvaluated += int64(nodesAffected)
	pm.TotalFieldsChanged += int64(fieldsChanged)

	pm.updateTiming(duration)
	pm.updateTimestamps()
	pm.recalculateAverages()
}

// RecordClassicPropagation enregistre une propagation classique (Retract+Insert).
//
// Paramètres :
//   - duration : temps pris
//   - totalNodes : nombre total de nœuds dans le réseau
func (pm *PropagationMetrics) RecordClassicPropagation(
	duration time.Duration,
	totalNodes int,
) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.TotalPropagations++
	pm.ClassicPropagations++
	pm.TotalNodesEvaluated += int64(totalNodes)

	pm.updateTiming(duration)
	pm.updateTimestamps()
	pm.recalculateAverages()
}

// RecordFailedPropagation enregistre une propagation échouée.
func (pm *PropagationMetrics) RecordFailedPropagation() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.TotalPropagations++
	pm.FailedPropagations++
	pm.updateTimestamps()
}

// RecordFallback enregistre un fallback vers mode classique.
//
// Paramètres :
//   - reason : raison du fallback ("ratio", "nodes", "pk", "error")
func (pm *PropagationMetrics) RecordFallback(reason string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	switch reason {
	case "fields":
		pm.FallbacksDueToFields++
	case "ratio":
		pm.FallbacksDueToRatio++
	case "nodes":
		pm.FallbacksDueToNodes++
	case "pk":
		pm.FallbacksDueToPK++
	case "error":
		pm.FallbacksDueToError++
	}
}

// RecordNodesSkipped enregistre des nœuds évités grâce au delta.
//
// Paramètres :
//   - count : nombre de nœuds évités
func (pm *PropagationMetrics) RecordNodesSkipped(count int) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.NodesSkippedByDelta += int64(count)
}

// GetSnapshot retourne un instantané des métriques actuelles.
//
// Retourne une copie thread-safe des métriques.
func (pm *PropagationMetrics) GetSnapshot() PropagationMetrics {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	return PropagationMetrics{
		TotalPropagations:       pm.TotalPropagations,
		DeltaPropagations:       pm.DeltaPropagations,
		ClassicPropagations:     pm.ClassicPropagations,
		FailedPropagations:      pm.FailedPropagations,
		TotalPropagationTime:    pm.TotalPropagationTime,
		AvgPropagationTime:      pm.AvgPropagationTime,
		MinPropagationTime:      pm.MinPropagationTime,
		MaxPropagationTime:      pm.MaxPropagationTime,
		TotalNodesEvaluated:     pm.TotalNodesEvaluated,
		NodesSkippedByDelta:     pm.NodesSkippedByDelta,
		AvgNodesPerPropagation:  pm.AvgNodesPerPropagation,
		TotalFieldsChanged:      pm.TotalFieldsChanged,
		AvgFieldsPerPropagation: pm.AvgFieldsPerPropagation,
		FallbacksDueToFields:    pm.FallbacksDueToFields,
		FallbacksDueToRatio:     pm.FallbacksDueToRatio,
		FallbacksDueToNodes:     pm.FallbacksDueToNodes,
		FallbacksDueToPK:        pm.FallbacksDueToPK,
		FallbacksDueToError:     pm.FallbacksDueToError,
		FirstPropagation:        pm.FirstPropagation,
		LastPropagation:         pm.LastPropagation,
	}
}

// GetEfficiencyRatio retourne le ratio d'efficacité de la propagation delta.
//
// Ratio = NodesSkipped / TotalNodesEvaluated
// Plus ce ratio est élevé, plus le système delta est efficace.
//
// Retourne une valeur entre 0.0 et 1.0, ou 0.0 si aucune propagation.
func (pm *PropagationMetrics) GetEfficiencyRatio() float64 {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	if pm.TotalNodesEvaluated == 0 {
		return 0.0
	}

	totalClassicNodes := pm.TotalNodesEvaluated + pm.NodesSkippedByDelta

	if totalClassicNodes == 0 {
		return 0.0
	}

	return float64(pm.NodesSkippedByDelta) / float64(totalClassicNodes)
}

// GetDeltaUsageRatio retourne le ratio d'utilisation du mode delta.
//
// Ratio = DeltaPropagations / TotalPropagations
//
// Retourne une valeur entre 0.0 et 1.0.
func (pm *PropagationMetrics) GetDeltaUsageRatio() float64 {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	if pm.TotalPropagations == 0 {
		return 0.0
	}

	return float64(pm.DeltaPropagations) / float64(pm.TotalPropagations)
}

// Reset réinitialise toutes les métriques à zéro.
func (pm *PropagationMetrics) Reset() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.TotalPropagations = 0
	pm.DeltaPropagations = 0
	pm.ClassicPropagations = 0
	pm.FailedPropagations = 0
	pm.TotalPropagationTime = 0
	pm.AvgPropagationTime = 0
	pm.MinPropagationTime = time.Duration(1<<63 - 1)
	pm.MaxPropagationTime = 0
	pm.TotalNodesEvaluated = 0
	pm.NodesSkippedByDelta = 0
	pm.AvgNodesPerPropagation = 0
	pm.TotalFieldsChanged = 0
	pm.AvgFieldsPerPropagation = 0
	pm.FallbacksDueToRatio = 0
	pm.FallbacksDueToNodes = 0
	pm.FallbacksDueToPK = 0
	pm.FallbacksDueToError = 0
	pm.FirstPropagation = time.Time{}
	pm.LastPropagation = time.Time{}
}

// updateTiming met à jour les statistiques de timing.
// ATTENTION : doit être appelé avec mutex déjà acquis.
func (pm *PropagationMetrics) updateTiming(duration time.Duration) {
	pm.TotalPropagationTime += duration

	if duration < pm.MinPropagationTime {
		pm.MinPropagationTime = duration
	}

	if duration > pm.MaxPropagationTime {
		pm.MaxPropagationTime = duration
	}
}

// updateTimestamps met à jour les timestamps.
// ATTENTION : doit être appelé avec mutex déjà acquis.
func (pm *PropagationMetrics) updateTimestamps() {
	now := time.Now()

	if pm.FirstPropagation.IsZero() {
		pm.FirstPropagation = now
	}

	pm.LastPropagation = now
}

// recalculateAverages recalcule les moyennes.
// ATTENTION : doit être appelé avec mutex déjà acquis.
func (pm *PropagationMetrics) recalculateAverages() {
	if pm.TotalPropagations > 0 {
		pm.AvgPropagationTime = time.Duration(
			int64(pm.TotalPropagationTime) / pm.TotalPropagations,
		)

		pm.AvgNodesPerPropagation = float64(pm.TotalNodesEvaluated) /
			float64(pm.TotalPropagations)

		pm.AvgFieldsPerPropagation = float64(pm.TotalFieldsChanged) /
			float64(pm.TotalPropagations)
	}
}
