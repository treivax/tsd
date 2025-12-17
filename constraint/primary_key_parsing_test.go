// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"testing"
)

func TestParsePrimaryKeyFields(t *testing.T) {
	t.Log("🧪 TEST PARSING PRIMARY KEY FIELDS")
	t.Log("====================================")

	tests := []struct {
		name           string
		input          string
		wantTypeName   string
		wantFieldCount int
		wantPKFields   []string // Noms des champs marqués comme clé primaire
	}{
		{
			name:           "clé primaire simple",
			input:          "type User(#login: string, name: string, age: number)",
			wantTypeName:   "User",
			wantFieldCount: 3,
			wantPKFields:   []string{"login"},
		},
		{
			name:           "clé primaire composite",
			input:          "type Person(#firstName: string, #lastName: string, age: number)",
			wantTypeName:   "Person",
			wantFieldCount: 3,
			wantPKFields:   []string{"firstName", "lastName"},
		},
		{
			name:           "sans clé primaire",
			input:          "type Document(title: string, content: string)",
			wantTypeName:   "Document",
			wantFieldCount: 2,
			wantPKFields:   []string{},
		},
		{
			name:           "tous les champs sont clé primaire",
			input:          "type Coordinate(#x: number, #y: number, #z: number)",
			wantTypeName:   "Coordinate",
			wantFieldCount: 3,
			wantPKFields:   []string{"x", "y", "z"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Test: %s", tt.name)
			t.Logf("Input: %s", tt.input)

			// Parser le contenu
			result, err := Parse("test.tsd", []byte(tt.input))
			if err != nil {
				t.Fatalf("❌ Erreur de parsing: %v", err)
			}

			// Convertir en Program
			program, err := ConvertResultToProgram(result)
			if err != nil {
				t.Fatalf("❌ Erreur de conversion: %v", err)
			}

			// Vérifier qu'on a bien un type
			if len(program.Types) != 1 {
				t.Fatalf("❌ Attendu 1 type, reçu %d", len(program.Types))
			}

			typeDef := program.Types[0]

			// Vérifier le nom du type
			if typeDef.Name != tt.wantTypeName {
				t.Errorf("❌ Nom du type: attendu '%s', reçu '%s'", tt.wantTypeName, typeDef.Name)
			}

			// Vérifier le nombre de champs
			if len(typeDef.Fields) != tt.wantFieldCount {
				t.Errorf("❌ Nombre de champs: attendu %d, reçu %d", tt.wantFieldCount, len(typeDef.Fields))
			}

			// Vérifier les champs marqués comme clé primaire
			foundPKFields := []string{}
			for _, field := range typeDef.Fields {
				if field.IsPrimaryKey {
					foundPKFields = append(foundPKFields, field.Name)
				}
			}

			if len(foundPKFields) != len(tt.wantPKFields) {
				t.Errorf("❌ Nombre de champs PK: attendu %d, reçu %d", len(tt.wantPKFields), len(foundPKFields))
			}

			// Vérifier chaque champ PK attendu
			for _, wantPK := range tt.wantPKFields {
				found := false
				for _, gotPK := range foundPKFields {
					if gotPK == wantPK {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("❌ Champ PK '%s' attendu mais non trouvé", wantPK)
				}
			}

			t.Log("✅ Test réussi")
		})
	}
}

func TestParsePrimaryKeyInvalidSyntax(t *testing.T) {
	t.Log("🧪 TEST PARSING PRIMARY KEY - CAS INVALIDES")
	t.Log("=============================================")

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "# isolé sans nom de champ",
			input: "type Bad(#: string)",
		},
		{
			name:  "# après le nom",
			input: "type Bad(name#: string)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Test: %s", tt.name)

			_, err := Parse("test.tsd", []byte(tt.input))
			if err == nil {
				t.Errorf("⚠️  Attendu une erreur de parsing, reçu nil")
			} else {
				t.Logf("✅ Erreur attendue: %v", err)
			}
		})
	}
}
