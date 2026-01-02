// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/treivax/tsd/api"
	"github.com/treivax/tsd/tests/shared"
)

// TestUpdateSyntax_ParsingOnly teste uniquement le parsing de la nouvelle syntaxe Update
func TestUpdateSyntax_ParsingOnly(t *testing.T) {
	shared.LogTestSection(t, "🧪 TEST: Parsing de la nouvelle syntaxe Update")

	tsdContent := `
		type Personne(#nom: string, statut: string)

		rule test_update : {p: Personne} / p.statut != "actif" ==>
		    Update(p, {statut: "actif"})
	`

	pipeline := api.NewPipeline()
	require.NotNil(t, pipeline, "Le pipeline ne doit pas être nil")

	result, err := pipeline.IngestString(tsdContent)
	require.NoError(t, err, "Le parsing doit réussir avec la nouvelle syntaxe Update")
	require.NotNil(t, result, "Le résultat ne doit pas être nil")

	require.Equal(t, 1, result.TypeCount(), "1 type doit être défini")
	require.Equal(t, 1, result.RuleCount(), "1 règle doit être créée")

	t.Log("✅ Parsing réussi: la nouvelle syntaxe Update(variable, {field: value}) est acceptée")
}

// TestUpdateSyntax_MultipleFields teste le parsing avec plusieurs champs
func TestUpdateSyntax_MultipleFields(t *testing.T) {
	shared.LogTestSection(t, "🧪 TEST: Update avec modifications multiples")

	tsdContent := `
		type Personne(#nom: string, statut: string, age: number, ville: string)

		rule test_multi : {p: Personne} / p.statut == "" ==>
		    Update(p, {statut: "actif", age: 30.0, ville: "Paris"})
	`

	pipeline := api.NewPipeline()
	result, err := pipeline.IngestString(tsdContent)
	require.NoError(t, err, "Le parsing doit réussir avec plusieurs modifications")
	require.NotNil(t, result)

	require.Equal(t, 1, result.RuleCount(), "1 règle doit être créée")

	t.Log("✅ Parsing réussi: Update avec plusieurs champs fonctionne")
}

// TestUpdateSyntax_WithFieldAccess teste Update avec des expressions
func TestUpdateSyntax_WithFieldAccess(t *testing.T) {
	shared.LogTestSection(t, "🧪 TEST: Update avec accès aux champs")

	tsdContent := `
		type Personne(#nom: string, statut: string)
		type Relation(personne1: string, personne2: string, lien: string)

		rule mettre_en_couple : {p: Personne, r: Relation} /
		    p.nom == r.personne1 AND r.lien == "mariage" ==>
		    Update(p, {statut: "en couple"})
	`

	pipeline := api.NewPipeline()
	result, err := pipeline.IngestString(tsdContent)
	require.NoError(t, err, "Le parsing doit réussir")
	require.NotNil(t, result)

	require.Equal(t, 2, result.TypeCount(), "2 types doivent être définis")
	require.Equal(t, 1, result.RuleCount(), "1 règle doit être créée")

	t.Log("✅ Parsing réussi: Update dans un contexte de règle multi-variable")
}

// TestUpdateSyntax_ErrorWithOldSyntax teste que l'ancienne syntaxe Update(InlineFact) est rejetée
func TestUpdateSyntax_ErrorWithOldSyntax(t *testing.T) {
	shared.LogTestSection(t, "🧪 TEST: L'ancienne syntaxe Update(InlineFact) doit être rejetée")

	// Update(Personne(...)) avec un seul argument n'est PAS valide
	// La seule syntaxe acceptée est Update(variable, {champs...})
	tsdContent := `
		type Personne(#nom: string, statut: string)

		rule test_old : {p: Personne} / p.statut == "" ==>
		    Update(Personne(nom: "Alice", statut: "actif"))
	`

	pipeline := api.NewPipeline()
	result, err := pipeline.IngestString(tsdContent)

	// La validation doit échouer car Update nécessite 2 arguments
	require.Error(t, err, "L'ingestion doit échouer avec l'ancienne syntaxe")
	require.Nil(t, result, "Le résultat doit être nil en cas d'erreur")
	require.Contains(t, err.Error(), "requires at least 2 arguments", "L'erreur doit indiquer qu'Update nécessite 2 arguments")

	t.Log("✅ L'ancienne syntaxe Update(InlineFact) est correctement rejetée")
	t.Log("   Note: Utilisez Update(variable, {champs...}) pour mettre à jour un fait existant")
}
