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
