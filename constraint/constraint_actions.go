// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
)

// ValidateAction vérifie qu'une action est valide dans le contexte d'une expression
func ValidateAction(program Program, action Action, expressionIndex int) error {
	if expressionIndex >= len(program.Expressions) {
		return fmt.Errorf("index d'expression invalide: %d", expressionIndex)
	}

	expression := program.Expressions[expressionIndex]

	// Créer une map des variables disponibles dans l'expression
	availableVars := make(map[string]bool)

	// Ajouter les variables du Set principal (ancien format, rétrocompatibilité)
	for _, variable := range expression.Set.Variables {
		availableVars[variable.Name] = true
	}

	// Ajouter les variables des Patterns multiples (nouveau format avec agrégation)
	for _, pattern := range expression.Patterns {
		for _, variable := range pattern.Variables {
			availableVars[variable.Name] = true
		}
	}

	// Obtenir tous les jobs (supporte ancien et nouveau format)
	jobs := action.GetJobs()

	// Vérifier que tous les arguments de chaque job référencent des variables valides
	for _, job := range jobs {
		for _, arg := range job.Args {
			// Extraire les variables utilisées dans l'argument
			vars := extractVariablesFromArg(arg)
			for _, varName := range vars {
				if !availableVars[varName] {
					return fmt.Errorf("action %s: argument contient la variable '%s' qui ne correspond à aucune variable de l'expression", job.Name, varName)
				}
			}
		}
	}

	return nil
}

// extractVariablesFromArg extrait les noms de variables utilisées dans un argument d'action
func extractVariablesFromArg(arg interface{}) []string {
	// Si c'est une string simple, c'est potentiellement un nom de variable
	if str, ok := arg.(string); ok {
		return []string{str}
	}

	// Si c'est un objet (map), extraire les variables selon le type
	argMap, ok := arg.(map[string]interface{})
	if !ok {
		return []string{}
	}

	argType, _ := argMap["type"].(string)
	return extractVariablesByType(argType, argMap)
}

// extractVariablesByType extrait les variables en fonction du type d'argument
func extractVariablesByType(argType string, argMap map[string]interface{}) []string {
	switch argType {
	case "fieldAccess":
		return extractFromFieldAccess(argMap)
	case "variable":
		return extractFromVariable(argMap)
	case ArgTypeStringLiteral, "string", ArgTypeNumberLiteral, "number", ArgTypeBoolLiteral, ValueTypeBoolean:
		return []string{} // Literals ne contiennent pas de variables
	case ArgTypeFunctionCall:
		return extractFromFunctionCall(argMap)
	default:
		if isBinaryOperationType(argType) {
			return extractFromBinaryOp(argMap)
		}
		return []string{}
	}
}

// extractFromFieldAccess extrait la variable d'un accès à un champ
func extractFromFieldAccess(argMap map[string]interface{}) []string {
	if object, ok := argMap["object"].(string); ok {
		return []string{object}
	}
	return []string{}
}

// extractFromVariable extrait le nom d'une variable explicitement typée
func extractFromVariable(argMap map[string]interface{}) []string {
	if name, ok := argMap["name"].(string); ok {
		return []string{name}
	}
	return []string{}
}

// extractFromFunctionCall extrait les variables des arguments d'un appel de fonction
func extractFromFunctionCall(argMap map[string]interface{}) []string {
	vars := []string{}
	if args, ok := argMap["args"].([]interface{}); ok {
		for _, funcArg := range args {
			vars = append(vars, extractVariablesFromArg(funcArg)...)
		}
	}
	return vars
}

// extractFromBinaryOp extrait les variables des opérandes d'une opération binaire
func extractFromBinaryOp(argMap map[string]interface{}) []string {
	vars := []string{}
	if left := argMap["left"]; left != nil {
		vars = append(vars, extractVariablesFromArg(left)...)
	}
	if right := argMap["right"]; right != nil {
		vars = append(vars, extractVariablesFromArg(right)...)
	}
	return vars
}
