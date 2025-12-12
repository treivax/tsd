// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

// evaluateArgument évalue un argument selon son type
func (ae *ActionExecutor) evaluateArgument(arg interface{}, ctx *ExecutionContext) (interface{}, error) {
	// Cas 1: Valeur littérale simple (string, number, bool)
	switch v := arg.(type) {
	case string, float64, bool, int, int64:
		return v, nil
	}

	// Cas 2: Map (objet structuré du parser)
	argMap, ok := arg.(map[string]interface{})
	if !ok {
		return arg, nil // Retourner tel quel si on ne peut pas le comprendre
	}

	argType, hasType := argMap["type"].(string)
	if !hasType {
		return arg, nil
	}

	switch argType {
	case "string", "number", "bool":
		// Valeur littérale typée
		if value, ok := argMap["value"]; ok {
			return value, nil
		}
		return arg, nil

	case "variable":
		// Cas 2: Fait complet référencé par variable
		varName, ok := argMap["name"].(string)
		if !ok {
			return nil, fmt.Errorf("nom de variable invalide")
		}
		fact := ctx.GetVariable(varName)
		if fact == nil {
			// Debug: list available variables
			availableVars := make([]string, 0)
			if ctx.varCache != nil {
				for k := range ctx.varCache {
					availableVars = append(availableVars, k)
				}
			}
			return nil, fmt.Errorf("variable '%s' non trouvée (variables disponibles: %v)", varName, availableVars)
		}
		return fact, nil

	case "fieldAccess":
		// Cas 3: Attribut de fait (variable.attribut)
		objectName, ok := argMap["object"].(string)
		if !ok {
			return nil, fmt.Errorf("nom d'objet invalide dans fieldAccess")
		}
		fieldName, ok := argMap["field"].(string)
		if !ok {
			return nil, fmt.Errorf("nom de champ invalide dans fieldAccess")
		}

		fact := ctx.GetVariable(objectName)
		if fact == nil {
			// Debug: list available variables
			availableVars := make([]string, 0)
			if ctx.varCache != nil {
				for k := range ctx.varCache {
					availableVars = append(availableVars, k)
				}
			}
			return nil, fmt.Errorf("variable '%s' non trouvée (variables disponibles: %v)", objectName, availableVars)
		}

		value, exists := fact.Fields[fieldName]
		if !exists {
			return nil, fmt.Errorf("champ '%s' non trouvé dans le fait %s", fieldName, objectName)
		}
		return value, nil

	case "factCreation":
		// Cas 4: Création de nouveau fait
		return ae.evaluateFactCreation(argMap, ctx)

	case "factModification":
		// Cas 5: Modification de fait
		return ae.evaluateFactModification(argMap, ctx)

	case "arithmetic":
		// Expression arithmétique (format legacy)
		return ae.evaluateArithmetic(argMap, ctx)

	case "binaryOperation", "binaryOp", "binary_operation":
		// Opération binaire (format du parser)
		return ae.evaluateBinaryOperation(argMap, ctx)

	case "cast":
		// Expression de cast
		return ae.evaluateCastExpression(argMap, ctx)

	default:
		return arg, nil
	}
}

// evaluateArithmetic évalue une expression arithmétique (format legacy)
func (ae *ActionExecutor) evaluateArithmetic(argMap map[string]interface{}, ctx *ExecutionContext) (interface{}, error) {
	operator, ok := argMap["operator"].(string)
	if !ok {
		return nil, fmt.Errorf("opérateur manquant dans arithmetic")
	}

	left, err := ae.evaluateArgument(argMap["left"], ctx)
	if err != nil {
		return nil, fmt.Errorf("erreur évaluation left: %w", err)
	}

	right, err := ae.evaluateArgument(argMap["right"], ctx)
	if err != nil {
		return nil, fmt.Errorf("erreur évaluation right: %w", err)
	}

	return ae.evaluateArithmeticOperation(left, operator, right)
}

// evaluateBinaryOperation évalue une opération binaire (format du parser)
func (ae *ActionExecutor) evaluateBinaryOperation(argMap map[string]interface{}, ctx *ExecutionContext) (interface{}, error) {
	// Extraire et normaliser l'opérateur en utilisant l'utilitaire centralisé
	operator, err := ExtractOperatorFromMap(argMap)
	if err != nil {
		return nil, fmt.Errorf("erreur extraction opérateur: %w", err)
	}

	left, err := ae.evaluateArgument(argMap["left"], ctx)
	if err != nil {
		return nil, fmt.Errorf("erreur évaluation left: %w", err)
	}

	right, err := ae.evaluateArgument(argMap["right"], ctx)
	if err != nil {
		return nil, fmt.Errorf("erreur évaluation right: %w", err)
	}

	// Distinguer les opérations arithmétiques des comparaisons
	switch operator {
	case "+", "-", "*", "/", "%":
		// Opération arithmétique - retourne une valeur numérique
		return ae.evaluateArithmeticOperation(left, operator, right)
	case "==", "!=", "<", "<=", ">", ">=":
		// Opération de comparaison - retourne un booléen
		// (utile si des actions ont besoin d'évaluer des booléens)
		return ae.evaluateComparison(left, operator, right)
	default:
		return nil, fmt.Errorf("opérateur binaire non supporté dans action: '%s'", operator)
	}
}

// evaluateArithmeticOperation effectue une opération arithmétique ou une concaténation de strings
func (ae *ActionExecutor) evaluateArithmeticOperation(left interface{}, operator string, right interface{}) (interface{}, error) {
	// Cas spécial pour l'opérateur + : si LES DEUX opérandes sont des strings, faire une concaténation
	if operator == "+" {
		leftStr, leftIsString := left.(string)
		rightStr, rightIsString := right.(string)

		// Si les deux sont des strings, concaténer
		if leftIsString && rightIsString {
			return leftStr + rightStr, nil
		}

		// Si un seul est une string, c'est une erreur - utiliser un cast explicite
		if leftIsString || rightIsString {
			return nil, fmt.Errorf("opération + avec types mixtes string/non-string (reçu: %T, %T). Utilisez un cast explicite: (string)valeur", left, right)
		}
	}

	// Pour tous les autres opérateurs (et + avec deux nombres), faire une opération arithmétique
	leftNum, okL := toNumber(left)
	rightNum, okR := toNumber(right)
	if !okL || !okR {
		return nil, fmt.Errorf("opération arithmétique nécessite des nombres (reçu: %T, %T)", left, right)
	}

	switch operator {
	case "+":
		return leftNum + rightNum, nil
	case "-":
		return leftNum - rightNum, nil
	case "*":
		return leftNum * rightNum, nil
	case "/":
		if rightNum == 0 {
			return nil, fmt.Errorf("division par zéro")
		}
		return leftNum / rightNum, nil
	case "%":
		if rightNum == 0 {
			return nil, fmt.Errorf("modulo par zéro")
		}
		return float64(int64(leftNum) % int64(rightNum)), nil
	default:
		return nil, fmt.Errorf("opérateur arithmétique inconnu: %s", operator)
	}
}

// evaluateComparison effectue une opération de comparaison
func (ae *ActionExecutor) evaluateComparison(left interface{}, operator string, right interface{}) (interface{}, error) {
	switch operator {
	case "==":
		return ae.areEqual(left, right), nil
	case "!=":
		return !ae.areEqual(left, right), nil
	case "<", "<=", ">", ">=":
		leftNum, okL := toNumber(left)
		rightNum, okR := toNumber(right)
		if !okL || !okR {
			return nil, fmt.Errorf("comparaison numérique nécessite des nombres (reçu: %T, %T)", left, right)
		}
		switch operator {
		case "<":
			return leftNum < rightNum, nil
		case "<=":
			return leftNum <= rightNum, nil
		case ">":
			return leftNum > rightNum, nil
		case ">=":
			return leftNum >= rightNum, nil
		}
	}
	return nil, fmt.Errorf("opérateur de comparaison inconnu: %s", operator)
}

// areEqual compare deux valeurs pour l'égalité
func (ae *ActionExecutor) areEqual(left, right interface{}) bool {
	// Normaliser les types numériques
	leftNum, leftIsNum := toNumber(left)
	rightNum, rightIsNum := toNumber(right)

	if leftIsNum && rightIsNum {
		return leftNum == rightNum
	}

	// Comparaison directe pour les autres types
	return left == right
}

// evaluateCastExpression évalue une expression de cast dans une action
func (ae *ActionExecutor) evaluateCastExpression(argMap map[string]interface{}, ctx *ExecutionContext) (interface{}, error) {
	// Extraire le type de cast
	castType, ok := argMap["castType"].(string)
	if !ok {
		return nil, fmt.Errorf("type de cast manquant ou invalide")
	}

	// Extraire l'expression à caster
	innerExpr, ok := argMap["expression"]
	if !ok {
		return nil, fmt.Errorf("expression à caster manquante")
	}

	// Évaluer l'expression interne
	value, err := ae.evaluateArgument(innerExpr, ctx)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'évaluation de l'expression à caster: %w", err)
	}

	// Appliquer le cast en utilisant les fonctions de rete/evaluator_cast.go
	result, err := EvaluateCast(castType, value)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du cast: %w", err)
	}

	return result, nil
}

// toNumber convertit une valeur en nombre flottant si possible
func toNumber(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case int32:
		return float64(n), true
	default:
		return 0, false
	}
}
