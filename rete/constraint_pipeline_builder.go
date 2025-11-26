package rete

import (
	"fmt"
	"strings"
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
			"type": "simple",
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

	// Créer le JoinNode
	leftVars := []string{variableNames[0]} // Variable primaire
	rightVars := variableNames[1:]         // Variables secondaires

	// Créer le mapping variable -> type
	varTypes := make(map[string]string)
	for i, varName := range variableNames {
		varTypes[varName] = variableTypes[i]
	}

	joinNode := NewJoinNode(ruleID+"_join", condition, leftVars, rightVars, varTypes, storage)
	joinNode.AddChild(terminalNode)

	// Stocker le JoinNode dans les BetaNodes du réseau
	network.BetaNodes[joinNode.ID] = joinNode

	// Créer des AlphaNodes pass-through qui ne filtrent pas mais transfèrent vers JoinNode
	for i, varName := range variableNames {
		varType := variableTypes[i]
		if varType != "" {
			if typeNode, exists := network.TypeNodes[varType]; exists {
				// Déterminer le côté (gauche/droite) selon l'architecture RETE
				side := "right"
				if i == 0 {
					side = "left" // Première variable va vers la gauche
				}

				// Créer un AlphaNode pass-through avec indication de côté
				passCondition := map[string]interface{}{
					"type": "passthrough",
					"side": side,
				}
				alphaNode := NewAlphaNode(ruleID+"_pass_"+varName, passCondition, varName, storage)

				// Connecter TypeNode -> AlphaPassthrough -> JoinNode
				typeNode.AddChild(alphaNode)
				alphaNode.AddChild(joinNode)

				fmt.Printf("   ✓ %s -> PassthroughAlpha_%s -> JoinNode_%s\n", varType, varName, ruleID)
			} else {
				fmt.Printf("   ⚠️ TypeNode %s introuvable!\n", varType)
			}
		} else {
			fmt.Printf("   ⚠️ Type vide pour variable %s\n", varName)
		}
	}

	fmt.Printf("   ✅ JoinNode %s créé pour jointure %s\n", joinNode.ID, strings.Join(variableNames, " ⋈ "))
	return nil
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
		"type":       "exists",
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
	// Variable principale → ActivateLeft
	if mainVarType != "" {
		if typeNode, exists := network.TypeNodes[mainVarType]; exists {
			mainAlphaCondition := map[string]interface{}{
				"type": "passthrough",
				"side": "left",
			}
			mainAlphaNode := NewAlphaNode(ruleID+"_pass_"+mainVariable, mainAlphaCondition, mainVariable, network.Storage)

			typeNode.AddChild(mainAlphaNode)
			mainAlphaNode.AddChild(existsNode)

			fmt.Printf("   ✓ %s -> PassthroughAlpha_%s -> ExistsNode_%s (LEFT)\n", mainVarType, mainVariable, ruleID)
		}
	}

	// Variable d'existence → ActivateRight
	if existsVarType != "" {
		if typeNode, exists := network.TypeNodes[existsVarType]; exists {
			existsAlphaCondition := map[string]interface{}{
				"type": "passthrough",
				"side": "right",
			}
			existsAlphaNode := NewAlphaNode(ruleID+"_pass_"+existsVariable, existsAlphaCondition, existsVariable, network.Storage)

			typeNode.AddChild(existsAlphaNode)
			existsAlphaNode.AddChild(existsNode)

			fmt.Printf("   ✓ %s -> PassthroughAlpha_%s -> ExistsNode_%s (RIGHT)\n", existsVarType, existsVariable, ruleID)
		}
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
		"type":     "comparison",
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

	// Connecter le TypeNode principal (Employee) à l'AccumulatorNode
	if typeNode, exists := network.TypeNodes[mainType]; exists {
		// Créer un AlphaNode passthrough pour la variable principale
		passCondition := map[string]interface{}{
			"type": "passthrough",
		}
		alphaNode := NewAlphaNode(ruleID+"_pass_"+mainVariable, passCondition, mainVariable, storage)

		typeNode.AddChild(alphaNode)
		alphaNode.AddChild(accumNode)

		fmt.Printf("   ✓ %s -> PassthroughAlpha -> AccumulatorNode[%s]\n", mainType, aggInfo.Function)
	}

	// CRUCIAL: Connecter aussi le TypeNode des faits à agréger (Performance) à l'AccumulatorNode
	if aggTypeNode, exists := network.TypeNodes[aggInfo.AggType]; exists {
		// Créer un AlphaNode passthrough pour la variable d'agrégation
		passConditionAgg := map[string]interface{}{
			"type": "passthrough",
		}
		alphaNodeAgg := NewAlphaNode(ruleID+"_pass_"+aggInfo.AggVariable, passConditionAgg, aggInfo.AggVariable, storage)

		aggTypeNode.AddChild(alphaNodeAgg)
		alphaNodeAgg.AddChild(accumNode)

		fmt.Printf("   ✓ %s -> PassthroughAlpha -> AccumulatorNode[%s] (pour agrégation)\n", aggInfo.AggType, aggInfo.Function)
	}

	fmt.Printf("   ✅ AccumulatorNode %s créé pour %s(%s.%s) %s %.2f\n",
		accumNode.ID, aggInfo.Function, aggInfo.AggVariable, aggInfo.Field, aggInfo.Operator, aggInfo.Threshold)
	return nil
}
