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

	// Extract join condition from constraints
	if constraintsData, hasConstraints := exprMap["constraints"]; hasConstraints {
		if constraintMap, ok := constraintsData.(map[string]interface{}); ok {
			aggInfo.JoinCondition = constraintMap

			// Extract join fields from comparison
			if constraintMap["type"] == "comparison" {
				// Left side
				if leftData, ok := constraintMap["left"].(map[string]interface{}); ok {
					if leftData["type"] == "fieldAccess" {
						if joinField, ok := leftData["field"].(string); ok {
							aggInfo.JoinField = joinField
						}
					}
				}

				// Right side
				if rightData, ok := constraintMap["right"].(map[string]interface{}); ok {
					if rightData["type"] == "fieldAccess" {
						if mainField, ok := rightData["field"].(string); ok {
							aggInfo.MainField = mainField
						}
					}
				}
			}
		}
	}

	// For aggregation with join syntax, we don't have a threshold comparison
	// Set default values
	aggInfo.Operator = ">="
	aggInfo.Threshold = 0

	return aggInfo, nil
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
