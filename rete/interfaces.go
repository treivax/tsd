// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

type Node interface {
	GetID() string
	GetType() string
	GetMemory() *WorkingMemory
	ActivateLeft(token *Token) error
	ActivateRight(fact *Fact) error
	ActivateRetract(factID string) error
	AddChild(child Node)
	GetChildren() []Node
}

// Storage interface pour la persistance
type Storage interface {
	SaveMemory(nodeID string, memory *WorkingMemory) error
	LoadMemory(nodeID string) (*WorkingMemory, error)
	DeleteMemory(nodeID string) error
	ListNodes() ([]string, error)
	Clear() error                   // Vider tous les faits du storage
	AddFact(fact *Fact) error       // Ajouter un fait au storage
	RemoveFact(factID string) error // Supprimer un fait du storage par son ID interne
	GetFact(factID string) *Fact    // Récupérer un fait par son ID interne
	GetAllFacts() []*Fact           // Récupérer tous les faits du storage
}
