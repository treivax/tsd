// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
)

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
