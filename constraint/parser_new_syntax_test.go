// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

func TestParseTypeWithUserDefinedField(t *testing.T) {
	t.Log("🧪 TEST: Parsing type with user-defined field type")
	t.Log("=================================================")

	input := `type Login(user: User, #email: string, password: string)`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur de parsing: %v", err)
	}

	program, err := convertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}

	if len(program.Types) != 1 {
		t.Fatalf("❌ Attendu 1 type, reçu %d", len(program.Types))
	}

	typeDef := program.Types[0]
	if typeDef.Name != "Login" {
		t.Errorf("❌ Nom attendu 'Login', reçu '%s'", typeDef.Name)
	}

	if len(typeDef.Fields) != 3 {
		t.Fatalf("❌ Attendu 3 champs, reçu %d", len(typeDef.Fields))
	}

	// Vérifier le champ user de type User
	userField := typeDef.Fields[0]
	if userField.Name != "user" {
		t.Errorf("❌ Champ 0: attendu 'user', reçu '%s'", userField.Name)
	}
	if userField.Type != "User" {
		t.Errorf("❌ Champ 0: type attendu 'User', reçu '%s'", userField.Type)
	}
	if userField.IsPrimaryKey {
		t.Error("❌ Champ user ne devrait pas être clé primaire")
	}

	// Vérifier email (clé primaire)
	emailField := typeDef.Fields[1]
	if emailField.Name != "email" {
		t.Errorf("❌ Champ 1: attendu 'email', reçu '%s'", emailField.Name)
	}
	if emailField.Type != "string" {
		t.Errorf("❌ Champ 1: type attendu 'string', reçu '%s'", emailField.Type)
	}
	if !emailField.IsPrimaryKey {
		t.Error("❌ Champ email devrait être clé primaire")
	}

	t.Log("✅ Type avec champ utilisateur parsé correctement")
}

func TestParseFactAssignment(t *testing.T) {
	t.Log("🧪 TEST: Parsing fact assignment")
	t.Log("=================================")

	input := `alice = User(name: "Alice", age: 30)`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur de parsing: %v", err)
	}

	program, err := convertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}

	if len(program.FactAssignments) != 1 {
		t.Fatalf("❌ Attendu 1 affectation, reçu %d", len(program.FactAssignments))
	}

	assignment := program.FactAssignments[0]
	if assignment.Variable != "alice" {
		t.Errorf("❌ Variable attendue 'alice', reçu '%s'", assignment.Variable)
	}

	fact := assignment.Fact
	if fact.TypeName != "User" {
		t.Errorf("❌ Type attendu 'User', reçu '%s'", fact.TypeName)
	}

	if len(fact.Fields) != 2 {
		t.Fatalf("❌ Attendu 2 champs, reçu %d", len(fact.Fields))
	}

	t.Log("✅ Affectation de fait parsée correctement")
}

func TestParseMultipleFactAssignments(t *testing.T) {
	t.Log("🧪 TEST: Parsing multiple fact assignments")
	t.Log("==========================================")

	input := `
		alice = User(name: "Alice", age: 30)
		bob = User(name: "Bob", age: 25)
		Login(user: alice, email: "alice@example.com", password: "pass")
	`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur de parsing: %v", err)
	}

	program, err := convertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}

	if len(program.FactAssignments) != 2 {
		t.Errorf("❌ Attendu 2 affectations, reçu %d", len(program.FactAssignments))
	}

	if len(program.Facts) != 1 {
		t.Errorf("❌ Attendu 1 fait direct, reçu %d", len(program.Facts))
	}

	t.Log("✅ Affectations multiples parsées correctement")
}

func TestParseFactWithVariableReference(t *testing.T) {
	t.Log("🧪 TEST: Parsing fact with variable reference")
	t.Log("==============================================")

	input := `Login(user: alice, email: "alice@example.com", password: "pass")`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur de parsing: %v", err)
	}

	program, err := convertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}

	if len(program.Facts) != 1 {
		t.Fatalf("❌ Attendu 1 fait, reçu %d", len(program.Facts))
	}

	fact := program.Facts[0]
	if len(fact.Fields) != 3 {
		t.Fatalf("❌ Attendu 3 champs, reçu %d", len(fact.Fields))
	}

	// Premier champ est une référence de variable
	userField := fact.Fields[0]
	if userField.Value.Type != "variableReference" {
		t.Errorf("❌ Type attendu 'variableReference', reçu '%s'", userField.Value.Type)
	}

	varName, ok := userField.Value.Value.(string)
	if !ok || varName != "alice" {
		t.Errorf("❌ Nom de variable attendu 'alice', reçu '%v'", userField.Value.Value)
	}

	t.Log("✅ Fait avec référence variable parsé correctement")
}

func TestParseType_InternalIDForbidden(t *testing.T) {
	t.Log("🧪 TEST: Interdiction de _id_ dans type")
	t.Log("========================================")

	input := `type User(_id_: string, name: string)`

	_, err := Parse("test", []byte(input))

	if err == nil {
		t.Fatal("❌ Attendu une erreur pour champ _id_")
	}

	if !strings.Contains(err.Error(), "réservé") {
		t.Errorf("❌ Message d'erreur attendu contenant 'réservé', reçu: %v", err)
	}

	t.Log("✅ _id_ correctement rejeté dans type")
}

func TestParseFact_InternalIDForbidden(t *testing.T) {
	t.Log("🧪 TEST: Interdiction de _id_ dans fact")
	t.Log("========================================")

	input := `User(_id_: "manual", name: "Alice")`

	_, err := Parse("test", []byte(input))

	if err == nil {
		t.Fatal("❌ Attendu une erreur pour affectation _id_")
	}

	if !strings.Contains(err.Error(), "réservé") {
		t.Errorf("❌ Message d'erreur attendu contenant 'réservé', reçu: %v", err)
	}

	t.Log("✅ _id_ correctement rejeté dans fact")
}

func TestParseFactComparison(t *testing.T) {
	t.Log("🧪 TEST: Parsing fact comparison in rule")
	t.Log("=========================================")

	input := `rule test_rule: {u: User, l: Login} / l.user == u ==> Log("test")`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur de parsing: %v", err)
	}

	program, err := convertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}

	if len(program.Expressions) != 1 {
		t.Fatalf("❌ Attendu 1 expression, reçu %d", len(program.Expressions))
	}

	t.Log("✅ Syntaxe p.user == u acceptée par le parser")
}

func TestValidateTypeReferences(t *testing.T) {
	t.Log("🧪 TEST: Validation des références de types")
	t.Log("============================================")

	tests := []struct {
		name    string
		program Program
		wantErr bool
	}{
		{
			name: "type valide référencé",
			program: Program{
				Types: []TypeDefinition{
					{
						Name: "User",
						Fields: []Field{
							{Name: "name", Type: "string", IsPrimaryKey: true},
						},
					},
					{
						Name: "Login",
						Fields: []Field{
							{Name: "user", Type: "User"},
							{Name: "password", Type: "string"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "type inconnu référencé",
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
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTypeReferences(tt.program)
			if (err != nil) != tt.wantErr {
				t.Errorf("❌ wantErr %v, got err %v", tt.wantErr, err)
			} else {
				t.Logf("✅ %s validé", tt.name)
			}
		})
	}
}

func TestValidateCircularReferences(t *testing.T) {
	t.Log("🧪 TEST: Validation références circulaires")
	t.Log("===========================================")

	program := Program{
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
	}

	err := validateNoCircularReferences(program)
	if err == nil {
		t.Error("❌ Attendu une erreur pour référence circulaire")
	} else {
		t.Log("✅ Référence circulaire détectée correctement")
	}
}

func TestParseAndValidate_Complete(t *testing.T) {
	t.Log("🧪 TEST: Parsing et validation complet")
	t.Log("======================================")

	input := `
		type User(#name: string, age: number)
		type Login(user: User, #email: string, password: string)
		
		alice = User(name: "Alice", age: 30)
		bob = User(name: "Bob", age: 25)
		
		Login(user: alice, email: "alice@example.com", password: "pass123")
		Login(user: bob, email: "bob@example.com", password: "secret")
	`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur de parsing: %v", err)
	}

	program, err := convertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}

	// Validation manuelle (sans ValidateProgram pour avoir un contrôle fin)
	if err := validateTypeReferences(program); err != nil {
		t.Fatalf("❌ Erreur validation types: %v", err)
	}

	if err := validateNoCircularReferences(program); err != nil {
		t.Fatalf("❌ Erreur validation circulaire: %v", err)
	}

	if err := validateVariableReferences(program); err != nil {
		t.Fatalf("❌ Erreur validation variables: %v", err)
	}

	// Vérifications
	if len(program.Types) != 2 {
		t.Errorf("❌ Attendu 2 types, reçu %d", len(program.Types))
	}

	if len(program.FactAssignments) != 2 {
		t.Errorf("❌ Attendu 2 affectations, reçu %d", len(program.FactAssignments))
	}

	if len(program.Facts) != 2 {
		t.Errorf("❌ Attendu 2 faits, reçu %d", len(program.Facts))
	}

	t.Log("✅ Programme complet parsé et validé")
}

func TestValidateVariableReferences_Undefined(t *testing.T) {
	t.Log("🧪 TEST: Validation variable non définie")
	t.Log("=========================================")

	program := Program{
		Types: []TypeDefinition{
			{
				Name: "User",
				Fields: []Field{
					{Name: "name", Type: "string"},
				},
			},
			{
				Name: "Login",
				Fields: []Field{
					{Name: "user", Type: "User"}, // Custom type, not primitive
					{Name: "email", Type: "string"},
				},
			},
		},
		Facts: []Fact{
			{
				TypeName: "Login",
				Fields: []FactField{
					{
						Name: "user",
						Value: FactValue{
							Type:  "variableReference",
							Value: "unknownVar",
						},
					},
				},
			},
		},
	}

	err := validateVariableReferences(program)
	if err == nil {
		t.Error("❌ Attendu une erreur pour variable non définie")
	} else {
		if !strings.Contains(err.Error(), "non définie") {
			t.Errorf("❌ Message d'erreur incorrect: %v", err)
		} else {
			t.Log("✅ Variable non définie détectée correctement")
		}
	}
}
