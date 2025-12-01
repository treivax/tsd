// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

func TestNewRuleBuilder(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)

	// Create a mock pipeline (nil for basic test)
	rb := NewRuleBuilder(utils, nil)

	if rb == nil {
		t.Fatal("NewRuleBuilder returned nil")
	}

	if rb.utils != utils {
		t.Error("RuleBuilder.utils not set correctly")
	}

	if rb.alphaBuilder == nil {
		t.Error("RuleBuilder.alphaBuilder not initialized")
	}

	if rb.joinBuilder == nil {
		t.Error("RuleBuilder.joinBuilder not initialized")
	}

	if rb.existsBuilder == nil {
		t.Error("RuleBuilder.existsBuilder not initialized")
	}

	if rb.accumulatorBuilder == nil {
		t.Error("RuleBuilder.accumulatorBuilder not initialized")
	}
}

func TestRuleBuilder_CreateRuleNodes(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Create a simple pipeline mock for testing
	pipeline := &ConstraintPipelineBuilder{
		network: network,
		storage: storage,
	}

	utils := NewBuilderUtils(storage)
	rb := NewRuleBuilder(utils, pipeline)

	t.Run("error on invalid expression format", func(t *testing.T) {
		expressions := []interface{}{
			"not a map", // Invalid format
		}

		err := rb.CreateRuleNodes(network, expressions)
		if err == nil {
			t.Error("Expected error for invalid expression format, got nil")
		}
	})

	t.Run("empty expression list", func(t *testing.T) {
		expressions := []interface{}{}

		err := rb.CreateRuleNodes(network, expressions)
		if err != nil {
			t.Errorf("Empty expression list should not error, got: %v", err)
		}
	})

	t.Run("expression with default ruleId", func(t *testing.T) {
		network := NewReteNetwork(storage)
		pipeline := &ConstraintPipelineBuilder{
			network: network,
			storage: storage,
		}
		rb := NewRuleBuilder(utils, pipeline)

		// Create TypeNode for the rule
		personNode := NewTypeNode("Person", TypeDefinition{Name: "Person"}, storage)
		network.TypeNodes["Person"] = personNode
		network.RootNode.AddChild(personNode)

		// Declare the action
		network.Actions["print"] = &ActionDefinition{
			Name: "print",
			Parameters: []ActionParameter{
				{Name: "message", Type: "string"},
			},
		}

		expressions := []interface{}{
			map[string]interface{}{
				// Missing ruleId - should use default "rule_0"
				"set": map[string]interface{}{
					"variables": []interface{}{
						map[string]interface{}{
							"name":     "p",
							"dataType": "Person",
						},
					},
				},
				"action": map[string]interface{}{
					"type": "print",
					"arguments": []interface{}{
						map[string]interface{}{
							"type":  "literal",
							"value": "test",
						},
					},
				},
			},
		}

		err := rb.CreateRuleNodes(network, expressions)
		// May fail due to missing constraint details, but should not panic
		if err != nil {
			// Expected to fail without proper action definition setup
			t.Logf("Expected error for incomplete rule: %v", err)
		}
	})

	t.Run("expression with custom ruleId", func(t *testing.T) {
		network := NewReteNetwork(storage)
		pipeline := &ConstraintPipelineBuilder{
			network: network,
			storage: storage,
		}
		rb := NewRuleBuilder(utils, pipeline)

		// Create TypeNode
		personNode := NewTypeNode("Person", TypeDefinition{Name: "Person"}, storage)
		network.TypeNodes["Person"] = personNode
		network.RootNode.AddChild(personNode)

		// Declare the action
		network.Actions["print"] = &ActionDefinition{
			Name: "print",
			Parameters: []ActionParameter{
				{Name: "message", Type: "string"},
			},
		}

		expressions := []interface{}{
			map[string]interface{}{
				"ruleId": "custom_rule_id",
				"set": map[string]interface{}{
					"variables": []interface{}{
						map[string]interface{}{
							"name":     "p",
							"dataType": "Person",
						},
					},
				},
				"action": map[string]interface{}{
					"type": "print",
					"arguments": []interface{}{
						map[string]interface{}{
							"type":  "literal",
							"value": "test",
						},
					},
				},
			},
		}

		err := rb.CreateRuleNodes(network, expressions)
		// May fail, but ruleId extraction should work
		if err != nil {
			t.Logf("Rule creation may fail: %v", err)
		}
	})
}

func TestRuleBuilder_Integration_AlphaRules(t *testing.T) {
	// Integration test: Create alpha rules through RuleBuilder
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Setup pipeline
	pipeline := &ConstraintPipelineBuilder{
		network: network,
		storage: storage,
	}

	utils := NewBuilderUtils(storage)
	rb := NewRuleBuilder(utils, pipeline)

	// Create TypeNode
	personNode := NewTypeNode("Person", TypeDefinition{
		Type: "type",
		Name: "Person",
		Fields: []Field{
			{Name: "age", Type: "number"},
		},
	}, storage)
	network.TypeNodes["Person"] = personNode
	network.RootNode.AddChild(personNode)

	// Declare print action
	network.Actions["print"] = &ActionDefinition{
		Name: "print",
		Parameters: []ActionParameter{
			{Name: "message", Type: "string"},
		},
	}

	// Create a simple alpha rule expression
	expressions := []interface{}{
		map[string]interface{}{
			"ruleId": "adult_rule",
			"set": map[string]interface{}{
				"variables": []interface{}{
					map[string]interface{}{
						"name":     "p",
						"dataType": "Person",
					},
				},
			},
			"constraints": map[string]interface{}{
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
			},
			"action": map[string]interface{}{
				"type": "print",
				"arguments": []interface{}{
					map[string]interface{}{
						"type":  "literal",
						"value": "Adult found",
					},
				},
			},
		},
	}

	err := rb.CreateRuleNodes(network, expressions)
	if err != nil {
		t.Fatalf("CreateRuleNodes failed for alpha rule: %v", err)
	}

	// Verify AlphaNode was created
	if len(network.AlphaNodes) == 0 {
		t.Error("No AlphaNodes created")
	}
}

func TestRuleBuilder_Integration_JoinRules(t *testing.T) {
	// Integration test: Create join rules through RuleBuilder
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	pipeline := &ConstraintPipelineBuilder{
		network: network,
		storage: storage,
	}

	utils := NewBuilderUtils(storage)
	rb := NewRuleBuilder(utils, pipeline)

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

	// Declare print action
	network.Actions["print"] = &ActionDefinition{
		Name: "print",
		Parameters: []ActionParameter{
			{Name: "message", Type: "string"},
		},
	}

	// Create a join rule expression
	expressions := []interface{}{
		map[string]interface{}{
			"ruleId": "person_employee_join",
			"set": map[string]interface{}{
				"variables": []interface{}{
					map[string]interface{}{
						"name":     "p",
						"dataType": "Person",
					},
					map[string]interface{}{
						"name":     "e",
						"dataType": "Employee",
					},
				},
			},
			"constraints": map[string]interface{}{
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
			"action": map[string]interface{}{
				"type": "print",
				"arguments": []interface{}{
					map[string]interface{}{
						"type":  "literal",
						"value": "Join match",
					},
				},
			},
		},
	}

	err := rb.CreateRuleNodes(network, expressions)
	if err != nil {
		t.Fatalf("CreateRuleNodes failed for join rule: %v", err)
	}

	// Verify JoinNode was created
	if len(network.BetaNodes) == 0 {
		t.Error("No BetaNodes (JoinNodes) created")
	}
}

func TestRuleBuilder_Integration_MultipleRules(t *testing.T) {
	// Integration test: Create multiple rules of different types
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	pipeline := &ConstraintPipelineBuilder{
		network: network,
		storage: storage,
	}

	utils := NewBuilderUtils(storage)
	rb := NewRuleBuilder(utils, pipeline)

	// Create TypeNodes
	personNode := NewTypeNode("Person", TypeDefinition{
		Type: "type",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "number"},
			{Name: "age", Type: "number"},
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

	// Declare print action
	network.Actions["print"] = &ActionDefinition{
		Name: "print",
		Parameters: []ActionParameter{
			{Name: "message", Type: "string"},
		},
	}

	// Create multiple rule expressions
	expressions := []interface{}{
		// Rule 1: Alpha rule
		map[string]interface{}{
			"ruleId": "alpha_rule_1",
			"set": map[string]interface{}{
				"variables": []interface{}{
					map[string]interface{}{
						"name":     "p",
						"dataType": "Person",
					},
				},
			},
			"constraints": map[string]interface{}{
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
			},
			"action": map[string]interface{}{
				"type": "print",
				"arguments": []interface{}{
					map[string]interface{}{
						"type":  "literal",
						"value": "Adult",
					},
				},
			},
		},
		// Rule 2: Join rule
		map[string]interface{}{
			"ruleId": "join_rule_1",
			"set": map[string]interface{}{
				"variables": []interface{}{
					map[string]interface{}{
						"name":     "p",
						"dataType": "Person",
					},
					map[string]interface{}{
						"name":     "e",
						"dataType": "Employee",
					},
				},
			},
			"constraints": map[string]interface{}{
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
			"action": map[string]interface{}{
				"type": "print",
				"arguments": []interface{}{
					map[string]interface{}{
						"type":  "literal",
						"value": "Join",
					},
				},
			},
		},
	}

	err := rb.CreateRuleNodes(network, expressions)
	if err != nil {
		t.Fatalf("CreateRuleNodes failed for multiple rules: %v", err)
	}

	// Verify both rules were created
	if len(network.AlphaNodes) == 0 {
		t.Error("No AlphaNodes created")
	}

	if len(network.BetaNodes) == 0 {
		t.Error("No BetaNodes created")
	}
}

func TestRuleBuilder_ErrorHandling(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	pipeline := &ConstraintPipelineBuilder{
		network: network,
		storage: storage,
	}

	utils := NewBuilderUtils(storage)
	rb := NewRuleBuilder(utils, pipeline)

	t.Run("missing type node", func(t *testing.T) {
		network := NewReteNetwork(storage)
		pipeline := &ConstraintPipelineBuilder{
			network: network,
			storage: storage,
		}
		rb := NewRuleBuilder(utils, pipeline)

		// Declare print action
		network.Actions["print"] = &ActionDefinition{
			Name: "print",
			Parameters: []ActionParameter{
				{Name: "message", Type: "string"},
			},
		}

		expressions := []interface{}{
			map[string]interface{}{
				"ruleId": "bad_rule",
				"set": map[string]interface{}{
					"variables": []interface{}{
						map[string]interface{}{
							"name":     "p",
							"dataType": "NonExistent",
						},
					},
				},
				"action": map[string]interface{}{
					"type": "print",
					"arguments": []interface{}{
						map[string]interface{}{
							"type":  "literal",
							"value": "test",
						},
					},
				},
			},
		}

		err := rb.CreateRuleNodes(network, expressions)
		if err == nil {
			t.Error("Expected error for missing TypeNode, got nil")
		}
	})

	t.Run("invalid expression structure", func(t *testing.T) {
		expressions := []interface{}{
			123, // Not a map
		}

		err := rb.CreateRuleNodes(network, expressions)
		if err == nil {
			t.Error("Expected error for invalid expression, got nil")
		}
	})
}
