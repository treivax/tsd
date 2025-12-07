// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import "fmt"

// aggregation_extraction.go contient les fonctions décomposées pour l'extraction
// des informations d'agrégation depuis la syntaxe multi-pattern.
// Ces fonctions ont été extraites de extractAggregationInfoFromVariables()
// pour réduire la complexité et améliorer la testabilité.

// parseAggregationExpression valide et extrait la structure de base (patterns, variables)
// Retourne le premier pattern et la liste de ses variables.
func (cp *ConstraintPipeline) parseAggregationExpression(exprMap map[string]interface{}) (map[string]interface{}, []interface{}, error) {
	// Valider et extraire le premier pattern
	firstPattern, err := getFirstPattern(exprMap)
	if err != nil {
		return nil, nil, err
	}

	// Extraire la liste des variables du premier pattern
	varsList, err := getVariablesList(firstPattern)
	if err != nil {
		return nil, nil, err
	}

	return firstPattern, varsList, nil
}

// extractAggregationFunction extrait la fonction d'agrégation (AVG, SUM, COUNT, etc.)
// depuis une variable d'agrégation. Supporte deux formats :
// 1. Direct: varMap["function"]
// 2. Nested: varMap["value"]["function"]
func (cp *ConstraintPipeline) extractAggregationFunction(varMap map[string]interface{}) (string, error) {
	// Format direct: varMap["function"]
	if function, ok := extractStringField(varMap, "function"); ok {
		return function, nil
	}

	// Format nested: varMap["value"]["function"]
	valueData, ok := extractMapField(varMap, "value")
	if !ok {
		return "", fmt.Errorf("no function found in aggregation variable")
	}

	// Vérifier le type de la valeur
	if isFunctionCallType(valueData) {
		if fnName, ok := extractStringField(valueData, "function"); ok {
			return fnName, nil
		}
	}

	return "", fmt.Errorf("no function found in aggregation variable value")
}

// extractAggregationField extrait le champ agrégé et la variable source.
// Supporte deux formats :
// 1. Direct field: varMap["field"]
// 2. Nested in value: varMap["value"]["arguments"][0]
func (cp *ConstraintPipeline) extractAggregationField(varMap map[string]interface{}) (aggVariable, field string, err error) {
	// Format direct: varMap["field"]
	if fieldData, ok := extractMapField(varMap, "field"); ok {
		if fieldObj, ok := extractStringField(fieldData, "object"); ok {
			aggVariable = fieldObj
		}
		if fieldName, ok := extractStringField(fieldData, "field"); ok {
			field = fieldName
		}
		if aggVariable != "" && field != "" {
			return aggVariable, field, nil
		}
	}

	// Format nested: varMap["value"]["arguments"][0]
	valueData, ok := extractMapField(varMap, "value")
	if !ok {
		return "", "", fmt.Errorf("no field information found in aggregation variable")
	}

	argsData, ok := extractListField(valueData, "arguments")
	if !ok || len(argsData) == 0 {
		return "", "", fmt.Errorf("no arguments found in aggregation value")
	}

	argMap, ok := argsData[0].(map[string]interface{})
	if !ok {
		return "", "", fmt.Errorf("first argument is not a map")
	}

	if !isFieldAccessType(argMap) {
		return "", "", fmt.Errorf("first argument is not a field access")
	}

	objName, hasObj := extractStringField(argMap, "object")
	fieldName, hasField := extractStringField(argMap, "field")

	if !hasObj || !hasField {
		return "", "", fmt.Errorf("field access missing object or field")
	}

	return objName, fieldName, nil
}

// extractSourceType extrait le type de la source d'agrégation depuis le second pattern
func (cp *ConstraintPipeline) extractSourceType(exprMap map[string]interface{}) (string, error) {
	secondPattern, err := getSecondPattern(exprMap)
	if err != nil {
		return "", err
	}

	varsData, hasVars := secondPattern["variables"]
	if !hasVars {
		return "", fmt.Errorf("no variables in second pattern")
	}

	varsList, ok := varsData.([]interface{})
	if !ok || len(varsList) == 0 {
		return "", fmt.Errorf("second pattern variables is not a list or is empty")
	}

	varMap, ok := varsList[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("first variable in second pattern is not a map")
	}

	aggType, ok := extractStringField(varMap, "dataType")
	if !ok {
		return "", fmt.Errorf("no dataType in second pattern variable")
	}

	return aggType, nil
}

// extractJoinFields extrait les champs de jointure depuis les conditions de jointure
func (cp *ConstraintPipeline) extractJoinFields(joinConditions map[string]interface{}) (joinField, mainField string) {
	if !isComparisonType(joinConditions) {
		return "", ""
	}

	// Left side: champ de jointure (e.g., e.deptId)
	if leftData, ok := extractMapField(joinConditions, "left"); ok {
		if isFieldAccessType(leftData) {
			if jf, ok := extractStringField(leftData, "field"); ok {
				joinField = jf
			}
		}
	}

	// Right side: champ principal (e.g., d.id)
	if rightData, ok := extractMapField(joinConditions, "right"); ok {
		if isFieldAccessType(rightData) {
			if mf, ok := extractStringField(rightData, "field"); ok {
				mainField = mf
			}
		}
	}

	return joinField, mainField
}

// extractThresholdConditions extrait l'opérateur et la valeur seuil depuis les conditions de seuil
func (cp *ConstraintPipeline) extractThresholdConditions(thresholdConditions []map[string]interface{}) (operator string, threshold float64) {
	if len(thresholdConditions) == 0 {
		// Pas de seuil - toujours déclencher (valeur par défaut)
		return DefaultThresholdOperator, DefaultThresholdValue
	}

	// Utiliser la première condition de seuil trouvée
	firstThreshold := thresholdConditions[0]

	// Extraire l'opérateur
	if op, ok := extractStringField(firstThreshold, "operator"); ok {
		operator = op
	} else {
		operator = DefaultThresholdOperator
	}

	// Extraire la valeur
	if rightData, ok := extractMapField(firstThreshold, "right"); ok {
		if value, ok := extractFloat64Field(rightData, "value"); ok {
			threshold = value
		} else {
			threshold = DefaultThresholdValue
		}
	} else {
		threshold = DefaultThresholdValue
	}

	return operator, threshold
}
