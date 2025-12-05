// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

// Ce fichier contient les fonctions core d'extraction de composants AST pour le système de contraintes RETE.
//
// Pour les fonctions spécialisées, voir :
// - Agrégations : constraint_pipeline_aggregation.go
// - Jointures : constraint_pipeline_join.go
// - Variables : constraint_pipeline_variables.go
// - Détection : constraint_pipeline_detection.go

// extractComponents extrait les types et expressions d'un AST parsé
// Note: La gestion des instructions reset est effectuée en amont dans buildNetworkWithResetSemantics
// Cette fonction suppose que le programme reçu a déjà la sémantique reset appliquée si nécessaire
// Retourne (types, expressions, error)
func (cp *ConstraintPipeline) extractComponents(resultMap map[string]interface{}) ([]interface{}, []interface{}, error) {
	// Extraire les types
	typesData, hasTypes := resultMap["types"]
	if !hasTypes {
		return nil, nil, fmt.Errorf("aucun type trouvé dans l'AST")
	}

	types, ok := typesData.([]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("format types invalide: %T", typesData)
	}

	// Extraire les expressions
	expressionsData, hasExpressions := resultMap["expressions"]
	if !hasExpressions {
		return nil, nil, fmt.Errorf("aucune expression trouvée dans l'AST")
	}

	expressions, ok := expressionsData.([]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("format expressions invalide: %T", expressionsData)
	}

	return types, expressions, nil
}

// extractAndStoreActions extracts action definitions from the AST and stores them in the network
func (cp *ConstraintPipeline) extractAndStoreActions(network *ReteNetwork, resultMap map[string]interface{}) error {
	// Extraire les actions si présentes
	actionsData, hasActions := resultMap["actions"]
	if !hasActions {
		// Pas d'actions dans ce fichier, ce n'est pas une erreur
		return nil
	}

	actions, ok := actionsData.([]interface{})
	if !ok {
		return fmt.Errorf("format actions invalide: %T", actionsData)
	}

	// Traiter chaque action
	for _, actionData := range actions {
		actionMap, ok := actionData.(map[string]interface{})
		if !ok {
			return fmt.Errorf("format action invalide: %T", actionData)
		}

		// Extraire le nom de l'action
		actionName, ok := actionMap["name"].(string)
		if !ok {
			return fmt.Errorf("nom d'action non trouvé ou invalide")
		}

		// Créer la définition d'action
		actionDef := ActionDefinition{
			Type: "actionDefinition",
			Name: actionName,
		}

		// Extraire les paramètres si présents
		if paramsData, hasParams := actionMap["parameters"]; hasParams {
			params, ok := paramsData.([]interface{})
			if !ok {
				return fmt.Errorf("format paramètres invalide pour l'action %s: %T", actionName, paramsData)
			}

			actionDef.Parameters = make([]Parameter, len(params))
			for i, paramData := range params {
				paramMap, ok := paramData.(map[string]interface{})
				if !ok {
					return fmt.Errorf("format paramètre invalide pour l'action %s: %T", actionName, paramData)
				}

				paramName, ok := paramMap["name"].(string)
				if !ok {
					return fmt.Errorf("nom de paramètre non trouvé pour l'action %s", actionName)
				}

				paramType, ok := paramMap["type"].(string)
				if !ok {
					return fmt.Errorf("type de paramètre non trouvé pour l'action %s", actionName)
				}

				actionDef.Parameters[i] = Parameter{
					Name: paramName,
					Type: paramType,
				}
			}
		}

		// Ajouter à network.Actions (éviter les doublons)
		actionExists := false
		for i, existingAction := range network.Actions {
			if existingAction.Name == actionName {
				// Remplacer l'action existante
				network.Actions[i] = actionDef
				actionExists = true
				break
			}
		}

		if !actionExists {
			network.Actions = append(network.Actions, actionDef)
		}
	}

	return nil
}

// analyzeConstraints analyse les contraintes pour détecter les négations
// Retourne (isNegation, negatedCondition, error)
func (cp *ConstraintPipeline) analyzeConstraints(constraints interface{}) (bool, interface{}, error) {
	constraintMap, ok := constraints.(map[string]interface{})
	if !ok {
		return false, constraints, nil
	}

	// Détecter contrainte NOT
	if constraintType, exists := constraintMap["type"].(string); exists {
		if constraintType == "notConstraint" {
			// Extraire la contrainte niée
			if negatedConstraint, hasNegated := constraintMap["constraint"]; hasNegated {
				return true, negatedConstraint, nil
			}
		}
	}

	return false, constraints, nil
}
