// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
	"strings"
)

// ValidatePrimaryKeyField vérifie qu'un champ peut être utilisé comme clé primaire.
// Les champs de clé primaire doivent être de type primitif (string, number, bool).
func ValidatePrimaryKeyField(field Field, typeName string) error {
	if !field.IsPrimaryKey {
		return nil
	}

	// Vérifier que le type est primitif
	if IsValidPrimitiveType(field.Type) {
		return nil
	}

	return fmt.Errorf(
		"type '%s', champ '%s': les champs de clé primaire doivent être de type primitif (string, number, bool), reçu '%s'",
		typeName, field.Name, field.Type)
}

// ValidateTypePrimaryKey valide la cohérence de la clé primaire d'un type.
func ValidateTypePrimaryKey(typeDef TypeDefinition) error {
	pkFields := typeDef.GetPrimaryKeyFields()

	// Pas de clé primaire définie = OK (ID sera généré par hash)
	if len(pkFields) == 0 {
		return nil
	}

	// Valider chaque champ de clé primaire
	for _, field := range pkFields {
		if err := ValidatePrimaryKeyField(field, typeDef.Name); err != nil {
			return err
		}
	}

	return nil
}

// ValidateFactPrimaryKey valide qu'un fait respecte les contraintes de clé primaire.
// - Le fait ne doit JAMAIS définir le champ '_id_' manuellement
// - Tous les champs de clé primaire doivent être fournis
func ValidateFactPrimaryKey(fact Fact, typeDef TypeDefinition) error {
	// Vérifier que _id_ n'est JAMAIS défini manuellement
	for _, factField := range fact.Fields {
		if factField.Name == FieldNameInternalID {
			return fmt.Errorf(
				"fait de type '%s': le champ '%s' est réservé au système et ne peut pas être défini manuellement",
				fact.TypeName,
				FieldNameInternalID,
			)
		}
	}

	// Si le type a une clé primaire, vérifier que tous les champs PK sont fournis
	pkFields := typeDef.GetPrimaryKeyFields()
	if len(pkFields) == 0 {
		return nil // Pas de clé primaire définie
	}

	// Créer une map des champs fournis dans le fait
	providedFields := make(map[string]bool)
	for _, factField := range fact.Fields {
		providedFields[factField.Name] = true
	}

	// Vérifier que chaque champ PK est fourni
	var missingFields []string
	for _, pkField := range pkFields {
		if !providedFields[pkField.Name] {
			missingFields = append(missingFields, pkField.Name)
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf(
			"fait de type '%s': champs de clé primaire manquants: %s",
			fact.TypeName, strings.Join(missingFields, ", "))
	}

	return nil
}

// ValidateFactPrimaryKeyValues vérifie que les valeurs des champs de clé primaire sont non-nulles.
// Note: La validation diffère selon les types pour respecter les sémantiques:
// - string vide "" est invalide (pas d'identifiant significatif)
// - number 0 est valide (valeur numérique légitime)
// - bool false est valide (état booléen légitime)
func ValidateFactPrimaryKeyValues(fact Fact, typeDef TypeDefinition) error {
	pkFieldNames := typeDef.GetPrimaryKeyFieldNames()
	if len(pkFieldNames) == 0 {
		return nil // Pas de clé primaire
	}

	// Créer une map des champs PK
	pkFieldMap := make(map[string]bool)
	for _, name := range pkFieldNames {
		pkFieldMap[name] = true
	}

	// Vérifier chaque champ du fait
	for _, factField := range fact.Fields {
		if !pkFieldMap[factField.Name] {
			continue // Pas un champ PK
		}

		// Vérifier que la valeur n'est pas nulle/vide
		if factField.Value.Value == nil {
			return fmt.Errorf(
				"fait de type '%s': le champ de clé primaire '%s' ne peut pas être nul",
				fact.TypeName, factField.Name)
		}

		// Pour les strings, vérifier qu'elles ne sont pas vides
		if factField.Value.Type == ValueTypeString || factField.Value.Type == ValueTypeIdentifier {
			if strVal, ok := factField.Value.Value.(string); ok && strVal == "" {
				return fmt.Errorf(
					"fait de type '%s': le champ de clé primaire '%s' ne peut pas être vide",
					fact.TypeName, factField.Name)
			}
		}
	}

	return nil
}
