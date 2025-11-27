// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"log"
)

// AlphaChain représente une chaîne d'AlphaNodes construite pour un ensemble de conditions
type AlphaChain struct {
	Nodes     []*AlphaNode `json:"nodes"`      // Liste ordonnée des nœuds alpha dans la chaîne
	Hashes    []string     `json:"hashes"`     // Hashes correspondants pour chaque nœud
	FinalNode *AlphaNode   `json:"final_node"` // Le dernier nœud de la chaîne
	RuleID    string       `json:"rule_id"`    // ID de la règle pour laquelle la chaîne a été construite
}

// AlphaChainBuilder construit des chaînes d'AlphaNodes avec partage automatique
type AlphaChainBuilder struct {
	network *ReteNetwork
	storage Storage
}

// NewAlphaChainBuilder crée un nouveau constructeur de chaînes alpha
func NewAlphaChainBuilder(network *ReteNetwork, storage Storage) *AlphaChainBuilder {
	return &AlphaChainBuilder{
		network: network,
		storage: storage,
	}
}

// BuildChain construit une chaîne d'AlphaNodes pour un ensemble de conditions
// avec partage automatique des nœuds identiques entre règles.
//
// Paramètres:
//   - conditions: liste de conditions simples dans l'ordre normalisé
//   - variableName: nom de la variable à laquelle les conditions s'appliquent
//   - parentNode: nœud parent auquel connecter le premier nœud alpha
//   - ruleID: identifiant de la règle pour le tracking du cycle de vie
//
// Retourne:
//   - *AlphaChain: la chaîne construite avec tous les nœuds
//   - error: erreur éventuelle lors de la construction
func (acb *AlphaChainBuilder) BuildChain(
	conditions []SimpleCondition,
	variableName string,
	parentNode Node,
	ruleID string,
) (*AlphaChain, error) {
	if len(conditions) == 0 {
		return nil, fmt.Errorf("impossible de construire une chaîne sans conditions")
	}

	if parentNode == nil {
		return nil, fmt.Errorf("le nœud parent ne peut pas être nil")
	}

	if acb.network.AlphaSharingManager == nil {
		return nil, fmt.Errorf("AlphaSharingManager non initialisé dans le réseau")
	}

	if acb.network.LifecycleManager == nil {
		return nil, fmt.Errorf("LifecycleManager non initialisé dans le réseau")
	}

	chain := &AlphaChain{
		Nodes:  make([]*AlphaNode, 0, len(conditions)),
		Hashes: make([]string, 0, len(conditions)),
		RuleID: ruleID,
	}

	currentParent := parentNode

	// Construire la chaîne condition par condition
	for i, condition := range conditions {
		// Convertir SimpleCondition en map pour la condition du nœud alpha
		conditionMap := map[string]interface{}{
			"type":     condition.Type,
			"left":     condition.Left,
			"operator": condition.Operator,
			"right":    condition.Right,
		}

		// Obtenir ou créer l'AlphaNode via le gestionnaire de partage
		alphaNode, hash, reused, err := acb.network.AlphaSharingManager.GetOrCreateAlphaNode(
			conditionMap,
			variableName,
			acb.storage,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la création/récupération du nœud alpha %d: %w", i, err)
		}

		// Ajouter le nœud et son hash à la chaîne
		chain.Nodes = append(chain.Nodes, alphaNode)
		chain.Hashes = append(chain.Hashes, hash)

		if reused {
			// Nœud réutilisé - vérifier la connexion au parent
			log.Printf("♻️  [AlphaChainBuilder] Réutilisation du nœud alpha %s pour la règle %s (condition %d/%d)",
				alphaNode.ID, ruleID, i+1, len(conditions))

			if !isAlreadyConnected(currentParent, alphaNode) {
				// Connecter au parent si pas déjà connecté
				currentParent.AddChild(alphaNode)
				log.Printf("🔗 [AlphaChainBuilder] Connexion du nœud réutilisé %s au parent %s",
					alphaNode.ID, currentParent.GetID())
			} else {
				log.Printf("✓  [AlphaChainBuilder] Nœud %s déjà connecté au parent %s",
					alphaNode.ID, currentParent.GetID())
			}
		} else {
			// Nouveau nœud - le connecter au parent et l'ajouter au réseau
			currentParent.AddChild(alphaNode)
			acb.network.AlphaNodes[alphaNode.ID] = alphaNode

			log.Printf("🆕 [AlphaChainBuilder] Nouveau nœud alpha %s créé pour la règle %s (condition %d/%d)",
				alphaNode.ID, ruleID, i+1, len(conditions))
			log.Printf("🔗 [AlphaChainBuilder] Connexion du nœud %s au parent %s",
				alphaNode.ID, currentParent.GetID())
		}

		// Enregistrer le nœud dans le LifecycleManager avec la règle
		lifecycle := acb.network.LifecycleManager.RegisterNode(alphaNode.ID, "alpha")
		lifecycle.AddRuleReference(ruleID, "") // RuleName peut être ajouté plus tard si nécessaire

		if reused {
			log.Printf("📊 [AlphaChainBuilder] Nœud %s maintenant utilisé par %d règle(s)",
				alphaNode.ID, lifecycle.GetRefCount())
		}

		// Le nœud actuel devient le parent pour le prochain nœud
		currentParent = alphaNode
	}

	// Le dernier nœud de la chaîne est le nœud final
	chain.FinalNode = chain.Nodes[len(chain.Nodes)-1]

	log.Printf("✅ [AlphaChainBuilder] Chaîne alpha complète construite pour la règle %s: %d nœud(s)",
		ruleID, len(chain.Nodes))

	return chain, nil
}

// isAlreadyConnected vérifie si un nœud enfant est déjà connecté à un nœud parent
func isAlreadyConnected(parent Node, child Node) bool {
	if parent == nil || child == nil {
		return false
	}

	children := parent.GetChildren()
	childID := child.GetID()

	for _, c := range children {
		if c.GetID() == childID {
			return true
		}
	}

	return false
}

// GetChainInfo retourne des informations sur la chaîne alpha
func (ac *AlphaChain) GetChainInfo() map[string]interface{} {
	if ac == nil {
		return map[string]interface{}{
			"error": "chain is nil",
		}
	}

	nodeIDs := make([]string, len(ac.Nodes))
	for i, node := range ac.Nodes {
		nodeIDs[i] = node.ID
	}

	finalNodeID := ""
	if ac.FinalNode != nil {
		finalNodeID = ac.FinalNode.ID
	}

	return map[string]interface{}{
		"rule_id":       ac.RuleID,
		"node_count":    len(ac.Nodes),
		"node_ids":      nodeIDs,
		"hashes":        ac.Hashes,
		"final_node_id": finalNodeID,
	}
}

// ValidateChain vérifie que la chaîne est valide et cohérente
func (ac *AlphaChain) ValidateChain() error {
	if ac == nil {
		return fmt.Errorf("chaîne alpha nil")
	}

	if len(ac.Nodes) == 0 {
		return fmt.Errorf("chaîne alpha vide")
	}

	if len(ac.Nodes) != len(ac.Hashes) {
		return fmt.Errorf("incohérence: %d nœuds mais %d hashes", len(ac.Nodes), len(ac.Hashes))
	}

	if ac.FinalNode == nil {
		return fmt.Errorf("nœud final nil")
	}

	// Vérifier que le nœud final est bien le dernier de la liste
	if ac.FinalNode != ac.Nodes[len(ac.Nodes)-1] {
		return fmt.Errorf("le nœud final ne correspond pas au dernier nœud de la liste")
	}

	// Vérifier que tous les nœuds sont non-nil
	for i, node := range ac.Nodes {
		if node == nil {
			return fmt.Errorf("nœud %d est nil", i)
		}
	}

	return nil
}

// CountSharedNodes retourne le nombre de nœuds partagés dans la chaîne
// (nœuds avec plus d'une référence dans le LifecycleManager)
func (acb *AlphaChainBuilder) CountSharedNodes(chain *AlphaChain) int {
	if chain == nil || acb.network.LifecycleManager == nil {
		return 0
	}

	sharedCount := 0
	for _, node := range chain.Nodes {
		if lifecycle, exists := acb.network.LifecycleManager.GetNodeLifecycle(node.ID); exists {
			if lifecycle.GetRefCount() > 1 {
				sharedCount++
			}
		}
	}

	return sharedCount
}

// GetChainStats retourne des statistiques détaillées sur la chaîne
func (acb *AlphaChainBuilder) GetChainStats(chain *AlphaChain) map[string]interface{} {
	if chain == nil {
		return map[string]interface{}{
			"error": "chain is nil",
		}
	}

	sharedNodes := acb.CountSharedNodes(chain)
	newNodes := len(chain.Nodes) - sharedNodes

	stats := map[string]interface{}{
		"total_nodes":  len(chain.Nodes),
		"shared_nodes": sharedNodes,
		"new_nodes":    newNodes,
		"rule_id":      chain.RuleID,
	}

	// Ajouter les détails de chaque nœud
	nodeDetails := make([]map[string]interface{}, len(chain.Nodes))
	for i, node := range chain.Nodes {
		refCount := 0
		if lifecycle, exists := acb.network.LifecycleManager.GetNodeLifecycle(node.ID); exists {
			refCount = lifecycle.GetRefCount()
		}

		nodeDetails[i] = map[string]interface{}{
			"index":     i,
			"node_id":   node.ID,
			"hash":      chain.Hashes[i],
			"ref_count": refCount,
			"is_shared": refCount > 1,
			"is_final":  node == chain.FinalNode,
		}
	}
	stats["node_details"] = nodeDetails

	return stats
}
