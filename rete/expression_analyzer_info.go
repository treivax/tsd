// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package rete fournit l'implémentation du réseau RETE pour l'évaluation de règles.
// Ce fichier contient les structures et fonctions pour l'analyse détaillée des expressions,
// incluant l'extraction d'informations complètes et le calcul de complexité.
package rete

import (
	"fmt"

	"github.com/treivax/tsd/constraint"
)

// ExpressionInfo contient des informations détaillées sur une expression analysée
type ExpressionInfo struct {
	Type            ExpressionType
	CanDecompose    bool
	ShouldNormalize bool
	Complexity      int
	RequiresBeta    bool
	// InnerInfo contient l'analyse de l'expression interne pour les expressions
	// imbriquées (NOT, parenthésées, etc.)
	InnerInfo *ExpressionInfo
	// OptimizationHints contient des suggestions d'optimisation pour cette expression
	OptimizationHints []string
}

// GetExpressionInfo retourne des informations complètes sur une expression
// Pour les expressions NOT et parenthésées, analyse également l'expression interne
func GetExpressionInfo(expr interface{}) (*ExpressionInfo, error) {
	exprType, err := AnalyzeExpression(expr)
	if err != nil {
		return nil, err
	}

	// Calculer la complexité réelle en fonction du nombre d'opérations
	complexity := calculateActualComplexity(expr, exprType)

	info := &ExpressionInfo{
		Type:            exprType,
		CanDecompose:    CanDecompose(exprType),
		ShouldNormalize: ShouldNormalize(exprType),
		Complexity:      complexity,
		RequiresBeta:    RequiresBetaNode(exprType),
	}

	// Analyser récursivement l'expression interne pour les NOT
	if exprType == ExprTypeNOT {
		innerExpr := extractInnerExpression(expr)
		if innerExpr != nil {
			innerInfo, err := GetExpressionInfo(innerExpr)
			if err == nil {
				info.InnerInfo = innerInfo
				// Ajuster la complexité en fonction de l'expression interne
				info.Complexity = 2 + innerInfo.Complexity
			}
		}
	}

	// Générer les hints d'optimisation
	info.OptimizationHints = generateOptimizationHints(expr, info)

	return info, nil
}

// extractInnerExpression extrait l'expression interne d'un NOT ou d'une expression parenthésée
func extractInnerExpression(expr interface{}) interface{} {
	switch e := expr.(type) {
	case constraint.NotConstraint:
		return e.Expression

	case map[string]interface{}:
		exprType, ok := e["type"].(string)
		if !ok {
			return nil
		}

		switch exprType {
		case "notConstraint", "not", "negation":
			// Extraire la contrainte interne
			if inner, ok := e["constraint"]; ok {
				return inner
			}
			if inner, ok := e["expr"]; ok {
				return inner
			}
			if inner, ok := e["expression"]; ok {
				return inner
			}

		case "parenthesized", "parenthesizedExpression", "group":
			// Extraire l'expression à l'intérieur des parenthèses
			if inner, ok := e["expression"]; ok {
				return inner
			}
			if inner, ok := e["expr"]; ok {
				return inner
			}
			if inner, ok := e["inner"]; ok {
				return inner
			}
		}
	}

	return nil
}

// AnalyzeInnerExpression analyse l'expression interne d'un NOT ou d'une expression parenthésée
// Retourne le type de l'expression interne et une erreur si l'analyse échoue
func AnalyzeInnerExpression(expr interface{}) (ExpressionType, error) {
	innerExpr := extractInnerExpression(expr)
	if innerExpr == nil {
		return ExprTypeSimple, fmt.Errorf("impossible d'extraire l'expression interne")
	}
	return AnalyzeExpression(innerExpr)
}

// calculateActualComplexity calcule la complexité réelle en fonction du nombre d'opérations
func calculateActualComplexity(expr interface{}, exprType ExpressionType) int {
	baseComplexity := GetExpressionComplexity(exprType)

	// Pour les expressions logiques, calculer en fonction du nombre d'opérations
	switch e := expr.(type) {
	case constraint.LogicalExpression:
		if exprType == ExprTypeAND || exprType == ExprTypeOR || exprType == ExprTypeMixed {
			// 1 pour le terme de gauche + 1 pour chaque opération
			return 1 + len(e.Operations)
		}
	case map[string]interface{}:
		if operations, ok := e["operations"].([]interface{}); ok {
			if exprType == ExprTypeAND || exprType == ExprTypeOR || exprType == ExprTypeMixed {
				return 1 + len(operations)
			}
		}
	}

	return baseComplexity
}
