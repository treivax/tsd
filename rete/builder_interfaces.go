// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// PipelineHelper defines the interface for pipeline operations needed by builders.
// This interface abstracts the constraint pipeline functionality that builders
// depend on, allowing for better separation of concerns and testability.
type PipelineHelper interface {
	// extractActionFromExpression extracts the action from a rule expression
	extractActionFromExpression(exprMap map[string]interface{}, ruleID string) (*Action, error)

	// detectAggregation detects if constraints contain aggregation operations
	detectAggregation(constraintsData interface{}) bool

	// extractAggregationInfo extracts aggregation information from constraints
	extractAggregationInfo(constraintsData interface{}) (*AggregationInfo, error)

	// extractAggregationInfoFromVariables extracts aggregation info from variable definitions
	extractAggregationInfoFromVariables(exprMap map[string]interface{}) (*AggregationInfo, error)

	// extractMultiSourceAggregationInfo extracts info for multi-source aggregations
	extractMultiSourceAggregationInfo(exprMap map[string]interface{}) (*AggregationInfo, error)

	// hasAggregationVariables checks if the expression has aggregation variables
	hasAggregationVariables(exprMap map[string]interface{}) bool

	// buildConditionFromConstraints builds a condition map from constraints
	buildConditionFromConstraints(constraintsData interface{}) (map[string]interface{}, error)

	// extractVariablesFromExpression extracts variables from rule expressions
	// Returns: (variables []map[string]interface{}, variableNames []string, variableTypes []string)
	extractVariablesFromExpression(exprMap map[string]interface{}) ([]map[string]interface{}, []string, []string)

	// determineRuleType determines the type of rule based on the expression
	// Returns one of: "alpha", "join", "exists", "accumulator"
	determineRuleType(exprMap map[string]interface{}, varCount int, hasAggregation bool) string

	// logRuleCreation logs information about rule creation
	logRuleCreation(ruleType string, ruleID string, variableNames []string)
}
