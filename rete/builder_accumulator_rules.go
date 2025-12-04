// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"github.com/treivax/tsd/tsdio"
	"fmt"
)

// AccumulatorRuleBuilder handles the creation of accumulator rules
type AccumulatorRuleBuilder struct {
	utils *BuilderUtils
}

// NewAccumulatorRuleBuilder creates a new AccumulatorRuleBuilder instance
func NewAccumulatorRuleBuilder(utils *BuilderUtils) *AccumulatorRuleBuilder {
	return &AccumulatorRuleBuilder{
		utils: utils,
	}
}

// CreateAccumulatorRule creates a rule with AccumulatorNode
func (arb *AccumulatorRuleBuilder) CreateAccumulatorRule(
	network *ReteNetwork,
	ruleID string,
	variables []map[string]interface{},
	variableNames []string,
	variableTypes []string,
	aggInfo *AggregationInfo,
	action *Action,
) error {
	// Extract the main variable and its type from variables
	if len(variables) == 0 || len(variableTypes) == 0 {
		return fmt.Errorf("aucune variable principale trouvée")
	}

	mainVariable := variableNames[0]
	mainType := variableTypes[0]

	// Store in aggInfo
	aggInfo.MainVariable = mainVariable
	aggInfo.MainType = mainType

	// Create the terminal node
	terminalNode := arb.utils.CreateTerminalNode(network, ruleID, action)

	// Create the comparison condition
	condition := map[string]interface{}{
		"type":     ConditionTypeComparison,
		"operator": aggInfo.Operator,
		"value":    aggInfo.Threshold,
	}

	// Create the AccumulatorNode with all parameters
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
		arb.utils.storage,
	)
	accumNode.AddChild(terminalNode)
	network.BetaNodes[accumNode.ID] = accumNode

	// Connect the TypeNodes to the AccumulatorNode
	arb.utils.ConnectTypeNodeToBetaNode(network, ruleID, mainVariable, mainType, accumNode, "")
	tsdio.Printf("   ✓ %s -> PassthroughAlpha -> AccumulatorNode[%s]\n", mainType, aggInfo.Function)

	arb.utils.ConnectTypeNodeToBetaNode(network, ruleID, aggInfo.AggVariable, aggInfo.AggType, accumNode, "")
	tsdio.Printf("   ✓ %s -> PassthroughAlpha -> AccumulatorNode[%s] (pour agrégation)\n", aggInfo.AggType, aggInfo.Function)

	tsdio.Printf("   ✅ AccumulatorNode %s créé pour %s(%s.%s) %s %.2f\n",
		accumNode.ID, aggInfo.Function, aggInfo.AggVariable, aggInfo.Field, aggInfo.Operator, aggInfo.Threshold)
	return nil
}

// IsMultiSourceAggregation checks if the rule has multiple aggregation sources
func (arb *AccumulatorRuleBuilder) IsMultiSourceAggregation(exprMap map[string]interface{}) bool {
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

// CreateMultiSourceAccumulatorRule creates a rule with multiple aggregation sources
// This is the refactored version - broken down into smaller, more manageable functions
func (arb *AccumulatorRuleBuilder) CreateMultiSourceAccumulatorRule(
	network *ReteNetwork,
	ruleID string,
	aggInfo *AggregationInfo,
	action *Action,
) error {
	tsdio.Printf("   🔗 Création règle multi-source avec %d sources et %d agrégations\n",
		len(aggInfo.SourcePatterns), len(aggInfo.AggregationVars))

	// Step 1: Create the join chain for all sources
	lastJoinNode, err := arb.createJoinChainForSources(network, ruleID, aggInfo)
	if err != nil {
		return err
	}

	// Step 2: Create the multi-source accumulator node
	accumulatorNode := arb.createMultiSourceAccumulatorNode(network, ruleID, aggInfo)

	// Step 3: Connect the join chain to the accumulator
	if lastJoinNode != nil {
		lastJoinNode.AddChild(accumulatorNode)
		tsdio.Printf("   ✓ JoinChain -> MultiSourceAccumulatorNode[%s]\n", accumulatorNode.ID)
	}

	// Step 4: Create and connect the terminal node
	return arb.connectAccumulatorToTerminal(network, ruleID, accumulatorNode, aggInfo, action)
}

// createJoinChainForSources creates a chain of join nodes for all source patterns
func (arb *AccumulatorRuleBuilder) createJoinChainForSources(
	network *ReteNetwork,
	ruleID string,
	aggInfo *AggregationInfo,
) (Node, error) {
	// Validate that the main type exists
	mainTypeNode, exists := network.TypeNodes[aggInfo.MainType]
	if !exists {
		return nil, fmt.Errorf("TypeNode pour %s non trouvé", aggInfo.MainType)
	}
	_ = mainTypeNode // Used for validation

	var lastJoinNode Node

	// Create a join node for each source pattern
	for i, sourcePattern := range aggInfo.SourcePatterns {
		// Validate source type node exists
		sourceTypeNode, exists := network.TypeNodes[sourcePattern.Type]
		if !exists {
			return nil, fmt.Errorf("TypeNode pour %s non trouvé", sourcePattern.Type)
		}
		_ = sourceTypeNode // Used for validation

		// Create the join node for this source
		joinNode, err := arb.createSourceJoinNode(network, ruleID, aggInfo, sourcePattern, i, lastJoinNode)
		if err != nil {
			return nil, err
		}

		// Connect to network
		arb.connectSourceJoinNode(network, ruleID, aggInfo, sourcePattern, joinNode, i, lastJoinNode)

		lastJoinNode = joinNode
	}

	return lastJoinNode, nil
}

// createSourceJoinNode creates a single join node for a source pattern
func (arb *AccumulatorRuleBuilder) createSourceJoinNode(
	network *ReteNetwork,
	ruleID string,
	aggInfo *AggregationInfo,
	sourcePattern SourcePattern,
	index int,
	lastJoinNode Node,
) (*JoinNode, error) {
	// Build condition map and variable lists
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

	// Build the join condition map
	if joinCondition != nil {
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

		tsdio.Printf("   ✓ Creating JoinNode: %s.%s == %s.%s\n",
			joinCondition.LeftVar, joinCondition.LeftField,
			joinCondition.RightVar, joinCondition.RightField)
	}

	// Build variable lists
	if index == 0 {
		// First join: main + first source
		leftVars = []string{aggInfo.MainVariable}
		varTypes[aggInfo.MainVariable] = aggInfo.MainType
	} else {
		// Subsequent joins: all previous variables on left side
		leftVars = []string{aggInfo.MainVariable}
		varTypes[aggInfo.MainVariable] = aggInfo.MainType
		for j := 0; j < index; j++ {
			leftVars = append(leftVars, aggInfo.SourcePatterns[j].Variable)
			varTypes[aggInfo.SourcePatterns[j].Variable] = aggInfo.SourcePatterns[j].Type
		}
	}
	rightVars = []string{sourcePattern.Variable}
	varTypes[sourcePattern.Variable] = sourcePattern.Type

	// Create the join node
	joinNodeID := fmt.Sprintf("%s_join_%d", ruleID, index)
	joinNode := NewJoinNode(joinNodeID, joinConditionMap, leftVars, rightVars, varTypes, arb.utils.storage)
	network.BetaNodes[joinNodeID] = joinNode

	return joinNode, nil
}

// connectSourceJoinNode connects a source join node to the network
func (arb *AccumulatorRuleBuilder) connectSourceJoinNode(
	network *ReteNetwork,
	ruleID string,
	aggInfo *AggregationInfo,
	sourcePattern SourcePattern,
	joinNode *JoinNode,
	index int,
	lastJoinNode Node,
) {
	if index == 0 {
		// First join: main type -> join node (left)
		arb.utils.ConnectTypeNodeToBetaNode(network, ruleID, aggInfo.MainVariable, aggInfo.MainType, joinNode, "left")
		// Source type -> join node (right)
		arb.utils.ConnectTypeNodeToBetaNode(network, ruleID, sourcePattern.Variable, sourcePattern.Type, joinNode, "right")
	} else {
		// Subsequent joins: previous join -> join node (left)
		if lastJoinNode != nil {
			lastJoinNode.AddChild(joinNode)
		}
		// Source type -> join node (right)
		arb.utils.ConnectTypeNodeToBetaNode(network, ruleID, sourcePattern.Variable, sourcePattern.Type, joinNode, "right")
	}
}

// createMultiSourceAccumulatorNode creates the accumulator node for multiple sources
func (arb *AccumulatorRuleBuilder) createMultiSourceAccumulatorNode(
	network *ReteNetwork,
	ruleID string,
	aggInfo *AggregationInfo,
) *MultiSourceAccumulatorNode {
	// Create MultiSourceAccumulatorNode to compute aggregations
	accumulatorNode := NewMultiSourceAccumulatorNode(
		ruleID+"_msaccum",
		aggInfo.MainVariable,
		aggInfo.MainType,
		aggInfo.AggregationVars,
		aggInfo.SourcePatterns,
		arb.utils.storage,
	)
	network.BetaNodes[accumulatorNode.ID] = accumulatorNode

	// Log aggregation details
	tsdio.Printf("   📊 MultiSourceAccumulatorNode créé avec %d agrégations\n", len(aggInfo.AggregationVars))
	for _, aggVar := range aggInfo.AggregationVars {
		thresholdStr := ""
		if aggVar.Operator != "" && (aggVar.Operator != ">=" || aggVar.Threshold != 0) {
			thresholdStr = fmt.Sprintf(" (threshold: %s %.2f)", aggVar.Operator, aggVar.Threshold)
		}
		tsdio.Printf("     • %s: %s(%s.%s)%s\n",
			aggVar.Name, aggVar.Function, aggVar.SourceVar, aggVar.Field, thresholdStr)
	}

	return accumulatorNode
}

// connectAccumulatorToTerminal creates and connects the terminal node to the accumulator
func (arb *AccumulatorRuleBuilder) connectAccumulatorToTerminal(
	network *ReteNetwork,
	ruleID string,
	accumulatorNode *MultiSourceAccumulatorNode,
	aggInfo *AggregationInfo,
	action *Action,
) error {
	// Create terminal node for action
	terminalNode := arb.utils.CreateTerminalNode(network, ruleID, action)

	// Connect the accumulator to the terminal
	accumulatorNode.AddChild(terminalNode)
	tsdio.Printf("   ✓ MultiSourceAccumulatorNode -> TerminalNode[%s]\n", terminalNode.ID)

	tsdio.Printf("   ✅ Multi-source accumulator rule créée: %s\n", ruleID)
	return nil
}
