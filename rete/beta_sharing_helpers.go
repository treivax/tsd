// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"encoding/json"
)

// ============================================================================
// Helper Functions (Standalone backward-compatible functions)
// ============================================================================

// NormalizeJoinCondition normalizes a join condition for consistent hashing.
// This is a standalone helper function for backward compatibility.
func NormalizeJoinCondition(condition map[string]interface{}) (map[string]interface{}, error) {
	if condition == nil {
		return map[string]interface{}{"type": "simple"}, nil
	}

	normalized := make(map[string]interface{})

	// Normalize type field
	normalizeTypeField(condition, normalized)

	// Handle operator and operands
	if operator, hasOp := condition["operator"].(string); hasOp {
		normalizeOperatorAndOperands(condition, normalized, operator)
	}

	// Copy remaining fields
	copyRemainingFields(condition, normalized)

	return normalized, nil
}

// normalizeTypeField normalizes the type field of a condition.
func normalizeTypeField(condition, normalized map[string]interface{}) {
	if condType, hasType := condition["type"]; hasType {
		if typeStr, ok := condType.(string); ok {
			switch typeStr {
			case "comparison", "binaryOperation":
				normalized["type"] = "binaryOperation"
			default:
				normalized["type"] = typeStr
			}
		} else {
			normalized["type"] = condType
		}
	}
}

// normalizeOperatorAndOperands handles operator normalization and operand ordering.
func normalizeOperatorAndOperands(condition, normalized map[string]interface{}, operator string) {
	normalized["operator"] = operator

	if isCommutativeOperator(operator) {
		normalizeCommutativeOperands(condition, normalized)
	} else {
		normalizeNonCommutativeOperands(condition, normalized)
	}
}

// isCommutativeOperator checks if an operator is commutative.
func isCommutativeOperator(operator string) bool {
	return operator == "==" || operator == "!="
}

// normalizeCommutativeOperands sorts operands for commutative operators.
func normalizeCommutativeOperands(condition, normalized map[string]interface{}) {
	left, hasLeft := condition["left"]
	right, hasRight := condition["right"]

	if !hasLeft || !hasRight {
		return
	}

	leftJSON, _ := json.Marshal(left)
	rightJSON, _ := json.Marshal(right)

	if string(leftJSON) > string(rightJSON) {
		normalized["left"] = right
		normalized["right"] = left
	} else {
		normalized["left"] = left
		normalized["right"] = right
	}
}

// normalizeNonCommutativeOperands preserves operand order for non-commutative operators.
func normalizeNonCommutativeOperands(condition, normalized map[string]interface{}) {
	if left, hasLeft := condition["left"]; hasLeft {
		normalized["left"] = left
	}
	if right, hasRight := condition["right"]; hasRight {
		normalized["right"] = right
	}
}

// copyRemainingFields copies all fields except already processed ones.
func copyRemainingFields(condition, normalized map[string]interface{}) {
	const (
		typeField     = "type"
		operatorField = "operator"
		leftField     = "left"
		rightField    = "right"
	)

	for key, value := range condition {
		if key != typeField && key != operatorField && key != leftField && key != rightField {
			normalized[key] = value
		}
	}
}

// ComputeJoinHash computes a hash for a join specification.
// This is a standalone helper function for backward compatibility.
//
// Deprecated: Use BetaSharingRegistry.GetOrCreateJoinNode instead which includes
// cascadeLevel for proper sharing semantics.
func ComputeJoinHash(condition map[string]interface{}, leftVars, rightVars []string, varTypes map[string]string) (string, error) {
	// Create canonical signature with default values for backward compatibility
	canonical := &CanonicalJoinSignature{
		Version:      "1.0",
		LeftVars:     sortStrings(leftVars),
		RightVars:    sortStrings(rightVars),
		AllVars:      sortStrings(append(leftVars, rightVars...)),
		VarTypes:     sortVarTypes(varTypes),
		Condition:    condition,
		CascadeLevel: 0, // Default cascade level for backward compatibility
	}

	// Normalize condition
	if condition != nil {
		normalizedCond, err := NormalizeJoinCondition(condition)
		if err != nil {
			return "", err
		}
		canonical.Condition = normalizedCond
	}

	return ComputeJoinNodeHash(canonical)
}
