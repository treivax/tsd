// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNetworkLifecycle_RemoveSimpleRule teste la suppression d'une règle simple
func TestNetworkLifecycle_RemoveSimpleRule(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	content := `type Person(id: string, age: number, name:string)
action adult_detected(id: string)
action not_retired(arg: string)
rule r1 : {p: Person} / p.age > 18 ==> adult_detected(p.id)
rule r2 : {p: Person} / p.age < 65 ==> not_retired(p.id)
`
	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// Vérifier l'état initial
	if len(network.AlphaNodes) != 2 {
		t.Errorf("Attendu 2 AlphaNodes, obtenu %d", len(network.AlphaNodes))
	}
	if len(network.TerminalNodes) != 2 {
		t.Errorf("Attendu 2 TerminalNodes, obtenu %d", len(network.TerminalNodes))
	}
	// Vérifier le lifecycle manager
	stats := network.GetNetworkStats()
	t.Logf("Stats avant suppression: %+v", stats)
	// Supprimer r1
	err = network.RemoveRule("r1")
	if err != nil {
		t.Fatalf("Erreur suppression règle: %v", err)
	}
	// Vérifier qu'un AlphaNode et un TerminalNode ont été supprimés
	if len(network.AlphaNodes) != 1 {
		t.Errorf("Après suppression, attendu 1 AlphaNode, obtenu %d", len(network.AlphaNodes))
	}
	if len(network.TerminalNodes) != 1 {
		t.Errorf("Après suppression, attendu 1 TerminalNode, obtenu %d", len(network.TerminalNodes))
	}
	// Le TypeNode doit être conservé (utilisé par rule_1)
	if len(network.TypeNodes) != 1 {
		t.Errorf("Le TypeNode devrait être conservé, obtenu %d", len(network.TypeNodes))
	}
	// Vérifier les stats après suppression
	stats = network.GetNetworkStats()
	t.Logf("Stats après suppression: %+v", stats)
}

// TestNetworkLifecycle_RemoveAllRulesForType teste la suppression de toutes les règles d'un type
func TestNetworkLifecycle_RemoveAllRulesForType(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	content := `type Person(id: string, age:number)
action adult(arg: string)
action young(arg: string)
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
rule r2 : {p: Person} / p.age < 65 ==> young(p.id)
`
	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// Supprimer les deux règles
	err = network.RemoveRule("r1")
	if err != nil {
		t.Fatalf("Erreur suppression r1: %v", err)
	}
	err = network.RemoveRule("r2")
	if err != nil {
		t.Fatalf("Erreur suppression r2: %v", err)
	}
	// Tous les nœuds sauf le TypeNode devraient être supprimés
	if len(network.AlphaNodes) != 0 {
		t.Errorf("Après suppression des règles, attendu 0 AlphaNode, obtenu %d", len(network.AlphaNodes))
	}
	if len(network.TerminalNodes) != 0 {
		t.Errorf("Après suppression des règles, attendu 0 TerminalNode, obtenu %d", len(network.TerminalNodes))
	}
	// Note: Le TypeNode reste car il est créé au niveau du type, pas de la règle
	// C'est un comportement intentionnel pour permettre l'ajout dynamique de règles
}

// TestNetworkLifecycle_SharedNodeNotRemoved teste qu'un nœud partagé n'est pas supprimé
func TestNetworkLifecycle_SharedNodeNotRemoved(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	content := `type Person(id: string, age:number)
action adult(arg: string)
action young(arg: string)
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
rule r2 : {p: Person} / p.age < 65 ==> young(p.id)
`
	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// Le TypeNode est partagé entre les deux règles
	personTypeNode := network.TypeNodes["Person"]
	if personTypeNode == nil {
		t.Fatal("TypeNode Person non trouvé")
	}
	// Vérifier le lifecycle du TypeNode
	if network.LifecycleManager != nil {
		lifecycle, exists := network.LifecycleManager.GetNodeLifecycle(personTypeNode.GetID())
		if exists {
			t.Logf("TypeNode a %d référence(s) initiale(s)", lifecycle.GetRefCount())
		}
	}
	// Supprimer r1
	err = network.RemoveRule("r1")
	if err != nil {
		t.Fatalf("Erreur suppression r1: %v", err)
	}
	// Le TypeNode doit toujours exister (partagé avec rule_1)
	if network.TypeNodes["Person"] == nil {
		t.Error("Le TypeNode Person ne devrait pas être supprimé (partagé)")
	}
}

// TestNetworkLifecycle_GetRuleInfo teste la récupération d'informations sur une règle
func TestNetworkLifecycle_GetRuleInfo(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	content := `type Person(id: string, age:number)
action adult(arg: string)
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
`
	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// Récupérer les informations de la règle
	info := network.GetRuleInfo("r1")
	if info == nil {
		t.Fatal("GetRuleInfo ne devrait pas retourner nil")
	}
	t.Logf("Informations de r1:")
	t.Logf("  RuleID: %s", info.RuleID)
	t.Logf("  RuleName: %s", info.RuleName)
	t.Logf("  NodeCount: %d", info.NodeCount)
	t.Logf("  NodeIDs: %v", info.NodeIDs)
	// La règle devrait avoir au moins 2 nœuds (AlphaNode + TerminalNode)
	if info.NodeCount < 2 {
		t.Errorf("Attendu au moins 2 nœuds pour la règle, obtenu %d", info.NodeCount)
	}
}

// TestNetworkLifecycle_GetNetworkStats teste les statistiques du réseau
func TestNetworkLifecycle_GetNetworkStats(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	content := `type Person(id: string, age:number)
type Company(id: string, revenue:number)
action adult(arg: string)
action big_company(arg: string)
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
rule r2 : {c: Company} / c.revenue > 1000000 ==> big_company(c.id)
`
	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	stats := network.GetNetworkStats()
	t.Log("Statistiques du réseau:")
	for key, value := range stats {
		t.Logf("  %s: %v", key, value)
	}
	// Vérifications basiques
	if stats["type_nodes"].(int) != 2 {
		t.Errorf("Attendu 2 TypeNodes, obtenu %d", stats["type_nodes"].(int))
	}
	if stats["alpha_nodes"].(int) != 2 {
		t.Errorf("Attendu 2 AlphaNodes, obtenu %d", stats["alpha_nodes"].(int))
	}
	if stats["terminal_nodes"].(int) != 2 {
		t.Errorf("Attendu 2 TerminalNodes, obtenu %d", stats["terminal_nodes"].(int))
	}
	// Vérifier les stats du lifecycle
	if lifecycle, ok := stats["lifecycle_total_nodes"]; ok {
		t.Logf("Total de nœuds trackés par le lifecycle: %v", lifecycle)
	}
}

// TestNetworkLifecycle_RemoveNonExistentRule teste la suppression d'une règle inexistante
func TestNetworkLifecycle_RemoveNonExistentRule(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Tenter de supprimer une règle qui n'existe pas
	err := network.RemoveRule("rule_999")
	if err == nil {
		t.Error("La suppression d'une règle inexistante devrait échouer")
	}
	t.Logf("Erreur attendue: %v", err)
}

// TestNetworkLifecycle_ResetClearsLifecycle teste que Reset nettoie le lifecycle
func TestNetworkLifecycle_ResetClearsLifecycle(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	content := `type Person(id: string, age:number)
action adult(arg: string)
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
`
	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// Vérifier que le lifecycle contient des nœuds
	statsBefore := network.GetNetworkStats()
	if lifecycle, ok := statsBefore["lifecycle_total_nodes"]; ok {
		if lifecycle.(int) == 0 {
			t.Error("Le lifecycle devrait contenir des nœuds avant Reset")
		}
	}
	// Réinitialiser le réseau
	network.Reset()
	// Vérifier que le lifecycle a été nettoyé
	statsAfter := network.GetNetworkStats()
	if lifecycle, ok := statsAfter["lifecycle_total_nodes"]; ok {
		if lifecycle.(int) != 0 {
			t.Errorf("Le lifecycle devrait être vide après Reset, contient %d nœuds", lifecycle.(int))
		}
	}
}

// TestNetworkLifecycle_MultipleRulesOnSameType teste la suppression partielle de règles
func TestNetworkLifecycle_MultipleRulesOnSameType(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	content := `type Person(id: string, age: number, salary:number)
action adult(arg: string)
action high_earner(arg: string)
action young(arg: string)
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
rule r2 : {p: Person} / p.age < 65 ==> young(p.id)
rule r3 : {p: Person} / p.salary > 50000 ==> high_earner(p.id)
`
	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// État initial: 3 AlphaNodes, 3 TerminalNodes
	initialAlpha := len(network.AlphaNodes)
	initialTerminal := len(network.TerminalNodes)
	t.Logf("État initial: %d AlphaNodes, %d TerminalNodes", initialAlpha, initialTerminal)
	if initialAlpha != 3 {
		t.Errorf("Attendu 3 AlphaNodes, obtenu %d", initialAlpha)
	}
	// Supprimer r2 (au milieu)
	err = network.RemoveRule("r2")
	if err != nil {
		t.Fatalf("Erreur suppression r2: %v", err)
	}
	// Vérifier qu'un nœud a été supprimé
	if len(network.AlphaNodes) != 2 {
		t.Errorf("Après suppression r2, attendu 2 AlphaNodes, obtenu %d", len(network.AlphaNodes))
	}
	if len(network.TerminalNodes) != 2 {
		t.Errorf("Après suppression r2, attendu 2 TerminalNodes, obtenu %d", len(network.TerminalNodes))
	}
	// Le TypeNode doit toujours exister
	if len(network.TypeNodes) != 1 {
		t.Errorf("Le TypeNode devrait toujours exister, obtenu %d", len(network.TypeNodes))
	}
	// Supprimer r2
	// Supprimer r1
	err = network.RemoveRule("r1")
	if err != nil {
		t.Fatalf("Erreur suppression r1: %v", err)
	}
	if len(network.AlphaNodes) != 1 {
		t.Errorf("Après suppression rule_0, attendu 1 AlphaNode, obtenu %d", len(network.AlphaNodes))
	}
	// Supprimer rule_2 (la dernière)
	// Supprimer r3
	err = network.RemoveRule("r3")
	if err != nil {
		t.Fatalf("Erreur suppression r3: %v", err)
	}
	if len(network.AlphaNodes) != 0 {
		t.Errorf("Après suppression de toutes les règles, attendu 0 AlphaNode, obtenu %d", len(network.AlphaNodes))
	}
	if len(network.TerminalNodes) != 0 {
		t.Errorf("Après suppression de toutes les règles, attendu 0 TerminalNode, obtenu %d", len(network.TerminalNodes))
	}
	t.Log("✅ Toutes les règles ont été supprimées avec succès")
}
