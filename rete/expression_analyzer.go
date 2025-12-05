// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package rete fournit l'implémentation du réseau RETE pour l'évaluation de règles.
// Ce fichier contient l'analyseur d'expressions qui détermine le type d'une expression
// et comment elle doit être traitée dans le réseau RETE.
//
// Exemple d'utilisation:
//
//	// Analyser une expression simple
//	expr := map[string]interface{}{
//		"type": "binaryOperation",
//		"left": map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
//		"operator": ">",
//		"right": map[string]interface{}{"type": "numberLiteral", "value": 18},
//	}
//
//	exprType, err := AnalyzeExpression(expr)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if CanDecompose(exprType) {
//		fmt.Println("Cette expression peut être décomposée en chaîne alpha")
//	}
//
//	if ShouldNormalize(exprType) {
//		fmt.Println("Cette expression nécessite une normalisation")
//	}
package rete

import (
	"fmt"

	"github.com/treivax/tsd/constraint"
)

// ExpressionType représente le type d'une expression analysée
type ExpressionType int

const (
	// ExprTypeSimple représente une condition atomique (p.age > 18)
	ExprTypeSimple ExpressionType = iota

	// ExprTypeAND représente une expression avec uniquement des opérateurs AND
	// Exemple: p.age > 18 AND p.salary >= 50000 AND p.active = true
	ExprTypeAND

	// ExprTypeOR représente une expression avec uniquement des opérateurs OR
	// Exemple: p.status = "active" OR p.status = "pending"
	ExprTypeOR

	// ExprTypeMixed représente une expression avec AND et OR mélangés
	// Exemple: (p.age > 18 AND p.salary >= 50000) OR p.vip = true
	ExprTypeMixed

	// ExprTypeArithmetic représente une chaîne d'opérations arithmétiques commutatives
	// Exemple: p.price * 1.2 + 5 (peut être décomposé en étapes)
	ExprTypeArithmetic

	// ExprTypeNOT représente une négation d'expression
	// Exemple: NOT p.active, NOT (p.age > 18 AND p.salary < 50000)
	ExprTypeNOT
)

// String retourne une représentation textuelle du type d'expression
func (et ExpressionType) String() string {
	switch et {
	case ExprTypeSimple:
		return "ExprTypeSimple"
	case ExprTypeAND:
		return "ExprTypeAND"
	case ExprTypeOR:
		return "ExprTypeOR"
	case ExprTypeMixed:
		return "ExprTypeMixed"
	case ExprTypeArithmetic:
		return "ExprTypeArithmetic"
	case ExprTypeNOT:
		return "ExprTypeNOT"
	default:
		return "Unknown"
	}
}

// AnalyzeExpression analyse une expression et détermine son type
// Retourne le type d'expression identifié et une erreur si l'analyse échoue
func AnalyzeExpression(expr interface{}) (ExpressionType, error) {
	if expr == nil {
		return ExprTypeSimple, fmt.Errorf("expression nil")
	}

	switch e := expr.(type) {
	case map[string]interface{}:
		return analyzeMapExpression(e)

	case constraint.BinaryOperation:
		// Une opération binaire simple peut être une condition ou une opération arithmétique
		if isArithmeticOperator(e.Operator) {
			return ExprTypeArithmetic, nil
		}
		return ExprTypeSimple, nil

	case constraint.LogicalExpression:
		return analyzeLogicalExpression(e)

	case constraint.Constraint:
		// Une contrainte simple
		if e.Left != nil && e.Operator != "" {
			if isArithmeticOperator(e.Operator) {
				return ExprTypeArithmetic, nil
			}
			return ExprTypeSimple, nil
		}
		return ExprTypeSimple, nil

	case constraint.FieldAccess:
		// Un accès de champ seul n'est pas vraiment une condition complète
		return ExprTypeSimple, nil

	case constraint.NumberLiteral, constraint.StringLiteral, constraint.BooleanLiteral:
		// Les littéraux seuls ne forment pas une expression évaluable
		return ExprTypeSimple, nil

	case SimpleCondition:
		// Une SimpleCondition déjà extraite
		return ExprTypeSimple, nil

	case constraint.NotConstraint:
		// Une contrainte de négation
		return ExprTypeNOT, nil

	default:
		return ExprTypeSimple, fmt.Errorf("type d'expression non supporté: %T", expr)
	}
}

// analyzeMapExpression analyse une expression sous forme de map
func analyzeMapExpression(expr map[string]interface{}) (ExpressionType, error) {
	exprType, ok := expr["type"].(string)
	if !ok {
		return ExprTypeSimple, fmt.Errorf("type d'expression manquant")
	}

	switch exprType {
	case "binaryOperation", "binary_op", "comparison":
		operator, ok := expr["operator"].(string)
		if !ok {
			if operator, ok = expr["op"].(string); !ok {
				return ExprTypeSimple, fmt.Errorf("opérateur manquant")
			}
		}
		if isArithmeticOperator(operator) {
			return ExprTypeArithmetic, nil
		}
		return ExprTypeSimple, nil

	case "logicalExpression", "logical_op", "logicalExpr":
		return analyzeLogicalExpressionMap(expr)

	case "constraint":
		// Analyser la contrainte pour déterminer si elle contient des opérations arithmétiques
		if operator, ok := expr["operator"].(string); ok && isArithmeticOperator(operator) {
			return ExprTypeArithmetic, nil
		}
		return ExprTypeSimple, nil

	case "arithmeticOperation", "arithmetic_op":
		return ExprTypeArithmetic, nil

	case "notConstraint", "not", "negation":
		return ExprTypeNOT, nil

	case "parenthesized", "parenthesizedExpression", "group":
		// Une expression parenthésée - analyser le contenu
		return analyzeParenthesizedExpression(expr)

	case "fieldAccess", "literal", "numberLiteral", "stringLiteral", "booleanLiteral":
		return ExprTypeSimple, nil

	default:
		return ExprTypeSimple, fmt.Errorf("type d'expression map non supporté: %s", exprType)
	}
}

// analyzeLogicalExpression analyse une expression logique structurée
func analyzeLogicalExpression(expr constraint.LogicalExpression) (ExpressionType, error) {
	if len(expr.Operations) == 0 {
		// Pas d'opérations - analyser le côté gauche uniquement
		return AnalyzeExpression(expr.Left)
	}

	// Vérifier tous les opérateurs
	hasAND := false
	hasOR := false

	for _, op := range expr.Operations {
		switch op.Op {
		case "AND", "and", "&&":
			hasAND = true
		case "OR", "or", "||":
			hasOR = true
		}
	}

	// Déterminer le type en fonction des opérateurs trouvés
	if hasAND && hasOR {
		return ExprTypeMixed, nil
	} else if hasOR {
		return ExprTypeOR, nil
	} else if hasAND {
		return ExprTypeAND, nil
	}

	// Par défaut, considérer comme simple
	return ExprTypeSimple, nil
}

// analyzeLogicalExpressionMap analyse une expression logique (format map)
func analyzeLogicalExpressionMap(expr map[string]interface{}) (ExpressionType, error) {
	operations, ok := expr["operations"]
	if !ok {
		// Pas d'opérations - analyser le côté gauche uniquement
		if left, ok := expr["left"]; ok {
			return AnalyzeExpression(left)
		}
		return ExprTypeSimple, nil
	}

	// Vérifier tous les opérateurs - supporter []interface{}, []map[string]interface{} et []constraint.LogicalOperation
	hasAND := false
	hasOR := false

	// Essayer []map[string]interface{} en premier (type le plus courant du parser)
	if opsMapList, ok := operations.([]map[string]interface{}); ok {
		if len(opsMapList) == 0 {
			if left, ok := expr["left"]; ok {
				return AnalyzeExpression(left)
			}
			return ExprTypeSimple, nil
		}

		for _, opMap := range opsMapList {
			operator, ok := opMap["op"].(string)
			if !ok {
				continue
			}

			switch operator {
			case "AND", "and", "&&":
				hasAND = true
			case "OR", "or", "||":
				hasOR = true
			}
		}
	} else if opsList, ok := operations.([]interface{}); ok {
		if len(opsList) == 0 {
			if left, ok := expr["left"]; ok {
				return AnalyzeExpression(left)
			}
			return ExprTypeSimple, nil
		}

		for _, opInterface := range opsList {
			opMap, ok := opInterface.(map[string]interface{})
			if !ok {
				continue
			}

			operator, ok := opMap["op"].(string)
			if !ok {
				continue
			}

			switch operator {
			case "AND", "and", "&&":
				hasAND = true
			case "OR", "or", "||":
				hasOR = true
			}
		}
	} else if logicalOps, ok := operations.([]constraint.LogicalOperation); ok {
		if len(logicalOps) == 0 {
			if left, ok := expr["left"]; ok {
				return AnalyzeExpression(left)
			}
			return ExprTypeSimple, nil
		}

		for _, op := range logicalOps {
			switch op.Op {
			case "AND", "and", "&&":
				hasAND = true
			case "OR", "or", "||":
				hasOR = true
			}
		}
	} else {
		return ExprTypeSimple, fmt.Errorf("operations doit être un tableau ([]interface{}, []map[string]interface{} ou []LogicalOperation)")
	}

	// Déterminer le type en fonction des opérateurs trouvés
	if hasAND && hasOR {
		return ExprTypeMixed, nil
	} else if hasOR {
		return ExprTypeOR, nil
	} else if hasAND {
		return ExprTypeAND, nil
	}

	// Par défaut, considérer comme simple
	return ExprTypeSimple, nil
}

// isArithmeticOperator détermine si un opérateur est arithmétique
func isArithmeticOperator(operator string) bool {
	switch operator {
	case "+", "-", "*", "/", "%", "**", "^":
		return true
	default:
		return false
	}
}

// analyzeParenthesizedExpression analyse une expression parenthésée
func analyzeParenthesizedExpression(expr map[string]interface{}) (ExpressionType, error) {
	// Extraire l'expression interne
	var innerExpr interface{}
	var ok bool

	if innerExpr, ok = expr["expression"]; !ok {
		if innerExpr, ok = expr["expr"]; !ok {
			if innerExpr, ok = expr["inner"]; !ok {
				return ExprTypeSimple, fmt.Errorf("expression interne manquante dans l'expression parenthésée")
			}
		}
	}

	// Analyser récursivement l'expression interne
	return AnalyzeExpression(innerExpr)
}
