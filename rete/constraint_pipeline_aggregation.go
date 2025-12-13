// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strings"
)

// Ce fichier contient toutes les fonctions de parsing d'agrégation.
// Il supporte trois formats d'agrégation :
// 1. Format legacy : extractAggregationInfo
// 2. Format multi-pattern simple : extractAggregationInfoFromVariables
// 3. Format multi-source avancé : extractMultiSourceAggregationInfo

// extractAggregationInfoFromVariables extracts aggregation info from the new multi-pattern syntax
// where aggregation variables are declared in the first pattern block.
// Cette fonction agit comme orchestrateur, déléguant l'extraction à des fonctions spécialisées
// pour réduire la complexité et améliorer la testabilité.
func (cp *ConstraintPipeline) extractAggregationInfoFromVariables(exprMap map[string]interface{}) (*AggregationInfo, error) {
	aggInfo := &AggregationInfo{}

	// Étape 1: Parser et valider la structure de base
	_, varsList, err := cp.parseAggregationExpression(exprMap)
	if err != nil {
		return nil, err
	}

	// Étape 2: Trouver la variable d'agrégation
	aggVar, found := findAggregationVariable(varsList)
	if !found {
		return nil, fmt.Errorf("no aggregation variable found in first pattern")
	}

	// Étape 3: Extraire la fonction d'agrégation (AVG, SUM, COUNT, etc.)
	function, err := cp.extractAggregationFunction(aggVar)
	if err != nil {
		return nil, err
	}
	aggInfo.Function = function

	// Étape 4: Extraire le champ agrégé et la variable source
	aggVariable, field, err := cp.extractAggregationField(aggVar)
	if err != nil {
		return nil, err
	}
	aggInfo.AggVariable = aggVariable
	aggInfo.Field = field

	// Étape 5: Extraire le type source depuis le second pattern
	aggType, err := cp.extractSourceType(exprMap)
	if err != nil {
		// Type source optionnel, continuer sans erreur
		aggType = ""
	}
	aggInfo.AggType = aggType

	// Étape 6: Extraire les conditions de jointure et seuil
	if constraintsData, hasConstraints := exprMap["constraints"]; hasConstraints {
		if constraintMap, ok := constraintsData.(map[string]interface{}); ok {
			aggInfo.JoinCondition = constraintMap

			// Obtenir les noms des variables d'agrégation
			aggVarNames := cp.getAggregationVariableNames(exprMap)

			// Séparer les conditions de jointure et de seuil
			joinConditions, thresholdConditions := cp.separateAggregationConstraints(constraintMap, aggVarNames)

			// Étape 7: Extraire les champs de jointure
			if joinConditions != nil {
				joinField, mainField := cp.extractJoinFields(joinConditions)
				aggInfo.JoinField = joinField
				aggInfo.MainField = mainField
			}

			// Étape 8: Extraire les conditions de seuil
			operator, threshold := cp.extractThresholdConditions(thresholdConditions)
			aggInfo.Operator = operator
			aggInfo.Threshold = threshold
		}
	} else {
		// Pas de contraintes - utiliser les valeurs par défaut
		aggInfo.Operator = DefaultThresholdOperator
		aggInfo.Threshold = DefaultThresholdValue
	}

	return aggInfo, nil
}

// extractMultiSourceAggregationInfo extracts aggregation info for multi-source aggregations
// This supports aggregating over multiple joined patterns
func (cp *ConstraintPipeline) extractMultiSourceAggregationInfo(exprMap map[string]interface{}) (*AggregationInfo, error) {
	// Delegate to orchestrated version
	return cp.extractMultiSourceAggregationInfoOrchestrated(exprMap)
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

// extractAggregationInfo extrait les informations d'agrégation d'une contrainte (format legacy)
func (cp *ConstraintPipeline) extractAggregationInfo(constraintsData interface{}) (*AggregationInfo, error) {
	constraintMap, ok := constraintsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("constraints n'est pas un map: %T", constraintsData)
	}

	aggInfo := &AggregationInfo{}

	// Extraire les composants de base
	if err := cp.extractBasicAggregationFields(constraintMap, aggInfo); err != nil {
		return nil, err
	}

	// Extraire les champs de jointure
	cp.extractAggregationJoinFields(constraintMap, aggInfo)

	return aggInfo, nil
}

// extractBasicAggregationFields extrait les champs de base d'une agrégation
func (cp *ConstraintPipeline) extractBasicAggregationFields(constraintMap map[string]interface{}, aggInfo *AggregationInfo) error {
	// Fonction d'agrégation
	if function, ok := constraintMap["function"].(string); ok {
		aggInfo.Function = function
	} else {
		return fmt.Errorf("fonction d'agrégation non trouvée")
	}

	// Opérateur de comparaison
	if operator, ok := constraintMap["operator"].(string); ok {
		aggInfo.Operator = operator
	} else {
		return fmt.Errorf("opérateur de comparaison non trouvé")
	}

	// Seuil
	if err := cp.extractAggregationThreshold(constraintMap, aggInfo); err != nil {
		return err
	}

	// Variable à agréger
	cp.extractAggregationVariable(constraintMap, aggInfo)

	// Champ à agréger
	if field, ok := constraintMap["field"].(string); ok {
		aggInfo.Field = field
	}

	return nil
}

// extractAggregationThreshold extrait le seuil d'agrégation
func (cp *ConstraintPipeline) extractAggregationThreshold(constraintMap map[string]interface{}, aggInfo *AggregationInfo) error {
	thresholdData, ok := constraintMap["threshold"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("seuil manquant")
	}

	if threshold, ok := thresholdData["value"].(float64); ok {
		aggInfo.Threshold = threshold
		return nil
	}

	if thresholdInt, ok := thresholdData["value"].(int); ok {
		aggInfo.Threshold = float64(thresholdInt)
		return nil
	}

	return fmt.Errorf("valeur de seuil non trouvée ou invalide")
}

// extractAggregationVariable extrait la variable à agréger
func (cp *ConstraintPipeline) extractAggregationVariable(constraintMap map[string]interface{}, aggInfo *AggregationInfo) {
	variableData, ok := constraintMap["variable"].(map[string]interface{})
	if !ok {
		return
	}

	if aggVar, ok := variableData["name"].(string); ok {
		aggInfo.AggVariable = aggVar
	}

	if aggType, ok := variableData["dataType"].(string); ok {
		aggInfo.AggType = aggType
	}
}

// extractAggregationJoinFields extrait les champs de jointure
func (cp *ConstraintPipeline) extractAggregationJoinFields(constraintMap map[string]interface{}, aggInfo *AggregationInfo) {
	// Condition de jointure complète
	if joinCond, ok := constraintMap["join"]; ok {
		aggInfo.JoinCondition = joinCond
	}

	// Extraire depuis la condition de comparaison
	conditionData, ok := constraintMap["condition"].(map[string]interface{})
	if !ok {
		return
	}

	aggInfo.JoinCondition = conditionData

	condType, ok := conditionData["type"].(string)
	if !ok || condType != "comparison" {
		return
	}

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

// detectAggregation détecte si une contrainte contient une agrégation
// Cette fonction utilise une détection simple par string matching
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
		if cp.checkPatternsForAggregation(patternsData) {
			return true
		}
	}

	// Check old single-pattern syntax for backward compatibility
	if setData, hasSet := exprMap["set"]; hasSet {
		return cp.checkSetForAggregation(setData)
	}

	return false
}

// checkPatternsForAggregation vérifie si des patterns contiennent des variables d'agrégation
func (cp *ConstraintPipeline) checkPatternsForAggregation(patternsData interface{}) bool {
	patternsList, ok := patternsData.([]interface{})
	if !ok {
		return false
	}

	for _, patternInterface := range patternsList {
		patternMap, ok := patternInterface.(map[string]interface{})
		if !ok {
			continue
		}

		if cp.checkVariablesForAggregation(patternMap) {
			return true
		}
	}

	return false
}

// checkSetForAggregation vérifie si un set contient des variables d'agrégation
func (cp *ConstraintPipeline) checkSetForAggregation(setData interface{}) bool {
	setMap, ok := setData.(map[string]interface{})
	if !ok {
		return false
	}

	return cp.checkVariablesForAggregation(setMap)
}

// checkVariablesForAggregation vérifie si une liste de variables contient des agrégations
func (cp *ConstraintPipeline) checkVariablesForAggregation(container map[string]interface{}) bool {
	varsData, hasVars := container["variables"]
	if !hasVars {
		return false
	}

	varsList, ok := varsData.([]interface{})
	if !ok {
		return false
	}

	for _, varInterface := range varsList {
		varMap, ok := varInterface.(map[string]interface{})
		if !ok {
			continue
		}

		if varType, ok := varMap["type"].(string); ok && varType == "aggregationVariable" {
			return true
		}
	}

	return false
}
