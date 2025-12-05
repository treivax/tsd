// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import "fmt"

// RemoveRule supprime une règle et tous ses nœuds qui ne sont plus utilisés
func (rn *ReteNetwork) RemoveRule(ruleID string) error {
	rn.logger.Info("🗑️ Suppression de la règle: %s", ruleID)

	if rn.LifecycleManager == nil {
		return fmt.Errorf("LifecycleManager non initialisé")
	}

	// Récupérer tous les nœuds utilisés par cette règle
	nodeIDs := rn.LifecycleManager.GetNodesForRule(ruleID)
	if len(nodeIDs) == 0 {
		return fmt.Errorf("règle %s non trouvée ou aucun nœud associé", ruleID)
	}

	rn.logger.Debug("   📊 Nœuds associés à la règle %s: %d", ruleID, len(nodeIDs))

	// Detect rule type and use appropriate removal strategy
	hasChain := false
	hasJoinNodes := false

	for _, nodeID := range nodeIDs {
		if rn.isPartOfChain(nodeID) {
			hasChain = true
		}
		if rn.isJoinNode(nodeID) {
			hasJoinNodes = true
		}
	}

	// Utiliser la suppression optimisée pour les chaînes avec joins
	if hasJoinNodes {
		rn.logger.Debug("   🔗 JoinNodes détectés, utilisation de la suppression avec lifecycle")
		return rn.removeRuleWithJoins(ruleID, nodeIDs)
	}

	// Utiliser la suppression optimisée pour les chaînes alpha
	if hasChain {
		rn.logger.Debug("   🔗 Chaîne d'AlphaNodes détectée, utilisation de la suppression optimisée")
		return rn.removeAlphaChain(ruleID)
	}

	// Comportement classique pour les règles simples
	return rn.removeSimpleRule(ruleID, nodeIDs)
}

// removeSimpleRule supprime une règle simple (sans chaîne)
func (rn *ReteNetwork) removeSimpleRule(ruleID string, nodeIDs []string) error {
	// Parcourir chaque nœud et retirer la référence à la règle
	nodesToDelete := make([]string, 0)
	for _, nodeID := range nodeIDs {
		shouldDelete, err := rn.LifecycleManager.RemoveRuleFromNode(nodeID, ruleID)
		if err != nil {
			rn.logger.Warn("   ⚠️  Erreur lors de la suppression de la règle du nœud %s: %v", nodeID, err)
			continue
		}

		if shouldDelete {
			nodesToDelete = append(nodesToDelete, nodeID)
			rn.logger.Debug("   ✓ Nœud %s marqué pour suppression (plus de références)", nodeID)
		} else {
			lifecycle, _ := rn.LifecycleManager.GetNodeLifecycle(nodeID)
			rn.logger.Debug("   ✓ Nœud %s conservé (%d référence(s) restante(s))", nodeID, lifecycle.GetRefCount())
		}
	}

	// Supprimer les nœuds qui n'ont plus de références
	for _, nodeID := range nodesToDelete {
		if err := rn.removeNodeFromNetwork(nodeID); err != nil {
			rn.logger.Warn("   ⚠️  Erreur lors de la suppression du nœud %s: %v", nodeID, err)
		} else {
			rn.logger.Debug("   🗑️  Nœud %s supprimé du réseau", nodeID)
		}
	}

	rn.logger.Info("✅ Règle %s supprimée avec succès (%d nœud(s) supprimé(s))", ruleID, len(nodesToDelete))
	return nil
}

// removeAlphaChain supprime une règle avec une chaîne d'AlphaNodes
// Remonte la chaîne en ordre inverse depuis le terminal pour supprimer les nœuds
func (rn *ReteNetwork) removeAlphaChain(ruleID string) error {
	// Récupérer tous les nœuds de la règle
	nodeIDs := rn.LifecycleManager.GetNodesForRule(ruleID)

	// Séparer les nœuds par type
	var terminalID string
	alphaNodes := make([]string, 0)
	otherNodes := make([]string, 0)

	for _, nodeID := range nodeIDs {
		lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(nodeID)
		if !exists {
			continue
		}

		switch lifecycle.NodeType {
		case "terminal":
			terminalID = nodeID
		case "alpha":
			alphaNodes = append(alphaNodes, nodeID)
		default:
			otherNodes = append(otherNodes, nodeID)
		}
	}

	// Supprimer le terminal en premier
	deletedCount := 0
	if terminalID != "" {
		if err := rn.removeNodeWithCheck(terminalID, ruleID); err == nil {
			deletedCount++
			rn.logger.Debug("   🗑️  TerminalNode %s supprimé", terminalID)
		}
	}

	// Ordonner les AlphaNodes dans l'ordre inverse de la chaîne (du terminal vers le TypeNode)
	orderedAlphaNodes := rn.orderAlphaNodesReverse(alphaNodes)

	// Parcourir les AlphaNodes en ordre inverse
	stopDeletion := false
	for i, nodeID := range orderedAlphaNodes {
		lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(nodeID)
		if !exists {
			continue
		}

		// Décrémenter RefCount pour tous les nœuds
		shouldDelete, err := rn.LifecycleManager.RemoveRuleFromNode(nodeID, ruleID)
		if err != nil {
			rn.logger.Warn("   ⚠️  Erreur lors de la suppression de la règle du nœud %s: %v", nodeID, err)
			continue
		}

		if !stopDeletion && shouldDelete {
			// RefCount == 0, on peut supprimer
			if err := rn.removeNodeFromNetwork(nodeID); err != nil {
				rn.logger.Warn("   ⚠️  Erreur suppression nœud %s: %v", nodeID, err)
			} else {
				deletedCount++
				rn.logger.Debug("   🗑️  AlphaNode %s supprimé (position %d dans la chaîne)", nodeID, len(orderedAlphaNodes)-i)
			}
		} else if !shouldDelete && !stopDeletion {
			// Premier nœud partagé rencontré - on arrête la suppression mais on continue à décrémenter
			refCount := lifecycle.GetRefCount()
			rn.logger.Debug("   ♻️  AlphaNode %s conservé (%d référence(s) restante(s)) - arrêt des suppressions", nodeID, refCount)
			rn.logger.Debug("   ℹ️  Décrémentation du RefCount des nœuds parents partagés")
			stopDeletion = true
		} else if stopDeletion {
			// Nœuds parents - juste décrémenter le RefCount
			refCount := lifecycle.GetRefCount()
			rn.logger.Debug("   ♻️  AlphaNode %s: RefCount décrémenté (%d référence(s) restante(s))", nodeID, refCount)
		}
	}

	// Supprimer les autres nœuds (TypeNodes, JoinNodes, etc.)
	for _, nodeID := range otherNodes {
		if err := rn.removeNodeWithCheck(nodeID, ruleID); err == nil {
			deletedCount++
			lifecycle, _ := rn.LifecycleManager.GetNodeLifecycle(nodeID)
			rn.logger.Debug("   🗑️  %s %s supprimé", lifecycle.NodeType, nodeID)
		}
	}

	rn.logger.Info("✅ Règle %s avec chaîne supprimée avec succès (%d nœud(s) supprimé(s))", ruleID, deletedCount)
	return nil
}

// removeRuleWithJoins removes a rule that contains join nodes
func (rn *ReteNetwork) removeRuleWithJoins(ruleID string, nodeIDs []string) error {
	rn.logger.Debug("   🔗 Removing rule with join nodes: %s", ruleID)

	// Separate nodes by type
	var terminalNodes []string
	var joinNodes []string
	var alphaNodes []string
	var typeNodes []string

	for _, nodeID := range nodeIDs {
		lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(nodeID)
		if !exists {
			continue
		}

		switch lifecycle.NodeType {
		case "terminal":
			terminalNodes = append(terminalNodes, nodeID)
		case "join":
			joinNodes = append(joinNodes, nodeID)
		case "alpha":
			alphaNodes = append(alphaNodes, nodeID)
		case "type":
			typeNodes = append(typeNodes, nodeID)
		}
	}

	deletedCount := 0

	// Step 1: Remove terminal nodes first
	for _, nodeID := range terminalNodes {
		if err := rn.removeNodeWithCheck(nodeID, ruleID); err == nil {
			deletedCount++
			rn.logger.Debug("   🗑️  TerminalNode %s removed", nodeID)
		}
	}

	// Step 2: Remove join nodes with reference counting
	for _, nodeID := range joinNodes {
		// Remove rule reference from join node
		if rn.BetaSharingRegistry != nil {
			canDelete, err := rn.BetaSharingRegistry.RemoveRuleFromJoinNode(nodeID, ruleID)
			if err != nil {
				rn.logger.Warn("   ⚠️  Error removing rule from join node %s: %v", nodeID, err)
				continue
			}

			if canDelete {
				// No more rules reference this join node - safe to delete
				if err := rn.removeJoinNodeFromNetwork(nodeID); err == nil {
					deletedCount++
					rn.logger.Debug("   🗑️  JoinNode %s removed (no more references)", nodeID)
				}
			} else {
				// Join node is still shared by other rules
				refCount := rn.BetaSharingRegistry.GetJoinNodeRefCount(nodeID)
				rn.logger.Debug("   ✓ JoinNode %s preserved (%d rule(s) remaining)", nodeID, refCount)
			}
		} else {
			// No sharing registry - use lifecycle manager
			if err := rn.removeNodeWithCheck(nodeID, ruleID); err == nil {
				deletedCount++
				rn.logger.Debug("   🗑️  JoinNode %s removed", nodeID)
			}
		}
	}

	// Step 3: Remove alpha nodes
	for _, nodeID := range alphaNodes {
		if err := rn.removeNodeWithCheck(nodeID, ruleID); err == nil {
			deletedCount++
			rn.logger.Debug("   🗑️  AlphaNode %s removed", nodeID)
		} else {
			lifecycle, _ := rn.LifecycleManager.GetNodeLifecycle(nodeID)
			if lifecycle != nil && lifecycle.HasReferences() {
				rn.logger.Debug("   ✓ AlphaNode %s preserved (%d reference(s))", nodeID, lifecycle.GetRefCount())
			}
		}
	}

	// Step 4: Type nodes are usually shared - only remove if no references
	for _, nodeID := range typeNodes {
		lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(nodeID)
		if !exists {
			continue
		}

		shouldDelete, err := rn.LifecycleManager.RemoveRuleFromNode(nodeID, ruleID)
		if err != nil {
			rn.logger.Warn("   ⚠️  Error removing rule from type node %s: %v", nodeID, err)
			continue
		}

		if shouldDelete {
			if err := rn.removeNodeFromNetwork(nodeID); err == nil {
				deletedCount++
				rn.logger.Debug("   🗑️  TypeNode %s removed", nodeID)
			}
		} else {
			rn.logger.Debug("   ✓ TypeNode %s preserved (%d reference(s))", nodeID, lifecycle.GetRefCount())
		}
	}

	rn.logger.Info("✅ Rule %s removed successfully (%d node(s) deleted)", ruleID, deletedCount)
	return nil
}

// removeNodeWithCheck supprime un nœud seulement si RefCount == 0
func (rn *ReteNetwork) removeNodeWithCheck(nodeID, ruleID string) error {
	shouldDelete, err := rn.LifecycleManager.RemoveRuleFromNode(nodeID, ruleID)
	if err != nil {
		return err
	}

	if shouldDelete {
		return rn.removeNodeFromNetwork(nodeID)
	}

	return fmt.Errorf("nœud %s encore référencé", nodeID)
}

// removeNodeFromNetwork supprime un nœud du réseau RETE
// Ne supprime que si RefCount == 0
func (rn *ReteNetwork) removeNodeFromNetwork(nodeID string) error {
	// Vérifier que le nœud existe et peut être supprimé
	lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(nodeID)
	if !exists {
		return fmt.Errorf("nœud %s non trouvé dans le LifecycleManager", nodeID)
	}

	// Ne pas supprimer si le nœud a encore des références
	if lifecycle.HasReferences() {
		return fmt.Errorf("impossible de supprimer le nœud %s: encore %d référence(s)",
			nodeID, lifecycle.GetRefCount())
	}

	// Déterminer le type de nœud et le supprimer de la map appropriée
	switch lifecycle.NodeType {
	case "type":
		// Trouver et supprimer le TypeNode
		for typeName, typeNode := range rn.TypeNodes {
			if typeNode.GetID() == nodeID {
				// Déconnecter du RootNode
				rn.removeChildFromNode(rn.RootNode, typeNode)
				delete(rn.TypeNodes, typeName)
				return rn.LifecycleManager.RemoveNode(nodeID)
			}
		}

	case "alpha":
		// Supprimer l'AlphaNode
		if alphaNode, exists := rn.AlphaNodes[nodeID]; exists {
			// Déconnecter des parents (TypeNodes ou autres AlphaNodes)
			parent := rn.getChainParent(alphaNode)
			if parent != nil {
				rn.removeChildFromNode(parent, alphaNode)
				rn.logger.Debug("   🔗 AlphaNode %s déconnecté de son parent %s", nodeID, parent.GetID())
			}

			delete(rn.AlphaNodes, nodeID)

			// Supprimer du registre de partage AlphaSharingManager
			if rn.AlphaSharingManager != nil {
				// Vérifier si c'est un nœud partagé (les nœuds partagés ont un ID qui commence par "alpha_")
				if len(nodeID) > 6 && nodeID[:6] == "alpha_" {
					if err := rn.AlphaSharingManager.RemoveAlphaNode(nodeID); err != nil {
						rn.logger.Warn("   ⚠️  Erreur suppression AlphaNode du registre de partage: %v", err)
					} else {
						rn.logger.Debug("   ✓ AlphaNode %s supprimé du AlphaSharingManager", nodeID)
					}
				}
			}

			return rn.LifecycleManager.RemoveNode(nodeID)
		}

	case "terminal":
		// Supprimer le TerminalNode
		if terminalNode, exists := rn.TerminalNodes[nodeID]; exists {
			// Déconnecter des parents (AlphaNodes ou JoinNodes)
			for _, alphaNode := range rn.AlphaNodes {
				rn.removeChildFromNode(alphaNode, terminalNode)
			}
			// Aussi déconnecter des BetaNodes si nécessaire
			for _, betaNode := range rn.BetaNodes {
				if node, ok := betaNode.(Node); ok {
					rn.removeChildFromNode(node, terminalNode)
				}
			}
			delete(rn.TerminalNodes, nodeID)
			return rn.LifecycleManager.RemoveNode(nodeID)
		}

	case "join", "exists", "accumulate":
		// Supprimer le BetaNode
		if betaNode, exists := rn.BetaNodes[nodeID]; exists {
			// Déconnecter des parents
			for _, typeNode := range rn.TypeNodes {
				if node, ok := betaNode.(Node); ok {
					rn.removeChildFromNode(typeNode, node)
				}
			}
			delete(rn.BetaNodes, nodeID)
			return rn.LifecycleManager.RemoveNode(nodeID)
		}
	}

	return fmt.Errorf("nœud %s non trouvé dans le réseau", nodeID)
}

// removeJoinNodeFromNetwork removes a join node and all its dependent nodes from the network.
// This should only be called when the join node has no remaining rule references.
func (rn *ReteNetwork) removeJoinNodeFromNetwork(nodeID string) error {
	// Get the join node
	joinNode, exists := rn.BetaNodes[nodeID]
	if !exists {
		return fmt.Errorf("join node %s not found in network", nodeID)
	}

	rn.logger.Debug("   🗑️  Removing join node %s from network", nodeID)

	// Convert join node to proper type first
	var node Node
	var jn *JoinNode
	var ok bool
	if jn, ok = joinNode.(*JoinNode); !ok {
		return fmt.Errorf("beta node %s is not a JoinNode", nodeID)
	}
	node = jn

	// Step 1: Find and remove all terminal nodes that depend on this join node
	// Check if any terminal nodes are children of this join node
	for terminalID := range rn.TerminalNodes {
		// Check if this terminal is in the join node's children list
		isChild := false
		for _, child := range jn.GetChildren() {
			if child.GetID() == terminalID {
				isChild = true
				break
			}
		}

		if isChild {
			delete(rn.TerminalNodes, terminalID)
			rn.logger.Debug("   🗑️  Removed terminal node %s (child of join node)", terminalID)

			// Remove from lifecycle manager
			if rn.LifecycleManager != nil {
				rn.LifecycleManager.RemoveNode(terminalID)
			}
		}
	}

	// Step 2: Disconnect from parent nodes using the disconnectChild helper

	// Join nodes can have alpha nodes as parents
	for _, alphaNode := range rn.AlphaNodes {
		rn.disconnectChild(alphaNode, node)
	}

	// Check all other beta nodes (for cascading joins)
	for betaNodeID, betaNode := range rn.BetaNodes {
		if betaNodeID != nodeID {
			if bn, ok := betaNode.(*JoinNode); ok {
				rn.disconnectChild(bn, node)
			}
		}
	}

	// Also check type nodes (join nodes can connect directly to type nodes)
	for _, typeNode := range rn.TypeNodes {
		rn.disconnectChild(typeNode, node)
	}

	// Step 3: Remove from beta nodes map
	delete(rn.BetaNodes, nodeID)

	// Step 4: Remove from lifecycle manager
	if rn.LifecycleManager != nil {
		if err := rn.LifecycleManager.RemoveNode(nodeID); err != nil {
			rn.logger.Warn("   ⚠️  Warning: failed to remove join node %s from lifecycle manager: %v", nodeID, err)
		}
	}

	// Step 5: Remove from beta sharing registry
	if rn.BetaSharingRegistry != nil {
		if err := rn.BetaSharingRegistry.UnregisterJoinNode(nodeID); err != nil {
			rn.logger.Warn("   ⚠️  Warning: failed to unregister join node %s from beta sharing: %v", nodeID, err)
		}
	}

	rn.logger.Info("   ✅ Join node %s successfully removed from network", nodeID)
	return nil
}

// removeChildFromNode retire un nœud enfant d'un nœud parent
func (rn *ReteNetwork) removeChildFromNode(parent Node, child Node) {
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

	// Mettre à jour les enfants (nécessite un cast vers le type concret)
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

// disconnectChild removes a child from a node's children list
func (rn *ReteNetwork) disconnectChild(parent Node, child Node) {
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

	// Update parent's children list (this assumes BaseNode structure)
	if baseNode, ok := parent.(interface{ SetChildren([]Node) }); ok {
		baseNode.SetChildren(newChildren)
	}
}

// orderAlphaNodesReverse ordonne les AlphaNodes en ordre inverse de la chaîne
// (du nœud le plus éloigné du TypeNode vers le TypeNode)
func (rn *ReteNetwork) orderAlphaNodesReverse(alphaNodeIDs []string) []string {
	if len(alphaNodeIDs) <= 1 {
		return alphaNodeIDs
	}

	// Construire un graphe parent->enfants pour trouver l'ordre
	childToParent := make(map[string]string)
	hasParent := make(map[string]bool)

	for _, nodeID := range alphaNodeIDs {
		alphaNode, exists := rn.AlphaNodes[nodeID]
		if !exists {
			continue
		}

		parent := rn.getChainParent(alphaNode)
		if parent != nil {
			parentID := parent.GetID()
			// Vérifier si le parent est aussi un AlphaNode de notre liste
			for _, candidateID := range alphaNodeIDs {
				if candidateID == parentID {
					childToParent[nodeID] = parentID
					hasParent[nodeID] = true
					break
				}
			}
		}
	}

	// Trouver le nœud terminal de la chaîne (celui qui n'est parent de personne)
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

	// Si pas de structure de chaîne détectée, retourner l'ordre original
	if terminalNode == "" {
		return alphaNodeIDs
	}

	// Remonter la chaîne depuis le terminal
	ordered := make([]string, 0, len(alphaNodeIDs))
	current := terminalNode
	visited := make(map[string]bool)

	for current != "" && !visited[current] {
		ordered = append(ordered, current)
		visited[current] = true
		current = childToParent[current]
	}

	// Ajouter les nœuds non visités (au cas où)
	for _, nodeID := range alphaNodeIDs {
		if !visited[nodeID] {
			ordered = append(ordered, nodeID)
		}
	}

	return ordered
}

// isPartOfChain détecte si un nœud fait partie d'une chaîne d'AlphaNodes
func (rn *ReteNetwork) isPartOfChain(nodeID string) bool {
	lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(nodeID)
	if !exists || lifecycle.NodeType != "alpha" {
		return false
	}

	alphaNode, exists := rn.AlphaNodes[nodeID]
	if !exists {
		return false
	}

	// Un AlphaNode fait partie d'une chaîne si:
	// 1. Son parent est un autre AlphaNode, OU
	// 2. Un de ses enfants est un autre AlphaNode

	parent := rn.getChainParent(alphaNode)
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

// getChainParent récupère le nœud parent d'un AlphaNode dans une chaîne
func (rn *ReteNetwork) getChainParent(alphaNode *AlphaNode) Node {
	if alphaNode == nil {
		return nil
	}

	alphaID := alphaNode.GetID()

	// Chercher dans les TypeNodes
	for _, typeNode := range rn.TypeNodes {
		for _, child := range typeNode.GetChildren() {
			if child.GetID() == alphaID {
				return typeNode
			}
		}
	}

	// Chercher dans les autres AlphaNodes
	for _, parentAlpha := range rn.AlphaNodes {
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

// isJoinNode checks if a node ID corresponds to a JoinNode
func (rn *ReteNetwork) isJoinNode(nodeID string) bool {
	_, exists := rn.BetaNodes[nodeID]
	return exists
}
