// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"
)

func TestNewAlphaRuleBuilder(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	arb := NewAlphaRuleBuilder(utils)
	if arb == nil {
		t.Fatal("NewAlphaRuleBuilder returned nil")
	}
	if arb.utils != utils {
		t.Error("AlphaRuleBuilder.utils not set correctly")
	}
}
func TestAlphaRuleBuilder_getVariableInfo(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	arb := NewAlphaRuleBuilder(utils)
	tests := []struct {
		name          string
		variables     []map[string]interface{}
		variableTypes []string
		wantName      string
		wantType      string
	}{
		{
			name: "single variable",
			variables: []map[string]interface{}{
				{"name": "p"},
			},
			variableTypes: []string{"Person"},
			wantName:      "p",
			wantType:      "Person",
		},
		{
			name: "multiple variables - returns first",
			variables: []map[string]interface{}{
				{"name": "p"},
				{"name": "e"},
			},
			variableTypes: []string{"Person", "Employee"},
			wantName:      "p",
			wantType:      "Person",
		},
		{
			name:          "empty variables",
			variables:     []map[string]interface{}{},
			variableTypes: []string{},
			wantName:      "",
			wantType:      "",
		},
		{
			name: "variable without name field",
			variables: []map[string]interface{}{
				{"type": "Person"},
			},
			variableTypes: []string{"Person"},
			wantName:      "",
			wantType:      "Person",
		},
		{
			name: "variable with non-string name",
			variables: []map[string]interface{}{
				{"name": 123},
			},
			variableTypes: []string{"Person"},
			wantName:      "",
			wantType:      "Person",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotType := arb.getVariableInfo(tt.variables, tt.variableTypes)
			if gotName != tt.wantName {
				t.Errorf("getVariableInfo() name = %q, want %q", gotName, tt.wantName)
			}
			if gotType != tt.wantType {
				t.Errorf("getVariableInfo() type = %q, want %q", gotType, tt.wantType)
			}
		})
	}
}
func TestAlphaRuleBuilder_createAlphaNodeWithTerminal(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	arb := NewAlphaRuleBuilder(utils)
	t.Run("create alpha node with terminal", func(t *testing.T) {
		network := NewReteNetwork(storage)
		// Create TypeNode first
		typeNode := NewTypeNode("Person", TypeDefinition{
			Type: "type",
			Name: "Person",
			Fields: []Field{
				{Name: "age", Type: "number"},
			},
		}, storage)
		network.TypeNodes["Person"] = typeNode
		network.RootNode.AddChild(typeNode)
		condition := map[string]interface{}{
			"type":     "comparison",
			"operator": ">",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "age",
			},
			"right": map[string]interface{}{
				"type":  "literal",
				"value": 18,
			},
		}
		action := &Action{
			Type: "print",
			Job:  &JobCall{Name: "print", Args: []interface{}{"Adult person found"}},
		}
		err := arb.createAlphaNodeWithTerminal(
			network,
			"test_rule",
			condition,
			"p",
			"Person",
			action,
		)
		if err != nil {
			t.Fatalf("createAlphaNodeWithTerminal failed: %v", err)
		}
		// Verify AlphaNode was created (now with hash-based ID via AlphaSharingManager)
		if len(network.AlphaNodes) != 1 {
			t.Fatalf("Expected 1 AlphaNode, got %d", len(network.AlphaNodes))
		}
		// Get the alpha node (hash-based ID)
		var alphaNode *AlphaNode
		for _, node := range network.AlphaNodes {
			alphaNode = node
			break
		}
		if alphaNode.VariableName != "p" {
			t.Errorf("AlphaNode.VariableName = %q, want 'p'", alphaNode.VariableName)
		}
		// Verify AlphaNode is connected to TypeNode
		foundChild := false
		for _, child := range typeNode.Children {
			if child == alphaNode {
				foundChild = true
				break
			}
		}
		if !foundChild {
			t.Error("AlphaNode not connected to TypeNode")
		}
		// Verify TerminalNode was created and connected
		if len(alphaNode.Children) != 1 {
			t.Fatalf("AlphaNode should have 1 child, got %d", len(alphaNode.Children))
		}
		terminalNode, ok := alphaNode.Children[0].(*TerminalNode)
		if !ok {
			t.Fatal("AlphaNode child is not a TerminalNode")
		}
		if terminalNode.Action.Type != "print" {
			t.Errorf("TerminalNode action type = %q, want 'print'", terminalNode.Action.Type)
		}
	})
	t.Run("error when type node not found", func(t *testing.T) {
		network := NewReteNetwork(storage)
		condition := map[string]interface{}{}
		action := &Action{Type: "print"}
		err := arb.createAlphaNodeWithTerminal(
			network,
			"test_rule",
			condition,
			"p",
			"NonExistent",
			action,
		)
		if err == nil {
			t.Error("Expected error when TypeNode not found, got nil")
		}
	})
}
func TestAlphaRuleBuilder_CreateAlphaRule(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	arb := NewAlphaRuleBuilder(utils)
	t.Run("create complete alpha rule", func(t *testing.T) {
		network := NewReteNetwork(storage)
		// Create TypeNode
		typeNode := NewTypeNode("Person", TypeDefinition{
			Type: "type",
			Name: "Person",
			Fields: []Field{
				{Name: "age", Type: "number"},
				{Name: "name", Type: "string"},
			},
		}, storage)
		network.TypeNodes["Person"] = typeNode
		network.RootNode.AddChild(typeNode)
		variables := []map[string]interface{}{
			{"name": "p", "dataType": "Person"},
		}
		variableNames := []string{"p"}
		variableTypes := []string{"Person"}
		condition := map[string]interface{}{
			"type":     "comparison",
			"operator": ">=",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "age",
			},
			"right": map[string]interface{}{
				"type":  "literal",
				"value": 21,
			},
		}
		action := &Action{
			Type: "print",
			Job:  &JobCall{Name: "print", Args: []interface{}{"Legal adult"}},
		}
		err := arb.CreateAlphaRule(
			network,
			"adult_rule",
			variables,
			variableNames,
			variableTypes,
			condition,
			action,
		)
		if err != nil {
			t.Fatalf("CreateAlphaRule failed: %v", err)
		}
		// Verify the rule structure (hash-based IDs)
		if len(network.AlphaNodes) != 1 {
			t.Fatalf("Expected 1 AlphaNode, got %d", len(network.AlphaNodes))
		}
		// Get the alpha node
		var alphaNode *AlphaNode
		for _, node := range network.AlphaNodes {
			alphaNode = node
			break
		}
		// Verify graph connectivity: TypeNode -> AlphaNode -> TerminalNode
		foundAlpha := false
		for _, child := range typeNode.Children {
			if child == alphaNode {
				foundAlpha = true
				break
			}
		}
		if !foundAlpha {
			t.Error("TypeNode -> AlphaNode connection missing")
		}
		if len(alphaNode.Children) == 0 {
			t.Fatal("AlphaNode has no children")
		}
		_, isTerminal := alphaNode.Children[0].(*TerminalNode)
		if !isTerminal {
			t.Error("AlphaNode -> TerminalNode connection missing")
		}
	})
	t.Run("create rule with no variables", func(t *testing.T) {
		network := NewReteNetwork(storage)
		// Create a TypeNode with empty type
		typeNode := NewTypeNode("", TypeDefinition{}, storage)
		network.TypeNodes[""] = typeNode
		network.RootNode.AddChild(typeNode)
		err := arb.CreateAlphaRule(
			network,
			"empty_rule",
			[]map[string]interface{}{},
			[]string{},
			[]string{""},
			map[string]interface{}{},
			&Action{Type: "noop"},
		)
		// Empty conditions are valid with the new sharing-based architecture
		if err != nil {
			t.Errorf("CreateAlphaRule failed: %v", err)
		}
		// Verify an AlphaNode was created even for empty conditions
		if len(network.AlphaNodes) != 1 {
			t.Errorf("Expected 1 AlphaNode, got %d", len(network.AlphaNodes))
		}
	})
	t.Run("create rule with multiple variables uses first", func(t *testing.T) {
		network := NewReteNetwork(storage)
		// Create TypeNodes
		personNode := NewTypeNode("Person", TypeDefinition{Name: "Person"}, storage)
		network.TypeNodes["Person"] = personNode
		network.RootNode.AddChild(personNode)
		employeeNode := NewTypeNode("Employee", TypeDefinition{Name: "Employee"}, storage)
		network.TypeNodes["Employee"] = employeeNode
		network.RootNode.AddChild(employeeNode)
		variables := []map[string]interface{}{
			{"name": "p"},
			{"name": "e"},
		}
		variableNames := []string{"p", "e"}
		variableTypes := []string{"Person", "Employee"}
		err := arb.CreateAlphaRule(
			network,
			"multi_var_rule",
			variables,
			variableNames,
			variableTypes,
			map[string]interface{}{},
			&Action{Type: "print"},
		)
		if err != nil {
			t.Fatalf("CreateAlphaRule failed: %v", err)
		}
		// Should use first variable (Person)
		if len(network.AlphaNodes) != 1 {
			t.Fatalf("Expected 1 AlphaNode, got %d", len(network.AlphaNodes))
		}
		// Get the alpha node
		var alphaNode *AlphaNode
		for _, node := range network.AlphaNodes {
			alphaNode = node
			break
		}
		if alphaNode.VariableName != "p" {
			t.Errorf("Expected variable 'p', got %q", alphaNode.VariableName)
		}
		// Should be connected to Person TypeNode, not Employee
		foundChild := false
		for _, child := range personNode.Children {
			if child == alphaNode {
				foundChild = true
				break
			}
		}
		if !foundChild {
			t.Error("AlphaNode should be connected to Person TypeNode")
		}
	})
}
func TestAlphaRuleBuilder_Integration(t *testing.T) {
	// Integration test: Create multiple alpha rules and verify they work together
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	arb := NewAlphaRuleBuilder(utils)
	network := NewReteNetwork(storage)
	// Setup: Create TypeNode
	typeNode := NewTypeNode("Person", TypeDefinition{
		Type: "type",
		Name: "Person",
		Fields: []Field{
			{Name: "age", Type: "number"},
			{Name: "income", Type: "number"},
		},
	}, storage)
	network.TypeNodes["Person"] = typeNode
	network.RootNode.AddChild(typeNode)
	// Rule 1: age > 18
	err := arb.CreateAlphaRule(
		network,
		"adult_rule",
		[]map[string]interface{}{{"name": "p"}},
		[]string{"p"},
		[]string{"Person"},
		map[string]interface{}{
			"type":     "comparison",
			"operator": ">",
			"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
			"right":    map[string]interface{}{"type": "literal", "value": 18},
		},
		&Action{Type: "print", Job: &JobCall{Name: "print", Args: []interface{}{"Adult"}}},
	)
	if err != nil {
		t.Fatalf("Failed to create adult_rule: %v", err)
	}
	// Rule 2: income > 50000
	err = arb.CreateAlphaRule(
		network,
		"high_income_rule",
		[]map[string]interface{}{{"name": "p"}},
		[]string{"p"},
		[]string{"Person"},
		map[string]interface{}{
			"type":     "comparison",
			"operator": ">",
			"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "income"},
			"right":    map[string]interface{}{"type": "literal", "value": 50000},
		},
		&Action{Type: "print", Job: &JobCall{Name: "print", Args: []interface{}{"High income"}}},
	)
	if err != nil {
		t.Fatalf("Failed to create high_income_rule: %v", err)
	}
	// Verify both rules created
	if len(network.AlphaNodes) != 2 {
		t.Errorf("Expected 2 AlphaNodes, got %d", len(network.AlphaNodes))
	}
	// Verify both are connected to the same TypeNode
	if len(typeNode.Children) != 2 {
		t.Errorf("TypeNode should have 2 children (AlphaNodes), got %d", len(typeNode.Children))
	}
	// Verify each AlphaNode has a TerminalNode
	for _, alphaNodeInterface := range network.AlphaNodes {
		alphaNode := alphaNodeInterface
		if len(alphaNode.Children) != 1 {
			t.Errorf("AlphaNode should have 1 child (TerminalNode), got %d", len(alphaNode.Children))
		}
		_, isTerminal := alphaNode.Children[0].(*TerminalNode)
		if !isTerminal {
			t.Error("AlphaNode child should be TerminalNode")
		}
	}
}
