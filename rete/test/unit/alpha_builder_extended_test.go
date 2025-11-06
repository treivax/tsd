package rete

import (
	"testing"
)

// TestAlphaBuilderPartialCoverage teste les parties non couvertes de alpha_builder.go
func TestAlphaBuilderPartialCoverage(t *testing.T) {
	builder := NewAlphaConditionBuilder()

	// Test AndMultiple avec liste vide - devrait retourner True()
	emptyCondition := builder.AndMultiple()
	emptyCondMap, ok := emptyCondition.(map[string]interface{})
	if !ok {
		t.Fatalf("AndMultiple vide devrait retourner un map")
	}
	if emptyCondMap["type"] != "booleanLiteral" || emptyCondMap["value"] != true {
		t.Fatalf("AndMultiple vide devrait retourner True()")
	}

	// Test OrMultiple avec liste vide - devrait retourner False()
	emptyOrCondition := builder.OrMultiple()
	emptyOrCondMap, ok := emptyOrCondition.(map[string]interface{})
	if !ok {
		t.Fatalf("OrMultiple vide devrait retourner un map")
	}
	if emptyOrCondMap["type"] != "booleanLiteral" || emptyOrCondMap["value"] != false {
		t.Fatalf("OrMultiple vide devrait retourner False()")
	}

	// Test AndMultiple avec une seule condition - retourne directement la condition
	singleCondition := builder.FieldEquals("obj", "field", "value")
	andSingle := builder.AndMultiple(singleCondition)

	// Vérifier que la condition est retournée directement
	andSingleMap, ok := andSingle.(map[string]interface{})
	if !ok {
		t.Fatal("AndMultiple devrait retourner un map")
	}
	// Avec une seule condition, AndMultiple retourne la condition elle-même
	if andSingleMap["type"] != "binaryOperation" {
		t.Fatalf("Type incorrect pour AndMultiple avec 1 condition: attendu 'binaryOperation', obtenu %v", andSingleMap["type"])
	}

	// Test OrMultiple avec une seule condition - retourne directement la condition
	orSingle := builder.OrMultiple(singleCondition)
	orSingleMap, ok := orSingle.(map[string]interface{})
	if !ok {
		t.Fatal("OrMultiple devrait retourner un map")
	}
	// Avec une seule condition, OrMultiple retourne la condition elle-même
	if orSingleMap["type"] != "binaryOperation" {
		t.Fatalf("Type incorrect pour OrMultiple avec 1 condition: attendu 'binaryOperation', obtenu %v", orSingleMap["type"])
	}

	// Test AndMultiple avec plusieurs conditions
	condition1 := builder.FieldEquals("obj", "field1", "value1")
	condition2 := builder.FieldEquals("obj", "field2", "value2")
	andMultiple := builder.AndMultiple(condition1, condition2)

	andMultipleMap, ok := andMultiple.(map[string]interface{})
	if !ok {
		t.Fatal("AndMultiple avec plusieurs conditions devrait retourner un map")
	}
	if andMultipleMap["type"] != "logicalExpression" {
		t.Fatalf("Type incorrect pour AndMultiple: attendu 'logicalExpression', obtenu %v", andMultipleMap["type"])
	}

	// Test OrMultiple avec plusieurs conditions
	orMultiple := builder.OrMultiple(condition1, condition2)
	orMultipleMap, ok := orMultiple.(map[string]interface{})
	if !ok {
		t.Fatal("OrMultiple avec plusieurs conditions devrait retourner un map")
	}
	if orMultipleMap["type"] != "logicalExpression" {
		t.Fatalf("Type incorrect pour OrMultiple: attendu 'logicalExpression', obtenu %v", orMultipleMap["type"])
	}
}

// TestFieldInNotInCoverage teste les cas non couverts de FieldIn et FieldNotIn
func TestFieldInNotInCoverage(t *testing.T) {
	builder := NewAlphaConditionBuilder()

	// Test FieldIn avec liste vide - devrait retourner False()
	emptyFieldIn := builder.FieldIn("obj", "status")
	fieldInMap, ok := emptyFieldIn.(map[string]interface{})
	if !ok {
		t.Fatal("FieldIn vide devrait retourner un map")
	}

	// FieldIn avec liste vide retourne False(), donc type "booleanLiteral"
	if fieldInMap["type"] != "booleanLiteral" {
		t.Fatalf("FieldIn vide devrait retourner False(), obtenu type: %v", fieldInMap["type"])
	}

	// Test FieldNotIn avec liste vide - devrait retourner True()
	emptyFieldNotIn := builder.FieldNotIn("obj", "status")
	fieldNotInMap, ok := emptyFieldNotIn.(map[string]interface{})
	if !ok {
		t.Fatal("FieldNotIn vide devrait retourner un map")
	}

	// FieldNotIn avec liste vide retourne True(), donc type "booleanLiteral"
	if fieldNotInMap["type"] != "booleanLiteral" {
		t.Fatalf("FieldNotIn vide devrait retourner True(), obtenu type: %v", fieldNotInMap["type"])
	}

	// Test FieldIn avec plusieurs valeurs de types différents
	mixedFieldIn := builder.FieldIn("obj", "mixed_field", "string", 42, 3.14, true)
	mixedMap, ok := mixedFieldIn.(map[string]interface{})
	if !ok {
		t.Fatal("FieldIn avec valeurs mixtes devrait retourner un map")
	}

	// FieldIn avec plusieurs valeurs retourne une expression logique OR
	if mixedMap["type"] != "logicalExpression" {
		t.Fatalf("FieldIn avec valeurs devrait retourner logicalExpression, obtenu: %v", mixedMap["type"])
	}

	operations, exists := mixedMap["operations"]
	if !exists {
		t.Fatal("FieldIn devrait avoir un champ 'operations'")
	}

	operationsList, ok := operations.([]interface{})
	if !ok || len(operationsList) == 0 {
		t.Fatal("FieldIn devrait avoir des opérations OR")
	}

	firstOp, ok := operationsList[0].(map[string]interface{})
	if !ok || firstOp["op"] != "OR" {
		t.Fatalf("FieldIn devrait utiliser l'opérateur OR, obtenu: %v", firstOp["op"])
	}

	// Test FieldNotIn avec plusieurs valeurs
	mixedFieldNotIn := builder.FieldNotIn("obj", "mixed_field", "string", 42)
	mixedNotInMap, ok := mixedFieldNotIn.(map[string]interface{})
	if !ok {
		t.Fatal("FieldNotIn avec valeurs devrait retourner un map")
	}

	// FieldNotIn avec plusieurs valeurs retourne une expression logique AND
	if mixedNotInMap["type"] != "logicalExpression" {
		t.Fatalf("FieldNotIn avec valeurs devrait retourner logicalExpression, obtenu: %v", mixedNotInMap["type"])
	}

	notInOperations, exists := mixedNotInMap["operations"]
	if !exists {
		t.Fatal("FieldNotIn devrait avoir un champ 'operations'")
	}

	notInOperationsList, ok := notInOperations.([]interface{})
	if !ok || len(notInOperationsList) == 0 {
		t.Fatal("FieldNotIn devrait avoir des opérations AND")
	}

	firstNotInOp, ok := notInOperationsList[0].(map[string]interface{})
	if !ok || firstNotInOp["op"] != "AND" {
		t.Fatalf("FieldNotIn devrait utiliser l'opérateur AND, obtenu: %v", firstNotInOp["op"])
	}
}

// TestCreateLiteralCoverage teste les cas non couverts de createLiteral
func TestCreateLiteralCoverage(t *testing.T) {
	builder := NewAlphaConditionBuilder()

	// Test avec des types de données non couverts auparavant
	testCases := []struct {
		value    interface{}
		expected string
	}{
		// Types supportés nativement
		{int32(42), "numberLiteral"},
		{int64(42), "numberLiteral"},
		{float32(3.14), "numberLiteral"},
		{float64(3.14159), "numberLiteral"},
		{true, "booleanLiteral"},
		{false, "booleanLiteral"},
		{"hello", "stringLiteral"},
		// Types non supportés qui tombent dans le cas default
		{int8(42), "stringLiteral"},
		{int16(42), "stringLiteral"},
		{uint(42), "stringLiteral"},
		{uint8(42), "stringLiteral"},
		{uint16(42), "stringLiteral"},
		{uint32(42), "stringLiteral"},
		{uint64(42), "stringLiteral"},
		{[]byte("hello"), "stringLiteral"}, // Les []byte sont convertis en string
		{nil, "stringLiteral"},             // nil est converti en string vide
	}

	for _, tc := range testCases {
		result := builder.createLiteral(tc.value)
		// result est déjà un map[string]interface{}, pas besoin de type assertion
		resultType, exists := result["type"]
		if !exists {
			t.Fatalf("createLiteral devrait avoir un champ 'type' pour %T", tc.value)
		}

		if resultType != tc.expected {
			t.Fatalf("Type incorrect pour %T: attendu %s, obtenu %v", tc.value, tc.expected, resultType)
		}
	}

	// Test avec un type non supporté (struct personnalisé)
	type CustomStruct struct {
		Field string
	}
	custom := CustomStruct{Field: "test"}

	result := builder.createLiteral(custom)
	// result est déjà un map[string]interface{}, pas besoin de type assertion
	resultType, exists := result["type"]
	if !exists {
		t.Fatal("createLiteral devrait avoir un champ 'type'")
	}

	// Les types non supportés sont convertis en stringLiteral
	if resultType != "stringLiteral" {
		t.Fatalf("Type non supporté devrait être converti en stringLiteral, obtenu %v", resultType)
	}
}

// TestComplexLogicalExpressions teste les expressions logiques complexes
func TestComplexLogicalExpressions(t *testing.T) {
	builder := NewAlphaConditionBuilder()

	// Créer des conditions de base
	cond1 := builder.FieldEquals("priority", "high", "high")
	cond2 := builder.FieldEquals("status", "active", "active")
	cond3 := builder.FieldGreaterThan("score", "threshold", 80.0)

	// Test de compositions complexes
	andCond := builder.And(cond1, cond2)
	orCond := builder.Or(andCond, cond3)

	// Vérifier que les compositions fonctionnent
	orCondMap, ok := orCond.(map[string]interface{})
	if !ok {
		t.Fatal("La composition complexe devrait retourner un map")
	}

	if orCondMap["type"] != "logicalExpression" {
		t.Fatalf("Type incorrect pour la composition: attendu 'logicalExpression', obtenu %v", orCondMap["type"])
	}

	// Test AndMultiple avec composition
	complexAnd := builder.AndMultiple(cond1, cond2, cond3)
	complexAndMap, ok := complexAnd.(map[string]interface{})
	if !ok {
		t.Fatal("AndMultiple complexe devrait retourner un map")
	}

	if complexAndMap["type"] != "logicalExpression" {
		t.Fatalf("AndMultiple complexe devrait être une logicalExpression, obtenu: %v", complexAndMap["type"])
	}

	operations, exists := complexAndMap["operations"]
	if !exists {
		t.Fatal("AndMultiple devrait avoir un champ 'operations'")
	}

	operationsList, ok := operations.([]interface{})
	if !ok {
		t.Fatal("Le champ 'operations' devrait être une liste")
	}

	// AndMultiple(cond1, cond2, cond3) crée: left=cond1, operations=[{op:AND, right:cond2}, {op:AND, right:cond3}]
	if len(operationsList) != 2 {
		t.Fatalf("Attendu 2 opérations, obtenu %d", len(operationsList))
	}
}
