// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// BenchmarkBuilderUtils_CreateTerminalNode benchmarks terminal node creation
func BenchmarkBuilderUtils_CreateTerminalNode(b *testing.B) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	network := NewReteNetwork(storage)
	action := &Action{Type: "print", Job: &JobCall{Name: "print", Args: []interface{}{"test"}}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.CreateTerminalNode(network, "test_rule", action)
	}
}

// BenchmarkBuilderUtils_ConnectTypeNodeToBetaNode benchmarks type node connections
func BenchmarkBuilderUtils_ConnectTypeNodeToBetaNode(b *testing.B) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	network := NewReteNetwork(storage)

	// Setup
	typeNode := NewTypeNode("Person", TypeDefinition{Name: "Person"}, storage)
	network.TypeNodes["Person"] = typeNode
	network.RootNode.AddChild(typeNode)

	betaNode := NewJoinNode("test_join", map[string]interface{}{}, []string{"p"}, []string{"e"}, map[string]string{"p": "Person"}, storage)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.ConnectTypeNodeToBetaNode(network, "test_rule", "p", "Person", betaNode, NodeSideLeft)
	}
}

// BenchmarkTypeBuilder_CreateTypeDefinition benchmarks type definition creation
func BenchmarkTypeBuilder_CreateTypeDefinition(b *testing.B) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	tb := NewTypeBuilder(utils)

	typeMap := map[string]interface{}{
		"name": "Employee",
		"fields": []interface{}{
			map[string]interface{}{"name": "id", "type": "number"},
			map[string]interface{}{"name": "name", "type": "string"},
			map[string]interface{}{"name": "salary", "type": "number"},
			map[string]interface{}{"name": "department", "type": "string"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tb.CreateTypeDefinition("Employee", typeMap)
	}
}

// BenchmarkTypeBuilder_CreateTypeNodes benchmarks creating multiple type nodes
func BenchmarkTypeBuilder_CreateTypeNodes(b *testing.B) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	tb := NewTypeBuilder(utils)

	types := []interface{}{
		map[string]interface{}{
			"name": "Person",
			"fields": []interface{}{
				map[string]interface{}{"name": "id", "type": "number"},
				map[string]interface{}{"name": "name", "type": "string"},
			},
		},
		map[string]interface{}{
			"name": "Employee",
			"fields": []interface{}{
				map[string]interface{}{"name": "person_id", "type": "number"},
				map[string]interface{}{"name": "salary", "type": "number"},
			},
		},
		map[string]interface{}{
			"name": "Department",
			"fields": []interface{}{
				map[string]interface{}{"name": "id", "type": "number"},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		network := NewReteNetwork(storage)
		_ = tb.CreateTypeNodes(network, types, storage)
	}
}

// BenchmarkAlphaRuleBuilder_CreateAlphaRule benchmarks alpha rule creation
func BenchmarkAlphaRuleBuilder_CreateAlphaRule(b *testing.B) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	arb := NewAlphaRuleBuilder(utils)

	variables := []map[string]interface{}{
		{"name": "p", "dataType": "Person"},
	}
	variableNames := []string{"p"}
	variableTypes := []string{"Person"}
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
		"right":    map[string]interface{}{"type": "literal", "value": 18},
	}
	action := &Action{Type: "print", Job: &JobCall{Name: "print", Args: []interface{}{"Adult"}}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		network := NewReteNetwork(storage)
		typeNode := NewTypeNode("Person", TypeDefinition{Name: "Person"}, storage)
		network.TypeNodes["Person"] = typeNode
		network.RootNode.AddChild(typeNode)

		_ = arb.CreateAlphaRule(network, "alpha_rule", variables, variableNames, variableTypes, condition, action)
	}
}

// BenchmarkJoinRuleBuilder_CreateBinaryJoin benchmarks binary join creation
func BenchmarkJoinRuleBuilder_CreateBinaryJoin(b *testing.B) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	jrb := NewJoinRuleBuilder(utils)

	variableNames := []string{"p", "e"}
	variableTypes := []string{"Person", "Employee"}
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
	}
	action := &Action{Type: "print"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		network := NewReteNetwork(storage)
		personNode := NewTypeNode("Person", TypeDefinition{Name: "Person"}, storage)
		network.TypeNodes["Person"] = personNode
		network.RootNode.AddChild(personNode)
		employeeNode := NewTypeNode("Employee", TypeDefinition{Name: "Employee"}, storage)
		network.TypeNodes["Employee"] = employeeNode
		network.RootNode.AddChild(employeeNode)

		_ = jrb.CreateJoinRule(network, "join_rule", variableNames, variableTypes, condition, action)
	}
}

// BenchmarkJoinRuleBuilder_CreateCascadeJoin benchmarks cascade join creation (3 variables)
func BenchmarkJoinRuleBuilder_CreateCascadeJoin(b *testing.B) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	jrb := NewJoinRuleBuilder(utils)

	variableNames := []string{"p", "e", "d"}
	variableTypes := []string{"Person", "Employee", "Department"}
	condition := map[string]interface{}{"type": "comparison"}
	action := &Action{Type: "print"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		network := NewReteNetwork(storage)
		types := []string{"Person", "Employee", "Department"}
		for _, typeName := range types {
			typeNode := NewTypeNode(typeName, TypeDefinition{Name: typeName}, storage)
			network.TypeNodes[typeName] = typeNode
			network.RootNode.AddChild(typeNode)
		}

		_ = jrb.CreateJoinRule(network, "cascade_rule", variableNames, variableTypes, condition, action)
	}
}

// BenchmarkJoinRuleBuilder_CreateCascadeJoin4Vars benchmarks cascade join with 4 variables
func BenchmarkJoinRuleBuilder_CreateCascadeJoin4Vars(b *testing.B) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	jrb := NewJoinRuleBuilder(utils)

	variableNames := []string{"v1", "v2", "v3", "v4"}
	variableTypes := []string{"T1", "T2", "T3", "T4"}
	condition := map[string]interface{}{}
	action := &Action{Type: "print"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		network := NewReteNetwork(storage)
		for _, typeName := range variableTypes {
			typeNode := NewTypeNode(typeName, TypeDefinition{Name: typeName}, storage)
			network.TypeNodes[typeName] = typeNode
			network.RootNode.AddChild(typeNode)
		}

		_ = jrb.CreateJoinRule(network, "cascade_4_rule", variableNames, variableTypes, condition, action)
	}
}

// BenchmarkExistsRuleBuilder_CreateExistsRule benchmarks EXISTS rule creation
func BenchmarkExistsRuleBuilder_CreateExistsRule(b *testing.B) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	erb := NewExistsRuleBuilder(utils)

	exprMap := map[string]interface{}{
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
	condition := map[string]interface{}{"type": "exists"}
	action := &Action{Type: "print"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		network := NewReteNetwork(storage)
		personNode := NewTypeNode("Person", TypeDefinition{Name: "Person"}, storage)
		network.TypeNodes["Person"] = personNode
		network.RootNode.AddChild(personNode)
		employeeNode := NewTypeNode("Employee", TypeDefinition{Name: "Employee"}, storage)
		network.TypeNodes["Employee"] = employeeNode
		network.RootNode.AddChild(employeeNode)

		_ = erb.CreateExistsRule(network, "exists_rule", exprMap, condition, action)
	}
}

// BenchmarkAccumulatorRuleBuilder_CreateAccumulatorRule benchmarks accumulator rule creation
func BenchmarkAccumulatorRuleBuilder_CreateAccumulatorRule(b *testing.B) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	arb := NewAccumulatorRuleBuilder(utils)

	variables := []map[string]interface{}{
		{"name": "e", "dataType": "Employee"},
		{"name": "p", "dataType": "Performance"},
	}
	variableNames := []string{"e", "p"}
	variableTypes := []string{"Employee", "Performance"}
	aggInfo := &AggregationInfo{
		MainVariable: "e",
		MainType:     "Employee",
		AggVariable:  "p",
		AggType:      "Performance",
		Field:        "score",
		JoinField:    "employee_id",
		MainField:    "id",
		Function:     "AVG",
		Operator:     ">",
		Threshold:    75.0,
	}
	action := &Action{Type: "print"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		network := NewReteNetwork(storage)
		employeeNode := NewTypeNode("Employee", TypeDefinition{Name: "Employee"}, storage)
		network.TypeNodes["Employee"] = employeeNode
		network.RootNode.AddChild(employeeNode)
		performanceNode := NewTypeNode("Performance", TypeDefinition{Name: "Performance"}, storage)
		network.TypeNodes["Performance"] = performanceNode
		network.RootNode.AddChild(performanceNode)

		_ = arb.CreateAccumulatorRule(network, "accum_rule", variables, variableNames, variableTypes, aggInfo, action)
	}
}

// BenchmarkAccumulatorRuleBuilder_CreateMultiSourceRule benchmarks multi-source accumulator
func BenchmarkAccumulatorRuleBuilder_CreateMultiSourceRule(b *testing.B) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	arb := NewAccumulatorRuleBuilder(utils)

	aggInfo := &AggregationInfo{
		MainVariable: "e",
		MainType:     "Employee",
		SourcePatterns: []SourcePattern{
			{Variable: "s", Type: "Sales"},
			{Variable: "p", Type: "Performance"},
		},
		JoinConditions: []JoinCondition{
			{LeftVar: "e", LeftField: "id", RightVar: "s", RightField: "employee_id", Operator: "=="},
			{LeftVar: "e", LeftField: "id", RightVar: "p", RightField: "employee_id", Operator: "=="},
		},
		AggregationVars: []AggregationVariable{
			{Name: "total_sales", Function: "SUM", SourceVar: "s", Field: "amount"},
			{Name: "avg_perf", Function: "AVG", SourceVar: "p", Field: "score"},
		},
	}
	action := &Action{Type: "print"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		network := NewReteNetwork(storage)
		types := []string{"Employee", "Sales", "Performance"}
		for _, typeName := range types {
			typeNode := NewTypeNode(typeName, TypeDefinition{Name: typeName}, storage)
			network.TypeNodes[typeName] = typeNode
			network.RootNode.AddChild(typeNode)
		}

		_ = arb.CreateMultiSourceAccumulatorRule(network, "multi_rule", aggInfo, action)
	}
}

// BenchmarkNetworkConstruction_SmallNetwork benchmarks building a small network (3 types, 2 rules)
func BenchmarkNetworkConstruction_SmallNetwork(b *testing.B) {
	storage := NewMemoryStorage()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		network := NewReteNetwork(storage)
		utils := NewBuilderUtils(storage)
		tb := NewTypeBuilder(utils)
		arb := NewAlphaRuleBuilder(utils)

		// Create 3 types
		types := []interface{}{
			map[string]interface{}{"name": "Person", "fields": []interface{}{
				map[string]interface{}{"name": "id", "type": "number"},
				map[string]interface{}{"name": "age", "type": "number"},
			}},
			map[string]interface{}{"name": "Employee", "fields": []interface{}{
				map[string]interface{}{"name": "person_id", "type": "number"},
			}},
			map[string]interface{}{"name": "Department", "fields": []interface{}{
				map[string]interface{}{"name": "id", "type": "number"},
			}},
		}
		_ = tb.CreateTypeNodes(network, types, storage)

		// Create 2 alpha rules
		action := &Action{Type: "print"}
		_ = arb.CreateAlphaRule(network, "rule1", []map[string]interface{}{{"name": "p"}}, []string{"p"}, []string{"Person"}, map[string]interface{}{}, action)
		_ = arb.CreateAlphaRule(network, "rule2", []map[string]interface{}{{"name": "e"}}, []string{"e"}, []string{"Employee"}, map[string]interface{}{}, action)
	}
}

// BenchmarkNetworkConstruction_MediumNetwork benchmarks building a medium network (5 types, 5 rules)
func BenchmarkNetworkConstruction_MediumNetwork(b *testing.B) {
	storage := NewMemoryStorage()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		network := NewReteNetwork(storage)
		utils := NewBuilderUtils(storage)
		tb := NewTypeBuilder(utils)
		arb := NewAlphaRuleBuilder(utils)
		jrb := NewJoinRuleBuilder(utils)

		// Create 5 types
		typeNames := []string{"Person", "Employee", "Department", "Performance", "Sales"}
		types := make([]interface{}, len(typeNames))
		for idx, name := range typeNames {
			types[idx] = map[string]interface{}{
				"name": name,
				"fields": []interface{}{
					map[string]interface{}{"name": "id", "type": "number"},
				},
			}
		}
		_ = tb.CreateTypeNodes(network, types, storage)

		// Create 3 alpha rules
		action := &Action{Type: "print"}
		for idx := 0; idx < 3; idx++ {
			typeName := typeNames[idx]
			_ = arb.CreateAlphaRule(network, "alpha_"+typeName, []map[string]interface{}{{"name": "x"}}, []string{"x"}, []string{typeName}, map[string]interface{}{}, action)
		}

		// Create 2 join rules
		_ = jrb.CreateJoinRule(network, "join1", []string{"p", "e"}, []string{"Person", "Employee"}, map[string]interface{}{}, action)
		_ = jrb.CreateJoinRule(network, "join2", []string{"e", "d"}, []string{"Employee", "Department"}, map[string]interface{}{}, action)
	}
}

// BenchmarkNetworkConstruction_LargeNetwork benchmarks building a large network (10 types, 15 rules)
func BenchmarkNetworkConstruction_LargeNetwork(b *testing.B) {
	storage := NewMemoryStorage()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		network := NewReteNetwork(storage)
		utils := NewBuilderUtils(storage)
		tb := NewTypeBuilder(utils)
		arb := NewAlphaRuleBuilder(utils)
		jrb := NewJoinRuleBuilder(utils)

		// Create 10 types
		typeNames := []string{"T1", "T2", "T3", "T4", "T5", "T6", "T7", "T8", "T9", "T10"}
		types := make([]interface{}, len(typeNames))
		for idx, name := range typeNames {
			types[idx] = map[string]interface{}{
				"name": name,
				"fields": []interface{}{
					map[string]interface{}{"name": "id", "type": "number"},
					map[string]interface{}{"name": "value", "type": "number"},
				},
			}
		}
		_ = tb.CreateTypeNodes(network, types, storage)

		action := &Action{Type: "print"}

		// Create 10 alpha rules (one per type)
		for idx, typeName := range typeNames {
			_ = arb.CreateAlphaRule(network, "alpha_"+typeName, []map[string]interface{}{{"name": "x"}}, []string{"x"}, []string{typeName}, map[string]interface{}{}, action)
			_ = idx // Use idx
		}

		// Create 5 join rules
		for idx := 0; idx < 5; idx++ {
			_ = jrb.CreateJoinRule(network, "join_"+typeNames[idx], []string{"a", "b"}, []string{typeNames[idx], typeNames[idx+1]}, map[string]interface{}{}, action)
		}
	}
}

// BenchmarkBetaSharing_WithSharing benchmarks join creation with beta sharing enabled
func BenchmarkBetaSharing_WithSharing(b *testing.B) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	jrb := NewJoinRuleBuilder(utils)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		network := NewReteNetwork(storage)
		network.Config = &ReteConfig{BetaSharingEnabled: true}
		network.BetaSharingRegistry = NewBetaSharingRegistry()

		personNode := NewTypeNode("Person", TypeDefinition{Name: "Person"}, storage)
		network.TypeNodes["Person"] = personNode
		network.RootNode.AddChild(personNode)
		employeeNode := NewTypeNode("Employee", TypeDefinition{Name: "Employee"}, storage)
		network.TypeNodes["Employee"] = employeeNode
		network.RootNode.AddChild(employeeNode)

		condition := map[string]interface{}{"type": "comparison"}
		action1 := &Action{Type: "print", Job: &JobCall{Name: "print", Args: []interface{}{"rule1"}}}
		action2 := &Action{Type: "print", Job: &JobCall{Name: "print", Args: []interface{}{"rule2"}}}

		// Create two rules with same join pattern - should share JoinNode
		_ = jrb.CreateJoinRule(network, "rule1", []string{"p", "e"}, []string{"Person", "Employee"}, condition, action1)
		_ = jrb.CreateJoinRule(network, "rule2", []string{"p", "e"}, []string{"Person", "Employee"}, condition, action2)
	}
}

// BenchmarkBetaSharing_WithoutSharing benchmarks join creation without beta sharing
func BenchmarkBetaSharing_WithoutSharing(b *testing.B) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	jrb := NewJoinRuleBuilder(utils)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		network := NewReteNetwork(storage)
		// No beta sharing configuration

		personNode := NewTypeNode("Person", TypeDefinition{Name: "Person"}, storage)
		network.TypeNodes["Person"] = personNode
		network.RootNode.AddChild(personNode)
		employeeNode := NewTypeNode("Employee", TypeDefinition{Name: "Employee"}, storage)
		network.TypeNodes["Employee"] = employeeNode
		network.RootNode.AddChild(employeeNode)

		condition := map[string]interface{}{"type": "comparison"}
		action1 := &Action{Type: "print", Job: &JobCall{Name: "print", Args: []interface{}{"rule1"}}}
		action2 := &Action{Type: "print", Job: &JobCall{Name: "print", Args: []interface{}{"rule2"}}}

		_ = jrb.CreateJoinRule(network, "rule1", []string{"p", "e"}, []string{"Person", "Employee"}, condition, action1)
		_ = jrb.CreateJoinRule(network, "rule2", []string{"p", "e"}, []string{"Person", "Employee"}, condition, action2)
	}
}
