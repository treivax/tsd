// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// TestE2EBindingsDebug is a focused test to debug the 3-variable join issue
// This test reproduces the scenario from beta_join_complex.tsd rule r2
func TestE2EBindingsDebug(t *testing.T) {
	// Enable debug logging
	logger := GetDebugLogger()
	logger.Enable()
	defer logger.Disable()

	logger.Log("========== E2E BINDINGS DEBUG TEST START ==========")

	// Create network
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Define types
	userType := TypeDefinition{
		Name: "User",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
			{Name: "age", Type: "number"},
			{Name: "city", Type: "string"},
			{Name: "status", Type: "string"},
		},
	}

	orderType := TypeDefinition{
		Name: "Order",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "user_id", Type: "string"},
			{Name: "product_id", Type: "string"},
			{Name: "amount", Type: "number"},
			{Name: "date", Type: "string"},
		},
	}

	productType := TypeDefinition{
		Name: "Product",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
			{Name: "category", Type: "string"},
			{Name: "price", Type: "number"},
		},
	}

	// Add type nodes
	userTypeNode := NewTypeNode("User", userType, storage)
	orderTypeNode := NewTypeNode("Order", orderType, storage)
	productTypeNode := NewTypeNode("Product", productType, storage)

	network.TypeNodes["User"] = userTypeNode
	network.TypeNodes["Order"] = orderTypeNode
	network.TypeNodes["Product"] = productTypeNode

	network.RootNode.AddChild(userTypeNode)
	network.RootNode.AddChild(orderTypeNode)
	network.RootNode.AddChild(productTypeNode)

	// Create passthrough alphas (simulating what the builder does)
	network.PassthroughRegistry = make(map[string]*AlphaNode)

	// Passthrough for User (variable u) - feeds LEFT side of JoinNode1
	passthroughUser := NewAlphaNode("passthrough_User_u", map[string]interface{}{
		"type":        "passthrough",
		"source_type": "User",
		"variable":    "u",
		"side":        "left",
	}, "u", storage)
	network.PassthroughRegistry["User_u"] = passthroughUser
	userTypeNode.AddChild(passthroughUser)

	// Passthrough for Order (variable o) - feeds RIGHT side of JoinNode1
	passthroughOrder := NewAlphaNode("passthrough_Order_o", map[string]interface{}{
		"type":        "passthrough",
		"source_type": "Order",
		"variable":    "o",
		"side":        "right",
	}, "o", storage)
	network.PassthroughRegistry["Order_o"] = passthroughOrder
	orderTypeNode.AddChild(passthroughOrder)

	// Passthrough for Product (variable p) - feeds RIGHT side of JoinNode2
	passthroughProduct := NewAlphaNode("passthrough_Product_p", map[string]interface{}{
		"type":        "passthrough",
		"source_type": "Product",
		"variable":    "p",
		"side":        "right",
	}, "p", storage)
	network.PassthroughRegistry["Product_p"] = passthroughProduct
	productTypeNode.AddChild(passthroughProduct)

	// Create JoinNode1: User ⋈ Order
	// Condition: u.status == "vip" AND o.user_id == u.id
	joinNode1 := NewJoinNode(
		"join_r2_step1",
		map[string]interface{}{
			"type": "logicalExpr",
			"left": map[string]interface{}{
				"type": "comparison",
				"left": map[string]interface{}{
					"type":     "fieldAccess",
					"variable": "u",
					"field":    "status",
				},
				"operator": "==",
				"right": map[string]interface{}{
					"type":  "literal",
					"value": "vip",
				},
			},
			"operator": "AND",
			"operations": []map[string]interface{}{
				{
					"operator": "AND",
					"right": map[string]interface{}{
						"type": "comparison",
						"left": map[string]interface{}{
							"type":     "fieldAccess",
							"variable": "o",
							"field":    "user_id",
						},
						"operator": "==",
						"right": map[string]interface{}{
							"type":     "fieldAccess",
							"variable": "u",
							"field":    "id",
						},
					},
				},
			},
		},
		[]string{"u"}, // left variables
		[]string{"o"}, // right variables
		map[string]string{ // variable types
			"u": "User",
			"o": "Order",
		},
		storage,
	)

	// Create JoinNode2: (User ⋈ Order) ⋈ Product
	// Condition: p.id == o.product_id AND p.category == "luxury"
	joinNode2 := NewJoinNode(
		"join_r2_step2",
		map[string]interface{}{
			"type": "logicalExpr",
			"left": map[string]interface{}{
				"type": "comparison",
				"left": map[string]interface{}{
					"type":     "fieldAccess",
					"variable": "p",
					"field":    "id",
				},
				"operator": "==",
				"right": map[string]interface{}{
					"type":     "fieldAccess",
					"variable": "o",
					"field":    "product_id",
				},
			},
			"operator": "AND",
			"operations": []map[string]interface{}{
				{
					"operator": "AND",
					"right": map[string]interface{}{
						"type": "comparison",
						"left": map[string]interface{}{
							"type":     "fieldAccess",
							"variable": "p",
							"field":    "category",
						},
						"operator": "==",
						"right": map[string]interface{}{
							"type":  "literal",
							"value": "luxury",
						},
					},
				},
			},
		},
		[]string{"u", "o"}, // left variables (result from joinNode1)
		[]string{"p"},      // right variables
		map[string]string{ // variable types
			"u": "User",
			"o": "Order",
			"p": "Product",
		},
		storage,
	)

	// Connect the network
	// User passthrough (side=left) -> JoinNode1 left input
	passthroughUser.AddChild(joinNode1)

	// Order passthrough (side=right) -> JoinNode1 right input
	passthroughOrder.AddChild(joinNode1)

	// JoinNode1 result -> JoinNode2 left input (via ActivateLeft from PropagateToChildren)
	joinNode1.AddChild(joinNode2)

	// Product passthrough (side=right) -> JoinNode2 right input
	passthroughProduct.AddChild(joinNode2)

	// Create terminal node
	terminalNode := NewTerminalNode("terminal_r2", &Action{
		Type: "action",
		Job: &JobCall{
			Type: "jobCall",
			Name: "vip_luxury_purchase",
			Args: []interface{}{
				map[string]interface{}{"type": "fieldAccess", "variable": "u", "field": "id"},
				map[string]interface{}{"type": "fieldAccess", "variable": "p", "field": "name"},
			},
		},
	}, storage)
	joinNode2.AddChild(terminalNode)

	logger.LogNetworkStructure(network)

	logger.Log("========== SUBMITTING FACTS ==========")

	// Submit facts in order
	// 1. User (VIP status)
	userFact := &Fact{
		ID:   "USER001",
		Type: "User",
		Fields: map[string]interface{}{
			"id":     "USER001",
			"name":   "Alice",
			"age":    30.0,
			"city":   "Paris",
			"status": "vip",
		},
	}
	logger.Log(">> Submitting User fact: %s", userFact.ID)
	if err := network.SubmitFact(userFact); err != nil {
		t.Fatalf("Failed to submit user fact: %v", err)
	}

	// 2. Order
	orderFact := &Fact{
		ID:   "ORD004",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":         "ORD004",
			"user_id":    "USER001",
			"product_id": "PROD003",
			"amount":     1.0,
			"date":       "2025-01-18",
		},
	}
	logger.Log(">> Submitting Order fact: %s", orderFact.ID)
	if err := network.SubmitFact(orderFact); err != nil {
		t.Fatalf("Failed to submit order fact: %v", err)
	}

	// 3. Product (luxury category)
	productFact := &Fact{
		ID:   "PROD003",
		Type: "Product",
		Fields: map[string]interface{}{
			"id":       "PROD003",
			"name":     "Luxury Watch",
			"category": "luxury",
			"price":    2500.0,
		},
	}
	logger.Log(">> Submitting Product fact: %s", productFact.ID)
	if err := network.SubmitFact(productFact); err != nil {
		t.Fatalf("Failed to submit product fact: %v", err)
	}

	logger.Log("========== CHECKING RESULTS ==========")

	// Check JoinNode1 result memory
	joinNode1ResultTokens := joinNode1.ResultMemory.GetTokens()
	logger.Log("JoinNode1 ResultMemory: %d tokens", len(joinNode1ResultTokens))
	for i, token := range joinNode1ResultTokens {
		logger.Log("  Token[%d]: vars=%v", i, token.GetVariables())
		logger.LogBindings("  ", token.Bindings)
	}

	// Check JoinNode2 result memory
	joinNode2ResultTokens := joinNode2.ResultMemory.GetTokens()
	logger.Log("JoinNode2 ResultMemory: %d tokens", len(joinNode2ResultTokens))
	for i, token := range joinNode2ResultTokens {
		logger.Log("  Token[%d]: vars=%v", i, token.GetVariables())
		logger.LogBindings("  ", token.Bindings)
	}

	// Check terminal node - TerminalNodes execute immediately and don't store tokens
	// So we check the execution count instead
	executionCount := terminalNode.GetExecutionCount()
	logger.Log("TerminalNode Execution Count: %d", executionCount)

	logger.Log("========== E2E BINDINGS DEBUG TEST END ==========")

	// Final assertion - check execution count instead of token count
	if executionCount != 1 {
		t.Errorf("Expected 1 terminal execution, got %d", executionCount)
	}
}
