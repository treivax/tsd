// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestArithmeticE2E_NetworkVisualization creates a comprehensive visualization
// of the RETE network structure for arithmetic alpha extraction, including
// analysis of node sharing between rules
func TestArithmeticE2E_NetworkVisualization(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "arithmetic_viz.tsd")
	// Multiple rules with different arithmetic conditions to test sharing
	content := `type Product(id: string, price: number, stock: number, weight: number)
type Order(id: string, productId: string, quantity: number)
action notify_expensive(productId: string)
action notify_heavy(productId: string)
action notify_low_stock(productId: string)
action notify_bulk(orderId: string)
action notify_match(productId: string, orderId: string)
// Rule 1: Expensive products (alpha condition on Product.price)
rule expensive_products : {p: Product} / p.price * 1.2 > 1000
    ==> notify_expensive(p.id)
// Rule 2: Heavy products (alpha condition on Product.weight)
rule heavy_products : {p: Product} / p.weight * 2.2 > 50
    ==> notify_heavy(p.id)
// Rule 3: Low stock (alpha condition on Product.stock)
rule low_stock : {p: Product} / p.stock < 10
    ==> notify_low_stock(p.id)
// Rule 4: Same condition as rule 1 (should share AlphaNode)
rule expensive_products_v2 : {p: Product} / p.price * 1.2 > 1000
    ==> notify_expensive(p.id)
// Rule 5: Bulk orders (alpha condition on Order.quantity)
rule bulk_orders : {o: Order} / o.quantity * 100 > 1000
    ==> notify_bulk(o.id)
// Rule 6: Join rule with alpha filters
rule expensive_bulk : {p: Product, o: Order} /
    p.id == o.productId AND
    p.price > 500 AND
    o.quantity > 5
    ==> notify_match(p.id, o.id)
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	// Print comprehensive network analysis
	printNetworkVisualization(t, network)
	// Analyze node sharing
	analyzeNodeSharing(t, network)
	// Print detailed node connections
	printNodeConnections(t, network)
	// Test actual behavior with facts
	testNetworkBehavior(t, network)
}

// printNetworkVisualization prints a visual representation of the network
func printNetworkVisualization(t *testing.T, network *ReteNetwork) {
	t.Log("\n" + strings.Repeat("=", 80))
	t.Log("RETE NETWORK VISUALIZATION")
	t.Log(strings.Repeat("=", 80))
	stats := network.GetNetworkStats()
	t.Logf("\n📊 NETWORK STATISTICS:")
	t.Logf("   • TypeNodes:     %d", stats["type_nodes"].(int))
	t.Logf("   • AlphaNodes:    %d", stats["alpha_nodes"].(int))
	t.Logf("   • BetaNodes:     %d", stats["beta_nodes"].(int))
	t.Logf("   • TerminalNodes: %d", stats["terminal_nodes"].(int))
	t.Log("\n🌳 NETWORK TOPOLOGY:")
	t.Log("   (Walking from TypeNodes to TerminalNodes)")
	t.Log("")
	// Walk each type node
	for typeName, typeNode := range network.TypeNodes {
		t.Logf("📦 TypeNode: %s", typeName)
		t.Logf("   ID: %s", typeNode.ID)
		t.Logf("   Children: %d", len(typeNode.GetChildren()))
		for i, child := range typeNode.GetChildren() {
			walkAndPrintNode(t, child, 1, i+1, len(typeNode.GetChildren()))
		}
		t.Log("")
	}
}

// walkAndPrintNode recursively walks and prints the node tree
func walkAndPrintNode(t *testing.T, node Node, depth int, childIndex int, totalChildren int) {
	indent := strings.Repeat("   ", depth)
	connector := "├──"
	if childIndex == totalChildren {
		connector = "└──"
	}
	nodeInfo := getNodeInfo(node)
	t.Logf("%s%s %s", indent, connector, nodeInfo)
	children := node.GetChildren()
	for i, child := range children {
		walkAndPrintNode(t, child, depth+1, i+1, len(children))
	}
}

// getNodeInfo returns detailed information about a node
func getNodeInfo(node Node) string {
	switch n := node.(type) {
	case *AlphaNode:
		// Determine if it's a passthrough or filter
		if cond, ok := n.Condition.(map[string]interface{}); ok {
			if condType, ok := cond["type"].(string); ok && condType == "passthrough" {
				return fmt.Sprintf("🔀 PassthroughAlpha: %s (side: %s)", n.ID, cond["side"])
			}
			// It's a filter node - extract condition info
			operator := ""
			if op, ok := cond["operator"].(string); ok {
				operator = op
			}
			varName := n.VariableName
			return fmt.Sprintf("🔍 AlphaNode[FILTER]: %s (var: %s, op: %s)", n.ID, varName, operator)
		}
		return fmt.Sprintf("🔍 AlphaNode: %s", n.ID)
	case *JoinNode:
		return fmt.Sprintf("⋈ JoinNode: %s (left: %v, right: %v)",
			n.ID, n.LeftVariables, n.RightVariables)
	case *TerminalNode:
		return fmt.Sprintf("🎯 TerminalNode: %s", n.ID)
	case *TypeNode:
		return fmt.Sprintf("📦 TypeNode: %s", n.ID)
	default:
		return fmt.Sprintf("❓ Unknown: %s", node.GetID())
	}
}

// analyzeNodeSharing analyzes which nodes are shared between rules
func analyzeNodeSharing(t *testing.T, network *ReteNetwork) {
	t.Log("\n" + strings.Repeat("=", 80))
	t.Log("NODE SHARING ANALYSIS")
	t.Log(strings.Repeat("=", 80))
	// Analyze AlphaNode sharing
	t.Log("\n🔍 ALPHANODE ANALYSIS:")
	// Group AlphaNodes by type
	filterNodes := make(map[string][]*AlphaNode)
	passthroughNodes := make(map[string][]*AlphaNode)
	for nodeID, alphaNode := range network.AlphaNodes {
		if cond, ok := alphaNode.Condition.(map[string]interface{}); ok {
			if condType, ok := cond["type"].(string); ok && condType == "passthrough" {
				passthroughNodes[nodeID] = append(passthroughNodes[nodeID], alphaNode)
			} else {
				filterNodes[nodeID] = append(filterNodes[nodeID], alphaNode)
			}
		}
	}
	t.Logf("   • Filter AlphaNodes: %d", len(filterNodes))
	t.Logf("   • Passthrough AlphaNodes: %d", len(passthroughNodes))
	// Analyze filter nodes in detail
	if len(filterNodes) > 0 {
		t.Log("\n   📋 Filter AlphaNode Details:")
		for nodeID, nodes := range filterNodes {
			for _, node := range nodes {
				children := node.GetChildren()
				rules := extractRulesFromChildren(children)
				t.Logf("      • %s", nodeID)
				t.Logf("        - Variable: %s", node.VariableName)
				if cond, ok := node.Condition.(map[string]interface{}); ok {
					if op, ok := cond["operator"].(string); ok {
						t.Logf("        - Operator: %s", op)
					}
				}
				t.Logf("        - Children: %d", len(children))
				t.Logf("        - Connected to rules: %v", rules)
				t.Logf("        - Memory size: %d facts", len(node.Memory.Facts))
			}
		}
	}
	// Analyze passthrough sharing
	if len(passthroughNodes) > 0 {
		t.Log("\n   📋 Passthrough AlphaNode Details (Per-Rule Isolation):")
		for nodeID := range passthroughNodes {
			// Extract rule name from passthrough ID
			parts := strings.Split(nodeID, "_")
			ruleName := "unknown"
			if len(parts) >= 2 {
				ruleName = parts[1]
			}
			t.Logf("      • %s (rule: %s)", nodeID, ruleName)
		}
	}
	// Analyze TypeNode sharing
	t.Log("\n📦 TYPENODE SHARING:")
	for typeName, typeNode := range network.TypeNodes {
		children := typeNode.GetChildren()
		rules := extractRulesFromChildren(children)
		t.Logf("   • %s: Shared by %d rules", typeName, len(rules))
		if len(rules) > 0 {
			t.Logf("     Rules: %v", rules)
		}
	}
	// Check for identical condition sharing
	t.Log("\n🔄 IDENTICAL CONDITION SHARING:")
	conditionGroups := groupNodesByCondition(network)
	if len(conditionGroups) > 0 {
		for condHash, nodes := range conditionGroups {
			if len(nodes) > 1 {
				t.Logf("   ⚠️  Condition hash %s: %d nodes with same condition", condHash[:8], len(nodes))
				for _, nodeID := range nodes {
					t.Logf("      - %s", nodeID)
				}
				t.Log("      💡 Opportunity: These nodes could potentially be shared")
			}
		}
	} else {
		t.Log("   ✅ No duplicate conditions found (all unique)")
	}
}

// extractRulesFromChildren extracts rule names from terminal nodes in children
func extractRulesFromChildren(children []Node) []string {
	rules := make(map[string]bool)
	var walk func(Node)
	walk = func(node Node) {
		if terminal, ok := node.(*TerminalNode); ok {
			// Extract rule name from terminal ID
			ruleName := strings.TrimSuffix(terminal.ID, "_terminal")
			rules[ruleName] = true
		}
		for _, child := range node.GetChildren() {
			walk(child)
		}
	}
	for _, child := range children {
		walk(child)
	}
	result := make([]string, 0, len(rules))
	for rule := range rules {
		result = append(result, rule)
	}
	return result
}

// groupNodesByCondition groups nodes that have identical conditions
func groupNodesByCondition(network *ReteNetwork) map[string][]string {
	groups := make(map[string][]string)
	for nodeID, alphaNode := range network.AlphaNodes {
		if cond, ok := alphaNode.Condition.(map[string]interface{}); ok {
			condHash := fmt.Sprintf("%v", cond)
			groups[condHash] = append(groups[condHash], nodeID)
		}
	}
	return groups
}

// printNodeConnections prints detailed connection information
func printNodeConnections(t *testing.T, network *ReteNetwork) {
	t.Log("\n" + strings.Repeat("=", 80))
	t.Log("DETAILED NODE CONNECTIONS")
	t.Log(strings.Repeat("=", 80))
	// AlphaNodes connections
	t.Log("\n🔍 ALPHANODE CONNECTIONS:")
	for nodeID, alphaNode := range network.AlphaNodes {
		children := alphaNode.GetChildren()
		if len(children) == 0 {
			continue
		}
		t.Logf("\n   %s:", nodeID)
		t.Logf("      Variable: %s", alphaNode.VariableName)
		t.Logf("      Children: %d", len(children))
		for _, child := range children {
			t.Logf("         → %s", child.GetID())
		}
	}
	// JoinNodes connections
	if len(network.BetaNodes) > 0 {
		t.Log("\n⋈ JOINNODE CONNECTIONS:")
		for nodeID, betaNode := range network.BetaNodes {
			if joinNode, ok := betaNode.(*JoinNode); ok {
				children := joinNode.GetChildren()
				t.Logf("\n   %s:", nodeID)
				t.Logf("      Left vars: %v", joinNode.LeftVariables)
				t.Logf("      Right vars: %v", joinNode.RightVariables)
				t.Logf("      Children: %d", len(children))
				for _, child := range children {
					t.Logf("         → %s", child.GetID())
				}
			}
		}
	}
	// TerminalNodes
	t.Log("\n🎯 TERMINAL NODES:")
	for nodeID, terminal := range network.TerminalNodes {
		t.Logf("   %s:", nodeID)
		t.Logf("      Activations: %d", len(terminal.GetMemory().Tokens))
	}
}

// testNetworkBehavior tests the actual behavior with sample facts
func testNetworkBehavior(t *testing.T, network *ReteNetwork) {
	t.Log("\n" + strings.Repeat("=", 80))
	t.Log("NETWORK BEHAVIOR TEST")
	t.Log(strings.Repeat("=", 80))
	// Test data
	facts := []*Fact{
		// Should match: expensive_products, expensive_products_v2 (price * 1.2 = 1200 > 1000)
		{
			ID:   "P1",
			Type: "Product",
			Fields: map[string]interface{}{
				"id":     "p1",
				"price":  1000.0,
				"stock":  5.0,
				"weight": 20.0,
			},
		},
		// Should match: heavy_products (weight * 2.2 = 55 > 50), low_stock (5 < 10)
		{
			ID:   "P2",
			Type: "Product",
			Fields: map[string]interface{}{
				"id":     "p2",
				"price":  500.0,
				"stock":  5.0,
				"weight": 25.0,
			},
		},
		// Should match: bulk_orders (quantity * 100 = 1500 > 1000)
		{
			ID:   "O1",
			Type: "Order",
			Fields: map[string]interface{}{
				"id":        "o1",
				"productId": "p1",
				"quantity":  15.0,
			},
		},
		// Should match: bulk_orders (quantity * 100 = 1000 = 1000, edge case)
		{
			ID:   "O2",
			Type: "Order",
			Fields: map[string]interface{}{
				"id":        "o2",
				"productId": "p2",
				"quantity":  10.0,
			},
		},
	}
	t.Log("\n📤 SUBMITTING FACTS:")
	for _, fact := range facts {
		t.Logf("   • %s (%s): %+v", fact.ID, fact.Type, fact.Fields)
		if err := network.SubmitFact(fact); err != nil {
			t.Errorf("Failed to submit fact %s: %v", fact.ID, err)
		}
	}
	t.Log("\n🎯 RULE ACTIVATIONS:")
	totalActivations := 0
	for ruleName, terminal := range network.TerminalNodes {
		tokens := terminal.GetMemory().Tokens
		activationCount := len(tokens)
		totalActivations += activationCount
		if activationCount > 0 {
			t.Logf("   ✅ %s: %d activation(s)",
				strings.TrimSuffix(ruleName, "_terminal"), activationCount)
			// Show which facts triggered the rule
			for _, token := range tokens {
				factIDs := make([]string, 0, len(token.Facts))
				for _, fact := range token.Facts {
					factIDs = append(factIDs, fact.ID)
				}
				t.Logf("      Facts: %v", factIDs)
			}
		} else {
			t.Logf("   ⭕ %s: no activations",
				strings.TrimSuffix(ruleName, "_terminal"))
		}
	}
	t.Logf("\n📊 TOTAL ACTIVATIONS: %d", totalActivations)
	// Verify expected behavior
	t.Log("\n✓ VERIFICATION:")
	expectedActivations := map[string]int{
		"expensive_products_terminal":    1, // P1
		"expensive_products_v2_terminal": 1, // P1 (same condition)
		"heavy_products_terminal":        1, // P2
		"low_stock_terminal":             2, // P1 and P2
		"bulk_orders_terminal":           1, // O1 only (quantity * 100 = 1500 > 1000), O2 has 10 * 100 = 1000 NOT > 1000
		"expensive_bulk_terminal":        1, // P1 + O1
	}
	allCorrect := true
	for terminalID, expected := range expectedActivations {
		if terminal, exists := network.TerminalNodes[terminalID]; exists {
			actual := len(terminal.GetMemory().Tokens)
			if actual != expected {
				t.Errorf("   ❌ %s: expected %d, got %d",
					strings.TrimSuffix(terminalID, "_terminal"), expected, actual)
				allCorrect = false
			} else {
				t.Logf("   ✅ %s: correct (%d activations)",
					strings.TrimSuffix(terminalID, "_terminal"), expected)
			}
		}
	}
	if allCorrect {
		t.Log("\n🎉 ALL VERIFICATIONS PASSED!")
	}
}

// TestArithmeticE2E_SharingOpportunities identifies potential sharing opportunities
func TestArithmeticE2E_SharingOpportunities(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "sharing.tsd")
	// Rules with intentional duplication to test sharing detection
	content := `type Item(id: string, value: number)
action log(msg: string)
rule expensive_1 : {i: Item} / i.value * 1.5 > 100 ==> log("Expensive 1")
rule expensive_2 : {i: Item} / i.value * 1.5 > 100 ==> log("Expensive 2")
rule expensive_3 : {i: Item} / i.value * 1.5 > 100 ==> log("Expensive 3")
rule cheap_1 : {i: Item} / i.value * 0.5 < 50 ==> log("Cheap 1")
rule cheap_2 : {i: Item} / i.value * 0.5 < 50 ==> log("Cheap 2")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	t.Log("\n" + strings.Repeat("=", 80))
	t.Log("SHARING OPPORTUNITIES ANALYSIS")
	t.Log(strings.Repeat("=", 80))
	// Count AlphaNodes with same conditions
	conditionMap := make(map[string][]string)
	for nodeID, alphaNode := range network.AlphaNodes {
		if cond, ok := alphaNode.Condition.(map[string]interface{}); ok {
			if condType, ok := cond["type"].(string); ok && condType != "passthrough" {
				condKey := fmt.Sprintf("%v", cond)
				conditionMap[condKey] = append(conditionMap[condKey], nodeID)
			}
		}
	}
	t.Log("\n🔍 DUPLICATE CONDITIONS FOUND:")
	duplicateCount := 0
	for _, nodeIDs := range conditionMap {
		if len(nodeIDs) > 1 {
			duplicateCount++
			t.Logf("\n   Duplicate set #%d: %d nodes with same condition", duplicateCount, len(nodeIDs))
			for _, nodeID := range nodeIDs {
				t.Logf("      • %s", nodeID)
			}
			t.Log("      💡 Optimization: These nodes could be merged into one shared AlphaNode")
			t.Logf("      💾 Memory savings: %d duplicate nodes eliminated", len(nodeIDs)-1)
		}
	}
	if duplicateCount == 0 {
		t.Log("   ✅ No duplicate conditions found")
	} else {
		t.Logf("\n📊 SUMMARY:")
		t.Logf("   • Duplicate condition sets: %d", duplicateCount)
		t.Logf("   • Total AlphaNodes: %d", len(network.AlphaNodes))
		t.Logf("   • Potential sharing ratio: %.1f%%",
			float64(duplicateCount*100)/float64(len(network.AlphaNodes)))
	}
}
