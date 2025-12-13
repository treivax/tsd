// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
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
	fmt.Println()
	fmt.Println(strings.Repeat("═", 120))
	fmt.Println("📊 DIAGRAMME DÉTAILLÉ DU RÉSEAU RETE")
	fmt.Println(strings.Repeat("═", 120))
	fmt.Println()

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

	fmt.Println()
	fmt.Println(strings.Repeat("═", 120))
	fmt.Println()
}

func (nd *NetworkDiagram) printTypeNodes() {
	fmt.Println("┌" + strings.Repeat("─", 118) + "┐")
	fmt.Println("│ 1️⃣  TYPE NODES (Routage par type)                                                                    │")
	fmt.Println("└" + strings.Repeat("─", 118) + "┘")
	fmt.Println()

	typeNames := make([]string, 0, len(nd.network.TypeNodes))
	for typeName := range nd.network.TypeNodes {
		typeNames = append(typeNames, typeName)
	}
	sort.Strings(typeNames)

	for _, typeName := range typeNames {
		node := nd.network.TypeNodes[typeName]
		fmt.Printf("   [T] type_%s\n", typeName)
		fmt.Printf("       │ Type: %s\n", typeName)
		fmt.Printf("       │ Enfants: %d nœuds\n", len(node.Children))
		fmt.Printf("       └─→ Propage tous les faits de type %s\n", typeName)
		fmt.Println()
	}
}

func (nd *NetworkDiagram) printAlphaNodes() {
	fmt.Println("┌" + strings.Repeat("─", 118) + "┐")
	fmt.Println("│ 2️⃣  ALPHA NODES (Filtres et calculs atomiques)                                                       │")
	fmt.Println("└" + strings.Repeat("─", 118) + "┘")
	fmt.Println()

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
		fmt.Printf("   📍 Variable: %s (%d nœuds)\n", varName, len(nodes))
		fmt.Println()

		// Trier les nœuds par ID pour un affichage cohérent
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].ID < nodes[j].ID
		})

		for i, node := range nodes {
			nd.printAlphaNodeDetail(node, i+1, len(nodes))
		}
		fmt.Println()
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

	fmt.Printf("      %s [α] %s\n", symbol, node.ID)

	// Extraire les détails de la condition
	if condMap, ok := node.Condition.(map[string]interface{}); ok {
		condType, _ := condMap["type"].(string)

		switch condType {
		case "passthrough":
			side, _ := condMap["side"].(string)
			fmt.Printf("         │ Type: PASSTHROUGH\n")
			if side != "" {
				fmt.Printf("         │ Side: %s\n", side)
			}
			fmt.Printf("         │ Opération: Propagation sans filtre\n")

		case "comparison":
			operator, _ := condMap["operator"].(string)
			left := nd.formatExpression(condMap["left"])
			right := nd.formatExpression(condMap["right"])
			fmt.Printf("         │ Type: COMPARISON\n")
			fmt.Printf("         │ Opérateur: %s\n", nd.symbolizeOperator(operator))
			fmt.Printf("         │ Expression: %s %s %s\n", left, nd.symbolizeOperator(operator), right)

		case "binaryOp":
			operator, _ := condMap["operator"].(string)
			left := nd.formatExpression(condMap["left"])
			right := nd.formatExpression(condMap["right"])
			fmt.Printf("         │ Type: BINARY OPERATION\n")
			fmt.Printf("         │ Opérateur: %s\n", nd.symbolizeOperator(operator))
			fmt.Printf("         │ Calcul: %s %s %s\n", left, nd.symbolizeOperator(operator), right)

		case "tempResult":
			stepName, _ := condMap["step_name"].(string)
			stepIdx, _ := condMap["step_idx"].(int)
			fmt.Printf("         │ Type: TEMP RESULT\n")
			fmt.Printf("         │ Step: %s (étape %d)\n", stepName, stepIdx)
			fmt.Printf("         │ Opération: Stockage résultat intermédiaire\n")

		default:
			fmt.Printf("         │ Type: %s\n", condType)
		}

		// Note: Les informations de partage sont affichées dans le résumé
	}

	fmt.Printf("         │ Enfants: %d\n", len(node.Children))
}

func (nd *NetworkDiagram) printPassthroughNodes() {
	fmt.Println("┌" + strings.Repeat("─", 118) + "┐")
	fmt.Println("│ 3️⃣  PASSTHROUGH NODES (Préparation pour jointure)                                                    │")
	fmt.Println("└" + strings.Repeat("─", 118) + "┘")
	fmt.Println()

	if len(nd.network.PassthroughRegistry) == 0 {
		fmt.Println("   (Aucun passthrough node)")
		fmt.Println()
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

	fmt.Println("   LEFT Side (tokens pour jointure gauche):")
	for _, key := range leftNodes {
		node := nd.network.PassthroughRegistry[key]
		fmt.Printf("      [⇒] %s\n", node.ID)
		fmt.Printf("          │ Rôle: Passthrough LEFT\n")
		fmt.Printf("          │ Enfants: %d\n", len(node.Children))
		fmt.Println()
	}

	fmt.Println("   RIGHT Side (tokens pour jointure droite):")
	for _, key := range rightNodes {
		node := nd.network.PassthroughRegistry[key]
		fmt.Printf("      [⇒] %s\n", node.ID)
		fmt.Printf("          │ Rôle: Passthrough RIGHT\n")
		fmt.Printf("          │ Enfants: %d\n", len(node.Children))
		fmt.Println()
	}
}

func (nd *NetworkDiagram) printJoinNodes() {
	fmt.Println("┌" + strings.Repeat("─", 118) + "┐")
	fmt.Println("│ 4️⃣  JOIN NODES (Jointures Beta)                                                                      │")
	fmt.Println("└" + strings.Repeat("─", 118) + "┘")
	fmt.Println()

	if len(nd.network.BetaNodes) == 0 {
		fmt.Println("   (Aucun join node)")
		fmt.Println()
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

		fmt.Printf("   %s [⋈] %s\n", symbol, joinNode.ID)
		fmt.Printf("      │ Type: JOIN NODE\n")
		fmt.Printf("      │ Variables LEFT: %v\n", joinNode.LeftVariables)
		fmt.Printf("      │ Variables RIGHT: %v\n", joinNode.RightVariables)

		// Afficher les JoinConditions
		if len(joinNode.JoinConditions) > 0 {
			fmt.Printf("      │ Conditions de jointure:\n")
			for _, jc := range joinNode.JoinConditions {
				fmt.Printf("      │   • %s.%s %s %s.%s\n",
					jc.LeftVar, jc.LeftField,
					nd.symbolizeOperator(jc.Operator),
					jc.RightVar, jc.RightField)
			}
		}

		// Vérifier si partagé
		rulesUsing := nd.findRulesUsingJoinNode(id)
		if len(rulesUsing) > 1 {
			fmt.Printf("      │ Partagé par: %d règles ✨ SHARED\n", len(rulesUsing))
			for _, rule := range rulesUsing {
				fmt.Printf("      │   - %s\n", rule)
			}
		} else if len(rulesUsing) == 1 {
			fmt.Printf("      │ Utilisé par: %s\n", rulesUsing[0])
		}

		fmt.Printf("      │ Enfants: %d\n", len(joinNode.Children))
		fmt.Println()
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

	fmt.Println("┌" + strings.Repeat("─", 118) + "┐")
	fmt.Println("│ 5️⃣  ROUTER NODES (Routage des tokens vers les règles)                                                │")
	fmt.Println("└" + strings.Repeat("─", 118) + "┘")
	fmt.Println()

	for _, router := range routers {
		fmt.Printf("   [🔀] %s\n", router.ID)
		fmt.Printf("       │ Type: RULE ROUTER\n")
		fmt.Printf("       │ Pour la règle: %s\n", router.RuleID)
		fmt.Printf("       │ Depuis JoinNode: %s\n", router.JoinNodeID)
		if router.TerminalNode != nil {
			fmt.Printf("       │ Vers TerminalNode: %s\n", router.TerminalNode.ID)
		}
		fmt.Println()
	}
}

func (nd *NetworkDiagram) printTerminalNodes() {
	fmt.Println("┌" + strings.Repeat("─", 118) + "┐")
	fmt.Println("│ 6️⃣  TERMINAL NODES (Actions)                                                                         │")
	fmt.Println("└" + strings.Repeat("─", 118) + "┘")
	fmt.Println()

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

		fmt.Printf("   %s [⚡] %s\n", symbol, terminal.ID)
		if terminal.Action != nil {
			fmt.Printf("      │ Action: %s\n", terminal.Action.Type)
		}
		fmt.Printf("      │ Tokens en mémoire: %d\n", len(terminal.Memory.Tokens))
		fmt.Println()
	}
}

func (nd *NetworkDiagram) printFlowDiagram() {
	printDiagramHeader("7️⃣  DIAGRAMME DE FLUX (Architecture complète)", diagramWidth)
	fmt.Println()

	printRulesExpression()
	printArchitectureDiagram()
	printDiagramLegend()
	printKeyPoints()
}

func (nd *NetworkDiagram) printSummary() {
	fmt.Println("┌" + strings.Repeat("─", 118) + "┐")
	fmt.Println("│ 📈 RÉSUMÉ DU RÉSEAU                                                                                  │")
	fmt.Println("└" + strings.Repeat("─", 118) + "┘")
	fmt.Println()

	uniqueJoins := make(map[string]bool)
	for _, node := range nd.network.BetaNodes {
		if joinNode, ok := node.(*JoinNode); ok {
			uniqueJoins[joinNode.ID] = true
		}
	}

	fmt.Printf("   Type Nodes:        %3d (routage par type)\n", len(nd.network.TypeNodes))
	fmt.Printf("   Alpha Nodes:       %3d (filtres et calculs atomiques)\n", len(nd.network.AlphaNodes))
	fmt.Printf("   Passthrough Nodes: %3d (préparation jointure)\n", len(nd.network.PassthroughRegistry))
	fmt.Printf("   Join Nodes:        %3d (jointures beta)\n", len(uniqueJoins))
	fmt.Printf("   Terminal Nodes:    %3d (actions)\n", len(nd.network.TerminalNodes))
	fmt.Println()

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

	fmt.Printf("   📊 Statistiques de partage:\n")
	fmt.Printf("      • AlphaNodes partagés:  %d / %d\n", sharedAlphaCount, len(nd.network.AlphaNodes))
	fmt.Printf("      • JoinNodes partagés:   %d / %d\n", sharedJoinCount, len(uniqueJoins))
	fmt.Println()
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
