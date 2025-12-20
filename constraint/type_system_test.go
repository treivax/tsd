// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"testing"
)

func TestTypeSystem_TypeChecks(t *testing.T) {
	t.Log("🧪 TEST TYPE SYSTEM - VÉRIFICATIONS DE TYPES")
	t.Log("=============================================")

	types := []TypeDefinition{
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
			},
		},
	}

	ts := NewTypeSystem(types)

	tests := []struct {
		name     string
		testFunc func() bool
		expected bool
	}{
		{
			name:     "string est primitif",
			testFunc: func() bool { return ts.IsPrimitiveType(ValueTypeString) },
			expected: true,
		},
		{
			name:     "User n'est pas primitif",
			testFunc: func() bool { return ts.IsPrimitiveType("User") },
			expected: false,
		},
		{
			name:     "User est user-defined",
			testFunc: func() bool { return ts.IsUserDefinedType("User") },
			expected: true,
		},
		{
			name:     "string n'est pas user-defined",
			testFunc: func() bool { return ts.IsUserDefinedType(ValueTypeString) },
			expected: false,
		},
		{
			name:     "User existe",
			testFunc: func() bool { return ts.TypeExists("User") },
			expected: true,
		},
		{
			name:     "string existe",
			testFunc: func() bool { return ts.TypeExists(ValueTypeString) },
			expected: true,
		},
		{
			name:     "Unknown n'existe pas",
			testFunc: func() bool { return ts.TypeExists("Unknown") },
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc()
			if result != tt.expected {
				t.Errorf("❌ Attendu %v, reçu %v", tt.expected, result)
			} else {
				t.Logf("✅ Vérification correcte")
			}
		})
	}
}

func TestTypeSystem_GetFieldType(t *testing.T) {
	t.Log("🧪 TEST TYPE SYSTEM - TYPE DE CHAMP")
	t.Log("====================================")

	types := []TypeDefinition{
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
			},
		},
	}

	ts := NewTypeSystem(types)

	tests := []struct {
		name      string
		typeName  string
		fieldName string
		wantType  string
		wantErr   bool
	}{
		{
			name:      "champ primitif",
			typeName:  "User",
			fieldName: "name",
			wantType:  ValueTypeString,
			wantErr:   false,
		},
		{
			name:      "champ de type fait",
			typeName:  "Login",
			fieldName: "user",
			wantType:  "User",
			wantErr:   false,
		},
		{
			name:      "champ _id_ interdit",
			typeName:  "User",
			fieldName: FieldNameInternalID,
			wantErr:   true,
		},
		{
			name:      "champ inexistant",
			typeName:  "User",
			fieldName: "unknown",
			wantErr:   true,
		},
		{
			name:      "type inexistant",
			typeName:  "Unknown",
			fieldName: "field",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldType, err := ts.GetFieldType(tt.typeName, tt.fieldName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("❌ Erreur inattendue: %v", err)
			}

			if fieldType != tt.wantType {
				t.Errorf("❌ Type attendu '%s', reçu '%s'", tt.wantType, fieldType)
			} else {
				t.Logf("✅ Type correct: %s", fieldType)
			}
		})
	}
}

func TestTypeSystem_Variables(t *testing.T) {
	t.Log("🧪 TEST TYPE SYSTEM - VARIABLES")
	t.Log("================================")

	types := []TypeDefinition{
		{
			Name: "User",
			Fields: []Field{
				{Name: "name", Type: ValueTypeString, IsPrimaryKey: true},
			},
		},
	}

	ts := NewTypeSystem(types)

	err := ts.RegisterVariable("alice", "User")
	if err != nil {
		t.Fatalf("❌ Erreur d'enregistrement: %v", err)
	}
	t.Log("✅ Variable 'alice' enregistrée")

	if !ts.VariableExists("alice") {
		t.Error("❌ Variable 'alice' devrait exister")
	}

	varType, err := ts.GetVariableType("alice")
	if err != nil {
		t.Fatalf("❌ Erreur de récupération: %v", err)
	}

	if varType != "User" {
		t.Errorf("❌ Type attendu 'User', reçu '%s'", varType)
	} else {
		t.Logf("✅ Type de variable correct: %s", varType)
	}

	_, err = ts.GetVariableType("bob")
	if err == nil {
		t.Error("❌ Attendu une erreur pour variable non définie")
	} else {
		t.Logf("✅ Erreur pour variable non définie: %v", err)
	}

	err = ts.RegisterVariable("invalid", "UnknownType")
	if err == nil {
		t.Error("❌ Attendu une erreur pour type inexistant")
	} else {
		t.Logf("✅ Erreur pour type inexistant: %v", err)
	}
}

func TestTypeSystem_CircularReferences(t *testing.T) {
	t.Log("🧪 TEST TYPE SYSTEM - RÉFÉRENCES CIRCULAIRES")
	t.Log("=============================================")

	types1 := []TypeDefinition{
		{
			Name: "User",
			Fields: []Field{
				{Name: "name", Type: ValueTypeString},
			},
		},
		{
			Name: "Login",
			Fields: []Field{
				{Name: "user", Type: "User"},
			},
		},
	}

	ts1 := NewTypeSystem(types1)
	err := ts1.ValidateCircularReferences()
	if err != nil {
		t.Errorf("❌ Pas de cycle attendu, erreur reçue: %v", err)
	} else {
		t.Log("✅ Aucun cycle détecté (correct)")
	}

	types2 := []TypeDefinition{
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
	}

	ts2 := NewTypeSystem(types2)
	err = ts2.ValidateCircularReferences()
	if err == nil {
		t.Error("❌ Cycle attendu, aucune erreur reçue")
	} else {
		t.Logf("✅ Cycle détecté: %v", err)
	}

	types3 := []TypeDefinition{
		{
			Name: "A",
			Fields: []Field{
				{Name: "b", Type: "B"},
			},
		},
		{
			Name: "B",
			Fields: []Field{
				{Name: "c", Type: "C"},
			},
		},
		{
			Name: "C",
			Fields: []Field{
				{Name: "a", Type: "A"},
			},
		},
	}

	ts3 := NewTypeSystem(types3)
	err = ts3.ValidateCircularReferences()
	if err == nil {
		t.Error("❌ Cycle indirect attendu, aucune erreur reçue")
	} else {
		t.Logf("✅ Cycle indirect détecté: %v", err)
	}
}

func TestTypeSystem_TypeCompatibility(t *testing.T) {
	t.Log("🧪 TEST TYPE SYSTEM - COMPATIBILITÉ")
	t.Log("====================================")

	types := []TypeDefinition{
		{
			Name: "User",
			Fields: []Field{
				{Name: "name", Type: ValueTypeString},
			},
		},
	}

	ts := NewTypeSystem(types)

	tests := []struct {
		name     string
		type1    string
		type2    string
		operator string
		expected bool
	}{
		{
			name:     "même type primitif",
			type1:    ValueTypeString,
			type2:    ValueTypeString,
			operator: OpEq,
			expected: true,
		},
		{
			name:     "bool et boolean compatibles",
			type1:    ValueTypeBool,
			type2:    ValueTypeBoolean,
			operator: OpEq,
			expected: true,
		},
		{
			name:     "types primitifs différents",
			type1:    ValueTypeString,
			type2:    ValueTypeNumber,
			operator: OpEq,
			expected: false,
		},
		{
			name:     "même type fait avec ==",
			type1:    "User",
			type2:    "User",
			operator: OpEq,
			expected: true,
		},
		{
			name:     "même type fait avec !=",
			type1:    "User",
			type2:    "User",
			operator: OpNeq,
			expected: true,
		},
		{
			name:     "type fait avec < interdit",
			type1:    "User",
			type2:    "User",
			operator: OpLt,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ts.AreTypesCompatible(tt.type1, tt.type2, tt.operator)
			if result != tt.expected {
				t.Errorf("❌ Attendu %v, reçu %v", tt.expected, result)
			} else {
				t.Logf("✅ Compatibilité correcte")
			}
		})
	}
}

func TestValidateTypeDefinition(t *testing.T) {
	t.Log("🧪 TEST VALIDATE TYPE DEFINITION")
	t.Log("=================================")

	tests := []struct {
		name    string
		typeDef TypeDefinition
		wantErr bool
	}{
		{
			name: "type valide",
			typeDef: TypeDefinition{
				Name: "User",
				Fields: []Field{
					{Name: "name", Type: ValueTypeString},
				},
			},
			wantErr: false,
		},
		{
			name: "nom vide",
			typeDef: TypeDefinition{
				Name:   "",
				Fields: []Field{{Name: "field", Type: ValueTypeString}},
			},
			wantErr: true,
		},
		{
			name: "pas de champs",
			typeDef: TypeDefinition{
				Name:   "Empty",
				Fields: []Field{},
			},
			wantErr: true,
		},
		{
			name: "champ _id_ interdit",
			typeDef: TypeDefinition{
				Name: "Bad",
				Fields: []Field{
					{Name: FieldNameInternalID, Type: ValueTypeString},
				},
			},
			wantErr: true,
		},
		{
			name: "champs dupliqués",
			typeDef: TypeDefinition{
				Name: "Duplicate",
				Fields: []Field{
					{Name: "field", Type: ValueTypeString},
					{Name: "field", Type: ValueTypeNumber},
				},
			},
			wantErr: true,
		},
		{
			name: "champ sans type",
			typeDef: TypeDefinition{
				Name: "NoType",
				Fields: []Field{
					{Name: "field", Type: ""},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTypeDefinition(tt.typeDef)

			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("❌ Erreur inattendue: %v", err)
				} else {
					t.Logf("✅ Validation réussie")
				}
			}
		})
	}
}
