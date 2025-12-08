// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// AlphaChainRemovalStrategy handles removal of rules with alpha node chains
type AlphaChainRemovalStrategy struct {
	network *ReteNetwork
	helpers *OptimizerHelpers
}

// NewAlphaChainRemovalStrategy creates a new alpha chain removal strategy
func NewAlphaChainRemovalStrategy(network *ReteNetwork) *AlphaChainRemovalStrategy {
	return &AlphaChainRemovalStrategy{
		network: network,
		helpers: NewOptimizerHelpers(network),
	}
}

// Name returns the name of this strategy
func (s *AlphaChainRemovalStrategy) Name() string {
	return "AlphaChainRemoval"
}

// CanHandle returns true if this strategy can handle the given rule
// Alpha chain rules have chains of alpha nodes but no join nodes
func (s *AlphaChainRemovalStrategy) CanHandle(ruleID string, nodeIDs []string) bool {
	hasChain := false
	hasJoinNodes := false

	for _, nodeID := range nodeIDs {
		if s.helpers.IsPartOfChain(nodeID) {
			hasChain = true
		}
		if s.helpers.IsJoinNode(nodeID) {
			hasJoinNodes = true
		}
	}

	// Can handle if there's a chain but no join nodes
	return hasChain && !hasJoinNodes
}

// RemoveRule removes a rule with an alpha node chain
// Walks up the chain in reverse order from terminal to type node to remove nodes
func (s *AlphaChainRemovalStrategy) RemoveRule(ruleID string, nodeIDs []string) (int, error) {
	s.network.logger.Debug("   🔗 Utilisation de la stratégie AlphaChainRemoval pour %s", ruleID)

	// Classify nodes by type
	classification := s.helpers.ClassifyNodes(nodeIDs)

	deletedCount := 0

	// Step 1: Remove terminal node first
	if len(classification.TerminalNodes) > 0 {
		for _, terminalID := range classification.TerminalNodes {
			if err := s.helpers.RemoveNodeWithCheck(terminalID, ruleID); err == nil {
				deletedCount++
				s.network.logger.Debug("   🗑️  TerminalNode %s supprimé", terminalID)
			}
		}
	}

	// Step 2: Order alpha nodes in reverse chain order (terminal to type node)
	orderedAlphaNodes := s.helpers.OrderAlphaNodesReverse(classification.AlphaNodes)

	// Step 3: Walk through alpha nodes in reverse order
	stopDeletion := false
	for i, nodeID := range orderedAlphaNodes {
		lifecycle, exists := s.network.LifecycleManager.GetNodeLifecycle(nodeID)
		if !exists {
			continue
		}

		// Decrement RefCount for all nodes
		shouldDelete, err := s.network.LifecycleManager.RemoveRuleFromNode(nodeID, ruleID)
		if err != nil {
			s.network.logger.Warn("   ⚠️  Erreur lors de la suppression de la règle du nœud %s: %v", nodeID, err)
			continue
		}

		if !stopDeletion && shouldDelete {
			// RefCount == 0, we can delete
			if err := s.helpers.RemoveNodeFromNetwork(nodeID); err != nil {
				s.network.logger.Warn("   ⚠️  Erreur suppression nœud %s: %v", nodeID, err)
			} else {
				deletedCount++
				s.network.logger.Debug("   🗑️  AlphaNode %s supprimé (position %d dans la chaîne)", nodeID, len(orderedAlphaNodes)-i)
			}
		} else if !shouldDelete && !stopDeletion {
			// First shared node encountered - stop deletion but continue decrementing
			refCount := lifecycle.GetRefCount()
			s.network.logger.Debug("   ♻️  AlphaNode %s conservé (%d référence(s) restante(s)) - arrêt des suppressions", nodeID, refCount)
			s.network.logger.Debug("   ℹ️  Décrémentation du RefCount des nœuds parents partagés")
			stopDeletion = true
		} else if stopDeletion {
			// Parent nodes - just decrement RefCount
			refCount := lifecycle.GetRefCount()
			s.network.logger.Debug("   ♻️  AlphaNode %s: RefCount décrémenté (%d référence(s) restante(s))", nodeID, refCount)
		}
	}

	// Step 4: Remove other nodes (TypeNodes, etc.)
	for _, nodeID := range classification.OtherNodes {
		if err := s.helpers.RemoveNodeWithCheck(nodeID, ruleID); err == nil {
			deletedCount++
			lifecycle, _ := s.network.LifecycleManager.GetNodeLifecycle(nodeID)
			if lifecycle != nil {
				s.network.logger.Debug("   🗑️  %s %s supprimé", lifecycle.NodeType, nodeID)
			}
		}
	}

	// Step 5: Handle type nodes
	for _, nodeID := range classification.TypeNodes {
		if err := s.helpers.RemoveNodeWithCheck(nodeID, ruleID); err == nil {
			deletedCount++
			s.network.logger.Debug("   🗑️  TypeNode %s supprimé", nodeID)
		}
	}

	s.network.logger.Info("✅ Règle %s avec chaîne supprimée avec succès (%d nœud(s) supprimé(s))", ruleID, deletedCount)
	return deletedCount, nil
}
