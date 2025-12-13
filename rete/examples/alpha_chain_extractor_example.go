// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

//go:build ignore
// +build ignore

// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"strings"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/rete"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("   Alpha Chain Extractor - Exemple")
	fmt.Println("========================================\n")

	// Exemple 1: Comparaison simple
	example1()

	// Exemple 2: Expression AND
	example2()

	// Exemple 3: Expression imbriquée complexe
	example3()

	// Exemple 4: Détection de partage de conditions
	example4()
}

// example1 démontre l'extraction d'une comparaison simple
func example1() {
	fmt.Println("📋 Exemple 1: Comparaison simple")
	fmt.Println("Expression: p.age > 18\n")

	expr := constraint.BinaryOperation{
		Type:     "binaryOperation",
		Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		Operator: ">",
		Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	}

	conditions, opType, err := rete.ExtractConditions(expr)
	if err != nil {
		fmt.Printf("❌ Erreur: %v\n\n", err)
		return
	}

	fmt.Printf("✅ Type d'opérateur: %s\n", opType)
	fmt.Printf("✅ Nombre de conditions: %d\n", len(conditions))

	for i, cond := range conditions {
		fmt.Printf("\nCondition %d:\n", i+1)
		fmt.Printf("  Type: %s\n", cond.Type)
		fmt.Printf("  Opérateur: %s\n", cond.Operator)
		fmt.Printf("  Canonique: %s\n", rete.CanonicalString(cond))
		fmt.Printf("  Hash: %s...\n", cond.Hash[:16])
	}

	fmt.Println("\n" + strings.Repeat("-", 60) + "\n")
}

// example2 démontre l'extraction d'une expression AND
func example2() {
	fmt.Println("📋 Exemple 2: Expression AND")
	fmt.Println("Expression: p.age > 18 AND p.salary >= 50000\n")

	expr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
			Operator: ">",
			Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
					Operator: ">=",
					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
				},
			},
		},
	}

	conditions, opType, err := rete.ExtractConditions(expr)
	if err != nil {
		fmt.Printf("❌ Erreur: %v\n\n", err)
		return
	}

	fmt.Printf("✅ Type d'opérateur: %s\n", opType)
	fmt.Printf("✅ Nombre de conditions: %d\n", len(conditions))

	for i, cond := range conditions {
		fmt.Printf("\nCondition %d:\n", i+1)
		canonical := rete.CanonicalString(cond)
		fmt.Printf("  Canonique: %s\n", canonical)
		fmt.Printf("  Hash: %s...\n", cond.Hash[:16])
	}

	fmt.Println("\n" + strings.Repeat("-", 60) + "\n")
}

// example3 démontre l'extraction d'une expression imbriquée complexe
func example3() {
	fmt.Println("📋 Exemple 3: Expression imbriquée complexe")
	fmt.Println("Expression: (p.age > 18 AND p.salary >= 50000) OR p.vip == true\n")

	// Expression interne: p.age > 18 AND p.salary >= 50000
	innerExpr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
			Operator: ">",
			Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
					Operator: ">=",
					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
				},
			},
		},
	}

	// Expression globale: (inner) OR p.vip == true
	expr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: innerExpr,
		Operations: []constraint.LogicalOperation{
			{
				Op: "OR",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "vip"},
					Operator: "==",
					Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
				},
			},
		},
	}

	conditions, opType, err := rete.ExtractConditions(expr)
	if err != nil {
		fmt.Printf("❌ Erreur: %v\n\n", err)
		return
	}

	fmt.Printf("✅ Type d'opérateur: %s\n", opType)
	fmt.Printf("✅ Nombre de conditions: %d\n", len(conditions))

	for i, cond := range conditions {
		fmt.Printf("\nCondition %d:\n", i+1)
		canonical := rete.CanonicalString(cond)
		fmt.Printf("  Canonique: %s\n", canonical)
		fmt.Printf("  Hash: %s...\n", cond.Hash[:16])
	}

	// Tester la déduplication
	fmt.Println("\n📊 Test de déduplication:")
	duplicated := append(conditions, conditions[0]) // Ajouter un doublon
	fmt.Printf("  Avant: %d conditions\n", len(duplicated))
	unique := rete.DeduplicateConditions(duplicated)
	fmt.Printf("  Après: %d conditions\n", len(unique))

	fmt.Println("\n" + strings.Repeat("-", 60) + "\n")
}

// example4 démontre la détection de partage de conditions entre règles
func example4() {
	fmt.Println("📋 Exemple 4: Détection de partage de conditions")
	fmt.Println("Règles:")
	fmt.Println("  - Règle 1: p.age > 18 AND p.active == true")
	fmt.Println("  - Règle 2: p.age > 18 AND p.salary >= 30000")
	fmt.Println("  - Règle 3: p.active == true AND p.department == 'Sales'\n")

	// Règle 1: p.age > 18 AND p.active == true
	rule1 := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
			Operator: ">",
			Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "active"},
					Operator: "==",
					Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
				},
			},
		},
	}

	// Règle 2: p.age > 18 AND p.salary >= 30000
	rule2 := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
			Operator: ">",
			Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
					Operator: ">=",
					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 30000},
				},
			},
		},
	}

	// Règle 3: p.active == true AND p.department == 'Sales'
	rule3 := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "active"},
			Operator: "==",
			Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "department"},
					Operator: "==",
					Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "Sales"},
				},
			},
		},
	}

	// Extraire toutes les conditions
	conds1, _, _ := rete.ExtractConditions(rule1)
	conds2, _, _ := rete.ExtractConditions(rule2)
	conds3, _, _ := rete.ExtractConditions(rule3)

	// Analyser le partage
	conditionUsage := make(map[string][]string)

	for _, cond := range conds1 {
		conditionUsage[cond.Hash] = append(conditionUsage[cond.Hash], "Règle 1")
	}
	for _, cond := range conds2 {
		conditionUsage[cond.Hash] = append(conditionUsage[cond.Hash], "Règle 2")
	}
	for _, cond := range conds3 {
		conditionUsage[cond.Hash] = append(conditionUsage[cond.Hash], "Règle 3")
	}

	// Afficher les résultats
	fmt.Println("📊 Analyse de partage:")
	sharedCount := 0
	for hash, rules := range conditionUsage {
		if len(rules) > 1 {
			sharedCount++
			fmt.Printf("\n✅ Condition partagée (hash: %s...):\n", hash[:16])
			fmt.Printf("   Utilisée par: %v\n", rules)

			// Trouver la condition pour afficher sa forme canonique
			for _, cond := range append(append(conds1, conds2...), conds3...) {
				if cond.Hash == hash {
					fmt.Printf("   Canonique: %s\n", rete.CanonicalString(cond))
					break
				}
			}
		}
	}

	totalConditions := len(conds1) + len(conds2) + len(conds3)
	uniqueConditions := len(conditionUsage)
	saved := totalConditions - uniqueConditions

	fmt.Printf("\n📈 Statistiques:\n")
	fmt.Printf("  Total de conditions: %d\n", totalConditions)
	fmt.Printf("  Conditions uniques: %d\n", uniqueConditions)
	fmt.Printf("  Conditions partagées: %d\n", sharedCount)
	fmt.Printf("  Économie potentielle: %d nœuds alpha\n", saved)

	fmt.Println("\n" + strings.Repeat("=", 60))
}
