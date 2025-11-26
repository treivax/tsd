// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package nodes

import (
	"sync"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// BaseNode implémente les fonctionnalités communes à tous les nœuds.
// Applique le principe DRY (Don't Repeat Yourself).
type BaseNode struct {
	id       string
	nodeType string
	memory   *domain.WorkingMemory
	children []domain.Node
	logger   domain.Logger
	mutex    sync.RWMutex
}

// NewBaseNode crée un nouveau nœud de base.
func NewBaseNode(id, nodeType string, logger domain.Logger) *BaseNode {
	return &BaseNode{
		id:       id,
		nodeType: nodeType,
		memory:   domain.NewWorkingMemory(id),
		children: make([]domain.Node, 0),
		logger:   logger,
	}
}

// ID retourne l'identifiant du nœud.
func (bn *BaseNode) ID() string {
	return bn.id
}

// Type retourne le type du nœud.
func (bn *BaseNode) Type() string {
	return bn.nodeType
}

// GetMemory retourne la mémoire de travail du nœud.
func (bn *BaseNode) GetMemory() *domain.WorkingMemory {
	bn.mutex.RLock()
	defer bn.mutex.RUnlock()
	return bn.memory
}

// AddChild ajoute un nœud enfant.
func (bn *BaseNode) AddChild(child domain.Node) {
	bn.mutex.Lock()
	defer bn.mutex.Unlock()
	bn.children = append(bn.children, child)
}

// GetChildren retourne les nœuds enfants.
func (bn *BaseNode) GetChildren() []domain.Node {
	bn.mutex.RLock()
	defer bn.mutex.RUnlock()

	// Retourne une copie pour éviter les modifications concurrentes
	children := make([]domain.Node, len(bn.children))
	copy(children, bn.children)
	return children
}

// logFactProcessing enregistre le traitement d'un fait.
func (bn *BaseNode) logFactProcessing(fact *domain.Fact, action string) {
	bn.logger.Debug("processing fact", map[string]interface{}{
		"node_id":   bn.id,
		"node_type": bn.nodeType,
		"fact_id":   fact.ID,
		"fact_type": fact.Type,
		"action":    action,
	})
}
