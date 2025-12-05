package constraint

import (
	"strings"
	"testing"
)

// TestValidateConstraintWithOperands_NilOperands tests the nil operand early return path
func TestValidateConstraintWithOperands_NilOperands(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{Name: "Person", Fields: []Field{{Name: "age", Type: "number"}}},
		},
	}

	tests := []struct {
		name       string
		constraint map[string]interface{}
	}{
		{
			name: "both nil",
			constraint: map[string]interface{}{
				"type":  ConstraintTypeComparison,
				"left":  nil,
				"right": nil,
			},
		},
		{
			name: "left nil",
			constraint: map[string]interface{}{
				"type":  ConstraintTypeComparison,
				"left":  nil,
				"right": "test",
			},
		},
		{
			name: "right nil",
			constraint: map[string]interface{}{
				"type":  ConstraintTypeComparison,
				"left":  "test",
				"right": nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConstraintWithOperands(program, tt.constraint, 0, true)
			if err != nil {
				t.Errorf("validateConstraintWithOperands() with nil operands should return nil, got: %v", err)
			}
		})
	}
}

// TestValidateConstraintWithOperands_NestedValidation tests recursive validation of operands
func TestValidateConstraintWithOperands_NestedValidation(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{Name: "Person", Fields: []Field{
				{Name: "age", Type: "number"},
				{Name: "name", Type: "string"},
			}},
		},
		Expressions: []Expression{
			{
				Patterns: []Set{
					{
						Variables: []TypedVariable{
							{Name: "p", DataType: "Person"},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name          string
		constraint    map[string]interface{}
		checkCompat   bool
		expectError   bool
		errorContains string
	}{
		{
			name: "valid field access comparison",
			constraint: map[string]interface{}{
				"type": ConstraintTypeComparison,
				"left": map[string]interface{}{
					"type":   ConstraintTypeFieldAccess,
					"object": "p",
					"field":  "age",
				},
				"right": float64(30),
			},
			checkCompat: true,
			expectError: false,
		},
		{
			name: "invalid left field access",
			constraint: map[string]interface{}{
				"type": ConstraintTypeComparison,
				"left": map[string]interface{}{
					"type":   ConstraintTypeFieldAccess,
					"object": "p",
					"field":  "nonexistent",
				},
				"right": float64(30),
			},
			checkCompat:   true,
			expectError:   true,
			errorContains: "non trouvé", // French: "not found"
		},
		{
			name: "invalid right field access",
			constraint: map[string]interface{}{
				"type": ConstraintTypeComparison,
				"left": map[string]interface{}{
					"type":   ConstraintTypeFieldAccess,
					"object": "p",
					"field":  "age",
				},
				"right": map[string]interface{}{
					"type":   ConstraintTypeFieldAccess,
					"object": "p",
					"field":  "badfield",
				},
			},
			checkCompat:   true,
			expectError:   true,
			errorContains: "non trouvé",
		},
		{
			name: "without compatibility check",
			constraint: map[string]interface{}{
				"type": ConstraintTypeComparison,
				"left": map[string]interface{}{
					"type":   ConstraintTypeFieldAccess,
					"object": "p",
					"field":  "age",
				},
				"right": "string_value",
			},
			checkCompat: false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConstraintWithOperands(program, tt.constraint, 0, tt.checkCompat)
			if tt.expectError {
				if err == nil {
					t.Errorf("validateConstraintWithOperands() expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("validateConstraintWithOperands() error = %v, should contain %q", err, tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("validateConstraintWithOperands() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateOperandTypeCompatibility_AllCombinations tests various operand type combinations
func TestValidateOperandTypeCompatibility_AllCombinations(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{Name: "Person", Fields: []Field{
				{Name: "age", Type: "number"},
				{Name: "name", Type: "string"},
				{Name: "active", Type: "bool"},
			}},
		},
		Expressions: []Expression{
			{
				Patterns: []Set{
					{
						Variables: []TypedVariable{
							{Name: "p", DataType: "Person"},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name          string
		left          interface{}
		right         interface{}
		expectError   bool
		errorContains string
	}{
		{
			name: "number type map vs number type map - compatible",
			left: map[string]interface{}{
				"type":  "number",
				"value": float64(10),
			},
			right: map[string]interface{}{
				"type":  "number",
				"value": float64(20),
			},
			expectError: false,
		},
		{
			name: "string type map vs string type map - compatible",
			left: map[string]interface{}{
				"type":  "string",
				"value": "hello",
			},
			right: map[string]interface{}{
				"type":  "string",
				"value": "world",
			},
			expectError: false,
		},
		{
			name: "bool type map vs bool type map - compatible",
			left: map[string]interface{}{
				"type":  "bool",
				"value": true,
			},
			right: map[string]interface{}{
				"type":  "bool",
				"value": false,
			},
			expectError: false,
		},
		{
			name: "number vs string type maps - incompatible",
			left: map[string]interface{}{
				"type":  "number",
				"value": float64(10),
			},
			right: map[string]interface{}{
				"type":  "string",
				"value": "hello",
			},
			expectError:   true,
			errorContains: "incompatibilité de types",
		},
		{
			name: "variable vs number - allowed (aggregation)",
			left: map[string]interface{}{
				"type":  "variable",
				"value": "count",
			},
			right: map[string]interface{}{
				"type":  "number",
				"value": float64(5),
			},
			expectError: false,
		},
		{
			name: "number vs variable - allowed (aggregation)",
			left: map[string]interface{}{
				"type":  "number",
				"value": float64(5),
			},
			right: map[string]interface{}{
				"type":  "variable",
				"value": "count",
			},
			expectError: false,
		},
		{
			name: "field access vs compatible literal",
			left: map[string]interface{}{
				"type":   ConstraintTypeFieldAccess,
				"object": "p",
				"field":  "age",
			},
			right: map[string]interface{}{
				"type":  "number",
				"value": float64(30),
			},
			expectError: false,
		},
		{
			name: "field access vs incompatible literal",
			left: map[string]interface{}{
				"type":   ConstraintTypeFieldAccess,
				"object": "p",
				"field":  "age",
			},
			right: map[string]interface{}{
				"type":  "string",
				"value": "string_value",
			},
			expectError:   true,
			errorContains: "incompatibilité de types",
		},
		{
			name: "two field accesses - same type",
			left: map[string]interface{}{
				"type":   ConstraintTypeFieldAccess,
				"object": "p",
				"field":  "name",
			},
			right: map[string]interface{}{
				"type":   ConstraintTypeFieldAccess,
				"object": "p",
				"field":  "name",
			},
			expectError: false,
		},
		{
			name: "two field accesses - different types",
			left: map[string]interface{}{
				"type":   ConstraintTypeFieldAccess,
				"object": "p",
				"field":  "age",
			},
			right: map[string]interface{}{
				"type":   ConstraintTypeFieldAccess,
				"object": "p",
				"field":  "name",
			},
			expectError:   true,
			errorContains: "incompatibilité de types",
		},
		{
			name: "field access error propagation",
			left: map[string]interface{}{
				"type":   ConstraintTypeFieldAccess,
				"object": "p",
				"field":  "nonexistent",
			},
			right: map[string]interface{}{
				"type":  "number",
				"value": float64(10),
			},
			expectError:   true,
			errorContains: "non trouvé",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOperandTypeCompatibility(program, tt.left, tt.right, 0)
			if tt.expectError {
				if err == nil {
					t.Errorf("validateOperandTypeCompatibility() expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("validateOperandTypeCompatibility() error = %v, should contain %q", err, tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("validateOperandTypeCompatibility() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateLogicalExpressionConstraint_EdgeCases tests nested logical expressions
func TestValidateLogicalExpressionConstraint_EdgeCases(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{Name: "Person", Fields: []Field{
				{Name: "age", Type: "number"},
				{Name: "name", Type: "string"},
			}},
		},
		Expressions: []Expression{
			{
				Patterns: []Set{
					{
						Variables: []TypedVariable{
							{Name: "p", DataType: "Person"},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name          string
		constraint    map[string]interface{}
		expectError   bool
		errorContains string
	}{
		{
			name: "no left operand",
			constraint: map[string]interface{}{
				"type": ConstraintTypeLogicalExpr,
				"left": nil,
				"operations": []interface{}{
					map[string]interface{}{
						"operator": "&&",
						"right": map[string]interface{}{
							"type":  "number",
							"value": float64(10),
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "valid left with operations",
			constraint: map[string]interface{}{
				"type": ConstraintTypeLogicalExpr,
				"left": map[string]interface{}{
					"type":   ConstraintTypeFieldAccess,
					"object": "p",
					"field":  "age",
				},
				"operations": []interface{}{
					map[string]interface{}{
						"operator": "&&",
						"right": map[string]interface{}{
							"type":   ConstraintTypeFieldAccess,
							"object": "p",
							"field":  "name",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "left is field access - validated recursively",
			constraint: map[string]interface{}{
				"type": ConstraintTypeLogicalExpr,
				"left": map[string]interface{}{
					"type":   ConstraintTypeFieldAccess,
					"object": "p",
					"field":  "age",
				},
				"operations": []interface{}{},
			},
			expectError: false,
		},
		{
			name: "right is field access in operation - validated recursively",
			constraint: map[string]interface{}{
				"type": ConstraintTypeLogicalExpr,
				"left": map[string]interface{}{
					"type":  "number",
					"value": float64(10),
				},
				"operations": []interface{}{
					map[string]interface{}{
						"operator": "||",
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
			name: "multiple operations - all valid",
			constraint: map[string]interface{}{
				"type": ConstraintTypeLogicalExpr,
				"left": map[string]interface{}{
					"type":  "number",
					"value": float64(10),
				},
				"operations": []interface{}{
					map[string]interface{}{
						"operator": "&&",
						"right": map[string]interface{}{
							"type":  "number",
							"value": float64(20),
						},
					},
					map[string]interface{}{
						"operator": "||",
						"right": map[string]interface{}{
							"type":  "number",
							"value": float64(30),
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "multiple operations - mixed types",
			constraint: map[string]interface{}{
				"type": ConstraintTypeLogicalExpr,
				"left": map[string]interface{}{
					"type":  "number",
					"value": float64(10),
				},
				"operations": []interface{}{
					map[string]interface{}{
						"operator": "&&",
						"right": map[string]interface{}{
							"type":  "number",
							"value": float64(20),
						},
					},
					map[string]interface{}{
						"operator": "||",
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
			name: "operations not a slice",
			constraint: map[string]interface{}{
				"type": ConstraintTypeLogicalExpr,
				"left": map[string]interface{}{
					"type":  "number",
					"value": float64(10),
				},
				"operations": "not_a_slice",
			},
			expectError: false,
		},
		{
			name: "operation not a map",
			constraint: map[string]interface{}{
				"type": ConstraintTypeLogicalExpr,
				"left": map[string]interface{}{
					"type":  "number",
					"value": float64(10),
				},
				"operations": []interface{}{
					"not_a_map",
				},
			},
			expectError: false,
		},
		{
			name: "operation with nil right",
			constraint: map[string]interface{}{
				"type": ConstraintTypeLogicalExpr,
				"left": map[string]interface{}{
					"type":  "number",
					"value": float64(10),
				},
				"operations": []interface{}{
					map[string]interface{}{
						"operator": "&&",
						"right":    nil,
					},
				},
			},
			expectError: false,
		},
		{
			name: "nested logical expression in left",
			constraint: map[string]interface{}{
				"type": ConstraintTypeLogicalExpr,
				"left": map[string]interface{}{
					"type": ConstraintTypeLogicalExpr,
					"left": map[string]interface{}{
						"type":   ConstraintTypeFieldAccess,
						"object": "p",
						"field":  "age",
					},
					"operations": []interface{}{},
				},
				"operations": []interface{}{
					map[string]interface{}{
						"operator": "&&",
						"right": map[string]interface{}{
							"type":  "number",
							"value": float64(10),
						},
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLogicalExpressionConstraint(program, tt.constraint, 0)
			if tt.expectError {
				if err == nil {
					t.Errorf("validateLogicalExpressionConstraint() expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("validateLogicalExpressionConstraint() error = %v, should contain %q", err, tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("validateLogicalExpressionConstraint() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateProgram_ErrorPaths tests various error paths in ValidateProgram
func TestValidateProgram_ErrorPaths(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		expectError   bool
		errorContains string
	}{
		{
			name:          "invalid input type - string",
			input:         "not_a_map",
			expectError:   true,
			errorContains: "erreur", // Will get JSON unmarshal error
		},
		{
			name: "fact validation error - undefined type",
			input: map[string]interface{}{
				"types": []interface{}{
					map[string]interface{}{
						"name": "Person",
						"fields": []interface{}{
							map[string]interface{}{"name": "age", "type": "number"},
						},
					},
				},
				"expressions": []interface{}{},
				"facts": []interface{}{
					map[string]interface{}{
						"typeName": "UndefinedType",
						"fields": []interface{}{
							map[string]interface{}{
								"name": "id",
								"value": map[string]interface{}{
									"type":  "identifier",
									"value": "F001",
								},
							},
						},
					},
				},
			},
			expectError:   true,
			errorContains: "erreur validation faits",
		},
		{
			name: "valid program with all components",
			input: map[string]interface{}{
				"types": []interface{}{
					map[string]interface{}{
						"name": "Person",
						"fields": []interface{}{
							map[string]interface{}{"name": "age", "type": "number"},
							map[string]interface{}{"name": "name", "type": "string"},
						},
					},
				},
				"expressions": []interface{}{
					map[string]interface{}{
						"patterns": []interface{}{
							map[string]interface{}{
								"variables": []interface{}{
									map[string]interface{}{
										"name":     "p",
										"dataType": "Person",
									},
								},
							},
						},
						"constraints": []interface{}{
							map[string]interface{}{
								"type":     ConstraintTypeComparison,
								"operator": ">",
								"left": map[string]interface{}{
									"type":   ConstraintTypeFieldAccess,
									"object": "p",
									"field":  "age",
								},
								"right": map[string]interface{}{
									"type":  "number",
									"value": float64(18),
								},
							},
						},
						"action": map[string]interface{}{
							"type": "log",
							"arguments": []interface{}{
								"Adult person",
							},
						},
					},
				},
				"facts": []interface{}{
					map[string]interface{}{
						"typeName": "Person",
						"fields": []interface{}{
							map[string]interface{}{
								"name": "age",
								"value": map[string]interface{}{
									"type":  "number",
									"value": float64(25),
								},
							},
							map[string]interface{}{
								"name": "name",
								"value": map[string]interface{}{
									"type":  "string",
									"value": "John",
								},
							},
						},
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProgram(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateProgram() expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("ValidateProgram() error = %v, should contain %q", err, tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateProgram() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateBinaryOpConstraint_NoCompatCheck tests the backward compatibility wrapper
func TestValidateBinaryOpConstraint_NoCompatCheck(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{Name: "Person", Fields: []Field{
				{Name: "age", Type: "number"},
			}},
		},
		Expressions: []Expression{
			{
				Patterns: []Set{
					{
						Variables: []TypedVariable{
							{Name: "p", DataType: "Person"},
						},
					},
				},
			},
		},
	}

	// Binary op should not check compatibility (checkCompatibility=false)
	constraint := map[string]interface{}{
		"type": ConstraintTypeBinaryOp,
		"left": map[string]interface{}{
			"type":   ConstraintTypeFieldAccess,
			"object": "p",
			"field":  "age",
		},
		"right": "incompatible_string", // This would fail with compatibility check
	}

	err := validateBinaryOpConstraint(program, constraint, 0)
	if err != nil {
		t.Errorf("validateBinaryOpConstraint() should not check type compatibility, got error: %v", err)
	}
}

// TestValidateConstraintWithOperands_DeepNesting tests deeply nested constraint structures
func TestValidateConstraintWithOperands_DeepNesting(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{Name: "Person", Fields: []Field{
				{Name: "age", Type: "number"},
				{Name: "score", Type: "number"},
			}},
		},
		Expressions: []Expression{
			{
				Patterns: []Set{
					{
						Variables: []TypedVariable{
							{Name: "p", DataType: "Person"},
						},
					},
				},
			},
		},
	}

	// Nested binary operations in left operand
	constraint := map[string]interface{}{
		"type": ConstraintTypeComparison,
		"left": map[string]interface{}{
			"type":     ConstraintTypeBinaryOp,
			"operator": "+",
			"left": map[string]interface{}{
				"type":   ConstraintTypeFieldAccess,
				"object": "p",
				"field":  "age",
			},
			"right": map[string]interface{}{
				"type":   ConstraintTypeFieldAccess,
				"object": "p",
				"field":  "score",
			},
		},
		"right": map[string]interface{}{
			"type":  "number",
			"value": float64(100),
		},
	}

	err := validateConstraintWithOperands(program, constraint, 0, true)
	if err != nil {
		t.Errorf("validateConstraintWithOperands() with nested binary op failed: %v", err)
	}
}

// TestValidateOperandTypeCompatibility_UnknownTypes tests handling of unknown types
func TestValidateOperandTypeCompatibility_UnknownTypes(t *testing.T) {
	program := Program{
		Types:       []TypeDefinition{},
		Expressions: []Expression{},
	}

	// Plain Go primitives without type maps return "unknown" from GetValueType
	// Unknown types should be compatible with anything
	err := validateOperandTypeCompatibility(program, "string", float64(10), 0)
	if err != nil {
		t.Errorf("validateOperandTypeCompatibility() with unknown types should be compatible, got error: %v", err)
	}

	err = validateOperandTypeCompatibility(program, true, "test", 0)
	if err != nil {
		t.Errorf("validateOperandTypeCompatibility() with unknown types should be compatible, got error: %v", err)
	}
}
