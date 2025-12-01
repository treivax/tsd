// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"log"
)

// AlphaRuleBuilder handles the creation of alpha rules
type AlphaRuleBuilder struct {
	utils *BuilderUtils
}

// NewAlphaRuleBuilder creates a new AlphaRuleBuilder instance
func NewAlphaRuleBuilder(utils *BuilderUtils) *AlphaRuleBuilder {
	return &AlphaRuleBuilder{
		utils: utils,
	}
}

// CreateAlphaRule creates a simple alpha rule with a single variable
// This method now checks if the condition can be decomposed into a chain
func (arb *AlphaRuleBuilder) CreateAlphaRule(
	network *ReteNetwork,
	ruleID string,
	variables []map[string]interface{},
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	action *Action,
) error {
	// Extract variable information
	variableName, variableType := arb.getVariableInfo(variables, variableTypes)

	// Check if this condition should be built as a chain
	shouldUseChain, actualCondition := arb.shouldBuildAsChain(condition)

	if shouldUseChain {
		log.Printf("   🔗 Multi-condition AND detected for rule %s, using AlphaChainBuilder", ruleID)
		return arb.createAlphaChainWithTerminal(
			network,
			ruleID,
			actualCondition,
			variableName,
			variableType,
			action,
		)
	}

	// Create a simple AlphaNode with its terminal
	return arb.createAlphaNodeWithTerminal(
		network,
		ruleID,
		condition,
		variableName,
		variableType,
		action,
	)
}

// getVariableInfo extracts the first variable name and type
func (arb *AlphaRuleBuilder) getVariableInfo(
	variables []map[string]interface{},
	variableTypes []string,
) (string, string) {
	var variableName, variableType string

	if len(variables) > 0 {
		if name, ok := variables[0]["name"].(string); ok {
			variableName = name
		}
	}

	if len(variableTypes) > 0 {
		variableType = variableTypes[0]
	}

	return variableName, variableType
}

// shouldBuildAsChain determines if a condition should be built as a chain
// Returns (shouldUseChain bool, actualCondition interface{})
func (arb *AlphaRuleBuilder) shouldBuildAsChain(condition map[string]interface{}) (bool, interface{}) {
	// Unwrap the condition if it's wrapped in a map
	actualCondition := interface{}(condition)

	if condType, hasType := condition["type"]; hasType {
		switch condType {
		case "constraint":
			if constraint, hasConstraint := condition["constraint"]; hasConstraint {
				actualCondition = constraint
			}
		case "negation", "simple", "passthrough":
			// These types don't use chains
			return false, actualCondition
		}
	}

	// Analyze the expression to determine its type
	exprType, err := AnalyzeExpression(actualCondition)
	if err != nil {
		// If we can't analyze it, use simple behavior
		return false, actualCondition
	}

	// Only decompose AND expressions
	if exprType != ExprTypeAND {
		return false, actualCondition
	}

	// Check if it can be decomposed
	if !CanDecompose(exprType) {
		return false, actualCondition
	}

	// Extract conditions to verify we have multiple
	conditions, _, err := ExtractConditions(actualCondition)
	if err != nil || len(conditions) <= 1 {
		return false, actualCondition
	}

	// This should be built as a chain
	return true, actualCondition
}

// createAlphaChainWithTerminal creates a chain of AlphaNodes for multi-condition AND rules
func (arb *AlphaRuleBuilder) createAlphaChainWithTerminal(
	network *ReteNetwork,
	ruleID string,
	condition interface{},
	variableName string,
	variableType string,
	action *Action,
) error {
	// Extract the conditions from the AND expression
	conditions, opType, err := ExtractConditions(condition)
	if err != nil {
		return fmt.Errorf("erreur extraction conditions: %w", err)
	}

	log.Printf("   🔗 Décomposition en chaîne: %d conditions détectées (opérateur: %s)", len(conditions), opType)

	// Normalize the conditions
	normalizedConditions := NormalizeConditions(conditions, opType)
	log.Printf("   📋 Conditions normalisées: %d condition(s)", len(normalizedConditions))

	// Find the TypeNode parent to connect the chain
	var parentNode Node
	if variableType != "" {
		if typeNode, exists := network.TypeNodes[variableType]; exists {
			parentNode = typeNode
		}
	}

	// If no TypeNode found, use the first available
	if parentNode == nil {
		for _, typeNode := range network.TypeNodes {
			parentNode = typeNode
			break
		}
	}

	if parentNode == nil {
		return fmt.Errorf("aucun TypeNode trouvé pour connecter la chaîne")
	}

	// Create the chain builder
	chainBuilder := NewAlphaChainBuilder(network, arb.utils.storage)

	// Build the chain of AlphaNodes
	chain, err := chainBuilder.BuildChain(normalizedConditions, variableName, parentNode, ruleID)
	if err != nil {
		return fmt.Errorf("erreur construction chaîne: %w", err)
	}

	// Validate the chain
	if err := chain.ValidateChain(); err != nil {
		return fmt.Errorf("chaîne invalide: %w", err)
	}

	// Get chain statistics
	stats := chainBuilder.GetChainStats(chain)
	sharedCount := 0
	if sc, ok := stats["shared_nodes"].(int); ok {
		sharedCount = sc
	}

	// Display construction statistics
	log.Printf("   ✅ Chaîne construite: %d nœud(s), %d partagé(s)", len(chain.Nodes), sharedCount)

	// Log details for each node
	for i, node := range chain.Nodes {
		if i < sharedCount {
			log.Printf("   ♻️  AlphaNode partagé réutilisé: %s (hash: %s)", node.ID, chain.Hashes[i])
		} else {
			log.Printf("   ✨ Nouveau AlphaNode créé: %s (hash: %s)", node.ID, chain.Hashes[i])
		}
	}

	// Create and attach the terminal to the last node of the chain
	terminalNode := arb.utils.CreateTerminalNode(network, ruleID, action)
	chain.FinalNode.AddChild(terminalNode)

	// Register terminal node with lifecycle manager
	if network.LifecycleManager != nil {
		network.LifecycleManager.RegisterNode(terminalNode.ID, "terminal")
		network.LifecycleManager.AddRuleToNode(terminalNode.ID, ruleID, ruleID)
	}

	log.Printf("   ✓ TerminalNode %s attaché au nœud final %s de la chaîne", terminalNode.ID, chain.FinalNode.ID)
	fmt.Printf("   ✓ Règle alpha avec chaîne créée pour: %s\n", ruleID)

	return nil
}

// createAlphaNodeWithTerminal creates a simple AlphaNode with a terminal node using AlphaSharingManager
func (arb *AlphaRuleBuilder) createAlphaNodeWithTerminal(
	network *ReteNetwork,
	ruleID string,
	condition map[string]interface{},
	variableName string,
	variableType string,
	action *Action,
) error {
	// Verify AlphaSharingManager is initialized
	if network.AlphaSharingManager == nil {
		return fmt.Errorf("AlphaSharingManager non initialisé dans le réseau")
	}

	// Get the TypeNode
	typeNode, exists := network.TypeNodes[variableType]
	if !exists {
		return fmt.Errorf("TypeNode pour %s non trouvé", variableType)
	}

	// Use AlphaSharingManager to get or create the AlphaNode (for sharing)
	alphaNode, alphaHash, wasShared, err := network.AlphaSharingManager.GetOrCreateAlphaNode(
		condition,
		variableName,
		arb.utils.storage,
	)
	if err != nil {
		return fmt.Errorf("erreur création AlphaNode partagé: %w", err)
	}

	if wasShared {
		fmt.Printf("   ♻️  AlphaNode partagé réutilisé: %s (hash: %s)\n", alphaNode.ID, alphaHash)
	} else {
		fmt.Printf("   ✨ Nouveau AlphaNode partageable créé: %s (hash: %s)\n", alphaNode.ID, alphaHash)

		// Connect TypeNode -> AlphaNode (only for new nodes)
		typeNode.AddChild(alphaNode)

		// Store the AlphaNode in the network's global registry
		network.AlphaNodes[alphaNode.ID] = alphaNode
	}

	// Register or update the AlphaNode in LifecycleManager
	if network.LifecycleManager != nil {
		lifecycle := network.LifecycleManager.RegisterNode(alphaNode.ID, "alpha")
		lifecycle.AddRuleReference(ruleID, ruleID)
	}

	// Create the terminal node (always rule-specific)
	terminalNode := arb.utils.CreateTerminalNode(network, ruleID, action)

	// Connect AlphaNode -> TerminalNode
	alphaNode.AddChild(terminalNode)

	// Register terminal node with lifecycle manager
	if network.LifecycleManager != nil {
		network.LifecycleManager.RegisterNode(terminalNode.ID, "terminal")
		network.LifecycleManager.AddRuleToNode(terminalNode.ID, ruleID, ruleID)
	}

	fmt.Printf("   ✓ Règle alpha simple créée pour: %s\n", ruleID)
	fmt.Printf("   ✓ %s -> AlphaNode[%s] -> TerminalNode[%s]\n",
		variableType, alphaNode.ID, terminalNode.ID)

	return nil
}
