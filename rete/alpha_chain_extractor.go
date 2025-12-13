// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package rete fournit l'implémentation du réseau RETE pour l'évaluation de règles.
// Ce fichier contient des utilitaires pour extraire et analyser les conditions d'expressions complexes.
//
// Ce fichier contient les fonctions core d'extraction de conditions depuis expressions complexes.
// Pour les fonctions spécialisées :
// - Représentation canonique : alpha_chain_canonical.go
// - Normalisation : alpha_chain_normalize.go
// - Reconstruction : alpha_chain_rebuild.go
// - Comparaison : alpha_chain_compare.go
//
// Exemple d'utilisation:
//
//	// Expression AND: p.age > 18 AND p.salary >= 50000
//	expr := constraint.LogicalExpression{
//		Type: "logicalExpr",
//		Left: constraint.BinaryOperation{
//			Type:     "binaryOperation",
//			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
//			Operator: ">",
//			Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
//		},
//		Operations: []constraint.LogicalOperation{
//			{
//				Op: "AND",
//				Right: constraint.BinaryOperation{
//					Type:     "binaryOperation",
//					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
//					Operator: ">=",
//					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
//				},
//			},
//		},
//	}
//
//	// Extraire les conditions
//	conditions, opType, err := ExtractConditions(expr)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Type d'opérateur: %s\n", opType) // Output: AND
//	fmt.Printf("Nombre de conditions: %d\n", len(conditions)) // Output: 2
//
//	// Générer des représentations canoniques
//	for _, cond := range conditions {
//		canonical := CanonicalString(cond)
//		fmt.Printf("Condition: %s\n", canonical)
//		fmt.Printf("Hash: %s\n", cond.Hash)
//	}
//
//	// Dédupliquer les conditions
//	uniqueConditions := DeduplicateConditions(conditions)
package rete

import (
	"fmt"

	"github.com/treivax/tsd/constraint"
)

// SimpleCondition représente une condition atomique extraite d'une expression complexe
type SimpleCondition struct {
	Type     string      `json:"type"`     // Type de condition: binaryOperation, comparison, arithmetic, etc.
	Left     interface{} `json:"left"`     // Opérande gauche
	Operator string      `json:"operator"` // Opérateur
	Right    interface{} `json:"right"`    // Opérande droite
	Hash     string      `json:"hash"`     // Hash unique calculé automatiquement
}

// DecomposedCondition extends SimpleCondition with decomposition metadata
// for supporting intermediate result propagation in alpha chains
type DecomposedCondition struct {
	SimpleCondition
	ResultName   string   `json:"result_name,omitempty"`  // Name of intermediate result produced (e.g., "temp_1")
	Dependencies []string `json:"dependencies,omitempty"` // Required intermediate results
	IsAtomic     bool     `json:"is_atomic,omitempty"`    // true if atomic operation
}

// NewSimpleCondition crée une nouvelle condition simple avec hash calculé
func NewSimpleCondition(condType string, left interface{}, operator string, right interface{}) SimpleCondition {
	cond := SimpleCondition{
		Type:     condType,
		Left:     left,
		Operator: operator,
		Right:    right,
	}
	cond.Hash = computeHash(cond)
	return cond
}

// ExtractConditions extrait toutes les conditions simples d'une expression complexe
// Retourne: liste de conditions, type d'opérateur principal (AND/OR/SINGLE), erreur
func ExtractConditions(expr interface{}) ([]SimpleCondition, string, error) {
	switch e := expr.(type) {
	case map[string]interface{}:
		return extractFromMap(e)

	case constraint.BinaryOperation:
		cond := NewSimpleCondition("binaryOperation", e.Left, e.Operator, e.Right)
		return []SimpleCondition{cond}, "SINGLE", nil

	case constraint.LogicalExpression:
		return extractFromLogicalExpression(e)

	case constraint.Constraint:
		return extractFromConstraint(e)

	case constraint.FieldAccess:
		// Un accès de champ seul n'est pas une condition
		return []SimpleCondition{}, "NONE", nil

	case constraint.NumberLiteral, constraint.StringLiteral, constraint.BooleanLiteral:
		// Les littéraux seuls ne sont pas des conditions
		return []SimpleCondition{}, "NONE", nil

	default:
		return nil, "", fmt.Errorf("type d'expression non supporté: %T", expr)
	}
}

// getMapType extrait le type d'une expression sous forme de map
func getMapType(m map[string]interface{}) (string, error) {
	exprType, ok := m["type"].(string)
	if !ok {
		return "", fmt.Errorf("type d'expression manquant")
	}
	return exprType, nil
}

// extractFromMap extrait les conditions d'une expression sous forme de map
func extractFromMap(expr map[string]interface{}) ([]SimpleCondition, string, error) {
	exprType, err := getMapType(expr)
	if err != nil {
		return nil, "", err
	}

	switch exprType {
	case "binaryOperation", "binary_op", "comparison":
		operator, ok := expr["operator"].(string)
		if !ok {
			if operator, ok = expr["op"].(string); !ok {
				return nil, "", fmt.Errorf("opérateur manquant")
			}
		}
		cond := NewSimpleCondition("binaryOperation", expr["left"], operator, expr["right"])
		return []SimpleCondition{cond}, "SINGLE", nil

	case "logicalExpression", "logical_op", "logicalExpr":
		return extractFromLogicalExpressionMap(expr)

	case "constraint":
		// Les contraintes dans les maps peuvent avoir différentes structures
		// On essaie d'extraire left/operator/right directement
		if left, ok := expr["left"]; ok {
			operator, _ := expr["operator"].(string)
			right := expr["right"]
			cond := NewSimpleCondition("constraint", left, operator, right)
			return []SimpleCondition{cond}, "SINGLE", nil
		}
		return []SimpleCondition{}, "NONE", nil

	case "fieldAccess":
		return []SimpleCondition{}, "NONE", nil

	case "literal", "numberLiteral", "stringLiteral", "booleanLiteral":
		return []SimpleCondition{}, "NONE", nil

	default:
		return nil, "", fmt.Errorf("type d'expression map non supporté: %s", exprType)
	}
}

// extractFromLogicalExpression extrait les conditions d'une expression logique
func extractFromLogicalExpression(expr constraint.LogicalExpression) ([]SimpleCondition, string, error) {
	allConditions := []SimpleCondition{}
	operatorType := ""

	// Extraire les conditions du côté gauche
	leftConds, _, err := ExtractConditions(expr.Left)
	if err != nil {
		return nil, "", fmt.Errorf("erreur extraction left: %w", err)
	}
	allConditions = append(allConditions, leftConds...)

	// Traiter toutes les opérations
	for _, op := range expr.Operations {
		rightConds, _, err := ExtractConditions(op.Right)
		if err != nil {
			return nil, "", fmt.Errorf("erreur extraction right: %w", err)
		}
		allConditions = append(allConditions, rightConds...)

		// Déterminer le type d'opérateur principal
		if operatorType == "" {
			operatorType = op.Op
		} else if operatorType != op.Op {
			// Mélange d'opérateurs - retourner "MIXED"
			operatorType = "MIXED"
		}
	}

	if operatorType == "" {
		operatorType = "SINGLE"
	}

	return allConditions, operatorType, nil
}

// extractFromLogicalExpressionMap extrait les conditions d'une expression logique (format map)
func extractFromLogicalExpressionMap(expr map[string]interface{}) ([]SimpleCondition, string, error) {
	allConditions := []SimpleCondition{}
	operatorType := ""

	// Extraire les conditions du côté gauche
	left, ok := expr["left"]
	if !ok {
		return nil, "", fmt.Errorf("left manquant dans logicalExpression")
	}

	leftConds, _, err := ExtractConditions(left)
	if err != nil {
		return nil, "", fmt.Errorf("erreur extraction left: %w", err)
	}
	allConditions = append(allConditions, leftConds...)

	// Traiter toutes les opérations
	operations, ok := expr["operations"]
	if !ok {
		return allConditions, "SINGLE", nil
	}

	// Supporter []interface{}, []map[string]interface{} et []constraint.LogicalOperation
	// Essayer []map[string]interface{} en premier (type le plus courant du parser)
	if opsMapList, ok := operations.([]map[string]interface{}); ok {
		for _, opMap := range opsMapList {
			op, ok := opMap["op"].(string)
			if !ok {
				return nil, "", fmt.Errorf("op manquant dans operation")
			}

			right, ok := opMap["right"]
			if !ok {
				return nil, "", fmt.Errorf("right manquant dans operation")
			}

			rightConds, _, err := ExtractConditions(right)
			if err != nil {
				return nil, "", fmt.Errorf("erreur extraction right: %w", err)
			}
			allConditions = append(allConditions, rightConds...)

			// Déterminer le type d'opérateur principal
			if operatorType == "" {
				operatorType = op
			} else if operatorType != op {
				operatorType = "MIXED"
			}
		}
	} else if opsList, ok := operations.([]interface{}); ok {
		for _, opInterface := range opsList {
			opMap, ok := opInterface.(map[string]interface{})
			if !ok {
				return nil, "", fmt.Errorf("operation doit être une map")
			}

			op, ok := opMap["op"].(string)
			if !ok {
				return nil, "", fmt.Errorf("op manquant dans operation")
			}

			right, ok := opMap["right"]
			if !ok {
				return nil, "", fmt.Errorf("right manquant dans operation")
			}

			rightConds, _, err := ExtractConditions(right)
			if err != nil {
				return nil, "", fmt.Errorf("erreur extraction right: %w", err)
			}
			allConditions = append(allConditions, rightConds...)

			// Déterminer le type d'opérateur principal
			if operatorType == "" {
				operatorType = op
			} else if operatorType != op {
				operatorType = "MIXED"
			}
		}
	} else if logicalOps, ok := operations.([]constraint.LogicalOperation); ok {
		// Supporter le type constraint.LogicalOperation directement
		for _, op := range logicalOps {
			rightConds, _, err := ExtractConditions(op.Right)
			if err != nil {
				return nil, "", fmt.Errorf("erreur extraction right: %w", err)
			}
			allConditions = append(allConditions, rightConds...)

			// Déterminer le type d'opérateur principal
			if operatorType == "" {
				operatorType = op.Op
			} else if operatorType != op.Op {
				operatorType = "MIXED"
			}
		}
	} else {
		return nil, "", fmt.Errorf("operations doit être un tableau ([]interface{}, []map[string]interface{} ou []LogicalOperation)")
	}

	if operatorType == "" {
		operatorType = "SINGLE"
	}

	return allConditions, operatorType, nil
}

// extractFromNOTConstraint extrait les conditions d'une contrainte NOT
func extractFromNOTConstraint(expr constraint.NotConstraint) ([]SimpleCondition, string, error) {
	// Pour les contraintes NOT, on retourne une condition spéciale
	// qui sera gérée différemment par le constructeur de chaîne
	cond := NewSimpleCondition("not", expr.Expression, "NOT", nil)
	return []SimpleCondition{cond}, "NOT", nil
}

// extractFromNOTConstraintMap extrait les conditions d'une contrainte NOT (format map)
func extractFromNOTConstraintMap(expr map[string]interface{}) ([]SimpleCondition, string, error) {
	expression, ok := expr["expression"]
	if !ok {
		return nil, "", fmt.Errorf("expression manquant dans notConstraint")
	}

	// Pour les contraintes NOT, on retourne une condition spéciale
	cond := NewSimpleCondition("not", expression, "NOT", nil)
	return []SimpleCondition{cond}, "NOT", nil
}

// extractFromConstraint extrait les conditions d'une contrainte
func extractFromConstraint(c constraint.Constraint) ([]SimpleCondition, string, error) {
	// Les contraintes peuvent avoir left/operator/right directement
	if c.Left != nil && c.Operator != "" {
		cond := NewSimpleCondition("constraint", c.Left, c.Operator, c.Right)
		return []SimpleCondition{cond}, "SINGLE", nil
	}
	return []SimpleCondition{}, "NONE", nil
}
