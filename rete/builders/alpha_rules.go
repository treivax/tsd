// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package builders

import (
	"fmt"

	"github.com/treivax/tsd/rete"
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
	network *rete.ReteNetwork,
	ruleID string,
	variables []map[string]interface{},
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	action *rete.Action,
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

// createAlphaNodeWithTerminal creates an AlphaNode with a terminal node
func (arb *AlphaRuleBuilder) createAlphaNodeWithTerminal(
	network *rete.ReteNetwork,
	ruleID string,
	condition map[string]interface{},
	variableName string,
	variableType string,
	action *rete.Action,
) error {
	// Create the terminal node
	terminalNode := arb.utils.CreateTerminalNode(network, ruleID, action)

	// Get the TypeNode
	typeNode, exists := network.TypeNodes[variableType]
	if !exists {
		return fmt.Errorf("TypeNode pour %s non trouvé", variableType)
	}

	// Create the AlphaNode
	alphaNode := rete.NewAlphaNode(ruleID+"_alpha", condition, variableName, arb.utils.storage)
	alphaNode.AddChild(terminalNode)

	// Connect TypeNode -> AlphaNode
	typeNode.AddChild(alphaNode)

	// Store the AlphaNode in the network
	network.AlphaNodes[alphaNode.ID] = alphaNode

	fmt.Printf("   ✓ %s -> AlphaNode[%s] -> TerminalNode[%s]\n",
		variableType, alphaNode.ID, terminalNode.ID)

	return nil
}
