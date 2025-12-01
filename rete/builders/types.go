// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package builders

import (
	"fmt"

	"github.com/treivax/tsd/rete"
)

// TypeBuilder handles the creation of TypeNodes and type definitions
type TypeBuilder struct {
	utils *BuilderUtils
}

// NewTypeBuilder creates a new TypeBuilder instance
func NewTypeBuilder(utils *BuilderUtils) *TypeBuilder {
	return &TypeBuilder{
		utils: utils,
	}
}

// CreateTypeNodes creates TypeNodes from type definitions
func (tb *TypeBuilder) CreateTypeNodes(network *rete.ReteNetwork, types []interface{}, storage rete.Storage) error {
	for _, typeInterface := range types {
		typeMap, ok := typeInterface.(map[string]interface{})
		if !ok {
			return fmt.Errorf("format type invalide: %T", typeInterface)
		}

		// Extract the type name
		typeName, ok := typeMap["name"].(string)
		if !ok {
			return fmt.Errorf("nom de type non trouvé")
		}

		// Create the type definition
		typeDef := tb.CreateTypeDefinition(typeName, typeMap)

		// Create the TypeNode
		typeNode := rete.NewTypeNode(typeName, typeDef, storage)
		network.TypeNodes[typeName] = typeNode

		// Register the TypeNode in the LifecycleManager
		if network.LifecycleManager != nil {
			network.LifecycleManager.RegisterNode(typeNode.GetID(), "type")
		}

		// CRUCIAL: Connect the TypeNode to the RootNode to enable fact propagation
		network.RootNode.AddChild(typeNode)

		fmt.Printf("   ✓ TypeNode créé: %s\n", typeName)
	}

	return nil
}

// CreateTypeDefinition creates a type definition from a map
func (tb *TypeBuilder) CreateTypeDefinition(typeName string, typeMap map[string]interface{}) rete.TypeDefinition {
	typeDef := rete.TypeDefinition{
		Type:   "type",
		Name:   typeName,
		Fields: []rete.Field{},
	}

	// Extract the fields
	fieldsData, hasFields := typeMap["fields"]
	if !hasFields {
		return typeDef
	}

	fields, ok := fieldsData.([]interface{})
	if !ok {
		return typeDef
	}

	for _, fieldInterface := range fields {
		fieldMap, ok := fieldInterface.(map[string]interface{})
		if !ok {
			continue
		}

		fieldName := GetStringField(fieldMap, "name", "")
		fieldType := GetStringField(fieldMap, "type", "")

		if fieldName != "" && fieldType != "" {
			typeDef.Fields = append(typeDef.Fields, rete.Field{
				Name: fieldName,
				Type: fieldType,
			})
		}
	}

	return typeDef
}
