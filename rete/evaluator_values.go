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

	// Dispatcher simplifié vers les évaluateurs spécifiques
	switch valType {
	case "fieldAccess", "field_access":
		return e.evaluateFieldAccessValue(val)
	case "variable":
		return e.evaluateVariableValue(val)
	case "numberLiteral", "number":
		return e.evaluateNumberLiteralValue(val)
	case "stringLiteral", "string":
		return e.evaluateStringLiteralValue(val)
	case "booleanLiteral", "boolean":
		return e.evaluateBooleanLiteralValue(val)
	case "functionCall", "function_call":
		return e.evaluateFunctionCallValue(val)
	case "arrayLiteral", "array_literal":
		return e.evaluateArrayLiteralValue(val)
	case "cast":
		return e.evaluateCastValue(val)
	case "binaryOp", "binary_operation", "binaryOperation":
		return e.evaluateBinaryOpValue(val)
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

	// Cas spécial : le champ '_id_' est INTERDIT dans les expressions
	// Ce check ne devrait jamais être atteint car la validation devrait rejeter
	// l'accès à _id_ avant l'évaluation, mais on le garde par sécurité.
	if field == FieldNameID {
		return nil, fmt.Errorf(
			"le champ '_id_' est interne et ne peut pas être accédé dans les expressions",
		)
	}

	// Cas spécial : le champ 'id' retourne l'ID interne du fait
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
