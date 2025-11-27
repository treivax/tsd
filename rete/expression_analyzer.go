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

// CanDecompose détermine si une expression peut être décomposée en chaîne alpha
// Retourne true pour les expressions qui peuvent être transformées en séquence
// de nœuds alpha chaînés.
//
// Les expressions décomposables:
//   - ExprTypeSimple: condition atomique directement utilisable
//   - ExprTypeAND: peut être décomposée en chaîne de conditions (p1 AND p2 AND p3 -> alpha1->alpha2->alpha3)
//   - ExprTypeArithmetic: opérations arithmétiques commutatives peuvent être chaînées
//
// Les expressions non-décomposables:
//   - ExprTypeOR: nécessite des branches multiples (beta nodes ou duplication)
//   - ExprTypeMixed: nécessite normalisation en forme normale disjonctive (DNF) ou conjonctive (CNF)
func CanDecompose(exprType ExpressionType) bool {
	switch exprType {
	case ExprTypeSimple:
		return true
	case ExprTypeAND:
		return true
	case ExprTypeArithmetic:
		return true
	case ExprTypeNOT:
		return true // NOT peut être décomposé en alpha node avec négation
	case ExprTypeOR:
		return false // Nécessite traitement spécial avec branches
	case ExprTypeMixed:
		return false // Nécessite normalisation d'abord
	default:
		return false
	}
}

// ShouldNormalize détermine si une expression nécessite une normalisation
// avant de pouvoir être traitée dans le réseau RETE.
//
// La normalisation est nécessaire pour:
//   - ExprTypeMixed: doit être convertie en DNF ou CNF
//   - ExprTypeOR: peut bénéficier de la normalisation pour optimisation
//
// La normalisation n'est pas nécessaire pour:
//   - ExprTypeSimple: déjà sous forme atomique
//   - ExprTypeAND: déjà sous forme conjonctive, facilement chaînable
//   - ExprTypeArithmetic: structure déjà linéaire
func ShouldNormalize(exprType ExpressionType) bool {
	switch exprType {
	case ExprTypeSimple:
		return false
	case ExprTypeAND:
		return false
	case ExprTypeArithmetic:
		return false
	case ExprTypeNOT:
		return false // NOT peut être géré directement, mais peut bénéficier de la normalisation de l'expression interne
	case ExprTypeOR:
		return true // Bénéficie de la normalisation pour optimisation
	case ExprTypeMixed:
		return true // Doit être normalisée
	default:
		return false
	}
}

// GetExpressionComplexity retourne une estimation de la complexité d'une expression
// Utile pour décider de stratégies d'optimisation
func GetExpressionComplexity(exprType ExpressionType) int {
	switch exprType {
	case ExprTypeSimple:
		return 1
	case ExprTypeAND:
		return 2
	case ExprTypeArithmetic:
		return 2
	case ExprTypeNOT:
		return 2 // La négation ajoute une couche de complexité
	case ExprTypeOR:
		return 3
	case ExprTypeMixed:
		return 4
	default:
		return 0
	}
}

// RequiresBetaNode détermine si une expression nécessite des nœuds beta
// pour son évaluation (jointures, branches multiples, etc.)
func RequiresBetaNode(exprType ExpressionType) bool {
	switch exprType {
	case ExprTypeSimple:
		return false
	case ExprTypeAND:
		return false // Peut être géré avec alpha nodes chaînés
	case ExprTypeArithmetic:
		return false
	case ExprTypeNOT:
		return false // Peut être géré avec un alpha node avec flag de négation
	case ExprTypeOR:
		return true // Nécessite branches ou duplication
	case ExprTypeMixed:
		return true // Nécessite structure complexe
	default:
		return false
	}
}

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

// generateOptimizationHints génère des suggestions d'optimisation basées sur l'analyse
func generateOptimizationHints(expr interface{}, info *ExpressionInfo) []string {
	hints := make([]string, 0)

	// Hint pour De Morgan transformation
	if info.Type == ExprTypeNOT && info.InnerInfo != nil {
		if info.InnerInfo.Type == ExprTypeOR {
			hints = append(hints, "apply_demorgan_not_or")
		} else if info.InnerInfo.Type == ExprTypeAND {
			hints = append(hints, "apply_demorgan_not_and")
		} else if info.InnerInfo.Type == ExprTypeMixed {
			hints = append(hints, "push_negation_down")
		}
	}

	// Hint pour normalisation
	if info.ShouldNormalize {
		if info.Type == ExprTypeMixed {
			hints = append(hints, "normalize_to_dnf")
		} else if info.Type == ExprTypeOR {
			hints = append(hints, "consider_dnf_expansion")
		}
	}

	// Hint pour partage d'alpha nodes
	if info.Type == ExprTypeAND && info.Complexity >= 3 {
		hints = append(hints, "alpha_sharing_opportunity")
	}

	// Hint pour réordonnancement
	if info.Type == ExprTypeAND && canBenefitFromReordering(expr) {
		hints = append(hints, "consider_reordering")
	}

	// Hint pour expressions complexes
	if info.Complexity >= 4 {
		hints = append(hints, "high_complexity_review")
	}

	// Hint pour expressions nécessitant beta nodes
	if info.RequiresBeta {
		hints = append(hints, "requires_beta_node")
	}

	// Hint pour simplification arithmétique
	if info.Type == ExprTypeArithmetic {
		hints = append(hints, "consider_arithmetic_simplification")
	}

	return hints
}

// canBenefitFromReordering détermine si une expression AND peut bénéficier d'un réordonnancement
func canBenefitFromReordering(expr interface{}) bool {
	switch e := expr.(type) {
	case constraint.LogicalExpression:
		// Si l'expression a plusieurs opérations AND, le réordonnancement peut aider
		return len(e.Operations) >= 2

	case map[string]interface{}:
		if operations, ok := e["operations"].([]interface{}); ok {
			return len(operations) >= 2
		}
	}

	return false
}

// ShouldApplyDeMorgan détermine si la transformation de De Morgan devrait être appliquée
// basé sur des critères d'optimisation
func ShouldApplyDeMorgan(expr interface{}) bool {
	info, err := GetExpressionInfo(expr)
	if err != nil {
		return false
	}

	// Appliquer De Morgan si c'est un NOT(OR) ou NOT(AND)
	if info.Type == ExprTypeNOT && info.InnerInfo != nil {
		innerType := info.InnerInfo.Type
		if innerType == ExprTypeOR || innerType == ExprTypeAND {
			// Appliquer De Morgan seulement si cela réduit la complexité
			// NOT(A OR B) -> (NOT A) AND (NOT B) est bénéfique car AND est décomposable
			if innerType == ExprTypeOR {
				return true
			}
			// NOT(A AND B) -> (NOT A) OR (NOT B) peut être moins optimal
			// car OR nécessite des branches, donc on applique seulement si l'expression est simple
			if innerType == ExprTypeAND && info.InnerInfo.Complexity <= 2 {
				return true
			}
		}
	}

	return false
}
