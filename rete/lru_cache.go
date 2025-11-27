// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"container/list"
	"sync"
	"time"
)

// LRUCache est un cache LRU (Least Recently Used) thread-safe
type LRUCache struct {
	capacity int
	ttl      time.Duration
	items    map[string]*lruItem
	order    *list.List
	mutex    sync.RWMutex

	// Statistiques
	hits      int64
	misses    int64
	evictions int64
	sets      int64
}

// lruItem représente un élément dans le cache
type lruItem struct {
	key       string
	value     interface{}
	element   *list.Element
	timestamp time.Time
}

// NewLRUCache crée un nouveau cache LRU
func NewLRUCache(capacity int, ttl time.Duration) *LRUCache {
	if capacity <= 0 {
		capacity = 1000 // Capacité par défaut
	}

	return &LRUCache{
		capacity: capacity,
		ttl:      ttl,
		items:    make(map[string]*lruItem),
		order:    list.New(),
	}
}

// Get récupère une valeur du cache
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, exists := c.items[key]
	if !exists {
		c.misses++
		return nil, false
	}

	// Vérifier l'expiration
	if c.ttl > 0 && time.Since(item.timestamp) > c.ttl {
		c.removeItem(item)
		c.misses++
		return nil, false
	}

	// Déplacer vers le front (le plus récemment utilisé)
	c.order.MoveToFront(item.element)
	c.hits++
	return item.value, true
}

// Set ajoute ou met à jour une valeur dans le cache
func (c *LRUCache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.sets++

	// Si l'élément existe déjà, le mettre à jour
	if item, exists := c.items[key]; exists {
		item.value = value
		item.timestamp = time.Now()
		c.order.MoveToFront(item.element)
		return
	}

	// Éviction si capacité atteinte
	if c.order.Len() >= c.capacity {
		c.evictOldest()
	}

	// Créer le nouvel élément
	item := &lruItem{
		key:       key,
		value:     value,
		timestamp: time.Now(),
	}
	item.element = c.order.PushFront(item)
	c.items[key] = item
}

// Delete supprime une valeur du cache
func (c *LRUCache) Delete(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, exists := c.items[key]
	if !exists {
		return false
	}

	c.removeItem(item)
	return true
}

// Clear vide complètement le cache
func (c *LRUCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items = make(map[string]*lruItem)
	c.order.Init()
}

// Len retourne le nombre d'éléments dans le cache
func (c *LRUCache) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.items)
}

// Capacity retourne la capacité maximale du cache
func (c *LRUCache) Capacity() int {
	return c.capacity
}

// GetStats retourne les statistiques du cache
func (c *LRUCache) GetStats() LRUCacheStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return LRUCacheStats{
		Hits:      c.hits,
		Misses:    c.misses,
		Evictions: c.evictions,
		Sets:      c.sets,
		Size:      len(c.items),
		Capacity:  c.capacity,
	}
}

// GetHitRate retourne le taux de hits (0.0 à 1.0)
func (c *LRUCache) GetHitRate() float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	total := c.hits + c.misses
	if total == 0 {
		return 0.0
	}
	return float64(c.hits) / float64(total)
}

// ResetStats réinitialise les statistiques
func (c *LRUCache) ResetStats() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.hits = 0
	c.misses = 0
	c.evictions = 0
	c.sets = 0
}

// CleanExpired supprime tous les éléments expirés
func (c *LRUCache) CleanExpired() int {
	if c.ttl == 0 {
		return 0 // Pas d'expiration configurée
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	expired := make([]*lruItem, 0)

	// Identifier les éléments expirés
	for _, item := range c.items {
		if now.Sub(item.timestamp) > c.ttl {
			expired = append(expired, item)
		}
	}

	// Supprimer les éléments expirés
	for _, item := range expired {
		c.removeItem(item)
	}

	return len(expired)
}

// Keys retourne toutes les clés du cache (dans l'ordre LRU)
func (c *LRUCache) Keys() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	keys := make([]string, 0, len(c.items))
	for e := c.order.Front(); e != nil; e = e.Next() {
		item := e.Value.(*lruItem)
		keys = append(keys, item.key)
	}
	return keys
}

// Oldest retourne la clé la plus ancienne (prochaine à être évincée)
func (c *LRUCache) Oldest() (string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.order.Len() == 0 {
		return "", false
	}

	item := c.order.Back().Value.(*lruItem)
	return item.key, true
}

// Newest retourne la clé la plus récente
func (c *LRUCache) Newest() (string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.order.Len() == 0 {
		return "", false
	}

	item := c.order.Front().Value.(*lruItem)
	return item.key, true
}

// Contains vérifie si une clé existe dans le cache sans la marquer comme utilisée
func (c *LRUCache) Contains(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return false
	}

	// Vérifier l'expiration
	if c.ttl > 0 && time.Since(item.timestamp) > c.ttl {
		return false
	}

	return true
}

// evictOldest évince l'élément le plus ancien (non thread-safe, doit être appelé avec le lock)
func (c *LRUCache) evictOldest() {
	if c.order.Len() == 0 {
		return
	}

	oldest := c.order.Back()
	if oldest != nil {
		item := oldest.Value.(*lruItem)
		c.removeItem(item)
		c.evictions++
	}
}

// removeItem supprime un élément du cache (non thread-safe, doit être appelé avec le lock)
func (c *LRUCache) removeItem(item *lruItem) {
	c.order.Remove(item.element)
	delete(c.items, item.key)
}

// LRUCacheStats contient les statistiques du cache LRU
type LRUCacheStats struct {
	Hits      int64 `json:"hits"`
	Misses    int64 `json:"misses"`
	Evictions int64 `json:"evictions"`
	Sets      int64 `json:"sets"`
	Size      int   `json:"size"`
	Capacity  int   `json:"capacity"`
}

// HitRate retourne le taux de hits
func (s LRUCacheStats) HitRate() float64 {
	total := s.Hits + s.Misses
	if total == 0 {
		return 0.0
	}
	return float64(s.Hits) / float64(total)
}

// EvictionRate retourne le taux d'évictions par rapport aux sets
func (s LRUCacheStats) EvictionRate() float64 {
	if s.Sets == 0 {
		return 0.0
	}
	return float64(s.Evictions) / float64(s.Sets)
}

// FillRate retourne le taux de remplissage du cache (0.0 à 1.0)
func (s LRUCacheStats) FillRate() float64 {
	if s.Capacity == 0 {
		return 0.0
	}
	return float64(s.Size) / float64(s.Capacity)
}
