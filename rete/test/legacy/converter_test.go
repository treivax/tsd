package rete

import (
	"testing"

	"github.com/treivax/tsd/constraint"
)

// TestASTConverter teste la création et utilisation du convertisseur AST
func TestASTConverter(t *testing.T) {
	// Test NewASTConverter
	converter := NewASTConverter()
	if converter == nil {
		t.Fatal("NewASTConverter devrait retourner un convertisseur non-nil")
	}

	// Test ConvertProgram avec programme valide
	constraintProgram := &constraint.Program{
		Types: []constraint.TypeDefinition{
			{
				Type: "struct",
				Name: "Person",
				Fields: []constraint.Field{
					{Name: "name", Type: "string"},
					{Name: "age", Type: "int"},
				},
			},
		},
		Expressions: []constraint.Expression{
			{
				Type: "rule",
				Set: constraint.Set{
					Type: "fact",
					Variables: []constraint.TypedVariable{
						{Type: "var", Name: "p", DataType: "Person"},
					},
				},
				Constraints: []interface{}{
					map[string]interface{}{
						"operator": "==",
						"left":     map[string]interface{}{"type": "field_access", "object": "p", "field": "name"},
						"right":    map[string]interface{}{"type": "literal", "value": "Alice"},
					},
				},
				Action: &constraint.Action{
					Type: "job_call",
					Job: constraint.JobCall{
						Type: "call",
						Name: "print",
						Args: []string{"Match found"},
					},
				},
			},
		},
	}

	result, err := converter.ConvertProgram(constraintProgram)
	if err != nil {
		t.Fatalf("ConvertProgram devrait réussir: %v", err)
	}

	if result == nil {
		t.Fatal("Le résultat ne devrait pas être nil")
	}

	// Vérifier la conversion des types
	if len(result.Types) != 1 {
		t.Errorf("Attendu 1 type, obtenu %d", len(result.Types))
	}

	if result.Types[0].Name != "Person" {
		t.Errorf("Attendu nom 'Person', obtenu '%s'", result.Types[0].Name)
	}

	if len(result.Types[0].Fields) != 2 {
		t.Errorf("Attendu 2 champs, obtenu %d", len(result.Types[0].Fields))
	}

	// Vérifier la conversion des expressions
	if len(result.Expressions) != 1 {
		t.Errorf("Attendu 1 expression, obtenu %d", len(result.Expressions))
	}

	if result.Expressions[0].Type != "rule" {
		t.Errorf("Attendu type 'rule', obtenu '%s'", result.Expressions[0].Type)
	}

	// Vérifier l'action convertie
	if result.Expressions[0].Action == nil {
		t.Error("L'action ne devrait pas être nil")
	} else {
		if result.Expressions[0].Action.Job.Name != "print" {
			t.Errorf("Attendu job name 'print', obtenu '%s'", result.Expressions[0].Action.Job.Name)
		}
	}
}

// TestConvertProgramInvalidType teste ConvertProgram avec type invalide
func TestConvertProgramInvalidType(t *testing.T) {
	converter := NewASTConverter()

	// Tester avec un type invalide
	invalidProgram := "not a constraint program"
	_, err := converter.ConvertProgram(invalidProgram)

	if err == nil {
		t.Error("ConvertProgram devrait échouer avec un type invalide")
	}

	expectedError := "type de programme AST non reconnu"
	if err.Error() != expectedError {
		t.Errorf("Attendu erreur '%s', obtenu '%s'", expectedError, err.Error())
	}
}

// TestConvertFields teste la conversion des champs
func TestConvertFields(t *testing.T) {
	converter := NewASTConverter()

	constraintFields := []constraint.Field{
		{Name: "id", Type: "int"},
		{Name: "name", Type: "string"},
		{Name: "active", Type: "bool"},
	}

	result := converter.convertFields(constraintFields)

	if len(result) != 3 {
		t.Errorf("Attendu 3 champs, obtenu %d", len(result))
	}

	expectedFields := []Field{
		{Name: "id", Type: "int"},
		{Name: "name", Type: "string"},
		{Name: "active", Type: "bool"},
	}

	for i, expected := range expectedFields {
		if result[i].Name != expected.Name {
			t.Errorf("Champ %d: attendu nom '%s', obtenu '%s'", i, expected.Name, result[i].Name)
		}
		if result[i].Type != expected.Type {
			t.Errorf("Champ %d: attendu type '%s', obtenu '%s'", i, expected.Type, result[i].Type)
		}
	}
}

// TestConvertTypedVariables teste la conversion des variables typées
func TestConvertTypedVariables(t *testing.T) {
	converter := NewASTConverter()

	constraintVars := []constraint.TypedVariable{
		{Type: "var", Name: "person", DataType: "Person"},
		{Type: "var", Name: "order", DataType: "Order"},
	}

	result := converter.convertTypedVariables(constraintVars)

	if len(result) != 2 {
		t.Errorf("Attendu 2 variables, obtenu %d", len(result))
	}

	// Vérifier première variable
	if result[0].Name != "person" || result[0].DataType != "Person" {
		t.Errorf("Variable 0: attendu person/Person, obtenu %s/%s", result[0].Name, result[0].DataType)
	}

	// Vérifier deuxième variable
	if result[1].Name != "order" || result[1].DataType != "Order" {
		t.Errorf("Variable 1: attendu order/Order, obtenu %s/%s", result[1].Name, result[1].DataType)
	}
}

// TestConvertAction teste la conversion des actions
func TestConvertAction(t *testing.T) {
	converter := NewASTConverter()

	constraintAction := constraint.Action{
		Type: "job_call",
		Job: constraint.JobCall{
			Type: "call",
			Name: "notify",
			Args: []string{"user", "message"},
		},
	}

	result, err := converter.convertAction(constraintAction)
	if err != nil {
		t.Fatalf("convertAction ne devrait pas échouer: %v", err)
	}

	if result == nil {
		t.Fatal("Le résultat ne devrait pas être nil")
	}

	if result.Type != "job_call" {
		t.Errorf("Attendu type 'job_call', obtenu '%s'", result.Type)
	}

	if result.Job.Name != "notify" {
		t.Errorf("Attendu job name 'notify', obtenu '%s'", result.Job.Name)
	}

	if len(result.Job.Args) != 2 {
		t.Errorf("Attendu 2 arguments, obtenu %d", len(result.Job.Args))
	}

	expectedArgs := []string{"user", "message"}
	for i, expected := range expectedArgs {
		if result.Job.Args[i] != expected {
			t.Errorf("Argument %d: attendu '%s', obtenu '%s'", i, expected, result.Job.Args[i])
		}
	}
}

// TestConvertExpressionWithoutAction teste une expression sans action
func TestConvertExpressionWithoutAction(t *testing.T) {
	converter := NewASTConverter()

	constraintExpr := constraint.Expression{
		Type: "rule",
		Set: constraint.Set{
			Type: "fact",
			Variables: []constraint.TypedVariable{
				{Type: "var", Name: "x", DataType: "Number"},
			},
		},
		Constraints: []interface{}{
			map[string]interface{}{
				"operator": ">",
				"left":     "x",
				"right":    10,
			},
		},
		Action: nil, // Pas d'action
	}

	result, err := converter.convertExpression(constraintExpr)
	if err != nil {
		t.Fatalf("convertExpression ne devrait pas échouer: %v", err)
	}

	if result == nil {
		t.Fatal("Le résultat ne devrait pas être nil")
	}

	if result.Action != nil {
		t.Error("L'action devrait être nil")
	}

	// Vérifier le set
	if result.Set.Type != "fact" {
		t.Errorf("Attendu set type 'fact', obtenu '%s'", result.Set.Type)
	}

	if len(result.Set.Variables) != 1 {
		t.Errorf("Attendu 1 variable, obtenu %d", len(result.Set.Variables))
	}
}
