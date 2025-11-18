package rete

import (
	"fmt"
	"sync"
	"time"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// HashJoinEngine fournit des capacités de jointure optimisées pour les nœuds Beta
type HashJoinEngine struct {
	// Hash tables pour les jointures
	leftHashTable  map[string][]*domain.Token
	rightHashTable map[string][]*domain.Fact

	// Configuration de jointure
	joinConfig JoinConfig

	// Statistiques de performance
	stats JoinStats

	// Verrou pour la concurrence
	mutex sync.RWMutex
}

// JoinConfig configure les paramètres de jointure
type JoinConfig struct {
	// Taille initiale des hash tables
	InitialHashSize int

	// Factor de croissance pour les hash tables
	GrowthFactor float64

	// Seuil pour déclencher une optimisation
	OptimizationThreshold int

	// Activer le cache des résultats de jointure
	EnableJoinCache bool

	// TTL pour le cache de jointure
	JoinCacheTTL time.Duration

	// Taille maximale du cache
	MaxCacheEntries int
}

// JoinStats contient les statistiques de performance des jointures
type JoinStats struct {
	TotalJoins       int64
	CacheHits        int64
	CacheMisses      int64
	AverageJoinTime  time.Duration
	HashTableResizes int64
	OptimizationRuns int64
}

// JoinCondition définit une condition de jointure optimisée
type OptimizedJoinCondition struct {
	LeftField    string
	RightField   string
	Operator     string
	HashFunction func(interface{}) string
}

// JoinResult contient le résultat d'une jointure
type JoinResult struct {
	LeftToken  *domain.Token
	RightFact  *domain.Fact
	JoinKey    string
	Confidence float64 // Score de confiance pour la jointure
}

// JoinCache cache les résultats de jointure fréquents
type JoinCache struct {
	cache     map[string]*CacheEntry
	maxSize   int
	ttl       time.Duration
	mutex     sync.RWMutex
	accessLog map[string]time.Time
}

// CacheEntry représente une entrée dans le cache de jointure
type CacheEntry struct {
	Results     []*JoinResult
	CreatedAt   time.Time
	AccessCount int64
}

// NewHashJoinEngine crée un nouveau moteur de jointure optimisé
func NewHashJoinEngine(config JoinConfig) *HashJoinEngine {
	if config.InitialHashSize <= 0 {
		config.InitialHashSize = 1024
	}
	if config.GrowthFactor <= 0 {
		config.GrowthFactor = 2.0
	}

	return &HashJoinEngine{
		leftHashTable:  make(map[string][]*domain.Token, config.InitialHashSize),
		rightHashTable: make(map[string][]*domain.Fact, config.InitialHashSize),
		joinConfig:     config,
		stats:          JoinStats{},
	}
}

// NewJoinCache crée un nouveau cache de jointure
func NewJoinCache(maxSize int, ttl time.Duration) *JoinCache {
	return &JoinCache{
		cache:     make(map[string]*CacheEntry),
		maxSize:   maxSize,
		ttl:       ttl,
		accessLog: make(map[string]time.Time),
	}
}

// AddLeftToken ajoute un token au côté gauche de la jointure
func (hje *HashJoinEngine) AddLeftToken(token *domain.Token, joinCondition *OptimizedJoinCondition) error {
	hje.mutex.Lock()
	defer hje.mutex.Unlock()

	// Calculer la clé de hash pour le token
	hashKey := hje.computeTokenHashKey(token, joinCondition)

	// Ajouter à la hash table gauche
	if hje.leftHashTable[hashKey] == nil {
		hje.leftHashTable[hashKey] = make([]*domain.Token, 0)
	}
	hje.leftHashTable[hashKey] = append(hje.leftHashTable[hashKey], token)

	return nil
}

// AddRightFact ajoute un fait au côté droit de la jointure
func (hje *HashJoinEngine) AddRightFact(fact *domain.Fact, joinCondition *OptimizedJoinCondition) error {
	hje.mutex.Lock()
	defer hje.mutex.Unlock()

	// Calculer la clé de hash pour le fait
	hashKey := hje.computeFactHashKey(fact, joinCondition)

	// Ajouter à la hash table droite
	if hje.rightHashTable[hashKey] == nil {
		hje.rightHashTable[hashKey] = make([]*domain.Fact, 0)
	}
	hje.rightHashTable[hashKey] = append(hje.rightHashTable[hashKey], fact)

	return nil
}

// PerformHashJoin exécute une jointure hash optimisée
func (hje *HashJoinEngine) PerformHashJoin(
	joinCondition *OptimizedJoinCondition,
	cache *JoinCache,
) ([]*JoinResult, error) {
	startTime := time.Now()
	defer func() {
		hje.updateJoinStats(time.Since(startTime))
	}()

	hje.mutex.RLock()
	defer hje.mutex.RUnlock()

	// Vérifier le cache si activé
	if cache != nil && hje.joinConfig.EnableJoinCache {
		cacheKey := hje.generateCacheKey(joinCondition)
		if cachedResults := cache.Get(cacheKey); cachedResults != nil {
			hje.stats.CacheHits++
			return cachedResults, nil
		}
		hje.stats.CacheMisses++
	}

	// Sélectionner la plus petite hash table comme table de build
	leftSize := len(hje.leftHashTable)
	rightSize := len(hje.rightHashTable)

	var results []*JoinResult
	if leftSize <= rightSize {
		// Utiliser la table gauche comme table de build
		results = hje.performLeftBuildJoin(joinCondition)
	} else {
		// Utiliser la table droite comme table de build
		results = hje.performRightBuildJoin(joinCondition)
	}

	// Mettre en cache les résultats si activé
	if cache != nil && hje.joinConfig.EnableJoinCache {
		cacheKey := hje.generateCacheKey(joinCondition)
		cache.Put(cacheKey, results)
	}

	hje.stats.TotalJoins++

	return results, nil
}

// performLeftBuildJoin exécute une jointure avec la table gauche comme table de build
func (hje *HashJoinEngine) performLeftBuildJoin(joinCondition *OptimizedJoinCondition) []*JoinResult {
	results := make([]*JoinResult, 0)

	// Parcourir la hash table gauche (table de build)
	for hashKey, leftTokens := range hje.leftHashTable {
		// Chercher les faits correspondants dans la table droite
		if rightFacts, exists := hje.rightHashTable[hashKey]; exists {
			// Évaluer toutes les combinaisons pour cette clé
			for _, leftToken := range leftTokens {
				for _, rightFact := range rightFacts {
					if hje.evaluateJoinCondition(leftToken, rightFact, joinCondition) {
						confidence := hje.calculateJoinConfidence(leftToken, rightFact, joinCondition)
						result := &JoinResult{
							LeftToken:  leftToken,
							RightFact:  rightFact,
							JoinKey:    hashKey,
							Confidence: confidence,
						}
						results = append(results, result)
					}
				}
			}
		}
	}

	return results
}

// performRightBuildJoin exécute une jointure avec la table droite comme table de build
func (hje *HashJoinEngine) performRightBuildJoin(joinCondition *OptimizedJoinCondition) []*JoinResult {
	results := make([]*JoinResult, 0)

	// Parcourir la hash table droite (table de build)
	for hashKey, rightFacts := range hje.rightHashTable {
		// Chercher les tokens correspondants dans la table gauche
		if leftTokens, exists := hje.leftHashTable[hashKey]; exists {
			// Évaluer toutes les combinaisons pour cette clé
			for _, rightFact := range rightFacts {
				for _, leftToken := range leftTokens {
					if hje.evaluateJoinCondition(leftToken, rightFact, joinCondition) {
						confidence := hje.calculateJoinConfidence(leftToken, rightFact, joinCondition)
						result := &JoinResult{
							LeftToken:  leftToken,
							RightFact:  rightFact,
							JoinKey:    hashKey,
							Confidence: confidence,
						}
						results = append(results, result)
					}
				}
			}
		}
	}

	return results
}

// computeTokenHashKey calcule la clé de hash pour un token
func (hje *HashJoinEngine) computeTokenHashKey(token *domain.Token, joinCondition *OptimizedJoinCondition) string {
	// Utiliser une fonction de hash personnalisée si fournie
	if joinCondition.HashFunction != nil {
		if value := hje.extractTokenFieldValue(token, joinCondition.LeftField); value != nil {
			return joinCondition.HashFunction(value)
		}
	}

	// Hash par défaut basé sur la valeur du champ
	if value := hje.extractTokenFieldValue(token, joinCondition.LeftField); value != nil {
		return fmt.Sprintf("%v", value)
	}

	return "NULL"
}

// computeFactHashKey calcule la clé de hash pour un fait
func (hje *HashJoinEngine) computeFactHashKey(fact *domain.Fact, joinCondition *OptimizedJoinCondition) string {
	// Utiliser une fonction de hash personnalisée si fournie
	if joinCondition.HashFunction != nil {
		if value, exists := fact.Fields[joinCondition.RightField]; exists {
			return joinCondition.HashFunction(value)
		}
	}

	// Hash par défaut basé sur la valeur du champ
	if value, exists := fact.Fields[joinCondition.RightField]; exists {
		return fmt.Sprintf("%v", value)
	}

	return "NULL"
}

// extractTokenFieldValue extrait la valeur d'un champ depuis un token
func (hje *HashJoinEngine) extractTokenFieldValue(token *domain.Token, fieldName string) interface{} {
	// Chercher dans tous les faits du token
	for _, fact := range token.Facts {
		if value, exists := fact.Fields[fieldName]; exists {
			return value
		}
	}
	return nil
}

// evaluateJoinCondition évalue si deux éléments satisfont la condition de jointure
func (hje *HashJoinEngine) evaluateJoinCondition(
	token *domain.Token,
	fact *domain.Fact,
	joinCondition *OptimizedJoinCondition,
) bool {
	leftValue := hje.extractTokenFieldValue(token, joinCondition.LeftField)
	rightValue, exists := fact.Fields[joinCondition.RightField]

	if leftValue == nil || !exists {
		return false
	}

	switch joinCondition.Operator {
	case "==":
		return fmt.Sprintf("%v", leftValue) == fmt.Sprintf("%v", rightValue)
	case "!=":
		return fmt.Sprintf("%v", leftValue) != fmt.Sprintf("%v", rightValue)
	// Ajouter d'autres opérateurs selon les besoins
	default:
		return false
	}
}

// calculateJoinConfidence calcule un score de confiance pour la jointure
func (hje *HashJoinEngine) calculateJoinConfidence(
	token *domain.Token,
	fact *domain.Fact,
	joinCondition *OptimizedJoinCondition,
) float64 {
	// Score de base
	confidence := 1.0

	// Ajuster basé sur la fraîcheur des données
	now := time.Now()
	for _, tokenFact := range token.Facts {
		age := now.Sub(tokenFact.Timestamp)
		if age > time.Hour {
			confidence *= 0.9 // Réduire la confiance pour les données anciennes
		}
	}

	factAge := now.Sub(fact.Timestamp)
	if factAge > time.Hour {
		confidence *= 0.9
	}

	return confidence
}

// generateCacheKey génère une clé de cache pour une condition de jointure
func (hje *HashJoinEngine) generateCacheKey(joinCondition *OptimizedJoinCondition) string {
	return fmt.Sprintf("join:%s_%s_%s",
		joinCondition.LeftField,
		joinCondition.Operator,
		joinCondition.RightField)
}

// updateJoinStats met à jour les statistiques de jointure
func (hje *HashJoinEngine) updateJoinStats(duration time.Duration) {
	// Calculer la moyenne mobile des temps de jointure
	totalJoins := hje.stats.TotalJoins
	if totalJoins == 0 {
		hje.stats.AverageJoinTime = duration
	} else {
		// Moyenne pondérée
		hje.stats.AverageJoinTime = time.Duration(
			(int64(hje.stats.AverageJoinTime)*totalJoins + int64(duration)) / (totalJoins + 1),
		)
	}
}

// Get récupère des résultats depuis le cache
func (jc *JoinCache) Get(key string) []*JoinResult {
	jc.mutex.RLock()
	defer jc.mutex.RUnlock()

	entry, exists := jc.cache[key]
	if !exists {
		return nil
	}

	// Vérifier la TTL
	if time.Since(entry.CreatedAt) > jc.ttl {
		delete(jc.cache, key)
		delete(jc.accessLog, key)
		return nil
	}

	// Mettre à jour les statistiques d'accès
	entry.AccessCount++
	jc.accessLog[key] = time.Now()

	return entry.Results
}

// Put stocke des résultats dans le cache
func (jc *JoinCache) Put(key string, results []*JoinResult) {
	jc.mutex.Lock()
	defer jc.mutex.Unlock()

	// Éviction si le cache est plein
	if len(jc.cache) >= jc.maxSize {
		jc.evictLRU()
	}

	// Stocker la nouvelle entrée
	jc.cache[key] = &CacheEntry{
		Results:     results,
		CreatedAt:   time.Now(),
		AccessCount: 1,
	}
	jc.accessLog[key] = time.Now()
}

// evictLRU évince l'entrée la moins récemment utilisée
func (jc *JoinCache) evictLRU() {
	var oldestKey string
	var oldestTime time.Time

	for key, lastAccess := range jc.accessLog {
		if oldestKey == "" || lastAccess.Before(oldestTime) {
			oldestKey = key
			oldestTime = lastAccess
		}
	}

	if oldestKey != "" {
		delete(jc.cache, oldestKey)
		delete(jc.accessLog, oldestKey)
	}
}

// Clear vide le cache
func (jc *JoinCache) Clear() {
	jc.mutex.Lock()
	defer jc.mutex.Unlock()

	jc.cache = make(map[string]*CacheEntry)
	jc.accessLog = make(map[string]time.Time)
}

// GetStats retourne les statistiques de performance
func (hje *HashJoinEngine) GetStats() JoinStats {
	hje.mutex.RLock()
	defer hje.mutex.RUnlock()

	return hje.stats
}

// OptimizeHashTables optimise les hash tables basées sur l'usage
func (hje *HashJoinEngine) OptimizeHashTables() {
	hje.mutex.Lock()
	defer hje.mutex.Unlock()

	// Redimensionner si nécessaire
	leftLoad := float64(len(hje.leftHashTable)) / float64(hje.joinConfig.InitialHashSize)
	rightLoad := float64(len(hje.rightHashTable)) / float64(hje.joinConfig.InitialHashSize)

	if leftLoad > 0.75 || rightLoad > 0.75 {
		// Redimensionner les tables
		hje.resizeHashTables()
		hje.stats.HashTableResizes++
	}

	hje.stats.OptimizationRuns++
}

// resizeHashTables redimensionne les hash tables
func (hje *HashJoinEngine) resizeHashTables() {
	newSize := int(float64(hje.joinConfig.InitialHashSize) * hje.joinConfig.GrowthFactor)

	// Créer de nouvelles tables plus grandes
	newLeftTable := make(map[string][]*domain.Token, newSize)
	newRightTable := make(map[string][]*domain.Fact, newSize)

	// Copier les données existantes
	for key, tokens := range hje.leftHashTable {
		newLeftTable[key] = tokens
	}
	for key, facts := range hje.rightHashTable {
		newRightTable[key] = facts
	}

	// Remplacer les anciennes tables
	hje.leftHashTable = newLeftTable
	hje.rightHashTable = newRightTable
	hje.joinConfig.InitialHashSize = newSize
}
