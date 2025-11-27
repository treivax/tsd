// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"os"
	"testing"

	"github.com/treivax/tsd/constraint"
)

// TestOR_SingleNode_NotDecomposed vérifie qu'une expression OR n'est pas décomposée
// en plusieurs AlphaNodes mais reste un seul nœud atomique
func TestOR_SingleNode_NotDecomposed(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Créer un TypeNode pour Person
	typeDef := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "status", Type: "string"},
			{Name: "age", Type: "number"},
		},
	}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode

	// Expression OR: p.status == "VIP" OR p.age > 18
	orExpr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type: "binaryOperation",
			Left: constraint.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "status",
			},
			Operator: "==",
			Right: constraint.StringLiteral{
				Type:  "stringLiteral",
				Value: "VIP",
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
						Value: 18.0,
					},
				},
			},
		},
	}

	// Créer une règle avec cette expression OR
	action := &Action{
		Type: "log",
		Job: JobCall{
			Name: "log",
			Args: []interface{}{"VIP or Adult detected"},
		},
	}

	cp := &ConstraintPipeline{}

	// Créer l'AlphaNode avec l'expression OR
	err := cp.createAlphaNodeWithTerminal(
		network,
		"rule_or_test",
		orExpr,
		"p",
		"Person",
		action,
		storage,
	)

	if err != nil {
		t.Fatalf("Erreur création AlphaNode OR: %v", err)
	}

	// Vérifier qu'un seul AlphaNode a été créé (pas de décomposition)
	alphaCount := len(network.AlphaNodes)
	if alphaCount != 1 {
		t.Errorf("Expression OR décomposée: attendu 1 AlphaNode, obtenu %d", alphaCount)
	}

	// Vérifier qu'un TerminalNode a été créé
	if len(network.TerminalNodes) != 1 {
		t.Errorf("Attendu 1 TerminalNode, obtenu %d", len(network.TerminalNodes))
	}

	// Vérifier la connexion TypeNode -> AlphaNode -> TerminalNode
	if len(typeNode.Children) != 1 {
		t.Errorf("TypeNode devrait avoir 1 enfant, a %d", len(typeNode.Children))
	}

	t.Logf("✓ Expression OR crée un seul AlphaNode (non décomposée)")
}

// TestOR_Normalization_OrderIndependent vérifie que l'ordre des termes OR
// n'affecte pas le hash après normalisation (partage possible)
func TestOR_Normalization_OrderIndependent(t *testing.T) {

	// Expression 1: p.status == "VIP" OR p.age > 18
	orExpr1 := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type: "binaryOperation",
			Left: constraint.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "status",
			},
			Operator: "==",
			Right: constraint.StringLiteral{
				Type:  "stringLiteral",
				Value: "VIP",
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
						Value: 18.0,
					},
				},
			},
		},
	}

	// Expression 2: p.age > 18 OR p.status == "VIP" (ordre inversé)
	orExpr2 := constraint.LogicalExpression{
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
				Value: 18.0,
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
						Field:  "status",
					},
					Operator: "==",
					Right: constraint.StringLiteral{
						Type:  "stringLiteral",
						Value: "VIP",
					},
				},
			},
		},
	}

	// Normaliser les deux expressions
	normalized1, err := NormalizeORExpression(orExpr1)
	if err != nil {
		t.Fatalf("Erreur normalisation expr1: %v", err)
	}

	normalized2, err := NormalizeORExpression(orExpr2)
	if err != nil {
		t.Fatalf("Erreur normalisation expr2: %v", err)
	}

	// Calculer les hashes
	condition1 := map[string]interface{}{
		"type":       "constraint",
		"constraint": normalized1,
	}

	condition2 := map[string]interface{}{
		"type":       "constraint",
		"constraint": normalized2,
	}

	hash1, err := ConditionHash(condition1, "p")
	if err != nil {
		t.Fatalf("Erreur calcul hash1: %v", err)
	}

	hash2, err := ConditionHash(condition2, "p")
	if err != nil {
		t.Fatalf("Erreur calcul hash2: %v", err)
	}

	// Les hashes doivent être identiques après normalisation
	if hash1 != hash2 {
		t.Errorf("Normalisation OR échouée: hashes différents\n  Hash1: %s\n  Hash2: %s", hash1, hash2)
	}

	t.Logf("✓ Normalisation OR indépendante de l'ordre des termes")
	t.Logf("  Hash commun: %s", hash1)
}

// TestMixedAND_OR_SingleNode vérifie qu'une expression mixte (AND + OR)
// crée un seul AlphaNode normalisé
func TestMixedAND_OR_SingleNode(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Créer un TypeNode pour Person
	typeDef := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "status", Type: "string"},
			{Name: "age", Type: "number"},
			{Name: "country", Type: "string"},
		},
	}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode

	// Expression mixte: (p.age > 18 OR p.status == "VIP") AND p.country == "FR"
	// Note: pour simplifier, on utilise une représentation avec AND de plus haut niveau
	mixedExpr := constraint.LogicalExpression{
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
				Value: 18.0,
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
						Field:  "status",
					},
					Operator: "==",
					Right: constraint.StringLiteral{
						Type:  "stringLiteral",
						Value: "VIP",
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
						Field:  "country",
					},
					Operator: "==",
					Right: constraint.StringLiteral{
						Type:  "stringLiteral",
						Value: "FR",
					},
				},
			},
		},
	}

	action := &Action{
		Type: "log",
		Job: JobCall{
			Name: "log",
			Args: []interface{}{"Mixed condition matched"},
		},
	}

	cp := &ConstraintPipeline{}

	// Créer l'AlphaNode avec l'expression mixte
	err := cp.createAlphaNodeWithTerminal(
		network,
		"rule_mixed_test",
		mixedExpr,
		"p",
		"Person",
		action,
		storage,
	)

	if err != nil {
		t.Fatalf("Erreur création AlphaNode Mixed: %v", err)
	}

	// Vérifier qu'un seul AlphaNode a été créé (pas de décomposition)
	alphaCount := len(network.AlphaNodes)
	if alphaCount != 1 {
		t.Errorf("Expression Mixed décomposée: attendu 1 AlphaNode, obtenu %d", alphaCount)
	}

	// Vérifier qu'un TerminalNode a été créé
	if len(network.TerminalNodes) != 1 {
		t.Errorf("Attendu 1 TerminalNode, obtenu %d", len(network.TerminalNodes))
	}

	t.Logf("✓ Expression mixte (AND+OR) crée un seul AlphaNode normalisé")
}

// TestOR_FactPropagation_Correct vérifie que les faits se propagent correctement
// à travers un AlphaNode contenant une expression OR
func TestOR_FactPropagation_Correct(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Créer un TypeNode pour Person
	typeDef := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "status", Type: "string"},
			{Name: "age", Type: "number"},
		},
	}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode

	// Expression OR: p.status == "VIP" OR p.age > 18
	orExpr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type: "binaryOperation",
			Left: constraint.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "status",
			},
			Operator: "==",
			Right: constraint.StringLiteral{
				Type:  "stringLiteral",
				Value: "VIP",
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
						Value: 18.0,
					},
				},
			},
		},
	}

	// Créer un fichier temporaire pour capturer les logs
	tmpFile, err := os.CreateTemp("", "or_test_*.log")
	if err != nil {
		t.Fatalf("Erreur création fichier temp: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	action := &Action{
		Type: "log",
		Job: JobCall{
			Name: "log_to_file",
			Args: []interface{}{tmpFile.Name(), "Match detected"},
		},
	}

	cp := &ConstraintPipeline{}

	// Créer l'AlphaNode avec l'expression OR
	err = cp.createAlphaNodeWithTerminal(
		network,
		"rule_or_propagation",
		orExpr,
		"p",
		"Person",
		action,
		storage,
	)

	if err != nil {
		t.Fatalf("Erreur création AlphaNode OR: %v", err)
	}

	// Créer des faits de test
	// Fait 1: satisfait la première condition (status == "VIP")
	fact1 := &Fact{
		ID:   "p1",
		Type: "Person",
		Fields: map[string]interface{}{
			"status": "VIP",
			"age":    15.0, // ne satisfait pas la deuxième condition
		},
	}

	// Fait 2: satisfait la deuxième condition (age > 18)
	fact2 := &Fact{
		ID:   "p2",
		Type: "Person",
		Fields: map[string]interface{}{
			"status": "Regular",
			"age":    25.0,
		},
	}

	// Fait 3: satisfait les deux conditions
	fact3 := &Fact{
		ID:   "p3",
		Type: "Person",
		Fields: map[string]interface{}{
			"status": "VIP",
			"age":    30.0,
		},
	}

	// Fait 4: ne satisfait aucune condition
	fact4 := &Fact{
		ID:   "p4",
		Type: "Person",
		Fields: map[string]interface{}{
			"status": "Regular",
			"age":    16.0,
		},
	}

	// Propager les faits
	typeNode.ActivateRight(fact1)
	typeNode.ActivateRight(fact2)
	typeNode.ActivateRight(fact3)
	typeNode.ActivateRight(fact4)

	// Obtenir l'AlphaNode créé
	var alphaNode *AlphaNode
	for _, node := range network.AlphaNodes {
		alphaNode = node
		break
	}

	if alphaNode == nil {
		t.Fatal("AlphaNode non trouvé")
	}

	// Vérifier les faits dans le working memory de l'AlphaNode
	facts := alphaNode.Memory.GetFacts()

	// Au moins les 3 premiers faits devraient être présents (satisfont au moins une condition OR)
	if len(facts) < 3 {
		t.Errorf("Propagation incorrecte: attendu au moins 3 faits, obtenu %d", len(facts))
		t.Logf("Faits présents dans l'AlphaNode:")
		for _, f := range facts {
			t.Logf("  - %v", f)
		}
	}

	// Vérifier que fact4 n'est PAS présent (ne satisfait aucune condition)
	fact4Present := false
	for _, f := range facts {
		if f.ID == "p4" {
			fact4Present = true
			break
		}
	}

	if fact4Present {
		t.Error("Fait 4 ne devrait pas être propagé (ne satisfait aucune condition OR)")
	}

	t.Logf("✓ Propagation correcte des faits à travers l'AlphaNode OR")
	t.Logf("  Faits propagés: %d sur 4", len(facts))
}

// TestOR_SharingBetweenRules vérifie que deux règles avec la même expression OR
// (mais dans un ordre différent) partagent le même AlphaNode après normalisation
func TestOR_SharingBetweenRules(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Créer un TypeNode pour Person
	typeDef := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "status", Type: "string"},
			{Name: "age", Type: "number"},
		},
	}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode

	// Règle 1: p.status == "VIP" OR p.age > 18
	orExpr1 := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type: "binaryOperation",
			Left: constraint.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "status",
			},
			Operator: "==",
			Right: constraint.StringLiteral{
				Type:  "stringLiteral",
				Value: "VIP",
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
						Value: 18.0,
					},
				},
			},
		},
	}

	// Règle 2: p.age > 18 OR p.status == "VIP" (ordre inversé)
	orExpr2 := constraint.LogicalExpression{
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
				Value: 18.0,
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
						Field:  "status",
					},
					Operator: "==",
					Right: constraint.StringLiteral{
						Type:  "stringLiteral",
						Value: "VIP",
					},
				},
			},
		},
	}

	action1 := &Action{
		Type: "log",
		Job: JobCall{
			Name: "log",
			Args: []interface{}{"Rule 1 matched"},
		},
	}

	action2 := &Action{
		Type: "log",
		Job: JobCall{
			Name: "log",
			Args: []interface{}{"Rule 2 matched"},
		},
	}

	cp := &ConstraintPipeline{}

	// Créer les deux règles
	err := cp.createAlphaNodeWithTerminal(
		network,
		"rule_or_1",
		orExpr1,
		"p",
		"Person",
		action1,
		storage,
	)
	if err != nil {
		t.Fatalf("Erreur création règle 1: %v", err)
	}

	err = cp.createAlphaNodeWithTerminal(
		network,
		"rule_or_2",
		orExpr2,
		"p",
		"Person",
		action2,
		storage,
	)
	if err != nil {
		t.Fatalf("Erreur création règle 2: %v", err)
	}

	// Vérifier qu'un seul AlphaNode a été créé (partage)
	alphaCount := len(network.AlphaNodes)
	if alphaCount != 1 {
		t.Errorf("Partage échoué: attendu 1 AlphaNode partagé, obtenu %d AlphaNodes", alphaCount)
	}

	// Vérifier que deux TerminalNodes ont été créés
	if len(network.TerminalNodes) != 2 {
		t.Errorf("Attendu 2 TerminalNodes, obtenu %d", len(network.TerminalNodes))
	}

	// Vérifier que l'AlphaNode a deux enfants (les deux TerminalNodes)
	var alphaNode *AlphaNode
	for _, node := range network.AlphaNodes {
		alphaNode = node
		break
	}

	if alphaNode == nil {
		t.Fatal("AlphaNode non trouvé")
	}

	if len(alphaNode.Children) != 2 {
		t.Errorf("AlphaNode devrait avoir 2 enfants, a %d", len(alphaNode.Children))
	}

	t.Logf("✓ Partage d'AlphaNode réussi entre règles OR avec ordre différent")
	t.Logf("  1 AlphaNode partagé -> 2 TerminalNodes")
}
