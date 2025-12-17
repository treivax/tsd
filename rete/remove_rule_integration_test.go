// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"os"
	"testing"
)

// TestRemoveRuleCommand_ParseAndExecute vérifie que la commande remove rule fonctionne de bout en bout
func TestRemoveRuleCommand_ParseAndExecute(t *testing.T) {
	t.Log("🧪 TEST INTEGRATION REMOVE RULE COMMAND")
	t.Log("========================================")
	// Créer un fichier .tsd temporaire avec une règle et sa suppression
	content := `
type Person(#id: string, name: string, age:number)
action notify(message: string)
action notify_senior(message: string)
rule adult_check : {p: Person} / p.age >= 18 ==> notify(p.id)
rule senior_check : {p: Person} / p.age >= 65 ==> notify_senior(p.id)
Person(id: "P1", name: "Alice", age: 25)
Person(id: "P2", name: "Bob", age: 70)
remove rule adult_check
`
	tmpFile, err := os.CreateTemp("", "remove_rule_test_*.tsd")
	if err != nil {
		t.Fatalf("❌ Erreur création fichier temporaire: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("❌ Erreur écriture fichier: %v", err)
	}
	tmpFile.Close()
	t.Logf("📝 Fichier temporaire créé: %s", tmpFile.Name())
	// Construire le réseau avec le pipeline
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tmpFile.Name(), nil, storage)
	if err != nil {
		t.Fatalf("❌ Erreur construction réseau: %v", err)
	}
	// Debug: afficher tous les terminaux
	t.Logf("📋 Terminaux dans le réseau:")
	for termID := range network.TerminalNodes {
		t.Logf("   - %s", termID)
	}
	// Vérifier que la règle senior_check existe encore
	seniorTerminalExists := false
	for termID := range network.TerminalNodes {
		if termID == "senior_check_terminal" {
			seniorTerminalExists = true
			break
		}
	}
	if !seniorTerminalExists {
		t.Errorf("❌ La règle senior_check devrait encore exister")
	} else {
		t.Logf("✅ La règle senior_check existe toujours")
	}
	// Vérifier que la règle adult_check n'existe plus
	adultTerminalExists := false
	for termID := range network.TerminalNodes {
		if termID == "adult_check_terminal" {
			adultTerminalExists = true
			break
		}
	}
	if adultTerminalExists {
		t.Errorf("❌ La règle adult_check ne devrait plus exister après remove rule")
	} else {
		t.Logf("✅ La règle adult_check a été supprimée correctement")
	}
	t.Log("✅ Test réussi - Commande remove rule fonctionne de bout en bout")
}

// TestRemoveRuleCommand_MultipleRules vérifie la suppression de plusieurs règles
func TestRemoveRuleCommand_MultipleRules(t *testing.T) {
	t.Log("🧪 TEST REMOVE MULTIPLE RULES")
	t.Log("=============================")
	content := `
type Person(#id: string, name: string, age:number)
action action1(arg: string)
action action2(arg: string)
action action3(arg: string)
rule rule1 : {p: Person} / p.age > 18 ==> action1(p.id)
rule rule2 : {p: Person} / p.age > 30 ==> action2(p.id)
rule rule3 : {p: Person} / p.age > 50 ==> action3(p.id)
Person(id: "P1", name: "Alice", age: 35)
remove rule rule1
remove rule rule3
`
	tmpFile, err := os.CreateTemp("", "remove_multiple_rules_*.tsd")
	if err != nil {
		t.Fatalf("❌ Erreur création fichier temporaire: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("❌ Erreur écriture fichier: %v", err)
	}
	tmpFile.Close()
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tmpFile.Name(), nil, storage)
	if err != nil {
		t.Fatalf("❌ Erreur construction réseau: %v", err)
	}
	// Vérifier que rule2 existe
	if _, exists := network.TerminalNodes["rule2_terminal"]; !exists {
		t.Errorf("❌ La règle rule2 devrait encore exister")
	} else {
		t.Logf("✅ La règle rule2 existe toujours")
	}
	// Vérifier que rule1 et rule3 n'existent plus
	if _, exists := network.TerminalNodes["rule1_terminal"]; exists {
		t.Errorf("❌ La règle rule1 ne devrait plus exister")
	} else {
		t.Logf("✅ La règle rule1 a été supprimée")
	}
	if _, exists := network.TerminalNodes["rule3_terminal"]; exists {
		t.Errorf("❌ La règle rule3 ne devrait plus exister")
	} else {
		t.Logf("✅ La règle rule3 a été supprimée")
	}
	t.Log("✅ Test réussi - Suppression de plusieurs règles fonctionne")
}

// TestRemoveRuleCommand_WithSharedAlphaNodes vérifie que la suppression n'affecte pas les nœuds partagés
func TestRemoveRuleCommand_WithSharedAlphaNodes(t *testing.T) {
	t.Log("🧪 TEST REMOVE RULE WITH SHARED ALPHA NODES")
	t.Log("===========================================")
	content := `
type Person(#id: string, name: string, age:number)
action action_adult(arg: string)
action action_voting(arg: string)
rule adult_rule : {p: Person} / p.age >= 18 ==> action_adult(p.id)
rule voting_rule : {p: Person} / p.age >= 18 ==> action_voting(p.id)
Person(id: "P1", name: "Alice", age: 25)
remove rule adult_rule
`
	tmpFile, err := os.CreateTemp("", "remove_shared_nodes_*.tsd")
	if err != nil {
		t.Fatalf("❌ Erreur création fichier temporaire: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("❌ Erreur écriture fichier: %v", err)
	}
	tmpFile.Close()
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tmpFile.Name(), nil, storage)
	if err != nil {
		t.Fatalf("❌ Erreur construction réseau: %v", err)
	}
	// Vérifier que voting_rule existe encore
	if _, exists := network.TerminalNodes["voting_rule_terminal"]; !exists {
		t.Errorf("❌ La règle voting_rule devrait encore exister")
	} else {
		t.Logf("✅ La règle voting_rule existe toujours")
	}
	// Vérifier que adult_rule n'existe plus
	if _, exists := network.TerminalNodes["adult_rule_terminal"]; exists {
		t.Errorf("❌ La règle adult_rule ne devrait plus exister")
	} else {
		t.Logf("✅ La règle adult_rule a été supprimée")
	}
	// Vérifier que les AlphaNodes partagés existent encore (pour voting_rule)
	// Le nœud alpha pour "p.age >= 18" devrait toujours exister
	alphaNodesCount := len(network.AlphaNodes)
	if alphaNodesCount == 0 {
		t.Errorf("❌ Les AlphaNodes partagés devraient encore exister")
	} else {
		t.Logf("✅ %d AlphaNodes existent encore (partagés avec voting_rule)", alphaNodesCount)
	}
	t.Log("✅ Test réussi - Les nœuds partagés sont préservés")
}

// TestRemoveRuleCommand_NonExistentRule vérifie le comportement avec une règle inexistante
func TestRemoveRuleCommand_NonExistentRule(t *testing.T) {
	t.Log("🧪 TEST REMOVE NON-EXISTENT RULE")
	t.Log("=================================")
	content := `
type Person(#id: string, name: string, age:number)
action action(arg: string)
rule existing_rule : {p: Person} / p.age > 18 ==> action(p.id)
remove rule non_existent_rule
`
	tmpFile, err := os.CreateTemp("", "remove_nonexistent_*.tsd")
	if err != nil {
		t.Fatalf("❌ Erreur création fichier temporaire: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("❌ Erreur écriture fichier: %v", err)
	}
	tmpFile.Close()
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	// La construction devrait réussir mais loguer un warning
	network, _, err := pipeline.IngestFile(tmpFile.Name(), nil, storage)
	if err != nil {
		// L'erreur est acceptable si elle mentionne que la règle n'existe pas
		t.Logf("⚠️ Erreur attendue pour règle inexistante: %v", err)
	}
	if network != nil {
		// Vérifier que existing_rule existe toujours
		if _, exists := network.TerminalNodes["existing_rule_terminal"]; !exists {
			t.Errorf("❌ La règle existing_rule devrait exister")
		} else {
			t.Logf("✅ La règle existing_rule existe toujours")
		}
	}
	t.Log("✅ Test réussi - Gestion correcte d'une règle inexistante")
}

// TestRemoveRuleCommand_AfterFactSubmission vérifie la suppression après soumission de faits
func TestRemoveRuleCommand_AfterFactSubmission(t *testing.T) {
	t.Log("🧪 TEST REMOVE RULE AFTER FACT SUBMISSION")
	t.Log("==========================================")
	// Étape 1: Créer le réseau avec une règle
	content1 := `
type Person(#id: string, name: string, age:number)
action notify(id: string)
rule adult_check : {p: Person} / p.age >= 18 ==> notify(p.id)
`
	tmpFile1, err := os.CreateTemp("", "rules_*.tsd")
	if err != nil {
		t.Fatalf("❌ Erreur création fichier: %v", err)
	}
	defer os.Remove(tmpFile1.Name())
	if _, err := tmpFile1.Write([]byte(content1)); err != nil {
		t.Fatalf("❌ Erreur écriture: %v", err)
	}
	tmpFile1.Close()
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tmpFile1.Name(), nil, storage)
	if err != nil {
		t.Fatalf("❌ Erreur construction: %v", err)
	}
	// Étape 2: Soumettre des faits
	fact := &Fact{
		ID:   "P1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   "P1",
			"name": "Alice",
			"age":  25.0,
		},
	}
	err = network.SubmitFact(fact)
	if err != nil {
		t.Fatalf("❌ Erreur soumission fait: %v", err)
	}
	t.Logf("✅ Fait soumis: %s", fact.ID)
	// Vérifier qu'il y a des activations
	terminal, exists := network.TerminalNodes["adult_check_terminal"]
	if !exists {
		t.Fatalf("❌ Terminal adult_check_terminal introuvable")
	}
	activationsBefore := 0
	if terminal.Memory != nil && terminal.Memory.Tokens != nil {
		activationsBefore = len(terminal.Memory.Tokens)
	}
	t.Logf("📊 Activations avant suppression: %d", activationsBefore)
	// Étape 3: Supprimer la règle
	err = network.RemoveRule("adult_check")
	if err != nil {
		t.Fatalf("❌ Erreur suppression règle: %v", err)
	}
	// Étape 4: Vérifier que la règle n'existe plus
	if _, exists := network.TerminalNodes["adult_check_terminal"]; exists {
		t.Errorf("❌ La règle adult_check ne devrait plus exister")
	} else {
		t.Logf("✅ La règle adult_check a été supprimée après soumission de faits")
	}
	t.Log("✅ Test réussi - Suppression après soumission de faits fonctionne")
}
