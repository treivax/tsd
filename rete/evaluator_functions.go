package rete

import (
	"fmt"
	"math"
	"strings"
)

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

	// Debug: vérifier si les arguments contiennent des nil
	for i, arg := range args {
		if arg == nil {
			return nil, fmt.Errorf("argument[%d] de %s est nil, args complets: %+v", i, functionName, args)
		}
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
	case "FLOOR":
		return e.evaluateFloor(evaluatedArgs)
	case "CEIL":
		return e.evaluateCeil(evaluatedArgs)
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

// evaluateFloor arrondit un nombre vers le bas
func (e *AlphaConditionEvaluator) evaluateFloor(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("FLOOR() attend 1 argument, reçu %d", len(args))
	}

	num, ok := args[0].(float64)
	if !ok {
		return nil, fmt.Errorf("FLOOR() nécessite un nombre, reçu: %T", args[0])
	}

	return math.Floor(num), nil
}

// evaluateCeil arrondit un nombre vers le haut
func (e *AlphaConditionEvaluator) evaluateCeil(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("CEIL() attend 1 argument, reçu %d", len(args))
	}

	num, ok := args[0].(float64)
	if !ok {
		return nil, fmt.Errorf("CEIL() nécessite un nombre, reçu: %T", args[0])
	}

	return math.Ceil(num), nil
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
