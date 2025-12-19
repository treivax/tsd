// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
)

// ValidateTypes vérifie que tous les types référencés dans les expressions sont définis
// et valide également la cohérence des clés primaires.
func ValidateTypes(program Program) error {
	definedTypes := make(map[string]bool)
	for _, typeDef := range program.Types {
		definedTypes[typeDef.Name] = true

		// Valider que _id_ n'est pas utilisé comme nom de champ
		for _, field := range typeDef.Fields {
			if field.Name == FieldNameInternalID {
				return fmt.Errorf(
					"type '%s': le champ '%s' est réservé au système et ne peut pas être utilisé",
					typeDef.Name,
					FieldNameInternalID,
				)
			}
		}

		// Valider la clé primaire du type
		if err := ValidateTypePrimaryKey(typeDef); err != nil {
			return err
		}
	}

	// Vérifier les variables typées dans toutes les expressions
	for i, expression := range program.Expressions {
		for _, variable := range expression.Set.Variables {
			if !definedTypes[variable.DataType] {
				return fmt.Errorf("expression %d: undefined type: %s for variable %s",
					i+1, sanitizeForLog(variable.DataType, 50), sanitizeForLog(variable.Name, 50))
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
	return nil, fmt.Errorf("type not found: %s", sanitizeForLog(typeName, 50))
}
