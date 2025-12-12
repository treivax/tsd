// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"github.com/treivax/tsd/constraint"
	"testing"
)

// TestIsCommutative_AllOperators vérifie que IsCommutative identifie correctement
// les opérateurs commutatifs et non-commutatifs
func TestIsCommutative_AllOperators(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		want     bool
	}{
		// Opérateurs commutatifs
		{"AND", "AND", true},
		{"OR", "OR", true},
		{"&&", "&&", true},
		{"||", "||", true},
		{"+", "+", true},
		{"*", "*", true},
		{"==", "==", true},
		{"!=", "!=", true},
		{"<>", "<>", true},
		// Opérateurs non-commutatifs
		{"-", "-", false},
		{"/", "/", false},
		{"<", "<", false},
		{">", ">", false},
		{"<=", "<=", false},
		{">=", ">=", false},
		{"XOR", "XOR", false},
		{"THEN", "THEN", false},
		{"SEQ", "SEQ", false},
		{"unknown", "unknown", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsCommutative(tt.operator)
			if got != tt.want {
				t.Errorf("IsCommutative(%q) = %v, want %v", tt.operator, got, tt.want)
			}
		})
	}
}

// TestNormalizeConditions_AND_OrderIndependent vérifie que A AND B et B AND A
// normalisent au même ordre
func TestNormalizeConditions_AND_OrderIndependent(t *testing.T) {
	// Créer deux conditions: age > 18 et salary >= 50000
	condAge := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	)
	condSalary := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
		">=",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
	)
	// Test avec ordre A, B
	conditionsAB := []SimpleCondition{condAge, condSalary}
	normalizedAB := NormalizeConditions(conditionsAB, "AND")
	// Test avec ordre B, A
	conditionsBA := []SimpleCondition{condSalary, condAge}
	normalizedBA := NormalizeConditions(conditionsBA, "AND")
	// Vérifier que les deux ordres produisent le même résultat
	if len(normalizedAB) != len(normalizedBA) {
		t.Fatalf("Longueurs différentes: %d vs %d", len(normalizedAB), len(normalizedBA))
	}
	for i := range normalizedAB {
		if !CompareConditions(normalizedAB[i], normalizedBA[i]) {
			t.Errorf("Position %d: conditions différentes\nAB: %s\nBA: %s",
				i,
				CanonicalString(normalizedAB[i]),
				CanonicalString(normalizedBA[i]),
			)
		}
	}
}

// TestNormalizeConditions_OR_OrderIndependent vérifie que A OR B et B OR A
// normalisent au même ordre
func TestNormalizeConditions_OR_OrderIndependent(t *testing.T) {
	// Créer deux conditions: status == "active" et verified == true
	condStatus := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
		"==",
		constraint.StringLiteral{Type: "stringLiteral", Value: "active"},
	)
	condVerified := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "verified"},
		"==",
		constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
	)
	// Test avec ordre A, B
	conditionsAB := []SimpleCondition{condStatus, condVerified}
	normalizedAB := NormalizeConditions(conditionsAB, "OR")
	// Test avec ordre B, A
	conditionsBA := []SimpleCondition{condVerified, condStatus}
	normalizedBA := NormalizeConditions(conditionsBA, "OR")
	// Vérifier que les deux ordres produisent le même résultat
	if len(normalizedAB) != len(normalizedBA) {
		t.Fatalf("Longueurs différentes: %d vs %d", len(normalizedAB), len(normalizedBA))
	}
	for i := range normalizedAB {
		if !CompareConditions(normalizedAB[i], normalizedBA[i]) {
			t.Errorf("Position %d: conditions différentes\nAB: %s\nBA: %s",
				i,
				CanonicalString(normalizedAB[i]),
				CanonicalString(normalizedBA[i]),
			)
		}
	}
}

// TestNormalizeConditions_NonCommutative_PreserveOrder vérifie que les opérateurs
// non-commutatifs préservent l'ordre original
func TestNormalizeConditions_NonCommutative_PreserveOrder(t *testing.T) {
	// Créer deux conditions pour des opérations séquentielles
	cond1 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "x", Field: "value"},
		"-",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 10},
	)
	cond2 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "y", Field: "value"},
		"-",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 5},
	)
	// L'ordre original
	original := []SimpleCondition{cond1, cond2}
	// Normaliser avec un opérateur non-commutatif
	normalized := NormalizeConditions(original, "SEQ")
	// Vérifier que l'ordre est préservé
	if len(normalized) != len(original) {
		t.Fatalf("Longueurs différentes: %d vs %d", len(normalized), len(original))
	}
	for i := range original {
		if !CompareConditions(normalized[i], original[i]) {
			t.Errorf("Position %d: ordre non préservé\nOriginal: %s\nNormalized: %s",
				i,
				CanonicalString(original[i]),
				CanonicalString(normalized[i]),
			)
		}
	}
}

// TestNormalizeConditions_EmptyAndSingle vérifie le comportement avec des cas limites
func TestNormalizeConditions_EmptyAndSingle(t *testing.T) {
	// Test avec liste vide
	empty := []SimpleCondition{}
	normalizedEmpty := NormalizeConditions(empty, "AND")
	if len(normalizedEmpty) != 0 {
		t.Errorf("Expected empty list, got %d conditions", len(normalizedEmpty))
	}
	// Test avec une seule condition
	single := []SimpleCondition{
		NewSimpleCondition(
			"binaryOperation",
			constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
			">",
			constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
		),
	}
	normalizedSingle := NormalizeConditions(single, "AND")
	if len(normalizedSingle) != 1 {
		t.Errorf("Expected 1 condition, got %d", len(normalizedSingle))
	}
	if !CompareConditions(normalizedSingle[0], single[0]) {
		t.Error("Single condition was modified")
	}
}

// TestNormalizeConditions_ThreeConditions vérifie la normalisation avec 3+ conditions
func TestNormalizeConditions_ThreeConditions(t *testing.T) {
	// Créer trois conditions: C, A, B
	condA := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "aaa"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 1},
	)
	condB := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "bbb"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 2},
	)
	condC := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "ccc"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 3},
	)
	// Test plusieurs permutations
	testCases := [][]SimpleCondition{
		{condC, condA, condB},
		{condB, condC, condA},
		{condA, condC, condB},
		{condC, condB, condA},
	}
	var firstNormalized []SimpleCondition
	for i, conditions := range testCases {
		normalized := NormalizeConditions(conditions, "AND")
		if i == 0 {
			firstNormalized = normalized
		} else {
			// Vérifier que toutes les permutations produisent le même ordre
			if len(normalized) != len(firstNormalized) {
				t.Errorf("Test case %d: longueur différente: %d vs %d", i, len(normalized), len(firstNormalized))
				continue
			}
			for j := range normalized {
				if !CompareConditions(normalized[j], firstNormalized[j]) {
					t.Errorf("Test case %d, position %d: conditions différentes\nGot: %s\nWant: %s",
						i, j,
						CanonicalString(normalized[j]),
						CanonicalString(firstNormalized[j]),
					)
				}
			}
		}
	}
}

// TestNormalizeExpression_ComplexNested vérifie la normalisation d'expressions complexes imbriquées
func TestNormalizeExpression_ComplexNested(t *testing.T) {
	// Expression: (age > 18 AND salary >= 50000) OR (status == "active")
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
	// Normaliser l'expression
	normalized, err := NormalizeExpression(expr)
	if err != nil {
		t.Fatalf("NormalizeExpression failed: %v", err)
	}
	// Vérifier que le résultat est du bon type
	_, ok := normalized.(constraint.LogicalExpression)
	if !ok {
		t.Errorf("Expected LogicalExpression, got %T", normalized)
	}
}

// TestNormalizeExpression_BinaryOperation vérifie la normalisation d'opérations binaires simples
func TestNormalizeExpression_BinaryOperation(t *testing.T) {
	expr := constraint.BinaryOperation{
		Type:     "binaryOperation",
		Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		Operator: ">",
		Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	}
	normalized, err := NormalizeExpression(expr)
	if err != nil {
		t.Fatalf("NormalizeExpression failed: %v", err)
	}
	// Vérifier que le résultat est inchangé pour une opération binaire simple
	binOp, ok := normalized.(constraint.BinaryOperation)
	if !ok {
		t.Fatalf("Expected BinaryOperation, got %T", normalized)
	}
	if binOp.Operator != ">" {
		t.Errorf("Expected operator '>', got '%s'", binOp.Operator)
	}
}

// TestNormalizeExpression_Map vérifie la normalisation d'expressions sous forme de map
func TestNormalizeExpression_Map(t *testing.T) {
	expr := map[string]interface{}{
		"type": "binaryOperation",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "p",
			"field":  "age",
		},
		"operator": ">",
		"right": map[string]interface{}{
			"type":  "numberLiteral",
			"value": 18,
		},
	}
	normalized, err := NormalizeExpression(expr)
	if err != nil {
		t.Fatalf("NormalizeExpression failed: %v", err)
	}
	// Vérifier que le résultat est une map
	normMap, ok := normalized.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map, got %T", normalized)
	}
	if normMap["type"] != "binaryOperation" {
		t.Errorf("Expected type 'binaryOperation', got '%v'", normMap["type"])
	}
}

// TestNormalizeExpression_Literals vérifie que les littéraux ne sont pas modifiés
func TestNormalizeExpression_Literals(t *testing.T) {
	tests := []struct {
		name string
		expr interface{}
	}{
		{
			name: "NumberLiteral",
			expr: constraint.NumberLiteral{Type: "numberLiteral", Value: 42},
		},
		{
			name: "StringLiteral",
			expr: constraint.StringLiteral{Type: "stringLiteral", Value: "test"},
		},
		{
			name: "BooleanLiteral",
			expr: constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
		},
		{
			name: "FieldAccess",
			expr: constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalized, err := NormalizeExpression(tt.expr)
			if err != nil {
				t.Fatalf("NormalizeExpression failed: %v", err)
			}
			// Les littéraux et field access devraient être retournés inchangés
			if normalized != tt.expr {
				t.Errorf("Expression was modified, expected unchanged")
			}
		})
	}
}

// TestNormalizeConditions_DeterministicOrder vérifie que l'ordre est vraiment déterministe
func TestNormalizeConditions_DeterministicOrder(t *testing.T) {
	// Créer plusieurs conditions
	conditions := []SimpleCondition{
		NewSimpleCondition(
			"binaryOperation",
			constraint.FieldAccess{Type: "fieldAccess", Object: "z", Field: "zzz"},
			">",
			constraint.NumberLiteral{Type: "numberLiteral", Value: 1},
		),
		NewSimpleCondition(
			"binaryOperation",
			constraint.FieldAccess{Type: "fieldAccess", Object: "a", Field: "aaa"},
			">",
			constraint.NumberLiteral{Type: "numberLiteral", Value: 2},
		),
		NewSimpleCondition(
			"binaryOperation",
			constraint.FieldAccess{Type: "fieldAccess", Object: "m", Field: "mmm"},
			">",
			constraint.NumberLiteral{Type: "numberLiteral", Value: 3},
		),
	}
	// Normaliser plusieurs fois
	first := NormalizeConditions(conditions, "AND")
	second := NormalizeConditions(conditions, "AND")
	third := NormalizeConditions(conditions, "AND")
	// Vérifier que tous les résultats sont identiques
	for i := range first {
		if !CompareConditions(first[i], second[i]) || !CompareConditions(first[i], third[i]) {
			t.Errorf("Position %d: ordre non déterministe", i)
		}
	}
}

// TestRebuildLogicalExpression_SingleCondition teste la reconstruction avec une seule condition
func TestRebuildLogicalExpression_SingleCondition(t *testing.T) {
	// Créer une condition simple
	cond := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	)
	// Reconstruire l'expression
	rebuilt, err := rebuildLogicalExpression([]SimpleCondition{cond}, "AND")
	if err != nil {
		t.Fatalf("rebuildLogicalExpression failed: %v", err)
	}
	// Vérifier la structure
	if rebuilt.Type != "logicalExpr" {
		t.Errorf("Expected type 'logicalExpr', got '%s'", rebuilt.Type)
	}
	if len(rebuilt.Operations) != 0 {
		t.Errorf("Expected 0 operations for single condition, got %d", len(rebuilt.Operations))
	}
	// Vérifier que Left est une BinaryOperation
	binOp, ok := rebuilt.Left.(constraint.BinaryOperation)
	if !ok {
		t.Fatalf("Expected Left to be BinaryOperation, got %T", rebuilt.Left)
	}
	if binOp.Operator != ">" {
		t.Errorf("Expected operator '>', got '%s'", binOp.Operator)
	}
}

// TestRebuildLogicalExpression_TwoConditions teste la reconstruction avec deux conditions
func TestRebuildLogicalExpression_TwoConditions(t *testing.T) {
	// Créer deux conditions
	condAge := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	)
	condSalary := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
		">=",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
	)
	// Reconstruire l'expression
	rebuilt, err := rebuildLogicalExpression([]SimpleCondition{condAge, condSalary}, "AND")
	if err != nil {
		t.Fatalf("rebuildLogicalExpression failed: %v", err)
	}
	// Vérifier la structure
	if rebuilt.Type != "logicalExpr" {
		t.Errorf("Expected type 'logicalExpr', got '%s'", rebuilt.Type)
	}
	if len(rebuilt.Operations) != 1 {
		t.Fatalf("Expected 1 operation, got %d", len(rebuilt.Operations))
	}
	// Vérifier l'opération
	if rebuilt.Operations[0].Op != "AND" {
		t.Errorf("Expected operator 'AND', got '%s'", rebuilt.Operations[0].Op)
	}
	// Vérifier Left
	leftOp, ok := rebuilt.Left.(constraint.BinaryOperation)
	if !ok {
		t.Fatalf("Expected Left to be BinaryOperation, got %T", rebuilt.Left)
	}
	if leftOp.Operator != ">" {
		t.Errorf("Expected Left operator '>', got '%s'", leftOp.Operator)
	}
	// Vérifier Right
	rightOp, ok := rebuilt.Operations[0].Right.(constraint.BinaryOperation)
	if !ok {
		t.Fatalf("Expected Right to be BinaryOperation, got %T", rebuilt.Operations[0].Right)
	}
	if rightOp.Operator != ">=" {
		t.Errorf("Expected Right operator '>=', got '%s'", rightOp.Operator)
	}
}

// TestRebuildLogicalExpression_ThreeConditions teste la reconstruction avec trois conditions
func TestRebuildLogicalExpression_ThreeConditions(t *testing.T) {
	// Créer trois conditions
	cond1 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "aaa"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 1},
	)
	cond2 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "bbb"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 2},
	)
	cond3 := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "ccc"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 3},
	)
	// Reconstruire l'expression
	rebuilt, err := rebuildLogicalExpression([]SimpleCondition{cond1, cond2, cond3}, "OR")
	if err != nil {
		t.Fatalf("rebuildLogicalExpression failed: %v", err)
	}
	// Vérifier la structure
	if len(rebuilt.Operations) != 2 {
		t.Fatalf("Expected 2 operations, got %d", len(rebuilt.Operations))
	}
	// Vérifier que tous les opérateurs sont OR
	for i, op := range rebuilt.Operations {
		if op.Op != "OR" {
			t.Errorf("Operation %d: expected 'OR', got '%s'", i, op.Op)
		}
	}
}

// TestRebuildLogicalExpression_Empty teste le cas d'erreur avec liste vide
func TestRebuildLogicalExpression_Empty(t *testing.T) {
	_, err := rebuildLogicalExpression([]SimpleCondition{}, "AND")
	if err == nil {
		t.Error("Expected error for empty conditions, got nil")
	}
}

// TestNormalizeExpression_WithReconstruction teste la normalisation complète avec reconstruction
func TestNormalizeExpression_WithReconstruction(t *testing.T) {
	// Expression : salary >= 50000 AND age > 18 (ordre inversé)
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
	// Normaliser l'expression
	normalized, err := NormalizeExpression(expr)
	if err != nil {
		t.Fatalf("NormalizeExpression failed: %v", err)
	}
	// Vérifier que c'est une LogicalExpression
	normExpr, ok := normalized.(constraint.LogicalExpression)
	if !ok {
		t.Fatalf("Expected LogicalExpression, got %T", normalized)
	}
	// Vérifier que Left est age > 18 (ordre canonique)
	leftOp, ok := normExpr.Left.(constraint.BinaryOperation)
	if !ok {
		t.Fatalf("Expected Left to be BinaryOperation, got %T", normExpr.Left)
	}
	leftField, ok := leftOp.Left.(constraint.FieldAccess)
	if !ok {
		t.Fatalf("Expected Left.Left to be FieldAccess, got %T", leftOp.Left)
	}
	// age vient avant salary dans l'ordre canonique
	if leftField.Field != "age" {
		t.Errorf("Expected first condition to be 'age', got '%s'", leftField.Field)
	}
	// Vérifier que Right est salary >= 50000
	if len(normExpr.Operations) != 1 {
		t.Fatalf("Expected 1 operation, got %d", len(normExpr.Operations))
	}
	rightOp, ok := normExpr.Operations[0].Right.(constraint.BinaryOperation)
	if !ok {
		t.Fatalf("Expected Right to be BinaryOperation, got %T", normExpr.Operations[0].Right)
	}
	rightField, ok := rightOp.Left.(constraint.FieldAccess)
	if !ok {
		t.Fatalf("Expected Right.Left to be FieldAccess, got %T", rightOp.Left)
	}
	if rightField.Field != "salary" {
		t.Errorf("Expected second condition to be 'salary', got '%s'", rightField.Field)
	}
}

// TestNormalizeExpression_PreservesSemantics teste que la normalisation préserve la sémantique
func TestNormalizeExpression_PreservesSemantics(t *testing.T) {
	// Créer deux expressions équivalentes dans des ordres différents
	expr1 := constraint.LogicalExpression{
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
	expr2 := constraint.LogicalExpression{
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
	// Normaliser les deux
	norm1, err1 := NormalizeExpression(expr1)
	norm2, err2 := NormalizeExpression(expr2)
	if err1 != nil {
		t.Fatalf("NormalizeExpression(expr1) failed: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("NormalizeExpression(expr2) failed: %v", err2)
	}
	// Extraire les conditions des deux expressions normalisées
	conds1, _, err := ExtractConditions(norm1)
	if err != nil {
		t.Fatalf("ExtractConditions(norm1) failed: %v", err)
	}
	conds2, _, err := ExtractConditions(norm2)
	if err != nil {
		t.Fatalf("ExtractConditions(norm2) failed: %v", err)
	}
	// Vérifier que les deux ont le même nombre de conditions
	if len(conds1) != len(conds2) {
		t.Fatalf("Different number of conditions: %d vs %d", len(conds1), len(conds2))
	}
	// Vérifier que les conditions sont dans le même ordre
	for i := range conds1 {
		if !CompareConditions(conds1[i], conds2[i]) {
			t.Errorf("Position %d: conditions differ\n  conds1: %s\n  conds2: %s",
				i,
				CanonicalString(conds1[i]),
				CanonicalString(conds2[i]),
			)
		}
	}
}

// TestRebuildLogicalExpressionMap_TwoConditions teste la reconstruction de map
func TestRebuildLogicalExpressionMap_TwoConditions(t *testing.T) {
	// Créer deux conditions
	condAge := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		">",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	)
	condSalary := NewSimpleCondition(
		"binaryOperation",
		constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
		">=",
		constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
	)
	// Reconstruire en map
	rebuilt, err := rebuildLogicalExpressionMap([]SimpleCondition{condAge, condSalary}, "AND")
	if err != nil {
		t.Fatalf("rebuildLogicalExpressionMap failed: %v", err)
	}
	// Vérifier le type
	exprType, ok := rebuilt["type"].(string)
	if !ok || exprType != "logicalExpr" {
		t.Errorf("Expected type 'logicalExpr', got '%v'", rebuilt["type"])
	}
	// Vérifier les opérations
	operations, ok := rebuilt["operations"].([]interface{})
	if !ok {
		t.Fatalf("Expected operations to be []interface{}, got %T", rebuilt["operations"])
	}
	if len(operations) != 1 {
		t.Errorf("Expected 1 operation, got %d", len(operations))
	}
}

// TestNormalizeExpressionMap_WithReconstruction teste la normalisation de map avec reconstruction
func TestNormalizeExpressionMap_WithReconstruction(t *testing.T) {
	// Expression map : salary >= 50000 AND age > 18
	expr := map[string]interface{}{
		"type": "logicalExpression",
		"left": map[string]interface{}{
			"type":     "binaryOperation",
			"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "salary"},
			"operator": ">=",
			"right":    map[string]interface{}{"type": "numberLiteral", "value": 50000},
		},
		"operations": []interface{}{
			map[string]interface{}{
				"op": "AND",
				"right": map[string]interface{}{
					"type":     "binaryOperation",
					"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
					"operator": ">",
					"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
				},
			},
		},
	}
	// Normaliser
	normalized, err := NormalizeExpression(expr)
	if err != nil {
		t.Fatalf("NormalizeExpression failed: %v", err)
	}
	// Vérifier que c'est une map
	normMap, ok := normalized.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map, got %T", normalized)
	}
	// Vérifier le type
	if normMap["type"] != "logicalExpr" {
		t.Errorf("Expected type 'logicalExpr', got '%v'", normMap["type"])
	}
	// Vérifier que Left existe
	if normMap["left"] == nil {
		t.Error("Left is nil")
	}
	// Vérifier les opérations
	operations, ok := normMap["operations"].([]interface{})
	if !ok {
		t.Fatalf("Expected operations to be []interface{}, got %T", normMap["operations"])
	}
	if len(operations) != 1 {
		t.Errorf("Expected 1 operation, got %d", len(operations))
	}
}
