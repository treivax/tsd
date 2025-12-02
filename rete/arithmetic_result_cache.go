// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// ArithmeticResultCache fournit un cache LRU thread-safe pour les résultats
// arithmétiques intermédiaires avec support TTL et statistiques détaillées
type ArithmeticResultCache struct {
	entries    map[string]*cacheEntry
	lruList    *lruNode
	maxSize    int
	ttl        time.Duration
	mutex      sync.RWMutex
	stats      CacheStatistics
	enabled    bool
	onEviction EvictionCallback
}

// cacheEntry représente une entrée dans le cache
type cacheEntry struct {
	key        string
	value      interface{}
	computedAt time.Time
	lastAccess time.Time
	hitCount   int64
	lruNode    *lruNode
}

// lruNode représente un nœud dans la liste doublement chaînée LRU
type lruNode struct {
	prev  *lruNode
	next  *lruNode
	entry *cacheEntry
}

// CacheStatistics collecte les statistiques du cache
type CacheStatistics struct {
	Hits              int64
	Misses            int64
	Evictions         int64
	Sets              int64
	CurrentSize       int
	TotalComputations int64
	AverageHitTime    time.Duration
	AverageMissTime   time.Duration
	totalHitTime      time.Duration
	totalMissTime     time.Duration
}

// EvictionCallback est appelé lors de l'éviction d'une entrée
type EvictionCallback func(key string, value interface{})

// CacheConfig configure le comportement du cache
type CacheConfig struct {
	MaxSize    int
	TTL        time.Duration
	Enabled    bool
	OnEviction EvictionCallback
}

// DefaultCacheConfig retourne une configuration par défaut
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		MaxSize: 1000,
		TTL:     5 * time.Minute,
		Enabled: true,
	}
}

// NewArithmeticResultCache crée un nouveau cache avec la configuration donnée
func NewArithmeticResultCache(config CacheConfig) *ArithmeticResultCache {
	// Créer le sentinel LRU (nœud bidirectionnel sentinelle)
	sentinel := &lruNode{}
	sentinel.prev = sentinel
	sentinel.next = sentinel

	return &ArithmeticResultCache{
		entries:    make(map[string]*cacheEntry),
		lruList:    sentinel,
		maxSize:    config.MaxSize,
		ttl:        config.TTL,
		enabled:    config.Enabled,
		onEviction: config.OnEviction,
		stats: CacheStatistics{
			CurrentSize: 0,
		},
	}
}

// GenerateCacheKey génère une clé de cache à partir du nom de résultat
// et des valeurs de dépendances
func GenerateCacheKey(resultName string, dependencies map[string]interface{}) string {
	// Sérialiser les dépendances de manière déterministe
	data := map[string]interface{}{
		"result": resultName,
		"deps":   dependencies,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		// Fallback: utiliser uniquement le nom de résultat
		return fmt.Sprintf("%s_error", resultName)
	}

	// Générer hash SHA-256
	hash := sha256.Sum256(jsonBytes)
	return hex.EncodeToString(hash[:])
}

// Get récupère une valeur du cache
func (arc *ArithmeticResultCache) Get(key string) (interface{}, bool) {
	if !arc.enabled {
		return nil, false
	}

	startTime := time.Now()
	arc.mutex.Lock()
	defer arc.mutex.Unlock()

	entry, exists := arc.entries[key]
	if !exists {
		arc.stats.Misses++
		arc.stats.totalMissTime += time.Since(startTime)
		if arc.stats.Misses > 0 {
			arc.stats.AverageMissTime = arc.stats.totalMissTime / time.Duration(arc.stats.Misses)
		}
		return nil, false
	}

	// Vérifier TTL
	if arc.ttl > 0 && time.Since(entry.computedAt) > arc.ttl {
		// Entrée expirée, la retirer
		arc.removeEntry(entry)
		arc.stats.Misses++
		arc.stats.totalMissTime += time.Since(startTime)
		if arc.stats.Misses > 0 {
			arc.stats.AverageMissTime = arc.stats.totalMissTime / time.Duration(arc.stats.Misses)
		}
		return nil, false
	}

	// Mise à jour des statistiques d'accès
	entry.lastAccess = time.Now()
	entry.hitCount++
	arc.stats.Hits++
	arc.stats.totalHitTime += time.Since(startTime)
	if arc.stats.Hits > 0 {
		arc.stats.AverageHitTime = arc.stats.totalHitTime / time.Duration(arc.stats.Hits)
	}

	// Déplacer en tête de LRU (le plus récemment utilisé)
	arc.moveToFront(entry)

	return entry.value, true
}

// GetWithDependencies récupère une valeur en générant automatiquement la clé
func (arc *ArithmeticResultCache) GetWithDependencies(
	resultName string,
	dependencies map[string]interface{},
) (interface{}, bool) {
	key := GenerateCacheKey(resultName, dependencies)
	return arc.Get(key)
}

// Set ajoute une valeur au cache
func (arc *ArithmeticResultCache) Set(key string, value interface{}) {
	if !arc.enabled {
		return
	}

	arc.mutex.Lock()
	defer arc.mutex.Unlock()

	arc.stats.Sets++
	arc.stats.TotalComputations++

	// Si l'entrée existe déjà, la mettre à jour
	if entry, exists := arc.entries[key]; exists {
		entry.value = value
		entry.computedAt = time.Now()
		entry.lastAccess = time.Now()
		arc.moveToFront(entry)
		return
	}

	// Éviction si nécessaire
	if len(arc.entries) >= arc.maxSize {
		arc.evictLRU()
	}

	// Créer nouvelle entrée
	entry := &cacheEntry{
		key:        key,
		value:      value,
		computedAt: time.Now(),
		lastAccess: time.Now(),
		hitCount:   0,
	}

	// Créer nœud LRU
	node := &lruNode{entry: entry}
	entry.lruNode = node

	// Ajouter à la map
	arc.entries[key] = entry

	// Ajouter en tête de la liste LRU
	arc.addToFront(node)

	arc.stats.CurrentSize = len(arc.entries)
}

// SetWithDependencies ajoute une valeur en générant automatiquement la clé
func (arc *ArithmeticResultCache) SetWithDependencies(
	resultName string,
	dependencies map[string]interface{},
	value interface{},
) {
	key := GenerateCacheKey(resultName, dependencies)
	arc.Set(key, value)
}

// Clear vide complètement le cache
func (arc *ArithmeticResultCache) Clear() {
	arc.mutex.Lock()
	defer arc.mutex.Unlock()

	arc.entries = make(map[string]*cacheEntry)

	// Réinitialiser la liste LRU
	arc.lruList.prev = arc.lruList
	arc.lruList.next = arc.lruList

	arc.stats.CurrentSize = 0
}

// GetStatistics retourne une copie thread-safe des statistiques
func (arc *ArithmeticResultCache) GetStatistics() CacheStatistics {
	arc.mutex.RLock()
	defer arc.mutex.RUnlock()

	stats := arc.stats
	stats.CurrentSize = len(arc.entries)
	return stats
}

// GetHitRate retourne le taux de succès du cache (0.0 à 1.0)
func (arc *ArithmeticResultCache) GetHitRate() float64 {
	arc.mutex.RLock()
	defer arc.mutex.RUnlock()

	total := arc.stats.Hits + arc.stats.Misses
	if total == 0 {
		return 0.0
	}
	return float64(arc.stats.Hits) / float64(total)
}

// GetSize retourne le nombre d'entrées actuellement dans le cache
func (arc *ArithmeticResultCache) GetSize() int {
	arc.mutex.RLock()
	defer arc.mutex.RUnlock()

	return len(arc.entries)
}

// SetEnabled active ou désactive le cache dynamiquement
func (arc *ArithmeticResultCache) SetEnabled(enabled bool) {
	arc.mutex.Lock()
	defer arc.mutex.Unlock()

	arc.enabled = enabled
}

// IsEnabled retourne l'état d'activation du cache
func (arc *ArithmeticResultCache) IsEnabled() bool {
	arc.mutex.RLock()
	defer arc.mutex.RUnlock()

	return arc.enabled
}

// ResetStatistics réinitialise les statistiques sans vider le cache
func (arc *ArithmeticResultCache) ResetStatistics() {
	arc.mutex.Lock()
	defer arc.mutex.Unlock()

	arc.stats = CacheStatistics{
		CurrentSize: len(arc.entries),
	}
}

// GetTopEntries retourne les N entrées les plus utilisées
func (arc *ArithmeticResultCache) GetTopEntries(n int) []CacheEntryInfo {
	arc.mutex.RLock()
	defer arc.mutex.RUnlock()

	// Collecter toutes les entrées
	entries := make([]*cacheEntry, 0, len(arc.entries))
	for _, entry := range arc.entries {
		entries = append(entries, entry)
	}

	// Tri par bulle simple (suffisant pour petites listes)
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].hitCount > entries[i].hitCount {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	// Limiter à n
	if n > len(entries) {
		n = len(entries)
	}

	// Convertir en CacheEntryInfo
	result := make([]CacheEntryInfo, n)
	for i := 0; i < n; i++ {
		result[i] = CacheEntryInfo{
			Key:        entries[i].key,
			HitCount:   entries[i].hitCount,
			ComputedAt: entries[i].computedAt,
			LastAccess: entries[i].lastAccess,
			Age:        time.Since(entries[i].computedAt),
		}
	}

	return result
}

// CacheEntryInfo contient des informations publiques sur une entrée
type CacheEntryInfo struct {
	Key        string
	HitCount   int64
	ComputedAt time.Time
	LastAccess time.Time
	Age        time.Duration
}

// GetSummary retourne un résumé formaté des statistiques
func (arc *ArithmeticResultCache) GetSummary() map[string]interface{} {
	stats := arc.GetStatistics()

	return map[string]interface{}{
		"enabled":       arc.enabled,
		"size":          stats.CurrentSize,
		"max_size":      arc.maxSize,
		"ttl":           arc.ttl.String(),
		"hits":          stats.Hits,
		"misses":        stats.Misses,
		"hit_rate":      arc.GetHitRate(),
		"evictions":     stats.Evictions,
		"sets":          stats.Sets,
		"avg_hit_time":  stats.AverageHitTime.String(),
		"avg_miss_time": stats.AverageMissTime.String(),
	}
}

// --- Méthodes privées pour gestion LRU ---

// moveToFront déplace une entrée en tête de la liste LRU
func (arc *ArithmeticResultCache) moveToFront(entry *cacheEntry) {
	node := entry.lruNode
	arc.removeNode(node)
	arc.addToFront(node)
}

// addToFront ajoute un nœud en tête de la liste LRU
func (arc *ArithmeticResultCache) addToFront(node *lruNode) {
	node.next = arc.lruList.next
	node.prev = arc.lruList
	arc.lruList.next.prev = node
	arc.lruList.next = node
}

// removeNode retire un nœud de la liste LRU
func (arc *ArithmeticResultCache) removeNode(node *lruNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// evictLRU évince l'entrée la moins récemment utilisée
func (arc *ArithmeticResultCache) evictLRU() {
	if arc.lruList.prev == arc.lruList {
		// Liste vide
		return
	}

	// Le dernier nœud est le moins récemment utilisé
	lru := arc.lruList.prev
	arc.removeEntry(lru.entry)
}

// removeEntry retire une entrée du cache
func (arc *ArithmeticResultCache) removeEntry(entry *cacheEntry) {
	// Retirer de la map
	delete(arc.entries, entry.key)

	// Retirer de la liste LRU
	arc.removeNode(entry.lruNode)

	// Mettre à jour statistiques
	arc.stats.Evictions++
	arc.stats.CurrentSize = len(arc.entries)

	// Callback d'éviction
	if arc.onEviction != nil {
		arc.onEviction(entry.key, entry.value)
	}
}

// EstimateMemoryUsage estime l'utilisation mémoire du cache en octets
func (arc *ArithmeticResultCache) EstimateMemoryUsage() int64 {
	arc.mutex.RLock()
	defer arc.mutex.RUnlock()

	// Estimation approximative :
	// - Map overhead : ~48 bytes par entrée
	// - cacheEntry struct : ~120 bytes
	// - lruNode : ~24 bytes
	// - Key (string) : variable, moyenne ~32 bytes
	// - Value : variable, moyenne ~64 bytes (dépend du type)

	const avgEntrySize = 48 + 120 + 24 + 32 + 64 // ~288 bytes
	return int64(len(arc.entries)) * avgEntrySize
}

// Purge retire toutes les entrées expirées (cleanup)
func (arc *ArithmeticResultCache) Purge() int {
	if arc.ttl == 0 {
		return 0 // Pas de TTL configuré
	}

	arc.mutex.Lock()
	defer arc.mutex.Unlock()

	now := time.Now()
	expiredCount := 0

	// Parcourir toutes les entrées et retirer les expirées
	for key, entry := range arc.entries {
		if now.Sub(entry.computedAt) > arc.ttl {
			delete(arc.entries, key)
			arc.removeNode(entry.lruNode)
			expiredCount++

			if arc.onEviction != nil {
				arc.onEviction(entry.key, entry.value)
			}
		}
	}

	if expiredCount > 0 {
		arc.stats.Evictions += int64(expiredCount)
		arc.stats.CurrentSize = len(arc.entries)
	}

	return expiredCount
}

// StartAutoPurge démarre une goroutine de nettoyage périodique
func (arc *ArithmeticResultCache) StartAutoPurge(interval time.Duration) chan struct{} {
	stopChan := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				arc.Purge()
			case <-stopChan:
				return
			}
		}
	}()

	return stopChan
}
