package rete

import (
	"fmt"

	"github.com/treivax/tsd/constraint"
)

// evaluateExpression évalue récursivement une expression
func (e *AlphaConditionEvaluator) evaluateExpression(expr interface{}) (bool, error) {
	switch condition := expr.(type) {
	case map[string]interface{}:
		return e.evaluateMapExpression(condition)
	case constraint.BinaryOperation:
		return e.evaluateBinaryOperation(condition)
	case constraint.LogicalExpression:
		return e.evaluateLogicalExpression(condition)
	case constraint.Constraint:
		return e.evaluateConstraint(condition)
	case constraint.BooleanLiteral:
		return condition.Value, nil
	default:
		return false, fmt.Errorf("type d'expression non supporté: %T", expr)
	}
}

// evaluateMapExpression évalue une expression sous forme de map (format JSON)
func (e *AlphaConditionEvaluator) evaluateMapExpression(expr map[string]interface{}) (bool, error) {
	exprType, ok := expr["type"].(string)
	if !ok {
		return false, fmt.Errorf("type d'expression manquant")
	}

	switch exprType {
	case "binaryOperation", "binary_op":
		return e.evaluateBinaryOperationMap(expr)
	case "logicalExpression", "logical_op", "logicalExpr":
		return e.evaluateLogicalExpressionMap(expr)
	case "constraint":
		return e.evaluateConstraintMap(expr)
	case "comparison":
		// Traitement des comparaisons directes
		return e.evaluateBinaryOperationMap(expr)
	case "negation":
		// Traitement spécial pour les contraintes de négation
		return e.evaluateNegationConstraint(expr)
	case "notConstraint":
		// Traitement spécial pour les contraintes NOT
		return e.evaluateNotConstraint(expr)
	case "existsConstraint":
		// Traitement spécial pour les contraintes EXISTS
		return e.evaluateExistsConstraint(expr)
	case "booleanLiteral":
		value, ok := expr["value"].(bool)
		if !ok {
			return false, fmt.Errorf("valeur booléenne invalide")
		}
		return value, nil
	case "simple":
		// Type simple: toujours vrai pour ce pipeline de base (contraintes simples sont filtrées par AlphaNodes)
		return true, nil
	default:
		return false, fmt.Errorf("type d'expression non supporté: %s", exprType)
	}
}

// evaluateBinaryOperation évalue une opération binaire
func (e *AlphaConditionEvaluator) evaluateBinaryOperation(op constraint.BinaryOperation) (bool, error) {
	left, err := e.evaluateValue(op.Left)
	if err != nil {
		return false, fmt.Errorf("erreur évaluation côté gauche: %w", err)
	}

	right, err := e.evaluateValue(op.Right)
	if err != nil {
		return false, fmt.Errorf("erreur évaluation côté droit: %w", err)
	}

	return e.compareValues(left, op.Operator, right)
}

// evaluateBinaryOperationMap évalue une opération binaire depuis une map
func (e *AlphaConditionEvaluator) evaluateBinaryOperationMap(expr map[string]interface{}) (bool, error) {
	// Supporter les deux formats: "operator" et "op"
	var operator string
	var ok bool

	if operator, ok = expr["operator"].(string); !ok {
		if operator, ok = expr["op"].(string); !ok {
			return false, fmt.Errorf("opérateur manquant (recherché 'operator' ou 'op')")
		}
	}

	// Debug: vérifier si left et right existent
	if expr["left"] == nil {
		return false, fmt.Errorf("côté gauche nil dans expr: %+v", expr)
	}
	if expr["right"] == nil {
		return false, fmt.Errorf("côté droit nil dans expr: %+v", expr)
	}

	left, err := e.evaluateValue(expr["left"])
	if err != nil {
		return false, fmt.Errorf("erreur évaluation côté gauche: %w", err)
	}

	right, err := e.evaluateValue(expr["right"])
	if err != nil {
		return false, fmt.Errorf("erreur évaluation côté droit: %w", err)
	}

	result, err := e.compareValues(left, operator, right)
	return result, err
}

// evaluateLogicalExpression évalue une expression logique (AND, OR)
func (e *AlphaConditionEvaluator) evaluateLogicalExpression(expr constraint.LogicalExpression) (bool, error) {
	leftResult, err := e.evaluateExpression(expr.Left)
	if err != nil {
		return false, fmt.Errorf("erreur évaluation côté gauche: %w", err)
	}

	result := leftResult
	for _, op := range expr.Operations {
		rightResult, err := e.evaluateExpression(op.Right)
		if err != nil {
			return false, fmt.Errorf("erreur évaluation opération %s: %w", op.Op, err)
		}

		switch op.Op {
		case "AND":
			result = result && rightResult
		case "OR":
			result = result || rightResult
		default:
			return false, fmt.Errorf("opérateur logique non supporté: %s", op.Op)
		}
	}

	return result, nil
}

// evaluateLogicalExpressionMap évalue une expression logique depuis une map
func (e *AlphaConditionEvaluator) evaluateLogicalExpressionMap(expr map[string]interface{}) (bool, error) {
	leftResult, err := e.evaluateExpression(expr["left"])
	if err != nil {
		return false, fmt.Errorf("erreur évaluation côté gauche: %w", err)
	}

	// Essayer d'extraire operations - supporter les deux types possibles
	var operations []interface{}
	var ok bool

	if opsRaw, hasOps := expr["operations"]; hasOps {
		// Essayer []interface{} d'abord
		operations, ok = opsRaw.([]interface{})
		if !ok {
			// Essayer []map[string]interface{} (structure retournée par parser PEG)
			if opsTyped, okTyped := opsRaw.([]map[string]interface{}); okTyped {
				// Convertir []map[string]interface{} en []interface{}
				operations = make([]interface{}, len(opsTyped))
				for i, op := range opsTyped {
					operations[i] = op
				}
				ok = true
			}
		}
	}

	if !ok {
		return leftResult, nil // Pas d'opérations supplémentaires
	}

	result := leftResult
	for _, opInterface := range operations {
		opMap, ok := opInterface.(map[string]interface{})
		if !ok {
			continue
		}

		operator, ok := opMap["op"].(string)
		if !ok {
			continue
		}

		rightResult, err := e.evaluateExpression(opMap["right"])
		if err != nil {
			return false, fmt.Errorf("erreur évaluation opération %s: %w", operator, err)
		}

		switch operator {
		case "AND":
			result = result && rightResult
		case "OR":
			result = result || rightResult
		default:
			return false, fmt.Errorf("opérateur logique non supporté: %s", operator)
		}
	}

	return result, nil
}
