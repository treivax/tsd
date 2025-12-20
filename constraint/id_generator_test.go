// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

func TestGenerateFactIDWithPrimaryKey(t *testing.T) {
	t.Log("🧪 TEST GENERATE FACT ID WITH PRIMARY KEY")
	t.Log("==========================================")

	tests := []struct {
		name    string
		fact    Fact
		typeDef TypeDefinition
		wantID  string
		wantErr bool
	}{
		{
			name: "clé primaire simple",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: "login", Value: FactValue{Type: "string", Value: "alice"}},
					{Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
				},
			},
			typeDef: TypeDefinition{
				Name: "User",
				Fields: []Field{
					{Name: "login", Type: "string", IsPrimaryKey: true},
					{Name: "name", Type: "string", IsPrimaryKey: false},
				},
			},
			wantID:  "User~alice",
			wantErr: false,
		},
		{
			name: "clé primaire composite",
			fact: Fact{
				TypeName: "Person",
				Fields: []FactField{
					{Name: "firstName", Value: FactValue{Type: "string", Value: "Jean-Claude"}},
					{Name: "lastName", Value: FactValue{Type: "string", Value: "Pignon"}},
					{Name: "age", Value: FactValue{Type: "number", Value: float64(27)}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Person",
				Fields: []Field{
					{Name: "firstName", Type: "string", IsPrimaryKey: true},
					{Name: "lastName", Type: "string", IsPrimaryKey: true},
					{Name: "age", Type: "number", IsPrimaryKey: false},
				},
			},
			wantID:  "Person~Jean-Claude_Pignon",
			wantErr: false,
		},
		{
			name: "clé primaire avec number",
			fact: Fact{
				TypeName: "Product",
				Fields: []FactField{
					{Name: "code", Value: FactValue{Type: "number", Value: float64(12345)}},
					{Name: "name", Value: FactValue{Type: "string", Value: "Widget"}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Product",
				Fields: []Field{
					{Name: "code", Type: "number", IsPrimaryKey: true},
					{Name: "name", Type: "string", IsPrimaryKey: false},
				},
			},
			wantID:  "Product~12345",
			wantErr: false,
		},
		{
			name: "clé primaire avec bool",
			fact: Fact{
				TypeName: "Flag",
				Fields: []FactField{
					{Name: "active", Value: FactValue{Type: "bool", Value: true}},
					{Name: "label", Value: FactValue{Type: "string", Value: "Test"}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Flag",
				Fields: []Field{
					{Name: "active", Type: "bool", IsPrimaryKey: true},
					{Name: "label", Type: "string", IsPrimaryKey: false},
				},
			},
			wantID:  "Flag~true",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := GenerateFactID(tt.fact, tt.typeDef, nil)

			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("❌ Erreur inattendue: %v", err)
				} else if id != tt.wantID {
					t.Errorf("❌ ID: attendu '%s', reçu '%s'", tt.wantID, id)
				} else {
					t.Logf("✅ ID généré: %s", id)
				}
			}
		})
	}
}

func TestGenerateFactIDWithHash(t *testing.T) {
	t.Log("🧪 TEST GENERATE FACT ID WITH HASH")
	t.Log("===================================")

	typeDef := TypeDefinition{
		Name: "Document",
		Fields: []Field{
			{Name: "title", Type: "string", IsPrimaryKey: false},
			{Name: "content", Type: "string", IsPrimaryKey: false},
		},
	}

	fact := Fact{
		TypeName: "Document",
		Fields: []FactField{
			{Name: "title", Value: FactValue{Type: "string", Value: "Doc1"}},
			{Name: "content", Value: FactValue{Type: "string", Value: "Content"}},
		},
	}

	id, err := GenerateFactID(fact, typeDef, nil)
	if err != nil {
		t.Fatalf("❌ Erreur inattendue: %v", err)
	}

	// Vérifier le format: Document~<16 caractères hex>
	if !strings.HasPrefix(id, "Document~") {
		t.Errorf("❌ ID devrait commencer par 'Document~', reçu '%s'", id)
	}

	hashPart := strings.TrimPrefix(id, "Document~")
	if len(hashPart) != IDHashLength {
		t.Errorf("❌ Hash devrait avoir %d caractères, reçu %d", IDHashLength, len(hashPart))
	}

	if !isHexString(hashPart) {
		t.Errorf("❌ Hash devrait être hexadécimal, reçu '%s'", hashPart)
	}

	// Vérifier la reproductibilité (même fait = même hash)
	id2, err := GenerateFactID(fact, typeDef, nil)
	if err != nil {
		t.Fatalf("❌ Erreur inattendue: %v", err)
	}

	if id != id2 {
		t.Errorf("❌ Hash non reproductible: '%s' != '%s'", id, id2)
	}

	t.Logf("✅ ID généré avec hash: %s", id)
}

func TestEscapeIDValue(t *testing.T) {
	t.Log("🧪 TEST ESCAPE ID VALUE")
	t.Log("========================")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "pas de caractères spéciaux",
			input:    "alice",
			expected: "alice",
		},
		{
			name:     "avec tilde",
			input:    "user~123",
			expected: "user%7E123",
		},
		{
			name:     "avec underscore",
			input:    "first_last",
			expected: "first%5Flast",
		},
		{
			name:     "avec les deux",
			input:    "user~name_123",
			expected: "user%7Ename%5F123",
		},
		{
			name:     "avec pourcent",
			input:    "discount%20",
			expected: "discount%2520",
		},
		{
			name:     "pourcent et tilde",
			input:    "%~test",
			expected: "%25%7Etest",
		},
		{
			name:     "séquence complexe",
			input:    "a%b~c_d",
			expected: "a%25b%7Ec%5Fd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeIDValue(tt.input)
			if result != tt.expected {
				t.Errorf("❌ Attendu '%s', reçu '%s'", tt.expected, result)
			} else {
				t.Logf("✅ '%s' → '%s'", tt.input, result)
			}

			unescaped := unescapeIDValue(result)
			if unescaped != tt.input {
				t.Errorf("❌ Unescape: attendu '%s', reçu '%s'", tt.input, unescaped)
			} else {
				t.Logf("✅ Unescape: '%s' → '%s'", result, unescaped)
			}
		})
	}
}

func TestParseFactID(t *testing.T) {
	t.Log("🧪 TEST PARSE FACT ID")
	t.Log("======================")

	tests := []struct {
		name         string
		id           string
		wantTypeName string
		wantPKValues []string
		wantIsHashID bool
		wantErr      bool
	}{
		{
			name:         "clé primaire simple",
			id:           "User~alice",
			wantTypeName: "User",
			wantPKValues: []string{"alice"},
			wantIsHashID: false,
			wantErr:      false,
		},
		{
			name:         "clé primaire composite",
			id:           "Person~Jean-Claude_Pignon",
			wantTypeName: "Person",
			wantPKValues: []string{"Jean-Claude", "Pignon"},
			wantIsHashID: false,
			wantErr:      false,
		},
		{
			name:         "hash ID",
			id:           "Document~a3f5b9c2e1d4f8a7",
			wantTypeName: "Document",
			wantPKValues: []string{"a3f5b9c2e1d4f8a7"},
			wantIsHashID: true,
			wantErr:      false,
		},
		{
			name:    "format invalide - pas de tilde",
			id:      "InvalidIDFormat",
			wantErr: true,
		},
		{
			name:    "format invalide - type vide",
			id:      "~value",
			wantErr: true,
		},
		{
			name:    "format invalide - valeur vide",
			id:      "Type~",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeName, pkValues, isHashID, err := ParseFactID(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("❌ Erreur inattendue: %v", err)
				return
			}

			if typeName != tt.wantTypeName {
				t.Errorf("❌ Type: attendu '%s', reçu '%s'", tt.wantTypeName, typeName)
			}

			if isHashID != tt.wantIsHashID {
				t.Errorf("❌ IsHashID: attendu %v, reçu %v", tt.wantIsHashID, isHashID)
			}

			if len(pkValues) != len(tt.wantPKValues) {
				t.Errorf("❌ Nombre de valeurs: attendu %d, reçu %d", len(tt.wantPKValues), len(pkValues))
			} else {
				for i, want := range tt.wantPKValues {
					if pkValues[i] != want {
						t.Errorf("❌ Valeur[%d]: attendu '%s', reçu '%s'", i, want, pkValues[i])
					}
				}
			}

			t.Log("✅ Test réussi")
		})
	}
}

func TestValueToString(t *testing.T) {
	t.Log("🧪 TEST VALUE TO STRING")
	t.Log("========================")

	tests := []struct {
		name     string
		value    interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "string",
			value:    "test",
			expected: "test",
			wantErr:  false,
		},
		{
			name:     "int",
			value:    42,
			expected: "42",
			wantErr:  false,
		},
		{
			name:     "int64",
			value:    int64(123),
			expected: "123",
			wantErr:  false,
		},
		{
			name:     "float64",
			value:    float64(3.14),
			expected: "3.14",
			wantErr:  false,
		},
		{
			name:     "bool true",
			value:    true,
			expected: "true",
			wantErr:  false,
		},
		{
			name:     "bool false",
			value:    false,
			expected: "false",
			wantErr:  false,
		},
		{
			name:    "nil",
			value:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := valueToString(tt.value)

			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("❌ Erreur inattendue: %v", err)
				} else if result != tt.expected {
					t.Errorf("❌ Attendu '%s', reçu '%s'", tt.expected, result)
				} else {
					t.Logf("✅ %v → '%s'", tt.value, result)
				}
			}
		})
	}
}
