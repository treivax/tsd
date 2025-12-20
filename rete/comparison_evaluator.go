// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"math"
)

const ComparisonEpsilon = 1e-9

// ComparisonEvaluator évalue les comparaisons entre valeurs avec support des types
type ComparisonEvaluator struct {
	resolver *FieldResolver
}

// NewComparisonEvaluator crée un nouvel évaluateur de comparaisons
func NewComparisonEvaluator(resolver *FieldResolver) *ComparisonEvaluator {
	return &ComparisonEvaluator{
		resolver: resolver,
	}
}

// EvaluateComparison évalue une comparaison entre deux valeurs
// Gère les comparaisons de primitifs ET les comparaisons de faits (via IDs)
func (ce *ComparisonEvaluator) EvaluateComparison(left, right interface{}, operator string, leftType, rightType string) (bool, error) {
	// Cas 1: Les deux valeurs sont des IDs de faits
	if leftType == FieldTypeFact && rightType == FieldTypeFact {
		return ce.compareFactIDs(left, right, operator)
	}

	// Cas 2: Les deux valeurs sont des primitifs
	if leftType == FieldTypePrimitive && rightType == FieldTypePrimitive {
		return ce.comparePrimitives(left, right, operator)
	}

	// Cas 3: Types incompatibles
	return false, fmt.Errorf("comparaison impossible entre types '%s' et '%s'", leftType, rightType)
}

// compareFactIDs compare deux IDs de faits
func (ce *ComparisonEvaluator) compareFactIDs(left, right interface{}, operator string) (bool, error) {
	leftID, ok1 := left.(string)
	rightID, ok2 := right.(string)

	if !ok1 || !ok2 {
		return false, fmt.Errorf("IDs de faits doivent être des strings")
	}

	switch operator {
	case "==":
		return leftID == rightID, nil
	case "!=":
		return leftID != rightID, nil
	default:
		return false, fmt.Errorf("opérateur '%s' non supporté pour les comparaisons de faits (seuls == et != sont autorisés)", operator)
	}
}

// comparePrimitives compare deux valeurs primitives
func (ce *ComparisonEvaluator) comparePrimitives(left, right interface{}, operator string) (bool, error) {
	// Essayer de comparer comme strings
	leftStr, leftIsStr := left.(string)
	rightStr, rightIsStr := right.(string)

	if leftIsStr && rightIsStr {
		return ce.compareStrings(leftStr, rightStr, operator)
	}

	// Essayer de comparer comme numbers
	leftNum, leftIsNum := convertToFloat64(left)
	rightNum, rightIsNum := convertToFloat64(right)

	if leftIsNum && rightIsNum {
		return ce.compareNumbers(leftNum, rightNum, operator)
	}

	// Essayer de comparer comme booleans
	leftBool, leftIsBool := left.(bool)
	rightBool, rightIsBool := right.(bool)

	if leftIsBool && rightIsBool {
		return ce.compareBooleans(leftBool, rightBool, operator)
	}

	// Types incompatibles
	return false, fmt.Errorf("types incompatibles pour comparaison: %T et %T", left, right)
}

// compareStrings compare deux strings
func (ce *ComparisonEvaluator) compareStrings(left, right, operator string) (bool, error) {
	switch operator {
	case "==":
		return left == right, nil
	case "!=":
		return left != right, nil
	case "<":
		return left < right, nil
	case "<=":
		return left <= right, nil
	case ">":
		return left > right, nil
	case ">=":
		return left >= right, nil
	default:
		return false, fmt.Errorf("opérateur '%s' non supporté pour strings", operator)
	}
}

// compareNumbers compare deux numbers
func (ce *ComparisonEvaluator) compareNumbers(left, right float64, operator string) (bool, error) {
	switch operator {
	case "==":
		return math.Abs(left-right) < ComparisonEpsilon, nil
	case "!=":
		return math.Abs(left-right) >= ComparisonEpsilon, nil
	case "<":
		return left < right, nil
	case "<=":
		return left <= right, nil
	case ">":
		return left > right, nil
	case ">=":
		return left >= right, nil
	default:
		return false, fmt.Errorf("opérateur '%s' non supporté pour numbers", operator)
	}
}

// compareBooleans compare deux booleans
func (ce *ComparisonEvaluator) compareBooleans(left, right bool, operator string) (bool, error) {
	switch operator {
	case "==":
		return left == right, nil
	case "!=":
		return left != right, nil
	default:
		return false, fmt.Errorf("opérateur '%s' non supporté pour booleans (seuls == et != sont autorisés)", operator)
	}
}
