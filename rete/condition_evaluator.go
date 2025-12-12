// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strings"
)

// ConditionEvaluator evaluates conditions with support for intermediate results.
// It can resolve references to temporary results stored in the EvaluationContext.
type ConditionEvaluator struct {
	storage Storage
}

// NewConditionEvaluator creates a new condition evaluator.
func NewConditionEvaluator(storage Storage) *ConditionEvaluator {
	return &ConditionEvaluator{
		storage: storage,
	}
}

// EvaluateWithContext evaluates a condition using the evaluation context.
// The context provides access to intermediate results from previous steps.
func (ce *ConditionEvaluator) EvaluateWithContext(
	condition interface{},
	fact *Fact,
	context *EvaluationContext,
) (interface{}, error) {
	condMap, ok := condition.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("condition must be a map, got %T", condition)
	}

	condType, _ := condMap["type"].(string)

	switch condType {
	case "binaryOp", "binaryOperation":
		return ce.evaluateBinaryOp(condMap, fact, context)

	case "comparison":
		return ce.evaluateComparison(condMap, fact, context)

	case "fieldAccess":
		return ce.evaluateFieldAccess(condMap, fact, context)

	case "number", "numberLiteral":
		// Extract numeric value
		if value, ok := condMap["value"]; ok {
			return value, nil
		}
		return nil, fmt.Errorf("number literal missing value")

	case "string", "stringLiteral":
		// Extract string value
		if value, ok := condMap["value"]; ok {
			return value, nil
		}
		return nil, fmt.Errorf("string literal missing value")

	case "boolean", "booleanLiteral":
		// Extract boolean value
		if value, ok := condMap["value"].(bool); ok {
			return value, nil
		}
		return nil, fmt.Errorf("boolean literal missing value")

	case "arrayLiteral", "array_literal":
		// Extract array value
		if elements, ok := condMap["elements"].([]interface{}); ok {
			// Evaluate each element
			evaluatedElements := make([]interface{}, len(elements))
			for i, element := range elements {
				evaluatedElement, err := ce.EvaluateWithContext(element, fact, context)
				if err != nil {
					return nil, fmt.Errorf("error evaluating array element %d: %w", i, err)
				}
				evaluatedElements[i] = evaluatedElement
			}
			return evaluatedElements, nil
		}
		return nil, fmt.Errorf("array literal missing elements")

	case "tempResult":
		// KEY FEATURE: Resolve intermediate result from context
		return ce.resolveTempResult(condMap, context)

	case "constraint":
		// Handle constraint type by evaluating its nested condition
		if innerCondition, ok := condMap["condition"]; ok {
			return ce.EvaluateWithContext(innerCondition, fact, context)
		}
		// If no nested condition, treat as always true
		return true, nil

	case "logicalExpression", "logical_op", "logicalExpr":
		// Handle logical expressions (AND, OR, NOT)
		return ce.evaluateLogicalExpression(condMap, fact, context)

	default:
		return nil, fmt.Errorf("unsupported condition type: %s", condType)
	}
}

// resolveTempResult resolves a reference to an intermediate result from the context.
func (ce *ConditionEvaluator) resolveTempResult(
	tempRef map[string]interface{},
	context *EvaluationContext,
) (interface{}, error) {
	// Extract the result name from the tempResult reference
	var resultName string

	if name, ok := tempRef["step_name"].(string); ok {
		resultName = name
	} else if hash, ok := tempRef["hash"].(string); ok {
		resultName = hash
	} else if name, ok := tempRef["name"].(string); ok {
		resultName = name
	} else {
		return nil, fmt.Errorf("tempResult missing identifier (step_name, hash, or name)")
	}

	// Retrieve from context
	value, exists := context.GetIntermediateResult(resultName)
	if !exists {
		return nil, fmt.Errorf("intermediate result %s not found in context", resultName)
	}

	return value, nil
}

// evaluateBinaryOp evaluates a binary arithmetic operation.
func (ce *ConditionEvaluator) evaluateBinaryOp(
	op map[string]interface{},
	fact *Fact,
	context *EvaluationContext,
) (interface{}, error) {
	// Recursively evaluate left and right operands
	left, err := ce.EvaluateWithContext(op["left"], fact, context)
	if err != nil {
		return nil, fmt.Errorf("error evaluating left operand: %w", err)
	}

	right, err := ce.EvaluateWithContext(op["right"], fact, context)
	if err != nil {
		return nil, fmt.Errorf("error evaluating right operand: %w", err)
	}

	operator, _ := op["operator"].(string)

	// Apply the operator
	return ce.applyOperator(left, operator, right)
}

// applyOperator applies an arithmetic operator to two operands.
func (ce *ConditionEvaluator) applyOperator(left interface{}, op string, right interface{}) (interface{}, error) {
	// Convert to float64 for calculations
	leftFloat, err := convertValueToFloat64(left)
	if err != nil {
		return nil, fmt.Errorf("error converting left operand: %w", err)
	}

	rightFloat, err := convertValueToFloat64(right)
	if err != nil {
		return nil, fmt.Errorf("error converting right operand: %w", err)
	}

	switch op {
	case "*", "Kg==":
		return leftFloat * rightFloat, nil
	case "+", "Kw==":
		return leftFloat + rightFloat, nil
	case "-", "LQ==":
		return leftFloat - rightFloat, nil
	case "/", "Lw==":
		if rightFloat == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return leftFloat / rightFloat, nil
	case "%":
		if rightFloat == 0 {
			return nil, fmt.Errorf("modulo by zero")
		}
		return float64(int64(leftFloat) % int64(rightFloat)), nil
	default:
		return nil, fmt.Errorf("unsupported operator: %s", op)
	}
}

// evaluateComparison evaluates a comparison operation.
func (ce *ConditionEvaluator) evaluateComparison(
	comp map[string]interface{},
	fact *Fact,
	context *EvaluationContext,
) (interface{}, error) {
	// Recursively evaluate left and right operands
	left, err := ce.EvaluateWithContext(comp["left"], fact, context)
	if err != nil {
		return nil, fmt.Errorf("error evaluating left operand: %w", err)
	}

	right, err := ce.EvaluateWithContext(comp["right"], fact, context)
	if err != nil {
		return nil, fmt.Errorf("error evaluating right operand: %w", err)
	}

	operator, _ := comp["operator"].(string)

	// Handle CONTAINS operator for strings
	if operator == "CONTAINS" {
		leftStr, leftOk := left.(string)
		rightStr, rightOk := right.(string)
		if leftOk && rightOk {
			return strings.Contains(leftStr, rightStr), nil
		}
		return nil, fmt.Errorf("CONTAINS operator requires string operands, got %T and %T", left, right)
	}

	// Handle IN operator for arrays
	if operator == "IN" {
		rightSlice, ok := right.([]interface{})
		if !ok {
			return nil, fmt.Errorf("IN operator requires array on right side, got %T", right)
		}
		for _, item := range rightSlice {
			if left == item {
				return true, nil
			}
		}
		return false, nil
	}

	leftFloat, leftErr := convertValueToFloat64(left)
	rightFloat, rightErr := convertValueToFloat64(right)

	if leftErr == nil && rightErr == nil {
		switch operator {
		case ">":
			return leftFloat > rightFloat, nil
		case "<":
			return leftFloat < rightFloat, nil
		case ">=":
			return leftFloat >= rightFloat, nil
		case "<=":
			return leftFloat <= rightFloat, nil
		case "==":
			return leftFloat == rightFloat, nil
		case "!=":
			return leftFloat != rightFloat, nil
		default:
			return nil, fmt.Errorf("unsupported comparison operator: %s", operator)
		}
	}

	switch operator {
	case "==":
		return left == right, nil
	case "!=":
		return left != right, nil
	default:
		return nil, fmt.Errorf("comparison operator %s requires numeric operands", operator)
	}
}

// evaluateFieldAccess extracts a field value from the fact.
func (ce *ConditionEvaluator) evaluateFieldAccess(
	fieldAccess map[string]interface{},
	fact *Fact,
	context *EvaluationContext,
) (interface{}, error) {
	fieldName, ok := fieldAccess["field"].(string)
	if !ok {
		return nil, fmt.Errorf("fieldAccess missing field name")
	}

	value, exists := fact.GetField(fieldName)
	if !exists {
		return nil, fmt.Errorf("field %s not found in fact %s", fieldName, fact.ID)
	}

	return value, nil
}

// evaluateLogicalExpression evaluates a logical expression (AND, OR, NOT).
func (ce *ConditionEvaluator) evaluateLogicalExpression(
	expr map[string]interface{},
	fact *Fact,
	context *EvaluationContext,
) (interface{}, error) {
	// Handle NOT operator
	if operator, ok := expr["operator"].(string); ok && operator == "NOT" {
		// Evaluate the inner expression
		inner, err := ce.EvaluateWithContext(expr["right"], fact, context)
		if err != nil {
			return nil, fmt.Errorf("error evaluating NOT expression: %w", err)
		}
		// Convert to boolean and negate
		if boolVal, ok := inner.(bool); ok {
			return !boolVal, nil
		}
		return nil, fmt.Errorf("NOT operator requires boolean operand, got %T", inner)
	}

	// Handle binary logical operators (AND, OR)
	left, err := ce.EvaluateWithContext(expr["left"], fact, context)
	if err != nil {
		return nil, fmt.Errorf("error evaluating left operand: %w", err)
	}

	// Check for operations array
	operations, hasOps := expr["operations"].([]interface{})
	if !hasOps || len(operations) == 0 {
		// No operations, just return left value
		if boolVal, ok := left.(bool); ok {
			return boolVal, nil
		}
		return nil, fmt.Errorf("logical expression requires boolean operand, got %T", left)
	}

	leftBool, ok := left.(bool)
	if !ok {
		return nil, fmt.Errorf("logical expression left operand must be boolean, got %T", left)
	}

	result := leftBool
	for _, opInterface := range operations {
		opMap, ok := opInterface.(map[string]interface{})
		if !ok {
			continue
		}

		operator, ok := opMap["op"].(string)
		if !ok {
			continue
		}

		right, err := ce.EvaluateWithContext(opMap["right"], fact, context)
		if err != nil {
			return nil, fmt.Errorf("error evaluating right operand for %s: %w", operator, err)
		}

		rightBool, ok := right.(bool)
		if !ok {
			return nil, fmt.Errorf("logical expression right operand must be boolean, got %T", right)
		}

		switch operator {
		case "AND":
			result = result && rightBool
		case "OR":
			result = result || rightBool
		default:
			return nil, fmt.Errorf("unsupported logical operator: %s", operator)
		}
	}

	return result, nil
}

// convertValueToFloat64 converts a value to float64 for numeric operations.
func convertValueToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}
