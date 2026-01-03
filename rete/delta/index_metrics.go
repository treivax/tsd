// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
	"time"
)

// IndexMetrics collecte et expose des métriques sur l'utilisation du DependencyIndex.
//
// Cette structure permet de suivre les performances et l'efficacité de l'indexation.
type IndexMetrics struct {
	// Compteurs d'opérations
	lookupCount  int64
	nodeAddCount int64
	clearCount   int64

	// Métriques de performance
	totalLookupTime   time.Duration
	avgNodesPerLookup float64

	// Dernière mise à jour
	lastUpdate time.Time
}

// NewIndexMetrics crée une nouvelle instance de métriques.
func NewIndexMetrics() *IndexMetrics {
	return &IndexMetrics{
		lastUpdate: time.Now(),
	}
}

// RecordLookup enregistre une opération de recherche.
//
// Paramètres :
//   - duration : temps pris pour la recherche
//   - nodesFound : nombre de nœuds trouvés
func (im *IndexMetrics) RecordLookup(duration time.Duration, nodesFound int) {
	im.lookupCount++
	im.totalLookupTime += duration

	// Calcul de la moyenne mobile
	alpha := 0.1
	im.avgNodesPerLookup = alpha*float64(nodesFound) + (1-alpha)*im.avgNodesPerLookup

	im.lastUpdate = time.Now()
}

// RecordNodeAdd enregistre l'ajout d'un nœud à l'index.
func (im *IndexMetrics) RecordNodeAdd() {
	im.nodeAddCount++
	im.lastUpdate = time.Now()
}

// RecordClear enregistre une opération de vidage de l'index.
func (im *IndexMetrics) RecordClear() {
	im.clearCount++
	im.lastUpdate = time.Now()
}

// GetAverageLookupTime retourne le temps moyen de recherche.
func (im *IndexMetrics) GetAverageLookupTime() time.Duration {
	if im.lookupCount == 0 {
		return 0
	}
	return time.Duration(int64(im.totalLookupTime) / im.lookupCount)
}

// GetLookupCount retourne le nombre total de recherches effectuées.
func (im *IndexMetrics) GetLookupCount() int64 {
	return im.lookupCount
}

// GetNodeAddCount retourne le nombre total d'ajouts de nœuds.
func (im *IndexMetrics) GetNodeAddCount() int64 {
	return im.nodeAddCount
}

// GetAverageNodesPerLookup retourne le nombre moyen de nœuds trouvés par recherche.
func (im *IndexMetrics) GetAverageNodesPerLookup() float64 {
	return im.avgNodesPerLookup
}

// Reset réinitialise toutes les métriques.
func (im *IndexMetrics) Reset() {
	im.lookupCount = 0
	im.nodeAddCount = 0
	im.clearCount = 0
	im.totalLookupTime = 0
	im.avgNodesPerLookup = 0
	im.lastUpdate = time.Now()
}

// String retourne une représentation string des métriques.
func (im *IndexMetrics) String() string {
	return fmt.Sprintf(
		"IndexMetrics[lookups=%d, adds=%d, avgTime=%v, avgNodes=%.2f]",
		im.lookupCount,
		im.nodeAddCount,
		im.GetAverageLookupTime(),
		im.avgNodesPerLookup,
	)
}

// MetricsSnapshot représente un instantané des métriques à un moment donné.
type MetricsSnapshot struct {
	LookupCount           int64
	NodeAddCount          int64
	ClearCount            int64
	AverageLookupTime     time.Duration
	AverageNodesPerLookup float64
	Timestamp             time.Time
}

// Snapshot crée un instantané des métriques actuelles.
func (im *IndexMetrics) Snapshot() MetricsSnapshot {
	return MetricsSnapshot{
		LookupCount:           im.lookupCount,
		NodeAddCount:          im.nodeAddCount,
		ClearCount:            im.clearCount,
		AverageLookupTime:     im.GetAverageLookupTime(),
		AverageNodesPerLookup: im.avgNodesPerLookup,
		Timestamp:             time.Now(),
	}
}

// String retourne une représentation string du snapshot.
func (ms MetricsSnapshot) String() string {
	return fmt.Sprintf(
		"Snapshot[lookups=%d, adds=%d, clears=%d, avgTime=%v, avgNodes=%.2f, at=%s]",
		ms.LookupCount,
		ms.NodeAddCount,
		ms.ClearCount,
		ms.AverageLookupTime,
		ms.AverageNodesPerLookup,
		ms.Timestamp.Format(time.RFC3339),
	)
}
