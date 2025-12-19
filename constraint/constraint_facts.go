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
	return validateFactFieldTypeValue(value, expectedType, typeName, fieldName)
}

// validateFactFieldTypeValue performs the actual validation of a fact field type
func validateFactFieldTypeValue(value FactValue, expectedType, typeName, fieldName string) error {
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
		// Type non primitif : accepter les types personnalisés
		// La validation complète des types personnalisés est faite par FactValidator
		// qui a accès au TypeSystem
		if !IsPrimitiveType(expectedType) {
			// Accepter les variableReference et les types personnalisés
			// La résolution et validation complète se fait plus tard
			return nil
		}
		// Type primitif non géré dans le switch - erreur de validation
		return fmt.Errorf("champ '%s' du type %s : type attendu '%s' non pris en charge", fieldName, typeName, expectedType)
	}
	return nil
}

// ConvertFactsToReteFormat convertit les faits parsés par la grammaire vers le format attendu par le réseau RETE.
// Gère les affectations de variables et les références entre faits.
func ConvertFactsToReteFormat(program Program) ([]map[string]interface{}, error) {
	// Normaliser les types de valeurs de faits
	normalizeFactValueTypes(&program)

	// Créer le contexte avec les types
	ctx := NewFactContext(program.Types)

	typeMap := buildTypeMap(program.Types)
	var reteFacts []map[string]interface{}

	// 1. Traiter d'abord les affectations de variables
	for i, assignment := range program.FactAssignments {
		typeDef, exists := typeMap[assignment.Fact.TypeName]
		if !exists {
			return nil, fmt.Errorf("affectation %d: type '%s' non défini", i+1, assignment.Fact.TypeName)
		}

		reteFact := createReteFact(assignment.Fact, typeDef, ctx)
		factID, err := ensureFactID(reteFact, assignment.Fact, typeDef, ctx)
		if err != nil {
			return nil, fmt.Errorf("affectation %d: %v", i+1, err)
		}

		reteFact[FieldNameInternalID] = factID
		reteFact[FieldNameReteType] = assignment.Fact.TypeName

		// Enregistrer la variable dans le contexte
		ctx.RegisterVariable(assignment.Variable, factID)

		reteFacts = append(reteFacts, reteFact)
	}

	// 2. Traiter les faits normaux (peuvent référencer les variables)
	for i, fact := range program.Facts {
		typeDef, exists := typeMap[fact.TypeName]
		if !exists {
			return nil, fmt.Errorf("fait %d: type '%s' non défini", i+1, fact.TypeName)
		}

		reteFact := createReteFact(fact, typeDef, ctx)
		factID, err := ensureFactID(reteFact, fact, typeDef, ctx)
		if err != nil {
			return nil, fmt.Errorf("fait %d: %v", i+1, err)
		}

		reteFact[FieldNameInternalID] = factID
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
// Résout les références de variables vers leurs IDs.
func createReteFact(fact Fact, typeDef TypeDefinition, ctx *FactContext) map[string]interface{} {
	reteFact := map[string]interface{}{
		FieldNameReteType: fact.TypeName,
	}
	convertFactFieldsToMap(fact.Fields, reteFact, ctx)
	return reteFact
}

// convertFactFieldsToMap converts fact fields to a map, handling value conversion and variable resolution.
func convertFactFieldsToMap(fields []FactField, targetMap map[string]interface{}, ctx *FactContext) {
	for _, field := range fields {
		// Si c'est une référence de variable, résoudre vers l'ID
		if field.Value.Type == ValueTypeVariableReference {
			if varName, ok := field.Value.Value.(string); ok && ctx != nil {
				if id, err := ctx.ResolveVariable(varName); err == nil {
					targetMap[field.Name] = id
					continue
				}
			}
		}
		// Sinon, utiliser la valeur unwrapped
		targetMap[field.Name] = field.Value.Unwrap()
	}
}

// ensureFactID generates an internal ID for a fact with support for variable resolution.
// The ID is ALWAYS generated, never provided manually.
func ensureFactID(reteFact map[string]interface{}, fact Fact, typeDef TypeDefinition, ctx *FactContext) (string, error) {
	// Vérifier que _id_ n'a PAS été fourni manuellement
	if _, exists := reteFact[FieldNameInternalID]; exists {
		return "", fmt.Errorf(
			"le champ '%s' ne peut pas être défini manuellement pour le type '%s'",
			FieldNameInternalID,
			fact.TypeName,
		)
	}

	// TOUJOURS générer l'ID avec le contexte
	id, err := GenerateFactID(fact, typeDef, ctx)
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
