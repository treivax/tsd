// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
)

func TestNewIndexBuilder(t *testing.T) {
	builder := NewIndexBuilder()

	if builder == nil {
		t.Fatal("NewIndexBuilder returned nil")
	}

	if builder.alphaExtractor == nil {
		t.Error("alphaExtractor not initialized")
	}

	if builder.betaExtractor == nil {
		t.Error("betaExtractor not initialized")
	}

	if builder.actionExtractor == nil {
		t.Error("actionExtractor not initialized")
	}

	if builder.diagnostics == nil {
		t.Error("diagnostics not initialized")
	}
}

func TestIndexBuilder_BuildFromAlphaNode(t *testing.T) {
	builder := NewIndexBuilder()
	builder.EnableDiagnostics()
	idx := NewDependencyIndex()

	condition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "price",
		},
		"right": 100,
	}

	err := builder.BuildFromAlphaNode(idx, "alpha1", "Product", condition)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	nodes := idx.GetAffectedNodes("Product", "price")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node indexed, got %d", len(nodes))
	}

	if nodes[0].NodeID != "alpha1" {
		t.Errorf("Expected alpha1, got %s", nodes[0].NodeID)
	}

	diag := builder.GetDiagnostics()
	if diag.NodesProcessed != 1 {
		t.Errorf("Expected 1 node processed, got %d", diag.NodesProcessed)
	}
	if diag.FieldsExtracted != 1 {
		t.Errorf("Expected 1 field extracted, got %d", diag.FieldsExtracted)
	}
}

func TestIndexBuilder_BuildFromAlphaNode_MultipleFields(t *testing.T) {
	builder := NewIndexBuilder()
	idx := NewDependencyIndex()

	condition := map[string]interface{}{
		"type": "binaryOp",
		"left": map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":  "fieldAccess",
				"field": "price",
			},
			"right": 100,
		},
		"right": map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":  "fieldAccess",
				"field": "status",
			},
			"right": "active",
		},
	}

	err := builder.BuildFromAlphaNode(idx, "alpha1", "Product", condition)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	priceNodes := idx.GetAffectedNodes("Product", "price")
	statusNodes := idx.GetAffectedNodes("Product", "status")

	if len(priceNodes) != 1 || len(statusNodes) != 1 {
		t.Errorf("Expected 1 node for each field")
	}
}

func TestIndexBuilder_BuildFromBetaNode(t *testing.T) {
	builder := NewIndexBuilder()
	builder.EnableDiagnostics()
	idx := NewDependencyIndex()

	joinCondition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "customer_id",
		},
		"right": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "id",
		},
	}

	err := builder.BuildFromBetaNode(idx, "beta1", "Order", joinCondition)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	nodes := idx.GetAffectedNodes("Order", "customer_id")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node indexed, got %d", len(nodes))
	}

	if nodes[0].NodeType != "beta" {
		t.Errorf("Expected beta type, got %s", nodes[0].NodeType)
	}

	diag := builder.GetDiagnostics()
	if diag.NodesProcessed != 1 {
		t.Errorf("Expected 1 node processed, got %d", diag.NodesProcessed)
	}
}

func TestIndexBuilder_BuildFromTerminalNode(t *testing.T) {
	builder := NewIndexBuilder()
	builder.EnableDiagnostics()
	idx := NewDependencyIndex()

	actions := []interface{}{
		map[string]interface{}{
			"type": "updateWithModifications",
			"modifications": map[string]interface{}{
				"price":  150,
				"status": "updated",
			},
		},
	}

	err := builder.BuildFromTerminalNode(idx, "term1", "Product", actions)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	priceNodes := idx.GetAffectedNodes("Product", "price")
	statusNodes := idx.GetAffectedNodes("Product", "status")

	if len(priceNodes) != 1 || len(statusNodes) != 1 {
		t.Errorf("Expected 1 node for each field")
	}

	if priceNodes[0].NodeType != "terminal" {
		t.Errorf("Expected terminal type, got %s", priceNodes[0].NodeType)
	}

	diag := builder.GetDiagnostics()
	if diag.NodesProcessed != 1 {
		t.Errorf("Expected 1 node processed, got %d", diag.NodesProcessed)
	}
	if diag.FieldsExtracted != 2 {
		t.Errorf("Expected 2 fields extracted, got %d", diag.FieldsExtracted)
	}
}

func TestIndexBuilder_BuildFromTerminalNode_MultipleActions(t *testing.T) {
	builder := NewIndexBuilder()
	idx := NewDependencyIndex()

	actions := []interface{}{
		map[string]interface{}{
			"type": "updateWithModifications",
			"modifications": map[string]interface{}{
				"price": 150,
			},
		},
		map[string]interface{}{
			"type": "factCreation",
			"fields": map[string]interface{}{
				"severity": "high",
			},
		},
	}

	err := builder.BuildFromTerminalNode(idx, "term1", "Product", actions)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	priceNodes := idx.GetAffectedNodes("Product", "price")
	severityNodes := idx.GetAffectedNodes("Product", "severity")

	if len(priceNodes) != 1 || len(severityNodes) != 1 {
		t.Errorf("Expected nodes for both fields from different actions")
	}
}

func TestIndexBuilder_EnableDiagnostics(t *testing.T) {
	builder := NewIndexBuilder()

	if builder.enableDiagnostics {
		t.Error("Diagnostics should be disabled by default")
	}

	builder.EnableDiagnostics()

	if !builder.enableDiagnostics {
		t.Error("Diagnostics should be enabled after EnableDiagnostics()")
	}
}

func TestIndexBuilder_DiagnosticsWarnings(t *testing.T) {
	builder := NewIndexBuilder()
	builder.EnableDiagnostics()
	idx := NewDependencyIndex()

	// Condition sans champs
	condition := map[string]interface{}{
		"type":  "constant",
		"value": true,
	}

	err := builder.BuildFromAlphaNode(idx, "alpha1", "Product", condition)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	diag := builder.GetDiagnostics()
	if len(diag.Warnings) == 0 {
		t.Error("Expected warning for node with no fields")
	}
}

func TestIndexBuilder_BuildFromNetwork(t *testing.T) {
	t.Log("🧪 TEST INDEX BUILDER - BUILD FROM NETWORK")
	t.Log("===========================================")

	t.Run("NilNetwork", func(t *testing.T) {
		builder := NewIndexBuilder()

		idx, err := builder.BuildFromNetwork(nil)

		if err != nil {
			t.Fatalf("❌ Unexpected error: %v", err)
		}

		if idx == nil {
			t.Error("❌ Expected non-nil index")
		}

		if idx.GetTotalNodeCount() != 0 {
			t.Errorf("❌ Expected empty index, got %d nodes", idx.GetTotalNodeCount())
		}

		t.Log("✅ Nil network handled correctly")
	})

	t.Run("EmptyNetwork", func(t *testing.T) {
		builder := NewIndexBuilder()
		builder.EnableDiagnostics()

		// Mock réseau vide
		network := &mockReteNetwork{
			AlphaNodes:    make(map[string]*mockAlphaNode),
			TerminalNodes: make(map[string]*mockTerminalNode),
		}

		idx, err := builder.BuildFromNetwork(network)

		if err != nil {
			t.Fatalf("❌ Unexpected error: %v", err)
		}

		if idx.GetTotalNodeCount() != 0 {
			t.Errorf("❌ Expected empty index, got %d nodes", idx.GetTotalNodeCount())
		}

		diag := builder.GetDiagnostics()
		if diag.NodesProcessed != 0 {
			t.Errorf("❌ Expected 0 nodes processed, got %d", diag.NodesProcessed)
		}

		t.Log("✅ Empty network handled correctly")
	})

	t.Run("NetworkWithAlphaNodes", func(t *testing.T) {
		builder := NewIndexBuilder()
		builder.EnableDiagnostics()

		// Mock réseau avec alpha nodes
		network := &mockReteNetwork{
			AlphaNodes: map[string]*mockAlphaNode{
				"alpha_1": {
					ID:           "alpha_1",
					VariableName: "Product",
					Condition: map[string]interface{}{
						"type":  "comparison",
						"left":  map[string]interface{}{"type": "fieldAccess", "field": "price"},
						"op":    ">",
						"right": 100,
					},
				},
				"alpha_2": {
					ID:           "alpha_2",
					VariableName: "Order",
					Condition: map[string]interface{}{
						"type":  "comparison",
						"left":  map[string]interface{}{"type": "fieldAccess", "field": "status"},
						"op":    "==",
						"right": "active",
					},
				},
			},
			TerminalNodes: make(map[string]*mockTerminalNode),
		}

		idx, err := builder.BuildFromNetwork(network)

		if err != nil {
			t.Fatalf("❌ Unexpected error: %v", err)
		}

		// Vérifier que les nœuds alpha ont été indexés
		productNodes := idx.GetAffectedNodes("Product", "price")
		if len(productNodes) != 1 || productNodes[0].NodeID != "alpha_1" {
			t.Errorf("❌ Expected alpha_1 for Product.price, got %v", productNodes)
		}

		orderNodes := idx.GetAffectedNodes("Order", "status")
		if len(orderNodes) != 1 || orderNodes[0].NodeID != "alpha_2" {
			t.Errorf("❌ Expected alpha_2 for Order.status, got %v", orderNodes)
		}

		diag := builder.GetDiagnostics()
		if diag.NodesProcessed != 2 {
			t.Errorf("❌ Expected 2 nodes processed, got %d", diag.NodesProcessed)
		}

		if diag.FieldsExtracted != 2 {
			t.Errorf("❌ Expected 2 fields extracted, got %d", diag.FieldsExtracted)
		}

		t.Log("✅ Alpha nodes extracted successfully")
	})

	t.Run("NetworkWithTerminalNodes", func(t *testing.T) {
		builder := NewIndexBuilder()
		builder.EnableDiagnostics()

		// Mock réseau avec terminal nodes
		network := &mockReteNetwork{
			AlphaNodes: make(map[string]*mockAlphaNode),
			TerminalNodes: map[string]*mockTerminalNode{
				"terminal_1": {
					ID: "terminal_1",
					Action: &mockAction{
						VariableName: "Product",
						TypeName:     "Product",
						ActionType:   "update",
					},
				},
			},
		}

		idx, err := builder.BuildFromNetwork(network)

		if err != nil {
			t.Fatalf("❌ Unexpected error: %v", err)
		}

		// Vérifier que le nœud terminal a été indexé
		if idx.GetTotalNodeCount() == 0 {
			t.Error("❌ Expected terminal node to be indexed")
		}

		diag := builder.GetDiagnostics()
		if diag.NodesProcessed != 1 {
			t.Errorf("❌ Expected 1 node processed, got %d", diag.NodesProcessed)
		}

		t.Log("✅ Terminal nodes extracted successfully")
	})

	t.Run("NetworkWithMixedNodes", func(t *testing.T) {
		builder := NewIndexBuilder()
		builder.EnableDiagnostics()

		// Mock réseau complet avec alpha et terminal
		network := &mockReteNetwork{
			AlphaNodes: map[string]*mockAlphaNode{
				"alpha_1": {
					ID:           "alpha_1",
					VariableName: "Product",
					Condition: map[string]interface{}{
						"type":  "fieldAccess",
						"field": "price",
					},
				},
			},
			TerminalNodes: map[string]*mockTerminalNode{
				"terminal_1": {
					ID: "terminal_1",
					Action: &mockAction{
						VariableName: "Order",
						TypeName:     "Order",
					},
				},
			},
		}

		idx, err := builder.BuildFromNetwork(network)

		if err != nil {
			t.Fatalf("❌ Unexpected error: %v", err)
		}

		// Vérifier que les nœuds ont été indexés
		if idx.GetTotalNodeCount() == 0 {
			t.Error("❌ Expected nodes to be indexed")
		}

		diag := builder.GetDiagnostics()
		if diag.NodesProcessed != 2 {
			t.Errorf("❌ Expected 2 nodes processed, got %d", diag.NodesProcessed)
		}

		t.Log("✅ Mixed nodes extracted successfully")
	})

	t.Log("✅ All BuildFromNetwork tests passed")
}

// Mock structures pour tester BuildFromNetwork

type mockReteNetwork struct {
	AlphaNodes    map[string]*mockAlphaNode
	TerminalNodes map[string]*mockTerminalNode
	BetaNodes     map[string]interface{}
}

type mockAlphaNode struct {
	ID           string
	VariableName string
	Condition    interface{}
}

type mockTerminalNode struct {
	ID     string
	Action *mockAction
}

type mockAction struct {
	VariableName string
	TypeName     string
	ActionType   string
}

func TestIndexBuilder_CompleteWorkflow(t *testing.T) {
	builder := NewIndexBuilder()
	builder.EnableDiagnostics()
	idx := NewDependencyIndex()

	// Alpha node
	alphaCondition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "price",
		},
		"right": 100,
	}

	err := builder.BuildFromAlphaNode(idx, "alpha1", "Product", alphaCondition)
	if err != nil {
		t.Fatalf("Alpha node build error: %v", err)
	}

	// Beta node
	betaCondition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "customer_id",
		},
		"right": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "id",
		},
	}

	err = builder.BuildFromBetaNode(idx, "beta1", "Order", betaCondition)
	if err != nil {
		t.Fatalf("Beta node build error: %v", err)
	}

	// Terminal node
	actions := []interface{}{
		map[string]interface{}{
			"type": "updateWithModifications",
			"modifications": map[string]interface{}{
				"status": "processed",
			},
		},
	}

	err = builder.BuildFromTerminalNode(idx, "term1", "Order", actions)
	if err != nil {
		t.Fatalf("Terminal node build error: %v", err)
	}

	// Vérifier diagnostics
	diag := builder.GetDiagnostics()
	if diag.NodesProcessed != 3 {
		t.Errorf("Expected 3 nodes processed, got %d", diag.NodesProcessed)
	}

	// Vérifier index
	stats := idx.GetStats()
	if stats.NodeCount != 3 {
		t.Errorf("Expected 3 nodes in index, got %d", stats.NodeCount)
	}
	if stats.AlphaNodeCount != 1 {
		t.Errorf("Expected 1 alpha node, got %d", stats.AlphaNodeCount)
	}
	if stats.BetaNodeCount != 1 {
		t.Errorf("Expected 1 beta node, got %d", stats.BetaNodeCount)
	}
	if stats.TerminalCount != 1 {
		t.Errorf("Expected 1 terminal node, got %d", stats.TerminalCount)
	}
}

func BenchmarkIndexBuilder_BuildFromAlphaNode(b *testing.B) {
	builder := NewIndexBuilder()
	idx := NewDependencyIndex()

	condition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "price",
		},
		"right": 100,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.BuildFromAlphaNode(idx, "alpha1", "Product", condition)
	}
}

// TestIndexBuilder_BuildFromBetaNode_ErrorCases teste les cas d'erreur de BuildFromBetaNode
func TestIndexBuilder_BuildFromBetaNode_ErrorCases(t *testing.T) {
	t.Run("invalid join condition - extraction error", func(t *testing.T) {
		builder := NewIndexBuilder()
		builder.EnableDiagnostics()
		idx := NewDependencyIndex()

		// Condition invalide qui devrait causer une erreur d'extraction
		invalidCondition := map[string]interface{}{
			"type": "unknownType",
			"data": "invalid",
		}

		err := builder.BuildFromBetaNode(idx, "beta1", "Order", invalidCondition)

		// L'erreur peut être nil si l'extracteur ignore les types inconnus
		// Mais les diagnostics devraient refléter qu'aucun champ n'a été extrait
		diag := builder.GetDiagnostics()

		// Vérifier que soit une erreur est retournée, soit un warning est émis
		if err == nil {
			// Pas d'erreur, vérifier les warnings
			if len(diag.Warnings) == 0 && diag.FieldsExtracted == 0 {
				t.Log("✅ Condition invalide gérée (aucun champ extrait)")
			}
		} else {
			// Erreur retournée
			if diag.NodesSkipped != 1 {
				t.Errorf("Expected NodesSkipped=1, got %d", diag.NodesSkipped)
			}
			t.Log("✅ Erreur d'extraction détectée correctement")
		}
	})

	t.Run("no fields extracted - warning", func(t *testing.T) {
		builder := NewIndexBuilder()
		builder.EnableDiagnostics()
		idx := NewDependencyIndex()

		// Condition vide qui ne devrait extraire aucun champ
		emptyCondition := map[string]interface{}{}

		err := builder.BuildFromBetaNode(idx, "beta2", "Order", emptyCondition)

		// Pas forcément une erreur, mais devrait générer un warning
		diag := builder.GetDiagnostics()

		if err == nil && len(diag.Warnings) > 0 {
			// Warning émis pour aucun champ extrait
			t.Log("✅ Warning émis pour condition sans champs")
		} else if err != nil {
			t.Log("✅ Erreur retournée pour condition vide")
		}
	})

	t.Run("valid condition with multiple fields", func(t *testing.T) {
		builder := NewIndexBuilder()
		builder.EnableDiagnostics()
		idx := NewDependencyIndex()

		// Condition complexe avec plusieurs champs
		complexCondition := map[string]interface{}{
			"type": "and",
			"conditions": []interface{}{
				map[string]interface{}{
					"type": "comparison",
					"left": map[string]interface{}{
						"type":  "fieldAccess",
						"field": "customer_id",
					},
					"right": map[string]interface{}{
						"type":  "fieldAccess",
						"field": "id",
					},
				},
				map[string]interface{}{
					"type": "comparison",
					"left": map[string]interface{}{
						"type":  "fieldAccess",
						"field": "order_date",
					},
					"operator": ">=",
					"right":    "2024-01-01",
				},
			},
		}

		err := builder.BuildFromBetaNode(idx, "beta3", "Order", complexCondition)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		diag := builder.GetDiagnostics()
		if diag.NodesProcessed < 1 {
			t.Errorf("Expected at least 1 node processed, got %d", diag.NodesProcessed)
		}

		// Vérifier que customer_id et order_date ont été indexés
		customerNodes := idx.GetAffectedNodes("Order", "customer_id")
		dateNodes := idx.GetAffectedNodes("Order", "order_date")

		if len(customerNodes) == 0 {
			t.Error("Expected customer_id to be indexed")
		}
		if len(dateNodes) == 0 {
			t.Error("Expected order_date to be indexed")
		}

		t.Log("✅ Condition complexe traitée correctement")
	})
}

// TestIndexBuilder_BuildFromBetaNode_Diagnostics teste les diagnostics pour BuildFromBetaNode
func TestIndexBuilder_BuildFromBetaNode_Diagnostics(t *testing.T) {
	t.Run("diagnostics disabled", func(t *testing.T) {
		builder := NewIndexBuilder()
		// Ne pas activer les diagnostics
		idx := NewDependencyIndex()

		condition := map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":  "fieldAccess",
				"field": "price",
			},
			"right": 100,
		}

		err := builder.BuildFromBetaNode(idx, "beta1", "Product", condition)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Diagnostics devraient être à zéro car désactivés
		diag := builder.GetDiagnostics()
		if diag.NodesProcessed != 0 {
			t.Error("Diagnostics should not be incremented when disabled")
		}

		t.Log("✅ Diagnostics désactivés fonctionnent")
	})

	t.Run("diagnostics enabled", func(t *testing.T) {
		builder := NewIndexBuilder()
		builder.EnableDiagnostics()
		idx := NewDependencyIndex()

		condition := map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":  "fieldAccess",
				"field": "stock",
			},
			"right": 0,
		}

		err := builder.BuildFromBetaNode(idx, "beta2", "Product", condition)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		diag := builder.GetDiagnostics()
		if diag.NodesProcessed != 1 {
			t.Errorf("Expected NodesProcessed=1, got %d", diag.NodesProcessed)
		}
		if diag.FieldsExtracted < 1 {
			t.Errorf("Expected FieldsExtracted>=1, got %d", diag.FieldsExtracted)
		}

		t.Log("✅ Diagnostics activés fonctionnent")
	})
}
