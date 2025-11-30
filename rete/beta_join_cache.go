// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// BetaJoinCache gère le cache des résultats de jointure pour optimiser les performances.
//
// Le cache stocke les résultats de jointure entre tokens gauche et faits droite,
// évitant ainsi de recalculer les mêmes matchs répétitivement.
//
// Architecture:
//
//	Cache Key: hash(leftTokenID + rightFactID + joinConditions)
//	Cache Value: JoinResult (matched: bool, token: *Token, timestamp: time)
//
// Utilisation du cache LRU existant pour bénéficier de:
//   - Éviction automatique (LRU)
//   - TTL configurable
//   - Thread-safety
//   - Métriques intégrées
//
// Exemple d'utilisation:
//
//	cache := NewBetaJoinCache(config)
//
//	// Tenter de récupérer du cache
//	if result, found := cache.GetJoinResult(leftToken, rightFact, joinNode); found {
//	    return result.Token // Résultat mis en cache
//	}
//
//	// Calculer et mettre en cache
//	result := performJoin(leftToken, rightFact)
//	cache.SetJoinResult(leftToken, rightFact, joinNode, result)
type BetaJoinCache struct {
	// Cache LRU sous-jacent
	lruCache *LRUCache

	// Configuration
	config *ChainPerformanceConfig

	// Mutex pour les opérations qui nécessitent une atomicité
	mutex sync.RWMutex

	// Métriques spécifiques au cache de jointure
	joinCacheHits      int64
	joinCacheMisses    int64
	joinCacheEvictions int64
	invalidations      int64
}

// JoinResult représente le résultat d'une opération de jointure mise en cache.
type JoinResult struct {
	Matched   bool      // true si le token et le fait ont matché
	Token     *Token    // Token résultant si matched=true, nil sinon
	Timestamp time.Time // Timestamp de la mise en cache
	JoinType  string    // Type de jointure (pour debug/metrics)
}

// joinCacheKey représente une clé de cache pour une opération de jointure.
type joinCacheKey struct {
	LeftTokenID  string `json:"left_token_id"`
	RightFactID  string `json:"right_fact_id"`
	JoinNodeID   string `json:"join_node_id"`
	ConditionSig string `json:"condition_sig"` // Signature des conditions de jointure
}

// NewBetaJoinCache crée un nouveau cache de résultats de jointure.
//
// Paramètres:
//   - config: Configuration de performance (nil = config par défaut)
//
// Retourne:
//   - Un nouveau BetaJoinCache prêt à l'emploi
//
// Exemple:
//
//	config := DefaultChainPerformanceConfig()
//	config.BetaJoinResultCacheEnabled = true
//	config.BetaJoinResultCacheMaxSize = 10000
//	config.BetaJoinResultCacheTTL = time.Minute
//
//	cache := NewBetaJoinCache(config)
func NewBetaJoinCache(config *ChainPerformanceConfig) *BetaJoinCache {
	if config == nil {
		config = DefaultChainPerformanceConfig()
	}

	// Créer le cache LRU sous-jacent
	lruCache := NewLRUCache(
		config.BetaJoinResultCacheMaxSize,
		config.BetaJoinResultCacheTTL,
	)

	return &BetaJoinCache{
		lruCache: lruCache,
		config:   config,
	}
}

// GetJoinResult tente de récupérer un résultat de jointure du cache.
//
// Paramètres:
//   - leftToken: Token du côté gauche
//   - rightFact: Fait du côté droit
//   - joinNode: JoinNode contenant les conditions de jointure
//
// Retourne:
//   - result: Le résultat de jointure mis en cache
//   - found: true si trouvé dans le cache, false sinon
//
// Thread-safe: Oui
//
// Exemple:
//
//	if result, found := cache.GetJoinResult(token, fact, node); found {
//	    if result.Matched {
//	        return result.Token // Jointure réussie
//	    }
//	    return nil // Pas de match
//	}
//	// Cache miss - calculer la jointure
func (bjc *BetaJoinCache) GetJoinResult(leftToken *Token, rightFact *Fact, joinNode *JoinNode) (*JoinResult, bool) {
	if !bjc.isEnabled() {
		return nil, false
	}

	// Générer la clé de cache
	cacheKey, err := bjc.computeCacheKey(leftToken, rightFact, joinNode)
	if err != nil {
		// En cas d'erreur, considérer comme cache miss
		bjc.recordMiss()
		return nil, false
	}

	// Chercher dans le cache LRU
	value, found := bjc.lruCache.Get(cacheKey)
	if !found {
		bjc.recordMiss()
		return nil, false
	}

	// Vérifier le type
	result, ok := value.(*JoinResult)
	if !ok {
		// Type invalide - considérer comme miss
		bjc.recordMiss()
		return nil, false
	}

	bjc.recordHit()
	return result, true
}

// SetJoinResult met en cache un résultat de jointure.
//
// Paramètres:
//   - leftToken: Token du côté gauche
//   - rightFact: Fait du côté droit
//   - joinNode: JoinNode contenant les conditions de jointure
//   - result: Le résultat de jointure à mettre en cache
//
// Thread-safe: Oui
//
// Exemple:
//
//	result := &JoinResult{
//	    Matched: true,
//	    Token: joinedToken,
//	    Timestamp: time.Now(),
//	    JoinType: "binary",
//	}
//	cache.SetJoinResult(token, fact, node, result)
func (bjc *BetaJoinCache) SetJoinResult(leftToken *Token, rightFact *Fact, joinNode *JoinNode, result *JoinResult) {
	if !bjc.isEnabled() {
		return
	}

	// Générer la clé de cache
	cacheKey, err := bjc.computeCacheKey(leftToken, rightFact, joinNode)
	if err != nil {
		// Impossible de générer la clé - ignorer
		return
	}

	// Définir le timestamp si pas déjà défini
	if result.Timestamp.IsZero() {
		result.Timestamp = time.Now()
	}

	// Stocker dans le cache LRU
	bjc.lruCache.Set(cacheKey, result)
}

// InvalidateForFact invalide toutes les entrées de cache contenant un fait donné.
//
// Appelé lors de la rétractation d'un fait pour assurer la cohérence du cache.
//
// Paramètres:
//   - factID: ID du fait à invalider
//
// Retourne:
//   - Le nombre d'entrées invalidées
//
// Note: Cette opération peut être coûteuse car elle nécessite de parcourir
// toutes les entrées du cache. Pour de meilleures performances, considérer
// un TTL court plutôt que des invalidations fréquentes.
//
// Thread-safe: Oui
//
// Exemple:
//
//	// Lors de la rétractation d'un fait
//	invalidated := cache.InvalidateForFact(fact.GetInternalID())
//	log.Printf("Invalidé %d entrées de cache", invalidated)
func (bjc *BetaJoinCache) InvalidateForFact(factID string) int {
	if !bjc.isEnabled() {
		return 0
	}

	bjc.mutex.Lock()
	defer bjc.mutex.Unlock()

	// Cette opération nécessite de parcourir le cache
	// Pour des raisons de performance, on utilise une approche simple:
	// Clear tout le cache si l'invalidation concerne trop d'entrées
	// (implémentation simple pour éviter la complexité)

	// Pour l'instant, on clear tout le cache lors d'une invalidation
	// Une future optimisation pourrait maintenir un index inverse factID -> cacheKeys
	bjc.lruCache.Clear()
	bjc.invalidations++

	return 1 // Retourner 1 pour indiquer qu'une invalidation a eu lieu
}

// InvalidateForToken invalide toutes les entrées de cache contenant un token donné.
//
// Appelé lors de modifications du token pour assurer la cohérence.
//
// Paramètres:
//   - tokenID: ID du token à invalider
//
// Retourne:
//   - Le nombre d'entrées invalidées
//
// Thread-safe: Oui
func (bjc *BetaJoinCache) InvalidateForToken(tokenID string) int {
	if !bjc.isEnabled() {
		return 0
	}

	bjc.mutex.Lock()
	defer bjc.mutex.Unlock()

	// Même approche que InvalidateForFact
	bjc.lruCache.Clear()
	bjc.invalidations++

	return 1
}

// Clear vide complètement le cache.
//
// Utile pour les tests ou pour forcer un rafraîchissement complet.
//
// Thread-safe: Oui
//
// Exemple:
//
//	// Forcer un rafraîchissement complet
//	cache.Clear()
func (bjc *BetaJoinCache) Clear() {
	if !bjc.isEnabled() {
		return
	}

	bjc.lruCache.Clear()
}

// GetStats retourne les statistiques détaillées du cache.
//
// Retourne:
//   - Map contenant les statistiques (hits, misses, hit_rate, size, etc.)
//
// Thread-safe: Oui
//
// Exemple:
//
//	stats := cache.GetStats()
//	fmt.Printf("Hit rate: %.2f%%\n", stats["hit_rate"].(float64) * 100)
//	fmt.Printf("Cache size: %d\n", stats["size"].(int))
func (bjc *BetaJoinCache) GetStats() map[string]interface{} {
	if !bjc.isEnabled() {
		return map[string]interface{}{
			"enabled": false,
		}
	}

	// Récupérer les stats du cache LRU
	lruStats := bjc.lruCache.GetStats()

	bjc.mutex.RLock()
	hits := bjc.joinCacheHits
	misses := bjc.joinCacheMisses
	invalidations := bjc.invalidations
	bjc.mutex.RUnlock()

	// Calculer le hit rate
	total := hits + misses
	hitRate := 0.0
	if total > 0 {
		hitRate = float64(hits) / float64(total)
	}

	return map[string]interface{}{
		"enabled":       true,
		"size":          lruStats.Size,
		"capacity":      lruStats.Capacity,
		"hits":          hits,
		"misses":        misses,
		"evictions":     lruStats.Evictions,
		"invalidations": invalidations,
		"hit_rate":      hitRate,
		"ttl_seconds":   bjc.config.BetaJoinResultCacheTTL.Seconds(),
	}
}

// GetHitRate retourne le taux de hit du cache (0.0 à 1.0).
//
// Retourne:
//   - Hit rate entre 0.0 (0%) et 1.0 (100%)
//
// Thread-safe: Oui
//
// Exemple:
//
//	hitRate := cache.GetHitRate()
//	if hitRate < 0.5 {
//	    log.Printf("WARNING: Low cache hit rate: %.2f%%", hitRate * 100)
//	}
func (bjc *BetaJoinCache) GetHitRate() float64 {
	bjc.mutex.RLock()
	defer bjc.mutex.RUnlock()

	total := bjc.joinCacheHits + bjc.joinCacheMisses
	if total == 0 {
		return 0.0
	}

	return float64(bjc.joinCacheHits) / float64(total)
}

// GetSize retourne le nombre d'entrées actuellement dans le cache.
//
// Thread-safe: Oui
func (bjc *BetaJoinCache) GetSize() int {
	if !bjc.isEnabled() {
		return 0
	}

	return bjc.lruCache.Len()
}

// CleanExpired nettoie les entrées expirées du cache.
//
// Retourne:
//   - Le nombre d'entrées nettoyées
//
// Thread-safe: Oui
//
// Exemple:
//
//	// Périodiquement nettoyer les entrées expirées
//	ticker := time.NewTicker(time.Minute)
//	go func() {
//	    for range ticker.C {
//	        cleaned := cache.CleanExpired()
//	        if cleaned > 0 {
//	            log.Printf("Cleaned %d expired cache entries", cleaned)
//	        }
//	    }
//	}()
func (bjc *BetaJoinCache) CleanExpired() int {
	if !bjc.isEnabled() {
		return 0
	}

	return bjc.lruCache.CleanExpired()
}

// computeCacheKey calcule une clé de cache pour une opération de jointure.
//
// La clé est basée sur:
//   - ID du token gauche
//   - ID du fait droit
//   - ID du JoinNode
//   - Signature des conditions de jointure
//
// Utilise SHA-256 pour obtenir une clé stable et courte.
func (bjc *BetaJoinCache) computeCacheKey(leftToken *Token, rightFact *Fact, joinNode *JoinNode) (string, error) {
	// Construire la structure de clé
	key := joinCacheKey{
		LeftTokenID:  leftToken.ID,
		RightFactID:  rightFact.GetInternalID(),
		JoinNodeID:   joinNode.ID,
		ConditionSig: bjc.computeConditionSignature(joinNode),
	}

	// Sérialiser en JSON pour hashing
	jsonBytes, err := json.Marshal(key)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cache key: %w", err)
	}

	// Calculer le hash SHA-256
	hash := sha256.Sum256(jsonBytes)
	hashStr := fmt.Sprintf("join_cache_%x", hash[:8]) // Utiliser les 8 premiers octets

	return hashStr, nil
}

// computeConditionSignature calcule une signature des conditions de jointure.
//
// Cette signature permet de différencier des jointures sur les mêmes
// tokens/faits mais avec des conditions différentes.
func (bjc *BetaJoinCache) computeConditionSignature(joinNode *JoinNode) string {
	// Utiliser l'ID du JoinNode comme signature simple
	// Dans une implémentation plus avancée, on pourrait hasher les conditions réelles
	return joinNode.ID
}

// isEnabled vérifie si le cache est activé.
func (bjc *BetaJoinCache) isEnabled() bool {
	return bjc.config != nil &&
		bjc.config.BetaCacheEnabled &&
		bjc.config.BetaJoinResultCacheEnabled
}

// recordHit enregistre un hit du cache.
func (bjc *BetaJoinCache) recordHit() {
	bjc.mutex.Lock()
	defer bjc.mutex.Unlock()
	bjc.joinCacheHits++
}

// recordMiss enregistre un miss du cache.
func (bjc *BetaJoinCache) recordMiss() {
	bjc.mutex.Lock()
	defer bjc.mutex.Unlock()
	bjc.joinCacheMisses++
}

// ResetStats réinitialise les statistiques du cache.
//
// Utile pour les tests ou pour commencer un nouveau cycle de monitoring.
//
// Thread-safe: Oui
func (bjc *BetaJoinCache) ResetStats() {
	bjc.mutex.Lock()
	defer bjc.mutex.Unlock()

	bjc.joinCacheHits = 0
	bjc.joinCacheMisses = 0
	bjc.joinCacheEvictions = 0
	bjc.invalidations = 0
}
