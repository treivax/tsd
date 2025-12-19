// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

func TestGenerateFactID_EdgeCases(t *testing.T) {
	t.Log("🧪 TEST GENERATE FACT ID - CAS LIMITES")
	t.Log("=======================================")

	tests := []struct {
		name    string
		fact    Fact
		typeDef TypeDefinition
		wantID  string
		wantErr bool
		checkFn func(*testing.T, string)
	}{
		{
			name: "PK avec chaîne vide",
			fact: Fact{
				TypeName: "Person",
				Fields: []FactField{
					{Name: "nom", Value: FactValue{Type: "string", Value: ""}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Person",
				Fields: []Field{
					{Name: "nom", Type: "string", IsPrimaryKey: true},
				},
			},
			wantID:  "Person~",
			wantErr: false,
		},
		{
			name: "PK avec zéro",
			fact: Fact{
				TypeName: "Item",
				Fields: []FactField{
					{Name: "id", Value: FactValue{Type: "number", Value: float64(0)}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Item",
				Fields: []Field{
					{Name: "id", Type: "number", IsPrimaryKey: true},
				},
			},
			wantID:  "Item~0",
			wantErr: false,
		},
		{
			name: "PK avec false",
			fact: Fact{
				TypeName: "Flag",
				Fields: []FactField{
					{Name: "enabled", Value: FactValue{Type: "bool", Value: false}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Flag",
				Fields: []Field{
					{Name: "enabled", Type: "bool", IsPrimaryKey: true},
				},
			},
			wantID:  "Flag~false",
			wantErr: false,
		},
		{
			name: "Float avec précision élevée",
			fact: Fact{
				TypeName: "Measurement",
				Fields: []FactField{
					{Name: "value", Value: FactValue{Type: "number", Value: 3.141592653589793}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Measurement",
				Fields: []Field{
					{Name: "value", Type: "number", IsPrimaryKey: true},
				},
			},
			wantID:  "Measurement~3.141592653589793",
			wantErr: false,
		},
		{
			name: "Float très petit",
			fact: Fact{
				TypeName: "Precision",
				Fields: []FactField{
					{Name: "val", Value: FactValue{Type: "number", Value: 0.000000001}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Precision",
				Fields: []Field{
					{Name: "val", Type: "number", IsPrimaryKey: true},
				},
			},
			checkFn: func(t *testing.T, id string) {
				if !strings.HasPrefix(id, "Precision~0.00000000") {
					t.Errorf("❌ ID devrait commencer par 'Precision~0.00000000', reçu '%s'", id)
				} else {
					t.Logf("✅ Float très petit formaté correctement: %s", id)
				}
			},
			wantErr: false,
		},
		{
			name: "Float très grand",
			fact: Fact{
				TypeName: "BigNumber",
				Fields: []FactField{
					{Name: "amount", Value: FactValue{Type: "number", Value: 999999999999.99}},
				},
			},
			typeDef: TypeDefinition{
				Name: "BigNumber",
				Fields: []Field{
					{Name: "amount", Type: "number", IsPrimaryKey: true},
				},
			},
			wantID:  "BigNumber~999999999999.99",
			wantErr: false,
		},
		{
			name: "PK composite avec >3 champs",
			fact: Fact{
				TypeName: "ComplexKey",
				Fields: []FactField{
					{Name: "field1", Value: FactValue{Type: "string", Value: "a"}},
					{Name: "field2", Value: FactValue{Type: "string", Value: "b"}},
					{Name: "field3", Value: FactValue{Type: "string", Value: "c"}},
					{Name: "field4", Value: FactValue{Type: "number", Value: 1.0}},
					{Name: "field5", Value: FactValue{Type: "bool", Value: true}},
				},
			},
			typeDef: TypeDefinition{
				Name: "ComplexKey",
				Fields: []Field{
					{Name: "field1", Type: "string", IsPrimaryKey: true},
					{Name: "field2", Type: "string", IsPrimaryKey: true},
					{Name: "field3", Type: "string", IsPrimaryKey: true},
					{Name: "field4", Type: "number", IsPrimaryKey: true},
					{Name: "field5", Type: "bool", IsPrimaryKey: true},
				},
			},
			wantID:  "ComplexKey~a_b_c_1_true",
			wantErr: false,
		},
		{
			name: "String avec caractères Unicode",
			fact: Fact{
				TypeName: "Unicode",
				Fields: []FactField{
					{Name: "text", Value: FactValue{Type: "string", Value: "Héllo 世界 🌍"}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Unicode",
				Fields: []Field{
					{Name: "text", Type: "string", IsPrimaryKey: true},
				},
			},
			checkFn: func(t *testing.T, id string) {
				if !strings.HasPrefix(id, "Unicode~") {
					t.Errorf("❌ ID devrait commencer par 'Unicode~', reçu '%s'", id)
				} else if !strings.Contains(id, "H") {
					t.Errorf("❌ ID devrait contenir le texte Unicode, reçu '%s'", id)
				} else {
					t.Logf("✅ Unicode géré correctement: %s", id)
				}
			},
			wantErr: false,
		},
		{
			name: "String avec tous types de caractères spéciaux",
			fact: Fact{
				TypeName: "Special",
				Fields: []FactField{
					{Name: "path", Value: FactValue{Type: "string", Value: "/path%with~special_chars"}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Special",
				Fields: []Field{
					{Name: "path", Type: "string", IsPrimaryKey: true},
				},
			},
			// Note: Seuls %, ~, et _ sont échappés (pas le /)
			wantID:  "Special~/path%25with%7Especial%5Fchars",
			wantErr: false,
		},
		{
			name: "Hash avec tous types de champs",
			fact: Fact{
				TypeName: "Mixed",
				Fields: []FactField{
					{Name: "str", Value: FactValue{Type: "string", Value: "text"}},
					{Name: "num", Value: FactValue{Type: "number", Value: 42.5}},
					{Name: "flag", Value: FactValue{Type: "bool", Value: true}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Mixed",
				Fields: []Field{
					{Name: "str", Type: "string", IsPrimaryKey: false},
					{Name: "num", Type: "number", IsPrimaryKey: false},
					{Name: "flag", Type: "bool", IsPrimaryKey: false},
				},
			},
			checkFn: func(t *testing.T, id string) {
				if !strings.HasPrefix(id, "Mixed~") {
					t.Errorf("❌ ID devrait commencer par 'Mixed~', reçu '%s'", id)
				}
				hashPart := strings.TrimPrefix(id, "Mixed~")
				if len(hashPart) != IDHashLength {
					t.Errorf("❌ Hash devrait avoir %d caractères, reçu %d", IDHashLength, len(hashPart))
				}
				if !isHexString(hashPart) {
					t.Errorf("❌ Hash devrait être hexadécimal, reçu '%s'", hashPart)
				}
				t.Logf("✅ Hash généré pour types mixtes: %s", id)
			},
			wantErr: false,
		},
		{
			name: "Déterminisme - même fait génère même ID",
			fact: Fact{
				TypeName: "Event",
				Fields: []FactField{
					{Name: "timestamp", Value: FactValue{Type: "number", Value: 1234567890.123}},
					{Name: "message", Value: FactValue{Type: "string", Value: "test message"}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Event",
				Fields: []Field{
					{Name: "timestamp", Type: "number", IsPrimaryKey: false},
					{Name: "message", Type: "string", IsPrimaryKey: false},
				},
			},
			checkFn: func(t *testing.T, id string) {
				// Générer plusieurs fois et vérifier que c'est toujours le même
				for i := 0; i < 5; i++ {
					id2, err := GenerateFactID(Fact{
						TypeName: "Event",
						Fields: []FactField{
							{Name: "timestamp", Value: FactValue{Type: "number", Value: 1234567890.123}},
							{Name: "message", Value: FactValue{Type: "string", Value: "test message"}},
						},
					}, TypeDefinition{
						Name: "Event",
						Fields: []Field{
							{Name: "timestamp", Type: "number", IsPrimaryKey: false},
							{Name: "message", Type: "string", IsPrimaryKey: false},
						},
					}, nil)
					if err != nil {
						t.Errorf("❌ Erreur lors de la génération %d: %v", i+1, err)
					}
					if id2 != id {
						t.Errorf("❌ ID non déterministe: tentative %d génère '%s' au lieu de '%s'", i+1, id2, id)
					}
				}
				t.Logf("✅ Déterminisme confirmé sur 5 générations: %s", id)
			},
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
				return
			}

			if err != nil {
				t.Errorf("❌ Erreur inattendue: %v", err)
				return
			}

			if tt.checkFn != nil {
				tt.checkFn(t, id)
			} else if tt.wantID != "" {
				if id != tt.wantID {
					t.Errorf("❌ ID: attendu '%s', reçu '%s'", tt.wantID, id)
				} else {
					t.Logf("✅ ID généré: %s", id)
				}
			}
		})
	}
}

func TestValueToString_FloatFormats(t *testing.T) {
	t.Log("🧪 TEST VALUE TO STRING - FORMATS FLOAT")
	t.Log("========================================")

	tests := []struct {
		name      string
		value     float64
		wantRegex string // Regex pattern si format exact non prévisible
	}{
		{
			name:  "entier représenté en float",
			value: 42.0,
		},
		{
			name:  "float simple",
			value: 3.14,
		},
		{
			name:  "float précision IEEE 754",
			value: 0.1 + 0.2, // = 0.30000000000000004
		},
		{
			name:  "très grand nombre",
			value: 1e15,
		},
		{
			name:  "très petit nombre",
			value: 1e-10,
		},
		{
			name:  "zéro",
			value: 0.0,
		},
		{
			name:  "nombre négatif",
			value: -123.456,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := valueToString(tt.value)
			if err != nil {
				t.Errorf("❌ Erreur inattendue: %v", err)
			} else {
				// Vérifier que le format ne contient pas d'exposant
				if strings.Contains(result, "e") || strings.Contains(result, "E") {
					t.Errorf("⚠️  Format avec exposant: %s (valeur: %f)", result, tt.value)
				} else {
					t.Logf("✅ %f → '%s' (sans notation scientifique)", tt.value, result)
				}

				// Vérifier la reproductibilité
				result2, _ := valueToString(tt.value)
				if result != result2 {
					t.Errorf("❌ Format non reproductible: '%s' != '%s'", result, result2)
				}
			}
		})
	}
}

func TestConstants_ThreadSafety(t *testing.T) {
	t.Log("🧪 TEST CONSTANTS - THREAD SAFETY")
	t.Log("==================================")

	// Test que l'initialisation paresseuse est thread-safe
	done := make(chan bool)

	// Lancer plusieurs goroutines qui utilisent les constantes simultanément
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				_ = IsValidOperator(OpEq)
				_ = IsValidPrimitiveType(ValueTypeString)
			}
			done <- true
		}()
	}

	// Attendre que toutes les goroutines finissent
	for i := 0; i < 10; i++ {
		<-done
	}

	t.Log("✅ Pas de race condition détectée")
}
