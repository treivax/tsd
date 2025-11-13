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
