// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"io"
	"strings"
	"testing"
)

func TestUnicodeSupport(t *testing.T) {
	// Test avec des caractères accentués français
	input := `Utilisateur(id:U001, nom:"François", prenom:"José", age:25)`

	result, err := parseConstraintGrammar(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Erreur de parsing avec caractères accentués: %v", err)
	}

	// Vérifier que le résultat contient bien les faits
	facts := result.(map[string]interface{})["facts"].([]interface{})
	if len(facts) != 1 {
		t.Fatalf("Attendu 1 fait, obtenu %d", len(facts))
	}

	fact := facts[0].(map[string]interface{})
	if fact["typeName"] != "Utilisateur" {
		t.Errorf("Attendu typeName 'Utilisateur', obtenu '%s'", fact["typeName"])
	}

	fields := fact["fields"].([]interface{})
	if len(fields) != 4 {
		t.Fatalf("Attendu 4 champs, obtenu %d", len(fields))
	}

	// Vérifier que les valeurs avec accents sont correctement parsées
	nomField := fields[1].(map[string]interface{})
	if nomField["name"] != "nom" {
		t.Errorf("Attendu nom de champ 'nom', obtenu '%s'", nomField["name"])
	}

	nomValue := nomField["value"].(map[string]interface{})
	if nomValue["value"] != "François" {
		t.Errorf("Attendu valeur 'François', obtenu '%s'", nomValue["value"])
	}

	t.Logf("✅ Test Unicode réussi: %+v", fact)
}

func TestTraitDUnionSupport(t *testing.T) {
	// Test avec des traits d'union dans les identifiants
	input := `Adresse(utilisateur-id:U001, rue:"Rue d'Alsace-Lorraine", ville:"Aix-en-Provence")`

	result, err := parseConstraintGrammar(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Erreur de parsing avec traits d'union: %v", err)
	}

	// Vérifier que le résultat contient bien les faits
	facts := result.(map[string]interface{})["facts"].([]interface{})
	if len(facts) != 1 {
		t.Fatalf("Attendu 1 fait, obtenu %d", len(facts))
	}

	fact := facts[0].(map[string]interface{})
	if fact["typeName"] != "Adresse" {
		t.Errorf("Attendu typeName 'Adresse', obtenu '%s'", fact["typeName"])
	}

	t.Logf("✅ Test trait d'union réussi: %+v", fact)
}

func TestCommentaireSupport(t *testing.T) {
	// Test avec différents types de commentaires
	input := `// Commentaire ligne
	/* Commentaire
	   multi-lignes */
	Utilisateur(id:U001, nom:"Test", age:25)  # Commentaire fin de ligne
	`

	result, err := parseConstraintGrammar(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Erreur de parsing avec commentaires: %v", err)
	}

	// Vérifier que le résultat contient bien les faits (commentaires ignorés)
	facts := result.(map[string]interface{})["facts"].([]interface{})
	if len(facts) != 1 {
		t.Fatalf("Attendu 1 fait, obtenu %d", len(facts))
	}

	fact := facts[0].(map[string]interface{})
	if fact["typeName"] != "Utilisateur" {
		t.Errorf("Attendu typeName 'Utilisateur', obtenu '%s'", fact["typeName"])
	}

	t.Logf("✅ Test commentaires réussi: %+v", fact)
}

func parseConstraintGrammar(input *strings.Reader) (interface{}, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	return Parse("test", data)
}
