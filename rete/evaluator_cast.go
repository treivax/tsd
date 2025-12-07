// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strconv"
	"strings"
)

// EvaluateCast évalue une expression de cast et retourne la valeur convertie
func EvaluateCast(castType string, value interface{}) (interface{}, error) {
	switch castType {
	case "number":
		return CastToNumber(value)
	case "string":
		return CastToString(value)
	case "bool":
		return CastToBool(value)
	default:
		return nil, fmt.Errorf("type de cast non supporté: %s", castType)
	}
}

// CastToNumber convertit une valeur en nombre (float64)
func CastToNumber(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		// Déjà un nombre
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		// Nettoyer les espaces
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			return 0, fmt.Errorf("cannot cast empty string to number")
		}
		num, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot cast '%s' to number: %w", v, err)
		}
		return num, nil
	case bool:
		// true -> 1, false -> 0
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot cast type %T to number", value)
	}
}

// CastToString convertit une valeur en chaîne de caractères
func CastToString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		// Déjà une chaîne
		return v, nil
	case float64:
		// Convertir le nombre en chaîne
		// Si c'est un entier, ne pas afficher de décimales
		if v == float64(int64(v)) {
			return strconv.FormatInt(int64(v), 10), nil
		}
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case bool:
		// true -> "true", false -> "false"
		if v {
			return "true", nil
		}
		return "false", nil
	default:
		return "", fmt.Errorf("cannot cast type %T to string", value)
	}
}

// CastToBool convertit une valeur en booléen
func CastToBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		// Déjà un booléen
		return v, nil
	case string:
		// Parse la chaîne en booléen
		trimmed := strings.TrimSpace(v)
		lower := strings.ToLower(trimmed)

		// Valeurs vraies
		if lower == "true" || lower == "1" {
			return true, nil
		}

		// Valeurs fausses
		if lower == "false" || lower == "0" || lower == "" {
			return false, nil
		}

		// Comportement permissif : toute autre chaîne est considérée comme false
		return false, nil
	case float64:
		// 0 -> false, non-zéro -> true
		return v != 0, nil
	case int:
		return v != 0, nil
	case int64:
		return v != 0, nil
	default:
		return false, fmt.Errorf("cannot cast type %T to bool", value)
	}
}

// evaluateCastExpression évalue une expression de cast depuis une map
func (e *AlphaConditionEvaluator) evaluateCastExpression(expr map[string]interface{}) (interface{}, error) {
	// Extraire le type de cast
	castType, ok := expr["castType"].(string)
	if !ok {
		return nil, fmt.Errorf("type de cast manquant ou invalide")
	}

	// Extraire l'expression à caster
	innerExpr, ok := expr["expression"]
	if !ok {
		return nil, fmt.Errorf("expression à caster manquante")
	}

	// Évaluer l'expression interne
	value, err := e.evaluateValue(innerExpr)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'évaluation de l'expression à caster: %w", err)
	}

	// Appliquer le cast
	result, err := EvaluateCast(castType, value)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du cast: %w", err)
	}

	return result, nil
}
