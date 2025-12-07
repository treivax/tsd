// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"regexp"
	"strings"
)

// evaluateArithmeticOperation évalue une opération arithmétique ou une concaténation de strings
func (e *AlphaConditionEvaluator) evaluateArithmeticOperation(left interface{}, operator string, right interface{}) (interface{}, error) {
	// Cas spécial pour l'opérateur + : si LES DEUX opérandes sont des strings, faire une concaténation
	if operator == "+" {
		leftVal := e.normalizeValue(left)
		rightVal := e.normalizeValue(right)

		leftStr, leftIsString := leftVal.(string)
		rightStr, rightIsString := rightVal.(string)

		// Si les deux sont des strings, concaténer
		if leftIsString && rightIsString {
			return leftStr + rightStr, nil
		}

		// Si un seul est une string, c'est une erreur - utiliser un cast explicite
		if leftIsString || rightIsString {
			return nil, fmt.Errorf("opération + avec types mixtes string/non-string (reçu: %T, %T). Utilisez un cast explicite: (string)valeur", leftVal, rightVal)
		}
	}

	// Pour tous les autres opérateurs (et + avec deux nombres), faire une opération arithmétique
	// Normaliser les valeurs numériques
	leftVal := e.normalizeValue(left)
	rightVal := e.normalizeValue(right)

	// Convertir en float64 pour les calculs
	leftNum, leftOk := leftVal.(float64)
	rightNum, rightOk := rightVal.(float64)

	if !leftOk || !rightOk {
		return nil, fmt.Errorf("opérations arithmétiques requièrent des valeurs numériques: gauche=%T, droite=%T", left, right)
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
		return nil, fmt.Errorf("opérateur arithmétique non supporté: %s", operator)
	}
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

	// D'abord remplacer les caractères LIKE par des placeholders temporaires
	tempPattern := strings.ReplaceAll(rightStr, "%", "PERCENTPLACEHOLDER")
	tempPattern = strings.ReplaceAll(tempPattern, "_", "UNDERSCOREPLACEHOLDER")

	// Échapper les caractères regex
	pattern := regexp.QuoteMeta(tempPattern)

	// Remplacer les placeholders par les équivalents regex
	pattern = strings.ReplaceAll(pattern, "PERCENTPLACEHOLDER", ".*")
	pattern = strings.ReplaceAll(pattern, "UNDERSCOREPLACEHOLDER", ".")
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
