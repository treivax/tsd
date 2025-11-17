package rete

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
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
	case "comparison":
		// Traitement des comparaisons directes
		return e.evaluateBinaryOperationMap(expr)
	case "negation":
		// Traitement spécial pour les contraintes de négation
		return e.evaluateNegationConstraint(expr)
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
	// Si l'expression a une clé "constraint", extraire la contrainte réelle
	var actualConstraint map[string]interface{}
	if constraintData, hasConstraint := expr["constraint"]; hasConstraint {
		if constraintMap, ok := constraintData.(map[string]interface{}); ok {
			actualConstraint = constraintMap
		} else {
			return false, fmt.Errorf("format contrainte invalide: %T", constraintData)
		}
	} else {
		// Utiliser directement l'expression si pas d'indirection
		actualConstraint = expr
	}

	operator, ok := actualConstraint["operator"].(string)
	if !ok {
		return false, fmt.Errorf("opérateur manquant")
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

	case "numberLiteral", "number":
		value, ok := val["value"].(float64)
		if !ok {
			// Essayer aussi avec int
			if intValue, ok := val["value"].(int); ok {
				return float64(intValue), nil
			}
			return nil, fmt.Errorf("valeur numérique invalide")
		}
		return value, nil

	case "stringLiteral", "string":
		value, ok := val["value"].(string)
		if !ok {
			return nil, fmt.Errorf("valeur de chaîne invalide")
		}
		return value, nil

	case "booleanLiteral", "boolean":
		value, ok := val["value"].(bool)
		if !ok {
			return nil, fmt.Errorf("valeur booléenne invalide")
		}
		return value, nil

	case "functionCall", "function_call":
		// Support des appels de fonction
		return e.evaluateFunctionCall(val)

	case "arrayLiteral", "array_literal":
		// Support des littéraux de tableau
		elements, ok := val["elements"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("éléments de tableau invalides")
		}

		// Évaluer chaque élément du tableau
		evaluatedElements := make([]interface{}, len(elements))
		for i, element := range elements {
			evaluatedElement, err := e.evaluateValue(element)
			if err != nil {
				return nil, fmt.Errorf("erreur évaluation élément tableau[%d]: %w", i, err)
			}
			evaluatedElements[i] = evaluatedElement
		}
		return evaluatedElements, nil

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
		// Debug info pour aider à diagnostiquer les problèmes de binding
		availableVars := make([]string, 0, len(e.variableBindings))
		for k := range e.variableBindings {
			availableVars = append(availableVars, k)
		}
		return nil, fmt.Errorf("variable non liée: %s (variables disponibles: %v)", object, availableVars)
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
	case "CONTAINS":
		return e.evaluateContains(leftVal, rightVal)
	case "IN":
		return e.evaluateIn(leftVal, rightVal)
	case "LIKE":
		return e.evaluateLike(leftVal, rightVal)
	case "MATCHES":
		return e.evaluateMatches(leftVal, rightVal)
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

// evaluateContains vérifie si une chaîne contient une sous-chaîne
func (e *AlphaConditionEvaluator) evaluateContains(left, right interface{}) (bool, error) {
	leftStr, ok := left.(string)
	if !ok {
		return false, fmt.Errorf("l'opérateur CONTAINS nécessite une chaîne à gauche, reçu: %T", left)
	}

	rightStr, ok := right.(string)
	if !ok {
		return false, fmt.Errorf("l'opérateur CONTAINS nécessite une chaîne à droite, reçu: %T", right)
	}

	return strings.Contains(leftStr, rightStr), nil
}

// evaluateIn vérifie si une valeur fait partie d'un tableau
func (e *AlphaConditionEvaluator) evaluateIn(left, right interface{}) (bool, error) {
	// Convertir le côté droit en slice
	var rightSlice []interface{}

	switch rightVal := right.(type) {
	case []interface{}:
		rightSlice = rightVal
	case []string:
		rightSlice = make([]interface{}, len(rightVal))
		for i, v := range rightVal {
			rightSlice[i] = v
		}
	case []int:
		rightSlice = make([]interface{}, len(rightVal))
		for i, v := range rightVal {
			rightSlice[i] = v
		}
	case []float64:
		rightSlice = make([]interface{}, len(rightVal))
		for i, v := range rightVal {
			rightSlice[i] = v
		}
	default:
		return false, fmt.Errorf("l'opérateur IN nécessite un tableau à droite, reçu: %T", right)
	}

	// Vérifier si la valeur de gauche existe dans le tableau
	for _, item := range rightSlice {
		if e.areEqual(left, item) {
			return true, nil
		}
	}

	return false, nil
}

// evaluateLike vérifie si une chaîne correspond à un pattern (SQL LIKE style)
func (e *AlphaConditionEvaluator) evaluateLike(left, right interface{}) (bool, error) {
	leftStr, ok := left.(string)
	if !ok {
		return false, fmt.Errorf("l'opérateur LIKE nécessite une chaîne à gauche, reçu: %T", left)
	}

	rightStr, ok := right.(string)
	if !ok {
		return false, fmt.Errorf("l'opérateur LIKE nécessite un pattern à droite, reçu: %T", right)
	}

	// Convertir pattern SQL LIKE en regex Go
	// % = .* (zéro ou plus de caractères)
	// _ = . (exactement un caractère)
	pattern := regexp.QuoteMeta(rightStr)
	pattern = strings.ReplaceAll(pattern, "\\%", ".*")
	pattern = strings.ReplaceAll(pattern, "\\_", ".")
	pattern = "^" + pattern + "$"

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Errorf("pattern LIKE invalide '%s': %w", rightStr, err)
	}

	return regex.MatchString(leftStr), nil
}

// evaluateMatches vérifie si une chaîne correspond à une expression régulière
func (e *AlphaConditionEvaluator) evaluateMatches(left, right interface{}) (bool, error) {
	leftStr, ok := left.(string)
	if !ok {
		return false, fmt.Errorf("l'opérateur MATCHES nécessite une chaîne à gauche, reçu: %T", left)
	}

	rightStr, ok := right.(string)
	if !ok {
		return false, fmt.Errorf("l'opérateur MATCHES nécessite un pattern regex à droite, reçu: %T", right)
	}

	regex, err := regexp.Compile(rightStr)
	if err != nil {
		return false, fmt.Errorf("pattern regex invalide '%s': %w", rightStr, err)
	}

	return regex.MatchString(leftStr), nil
}

// evaluateFunctionCall évalue un appel de fonction
func (e *AlphaConditionEvaluator) evaluateFunctionCall(val map[string]interface{}) (interface{}, error) {
	functionName, ok := val["name"].(string)
	if !ok {
		return nil, fmt.Errorf("nom de fonction invalide")
	}

	args, ok := val["args"].([]interface{})
	if !ok {
		// Pas d'arguments
		args = []interface{}{}
	}

	// Évaluer les arguments
	evaluatedArgs := make([]interface{}, len(args))
	for i, arg := range args {
		evaluatedArg, err := e.evaluateValue(arg)
		if err != nil {
			return nil, fmt.Errorf("erreur évaluation argument[%d] pour %s: %w", i, functionName, err)
		}
		evaluatedArgs[i] = evaluatedArg
	}

	// Appeler la fonction appropriée
	switch functionName {
	case "LENGTH":
		return e.evaluateLength(evaluatedArgs)
	case "UPPER":
		return e.evaluateUpper(evaluatedArgs)
	case "LOWER":
		return e.evaluateLower(evaluatedArgs)
	case "ABS":
		return e.evaluateAbs(evaluatedArgs)
	case "ROUND":
		return e.evaluateRound(evaluatedArgs)
	case "SUBSTRING":
		return e.evaluateSubstring(evaluatedArgs)
	case "TRIM":
		return e.evaluateTrim(evaluatedArgs)
	default:
		return nil, fmt.Errorf("fonction non supportée: %s", functionName)
	}
}

// evaluateLength retourne la longueur d'une chaîne
func (e *AlphaConditionEvaluator) evaluateLength(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("LENGTH() attend 1 argument, reçu %d", len(args))
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("LENGTH() nécessite une chaîne, reçu: %T", args[0])
	}

	return float64(len(str)), nil
}

// evaluateUpper convertit une chaîne en majuscules
func (e *AlphaConditionEvaluator) evaluateUpper(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("UPPER() attend 1 argument, reçu %d", len(args))
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("UPPER() nécessite une chaîne, reçu: %T", args[0])
	}

	return strings.ToUpper(str), nil
}

// evaluateLower convertit une chaîne en minuscules
func (e *AlphaConditionEvaluator) evaluateLower(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("LOWER() attend 1 argument, reçu %d", len(args))
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("LOWER() nécessite une chaîne, reçu: %T", args[0])
	}

	return strings.ToLower(str), nil
}

// evaluateAbs retourne la valeur absolue d'un nombre
func (e *AlphaConditionEvaluator) evaluateAbs(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ABS() attend 1 argument, reçu %d", len(args))
	}

	num, ok := args[0].(float64)
	if !ok {
		return nil, fmt.Errorf("ABS() nécessite un nombre, reçu: %T", args[0])
	}

	return math.Abs(num), nil
}

// evaluateRound arrondit un nombre
func (e *AlphaConditionEvaluator) evaluateRound(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ROUND() attend 1 argument, reçu %d", len(args))
	}

	num, ok := args[0].(float64)
	if !ok {
		return nil, fmt.Errorf("ROUND() nécessite un nombre, reçu: %T", args[0])
	}

	return math.Round(num), nil
}

// evaluateSubstring extrait une sous-chaîne
func (e *AlphaConditionEvaluator) evaluateSubstring(args []interface{}) (interface{}, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("SUBSTRING() attend 2 ou 3 arguments, reçu %d", len(args))
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("SUBSTRING() nécessite une chaîne comme premier argument, reçu: %T", args[0])
	}

	start, ok := args[1].(float64)
	if !ok {
		return nil, fmt.Errorf("SUBSTRING() nécessite un nombre comme deuxième argument, reçu: %T", args[1])
	}

	startInt := int(start)
	if startInt < 0 || startInt >= len(str) {
		return "", nil // Retourner chaîne vide si index hors limites
	}

	if len(args) == 3 {
		length, ok := args[2].(float64)
		if !ok {
			return nil, fmt.Errorf("SUBSTRING() nécessite un nombre comme troisième argument, reçu: %T", args[2])
		}

		lengthInt := int(length)
		endInt := startInt + lengthInt
		if endInt > len(str) {
			endInt = len(str)
		}
		return str[startInt:endInt], nil
	}

	return str[startInt:], nil
}

// evaluateTrim supprime les espaces en début et fin
func (e *AlphaConditionEvaluator) evaluateTrim(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("TRIM() attend 1 argument, reçu %d", len(args))
	}

	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("TRIM() nécessite une chaîne, reçu: %T", args[0])
	}

	return strings.TrimSpace(str), nil
}
