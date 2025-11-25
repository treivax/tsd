package rete

import (
	"fmt"
)

type RootNode struct {
	BaseNode
}

// NewRootNode crée un nouveau nœud racine
func NewRootNode(storage Storage) *RootNode {
	return &RootNode{
		BaseNode: BaseNode{
			ID:       "root",
			Type:     "root",
			Memory:   &WorkingMemory{NodeID: "root", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children: make([]Node, 0),
			Storage:  storage,
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
	rn.mutex.Lock()
	rn.Memory.AddFact(fact)
	rn.mutex.Unlock()

	// Log désactivé pour les performances
	// fmt.Printf("[ROOT] Reçu fait: %s\n", fact.String())

	// Persistance désactivée pour les performances

	// Propager aux enfants (TypeNodes)
	return rn.PropagateToChildren(fact, nil)
}
