// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import "fmt"

// evaluateFieldAccessValue évalue une valeur de type "fieldAccess" ou "field_access"
func (e *AlphaConditionEvaluator) evaluateFieldAccessValue(val map[string]interface{}) (interface{}, error) {
	// Supporter les deux formats: object/field et variable/field
	var objectOrVariable, field string
	var ok bool

	if objectOrVariable, ok = val["object"].(string); ok {
		// Format: object + field
	} else if objectOrVariable, ok = val["variable"].(string); ok {
		// Format: variable + field
	} else {
		return nil, fmt.Errorf("objet ou variable d'accès de champ invalide")
	}

	if field, ok = val["field"].(string); !ok {
		return nil, fmt.Errorf("champ d'accès invalide")
	}

	return e.evaluateFieldAccessByName(objectOrVariable, field)
}

// evaluateVariableValue évalue une valeur de type "variable"
func (e *AlphaConditionEvaluator) evaluateVariableValue(val map[string]interface{}) (interface{}, error) {
	name, ok := val["name"].(string)
	if !ok {
		return nil, fmt.Errorf("nom de variable invalide")
	}
	return e.evaluateVariableByName(name)
}

// evaluateNumberLiteralValue évalue une valeur de type "numberLiteral" ou "number"
func (e *AlphaConditionEvaluator) evaluateNumberLiteralValue(val map[string]interface{}) (interface{}, error) {
	value, ok := val["value"].(float64)
	if !ok {
		// Essayer aussi avec int
		if intValue, ok := val["value"].(int); ok {
			return float64(intValue), nil
		}
		return nil, fmt.Errorf("valeur numérique invalide")
	}
	return value, nil
}

// evaluateStringLiteralValue évalue une valeur de type "stringLiteral" ou "string"
func (e *AlphaConditionEvaluator) evaluateStringLiteralValue(val map[string]interface{}) (interface{}, error) {
	value, ok := val["value"].(string)
	if !ok {
		return nil, fmt.Errorf("valeur de chaîne invalide")
	}
	return value, nil
}

// evaluateBooleanLiteralValue évalue une valeur de type "booleanLiteral" ou "boolean"
func (e *AlphaConditionEvaluator) evaluateBooleanLiteralValue(val map[string]interface{}) (interface{}, error) {
	value, ok := val["value"].(bool)
	if !ok {
		return nil, fmt.Errorf("valeur booléenne invalide")
	}
	return value, nil
}

// evaluateFunctionCallValue évalue une valeur de type "functionCall" ou "function_call"
func (e *AlphaConditionEvaluator) evaluateFunctionCallValue(val map[string]interface{}) (interface{}, error) {
	// Déléguer à la fonction existante evaluateFunctionCall
	return e.evaluateFunctionCall(val)
}

// evaluateArrayLiteralValue évalue une valeur de type "arrayLiteral" ou "array_literal"
func (e *AlphaConditionEvaluator) evaluateArrayLiteralValue(val map[string]interface{}) (interface{}, error) {
	elements, ok := val["elements"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("éléments de tableau invalides")
	}

	// Évaluer chaque élément du tableau
	evaluatedElements := make([]interface{}, len(elements))
	for i, element := range elements {
		evaluatedElement, err := e.evaluateValue(element)
		if err != nil {
			return nil, fmt.Errorf("erreur évaluation élément tableau[%d]: %w", i, err)
		}
		evaluatedElements[i] = evaluatedElement
	}
	return evaluatedElements, nil
}

// evaluateCastValue évalue une valeur de type "cast"
func (e *AlphaConditionEvaluator) evaluateCastValue(val map[string]interface{}) (interface{}, error) {
	// Déléguer à la fonction existante evaluateCastExpression
	return e.evaluateCastExpression(val)
}

// evaluateBinaryOpValue évalue une valeur de type "binaryOp", "binary_operation", ou "binaryOperation"
func (e *AlphaConditionEvaluator) evaluateBinaryOpValue(val map[string]interface{}) (interface{}, error) {
	// Extraire et normaliser l'opérateur en utilisant l'utilitaire centralisé
	operator, err := ExtractOperatorFromMap(val)
	if err != nil {
		return nil, fmt.Errorf("erreur extraction opérateur: %w", err)
	}

	left, evalErr := e.evaluateValue(val["left"])
	if evalErr != nil {
		return nil, fmt.Errorf("erreur évaluation côté gauche (binaryOp %s): %w", operator, evalErr)
	}

	right, evalErr := e.evaluateValue(val["right"])
	if evalErr != nil {
		return nil, fmt.Errorf("erreur évaluation côté droit (binaryOp %s): %w", operator, evalErr)
	}

	// Distinguer les opérations arithmétiques des comparaisons
	switch operator {
	case "+", "-", "*", "/", "%":
		// Opération arithmétique - retourne une valeur numérique
		result, err := e.evaluateArithmeticOperation(left, operator, right)
		if err != nil {
			return nil, fmt.Errorf("erreur opération arithmétique %s: %w", operator, err)
		}
		return result, nil
	case "==", "!=", "<", "<=", ">", ">=", "CONTAINS", "IN", "LIKE", "MATCHES":
		// Opération de comparaison - retourne un booléen
		return e.compareValues(left, operator, right)
	default:
		return nil, fmt.Errorf("opérateur binaire non supporté: '%s' (bytes: %v)", operator, []byte(operator))
	}
}
