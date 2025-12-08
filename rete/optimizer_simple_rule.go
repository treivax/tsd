// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// SimpleRuleRemovalStrategy handles removal of simple rules without chains
type SimpleRuleRemovalStrategy struct {
	network *ReteNetwork
	helpers *OptimizerHelpers
}

// NewSimpleRuleRemovalStrategy creates a new simple rule removal strategy
func NewSimpleRuleRemovalStrategy(network *ReteNetwork) *SimpleRuleRemovalStrategy {
	return &SimpleRuleRemovalStrategy{
		network: network,
		helpers: NewOptimizerHelpers(network),
	}
}

// Name returns the name of this strategy
func (s *SimpleRuleRemovalStrategy) Name() string {
	return "SimpleRuleRemoval"
}

// CanHandle returns true if this strategy can handle the given rule
// Simple rules are those without chains or join nodes
func (s *SimpleRuleRemovalStrategy) CanHandle(ruleID string, nodeIDs []string) bool {
	// If there are join nodes or chains, this strategy cannot handle it
	for _, nodeID := range nodeIDs {
		if s.helpers.IsPartOfChain(nodeID) {
			return false
		}
		if s.helpers.IsJoinNode(nodeID) {
			return false
		}
	}
	return true
}

// RemoveRule removes a simple rule (without chains or joins)
func (s *SimpleRuleRemovalStrategy) RemoveRule(ruleID string, nodeIDs []string) (int, error) {
	s.network.logger.Debug("   📝 Utilisation de la stratégie SimpleRuleRemoval pour %s", ruleID)

	// Iterate through each node and remove the rule reference
	nodesToDelete := make([]string, 0)
	for _, nodeID := range nodeIDs {
		shouldDelete, err := s.network.LifecycleManager.RemoveRuleFromNode(nodeID, ruleID)
		if err != nil {
			s.network.logger.Warn("   ⚠️  Erreur lors de la suppression de la règle du nœud %s: %v", nodeID, err)
			continue
		}

		if shouldDelete {
			nodesToDelete = append(nodesToDelete, nodeID)
			s.network.logger.Debug("   ✓ Nœud %s marqué pour suppression (plus de références)", nodeID)
		} else {
			lifecycle, _ := s.network.LifecycleManager.GetNodeLifecycle(nodeID)
			s.network.logger.Debug("   ✓ Nœud %s conservé (%d référence(s) restante(s))", nodeID, lifecycle.GetRefCount())
		}
	}

	// Remove nodes that have no more references
	deletedCount := 0
	for _, nodeID := range nodesToDelete {
		if err := s.helpers.RemoveNodeFromNetwork(nodeID); err != nil {
			s.network.logger.Warn("   ⚠️  Erreur lors de la suppression du nœud %s: %v", nodeID, err)
		} else {
			s.network.logger.Debug("   🗑️  Nœud %s supprimé du réseau", nodeID)
			deletedCount++
		}
	}

	if deletedCount > 0 {
		s.network.logger.Info("✅ Règle %s supprimée avec succès (%d nœud(s) supprimé(s))", ruleID, deletedCount)
	} else {
		s.network.logger.Info("✅ Règle %s supprimée (aucun nœud à supprimer, tous partagés)", ruleID)
	}

	return deletedCount, nil
}
