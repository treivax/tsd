// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

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

// IsValidOperator checks if an operator string is a valid operator
func IsValidOperator(op string) bool {
	validOps := getValidOperators()
	return validOps[op]
}

// getValidOperators returns the set of valid operators
// This function creates the map on each call to ensure immutability
func getValidOperators() map[string]bool {
	return map[string]bool{
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

// IsValidPrimitiveType checks if a type string is a valid primitive type
func IsValidPrimitiveType(typeName string) bool {
	validTypes := getValidPrimitiveTypes()
	return validTypes[typeName]
}

// getValidPrimitiveTypes returns the set of valid primitive types
// This function creates the map on each call to ensure immutability
func getValidPrimitiveTypes() map[string]bool {
	return map[string]bool{
		ValueTypeString:  true,
		ValueTypeNumber:  true,
		ValueTypeBool:    true,
		ValueTypeBoolean: true,
		"integer":        true,
	}
}

// ValidOperators is deprecated: use IsValidOperator instead
// Kept for backward compatibility
var ValidOperators = getValidOperators()

// ValidPrimitiveTypes is deprecated: use IsValidPrimitiveType instead
// Kept for backward compatibility
var ValidPrimitiveTypes = getValidPrimitiveTypes()
