// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

// AlphaSharingRegistry gère le partage des AlphaNodes entre plusieurs règles
// qui ont des conditions alpha identiques
type AlphaSharingRegistry struct {
	// Map[ConditionHash] -> AlphaNode
	sharedAlphaNodes map[string]*AlphaNode
	// Cache pour les hash de conditions (Map[conditionJSON] -> hash)
	hashCache map[string]string
	// Cache LRU pour les hash (si configuré)
	lruHashCache *LRUCache
	// Configuration de performance
	config *ChainPerformanceConfig
	// Métriques de performance
	metrics *ChainBuildMetrics
	mutex   sync.RWMutex
}

// NewAlphaSharingRegistry crée un nouveau registre de partage d'AlphaNodes
// Utilise la configuration par défaut
func NewAlphaSharingRegistry() *AlphaSharingRegistry {
	config := DefaultChainPerformanceConfig()
	return NewAlphaSharingRegistryWithConfig(config, NewChainBuildMetrics())
}

// NewAlphaSharingRegistryWithMetrics crée un registre avec des métriques partagées
// Utilise la configuration par défaut
func NewAlphaSharingRegistryWithMetrics(metrics *ChainBuildMetrics) *AlphaSharingRegistry {
	config := DefaultChainPerformanceConfig()
	return NewAlphaSharingRegistryWithConfig(config, metrics)
}

// NewAlphaSharingRegistryWithConfig crée un registre avec une configuration personnalisée
func NewAlphaSharingRegistryWithConfig(config *ChainPerformanceConfig, metrics *ChainBuildMetrics) *AlphaSharingRegistry {
	if config == nil {
		config = DefaultChainPerformanceConfig()
	}
	if metrics == nil {
		metrics = NewChainBuildMetrics()
	}

	asr := &AlphaSharingRegistry{
		sharedAlphaNodes: make(map[string]*AlphaNode),
		config:           config,
		metrics:          metrics,
	}

	// Initialiser le cache approprié selon la configuration
	if config.HashCacheEnabled {
		if config.HashCacheEviction == EvictionPolicyLRU {
			// Utiliser le cache LRU
			asr.lruHashCache = NewLRUCache(config.HashCacheMaxSize, config.HashCacheTTL)
		} else {
			// Utiliser le simple map pour EvictionPolicyNone
			asr.hashCache = make(map[string]string)
		}
	}

	return asr
}

// ConditionHash calcule un hash unique pour une condition alpha
// Deux conditions identiques produiront le même hash
func ConditionHash(condition interface{}, variableName string) (string, error) {
	// Déballer la condition si elle est wrappée (pour le partage entre règles simples et chaînes)
	unwrapped := normalizeConditionForSharing(condition)

	// Normaliser la condition pour assurer un hash cohérent
	normalized, err := normalizeCondition(unwrapped)
	if err != nil {
		return "", fmt.Errorf("erreur normalisation condition: %w", err)
	}

	// Créer une structure canonique incluant le nom de variable
	// car une condition sur "p" vs "q" est différente
	canonical := map[string]interface{}{
		"condition": normalized,
		"variable":  variableName,
	}

	// Sérialiser en JSON avec clés triées pour cohérence
	jsonBytes, err := json.Marshal(canonical)
	if err != nil {
		return "", fmt.Errorf("erreur sérialisation condition: %w", err)
	}

	// Calculer le hash SHA-256
	hash := sha256.Sum256(jsonBytes)
	return fmt.Sprintf("alpha_%x", hash[:8]), nil // Utiliser les 8 premiers octets pour l'ID
}

// ConditionHashCached calcule un hash avec cache pour améliorer les performances
func (asr *AlphaSharingRegistry) ConditionHashCached(condition interface{}, variableName string) (string, error) {
	// Si le cache n'est pas activé, calculer directement
	if !asr.isCacheEnabled() {
		return ConditionHash(condition, variableName)
	}

	// Déballer et normaliser pour créer la clé de cache
	unwrapped := normalizeConditionForSharing(condition)
	normalized, err := normalizeCondition(unwrapped)
	if err != nil {
		return "", fmt.Errorf("erreur normalisation condition: %w", err)
	}

	canonical := map[string]interface{}{
		"condition": normalized,
		"variable":  variableName,
	}

	// Créer la clé de cache
	jsonBytes, err := json.Marshal(canonical)
	if err != nil {
		return "", fmt.Errorf("erreur sérialisation condition: %w", err)
	}
	cacheKey := string(jsonBytes)

	// Utiliser le cache LRU si configuré
	if asr.lruHashCache != nil {
		if cachedValue, found := asr.lruHashCache.Get(cacheKey); found {
			if asr.metrics != nil {
				asr.metrics.RecordHashCacheHit()
			}
			return cachedValue.(string), nil
		}

		// Cache miss
		if asr.metrics != nil {
			asr.metrics.RecordHashCacheMiss()
		}

		// Calculer le hash
		hash := sha256.Sum256(jsonBytes)
		hashStr := fmt.Sprintf("alpha_%x", hash[:8])

		// Stocker dans le cache LRU
		asr.lruHashCache.Set(cacheKey, hashStr)
		if asr.metrics != nil {
			asr.metrics.UpdateHashCacheSize(asr.lruHashCache.Len())
		}

		return hashStr, nil
	}

	// Utiliser le simple map cache (fallback)
	asr.mutex.RLock()
	if cachedHash, exists := asr.hashCache[cacheKey]; exists {
		asr.mutex.RUnlock()
		if asr.metrics != nil {
			asr.metrics.RecordHashCacheHit()
		}
		return cachedHash, nil
	}
	asr.mutex.RUnlock()

	// Cache miss - calculer le hash
	if asr.metrics != nil {
		asr.metrics.RecordHashCacheMiss()
	}

	hash := sha256.Sum256(jsonBytes)
	hashStr := fmt.Sprintf("alpha_%x", hash[:8])

	// Stocker dans le cache
	asr.mutex.Lock()
	asr.hashCache[cacheKey] = hashStr
	if asr.metrics != nil {
		asr.metrics.UpdateHashCacheSize(len(asr.hashCache))
	}
	asr.mutex.Unlock()

	return hashStr, nil
}

// isCacheEnabled vérifie si le cache est activé
func (asr *AlphaSharingRegistry) isCacheEnabled() bool {
	return asr.config != nil && asr.config.HashCacheEnabled
}

// normalizeConditionForSharing déballe les conditions wrappées pour permettre le partage
// entre règles simples (qui wrappent dans {"type": "constraint", "constraint": X})
// et chaînes (qui utilisent directement la condition décomposée)
func normalizeConditionForSharing(condition interface{}) interface{} {
	// Si la condition est une map
	if condMap, ok := condition.(map[string]interface{}); ok {
		// Vérifier si c'est une condition wrappée dans un type "constraint"
		if condType, hasType := condMap["type"]; hasType {
			if condTypeStr, ok := condType.(string); ok && condTypeStr == "constraint" {
				// Déballer la condition interne
				if innerCond, hasConstraint := condMap["constraint"]; hasConstraint {
					// Récursion pour déballer plusieurs niveaux si nécessaire
					return normalizeConditionForSharing(innerCond)
				}
			}
		}

		// Normaliser les types équivalents pour le partage
		// "comparison" et "binaryOperation" sont des synonymes
		normalized := make(map[string]interface{})
		for key, value := range condMap {
			if key == "type" {
				if typeStr, ok := value.(string); ok {
					// Normaliser "comparison" vers "binaryOperation"
					if typeStr == "comparison" {
						normalized[key] = "binaryOperation"
					} else {
						normalized[key] = value
					}
				} else {
					normalized[key] = value
				}
			} else {
				// Normaliser récursivement les valeurs imbriquées
				normalized[key] = normalizeConditionForSharing(value)
			}
		}
		return normalized
	}

	// Si c'est un slice, normaliser chaque élément
	if slice, ok := condition.([]interface{}); ok {
		normalized := make([]interface{}, len(slice))
		for i, item := range slice {
			normalized[i] = normalizeConditionForSharing(item)
		}
		return normalized
	}

	// Sinon, retourner la condition telle quelle
	return condition
}

// normalizeCondition normalise une condition pour assurer un hash cohérent
func normalizeCondition(condition interface{}) (interface{}, error) {
	if condition == nil {
		return map[string]interface{}{"type": "simple"}, nil
	}

	switch cond := condition.(type) {
	case map[string]interface{}:
		// Copier la map et normaliser récursivement
		normalized := make(map[string]interface{})
		for key, value := range cond {
			normalizedValue, err := normalizeCondition(value)
			if err != nil {
				return nil, err
			}
			normalized[key] = normalizedValue
		}
		return normalized, nil

	case []interface{}:
		// Normaliser chaque élément du slice
		normalized := make([]interface{}, len(cond))
		for i, item := range cond {
			normalizedItem, err := normalizeCondition(item)
			if err != nil {
				return nil, err
			}
			normalized[i] = normalizedItem
		}
		return normalized, nil

	default:
		// Types primitifs: retourner tel quel
		return condition, nil
	}
}

// GetOrCreateAlphaNode récupère un AlphaNode existant ou en crée un nouveau
// Si un AlphaNode avec la même condition existe, il est réutilisé
func (asr *AlphaSharingRegistry) GetOrCreateAlphaNode(
	condition interface{},
	variableName string,
	storage Storage,
) (*AlphaNode, string, bool, error) {
	// Calculer le hash de la condition avec cache
	hash, err := asr.ConditionHashCached(condition, variableName)
	if err != nil {
		return nil, "", false, fmt.Errorf("erreur calcul hash condition: %w", err)
	}

	// Vérifier si un AlphaNode existe déjà pour cette condition
	asr.mutex.RLock()
	existingNode, exists := asr.sharedAlphaNodes[hash]
	asr.mutex.RUnlock()

	if exists {
		// AlphaNode partagé trouvé
		return existingNode, hash, true, nil
	}

	// Créer un nouveau AlphaNode
	asr.mutex.Lock()
	defer asr.mutex.Unlock()

	// Double-check après avoir acquis le verrou d'écriture
	if existingNode, exists := asr.sharedAlphaNodes[hash]; exists {
		return existingNode, hash, true, nil
	}

	// Créer le nouveau nœud avec l'ID basé sur le hash
	alphaNode := NewAlphaNode(hash, condition, variableName, storage)
	asr.sharedAlphaNodes[hash] = alphaNode

	return alphaNode, hash, false, nil
}

// GetAlphaNode récupère un AlphaNode partagé par son hash
func (asr *AlphaSharingRegistry) GetAlphaNode(hash string) (*AlphaNode, bool) {
	asr.mutex.RLock()
	defer asr.mutex.RUnlock()

	node, exists := asr.sharedAlphaNodes[hash]
	return node, exists
}

// RemoveAlphaNode supprime un AlphaNode du registre
// Cette méthode doit être appelée uniquement quand plus aucune règle n'utilise ce nœud
func (asr *AlphaSharingRegistry) RemoveAlphaNode(hash string) error {
	asr.mutex.Lock()
	defer asr.mutex.Unlock()

	if _, exists := asr.sharedAlphaNodes[hash]; !exists {
		return fmt.Errorf("AlphaNode %s non trouvé dans le registre", hash)
	}

	delete(asr.sharedAlphaNodes, hash)
	return nil
}

// GetStats retourne des statistiques sur le partage des AlphaNodes
func (asr *AlphaSharingRegistry) GetStats() map[string]interface{} {
	asr.mutex.RLock()
	defer asr.mutex.RUnlock()

	totalNodes := len(asr.sharedAlphaNodes)

	// Compter le nombre total de règles utilisant ces nœuds
	totalRuleReferences := 0
	childCounts := make([]int, 0, totalNodes)

	for _, node := range asr.sharedAlphaNodes {
		childCount := len(node.GetChildren())
		childCounts = append(childCounts, childCount)
		totalRuleReferences += childCount
	}

	// Calculer la moyenne
	avgSharing := 0.0
	if totalNodes > 0 {
		avgSharing = float64(totalRuleReferences) / float64(totalNodes)
	}

	return map[string]interface{}{
		"total_shared_alpha_nodes": totalNodes,
		"total_rule_references":    totalRuleReferences,
		"average_sharing_ratio":    avgSharing,
	}
}

// ListSharedAlphaNodes retourne la liste de tous les AlphaNodes partagés
func (asr *AlphaSharingRegistry) ListSharedAlphaNodes() []string {
	asr.mutex.RLock()
	defer asr.mutex.RUnlock()

	hashes := make([]string, 0, len(asr.sharedAlphaNodes))
	for hash := range asr.sharedAlphaNodes {
		hashes = append(hashes, hash)
	}

	// Trier pour avoir un ordre déterministe
	sort.Strings(hashes)
	return hashes
}

// Reset réinitialise complètement le registre
func (asr *AlphaSharingRegistry) Reset() {
	asr.mutex.Lock()
	defer asr.mutex.Unlock()

	asr.sharedAlphaNodes = make(map[string]*AlphaNode)

	if asr.lruHashCache != nil {
		asr.lruHashCache.Clear()
	} else if asr.hashCache != nil {
		asr.hashCache = make(map[string]string)
	}

	if asr.metrics != nil {
		asr.metrics.UpdateHashCacheSize(0)
	}
}

// ClearHashCache vide uniquement le cache de hash
func (asr *AlphaSharingRegistry) ClearHashCache() {
	if asr.lruHashCache != nil {
		asr.lruHashCache.Clear()
		if asr.metrics != nil {
			asr.metrics.UpdateHashCacheSize(0)
		}
		return
	}

	asr.mutex.Lock()
	defer asr.mutex.Unlock()

	if asr.hashCache != nil {
		asr.hashCache = make(map[string]string)
	}
	if asr.metrics != nil {
		asr.metrics.UpdateHashCacheSize(0)
	}
}

// GetHashCacheSize retourne la taille actuelle du cache de hash
func (asr *AlphaSharingRegistry) GetHashCacheSize() int {
	if asr.lruHashCache != nil {
		return asr.lruHashCache.Len()
	}

	asr.mutex.RLock()
	defer asr.mutex.RUnlock()
	return len(asr.hashCache)
}

// GetMetrics retourne les métriques de performance
func (asr *AlphaSharingRegistry) GetMetrics() *ChainBuildMetrics {
	return asr.metrics
}

// GetHashCacheStats retourne les statistiques détaillées du cache de hash
func (asr *AlphaSharingRegistry) GetHashCacheStats() map[string]interface{} {
	if asr.lruHashCache != nil {
		stats := asr.lruHashCache.GetStats()
		return map[string]interface{}{
			"type":          "lru",
			"size":          stats.Size,
			"capacity":      stats.Capacity,
			"hits":          stats.Hits,
			"misses":        stats.Misses,
			"evictions":     stats.Evictions,
			"sets":          stats.Sets,
			"hit_rate":      stats.HitRate(),
			"eviction_rate": stats.EvictionRate(),
			"fill_rate":     stats.FillRate(),
		}
	}

	asr.mutex.RLock()
	defer asr.mutex.RUnlock()

	size := 0
	if asr.hashCache != nil {
		size = len(asr.hashCache)
	}

	return map[string]interface{}{
		"type": "simple_map",
		"size": size,
	}
}

// GetConfig retourne la configuration actuelle
func (asr *AlphaSharingRegistry) GetConfig() *ChainPerformanceConfig {
	return asr.config
}

// CleanExpiredHashCache nettoie les entrées expirées du cache (si LRU avec TTL)
func (asr *AlphaSharingRegistry) CleanExpiredHashCache() int {
	if asr.lruHashCache != nil {
		cleaned := asr.lruHashCache.CleanExpired()
		if cleaned > 0 && asr.metrics != nil {
			asr.metrics.UpdateHashCacheSize(asr.lruHashCache.Len())
		}
		return cleaned
	}
	return 0
}

// GetSharedAlphaNodeDetails retourne les détails d'un AlphaNode partagé
func (asr *AlphaSharingRegistry) GetSharedAlphaNodeDetails(hash string) map[string]interface{} {
	asr.mutex.RLock()
	defer asr.mutex.RUnlock()

	node, exists := asr.sharedAlphaNodes[hash]
	if !exists {
		return nil
	}

	children := node.GetChildren()
	childIDs := make([]string, len(children))
	for i, child := range children {
		childIDs[i] = child.GetID()
	}

	return map[string]interface{}{
		"hash":          hash,
		"node_id":       node.GetID(),
		"variable_name": node.VariableName,
		"condition":     node.Condition,
		"child_count":   len(children),
		"child_ids":     childIDs,
	}
}

// Clear vide tous les caches et nodes partagés
func (asr *AlphaSharingRegistry) Clear() {
	asr.mutex.Lock()
	defer asr.mutex.Unlock()

	asr.sharedAlphaNodes = make(map[string]*AlphaNode)
	asr.hashCache = make(map[string]string)

	if asr.lruHashCache != nil {
		asr.lruHashCache.Clear()
	}
}
