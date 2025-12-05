// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
)

// ValidateFieldAccess vérifie qu'un accès aux champs est valide dans une expression donnée
func ValidateFieldAccess(program Program, fieldAccess FieldAccess, expressionIndex int) error {
	if expressionIndex >= len(program.Expressions) {
		return fmt.Errorf("index d'expression invalide: %d", expressionIndex)
	}

	expr := program.Expressions[expressionIndex]

	// Trouver le type de l'objet dans l'expression spécifiée
	var objectType string

	// Check new multi-pattern syntax first
	if len(expr.Patterns) > 0 {
		for _, pattern := range expr.Patterns {
			for _, variable := range pattern.Variables {
				if variable.Name == fieldAccess.Object {
					objectType = variable.DataType
					break
				}
			}
			if objectType != "" {
				break
			}
		}
	} else {
		// Old single-pattern syntax (backward compatibility)
		for _, variable := range expr.Set.Variables {
			if variable.Name == fieldAccess.Object {
				objectType = variable.DataType
				break
			}
		}
	}

	if objectType == "" {
		return fmt.Errorf("variable non trouvée: %s dans l'expression %d", fieldAccess.Object, expressionIndex+1)
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

	return fmt.Errorf("champ %s non trouvé dans le type %s", fieldAccess.Field, objectType)
}

// ValidateConstraintFieldAccess parcourt récursivement les contraintes pour valider les accès aux champs
func ValidateConstraintFieldAccess(program Program, constraint interface{}, expressionIndex int) error {
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
			return validateFieldAccessInOperands(program, c, expressionIndex)
		case ConstraintTypeLogicalExpr:
			return validateFieldAccessInLogicalExpr(program, c, expressionIndex)
		}
	}
	return nil
}

// validateFieldAccessInOperands validates field access in left/right operands
func validateFieldAccessInOperands(program Program, c map[string]interface{}, expressionIndex int) error {
	if left := c["left"]; left != nil {
		if err := ValidateConstraintFieldAccess(program, left, expressionIndex); err != nil {
			return err
		}
	}
	if right := c["right"]; right != nil {
		if err := ValidateConstraintFieldAccess(program, right, expressionIndex); err != nil {
			return err
		}
	}
	return nil
}

// validateFieldAccessInLogicalExpr validates field access in logical expressions
func validateFieldAccessInLogicalExpr(program Program, c map[string]interface{}, expressionIndex int) error {
	if left := c["left"]; left != nil {
		if err := ValidateConstraintFieldAccess(program, left, expressionIndex); err != nil {
			return err
		}
	}
	if operations, ok := c["operations"].([]interface{}); ok {
		for _, op := range operations {
			if opMap, ok := op.(map[string]interface{}); ok {
				if right := opMap["right"]; right != nil {
					if err := ValidateConstraintFieldAccess(program, right, expressionIndex); err != nil {
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
		return "", fmt.Errorf("index d'expression invalide: %d", expressionIndex)
	}

	expr := program.Expressions[expressionIndex]

	// Trouver le type de l'objet
	var objectType string

	// Check new multi-pattern syntax first
	if len(expr.Patterns) > 0 {
		for _, pattern := range expr.Patterns {
			for _, variable := range pattern.Variables {
				if variable.Name == object {
					objectType = variable.DataType
					break
				}
			}
			if objectType != "" {
				break
			}
		}
	} else {
		// Old single-pattern syntax (backward compatibility)
		for _, variable := range expr.Set.Variables {
			if variable.Name == object {
				objectType = variable.DataType
				break
			}
		}
	}

	if objectType == "" {
		return "", fmt.Errorf("variable non trouvée: %s", object)
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

	return "", fmt.Errorf("champ %s non trouvé dans le type %s", field, objectType)
}
