// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"github.com/treivax/tsd/constraint"
	"testing"
)

// TestRemoveChain_AllNodesUnique_DeletesAll vérifie qu'une chaîne avec des nœuds
// uniques (non partagés) est complètement supprimée
func TestRemoveChain_AllNodesUnique_DeletesAll(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.AlphaSharingManager = NewAlphaSharingRegistry()
	network.LifecycleManager = NewLifecycleManager()
	// Créer un TypeNode pour Person
	typeDef := TypeDefinition{
		Name:   "Person",
		Fields: []Field{},
	}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode
	// Créer une règle avec chaîne: p.age > 18 AND p.salary >= 50000
	condition := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type: "binaryOperation",
			Left: constraint.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "age",
			},
			Operator: ">",
			Right: constraint.NumberLiteral{
				Type:  "numberLiteral",
				Value: 18,
			},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type: "binaryOperation",
					Left: constraint.FieldAccess{
						Type:   "fieldAccess",
						Object: "p",
						Field:  "salary",
					},
					Operator: ">=",
					Right: constraint.NumberLiteral{
						Type:  "numberLiteral",
						Value: 50000,
					},
				},
			},
		},
	}
	action := &Action{
		Type: "print",
		Job: &JobCall{
			Name: "print",
			Args: []interface{}{},
		},
	}
	cp := &ConstraintPipeline{}
	err := cp.createAlphaNodeWithTerminal(network, "rule_unique", condition, "p", "Person", action, storage)
	if err != nil {
		t.Fatalf("Erreur création règle: %v", err)
	}
	// Vérifier que les nœuds ont été créés
	initialAlphaCount := len(network.AlphaNodes)
	initialTerminalCount := len(network.TerminalNodes)
	if initialAlphaCount != 2 {
		t.Errorf("Attendu 2 AlphaNodes, obtenu %d", initialAlphaCount)
	}
	if initialTerminalCount != 1 {
		t.Errorf("Attendu 1 TerminalNode, obtenu %d", initialTerminalCount)
	}
	// Supprimer la règle
	err = network.RemoveRule("rule_unique")
	if err != nil {
		t.Fatalf("Erreur suppression règle: %v", err)
	}
	// Vérifier que tous les nœuds ont été supprimés
	if len(network.AlphaNodes) != 0 {
		t.Errorf("Attendu 0 AlphaNodes après suppression, obtenu %d", len(network.AlphaNodes))
	}
	if len(network.TerminalNodes) != 0 {
		t.Errorf("Attendu 0 TerminalNodes après suppression, obtenu %d", len(network.TerminalNodes))
	}
	// Vérifier que les nœuds ont été supprimés du LifecycleManager
	for nodeID := range network.AlphaNodes {
		if _, exists := network.LifecycleManager.GetNodeLifecycle(nodeID); exists {
			t.Errorf("Nœud %s encore présent dans LifecycleManager", nodeID)
		}
	}
	t.Logf("✓ Chaîne unique supprimée complètement: %d AlphaNodes + %d TerminalNode",
		initialAlphaCount, initialTerminalCount)
}

// TestRemoveChain_PartialSharing_DeletesOnlyUnused vérifie que seuls les nœuds
// non partagés sont supprimés lors de la suppression d'une règle
func TestRemoveChain_PartialSharing_DeletesOnlyUnused(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.AlphaSharingManager = NewAlphaSharingRegistry()
	network.LifecycleManager = NewLifecycleManager()
	// Créer un TypeNode pour Person
	typeDef := TypeDefinition{
		Name:   "Person",
		Fields: []Field{},
	}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode
	// Règle 1: p.age > 18 AND p.salary >= 50000
	createCondition1 := func() interface{} {
		return constraint.LogicalExpression{
			Type: "logicalExpr",
			Left: constraint.BinaryOperation{
				Type: "binaryOperation",
				Left: constraint.FieldAccess{
					Type:   "fieldAccess",
					Object: "p",
					Field:  "age",
				},
				Operator: ">",
				Right: constraint.NumberLiteral{
					Type:  "numberLiteral",
					Value: 18,
				},
			},
			Operations: []constraint.LogicalOperation{
				{
					Op: "AND",
					Right: constraint.BinaryOperation{
						Type: "binaryOperation",
						Left: constraint.FieldAccess{
							Type:   "fieldAccess",
							Object: "p",
							Field:  "salary",
						},
						Operator: ">=",
						Right: constraint.NumberLiteral{
							Type:  "numberLiteral",
							Value: 50000,
						},
					},
				},
			},
		}
	}
	// Règle 2: p.age > 18 AND p.experience > 5 (partage p.age > 18)
	createCondition2 := func() interface{} {
		return constraint.LogicalExpression{
			Type: "logicalExpr",
			Left: constraint.BinaryOperation{
				Type: "binaryOperation",
				Left: constraint.FieldAccess{
					Type:   "fieldAccess",
					Object: "p",
					Field:  "age",
				},
				Operator: ">",
				Right: constraint.NumberLiteral{
					Type:  "numberLiteral",
					Value: 18,
				},
			},
			Operations: []constraint.LogicalOperation{
				{
					Op: "AND",
					Right: constraint.BinaryOperation{
						Type: "binaryOperation",
						Left: constraint.FieldAccess{
							Type:   "fieldAccess",
							Object: "p",
							Field:  "experience",
						},
						Operator: ">",
						Right: constraint.NumberLiteral{
							Type:  "numberLiteral",
							Value: 5,
						},
					},
				},
			},
		}
	}
	action1 := &Action{Type: "print", Job: &JobCall{Name: "action1", Args: []interface{}{}}}
	action2 := &Action{Type: "print", Job: &JobCall{Name: "action2", Args: []interface{}{}}}
	cp := &ConstraintPipeline{}
	// Créer les deux règles
	err := cp.createAlphaNodeWithTerminal(network, "rule_partial_1", createCondition1(), "p", "Person", action1, storage)
	if err != nil {
		t.Fatalf("Erreur création règle 1: %v", err)
	}
	alphaCountAfterRule1 := len(network.AlphaNodes)
	err = cp.createAlphaNodeWithTerminal(network, "rule_partial_2", createCondition2(), "p", "Person", action2, storage)
	if err != nil {
		t.Fatalf("Erreur création règle 2: %v", err)
	}
	alphaCountAfterRule2 := len(network.AlphaNodes)
	t.Logf("AlphaNodes après règle 1: %d, après règle 2: %d", alphaCountAfterRule1, alphaCountAfterRule2)
	// Identifier le nœud partagé (p.age > 18)
	var sharedNodeID string
	for nodeID := range network.AlphaNodes {
		lifecycle, exists := network.LifecycleManager.GetNodeLifecycle(nodeID)
		if exists && lifecycle.GetRefCount() == 2 {
			sharedNodeID = nodeID
			break
		}
	}
	if sharedNodeID == "" {
		t.Fatal("Aucun nœud partagé trouvé")
	}
	t.Logf("Nœud partagé identifié: %s", sharedNodeID)
	// Supprimer la première règle
	err = network.RemoveRule("rule_partial_1")
	if err != nil {
		t.Fatalf("Erreur suppression règle 1: %v", err)
	}
	// Vérifier que le nœud partagé existe toujours
	if _, exists := network.AlphaNodes[sharedNodeID]; !exists {
		t.Errorf("Le nœud partagé %s a été supprimé à tort", sharedNodeID)
	}
	// Vérifier que le nœud partagé a maintenant RefCount == 1
	lifecycle, exists := network.LifecycleManager.GetNodeLifecycle(sharedNodeID)
	if !exists {
		t.Errorf("Le nœud partagé %s n'existe plus dans LifecycleManager", sharedNodeID)
	} else if lifecycle.GetRefCount() != 1 {
		t.Errorf("Le nœud partagé devrait avoir RefCount=1, obtenu %d", lifecycle.GetRefCount())
	}
	// Vérifier que les nœuds spécifiques à la règle 1 ont été supprimés
	alphaCountAfterRemoval := len(network.AlphaNodes)
	expectedCount := alphaCountAfterRule2 - 1 // Un nœud supprimé (salary)
	if alphaCountAfterRemoval != expectedCount {
		t.Errorf("Attendu %d AlphaNodes après suppression, obtenu %d", expectedCount, alphaCountAfterRemoval)
	}
	t.Logf("✓ Suppression partielle correcte: nœud partagé conservé, nœud unique supprimé")
}

// TestRemoveChain_CompleteSharing_DeletesNone vérifie que si deux règles
// partagent tous leurs nœuds, la suppression d'une règle ne supprime aucun nœud
func TestRemoveChain_CompleteSharing_DeletesNone(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.AlphaSharingManager = NewAlphaSharingRegistry()
	network.LifecycleManager = NewLifecycleManager()
	// Créer un TypeNode pour Person
	typeDef := TypeDefinition{
		Name:   "Person",
		Fields: []Field{},
	}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode
	// Même condition pour les deux règles
	createCondition := func() interface{} {
		return constraint.LogicalExpression{
			Type: "logicalExpr",
			Left: constraint.BinaryOperation{
				Type: "binaryOperation",
				Left: constraint.FieldAccess{
					Type:   "fieldAccess",
					Object: "p",
					Field:  "age",
				},
				Operator: ">",
				Right: constraint.NumberLiteral{
					Type:  "numberLiteral",
					Value: 18,
				},
			},
			Operations: []constraint.LogicalOperation{
				{
					Op: "AND",
					Right: constraint.BinaryOperation{
						Type: "binaryOperation",
						Left: constraint.FieldAccess{
							Type:   "fieldAccess",
							Object: "p",
							Field:  "salary",
						},
						Operator: ">=",
						Right: constraint.NumberLiteral{
							Type:  "numberLiteral",
							Value: 50000,
						},
					},
				},
			},
		}
	}
	action1 := &Action{Type: "print", Job: &JobCall{Name: "action1", Args: []interface{}{}}}
	action2 := &Action{Type: "print", Job: &JobCall{Name: "action2", Args: []interface{}{}}}
	cp := &ConstraintPipeline{}
	// Créer deux règles avec la même condition
	err := cp.createAlphaNodeWithTerminal(network, "rule_complete_1", createCondition(), "p", "Person", action1, storage)
	if err != nil {
		t.Fatalf("Erreur création règle 1: %v", err)
	}
	err = cp.createAlphaNodeWithTerminal(network, "rule_complete_2", createCondition(), "p", "Person", action2, storage)
	if err != nil {
		t.Fatalf("Erreur création règle 2: %v", err)
	}
	alphaCountBefore := len(network.AlphaNodes)
	terminalCountBefore := len(network.TerminalNodes)
	// Vérifier que tous les AlphaNodes ont RefCount == 2
	for nodeID := range network.AlphaNodes {
		lifecycle, exists := network.LifecycleManager.GetNodeLifecycle(nodeID)
		if !exists {
			t.Errorf("Nœud %s non trouvé dans LifecycleManager", nodeID)
			continue
		}
		if lifecycle.GetRefCount() != 2 {
			t.Errorf("Nœud %s devrait avoir RefCount=2, obtenu %d", nodeID, lifecycle.GetRefCount())
		}
	}
	// Supprimer la première règle
	err = network.RemoveRule("rule_complete_1")
	if err != nil {
		t.Fatalf("Erreur suppression règle 1: %v", err)
	}
	alphaCountAfter := len(network.AlphaNodes)
	terminalCountAfter := len(network.TerminalNodes)
	// Vérifier qu'aucun AlphaNode n'a été supprimé
	if alphaCountAfter != alphaCountBefore {
		t.Errorf("Les AlphaNodes ne devraient pas être supprimés. Avant: %d, Après: %d",
			alphaCountBefore, alphaCountAfter)
	}
	// Vérifier qu'un TerminalNode a été supprimé (celui de la règle 1)
	expectedTerminals := terminalCountBefore - 1
	if terminalCountAfter != expectedTerminals {
		t.Errorf("Attendu %d TerminalNodes, obtenu %d", expectedTerminals, terminalCountAfter)
	}
	// Vérifier que tous les AlphaNodes ont maintenant RefCount == 1
	for nodeID := range network.AlphaNodes {
		lifecycle, exists := network.LifecycleManager.GetNodeLifecycle(nodeID)
		if !exists {
			t.Errorf("Nœud %s non trouvé dans LifecycleManager", nodeID)
			continue
		}
		if lifecycle.GetRefCount() != 1 {
			t.Errorf("Nœud %s devrait avoir RefCount=1 après suppression, obtenu %d",
				nodeID, lifecycle.GetRefCount())
		}
	}
	t.Logf("✓ Partage complet: aucun AlphaNode supprimé, RefCount correctement décrémenté")
}

// TestRemoveRule_WithChain_CorrectCleanup vérifie que RemoveRule() gère
// correctement la suppression d'une règle avec chaîne
func TestRemoveRule_WithChain_CorrectCleanup(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.AlphaSharingManager = NewAlphaSharingRegistry()
	network.LifecycleManager = NewLifecycleManager()
	// Créer un TypeNode pour Person
	typeDef := TypeDefinition{
		Name:   "Person",
		Fields: []Field{},
	}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode
	// Règle avec chaîne de 3 conditions
	condition := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type: "binaryOperation",
			Left: constraint.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "age",
			},
			Operator: ">",
			Right: constraint.NumberLiteral{
				Type:  "numberLiteral",
				Value: 18,
			},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type: "binaryOperation",
					Left: constraint.FieldAccess{
						Type:   "fieldAccess",
						Object: "p",
						Field:  "salary",
					},
					Operator: ">=",
					Right: constraint.NumberLiteral{
						Type:  "numberLiteral",
						Value: 50000,
					},
				},
			},
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type: "binaryOperation",
					Left: constraint.FieldAccess{
						Type:   "fieldAccess",
						Object: "p",
						Field:  "experience",
					},
					Operator: ">",
					Right: constraint.NumberLiteral{
						Type:  "numberLiteral",
						Value: 5,
					},
				},
			},
		},
	}
	action := &Action{
		Type: "print",
		Job: &JobCall{
			Name: "print",
			Args: []interface{}{},
		},
	}
	cp := &ConstraintPipeline{}
	err := cp.createAlphaNodeWithTerminal(network, "rule_cleanup", condition, "p", "Person", action, storage)
	if err != nil {
		t.Fatalf("Erreur création règle: %v", err)
	}
	// Vérifier la création
	if len(network.AlphaNodes) != 3 {
		t.Errorf("Attendu 3 AlphaNodes, obtenu %d", len(network.AlphaNodes))
	}
	// Récupérer les IDs des nœuds avant suppression
	nodeIDs := make([]string, 0, len(network.AlphaNodes))
	for nodeID := range network.AlphaNodes {
		nodeIDs = append(nodeIDs, nodeID)
	}
	// Supprimer la règle
	err = network.RemoveRule("rule_cleanup")
	if err != nil {
		t.Fatalf("Erreur suppression règle: %v", err)
	}
	// Vérifier la suppression complète
	if len(network.AlphaNodes) != 0 {
		t.Errorf("Attendu 0 AlphaNodes après suppression, obtenu %d", len(network.AlphaNodes))
	}
	if len(network.TerminalNodes) != 0 {
		t.Errorf("Attendu 0 TerminalNodes après suppression, obtenu %d", len(network.TerminalNodes))
	}
	// Vérifier que les nœuds ont été supprimés du LifecycleManager
	for _, nodeID := range nodeIDs {
		if _, exists := network.LifecycleManager.GetNodeLifecycle(nodeID); exists {
			t.Errorf("Nœud %s encore présent dans LifecycleManager après suppression", nodeID)
		}
	}
	// Vérifier que les nœuds ont été supprimés du AlphaSharingManager
	if network.AlphaSharingManager != nil {
		stats := network.AlphaSharingManager.GetStats()
		if totalNodes, ok := stats["total_nodes"].(int); ok && totalNodes != 0 {
			t.Errorf("AlphaSharingManager devrait être vide, contient %d nœud(s)", totalNodes)
		}
	}
	t.Logf("✓ Nettoyage complet: tous les nœuds supprimés du réseau, LifecycleManager et AlphaSharingManager")
}

// TestRemoveRule_MultipleChains_IndependentCleanup vérifie que la suppression
// de plusieurs règles avec chaînes se fait indépendamment
func TestRemoveRule_MultipleChains_IndependentCleanup(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.AlphaSharingManager = NewAlphaSharingRegistry()
	network.LifecycleManager = NewLifecycleManager()
	// Créer un TypeNode pour Person
	typeDef := TypeDefinition{
		Name:   "Person",
		Fields: []Field{},
	}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode
	// Règle 1: p.age > 18 AND p.salary >= 50000
	condition1 := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type: "binaryOperation",
			Left: constraint.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "age",
			},
			Operator: ">",
			Right: constraint.NumberLiteral{
				Type:  "numberLiteral",
				Value: 18,
			},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type: "binaryOperation",
					Left: constraint.FieldAccess{
						Type:   "fieldAccess",
						Object: "p",
						Field:  "salary",
					},
					Operator: ">=",
					Right: constraint.NumberLiteral{
						Type:  "numberLiteral",
						Value: 50000,
					},
				},
			},
		},
	}
	// Règle 2: p.name == "John" AND p.city == "NYC"
	condition2 := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type: "binaryOperation",
			Left: constraint.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "name",
			},
			Operator: "==",
			Right: constraint.StringLiteral{
				Type:  "stringLiteral",
				Value: "John",
			},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type: "binaryOperation",
					Left: constraint.FieldAccess{
						Type:   "fieldAccess",
						Object: "p",
						Field:  "city",
					},
					Operator: "==",
					Right: constraint.StringLiteral{
						Type:  "stringLiteral",
						Value: "NYC",
					},
				},
			},
		},
	}
	action1 := &Action{Type: "print", Job: &JobCall{Name: "action1", Args: []interface{}{}}}
	action2 := &Action{Type: "print", Job: &JobCall{Name: "action2", Args: []interface{}{}}}
	cp := &ConstraintPipeline{}
	// Créer les deux règles
	err := cp.createAlphaNodeWithTerminal(network, "rule_independent_1", condition1, "p", "Person", action1, storage)
	if err != nil {
		t.Fatalf("Erreur création règle 1: %v", err)
	}
	err = cp.createAlphaNodeWithTerminal(network, "rule_independent_2", condition2, "p", "Person", action2, storage)
	if err != nil {
		t.Fatalf("Erreur création règle 2: %v", err)
	}
	totalAlphaNodes := len(network.AlphaNodes)
	totalTerminalNodes := len(network.TerminalNodes)
	if totalAlphaNodes != 4 {
		t.Errorf("Attendu 4 AlphaNodes (2 par règle), obtenu %d", totalAlphaNodes)
	}
	if totalTerminalNodes != 2 {
		t.Errorf("Attendu 2 TerminalNodes, obtenu %d", totalTerminalNodes)
	}
	// Identifier les nœuds de la règle 1
	rule1Nodes := network.LifecycleManager.GetNodesForRule("rule_independent_1")
	rule1AlphaCount := 0
	for _, nodeID := range rule1Nodes {
		lifecycle, _ := network.LifecycleManager.GetNodeLifecycle(nodeID)
		if lifecycle.NodeType == "alpha" {
			rule1AlphaCount++
		}
	}
	// Supprimer la règle 1
	err = network.RemoveRule("rule_independent_1")
	if err != nil {
		t.Fatalf("Erreur suppression règle 1: %v", err)
	}
	// Vérifier que seuls les nœuds de la règle 1 ont été supprimés
	expectedAlphaNodes := totalAlphaNodes - rule1AlphaCount
	if len(network.AlphaNodes) != expectedAlphaNodes {
		t.Errorf("Attendu %d AlphaNodes après suppression règle 1, obtenu %d",
			expectedAlphaNodes, len(network.AlphaNodes))
	}
	expectedTerminalNodes := totalTerminalNodes - 1
	if len(network.TerminalNodes) != expectedTerminalNodes {
		t.Errorf("Attendu %d TerminalNodes après suppression règle 1, obtenu %d",
			expectedTerminalNodes, len(network.TerminalNodes))
	}
	// Vérifier que les nœuds de la règle 2 sont intacts
	rule2Nodes := network.LifecycleManager.GetNodesForRule("rule_independent_2")
	if len(rule2Nodes) == 0 {
		t.Error("Les nœuds de la règle 2 ne devraient pas être affectés")
	}
	for _, nodeID := range rule2Nodes {
		lifecycle, exists := network.LifecycleManager.GetNodeLifecycle(nodeID)
		if !exists {
			t.Errorf("Nœud %s de la règle 2 a été supprimé à tort", nodeID)
		} else if lifecycle.GetRefCount() != 1 {
			t.Errorf("Nœud %s de la règle 2 devrait avoir RefCount=1, obtenu %d",
				nodeID, lifecycle.GetRefCount())
		}
	}
	// Supprimer la règle 2
	err = network.RemoveRule("rule_independent_2")
	if err != nil {
		t.Fatalf("Erreur suppression règle 2: %v", err)
	}
	// Vérifier que tout est supprimé
	if len(network.AlphaNodes) != 0 {
		t.Errorf("Attendu 0 AlphaNodes après suppression des deux règles, obtenu %d", len(network.AlphaNodes))
	}
	if len(network.TerminalNodes) != 0 {
		t.Errorf("Attendu 0 TerminalNodes après suppression des deux règles, obtenu %d", len(network.TerminalNodes))
	}
	t.Logf("✓ Suppression indépendante: chaque règle supprimée sans affecter l'autre")
}

// TestRemoveRule_SimpleCondition_BackwardCompatibility vérifie que la suppression
// de règles simples (sans chaîne) fonctionne toujours correctement
func TestRemoveRule_SimpleCondition_BackwardCompatibility(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.AlphaSharingManager = NewAlphaSharingRegistry()
	network.LifecycleManager = NewLifecycleManager()
	// Créer un TypeNode pour Person
	typeDef := TypeDefinition{
		Name:   "Person",
		Fields: []Field{},
	}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode
	// Condition simple: p.age > 18
	condition := map[string]interface{}{
		"type": "binaryOperation",
		"left": constraint.FieldAccess{
			Type:   "fieldAccess",
			Object: "p",
			Field:  "age",
		},
		"operator": ">",
		"right": constraint.NumberLiteral{
			Type:  "numberLiteral",
			Value: 18,
		},
	}
	action := &Action{
		Type: "print",
		Job: &JobCall{
			Name: "print",
			Args: []interface{}{},
		},
	}
	cp := &ConstraintPipeline{}
	err := cp.createAlphaNodeWithTerminal(network, "rule_simple", condition, "p", "Person", action, storage)
	if err != nil {
		t.Fatalf("Erreur création règle simple: %v", err)
	}
	// Vérifier qu'un seul AlphaNode a été créé
	if len(network.AlphaNodes) != 1 {
		t.Errorf("Attendu 1 AlphaNode pour règle simple, obtenu %d", len(network.AlphaNodes))
	}
	// Supprimer la règle
	err = network.RemoveRule("rule_simple")
	if err != nil {
		t.Fatalf("Erreur suppression règle simple: %v", err)
	}
	// Vérifier la suppression complète
	if len(network.AlphaNodes) != 0 {
		t.Errorf("Attendu 0 AlphaNodes après suppression, obtenu %d", len(network.AlphaNodes))
	}
	if len(network.TerminalNodes) != 0 {
		t.Errorf("Attendu 0 TerminalNodes après suppression, obtenu %d", len(network.TerminalNodes))
	}
	t.Logf("✓ Rétrocompatibilité: règle simple supprimée correctement")
}
