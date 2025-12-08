// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import "fmt"

// aggregation_helpers.go contient des fonctions utilitaires réutilisables
// pour l'extraction et la manipulation des informations d'agrégation.

// Constants for aggregation types and default values
const (
	// Aggregation variable type identifier
	AggregationVariableType = "aggregationVariable"

	// Field access type identifier
	FieldAccessType = "fieldAccess"

	// Function call type identifiers
	FunctionCallType    = "functionCall"
	AggregationCallType = "aggregationCall"

	// Comparison type identifier
	ComparisonType = "comparison"

	// Default threshold values when no threshold is specified
	DefaultThresholdOperator = ">="
	DefaultThresholdValue    = 0.0
)

// getFirstPattern extracts and validates the first pattern from patterns list
func getFirstPattern(exprMap map[string]interface{}) (map[string]interface{}, error) {
	patternsData, hasPatterns := exprMap["patterns"]
	if !hasPatterns {
		return nil, fmt.Errorf("no patterns field found")
	}

	patternsList, ok := patternsData.([]interface{})
	if !ok || len(patternsList) < 1 {
		return nil, fmt.Errorf("patterns is not a list or is empty")
	}

	firstPattern, ok := patternsList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("first pattern is not a map")
	}

	return firstPattern, nil
}

// getSecondPattern extracts and validates the second pattern from patterns list
func getSecondPattern(exprMap map[string]interface{}) (map[string]interface{}, error) {
	patternsData, hasPatterns := exprMap["patterns"]
	if !hasPatterns {
		return nil, fmt.Errorf("no patterns field found")
	}

	patternsList, ok := patternsData.([]interface{})
	if !ok || len(patternsList) < 2 {
		return nil, fmt.Errorf("expected at least 2 pattern blocks for aggregation with join")
	}

	secondPattern, ok := patternsList[1].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("second pattern is not a map")
	}

	return secondPattern, nil
}

// getVariablesList extracts the variables list from a pattern
func getVariablesList(pattern map[string]interface{}) ([]interface{}, error) {
	varsData, hasVars := pattern["variables"]
	if !hasVars {
		return nil, fmt.Errorf("no variables in pattern")
	}

	varsList, ok := varsData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("variables is not a list")
	}

	return varsList, nil
}

// findAggregationVariable finds the first aggregation variable in a variables list
func findAggregationVariable(varsList []interface{}) (map[string]interface{}, bool) {
	for _, varInterface := range varsList {
		varMap, ok := varInterface.(map[string]interface{})
		if !ok {
			continue
		}

		varType, _ := varMap["type"].(string)
		if varType == AggregationVariableType {
			return varMap, true
		}
	}
	return nil, false
}

// extractStringField safely extracts a string field from a map
func extractStringField(m map[string]interface{}, key string) (string, bool) {
	if value, ok := m[key].(string); ok {
		return value, true
	}
	return "", false
}

// extractMapField safely extracts a map field from a map
func extractMapField(m map[string]interface{}, key string) (map[string]interface{}, bool) {
	if value, ok := m[key].(map[string]interface{}); ok {
		return value, true
	}
	return nil, false
}

// extractListField safely extracts a list field from a map
func extractListField(m map[string]interface{}, key string) ([]interface{}, bool) {
	if value, ok := m[key].([]interface{}); ok {
		return value, true
	}
	return nil, false
}

// extractFloat64Field safely extracts a float64 field from a map
func extractFloat64Field(m map[string]interface{}, key string) (float64, bool) {
	if value, ok := m[key].(float64); ok {
		return value, true
	}
	return 0, false
}

// isFieldAccessType checks if a map has type "fieldAccess"
func isFieldAccessType(m map[string]interface{}) bool {
	typeStr, ok := extractStringField(m, "type")
	return ok && typeStr == FieldAccessType
}

// isComparisonType checks if a map has type "comparison"
func isComparisonType(m map[string]interface{}) bool {
	typeStr, ok := extractStringField(m, "type")
	return ok && typeStr == ComparisonType
}

// isFunctionCallType checks if a map is a function call or aggregation call
func isFunctionCallType(m map[string]interface{}) bool {
	typeStr, ok := extractStringField(m, "type")
	return ok && (typeStr == FunctionCallType || typeStr == AggregationCallType)
}
