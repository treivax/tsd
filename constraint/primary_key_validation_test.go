// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

func TestValidatePrimaryKeyField(t *testing.T) {
	t.Log("🧪 TEST VALIDATE PRIMARY KEY FIELD")
	t.Log("===================================")

	tests := []struct {
		name     string
		field    Field
		typeName string
		wantErr  bool
		errMsg   string
	}{
		{
			name: "string PK - valide",
			field: Field{
				Name:         "id",
				Type:         ValueTypeString,
				IsPrimaryKey: true,
			},
			typeName: "User",
			wantErr:  false,
		},
		{
			name: "number PK - valide",
			field: Field{
				Name:         "code",
				Type:         ValueTypeNumber,
				IsPrimaryKey: true,
			},
			typeName: "Product",
			wantErr:  false,
		},
		{
			name: "bool PK - valide",
			field: Field{
				Name:         "flag",
				Type:         ValueTypeBool,
				IsPrimaryKey: true,
			},
			typeName: "Flag",
			wantErr:  false,
		},
		{
			name: "boolean PK - valide",
			field: Field{
				Name:         "active",
				Type:         ValueTypeBoolean,
				IsPrimaryKey: true,
			},
			typeName: "Status",
			wantErr:  false,
		},
		{
			name: "champ non-PK - toujours valide",
			field: Field{
				Name:         "data",
				Type:         "CustomType",
				IsPrimaryKey: false,
			},
			typeName: "Record",
			wantErr:  false,
		},
		{
			name: "type complexe PK - invalide",
			field: Field{
				Name:         "obj",
				Type:         "CustomObject",
				IsPrimaryKey: true,
			},
			typeName: "Entity",
			wantErr:  true,
			errMsg:   "doivent être de type primitif",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePrimaryKeyField(tt.field, tt.typeName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("❌ Message d'erreur attendu contenant '%s', reçu '%s'",
						tt.errMsg, err.Error())
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("❌ Erreur inattendue: %v", err)
				} else {
					t.Log("✅ Test réussi")
				}
			}
		})
	}
}

func TestValidateTypePrimaryKey(t *testing.T) {
	t.Log("🧪 TEST VALIDATE TYPE PRIMARY KEY")
	t.Log("==================================")

	tests := []struct {
		name    string
		typeDef TypeDefinition
		wantErr bool
		errMsg  string
	}{
		{
			name: "sans clé primaire - valide",
			typeDef: TypeDefinition{
				Name: "Document",
				Fields: []Field{
					{Name: "title", Type: ValueTypeString, IsPrimaryKey: false},
					{Name: "content", Type: ValueTypeString, IsPrimaryKey: false},
				},
			},
			wantErr: false,
		},
		{
			name: "clé primaire string - valide",
			typeDef: TypeDefinition{
				Name: "User",
				Fields: []Field{
					{Name: "login", Type: ValueTypeString, IsPrimaryKey: true},
					{Name: "name", Type: ValueTypeString, IsPrimaryKey: false},
				},
			},
			wantErr: false,
		},
		{
			name: "clé primaire composite - valide",
			typeDef: TypeDefinition{
				Name: "Person",
				Fields: []Field{
					{Name: "firstName", Type: ValueTypeString, IsPrimaryKey: true},
					{Name: "lastName", Type: ValueTypeString, IsPrimaryKey: true},
					{Name: "age", Type: ValueTypeNumber, IsPrimaryKey: false},
				},
			},
			wantErr: false,
		},
		{
			name: "clé primaire type complexe - invalide",
			typeDef: TypeDefinition{
				Name: "Entity",
				Fields: []Field{
					{Name: "complexKey", Type: "CustomObject", IsPrimaryKey: true},
					{Name: "data", Type: ValueTypeString, IsPrimaryKey: false},
				},
			},
			wantErr: true,
			errMsg:  "type primitif",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTypePrimaryKey(tt.typeDef)

			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("❌ Message d'erreur attendu contenant '%s', reçu '%s'",
						tt.errMsg, err.Error())
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("❌ Erreur inattendue: %v", err)
				} else {
					t.Log("✅ Test réussi")
				}
			}
		})
	}
}

func TestValidateFactPrimaryKey(t *testing.T) {
	t.Log("🧪 TEST VALIDATE FACT PRIMARY KEY")
	t.Log("==================================")

	tests := []struct {
		name    string
		fact    Fact
		typeDef TypeDefinition
		wantErr bool
		errMsg  string
	}{
		{
			name: "fait valide sans PK",
			fact: Fact{
				TypeName: "Document",
				Fields: []FactField{
					{Name: "title", Value: FactValue{Type: ValueTypeString, Value: "Doc1"}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Document",
				Fields: []Field{
					{Name: "title", Type: ValueTypeString, IsPrimaryKey: false},
				},
			},
			wantErr: false,
		},
		{
			name: "fait valide avec PK simple",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: "login", Value: FactValue{Type: ValueTypeString, Value: "alice"}},
					{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
				},
			},
			typeDef: TypeDefinition{
				Name: "User",
				Fields: []Field{
					{Name: "login", Type: ValueTypeString, IsPrimaryKey: true},
					{Name: "name", Type: ValueTypeString, IsPrimaryKey: false},
				},
			},
			wantErr: false,
		},
		{
			name: "fait avec _id_ manuel - invalide",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: FieldNameInternalID, Value: FactValue{Type: ValueTypeString, Value: "manual-id"}},
					{Name: "login", Value: FactValue{Type: ValueTypeString, Value: "alice"}},
				},
			},
			typeDef: TypeDefinition{
				Name: "User",
				Fields: []Field{
					{Name: "login", Type: ValueTypeString, IsPrimaryKey: true},
					{Name: "name", Type: ValueTypeString, IsPrimaryKey: false},
				},
			},
			wantErr: true,
			errMsg:  "réservé au système",
		},
		{
			name: "fait sans champ PK - invalide",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
				},
			},
			typeDef: TypeDefinition{
				Name: "User",
				Fields: []Field{
					{Name: "login", Type: ValueTypeString, IsPrimaryKey: true},
					{Name: "name", Type: ValueTypeString, IsPrimaryKey: false},
				},
			},
			wantErr: true,
			errMsg:  "manquants",
		},
		{
			name: "fait avec PK composite complet - valide",
			fact: Fact{
				TypeName: "Person",
				Fields: []FactField{
					{Name: "firstName", Value: FactValue{Type: ValueTypeString, Value: "John"}},
					{Name: "lastName", Value: FactValue{Type: ValueTypeString, Value: "Doe"}},
					{Name: "age", Value: FactValue{Type: ValueTypeNumber, Value: float64(30)}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Person",
				Fields: []Field{
					{Name: "firstName", Type: ValueTypeString, IsPrimaryKey: true},
					{Name: "lastName", Type: ValueTypeString, IsPrimaryKey: true},
					{Name: "age", Type: ValueTypeNumber, IsPrimaryKey: false},
				},
			},
			wantErr: false,
		},
		{
			name: "fait avec PK composite partiel - invalide",
			fact: Fact{
				TypeName: "Person",
				Fields: []FactField{
					{Name: "firstName", Value: FactValue{Type: ValueTypeString, Value: "John"}},
					{Name: "age", Value: FactValue{Type: ValueTypeNumber, Value: float64(30)}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Person",
				Fields: []Field{
					{Name: "firstName", Type: ValueTypeString, IsPrimaryKey: true},
					{Name: "lastName", Type: ValueTypeString, IsPrimaryKey: true},
					{Name: "age", Type: ValueTypeNumber, IsPrimaryKey: false},
				},
			},
			wantErr: true,
			errMsg:  "lastName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFactPrimaryKey(tt.fact, tt.typeDef)

			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("❌ Message d'erreur attendu contenant '%s', reçu '%s'",
						tt.errMsg, err.Error())
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("❌ Erreur inattendue: %v", err)
				} else {
					t.Log("✅ Test réussi")
				}
			}
		})
	}
}

func TestValidateFactPrimaryKeyValues(t *testing.T) {
	t.Log("🧪 TEST VALIDATE FACT PRIMARY KEY VALUES")
	t.Log("=========================================")

	typeDef := TypeDefinition{
		Name: "User",
		Fields: []Field{
			{Name: "login", Type: ValueTypeString, IsPrimaryKey: true},
			{Name: "name", Type: ValueTypeString, IsPrimaryKey: false},
		},
	}

	tests := []struct {
		name    string
		fact    Fact
		wantErr bool
		errMsg  string
	}{
		{
			name: "valeur PK valide",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: "login", Value: FactValue{Type: ValueTypeString, Value: "alice"}},
					{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
				},
			},
			wantErr: false,
		},
		{
			name: "valeur PK nulle - invalide",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: "login", Value: FactValue{Type: ValueTypeString, Value: nil}},
					{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
				},
			},
			wantErr: true,
			errMsg:  "ne peut pas être nul",
		},
		{
			name: "valeur PK string vide - invalide",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: "login", Value: FactValue{Type: ValueTypeString, Value: ""}},
					{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
				},
			},
			wantErr: true,
			errMsg:  "ne peut pas être vide",
		},
		{
			name: "valeur PK identifier vide - invalide",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: "login", Value: FactValue{Type: ValueTypeIdentifier, Value: ""}},
					{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
				},
			},
			wantErr: true,
			errMsg:  "ne peut pas être vide",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFactPrimaryKeyValues(tt.fact, typeDef)

			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("❌ Message d'erreur attendu contenant '%s', reçu '%s'",
						tt.errMsg, err.Error())
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("❌ Erreur inattendue: %v", err)
				} else {
					t.Log("✅ Test réussi")
				}
			}
		})
	}
}

func TestValidatePrimaryKeyNumberField(t *testing.T) {
	t.Log("🧪 TEST VALIDATE PRIMARY KEY NUMBER FIELD")
	t.Log("==========================================")

	typeDef := TypeDefinition{
		Name: "Product",
		Fields: []Field{
			{Name: "code", Type: ValueTypeNumber, IsPrimaryKey: true},
			{Name: "name", Type: ValueTypeString, IsPrimaryKey: false},
		},
	}

	tests := []struct {
		name    string
		fact    Fact
		wantErr bool
		errMsg  string
	}{
		{
			name: "valeur number PK valide",
			fact: Fact{
				TypeName: "Product",
				Fields: []FactField{
					{Name: "code", Value: FactValue{Type: ValueTypeNumber, Value: float64(123)}},
					{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Widget"}},
				},
			},
			wantErr: false,
		},
		{
			name: "valeur number PK zéro - valide",
			fact: Fact{
				TypeName: "Product",
				Fields: []FactField{
					{Name: "code", Value: FactValue{Type: ValueTypeNumber, Value: float64(0)}},
					{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Zero"}},
				},
			},
			wantErr: false,
		},
		{
			name: "valeur number PK nil - invalide",
			fact: Fact{
				TypeName: "Product",
				Fields: []FactField{
					{Name: "code", Value: FactValue{Type: ValueTypeNumber, Value: nil}},
					{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Null"}},
				},
			},
			wantErr: true,
			errMsg:  "ne peut pas être nul",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFactPrimaryKeyValues(tt.fact, typeDef)

			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("❌ Message d'erreur attendu contenant '%s', reçu '%s'",
						tt.errMsg, err.Error())
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("❌ Erreur inattendue: %v", err)
				} else {
					t.Log("✅ Test réussi")
				}
			}
		})
	}
}
