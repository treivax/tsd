// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"sync"
	"time"
)

// BetaChainMetricsData contient les données de métriques sans le mutex
type BetaChainMetricsData struct {
	// Métriques de chaînes
	TotalChainsBuilt   int     `json:"total_chains_built"`
	TotalNodesCreated  int     `json:"total_nodes_created"`
	TotalNodesReused   int     `json:"total_nodes_reused"`
	AverageChainLength float64 `json:"average_chain_length"`
	SharingRatio       float64 `json:"sharing_ratio"`

	// Métriques de partage (compatibilité avec BetaBuildMetrics)
	TotalJoinNodesRequested int64 `json:"total_join_nodes_requested"`
	SharedJoinNodesReused   int64 `json:"shared_join_nodes_reused"`
	UniqueJoinNodesCreated  int64 `json:"unique_join_nodes_created"`

	// Métriques de jointures
	TotalJoinsExecuted     int           `json:"total_joins_executed"`
	AverageJoinTime        time.Duration `json:"average_join_time_ns"`
	TotalJoinTime          time.Duration `json:"total_join_time_ns"`
	AverageJoinSelectivity float64       `json:"average_join_selectivity"`
	AverageResultSize      float64       `json:"average_result_size"`

	// Métriques de cache de hash
	HashCacheHits   int `json:"hash_cache_hits"`
	HashCacheMisses int `json:"hash_cache_misses"`
	HashCacheSize   int `json:"hash_cache_size"`

	// Métriques de cache de jointure
	JoinCacheHits      int `json:"join_cache_hits"`
	JoinCacheMisses    int `json:"join_cache_misses"`
	JoinCacheSize      int `json:"join_cache_size"`
	JoinCacheEvictions int `json:"join_cache_evictions"`

	// Métriques de cache de connexion
	ConnectionCacheHits   int `json:"connection_cache_hits"`
	ConnectionCacheMisses int `json:"connection_cache_misses"`

	// Métriques de cache de préfixe
	PrefixCacheHits   int `json:"prefix_cache_hits"`
	PrefixCacheMisses int `json:"prefix_cache_misses"`
	PrefixCacheSize   int `json:"prefix_cache_size"`

	// Métriques de temps
	TotalBuildTime       time.Duration `json:"total_build_time_ns"`
	AverageBuildTime     time.Duration `json:"average_build_time_ns"`
	TotalHashComputeTime time.Duration `json:"total_hash_compute_time_ns"`

	// Détails par chaîne
	ChainDetails []BetaChainMetricDetail `json:"chain_details,omitempty"`
}

// BetaChainMetrics collecte les métriques de performance pour la construction des chaînes beta
type BetaChainMetrics struct {
	BetaChainMetricsData
	mutex sync.RWMutex
}

// BetaChainMetricDetail contient les détails métriques pour une chaîne beta individuelle
type BetaChainMetricDetail struct {
	RuleID          string        `json:"rule_id"`
	ChainLength     int           `json:"chain_length"`
	NodesCreated    int           `json:"nodes_created"`
	NodesReused     int           `json:"nodes_reused"`
	BuildTime       time.Duration `json:"build_time_ns"`
	Timestamp       time.Time     `json:"timestamp"`
	HashesGenerated []string      `json:"hashes_generated,omitempty"`

	// Métriques de jointure spécifiques
	JoinsExecuted      int           `json:"joins_executed"`
	TotalJoinTime      time.Duration `json:"total_join_time_ns"`
	AverageSelectivity float64       `json:"average_selectivity"`
	TotalResultSize    int           `json:"total_result_size"`
}

// JoinMetricDetail contient les métriques pour une jointure individuelle
type JoinMetricDetail struct {
	JoinNodeID     string        `json:"join_node_id"`
	LeftInputSize  int           `json:"left_input_size"`
	RightInputSize int           `json:"right_input_size"`
	ResultSize     int           `json:"result_size"`
	Selectivity    float64       `json:"selectivity"`
	ExecutionTime  time.Duration `json:"execution_time_ns"`
	Timestamp      time.Time     `json:"timestamp"`
}

// NewBetaChainMetrics crée une nouvelle instance de métriques
func NewBetaChainMetrics() *BetaChainMetrics {
	return &BetaChainMetrics{
		BetaChainMetricsData: BetaChainMetricsData{
			ChainDetails: make([]BetaChainMetricDetail, 0),
		},
	}
}

// RecordChainBuild enregistre les métriques pour une chaîne construite
func (m *BetaChainMetrics) RecordChainBuild(detail BetaChainMetricDetail) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.TotalChainsBuilt++
	m.TotalNodesCreated += detail.NodesCreated
	m.TotalNodesReused += detail.NodesReused
	m.TotalBuildTime += detail.BuildTime
	m.TotalJoinsExecuted += detail.JoinsExecuted
	m.TotalJoinTime += detail.TotalJoinTime

	// Update sharing metrics
	m.TotalJoinNodesRequested += int64(detail.NodesCreated + detail.NodesReused)
	m.SharedJoinNodesReused += int64(detail.NodesReused)
	m.UniqueJoinNodesCreated += int64(detail.NodesCreated)

	// Mettre à jour les moyennes
	if m.TotalChainsBuilt > 0 {
		totalNodes := m.TotalNodesCreated + m.TotalNodesReused
		m.AverageChainLength = float64(totalNodes) / float64(m.TotalChainsBuilt)
		m.AverageBuildTime = m.TotalBuildTime / time.Duration(m.TotalChainsBuilt)

		if totalNodes > 0 {
			m.SharingRatio = float64(m.TotalNodesReused) / float64(totalNodes)
		}
	}

	if m.TotalJoinsExecuted > 0 {
		m.AverageJoinTime = m.TotalJoinTime / time.Duration(m.TotalJoinsExecuted)
	}

	m.ChainDetails = append(m.ChainDetails, detail)
}

// RecordJoinExecution enregistre l'exécution d'une jointure
func (m *BetaChainMetrics) RecordJoinExecution(leftSize, rightSize, resultSize int, duration time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.TotalJoinsExecuted++
	m.TotalJoinTime += duration

	// Calculer la sélectivité: résultat / (gauche * droite)
	maxPossible := leftSize * rightSize
	if maxPossible > 0 {
		selectivity := float64(resultSize) / float64(maxPossible)

		// Mettre à jour la moyenne de sélectivité
		totalSelectivity := m.AverageJoinSelectivity * float64(m.TotalJoinsExecuted-1)
		m.AverageJoinSelectivity = (totalSelectivity + selectivity) / float64(m.TotalJoinsExecuted)
	}

	// Mettre à jour la taille moyenne des résultats
	totalResultSize := m.AverageResultSize * float64(m.TotalJoinsExecuted-1)
	m.AverageResultSize = (totalResultSize + float64(resultSize)) / float64(m.TotalJoinsExecuted)

	// Mettre à jour le temps moyen
	m.AverageJoinTime = m.TotalJoinTime / time.Duration(m.TotalJoinsExecuted)
}

// RecordHashCacheHit enregistre un hit du cache de hash
func (m *BetaChainMetrics) RecordHashCacheHit() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.HashCacheHits++
}

// RecordHashCacheMiss enregistre un miss du cache de hash
func (m *BetaChainMetrics) RecordHashCacheMiss() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.HashCacheMisses++
}

// UpdateHashCacheSize met à jour la taille du cache de hash
func (m *BetaChainMetrics) UpdateHashCacheSize(size int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.HashCacheSize = size
}

// RecordJoinCacheHit enregistre un hit du cache de jointure
func (m *BetaChainMetrics) RecordJoinCacheHit() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.JoinCacheHits++
}

// RecordJoinCacheMiss enregistre un miss du cache de jointure
func (m *BetaChainMetrics) RecordJoinCacheMiss() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.JoinCacheMisses++
}

// UpdateJoinCacheSize met à jour la taille du cache de jointure
func (m *BetaChainMetrics) UpdateJoinCacheSize(size int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.JoinCacheSize = size
}

// RecordJoinCacheEviction enregistre une éviction du cache de jointure
func (m *BetaChainMetrics) RecordJoinCacheEviction() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.JoinCacheEvictions++
}

// RecordConnectionCacheHit enregistre un hit du cache de connexion
func (m *BetaChainMetrics) RecordConnectionCacheHit() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.ConnectionCacheHits++
}

// RecordConnectionCacheMiss enregistre un miss du cache de connexion
func (m *BetaChainMetrics) RecordConnectionCacheMiss() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.ConnectionCacheMisses++
}

// RecordPrefixCacheHit enregistre un hit du cache de préfixe
func (m *BetaChainMetrics) RecordPrefixCacheHit() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.PrefixCacheHits++
}

// RecordPrefixCacheMiss enregistre un miss du cache de préfixe
func (m *BetaChainMetrics) RecordPrefixCacheMiss() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.PrefixCacheMisses++
}

// UpdatePrefixCacheSize met à jour la taille du cache de préfixe
func (m *BetaChainMetrics) UpdatePrefixCacheSize(size int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.PrefixCacheSize = size
}

// AddHashComputeTime ajoute du temps de calcul de hash
func (m *BetaChainMetrics) AddHashComputeTime(duration time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.TotalHashComputeTime += duration
}

// GetSnapshot retourne une copie thread-safe des métriques actuelles
func (m *BetaChainMetrics) GetSnapshot() BetaChainMetricsData {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Copie profonde de la structure SANS le mutex
	snapshot := BetaChainMetricsData{
		TotalChainsBuilt:        m.TotalChainsBuilt,
		TotalNodesCreated:       m.TotalNodesCreated,
		TotalNodesReused:        m.TotalNodesReused,
		AverageChainLength:      m.AverageChainLength,
		SharingRatio:            m.SharingRatio,
		TotalJoinNodesRequested: m.TotalJoinNodesRequested,
		SharedJoinNodesReused:   m.SharedJoinNodesReused,
		UniqueJoinNodesCreated:  m.UniqueJoinNodesCreated,
		TotalJoinsExecuted:      m.TotalJoinsExecuted,
		AverageJoinTime:         m.AverageJoinTime,
		TotalJoinTime:           m.TotalJoinTime,
		AverageJoinSelectivity:  m.AverageJoinSelectivity,
		AverageResultSize:       m.AverageResultSize,
		HashCacheHits:           m.HashCacheHits,
		HashCacheMisses:         m.HashCacheMisses,
		HashCacheSize:           m.HashCacheSize,
		JoinCacheHits:           m.JoinCacheHits,
		JoinCacheMisses:         m.JoinCacheMisses,
		JoinCacheSize:           m.JoinCacheSize,
		JoinCacheEvictions:      m.JoinCacheEvictions,
		ConnectionCacheHits:     m.ConnectionCacheHits,
		ConnectionCacheMisses:   m.ConnectionCacheMisses,
		PrefixCacheHits:         m.PrefixCacheHits,
		PrefixCacheMisses:       m.PrefixCacheMisses,
		PrefixCacheSize:         m.PrefixCacheSize,
		TotalBuildTime:          m.TotalBuildTime,
		AverageBuildTime:        m.AverageBuildTime,
		TotalHashComputeTime:    m.TotalHashComputeTime,
		ChainDetails:            make([]BetaChainMetricDetail, len(m.ChainDetails)),
	}

	copy(snapshot.ChainDetails, m.ChainDetails)
	return snapshot
}

// Reset réinitialise toutes les métriques
func (m *BetaChainMetrics) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.TotalChainsBuilt = 0
	m.TotalNodesCreated = 0
	m.TotalNodesReused = 0
	m.AverageChainLength = 0.0
	m.SharingRatio = 0.0
	m.TotalJoinNodesRequested = 0
	m.SharedJoinNodesReused = 0
	m.UniqueJoinNodesCreated = 0
	m.TotalJoinsExecuted = 0
	m.AverageJoinTime = 0
	m.TotalJoinTime = 0
	m.AverageJoinSelectivity = 0.0
	m.AverageResultSize = 0.0
	m.HashCacheHits = 0
	m.HashCacheMisses = 0
	m.HashCacheSize = 0
	m.JoinCacheHits = 0
	m.JoinCacheMisses = 0
	m.JoinCacheSize = 0
	m.JoinCacheEvictions = 0
	m.ConnectionCacheHits = 0
	m.ConnectionCacheMisses = 0
	m.PrefixCacheHits = 0
	m.PrefixCacheMisses = 0
	m.PrefixCacheSize = 0
	m.TotalBuildTime = 0
	m.AverageBuildTime = 0
	m.TotalHashComputeTime = 0
	m.ChainDetails = make([]BetaChainMetricDetail, 0)
}

// GetHashCacheEfficiency retourne l'efficacité du cache de hash (0.0 à 1.0)
func (m *BetaChainMetrics) GetHashCacheEfficiency() float64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	total := m.HashCacheHits + m.HashCacheMisses
	if total == 0 {
		return 0.0
	}
	return float64(m.HashCacheHits) / float64(total)
}

// GetJoinCacheEfficiency retourne l'efficacité du cache de jointure (0.0 à 1.0)
func (m *BetaChainMetrics) GetJoinCacheEfficiency() float64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	total := m.JoinCacheHits + m.JoinCacheMisses
	if total == 0 {
		return 0.0
	}
	return float64(m.JoinCacheHits) / float64(total)
}

// GetConnectionCacheEfficiency retourne l'efficacité du cache de connexion (0.0 à 1.0)
func (m *BetaChainMetrics) GetConnectionCacheEfficiency() float64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	total := m.ConnectionCacheHits + m.ConnectionCacheMisses
	if total == 0 {
		return 0.0
	}
	return float64(m.ConnectionCacheHits) / float64(total)
}

// GetPrefixCacheEfficiency retourne l'efficacité du cache de préfixe (0.0 à 1.0)
func (m *BetaChainMetrics) GetPrefixCacheEfficiency() float64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	total := m.PrefixCacheHits + m.PrefixCacheMisses
	if total == 0 {
		return 0.0
	}
	return float64(m.PrefixCacheHits) / float64(total)
}

// GetSummary retourne un résumé formaté des métriques
func (m *BetaChainMetrics) GetSummary() map[string]interface{} {
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
		"joins": map[string]interface{}{
			"total_executed":      snapshot.TotalJoinsExecuted,
			"total_time":          snapshot.TotalJoinTime.String(),
			"average_time":        snapshot.AverageJoinTime.String(),
			"average_selectivity": snapshot.AverageJoinSelectivity,
			"average_result_size": snapshot.AverageResultSize,
			"selectivity_pct":     snapshot.AverageJoinSelectivity * 100,
		},
		"hash_cache": map[string]interface{}{
			"hits":           snapshot.HashCacheHits,
			"misses":         snapshot.HashCacheMisses,
			"size":           snapshot.HashCacheSize,
			"efficiency":     m.GetHashCacheEfficiency(),
			"total_time":     snapshot.TotalHashComputeTime.String(),
			"efficiency_pct": m.GetHashCacheEfficiency() * 100,
		},
		"join_cache": map[string]interface{}{
			"hits":           snapshot.JoinCacheHits,
			"misses":         snapshot.JoinCacheMisses,
			"size":           snapshot.JoinCacheSize,
			"evictions":      snapshot.JoinCacheEvictions,
			"efficiency":     m.GetJoinCacheEfficiency(),
			"efficiency_pct": m.GetJoinCacheEfficiency() * 100,
		},
		"connection_cache": map[string]interface{}{
			"hits":           snapshot.ConnectionCacheHits,
			"misses":         snapshot.ConnectionCacheMisses,
			"efficiency":     m.GetConnectionCacheEfficiency(),
			"efficiency_pct": m.GetConnectionCacheEfficiency() * 100,
		},
		"prefix_cache": map[string]interface{}{
			"hits":           snapshot.PrefixCacheHits,
			"misses":         snapshot.PrefixCacheMisses,
			"size":           snapshot.PrefixCacheSize,
			"efficiency":     m.GetPrefixCacheEfficiency(),
			"efficiency_pct": m.GetPrefixCacheEfficiency() * 100,
		},
	}
}

// GetTopChainsByBuildTime retourne les N chaînes avec les temps de construction les plus longs
func (m *BetaChainMetrics) GetTopChainsByBuildTime(n int) []BetaChainMetricDetail {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.ChainDetails) == 0 {
		return []BetaChainMetricDetail{}
	}

	// Copier et trier
	chains := make([]BetaChainMetricDetail, len(m.ChainDetails))
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
func (m *BetaChainMetrics) GetTopChainsByLength(n int) []BetaChainMetricDetail {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.ChainDetails) == 0 {
		return []BetaChainMetricDetail{}
	}

	// Copier et trier
	chains := make([]BetaChainMetricDetail, len(m.ChainDetails))
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

// GetTopChainsByJoinTime retourne les N chaînes avec les temps de jointure les plus longs
func (m *BetaChainMetrics) GetTopChainsByJoinTime(n int) []BetaChainMetricDetail {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.ChainDetails) == 0 {
		return []BetaChainMetricDetail{}
	}

	// Copier et trier
	chains := make([]BetaChainMetricDetail, len(m.ChainDetails))
	copy(chains, m.ChainDetails)

	// Tri par bulle simple
	for i := 0; i < len(chains); i++ {
		for j := i + 1; j < len(chains); j++ {
			if chains[j].TotalJoinTime > chains[i].TotalJoinTime {
				chains[i], chains[j] = chains[j], chains[i]
			}
		}
	}

	if n > len(chains) {
		n = len(chains)
	}

	return chains[:n]
}

// GetJoinPerformanceStats retourne des statistiques détaillées sur les performances de jointure
func (m *BetaChainMetrics) GetJoinPerformanceStats() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := map[string]interface{}{
		"total_joins":         m.TotalJoinsExecuted,
		"total_time":          m.TotalJoinTime.String(),
		"average_time":        m.AverageJoinTime.String(),
		"average_selectivity": m.AverageJoinSelectivity,
		"average_result_size": m.AverageResultSize,
	}

	if m.TotalJoinsExecuted > 0 {
		stats["average_time_per_join_ms"] = float64(m.AverageJoinTime.Microseconds()) / 1000.0
		stats["throughput_joins_per_sec"] = float64(m.TotalJoinsExecuted) / m.TotalJoinTime.Seconds()
	}

	return stats
}

// GetCacheStats retourne un résumé de tous les caches
func (m *BetaChainMetrics) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"hash_cache": map[string]interface{}{
			"hits":           m.HashCacheHits,
			"misses":         m.HashCacheMisses,
			"size":           m.HashCacheSize,
			"efficiency":     m.GetHashCacheEfficiency(),
			"efficiency_pct": m.GetHashCacheEfficiency() * 100,
		},
		"join_cache": map[string]interface{}{
			"hits":           m.JoinCacheHits,
			"misses":         m.JoinCacheMisses,
			"size":           m.JoinCacheSize,
			"evictions":      m.JoinCacheEvictions,
			"efficiency":     m.GetJoinCacheEfficiency(),
			"efficiency_pct": m.GetJoinCacheEfficiency() * 100,
		},
		"connection_cache": map[string]interface{}{
			"hits":           m.ConnectionCacheHits,
			"misses":         m.ConnectionCacheMisses,
			"efficiency":     m.GetConnectionCacheEfficiency(),
			"efficiency_pct": m.GetConnectionCacheEfficiency() * 100,
		},
		"prefix_cache": map[string]interface{}{
			"hits":           m.PrefixCacheHits,
			"misses":         m.PrefixCacheMisses,
			"size":           m.PrefixCacheSize,
			"efficiency":     m.GetPrefixCacheEfficiency(),
			"efficiency_pct": m.GetPrefixCacheEfficiency() * 100,
		},
	}
}
