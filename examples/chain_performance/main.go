//go:build ignore
// +build ignore

// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"time"

	"github.com/treivax/tsd/rete"
)

func main() {
	fmt.Println("🚀 Exemple d'Optimisations de Performance des Chaînes Alpha")
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println()

	// Exemple 1: Construction basique avec métriques
	example1_BasicMetrics()

	fmt.Println()

	// Exemple 2: Comparaison règles similaires vs variées
	example2_ComparePatterns()

	fmt.Println()

	// Exemple 3: Analyse des chaînes
	example3_AnalyzeChains()

	fmt.Println()

	// Exemple 4: Monitoring continu
	example4_ContinuousMonitoring()
}

// Exemple 1: Construction basique avec métriques
func example1_BasicMetrics() {
	fmt.Println("📊 Exemple 1: Construction Basique avec Métriques")
	fmt.Println("-" + string(make([]byte, 60)))

	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)

	// Construire 50 règles similaires
	startTime := time.Now()
	for i := 0; i < 50; i++ {
		conditions := []rete.SimpleCondition{
			{
				Type:     "binaryOperation",
				Left:     map[string]interface{}{"type": "variable", "name": "age"},
				Operator: ">",
				Right:    map[string]interface{}{"type": "literal", "value": float64(18)},
			},
			{
				Type:     "binaryOperation",
				Left:     map[string]interface{}{"type": "variable", "name": "status"},
				Operator: "==",
				Right:    map[string]interface{}{"type": "literal", "value": "active"},
			},
		}

		builder := rete.NewAlphaChainBuilderWithMetrics(network, storage, network.ChainMetrics)
		ruleID := fmt.Sprintf("rule_%d", i)
		_, err := builder.BuildChain(conditions, "person", network.RootNode, ruleID)
		if err != nil {
			panic(err)
		}
	}
	elapsed := time.Since(startTime)

	// Afficher les métriques
	metrics := network.GetChainMetrics()
	summary := metrics.GetSummary()

	chains := summary["chains"].(map[string]interface{})
	nodes := summary["nodes"].(map[string]interface{})
	hashCache := summary["hash_cache"].(map[string]interface{})

	fmt.Printf("  Temps total:            %v\n", elapsed)
	fmt.Printf("  Chaînes construites:    %d\n", chains["total_built"])
	fmt.Printf("  Nœuds créés:            %d\n", nodes["total_created"])
	fmt.Printf("  Nœuds réutilisés:       %d\n", nodes["total_reused"])
	fmt.Printf("  Ratio de partage:       %.2f%%\n", nodes["reuse_rate_pct"])
	fmt.Printf("  Efficacité cache hash:  %.2f%%\n", hashCache["efficiency_pct"])
	fmt.Printf("  Temps moyen/règle:      %s\n", chains["average_build_time"])
}

// Exemple 2: Comparaison règles similaires vs variées
func example2_ComparePatterns() {
	fmt.Println("🔄 Exemple 2: Comparaison de Patterns")
	fmt.Println("-" + string(make([]byte, 60)))

	// Test avec règles similaires
	fmt.Println("\n  🔹 Test 1: 100 règles similaires")
	testPattern(100, "similar")

	// Test avec règles variées
	fmt.Println("\n  🔹 Test 2: 100 règles variées")
	testPattern(100, "varied")
}

func testPattern(count int, pattern string) {
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)
	builder := rete.NewAlphaChainBuilderWithMetrics(network, storage, network.ChainMetrics)

	startTime := time.Now()
	for i := 0; i < count; i++ {
		var conditions []rete.SimpleCondition

		if pattern == "similar" {
			// Règles similaires (favorise le partage)
			conditions = []rete.SimpleCondition{
				{
					Type:     "binaryOperation",
					Left:     map[string]interface{}{"type": "variable", "name": "value"},
					Operator: ">",
					Right:    map[string]interface{}{"type": "literal", "value": float64(i % 5)},
				},
			}
		} else {
			// Règles variées (moins de partage)
			conditions = []rete.SimpleCondition{
				{
					Type:     "binaryOperation",
					Left:     map[string]interface{}{"type": "variable", "name": fmt.Sprintf("field%d", i)},
					Operator: selectOperator(i),
					Right:    map[string]interface{}{"type": "literal", "value": float64(i)},
				},
			}
		}

		ruleID := fmt.Sprintf("rule_%d", i)
		_, err := builder.BuildChain(conditions, "entity", network.RootNode, ruleID)
		if err != nil {
			panic(err)
		}
	}
	elapsed := time.Since(startTime)

	metrics := network.GetChainMetrics()
	snapshot := metrics.GetSnapshot()

	fmt.Printf("     Temps total:         %v\n", elapsed)
	fmt.Printf("     Nœuds créés:         %d\n", snapshot.TotalNodesCreated)
	fmt.Printf("     Nœuds réutilisés:    %d\n", snapshot.TotalNodesReused)
	fmt.Printf("     Ratio de partage:    %.2f%%\n", snapshot.SharingRatio*100)
	fmt.Printf("     Cache hash:          %.2f%%\n", metrics.GetHashCacheEfficiency()*100)
}

// Exemple 3: Analyse des chaînes
func example3_AnalyzeChains() {
	fmt.Println("🔍 Exemple 3: Analyse Détaillée des Chaînes")
	fmt.Println("-" + string(make([]byte, 60)))

	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)
	builder := rete.NewAlphaChainBuilderWithMetrics(network, storage, network.ChainMetrics)

	// Créer des chaînes de différentes tailles
	for i := 1; i <= 5; i++ {
		conditions := make([]rete.SimpleCondition, i)
		for j := 0; j < i; j++ {
			conditions[j] = rete.SimpleCondition{
				Type:     "binaryOperation",
				Left:     map[string]interface{}{"type": "variable", "name": fmt.Sprintf("field%d", j)},
				Operator: "==",
				Right:    map[string]interface{}{"type": "literal", "value": float64(j)},
			}
		}

		ruleID := fmt.Sprintf("chain_%d_nodes", i)
		_, err := builder.BuildChain(conditions, "test", network.RootNode, ruleID)
		if err != nil {
			panic(err)
		}
	}

	metrics := network.GetChainMetrics()

	// Top chaînes par longueur
	fmt.Println("\n  📏 Top 3 des chaînes les plus longues:")
	topLong := metrics.GetTopChainsByLength(3)
	for i, chain := range topLong {
		fmt.Printf("     %d. %s - %d nœuds (%d créés, %d réutilisés)\n",
			i+1, chain.RuleID, chain.ChainLength,
			chain.NodesCreated, chain.NodesReused)
	}

	// Top chaînes par temps de construction
	fmt.Println("\n  ⏱️  Top 3 des chaînes les plus lentes:")
	topSlow := metrics.GetTopChainsByBuildTime(3)
	for i, chain := range topSlow {
		fmt.Printf("     %d. %s - %v (longueur: %d)\n",
			i+1, chain.RuleID, chain.BuildTime, chain.ChainLength)
	}
}

// Exemple 4: Monitoring continu
func example4_ContinuousMonitoring() {
	fmt.Println("📡 Exemple 4: Monitoring Continu (simulation)")
	fmt.Println("-" + string(make([]byte, 60)))

	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)
	builder := rete.NewAlphaChainBuilderWithMetrics(network, storage, network.ChainMetrics)

	// Simuler la construction de règles en plusieurs vagues
	for wave := 1; wave <= 3; wave++ {
		fmt.Printf("\n  🌊 Vague %d - Construction de 20 règles\n", wave)

		for i := 0; i < 20; i++ {
			conditions := []rete.SimpleCondition{
				{
					Type:     "binaryOperation",
					Left:     map[string]interface{}{"type": "variable", "name": "x"},
					Operator: "==",
					Right:    map[string]interface{}{"type": "literal", "value": float64(i % 10)},
				},
			}

			ruleID := fmt.Sprintf("wave%d_rule%d", wave, i)
			_, err := builder.BuildChain(conditions, "obj", network.RootNode, ruleID)
			if err != nil {
				panic(err)
			}
		}

		// Afficher les métriques après chaque vague
		metrics := network.GetChainMetrics()
		snapshot := metrics.GetSnapshot()

		fmt.Printf("     Total chaînes:       %d\n", snapshot.TotalChainsBuilt)
		fmt.Printf("     Partage:             %.2f%%\n", snapshot.SharingRatio*100)
		fmt.Printf("     Cache hash:          %.2f%% (%d entrées)\n",
			metrics.GetHashCacheEfficiency()*100, snapshot.HashCacheSize)

		// Simuler un délai
		time.Sleep(100 * time.Millisecond)
	}

	// Résumé final
	fmt.Println("\n  📊 Résumé Final:")
	summary := network.GetChainMetrics().GetSummary()
	chains := summary["chains"].(map[string]interface{})
	nodes := summary["nodes"].(map[string]interface{})

	fmt.Printf("     Total chaînes:       %d\n", chains["total_built"])
	fmt.Printf("     Total nœuds:         %d (créés) + %d (réutilisés)\n",
		nodes["total_created"], nodes["total_reused"])
	fmt.Printf("     Ratio de partage:    %.2f%%\n", nodes["reuse_rate_pct"])
}

// Helper: sélectionner un opérateur basé sur l'index
func selectOperator(i int) string {
	operators := []string{"==", "!=", ">", "<", ">=", "<="}
	return operators[i%len(operators)]
}
