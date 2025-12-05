// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
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
