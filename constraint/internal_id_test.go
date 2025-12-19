// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

func TestInternalID_ReservedName(t *testing.T) {
	t.Log("🧪 TEST INTERNAL ID - NOM RÉSERVÉ")
	t.Log("===================================")

	tests := []struct {
		name    string
		testFn  func() error
		wantErr bool
	}{
		{
			name: "interdire _id_ dans définition de type",
			testFn: func() error {
				program := Program{
					Types: []TypeDefinition{
						{
							Name: "User",
							Fields: []Field{
								{Name: FieldNameInternalID, Type: ValueTypeString},
								{Name: "name", Type: ValueTypeString},
							},
						},
					},
				}
				return ValidateTypes(program)
			},
			wantErr: true,
		},
		{
			name: "interdire _id_ dans définition de fait",
			testFn: func() error {
				typeDef := TypeDefinition{
					Name: "User",
					Fields: []Field{
						{Name: "name", Type: ValueTypeString, IsPrimaryKey: true},
					},
				}
				fact := Fact{
					TypeName: "User",
					Fields: []FactField{
						{Name: FieldNameInternalID, Value: FactValue{Type: ValueTypeString, Value: "test"}},
						{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
					},
				}
				return ValidateFactPrimaryKey(fact, typeDef)
			},
			wantErr: true,
		},
		{
			name: "type valide sans _id_",
			testFn: func() error {
				program := Program{
					Types: []TypeDefinition{
						{
							Name: "User",
							Fields: []Field{
								{Name: "name", Type: ValueTypeString, IsPrimaryKey: true},
							},
						},
					},
				}
				return ValidateTypes(program)
			},
			wantErr: false,
		},
		{
			name: "fait valide sans _id_",
			testFn: func() error {
				typeDef := TypeDefinition{
					Name: "User",
					Fields: []Field{
						{Name: "name", Type: ValueTypeString, IsPrimaryKey: true},
					},
				}
				fact := Fact{
					TypeName: "User",
					Fields: []FactField{
						{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
					},
				}
				return ValidateFactPrimaryKey(fact, typeDef)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testFn()

			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
					if !strings.Contains(err.Error(), "réservé") {
						t.Errorf("⚠️  Message d'erreur ne contient pas 'réservé': %v", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("❌ Erreur inattendue: %v", err)
				} else {
					t.Logf("✅ Pas d'erreur comme attendu")
				}
			}
		})
	}
}

func TestInternalID_AlwaysGenerated(t *testing.T) {
	t.Log("🧪 TEST INTERNAL ID - TOUJOURS GÉNÉRÉ")
	t.Log("======================================")

	fact := Fact{
		TypeName: "User",
		Fields: []FactField{
			{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
		},
	}

	typeDef := TypeDefinition{
		Name: "User",
		Fields: []Field{
			{Name: "name", Type: ValueTypeString, IsPrimaryKey: true},
		},
	}

	ctx := NewFactContext([]TypeDefinition{typeDef})
	reteFact := createReteFact(fact, typeDef, ctx)
	id, err := ensureFactID(reteFact, fact, typeDef, ctx)

	if err != nil {
		t.Fatalf("❌ Erreur inattendue: %v", err)
	}

	if id == "" {
		t.Error("❌ ID généré est vide")
	}

	expectedID := "User~Alice"
	if id != expectedID {
		t.Errorf("❌ ID attendu '%s', reçu '%s'", expectedID, id)
	}

	t.Logf("✅ ID généré correctement: %s", id)
}

func TestInternalID_NeverManual(t *testing.T) {
	t.Log("🧪 TEST INTERNAL ID - JAMAIS MANUEL")
	t.Log("====================================")

	fact := Fact{
		TypeName: "User",
		Fields: []FactField{
			{Name: "name", Value: FactValue{Type: ValueTypeString, Value: "Alice"}},
		},
	}

	typeDef := TypeDefinition{
		Name: "User",
		Fields: []Field{
			{Name: "name", Type: ValueTypeString, IsPrimaryKey: true},
		},
	}

	ctx := NewFactContext([]TypeDefinition{typeDef})
	reteFact := createReteFact(fact, typeDef, ctx)

	// Simuler une tentative de définition manuelle
	reteFact[FieldNameInternalID] = "manual_id"

	_, err := ensureFactID(reteFact, fact, typeDef, ctx)

	if err == nil {
		t.Error("❌ Attendu une erreur pour ID manuel")
	} else {
		t.Logf("✅ Erreur attendue pour ID manuel: %v", err)
		if !strings.Contains(err.Error(), "ne peut pas être défini manuellement") {
			t.Errorf("⚠️  Message d'erreur inattendu: %v", err)
		}
	}
}

func TestInternalID_GenerationFromHash(t *testing.T) {
	t.Log("🧪 TEST INTERNAL ID - GÉNÉRATION PAR HASH")
	t.Log("==========================================")

	// Type sans clé primaire - devrait générer un hash
	fact := Fact{
		TypeName: "Event",
		Fields: []FactField{
			{Name: "timestamp", Value: FactValue{Type: ValueTypeNumber, Value: float64(1234567890)}},
			{Name: "message", Value: FactValue{Type: ValueTypeString, Value: "Test event"}},
		},
	}

	typeDef := TypeDefinition{
		Name: "Event",
		Fields: []Field{
			{Name: "timestamp", Type: ValueTypeNumber},
			{Name: "message", Type: ValueTypeString},
		},
		// Pas de clé primaire définie
	}

	ctx := NewFactContext([]TypeDefinition{typeDef})
	reteFact := createReteFact(fact, typeDef, ctx)
	id, err := ensureFactID(reteFact, fact, typeDef, ctx)

	if err != nil {
		t.Fatalf("❌ Erreur inattendue: %v", err)
	}

	if id == "" {
		t.Fatal("❌ ID généré est vide")
	}

	// L'ID devrait commencer par le nom du type
	if !strings.HasPrefix(id, "Event~") {
		t.Errorf("❌ ID devrait commencer par 'Event~', reçu: %s", id)
	}

	// La partie après ~ devrait être un hash (16 caractères hexadécimaux)
	parts := strings.SplitN(id, "~", 2)
	if len(parts) != 2 {
		t.Fatalf("❌ Format d'ID invalide: %s", id)
	}

	hashPart := parts[1]
	if len(hashPart) != 16 {
		t.Errorf("❌ Hash devrait faire 16 caractères, reçu %d: %s", len(hashPart), hashPart)
	}

	// Vérifier que c'est hexadécimal
	for _, c := range hashPart {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("❌ Hash devrait être hexadécimal, caractère invalide: %c", c)
			break
		}
	}

	t.Logf("✅ ID généré par hash: %s", id)
}

func TestInternalID_AccessForbidden(t *testing.T) {
	t.Log("🧪 TEST INTERNAL ID - ACCÈS INTERDIT")
	t.Log("=====================================")

	program := Program{
		Types: []TypeDefinition{
			{
				Name: "User",
				Fields: []Field{
					{Name: "name", Type: ValueTypeString, IsPrimaryKey: true},
				},
			},
		},
		Expressions: []Expression{
			{
				Set: Set{
					Variables: []TypedVariable{
						{Name: "u", DataType: "User"},
					},
				},
			},
		},
	}

	// Tenter d'obtenir le type du champ _id_
	_, err := GetFieldType(program, "u", FieldNameInternalID, 0)

	if err == nil {
		t.Error("❌ Attendu une erreur pour accès à _id_")
	} else {
		t.Logf("✅ Erreur attendue pour accès à _id_: %v", err)
		if !strings.Contains(err.Error(), "interne") {
			t.Errorf("⚠️  Message d'erreur ne contient pas 'interne': %v", err)
		}
	}
}

func TestInternalID_ConversionToRete(t *testing.T) {
	t.Log("🧪 TEST INTERNAL ID - CONVERSION RETE")
	t.Log("======================================")

	program := Program{
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
					{Name: "age", Value: FactValue{Type: ValueTypeNumber, Value: float64(30)}},
				},
			},
		},
	}

	reteFacts, err := ConvertFactsToReteFormat(program)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}

	if len(reteFacts) != 1 {
		t.Fatalf("❌ Attendu 1 fait RETE, reçu %d", len(reteFacts))
	}

	reteFact := reteFacts[0]

	// Vérifier que _id_ est présent
	id, exists := reteFact[FieldNameInternalID]
	if !exists {
		t.Fatal("❌ Champ _id_ manquant dans le fait RETE")
	}

	expectedID := "User~Alice"
	if id != expectedID {
		t.Errorf("❌ ID attendu '%s', reçu '%v'", expectedID, id)
	}

	// Vérifier que les autres champs sont présents
	if name, exists := reteFact["name"]; !exists || name != "Alice" {
		t.Errorf("❌ Champ 'name' incorrect: %v", name)
	}

	if age, exists := reteFact["age"]; !exists || age != float64(30) {
		t.Errorf("❌ Champ 'age' incorrect: %v", age)
	}

	t.Logf("✅ Conversion RETE réussie avec _id_ = %s", id)
}
