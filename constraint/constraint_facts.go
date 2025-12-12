// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
)

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
	case ValueTypeString:
		if value.Type != ValueTypeString && value.Type != ValueTypeIdentifier {
			return fmt.Errorf("champ '%s' du type %s attend une valeur string, reçu %s", fieldName, typeName, value.Type)
		}
	case ValueTypeNumber:
		if value.Type != ValueTypeNumber {
			return fmt.Errorf("champ '%s' du type %s attend une valeur number, reçu %s", fieldName, typeName, value.Type)
		}
	case ValueTypeBool, ValueTypeBoolean:
		if value.Type != ValueTypeBoolean {
			return fmt.Errorf("champ '%s' du type %s attend une valeur boolean, reçu %s", fieldName, typeName, value.Type)
		}
	default:
		// Type non primitif : vérifier si c'est un type valide défini
		// Les types personnalisés et futurs types primitifs doivent être validés explicitement
		if !ValidPrimitiveTypes[expectedType] {
			// Type personnalisé ou non standard accepté pour extensibilité
			// TODO: Valider que le type personnalisé existe dans le programme
			return nil
		}
		// Type primitif non géré dans le switch - erreur de validation
		return fmt.Errorf("champ '%s' du type %s : type attendu '%s' non pris en charge", fieldName, typeName, expectedType)
	}
	return nil
}

// ConvertFactsToReteFormat convertit les faits parsés par la grammaire vers le format attendu par le réseau RETE
func ConvertFactsToReteFormat(program Program) []map[string]interface{} {
	var reteFacts []map[string]interface{}

	for i, fact := range program.Facts {
		reteFact := map[string]interface{}{
			FieldNameReteType: fact.TypeName, // Type RETE (ex: "Balance")
		}

		// Convertir les champs
		convertFactFieldsToMap(fact.Fields, reteFact)

		// Gérer l'ID du fait
		factID := ensureFactID(reteFact, i)
		reteFact[FieldNameID] = factID

		// CORRECTION CRITIQUE: Assurer que le type RETE est toujours préservé
		reteFact[FieldNameReteType] = fact.TypeName

		reteFacts = append(reteFacts, reteFact)
	}

	return reteFacts
}

// convertFactFieldsToMap converts fact fields to a map, handling value conversion
func convertFactFieldsToMap(fields []FactField, targetMap map[string]interface{}) {
	for _, field := range fields {
		convertedValue := convertFactFieldValue(field.Value)
		targetMap[field.Name] = convertedValue
	}
}

// ensureFactID ensures a fact has an ID, generating one if necessary
func ensureFactID(reteFact map[string]interface{}, factIndex int) string {
	// Check if ID was explicitly provided in the fields
	if id, exists := reteFact[FieldNameID]; exists {
		if idStr, ok := id.(string); ok && idStr != "" {
			return idStr
		}
	}

	// Generate ID if not provided
	return fmt.Sprintf("parsed_fact_%d", factIndex+1)
}

// convertFactFieldValue converts a fact field value to its appropriate Go type
func convertFactFieldValue(value FactValue) interface{} {
	switch value.Type {
	case ValueTypeString, ValueTypeNumber, ValueTypeBoolean:
		if valMap, ok := value.Value.(map[string]interface{}); ok {
			return valMap["value"]
		}
		return value.Value
	case ValueTypeIdentifier:
		// Les identifiants non-quotés sont traités comme des strings
		return value.Value
	default:
		return value.Value
	}
}
