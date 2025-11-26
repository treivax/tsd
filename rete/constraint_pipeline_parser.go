// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strings"
)

// extractComponents extrait les types et expressions d'un AST parsé
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
