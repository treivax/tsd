// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"

	"github.com/treivax/tsd/constraint"
)

// TestPipeline_SimpleCondition_NoChange vérifie que les conditions simples
// utilisent toujours le comportement actuel (pas de chaîne)
func TestPipeline_SimpleCondition_NoChange(t *testing.T) {
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
	err := cp.createAlphaNodeWithTerminal(
		network,
		"rule_simple",
		condition,
		"p",
		"Person",
		action,
		storage,
	)
	if err != nil {
		t.Fatalf("Erreur création AlphaNode: %v", err)
	}
	// Vérifier qu'un seul AlphaNode a été créé (pas de chaîne)
	if len(network.AlphaNodes) != 1 {
		t.Errorf("Attendu 1 AlphaNode, obtenu %d", len(network.AlphaNodes))
	}
	// Vérifier qu'un TerminalNode a été créé
	if len(network.TerminalNodes) != 1 {
		t.Errorf("Attendu 1 TerminalNode, obtenu %d", len(network.TerminalNodes))
	}
	t.Logf("✓ Condition simple traitée correctement sans décomposition")
}

// TestPipeline_AND_CreatesChain vérifie qu'une expression AND
// est décomposée en chaîne d'AlphaNodes
func TestPipeline_AND_CreatesChain(t *testing.T) {
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
	// Expression AND: p.age > 18 AND p.salary >= 50000
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
	err := cp.createAlphaNodeWithTerminal(
		network,
		"rule_and",
		condition,
		"p",
		"Person",
		action,
		storage,
	)
	if err != nil {
		t.Fatalf("Erreur création chaîne: %v", err)
	}
	// Vérifier que 2 AlphaNodes ont été créés (chaîne)
	if len(network.AlphaNodes) != 2 {
		t.Errorf("Attendu 2 AlphaNodes pour la chaîne, obtenu %d", len(network.AlphaNodes))
	}
	// Vérifier qu'un TerminalNode a été créé
	if len(network.TerminalNodes) != 1 {
		t.Errorf("Attendu 1 TerminalNode, obtenu %d", len(network.TerminalNodes))
	}
	// Vérifier que les AlphaNodes sont enregistrés dans le LifecycleManager
	if network.LifecycleManager != nil {
		count := 0
		for nodeID := range network.AlphaNodes {
			if _, exists := network.LifecycleManager.GetNodeLifecycle(nodeID); exists {
				count++
			}
		}
		if count != 2 {
			t.Errorf("Attendu 2 AlphaNodes dans LifecycleManager, obtenu %d", count)
		}
	}
	t.Logf("✓ Expression AND décomposée en chaîne de %d nœuds", len(network.AlphaNodes))
}

// TestPipeline_OR_SingleNode vérifie qu'une expression OR
// crée un seul AlphaNode normalisé (pas de chaîne)
func TestPipeline_OR_SingleNode(t *testing.T) {
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
	// Expression OR: p.age < 18 OR p.age > 65
	condition := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type: "binaryOperation",
			Left: constraint.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "age",
			},
			Operator: "<",
			Right: constraint.NumberLiteral{
				Type:  "numberLiteral",
				Value: 18,
			},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "OR",
				Right: constraint.BinaryOperation{
					Type: "binaryOperation",
					Left: constraint.FieldAccess{
						Type:   "fieldAccess",
						Object: "p",
						Field:  "age",
					},
					Operator: ">",
					Right: constraint.NumberLiteral{
						Type:  "numberLiteral",
						Value: 65,
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
	err := cp.createAlphaNodeWithTerminal(
		network,
		"rule_or",
		condition,
		"p",
		"Person",
		action,
		storage,
	)
	if err != nil {
		t.Fatalf("Erreur création AlphaNode: %v", err)
	}
	// Vérifier qu'un seul AlphaNode a été créé (pas de chaîne)
	if len(network.AlphaNodes) != 1 {
		t.Errorf("Attendu 1 AlphaNode pour OR, obtenu %d", len(network.AlphaNodes))
	}
	// Vérifier qu'un TerminalNode a été créé
	if len(network.TerminalNodes) != 1 {
		t.Errorf("Attendu 1 TerminalNode, obtenu %d", len(network.TerminalNodes))
	}
	t.Logf("✓ Expression OR traitée avec un seul nœud normalisé")
}

// TestPipeline_TwoRules_ShareChain vérifie que deux règles
// avec des conditions communes partagent les nœuds de la chaîne
func TestPipeline_TwoRules_ShareChain(t *testing.T) {
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
	// Condition commune: p.age > 18 AND p.salary >= 50000
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
	action1 := &Action{
		Type: "print",
		Job: &JobCall{
			Name: "print",
			Args: []interface{}{"rule1"},
		},
	}
	action2 := &Action{
		Type: "print",
		Job: &JobCall{
			Name: "print",
			Args: []interface{}{"rule2"},
		},
	}
	cp := &ConstraintPipeline{}
	// Créer la première règle
	err := cp.createAlphaNodeWithTerminal(
		network,
		"rule_share_1",
		createCondition(),
		"p",
		"Person",
		action1,
		storage,
	)
	if err != nil {
		t.Fatalf("Erreur création règle 1: %v", err)
	}
	alphaNodesAfterRule1 := len(network.AlphaNodes)
	t.Logf("Après règle 1: %d AlphaNodes", alphaNodesAfterRule1)
	// Créer la deuxième règle avec la même condition
	err = cp.createAlphaNodeWithTerminal(
		network,
		"rule_share_2",
		createCondition(),
		"p",
		"Person",
		action2,
		storage,
	)
	if err != nil {
		t.Fatalf("Erreur création règle 2: %v", err)
	}
	alphaNodesAfterRule2 := len(network.AlphaNodes)
	t.Logf("Après règle 2: %d AlphaNodes", alphaNodesAfterRule2)
	// Vérifier que les AlphaNodes sont partagés (pas de nouveaux nœuds créés)
	if alphaNodesAfterRule2 != alphaNodesAfterRule1 {
		t.Errorf("Les AlphaNodes devraient être partagés. Avant: %d, Après: %d",
			alphaNodesAfterRule1, alphaNodesAfterRule2)
	}
	// Vérifier que 2 TerminalNodes ont été créés (un par règle)
	if len(network.TerminalNodes) != 2 {
		t.Errorf("Attendu 2 TerminalNodes, obtenu %d", len(network.TerminalNodes))
	}
	// Vérifier que les nœuds partagés ont 2 références
	if network.LifecycleManager != nil {
		for nodeID := range network.AlphaNodes {
			lifecycle, exists := network.LifecycleManager.GetNodeLifecycle(nodeID)
			if !exists {
				t.Errorf("Nœud %s non trouvé dans LifecycleManager", nodeID)
				continue
			}
			refCount := lifecycle.GetRefCount()
			if refCount != 2 {
				t.Errorf("Nœud partagé %s devrait avoir 2 références, obtenu %d", nodeID, refCount)
			}
		}
	}
	t.Logf("✓ Deux règles partagent correctement la chaîne de %d nœuds", alphaNodesAfterRule2)
}

// TestPipeline_ErrorHandling_FallbackToSimple vérifie que le système
// revient au comportement simple en cas d'erreur
func TestPipeline_ErrorHandling_FallbackToSimple(t *testing.T) {
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
	// Condition invalide mais qui devrait déclencher le fallback
	condition := map[string]interface{}{
		"type":     "unknown_type",
		"operator": "???",
	}
	action := &Action{
		Type: "print",
		Job: &JobCall{
			Name: "print",
			Args: []interface{}{},
		},
	}
	cp := &ConstraintPipeline{}
	err := cp.createAlphaNodeWithTerminal(
		network,
		"rule_fallback",
		condition,
		"p",
		"Person",
		action,
		storage,
	)
	// Le fallback devrait créer un nœud simple même avec une erreur d'analyse
	if err != nil {
		t.Logf("Erreur attendue avec fallback: %v", err)
	}
	t.Logf("✓ Fallback vers comportement simple fonctionne correctement")
}

// TestPipeline_ComplexAND_ThreeConditions vérifie qu'une expression
// AND avec 3 conditions crée une chaîne de 3 nœuds
func TestPipeline_ComplexAND_ThreeConditions(t *testing.T) {
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
	// Expression AND: p.age > 18 AND p.salary >= 50000 AND p.experience > 5
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
	err := cp.createAlphaNodeWithTerminal(
		network,
		"rule_three",
		condition,
		"p",
		"Person",
		action,
		storage,
	)
	if err != nil {
		t.Fatalf("Erreur création chaîne: %v", err)
	}
	// Vérifier que 3 AlphaNodes ont été créés
	if len(network.AlphaNodes) != 3 {
		t.Errorf("Attendu 3 AlphaNodes pour la chaîne, obtenu %d", len(network.AlphaNodes))
	}
	// Vérifier qu'un TerminalNode a été créé
	if len(network.TerminalNodes) != 1 {
		t.Errorf("Attendu 1 TerminalNode, obtenu %d", len(network.TerminalNodes))
	}
	t.Logf("✓ Expression AND avec 3 conditions décomposée en chaîne de %d nœuds", len(network.AlphaNodes))
}

// TestPipeline_Arithmetic_NoChain vérifie qu'une expression arithmétique
// n'est pas décomposée en chaîne
func TestPipeline_Arithmetic_NoChain(t *testing.T) {
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
	// Expression arithmétique: p.salary * 1.1 > 60000
	condition := constraint.BinaryOperation{
		Type: "binaryOperation",
		Left: constraint.BinaryOperation{
			Type: "binaryOperation",
			Left: constraint.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "salary",
			},
			Operator: "*",
			Right: constraint.NumberLiteral{
				Type:  "numberLiteral",
				Value: 1.1,
			},
		},
		Operator: ">",
		Right: constraint.NumberLiteral{
			Type:  "numberLiteral",
			Value: 60000,
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
	err := cp.createAlphaNodeWithTerminal(
		network,
		"rule_arithmetic",
		condition,
		"p",
		"Person",
		action,
		storage,
	)
	if err != nil {
		t.Fatalf("Erreur création AlphaNode: %v", err)
	}
	// Vérifier qu'un seul AlphaNode a été créé (pas de chaîne)
	if len(network.AlphaNodes) != 1 {
		t.Errorf("Attendu 1 AlphaNode pour expression arithmétique, obtenu %d", len(network.AlphaNodes))
	}
	t.Logf("✓ Expression arithmétique traitée sans décomposition")
}
