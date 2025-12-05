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

	// Handle commutative operators
	if operator, hasOp := condition["operator"].(string); hasOp {
		normalized["operator"] = operator

		if operator == "==" || operator == "!=" {
			// For commutative operators, sort left/right by lexicographic order
			left, hasLeft := condition["left"]
			right, hasRight := condition["right"]

			if hasLeft && hasRight {
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
		} else {
			// Non-commutative operators: preserve order
			if left, hasLeft := condition["left"]; hasLeft {
				normalized["left"] = left
			}
			if right, hasRight := condition["right"]; hasRight {
				normalized["right"] = right
			}
		}
	}

	// Copy remaining fields
	for key, value := range condition {
		if key != "type" && key != "operator" && key != "left" && key != "right" {
			normalized[key] = value
		}
	}

	return normalized, nil
}

// ComputeJoinHash computes a hash for a join specification.
// This is a standalone helper function for backward compatibility.
func ComputeJoinHash(condition map[string]interface{}, leftVars, rightVars []string, varTypes map[string]string) (string, error) {
	// Create canonical signature
	canonical := &CanonicalJoinSignature{
		Version:   "1.0",
		LeftVars:  sortStrings(leftVars),
		RightVars: sortStrings(rightVars),
		AllVars:   sortStrings(append(leftVars, rightVars...)),
		VarTypes:  sortVarTypes(varTypes),
		Condition: condition,
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
