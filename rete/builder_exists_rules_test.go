// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

func TestNewExistsRuleBuilder(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)

	erb := NewExistsRuleBuilder(utils)

	if erb == nil {
		t.Fatal("NewExistsRuleBuilder returned nil")
	}

	if erb.utils != utils {
		t.Error("ExistsRuleBuilder.utils not set correctly")
	}
}

func TestExistsRuleBuilder_ExtractExistsVariables(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	erb := NewExistsRuleBuilder(utils)

	tests := []struct {
		name           string
		exprMap        map[string]interface{}
		wantMainVar    string
		wantExistsVar  string
		wantMainType   string
		wantExistsType string
		wantErr        bool
	}{
		{
			name: "valid EXISTS expression",
			exprMap: map[string]interface{}{
				"set": map[string]interface{}{
					"variables": []interface{}{
						map[string]interface{}{
							"name":     "p",
							"dataType": "Person",
						},
					},
				},
				"constraints": map[string]interface{}{
					"variable": map[string]interface{}{
						"name":     "e",
						"dataType": "Employee",
					},
				},
			},
			wantMainVar:    "p",
			wantExistsVar:  "e",
			wantMainType:   "Person",
			wantExistsType: "Employee",
			wantErr:        false,
		},
		{
			name: "missing main variable",
			exprMap: map[string]interface{}{
				"constraints": map[string]interface{}{
					"variable": map[string]interface{}{
						"name":     "e",
						"dataType": "Employee",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing exists variable",
			exprMap: map[string]interface{}{
				"set": map[string]interface{}{
					"variables": []interface{}{
						map[string]interface{}{
							"name":     "p",
							"dataType": "Person",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty variable name",
			exprMap: map[string]interface{}{
				"set": map[string]interface{}{
					"variables": []interface{}{
						map[string]interface{}{
							"name":     "",
							"dataType": "Person",
						},
					},
				},
				"constraints": map[string]interface{}{
					"variable": map[string]interface{}{
						"name":     "e",
						"dataType": "Employee",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "variables without types",
			exprMap: map[string]interface{}{
				"set": map[string]interface{}{
					"variables": []interface{}{
						map[string]interface{}{
							"name": "p",
						},
					},
				},
				"constraints": map[string]interface{}{
					"variable": map[string]interface{}{
						"name": "e",
					},
				},
			},
			wantMainVar:    "p",
			wantExistsVar:  "e",
			wantMainType:   "",
			wantExistsType: "",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mainVar, existsVar, mainType, existsType, err := erb.ExtractExistsVariables(tt.exprMap)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractExistsVariables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if mainVar != tt.wantMainVar {
					t.Errorf("mainVar = %q, want %q", mainVar, tt.wantMainVar)
				}
				if existsVar != tt.wantExistsVar {
					t.Errorf("existsVar = %q, want %q", existsVar, tt.wantExistsVar)
				}
				if mainType != tt.wantMainType {
					t.Errorf("mainType = %q, want %q", mainType, tt.wantMainType)
				}
				if existsType != tt.wantExistsType {
					t.Errorf("existsType = %q, want %q", existsType, tt.wantExistsType)
				}
			}
		})
	}
}

func TestExistsRuleBuilder_ExtractExistsConditions(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	erb := NewExistsRuleBuilder(utils)

	tests := []struct {
		name      string
		exprMap   map[string]interface{}
		wantCount int
		wantErr   bool
	}{
		{
			name: "single condition",
			exprMap: map[string]interface{}{
				"constraints": map[string]interface{}{
					"condition": map[string]interface{}{
						"type":     "comparison",
						"operator": "==",
					},
				},
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "multiple conditions",
			exprMap: map[string]interface{}{
				"constraints": map[string]interface{}{
					"conditions": []interface{}{
						map[string]interface{}{
							"type":     "comparison",
							"operator": "==",
						},
						map[string]interface{}{
							"type":     "comparison",
							"operator": ">",
						},
					},
				},
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "no conditions",
			exprMap:   map[string]interface{}{},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "empty constraints",
			exprMap: map[string]interface{}{
				"constraints": map[string]interface{}{},
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conditions, err := erb.ExtractExistsConditions(tt.exprMap)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractExistsConditions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(conditions) != tt.wantCount {
				t.Errorf("len(conditions) = %d, want %d", len(conditions), tt.wantCount)
			}
		})
	}
}

func TestExistsRuleBuilder_ConnectExistsNodeToTypeNodes(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	erb := NewExistsRuleBuilder(utils)

	t.Run("connect both variables", func(t *testing.T) {
		network := NewReteNetwork(storage)

		// Create TypeNodes
		personNode := NewTypeNode("Person", TypeDefinition{Name: "Person"}, storage)
		network.TypeNodes["Person"] = personNode
		network.RootNode.AddChild(personNode)

		employeeNode := NewTypeNode("Employee", TypeDefinition{Name: "Employee"}, storage)
		network.TypeNodes["Employee"] = employeeNode
		network.RootNode.AddChild(employeeNode)

		// Create ExistsNode
		varTypes := map[string]string{
			"p": "Person",
			"e": "Employee",
		}
		existsNode := NewExistsNode("test_exists", map[string]interface{}{}, "p", "e", varTypes, storage)

		// Connect
		erb.ConnectExistsNodeToTypeNodes(network, "test_rule", existsNode, "p", "Person", "e", "Employee")

		// Verify connections were made via pass-through alpha nodes
		// TypeNode should have pass-through alpha as child
		if len(personNode.Children) != 1 {
			t.Errorf("Person TypeNode should have 1 child (pass-through alpha), got %d", len(personNode.Children))
		}

		if len(employeeNode.Children) != 1 {
			t.Errorf("Employee TypeNode should have 1 child (pass-through alpha), got %d", len(employeeNode.Children))
		}
	})

	t.Run("connect with empty types", func(t *testing.T) {
		network := NewReteNetwork(storage)

		varTypes := map[string]string{}
		existsNode := NewExistsNode("test_exists", map[string]interface{}{}, "p", "e", varTypes, storage)

		// Should not panic with empty types
		erb.ConnectExistsNodeToTypeNodes(network, "test_rule", existsNode, "p", "", "e", "")
	})
}

func TestExistsRuleBuilder_CreateExistsRule(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	erb := NewExistsRuleBuilder(utils)

	t.Run("create complete EXISTS rule", func(t *testing.T) {
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

		exprMap := map[string]interface{}{
			"set": map[string]interface{}{
				"variables": []interface{}{
					map[string]interface{}{
						"name":     "p",
						"dataType": "Person",
					},
				},
			},
			"constraints": map[string]interface{}{
				"variable": map[string]interface{}{
					"name":     "e",
					"dataType": "Employee",
				},
				"condition": map[string]interface{}{
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
				},
			},
		}

		condition := map[string]interface{}{
			"type": "exists",
		}

		action := &Action{
			Type: "print",
			Job:  &JobCall{Name: "print", Args: []interface{}{"Person with employee record exists"}},
		}

		err := erb.CreateExistsRule(network, "exists_rule", exprMap, condition, action)

		if err != nil {
			t.Fatalf("CreateExistsRule failed: %v", err)
		}

		// Verify ExistsNode was created
		existsNode, exists := network.BetaNodes["exists_rule_exists"]
		if !exists {
			t.Fatal("ExistsNode not created")
		}

		existsNodeTyped, ok := existsNode.(*ExistsNode)
		if !ok {
			t.Fatal("BetaNode is not an ExistsNode")
		}

		if existsNodeTyped.MainVariable != "p" {
			t.Errorf("MainVariable = %q, want 'p'", existsNodeTyped.MainVariable)
		}

		if existsNodeTyped.ExistsVariable != "e" {
			t.Errorf("ExistsVariable = %q, want 'e'", existsNodeTyped.ExistsVariable)
		}

		// Verify TerminalNode connection
		if len(existsNodeTyped.Children) != 1 {
			t.Fatalf("ExistsNode should have 1 child (TerminalNode), got %d", len(existsNodeTyped.Children))
		}

		_, isTerminal := existsNodeTyped.Children[0].(*TerminalNode)
		if !isTerminal {
			t.Error("ExistsNode child should be TerminalNode")
		}

		// Verify connections to TypeNodes
		if len(personNode.Children) == 0 {
			t.Error("Person TypeNode should have children (pass-through alpha)")
		}

		if len(employeeNode.Children) == 0 {
			t.Error("Employee TypeNode should have children (pass-through alpha)")
		}
	})

	t.Run("error on missing variables", func(t *testing.T) {
		network := NewReteNetwork(storage)

		exprMap := map[string]interface{}{
			// Missing variable definitions
		}

		err := erb.CreateExistsRule(network, "bad_rule", exprMap, map[string]interface{}{}, &Action{})

		if err == nil {
			t.Error("Expected error for missing variables, got nil")
		}
	})

	t.Run("with multiple conditions", func(t *testing.T) {
		network := NewReteNetwork(storage)

		personNode := NewTypeNode("Person", TypeDefinition{Name: "Person"}, storage)
		network.TypeNodes["Person"] = personNode
		network.RootNode.AddChild(personNode)

		employeeNode := NewTypeNode("Employee", TypeDefinition{Name: "Employee"}, storage)
		network.TypeNodes["Employee"] = employeeNode
		network.RootNode.AddChild(employeeNode)

		exprMap := map[string]interface{}{
			"set": map[string]interface{}{
				"variables": []interface{}{
					map[string]interface{}{
						"name":     "p",
						"dataType": "Person",
					},
				},
			},
			"constraints": map[string]interface{}{
				"variable": map[string]interface{}{
					"name":     "e",
					"dataType": "Employee",
				},
				"conditions": []interface{}{
					map[string]interface{}{
						"type":     "comparison",
						"operator": "==",
					},
					map[string]interface{}{
						"type":     "comparison",
						"operator": ">",
					},
				},
			},
		}

		err := erb.CreateExistsRule(network, "multi_cond_rule", exprMap, map[string]interface{}{}, &Action{Type: "print"})

		if err != nil {
			t.Fatalf("CreateExistsRule failed: %v", err)
		}

		existsNode := network.BetaNodes["multi_cond_rule_exists"]
		if existsNode == nil {
			t.Fatal("ExistsNode not created for multi-condition rule")
		}
	})
}

func TestExistsRuleBuilder_Integration(t *testing.T) {
	// Integration test: Create multiple EXISTS rules and verify the network structure
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	erb := NewExistsRuleBuilder(utils)
	network := NewReteNetwork(storage)

	// Setup: Create TypeNodes
	personNode := NewTypeNode("Person", TypeDefinition{
		Type: "type",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "number"},
			{Name: "name", Type: "string"},
		},
	}, storage)
	network.TypeNodes["Person"] = personNode
	network.RootNode.AddChild(personNode)

	employeeNode := NewTypeNode("Employee", TypeDefinition{
		Type: "type",
		Name: "Employee",
		Fields: []Field{
			{Name: "person_id", Type: "number"},
			{Name: "department", Type: "string"},
		},
	}, storage)
	network.TypeNodes["Employee"] = employeeNode
	network.RootNode.AddChild(employeeNode)

	managerNode := NewTypeNode("Manager", TypeDefinition{
		Type: "type",
		Name: "Manager",
		Fields: []Field{
			{Name: "employee_id", Type: "number"},
		},
	}, storage)
	network.TypeNodes["Manager"] = managerNode
	network.RootNode.AddChild(managerNode)

	// Rule 1: Person EXISTS Employee
	exprMap1 := map[string]interface{}{
		"set": map[string]interface{}{
			"variables": []interface{}{
				map[string]interface{}{"name": "p", "dataType": "Person"},
			},
		},
		"constraints": map[string]interface{}{
			"variable":  map[string]interface{}{"name": "e", "dataType": "Employee"},
			"condition": map[string]interface{}{"type": "comparison", "operator": "=="},
		},
	}

	err := erb.CreateExistsRule(network, "person_has_employee", exprMap1, map[string]interface{}{}, &Action{Type: "print"})
	if err != nil {
		t.Fatalf("Failed to create person_has_employee rule: %v", err)
	}

	// Rule 2: Employee EXISTS Manager
	exprMap2 := map[string]interface{}{
		"set": map[string]interface{}{
			"variables": []interface{}{
				map[string]interface{}{"name": "e", "dataType": "Employee"},
			},
		},
		"constraints": map[string]interface{}{
			"variable":  map[string]interface{}{"name": "m", "dataType": "Manager"},
			"condition": map[string]interface{}{"type": "comparison", "operator": "=="},
		},
	}

	err = erb.CreateExistsRule(network, "employee_has_manager", exprMap2, map[string]interface{}{}, &Action{Type: "print"})
	if err != nil {
		t.Fatalf("Failed to create employee_has_manager rule: %v", err)
	}

	// Verify both ExistsNodes created
	if len(network.BetaNodes) != 2 {
		t.Errorf("Expected 2 BetaNodes (ExistsNodes), got %d", len(network.BetaNodes))
	}

	// Verify specific ExistsNodes exist
	if _, exists := network.BetaNodes["person_has_employee_exists"]; !exists {
		t.Error("person_has_employee ExistsNode not found")
	}

	if _, exists := network.BetaNodes["employee_has_manager_exists"]; !exists {
		t.Error("employee_has_manager ExistsNode not found")
	}

	// Verify each ExistsNode has a TerminalNode
	for nodeID, betaNode := range network.BetaNodes {
		existsNode, ok := betaNode.(*ExistsNode)
		if !ok {
			t.Errorf("BetaNode %s is not an ExistsNode", nodeID)
			continue
		}

		if len(existsNode.Children) != 1 {
			t.Errorf("ExistsNode %s should have 1 child, got %d", nodeID, len(existsNode.Children))
		}

		_, isTerminal := existsNode.Children[0].(*TerminalNode)
		if !isTerminal {
			t.Errorf("ExistsNode %s child is not a TerminalNode", nodeID)
		}
	}
}
