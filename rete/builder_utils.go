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

// PassthroughNodeKey generates a registry key for a passthrough node
// The key includes rule name to prevent incorrect sharing when rules have different alpha filters
func PassthroughNodeKey(ruleName, typeName, varName, side string) string {
	if side != "" {
		return fmt.Sprintf("passthrough_%s_%s_%s_%s", ruleName, varName, typeName, side)
	}
	return fmt.Sprintf("passthrough_%s_%s_%s", ruleName, varName, typeName)
}

// CreatePassthroughAlphaNode creates a passthrough AlphaNode with optional side specification
// DEPRECATED: Use GetOrCreatePassthroughAlphaNode instead for proper sharing
func (bu *BuilderUtils) CreatePassthroughAlphaNode(ruleID, varName, side string) *AlphaNode {
	passCondition := map[string]interface{}{
		"type": ConditionTypePassthrough,
	}
	if side != "" {
		passCondition["side"] = side
	}
	return NewAlphaNode(ruleID+"_pass_"+varName, passCondition, varName, bu.storage)
}

// GetOrCreatePassthroughAlphaNode gets or creates a passthrough AlphaNode per rule
// Each rule gets its own passthrough to prevent incorrect sharing when alpha filters differ
func (bu *BuilderUtils) GetOrCreatePassthroughAlphaNode(
	network *ReteNetwork,
	ruleName string,
	typeName string,
	varName string,
	side string,
) *AlphaNode {
	// Generate registry key based on rule, variable, type and side
	key := PassthroughNodeKey(ruleName, typeName, varName, side)

	// Check if passthrough already exists in registry
	if existingNode, exists := network.PassthroughRegistry[key]; exists {
		return existingNode // ✅ Reuse existing node
	}

	// Create new passthrough node
	passCondition := map[string]interface{}{
		"type": ConditionTypePassthrough,
	}
	if side != "" {
		passCondition["side"] = side
	}

	alphaNode := NewAlphaNode(key, passCondition, varName, bu.storage)

	// Register for future reuse
	network.PassthroughRegistry[key] = alphaNode

	return alphaNode
}

// ConnectTypeNodeToBetaNode connects a TypeNode to a BetaNode via a passthrough AlphaNode
// Creates per-rule passthrough nodes to prevent incorrect sharing
func (bu *BuilderUtils) ConnectTypeNodeToBetaNode(
	network *ReteNetwork,
	ruleID string,
	varName string,
	varType string,
	betaNode Node,
	side string,
) {
	if typeNode, exists := network.TypeNodes[varType]; exists {
		// Get or create per-rule passthrough node
		alphaNode := bu.GetOrCreatePassthroughAlphaNode(network, ruleID, varType, varName, side)

		// Connect TypeNode -> AlphaNode (if not already connected)
		if !bu.hasChild(typeNode, alphaNode) {
			typeNode.AddChild(alphaNode)
		}

		// Connect AlphaNode -> BetaNode
		alphaNode.AddChild(betaNode)

		sideInfo := ""
		if side != "" {
			sideInfo = fmt.Sprintf(" (%s)", strings.ToUpper(side))
		}
		fmt.Printf("   ✓ %s -> PassthroughAlpha[%s] -> %s%s\n", varType, alphaNode.GetID(), betaNode.GetID(), sideInfo)
	}
}

// hasChild checks if a node already has a specific child
func (bu *BuilderUtils) hasChild(parent Node, child Node) bool {
	for _, c := range parent.GetChildren() {
		if c.GetID() == child.GetID() {
			return true
		}
	}
	return false
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
