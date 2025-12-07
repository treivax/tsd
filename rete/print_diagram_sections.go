// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strings"
)

// Constantes pour le formatage du diagramme
const (
	diagramWidth         = 118
	diagramSeparatorChar = "─"
	diagramSeparator     = "═════════════════════════════════════════════════════════════════════════════════════════"
)

// printDiagramHeader affiche l'en-tête d'une section du diagramme
func printDiagramHeader(title string, width int) {
	fmt.Println("┌" + strings.Repeat(diagramSeparatorChar, width) + "┐")
	// Calculer le padding pour centrer le titre
	padding := width - len(title)
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("│ %s%s│\n", title, strings.Repeat(" ", padding-2))
	fmt.Println("└" + strings.Repeat(diagramSeparatorChar, width) + "┘")
}

// printRulesExpression affiche la section d'expression des règles TSD
func printRulesExpression() {
	fmt.Println("   Expression TSD des règles:")
	fmt.Println("   " + diagramSeparator)
	fmt.Println()
	fmt.Println("   R1: c.produit_id == p.id AND (c.qte × 23 - 10 + c.remise × 43) > 0  → facture_calculee")
	fmt.Println("   R2: c.produit_id == p.id AND (c.qte × 23 - 10 + c.remise × 43) < 0  → facture_speciale")
	fmt.Println("   R3: c.produit_id == p.id AND (c.qte × 23 - 10 + c.remise × 43) > 0  → facture_speciale")
	fmt.Println()
}

// printArchitectureDiagram affiche le diagramme d'architecture ASCII complet
func printArchitectureDiagram() {
	fmt.Println("   Architecture RETE:")
	fmt.Println("   " + diagramSeparator)
	fmt.Println()
	fmt.Println("                                     ┌─────────────────┐")
	fmt.Println("                                     │  [T] Commande   │")
	fmt.Println("                                     │   type_Commande │")
	fmt.Println("                                     └────────┬────────┘")
	fmt.Println("                                              │")
	fmt.Println("                      ┌───────────────────────┼───────────────────────┐")
	fmt.Println("                      │                       │                       │")
	fmt.Println("   ┌──────────────────▼──────────────┐        │        ┌──────────────▼──────────────┐")
	fmt.Println("   │ [α] c.qte × 23                  │◄───────┴───────►│ [α] c.qte × 23              │")
	fmt.Println("   │  alpha_1362ff5a962dca07         │  PARTAGÉ R1-R3  │  alpha_1362ff5a962dca07     │")
	fmt.Println("   └────────────┬────────────────────┘                 └────────────┬────────────────┘")
	fmt.Println("                │                                                   │")
	fmt.Println("   ┌────────────▼────────────────────┐                 ┌────────────▼────────────────┐")
	fmt.Println("   │ [α] <temp_1> - 10               │◄───────────────►│ [α] <temp_1> - 10           │")
	fmt.Println("   │  alpha_e2ae7bbb66d00288         │  PARTAGÉ R1-R3  │  alpha_e2ae7bbb66d00288     │")
	fmt.Println("   └────────────┬────────────────────┘                 └────────────┬────────────────┘")
	fmt.Println("                │                                                   │")
	fmt.Println("   ┌────────────▼────────────────────┐                 ┌────────────▼────────────────┐")
	fmt.Println("   │ [α] c.remise × 43               │◄───────────────►│ [α] c.remise × 43           │")
	fmt.Println("   │  alpha_c4780a7d3c271103         │  PARTAGÉ R1-R3  │  alpha_c4780a7d3c271103     │")
	fmt.Println("   └────────────┬────────────────────┘                 └────────────┬────────────────┘")
	fmt.Println("                │                                                   │")
	fmt.Println("   ┌────────────▼────────────────────┐                 ┌────────────▼────────────────┐")
	fmt.Println("   │ [α] <temp_2> + <temp_3>         │◄───────────────►│ [α] <temp_2> + <temp_3>     │")
	fmt.Println("   │  alpha_e03528dec0e1f043         │  PARTAGÉ R1-R3  │  alpha_e03528dec0e1f043     │")
	fmt.Println("   └────────┬───────────┬────────────┘                 └────────┬───────────┬────────┘")
	fmt.Println("            │           │                                       │           │")
	fmt.Println("   ┌────────▼──────┐    │                              ┌────────▼──────┐    │")
	fmt.Println("   │ [α] > 0       │    │                              │ [α] < 0       │    │")
	fmt.Println("   │ R1 & R3       │    │                              │ R2 seule      │    │")
	fmt.Println("   │ alpha_2913... │    │                              │ alpha_81a5... │    │")
	fmt.Println("   └───────┬───────┘    │                              └───────┬───────┘    │")
	fmt.Println("           │            │                                      │            │")
	fmt.Println("   ┌───────▼────────┐   │                              ┌───────▼────────┐   │")
	fmt.Println("   │ [⇒] Passthrough│   │                              │ [⇒] Passthrough│   │")
	fmt.Println("   │ RIGHT R1       │   │                              │ RIGHT R2       │   │")
	fmt.Println("   └───────┬────────┘   │                              └───────┬────────┘   │")
	fmt.Println("           │            │                                      │            │")
	fmt.Println("           │    ┌───────▼────────┐                            │    ┌───────▼────────┐")
	fmt.Println("           │    │ [⇒] Passthrough│                            │    │ [⇒] Passthrough│")
	fmt.Println("           │    │ RIGHT R3       │                            │    │ RIGHT R3 (skip)│")
	fmt.Println("           │    └───────┬────────┘                            │    └───────┬────────┘")
	fmt.Println("           │            │                                      │            │")
	fmt.Println("           └────────┬───┘                                      └────────┬───┘")
	fmt.Println("                    │                                                   │")
	fmt.Println("   ┌─────────────────────────────────┐                 ┌─────────────────────────────────┐")
	fmt.Println("   │  [T] Produit (LEFT)             │                 │  [T] Produit (LEFT)             │")
	fmt.Println("   └────────┬────────────────────────┘                 └────────┬────────────────────────┘")
	fmt.Println("            │                                                   │")
	fmt.Println("   ┌────────▼────────────────────┐                    ┌────────▼────────────────────┐")
	fmt.Println("   │ [⇒] Passthrough LEFT R1     │                    │ [⇒] Passthrough LEFT R2     │")
	fmt.Println("   └────────┬────────────────────┘                    └────────┬────────────────────┘")
	fmt.Println("            │                                                  │")
	fmt.Println("            ├──────────────────────────────┐                  │")
	fmt.Println("            │                              │                  │")
	fmt.Println("   ┌────────▼─────────────────────┐  ┌────▼───────────────────────┐")
	fmt.Println("   │ [⋈] JoinNode SHARED          │  │ [⋈] JoinNode R2            │")
	fmt.Println("   │ join_514c9d1bff12fa4f        │  │ join_118236e6b5bc9f95      │")
	fmt.Println("   │ c.produit_id == p.id         │  │ c.produit_id == p.id       │")
	fmt.Println("   │ Partagé: R1 & R3             │  │ Dédié: R2                  │")
	fmt.Println("   └────┬───────────────────┬─────┘  └────────┬───────────────────┘")
	fmt.Println("        │                   │                 │")
	fmt.Println("   ┌────▼────────┐   ┌──────▼─────────┐   ┌──▼───────────────┐")
	fmt.Println("   │ [⚡] R1      │   │ [🔀] Router R3 │   │ [⚡] R2           │")
	fmt.Println("   │ Terminal    │   │ RuleRouterNode │   │ Terminal         │")
	fmt.Println("   │ 3 tokens ✓  │   └──────┬─────────┘   │ 0 tokens ✓       │")
	fmt.Println("   └─────────────┘          │             └──────────────────┘")
	fmt.Println("                     ┌──────▼─────────┐")
	fmt.Println("                     │ [⚡] R3         │")
	fmt.Println("                     │ Terminal       │")
	fmt.Println("                     │ 3 tokens ✓     │")
	fmt.Println("                     └────────────────┘")
	fmt.Println()
}

// printDiagramLegend affiche la légende du diagramme
func printDiagramLegend() {
	fmt.Println("   Légende:")
	fmt.Println("   ────────")
	fmt.Println("   [T]  = TypeNode     (routage par type)")
	fmt.Println("   [α]  = AlphaNode    (filtrage/calcul atomique)")
	fmt.Println("   [⇒]  = Passthrough  (préparation jointure)")
	fmt.Println("   [⋈]  = JoinNode     (jointure beta)")
	fmt.Println("   [🔀] = RouterNode   (routage tokens)")
	fmt.Println("   [⚡] = TerminalNode (action)")
	fmt.Println()
}

// printKeyPoints affiche les points clés du diagramme
func printKeyPoints() {
	fmt.Println("   Points clés:")
	fmt.Println("   ───────────")
	fmt.Println("   • AlphaNodes décomposés partagés entre R1 et R3 (mêmes conditions)")
	fmt.Println("   • JoinNode partagé entre R1 et R3 (mêmes conditions alpha + beta)")
	fmt.Println("   • R2 a son propre JoinNode (condition alpha différente: < au lieu de >)")
	fmt.Println("   • RuleRouterNode route les tokens du JoinNode partagé vers R3")
	fmt.Println("   • R1 connectée directement au JoinNode (première règle)")
	fmt.Println()
}
