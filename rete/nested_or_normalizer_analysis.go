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
	updateNestingDepth(analysis, depth)

	hasOR, hasAND := scanLogicalOperations(expr, analysis, depth)
	analyzeLogicalLeftSide(expr, analysis, depth)
	updateComplexityBasedOnOperators(analysis, hasOR, hasAND, depth)
}

// scanLogicalOperations parcourt les opérations logiques et met à jour les compteurs
func scanLogicalOperations(expr constraint.LogicalExpression, analysis *NestedORAnalysis, depth int) (hasOR bool, hasAND bool) {
	for _, op := range expr.Operations {
		opStr := strings.ToUpper(op.Op)

		if opStr == "OR" || opStr == "||" {
			hasOR = true
			analysis.ORTermCount++
		} else if opStr == "AND" || opStr == "&&" {
			hasAND = true
			analysis.ANDTermCount++
		}

		analyzeLogicalRightSide(op, analysis, depth)
	}

	return hasOR, hasAND
}

// analyzeLogicalRightSide analyse récursivement le côté droit d'une opération logique
func analyzeLogicalRightSide(op constraint.LogicalOperation, analysis *NestedORAnalysis, depth int) {
	if rightLogical, ok := op.Right.(constraint.LogicalExpression); ok {
		analyzeLogicalExpressionNesting(rightLogical, analysis, depth+1)
	}
}

// analyzeLogicalLeftSide analyse le côté gauche de l'expression logique
func analyzeLogicalLeftSide(expr constraint.LogicalExpression, analysis *NestedORAnalysis, depth int) {
	if leftLogical, ok := expr.Left.(constraint.LogicalExpression); ok {
		analyzeLogicalExpressionNesting(leftLogical, analysis, depth+1)
	}
}

// analyzeMapExpressionNesting analyse l'imbrication d'une expression map
func analyzeMapExpressionNesting(expr map[string]interface{}, analysis *NestedORAnalysis, depth int) {
	updateNestingDepth(analysis, depth)

	hasOR, hasAND := scanMapOperationsForNesting(expr, analysis, depth)
	analyzeMapLeftSide(expr, analysis, depth)
	updateComplexityBasedOnOperators(analysis, hasOR, hasAND, depth)
}

// updateNestingDepth met à jour la profondeur d'imbrication maximale
func updateNestingDepth(analysis *NestedORAnalysis, depth int) {
	if depth > analysis.NestingDepth {
		analysis.NestingDepth = depth
	}
}

// scanMapOperationsForNesting parcourt les opérations et met à jour les compteurs
func scanMapOperationsForNesting(expr map[string]interface{}, analysis *NestedORAnalysis, depth int) (hasOR bool, hasAND bool) {
	operations, ok := expr["operations"]
	if !ok {
		return false, false
	}

	opsList, ok := operations.([]interface{})
	if !ok {
		return false, false
	}

	for _, opInterface := range opsList {
		opMap, ok := opInterface.(map[string]interface{})
		if !ok {
			continue
		}

		op, ok := opMap["op"].(string)
		if !ok {
			continue
		}

		opStr := strings.ToUpper(op)
		if opStr == "OR" || opStr == "||" {
			hasOR = true
			analysis.ORTermCount++
		} else if opStr == "AND" || opStr == "&&" {
			hasAND = true
			analysis.ANDTermCount++
		}

		analyzeRightSideRecursively(opMap, analysis, depth)
	}

	return hasOR, hasAND
}

// analyzeRightSideRecursively analyse récursivement le côté droit d'une opération
func analyzeRightSideRecursively(opMap map[string]interface{}, analysis *NestedORAnalysis, depth int) {
	right, ok := opMap["right"].(map[string]interface{})
	if !ok {
		return
	}

	rightType, ok := right["type"].(string)
	if ok && rightType == "logicalExpr" {
		analyzeMapExpressionNesting(right, analysis, depth+1)
	}
}

// analyzeMapLeftSide analyse le côté gauche de l'expression
func analyzeMapLeftSide(expr map[string]interface{}, analysis *NestedORAnalysis, depth int) {
	left, ok := expr["left"].(map[string]interface{})
	if !ok {
		return
	}

	leftType, ok := left["type"].(string)
	if ok && leftType == "logicalExpr" {
		analyzeMapExpressionNesting(left, analysis, depth+1)
	}
}

// updateComplexityBasedOnOperators détermine et met à jour la complexité
func updateComplexityBasedOnOperators(analysis *NestedORAnalysis, hasOR bool, hasAND bool, depth int) {
	if hasOR && hasAND {
		updateMixedComplexity(analysis)
	} else if hasOR {
		updateORComplexity(analysis, depth)
	}
}

// updateMixedComplexity met à jour la complexité pour les expressions mixtes AND/OR
func updateMixedComplexity(analysis *NestedORAnalysis) {
	if analysis.ORTermCount >= 2 && analysis.ANDTermCount >= 1 {
		if analysis.Complexity < ComplexityDNFCandidate {
			analysis.Complexity = ComplexityDNFCandidate
		}
	} else if analysis.Complexity < ComplexityMixedANDOR {
		analysis.Complexity = ComplexityMixedANDOR
	}
}

// updateORComplexity met à jour la complexité pour les expressions avec OR
func updateORComplexity(analysis *NestedORAnalysis, depth int) {
	if depth > 0 {
		if analysis.Complexity < ComplexityNestedOR {
			analysis.Complexity = ComplexityNestedOR
		}
	} else if analysis.Complexity < ComplexityFlat {
		analysis.Complexity = ComplexityFlat
	}
}
