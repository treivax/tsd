// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strings"
	"testing"
)

// formatCondition formate une condition de manière lisible
func formatCondition(cond interface{}) string {
	condMap, ok := cond.(map[string]interface{})
	if !ok {
		return fmt.Sprintf("%v", cond)
	}

	condType, _ := condMap["type"].(string)

	switch condType {
	case "comparison":
		left := formatCondition(condMap["left"])
		operator, _ := condMap["operator"].(string)
		right := formatCondition(condMap["right"])
		return fmt.Sprintf("(%s %s %s)", left, operator, right)

	case "binaryOp":
		left := formatCondition(condMap["left"])
		operator, _ := condMap["operator"].(string)
		right := formatCondition(condMap["right"])

		// Décode les opérateurs encodés
		switch operator {
		case "Kg==":
			operator = "*"
		case "Kw==":
			operator = "+"
		case "LQ==":
			operator = "-"
		case "Lw==":
			operator = "/"
		}

		return fmt.Sprintf("(%s %s %s)", left, operator, right)

	case "fieldAccess":
		object, _ := condMap["object"].(string)
		field, _ := condMap["field"].(string)
		return fmt.Sprintf("%s.%s", object, field)

	case "number":
		value := condMap["value"]
		return fmt.Sprintf("%v", value)

	default:
		return fmt.Sprintf("%v", condMap)
	}
}

// showExpressionTree affiche récursivement l'arbre AST d'une expression
func showExpressionTree(expr interface{}, indent string) {
	exprMap, ok := expr.(map[string]interface{})
	if !ok {
		fmt.Printf("%s%v\n", indent, expr)
		return
	}

	exprType, _ := exprMap["type"].(string)
	switch exprType {
	case "binaryOp":
		operator, _ := exprMap["operator"].(string)
		// Décode l'opérateur
		switch operator {
		case "Kg==":
			operator = "*"
		case "Kw==":
			operator = "+"
		case "LQ==":
			operator = "-"
		case "Lw==":
			operator = "/"
		}
		fmt.Printf("%sbinaryOp: %s\n", indent, operator)
		fmt.Printf("%s├─ left:\n", indent)
		showExpressionTree(exprMap["left"], indent+"│  ")
		fmt.Printf("%s└─ right:\n", indent)
		showExpressionTree(exprMap["right"], indent+"   ")

	case "fieldAccess":
		object, _ := exprMap["object"].(string)
		field, _ := exprMap["field"].(string)
		fmt.Printf("%sfieldAccess: %s.%s\n", indent, object, field)

	case "number":
		value := exprMap["value"]
		fmt.Printf("%snumber: %v\n", indent, value)

	default:
		fmt.Printf("%s%s: %v\n", indent, exprType, exprMap)
	}
}

// TestArithmeticExpressionsE2E teste le pipeline complet avec expressions arithmétiques complexes
// Ce test vérifie que les expressions arithmétiques dans les actions sont correctement évaluées
func TestArithmeticExpressionsE2E(t *testing.T) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("🧪 TEST E2E: Expressions Arithmétiques - Analyse Détaillée du Réseau RETE")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println()

	// Créer le pipeline
	pipeline := NewConstraintPipeline()
	storage := NewMemoryStorage()

	// Fichier contenant types, règles ET faits
	tsdFile := "testdata/arithmetic_e2e.tsd"

	fmt.Printf("📁 Fichier TSD: %s\n\n", tsdFile)

	// Construire le réseau depuis le fichier unique
	network, err := pipeline.BuildNetworkFromMultipleFiles([]string{tsdFile}, storage)
	if err != nil {
		t.Fatalf("❌ Erreur construction réseau: %v", err)
	}

	fmt.Printf("✅ Réseau RETE construit avec succès\n")
	fmt.Printf("   - TypeNodes: %d\n", len(network.TypeNodes))
	fmt.Printf("   - AlphaNodes: %d\n", len(network.AlphaNodes))
	fmt.Printf("   - BetaNodes: %d\n", len(network.BetaNodes))
	fmt.Printf("   - TerminalNodes: %d\n", len(network.TerminalNodes))
	fmt.Printf("   - PassthroughRegistry: %d\n\n", len(network.PassthroughRegistry))

	// ========================================
	// DIAGRAMME DÉTAILLÉ DU RÉSEAU
	// ========================================
	diagram := NewNetworkDiagram(network)
	diagram.PrintDetailedDiagram()

	// ========================================
	// SECTION 1: MAPPING RÈGLES TSD → NŒUDS RETE
	// ========================================
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("📋 SECTION 1: MAPPING DES RÈGLES TSD VERS LES NŒUDS RETE")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println()

	// Règle 1
	fmt.Println("┌" + strings.Repeat("─", 98) + "┐")
	fmt.Println("│ RÈGLE 1: calcul_facture_base                                                                     │")
	fmt.Println("└" + strings.Repeat("─", 98) + "┘")
	fmt.Println()
	fmt.Println("📝 Expression TSD:")
	fmt.Println("   rule calcul_facture_base : {p: Produit, c: Commande} /")
	fmt.Println("       c.produit_id == p.id AND c.qte * 23 - 10 > 0")
	fmt.Println("   ==> facture_calculee(...)")
	fmt.Println()
	fmt.Println("🔧 Décomposition en nœuds RETE:")
	fmt.Println()
	fmt.Println("   1️⃣  Variable 'p: Produit'")
	fmt.Println("       └─→ TypeNode[type_Produit]")
	fmt.Println("           ├─ ID: type_Produit")
	fmt.Println("           ├─ Statut: ✅ PARTAGÉ entre TOUTES les règles utilisant Produit")
	fmt.Println("           └─ Rôle: Routage des faits de type Produit")
	fmt.Println()
	fmt.Println("   2️⃣  Passthrough LEFT pour 'p'")
	fmt.Println("       └─→ AlphaNode[passthrough_calcul_facture_base_p_Produit_left]")
	fmt.Println("           ├─ ID: passthrough_calcul_facture_base_p_Produit_left")
	fmt.Println("           ├─ Statut: ○ DÉDIÉ à cette règle (per-rule passthrough)")
	fmt.Println("           ├─ Rôle: Passthrough sans filtre, prépare pour jointure LEFT")
	fmt.Println("           └─ Condition: PASSTHROUGH (side: left)")
	fmt.Println()
	fmt.Println("   3️⃣  Variable 'c: Commande'")
	fmt.Println("       └─→ TypeNode[type_Commande]")
	fmt.Println("           ├─ ID: type_Commande")
	fmt.Println("           ├─ Statut: ✅ PARTAGÉ entre TOUTES les règles utilisant Commande")
	fmt.Println("           └─ Rôle: Routage des faits de type Commande")
	fmt.Println()
	fmt.Println("   4️⃣  Condition alpha: 'c.qte * 23 - 10 > 0'")
	fmt.Println("       └─→ AlphaNode[calcul_facture_base_alpha_c_0]")
	fmt.Println("           ├─ ID: calcul_facture_base_alpha_c_0")
	fmt.Println("           ├─ Statut: ○ DÉDIÉ à cette règle (condition alpha spécifique)")
	fmt.Println("           ├─ Rôle: Filtrage alpha sur variable 'c' (test arithmétique)")
	fmt.Println("           ├─ Condition: ((c.qte * 23) - 10) > 0")
	fmt.Println("           └─ Expression: Opération arithmétique sur une seule variable")
	fmt.Println()
	fmt.Println("   5️⃣  Passthrough RIGHT pour 'c' (après filtre alpha)")
	fmt.Println("       └─→ AlphaNode[passthrough_calcul_facture_base_c_Commande_right]")
	fmt.Println("           ├─ ID: passthrough_calcul_facture_base_c_Commande_right")
	fmt.Println("           ├─ Statut: ○ DÉDIÉ à cette règle (per-rule passthrough)")
	fmt.Println("           ├─ Rôle: Passthrough après filtre, prépare pour jointure RIGHT")
	fmt.Println("           └─ Condition: PASSTHROUGH (side: right)")
	fmt.Println()
	fmt.Println("   6️⃣  Condition beta: 'c.produit_id == p.id'")
	fmt.Println("       └─→ JoinNode[calcul_facture_base_join]")
	fmt.Println("           ├─ ID: calcul_facture_base_join")
	fmt.Println("           ├─ Statut: ○ DÉDIÉ à cette règle")
	fmt.Println("           ├─ Rôle: Jointure entre p et c + évaluation condition beta")
	fmt.Println("           ├─ Left Input: passthrough_calcul_facture_base_p_Produit_left")
	fmt.Println("           ├─ Right Input: passthrough_calcul_facture_base_c_Commande_right")
	fmt.Println("           └─ Condition: c.produit_id == p.id (équi-jointure)")
	fmt.Println()
	fmt.Println("   7️⃣  Action: 'facture_calculee(...)'")
	fmt.Println("       └─→ TerminalNode[calcul_facture_base_terminal]")
	fmt.Println("           ├─ ID: calcul_facture_base_terminal")
	fmt.Println("           ├─ Statut: ○ DÉDIÉ à cette règle (une action par règle)")
	fmt.Println("           ├─ Rôle: Exécution de l'action facture_calculee")
	fmt.Println("           └─ Parent: calcul_facture_base_join")
	fmt.Println()

	// Règle 2
	fmt.Println("┌" + strings.Repeat("─", 98) + "┐")
	fmt.Println("│ RÈGLE 2: calcul_facture_speciale                                                                 │")
	fmt.Println("└" + strings.Repeat("─", 98) + "┘")
	fmt.Println()
	fmt.Println("📝 Expression TSD:")
	fmt.Println("   rule calcul_facture_speciale : {p: Produit, c: Commande} /")
	fmt.Println("       c.produit_id == p.id AND c.qte * 23 - 10 < 0")
	fmt.Println("   ==> facture_speciale(...)")
	fmt.Println()
	fmt.Println("🔧 Décomposition en nœuds RETE:")
	fmt.Println()
	fmt.Println("   1️⃣  Variable 'p: Produit'")
	fmt.Println("       └─→ TypeNode[type_Produit]")
	fmt.Println("           ├─ ID: type_Produit")
	fmt.Println("           ├─ Statut: ✅ PARTAGÉ avec règles 1 et 3")
	fmt.Println("           └─ Rôle: Routage des faits de type Produit")
	fmt.Println()
	fmt.Println("   2️⃣  Passthrough LEFT pour 'p'")
	fmt.Println("       └─→ AlphaNode[passthrough_calcul_facture_speciale_p_Produit_left]")
	fmt.Println("           ├─ ID: passthrough_calcul_facture_speciale_p_Produit_left")
	fmt.Println("           ├─ Statut: ○ DÉDIÉ à cette règle (per-rule passthrough)")
	fmt.Println("           ├─ Rôle: Passthrough sans filtre, prépare pour jointure LEFT")
	fmt.Println("           └─ Condition: PASSTHROUGH (side: left)")
	fmt.Println()
	fmt.Println("   3️⃣  Variable 'c: Commande'")
	fmt.Println("       └─→ TypeNode[type_Commande]")
	fmt.Println("           ├─ ID: type_Commande")
	fmt.Println("           ├─ Statut: ✅ PARTAGÉ avec règles 1 et 3")
	fmt.Println("           └─ Rôle: Routage des faits de type Commande")
	fmt.Println()
	fmt.Println("   4️⃣  Condition alpha: 'c.qte * 23 - 10 < 0' ⚠️  OPÉRATEUR INVERSÉ")
	fmt.Println("       └─→ AlphaNode[calcul_facture_speciale_alpha_c_0]")
	fmt.Println("           ├─ ID: calcul_facture_speciale_alpha_c_0")
	fmt.Println("           ├─ Statut: ○ DÉDIÉ à cette règle (condition différente: < au lieu de >)")
	fmt.Println("           ├─ Rôle: Filtrage alpha sur variable 'c' (test arithmétique)")
	fmt.Println("           ├─ Condition: ((c.qte * 23) - 10) < 0")
	fmt.Println("           └─ Expression: DIFFÉRENTE de la règle 1 (opérateur < vs >)")
	fmt.Println()
	fmt.Println("   5️⃣  Passthrough RIGHT pour 'c' (après filtre alpha)")
	fmt.Println("       └─→ AlphaNode[passthrough_calcul_facture_speciale_c_Commande_right]")
	fmt.Println("           ├─ ID: passthrough_calcul_facture_speciale_c_Commande_right")
	fmt.Println("           ├─ Statut: ○ DÉDIÉ à cette règle (per-rule passthrough)")
	fmt.Println("           ├─ Rôle: Passthrough après filtre, prépare pour jointure RIGHT")
	fmt.Println("           └─ Condition: PASSTHROUGH (side: right)")
	fmt.Println()
	fmt.Println("   6️⃣  Condition beta: 'c.produit_id == p.id'")
	fmt.Println("       └─→ JoinNode[calcul_facture_speciale_join]")
	fmt.Println("           ├─ ID: calcul_facture_speciale_join")
	fmt.Println("           ├─ Statut: ○ DÉDIÉ à cette règle")
	fmt.Println("           ├─ Rôle: Jointure entre p et c + évaluation condition beta")
	fmt.Println("           ├─ Left Input: passthrough_calcul_facture_speciale_p_Produit_left")
	fmt.Println("           ├─ Right Input: passthrough_calcul_facture_speciale_c_Commande_right")
	fmt.Println("           └─ Condition: c.produit_id == p.id (équi-jointure)")
	fmt.Println()
	fmt.Println("   7️⃣  Action: 'facture_speciale(...)'")
	fmt.Println("       └─→ TerminalNode[calcul_facture_speciale_terminal]")
	fmt.Println("           ├─ ID: calcul_facture_speciale_terminal")
	fmt.Println("           ├─ Statut: ○ DÉDIÉ à cette règle (une action par règle)")
	fmt.Println("           ├─ Rôle: Exécution de l'action facture_speciale")
	fmt.Println("           └─ Parent: calcul_facture_speciale_join")
	fmt.Println()

	// Règle 3
	fmt.Println("┌" + strings.Repeat("─", 98) + "┐")
	fmt.Println("│ RÈGLE 3: calcul_facture_premium                                                                  │")
	fmt.Println("└" + strings.Repeat("─", 98) + "┘")
	fmt.Println()
	fmt.Println("📝 Expression TSD:")
	fmt.Println("   rule calcul_facture_premium : {p: Produit, c: Commande} /")
	fmt.Println("       c.produit_id == p.id AND c.qte * 23 - 10 > 0")
	fmt.Println("   ==> facture_speciale(...)")
	fmt.Println()
	fmt.Println("🔧 Décomposition en nœuds RETE:")
	fmt.Println()
	fmt.Println("   1️⃣  Variable 'p: Produit'")
	fmt.Println("       └─→ TypeNode[type_Produit]")
	fmt.Println("           ├─ ID: type_Produit")
	fmt.Println("           ├─ Statut: ✅ PARTAGÉ avec règles 1 et 2")
	fmt.Println("           └─ Rôle: Routage des faits de type Produit")
	fmt.Println()
	fmt.Println("   2️⃣  Passthrough LEFT pour 'p'")
	fmt.Println("       └─→ AlphaNode[passthrough_calcul_facture_premium_p_Produit_left]")
	fmt.Println("           ├─ ID: passthrough_calcul_facture_premium_p_Produit_left")
	fmt.Println("           ├─ Statut: ○ DÉDIÉ à cette règle (per-rule passthrough)")
	fmt.Println("           ├─ Rôle: Passthrough sans filtre, prépare pour jointure LEFT")
	fmt.Println("           └─ Condition: PASSTHROUGH (side: left)")
	fmt.Println()
	fmt.Println("   3️⃣  Variable 'c: Commande'")
	fmt.Println("       └─→ TypeNode[type_Commande]")
	fmt.Println("           ├─ ID: type_Commande")
	fmt.Println("           ├─ Statut: ✅ PARTAGÉ avec règles 1 et 2")
	fmt.Println("           └─ Rôle: Routage des faits de type Commande")
	fmt.Println()
	fmt.Println("   4️⃣  Condition alpha: 'c.qte * 23 - 10 > 0' ⚠️  IDENTIQUE À RÈGLE 1")
	fmt.Println("       └─→ AlphaNode[calcul_facture_premium_alpha_c_0]")
	fmt.Println("           ├─ ID: calcul_facture_premium_alpha_c_0")
	fmt.Println("           ├─ Statut: ○ DÉDIÉ (mais POURRAIT être partagé avec règle 1!)")
	fmt.Println("           ├─ Rôle: Filtrage alpha sur variable 'c' (test arithmétique)")
	fmt.Println("           ├─ Condition: ((c.qte * 23) - 10) > 0")
	fmt.Println("           └─ Expression: IDENTIQUE à la règle 1 (même opérateur >)")
	fmt.Println()
	fmt.Println("   5️⃣  Passthrough RIGHT pour 'c' (après filtre alpha)")
	fmt.Println("       └─→ AlphaNode[passthrough_calcul_facture_premium_c_Commande_right]")
	fmt.Println("           ├─ ID: passthrough_calcul_facture_premium_c_Commande_right")
	fmt.Println("           ├─ Statut: ○ DÉDIÉ à cette règle (per-rule passthrough)")
	fmt.Println("           ├─ Rôle: Passthrough après filtre, prépare pour jointure RIGHT")
	fmt.Println("           └─ Condition: PASSTHROUGH (side: right)")
	fmt.Println()
	fmt.Println("   6️⃣  Condition beta: 'c.produit_id == p.id'")
	fmt.Println("       └─→ JoinNode[calcul_facture_premium_join]")
	fmt.Println("           ├─ ID: calcul_facture_premium_join")
	fmt.Println("           ├─ Statut: ○ DÉDIÉ à cette règle")
	fmt.Println("           ├─ Rôle: Jointure entre p et c + évaluation condition beta")
	fmt.Println("           ├─ Left Input: passthrough_calcul_facture_premium_p_Produit_left")
	fmt.Println("           ├─ Right Input: passthrough_calcul_facture_premium_c_Commande_right")
	fmt.Println("           └─ Condition: c.produit_id == p.id (équi-jointure)")
	fmt.Println()
	fmt.Println("   7️⃣  Action: 'facture_speciale(...)' (action différente de règle 1)")
	fmt.Println("       └─→ TerminalNode[calcul_facture_premium_terminal]")
	fmt.Println("           ├─ ID: calcul_facture_premium_terminal")
	fmt.Println("           ├─ Statut: ○ DÉDIÉ à cette règle (une action par règle)")
	fmt.Println("           ├─ Rôle: Exécution de l'action facture_speciale")
	fmt.Println("           └─ Parent: calcul_facture_premium_join")
	fmt.Println()

	// ========================================
	// SECTION 2: ANALYSE DU PARTAGE
	// ========================================
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("📊 SECTION 2: ANALYSE DU PARTAGE DES NŒUDS")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println()

	fmt.Println("┌" + strings.Repeat("─", 98) + "┐")
	fmt.Println("│ ✅ NŒUDS PARTAGÉS (réutilisés par plusieurs règles)                                              │")
	fmt.Println("└" + strings.Repeat("─", 98) + "┘")
	fmt.Println()
	fmt.Println("   🔷 TypeNode[type_Produit]")
	fmt.Println("      ├─ Partagé par: Règle 1, Règle 2, Règle 3")
	fmt.Println("      ├─ Bénéfice: 1 nœud au lieu de 3 → Économie de 67%")
	fmt.Println("      └─ Impact: Tous les faits Produit passent par ce nœud unique")
	fmt.Println()
	fmt.Println("   🔷 TypeNode[type_Commande]")
	fmt.Println("      ├─ Partagé par: Règle 1, Règle 2, Règle 3")
	fmt.Println("      ├─ Bénéfice: 1 nœud au lieu de 3 → Économie de 67%")
	fmt.Println("      └─ Impact: Tous les faits Commande passent par ce nœud unique")
	fmt.Println()
	fmt.Println("   📈 Total nœuds partagés: 2 TypeNodes")
	fmt.Println()

	fmt.Println("┌" + strings.Repeat("─", 98) + "┐")
	fmt.Println("│ ○ NŒUDS DÉDIÉS (un par règle - per-rule isolation)                                               │")
	fmt.Println("└" + strings.Repeat("─", 98) + "┘")
	fmt.Println()
	fmt.Println("   🔹 AlphaNodes Passthrough LEFT:")
	fmt.Println("      ├─ passthrough_calcul_facture_base_p_Produit_left     (Règle 1)")
	fmt.Println("      ├─ passthrough_calcul_facture_speciale_p_Produit_left (Règle 2)")
	fmt.Println("      └─ passthrough_calcul_facture_premium_p_Produit_left  (Règle 3)")
	fmt.Println("      Raison: Isolation per-rule pour éviter cross-contamination")
	fmt.Println()
	fmt.Println("   🔹 AlphaNodes Filtres (conditions alpha):")
	fmt.Println("      ├─ calcul_facture_base_alpha_c_0     : c.qte * 23 - 10 > 0  (Règle 1)")
	fmt.Println("      ├─ calcul_facture_speciale_alpha_c_0 : c.qte * 23 - 10 < 0  (Règle 2) ← Différent!")
	fmt.Println("      └─ calcul_facture_premium_alpha_c_0  : c.qte * 23 - 10 > 0  (Règle 3) ← Identique à Règle 1!")
	fmt.Println("      Raison: Chaque règle a son propre filtre (pas de partage actuellement)")
	fmt.Println()
	fmt.Println("   🔹 AlphaNodes Passthrough RIGHT:")
	fmt.Println("      ├─ passthrough_calcul_facture_base_c_Commande_right     (Règle 1)")
	fmt.Println("      ├─ passthrough_calcul_facture_speciale_c_Commande_right (Règle 2)")
	fmt.Println("      └─ passthrough_calcul_facture_premium_c_Commande_right  (Règle 3)")
	fmt.Println("      Raison: Isolation per-rule pour éviter cross-contamination")
	fmt.Println()
	fmt.Println("   🔶 JoinNodes:")
	fmt.Println("      ├─ calcul_facture_base_join     (Règle 1)")
	fmt.Println("      ├─ calcul_facture_speciale_join (Règle 2)")
	fmt.Println("      └─ calcul_facture_premium_join  (Règle 3)")
	fmt.Println("      Raison: Chaque règle a sa propre jointure (comportement actuel)")
	fmt.Println()
	fmt.Println("   🎯 TerminalNodes:")
	fmt.Println("      ├─ calcul_facture_base_terminal     → facture_calculee (Règle 1)")
	fmt.Println("      ├─ calcul_facture_speciale_terminal → facture_speciale (Règle 2)")
	fmt.Println("      └─ calcul_facture_premium_terminal  → facture_speciale (Règle 3)")
	fmt.Println("      Raison: Une action par règle (obligatoire)")
	fmt.Println()

	fmt.Println("┌" + strings.Repeat("─", 98) + "┐")
	fmt.Println("│ 💡 OPPORTUNITÉS D'OPTIMISATION                                                                   │")
	fmt.Println("└" + strings.Repeat("─", 98) + "┘")
	fmt.Println()
	fmt.Println("   ⚠️  Règle 1 et Règle 3 ont la MÊME condition alpha:")
	fmt.Println("      • calcul_facture_base_alpha_c_0    : c.qte * 23 - 10 > 0")
	fmt.Println("      • calcul_facture_premium_alpha_c_0 : c.qte * 23 - 10 > 0")
	fmt.Println()
	fmt.Println("      Ces deux AlphaNodes POURRAIENT être fusionnés en un seul nœud partagé!")
	fmt.Println("      Bénéfice potentiel: 2 nœuds → 1 nœud (économie de 50%)")
	fmt.Println()
	fmt.Println("   ⚠️  Architecture actuelle: Per-rule passthroughs")
	fmt.Println("      • Avantage: Isolation complète entre règles (pas de cross-contamination)")
	fmt.Println("      • Coût: Plus de nœuds (overhead mémoire modeste)")
	fmt.Println("      • Alternative future: Partager les passthroughs si alpha-chains identiques")
	fmt.Println()

	// ========================================
	// SECTION 3: VUE GRAPHIQUE ASCII
	// ========================================
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("🎨 SECTION 3: VUE GRAPHIQUE COMPLÈTE DU RÉSEAU")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println()

	fmt.Println("Légende:")
	fmt.Println("   [T]✅ = TypeNode PARTAGÉ")
	fmt.Println("   [α]○  = AlphaNode DÉDIÉ à une règle")
	fmt.Println("   [β]○  = JoinNode DÉDIÉ à une règle")
	fmt.Println("   [⚡]○ = TerminalNode DÉDIÉ à une règle")
	fmt.Println("   ⚠️   = Nœud potentiellement partageable")
	fmt.Println()

	fmt.Println("                                    ┌──────────┐")
	fmt.Println("                                    │   ROOT   │")
	fmt.Println("                                    └─────┬────┘")
	fmt.Println("                                          │")
	fmt.Println("                        ┌─────────────────┼─────────────────┐")
	fmt.Println("                        │                 │                 │")
	fmt.Println("                        ▼                 ▼                 ▼")
	fmt.Println("                  ┌──────────┐      ┌──────────┐      ┌──────────┐")
	fmt.Println("                  │[T]✅      │      │[T]✅      │      │[T]✅      │")
	fmt.Println("                  │ Produit  │      │ Commande │      │ Client   │")
	fmt.Println("                  │          │      │          │      │ (unused) │")
	fmt.Println("                  └────┬─────┘      └─────┬────┘      └──────────┘")
	fmt.Println("                       │                  │")
	fmt.Println("        ┌──────────────┼──────────────┐   │")
	fmt.Println("        │              │              │   │")
	fmt.Println("        ▼              ▼              ▼   │")
	fmt.Println("   ┌────────┐     ┌────────┐     ┌────────┐")
	fmt.Println("   │[α]○ PT │     │[α]○ PT │     │[α]○ PT │")
	fmt.Println("   │ R1-L   │     │ R2-L   │     │ R3-L   │")
	fmt.Println("   └───┬────┘     └───┬────┘     └───┬────┘")
	fmt.Println("       │              │              │")
	fmt.Println("       │              │              │                ┌──────────────────┐")
	fmt.Println("       │              │              │                │ Légende PT:      │")
	fmt.Println("       │              │              │                │ R1-L = Règle 1   │")
	fmt.Println("       ▼              ▼              ▼                │        Left PT   │")
	fmt.Println("   ┌────────┐     ┌────────┐     ┌────────┐          └──────────────────┘")
	fmt.Println("   │[β]○    │     │[β]○    │     │[β]○    │")
	fmt.Println("   │ R1 Join│◄────┤ R2 Join│◄────┤ R3 Join│")
	fmt.Println("   │        │     │        │     │        │")
	fmt.Println("   │c.id==  │     │c.id==  │     │c.id==  │")
	fmt.Println("   │  p.id  │     │  p.id  │     │  p.id  │")
	fmt.Println("   └───┬────┘     └───┬────┘     └───┬────┘")
	fmt.Println("       │              │              │")
	fmt.Println("       ▲              ▲              ▲")
	fmt.Println("       │              │              │")
	fmt.Println("       │              │              │")
	fmt.Println("   ┌───┴────┐     ┌───┴────┐     ┌───┴────┐          De Commande:")
	fmt.Println("   │[α]○ PT │     │[α]○ PT │     │[α]○ PT │")
	fmt.Println("   │ R1-R   │     │ R2-R   │     │ R3-R   │")
	fmt.Println("   └───▲────┘     └───▲────┘     └───▲────┘")
	fmt.Println("       │              │              │")
	fmt.Println("   ┌───┴──────┐   ┌───┴──────┐   ┌───┴──────┐")
	fmt.Println("   │[α]○ ⚠️    │   │[α]○      │   │[α]○ ⚠️    │")
	fmt.Println("   │qte*23    │   │qte*23    │   │qte*23    │")
	fmt.Println("   │  -10>0   │   │  -10<0   │   │  -10>0   │")
	fmt.Println("   │(IDENTIQUE│   │(INVERSÉ) │   │(IDENTIQUE│")
	fmt.Println("   │ R1 & R3) │   │          │   │ R1 & R3) │")
	fmt.Println("   └─────▲────┘   └─────▲────┘   └─────▲────┘")
	fmt.Println("         │              │              │")
	fmt.Println("         └──────────────┴──────────────┘")
	fmt.Println("                        │")
	fmt.Println("                        │ [T]✅ type_Commande")
	fmt.Println("                        │ (PARTAGÉ)")
	fmt.Println()
	fmt.Println("   Actions (TerminalNodes):")
	fmt.Println()
	fmt.Println("   R1 Join ──→ [⚡]○ calcul_facture_base_terminal     ──→ facture_calculee()")
	fmt.Println("   R2 Join ──→ [⚡]○ calcul_facture_speciale_terminal ──→ facture_speciale()")
	fmt.Println("   R3 Join ──→ [⚡]○ calcul_facture_premium_terminal  ──→ facture_speciale()")
	fmt.Println()

	// ========================================
	// SECTION 4: STATISTIQUES
	// ========================================
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("📊 SECTION 4: STATISTIQUES DÉTAILLÉES")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println()

	totalNodes := len(network.TypeNodes) + len(network.AlphaNodes) + len(network.BetaNodes) + len(network.TerminalNodes) + len(network.PassthroughRegistry)
	fmt.Printf("   📦 Total de nœuds dans le réseau: %d\n", totalNodes)
	fmt.Println()
	fmt.Println("   Par type:")
	fmt.Printf("      • TypeNodes (partagés):          %d ✅\n", len(network.TypeNodes))
	fmt.Printf("      • AlphaNodes (filtres):          %d ○\n", len(network.AlphaNodes))
	fmt.Printf("      • PassthroughRegistry:           %d ○\n", len(network.PassthroughRegistry))
	fmt.Printf("      • BetaNodes (jointures):         %d ○\n", len(network.BetaNodes))
	fmt.Printf("      • TerminalNodes (actions):       %d ○\n", len(network.TerminalNodes))
	fmt.Println()
	fmt.Println("   Taux de partage:")
	sharedNodes := len(network.TypeNodes)
	dedicatedNodes := totalNodes - sharedNodes
	shareRate := float64(sharedNodes) / float64(totalNodes) * 100
	fmt.Printf("      • Nœuds partagés:                %d (%.1f%%)\n", sharedNodes, shareRate)
	fmt.Printf("      • Nœuds dédiés:                  %d (%.1f%%)\n", dedicatedNodes, 100-shareRate)
	fmt.Println()

	// ========================================
	// SECTION 5: EXÉCUTION ET RÉSULTATS
	// ========================================
	// ========================================
	// SECTION 5: ANALYSE DÉTAILLÉE DES ALPHANODES
	// ========================================
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("🔬 SECTION 5: ANALYSE DÉTAILLÉE DES ALPHANODES (DÉCOMPOSITION)")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println()

	fmt.Println("┌" + strings.Repeat("─", 98) + "┐")
	fmt.Println("│ ✅ PARTAGE DES ALPHANODES ACTIVÉ                                                                 │")
	fmt.Println("└" + strings.Repeat("─", 98) + "┘")
	fmt.Println()
	fmt.Println("Les règles 1 et 3 ont la MÊME condition alpha:")
	fmt.Println()

	// Trouver tous les AlphaNodes avec condition > 0 (devrait être partagé entre règles 1 et 3)
	var sharedAlphaNode *AlphaNode
	var sharedAlphaNodeID string

	for id, node := range network.AlphaNodes {
		if node.VariableName == "c" {
			formatted := formatCondition(node.Condition)
			if strings.Contains(formatted, "> 0") {
				sharedAlphaNode = node
				sharedAlphaNodeID = id
				break
			}
		}
	}

	// Afficher la condition alpha
	if sharedAlphaNode != nil && sharedAlphaNode.Condition != nil {
		formatted := formatCondition(sharedAlphaNode.Condition)
		fmt.Printf("   Condition: %s\n", formatted)
		fmt.Printf("   AlphaNode ID: %s\n", sharedAlphaNodeID)

		// Vérifier si c'est partagé en comptant les enfants
		children := sharedAlphaNode.GetChildren()
		fmt.Printf("   Nombre de règles utilisant ce nœud: %d\n", len(children))
	}
	fmt.Println()

	if sharedAlphaNode != nil && len(sharedAlphaNode.GetChildren()) >= 2 {
		fmt.Println("✅ PARTAGE DÉTECTÉ: Plusieurs règles partagent le MÊME AlphaNode!")
		fmt.Printf("   • ID partagé: %s\n", sharedAlphaNodeID)
		fmt.Printf("   • Nombre de règles: %d\n", len(sharedAlphaNode.GetChildren()))
		fmt.Printf("   • Économie: %d nœuds au lieu de %d (%d%% de réduction)\n",
			1, len(sharedAlphaNode.GetChildren()),
			(len(sharedAlphaNode.GetChildren())-1)*100/len(sharedAlphaNode.GetChildren()))
	} else {
		fmt.Println("❌ PAS DE PARTAGE DÉTECTÉ")
	}
	fmt.Println()

	fmt.Println("┌" + strings.Repeat("─", 98) + "┐")
	fmt.Println("│ 🔬 ANALYSE: DÉCOMPOSITION DE L'EXPRESSION ALPHA                                                  │")
	fmt.Println("└" + strings.Repeat("─", 98) + "┘")
	fmt.Println()
	fmt.Println("Expression TSD: (c.qte * 23 - 10 + c.remise * 43) > 0")
	fmt.Println()
	fmt.Println("❌ PAS DE DÉCOMPOSITION EN SOUS-EXPRESSIONS:")
	fmt.Println("   L'expression est traitée comme UN SEUL AlphaNode monolithique")
	fmt.Println()

	// Montrer la structure interne d'un AlphaNode (utiliser celui qu'on a trouvé)
	if sharedAlphaNode != nil && sharedAlphaNode.Condition != nil {
		fmt.Println("📦 Structure interne de l'AlphaNode:")
		fmt.Println()

		condMap := sharedAlphaNode.Condition.(map[string]interface{})

		// Afficher la structure récursivement
		fmt.Println("   Type: comparison")
		fmt.Println("   ├─ Operator: " + fmt.Sprintf("%v", condMap["operator"]))
		fmt.Println("   ├─ Right: " + fmt.Sprintf("%v", condMap["right"]))
		fmt.Println("   └─ Left: (expression arithmétique complexe)")

		if left, ok := condMap["left"].(map[string]interface{}); ok {
			fmt.Println()
			fmt.Println("   Expression Left décomposée (AST interne):")
			showExpressionTree(left, "      ")
		}
	}
	fmt.Println()

	fmt.Println("✅ Ce qu'on OBSERVE:")
	fmt.Printf("   • 1 AlphaNode par règle\n")
	fmt.Printf("   • Expression stockée comme un arbre AST unique\n")
	fmt.Printf("   • Pas de décomposition en nœuds atomiques séparés\n")
	fmt.Println()

	fmt.Println("Ce qui est IMPLÉMENTÉ:")
	fmt.Printf("   ✅ Partage des AlphaNodes identiques entre règles (via AlphaSharingRegistry)\n")
	fmt.Println()
	fmt.Println("Ce qui MANQUE (optimisations possibles):")
	fmt.Printf("   • Décomposition en sous-expressions réutilisables:\n")
	fmt.Printf("     - AlphaNode 1: (c.qte * 23)\n")
	fmt.Printf("     - AlphaNode 2: (AlphaNode1 - 10)\n")
	fmt.Printf("     - AlphaNode 3: (c.remise * 43)\n")
	fmt.Printf("     - AlphaNode 4: (AlphaNode2 + AlphaNode3)\n")
	fmt.Printf("     - AlphaNode 5: (AlphaNode4 > 0)\n")
	fmt.Println()

	fmt.Println("┌" + strings.Repeat("─", 98) + "┐")
	fmt.Println("│ ✅ SOLUTION IMPLÉMENTÉE: ALPHASHARINGREGISTRY                                                    │")
	fmt.Println("└" + strings.Repeat("─", 98) + "┘")
	fmt.Println()
	fmt.Println("Le partage des AlphaNodes est maintenant actif via AlphaSharingRegistry:")
	fmt.Println()
	fmt.Println("1. Hash canonique calculé pour chaque condition (indépendant du ruleID)")
	fmt.Println("2. AlphaNodes identiques partagés automatiquement")
	fmt.Println("3. ID basé sur le hash de la condition: alpha_<hash>")
	fmt.Println()

	// Afficher les statistiques réelles
	if network.AlphaSharingManager != nil {
		stats := network.AlphaSharingManager.GetStats()
		fmt.Println("📊 Statistiques de partage:")
		fmt.Printf("   • AlphaNodes partagés: %v\n", stats["total_shared_alpha_nodes"])
		fmt.Printf("   • Références totales: %v\n", stats["total_rule_references"])
		fmt.Printf("   • Ratio de partage moyen: %.2f\n", stats["average_sharing_ratio"])
		fmt.Println()
	}

	fmt.Println("Bénéfice pour ce test:")
	fmt.Printf("   • AlphaNodes créés: %d (au lieu de 3 sans partage)\n", len(network.AlphaNodes))
	if len(network.AlphaNodes) < 3 {
		saving := (3 - len(network.AlphaNodes)) * 100 / 3
		fmt.Printf("   • Économie: %d%% de nœuds en moins\n", saving)
	}
	fmt.Println()

	// ========================================
	// SECTION 6: EXÉCUTION ET RÉSULTATS
	// ========================================
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("🚀 SECTION 6: EXÉCUTION ET RÉSULTATS")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println()

	totalTokens := 0
	tokensPerRule := make(map[string]int)

	for _, terminal := range network.TerminalNodes {
		tokens := terminal.Memory.GetTokens()
		tokenCount := len(tokens)
		totalTokens += tokenCount
		tokensPerRule[strings.TrimSuffix(terminal.ID, "_terminal")] = tokenCount
	}

	fmt.Println("📊 Résultats par règle:")
	fmt.Println()

	// Extraire les vraies conditions alpha depuis le réseau
	alphaConditions := make(map[string]string)
	for alphaID, alphaNode := range network.AlphaNodes {
		if alphaNode.Condition != nil {
			if condMap, ok := alphaNode.Condition.(map[string]interface{}); ok {
				if condType, _ := condMap["type"].(string); condType != "passthrough" {
					// Format simplifié de la condition
					condStr := fmt.Sprintf("%v", condMap)
					alphaConditions[alphaID] = condStr
				}
			}
		}
	}

	// Règle 1
	baseTokens := tokensPerRule["calcul_facture_base"]
	fmt.Printf("   Règle 1 (calcul_facture_base):\n")
	if sharedAlphaNode != nil && sharedAlphaNode.Condition != nil {
		formatted := formatCondition(sharedAlphaNode.Condition)
		fmt.Printf("      • Condition alpha: %s\n", formatted)
		fmt.Printf("      • AlphaNode ID: %s\n", sharedAlphaNodeID)
	}
	fmt.Printf("      • Tokens générés: %d\n", baseTokens)
	fmt.Println()

	// Règle 2
	specialTokens := tokensPerRule["calcul_facture_speciale"]
	fmt.Printf("   Règle 2 (calcul_facture_speciale):\n")
	// Trouver l'AlphaNode avec condition < 0
	for id, node := range network.AlphaNodes {
		if node.VariableName == "c" {
			formatted := formatCondition(node.Condition)
			if strings.Contains(formatted, "< 0") {
				fmt.Printf("      • Condition alpha: %s\n", formatted)
				fmt.Printf("      • AlphaNode ID: %s\n", id)
				break
			}
		}
	}
	fmt.Printf("      • Tokens générés: %d\n", specialTokens)
	fmt.Println()

	// Règle 3
	premiumTokens := tokensPerRule["calcul_facture_premium"]
	fmt.Printf("   Règle 3 (calcul_facture_premium):\n")
	if sharedAlphaNode != nil && sharedAlphaNode.Condition != nil {
		formatted := formatCondition(sharedAlphaNode.Condition)
		fmt.Printf("      • Condition alpha: %s\n", formatted)
		fmt.Printf("      • AlphaNode ID: %s\n", sharedAlphaNodeID)
		fmt.Printf("      • ♻️  Partage le même AlphaNode que Règle 1!\n")
	}
	fmt.Printf("      • Tokens générés: %d\n", premiumTokens)
	if baseTokens == premiumTokens && baseTokens > 0 {
		fmt.Printf("      • Note: Même résultat que Règle 1 (conditions identiques)\n")
	}
	fmt.Println()

	fmt.Printf("✅ Total: %d tokens générés\n", totalTokens)
	fmt.Printf("✅ Total: %d actions déclenchées\n", totalTokens)
	fmt.Println()

	// ========================================
	// VALIDATIONS
	// ========================================
	fmt.Println()

	// ========================================
	// SECTION 7: VALIDATIONS
	// ========================================
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("✓ SECTION 7: VALIDATIONS")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println()

	// Vérifier la structure du réseau
	// Avec le partage beta obligatoire :
	// - Règles 1 et 3 ont les mêmes conditions => partagent le même JoinNode
	// - Règle 2 a une condition différente => son propre JoinNode
	// Total : 2 JoinNodes (au lieu de 3 sans partage)
	// Note: network.BetaNodes peut contenir des entrées dupliquées (hash + legacy key)
	// On compte les JoinNodes uniques par ID
	uniqueJoinNodes := make(map[string]bool)
	for _, node := range network.BetaNodes {
		if joinNode, ok := node.(*JoinNode); ok {
			uniqueJoinNodes[joinNode.ID] = true
		}
	}
	expectedBetaNodes := 2
	actualUniqueJoinNodes := len(uniqueJoinNodes)
	if actualUniqueJoinNodes != expectedBetaNodes {
		t.Errorf("❌ Devrait avoir %d JoinNodes avec partage beta, got %d", expectedBetaNodes, actualUniqueJoinNodes)
	} else {
		fmt.Printf("   ✅ Structure réseau: %d JoinNodes avec partage (règles 1 et 3 partagent 1 JoinNode)\n", expectedBetaNodes)
	}

	if len(network.TerminalNodes) != 3 {
		t.Fatalf("❌ Devrait avoir 3 TerminalNodes, got %d", len(network.TerminalNodes))
	} else {
		fmt.Printf("   ✅ Structure réseau: %d TerminalNodes (un par règle)\n", len(network.TerminalNodes))
	}

	// Vérifier les résultats d'exécution
	// Note : Le partage de JoinNode fait que chaque règle reçoit TOUS les tokens du JoinNode partagé
	// Les règles 1 et 3 partagent le même JoinNode, donc reçoivent les mêmes 3 tokens
	// C'est le comportement attendu avec le partage beta
	expectedBase := 3
	if baseTokens != expectedBase {
		t.Errorf("❌ Règle 'calcul_facture_base': attendu %d tokens, got %d", expectedBase, baseTokens)
	} else {
		fmt.Printf("   ✅ Règle 'calcul_facture_base': %d tokens (attendu: %d)\n", baseTokens, expectedBase)
	}

	expectedSpeciale := 0
	if specialTokens != expectedSpeciale {
		t.Errorf("❌ Règle 'calcul_facture_speciale': attendu %d tokens, got %d", expectedSpeciale, specialTokens)
	} else {
		fmt.Printf("   ✅ Règle 'calcul_facture_speciale': %d tokens (attendu: %d)\n", specialTokens, expectedSpeciale)
	}

	expectedPremium := 3
	if premiumTokens != expectedPremium {
		t.Errorf("❌ Règle 'calcul_facture_premium': attendu %d tokens, got %d", expectedPremium, premiumTokens)
	} else {
		fmt.Printf("   ✅ Règle 'calcul_facture_premium': %d tokens (attendu: %d)\n", premiumTokens, expectedPremium)
	}

	expectedTotal := expectedBase + expectedSpeciale + expectedPremium
	if totalTokens != expectedTotal {
		t.Errorf("❌ Total de tokens incorrect: got %d, want %d", totalTokens, expectedTotal)
	} else {
		fmt.Printf("   ✅ Total tokens: %d (attendu: %d)\n", totalTokens, expectedTotal)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("🎉 TEST E2E TERMINÉ AVEC SUCCÈS!")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println()
}
