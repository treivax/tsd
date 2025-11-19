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
		// Extraire les variables utilisées dans l'argument
		vars := extractVariablesFromArg(arg)
		for _, varName := range vars {
			if !availableVars[varName] {
				return fmt.Errorf("action %s: argument contient la variable '%s' qui ne correspond à aucune variable de l'expression", action.Job.Name, varName)
			}
		}
	}

	return nil
}

// extractVariablesFromArg extrait les noms de variables utilisées dans un argument d'action
func extractVariablesFromArg(arg interface{}) []string {
	var vars []string

	// Si c'est une string simple, c'est probablement un nom de variable
	if str, ok := arg.(string); ok {
		vars = append(vars, str)
		return vars
	}

	// Si c'est un objet (map), extraire les variables selon le type
	if argMap, ok := arg.(map[string]interface{}); ok {
		argType, _ := argMap["type"].(string)
		switch argType {
		case "fieldAccess":
			if object, ok := argMap["object"].(string); ok {
				vars = append(vars, object)
			}
		case "string":
			// Les string literals ne contiennent pas de variables
		case "number":
			// Les number literals ne contiennent pas de variables
		default:
			// Pour d'autres types, on peut chercher récursivement
			// mais pour l'instant on ignore
		}
	}

	return vars
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
// ValidateTypeCompatibility validates type compatibility within constraints
func ValidateTypeCompatibility(program Program, constraint interface{}, expressionIndex int) error {
	constraintMap, ok := constraint.(map[string]interface{})
	if !ok {
		return nil
	}

	constraintType, ok := constraintMap["type"].(string)
	if !ok {
		return nil
	}

	switch constraintType {
	case "comparison":
		return validateComparisonConstraint(program, constraintMap, expressionIndex)
	case "logicalExpr":
		return validateLogicalExpressionConstraint(program, constraintMap, expressionIndex)
	case "binaryOp":
		return validateBinaryOpConstraint(program, constraintMap, expressionIndex)
	}
	return nil
}

// validateComparisonConstraint handles comparison constraint validation
func validateComparisonConstraint(program Program, c map[string]interface{}, expressionIndex int) error {
	left := c["left"]
	right := c["right"]

	if left == nil || right == nil {
		return nil
	}

	// Validate type compatibility between operands
	if err := validateOperandTypeCompatibility(program, left, right, expressionIndex); err != nil {
		return err
	}

	// Recursive validation for operands
	if err := ValidateTypeCompatibility(program, left, expressionIndex); err != nil {
		return err
	}
	if err := ValidateTypeCompatibility(program, right, expressionIndex); err != nil {
		return err
	}

	return nil
}

// validateOperandTypeCompatibility checks if two operands have compatible types
func validateOperandTypeCompatibility(program Program, left, right interface{}, expressionIndex int) error {
	leftType, err := getOperandType(program, left, expressionIndex)
	if err != nil {
		return err
	}

	rightType, err := getOperandType(program, right, expressionIndex)
	if err != nil {
		return err
	}

	// Check compatibility
	if leftType != "unknown" && rightType != "unknown" && rightType != "variable" {
		if leftType != rightType {
			return fmt.Errorf("incompatibilité de types dans la comparaison: %s vs %s", leftType, rightType)
		}
	}

	return nil
}

// getOperandType determines the type of an operand in a constraint
func getOperandType(program Program, operand interface{}, expressionIndex int) (string, error) {
	operandMap, ok := operand.(map[string]interface{})
	if !ok {
		return GetValueType(operand), nil
	}

	if operandMap["type"] == "fieldAccess" {
		object := operandMap["object"].(string)
		field := operandMap["field"].(string)
		return GetFieldType(program, object, field, expressionIndex)
	}

	return GetValueType(operand), nil
}

// validateLogicalExpressionConstraint handles logical expression validation
func validateLogicalExpressionConstraint(program Program, c map[string]interface{}, expressionIndex int) error {
	if left := c["left"]; left != nil {
		if err := ValidateTypeCompatibility(program, left, expressionIndex); err != nil {
			return err
		}
	}

	operations, ok := c["operations"].([]interface{})
	if !ok {
		return nil
	}

	for _, op := range operations {
		opMap, ok := op.(map[string]interface{})
		if !ok {
			continue
		}

		if right := opMap["right"]; right != nil {
			if err := ValidateTypeCompatibility(program, right, expressionIndex); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateBinaryOpConstraint handles binary operation validation
func validateBinaryOpConstraint(program Program, c map[string]interface{}, expressionIndex int) error {
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

	// Validation des faits
	err = ValidateFacts(program)
	if err != nil {
		return fmt.Errorf("erreur validation faits: %v", err)
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

	fmt.Printf("✓ Programme valide avec %d type(s), %d expression(s) et %d fait(s)\n", len(program.Types), len(program.Expressions), len(program.Facts))
	return nil
}

// ValidateFacts vérifie que tous les faits parsés sont cohérents avec les définitions de types
func ValidateFacts(program Program) error {
	definedTypes := make(map[string]TypeDefinition)
	for _, typeDef := range program.Types {
		definedTypes[typeDef.Name] = typeDef
	}

	for i, fact := range program.Facts {
		// Vérifier que le type du fait existe
		typeDef, exists := definedTypes[fact.TypeName]
		if !exists {
			return fmt.Errorf("fait %d: type non défini: %s", i+1, fact.TypeName)
		}

		// Créer une map des champs définis pour ce type
		definedFields := make(map[string]string)
		for _, field := range typeDef.Fields {
			definedFields[field.Name] = field.Type
		}

		// Vérifier chaque champ du fait
		for j, factField := range fact.Fields {
			// Vérifier que le champ existe dans le type
			expectedType, exists := definedFields[factField.Name]
			if !exists {
				return fmt.Errorf("fait %d, champ %d: champ '%s' non défini dans le type %s", i+1, j+1, factField.Name, fact.TypeName)
			}

			// Vérifier la compatibilité du type de la valeur
			err := ValidateFactFieldType(factField.Value, expectedType, fact.TypeName, factField.Name)
			if err != nil {
				return fmt.Errorf("fait %d, champ %d: %v", i+1, j+1, err)
			}
		}
	}

	return nil
}

// ValidateFactFieldType vérifie que la valeur d'un champ de fait correspond au type attendu
func ValidateFactFieldType(value FactValue, expectedType, typeName, fieldName string) error {
	switch expectedType {
	case "string":
		if value.Type != "string" && value.Type != "identifier" {
			return fmt.Errorf("champ '%s' du type %s attend une valeur string, reçu %s", fieldName, typeName, value.Type)
		}
	case "number":
		if value.Type != "number" {
			return fmt.Errorf("champ '%s' du type %s attend une valeur number, reçu %s", fieldName, typeName, value.Type)
		}
	case "bool", "boolean":
		if value.Type != "boolean" {
			return fmt.Errorf("champ '%s' du type %s attend une valeur boolean, reçu %s", fieldName, typeName, value.Type)
		}
	default:
		// Type non reconnu, on accepte pour l'instant
		return nil
	}
	return nil
}

// ConvertFactsToReteFormat convertit les faits parsés par la grammaire vers le format attendu par le réseau RETE
func ConvertFactsToReteFormat(program Program) []map[string]interface{} {
	var reteFacts []map[string]interface{}

	for i, fact := range program.Facts {
		reteFact := map[string]interface{}{
			"reteType": fact.TypeName, // Type RETE (ex: "Balance")
		}

		// Variables pour gérer l'ID du fait
		var factID string
		hasExplicitID := false

		// Convertir les champs
		for _, field := range fact.Fields {
			var convertedValue interface{}
			switch field.Value.Type {
			case "string":
				if strVal, ok := field.Value.Value.(map[string]interface{}); ok {
					convertedValue = strVal["value"]
				} else {
					convertedValue = field.Value.Value
				}
			case "number":
				if numVal, ok := field.Value.Value.(map[string]interface{}); ok {
					convertedValue = numVal["value"]
				} else {
					convertedValue = field.Value.Value
				}
			case "boolean":
				if boolVal, ok := field.Value.Value.(map[string]interface{}); ok {
					convertedValue = boolVal["value"]
				} else {
					convertedValue = field.Value.Value
				}
			case "identifier":
				// Les identifiants non-quotés sont traités comme des strings
				convertedValue = field.Value.Value
			default:
				convertedValue = field.Value.Value
			}

			// Ajouter le champ au fact, y compris l'ID s'il existe
			reteFact[field.Name] = convertedValue

			// Vérifier si c'est un champ ID
			if field.Name == "id" {
				factID = convertedValue.(string)
				hasExplicitID = true
			}
		}

		// Générer un ID si pas fourni explicitement
		if !hasExplicitID {
			factID = fmt.Sprintf("parsed_fact_%d", i+1)
			reteFact["id"] = factID
		}

		// Définir l'ID du fait (nécessaire pour le réseau RETE)
		reteFact["id"] = factID

		// CORRECTION CRITIQUE: Assurer que le type RETE est toujours préservé
		reteFact["reteType"] = fact.TypeName

		reteFacts = append(reteFacts, reteFact)
	}

	return reteFacts
}
