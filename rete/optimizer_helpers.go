// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import "fmt"

// OptimizerHelpers provides common utility functions for network optimization
type OptimizerHelpers struct {
	network *ReteNetwork
}

// NewOptimizerHelpers creates a new helper instance
func NewOptimizerHelpers(network *ReteNetwork) *OptimizerHelpers {
	return &OptimizerHelpers{network: network}
}

// RemoveNodeWithCheck removes a node only if RefCount == 0
func (h *OptimizerHelpers) RemoveNodeWithCheck(nodeID, ruleID string) error {
	shouldDelete, err := h.network.LifecycleManager.RemoveRuleFromNode(nodeID, ruleID)
	if err != nil {
		return err
	}

	if shouldDelete {
		return h.RemoveNodeFromNetwork(nodeID)
	}

	return fmt.Errorf("nœud %s encore référencé", nodeID)
}

// RemoveNodeFromNetwork removes a node from the RETE network
// Only removes if RefCount == 0
func (h *OptimizerHelpers) RemoveNodeFromNetwork(nodeID string) error {
	// Verify node exists and can be removed
	lifecycle, exists := h.network.LifecycleManager.GetNodeLifecycle(nodeID)
	if !exists {
		return fmt.Errorf("nœud %s non trouvé dans le LifecycleManager", nodeID)
	}

	// Don't remove if node still has references
	if lifecycle.HasReferences() {
		return fmt.Errorf("impossible de supprimer le nœud %s: encore %d référence(s)",
			nodeID, lifecycle.GetRefCount())
	}

	// Determine node type and remove from appropriate map
	switch lifecycle.NodeType {
	case "type":
		return h.removeTypeNode(nodeID)
	case "alpha":
		return h.removeAlphaNode(nodeID)
	case "terminal":
		return h.removeTerminalNode(nodeID)
	case "join", "exists", "accumulate":
		return h.removeBetaNode(nodeID)
	}

	return fmt.Errorf("nœud %s non trouvé dans le réseau", nodeID)
}

// removeTypeNode removes a type node from the network
func (h *OptimizerHelpers) removeTypeNode(nodeID string) error {
	// Find and remove TypeNode
	for typeName, typeNode := range h.network.TypeNodes {
		if typeNode.GetID() == nodeID {
			// Disconnect from RootNode
			h.RemoveChildFromNode(h.network.RootNode, typeNode)
			delete(h.network.TypeNodes, typeName)
			return h.network.LifecycleManager.RemoveNode(nodeID)
		}
	}
	return fmt.Errorf("type node %s not found", nodeID)
}

// removeAlphaNode removes an alpha node from the network
func (h *OptimizerHelpers) removeAlphaNode(nodeID string) error {
	alphaNode, exists := h.network.AlphaNodes[nodeID]
	if !exists {
		return fmt.Errorf("alpha node %s not found", nodeID)
	}

	// Disconnect from parents (TypeNodes or other AlphaNodes)
	parent := h.GetChainParent(alphaNode)
	if parent != nil {
		h.RemoveChildFromNode(parent, alphaNode)
		h.network.logger.Debug("   🔗 AlphaNode %s déconnecté de son parent %s", nodeID, parent.GetID())
	}

	delete(h.network.AlphaNodes, nodeID)

	// Remove from AlphaSharingManager if it's a shared node
	if h.network.AlphaSharingManager != nil {
		if len(nodeID) > 6 && nodeID[:6] == "alpha_" {
			if err := h.network.AlphaSharingManager.RemoveAlphaNode(nodeID); err != nil {
				h.network.logger.Warn("   ⚠️  Erreur suppression AlphaNode du registre de partage: %v", err)
			} else {
				h.network.logger.Debug("   ✓ AlphaNode %s supprimé du AlphaSharingManager", nodeID)
			}
		}
	}

	return h.network.LifecycleManager.RemoveNode(nodeID)
}

// removeTerminalNode removes a terminal node from the network
func (h *OptimizerHelpers) removeTerminalNode(nodeID string) error {
	terminalNode, exists := h.network.TerminalNodes[nodeID]
	if !exists {
		return fmt.Errorf("terminal node %s not found", nodeID)
	}

	// Disconnect from parents (AlphaNodes or JoinNodes)
	for _, alphaNode := range h.network.AlphaNodes {
		h.RemoveChildFromNode(alphaNode, terminalNode)
	}

	// Also disconnect from BetaNodes if necessary
	for _, betaNode := range h.network.BetaNodes {
		if node, ok := betaNode.(Node); ok {
			h.RemoveChildFromNode(node, terminalNode)
		}
	}

	delete(h.network.TerminalNodes, nodeID)
	return h.network.LifecycleManager.RemoveNode(nodeID)
}

// removeBetaNode removes a beta node from the network
func (h *OptimizerHelpers) removeBetaNode(nodeID string) error {
	betaNode, exists := h.network.BetaNodes[nodeID]
	if !exists {
		return fmt.Errorf("beta node %s not found", nodeID)
	}

	// Disconnect from parents
	for _, typeNode := range h.network.TypeNodes {
		if node, ok := betaNode.(Node); ok {
			h.RemoveChildFromNode(typeNode, node)
		}
	}

	delete(h.network.BetaNodes, nodeID)
	return h.network.LifecycleManager.RemoveNode(nodeID)
}

// RemoveJoinNodeFromNetwork removes a join node and all its dependent nodes
func (h *OptimizerHelpers) RemoveJoinNodeFromNetwork(nodeID string) error {
	// Get the join node
	joinNode, exists := h.network.BetaNodes[nodeID]
	if !exists {
		return fmt.Errorf("join node %s not found in network", nodeID)
	}

	h.network.logger.Debug("   🗑️  Removing join node %s from network", nodeID)

	// Convert to JoinNode type
	jn, ok := joinNode.(*JoinNode)
	if !ok {
		return fmt.Errorf("beta node %s is not a JoinNode", nodeID)
	}

	// Step 1: Remove dependent terminal nodes
	for terminalID := range h.network.TerminalNodes {
		isChild := false
		for _, child := range jn.GetChildren() {
			if child.GetID() == terminalID {
				isChild = true
				break
			}
		}

		if isChild {
			delete(h.network.TerminalNodes, terminalID)
			h.network.logger.Debug("   🗑️  Removed terminal node %s (child of join node)", terminalID)

			if h.network.LifecycleManager != nil {
				h.network.LifecycleManager.RemoveNode(terminalID)
			}
		}
	}

	// Step 2: Disconnect from parent nodes
	var node Node = jn

	// Disconnect from alpha nodes
	for _, alphaNode := range h.network.AlphaNodes {
		h.DisconnectChild(alphaNode, node)
	}

	// Disconnect from other beta nodes (cascading joins)
	for betaNodeID, betaNode := range h.network.BetaNodes {
		if betaNodeID != nodeID {
			if bn, ok := betaNode.(*JoinNode); ok {
				h.DisconnectChild(bn, node)
			}
		}
	}

	// Disconnect from type nodes
	for _, typeNode := range h.network.TypeNodes {
		h.DisconnectChild(typeNode, node)
	}

	// Step 3: Remove from beta nodes map
	delete(h.network.BetaNodes, nodeID)

	// Step 4: Remove from lifecycle manager
	if h.network.LifecycleManager != nil {
		if err := h.network.LifecycleManager.RemoveNode(nodeID); err != nil {
			h.network.logger.Warn("   ⚠️  Warning: failed to remove join node %s from lifecycle manager: %v", nodeID, err)
		}
	}

	// Step 5: Remove from beta sharing registry (always initialized)
	if err := h.network.BetaSharingRegistry.UnregisterJoinNode(nodeID); err != nil {
		h.network.logger.Warn("   ⚠️  Warning: failed to unregister join node %s from beta sharing: %v", nodeID, err)
	}

	h.network.logger.Info("   ✅ Join node %s successfully removed from network", nodeID)
	return nil
}

// RemoveChildFromNode removes a child node from a parent node
func (h *OptimizerHelpers) RemoveChildFromNode(parent Node, child Node) {
	if parent == nil || child == nil {
		return
	}

	children := parent.GetChildren()
	newChildren := make([]Node, 0, len(children))
	for _, c := range children {
		if c.GetID() != child.GetID() {
			newChildren = append(newChildren, c)
		}
	}

	// Update children (requires cast to concrete type)
	switch p := parent.(type) {
	case *RootNode:
		p.Children = newChildren
	case *TypeNode:
		p.Children = newChildren
	case *AlphaNode:
		p.Children = newChildren
	case *JoinNode:
		p.Children = newChildren
	case *ExistsNode:
		p.Children = newChildren
	}
}

// DisconnectChild removes a child from a node's children list
func (h *OptimizerHelpers) DisconnectChild(parent Node, child Node) {
	if parent == nil || child == nil {
		return
	}

	children := parent.GetChildren()
	newChildren := make([]Node, 0, len(children))
	for _, c := range children {
		if c.GetID() != child.GetID() {
			newChildren = append(newChildren, c)
		}
	}

	// Update parent's children list using SetChildren interface
	if baseNode, ok := parent.(interface{ SetChildren([]Node) }); ok {
		baseNode.SetChildren(newChildren)
	}
}

// IsPartOfChain detects if a node is part of an alpha node chain
func (h *OptimizerHelpers) IsPartOfChain(nodeID string) bool {
	lifecycle, exists := h.network.LifecycleManager.GetNodeLifecycle(nodeID)
	if !exists || lifecycle.NodeType != "alpha" {
		return false
	}

	alphaNode, exists := h.network.AlphaNodes[nodeID]
	if !exists {
		return false
	}

	// An AlphaNode is part of a chain if:
	// 1. Its parent is another AlphaNode, OR
	// 2. One of its children is another AlphaNode

	parent := h.GetChainParent(alphaNode)
	if parent != nil && parent.GetType() == "alpha" {
		return true
	}

	children := alphaNode.GetChildren()
	for _, child := range children {
		if child.GetType() == "alpha" {
			return true
		}
	}

	return false
}

// GetChainParent retrieves the parent node of an AlphaNode in a chain
func (h *OptimizerHelpers) GetChainParent(alphaNode *AlphaNode) Node {
	if alphaNode == nil {
		return nil
	}

	alphaID := alphaNode.GetID()

	// Search in TypeNodes
	for _, typeNode := range h.network.TypeNodes {
		for _, child := range typeNode.GetChildren() {
			if child.GetID() == alphaID {
				return typeNode
			}
		}
	}

	// Search in other AlphaNodes
	for _, parentAlpha := range h.network.AlphaNodes {
		if parentAlpha.GetID() == alphaID {
			continue
		}
		for _, child := range parentAlpha.GetChildren() {
			if child.GetID() == alphaID {
				return parentAlpha
			}
		}
	}

	return nil
}

// OrderAlphaNodesReverse orders alpha nodes in reverse chain order
// (from the node furthest from TypeNode towards the TypeNode)
func (h *OptimizerHelpers) OrderAlphaNodesReverse(alphaNodeIDs []string) []string {
	if len(alphaNodeIDs) <= 1 {
		return alphaNodeIDs
	}

	// Build parent-child graph to find ordering
	childToParent := make(map[string]string)
	hasParent := make(map[string]bool)

	for _, nodeID := range alphaNodeIDs {
		alphaNode, exists := h.network.AlphaNodes[nodeID]
		if !exists {
			continue
		}

		parent := h.GetChainParent(alphaNode)
		if parent != nil {
			parentID := parent.GetID()
			// Check if parent is also an AlphaNode in our list
			for _, candidateID := range alphaNodeIDs {
				if candidateID == parentID {
					childToParent[nodeID] = parentID
					hasParent[nodeID] = true
					break
				}
			}
		}
	}

	// Find terminal node of the chain (one that is not a parent of anyone)
	var terminalNode string
	for _, nodeID := range alphaNodeIDs {
		isParent := false
		for _, parentID := range childToParent {
			if parentID == nodeID {
				isParent = true
				break
			}
		}
		if !isParent {
			terminalNode = nodeID
			break
		}
	}

	// If no chain structure detected, return original order
	if terminalNode == "" {
		return alphaNodeIDs
	}

	// Walk up the chain from terminal
	ordered := make([]string, 0, len(alphaNodeIDs))
	current := terminalNode
	visited := make(map[string]bool)

	for current != "" && !visited[current] {
		ordered = append(ordered, current)
		visited[current] = true
		current = childToParent[current]
	}

	// Add unvisited nodes (in case of disconnected nodes)
	for _, nodeID := range alphaNodeIDs {
		if !visited[nodeID] {
			ordered = append(ordered, nodeID)
		}
	}

	return ordered
}

// IsJoinNode checks if a node ID corresponds to a JoinNode
func (h *OptimizerHelpers) IsJoinNode(nodeID string) bool {
	_, exists := h.network.BetaNodes[nodeID]
	return exists
}

// ClassifyNodes separates nodes by type
func (h *OptimizerHelpers) ClassifyNodes(nodeIDs []string) *NodeClassification {
	classification := &NodeClassification{
		TerminalNodes: make([]string, 0),
		JoinNodes:     make([]string, 0),
		AlphaNodes:    make([]string, 0),
		TypeNodes:     make([]string, 0),
		OtherNodes:    make([]string, 0),
	}

	for _, nodeID := range nodeIDs {
		lifecycle, exists := h.network.LifecycleManager.GetNodeLifecycle(nodeID)
		if !exists {
			continue
		}

		switch lifecycle.NodeType {
		case "terminal":
			classification.TerminalNodes = append(classification.TerminalNodes, nodeID)
		case "join":
			classification.JoinNodes = append(classification.JoinNodes, nodeID)
		case "alpha":
			classification.AlphaNodes = append(classification.AlphaNodes, nodeID)
		case "type":
			classification.TypeNodes = append(classification.TypeNodes, nodeID)
		default:
			classification.OtherNodes = append(classification.OtherNodes, nodeID)
		}
	}

	return classification
}
