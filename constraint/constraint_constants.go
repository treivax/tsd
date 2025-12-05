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
	ValueTypeBoolean    = "boolean"
	ValueTypeBool       = "bool"
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
