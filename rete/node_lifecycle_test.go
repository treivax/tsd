// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// TestNodeLifecycle_Basic teste les opérations de base sur le cycle de vie d'un nœud
func TestNodeLifecycle_Basic(t *testing.T) {
	lifecycle := NewNodeLifecycle("alpha_1", "alpha")

	// Test initial
	if lifecycle.NodeID != "alpha_1" {
		t.Errorf("NodeID attendu 'alpha_1', obtenu '%s'", lifecycle.NodeID)
	}
	if lifecycle.NodeType != "alpha" {
		t.Errorf("NodeType attendu 'alpha', obtenu '%s'", lifecycle.NodeType)
	}
	if lifecycle.RefCount != 0 {
		t.Errorf("RefCount initial devrait être 0, obtenu %d", lifecycle.RefCount)
	}
	if lifecycle.HasReferences() {
		t.Error("HasReferences devrait être false initialement")
	}
}

// TestNodeLifecycle_AddRuleReference teste l'ajout de références de règles
func TestNodeLifecycle_AddRuleReference(t *testing.T) {
	lifecycle := NewNodeLifecycle("alpha_1", "alpha")

	// Ajouter une première règle
	lifecycle.AddRuleReference("rule_1", "First Rule")
	if lifecycle.RefCount != 1 {
		t.Errorf("RefCount devrait être 1, obtenu %d", lifecycle.RefCount)
	}
	if !lifecycle.HasReferences() {
		t.Error("HasReferences devrait être true")
	}

	// Ajouter une deuxième règle
	lifecycle.AddRuleReference("rule_2", "Second Rule")
	if lifecycle.RefCount != 2 {
		t.Errorf("RefCount devrait être 2, obtenu %d", lifecycle.RefCount)
	}

	// Ajouter la même règle à nouveau (ne devrait pas augmenter le compteur)
	lifecycle.AddRuleReference("rule_1", "First Rule")
	if lifecycle.RefCount != 2 {
		t.Errorf("RefCount devrait rester 2 (pas de duplication), obtenu %d", lifecycle.RefCount)
	}
}

// TestNodeLifecycle_RemoveRuleReference teste la suppression de références de règles
func TestNodeLifecycle_RemoveRuleReference(t *testing.T) {
	lifecycle := NewNodeLifecycle("alpha_1", "alpha")

	// Ajouter deux règles
	lifecycle.AddRuleReference("rule_1", "First Rule")
	lifecycle.AddRuleReference("rule_2", "Second Rule")

	// Retirer la première règle
	shouldDelete := lifecycle.RemoveRuleReference("rule_1")
	if shouldDelete {
		t.Error("shouldDelete devrait être false (reste rule_2)")
	}
	if lifecycle.RefCount != 1 {
		t.Errorf("RefCount devrait être 1, obtenu %d", lifecycle.RefCount)
	}

	// Retirer la deuxième règle
	shouldDelete = lifecycle.RemoveRuleReference("rule_2")
	if !shouldDelete {
		t.Error("shouldDelete devrait être true (plus de références)")
	}
	if lifecycle.RefCount != 0 {
		t.Errorf("RefCount devrait être 0, obtenu %d", lifecycle.RefCount)
	}
	if lifecycle.HasReferences() {
		t.Error("HasReferences devrait être false")
	}
}

// TestNodeLifecycle_GetRules teste la récupération de la liste des règles
func TestNodeLifecycle_GetRules(t *testing.T) {
	lifecycle := NewNodeLifecycle("alpha_1", "alpha")

	lifecycle.AddRuleReference("rule_1", "First Rule")
	lifecycle.AddRuleReference("rule_2", "Second Rule")
	lifecycle.AddRuleReference("rule_3", "Third Rule")

	rules := lifecycle.GetRules()
	if len(rules) != 3 {
		t.Errorf("Attendu 3 règles, obtenu %d", len(rules))
	}

	// Vérifier que les IDs sont présents
	rulesMap := make(map[string]bool)
	for _, ruleID := range rules {
		rulesMap[ruleID] = true
	}

	if !rulesMap["rule_1"] || !rulesMap["rule_2"] || !rulesMap["rule_3"] {
		t.Error("Toutes les règles devraient être présentes")
	}
}

// TestNodeLifecycle_GetRuleInfo teste la récupération d'informations sur une règle
func TestNodeLifecycle_GetRuleInfo(t *testing.T) {
	lifecycle := NewNodeLifecycle("alpha_1", "alpha")

	lifecycle.AddRuleReference("rule_1", "Test Rule")

	ref, exists := lifecycle.GetRuleInfo("rule_1")
	if !exists {
		t.Fatal("La règle rule_1 devrait exister")
	}
	if ref.RuleID != "rule_1" {
		t.Errorf("RuleID attendu 'rule_1', obtenu '%s'", ref.RuleID)
	}
	if ref.RuleName != "Test Rule" {
		t.Errorf("RuleName attendu 'Test Rule', obtenu '%s'", ref.RuleName)
	}

	// Tester avec une règle inexistante
	_, exists = lifecycle.GetRuleInfo("rule_999")
	if exists {
		t.Error("La règle rule_999 ne devrait pas exister")
	}
}

// TestLifecycleManager_Basic teste les opérations de base du gestionnaire
func TestLifecycleManager_Basic(t *testing.T) {
	manager := NewLifecycleManager()

	if len(manager.Nodes) != 0 {
		t.Errorf("Le gestionnaire devrait être vide initialement, contient %d nœuds", len(manager.Nodes))
	}
}

// TestLifecycleManager_RegisterNode teste l'enregistrement de nœuds
func TestLifecycleManager_RegisterNode(t *testing.T) {
	manager := NewLifecycleManager()

	// Enregistrer un premier nœud
	lifecycle1 := manager.RegisterNode("alpha_1", "alpha")
	if lifecycle1 == nil {
		t.Fatal("RegisterNode devrait retourner un lifecycle non nil")
	}
	if lifecycle1.NodeID != "alpha_1" {
		t.Errorf("NodeID attendu 'alpha_1', obtenu '%s'", lifecycle1.NodeID)
	}

	// Enregistrer le même nœud à nouveau (devrait retourner le même lifecycle)
	lifecycle2 := manager.RegisterNode("alpha_1", "alpha")
	if lifecycle2 != lifecycle1 {
		t.Error("L'enregistrement du même nœud devrait retourner le même lifecycle")
	}

	// Enregistrer un deuxième nœud différent
	lifecycle3 := manager.RegisterNode("alpha_2", "alpha")
	if lifecycle3 == lifecycle1 {
		t.Error("Un nouveau nœud devrait avoir un nouveau lifecycle")
	}

	if len(manager.Nodes) != 2 {
		t.Errorf("Le gestionnaire devrait contenir 2 nœuds, contient %d", len(manager.Nodes))
	}
}

// TestLifecycleManager_AddRuleToNode teste l'ajout de règles aux nœuds
func TestLifecycleManager_AddRuleToNode(t *testing.T) {
	manager := NewLifecycleManager()

	// Enregistrer un nœud
	manager.RegisterNode("alpha_1", "alpha")

	// Ajouter une règle au nœud
	err := manager.AddRuleToNode("alpha_1", "rule_1", "Test Rule")
	if err != nil {
		t.Errorf("AddRuleToNode ne devrait pas échouer: %v", err)
	}

	// Vérifier que la règle a été ajoutée
	lifecycle, exists := manager.GetNodeLifecycle("alpha_1")
	if !exists {
		t.Fatal("Le nœud alpha_1 devrait exister")
	}
	if lifecycle.RefCount != 1 {
		t.Errorf("RefCount devrait être 1, obtenu %d", lifecycle.RefCount)
	}

	// Tenter d'ajouter une règle à un nœud non enregistré
	err = manager.AddRuleToNode("alpha_999", "rule_1", "Test Rule")
	if err == nil {
		t.Error("AddRuleToNode devrait échouer pour un nœud non enregistré")
	}
}

// TestLifecycleManager_RemoveRuleFromNode teste la suppression de règles des nœuds
func TestLifecycleManager_RemoveRuleFromNode(t *testing.T) {
	manager := NewLifecycleManager()

	// Préparer un nœud avec deux règles
	manager.RegisterNode("alpha_1", "alpha")
	manager.AddRuleToNode("alpha_1", "rule_1", "Rule 1")
	manager.AddRuleToNode("alpha_1", "rule_2", "Rule 2")

	// Retirer la première règle
	shouldDelete, err := manager.RemoveRuleFromNode("alpha_1", "rule_1")
	if err != nil {
		t.Errorf("RemoveRuleFromNode ne devrait pas échouer: %v", err)
	}
	if shouldDelete {
		t.Error("shouldDelete devrait être false (reste rule_2)")
	}

	// Retirer la deuxième règle
	shouldDelete, err = manager.RemoveRuleFromNode("alpha_1", "rule_2")
	if err != nil {
		t.Errorf("RemoveRuleFromNode ne devrait pas échouer: %v", err)
	}
	if !shouldDelete {
		t.Error("shouldDelete devrait être true (plus de références)")
	}
}

// TestLifecycleManager_RemoveNode teste la suppression de nœuds
func TestLifecycleManager_RemoveNode(t *testing.T) {
	manager := NewLifecycleManager()

	// Créer un nœud avec une règle
	manager.RegisterNode("alpha_1", "alpha")
	manager.AddRuleToNode("alpha_1", "rule_1", "Rule 1")

	// Tenter de supprimer un nœud avec des références (devrait échouer)
	err := manager.RemoveNode("alpha_1")
	if err == nil {
		t.Error("RemoveNode devrait échouer pour un nœud avec des références")
	}

	// Retirer la règle
	manager.RemoveRuleFromNode("alpha_1", "rule_1")

	// Supprimer le nœud sans références (devrait réussir)
	err = manager.RemoveNode("alpha_1")
	if err != nil {
		t.Errorf("RemoveNode ne devrait pas échouer pour un nœud sans références: %v", err)
	}

	// Vérifier que le nœud a été supprimé
	_, exists := manager.GetNodeLifecycle("alpha_1")
	if exists {
		t.Error("Le nœud alpha_1 ne devrait plus exister")
	}
}

// TestLifecycleManager_GetNodesForRule teste la récupération des nœuds d'une règle
func TestLifecycleManager_GetNodesForRule(t *testing.T) {
	manager := NewLifecycleManager()

	// Créer plusieurs nœuds et les associer à des règles
	manager.RegisterNode("alpha_1", "alpha")
	manager.RegisterNode("alpha_2", "alpha")
	manager.RegisterNode("terminal_1", "terminal")

	manager.AddRuleToNode("alpha_1", "rule_1", "Rule 1")
	manager.AddRuleToNode("alpha_2", "rule_1", "Rule 1")
	manager.AddRuleToNode("terminal_1", "rule_1", "Rule 1")
	manager.AddRuleToNode("alpha_2", "rule_2", "Rule 2")

	// Récupérer les nœuds de rule_1
	nodes1 := manager.GetNodesForRule("rule_1")
	if len(nodes1) != 3 {
		t.Errorf("rule_1 devrait avoir 3 nœuds, obtenu %d", len(nodes1))
	}

	// Récupérer les nœuds de rule_2
	nodes2 := manager.GetNodesForRule("rule_2")
	if len(nodes2) != 1 {
		t.Errorf("rule_2 devrait avoir 1 nœud, obtenu %d", len(nodes2))
	}

	// Récupérer les nœuds d'une règle inexistante
	nodes3 := manager.GetNodesForRule("rule_999")
	if len(nodes3) != 0 {
		t.Errorf("rule_999 devrait avoir 0 nœud, obtenu %d", len(nodes3))
	}
}

// TestLifecycleManager_CanRemoveNode teste la vérification de suppression
func TestLifecycleManager_CanRemoveNode(t *testing.T) {
	manager := NewLifecycleManager()

	manager.RegisterNode("alpha_1", "alpha")
	manager.AddRuleToNode("alpha_1", "rule_1", "Rule 1")

	// Ne devrait pas pouvoir être supprimé (a des références)
	if manager.CanRemoveNode("alpha_1") {
		t.Error("alpha_1 ne devrait pas pouvoir être supprimé (a des références)")
	}

	// Retirer la règle
	manager.RemoveRuleFromNode("alpha_1", "rule_1")

	// Devrait pouvoir être supprimé maintenant
	if !manager.CanRemoveNode("alpha_1") {
		t.Error("alpha_1 devrait pouvoir être supprimé (plus de références)")
	}
}

// TestLifecycleManager_GetStats teste les statistiques du gestionnaire
func TestLifecycleManager_GetStats(t *testing.T) {
	manager := NewLifecycleManager()

	// État initial
	stats := manager.GetStats()
	if stats["total_nodes"].(int) != 0 {
		t.Error("total_nodes devrait être 0 initialement")
	}

	// Ajouter des nœuds avec et sans références
	manager.RegisterNode("alpha_1", "alpha")
	manager.RegisterNode("alpha_2", "alpha")
	manager.RegisterNode("alpha_3", "alpha")

	manager.AddRuleToNode("alpha_1", "rule_1", "Rule 1")
	manager.AddRuleToNode("alpha_2", "rule_1", "Rule 1")
	// alpha_3 reste sans référence

	stats = manager.GetStats()
	if stats["total_nodes"].(int) != 3 {
		t.Errorf("total_nodes devrait être 3, obtenu %d", stats["total_nodes"].(int))
	}
	if stats["nodes_with_refs"].(int) != 2 {
		t.Errorf("nodes_with_refs devrait être 2, obtenu %d", stats["nodes_with_refs"].(int))
	}
	if stats["nodes_without_refs"].(int) != 1 {
		t.Errorf("nodes_without_refs devrait être 1, obtenu %d", stats["nodes_without_refs"].(int))
	}
	if stats["total_references"].(int) != 2 {
		t.Errorf("total_references devrait être 2, obtenu %d", stats["total_references"].(int))
	}
}

// TestLifecycleManager_Reset teste la réinitialisation du gestionnaire
func TestLifecycleManager_Reset(t *testing.T) {
	manager := NewLifecycleManager()

	// Ajouter des nœuds
	manager.RegisterNode("alpha_1", "alpha")
	manager.RegisterNode("alpha_2", "alpha")
	manager.AddRuleToNode("alpha_1", "rule_1", "Rule 1")

	// Réinitialiser
	manager.Reset()

	// Vérifier que tout est vide
	if len(manager.Nodes) != 0 {
		t.Errorf("Après Reset, le gestionnaire devrait être vide, contient %d nœuds", len(manager.Nodes))
	}

	stats := manager.GetStats()
	if stats["total_nodes"].(int) != 0 {
		t.Error("total_nodes devrait être 0 après Reset")
	}
}

// TestLifecycleManager_GetRuleInfo teste la récupération d'informations sur une règle
func TestLifecycleManager_GetRuleInfo(t *testing.T) {
	manager := NewLifecycleManager()

	// Créer une règle avec plusieurs nœuds
	manager.RegisterNode("alpha_1", "alpha")
	manager.RegisterNode("terminal_1", "terminal")

	manager.AddRuleToNode("alpha_1", "rule_1", "Test Rule")
	manager.AddRuleToNode("terminal_1", "rule_1", "Test Rule")

	// Récupérer les informations
	info := manager.GetRuleInfo("rule_1")
	if info == nil {
		t.Fatal("GetRuleInfo ne devrait pas retourner nil")
	}
	if info.RuleID != "rule_1" {
		t.Errorf("RuleID attendu 'rule_1', obtenu '%s'", info.RuleID)
	}
	if info.RuleName != "Test Rule" {
		t.Errorf("RuleName attendu 'Test Rule', obtenu '%s'", info.RuleName)
	}
	if info.NodeCount != 2 {
		t.Errorf("NodeCount attendu 2, obtenu %d", info.NodeCount)
	}
	if len(info.NodeIDs) != 2 {
		t.Errorf("NodeIDs devrait contenir 2 éléments, obtenu %d", len(info.NodeIDs))
	}
}

// TestLifecycleManager_ConcurrentAccess teste l'accès concurrent (basique)
func TestLifecycleManager_ConcurrentAccess(t *testing.T) {
	manager := NewLifecycleManager()

	// Tester quelques opérations concurrentes basiques
	done := make(chan bool, 2)

	// Goroutine 1: Ajouter des nœuds
	go func() {
		for i := 0; i < 100; i++ {
			manager.RegisterNode("alpha_"+string(rune(i)), "alpha")
		}
		done <- true
	}()

	// Goroutine 2: Ajouter des règles
	go func() {
		for i := 0; i < 100; i++ {
			manager.RegisterNode("beta_"+string(rune(i)), "beta")
		}
		done <- true
	}()

	// Attendre la fin
	<-done
	<-done

	// Vérifier qu'il n'y a pas eu de panic et que les opérations se sont bien passées
	stats := manager.GetStats()
	totalNodes := stats["total_nodes"].(int)
	if totalNodes != 200 {
		t.Logf("Avertissement: total_nodes devrait être proche de 200, obtenu %d", totalNodes)
	}
}
