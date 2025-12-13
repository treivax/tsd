// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

// RuleBuilder orchestrates the creation of all types of rules
type RuleBuilder struct {
	alphaBuilder       *AlphaRuleBuilder
	joinBuilder        *JoinRuleBuilder
	existsBuilder      *ExistsRuleBuilder
	accumulatorBuilder *AccumulatorRuleBuilder
	utils              *BuilderUtils

	// Pipeline helper for accessing constraint pipeline functionality
	pipeline PipelineHelper
}

// NewRuleBuilder creates a new RuleBuilder instance
func NewRuleBuilder(utils *BuilderUtils, pipeline PipelineHelper) *RuleBuilder {
	return &RuleBuilder{
		alphaBuilder:       NewAlphaRuleBuilder(utils),
		joinBuilder:        NewJoinRuleBuilder(utils),
		existsBuilder:      NewExistsRuleBuilder(utils),
		accumulatorBuilder: NewAccumulatorRuleBuilder(utils),
		utils:              utils,
		pipeline:           pipeline,
	}
}

// CreateRuleNodes creates rule nodes (Alpha, Beta, Terminal) from expressions
func (rb *RuleBuilder) CreateRuleNodes(network *ReteNetwork, expressions []interface{}) error {
	for i, exprInterface := range expressions {
		exprMap, ok := exprInterface.(map[string]interface{})
		if !ok {
			return fmt.Errorf("format expression invalide: %T", exprInterface)
		}

		// Extract the ruleId from the expression
		ruleID := fmt.Sprintf("rule_%d", i) // Default fallback
		if ruleIdValue, ok := exprMap["ruleId"]; ok {
			if ruleIdStr, ok := ruleIdValue.(string); ok && ruleIdStr != "" {
				ruleID = ruleIdStr
			}
		}

		// Create the rule
		err := rb.CreateSingleRule(network, ruleID, exprMap)
		if err != nil {
			return fmt.Errorf("erreur création règle %s: %w", ruleID, err)
		}

		fmt.Printf("   ✓ Règle créée: %s\n", ruleID)
	}

	return nil
}

// CreateSingleRule creates a single rule (refactored into small functions)
func (rb *RuleBuilder) CreateSingleRule(network *ReteNetwork, ruleID string, exprMap map[string]interface{}) error {
	// Validate pipeline is not nil
	if rb.pipeline == nil {
		return fmt.Errorf("pipeline is nil - cannot create rule")
	}

	// Step 1: Extract the action
	action, err := rb.pipeline.extractActionFromExpression(exprMap, ruleID)
	if err != nil {
		return err
	}

	// Step 2: Extract and analyze constraints
	constraintsData, hasConstraints := exprMap["constraints"]
	var condition map[string]interface{}
	var hasAggregation bool

	if hasConstraints {
		// Detect if this is an aggregation (from constraints)
		hasAggregation = rb.pipeline.detectAggregation(constraintsData)

		// Build the appropriate condition
		condition, err = rb.pipeline.buildConditionFromConstraints(constraintsData)
		if err != nil {
			return fmt.Errorf("erreur construction condition pour règle %s: %w", ruleID, err)
		}
	} else {
		condition = map[string]interface{}{
			"type": ConditionTypeSimple,
		}
	}

	// Step 3: Extract variables
	variables, variableNames, variableTypes := rb.pipeline.extractVariablesFromExpression(exprMap)

	// Also check if any variables are aggregation variables (new syntax)
	if !hasAggregation {
		hasAggregation = rb.pipeline.hasAggregationVariables(exprMap)
	}

	// Step 4: Determine the rule type and create it
	ruleType := rb.pipeline.determineRuleType(exprMap, len(variables), hasAggregation)
	rb.pipeline.logRuleCreation(ruleType, ruleID, variableNames)

	// Delegate to the appropriate specialized builder
	return rb.createRuleByType(network, ruleID, ruleType, exprMap, condition, action,
		variables, variableNames, variableTypes, constraintsData, hasConstraints)
}

// createRuleByType delegates rule creation to the appropriate builder based on type
func (rb *RuleBuilder) createRuleByType(
	network *ReteNetwork,
	ruleID string,
	ruleType string,
	exprMap map[string]interface{},
	condition map[string]interface{},
	action *Action,
	variables []map[string]interface{},
	variableNames []string,
	variableTypes []string,
	constraintsData interface{},
	hasConstraints bool,
) error {
	switch ruleType {
	case "exists":
		return rb.existsBuilder.CreateExistsRule(network, ruleID, exprMap, condition, action)

	case "accumulator":
		return rb.createAccumulatorRuleWithInfo(network, ruleID, exprMap, condition, action,
			variables, variableNames, variableTypes, constraintsData, hasConstraints)

	case "join":
		return rb.joinBuilder.CreateJoinRule(network, ruleID, variableNames, variableTypes, condition, action)

	case "alpha":
		return rb.alphaBuilder.CreateAlphaRule(network, ruleID, variables, variableNames, variableTypes, condition, action)

	default:
		return fmt.Errorf("type de règle inconnu: %s", ruleType)
	}
}

// createAccumulatorRuleWithInfo handles accumulator rule creation with aggregation info extraction
func (rb *RuleBuilder) createAccumulatorRuleWithInfo(
	network *ReteNetwork,
	ruleID string,
	exprMap map[string]interface{},
	condition map[string]interface{},
	action *Action,
	variables []map[string]interface{},
	variableNames []string,
	variableTypes []string,
	constraintsData interface{},
	hasConstraints bool,
) error {
	var aggInfo *AggregationInfo
	var err error

	// Check if this is the new multi-pattern aggregation syntax
	if _, hasPatterns := exprMap["patterns"]; hasPatterns {
		// Check if this is multi-source aggregation
		if rb.accumulatorBuilder.IsMultiSourceAggregation(exprMap) {
			fmt.Printf("   📊 Multi-source aggregation détectée pour: %s\n", ruleID)
			aggInfo, err = rb.pipeline.extractMultiSourceAggregationInfo(exprMap)
		} else {
			aggInfo, err = rb.pipeline.extractAggregationInfoFromVariables(exprMap)
		}
	} else {
		// Old AccumulateConstraint syntax
		aggInfo, err = rb.pipeline.extractAggregationInfo(constraintsData)
	}

	if err != nil {
		fmt.Printf("   ⚠️  Impossible d'extraire info agrégation: %v, utilisation JoinNode standard\n", err)
		return rb.joinBuilder.CreateJoinRule(network, ruleID, variableNames, variableTypes, condition, action)
	}

	// Check if we need multi-source accumulator
	if len(aggInfo.SourcePatterns) > 1 || len(aggInfo.AggregationVars) > 1 {
		return rb.accumulatorBuilder.CreateMultiSourceAccumulatorRule(network, ruleID, aggInfo, action)
	}

	return rb.accumulatorBuilder.CreateAccumulatorRule(network, ruleID, variables, variableNames, variableTypes, aggInfo, action)
}
