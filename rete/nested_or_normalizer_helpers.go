// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

// NormalizeNestedOR normalise une expression avec OR imbriqués
// Combine aplatissement et normalisation canonique
func NormalizeNestedOR(expr interface{}) (interface{}, error) {
	if expr == nil {
		return nil, fmt.Errorf("expression nil")
	}

	// Étape 1: Analyser la structure
	analysis, err := AnalyzeNestedOR(expr)
	if err != nil {
		return nil, fmt.Errorf("erreur analyse: %w", err)
	}

	// Étape 2: Aplatir si nécessaire
	if analysis.RequiresFlattening {
		expr, err = FlattenNestedOR(expr)
		if err != nil {
			return nil, fmt.Errorf("erreur aplatissement: %w", err)
		}
	}

	// Étape 3: Transformer en DNF si nécessaire
	if analysis.RequiresDNF {
		expr, err = TransformToDNF(expr)
		if err != nil {
			return nil, fmt.Errorf("erreur transformation DNF: %w", err)
		}
	}

	// Étape 4: Normalisation canonique finale
	normalized, err := NormalizeORExpression(expr)
	if err != nil {
		return nil, fmt.Errorf("erreur normalisation: %w", err)
	}

	return normalized, nil
}
