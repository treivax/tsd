// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/treivax/tsd/constraint"
	"testing"
)

func TestNewASTConverter(t *testing.T) {
	converter := NewASTConverter()
	assert.NotNil(t, converter)
}
func TestASTConverter_ConvertProgram(t *testing.T) {
	converter := NewASTConverter()
	t.Run("convert empty program", func(t *testing.T) {
		constraintProgram := &constraint.Program{
			Types:       []constraint.TypeDefinition{},
			Expressions: []constraint.Expression{},
		}
		result, err := converter.ConvertProgram(constraintProgram)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Types, 0)
		assert.Len(t, result.Expressions, 0)
	})
	t.Run("convert program with types", func(t *testing.T) {
		constraintProgram := &constraint.Program{
			Types: []constraint.TypeDefinition{
				{
					Type: "typeDefinition",
					Name: "Person",
					Fields: []constraint.Field{
						{Name: "name", Type: "string"},
						{Name: "age", Type: "number"},
					},
				},
			},
			Expressions: []constraint.Expression{},
		}
		result, err := converter.ConvertProgram(constraintProgram)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Types, 1)
		typeDef := result.Types[0]
		assert.Equal(t, "typeDefinition", typeDef.Type)
		assert.Equal(t, "Person", typeDef.Name)
		assert.Len(t, typeDef.Fields, 2)
		assert.Equal(t, "name", typeDef.Fields[0].Name)
		assert.Equal(t, "string", typeDef.Fields[0].Type)
		assert.Equal(t, "age", typeDef.Fields[1].Name)
		assert.Equal(t, "number", typeDef.Fields[1].Type)
	})
	t.Run("convert program with multiple types", func(t *testing.T) {
		constraintProgram := &constraint.Program{
			Types: []constraint.TypeDefinition{
				{
					Type: "typeDefinition",
					Name: "Person",
					Fields: []constraint.Field{
						{Name: "name", Type: "string"},
					},
				},
				{
					Type: "typeDefinition",
					Name: "Company",
					Fields: []constraint.Field{
						{Name: "name", Type: "string"},
						{Name: "employees", Type: "number"},
					},
				},
			},
			Expressions: []constraint.Expression{},
		}
		result, err := converter.ConvertProgram(constraintProgram)
		require.NoError(t, err)
		assert.Len(t, result.Types, 2)
		assert.Equal(t, "Person", result.Types[0].Name)
		assert.Equal(t, "Company", result.Types[1].Name)
	})
	t.Run("convert program with expression and action", func(t *testing.T) {
		constraintProgram := &constraint.Program{
			Types: []constraint.TypeDefinition{},
			Expressions: []constraint.Expression{
				{
					Type: "expression",
					Set: constraint.Set{
						Type: "set",
						Variables: []constraint.TypedVariable{
							{
								Type:     "typedVariable",
								Name:     "p",
								DataType: "Person",
							},
						},
					},
					Constraints: map[string]interface{}{
						"type": "simple",
					},
					Action: &constraint.Action{
						Type: "action",
						Job: &constraint.JobCall{
							Type: "jobCall",
							Name: "print",
							Args: []interface{}{"Hello"},
						},
					},
				},
			},
		}
		result, err := converter.ConvertProgram(constraintProgram)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Expressions, 1)
		expr := result.Expressions[0]
		assert.Equal(t, "expression", expr.Type)
		assert.Len(t, expr.Set.Variables, 1)
		assert.Equal(t, "p", expr.Set.Variables[0].Name)
		assert.Equal(t, "Person", expr.Set.Variables[0].DataType)
		assert.NotNil(t, expr.Action)
		assert.NotNil(t, expr.Action.Job)
		assert.Equal(t, "print", expr.Action.Job.Name)
		assert.Len(t, expr.Action.Job.Args, 1)
	})
	t.Run("convert program with multiple jobs", func(t *testing.T) {
		constraintProgram := &constraint.Program{
			Types: []constraint.TypeDefinition{},
			Expressions: []constraint.Expression{
				{
					Type: "expression",
					Set: constraint.Set{
						Type: "set",
						Variables: []constraint.TypedVariable{
							{
								Type:     "typedVariable",
								Name:     "p",
								DataType: "Person",
							},
						},
					},
					Constraints: map[string]interface{}{
						"type": "simple",
					},
					Action: &constraint.Action{
						Type: "action",
						Jobs: []constraint.JobCall{
							{
								Type: "jobCall",
								Name: "log",
								Args: []interface{}{"First action"},
							},
							{
								Type: "jobCall",
								Name: "notify",
								Args: []interface{}{"Second action"},
							},
						},
					},
				},
			},
		}
		result, err := converter.ConvertProgram(constraintProgram)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Expressions, 1)
		expr := result.Expressions[0]
		assert.NotNil(t, expr.Action)
		assert.Len(t, expr.Action.Jobs, 2)
		assert.Equal(t, "log", expr.Action.Jobs[0].Name)
		assert.Equal(t, "notify", expr.Action.Jobs[1].Name)
	})
	t.Run("convert program with expression but no action", func(t *testing.T) {
		constraintProgram := &constraint.Program{
			Types: []constraint.TypeDefinition{},
			Expressions: []constraint.Expression{
				{
					Type: "expression",
					Set: constraint.Set{
						Type: "set",
						Variables: []constraint.TypedVariable{
							{
								Type:     "typedVariable",
								Name:     "p",
								DataType: "Person",
							},
						},
					},
					Constraints: map[string]interface{}{
						"type": "simple",
					},
					Action: nil, // No action
				},
			},
		}
		result, err := converter.ConvertProgram(constraintProgram)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "action manquante")
	})
	t.Run("convert program with multiple expressions", func(t *testing.T) {
		constraintProgram := &constraint.Program{
			Types: []constraint.TypeDefinition{},
			Expressions: []constraint.Expression{
				{
					Type: "expression",
					Set: constraint.Set{
						Type: "set",
						Variables: []constraint.TypedVariable{
							{Type: "typedVariable", Name: "p", DataType: "Person"},
						},
					},
					Constraints: map[string]interface{}{"type": "simple"},
					Action: &constraint.Action{
						Type: "action",
						Job:  &constraint.JobCall{Type: "jobCall", Name: "action1"},
					},
				},
				{
					Type: "expression",
					Set: constraint.Set{
						Type: "set",
						Variables: []constraint.TypedVariable{
							{Type: "typedVariable", Name: "c", DataType: "Company"},
						},
					},
					Constraints: map[string]interface{}{"type": "simple"},
					Action: &constraint.Action{
						Type: "action",
						Job:  &constraint.JobCall{Type: "jobCall", Name: "action2"},
					},
				},
			},
		}
		result, err := converter.ConvertProgram(constraintProgram)
		require.NoError(t, err)
		assert.Len(t, result.Expressions, 2)
		assert.Equal(t, "p", result.Expressions[0].Set.Variables[0].Name)
		assert.Equal(t, "c", result.Expressions[1].Set.Variables[0].Name)
	})
	t.Run("convert complete program", func(t *testing.T) {
		constraintProgram := &constraint.Program{
			Types: []constraint.TypeDefinition{
				{
					Type: "typeDefinition",
					Name: "Person",
					Fields: []constraint.Field{
						{Name: "name", Type: "string"},
						{Name: "age", Type: "number"},
					},
				},
			},
			Expressions: []constraint.Expression{
				{
					Type: "expression",
					Set: constraint.Set{
						Type: "set",
						Variables: []constraint.TypedVariable{
							{
								Type:     "typedVariable",
								Name:     "p",
								DataType: "Person",
							},
						},
					},
					Constraints: map[string]interface{}{
						"type":     "binaryOperation",
						"operator": ">",
						"left": map[string]interface{}{
							"type":   "fieldAccess",
							"object": "p",
							"field":  "age",
						},
						"right": map[string]interface{}{
							"type":  "numberLiteral",
							"value": 18.0,
						},
					},
					Action: &constraint.Action{
						Type: "action",
						Job: &constraint.JobCall{
							Type: "jobCall",
							Name: "print",
							Args: []interface{}{"Adult person"},
						},
					},
				},
			},
		}
		result, err := converter.ConvertProgram(constraintProgram)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Types, 1)
		assert.Len(t, result.Expressions, 1)
		// Verify type
		assert.Equal(t, "Person", result.Types[0].Name)
		assert.Len(t, result.Types[0].Fields, 2)
		// Verify expression
		expr := result.Expressions[0]
		assert.Equal(t, "expression", expr.Type)
		assert.NotNil(t, expr.Constraints)
		assert.NotNil(t, expr.Action)
	})
	t.Run("invalid program type", func(t *testing.T) {
		invalidProgram := "not a constraint program"
		result, err := converter.ConvertProgram(invalidProgram)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "type de programme AST non reconnu")
	})
	t.Run("nil program", func(t *testing.T) {
		result, err := converter.ConvertProgram(nil)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
func TestASTConverter_convertFields(t *testing.T) {
	converter := NewASTConverter()
	t.Run("convert empty fields", func(t *testing.T) {
		constraintFields := []constraint.Field{}
		result := converter.convertFields(constraintFields)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
	})
	t.Run("convert single field", func(t *testing.T) {
		constraintFields := []constraint.Field{
			{Name: "name", Type: "string"},
		}
		result := converter.convertFields(constraintFields)
		assert.Len(t, result, 1)
		assert.Equal(t, "name", result[0].Name)
		assert.Equal(t, "string", result[0].Type)
	})
	t.Run("convert multiple fields", func(t *testing.T) {
		constraintFields := []constraint.Field{
			{Name: "name", Type: "string"},
			{Name: "age", Type: "number"},
			{Name: "active", Type: "bool"},
		}
		result := converter.convertFields(constraintFields)
		assert.Len(t, result, 3)
		assert.Equal(t, "name", result[0].Name)
		assert.Equal(t, "age", result[1].Name)
		assert.Equal(t, "active", result[2].Name)
	})
}
func TestASTConverter_convertTypedVariables(t *testing.T) {
	converter := NewASTConverter()
	t.Run("convert empty variables", func(t *testing.T) {
		constraintVars := []constraint.TypedVariable{}
		result := converter.convertTypedVariables(constraintVars)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
	})
	t.Run("convert single variable", func(t *testing.T) {
		constraintVars := []constraint.TypedVariable{
			{
				Type:     "typedVariable",
				Name:     "p",
				DataType: "Person",
			},
		}
		result := converter.convertTypedVariables(constraintVars)
		assert.Len(t, result, 1)
		assert.Equal(t, "typedVariable", result[0].Type)
		assert.Equal(t, "p", result[0].Name)
		assert.Equal(t, "Person", result[0].DataType)
	})
	t.Run("convert multiple variables", func(t *testing.T) {
		constraintVars := []constraint.TypedVariable{
			{Type: "typedVariable", Name: "p", DataType: "Person"},
			{Type: "typedVariable", Name: "c", DataType: "Company"},
			{Type: "typedVariable", Name: "o", DataType: "Order"},
		}
		result := converter.convertTypedVariables(constraintVars)
		assert.Len(t, result, 3)
		assert.Equal(t, "p", result[0].Name)
		assert.Equal(t, "c", result[1].Name)
		assert.Equal(t, "o", result[2].Name)
	})
}
func TestASTConverter_convertAction(t *testing.T) {
	converter := NewASTConverter()
	t.Run("convert action with single job (old format)", func(t *testing.T) {
		constraintAction := constraint.Action{
			Type: "action",
			Job: &constraint.JobCall{
				Type: "jobCall",
				Name: "print",
				Args: []interface{}{"Hello", "World"},
			},
		}
		result, err := converter.convertAction(constraintAction)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "action", result.Type)
		assert.NotNil(t, result.Job)
		assert.Equal(t, "print", result.Job.Name)
		assert.Len(t, result.Job.Args, 2)
	})
	t.Run("convert action with multiple jobs (new format)", func(t *testing.T) {
		constraintAction := constraint.Action{
			Type: "action",
			Jobs: []constraint.JobCall{
				{
					Type: "jobCall",
					Name: "log",
					Args: []interface{}{"Logging"},
				},
				{
					Type: "jobCall",
					Name: "notify",
					Args: []interface{}{"Notifying"},
				},
			},
		}
		result, err := converter.convertAction(constraintAction)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "action", result.Type)
		assert.Len(t, result.Jobs, 2)
		assert.Equal(t, "log", result.Jobs[0].Name)
		assert.Equal(t, "notify", result.Jobs[1].Name)
	})
	t.Run("convert action with no jobs", func(t *testing.T) {
		constraintAction := constraint.Action{
			Type: "action",
		}
		result, err := converter.convertAction(constraintAction)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "action", result.Type)
		assert.Nil(t, result.Job)
		assert.Len(t, result.Jobs, 0)
	})
	t.Run("convert action with job and args", func(t *testing.T) {
		constraintAction := constraint.Action{
			Type: "action",
			Job: &constraint.JobCall{
				Type: "jobCall",
				Name: "calculate",
				Args: []interface{}{42, "test", true, 3.14},
			},
		}
		result, err := converter.convertAction(constraintAction)
		require.NoError(t, err)
		assert.NotNil(t, result.Job)
		assert.Len(t, result.Job.Args, 4)
		assert.Equal(t, 42, result.Job.Args[0])
		assert.Equal(t, "test", result.Job.Args[1])
		assert.Equal(t, true, result.Job.Args[2])
		assert.Equal(t, 3.14, result.Job.Args[3])
	})
}
func TestASTConverter_Integration(t *testing.T) {
	converter := NewASTConverter()
	t.Run("full pipeline conversion", func(t *testing.T) {
		// Simulate a realistic constraint program
		constraintProgram := &constraint.Program{
			Types: []constraint.TypeDefinition{
				{
					Type: "typeDefinition",
					Name: "Employee",
					Fields: []constraint.Field{
						{Name: "id", Type: "string"},
						{Name: "name", Type: "string"},
						{Name: "salary", Type: "number"},
						{Name: "active", Type: "bool"},
					},
				},
			},
			Expressions: []constraint.Expression{
				{
					Type: "expression",
					Set: constraint.Set{
						Type: "set",
						Variables: []constraint.TypedVariable{
							{
								Type:     "typedVariable",
								Name:     "e",
								DataType: "Employee",
							},
						},
					},
					Constraints: map[string]interface{}{
						"type": "logicalExpression",
						"left": map[string]interface{}{
							"type":     "binaryOperation",
							"operator": ">",
							"left": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "e",
								"field":  "salary",
							},
							"right": map[string]interface{}{
								"type":  "numberLiteral",
								"value": 50000.0,
							},
						},
						"operations": []interface{}{
							map[string]interface{}{
								"op": "AND",
								"right": map[string]interface{}{
									"type":     "binaryOperation",
									"operator": "==",
									"left": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "e",
										"field":  "active",
									},
									"right": map[string]interface{}{
										"type":  "booleanLiteral",
										"value": true,
									},
								},
							},
						},
					},
					Action: &constraint.Action{
						Type: "action",
						Jobs: []constraint.JobCall{
							{
								Type: "jobCall",
								Name: "grantBonus",
								Args: []interface{}{"e.id"},
							},
							{
								Type: "jobCall",
								Name: "log",
								Args: []interface{}{"High earner bonus granted"},
							},
						},
					},
				},
			},
		}
		result, err := converter.ConvertProgram(constraintProgram)
		require.NoError(t, err)
		assert.NotNil(t, result)
		// Verify types
		assert.Len(t, result.Types, 1)
		employeeType := result.Types[0]
		assert.Equal(t, "Employee", employeeType.Name)
		assert.Len(t, employeeType.Fields, 4)
		// Verify expressions
		assert.Len(t, result.Expressions, 1)
		expr := result.Expressions[0]
		assert.Equal(t, "expression", expr.Type)
		assert.Len(t, expr.Set.Variables, 1)
		assert.Equal(t, "e", expr.Set.Variables[0].Name)
		assert.Equal(t, "Employee", expr.Set.Variables[0].DataType)
		// Verify action
		assert.NotNil(t, expr.Action)
		assert.Len(t, expr.Action.Jobs, 2)
		assert.Equal(t, "grantBonus", expr.Action.Jobs[0].Name)
		assert.Equal(t, "log", expr.Action.Jobs[1].Name)
		// Verify constraints are preserved
		assert.NotNil(t, expr.Constraints)
	})
}
