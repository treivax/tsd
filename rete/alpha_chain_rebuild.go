// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"

	"github.com/treivax/tsd/constraint"
)

// Ce fichier contient les fonctions de reconstruction d'expressions logiques
// à partir de conditions normalisées pour les chaînes alpha du système RETE.
//
// Le processus de reconstruction permet de reconstituer une expression complexe
// après avoir extrait, normalisé et trié ses conditions individuelles.

// rebuildLogicalExpression reconstruit une expression logique à partir de conditions normalisées
// Prend une liste de conditions et un opérateur logique (AND/OR) et reconstruit
// une LogicalExpression structurée.
//
// Format résultant:
//
//	Left: première condition
//	Operations: [Op: operator, Right: condition_i] pour i=1..n
func rebuildLogicalExpression(conditions []SimpleCondition, operator string) (constraint.LogicalExpression, error) {
	if len(conditions) == 0 {
		return constraint.LogicalExpression{}, fmt.Errorf("cannot rebuild expression from empty conditions")
	}

	// Cas simple : une seule condition
	if len(conditions) == 1 {
		cond := conditions[0]
		return constraint.LogicalExpression{
			Type:       "logicalExpr",
			Left:       rebuildConditionAsExpression(cond),
			Operations: []constraint.LogicalOperation{},
		}, nil
	}

	// Cas avec plusieurs conditions : créer une chaîne d'opérations
	// Le premier élément devient Left, les autres deviennent Operations
	rebuiltExpr := constraint.LogicalExpression{
		Type:       "logicalExpr",
		Left:       rebuildConditionAsExpression(conditions[0]),
		Operations: make([]constraint.LogicalOperation, 0, len(conditions)-1),
	}

	// Ajouter les conditions restantes comme opérations
	for i := 1; i < len(conditions); i++ {
		rebuiltExpr.Operations = append(rebuiltExpr.Operations, constraint.LogicalOperation{
			Op:    operator,
			Right: rebuildConditionAsExpression(conditions[i]),
		})
	}

	return rebuiltExpr, nil
}

// rebuildConditionAsExpression convertit une SimpleCondition en expression utilisable
// Crée une BinaryOperation à partir d'une SimpleCondition pour l'intégration
// dans une expression logique plus large.
func rebuildConditionAsExpression(cond SimpleCondition) interface{} {
	// Créer une BinaryOperation à partir de la SimpleCondition
	return constraint.BinaryOperation{
		Type:     cond.Type,
		Left:     cond.Left,
		Operator: cond.Operator,
		Right:    cond.Right,
	}
}

// rebuildLogicalExpressionMap reconstruit une expression map à partir de conditions normalisées
// Version map de rebuildLogicalExpression pour compatibilité avec le format map du parser.
//
// Retourne une map avec la structure:
//
//	type: "logicalExpr"
//	left: première condition sous forme de map
//	operations: liste d'opérations [{op: operator, right: condition_map}]
func rebuildLogicalExpressionMap(conditions []SimpleCondition, operator string) (map[string]interface{}, error) {
	if len(conditions) == 0 {
		return nil, fmt.Errorf("cannot rebuild expression from empty conditions")
	}

	// Cas simple : une seule condition
	if len(conditions) == 1 {
		cond := conditions[0]
		return map[string]interface{}{
			"type":       "logicalExpr",
			"left":       rebuildConditionAsMap(cond),
			"operations": []interface{}{},
		}, nil
	}

	// Cas avec plusieurs conditions
	operations := make([]interface{}, 0, len(conditions)-1)
	for i := 1; i < len(conditions); i++ {
		operations = append(operations, map[string]interface{}{
			"op":    operator,
			"right": rebuildConditionAsMap(conditions[i]),
		})
	}

	return map[string]interface{}{
		"type":       "logicalExpr",
		"left":       rebuildConditionAsMap(conditions[0]),
		"operations": operations,
	}, nil
}

// rebuildConditionAsMap convertit une SimpleCondition en map
// Crée une représentation map d'une condition pour compatibilité avec
// le format utilisé par le parser.
func rebuildConditionAsMap(cond SimpleCondition) map[string]interface{} {
	return map[string]interface{}{
		"type":     cond.Type,
		"left":     cond.Left,
		"operator": cond.Operator,
		"right":    cond.Right,
	}
}
