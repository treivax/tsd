package rete

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdStorage implémente l'interface Storage avec etcd
type EtcdStorage struct {
	client    *clientv3.Client
	keyPrefix string
	timeout   time.Duration
}

// NewEtcdStorage crée un nouveau storage etcd
func NewEtcdStorage(endpoints []string, keyPrefix string) (*EtcdStorage, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("erreur connexion etcd: %w", err)
	}

	storage := &EtcdStorage{
		client:    client,
		keyPrefix: keyPrefix,
		timeout:   5 * time.Second,
	}

	// Tester la connexion
	ctx, cancel := context.WithTimeout(context.Background(), storage.timeout)
	defer cancel()

	_, err = client.Status(ctx, endpoints[0])
	if err != nil {
		return nil, fmt.Errorf("test connexion etcd échoué: %w", err)
	}

	fmt.Printf("✅ Connexion etcd établie (prefix: %s)\n", keyPrefix)
	return storage, nil
}

// Close ferme la connexion etcd
func (es *EtcdStorage) Close() error {
	return es.client.Close()
}

// SaveMemory sauvegarde la mémoire d'un nœud dans etcd
func (es *EtcdStorage) SaveMemory(nodeID string, memory *WorkingMemory) error {
	ctx, cancel := context.WithTimeout(context.Background(), es.timeout)
	defer cancel()

	// Sérialiser la mémoire en JSON
	data, err := json.Marshal(memory)
	if err != nil {
		return fmt.Errorf("erreur sérialisation mémoire: %w", err)
	}

	// Clé etcd pour ce nœud
	key := fmt.Sprintf("%s/nodes/%s/memory", es.keyPrefix, nodeID)

	// Sauvegarder dans etcd avec un timestamp
	_, err = es.client.Put(ctx, key, string(data))
	if err != nil {
		return fmt.Errorf("erreur sauvegarde etcd: %w", err)
	}

	// Également sauvegarder un timestamp de dernière modification
	timestampKey := fmt.Sprintf("%s/nodes/%s/timestamp", es.keyPrefix, nodeID)
	_, err = es.client.Put(ctx, timestampKey, time.Now().Format(time.RFC3339))
	if err != nil {
		fmt.Printf("⚠️  Erreur sauvegarde timestamp: %v\n", err) // Non bloquant
	}

	return nil
}

// LoadMemory charge la mémoire d'un nœud depuis etcd
func (es *EtcdStorage) LoadMemory(nodeID string) (*WorkingMemory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), es.timeout)
	defer cancel()

	// Clé etcd pour ce nœud
	key := fmt.Sprintf("%s/nodes/%s/memory", es.keyPrefix, nodeID)

	// Récupérer depuis etcd
	resp, err := es.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture etcd: %w", err)
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("mémoire non trouvée pour le nœud %s", nodeID)
	}

	// Désérialiser la mémoire
	var memory WorkingMemory
	err = json.Unmarshal(resp.Kvs[0].Value, &memory)
	if err != nil {
		return nil, fmt.Errorf("erreur désérialisation mémoire: %w", err)
	}

	return &memory, nil
}

// DeleteMemory supprime la mémoire d'un nœud d'etcd
func (es *EtcdStorage) DeleteMemory(nodeID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), es.timeout)
	defer cancel()

	// Supprimer la mémoire et le timestamp
	memoryKey := fmt.Sprintf("%s/nodes/%s/memory", es.keyPrefix, nodeID)
	timestampKey := fmt.Sprintf("%s/nodes/%s/timestamp", es.keyPrefix, nodeID)

	_, err := es.client.Delete(ctx, memoryKey)
	if err != nil {
		return fmt.Errorf("erreur suppression mémoire etcd: %w", err)
	}

	_, err = es.client.Delete(ctx, timestampKey)
	if err != nil {
		fmt.Printf("⚠️  Erreur suppression timestamp: %v\n", err) // Non bloquant
	}

	return nil
}

// ListNodes retourne la liste des nœuds stockés
func (es *EtcdStorage) ListNodes() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), es.timeout)
	defer cancel()

	// Prefix pour tous les nœuds
	prefix := fmt.Sprintf("%s/nodes/", es.keyPrefix)

	// Récupérer toutes les clés avec ce préfixe
	resp, err := es.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return nil, fmt.Errorf("erreur liste nœuds etcd: %w", err)
	}

	// Extraire les IDs des nœuds
	nodeIDs := make(map[string]bool)
	for _, kv := range resp.Kvs {
		key := string(kv.Key)
		// Extraire l'ID du nœud depuis la clé
		// Format: prefix/nodes/{nodeID}/memory ou prefix/nodes/{nodeID}/timestamp
		parts := strings.Split(strings.TrimPrefix(key, prefix), "/")
		if len(parts) >= 1 {
			nodeIDs[parts[0]] = true
		}
	}

	// Convertir en slice
	nodes := make([]string, 0, len(nodeIDs))
	for nodeID := range nodeIDs {
		nodes = append(nodes, nodeID)
	}

	return nodes, nil
}

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
