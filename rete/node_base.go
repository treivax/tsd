// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sync"
)

type BaseNode struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	Memory   *WorkingMemory `json:"memory"`
	Children []Node         `json:"children"`
	Storage  Storage        `json:"-"`
	network  *ReteNetwork   `json:"-"` // Référence au réseau RETE parent
	mutex    sync.RWMutex   `json:"-"`
}

// GetID retourne l'ID du nœud
func (bn *BaseNode) GetID() string {
	return bn.ID
}

// GetType retourne le type du nœud
func (bn *BaseNode) GetType() string {
	return bn.Type
}

// GetMemory retourne la mémoire de travail du nœud
func (bn *BaseNode) GetMemory() *WorkingMemory {
	bn.mutex.RLock()
	defer bn.mutex.RUnlock()
	return bn.Memory
}

// SetNetwork définit la référence au réseau RETE parent
func (bn *BaseNode) SetNetwork(network *ReteNetwork) {
	bn.mutex.Lock()
	defer bn.mutex.Unlock()
	bn.network = network
}

// GetNetwork retourne la référence au réseau RETE parent
func (bn *BaseNode) GetNetwork() *ReteNetwork {
	bn.mutex.RLock()
	defer bn.mutex.RUnlock()
	return bn.network
}

// AddChild ajoute un nœud enfant
func (bn *BaseNode) AddChild(child Node) {
	bn.mutex.Lock()
	defer bn.mutex.Unlock()
	bn.Children = append(bn.Children, child)
}

// GetChildren retourne les nœuds enfants
func (bn *BaseNode) GetChildren() []Node {
	bn.mutex.RLock()
	defer bn.mutex.RUnlock()
	return bn.Children
}

// SetChildren sets the children nodes
func (bn *BaseNode) SetChildren(children []Node) {
	bn.mutex.Lock()
	defer bn.mutex.Unlock()
	bn.Children = children
}

// PropagateToChildren propage un fait ou token aux enfants
func (bn *BaseNode) PropagateToChildren(fact *Fact, token *Token) error {
	for _, child := range bn.GetChildren() {
		if fact != nil {
			if err := child.ActivateRight(fact); err != nil {
				return fmt.Errorf("erreur propagation fait vers %s: %w", child.GetID(), err)
			}
		}
		if token != nil {
			if err := child.ActivateLeft(token); err != nil {
				return fmt.Errorf("erreur propagation token vers %s: %w", child.GetID(), err)
			}
		}
	}
	return nil
}

// PropagateRetractToChildren propage la rétractation d'un fait aux nœuds enfants
func (bn *BaseNode) PropagateRetractToChildren(factID string) error {
	for _, child := range bn.GetChildren() {
		if err := child.ActivateRetract(factID); err != nil {
			return fmt.Errorf("erreur propagation rétractation vers %s: %w", child.GetID(), err)
		}
	}
	return nil
}

// SaveMemory sauvegarde la mémoire du nœud
func (bn *BaseNode) SaveMemory() error {
	if bn.Storage != nil {
		return bn.Storage.SaveMemory(bn.ID, bn.Memory)
	}
	return nil
}
