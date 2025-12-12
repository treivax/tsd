package constraint

import (
	"testing"
)

func TestValidateFieldAccess_Coverage(t *testing.T) {
	// Create a program with types and expressions
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "name", Type: "string"},
					{Name: "age", Type: "number"},
					{Name: "active", Type: "bool"},
				},
			},
			{
				Name: "Order",
				Fields: []Field{
					{Name: "id", Type: "string"},
					{Name: "total", Type: "number"},
				},
			},
		},
		Expressions: []Expression{
			{
				// Old single-pattern syntax
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
						{Name: "count", DataType: "number"},
					},
				},
			},
			{
				// New multi-pattern syntax
				Patterns: []Set{
					{
						Variables: []TypedVariable{
							{Name: "o", DataType: "Order"},
							{Name: "status", DataType: "string"},
						},
					},
				},
			},
			{
				// Multi-pattern with multiple patterns
				Patterns: []Set{
					{
						Variables: []TypedVariable{
							{Name: "x", DataType: "Person"},
						},
					},
					{
						Variables: []TypedVariable{
							{Name: "y", DataType: "Order"},
						},
					},
				},
			},
			{
				// Empty expression
				Patterns: []Set{},
			},
		},
	}

	tests := []struct {
		name            string
		fieldAccess     FieldAccess
		expressionIndex int
		expectError     bool
		errorContains   string
	}{
		{
			name: "invalid expression index",
			fieldAccess: FieldAccess{
				Object: "p",
				Field:  "name",
			},
			expressionIndex: 99,
			expectError:     true,
			errorContains:   "invalid expression index",
		},
		{
			name: "variable not found in old syntax",
			fieldAccess: FieldAccess{
				Object: "unknown",
				Field:  "name",
			},
			expressionIndex: 0,
			expectError:     true,
			errorContains:   "not found",
		},
		{
			name: "variable not found in new syntax",
			fieldAccess: FieldAccess{
				Object: "unknown",
				Field:  "id",
			},
			expressionIndex: 1,
			expectError:     true,
			errorContains:   "not found",
		},
		{
			name: "type not found",
			fieldAccess: FieldAccess{
				Object: "count",
				Field:  "value",
			},
			expressionIndex: 0,
			expectError:     true,
			errorContains:   "type not found: number",
		},
		{
			name: "field not found in type",
			fieldAccess: FieldAccess{
				Object: "p",
				Field:  "unknownField",
			},
			expressionIndex: 0,
			expectError:     true,
			errorContains:   "field unknownField not found in type Person",
		},
		{
			name: "valid field access old syntax",
			fieldAccess: FieldAccess{
				Object: "p",
				Field:  "name",
			},
			expressionIndex: 0,
			expectError:     false,
		},
		{
			name: "valid field access new syntax",
			fieldAccess: FieldAccess{
				Object: "o",
				Field:  "total",
			},
			expressionIndex: 1,
			expectError:     false,
		},
		{
			name: "valid field in first pattern of multi-pattern",
			fieldAccess: FieldAccess{
				Object: "x",
				Field:  "age",
			},
			expressionIndex: 2,
			expectError:     false,
		},
		{
			name: "valid field in second pattern of multi-pattern",
			fieldAccess: FieldAccess{
				Object: "y",
				Field:  "id",
			},
			expressionIndex: 2,
			expectError:     false,
		},
		{
			name: "variable not found in empty expression",
			fieldAccess: FieldAccess{
				Object: "p",
				Field:  "name",
			},
			expressionIndex: 3,
			expectError:     true,
			errorContains:   "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFieldAccess(program, tt.fieldAccess, tt.expressionIndex)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorContains != "" && !stringContains(err.Error(), tt.errorContains) {
					t.Errorf("expected error to contain %q, got %q", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateTypeCompatibility_Coverage(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "name", Type: "string"},
					{Name: "age", Type: "number"},
				},
			},
		},
		Expressions: []Expression{
			{
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
						{Name: "x", DataType: "number"},
					},
				},
			},
		},
	}

	tests := []struct {
		name        string
		constraint  interface{}
		expectError bool
	}{
		{
			name:        "non-map constraint returns nil",
			constraint:  "string constraint",
			expectError: false,
		},
		{
			name: "map without type field returns nil",
			constraint: map[string]interface{}{
				"value": "test",
			},
			expectError: false,
		},
		{
			name: "comparison constraint with valid operands",
			constraint: map[string]interface{}{
				"type": ConstraintTypeComparison,
				"left": map[string]interface{}{
					"type":  "numberLiteral",
					"value": 42,
				},
				"right": map[string]interface{}{
					"type":  "numberLiteral",
					"value": 100,
				},
				"operator": "==",
			},
			expectError: false,
		},
		{
			name: "logical expression constraint",
			constraint: map[string]interface{}{
				"type":     ConstraintTypeLogicalExpr,
				"operator": "AND",
				"constraints": []interface{}{
					map[string]interface{}{
						"type": ConstraintTypeComparison,
						"left": map[string]interface{}{
							"type":  "numberLiteral",
							"value": 1,
						},
						"right": map[string]interface{}{
							"type":  "numberLiteral",
							"value": 2,
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "binary operation constraint",
			constraint: map[string]interface{}{
				"type": ConstraintTypeBinaryOp,
				"left": map[string]interface{}{
					"type":  "numberLiteral",
					"value": 5,
				},
				"right": map[string]interface{}{
					"type":  "numberLiteral",
					"value": 3,
				},
				"operator": "+",
			},
			expectError: false,
		},
		{
			name: "unknown constraint type returns nil",
			constraint: map[string]interface{}{
				"type": "unknownType",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTypeCompatibility(program, tt.constraint, 0)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateConstraintWithOperands_Coverage(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "name", Type: "string"},
					{Name: "age", Type: "number"},
				},
			},
		},
		Expressions: []Expression{
			{
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
					},
				},
			},
		},
	}

	tests := []struct {
		name               string
		constraint         map[string]interface{}
		checkCompatibility bool
		expectError        bool
	}{
		{
			name: "nil left operand",
			constraint: map[string]interface{}{
				"left":  nil,
				"right": map[string]interface{}{"type": "numberLiteral", "value": 1},
			},
			checkCompatibility: true,
			expectError:        false,
		},
		{
			name: "nil right operand",
			constraint: map[string]interface{}{
				"left":  map[string]interface{}{"type": "numberLiteral", "value": 1},
				"right": nil,
			},
			checkCompatibility: true,
			expectError:        false,
		},
		{
			name: "compatible types with checking",
			constraint: map[string]interface{}{
				"left":  map[string]interface{}{"type": "numberLiteral", "value": 1},
				"right": map[string]interface{}{"type": "numberLiteral", "value": 2},
			},
			checkCompatibility: true,
			expectError:        false,
		},
		{
			name: "no compatibility check",
			constraint: map[string]interface{}{
				"left":  map[string]interface{}{"type": "numberLiteral", "value": 1},
				"right": map[string]interface{}{"type": "stringLiteral", "value": "test"},
			},
			checkCompatibility: false,
			expectError:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConstraintWithOperands(program, tt.constraint, 0, tt.checkCompatibility, 0)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateOperandTypeCompatibility_Coverage(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "name", Type: "string"},
					{Name: "age", Type: "number"},
				},
			},
		},
		Expressions: []Expression{
			{
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
						{Name: "count", DataType: "number"},
					},
				},
			},
		},
	}

	tests := []struct {
		name        string
		left        interface{}
		right       interface{}
		expectError bool
	}{
		{
			name:        "variable vs number - compatible",
			left:        map[string]interface{}{"type": "variable", "name": "count"},
			right:       map[string]interface{}{"type": "numberLiteral", "value": 10},
			expectError: false,
		},
		{
			name:        "number vs variable - compatible",
			left:        map[string]interface{}{"type": "numberLiteral", "value": 10},
			right:       map[string]interface{}{"type": "variable", "name": "count"},
			expectError: false,
		},
		{
			name:        "same types - compatible",
			left:        map[string]interface{}{"type": "numberLiteral", "value": 10},
			right:       map[string]interface{}{"type": "numberLiteral", "value": 20},
			expectError: false,
		},
		{
			name:        "unknown type - compatible",
			left:        map[string]interface{}{"type": "unknown"},
			right:       map[string]interface{}{"type": "numberLiteral", "value": 10},
			expectError: false,
		},
		{
			name:        "string vs number with explicit types",
			left:        map[string]interface{}{"type": "string"},
			right:       map[string]interface{}{"type": "number"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOperandTypeCompatibility(program, tt.left, tt.right, 0)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGetOperandType_Coverage(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "name", Type: "string"},
					{Name: "age", Type: "number"},
				},
			},
		},
		Expressions: []Expression{
			{
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
					},
				},
			},
		},
	}

	tests := []struct {
		name         string
		operand      interface{}
		expectedType string
		expectError  bool
	}{
		{
			name:         "non-map operand",
			operand:      "string value",
			expectedType: ValueTypeUnknown,
			expectError:  false,
		},
		{
			name: "field access operand",
			operand: map[string]interface{}{
				"type":   ConstraintTypeFieldAccess,
				"object": "p",
				"field":  "name",
			},
			expectedType: ValueTypeString,
			expectError:  false,
		},
		{
			name: "number type operand",
			operand: map[string]interface{}{
				"type":  "number",
				"value": 42,
			},
			expectedType: ValueTypeNumber,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, err := getOperandType(program, tt.operand, 0)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if gotType != tt.expectedType {
					t.Errorf("expected type %q, got %q", tt.expectedType, gotType)
				}
			}
		})
	}
}

func TestValidateExpressionConstraints_Coverage(t *testing.T) {
	tests := []struct {
		name        string
		program     Program
		expectError bool
	}{
		{
			name: "no constraints - valid",
			program: Program{
				Expressions: []Expression{
					{
						Set: Set{
							Variables: []TypedVariable{
								{Name: "x", DataType: "number"},
							},
						},
						Constraints: nil,
					},
				},
			},
			expectError: false,
		},
		{
			name: "valid constraints",
			program: Program{
				Types: []TypeDefinition{
					{
						Name: "Person",
						Fields: []Field{
							{Name: "age", Type: "number"},
						},
					},
				},
				Expressions: []Expression{
					{
						Set: Set{
							Variables: []TypedVariable{
								{Name: "p", DataType: "Person"},
							},
						},
						Constraints: map[string]interface{}{
							"type": ConstraintTypeComparison,
							"left": map[string]interface{}{
								"type":   ConstraintTypeFieldAccess,
								"object": "p",
								"field":  "age",
							},
							"right": map[string]interface{}{
								"type":  "numberLiteral",
								"value": 18,
							},
							"operator": ">",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "invalid field access",
			program: Program{
				Expressions: []Expression{
					{
						Set: Set{
							Variables: []TypedVariable{
								{Name: "p", DataType: "Person"},
							},
						},
						Constraints: map[string]interface{}{
							"type": ConstraintTypeComparison,
							"left": map[string]interface{}{
								"type":   ConstraintTypeFieldAccess,
								"object": "unknown",
								"field":  "age",
							},
							"right": map[string]interface{}{
								"type":  "numberLiteral",
								"value": 18,
							},
						},
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateExpressionConstraints(tt.program)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateExpressionActions_Coverage(t *testing.T) {
	tests := []struct {
		name        string
		program     Program
		expectError bool
	}{
		{
			name: "expression with action",
			program: Program{
				Types: []TypeDefinition{
					{
						Name: "TestType",
						Fields: []Field{
							{Name: "id", Type: "string"},
						},
					},
				},
				Actions: []ActionDefinition{
					{
						Name:       "testAction",
						Parameters: []Parameter{},
					},
				},
				Expressions: []Expression{
					{
						Set: Set{
							Variables: []TypedVariable{
								{Name: "x", DataType: "TestType"},
							},
						},
						Action: &Action{
							Type: "action",
							Job: &JobCall{
								Type: "jobCall",
								Name: "testAction",
								Args: []interface{}{},
							},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "expression without action - error",
			program: Program{
				Expressions: []Expression{
					{
						Set: Set{
							Variables: []TypedVariable{
								{Name: "x", DataType: "string"},
							},
						},
						Action: nil,
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateExpressionActions(tt.program)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateLogicalExpressionConstraint_Coverage(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "age", Type: "number"},
					{Name: "name", Type: "string"},
				},
			},
		},
		Expressions: []Expression{
			{
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
					},
				},
			},
		},
	}

	tests := []struct {
		name        string
		constraint  map[string]interface{}
		expectError bool
	}{
		{
			name: "logical expression with left and operations",
			constraint: map[string]interface{}{
				"type": ConstraintTypeLogicalExpr,
				"left": map[string]interface{}{
					"type":  "number",
					"value": 10,
				},
				"operations": []interface{}{
					map[string]interface{}{
						"operator": "AND",
						"right": map[string]interface{}{
							"type":  "number",
							"value": 20,
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "logical expression with only left",
			constraint: map[string]interface{}{
				"type": ConstraintTypeLogicalExpr,
				"left": map[string]interface{}{
					"type": ConstraintTypeComparison,
					"left": map[string]interface{}{
						"type":   ConstraintTypeFieldAccess,
						"object": "p",
						"field":  "age",
					},
					"right": map[string]interface{}{
						"type":  "number",
						"value": 18,
					},
				},
			},
			expectError: false,
		},
		{
			name: "logical expression with multiple operations",
			constraint: map[string]interface{}{
				"type": ConstraintTypeLogicalExpr,
				"left": map[string]interface{}{
					"type":  "number",
					"value": 1,
				},
				"operations": []interface{}{
					map[string]interface{}{
						"operator": "AND",
						"right": map[string]interface{}{
							"type":  "number",
							"value": 2,
						},
					},
					map[string]interface{}{
						"operator": "OR",
						"right": map[string]interface{}{
							"type":  "number",
							"value": 3,
						},
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLogicalExpressionConstraint(program, tt.constraint, 0, 0)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateFieldAccessInOperands_Coverage(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "name", Type: "string"},
					{Name: "age", Type: "number"},
				},
			},
		},
		Expressions: []Expression{
			{
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
					},
				},
			},
		},
	}

	tests := []struct {
		name        string
		constraint  map[string]interface{}
		expectError bool
	}{
		{
			name: "valid field access in left operand",
			constraint: map[string]interface{}{
				"left": map[string]interface{}{
					"type":   ConstraintTypeFieldAccess,
					"object": "p",
					"field":  "name",
				},
				"right": map[string]interface{}{
					"type":  "string",
					"value": "John",
				},
			},
			expectError: false,
		},
		{
			name: "valid field access in right operand",
			constraint: map[string]interface{}{
				"left": map[string]interface{}{
					"type":  "number",
					"value": 25,
				},
				"right": map[string]interface{}{
					"type":   ConstraintTypeFieldAccess,
					"object": "p",
					"field":  "age",
				},
			},
			expectError: false,
		},
		{
			name: "invalid field access in left operand",
			constraint: map[string]interface{}{
				"left": map[string]interface{}{
					"type":   ConstraintTypeFieldAccess,
					"object": "unknown",
					"field":  "name",
				},
				"right": map[string]interface{}{
					"type":  "string",
					"value": "test",
				},
			},
			expectError: true,
		},
		{
			name: "both operands valid field access",
			constraint: map[string]interface{}{
				"left": map[string]interface{}{
					"type":   ConstraintTypeFieldAccess,
					"object": "p",
					"field":  "name",
				},
				"right": map[string]interface{}{
					"type":   ConstraintTypeFieldAccess,
					"object": "p",
					"field":  "name",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFieldAccessInOperands(program, tt.constraint, 0, 0)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateFieldAccessInLogicalExpr_Coverage(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "age", Type: "number"},
				},
			},
		},
		Expressions: []Expression{
			{
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
					},
				},
			},
		},
	}

	tests := []struct {
		name        string
		constraint  map[string]interface{}
		expectError bool
	}{
		{
			name: "valid field access in left operand",
			constraint: map[string]interface{}{
				"type": ConstraintTypeLogicalExpr,
				"left": map[string]interface{}{
					"type":   ConstraintTypeFieldAccess,
					"object": "p",
					"field":  "age",
				},
				"operations": []interface{}{
					map[string]interface{}{
						"operator": "AND",
						"right": map[string]interface{}{
							"type":   ConstraintTypeFieldAccess,
							"object": "p",
							"field":  "age",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "invalid field access in operations",
			constraint: map[string]interface{}{
				"type": ConstraintTypeLogicalExpr,
				"left": map[string]interface{}{
					"type":  "number",
					"value": 1,
				},
				"operations": []interface{}{
					map[string]interface{}{
						"operator": "OR",
						"right": map[string]interface{}{
							"type":   ConstraintTypeFieldAccess,
							"object": "unknown",
							"field":  "age",
						},
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFieldAccessInLogicalExpr(program, tt.constraint, 0, 0)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateAction_Coverage(t *testing.T) {
	tests := []struct {
		name        string
		program     Program
		action      Action
		exprIndex   int
		expectError bool
	}{
		{
			name: "invalid expression index",
			program: Program{
				Expressions: []Expression{},
			},
			action: Action{
				Type: "action",
			},
			exprIndex:   99,
			expectError: true,
		},
		{
			name: "multiple jobs with valid variables",
			program: Program{
				Types: []TypeDefinition{
					{
						Name: "Person",
						Fields: []Field{
							{Name: "id", Type: "string"},
							{Name: "name", Type: "string"},
						},
					},
				},
				Expressions: []Expression{
					{
						Set: Set{
							Variables: []TypedVariable{
								{Name: "p", DataType: "Person"},
								{Name: "count", DataType: "number"},
							},
						},
					},
				},
			},
			action: Action{
				Type: "action",
				Jobs: []JobCall{
					{
						Type: "jobCall",
						Name: "job1",
						Args: []interface{}{
							"p",
							"count",
						},
					},
					{
						Type: "jobCall",
						Name: "job2",
						Args: []interface{}{
							map[string]interface{}{
								"type":   "fieldAccess",
								"object": "p",
								"field":  "id",
							},
						},
					},
				},
			},
			exprIndex:   0,
			expectError: false,
		},
		{
			name: "job with invalid variable",
			program: Program{
				Expressions: []Expression{
					{
						Set: Set{
							Variables: []TypedVariable{
								{Name: "p", DataType: "Person"},
							},
						},
					},
				},
			},
			action: Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "testJob",
					Args: []interface{}{"unknownVar"},
				},
			},
			exprIndex:   0,
			expectError: true,
		},
		{
			name: "action with string literal args - valid",
			program: Program{
				Expressions: []Expression{
					{
						Set: Set{
							Variables: []TypedVariable{
								{Name: "x", DataType: "number"},
							},
						},
					},
				},
			},
			action: Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "testJob",
					Args: []interface{}{
						map[string]interface{}{
							"type":  "string",
							"value": "test",
						},
					},
				},
			},
			exprIndex:   0,
			expectError: false,
		},
		{
			name: "action with number literal args - valid",
			program: Program{
				Expressions: []Expression{
					{
						Set: Set{
							Variables: []TypedVariable{
								{Name: "x", DataType: "number"},
							},
						},
					},
				},
			},
			action: Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "testJob",
					Args: []interface{}{
						map[string]interface{}{
							"type":  "number",
							"value": 42,
						},
					},
				},
			},
			exprIndex:   0,
			expectError: false,
		},
		{
			name: "variables from patterns (multi-pattern syntax)",
			program: Program{
				Expressions: []Expression{
					{
						Patterns: []Set{
							{
								Variables: []TypedVariable{
									{Name: "p1", DataType: "Person"},
								},
							},
							{
								Variables: []TypedVariable{
									{Name: "p2", DataType: "Person"},
								},
							},
						},
					},
				},
			},
			action: Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "testJob",
					Args: []interface{}{"p1", "p2"},
				},
			},
			exprIndex:   0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAction(tt.program, tt.action, tt.exprIndex)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGetFieldType_Coverage(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "id", Type: "string"},
					{Name: "name", Type: "string"},
					{Name: "age", Type: "number"},
				},
			},
			{
				Name: "Order",
				Fields: []Field{
					{Name: "id", Type: "string"},
					{Name: "total", Type: "number"},
				},
			},
		},
		Expressions: []Expression{
			{
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
					},
				},
			},
			{
				Patterns: []Set{
					{
						Variables: []TypedVariable{
							{Name: "o", DataType: "Order"},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name         string
		object       string
		field        string
		exprIndex    int
		expectedType string
		expectError  bool
	}{
		{
			name:         "valid field in old syntax",
			object:       "p",
			field:        "name",
			exprIndex:    0,
			expectedType: "string",
			expectError:  false,
		},
		{
			name:         "valid field number type",
			object:       "p",
			field:        "age",
			exprIndex:    0,
			expectedType: "number",
			expectError:  false,
		},
		{
			name:         "valid field in new pattern syntax",
			object:       "o",
			field:        "total",
			exprIndex:    1,
			expectedType: "number",
			expectError:  false,
		},
		{
			name:        "invalid expression index",
			object:      "p",
			field:       "name",
			exprIndex:   99,
			expectError: true,
		},
		{
			name:        "variable not found",
			object:      "unknownVar",
			field:       "name",
			exprIndex:   0,
			expectError: true,
		},
		{
			name:        "field not found in type",
			object:      "p",
			field:       "unknownField",
			exprIndex:   0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldType, err := GetFieldType(program, tt.object, tt.field, tt.exprIndex)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if fieldType != tt.expectedType {
					t.Errorf("Expected type %q, got %q", tt.expectedType, fieldType)
				}
			}
		})
	}
}

func TestValidateProgram_Coverage(t *testing.T) {
	tests := []struct {
		name        string
		program     Program
		expectError bool
	}{
		{
			name: "valid complete program",
			program: Program{
				Types: []TypeDefinition{
					{
						Name: "Person",
						Fields: []Field{
							{Name: "id", Type: "string"},
							{Name: "age", Type: "number"},
						},
					},
				},
				Expressions: []Expression{
					{
						Set: Set{
							Variables: []TypedVariable{
								{Name: "p", DataType: "Person"},
							},
						},
						Constraints: map[string]interface{}{
							"type": ConstraintTypeComparison,
							"left": map[string]interface{}{
								"type":   ConstraintTypeFieldAccess,
								"object": "p",
								"field":  "age",
							},
							"right": map[string]interface{}{
								"type":  "number",
								"value": 18,
							},
							"operator": ">",
						},
						Action: &Action{
							Type: "action",
							Job: &JobCall{
								Type: "jobCall",
								Name: "approve",
								Args: []interface{}{},
							},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "program with constraint validation error",
			program: Program{
				Expressions: []Expression{
					{
						Set: Set{
							Variables: []TypedVariable{
								{Name: "p", DataType: "Person"},
							},
						},
						Constraints: map[string]interface{}{
							"type": ConstraintTypeComparison,
							"left": map[string]interface{}{
								"type":   ConstraintTypeFieldAccess,
								"object": "unknownVar",
								"field":  "age",
							},
							"right": map[string]interface{}{
								"type":  "number",
								"value": 18,
							},
						},
						Action: &Action{
							Type: "action",
							Job: &JobCall{
								Type: "jobCall",
								Name: "test",
								Args: []interface{}{},
							},
						},
					},
				},
			},
			expectError: true,
		},
		{
			name: "program with action validation error",
			program: Program{
				Types: []TypeDefinition{
					{
						Name: "Person",
						Fields: []Field{
							{Name: "id", Type: "string"},
						},
					},
				},
				Expressions: []Expression{
					{
						Set: Set{
							Variables: []TypedVariable{
								{Name: "p", DataType: "Person"},
							},
						},
						Action: &Action{
							Type: "action",
							Job: &JobCall{
								Type: "jobCall",
								Name: "test",
								Args: []interface{}{"invalidVar"},
							},
						},
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProgram(tt.program)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// Helper function for string contains check
func stringContains(s, substr string) bool {
	return len(s) >= len(substr) && stringContainsAt(s, substr)
}

func stringContainsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
