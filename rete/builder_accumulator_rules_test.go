// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"
)

func TestNewAccumulatorRuleBuilder(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	arb := NewAccumulatorRuleBuilder(utils)
	if arb == nil {
		t.Fatal("NewAccumulatorRuleBuilder returned nil")
	}
	if arb.utils != utils {
		t.Error("AccumulatorRuleBuilder.utils not set correctly")
	}
}
func TestAccumulatorRuleBuilder_IsMultiSourceAggregation(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	arb := NewAccumulatorRuleBuilder(utils)
	tests := []struct {
		name    string
		exprMap map[string]interface{}
		want    bool
	}{
		{
			name: "no patterns",
			exprMap: map[string]interface{}{
				"other": "data",
			},
			want: false,
		},
		{
			name: "single pattern",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{
						"variables": []interface{}{},
					},
				},
			},
			want: false,
		},
		{
			name: "two patterns",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{},
					map[string]interface{}{},
				},
			},
			want: false,
		},
		{
			name: "three patterns - multi-source",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{},
					map[string]interface{}{},
					map[string]interface{}{},
				},
			},
			want: true,
		},
		{
			name: "multiple aggregation variables",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{
						"variables": []interface{}{
							map[string]interface{}{
								"type": "aggregationVariable",
								"name": "agg1",
							},
							map[string]interface{}{
								"type": "aggregationVariable",
								"name": "agg2",
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "single aggregation variable",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{
						"variables": []interface{}{
							map[string]interface{}{
								"type": "aggregationVariable",
								"name": "agg1",
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "no aggregation variables",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{
						"variables": []interface{}{
							map[string]interface{}{
								"type": "regular",
								"name": "var1",
							},
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := arb.IsMultiSourceAggregation(tt.exprMap)
			if got != tt.want {
				t.Errorf("IsMultiSourceAggregation() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestAccumulatorRuleBuilder_CreateAccumulatorRule(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	arb := NewAccumulatorRuleBuilder(utils)
	t.Run("create simple accumulator rule", func(t *testing.T) {
		network := NewReteNetwork(storage)
		// Create TypeNodes
		employeeNode := NewTypeNode("Employee", TypeDefinition{
			Type: "type",
			Name: "Employee",
			Fields: []Field{
				{Name: "id", Type: "number"},
			},
		}, storage)
		network.TypeNodes["Employee"] = employeeNode
		network.RootNode.AddChild(employeeNode)
		performanceNode := NewTypeNode("Performance", TypeDefinition{
			Type: "type",
			Name: "Performance",
			Fields: []Field{
				{Name: "employee_id", Type: "number"},
				{Name: "score", Type: "number"},
			},
		}, storage)
		network.TypeNodes["Performance"] = performanceNode
		network.RootNode.AddChild(performanceNode)
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
		action := &Action{
			Type: "print",
			Job:  &JobCall{Name: "print", Args: []interface{}{"High performer"}},
		}
		err := arb.CreateAccumulatorRule(
			network,
			"high_performer_rule",
			variables,
			variableNames,
			variableTypes,
			aggInfo,
			action,
		)
		if err != nil {
			t.Fatalf("CreateAccumulatorRule failed: %v", err)
		}
		// Verify AccumulatorNode was created
		accumNode, exists := network.BetaNodes["high_performer_rule_accum"]
		if !exists {
			t.Fatal("AccumulatorNode not created")
		}
		accumNodeTyped, ok := accumNode.(*AccumulatorNode)
		if !ok {
			t.Fatal("BetaNode is not an AccumulatorNode")
		}
		// Verify AccumulatorNode configuration
		if accumNodeTyped.MainVariable != "e" {
			t.Errorf("MainVariable = %q, want 'e'", accumNodeTyped.MainVariable)
		}
		if accumNodeTyped.AggVariable != "p" {
			t.Errorf("AggVariable = %q, want 'p'", accumNodeTyped.AggVariable)
		}
		if accumNodeTyped.AggregateFunc != "AVG" {
			t.Errorf("Function = %q, want 'AVG'", accumNodeTyped.AggregateFunc)
		}
		if accumNodeTyped.Field != "score" {
			t.Errorf("Field = %q, want 'score'", accumNodeTyped.Field)
		}
		// Verify TerminalNode connection
		if len(accumNodeTyped.Children) != 1 {
			t.Fatalf("AccumulatorNode should have 1 child (TerminalNode), got %d", len(accumNodeTyped.Children))
		}
		_, isTerminal := accumNodeTyped.Children[0].(*TerminalNode)
		if !isTerminal {
			t.Error("AccumulatorNode child should be TerminalNode")
		}
		// Verify TypeNode connections (via pass-through alphas)
		if len(employeeNode.Children) == 0 {
			t.Error("Employee TypeNode should have children")
		}
		if len(performanceNode.Children) == 0 {
			t.Error("Performance TypeNode should have children")
		}
	})
	t.Run("error on missing variables", func(t *testing.T) {
		network := NewReteNetwork(storage)
		aggInfo := &AggregationInfo{
			Function: "SUM",
		}
		err := arb.CreateAccumulatorRule(
			network,
			"bad_rule",
			[]map[string]interface{}{},
			[]string{},
			[]string{},
			aggInfo,
			&Action{Type: "print"},
		)
		if err == nil {
			t.Error("Expected error for missing variables, got nil")
		}
	})
	t.Run("with SUM function", func(t *testing.T) {
		network := NewReteNetwork(storage)
		accountNode := NewTypeNode("Account", TypeDefinition{Name: "Account"}, storage)
		network.TypeNodes["Account"] = accountNode
		network.RootNode.AddChild(accountNode)
		transactionNode := NewTypeNode("Transaction", TypeDefinition{Name: "Transaction"}, storage)
		network.TypeNodes["Transaction"] = transactionNode
		network.RootNode.AddChild(transactionNode)
		aggInfo := &AggregationInfo{
			MainVariable: "a",
			MainType:     "Account",
			AggVariable:  "t",
			AggType:      "Transaction",
			Field:        "amount",
			JoinField:    "account_id",
			MainField:    "id",
			Function:     "SUM",
			Operator:     ">=",
			Threshold:    1000.0,
		}
		err := arb.CreateAccumulatorRule(
			network,
			"high_balance_rule",
			[]map[string]interface{}{{"name": "a"}, {"name": "t"}},
			[]string{"a", "t"},
			[]string{"Account", "Transaction"},
			aggInfo,
			&Action{Type: "print"},
		)
		if err != nil {
			t.Fatalf("CreateAccumulatorRule with SUM failed: %v", err)
		}
		accumNode := network.BetaNodes["high_balance_rule_accum"]
		if accumNode == nil {
			t.Fatal("AccumulatorNode not created")
		}
		accumNodeTyped := accumNode.(*AccumulatorNode)
		if accumNodeTyped.AggregateFunc != "SUM" {
			t.Errorf("Function = %q, want 'SUM'", accumNodeTyped.AggregateFunc)
		}
	})
	t.Run("with COUNT function", func(t *testing.T) {
		network := NewReteNetwork(storage)
		customerNode := NewTypeNode("Customer", TypeDefinition{Name: "Customer"}, storage)
		network.TypeNodes["Customer"] = customerNode
		network.RootNode.AddChild(customerNode)
		orderNode := NewTypeNode("Order", TypeDefinition{Name: "Order"}, storage)
		network.TypeNodes["Order"] = orderNode
		network.RootNode.AddChild(orderNode)
		aggInfo := &AggregationInfo{
			MainVariable: "c",
			MainType:     "Customer",
			AggVariable:  "o",
			AggType:      "Order",
			Field:        "id",
			JoinField:    "customer_id",
			MainField:    "id",
			Function:     "COUNT",
			Operator:     ">",
			Threshold:    5.0,
		}
		err := arb.CreateAccumulatorRule(
			network,
			"frequent_customer_rule",
			[]map[string]interface{}{{"name": "c"}, {"name": "o"}},
			[]string{"c", "o"},
			[]string{"Customer", "Order"},
			aggInfo,
			&Action{Type: "print"},
		)
		if err != nil {
			t.Fatalf("CreateAccumulatorRule with COUNT failed: %v", err)
		}
		accumNode := network.BetaNodes["frequent_customer_rule_accum"]
		if accumNode == nil {
			t.Fatal("AccumulatorNode not created")
		}
	})
}
func TestAccumulatorRuleBuilder_CreateMultiSourceAccumulatorRule(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	arb := NewAccumulatorRuleBuilder(utils)
	t.Run("create multi-source accumulator rule", func(t *testing.T) {
		network := NewReteNetwork(storage)
		// Create TypeNodes
		employeeNode := NewTypeNode("Employee", TypeDefinition{Name: "Employee"}, storage)
		network.TypeNodes["Employee"] = employeeNode
		network.RootNode.AddChild(employeeNode)
		salesNode := NewTypeNode("Sales", TypeDefinition{Name: "Sales"}, storage)
		network.TypeNodes["Sales"] = salesNode
		network.RootNode.AddChild(salesNode)
		performanceNode := NewTypeNode("Performance", TypeDefinition{Name: "Performance"}, storage)
		network.TypeNodes["Performance"] = performanceNode
		network.RootNode.AddChild(performanceNode)
		aggInfo := &AggregationInfo{
			MainVariable: "e",
			MainType:     "Employee",
			SourcePatterns: []SourcePattern{
				{
					Variable: "s",
					Type:     "Sales",
				},
				{
					Variable: "p",
					Type:     "Performance",
				},
			},
			JoinConditions: []JoinCondition{
				{
					LeftVar:    "e",
					LeftField:  "id",
					RightVar:   "s",
					RightField: "employee_id",
					Operator:   "==",
				},
				{
					LeftVar:    "e",
					LeftField:  "id",
					RightVar:   "p",
					RightField: "employee_id",
					Operator:   "==",
				},
			},
			AggregationVars: []AggregationVariable{
				{
					Name:      "total_sales",
					Function:  "SUM",
					SourceVar: "s",
					Field:     "amount",
					Operator:  ">",
					Threshold: 10000.0,
				},
				{
					Name:      "avg_perf",
					Function:  "AVG",
					SourceVar: "p",
					Field:     "score",
					Operator:  ">=",
					Threshold: 4.5,
				},
			},
		}
		action := &Action{
			Type: "print",
			Job:  &JobCall{Name: "print", Args: []interface{}{"Top performer found"}},
		}
		err := arb.CreateMultiSourceAccumulatorRule(network, "top_performer_rule", aggInfo, action)
		if err != nil {
			t.Fatalf("CreateMultiSourceAccumulatorRule failed: %v", err)
		}
		// Verify MultiSourceAccumulatorNode was created
		var msAccumNode *MultiSourceAccumulatorNode
		for _, betaNode := range network.BetaNodes {
			if node, ok := betaNode.(*MultiSourceAccumulatorNode); ok {
				msAccumNode = node
				break
			}
		}
		if msAccumNode == nil {
			t.Fatal("MultiSourceAccumulatorNode not created")
		}
		// Verify configuration
		if msAccumNode.MainVariable != "e" {
			t.Errorf("MainVariable = %q, want 'e'", msAccumNode.MainVariable)
		}
		if len(msAccumNode.AggregationVars) != 2 {
			t.Errorf("AggregationVars count = %d, want 2", len(msAccumNode.AggregationVars))
		}
		// Verify JoinNodes were created for sources
		joinNodeCount := 0
		for _, betaNode := range network.BetaNodes {
			if _, ok := betaNode.(*JoinNode); ok {
				joinNodeCount++
			}
		}
		if joinNodeCount < 2 {
			t.Errorf("Expected at least 2 JoinNodes for 2 sources, got %d", joinNodeCount)
		}
		// Verify TerminalNode connection
		if len(msAccumNode.Children) != 1 {
			t.Fatalf("MultiSourceAccumulatorNode should have 1 child, got %d", len(msAccumNode.Children))
		}
		_, isTerminal := msAccumNode.Children[0].(*TerminalNode)
		if !isTerminal {
			t.Error("MultiSourceAccumulatorNode child should be TerminalNode")
		}
	})
	t.Run("error on missing main type", func(t *testing.T) {
		network := NewReteNetwork(storage)
		aggInfo := &AggregationInfo{
			MainVariable: "e",
			MainType:     "NonExistent",
			SourcePatterns: []SourcePattern{
				{Variable: "s", Type: "Sales"},
			},
		}
		err := arb.CreateMultiSourceAccumulatorRule(network, "bad_rule", aggInfo, &Action{})
		if err == nil {
			t.Error("Expected error for missing main type, got nil")
		}
	})
	t.Run("error on missing source type", func(t *testing.T) {
		network := NewReteNetwork(storage)
		employeeNode := NewTypeNode("Employee", TypeDefinition{Name: "Employee"}, storage)
		network.TypeNodes["Employee"] = employeeNode
		network.RootNode.AddChild(employeeNode)
		aggInfo := &AggregationInfo{
			MainVariable: "e",
			MainType:     "Employee",
			SourcePatterns: []SourcePattern{
				{Variable: "s", Type: "NonExistent"},
			},
		}
		err := arb.CreateMultiSourceAccumulatorRule(network, "bad_rule", aggInfo, &Action{})
		if err == nil {
			t.Error("Expected error for missing source type, got nil")
		}
	})
	t.Run("with three sources", func(t *testing.T) {
		network := NewReteNetwork(storage)
		// Create TypeNodes
		types := []string{"Main", "Source1", "Source2", "Source3"}
		for _, typeName := range types {
			typeNode := NewTypeNode(typeName, TypeDefinition{Name: typeName}, storage)
			network.TypeNodes[typeName] = typeNode
			network.RootNode.AddChild(typeNode)
		}
		aggInfo := &AggregationInfo{
			MainVariable: "m",
			MainType:     "Main",
			SourcePatterns: []SourcePattern{
				{Variable: "s1", Type: "Source1"},
				{Variable: "s2", Type: "Source2"},
				{Variable: "s3", Type: "Source3"},
			},
			JoinConditions: []JoinCondition{
				{LeftVar: "m", LeftField: "id", RightVar: "s1", RightField: "main_id", Operator: "=="},
				{LeftVar: "m", LeftField: "id", RightVar: "s2", RightField: "main_id", Operator: "=="},
				{LeftVar: "m", LeftField: "id", RightVar: "s3", RightField: "main_id", Operator: "=="},
			},
			AggregationVars: []AggregationVariable{
				{Name: "agg1", Function: "SUM", SourceVar: "s1", Field: "value"},
				{Name: "agg2", Function: "AVG", SourceVar: "s2", Field: "value"},
				{Name: "agg3", Function: "COUNT", SourceVar: "s3", Field: "id"},
			},
		}
		err := arb.CreateMultiSourceAccumulatorRule(network, "three_source_rule", aggInfo, &Action{Type: "print"})
		if err != nil {
			t.Fatalf("CreateMultiSourceAccumulatorRule with 3 sources failed: %v", err)
		}
		// Should create 3 JoinNodes (one for each source)
		joinNodeCount := 0
		for _, betaNode := range network.BetaNodes {
			if _, ok := betaNode.(*JoinNode); ok {
				joinNodeCount++
			}
		}
		if joinNodeCount != 3 {
			t.Errorf("Expected 3 JoinNodes for 3 sources, got %d", joinNodeCount)
		}
	})
}
func TestAccumulatorRuleBuilder_Integration(t *testing.T) {
	// Integration test: Create multiple accumulator rules and verify network structure
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	arb := NewAccumulatorRuleBuilder(utils)
	network := NewReteNetwork(storage)
	// Setup: Create TypeNodes
	employeeNode := NewTypeNode("Employee", TypeDefinition{
		Type: "type",
		Name: "Employee",
		Fields: []Field{
			{Name: "id", Type: "number"},
			{Name: "department", Type: "string"},
		},
	}, storage)
	network.TypeNodes["Employee"] = employeeNode
	network.RootNode.AddChild(employeeNode)
	performanceNode := NewTypeNode("Performance", TypeDefinition{
		Type: "type",
		Name: "Performance",
		Fields: []Field{
			{Name: "employee_id", Type: "number"},
			{Name: "score", Type: "number"},
		},
	}, storage)
	network.TypeNodes["Performance"] = performanceNode
	network.RootNode.AddChild(performanceNode)
	salesNode := NewTypeNode("Sales", TypeDefinition{
		Type: "type",
		Name: "Sales",
		Fields: []Field{
			{Name: "employee_id", Type: "number"},
			{Name: "amount", Type: "number"},
		},
	}, storage)
	network.TypeNodes["Sales"] = salesNode
	network.RootNode.AddChild(salesNode)
	// Rule 1: Simple accumulator - AVG performance score
	aggInfo1 := &AggregationInfo{
		MainVariable: "e",
		MainType:     "Employee",
		AggVariable:  "p",
		AggType:      "Performance",
		Field:        "score",
		JoinField:    "employee_id",
		MainField:    "id",
		Function:     "AVG",
		Operator:     ">=",
		Threshold:    4.0,
	}
	err := arb.CreateAccumulatorRule(
		network,
		"high_performer",
		[]map[string]interface{}{{"name": "e"}, {"name": "p"}},
		[]string{"e", "p"},
		[]string{"Employee", "Performance"},
		aggInfo1,
		&Action{Type: "print", Job: &JobCall{Name: "print", Args: []interface{}{"High performer"}}},
	)
	if err != nil {
		t.Fatalf("Failed to create high_performer rule: %v", err)
	}
	// Rule 2: Simple accumulator - SUM sales amount
	aggInfo2 := &AggregationInfo{
		MainVariable: "e",
		MainType:     "Employee",
		AggVariable:  "s",
		AggType:      "Sales",
		Field:        "amount",
		JoinField:    "employee_id",
		MainField:    "id",
		Function:     "SUM",
		Operator:     ">",
		Threshold:    50000.0,
	}
	err = arb.CreateAccumulatorRule(
		network,
		"high_sales",
		[]map[string]interface{}{{"name": "e"}, {"name": "s"}},
		[]string{"e", "s"},
		[]string{"Employee", "Sales"},
		aggInfo2,
		&Action{Type: "print", Job: &JobCall{Name: "print", Args: []interface{}{"Top seller"}}},
	)
	if err != nil {
		t.Fatalf("Failed to create high_sales rule: %v", err)
	}
	// Verify both AccumulatorNodes created
	accumCount := 0
	for _, betaNode := range network.BetaNodes {
		if _, ok := betaNode.(*AccumulatorNode); ok {
			accumCount++
		}
	}
	if accumCount != 2 {
		t.Errorf("Expected 2 AccumulatorNodes, got %d", accumCount)
	}
	// Verify each AccumulatorNode has a TerminalNode
	for nodeID, betaNode := range network.BetaNodes {
		accumNode, ok := betaNode.(*AccumulatorNode)
		if !ok {
			continue
		}
		if len(accumNode.Children) != 1 {
			t.Errorf("AccumulatorNode %s should have 1 child, got %d", nodeID, len(accumNode.Children))
		}
		_, isTerminal := accumNode.Children[0].(*TerminalNode)
		if !isTerminal {
			t.Errorf("AccumulatorNode %s child is not a TerminalNode", nodeID)
		}
	}
	// Verify TypeNodes have connections
	if len(employeeNode.Children) == 0 {
		t.Error("Employee TypeNode should have children (pass-through alphas)")
	}
	if len(performanceNode.Children) == 0 {
		t.Error("Performance TypeNode should have children")
	}
	if len(salesNode.Children) == 0 {
		t.Error("Sales TypeNode should have children")
	}
}
