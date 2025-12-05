// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

// BetaChain représente une chaîne de JoinNodes construite pour un ensemble de patterns.
//
// Une chaîne beta est une séquence ordonnée de nœuds de jointure qui combinent
// progressivement plusieurs variables. Chaque JoinNode évalue une condition de
// jointure et propage les tokens combinés au nœud suivant dans la chaîne.
//
// Structure de chaîne typique (cascade pour 3+ variables):
//
//	TypeNode(Person)  TypeNode(Order)
//	       └──────────┴────────┐
//	                    JoinNode(p ⋈ o)
//	                         │
//	                    TypeNode(Payment)
//	                         │
//	                    JoinNode((p,o) ⋈ pay)
//	                         │
//	                    TerminalNode(rule_terminal)
//
// Propriétés:
//   - len(Nodes) == len(Hashes) (toujours maintenu)
//   - FinalNode == Nodes[len(Nodes)-1] (si non vide)
//   - Ordre des jointures est optimisé pour la sélectivité
//
// Exemple d'utilisation:
//
//	patterns := []JoinPattern{
//	    {LeftVar: "p", RightVar: "o", Condition: ...},
//	    {LeftVar: "p,o", RightVar: "pay", Condition: ...},
//	}
//	chain, err := builder.BuildChain(patterns, "myRule")
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Chaîne construite: %d nœuds\n", len(chain.Nodes))
type BetaChain struct {
	Nodes     []*JoinNode `json:"nodes"`      // Liste ordonnée des JoinNodes dans la chaîne
	Hashes    []string    `json:"hashes"`     // Hashes correspondants pour chaque nœud
	FinalNode *JoinNode   `json:"final_node"` // Le dernier nœud de la chaîne
	RuleID    string      `json:"rule_id"`    // ID de la règle pour laquelle la chaîne a été construite
}

// GetChainInfo retourne des informations détaillées sur la chaîne.
//
// Les informations incluent:
//   - rule_id: ID de la règle associée
//   - chain_length: nombre de nœuds dans la chaîne
//   - node_ids: liste des IDs de nœuds
//   - hashes: liste des hashes de nœuds
//   - final_node_id: ID du nœud final
//   - summary: résumé textuel de la chaîne
//
// Exemple:
//
//	info := chain.GetChainInfo()
//	fmt.Printf("Chaîne: %s\n", info["summary"])
//	fmt.Printf("Longueur: %d\n", info["chain_length"])
func (bc *BetaChain) GetChainInfo() map[string]interface{} {
	info := make(map[string]interface{})

	nodeIDs := make([]string, len(bc.Nodes))
	for i, node := range bc.Nodes {
		nodeIDs[i] = node.ID
	}

	info["rule_id"] = bc.RuleID
	info["chain_length"] = len(bc.Nodes)
	info["node_ids"] = nodeIDs
	info["hashes"] = bc.Hashes

	if bc.FinalNode != nil {
		info["final_node_id"] = bc.FinalNode.ID
	}

	summary := fmt.Sprintf("BetaChain[%s]: %d nœuds", bc.RuleID, len(bc.Nodes))
	info["summary"] = summary

	return info
}

// ValidateChain valide la cohérence d'une chaîne beta.
//
// Vérifie:
//   - Longueurs cohérentes (nodes, hashes)
//   - FinalNode correspond au dernier nœud
//   - Tous les nœuds sont non-nil
//   - Tous les hashes sont non-vides
//
// Retourne une erreur si la validation échoue.
//
// Exemple:
//
//	if err := chain.ValidateChain(); err != nil {
//	    fmt.Printf("Chaîne invalide: %v", err)
//	}
func (bc *BetaChain) ValidateChain() error {
	if len(bc.Nodes) != len(bc.Hashes) {
		return fmt.Errorf("incohérence: %d nœuds mais %d hashes", len(bc.Nodes), len(bc.Hashes))
	}

	if len(bc.Nodes) == 0 {
		return fmt.Errorf("chaîne vide")
	}

	for i, node := range bc.Nodes {
		if node == nil {
			return fmt.Errorf("nœud %d est nil", i)
		}
		if bc.Hashes[i] == "" {
			return fmt.Errorf("hash %d est vide", i)
		}
	}

	if bc.FinalNode != bc.Nodes[len(bc.Nodes)-1] {
		return fmt.Errorf("FinalNode ne correspond pas au dernier nœud de la chaîne")
	}

	return nil
}
