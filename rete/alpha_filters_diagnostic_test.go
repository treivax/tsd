// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"os"
	"path/filepath"
	"testing"
)

// TestAlphaFiltersDiagnostic_JoinRules verifies that alpha filters are created
// for single-variable conditions in join rules
func TestAlphaFiltersDiagnostic_JoinRules(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "test.tsd")
	// Two join rules with different alpha filters on the same variable
	content := `type Person(#id: string, age: number)
type Order(#id: string, personId: string, amount: number)
action print(msg: string)
rule large_orders : {p: Person, o: Order} / p.id == o.personId AND o.amount > 100
    ==> print("Large order")
rule very_large_orders : {p: Person, o: Order} / p.id == o.personId AND o.amount > 500
    ==> print("Very large order")
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
	// Diagnostic: Print network structure
	t.Log("=== NETWORK STRUCTURE ===")
	// Check TypeNodes
	t.Logf("TypeNodes: %d", len(network.TypeNodes))
	for typeName, typeNode := range network.TypeNodes {
		t.Logf("  Type: %s", typeName)
		t.Logf("    Children: %d", len(typeNode.GetChildren()))
		for _, child := range typeNode.GetChildren() {
			t.Logf("      -> %s (type: %T)", child.GetID(), child)
		}
	}
	// Check for AlphaNodes in the network
	t.Log("\n=== CHECKING FOR ALPHA FILTERS ===")
	alphaCount := 0
	passthroughCount := 0
	// Walk the network from TypeNodes
	for typeName, typeNode := range network.TypeNodes {
		t.Logf("Walking from TypeNode: %s", typeName)
		walkNode(t, typeNode, 1, &alphaCount, &passthroughCount)
	}
	t.Logf("\nTotal AlphaNodes found: %d", alphaCount)
	t.Logf("Total PassthroughNodes found: %d", passthroughCount)
	stats := network.GetNetworkStats()
	t.Logf("\nNetwork stats: %+v", stats)
	// Expected: We should have AlphaNodes for the alpha filters
	// large_orders: o.amount > 100
	// very_large_orders: o.amount > 500
	// So we expect at least 2 alpha filter nodes
	reportedAlphaNodes := stats["alpha_nodes"].(int)
	if reportedAlphaNodes < 2 {
		t.Errorf("Expected at least 2 AlphaNodes for alpha filters, got %d", reportedAlphaNodes)
		t.Log("This suggests alpha filters are NOT being created for single-variable conditions in join rules")
	}
	// Test behavior: submit facts and verify correct filtering
	t.Log("\n=== TESTING FACT FILTERING ===")
	person := &Fact{
		ID:   "P1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":  "p1",
			"age": 30,
		},
	}
	// order1: amount=150 should match large_orders ONLY
	order1 := &Fact{
		ID:   "O1",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":       "o1",
			"personId": "p1",
			"amount":   150.0,
		},
	}
	// order2: amount=600 should match BOTH rules
	order2 := &Fact{
		ID:   "O2",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":       "o2",
			"personId": "p1",
			"amount":   600.0,
		},
	}
	network.SubmitFact(person)
	network.SubmitFact(order1)
	network.SubmitFact(order2)
	// Check activations
	largeTokens := network.TerminalNodes["large_orders_terminal"].GetMemory().Tokens
	veryLargeTokens := network.TerminalNodes["very_large_orders_terminal"].GetMemory().Tokens
	t.Logf("\nlarge_orders activations: %d", len(largeTokens))
	t.Logf("very_large_orders activations: %d", len(veryLargeTokens))
	// Expected:
	// - large_orders: 2 activations (order1 and order2 both > 100)
	// - very_large_orders: 1 activation (only order2 > 500)
	if len(largeTokens) != 2 {
		t.Errorf("Expected 2 activations for large_orders, got %d", len(largeTokens))
	}
	if len(veryLargeTokens) != 1 {
		t.Errorf("Expected 1 activation for very_large_orders, got %d", len(veryLargeTokens))
	}
}

// walkNode recursively walks the network and counts node types
func walkNode(t *testing.T, node Node, depth int, alphaCount, passthroughCount *int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	nodeType := "Unknown"
	switch n := node.(type) {
	case *AlphaNode:
		nodeType = "AlphaNode"
		// Check if it's a passthrough or a real filter
		if cond, ok := n.Condition.(map[string]interface{}); ok {
			if condType, ok := cond["type"].(string); ok && condType == "passthrough" {
				nodeType = "PassthroughAlpha"
				*passthroughCount++
			} else {
				*alphaCount++
				t.Logf("%s[AlphaNode] %s - Condition: %+v", indent, n.ID, n.Condition)
			}
		}
	case *JoinNode:
		nodeType = "JoinNode"
	case *TerminalNode:
		nodeType = "TerminalNode"
	case *TypeNode:
		nodeType = "TypeNode"
	}
	t.Logf("%s%s: %s", indent, nodeType, node.GetID())
	// Recurse to children
	for _, child := range node.GetChildren() {
		walkNode(t, child, depth+1, alphaCount, passthroughCount)
	}
}
