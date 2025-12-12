// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"testing"
	"time"
)
// TestChainBuildMetrics_NewAndBasic teste la création et les opérations de base
func TestChainBuildMetrics_NewAndBasic(t *testing.T) {
	metrics := NewChainBuildMetrics()
	if metrics == nil {
		t.Fatal("NewChainBuildMetrics() retourne nil")
	}
	// Vérifier l'état initial
	snapshot := metrics.GetSnapshot()
	if snapshot.TotalChainsBuilt != 0 {
		t.Errorf("TotalChainsBuilt initial devrait être 0, obtenu %d", snapshot.TotalChainsBuilt)
	}
	if snapshot.TotalNodesCreated != 0 {
		t.Errorf("TotalNodesCreated initial devrait être 0, obtenu %d", snapshot.TotalNodesCreated)
	}
	if snapshot.TotalNodesReused != 0 {
		t.Errorf("TotalNodesReused initial devrait être 0, obtenu %d", snapshot.TotalNodesReused)
	}
	if snapshot.AverageChainLength != 0.0 {
		t.Errorf("AverageChainLength initial devrait être 0.0, obtenu %.2f", snapshot.AverageChainLength)
	}
	if snapshot.SharingRatio != 0.0 {
		t.Errorf("SharingRatio initial devrait être 0.0, obtenu %.2f", snapshot.SharingRatio)
	}
}
// TestChainBuildMetrics_RecordChainBuild teste l'enregistrement de construction de chaîne
func TestChainBuildMetrics_RecordChainBuild(t *testing.T) {
	metrics := NewChainBuildMetrics()
	detail := ChainMetricDetail{
		RuleID:          "test_rule_1",
		ChainLength:     3,
		NodesCreated:    2,
		NodesReused:     1,
		BuildTime:       100 * time.Microsecond,
		Timestamp:       time.Now(),
		HashesGenerated: []string{"hash1", "hash2", "hash3"},
	}
	metrics.RecordChainBuild(detail)
	snapshot := metrics.GetSnapshot()
	if snapshot.TotalChainsBuilt != 1 {
		t.Errorf("TotalChainsBuilt devrait être 1, obtenu %d", snapshot.TotalChainsBuilt)
	}
	if snapshot.TotalNodesCreated != 2 {
		t.Errorf("TotalNodesCreated devrait être 2, obtenu %d", snapshot.TotalNodesCreated)
	}
	if snapshot.TotalNodesReused != 1 {
		t.Errorf("TotalNodesReused devrait être 1, obtenu %d", snapshot.TotalNodesReused)
	}
	if snapshot.AverageChainLength != 3.0 {
		t.Errorf("AverageChainLength devrait être 3.0, obtenu %.2f", snapshot.AverageChainLength)
	}
	expectedRatio := 1.0 / 3.0
	if snapshot.SharingRatio != expectedRatio {
		t.Errorf("SharingRatio devrait être %.2f, obtenu %.2f", expectedRatio, snapshot.SharingRatio)
	}
	if snapshot.TotalBuildTime != 100*time.Microsecond {
		t.Errorf("TotalBuildTime incorrect: %v", snapshot.TotalBuildTime)
	}
	if snapshot.AverageBuildTime != 100*time.Microsecond {
		t.Errorf("AverageBuildTime incorrect: %v", snapshot.AverageBuildTime)
	}
	// Vérifier que les détails sont enregistrés
	if len(snapshot.ChainDetails) != 1 {
		t.Errorf("ChainDetails devrait avoir 1 élément, obtenu %d", len(snapshot.ChainDetails))
	}
	if snapshot.ChainDetails[0].RuleID != "test_rule_1" {
		t.Errorf("RuleID incorrect dans les détails: %s", snapshot.ChainDetails[0].RuleID)
	}
}
// TestChainBuildMetrics_MultipleRecords teste plusieurs enregistrements
func TestChainBuildMetrics_MultipleRecords(t *testing.T) {
	metrics := NewChainBuildMetrics()
	// Enregistrer 3 chaînes
	details := []ChainMetricDetail{
		{
			RuleID:       "rule_1",
			ChainLength:  2,
			NodesCreated: 2,
			NodesReused:  0,
			BuildTime:    100 * time.Microsecond,
			Timestamp:    time.Now(),
		},
		{
			RuleID:       "rule_2",
			ChainLength:  2,
			NodesCreated: 0,
			NodesReused:  2,
			BuildTime:    50 * time.Microsecond,
			Timestamp:    time.Now(),
		},
		{
			RuleID:       "rule_3",
			ChainLength:  3,
			NodesCreated: 1,
			NodesReused:  2,
			BuildTime:    75 * time.Microsecond,
			Timestamp:    time.Now(),
		},
	}
	for _, detail := range details {
		metrics.RecordChainBuild(detail)
	}
	snapshot := metrics.GetSnapshot()
	// Vérifications
	if snapshot.TotalChainsBuilt != 3 {
		t.Errorf("TotalChainsBuilt devrait être 3, obtenu %d", snapshot.TotalChainsBuilt)
	}
	if snapshot.TotalNodesCreated != 3 {
		t.Errorf("TotalNodesCreated devrait être 3, obtenu %d", snapshot.TotalNodesCreated)
	}
	if snapshot.TotalNodesReused != 4 {
		t.Errorf("TotalNodesReused devrait être 4, obtenu %d", snapshot.TotalNodesReused)
	}
	// Average chain length = (2+2+3) / 3 = 7/3 = 2.33
	expectedAvgLength := 7.0 / 3.0
	if snapshot.AverageChainLength != expectedAvgLength {
		t.Errorf("AverageChainLength devrait être %.2f, obtenu %.2f",
			expectedAvgLength, snapshot.AverageChainLength)
	}
	// Sharing ratio = 4 / (3+4) = 4/7 = 0.571...
	expectedRatio := 4.0 / 7.0
	if snapshot.SharingRatio != expectedRatio {
		t.Errorf("SharingRatio devrait être %.3f, obtenu %.3f",
			expectedRatio, snapshot.SharingRatio)
	}
	// Total build time = 100 + 50 + 75 = 225 µs
	expectedTotalTime := 225 * time.Microsecond
	if snapshot.TotalBuildTime != expectedTotalTime {
		t.Errorf("TotalBuildTime devrait être %v, obtenu %v",
			expectedTotalTime, snapshot.TotalBuildTime)
	}
	// Average build time = 225 / 3 = 75 µs
	expectedAvgTime := 75 * time.Microsecond
	if snapshot.AverageBuildTime != expectedAvgTime {
		t.Errorf("AverageBuildTime devrait être %v, obtenu %v",
			expectedAvgTime, snapshot.AverageBuildTime)
	}
}
// TestChainBuildMetrics_HashCache teste les métriques de cache de hash
func TestChainBuildMetrics_HashCache(t *testing.T) {
	metrics := NewChainBuildMetrics()
	// Enregistrer des hits et misses
	metrics.RecordHashCacheMiss()
	metrics.RecordHashCacheHit()
	metrics.RecordHashCacheHit()
	metrics.RecordHashCacheMiss()
	metrics.RecordHashCacheHit()
	metrics.UpdateHashCacheSize(10)
	snapshot := metrics.GetSnapshot()
	if snapshot.HashCacheHits != 3 {
		t.Errorf("HashCacheHits devrait être 3, obtenu %d", snapshot.HashCacheHits)
	}
	if snapshot.HashCacheMisses != 2 {
		t.Errorf("HashCacheMisses devrait être 2, obtenu %d", snapshot.HashCacheMisses)
	}
	if snapshot.HashCacheSize != 10 {
		t.Errorf("HashCacheSize devrait être 10, obtenu %d", snapshot.HashCacheSize)
	}
	// Efficacité = 3 / (3+2) = 0.6
	expectedEfficiency := 0.6
	efficiency := metrics.GetHashCacheEfficiency()
	if efficiency != expectedEfficiency {
		t.Errorf("Hash cache efficiency devrait être %.2f, obtenu %.2f",
			expectedEfficiency, efficiency)
	}
}
// TestChainBuildMetrics_ConnectionCache teste les métriques de cache de connexion
func TestChainBuildMetrics_ConnectionCache(t *testing.T) {
	metrics := NewChainBuildMetrics()
	metrics.RecordConnectionCacheMiss()
	metrics.RecordConnectionCacheMiss()
	metrics.RecordConnectionCacheHit()
	metrics.RecordConnectionCacheHit()
	metrics.RecordConnectionCacheHit()
	metrics.RecordConnectionCacheHit()
	snapshot := metrics.GetSnapshot()
	if snapshot.ConnectionCacheHits != 4 {
		t.Errorf("ConnectionCacheHits devrait être 4, obtenu %d", snapshot.ConnectionCacheHits)
	}
	if snapshot.ConnectionCacheMisses != 2 {
		t.Errorf("ConnectionCacheMisses devrait être 2, obtenu %d", snapshot.ConnectionCacheMisses)
	}
	// Efficacité = 4 / (4+2) = 0.666...
	expectedEfficiency := 4.0 / 6.0
	efficiency := metrics.GetConnectionCacheEfficiency()
	if efficiency != expectedEfficiency {
		t.Errorf("Connection cache efficiency devrait être %.3f, obtenu %.3f",
			expectedEfficiency, efficiency)
	}
}
// TestChainBuildMetrics_Reset teste la réinitialisation des métriques
func TestChainBuildMetrics_Reset(t *testing.T) {
	metrics := NewChainBuildMetrics()
	// Ajouter des données
	detail := ChainMetricDetail{
		RuleID:       "test_rule",
		ChainLength:  5,
		NodesCreated: 3,
		NodesReused:  2,
		BuildTime:    100 * time.Microsecond,
		Timestamp:    time.Now(),
	}
	metrics.RecordChainBuild(detail)
	metrics.RecordHashCacheHit()
	metrics.RecordHashCacheMiss()
	metrics.RecordConnectionCacheHit()
	metrics.UpdateHashCacheSize(5)
	metrics.AddHashComputeTime(50 * time.Microsecond)
	// Vérifier que les données sont présentes
	snapshot := metrics.GetSnapshot()
	if snapshot.TotalChainsBuilt == 0 {
		t.Error("Les données devraient être présentes avant Reset")
	}
	// Reset
	metrics.Reset()
	// Vérifier que tout est réinitialisé
	snapshot = metrics.GetSnapshot()
	if snapshot.TotalChainsBuilt != 0 {
		t.Errorf("TotalChainsBuilt devrait être 0 après Reset, obtenu %d", snapshot.TotalChainsBuilt)
	}
	if snapshot.TotalNodesCreated != 0 {
		t.Errorf("TotalNodesCreated devrait être 0 après Reset, obtenu %d", snapshot.TotalNodesCreated)
	}
	if snapshot.TotalNodesReused != 0 {
		t.Errorf("TotalNodesReused devrait être 0 après Reset, obtenu %d", snapshot.TotalNodesReused)
	}
	if snapshot.HashCacheHits != 0 {
		t.Errorf("HashCacheHits devrait être 0 après Reset, obtenu %d", snapshot.HashCacheHits)
	}
	if snapshot.HashCacheMisses != 0 {
		t.Errorf("HashCacheMisses devrait être 0 après Reset, obtenu %d", snapshot.HashCacheMisses)
	}
	if snapshot.HashCacheSize != 0 {
		t.Errorf("HashCacheSize devrait être 0 après Reset, obtenu %d", snapshot.HashCacheSize)
	}
	if snapshot.ConnectionCacheHits != 0 {
		t.Errorf("ConnectionCacheHits devrait être 0 après Reset, obtenu %d", snapshot.ConnectionCacheHits)
	}
	if snapshot.ConnectionCacheMisses != 0 {
		t.Errorf("ConnectionCacheMisses devrait être 0 après Reset, obtenu %d", snapshot.ConnectionCacheMisses)
	}
	if snapshot.TotalBuildTime != 0 {
		t.Errorf("TotalBuildTime devrait être 0 après Reset, obtenu %v", snapshot.TotalBuildTime)
	}
	if snapshot.TotalHashComputeTime != 0 {
		t.Errorf("TotalHashComputeTime devrait être 0 après Reset, obtenu %v", snapshot.TotalHashComputeTime)
	}
	if len(snapshot.ChainDetails) != 0 {
		t.Errorf("ChainDetails devrait être vide après Reset, obtenu %d éléments", len(snapshot.ChainDetails))
	}
}
// TestChainBuildMetrics_GetSummary teste le format du résumé
func TestChainBuildMetrics_GetSummary(t *testing.T) {
	metrics := NewChainBuildMetrics()
	detail := ChainMetricDetail{
		RuleID:       "test_rule",
		ChainLength:  3,
		NodesCreated: 2,
		NodesReused:  1,
		BuildTime:    100 * time.Microsecond,
		Timestamp:    time.Now(),
	}
	metrics.RecordChainBuild(detail)
	metrics.RecordHashCacheHit()
	metrics.RecordHashCacheMiss()
	metrics.UpdateHashCacheSize(5)
	summary := metrics.GetSummary()
	if summary == nil {
		t.Fatal("GetSummary() retourne nil")
	}
	// Vérifier la présence des sections principales
	sections := []string{"chains", "nodes", "hash_cache", "connection_cache"}
	for _, section := range sections {
		if _, exists := summary[section]; !exists {
			t.Errorf("Section '%s' manquante dans le résumé", section)
		}
	}
	// Vérifier quelques valeurs dans la section chains
	chains := summary["chains"].(map[string]interface{})
	if totalBuilt, ok := chains["total_built"].(int); !ok || totalBuilt != 1 {
		t.Errorf("chains.total_built incorrect: %v", chains["total_built"])
	}
	// Vérifier la section nodes
	nodes := summary["nodes"].(map[string]interface{})
	if totalCreated, ok := nodes["total_created"].(int); !ok || totalCreated != 2 {
		t.Errorf("nodes.total_created incorrect: %v", nodes["total_created"])
	}
	if totalReused, ok := nodes["total_reused"].(int); !ok || totalReused != 1 {
		t.Errorf("nodes.total_reused incorrect: %v", nodes["total_reused"])
	}
	// Vérifier la section hash_cache
	hashCache := summary["hash_cache"].(map[string]interface{})
	if hits, ok := hashCache["hits"].(int); !ok || hits != 1 {
		t.Errorf("hash_cache.hits incorrect: %v", hashCache["hits"])
	}
	if size, ok := hashCache["size"].(int); !ok || size != 5 {
		t.Errorf("hash_cache.size incorrect: %v", hashCache["size"])
	}
}
// TestChainBuildMetrics_TopChainsByBuildTime teste le classement par temps de construction
func TestChainBuildMetrics_TopChainsByBuildTime(t *testing.T) {
	metrics := NewChainBuildMetrics()
	// Ajouter des chaînes avec différents temps de construction
	details := []ChainMetricDetail{
		{RuleID: "fast", ChainLength: 1, BuildTime: 10 * time.Microsecond},
		{RuleID: "slow", ChainLength: 1, BuildTime: 100 * time.Microsecond},
		{RuleID: "medium", ChainLength: 1, BuildTime: 50 * time.Microsecond},
		{RuleID: "very_slow", ChainLength: 1, BuildTime: 200 * time.Microsecond},
	}
	for _, detail := range details {
		metrics.RecordChainBuild(detail)
	}
	// Obtenir le top 3
	top3 := metrics.GetTopChainsByBuildTime(3)
	if len(top3) != 3 {
		t.Errorf("Attendu 3 chaînes, obtenu %d", len(top3))
	}
	// Vérifier l'ordre (décroissant)
	if top3[0].RuleID != "very_slow" {
		t.Errorf("Premier devrait être 'very_slow', obtenu '%s'", top3[0].RuleID)
	}
	if top3[1].RuleID != "slow" {
		t.Errorf("Deuxième devrait être 'slow', obtenu '%s'", top3[1].RuleID)
	}
	if top3[2].RuleID != "medium" {
		t.Errorf("Troisième devrait être 'medium', obtenu '%s'", top3[2].RuleID)
	}
	// Vérifier que c'est trié
	for i := 0; i < len(top3)-1; i++ {
		if top3[i].BuildTime < top3[i+1].BuildTime {
			t.Errorf("Liste non triée: index %d (%v) < index %d (%v)",
				i, top3[i].BuildTime, i+1, top3[i+1].BuildTime)
		}
	}
}
// TestChainBuildMetrics_TopChainsByLength teste le classement par longueur
func TestChainBuildMetrics_TopChainsByLength(t *testing.T) {
	metrics := NewChainBuildMetrics()
	// Ajouter des chaînes avec différentes longueurs
	details := []ChainMetricDetail{
		{RuleID: "short", ChainLength: 2, BuildTime: 10 * time.Microsecond},
		{RuleID: "long", ChainLength: 5, BuildTime: 10 * time.Microsecond},
		{RuleID: "medium", ChainLength: 3, BuildTime: 10 * time.Microsecond},
		{RuleID: "very_long", ChainLength: 10, BuildTime: 10 * time.Microsecond},
	}
	for _, detail := range details {
		metrics.RecordChainBuild(detail)
	}
	// Obtenir le top 3
	top3 := metrics.GetTopChainsByLength(3)
	if len(top3) != 3 {
		t.Errorf("Attendu 3 chaînes, obtenu %d", len(top3))
	}
	// Vérifier l'ordre (décroissant)
	if top3[0].RuleID != "very_long" {
		t.Errorf("Premier devrait être 'very_long', obtenu '%s'", top3[0].RuleID)
	}
	if top3[1].RuleID != "long" {
		t.Errorf("Deuxième devrait être 'long', obtenu '%s'", top3[1].RuleID)
	}
	if top3[2].RuleID != "medium" {
		t.Errorf("Troisième devrait être 'medium', obtenu '%s'", top3[2].RuleID)
	}
	// Vérifier que c'est trié
	for i := 0; i < len(top3)-1; i++ {
		if top3[i].ChainLength < top3[i+1].ChainLength {
			t.Errorf("Liste non triée: index %d (len=%d) < index %d (len=%d)",
				i, top3[i].ChainLength, i+1, top3[i+1].ChainLength)
		}
	}
}
// TestChainBuildMetrics_EmptyTopChains teste les fonctions top avec aucune donnée
func TestChainBuildMetrics_EmptyTopChains(t *testing.T) {
	metrics := NewChainBuildMetrics()
	topByTime := metrics.GetTopChainsByBuildTime(5)
	if len(topByTime) != 0 {
		t.Errorf("topByTime devrait être vide, obtenu %d éléments", len(topByTime))
	}
	topByLength := metrics.GetTopChainsByLength(5)
	if len(topByLength) != 0 {
		t.Errorf("topByLength devrait être vide, obtenu %d éléments", len(topByLength))
	}
}
// TestChainBuildMetrics_AddHashComputeTime teste l'accumulation du temps de calcul de hash
func TestChainBuildMetrics_AddHashComputeTime(t *testing.T) {
	metrics := NewChainBuildMetrics()
	metrics.AddHashComputeTime(50 * time.Microsecond)
	metrics.AddHashComputeTime(30 * time.Microsecond)
	metrics.AddHashComputeTime(20 * time.Microsecond)
	snapshot := metrics.GetSnapshot()
	expected := 100 * time.Microsecond
	if snapshot.TotalHashComputeTime != expected {
		t.Errorf("TotalHashComputeTime devrait être %v, obtenu %v",
			expected, snapshot.TotalHashComputeTime)
	}
}
// TestChainBuildMetrics_CacheEfficiencyZero teste l'efficacité avec aucune donnée
func TestChainBuildMetrics_CacheEfficiencyZero(t *testing.T) {
	metrics := NewChainBuildMetrics()
	hashEfficiency := metrics.GetHashCacheEfficiency()
	if hashEfficiency != 0.0 {
		t.Errorf("Hash cache efficiency devrait être 0.0, obtenu %.2f", hashEfficiency)
	}
	connEfficiency := metrics.GetConnectionCacheEfficiency()
	if connEfficiency != 0.0 {
		t.Errorf("Connection cache efficiency devrait être 0.0, obtenu %.2f", connEfficiency)
	}
}
// TestChainBuildMetrics_ThreadSafety teste la sécurité des threads (basic)
func TestChainBuildMetrics_ThreadSafety(t *testing.T) {
	metrics := NewChainBuildMetrics()
	// Lancer plusieurs goroutines qui modifient les métriques
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			detail := ChainMetricDetail{
				RuleID:       "concurrent_rule",
				ChainLength:  1,
				NodesCreated: 1,
				BuildTime:    10 * time.Microsecond,
				Timestamp:    time.Now(),
			}
			metrics.RecordChainBuild(detail)
			metrics.RecordHashCacheHit()
			metrics.RecordHashCacheMiss()
			metrics.RecordConnectionCacheHit()
			_ = metrics.GetSnapshot()
			done <- true
		}(i)
	}
	// Attendre que toutes les goroutines finissent
	for i := 0; i < 10; i++ {
		<-done
	}
	snapshot := metrics.GetSnapshot()
	if snapshot.TotalChainsBuilt != 10 {
		t.Errorf("TotalChainsBuilt devrait être 10, obtenu %d", snapshot.TotalChainsBuilt)
	}
}