package constraint

import (
	"encoding/json"
	"fmt"
)

// Fonctions utilitaires pour traiter l'AST du parser de contraintes

// ValidateTypes vérifie que tous les types référencés dans les expressions sont définis
func ValidateTypes(program Program) error {
	definedTypes := make(map[string]bool)
	for _, typeDef := range program.Types {
		definedTypes[typeDef.Name] = true
	}

	// Vérifier les variables typées dans toutes les expressions
	for i, expression := range program.Expressions {
		for _, variable := range expression.Set.Variables {
			if !definedTypes[variable.DataType] {
				return fmt.Errorf("expression %d: type non défini: %s pour la variable %s", i+1, variable.DataType, variable.Name)
			}
		}
	}

	return nil
}

// GetTypeFields retourne les champs d'un type donné
func GetTypeFields(program Program, typeName string) ([]Field, error) {
	for _, typeDef := range program.Types {
		if typeDef.Name == typeName {
			return typeDef.Fields, nil
		}
	}
	return nil, fmt.Errorf("type non trouvé: %s", typeName)
}

// ValidateFieldAccess vérifie qu'un accès aux champs est valide dans une expression donnée
func ValidateFieldAccess(program Program, fieldAccess FieldAccess, expressionIndex int) error {
	if expressionIndex >= len(program.Expressions) {
		return fmt.Errorf("index d'expression invalide: %d", expressionIndex)
	}

	// Trouver le type de l'objet dans l'expression spécifiée
	var objectType string
	for _, variable := range program.Expressions[expressionIndex].Set.Variables {
		if variable.Name == fieldAccess.Object {
			objectType = variable.DataType
			break
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

// ValidateAction vérifie qu'une action est valide dans le contexte d'une expression
func ValidateAction(program Program, action Action, expressionIndex int) error {
	if expressionIndex >= len(program.Expressions) {
		return fmt.Errorf("index d'expression invalide: %d", expressionIndex)
	}

	expression := program.Expressions[expressionIndex]

	// Créer une map des variables disponibles dans l'expression
	availableVars := make(map[string]bool)
	for _, variable := range expression.Set.Variables {
		availableVars[variable.Name] = true
	}

	// Vérifier que tous les arguments de l'action référencent des variables valides
	for _, arg := range action.Job.Args {
		if !availableVars[arg] {
			return fmt.Errorf("action %s: argument '%s' ne correspond à aucune variable de l'expression", action.Job.Name, arg)
		}
	}

	return nil
}

// GetFieldType retourne le type d'un champ spécifique d'un objet dans une expression
func GetFieldType(program Program, object string, field string, expressionIndex int) (string, error) {
	if expressionIndex >= len(program.Expressions) {
		return "", fmt.Errorf("index d'expression invalide: %d", expressionIndex)
	}

	// Trouver le type de l'objet
	var objectType string
	for _, variable := range program.Expressions[expressionIndex].Set.Variables {
		if variable.Name == object {
			objectType = variable.DataType
			break
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

// GetValueType retourne le type d'une valeur dans l'AST
func GetValueType(value interface{}) string {
	switch v := value.(type) {
	case map[string]interface{}:
		valueType, ok := v["type"].(string)
		if !ok {
			return "unknown"
		}
		switch valueType {
		case "number":
			return "number"
		case "string":
			return "string"
		case "boolean":
			return "bool"
		case "variable":
			// Pour les variables comme "true", "false" qui sont parsées comme variables
			name, ok := v["name"].(string)
			if ok {
				switch name {
				case "true", "false":
					return "bool"
				}
			}
			return "variable" // Type non déterminable sans contexte
		}
	}
	return "unknown"
}

// ValidateTypeCompatibility vérifie la compatibilité des types dans les comparaisons
func ValidateTypeCompatibility(program Program, constraint interface{}, expressionIndex int) error {
	switch c := constraint.(type) {
	case map[string]interface{}:
		constraintType, ok := c["type"].(string)
		if !ok {
			return nil
		}

		switch constraintType {
		case "comparison":
			left := c["left"]
			right := c["right"]

			if left != nil && right != nil {
				var leftType, rightType string
				var err error

				// Déterminer le type de la partie gauche
				if leftMap, ok := left.(map[string]interface{}); ok {
					if leftMap["type"] == "fieldAccess" {
						object := leftMap["object"].(string)
						field := leftMap["field"].(string)
						leftType, err = GetFieldType(program, object, field, expressionIndex)
						if err != nil {
							return err
						}
					} else {
						leftType = GetValueType(left)
					}
				}

				// Déterminer le type de la partie droite
				if rightMap, ok := right.(map[string]interface{}); ok {
					if rightMap["type"] == "fieldAccess" {
						object := rightMap["object"].(string)
						field := rightMap["field"].(string)
						rightType, err = GetFieldType(program, object, field, expressionIndex)
						if err != nil {
							return err
						}
					} else {
						rightType = GetValueType(right)
					}
				} else {
					rightType = GetValueType(right)
				}

				// Vérifier la compatibilité
				if leftType != "unknown" && rightType != "unknown" && rightType != "variable" {
					if leftType != rightType {
						return fmt.Errorf("incompatibilité de types dans la comparaison: %s vs %s", leftType, rightType)
					}
				}
			}

			// Validation récursive pour les opérands
			if left != nil {
				if err := ValidateTypeCompatibility(program, left, expressionIndex); err != nil {
					return err
				}
			}
			if right != nil {
				if err := ValidateTypeCompatibility(program, right, expressionIndex); err != nil {
					return err
				}
			}

		case "logicalExpr":
			if left := c["left"]; left != nil {
				if err := ValidateTypeCompatibility(program, left, expressionIndex); err != nil {
					return err
				}
			}
			if operations, ok := c["operations"].([]interface{}); ok {
				for _, op := range operations {
					if opMap, ok := op.(map[string]interface{}); ok {
						if right := opMap["right"]; right != nil {
							if err := ValidateTypeCompatibility(program, right, expressionIndex); err != nil {
								return err
							}
						}
					}
				}
			}

		case "binaryOp":
			if left := c["left"]; left != nil {
				if err := ValidateTypeCompatibility(program, left, expressionIndex); err != nil {
					return err
				}
			}
			if right := c["right"]; right != nil {
				if err := ValidateTypeCompatibility(program, right, expressionIndex); err != nil {
					return err
				}
			}
		}
	}
	return nil
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
		case "fieldAccess":
			object, objOk := c["object"].(string)
			field, fieldOk := c["field"].(string)
			if objOk && fieldOk {
				fieldAccess := FieldAccess{
					Type:   "fieldAccess",
					Object: object,
					Field:  field,
				}
				return ValidateFieldAccess(program, fieldAccess, expressionIndex)
			}
		case "comparison":
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
		case "logicalExpr":
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
		case "binaryOp":
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
		}
	}
	return nil
}

// ValidateProgram effectue une validation complète du programme parsé
func ValidateProgram(result interface{}) error {
	// Convertir le résultat en structure Program
	jsonData, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("erreur conversion JSON: %v", err)
	}

	var program Program
	err = json.Unmarshal(jsonData, &program)
	if err != nil {
		return fmt.Errorf("erreur parsing JSON: %v", err)
	}

	// Validation des types
	err = ValidateTypes(program)
	if err != nil {
		return fmt.Errorf("erreur validation types: %v", err)
	}

	// Validation des accès aux champs dans les contraintes
	for i, expression := range program.Expressions {
		if expression.Constraints != nil {
			err = ValidateConstraintFieldAccess(program, expression.Constraints, i)
			if err != nil {
				return fmt.Errorf("erreur validation champs dans l'expression %d: %v", i+1, err)
			}
		}
	}

	// Validation des types dans les comparaisons
	for i, expression := range program.Expressions {
		if expression.Constraints != nil {
			err = ValidateTypeCompatibility(program, expression.Constraints, i)
			if err != nil {
				return fmt.Errorf("erreur validation types dans l'expression %d: %v", i+1, err)
			}
		}
	}

	// Validation des actions (maintenant obligatoires)
	for i, expression := range program.Expressions {
		if expression.Action != nil {
			err = ValidateAction(program, *expression.Action, i)
			if err != nil {
				return fmt.Errorf("erreur validation action dans l'expression %d: %v", i+1, err)
			}
		} else {
			// Avec la nouvelle grammaire, cette condition ne devrait plus arriver
			return fmt.Errorf("action manquante dans l'expression %d: chaque règle doit avoir une action définie", i+1)
		}
	}

	fmt.Printf("✓ Programme valide avec %d type(s) et %d expression(s)\n", len(program.Types), len(program.Expressions))
	return nil
}
