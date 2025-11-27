// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"

	"github.com/treivax/tsd/constraint"
)

// TestExtractConditions_SimpleComparison teste l'extraction d'une comparaison simple
func TestExtractConditions_SimpleComparison(t *testing.T) {
	// Test avec BinaryOperation
	expr := constraint.BinaryOperation{
		Type:     "binaryOperation",
		Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		Operator: ">",
		Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	}

	conditions, opType, err := ExtractConditions(expr)

	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}

	if opType != "SINGLE" {
		t.Errorf("Attendu opType='SINGLE', obtenu '%s'", opType)
	}

	if len(conditions) != 1 {
		t.Fatalf("Attendu 1 condition, obtenu %d", len(conditions))
	}

	cond := conditions[0]
	if cond.Type != "binaryOperation" {
		t.Errorf("Attendu type='binaryOperation', obtenu '%s'", cond.Type)
	}
	if cond.Operator != ">" {
		t.Errorf("Attendu operator='>', obtenu '%s'", cond.Operator)
	}
	if cond.Hash == "" {
		t.Error("Hash ne doit pas être vide")
	}
}

// TestExtractConditions_SimpleComparison_Map teste l'extraction d'une comparaison (format map)
func TestExtractConditions_SimpleComparison_Map(t *testing.T) {
	expr := map[string]interface{}{
		"type":     "binaryOperation",
		"operator": "==",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "p",
			"field":  "name",
		},
		"right": map[string]interface{}{
			"type":  "stringLiteral",
			"value": "Alice",
		},
	}

	conditions, opType, err := ExtractConditions(expr)

	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}

	if opType != "SINGLE" {
		t.Errorf("Attendu opType='SINGLE', obtenu '%s'", opType)
	}

	if len(conditions) != 1 {
		t.Fatalf("Attendu 1 condition, obtenu %d", len(conditions))
	}

	cond := conditions[0]
	if cond.Operator != "==" {
		t.Errorf("Attendu operator='==', obtenu '%s'", cond.Operator)
	}
}

// TestExtractConditions_LogicalAND teste l'extraction d'une expression AND
func TestExtractConditions_LogicalAND(t *testing.T) {
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

	conditions, opType, err := ExtractConditions(expr)

	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}

	if opType != "AND" {
		t.Errorf("Attendu opType='AND', obtenu '%s'", opType)
	}

	if len(conditions) != 2 {
		t.Fatalf("Attendu 2 conditions, obtenu %d", len(conditions))
	}

	// Vérifier la première condition
	if conditions[0].Operator != ">" {
		t.Errorf("Condition[0]: attendu operator='>', obtenu '%s'", conditions[0].Operator)
	}

	// Vérifier la deuxième condition
	if conditions[1].Operator != ">=" {
		t.Errorf("Condition[1]: attendu operator='>=', obtenu '%s'", conditions[1].Operator)
	}

	// Vérifier que chaque condition a un hash unique
	if conditions[0].Hash == conditions[1].Hash {
		t.Error("Les deux conditions ne doivent pas avoir le même hash")
	}
}

// TestExtractConditions_LogicalOR teste l'extraction d'une expression OR
func TestExtractConditions_LogicalOR(t *testing.T) {
	expr := map[string]interface{}{
		"type": "logicalExpression",
		"left": map[string]interface{}{
			"type":     "binaryOperation",
			"operator": "<",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "age",
			},
			"right": map[string]interface{}{
				"type":  "numberLiteral",
				"value": 25,
			},
		},
		"operations": []interface{}{
			map[string]interface{}{
				"op": "OR",
				"right": map[string]interface{}{
					"type":     "binaryOperation",
					"operator": ">",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "p",
						"field":  "age",
					},
					"right": map[string]interface{}{
						"type":  "numberLiteral",
						"value": 65,
					},
				},
			},
		},
	}

	conditions, opType, err := ExtractConditions(expr)

	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}

	if opType != "OR" {
		t.Errorf("Attendu opType='OR', obtenu '%s'", opType)
	}

	if len(conditions) != 2 {
		t.Fatalf("Attendu 2 conditions, obtenu %d", len(conditions))
	}
}

// TestExtractConditions_NestedExpressions teste l'extraction d'expressions imbriquées
func TestExtractConditions_NestedExpressions(t *testing.T) {
	// (p.age > 18 AND p.salary >= 50000) AND p.active == true
	innerLogical := constraint.LogicalExpression{
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

	expr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: innerLogical,
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

	conditions, opType, err := ExtractConditions(expr)

	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}

	if opType != "AND" {
		t.Errorf("Attendu opType='AND', obtenu '%s'", opType)
	}

	if len(conditions) != 3 {
		t.Fatalf("Attendu 3 conditions, obtenu %d", len(conditions))
	}

	// Vérifier que toutes les conditions ont des hash uniques
	hashSet := make(map[string]bool)
	for _, cond := range conditions {
		if hashSet[cond.Hash] {
			t.Errorf("Hash dupliqué trouvé: %s", cond.Hash)
		}
		hashSet[cond.Hash] = true
	}
}

// TestExtractConditions_MixedOperators teste des opérateurs mélangés (AND et OR)
func TestExtractConditions_MixedOperators(t *testing.T) {
	expr := map[string]interface{}{
		"type": "logicalExpression",
		"left": map[string]interface{}{
			"type":     "binaryOperation",
			"operator": ">",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "age",
			},
			"right": map[string]interface{}{
				"type":  "numberLiteral",
				"value": 18,
			},
		},
		"operations": []interface{}{
			map[string]interface{}{
				"op": "AND",
				"right": map[string]interface{}{
					"type":     "binaryOperation",
					"operator": "<",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "p",
						"field":  "age",
					},
					"right": map[string]interface{}{
						"type":  "numberLiteral",
						"value": 65,
					},
				},
			},
			map[string]interface{}{
				"op": "OR",
				"right": map[string]interface{}{
					"type":     "binaryOperation",
					"operator": "==",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "p",
						"field":  "vip",
					},
					"right": map[string]interface{}{
						"type":  "booleanLiteral",
						"value": true,
					},
				},
			},
		},
	}

	conditions, opType, err := ExtractConditions(expr)

	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}

	if opType != "MIXED" {
		t.Errorf("Attendu opType='MIXED', obtenu '%s'", opType)
	}

	if len(conditions) != 3 {
		t.Fatalf("Attendu 3 conditions, obtenu %d", len(conditions))
	}
}

// TestExtractConditions_ArithmeticOperations teste l'extraction d'opérations arithmétiques (via map)
func TestExtractConditions_ArithmeticOperations(t *testing.T) {
	// Les opérations arithmétiques sont traitées comme des BinaryOperations
	expr := constraint.BinaryOperation{
		Type:     "binaryOperation",
		Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
		Operator: "+",
		Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 100},
	}

	conditions, opType, err := ExtractConditions(expr)

	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}

	if opType != "SINGLE" {
		t.Errorf("Attendu opType='SINGLE', obtenu '%s'", opType)
	}

	if len(conditions) != 1 {
		t.Fatalf("Attendu 1 condition, obtenu %d", len(conditions))
	}

	cond := conditions[0]
	if cond.Type != "binaryOperation" {
		t.Errorf("Attendu type='binaryOperation', obtenu '%s'", cond.Type)
	}
	if cond.Operator != "+" {
		t.Errorf("Attendu operator='+', obtenu '%s'", cond.Operator)
	}
}

// TestExtractConditions_ArithmeticInComparison teste une opération arithmétique dans une comparaison
func TestExtractConditions_ArithmeticInComparison(t *testing.T) {
	expr := constraint.BinaryOperation{
		Type: "binaryOperation",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
			Operator: "+",
			Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 1000},
		},
		Operator: ">",
		Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 60000},
	}

	conditions, opType, err := ExtractConditions(expr)

	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}

	if opType != "SINGLE" {
		t.Errorf("Attendu opType='SINGLE', obtenu '%s'", opType)
	}

	if len(conditions) != 1 {
		t.Fatalf("Attendu 1 condition, obtenu %d", len(conditions))
	}

	// La condition contient l'opération arithmétique dans son left
	cond := conditions[0]
	if cond.Type != "binaryOperation" {
		t.Errorf("Attendu type='binaryOperation', obtenu '%s'", cond.Type)
	}
}

// TestCanonicalString_Deterministic teste que CanonicalString est déterministe
func TestCanonicalString_Deterministic(t *testing.T) {
	// Créer la même condition deux fois
	cond1 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	)

	cond2 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	)

	str1 := CanonicalString(cond1)
	str2 := CanonicalString(cond2)

	if str1 != str2 {
		t.Errorf("CanonicalString doit être déterministe:\n  str1='%s'\n  str2='%s'", str1, str2)
	}

	// Les hash doivent aussi être identiques
	if cond1.Hash != cond2.Hash {
		t.Errorf("Hash doit être identique pour des conditions identiques:\n  hash1='%s'\n  hash2='%s'", cond1.Hash, cond2.Hash)
	}
}

// TestCanonicalString_Uniqueness teste que CanonicalString génère des strings uniques
func TestCanonicalString_Uniqueness(t *testing.T) {
	// Condition 1: p.age > 18
	cond1 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	)

	// Condition 2: p.age > 21 (différente valeur)
	cond2 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 21},
	)

	// Condition 3: p.salary > 18 (différent champ)
	cond3 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	)

	// Condition 4: p.age >= 18 (différent opérateur)
	cond4 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		">=",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	)

	str1 := CanonicalString(cond1)
	str2 := CanonicalString(cond2)
	str3 := CanonicalString(cond3)
	str4 := CanonicalString(cond4)

	// Toutes les strings doivent être différentes
	if str1 == str2 {
		t.Error("cond1 et cond2 doivent avoir des strings différentes")
	}
	if str1 == str3 {
		t.Error("cond1 et cond3 doivent avoir des strings différentes")
	}
	if str1 == str4 {
		t.Error("cond1 et cond4 doivent avoir des strings différentes")
	}

	// Les hash doivent aussi être différents
	if cond1.Hash == cond2.Hash {
		t.Error("cond1 et cond2 doivent avoir des hash différents")
	}
	if cond1.Hash == cond3.Hash {
		t.Error("cond1 et cond3 doivent avoir des hash différents")
	}
	if cond1.Hash == cond4.Hash {
		t.Error("cond1 et cond4 doivent avoir des hash différents")
	}
}

// TestCanonicalString_Format teste le format de la string canonique
func TestCanonicalString_Format(t *testing.T) {
	cond := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	)

	str := CanonicalString(cond)

	// Vérifier que la string contient les éléments clés
	if !contains(str, "binaryOperation") {
		t.Error("String doit contenir 'binaryOperation'")
	}
	if !contains(str, "fieldAccess") {
		t.Error("String doit contenir 'fieldAccess'")
	}
	if !contains(str, "p") {
		t.Error("String doit contenir 'p'")
	}
	if !contains(str, "age") {
		t.Error("String doit contenir 'age'")
	}
	if !contains(str, ">") {
		t.Error("String doit contenir '>'")
	}
	if !contains(str, "18") {
		t.Error("String doit contenir '18'")
	}
}

// TestCanonicalString_MapFormat teste le format canonique avec des maps
func TestCanonicalString_MapFormat(t *testing.T) {
	cond := NewSimpleCondition(
		"binaryOperation",
		map[string]interface{}{
			"type":   "fieldAccess",
			"object": "p",
			"field":  "salary",
		},
		"+",
		map[string]interface{}{
			"type":  "numberLiteral",
			"value": 100,
		},
	)

	str := CanonicalString(cond)

	// Vérifier que la string contient les éléments clés
	if !contains(str, "fieldAccess") {
		t.Error("String doit contenir 'fieldAccess'")
	}
	if !contains(str, "salary") {
		t.Error("String doit contenir 'salary'")
	}
	if !contains(str, "+") {
		t.Error("String doit contenir '+'")
	}
	if !contains(str, "100") {
		t.Error("String doit contenir '100'")
	}
}

// TestCompareConditions teste la comparaison de conditions
func TestCompareConditions(t *testing.T) {
	cond1 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	)

	cond2 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	)

	cond3 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 21},
	)

	if !CompareConditions(cond1, cond2) {
		t.Error("cond1 et cond2 doivent être égales")
	}

	if CompareConditions(cond1, cond3) {
		t.Error("cond1 et cond3 ne doivent pas être égales")
	}
}

// TestDeduplicateConditions teste la déduplication de conditions
func TestDeduplicateConditions(t *testing.T) {
	cond1 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	)

	cond2 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	)

	cond3 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
		">=",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
	)

	conditions := []SimpleCondition{cond1, cond2, cond3}
	deduplicated := DeduplicateConditions(conditions)

	if len(deduplicated) != 2 {
		t.Errorf("Attendu 2 conditions après déduplication, obtenu %d", len(deduplicated))
	}
}

// TestExtractConditions_Constraint teste l'extraction depuis un Constraint
func TestExtractConditions_Constraint(t *testing.T) {
	expr := constraint.Constraint{
		Type:     "constraint",
		Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		Operator: ">",
		Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	}

	conditions, opType, err := ExtractConditions(expr)

	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}

	if opType != "SINGLE" {
		t.Errorf("Attendu opType='SINGLE', obtenu '%s'", opType)
	}

	if len(conditions) != 1 {
		t.Fatalf("Attendu 1 condition, obtenu %d", len(conditions))
	}
}

// TestExtractConditions_EmptyExpression teste les expressions vides
func TestExtractConditions_EmptyExpression(t *testing.T) {
	// Test avec un fieldAccess seul (pas une condition)
	expr := constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"}

	conditions, opType, err := ExtractConditions(expr)

	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}

	if opType != "NONE" {
		t.Errorf("Attendu opType='NONE', obtenu '%s'", opType)
	}

	if len(conditions) != 0 {
		t.Errorf("Attendu 0 conditions, obtenu %d", len(conditions))
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
