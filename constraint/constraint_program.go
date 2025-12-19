// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"encoding/json"
	"fmt"

	"github.com/treivax/tsd/tsdio"
)

// ValidateProgram effectue une validation complète du programme parsé
func ValidateProgram(result interface{}) error {
	// Convertir le résultat en structure Program
	program, err := convertResultToProgram(result)
	if err != nil {
		return err
	}

	// Normaliser les types de valeurs de faits (avant validation)
	normalizeFactValueTypes(&program)

	// Validation des types
	if err := ValidateTypes(program); err != nil {
		return fmt.Errorf("erreur validation types: %v", err)
	}

	// Validation des références de types utilisateur dans les champs
	if err := validateTypeReferences(program); err != nil {
		return fmt.Errorf("erreur validation références types: %v", err)
	}

	// Validation des références circulaires entre types
	if err := validateNoCircularReferences(program); err != nil {
		return fmt.Errorf("erreur validation références circulaires: %v", err)
	}

	// Validation des variables dans les faits
	if err := validateVariableReferences(program); err != nil {
		return fmt.Errorf("erreur validation références variables: %v", err)
	}

	// Validation des faits
	if err := ValidateFacts(program); err != nil {
		return fmt.Errorf("erreur validation faits: %v", err)
	}

	// Validation des xuple-spaces
	if err := validateXupleSpaces(program); err != nil {
		return fmt.Errorf("erreur validation xuple-spaces: %v", err)
	}

	// Validation des contraintes dans les expressions
	if err := validateExpressionConstraints(program); err != nil {
		return err
	}

	// Validation des actions dans les expressions
	if err := validateExpressionActions(program); err != nil {
		return err
	}

	tsdio.Printf("✓ Programme valide avec %d type(s), %d expression(s), %d fait(s), %d affectation(s) et %d xuple-space(s)\n",
		len(program.Types), len(program.Expressions), len(program.Facts), len(program.FactAssignments), len(program.XupleSpaces))
	return nil
}

// convertResultToProgram converts the parser result to a Program structure
func convertResultToProgram(result interface{}) (Program, error) {
	var program Program

	jsonData, err := json.Marshal(result)
	if err != nil {
		return program, fmt.Errorf("erreur conversion JSON: %v", err)
	}

	err = json.Unmarshal(jsonData, &program)
	if err != nil {
		return program, fmt.Errorf("erreur parsing JSON: %v", err)
	}

	return program, nil
}

// validateExpressionConstraints validates field access and type compatibility in all constraints
func validateExpressionConstraints(program Program) error {
	for i, expression := range program.Expressions {
		if expression.Constraints != nil {
			// Validation des accès aux champs
			if err := ValidateConstraintFieldAccess(program, expression.Constraints, i); err != nil {
				return fmt.Errorf("erreur validation champs dans l'expression %d: %v", i+1, err)
			}

			// Validation des types dans les comparaisons
			if err := ValidateTypeCompatibility(program, expression.Constraints, i); err != nil {
				return fmt.Errorf("erreur validation types dans l'expression %d: %v", i+1, err)
			}
		}
	}
	return nil
}

// validateExpressionActions validates that all expressions have valid actions
func validateExpressionActions(program Program) error {
	for i, expression := range program.Expressions {
		if expression.Action != nil {
			if err := ValidateAction(program, *expression.Action, i); err != nil {
				return fmt.Errorf("erreur validation action dans l'expression %d: %v", i+1, err)
			}
		} else {
			// Avec la nouvelle grammaire, cette condition ne devrait plus arriver
			return fmt.Errorf("action manquante dans l'expression %d: chaque règle doit avoir une action définie", i+1)
		}
	}
	return nil
}

// validateXupleSpaces valide les déclarations de xuple-spaces
func validateXupleSpaces(program Program) error {
	for _, xs := range program.XupleSpaces {
		if err := ValidateXupleSpaceDeclaration(&xs); err != nil {
			return err
		}
	}
	return nil
}

// validateTypeReferences validates that all user-defined types referenced in fields exist.
func validateTypeReferences(program Program) error {
	typeMap := make(map[string]bool)
	for _, typeDef := range program.Types {
		typeMap[typeDef.Name] = true
	}

	primitiveTypes := GetPrimitiveTypesSet()

	for _, typeDef := range program.Types {
		for _, field := range typeDef.Fields {
			if !primitiveTypes[field.Type] && !typeMap[field.Type] {
				return fmt.Errorf(
					"type '%s': champ '%s' référence un type inconnu '%s'",
					typeDef.Name,
					field.Name,
					field.Type,
				)
			}
		}
	}

	return nil
}

// validateVariableReferences validates that all variable references in facts are defined.
// Only validates fields that are supposed to be references (custom types, not primitives).
func validateVariableReferences(program Program) error {
	varMap := buildVariableMap(program)
	typeDefMap := buildTypeDefinitionMap(program)
	primitiveTypes := GetPrimitiveTypesSet()

	for i, fact := range program.Facts {
		if err := validateFactVariableReferences(fact, i, varMap, typeDefMap, primitiveTypes); err != nil {
			return err
		}
	}

	return nil
}

// buildVariableMap builds a map of variable names to their fact types
func buildVariableMap(program Program) map[string]string {
	varMap := make(map[string]string)
	for _, assignment := range program.FactAssignments {
		varMap[assignment.Variable] = assignment.Fact.TypeName
	}
	return varMap
}

// buildTypeDefinitionMap builds a map of type names to their definitions
func buildTypeDefinitionMap(program Program) map[string]TypeDefinition {
	typeDefMap := make(map[string]TypeDefinition)
	for _, typeDef := range program.Types {
		typeDefMap[typeDef.Name] = typeDef
	}
	return typeDefMap
}

// buildFieldTypeMap builds a map of field names to their types for a given type definition
func buildFieldTypeMap(typeDef TypeDefinition) map[string]string {
	fieldTypeMap := make(map[string]string)
	for _, fieldDef := range typeDef.Fields {
		fieldTypeMap[fieldDef.Name] = fieldDef.Type
	}
	return fieldTypeMap
}

// validateFactVariableReferences validates variable references in a single fact
func validateFactVariableReferences(
	fact Fact,
	factIndex int,
	varMap map[string]string,
	typeDefMap map[string]TypeDefinition,
	primitiveTypes map[string]bool,
) error {
	typeDef, exists := typeDefMap[fact.TypeName]
	if !exists {
		return nil // Type validation will catch this
	}

	fieldTypeMap := buildFieldTypeMap(typeDef)

	for j, field := range fact.Fields {
		if err := validateFieldVariableReference(
			field, j, factIndex, fact, fieldTypeMap, varMap, primitiveTypes,
		); err != nil {
			return err
		}
	}

	return nil
}

// validateFieldVariableReference validates a single field's variable reference
func validateFieldVariableReference(
	field FactField,
	fieldIndex int,
	factIndex int,
	fact Fact,
	fieldTypeMap map[string]string,
	varMap map[string]string,
	primitiveTypes map[string]bool,
) error {
	if field.Value.Type != ValueTypeVariableReference {
		return nil
	}

	// Check if the field is supposed to be a custom type
	fieldType, exists := fieldTypeMap[field.Name]
	if !exists {
		return nil // Field validation will catch this
	}

	// If it's a primitive type, this shouldn't be a variableReference
	// The parser might have incorrectly set this
	if primitiveTypes[fieldType] {
		return nil // Skip - this is an identifier for a primitive field
	}

	// OK, this is legitimately a variable reference for a custom type field
	varName, ok := field.Value.Value.(string)
	if !ok || varName == "" {
		return fmt.Errorf("fait %d, champ %d: référence de variable invalide", factIndex+1, fieldIndex+1)
	}

	if _, exists := varMap[varName]; !exists {
		return fmt.Errorf(
			"fait %d (%s), champ %d (%s): variable '%s' non définie",
			factIndex+1,
			fact.TypeName,
			fieldIndex+1,
			field.Name,
			varName,
		)
	}

	return nil
}

// validateNoCircularReferences detects circular type dependencies.
func validateNoCircularReferences(program Program) error {
	typeGraph := make(map[string][]string)
	primitiveTypes := GetPrimitiveTypesSet()

	for _, typeDef := range program.Types {
		for _, field := range typeDef.Fields {
			if !primitiveTypes[field.Type] {
				typeGraph[typeDef.Name] = append(typeGraph[typeDef.Name], field.Type)
			}
		}
	}

	for typeName := range typeGraph {
		visited := make(map[string]bool)
		recStack := make(map[string]bool)
		if hasCycle(typeName, typeGraph, visited, recStack) {
			return fmt.Errorf("référence circulaire détectée impliquant le type '%s'", typeName)
		}
	}

	return nil
}

// hasCycle performs depth-first search to detect cycles in the type dependency graph.
func hasCycle(node string, graph map[string][]string, visited, recStack map[string]bool) bool {
	visited[node] = true
	recStack[node] = true

	for _, neighbor := range graph[node] {
		if !visited[neighbor] {
			if hasCycle(neighbor, graph, visited, recStack) {
				return true
			}
		} else if recStack[neighbor] {
			return true
		}
	}

	recStack[node] = false
	return false
}

// normalizeFactValueTypes normalizes fact value types based on field type definitions.
// Converts "variableReference" to "identifier" for primitive-typed fields since
// the parser can't distinguish between them at parse time.
func normalizeFactValueTypes(program *Program) {
	typeDefMap := buildTypeDefinitionMap(*program)
	primitiveTypes := GetPrimitiveTypesSet()

	// Normalize fact assignments
	for i := range program.FactAssignments {
		fact := &program.FactAssignments[i].Fact
		if typeDef, exists := typeDefMap[fact.TypeName]; exists {
			normalizeFactFields(fact, typeDef, primitiveTypes)
		}
	}

	// Normalize regular facts
	for i := range program.Facts {
		fact := &program.Facts[i]
		if typeDef, exists := typeDefMap[fact.TypeName]; exists {
			normalizeFactFields(fact, typeDef, primitiveTypes)
		}
	}
}

// normalizeFactFields normalizes the field value types for a single fact.
func normalizeFactFields(fact *Fact, typeDef TypeDefinition, primitiveTypes map[string]bool) {
	fieldTypeMap := buildFieldTypeMap(typeDef)

	// Normalize each field
	for j := range fact.Fields {
		field := &fact.Fields[j]

		// If the value type is "variableReference" but the field type is primitive,
		// change it to "identifier" (which is what it really is)
		if field.Value.Type == ValueTypeVariableReference {
			if fieldType, exists := fieldTypeMap[field.Name]; exists && primitiveTypes[fieldType] {
				field.Value.Type = ValueTypeIdentifier
			}
		}
	}
}
