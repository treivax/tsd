// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
)

// GetValueType retourne le type d'une valeur dans l'AST
func GetValueType(value interface{}) string {
	switch v := value.(type) {
	case map[string]interface{}:
		valueType, ok := v["type"].(string)
		if !ok {
			return ValueTypeUnknown
		}
		switch valueType {
		case ValueTypeNumber:
			return ValueTypeNumber
		case ValueTypeString:
			return ValueTypeString
		case ValueTypeBoolean:
			return ValueTypeBool
		case ValueTypeVariable:
			// Pour les variables comme "true", "false" qui sont parsées comme variables
			name, ok := v["name"].(string)
			if ok {
				switch name {
				case "true", "false":
					return ValueTypeBool
				}
			}
			return ValueTypeVariable // Type non déterminable sans contexte
		}
	}
	return ValueTypeUnknown
}

// ValidateTypeCompatibility vérifie la compatibilité des types dans les comparaisons
func ValidateTypeCompatibility(program Program, constraint interface{}, expressionIndex int) error {
	constraintMap, ok := constraint.(map[string]interface{})
	if !ok {
		return nil
	}

	constraintType, ok := constraintMap["type"].(string)
	if !ok {
		return nil
	}

	switch constraintType {
	case ConstraintTypeComparison:
		return validateConstraintWithOperands(program, constraintMap, expressionIndex, true)
	case ConstraintTypeLogicalExpr:
		return validateLogicalExpressionConstraint(program, constraintMap, expressionIndex)
	case ConstraintTypeBinaryOp:
		return validateConstraintWithOperands(program, constraintMap, expressionIndex, false)
	}
	return nil
}

// validateConstraintWithOperands handles validation for constraints with left/right operands
func validateConstraintWithOperands(program Program, c map[string]interface{}, expressionIndex int, checkCompatibility bool) error {
	left := c["left"]
	right := c["right"]

	if left == nil || right == nil {
		return nil
	}

	// Validate type compatibility between operands (only for comparisons)
	if checkCompatibility {
		if err := validateOperandTypeCompatibility(program, left, right, expressionIndex); err != nil {
			return err
		}
	}

	// Recursive validation for operands
	if err := ValidateTypeCompatibility(program, left, expressionIndex); err != nil {
		return err
	}
	if err := ValidateTypeCompatibility(program, right, expressionIndex); err != nil {
		return err
	}

	return nil
}

// validateOperandTypeCompatibility checks if two operands have compatible types
func validateOperandTypeCompatibility(program Program, left, right interface{}, expressionIndex int) error {
	leftType, err := getOperandType(program, left, expressionIndex)
	if err != nil {
		return err
	}

	rightType, err := getOperandType(program, right, expressionIndex)
	if err != nil {
		return err
	}

	// Skip type compatibility check for variable vs number comparisons
	// This handles aggregation variables which are always numeric
	if (leftType == ValueTypeVariable && rightType == ValueTypeNumber) ||
		(leftType == ValueTypeNumber && rightType == ValueTypeVariable) {
		return nil
	}

	// Check compatibility
	if leftType != ValueTypeUnknown && rightType != ValueTypeUnknown && rightType != ValueTypeVariable {
		if leftType != rightType {
			return fmt.Errorf("incompatibilité de types dans la comparaison: %s vs %s", leftType, rightType)
		}
	}

	return nil
}

// getOperandType determines the type of an operand in a constraint
func getOperandType(program Program, operand interface{}, expressionIndex int) (string, error) {
	operandMap, ok := operand.(map[string]interface{})
	if !ok {
		return GetValueType(operand), nil
	}

	if operandMap["type"] == ConstraintTypeFieldAccess {
		object := operandMap["object"].(string)
		field := operandMap["field"].(string)
		return GetFieldType(program, object, field, expressionIndex)
	}

	return GetValueType(operand), nil
}

// validateLogicalExpressionConstraint handles logical expression validation
func validateLogicalExpressionConstraint(program Program, c map[string]interface{}, expressionIndex int) error {
	if left := c["left"]; left != nil {
		if err := ValidateTypeCompatibility(program, left, expressionIndex); err != nil {
			return err
		}
	}

	operations, ok := c["operations"].([]interface{})
	if !ok {
		return nil
	}

	for _, op := range operations {
		opMap, ok := op.(map[string]interface{})
		if !ok {
			continue
		}

		if right := opMap["right"]; right != nil {
			if err := ValidateTypeCompatibility(program, right, expressionIndex); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateBinaryOpConstraint handles binary operation validation (wrapper for backward compatibility)
func validateBinaryOpConstraint(program Program, c map[string]interface{}, expressionIndex int) error {
	return validateConstraintWithOperands(program, c, expressionIndex, false)
}
