// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"

)

// ExistsRuleBuilder handles the creation of EXISTS rules
type ExistsRuleBuilder struct {
	utils *BuilderUtils
}

// NewExistsRuleBuilder creates a new ExistsRuleBuilder instance
func NewExistsRuleBuilder(utils *BuilderUtils) *ExistsRuleBuilder {
	return &ExistsRuleBuilder{
		utils: utils,
	}
}

// CreateExistsRule creates an EXISTS rule with ExistsNode
func (erb *ExistsRuleBuilder) CreateExistsRule(
	network *ReteNetwork,
	ruleID string,
	exprMap map[string]interface{},
	condition map[string]interface{},
	action *Action,
) error {
	// Create the terminal node for this rule
	terminalNode := erb.utils.CreateTerminalNode(network, ruleID, action)

	// Extract variables
	mainVariable, existsVariable, mainVarType, existsVarType, err := erb.ExtractExistsVariables(exprMap)
	if err != nil {
		return err
	}

	// Extract EXISTS conditions
	existsConditions, err := erb.ExtractExistsConditions(exprMap)
	if err != nil {
		return err
	}

	// Create the condition object for the ExistsNode
	existsConditionObj := map[string]interface{}{
		"type":       ConditionTypeExists,
		"conditions": existsConditions,
	}

	// Create the variable -> type mapping for the ExistsNode
	varTypes := make(map[string]string)
	varTypes[mainVariable] = mainVarType
	varTypes[existsVariable] = existsVarType

	// Create the ExistsNode with the actual conditions
	existsNode := NewExistsNode(ruleID+"_exists", existsConditionObj, mainVariable, existsVariable, varTypes, erb.utils.storage)
	existsNode.AddChild(terminalNode)

	// Store the ExistsNode in the network's BetaNodes
	network.BetaNodes[existsNode.ID] = existsNode

	// Create pass-through AlphaNodes for both variables
	erb.ConnectExistsNodeToTypeNodes(network, ruleID, existsNode, mainVariable, mainVarType, existsVariable, existsVarType)

	fmt.Printf("   ✅ ExistsNode %s créé pour %s EXISTS %s\n", existsNode.ID, mainVariable, existsVariable)
	return nil
}

// ExtractExistsVariables extracts variables from an EXISTS rule
func (erb *ExistsRuleBuilder) ExtractExistsVariables(exprMap map[string]interface{}) (string, string, string, string, error) {
	var mainVariable, existsVariable string
	var mainVarType, existsVarType string

	// Extract the main variable from "set"
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

	// Extract the existence variable from constraints
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

// ExtractExistsConditions extracts conditions from an EXISTS rule
func (erb *ExistsRuleBuilder) ExtractExistsConditions(exprMap map[string]interface{}) ([]map[string]interface{}, error) {
	var existsConditions []map[string]interface{}

	if constraintsData, hasConstraints := exprMap["constraints"]; hasConstraints {
		if constraintMap, ok := constraintsData.(map[string]interface{}); ok {
			// Try "condition" (singular) first
			if conditionData, hasCondition := constraintMap["condition"]; hasCondition {
				if conditionObj, ok := conditionData.(map[string]interface{}); ok {
					existsConditions = append(existsConditions, conditionObj)
				}
			}
			// Then try "conditions" (plural) if not found
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

// ConnectExistsNodeToTypeNodes connects an ExistsNode to the appropriate TypeNodes
func (erb *ExistsRuleBuilder) ConnectExistsNodeToTypeNodes(
	network *ReteNetwork,
	ruleID string,
	existsNode *ExistsNode,
	mainVariable string,
	mainVarType string,
	existsVariable string,
	existsVarType string,
) {
	// Connect the main and existence variables to the ExistsNode
	if mainVarType != "" {
		erb.utils.ConnectTypeNodeToBetaNode(network, ruleID, mainVariable, mainVarType, existsNode, NodeSideLeft)
	}
	if existsVarType != "" {
		erb.utils.ConnectTypeNodeToBetaNode(network, ruleID, existsVariable, existsVarType, existsNode, NodeSideRight)
	}
}
