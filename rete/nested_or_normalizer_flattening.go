// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strings"

	"github.com/treivax/tsd/constraint"
)

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
// collectORTermsFromMap collecte récursivement tous les termes OR d'une expression map
func collectORTermsFromMap(expr map[string]interface{}) []interface{} {
	terms := []interface{}{}

	collectLeftTerms(expr, &terms)
	collectOperationTerms(expr, &terms)

	return terms
}

// collectLeftTerms collecte les termes du côté gauche de l'expression
func collectLeftTerms(expr map[string]interface{}, terms *[]interface{}) {
	left, ok := expr["left"]
	if !ok {
		return
	}

	leftMap, ok := left.(map[string]interface{})
	if !ok {
		*terms = append(*terms, left)
		return
	}

	if shouldFlattenLeftMap(leftMap) {
		*terms = append(*terms, collectORTermsFromMap(leftMap)...)
	} else {
		*terms = append(*terms, left)
	}
}

// shouldFlattenLeftMap détermine si une map left doit être aplatie
func shouldFlattenLeftMap(leftMap map[string]interface{}) bool {
	leftType, ok := leftMap["type"].(string)
	return ok && leftType == "logicalExpr" && hasOnlyORInMap(leftMap)
}

// collectOperationTerms collecte les termes des opérations
func collectOperationTerms(expr map[string]interface{}, terms *[]interface{}) {
	operations, ok := expr["operations"]
	if !ok {
		return
	}

	opsList, ok := operations.([]interface{})
	if !ok {
		return
	}

	for _, opInterface := range opsList {
		collectSingleOperationTerm(opInterface, terms)
	}
}

// collectSingleOperationTerm collecte un terme d'une seule opération
func collectSingleOperationTerm(opInterface interface{}, terms *[]interface{}) {
	opMap, ok := opInterface.(map[string]interface{})
	if !ok {
		return
	}

	op, ok := opMap["op"].(string)
	if !ok {
		return
	}

	opStr := strings.ToUpper(op)
	if opStr != "OR" && opStr != "||" {
		return
	}

	collectRightTerm(opMap, terms)
}

// collectRightTerm collecte le terme du côté droit d'une opération
func collectRightTerm(opMap map[string]interface{}, terms *[]interface{}) {
	right, ok := opMap["right"]
	if !ok {
		return
	}

	rightMap, ok := right.(map[string]interface{})
	if !ok {
		*terms = append(*terms, right)
		return
	}

	if shouldFlattenRightMap(rightMap) {
		*terms = append(*terms, collectORTermsFromMap(rightMap)...)
	} else {
		*terms = append(*terms, right)
	}
}

// shouldFlattenRightMap détermine si une map right doit être aplatie
func shouldFlattenRightMap(rightMap map[string]interface{}) bool {
	rightType, ok := rightMap["type"].(string)
	return ok && rightType == "logicalExpr" && hasOnlyORInMap(rightMap)
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
