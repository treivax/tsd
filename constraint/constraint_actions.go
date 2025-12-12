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
	var vars []string

	// Si c'est une string simple, c'est potentiellement un nom de variable
	// (le parser peut produire des simples strings pour les variables dans certains contextes)
	if str, ok := arg.(string); ok {
		vars = append(vars, str)
		return vars
	}

	// Si c'est un objet (map), extraire les variables selon le type
	if argMap, ok := arg.(map[string]interface{}); ok {
		argType, _ := argMap["type"].(string)
		switch argType {
		case "fieldAccess":
			// Accès à un champ d'objet : l'objet est une variable
			if object, ok := argMap["object"].(string); ok {
				vars = append(vars, object)
			}
		case "variable":
			// Variable explicitement typée
			if name, ok := argMap["name"].(string); ok {
				vars = append(vars, name)
			}
		case ArgTypeStringLiteral, "string":
			// String literals ne contiennent pas de variables
		case ArgTypeNumberLiteral, "number":
			// Number literals ne contiennent pas de variables
		case ArgTypeBoolLiteral, ValueTypeBoolean:
			// Boolean literals ne contiennent pas de variables
		default:
			// Vérifier si c'est un type d'opération binaire (plusieurs variantes possibles)
			if isBinaryOperationType(argType) {
				// Pour les opérations binaires, extraire récursivement des opérandes
				if left := argMap["left"]; left != nil {
					vars = append(vars, extractVariablesFromArg(left)...)
				}
				if right := argMap["right"]; right != nil {
					vars = append(vars, extractVariablesFromArg(right)...)
				}
			} else if argType == ArgTypeFunctionCall {
				// Pour les appels de fonction, extraire des arguments
				if args, ok := argMap["args"].([]interface{}); ok {
					for _, funcArg := range args {
						vars = append(vars, extractVariablesFromArg(funcArg)...)
					}
				}
			}
		}
	}

	return vars
}
