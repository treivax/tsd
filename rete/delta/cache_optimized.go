// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"sync"
	"sync/atomic"
	"time"
)

// OptimizedCache est un cache haute performance avec éviction LRU.
//
// Optimisations :
// - LRU list avec pointeurs pour O(1) operations
// - Éviction automatique basée sur taille et TTL
// - Lock granulaire pour minimiser contention
// - Métriques atomiques pour zero-overhead
type OptimizedCache struct {
	entries    map[string]*cacheEntry
	lruList    *lruList
	maxSize    int
	ttl        time.Duration
	hits       atomic.Int64
	misses     atomic.Int64
	evictions  atomic.Int64
	mutex      sync.RWMutex
	cleanupMux sync.Mutex
}

// lruList est une liste doublement chaînée pour LRU.
type lruList struct {
	head *lruNode
	tail *lruNode
	size int
}

// lruNode représente un nœud dans la liste LRU.
type lruNode struct {
	key   string
	entry *cacheEntry
	prev  *lruNode
	next  *lruNode
}

// cacheEntry représente une entrée de cache.
type cacheEntry struct {
	delta       *FactDelta
	createdAt   time.Time
	accessedAt  time.Time
	accessCount atomic.Int64
	node        *lruNode // Référence au nœud LRU pour O(1) access
}

// NewOptimizedCache crée un cache optimisé.
func NewOptimizedCache(maxSize int, ttl time.Duration) *OptimizedCache {
	if maxSize <= 0 {
		maxSize = defaultCacheSize
	}
	if ttl <= 0 {
		ttl = DefaultCacheTTL
	}

	return &OptimizedCache{
		entries: make(map[string]*cacheEntry, maxSize),
		lruList: &lruList{},
		maxSize: maxSize,
		ttl:     ttl,
	}
}

// Get récupère une valeur depuis le cache.
func (oc *OptimizedCache) Get(key string) (*FactDelta, bool) {
	oc.mutex.RLock()
	entry, exists := oc.entries[key]
	oc.mutex.RUnlock()

	if !exists {
		oc.misses.Add(1)
		return nil, false
	}

	// Vérifier expiration sans lock
	if time.Since(entry.createdAt) > oc.ttl {
		// Expiration détectée - cleanup asynchrone
		go oc.removeExpired(key)
		oc.misses.Add(1)
		return nil, false
	}

	// Mettre à jour stats
	entry.accessedAt = time.Now()
	entry.accessCount.Add(1)

	// Mettre à jour LRU (avec write lock)
	oc.mutex.Lock()
	oc.lruList.moveToFront(entry.node)
	oc.mutex.Unlock()

	oc.hits.Add(1)
	return entry.delta, true
}

// Put ajoute une valeur au cache.
func (oc *OptimizedCache) Put(key string, delta *FactDelta) {
	oc.mutex.Lock()
	defer oc.mutex.Unlock()

	// Si existe déjà, mettre à jour
	if entry, exists := oc.entries[key]; exists {
		entry.delta = delta
		entry.accessedAt = time.Now()
		oc.lruList.moveToFront(entry.node)
		return
	}

	// Si plein, évincer LRU
	if len(oc.entries) >= oc.maxSize {
		oc.evictLRU()
	}

	// Créer nouvelle entrée
	now := time.Now()
	entry := &cacheEntry{
		delta:      delta,
		createdAt:  now,
		accessedAt: now,
	}

	// Ajouter à la liste LRU
	node := oc.lruList.addToFront(key, entry)
	entry.node = node

	// Ajouter à la map
	oc.entries[key] = entry
}

// evictLRU évince l'entrée la moins récemment utilisée.
func (oc *OptimizedCache) evictLRU() {
	if oc.lruList.tail == nil {
		return
	}

	key := oc.lruList.tail.key
	delete(oc.entries, key)
	oc.lruList.removeTail()
	oc.evictions.Add(1)
}

// removeExpired supprime une entrée expirée.
func (oc *OptimizedCache) removeExpired(key string) {
	oc.cleanupMux.Lock()
	defer oc.cleanupMux.Unlock()

	oc.mutex.Lock()
	defer oc.mutex.Unlock()

	entry, exists := oc.entries[key]
	if !exists {
		return
	}

	// Double-check expiration
	if time.Since(entry.createdAt) > oc.ttl {
		delete(oc.entries, key)
		oc.lruList.removeNode(entry.node)
		oc.evictions.Add(1)
	}
}

// Clear vide complètement le cache.
func (oc *OptimizedCache) Clear() {
	oc.mutex.Lock()
	defer oc.mutex.Unlock()

	oc.entries = make(map[string]*cacheEntry, oc.maxSize)
	oc.lruList = &lruList{}
}

// GetStats retourne les statistiques du cache.
func (oc *OptimizedCache) GetStats() CacheStats {
	oc.mutex.RLock()
	size := len(oc.entries)
	oc.mutex.RUnlock()

	hits := oc.hits.Load()
	misses := oc.misses.Load()
	evictions := oc.evictions.Load()

	total := hits + misses
	hitRate := 0.0
	if total > 0 {
		hitRate = float64(hits) / float64(total)
	}

	return CacheStats{
		Size:      size,
		Hits:      hits,
		Misses:    misses,
		Evictions: evictions,
		HitRate:   hitRate,
	}
}

// CacheStats contient les statistiques du cache.
type CacheStats struct {
	Size      int
	Hits      int64
	Misses    int64
	Evictions int64
	HitRate   float64
}

// --- Implémentation lruList ---

// addToFront ajoute un nœud en tête de liste.
func (l *lruList) addToFront(key string, entry *cacheEntry) *lruNode {
	node := &lruNode{
		key:   key,
		entry: entry,
	}

	if l.head == nil {
		l.head = node
		l.tail = node
	} else {
		node.next = l.head
		l.head.prev = node
		l.head = node
	}

	l.size++
	return node
}

// moveToFront déplace un nœud en tête de liste.
func (l *lruList) moveToFront(node *lruNode) {
	if node == nil || node == l.head {
		return
	}

	// Détacher du milieu/fin
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if node == l.tail {
		l.tail = node.prev
	}

	// Mettre en tête
	node.prev = nil
	node.next = l.head
	if l.head != nil {
		l.head.prev = node
	}
	l.head = node
}

// removeTail supprime le nœud de queue.
func (l *lruList) removeTail() {
	if l.tail == nil {
		return
	}

	if l.tail.prev != nil {
		l.tail.prev.next = nil
		l.tail = l.tail.prev
	} else {
		// Un seul élément
		l.head = nil
		l.tail = nil
	}

	l.size--
}

// removeNode supprime un nœud spécifique.
func (l *lruList) removeNode(node *lruNode) {
	if node == nil {
		return
	}

	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		l.tail = node.prev
	}

	l.size--
}

// Size constants
const (
	defaultCacheSize = 1000
)
