package rete

import (
	"fmt"
)

// createAction crée une action RETE à partir d'une map d'action parsée
func (cp *ConstraintPipeline) createAction(actionMap map[string]interface{}) *Action {
	actionType := getStringField(actionMap, "type", "print")

	action := &Action{
		Type: actionType,
		Job: JobCall{
			Name: getStringField(actionMap, "name", actionType),
			Args: []interface{}{},
		},
	}

	// Extraire les arguments si présents
	if argsData, hasArgs := actionMap["args"]; hasArgs {
		if argsList, ok := argsData.([]interface{}); ok {
			action.Job.Args = argsList
		}
	}

	return action
}

// buildConditionFromConstraints construit une condition appropriée à partir de contraintes
func (cp *ConstraintPipeline) buildConditionFromConstraints(constraintsData interface{}) (map[string]interface{}, error) {
	if constraintsData == nil {
		return map[string]interface{}{
			"type": "simple",
		}, nil
	}

	// Vérifier si c'est une agrégation
	if cp.detectAggregation(constraintsData) {
		return map[string]interface{}{
			"type": "passthrough",
		}, nil
	}

	// Analyser les contraintes pour détecter les négations
	isNegation, negatedCondition, err := cp.analyzeConstraints(constraintsData)
	if err != nil {
		return nil, fmt.Errorf("erreur analyse contraintes: %w", err)
	}

	if isNegation {
		fmt.Printf("   🚫 Détection contrainte NOT - création d'un AlphaNode de négation\n")
		return map[string]interface{}{
			"type":      "negation",
			"negated":   true,
			"condition": negatedCondition,
		}, nil
	}

	return map[string]interface{}{
		"type":       "constraint",
		"constraint": constraintsData,
	}, nil
}

// extractActionFromExpression extrait l'action d'une expression de règle
func (cp *ConstraintPipeline) extractActionFromExpression(exprMap map[string]interface{}, ruleID string) (*Action, error) {
	actionData, hasAction := exprMap["action"]
	if !hasAction {
		return nil, fmt.Errorf("aucune action trouvée pour règle %s", ruleID)
	}

	actionMap, ok := actionData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("format action invalide pour règle %s: %T", ruleID, actionData)
	}

	return cp.createAction(actionMap), nil
}

// determineRuleType détermine le type de règle (alpha, join, exists, accumulator)
// Retourne (ruleType string, shouldProcess bool)
func (cp *ConstraintPipeline) determineRuleType(
	exprMap map[string]interface{},
	variableCount int,
	hasAggregation bool,
) string {
	// Vérifier si c'est une contrainte EXISTS
	if constraintsData, hasConstraints := exprMap["constraints"]; hasConstraints {
		if cp.isExistsConstraint(constraintsData) {
			return "exists"
		}
	}

	// Si c'est une agrégation
	if hasAggregation {
		return "accumulator"
	}

	// Si plus d'une variable, c'est une jointure
	if variableCount > 1 {
		return "join"
	}

	// Sinon, c'est une règle alpha simple
	return "alpha"
}

// getVariableInfo extrait les informations de la première variable
// Retourne (variableName, variableType)
func (cp *ConstraintPipeline) getVariableInfo(variables []map[string]interface{}, variableTypes []string) (string, string) {
	variableName := "p" // défaut
	variableType := ""

	if len(variables) > 0 {
		if name, ok := variables[0]["name"].(string); ok {
			variableName = name
		}
		if len(variableTypes) > 0 {
			variableType = variableTypes[0]
		}
	}

	return variableName, variableType
}

// connectAlphaNodeToTypeNode connecte un AlphaNode au TypeNode approprié
func (cp *ConstraintPipeline) connectAlphaNodeToTypeNode(
	network *ReteNetwork,
	alphaNode *AlphaNode,
	variableType string,
	variableName string,
) {
	if variableType != "" {
		// Les TypeNodes sont stockés avec leur nom direct, pas avec "type_" préfixe
		if typeNode, exists := network.TypeNodes[variableType]; exists {
			typeNode.AddChild(alphaNode)
			fmt.Printf("   ✓ AlphaNode %s connecté au TypeNode %s\n", alphaNode.ID, variableType)
			return
		}
		fmt.Printf("   ⚠️  TypeNode %s non trouvé pour variable %s\n", variableType, variableName)
	} else {
		fmt.Printf("   ⚠️  Type de variable non trouvé pour %s, fallback\n", variableName)
	}

	// Fallback: connecter au premier type node trouvé
	for _, typeNode := range network.TypeNodes {
		typeNode.AddChild(alphaNode)
		break
	}
}

// createAlphaNodeWithTerminal crée un AlphaNode et son nœud terminal associé
func (cp *ConstraintPipeline) createAlphaNodeWithTerminal(
	network *ReteNetwork,
	ruleID string,
	condition map[string]interface{},
	variableName string,
	variableType string,
	action *Action,
	storage Storage,
) error {
	// Créer un nœud Alpha avec la condition appropriée
	alphaNode := NewAlphaNode(ruleID+"_alpha", condition, variableName, storage)

	// Connecter au type node approprié
	cp.connectAlphaNodeToTypeNode(network, alphaNode, variableType, variableName)

	network.AlphaNodes[alphaNode.ID] = alphaNode

	// Créer le terminal
	terminalNode := NewTerminalNode(ruleID+"_terminal", action, storage)
	alphaNode.AddChild(terminalNode)
	network.TerminalNodes[terminalNode.ID] = terminalNode

	if condition["type"] == "negation" {
		fmt.Printf("   ✓ AlphaNode de négation créé: %s -> %s\n", alphaNode.ID, terminalNode.ID)
	}

	return nil
}

// logRuleCreation affiche un message de log pour la création d'une règle
func (cp *ConstraintPipeline) logRuleCreation(ruleType string, ruleID string, variableNames []string) {
	switch ruleType {
	case "join":
		fmt.Printf("   📍 Règle multi-variables détectée (%d variables): %v\n", len(variableNames), variableNames)
	case "exists":
		fmt.Printf("   🔍 Règle EXISTS détectée pour: %s\n", ruleID)
	case "accumulator":
		fmt.Printf("   📊 Règle d'agrégation détectée pour: %s\n", ruleID)
	case "alpha":
		fmt.Printf("   ✓ Règle alpha simple créée pour: %s\n", ruleID)
	}
}
