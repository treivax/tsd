// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

// ReteNetwork représente le réseau RETE complet
type ReteNetwork struct {
	RootNode            *RootNode                `json:"root_node"`
	TypeNodes           map[string]*TypeNode     `json:"type_nodes"`
	AlphaNodes          map[string]*AlphaNode    `json:"alpha_nodes"`
	BetaNodes           map[string]interface{}   `json:"beta_nodes"` // Nœuds Beta pour les jointures multi-faits
	TerminalNodes       map[string]*TerminalNode `json:"terminal_nodes"`
	Storage             Storage                  `json:"-"`
	Types               []TypeDefinition         `json:"types"`
	BetaBuilder         interface{}              `json:"-"` // Constructeur de réseau Beta
	LifecycleManager    *LifecycleManager        `json:"-"` // Gestionnaire du cycle de vie des nœuds
	AlphaSharingManager *AlphaSharingRegistry    `json:"-"` // Gestionnaire du partage des AlphaNodes
}

// NewReteNetwork crée un nouveau réseau RETE
func NewReteNetwork(storage Storage) *ReteNetwork {
	rootNode := NewRootNode(storage)

	return &ReteNetwork{
		RootNode:            rootNode,
		TypeNodes:           make(map[string]*TypeNode),
		AlphaNodes:          make(map[string]*AlphaNode),
		BetaNodes:           make(map[string]interface{}),
		TerminalNodes:       make(map[string]*TerminalNode),
		Storage:             storage,
		Types:               make([]TypeDefinition, 0),
		BetaBuilder:         nil, // Sera initialisé si nécessaire
		LifecycleManager:    NewLifecycleManager(),
		AlphaSharingManager: NewAlphaSharingRegistry(),
	}
}

// SubmitFact soumet un nouveau fait au réseau
func (rn *ReteNetwork) SubmitFact(fact *Fact) error {
	fmt.Printf("🔥 Soumission d'un nouveau fait au réseau RETE: %s\n", fact.String())

	// Propager le fait depuis le nœud racine
	return rn.RootNode.ActivateRight(fact)
}

// SubmitFactsFromGrammar soumet plusieurs faits depuis la grammaire au réseau
func (rn *ReteNetwork) SubmitFactsFromGrammar(facts []map[string]interface{}) error {
	for i, factMap := range facts {
		// Convertir le map en Fact
		factID := fmt.Sprintf("fact_%d", i)
		if id, ok := factMap["id"].(string); ok {
			factID = id
		}

		factType := "unknown"
		if typ, ok := factMap["type"].(string); ok {
			factType = typ
		}

		fact := &Fact{
			ID:     factID,
			Type:   factType,
			Fields: make(map[string]interface{}),
		}

		// Copier tous les champs
		for key, value := range factMap {
			if key != "id" && key != "type" {
				fact.Fields[key] = value
			}
		}

		if err := rn.SubmitFact(fact); err != nil {
			return fmt.Errorf("erreur soumission fait %s: %w", fact.ID, err)
		}
	}
	return nil
}

// RetractFact retire un fait du réseau et propage la rétractation
// factID doit être l'identifiant interne (Type_ID)
func (rn *ReteNetwork) RetractFact(factID string) error {
	fmt.Printf("🗑️  Rétractation du fait: %s\n", factID)

	// Vérifier que le fait existe dans le réseau
	memory := rn.RootNode.GetMemory()
	if _, exists := memory.GetFact(factID); !exists {
		return fmt.Errorf("fait %s introuvable dans le réseau", factID)
	}

	// Propager la rétractation depuis le nœud racine
	return rn.RootNode.ActivateRetract(factID)
}

// Reset clears the entire RETE network and resets it to an empty state.
// This removes all facts, rules, types, and network nodes.
// After calling Reset, the network is ready to accept new definitions from scratch.
func (rn *ReteNetwork) Reset() {
	fmt.Println("🧹 Réinitialisation complète du réseau RETE")

	// Clear all node collections
	rn.TypeNodes = make(map[string]*TypeNode)
	rn.AlphaNodes = make(map[string]*AlphaNode)
	rn.BetaNodes = make(map[string]interface{})
	rn.TerminalNodes = make(map[string]*TerminalNode)
	rn.Types = make([]TypeDefinition, 0)
	rn.BetaBuilder = nil

	// Reset lifecycle manager
	if rn.LifecycleManager != nil {
		rn.LifecycleManager.Reset()
	} else {
		rn.LifecycleManager = NewLifecycleManager()
	}

	// Reset alpha sharing manager
	if rn.AlphaSharingManager != nil {
		rn.AlphaSharingManager.Reset()
	} else {
		rn.AlphaSharingManager = NewAlphaSharingRegistry()
	}

	// Recreate a fresh root node with the existing storage
	rn.RootNode = NewRootNode(rn.Storage)

	fmt.Println("✅ Réseau RETE réinitialisé avec succès")
}

// RemoveRule supprime une règle et tous ses nœuds qui ne sont plus utilisés
func (rn *ReteNetwork) RemoveRule(ruleID string) error {
	fmt.Printf("🗑️  Suppression de la règle: %s\n", ruleID)

	if rn.LifecycleManager == nil {
		return fmt.Errorf("LifecycleManager non initialisé")
	}

	// Récupérer tous les nœuds utilisés par cette règle
	nodeIDs := rn.LifecycleManager.GetNodesForRule(ruleID)
	if len(nodeIDs) == 0 {
		return fmt.Errorf("règle %s non trouvée ou aucun nœud associé", ruleID)
	}

	fmt.Printf("   📊 Nœuds associés à la règle: %d\n", len(nodeIDs))

	// Parcourir chaque nœud et retirer la référence à la règle
	nodesToDelete := make([]string, 0)
	for _, nodeID := range nodeIDs {
		shouldDelete, err := rn.LifecycleManager.RemoveRuleFromNode(nodeID, ruleID)
		if err != nil {
			fmt.Printf("   ⚠️  Erreur lors de la suppression de la règle du nœud %s: %v\n", nodeID, err)
			continue
		}

		if shouldDelete {
			nodesToDelete = append(nodesToDelete, nodeID)
			fmt.Printf("   ✓ Nœud %s marqué pour suppression (plus de références)\n", nodeID)
		} else {
			lifecycle, _ := rn.LifecycleManager.GetNodeLifecycle(nodeID)
			fmt.Printf("   ✓ Nœud %s conservé (%d référence(s) restante(s))\n", nodeID, lifecycle.GetRefCount())
		}
	}

	// Supprimer les nœuds qui n'ont plus de références
	for _, nodeID := range nodesToDelete {
		if err := rn.removeNodeFromNetwork(nodeID); err != nil {
			fmt.Printf("   ⚠️  Erreur lors de la suppression du nœud %s: %v\n", nodeID, err)
		} else {
			fmt.Printf("   🗑️  Nœud %s supprimé du réseau\n", nodeID)
		}
	}

	fmt.Printf("✅ Règle %s supprimée avec succès (%d nœud(s) supprimé(s))\n", ruleID, len(nodesToDelete))
	return nil
}

// removeNodeFromNetwork supprime un nœud du réseau RETE
func (rn *ReteNetwork) removeNodeFromNetwork(nodeID string) error {
	// Déterminer le type de nœud et le supprimer de la map appropriée
	if lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(nodeID); exists {
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
				// Déconnecter des parents (TypeNodes ou autres)
				for _, typeNode := range rn.TypeNodes {
					rn.removeChildFromNode(typeNode, alphaNode)
				}
				delete(rn.AlphaNodes, nodeID)

				// Supprimer du registre de partage
				// Le nodeID est le hash de la condition pour les nœuds partagés
				if rn.AlphaSharingManager != nil {
					// Vérifier si c'est un nœud partagé (les nœuds partagés ont un ID qui commence par "alpha_")
					if len(nodeID) > 6 && nodeID[:6] == "alpha_" {
						if err := rn.AlphaSharingManager.RemoveAlphaNode(nodeID); err != nil {
							fmt.Printf("   ⚠️  Erreur suppression AlphaNode du registre de partage: %v\n", err)
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
	}

	return fmt.Errorf("nœud %s non trouvé dans le réseau", nodeID)
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

// GetRuleInfo retourne les informations d'une règle
func (rn *ReteNetwork) GetRuleInfo(ruleID string) *RuleInfo {
	if rn.LifecycleManager == nil {
		return &RuleInfo{
			RuleID:    ruleID,
			NodeIDs:   []string{},
			NodeCount: 0,
		}
	}
	return rn.LifecycleManager.GetRuleInfo(ruleID)
}

// GetNetworkStats retourne des statistiques sur le réseau
func (rn *ReteNetwork) GetNetworkStats() map[string]interface{} {
	stats := map[string]interface{}{
		"type_nodes":     len(rn.TypeNodes),
		"alpha_nodes":    len(rn.AlphaNodes),
		"beta_nodes":     len(rn.BetaNodes),
		"terminal_nodes": len(rn.TerminalNodes),
	}

	if rn.LifecycleManager != nil {
		lifecycleStats := rn.LifecycleManager.GetStats()
		for k, v := range lifecycleStats {
			stats["lifecycle_"+k] = v
		}
	}

	if rn.AlphaSharingManager != nil {
		alphaStats := rn.AlphaSharingManager.GetStats()
		for k, v := range alphaStats {
			stats["sharing_"+k] = v
		}
	}

	return stats
}
