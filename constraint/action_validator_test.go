// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionValidator_inferFunctionReturnType(t *testing.T) {
	av := NewActionValidator(nil, nil)

	tests := []struct {
		name         string
		funcName     string
		expectedType string
	}{
		// String functions
		{
			name:         "LENGTH returns number",
			funcName:     "LENGTH",
			expectedType: "number",
		},
		{
			name:         "length lowercase returns number",
			funcName:     "length",
			expectedType: "number",
		},
		{
			name:         "SUBSTRING returns string",
			funcName:     "SUBSTRING",
			expectedType: "string",
		},
		{
			name:         "substring lowercase returns string",
			funcName:     "substring",
			expectedType: "string",
		},
		{
			name:         "UPPER returns string",
			funcName:     "UPPER",
			expectedType: "string",
		},
		{
			name:         "LOWER returns string",
			funcName:     "LOWER",
			expectedType: "string",
		},
		{
			name:         "TRIM returns string",
			funcName:     "TRIM",
			expectedType: "string",
		},
		// Math functions
		{
			name:         "ABS returns number",
			funcName:     "ABS",
			expectedType: "number",
		},
		{
			name:         "abs lowercase returns number",
			funcName:     "abs",
			expectedType: "number",
		},
		{
			name:         "ROUND returns number",
			funcName:     "ROUND",
			expectedType: "number",
		},
		{
			name:         "FLOOR returns number",
			funcName:     "FLOOR",
			expectedType: "number",
		},
		{
			name:         "CEIL returns number",
			funcName:     "CEIL",
			expectedType: "number",
		},
		// Unknown/default
		{
			name:         "unknown function defaults to string",
			funcName:     "UNKNOWN_FUNC",
			expectedType: "string",
		},
		{
			name:         "empty function name defaults to string",
			funcName:     "",
			expectedType: "string",
		},
		{
			name:         "custom function defaults to string",
			funcName:     "MY_CUSTOM_FUNCTION",
			expectedType: "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := av.inferFunctionReturnType(tt.funcName)
			assert.Equal(t, tt.expectedType, result)
		})
	}
}

func TestActionValidator_GetActionDefinition(t *testing.T) {
	t.Run("get existing action", func(t *testing.T) {
		actions := []ActionDefinition{
			{
				Name: "print",
				Parameters: []Parameter{
					{Name: "message", Type: "string"},
				},
			},
			{
				Name: "log",
				Parameters: []Parameter{
					{Name: "level", Type: "string"},
					{Name: "msg", Type: "string"},
				},
			},
		}

		av := NewActionValidator(actions, nil)

		// Get first action
		action, exists := av.GetActionDefinition("print")
		require.True(t, exists, "action 'print' should exist")
		assert.NotNil(t, action)
		assert.Equal(t, "print", action.Name)
		assert.Len(t, action.Parameters, 1)
		assert.Equal(t, "message", action.Parameters[0].Name)

		// Get second action
		action, exists = av.GetActionDefinition("log")
		require.True(t, exists, "action 'log' should exist")
		assert.NotNil(t, action)
		assert.Equal(t, "log", action.Name)
		assert.Len(t, action.Parameters, 2)
	})

	t.Run("get non-existent action", func(t *testing.T) {
		actions := []ActionDefinition{
			{Name: "print", Parameters: []Parameter{}},
		}

		av := NewActionValidator(actions, nil)

		action, exists := av.GetActionDefinition("nonexistent")
		assert.False(t, exists, "nonexistent action should not exist")
		assert.Nil(t, action)
	})

	t.Run("get action from empty validator", func(t *testing.T) {
		av := NewActionValidator(nil, nil)

		action, exists := av.GetActionDefinition("print")
		assert.False(t, exists)
		assert.Nil(t, action)
	})

	t.Run("get action with empty name", func(t *testing.T) {
		actions := []ActionDefinition{
			{Name: "print", Parameters: []Parameter{}},
		}

		av := NewActionValidator(actions, nil)

		action, exists := av.GetActionDefinition("")
		assert.False(t, exists)
		assert.Nil(t, action)
	})

	t.Run("action names are case-sensitive", func(t *testing.T) {
		actions := []ActionDefinition{
			{Name: "Print", Parameters: []Parameter{}},
		}

		av := NewActionValidator(actions, nil)

		// Exact match should work
		action, exists := av.GetActionDefinition("Print")
		assert.True(t, exists)
		assert.NotNil(t, action)

		// Different case should not work
		action, exists = av.GetActionDefinition("print")
		assert.False(t, exists)
		assert.Nil(t, action)
	})
}

func TestActionValidator_GetTypeDefinition(t *testing.T) {
	t.Run("get existing type", func(t *testing.T) {
		types := []TypeDefinition{
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
				},
			},
		}

		av := NewActionValidator(nil, types)

		// Get first type
		typeDef, exists := av.GetTypeDefinition("Person")
		require.True(t, exists, "type 'Person' should exist")
		assert.NotNil(t, typeDef)
		assert.Equal(t, "Person", typeDef.Name)
		assert.Len(t, typeDef.Fields, 2)
		assert.Equal(t, "name", typeDef.Fields[0].Name)
		assert.Equal(t, "string", typeDef.Fields[0].Type)

		// Get second type
		typeDef, exists = av.GetTypeDefinition("Company")
		require.True(t, exists, "type 'Company' should exist")
		assert.NotNil(t, typeDef)
		assert.Equal(t, "Company", typeDef.Name)
		assert.Len(t, typeDef.Fields, 1)
	})

	t.Run("get non-existent type", func(t *testing.T) {
		types := []TypeDefinition{
			{Name: "Person", Fields: []Field{}},
		}

		av := NewActionValidator(nil, types)

		typeDef, exists := av.GetTypeDefinition("NonExistent")
		assert.False(t, exists, "nonexistent type should not exist")
		assert.Nil(t, typeDef)
	})

	t.Run("get type from empty validator", func(t *testing.T) {
		av := NewActionValidator(nil, nil)

		typeDef, exists := av.GetTypeDefinition("Person")
		assert.False(t, exists)
		assert.Nil(t, typeDef)
	})

	t.Run("get type with empty name", func(t *testing.T) {
		types := []TypeDefinition{
			{Name: "Person", Fields: []Field{}},
		}

		av := NewActionValidator(nil, types)

		typeDef, exists := av.GetTypeDefinition("")
		assert.False(t, exists)
		assert.Nil(t, typeDef)
	})

	t.Run("type names are case-sensitive", func(t *testing.T) {
		types := []TypeDefinition{
			{Name: "Person", Fields: []Field{}},
		}

		av := NewActionValidator(nil, types)

		// Exact match should work
		typeDef, exists := av.GetTypeDefinition("Person")
		assert.True(t, exists)
		assert.NotNil(t, typeDef)

		// Different case should not work
		typeDef, exists = av.GetTypeDefinition("person")
		assert.False(t, exists)
		assert.Nil(t, typeDef)
	})

	t.Run("get multiple types", func(t *testing.T) {
		types := []TypeDefinition{
			{Name: "Person", Fields: []Field{{Name: "name", Type: "string"}}},
			{Name: "Company", Fields: []Field{{Name: "name", Type: "string"}}},
			{Name: "Product", Fields: []Field{{Name: "price", Type: "number"}}},
		}

		av := NewActionValidator(nil, types)

		// All should be accessible
		for _, expectedType := range types {
			typeDef, exists := av.GetTypeDefinition(expectedType.Name)
			assert.True(t, exists, "type '%s' should exist", expectedType.Name)
			assert.Equal(t, expectedType.Name, typeDef.Name)
		}
	})
}

func TestActionValidator_Integration(t *testing.T) {
	t.Run("validator with both types and actions", func(t *testing.T) {
		types := []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "name", Type: "string"},
					{Name: "age", Type: "number"},
				},
			},
		}

		actions := []ActionDefinition{
			{
				Name: "notify",
				Parameters: []Parameter{
					{Name: "person", Type: "Person"},
					{Name: "message", Type: "string"},
				},
			},
		}

		av := NewActionValidator(actions, types)

		// Should be able to get both
		typeDef, typeExists := av.GetTypeDefinition("Person")
		assert.True(t, typeExists)
		assert.NotNil(t, typeDef)

		action, actionExists := av.GetActionDefinition("notify")
		assert.True(t, actionExists)
		assert.NotNil(t, action)

		// Action should reference the type
		assert.Equal(t, "Person", action.Parameters[0].Type)
	})

	t.Run("function return type inference", func(t *testing.T) {
		av := NewActionValidator(nil, nil)

		// Test a sequence of function calls
		functions := map[string]string{
			"LENGTH":    "number",
			"UPPER":     "string",
			"ABS":       "number",
			"SUBSTRING": "string",
			"ROUND":     "number",
			"UNKNOWN":   "string", // defaults
		}

		for funcName, expectedType := range functions {
			result := av.inferFunctionReturnType(funcName)
			assert.Equal(t, expectedType, result,
				"Function %s should return type %s", funcName, expectedType)
		}
	})
}
