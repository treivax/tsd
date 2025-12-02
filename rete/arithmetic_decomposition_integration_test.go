// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// TestArithmeticDecomposition_IntegrationSimple tests decomposition with context evaluation
func TestArithmeticDecomposition_IntegrationSimple(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Define Order type
	orderType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Order",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "qte", Type: "number"},
		},
	}

	// Create TypeNode
	typeNode := NewTypeNode("Order", orderType, storage)
	network.TypeNodes["Order"] = typeNode

	// Create complex arithmetic expression: (c.qte * 23 - 10) > 100
	// This should decompose into:
	// Step 1: c.qte * 23 → temp_1
	// Step 2: temp_1 - 10 → temp_2
	// Step 3: temp_2 > 100 → temp_3 (final)

	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "-",
			"left": map[string]interface{}{
				"type":     "binaryOp",
				"operator": "*",
				"left": map[string]interface{}{
					"type":  "fieldAccess",
					"field": "qte",
				},
				"right": map[string]interface{}{
					"type":  "number",
					"value": 23,
				},
			},
			"right": map[string]interface{}{
				"type":  "number",
				"value": 10,
			},
		},
		"right": map[string]interface{}{
			"type":  "number",
			"value": 100,
		},
	}

	// Decompose the expression
	decomposer := NewArithmeticExpressionDecomposer()
	decomposedSteps, err := decomposer.DecomposeToDecomposedConditions(condition)
	if err != nil {
		t.Fatalf("Decomposition failed: %v", err)
	}

	t.Logf("✅ Decomposed into %d steps", len(decomposedSteps))
	for i, step := range decomposedSteps {
		t.Logf("  Step %d: %s (deps: %v)", i+1, step.ResultName, step.Dependencies)
	}

	// Verify decomposition
	if len(decomposedSteps) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(decomposedSteps))
	}

	// Build the decomposed chain
	chainBuilder := NewAlphaChainBuilder(network, storage)
	alphaChain, err := chainBuilder.BuildDecomposedChain(
		decomposedSteps,
		"c",
		typeNode,
		"test_rule",
	)
	if err != nil {
		t.Fatalf("BuildDecomposedChain failed: %v", err)
	}

	t.Logf("✅ Alpha chain built with %d nodes", len(alphaChain.Nodes))

	// Verify metadata on nodes
	for i, node := range alphaChain.Nodes {
		if node.ResultName == "" {
			t.Errorf("Node %d missing ResultName", i)
		}
		if !node.IsAtomic {
			t.Errorf("Node %d not marked as atomic", i)
		}
		if i == 0 && len(node.Dependencies) != 0 {
			t.Errorf("Node 0 should have no dependencies, got %v", node.Dependencies)
		}
		if i > 0 && len(node.Dependencies) == 0 {
			t.Errorf("Node %d should have dependencies", i)
		}
		t.Logf("  Node %d: %s (atomic=%v, deps=%v)", i, node.ResultName, node.IsAtomic, node.Dependencies)
	}

	// Test evaluation with context
	testCases := []struct {
		name       string
		qte        int
		shouldPass bool
		expected   float64 // qte * 23 - 10
	}{
		{"below_threshold", 4, false, 82},  // 4*23-10 = 82 < 100
		{"at_threshold", 5, true, 105},     // 5*23-10 = 105 > 100
		{"above_threshold", 10, true, 220}, // 10*23-10 = 220 > 100
		{"high_quantity", 20, true, 450},   // 20*23-10 = 450 > 100
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create fact
			fact := &Fact{
				ID:   tc.name,
				Type: "Order",
				Fields: map[string]interface{}{
					"qte": tc.qte,
				},
			}

			// Create evaluation context
			ctx := NewEvaluationContext(fact)

			t.Logf("Testing: qte=%d, expected_calc=%.0f, threshold=100", tc.qte, tc.expected)

			// Manually evaluate the chain step by step
			// Step 1: c.qte * 23 → temp_1
			node1 := alphaChain.Nodes[0]
			evaluator := NewConditionEvaluator(storage)
			result1, err := evaluator.EvaluateWithContext(node1.Condition, fact, ctx)
			if err != nil {
				t.Fatalf("Step 1 evaluation failed: %v", err)
			}
			ctx.SetIntermediateResult(node1.ResultName, result1)
			t.Logf("  Step 1: %s = %.0f", node1.ResultName, result1)

			// Step 2: temp_1 - 10 → temp_2
			node2 := alphaChain.Nodes[1]
			result2, err := evaluator.EvaluateWithContext(node2.Condition, fact, ctx)
			if err != nil {
				t.Fatalf("Step 2 evaluation failed: %v", err)
			}
			ctx.SetIntermediateResult(node2.ResultName, result2)
			t.Logf("  Step 2: %s = %.0f", node2.ResultName, result2)

			if result2 != tc.expected {
				t.Errorf("Step 2: Expected %.0f, got %.0f", tc.expected, result2)
			}

			// Step 3: temp_2 > 100 → temp_3 (comparison)
			node3 := alphaChain.Nodes[2]
			result3, err := evaluator.EvaluateWithContext(node3.Condition, fact, ctx)
			if err != nil {
				t.Fatalf("Step 3 evaluation failed: %v", err)
			}
			boolResult := result3.(bool)
			t.Logf("  Step 3: %s = %v", node3.ResultName, boolResult)

			if boolResult != tc.shouldPass {
				t.Errorf("Expected %v, got %v (%.0f > 100)", tc.shouldPass, boolResult, tc.expected)
			}
		})
	}
}

// TestArithmeticDecomposition_ActivateWithContext tests ActivateWithContext on decomposed chain
func TestArithmeticDecomposition_ActivateWithContext(t *testing.T) {
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

	// Simple expression: c.qte * 23 > 100
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "*",
			"left":     map[string]interface{}{"type": "fieldAccess", "field": "qte"},
			"right":    map[string]interface{}{"type": "number", "value": 23},
		},
		"right": map[string]interface{}{"type": "number", "value": 100},
	}

	// Decompose
	decomposer := NewArithmeticExpressionDecomposer()
	decomposedSteps, _ := decomposer.DecomposeToDecomposedConditions(condition)

	// Build chain
	chainBuilder := NewAlphaChainBuilder(network, storage)
	alphaChain, _ := chainBuilder.BuildDecomposedChain(decomposedSteps, "c", typeNode, "test_rule")

	// Test with ActivateWithContext
	fact := &Fact{
		ID:   "test",
		Type: "Order",
		Fields: map[string]interface{}{
			"qte": 5, // 5 * 23 = 115 > 100 → true
		},
	}

	ctx := NewEvaluationContext(fact)

	// Activate the first node with context
	err := alphaChain.Nodes[0].ActivateWithContext(fact, ctx)
	if err != nil {
		t.Fatalf("ActivateWithContext failed: %v", err)
	}

	// Verify intermediate results were stored
	for i, node := range alphaChain.Nodes {
		if !ctx.HasIntermediateResult(node.ResultName) {
			t.Errorf("Node %d: intermediate result %s not found in context", i, node.ResultName)
		} else {
			value, _ := ctx.GetIntermediateResult(node.ResultName)
			t.Logf("Node %d: %s = %v", i, node.ResultName, value)
		}
	}

	// Verify final result
	finalResult, exists := ctx.GetIntermediateResult(alphaChain.Nodes[len(alphaChain.Nodes)-1].ResultName)
	if !exists {
		t.Fatal("Final result not found in context")
	}

	if boolResult, ok := finalResult.(bool); !ok || !boolResult {
		t.Errorf("Expected final result to be true, got %v", finalResult)
	}

	t.Log("✅ ActivateWithContext successfully propagated intermediate results")
}

// TestArithmeticDecomposition_TypeNodeActivation tests TypeNode creating context for decomposed chains
func TestArithmeticDecomposition_TypeNodeActivation(t *testing.T) {
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

	// Expression: c.qte * 23 > 100
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "*",
			"left":     map[string]interface{}{"type": "fieldAccess", "field": "qte"},
			"right":    map[string]interface{}{"type": "number", "value": 23},
		},
		"right": map[string]interface{}{"type": "number", "value": 100},
	}

	// Build decomposed chain
	decomposer := NewArithmeticExpressionDecomposer()
	decomposedSteps, _ := decomposer.DecomposeToDecomposedConditions(condition)
	chainBuilder := NewAlphaChainBuilder(network, storage)
	alphaChain, _ := chainBuilder.BuildDecomposedChain(decomposedSteps, "c", typeNode, "test_rule")

	t.Logf("Built chain with %d nodes", len(alphaChain.Nodes))

	// Submit fact through TypeNode (should trigger context creation)
	fact := &Fact{
		ID:   "test",
		Type: "Order",
		Fields: map[string]interface{}{
			"qte": 10, // 10 * 23 = 230 > 100 → true
		},
	}

	err := typeNode.ActivateRight(fact)
	if err != nil {
		t.Fatalf("TypeNode.ActivateRight failed: %v", err)
	}

	// Check that facts propagated through the chain
	for i, node := range alphaChain.Nodes {
		factCount := len(node.Memory.Facts)
		t.Logf("Node %d (%s): %d fact(s) in memory", i, node.ResultName, factCount)
		if factCount != 1 {
			t.Errorf("Node %d: Expected 1 fact in memory, got %d", i, factCount)
		}
	}

	t.Log("✅ TypeNode successfully activated decomposed chain with context")
}

// TestArithmeticDecomposition_WithJoin tests decomposition in a 2-variable join rule
func TestArithmeticDecomposition_WithJoin(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Define types
	produitType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Produit",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "prix", Type: "number"},
		},
	}

	commandeType := TypeDefinition{
		Type: "typeDefinition",
		Name: "Commande",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "produit_id", Type: "string"},
			{Name: "qte", Type: "number"},
		},
	}

	// Create TypeNodes
	produitTypeNode := NewTypeNode("Produit", produitType, storage)
	commandeTypeNode := NewTypeNode("Commande", commandeType, storage)
	network.TypeNodes["Produit"] = produitTypeNode
	network.TypeNodes["Commande"] = commandeTypeNode

	// Create terminal node (dummy)
	terminalNode := &TerminalNode{
		BaseNode: BaseNode{
			ID:       "test_terminal",
			Type:     "terminal",
			Memory:   &WorkingMemory{NodeID: "test_terminal", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children: make([]Node, 0),
			Storage:  storage,
		},
		Action: &Action{
			Type: "action",
			Job: &JobCall{
				Name: "test_action",
				Args: []interface{}{},
			},
		},
	}

	// Rule: {p: Produit, c: Commande} / c.produit_id == p.id AND (c.qte * 23) > 100
	// Alpha condition: (c.qte * 23) > 100
	// Beta condition: c.produit_id == p.id

	alphaCondition := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "*",
			"left":     map[string]interface{}{"type": "fieldAccess", "field": "qte"},
			"right":    map[string]interface{}{"type": "number", "value": 23},
		},
		"right": map[string]interface{}{"type": "number", "value": 100},
	}

	betaCondition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
		"left":     map[string]interface{}{"type": "fieldAccess", "object": "c", "field": "produit_id"},
		"right":    map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "id"},
	}

	// Decompose alpha condition
	decomposer := NewArithmeticExpressionDecomposer()
	decomposedSteps, err := decomposer.DecomposeToDecomposedConditions(alphaCondition)
	if err != nil {
		t.Fatalf("Decomposition failed: %v", err)
	}

	t.Logf("Decomposed into %d steps", len(decomposedSteps))

	// Build decomposed chain for variable c
	chainBuilder := NewAlphaChainBuilder(network, storage)
	alphaChain, err := chainBuilder.BuildDecomposedChain(
		decomposedSteps,
		"c",
		commandeTypeNode,
		"test_rule",
	)
	if err != nil {
		t.Fatalf("BuildDecomposedChain failed: %v", err)
	}

	t.Logf("Chain built with %d nodes", len(alphaChain.Nodes))

	// Create passthrough for p (left side - no alpha filter)
	passthroughP := NewAlphaNode(
		"passthrough_p",
		map[string]interface{}{"type": "passthrough", "side": "left"},
		"p",
		storage,
	)
	produitTypeNode.AddChild(passthroughP)

	// Create passthrough for c (right side - after alpha chain)
	passthroughC := NewAlphaNode(
		"passthrough_c",
		map[string]interface{}{"type": "passthrough", "side": "right"},
		"c",
		storage,
	)
	alphaChain.FinalNode.AddChild(passthroughC)

	// Create JoinNode
	varTypes := map[string]string{"p": "Produit", "c": "Commande"}
	joinNode := NewJoinNode(
		"test_join",
		betaCondition,
		[]string{"p"},
		[]string{"c"},
		varTypes,
		storage,
	)

	passthroughP.AddChild(joinNode)
	passthroughC.AddChild(joinNode)
	joinNode.AddChild(terminalNode)

	network.BetaNodes["test_join"] = joinNode

	t.Log("Network structure built")

	// Submit facts IN ORDER: Produit first, then Commande
	produit := &Fact{
		ID:   "PROD001",
		Type: "Produit",
		Fields: map[string]interface{}{
			"id":   "PROD001",
			"prix": 100,
		},
	}

	commande := &Fact{
		ID:   "CMD001",
		Type: "Commande",
		Fields: map[string]interface{}{
			"produit_id": "PROD001",
			"qte":        5, // 5 * 23 = 115 > 100 → true
		},
	}

	t.Log("Submitting Produit fact...")
	t.Logf("ProduitTypeNode has %d children", len(produitTypeNode.GetChildren()))
	t.Logf("PassthroughP has %d children", len(passthroughP.GetChildren()))
	err = produitTypeNode.ActivateRight(produit)
	if err != nil {
		t.Fatalf("Error submitting Produit: %v", err)
	}
	t.Logf("After Produit submission - JoinNode LeftMemory: %d tokens", len(joinNode.LeftMemory.Tokens))

	t.Log("Submitting Commande fact...")
	t.Logf("CommandeTypeNode has %d children", len(commandeTypeNode.GetChildren()))
	t.Logf("AlphaChain.FinalNode has %d children", len(alphaChain.FinalNode.GetChildren()))
	t.Logf("PassthroughC has %d children", len(passthroughC.GetChildren()))
	err = commandeTypeNode.ActivateRight(commande)
	if err != nil {
		t.Fatalf("Error submitting Commande: %v", err)
	}
	t.Logf("After Commande submission - JoinNode RightMemory: %d tokens", len(joinNode.RightMemory.Tokens))

	// Check terminal tokens
	tokenCount := len(terminalNode.Memory.Tokens)
	t.Logf("Terminal node has %d token(s)", tokenCount)

	if tokenCount != 1 {
		t.Errorf("Expected 1 token in terminal, got %d", tokenCount)
		t.Logf("JoinNode memory: %d left tokens, %d right tokens",
			len(joinNode.LeftMemory.Tokens), len(joinNode.RightMemory.Tokens))
		t.Logf("PassthroughP memory: %d facts", len(passthroughP.Memory.Facts))
		t.Logf("PassthroughC memory: %d facts", len(passthroughC.Memory.Facts))

		// Debug: check token bindings
		for tokenID, token := range joinNode.LeftMemory.Tokens {
			t.Logf("Left token %s: bindings = %v, facts = %d", tokenID, getBindingKeys(token.Bindings), len(token.Facts))
			for varName, fact := range token.Bindings {
				t.Logf("  Binding %s -> %s (type: %s, fields: %v)", varName, fact.ID, fact.Type, fact.Fields)
			}
		}
		for tokenID, token := range joinNode.RightMemory.Tokens {
			t.Logf("Right token %s: bindings = %v, facts = %d", tokenID, getBindingKeys(token.Bindings), len(token.Facts))
			for varName, fact := range token.Bindings {
				t.Logf("  Binding %s -> %s (type: %s, fields: %v)", varName, fact.ID, fact.Type, fact.Fields)
			}
		}

		// Check beta condition
		t.Logf("Beta condition: %v", betaCondition)
	} else {
		t.Log("✅ Join successful! Token created and propagated to terminal")
	}
}

func getBindingKeys(bindings map[string]*Fact) []string {
	keys := make([]string, 0, len(bindings))
	for k := range bindings {
		keys = append(keys, k)
	}
	return keys
}
