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

	// Si c'est une string simple, c'est probablement un nom de variable
	if str, ok := arg.(string); ok {
		vars = append(vars, str)
		return vars
	}

	// Si c'est un objet (map), extraire les variables selon le type
	if argMap, ok := arg.(map[string]interface{}); ok {
		argType, _ := argMap["type"].(string)
		switch argType {
		case "fieldAccess":
			if object, ok := argMap["object"].(string); ok {
				vars = append(vars, object)
			}
		case "string":
			// Les string literals ne contiennent pas de variables
		case "number":
			// Les number literals ne contiennent pas de variables
		default:
			// Pour d'autres types, on peut chercher récursivement
			// mais pour l'instant on ignore
		}
	}

	return vars
}
