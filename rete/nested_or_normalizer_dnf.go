// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sort"
	"strings"

	"github.com/treivax/tsd/constraint"
)

// TransformToDNF transforme une expression en forme normale disjonctive (DNF)
// Exemple: (A OR B) AND (C OR D) -> (A AND C) OR (A AND D) OR (B AND C) OR (B AND D)
func TransformToDNF(expr interface{}) (interface{}, error) {
	if expr == nil {
		return nil, fmt.Errorf("expression nil")
	}

	// Analyser d'abord la complexité
	analysis, err := AnalyzeNestedOR(expr)
	if err != nil {
		return nil, fmt.Errorf("erreur analyse: %w", err)
	}

	// Si pas besoin de DNF, retourner tel quel
	if !analysis.RequiresDNF && analysis.Complexity != ComplexityDNFCandidate {
		return expr, nil
	}

	switch e := expr.(type) {
	case constraint.LogicalExpression:
		return transformLogicalExpressionToDNF(e)
	case map[string]interface{}:
		return transformMapToDNF(e)
	default:
		return expr, nil
	}
}

// transformLogicalExpressionToDNF transforme une LogicalExpression en DNF
func transformLogicalExpressionToDNF(expr constraint.LogicalExpression) (constraint.LogicalExpression, error) {
	// Extraire les termes AND et OR
	andGroups := extractANDGroups(expr)

	if len(andGroups) <= 1 {
		// Pas de structure AND/OR complexe
		return expr, nil
	}

	// Générer le produit cartésien pour créer la DNF
	dnfTerms := generateDNFTerms(andGroups)

	// Normaliser et trier les termes DNF
	normalizedTerms := normalizeDNFTerms(dnfTerms)

	// Construire l'expression DNF finale
	if len(normalizedTerms) == 0 {
		return expr, fmt.Errorf("aucun terme DNF généré")
	}

	dnfExpr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: normalizedTerms[0],
	}

	for i := 1; i < len(normalizedTerms); i++ {
		dnfExpr.Operations = append(dnfExpr.Operations, constraint.LogicalOperation{
			Op:    "OR",
			Right: normalizedTerms[i],
		})
	}

	return dnfExpr, nil
}

// extractANDGroups extrait les groupes de termes liés par AND
func extractANDGroups(expr constraint.LogicalExpression) [][]interface{} {
	groups := [][]interface{}{}
	currentGroup := []interface{}{expr.Left}

	for _, op := range expr.Operations {
		opStr := strings.ToUpper(op.Op)

		if opStr == "AND" || opStr == "&&" {
			currentGroup = append(currentGroup, op.Right)
		} else if opStr == "OR" || opStr == "||" {
			// Nouveau groupe
			if len(currentGroup) > 0 {
				groups = append(groups, currentGroup)
			}
			currentGroup = []interface{}{op.Right}
		}
	}

	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups
}

// generateDNFTerms génère les termes DNF par produit cartésien
func generateDNFTerms(andGroups [][]interface{}) [][]interface{} {
	if len(andGroups) == 0 {
		return nil
	}

	// Extraire les termes OR de chaque groupe AND
	orTermsByGroup := make([][][]interface{}, len(andGroups))

	for i, group := range andGroups {
		orTermsByGroup[i] = [][]interface{}{}
		for _, term := range group {
			if logicalTerm, ok := term.(constraint.LogicalExpression); ok {
				// Si le terme est une expression OR, extraire ses termes
				orTerms := collectORTermsRecursive(logicalTerm)
				orTermsByGroup[i] = append(orTermsByGroup[i], orTerms)
			} else {
				// Terme simple
				orTermsByGroup[i] = append(orTermsByGroup[i], []interface{}{term})
			}
		}
	}

	// Générer le produit cartésien
	result := [][]interface{}{{}}

	for _, orTerms := range orTermsByGroup {
		for _, termList := range orTerms {
			newResult := [][]interface{}{}
			for _, existing := range result {
				for _, term := range termList {
					combination := append([]interface{}{}, existing...)
					combination = append(combination, term)
					newResult = append(newResult, combination)
				}
			}
			result = newResult
		}
	}

	return result
}

// normalizeDNFTerms normalise et trie les termes DNF
func normalizeDNFTerms(dnfTerms [][]interface{}) []interface{} {
	normalized := make([]interface{}, 0, len(dnfTerms))

	for _, terms := range dnfTerms {
		if len(terms) == 1 {
			// Terme simple
			normalized = append(normalized, terms[0])
		} else {
			// Terme AND composé - créer une LogicalExpression
			andExpr := constraint.LogicalExpression{
				Type: "logicalExpr",
				Left: terms[0],
			}

			for i := 1; i < len(terms); i++ {
				andExpr.Operations = append(andExpr.Operations, constraint.LogicalOperation{
					Op:    "AND",
					Right: terms[i],
				})
			}

			normalized = append(normalized, andExpr)
		}
	}

	// Trier par représentation canonique
	type termWithCanonical struct {
		term      interface{}
		canonical string
	}

	termsWithCanonical := make([]termWithCanonical, len(normalized))
	for i, term := range normalized {
		canonical := canonicalValue(term)
		termsWithCanonical[i] = termWithCanonical{
			term:      term,
			canonical: canonical,
		}
	}

	sort.Slice(termsWithCanonical, func(i, j int) bool {
		return termsWithCanonical[i].canonical < termsWithCanonical[j].canonical
	})

	result := make([]interface{}, len(termsWithCanonical))
	for i, twc := range termsWithCanonical {
		result[i] = twc.term
	}

	return result
}

// transformMapToDNF transforme une expression map en DNF
func transformMapToDNF(expr map[string]interface{}) (map[string]interface{}, error) {
	// Conversion map -> LogicalExpression pour simplifier
	// (implémentation simplifiée, peut être étendue)
	return expr, nil
}
