// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"

	"github.com/treivax/tsd/constraint"
)

func TestFieldResolver_ResolveFieldValue(t *testing.T) {
	t.Log("🧪 TEST FIELD RESOLVER - RÉSOLUTION VALEURS")
	t.Log("============================================")

	types := []constraint.TypeDefinition{
		{
			Name: "User",
			Fields: []constraint.Field{
				{Name: "name", Type: "string", IsPrimaryKey: true},
				{Name: "age", Type: "number"},
			},
		},
		{
			Name: "Login",
			Fields: []constraint.Field{
				{Name: "user", Type: "User"},
				{Name: "email", Type: "string", IsPrimaryKey: true},
				{Name: "password", Type: "string"},
			},
		},
	}

	resolver := NewFieldResolver(types)

	tests := []struct {
		name          string
		fact          *Fact
		fieldName     string
		expectedValue interface{}
		expectedType  string
		wantErr       bool
	}{
		{
			name: "champ primitif string",
			fact: &Fact{
				ID:   "Login~alice@ex.com",
				Type: "Login",
				Fields: map[string]interface{}{
					"user":     "User~Alice",
					"email":    "alice@ex.com",
					"password": "secret",
				},
			},
			fieldName:     "email",
			expectedValue: "alice@ex.com",
			expectedType:  FieldTypePrimitive,
			wantErr:       false,
		},
		{
			name: "champ de type fait",
			fact: &Fact{
				ID:   "Login~alice@ex.com",
				Type: "Login",
				Fields: map[string]interface{}{
					"user":     "User~Alice",
					"email":    "alice@ex.com",
					"password": "secret",
				},
			},
			fieldName:     "user",
			expectedValue: "User~Alice",
			expectedType:  FieldTypeFact,
			wantErr:       false,
		},
		{
			name: "champ _id_ interdit",
			fact: &Fact{
				ID:   "User~Alice",
				Type: "User",
				Fields: map[string]interface{}{
					"name": "Alice",
					"age":  30.0,
				},
			},
			fieldName: constraint.FieldNameInternalID,
			wantErr:   true,
		},
		{
			name: "champ non existant",
			fact: &Fact{
				ID:   "User~Alice",
				Type: "User",
				Fields: map[string]interface{}{
					"name": "Alice",
					"age":  30.0,
				},
			},
			fieldName: "unknown",
			wantErr:   true,
		},
		{
			name: "champ primitif number",
			fact: &Fact{
				ID:   "User~Bob",
				Type: "User",
				Fields: map[string]interface{}{
					"name": "Bob",
					"age":  25.0,
				},
			},
			fieldName:     "age",
			expectedValue: 25.0,
			expectedType:  FieldTypePrimitive,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, fieldType, err := resolver.ResolveFieldValue(tt.fact, tt.fieldName)

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

			if value != tt.expectedValue {
				t.Errorf("❌ Valeur attendue %v, reçu %v", tt.expectedValue, value)
			}

			if fieldType != tt.expectedType {
				t.Errorf("❌ Type attendu '%s', reçu '%s'", tt.expectedType, fieldType)
			}

			t.Logf("✅ Résolution correcte: valeur=%v, type=%s", value, fieldType)
		})
	}
}

func TestFieldResolver_ResolveFactID(t *testing.T) {
	t.Log("🧪 TEST FIELD RESOLVER - RÉSOLUTION ID")
	t.Log("=======================================")

	resolver := NewFieldResolver(nil)

	fact := &Fact{
		ID:   "User~Alice",
		Type: "User",
		Fields: map[string]interface{}{
			"name": "Alice",
			"age":  30.0,
		},
	}

	id := resolver.ResolveFactID(fact)

	if id != "User~Alice" {
		t.Errorf("❌ ID attendu 'User~Alice', reçu '%s'", id)
	} else {
		t.Logf("✅ ID résolu correctement: %s", id)
	}
}

func TestFieldResolver_GetFieldType(t *testing.T) {
	t.Log("🧪 TEST FIELD RESOLVER - DÉTECTION TYPE")
	t.Log("========================================")

	types := []constraint.TypeDefinition{
		{
			Name: "User",
			Fields: []constraint.Field{
				{Name: "name", Type: "string"},
			},
		},
	}

	resolver := NewFieldResolver(types)

	tests := []struct {
		name         string
		typeName     string
		expectedType string
	}{
		{"type primitif string", "string", FieldTypePrimitive},
		{"type primitif number", "number", FieldTypePrimitive},
		{"type primitif bool", "bool", FieldTypePrimitive},
		{"type primitif boolean", "boolean", FieldTypePrimitive},
		{"type fait User", "User", FieldTypeFact},
		{"type inconnu", "UnknownType", FieldTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.getFieldType(tt.typeName)

			if result != tt.expectedType {
				t.Errorf("❌ Type attendu '%s', reçu '%s'", tt.expectedType, result)
			} else {
				t.Logf("✅ Type détecté correctement: %s -> %s", tt.typeName, result)
			}
		})
	}
}
