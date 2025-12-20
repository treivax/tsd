// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

func TestProgramValidator_ValidProgram(t *testing.T) {
	t.Log("🧪 TEST PROGRAM VALIDATOR - PROGRAMME VALIDE")
	t.Log("=============================================")

	program := Program{
		Types: []TypeDefinition{
			{
				Name: "User",
				Fields: []Field{
					{Name: "name", Type: ValueTypeString, IsPrimaryKey: true},
					{Name: "age", Type: ValueTypeNumber},
				},
			},
			{
				Name: "Login",
				Fields: []Field{
					{Name: "user", Type: "User"},
					{Name: "email", Type: ValueTypeString, IsPrimaryKey: true},
					{Name: "password", Type: ValueTypeString},
				},
			},
		},
		FactAssignments: []FactAssignment{
			{
				Variable: "alice",
				Fact: Fact{
					TypeName: "User",
					Fields: []FactField{
						{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
						{Name: "age", Value: FactValue{Type: ValueTypeNumber, Value: 30.0}},
					},
				},
			},
		},
		Facts: []Fact{
			{
				TypeName: "Login",
				Fields: []FactField{
					{Name: "user", Value: FactValue{Type: "variableReference", Value: "alice"}},
					{Name: "email", Value: FactValue{Type: ValueTypeString, Value: "alice@ex.com"}},
					{Name: "password", Value: FactValue{Type: ValueTypeString, Value: "secret"}},
				},
			},
		},
	}

	validator := NewProgramValidator()
	err := validator.Validate(program)

	if err != nil {
		t.Errorf("❌ Erreur inattendue: %v", err)
	} else {
		t.Log("✅ Programme valide")
	}
}

func TestProgramValidator_InvalidPrograms(t *testing.T) {
	t.Log("🧪 TEST PROGRAM VALIDATOR - PROGRAMMES INVALIDES")
	t.Log("=================================================")

	tests := []struct {
		name    string
		program Program
		errMsg  string
	}{
		{
			name: "référence circulaire",
			program: Program{
				Types: []TypeDefinition{
					{
						Name: "A",
						Fields: []Field{
							{Name: "b", Type: "B"},
						},
					},
					{
						Name: "B",
						Fields: []Field{
							{Name: "a", Type: "A"},
						},
					},
				},
			},
			errMsg: "circulaire",
		},
		{
			name: "type inexistant",
			program: Program{
				Types: []TypeDefinition{
					{
						Name: "Login",
						Fields: []Field{
							{Name: "user", Type: "UnknownType"},
						},
					},
				},
			},
			errMsg: "UnknownType",
		},
		{
			name: "variable non définie",
			program: Program{
				Types: []TypeDefinition{
					{
						Name: "Login",
						Fields: []Field{
							{Name: "user", Type: "User"},
							{Name: "email", Type: ValueTypeString, IsPrimaryKey: true},
							{Name: "password", Type: ValueTypeString},
						},
					},
					{
						Name: "User",
						Fields: []Field{
							{Name: "name", Type: ValueTypeString, IsPrimaryKey: true},
						},
					},
				},
				Facts: []Fact{
					{
						TypeName: "Login",
						Fields: []FactField{
							{Name: "user", Value: FactValue{Type: "variableReference", Value: "unknown"}},
							{Name: "email", Value: FactValue{Type: ValueTypeString, Value: "test@ex.com"}},
							{Name: "password", Value: FactValue{Type: ValueTypeString, Value: "pw"}},
						},
					},
				},
			},
			errMsg: "non définie",
		},
		{
			name: "type définition invalide - nom vide",
			program: Program{
				Types: []TypeDefinition{
					{
						Name:   "",
						Fields: []Field{{Name: "field", Type: ValueTypeString}},
					},
				},
			},
			errMsg: "nom du type ne peut pas être vide",
		},
		{
			name: "champ manquant dans fait",
			program: Program{
				Types: []TypeDefinition{
					{
						Name: "User",
						Fields: []Field{
							{Name: "name", Type: ValueTypeString, IsPrimaryKey: true},
							{Name: "age", Type: ValueTypeNumber},
						},
					},
				},
				Facts: []Fact{
					{
						TypeName: "User",
						Fields: []FactField{
							{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
						},
					},
				},
			},
			errMsg: "manquant",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewProgramValidator()
			err := validator.Validate(tt.program)

			if err == nil {
				t.Errorf("❌ Attendu une erreur, reçu nil")
			} else {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("❌ Erreur attendue contenant '%s', reçu: %v", tt.errMsg, err)
				} else {
					t.Logf("✅ Erreur détectée: %v", err)
				}
			}
		})
	}
}

func TestProgramValidator_ValidateExpression(t *testing.T) {
	t.Log("🧪 TEST PROGRAM VALIDATOR - VALIDATION EXPRESSION")
	t.Log("==================================================")

	program := Program{
		Types: []TypeDefinition{
			{
				Name: "User",
				Fields: []Field{
					{Name: "name", Type: ValueTypeString},
					{Name: "age", Type: ValueTypeNumber},
				},
			},
		},
	}

	validator := NewProgramValidator()
	validator.typeSystem = NewTypeSystem(program.Types)

	expr := Expression{
		RuleId: "test_rule",
		Set: Set{
			Type: "set",
			Variables: []TypedVariable{
				{Type: "typedVariable", Name: "u", DataType: "User"},
			},
		},
		Constraints: map[string]interface{}{
			"type":     ConstraintTypeComparison,
			"operator": OpGt,
			"left": map[string]interface{}{
				"type":   ConstraintTypeFieldAccess,
				"object": "u",
				"field":  "age",
			},
			"right": map[string]interface{}{
				"type":  ValueTypeNumber,
				"value": 18.0,
			},
		},
	}

	err := validator.validateExpression(expr)
	if err != nil {
		t.Errorf("❌ Erreur inattendue: %v", err)
	} else {
		t.Log("✅ Expression valide")
	}
}

func TestProgramValidator_InvalidExpression(t *testing.T) {
	t.Log("🧪 TEST PROGRAM VALIDATOR - EXPRESSION INVALIDE")
	t.Log("================================================")

	program := Program{
		Types: []TypeDefinition{
			{
				Name: "User",
				Fields: []Field{
					{Name: "name", Type: ValueTypeString},
					{Name: "age", Type: ValueTypeNumber},
				},
			},
		},
	}

	validator := NewProgramValidator()
	validator.typeSystem = NewTypeSystem(program.Types)

	expr := Expression{
		RuleId: "test_rule",
		Set: Set{
			Type: "set",
			Variables: []TypedVariable{
				{Type: "typedVariable", Name: "u", DataType: "User"},
			},
		},
		Constraints: map[string]interface{}{
			"type":     ConstraintTypeComparison,
			"operator": OpGt,
			"left": map[string]interface{}{
				"type":   ConstraintTypeFieldAccess,
				"object": "u",
				"field":  "name",
			},
			"right": map[string]interface{}{
				"type":  ValueTypeNumber,
				"value": 18.0,
			},
		},
	}

	err := validator.validateExpression(expr)
	if err == nil {
		t.Error("❌ Attendu une erreur pour types incompatibles")
	} else {
		if !strings.Contains(err.Error(), "incompatibles") {
			t.Errorf("❌ Message attendu avec 'incompatibles', reçu: %v", err)
		} else {
			t.Logf("✅ Erreur détectée: %v", err)
		}
	}
}
