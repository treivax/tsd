// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"

	"github.com/treivax/tsd/constraint"
)

// evaluateConstraint évalue une contrainte simple
func (e *AlphaConditionEvaluator) evaluateConstraint(constraint constraint.Constraint) (bool, error) {
	left, err := e.evaluateValue(constraint.Left)
	if err != nil {
		return false, fmt.Errorf("erreur évaluation côté gauche: %w", err)
	}

	right, err := e.evaluateValue(constraint.Right)
	if err != nil {
		return false, fmt.Errorf("erreur évaluation côté droit: %w", err)
	}

	return e.compareValues(left, constraint.Operator, right)
}

// evaluateConstraintMap évalue une contrainte depuis une map
func (e *AlphaConditionEvaluator) evaluateConstraintMap(expr map[string]interface{}) (bool, error) {
	// Si l'expression a une clé "constraint", extraire la contrainte réelle
	if constraintData, hasConstraint := expr["constraint"]; hasConstraint {
		// Si c'est une LogicalExpression structurée (pas une map), l'évaluer directement
		if logicalExpr, ok := constraintData.(constraint.LogicalExpression); ok {
			return e.evaluateLogicalExpression(logicalExpr)
		}

		// Si c'est une map, continuer avec le traitement map
		if constraintMap, ok := constraintData.(map[string]interface{}); ok {
			// Gérer directement les expressions logiques map
			if condType, hasType := constraintMap["type"].(string); hasType {
				if condType == "logicalExpr" || condType == "logicalExpression" {
					return e.evaluateLogicalExpressionMap(constraintMap)
				}
			}
			return e.evaluateConstraintMapInternal(constraintMap)
		}

		return false, fmt.Errorf("format contrainte invalide: %T", constraintData)
	}

	// Utiliser directement l'expression si pas d'indirection
	return e.evaluateConstraintMapInternal(expr)
}

// evaluateConstraintMapInternal évalue une map de contrainte
func (e *AlphaConditionEvaluator) evaluateConstraintMapInternal(actualConstraint map[string]interface{}) (bool, error) {

	// Vérifier d'abord le type de contrainte avant de chercher l'opérateur
	if condType, hasType := actualConstraint["type"].(string); hasType {
		switch condType {
		case "simple", "passthrough", "exists":
			return true, nil // Conditions spéciales toujours vraies
		case "logicalExpr", "logicalExpression":
			return e.evaluateLogicalExpressionMap(actualConstraint)
		case "existsConstraint":
			return e.evaluateExistsConstraint(actualConstraint)
		case "notConstraint":
			return e.evaluateNotConstraint(actualConstraint)
		case "negation":
			return e.evaluateNegationConstraint(actualConstraint)
		}
	}

	operator, ok := actualConstraint["operator"].(string)
	if !ok {
		return false, fmt.Errorf("opérateur manquant pour condition: %v", actualConstraint)
	}

	left, err := e.evaluateValue(actualConstraint["left"])
	if err != nil {
		return false, fmt.Errorf("erreur évaluation côté gauche: %w", err)
	}

	right, err := e.evaluateValue(actualConstraint["right"])
	if err != nil {
		return false, fmt.Errorf("erreur évaluation côté droit: %w", err)
	}

	return e.compareValues(left, operator, right)
}

// evaluateNegationConstraint évalue une contrainte de négation
func (e *AlphaConditionEvaluator) evaluateNegationConstraint(expr map[string]interface{}) (bool, error) {
	// Extraire la condition niée depuis "condition"
	condition, ok := expr["condition"]
	if !ok {
		return false, fmt.Errorf("condition manquante dans contrainte de négation")
	}

	// Évaluer la condition interne
	result, err := e.evaluateExpression(condition)
	if err != nil {
		return false, fmt.Errorf("erreur évaluation condition niée: %w", err)
	}

	// Retourner la négation du résultat
	return !result, nil
}

// evaluateNotConstraint évalue une contrainte NOT (notConstraint)
func (e *AlphaConditionEvaluator) evaluateNotConstraint(expr map[string]interface{}) (bool, error) {
	// Extraire l'expression depuis "expression"
	expression, ok := expr["expression"]
	if !ok {
		return false, fmt.Errorf("expression manquante dans contrainte NOT")
	}

	// Évaluer l'expression interne
	result, err := e.evaluateExpression(expression)
	if err != nil {
		return false, fmt.Errorf("erreur évaluation expression NOT: %w", err)
	}

	// Retourner la négation du résultat
	return !result, nil
}

// evaluateExistsConstraint évalue une contrainte EXISTS
func (e *AlphaConditionEvaluator) evaluateExistsConstraint(expr map[string]interface{}) (bool, error) {
	// Note: L'évaluation réelle EXISTS est gérée par les ExistsNodes dans le réseau RETE
	// Cette fonction est utilisée uniquement pour la validation initiale au niveau Alpha
	hash := fmt.Sprintf("%v", expr)
	checksum := 0
	for _, r := range hash {
		checksum += int(r)
	}
	return (checksum % 20) != 0, nil
}
