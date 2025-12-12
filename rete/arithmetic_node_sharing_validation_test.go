// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"fmt"
	"testing"
)
// TestArithmeticDecomposition_NodeSharingValidation validates that identical
// arithmetic sub-expressions in decomposed chains are properly shared across rules
func TestArithmeticDecomposition_NodeSharingValidation(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	orderType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Order",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "qte", Type: "number"},
		},
	}
	typeNode := NewTypeNode("Order", orderType, storage)
	network.TypeNodes["Order"] = typeNode
	// Two rules with identical sub-expressions:
	// Rule 1: (c.qte * 23 - 10) > 100
	// Rule 2: (c.qte * 23 - 10) < 50
	//
	// Expected sharing:
	// - Step 1: c.qte * 23 → temp_1 (SHARED)
	// - Step 2: temp_1 - 10 → temp_2 (SHARED)
	// - Step 3a: temp_2 > 100 (Rule 1 only)
	// - Step 3b: temp_2 < 50 (Rule 2 only)
	condition1 := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "-",
			"left": map[string]interface{}{
				"type":     "binaryOp",
				"operator": "*",
				"left":     map[string]interface{}{"type": "fieldAccess", "field": "qte"},
				"right":    map[string]interface{}{"type": "number", "value": 23},
			},
			"right": map[string]interface{}{"type": "number", "value": 10},
		},
		"right": map[string]interface{}{"type": "number", "value": 100},
	}
	condition2 := map[string]interface{}{
		"type":     "comparison",
		"operator": "<",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "-",
			"left": map[string]interface{}{
				"type":     "binaryOp",
				"operator": "*",
				"left":     map[string]interface{}{"type": "fieldAccess", "field": "qte"},
				"right":    map[string]interface{}{"type": "number", "value": 23},
			},
			"right": map[string]interface{}{"type": "number", "value": 10},
		},
		"right": map[string]interface{}{"type": "number", "value": 50},
	}
	decomposer := NewArithmeticExpressionDecomposer()
	chainBuilder := NewAlphaChainBuilder(network, storage)
	// Build Rule 1
	steps1, err := decomposer.DecomposeToDecomposedConditions(condition1)
	if err != nil {
		t.Fatalf("Decompose condition1 failed: %v", err)
	}
	chain1, err := chainBuilder.BuildDecomposedChain(steps1, "c", typeNode, "rule1")
	if err != nil {
		t.Fatalf("BuildDecomposedChain rule1 failed: %v", err)
	}
	t.Logf("Chain 1: %d nodes", len(chain1.Nodes))
	for i, node := range chain1.Nodes {
		t.Logf("  Node %d: %s (%s)", i, node.ID, node.ResultName)
	}
	// Build Rule 2
	steps2, err := decomposer.DecomposeToDecomposedConditions(condition2)
	if err != nil {
		t.Fatalf("Decompose condition2 failed: %v", err)
	}
	chain2, err := chainBuilder.BuildDecomposedChain(steps2, "c", typeNode, "rule2")
	if err != nil {
		t.Fatalf("BuildDecomposedChain rule2 failed: %v", err)
	}
	t.Logf("Chain 2: %d nodes", len(chain2.Nodes))
	for i, node := range chain2.Nodes {
		t.Logf("  Node %d: %s (%s)", i, node.ID, node.ResultName)
	}
	// VALIDATION 1: Step 1 (c.qte * 23) should be shared
	if chain1.Nodes[0].ID != chain2.Nodes[0].ID {
		t.Errorf("❌ Step 1 (c.qte * 23) NOT shared: rule1=%s, rule2=%s",
			chain1.Nodes[0].ID, chain2.Nodes[0].ID)
	} else {
		t.Logf("✅ Step 1 (c.qte * 23) SHARED: %s", chain1.Nodes[0].ID)
	}
	// VALIDATION 2: Step 2 (temp_1 - 10) should be shared
	if chain1.Nodes[1].ID != chain2.Nodes[1].ID {
		t.Errorf("❌ Step 2 (temp_1 - 10) NOT shared: rule1=%s, rule2=%s",
			chain1.Nodes[1].ID, chain2.Nodes[1].ID)
	} else {
		t.Logf("✅ Step 2 (temp_1 - 10) SHARED: %s", chain1.Nodes[1].ID)
	}
	// VALIDATION 3: Step 3 (final comparison) should be DIFFERENT
	if chain1.Nodes[2].ID == chain2.Nodes[2].ID {
		t.Errorf("❌ Step 3 (comparison) should NOT be shared: both=%s", chain1.Nodes[2].ID)
	} else {
		t.Logf("✅ Step 3 (comparison) NOT shared (expected):")
		t.Logf("   Rule1: %s (temp_2 > 100)", chain1.Nodes[2].ID)
		t.Logf("   Rule2: %s (temp_2 < 50)", chain2.Nodes[2].ID)
	}
	// VALIDATION 4: Shared nodes have correct child count
	sharedNode1 := chain1.Nodes[0]
	if len(sharedNode1.Children) != 1 {
		t.Errorf("❌ Shared node 1 should have 1 child (step 2), got %d", len(sharedNode1.Children))
	} else {
		t.Logf("✅ Shared node 1 has correct child count: %d", len(sharedNode1.Children))
	}
	sharedNode2 := chain1.Nodes[1]
	if len(sharedNode2.Children) != 2 {
		t.Errorf("❌ Shared node 2 should have 2 children (both rule final steps), got %d", len(sharedNode2.Children))
	} else {
		t.Logf("✅ Shared node 2 has correct child count: %d (branches to both rules)", len(sharedNode2.Children))
	}
	// VALIDATION 5: Test execution with shared nodes
	// Note: Each rule evaluation gets its own context, but the shared nodes are reused
	testCases := []struct {
		name            string
		qte             int
		expectedCalc    float64
		rule1ShouldPass bool // > 100
		rule2ShouldPass bool // < 50
	}{
		{"very_low", 1, 13, false, true}, // 1*23-10 = 13 < 50 (rule2 only)
		{"low", 3, 59, false, false},     // 3*23-10 = 59 (neither)
		{"medium", 5, 105, true, false},  // 5*23-10 = 105 > 100 (rule1 only)
		{"high", 10, 220, true, false},   // 10*23-10 = 220 > 100 (rule1 only)
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fact := &Fact{
				ID:   tc.name,
				Type: "Order",
				Fields: map[string]interface{}{
					"qte": tc.qte,
				},
			}
			// Test Rule 1 evaluation
			evaluator := NewConditionEvaluator(storage)
			// Manually evaluate chain 1 to verify correctness
			ctx1 := NewEvaluationContext(fact)
			for _, node := range chain1.Nodes {
				result, err := evaluator.EvaluateWithContext(node.Condition, fact, ctx1)
				if err != nil {
					t.Fatalf("Evaluate node %s failed: %v", node.ID, err)
				}
				ctx1.SetIntermediateResult(node.ResultName, result)
			}
			// Check final result for rule 1
			finalResult1, exists := ctx1.GetIntermediateResult("temp_3")
			if !exists {
				t.Fatal("temp_3 (rule1) not found in context")
			}
			boolResult1 := finalResult1.(bool)
			if boolResult1 != tc.rule1ShouldPass {
				t.Errorf("Rule1 (> 100): expected %v, got %v (calc=%.0f)", tc.rule1ShouldPass, boolResult1, tc.expectedCalc)
			}
			// Test Rule 2 evaluation (separate context)
			ctx2 := NewEvaluationContext(fact)
			for _, node := range chain2.Nodes {
				result, err := evaluator.EvaluateWithContext(node.Condition, fact, ctx2)
				if err != nil {
					t.Fatalf("Evaluate node %s failed: %v", node.ID, err)
				}
				ctx2.SetIntermediateResult(node.ResultName, result)
			}
			// Check final result for rule 2
			finalResult2, exists := ctx2.GetIntermediateResult("temp_3")
			if !exists {
				t.Fatal("temp_3 (rule2) not found in context")
			}
			boolResult2 := finalResult2.(bool)
			if boolResult2 != tc.rule2ShouldPass {
				t.Errorf("Rule2 (< 50): expected %v, got %v (calc=%.0f)", tc.rule2ShouldPass, boolResult2, tc.expectedCalc)
			}
			// Verify intermediate calculation is correct
			temp2_rule1, _ := ctx1.GetIntermediateResult("temp_2")
			if temp2_rule1 != tc.expectedCalc {
				t.Errorf("Rule1 temp_2: expected %.0f, got %.0f", tc.expectedCalc, temp2_rule1)
			}
			temp2_rule2, _ := ctx2.GetIntermediateResult("temp_2")
			if temp2_rule2 != tc.expectedCalc {
				t.Errorf("Rule2 temp_2: expected %.0f, got %.0f", tc.expectedCalc, temp2_rule2)
			}
			t.Logf("✅ qte=%d → %.0f: Rule1=%v, Rule2=%v (shared nodes, separate contexts)",
				tc.qte, tc.expectedCalc, boolResult1, boolResult2)
		})
	}
	t.Log("✅ Node sharing in decomposed chains validated successfully")
}
// TestArithmeticDecomposition_MultiRuleSharing tests sharing across 3+ rules
func TestArithmeticDecomposition_MultiRuleSharing(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	orderType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Order",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "qte", Type: "number"},
		},
	}
	typeNode := NewTypeNode("Order", orderType, storage)
	network.TypeNodes["Order"] = typeNode
	// Three rules sharing the same base expression:
	// Rule 1: (c.qte * 23) > 50
	// Rule 2: (c.qte * 23) > 100
	// Rule 3: (c.qte * 23) < 30
	//
	// All should share: c.qte * 23 → temp_1
	baseExpr := map[string]interface{}{
		"type":     "binaryOp",
		"operator": "*",
		"left":     map[string]interface{}{"type": "fieldAccess", "field": "qte"},
		"right":    map[string]interface{}{"type": "number", "value": 23},
	}
	conditions := []map[string]interface{}{
		{
			"type":     "comparison",
			"operator": ">",
			"left":     baseExpr,
			"right":    map[string]interface{}{"type": "number", "value": 50},
		},
		{
			"type":     "comparison",
			"operator": ">",
			"left":     baseExpr,
			"right":    map[string]interface{}{"type": "number", "value": 100},
		},
		{
			"type":     "comparison",
			"operator": "<",
			"left":     baseExpr,
			"right":    map[string]interface{}{"type": "number", "value": 30},
		},
	}
	decomposer := NewArithmeticExpressionDecomposer()
	chainBuilder := NewAlphaChainBuilder(network, storage)
	chains := make([]*AlphaChain, 0, len(conditions))
	// Build chains for all rules
	for i, cond := range conditions {
		steps, err := decomposer.DecomposeToDecomposedConditions(cond)
		if err != nil {
			t.Fatalf("Decompose condition %d failed: %v", i+1, err)
		}
		chain, err := chainBuilder.BuildDecomposedChain(steps, "c", typeNode, fmt.Sprintf("rule%d", i+1))
		if err != nil {
			t.Fatalf("BuildDecomposedChain rule%d failed: %v", i+1, err)
		}
		chains = append(chains, chain)
		t.Logf("Rule %d: %d nodes", i+1, len(chain.Nodes))
	}
	// VALIDATION: All three rules share the first step
	firstNodeID := chains[0].Nodes[0].ID
	for i := 1; i < len(chains); i++ {
		if chains[i].Nodes[0].ID != firstNodeID {
			t.Errorf("❌ Rule %d does NOT share first node: expected %s, got %s",
				i+1, firstNodeID, chains[i].Nodes[0].ID)
		}
	}
	t.Logf("✅ All 3 rules share first node: %s", firstNodeID)
	// VALIDATION: The shared node has 3 children (one for each rule's final comparison)
	sharedNode := chains[0].Nodes[0]
	if len(sharedNode.Children) != 3 {
		t.Errorf("❌ Shared node should have 3 children, got %d", len(sharedNode.Children))
	} else {
		t.Logf("✅ Shared node branches to 3 different comparisons")
	}
	// VALIDATION: Calculate sharing statistics
	allNodeIDs := make([]string, 0)
	for _, chain := range chains {
		for _, node := range chain.Nodes {
			allNodeIDs = append(allNodeIDs, node.ID)
		}
	}
	uniqueNodes := make(map[string]bool)
	for _, id := range allNodeIDs {
		uniqueNodes[id] = true
	}
	t.Logf("\n=== SHARING STATISTICS ===")
	t.Logf("Total references: %d", len(allNodeIDs))
	t.Logf("Unique nodes: %d", len(uniqueNodes))
	t.Logf("Shared instances: %d", len(allNodeIDs)-len(uniqueNodes))
	t.Logf("Efficiency: %.1f%%", float64(len(uniqueNodes))/float64(len(allNodeIDs))*100)
	if len(allNodeIDs)-len(uniqueNodes) < 2 {
		t.Errorf("❌ Expected at least 2 shared instances (step1 shared 3 times)")
	} else {
		t.Log("✅ Sharing working correctly for multi-rule scenario")
	}
}
// TestArithmeticDecomposition_PartialSharing tests rules with partial overlap
func TestArithmeticDecomposition_PartialSharing(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	orderType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Order",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "qte", Type: "number"},
		},
	}
	typeNode := NewTypeNode("Order", orderType, storage)
	network.TypeNodes["Order"] = typeNode
	// Rule 1: (c.qte * 23 - 10) > 100
	// Rule 2: (c.qte * 23 + 5) > 100
	//
	// Expected:
	// - Step 1: c.qte * 23 → temp_1 (SHARED)
	// - Step 2a: temp_1 - 10 → temp_2 (Rule 1 only)
	// - Step 2b: temp_1 + 5 → temp_2 (Rule 2 only)
	// - Step 3: temp_2 > 100 (both, but different temp_2 values)
	condition1 := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "-",
			"left": map[string]interface{}{
				"type":     "binaryOp",
				"operator": "*",
				"left":     map[string]interface{}{"type": "fieldAccess", "field": "qte"},
				"right":    map[string]interface{}{"type": "number", "value": 23},
			},
			"right": map[string]interface{}{"type": "number", "value": 10},
		},
		"right": map[string]interface{}{"type": "number", "value": 100},
	}
	condition2 := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "+",
			"left": map[string]interface{}{
				"type":     "binaryOp",
				"operator": "*",
				"left":     map[string]interface{}{"type": "fieldAccess", "field": "qte"},
				"right":    map[string]interface{}{"type": "number", "value": 23},
			},
			"right": map[string]interface{}{"type": "number", "value": 5},
		},
		"right": map[string]interface{}{"type": "number", "value": 100},
	}
	decomposer := NewArithmeticExpressionDecomposer()
	steps1, _ := decomposer.DecomposeToDecomposedConditions(condition1)
	steps2, _ := decomposer.DecomposeToDecomposedConditions(condition2)
	chainBuilder := NewAlphaChainBuilder(network, storage)
	chain1, _ := chainBuilder.BuildDecomposedChain(steps1, "c", typeNode, "rule1")
	chain2, _ := chainBuilder.BuildDecomposedChain(steps2, "c", typeNode, "rule2")
	// VALIDATION: Verify partial sharing
	if chain1.Nodes[0].ID != chain2.Nodes[0].ID {
		t.Errorf("❌ Step 1 (c.qte * 23) should be shared")
	} else {
		t.Logf("✅ Step 1 shared: %s", chain1.Nodes[0].ID)
	}
	if chain1.Nodes[1].ID == chain2.Nodes[1].ID {
		t.Errorf("❌ Step 2 should NOT be shared (different operations: - vs +)")
	} else {
		t.Logf("✅ Step 2 NOT shared (different ops):")
		t.Logf("   Rule1: %s (temp_1 - 10)", chain1.Nodes[1].ID)
		t.Logf("   Rule2: %s (temp_1 + 5)", chain2.Nodes[1].ID)
	}
	// The shared node should have 2 children (branches to both step2 variants)
	sharedNode := chain1.Nodes[0]
	if len(sharedNode.Children) != 2 {
		t.Errorf("❌ Shared node should have 2 children, got %d", len(sharedNode.Children))
	} else {
		t.Logf("✅ Shared node branches to 2 different step2 nodes")
	}
	t.Log("✅ Partial sharing validated")
}