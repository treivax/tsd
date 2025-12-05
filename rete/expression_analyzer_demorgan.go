// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package rete fournit l'implémentation du réseau RETE pour l'évaluation de règles.
// Ce fichier contient les fonctions de transformation De Morgan pour la manipulation
// d'expressions logiques avec négation (NOT(A AND B) ↔ (NOT A) OR (NOT B)).
package rete

import (
	"github.com/treivax/tsd/constraint"
)

// ApplyDeMorganTransformation applique la loi de De Morgan à une expression NOT
// Transforme NOT(A OR B) en (NOT A) AND (NOT B)
// Transforme NOT(A AND B) en (NOT A) OR (NOT B)
// Retourne l'expression transformée et un booléen indiquant si une transformation a eu lieu
func ApplyDeMorganTransformation(expr interface{}) (interface{}, bool) {
	// Vérifier si c'est une expression NOT
	exprType, err := AnalyzeExpression(expr)
	if err != nil || exprType != ExprTypeNOT {
		return expr, false
	}

	// Extraire l'expression interne
	innerExpr := extractInnerExpression(expr)
	if innerExpr == nil {
		return expr, false
	}

	// Analyser l'expression interne
	innerType, err := AnalyzeExpression(innerExpr)
	if err != nil {
		return expr, false
	}

	// Appliquer De Morgan seulement si l'expression interne est AND ou OR
	if innerType != ExprTypeAND && innerType != ExprTypeOR {
		return expr, false
	}

	// Appliquer la transformation selon le type
	switch innerType {
	case ExprTypeAND:
		// NOT(A AND B) -> (NOT A) OR (NOT B)
		return transformNotAnd(innerExpr), true
	case ExprTypeOR:
		// NOT(A OR B) -> (NOT A) AND (NOT B)
		return transformNotOr(innerExpr), true
	}

	return expr, false
}

// transformNotAnd transforme NOT(A AND B) en (NOT A) OR (NOT B)
func transformNotAnd(expr interface{}) interface{} {
	switch e := expr.(type) {
	case constraint.LogicalExpression:
		// Créer une nouvelle expression logique avec OR
		newExpr := constraint.LogicalExpression{
			Left: wrapInNot(e.Left),
		}

		// Transformer toutes les opérations AND en OR et nier les opérandes
		for _, op := range e.Operations {
			newExpr.Operations = append(newExpr.Operations, constraint.LogicalOperation{
				Op:    convertAndToOr(op.Op),
				Right: wrapInNot(op.Right),
			})
		}

		return newExpr

	case map[string]interface{}:
		return transformNotAndMap(e)
	}

	return expr
}

// transformNotOr transforme NOT(A OR B) en (NOT A) AND (NOT B)
func transformNotOr(expr interface{}) interface{} {
	switch e := expr.(type) {
	case constraint.LogicalExpression:
		// Créer une nouvelle expression logique avec AND
		newExpr := constraint.LogicalExpression{
			Left: wrapInNot(e.Left),
		}

		// Transformer toutes les opérations OR en AND et nier les opérandes
		for _, op := range e.Operations {
			newExpr.Operations = append(newExpr.Operations, constraint.LogicalOperation{
				Op:    convertOrToAnd(op.Op),
				Right: wrapInNot(op.Right),
			})
		}

		return newExpr

	case map[string]interface{}:
		return transformNotOrMap(e)
	}

	return expr
}

// transformNotAndMap transforme NOT(A AND B) en (NOT A) OR (NOT B) pour format map
func transformNotAndMap(expr map[string]interface{}) interface{} {
	result := map[string]interface{}{
		"type": "logicalExpression",
	}

	// Extraire le côté gauche
	if left, ok := expr["left"]; ok {
		result["left"] = wrapInNotMap(left)
	}

	// Transformer les opérations
	if operations, ok := expr["operations"].([]interface{}); ok {
		newOps := make([]interface{}, 0, len(operations))
		for _, opInterface := range operations {
			if opMap, ok := opInterface.(map[string]interface{}); ok {
				newOp := map[string]interface{}{
					"op": convertAndToOr(getOperatorFromMap(opMap)),
				}
				if right, ok := opMap["right"]; ok {
					newOp["right"] = wrapInNotMap(right)
				}
				newOps = append(newOps, newOp)
			}
		}
		result["operations"] = newOps
	}

	return result
}

// transformNotOrMap transforme NOT(A OR B) en (NOT A) AND (NOT B) pour format map
func transformNotOrMap(expr map[string]interface{}) interface{} {
	result := map[string]interface{}{
		"type": "logicalExpression",
	}

	// Extraire le côté gauche
	if left, ok := expr["left"]; ok {
		result["left"] = wrapInNotMap(left)
	}

	// Transformer les opérations
	if operations, ok := expr["operations"].([]interface{}); ok {
		newOps := make([]interface{}, 0, len(operations))
		for _, opInterface := range operations {
			if opMap, ok := opInterface.(map[string]interface{}); ok {
				newOp := map[string]interface{}{
					"op": convertOrToAnd(getOperatorFromMap(opMap)),
				}
				if right, ok := opMap["right"]; ok {
					newOp["right"] = wrapInNotMap(right)
				}
				newOps = append(newOps, newOp)
			}
		}
		result["operations"] = newOps
	}

	return result
}

// wrapInNot enveloppe une expression dans un NOT
func wrapInNot(expr interface{}) interface{} {
	return constraint.NotConstraint{
		Expression: expr,
	}
}

// wrapInNotMap enveloppe une expression dans un NOT (format map)
func wrapInNotMap(expr interface{}) interface{} {
	return map[string]interface{}{
		"type":       "notConstraint",
		"expression": expr,
	}
}

// convertAndToOr convertit un opérateur AND en OR
func convertAndToOr(op string) string {
	switch op {
	case "AND", "and":
		return "OR"
	case "&&":
		return "||"
	default:
		return "OR"
	}
}

// convertOrToAnd convertit un opérateur OR en AND
func convertOrToAnd(op string) string {
	switch op {
	case "OR", "or":
		return "AND"
	case "||":
		return "&&"
	default:
		return "AND"
	}
}

// getOperatorFromMap extrait l'opérateur d'une map d'opération
func getOperatorFromMap(opMap map[string]interface{}) string {
	if op, ok := opMap["op"].(string); ok {
		return op
	}
	if op, ok := opMap["operator"].(string); ok {
		return op
	}
	return ""
}
