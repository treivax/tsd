// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAlphaConditionBuilder(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	assert.NotNil(t, builder)
}
func TestAlphaConditionBuilder_FieldEquals(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	tests := []struct {
		name     string
		variable string
		field    string
		value    interface{}
	}{
		{
			name:     "string value",
			variable: "x",
			field:    "name",
			value:    "test",
		},
		{
			name:     "int value",
			variable: "obj",
			field:    "age",
			value:    42,
		},
		{
			name:     "float value",
			variable: "item",
			field:    "price",
			value:    19.99,
		},
		{
			name:     "bool value",
			variable: "flag",
			field:    "enabled",
			value:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.FieldEquals(tt.variable, tt.field, tt.value)
			resultMap, ok := result.(map[string]interface{})
			require.True(t, ok, "result should be a map")
			assert.Equal(t, "binaryOperation", resultMap["type"])
			assert.Equal(t, "==", resultMap["operator"])
			left, ok := resultMap["left"].(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, "fieldAccess", left["type"])
			assert.Equal(t, tt.variable, left["object"])
			assert.Equal(t, tt.field, left["field"])
			right, ok := resultMap["right"].(map[string]interface{})
			require.True(t, ok)
			assert.NotNil(t, right["type"])
			assert.NotNil(t, right["value"])
		})
	}
}
func TestAlphaConditionBuilder_FieldNotEquals(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	result := builder.FieldNotEquals("x", "status", "inactive")
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "binaryOperation", resultMap["type"])
	assert.Equal(t, "!=", resultMap["operator"])
	left, ok := resultMap["left"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "fieldAccess", left["type"])
	assert.Equal(t, "x", left["object"])
	assert.Equal(t, "status", left["field"])
}
func TestAlphaConditionBuilder_FieldLessThan(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	result := builder.FieldLessThan("x", "age", 30)
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "binaryOperation", resultMap["type"])
	assert.Equal(t, "<", resultMap["operator"])
}
func TestAlphaConditionBuilder_FieldLessOrEqual(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	result := builder.FieldLessOrEqual("x", "score", 100)
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "binaryOperation", resultMap["type"])
	assert.Equal(t, "<=", resultMap["operator"])
}
func TestAlphaConditionBuilder_FieldGreaterThan(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	result := builder.FieldGreaterThan("x", "balance", 0)
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "binaryOperation", resultMap["type"])
	assert.Equal(t, ">", resultMap["operator"])
}
func TestAlphaConditionBuilder_FieldGreaterOrEqual(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	result := builder.FieldGreaterOrEqual("x", "rating", 4.5)
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "binaryOperation", resultMap["type"])
	assert.Equal(t, ">=", resultMap["operator"])
}
func TestAlphaConditionBuilder_And(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	left := builder.FieldEquals("x", "name", "test")
	right := builder.FieldEquals("x", "age", 30)
	result := builder.And(left, right)
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "logicalExpression", resultMap["type"])
	assert.NotNil(t, resultMap["left"])
	operations, ok := resultMap["operations"].([]interface{})
	require.True(t, ok)
	require.Len(t, operations, 1)
	op, ok := operations[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "AND", op["op"])
	assert.NotNil(t, op["right"])
}
func TestAlphaConditionBuilder_Or(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	left := builder.FieldEquals("x", "status", "active")
	right := builder.FieldEquals("x", "status", "pending")
	result := builder.Or(left, right)
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "logicalExpression", resultMap["type"])
	assert.NotNil(t, resultMap["left"])
	operations, ok := resultMap["operations"].([]interface{})
	require.True(t, ok)
	require.Len(t, operations, 1)
	op, ok := operations[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "OR", op["op"])
	assert.NotNil(t, op["right"])
}
func TestAlphaConditionBuilder_AndMultiple(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	tests := []struct {
		name           string
		conditions     []interface{}
		expectedOpsLen int
		shouldBeTrue   bool
		shouldBeSimple bool
	}{
		{
			name:         "no conditions",
			conditions:   []interface{}{},
			shouldBeTrue: true,
		},
		{
			name:           "single condition",
			conditions:     []interface{}{builder.FieldEquals("x", "name", "test")},
			shouldBeSimple: true,
		},
		{
			name: "two conditions",
			conditions: []interface{}{
				builder.FieldEquals("x", "name", "test"),
				builder.FieldEquals("x", "age", 30),
			},
			expectedOpsLen: 1,
		},
		{
			name: "three conditions",
			conditions: []interface{}{
				builder.FieldEquals("x", "name", "test"),
				builder.FieldEquals("x", "age", 30),
				builder.FieldGreaterThan("x", "score", 80),
			},
			expectedOpsLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.AndMultiple(tt.conditions...)
			if tt.shouldBeTrue {
				resultMap, ok := result.(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "booleanLiteral", resultMap["type"])
				assert.Equal(t, true, resultMap["value"])
				return
			}
			if tt.shouldBeSimple {
				// Should return the single condition as-is
				assert.Equal(t, tt.conditions[0], result)
				return
			}
			resultMap, ok := result.(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, "logicalExpression", resultMap["type"])
			assert.NotNil(t, resultMap["left"])
			operations, ok := resultMap["operations"].([]interface{})
			require.True(t, ok)
			assert.Len(t, operations, tt.expectedOpsLen)
			for _, op := range operations {
				opMap, ok := op.(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "AND", opMap["op"])
				assert.NotNil(t, opMap["right"])
			}
		})
	}
}
func TestAlphaConditionBuilder_OrMultiple(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	tests := []struct {
		name           string
		conditions     []interface{}
		expectedOpsLen int
		shouldBeFalse  bool
		shouldBeSimple bool
	}{
		{
			name:          "no conditions",
			conditions:    []interface{}{},
			shouldBeFalse: true,
		},
		{
			name:           "single condition",
			conditions:     []interface{}{builder.FieldEquals("x", "status", "active")},
			shouldBeSimple: true,
		},
		{
			name: "two conditions",
			conditions: []interface{}{
				builder.FieldEquals("x", "status", "active"),
				builder.FieldEquals("x", "status", "pending"),
			},
			expectedOpsLen: 1,
		},
		{
			name: "four conditions",
			conditions: []interface{}{
				builder.FieldEquals("x", "type", "A"),
				builder.FieldEquals("x", "type", "B"),
				builder.FieldEquals("x", "type", "C"),
				builder.FieldEquals("x", "type", "D"),
			},
			expectedOpsLen: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.OrMultiple(tt.conditions...)
			if tt.shouldBeFalse {
				resultMap, ok := result.(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "booleanLiteral", resultMap["type"])
				assert.Equal(t, false, resultMap["value"])
				return
			}
			if tt.shouldBeSimple {
				// Should return the single condition as-is
				assert.Equal(t, tt.conditions[0], result)
				return
			}
			resultMap, ok := result.(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, "logicalExpression", resultMap["type"])
			assert.NotNil(t, resultMap["left"])
			operations, ok := resultMap["operations"].([]interface{})
			require.True(t, ok)
			assert.Len(t, operations, tt.expectedOpsLen)
			for _, op := range operations {
				opMap, ok := op.(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "OR", opMap["op"])
				assert.NotNil(t, opMap["right"])
			}
		})
	}
}
func TestAlphaConditionBuilder_True(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	result := builder.True()
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "booleanLiteral", resultMap["type"])
	assert.Equal(t, true, resultMap["value"])
}
func TestAlphaConditionBuilder_False(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	result := builder.False()
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "booleanLiteral", resultMap["type"])
	assert.Equal(t, false, resultMap["value"])
}
func TestAlphaConditionBuilder_FieldRange(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	result := builder.FieldRange("x", "age", 18, 65)
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok)
	// Should be a logical expression combining >= and <=
	assert.Equal(t, "logicalExpression", resultMap["type"])
	assert.NotNil(t, resultMap["left"])
	operations, ok := resultMap["operations"].([]interface{})
	require.True(t, ok)
	require.Len(t, operations, 1)
	op, ok := operations[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "AND", op["op"])
}
func TestAlphaConditionBuilder_FieldIn(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	tests := []struct {
		name           string
		variable       string
		field          string
		values         []interface{}
		shouldBeFalse  bool
		expectedOpsLen int
	}{
		{
			name:          "no values",
			variable:      "x",
			field:         "status",
			values:        []interface{}{},
			shouldBeFalse: true,
		},
		{
			name:           "single value",
			variable:       "x",
			field:          "status",
			values:         []interface{}{"active"},
			expectedOpsLen: 0, // Single condition, no OR needed
		},
		{
			name:           "two values",
			variable:       "x",
			field:          "status",
			values:         []interface{}{"active", "pending"},
			expectedOpsLen: 1,
		},
		{
			name:           "three values",
			variable:       "x",
			field:          "type",
			values:         []interface{}{"A", "B", "C"},
			expectedOpsLen: 2,
		},
		{
			name:           "mixed types",
			variable:       "x",
			field:          "value",
			values:         []interface{}{1, "test", true},
			expectedOpsLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.FieldIn(tt.variable, tt.field, tt.values...)
			if tt.shouldBeFalse {
				resultMap, ok := result.(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "booleanLiteral", resultMap["type"])
				assert.Equal(t, false, resultMap["value"])
				return
			}
			// For single value, might be returned directly
			if tt.expectedOpsLen == 0 {
				resultMap, ok := result.(map[string]interface{})
				require.True(t, ok)
				// Should be a binary operation (equality)
				assert.Equal(t, "binaryOperation", resultMap["type"])
				return
			}
			resultMap, ok := result.(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, "logicalExpression", resultMap["type"])
			operations, ok := resultMap["operations"].([]interface{})
			require.True(t, ok)
			assert.Len(t, operations, tt.expectedOpsLen)
		})
	}
}
func TestAlphaConditionBuilder_FieldNotIn(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	tests := []struct {
		name           string
		variable       string
		field          string
		values         []interface{}
		shouldBeTrue   bool
		expectedOpsLen int
	}{
		{
			name:         "no values",
			variable:     "x",
			field:        "status",
			values:       []interface{}{},
			shouldBeTrue: true,
		},
		{
			name:           "single value",
			variable:       "x",
			field:          "status",
			values:         []interface{}{"inactive"},
			expectedOpsLen: 0, // Single condition
		},
		{
			name:           "two values",
			variable:       "x",
			field:          "status",
			values:         []interface{}{"inactive", "deleted"},
			expectedOpsLen: 1,
		},
		{
			name:           "multiple values",
			variable:       "x",
			field:          "type",
			values:         []interface{}{"X", "Y", "Z", "W"},
			expectedOpsLen: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.FieldNotIn(tt.variable, tt.field, tt.values...)
			if tt.shouldBeTrue {
				resultMap, ok := result.(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "booleanLiteral", resultMap["type"])
				assert.Equal(t, true, resultMap["value"])
				return
			}
			// For single value, might be returned directly
			if tt.expectedOpsLen == 0 {
				resultMap, ok := result.(map[string]interface{})
				require.True(t, ok)
				// Should be a binary operation (not equals)
				assert.Equal(t, "binaryOperation", resultMap["type"])
				return
			}
			resultMap, ok := result.(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, "logicalExpression", resultMap["type"])
			operations, ok := resultMap["operations"].([]interface{})
			require.True(t, ok)
			assert.Len(t, operations, tt.expectedOpsLen)
		})
	}
}
func TestAlphaConditionBuilder_createLiteral(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	tests := []struct {
		name         string
		value        interface{}
		expectedType string
	}{
		{
			name:         "string",
			value:        "hello",
			expectedType: "stringLiteral",
		},
		{
			name:         "int",
			value:        42,
			expectedType: "numberLiteral",
		},
		{
			name:         "int32",
			value:        int32(42),
			expectedType: "numberLiteral",
		},
		{
			name:         "int64",
			value:        int64(42),
			expectedType: "numberLiteral",
		},
		{
			name:         "float32",
			value:        float32(3.14),
			expectedType: "numberLiteral",
		},
		{
			name:         "float64",
			value:        3.14159,
			expectedType: "numberLiteral",
		},
		{
			name:         "bool true",
			value:        true,
			expectedType: "booleanLiteral",
		},
		{
			name:         "bool false",
			value:        false,
			expectedType: "booleanLiteral",
		},
		{
			name:         "nil fallback to string",
			value:        nil,
			expectedType: "stringLiteral",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.createLiteral(tt.value)
			assert.Equal(t, tt.expectedType, result["type"])
			assert.NotNil(t, result["value"])
		})
	}
}
func TestAlphaConditionBuilder_CreateConstraintFromAST(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	t.Run("map constraint", func(t *testing.T) {
		constraint := map[string]interface{}{
			"type":  "test",
			"value": 123,
		}
		result := builder.CreateConstraintFromAST(constraint)
		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "test", resultMap["type"])
		assert.Equal(t, 123, resultMap["value"])
	})
	t.Run("non-map constraint", func(t *testing.T) {
		constraint := "some constraint"
		result := builder.CreateConstraintFromAST(constraint)
		assert.Equal(t, constraint, result)
	})
}
func TestAlphaConditionBuilder_Integration(t *testing.T) {
	builder := NewAlphaConditionBuilder()
	t.Run("complex condition with multiple operators", func(t *testing.T) {
		// (age >= 18 AND age <= 65) AND (status IN ["active", "pending"])
		ageRange := builder.FieldRange("person", "age", 18, 65)
		statusIn := builder.FieldIn("person", "status", "active", "pending")
		condition := builder.And(ageRange, statusIn)
		resultMap, ok := condition.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "logicalExpression", resultMap["type"])
	})
	t.Run("OR with NOT IN", func(t *testing.T) {
		// status = "active" OR (type NOT IN ["X", "Y"])
		activeStatus := builder.FieldEquals("item", "status", "active")
		typeNotIn := builder.FieldNotIn("item", "type", "X", "Y")
		condition := builder.Or(activeStatus, typeNotIn)
		resultMap, ok := condition.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "logicalExpression", resultMap["type"])
	})
}
