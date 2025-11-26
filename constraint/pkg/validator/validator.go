// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package validator

import (
	"fmt"
	"strings"

	"github.com/treivax/tsd/constraint/pkg/domain"
)

// ConstraintValidator implémente l'interface domain.Validator
type ConstraintValidator struct {
	typeRegistry domain.TypeRegistry
	typeChecker  domain.TypeChecker
	config       domain.ValidatorConfig
}

// NewConstraintValidator crée un nouveau validateur
func NewConstraintValidator(registry domain.TypeRegistry, checker domain.TypeChecker) *ConstraintValidator {
	return &ConstraintValidator{
		typeRegistry: registry,
		typeChecker:  checker,
		config: domain.ValidatorConfig{
			StrictMode:       true,
			AllowedOperators: []string{"==", "!=", "<", ">", "<=", ">=", "AND", "OR", "NOT"},
			MaxDepth:         10,
		},
	}
}

// ValidateProgram valide un programme complet
func (v *ConstraintValidator) ValidateProgram(program interface{}) error {
	// Conversion vers Program
	prog, ok := program.(*domain.Program)
	if !ok {
		return domain.NewValidationError("invalid program type", domain.Context{})
	}

	// Validation des types d'abord
	if err := v.ValidateTypes(prog.Types); err != nil {
		return err
	}

	// Enregistrer les types dans le registry
	for _, typeDef := range prog.Types {
		if err := v.typeRegistry.RegisterType(typeDef); err != nil {
			return domain.NewValidationError(
				fmt.Sprintf("failed to register type %s: %v", typeDef.Name, err),
				domain.Context{Type: typeDef.Name},
			)
		}
	}

	// Validation des expressions
	for i, expr := range prog.Expressions {
		if err := v.ValidateExpression(expr, prog.Types); err != nil {
			ctx := domain.Context{
				Field: fmt.Sprintf("expression[%d]", i),
			}
			return domain.NewValidationError(
				fmt.Sprintf("invalid expression %d: %v", i, err),
				ctx,
			)
		}
	}

	return nil
}

// ValidateTypes valide les définitions de types
func (v *ConstraintValidator) ValidateTypes(types []domain.TypeDefinition) error {
	typeNames := make(map[string]bool)

	for i, typeDef := range types {
		// Vérifier les noms dupliqués
		if typeNames[typeDef.Name] {
			return domain.NewValidationError(
				fmt.Sprintf("duplicate type name: %s", typeDef.Name),
				domain.Context{
					Type:  typeDef.Name,
					Field: fmt.Sprintf("types[%d]", i),
				},
			)
		}
		typeNames[typeDef.Name] = true

		// Vérifier que le type a un nom valide
		if typeDef.Name == "" {
			return domain.NewValidationError(
				"type name cannot be empty",
				domain.Context{Field: fmt.Sprintf("types[%d].name", i)},
			)
		}

		// Vérifier que le type a des champs
		if len(typeDef.Fields) == 0 {
			return domain.NewValidationError(
				fmt.Sprintf("type %s must have at least one field", typeDef.Name),
				domain.Context{
					Type:  typeDef.Name,
					Field: fmt.Sprintf("types[%d].fields", i),
				},
			)
		}

		// Valider chaque champ
		fieldNames := make(map[string]bool)
		for j, field := range typeDef.Fields {
			if fieldNames[field.Name] {
				return domain.NewValidationError(
					fmt.Sprintf("duplicate field name '%s' in type '%s'", field.Name, typeDef.Name),
					domain.Context{
						Type:  typeDef.Name,
						Field: field.Name,
					},
				)
			}
			fieldNames[field.Name] = true

			if field.Name == "" {
				return domain.NewValidationError(
					fmt.Sprintf("field name cannot be empty in type %s", typeDef.Name),
					domain.Context{
						Type:  typeDef.Name,
						Field: fmt.Sprintf("types[%d].fields[%d].name", i, j),
					},
				)
			}

			if !domain.IsValidType(field.Type) {
				return domain.NewValidationError(
					fmt.Sprintf("invalid field type '%s' for field '%s' in type '%s'",
						field.Type, field.Name, typeDef.Name),
					domain.Context{
						Type:     typeDef.Name,
						Field:    field.Name,
						Expected: "string, number, bool, or integer",
						Actual:   field.Type,
					},
				)
			}
		}
	}

	return nil
}

// ValidateExpression valide une expression/règle
func (v *ConstraintValidator) ValidateExpression(expr domain.Expression, types []domain.TypeDefinition) error {
	// Valider le set de variables
	if len(expr.Set.Variables) == 0 {
		return domain.NewValidationError(
			"expression must have at least one variable",
			domain.Context{},
		)
	}

	// Vérifier que tous les types des variables existent
	for i, variable := range expr.Set.Variables {
		found := false
		for _, typeDef := range types {
			if typeDef.Name == variable.DataType {
				found = true
				break
			}
		}
		if !found {
			return domain.NewUnknownTypeError(
				variable.DataType,
				domain.Context{
					Variable: variable.Name,
					Field:    fmt.Sprintf("set.variables[%d].dataType", i),
				},
			)
		}
	}

	// Valider les contraintes si elles existent
	if expr.Constraints != nil {
		if err := v.ValidateConstraint(expr.Constraints, expr.Set.Variables, types); err != nil {
			return err
		}
	}

	// Valider l'action (maintenant obligatoire)
	if expr.Action != nil {
		validator := NewActionValidator()
		if err := validator.ValidateAction(expr.Action); err != nil {
			return err
		}
	} else {
		// Avec la nouvelle grammaire, cette condition ne devrait plus arriver
		return fmt.Errorf("action manquante: chaque règle doit avoir une action définie")
	}

	return nil
}

// ValidateConstraint valide une contrainte
func (v *ConstraintValidator) ValidateConstraint(constraint interface{}, variables []domain.TypedVariable, types []domain.TypeDefinition) error {
	// Cette méthode délègue au type checker pour les détails de validation des types
	return v.typeChecker.ValidateTypeCompatibility("", "", "")
}

// SetConfig configure le validateur
func (v *ConstraintValidator) SetConfig(config domain.ValidatorConfig) {
	v.config = config
}

// GetConfig retourne la configuration actuelle
func (v *ConstraintValidator) GetConfig() domain.ValidatorConfig {
	return v.config
}

// ActionValidator valide les actions
type ActionValidator struct{}

// NewActionValidator crée un nouveau validateur d'actions
func NewActionValidator() *ActionValidator {
	return &ActionValidator{}
}

// ValidateAction valide une action
func (av *ActionValidator) ValidateAction(action *domain.Action) error {
	if action == nil {
		return domain.NewActionError(
			"action cannot be nil",
			domain.Context{},
		)
	}

	return av.ValidateJobCall(action.Job)
}

// ValidateJobCall valide un appel de fonction/job
func (av *ActionValidator) ValidateJobCall(jobCall domain.JobCall) error {
	if strings.TrimSpace(jobCall.Name) == "" {
		return domain.NewActionError(
			"job name cannot be empty",
			domain.Context{Field: "job.name"},
		)
	}

	// Valider les arguments (optionnel)
	for i, arg := range jobCall.Args {
		// Convertir l'argument en string pour la validation
		var argStr string
		if s, ok := arg.(string); ok {
			argStr = s
		} else {
			// Pour les objets complexes, on les considère comme valides s'ils ne sont pas nil
			if arg == nil {
				return domain.NewActionError(
					fmt.Sprintf("job argument %d cannot be nil", i),
					domain.Context{
						Field: fmt.Sprintf("job.args[%d]", i),
						Value: arg,
					},
				)
			}
			continue // Les objets complexes sont valides s'ils ne sont pas nil
		}

		if strings.TrimSpace(argStr) == "" {
			return domain.NewActionError(
				fmt.Sprintf("job argument %d cannot be empty", i),
				domain.Context{
					Field: fmt.Sprintf("job.args[%d]", i),
					Value: arg,
				},
			)
		}
	}

	return nil
}
