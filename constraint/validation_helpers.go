// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"encoding/base64"
	"fmt"
	"unicode/utf8"
)

// findVariableType searches for a variable in an expression and returns its type
// It handles both old single-pattern (Set) and new multi-pattern (Patterns) syntax
func findVariableType(expr Expression, varName string) (string, error) {
	if varName == "" {
		return "", fmt.Errorf("variable name cannot be empty")
	}

	// Check new multi-pattern syntax first
	if len(expr.Patterns) > 0 {
		for _, pattern := range expr.Patterns {
			for _, variable := range pattern.Variables {
				if variable.Name == varName {
					return variable.DataType, nil
				}
			}
		}
	} else if expr.Set.Variables != nil {
		// Old single-pattern syntax (backward compatibility)
		for _, variable := range expr.Set.Variables {
			if variable.Name == varName {
				return variable.DataType, nil
			}
		}
	}

	return "", fmt.Errorf("variable '%s' not found in expression", varName)
}

// sanitizeForLog sanitizes a string for safe logging by removing control characters
// and limiting length to prevent log injection attacks
func sanitizeForLog(s string, maxLen int) string {
	if maxLen <= 0 {
		maxLen = 256
	}

	// Limit length first
	if len(s) > maxLen {
		if utf8.ValidString(s) {
			// Count runes for proper UTF-8 truncation
			runes := []rune(s)
			if len(runes) > maxLen {
				s = string(runes[:maxLen]) + "..."
			}
		} else {
			// Binary data, truncate bytes
			s = s[:maxLen] + "..."
		}
	}

	// Remove control characters except newline and tab
	result := make([]rune, 0, len(s))
	for _, r := range s {
		// Allow printable characters, newline, tab
		if r == '\n' || r == '\t' || (r >= 32 && r < 127) || r >= 128 {
			result = append(result, r)
		} else {
			result = append(result, '�') // Replacement character
		}
	}

	return string(result)
}

// validateInputNotNil checks if critical inputs are not nil and returns error if they are
func validateInputNotNil(inputs map[string]interface{}) error {
	for name, value := range inputs {
		if value == nil {
			return fmt.Errorf("required input '%s' cannot be nil", name)
		}
	}
	return nil
}

// isArithmeticOperator checks if an operator is arithmetic
func isArithmeticOperator(op string) bool {
	switch op {
	case OpAdd, OpSub, OpMul, OpDiv, OpMod:
		return true
	default:
		return false
	}
}

// isComparisonOperator checks if an operator is a comparison operator
func isComparisonOperator(op string) bool {
	switch op {
	case OpEq, OpNeq, OpLt, OpGt, OpLte, OpGte:
		return true
	default:
		return false
	}
}

// recursiveConstraintValidator defines the interface for recursive constraint validation
type recursiveConstraintValidator func(program Program, constraint interface{}, expressionIndex int, depth int) error

// validateConstraintRecursive provides a generic recursive validation with depth tracking
// This function prevents stack overflow and applies a validator function recursively
func validateConstraintRecursive(
	program Program,
	constraint interface{},
	expressionIndex int,
	depth int,
	validator recursiveConstraintValidator,
) error {
	// Prevent stack overflow
	if depth > MaxValidationDepth {
		return fmt.Errorf("maximum validation depth exceeded (%d)", MaxValidationDepth)
	}

	constraintMap, ok := constraint.(map[string]interface{})
	if !ok {
		return nil // Not a map, nothing to validate
	}

	_, ok = constraintMap["type"].(string)
	if !ok {
		return nil // No type, skip validation
	}

	// Apply the specific validator
	return validator(program, constraint, expressionIndex, depth)
}

// validateOperands validates left and right operands recursively
func validateOperands(
	program Program,
	constraintMap map[string]interface{},
	expressionIndex int,
	depth int,
	validator recursiveConstraintValidator,
) error {
	if left := constraintMap["left"]; left != nil {
		if err := validator(program, left, expressionIndex, depth+1); err != nil {
			return err
		}
	}
	if right := constraintMap["right"]; right != nil {
		if err := validator(program, right, expressionIndex, depth+1); err != nil {
			return err
		}
	}
	return nil
}

// validateLogicalOperations validates operations in a logical expression
func validateLogicalOperations(
	program Program,
	constraintMap map[string]interface{},
	expressionIndex int,
	depth int,
	validator recursiveConstraintValidator,
) error {
	operations, ok := constraintMap["operations"].([]interface{})
	if !ok {
		return nil
	}

	for _, op := range operations {
		opMap, ok := op.(map[string]interface{})
		if !ok {
			continue
		}
		if right := opMap["right"]; right != nil {
			if err := validator(program, right, expressionIndex, depth+1); err != nil {
				return err
			}
		}
	}
	return nil
}

// safeBase64Decode safely decodes a base64 string with size limits
func safeBase64Decode(encoded string) (string, error) {
	// Check encoded string length first (base64 is ~4/3 of original)
	if len(encoded) > MaxBase64DecodeSize*4/3 {
		return "", fmt.Errorf("base64 encoded string exceeds maximum size")
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	// Check decoded size
	if len(decoded) > MaxBase64DecodeSize {
		return "", fmt.Errorf("base64 decoded data exceeds maximum size of %d bytes", MaxBase64DecodeSize)
	}

	// Validate UTF-8 encoding
	if !utf8.Valid(decoded) {
		return "", fmt.Errorf("base64 decoded data is not valid UTF-8")
	}

	return string(decoded), nil
}
