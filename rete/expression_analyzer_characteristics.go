// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package rete fournit l'implémentation du réseau RETE pour l'évaluation de règles.
// Ce fichier contient les fonctions d'analyse des caractéristiques structurelles
// des expressions (décomposabilité, normalisation, complexité, etc.).
package rete

// CanDecompose détermine si une expression peut être décomposée en chaîne alpha
// Retourne true pour les expressions qui peuvent être transformées en séquence
// de nœuds alpha chaînés.
//
// Les expressions décomposables:
//   - ExprTypeSimple: condition atomique directement utilisable
//   - ExprTypeAND: peut être décomposée en chaîne de conditions (p1 AND p2 AND p3 -> alpha1->alpha2->alpha3)
//   - ExprTypeArithmetic: opérations arithmétiques commutatives peuvent être chaînées
//
// Les expressions non-décomposables:
//   - ExprTypeOR: nécessite des branches multiples (beta nodes ou duplication)
//   - ExprTypeMixed: nécessite normalisation en forme normale disjonctive (DNF) ou conjonctive (CNF)
func CanDecompose(exprType ExpressionType) bool {
	switch exprType {
	case ExprTypeSimple:
		return true
	case ExprTypeAND:
		return true
	case ExprTypeArithmetic:
		return true
	case ExprTypeNOT:
		return true // NOT peut être décomposé en alpha node avec négation
	case ExprTypeOR:
		return false // Nécessite traitement spécial avec branches
	case ExprTypeMixed:
		return false // Nécessite normalisation d'abord
	default:
		return false
	}
}

// ShouldNormalize détermine si une expression nécessite une normalisation
// avant de pouvoir être traitée dans le réseau RETE.
//
// La normalisation est nécessaire pour:
//   - ExprTypeMixed: doit être convertie en DNF ou CNF
//   - ExprTypeOR: peut bénéficier de la normalisation pour optimisation
//
// La normalisation n'est pas nécessaire pour:
//   - ExprTypeSimple: déjà sous forme atomique
//   - ExprTypeAND: déjà sous forme conjonctive, facilement chaînable
//   - ExprTypeArithmetic: structure déjà linéaire
func ShouldNormalize(exprType ExpressionType) bool {
	switch exprType {
	case ExprTypeSimple:
		return false
	case ExprTypeAND:
		return false
	case ExprTypeArithmetic:
		return false
	case ExprTypeNOT:
		return false // NOT peut être géré directement, mais peut bénéficier de la normalisation de l'expression interne
	case ExprTypeOR:
		return true // Bénéficie de la normalisation pour optimisation
	case ExprTypeMixed:
		return true // Doit être normalisée
	default:
		return false
	}
}

// GetExpressionComplexity retourne une estimation de la complexité d'une expression
// Utile pour décider de stratégies d'optimisation
func GetExpressionComplexity(exprType ExpressionType) int {
	switch exprType {
	case ExprTypeSimple:
		return 1
	case ExprTypeAND:
		return 2
	case ExprTypeArithmetic:
		return 2
	case ExprTypeNOT:
		return 2 // La négation ajoute une couche de complexité
	case ExprTypeOR:
		return 3
	case ExprTypeMixed:
		return 4
	default:
		return 0
	}
}

// RequiresBetaNode détermine si une expression nécessite des nœuds beta
// pour son évaluation (jointures, branches multiples, etc.)
func RequiresBetaNode(exprType ExpressionType) bool {
	switch exprType {
	case ExprTypeSimple:
		return false
	case ExprTypeAND:
		return false // Peut être géré avec alpha nodes chaînés
	case ExprTypeArithmetic:
		return false
	case ExprTypeNOT:
		return false // Peut être géré avec un alpha node avec flag de négation
	case ExprTypeOR:
		return true // Nécessite branches ou duplication
	case ExprTypeMixed:
		return true // Nécessite structure complexe
	default:
		return false
	}
}
