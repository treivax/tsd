// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta_test

import (
	"testing"

	"github.com/treivax/tsd/rete/delta"
)

// TestIndexation_IntegrationScenario teste un scénario complet d'indexation.
//
// Ce test montre comment construire un index de dépendances depuis zéro,
// ajouter différents types de nœuds, et interroger l'index pour trouver
// les nœuds affectés par des changements.
func TestIndexation_IntegrationScenario(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION - Scénario d'indexation complet")
	t.Log("=" + "===========================================")

	// Étape 1 : Créer l'index et le builder
	idx, builder := setupIndexAndBuilder(t)

	// Étape 2 : Ajouter des nœuds alpha
	addAlphaNodes(t, idx, builder)

	// Étape 3 : Ajouter un nœud beta
	addBetaNode(t, idx, builder)

	// Étape 4 : Ajouter un nœud terminal
	addTerminalNode(t, idx, builder)

	// Étape 5 : Vérifier les statistiques
	verifyIndexStats(t, idx)

	// Étape 6 : Tester les requêtes d'affectation
	testAffectedNodeQueries(t, idx)

	// Étape 7 : Tester avec un FactDelta
	testFactDeltaQuery(t, idx)

	// Étape 8 : Vérifier les diagnostics
	verifyBuilderDiagnostics(t, builder)

	// Étape 9 : Test de Clear
	testClearIndex(t, idx)

	t.Log("\n🎉 TEST COMPLET RÉUSSI - Tous les scénarios validés!")
}

// setupIndexAndBuilder crée et initialise l'index et le builder.
func setupIndexAndBuilder(t *testing.T) (*delta.DependencyIndex, *delta.IndexBuilder) {
	t.Log("\n📋 Étape 1 : Création de l'index de dépendances")
	idx := delta.NewDependencyIndex()
	if idx == nil {
		t.Fatal("❌ Échec de création de l'index")
	}
	t.Log("✅ Index créé avec succès")

	t.Log("\n📋 Étape 2 : Création du builder avec diagnostics")
	builder := delta.NewIndexBuilder()
	builder.EnableDiagnostics()
	t.Log("✅ Builder créé et diagnostics activés")

	return idx, builder
}

// addAlphaNodes ajoute des nœuds alpha au builder.
func addAlphaNodes(t *testing.T, idx *delta.DependencyIndex, builder *delta.IndexBuilder) {
	t.Log("\n📋 Étape 3 : Indexation de nœuds alpha")

	// Alpha 1 : Product.price > 100
	alphaCondition1 := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "price",
		},
		"right": 100,
	}

	err := builder.BuildFromAlphaNode(idx, "alpha_price", "Product", alphaCondition1)
	if err != nil {
		t.Fatalf("❌ Erreur alpha node 1: %v", err)
	}
	t.Log("  ✅ Alpha node 'alpha_price' indexé (champ: price)")

	// Alpha 2 : Product.status == "active" && Product.stock > 0
	alphaCondition2 := map[string]interface{}{
		"type": "binaryOp",
		"left": map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":  "fieldAccess",
				"field": "status",
			},
			"right": "active",
		},
		"right": map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":  "fieldAccess",
				"field": "stock",
			},
			"right": 0,
		},
	}

	err = builder.BuildFromAlphaNode(idx, "alpha_status_stock", "Product", alphaCondition2)
	if err != nil {
		t.Fatalf("❌ Erreur alpha node 2: %v", err)
	}
	t.Log("  ✅ Alpha node 'alpha_status_stock' indexé (champs: status, stock)")
}

// addBetaNode ajoute un nœud beta au builder.
func addBetaNode(t *testing.T, idx *delta.DependencyIndex, builder *delta.IndexBuilder) {
	t.Log("\n📋 Étape 4 : Indexation d'un nœud beta")

	betaCondition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "product_id",
		},
		"right": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "id",
		},
	}

	err := builder.BuildFromBetaNode(idx, "beta_order_product", "Order", betaCondition)
	if err != nil {
		t.Fatalf("❌ Erreur beta node: %v", err)
	}
	t.Log("  ✅ Beta node 'beta_order_product' indexé (champs: product_id, id)")
}

// addTerminalNode ajoute un nœud terminal au builder.
func addTerminalNode(t *testing.T, idx *delta.DependencyIndex, builder *delta.IndexBuilder) {
	t.Log("\n📋 Étape 5 : Indexation d'un nœud terminal")

	actions := []interface{}{
		map[string]interface{}{
			"type": "updateWithModifications",
			"modifications": map[string]interface{}{
				"price":  150,
				"status": "updated",
			},
		},
	}

	err := builder.BuildFromTerminalNode(idx, "terminal_update_product", "Product", actions)
	if err != nil {
		t.Fatalf("❌ Erreur terminal node: %v", err)
	}
	t.Log("  ✅ Terminal node 'terminal_update_product' indexé (champs: price, status)")
}

// verifyIndexStats vérifie les statistiques de l'index.
func verifyIndexStats(t *testing.T, idx *delta.DependencyIndex) {
	t.Log("\n📋 Étape 6 : Vérification des statistiques")

	stats := idx.GetStats()
	t.Logf("  📊 Nœuds indexés : %d", stats.NodeCount)
	t.Logf("  📊 Entrées de champs : %d", stats.FieldCount)
	t.Logf("  📊 Nœuds alpha : %d", stats.AlphaNodeCount)
	t.Logf("  📊 Nœuds beta : %d", stats.BetaNodeCount)
	t.Logf("  📊 Nœuds terminaux : %d", stats.TerminalCount)
	t.Logf("  📊 Types de faits : %v", stats.FactTypes)
	t.Logf("  📊 Estimation mémoire : %d bytes", stats.MemoryEstimate)

	if stats.NodeCount != 4 {
		t.Errorf("❌ Attendu 4 nœuds, obtenu %d", stats.NodeCount)
	}
	if stats.AlphaNodeCount != 2 {
		t.Errorf("❌ Attendu 2 nœuds alpha, obtenu %d", stats.AlphaNodeCount)
	}
	if stats.BetaNodeCount != 1 {
		t.Errorf("❌ Attendu 1 nœud beta, obtenu %d", stats.BetaNodeCount)
	}
	if stats.TerminalCount != 1 {
		t.Errorf("❌ Attendu 1 nœud terminal, obtenu %d", stats.TerminalCount)
	}
}

// testAffectedNodeQueries teste les requêtes de nœuds affectés.
func testAffectedNodeQueries(t *testing.T, idx *delta.DependencyIndex) {
	t.Log("\n📋 Étape 7 : Requêtes de nœuds affectés")

	// Requête 1 : Qui est affecté par Product.price ?
	testSingleFieldQuery(t, idx, "Product", "price", 2)

	// Requête 2 : Qui est affecté par Product.status ?
	testSingleFieldQuery(t, idx, "Product", "status", 2)

	// Requête 3 : Qui est affecté par Order.product_id ?
	testSingleFieldQuery(t, idx, "Order", "product_id", 1)
}

// testSingleFieldQuery teste une requête pour un champ unique.
func testSingleFieldQuery(t *testing.T, idx *delta.DependencyIndex, factType, field string, expectedCount int) {
	t.Logf("\n  🔍 Requête : Qui est affecté par %s.%s ?", factType, field)
	affected := idx.GetAffectedNodes(factType, field)
	t.Logf("    Nœuds trouvés : %d", len(affected))
	for _, node := range affected {
		t.Logf("    - %s", node.String())
	}

	if len(affected) != expectedCount {
		t.Errorf("❌ Attendu %d nœuds affectés par %s, obtenu %d", expectedCount, field, len(affected))
	}
}

// testFactDeltaQuery teste une requête avec un FactDelta.
func testFactDeltaQuery(t *testing.T, idx *delta.DependencyIndex) {
	t.Log("\n📋 Étape 8 : Requête avec FactDelta")

	factDelta := delta.NewFactDelta("Product~p123", "Product")
	factDelta.AddFieldChange("price", 100.0, 150.0)
	factDelta.AddFieldChange("status", "active", "inactive")

	t.Log("  📝 Delta créé :")
	t.Logf("    - Fact ID : %s", factDelta.FactID)
	t.Logf("    - Type : %s", factDelta.FactType)
	t.Logf("    - Champs modifiés : %d", len(factDelta.Fields))

	affectedByDelta := idx.GetAffectedNodesForDelta(factDelta)
	t.Logf("\n  🔍 Nœuds affectés par le delta : %d", len(affectedByDelta))
	for _, node := range affectedByDelta {
		t.Logf("    - %s", node.String())
	}

	// On devrait avoir alpha_price, alpha_status_stock, terminal_update_product
	if len(affectedByDelta) != 3 {
		t.Errorf("❌ Attendu 3 nœuds affectés par delta, obtenu %d", len(affectedByDelta))
	}
}

// verifyBuilderDiagnostics vérifie les diagnostics du builder.
func verifyBuilderDiagnostics(t *testing.T, builder *delta.IndexBuilder) {
	t.Log("\n📋 Étape 9 : Diagnostics du builder")

	diag := builder.GetDiagnostics()
	t.Logf("  📊 Nœuds traités : %d", diag.NodesProcessed)
	t.Logf("  📊 Nœuds ignorés : %d", diag.NodesSkipped)
	t.Logf("  📊 Champs extraits : %d", diag.FieldsExtracted)
	t.Logf("  📊 Erreurs : %d", len(diag.Errors))
	t.Logf("  📊 Avertissements : %d", len(diag.Warnings))

	if diag.NodesProcessed != 4 {
		t.Errorf("❌ Attendu 4 nœuds traités, obtenu %d", diag.NodesProcessed)
	}
}

// testClearIndex teste le vidage de l'index.
func testClearIndex(t *testing.T, idx *delta.DependencyIndex) {
	t.Log("\n📋 Étape 10 : Test de vidage de l'index")

	idx.Clear()
	statsAfterClear := idx.GetStats()
	t.Logf("  📊 Nœuds après clear : %d", statsAfterClear.NodeCount)

	if statsAfterClear.NodeCount != 0 {
		t.Errorf("❌ L'index devrait être vide après Clear, mais contient %d nœuds", statsAfterClear.NodeCount)
	}

	t.Log("  ✅ Index vidé avec succès")
}

// TestIndexation_Performance teste les performances de base de l'indexation.
func TestIndexation_Performance(t *testing.T) {
	t.Log("🧪 TEST PERFORMANCE - Indexation à grande échelle")

	idx := delta.NewDependencyIndex()
	builder := delta.NewIndexBuilder()

	// Ajouter 100 nœuds alpha
	numNodes := 100
	condition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "price",
		},
		"right": 100,
	}

	for i := 0; i < numNodes; i++ {
		nodeID := "alpha_" + string(rune('0'+i%10)) + string(rune('0'+i/10))
		err := builder.BuildFromAlphaNode(idx, nodeID, "Product", condition)
		if err != nil {
			t.Fatalf("❌ Erreur lors de l'ajout du nœud %d: %v", i, err)
		}
	}

	stats := idx.GetStats()
	t.Logf("📊 %d nœuds indexés", stats.NodeCount)
	t.Logf("📊 Estimation mémoire : %d bytes (%.2f KB)", stats.MemoryEstimate, float64(stats.MemoryEstimate)/1024.0)

	// Test de recherche
	affected := idx.GetAffectedNodes("Product", "price")
	t.Logf("🔍 Nœuds trouvés pour Product.price : %d", len(affected))

	if len(affected) != numNodes {
		t.Errorf("❌ Attendu %d nœuds, obtenu %d", numNodes, len(affected))
	}

	t.Log("✅ Test de performance réussi")
}
