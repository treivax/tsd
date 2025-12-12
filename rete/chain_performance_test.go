// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"fmt"
	"testing"
	"time"
)

// TestPerformance_LargeRuleset_100Rules teste les performances avec 100 règles similaires
func TestPerformance_LargeRuleset_100Rules(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Créer 100 règles similaires avec des conditions légèrement différentes
	startTime := time.Now()
	for i := 0; i < 100; i++ {
		ruleID := fmt.Sprintf("rule_%d", i)
		// Créer des conditions similaires pour favoriser le partage
		conditions := []SimpleCondition{
			{
				Type:     "binaryOperation",
				Left:     map[string]interface{}{"type": "variable", "name": "age"},
				Operator: ">",
				Right:    map[string]interface{}{"type": "literal", "value": float64(i % 10)},
			},
			{
				Type:     "binaryOperation",
				Left:     map[string]interface{}{"type": "variable", "name": "status"},
				Operator: "==",
				Right:    map[string]interface{}{"type": "literal", "value": "active"},
			},
		}
		// Construire la chaîne
		builder := NewAlphaChainBuilderWithMetrics(network, storage, network.ChainMetrics)
		chain, err := builder.BuildChain(conditions, "person", network.RootNode, ruleID)
		if err != nil {
			t.Fatalf("Erreur construction chaîne pour %s: %v", ruleID, err)
		}
		if chain == nil {
			t.Fatalf("Chaîne nil pour %s", ruleID)
		}
	}
	elapsed := time.Since(startTime)
	t.Logf("⏱️  Construction de 100 règles terminée en %v", elapsed)
	// Vérifier les métriques
	metrics := network.GetChainMetrics()
	if metrics == nil {
		t.Fatal("Métriques nil")
	}
	snapshot := metrics.GetSnapshot()
	t.Logf("📊 Métriques de performance:")
	t.Logf("  - Total chaînes construites: %d", snapshot.TotalChainsBuilt)
	t.Logf("  - Nœuds créés: %d", snapshot.TotalNodesCreated)
	t.Logf("  - Nœuds réutilisés: %d", snapshot.TotalNodesReused)
	t.Logf("  - Longueur moyenne de chaîne: %.2f", snapshot.AverageChainLength)
	t.Logf("  - Ratio de partage: %.2f%%", snapshot.SharingRatio*100)
	t.Logf("  - Temps moyen de construction: %v", snapshot.AverageBuildTime)
	// Vérifier les métriques de cache
	t.Logf("  - Cache hash - hits: %d, misses: %d, efficacité: %.2f%%",
		snapshot.HashCacheHits, snapshot.HashCacheMisses,
		metrics.GetHashCacheEfficiency()*100)
	t.Logf("  - Cache connexion - hits: %d, misses: %d, efficacité: %.2f%%",
		snapshot.ConnectionCacheHits, snapshot.ConnectionCacheMisses,
		metrics.GetConnectionCacheEfficiency()*100)
	// Assertions
	if snapshot.TotalChainsBuilt != 100 {
		t.Errorf("Attendu 100 chaînes, obtenu %d", snapshot.TotalChainsBuilt)
	}
	// Vérifier que le partage fonctionne (au moins 50% de réutilisation)
	if snapshot.SharingRatio < 0.5 {
		t.Errorf("Ratio de partage trop faible: %.2f%% (attendu >= 50%%)", snapshot.SharingRatio*100)
	}
	// Vérifier que le cache de hash est efficace
	if metrics.GetHashCacheEfficiency() < 0.3 {
		t.Errorf("Efficacité du cache de hash trop faible: %.2f%% (attendu >= 30%%)",
			metrics.GetHashCacheEfficiency()*100)
	}
}

// TestPerformance_LargeRuleset_1000Rules teste les performances avec 1000 règles variées
func TestPerformance_LargeRuleset_1000Rules(t *testing.T) {
	if testing.Short() {
		t.Skip("Test long ignoré en mode short")
	}
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	startTime := time.Now()
	// Créer 1000 règles avec une plus grande variété
	for i := 0; i < 1000; i++ {
		ruleID := fmt.Sprintf("rule_%d", i)
		// Varier les conditions pour tester différents scénarios
		numConditions := (i % 5) + 1 // 1 à 5 conditions
		conditions := make([]SimpleCondition, numConditions)
		for j := 0; j < numConditions; j++ {
			conditions[j] = SimpleCondition{
				Type:     "binaryOperation",
				Left:     map[string]interface{}{"type": "variable", "name": fmt.Sprintf("field%d", j)},
				Operator: selectOperator(i, j),
				Right:    map[string]interface{}{"type": "literal", "value": float64((i + j) % 100)},
			}
		}
		builder := NewAlphaChainBuilderWithMetrics(network, storage, network.ChainMetrics)
		chain, err := builder.BuildChain(conditions, "entity", network.RootNode, ruleID)
		if err != nil {
			t.Fatalf("Erreur construction chaîne pour %s: %v", ruleID, err)
		}
		if chain == nil {
			t.Fatalf("Chaîne nil pour %s", ruleID)
		}
	}
	elapsed := time.Since(startTime)
	t.Logf("⏱️  Construction de 1000 règles terminée en %v", elapsed)
	metrics := network.GetChainMetrics()
	snapshot := metrics.GetSnapshot()
	t.Logf("📊 Métriques de performance (1000 règles):")
	t.Logf("  - Total chaînes construites: %d", snapshot.TotalChainsBuilt)
	t.Logf("  - Nœuds créés: %d", snapshot.TotalNodesCreated)
	t.Logf("  - Nœuds réutilisés: %d", snapshot.TotalNodesReused)
	t.Logf("  - Longueur moyenne de chaîne: %.2f", snapshot.AverageChainLength)
	t.Logf("  - Ratio de partage: %.2f%%", snapshot.SharingRatio*100)
	t.Logf("  - Temps total de construction: %v", snapshot.TotalBuildTime)
	t.Logf("  - Temps moyen de construction: %v", snapshot.AverageBuildTime)
	t.Logf("  - Temps total de calcul de hash: %v", snapshot.TotalHashComputeTime)
	t.Logf("  - Cache hash - efficacité: %.2f%% (hits: %d, misses: %d, taille: %d)",
		metrics.GetHashCacheEfficiency()*100,
		snapshot.HashCacheHits, snapshot.HashCacheMisses, snapshot.HashCacheSize)
	t.Logf("  - Cache connexion - efficacité: %.2f%% (hits: %d, misses: %d)",
		metrics.GetConnectionCacheEfficiency()*100,
		snapshot.ConnectionCacheHits, snapshot.ConnectionCacheMisses)
	// Assertions
	if snapshot.TotalChainsBuilt != 1000 {
		t.Errorf("Attendu 1000 chaînes, obtenu %d", snapshot.TotalChainsBuilt)
	}
	// Avec des règles variées, le partage devrait être plus faible mais toujours présent
	if snapshot.SharingRatio < 0.1 {
		t.Errorf("Ratio de partage trop faible: %.2f%% (attendu >= 10%%)", snapshot.SharingRatio*100)
	}
	// Performance: construction moyenne < 1ms par règle
	avgBuildTimeMs := float64(snapshot.AverageBuildTime.Nanoseconds()) / 1e6
	if avgBuildTimeMs > 1.0 {
		t.Errorf("Temps moyen de construction trop élevé: %.2fms (attendu < 1ms)", avgBuildTimeMs)
	}
}

// TestMetrics_Accurate vérifie la précision des métriques
func TestMetrics_Accurate(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Reset pour commencer avec des métriques propres
	network.ResetChainMetrics()
	// Construire 10 chaînes simples
	for i := 0; i < 10; i++ {
		ruleID := fmt.Sprintf("test_rule_%d", i)
		conditions := []SimpleCondition{
			{
				Type:     "binaryOperation",
				Left:     map[string]interface{}{"type": "variable", "name": "x"},
				Operator: "==",
				Right:    map[string]interface{}{"type": "literal", "value": float64(i % 3)},
			},
		}
		builder := NewAlphaChainBuilderWithMetrics(network, storage, network.ChainMetrics)
		_, err := builder.BuildChain(conditions, "obj", network.RootNode, ruleID)
		if err != nil {
			t.Fatalf("Erreur: %v", err)
		}
	}
	metrics := network.GetChainMetrics()
	snapshot := metrics.GetSnapshot()
	// Vérifications précises
	if snapshot.TotalChainsBuilt != 10 {
		t.Errorf("TotalChainsBuilt incorrect: attendu 10, obtenu %d", snapshot.TotalChainsBuilt)
	}
	totalNodes := snapshot.TotalNodesCreated + snapshot.TotalNodesReused
	if totalNodes != 10 {
		t.Errorf("Total de nœuds incorrect: attendu 10, obtenu %d", totalNodes)
	}
	// Vérifier que les détails de chaîne sont enregistrés
	if len(snapshot.ChainDetails) != 10 {
		t.Errorf("Nombre de détails de chaîne incorrect: attendu 10, obtenu %d", len(snapshot.ChainDetails))
	}
	// Vérifier que chaque détail a des informations valides
	for i, detail := range snapshot.ChainDetails {
		if detail.RuleID == "" {
			t.Errorf("Détail %d: RuleID vide", i)
		}
		if detail.ChainLength != 1 {
			t.Errorf("Détail %d: ChainLength incorrect: attendu 1, obtenu %d", i, detail.ChainLength)
		}
		if detail.NodesCreated+detail.NodesReused != 1 {
			t.Errorf("Détail %d: Total nœuds incorrect", i)
		}
	}
	// Vérifier les ratios
	expectedAvgLength := float64(totalNodes) / float64(snapshot.TotalChainsBuilt)
	if snapshot.AverageChainLength != expectedAvgLength {
		t.Errorf("AverageChainLength incorrect: attendu %.2f, obtenu %.2f",
			expectedAvgLength, snapshot.AverageChainLength)
	}
	expectedSharingRatio := float64(snapshot.TotalNodesReused) / float64(totalNodes)
	if snapshot.SharingRatio != expectedSharingRatio {
		t.Errorf("SharingRatio incorrect: attendu %.2f, obtenu %.2f",
			expectedSharingRatio, snapshot.SharingRatio)
	}
}

// TestMetrics_HashCache vérifie le fonctionnement du cache de hash
func TestMetrics_HashCache(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.ResetChainMetrics()
	// Créer deux règles avec des conditions identiques
	conditions := []SimpleCondition{
		{
			Type:     "binaryOperation",
			Left:     map[string]interface{}{"type": "variable", "name": "value"},
			Operator: ">",
			Right:    map[string]interface{}{"type": "literal", "value": 100.0},
		},
	}
	builder := NewAlphaChainBuilderWithMetrics(network, storage, network.ChainMetrics)
	// Première règle
	_, err := builder.BuildChain(conditions, "item", network.RootNode, "rule_1")
	if err != nil {
		t.Fatalf("Erreur règle 1: %v", err)
	}
	metrics := network.GetChainMetrics()
	initialMisses := metrics.GetSnapshot().HashCacheMisses
	// Deuxième règle (conditions identiques)
	_, err = builder.BuildChain(conditions, "item", network.RootNode, "rule_2")
	if err != nil {
		t.Fatalf("Erreur règle 2: %v", err)
	}
	snapshot := metrics.GetSnapshot()
	t.Logf("Cache hash - hits: %d, misses: %d, taille: %d",
		snapshot.HashCacheHits, snapshot.HashCacheMisses, snapshot.HashCacheSize)
	// La deuxième règle devrait avoir un hit de cache
	if snapshot.HashCacheHits == 0 {
		t.Error("Attendu au moins un hit de cache pour la deuxième règle identique")
	}
	// Il ne devrait pas y avoir de nouveau miss pour la deuxième règle
	if snapshot.HashCacheMisses > initialMisses {
		t.Errorf("Nouveau miss de cache détecté: initial=%d, final=%d",
			initialMisses, snapshot.HashCacheMisses)
	}
	// Vérifier l'efficacité du cache
	efficiency := metrics.GetHashCacheEfficiency()
	t.Logf("Efficacité du cache de hash: %.2f%%", efficiency*100)
	if efficiency < 0.3 {
		t.Errorf("Efficacité du cache trop faible: %.2f%% (attendu >= 30%%)", efficiency*100)
	}
}

// TestMetrics_GetSummary vérifie le format du résumé des métriques
func TestMetrics_GetSummary(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.ResetChainMetrics()
	// Construire quelques chaînes
	for i := 0; i < 5; i++ {
		conditions := []SimpleCondition{
			{
				Type:     "binaryOperation",
				Left:     map[string]interface{}{"type": "variable", "name": "n"},
				Operator: "==",
				Right:    map[string]interface{}{"type": "literal", "value": float64(i)},
			},
		}
		builder := NewAlphaChainBuilderWithMetrics(network, storage, network.ChainMetrics)
		_, err := builder.BuildChain(conditions, "test", network.RootNode, fmt.Sprintf("r%d", i))
		if err != nil {
			t.Fatalf("Erreur: %v", err)
		}
	}
	metrics := network.GetChainMetrics()
	summary := metrics.GetSummary()
	// Vérifier la structure du résumé
	if summary == nil {
		t.Fatal("Résumé nil")
	}
	// Vérifier les sections
	chains, ok := summary["chains"].(map[string]interface{})
	if !ok {
		t.Fatal("Section 'chains' manquante ou mal formatée")
	}
	nodes, ok := summary["nodes"].(map[string]interface{})
	if !ok {
		t.Fatal("Section 'nodes' manquante ou mal formatée")
	}
	_, ok = summary["hash_cache"].(map[string]interface{})
	if !ok {
		t.Fatal("Section 'hash_cache' manquante ou mal formatée")
	}
	_, ok = summary["connection_cache"].(map[string]interface{})
	if !ok {
		t.Fatal("Section 'connection_cache' manquante ou mal formatée")
	}
	// Vérifier quelques valeurs
	if totalBuilt, ok := chains["total_built"].(int); !ok || totalBuilt != 5 {
		t.Errorf("total_built incorrect: %v", chains["total_built"])
	}
	if sharingRatio, ok := nodes["sharing_ratio"].(float64); !ok || sharingRatio < 0 || sharingRatio > 1 {
		t.Errorf("sharing_ratio invalide: %v", nodes["sharing_ratio"])
	}
	t.Logf("Résumé des métriques: %+v", summary)
}

// TestMetrics_TopChains vérifie les fonctions de classement des chaînes
func TestMetrics_TopChains(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.ResetChainMetrics()
	// Créer des chaînes de différentes longueurs
	for i := 1; i <= 5; i++ {
		conditions := make([]SimpleCondition, i) // i conditions
		for j := 0; j < i; j++ {
			conditions[j] = SimpleCondition{
				Type:     "binaryOperation",
				Left:     map[string]interface{}{"type": "variable", "name": fmt.Sprintf("f%d", j)},
				Operator: "==",
				Right:    map[string]interface{}{"type": "literal", "value": float64(j)},
			}
		}
		builder := NewAlphaChainBuilderWithMetrics(network, storage, network.ChainMetrics)
		_, err := builder.BuildChain(conditions, "obj", network.RootNode, fmt.Sprintf("chain_%d", i))
		if err != nil {
			t.Fatalf("Erreur: %v", err)
		}
	}
	metrics := network.GetChainMetrics()
	// Obtenir les 3 chaînes les plus longues
	topByLength := metrics.GetTopChainsByLength(3)
	if len(topByLength) != 3 {
		t.Errorf("Attendu 3 chaînes, obtenu %d", len(topByLength))
	}
	// Vérifier que c'est trié par longueur décroissante
	for i := 0; i < len(topByLength)-1; i++ {
		if topByLength[i].ChainLength < topByLength[i+1].ChainLength {
			t.Errorf("Tri incorrect: chaîne %d (len=%d) avant chaîne %d (len=%d)",
				i, topByLength[i].ChainLength, i+1, topByLength[i+1].ChainLength)
		}
	}
	t.Logf("Top 3 chaînes par longueur:")
	for i, chain := range topByLength {
		t.Logf("  %d. %s - longueur: %d", i+1, chain.RuleID, chain.ChainLength)
	}
	// Obtenir les 3 chaînes avec le temps de construction le plus long
	topByTime := metrics.GetTopChainsByBuildTime(3)
	if len(topByTime) != 3 {
		t.Errorf("Attendu 3 chaînes, obtenu %d", len(topByTime))
	}
	// Vérifier que c'est trié par temps décroissant
	for i := 0; i < len(topByTime)-1; i++ {
		if topByTime[i].BuildTime < topByTime[i+1].BuildTime {
			t.Errorf("Tri incorrect par temps")
		}
	}
	t.Logf("Top 3 chaînes par temps de construction:")
	for i, chain := range topByTime {
		t.Logf("  %d. %s - temps: %v", i+1, chain.RuleID, chain.BuildTime)
	}
}

// Fonction helper pour sélectionner un opérateur basé sur les indices
func selectOperator(i, j int) string {
	operators := []string{"==", "!=", ">", "<", ">=", "<="}
	return operators[(i+j)%len(operators)]
}

// Benchmarks
// BenchmarkChainBuild_SimilarRules benchmark la construction avec règles similaires
func BenchmarkChainBuild_SimilarRules(b *testing.B) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	conditions := []SimpleCondition{
		{
			Type:     "binaryOperation",
			Left:     map[string]interface{}{"type": "variable", "name": "age"},
			Operator: ">",
			Right:    map[string]interface{}{"type": "literal", "value": 18.0},
		},
		{
			Type:     "binaryOperation",
			Left:     map[string]interface{}{"type": "variable", "name": "status"},
			Operator: "==",
			Right:    map[string]interface{}{"type": "literal", "value": "active"},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ruleID := fmt.Sprintf("bench_rule_%d", i)
		builder := NewAlphaChainBuilderWithMetrics(network, storage, network.ChainMetrics)
		_, err := builder.BuildChain(conditions, "person", network.RootNode, ruleID)
		if err != nil {
			b.Fatalf("Erreur: %v", err)
		}
	}
}

// BenchmarkChainBuild_VariedRules benchmark la construction avec règles variées
func BenchmarkChainBuild_VariedRules(b *testing.B) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ruleID := fmt.Sprintf("bench_rule_%d", i)
		// Varier les conditions
		conditions := []SimpleCondition{
			{
				Type:     "binaryOperation",
				Left:     map[string]interface{}{"type": "variable", "name": fmt.Sprintf("field%d", i%10)},
				Operator: selectOperator(i, 0),
				Right:    map[string]interface{}{"type": "literal", "value": float64(i % 100)},
			},
		}
		builder := NewAlphaChainBuilderWithMetrics(network, storage, network.ChainMetrics)
		_, err := builder.BuildChain(conditions, "entity", network.RootNode, ruleID)
		if err != nil {
			b.Fatalf("Erreur: %v", err)
		}
	}
}

// BenchmarkHashCompute benchmark le calcul de hash
func BenchmarkHashCompute(b *testing.B) {
	condition := map[string]interface{}{
		"type":     "binaryOperation",
		"left":     map[string]interface{}{"type": "variable", "name": "x"},
		"operator": "==",
		"right":    map[string]interface{}{"type": "literal", "value": 42.0},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ConditionHash(condition, "var")
		if err != nil {
			b.Fatalf("Erreur: %v", err)
		}
	}
}

// BenchmarkHashComputeCached benchmark le calcul de hash avec cache
func BenchmarkHashComputeCached(b *testing.B) {
	registry := NewAlphaSharingRegistry()
	condition := map[string]interface{}{
		"type":     "binaryOperation",
		"left":     map[string]interface{}{"type": "variable", "name": "x"},
		"operator": "==",
		"right":    map[string]interface{}{"type": "literal", "value": 42.0},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := registry.ConditionHashCached(condition, "var")
		if err != nil {
			b.Fatalf("Erreur: %v", err)
		}
	}
}
