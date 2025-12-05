// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package rete fournit l'implémentation du réseau RETE pour l'évaluation de règles.
// Ce fichier contient les fonctions de génération de hints d'optimisation pour les expressions,
// permettant de suggérer des transformations et améliorations possibles.
package rete

import (
	"github.com/treivax/tsd/constraint"
)

// generateOptimizationHints génère des suggestions d'optimisation basées sur l'analyse
func generateOptimizationHints(expr interface{}, info *ExpressionInfo) []string {
	hints := make([]string, 0)

	// Hint pour De Morgan transformation
	if info.Type == ExprTypeNOT && info.InnerInfo != nil {
		if info.InnerInfo.Type == ExprTypeOR {
			hints = append(hints, "apply_demorgan_not_or")
		} else if info.InnerInfo.Type == ExprTypeAND {
			hints = append(hints, "apply_demorgan_not_and")
		} else if info.InnerInfo.Type == ExprTypeMixed {
			hints = append(hints, "push_negation_down")
		}
	}

	// Hint pour normalisation
	if info.ShouldNormalize {
		if info.Type == ExprTypeMixed {
			hints = append(hints, "normalize_to_dnf")
		} else if info.Type == ExprTypeOR {
			hints = append(hints, "consider_dnf_expansion")
		}
	}

	// Hint pour partage d'alpha nodes
	if info.Type == ExprTypeAND && info.Complexity >= 3 {
		hints = append(hints, "alpha_sharing_opportunity")
	}

	// Hint pour réordonnancement
	if info.Type == ExprTypeAND && canBenefitFromReordering(expr) {
		hints = append(hints, "consider_reordering")
	}

	// Hint pour expressions complexes
	if info.Complexity >= 4 {
		hints = append(hints, "high_complexity_review")
	}

	// Hint pour expressions nécessitant beta nodes
	if info.RequiresBeta {
		hints = append(hints, "requires_beta_node")
	}

	// Hint pour simplification arithmétique
	if info.Type == ExprTypeArithmetic {
		hints = append(hints, "consider_arithmetic_simplification")
	}

	return hints
}

// canBenefitFromReordering détermine si une expression AND peut bénéficier d'un réordonnancement
func canBenefitFromReordering(expr interface{}) bool {
	switch e := expr.(type) {
	case constraint.LogicalExpression:
		// Si l'expression a plusieurs opérations AND, le réordonnancement peut aider
		return len(e.Operations) >= 2

	case map[string]interface{}:
		if operations, ok := e["operations"].([]interface{}); ok {
			return len(operations) >= 2
		}
	}

	return false
}

// ShouldApplyDeMorgan détermine si la transformation de De Morgan devrait être appliquée
// basé sur des critères d'optimisation
func ShouldApplyDeMorgan(expr interface{}) bool {
	info, err := GetExpressionInfo(expr)
	if err != nil {
		return false
	}

	// Appliquer De Morgan si c'est un NOT(OR) ou NOT(AND)
	if info.Type == ExprTypeNOT && info.InnerInfo != nil {
		innerType := info.InnerInfo.Type
		if innerType == ExprTypeOR || innerType == ExprTypeAND {
			// Appliquer De Morgan seulement si cela réduit la complexité
			// NOT(A OR B) -> (NOT A) AND (NOT B) est bénéfique car AND est décomposable
			if innerType == ExprTypeOR {
				return true
			}
			// NOT(A AND B) -> (NOT A) OR (NOT B) peut être moins optimal
			// car OR nécessite des branches, donc on applique seulement si l'expression est simple
			if innerType == ExprTypeAND && info.InnerInfo.Complexity <= 2 {
				return true
			}
		}
	}

	return false
}
