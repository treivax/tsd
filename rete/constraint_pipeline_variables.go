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
	variables := []map[string]interface{}{}
	variableNames := []string{}
	variableTypes := []string{}

	// Check for multi-pattern aggregation syntax (new format with "patterns" field)
	if patternsData, hasPatterns := exprMap["patterns"]; hasPatterns {
		if patternsList, ok := patternsData.([]interface{}); ok {
			// Process all pattern blocks
			for _, patternInterface := range patternsList {
				if patternMap, ok := patternInterface.(map[string]interface{}); ok {
					if varsData, hasVars := patternMap["variables"]; hasVars {
						if varsList, ok := varsData.([]interface{}); ok {
							for _, varInterface := range varsList {
								if varMap, ok := varInterface.(map[string]interface{}); ok {
									variables = append(variables, varMap)

									if name, ok := varMap["name"].(string); ok {
										variableNames = append(variableNames, name)
									}

									// Extract variable type - handle both regular and aggregation variables
									var varType string
									varTypeStr := varMap["type"].(string)
									if varTypeStr == "aggregationVariable" {
										// For aggregation variables, use a placeholder type
										// The actual aggregation result type will be determined at runtime
										varType = "AggregationResult"
									} else if dataType, ok := varMap["dataType"].(string); ok {
										varType = dataType
									} else if typeField, ok := varMap["type"].(string); ok {
										varType = typeField
									}
									variableTypes = append(variableTypes, varType)
								}
							}
						}
					}
				}
			}
		}
		return variables, variableNames, variableTypes
	}

	// Original single-pattern syntax (backward compatibility)
	if setData, hasSet := exprMap["set"]; hasSet {
		if setMap, ok := setData.(map[string]interface{}); ok {
			if varsData, hasVars := setMap["variables"]; hasVars {
				if varsList, ok := varsData.([]interface{}); ok && len(varsList) > 0 {
					// Extraire toutes les variables
					for _, varInterface := range varsList {
						if varMap, ok := varInterface.(map[string]interface{}); ok {
							variables = append(variables, varMap)

							if name, ok := varMap["name"].(string); ok {
								variableNames = append(variableNames, name)
							}

							// Extraire le type de la variable
							var varType string
							if dataType, ok := varMap["dataType"].(string); ok {
								varType = dataType
							} else if typeField, ok := varMap["type"].(string); ok {
								varType = typeField
							}
							variableTypes = append(variableTypes, varType)
						}
					}
				}
			}
		}
	}

	return variables, variableNames, variableTypes
}
