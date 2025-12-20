// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

func TestPrimaryKeyIntegration(t *testing.T) {
	t.Log("🧪 TEST PRIMARY KEY INTEGRATION")
	t.Log("================================")

	t.Run("type avec clé primaire simple - valide", func(t *testing.T) {
		input := `
type User(#login: string, name: string)
User(login: "alice", name: "Alice")
`
		result, err := Parse("test.tsd", []byte(input))
		if err != nil {
			t.Fatalf("❌ Erreur de parsing: %v", err)
		}

		// Valider le programme
		err = ValidateConstraintProgram(result)
		if err != nil {
			t.Fatalf("❌ Erreur de validation: %v", err)
		}

		// Vérifier que le programme est correct
		program, err := ConvertResultToProgram(result)
		if err != nil {
			t.Fatalf("❌ Erreur de conversion: %v", err)
		}

		// Vérifier que le type a bien une clé primaire
		if len(program.Types) != 1 {
			t.Fatalf("❌ Attendu 1 type, reçu %d", len(program.Types))
		}

		if !program.Types[0].HasPrimaryKey() {
			t.Error("❌ Le type devrait avoir une clé primaire")
		}

		pkFields := program.Types[0].GetPrimaryKeyFieldNames()
		if len(pkFields) != 1 || pkFields[0] != "login" {
			t.Errorf("❌ Attendu clé primaire 'login', reçu %v", pkFields)
		}

		t.Log("✅ Test réussi")
	})

	t.Run("type avec clé primaire composite - valide", func(t *testing.T) {
		input := `
type Person(#firstName: string, #lastName: string, age: number)
Person(firstName: "John", lastName: "Doe", age: 30)
`
		result, err := Parse("test.tsd", []byte(input))
		if err != nil {
			t.Fatalf("❌ Erreur de parsing: %v", err)
		}

		err = ValidateConstraintProgram(result)
		if err != nil {
			t.Fatalf("❌ Erreur de validation: %v", err)
		}

		program, err := ConvertResultToProgram(result)
		if err != nil {
			t.Fatalf("❌ Erreur de conversion: %v", err)
		}

		pkFields := program.Types[0].GetPrimaryKeyFieldNames()
		if len(pkFields) != 2 {
			t.Errorf("❌ Attendu 2 champs de clé primaire, reçu %d", len(pkFields))
		}

		t.Log("✅ Test réussi")
	})

	t.Run("fait sans champ PK - invalide", func(t *testing.T) {
		input := `
type User(#login: string, name: string)
User(name: "Bob")
`
		result, err := Parse("test.tsd", []byte(input))
		if err != nil {
			t.Fatalf("❌ Erreur de parsing: %v", err)
		}

		err = ValidateConstraintProgram(result)
		if err == nil {
			t.Error("❌ Attendu une erreur de validation, reçu nil")
		} else if !strings.Contains(err.Error(), "manquants") {
			t.Errorf("❌ Erreur inattendue: %v", err)
		} else {
			t.Logf("✅ Erreur attendue: %v", err)
		}
	})

	t.Run("fait avec id manuel - invalide", func(t *testing.T) {
		input := `
type User(#login: string, name: string)
User(id: "manual-id", login: "alice", name: "Alice")
`
		result, err := Parse("test.tsd", []byte(input))
		if err != nil {
			t.Fatalf("❌ Erreur de parsing: %v", err)
		}

		err = ValidateConstraintProgram(result)
		if err == nil {
			t.Error("❌ Attendu une erreur de validation, reçu nil")
		} else if !strings.Contains(err.Error(), "non défini dans le type") {
			t.Errorf("❌ Erreur inattendue: %v", err)
		} else {
			t.Logf("✅ Erreur attendue: %v", err)
		}
	})

	t.Run("fait avec valeur PK nulle - invalide", func(t *testing.T) {
		// Note: Le parser ne permet pas les valeurs null,
		// mais on teste avec une string vide
		input := `
type User(#login: string, name: string)
User(login: "", name: "Bob")
`
		result, err := Parse("test.tsd", []byte(input))
		if err != nil {
			t.Fatalf("❌ Erreur de parsing: %v", err)
		}

		err = ValidateConstraintProgram(result)
		if err == nil {
			t.Error("❌ Attendu une erreur de validation, reçu nil")
		} else if !strings.Contains(err.Error(), "ne peut pas être vide") {
			t.Errorf("❌ Erreur inattendue: %v", err)
		} else {
			t.Logf("✅ Erreur attendue: %v", err)
		}
	})

	t.Run("type avec clé primaire de type complexe - invalide", func(t *testing.T) {
		// On ne peut pas vraiment tester ceci avec le parser car il faudrait
		// définir d'abord un type complexe. Mais la validation fonctionnera
		// si jamais un tel cas se présente.
		t.Skip("Validation de type complexe PK testée dans les tests unitaires")
	})

	t.Run("type sans clé primaire - valide", func(t *testing.T) {
		input := `
type Document(title: string, content: string)
Document(title: "Doc1", content: "Content")
`
		result, err := Parse("test.tsd", []byte(input))
		if err != nil {
			t.Fatalf("❌ Erreur de parsing: %v", err)
		}

		err = ValidateConstraintProgram(result)
		if err != nil {
			t.Fatalf("❌ Erreur de validation: %v", err)
		}

		program, err := ConvertResultToProgram(result)
		if err != nil {
			t.Fatalf("❌ Erreur de conversion: %v", err)
		}

		if program.Types[0].HasPrimaryKey() {
			t.Error("❌ Le type ne devrait pas avoir de clé primaire")
		}

		t.Log("✅ Test réussi")
	})
}
