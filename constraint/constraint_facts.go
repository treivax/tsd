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

		// Valider la clé primaire du fait
		if err := ValidateFactPrimaryKey(fact, typeDef); err != nil {
			return fmt.Errorf("fait %d: %v", i+1, err)
		}

		// Valider les valeurs des champs PK
		if err := ValidateFactPrimaryKeyValues(fact, typeDef); err != nil {
			return fmt.Errorf("fait %d: %v", i+1, err)
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
func ConvertFactsToReteFormat(program Program) ([]map[string]interface{}, error) {
	reteFacts := make([]map[string]interface{}, 0, len(program.Facts))
	typeMap := buildTypeMap(program.Types)

	for i, fact := range program.Facts {
		typeDef, exists := typeMap[fact.TypeName]
		if !exists {
			return nil, fmt.Errorf("fait %d: type '%s' non défini", i+1, fact.TypeName)
		}

		reteFact := createReteFact(fact, typeDef)
		factID, err := ensureFactID(reteFact, fact, typeDef)
		if err != nil {
			return nil, fmt.Errorf("fait %d: %v", i+1, err)
		}
		reteFact[FieldNameID] = factID
		reteFact[FieldNameReteType] = fact.TypeName

		reteFacts = append(reteFacts, reteFact)
	}

	return reteFacts, nil
}

// buildTypeMap crée une map des types pour lookup rapide.
func buildTypeMap(types []TypeDefinition) map[string]TypeDefinition {
	typeMap := make(map[string]TypeDefinition, len(types))
	for _, typeDef := range types {
		typeMap[typeDef.Name] = typeDef
	}
	return typeMap
}

// createReteFact crée un fait RETE avec les champs convertis.
func createReteFact(fact Fact, typeDef TypeDefinition) map[string]interface{} {
	reteFact := map[string]interface{}{
		FieldNameReteType: fact.TypeName,
	}
	convertFactFieldsToMap(fact.Fields, reteFact)
	return reteFact
}

// convertFactFieldsToMap converts fact fields to a map, handling value conversion.
func convertFactFieldsToMap(fields []FactField, targetMap map[string]interface{}) {
	for _, field := range fields {
		targetMap[field.Name] = field.Value.Unwrap()
	}
}

// ensureFactID ensures a fact has an ID, generating one if necessary using primary keys or hash.
func ensureFactID(reteFact map[string]interface{}, fact Fact, typeDef TypeDefinition) (string, error) {
	// Check if ID was explicitly provided (should be prevented by validation)
	if id, exists := reteFact[FieldNameID]; exists {
		if idStr, ok := id.(string); ok && idStr != "" {
			// ID was provided, this should have been caught by validation
			// but we allow it for backward compatibility in some cases
			return idStr, nil
		}
	}

	// Generate ID based on primary key or hash
	id, err := GenerateFactID(fact, typeDef)
	if err != nil {
		return "", fmt.Errorf("génération d'ID pour le fait de type '%s': %v", fact.TypeName, err)
	}

	return id, nil
}

// convertFactFieldValue converts a fact field value to its appropriate Go type.
// Deprecated: Use FactValue.Unwrap() method instead.
// This function is kept for backward compatibility.
func convertFactFieldValue(value FactValue) interface{} {
	return value.Unwrap()
}
