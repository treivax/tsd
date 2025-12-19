// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"

	"github.com/treivax/tsd/constraint"
)

const (
	FieldTypeUnknown   = "unknown"
	FieldTypePrimitive = "primitive"
	FieldTypeFact      = "fact"
)

// FieldResolver résout les valeurs de champs en tenant compte de leur type
type FieldResolver struct {
	TypeMap map[string]constraint.TypeDefinition
}

// NewFieldResolver crée un nouveau résolveur de champs
func NewFieldResolver(types []constraint.TypeDefinition) *FieldResolver {
	typeMap := make(map[string]constraint.TypeDefinition)
	for _, t := range types {
		typeMap[t.Name] = t
	}

	return &FieldResolver{
		TypeMap: typeMap,
	}
}

// ResolveFieldValue résout la valeur d'un champ d'un fait
// Pour les types primitifs, retourne la valeur directement
// Pour les types de faits, retourne l'ID interne du fait référencé
func (fr *FieldResolver) ResolveFieldValue(fact *Fact, fieldName string) (interface{}, string, error) {
	// Le champ _id_ est interdit
	if fieldName == constraint.FieldNameInternalID {
		return nil, "", fmt.Errorf("le champ '%s' est interne et ne peut pas être accédé", constraint.FieldNameInternalID)
	}

	// Vérifier que le champ existe dans le fait
	value, exists := fact.Fields[fieldName]
	if !exists {
		return nil, "", fmt.Errorf("champ '%s' non trouvé dans le fait de type '%s'", fieldName, fact.Type)
	}

	// Obtenir la définition du type pour connaître le type du champ
	typeDef, exists := fr.TypeMap[fact.Type]
	if !exists {
		return nil, "", fmt.Errorf("type '%s' non trouvé dans le type map", fact.Type)
	}

	// Trouver le champ dans la définition du type
	var fieldDef constraint.Field
	found := false
	for _, f := range typeDef.Fields {
		if f.Name == fieldName {
			fieldDef = f
			found = true
			break
		}
	}

	if !found {
		return nil, "", fmt.Errorf("champ '%s' non défini dans le type '%s'", fieldName, fact.Type)
	}

	// Déterminer le type du champ
	fieldType := fr.getFieldType(fieldDef.Type)

	return value, fieldType, nil
}

// getFieldType retourne le type d'un champ (primitive ou user-defined)
func (fr *FieldResolver) getFieldType(typeName string) string {
	if constraint.IsValidPrimitiveType(typeName) {
		return FieldTypePrimitive
	}

	// Vérifier si c'est un type utilisateur défini
	if _, exists := fr.TypeMap[typeName]; exists {
		return FieldTypeFact
	}

	return FieldTypeUnknown
}

// ResolveFactID résout une variable de fait vers son ID interne
func (fr *FieldResolver) ResolveFactID(fact *Fact) string {
	return fact.ID
}
