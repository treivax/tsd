package constraint

import (
	"encoding/base64"
	"testing"
)

func TestInferArgumentType_Coverage(t *testing.T) {
	// Setup ActionValidator with test types
	types := map[string]*TypeDefinition{
		"Person": {
			Name: "Person",
			Fields: []Field{
				{Name: "name", Type: "string"},
				{Name: "age", Type: "number"},
				{Name: "active", Type: "bool"},
			},
		},
		"Order": {
			Name: "Order",
			Fields: []Field{
				{Name: "id", Type: "string"},
				{Name: "total", Type: "number"},
			},
		},
	}

	av := &ActionValidator{types: types}

	ruleVars := map[string]string{
		"p":      "Person",
		"o":      "Order",
		"count":  "number",
		"status": "string",
		"flag":   "bool",
	}

	tests := []struct {
		name        string
		arg         interface{}
		ruleVars    map[string]string
		expectType  string
		expectError bool
	}{
		// Missing type field
		{
			name: "map without type field",
			arg: map[string]interface{}{
				"value": "test",
			},
			ruleVars:    ruleVars,
			expectError: true,
		},
		// Variable cases
		{
			name: "variable without name",
			arg: map[string]interface{}{
				"type": "variable",
			},
			ruleVars:    ruleVars,
			expectError: true,
		},
		{
			name: "variable not found in rule",
			arg: map[string]interface{}{
				"type": "variable",
				"name": "unknown",
			},
			ruleVars:    ruleVars,
			expectError: true,
		},
		{
			name: "valid variable",
			arg: map[string]interface{}{
				"type": "variable",
				"name": "count",
			},
			ruleVars:   ruleVars,
			expectType: "number",
		},
		// FieldAccess cases - missing object
		{
			name: "fieldAccess without object",
			arg: map[string]interface{}{
				"type":  "fieldAccess",
				"field": "name",
			},
			ruleVars:    ruleVars,
			expectError: true,
		},
		// FieldAccess cases - missing field
		{
			name: "fieldAccess without field",
			arg: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
			},
			ruleVars:    ruleVars,
			expectError: true,
		},
		// FieldAccess cases - object not found
		{
			name: "fieldAccess object not in rule",
			arg: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "unknownObj",
				"field":  "name",
			},
			ruleVars:    ruleVars,
			expectError: true,
		},
		// FieldAccess cases - type not found
		{
			name: "fieldAccess type definition not found",
			arg: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "name",
			},
			ruleVars: map[string]string{
				"p": "UnknownType",
			},
			expectError: true,
		},
		// FieldAccess cases - field not found in type
		{
			name: "fieldAccess field not in type",
			arg: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "unknownField",
			},
			ruleVars:    ruleVars,
			expectError: true,
		},
		// FieldAccess cases - valid
		{
			name: "valid fieldAccess string field",
			arg: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "name",
			},
			ruleVars:   ruleVars,
			expectType: "string",
		},
		{
			name: "valid fieldAccess number field",
			arg: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "age",
			},
			ruleVars:   ruleVars,
			expectType: "number",
		},
		{
			name: "valid fieldAccess bool field",
			arg: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "active",
			},
			ruleVars:   ruleVars,
			expectType: "bool",
		},
		// BinaryOp cases - missing operator
		{
			name: "binaryOp without operator",
			arg: map[string]interface{}{
				"type": "binaryOp",
			},
			ruleVars:    ruleVars,
			expectError: true,
		},
		// BinaryOp cases - arithmetic operators
		{
			name: "binaryOp addition",
			arg: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "+",
			},
			ruleVars:   ruleVars,
			expectType: "number",
		},
		{
			name: "binaryOp subtraction",
			arg: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "-",
			},
			ruleVars:   ruleVars,
			expectType: "number",
		},
		{
			name: "binaryOp multiplication",
			arg: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "*",
			},
			ruleVars:   ruleVars,
			expectType: "number",
		},
		{
			name: "binaryOp division",
			arg: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "/",
			},
			ruleVars:   ruleVars,
			expectType: "number",
		},
		{
			name: "binaryOp modulo",
			arg: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "%",
			},
			ruleVars:   ruleVars,
			expectType: "number",
		},
		// BinaryOp cases - comparison operators
		{
			name: "binaryOp equals",
			arg: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "==",
			},
			ruleVars:   ruleVars,
			expectType: "bool",
		},
		{
			name: "binaryOp not equals",
			arg: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "!=",
			},
			ruleVars:   ruleVars,
			expectType: "bool",
		},
		{
			name: "binaryOp less than",
			arg: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "<",
			},
			ruleVars:   ruleVars,
			expectType: "bool",
		},
		{
			name: "binaryOp greater than",
			arg: map[string]interface{}{
				"type":     "binaryOp",
				"operator": ">",
			},
			ruleVars:   ruleVars,
			expectType: "bool",
		},
		{
			name: "binaryOp less or equal",
			arg: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "<=",
			},
			ruleVars:   ruleVars,
			expectType: "bool",
		},
		{
			name: "binaryOp greater or equal",
			arg: map[string]interface{}{
				"type":     "binaryOp",
				"operator": ">=",
			},
			ruleVars:   ruleVars,
			expectType: "bool",
		},
		// BinaryOp cases - base64 encoded operator
		{
			name: "binaryOp base64 encoded addition",
			arg: map[string]interface{}{
				"type":     "binaryOp",
				"operator": base64.StdEncoding.EncodeToString([]byte("+")),
			},
			ruleVars:   ruleVars,
			expectType: "number",
		},
		{
			name: "binaryOp base64 encoded equals",
			arg: map[string]interface{}{
				"type":     "binaryOp",
				"operator": base64.StdEncoding.EncodeToString([]byte("==")),
			},
			ruleVars:   ruleVars,
			expectType: "bool",
		},
		// BinaryOp cases - unknown operator
		{
			name: "binaryOp unknown operator",
			arg: map[string]interface{}{
				"type":     "binaryOp",
				"operator": "??",
			},
			ruleVars:    ruleVars,
			expectError: true,
		},
		// Alternative type names for binaryOp
		{
			name: "binaryOperation type",
			arg: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "+",
			},
			ruleVars:   ruleVars,
			expectType: "number",
		},
		{
			name: "binary_operation type",
			arg: map[string]interface{}{
				"type":     "binary_operation",
				"operator": "*",
			},
			ruleVars:   ruleVars,
			expectType: "number",
		},
		// FunctionCall cases - missing name
		{
			name: "functionCall without name",
			arg: map[string]interface{}{
				"type": "functionCall",
			},
			ruleVars:    ruleVars,
			expectError: true,
		},
		// FunctionCall cases - various functions
		{
			name: "functionCall LENGTH",
			arg: map[string]interface{}{
				"type": "functionCall",
				"name": "LENGTH",
			},
			ruleVars:   ruleVars,
			expectType: "number",
		},
		{
			name: "functionCall SUBSTRING",
			arg: map[string]interface{}{
				"type": "functionCall",
				"name": "SUBSTRING",
			},
			ruleVars:   ruleVars,
			expectType: "string",
		},
		{
			name: "functionCall UPPER",
			arg: map[string]interface{}{
				"type": "functionCall",
				"name": "UPPER",
			},
			ruleVars:   ruleVars,
			expectType: "string",
		},
		{
			name: "functionCall LOWER",
			arg: map[string]interface{}{
				"type": "functionCall",
				"name": "LOWER",
			},
			ruleVars:   ruleVars,
			expectType: "string",
		},
		{
			name: "functionCall TRIM",
			arg: map[string]interface{}{
				"type": "functionCall",
				"name": "TRIM",
			},
			ruleVars:   ruleVars,
			expectType: "string",
		},
		{
			name: "functionCall ABS",
			arg: map[string]interface{}{
				"type": "functionCall",
				"name": "ABS",
			},
			ruleVars:   ruleVars,
			expectType: "number",
		},
		{
			name: "functionCall ROUND",
			arg: map[string]interface{}{
				"type": "functionCall",
				"name": "ROUND",
			},
			ruleVars:   ruleVars,
			expectType: "number",
		},
		{
			name: "functionCall FLOOR",
			arg: map[string]interface{}{
				"type": "functionCall",
				"name": "FLOOR",
			},
			ruleVars:   ruleVars,
			expectType: "number",
		},
		{
			name: "functionCall CEIL",
			arg: map[string]interface{}{
				"type": "functionCall",
				"name": "CEIL",
			},
			ruleVars:   ruleVars,
			expectType: "number",
		},
		{
			name: "functionCall unknown defaults to string",
			arg: map[string]interface{}{
				"type": "functionCall",
				"name": "CUSTOM_FUNC",
			},
			ruleVars:   ruleVars,
			expectType: "string",
		},
		// Unknown type
		{
			name: "unknown argument type",
			arg: map[string]interface{}{
				"type": "unknownType",
			},
			ruleVars:    ruleVars,
			expectError: true,
		},
		// Non-map argument
		{
			name:        "non-map argument",
			arg:         "plain string",
			ruleVars:    ruleVars,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, err := av.inferArgumentType(tt.arg, tt.ruleVars)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if gotType != tt.expectType {
					t.Errorf("expected type %q, got %q", tt.expectType, gotType)
				}
			}
		})
	}
}

func TestIsTypeCompatible_Coverage(t *testing.T) {
	av := &ActionValidator{}

	tests := []struct {
		name      string
		argType   string
		paramType string
		expected  bool
	}{
		{"exact match string", "string", "string", true},
		{"exact match number", "number", "number", true},
		{"exact match bool", "bool", "bool", true},
		{"mismatch string vs number", "string", "number", false},
		{"mismatch number vs bool", "number", "bool", false},
		{"mismatch bool vs string", "bool", "string", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := av.isTypeCompatible(tt.argType, tt.paramType)
			if result != tt.expected {
				t.Errorf("isTypeCompatible(%q, %q) = %v, want %v",
					tt.argType, tt.paramType, result, tt.expected)
			}
		})
	}
}

func TestInferFunctionReturnType(t *testing.T) {
	av := &ActionValidator{}

	tests := []struct {
		funcName string
		expected string
	}{
		{"LENGTH", "number"},
		{"length", "number"},
		{"SUBSTRING", "string"},
		{"substring", "string"},
		{"UPPER", "string"},
		{"upper", "string"},
		{"LOWER", "string"},
		{"lower", "string"},
		{"TRIM", "string"},
		{"trim", "string"},
		{"ABS", "number"},
		{"abs", "number"},
		{"ROUND", "number"},
		{"round", "number"},
		{"FLOOR", "number"},
		{"floor", "number"},
		{"CEIL", "number"},
		{"ceil", "number"},
		{"UNKNOWN_FUNC", "string"},
		{"customFunction", "string"},
	}

	for _, tt := range tests {
		t.Run(tt.funcName, func(t *testing.T) {
			result := av.inferFunctionReturnType(tt.funcName)
			if result != tt.expected {
				t.Errorf("inferFunctionReturnType(%q) = %q, want %q",
					tt.funcName, result, tt.expected)
			}
		})
	}
}

func TestIsTypeCompatible_WithUserDefinedTypes(t *testing.T) {
	types := map[string]*TypeDefinition{
		"Person": {
			Name: "Person",
			Fields: []Field{
				{Name: "name", Type: "string"},
			},
		},
		"Order": {
			Name: "Order",
			Fields: []Field{
				{Name: "id", Type: "string"},
			},
		},
	}

	av := &ActionValidator{types: types}

	tests := []struct {
		name      string
		argType   string
		paramType string
		expected  bool
	}{
		{"user type exact match", "Person", "Person", true},
		{"user type mismatch", "Person", "Order", false},
		{"user type vs primitive", "Person", "string", false},
		{"primitive vs user type", "string", "Person", false},
		{"user type exists in registry", "Order", "Order", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := av.isTypeCompatible(tt.argType, tt.paramType)
			if result != tt.expected {
				t.Errorf("isTypeCompatible(%q, %q) = %v, want %v",
					tt.argType, tt.paramType, result, tt.expected)
			}
		})
	}
}

func TestValidateActionDefinitions_Coverage(t *testing.T) {
	types := map[string]*TypeDefinition{
		"Person": {
			Name: "Person",
			Fields: []Field{
				{Name: "name", Type: "string"},
			},
		},
	}

	tests := []struct {
		name        string
		actions     map[string]*ActionDefinition
		expectError bool
	}{
		{
			name: "valid primitive parameter types",
			actions: map[string]*ActionDefinition{
				"testAction": {
					Name: "testAction",
					Parameters: []Parameter{
						{Name: "name", Type: "string"},
						{Name: "age", Type: "number"},
						{Name: "active", Type: "bool"},
					},
				},
			},
			expectError: false,
		},
		{
			name: "valid user-defined parameter type",
			actions: map[string]*ActionDefinition{
				"createPerson": {
					Name: "createPerson",
					Parameters: []Parameter{
						{Name: "person", Type: "Person"},
					},
				},
			},
			expectError: false,
		},
		{
			name: "invalid parameter type",
			actions: map[string]*ActionDefinition{
				"invalidAction": {
					Name: "invalidAction",
					Parameters: []Parameter{
						{Name: "param", Type: "UnknownType"},
					},
				},
			},
			expectError: true,
		},
		{
			name: "parameter with default value - type match",
			actions: map[string]*ActionDefinition{
				"actionWithDefault": {
					Name: "actionWithDefault",
					Parameters: []Parameter{
						{Name: "count", Type: "number", DefaultValue: float64(10)},
					},
				},
			},
			expectError: false,
		},
		{
			name: "parameter with default value - type mismatch",
			actions: map[string]*ActionDefinition{
				"actionBadDefault": {
					Name: "actionBadDefault",
					Parameters: []Parameter{
						{Name: "name", Type: "string", DefaultValue: 42},
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := &ActionValidator{types: types, actions: tt.actions}
			errors := av.ValidateActionDefinitions()

			if tt.expectError {
				if len(errors) == 0 {
					t.Errorf("expected error but got none")
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("unexpected errors: %v", errors)
				}
			}
		})
	}
}

func TestInferDefaultValueType_Coverage(t *testing.T) {
	av := &ActionValidator{}

	tests := []struct {
		name         string
		defaultValue interface{}
		expectedType string
	}{
		{"string value", "hello", "string"},
		{"float64 number", float64(42.5), "number"},
		{"int number", int(42), "number"},
		{"int64 number", int64(100), "number"},
		{"bool true", true, "bool"},
		{"bool false", false, "bool"},
		{"int32 unsupported", int32(10), "unknown"},
		{"float32 unsupported", float32(3.14), "unknown"},
		{"nil value", nil, "unknown"},
		{"map value", map[string]interface{}{"key": "value"}, "unknown"},
		{"slice value", []interface{}{1, 2, 3}, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := av.inferDefaultValueType(tt.defaultValue)
			if result != tt.expectedType {
				t.Errorf("inferDefaultValueType(%v) = %q, want %q",
					tt.defaultValue, result, tt.expectedType)
			}
		})
	}
}
