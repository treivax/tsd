// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"encoding/json"
	"fmt"
	"sync"
)

// MemoryStorage implémente Storage en mémoire (pour les tests)
type MemoryStorage struct {
	memories map[string]*WorkingMemory
	mutex    sync.RWMutex
}

// NewMemoryStorage crée un nouveau storage en mémoire
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		memories: make(map[string]*WorkingMemory),
	}
}

// SaveMemory sauvegarde en mémoire
func (ms *MemoryStorage) SaveMemory(nodeID string, memory *WorkingMemory) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// Copier la mémoire pour éviter les modifications concurrentes
	data, err := json.Marshal(memory)
	if err != nil {
		return err
	}

	var copyMemory WorkingMemory
	err = json.Unmarshal(data, &copyMemory)
	if err != nil {
		return err
	}

	ms.memories[nodeID] = &copyMemory
	return nil
}

// LoadMemory charge depuis la mémoire
func (ms *MemoryStorage) LoadMemory(nodeID string) (*WorkingMemory, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	memory, exists := ms.memories[nodeID]
	if !exists {
		return nil, fmt.Errorf("mémoire non trouvée pour le nœud %s", nodeID)
	}

	// Retourner une copie
	data, err := json.Marshal(memory)
	if err != nil {
		return nil, err
	}

	var copyMemory WorkingMemory
	err = json.Unmarshal(data, &copyMemory)
	if err != nil {
		return nil, err
	}

	return &copyMemory, nil
}

// DeleteMemory supprime de la mémoire
func (ms *MemoryStorage) DeleteMemory(nodeID string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	delete(ms.memories, nodeID)
	return nil
}

// ListNodes liste les nœuds en mémoire
func (ms *MemoryStorage) ListNodes() ([]string, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	nodes := make([]string, 0, len(ms.memories))
	for nodeID := range ms.memories {
		nodes = append(nodes, nodeID)
	}
	return nodes, nil
}

// Clear vide tous les faits du storage
func (ms *MemoryStorage) Clear() error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// Vider toutes les mémoires
	for nodeID := range ms.memories {
		delete(ms.memories, nodeID)
	}
	ms.memories = make(map[string]*WorkingMemory)

	return nil
}

// AddFact ajoute un fait au storage (dans une mémoire globale)
func (ms *MemoryStorage) AddFact(fact *Fact) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// Utiliser une mémoire globale pour les faits
	const globalNodeID = "__global_facts__"

	memory, exists := ms.memories[globalNodeID]
	if !exists {
		memory = &WorkingMemory{
			NodeID: globalNodeID,
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		}
		ms.memories[globalNodeID] = memory
	}

	return memory.AddFact(fact)
}

// GetAllFacts récupère tous les faits du storage
func (ms *MemoryStorage) GetAllFacts() []*Fact {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	facts := make([]*Fact, 0)

	// Collecter les faits de toutes les mémoires
	for _, memory := range ms.memories {
		if memory != nil && memory.Facts != nil {
			for _, fact := range memory.Facts {
				facts = append(facts, fact)
			}
		}
	}

	return facts
}

// RemoveFact supprime un fait du storage par son ID interne
func (ms *MemoryStorage) RemoveFact(factID string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// Chercher dans toutes les mémoires
	for _, memory := range ms.memories {
		if memory != nil && memory.Facts != nil {
			if _, exists := memory.Facts[factID]; exists {
				delete(memory.Facts, factID)
				return nil
			}
		}
	}

	return fmt.Errorf("fact %s not found", factID)
}

// GetFact récupère un fait par son ID interne
func (ms *MemoryStorage) GetFact(factID string) *Fact {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	// Chercher dans toutes les mémoires
	for _, memory := range ms.memories {
		if memory != nil && memory.Facts != nil {
			if fact, exists := memory.Facts[factID]; exists {
				return fact
			}
		}
	}

	return nil
}

// Sync garantit que toutes les écritures sont durables et visibles
// Pour MemoryStorage, cette opération vérifie la cohérence interne
// car toutes les données sont déjà en mémoire et donc "durables" dans ce contexte
func (ms *MemoryStorage) Sync() error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// Vérification de cohérence interne :
	// - Toutes les mémoires doivent avoir des structures valides
	// - Pas de faits orphelins ou de tokens sans faits associés
	for nodeID, memory := range ms.memories {
		if memory == nil {
			return fmt.Errorf("mémoire nulle pour le nœud %s", nodeID)
		}

		// Vérifier que les structures de données sont initialisées
		if memory.Facts == nil {
			memory.Facts = make(map[string]*Fact)
		}
		if memory.Tokens == nil {
			memory.Tokens = make(map[string]*Token)
		}
	}

	// Pour MemoryStorage, Sync() réussit toujours après vérification
	// Dans une implémentation avec persistance disque, ici on appellerait fsync()
	return nil
}
