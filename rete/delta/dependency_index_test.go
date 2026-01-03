// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
	"testing"
	"time"
)

func TestNewDependencyIndex(t *testing.T) {
	idx := NewDependencyIndex()

	if idx == nil {
		t.Fatal("NewDependencyIndex returned nil")
	}

	if idx.alphaIndex == nil || idx.betaIndex == nil || idx.terminalIndex == nil {
		t.Error("Indexes not initialized")
	}

	if idx.nodeReferences == nil {
		t.Error("nodeReferences not initialized")
	}

	if time.Since(idx.builtAt) > time.Second {
		t.Error("builtAt timestamp incorrect")
	}
}

func TestDependencyIndex_AddAlphaNode(t *testing.T) {
	idx := NewDependencyIndex()

	idx.AddAlphaNode("alpha1", "Product", []string{"price", "status"})

	nodes := idx.GetAffectedNodes("Product", "price")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(nodes))
	}

	if nodes[0].NodeID != "alpha1" {
		t.Errorf("Expected alpha1, got %s", nodes[0].NodeID)
	}

	if nodes[0].NodeType != "alpha" {
		t.Errorf("Expected alpha type, got %s", nodes[0].NodeType)
	}

	nodes = idx.GetAffectedNodes("Product", "status")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node for status, got %d", len(nodes))
	}
}

func TestDependencyIndex_AddBetaNode(t *testing.T) {
	idx := NewDependencyIndex()

	idx.AddBetaNode("beta1", "Order", []string{"customer_id"})

	nodes := idx.GetAffectedNodes("Order", "customer_id")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(nodes))
	}

	if nodes[0].NodeType != "beta" {
		t.Errorf("Expected beta type, got %s", nodes[0].NodeType)
	}
}

func TestDependencyIndex_AddTerminalNode(t *testing.T) {
	idx := NewDependencyIndex()

	idx.AddTerminalNode("term1", "Alert", []string{"severity", "message"})

	nodes := idx.GetAffectedNodes("Alert", "severity")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(nodes))
	}

	if nodes[0].NodeType != "terminal" {
		t.Errorf("Expected terminal type, got %s", nodes[0].NodeType)
	}
}

func TestDependencyIndex_MultipleNodesPerField(t *testing.T) {
	idx := NewDependencyIndex()

	idx.AddAlphaNode("alpha1", "Product", []string{"price"})
	idx.AddAlphaNode("alpha2", "Product", []string{"price", "category"})
	idx.AddBetaNode("beta1", "Product", []string{"price"})

	nodes := idx.GetAffectedNodes("Product", "price")

	if len(nodes) != 3 {
		t.Fatalf("Expected 3 nodes for price, got %d", len(nodes))
	}

	nodeIDs := make(map[string]bool)
	for _, node := range nodes {
		nodeIDs[node.NodeID] = true
	}

	if !nodeIDs["alpha1"] || !nodeIDs["alpha2"] || !nodeIDs["beta1"] {
		t.Errorf("Missing expected nodes: %v", nodeIDs)
	}
}

func TestDependencyIndex_NoDuplicateIndexing(t *testing.T) {
	idx := NewDependencyIndex()

	idx.AddAlphaNode("alpha1", "Product", []string{"price"})
	idx.AddAlphaNode("alpha1", "Product", []string{"price"})

	nodes := idx.GetAffectedNodes("Product", "price")

	if len(nodes) != 1 {
		t.Errorf("Expected 1 node (no duplicate), got %d", len(nodes))
	}
}

func TestDependencyIndex_GetAffectedNodesForDelta(t *testing.T) {
	idx := NewDependencyIndex()

	idx.AddAlphaNode("alpha1", "Product", []string{"price"})
	idx.AddAlphaNode("alpha2", "Product", []string{"status"})
	idx.AddBetaNode("beta1", "Product", []string{"price", "category"})
	idx.AddTerminalNode("term1", "Product", []string{"status"})

	delta := NewFactDelta("Product~123", "Product")
	delta.AddFieldChange("price", 100.0, 150.0)
	delta.AddFieldChange("status", "active", "inactive")

	nodes := idx.GetAffectedNodesForDelta(delta)

	if len(nodes) != 4 {
		t.Fatalf("Expected 4 affected nodes, got %d", len(nodes))
	}

	nodeIDs := make(map[string]bool)
	for _, node := range nodes {
		nodeIDs[node.NodeID] = true
	}

	expected := []string{"alpha1", "alpha2", "beta1", "term1"}
	for _, expectedID := range expected {
		if !nodeIDs[expectedID] {
			t.Errorf("Missing expected node: %s", expectedID)
		}
	}
}

func TestDependencyIndex_DifferentFactTypes(t *testing.T) {
	idx := NewDependencyIndex()

	idx.AddAlphaNode("alpha_product", "Product", []string{"id"})
	idx.AddAlphaNode("alpha_order", "Order", []string{"id"})

	nodes := idx.GetAffectedNodes("Product", "id")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node for Product.id, got %d", len(nodes))
	}
	if nodes[0].NodeID != "alpha_product" {
		t.Errorf("Expected alpha_product, got %s", nodes[0].NodeID)
	}

	nodes = idx.GetAffectedNodes("Order", "id")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node for Order.id, got %d", len(nodes))
	}
	if nodes[0].NodeID != "alpha_order" {
		t.Errorf("Expected alpha_order, got %s", nodes[0].NodeID)
	}
}

func TestDependencyIndex_Clear(t *testing.T) {
	idx := NewDependencyIndex()

	idx.AddAlphaNode("alpha1", "Product", []string{"price"})
	idx.AddBetaNode("beta1", "Order", []string{"total"})

	if len(idx.GetAffectedNodes("Product", "price")) == 0 {
		t.Fatal("Nodes should exist before Clear")
	}

	idx.Clear()

	if len(idx.GetAffectedNodes("Product", "price")) != 0 {
		t.Error("Index should be empty after Clear")
	}

	stats := idx.GetStats()
	if stats.NodeCount != 0 || stats.FieldCount != 0 {
		t.Errorf("Stats should be zero after Clear: %+v", stats)
	}
}

func TestDependencyIndex_GetStats(t *testing.T) {
	idx := NewDependencyIndex()

	idx.AddAlphaNode("alpha1", "Product", []string{"price", "status"})
	idx.AddAlphaNode("alpha2", "Product", []string{"category"})
	idx.AddBetaNode("beta1", "Order", []string{"customer_id"})
	idx.AddTerminalNode("term1", "Alert", []string{"severity"})

	stats := idx.GetStats()

	if stats.NodeCount != 4 {
		t.Errorf("Expected 4 nodes, got %d", stats.NodeCount)
	}

	if stats.AlphaNodeCount != 2 {
		t.Errorf("Expected 2 alpha nodes, got %d", stats.AlphaNodeCount)
	}

	if stats.BetaNodeCount != 1 {
		t.Errorf("Expected 1 beta node, got %d", stats.BetaNodeCount)
	}

	if stats.TerminalCount != 1 {
		t.Errorf("Expected 1 terminal node, got %d", stats.TerminalCount)
	}

	if len(stats.FactTypes) != 3 {
		t.Errorf("Expected 3 fact types, got %d", len(stats.FactTypes))
	}

	if stats.MemoryEstimate <= 0 {
		t.Error("MemoryEstimate should be positive")
	}
}

func TestDependencyIndex_String(t *testing.T) {
	idx := NewDependencyIndex()
	idx.AddAlphaNode("alpha1", "Product", []string{"price"})

	str := idx.String()
	if str == "" {
		t.Error("String() should not be empty")
	}

	if len(str) < 20 {
		t.Errorf("String() too short: %s", str)
	}
}

func TestNodeReference_String(t *testing.T) {
	ref := NodeReference{
		NodeID:   "alpha1",
		NodeType: "alpha",
		FactType: "Product",
		Fields:   []string{"price"},
	}

	str := ref.String()
	expectedMatch := "alpha[alpha1](Product)"
	if str != expectedMatch {
		t.Errorf("String() = %s, want %s", str, expectedMatch)
	}
}

func TestDependencyIndex_ConcurrentAccess(t *testing.T) {
	idx := NewDependencyIndex()

	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			idx.AddAlphaNode("alpha1", "Product", []string{"price"})
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			_ = idx.GetAffectedNodes("Product", "price")
		}
		done <- true
	}()

	<-done
	<-done

	nodes := idx.GetAffectedNodes("Product", "price")
	if len(nodes) != 1 {
		t.Errorf("Expected 1 node after concurrent access, got %d", len(nodes))
	}
}

func BenchmarkDependencyIndex_AddAlphaNode(b *testing.B) {
	idx := NewDependencyIndex()
	fields := []string{"field1", "field2", "field3"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nodeID := "alpha" + string(rune(i))
		idx.AddAlphaNode(nodeID, "TestType", fields)
	}
}

func BenchmarkDependencyIndex_GetAffectedNodes(b *testing.B) {
	idx := NewDependencyIndex()

	for i := 0; i < 100; i++ {
		nodeID := "alpha" + string(rune(i))
		idx.AddAlphaNode(nodeID, "Product", []string{"price"})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = idx.GetAffectedNodes("Product", "price")
	}
}

func BenchmarkDependencyIndex_GetAffectedNodesForDelta(b *testing.B) {
	idx := NewDependencyIndex()

	for i := 0; i < 50; i++ {
		idx.AddAlphaNode("alpha"+string(rune(i)), "Product", []string{"price", "status"})
	}

	delta := NewFactDelta("Product~123", "Product")
	delta.AddFieldChange("price", 100, 150)
	delta.AddFieldChange("status", "active", "inactive")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = idx.GetAffectedNodesForDelta(delta)
	}
}

// Test lecture/écriture concurrente intensive
func TestDependencyIndex_ConcurrentReadWrite(t *testing.T) {
	idx := NewDependencyIndex()

	// Setup initial
	for i := 0; i < 10; i++ {
		nodeID := fmt.Sprintf("node%d", i)
		idx.AddAlphaNode(nodeID, "Product", []string{"field1", "field2"})
	}

	done := make(chan bool, 20)

	// 10 writers
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < TestConcurrentOps; j++ {
				nodeID := fmt.Sprintf("concurrent_node_%d_%d", id, j)
				idx.AddAlphaNode(nodeID, "Product", []string{"price", "status"})
			}
			done <- true
		}(i)
	}

	// 10 readers
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < TestConcurrentOps; j++ {
				_ = idx.GetAffectedNodes("Product", "price")
			}
			done <- true
		}()
	}

	// Attendre fin
	for i := 0; i < 20; i++ {
		<-done
	}

	// Vérifier cohérence finale
	stats := idx.GetStats()
	if stats.NodeCount == 0 {
		t.Error("❌ Expected nodes in index after concurrent operations")
	}
	t.Logf("✅ Concurrent ops completed. Nodes: %d", stats.NodeCount)
}

// Test clear pendant lecture
func TestDependencyIndex_ClearDuringRead(t *testing.T) {
	idx := NewDependencyIndex()

	for i := 0; i < TestConcurrentOps; i++ {
		idx.AddAlphaNode(fmt.Sprintf("node%d", i), "Product", []string{"price"})
	}

	done := make(chan bool, 2)

	// Reader continu
	go func() {
		for i := 0; i < TestStressRuns; i++ {
			_ = idx.GetAffectedNodes("Product", "price")
		}
		done <- true
	}()

	// Clear au milieu
	go func() {
		time.Sleep(5 * time.Millisecond)
		idx.Clear()
		done <- true
	}()

	<-done
	<-done

	// Ne devrait pas paniquer
	stats := idx.GetStats()
	if stats.NodeCount != 0 {
		t.Errorf("❌ Expected 0 nodes after clear, got %d", stats.NodeCount)
	}
	t.Log("✅ Clear during read completed without panic")
}

// Test ajout concurrent de même nœud
func TestDependencyIndex_ConcurrentSameNode(t *testing.T) {
	idx := NewDependencyIndex()

	done := make(chan bool, 10)
	const nodeID = "shared_node"

	// Plusieurs goroutines ajoutent le même nœud
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				idx.AddAlphaNode(nodeID, "Product", []string{"price"})
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// Le nœud ne devrait être présent qu'une fois
	nodes := idx.GetAffectedNodes("Product", "price")
	count := 0
	for _, node := range nodes {
		if node.NodeID == nodeID {
			count++
		}
	}

	if count != 1 {
		t.Errorf("❌ Expected node to appear once, found %d times", count)
	}
}

// Test GetAffectedNodesForDelta concurrent
func TestDependencyIndex_ConcurrentDeltaQuery(t *testing.T) {
	idx := NewDependencyIndex()

	// Setup
	for i := 0; i < 50; i++ {
		idx.AddAlphaNode(fmt.Sprintf("alpha_%d", i), "Product", []string{"price", "status"})
	}

	delta := NewFactDelta("Product~123", "Product")
	delta.AddFieldChange("price", 100, 150)
	delta.AddFieldChange("status", "active", "inactive")

	done := make(chan bool, 10)

	// Lecture concurrente
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < TestConcurrentOps; j++ {
				nodes := idx.GetAffectedNodesForDelta(delta)
				if len(nodes) == 0 {
					t.Error("❌ Expected nodes for delta")
				}
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	t.Log("✅ Concurrent delta queries completed")
}

// Test avec grande quantité de nœuds
func TestDependencyIndex_LargeScale(t *testing.T) {
	idx := NewDependencyIndex()

	const nodeCount = 10000
	t.Logf("📊 Adding %d nodes...", nodeCount)

	start := time.Now()
	for i := 0; i < nodeCount; i++ {
		nodeID := fmt.Sprintf("node_%d", i)
		factType := fmt.Sprintf("Type_%d", i%10)
		fields := []string{"field1", "field2", "field3"}
		idx.AddAlphaNode(nodeID, factType, fields)
	}
	duration := time.Since(start)

	stats := idx.GetStats()
	if stats.NodeCount != nodeCount {
		t.Errorf("❌ Expected %d nodes, got %d", nodeCount, stats.NodeCount)
	}

	t.Logf("✅ Added %d nodes in %v", nodeCount, duration)
	t.Logf("📊 Memory estimate: %d bytes (%.2f MB)",
		stats.MemoryEstimate, float64(stats.MemoryEstimate)/(1024*1024))

	// Test query performance
	start = time.Now()
	nodes := idx.GetAffectedNodes("Type_0", "field1")
	queryDuration := time.Since(start)

	t.Logf("✅ Query returned %d nodes in %v", len(nodes), queryDuration)
}

// Test stats accuracy
func TestDependencyIndex_StatsAccuracy(t *testing.T) {
	idx := NewDependencyIndex()

	// Ajouter différents types de nœuds
	idx.AddAlphaNode("alpha1", "Product", []string{"price"})
	idx.AddAlphaNode("alpha2", "Product", []string{"status"})
	idx.AddBetaNode("beta1", "Order", []string{"product_id"})
	idx.AddTerminalNode("term1", "Product", []string{"price", "status"})

	stats := idx.GetStats()

	// Test stats accuracy
	if stats.TerminalCount != 1 {
		t.Errorf("❌ Expected 1 terminal node, got %d", stats.TerminalCount)
	}
	if stats.NodeCount != 4 {
		t.Errorf("❌ Expected 4 total nodes, got %d", stats.NodeCount)
	}

	// Vérifier fact types
	if len(stats.FactTypes) != 2 {
		t.Errorf("❌ Expected 2 fact types, got %d", len(stats.FactTypes))
	}

	// Vérifier estimation mémoire
	if stats.MemoryEstimate == 0 {
		t.Error("❌ Expected non-zero memory estimate")
	}

	t.Logf("✅ Stats: %d nodes (%d alpha, %d beta, %d terminal)",
		stats.NodeCount, stats.AlphaNodeCount, stats.BetaNodeCount, stats.TerminalCount)
	t.Logf("📊 Memory estimate: %d bytes", stats.MemoryEstimate)
}

// Test multiples clear
func TestDependencyIndex_MultipleClear(t *testing.T) {
	idx := NewDependencyIndex()

	for iteration := 0; iteration < 5; iteration++ {
		// Ajouter nœuds
		for i := 0; i < 20; i++ {
			idx.AddAlphaNode(fmt.Sprintf("node_%d", i), "Product", []string{"price"})
		}

		stats := idx.GetStats()
		if stats.NodeCount != 20 {
			t.Errorf("❌ Iteration %d: expected 20 nodes, got %d", iteration, stats.NodeCount)
		}

		// Clear
		idx.Clear()

		stats = idx.GetStats()
		if stats.NodeCount != 0 {
			t.Errorf("❌ Iteration %d: expected 0 nodes after clear, got %d", iteration, stats.NodeCount)
		}
	}

	t.Log("✅ Multiple clear operations successful")
}
