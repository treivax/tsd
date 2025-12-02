// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strings"

	"github.com/treivax/tsd/constraint"
)

// evaluateValue évalue une valeur (littéral, accès de champ, variable)
func (e *AlphaConditionEvaluator) evaluateValue(value interface{}) (interface{}, error) {
	// Check for nil values
	if value == nil {
		// Essayer d'obtenir un stacktrace pour débugger
		return nil, fmt.Errorf("valeur nil reçue dans evaluateValue (check caller)")
	}

	switch val := value.(type) {
	case map[string]interface{}:
		return e.evaluateValueFromMap(val)
	case constraint.FieldAccess:
		return e.evaluateFieldAccess(val)
	case constraint.Variable:
		return e.evaluateVariable(val)
	case constraint.NumberLiteral:
		return val.Value, nil
	case constraint.StringLiteral:
		return val.Value, nil
	case constraint.BooleanLiteral:
		return val.Value, nil
	case string:
		return val, nil
	case int, int32, int64:
		return val, nil
	case float32, float64:
		return val, nil
	case bool:
		return val, nil
	default:
		return nil, fmt.Errorf("type de valeur non supporté: %T", value)
	}
}

// evaluateValueFromMap évalue une valeur depuis une map
func (e *AlphaConditionEvaluator) evaluateValueFromMap(val map[string]interface{}) (interface{}, error) {
	valType, ok := val["type"].(string)
	if !ok {
		return nil, fmt.Errorf("type de valeur manquant dans map: %+v", val)
	}

	switch valType {
	case "fieldAccess", "field_access":
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

	case "variable":
		name, ok := val["name"].(string)
		if !ok {
			return nil, fmt.Errorf("nom de variable invalide")
		}
		return e.evaluateVariableByName(name)

	case "numberLiteral", "number":
		value, ok := val["value"].(float64)
		if !ok {
			// Essayer aussi avec int
			if intValue, ok := val["value"].(int); ok {
				return float64(intValue), nil
			}
			return nil, fmt.Errorf("valeur numérique invalide")
		}
		return value, nil

	case "stringLiteral", "string":
		value, ok := val["value"].(string)
		if !ok {
			return nil, fmt.Errorf("valeur de chaîne invalide")
		}
		return value, nil

	case "booleanLiteral", "boolean":
		value, ok := val["value"].(bool)
		if !ok {
			return nil, fmt.Errorf("valeur booléenne invalide")
		}
		return value, nil

	case "functionCall", "function_call":
		// Support des appels de fonction
		return e.evaluateFunctionCall(val)

	case "arrayLiteral", "array_literal":
		// Support des littéraux de tableau
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

	case "binaryOp", "binary_operation", "binaryOperation":
		// Support des opérations binaires
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

	default:
		return nil, fmt.Errorf("type de valeur non supporté: %s", valType)
	}
}

// evaluateFieldAccess évalue l'accès à un champ d'une variable
func (e *AlphaConditionEvaluator) evaluateFieldAccess(fa constraint.FieldAccess) (interface{}, error) {
	return e.evaluateFieldAccessByName(fa.Object, fa.Field)
}

// evaluateFieldAccessByName évalue l'accès à un champ par nom
func (e *AlphaConditionEvaluator) evaluateFieldAccessByName(object, field string) (interface{}, error) {
	fact, exists := e.variableBindings[object]
	if !exists {
		// En mode d'évaluation partielle (jointures en cascade), retourner nil sans erreur
		// pour permettre l'évaluation de continuer avec les variables disponibles
		if e.partialEvalMode {
			return nil, nil // Sentinel value indiquant que la variable n'est pas encore liée
		}
		// Debug info pour aider à diagnostiquer les problèmes de binding
		availableVars := make([]string, 0, len(e.variableBindings))
		for k := range e.variableBindings {
			availableVars = append(availableVars, k)
		}
		return nil, fmt.Errorf("variable non liée: %s (variables disponibles: %v)", object, availableVars)
	}

	// Cas spécial : le champ 'id' est stocké dans fact.ID, pas dans fact.Fields
	if field == "id" {
		return fact.ID, nil
	}

	value, exists := fact.Fields[field]
	if !exists {
		return nil, fmt.Errorf("champ inexistant: %s.%s", object, field)
	}

	return value, nil
}

// evaluateVariable évalue une variable
func (e *AlphaConditionEvaluator) evaluateVariable(variable constraint.Variable) (interface{}, error) {
	return e.evaluateVariableByName(variable.Name)
}

// evaluateVariableByName évalue une variable par nom
func (e *AlphaConditionEvaluator) evaluateVariableByName(name string) (interface{}, error) {
	// CORRECTION: Traiter les accès aux champs mal parsés comme "d.name"
	if strings.Contains(name, ".") {
		parts := strings.Split(name, ".")
		if len(parts) == 2 {
			objectName := parts[0]
			fieldName := parts[1]
			return e.evaluateFieldAccessByName(objectName, fieldName)
		}
	}

	fact, exists := e.variableBindings[name]
	if !exists {
		return nil, fmt.Errorf("variable non liée: %s", name)
	}

	// Retourner le fait entier ou une représentation
	return fact, nil
}
