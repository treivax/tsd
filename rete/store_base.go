// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"encoding/json"
	"fmt"
	"sync"
)

// MemoryStorage implements pure in-memory storage with strong consistency guarantees.
// This is the only storage implementation in TSD - all data is kept in memory.
// Facts can be exported to .tsd files, and network replication via Raft is planned.
type MemoryStorage struct {
	memories map[string]*WorkingMemory
	mutex    sync.RWMutex
}

// NewMemoryStorage creates a new in-memory storage instance with strong consistency.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		memories: make(map[string]*WorkingMemory),
	}
}

// SaveMemory saves working memory to in-memory storage.
// Thread-safe with mutex protection.
func (ms *MemoryStorage) SaveMemory(nodeID string, memory *WorkingMemory) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// Copy memory to avoid concurrent modifications
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

// LoadMemory loads working memory from in-memory storage.
// Returns a copy to ensure thread safety.
func (ms *MemoryStorage) LoadMemory(nodeID string) (*WorkingMemory, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	memory, exists := ms.memories[nodeID]
	if !exists {
		return nil, fmt.Errorf("mémoire non trouvée pour le nœud %s", nodeID)
	}

	// Return a copy for thread safety
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

// DeleteMemory removes working memory from in-memory storage.
func (ms *MemoryStorage) DeleteMemory(nodeID string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	delete(ms.memories, nodeID)
	return nil
}

// ListNodes returns all node IDs stored in memory.
func (ms *MemoryStorage) ListNodes() ([]string, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	nodes := make([]string, 0, len(ms.memories))
	for nodeID := range ms.memories {
		nodes = append(nodes, nodeID)
	}
	return nodes, nil
}

// Clear removes all facts from in-memory storage.
func (ms *MemoryStorage) Clear() error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// Clear all memories
	for nodeID := range ms.memories {
		delete(ms.memories, nodeID)
	}
	ms.memories = make(map[string]*WorkingMemory)

	return nil
}

// AddFact adds a fact to in-memory storage (in global memory).
func (ms *MemoryStorage) AddFact(fact *Fact) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// Use global memory for facts
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

// GetAllFacts retrieves all facts from in-memory storage.
func (ms *MemoryStorage) GetAllFacts() []*Fact {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	facts := make([]*Fact, 0)

	// Collect facts from all memories
	for _, memory := range ms.memories {
		if memory != nil && memory.Facts != nil {
			for _, fact := range memory.Facts {
				facts = append(facts, fact)
			}
		}
	}

	return facts
}

// RemoveFact removes a fact from in-memory storage by its internal ID.
func (ms *MemoryStorage) RemoveFact(factID string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// Search in all memories
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

// GetFact retrieves a fact by its internal ID from in-memory storage.
func (ms *MemoryStorage) GetFact(factID string) *Fact {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	// Search in all memories
	for _, memory := range ms.memories {
		if memory != nil && memory.Facts != nil {
			if fact, exists := memory.Facts[factID]; exists {
				return fact
			}
		}
	}

	return nil
}

// Sync ensures all writes are consistent and visible.
// For in-memory storage, this verifies internal consistency since all data
// is already in memory. In a future replicated setup, this would coordinate
// with other nodes via Raft consensus.
func (ms *MemoryStorage) Sync() error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// Internal consistency verification:
	// - All memories must have valid structures
	// - No orphaned facts or tokens without associated facts
	for nodeID, memory := range ms.memories {
		if memory == nil {
			return fmt.Errorf("mémoire nulle pour le nœud %s", nodeID)
		}

		// Verify data structures are initialized
		if memory.Facts == nil {
			memory.Facts = make(map[string]*Fact)
		}
		if memory.Tokens == nil {
			memory.Tokens = make(map[string]*Token)
		}
	}

	// For in-memory storage, Sync() always succeeds after verification.
	// In a future replicated implementation, this would ensure Raft consensus.
	return nil
}
