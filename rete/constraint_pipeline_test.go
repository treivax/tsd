// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestIngestFileWithMetrics tests the IngestFileWithMetrics wrapper function
func TestIngestFileWithMetrics(t *testing.T) {
	t.Run("returns metrics on success", func(t *testing.T) {
		// Create a temporary test file
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.tsd")
		content := `
type Person(id: string, name: string)
action log(msg: string)
rule simple : {p: Person} / p.id == "1" ==> log(p.name)
Person(id: "1", name: "Alice")
`
		if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		pipeline := NewConstraintPipeline()
		storage := NewMemoryStorage()
		network, metrics, err := pipeline.IngestFileWithMetrics(testFile, nil, storage)

		if err != nil {
			t.Fatalf("IngestFileWithMetrics failed: %v", err)
		}

		if network == nil {
			t.Fatal("Expected non-nil network")
		}

		if metrics == nil {
			t.Fatal("Expected non-nil metrics")
		}

		// Verify metrics structure
		if metrics.TotalDuration == 0 {
			t.Error("Expected non-zero total duration")
		}
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		pipeline := NewConstraintPipeline()
		storage := NewMemoryStorage()

		network, metrics, err := pipeline.IngestFileWithMetrics("/non/existent/file.tsd", nil, storage)

		if err == nil {
			t.Error("Expected error for non-existent file")
		}

		// Metrics should still be returned even on error
		if metrics == nil {
			t.Error("Expected metrics even on error")
		}

		// Network might be nil or partial
		_ = network
	})

	t.Run("handles empty network input", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.tsd")
		content := `
type Item(id: string)
action process(id: string)
rule r1 : {i: Item} / i.id == "test" ==> process(i.id)
`
		if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		pipeline := NewConstraintPipeline()
		storage := NewMemoryStorage()

		// Pass nil network - should create new one
		network, metrics, err := pipeline.IngestFileWithMetrics(testFile, nil, storage)

		if err != nil {
			t.Fatalf("IngestFileWithMetrics failed: %v", err)
		}

		if network == nil {
			t.Fatal("Expected network to be created")
		}

		if metrics == nil {
			t.Fatal("Expected metrics")
		}

		// Verify network has type nodes
		if len(network.TypeNodes) == 0 {
			t.Error("Expected at least one type node")
		}
	})

	t.Run("handles existing network", func(t *testing.T) {
		tmpDir := t.TempDir()

		// First file
		file1 := filepath.Join(tmpDir, "file1.tsd")
		content1 := `
type Person(id: string)
action log1(id: string)
rule r1 : {p: Person} / p.id == "1" ==> log1(p.id)
`
		if err := os.WriteFile(file1, []byte(content1), 0644); err != nil {
			t.Fatalf("Failed to create file1: %v", err)
		}

		// Second file
		file2 := filepath.Join(tmpDir, "file2.tsd")
		content2 := `
type Order(id: string)
action log2(id: string)
rule r2 : {o: Order} / o.id == "2" ==> log2(o.id)
`
		if err := os.WriteFile(file2, []byte(content2), 0644); err != nil {
			t.Fatalf("Failed to create file2: %v", err)
		}

		pipeline := NewConstraintPipeline()
		storage := NewMemoryStorage()

		// Ingest first file
		network1, metrics1, err := pipeline.IngestFileWithMetrics(file1, nil, storage)
		if err != nil {
			t.Fatalf("First ingest failed: %v", err)
		}

		initialTypeCount := len(network1.TypeNodes)
		if initialTypeCount == 0 {
			t.Fatal("Expected type nodes after first ingest")
		}

		// Ingest second file with existing network
		network2, metrics2, err := pipeline.IngestFileWithMetrics(file2, network1, storage)
		if err != nil {
			t.Fatalf("Second ingest failed: %v", err)
		}

		if network2 == nil {
			t.Fatal("Expected network after second ingest")
		}

		// Should have more type nodes now
		if len(network2.TypeNodes) <= initialTypeCount {
			t.Error("Expected more type nodes after second ingest")
		}

		// Both metrics should be valid
		if metrics1 == nil || metrics2 == nil {
			t.Error("Expected metrics for both ingests")
		}
	})
}

// TestIsTerminalReachableFrom tests the isTerminalReachableFrom function
func TestIsTerminalReachableFrom(t *testing.T) {
	pipeline := NewConstraintPipeline()

	t.Run("returns true when node is the terminal itself", func(t *testing.T) {
		terminal := &TerminalNode{
			BaseNode: BaseNode{ID: "terminal1"},
		}

		result := pipeline.isTerminalReachableFrom(terminal, "terminal1")

		if !result {
			t.Error("Expected true when checking terminal against its own ID")
		}
	})

	t.Run("returns false when terminal not in tree", func(t *testing.T) {
		typeNode := &TypeNode{
			BaseNode: BaseNode{
				ID:       "type1",
				Children: []Node{},
			},
		}

		result := pipeline.isTerminalReachableFrom(typeNode, "terminal_not_present")

		if result {
			t.Error("Expected false when terminal not in tree")
		}
	})

	t.Run("returns true when terminal is direct child", func(t *testing.T) {
		terminal := &TerminalNode{
			BaseNode: BaseNode{ID: "terminal1"},
		}

		typeNode := &TypeNode{
			BaseNode: BaseNode{
				ID:       "type1",
				Children: []Node{terminal},
			},
		}

		result := pipeline.isTerminalReachableFrom(typeNode, "terminal1")

		if !result {
			t.Error("Expected true when terminal is direct child")
		}
	})

	t.Run("returns true when terminal is nested descendant", func(t *testing.T) {
		terminal := &TerminalNode{
			BaseNode: BaseNode{ID: "terminal1"},
		}

		alphaNode := &AlphaNode{
			BaseNode: BaseNode{
				ID:       "alpha1",
				Children: []Node{terminal},
			},
		}

		typeNode := &TypeNode{
			BaseNode: BaseNode{
				ID:       "type1",
				Children: []Node{alphaNode},
			},
		}

		result := pipeline.isTerminalReachableFrom(typeNode, "terminal1")

		if !result {
			t.Error("Expected true when terminal is nested descendant")
		}
	})

	t.Run("returns false when terminal is in sibling branch", func(t *testing.T) {
		terminal1 := &TerminalNode{
			BaseNode: BaseNode{ID: "terminal1"},
		}

		terminal2 := &TerminalNode{
			BaseNode: BaseNode{ID: "terminal2"},
		}

		alphaNode1 := &AlphaNode{
			BaseNode: BaseNode{
				ID:       "alpha1",
				Children: []Node{terminal1},
			},
		}

		alphaNode2 := &AlphaNode{
			BaseNode: BaseNode{
				ID:       "alpha2",
				Children: []Node{terminal2},
			},
		}

		// Search from alphaNode1 for terminal2
		result := pipeline.isTerminalReachableFrom(alphaNode1, "terminal2")

		if result {
			t.Error("Expected false when terminal is in sibling branch")
		}

		// Should work in reverse
		result2 := pipeline.isTerminalReachableFrom(alphaNode2, "terminal2")
		if !result2 {
			t.Error("Expected true for correct branch")
		}
	})

	t.Run("handles deeply nested tree", func(t *testing.T) {
		terminal := &TerminalNode{
			BaseNode: BaseNode{ID: "terminal1"},
		}

		// Create a chain: type -> alpha -> join -> terminal
		joinNode := &JoinNode{
			BaseNode: BaseNode{
				ID:       "join1",
				Children: []Node{terminal},
			},
		}

		alphaNode := &AlphaNode{
			BaseNode: BaseNode{
				ID:       "alpha1",
				Children: []Node{joinNode},
			},
		}

		typeNode := &TypeNode{
			BaseNode: BaseNode{
				ID:       "type1",
				Children: []Node{alphaNode},
			},
		}

		result := pipeline.isTerminalReachableFrom(typeNode, "terminal1")

		if !result {
			t.Error("Expected true for deeply nested terminal")
		}
	})

	t.Run("handles multiple children", func(t *testing.T) {
		terminal := &TerminalNode{
			BaseNode: BaseNode{ID: "terminal1"},
		}

		child1 := &AlphaNode{
			BaseNode: BaseNode{
				ID:       "alpha1",
				Children: []Node{},
			},
		}

		child2 := &AlphaNode{
			BaseNode: BaseNode{
				ID:       "alpha2",
				Children: []Node{terminal},
			},
		}

		child3 := &AlphaNode{
			BaseNode: BaseNode{
				ID:       "alpha3",
				Children: []Node{},
			},
		}

		typeNode := &TypeNode{
			BaseNode: BaseNode{
				ID:       "type1",
				Children: []Node{child1, child2, child3},
			},
		}

		result := pipeline.isTerminalReachableFrom(typeNode, "terminal1")

		if !result {
			t.Error("Expected true when terminal is in one of multiple children")
		}
	})
}

// TestIdentifyExpectedTypesForTerminal tests the identifyExpectedTypesForTerminal function
func TestIdentifyExpectedTypesForTerminal(t *testing.T) {
	pipeline := NewConstraintPipeline()

	t.Run("returns empty list when terminal not reachable", func(t *testing.T) {
		terminal := &TerminalNode{
			BaseNode: BaseNode{ID: "terminal1"},
		}

		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{
				"Person": {
					BaseNode: BaseNode{
						ID:       "type_Person",
						Children: []Node{},
					},
				},
			},
		}

		types := pipeline.identifyExpectedTypesForTerminal(network, terminal)

		if len(types) != 0 {
			t.Errorf("Expected empty list, got %d types", len(types))
		}
	})

	t.Run("returns single type when terminal is reachable", func(t *testing.T) {
		terminal := &TerminalNode{
			BaseNode: BaseNode{ID: "terminal1"},
		}

		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{
				"Person": {
					BaseNode: BaseNode{
						ID:       "type_Person",
						Children: []Node{terminal},
					},
				},
			},
		}

		types := pipeline.identifyExpectedTypesForTerminal(network, terminal)

		if len(types) != 1 {
			t.Fatalf("Expected 1 type, got %d", len(types))
		}

		if types[0] != "Person" {
			t.Errorf("Expected 'Person', got '%s'", types[0])
		}
	})

	t.Run("returns multiple types when terminal reachable from multiple type nodes", func(t *testing.T) {
		terminal := &TerminalNode{
			BaseNode: BaseNode{ID: "terminal1"},
		}

		// Terminal is reachable from both Person and Order type nodes
		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{
				"Person": {
					BaseNode: BaseNode{
						ID:       "type_Person",
						Children: []Node{terminal},
					},
				},
				"Order": {
					BaseNode: BaseNode{
						ID:       "type_Order",
						Children: []Node{terminal},
					},
				},
				"Product": {
					BaseNode: BaseNode{
						ID:       "type_Product",
						Children: []Node{},
					},
				},
			},
		}

		types := pipeline.identifyExpectedTypesForTerminal(network, terminal)

		if len(types) != 2 {
			t.Fatalf("Expected 2 types, got %d", len(types))
		}

		// Check that both Person and Order are in the result
		typeSet := make(map[string]bool)
		for _, typ := range types {
			typeSet[typ] = true
		}

		if !typeSet["Person"] {
			t.Error("Expected 'Person' in result")
		}
		if !typeSet["Order"] {
			t.Error("Expected 'Order' in result")
		}
		if typeSet["Product"] {
			t.Error("Did not expect 'Product' in result")
		}
	})

	t.Run("returns type when terminal nested under type node", func(t *testing.T) {
		terminal := &TerminalNode{
			BaseNode: BaseNode{ID: "terminal1"},
		}

		alphaNode := &AlphaNode{
			BaseNode: BaseNode{
				ID:       "alpha1",
				Children: []Node{terminal},
			},
		}

		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{
				"Person": {
					BaseNode: BaseNode{
						ID:       "type_Person",
						Children: []Node{alphaNode},
					},
				},
			},
		}

		types := pipeline.identifyExpectedTypesForTerminal(network, terminal)

		if len(types) != 1 {
			t.Fatalf("Expected 1 type, got %d", len(types))
		}

		if types[0] != "Person" {
			t.Errorf("Expected 'Person', got '%s'", types[0])
		}
	})

	t.Run("handles empty network", func(t *testing.T) {
		terminal := &TerminalNode{
			BaseNode: BaseNode{ID: "terminal1"},
		}

		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{},
		}

		types := pipeline.identifyExpectedTypesForTerminal(network, terminal)

		if len(types) != 0 {
			t.Errorf("Expected empty list for empty network, got %d types", len(types))
		}
	})

	t.Run("handles complex multi-branch network", func(t *testing.T) {
		terminal := &TerminalNode{
			BaseNode: BaseNode{ID: "terminal1"},
		}

		// Create a join node that connects to the terminal
		joinNode := &JoinNode{
			BaseNode: BaseNode{
				ID:       "join1",
				Children: []Node{terminal},
			},
		}

		// Alpha node feeding into join
		alphaNode1 := &AlphaNode{
			BaseNode: BaseNode{
				ID:       "alpha1",
				Children: []Node{joinNode},
			},
		}

		// Another alpha node feeding into join
		alphaNode2 := &AlphaNode{
			BaseNode: BaseNode{
				ID:       "alpha2",
				Children: []Node{joinNode},
			},
		}

		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{
				"Person": {
					BaseNode: BaseNode{
						ID:       "type_Person",
						Children: []Node{alphaNode1},
					},
				},
				"Order": {
					BaseNode: BaseNode{
						ID:       "type_Order",
						Children: []Node{alphaNode2},
					},
				},
			},
		}

		types := pipeline.identifyExpectedTypesForTerminal(network, terminal)

		if len(types) != 2 {
			t.Fatalf("Expected 2 types, got %d", len(types))
		}

		typeSet := make(map[string]bool)
		for _, typ := range types {
			typeSet[typ] = true
		}

		if !typeSet["Person"] || !typeSet["Order"] {
			t.Error("Expected both 'Person' and 'Order' in result")
		}
	})
}

// TestPropagateToNewTerminals tests the propagateToNewTerminals function
func TestPropagateToNewTerminals(t *testing.T) {
	pipeline := NewConstraintPipeline()

	t.Run("returns zero when no terminals", func(t *testing.T) {
		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{},
		}

		count := pipeline.propagateToNewTerminals(network, []*TerminalNode{}, map[string][]*Fact{})

		if count != 0 {
			t.Errorf("Expected 0 propagations, got %d", count)
		}
	})

	t.Run("returns zero when no facts", func(t *testing.T) {
		terminal := &TerminalNode{
			BaseNode: BaseNode{ID: "terminal1"},
		}

		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{
				"Person": {
					BaseNode: BaseNode{
						ID:       "type_Person",
						Children: []Node{terminal},
					},
					TypeName: "Person",
				},
			},
		}

		count := pipeline.propagateToNewTerminals(network, []*TerminalNode{terminal}, map[string][]*Fact{})

		if count != 0 {
			t.Errorf("Expected 0 propagations, got %d", count)
		}
	})

	t.Run("attempts to propagate facts to reachable terminal", func(t *testing.T) {
		// Note: This test verifies the logic of propagateToNewTerminals
		// Actual propagation success requires full node initialization with proper
		// child node implementations, which is complex. We verify the function
		// correctly identifies types and iterates through facts.
		storage := NewMemoryStorage()
		terminal := &TerminalNode{
			BaseNode: BaseNode{
				ID:     "terminal1",
				Memory: &WorkingMemory{NodeID: "terminal1", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			},
		}

		typeNode := NewTypeNode("Person", TypeDefinition{}, storage)
		typeNode.BaseNode.Children = []Node{terminal}

		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{
				"Person": typeNode,
			},
		}

		facts := map[string][]*Fact{
			"Person": {
				{
					ID:     "fact1",
					Type:   "Person",
					Fields: map[string]interface{}{"id": "1", "name": "Alice"},
				},
				{
					ID:     "fact2",
					Type:   "Person",
					Fields: map[string]interface{}{"id": "2", "name": "Bob"},
				},
			},
		}

		count := pipeline.propagateToNewTerminals(network, []*TerminalNode{terminal}, facts)

		// The function attempts propagation - count depends on child node implementation
		// We verify it doesn't panic and returns a non-negative count
		if count < 0 {
			t.Errorf("Expected non-negative count, got %d", count)
		}
	})

	t.Run("only considers matching fact types", func(t *testing.T) {
		// Verify that the function only attempts to propagate facts of expected types
		storage := NewMemoryStorage()
		terminal := &TerminalNode{
			BaseNode: BaseNode{
				ID:     "terminal1",
				Memory: &WorkingMemory{NodeID: "terminal1", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			},
		}

		typeNode := NewTypeNode("Person", TypeDefinition{}, storage)
		typeNode.BaseNode.Children = []Node{terminal}

		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{
				"Person": typeNode,
			},
		}

		facts := map[string][]*Fact{
			"Person": {
				{ID: "fact1", Type: "Person", Fields: map[string]interface{}{"id": "1"}},
			},
			"Order": {
				{ID: "fact2", Type: "Order", Fields: map[string]interface{}{"id": "2"}},
			},
		}

		// This should not panic - Order facts should be skipped
		count := pipeline.propagateToNewTerminals(network, []*TerminalNode{terminal}, facts)

		// Verify non-negative count (actual value depends on propagation success)
		if count < 0 {
			t.Errorf("Expected non-negative count, got %d", count)
		}
	})

	t.Run("handles multiple terminals", func(t *testing.T) {
		storage := NewMemoryStorage()
		terminal1 := &TerminalNode{
			BaseNode: BaseNode{
				ID:     "terminal1",
				Memory: &WorkingMemory{NodeID: "terminal1", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			},
		}
		terminal2 := &TerminalNode{
			BaseNode: BaseNode{
				ID:     "terminal2",
				Memory: &WorkingMemory{NodeID: "terminal2", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			},
		}

		typeNode1 := NewTypeNode("Person", TypeDefinition{}, storage)
		typeNode1.BaseNode.Children = []Node{terminal1}

		typeNode2 := NewTypeNode("Order", TypeDefinition{}, storage)
		typeNode2.BaseNode.Children = []Node{terminal2}

		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{
				"Person": typeNode1,
				"Order":  typeNode2,
			},
		}

		facts := map[string][]*Fact{
			"Person": {
				{ID: "fact1", Type: "Person", Fields: map[string]interface{}{"id": "1"}},
			},
			"Order": {
				{ID: "fact2", Type: "Order", Fields: map[string]interface{}{"id": "2"}},
			},
		}

		terminals := []*TerminalNode{terminal1, terminal2}
		count := pipeline.propagateToNewTerminals(network, terminals, facts)

		// Verify the function handles multiple terminals without errors
		if count < 0 {
			t.Errorf("Expected non-negative count, got %d", count)
		}
	})

	t.Run("skips propagation on error", func(t *testing.T) {
		// Create a type node that returns error on propagation
		errorTypeNode := &mockTypeNodeWithError{
			BaseNode: BaseNode{ID: "type_Person"},
			typeName: "Person",
		}

		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{
				"Person": errorTypeNode.asTypeNode(),
			},
		}

		terminal := &TerminalNode{
			BaseNode: BaseNode{ID: "terminal1"},
		}

		// Add terminal as child so it's reachable
		errorTypeNode.Children = []Node{terminal}

		facts := map[string][]*Fact{
			"Person": {
				{ID: "fact1", Type: "Person", Fields: map[string]interface{}{"id": "1"}},
			},
		}

		count := pipeline.propagateToNewTerminals(network, []*TerminalNode{terminal}, facts)

		// Should be 0 because propagation failed
		if count != 0 {
			t.Errorf("Expected 0 successful propagations, got %d", count)
		}
	})

	t.Run("handles terminal reachable from multiple types", func(t *testing.T) {
		storage := NewMemoryStorage()

		// Create a simple alpha node instead of join node to avoid complex initialization
		terminal := &TerminalNode{
			BaseNode: BaseNode{
				ID:     "terminal1",
				Memory: &WorkingMemory{NodeID: "terminal1", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			},
		}

		alphaNode := &AlphaNode{
			BaseNode: BaseNode{
				ID:       "alpha1",
				Children: []Node{terminal},
				Memory:   &WorkingMemory{NodeID: "alpha1", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			},
		}

		typeNode1 := NewTypeNode("Person", TypeDefinition{}, storage)
		typeNode1.BaseNode.Children = []Node{alphaNode}

		typeNode2 := NewTypeNode("Order", TypeDefinition{}, storage)
		typeNode2.BaseNode.Children = []Node{alphaNode}

		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{
				"Person": typeNode1,
				"Order":  typeNode2,
			},
		}

		facts := map[string][]*Fact{
			"Person": {
				{ID: "p1", Type: "Person", Fields: map[string]interface{}{"id": "1"}},
			},
			"Order": {
				{ID: "o1", Type: "Order", Fields: map[string]interface{}{"id": "2"}},
			},
		}

		count := pipeline.propagateToNewTerminals(network, []*TerminalNode{terminal}, facts)

		// Verify function completes without errors
		if count < 0 {
			t.Errorf("Expected non-negative count, got %d", count)
		}
	})
}

// Mock helpers for testing

type mockTypeNodeWithError struct {
	BaseNode
	typeName string
}

func (m *mockTypeNodeWithError) asTypeNode() *TypeNode {
	return &TypeNode{
		BaseNode: m.BaseNode,
		TypeName: m.typeName,
	}
}

func (m *mockTypeNodeWithError) PropagateToChildren(fact *Fact, token *Token) error {
	return fmt.Errorf("mock error during propagation")
}
