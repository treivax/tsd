//go:build ignore
// +build ignore

// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"log"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/rete"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("🔗 Constraint Pipeline Chain Decomposition Demo")
	fmt.Println("========================================\n")

	// Créer le réseau RETE
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)
	network.AlphaSharingManager = rete.NewAlphaSharingRegistry()
	network.LifecycleManager = rete.NewLifecycleManager()

	// Définir le type Person
	typeDef := rete.TypeDefinition{
		Type: "type",
		Name: "Person",
		Fields: []rete.Field{
			{Name: "age", Type: "number"},
			{Name: "salary", Type: "number"},
			{Name: "experience", Type: "number"},
			{Name: "performance", Type: "number"},
		},
	}

	// Créer le TypeNode
	typeNode := rete.NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode

	cp := &rete.ConstraintPipeline{}

	fmt.Println("📋 Exemple 1: Condition Simple (pas de décomposition)")
	fmt.Println("───────────────────────────────────────────────────")

	// Condition simple: p.age > 18
	simpleCondition := map[string]interface{}{
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

	action1 := &rete.Action{
		Type: "print",
		Job: rete.JobCall{
			Name: "check_age",
			Args: []interface{}{"Person is adult"},
		},
	}

	err := cp.createAlphaNodeWithTerminal(network, "rule_simple", simpleCondition, "p", "Person", action1, storage)
	if err != nil {
		log.Fatalf("Erreur règle simple: %v", err)
	}

	fmt.Printf("\n✅ Résultat: %d AlphaNode(s) créé(s)\n\n", len(network.AlphaNodes))

	// Réinitialiser pour le prochain exemple
	network.AlphaNodes = make(map[string]*rete.AlphaNode)
	network.TerminalNodes = make(map[string]*rete.TerminalNode)

	fmt.Println("📋 Exemple 2: Expression AND (décomposition en chaîne)")
	fmt.Println("─────────────────────────────────────────────────────")

	// Expression AND: p.age > 18 AND p.salary >= 50000
	andCondition := constraint.LogicalExpression{
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

	action2 := &rete.Action{
		Type: "print",
		Job: rete.JobCall{
			Name: "eligible_hire",
			Args: []interface{}{"Person eligible for hiring"},
		},
	}

	err = cp.createAlphaNodeWithTerminal(network, "rule_and", andCondition, "p", "Person", action2, storage)
	if err != nil {
		log.Fatalf("Erreur règle AND: %v", err)
	}

	fmt.Printf("\n✅ Résultat: %d AlphaNode(s) créé(s) (chaîne de 2 nœuds)\n\n", len(network.AlphaNodes))

	fmt.Println("📋 Exemple 3: Deux Règles avec Partage")
	fmt.Println("───────────────────────────────────────")

	// Même expression AND pour une deuxième règle
	action3 := &rete.Action{
		Type: "notify",
		Job: rete.JobCall{
			Name: "send_offer",
			Args: []interface{}{"Send job offer"},
		},
	}

	err = cp.createAlphaNodeWithTerminal(network, "rule_and_2", andCondition, "p", "Person", action3, storage)
	if err != nil {
		log.Fatalf("Erreur règle AND 2: %v", err)
	}

	fmt.Printf("\n✅ Résultat: %d AlphaNode(s) total (aucun nouveau nœud, réutilisation complète!)\n", len(network.AlphaNodes))
	fmt.Printf("✅ %d TerminalNode(s) créé(s) (un par règle)\n\n", len(network.TerminalNodes))

	fmt.Println("📋 Exemple 4: Expression AND Complexe (3 conditions)")
	fmt.Println("────────────────────────────────────────────────────")

	// Réinitialiser
	network.AlphaNodes = make(map[string]*rete.AlphaNode)
	network.TerminalNodes = make(map[string]*rete.TerminalNode)

	// Expression AND: p.age > 18 AND p.salary >= 50000 AND p.experience > 5
	complexAndCondition := constraint.LogicalExpression{
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

	action4 := &rete.Action{
		Type: "print",
		Job: rete.JobCall{
			Name: "senior_candidate",
			Args: []interface{}{"Senior candidate identified"},
		},
	}

	err = cp.createAlphaNodeWithTerminal(network, "rule_complex", complexAndCondition, "p", "Person", action4, storage)
	if err != nil {
		log.Fatalf("Erreur règle complexe: %v", err)
	}

	fmt.Printf("\n✅ Résultat: %d AlphaNode(s) créé(s) (chaîne de 3 nœuds)\n\n", len(network.AlphaNodes))

	fmt.Println("📋 Exemple 5: Expression OR (pas de décomposition)")
	fmt.Println("──────────────────────────────────────────────────")

	// Réinitialiser
	network.AlphaNodes = make(map[string]*rete.AlphaNode)
	network.TerminalNodes = make(map[string]*rete.TerminalNode)

	// Expression OR: p.age < 18 OR p.age > 65
	orCondition := constraint.LogicalExpression{
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

	action5 := &rete.Action{
		Type: "print",
		Job: rete.JobCall{
			Name: "age_exception",
			Args: []interface{}{"Age exception detected"},
		},
	}

	err = cp.createAlphaNodeWithTerminal(network, "rule_or", orCondition, "p", "Person", action5, storage)
	if err != nil {
		log.Fatalf("Erreur règle OR: %v", err)
	}

	fmt.Printf("\n✅ Résultat: %d AlphaNode(s) créé(s) (un seul nœud normalisé, pas de chaîne)\n\n", len(network.AlphaNodes))

	fmt.Println("📊 Statistiques Globales")
	fmt.Println("────────────────────────")
	fmt.Println("✅ Tous les exemples exécutés avec succès!")
	fmt.Println("\n🎯 Points Clés:")
	fmt.Println("   • Conditions simples → 1 nœud alpha")
	fmt.Println("   • Expressions AND → Chaîne de nœuds (1 nœud par condition)")
	fmt.Println("   • Expressions OR → 1 nœud normalisé (pas de chaîne)")
	fmt.Println("   • Partage automatique → Nœuds identiques réutilisés entre règles")
	fmt.Println("   • Fallback robuste → Erreurs gérées gracieusement")
	fmt.Println("\n📚 Documentation complète dans:")
	fmt.Println("   tsd/rete/docs/CONSTRAINT_PIPELINE_CHAIN_DECOMPOSITION.md")
	fmt.Println("\n========================================")
	fmt.Println("✅ Demo terminée avec succès!")
	fmt.Println("========================================")
}
