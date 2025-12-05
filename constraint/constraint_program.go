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

	// Validation des types
	if err := ValidateTypes(program); err != nil {
		return fmt.Errorf("erreur validation types: %v", err)
	}

	// Validation des faits
	if err := ValidateFacts(program); err != nil {
		return fmt.Errorf("erreur validation faits: %v", err)
	}

	// Validation des contraintes dans les expressions
	if err := validateExpressionConstraints(program); err != nil {
		return err
	}

	// Validation des actions dans les expressions
	if err := validateExpressionActions(program); err != nil {
		return err
	}

	tsdio.Printf("✓ Programme valide avec %d type(s), %d expression(s) et %d fait(s)\n", len(program.Types), len(program.Expressions), len(program.Facts))
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
