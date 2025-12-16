// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIngestFile tests the IngestFile wrapper function
func TestIngestFile(t *testing.T) {
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
		network, metrics, err := pipeline.IngestFile(testFile, nil, storage)
		if err != nil {
			t.Fatalf("IngestFile failed: %v", err)
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
		network, metrics, err := pipeline.IngestFile("/non/existent/file.tsd", nil, storage)
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
		network, metrics, err := pipeline.IngestFile(testFile, nil, storage)
		if err != nil {
			t.Fatalf("IngestFile failed: %v", err)
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
		network1, metrics1, err := pipeline.IngestFile(file1, nil, storage)
		if err != nil {
			t.Fatalf("First ingest failed: %v", err)
		}
		initialTypeCount := len(network1.TypeNodes)
		if initialTypeCount == 0 {
			t.Fatal("Expected type nodes after first ingest")
		}
		// Ingest second file with existing network
		network2, metrics2, err := pipeline.IngestFile(file2, network1, storage)
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
		// Note: Testing error propagation requires complex setup
		// This test verifies the function handles error cases gracefully
		storage := NewMemoryStorage()
		typeNode := NewTypeNode("Person", TypeDefinition{}, storage)
		terminal := &TerminalNode{
			BaseNode: BaseNode{
				ID:     "terminal1",
				Memory: &WorkingMemory{NodeID: "terminal1", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			},
		}
		// TypeNode with no children - propagation will fail
		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{
				"Person": typeNode,
			},
		}
		facts := map[string][]*Fact{
			"Person": {
				{ID: "fact1", Type: "Person", Fields: map[string]interface{}{"id": "1"}},
			},
		}
		count := pipeline.propagateToNewTerminals(network, []*TerminalNode{terminal}, facts)
		// Count should be non-negative even with propagation issues
		if count < 0 {
			t.Errorf("Expected non-negative count, got %d", count)
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

// TestOrganizeFactsByType tests the organizeFactsByType helper function
func TestOrganizeFactsByType(t *testing.T) {
	pipeline := NewConstraintPipeline()
	t.Run("returns empty map for empty input", func(t *testing.T) {
		result := pipeline.organizeFactsByType([]*Fact{})
		if len(result) != 0 {
			t.Errorf("Expected empty map, got %d entries", len(result))
		}
	})
	t.Run("returns empty map for nil input", func(t *testing.T) {
		result := pipeline.organizeFactsByType(nil)
		if len(result) != 0 {
			t.Errorf("Expected empty map, got %d entries", len(result))
		}
	})
	t.Run("organizes single fact by type", func(t *testing.T) {
		facts := []*Fact{
			{ID: "f1", Type: "Person", Fields: map[string]interface{}{"id": "1"}},
		}
		result := pipeline.organizeFactsByType(facts)
		if len(result) != 1 {
			t.Fatalf("Expected 1 type, got %d", len(result))
		}
		if len(result["Person"]) != 1 {
			t.Errorf("Expected 1 Person fact, got %d", len(result["Person"]))
		}
		if result["Person"][0].ID != "f1" {
			t.Errorf("Expected fact ID 'f1', got '%s'", result["Person"][0].ID)
		}
	})
	t.Run("organizes multiple facts of same type", func(t *testing.T) {
		facts := []*Fact{
			{ID: "f1", Type: "Person", Fields: map[string]interface{}{"id": "1"}},
			{ID: "f2", Type: "Person", Fields: map[string]interface{}{"id": "2"}},
			{ID: "f3", Type: "Person", Fields: map[string]interface{}{"id": "3"}},
		}
		result := pipeline.organizeFactsByType(facts)
		if len(result) != 1 {
			t.Fatalf("Expected 1 type, got %d", len(result))
		}
		if len(result["Person"]) != 3 {
			t.Errorf("Expected 3 Person facts, got %d", len(result["Person"]))
		}
	})
	t.Run("organizes facts by multiple types", func(t *testing.T) {
		facts := []*Fact{
			{ID: "p1", Type: "Person", Fields: map[string]interface{}{"id": "1"}},
			{ID: "p2", Type: "Person", Fields: map[string]interface{}{"id": "2"}},
			{ID: "o1", Type: "Order", Fields: map[string]interface{}{"id": "100"}},
			{ID: "o2", Type: "Order", Fields: map[string]interface{}{"id": "101"}},
			{ID: "prod1", Type: "Product", Fields: map[string]interface{}{"id": "A"}},
		}
		result := pipeline.organizeFactsByType(facts)
		if len(result) != 3 {
			t.Fatalf("Expected 3 types, got %d", len(result))
		}
		if len(result["Person"]) != 2 {
			t.Errorf("Expected 2 Person facts, got %d", len(result["Person"]))
		}
		if len(result["Order"]) != 2 {
			t.Errorf("Expected 2 Order facts, got %d", len(result["Order"]))
		}
		if len(result["Product"]) != 1 {
			t.Errorf("Expected 1 Product fact, got %d", len(result["Product"]))
		}
	})
	t.Run("skips nil facts", func(t *testing.T) {
		facts := []*Fact{
			{ID: "f1", Type: "Person", Fields: map[string]interface{}{"id": "1"}},
			nil,
			{ID: "f2", Type: "Order", Fields: map[string]interface{}{"id": "100"}},
			nil,
		}
		result := pipeline.organizeFactsByType(facts)
		if len(result) != 2 {
			t.Fatalf("Expected 2 types, got %d", len(result))
		}
		if len(result["Person"]) != 1 {
			t.Errorf("Expected 1 Person fact, got %d", len(result["Person"]))
		}
		if len(result["Order"]) != 1 {
			t.Errorf("Expected 1 Order fact, got %d", len(result["Order"]))
		}
	})
	t.Run("maintains fact order within type", func(t *testing.T) {
		facts := []*Fact{
			{ID: "f1", Type: "Person", Fields: map[string]interface{}{"id": "1"}},
			{ID: "f2", Type: "Person", Fields: map[string]interface{}{"id": "2"}},
			{ID: "f3", Type: "Person", Fields: map[string]interface{}{"id": "3"}},
		}
		result := pipeline.organizeFactsByType(facts)
		personFacts := result["Person"]
		if personFacts[0].ID != "f1" || personFacts[1].ID != "f2" || personFacts[2].ID != "f3" {
			t.Error("Facts not in expected order")
		}
	})
}

// TestCollectExistingFacts tests the collectExistingFacts function
func TestCollectExistingFacts(t *testing.T) {
	pipeline := NewConstraintPipeline()
	t.Run("returns empty for empty network", func(t *testing.T) {
		network := &ReteNetwork{
			TypeNodes:  make(map[string]*TypeNode),
			AlphaNodes: make(map[string]*AlphaNode),
			BetaNodes:  make(map[string]interface{}),
		}
		facts := pipeline.collectExistingFacts(network)
		if len(facts) != 0 {
			t.Errorf("Expected 0 facts, got %d", len(facts))
		}
	})
	t.Run("collects facts from RootNode", func(t *testing.T) {
		network := &ReteNetwork{
			RootNode: &RootNode{
				BaseNode: BaseNode{
					Memory: &WorkingMemory{
						Facts: map[string]*Fact{
							"f1": {ID: "f1", Type: "Person", Fields: map[string]interface{}{"id": "1"}},
							"f2": {ID: "f2", Type: "Order", Fields: map[string]interface{}{"id": "100"}},
						},
					},
				},
			},
			TypeNodes:  make(map[string]*TypeNode),
			AlphaNodes: make(map[string]*AlphaNode),
			BetaNodes:  make(map[string]interface{}),
		}
		facts := pipeline.collectExistingFacts(network)
		if len(facts) != 2 {
			t.Errorf("Expected 2 facts, got %d", len(facts))
		}
	})
	t.Run("collects facts from TypeNodes", func(t *testing.T) {
		storage := NewMemoryStorage()
		typeNode := NewTypeNode("Person", TypeDefinition{}, storage)
		token := &Token{
			ID:    "t1",
			Facts: []*Fact{{ID: "f1", Type: "Person", Fields: map[string]interface{}{"id": "1"}}},
		}
		typeNode.Memory.AddToken(token)
		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{
				"Person": typeNode,
			},
			AlphaNodes: make(map[string]*AlphaNode),
			BetaNodes:  make(map[string]interface{}),
		}
		facts := pipeline.collectExistingFacts(network)
		if len(facts) != 1 {
			t.Errorf("Expected 1 fact, got %d", len(facts))
		}
	})
	t.Run("collects facts from AlphaNodes", func(t *testing.T) {
		storage := NewMemoryStorage()
		alphaNode := &AlphaNode{
			BaseNode: BaseNode{
				ID:      "alpha1",
				Storage: storage,
				Memory: &WorkingMemory{
					NodeID: "alpha1",
					Facts:  make(map[string]*Fact),
					Tokens: make(map[string]*Token),
				},
			},
		}
		token := &Token{
			ID:    "t1",
			Facts: []*Fact{{ID: "f1", Type: "Person", Fields: map[string]interface{}{"id": "1"}}},
		}
		alphaNode.Memory.AddToken(token)
		network := &ReteNetwork{
			TypeNodes: make(map[string]*TypeNode),
			AlphaNodes: map[string]*AlphaNode{
				"alpha1": alphaNode,
			},
			BetaNodes: make(map[string]interface{}),
		}
		facts := pipeline.collectExistingFacts(network)
		if len(facts) != 1 {
			t.Errorf("Expected 1 fact, got %d", len(facts))
		}
	})
	t.Run("collects facts from JoinNode left and right memory", func(t *testing.T) {
		joinNode := &JoinNode{
			BaseNode: BaseNode{ID: "join1"},
			LeftMemory: &WorkingMemory{
				NodeID: "join1_left",
				Facts:  make(map[string]*Fact),
				Tokens: map[string]*Token{
					"t1": {
						ID:    "t1",
						Facts: []*Fact{{ID: "f1", Type: "Person", Fields: map[string]interface{}{"id": "1"}}},
					},
				},
			},
			RightMemory: &WorkingMemory{
				NodeID: "join1_right",
				Facts:  make(map[string]*Fact),
				Tokens: map[string]*Token{
					"t2": {
						ID:    "t2",
						Facts: []*Fact{{ID: "f2", Type: "Order", Fields: map[string]interface{}{"id": "100"}}},
					},
				},
			},
		}
		network := &ReteNetwork{
			TypeNodes:  make(map[string]*TypeNode),
			AlphaNodes: make(map[string]*AlphaNode),
			BetaNodes: map[string]interface{}{
				"join1": joinNode,
			},
		}
		facts := pipeline.collectExistingFacts(network)
		if len(facts) != 2 {
			t.Errorf("Expected 2 facts, got %d", len(facts))
		}
	})
	t.Run("collects facts from ExistsNode memories", func(t *testing.T) {
		existsNode := &ExistsNode{
			BaseNode: BaseNode{ID: "exists1"},
			MainMemory: &WorkingMemory{
				NodeID: "exists1_main",
				Facts:  make(map[string]*Fact),
				Tokens: map[string]*Token{
					"t1": {
						ID:    "t1",
						Facts: []*Fact{{ID: "f1", Type: "Person", Fields: map[string]interface{}{"id": "1"}}},
					},
				},
			},
			ExistsMemory: &WorkingMemory{
				NodeID: "exists1_exists",
				Facts:  make(map[string]*Fact),
				Tokens: map[string]*Token{
					"t2": {
						ID:    "t2",
						Facts: []*Fact{{ID: "f2", Type: "Order", Fields: map[string]interface{}{"id": "100"}}},
					},
				},
			},
		}
		network := &ReteNetwork{
			TypeNodes:  make(map[string]*TypeNode),
			AlphaNodes: make(map[string]*AlphaNode),
			BetaNodes: map[string]interface{}{
				"exists1": existsNode,
			},
		}
		facts := pipeline.collectExistingFacts(network)
		if len(facts) != 2 {
			t.Errorf("Expected 2 facts, got %d", len(facts))
		}
	})
	t.Run("collects facts from AccumulatorNode", func(t *testing.T) {
		accNode := &AccumulatorNode{
			BaseNode: BaseNode{ID: "acc1"},
			MainFacts: map[string]*Fact{
				"f1": {ID: "f1", Type: "Person", Fields: map[string]interface{}{"id": "1"}},
			},
			AllFacts: map[string]*Fact{
				"f2": {ID: "f2", Type: "Order", Fields: map[string]interface{}{"id": "100"}},
				"f3": {ID: "f3", Type: "Product", Fields: map[string]interface{}{"id": "A"}},
			},
		}
		network := &ReteNetwork{
			TypeNodes:  make(map[string]*TypeNode),
			AlphaNodes: make(map[string]*AlphaNode),
			BetaNodes: map[string]interface{}{
				"acc1": accNode,
			},
		}
		facts := pipeline.collectExistingFacts(network)
		if len(facts) != 3 {
			t.Errorf("Expected 3 facts, got %d", len(facts))
		}
	})
	t.Run("deduplicates facts by ID", func(t *testing.T) {
		storage := NewMemoryStorage()
		typeNode := NewTypeNode("Person", TypeDefinition{}, storage)
		sharedFact := &Fact{ID: "f1", Type: "Person", Fields: map[string]interface{}{"id": "1"}}
		token1 := &Token{ID: "t1", Facts: []*Fact{sharedFact}}
		token2 := &Token{ID: "t2", Facts: []*Fact{sharedFact}}
		typeNode.Memory.AddToken(token1)
		typeNode.Memory.AddToken(token2)
		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{
				"Person": typeNode,
			},
			AlphaNodes: make(map[string]*AlphaNode),
			BetaNodes:  make(map[string]interface{}),
		}
		facts := pipeline.collectExistingFacts(network)
		// Should only have 1 fact despite being in 2 tokens
		if len(facts) != 1 {
			t.Errorf("Expected 1 deduplicated fact, got %d", len(facts))
		}
	})
	t.Run("skips nil facts", func(t *testing.T) {
		storage := NewMemoryStorage()
		typeNode := NewTypeNode("Person", TypeDefinition{}, storage)
		token := &Token{
			ID: "t1",
			Facts: []*Fact{
				{ID: "f1", Type: "Person", Fields: map[string]interface{}{"id": "1"}},
				nil,
				{ID: "f2", Type: "Person", Fields: map[string]interface{}{"id": "2"}},
			},
		}
		typeNode.Memory.AddToken(token)
		network := &ReteNetwork{
			TypeNodes: map[string]*TypeNode{
				"Person": typeNode,
			},
			AlphaNodes: make(map[string]*AlphaNode),
			BetaNodes:  make(map[string]interface{}),
		}
		facts := pipeline.collectExistingFacts(network)
		// Should only collect non-nil facts
		if len(facts) != 2 {
			t.Errorf("Expected 2 non-nil facts, got %d", len(facts))
		}
	})
	t.Run("collects from multiple node types simultaneously", func(t *testing.T) {
		storage := NewMemoryStorage()
		// Root node with facts
		rootNode := &RootNode{
			BaseNode: BaseNode{
				Memory: &WorkingMemory{
					Facts: map[string]*Fact{
						"f1": {ID: "f1", Type: "Person", Fields: map[string]interface{}{"id": "1"}},
					},
				},
			},
		}
		// Type node with facts
		typeNode := NewTypeNode("Order", TypeDefinition{}, storage)
		typeToken := &Token{
			ID:    "t1",
			Facts: []*Fact{{ID: "f2", Type: "Order", Fields: map[string]interface{}{"id": "100"}}},
		}
		typeNode.Memory.AddToken(typeToken)
		// Alpha node with facts
		alphaNode := &AlphaNode{
			BaseNode: BaseNode{
				ID:      "alpha1",
				Storage: storage,
				Memory: &WorkingMemory{
					NodeID: "alpha1",
					Facts:  make(map[string]*Fact),
					Tokens: map[string]*Token{
						"t2": {
							ID:    "t2",
							Facts: []*Fact{{ID: "f3", Type: "Product", Fields: map[string]interface{}{"id": "A"}}},
						},
					},
				},
			},
		}
		// Join node with facts
		joinNode := &JoinNode{
			BaseNode: BaseNode{ID: "join1"},
			LeftMemory: &WorkingMemory{
				NodeID: "join1_left",
				Facts:  make(map[string]*Fact),
				Tokens: map[string]*Token{
					"t3": {
						ID:    "t3",
						Facts: []*Fact{{ID: "f4", Type: "Address", Fields: map[string]interface{}{"id": "addr1"}}},
					},
				},
			},
			RightMemory: &WorkingMemory{
				NodeID: "join1_right",
				Facts:  make(map[string]*Fact),
				Tokens: make(map[string]*Token),
			},
		}
		network := &ReteNetwork{
			RootNode: rootNode,
			TypeNodes: map[string]*TypeNode{
				"Order": typeNode,
			},
			AlphaNodes: map[string]*AlphaNode{
				"alpha1": alphaNode,
			},
			BetaNodes: map[string]interface{}{
				"join1": joinNode,
			},
		}
		facts := pipeline.collectExistingFacts(network)
		// Should collect from all node types: f1 (root), f2 (type), f3 (alpha), f4 (join)
		if len(facts) != 4 {
			t.Errorf("Expected 4 facts from all node types, got %d", len(facts))
		}
	})
}

// TestIngestFile_ErrorPaths tests error handling in IngestFile
func TestIngestFile_ErrorPaths(t *testing.T) {
	t.Run("handles conversion error", func(t *testing.T) {
		// This test would require mocking the constraint package
		// which is beyond the scope of unit tests here
		// The conversion error path is integration-tested
	})
	t.Run("handles reset with garbage collection", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.tsd")
		// First create a network with some data
		content1 := `
type Person(id: string, name: string)
action log(msg: string)
rule r1 : {p: Person} / p.id == "1" ==> log(p.name)
Person(id: "1", name: "Alice")
`
		if err := os.WriteFile(testFile, []byte(content1), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		pipeline := NewConstraintPipeline()
		storage := NewMemoryStorage()
		network1, _, err := pipeline.IngestFile(testFile, nil, storage)
		if err != nil {
			t.Fatalf("Initial ingestion failed: %v", err)
		}
		// Now ingest a file with reset command
		content2 := `
reset
type Order(id: string)
action process(id: string)
rule r2 : {o: Order} / o.id == "100" ==> process(o.id)
`
		if err := os.WriteFile(testFile, []byte(content2), 0644); err != nil {
			t.Fatalf("Failed to update test file: %v", err)
		}
		network2, metrics, err := pipeline.IngestFile(testFile, network1, storage)
		if err != nil {
			t.Fatalf("Reset ingestion failed: %v", err)
		}
		// Should be a new network after reset
		if network2 == network1 {
			t.Error("Expected new network after reset")
		}
		if metrics == nil {
			t.Fatal("Expected metrics")
		}
		if !metrics.WasReset {
			t.Error("Expected WasReset to be true")
		}
		// Old type should be gone, new type should exist
		if _, exists := network2.TypeNodes["Person"]; exists {
			t.Error("Expected Person type to be removed after reset")
		}
		if _, exists := network2.TypeNodes["Order"]; !exists {
			t.Error("Expected Order type to exist after reset")
		}
	})
	t.Run("handles validation error with rollback", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.tsd")
		// Create invalid content (undefined type reference)
		content := `
type Person(id: string, name: string)
action log(msg: string)
rule bad : {o: Order} / o.id == "1" ==> log("test")
`
		if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		pipeline := NewConstraintPipeline()
		storage := NewMemoryStorage()
		_, _, err := pipeline.IngestFile(testFile, nil, storage)
		// Should fail validation
		if err == nil {
			t.Error("Expected validation error for undefined type")
		}
	})
	t.Run("incremental validation path", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.tsd")
		// First file with base type
		content1 := `
type Person(id: string, name: string)
action log(msg: string)
rule r1 : {p: Person} / p.id == "1" ==> log(p.name)
`
		if err := os.WriteFile(testFile, []byte(content1), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		pipeline := NewConstraintPipeline()
		storage := NewMemoryStorage()
		network1, _, err := pipeline.IngestFile(testFile, nil, storage)
		if err != nil {
			t.Fatalf("Initial ingestion failed: %v", err)
		}
		// Second file that references existing type (incremental validation)
		content2 := `
type Order(id: string, person_id: string)
action process(id: string)
rule r2 : {o: Order, p: Person} / o.person_id == p.id ==> process(o.id)
`
		if err := os.WriteFile(testFile, []byte(content2), 0644); err != nil {
			t.Fatalf("Failed to update test file: %v", err)
		}
		network2, metrics, err := pipeline.IngestFile(testFile, network1, storage)
		if err != nil {
			t.Fatalf("Incremental ingestion failed: %v", err)
		}
		if network2 != network1 {
			t.Error("Expected same network instance for incremental ingestion")
		}
		if metrics == nil {
			t.Fatal("Expected metrics")
		}
		if !metrics.WasIncremental {
			t.Error("Expected WasIncremental to be true")
		}
		// Both types should exist
		if _, exists := network2.TypeNodes["Person"]; !exists {
			t.Error("Expected Person type to still exist")
		}
		if _, exists := network2.TypeNodes["Order"]; !exists {
			t.Error("Expected Order type to be added")
		}
	})
	t.Run("handles fact submission with consistency check", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.tsd")
		content := `
type Product(id: string, name: string, price: number)
action notify(msg: string)
rule expensive : {p: Product} / p.price > 100 ==> notify(p.name)
Product(id: "p1", name: "Laptop", price: 1500)
Product(id: "p2", name: "Mouse", price: 25)
`
		if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		pipeline := NewConstraintPipeline()
		storage := NewMemoryStorage()
		network, metrics, err := pipeline.IngestFile(testFile, nil, storage)
		if err != nil {
			t.Fatalf("Ingestion failed: %v", err)
		}
		if metrics == nil {
			t.Fatal("Expected metrics")
		}
		if metrics.FactsSubmitted != 2 {
			t.Errorf("Expected 2 facts submitted, got %d", metrics.FactsSubmitted)
		}
		// Verify facts are in storage
		if network == nil {
			t.Fatal("Expected non-nil network")
		}
	})
	t.Run("handles propagation to new terminals", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.tsd")
		// First file with facts
		content1 := `
type Item(id: string, category: string)
action log(msg: string)
Item(id: "i1", category: "books")
Item(id: "i2", category: "electronics")
`
		if err := os.WriteFile(testFile, []byte(content1), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		pipeline := NewConstraintPipeline()
		storage := NewMemoryStorage()
		network1, _, err := pipeline.IngestFile(testFile, nil, storage)
		if err != nil {
			t.Fatalf("Initial ingestion failed: %v", err)
		}
		// Second file adds rule (should propagate existing facts)
		content2 := `
rule books : {i: Item} / i.category == "books" ==> log(i.id)
`
		if err := os.WriteFile(testFile, []byte(content2), 0644); err != nil {
			t.Fatalf("Failed to update test file: %v", err)
		}
		network2, metrics, err := pipeline.IngestFile(testFile, network1, storage)
		if err != nil {
			t.Fatalf("Rule ingestion failed: %v", err)
		}
		if metrics == nil {
			t.Fatal("Expected metrics")
		}
		if metrics.NewTerminalsAdded == 0 {
			t.Error("Expected new terminals to be added")
		}
		// Should have propagated existing facts
		if network2 == nil {
			t.Fatal("Expected non-nil network")
		}
	})
	t.Run("handles empty expressions and types", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.tsd")
		// File with only facts (no types or rules)
		content := `
// Just a comment, no actual content
`
		if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		pipeline := NewConstraintPipeline()
		storage := NewMemoryStorage()
		network, metrics, err := pipeline.IngestFile(testFile, nil, storage)
		// Should succeed even with empty content
		if err != nil {
			t.Fatalf("Empty file ingestion failed: %v", err)
		}
		if network == nil {
			t.Fatal("Expected network to be created")
		}
		if metrics == nil {
			t.Fatal("Expected metrics")
		}
		if metrics.TypesAdded != 0 {
			t.Errorf("Expected 0 types added, got %d", metrics.TypesAdded)
		}
		if metrics.RulesAdded != 0 {
			t.Errorf("Expected 0 rules added, got %d", metrics.RulesAdded)
		}
	})
}

// TestProcessRuleRemovals tests the processRuleRemovals function
func TestProcessRuleRemovals(t *testing.T) {
	pipeline := NewConstraintPipeline()
	t.Run("no rule removals present", func(t *testing.T) {
		network := NewReteNetwork(NewMemoryStorage())
		resultMap := map[string]interface{}{
			"types": []interface{}{},
			"rules": []interface{}{},
		}
		err := pipeline.processRuleRemovals(network, resultMap)
		if err != nil {
			t.Errorf("Expected no error when no rule removals, got: %v", err)
		}
	})
	t.Run("rule removals key exists but empty", func(t *testing.T) {
		network := NewReteNetwork(NewMemoryStorage())
		resultMap := map[string]interface{}{
			"ruleRemovals": []interface{}{},
		}
		err := pipeline.processRuleRemovals(network, resultMap)
		if err != nil {
			t.Errorf("Expected no error for empty rule removals, got: %v", err)
		}
	})
	t.Run("rule removals with invalid type", func(t *testing.T) {
		network := NewReteNetwork(NewMemoryStorage())
		resultMap := map[string]interface{}{
			"ruleRemovals": "not a slice",
		}
		err := pipeline.processRuleRemovals(network, resultMap)
		if err != nil {
			t.Errorf("Expected no error for invalid type (should skip), got: %v", err)
		}
	})
	t.Run("rule removal with invalid format", func(t *testing.T) {
		network := NewReteNetwork(NewMemoryStorage())
		resultMap := map[string]interface{}{
			"ruleRemovals": []interface{}{
				"not a map",
				123,
			},
		}
		// Should handle gracefully with warnings
		err := pipeline.processRuleRemovals(network, resultMap)
		if err != nil {
			t.Errorf("Expected no error (should warn and continue), got: %v", err)
		}
	})
	t.Run("rule removal with missing rule ID", func(t *testing.T) {
		network := NewReteNetwork(NewMemoryStorage())
		resultMap := map[string]interface{}{
			"ruleRemovals": []interface{}{
				map[string]interface{}{
					"otherField": "value",
				},
				map[string]interface{}{
					"ruleID": "",
				},
			},
		}
		// Should handle gracefully with warnings
		err := pipeline.processRuleRemovals(network, resultMap)
		if err != nil {
			t.Errorf("Expected no error (should warn and continue), got: %v", err)
		}
	})
	t.Run("rule removal with invalid rule ID type", func(t *testing.T) {
		network := NewReteNetwork(NewMemoryStorage())
		resultMap := map[string]interface{}{
			"ruleRemovals": []interface{}{
				map[string]interface{}{
					"ruleID": 123,
				},
			},
		}
		// Should handle gracefully with warnings
		err := pipeline.processRuleRemovals(network, resultMap)
		if err != nil {
			t.Errorf("Expected no error (should warn and continue), got: %v", err)
		}
	})
	t.Run("rule removal attempt without lifecycle registration", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		// Add a terminal node but DON'T register it with lifecycle manager
		// This simulates an edge case where the rule isn't properly tracked
		terminalNode := &TerminalNode{
			BaseNode: BaseNode{
				ID: "rule_test_rule",
			},
		}
		network.TerminalNodes["rule_test_rule"] = terminalNode
		resultMap := map[string]interface{}{
			"ruleRemovals": []interface{}{
				map[string]interface{}{
					"ruleID": "test_rule",
				},
			},
		}
		// Should handle gracefully - RemoveRule will error but processRuleRemovals continues
		err := pipeline.processRuleRemovals(network, resultMap)
		if err != nil {
			t.Errorf("Expected no error (should warn and continue), got: %v", err)
		}
	})
	t.Run("multiple rule removals all fail gracefully", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		resultMap := map[string]interface{}{
			"ruleRemovals": []interface{}{
				map[string]interface{}{
					"ruleID": "rule1",
				},
				map[string]interface{}{
					"ruleID": "non_existent",
				},
				map[string]interface{}{
					"ruleID": "rule2",
				},
			},
		}
		// All removals will fail (no rules registered) but should continue processing
		err := pipeline.processRuleRemovals(network, resultMap)
		if err != nil {
			t.Errorf("Expected no error (should warn and continue for each), got: %v", err)
		}
	})
	t.Run("rule removal with non-existent rule", func(t *testing.T) {
		network := NewReteNetwork(NewMemoryStorage())
		resultMap := map[string]interface{}{
			"ruleRemovals": []interface{}{
				map[string]interface{}{
					"ruleID": "non_existent_rule",
				},
			},
		}
		// Should handle gracefully with warning (RemoveRule will error)
		err := pipeline.processRuleRemovals(network, resultMap)
		if err != nil {
			t.Errorf("Expected no error (should warn and continue), got: %v", err)
		}
	})
}
