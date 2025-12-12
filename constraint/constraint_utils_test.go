// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFieldTypeExtended(t *testing.T) {
	t.Run("get field type from single pattern", func(t *testing.T) {
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

		fieldType, err := GetFieldType(program, "p", "name", 0)
		require.NoError(t, err)
		assert.Equal(t, "string", fieldType)

		fieldType, err = GetFieldType(program, "p", "age", 0)
		require.NoError(t, err)
		assert.Equal(t, "number", fieldType)
	})

	t.Run("get field type from multi-pattern", func(t *testing.T) {
		program := Program{
			Types: []TypeDefinition{
				{
					Name: "Company",
					Fields: []Field{
						{Name: "name", Type: "string"},
						{Name: "revenue", Type: "number"},
					},
				},
			},
			Expressions: []Expression{
				{
					Patterns: []Set{
						{
							Variables: []TypedVariable{
								{Name: "c", DataType: "Company"},
							},
						},
					},
				},
			},
		}

		fieldType, err := GetFieldType(program, "c", "name", 0)
		require.NoError(t, err)
		assert.Equal(t, "string", fieldType)

		fieldType, err = GetFieldType(program, "c", "revenue", 0)
		require.NoError(t, err)
		assert.Equal(t, "number", fieldType)
	})

	t.Run("invalid expression index", func(t *testing.T) {
		program := Program{
			Expressions: []Expression{},
		}

		fieldType, err := GetFieldType(program, "p", "name", 0)
		assert.Error(t, err)
		assert.Equal(t, "", fieldType)
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("variable not found", func(t *testing.T) {
		program := Program{
			Types: []TypeDefinition{
				{
					Name: "Person",
					Fields: []Field{
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

		fieldType, err := GetFieldType(program, "unknown", "name", 0)
		assert.Error(t, err)
		assert.Equal(t, "", fieldType)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("field not found in type", func(t *testing.T) {
		program := Program{
			Types: []TypeDefinition{
				{
					Name: "Person",
					Fields: []Field{
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

		fieldType, err := GetFieldType(program, "p", "nonexistent", 0)
		assert.Error(t, err)
		assert.Equal(t, "", fieldType)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("type definition not found", func(t *testing.T) {
		program := Program{
			Types: []TypeDefinition{},
			Expressions: []Expression{
				{
					Set: Set{
						Variables: []TypedVariable{
							{Name: "p", DataType: "UnknownType"},
						},
					},
				},
			},
		}

		fieldType, err := GetFieldType(program, "p", "name", 0)
		assert.Error(t, err)
		assert.Equal(t, "", fieldType)
	})
}

func TestGetValueTypeExtended(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name: "string value",
			value: map[string]interface{}{
				"type":  "string",
				"value": "hello",
			},
			expected: "string",
		},
		{
			name: "number value",
			value: map[string]interface{}{
				"type":  "number",
				"value": 42.0,
			},
			expected: "number",
		},
		{
			name: "boolean value",
			value: map[string]interface{}{
				"type":  "boolean",
				"value": true,
			},
			expected: "bool",
		},
		{
			name: "field access returns unknown",
			value: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "name",
			},
			expected: "unknown",
		},
		{
			name: "function call returns unknown",
			value: map[string]interface{}{
				"type": "functionCall",
				"name": "concat",
			},
			expected: "unknown",
		},
		{
			name: "variable reference",
			value: map[string]interface{}{
				"type": "variable",
				"name": "x",
			},
			expected: "variable",
		},
		{
			name:     "nil value",
			value:    nil,
			expected: "unknown",
		},
		{
			name:     "non-map value",
			value:    "plain string",
			expected: "unknown",
		},
		{
			name: "map without type field",
			value: map[string]interface{}{
				"value": "test",
			},
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetValueType(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateActionExtended(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "name", Type: "string"},
				},
			},
		},
		Expressions: []Expression{
			{
				Type:   "expression",
				RuleId: "TestRule",
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
					},
				},
			},
		},
	}

	t.Run("valid action with variable argument", func(t *testing.T) {
		action := Action{
			Type: "action",
			Job: &JobCall{
				Type: "jobCall",
				Name: "log",
				Args: []interface{}{"p"},
			},
		}

		err := ValidateAction(program, action, 0)
		require.NoError(t, err)
	})

	t.Run("valid action with field access", func(t *testing.T) {
		action := Action{
			Type: "action",
			Job: &JobCall{
				Type: "jobCall",
				Name: "log",
				Args: []interface{}{
					map[string]interface{}{
						"type":   "fieldAccess",
						"object": "p",
						"field":  "name",
					},
				},
			},
		}

		err := ValidateAction(program, action, 0)
		require.NoError(t, err)
	})

	t.Run("invalid expression index", func(t *testing.T) {
		action := Action{
			Type: "action",
			Job: &JobCall{
				Type: "jobCall",
				Name: "log",
				Args: []interface{}{"p"},
			},
		}

		err := ValidateAction(program, action, 999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("action with undefined variable", func(t *testing.T) {
		action := Action{
			Type: "action",
			Job: &JobCall{
				Type: "jobCall",
				Name: "log",
				Args: []interface{}{"unknownVar"},
			},
		}

		err := ValidateAction(program, action, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ne correspond à aucune variable")
	})

	t.Run("action with multiple jobs", func(t *testing.T) {
		action := Action{
			Type: "action",
			Jobs: []JobCall{
				{
					Type: "jobCall",
					Name: "log",
					Args: []interface{}{"p"},
				},
				{
					Type: "jobCall",
					Name: "notify",
					Args: []interface{}{
						map[string]interface{}{
							"type":   "fieldAccess",
							"object": "p",
							"field":  "name",
						},
					},
				},
			},
		}

		err := ValidateAction(program, action, 0)
		require.NoError(t, err)
	})

	t.Run("action with literal args", func(t *testing.T) {
		action := Action{
			Type: "action",
			Job: &JobCall{
				Type: "jobCall",
				Name: "log",
				Args: []interface{}{
					map[string]interface{}{
						"type":  "string",
						"value": "Hello",
					},
				},
			},
		}

		err := ValidateAction(program, action, 0)
		// String literals should not cause errors
		require.NoError(t, err)
	})
}

func TestValidateFieldAccessExtended(t *testing.T) {
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "name", Type: "string"},
					{Name: "age", Type: "number"},
				},
			},
			{
				Name: "Company",
				Fields: []Field{
					{Name: "name", Type: "string"},
					{Name: "employees", Type: "number"},
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

	t.Run("valid field access", func(t *testing.T) {
		fa := FieldAccess{
			Type:   "fieldAccess",
			Object: "p",
			Field:  "name",
		}
		err := ValidateFieldAccess(program, fa, 0)
		require.NoError(t, err)
	})

	t.Run("valid field access on age", func(t *testing.T) {
		fa := FieldAccess{
			Type:   "fieldAccess",
			Object: "p",
			Field:  "age",
		}
		err := ValidateFieldAccess(program, fa, 0)
		require.NoError(t, err)
	})

	t.Run("undefined variable", func(t *testing.T) {
		fa := FieldAccess{
			Type:   "fieldAccess",
			Object: "unknown",
			Field:  "name",
		}
		err := ValidateFieldAccess(program, fa, 0)
		assert.Error(t, err)
	})

	t.Run("undefined field", func(t *testing.T) {
		fa := FieldAccess{
			Type:   "fieldAccess",
			Object: "p",
			Field:  "nonexistent",
		}
		err := ValidateFieldAccess(program, fa, 0)
		assert.Error(t, err)
	})

	t.Run("invalid expression index", func(t *testing.T) {
		fa := FieldAccess{
			Type:   "fieldAccess",
			Object: "p",
			Field:  "name",
		}
		err := ValidateFieldAccess(program, fa, 999)
		assert.Error(t, err)
	})
}

func TestGetTypeFieldsExtended(t *testing.T) {
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
				Name:   "EmptyType",
				Fields: []Field{},
			},
		},
	}

	t.Run("get fields from existing type", func(t *testing.T) {
		fields, err := GetTypeFields(program, "Person")
		require.NoError(t, err)
		assert.Len(t, fields, 3)
		assert.Equal(t, "name", fields[0].Name)
		assert.Equal(t, "age", fields[1].Name)
		assert.Equal(t, "active", fields[2].Name)
	})

	t.Run("get fields from empty type", func(t *testing.T) {
		fields, err := GetTypeFields(program, "EmptyType")
		require.NoError(t, err)
		assert.Len(t, fields, 0)
	})

	t.Run("get fields from non-existent type", func(t *testing.T) {
		fields, err := GetTypeFields(program, "NonExistent")
		assert.Error(t, err)
		assert.Nil(t, fields)
	})

	t.Run("empty type name", func(t *testing.T) {
		fields, err := GetTypeFields(program, "")
		assert.Error(t, err)
		assert.Nil(t, fields)
	})
}

func TestValidateProgram_Integration(t *testing.T) {
	t.Run("validate complete program", func(t *testing.T) {
		program := map[string]interface{}{
			"types": []interface{}{
				map[string]interface{}{
					"name": "Person",
					"fields": []interface{}{
						map[string]interface{}{"name": "name", "type": "string"},
						map[string]interface{}{"name": "age", "type": "number"},
					},
				},
			},
			"actions": []interface{}{
				map[string]interface{}{
					"name": "greet",
					"parameters": []interface{}{
						map[string]interface{}{"name": "name", "type": "string"},
					},
				},
			},
			"expressions": []interface{}{},
		}

		err := ValidateProgram(program)
		// Should not panic
		_ = err
	})

	t.Run("validate program with invalid structure", func(t *testing.T) {
		program := "not a valid program"

		err := ValidateProgram(program)
		assert.Error(t, err)
	})
}
