package rete

import (
	"testing"

	"github.com/treivax/tsd/constraint"
)

// TestEvaluatorUncoveredFunctions teste les fonctions non couvertes de l'évaluateur
func TestEvaluatorUncoveredFunctions(t *testing.T) {
	evaluator := NewAlphaConditionEvaluator()

	// Préparer des faits de test
	person := &Fact{
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Alice",
			"age":  25,
			"city": "Paris",
		},
	}

	// Lier les variables
	evaluator.variableBindings = map[string]*Fact{
		"p": person,
	}

	t.Run("evaluateBinaryOperation", func(t *testing.T) {
		// Créer une opération binaire simple
		binOp := constraint.BinaryOperation{
			Type:     "binary_operation",
			Left:     "Alice", // literal direct
			Operator: "==",
			Right:    "Alice", // literal direct
		}

		result, err := evaluator.evaluateBinaryOperation(binOp)
		if err != nil {
			t.Fatalf("evaluateBinaryOperation ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour Alice == Alice")
		}

		// Test avec des nombres
		binOpNum := constraint.BinaryOperation{
			Type:     "binary_operation",
			Left:     25,
			Operator: ">",
			Right:    20,
		}

		result, err = evaluator.evaluateBinaryOperation(binOpNum)
		if err != nil {
			t.Fatalf("evaluateBinaryOperation (nombres) ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour 25 > 20")
		}
	})

	t.Run("evaluateLogicalExpression", func(t *testing.T) {
		// Créer une expression logique AND avec des constraints valides
		logExpr := constraint.LogicalExpression{
			Type: "logical_expression",
			Left: constraint.BinaryOperation{
				Type:     "binary_operation",
				Left:     10,
				Operator: ">",
				Right:    5,
			},
			Operations: []constraint.LogicalOperation{
				{
					Op: "AND",
					Right: constraint.BinaryOperation{
						Type:     "binary_operation",
						Left:     20,
						Operator: "<",
						Right:    30,
					},
				},
			},
		}

		result, err := evaluator.evaluateLogicalExpression(logExpr)
		if err != nil {
			t.Fatalf("evaluateLogicalExpression ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour (10 > 5) AND (20 < 30)")
		}

		// Test avec OR
		logExprOr := constraint.LogicalExpression{
			Type: "logical_expression",
			Left: constraint.BinaryOperation{
				Type:     "binary_operation",
				Left:     5,
				Operator: ">",
				Right:    10, // false
			},
			Operations: []constraint.LogicalOperation{
				{
					Op: "OR",
					Right: constraint.BinaryOperation{
						Type:     "binary_operation",
						Left:     20,
						Operator: ">",
						Right:    10, // true
					},
				},
			},
		}

		result, err = evaluator.evaluateLogicalExpression(logExprOr)
		if err != nil {
			t.Fatalf("evaluateLogicalExpression (OR) ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour (5 > 10) OR (20 > 10)")
		}

		// Test opérateur non supporté
		logExprInvalid := constraint.LogicalExpression{
			Type: "logical_expression",
			Left: constraint.BinaryOperation{
				Type:     "binary_operation",
				Left:     1,
				Operator: "==",
				Right:    1,
			},
			Operations: []constraint.LogicalOperation{
				{
					Op: "XOR", // Opérateur non supporté
					Right: constraint.BinaryOperation{
						Type:     "binary_operation",
						Left:     1,
						Operator: "==",
						Right:    1,
					},
				},
			},
		}

		_, err = evaluator.evaluateLogicalExpression(logExprInvalid)
		if err == nil {
			t.Error("evaluateLogicalExpression devrait échouer avec opérateur XOR non supporté")
		}
	})

	t.Run("evaluateConstraint", func(t *testing.T) {
		// Test contrainte simple
		constr := constraint.Constraint{
			Type:     "constraint",
			Left:     "test",
			Operator: "==",
			Right:    "test",
		}

		result, err := evaluator.evaluateConstraint(constr)
		if err != nil {
			t.Fatalf("evaluateConstraint ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour test == test")
		}

		// Test contrainte avec nombres
		constrNum := constraint.Constraint{
			Type:     "constraint",
			Left:     100,
			Operator: ">=",
			Right:    50,
		}

		result, err = evaluator.evaluateConstraint(constrNum)
		if err != nil {
			t.Fatalf("evaluateConstraint (nombres) ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour 100 >= 50")
		}
	})

	t.Run("evaluateConstraintMap", func(t *testing.T) {
		// Test contrainte depuis map
		constraintMap := map[string]interface{}{
			"operator": "!=",
			"left":     "hello",
			"right":    "world",
		}

		result, err := evaluator.evaluateConstraintMap(constraintMap)
		if err != nil {
			t.Fatalf("evaluateConstraintMap ne devrait pas échouer: %v", err)
		}

		if !result {
			t.Error("Attendu true pour hello != world")
		}

		// Test map sans opérateur
		invalidMap := map[string]interface{}{
			"left":  "test",
			"right": "test",
		}

		_, err = evaluator.evaluateConstraintMap(invalidMap)
		if err == nil {
			t.Error("evaluateConstraintMap devrait échouer sans opérateur")
		}
	})

	t.Run("evaluateFieldAccess", func(t *testing.T) {
		// Test accès à un champ valide
		fieldAccess := constraint.FieldAccess{
			Type:   "field_access",
			Object: "p",
			Field:  "name",
		}

		result, err := evaluator.evaluateFieldAccess(fieldAccess)
		if err != nil {
			t.Fatalf("evaluateFieldAccess ne devrait pas échouer: %v", err)
		}

		if result != "Alice" {
			t.Errorf("Attendu 'Alice', obtenu '%v'", result)
		}

		// Test accès à un champ numérique
		fieldAccessAge := constraint.FieldAccess{
			Type:   "field_access",
			Object: "p",
			Field:  "age",
		}

		result, err = evaluator.evaluateFieldAccess(fieldAccessAge)
		if err != nil {
			t.Fatalf("evaluateFieldAccess (age) ne devrait pas échouer: %v", err)
		}

		if result != 25 {
			t.Errorf("Attendu 25, obtenu %v", result)
		}

		// Test accès à un objet inexistant
		fieldAccessInvalid := constraint.FieldAccess{
			Type:   "field_access",
			Object: "inexistant",
			Field:  "field",
		}

		_, err = evaluator.evaluateFieldAccess(fieldAccessInvalid)
		if err == nil {
			t.Error("evaluateFieldAccess devrait échouer avec objet inexistant")
		}

		// Test accès à un champ inexistant
		fieldAccessInvalidField := constraint.FieldAccess{
			Type:   "field_access",
			Object: "p",
			Field:  "inexistant",
		}

		_, err = evaluator.evaluateFieldAccess(fieldAccessInvalidField)
		if err == nil {
			t.Error("evaluateFieldAccess devrait échouer avec champ inexistant")
		}
	})

	t.Run("evaluateVariable", func(t *testing.T) {
		// Test évaluation variable valide
		variable := constraint.Variable{
			Type: "variable",
			Name: "p",
		}

		result, err := evaluator.evaluateVariable(variable)
		if err != nil {
			t.Fatalf("evaluateVariable ne devrait pas échouer: %v", err)
		}

		if result == nil {
			t.Error("Le résultat ne devrait pas être nil")
		}

		// Vérifier que c'est le bon fait
		if fact, ok := result.(*Fact); ok {
			if fact.Type != "Person" {
				t.Errorf("Attendu type Person, obtenu %s", fact.Type)
			}
		} else {
			t.Error("Le résultat devrait être un *Fact")
		}

		// Test variable inexistante
		variableInvalid := constraint.Variable{
			Type: "variable",
			Name: "inexistant",
		}

		_, err = evaluator.evaluateVariable(variableInvalid)
		if err == nil {
			t.Error("evaluateVariable devrait échouer avec variable inexistante")
		}
	})

	t.Run("evaluateVariableByName", func(t *testing.T) {
		// Test évaluation par nom valide
		result, err := evaluator.evaluateVariableByName("p")
		if err != nil {
			t.Fatalf("evaluateVariableByName ne devrait pas échouer: %v", err)
		}

		if result == nil {
			t.Error("Le résultat ne devrait pas être nil")
		}

		// Vérifier que c'est le bon fait
		if fact, ok := result.(*Fact); ok {
			if fact.Type != "Person" {
				t.Errorf("Attendu type Person, obtenu %s", fact.Type)
			}
		} else {
			t.Error("Le résultat devrait être un *Fact")
		}

		// Test nom inexistant
		_, err = evaluator.evaluateVariableByName("variable_inexistante")
		if err == nil {
			t.Error("evaluateVariableByName devrait échouer avec nom inexistant")
		}
	})
}
