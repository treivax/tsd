// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

func TestNewJoinRuleBuilder(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)

	jrb := NewJoinRuleBuilder(utils)

	if jrb == nil {
		t.Fatal("NewJoinRuleBuilder returned nil")
	}

	if jrb.utils != utils {
		t.Error("JoinRuleBuilder.utils not set correctly")
	}
}

func TestJoinRuleBuilder_createBinaryJoinRule(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	jrb := NewJoinRuleBuilder(utils)

	t.Run("create simple binary join", func(t *testing.T) {
		network := NewReteNetwork(storage)

		// Create TypeNodes
		personNode := NewTypeNode("Person", TypeDefinition{
			Type: "type",
			Name: "Person",
			Fields: []Field{
				{Name: "id", Type: "number"},
			},
		}, storage)
		network.TypeNodes["Person"] = personNode
		network.RootNode.AddChild(personNode)

		employeeNode := NewTypeNode("Employee", TypeDefinition{
			Type: "type",
			Name: "Employee",
			Fields: []Field{
				{Name: "person_id", Type: "number"},
			},
		}, storage)
		network.TypeNodes["Employee"] = employeeNode
		network.RootNode.AddChild(employeeNode)

		variableNames := []string{"p", "e"}
		variableTypes := []string{"Person", "Employee"}

		condition := map[string]interface{}{
			"type":     "comparison",
			"operator": "==",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "id",
			},
			"right": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "e",
				"field":  "person_id",
			},
		}

		terminalNode := utils.CreateTerminalNode(network, "test_rule", &Action{Type: "print"})

		err := jrb.createBinaryJoinRule(
			network,
			"test_rule",
			variableNames,
			variableTypes,
			condition,
			terminalNode,
		)

		if err != nil {
			t.Fatalf("createBinaryJoinRule failed: %v", err)
		}

		// Verify JoinNode was created
		joinNode, exists := network.BetaNodes["test_rule_join"]
		if !exists {
			t.Fatal("JoinNode not created")
		}

		joinNodeTyped, ok := joinNode.(*JoinNode)
		if !ok {
			t.Fatal("BetaNode is not a JoinNode")
		}

		// Verify JoinNode configuration
		if len(joinNodeTyped.LeftVariables) != 1 || joinNodeTyped.LeftVariables[0] != "p" {
			t.Errorf("LeftVariables = %v, want ['p']", joinNodeTyped.LeftVariables)
		}

		if len(joinNodeTyped.RightVariables) != 1 || joinNodeTyped.RightVariables[0] != "e" {
			t.Errorf("RightVariables = %v, want ['e']", joinNodeTyped.RightVariables)
		}

		// Verify TerminalNode connection
		if len(joinNodeTyped.Children) != 1 {
			t.Fatalf("JoinNode should have 1 child (TerminalNode), got %d", len(joinNodeTyped.Children))
		}

		if joinNodeTyped.Children[0] != terminalNode {
			t.Error("JoinNode not connected to TerminalNode")
		}

		// Verify TypeNode connections (via pass-through alphas)
		if len(personNode.Children) == 0 {
			t.Error("Person TypeNode should have children")
		}

		if len(employeeNode.Children) == 0 {
			t.Error("Employee TypeNode should have children")
		}
	})

	t.Run("with beta sharing enabled", func(t *testing.T) {
		network := NewReteNetwork(storage)
		config := BetaSharingConfig{Enabled: true, HashCacheSize: 1000, MaxSharedNodes: 10000}
		lifecycle := NewLifecycleManager()
		network.BetaSharingRegistry = NewBetaSharingRegistry(config, lifecycle)

		personNode := NewTypeNode("Person", TypeDefinition{Name: "Person"}, storage)
		network.TypeNodes["Person"] = personNode
		network.RootNode.AddChild(personNode)

		employeeNode := NewTypeNode("Employee", TypeDefinition{Name: "Employee"}, storage)
		network.TypeNodes["Employee"] = employeeNode
		network.RootNode.AddChild(employeeNode)

		variableNames := []string{"p", "e"}
		variableTypes := []string{"Person", "Employee"}
		condition := map[string]interface{}{"type": "comparison"}
		terminalNode := utils.CreateTerminalNode(network, "rule1", &Action{Type: "print"})

		err := jrb.createBinaryJoinRule(network, "rule1", variableNames, variableTypes, condition, terminalNode)
		if err != nil {
			t.Fatalf("createBinaryJoinRule failed: %v", err)
		}

		// Create another rule with same pattern - should reuse JoinNode
		terminalNode2 := utils.CreateTerminalNode(network, "rule2", &Action{Type: "print"})
		err = jrb.createBinaryJoinRule(network, "rule2", variableNames, variableTypes, condition, terminalNode2)
		if err != nil {
			t.Fatalf("createBinaryJoinRule failed: %v", err)
		}

		// Should still have only 1 JoinNode in registry (shared)
		// But 2 entries in BetaNodes (rule1_join and rule2_join)
		if len(network.BetaNodes) != 2 {
			t.Errorf("Expected 2 BetaNode entries, got %d", len(network.BetaNodes))
		}
	})

	t.Run("error on missing type node", func(t *testing.T) {
		network := NewReteNetwork(storage)

		variableNames := []string{"p", "e"}
		variableTypes := []string{"NonExistent", "Employee"}
		condition := map[string]interface{}{}
		terminalNode := utils.CreateTerminalNode(network, "bad_rule", &Action{Type: "print"})

		err := jrb.createBinaryJoinRule(network, "bad_rule", variableNames, variableTypes, condition, terminalNode)

		if err == nil {
			t.Error("Expected error for missing TypeNode, got nil")
		}
	})
}

func TestJoinRuleBuilder_createCascadeJoinRuleLegacy(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	jrb := NewJoinRuleBuilder(utils)

	t.Run("create 3-variable cascade", func(t *testing.T) {
		network := NewReteNetwork(storage)

		// Create TypeNodes
		types := []string{"Person", "Employee", "Department"}
		for _, typeName := range types {
			typeNode := NewTypeNode(typeName, TypeDefinition{Name: typeName}, storage)
			network.TypeNodes[typeName] = typeNode
			network.RootNode.AddChild(typeNode)
		}

		variableNames := []string{"p", "e", "d"}
		variableTypes := []string{"Person", "Employee", "Department"}
		condition := map[string]interface{}{"type": "comparison"}
		terminalNode := utils.CreateTerminalNode(network, "cascade_rule", &Action{Type: "print"})

		err := jrb.createCascadeJoinRuleLegacy(
			network,
			"cascade_rule",
			variableNames,
			variableTypes,
			condition,
			terminalNode,
		)

		if err != nil {
			t.Fatalf("createCascadeJoinRuleLegacy failed: %v", err)
		}

		// Should create 2 JoinNodes for 3 variables
		// cascade_rule_join_0_1 (p ⋈ e)
		// cascade_rule_join_2 ((p ⋈ e) ⋈ d)
		expectedNodes := []string{"cascade_rule_join_0_1", "cascade_rule_join_2"}
		for _, nodeID := range expectedNodes {
			if _, exists := network.BetaNodes[nodeID]; !exists {
				t.Errorf("Expected JoinNode %s not found", nodeID)
			}
		}

		if len(network.BetaNodes) != 2 {
			t.Errorf("Expected 2 JoinNodes for cascade, got %d", len(network.BetaNodes))
		}

		// Verify the final JoinNode is connected to TerminalNode
		finalJoin := network.BetaNodes["cascade_rule_join_2"]
		if finalJoin == nil {
			t.Fatal("Final JoinNode not found")
		}

		finalJoinTyped := finalJoin.(*JoinNode)
		if len(finalJoinTyped.Children) != 1 || finalJoinTyped.Children[0] != terminalNode {
			t.Error("Final JoinNode not connected to TerminalNode")
		}
	})

	t.Run("create 4-variable cascade", func(t *testing.T) {
		network := NewReteNetwork(storage)

		// Create TypeNodes
		types := []string{"T1", "T2", "T3", "T4"}
		for _, typeName := range types {
			typeNode := NewTypeNode(typeName, TypeDefinition{Name: typeName}, storage)
			network.TypeNodes[typeName] = typeNode
			network.RootNode.AddChild(typeNode)
		}

		variableNames := []string{"v1", "v2", "v3", "v4"}
		variableTypes := []string{"T1", "T2", "T3", "T4"}
		condition := map[string]interface{}{}
		terminalNode := utils.CreateTerminalNode(network, "four_var_rule", &Action{Type: "print"})

		err := jrb.createCascadeJoinRuleLegacy(
			network,
			"four_var_rule",
			variableNames,
			variableTypes,
			condition,
			terminalNode,
		)

		if err != nil {
			t.Fatalf("createCascadeJoinRuleLegacy failed: %v", err)
		}

		// Should create 3 JoinNodes for 4 variables
		// four_var_rule_join_0_1 (v1 ⋈ v2)
		// four_var_rule_join_2 ((v1 ⋈ v2) ⋈ v3)
		// four_var_rule_join_3 (((v1 ⋈ v2) ⋈ v3) ⋈ v4)
		if len(network.BetaNodes) != 3 {
			t.Errorf("Expected 3 JoinNodes for 4 variables, got %d", len(network.BetaNodes))
		}
	})
}

func TestJoinRuleBuilder_buildJoinPatterns(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	jrb := NewJoinRuleBuilder(utils)

	t.Run("build patterns for 2 variables", func(t *testing.T) {
		variableNames := []string{"p", "e"}
		variableTypes := []string{"Person", "Employee"}
		condition := map[string]interface{}{"type": "comparison"}

		patterns := jrb.buildJoinPatterns(variableNames, variableTypes, condition)

		if len(patterns) != 1 {
			t.Fatalf("Expected 1 pattern for 2 variables, got %d", len(patterns))
		}

		pattern := patterns[0]
		if len(pattern.LeftVars) != 1 || pattern.LeftVars[0] != "p" {
			t.Errorf("Pattern.LeftVars = %v, want ['p']", pattern.LeftVars)
		}

		if len(pattern.RightVars) != 1 || pattern.RightVars[0] != "e" {
			t.Errorf("Pattern.RightVars = %v, want ['e']", pattern.RightVars)
		}

		if len(pattern.AllVars) != 2 {
			t.Errorf("Pattern.AllVars length = %d, want 2", len(pattern.AllVars))
		}
	})

	t.Run("build patterns for 3 variables", func(t *testing.T) {
		variableNames := []string{"p", "e", "d"}
		variableTypes := []string{"Person", "Employee", "Department"}
		condition := map[string]interface{}{}

		patterns := jrb.buildJoinPatterns(variableNames, variableTypes, condition)

		if len(patterns) != 2 {
			t.Fatalf("Expected 2 patterns for 3 variables, got %d", len(patterns))
		}

		// First pattern: p ⋈ e
		if len(patterns[0].LeftVars) != 1 || patterns[0].LeftVars[0] != "p" {
			t.Errorf("Pattern[0].LeftVars = %v, want ['p']", patterns[0].LeftVars)
		}
		if len(patterns[0].RightVars) != 1 || patterns[0].RightVars[0] != "e" {
			t.Errorf("Pattern[0].RightVars = %v, want ['e']", patterns[0].RightVars)
		}

		// Second pattern: (p,e) ⋈ d
		if len(patterns[1].LeftVars) != 2 {
			t.Errorf("Pattern[1].LeftVars length = %d, want 2", len(patterns[1].LeftVars))
		}
		if len(patterns[1].RightVars) != 1 || patterns[1].RightVars[0] != "d" {
			t.Errorf("Pattern[1].RightVars = %v, want ['d']", patterns[1].RightVars)
		}
		if len(patterns[1].AllVars) != 3 {
			t.Errorf("Pattern[1].AllVars length = %d, want 3", len(patterns[1].AllVars))
		}
	})

	t.Run("build patterns for 4 variables", func(t *testing.T) {
		variableNames := []string{"v1", "v2", "v3", "v4"}
		variableTypes := []string{"T1", "T2", "T3", "T4"}
		condition := map[string]interface{}{}

		patterns := jrb.buildJoinPatterns(variableNames, variableTypes, condition)

		if len(patterns) != 3 {
			t.Fatalf("Expected 3 patterns for 4 variables, got %d", len(patterns))
		}

		// Verify progressive accumulation
		for i, pattern := range patterns {
			if len(pattern.RightVars) != 1 {
				t.Errorf("Pattern[%d].RightVars length = %d, want 1", i, len(pattern.RightVars))
			}

			expectedAllVarsSize := i + 2
			if len(pattern.AllVars) != expectedAllVarsSize {
				t.Errorf("Pattern[%d].AllVars length = %d, want %d", i, len(pattern.AllVars), expectedAllVarsSize)
			}
		}
	})
}

func TestJoinRuleBuilder_CreateJoinRule(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	jrb := NewJoinRuleBuilder(utils)

	t.Run("create 2-variable join rule", func(t *testing.T) {
		network := NewReteNetwork(storage)

		personNode := NewTypeNode("Person", TypeDefinition{Name: "Person"}, storage)
		network.TypeNodes["Person"] = personNode
		network.RootNode.AddChild(personNode)

		employeeNode := NewTypeNode("Employee", TypeDefinition{Name: "Employee"}, storage)
		network.TypeNodes["Employee"] = employeeNode
		network.RootNode.AddChild(employeeNode)

		variableNames := []string{"p", "e"}
		variableTypes := []string{"Person", "Employee"}
		condition := map[string]interface{}{"type": "comparison"}
		action := &Action{Type: "print", Job: &JobCall{Name: "print", Args: []interface{}{"Join found"}}}

		err := jrb.CreateJoinRule(network, "join_rule", variableNames, variableTypes, condition, action)

		if err != nil {
			t.Fatalf("CreateJoinRule failed: %v", err)
		}

		// Should create binary join
		joinNode := network.BetaNodes["join_rule_join"]
		if joinNode == nil {
			t.Fatal("JoinNode not created")
		}
	})

	t.Run("create 3-variable join rule - cascade", func(t *testing.T) {
		network := NewReteNetwork(storage)

		types := []string{"Person", "Employee", "Department"}
		for _, typeName := range types {
			typeNode := NewTypeNode(typeName, TypeDefinition{Name: typeName}, storage)
			network.TypeNodes[typeName] = typeNode
			network.RootNode.AddChild(typeNode)
		}

		variableNames := []string{"p", "e", "d"}
		variableTypes := []string{"Person", "Employee", "Department"}
		condition := map[string]interface{}{}
		action := &Action{Type: "print"}

		err := jrb.CreateJoinRule(network, "cascade_rule", variableNames, variableTypes, condition, action)

		if err != nil {
			t.Fatalf("CreateJoinRule failed: %v", err)
		}

		// Should create cascade (2 JoinNodes for 3 variables)
		if len(network.BetaNodes) < 2 {
			t.Errorf("Expected at least 2 JoinNodes for cascade, got %d", len(network.BetaNodes))
		}
	})

	t.Run("with BetaChainBuilder enabled", func(t *testing.T) {
		network := NewReteNetwork(storage)
		config := BetaSharingConfig{Enabled: true, HashCacheSize: 1000, MaxSharedNodes: 10000}
		lifecycle := NewLifecycleManager()
		network.BetaSharingRegistry = NewBetaSharingRegistry(config, lifecycle)
		network.BetaChainBuilder = NewBetaChainBuilderWithRegistry(network, storage, network.BetaSharingRegistry)

		types := []string{"T1", "T2", "T3"}
		for _, typeName := range types {
			typeNode := NewTypeNode(typeName, TypeDefinition{Name: typeName}, storage)
			network.TypeNodes[typeName] = typeNode
			network.RootNode.AddChild(typeNode)
		}

		variableNames := []string{"v1", "v2", "v3"}
		variableTypes := []string{"T1", "T2", "T3"}
		condition := map[string]interface{}{}
		action := &Action{Type: "print"}

		err := jrb.CreateJoinRule(network, "builder_rule", variableNames, variableTypes, condition, action)

		if err != nil {
			t.Fatalf("CreateJoinRule with BetaChainBuilder failed: %v", err)
		}

		// Should use BetaChainBuilder for cascade
		if len(network.BetaNodes) == 0 {
			t.Error("No BetaNodes created with BetaChainBuilder")
		}
	})
}

func TestJoinRuleBuilder_Integration(t *testing.T) {
	// Integration test: Create multiple join rules and verify network structure
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	jrb := NewJoinRuleBuilder(utils)
	network := NewReteNetwork(storage)

	// Setup: Create TypeNodes
	types := map[string][]Field{
		"Person": {
			{Name: "id", Type: "number"},
			{Name: "name", Type: "string"},
		},
		"Employee": {
			{Name: "person_id", Type: "number"},
			{Name: "department_id", Type: "number"},
		},
		"Department": {
			{Name: "id", Type: "number"},
			{Name: "name", Type: "string"},
		},
	}

	for typeName, fields := range types {
		typeNode := NewTypeNode(typeName, TypeDefinition{
			Type:   "type",
			Name:   typeName,
			Fields: fields,
		}, storage)
		network.TypeNodes[typeName] = typeNode
		network.RootNode.AddChild(typeNode)
	}

	// Rule 1: Simple 2-variable join
	err := jrb.CreateJoinRule(
		network,
		"person_employee_join",
		[]string{"p", "e"},
		[]string{"Person", "Employee"},
		map[string]interface{}{
			"type":     "comparison",
			"operator": "==",
		},
		&Action{Type: "print", Job: &JobCall{Name: "print", Args: []interface{}{"Person-Employee match"}}},
	)
	if err != nil {
		t.Fatalf("Failed to create person_employee_join: %v", err)
	}

	// Rule 2: 3-variable cascade join
	err = jrb.CreateJoinRule(
		network,
		"person_employee_dept_join",
		[]string{"p", "e", "d"},
		[]string{"Person", "Employee", "Department"},
		map[string]interface{}{
			"type": "comparison",
		},
		&Action{Type: "print", Job: &JobCall{Name: "print", Args: []interface{}{"Full hierarchy match"}}},
	)
	if err != nil {
		t.Fatalf("Failed to create person_employee_dept_join: %v", err)
	}

	// Verify network structure
	if len(network.BetaNodes) < 2 {
		t.Errorf("Expected at least 2 BetaNodes, got %d", len(network.BetaNodes))
	}

	// Verify all JoinNodes have TerminalNodes
	terminalCount := 0
	for _, betaNode := range network.BetaNodes {
		joinNode, ok := betaNode.(*JoinNode)
		if !ok {
			continue
		}

		// Walk the tree to find terminal nodes
		for _, child := range joinNode.Children {
			if _, isTerminal := child.(*TerminalNode); isTerminal {
				terminalCount++
			} else if childJoin, isJoin := child.(*JoinNode); isJoin {
				// Check child join's children
				for _, grandchild := range childJoin.Children {
					if _, isTerminal := grandchild.(*TerminalNode); isTerminal {
						terminalCount++
					}
				}
			}
		}
	}

	if terminalCount < 2 {
		t.Errorf("Expected at least 2 TerminalNodes in the network, found %d", terminalCount)
	}
}
