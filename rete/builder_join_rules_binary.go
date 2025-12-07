// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// ============================================================================
// Binary Join Rule Creation (2 variables)
// ============================================================================

// createBinaryJoinRule creates a simple binary join rule (2 variables)
func (jrb *JoinRuleBuilder) createBinaryJoinRule(
	network *ReteNetwork,
	ruleID string,
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	terminalNode *TerminalNode,
) error {
	// Delegate to orchestrated version
	return jrb.createBinaryJoinRuleOrchestrated(
		network, ruleID, variableNames, variableTypes, condition, terminalNode,
	)
}
