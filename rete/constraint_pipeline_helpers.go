// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

// createAction crée une action RETE à partir d'une map d'action parsée
func (cp *ConstraintPipeline) createAction(actionMap map[string]interface{}) *Action {
	actionType := getStringField(actionMap, "type", "print")

	// Extraire le job depuis l'action
	jobData, hasJob := actionMap["job"]
	if !hasJob {
		// Fallback: action simple sans job (ne devrait pas arriver avec le nouveau parser)
		return &Action{
			Type: actionType,
			Job: JobCall{
				Name: actionType,
				Args: []interface{}{},
			},
		}
	}

	jobMap, ok := jobData.(map[string]interface{})
	if !ok {
		return &Action{
			Type: actionType,
			Job: JobCall{
				Name: actionType,
				Args: []interface{}{},
			},
		}
	}

	// Extraire le nom du job
	jobName := getStringField(jobMap, "name", actionType)

	action := &Action{
		Type: actionType,
		Job: JobCall{
			Name: jobName,
			Args: []interface{}{},
		},
	}

	// Extraire les arguments du job (pas de l'action)
	if argsData, hasArgs := jobMap["args"]; hasArgs {
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

// createAlphaNodeWithTerminal crée un AlphaNode (partagé si possible) et son nœud terminal associé
// Cette fonction analyse l'expression et construit une chaîne si possible, sinon utilise le comportement simple
func (cp *ConstraintPipeline) createAlphaNodeWithTerminal(
	network *ReteNetwork,
	ruleID string,
	condition interface{},
	variableName string,
	variableType string,
	action *Action,
	storage Storage,
) error {
	// Déballer la condition si elle est wrappée dans une map
	actualCondition := condition
	if condMap, ok := condition.(map[string]interface{}); ok {
		if condType, hasType := condMap["type"]; hasType {
			if condType == "constraint" {
				if constraint, hasConstraint := condMap["constraint"]; hasConstraint {
					actualCondition = constraint
				}
			} else if condType == "negation" {
				// Pour les négations, utiliser le comportement simple
				return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
			} else if condType == "simple" || condType == "passthrough" {
				// Pas de vraie condition, utiliser le comportement simple
				return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
			}
		}
	}

	// Analyser l'expression pour déterminer son type
	exprType, err := AnalyzeExpression(actualCondition)
	if err != nil {
		fmt.Printf("   ⚠️  Erreur analyse expression: %v, fallback vers comportement simple\n", err)
		return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
	}

	// Cas spécial: expressions OR et mixtes - normalisation avancée avec support des OR imbriqués
	// DOIT être traité AVANT le check CanDecompose car OR n'est pas décomposable
	if exprType == ExprTypeOR || exprType == ExprTypeMixed {
		if exprType == ExprTypeOR {
			fmt.Printf("   ℹ️  Expression OR détectée, normalisation avancée et création d'un nœud alpha unique\n")
		} else {
			fmt.Printf("   ℹ️  Expression mixte (AND+OR) détectée, normalisation avancée et création d'un nœud alpha unique\n")
		}

		// Analyser la complexité de l'expression pour déterminer la stratégie de normalisation
		analysis, err := AnalyzeNestedOR(actualCondition)
		if err != nil {
			fmt.Printf("   ⚠️  Erreur analyse OR imbriqué: %v, fallback vers normalisation simple\n", err)
			// Fallback vers normalisation simple
			normalizedExpr, err := NormalizeORExpression(actualCondition)
			if err != nil {
				fmt.Printf("   ⚠️  Erreur normalisation simple: %v, fallback vers comportement simple\n", err)
				return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
			}
			normalizedCondition := map[string]interface{}{
				"type":       "constraint",
				"constraint": normalizedExpr,
			}
			return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, normalizedCondition, variableName, variableType, action, storage)
		}

		// Afficher les informations d'analyse
		fmt.Printf("   📊 Analyse OR: Complexité=%v, Profondeur=%d, OR=%d, AND=%d\n",
			analysis.Complexity, analysis.NestingDepth, analysis.ORTermCount, analysis.ANDTermCount)

		if analysis.OptimizationHint != "" {
			fmt.Printf("   💡 Suggestion: %s\n", analysis.OptimizationHint)
		}

		// Utiliser la normalisation avancée pour les expressions complexes
		var normalizedExpr interface{}
		if analysis.RequiresFlattening || analysis.RequiresDNF {
			fmt.Printf("   🔧 Application de la normalisation avancée (aplatissement=%v, DNF=%v)\n",
				analysis.RequiresFlattening, analysis.RequiresDNF)

			normalizedExpr, err = NormalizeNestedOR(actualCondition)
			if err != nil {
				fmt.Printf("   ⚠️  Erreur normalisation avancée: %v, fallback vers normalisation simple\n", err)
				// Fallback vers normalisation simple
				normalizedExpr, err = NormalizeORExpression(actualCondition)
				if err != nil {
					fmt.Printf("   ⚠️  Erreur normalisation simple: %v, utilisation expression originale\n", err)
					normalizedExpr = actualCondition
				}
			} else {
				fmt.Printf("   ✅ Normalisation avancée réussie\n")
			}
		} else {
			// Pour les expressions simples, utiliser la normalisation standard
			fmt.Printf("   🔧 Application de la normalisation standard\n")
			normalizedExpr, err = NormalizeORExpression(actualCondition)
			if err != nil {
				fmt.Printf("   ⚠️  Erreur normalisation: %v, utilisation expression originale\n", err)
				normalizedExpr = actualCondition
			}
		}

		// Créer un seul AlphaNode avec l'expression normalisée
		normalizedCondition := map[string]interface{}{
			"type":       "constraint",
			"constraint": normalizedExpr,
		}

		return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, normalizedCondition, variableName, variableType, action, storage)
	}

	// Vérifier si l'expression peut être décomposée
	if !CanDecompose(exprType) {
		fmt.Printf("   ℹ️  Expression de type %s non décomposable, utilisation du nœud simple\n", exprType)
		return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
	}

	// Cas spécial: expressions simples - utiliser le comportement actuel
	if exprType == ExprTypeSimple || exprType == ExprTypeArithmetic {
		return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
	}

	// Expressions AND ou NOT - tenter la décomposition en chaîne
	fmt.Printf("   🔍 Expression de type %s détectée, tentative de décomposition...\n", exprType)

	// Extraire les conditions de l'expression (utiliser la condition déballée)
	conditions, opType, err := ExtractConditions(actualCondition)
	if err != nil {
		fmt.Printf("   ⚠️  Erreur extraction conditions: %v, fallback vers comportement simple\n", err)
		return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
	}

	// Si une seule condition, pas besoin de chaîne
	if len(conditions) <= 1 {
		fmt.Printf("   ℹ️  Une seule condition extraite, utilisation du nœud simple\n")
		return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
	}

	fmt.Printf("   🔗 Décomposition en chaîne: %d conditions détectées (opérateur: %s)\n", len(conditions), opType)

	// Normaliser les conditions
	normalizedConditions := NormalizeConditions(conditions, opType)
	fmt.Printf("   📋 Conditions normalisées: %d condition(s)\n", len(normalizedConditions))

	// Trouver le TypeNode parent pour connecter la chaîne
	var parentNode Node
	if variableType != "" {
		if typeNode, exists := network.TypeNodes[variableType]; exists {
			parentNode = typeNode
		}
	}

	// Si pas de TypeNode trouvé, utiliser le premier disponible
	if parentNode == nil {
		for _, typeNode := range network.TypeNodes {
			parentNode = typeNode
			break
		}
	}

	if parentNode == nil {
		fmt.Printf("   ⚠️  Aucun TypeNode trouvé, fallback vers comportement simple\n")
		return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
	}

	// Créer le constructeur de chaîne
	chainBuilder := NewAlphaChainBuilder(network, storage)

	// Construire la chaîne d'AlphaNodes
	chain, err := chainBuilder.BuildChain(normalizedConditions, variableName, parentNode, ruleID)
	if err != nil {
		fmt.Printf("   ⚠️  Erreur construction chaîne: %v, fallback vers comportement simple\n", err)
		return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
	}

	// Valider la chaîne
	if err := chain.ValidateChain(); err != nil {
		fmt.Printf("   ⚠️  Chaîne invalide: %v, fallback vers comportement simple\n", err)
		return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
	}

	// Obtenir les statistiques de la chaîne
	stats := chainBuilder.GetChainStats(chain)
	sharedCount := 0
	if sc, ok := stats["shared_nodes"].(int); ok {
		sharedCount = sc
	}

	// Afficher les statistiques de construction
	fmt.Printf("   ✅ Chaîne construite: %d nœud(s), %d partagé(s)\n", len(chain.Nodes), sharedCount)

	// Logger les détails de chaque nœud
	for i, node := range chain.Nodes {
		if i < sharedCount {
			fmt.Printf("   ♻️  AlphaNode partagé réutilisé: %s (hash: %s)\n", node.ID, chain.Hashes[i])
		} else {
			fmt.Printf("   ✨ Nouveau AlphaNode créé: %s (hash: %s)\n", node.ID, chain.Hashes[i])
		}
	}

	// Créer et attacher le terminal au dernier nœud de la chaîne
	terminalNode := NewTerminalNode(ruleID+"_terminal", action, storage)
	chain.FinalNode.AddChild(terminalNode)
	network.TerminalNodes[terminalNode.ID] = terminalNode

	// Register terminal node with lifecycle manager
	if network.LifecycleManager != nil {
		network.LifecycleManager.RegisterNode(terminalNode.ID, "terminal")
		network.LifecycleManager.AddRuleToNode(terminalNode.ID, ruleID, ruleID)
	}

	fmt.Printf("   ✓ TerminalNode %s attaché au nœud final %s de la chaîne\n", terminalNode.ID, chain.FinalNode.ID)

	return nil
}

// createSimpleAlphaNodeWithTerminal crée un AlphaNode simple (partagé si possible) et son nœud terminal associé
// Cette fonction implémente le comportement original pour les conditions simples ou non-décomposables
func (cp *ConstraintPipeline) createSimpleAlphaNodeWithTerminal(
	network *ReteNetwork,
	ruleID string,
	condition interface{},
	variableName string,
	variableType string,
	action *Action,
	storage Storage,
) error {
	// Convertir la condition en map si nécessaire
	var conditionMap map[string]interface{}
	switch c := condition.(type) {
	case map[string]interface{}:
		conditionMap = c
	default:
		// Pour les types structurés (constraint.*), les passer directement
		conditionMap = map[string]interface{}{
			"type":       "constraint",
			"constraint": condition,
		}
	}
	// Utiliser le gestionnaire de partage pour obtenir ou créer un AlphaNode
	alphaNode, alphaHash, wasShared, err := network.AlphaSharingManager.GetOrCreateAlphaNode(
		conditionMap,
		variableName,
		storage,
	)
	if err != nil {
		return fmt.Errorf("erreur création AlphaNode partagé: %w", err)
	}

	if wasShared {
		fmt.Printf("   ♻️  AlphaNode partagé réutilisé: %s (hash: %s)\n", alphaNode.ID, alphaHash)
	} else {
		fmt.Printf("   ✨ Nouveau AlphaNode partageable créé: %s (hash: %s)\n", alphaNode.ID, alphaHash)

		// Connecter au type node approprié (seulement pour les nouveaux nœuds)
		cp.connectAlphaNodeToTypeNode(network, alphaNode, variableType, variableName)

		// Ajouter au registre global des AlphaNodes du réseau
		network.AlphaNodes[alphaNode.ID] = alphaNode
	}

	// Enregistrer ou mettre à jour l'AlphaNode dans le LifecycleManager
	if network.LifecycleManager != nil {
		lifecycle := network.LifecycleManager.RegisterNode(alphaNode.ID, "alpha")
		lifecycle.AddRuleReference(ruleID, ruleID)
	}

	// Créer le terminal (toujours spécifique à la règle)
	terminalNode := NewTerminalNode(ruleID+"_terminal", action, storage)
	alphaNode.AddChild(terminalNode)
	network.TerminalNodes[terminalNode.ID] = terminalNode

	// Register terminal node with lifecycle manager
	if network.LifecycleManager != nil {
		network.LifecycleManager.RegisterNode(terminalNode.ID, "terminal")
		network.LifecycleManager.AddRuleToNode(terminalNode.ID, ruleID, ruleID)
	}

	if conditionMap["type"] == "negation" {
		fmt.Printf("   ✓ AlphaNode de négation créé: %s -> %s\n", alphaNode.ID, terminalNode.ID)
	} else if wasShared {
		fmt.Printf("   ✓ Règle %s attachée à l'AlphaNode partagé %s via terminal %s\n",
			ruleID, alphaNode.ID, terminalNode.ID)
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
