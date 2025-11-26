// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRETENetwork_TypesAndFactsOnly tests creating a RETE network with only types and facts (no rules)
func TestRETENetwork_TypesAndFactsOnly(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "rete_no_rules_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create file with types and facts only (no rules)
	tsdFile := filepath.Join(tempDir, "data.tsd")
	content := `
type Person : <id: string, name: string, age: number>
type Product : <id: string, name: string, price: number>

Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)
Product(id: "PR001", name: "Laptop", price: 999.99)
Product(id: "PR002", name: "Mouse", price: 29.99)
`
	err = os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Build RETE network
	pipeline := NewConstraintPipeline()
	storage := NewMemoryStorage()

	network, err := pipeline.BuildNetworkFromConstraintFile(tsdFile, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}

	if network == nil {
		t.Fatal("Network is nil")
	}

	// Verify TypeNodes were created
	if len(network.TypeNodes) != 2 {
		t.Errorf("Expected 2 TypeNodes, got %d", len(network.TypeNodes))
	}

	// Verify Person TypeNode exists
	personNode, exists := network.TypeNodes["Person"]
	if !exists {
		t.Fatal("Person TypeNode not found")
	}

	if personNode.TypeName != "Person" {
		t.Errorf("Expected TypeNode name 'Person', got '%s'", personNode.TypeName)
	}

	// Verify Product TypeNode exists
	productNode, exists := network.TypeNodes["Product"]
	if !exists {
		t.Fatal("Product TypeNode not found")
	}

	if productNode.TypeName != "Product" {
		t.Errorf("Expected TypeNode name 'Product', got '%s'", productNode.TypeName)
	}

	t.Logf("✅ RETE Network created successfully with types and facts only")
	t.Logf("   - TypeNodes: %d", len(network.TypeNodes))
	t.Logf("   - Person TypeNode: %v", personNode != nil)
	t.Logf("   - Product TypeNode: %v", productNode != nil)
}

// TestRETENetwork_OnlyTypes tests creating a network with only type definitions
func TestRETENetwork_OnlyTypes(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "rete_types_only_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create file with types only
	tsdFile := filepath.Join(tempDir, "types.tsd")
	content := `
type Customer : <id: string, name: string, email: string>
type Order : <id: string, customer_id: string, total: number>
type Product : <id: string, name: string, price: number>
`
	err = os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Build RETE network
	pipeline := NewConstraintPipeline()
	storage := NewMemoryStorage()

	network, err := pipeline.BuildNetworkFromConstraintFile(tsdFile, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}

	if network == nil {
		t.Fatal("Network is nil")
	}

	// Verify all TypeNodes were created
	if len(network.TypeNodes) != 3 {
		t.Errorf("Expected 3 TypeNodes, got %d", len(network.TypeNodes))
	}

	// Verify each type
	expectedTypes := []string{"Customer", "Order", "Product"}
	for _, typeName := range expectedTypes {
		node, exists := network.TypeNodes[typeName]
		if !exists {
			t.Errorf("TypeNode for '%s' not found", typeName)
		} else if node.TypeName != typeName {
			t.Errorf("TypeNode name mismatch: expected '%s', got '%s'", typeName, node.TypeName)
		}
	}

	t.Logf("✅ RETE Network created successfully with types only")
	t.Logf("   - TypeNodes created: %d", len(network.TypeNodes))
	for typeName := range network.TypeNodes {
		t.Logf("     • %s", typeName)
	}
}

// TestRETENetwork_IncrementalTypesAndFacts tests incremental loading of types and facts without rules
func TestRETENetwork_IncrementalTypesAndFacts(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "rete_incremental_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// File 1: Types
	typesFile := filepath.Join(tempDir, "types.tsd")
	typesContent := `
type Person : <id: string, name: string, age: number>
type Company : <id: string, name: string, employees: number>
`
	err = os.WriteFile(typesFile, []byte(typesContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write types file: %v", err)
	}

	// File 2: Facts
	factsFile := filepath.Join(tempDir, "facts.tsd")
	factsContent := `
Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)
Company(id: "C001", name: "TechCorp", employees: 250)
`
	err = os.WriteFile(factsFile, []byte(factsContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write facts file: %v", err)
	}

	// Build network from multiple files
	pipeline := NewConstraintPipeline()
	storage := NewMemoryStorage()

	files := []string{typesFile, factsFile}
	network, err := pipeline.BuildNetworkFromMultipleFiles(files, storage)
	if err != nil {
		t.Fatalf("Failed to build network from multiple files: %v", err)
	}

	if network == nil {
		t.Fatal("Network is nil")
	}

	// Verify TypeNodes
	if len(network.TypeNodes) != 2 {
		t.Errorf("Expected 2 TypeNodes, got %d", len(network.TypeNodes))
	}

	// Verify Person TypeNode
	_, exists := network.TypeNodes["Person"]
	if !exists {
		t.Error("Person TypeNode not found")
	}

	// Verify Company TypeNode
	_, exists = network.TypeNodes["Company"]
	if !exists {
		t.Error("Company TypeNode not found")
	}

	t.Logf("✅ RETE Network created from multiple files (types and facts)")
	t.Logf("   - Files parsed: %d", len(files))
	t.Logf("   - TypeNodes: %d", len(network.TypeNodes))
}

// TestRETENetwork_EmptyFile tests creating a network from an empty file
func TestRETENetwork_EmptyFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "rete_empty_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create empty file
	tsdFile := filepath.Join(tempDir, "empty.tsd")
	err = os.WriteFile(tsdFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Build RETE network
	pipeline := NewConstraintPipeline()
	storage := NewMemoryStorage()

	network, err := pipeline.BuildNetworkFromConstraintFile(tsdFile, storage)

	// Should succeed but create an empty network
	if err != nil {
		t.Fatalf("Failed to build network from empty file: %v", err)
	}

	if network == nil {
		t.Fatal("Network is nil")
	}

	// Verify no TypeNodes were created
	if len(network.TypeNodes) != 0 {
		t.Errorf("Expected 0 TypeNodes for empty file, got %d", len(network.TypeNodes))
	}

	t.Logf("✅ RETE Network created from empty file")
	t.Logf("   - TypeNodes: %d (expected 0)", len(network.TypeNodes))
}

// TestRETENetwork_TypesAndFactsSeparateFiles tests types and facts in separate files
func TestRETENetwork_TypesAndFactsSeparateFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "rete_separate_files_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// File 1: Only types
	typesFile := filepath.Join(tempDir, "schema.tsd")
	typesContent := `
type User : <id: string, username: string, email: string>
`
	err = os.WriteFile(typesFile, []byte(typesContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write types file: %v", err)
	}

	// File 2: Only facts
	factsFile := filepath.Join(tempDir, "users.tsd")
	factsContent := `
User(id: "U001", username: "alice", email: "alice@example.com")
User(id: "U002", username: "bob", email: "bob@example.com")
User(id: "U003", username: "charlie", email: "charlie@example.com")
`
	err = os.WriteFile(factsFile, []byte(factsContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write facts file: %v", err)
	}

	// Build network from both files
	pipeline := NewConstraintPipeline()
	storage := NewMemoryStorage()

	files := []string{typesFile, factsFile}
	network, err := pipeline.BuildNetworkFromMultipleFiles(files, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}

	// Verify network structure
	if len(network.TypeNodes) != 1 {
		t.Errorf("Expected 1 TypeNode, got %d", len(network.TypeNodes))
	}

	userNode, exists := network.TypeNodes["User"]
	if !exists {
		t.Fatal("User TypeNode not found")
	}

	if userNode.TypeName != "User" {
		t.Errorf("Expected TypeNode name 'User', got '%s'", userNode.TypeName)
	}

	t.Logf("✅ Successfully built network from separate type and fact files")
	t.Logf("   - TypeNodes: %d", len(network.TypeNodes))
}
