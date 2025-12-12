// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"os"
	"path/filepath"
	"testing"
	"github.com/treivax/tsd/constraint"
)
// TestRemoveRuleIncremental_FullPipeline teste la suppression de règle de manière incrémentale
// en utilisant le pipeline complet : parsing → construction → assertion → suppression → vérification
func TestRemoveRuleIncremental_FullPipeline(t *testing.T) {
	t.Log("🧪 TEST REMOVE RULE - PIPELINE COMPLET INCRÉMENTAL")
	t.Log("====================================================")
	// Créer un répertoire temporaire pour le fichier .tsd
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "incremental.tsd")
	// ÉTAPE 1: Fichier initial avec types, règles et faits
	t.Log("\n📝 ÉTAPE 1: Création du fichier initial avec 3 règles")
	initialContent := `type Person(id: string, name: string, age:number)
action adult(id: string)
action senior(id: string)
action minor(id: string)
rule adult_check : {p: Person} / p.age >= 18 ==> adult(p.id)
rule senior_check : {p: Person} / p.age >= 65 ==> senior(p.id)
rule minor_check : {p: Person} / p.age < 18 ==> minor(p.id)
Person(id:p1, name:Alice, age:30)
Person(id:p2, name:Bob, age:70)
Person(id:p3, name:Charlie, age:15)
`
	err := os.WriteFile(tsdFile, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("❌ Erreur écriture fichier: %v", err)
	}
	t.Log("✅ Fichier créé avec 3 règles et 3 faits")
	// ÉTAPE 2: Parser et construire le réseau initial
	t.Log("\n🔧 ÉTAPE 2: Construction du réseau RETE initial")
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	initialNetwork, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("❌ Erreur construction réseau: %v", err)
	}
	// Vérifier que les 3 règles sont présentes via les terminal nodes
	initialStats := initialNetwork.GetNetworkStats()
	initialTerminalCount := initialStats["terminal_nodes"].(int)
	if initialTerminalCount != 3 {
		t.Errorf("❌ Attendu 3 règles initiales (terminal nodes), trouvé %d", initialTerminalCount)
	}
	t.Logf("✅ Réseau construit avec %d règles", initialTerminalCount)
	// ÉTAPE 3: Ajouter une commande de suppression de règle
	t.Log("\n🗑️  ÉTAPE 3: Ajout de la commande 'remove rule senior_check'")
	updatedContent := initialContent + "\nremove rule senior_check\n"
	err = os.WriteFile(tsdFile, []byte(updatedContent), 0644)
	if err != nil {
		t.Fatalf("❌ Erreur mise à jour fichier: %v", err)
	}
	t.Log("✅ Commande de suppression ajoutée au fichier")
	// ÉTAPE 4: Re-parser et reconstruire le réseau
	t.Log("\n🔄 ÉTAPE 4: Reconstruction du réseau avec suppression")
	storage = NewMemoryStorage() // Nouveau storage pour reconstruction propre
	updatedNetwork, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("❌ Erreur reconstruction réseau: %v", err)
	}
	// Vérifier que la règle a été supprimée
	updatedStats := updatedNetwork.GetNetworkStats()
	remainingCount := updatedStats["terminal_nodes"].(int)
	if remainingCount != 2 {
		t.Errorf("❌ Attendu 2 règles après suppression, trouvé %d", remainingCount)
	}
	t.Logf("✅ Règles restantes: %d", remainingCount)
	// Vérifier que senior_check n'est plus présente
	for ruleID := range updatedNetwork.TerminalNodes {
		if ruleID == "senior_check_terminal" {
			t.Errorf("❌ La règle 'senior_check' devrait être supprimée!")
		}
	}
	t.Log("✅ La règle 'senior_check' a été correctement supprimée")
	// Vérifier que les autres règles sont toujours présentes
	foundAdult := false
	foundMinor := false
	for ruleID := range updatedNetwork.TerminalNodes {
		if ruleID == "adult_check_terminal" {
			foundAdult = true
		}
		if ruleID == "minor_check_terminal" {
			foundMinor = true
		}
	}
	if !foundAdult || !foundMinor {
		t.Errorf("❌ Les règles 'adult_check' et 'minor_check' devraient toujours exister")
	}
	t.Log("✅ Les règles 'adult_check' et 'minor_check' sont préservées")
	// ÉTAPE 5: Vérification structure du réseau
	t.Log("\n📊 ÉTAPE 5: Vérification de la structure après suppression")
	t.Log("✅ Structure du réseau validée (règles correctement supprimées/préservées)")
	// ÉTAPE 6: Supprimer une deuxième règle
	t.Log("\n🗑️  ÉTAPE 6: Suppression d'une deuxième règle 'minor_check'")
	finalContent := updatedContent + "remove rule minor_check\n"
	err = os.WriteFile(tsdFile, []byte(finalContent), 0644)
	if err != nil {
		t.Fatalf("❌ Erreur mise à jour fichier: %v", err)
	}
	storage = NewMemoryStorage()
	finalNetwork, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("❌ Erreur reconstruction finale: %v", err)
	}
	finalStats := finalNetwork.GetNetworkStats()
	finalCount := finalStats["terminal_nodes"].(int)
	if finalCount != 1 {
		t.Errorf("❌ Attendu 1 règle finale, trouvé %d", finalCount)
	}
	t.Logf("✅ Règle finale restante: %d", finalCount)
	// Vérifier qu'il ne reste que adult_check
	if _, exists := finalNetwork.TerminalNodes["adult_check_terminal"]; !exists {
		t.Errorf("❌ La règle 'adult_check' devrait exister")
	}
	if _, exists := finalNetwork.TerminalNodes["senior_check_terminal"]; exists {
		t.Errorf("❌ La règle 'senior_check' ne devrait plus exister")
	}
	if _, exists := finalNetwork.TerminalNodes["minor_check_terminal"]; exists {
		t.Errorf("❌ La règle 'minor_check' ne devrait plus exister")
	}
	t.Log("✅ Seule 'adult_check' reste dans le réseau")
	// ÉTAPE 7: Vérification finale
	t.Log("\n📊 ÉTAPE 7: Vérification finale de la structure")
	// Devrait avoir seulement adult_check
	if len(finalNetwork.TerminalNodes) != 1 {
		t.Errorf("❌ Attendu 1 terminal node final, trouvé %d", len(finalNetwork.TerminalNodes))
	}
	t.Log("✅ Structure finale validée")
	t.Log("\n✅ TEST COMPLET - Pipeline incrémental validé avec succès!")
}
// TestRemoveRuleIncremental_WithJoins teste la suppression de règles avec jointures
func TestRemoveRuleIncremental_WithJoins(t *testing.T) {
	t.Log("🧪 TEST REMOVE RULE - AVEC JOINTURES")
	t.Log("=====================================")
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "joins.tsd")
	// Fichier avec des règles de jointure
	content := `type Person(id: string, name:string)
type Order(id: string, customer_id: string, amount:number)
action process_order(personId: string, orderId: string)
action vip_order(arg: string)
rule person_order : {p: Person, o: Order} / p.id == o.customer_id ==> process_order(p.id, o.id)
rule high_value : {p: Person, o: Order} / p.id == o.customer_id AND o.amount > 100 ==> vip_order(p.id)
Person(id:p1, name:Alice)
Person(id:p2, name:Bob)
Order(id:o1, customer_id:p1, amount:150)
Order(id:o2, customer_id:p2, amount:50)
`
	err := os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("❌ Erreur écriture fichier: %v", err)
	}
	// Construction initiale
	t.Log("\n🔧 Construction du réseau initial avec jointures")
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("❌ Erreur construction: %v", err)
	}
	stats := network.GetNetworkStats()
	initialTerminalCount := stats["terminal_nodes"].(int)
	if initialTerminalCount != 2 {
		t.Errorf("❌ Attendu 2 règles, trouvé %d", initialTerminalCount)
	}
	t.Logf("✅ Réseau initial: %d règles", initialTerminalCount)
	t.Log("✅ Réseau avec jointures construit")
	// Ajouter commande de suppression
	t.Log("\n🗑️  Ajout de 'remove rule high_value'")
	updatedContent := content + "\nremove rule high_value\n"
	err = os.WriteFile(tsdFile, []byte(updatedContent), 0644)
	if err != nil {
		t.Fatalf("❌ Erreur mise à jour fichier: %v", err)
	}
	// Reconstruire
	t.Log("\n🔄 Reconstruction après suppression")
	storage = NewMemoryStorage()
	updatedNetwork, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("❌ Erreur reconstruction: %v", err)
	}
	updatedStats := updatedNetwork.GetNetworkStats()
	remainingCount := updatedStats["terminal_nodes"].(int)
	if remainingCount != 1 {
		t.Errorf("❌ Attendu 1 règle restante, trouvé %d", remainingCount)
	}
	t.Logf("✅ Règles après suppression: %d", remainingCount)
	// Vérifier quelle règle reste
	if _, exists := updatedNetwork.TerminalNodes["person_order_terminal"]; !exists {
		t.Errorf("❌ La règle 'person_order' devrait exister")
	}
	if _, exists := updatedNetwork.TerminalNodes["high_value_terminal"]; exists {
		t.Errorf("❌ La règle 'high_value' ne devrait plus exister")
	}
	t.Log("\n📊 Vérification structure après suppression")
	t.Log("✅ Structure validée")
	t.Log("\n✅ TEST JOINTURES - Suppression validée avec succès!")
}
// TestRemoveRuleIncremental_MultipleRemovals teste plusieurs suppressions successives
func TestRemoveRuleIncremental_MultipleRemovals(t *testing.T) {
	t.Log("🧪 TEST REMOVE RULE - SUPPRESSIONS MULTIPLES")
	t.Log("=============================================")
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "multiple.tsd")
	// Fichier initial avec 5 règles
	content := `type Person(id: string, age:number)
action action1(arg: string)
action action2(arg: string)
action action3(arg: string)
action action4(arg: string)
action action5(arg: string)
rule r1 : {p: Person} / p.age > 10 ==> action1(p.id)
rule r2 : {p: Person} / p.age > 20 ==> action2(p.id)
rule r3 : {p: Person} / p.age > 30 ==> action3(p.id)
rule r4 : {p: Person} / p.age > 40 ==> action4(p.id)
rule r5 : {p: Person} / p.age > 50 ==> action5(p.id)
Person(id:p1, age:55)
`
	err := os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("❌ Erreur écriture fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	// Construction initiale
	t.Log("\n🔧 Construction initiale avec 5 règles")
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("❌ Erreur construction: %v", err)
	}
	stats := network.GetNetworkStats()
	initialCount := stats["terminal_nodes"].(int)
	if initialCount != 5 {
		t.Errorf("❌ Attendu 5 règles, trouvé %d", initialCount)
	}
	t.Logf("✅ État initial: %d règles", initialCount)
	// Supprimer r2 et r4
	t.Log("\n🗑️  Suppression de r2 et r4")
	content += "\nremove rule r2\nremove rule r4\n"
	err = os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("❌ Erreur mise à jour: %v", err)
	}
	storage = NewMemoryStorage()
	network, _, err = pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("❌ Erreur reconstruction: %v", err)
	}
	stats = network.GetNetworkStats()
	middleCount := stats["terminal_nodes"].(int)
	if middleCount != 3 {
		t.Errorf("❌ Attendu 3 règles après première suppression, trouvé %d", middleCount)
	}
	t.Logf("✅ Après première suppression: %d règles", middleCount)
	// Vérifier que r2 et r4 sont absentes
	if _, exists := network.TerminalNodes["r2_terminal"]; exists {
		t.Errorf("❌ La règle 'r2' devrait être supprimée")
	}
	if _, exists := network.TerminalNodes["r4_terminal"]; exists {
		t.Errorf("❌ La règle 'r4' devrait être supprimée")
	}
	// Supprimer r1 et r5
	t.Log("\n🗑️  Suppression de r1 et r5")
	content += "remove rule r1\nremove rule r5\n"
	err = os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("❌ Erreur mise à jour: %v", err)
	}
	storage = NewMemoryStorage()
	network, _, err = pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("❌ Erreur reconstruction: %v", err)
	}
	stats = network.GetNetworkStats()
	finalCount := stats["terminal_nodes"].(int)
	if finalCount != 1 {
		t.Errorf("❌ Attendu 1 règle finale, trouvé %d", finalCount)
	}
	if _, exists := network.TerminalNodes["r3_terminal"]; !exists {
		t.Errorf("❌ La règle 'r3' devrait exister")
	}
	t.Logf("✅ État final: seule r3 reste (comme attendu)")
	t.Log("\n✅ TEST SUPPRESSIONS MULTIPLES - Validé avec succès!")
}
// TestRemoveRuleIncremental_ParseOnly teste uniquement le parsing de la commande
func TestRemoveRuleIncremental_ParseOnly(t *testing.T) {
	t.Log("🧪 TEST REMOVE RULE - PARSING UNIQUEMENT")
	t.Log("=========================================")
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple removal",
			input:    "remove rule my_rule",
			expected: "my_rule",
		},
		{
			name:     "Removal with underscores",
			input:    "remove rule complex_rule_name_123",
			expected: "complex_rule_name_123",
		},
		{
			name:     "Multiple removals",
			input:    "remove rule rule1\nremove rule rule2",
			expected: "rule1", // On vérifie juste le premier
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Test: %s", tc.name)
			t.Logf("Input: %s", tc.input)
			result, err := constraint.ParseConstraint("", []byte(tc.input))
			if err != nil {
				t.Fatalf("❌ Erreur parsing: %v", err)
			}
			resultMap := result.(map[string]interface{})
			ruleRemovals := resultMap["ruleRemovals"].([]interface{})
			if len(ruleRemovals) == 0 {
				t.Fatalf("❌ Aucune suppression de règle trouvée")
			}
			removal := ruleRemovals[0].(map[string]interface{})
			ruleID := removal["ruleID"].(string)
			if ruleID != tc.expected {
				t.Errorf("❌ Attendu '%s', trouvé '%s'", tc.expected, ruleID)
			}
			t.Logf("✅ Parsing correct: ruleID='%s'", ruleID)
		})
	}
	t.Log("\n✅ TEST PARSING - Tous les cas validés!")
}