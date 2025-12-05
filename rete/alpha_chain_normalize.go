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

// Ce fichier contient les fonctions de normalisation d'expressions logiques
// pour les chaînes alpha du système RETE.
//
// La normalisation consiste à réorganiser les termes d'une expression dans un ordre
// canonique déterministe, permettant ainsi la comparaison et l'optimisation d'expressions
// équivalentes. Par exemple, "A AND B" devient toujours identique à "B AND A" après normalisation.

// NormalizeORExpression normalise une expression OR en triant ses termes dans un ordre canonique.
// Cette fonction extrait les termes OR, les trie et reconstruit l'expression normalisée
// SANS la décomposer en chaîne - l'expression OR reste un seul nœud atomique.
//
// Paramètres:
//   - expr: expression à normaliser (LogicalExpression ou map avec operations OR)
//
// Retourne:
//   - interface{}: expression OR normalisée avec termes triés
//   - error: erreur si l'expression n'est pas une expression OR valide
//
// Exemple:
//
//	Input:  p.status == "VIP" OR p.age > 18
//	Output: p.age > 18 OR p.status == "VIP"  (ordre alphabétique des champs)
func NormalizeORExpression(expr interface{}) (interface{}, error) {
	if expr == nil {
		return nil, fmt.Errorf("expression nil")
	}

	// Vérifier que c'est bien une expression OR
	exprType, err := AnalyzeExpression(expr)
	if err != nil {
		return nil, fmt.Errorf("erreur analyse expression: %w", err)
	}

	if exprType != ExprTypeOR && exprType != ExprTypeMixed {
		return nil, fmt.Errorf("expression n'est pas de type OR ou Mixed: %s", exprType)
	}

	switch e := expr.(type) {
	case constraint.LogicalExpression:
		return normalizeORLogicalExpression(e)
	case map[string]interface{}:
		return normalizeORExpressionMap(e)
	default:
		return nil, fmt.Errorf("type d'expression non supporté pour normalisation OR: %T", expr)
	}
}

// normalizeORLogicalExpression normalise une LogicalExpression contenant des OR
func normalizeORLogicalExpression(expr constraint.LogicalExpression) (constraint.LogicalExpression, error) {
	// Extraire tous les termes OR
	terms := []interface{}{expr.Left}

	for _, op := range expr.Operations {
		opStr := strings.ToUpper(op.Op)
		if opStr == "OR" || opStr == "||" {
			terms = append(terms, op.Right)
		}
	}

	// Convertir chaque terme en string canonique pour le tri
	type termWithCanonical struct {
		term      interface{}
		canonical string
	}

	termsWithCanonical := make([]termWithCanonical, len(terms))
	for i, term := range terms {
		// Créer une représentation canonique pour le tri
		canonical := canonicalValue(term)
		termsWithCanonical[i] = termWithCanonical{
			term:      term,
			canonical: canonical,
		}
	}

	// Trier par représentation canonique
	sort.Slice(termsWithCanonical, func(i, j int) bool {
		return termsWithCanonical[i].canonical < termsWithCanonical[j].canonical
	})

	// Reconstruire l'expression avec les termes triés
	if len(termsWithCanonical) == 0 {
		return constraint.LogicalExpression{}, fmt.Errorf("aucun terme trouvé")
	}

	normalized := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: termsWithCanonical[0].term,
	}

	// Reconstruire les opérations OR
	for i := 1; i < len(termsWithCanonical); i++ {
		normalized.Operations = append(normalized.Operations, constraint.LogicalOperation{
			Op:    "OR",
			Right: termsWithCanonical[i].term,
		})
	}

	return normalized, nil
}

// normalizeORExpressionMap normalise une expression OR au format map
func normalizeORExpressionMap(expr map[string]interface{}) (map[string]interface{}, error) {
	// Extraire tous les termes OR
	terms := []interface{}{}

	if left, ok := expr["left"]; ok {
		terms = append(terms, left)
	}

	if operations, ok := expr["operations"]; ok {
		// Supporter différents formats d'opérations
		if opsList, ok := operations.([]interface{}); ok {
			for _, opInterface := range opsList {
				if opMap, ok := opInterface.(map[string]interface{}); ok {
					if op, ok := opMap["op"].(string); ok {
						opStr := strings.ToUpper(op)
						if opStr == "OR" || opStr == "||" {
							if right, ok := opMap["right"]; ok {
								terms = append(terms, right)
							}
						}
					}
				}
			}
		} else if opsMapList, ok := operations.([]map[string]interface{}); ok {
			for _, opMap := range opsMapList {
				if op, ok := opMap["op"].(string); ok {
					opStr := strings.ToUpper(op)
					if opStr == "OR" || opStr == "||" {
						if right, ok := opMap["right"]; ok {
							terms = append(terms, right)
						}
					}
				}
			}
		}
	}

	if len(terms) == 0 {
		return nil, fmt.Errorf("aucun terme OR trouvé")
	}

	// Convertir chaque terme en string canonique pour le tri
	type termWithCanonical struct {
		term      interface{}
		canonical string
	}

	termsWithCanonical := make([]termWithCanonical, len(terms))
	for i, term := range terms {
		canonical := canonicalValue(term)
		termsWithCanonical[i] = termWithCanonical{
			term:      term,
			canonical: canonical,
		}
	}

	// Trier par représentation canonique
	sort.Slice(termsWithCanonical, func(i, j int) bool {
		return termsWithCanonical[i].canonical < termsWithCanonical[j].canonical
	})

	// Reconstruire l'expression map avec les termes triés
	normalized := map[string]interface{}{
		"type": expr["type"],
		"left": termsWithCanonical[0].term,
	}

	// Reconstruire les opérations OR
	normalizedOps := make([]map[string]interface{}, 0)
	for i := 1; i < len(termsWithCanonical); i++ {
		normalizedOps = append(normalizedOps, map[string]interface{}{
			"op":    "OR",
			"right": termsWithCanonical[i].term,
		})
	}

	if len(normalizedOps) > 0 {
		normalized["operations"] = normalizedOps
	}

	return normalized, nil
}

// NormalizeExpression normalise une expression en appliquant un ordre canonique
// aux conditions quand l'opérateur est commutatif.
//
// Cette fonction gère différents types d'expressions:
//   - LogicalExpression: normalise les opérations logiques
//   - BinaryOperation: retourne tel quel (déjà atomique)
//   - Map: délègue à normalizeExpressionMap
//   - Autres types: retourne tel quel
func NormalizeExpression(expr interface{}) (interface{}, error) {
	switch e := expr.(type) {
	case constraint.LogicalExpression:
		return normalizeLogicalExpression(e)

	case constraint.BinaryOperation:
		// Les opérations binaires simples ne nécessitent pas de normalisation
		// mais on peut normaliser récursivement les sous-expressions
		return e, nil

	case constraint.Constraint:
		// Les contraintes simples ne nécessitent pas de normalisation
		return e, nil

	case map[string]interface{}:
		return normalizeExpressionMap(e)

	default:
		// Pour les autres types (literals, field access), retourner tel quel
		return expr, nil
	}
}

// normalizeLogicalExpression normalise une expression logique
// Si tous les opérateurs sont identiques et commutatifs (AND/OR), les conditions
// sont extraites, triées et reconstruites dans un ordre canonique.
func normalizeLogicalExpression(expr constraint.LogicalExpression) (constraint.LogicalExpression, error) {
	// Si pas d'opérations, retourner tel quel
	if len(expr.Operations) == 0 {
		return expr, nil
	}

	// Déterminer si tous les opérateurs sont identiques et commutatifs
	firstOp := expr.Operations[0].Op
	allSame := true
	for _, op := range expr.Operations {
		if op.Op != firstOp {
			allSame = false
			break
		}
	}

	// Si les opérateurs ne sont pas tous identiques ou si non-commutatif, retourner tel quel
	if !allSame || !IsCommutative(firstOp) {
		return expr, nil
	}

	// Extraire toutes les conditions
	conditions, _, err := extractFromLogicalExpression(expr)
	if err != nil {
		return expr, err
	}

	// Normaliser l'ordre des conditions
	normalized := NormalizeConditions(conditions, firstOp)

	// Reconstruire l'expression logique avec les conditions normalisées
	if len(normalized) == 0 {
		return expr, nil
	}

	// Reconstruire l'expression avec les conditions normalisées
	rebuiltExpr, err := rebuildLogicalExpression(normalized, firstOp)
	if err != nil {
		return expr, err
	}

	return rebuiltExpr, nil
}

// normalizeExpressionMap normalise une expression sous forme de map
// Détecte le type d'expression et applique la normalisation appropriée.
func normalizeExpressionMap(expr map[string]interface{}) (map[string]interface{}, error) {
	exprType, ok := expr["type"].(string)
	if !ok {
		return expr, nil
	}

	switch exprType {
	case "logicalExpression", "logical_op", "logicalExpr":
		// Extraire les conditions
		conditions, opType, err := extractFromLogicalExpressionMap(expr)
		if err != nil {
			return expr, err
		}

		// Si l'opérateur est commutatif, normaliser et reconstruire
		if IsCommutative(opType) {
			normalized := NormalizeConditions(conditions, opType)
			rebuiltExpr, err := rebuildLogicalExpressionMap(normalized, opType)
			if err != nil {
				return expr, err
			}
			return rebuiltExpr, nil
		}
		return expr, nil

	case "binaryOperation", "binary_op", "comparison":
		// Les opérations binaires simples ne nécessitent pas de normalisation
		return expr, nil

	default:
		return expr, nil
	}
}
