// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"os"
	"path/filepath"
	"testing"
)

// TestAlphaSharingIntegration_TwoRulesSameCondition teste le partage d'AlphaNodes
// entre deux règles ayant la même condition
func TestAlphaSharingIntegration_TwoRulesSameCondition(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	// Deux règles avec la même condition: p.age > 18
	content := `type Person(#id: string, age:number)
action print(message: string)
rule adult_check : {p: Person} / p.age > 18 ==> print("Adult")
rule voting_check : {p: Person} / p.age > 18 ==> print("Can vote")
`
	if err := os.WriteFile(constraintFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur écriture fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// Vérifier les statistiques de partage
	stats := network.GetNetworkStats()
	// Devrait avoir 1 AlphaNode partagé (même condition: p.age > 18)
	totalAlphaNodes := stats["alpha_nodes"].(int)
	if totalAlphaNodes != 1 {
		t.Errorf("Devrait avoir 1 AlphaNode partagé, got %d", totalAlphaNodes)
	}
	// Devrait avoir 2 TerminalNodes (un par règle)
	totalTerminalNodes := stats["terminal_nodes"].(int)
	if totalTerminalNodes != 2 {
		t.Errorf("Devrait avoir 2 TerminalNodes, got %d", totalTerminalNodes)
	}
	// Vérifier les stats de partage
	if sharingStats, ok := stats["sharing_total_shared_alpha_nodes"]; ok {
		sharedCount := sharingStats.(int)
		if sharedCount != 1 {
			t.Errorf("Devrait avoir 1 nœud dans le registre de partage, got %d", sharedCount)
		}
	} else {
		t.Error("Les stats de partage devraient être présentes")
	}
	if ruleRefs, ok := stats["sharing_total_rule_references"]; ok {
		refCount := ruleRefs.(int)
		if refCount != 2 {
			t.Errorf("Devrait avoir 2 références de règles, got %d", refCount)
		}
	}
	// Vérifier le ratio de partage (2 règles / 1 nœud = 2.0)
	if avgRatio, ok := stats["sharing_average_sharing_ratio"]; ok {
		ratio := avgRatio.(float64)
		if ratio != 2.0 {
			t.Errorf("Le ratio de partage devrait être 2.0, got %v", ratio)
		}
	}
}

// TestAlphaSharingIntegration_ThreeRulesMixedConditions teste le partage
// avec des conditions mixtes (certaines identiques, d'autres différentes)
func TestAlphaSharingIntegration_ThreeRulesMixedConditions(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	// Trois règles: deux avec p.age > 18, une avec p.age > 21
	content := `type Person(#id: string, age:number)
action print(message: string)
rule adult_check : {p: Person} / p.age > 18 ==> print("Adult")
rule voting_check : {p: Person} / p.age > 18 ==> print("Can vote")
rule drinking_check : {p: Person} / p.age > 21 ==> print("Can drink")
`
	if err := os.WriteFile(constraintFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur écriture fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	stats := network.GetNetworkStats()
	// Devrait avoir 2 AlphaNodes (un pour >18, un pour >21)
	totalAlphaNodes := stats["alpha_nodes"].(int)
	if totalAlphaNodes != 2 {
		t.Errorf("Devrait avoir 2 AlphaNodes, got %d", totalAlphaNodes)
	}
	// Devrait avoir 3 TerminalNodes (un par règle)
	totalTerminalNodes := stats["terminal_nodes"].(int)
	if totalTerminalNodes != 3 {
		t.Errorf("Devrait avoir 3 TerminalNodes, got %d", totalTerminalNodes)
	}
	// Vérifier les stats de partage
	if sharingStats, ok := stats["sharing_total_shared_alpha_nodes"]; ok {
		sharedCount := sharingStats.(int)
		if sharedCount != 2 {
			t.Errorf("Devrait avoir 2 nœuds dans le registre de partage, got %d", sharedCount)
		}
	}
	if ruleRefs, ok := stats["sharing_total_rule_references"]; ok {
		refCount := ruleRefs.(int)
		if refCount != 3 {
			t.Errorf("Devrait avoir 3 références de règles, got %d", refCount)
		}
	}
	// Vérifier le ratio de partage (3 règles / 2 nœuds = 1.5)
	if avgRatio, ok := stats["sharing_average_sharing_ratio"]; ok {
		ratio := avgRatio.(float64)
		if ratio != 1.5 {
			t.Errorf("Le ratio de partage devrait être 1.5, got %v", ratio)
		}
	}
}

// TestAlphaSharingIntegration_FactPropagation teste que les faits se propagent
// correctement à travers les AlphaNodes partagés vers tous les TerminalNodes
func TestAlphaSharingIntegration_FactPropagation(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	content := `type Person(#id: string, age:number)
action print(message: string)
rule adult_check : {p: Person} / p.age > 18 ==> print("Adult")
rule voting_check : {p: Person} / p.age > 18 ==> print("Can vote")
`
	if err := os.WriteFile(constraintFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur écriture fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// Soumettre un fait qui satisfait la condition
	fact := &Fact{
		ID:   "person1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":  "p1",
			"age": 25.0,
		},
	}
	err = network.SubmitFact(fact)
	if err != nil {
		t.Fatalf("Erreur soumission fait: %v", err)
	}
	// Vérifier que les deux TerminalNodes ont reçu le fait
	terminalCount := 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		tokens := memory.Tokens
		if len(tokens) > 0 {
			terminalCount++
		}
	}
	if terminalCount != 2 {
		t.Errorf("Les 2 TerminalNodes devraient avoir reçu le fait, got %d", terminalCount)
	}
	// Soumettre un fait qui ne satisfait pas la condition
	fact2 := &Fact{
		ID:   "person2",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":  "p2",
			"age": 15.0,
		},
	}
	err = network.SubmitFact(fact2)
	if err != nil {
		t.Fatalf("Erreur soumission fait2: %v", err)
	}
	// Vérifier que les TerminalNodes n'ont toujours qu'un seul token chacun
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		tokens := memory.Tokens
		if len(tokens) != 1 {
			t.Errorf("Chaque TerminalNode devrait avoir exactement 1 token, got %d", len(tokens))
		}
	}
}

// TestAlphaSharingIntegration_RuleRemoval teste la suppression de règles
// avec gestion correcte des AlphaNodes partagés
func TestAlphaSharingIntegration_RuleRemoval(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	content := `type Person(#id: string, age:number)
action print(message: string)
rule adult_check : {p: Person} / p.age > 18 ==> print("Adult")
rule voting_check : {p: Person} / p.age > 18 ==> print("Can vote")
rule senior_check : {p: Person} / p.age > 65 ==> print("Senior")
`
	if err := os.WriteFile(constraintFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur écriture fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// État initial: 2 AlphaNodes partagés, 3 TerminalNodes
	stats := network.GetNetworkStats()
	initialAlphaCount := stats["alpha_nodes"].(int)
	if initialAlphaCount != 2 {
		t.Errorf("Devrait avoir 2 AlphaNodes initialement, got %d", initialAlphaCount)
	}
	// Supprimer adult_check (partage l'AlphaNode avec voting_check)
	err = network.RemoveRule("adult_check")
	if err != nil {
		t.Fatalf("Erreur suppression règle adult_check: %v", err)
	}
	// L'AlphaNode pour >18 devrait toujours exister (utilisé par voting_check)
	stats = network.GetNetworkStats()
	afterRemovalAlphaCount := stats["alpha_nodes"].(int)
	if afterRemovalAlphaCount != 2 {
		t.Errorf("Devrait toujours avoir 2 AlphaNodes après suppression de adult_check, got %d", afterRemovalAlphaCount)
	}
	// Devrait avoir 2 TerminalNodes restants
	afterRemovalTerminalCount := stats["terminal_nodes"].(int)
	if afterRemovalTerminalCount != 2 {
		t.Errorf("Devrait avoir 2 TerminalNodes après suppression, got %d", afterRemovalTerminalCount)
	}
	// Supprimer voting_check (dernière règle utilisant l'AlphaNode >18)
	err = network.RemoveRule("voting_check")
	if err != nil {
		t.Fatalf("Erreur suppression règle voting_check: %v", err)
	}
	// L'AlphaNode pour >18 devrait maintenant être supprimé
	stats = network.GetNetworkStats()
	finalAlphaCount := stats["alpha_nodes"].(int)
	if finalAlphaCount != 1 {
		t.Errorf("Devrait avoir 1 AlphaNode après suppression des deux règles partageant >18, got %d", finalAlphaCount)
	}
	// Devrait avoir 1 TerminalNode restant (senior_check)
	finalTerminalCount := stats["terminal_nodes"].(int)
	if finalTerminalCount != 1 {
		t.Errorf("Devrait avoir 1 TerminalNode final, got %d", finalTerminalCount)
	}
}

// TestAlphaSharingIntegration_DifferentTypes teste que les conditions
// sur des types différents ne sont pas partagées (variables différentes)
func TestAlphaSharingIntegration_DifferentTypes(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	// Même attribut et opérateur, mais types et variables différents
	content := `type Person(#id: string, age:number)
type Animal(#id: string, age:number)
action print(message: string)
rule person_adult : {p: Person} / p.age > 18 ==> print("Adult person")
rule animal_old : {a: Animal} / a.age > 18 ==> print("Old animal")
`
	if err := os.WriteFile(constraintFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur écriture fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	stats := network.GetNetworkStats()
	// Devrait avoir 2 TypeNodes (Person et Animal)
	totalTypeNodes := stats["type_nodes"].(int)
	if totalTypeNodes != 2 {
		t.Errorf("Devrait avoir 2 TypeNodes, got %d", totalTypeNodes)
	}
	// Les conditions sont identiques mais avec des variables différentes (p vs a)
	// Donc elles ne partagent pas le même AlphaNode
	totalAlphaNodes := stats["alpha_nodes"].(int)
	if totalAlphaNodes != 2 {
		t.Errorf("Devrait avoir 2 AlphaNodes (variables différentes: p vs a), got %d", totalAlphaNodes)
	}
}

// TestAlphaSharingIntegration_NetworkReset teste que le reset
// nettoie correctement le registre de partage
func TestAlphaSharingIntegration_NetworkReset(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	content := `type Person(#id: string, age:number)
action print(message: string)
rule adult_check : {p: Person} / p.age > 18 ==> print("Adult")
rule voting_check : {p: Person} / p.age > 18 ==> print("Can vote")
`
	if err := os.WriteFile(constraintFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur écriture fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// Vérifier l'état initial
	stats := network.GetNetworkStats()
	if stats["alpha_nodes"].(int) != 1 {
		t.Error("Devrait avoir 1 AlphaNode avant reset")
	}
	if stats["sharing_total_shared_alpha_nodes"].(int) != 1 {
		t.Error("Devrait avoir 1 nœud partagé avant reset")
	}
	// Reset du réseau
	network.Reset()
	// Vérifier que tout est nettoyé
	stats = network.GetNetworkStats()
	if stats["alpha_nodes"].(int) != 0 {
		t.Error("Devrait avoir 0 AlphaNodes après reset")
	}
	if stats["terminal_nodes"].(int) != 0 {
		t.Error("Devrait avoir 0 TerminalNodes après reset")
	}
	if stats["sharing_total_shared_alpha_nodes"].(int) != 0 {
		t.Error("Le registre de partage devrait être vide après reset")
	}
	if stats["sharing_total_rule_references"].(int) != 0 {
		t.Error("Devrait avoir 0 références de règles après reset")
	}
}

// TestAlphaSharingIntegration_ComplexConditions teste le partage
// avec des conditions plus complexes
func TestAlphaSharingIntegration_ComplexConditions(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	// Conditions identiques mais actions différentes
	content := `type Person(#id: string, age: number, salary:number)
action alert(arg: string)
action print(message: string)
rule high_earner1 : {p: Person} / p.salary > 100000 ==> print("High earner")
rule high_earner2 : {p: Person} / p.salary > 100000 ==> alert("Very high salary")
rule mid_earner : {p: Person} / p.salary > 50000 ==> print("Mid earner")
`
	if err := os.WriteFile(constraintFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur écriture fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	stats := network.GetNetworkStats()
	// Devrait avoir 2 AlphaNodes (>100000 et >50000)
	totalAlphaNodes := stats["alpha_nodes"].(int)
	if totalAlphaNodes != 2 {
		t.Errorf("Devrait avoir 2 AlphaNodes, got %d", totalAlphaNodes)
	}
	// Devrait avoir 3 TerminalNodes (un par règle)
	totalTerminalNodes := stats["terminal_nodes"].(int)
	if totalTerminalNodes != 3 {
		t.Errorf("Devrait avoir 3 TerminalNodes, got %d", totalTerminalNodes)
	}
	// Le nœud pour >100000 est partagé par 2 règles
	if avgRatio, ok := stats["sharing_average_sharing_ratio"]; ok {
		ratio := avgRatio.(float64)
		if ratio != 1.5 { // 3 règles / 2 nœuds = 1.5
			t.Errorf("Le ratio de partage devrait être 1.5, got %v", ratio)
		}
	}
	// Soumettre un fait pour tester
	fact := &Fact{
		ID:   "person1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":     "p1",
			"age":    30.0,
			"salary": 120000.0,
		},
	}
	err = network.SubmitFact(fact)
	if err != nil {
		t.Fatalf("Erreur soumission fait: %v", err)
	}
	// Vérifier que les 3 règles sont activées
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		tokens := memory.Tokens
		if len(tokens) > 0 {
			activatedCount++
		}
	}
	if activatedCount != 3 {
		t.Errorf("Les 3 règles devraient être activées, got %d", activatedCount)
	}
}
