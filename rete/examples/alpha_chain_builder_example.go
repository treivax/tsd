//go:build ignore
// +build ignore

// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"log"

	"github.com/treivax/tsd/rete"
)

func main() {
	fmt.Println("=== Alpha Chain Builder - Exemple d'utilisation ===")
	fmt.Println()

	// 1. Initialiser le réseau RETE
	fmt.Println("1. Initialisation du réseau RETE")
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)
	fmt.Println("   ✓ Réseau RETE créé")

	// 2. Définir le type de données Person
	fmt.Println("\n2. Définition du type Person")
	typeDef := rete.TypeDefinition{
		Type: "type",
		Name: "Person",
		Fields: []rete.Field{
			{Name: "id", Type: "string"},
			{Name: "age", Type: "number"},
			{Name: "name", Type: "string"},
			{Name: "city", Type: "string"},
			{Name: "salary", Type: "number"},
		},
	}
	parentNode := rete.NewTypeNode("person", typeDef, storage)
	network.TypeNodes["Person"] = parentNode
	fmt.Println("   ✓ Type Person défini avec 5 champs")

	// 3. Créer le builder
	fmt.Println("\n3. Création de l'Alpha Chain Builder")
	builder := rete.NewAlphaChainBuilder(network, storage)
	fmt.Println("   ✓ Builder créé")

	// 4. Définir les règles avec conditions
	fmt.Println("\n4. Définition des règles")

	// Règle 1: Personnes majeures nommées Alice
	rule1Conditions := []rete.SimpleCondition{
		rete.NewSimpleCondition("comparison", "p.age", ">", 18),
		rete.NewSimpleCondition("comparison", "p.name", "==", "Alice"),
	}
	fmt.Println("   Rule1: age > 18 AND name == 'Alice'")

	// Règle 2: Personnes majeures à Paris
	rule2Conditions := []rete.SimpleCondition{
		rete.NewSimpleCondition("comparison", "p.age", ">", 18),
		rete.NewSimpleCondition("comparison", "p.city", "==", "Paris"),
	}
	fmt.Println("   Rule2: age > 18 AND city == 'Paris'")

	// Règle 3: Identique à Rule1 (test de réutilisation complète)
	rule3Conditions := []rete.SimpleCondition{
		rete.NewSimpleCondition("comparison", "p.age", ">", 18),
		rete.NewSimpleCondition("comparison", "p.name", "==", "Alice"),
	}
	fmt.Println("   Rule3: age > 18 AND name == 'Alice' (identique à Rule1)")

	// Règle 4: Personnes bien payées
	rule4Conditions := []rete.SimpleCondition{
		rete.NewSimpleCondition("comparison", "p.salary", ">", 50000),
	}
	fmt.Println("   Rule4: salary > 50000 (aucun partage)")

	// 5. Construire les chaînes
	fmt.Println("\n5. Construction des chaînes alpha")
	fmt.Println("   (Regardez les logs détaillés ci-dessous)")
	fmt.Println()

	chain1, err := builder.BuildChain(rule1Conditions, "p", parentNode, "rule1")
	if err != nil {
		log.Fatalf("Erreur Rule1: %v", err)
	}

	chain2, err := builder.BuildChain(rule2Conditions, "p", parentNode, "rule2")
	if err != nil {
		log.Fatalf("Erreur Rule2: %v", err)
	}

	chain3, err := builder.BuildChain(rule3Conditions, "p", parentNode, "rule3")
	if err != nil {
		log.Fatalf("Erreur Rule3: %v", err)
	}

	chain4, err := builder.BuildChain(rule4Conditions, "p", parentNode, "rule4")
	if err != nil {
		log.Fatalf("Erreur Rule4: %v", err)
	}

	// 6. Analyser les résultats
	fmt.Println("\n6. Analyse des résultats")
	fmt.Println()

	// Statistiques par règle
	fmt.Println("📊 Statistiques par règle:")
	fmt.Println()

	analyzeChain(builder, chain1, "Rule1")
	analyzeChain(builder, chain2, "Rule2")
	analyzeChain(builder, chain3, "Rule3")
	analyzeChain(builder, chain4, "Rule4")

	// Statistiques réseau globales
	fmt.Println("📈 Statistiques du réseau RETE:")
	fmt.Println()
	netStats := network.GetNetworkStats()
	fmt.Printf("   Total alpha nodes: %d\n", netStats["alpha_nodes"])
	fmt.Printf("   Shared alpha nodes: %d\n", netStats["sharing_total_shared_alpha_nodes"])
	fmt.Printf("   Average sharing ratio: %.2f\n", netStats["sharing_average_sharing_ratio"])
	fmt.Printf("   Lifecycle nodes: %d\n", netStats["lifecycle_total_nodes"])
	fmt.Printf("   Total references: %d\n", netStats["lifecycle_total_references"])
	fmt.Println()

	// Calcul de l'économie
	totalConditions := len(rule1Conditions) + len(rule2Conditions) +
		len(rule3Conditions) + len(rule4Conditions)
	actualNodes := netStats["alpha_nodes"].(int)
	savedNodes := totalConditions - actualNodes
	savingsPercent := float64(savedNodes) / float64(totalConditions) * 100

	fmt.Printf("💰 Économie de mémoire:\n")
	fmt.Printf("   Sans partage: %d nœuds\n", totalConditions)
	fmt.Printf("   Avec partage: %d nœuds\n", actualNodes)
	fmt.Printf("   Économisés: %d nœuds (%.1f%%)\n", savedNodes, savingsPercent)
	fmt.Println()

	// 7. Détails des nœuds partagés
	fmt.Println("🔗 Détails des nœuds partagés:")
	fmt.Println()

	sharedNodeIDs := network.AlphaSharingManager.ListSharedAlphaNodes()
	for i, nodeID := range sharedNodeIDs {
		details := network.AlphaSharingManager.GetSharedAlphaNodeDetails(nodeID)
		if details != nil {
			fmt.Printf("   %d. Node %s\n", i+1, nodeID)
			fmt.Printf("      Variable: %s\n", details["variable_name"])
			fmt.Printf("      Child count: %d\n", details["child_count"])

			lifecycle, exists := network.LifecycleManager.GetNodeLifecycle(nodeID)
			if exists {
				fmt.Printf("      Referenced by: %d rule(s)\n", lifecycle.GetRefCount())
				rules := lifecycle.GetRules()
				fmt.Printf("      Rules: %v\n", rules)
			}
			fmt.Println()
		}
	}

	// 8. Validation des chaînes
	fmt.Println("✓ Validation des chaînes:")
	fmt.Println()

	chains := []*rete.AlphaChain{chain1, chain2, chain3, chain4}
	chainNames := []string{"Rule1", "Rule2", "Rule3", "Rule4"}

	for i, chain := range chains {
		if err := chain.ValidateChain(); err != nil {
			fmt.Printf("   ✗ %s: INVALIDE - %v\n", chainNames[i], err)
		} else {
			fmt.Printf("   ✓ %s: VALIDE\n", chainNames[i])
		}
	}
	fmt.Println()

	fmt.Println("=== Exemple terminé avec succès ===")
}

// analyzeChain affiche les statistiques détaillées d'une chaîne
func analyzeChain(builder *rete.AlphaChainBuilder, chain *rete.AlphaChain, name string) {
	stats := builder.GetChainStats(chain)

	fmt.Printf("   %s:\n", name)
	fmt.Printf("      Total nodes: %d\n", stats["total_nodes"])
	fmt.Printf("      Shared nodes: %d\n", stats["shared_nodes"])
	fmt.Printf("      New nodes: %d\n", stats["new_nodes"])

	// Afficher les détails de chaque nœud
	nodeDetails := stats["node_details"].([]map[string]interface{})
	fmt.Printf("      Details:\n")
	for _, detail := range nodeDetails {
		shared := ""
		if detail["is_shared"].(bool) {
			shared = " (SHARED)"
		}
		final := ""
		if detail["is_final"].(bool) {
			final = " [FINAL]"
		}
		fmt.Printf("         [%d] %s - refs=%d%s%s\n",
			detail["index"],
			detail["node_id"],
			detail["ref_count"],
			shared,
			final)
	}
	fmt.Println()
}
