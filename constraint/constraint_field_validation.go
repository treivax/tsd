// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
)

// ValidateFieldAccess vérifie qu'un accès aux champs est valide dans une expression donnée
func ValidateFieldAccess(program Program, fieldAccess FieldAccess, expressionIndex int) error {
	// Validate inputs
	if err := validateInputNotNil(map[string]interface{}{
		"fieldAccess": fieldAccess,
	}); err != nil {
		return err
	}

	if expressionIndex >= len(program.Expressions) {
		return fmt.Errorf("invalid expression index: %d", expressionIndex)
	}

	expr := program.Expressions[expressionIndex]

	// Use helper to find variable type
	objectType, err := findVariableType(expr, fieldAccess.Object)
	if err != nil {
		return fmt.Errorf("in expression %d: %v", expressionIndex+1, err)
	}

	// Vérifier que le champ existe dans le type
	fields, err := GetTypeFields(program, objectType)
	if err != nil {
		return err
	}

	for _, field := range fields {
		if field.Name == fieldAccess.Field {
			return nil // Champ trouvé
		}
	}

	return fmt.Errorf("field %s not found in type %s",
		sanitizeForLog(fieldAccess.Field, 50), sanitizeForLog(objectType, 50))
}

// ValidateConstraintFieldAccess parcourt récursivement les contraintes pour valider les accès aux champs
func ValidateConstraintFieldAccess(program Program, constraint interface{}, expressionIndex int) error {
	return validateConstraintFieldAccessWithDepth(program, constraint, expressionIndex, 0)
}

// validateConstraintFieldAccessWithDepth validates field access with recursion depth tracking
func validateConstraintFieldAccessWithDepth(program Program, constraint interface{}, expressionIndex int, depth int) error {
	// Prevent stack overflow
	if depth > MaxValidationDepth {
		return fmt.Errorf("maximum validation depth exceeded (%d)", MaxValidationDepth)
	}

	switch c := constraint.(type) {
	case map[string]interface{}:
		constraintType, ok := c["type"].(string)
		if !ok {
			return nil
		}

		switch constraintType {
		case ConstraintTypeFieldAccess:
			object, objOk := c["object"].(string)
			field, fieldOk := c["field"].(string)
			if objOk && fieldOk {
				fieldAccess := FieldAccess{
					Type:   ConstraintTypeFieldAccess,
					Object: object,
					Field:  field,
				}
				return ValidateFieldAccess(program, fieldAccess, expressionIndex)
			}
		case ConstraintTypeComparison, ConstraintTypeBinaryOp:
			return validateFieldAccessInOperands(program, c, expressionIndex, depth)
		case ConstraintTypeLogicalExpr:
			return validateFieldAccessInLogicalExpr(program, c, expressionIndex, depth)
		}
	}
	return nil
}

// validateFieldAccessInOperands validates field access in left/right operands
func validateFieldAccessInOperands(program Program, c map[string]interface{}, expressionIndex int, depth int) error {
	if left := c["left"]; left != nil {
		if err := validateConstraintFieldAccessWithDepth(program, left, expressionIndex, depth+1); err != nil {
			return err
		}
	}
	if right := c["right"]; right != nil {
		if err := validateConstraintFieldAccessWithDepth(program, right, expressionIndex, depth+1); err != nil {
			return err
		}
	}
	return nil
}

// validateFieldAccessInLogicalExpr validates field access in logical expressions
func validateFieldAccessInLogicalExpr(program Program, c map[string]interface{}, expressionIndex int, depth int) error {
	if left := c["left"]; left != nil {
		if err := validateConstraintFieldAccessWithDepth(program, left, expressionIndex, depth+1); err != nil {
			return err
		}
	}
	if operations, ok := c["operations"].([]interface{}); ok {
		for _, op := range operations {
			if opMap, ok := op.(map[string]interface{}); ok {
				if right := opMap["right"]; right != nil {
					if err := validateConstraintFieldAccessWithDepth(program, right, expressionIndex, depth+1); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// GetFieldType retourne le type d'un champ spécifique d'un objet dans une expression
func GetFieldType(program Program, object string, field string, expressionIndex int) (string, error) {
	if expressionIndex >= len(program.Expressions) {
		return "", fmt.Errorf("invalid expression index: %d", expressionIndex)
	}

	expr := program.Expressions[expressionIndex]

	// Use helper to find variable type
	objectType, err := findVariableType(expr, object)
	if err != nil {
		return "", err
	}

	// Trouver le type du champ dans la définition du type
	fields, err := GetTypeFields(program, objectType)
	if err != nil {
		return "", err
	}

	for _, f := range fields {
		if f.Name == field {
			return f.Type, nil
		}
	}

	return "", fmt.Errorf("field %s not found in type %s",
		sanitizeForLog(field, 50), sanitizeForLog(objectType, 50))
}
