// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"sync"
	"time"
)

// ChainBuildMetricsData contient les données de métriques sans le mutex
type ChainBuildMetricsData struct {
	TotalChainsBuilt   int     `json:"total_chains_built"`
	TotalNodesCreated  int     `json:"total_nodes_created"`
	TotalNodesReused   int     `json:"total_nodes_reused"`
	AverageChainLength float64 `json:"average_chain_length"`
	SharingRatio       float64 `json:"sharing_ratio"`

	// Métriques de cache
	HashCacheHits   int `json:"hash_cache_hits"`
	HashCacheMisses int `json:"hash_cache_misses"`
	HashCacheSize   int `json:"hash_cache_size"`

	// Métriques de connexion
	ConnectionCacheHits   int `json:"connection_cache_hits"`
	ConnectionCacheMisses int `json:"connection_cache_misses"`

	// Métriques de temps
	TotalBuildTime       time.Duration `json:"total_build_time_ns"`
	AverageBuildTime     time.Duration `json:"average_build_time_ns"`
	TotalHashComputeTime time.Duration `json:"total_hash_compute_time_ns"`

	// Détails par règle
	ChainDetails []ChainMetricDetail `json:"chain_details,omitempty"`
}

// ChainBuildMetrics collecte les métriques de performance pour la construction des chaînes alpha
type ChainBuildMetrics struct {
	ChainBuildMetricsData
	mutex sync.RWMutex
}

// ChainMetricDetail contient les détails métriques pour une chaîne individuelle
type ChainMetricDetail struct {
	RuleID          string        `json:"rule_id"`
	ChainLength     int           `json:"chain_length"`
	NodesCreated    int           `json:"nodes_created"`
	NodesReused     int           `json:"nodes_reused"`
	BuildTime       time.Duration `json:"build_time_ns"`
	Timestamp       time.Time     `json:"timestamp"`
	HashesGenerated []string      `json:"hashes_generated,omitempty"`
}

// NewChainBuildMetrics crée une nouvelle instance de métriques
func NewChainBuildMetrics() *ChainBuildMetrics {
	return &ChainBuildMetrics{
		ChainBuildMetricsData: ChainBuildMetricsData{
			ChainDetails: make([]ChainMetricDetail, 0),
		},
	}
}

// RecordChainBuild enregistre les métriques pour une chaîne construite
func (m *ChainBuildMetrics) RecordChainBuild(detail ChainMetricDetail) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.TotalChainsBuilt++
	m.TotalNodesCreated += detail.NodesCreated
	m.TotalNodesReused += detail.NodesReused
	m.TotalBuildTime += detail.BuildTime

	// Mettre à jour les moyennes
	if m.TotalChainsBuilt > 0 {
		totalNodes := m.TotalNodesCreated + m.TotalNodesReused
		m.AverageChainLength = float64(totalNodes) / float64(m.TotalChainsBuilt)
		m.AverageBuildTime = m.TotalBuildTime / time.Duration(m.TotalChainsBuilt)

		if totalNodes > 0 {
			m.SharingRatio = float64(m.TotalNodesReused) / float64(totalNodes)
		}
	}

	m.ChainDetails = append(m.ChainDetails, detail)
}

// RecordHashCacheHit enregistre un hit du cache de hash
func (m *ChainBuildMetrics) RecordHashCacheHit() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.HashCacheHits++
}

// RecordHashCacheMiss enregistre un miss du cache de hash
func (m *ChainBuildMetrics) RecordHashCacheMiss() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.HashCacheMisses++
}

// UpdateHashCacheSize met à jour la taille du cache de hash
func (m *ChainBuildMetrics) UpdateHashCacheSize(size int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.HashCacheSize = size
}

// RecordConnectionCacheHit enregistre un hit du cache de connexion
func (m *ChainBuildMetrics) RecordConnectionCacheHit() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.ConnectionCacheHits++
}

// RecordConnectionCacheMiss enregistre un miss du cache de connexion
func (m *ChainBuildMetrics) RecordConnectionCacheMiss() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.ConnectionCacheMisses++
}

// AddHashComputeTime ajoute du temps de calcul de hash
func (m *ChainBuildMetrics) AddHashComputeTime(duration time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.TotalHashComputeTime += duration
}

// GetSnapshot retourne une copie thread-safe des métriques actuelles
func (m *ChainBuildMetrics) GetSnapshot() ChainBuildMetricsData {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Copie profonde de la structure SANS le mutex
	snapshot := ChainBuildMetricsData{
		TotalChainsBuilt:      m.TotalChainsBuilt,
		TotalNodesCreated:     m.TotalNodesCreated,
		TotalNodesReused:      m.TotalNodesReused,
		AverageChainLength:    m.AverageChainLength,
		SharingRatio:          m.SharingRatio,
		HashCacheHits:         m.HashCacheHits,
		HashCacheMisses:       m.HashCacheMisses,
		HashCacheSize:         m.HashCacheSize,
		ConnectionCacheHits:   m.ConnectionCacheHits,
		ConnectionCacheMisses: m.ConnectionCacheMisses,
		TotalBuildTime:        m.TotalBuildTime,
		AverageBuildTime:      m.AverageBuildTime,
		TotalHashComputeTime:  m.TotalHashComputeTime,
		ChainDetails:          make([]ChainMetricDetail, len(m.ChainDetails)),
	}

	copy(snapshot.ChainDetails, m.ChainDetails)
	return snapshot
}

// Reset réinitialise toutes les métriques
func (m *ChainBuildMetrics) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.TotalChainsBuilt = 0
	m.TotalNodesCreated = 0
	m.TotalNodesReused = 0
	m.AverageChainLength = 0.0
	m.SharingRatio = 0.0
	m.HashCacheHits = 0
	m.HashCacheMisses = 0
	m.HashCacheSize = 0
	m.ConnectionCacheHits = 0
	m.ConnectionCacheMisses = 0
	m.TotalBuildTime = 0
	m.AverageBuildTime = 0
	m.TotalHashComputeTime = 0
	m.ChainDetails = make([]ChainMetricDetail, 0)
}

// GetHashCacheEfficiency retourne l'efficacité du cache de hash (0.0 à 1.0)
func (m *ChainBuildMetrics) GetHashCacheEfficiency() float64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	total := m.HashCacheHits + m.HashCacheMisses
	if total == 0 {
		return 0.0
	}
	return float64(m.HashCacheHits) / float64(total)
}

// GetConnectionCacheEfficiency retourne l'efficacité du cache de connexion (0.0 à 1.0)
func (m *ChainBuildMetrics) GetConnectionCacheEfficiency() float64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	total := m.ConnectionCacheHits + m.ConnectionCacheMisses
	if total == 0 {
		return 0.0
	}
	return float64(m.ConnectionCacheHits) / float64(total)
}

// GetSummary retourne un résumé formaté des métriques
func (m *ChainBuildMetrics) GetSummary() map[string]interface{} {
	snapshot := m.GetSnapshot()

	return map[string]interface{}{
		"chains": map[string]interface{}{
			"total_built":        snapshot.TotalChainsBuilt,
			"average_length":     snapshot.AverageChainLength,
			"total_build_time":   snapshot.TotalBuildTime.String(),
			"average_build_time": snapshot.AverageBuildTime.String(),
		},
		"nodes": map[string]interface{}{
			"total_created":  snapshot.TotalNodesCreated,
			"total_reused":   snapshot.TotalNodesReused,
			"sharing_ratio":  snapshot.SharingRatio,
			"reuse_rate_pct": snapshot.SharingRatio * 100,
		},
		"hash_cache": map[string]interface{}{
			"hits":           snapshot.HashCacheHits,
			"misses":         snapshot.HashCacheMisses,
			"size":           snapshot.HashCacheSize,
			"efficiency":     m.GetHashCacheEfficiency(),
			"total_time":     snapshot.TotalHashComputeTime.String(),
			"efficiency_pct": m.GetHashCacheEfficiency() * 100,
		},
		"connection_cache": map[string]interface{}{
			"hits":           snapshot.ConnectionCacheHits,
			"misses":         snapshot.ConnectionCacheMisses,
			"efficiency":     m.GetConnectionCacheEfficiency(),
			"efficiency_pct": m.GetConnectionCacheEfficiency() * 100,
		},
	}
}

// GetTopChainsByBuildTime retourne les N chaînes avec les temps de construction les plus longs
func (m *ChainBuildMetrics) GetTopChainsByBuildTime(n int) []ChainMetricDetail {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.ChainDetails) == 0 {
		return []ChainMetricDetail{}
	}

	// Copier et trier
	chains := make([]ChainMetricDetail, len(m.ChainDetails))
	copy(chains, m.ChainDetails)

	// Tri par bulle simple (suffisant pour de petites listes)
	for i := 0; i < len(chains); i++ {
		for j := i + 1; j < len(chains); j++ {
			if chains[j].BuildTime > chains[i].BuildTime {
				chains[i], chains[j] = chains[j], chains[i]
			}
		}
	}

	if n > len(chains) {
		n = len(chains)
	}

	return chains[:n]
}

// GetTopChainsByLength retourne les N chaînes les plus longues
func (m *ChainBuildMetrics) GetTopChainsByLength(n int) []ChainMetricDetail {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.ChainDetails) == 0 {
		return []ChainMetricDetail{}
	}

	// Copier et trier
	chains := make([]ChainMetricDetail, len(m.ChainDetails))
	copy(chains, m.ChainDetails)

	// Tri par bulle simple
	for i := 0; i < len(chains); i++ {
		for j := i + 1; j < len(chains); j++ {
			if chains[j].ChainLength > chains[i].ChainLength {
				chains[i], chains[j] = chains[j], chains[i]
			}
		}
	}

	if n > len(chains) {
		n = len(chains)
	}

	return chains[:n]
}
