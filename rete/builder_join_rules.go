// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// JoinRuleBuilder handles the creation of join rules
type JoinRuleBuilder struct {
	utils                   *BuilderUtils
	enableDecomposition     bool // Enable arithmetic decomposition
	decompositionComplexity int  // Minimum complexity to trigger decomposition
}

// NewJoinRuleBuilder creates a new JoinRuleBuilder instance
func NewJoinRuleBuilder(utils *BuilderUtils) *JoinRuleBuilder {
	return &JoinRuleBuilder{
		utils:                   utils,
		enableDecomposition:     true, // Always enable decomposition
		decompositionComplexity: 1,    // Decompose all arithmetic expressions (even single operations)
	}
}

// SetDecompositionEnabled enables or disables arithmetic decomposition
func (jrb *JoinRuleBuilder) SetDecompositionEnabled(enabled bool) {
	jrb.enableDecomposition = enabled
}

// SetDecompositionComplexity sets the minimum complexity threshold for decomposition
func (jrb *JoinRuleBuilder) SetDecompositionComplexity(complexity int) {
	jrb.decompositionComplexity = complexity
}

// CreateJoinRule creates a join rule with JoinNode
func (jrb *JoinRuleBuilder) CreateJoinRule(
	network *ReteNetwork,
	ruleID string,
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	action *Action,
) error {
	// Create the terminal node for this rule
	terminalNode := jrb.utils.CreateTerminalNode(network, ruleID, action)

	// Delegate to the appropriate function based on the number of variables
	if len(variableNames) > 2 {
		return jrb.createCascadeJoinRule(network, ruleID, variableNames, variableTypes, condition, terminalNode)
	}

	return jrb.createBinaryJoinRule(network, ruleID, variableNames, variableTypes, condition, terminalNode)
}

// Note: Binary join creation has been extracted to builder_join_rules_binary.go

// Note: Cascade join creation and helper functions have been extracted to builder_join_rules_cascade.go
