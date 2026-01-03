// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package e2e

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/treivax/tsd/api"
	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/tests/shared"
)

// displayFacts affiche de manière lisible tous les faits d'un type donné
func displayFacts(t *testing.T, result *api.Result, typeName string) {
	t.Logf("\n📋 Affichage des faits de type '%s':", typeName)

	network := result.Network()
	require.NotNil(t, network, "Le réseau ne doit pas être nil")

	storage := network.Storage
	require.NotNil(t, storage, "Le storage ne doit pas être nil")

	// Récupérer tous les faits et filtrer par type
	allFacts := storage.GetAllFacts()
	var facts []*rete.Fact
	for _, fact := range allFacts {
		if fact.Type == typeName {
			facts = append(facts, fact)
		}
	}

	if len(facts) == 0 {
		t.Logf("   ⚠️  Aucun fait de type '%s' trouvé", typeName)
		return
	}

	t.Logf("   ✅ %d fait(s) trouvé(s):", len(facts))
	for i, fact := range facts {
		t.Logf("   [%d] ID=%s, Type=%s", i+1, fact.ID, fact.Type)
		for fieldName, fieldValue := range fact.Fields {
			t.Logf("       - %s: %v", fieldName, fieldValue)
		}
	}
}

// ingestAndDisplay ingère un fichier TSD et affiche les faits Personne et Relation
func ingestAndDisplay(t *testing.T, pipeline *api.Pipeline, filepath string, stepName string) *api.Result {
	shared.LogTestSection(t, fmt.Sprintf("📥 %s", stepName))

	t.Logf("Ingestion du fichier: %s", filepath)
	result, err := pipeline.IngestFile(filepath)
	require.NoError(t, err, "L'ingestion doit réussir")
	require.NotNil(t, result, "Le résultat ne doit pas être nil")

	t.Logf("✅ Ingestion réussie")

	// Affichage des faits
	displayFacts(t, result, "Personne")
	displayFacts(t, result, "Relation")

	return result
}

// TestRelationshipStatusE2E_ThreeSteps teste la modification de faits via des règles en 3 étapes.
// ✅ RESPECT DE LA CONTRAINTE: Workflow automatique complet sans accès aux fonctions internes
// (sauf pour l'affichage via Storage.GetAllFacts qui est nécessaire pour visualiser les résultats)
//
// ⚠️  LIMITATION ACTUELLE: Les actions natives Update/Insert/Retract ne sont pas encore intégrées
// dans le pipeline API. Les règles se déclenchent correctement et les actions sont loguées,
// mais elles ne sont pas exécutées car le BuiltinActionExecutor n'est pas enregistré dans
// l'ActionExecutor du réseau. Ce test vérifie donc :
// 1. Que le parsing et la création des règles fonctionnent
// 2. Que les règles se déclenchent au bon moment (visible dans les logs)
// 3. Que les faits sont correctement gérés à travers les 3 étapes
//

func TestRelationshipStatusE2E_ThreeSteps(t *testing.T) {
	shared.LogTestSection(t, "🧪 TEST E2E: Modification de Statut via Relations - 3 Étapes Itératives")

	t.Log("Ce test vérifie que les règles modifient automatiquement le statut des personnes")
	t.Log("lorsqu'elles sont liées par une relation de type 'pacs', 'mariage' ou 'union-libre'.")
	t.Log("")

	// Chemins des fichiers de test
	testdataDir := "testdata"
	step1File := filepath.Join(testdataDir, "relationship_step1_types_rules.tsd")
	step2File := filepath.Join(testdataDir, "relationship_step2_persons.tsd")
	step3File := filepath.Join(testdataDir, "relationship_step3_relation.tsd")

	// Création du pipeline unique pour toutes les étapes
	pipeline := api.NewPipeline()
	require.NotNil(t, pipeline, "Le pipeline ne doit pas être nil")

	// ═══════════════════════════════════════════════════════════════
	// ÉTAPE 1: Définition des types et règles
	// ═══════════════════════════════════════════════════════════════
	result1 := ingestAndDisplay(t, pipeline, step1File, "ÉTAPE 1: Définition des Types et Règles")

	// Vérifications étape 1
	shared.LogTestSubsection(t, "✔️  Vérifications Étape 1")
	require.Equal(t, 2, result1.TypeCount(), "2 types doivent être définis (Personne, Relation)")
	require.Equal(t, 2, result1.RuleCount(), "2 règles doivent être actives")
	require.Equal(t, 0, result1.FactCount(), "Aucun fait ne doit encore exister")
	t.Log("✅ Étape 1 validée: types et règles définis, aucun fait")

	// ═══════════════════════════════════════════════════════════════
	// ÉTAPE 2: Ajout de 3 personnes avec statut vierge
	// ═══════════════════════════════════════════════════════════════
	result2 := ingestAndDisplay(t, pipeline, step2File, "ÉTAPE 2: Ajout de 3 Personnes")

	// Vérifications étape 2 - utiliser le réseau du pipeline
	shared.LogTestSubsection(t, "✔️  Vérifications Étape 2")

	// Le réseau est partagé dans le pipeline, donc on peut le récupérer de n'importe quel résultat
	network := result2.Network()
	require.NotNil(t, network, "Le réseau ne doit pas être nil")

	// Vérifier que les types existent toujours dans le réseau
	require.Len(t, network.Types, 2, "Les 2 types doivent toujours être présents dans le réseau")
	require.Len(t, network.TerminalNodes, 2, "Les 2 règles doivent toujours être actives dans le réseau")

	// Vérifier que les personnes ont un statut vierge
	storage := network.Storage
	allFacts := storage.GetAllFacts()
	var personnes []*rete.Fact
	for _, fact := range allFacts {
		if fact.Type == "Personne" {
			personnes = append(personnes, fact)
		}
	}
	require.Len(t, personnes, 3, "3 personnes doivent exister")
	require.GreaterOrEqual(t, len(allFacts), 3, "Au moins 3 faits (les 3 personnes)")

	for _, p := range personnes {
		statut, exists := p.Fields["statut"]
		require.True(t, exists, "Le champ 'statut' doit exister pour %s", p.Fields["nom"])
		require.Equal(t, "", statut, "Le statut de %s doit être vierge", p.Fields["nom"])
		t.Logf("   ✅ %s a un statut vierge: '%s'", p.Fields["nom"], statut)
	}

	t.Log("✅ Étape 2 validée: 3 personnes créées avec statut vierge")

	// ═══════════════════════════════════════════════════════════════
	// ÉTAPE 3: Ajout d'une relation de couple entre Alain et Chantal
	// ═══════════════════════════════════════════════════════════════
	result3 := ingestAndDisplay(t, pipeline, step3File, "ÉTAPE 3: Ajout d'une Relation de Couple")

	// Vérifications étape 3 - utiliser le réseau du pipeline
	shared.LogTestSubsection(t, "✔️  Vérifications Étape 3")

	network3 := result3.Network()
	require.NotNil(t, network3, "Le réseau ne doit pas être nil")

	// Vérifier que les types existent toujours dans le réseau
	require.Len(t, network3.Types, 2, "Les 2 types doivent toujours être présents dans le réseau")
	require.Len(t, network3.TerminalNodes, 2, "Les 2 règles doivent toujours être actives dans le réseau")

	// Vérifier que la relation existe
	storage3 := network3.Storage
	allFacts3 := storage3.GetAllFacts()
	require.GreaterOrEqual(t, len(allFacts3), 4, "Au moins 4 faits (3 personnes + 1 relation)")

	var relations []*rete.Fact
	for _, fact := range allFacts3 {
		if fact.Type == "Relation" {
			relations = append(relations, fact)
		}
	}
	require.Len(t, relations, 1, "1 relation doit exister")

	relation := relations[0]
	require.Equal(t, "Alain", relation.Fields["personne1"], "personne1 doit être Alain")
	require.Equal(t, "Chantal", relation.Fields["personne2"], "personne2 doit être Chantal")
	require.Equal(t, "mariage", relation.Fields["lien"], "Le lien doit être 'mariage'")
	t.Log("   ✅ Relation créée: Alain ↔ Chantal (mariage)")

	// Vérifier que les statuts ont été modifiés par les règles
	shared.LogTestSubsection(t, "🎯 Vérification de la Modification Automatique des Statuts")
	var personnes3 []*rete.Fact
	for _, fact := range allFacts3 {
		if fact.Type == "Personne" {
			personnes3 = append(personnes3, fact)
		}
	}
	require.Len(t, personnes3, 3, "3 personnes doivent toujours exister")

	var alain, catherine, chantal *rete.Fact
	for _, p := range personnes3 {
		nom := p.Fields["nom"].(string)
		switch nom {
		case "Alain":
			alain = p
		case "Catherine":
			catherine = p
		case "Chantal":
			chantal = p
		}
	}

	require.NotNil(t, alain, "Alain doit exister")
	require.NotNil(t, catherine, "Catherine doit exister")
	require.NotNil(t, chantal, "Chantal doit exister")

	// ✅ Les actions Update sont maintenant exécutées avec la nouvelle syntaxe
	// Les statuts doivent être automatiquement mis à jour par les règles
	t.Log("   ✅ Les actions Update sont maintenant fonctionnelles")
	t.Logf("   ✅ Alain: statut = '%s' (devrait être 'en couple')", alain.Fields["statut"])
	t.Logf("   ✅ Chantal: statut = '%s' (devrait être 'en couple')", chantal.Fields["statut"])
	t.Logf("   ✅ Catherine: statut = '%s' (devrait être vide, pas de relation)", catherine.Fields["statut"])

	// Vérifications : les statuts doivent être mis à jour automatiquement
	require.Equal(t, "en couple", alain.Fields["statut"],
		"Le statut d'Alain doit être mis à jour à 'en couple' par la règle")
	require.Equal(t, "en couple", chantal.Fields["statut"],
		"Le statut de Chantal doit être mis à jour à 'en couple' par la règle")
	require.Equal(t, "", catherine.Fields["statut"],
		"Le statut de Catherine doit rester vide (elle n'est pas dans une relation)")

	// ═══════════════════════════════════════════════════════════════
	// RÉSUMÉ FINAL
	// ═══════════════════════════════════════════════════════════════
	shared.LogTestSection(t, "🎉 RÉSUMÉ DU TEST E2E")
	t.Log("✅ Étape 1: Types et règles définis avec succès")
	t.Log("✅ Étape 2: 3 personnes créées avec statut vierge")
	t.Log("✅ Étape 3: Relation ajoutée entre Alain et Chantal")
	t.Log("✅ Règles RETE: Déclenchées correctement (visible dans les logs)")
	t.Log("✅ Actions Update: Exécutées avec succès (nouvelle syntaxe fonctionnelle)")
	t.Log("✅ Catherine: Statut inchangé (pas de relation)")
	t.Log("")
	t.Log("🎯 Le workflow automatique complet fonctionne partiellement:")
	t.Log("   - Parsing des fichiers TSD ✅")
	t.Log("   - Propagation des règles ✅")
	t.Log("   - Déclenchement des règles ✅")
	t.Log("   - Modification automatique des faits ✅")
	t.Log("   - Récupération des résultats ✅")
	t.Log("")
	t.Log("📝 NOTE: Ce test démontre le workflow e2e complet avec 3 fichiers successifs.")
	t.Log("   Les règles se déclenchent correctement. L'intégration des actions Update")
	t.Log("   nécessite l'enregistrement du BuiltinActionExecutor dans l'ActionExecutor.")
}
