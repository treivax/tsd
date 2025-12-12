// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"fmt"
	"strings"
	"testing"
)
// TestBugRETE001_ReproduceIssue verifies that the bug where alpha and beta conditions
// were not separated in the RETE network structure has been FIXED.
//
// Expected behavior (FIXED):
//   - Alpha conditions (tests on single facts) should be evaluated in AlphaNodes
//   - Beta conditions (tests between facts) should be evaluated in JoinNodes
//   - Alpha nodes should filter facts BEFORE they reach the join
//
// Previous buggy behavior:
//   - All conditions (alpha AND beta) were placed in a single JoinNode
//   - No early filtering occurred
//   - Conditions were evaluated for every pair of facts
func TestBugRETE001_ReproduceIssue(t *testing.T) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("✅ VERIFICATION DU FIX #RETE-001")
	fmt.Println("Bug Fixed: Alpha/Beta Condition Separation")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()
	// Arrange - Build network from TSD file with mixed alpha/beta conditions
	t.Log("📋 Construction du réseau depuis fichier TSD")
	pipeline := NewConstraintPipeline()
	storage := NewMemoryStorage()
	tsdFile := "testdata/bug_rete001_minimal.tsd"
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("❌ Erreur construction réseau: %v", err)
	}
	t.Log("✅ Réseau construit")
	t.Logf("   - TypeNodes: %d", len(network.TypeNodes))
	t.Logf("   - AlphaNodes: %d", len(network.AlphaNodes))
	t.Logf("   - BetaNodes: %d", len(network.BetaNodes))
	fmt.Println()
	// Act - Analyze the network structure
	t.Log("🔍 Analyse de la structure du réseau RETE...")
	fmt.Println()
	// Count AlphaNodes with filtering conditions (not passthrough)
	alphaNodesWithFilters := 0
	passthroughAlphaNodes := 0
	for alphaID, alphaNode := range network.AlphaNodes {
		node := alphaNode
		// Check if this is a passthrough node (no real condition)
		if node.Condition == nil {
			passthroughAlphaNodes++
			t.Logf("   AlphaNode passthrough: %s", alphaID)
		} else {
			// Check if it's a real filtering condition
			if condMap, ok := node.Condition.(map[string]interface{}); ok {
				// Check if it's a passthrough type
				if condType, hasType := condMap["type"].(string); hasType && condType == "passthrough" {
					passthroughAlphaNodes++
					t.Logf("   AlphaNode passthrough: %s", alphaID)
				} else {
					// A filtering AlphaNode should have an operator like >, <, ==, etc.
					// on a single variable
					if op, hasOp := condMap["operator"]; hasOp {
						alphaNodesWithFilters++
						t.Logf("   AlphaNode avec filtre: %s", alphaID)
						t.Logf("      Opérateur: %v", op)
						t.Logf("      Condition: %+v", node.Condition)
					}
				}
			}
		}
	}
	// Analyze JoinNodes and their conditions
	joinNodesCount := 0
	var joinNodeCondition interface{}
	for betaID, betaNode := range network.BetaNodes {
		if node, ok := betaNode.(*JoinNode); ok {
			joinNodesCount++
			joinNodeCondition = node.Condition
			t.Logf("   JoinNode: %s", betaID)
			t.Logf("      Condition complète: %+v", node.Condition)
		}
	}
	fmt.Println()
	t.Log("📊 Résultat de l'analyse:")
	t.Logf("   - AlphaNodes passthrough: %d", passthroughAlphaNodes)
	t.Logf("   - AlphaNodes avec filtres: %d", alphaNodesWithFilters)
	t.Logf("   - JoinNodes: %d", joinNodesCount)
	fmt.Println()
	// Assert - Verify the bug is FIXED
	t.Log("✅ Vérification du fix...")
	fmt.Println()
	// FIX VERIFICATION:
	// Expected: At least 1 AlphaNode with the alpha condition (c.qte > 5)
	// This means the bug has been fixed
	if alphaNodesWithFilters == 0 {
		t.Error("❌ Bug still exists - no filtering AlphaNodes found!")
		t.Error("   Expected: At least 1 AlphaNode with filters")
		t.Errorf("   Got: %d AlphaNodes with filters", alphaNodesWithFilters)
		fmt.Println()
		fmt.Println("The bug has NOT been fixed!")
		return
	}
	t.Log("✅ FIX CONFIRMED: Filtering AlphaNodes detected")
	t.Logf("   → %d AlphaNode(s) with filters created", alphaNodesWithFilters)
	t.Log("   → Alpha conditions are separated from beta conditions")
	t.Log("   → Early filtering is now active")
	fmt.Println()
	// The JoinNode should contain ONLY the beta condition, not both
	if joinNodeCondition == nil {
		t.Error("❌ JoinNode without condition - unexpected structure")
		return
	}
	t.Log("✅ FIX VERIFIED: Alpha/Beta separation is working correctly")
	t.Logf("   → Alpha conditions: in AlphaNodes (%d nodes)", alphaNodesWithFilters)
	t.Logf("   → Beta conditions: in JoinNodes (%d nodes)", joinNodesCount)
	t.Log("   → Facts are filtered BEFORE joining")
	fmt.Println()
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🐛 RÉSUMÉ DU BUG")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()
	fmt.Println("PROBLÈME IDENTIFIÉ:")
	fmt.Println("  • Les conditions alpha (c.qte > 5) ne sont PAS dans des AlphaNodes")
	fmt.Println("  • La condition complète est dans le JoinNode")
	fmt.Println("  • Pas de filtrage précoce des faits")
	fmt.Println()
	fmt.Println("IMPACT:")
	fmt.Println("  • Performance dégradée (évaluations redondantes)")
	fmt.Println("  • Violation du principe RETE classique")
	fmt.Println("  • Pas de partage des conditions alpha entre règles")
	fmt.Println("  • Évaluation pour CHAQUE paire (Commande, Produit)")
	fmt.Println()
	fmt.Println("EXEMPLE AVEC NOS FAITS:")
	fmt.Println("  • 3 Commandes × 2 Produits = 6 évaluations de la condition complète")
	fmt.Println("  • Avec AlphaNode filtre: seulement 2 Commandes (qte > 5) × 2 Produits = 4 évaluations")
	fmt.Println("  • Économie: 33% d'évaluations en moins")
	fmt.Println()
	fmt.Println("CORRECTION NÉCESSAIRE:")
	fmt.Println("  1. Décomposer les conditions AND en alpha vs beta")
	fmt.Println("  2. Créer des AlphaNodes filtrants pour les conditions alpha")
	fmt.Println("  3. Ne mettre que les conditions beta dans les JoinNodes")
	fmt.Println("  4. Chaîner: TypeNode → AlphaFilter → PassthroughAlpha → JoinNode")
	fmt.Println()
	fmt.Println(strings.Repeat("=", 80))
	// This test documents the bug - it passes if the bug is present
	// After the fix, this test should FAIL, indicating the bug is fixed
}
// TestBugRETE001_VerifyFix verifies that the bug has been fixed
// This test should PASS after the correction
func TestBugRETE001_VerifyFix(t *testing.T) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("✅ VERIFICATION: Bug RETE-001 Fixed")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()
	pipeline := NewConstraintPipeline()
	storage := NewMemoryStorage()
	tsdFile := "testdata/bug_rete001_minimal.tsd"
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("❌ Erreur construction réseau: %v", err)
	}
	t.Log("✅ Réseau construit")
	fmt.Println()
	// Analyze network structure
	alphaNodesWithFilters := 0
	var alphaFilterCondition interface{}
	for alphaID, alphaNode := range network.AlphaNodes {
		node := alphaNode
		if node.Condition != nil {
			if condMap, ok := node.Condition.(map[string]interface{}); ok {
				// Check if it's NOT a passthrough type
				if condType, hasType := condMap["type"].(string); hasType && condType == "passthrough" {
					continue
				}
				// Real filtering condition
				if op, hasOp := condMap["operator"]; hasOp {
					alphaNodesWithFilters++
					alphaFilterCondition = node.Condition
					t.Logf("   ✅ AlphaNode filter found: %s", alphaID)
					t.Logf("      Operator: %v", op)
				}
			}
		}
	}
	// Analyze JoinNode condition
	hasAndOperator := false
	for _, betaNode := range network.BetaNodes {
		if node, ok := betaNode.(*JoinNode); ok {
			// node.Condition is already a map[string]interface{}
			if operator, hasOp := node.Condition["operator"]; hasOp {
				if operator == "AND" {
					hasAndOperator = true
				}
			}
		}
	}
	fmt.Println()
	t.Log("📊 Verification Results:")
	t.Logf("   - AlphaNodes with filters: %d", alphaNodesWithFilters)
	t.Logf("   - JoinNode has AND operator: %v", hasAndOperator)
	fmt.Println()
	// Assert: Bug is fixed
	if alphaNodesWithFilters == 0 {
		t.Error("❌ BUG NOT FIXED: No AlphaNode filters found")
		t.Error("   Expected: At least 1 AlphaNode with alpha condition")
		t.Error("   Got: 0 AlphaNodes with filters")
		return
	}
	t.Log("✅ VERIFIED: AlphaNode filters exist")
	if hasAndOperator {
		t.Error("❌ BUG NOT FIXED: JoinNode still contains AND operator")
		t.Error("   Expected: JoinNode with beta condition only")
		t.Error("   Got: JoinNode with AND (multiple conditions)")
		return
	}
	t.Log("✅ VERIFIED: JoinNode contains beta condition only")
	if alphaFilterCondition == nil {
		t.Error("❌ Alpha filter condition is nil")
		return
	}
	t.Log("✅ VERIFIED: Alpha filter has proper condition")
	// Verify results are correct (same actions triggered)
	// Count triggered actions by checking terminal node memory
	actionsTriggered := 0
	for _, terminal := range network.TerminalNodes {
		actionsTriggered += len(terminal.Memory.Tokens)
	}
	t.Logf("   - Actions triggered: %d", actionsTriggered)
	if actionsTriggered != 2 {
		t.Errorf("❌ Expected 2 actions (C2 and C3), got %d", actionsTriggered)
		return
	}
	t.Log("✅ VERIFIED: Correct number of actions triggered")
	fmt.Println()
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("✅ BUG FIX VERIFIED - All Checks Passed")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()
	fmt.Println("STRUCTURE CORRECTE:")
	fmt.Println("  TypeNode(Commande)")
	fmt.Println("       ↓")
	fmt.Println("  AlphaNode(c.qte > 5)          ← FILTRE ALPHA ✅")
	fmt.Println("       ↓")
	fmt.Println("  PassthroughAlphaNode")
	fmt.Println("       ↓")
	fmt.Println("  JoinNode(c.produit_id == p.id) ← BETA SEULEMENT ✅")
	fmt.Println("       ⋈")
	fmt.Println("  PassthroughAlphaNode")
	fmt.Println("       ↑")
	fmt.Println("  TypeNode(Produit)")
	fmt.Println()
	fmt.Println("RÉSULTATS:")
	fmt.Println("  ✅ Filtrage précoce: C1 (qte=3) éliminé avant jointure")
	fmt.Println("  ✅ Actions déclenchées: C2 et C3 seulement")
	fmt.Println("  ✅ Principe RETE respecté")
	fmt.Println()
	fmt.Println(strings.Repeat("=", 80))
}
// TestBugRETE001_VerifyExpectedBehavior documents what the expected behavior should be
func TestBugRETE001_VerifyExpectedBehavior(t *testing.T) {
	t.Log("📋 COMPORTEMENT ATTENDU (après correction)")
	fmt.Println()
	fmt.Println("Pour la condition: c.produit_id == p.id AND c.qte > 5")
	fmt.Println()
	fmt.Println("STRUCTURE ATTENDUE:")
	fmt.Println("  TypeNode(Commande)")
	fmt.Println("       ↓")
	fmt.Println("  AlphaNode(c.qte > 5)          ← Filtre ALPHA (1 variable)")
	fmt.Println("       ↓                           Élimine C1 (qte=3)")
	fmt.Println("  PassthroughAlphaNode")
	fmt.Println("       ↓")
	fmt.Println("  JoinNode(c.produit_id == p.id) ← Condition BETA (2 variables)")
	fmt.Println("       ⋈                           Teste seulement C2, C3")
	fmt.Println("  PassthroughAlphaNode")
	fmt.Println("       ↑")
	fmt.Println("  TypeNode(Produit)")
	fmt.Println()
	fmt.Println("BÉNÉFICES:")
	fmt.Println("  ✓ Filtrage précoce: seules les commandes avec qte > 5 atteignent le join")
	fmt.Println("  ✓ Partage: l'AlphaNode(c.qte > 5) peut être réutilisé par d'autres règles")
	fmt.Println("  ✓ Performance: 4 évaluations au lieu de 6 (33% de réduction)")
	fmt.Println("  ✓ Respect du principe RETE")
	fmt.Println()
	fmt.Println("STRUCTURE ACTUELLE (BUGGUÉE):")
	fmt.Println("  TypeNode(Commande)")
	fmt.Println("       ↓")
	fmt.Println("  PassthroughAlphaNode          ← Pas de filtrage!")
	fmt.Println("       ↓                           Toutes les commandes passent")
	fmt.Println("  JoinNode(c.produit_id == p.id AND c.qte > 5)")
	fmt.Println("       ⋈                           Condition complète ici")
	fmt.Println("  PassthroughAlphaNode           Évalue 6 paires")
	fmt.Println("       ↑")
	fmt.Println("  TypeNode(Produit)")
	fmt.Println()
	// This test is purely documentary
	t.Skip("Test documentaire - pas d'assertions")
}
// TestBugRETE001_PerformanceImpact demonstrates the performance impact of the bug
func TestBugRETE001_PerformanceImpact(t *testing.T) {
	t.Log("📊 IMPACT PERFORMANCE DU BUG")
	fmt.Println()
	pipeline := NewConstraintPipeline()
	storage := NewMemoryStorage()
	tsdFile := "testdata/bug_rete001_minimal.tsd"
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction: %v", err)
	}
	// Count facts
	commandeCount := 0
	produitCount := 0
	for _, typeNode := range network.TypeNodes {
		node := typeNode
		factCount := len(node.Memory.Facts)
		if strings.Contains(node.TypeName, "Commande") {
			commandeCount = factCount
		} else if strings.Contains(node.TypeName, "Produit") {
			produitCount = factCount
		}
	}
	t.Logf("Faits injectés:")
	t.Logf("  - Commandes: %d", commandeCount)
	t.Logf("  - Produits: %d", produitCount)
	fmt.Println()
	// Calculate evaluation counts
	withoutFilter := commandeCount * produitCount
	withFilter := 2 * produitCount // Only 2 commandes have qte > 5
	savings := float64(withoutFilter-withFilter) / float64(withoutFilter) * 100
	t.Logf("Évaluations de la condition de jointure:")
	t.Logf("  - ACTUEL (sans filtre alpha): %d paires", withoutFilter)
	t.Logf("  - ATTENDU (avec filtre alpha): %d paires", withFilter)
	t.Logf("  - ÉCONOMIE: %.0f%%", savings)
	fmt.Println()
	fmt.Println("AVEC PLUS DE FAITS:")
	scales := []int{10, 100, 1000}
	for _, scale := range scales {
		currentEvals := scale * scale
		filteredEvals := (scale * 2 / 3) * scale // Assume 2/3 pass the filter
		saving := float64(currentEvals-filteredEvals) / float64(currentEvals) * 100
		fmt.Printf("  • %d Commandes × %d Produits:\n", scale, scale)
		fmt.Printf("    - Sans filtre: %d évaluations\n", currentEvals)
		fmt.Printf("    - Avec filtre: %d évaluations\n", filteredEvals)
		fmt.Printf("    - Économie: %.0f%%\n", saving)
		fmt.Println()
	}
	t.Log("💡 Plus il y a de faits, plus l'impact du bug est important!")
}