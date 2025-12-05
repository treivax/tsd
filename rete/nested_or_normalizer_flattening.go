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
