// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import "sync"

// Constraint type constants define the different types of constraints
// supported in the constraint system.
const (
	ConstraintTypeFieldAccess = "fieldAccess"
	ConstraintTypeComparison  = "comparison"
	ConstraintTypeLogicalExpr = "logicalExpr"
	ConstraintTypeBinaryOp    = "binaryOp"
)

// Value type constants define the different types of values
// that can be used in constraints and facts.
const (
	ValueTypeString     = "string"
	ValueTypeNumber     = "number"
	ValueTypeBoolean    = "boolean" // Primary format
	ValueTypeBool       = "bool"    // Alias for backward compatibility
	ValueTypeIdentifier = "identifier"
	ValueTypeVariable   = "variable"
	ValueTypeUnknown    = "unknown"
)

// Special field name constants are reserved field names used internally
// by the RETE network for identification and typing.
const (
	FieldNameID       = "id"
	FieldNameReteType = "reteType"
)

// JSON key constants define the keys used in JSON serialization and parsing
const (
	JSONKeyType         = "type"
	JSONKeyFieldAccess  = "fieldAccess"
	JSONKeyObject       = "object"
	JSONKeyField        = "field"
	JSONKeyTypes        = "types"
	JSONKeyActions      = "actions"
	JSONKeyExpressions  = "expressions"
	JSONKeyRuleRemovals = "ruleRemovals"
)

// Argument type constants for action parameter validation
const (
	ArgTypeStringLiteral = "stringLiteral"
	ArgTypeNumberLiteral = "numberLiteral"
	ArgTypeBoolLiteral   = "booleanLiteral"
	ArgTypeFunctionCall  = "functionCall"
)

// Binary operation type constants
// Multiple variants exist due to parser evolution and backward compatibility
const (
	ArgTypeBinaryOp = "binaryOp" // Primary format
)

// Legacy binary operation type aliases for backward compatibility
// These are kept to support older parsed ASTs
const (
	ArgTypeBinaryOp2 = "binaryOperation"
	ArgTypeBinaryOp3 = "binary_operation"
)

// isBinaryOperationType checks if a type string represents a binary operation
// Handles all variants for backward compatibility
func isBinaryOperationType(t string) bool {
	return t == ArgTypeBinaryOp || t == ArgTypeBinaryOp2 || t == ArgTypeBinaryOp3
}

// Operator constants for binary operations
const (
	OpAdd = "+"
	OpSub = "-"
	OpMul = "*"
	OpDiv = "/"
	OpMod = "%"
	OpEq  = "=="
	OpNeq = "!="
	OpLt  = "<"
	OpGt  = ">"
	OpLte = "<="
	OpGte = ">="
)

// Logical operator constants
const (
	OpAnd = "AND"
	OpOr  = "OR"
	OpNot = "NOT"
)

// Validation limits
const (
	// MaxValidationDepth is the maximum recursion depth for constraint validation
	// to prevent stack overflow attacks with deeply nested structures
	MaxValidationDepth = 100

	// MaxBase64DecodeSize is the maximum size in bytes for base64 decoded strings
	// to prevent DoS attacks with large payloads
	MaxBase64DecodeSize = 1024 * 1024 // 1MB
)

// validOperatorsMap holds the set of valid operators.
// Initialized once using sync.Once for thread-safe lazy loading.
var (
	validOperatorsMap     map[string]bool
	validOperatorsMapOnce sync.Once
)

// validPrimitiveTypesMap holds the set of valid primitive types.
// Initialized once using sync.Once for thread-safe lazy loading.
var (
	validPrimitiveTypesMap     map[string]bool
	validPrimitiveTypesMapOnce sync.Once
)

// initValidOperatorsMap initializes the validOperatorsMap exactly once.
func initValidOperatorsMap() {
	validOperatorsMap = map[string]bool{
		OpEq:  true,
		OpNeq: true,
		OpLt:  true,
		OpGt:  true,
		OpLte: true,
		OpGte: true,
		OpAnd: true,
		OpOr:  true,
		OpNot: true,
		OpAdd: true,
		OpSub: true,
		OpMul: true,
		OpDiv: true,
		OpMod: true,
	}
}

// initValidPrimitiveTypesMap initializes the validPrimitiveTypesMap exactly once.
func initValidPrimitiveTypesMap() {
	validPrimitiveTypesMap = map[string]bool{
		ValueTypeString:  true,
		ValueTypeNumber:  true,
		ValueTypeBool:    true,
		ValueTypeBoolean: true,
		"integer":        true,
	}
}

// IsValidOperator checks if an operator string is a valid operator.
// Thread-safe and efficient - the map is initialized only once.
func IsValidOperator(op string) bool {
	validOperatorsMapOnce.Do(initValidOperatorsMap)
	return validOperatorsMap[op]
}

// IsValidPrimitiveType checks if a type string is a valid primitive type.
// Thread-safe and efficient - the map is initialized only once.
func IsValidPrimitiveType(typeName string) bool {
	validPrimitiveTypesMapOnce.Do(initValidPrimitiveTypesMap)
	return validPrimitiveTypesMap[typeName]
}

// getValidOperators returns the set of valid operators.
// Deprecated: This function is kept for backward compatibility.
// Prefer using IsValidOperator for individual checks.
func getValidOperators() map[string]bool {
	validOperatorsMapOnce.Do(initValidOperatorsMap)
	// Return a copy to prevent external modification
	result := make(map[string]bool, len(validOperatorsMap))
	for k, v := range validOperatorsMap {
		result[k] = v
	}
	return result
}

// getValidPrimitiveTypes returns the set of valid primitive types.
// Deprecated: This function is kept for backward compatibility.
// Prefer using IsValidPrimitiveType for individual checks.
func getValidPrimitiveTypes() map[string]bool {
	validPrimitiveTypesMapOnce.Do(initValidPrimitiveTypesMap)
	// Return a copy to prevent external modification
	result := make(map[string]bool, len(validPrimitiveTypesMap))
	for k, v := range validPrimitiveTypesMap {
		result[k] = v
	}
	return result
}

// ValidOperators is deprecated: use IsValidOperator instead
// Kept for backward compatibility
var ValidOperators = getValidOperators()

// ValidPrimitiveTypes is deprecated: use IsValidPrimitiveType instead
// Kept for backward compatibility
var ValidPrimitiveTypes = getValidPrimitiveTypes()
