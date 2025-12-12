// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package domain

import (
	"github.com/treivax/tsd/constraint"
)

// Type aliases to eliminate duplication with constraint_types.go
// The canonical definitions are in the constraint package root.
// This allows internal packages like validator to use domain.TypeDefinition
// while avoiding code duplication.

type (
	// Core AST types - Aliases to constraint package
	Program        = constraint.Program
	TypeDefinition = constraint.TypeDefinition
	Field          = constraint.Field
	Expression     = constraint.Expression
	Set            = constraint.Set
	TypedVariable  = constraint.TypedVariable

	// Constraint types
	Constraint          = constraint.Constraint
	LogicalExpression   = constraint.LogicalExpression
	LogicalOperation    = constraint.LogicalOperation
	BinaryOperation     = constraint.BinaryOperation
	NotConstraint       = constraint.NotConstraint
	ExistsConstraint    = constraint.ExistsConstraint
	AggregateConstraint = constraint.AggregateConstraint
	FunctionCall        = constraint.FunctionCall
	ArrayLiteral        = constraint.ArrayLiteral

	// Value types
	FieldAccess    = constraint.FieldAccess
	Variable       = constraint.Variable
	NumberLiteral  = constraint.NumberLiteral
	StringLiteral  = constraint.StringLiteral
	BooleanLiteral = constraint.BooleanLiteral

	// Action types
	Action  = constraint.Action
	JobCall = constraint.JobCall

	// Fact types
	Fact      = constraint.Fact
	FactField = constraint.FactField
	FactValue = constraint.FactValue

	// Command types
	Reset            = constraint.Reset
	RuleRemoval      = constraint.RuleRemoval
	ActionDefinition = constraint.ActionDefinition
	Parameter        = constraint.Parameter
)

// IsValidOperator vérifie si un opérateur est valide
func IsValidOperator(op string) bool {
	return constraint.ValidOperators[op]
}

// IsValidType vérifie si un type primitif est valide
func IsValidType(t string) bool {
	return constraint.ValidPrimitiveTypes[t]
}
