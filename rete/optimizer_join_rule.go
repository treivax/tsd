// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// JoinRuleRemovalStrategy handles removal of rules containing join nodes
type JoinRuleRemovalStrategy struct {
	network *ReteNetwork
	helpers *OptimizerHelpers
}

// NewJoinRuleRemovalStrategy creates a new join rule removal strategy
func NewJoinRuleRemovalStrategy(network *ReteNetwork) *JoinRuleRemovalStrategy {
	return &JoinRuleRemovalStrategy{
		network: network,
		helpers: NewOptimizerHelpers(network),
	}
}

// Name returns the name of this strategy
func (s *JoinRuleRemovalStrategy) Name() string {
	return "JoinRuleRemoval"
}

// CanHandle returns true if this strategy can handle the given rule
// Join rules contain one or more join nodes
func (s *JoinRuleRemovalStrategy) CanHandle(ruleID string, nodeIDs []string) bool {
	for _, nodeID := range nodeIDs {
		if s.helpers.IsJoinNode(nodeID) {
			return true
		}
	}
	return false
}

// RemoveRule removes a rule that contains join nodes
func (s *JoinRuleRemovalStrategy) RemoveRule(ruleID string, nodeIDs []string) (int, error) {
	s.network.logger.Debug("   🔗 Utilisation de la stratégie JoinRuleRemoval pour %s", ruleID)

	// Classify nodes by type
	classification := s.helpers.ClassifyNodes(nodeIDs)

	deletedCount := 0

	// Step 1: Remove terminal nodes first
	for _, nodeID := range classification.TerminalNodes {
		if err := s.helpers.RemoveNodeWithCheck(nodeID, ruleID); err == nil {
			deletedCount++
			s.network.logger.Debug("   🗑️  TerminalNode %s removed", nodeID)
		}
	}

	// Step 2: Remove join nodes with reference counting
	for _, nodeID := range classification.JoinNodes {
		// Remove rule reference from join node (BetaSharingRegistry always initialized)
		canDelete, err := s.network.BetaSharingRegistry.RemoveRuleFromJoinNode(nodeID, ruleID)
		if err != nil {
			s.network.logger.Warn("   ⚠️  Error removing rule from join node %s: %v", nodeID, err)
			continue
		}

		if canDelete {
			// No more rules reference this join node - safe to delete
			if err := s.helpers.RemoveJoinNodeFromNetwork(nodeID); err == nil {
				deletedCount++
				s.network.logger.Debug("   🗑️  JoinNode %s removed (no more references)", nodeID)
			}
		} else {
			// Join node is still shared by other rules
			refCount := s.network.BetaSharingRegistry.GetJoinNodeRefCount(nodeID)
			s.network.logger.Debug("   ✓ JoinNode %s preserved (%d rule(s) remaining)", nodeID, refCount)
		}
	}

	// Step 3: Remove alpha nodes
	for _, nodeID := range classification.AlphaNodes {
		if err := s.helpers.RemoveNodeWithCheck(nodeID, ruleID); err == nil {
			deletedCount++
			s.network.logger.Debug("   🗑️  AlphaNode %s removed", nodeID)
		} else {
			lifecycle, _ := s.network.LifecycleManager.GetNodeLifecycle(nodeID)
			if lifecycle != nil && lifecycle.HasReferences() {
				s.network.logger.Debug("   ✓ AlphaNode %s preserved (%d reference(s))", nodeID, lifecycle.GetRefCount())
			}
		}
	}

	// Step 4: Type nodes are usually shared - only remove if no references
	for _, nodeID := range classification.TypeNodes {
		lifecycle, exists := s.network.LifecycleManager.GetNodeLifecycle(nodeID)
		if !exists {
			continue
		}

		shouldDelete, err := s.network.LifecycleManager.RemoveRuleFromNode(nodeID, ruleID)
		if err != nil {
			s.network.logger.Warn("   ⚠️  Error removing rule from type node %s: %v", nodeID, err)
			continue
		}

		if shouldDelete {
			if err := s.helpers.RemoveNodeFromNetwork(nodeID); err == nil {
				deletedCount++
				s.network.logger.Debug("   🗑️  TypeNode %s removed", nodeID)
			}
		} else {
			s.network.logger.Debug("   ✓ TypeNode %s preserved (%d reference(s))", nodeID, lifecycle.GetRefCount())
		}
	}

	s.network.logger.Info("✅ Rule %s removed successfully (%d node(s) deleted)", ruleID, deletedCount)
	return deletedCount, nil
}
