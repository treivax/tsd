// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package rete fournit l'implémentation du réseau RETE pour l'évaluation de règles.
// Ce fichier contient les fonctions de statistiques, validation et introspection
// pour les chaînes alpha, permettant le monitoring et le debugging.
package rete

import (
	"fmt"
)

// GetChainInfo retourne des informations détaillées sur la chaîne alpha.
//
// Utile pour debugging, logging, et visualisation de la structure de la chaîne.
//
// Retourne:
//   - Map contenant: rule_id, node_count, node_ids, hashes, final_node_id
//   - Map avec clé "error" si chaîne nil
//
// Exemple:
//
//	info := chain.GetChainInfo()
//	fmt.Printf("Chaîne pour règle: %s\n", info["rule_id"])
//	fmt.Printf("Longueur: %d nœuds\n", info["node_count"])
//	fmt.Printf("IDs: %v\n", info["node_ids"])
//	fmt.Printf("Hashes: %v\n", info["hashes"])
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

// ValidateChain vérifie que la chaîne alpha est valide et cohérente.
//
// Vérifie:
//   - Chaîne non nil
//   - Au moins un nœud présent
//   - len(Nodes) == len(Hashes)
//   - FinalNode correspond au dernier élément de Nodes
//   - Tous les nœuds ont un ID non vide
//
// Retourne:
//   - nil si chaîne valide
//   - error décrivant le problème si invalide
//
// Exemple:
//
//	chain, err := builder.BuildChain(...)
//	if err := chain.ValidateChain(); err != nil {
//	    log.Fatalf("Chaîne invalide: %v", err)
//	}
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
	if chain == nil {
		return 0
	}

	sharedCount := 0
	// LifecycleManager is always initialized
	for _, node := range chain.Nodes {
		if lifecycle, exists := acb.network.LifecycleManager.GetNodeLifecycle(node.ID); exists {
			if lifecycle.GetRefCount() > 1 {
				sharedCount++
			}
		}
	}
	return sharedCount
}

// GetChainStats retourne des statistiques détaillées sur une chaîne alpha.
//
// Calcule et retourne:
//   - chain_length: Nombre total de nœuds dans la chaîne
//   - shared_nodes: Nœuds avec RefCount > 1
//   - new_nodes: Nœuds avec RefCount == 1
//   - sharing_ratio: Proportion de nœuds partagés (0.0 à 1.0)
//   - node_details: Liste des infos par nœud (ID, RefCount, is_shared)
//
// Paramètres:
//   - chain: Chaîne alpha à analyser
//
// Retourne:
//   - Map avec statistiques détaillées
//   - Map avec clé "error" si chaîne nil
//
// Exemple:
//
//	chain, _ := builder.BuildChain(...)
//	stats := builder.GetChainStats(chain)
//	fmt.Printf("Longueur: %d\n", stats["chain_length"])
//	fmt.Printf("Partagés: %d\n", stats["shared_nodes"])
//	fmt.Printf("Nouveaux: %d\n", stats["new_nodes"])
//	fmt.Printf("Ratio: %.1f%%\n", stats["sharing_ratio"].(float64) * 100)
//
//	// Détails par nœud
//	for _, detail := range stats["node_details"].([]map[string]interface{}) {
//	    fmt.Printf("  Nœud %s: RefCount=%d, Partagé=%v\n",
//	        detail["node_id"], detail["ref_count"], detail["is_shared"])
//	}
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
