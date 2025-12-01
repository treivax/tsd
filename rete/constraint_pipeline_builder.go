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

		// Extraire le ruleId de l'expression
		ruleID := fmt.Sprintf("rule_%d", i) // Default fallback
		if ruleIdValue, ok := exprMap["ruleId"]; ok {
			if ruleIdStr, ok := ruleIdValue.(string); ok && ruleIdStr != "" {
				ruleID = ruleIdStr
			}
		}

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
		// Détecter si c'est une agrégation (from constraints)
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

	// Also check if any variables are aggregation variables (new syntax)
	if !hasAggregation {
		hasAggregation = cp.hasAggregationVariables(exprMap)
	}

	// Étape 4: Déterminer le type de règle et la créer
	ruleType := cp.determineRuleType(exprMap, len(variables), hasAggregation)
	cp.logRuleCreation(ruleType, ruleID, variableNames)

	switch ruleType {
	case "exists":
		return cp.createExistsRule(network, ruleID, exprMap, condition, action, storage)

	case "accumulator":
		var aggInfo *AggregationInfo
		var err error

		// Check if this is the new multi-pattern aggregation syntax
		if _, hasPatterns := exprMap["patterns"]; hasPatterns {
			// Check if this is multi-source aggregation
			if cp.isMultiSourceAggregation(exprMap) {
				fmt.Printf("   📊 Multi-source aggregation détectée pour: %s\n", ruleID)
				aggInfo, err = cp.extractMultiSourceAggregationInfo(exprMap)
			} else {
				aggInfo, err = cp.extractAggregationInfoFromVariables(exprMap)
			}
		} else {
			// Old AccumulateConstraint syntax
			aggInfo, err = cp.extractAggregationInfo(constraintsData)
		}

		if err != nil {
			fmt.Printf("   ⚠️  Impossible d'extraire info agrégation: %v, utilisation JoinNode standard\n", err)
			return cp.createJoinRule(network, ruleID, variables, variableNames, variableTypes, condition, action, storage)
		}

		// Check if we need multi-source accumulator
		if len(aggInfo.SourcePatterns) > 1 || len(aggInfo.AggregationVars) > 1 {
			return cp.createMultiSourceAccumulatorRule(network, ruleID, aggInfo, action, storage)
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

// isMultiSourceAggregation checks if the rule has multiple aggregation sources
func (cp *ConstraintPipeline) isMultiSourceAggregation(exprMap map[string]interface{}) bool {
	patternsData, hasPatterns := exprMap["patterns"]
	if !hasPatterns {
		return false
	}

	patternsList, ok := patternsData.([]interface{})
	if !ok {
		return false
	}

	// Multi-source if we have more than 2 pattern blocks
	if len(patternsList) > 2 {
		return true
	}

	// Or if we have multiple aggregation variables
	if len(patternsList) >= 1 {
		firstPattern, ok := patternsList[0].(map[string]interface{})
		if !ok {
			return false
		}

		varsData, hasVars := firstPattern["variables"]
		if !hasVars {
			return false
		}

		varsList, ok := varsData.([]interface{})
		if !ok {
			return false
		}

		aggVarCount := 0
		for _, varInterface := range varsList {
			if varMap, ok := varInterface.(map[string]interface{}); ok {
				if varType, ok := varMap["type"].(string); ok && varType == "aggregationVariable" {
					aggVarCount++
					if aggVarCount > 1 {
						return true
					}
				}
			}
		}
	}

	return false
}

// createMultiSourceAccumulatorRule creates a rule with multiple aggregation sources
func (cp *ConstraintPipeline) createMultiSourceAccumulatorRule(
	network *ReteNetwork,
	ruleID string,
	aggInfo *AggregationInfo,
	action *Action,
	storage Storage,
) error {
	fmt.Printf("   🔗 Création règle multi-source avec %d sources et %d agrégations\n",
		len(aggInfo.SourcePatterns), len(aggInfo.AggregationVars))

	// Create a join chain that feeds into a MultiSourceAccumulatorNode
	// The join chain combines all sources, then the accumulator computes aggregations

	// Create a join chain for all source patterns
	var lastJoinNode Node

	// Start with the main variable's type node
	mainTypeNode, exists := network.TypeNodes[aggInfo.MainType]
	if !exists {
		return fmt.Errorf("TypeNode pour %s non trouvé", aggInfo.MainType)
	}

	// For each source pattern, create a join node
	for i, sourcePattern := range aggInfo.SourcePatterns {
		sourceTypeNode, exists := network.TypeNodes[sourcePattern.Type]
		if !exists {
			return fmt.Errorf("TypeNode pour %s non trouvé", sourcePattern.Type)
		}

		// Create join node - build condition map and variable lists
		joinConditionMap := make(map[string]interface{})
		var leftVars, rightVars []string
		varTypes := make(map[string]string)

		// Find the join condition for this source
		var joinCondition *JoinCondition
		for j := range aggInfo.JoinConditions {
			cond := &aggInfo.JoinConditions[j]
			if cond.LeftVar == sourcePattern.Variable || cond.RightVar == sourcePattern.Variable {
				joinCondition = cond
				break
			}
		}

		if joinCondition != nil {
			// Build a simple comparison condition
			joinConditionMap = map[string]interface{}{
				"type":     "comparison",
				"operator": joinCondition.Operator,
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": joinCondition.LeftVar,
					"field":  joinCondition.LeftField,
				},
				"right": map[string]interface{}{
					"type":   "fieldAccess",
					"object": joinCondition.RightVar,
					"field":  joinCondition.RightField,
				},
			}

			fmt.Printf("   ✓ Creating JoinNode: %s.%s == %s.%s\n",
				joinCondition.LeftVar, joinCondition.LeftField,
				joinCondition.RightVar, joinCondition.RightField)
		}

		// Build variable lists - always accumulate variables from left side
		if i == 0 {
			// First join: main + first source
			leftVars = []string{aggInfo.MainVariable}
			varTypes[aggInfo.MainVariable] = aggInfo.MainType
		} else {
			// Subsequent joins: all previous variables on left side
			leftVars = []string{aggInfo.MainVariable}
			varTypes[aggInfo.MainVariable] = aggInfo.MainType
			for j := 0; j < i; j++ {
				leftVars = append(leftVars, aggInfo.SourcePatterns[j].Variable)
				varTypes[aggInfo.SourcePatterns[j].Variable] = aggInfo.SourcePatterns[j].Type
			}
		}
		rightVars = []string{sourcePattern.Variable}
		varTypes[sourcePattern.Variable] = sourcePattern.Type

		joinNodeID := fmt.Sprintf("%s_join_%d", ruleID, i)
		joinNode := NewJoinNode(joinNodeID, joinConditionMap, leftVars, rightVars, varTypes, storage)
		network.BetaNodes[joinNodeID] = joinNode

		// Connect nodes
		if i == 0 {
			// First join: main type -> join node (left)
			cp.connectTypeNodeToBetaNode(network, ruleID, aggInfo.MainVariable, aggInfo.MainType, joinNode, "left")
			// Source type -> join node (right)
			cp.connectTypeNodeToBetaNode(network, ruleID, sourcePattern.Variable, sourcePattern.Type, joinNode, "right")
		} else {
			// Subsequent joins: previous join -> join node (left)
			if lastJoinNode != nil {
				lastJoinNode.AddChild(joinNode)
			}
			// Source type -> join node (right)
			cp.connectTypeNodeToBetaNode(network, ruleID, sourcePattern.Variable, sourcePattern.Type, joinNode, "right")
		}

		lastJoinNode = joinNode

		// Update mainTypeNode reference for logging
		_ = mainTypeNode
		_ = sourceTypeNode
	}

	// Create MultiSourceAccumulatorNode to compute aggregations
	accumulatorNode := NewMultiSourceAccumulatorNode(
		ruleID+"_msaccum",
		aggInfo.MainVariable,
		aggInfo.MainType,
		aggInfo.AggregationVars,
		aggInfo.SourcePatterns,
		storage,
	)
	network.BetaNodes[accumulatorNode.ID] = accumulatorNode

	// Connect the last join node to the accumulator
	if lastJoinNode != nil {
		lastJoinNode.AddChild(accumulatorNode)
		fmt.Printf("   ✓ JoinChain -> MultiSourceAccumulatorNode[%s]\n", accumulatorNode.ID)
	}

	fmt.Printf("   📊 MultiSourceAccumulatorNode créé avec %d agrégations\n", len(aggInfo.AggregationVars))
	for _, aggVar := range aggInfo.AggregationVars {
		thresholdStr := ""
		if aggVar.Operator != "" && (aggVar.Operator != ">=" || aggVar.Threshold != 0) {
			thresholdStr = fmt.Sprintf(" (threshold: %s %.2f)", aggVar.Operator, aggVar.Threshold)
		}
		fmt.Printf("     • %s: %s(%s.%s)%s\n",
			aggVar.Name, aggVar.Function, aggVar.SourceVar, aggVar.Field, thresholdStr)
	}

	// Create terminal node for action
	terminalNode := NewTerminalNode(ruleID+"_terminal", action, storage)
	terminalNode.SetNetwork(network)
	network.TerminalNodes[terminalNode.ID] = terminalNode

	// Register terminal node with lifecycle manager
	if network.LifecycleManager != nil {
		network.LifecycleManager.RegisterNode(terminalNode.ID, "terminal")
		network.LifecycleManager.AddRuleToNode(terminalNode.ID, ruleID, ruleID)
	}

	// Determine which accumulator type to use
	// Connect the accumulator to the terminal
	accumulatorNode.AddChild(terminalNode)
	fmt.Printf("   ✓ MultiSourceAccumulatorNode -> TerminalNode[%s]\n", terminalNode.ID)

	fmt.Printf("   ✅ Multi-source accumulator rule créée: %s\n", ruleID)
	return nil
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
	terminalNode.SetNetwork(network)
	network.TerminalNodes[terminalNode.ID] = terminalNode

	// Register terminal node with lifecycle manager
	if network.LifecycleManager != nil {
		network.LifecycleManager.RegisterNode(terminalNode.ID, "terminal")
		network.LifecycleManager.AddRuleToNode(terminalNode.ID, ruleID, ruleID)
	}

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
	terminalNode.SetNetwork(network)
	network.TerminalNodes[terminalNode.ID] = terminalNode

	// Register terminal node with lifecycle manager
	if network.LifecycleManager != nil {
		network.LifecycleManager.RegisterNode(terminalNode.ID, "terminal")
		network.LifecycleManager.AddRuleToNode(terminalNode.ID, ruleID, ruleID)
	}

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
	terminalNode.SetNetwork(network)
	network.TerminalNodes[terminalNode.ID] = terminalNode

	// Register terminal node with lifecycle manager
	if network.LifecycleManager != nil {
		network.LifecycleManager.RegisterNode(terminalNode.ID, "terminal")
		network.LifecycleManager.AddRuleToNode(terminalNode.ID, ruleID, ruleID)
	}

	// Extract variable name and type
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

	var joinNode *JoinNode
	var wasShared bool

	// Try to use BetaSharingRegistry if available and enabled
	if network.BetaSharingRegistry != nil && network.Config != nil && network.Config.BetaSharingEnabled {
		allVars := []string{variableNames[0], variableNames[1]}
		node, hash, shared, err := network.BetaSharingRegistry.GetOrCreateJoinNode(
			condition,
			leftVars,
			rightVars,
			allVars,
			varTypes,
			storage,
		)
		if err != nil {
			// Fallback to direct creation on error
			fmt.Printf("   ⚠️ Beta sharing failed: %v, falling back to direct creation\n", err)
			joinNode = NewJoinNode(ruleID+"_join", condition, leftVars, rightVars, varTypes, storage)
		} else {
			joinNode = node
			wasShared = shared
			if shared {
				fmt.Printf("   ♻️  Reused shared JoinNode %s (hash: %s)\n", joinNode.ID, hash)
			} else {
				fmt.Printf("   ✨ Created new shared JoinNode %s (hash: %s)\n", joinNode.ID, hash)
			}
			// Register with lifecycle manager
			if network.LifecycleManager != nil {
				network.LifecycleManager.RegisterNode(hash, "JoinNode")
			}
		}
	} else {
		// Legacy mode: direct creation
		joinNode = NewJoinNode(ruleID+"_join", condition, leftVars, rightVars, varTypes, storage)
	}

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

	if wasShared {
		fmt.Printf("   ✅ JoinNode %s réutilisé pour jointure %s\n", joinNode.ID, strings.Join(variableNames, " ⋈ "))
	} else {
		fmt.Printf("   ✅ JoinNode %s créé pour jointure %s\n", joinNode.ID, strings.Join(variableNames, " ⋈ "))
	}
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

	// Try to use BetaChainBuilder if available and enabled
	if network.BetaChainBuilder != nil && network.Config != nil && network.Config.BetaSharingEnabled {
		return cp.createCascadeJoinRuleWithBuilder(network, ruleID, variableNames, variableTypes, condition, terminalNode, storage)
	}

	// Fallback to legacy cascade implementation
	fmt.Printf("   🔧 Construction d'architecture en cascade de JoinNodes (legacy mode)\n")

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

// createCascadeJoinRuleWithBuilder creates a cascade using BetaChainBuilder with sharing support
func (cp *ConstraintPipeline) createCascadeJoinRuleWithBuilder(
	network *ReteNetwork,
	ruleID string,
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	terminalNode *TerminalNode,
	storage Storage,
) error {
	fmt.Printf("   🔧 Construction avec BetaChainBuilder (sharing enabled)\n")

	// Créer le mapping variable -> type
	varTypes := make(map[string]string)
	for i, varName := range variableNames {
		varTypes[varName] = variableTypes[i]
	}

	// Build join patterns for the chain
	patterns := make([]JoinPattern, 0, len(variableNames)-1)

	// Pattern 1: First two variables
	patterns = append(patterns, JoinPattern{
		LeftVars:    []string{variableNames[0]},
		RightVars:   []string{variableNames[1]},
		AllVars:     []string{variableNames[0], variableNames[1]},
		VarTypes:    varTypes,
		Condition:   condition,
		Selectivity: 0.5, // Default selectivity
	})

	// Patterns 2+: Each subsequent variable joins with accumulated results
	for i := 2; i < len(variableNames); i++ {
		accumulatedVars := make([]string, i)
		copy(accumulatedVars, variableNames[0:i])

		allVars := make([]string, i+1)
		copy(allVars, variableNames[0:i+1])

		patterns = append(patterns, JoinPattern{
			LeftVars:    accumulatedVars,
			RightVars:   []string{variableNames[i]},
			AllVars:     allVars,
			VarTypes:    varTypes,
			Condition:   condition,
			Selectivity: 0.5, // Default selectivity
		})
	}

	// Build the chain using BetaChainBuilder
	chain, err := network.BetaChainBuilder.BuildChain(patterns, ruleID)
	if err != nil {
		return fmt.Errorf("failed to build beta chain: %w", err)
	}

	// Add all nodes to network's BetaNodes map
	for _, node := range chain.Nodes {
		network.BetaNodes[node.ID] = node
	}

	// Connect type nodes to the first join node (for first two variables)
	if len(chain.Nodes) > 0 {
		firstJoin := chain.Nodes[0]
		for i := 0; i < 2 && i < len(variableNames); i++ {
			varName := variableNames[i]
			varType := variableTypes[i]
			side := NodeSideRight
			if i == 0 {
				side = NodeSideLeft
			}
			cp.connectTypeNodeToBetaNode(network, ruleID, varName, varType, firstJoin, side)
		}
		fmt.Printf("   ✓ Connected first two variables to initial JoinNode\n")
	}

	// Connect subsequent variables to their respective join nodes
	for i := 2; i < len(variableNames) && i-1 < len(chain.Nodes); i++ {
		joinNode := chain.Nodes[i-1]
		varName := variableNames[i]
		varType := variableTypes[i]
		cp.connectTypeNodeToBetaNode(network, ruleID, varName, varType, joinNode, NodeSideRight)
		fmt.Printf("   ✓ Connected variable %s to cascade level %d\n", varName, i)
	}

	// Connect the final node to the terminal
	if chain.FinalNode != nil {
		chain.FinalNode.AddChild(terminalNode)
		fmt.Printf("   ✅ Chain complete: %d JoinNodes, %d shared\n",
			len(chain.Nodes),
			network.BetaChainBuilder.GetMetrics().SharedJoinNodesReused)
	}

	fmt.Printf("   ✅ Architecture en cascade complète avec partage: %s\n", strings.Join(variableNames, " ⋈ "))

	return nil
}
