// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strings"
)

// Node condition type constants
const (
	ConditionTypePassthrough = "passthrough"
	ConditionTypeSimple      = "simple"
	ConditionTypeExists      = "exists"
	ConditionTypeComparison  = "comparison"
)

// Node side constants for beta nodes
const (
	NodeSideLeft  = "left"
	NodeSideRight = "right"
)

// buildNetwork construit le réseau RETE à partir des types et expressions parsés
func (cp *ConstraintPipeline) buildNetwork(storage Storage, types []interface{}, expressions []interface{}) (*ReteNetwork, error) {
	// Créer le réseau
	network := NewReteNetwork(storage)

	// ÉTAPE 1: Créer les TypeNodes
	err := cp.createTypeNodes(network, types, storage)
	if err != nil {
		return nil, fmt.Errorf("erreur création TypeNodes: %w", err)
	}

	// ÉTAPE 2: Créer les règles (AlphaNodes, BetaNodes, TerminalNodes)
	err = cp.createRuleNodes(network, expressions, storage)
	if err != nil {
		return nil, fmt.Errorf("erreur création règles: %w", err)
	}

	return network, nil
}

// createTypeNodes crée les TypeNodes à partir des définitions de types
func (cp *ConstraintPipeline) createTypeNodes(network *ReteNetwork, types []interface{}, storage Storage) error {
	for _, typeInterface := range types {
		typeMap, ok := typeInterface.(map[string]interface{})
		if !ok {
			return fmt.Errorf("format type invalide: %T", typeInterface)
		}

		// Extraire le nom du type
		typeName, ok := typeMap["name"].(string)
		if !ok {
			return fmt.Errorf("nom de type non trouvé")
		}

		// Créer la définition de type
		typeDef := cp.createTypeDefinition(typeName, typeMap)

		// Créer le TypeNode
		typeNode := NewTypeNode(typeName, typeDef, storage)
		network.TypeNodes[typeName] = typeNode

		// Enregistrer le TypeNode dans le LifecycleManager
		if network.LifecycleManager != nil {
			network.LifecycleManager.RegisterNode(typeNode.GetID(), "type")
		}

		// CRUCIAL: Connecter le TypeNode au RootNode pour permettre la propagation des faits
		network.RootNode.AddChild(typeNode)

		fmt.Printf("   ✓ TypeNode créé: %s\n", typeName)
	}

	return nil
}

// createTypeDefinition crée une définition de type à partir d'une map
func (cp *ConstraintPipeline) createTypeDefinition(typeName string, typeMap map[string]interface{}) TypeDefinition {
	typeDef := TypeDefinition{
		Type:   "type",
		Name:   typeName,
		Fields: []Field{},
	}

	// Extraire les champs
	fieldsData, hasFields := typeMap["fields"]
	if !hasFields {
		return typeDef
	}

	fields, ok := fieldsData.([]interface{})
	if !ok {
		return typeDef
	}

	for _, fieldInterface := range fields {
		fieldMap, ok := fieldInterface.(map[string]interface{})
		if !ok {
			continue
		}

		fieldName := getStringField(fieldMap, "name", "")
		fieldType := getStringField(fieldMap, "type", "")

		if fieldName != "" && fieldType != "" {
			typeDef.Fields = append(typeDef.Fields, Field{
				Name: fieldName,
				Type: fieldType,
			})
		}
	}

	return typeDef
}

// createRuleNodes crée les nœuds de règles (Alpha, Beta, Terminal) à partir des expressions
func (cp *ConstraintPipeline) createRuleNodes(network *ReteNetwork, expressions []interface{}, storage Storage) error {
	for i, exprInterface := range expressions {
		exprMap, ok := exprInterface.(map[string]interface{})
		if !ok {
			return fmt.Errorf("format expression invalide: %T", exprInterface)
		}

		// Générer un ID de règle
		ruleID := fmt.Sprintf("rule_%d", i)

		// Créer la règle
		err := cp.createSingleRule(network, ruleID, exprMap, storage)
		if err != nil {
			return fmt.Errorf("erreur création règle %s: %w", ruleID, err)
		}

		fmt.Printf("   ✓ Règle créée: %s\n", ruleID)
	}

	return nil
}

// createSingleRule crée une règle unique (refactorisée en petites fonctions)
func (cp *ConstraintPipeline) createSingleRule(network *ReteNetwork, ruleID string, exprMap map[string]interface{}, storage Storage) error {
	// Étape 1: Extraire l'action
	action, err := cp.extractActionFromExpression(exprMap, ruleID)
	if err != nil {
		return err
	}

	// Étape 2: Extraire et analyser les contraintes
	constraintsData, hasConstraints := exprMap["constraints"]
	var condition map[string]interface{}
	var hasAggregation bool

	if hasConstraints {
		// Détecter si c'est une agrégation
		hasAggregation = cp.detectAggregation(constraintsData)

		// Construire la condition appropriée
		condition, err = cp.buildConditionFromConstraints(constraintsData)
		if err != nil {
			return fmt.Errorf("erreur construction condition pour règle %s: %w", ruleID, err)
		}
	} else {
		condition = map[string]interface{}{
			"type": ConditionTypeSimple,
		}
	}

	// Étape 3: Extraire les variables
	variables, variableNames, variableTypes := cp.extractVariablesFromExpression(exprMap)

	// Étape 4: Déterminer le type de règle et la créer
	ruleType := cp.determineRuleType(exprMap, len(variables), hasAggregation)
	cp.logRuleCreation(ruleType, ruleID, variableNames)

	switch ruleType {
	case "exists":
		return cp.createExistsRule(network, ruleID, exprMap, condition, action, storage)

	case "accumulator":
		aggInfo, err := cp.extractAggregationInfo(constraintsData)
		if err != nil {
			fmt.Printf("   ⚠️  Impossible d'extraire info agrégation: %v, utilisation JoinNode standard\n", err)
			return cp.createJoinRule(network, ruleID, variables, variableNames, variableTypes, condition, action, storage)
		}
		return cp.createAccumulatorRule(network, ruleID, variables, variableNames, variableTypes, aggInfo, action, storage)

	case "join":
		return cp.createJoinRule(network, ruleID, variables, variableNames, variableTypes, condition, action, storage)

	case "alpha":
		return cp.createAlphaRule(network, ruleID, variables, variableNames, variableTypes, condition, action, storage)

	default:
		return fmt.Errorf("type de règle inconnu: %s", ruleType)
	}
}

// createAlphaRule crée une règle alpha simple avec une seule variable
func (cp *ConstraintPipeline) createAlphaRule(
	network *ReteNetwork,
	ruleID string,
	variables []map[string]interface{},
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	action *Action,
	storage Storage,
) error {
	// Extraire les informations de la variable
	variableName, variableType := cp.getVariableInfo(variables, variableTypes)

	// Créer l'AlphaNode avec son terminal
	return cp.createAlphaNodeWithTerminal(
		network,
		ruleID,
		condition,
		variableName,
		variableType,
		action,
		storage,
	)
}

// createJoinRule crée une règle de jointure avec JoinNode
func (cp *ConstraintPipeline) createJoinRule(
	network *ReteNetwork,
	ruleID string,
	variables []map[string]interface{},
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	action *Action,
	storage Storage,
) error {
	// Créer le nœud terminal pour cette règle
	terminalNode := NewTerminalNode(ruleID+"_terminal", action, storage)
	network.TerminalNodes[terminalNode.ID] = terminalNode

	// Déléguer à la fonction appropriée selon le nombre de variables
	if len(variableNames) > 2 {
		return cp.createCascadeJoinRule(network, ruleID, variableNames, variableTypes, condition, terminalNode, storage)
	}

	return cp.createBinaryJoinRule(network, ruleID, variableNames, variableTypes, condition, terminalNode, storage)
}

// createExistsRule crée une règle EXISTS avec ExistsNode
func (cp *ConstraintPipeline) createExistsRule(
	network *ReteNetwork,
	ruleID string,
	exprMap map[string]interface{},
	condition map[string]interface{},
	action *Action,
	storage Storage,
) error {
	// Créer le nœud terminal pour cette règle
	terminalNode := NewTerminalNode(ruleID+"_terminal", action, storage)
	network.TerminalNodes[terminalNode.ID] = terminalNode

	// Extraire les variables
	mainVariable, existsVariable, mainVarType, existsVarType, err := cp.extractExistsVariables(exprMap)
	if err != nil {
		return err
	}

	// Extraire les conditions d'EXISTS
	existsConditions, err := cp.extractExistsConditions(exprMap)
	if err != nil {
		return err
	}

	// Créer l'objet condition pour l'ExistsNode
	existsConditionObj := map[string]interface{}{
		"type":       ConditionTypeExists,
		"conditions": existsConditions,
	}

	// Créer le mapping variable -> type pour l'ExistsNode
	varTypes := make(map[string]string)
	varTypes[mainVariable] = mainVarType
	varTypes[existsVariable] = existsVarType

	// Créer l'ExistsNode avec les vraies conditions
	existsNode := NewExistsNode(ruleID+"_exists", existsConditionObj, mainVariable, existsVariable, varTypes, storage)
	existsNode.AddChild(terminalNode)

	// Stocker l'ExistsNode dans les BetaNodes du réseau
	network.BetaNodes[existsNode.ID] = existsNode

	// Créer des AlphaNodes pass-through pour les deux variables
	cp.connectExistsNodeToTypeNodes(network, ruleID, existsNode, mainVariable, mainVarType, existsVariable, existsVarType)

	fmt.Printf("   ✅ ExistsNode %s créé pour %s EXISTS %s\n", existsNode.ID, mainVariable, existsVariable)
	return nil
}

// extractExistsVariables extrait les variables d'une règle EXISTS
func (cp *ConstraintPipeline) extractExistsVariables(exprMap map[string]interface{}) (string, string, string, string, error) {
	var mainVariable, existsVariable string
	var mainVarType, existsVarType string

	// Extraire la variable principale depuis "set"
	if setData, hasSet := exprMap["set"]; hasSet {
		if setMap, ok := setData.(map[string]interface{}); ok {
			if varsData, hasVars := setMap["variables"]; hasVars {
				if varsList, ok := varsData.([]interface{}); ok && len(varsList) > 0 {
					if varMap, ok := varsList[0].(map[string]interface{}); ok {
						if name, ok := varMap["name"].(string); ok {
							mainVariable = name
						}
						if dataType, ok := varMap["dataType"].(string); ok {
							mainVarType = dataType
						}
					}
				}
			}
		}
	}

	// Extraire la variable d'existence depuis les contraintes
	if constraintsData, hasConstraints := exprMap["constraints"]; hasConstraints {
		if constraintMap, ok := constraintsData.(map[string]interface{}); ok {
			if variable, hasVar := constraintMap["variable"]; hasVar {
				if varMap, ok := variable.(map[string]interface{}); ok {
					if name, ok := varMap["name"].(string); ok {
						existsVariable = name
					}
					if dataType, ok := varMap["dataType"].(string); ok {
						existsVarType = dataType
					}
				}
			}
		}
	}

	if mainVariable == "" || existsVariable == "" {
		return "", "", "", "", fmt.Errorf("variables EXISTS non trouvées: main=%s, exists=%s", mainVariable, existsVariable)
	}

	return mainVariable, existsVariable, mainVarType, existsVarType, nil
}

// extractExistsConditions extrait les conditions d'une règle EXISTS
func (cp *ConstraintPipeline) extractExistsConditions(exprMap map[string]interface{}) ([]map[string]interface{}, error) {
	var existsConditions []map[string]interface{}

	if constraintsData, hasConstraints := exprMap["constraints"]; hasConstraints {
		if constraintMap, ok := constraintsData.(map[string]interface{}); ok {
			// Essayer d'abord "condition" (au singulier)
			if conditionData, hasCondition := constraintMap["condition"]; hasCondition {
				if conditionObj, ok := conditionData.(map[string]interface{}); ok {
					existsConditions = append(existsConditions, conditionObj)
				}
			}
			// Puis essayer "conditions" (au pluriel) si pas trouvé
			if len(existsConditions) == 0 {
				if conditionsData, hasConditions := constraintMap["conditions"]; hasConditions {
					if conditionsList, ok := conditionsData.([]interface{}); ok {
						for _, conditionData := range conditionsList {
							if conditionObj, ok := conditionData.(map[string]interface{}); ok {
								existsConditions = append(existsConditions, conditionObj)
							}
						}
					}
				}
			}
		}
	}

	return existsConditions, nil
}

// connectExistsNodeToTypeNodes connecte un ExistsNode aux TypeNodes appropriés
func (cp *ConstraintPipeline) connectExistsNodeToTypeNodes(
	network *ReteNetwork,
	ruleID string,
	existsNode *ExistsNode,
	mainVariable string,
	mainVarType string,
	existsVariable string,
	existsVarType string,
) {
	// Connecter les variables principale et d'existence à l'ExistsNode
	if mainVarType != "" {
		cp.connectTypeNodeToBetaNode(network, ruleID, mainVariable, mainVarType, existsNode, NodeSideLeft)
	}
	if existsVarType != "" {
		cp.connectTypeNodeToBetaNode(network, ruleID, existsVariable, existsVarType, existsNode, NodeSideRight)
	}
}

// createAccumulatorRule crée une règle avec AccumulatorNode
func (cp *ConstraintPipeline) createAccumulatorRule(
	network *ReteNetwork,
	ruleID string,
	variables []map[string]interface{},
	variableNames []string,
	variableTypes []string,
	aggInfo *AggregationInfo,
	action *Action,
	storage Storage,
) error {
	// Extraire la variable principale et son type depuis variables
	if len(variables) == 0 || len(variableTypes) == 0 {
		return fmt.Errorf("aucune variable principale trouvée")
	}

	mainVariable := variableNames[0]
	mainType := variableTypes[0]

	// Stocker dans aggInfo
	aggInfo.MainVariable = mainVariable
	aggInfo.MainType = mainType

	// Créer le nœud terminal
	terminalNode := NewTerminalNode(ruleID+"_terminal", action, storage)
	network.TerminalNodes[terminalNode.ID] = terminalNode

	// Créer la condition de comparaison
	condition := map[string]interface{}{
		"type":     ConditionTypeComparison,
		"operator": aggInfo.Operator,
		"value":    aggInfo.Threshold,
	}

	// Créer l'AccumulatorNode avec tous les paramètres
	accumNode := NewAccumulatorNode(
		ruleID+"_accum",
		aggInfo.MainVariable, // "e"
		aggInfo.MainType,     // "Employee"
		aggInfo.AggVariable,  // "p"
		aggInfo.AggType,      // "Performance"
		aggInfo.Field,        // "score"
		aggInfo.JoinField,    // "employee_id"
		aggInfo.MainField,    // "id"
		aggInfo.Function,     // "AVG"
		condition,
		storage,
	)
	accumNode.AddChild(terminalNode)
	network.BetaNodes[accumNode.ID] = accumNode

	// Connecter les TypeNodes à l'AccumulatorNode
	cp.connectTypeNodeToBetaNode(network, ruleID, mainVariable, mainType, accumNode, "")
	fmt.Printf("   ✓ %s -> PassthroughAlpha -> AccumulatorNode[%s]\n", mainType, aggInfo.Function)

	cp.connectTypeNodeToBetaNode(network, ruleID, aggInfo.AggVariable, aggInfo.AggType, accumNode, "")
	fmt.Printf("   ✓ %s -> PassthroughAlpha -> AccumulatorNode[%s] (pour agrégation)\n", aggInfo.AggType, aggInfo.Function)

	fmt.Printf("   ✅ AccumulatorNode %s créé pour %s(%s.%s) %s %.2f\n",
		accumNode.ID, aggInfo.Function, aggInfo.AggVariable, aggInfo.Field, aggInfo.Operator, aggInfo.Threshold)
	return nil
}

// createPassthroughAlphaNode creates a passthrough AlphaNode with optional side specification
func (cp *ConstraintPipeline) createPassthroughAlphaNode(ruleID, varName, side string, storage Storage) *AlphaNode {
	passCondition := map[string]interface{}{
		"type": ConditionTypePassthrough,
	}
	if side != "" {
		passCondition["side"] = side
	}
	return NewAlphaNode(ruleID+"_pass_"+varName, passCondition, varName, storage)
}

// connectTypeNodeToBetaNode connects a TypeNode to a BetaNode via a passthrough AlphaNode
func (cp *ConstraintPipeline) connectTypeNodeToBetaNode(
	network *ReteNetwork,
	ruleID string,
	varName string,
	varType string,
	betaNode Node,
	side string,
) {
	if typeNode, exists := network.TypeNodes[varType]; exists {
		alphaNode := cp.createPassthroughAlphaNode(ruleID, varName, side, network.Storage)
		typeNode.AddChild(alphaNode)
		alphaNode.AddChild(betaNode)

		sideInfo := ""
		if side != "" {
			sideInfo = fmt.Sprintf(" (%s)", strings.ToUpper(side))
		}
		fmt.Printf("   ✓ %s -> PassthroughAlpha_%s -> %s%s\n", varType, varName, betaNode.GetID(), sideInfo)
	}
}

// createBinaryJoinRule creates a simple binary join rule (2 variables)
func (cp *ConstraintPipeline) createBinaryJoinRule(
	network *ReteNetwork,
	ruleID string,
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	terminalNode *TerminalNode,
	storage Storage,
) error {
	leftVars := []string{variableNames[0]}
	rightVars := []string{variableNames[1]}

	// Créer le mapping variable -> type
	varTypes := make(map[string]string)
	for i, varName := range variableNames {
		varTypes[varName] = variableTypes[i]
	}

	joinNode := NewJoinNode(ruleID+"_join", condition, leftVars, rightVars, varTypes, storage)
	joinNode.AddChild(terminalNode)

	// Stocker le JoinNode dans les BetaNodes du réseau
	network.BetaNodes[joinNode.ID] = joinNode

	// Connecter les TypeNodes via des AlphaNodes pass-through
	for i, varName := range variableNames {
		varType := variableTypes[i]
		if varType != "" {
			side := NodeSideRight
			if i == 0 {
				side = NodeSideLeft
			}
			cp.connectTypeNodeToBetaNode(network, ruleID, varName, varType, joinNode, side)
		} else {
			fmt.Printf("   ⚠️ Type vide pour variable %s\n", varName)
		}
	}

	fmt.Printf("   ✅ JoinNode %s créé pour jointure %s\n", joinNode.ID, strings.Join(variableNames, " ⋈ "))
	return nil
}

// createCascadeJoinRule creates a cascade of join nodes for multi-variable rules (3+ variables)
func (cp *ConstraintPipeline) createCascadeJoinRule(
	network *ReteNetwork,
	ruleID string,
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	terminalNode *TerminalNode,
	storage Storage,
) error {
	fmt.Printf("   📍 Règle multi-variables détectée (%d variables): %v\n", len(variableNames), variableNames)
	fmt.Printf("   🔧 Construction d'architecture en cascade de JoinNodes\n")

	// Créer le mapping variable -> type
	varTypes := make(map[string]string)
	for i, varName := range variableNames {
		varTypes[varName] = variableTypes[i]
	}

	// Étape 1: Créer le premier JoinNode pour les 2 premières variables
	leftVars := []string{variableNames[0]}
	rightVars := []string{variableNames[1]}
	currentVarTypes := map[string]string{
		variableNames[0]: variableTypes[0],
		variableNames[1]: variableTypes[1],
	}

	currentJoinNode := NewJoinNode(
		fmt.Sprintf("%s_join_%d_%d", ruleID, 0, 1),
		condition,
		leftVars,
		rightVars,
		currentVarTypes,
		storage,
	)
	network.BetaNodes[currentJoinNode.ID] = currentJoinNode

	// Connecter les 2 premières variables au premier JoinNode
	for i := 0; i < 2; i++ {
		varName := variableNames[i]
		varType := variableTypes[i]
		side := NodeSideRight
		if i == 0 {
			side = NodeSideLeft
		}
		cp.connectTypeNodeToBetaNode(network, ruleID, varName, varType, currentJoinNode, side)
		fmt.Printf("   ✓ Cascade level 1 connection\n")
	}

	fmt.Printf("   ✅ JoinNode cascade level 1: %s ⋈ %s\n", variableNames[0], variableNames[1])

	// Étape 2+: Joindre chaque variable suivante au résultat précédent
	for i := 2; i < len(variableNames); i++ {
		nextVarName := variableNames[i]
		nextVarType := variableTypes[i]

		// Variables accumulées jusqu'ici
		accumulatedVars := variableNames[0:i]
		accumulatedVarTypes := make(map[string]string)
		for j := 0; j < i; j++ {
			accumulatedVarTypes[variableNames[j]] = variableTypes[j]
		}
		accumulatedVarTypes[nextVarName] = nextVarType

		// Créer le prochain JoinNode
		nextJoinNode := NewJoinNode(
			fmt.Sprintf("%s_join_%d", ruleID, i),
			condition,
			accumulatedVars,
			[]string{nextVarName},
			accumulatedVarTypes,
			storage,
		)
		network.BetaNodes[nextJoinNode.ID] = nextJoinNode

		// Connecter le JoinNode précédent au nouveau JoinNode
		currentJoinNode.AddChild(nextJoinNode)

		// Connecter la nouvelle variable au JoinNode
		cp.connectTypeNodeToBetaNode(network, ruleID, nextVarName, nextVarType, nextJoinNode, NodeSideRight)
		fmt.Printf("   ✓ Cascade level %d connection\n", i)

		fmt.Printf("   ✅ JoinNode cascade level %d: (%s) ⋈ %s\n", i, strings.Join(accumulatedVars, " ⋈ "), nextVarName)

		currentJoinNode = nextJoinNode
	}

	// Connecter le dernier JoinNode au terminal
	currentJoinNode.AddChild(terminalNode)
	fmt.Printf("   ✅ Architecture en cascade complète: %s\n", strings.Join(variableNames, " ⋈ "))

	return nil
}
