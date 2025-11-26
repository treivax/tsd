// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"

	"github.com/treivax/tsd/constraint"
)

// ASTConverter convertit l'AST du parser constraint vers les types RETE
type ASTConverter struct{}

// NewASTConverter crée un nouveau convertisseur AST
func NewASTConverter() *ASTConverter {
	return &ASTConverter{}
}

// ConvertProgram convertit un constraint.Program vers un rete.Program
func (ac *ASTConverter) ConvertProgram(constraintProgram interface{}) (*Program, error) {
	// Essayer de caster vers constraint.Program
	program, ok := constraintProgram.(*constraint.Program)
	if !ok {
		return nil, fmt.Errorf("type de programme AST non reconnu")
	}

	reteProgram := &Program{
		Types:       make([]TypeDefinition, len(program.Types)),
		Expressions: make([]Expression, len(program.Expressions)),
	}

	// Convertir les types
	for i, constraintType := range program.Types {
		reteProgram.Types[i] = TypeDefinition{
			Type:   constraintType.Type,
			Name:   constraintType.Name,
			Fields: ac.convertFields(constraintType.Fields),
		}
	}

	// Convertir les expressions
	for i, constraintExpr := range program.Expressions {
		reteExpr, err := ac.convertExpression(constraintExpr)
		if err != nil {
			return nil, fmt.Errorf("erreur conversion expression %d: %w", i, err)
		}
		reteProgram.Expressions[i] = *reteExpr
	}

	return reteProgram, nil
}

// convertFields convertit les champs
func (ac *ASTConverter) convertFields(constraintFields []constraint.Field) []Field {
	fields := make([]Field, len(constraintFields))
	for i, field := range constraintFields {
		fields[i] = Field{
			Name: field.Name,
			Type: field.Type,
		}
	}
	return fields
}

// convertExpression convertit une expression
func (ac *ASTConverter) convertExpression(constraintExpr constraint.Expression) (*Expression, error) {
	expr := &Expression{
		Type:        constraintExpr.Type,
		Constraints: constraintExpr.Constraints,
	}

	// Convertir le set
	expr.Set = Set{
		Type:      constraintExpr.Set.Type,
		Variables: ac.convertTypedVariables(constraintExpr.Set.Variables),
	}

	// Convertir l'action (maintenant obligatoire)
	if constraintExpr.Action != nil {
		action, err := ac.convertAction(*constraintExpr.Action)
		if err != nil {
			return nil, fmt.Errorf("erreur conversion action: %w", err)
		}
		expr.Action = action
	} else {
		// Cette condition ne devrait plus arriver avec la nouvelle grammaire
		return nil, fmt.Errorf("action manquante: chaque règle doit avoir une action définie")
	}

	return expr, nil
}

// convertTypedVariables convertit les variables typées
func (ac *ASTConverter) convertTypedVariables(constraintVars []constraint.TypedVariable) []TypedVariable {
	vars := make([]TypedVariable, len(constraintVars))
	for i, variable := range constraintVars {
		vars[i] = TypedVariable{
			Type:     variable.Type,
			Name:     variable.Name,
			DataType: variable.DataType,
		}
	}
	return vars
}

// convertAction convertit une action
func (ac *ASTConverter) convertAction(constraintAction constraint.Action) (*Action, error) {
	action := &Action{
		Type: constraintAction.Type,
		Job: JobCall{
			Type: constraintAction.Job.Type,
			Name: constraintAction.Job.Name,
			Args: constraintAction.Job.Args,
		},
	}
	return action, nil
}
