// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import "sort"

// Ce fichier contient les utilitaires de comparaison, déduplication et normalisation
// de conditions pour les chaînes alpha du système RETE.

// CompareConditions compare deux conditions pour l'égalité
// La comparaison se fait via les hash calculés pour garantir une égalité sémantique
func CompareConditions(c1, c2 SimpleCondition) bool {
	return c1.Hash == c2.Hash
}

// DeduplicateConditions supprime les conditions dupliquées d'une liste
// Utilise les hash des conditions pour détecter les doublons
func DeduplicateConditions(conditions []SimpleCondition) []SimpleCondition {
	seen := make(map[string]bool)
	result := []SimpleCondition{}

	for _, cond := range conditions {
		if !seen[cond.Hash] {
			seen[cond.Hash] = true
			result = append(result, cond)
		}
	}

	return result
}

// IsCommutative retourne true si l'opérateur est commutatif
// Les opérateurs commutatifs (AND, OR, +, *, ==, !=) peuvent être réordonnés
// Les opérateurs non-commutatifs (-, /, <, >, <=, >=, séquences) doivent préserver l'ordre
func IsCommutative(operator string) bool {
	commutativeOps := map[string]bool{
		"AND": true,
		"OR":  true,
		"&&":  true,
		"||":  true,
		"+":   true,
		"*":   true,
		"==":  true,
		"!=":  true,
		"<>":  true,
	}
	return commutativeOps[operator]
}

// NormalizeConditions trie les conditions dans un ordre canonique déterministe
// Si l'opérateur est commutatif (AND, OR), les conditions sont triées
// Si l'opérateur est non-commutatif, l'ordre est préservé
func NormalizeConditions(conditions []SimpleCondition, operator string) []SimpleCondition {
	// Si pas de conditions ou une seule condition, retourner tel quel
	if len(conditions) <= 1 {
		return conditions
	}

	// Si l'opérateur n'est pas commutatif, préserver l'ordre original
	if !IsCommutative(operator) {
		return conditions
	}

	// Créer une copie pour ne pas modifier l'original
	normalized := make([]SimpleCondition, len(conditions))
	copy(normalized, conditions)

	// Trier par représentation canonique pour un ordre déterministe
	sort.Slice(normalized, func(i, j int) bool {
		return CanonicalString(normalized[i]) < CanonicalString(normalized[j])
	})

	return normalized
}
