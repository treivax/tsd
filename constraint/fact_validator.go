// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
)

// FactValidator valide les faits selon leur définition de type.
type FactValidator struct {
	typeSystem *TypeSystem
}

// NewFactValidator crée un nouveau validateur de faits.
func NewFactValidator(ts *TypeSystem) *FactValidator {
	return &FactValidator{
		typeSystem: ts,
	}
}

// ValidateFact valide un fait complet.
func (fv *FactValidator) ValidateFact(fact Fact) error {
	if !fv.typeSystem.TypeExists(fact.TypeName) {
		return fmt.Errorf(
			"type '%s' non défini",
			fact.TypeName,
		)
	}

	typeDef := fv.typeSystem.types[fact.TypeName]

	if err := fv.validateRequiredFields(fact, typeDef); err != nil {
		return err
	}

	if err := fv.validateFieldDefinitions(fact, typeDef); err != nil {
		return err
	}

	if err := fv.validateFieldValues(fact, typeDef); err != nil {
		return err
	}

	if err := ValidateFactPrimaryKey(fact, typeDef); err != nil {
		return err
	}

	return nil
}

// validateRequiredFields vérifie que tous les champs requis sont présents.
func (fv *FactValidator) validateRequiredFields(fact Fact, typeDef TypeDefinition) error {
	providedFields := make(map[string]bool)
	for _, field := range fact.Fields {
		providedFields[field.Name] = true
	}

	for _, fieldDef := range typeDef.Fields {
		if !providedFields[fieldDef.Name] {
			return fmt.Errorf(
				"fait de type '%s': champ requis '%s' manquant",
				fact.TypeName,
				fieldDef.Name,
			)
		}
	}

	return nil
}

// validateFieldDefinitions vérifie que les champs fournis sont définis.
func (fv *FactValidator) validateFieldDefinitions(fact Fact, typeDef TypeDefinition) error {
	definedFields := make(map[string]Field)
	for _, fieldDef := range typeDef.Fields {
		definedFields[fieldDef.Name] = fieldDef
	}

	for _, factField := range fact.Fields {
		if factField.Name == FieldNameInternalID {
			return fmt.Errorf(
				"fait de type '%s': le champ '%s' est réservé et ne peut pas être défini",
				fact.TypeName,
				FieldNameInternalID,
			)
		}

		if _, exists := definedFields[factField.Name]; !exists {
			return fmt.Errorf(
				"fait de type '%s': champ '%s' non défini dans le type",
				fact.TypeName,
				factField.Name,
			)
		}
	}

	return nil
}

// validateFieldValues vérifie que les valeurs des champs ont le bon type.
func (fv *FactValidator) validateFieldValues(fact Fact, typeDef TypeDefinition) error {
	fieldTypes := make(map[string]string)
	for _, fieldDef := range typeDef.Fields {
		fieldTypes[fieldDef.Name] = fieldDef.Type
	}

	for _, factField := range fact.Fields {
		expectedType := fieldTypes[factField.Name]

		if err := fv.validateFieldValue(factField, expectedType); err != nil {
			return fmt.Errorf(
				"fait de type '%s', champ '%s': %v",
				fact.TypeName,
				factField.Name,
				err,
			)
		}
	}

	return nil
}

// validateFieldValue valide une valeur de champ.
func (fv *FactValidator) validateFieldValue(field FactField, expectedType string) error {
	value := field.Value

	if value.Type == "variableReference" {
		varName, ok := value.Value.(string)
		if !ok {
			return fmt.Errorf("référence de variable invalide")
		}

		if !fv.typeSystem.VariableExists(varName) {
			return fmt.Errorf("variable '%s' non définie", varName)
		}

		varType, _ := fv.typeSystem.GetVariableType(varName)
		if varType != expectedType {
			return fmt.Errorf(
				"type incompatible: attendu '%s', la variable '%s' est de type '%s'",
				expectedType,
				varName,
				varType,
			)
		}

		return nil
	}

	return fv.validatePrimitiveValue(value, expectedType)
}

// validatePrimitiveValue valide une valeur primitive.
func (fv *FactValidator) validatePrimitiveValue(value FactValue, expectedType string) error {
	typeMapping := map[string][]string{
		ValueTypeString:  {ValueTypeString},
		ValueTypeNumber:  {ValueTypeNumber},
		ValueTypeBoolean: {ValueTypeBool, ValueTypeBoolean},
		ValueTypeBool:    {ValueTypeBool, ValueTypeBoolean},
	}

	validTypes, exists := typeMapping[value.Type]
	if !exists {
		return fmt.Errorf("type de valeur '%s' non supporté", value.Type)
	}

	for _, validType := range validTypes {
		if expectedType == validType {
			return nil
		}
	}

	return fmt.Errorf(
		"type incompatible: attendu '%s', reçu '%s'",
		expectedType,
		value.Type,
	)
}
