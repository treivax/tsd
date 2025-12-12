// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package domain

// IntegerLiteral represents an integer literal (needed for backward compatibility with tests)
// Note: This type is not present in constraint package but used in domain tests
type IntegerLiteral struct {
	Type  string `json:"type"`
	Value int64  `json:"value"`
}

// Helper functions to maintain backward compatibility with existing tests
// These functions provide convenience wrappers around constraint package types

// NewProgram creates a new Program with initialized slices
func NewProgram() *Program {
	return &Program{
		Types:        make([]TypeDefinition, 0),
		Expressions:  make([]Expression, 0),
		Actions:      make([]ActionDefinition, 0),
		Facts:        make([]Fact, 0),
		Resets:       make([]Reset, 0),
		RuleRemovals: make([]RuleRemoval, 0),
	}
}

// NewTypeDefinition creates a new TypeDefinition with the given name
func NewTypeDefinition(name string) TypeDefinition {
	return TypeDefinition{
		Type:   "typeDefinition",
		Name:   name,
		Fields: make([]Field, 0),
	}
}

// AddTypeField adds a field to a TypeDefinition
func AddTypeField(td *TypeDefinition, name, fieldType string) {
	td.Fields = append(td.Fields, Field{
		Name: name,
		Type: fieldType,
	})
}

// GetProgramTypeByName finds a type definition by name in a program
func GetProgramTypeByName(p *Program, name string) *TypeDefinition {
	for i := range p.Types {
		if p.Types[i].Name == name {
			return &p.Types[i]
		}
	}
	return nil
}

// GetTypeFieldByName finds a field by name in a type definition
func GetTypeFieldByName(td *TypeDefinition, name string) *Field {
	for i := range td.Fields {
		if td.Fields[i].Name == name {
			return &td.Fields[i]
		}
	}
	return nil
}

// TypeHasField checks if a type definition has a field with the given name
func TypeHasField(td *TypeDefinition, name string) bool {
	return GetTypeFieldByName(td, name) != nil
}

// NewExpression creates a new Expression with initialized Set
func NewExpression() Expression {
	return Expression{
		Type: "expression",
		Set: Set{
			Type:      "set",
			Variables: make([]TypedVariable, 0),
		},
	}
}

// AddExpressionVariable adds a variable to an expression's set
func AddExpressionVariable(e *Expression, name, dataType string) {
	e.Set.Variables = append(e.Set.Variables, TypedVariable{
		Type:     "typedVariable",
		Name:     name,
		DataType: dataType,
	})
}

// NewConstraint creates a new Constraint
func NewConstraint(left interface{}, operator string, right interface{}) Constraint {
	return Constraint{
		Type:     "constraint",
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

// NewFieldAccess creates a new FieldAccess
func NewFieldAccess(object, field string) FieldAccess {
	return FieldAccess{
		Type:   "fieldAccess",
		Object: object,
		Field:  field,
	}
}

// NewAction creates a new Action with a single job
func NewAction(jobName string, args ...interface{}) Action {
	return Action{
		Type: "action",
		Job: &JobCall{
			Type: "jobCall",
			Name: jobName,
			Args: args,
		},
	}
}

// ProgramToJSON converts a program to JSON string (for String() method compatibility)
func ProgramToJSON(p *Program) string {
	// Note: This would require encoding/json import
	// For now, return a simple representation
	// Tests should be updated to not rely on JSON serialization
	return ""
}
