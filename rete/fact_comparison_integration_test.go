// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"

	"github.com/treivax/tsd/constraint"
)

func TestFactComparison_DirectVariableComparison(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION - COMPARAISON DIRECTE DE FAITS")
	t.Log("===================================================")

	// Créer les types
	types := []constraint.TypeDefinition{
		{
			Name: "User",
			Fields: []constraint.Field{
				{Name: "name", Type: "string", IsPrimaryKey: true},
				{Name: "age", Type: "number"},
			},
		},
		{
			Name: "Login",
			Fields: []constraint.Field{
				{Name: "user", Type: "User"},
				{Name: "email", Type: "string", IsPrimaryKey: true},
			},
		},
	}

	// Créer les faits
	userAlice := &Fact{
		ID:   "User~Alice",
		Type: "User",
		Fields: map[string]interface{}{
			"name": "Alice",
			"age":  30.0,
		},
	}

	userBob := &Fact{
		ID:   "User~Bob",
		Type: "User",
		Fields: map[string]interface{}{
			"name": "Bob",
			"age":  25.0,
		},
	}

	loginAlice := &Fact{
		ID:   "Login~alice@ex.com",
		Type: "Login",
		Fields: map[string]interface{}{
			"user":  "User~Alice", // Référence au User Alice
			"email": "alice@ex.com",
		},
	}

	loginBob := &Fact{
		ID:   "Login~bob@ex.com",
		Type: "Login",
		Fields: map[string]interface{}{
			"user":  "User~Bob", // Référence au User Bob
			"email": "bob@ex.com",
		},
	}

	// Créer les résolveurs
	resolver := NewFieldResolver(types)
	compEvaluator := NewComparisonEvaluator(resolver)

	// Créer l'évaluateur et configurer le contexte de types
	evaluator := NewAlphaConditionEvaluator()
	evaluator.SetTypeContext(resolver, compEvaluator)

	t.Run("comparaison variable == variable (même fait)", func(t *testing.T) {
		// Lier les variables
		evaluator.variableBindings["u"] = userAlice
		evaluator.variableBindings["u2"] = userAlice

		// Comparer u == u2 (doivent être égaux car c'est le même fait)
		left := userAlice
		right := userAlice
		result, err := evaluator.compareValues(left, "==", right)

		if err != nil {
			t.Fatalf("❌ Erreur inattendue: %v", err)
		}

		if !result {
			t.Errorf("❌ Attendu true (même fait), reçu false")
		} else {
			t.Logf("✅ Comparaison u == u2 correcte (même ID)")
		}

		evaluator.ClearBindings()
	})

	t.Run("comparaison variable == variable (faits différents)", func(t *testing.T) {
		// Lier les variables
		evaluator.variableBindings["u"] = userAlice
		evaluator.variableBindings["u2"] = userBob

		// Comparer u == u2 (doivent être différents)
		left := userAlice
		right := userBob
		result, err := evaluator.compareValues(left, "==", right)

		if err != nil {
			t.Fatalf("❌ Erreur inattendue: %v", err)
		}

		if result {
			t.Errorf("❌ Attendu false (faits différents), reçu true")
		} else {
			t.Logf("✅ Comparaison u == u2 correcte (IDs différents)")
		}

		evaluator.ClearBindings()
	})

	t.Run("comparaison variable != variable", func(t *testing.T) {
		// Lier les variables
		evaluator.variableBindings["u"] = userAlice
		evaluator.variableBindings["u2"] = userBob

		// Comparer u != u2 (doivent être différents)
		left := userAlice
		right := userBob
		result, err := evaluator.compareValues(left, "!=", right)

		if err != nil {
			t.Fatalf("❌ Erreur inattendue: %v", err)
		}

		if !result {
			t.Errorf("❌ Attendu true (faits différents), reçu false")
		} else {
			t.Logf("✅ Comparaison u != u2 correcte")
		}

		evaluator.ClearBindings()
	})

	t.Run("comparaison field access == variable", func(t *testing.T) {
		// Tester l.user == u
		evaluator.variableBindings["l"] = loginAlice
		evaluator.variableBindings["u"] = userAlice

		// Résoudre l.user pour obtenir l'ID du fait référencé
		userID, fieldType, err := resolver.ResolveFieldValue(loginAlice, "user")
		if err != nil {
			t.Fatalf("❌ Erreur lors de la résolution de l.user: %v", err)
		}

		if fieldType != FieldTypeFact {
			t.Fatalf("❌ Type attendu 'fact', reçu '%s'", fieldType)
		}

		// Comparer l'ID du champ avec l'ID du fait
		result, err := compEvaluator.EvaluateComparison(
			userID,
			userAlice.ID,
			"==",
			FieldTypeFact,
			FieldTypeFact,
		)

		if err != nil {
			t.Fatalf("❌ Erreur inattendue: %v", err)
		}

		if !result {
			t.Errorf("❌ Attendu true (l.user référence u), reçu false")
		} else {
			t.Logf("✅ Comparaison l.user == u correcte")
		}

		evaluator.ClearBindings()
	})

	t.Run("comparaison field access == variable (pas de match)", func(t *testing.T) {
		// Tester l.user == u où l référence Bob mais u est Alice
		evaluator.variableBindings["l"] = loginBob
		evaluator.variableBindings["u"] = userAlice

		// Résoudre l.user pour obtenir l'ID du fait référencé
		userID, fieldType, err := resolver.ResolveFieldValue(loginBob, "user")
		if err != nil {
			t.Fatalf("❌ Erreur lors de la résolution de l.user: %v", err)
		}

		// Comparer l'ID du champ avec l'ID du fait
		result, err := compEvaluator.EvaluateComparison(
			userID,
			userAlice.ID,
			"==",
			fieldType,
			FieldTypeFact,
		)

		if err != nil {
			t.Fatalf("❌ Erreur inattendue: %v", err)
		}

		if result {
			t.Errorf("❌ Attendu false (l.user référence Bob, pas Alice), reçu true")
		} else {
			t.Logf("✅ Comparaison l.user == u correcte (pas de match)")
		}

		evaluator.ClearBindings()
	})
}

func TestFactComparison_WithEvaluator(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION - ÉVALUATION COMPLÈTE")
	t.Log("=========================================")

	// Créer les types
	types := []constraint.TypeDefinition{
		{
			Name: "User",
			Fields: []constraint.Field{
				{Name: "name", Type: "string", IsPrimaryKey: true},
			},
		},
		{
			Name: "Login",
			Fields: []constraint.Field{
				{Name: "user", Type: "User"},
				{Name: "email", Type: "string", IsPrimaryKey: true},
			},
		},
	}

	// Créer les faits
	userAlice := &Fact{
		ID:   "User~Alice",
		Type: "User",
		Fields: map[string]interface{}{
			"name": "Alice",
		},
	}

	loginAlice := &Fact{
		ID:   "Login~alice@ex.com",
		Type: "Login",
		Fields: map[string]interface{}{
			"user":  "User~Alice",
			"email": "alice@ex.com",
		},
	}

	// Créer les résolveurs
	resolver := NewFieldResolver(types)
	compEvaluator := NewComparisonEvaluator(resolver)

	// Créer l'évaluateur et configurer le contexte
	evaluator := NewAlphaConditionEvaluator()
	evaluator.SetTypeContext(resolver, compEvaluator)

	// Lier les variables
	evaluator.variableBindings["u"] = userAlice
	evaluator.variableBindings["l"] = loginAlice

	t.Run("évaluation expression binaire u == u", func(t *testing.T) {
		// Créer une expression binaire pour comparer u avec lui-même
		expr := map[string]interface{}{
			"type":     "binaryOperation",
			"left":     constraint.Variable{Type: "variable", Name: "u"},
			"operator": "==",
			"right":    constraint.Variable{Type: "variable", Name: "u"},
		}

		result, err := evaluator.evaluateExpression(expr)
		if err != nil {
			t.Fatalf("❌ Erreur lors de l'évaluation: %v", err)
		}

		if !result {
			t.Errorf("❌ Attendu true (u == u), reçu false")
		} else {
			t.Logf("✅ Évaluation u == u correcte")
		}
	})

	evaluator.ClearBindings()
}
