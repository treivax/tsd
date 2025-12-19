// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

func TestFactValidator_ValidateFact(t *testing.T) {
	t.Log("🧪 TEST FACT VALIDATOR - VALIDATION COMPLÈTE")
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
				{Name: "password", Type: ValueTypeString},
			},
		},
	}

	ts := NewTypeSystem(types)
	ts.RegisterVariable("alice", "User")

	validator := NewFactValidator(ts)

	tests := []struct {
		name    string
		fact    Fact
		wantErr bool
		errMsg  string
	}{
		{
			name: "fait valide avec primitifs",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
					{Name: "age", Value: FactValue{Type: ValueTypeNumber, Value: 30.0}},
				},
			},
			wantErr: false,
		},
		{
			name: "fait valide avec variable",
			fact: Fact{
				TypeName: "Login",
				Fields: []FactField{
					{Name: "user", Value: FactValue{Type: "variableReference", Value: "alice"}},
					{Name: "email", Value: FactValue{Type: ValueTypeString, Value: "alice@ex.com"}},
					{Name: "password", Value: FactValue{Type: ValueTypeString, Value: "secret"}},
				},
			},
			wantErr: false,
		},
		{
			name: "type inexistant",
			fact: Fact{
				TypeName: "UnknownType",
				Fields:   []FactField{},
			},
			wantErr: true,
			errMsg:  "non défini",
		},
		{
			name: "champ manquant",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
				},
			},
			wantErr: true,
			errMsg:  "manquant",
		},
		{
			name: "champ non défini",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
					{Name: "age", Value: FactValue{Type: ValueTypeNumber, Value: 30.0}},
					{Name: "unknown", Value: FactValue{Type: ValueTypeString, Value: "test"}},
				},
			},
			wantErr: true,
			errMsg:  "non défini dans le type",
		},
		{
			name: "champ _id_ interdit",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: FieldNameInternalID, Value: FactValue{Type: ValueTypeString, Value: "manual"}},
					{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
					{Name: "age", Value: FactValue{Type: ValueTypeNumber, Value: 30.0}},
				},
			},
			wantErr: true,
			errMsg:  "réservé",
		},
		{
			name: "type de valeur incompatible",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
					{Name: "age", Value: FactValue{Type: ValueTypeString, Value: "thirty"}},
				},
			},
			wantErr: true,
			errMsg:  "type incompatible",
		},
		{
			name: "variable non définie",
			fact: Fact{
				TypeName: "Login",
				Fields: []FactField{
					{Name: "user", Value: FactValue{Type: "variableReference", Value: "bob"}},
					{Name: "email", Value: FactValue{Type: ValueTypeString, Value: "test@ex.com"}},
					{Name: "password", Value: FactValue{Type: ValueTypeString, Value: "secret"}},
				},
			},
			wantErr: true,
			errMsg:  "non définie",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateFact(tt.fact)

			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else {
					if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
						t.Errorf("❌ Message d'erreur attendu contenant '%s', reçu: %v", tt.errMsg, err)
					} else {
						t.Logf("✅ Erreur attendue: %v", err)
					}
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

func TestFactValidator_ValidateFieldValue(t *testing.T) {
	t.Log("🧪 TEST FACT VALIDATOR - VALIDATION VALEUR CHAMP")
	t.Log("=================================================")

	types := []TypeDefinition{
		{
			Name: "TestType",
			Fields: []Field{
				{Name: "strField", Type: ValueTypeString},
				{Name: "numField", Type: ValueTypeNumber},
				{Name: "boolField", Type: ValueTypeBool},
			},
		},
	}

	ts := NewTypeSystem(types)
	validator := NewFactValidator(ts)

	tests := []struct {
		name         string
		field        FactField
		expectedType string
		wantErr      bool
	}{
		{
			name: "string valide",
			field: FactField{
				Name:  "strField",
				Value: FactValue{Type: ValueTypeString, Value: "test"},
			},
			expectedType: ValueTypeString,
			wantErr:      false,
		},
		{
			name: "number valide",
			field: FactField{
				Name:  "numField",
				Value: FactValue{Type: ValueTypeNumber, Value: 42.0},
			},
			expectedType: ValueTypeNumber,
			wantErr:      false,
		},
		{
			name: "boolean valide",
			field: FactField{
				Name:  "boolField",
				Value: FactValue{Type: ValueTypeBoolean, Value: true},
			},
			expectedType: ValueTypeBool,
			wantErr:      false,
		},
		{
			name: "type incompatible",
			field: FactField{
				Name:  "numField",
				Value: FactValue{Type: ValueTypeString, Value: "not a number"},
			},
			expectedType: ValueTypeNumber,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateFieldValue(tt.field, tt.expectedType)

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
