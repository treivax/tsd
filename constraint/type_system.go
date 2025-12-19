// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
	"strings"
)

// TypeSystem gère le système de types du langage TSD.
// Il maintient les définitions de types et fournit des utilitaires de validation.
type TypeSystem struct {
	types     map[string]TypeDefinition
	variables map[string]string
}

// NewTypeSystem crée un nouveau système de types.
func NewTypeSystem(types []TypeDefinition) *TypeSystem {
	typeMap := make(map[string]TypeDefinition, len(types))
	for _, t := range types {
		typeMap[t.Name] = t
	}

	return &TypeSystem{
		types:     typeMap,
		variables: make(map[string]string),
	}
}

// IsPrimitiveType vérifie si un type est primitif.
func (ts *TypeSystem) IsPrimitiveType(typeName string) bool {
	return IsValidPrimitiveType(typeName)
}

// IsUserDefinedType vérifie si un type est défini par l'utilisateur.
func (ts *TypeSystem) IsUserDefinedType(typeName string) bool {
	_, exists := ts.types[typeName]
	return exists
}

// TypeExists vérifie qu'un type existe (primitif ou user-defined).
func (ts *TypeSystem) TypeExists(typeName string) bool {
	return ts.IsPrimitiveType(typeName) || ts.IsUserDefinedType(typeName)
}

// GetFieldType retourne le type d'un champ dans un type donné.
func (ts *TypeSystem) GetFieldType(typeName, fieldName string) (string, error) {
	if fieldName == FieldNameInternalID {
		return "", fmt.Errorf(
			"le champ '%s' est interne et ne peut pas être accédé",
			FieldNameInternalID,
		)
	}

	typeDef, exists := ts.types[typeName]
	if !exists {
		return "", fmt.Errorf("type '%s' non trouvé", typeName)
	}

	for _, field := range typeDef.Fields {
		if field.Name == fieldName {
			return field.Type, nil
		}
	}

	return "", fmt.Errorf(
		"champ '%s' non trouvé dans le type '%s'",
		fieldName,
		typeName,
	)
}

// ValidateFieldType valide qu'un champ existe et retourne son type.
func (ts *TypeSystem) ValidateFieldType(typeName, fieldName string) (string, error) {
	fieldType, err := ts.GetFieldType(typeName, fieldName)
	if err != nil {
		return "", err
	}

	if !ts.TypeExists(fieldType) {
		return "", fmt.Errorf(
			"type '%s' du champ '%s.%s' n'existe pas",
			fieldType,
			typeName,
			fieldName,
		)
	}

	return fieldType, nil
}

// RegisterVariable enregistre une variable avec son type.
func (ts *TypeSystem) RegisterVariable(varName, typeName string) error {
	if !ts.TypeExists(typeName) {
		return fmt.Errorf(
			"impossible d'enregistrer la variable '%s': type '%s' n'existe pas",
			varName,
			typeName,
		)
	}

	ts.variables[varName] = typeName
	return nil
}

// GetVariableType retourne le type d'une variable.
func (ts *TypeSystem) GetVariableType(varName string) (string, error) {
	typeName, exists := ts.variables[varName]
	if !exists {
		return "", fmt.Errorf("variable '%s' non définie", varName)
	}
	return typeName, nil
}

// VariableExists vérifie qu'une variable existe.
func (ts *TypeSystem) VariableExists(varName string) bool {
	_, exists := ts.variables[varName]
	return exists
}

// AreTypesCompatible vérifie si deux types sont compatibles pour une opération.
func (ts *TypeSystem) AreTypesCompatible(type1, type2 string, operator string) bool {
	if type1 == type2 {
		if ts.IsUserDefinedType(type1) {
			return operator == OpEq || operator == OpNeq
		}
		return true
	}

	if (type1 == ValueTypeBool || type1 == ValueTypeBoolean) &&
		(type2 == ValueTypeBool || type2 == ValueTypeBoolean) {
		return true
	}

	return false
}

// ValidateCircularReferences détecte les références circulaires dans les types.
func (ts *TypeSystem) ValidateCircularReferences() error {
	graph := make(map[string][]string)

	for typeName, typeDef := range ts.types {
		for _, field := range typeDef.Fields {
			if ts.IsUserDefinedType(field.Type) {
				graph[typeName] = append(graph[typeName], field.Type)
			}
		}
	}

	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var hasCycle func(string) bool
	hasCycle = func(node string) bool {
		visited[node] = true
		recStack[node] = true

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				if hasCycle(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				return true
			}
		}

		recStack[node] = false
		return false
	}

	for typeName := range ts.types {
		if !visited[typeName] {
			if hasCycle(typeName) {
				return fmt.Errorf(
					"référence circulaire détectée impliquant le type '%s'",
					typeName,
				)
			}
		}
	}

	return nil
}

// GetTypePath retourne le chemin de types pour une expression de field access.
// Ex: login.user.name -> [Login, User, string]
func (ts *TypeSystem) GetTypePath(rootType, fieldPath string) ([]string, error) {
	parts := strings.Split(fieldPath, ".")
	path := []string{rootType}
	currentType := rootType

	for _, fieldName := range parts {
		fieldType, err := ts.GetFieldType(currentType, fieldName)
		if err != nil {
			return nil, err
		}

		path = append(path, fieldType)
		currentType = fieldType
	}

	return path, nil
}

// ValidateTypeDefinition valide la définition d'un type.
func ValidateTypeDefinition(typeDef TypeDefinition) error {
	if typeDef.Name == "" {
		return fmt.Errorf("le nom du type ne peut pas être vide")
	}

	if len(typeDef.Fields) == 0 {
		return fmt.Errorf("type '%s': doit avoir au moins un champ", typeDef.Name)
	}

	fieldNames := make(map[string]bool)
	for i, field := range typeDef.Fields {
		if field.Name == "" {
			return fmt.Errorf("type '%s': le champ %d a un nom vide", typeDef.Name, i+1)
		}

		if field.Name == FieldNameInternalID {
			return fmt.Errorf(
				"type '%s': le champ '%s' est réservé au système et ne peut pas être utilisé",
				typeDef.Name,
				FieldNameInternalID,
			)
		}

		if fieldNames[field.Name] {
			return fmt.Errorf("type '%s': champ '%s' défini plusieurs fois", typeDef.Name, field.Name)
		}
		fieldNames[field.Name] = true

		if field.Type == "" {
			return fmt.Errorf("type '%s': le champ '%s' n'a pas de type", typeDef.Name, field.Name)
		}
	}

	return nil
}
