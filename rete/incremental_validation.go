// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"

	"github.com/treivax/tsd/constraint"
)

// IncrementalValidator gère la validation sémantique incrémentale avec contexte
type IncrementalValidator struct {
	network *ReteNetwork
}

// NewIncrementalValidator crée un nouveau validateur incrémental
func NewIncrementalValidator(network *ReteNetwork) *IncrementalValidator {
	return &IncrementalValidator{
		network: network,
	}
}

// ValidateWithContext valide un programme en tenant compte du contexte du réseau existant
func (iv *IncrementalValidator) ValidateWithContext(parsedAST interface{}) error {
	if iv.network == nil {
		// Pas de contexte, validation standard
		return constraint.ValidateConstraintProgram(parsedAST)
	}

	// Extraire les types existants du réseau
	// Extract existing types and actions from the network
	existingTypes := iv.extractExistingTypes()
	existingActions := iv.extractExistingActions()

	// Convertir l'AST en programme
	program, err := constraint.ConvertResultToProgram(parsedAST)
	if err != nil {
		return fmt.Errorf("erreur conversion programme: %w", err)
	}

	// Fusionner les types et actions existants avec les nouveaux
	mergedProgram := iv.mergePrograms(existingTypes, existingActions, program)

	// Créer un AST fusionné pour la validation
	mergedAST := iv.programToAST(mergedProgram)

	// Valider le programme complet
	err = constraint.ValidateConstraintProgram(mergedAST)
	if err != nil {
		return fmt.Errorf("validation incrémentale échouée: %w", err)
	}

	// Validation supplémentaire: vérifier la cohérence inter-fichiers
	err = iv.validateCrossFileConsistency(program, existingTypes)
	if err != nil {
		return fmt.Errorf("incohérence inter-fichiers: %w", err)
	}

	return nil
}

// extractExistingTypes extrait les types déjà présents dans le réseau
func (iv *IncrementalValidator) extractExistingTypes() []constraint.TypeDefinition {
	types := make([]constraint.TypeDefinition, 0, len(iv.network.Types))

	for _, typeDef := range iv.network.Types {
		// Convertir TypeDefinition RETE vers TypeDefinition constraint
		fields := make([]constraint.Field, len(typeDef.Fields))
		for i, field := range typeDef.Fields {
			fields[i] = constraint.Field{
				Name:         field.Name,
				Type:         field.Type,
				IsPrimaryKey: field.IsPrimaryKey,
			}
		}

		types = append(types, constraint.TypeDefinition{
			Type:   "typeDefinition",
			Name:   typeDef.Name,
			Fields: fields,
		})
	}

	return types
}

// extractExistingActions extrait les actions déjà présentes dans le réseau
func (iv *IncrementalValidator) extractExistingActions() []constraint.ActionDefinition {
	actions := make([]constraint.ActionDefinition, 0, len(iv.network.Actions))

	for _, actionDef := range iv.network.Actions {
		// Convertir ActionDefinition RETE vers ActionDefinition constraint
		params := make([]constraint.Parameter, len(actionDef.Parameters))
		for i, param := range actionDef.Parameters {
			params[i] = constraint.Parameter{
				Name: param.Name,
				Type: param.Type,
			}
		}

		actions = append(actions, constraint.ActionDefinition{
			Type:       actionDef.Type,
			Name:       actionDef.Name,
			Parameters: params,
		})
	}

	return actions
}

// mergePrograms fusionne les types et actions existants avec le nouveau programme
func (iv *IncrementalValidator) mergePrograms(
	existingTypes []constraint.TypeDefinition,
	existingActions []constraint.ActionDefinition,
	newProgram *constraint.Program,
) *constraint.Program {
	merged := &constraint.Program{
		Types:       make([]constraint.TypeDefinition, 0),
		Actions:     make([]constraint.ActionDefinition, 0),
		Expressions: newProgram.Expressions,
		Facts:       newProgram.Facts,
	}

	// Créer un index des types existants
	typeIndex := make(map[string]constraint.TypeDefinition)
	for _, typeDef := range existingTypes {
		typeIndex[typeDef.Name] = typeDef
	}

	// Ajouter les types existants qui ne sont pas redéfinis
	for _, typeDef := range existingTypes {
		// Vérifier si ce type est redéfini dans le nouveau programme
		redefined := false
		for _, newType := range newProgram.Types {
			if newType.Name == typeDef.Name {
				redefined = true
				break
			}
		}

		if !redefined {
			merged.Types = append(merged.Types, typeDef)
		}
	}

	// Ajouter les nouveaux types (écrasent les anciens si même nom)
	merged.Types = append(merged.Types, newProgram.Types...)

	// Créer un index des actions existantes
	actionIndex := make(map[string]constraint.ActionDefinition)
	for _, actionDef := range existingActions {
		actionIndex[actionDef.Name] = actionDef
	}

	// Ajouter les actions existantes qui ne sont pas redéfinies
	for _, actionDef := range existingActions {
		// Vérifier si cette action est redéfinie dans le nouveau programme
		redefined := false
		for _, newAction := range newProgram.Actions {
			if newAction.Name == actionDef.Name {
				redefined = true
				break
			}
		}

		if !redefined {
			merged.Actions = append(merged.Actions, actionDef)
		}
	}

	// Ajouter les nouvelles actions (écrasent les anciennes si même nom)
	merged.Actions = append(merged.Actions, newProgram.Actions...)

	return merged
}

// programToAST convertit un programme en AST pour la validation
func (iv *IncrementalValidator) programToAST(program *constraint.Program) interface{} {
	result := make(map[string]interface{})

	// Convertir les types
	types := make([]interface{}, len(program.Types))
	for i, typeDef := range program.Types {
		fields := make([]interface{}, len(typeDef.Fields))
		for j, field := range typeDef.Fields {
			fields[j] = map[string]interface{}{
				"name": field.Name,
				"type": field.Type,
			}
		}

		types[i] = map[string]interface{}{
			"type":   "typeDefinition",
			"name":   typeDef.Name,
			"fields": fields,
		}
	}
	result["types"] = types

	// Convertir les actions
	actions := make([]interface{}, len(program.Actions))
	for i, actionDef := range program.Actions {
		params := make([]interface{}, len(actionDef.Parameters))
		for j, param := range actionDef.Parameters {
			params[j] = map[string]interface{}{
				"name": param.Name,
				"type": param.Type,
			}
		}

		actions[i] = map[string]interface{}{
			"type":       actionDef.Type,
			"name":       actionDef.Name,
			"parameters": params,
		}
	}
	result["actions"] = actions

	// Convertir les expressions
	expressions := make([]interface{}, len(program.Expressions))
	for i, expr := range program.Expressions {
		// Conversion complète pour la validation
		exprMap := map[string]interface{}{
			"type": expr.Type,
			"set":  expr.Set,
		}

		// Ajouter les contraintes si présentes
		if expr.Constraints != nil {
			exprMap["constraints"] = expr.Constraints
		}

		// Ajouter l'action si présente (CRUCIAL pour la validation)
		if expr.Action != nil {
			exprMap["action"] = expr.Action
		}

		// Ajouter les patterns si présents (pour agrégation)
		if expr.Patterns != nil {
			exprMap["patterns"] = expr.Patterns
		}

		// Ajouter le ruleId si présent
		if expr.RuleId != "" {
			exprMap["ruleId"] = expr.RuleId
		}

		expressions[i] = exprMap
	}
	result["expressions"] = expressions

	// Ajouter les faits s'ils existent
	if len(program.Facts) > 0 {
		facts := make([]interface{}, len(program.Facts))
		for i, fact := range program.Facts {
			facts[i] = map[string]interface{}{
				"type":     fact.Type,
				"typeName": fact.TypeName,
				"fields":   fact.Fields,
			}
		}
		result["facts"] = facts
	}

	return result
}

// validateCrossFileConsistency vérifie la cohérence entre fichiers
func (iv *IncrementalValidator) validateCrossFileConsistency(
	newProgram *constraint.Program,
	existingTypes []constraint.TypeDefinition,
) error {
	// Créer un index des types existants
	typeIndex := make(map[string]constraint.TypeDefinition)
	for _, typeDef := range existingTypes {
		typeIndex[typeDef.Name] = typeDef
	}

	// Vérifier que les types référencés dans les nouvelles expressions existent
	for _, expr := range newProgram.Expressions {
		// Extraire les types utilisés dans l'expression
		usedTypes := iv.extractTypesFromExpression(expr)

		for _, typeName := range usedTypes {
			// Vérifier si le type existe dans les types existants ou nouveaux
			existsInExisting := false
			existsInNew := false

			if _, exists := typeIndex[typeName]; exists {
				existsInExisting = true
			}

			for _, newType := range newProgram.Types {
				if newType.Name == typeName {
					existsInNew = true
					break
				}
			}

			if !existsInExisting && !existsInNew {
				return fmt.Errorf("type '%s' référencé mais non défini", typeName)
			}
		}
	}

	// Vérifier que les champs référencés existent dans les types
	err := iv.validateFieldReferences(newProgram, typeIndex)
	if err != nil {
		return err
	}

	// Vérifier que les types redéfinis sont compatibles
	err = iv.validateTypeCompatibility(newProgram, existingTypes)
	if err != nil {
		return err
	}

	return nil
}

// extractTypesFromExpression extrait les noms de types utilisés dans une expression
func (iv *IncrementalValidator) extractTypesFromExpression(expr constraint.Expression) []string {
	types := make([]string, 0)

	// Extraire les types des variables
	for _, variable := range expr.Set.Variables {
		types = append(types, variable.DataType)
	}

	return types
}

// validateFieldReferences valide que tous les champs référencés existent
func (iv *IncrementalValidator) validateFieldReferences(
	program *constraint.Program,
	existingTypeIndex map[string]constraint.TypeDefinition,
) error {
	// Créer un index combiné de tous les types (existants + nouveaux)
	allTypes := make(map[string]constraint.TypeDefinition)

	// Ajouter les types existants
	for name, typeDef := range existingTypeIndex {
		allTypes[name] = typeDef
	}

	// Ajouter les nouveaux types (écrasent les anciens si même nom)
	for _, typeDef := range program.Types {
		allTypes[typeDef.Name] = typeDef
	}

	// Pour chaque fait, vérifier que les champs existent dans le type
	for _, fact := range program.Facts {
		typeDef, exists := allTypes[fact.TypeName]
		if !exists {
			return fmt.Errorf("fait de type '%s' mais type non défini", fact.TypeName)
		}

		// Créer un index des champs du type
		fieldIndex := make(map[string]bool)
		for _, field := range typeDef.Fields {
			fieldIndex[field.Name] = true
		}

		// Vérifier que tous les champs du fait existent dans le type
		for _, factField := range fact.Fields {
			if !fieldIndex[factField.Name] && factField.Name != "id" {
				return fmt.Errorf(
					"champ '%s' n'existe pas dans le type '%s'",
					factField.Name,
					fact.TypeName,
				)
			}
		}
	}

	return nil
}

// validateTypeCompatibility vérifie que les types redéfinis sont compatibles
func (iv *IncrementalValidator) validateTypeCompatibility(
	newProgram *constraint.Program,
	existingTypes []constraint.TypeDefinition,
) error {
	// Créer un index des types existants
	existingTypeIndex := make(map[string]constraint.TypeDefinition)
	for _, typeDef := range existingTypes {
		existingTypeIndex[typeDef.Name] = typeDef
	}

	// Vérifier chaque nouveau type
	for _, newType := range newProgram.Types {
		existingType, exists := existingTypeIndex[newType.Name]
		if !exists {
			// Nouveau type, pas de problème de compatibilité
			continue
		}

		// Type redéfini, vérifier la compatibilité
		// Pour l'instant, on autorise uniquement les redéfinitions identiques
		if !iv.areTypesEqual(newType, existingType) {
			return fmt.Errorf(
				"type '%s' redéfini de manière incompatible (modification de types existants non supportée)",
				newType.Name,
			)
		}
	}

	return nil
}

// areTypesEqual vérifie si deux définitions de type sont identiques
func (iv *IncrementalValidator) areTypesEqual(
	type1, type2 constraint.TypeDefinition,
) bool {
	if type1.Name != type2.Name {
		return false
	}

	if len(type1.Fields) != len(type2.Fields) {
		return false
	}

	// Créer un index des champs du type2
	type2Fields := make(map[string]string)
	for _, field := range type2.Fields {
		type2Fields[field.Name] = field.Type
	}

	// Vérifier que tous les champs de type1 existent dans type2 avec le même type
	for _, field := range type1.Fields {
		fieldType, exists := type2Fields[field.Name]
		if !exists || fieldType != field.Type {
			return false
		}
	}

	return true
}
