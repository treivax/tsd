package rete

import (
	"testing"

	"github.com/treivax/tsd/constraint"
)

// TestEvaluatorPartialCoverage améliore la couverture des fonctions partiellement testées
func TestEvaluatorPartialCoverage(t *testing.T) {
	evaluator := NewAlphaConditionEvaluator()

	// Préparer des faits de test
	person := &Fact{
		Type: "Person",
		Fields: map[string]interface{}{
			"name":   "Bob",
			"age":    30,
			"active": true,
		},
	}

	// Lier les variables
	evaluator.variableBindings = map[string]*Fact{
		"person": person,
	}

	t.Run("evaluateValue avec tous les types", func(t *testing.T) {
		// Test stringLiteral
		literal := map[string]interface{}{
			"type":  "stringLiteral",
			"value": "test_string",
		}

		result, err := evaluator.evaluateValue(literal)
		if err != nil {
			t.Fatalf("evaluateValue (stringLiteral) ne devrait pas échouer: %v", err)
		}

		if result != "test_string" {
			t.Errorf("Attendu 'test_string', obtenu '%v'", result)
		}

		// Test fieldAccess
		fieldAccess := map[string]interface{}{
			"type":   "fieldAccess",
			"object": "person",
			"field":  "name",
		}

		result, err = evaluator.evaluateValue(fieldAccess)
		if err != nil {
			t.Fatalf("evaluateValue (fieldAccess) ne devrait pas échouer: %v", err)
		}

		if result != "Bob" {
			t.Errorf("Attendu 'Bob', obtenu '%v'", result)
		}

		// Test variable
		variable := map[string]interface{}{
			"type": "variable",
			"name": "person",
		}

		result, err = evaluator.evaluateValue(variable)
		if err != nil {
			t.Fatalf("evaluateValue (variable) ne devrait pas échouer: %v", err)
		}

		if result == nil {
			t.Error("Le résultat ne devrait pas être nil")
		}

		// Test numberLiteral
		numberLiteral := map[string]interface{}{
			"type":  "numberLiteral",
			"value": float64(42.5),
		}

		result, err = evaluator.evaluateValue(numberLiteral)
		if err != nil {
			t.Fatalf("evaluateValue (numberLiteral) ne devrait pas échouer: %v", err)
		}

		if result != 42.5 {
			t.Errorf("Attendu 42.5, obtenu %v", result)
		}

		// Test booleanLiteral
		boolLiteral := map[string]interface{}{
			"type":  "booleanLiteral",
			"value": true,
		}

		result, err = evaluator.evaluateValue(boolLiteral)
		if err != nil {
			t.Fatalf("evaluateValue (booleanLiteral) ne devrait pas échouer: %v", err)
		}

		if result != true {
			t.Errorf("Attendu true, obtenu %v", result)
		}

		// Test valeurs directes
		result, err = evaluator.evaluateValue("direct_string")
		if err != nil {
			t.Fatalf("evaluateValue (string direct) ne devrait pas échouer: %v", err)
		}

		if result != "direct_string" {
			t.Errorf("Attendu 'direct_string', obtenu %v", result)
		}

		result, err = evaluator.evaluateValue(123)
		if err != nil {
			t.Fatalf("evaluateValue (int direct) ne devrait pas échouer: %v", err)
		}

		if result != 123 {
			t.Errorf("Attendu 123, obtenu %v", result)
		}

		result, err = evaluator.evaluateValue(3.14)
		if err != nil {
			t.Fatalf("evaluateValue (float direct) ne devrait pas échouer: %v", err)
		}

		if result != 3.14 {
			t.Errorf("Attendu 3.14, obtenu %v", result)
		}

		result, err = evaluator.evaluateValue(false)
		if err != nil {
			t.Fatalf("evaluateValue (bool direct) ne devrait pas échouer: %v", err)
		}

		if result != false {
			t.Errorf("Attendu false, obtenu %v", result)
		}

		// Test type non supporté dans map
		unsupported := map[string]interface{}{
			"type": "unknown_type",
		}

		_, err = evaluator.evaluateValue(unsupported)
		if err == nil {
			t.Error("evaluateValue devrait échouer avec type non supporté")
		}

		// Test type complètement non supporté
		_, err = evaluator.evaluateValue([]int{1, 2, 3}) // slice non supporté
		if err == nil {
			t.Error("evaluateValue devrait échouer avec slice")
		}
	})

	t.Run("evaluateExpression avec différents types", func(t *testing.T) {
		// Test avec map[string]interface{} - binaryOperation
		exprMap := map[string]interface{}{
			"type":     "binaryOperation",
			"left":     25,
			"operator": ">=",
			"right":    18,
		}

		result, err := evaluator.evaluateExpression(exprMap)
		if err != nil {
			t.Fatalf("evaluateExpression (map) ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour 25 >= 18")
		}

		// Test avec map constraint
		constraintMap := map[string]interface{}{
			"type":     "constraint",
			"left":     "hello",
			"operator": "==",
			"right":    "hello",
		}

		result, err = evaluator.evaluateExpression(constraintMap)
		if err != nil {
			t.Fatalf("evaluateExpression (constraint map) ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour hello == hello")
		}

		// Test avec map booleanLiteral
		boolMap := map[string]interface{}{
			"type":  "booleanLiteral",
			"value": true,
		}

		result, err = evaluator.evaluateExpression(boolMap)
		if err != nil {
			t.Fatalf("evaluateExpression (booleanLiteral map) ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour booleanLiteral")
		}

		// Test avec constraint.BinaryOperation
		binaryOp := constraint.BinaryOperation{
			Type:     "binary_operation",
			Left:     "test",
			Operator: "!=",
			Right:    "other",
		}

		result, err = evaluator.evaluateExpression(binaryOp)
		if err != nil {
			t.Fatalf("evaluateExpression (BinaryOperation) ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour test != other")
		}

		// Test avec constraint.LogicalExpression
		logicalExpr := constraint.LogicalExpression{
			Type: "logical_expression",
			Left: constraint.BinaryOperation{
				Type:     "binary_operation",
				Left:     5,
				Operator: "<",
				Right:    10,
			},
			Operations: []constraint.LogicalOperation{
				{
					Op: "OR",
					Right: constraint.BinaryOperation{
						Type:     "binary_operation",
						Left:     15,
						Operator: ">",
						Right:    20,
					},
				},
			},
		}

		result, err = evaluator.evaluateExpression(logicalExpr)
		if err != nil {
			t.Fatalf("evaluateExpression (LogicalExpression) ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour (5 < 10) OR (15 > 20)")
		}

		// Test avec constraint.BooleanLiteral direct
		boolLiteralStruct := constraint.BooleanLiteral{
			Type:  "booleanLiteral",
			Value: true,
		}

		result, err = evaluator.evaluateExpression(boolLiteralStruct)
		if err != nil {
			t.Fatalf("evaluateExpression (BooleanLiteral struct) ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour BooleanLiteral struct")
		}

		// Test avec type non supporté
		_, err = evaluator.evaluateExpression(123) // int non supporté directement
		if err == nil {
			t.Error("evaluateExpression devrait échouer avec type non supporté")
		}

		// Test avec map de type non supporté
		unsupportedMap := map[string]interface{}{
			"type": "unknownType",
		}

		_, err = evaluator.evaluateExpression(unsupportedMap)
		if err == nil {
			t.Error("evaluateExpression devrait échouer avec type de map non supporté")
		}
	})

	t.Run("normalizeValue avec différents types numériques", func(t *testing.T) {
		// Test avec int
		result := evaluator.normalizeValue(42)
		if result != float64(42) {
			t.Errorf("Attendu 42.0, obtenu %v", result)
		}

		// Test avec int32
		result = evaluator.normalizeValue(int32(25))
		if result != float64(25) {
			t.Errorf("Attendu 25.0, obtenu %v", result)
		}

		// Test avec int64
		result = evaluator.normalizeValue(int64(100))
		if result != float64(100) {
			t.Errorf("Attendu 100.0, obtenu %v", result)
		}

		// Test avec float32
		result = evaluator.normalizeValue(float32(3.14))
		expected := float64(float32(3.14)) // Conversion précise
		if result != expected {
			t.Errorf("Attendu %v, obtenu %v", expected, result)
		}

		// Test avec float64
		result = evaluator.normalizeValue(float64(2.718))
		if result != 2.718 {
			t.Errorf("Attendu 2.718, obtenu %v", result)
		}

		// Test avec string (non numérique)
		result = evaluator.normalizeValue("hello")
		if result != "hello" {
			t.Errorf("Attendu 'hello', obtenu %v", result)
		}
	})

	t.Run("compareValues avec cas manqués", func(t *testing.T) {
		// Test opérateur <= avec valeurs égales
		result, err := evaluator.compareValues(10, "<=", 10)
		if err != nil {
			t.Fatalf("compareValues (<=) ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour 10 <= 10")
		}

		// Test opérateur <= avec valeur inférieure
		result, err = evaluator.compareValues(5, "<=", 10)
		if err != nil {
			t.Fatalf("compareValues (<=) ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour 5 <= 10")
		}

		// Test opérateur >= avec valeurs égales
		result, err = evaluator.compareValues(20, ">=", 20)
		if err != nil {
			t.Fatalf("compareValues (>=) ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour 20 >= 20")
		}

		// Test opérateur >= avec valeur supérieure
		result, err = evaluator.compareValues(25, ">=", 20)
		if err != nil {
			t.Fatalf("compareValues (>=) ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour 25 >= 20")
		}

		// Test opérateur non supporté
		_, err = evaluator.compareValues(10, "~=", 10)
		if err == nil {
			t.Error("compareValues devrait échouer avec opérateur non supporté")
		}
	})
}
