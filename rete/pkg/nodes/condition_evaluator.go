// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package nodes

import (
	"fmt"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// ConditionEvaluator provides shared utilities for evaluating conditions
// across different node types (NOT, EXISTS, etc.)
type ConditionEvaluator struct{}

// NewConditionEvaluator creates a new condition evaluator
func NewConditionEvaluator() *ConditionEvaluator {
	return &ConditionEvaluator{}
}

// EvaluateCondition evaluates a condition against a token and fact
func (ce *ConditionEvaluator) EvaluateCondition(condition interface{}, token *domain.Token, fact *domain.Fact) (bool, error) {
	// Check if it's a basic condition with comparison
	if conditionMap, ok := condition.(map[string]interface{}); ok {
		if conditionType, hasType := conditionMap["type"]; hasType {
			switch conditionType {
			case "binaryOperation", "binary_op":
				return ce.EvaluateBinaryCondition(conditionMap, token, fact)
			case "simple":
				return ce.EvaluateBinaryCondition(conditionMap, token, fact)
			}
		}
	}
	return false, fmt.Errorf("format de condition non reconnu")
}

// EvaluateBinaryCondition evaluates a binary operation condition
// Expected format: {left: {variable: "p", field: "age"}, operator: "==", right: {value: 0}}
func (ce *ConditionEvaluator) EvaluateBinaryCondition(conditionMap map[string]interface{}, token *domain.Token, fact *domain.Fact) (bool, error) {
	// Extract operator
	operator, ok := conditionMap["operator"].(string)
	if !ok {
		operator, ok = conditionMap["op"].(string)
		if !ok {
			return false, fmt.Errorf("opérateur manquant")
		}
	}

	// Extract left side (field value)
	leftExpr, hasLeft := conditionMap["left"]
	if !hasLeft {
		return false, fmt.Errorf("expression gauche manquante")
	}
	leftValue, err := ce.ExtractFieldValue(leftExpr, token, fact)
	if err != nil {
		return false, err
	}

	// Extract right side (constant or field value)
	rightExpr, hasRight := conditionMap["right"]
	if !hasRight {
		return false, fmt.Errorf("expression droite manquante")
	}
	rightValue, err := ce.ExtractConstantValue(rightExpr)
	if err != nil {
		return false, err
	}

	// Compare values
	return ce.CompareValues(leftValue, operator, rightValue)
}

// ExtractFieldValue extracts a field value from token or fact
func (ce *ConditionEvaluator) ExtractFieldValue(leftExpr interface{}, token *domain.Token, fact *domain.Fact) (interface{}, error) {
	if leftMap, ok := leftExpr.(map[string]interface{}); ok {
		if fieldName, hasField := leftMap["field"].(string); hasField {
			// Look for value in the primary fact of the token
			if len(token.Facts) > 0 {
				primaryFact := token.Facts[0]
				if value, exists := primaryFact.Fields[fieldName]; exists {
					return value, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("impossible d'extraire la valeur du champ")
}

// ExtractConstantValue extracts a constant value from expression
func (ce *ConditionEvaluator) ExtractConstantValue(rightExpr interface{}) (interface{}, error) {
	if rightMap, ok := rightExpr.(map[string]interface{}); ok {
		if value, hasValue := rightMap["value"]; hasValue {
			return value, nil
		}
	}
	// If it's directly the value
	return rightExpr, nil
}

// CompareValues compares two values using the given operator
func (ce *ConditionEvaluator) CompareValues(left interface{}, operator string, right interface{}) (bool, error) {
	switch operator {
	case "==":
		return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right), nil
	case "!=":
		return fmt.Sprintf("%v", left) != fmt.Sprintf("%v", right), nil
	case "<":
		return ce.NumericCompare(left, right, func(l, r float64) bool { return l < r })
	case ">":
		return ce.NumericCompare(left, right, func(l, r float64) bool { return l > r })
	case "<=":
		return ce.NumericCompare(left, right, func(l, r float64) bool { return l <= r })
	case ">=":
		return ce.NumericCompare(left, right, func(l, r float64) bool { return l >= r })
	default:
		return false, fmt.Errorf("opérateur non supporté: %s", operator)
	}
}

// NumericCompare performs numeric comparison using a comparison function
func (ce *ConditionEvaluator) NumericCompare(left, right interface{}, compareFunc func(float64, float64) bool) (bool, error) {
	leftFloat, err := ce.ToFloat64(left)
	if err != nil {
		return false, err
	}
	rightFloat, err := ce.ToFloat64(right)
	if err != nil {
		return false, err
	}
	return compareFunc(leftFloat, rightFloat), nil
}

// ToFloat64 converts a value to float64
func (ce *ConditionEvaluator) ToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("impossible de convertir %T en float64", value)
	}
}
