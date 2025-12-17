// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"encoding/json"
	"testing"
)

func TestFieldPrimaryKey(t *testing.T) {
	t.Log("🧪 TEST FIELD PRIMARY KEY")
	t.Log("=========================")

	tests := []struct {
		name     string
		field    Field
		wantIsPK bool
	}{
		{
			name: "champ standard",
			field: Field{
				Name:         "age",
				Type:         "number",
				IsPrimaryKey: false,
			},
			wantIsPK: false,
		},
		{
			name: "champ clé primaire",
			field: Field{
				Name:         "id",
				Type:         "string",
				IsPrimaryKey: true,
			},
			wantIsPK: true,
		},
		{
			name: "valeur par défaut",
			field: Field{
				Name: "name",
				Type: "string",
				// IsPrimaryKey non spécifié = false
			},
			wantIsPK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.field.IsPrimaryKey != tt.wantIsPK {
				t.Errorf("❌ IsPrimaryKey: attendu %v, reçu %v",
					tt.wantIsPK, tt.field.IsPrimaryKey)
			} else {
				t.Log("✅ Test réussi")
			}
		})
	}
}

func TestFieldJSONSerialization(t *testing.T) {
	t.Log("🧪 TEST FIELD JSON SERIALIZATION")
	t.Log("=================================")

	tests := []struct {
		name       string
		field      Field
		wantJSON   string
		shouldOmit bool
	}{
		{
			name: "field with primary key true",
			field: Field{
				Name:         "id",
				Type:         "string",
				IsPrimaryKey: true,
			},
			wantJSON:   `{"name":"id","type":"string","isPrimaryKey":true}`,
			shouldOmit: false,
		},
		{
			name: "field with primary key false should omit",
			field: Field{
				Name:         "age",
				Type:         "number",
				IsPrimaryKey: false,
			},
			wantJSON:   `{"name":"age","type":"number"}`,
			shouldOmit: true,
		},
		{
			name: "field without primary key specified",
			field: Field{
				Name: "name",
				Type: "string",
			},
			wantJSON:   `{"name":"name","type":"string"}`,
			shouldOmit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test serialization
			data, err := json.Marshal(tt.field)
			if err != nil {
				t.Fatalf("❌ Erreur de sérialisation: %v", err)
			}

			if string(data) != tt.wantJSON {
				t.Errorf("❌ JSON: attendu %s, reçu %s", tt.wantJSON, string(data))
			}

			// Test deserialization
			var field Field
			err = json.Unmarshal(data, &field)
			if err != nil {
				t.Fatalf("❌ Erreur de désérialisation: %v", err)
			}

			if field.Name != tt.field.Name {
				t.Errorf("❌ Name: attendu %s, reçu %s", tt.field.Name, field.Name)
			}
			if field.Type != tt.field.Type {
				t.Errorf("❌ Type: attendu %s, reçu %s", tt.field.Type, field.Type)
			}
			if field.IsPrimaryKey != tt.field.IsPrimaryKey {
				t.Errorf("❌ IsPrimaryKey: attendu %v, reçu %v", tt.field.IsPrimaryKey, field.IsPrimaryKey)
			}

			t.Log("✅ Test réussi")
		})
	}
}

func TestTypeDefinitionPrimaryKeyMethods(t *testing.T) {
	t.Log("🧪 TEST TYPE DEFINITION PRIMARY KEY METHODS")
	t.Log("============================================")

	tests := []struct {
		name             string
		typeDef          TypeDefinition
		wantHasPK        bool
		wantPKFieldCount int
		wantPKFieldNames []string
	}{
		{
			name: "sans clé primaire",
			typeDef: TypeDefinition{
				Name: "Document",
				Fields: []Field{
					{Name: "title", Type: "string", IsPrimaryKey: false},
					{Name: "content", Type: "string", IsPrimaryKey: false},
				},
			},
			wantHasPK:        false,
			wantPKFieldCount: 0,
			wantPKFieldNames: []string{},
		},
		{
			name: "clé primaire simple",
			typeDef: TypeDefinition{
				Name: "User",
				Fields: []Field{
					{Name: "login", Type: "string", IsPrimaryKey: true},
					{Name: "name", Type: "string", IsPrimaryKey: false},
					{Name: "age", Type: "number", IsPrimaryKey: false},
				},
			},
			wantHasPK:        true,
			wantPKFieldCount: 1,
			wantPKFieldNames: []string{"login"},
		},
		{
			name: "clé primaire composite",
			typeDef: TypeDefinition{
				Name: "Person",
				Fields: []Field{
					{Name: "firstName", Type: "string", IsPrimaryKey: true},
					{Name: "lastName", Type: "string", IsPrimaryKey: true},
					{Name: "age", Type: "number", IsPrimaryKey: false},
				},
			},
			wantHasPK:        true,
			wantPKFieldCount: 2,
			wantPKFieldNames: []string{"firstName", "lastName"},
		},
		{
			name: "tous les champs sont PK",
			typeDef: TypeDefinition{
				Name: "Coordinate",
				Fields: []Field{
					{Name: "x", Type: "number", IsPrimaryKey: true},
					{Name: "y", Type: "number", IsPrimaryKey: true},
				},
			},
			wantHasPK:        true,
			wantPKFieldCount: 2,
			wantPKFieldNames: []string{"x", "y"},
		},
		{
			name: "type vide",
			typeDef: TypeDefinition{
				Name:   "Empty",
				Fields: []Field{},
			},
			wantHasPK:        false,
			wantPKFieldCount: 0,
			wantPKFieldNames: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test HasPrimaryKey
			if got := tt.typeDef.HasPrimaryKey(); got != tt.wantHasPK {
				t.Errorf("❌ HasPrimaryKey(): attendu %v, reçu %v", tt.wantHasPK, got)
			}

			// Test GetPrimaryKeyFields
			pkFields := tt.typeDef.GetPrimaryKeyFields()
			if len(pkFields) != tt.wantPKFieldCount {
				t.Errorf("❌ GetPrimaryKeyFields() count: attendu %d, reçu %d",
					tt.wantPKFieldCount, len(pkFields))
			}

			// Verify that all returned fields have IsPrimaryKey = true
			for i, field := range pkFields {
				if !field.IsPrimaryKey {
					t.Errorf("❌ GetPrimaryKeyFields()[%d]: le champ '%s' n'a pas IsPrimaryKey=true",
						i, field.Name)
				}
			}

			// Test GetPrimaryKeyFieldNames
			pkNames := tt.typeDef.GetPrimaryKeyFieldNames()
			if len(pkNames) != len(tt.wantPKFieldNames) {
				t.Errorf("❌ GetPrimaryKeyFieldNames() count: attendu %d, reçu %d",
					len(tt.wantPKFieldNames), len(pkNames))
			}

			// Vérifier l'ordre et les noms
			for i, wantName := range tt.wantPKFieldNames {
				if i >= len(pkNames) {
					t.Errorf("❌ Manque le champ PK '%s' à l'index %d", wantName, i)
					continue
				}
				if pkNames[i] != wantName {
					t.Errorf("❌ Champ PK[%d]: attendu '%s', reçu '%s'",
						i, wantName, pkNames[i])
				}
			}

			t.Log("✅ Test réussi")
		})
	}
}

func TestTypeDefinitionClone(t *testing.T) {
	t.Log("🧪 TEST TYPE DEFINITION CLONE WITH PRIMARY KEY")
	t.Log("===============================================")

	original := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "id", Type: "string", IsPrimaryKey: true},
			{Name: "name", Type: "string", IsPrimaryKey: false},
			{Name: "email", Type: "string", IsPrimaryKey: false},
		},
	}

	cloned := original.Clone()

	// Vérifier que le clone a les mêmes valeurs
	if cloned.Type != original.Type {
		t.Errorf("❌ Type: attendu '%s', reçu '%s'", original.Type, cloned.Type)
	}
	if cloned.Name != original.Name {
		t.Errorf("❌ Name: attendu '%s', reçu '%s'", original.Name, cloned.Name)
	}

	if len(cloned.Fields) != len(original.Fields) {
		t.Fatalf("❌ Nombre de champs: attendu %d, reçu %d",
			len(original.Fields), len(cloned.Fields))
	}

	for i := range original.Fields {
		if cloned.Fields[i].Name != original.Fields[i].Name {
			t.Errorf("❌ Field[%d].Name: attendu '%s', reçu '%s'",
				i, original.Fields[i].Name, cloned.Fields[i].Name)
		}
		if cloned.Fields[i].Type != original.Fields[i].Type {
			t.Errorf("❌ Field[%d].Type: attendu '%s', reçu '%s'",
				i, original.Fields[i].Type, cloned.Fields[i].Type)
		}
		if cloned.Fields[i].IsPrimaryKey != original.Fields[i].IsPrimaryKey {
			t.Errorf("❌ Field[%d].IsPrimaryKey: attendu %v, reçu %v",
				i, original.Fields[i].IsPrimaryKey, cloned.Fields[i].IsPrimaryKey)
		}
	}

	// Modifier le clone et vérifier que l'original n'est pas affecté
	cloned.Fields[0].IsPrimaryKey = false
	cloned.Fields[0].Name = "modified_id"
	cloned.Name = "ModifiedPerson"

	// Vérifier que l'original n'est pas modifié
	if !original.Fields[0].IsPrimaryKey {
		t.Error("❌ Modification du clone a affecté l'original IsPrimaryKey (copie non profonde)")
	}
	if original.Fields[0].Name != "id" {
		t.Error("❌ Modification du clone a affecté l'original Name (copie non profonde)")
	}
	if original.Name != "Person" {
		t.Error("❌ Modification du clone a affecté l'original TypeName (copie non profonde)")
	}

	t.Log("✅ Clone test réussi")
}

func TestTypeDefinitionEmptyPrimaryKey(t *testing.T) {
	t.Log("🧪 TEST TYPE DEFINITION EMPTY PRIMARY KEY")
	t.Log("==========================================")

	typeDef := TypeDefinition{
		Name:   "Test",
		Fields: []Field{},
	}

	if typeDef.HasPrimaryKey() {
		t.Error("❌ Type vide ne devrait pas avoir de clé primaire")
	}

	pkFields := typeDef.GetPrimaryKeyFields()
	if len(pkFields) != 0 {
		t.Errorf("❌ GetPrimaryKeyFields() devrait retourner une slice vide, reçu %d éléments", len(pkFields))
	}

	pkNames := typeDef.GetPrimaryKeyFieldNames()
	if len(pkNames) != 0 {
		t.Errorf("❌ GetPrimaryKeyFieldNames() devrait retourner une slice vide, reçu %d éléments", len(pkNames))
	}

	t.Log("✅ Test réussi")
}

func TestTypeDefinitionPrimaryKeyOrder(t *testing.T) {
	t.Log("🧪 TEST TYPE DEFINITION PRIMARY KEY ORDER")
	t.Log("==========================================")

	// L'ordre des champs de clé primaire DOIT être préservé
	typeDef := TypeDefinition{
		Name: "CompositeKey",
		Fields: []Field{
			{Name: "country", Type: "string", IsPrimaryKey: true},
			{Name: "region", Type: "string", IsPrimaryKey: true},
			{Name: "city", Type: "string", IsPrimaryKey: true},
			{Name: "name", Type: "string", IsPrimaryKey: false},
		},
	}

	expectedOrder := []string{"country", "region", "city"}
	pkNames := typeDef.GetPrimaryKeyFieldNames()

	if len(pkNames) != len(expectedOrder) {
		t.Fatalf("❌ Nombre de clés primaires: attendu %d, reçu %d",
			len(expectedOrder), len(pkNames))
	}

	for i, expected := range expectedOrder {
		if pkNames[i] != expected {
			t.Errorf("❌ Ordre PK[%d]: attendu '%s', reçu '%s'",
				i, expected, pkNames[i])
		}
	}

	t.Log("✅ L'ordre des clés primaires est préservé")
}

func TestBackwardCompatibilityJSON(t *testing.T) {
	t.Log("🧪 TEST BACKWARD COMPATIBILITY JSON")
	t.Log("====================================")

	// Ancien JSON sans isPrimaryKey
	oldJSON := `{
		"type": "typeDefinition",
		"name": "OldType",
		"fields": [
			{"name": "id", "type": "string"},
			{"name": "name", "type": "string"}
		]
	}`

	var typeDef TypeDefinition
	err := json.Unmarshal([]byte(oldJSON), &typeDef)
	if err != nil {
		t.Fatalf("❌ Erreur de désérialisation: %v", err)
	}

	// Tous les champs doivent avoir IsPrimaryKey = false par défaut
	for i, field := range typeDef.Fields {
		if field.IsPrimaryKey {
			t.Errorf("❌ Field[%d] '%s': IsPrimaryKey devrait être false par défaut, reçu true",
				i, field.Name)
		}
	}

	if typeDef.HasPrimaryKey() {
		t.Error("❌ Type sans isPrimaryKey ne devrait pas avoir de clé primaire")
	}

	t.Log("✅ Rétrocompatibilité JSON préservée")
}
