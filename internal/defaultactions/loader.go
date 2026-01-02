// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package defaultactions

import (
	"fmt"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/internal/actiondefs"
)

// LoadDefaultActions parse le fichier defaults.tsd et retourne les actions.
//
// Cette fonction charge les définitions des actions système embarquées dans le binaire,
// les parse via le parser TSD standard, et marque chaque action comme "par défaut"
// pour empêcher leur redéfinition par l'utilisateur.
//
// Retourne:
//   - []constraint.ActionDefinition: les 6 actions par défaut avec IsDefault=true
//   - error: en cas d'erreur de parsing ou si le nombre d'actions ne correspond pas
func LoadDefaultActions() ([]constraint.ActionDefinition, error) {
	// Parser le fichier embarqué depuis actiondefs (package intermédiaire)
	result, err := constraint.ParseConstraint("defaults.tsd", []byte(actiondefs.DefaultActionsTSD))
	if err != nil {
		return nil, fmt.Errorf("failed to parse default actions: %w", err)
	}

	// Convertir le résultat en Program
	program, err := constraint.ConvertResultToProgram(result)
	if err != nil {
		return nil, fmt.Errorf("failed to convert default actions to program: %w", err)
	}

	// Vérifier que toutes les actions attendues sont présentes
	if len(program.Actions) != len(actiondefs.DefaultActionNames) {
		return nil, fmt.Errorf("expected %d default actions, got %d",
			len(actiondefs.DefaultActionNames), len(program.Actions))
	}

	// Marquer chaque action comme "par défaut"
	for i := range program.Actions {
		program.Actions[i].IsDefault = true
	}

	return program.Actions, nil
}
