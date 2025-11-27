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

// NestedORComplexity représente la complexité d'une expression OR imbriquée
type NestedORComplexity int

const (
	// ComplexitySimple indique une expression sans imbrication
	ComplexitySimple NestedORComplexity = iota
	// ComplexityFlat indique des OR au même niveau (A OR B OR C)
	ComplexityFlat
	// ComplexityNestedOR indique des OR imbriqués (A OR (B OR C))
	ComplexityNestedOR
	// ComplexityMixedANDOR indique un mélange AND/OR ((A OR B) AND C)
	ComplexityMixedANDOR
	// ComplexityDNFCandidate indique une structure candidate pour DNF ((A OR B) AND (C OR D))
	ComplexityDNFCandidate
)

// NestedORAnalysis contient l'analyse d'une expression avec OR imbriqués
type NestedORAnalysis struct {
	Complexity         NestedORComplexity
	NestingDepth       int
	RequiresDNF        bool
	RequiresFlattening bool
	ORTermCount        int
	ANDTermCount       int
	OptimizationHint   string
}

// AnalyzeNestedOR analyse la complexité d'une expression avec OR potentiellement imbriqués
func AnalyzeNestedOR(expr interface{}) (*NestedORAnalysis, error) {
	if expr == nil {
		return nil, fmt.Errorf("expression nil")
	}

	analysis := &NestedORAnalysis{
		Complexity:   ComplexitySimple,
		NestingDepth: 0,
	}

	switch e := expr.(type) {
	case constraint.LogicalExpression:
		analyzeLogicalExpressionNesting(e, analysis, 0)
	case map[string]interface{}:
		analyzeMapExpressionNesting(e, analysis, 0)
	default:
		// Expression simple, pas d'imbrication
		return analysis, nil
	}

	// Déterminer si transformation DNF est recommandée
	if analysis.Complexity == ComplexityDNFCandidate {
		analysis.RequiresDNF = true
		analysis.OptimizationHint = "DNF transformation recommended for better node sharing"
	}

	// Déterminer si aplatissement est nécessaire
	if analysis.Complexity == ComplexityNestedOR && analysis.NestingDepth > 0 {
		analysis.RequiresFlattening = true
		analysis.OptimizationHint = "OR flattening required to normalize expression"
	}

	return analysis, nil
}

// analyzeLogicalExpressionNesting analyse l'imbrication d'une LogicalExpression
func analyzeLogicalExpressionNesting(expr constraint.LogicalExpression, analysis *NestedORAnalysis, depth int) {
	if depth > analysis.NestingDepth {
		analysis.NestingDepth = depth
	}

	hasOR := false
	hasAND := false

	// Analyser les opérations
	for _, op := range expr.Operations {
		opStr := strings.ToUpper(op.Op)

		if opStr == "OR" || opStr == "||" {
			hasOR = true
			analysis.ORTermCount++
		} else if opStr == "AND" || opStr == "&&" {
			hasAND = true
			analysis.ANDTermCount++
		}

		// Récursion pour analyser les sous-expressions
		if rightLogical, ok := op.Right.(constraint.LogicalExpression); ok {
			analyzeLogicalExpressionNesting(rightLogical, analysis, depth+1)
		}
	}

	// Analyser le left également
	if leftLogical, ok := expr.Left.(constraint.LogicalExpression); ok {
		analyzeLogicalExpressionNesting(leftLogical, analysis, depth+1)
	}

	// Déterminer la complexité
	if hasOR && hasAND {
		// Vérifier si c'est un candidat DNF (plusieurs groupes OR liés par AND)
		if analysis.ORTermCount >= 2 && analysis.ANDTermCount >= 1 {
			if analysis.Complexity < ComplexityDNFCandidate {
				analysis.Complexity = ComplexityDNFCandidate
			}
		} else if analysis.Complexity < ComplexityMixedANDOR {
			analysis.Complexity = ComplexityMixedANDOR
		}
	} else if hasOR {
		// Vérifier si c'est un OR imbriqué (profondeur > 0 signifie imbrication)
		if depth > 0 {
			if analysis.Complexity < ComplexityNestedOR {
				analysis.Complexity = ComplexityNestedOR
			}
		} else if analysis.Complexity < ComplexityFlat {
			analysis.Complexity = ComplexityFlat
		}
	}
}

// analyzeMapExpressionNesting analyse l'imbrication d'une expression map
func analyzeMapExpressionNesting(expr map[string]interface{}, analysis *NestedORAnalysis, depth int) {
	if depth > analysis.NestingDepth {
		analysis.NestingDepth = depth
	}

	hasOR := false
	hasAND := false

	// Analyser les opérations
	if operations, ok := expr["operations"]; ok {
		if opsList, ok := operations.([]interface{}); ok {
			for _, opInterface := range opsList {
				if opMap, ok := opInterface.(map[string]interface{}); ok {
					if op, ok := opMap["op"].(string); ok {
						opStr := strings.ToUpper(op)

						if opStr == "OR" || opStr == "||" {
							hasOR = true
							analysis.ORTermCount++
						} else if opStr == "AND" || opStr == "&&" {
							hasAND = true
							analysis.ANDTermCount++
						}

						// Récursion sur right
						if right, ok := opMap["right"].(map[string]interface{}); ok {
							if rightType, ok := right["type"].(string); ok && rightType == "logicalExpr" {
								analyzeMapExpressionNesting(right, analysis, depth+1)
							}
						}
					}
				}
			}
		}
	}

	// Analyser left
	if left, ok := expr["left"].(map[string]interface{}); ok {
		if leftType, ok := left["type"].(string); ok && leftType == "logicalExpr" {
			analyzeMapExpressionNesting(left, analysis, depth+1)
		}
	}

	// Déterminer la complexité
	if hasOR && hasAND {
		// Vérifier si c'est un candidat DNF (plusieurs groupes OR liés par AND)
		if analysis.ORTermCount >= 2 && analysis.ANDTermCount >= 1 {
			if analysis.Complexity < ComplexityDNFCandidate {
				analysis.Complexity = ComplexityDNFCandidate
			}
		} else if analysis.Complexity < ComplexityMixedANDOR {
			analysis.Complexity = ComplexityMixedANDOR
		}
	} else if hasOR {
		// Vérifier si c'est un OR imbriqué (profondeur > 0 signifie imbrication)
		if depth > 0 {
			if analysis.Complexity < ComplexityNestedOR {
				analysis.Complexity = ComplexityNestedOR
			}
		} else if analysis.Complexity < ComplexityFlat {
			analysis.Complexity = ComplexityFlat
		}
	}
}

// FlattenNestedOR aplatit les OR imbriqués (A OR (B OR C)) en (A OR B OR C)
func FlattenNestedOR(expr interface{}) (interface{}, error) {
	if expr == nil {
		return nil, fmt.Errorf("expression nil")
	}

	switch e := expr.(type) {
	case constraint.LogicalExpression:
		return flattenLogicalExpression(e)
	case map[string]interface{}:
		return flattenMapExpression(e)
	default:
		// Expression simple, rien à aplatir
		return expr, nil
	}
}

// flattenLogicalExpression aplatit une LogicalExpression avec OR imbriqués
func flattenLogicalExpression(expr constraint.LogicalExpression) (constraint.LogicalExpression, error) {
	// Collecter tous les termes OR à tous les niveaux
	orTerms := collectORTermsRecursive(expr)

	if len(orTerms) <= 1 {
		return expr, nil
	}

	// Reconstruire une expression plate
	flattened := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: orTerms[0],
	}

	for i := 1; i < len(orTerms); i++ {
		flattened.Operations = append(flattened.Operations, constraint.LogicalOperation{
			Op:    "OR",
			Right: orTerms[i],
		})
	}

	return flattened, nil
}

// collectORTermsRecursive collecte tous les termes OR de manière récursive
func collectORTermsRecursive(expr constraint.LogicalExpression) []interface{} {
	terms := []interface{}{}

	// Traiter left
	if leftLogical, ok := expr.Left.(constraint.LogicalExpression); ok {
		// Si left est lui-même une expression OR, récurser
		if hasOnlyOR(leftLogical) {
			terms = append(terms, collectORTermsRecursive(leftLogical)...)
		} else {
			terms = append(terms, expr.Left)
		}
	} else {
		terms = append(terms, expr.Left)
	}

	// Traiter operations
	for _, op := range expr.Operations {
		opStr := strings.ToUpper(op.Op)

		if opStr == "OR" || opStr == "||" {
			if rightLogical, ok := op.Right.(constraint.LogicalExpression); ok {
				// Si right est lui-même une expression OR, récurser
				if hasOnlyOR(rightLogical) {
					terms = append(terms, collectORTermsRecursive(rightLogical)...)
				} else {
					terms = append(terms, op.Right)
				}
			} else {
				terms = append(terms, op.Right)
			}
		} else {
			// Pas un OR, ne pas aplatir
			terms = append(terms, op.Right)
		}
	}

	return terms
}

// hasOnlyOR vérifie si une expression ne contient que des opérateurs OR
func hasOnlyOR(expr constraint.LogicalExpression) bool {
	for _, op := range expr.Operations {
		opStr := strings.ToUpper(op.Op)
		if opStr != "OR" && opStr != "||" {
			return false
		}
	}
	return true
}

// flattenMapExpression aplatit une expression map avec OR imbriqués
func flattenMapExpression(expr map[string]interface{}) (map[string]interface{}, error) {
	orTerms := collectORTermsFromMap(expr)

	if len(orTerms) <= 1 {
		return expr, nil
	}

	// Reconstruire une expression map plate
	flattened := map[string]interface{}{
		"type": "logicalExpr",
		"left": orTerms[0],
	}

	operations := make([]map[string]interface{}, 0)
	for i := 1; i < len(orTerms); i++ {
		operations = append(operations, map[string]interface{}{
			"op":    "OR",
			"right": orTerms[i],
		})
	}

	if len(operations) > 0 {
		flattened["operations"] = operations
	}

	return flattened, nil
}

// collectORTermsFromMap collecte les termes OR d'une expression map
func collectORTermsFromMap(expr map[string]interface{}) []interface{} {
	terms := []interface{}{}

	// Traiter left
	if left, ok := expr["left"]; ok {
		if leftMap, ok := left.(map[string]interface{}); ok {
			if leftType, ok := leftMap["type"].(string); ok && leftType == "logicalExpr" {
				if hasOnlyORInMap(leftMap) {
					terms = append(terms, collectORTermsFromMap(leftMap)...)
				} else {
					terms = append(terms, left)
				}
			} else {
				terms = append(terms, left)
			}
		} else {
			terms = append(terms, left)
		}
	}

	// Traiter operations
	if operations, ok := expr["operations"]; ok {
		if opsList, ok := operations.([]interface{}); ok {
			for _, opInterface := range opsList {
				if opMap, ok := opInterface.(map[string]interface{}); ok {
					if op, ok := opMap["op"].(string); ok {
						opStr := strings.ToUpper(op)

						if opStr == "OR" || opStr == "||" {
							if right, ok := opMap["right"]; ok {
								if rightMap, ok := right.(map[string]interface{}); ok {
									if rightType, ok := rightMap["type"].(string); ok && rightType == "logicalExpr" {
										if hasOnlyORInMap(rightMap) {
											terms = append(terms, collectORTermsFromMap(rightMap)...)
										} else {
											terms = append(terms, right)
										}
									} else {
										terms = append(terms, right)
									}
								} else {
									terms = append(terms, right)
								}
							}
						}
					}
				}
			}
		}
	}

	return terms
}

// hasOnlyORInMap vérifie si une expression map ne contient que des OR
func hasOnlyORInMap(expr map[string]interface{}) bool {
	if operations, ok := expr["operations"]; ok {
		if opsList, ok := operations.([]interface{}); ok {
			for _, opInterface := range opsList {
				if opMap, ok := opInterface.(map[string]interface{}); ok {
					if op, ok := opMap["op"].(string); ok {
						opStr := strings.ToUpper(op)
						if opStr != "OR" && opStr != "||" {
							return false
						}
					}
				}
			}
		}
	}
	return true
}

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
