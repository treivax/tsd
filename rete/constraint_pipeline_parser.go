// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strings"
)

// extractComponents extrait les types et expressions d'un AST parsé
// Note: La gestion des instructions reset est effectuée en amont dans buildNetworkWithResetSemantics
// Cette fonction suppose que le programme reçu a déjà la sémantique reset appliquée si nécessaire
// Retourne (types, expressions, error)
func (cp *ConstraintPipeline) extractComponents(resultMap map[string]interface{}) ([]interface{}, []interface{}, error) {
	// Extraire les types
	typesData, hasTypes := resultMap["types"]
	if !hasTypes {
		return nil, nil, fmt.Errorf("aucun type trouvé dans l'AST")
	}

	types, ok := typesData.([]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("format types invalide: %T", typesData)
	}

	// Extraire les expressions
	expressionsData, hasExpressions := resultMap["expressions"]
	if !hasExpressions {
		return nil, nil, fmt.Errorf("aucune expression trouvée dans l'AST")
	}

	expressions, ok := expressionsData.([]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("format expressions invalide: %T", expressionsData)
	}

	return types, expressions, nil
}

// analyzeConstraints analyse les contraintes pour détecter les négations
// Retourne (isNegation, negatedCondition, error)
func (cp *ConstraintPipeline) analyzeConstraints(constraints interface{}) (bool, interface{}, error) {
	constraintMap, ok := constraints.(map[string]interface{})
	if !ok {
		return false, constraints, nil
	}

	// Détecter contrainte NOT
	if constraintType, exists := constraintMap["type"].(string); exists {
		if constraintType == "notConstraint" {
			// Extraire la contrainte niée
			if negatedConstraint, hasNegated := constraintMap["constraint"]; hasNegated {
				return true, negatedConstraint, nil
			}
		}
	}

	return false, constraints, nil
}

// extractAggregationInfoFromVariables extracts aggregation info from the new multi-pattern syntax
// where aggregation variables are declared in the first pattern block
func (cp *ConstraintPipeline) extractAggregationInfoFromVariables(exprMap map[string]interface{}) (*AggregationInfo, error) {
	aggInfo := &AggregationInfo{}

	// Check for multi-pattern syntax
	patternsData, hasPatterns := exprMap["patterns"]
	if !hasPatterns {
		return nil, fmt.Errorf("no patterns field found")
	}

	patternsList, ok := patternsData.([]interface{})
	if !ok || len(patternsList) < 2 {
		return nil, fmt.Errorf("expected at least 2 pattern blocks for aggregation with join")
	}

	// First pattern block should contain the aggregation variable(s)
	firstPattern, ok := patternsList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("first pattern is not a map")
	}

	varsData, hasVars := firstPattern["variables"]
	if !hasVars {
		return nil, fmt.Errorf("no variables in first pattern")
	}

	varsList, ok := varsData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("variables is not a list")
	}

	// Find the aggregation variable
	for _, varInterface := range varsList {
		varMap, ok := varInterface.(map[string]interface{})
		if !ok {
			continue
		}

		varType, _ := varMap["type"].(string)
		if varType == "aggregationVariable" {
			// Extract aggregation function
			if function, ok := varMap["function"].(string); ok {
				aggInfo.Function = function
			}

			// Extract field being aggregated
			if fieldData, ok := varMap["field"].(map[string]interface{}); ok {
				if fieldObj, ok := fieldData["object"].(string); ok {
					aggInfo.AggVariable = fieldObj
				}
				if fieldName, ok := fieldData["field"].(string); ok {
					aggInfo.Field = fieldName
				}
			}
			break
		}
	}

	// Second pattern block contains the source type
	secondPattern, ok := patternsList[1].(map[string]interface{})
	if ok {
		if varsData2, hasVars2 := secondPattern["variables"]; hasVars2 {
			if varsList2, ok := varsData2.([]interface{}); ok && len(varsList2) > 0 {
				if varMap2, ok := varsList2[0].(map[string]interface{}); ok {
					if aggType, ok := varMap2["dataType"].(string); ok {
						aggInfo.AggType = aggType
					}
				}
			}
		}
	}

	// Extract join condition and threshold from constraints
	if constraintsData, hasConstraints := exprMap["constraints"]; hasConstraints {
		if constraintMap, ok := constraintsData.(map[string]interface{}); ok {
			aggInfo.JoinCondition = constraintMap

			// Get list of aggregation variable names from first pattern
			aggVarNames := cp.getAggregationVariableNames(exprMap)

			// Separate join conditions and threshold conditions
			joinConditions, thresholdConditions := cp.separateAggregationConstraints(constraintMap, aggVarNames)

			// Extract join fields from join conditions
			if joinConditions != nil {
				if joinConditions["type"] == "comparison" {
					// Left side
					if leftData, ok := joinConditions["left"].(map[string]interface{}); ok {
						if leftData["type"] == "fieldAccess" {
							if joinField, ok := leftData["field"].(string); ok {
								aggInfo.JoinField = joinField
							}
						}
					}

					// Right side
					if rightData, ok := joinConditions["right"].(map[string]interface{}); ok {
						if rightData["type"] == "fieldAccess" {
							if mainField, ok := rightData["field"].(string); ok {
								aggInfo.MainField = mainField
							}
						}
					}
				}
			}

			// Extract threshold from threshold conditions
			if len(thresholdConditions) > 0 {
				// Use the first threshold condition found
				threshold := thresholdConditions[0]
				if operator, ok := threshold["operator"].(string); ok {
					aggInfo.Operator = operator
				}
				if rightData, ok := threshold["right"].(map[string]interface{}); ok {
					if value, ok := rightData["value"].(float64); ok {
						aggInfo.Threshold = value
					}
				}
			} else {
				// No threshold - always fire (use >= 0 as default)
				aggInfo.Operator = ">="
				aggInfo.Threshold = 0
			}
		}
	} else {
		// No constraints - set default threshold
		aggInfo.Operator = ">="
		aggInfo.Threshold = 0
	}

	return aggInfo, nil
}

// extractMultiSourceAggregationInfo extracts aggregation info for multi-source aggregations
// This supports aggregating over multiple joined patterns
func (cp *ConstraintPipeline) extractMultiSourceAggregationInfo(exprMap map[string]interface{}) (*AggregationInfo, error) {
	aggInfo := &AggregationInfo{}
	aggInfo.AggregationVars = []AggregationVariable{}
	aggInfo.SourcePatterns = []SourcePattern{}
	aggInfo.JoinConditions = []JoinCondition{}

	// Check for multi-pattern syntax
	patternsData, hasPatterns := exprMap["patterns"]
	if !hasPatterns {
		return nil, fmt.Errorf("no patterns field found")
	}

	patternsList, ok := patternsData.([]interface{})
	if !ok || len(patternsList) < 2 {
		return nil, fmt.Errorf("expected at least 2 pattern blocks for aggregation")
	}

	// First pattern block contains main variable and aggregation variable(s)
	firstPattern, ok := patternsList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("first pattern is not a map")
	}

	varsData, hasVars := firstPattern["variables"]
	if !hasVars {
		return nil, fmt.Errorf("no variables in first pattern")
	}

	varsList, ok := varsData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("variables is not a list")
	}

	// Extract main variable and aggregation variables from first pattern
	for _, varInterface := range varsList {
		varMap, ok := varInterface.(map[string]interface{})
		if !ok {
			continue
		}

		varType, _ := varMap["type"].(string)
		if varType == "aggregationVariable" {
			// This is an aggregation variable
			aggVar := AggregationVariable{}

			if name, ok := varMap["name"].(string); ok {
				aggVar.Name = name
			}

			if function, ok := varMap["function"].(string); ok {
				aggVar.Function = function
			}

			// Extract field being aggregated
			if fieldData, ok := varMap["field"].(map[string]interface{}); ok {
				if fieldObj, ok := fieldData["object"].(string); ok {
					aggVar.SourceVar = fieldObj
				}
				if fieldName, ok := fieldData["field"].(string); ok {
					aggVar.Field = fieldName
				}
			}

			aggInfo.AggregationVars = append(aggInfo.AggregationVars, aggVar)
		} else {
			// This is the main variable
			if name, ok := varMap["name"].(string); ok {
				aggInfo.MainVariable = name
			}
			if dataType, ok := varMap["dataType"].(string); ok {
				aggInfo.MainType = dataType
			}
		}
	}

	// Extract source patterns from remaining pattern blocks
	for i := 1; i < len(patternsList); i++ {
		pattern, ok := patternsList[i].(map[string]interface{})
		if !ok {
			continue
		}

		if varsData, hasVars := pattern["variables"]; hasVars {
			if varsList, ok := varsData.([]interface{}); ok && len(varsList) > 0 {
				if varMap, ok := varsList[0].(map[string]interface{}); ok {
					sourcePattern := SourcePattern{}

					if varName, ok := varMap["name"].(string); ok {
						sourcePattern.Variable = varName
					}
					if varType, ok := varMap["dataType"].(string); ok {
						sourcePattern.Type = varType
					}

					aggInfo.SourcePatterns = append(aggInfo.SourcePatterns, sourcePattern)

					// Also update legacy fields if this is the first/primary aggregation source
					if i == 1 && len(aggInfo.AggregationVars) > 0 {
						aggInfo.AggVariable = sourcePattern.Variable
						aggInfo.AggType = sourcePattern.Type
						aggInfo.Function = aggInfo.AggregationVars[0].Function
						aggInfo.Field = aggInfo.AggregationVars[0].Field
					}
				}
			}
		}
	}

	// Extract join conditions and thresholds from constraints
	if constraintsData, hasConstraints := exprMap["constraints"]; hasConstraints {
		if constraintMap, ok := constraintsData.(map[string]interface{}); ok {
			aggInfo.JoinCondition = constraintMap

			// Get list of aggregation variable names
			aggVarNames := cp.getAggregationVariableNames(exprMap)

			// Separate join conditions and threshold conditions
			joinConditionsMap, thresholdConditions := cp.separateAggregationConstraints(constraintMap, aggVarNames)

			// Extract all join conditions
			cp.extractJoinConditionsRecursive(constraintMap, aggVarNames, &aggInfo.JoinConditions)

			// Extract join fields from the first join condition (for backward compatibility)
			if joinConditionsMap != nil && joinConditionsMap["type"] == "comparison" {
				if leftData, ok := joinConditionsMap["left"].(map[string]interface{}); ok {
					if leftData["type"] == "fieldAccess" {
						if leftObj, ok := leftData["object"].(string); ok {
							if field, ok := leftData["field"].(string); ok {
								// Determine if this is the main or agg side
								if leftObj == aggInfo.MainVariable {
									aggInfo.MainField = field
								} else {
									aggInfo.JoinField = field
								}
							}
						}
					}
				}

				if rightData, ok := joinConditionsMap["right"].(map[string]interface{}); ok {
					if rightData["type"] == "fieldAccess" {
						if rightObj, ok := rightData["object"].(string); ok {
							if field, ok := rightData["field"].(string); ok {
								if rightObj == aggInfo.MainVariable {
									aggInfo.MainField = field
								} else {
									aggInfo.JoinField = field
								}
							}
						}
					}
				}
			}

			// Extract thresholds and apply to aggregation variables
			for _, threshold := range thresholdConditions {
				if leftData, ok := threshold["left"].(map[string]interface{}); ok {
					if leftData["type"] == "variable" {
						if aggVarName, ok := leftData["name"].(string); ok {
							// Find the matching aggregation variable
							for i := range aggInfo.AggregationVars {
								if aggInfo.AggregationVars[i].Name == aggVarName {
									if operator, ok := threshold["operator"].(string); ok {
										aggInfo.AggregationVars[i].Operator = operator
									}
									if rightData, ok := threshold["right"].(map[string]interface{}); ok {
										if value, ok := rightData["value"].(float64); ok {
											aggInfo.AggregationVars[i].Threshold = value
										}
									}
									break
								}
							}
						}
					}
				}
			}

			// Set default threshold for first aggregation variable (backward compatibility)
			if len(aggInfo.AggregationVars) > 0 {
				if aggInfo.AggregationVars[0].Operator != "" {
					aggInfo.Operator = aggInfo.AggregationVars[0].Operator
					aggInfo.Threshold = aggInfo.AggregationVars[0].Threshold
				} else {
					aggInfo.Operator = ">="
					aggInfo.Threshold = 0
					aggInfo.AggregationVars[0].Operator = ">="
					aggInfo.AggregationVars[0].Threshold = 0
				}
			}
		}
	} else {
		// No constraints - set default threshold
		aggInfo.Operator = ">="
		aggInfo.Threshold = 0
		for i := range aggInfo.AggregationVars {
			aggInfo.AggregationVars[i].Operator = ">="
			aggInfo.AggregationVars[i].Threshold = 0
		}
	}

	return aggInfo, nil
}

// extractJoinConditionsRecursive recursively extracts all join conditions from constraint tree
func (cp *ConstraintPipeline) extractJoinConditionsRecursive(constraints map[string]interface{}, aggVarNames map[string]bool, joinConditions *[]JoinCondition) {
	constraintType, _ := constraints["type"].(string)

	if constraintType == "comparison" {
		// Check if this is a join condition (not a threshold)
		if !cp.isThresholdCondition(constraints, aggVarNames) {
			// Extract join condition details
			joinCond := JoinCondition{}

			if operator, ok := constraints["operator"].(string); ok {
				joinCond.Operator = operator
			}

			if leftData, ok := constraints["left"].(map[string]interface{}); ok {
				if leftData["type"] == "fieldAccess" {
					if obj, ok := leftData["object"].(string); ok {
						joinCond.LeftVar = obj
					}
					if field, ok := leftData["field"].(string); ok {
						joinCond.LeftField = field
					}
				}
			}

			if rightData, ok := constraints["right"].(map[string]interface{}); ok {
				if rightData["type"] == "fieldAccess" {
					if obj, ok := rightData["object"].(string); ok {
						joinCond.RightVar = obj
					}
					if field, ok := rightData["field"].(string); ok {
						joinCond.RightField = field
					}
				}
			}

			*joinConditions = append(*joinConditions, joinCond)
		}
	} else if constraintType == "logicalExpr" {
		// Recursively process left side
		if leftData, ok := constraints["left"].(map[string]interface{}); ok {
			cp.extractJoinConditionsRecursive(leftData, aggVarNames, joinConditions)
		}

		// Recursively process operations
		if ops, ok := constraints["operations"].([]interface{}); ok {
			for _, opInterface := range ops {
				if opMap, ok := opInterface.(map[string]interface{}); ok {
					if rightData, ok := opMap["right"].(map[string]interface{}); ok {
						cp.extractJoinConditionsRecursive(rightData, aggVarNames, joinConditions)
					}
				}
			}
		}
	}
}

// getAggregationVariableNames extracts the names of all aggregation variables from patterns
func (cp *ConstraintPipeline) getAggregationVariableNames(exprMap map[string]interface{}) map[string]bool {
	aggVarNames := make(map[string]bool)

	if patternsData, hasPatterns := exprMap["patterns"]; hasPatterns {
		if patternsList, ok := patternsData.([]interface{}); ok {
			for _, patternInterface := range patternsList {
				if patternMap, ok := patternInterface.(map[string]interface{}); ok {
					if varsData, hasVars := patternMap["variables"]; hasVars {
						if varsList, ok := varsData.([]interface{}); ok {
							for _, varInterface := range varsList {
								if varMap, ok := varInterface.(map[string]interface{}); ok {
									if varType, ok := varMap["type"].(string); ok && varType == "aggregationVariable" {
										if name, ok := varMap["name"].(string); ok {
											aggVarNames[name] = true
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return aggVarNames
}

// separateAggregationConstraints separates join conditions from threshold conditions
// Returns (joinConditions, thresholdConditions)
func (cp *ConstraintPipeline) separateAggregationConstraints(constraints map[string]interface{}, aggVarNames map[string]bool) (map[string]interface{}, []map[string]interface{}) {
	var joinConditions map[string]interface{}
	var thresholdConditions []map[string]interface{}

	constraintType, _ := constraints["type"].(string)

	if constraintType == "comparison" {
		// Single comparison - check if it's a threshold or join
		if cp.isThresholdCondition(constraints, aggVarNames) {
			thresholdConditions = append(thresholdConditions, constraints)
		} else {
			joinConditions = constraints
		}
	} else if constraintType == "logicalExpr" {
		// Logical expression (AND/OR) - separate conditions
		leftData, _ := constraints["left"]

		// Handle operations field - it might be []interface{} or []map[string]interface{}
		var operations []interface{}
		if ops, ok := constraints["operations"].([]interface{}); ok {
			operations = ops
		} else if ops, ok := constraints["operations"].([]map[string]interface{}); ok {
			// Convert to []interface{}
			for _, op := range ops {
				operations = append(operations, op)
			}
		}

		// Check left condition
		if leftMap, ok := leftData.(map[string]interface{}); ok {
			if cp.isThresholdCondition(leftMap, aggVarNames) {
				thresholdConditions = append(thresholdConditions, leftMap)
			} else {
				joinConditions = leftMap
			}
		}

		// Check operations
		for _, opInterface := range operations {
			if opMap, ok := opInterface.(map[string]interface{}); ok {
				if rightData, ok := opMap["right"].(map[string]interface{}); ok {
					if cp.isThresholdCondition(rightData, aggVarNames) {
						thresholdConditions = append(thresholdConditions, rightData)
					} else if joinConditions == nil {
						joinConditions = rightData
					}
				}
			}
		}
	}

	return joinConditions, thresholdConditions
}

// isThresholdCondition checks if a comparison references an aggregation variable
func (cp *ConstraintPipeline) isThresholdCondition(condition map[string]interface{}, aggVarNames map[string]bool) bool {
	condType, _ := condition["type"].(string)
	if condType != "comparison" {
		return false
	}

	// Check if left side is an aggregation variable
	if leftData, ok := condition["left"].(map[string]interface{}); ok {
		if leftData["type"] == "variable" {
			if varName, ok := leftData["name"].(string); ok {
				return aggVarNames[varName]
			}
		}
	}

	// Check if right side is an aggregation variable
	if rightData, ok := condition["right"].(map[string]interface{}); ok {
		if rightData["type"] == "variable" {
			if varName, ok := rightData["name"].(string); ok {
				return aggVarNames[varName]
			}
		}
	}

	return false
}

// extractAggregationInfo extrait les informations d'agrégation d'une contrainte
func (cp *ConstraintPipeline) extractAggregationInfo(constraintsData interface{}) (*AggregationInfo, error) {
	constraintMap, ok := constraintsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("constraints n'est pas un map: %T", constraintsData)
	}

	aggInfo := &AggregationInfo{}

	// Extraire la fonction d'agrégation (AVG, SUM, COUNT, etc.)
	if function, ok := constraintMap["function"].(string); ok {
		aggInfo.Function = function
	} else {
		return nil, fmt.Errorf("fonction d'agrégation non trouvée")
	}

	// Extraire l'opérateur de comparaison
	if operator, ok := constraintMap["operator"].(string); ok {
		aggInfo.Operator = operator
	} else {
		return nil, fmt.Errorf("opérateur de comparaison non trouvé")
	}

	// Extraire le seuil (threshold) depuis constraintMap["threshold"]["value"]
	if thresholdData, ok := constraintMap["threshold"].(map[string]interface{}); ok {
		if threshold, ok := thresholdData["value"].(float64); ok {
			aggInfo.Threshold = threshold
		} else if thresholdInt, ok := thresholdData["value"].(int); ok {
			aggInfo.Threshold = float64(thresholdInt)
		} else {
			return nil, fmt.Errorf("valeur de seuil non trouvée ou invalide")
		}
	} else {
		return nil, fmt.Errorf("seuil manquant")
	}

	// Extraire la condition de jointure complète
	if joinCond, ok := constraintMap["join"]; ok {
		aggInfo.JoinCondition = joinCond
	}

	// Extraire la variable à agréger depuis constraintMap["variable"]
	if variableData, ok := constraintMap["variable"].(map[string]interface{}); ok {
		if aggVar, ok := variableData["name"].(string); ok {
			aggInfo.AggVariable = aggVar
		}
		if aggType, ok := variableData["dataType"].(string); ok {
			aggInfo.AggType = aggType
		}
	}

	// Extraire le champ à agréger
	if field, ok := constraintMap["field"].(string); ok {
		aggInfo.Field = field
	}

	// Extraire les informations de jointure depuis la condition
	if conditionData, ok := constraintMap["condition"].(map[string]interface{}); ok {
		aggInfo.JoinCondition = conditionData

		// Extraire les champs de jointure depuis la condition de type comparison
		if condType, ok := conditionData["type"].(string); ok && condType == "comparison" {
			// Left side: p.employee_id
			if leftData, ok := conditionData["left"].(map[string]interface{}); ok {
				if leftType, ok := leftData["type"].(string); ok && leftType == "fieldAccess" {
					if joinField, ok := leftData["field"].(string); ok {
						aggInfo.JoinField = joinField
					}
				}
			}

			// Right side: e.id
			if rightData, ok := conditionData["right"].(map[string]interface{}); ok {
				if rightType, ok := rightData["type"].(string); ok && rightType == "fieldAccess" {
					if mainField, ok := rightData["field"].(string); ok {
						aggInfo.MainField = mainField
					}
				}
			}
		}
	}

	return aggInfo, nil
}

// extractVariablesFromExpression extrait les variables d'une expression
// Retourne (variables, variableNames, variableTypes)
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

// detectAggregation détecte si une contrainte contient une agrégation
func (cp *ConstraintPipeline) detectAggregation(constraintsData interface{}) bool {
	if constraintStr := fmt.Sprintf("%v", constraintsData); constraintStr != "" {
		return strings.Contains(constraintStr, "AVG") ||
			strings.Contains(constraintStr, "SUM") ||
			strings.Contains(constraintStr, "COUNT") ||
			strings.Contains(constraintStr, "MIN") ||
			strings.Contains(constraintStr, "MAX") ||
			strings.Contains(constraintStr, "ACCUMULATE")
	}
	return false
}

// hasAggregationVariables checks if any variables in the expression are aggregation variables
func (cp *ConstraintPipeline) hasAggregationVariables(exprMap map[string]interface{}) bool {
	// Check new multi-pattern syntax
	if patternsData, hasPatterns := exprMap["patterns"]; hasPatterns {
		if patternsList, ok := patternsData.([]interface{}); ok {
			for _, patternInterface := range patternsList {
				if patternMap, ok := patternInterface.(map[string]interface{}); ok {
					if varsData, hasVars := patternMap["variables"]; hasVars {
						if varsList, ok := varsData.([]interface{}); ok {
							for _, varInterface := range varsList {
								if varMap, ok := varInterface.(map[string]interface{}); ok {
									if varType, ok := varMap["type"].(string); ok && varType == "aggregationVariable" {
										return true
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Check old single-pattern syntax for backward compatibility
	if setData, hasSet := exprMap["set"]; hasSet {
		if setMap, ok := setData.(map[string]interface{}); ok {
			if varsData, hasVars := setMap["variables"]; hasVars {
				if varsList, ok := varsData.([]interface{}); ok {
					for _, varInterface := range varsList {
						if varMap, ok := varInterface.(map[string]interface{}); ok {
							if varType, ok := varMap["type"].(string); ok && varType == "aggregationVariable" {
								return true
							}
						}
					}
				}
			}
		}
	}

	return false
}

// isExistsConstraint vérifie si une contrainte est de type EXISTS
func (cp *ConstraintPipeline) isExistsConstraint(constraintsData interface{}) bool {
	if constraintMap, ok := constraintsData.(map[string]interface{}); ok {
		if constraintType, exists := constraintMap["type"].(string); exists && constraintType == "existsConstraint" {
			return true
		}
	}
	return false
}

// getStringField extrait un champ string d'une map avec une valeur par défaut
func getStringField(m map[string]interface{}, key, defaultValue string) string {
	if value, ok := m[key].(string); ok {
		return value
	}
	return defaultValue
}
