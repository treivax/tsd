// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"reflect"
	"strings"
)

// compareValues compare deux valeurs avec un opérateur
func (e *AlphaConditionEvaluator) compareValues(left interface{}, operator string, right interface{}) (bool, error) {
	// En mode d'évaluation partielle, si l'une des valeurs est nil (variable non liée),
	// retourner true pour permettre l'évaluation de continuer
	if e.partialEvalMode && (left == nil || right == nil) {
		return true, nil
	}

	// Gérer les opérations arithmétiques qui retournent une valeur
	switch operator {
	case "+", "-", "*", "/", "%":
		return false, fmt.Errorf("opération arithmétique %s ne peut pas retourner un booléen", operator)
	}

	// Cas spécial: Comparaison de faits via IDs si les comparateurs avancés sont disponibles
	// Détecte si on compare deux faits (ex: l.user == u)
	leftFact, leftIsFact := left.(*Fact)
	rightFact, rightIsFact := right.(*Fact)

	if leftIsFact && rightIsFact && e.comparisonEvaluator != nil {
		// Comparaison de deux faits via leurs IDs
		return e.comparisonEvaluator.EvaluateComparison(
			leftFact.ID,
			rightFact.ID,
			operator,
			FieldTypeFact,
			FieldTypeFact,
		)
	}

	// Normaliser les valeurs numériques
	leftVal := e.normalizeValue(left)
	rightVal := e.normalizeValue(right)

	switch operator {
	case "==":
		return e.areEqual(leftVal, rightVal), nil
	case "!=", "<>":
		return !e.areEqual(leftVal, rightVal), nil
	case "<":
		return e.isLess(leftVal, rightVal)
	case "<=":
		equal := e.areEqual(leftVal, rightVal)
		less, err := e.isLess(leftVal, rightVal)
		return equal || less, err
	case ">":
		return e.isGreater(leftVal, rightVal)
	case ">=":
		equal := e.areEqual(leftVal, rightVal)
		greater, err := e.isGreater(leftVal, rightVal)
		return equal || greater, err
	case "CONTAINS":
		return e.evaluateContains(leftVal, rightVal)
	case "IN":
		return e.evaluateIn(leftVal, rightVal)
	case "LIKE":
		return e.evaluateLike(leftVal, rightVal)
	case "MATCHES":
		return e.evaluateMatches(leftVal, rightVal)
	default:
		return false, fmt.Errorf("opérateur non supporté: %s", operator)
	}
}

// normalizeValue normalise une valeur pour la comparaison
func (e *AlphaConditionEvaluator) normalizeValue(value interface{}) interface{} {
	switch v := value.(type) {
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case float32:
		return float64(v)
	default:
		return value
	}
}

// areEqual vérifie si deux valeurs sont égales
func (e *AlphaConditionEvaluator) areEqual(left, right interface{}) bool {
	return reflect.DeepEqual(left, right)
}

// isLess vérifie si left < right
func (e *AlphaConditionEvaluator) isLess(left, right interface{}) (bool, error) {
	switch leftVal := left.(type) {
	case float64:
		if rightVal, ok := right.(float64); ok {
			return leftVal < rightVal, nil
		}
	case string:
		if rightVal, ok := right.(string); ok {
			return strings.Compare(leftVal, rightVal) < 0, nil
		}
	}
	return false, fmt.Errorf("impossible de comparer %T avec %T", left, right)
}

// isGreater vérifie si left > right
func (e *AlphaConditionEvaluator) isGreater(left, right interface{}) (bool, error) {
	switch leftVal := left.(type) {
	case float64:
		if rightVal, ok := right.(float64); ok {
			return leftVal > rightVal, nil
		}
	case string:
		if rightVal, ok := right.(string); ok {
			return strings.Compare(leftVal, rightVal) > 0, nil
		}
	}
	return false, fmt.Errorf("impossible de comparer %T avec %T", left, right)
}
