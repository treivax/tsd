// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

// Program represents the complete AST of a constraint program including types, expressions and facts.
// It serves as the root structure for parsed constraint files.
type Program struct {
	Types       []TypeDefinition `json:"types"`       // Type definitions declared in the program
	Expressions []Expression     `json:"expressions"` // Constraint expressions/rules
	Facts       []Fact           `json:"facts"`       // Facts parsed from the program
}

// TypeDefinition represents a user-defined type with its fields.
// Example: type Person : <id: string, name: string, age: number>
type TypeDefinition struct {
	Type   string  `json:"type"`   // Always "typeDefinition"
	Name   string  `json:"name"`   // The type name (e.g., "Person")
	Fields []Field `json:"fields"` // List of fields in the type
}

// Field represents a single field within a type definition.
// It contains the field name and its type.
type Field struct {
	Name string `json:"name"` // Field name (e.g., "id", "name")
	Type string `json:"type"` // Field type (e.g., "string", "number", "bool")
}

// Expression represents a constraint expression or rule in the system.
// It defines variables, constraints on those variables, and actions to execute when matched.
type Expression struct {
	Type        string      `json:"type"`             // Always "expression"
	Set         Set         `json:"set"`              // Set of variables used in the expression
	Constraints interface{} `json:"constraints"`      // Constraints to evaluate
	Action      *Action     `json:"action,omitempty"` // Action to execute when constraints match
}

// Set represents a collection of typed variables used in an expression.
// It defines the scope of variables that constraints can reference.
type Set struct {
	Type      string          `json:"type"`      // Always "set"
	Variables []TypedVariable `json:"variables"` // List of variables in the set
}

// TypedVariable represents a variable with its associated type.
// Example: p: Person where 'p' is the name and 'Person' is the dataType.
type TypedVariable struct {
	Type     string `json:"type"`     // Always "typedVariable"
	Name     string `json:"name"`     // Variable name (e.g., "p", "order")
	DataType string `json:"dataType"` // Associated type (e.g., "Person", "Order")
}

// Constraint represents a basic constraint in the system.
// It can be a comparison, field access, or other constraint type.
type Constraint struct {
	Type     string      `json:"type"`               // Type of constraint
	Left     interface{} `json:"left,omitempty"`     // Left operand
	Operator string      `json:"operator,omitempty"` // Comparison operator (==, !=, <, >, etc.)
	Right    interface{} `json:"right,omitempty"`    // Right operand
}

// LogicalExpression represents a complex logical expression with AND/OR operations.
// It consists of a left operand and a series of logical operations.
type LogicalExpression struct {
	Type       string             `json:"type"`       // Always "logicalExpr"
	Left       interface{}        `json:"left"`       // Left operand of the expression
	Operations []LogicalOperation `json:"operations"` // Chain of AND/OR operations
}

// LogicalOperation represents a single logical operation (AND/OR) in a chain.
// It combines with the previous expression using the specified operator.
type LogicalOperation struct {
	Op    string      `json:"op"`    // Logical operator ("AND" or "OR")
	Right interface{} `json:"right"` // Right operand for this operation
}

// BinaryOperation represents a binary operation between two operands.
// Common operations include arithmetic (+, -, *, /) and comparisons (==, !=, <, >).
type BinaryOperation struct {
	Type     string      `json:"type"`     // Always "binaryOperation" or "comparison"
	Left     interface{} `json:"left"`     // Left operand
	Operator string      `json:"operator"` // Operation symbol
	Right    interface{} `json:"right"`    // Right operand
}

// FieldAccess represents accessing a field of an object/variable.
// Example: p.age where 'p' is the object and 'age' is the field.
type FieldAccess struct {
	Type   string `json:"type"`   // Always "fieldAccess"
	Object string `json:"object"` // Variable name (e.g., "p")
	Field  string `json:"field"`  // Field name (e.g., "age", "name")
}

// Variable represents a simple variable reference in expressions.
// It's used when referencing a variable by name.
type Variable struct {
	Type string `json:"type"` // Always "variable"
	Name string `json:"name"` // Variable name
}

// NumberLiteral represents a numeric literal value in expressions.
// It supports both integer and floating-point numbers.
type NumberLiteral struct {
	Type  string  `json:"type"`  // Always "numberLiteral" or "number"
	Value float64 `json:"value"` // Numeric value
}

// StringLiteral represents a string literal value in expressions.
// It contains text values enclosed in quotes.
type StringLiteral struct {
	Type  string `json:"type"`  // Always "stringLiteral" or "string"
	Value string `json:"value"` // String content
}

// BooleanLiteral represents a boolean literal value (true/false).
// It's used for boolean constants in expressions.
type BooleanLiteral struct {
	Type  string `json:"type"`  // Always "booleanLiteral" or "bool"
	Value bool   `json:"value"` // Boolean value
}

// NotConstraint represents a negation constraint (NOT operator).
// It negates the result of the contained expression.
type NotConstraint struct {
	Type       string      `json:"type"`       // Always "notConstraint"
	Expression interface{} `json:"expression"` // Expression to negate
}

// ExistsConstraint represents an existential quantifier (EXISTS).
// It checks if there exists at least one instance satisfying the condition.
type ExistsConstraint struct {
	Type      string        `json:"type"`      // Always "existsConstraint"
	Variable  TypedVariable `json:"variable"`  // Variable to quantify over
	Condition interface{}   `json:"condition"` // Condition that must be satisfied
}

// AggregateConstraint represents aggregate operations (SUM, COUNT, AVG, MIN, MAX).
// It performs calculations over sets of data and compares the result.
type AggregateConstraint struct {
	Type       string      `json:"type"`       // Always "aggregateConstraint"
	Function   string      `json:"function"`   // Aggregate function (SUM, COUNT, AVG, MIN, MAX)
	Expression interface{} `json:"expression"` // Expression to aggregate
	Operator   string      `json:"operator"`   // Comparison operator
	Value      interface{} `json:"value"`      // Value to compare against
}

// FunctionCall represents a function call in expressions.
// It includes the function name and its arguments.
type FunctionCall struct {
	Type string        `json:"type"` // Always "functionCall"
	Name string        `json:"name"` // Function name (e.g., "LENGTH", "UPPER")
	Args []interface{} `json:"args"` // Function arguments
}

// ArrayLiteral represents an array/list literal in expressions.
// It contains a collection of elements of potentially different types.
type ArrayLiteral struct {
	Type     string        `json:"type"`     // Always "arrayLiteral"
	Elements []interface{} `json:"elements"` // Array elements
}

// Action represents an action to execute when constraints are satisfied.
// It defines what job should be performed and with what parameters.
type Action struct {
	Type string  `json:"type"` // Always "action"
	Job  JobCall `json:"job"`  // Job to execute
}

// JobCall represents a specific job/function call within an action.
// It specifies the job name and arguments to pass.
type JobCall struct {
	Type string        `json:"type"` // Always "jobCall"
	Name string        `json:"name"` // Job/function name
	Args []interface{} `json:"args"` // Arguments to pass to the job
}

// Fact represents a fact parsed from constraint files.
// Facts are data instances that can be inserted into the RETE network.
// Example: Person(id: "P001", name: "Alice", age: 25)
type Fact struct {
	Type     string      `json:"type"`     // Always "fact"
	TypeName string      `json:"typeName"` // Type name (e.g., "Person", "Order")
	Fields   []FactField `json:"fields"`   // List of field assignments
}

// FactField represents a field assignment within a fact.
// It pairs a field name with its assigned value.
type FactField struct {
	Name  string    `json:"name"`  // Field name (e.g., "id", "name")
	Value FactValue `json:"value"` // Assigned value
}

// FactValue represents a value assigned to a fact field.
// It wraps the actual value with type information.
type FactValue struct {
	Type  string      `json:"type"`  // Value type ("string", "number", "bool")
	Value interface{} `json:"value"` // Actual value
}
