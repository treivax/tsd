// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
)

// NormalizationCache stocke les expressions normalisées pour améliorer les performances
type NormalizationCache struct {
	enabled  bool
	maxSize  int
	cache    map[string]interface{}
	mutex    sync.RWMutex
	hits     atomic.Int64
	misses   atomic.Int64
	eviction string // Stratégie d'éviction: "lru", "fifo", "none"
	lru      *lruTracker
}

// CacheStats contient les statistiques du cache
type CacheStats struct {
	Hits     int64   `json:"hits"`
	Misses   int64   `json:"misses"`
	Size     int     `json:"size"`
	MaxSize  int     `json:"maxSize"`
	HitRate  float64 `json:"hitRate"`
	Enabled  bool    `json:"enabled"`
	Eviction string  `json:"eviction"`
}

// lruTracker garde la trace de l'ordre d'accès pour l'éviction LRU
type lruTracker struct {
	order []string
	index map[string]int
	mutex sync.Mutex
}

// newLRUTracker crée un nouveau tracker LRU
func newLRUTracker() *lruTracker {
	return &lruTracker{
		order: make([]string, 0),
		index: make(map[string]int),
	}
}

// touch marque une clé comme récemment utilisée
func (lru *lruTracker) touch(key string) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	// Si la clé existe déjà, la retirer de sa position actuelle
	if idx, exists := lru.index[key]; exists {
		lru.order = append(lru.order[:idx], lru.order[idx+1:]...)
		// Mettre à jour les index
		for i := idx; i < len(lru.order); i++ {
			lru.index[lru.order[i]] = i
		}
	}

	// Ajouter la clé à la fin (plus récente)
	lru.order = append(lru.order, key)
	lru.index[key] = len(lru.order) - 1
}

// getLeastRecentlyUsed retourne la clé la moins récemment utilisée
func (lru *lruTracker) getLeastRecentlyUsed() string {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if len(lru.order) == 0 {
		return ""
	}
	return lru.order[0]
}

// remove retire une clé du tracker
func (lru *lruTracker) remove(key string) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if idx, exists := lru.index[key]; exists {
		lru.order = append(lru.order[:idx], lru.order[idx+1:]...)
		delete(lru.index, key)
		// Mettre à jour les index
		for i := idx; i < len(lru.order); i++ {
			lru.index[lru.order[i]] = i
		}
	}
}

// clear vide le tracker
func (lru *lruTracker) clear() {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()
	lru.order = make([]string, 0)
	lru.index = make(map[string]int)
}

// globalNormalizationCache est l'instance globale du cache (optionnelle)
var globalNormalizationCache *NormalizationCache

// NewNormalizationCache crée un nouveau cache de normalisation
func NewNormalizationCache(maxSize int) *NormalizationCache {
	return &NormalizationCache{
		enabled:  true,
		maxSize:  maxSize,
		cache:    make(map[string]interface{}),
		eviction: "lru",
		lru:      newLRUTracker(),
	}
}

// NewNormalizationCacheWithEviction crée un cache avec une stratégie d'éviction spécifique
func NewNormalizationCacheWithEviction(maxSize int, eviction string) *NormalizationCache {
	cache := &NormalizationCache{
		enabled:  true,
		maxSize:  maxSize,
		cache:    make(map[string]interface{}),
		eviction: eviction,
	}
	if eviction == "lru" {
		cache.lru = newLRUTracker()
	}
	return cache
}

// SetGlobalCache définit le cache global de normalisation
func SetGlobalCache(cache *NormalizationCache) {
	globalNormalizationCache = cache
}

// GetGlobalCache retourne le cache global de normalisation
func GetGlobalCache() *NormalizationCache {
	return globalNormalizationCache
}

// Enable active le cache
func (nc *NormalizationCache) Enable() {
	nc.mutex.Lock()
	defer nc.mutex.Unlock()
	nc.enabled = true
}

// Disable désactive le cache
func (nc *NormalizationCache) Disable() {
	nc.mutex.Lock()
	defer nc.mutex.Unlock()
	nc.enabled = false
}

// IsEnabled retourne true si le cache est activé
func (nc *NormalizationCache) IsEnabled() bool {
	nc.mutex.RLock()
	defer nc.mutex.RUnlock()
	return nc.enabled
}

// Get récupère une expression normalisée du cache
func (nc *NormalizationCache) Get(key string) (interface{}, bool) {
	if !nc.IsEnabled() {
		return nil, false
	}

	nc.mutex.RLock()
	value, exists := nc.cache[key]
	nc.mutex.RUnlock()

	if exists {
		nc.hits.Add(1)
		if nc.eviction == "lru" && nc.lru != nil {
			nc.lru.touch(key)
		}
		return value, true
	}

	nc.misses.Add(1)
	return nil, false
}

// Set ajoute une expression normalisée au cache
func (nc *NormalizationCache) Set(key string, value interface{}) {
	if !nc.IsEnabled() {
		return
	}

	nc.mutex.Lock()
	defer nc.mutex.Unlock()

	// Si le cache est plein, appliquer l'éviction
	if len(nc.cache) >= nc.maxSize && nc.maxSize > 0 {
		nc.evict()
	}

	nc.cache[key] = value

	if nc.eviction == "lru" && nc.lru != nil {
		nc.lru.touch(key)
	}
}

// evict retire une entrée du cache selon la stratégie d'éviction
// DOIT être appelé avec le mutex verrouillé
func (nc *NormalizationCache) evict() {
	if len(nc.cache) == 0 {
		return
	}

	switch nc.eviction {
	case "lru":
		if nc.lru != nil {
			keyToEvict := nc.lru.getLeastRecentlyUsed()
			if keyToEvict != "" {
				delete(nc.cache, keyToEvict)
				nc.lru.remove(keyToEvict)
			}
		}
	case "fifo":
		// Pour FIFO, on retire la première clé (arbitraire dans une map Go)
		for key := range nc.cache {
			delete(nc.cache, key)
			break
		}
	default:
		// Pas d'éviction, on refuse d'ajouter
		return
	}
}

// Clear vide complètement le cache
func (nc *NormalizationCache) Clear() {
	nc.mutex.Lock()
	defer nc.mutex.Unlock()

	nc.cache = make(map[string]interface{})
	if nc.lru != nil {
		nc.lru.clear()
	}
}

// ResetStats réinitialise les statistiques du cache
func (nc *NormalizationCache) ResetStats() {
	nc.hits.Store(0)
	nc.misses.Store(0)
}

// GetStats retourne les statistiques du cache
func (nc *NormalizationCache) GetStats() CacheStats {
	nc.mutex.RLock()
	size := len(nc.cache)
	enabled := nc.enabled
	eviction := nc.eviction
	nc.mutex.RUnlock()

	hits := nc.hits.Load()
	misses := nc.misses.Load()
	total := hits + misses

	hitRate := 0.0
	if total > 0 {
		hitRate = float64(hits) / float64(total)
	}

	return CacheStats{
		Hits:     hits,
		Misses:   misses,
		Size:     size,
		MaxSize:  nc.maxSize,
		HitRate:  hitRate,
		Enabled:  enabled,
		Eviction: eviction,
	}
}

// Size retourne le nombre d'entrées dans le cache
func (nc *NormalizationCache) Size() int {
	nc.mutex.RLock()
	defer nc.mutex.RUnlock()
	return len(nc.cache)
}

// computeCacheKey calcule une clé de cache unique pour une expression
func computeCacheKey(expr interface{}) string {
	// Sérialiser l'expression en JSON pour avoir une représentation unique
	jsonBytes, err := json.Marshal(expr)
	if err != nil {
		// En cas d'erreur, utiliser une représentation string
		return fmt.Sprintf("%T:%v", expr, expr)
	}

	// Calculer le hash SHA-256
	hash := sha256.Sum256(jsonBytes)
	return fmt.Sprintf("%x", hash)
}

// NormalizeExpressionWithCache normalise une expression en utilisant le cache
func NormalizeExpressionWithCache(expr interface{}, cache *NormalizationCache) (interface{}, error) {
	// Si pas de cache ou cache désactivé, utiliser la normalisation directe
	if cache == nil || !cache.IsEnabled() {
		return NormalizeExpression(expr)
	}

	// Calculer la clé du cache
	key := computeCacheKey(expr)

	// Chercher dans le cache
	if cached, found := cache.Get(key); found {
		return cached, nil
	}

	// Pas dans le cache, normaliser
	normalized, err := NormalizeExpression(expr)
	if err != nil {
		return nil, err
	}

	// Stocker dans le cache
	cache.Set(key, normalized)

	return normalized, nil
}

// NormalizeExpressionCached normalise une expression en utilisant le cache global
func NormalizeExpressionCached(expr interface{}) (interface{}, error) {
	return NormalizeExpressionWithCache(expr, globalNormalizationCache)
}

// SetCacheMaxSize change la taille maximum du cache
func (nc *NormalizationCache) SetCacheMaxSize(maxSize int) {
	nc.mutex.Lock()
	defer nc.mutex.Unlock()
	nc.maxSize = maxSize

	// Si le cache est maintenant trop grand, évincer des entrées
	for len(nc.cache) > maxSize && maxSize > 0 {
		nc.evict()
	}
}

// SetEvictionStrategy change la stratégie d'éviction
func (nc *NormalizationCache) SetEvictionStrategy(strategy string) {
	nc.mutex.Lock()
	defer nc.mutex.Unlock()

	nc.eviction = strategy

	// Initialiser le tracker LRU si nécessaire
	if strategy == "lru" && nc.lru == nil {
		nc.lru = newLRUTracker()
		// Peupler le tracker avec les clés existantes
		for key := range nc.cache {
			nc.lru.touch(key)
		}
	}
}

// GetHitRate retourne le taux de succès du cache (0.0 à 1.0)
func (nc *NormalizationCache) GetHitRate() float64 {
	hits := nc.hits.Load()
	misses := nc.misses.Load()
	total := hits + misses

	if total == 0 {
		return 0.0
	}

	return float64(hits) / float64(total)
}

// String retourne une représentation string des statistiques du cache
func (cs CacheStats) String() string {
	return fmt.Sprintf(
		"CacheStats{Hits: %d, Misses: %d, Size: %d/%d, HitRate: %.2f%%, Enabled: %v, Eviction: %s}",
		cs.Hits,
		cs.Misses,
		cs.Size,
		cs.MaxSize,
		cs.HitRate*100,
		cs.Enabled,
		cs.Eviction,
	)
}
