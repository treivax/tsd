// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sync"
)

// RuleReference représente une référence à une règle utilisant un nœud
type RuleReference struct {
	RuleID   string `json:"rule_id"`
	RuleName string `json:"rule_name,omitempty"`
}

// NodeLifecycle gère le cycle de vie d'un nœud avec tracking des règles qui l'utilisent
type NodeLifecycle struct {
	NodeID         string                    `json:"node_id"`
	NodeType       string                    `json:"node_type"`
	Rules          map[string]*RuleReference `json:"rules"` // Map[RuleID] -> RuleReference
	RefCount       int                       `json:"ref_count"`
	CreatedByRules []string                  `json:"created_by_rules,omitempty"`
	mutex          sync.RWMutex              `json:"-"`
}

// NewNodeLifecycle crée une nouvelle instance de gestion du cycle de vie
func NewNodeLifecycle(nodeID, nodeType string) *NodeLifecycle {
	return &NodeLifecycle{
		NodeID:         nodeID,
		NodeType:       nodeType,
		Rules:          make(map[string]*RuleReference),
		RefCount:       0,
		CreatedByRules: make([]string, 0),
	}
}

// AddRuleReference ajoute une référence de règle à ce nœud
func (nl *NodeLifecycle) AddRuleReference(ruleID string, ruleName string) {
	nl.mutex.Lock()
	defer nl.mutex.Unlock()

	if _, exists := nl.Rules[ruleID]; !exists {
		nl.Rules[ruleID] = &RuleReference{
			RuleID:   ruleID,
			RuleName: ruleName,
		}
		nl.RefCount++
		nl.CreatedByRules = append(nl.CreatedByRules, ruleID)
	}
}

// RemoveRuleReference retire une référence de règle de ce nœud
func (nl *NodeLifecycle) RemoveRuleReference(ruleID string) bool {
	nl.mutex.Lock()
	defer nl.mutex.Unlock()

	if _, exists := nl.Rules[ruleID]; exists {
		delete(nl.Rules, ruleID)
		nl.RefCount--

		// Retirer de la liste des règles créatrices
		for i, id := range nl.CreatedByRules {
			if id == ruleID {
				nl.CreatedByRules = append(nl.CreatedByRules[:i], nl.CreatedByRules[i+1:]...)
				break
			}
		}

		return nl.RefCount == 0 // Retourne true si plus aucune référence
	}

	return false
}

// HasReferences retourne true si le nœud a encore des références
func (nl *NodeLifecycle) HasReferences() bool {
	nl.mutex.RLock()
	defer nl.mutex.RUnlock()
	return nl.RefCount > 0
}

// GetRefCount retourne le nombre de références
func (nl *NodeLifecycle) GetRefCount() int {
	nl.mutex.RLock()
	defer nl.mutex.RUnlock()
	return nl.RefCount
}

// GetRules retourne la liste des IDs de règles référençant ce nœud
func (nl *NodeLifecycle) GetRules() []string {
	nl.mutex.RLock()
	defer nl.mutex.RUnlock()

	rules := make([]string, 0, len(nl.Rules))
	for ruleID := range nl.Rules {
		rules = append(rules, ruleID)
	}
	return rules
}

// GetRuleInfo retourne les informations d'une règle spécifique
func (nl *NodeLifecycle) GetRuleInfo(ruleID string) (*RuleReference, bool) {
	nl.mutex.RLock()
	defer nl.mutex.RUnlock()

	ref, exists := nl.Rules[ruleID]
	return ref, exists
}

// LifecycleManager gère le cycle de vie de tous les nœuds du réseau
type LifecycleManager struct {
	Nodes map[string]*NodeLifecycle `json:"nodes"` // Map[NodeID] -> NodeLifecycle
	mutex sync.RWMutex              `json:"-"`
}

// NewLifecycleManager crée un nouveau gestionnaire de cycle de vie
func NewLifecycleManager() *LifecycleManager {
	return &LifecycleManager{
		Nodes: make(map[string]*NodeLifecycle),
	}
}

// RegisterNode enregistre un nouveau nœud dans le gestionnaire
func (lm *LifecycleManager) RegisterNode(nodeID, nodeType string) *NodeLifecycle {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	if lifecycle, exists := lm.Nodes[nodeID]; exists {
		return lifecycle
	}

	lifecycle := NewNodeLifecycle(nodeID, nodeType)
	lm.Nodes[nodeID] = lifecycle
	return lifecycle
}

// AddRuleToNode associe une règle à un nœud
func (lm *LifecycleManager) AddRuleToNode(nodeID, ruleID, ruleName string) error {
	lm.mutex.RLock()
	lifecycle, exists := lm.Nodes[nodeID]
	lm.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("nœud %s non enregistré dans le lifecycle manager", nodeID)
	}

	lifecycle.AddRuleReference(ruleID, ruleName)
	return nil
}

// RemoveRuleFromNode retire une règle d'un nœud
// Retourne true si le nœud n'a plus de références et peut être supprimé
func (lm *LifecycleManager) RemoveRuleFromNode(nodeID, ruleID string) (bool, error) {
	lm.mutex.RLock()
	lifecycle, exists := lm.Nodes[nodeID]
	lm.mutex.RUnlock()

	if !exists {
		return false, fmt.Errorf("nœud %s non trouvé", nodeID)
	}

	shouldDelete := lifecycle.RemoveRuleReference(ruleID)
	return shouldDelete, nil
}

// RemoveNode supprime complètement un nœud du gestionnaire
func (lm *LifecycleManager) RemoveNode(nodeID string) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	if lifecycle, exists := lm.Nodes[nodeID]; exists {
		if lifecycle.HasReferences() {
			return fmt.Errorf("impossible de supprimer le nœud %s: encore %d référence(s)",
				nodeID, lifecycle.RefCount)
		}
		delete(lm.Nodes, nodeID)
		return nil
	}

	return fmt.Errorf("nœud %s non trouvé", nodeID)
}

// GetNodeLifecycle récupère le lifecycle d'un nœud
func (lm *LifecycleManager) GetNodeLifecycle(nodeID string) (*NodeLifecycle, bool) {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()

	lifecycle, exists := lm.Nodes[nodeID]
	return lifecycle, exists
}

// GetNodesForRule retourne tous les nœuds utilisés par une règle donnée
func (lm *LifecycleManager) GetNodesForRule(ruleID string) []string {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()

	nodes := make([]string, 0)
	for nodeID, lifecycle := range lm.Nodes {
		if _, exists := lifecycle.Rules[ruleID]; exists {
			nodes = append(nodes, nodeID)
		}
	}
	return nodes
}

// CanRemoveNode vérifie si un nœud peut être supprimé (pas de références)
func (lm *LifecycleManager) CanRemoveNode(nodeID string) bool {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()

	if lifecycle, exists := lm.Nodes[nodeID]; exists {
		return !lifecycle.HasReferences()
	}
	return false
}

// GetStats retourne des statistiques sur le gestionnaire
func (lm *LifecycleManager) GetStats() map[string]interface{} {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()

	totalNodes := len(lm.Nodes)
	totalRefs := 0
	nodesWithoutRefs := 0

	for _, lifecycle := range lm.Nodes {
		totalRefs += lifecycle.RefCount
		if lifecycle.RefCount == 0 {
			nodesWithoutRefs++
		}
	}

	return map[string]interface{}{
		"total_nodes":        totalNodes,
		"total_references":   totalRefs,
		"nodes_without_refs": nodesWithoutRefs,
		"nodes_with_refs":    totalNodes - nodesWithoutRefs,
	}
}

// Reset réinitialise complètement le gestionnaire
func (lm *LifecycleManager) Reset() {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	lm.Nodes = make(map[string]*NodeLifecycle)
}

// RuleInfo contient les informations d'une règle et ses nœuds
type RuleInfo struct {
	RuleID    string   `json:"rule_id"`
	RuleName  string   `json:"rule_name,omitempty"`
	NodeIDs   []string `json:"node_ids"`
	NodeCount int      `json:"node_count"`
}

// GetRuleInfo retourne les informations complètes d'une règle
func (lm *LifecycleManager) GetRuleInfo(ruleID string) *RuleInfo {
	nodes := lm.GetNodesForRule(ruleID)

	info := &RuleInfo{
		RuleID:    ruleID,
		NodeIDs:   nodes,
		NodeCount: len(nodes),
	}

	// Essayer de récupérer le nom de la règle depuis le premier nœud trouvé
	if len(nodes) > 0 {
		if lifecycle, exists := lm.GetNodeLifecycle(nodes[0]); exists {
			if ref, ok := lifecycle.GetRuleInfo(ruleID); ok {
				info.RuleName = ref.RuleName
			}
		}
	}

	return info
}


// Cleanup nettoie toutes les ressources du LifecycleManager
func (lm *LifecycleManager) Cleanup() {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	
	// Vider toutes les nodes
	lm.Nodes = make(map[string]*NodeLifecycle)
}
