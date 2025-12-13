// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"time"
)

type RootNode struct {
	BaseNode
}

// NewRootNode crée un nouveau nœud racine
func NewRootNode(storage Storage) *RootNode {
	return &RootNode{
		BaseNode: BaseNode{
			ID:        "root",
			Type:      "root",
			Memory:    &WorkingMemory{NodeID: "root", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children:  make([]Node, 0),
			Storage:   storage,
			createdAt: time.Now(),
		},
	}
}

// ActivateLeft (non utilisé pour le nœud racine)
func (rn *RootNode) ActivateLeft(token *Token) error {
	return fmt.Errorf("le nœud racine ne peut pas recevoir de tokens")
}

// ActivateRetract retire le fait de la mémoire racine et propage aux enfants
func (rn *RootNode) ActivateRetract(factID string) error {
	rn.mutex.Lock()
	rn.Memory.RemoveFact(factID)
	rn.mutex.Unlock()
	fmt.Printf("🗑️  [ROOT] Rétractation du fait: %s\n", factID)
	return rn.PropagateRetractToChildren(factID)
}

// ActivateRight distribue les faits aux nœuds de type
func (rn *RootNode) ActivateRight(fact *Fact) error {
	// Enregistrer l'activation
	rn.recordActivation()
	
	rn.mutex.Lock()
	if err := rn.Memory.AddFact(fact); err != nil {
		rn.mutex.Unlock()
		return fmt.Errorf("erreur ajout fait dans root node: %w", err)
	}
	rn.mutex.Unlock()

	// Log désactivé pour les performances
	// fmt.Printf("[ROOT] Reçu fait: %s\n", fact.String())

	// Persistance désactivée pour les performances

	// Propager aux enfants (TypeNodes)
	return rn.PropagateToChildren(fact, nil)
}
