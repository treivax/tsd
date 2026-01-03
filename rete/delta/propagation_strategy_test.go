// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
)

func TestSequentialStrategy_GetName(t *testing.T) {
	strategy := &SequentialStrategy{}
	if strategy.GetName() != "Sequential" {
		t.Errorf("Expected 'Sequential', got '%s'", strategy.GetName())
	}
}

func TestSequentialStrategy_ShouldPropagate(t *testing.T) {
	strategy := &SequentialStrategy{}
	delta := NewFactDelta("Test~1", "Test")

	t.Run("with nodes", func(t *testing.T) {
		nodes := []NodeReference{{NodeID: "n1", NodeType: "alpha"}}
		if !strategy.ShouldPropagate(delta, nodes) {
			t.Error("Expected true when nodes present")
		}
	})

	t.Run("without nodes", func(t *testing.T) {
		nodes := []NodeReference{}
		if strategy.ShouldPropagate(delta, nodes) {
			t.Error("Expected false when no nodes")
		}
	})
}

func TestSequentialStrategy_GetPropagationOrder(t *testing.T) {
	strategy := &SequentialStrategy{}

	nodes := []NodeReference{
		{NodeID: "term1", NodeType: "terminal", FactType: "Test"},
		{NodeID: "alpha1", NodeType: "alpha", FactType: "Test"},
		{NodeID: "beta1", NodeType: "beta", FactType: "Test"},
		{NodeID: "alpha2", NodeType: "alpha", FactType: "Test"},
		{NodeID: "term2", NodeType: "terminal", FactType: "Test"},
	}

	ordered := strategy.GetPropagationOrder(nodes)

	if len(ordered) != len(nodes) {
		t.Fatalf("Expected %d nodes, got %d", len(nodes), len(ordered))
	}

	alphaCount := 0
	for i, node := range ordered {
		if node.NodeType == "alpha" {
			alphaCount++
		} else if node.NodeType == "beta" {
			if alphaCount != 2 {
				t.Errorf("Beta at position %d before all alphas", i)
			}
		}
	}

	if alphaCount != 2 {
		t.Errorf("Expected 2 alpha nodes, got %d", alphaCount)
	}
}

func TestTopologicalStrategy_GetName(t *testing.T) {
	strategy := NewTopologicalStrategy()
	if strategy.GetName() != "Topological" {
		t.Errorf("Expected 'Topological', got '%s'", strategy.GetName())
	}
}

func TestTopologicalStrategy_SetNodeDepth(t *testing.T) {
	strategy := NewTopologicalStrategy()

	strategy.SetNodeDepth("node1", 1)
	strategy.SetNodeDepth("node2", 2)

	if strategy.getDepth("node1") != 1 {
		t.Errorf("Expected depth 1, got %d", strategy.getDepth("node1"))
	}

	if strategy.getDepth("node2") != 2 {
		t.Errorf("Expected depth 2, got %d", strategy.getDepth("node2"))
	}

	if strategy.getDepth("unknown") != 0 {
		t.Errorf("Expected depth 0 for unknown node, got %d", strategy.getDepth("unknown"))
	}
}

func TestTopologicalStrategy_GetPropagationOrder(t *testing.T) {
	strategy := NewTopologicalStrategy()

	strategy.SetNodeDepth("node1", 3)
	strategy.SetNodeDepth("node2", 1)
	strategy.SetNodeDepth("node3", 2)

	nodes := []NodeReference{
		{NodeID: "node1", NodeType: "alpha"},
		{NodeID: "node2", NodeType: "alpha"},
		{NodeID: "node3", NodeType: "beta"},
	}

	ordered := strategy.GetPropagationOrder(nodes)

	if ordered[0].NodeID != "node2" {
		t.Errorf("Expected node2 first, got %s", ordered[0].NodeID)
	}
	if ordered[1].NodeID != "node3" {
		t.Errorf("Expected node3 second, got %s", ordered[1].NodeID)
	}
	if ordered[2].NodeID != "node1" {
		t.Errorf("Expected node1 third, got %s", ordered[2].NodeID)
	}
}

func TestTopologicalStrategy_GetPropagationOrder_NoDepths(t *testing.T) {
	strategy := NewTopologicalStrategy()

	nodes := []NodeReference{
		{NodeID: "term1", NodeType: "terminal"},
		{NodeID: "alpha1", NodeType: "alpha"},
	}

	ordered := strategy.GetPropagationOrder(nodes)

	if len(ordered) != 2 {
		t.Fatalf("Expected 2 nodes, got %d", len(ordered))
	}

	if ordered[0].NodeType != "alpha" {
		t.Error("Expected alpha first in fallback mode")
	}
}

func TestOptimizedStrategy_GetName(t *testing.T) {
	strategy := &OptimizedStrategy{}
	if strategy.GetName() != "Optimized" {
		t.Errorf("Expected 'Optimized', got '%s'", strategy.GetName())
	}
}

func TestOptimizedStrategy_ShouldPropagate(t *testing.T) {
	strategy := &OptimizedStrategy{}

	t.Run("empty delta", func(t *testing.T) {
		delta := NewFactDelta("Test~1", "Test")
		nodes := []NodeReference{{NodeID: "n1"}}

		if strategy.ShouldPropagate(delta, nodes) {
			t.Error("Expected false for empty delta")
		}
	})

	t.Run("no nodes", func(t *testing.T) {
		delta := NewFactDelta("Test~1", "Test")
		delta.AddFieldChange("field", "old", "new")
		nodes := []NodeReference{}

		if strategy.ShouldPropagate(delta, nodes) {
			t.Error("Expected false when no nodes")
		}
	})

	t.Run("valid propagation", func(t *testing.T) {
		delta := NewFactDelta("Test~1", "Test")
		delta.AddFieldChange("field", "old", "new")
		nodes := []NodeReference{{NodeID: "n1"}}

		if !strategy.ShouldPropagate(delta, nodes) {
			t.Error("Expected true for valid propagation")
		}
	})
}

func TestOptimizedStrategy_GetPropagationOrder(t *testing.T) {
	strategy := &OptimizedStrategy{}

	nodes := []NodeReference{
		{NodeID: "term1", NodeType: "terminal", FactType: "Product"},
		{NodeID: "alpha1", NodeType: "alpha", FactType: "Product"},
		{NodeID: "beta1", NodeType: "beta", FactType: "Order"},
		{NodeID: "alpha2", NodeType: "alpha", FactType: "Order"},
		{NodeID: "term2", NodeType: "terminal", FactType: "Product"},
	}

	ordered := strategy.GetPropagationOrder(nodes)

	if len(ordered) != len(nodes) {
		t.Fatalf("Expected %d nodes, got %d", len(nodes), len(ordered))
	}

	seenBeta := false
	seenTerminal := false

	for _, node := range ordered {
		if node.NodeType == "beta" {
			seenBeta = true
		}
		if node.NodeType == "terminal" {
			seenTerminal = true
		}

		if node.NodeType == "alpha" && (seenBeta || seenTerminal) {
			t.Error("Alpha node found after beta or terminal")
		}
	}
}

func TestOptimizedStrategy_GetPropagationOrder_Empty(t *testing.T) {
	strategy := &OptimizedStrategy{}

	ordered := strategy.GetPropagationOrder([]NodeReference{})

	if len(ordered) != 0 {
		t.Errorf("Expected empty result, got %d nodes", len(ordered))
	}
}
