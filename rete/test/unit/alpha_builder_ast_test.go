package rete

import (
	"testing"

	"github.com/treivax/tsd/constraint"
)

// TestCreateConstraintFromAST teste la fonction CreateConstraintFromAST
func TestCreateConstraintFromAST(t *testing.T) {
	builder := NewAlphaConditionBuilder()

	t.Run("avec map[string]interface{}", func(t *testing.T) {
		// Test avec une map - devrait la retourner telle quelle
		constraintMap := map[string]interface{}{
			"operator": "==",
			"left": map[string]interface{}{
				"type":   "field_access",
				"object": "person",
				"field":  "age",
			},
			"right": map[string]interface{}{
				"type":  "literal",
				"value": 25,
			},
		}

		result := builder.CreateConstraintFromAST(constraintMap)

		// Le résultat devrait être la même map
		if resultMap, ok := result.(map[string]interface{}); ok {
			if resultMap["operator"] != "==" {
				t.Errorf("Attendu opérateur '==', obtenu '%v'", resultMap["operator"])
			}
		} else {
			t.Error("Le résultat devrait être une map[string]interface{}")
		}

		// Les maps ne peuvent pas être comparées directement en Go
		// On vérifie juste que le type est correct
	})

	t.Run("avec struct quelconque", func(t *testing.T) {
		// Test avec une structure quelconque - devrait la retourner telle quelle
		type TestStruct struct {
			Field1 string
			Field2 int
		}

		testStruct := TestStruct{
			Field1: "test",
			Field2: 42,
		}

		result := builder.CreateConstraintFromAST(testStruct)

		// Le résultat devrait être la même structure
		if resultStruct, ok := result.(TestStruct); ok {
			if resultStruct.Field1 != "test" {
				t.Errorf("Attendu Field1 'test', obtenu '%s'", resultStruct.Field1)
			}
			if resultStruct.Field2 != 42 {
				t.Errorf("Attendu Field2 42, obtenu %d", resultStruct.Field2)
			}
		} else {
			t.Error("Le résultat devrait être un TestStruct")
		}
	})

	t.Run("avec string", func(t *testing.T) {
		// Test avec un string - devrait le retourner tel quel
		constraintString := "simple string constraint"

		result := builder.CreateConstraintFromAST(constraintString)

		if result != constraintString {
			t.Errorf("Attendu '%s', obtenu '%v'", constraintString, result)
		}
	})

	t.Run("avec nombre", func(t *testing.T) {
		// Test avec un nombre - devrait le retourner tel quel
		constraintNumber := 42

		result := builder.CreateConstraintFromAST(constraintNumber)

		if result != constraintNumber {
			t.Errorf("Attendu %d, obtenu %v", constraintNumber, result)
		}
	})

	t.Run("avec boolean", func(t *testing.T) {
		// Test avec un boolean - devrait le retourner tel quel
		constraintBool := true

		result := builder.CreateConstraintFromAST(constraintBool)

		if result != constraintBool {
			t.Errorf("Attendu %t, obtenu %v", constraintBool, result)
		}
	})

	t.Run("avec slice", func(t *testing.T) {
		// Test avec un slice - devrait le retourner tel quel
		constraintSlice := []string{"a", "b", "c"}

		result := builder.CreateConstraintFromAST(constraintSlice)

		if resultSlice, ok := result.([]string); ok {
			if len(resultSlice) != 3 {
				t.Errorf("Attendu slice de longueur 3, obtenu %d", len(resultSlice))
			}
			if resultSlice[0] != "a" || resultSlice[1] != "b" || resultSlice[2] != "c" {
				t.Errorf("Attendu [a, b, c], obtenu %v", resultSlice)
			}
		} else {
			t.Error("Le résultat devrait être un []string")
		}
	})

	t.Run("avec nil", func(t *testing.T) {
		// Test avec nil - devrait retourner nil
		result := builder.CreateConstraintFromAST(nil)

		if result != nil {
			t.Errorf("Attendu nil, obtenu %v", result)
		}
	})

	t.Run("avec constraint.BinaryOperation", func(t *testing.T) {
		// Test avec une BinaryOperation - devrait la retourner telle quelle
		binOp := constraint.BinaryOperation{
			Type:     "binary_operation",
			Left:     "left_value",
			Operator: "!=",
			Right:    "right_value",
		}

		result := builder.CreateConstraintFromAST(binOp)

		if resultBinOp, ok := result.(constraint.BinaryOperation); ok {
			if resultBinOp.Operator != "!=" {
				t.Errorf("Attendu opérateur '!=', obtenu '%s'", resultBinOp.Operator)
			}
			if resultBinOp.Type != "binary_operation" {
				t.Errorf("Attendu type 'binary_operation', obtenu '%s'", resultBinOp.Type)
			}
		} else {
			t.Error("Le résultat devrait être un constraint.BinaryOperation")
		}
	})
}
