package rete

import (
	"github.com/treivax/tsd/tsdio"
	"fmt"
	"sort"
	"strings"
)

// NetworkDiagram génère un diagramme ASCII détaillé du réseau RETE
type NetworkDiagram struct {
	network *ReteNetwork
}

// NewNetworkDiagram crée un nouveau générateur de diagramme
func NewNetworkDiagram(network *ReteNetwork) *NetworkDiagram {
	return &NetworkDiagram{network: network}
}

// PrintDetailedDiagram affiche un diagramme complet avec les opérateurs
func (nd *NetworkDiagram) PrintDetailedDiagram() {
	tsdio.Println()
	tsdio.Println(strings.Repeat("═", 120))
	tsdio.Println("📊 DIAGRAMME DÉTAILLÉ DU RÉSEAU RETE")
	tsdio.Println(strings.Repeat("═", 120))
	tsdio.Println()

	// 1. Type Nodes
	nd.printTypeNodes()

	// 2. Alpha Nodes (décomposés)
	nd.printAlphaNodes()

	// 3. Passthrough Nodes
	nd.printPassthroughNodes()

	// 4. Join Nodes (Beta)
	nd.printJoinNodes()

	// 5. Router Nodes
	nd.printRouterNodes()

	// 6. Terminal Nodes
	nd.printTerminalNodes()

	// 7. Flow Diagram
	nd.printFlowDiagram()

	// 8. Summary
	nd.printSummary()

	tsdio.Println()
	tsdio.Println(strings.Repeat("═", 120))
	tsdio.Println()
}

func (nd *NetworkDiagram) printTypeNodes() {
	tsdio.Println("┌" + strings.Repeat("─", 118) + "┐")
	tsdio.Println("│ 1️⃣  TYPE NODES (Routage par type)                                                                    │")
	tsdio.Println("└" + strings.Repeat("─", 118) + "┘")
	tsdio.Println()

	typeNames := make([]string, 0, len(nd.network.TypeNodes))
	for typeName := range nd.network.TypeNodes {
		typeNames = append(typeNames, typeName)
	}
	sort.Strings(typeNames)

	for _, typeName := range typeNames {
		node := nd.network.TypeNodes[typeName]
		tsdio.Printf("   [T] type_%s\n", typeName)
		tsdio.Printf("       │ Type: %s\n", typeName)
		tsdio.Printf("       │ Enfants: %d nœuds\n", len(node.Children))
		tsdio.Printf("       └─→ Propage tous les faits de type %s\n", typeName)
		tsdio.Println()
	}
}

func (nd *NetworkDiagram) printAlphaNodes() {
	tsdio.Println("┌" + strings.Repeat("─", 118) + "┐")
	tsdio.Println("│ 2️⃣  ALPHA NODES (Filtres et calculs atomiques)                                                       │")
	tsdio.Println("└" + strings.Repeat("─", 118) + "┘")
	tsdio.Println()

	// Grouper par variable
	alphasByVar := make(map[string][]*AlphaNode)
	for _, node := range nd.network.AlphaNodes {
		if node.VariableName != "" {
			alphasByVar[node.VariableName] = append(alphasByVar[node.VariableName], node)
		}
	}

	vars := make([]string, 0, len(alphasByVar))
	for v := range alphasByVar {
		vars = append(vars, v)
	}
	sort.Strings(vars)

	for _, varName := range vars {
		nodes := alphasByVar[varName]
		tsdio.Printf("   📍 Variable: %s (%d nœuds)\n", varName, len(nodes))
		tsdio.Println()

		// Trier les nœuds par ID pour un affichage cohérent
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].ID < nodes[j].ID
		})

		for i, node := range nodes {
			nd.printAlphaNodeDetail(node, i+1, len(nodes))
		}
		tsdio.Println()
	}
}

func (nd *NetworkDiagram) printAlphaNodeDetail(node *AlphaNode, index, total int) {
	var symbol string
	if index == 1 {
		symbol = "┌─"
	} else if index == total {
		symbol = "└─"
	} else {
		symbol = "├─"
	}

	tsdio.Printf("      %s [α] %s\n", symbol, node.ID)

	// Extraire les détails de la condition
	if condMap, ok := node.Condition.(map[string]interface{}); ok {
		condType, _ := condMap["type"].(string)

		switch condType {
		case "passthrough":
			side, _ := condMap["side"].(string)
			tsdio.Printf("         │ Type: PASSTHROUGH\n")
			if side != "" {
				tsdio.Printf("         │ Side: %s\n", side)
			}
			tsdio.Printf("         │ Opération: Propagation sans filtre\n")

		case "comparison":
			operator, _ := condMap["operator"].(string)
			left := nd.formatExpression(condMap["left"])
			right := nd.formatExpression(condMap["right"])
			tsdio.Printf("         │ Type: COMPARISON\n")
			tsdio.Printf("         │ Opérateur: %s\n", nd.symbolizeOperator(operator))
			tsdio.Printf("         │ Expression: %s %s %s\n", left, nd.symbolizeOperator(operator), right)

		case "binaryOp":
			operator, _ := condMap["operator"].(string)
			left := nd.formatExpression(condMap["left"])
			right := nd.formatExpression(condMap["right"])
			tsdio.Printf("         │ Type: BINARY OPERATION\n")
			tsdio.Printf("         │ Opérateur: %s\n", nd.symbolizeOperator(operator))
			tsdio.Printf("         │ Calcul: %s %s %s\n", left, nd.symbolizeOperator(operator), right)

		case "tempResult":
			stepName, _ := condMap["step_name"].(string)
			stepIdx, _ := condMap["step_idx"].(int)
			tsdio.Printf("         │ Type: TEMP RESULT\n")
			tsdio.Printf("         │ Step: %s (étape %d)\n", stepName, stepIdx)
			tsdio.Printf("         │ Opération: Stockage résultat intermédiaire\n")

		default:
			tsdio.Printf("         │ Type: %s\n", condType)
		}

		// Note: Les informations de partage sont affichées dans le résumé
	}

	tsdio.Printf("         │ Enfants: %d\n", len(node.Children))
}

func (nd *NetworkDiagram) printPassthroughNodes() {
	tsdio.Println("┌" + strings.Repeat("─", 118) + "┐")
	tsdio.Println("│ 3️⃣  PASSTHROUGH NODES (Préparation pour jointure)                                                    │")
	tsdio.Println("└" + strings.Repeat("─", 118) + "┘")
	tsdio.Println()

	if len(nd.network.PassthroughRegistry) == 0 {
		tsdio.Println("   (Aucun passthrough node)")
		tsdio.Println()
		return
	}

	// Grouper par side
	leftNodes := make([]string, 0)
	rightNodes := make([]string, 0)

	for key := range nd.network.PassthroughRegistry {
		if strings.Contains(key, "_left") {
			leftNodes = append(leftNodes, key)
		} else if strings.Contains(key, "_right") {
			rightNodes = append(rightNodes, key)
		}
	}

	sort.Strings(leftNodes)
	sort.Strings(rightNodes)

	tsdio.Println("   LEFT Side (tokens pour jointure gauche):")
	for _, key := range leftNodes {
		node := nd.network.PassthroughRegistry[key]
		tsdio.Printf("      [⇒] %s\n", node.ID)
		tsdio.Printf("          │ Rôle: Passthrough LEFT\n")
		tsdio.Printf("          │ Enfants: %d\n", len(node.Children))
		tsdio.Println()
	}

	tsdio.Println("   RIGHT Side (tokens pour jointure droite):")
	for _, key := range rightNodes {
		node := nd.network.PassthroughRegistry[key]
		tsdio.Printf("      [⇒] %s\n", node.ID)
		tsdio.Printf("          │ Rôle: Passthrough RIGHT\n")
		tsdio.Printf("          │ Enfants: %d\n", len(node.Children))
		tsdio.Println()
	}
}

func (nd *NetworkDiagram) printJoinNodes() {
	tsdio.Println("┌" + strings.Repeat("─", 118) + "┐")
	tsdio.Println("│ 4️⃣  JOIN NODES (Jointures Beta)                                                                      │")
	tsdio.Println("└" + strings.Repeat("─", 118) + "┘")
	tsdio.Println()

	if len(nd.network.BetaNodes) == 0 {
		tsdio.Println("   (Aucun join node)")
		tsdio.Println()
		return
	}

	// Compter les JoinNodes uniques
	uniqueJoins := make(map[string]*JoinNode)
	for _, node := range nd.network.BetaNodes {
		if joinNode, ok := node.(*JoinNode); ok {
			if _, exists := uniqueJoins[joinNode.ID]; !exists {
				uniqueJoins[joinNode.ID] = joinNode
			}
		}
	}

	// Trier par ID
	joinIDs := make([]string, 0, len(uniqueJoins))
	for id := range uniqueJoins {
		joinIDs = append(joinIDs, id)
	}
	sort.Strings(joinIDs)

	for i, id := range joinIDs {
		joinNode := uniqueJoins[id]
		symbol := "├─"
		if i == len(joinIDs)-1 {
			symbol = "└─"
		}

		tsdio.Printf("   %s [⋈] %s\n", symbol, joinNode.ID)
		tsdio.Printf("      │ Type: JOIN NODE\n")
		tsdio.Printf("      │ Variables LEFT: %v\n", joinNode.LeftVariables)
		tsdio.Printf("      │ Variables RIGHT: %v\n", joinNode.RightVariables)

		// Afficher les JoinConditions
		if len(joinNode.JoinConditions) > 0 {
			tsdio.Printf("      │ Conditions de jointure:\n")
			for _, jc := range joinNode.JoinConditions {
				tsdio.Printf("      │   • %s.%s %s %s.%s\n",
					jc.LeftVar, jc.LeftField,
					nd.symbolizeOperator(jc.Operator),
					jc.RightVar, jc.RightField)
			}
		}

		// Vérifier si partagé
		rulesUsing := nd.findRulesUsingJoinNode(id)
		if len(rulesUsing) > 1 {
			tsdio.Printf("      │ Partagé par: %d règles ✨ SHARED\n", len(rulesUsing))
			for _, rule := range rulesUsing {
				tsdio.Printf("      │   - %s\n", rule)
			}
		} else if len(rulesUsing) == 1 {
			tsdio.Printf("      │ Utilisé par: %s\n", rulesUsing[0])
		}

		tsdio.Printf("      │ Enfants: %d\n", len(joinNode.Children))
		tsdio.Println()
	}
}

func (nd *NetworkDiagram) printRouterNodes() {
	// Chercher les RuleRouterNodes dans les enfants des JoinNodes
	routers := make([]*RuleRouterNode, 0)
	seenRouters := make(map[string]bool)

	for _, betaNode := range nd.network.BetaNodes {
		if joinNode, ok := betaNode.(*JoinNode); ok {
			for _, child := range joinNode.Children {
				if router, ok := child.(*RuleRouterNode); ok {
					// Éviter les doublons (le même JoinNode peut être dans BetaNodes plusieurs fois avec des clés différentes)
					if !seenRouters[router.ID] {
						routers = append(routers, router)
						seenRouters[router.ID] = true
					}
				}
			}
		}
	}

	if len(routers) == 0 {
		return
	}

	tsdio.Println("┌" + strings.Repeat("─", 118) + "┐")
	tsdio.Println("│ 5️⃣  ROUTER NODES (Routage des tokens vers les règles)                                                │")
	tsdio.Println("└" + strings.Repeat("─", 118) + "┘")
	tsdio.Println()

	for _, router := range routers {
		tsdio.Printf("   [🔀] %s\n", router.ID)
		tsdio.Printf("       │ Type: RULE ROUTER\n")
		tsdio.Printf("       │ Pour la règle: %s\n", router.RuleID)
		tsdio.Printf("       │ Depuis JoinNode: %s\n", router.JoinNodeID)
		if router.TerminalNode != nil {
			tsdio.Printf("       │ Vers TerminalNode: %s\n", router.TerminalNode.ID)
		}
		tsdio.Println()
	}
}

func (nd *NetworkDiagram) printTerminalNodes() {
	tsdio.Println("┌" + strings.Repeat("─", 118) + "┐")
	tsdio.Println("│ 6️⃣  TERMINAL NODES (Actions)                                                                         │")
	tsdio.Println("└" + strings.Repeat("─", 118) + "┘")
	tsdio.Println()

	terminalIDs := make([]string, 0, len(nd.network.TerminalNodes))
	for id := range nd.network.TerminalNodes {
		terminalIDs = append(terminalIDs, id)
	}
	sort.Strings(terminalIDs)

	for i, id := range terminalIDs {
		terminal := nd.network.TerminalNodes[id]
		symbol := "├─"
		if i == len(terminalIDs)-1 {
			symbol = "└─"
		}

		tsdio.Printf("   %s [⚡] %s\n", symbol, terminal.ID)
		if terminal.Action != nil {
			tsdio.Printf("      │ Action: %s\n", terminal.Action.Type)
		}
		tsdio.Printf("      │ Tokens en mémoire: %d\n", len(terminal.Memory.Tokens))
		tsdio.Println()
	}
}

func (nd *NetworkDiagram) printFlowDiagram() {
	tsdio.Println("┌" + strings.Repeat("─", 118) + "┐")
	tsdio.Println("│ 7️⃣  DIAGRAMME DE FLUX (Architecture complète)                                                        │")
	tsdio.Println("└" + strings.Repeat("─", 118) + "┘")
	tsdio.Println()

	tsdio.Println("   Expression TSD des règles:")
	tsdio.Println("   ═════════════════════════════════════════════════════════════════════════════════════════")
	tsdio.Println()
	tsdio.Println("   R1: c.produit_id == p.id AND (c.qte × 23 - 10 + c.remise × 43) > 0  → facture_calculee")
	tsdio.Println("   R2: c.produit_id == p.id AND (c.qte × 23 - 10 + c.remise × 43) < 0  → facture_speciale")
	tsdio.Println("   R3: c.produit_id == p.id AND (c.qte × 23 - 10 + c.remise × 43) > 0  → facture_speciale")
	tsdio.Println()
	tsdio.Println("   Architecture RETE:")
	tsdio.Println("   ═════════════════════════════════════════════════════════════════════════════════════════")
	tsdio.Println()
	tsdio.Println("                                     ┌─────────────────┐")
	tsdio.Println("                                     │  [T] Commande   │")
	tsdio.Println("                                     │   type_Commande │")
	tsdio.Println("                                     └────────┬────────┘")
	tsdio.Println("                                              │")
	tsdio.Println("                      ┌───────────────────────┼───────────────────────┐")
	tsdio.Println("                      │                       │                       │")
	tsdio.Println("   ┌──────────────────▼──────────────┐        │        ┌──────────────▼──────────────┐")
	tsdio.Println("   │ [α] c.qte × 23                  │◄───────┴───────►│ [α] c.qte × 23              │")
	tsdio.Println("   │  alpha_1362ff5a962dca07         │  PARTAGÉ R1-R3  │  alpha_1362ff5a962dca07     │")
	tsdio.Println("   └────────────┬────────────────────┘                 └────────────┬────────────────┘")
	tsdio.Println("                │                                                   │")
	tsdio.Println("   ┌────────────▼────────────────────┐                 ┌────────────▼────────────────┐")
	tsdio.Println("   │ [α] <temp_1> - 10               │◄───────────────►│ [α] <temp_1> - 10           │")
	tsdio.Println("   │  alpha_e2ae7bbb66d00288         │  PARTAGÉ R1-R3  │  alpha_e2ae7bbb66d00288     │")
	tsdio.Println("   └────────────┬────────────────────┘                 └────────────┬────────────────┘")
	tsdio.Println("                │                                                   │")
	tsdio.Println("   ┌────────────▼────────────────────┐                 ┌────────────▼────────────────┐")
	tsdio.Println("   │ [α] c.remise × 43               │◄───────────────►│ [α] c.remise × 43           │")
	tsdio.Println("   │  alpha_c4780a7d3c271103         │  PARTAGÉ R1-R3  │  alpha_c4780a7d3c271103     │")
	tsdio.Println("   └────────────┬────────────────────┘                 └────────────┬────────────────┘")
	tsdio.Println("                │                                                   │")
	tsdio.Println("   ┌────────────▼────────────────────┐                 ┌────────────▼────────────────┐")
	tsdio.Println("   │ [α] <temp_2> + <temp_3>         │◄───────────────►│ [α] <temp_2> + <temp_3>     │")
	tsdio.Println("   │  alpha_e03528dec0e1f043         │  PARTAGÉ R1-R3  │  alpha_e03528dec0e1f043     │")
	tsdio.Println("   └────────┬───────────┬────────────┘                 └────────┬───────────┬────────┘")
	tsdio.Println("            │           │                                       │           │")
	tsdio.Println("   ┌────────▼──────┐    │                              ┌────────▼──────┐    │")
	tsdio.Println("   │ [α] > 0       │    │                              │ [α] < 0       │    │")
	tsdio.Println("   │ R1 & R3       │    │                              │ R2 seule      │    │")
	tsdio.Println("   │ alpha_2913... │    │                              │ alpha_81a5... │    │")
	tsdio.Println("   └───────┬───────┘    │                              └───────┬───────┘    │")
	tsdio.Println("           │            │                                      │            │")
	tsdio.Println("   ┌───────▼────────┐   │                              ┌───────▼────────┐   │")
	tsdio.Println("   │ [⇒] Passthrough│   │                              │ [⇒] Passthrough│   │")
	tsdio.Println("   │ RIGHT R1       │   │                              │ RIGHT R2       │   │")
	tsdio.Println("   └───────┬────────┘   │                              └───────┬────────┘   │")
	tsdio.Println("           │            │                                      │            │")
	tsdio.Println("           │    ┌───────▼────────┐                            │    ┌───────▼────────┐")
	tsdio.Println("           │    │ [⇒] Passthrough│                            │    │ [⇒] Passthrough│")
	tsdio.Println("           │    │ RIGHT R3       │                            │    │ RIGHT R3 (skip)│")
	tsdio.Println("           │    └───────┬────────┘                            │    └───────┬────────┘")
	tsdio.Println("           │            │                                      │            │")
	tsdio.Println("           └────────┬───┘                                      └────────┬───┘")
	tsdio.Println("                    │                                                   │")
	tsdio.Println("   ┌─────────────────────────────────┐                 ┌─────────────────────────────────┐")
	tsdio.Println("   │  [T] Produit (LEFT)             │                 │  [T] Produit (LEFT)             │")
	tsdio.Println("   └────────┬────────────────────────┘                 └────────┬────────────────────────┘")
	tsdio.Println("            │                                                   │")
	tsdio.Println("   ┌────────▼────────────────────┐                    ┌────────▼────────────────────┐")
	tsdio.Println("   │ [⇒] Passthrough LEFT R1     │                    │ [⇒] Passthrough LEFT R2     │")
	tsdio.Println("   └────────┬────────────────────┘                    └────────┬────────────────────┘")
	tsdio.Println("            │                                                  │")
	tsdio.Println("            ├──────────────────────────────┐                  │")
	tsdio.Println("            │                              │                  │")
	tsdio.Println("   ┌────────▼─────────────────────┐  ┌────▼───────────────────────┐")
	tsdio.Println("   │ [⋈] JoinNode SHARED          │  │ [⋈] JoinNode R2            │")
	tsdio.Println("   │ join_514c9d1bff12fa4f        │  │ join_118236e6b5bc9f95      │")
	tsdio.Println("   │ c.produit_id == p.id         │  │ c.produit_id == p.id       │")
	tsdio.Println("   │ Partagé: R1 & R3             │  │ Dédié: R2                  │")
	tsdio.Println("   └────┬───────────────────┬─────┘  └────────┬───────────────────┘")
	tsdio.Println("        │                   │                 │")
	tsdio.Println("   ┌────▼────────┐   ┌──────▼─────────┐   ┌──▼───────────────┐")
	tsdio.Println("   │ [⚡] R1      │   │ [🔀] Router R3 │   │ [⚡] R2           │")
	tsdio.Println("   │ Terminal    │   │ RuleRouterNode │   │ Terminal         │")
	tsdio.Println("   │ 3 tokens ✓  │   └──────┬─────────┘   │ 0 tokens ✓       │")
	tsdio.Println("   └─────────────┘          │             └──────────────────┘")
	tsdio.Println("                     ┌──────▼─────────┐")
	tsdio.Println("                     │ [⚡] R3         │")
	tsdio.Println("                     │ Terminal       │")
	tsdio.Println("                     │ 3 tokens ✓     │")
	tsdio.Println("                     └────────────────┘")
	tsdio.Println()
	tsdio.Println("   Légende:")
	tsdio.Println("   ────────")
	tsdio.Println("   [T]  = TypeNode     (routage par type)")
	tsdio.Println("   [α]  = AlphaNode    (filtrage/calcul atomique)")
	tsdio.Println("   [⇒]  = Passthrough  (préparation jointure)")
	tsdio.Println("   [⋈]  = JoinNode     (jointure beta)")
	tsdio.Println("   [🔀] = RouterNode   (routage tokens)")
	tsdio.Println("   [⚡] = TerminalNode (action)")
	tsdio.Println()
	tsdio.Println("   Points clés:")
	tsdio.Println("   ───────────")
	tsdio.Println("   • AlphaNodes décomposés partagés entre R1 et R3 (mêmes conditions)")
	tsdio.Println("   • JoinNode partagé entre R1 et R3 (mêmes conditions alpha + beta)")
	tsdio.Println("   • R2 a son propre JoinNode (condition alpha différente: < au lieu de >)")
	tsdio.Println("   • RuleRouterNode route les tokens du JoinNode partagé vers R3")
	tsdio.Println("   • R1 connectée directement au JoinNode (première règle)")
	tsdio.Println()
}

func (nd *NetworkDiagram) printSummary() {
	tsdio.Println("┌" + strings.Repeat("─", 118) + "┐")
	tsdio.Println("│ 📈 RÉSUMÉ DU RÉSEAU                                                                                  │")
	tsdio.Println("└" + strings.Repeat("─", 118) + "┘")
	tsdio.Println()

	uniqueJoins := make(map[string]bool)
	for _, node := range nd.network.BetaNodes {
		if joinNode, ok := node.(*JoinNode); ok {
			uniqueJoins[joinNode.ID] = true
		}
	}

	tsdio.Printf("   Type Nodes:        %3d (routage par type)\n", len(nd.network.TypeNodes))
	tsdio.Printf("   Alpha Nodes:       %3d (filtres et calculs atomiques)\n", len(nd.network.AlphaNodes))
	tsdio.Printf("   Passthrough Nodes: %3d (préparation jointure)\n", len(nd.network.PassthroughRegistry))
	tsdio.Printf("   Join Nodes:        %3d (jointures beta)\n", len(uniqueJoins))
	tsdio.Printf("   Terminal Nodes:    %3d (actions)\n", len(nd.network.TerminalNodes))
	tsdio.Println()

	// Calculer les statistiques de partage alpha
	sharedAlphaCount := 0
	if nd.network.AlphaSharingManager != nil {
		stats := nd.network.AlphaSharingManager.GetStats()
		if totalShared, ok := stats["total_shared_alpha_nodes"].(int); ok {
			sharedAlphaCount = totalShared
		}
	}

	sharedJoinCount := 0
	for id := range uniqueJoins {
		rules := nd.findRulesUsingJoinNode(id)
		if len(rules) > 1 {
			sharedJoinCount++
		}
	}

	tsdio.Printf("   📊 Statistiques de partage:\n")
	tsdio.Printf("      • AlphaNodes partagés:  %d / %d\n", sharedAlphaCount, len(nd.network.AlphaNodes))
	tsdio.Printf("      • JoinNodes partagés:   %d / %d\n", sharedJoinCount, len(uniqueJoins))
	tsdio.Println()
}

// Helper functions

func (nd *NetworkDiagram) symbolizeOperator(op string) string {
	symbols := map[string]string{
		"==":  "==",
		"!=":  "≠",
		"<":   "<",
		">":   ">",
		"<=":  "≤",
		">=":  "≥",
		"+":   "+",
		"-":   "-",
		"*":   "×",
		"/":   "÷",
		"AND": "∧",
		"OR":  "∨",
	}

	if sym, ok := symbols[op]; ok {
		return sym
	}
	return op
}

func (nd *NetworkDiagram) formatExpression(expr interface{}) string {
	if exprMap, ok := expr.(map[string]interface{}); ok {
		exprType, _ := exprMap["type"].(string)

		switch exprType {
		case "fieldAccess":
			obj, _ := exprMap["object"].(string)
			field, _ := exprMap["field"].(string)
			return fmt.Sprintf("%s.%s", obj, field)

		case "number":
			val, _ := exprMap["value"].(float64)
			return fmt.Sprintf("%.0f", val)

		case "tempResult":
			stepName, _ := exprMap["step_name"].(string)
			return fmt.Sprintf("<%s>", stepName)

		case "binaryOp":
			op, _ := exprMap["operator"].(string)
			left := nd.formatExpression(exprMap["left"])
			right := nd.formatExpression(exprMap["right"])
			return fmt.Sprintf("(%s %s %s)", left, nd.symbolizeOperator(op), right)
		}
	}

	return fmt.Sprintf("%v", expr)
}

func (nd *NetworkDiagram) findRulesUsingJoinNode(joinNodeID string) []string {
	rules := make([]string, 0)

	// Chercher dans BetaNodes avec les legacy keys
	for key := range nd.network.BetaNodes {
		if strings.HasSuffix(key, "_join") {
			// Extraire le ruleID de la legacy key
			ruleID := strings.TrimSuffix(key, "_join")
			// Vérifier si ce ruleID utilise ce JoinNode
			if node, exists := nd.network.BetaNodes[key]; exists {
				if joinNode, ok := node.(*JoinNode); ok {
					if joinNode.ID == joinNodeID {
						rules = append(rules, ruleID)
					}
				}
			}
		}
	}

	return rules
}
