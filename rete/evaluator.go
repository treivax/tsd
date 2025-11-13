package rete

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/treivax/tsd/constraint"
)

// AlphaConditionEvaluator évalue les conditions Alpha sur les faits
type AlphaConditionEvaluator struct {
	variableBindings map[string]*Fact
}

// NewAlphaConditionEvaluator crée un nouvel évaluateur de conditions
func NewAlphaConditionEvaluator() *AlphaConditionEvaluator {
	return &AlphaConditionEvaluator{
		variableBindings: make(map[string]*Fact),
	}
}

// EvaluateCondition évalue une condition sur un fait
func (e *AlphaConditionEvaluator) EvaluateCondition(condition interface{}, fact *Fact, variableName string) (bool, error) {
	// Lier la variable au fait pour l'évaluation
	if variableName != "" {
		e.variableBindings[variableName] = fact
	}

	return e.evaluateExpression(condition)
}

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
	case "logicalExpression", "logical_op":
		return e.evaluateLogicalExpressionMap(expr)
	case "constraint":
		return e.evaluateConstraintMap(expr)
	case "booleanLiteral":
		value, ok := expr["value"].(bool)
		if !ok {
			return false, fmt.Errorf("valeur booléenne invalide")
		}
		return value, nil
	case "simple":
		// Type simple: toujours vrai pour ce pipeline de base
		// TODO: Implémenter l'évaluation réelle des contraintes simples
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

	left, err := e.evaluateValue(expr["left"])
	if err != nil {
		return false, fmt.Errorf("erreur évaluation côté gauche: %w", err)
	}

	right, err := e.evaluateValue(expr["right"])
	if err != nil {
		return false, fmt.Errorf("erreur évaluation côté droit: %w", err)
	}

	return e.compareValues(left, operator, right)
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

	operations, ok := expr["operations"].([]interface{})
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
	operator, ok := expr["operator"].(string)
	if !ok {
		return false, fmt.Errorf("opérateur manquant")
	}

	left, err := e.evaluateValue(expr["left"])
	if err != nil {
		return false, fmt.Errorf("erreur évaluation côté gauche: %w", err)
	}

	right, err := e.evaluateValue(expr["right"])
	if err != nil {
		return false, fmt.Errorf("erreur évaluation côté droit: %w", err)
	}

	return e.compareValues(left, operator, right)
}

// evaluateValue évalue une valeur (littéral, accès de champ, variable)
func (e *AlphaConditionEvaluator) evaluateValue(value interface{}) (interface{}, error) {
	switch val := value.(type) {
	case map[string]interface{}:
		return e.evaluateValueFromMap(val)
	case constraint.FieldAccess:
		return e.evaluateFieldAccess(val)
	case constraint.Variable:
		return e.evaluateVariable(val)
	case constraint.NumberLiteral:
		return val.Value, nil
	case constraint.StringLiteral:
		return val.Value, nil
	case constraint.BooleanLiteral:
		return val.Value, nil
	case string:
		return val, nil
	case int, int32, int64:
		return val, nil
	case float32, float64:
		return val, nil
	case bool:
		return val, nil
	default:
		return nil, fmt.Errorf("type de valeur non supporté: %T", value)
	}
}

// evaluateValueFromMap évalue une valeur depuis une map
func (e *AlphaConditionEvaluator) evaluateValueFromMap(val map[string]interface{}) (interface{}, error) {
	valType, ok := val["type"].(string)
	if !ok {
		return nil, fmt.Errorf("type de valeur manquant")
	}

	switch valType {
	case "fieldAccess", "field_access":
		// Supporter les deux formats: object/field et variable/field
		var objectOrVariable, field string
		var ok bool

		if objectOrVariable, ok = val["object"].(string); ok {
			// Format: object + field
		} else if objectOrVariable, ok = val["variable"].(string); ok {
			// Format: variable + field
		} else {
			return nil, fmt.Errorf("objet ou variable d'accès de champ invalide")
		}

		if field, ok = val["field"].(string); !ok {
			return nil, fmt.Errorf("champ d'accès invalide")
		}
		return e.evaluateFieldAccessByName(objectOrVariable, field)

	case "variable":
		name, ok := val["name"].(string)
		if !ok {
			return nil, fmt.Errorf("nom de variable invalide")
		}
		return e.evaluateVariableByName(name)

	case "numberLiteral":
		value, ok := val["value"].(float64)
		if !ok {
			return nil, fmt.Errorf("valeur numérique invalide")
		}
		return value, nil

	case "stringLiteral":
		value, ok := val["value"].(string)
		if !ok {
			return nil, fmt.Errorf("valeur de chaîne invalide")
		}
		return value, nil

	case "booleanLiteral":
		value, ok := val["value"].(bool)
		if !ok {
			return nil, fmt.Errorf("valeur booléenne invalide")
		}
		return value, nil

	default:
		return nil, fmt.Errorf("type de valeur non supporté: %s", valType)
	}
}

// evaluateFieldAccess évalue l'accès à un champ d'une variable
func (e *AlphaConditionEvaluator) evaluateFieldAccess(fa constraint.FieldAccess) (interface{}, error) {
	return e.evaluateFieldAccessByName(fa.Object, fa.Field)
}

// evaluateFieldAccessByName évalue l'accès à un champ par nom
func (e *AlphaConditionEvaluator) evaluateFieldAccessByName(object, field string) (interface{}, error) {
	fact, exists := e.variableBindings[object]
	if !exists {
		return nil, fmt.Errorf("variable non liée: %s", object)
	}

	value, exists := fact.Fields[field]
	if !exists {
		return nil, fmt.Errorf("champ inexistant: %s.%s", object, field)
	}

	return value, nil
}

// evaluateVariable évalue une variable
func (e *AlphaConditionEvaluator) evaluateVariable(variable constraint.Variable) (interface{}, error) {
	return e.evaluateVariableByName(variable.Name)
}

// evaluateVariableByName évalue une variable par nom
func (e *AlphaConditionEvaluator) evaluateVariableByName(name string) (interface{}, error) {
	fact, exists := e.variableBindings[name]
	if !exists {
		return nil, fmt.Errorf("variable non liée: %s", name)
	}

	// Retourner le fait entier ou une représentation
	return fact, nil
}

// compareValues compare deux valeurs avec un opérateur
func (e *AlphaConditionEvaluator) compareValues(left interface{}, operator string, right interface{}) (bool, error) {
	// Normaliser les valeurs numériques
	leftVal := e.normalizeValue(left)
	rightVal := e.normalizeValue(right)

	switch operator {
	case "==", "=":
		return e.areEqual(leftVal, rightVal), nil
	case "!=", "<>":
		return !e.areEqual(leftVal, rightVal), nil
	case "<":
		return e.isLess(leftVal, rightVal)
	case "<=":
		equal := e.areEqual(leftVal, rightVal)
		less, err := e.isLess(leftVal, rightVal)
		return equal || less, err
	case ">":
		return e.isGreater(leftVal, rightVal)
	case ">=":
		equal := e.areEqual(leftVal, rightVal)
		greater, err := e.isGreater(leftVal, rightVal)
		return equal || greater, err
	default:
		return false, fmt.Errorf("opérateur non supporté: %s", operator)
	}
}

// normalizeValue normalise une valeur pour la comparaison
func (e *AlphaConditionEvaluator) normalizeValue(value interface{}) interface{} {
	switch v := value.(type) {
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case float32:
		return float64(v)
	default:
		return value
	}
}

// areEqual vérifie si deux valeurs sont égales
func (e *AlphaConditionEvaluator) areEqual(left, right interface{}) bool {
	return reflect.DeepEqual(left, right)
}

// isLess vérifie si left < right
func (e *AlphaConditionEvaluator) isLess(left, right interface{}) (bool, error) {
	switch leftVal := left.(type) {
	case float64:
		if rightVal, ok := right.(float64); ok {
			return leftVal < rightVal, nil
		}
	case string:
		if rightVal, ok := right.(string); ok {
			return strings.Compare(leftVal, rightVal) < 0, nil
		}
	}
	return false, fmt.Errorf("impossible de comparer %T avec %T", left, right)
}

// isGreater vérifie si left > right
func (e *AlphaConditionEvaluator) isGreater(left, right interface{}) (bool, error) {
	switch leftVal := left.(type) {
	case float64:
		if rightVal, ok := right.(float64); ok {
			return leftVal > rightVal, nil
		}
	case string:
		if rightVal, ok := right.(string); ok {
			return strings.Compare(leftVal, rightVal) > 0, nil
		}
	}
	return false, fmt.Errorf("impossible de comparer %T avec %T", left, right)
}

// ClearBindings efface les liaisons de variables
func (e *AlphaConditionEvaluator) ClearBindings() {
	e.variableBindings = make(map[string]*Fact)
}

// GetBindings retourne les liaisons actuelles
func (e *AlphaConditionEvaluator) GetBindings() map[string]*Fact {
	return e.variableBindings
}
