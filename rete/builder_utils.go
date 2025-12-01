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

// BuilderUtils provides common utility functions for all builders
type BuilderUtils struct {
	storage Storage
}

// NewBuilderUtils creates a new BuilderUtils instance
func NewBuilderUtils(storage Storage) *BuilderUtils {
	return &BuilderUtils{
		storage: storage,
	}
}

// CreatePassthroughAlphaNode creates a passthrough AlphaNode with optional side specification
func (bu *BuilderUtils) CreatePassthroughAlphaNode(ruleID, varName, side string) *AlphaNode {
	passCondition := map[string]interface{}{
		"type": ConditionTypePassthrough,
	}
	if side != "" {
		passCondition["side"] = side
	}
	return NewAlphaNode(ruleID+"_pass_"+varName, passCondition, varName, bu.storage)
}

// ConnectTypeNodeToBetaNode connects a TypeNode to a BetaNode via a passthrough AlphaNode
func (bu *BuilderUtils) ConnectTypeNodeToBetaNode(
	network *ReteNetwork,
	ruleID string,
	varName string,
	varType string,
	betaNode Node,
	side string,
) {
	if typeNode, exists := network.TypeNodes[varType]; exists {
		alphaNode := bu.CreatePassthroughAlphaNode(ruleID, varName, side)
		typeNode.AddChild(alphaNode)
		alphaNode.AddChild(betaNode)

		sideInfo := ""
		if side != "" {
			sideInfo = fmt.Sprintf(" (%s)", strings.ToUpper(side))
		}
		fmt.Printf("   ✓ %s -> PassthroughAlpha_%s -> %s%s\n", varType, varName, betaNode.GetID(), sideInfo)
	}
}

// GetStringField safely extracts a string field from a map
func GetStringField(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

// GetIntField safely extracts an int field from a map
func GetIntField(m map[string]interface{}, key string, defaultValue int) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		}
	}
	return defaultValue
}

// GetBoolField safely extracts a bool field from a map
func GetBoolField(m map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := m[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return defaultValue
}

// GetMapField safely extracts a map field from a map
func GetMapField(m map[string]interface{}, key string) (map[string]interface{}, bool) {
	if val, ok := m[key]; ok {
		if mapVal, ok := val.(map[string]interface{}); ok {
			return mapVal, true
		}
	}
	return nil, false
}

// GetListField safely extracts a list field from a map
func GetListField(m map[string]interface{}, key string) ([]interface{}, bool) {
	if val, ok := m[key]; ok {
		if listVal, ok := val.([]interface{}); ok {
			return listVal, true
		}
	}
	return nil, false
}

// CreateTerminalNode creates a terminal node and registers it with the network
func (bu *BuilderUtils) CreateTerminalNode(
	network *ReteNetwork,
	ruleID string,
	action *Action,
) *TerminalNode {
	terminalNode := NewTerminalNode(ruleID+"_terminal", action, bu.storage)
	terminalNode.SetNetwork(network)
	network.TerminalNodes[terminalNode.ID] = terminalNode

	// Register terminal node with lifecycle manager
	if network.LifecycleManager != nil {
		network.LifecycleManager.RegisterNode(terminalNode.ID, "terminal")
		network.LifecycleManager.AddRuleToNode(terminalNode.ID, ruleID, ruleID)
	}

	return terminalNode
}

// BuildVarTypesMap creates a map of variable names to their types
func BuildVarTypesMap(variableNames, variableTypes []string) map[string]string {
	varTypes := make(map[string]string)
	for i, varName := range variableNames {
		if i < len(variableTypes) {
			varTypes[varName] = variableTypes[i]
		}
	}
	return varTypes
}
