// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"strings"
	"testing"
)

// TestValidateNetwork tests the validateNetwork function
func TestValidateNetwork(t *testing.T) {
	pipeline := NewConstraintPipeline()

	t.Run("returns error for nil network", func(t *testing.T) {
		err := pipeline.validateNetwork(nil)

		if err == nil {
			t.Fatal("Expected error for nil network")
		}

		if !strings.Contains(err.Error(), "nil") {
			t.Errorf("Expected error message to mention 'nil', got: %s", err.Error())
		}
	})

	t.Run("accepts empty network", func(t *testing.T) {
		network := &ReteNetwork{
			TypeNodes:     make(map[string]*TypeNode),
			TerminalNodes: make(map[string]*TerminalNode),
		}

		err := pipeline.validateNetwork(network)

		if err != nil {
			t.Errorf("Expected no error for empty network, got: %v", err)
		}
	})

	t.Run("accepts network with terminal nodes that have actions", func(t *testing.T) {
		network := &ReteNetwork{
			TypeNodes: make(map[string]*TypeNode),
			TerminalNodes: map[string]*TerminalNode{
				"term1": {
					BaseNode: BaseNode{ID: "term1"},
					Action: &Action{
						Jobs: []JobCall{{Name: "log", Args: []interface{}{"test"}}},
					},
				},
			},
		}

		err := pipeline.validateNetwork(network)

		if err != nil {
			t.Errorf("Expected no error for valid network, got: %v", err)
		}
	})

	t.Run("returns error for terminal node without action", func(t *testing.T) {
		network := &ReteNetwork{
			TypeNodes: make(map[string]*TypeNode),
			TerminalNodes: map[string]*TerminalNode{
				"term1": {
					BaseNode: BaseNode{ID: "term1"},
					Action:   nil,
				},
			},
		}

		err := pipeline.validateNetwork(network)

		if err == nil {
			t.Fatal("Expected error for terminal without action")
		}

		if !strings.Contains(err.Error(), "sans action") {
			t.Errorf("Expected error about missing action, got: %s", err.Error())
		}
	})

	t.Run("validates all terminal nodes", func(t *testing.T) {
		network := &ReteNetwork{
			TypeNodes: make(map[string]*TypeNode),
			TerminalNodes: map[string]*TerminalNode{
				"term1": {
					BaseNode: BaseNode{ID: "term1"},
					Action: &Action{
						Jobs: []JobCall{{Name: "log"}},
					},
				},
				"term2": {
					BaseNode: BaseNode{ID: "term2"},
					Action:   nil, // This one is invalid
				},
			},
		}

		err := pipeline.validateNetwork(network)

		if err == nil {
			t.Fatal("Expected error for network with invalid terminal")
		}
	})
}

// TestValidateAction tests the validateAction function
func TestValidateAction(t *testing.T) {
	pipeline := NewConstraintPipeline()

	t.Run("returns error for nil action", func(t *testing.T) {
		err := pipeline.validateAction(nil)

		if err == nil {
			t.Fatal("Expected error for nil action")
		}

		if !strings.Contains(err.Error(), "nil") {
			t.Errorf("Expected error about nil, got: %s", err.Error())
		}
	})

	t.Run("returns error when type is missing", func(t *testing.T) {
		action := map[string]interface{}{
			"name": "someAction",
		}

		err := pipeline.validateAction(action)

		if err == nil {
			t.Fatal("Expected error for missing type")
		}

		if !strings.Contains(err.Error(), "type") {
			t.Errorf("Expected error about type, got: %s", err.Error())
		}
	})

	t.Run("accepts print action with message", func(t *testing.T) {
		action := map[string]interface{}{
			"type":    "print",
			"message": "Hello",
		}

		err := pipeline.validateAction(action)

		if err != nil {
			t.Errorf("Expected no error for valid print action, got: %v", err)
		}
	})

	t.Run("accepts print action with expression", func(t *testing.T) {
		action := map[string]interface{}{
			"type":       "print",
			"expression": "p.name",
		}

		err := pipeline.validateAction(action)

		if err != nil {
			t.Errorf("Expected no error for valid print action, got: %v", err)
		}
	})

	t.Run("accepts PRINT action (uppercase)", func(t *testing.T) {
		action := map[string]interface{}{
			"type":    "PRINT",
			"message": "Hello",
		}

		err := pipeline.validateAction(action)

		if err != nil {
			t.Errorf("Expected no error for uppercase PRINT, got: %v", err)
		}
	})

	t.Run("returns error for print without message or expression", func(t *testing.T) {
		action := map[string]interface{}{
			"type": "print",
		}

		err := pipeline.validateAction(action)

		if err == nil {
			t.Fatal("Expected error for print without message/expression")
		}
	})

	t.Run("accepts assert action with fact", func(t *testing.T) {
		action := map[string]interface{}{
			"type": "assert",
			"fact": map[string]interface{}{
				"type":   "Person",
				"fields": []interface{}{},
			},
		}

		err := pipeline.validateAction(action)

		if err != nil {
			t.Errorf("Expected no error for valid assert action, got: %v", err)
		}
	})

	t.Run("accepts ASSERT action (uppercase)", func(t *testing.T) {
		action := map[string]interface{}{
			"type": "ASSERT",
			"fact": map[string]interface{}{
				"type": "Order",
			},
		}

		err := pipeline.validateAction(action)

		if err != nil {
			t.Errorf("Expected no error for uppercase ASSERT, got: %v", err)
		}
	})

	t.Run("returns error for assert without fact", func(t *testing.T) {
		action := map[string]interface{}{
			"type": "assert",
		}

		err := pipeline.validateAction(action)

		if err == nil {
			t.Fatal("Expected error for assert without fact")
		}

		if !strings.Contains(err.Error(), "assert") && !strings.Contains(err.Error(), "fait") {
			t.Errorf("Expected error about assert/fact, got: %s", err.Error())
		}
	})

	t.Run("accepts retract action with fact", func(t *testing.T) {
		action := map[string]interface{}{
			"type": "retract",
			"fact": "p",
		}

		err := pipeline.validateAction(action)

		if err != nil {
			t.Errorf("Expected no error for valid retract action, got: %v", err)
		}
	})

	t.Run("accepts RETRACT action (uppercase)", func(t *testing.T) {
		action := map[string]interface{}{
			"type": "RETRACT",
			"fact": "o",
		}

		err := pipeline.validateAction(action)

		if err != nil {
			t.Errorf("Expected no error for uppercase RETRACT, got: %v", err)
		}
	})

	t.Run("returns error for retract without fact", func(t *testing.T) {
		action := map[string]interface{}{
			"type": "retract",
		}

		err := pipeline.validateAction(action)

		if err == nil {
			t.Fatal("Expected error for retract without fact")
		}
	})

	t.Run("accepts other action types without specific validation", func(t *testing.T) {
		action := map[string]interface{}{
			"type": "custom_action",
			"data": "some data",
		}

		err := pipeline.validateAction(action)

		if err != nil {
			t.Errorf("Expected no error for custom action type, got: %v", err)
		}
	})
}

// TestValidateRuleExpression tests the validateRuleExpression function
func TestValidateRuleExpression(t *testing.T) {
	pipeline := NewConstraintPipeline()

	t.Run("returns error for nil expression", func(t *testing.T) {
		err := pipeline.validateRuleExpression(nil)

		if err == nil {
			t.Fatal("Expected error for nil expression")
		}
	})

	t.Run("returns error when action is missing", func(t *testing.T) {
		expr := map[string]interface{}{
			"set": []interface{}{"p: Person"},
		}

		err := pipeline.validateRuleExpression(expr)

		if err == nil {
			t.Fatal("Expected error for missing action")
		}

		if !strings.Contains(err.Error(), "action") {
			t.Errorf("Expected error about action, got: %s", err.Error())
		}
	})

	t.Run("returns error when action is not a map", func(t *testing.T) {
		expr := map[string]interface{}{
			"action": "not_a_map",
			"set":    []interface{}{"p: Person"},
		}

		err := pipeline.validateRuleExpression(expr)

		if err == nil {
			t.Fatal("Expected error for invalid action format")
		}

		if !strings.Contains(err.Error(), "format") && !strings.Contains(err.Error(), "invalide") {
			t.Errorf("Expected error about format, got: %s", err.Error())
		}
	})

	t.Run("returns error when action is invalid", func(t *testing.T) {
		expr := map[string]interface{}{
			"action": map[string]interface{}{
				// Missing type
			},
			"set": []interface{}{"p: Person"},
		}

		err := pipeline.validateRuleExpression(expr)

		if err == nil {
			t.Fatal("Expected error for invalid action")
		}
	})

	t.Run("returns error when set is missing", func(t *testing.T) {
		expr := map[string]interface{}{
			"action": map[string]interface{}{
				"type":    "print",
				"message": "test",
			},
		}

		err := pipeline.validateRuleExpression(expr)

		if err == nil {
			t.Fatal("Expected error for missing set")
		}

		if !strings.Contains(err.Error(), "set") && !strings.Contains(err.Error(), "variable") {
			t.Errorf("Expected error about set/variables, got: %s", err.Error())
		}
	})

	t.Run("accepts valid rule expression", func(t *testing.T) {
		expr := map[string]interface{}{
			"action": map[string]interface{}{
				"type":    "print",
				"message": "Found person",
			},
			"set": []interface{}{"p: Person"},
		}

		err := pipeline.validateRuleExpression(expr)

		if err != nil {
			t.Errorf("Expected no error for valid expression, got: %v", err)
		}
	})

	t.Run("validates action within expression", func(t *testing.T) {
		expr := map[string]interface{}{
			"action": map[string]interface{}{
				"type": "print",
				// Missing message and expression
			},
			"set": []interface{}{"p: Person"},
		}

		err := pipeline.validateRuleExpression(expr)

		if err == nil {
			t.Fatal("Expected error for invalid action in expression")
		}

		if !strings.Contains(err.Error(), "action") {
			t.Errorf("Expected error about action, got: %s", err.Error())
		}
	})
}

// TestValidateTypeDefinition tests the validateTypeDefinition function
func TestValidateTypeDefinition(t *testing.T) {
	pipeline := NewConstraintPipeline()

	t.Run("returns error for empty type name", func(t *testing.T) {
		typeDef := map[string]interface{}{
			"fields": []interface{}{},
		}

		err := pipeline.validateTypeDefinition("", typeDef)

		if err == nil {
			t.Fatal("Expected error for empty type name")
		}

		if !strings.Contains(err.Error(), "vide") {
			t.Errorf("Expected error about empty name, got: %s", err.Error())
		}
	})

	t.Run("returns error for nil type definition", func(t *testing.T) {
		err := pipeline.validateTypeDefinition("Person", nil)

		if err == nil {
			t.Fatal("Expected error for nil type definition")
		}

		if !strings.Contains(err.Error(), "nil") {
			t.Errorf("Expected error about nil, got: %s", err.Error())
		}
	})

	t.Run("returns error when fields are missing", func(t *testing.T) {
		typeDef := map[string]interface{}{
			"name": "Person",
		}

		err := pipeline.validateTypeDefinition("Person", typeDef)

		if err == nil {
			t.Fatal("Expected error for missing fields")
		}

		if !strings.Contains(err.Error(), "champs") && !strings.Contains(err.Error(), "fields") {
			t.Errorf("Expected error about fields, got: %s", err.Error())
		}
	})

	t.Run("returns error when fields is not an array", func(t *testing.T) {
		typeDef := map[string]interface{}{
			"fields": "not_an_array",
		}

		err := pipeline.validateTypeDefinition("Person", typeDef)

		if err == nil {
			t.Fatal("Expected error for invalid fields format")
		}

		if !strings.Contains(err.Error(), "format") && !strings.Contains(err.Error(), "invalide") {
			t.Errorf("Expected error about format, got: %s", err.Error())
		}
	})

	t.Run("returns error for empty fields array", func(t *testing.T) {
		typeDef := map[string]interface{}{
			"fields": []interface{}{},
		}

		err := pipeline.validateTypeDefinition("Person", typeDef)

		if err == nil {
			t.Fatal("Expected error for empty fields")
		}

		if !strings.Contains(err.Error(), "vide") {
			t.Errorf("Expected error about empty fields, got: %s", err.Error())
		}
	})

	t.Run("returns error when field is not a map", func(t *testing.T) {
		typeDef := map[string]interface{}{
			"fields": []interface{}{
				"not_a_map",
			},
		}

		err := pipeline.validateTypeDefinition("Person", typeDef)

		if err == nil {
			t.Fatal("Expected error for invalid field format")
		}
	})

	t.Run("returns error when field has no name", func(t *testing.T) {
		typeDef := map[string]interface{}{
			"fields": []interface{}{
				map[string]interface{}{
					"type": "string",
				},
			},
		}

		err := pipeline.validateTypeDefinition("Person", typeDef)

		if err == nil {
			t.Fatal("Expected error for field without name")
		}

		if !strings.Contains(err.Error(), "nom") && !strings.Contains(err.Error(), "name") {
			t.Errorf("Expected error about name, got: %s", err.Error())
		}
	})

	t.Run("returns error when field name is empty", func(t *testing.T) {
		typeDef := map[string]interface{}{
			"fields": []interface{}{
				map[string]interface{}{
					"name": "",
					"type": "string",
				},
			},
		}

		err := pipeline.validateTypeDefinition("Person", typeDef)

		if err == nil {
			t.Fatal("Expected error for empty field name")
		}
	})

	t.Run("returns error when field has no type", func(t *testing.T) {
		typeDef := map[string]interface{}{
			"fields": []interface{}{
				map[string]interface{}{
					"name": "id",
				},
			},
		}

		err := pipeline.validateTypeDefinition("Person", typeDef)

		if err == nil {
			t.Fatal("Expected error for field without type")
		}

		if !strings.Contains(err.Error(), "type") {
			t.Errorf("Expected error about type, got: %s", err.Error())
		}
	})

	t.Run("returns error when field type is empty", func(t *testing.T) {
		typeDef := map[string]interface{}{
			"fields": []interface{}{
				map[string]interface{}{
					"name": "id",
					"type": "",
				},
			},
		}

		err := pipeline.validateTypeDefinition("Person", typeDef)

		if err == nil {
			t.Fatal("Expected error for empty field type")
		}
	})

	t.Run("accepts valid type definition with single field", func(t *testing.T) {
		typeDef := map[string]interface{}{
			"fields": []interface{}{
				map[string]interface{}{
					"name": "id",
					"type": "string",
				},
			},
		}

		err := pipeline.validateTypeDefinition("Person", typeDef)

		if err != nil {
			t.Errorf("Expected no error for valid type definition, got: %v", err)
		}
	})

	t.Run("accepts valid type definition with multiple fields", func(t *testing.T) {
		typeDef := map[string]interface{}{
			"fields": []interface{}{
				map[string]interface{}{
					"name": "id",
					"type": "string",
				},
				map[string]interface{}{
					"name": "name",
					"type": "string",
				},
				map[string]interface{}{
					"name": "age",
					"type": "number",
				},
			},
		}

		err := pipeline.validateTypeDefinition("Person", typeDef)

		if err != nil {
			t.Errorf("Expected no error for valid type definition, got: %v", err)
		}
	})

	t.Run("validates all fields", func(t *testing.T) {
		typeDef := map[string]interface{}{
			"fields": []interface{}{
				map[string]interface{}{
					"name": "id",
					"type": "string",
				},
				map[string]interface{}{
					"name": "invalid",
					// Missing type
				},
			},
		}

		err := pipeline.validateTypeDefinition("Person", typeDef)

		if err == nil {
			t.Fatal("Expected error for invalid field in type definition")
		}
	})
}

// TestValidateAggregationInfo tests the validateAggregationInfo function
func TestValidateAggregationInfo(t *testing.T) {
	pipeline := NewConstraintPipeline()

	t.Run("returns error for nil aggregation info", func(t *testing.T) {
		err := pipeline.validateAggregationInfo(nil)

		if err == nil {
			t.Fatal("Expected error for nil aggregation info")
		}

		if !strings.Contains(err.Error(), "nil") {
			t.Errorf("Expected error about nil, got: %s", err.Error())
		}
	})

	t.Run("returns error for empty function", func(t *testing.T) {
		aggInfo := &AggregationInfo{
			Function: "",
			Operator: ">=",
		}

		err := pipeline.validateAggregationInfo(aggInfo)

		if err == nil {
			t.Fatal("Expected error for empty function")
		}

		if !strings.Contains(err.Error(), "fonction") && !strings.Contains(err.Error(), "vide") {
			t.Errorf("Expected error about empty function, got: %s", err.Error())
		}
	})

	t.Run("returns error for invalid function", func(t *testing.T) {
		aggInfo := &AggregationInfo{
			Function: "INVALID",
			Operator: ">=",
		}

		err := pipeline.validateAggregationInfo(aggInfo)

		if err == nil {
			t.Fatal("Expected error for invalid function")
		}

		if !strings.Contains(err.Error(), "fonction") && !strings.Contains(err.Error(), "invalide") {
			t.Errorf("Expected error about invalid function, got: %s", err.Error())
		}
	})

	t.Run("accepts AVG function", func(t *testing.T) {
		aggInfo := &AggregationInfo{
			Function: "AVG",
			Operator: ">=",
		}

		err := pipeline.validateAggregationInfo(aggInfo)

		if err != nil {
			t.Errorf("Expected no error for AVG function, got: %v", err)
		}
	})

	t.Run("accepts SUM function", func(t *testing.T) {
		aggInfo := &AggregationInfo{
			Function: "SUM",
			Operator: "<=",
		}

		err := pipeline.validateAggregationInfo(aggInfo)

		if err != nil {
			t.Errorf("Expected no error for SUM function, got: %v", err)
		}
	})

	t.Run("accepts COUNT function", func(t *testing.T) {
		aggInfo := &AggregationInfo{
			Function: "COUNT",
			Operator: ">",
		}

		err := pipeline.validateAggregationInfo(aggInfo)

		if err != nil {
			t.Errorf("Expected no error for COUNT function, got: %v", err)
		}
	})

	t.Run("accepts MIN function", func(t *testing.T) {
		aggInfo := &AggregationInfo{
			Function: "MIN",
			Operator: "<",
		}

		err := pipeline.validateAggregationInfo(aggInfo)

		if err != nil {
			t.Errorf("Expected no error for MIN function, got: %v", err)
		}
	})

	t.Run("accepts MAX function", func(t *testing.T) {
		aggInfo := &AggregationInfo{
			Function: "MAX",
			Operator: "==",
		}

		err := pipeline.validateAggregationInfo(aggInfo)

		if err != nil {
			t.Errorf("Expected no error for MAX function, got: %v", err)
		}
	})

	t.Run("accepts ACCUMULATE function", func(t *testing.T) {
		aggInfo := &AggregationInfo{
			Function: "ACCUMULATE",
			Operator: "!=",
		}

		err := pipeline.validateAggregationInfo(aggInfo)

		if err != nil {
			t.Errorf("Expected no error for ACCUMULATE function, got: %v", err)
		}
	})

	t.Run("returns error for empty operator", func(t *testing.T) {
		aggInfo := &AggregationInfo{
			Function: "AVG",
			Operator: "",
		}

		err := pipeline.validateAggregationInfo(aggInfo)

		if err == nil {
			t.Fatal("Expected error for empty operator")
		}

		if !strings.Contains(err.Error(), "opérateur") && !strings.Contains(err.Error(), "vide") {
			t.Errorf("Expected error about empty operator, got: %s", err.Error())
		}
	})

	t.Run("returns error for invalid operator", func(t *testing.T) {
		aggInfo := &AggregationInfo{
			Function: "AVG",
			Operator: "invalid",
		}

		err := pipeline.validateAggregationInfo(aggInfo)

		if err == nil {
			t.Fatal("Expected error for invalid operator")
		}

		if !strings.Contains(err.Error(), "opérateur") && !strings.Contains(err.Error(), "invalide") {
			t.Errorf("Expected error about invalid operator, got: %s", err.Error())
		}
	})

	t.Run("accepts all valid operators", func(t *testing.T) {
		operators := []string{">=", "<=", ">", "<", "==", "!="}

		for _, op := range operators {
			aggInfo := &AggregationInfo{
				Function: "COUNT",
				Operator: op,
			}

			err := pipeline.validateAggregationInfo(aggInfo)

			if err != nil {
				t.Errorf("Expected no error for operator %s, got: %v", op, err)
			}
		}
	})

	t.Run("returns error when JoinField without MainField", func(t *testing.T) {
		aggInfo := &AggregationInfo{
			Function:  "AVG",
			Operator:  ">=",
			JoinField: "employee_id",
			MainField: "",
		}

		err := pipeline.validateAggregationInfo(aggInfo)

		if err == nil {
			t.Fatal("Expected error for JoinField without MainField")
		}

		if !strings.Contains(err.Error(), "principal") && !strings.Contains(err.Error(), "manquant") {
			t.Errorf("Expected error about missing main field, got: %s", err.Error())
		}
	})

	t.Run("returns error when MainField without JoinField", func(t *testing.T) {
		aggInfo := &AggregationInfo{
			Function:  "AVG",
			Operator:  ">=",
			JoinField: "",
			MainField: "id",
		}

		err := pipeline.validateAggregationInfo(aggInfo)

		if err == nil {
			t.Fatal("Expected error for MainField without JoinField")
		}

		if !strings.Contains(err.Error(), "agrégé") && !strings.Contains(err.Error(), "manquant") {
			t.Errorf("Expected error about missing agg field, got: %s", err.Error())
		}
	})

	t.Run("accepts both JoinField and MainField together", func(t *testing.T) {
		aggInfo := &AggregationInfo{
			Function:  "AVG",
			Operator:  ">=",
			JoinField: "employee_id",
			MainField: "id",
		}

		err := pipeline.validateAggregationInfo(aggInfo)

		if err != nil {
			t.Errorf("Expected no error for valid join fields, got: %v", err)
		}
	})

	t.Run("accepts empty join fields", func(t *testing.T) {
		aggInfo := &AggregationInfo{
			Function:  "COUNT",
			Operator:  ">",
			JoinField: "",
			MainField: "",
		}

		err := pipeline.validateAggregationInfo(aggInfo)

		if err != nil {
			t.Errorf("Expected no error for no join fields, got: %v", err)
		}
	})
}

// TestValidateJoinCondition tests the validateJoinCondition function
func TestValidateJoinCondition(t *testing.T) {
	pipeline := NewConstraintPipeline()

	t.Run("returns error for nil condition", func(t *testing.T) {
		err := pipeline.validateJoinCondition(nil)

		if err == nil {
			t.Fatal("Expected error for nil condition")
		}

		if !strings.Contains(err.Error(), "nil") {
			t.Errorf("Expected error about nil, got: %s", err.Error())
		}
	})

	t.Run("returns error when type is missing", func(t *testing.T) {
		condition := map[string]interface{}{
			"data": "some data",
		}

		err := pipeline.validateJoinCondition(condition)

		if err == nil {
			t.Fatal("Expected error for missing type")
		}

		if !strings.Contains(err.Error(), "type") {
			t.Errorf("Expected error about type, got: %s", err.Error())
		}
	})

	t.Run("accepts simple condition", func(t *testing.T) {
		condition := map[string]interface{}{
			"type": "simple",
		}

		err := pipeline.validateJoinCondition(condition)

		if err != nil {
			t.Errorf("Expected no error for simple condition, got: %v", err)
		}
	})

	t.Run("accepts passthrough condition", func(t *testing.T) {
		condition := map[string]interface{}{
			"type": "passthrough",
		}

		err := pipeline.validateJoinCondition(condition)

		if err != nil {
			t.Errorf("Expected no error for passthrough condition, got: %v", err)
		}
	})

	t.Run("accepts constraint condition with constraint", func(t *testing.T) {
		condition := map[string]interface{}{
			"type":       "constraint",
			"constraint": map[string]interface{}{"field": "id"},
		}

		err := pipeline.validateJoinCondition(condition)

		if err != nil {
			t.Errorf("Expected no error for valid constraint condition, got: %v", err)
		}
	})

	t.Run("returns error for constraint condition without constraint", func(t *testing.T) {
		condition := map[string]interface{}{
			"type": "constraint",
		}

		err := pipeline.validateJoinCondition(condition)

		if err == nil {
			t.Fatal("Expected error for constraint condition without constraint")
		}

		if !strings.Contains(err.Error(), "constraint") {
			t.Errorf("Expected error about constraint, got: %s", err.Error())
		}
	})

	t.Run("accepts negation condition with all required fields", func(t *testing.T) {
		condition := map[string]interface{}{
			"type":      "negation",
			"negated":   true,
			"condition": map[string]interface{}{"type": "simple"},
		}

		err := pipeline.validateJoinCondition(condition)

		if err != nil {
			t.Errorf("Expected no error for valid negation condition, got: %v", err)
		}
	})

	t.Run("returns error for negation without negated flag", func(t *testing.T) {
		condition := map[string]interface{}{
			"type":      "negation",
			"condition": map[string]interface{}{"type": "simple"},
		}

		err := pipeline.validateJoinCondition(condition)

		if err == nil {
			t.Fatal("Expected error for negation without negated flag")
		}

		if !strings.Contains(err.Error(), "negated") {
			t.Errorf("Expected error about negated flag, got: %s", err.Error())
		}
	})

	t.Run("returns error for negation without condition", func(t *testing.T) {
		condition := map[string]interface{}{
			"type":    "negation",
			"negated": true,
		}

		err := pipeline.validateJoinCondition(condition)

		if err == nil {
			t.Fatal("Expected error for negation without condition")
		}

		if !strings.Contains(err.Error(), "condition") {
			t.Errorf("Expected error about condition, got: %s", err.Error())
		}
	})

	t.Run("returns error for unknown condition type", func(t *testing.T) {
		condition := map[string]interface{}{
			"type": "unknown_type",
		}

		err := pipeline.validateJoinCondition(condition)

		if err == nil {
			t.Fatal("Expected error for unknown condition type")
		}

		if !strings.Contains(err.Error(), "inconnu") && !strings.Contains(err.Error(), "unknown") {
			t.Errorf("Expected error about unknown type, got: %s", err.Error())
		}
	})
}
