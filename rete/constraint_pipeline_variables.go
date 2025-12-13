// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// extractVariablesFromExpression extrait les variables d'une expression
// Retourne (variables, variableNames, variableTypes)
//
// Cette fonction supporte deux formats d'expressions :
// 1. Format multi-pattern moderne (field "patterns")
// 2. Format single-pattern legacy (field "set")
func (cp *ConstraintPipeline) extractVariablesFromExpression(exprMap map[string]interface{}) ([]map[string]interface{}, []string, []string) {
	// Check for multi-pattern aggregation syntax (new format with "patterns" field)
	if patternsData, hasPatterns := exprMap["patterns"]; hasPatterns {
		return cp.extractVariablesFromPatterns(patternsData)
	}

	// Original single-pattern syntax (backward compatibility)
	if setData, hasSet := exprMap["set"]; hasSet {
		return cp.extractVariablesFromSet(setData)
	}

	return []map[string]interface{}{}, []string{}, []string{}
}

// extractVariablesFromPatterns extrait les variables du format multi-pattern
func (cp *ConstraintPipeline) extractVariablesFromPatterns(patternsData interface{}) ([]map[string]interface{}, []string, []string) {
	variables := []map[string]interface{}{}
	variableNames := []string{}
	variableTypes := []string{}

	patternsList, ok := patternsData.([]interface{})
	if !ok {
		return variables, variableNames, variableTypes
	}

	// Process all pattern blocks
	for _, patternInterface := range patternsList {
		patternMap, ok := patternInterface.(map[string]interface{})
		if !ok {
			continue
		}

		varsData, hasVars := patternMap["variables"]
		if !hasVars {
			continue
		}

		varsList, ok := varsData.([]interface{})
		if !ok {
			continue
		}

		for _, varInterface := range varsList {
			varMap, ok := varInterface.(map[string]interface{})
			if !ok {
				continue
			}

			variables = append(variables, varMap)

			if name, ok := varMap["name"].(string); ok {
				variableNames = append(variableNames, name)
			}

			varType := cp.extractVariableType(varMap)
			variableTypes = append(variableTypes, varType)
		}
	}

	return variables, variableNames, variableTypes
}

// extractVariablesFromSet extrait les variables du format single-pattern legacy
func (cp *ConstraintPipeline) extractVariablesFromSet(setData interface{}) ([]map[string]interface{}, []string, []string) {
	variables := []map[string]interface{}{}
	variableNames := []string{}
	variableTypes := []string{}

	setMap, ok := setData.(map[string]interface{})
	if !ok {
		return variables, variableNames, variableTypes
	}

	varsData, hasVars := setMap["variables"]
	if !hasVars {
		return variables, variableNames, variableTypes
	}

	varsList, ok := varsData.([]interface{})
	if !ok || len(varsList) == 0 {
		return variables, variableNames, variableTypes
	}

	// Extraire toutes les variables
	for _, varInterface := range varsList {
		varMap, ok := varInterface.(map[string]interface{})
		if !ok {
			continue
		}

		variables = append(variables, varMap)

		if name, ok := varMap["name"].(string); ok {
			variableNames = append(variableNames, name)
		}

		varType := cp.extractVariableType(varMap)
		variableTypes = append(variableTypes, varType)
	}

	return variables, variableNames, variableTypes
}

// extractVariableType extrait le type d'une variable
func (cp *ConstraintPipeline) extractVariableType(varMap map[string]interface{}) string {
	// Handle aggregation variables specially
	if varTypeStr, ok := varMap["type"].(string); ok && varTypeStr == "aggregationVariable" {
		// For aggregation variables, use a placeholder type
		// The actual aggregation result type will be determined at runtime
		return "AggregationResult"
	}

	// Try dataType field first
	if dataType, ok := varMap["dataType"].(string); ok {
		return dataType
	}

	// Fall back to type field
	if typeField, ok := varMap["type"].(string); ok {
		return typeField
	}

	return ""
}
