// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

// validateFactFields valide que les champs d'un fait correspondent à sa définition de type
func (ae *ActionExecutor) validateFactFields(typeDef *TypeDefinition, fields map[string]interface{}) error {
	// Vérifier que tous les champs définis sont présents
	for _, fieldDef := range typeDef.Fields {
		value, exists := fields[fieldDef.Name]
		if !exists {
			return fmt.Errorf("champ requis '%s' manquant", fieldDef.Name)
		}

		if err := ae.validateFieldType(fieldDef.Type, value); err != nil {
			return fmt.Errorf("champ '%s': %w", fieldDef.Name, err)
		}
	}

	// Vérifier qu'il n'y a pas de champs non définis
	for fieldName := range fields {
		found := false
		for _, fieldDef := range typeDef.Fields {
			if fieldDef.Name == fieldName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("champ '%s' non défini dans le type", fieldName)
		}
	}

	return nil
}

// validateFieldValue valide qu'une valeur de champ est cohérente avec sa définition
func (ae *ActionExecutor) validateFieldValue(typeDef *TypeDefinition, fieldName string, value interface{}) error {
	// Trouver la définition du champ
	var fieldDef *Field
	for i := range typeDef.Fields {
		if typeDef.Fields[i].Name == fieldName {
			fieldDef = &typeDef.Fields[i]
			break
		}
	}

	if fieldDef == nil {
		return fmt.Errorf("champ '%s' non défini dans le type '%s'", fieldName, typeDef.Name)
	}

	return ae.validateFieldType(fieldDef.Type, value)
}

// validateFieldType valide qu'une valeur correspond au type attendu
func (ae *ActionExecutor) validateFieldType(expectedType string, value interface{}) error {
	switch expectedType {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("type attendu: string, reçu: %T", value)
		}
	case "number":
		if _, ok := toNumber(value); !ok {
			return fmt.Errorf("type attendu: number, reçu: %T", value)
		}
	case "bool":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("type attendu: bool, reçu: %T", value)
		}
	default:
		// Type personnalisé ou non reconnu
		return nil
	}
	return nil
}
