// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"time"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/rete"
)

func main() {
	fmt.Println("=== Démonstration de la Normalisation des Conditions ===")

	// Exemple 1: AND Order Independence
	demonstrateANDNormalization()

	// Exemple 2: OR Order Independence
	demonstrateORNormalization()

	// Exemple 3: Non-Commutative Operations
	demonstrateNonCommutativeOperations()

	// Exemple 4: Complex Expression Normalization
	demonstrateComplexNormalization()

	// Exemple 5: Expression Reconstruction
	demonstrateExpressionReconstruction()

	// Exemple 6: Cache Performance
	demonstrateCachePerformance()
}

func demonstrateANDNormalization() {
	fmt.Println("📋 Exemple 1: Normalisation AND (opérateur commutatif)")
	fmt.Println("=" + repeat("=", 60))

	// Créer deux conditions: age > 18 ET salary >= 50000
	condAge := rete.NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	)

	condSalary := rete.NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
		">=",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
	)

	// Tester les deux ordres
	fmt.Println("\n🔄 Ordre A: age > 18 AND salary >= 50000")
	conditionsAB := []rete.SimpleCondition{condAge, condSalary}
	normalizedAB := rete.NormalizeConditions(conditionsAB, "AND")
	printConditions(normalizedAB)

	fmt.Println("\n🔄 Ordre B: salary >= 50000 AND age > 18")
	conditionsBA := []rete.SimpleCondition{condSalary, condAge}
	normalizedBA := rete.NormalizeConditions(conditionsBA, "AND")
	printConditions(normalizedBA)

	// Vérifier l'équivalence
	fmt.Println("\n✅ Vérification:")
	if areConditionsEqual(normalizedAB, normalizedBA) {
		fmt.Println("   Les deux ordres produisent le MÊME ordre canonique!")
		fmt.Println("   Ordre canonique:", rete.CanonicalString(normalizedAB[0]))
	} else {
		fmt.Println("   ❌ ERREUR: Les ordres sont différents")
	}
	fmt.Println()
}

func demonstrateORNormalization() {
	fmt.Println("📋 Exemple 2: Normalisation OR (opérateur commutatif)")
	fmt.Println("=" + repeat("=", 60))

	// Créer deux conditions: status == "active" OU verified == true
	condStatus := rete.NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "user", Field: "status"},
		"==",
		constraint.StringLiteral{Type: "stringLiteral", Value: "active"},
	)

	condVerified := rete.NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "user", Field: "verified"},
		"==",
		constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
	)

	fmt.Println("\n🔄 Ordre A: status == 'active' OR verified == true")
	conditionsAB := []rete.SimpleCondition{condStatus, condVerified}
	normalizedAB := rete.NormalizeConditions(conditionsAB, "OR")
	printConditions(normalizedAB)

	fmt.Println("\n🔄 Ordre B: verified == true OR status == 'active'")
	conditionsBA := []rete.SimpleCondition{condVerified, condStatus}
	normalizedBA := rete.NormalizeConditions(conditionsBA, "OR")
	printConditions(normalizedBA)

	fmt.Println("\n✅ Vérification:")
	if areConditionsEqual(normalizedAB, normalizedBA) {
		fmt.Println("   Les deux ordres produisent le MÊME ordre canonique!")
	}
	fmt.Println()
}

func demonstrateNonCommutativeOperations() {
	fmt.Println("📋 Exemple 3: Opérations Non-Commutatives (préservation de l'ordre)")
	fmt.Println("=" + repeat("=", 60))

	// Créer deux conditions avec soustraction (non-commutatif)
	cond1 := rete.NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "x", Field: "balance"},
		"-",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 100},
	)

	cond2 := rete.NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "x", Field: "fee"},
		"-",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 10},
	)

	fmt.Println("\n🔒 Opérateur '-' est NON-COMMUTATIF")
	fmt.Printf("   IsCommutative('-') = %v\n", rete.IsCommutative("-"))

	fmt.Println("\n📌 Ordre original: [cond1, cond2]")
	original := []rete.SimpleCondition{cond1, cond2}
	printConditions(original)

	fmt.Println("\n🔒 Après normalisation avec opérateur SEQ (non-commutatif):")
	normalized := rete.NormalizeConditions(original, "SEQ")
	printConditions(normalized)

	fmt.Println("\n✅ Vérification:")
	if areConditionsEqual(original, normalized) {
		fmt.Println("   L'ordre original est PRÉSERVÉ (pas de tri)")
	}
	fmt.Println()
}

func demonstrateComplexNormalization() {
	fmt.Println("📋 Exemple 4: Normalisation d'Expressions Complexes")
	fmt.Println("=" + repeat("=", 60))

	// Expression: age > 18 AND salary >= 50000
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

	fmt.Println("\n🔍 Expression originale:")
	fmt.Println("   (p.age > 18) AND (p.salary >= 50000)")

	// Extraire les conditions
	conditions, opType, err := rete.ExtractConditions(expr)
	if err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
		return
	}

	fmt.Printf("\n📊 Conditions extraites (opérateur: %s):\n", opType)
	printConditions(conditions)

	// Normaliser
	normalized := rete.NormalizeConditions(conditions, opType)
	fmt.Println("\n✨ Conditions normalisées:")
	printConditions(normalized)

	// Normaliser l'expression complète
	normalizedExpr, err := rete.NormalizeExpression(expr)
	if err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
		return
	}

	fmt.Println("\n✅ Expression normalisée avec succès!")
	fmt.Printf("   Type: %T\n", normalizedExpr)
	fmt.Println()
}

func demonstrateExpressionReconstruction() {
	fmt.Println("📋 Exemple 5: Reconstruction d'Expressions Normalisées")
	fmt.Println("=" + repeat("=", 60))

	// Expression originale : salary >= 50000 AND age > 18 (ordre inversé)
	fmt.Println("\n🔍 Expression originale (ordre inversé):")
	fmt.Println("   (p.salary >= 50000) AND (p.age > 18)")

	expr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
			Operator: ">=",
			Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
					Operator: ">",
					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
				},
			},
		},
	}

	// Extraire les conditions AVANT normalisation
	fmt.Println("\n📊 Conditions AVANT normalisation:")
	condsBefore, _, _ := rete.ExtractConditions(expr)
	printConditions(condsBefore)

	// Normaliser l'expression (reconstruction automatique)
	fmt.Println("\n✨ Normalisation avec RECONSTRUCTION automatique...")
	normalized, err := rete.NormalizeExpression(expr)
	if err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
		return
	}

	// Vérifier que c'est bien une LogicalExpression reconstruite
	normExpr, ok := normalized.(constraint.LogicalExpression)
	if !ok {
		fmt.Printf("❌ Type inattendu: %T\n", normalized)
		return
	}

	// Extraire les conditions APRÈS normalisation
	fmt.Println("\n📊 Conditions APRÈS normalisation et reconstruction:")
	condsAfter, _, _ := rete.ExtractConditions(normExpr)
	printConditions(condsAfter)

	// Vérifier l'ordre canonique
	fmt.Println("\n🔍 Vérification de l'ordre canonique:")

	// Left devrait être age > 18 (vient avant salary alphabétiquement)
	leftOp, ok := normExpr.Left.(constraint.BinaryOperation)
	if ok {
		leftField, ok := leftOp.Left.(constraint.FieldAccess)
		if ok {
			fmt.Printf("   ✓ Premier élément (Left): p.%s %s ...\n", leftField.Field, leftOp.Operator)
			if leftField.Field == "age" {
				fmt.Println("     ✅ Correct ! 'age' vient avant 'salary' en ordre canonique")
			} else {
				fmt.Println("     ❌ Attendu: 'age'")
			}
		}
	}

	// Right devrait être salary >= 50000
	if len(normExpr.Operations) > 0 {
		rightOp, ok := normExpr.Operations[0].Right.(constraint.BinaryOperation)
		if ok {
			rightField, ok := rightOp.Left.(constraint.FieldAccess)
			if ok {
				fmt.Printf("   ✓ Deuxième élément (Operations[0]): p.%s %s ...\n", rightField.Field, rightOp.Operator)
				if rightField.Field == "salary" {
					fmt.Println("     ✅ Correct ! 'salary' vient après 'age' en ordre canonique")
				} else {
					fmt.Println("     ❌ Attendu: 'salary'")
				}
			}
		}
	}

	// Démonstration avec deux expressions différentes
	fmt.Println("\n🔄 Démonstration: Deux ordres différents → Même structure reconstruite")

	// Expression 2 : age > 18 AND salary >= 50000 (ordre normal)
	expr2 := constraint.LogicalExpression{
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

	// Normaliser la deuxième expression
	normalized2, _ := rete.NormalizeExpression(expr2)
	normExpr2, _ := normalized2.(constraint.LogicalExpression)

	// Extraire et comparer
	conds1, _, _ := rete.ExtractConditions(normExpr)
	conds2, _, _ := rete.ExtractConditions(normExpr2)

	fmt.Println("\n   Expression 1 (inversée) après normalisation:")
	for i, c := range conds1 {
		canonical := rete.CanonicalString(c)
		if len(canonical) > 50 {
			canonical = canonical[:50] + "..."
		}
		fmt.Printf("     [%d] %s\n", i, canonical)
	}

	fmt.Println("\n   Expression 2 (normale) après normalisation:")
	for i, c := range conds2 {
		canonical := rete.CanonicalString(c)
		if len(canonical) > 50 {
			canonical = canonical[:50] + "..."
		}
		fmt.Printf("     [%d] %s\n", i, canonical)
	}

	// Vérifier l'égalité
	allEqual := true
	if len(conds1) == len(conds2) {
		for i := range conds1 {
			if !rete.CompareConditions(conds1[i], conds2[i]) {
				allEqual = false
				break
			}
		}
	} else {
		allEqual = false
	}

	fmt.Println("\n✅ Résultat:")
	if allEqual {
		fmt.Println("   🎉 Les deux expressions ont été reconstruites avec le MÊME ordre canonique!")
		fmt.Println("   → Le partage de nœuds Alpha sera maximal")
	} else {
		fmt.Println("   ❌ Les ordres diffèrent")
	}

	fmt.Println()
}

func demonstrateCachePerformance() {
	fmt.Println("📋 Exemple 6: Cache de Normalisation (Performance)")
	fmt.Println("=" + repeat("=", 60))

	// Créer un cache
	cache := rete.NewNormalizationCache(100)

	fmt.Println("\n🔧 Configuration du cache:")
	fmt.Printf("   Taille max: %d entrées\n", 100)
	fmt.Printf("   Stratégie d'éviction: LRU\n")
	fmt.Printf("   Status: %v\n", cache.IsEnabled())

	// Créer une expression à normaliser
	expr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
			Operator: ">=",
			Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
					Operator: ">",
					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
				},
			},
		},
	}

	fmt.Println("\n🔄 Test 1: Première normalisation (cache MISS)")
	_, _ = rete.NormalizeExpressionWithCache(expr, cache)
	stats := cache.GetStats()
	fmt.Printf("   Hits: %d, Misses: %d, Hit Rate: %.1f%%\n", stats.Hits, stats.Misses, stats.HitRate*100)

	fmt.Println("\n🔄 Test 2: Deuxième normalisation (cache HIT)")
	_, _ = rete.NormalizeExpressionWithCache(expr, cache)
	stats = cache.GetStats()
	fmt.Printf("   Hits: %d, Misses: %d, Hit Rate: %.1f%%\n", stats.Hits, stats.Misses, stats.HitRate*100)

	fmt.Println("\n🔄 Test 3: Normalisation répétée (10x)")
	for i := 0; i < 10; i++ {
		_, _ = rete.NormalizeExpressionWithCache(expr, cache)
	}
	stats = cache.GetStats()
	fmt.Printf("   Hits: %d, Misses: %d, Hit Rate: %.1f%%\n", stats.Hits, stats.Misses, stats.HitRate*100)

	fmt.Println("\n📊 Statistiques finales du cache:")
	fmt.Printf("   %s\n", stats.String())

	// Benchmark simple
	fmt.Println("\n⚡ Benchmark de performance (1000 itérations):")

	iterations := 1000

	// Sans cache
	start := timeNow()
	for i := 0; i < iterations; i++ {
		_, _ = rete.NormalizeExpression(expr)
	}
	durationNoCache := timeSince(start)

	// Avec cache (nouveau cache pour reset)
	cacheNew := rete.NewNormalizationCache(100)
	start = timeNow()
	for i := 0; i < iterations; i++ {
		_, _ = rete.NormalizeExpressionWithCache(expr, cacheNew)
	}
	durationWithCache := timeSince(start)

	fmt.Printf("   Sans cache:  %v\n", durationNoCache)
	fmt.Printf("   Avec cache:  %v\n", durationWithCache)

	if durationNoCache > durationWithCache {
		speedup := float64(durationNoCache) / float64(durationWithCache)
		fmt.Printf("   ⚡ Speedup:   %.2fx plus rapide!\n", speedup)
	}

	finalStats := cacheNew.GetStats()
	fmt.Printf("\n   Cache final: %d hits, %d miss, taux de succès %.1f%%\n",
		finalStats.Hits, finalStats.Misses, finalStats.HitRate*100)

	fmt.Println("\n✅ Conclusion:")
	fmt.Println("   Le cache améliore significativement les performances")
	fmt.Println("   pour les expressions normalisées fréquemment utilisées.")
	fmt.Println()
}

// Fonctions utilitaires

func timeNow() time.Time {
	return time.Now()
}

func timeSince(start time.Time) time.Duration {
	return time.Since(start)
}

func printConditions(conditions []rete.SimpleCondition) {
	for i, cond := range conditions {
		canonical := rete.CanonicalString(cond)
		fmt.Printf("   [%d] %s\n", i, canonical)
		fmt.Printf("       Hash: %s\n", cond.Hash[:16]+"...")
	}
}

func areConditionsEqual(c1, c2 []rete.SimpleCondition) bool {
	if len(c1) != len(c2) {
		return false
	}
	for i := range c1 {
		if !rete.CompareConditions(c1[i], c2[i]) {
			return false
		}
	}
	return true
}

func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
