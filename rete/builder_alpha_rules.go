// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
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

	// Create the AlphaNode with its terminal
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

// createAlphaNodeWithTerminal creates an AlphaNode with a terminal node using AlphaSharingManager
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
