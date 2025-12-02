// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"encoding/base64"
	"fmt"
)

// OperatorUtils provides utilities for working with operators in the RETE engine.
// The parser sometimes encodes operators as Base64 strings (e.g., "Kw==" for "+", "Kg==" for "*").
// These utilities centralize the decoding and validation logic.

// DecodeOperator attempts to decode a Base64-encoded operator string.
// If the input is not Base64-encoded, it returns the original string.
// This function is idempotent: decoding an already-decoded operator returns it unchanged.
//
// Examples:
//   - "Kw==" -> "+"
//   - "Kg==" -> "*"
//   - "+" -> "+"
func DecodeOperator(operator string) string {
	if operator == "" {
		return operator
	}

	// Try to decode from Base64
	if decoded, err := base64.StdEncoding.DecodeString(operator); err == nil {
		decodedStr := string(decoded)
		// Only return the decoded value if it looks like an operator
		// (to avoid false positives on strings that happen to be valid Base64)
		if IsValidOperator(decodedStr) {
			return decodedStr
		}
	}

	// If decoding failed or result is not a valid operator, return original
	return operator
}

// IsValidOperator checks if the given string is a known operator.
// This includes arithmetic, comparison, and logical operators.
func IsValidOperator(op string) bool {
	switch op {
	// Arithmetic operators
	case "+", "-", "*", "/", "%":
		return true
	// Comparison operators
	case "==", "!=", "<", "<=", ">", ">=":
		return true
	// String/collection operators
	case "CONTAINS", "IN", "LIKE", "MATCHES":
		return true
	// Logical operators
	case "AND", "OR", "NOT":
		return true
	default:
		return false
	}
}

// IsArithmeticOperator checks if the operator is arithmetic (+, -, *, /, %).
func IsArithmeticOperator(op string) bool {
	switch op {
	case "+", "-", "*", "/", "%":
		return true
	default:
		return false
	}
}

// IsComparisonOperator checks if the operator is a comparison operator.
func IsComparisonOperator(op string) bool {
	switch op {
	case "==", "!=", "<", "<=", ">", ">=":
		return true
	default:
		return false
	}
}

// IsStringOperator checks if the operator works on strings/collections.
func IsStringOperator(op string) bool {
	switch op {
	case "CONTAINS", "IN", "LIKE", "MATCHES":
		return true
	default:
		return false
	}
}

// IsLogicalOperator checks if the operator is a logical operator.
func IsLogicalOperator(op string) bool {
	switch op {
	case "AND", "OR", "NOT":
		return true
	default:
		return false
	}
}

// NormalizeOperator decodes and validates an operator, returning an error if invalid.
func NormalizeOperator(operator string) (string, error) {
	decoded := DecodeOperator(operator)
	if !IsValidOperator(decoded) {
		return "", fmt.Errorf("invalid or unsupported operator: '%s' (decoded: '%s')", operator, decoded)
	}
	return decoded, nil
}

// ExtractOperatorFromMap extracts an operator from a map[string]interface{} value.
// This handles the various ways operators can be represented in parsed expressions:
//   - As a string: map["operator"] = "+"
//   - As a byte array: map["operator"] = []uint8{'+'}
//   - As Base64: map["operator"] = "Kw=="
//
// Returns the decoded operator string and any error encountered.
func ExtractOperatorFromMap(m map[string]interface{}) (string, error) {
	opVal, ok := m["operator"]
	if !ok {
		return "", fmt.Errorf("operator key not found in map")
	}

	var operator string

	switch op := opVal.(type) {
	case string:
		operator = op
	case []byte:
		operator = string(op)
	default:
		return "", fmt.Errorf("unsupported operator type: %T (value: %+v)", opVal, opVal)
	}

	// Decode if Base64-encoded
	return DecodeOperator(operator), nil
}
