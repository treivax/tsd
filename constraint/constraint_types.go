package constraint

// Structures pour représenter l'AST
type Program struct {
	Types       []TypeDefinition `json:"types"`
	Expressions []Expression     `json:"expressions"`
}

type TypeDefinition struct {
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Expression struct {
	Type        string      `json:"type"`
	Set         Set         `json:"set"`
	Constraints interface{} `json:"constraints"`
	Action      *Action     `json:"action,omitempty"`
}

type Set struct {
	Type      string          `json:"type"`
	Variables []TypedVariable `json:"variables"`
}

type TypedVariable struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	DataType string `json:"dataType"`
}

type Constraint struct {
	Type     string      `json:"type"`
	Left     interface{} `json:"left,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Right    interface{} `json:"right,omitempty"`
}

type LogicalExpression struct {
	Type       string             `json:"type"`
	Left       interface{}        `json:"left"`
	Operations []LogicalOperation `json:"operations"`
}

type LogicalOperation struct {
	Op    string      `json:"op"`
	Right interface{} `json:"right"`
}

type BinaryOperation struct {
	Type     string      `json:"type"`
	Left     interface{} `json:"left"`
	Operator string      `json:"operator"`
	Right    interface{} `json:"right"`
}

type FieldAccess struct {
	Type   string `json:"type"`
	Object string `json:"object"`
	Field  string `json:"field"`
}

type Variable struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type NumberLiteral struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type StringLiteral struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type BooleanLiteral struct {
	Type  string `json:"type"`
	Value bool   `json:"value"`
}

// Structures pour les actions
type Action struct {
	Type string  `json:"type"`
	Job  JobCall `json:"job"`
}

type JobCall struct {
	Type string   `json:"type"`
	Name string   `json:"name"`
	Args []string `json:"args"`
}
