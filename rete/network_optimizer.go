// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import "fmt"

// RemoveRule removes a rule and all its nodes that are no longer used
// This is the main entry point for rule removal, which delegates to appropriate strategies
func (rn *ReteNetwork) RemoveRule(ruleID string) error {
	rn.logger.Info("🗑️ Suppression de la règle: %s", ruleID)

	if rn.LifecycleManager == nil {
		return fmt.Errorf("LifecycleManager non initialisé")
	}

	// Get all nodes used by this rule
	nodeIDs := rn.LifecycleManager.GetNodesForRule(ruleID)
	if len(nodeIDs) == 0 {
		return fmt.Errorf("règle %s non trouvée ou aucun nœud associé", ruleID)
	}

	rn.logger.Debug("   📊 Nœuds associés à la règle %s: %d", ruleID, len(nodeIDs))

	// Create strategies
	simpleStrategy := NewSimpleRuleRemovalStrategy(rn)
	alphaChainStrategy := NewAlphaChainRemovalStrategy(rn)
	joinStrategy := NewJoinRuleRemovalStrategy(rn)

	// Create strategy selector
	selector := NewDefaultStrategySelector(rn, simpleStrategy, alphaChainStrategy, joinStrategy)

	// Select and execute the appropriate strategy
	strategy := selector.SelectStrategy(ruleID, nodeIDs)
	rn.logger.Debug("   🎯 Stratégie sélectionnée: %s", strategy.Name())

	deletedCount, err := strategy.RemoveRule(ruleID, nodeIDs)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression de la règle %s: %w", ruleID, err)
	}

	rn.logger.Info("✅ Règle %s supprimée avec succès (%d nœud(s) supprimé(s))", ruleID, deletedCount)
	return nil
}

// removeNodeWithCheck removes a node only if RefCount == 0
// This is a convenience method that delegates to helpers
func (rn *ReteNetwork) removeNodeWithCheck(nodeID, ruleID string) error {
	helpers := NewOptimizerHelpers(rn)
	return helpers.RemoveNodeWithCheck(nodeID, ruleID)
}

// removeNodeFromNetwork removes a node from the RETE network
// This is a convenience method that delegates to helpers
func (rn *ReteNetwork) removeNodeFromNetwork(nodeID string) error {
	helpers := NewOptimizerHelpers(rn)
	return helpers.RemoveNodeFromNetwork(nodeID)
}

// removeJoinNodeFromNetwork removes a join node and all its dependent nodes from the network
// This is a convenience method that delegates to helpers
func (rn *ReteNetwork) removeJoinNodeFromNetwork(nodeID string) error {
	helpers := NewOptimizerHelpers(rn)
	return helpers.RemoveJoinNodeFromNetwork(nodeID)
}

// removeChildFromNode removes a child node from a parent node
// This is a convenience method that delegates to helpers
func (rn *ReteNetwork) removeChildFromNode(parent Node, child Node) {
	helpers := NewOptimizerHelpers(rn)
	helpers.RemoveChildFromNode(parent, child)
}

// disconnectChild removes a child from a node's children list
// This is a convenience method that delegates to helpers
func (rn *ReteNetwork) disconnectChild(parent Node, child Node) {
	helpers := NewOptimizerHelpers(rn)
	helpers.DisconnectChild(parent, child)
}

// orderAlphaNodesReverse orders alpha nodes in reverse chain order
// This is a convenience method that delegates to helpers
func (rn *ReteNetwork) orderAlphaNodesReverse(alphaNodeIDs []string) []string {
	helpers := NewOptimizerHelpers(rn)
	return helpers.OrderAlphaNodesReverse(alphaNodeIDs)
}

// isPartOfChain detects if a node is part of an alpha node chain
// This is a convenience method that delegates to helpers
func (rn *ReteNetwork) isPartOfChain(nodeID string) bool {
	helpers := NewOptimizerHelpers(rn)
	return helpers.IsPartOfChain(nodeID)
}

// getChainParent retrieves the parent node of an AlphaNode in a chain
// This is a convenience method that delegates to helpers
func (rn *ReteNetwork) getChainParent(alphaNode *AlphaNode) Node {
	helpers := NewOptimizerHelpers(rn)
	return helpers.GetChainParent(alphaNode)
}

// isJoinNode checks if a node ID corresponds to a JoinNode
// This is a convenience method that delegates to helpers
func (rn *ReteNetwork) isJoinNode(nodeID string) bool {
	helpers := NewOptimizerHelpers(rn)
	return helpers.IsJoinNode(nodeID)
}
