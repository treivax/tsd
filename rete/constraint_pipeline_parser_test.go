// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// TestExtractComponents tests the extractComponents function
func TestExtractComponents(t *testing.T) {
	pipeline := NewConstraintPipeline()

	t.Run("extracts types and expressions successfully", func(t *testing.T) {
		resultMap := map[string]interface{}{
			"types": []interface{}{
				map[string]interface{}{"name": "Person"},
				map[string]interface{}{"name": "Order"},
			},
			"expressions": []interface{}{
				map[string]interface{}{"rule": "r1"},
				map[string]interface{}{"rule": "r2"},
			},
		}

		types, expressions, err := pipeline.extractComponents(resultMap)
		if err != nil {
			t.Fatalf("extractComponents failed: %v", err)
		}

		if len(types) != 2 {
			t.Errorf("Expected 2 types, got %d", len(types))
		}

		if len(expressions) != 2 {
			t.Errorf("Expected 2 expressions, got %d", len(expressions))
		}
	})

	t.Run("returns error when types missing", func(t *testing.T) {
		resultMap := map[string]interface{}{
			"expressions": []interface{}{
				map[string]interface{}{"rule": "r1"},
			},
		}

		_, _, err := pipeline.extractComponents(resultMap)
		if err == nil {
			t.Error("Expected error when types missing")
		}
	})

	t.Run("returns error when types invalid format", func(t *testing.T) {
		resultMap := map[string]interface{}{
			"types": "not a slice",
			"expressions": []interface{}{
				map[string]interface{}{"rule": "r1"},
			},
		}

		_, _, err := pipeline.extractComponents(resultMap)
		if err == nil {
			t.Error("Expected error when types has invalid format")
		}
	})

	t.Run("returns error when expressions missing", func(t *testing.T) {
		resultMap := map[string]interface{}{
			"types": []interface{}{
				map[string]interface{}{"name": "Person"},
			},
		}

		_, _, err := pipeline.extractComponents(resultMap)
		if err == nil {
			t.Error("Expected error when expressions missing")
		}
	})

	t.Run("returns error when expressions invalid format", func(t *testing.T) {
		resultMap := map[string]interface{}{
			"types": []interface{}{
				map[string]interface{}{"name": "Person"},
			},
			"expressions": "not a slice",
		}

		_, _, err := pipeline.extractComponents(resultMap)
		if err == nil {
			t.Error("Expected error when expressions has invalid format")
		}
	})
}

// TestExtractAndStoreActions tests the extractAndStoreActions function
func TestExtractAndStoreActions(t *testing.T) {
	pipeline := NewConstraintPipeline()

	t.Run("no actions present is not an error", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		resultMap := map[string]interface{}{
			"types": []interface{}{},
		}

		err := pipeline.extractAndStoreActions(network, resultMap)
		if err != nil {
			t.Errorf("Expected no error when actions missing, got: %v", err)
		}
	})

	t.Run("extracts and stores actions successfully", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		resultMap := map[string]interface{}{
			"actions": []interface{}{
				map[string]interface{}{
					"name": "log",
					"parameters": []interface{}{
						map[string]interface{}{
							"name": "msg",
							"type": "string",
						},
					},
				},
			},
		}

		err := pipeline.extractAndStoreActions(network, resultMap)
		if err != nil {
			t.Fatalf("extractAndStoreActions failed: %v", err)
		}

		// Verify action was stored
		if len(network.Actions) != 1 {
			t.Errorf("Expected 1 action definition, got %d", len(network.Actions))
		}

		if network.Actions[0].Name != "log" {
			t.Errorf("Expected action name 'log', got '%s'", network.Actions[0].Name)
		}

		if len(network.Actions[0].Parameters) != 1 {
			t.Errorf("Expected 1 parameter, got %d", len(network.Actions[0].Parameters))
		}
	})

	t.Run("returns error for invalid actions format", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		resultMap := map[string]interface{}{
			"actions": "not a slice",
		}

		err := pipeline.extractAndStoreActions(network, resultMap)
		if err == nil {
			t.Error("Expected error for invalid actions format")
		}
	})

	t.Run("returns error for invalid action item format", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		resultMap := map[string]interface{}{
			"actions": []interface{}{
				"not a map",
			},
		}

		err := pipeline.extractAndStoreActions(network, resultMap)
		if err == nil {
			t.Error("Expected error for invalid action item format")
		}
	})

	t.Run("returns error when action name missing", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		resultMap := map[string]interface{}{
			"actions": []interface{}{
				map[string]interface{}{
					"parameters": []interface{}{},
				},
			},
		}

		err := pipeline.extractAndStoreActions(network, resultMap)
		if err == nil {
			t.Error("Expected error when action name missing")
		}
	})

	t.Run("returns error for invalid parameters format", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		resultMap := map[string]interface{}{
			"actions": []interface{}{
				map[string]interface{}{
					"name":       "log",
					"parameters": "not a slice",
				},
			},
		}

		err := pipeline.extractAndStoreActions(network, resultMap)
		if err == nil {
			t.Error("Expected error for invalid parameters format")
		}
	})

	t.Run("handles action without parameters", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		resultMap := map[string]interface{}{
			"actions": []interface{}{
				map[string]interface{}{
					"name": "notify",
				},
			},
		}

		err := pipeline.extractAndStoreActions(network, resultMap)
		if err != nil {
			t.Fatalf("extractAndStoreActions failed: %v", err)
		}

		if len(network.Actions) != 1 {
			t.Errorf("Expected 1 action definition, got %d", len(network.Actions))
		}
	})

	t.Run("replaces existing action with same name", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)

		// Add initial action
		network.Actions = append(network.Actions, ActionDefinition{
			Name: "log",
			Type: "actionDefinition",
		})

		resultMap := map[string]interface{}{
			"actions": []interface{}{
				map[string]interface{}{
					"name": "log",
					"parameters": []interface{}{
						map[string]interface{}{
							"name": "msg",
							"type": "string",
						},
					},
				},
			},
		}

		err := pipeline.extractAndStoreActions(network, resultMap)
		if err != nil {
			t.Fatalf("extractAndStoreActions failed: %v", err)
		}

		// Should still have only 1 action (replaced, not added)
		if len(network.Actions) != 1 {
			t.Errorf("Expected 1 action definition (replaced), got %d", len(network.Actions))
		}

		// Verify it has parameters now
		if len(network.Actions[0].Parameters) != 1 {
			t.Errorf("Expected 1 parameter in replaced action, got %d", len(network.Actions[0].Parameters))
		}
	})
}

// TestAnalyzeConstraints tests the analyzeConstraints function
func TestAnalyzeConstraints(t *testing.T) {
	pipeline := NewConstraintPipeline()

	t.Run("detects NOT constraint", func(t *testing.T) {
		constraint := map[string]interface{}{
			"type": "notConstraint",
			"constraint": map[string]interface{}{
				"type":     "comparison",
				"operator": "==",
			},
		}

		isNegation, negatedConstraint, err := pipeline.analyzeConstraints(constraint)
		if err != nil {
			t.Fatalf("analyzeConstraints failed: %v", err)
		}

		if !isNegation {
			t.Error("Expected negation to be detected")
		}

		if negatedConstraint == nil {
			t.Error("Expected negated constraint to be returned")
		}
	})

	t.Run("returns false for regular constraint", func(t *testing.T) {
		constraint := map[string]interface{}{
			"type":     "comparison",
			"operator": "==",
		}

		isNegation, _, err := pipeline.analyzeConstraints(constraint)
		if err != nil {
			t.Fatalf("analyzeConstraints failed: %v", err)
		}

		if isNegation {
			t.Error("Expected no negation")
		}
	})

	t.Run("handles non-map constraint", func(t *testing.T) {
		constraint := "not a map"

		isNegation, _, err := pipeline.analyzeConstraints(constraint)
		if err != nil {
			t.Fatalf("analyzeConstraints failed: %v", err)
		}

		if isNegation {
			t.Error("Expected no negation for non-map")
		}
	})
}

// TestDetectAggregation tests the detectAggregation function
func TestDetectAggregation(t *testing.T) {
	pipeline := NewConstraintPipeline()

	t.Run("detects SUM aggregation", func(t *testing.T) {
		constraintsData := "SUM(amount)"

		hasAgg := pipeline.detectAggregation(constraintsData)

		if !hasAgg {
			t.Error("Expected SUM aggregation to be detected")
		}
	})

	t.Run("detects AVG aggregation", func(t *testing.T) {
		constraintsData := "AVG(price)"

		hasAgg := pipeline.detectAggregation(constraintsData)

		if !hasAgg {
			t.Error("Expected AVG aggregation to be detected")
		}
	})

	t.Run("detects COUNT aggregation", func(t *testing.T) {
		constraintsData := "COUNT(*)"

		hasAgg := pipeline.detectAggregation(constraintsData)

		if !hasAgg {
			t.Error("Expected COUNT aggregation to be detected")
		}
	})

	t.Run("detects MIN aggregation", func(t *testing.T) {
		constraintsData := "MIN(value)"

		hasAgg := pipeline.detectAggregation(constraintsData)

		if !hasAgg {
			t.Error("Expected MIN aggregation to be detected")
		}
	})

	t.Run("detects MAX aggregation", func(t *testing.T) {
		constraintsData := "MAX(value)"

		hasAgg := pipeline.detectAggregation(constraintsData)

		if !hasAgg {
			t.Error("Expected MAX aggregation to be detected")
		}
	})

	t.Run("detects ACCUMULATE", func(t *testing.T) {
		constraintsData := "ACCUMULATE(items)"

		hasAgg := pipeline.detectAggregation(constraintsData)

		if !hasAgg {
			t.Error("Expected ACCUMULATE to be detected")
		}
	})

	t.Run("returns false for non-aggregation", func(t *testing.T) {
		constraintsData := "p.age > 18"

		hasAgg := pipeline.detectAggregation(constraintsData)

		if hasAgg {
			t.Error("Expected no aggregation")
		}
	})

	t.Run("handles empty data", func(t *testing.T) {
		constraintsData := ""

		hasAgg := pipeline.detectAggregation(constraintsData)

		if hasAgg {
			t.Error("Expected no aggregation for empty data")
		}
	})
}

// TestIsExistsConstraint tests the isExistsConstraint function
func TestIsExistsConstraint(t *testing.T) {
	pipeline := NewConstraintPipeline()

	t.Run("detects exists constraint", func(t *testing.T) {
		constraint := map[string]interface{}{
			"type": "existsConstraint",
		}

		result := pipeline.isExistsConstraint(constraint)

		if !result {
			t.Error("Expected exists constraint to be detected")
		}
	})

	t.Run("returns false for non-exists constraint", func(t *testing.T) {
		constraint := map[string]interface{}{
			"type":     "comparison",
			"operator": "==",
		}

		result := pipeline.isExistsConstraint(constraint)

		if result {
			t.Error("Expected non-exists constraint")
		}
	})

	t.Run("handles invalid constraint format", func(t *testing.T) {
		constraint := "not a map"

		result := pipeline.isExistsConstraint(constraint)

		if result {
			t.Error("Expected false for invalid format")
		}
	})

	t.Run("handles map without type field", func(t *testing.T) {
		constraint := map[string]interface{}{
			"operator": "==",
		}

		result := pipeline.isExistsConstraint(constraint)

		if result {
			t.Error("Expected false for map without type field")
		}
	})
}

// TestGetStringField tests the standalone getStringField function
func TestGetStringField(t *testing.T) {
	t.Run("gets string field successfully", func(t *testing.T) {
		data := map[string]interface{}{
			"name": "Alice",
			"age":  30,
		}

		name := getStringField(data, "name", "")

		if name != "Alice" {
			t.Errorf("Expected 'Alice', got '%s'", name)
		}
	})

	t.Run("returns default value for missing field", func(t *testing.T) {
		data := map[string]interface{}{
			"age": 30,
		}

		name := getStringField(data, "name", "default")

		if name != "default" {
			t.Errorf("Expected 'default', got '%s'", name)
		}
	})

	t.Run("returns default value for non-string field", func(t *testing.T) {
		data := map[string]interface{}{
			"age": 30,
		}

		age := getStringField(data, "age", "unknown")

		if age != "unknown" {
			t.Errorf("Expected 'unknown', got '%s'", age)
		}
	})
}
