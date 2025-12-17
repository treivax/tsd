// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

// Program represents the complete AST of a constraint program including types, actions, expressions, facts and resets.
// It serves as the root structure for parsed constraint files.
type Program struct {
	Types        []TypeDefinition   `json:"types"`        // Type definitions declared in the program
	Actions      []ActionDefinition `json:"actions"`      // Action definitions with their signatures
	Expressions  []Expression       `json:"expressions"`  // Constraint expressions/rules
	Facts        []Fact             `json:"facts"`        // Facts parsed from the program
	Resets       []Reset            `json:"resets"`       // Reset instructions to clear the system
	RuleRemovals []RuleRemoval      `json:"ruleRemovals"` // Rule removal commands
}

// TypeDefinition represents a user-defined type with its fields.
// Example: type Person(#id: string, name: string, age: number)
type TypeDefinition struct {
	Type   string  `json:"type"`   // Always "typeDefinition"
	Name   string  `json:"name"`   // The type name (e.g., "Person")
	Fields []Field `json:"fields"` // List of fields in the type
}

// Field represents a single field within a type definition.
// It contains the field name, its type, and whether it's part of the primary key.
type Field struct {
	Name         string `json:"name"`                   // Field name (e.g., "id", "name")
	Type         string `json:"type"`                   // Field type (e.g., "string", "number", "bool")
	IsPrimaryKey bool   `json:"isPrimaryKey,omitempty"` // True if field is part of primary key (marked with #)
}

// ActionDefinition represents a user-defined action with its signature.
// Example: action notify(recipient: string, message: string, priority: number = 1)
type ActionDefinition struct {
	Type       string      `json:"type"`       // Always "actionDefinition"
	Name       string      `json:"name"`       // The action name (e.g., "notify")
	Parameters []Parameter `json:"parameters"` // List of parameters for the action
}

// Parameter represents a single parameter within an action definition.
// It contains the parameter name, type, whether it's optional, and an optional default value.
type Parameter struct {
	Name         string      `json:"name"`                   // Parameter name (e.g., "recipient", "priority")
	Type         string      `json:"type"`                   // Parameter type (e.g., "string", "number", "bool", or a user-defined type like "Person")
	Optional     bool        `json:"optional"`               // Whether the parameter is optional (marked with ?)
	DefaultValue interface{} `json:"defaultValue,omitempty"` // Default value if provided
}

// Expression represents a constraint expression or rule in the system.
// It defines variables, constraints on those variables, and actions to execute when matched.
// Each expression must have a unique identifier for management purposes (e.g., deletion).
type Expression struct {
	Type        string      `json:"type"`               // Always "expression"
	RuleId      string      `json:"ruleId"`             // Unique identifier for the rule
	Set         Set         `json:"set,omitempty"`      // Set of variables (single pattern, backward compatibility)
	Patterns    []Set       `json:"patterns,omitempty"` // Multiple pattern blocks (aggregation with joins)
	Constraints interface{} `json:"constraints"`        // Constraints to evaluate
	Action      *Action     `json:"action,omitempty"`   // Action to execute when constraints match
}

// Set represents a collection of typed variables used in an expression.
// It defines the scope of variables that constraints can reference.
type Set struct {
	Type      string          `json:"type"`      // Always "set"
	Variables []TypedVariable `json:"variables"` // List of variables in the set
}

// TypedVariable represents a variable with its associated type.
// Example: p: Person where 'p' is the name and 'Person' is the dataType.
// For aggregation variables (type="aggregationVariable"), it also includes Function and Field.
type TypedVariable struct {
	Type     string      `json:"type"`               // "typedVariable" or "aggregationVariable"
	Name     string      `json:"name"`               // Variable name (e.g., "p", "order", "avg_sal")
	DataType string      `json:"dataType"`           // Associated type (e.g., "Person", "Order")
	Function string      `json:"function,omitempty"` // Aggregation function (e.g., "AVG", "COUNT", "SUM") for aggregation variables
	Field    interface{} `json:"field,omitempty"`    // Field being aggregated (map with object/field/type) for aggregation variables
	Value    interface{} `json:"value,omitempty"`    // Optional value field for complex variable definitions
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
// It defines what job(s) should be performed and with what parameters.
// Supports both single action (Job field, for backward compatibility) and
// multiple actions (Jobs field, new format).
type Action struct {
	Type string    `json:"type"`           // Always "action"
	Job  *JobCall  `json:"job,omitempty"`  // Single job (backward compatibility)
	Jobs []JobCall `json:"jobs,omitempty"` // Multiple jobs (new format)
}

// GetJobs returns the list of jobs to execute.
// It handles both the old format (single Job) and new format (multiple Jobs).
func (a *Action) GetJobs() []JobCall {
	if len(a.Jobs) > 0 {
		return a.Jobs
	}
	if a.Job != nil {
		return []JobCall{*a.Job}
	}
	return []JobCall{}
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

// BuildFieldMap creates a map of field values indexed by field name.
// This is useful for efficient field lookup when validating or processing facts.
func (f Fact) BuildFieldMap() map[string]FactValue {
	fieldMap := make(map[string]FactValue, len(f.Fields))
	for _, field := range f.Fields {
		fieldMap[field.Name] = field.Value
	}
	return fieldMap
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

// Unwrap extracts the underlying value from a FactValue.
// Handles nested map structures from parser output where the value
// might be wrapped in a map with a "value" key.
// Returns the unwrapped value ready for use in the RETE network.
func (fv FactValue) Unwrap() interface{} {
	if valMap, ok := fv.Value.(map[string]interface{}); ok {
		if val, exists := valMap["value"]; exists {
			return val
		}
	}
	return fv.Value
}

// Reset represents a reset instruction that clears the entire system.
// When executed, it removes all facts, rules, types, and the RETE network.
// Example: reset
type Reset struct {
	Type string `json:"type"` // Always "reset"
}

// RuleRemoval represents a command to remove a rule from the system.
// Example: remove rule my_rule
type RuleRemoval struct {
	Type   string `json:"type"`   // Always "ruleRemoval"
	RuleID string `json:"ruleID"` // ID of the rule to remove
}

// GetPrimaryKeyFields returns the list of fields that are part of the primary key.
// Fields are returned in the order they appear in the type definition.
func (td TypeDefinition) GetPrimaryKeyFields() []Field {
	var pkFields []Field
	for _, field := range td.Fields {
		if field.IsPrimaryKey {
			pkFields = append(pkFields, field)
		}
	}
	return pkFields
}

// HasPrimaryKey returns true if the type has at least one primary key field.
func (td TypeDefinition) HasPrimaryKey() bool {
	for _, field := range td.Fields {
		if field.IsPrimaryKey {
			return true
		}
	}
	return false
}

// GetPrimaryKeyFieldNames returns the names of primary key fields in definition order.
func (td TypeDefinition) GetPrimaryKeyFieldNames() []string {
	var names []string
	for _, field := range td.Fields {
		if field.IsPrimaryKey {
			names = append(names, field.Name)
		}
	}
	return names
}

// Clone creates a deep copy of the TypeDefinition.
// All fields including IsPrimaryKey are copied.
func (td TypeDefinition) Clone() TypeDefinition {
	clone := TypeDefinition{
		Type:   td.Type,
		Name:   td.Name,
		Fields: make([]Field, len(td.Fields)),
	}

	// Copy all fields (copy() copies all struct fields including IsPrimaryKey)
	copy(clone.Fields, td.Fields)

	return clone
}
