// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
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
// Délègue au TypeBuilder pour la création effective
func (cp *ConstraintPipeline) createTypeNodes(network *ReteNetwork, types []interface{}, storage Storage) error {
	utils := NewBuilderUtils(storage)
	typeBuilder := NewTypeBuilder(utils)
	return typeBuilder.CreateTypeNodes(network, types, storage)
}

// createTypeDefinition crée une définition de type à partir d'une map
// Wrapper legacy pour compatibilité
func (cp *ConstraintPipeline) createTypeDefinition(typeName string, typeMap map[string]interface{}) TypeDefinition {
	utils := NewBuilderUtils(nil)
	typeBuilder := NewTypeBuilder(utils)
	return typeBuilder.CreateTypeDefinition(typeName, typeMap)
}

// createRuleNodes crée les nœuds de règles (Alpha, Beta, Terminal) à partir des expressions
// Délègue au RuleBuilder pour l'orchestration
func (cp *ConstraintPipeline) createRuleNodes(network *ReteNetwork, expressions []interface{}, storage Storage) error {
	utils := NewBuilderUtils(storage)
	ruleBuilder := NewRuleBuilder(utils, cp)

	return ruleBuilder.CreateRuleNodes(network, expressions)
}

// createSingleRule crée une règle unique
// Délègue au RuleBuilder
func (cp *ConstraintPipeline) createSingleRule(network *ReteNetwork, ruleID string, exprMap map[string]interface{}, storage Storage) error {
	utils := NewBuilderUtils(storage)
	ruleBuilder := NewRuleBuilder(utils, cp)

	return ruleBuilder.CreateSingleRule(network, ruleID, exprMap)
}

// isMultiSourceAggregation checks if the rule has multiple aggregation sources
// Wrapper pour compatibilité
func (cp *ConstraintPipeline) isMultiSourceAggregation(exprMap map[string]interface{}) bool {
	utils := NewBuilderUtils(nil)
	accumulatorBuilder := NewAccumulatorRuleBuilder(utils)
	return accumulatorBuilder.IsMultiSourceAggregation(exprMap)
}

// createMultiSourceAccumulatorRule creates a rule with multiple aggregation sources
// Délègue au AccumulatorRuleBuilder
func (cp *ConstraintPipeline) createMultiSourceAccumulatorRule(
	network *ReteNetwork,
	ruleID string,
	aggInfo *AggregationInfo,
	action *Action,
	storage Storage,
) error {
	utils := NewBuilderUtils(storage)
	accumulatorBuilder := NewAccumulatorRuleBuilder(utils)
	return accumulatorBuilder.CreateMultiSourceAccumulatorRule(network, ruleID, aggInfo, action)
}

// createAlphaRule crée une règle alpha simple avec une seule variable
// Délègue au AlphaRuleBuilder
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
	utils := NewBuilderUtils(storage)
	alphaBuilder := NewAlphaRuleBuilder(utils)
	return alphaBuilder.CreateAlphaRule(network, ruleID, variables, variableNames, variableTypes, condition, action)
}

// createJoinRule crée une règle de jointure avec JoinNode
// Délègue au JoinRuleBuilder
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
	utils := NewBuilderUtils(storage)
	joinBuilder := NewJoinRuleBuilder(utils)
	return joinBuilder.CreateJoinRule(network, ruleID, variableNames, variableTypes, condition, action)
}

// createExistsRule crée une règle EXISTS avec ExistsNode
// Délègue au ExistsRuleBuilder
func (cp *ConstraintPipeline) createExistsRule(
	network *ReteNetwork,
	ruleID string,
	exprMap map[string]interface{},
	condition map[string]interface{},
	action *Action,
	storage Storage,
) error {
	utils := NewBuilderUtils(storage)
	existsBuilder := NewExistsRuleBuilder(utils)
	return existsBuilder.CreateExistsRule(network, ruleID, exprMap, condition, action)
}

// extractExistsVariables extrait les variables d'une règle EXISTS
// Délègue au ExistsRuleBuilder
func (cp *ConstraintPipeline) extractExistsVariables(exprMap map[string]interface{}) (string, string, string, string, error) {
	utils := NewBuilderUtils(nil)
	existsBuilder := NewExistsRuleBuilder(utils)
	return existsBuilder.ExtractExistsVariables(exprMap)
}

// extractExistsConditions extrait les conditions d'une règle EXISTS
// Délègue au ExistsRuleBuilder
func (cp *ConstraintPipeline) extractExistsConditions(exprMap map[string]interface{}) ([]map[string]interface{}, error) {
	utils := NewBuilderUtils(nil)
	existsBuilder := NewExistsRuleBuilder(utils)
	return existsBuilder.ExtractExistsConditions(exprMap)
}

// connectExistsNodeToTypeNodes connecte un ExistsNode aux TypeNodes appropriés
// Délègue au ExistsRuleBuilder
func (cp *ConstraintPipeline) connectExistsNodeToTypeNodes(
	network *ReteNetwork,
	ruleID string,
	existsNode *ExistsNode,
	mainVariable string,
	mainVarType string,
	existsVariable string,
	existsVarType string,
) {
	utils := NewBuilderUtils(network.Storage)
	existsBuilder := NewExistsRuleBuilder(utils)
	existsBuilder.ConnectExistsNodeToTypeNodes(network, ruleID, existsNode, mainVariable, mainVarType, existsVariable, existsVarType)
}

// createAccumulatorRule crée une règle avec AccumulatorNode
// Délègue au AccumulatorRuleBuilder
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
	utils := NewBuilderUtils(storage)
	accumulatorBuilder := NewAccumulatorRuleBuilder(utils)
	return accumulatorBuilder.CreateAccumulatorRule(network, ruleID, variables, variableNames, variableTypes, aggInfo, action)
}

// createPassthroughAlphaNode creates a passthrough AlphaNode with optional side specification
// Délègue au BuilderUtils
func (cp *ConstraintPipeline) createPassthroughAlphaNode(ruleID, varName, side string, storage Storage) *AlphaNode {
	utils := NewBuilderUtils(storage)
	return utils.CreatePassthroughAlphaNode(ruleID, varName, side)
}

// connectTypeNodeToBetaNode connects a TypeNode to a BetaNode via a passthrough AlphaNode
// Délègue au BuilderUtils
func (cp *ConstraintPipeline) connectTypeNodeToBetaNode(
	network *ReteNetwork,
	ruleID string,
	varName string,
	varType string,
	betaNode Node,
	side string,
) {
	utils := NewBuilderUtils(network.Storage)
	utils.ConnectTypeNodeToBetaNode(network, ruleID, varName, varType, betaNode, side)
}
