package rete

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// EvaluationCache fournit un cache intelligent pour les résultats d'évaluation
type EvaluationCache struct {
	// Cache principal pour les résultats d'évaluation
	cache map[string]*EvalCacheEntry

	// Index par type de condition pour des recherches rapides
	conditionTypeIndex map[string][]string

	// Cache LRU pour éviction intelligente
	lruList *LRUList

	// Configuration du cache
	config CacheConfig

	// Statistiques de performance
	stats CacheStats

	// Verrou pour la concurrence
	mutex sync.RWMutex

	// Timer pour le nettoyage automatique
	cleanupTimer *time.Timer
}

// CacheConfig configure les paramètres du cache
type CacheConfig struct {
	// Taille maximale du cache
	MaxSize int

	// TTL par défaut pour les entrées
	DefaultTTL time.Duration

	// TTL spécifique par type de condition
	TypeSpecificTTL map[string]time.Duration

	// Intervalle de nettoyage automatique
	CleanupInterval time.Duration

	// Seuil pour le pré-calcul intelligent
	PrecomputeThreshold int

	// Activer la compression des clés
	EnableKeyCompression bool

	// Taille maximale des clés avant compression
	MaxKeyLength int
}

// CacheStats contient les statistiques du cache
type CacheStats struct {
	Hits            int64
	Misses          int64
	Evictions       int64
	Cleanups        int64
	Compressions    int64
	PrecomputeHits  int64
	AverageHitTime  time.Duration
	AverageMissTime time.Duration
}

// EvalCacheEntry représente une entrée dans le cache d'évaluation
type EvalCacheEntry struct {
	// Résultat de l'évaluation
	Result bool

	// Erreur associée (si any)
	Error error

	// Métadonnées de cache
	CreatedAt    time.Time
	LastAccessed time.Time
	AccessCount  int64
	TTL          time.Duration

	// Score de confiance du résultat
	Confidence float64

	// Coût de calcul (pour priorité d'éviction)
	ComputeCost time.Duration

	// Tags pour classification
	Tags []string
}

// LRUList implémente une liste doublement liée pour LRU
type LRUList struct {
	head *LRUNode
	tail *LRUNode
	size int
}

// LRUNode représente un nœud dans la liste LRU
type LRUNode struct {
	key  string
	prev *LRUNode
	next *LRUNode
}

// EvaluationKey représente une clé d'évaluation structurée
type EvaluationKey struct {
	ConditionType string
	FactType      string
	FieldName     string
	Operator      string
	Value         interface{}
	FactID        string
}

// NewEvaluationCache crée un nouveau cache d'évaluation intelligent
func NewEvaluationCache(config CacheConfig) *EvaluationCache {
	if config.MaxSize <= 0 {
		config.MaxSize = 10000
	}
	if config.DefaultTTL <= 0 {
		config.DefaultTTL = 5 * time.Minute
	}
	if config.CleanupInterval <= 0 {
		config.CleanupInterval = time.Minute
	}

	ec := &EvaluationCache{
		cache:              make(map[string]*EvalCacheEntry),
		conditionTypeIndex: make(map[string][]string),
		lruList:            NewLRUList(),
		config:             config,
		stats:              CacheStats{},
	}

	// Démarrer le nettoyage automatique
	ec.startCleanupTimer()

	return ec
}

// NewLRUList crée une nouvelle liste LRU
func NewLRUList() *LRUList {
	head := &LRUNode{}
	tail := &LRUNode{}
	head.next = tail
	tail.prev = head

	return &LRUList{
		head: head,
		tail: tail,
		size: 0,
	}
}

// Get récupère un résultat d'évaluation depuis le cache
func (ec *EvaluationCache) Get(key *EvaluationKey) (bool, error, bool) {
	startTime := time.Now()

	ec.mutex.RLock()
	cacheKey := ec.generateCacheKey(key)
	entry, exists := ec.cache[cacheKey]
	ec.mutex.RUnlock()

	if !exists {
		ec.stats.Misses++
		ec.stats.AverageMissTime = ec.updateAverageTime(ec.stats.AverageMissTime, time.Since(startTime), ec.stats.Misses)
		return false, nil, false
	}

	// Vérifier la TTL
	if time.Since(entry.CreatedAt) > entry.TTL {
		ec.mutex.Lock()
		ec.removeEntry(cacheKey)
		ec.mutex.Unlock()
		ec.stats.Misses++
		ec.stats.AverageMissTime = ec.updateAverageTime(ec.stats.AverageMissTime, time.Since(startTime), ec.stats.Misses)
		return false, nil, false
	}

	// Mettre à jour les statistiques d'accès
	ec.mutex.Lock()
	entry.LastAccessed = time.Now()
	entry.AccessCount++
	ec.lruList.MoveToFront(cacheKey)
	ec.mutex.Unlock()

	ec.stats.Hits++
	ec.stats.AverageHitTime = ec.updateAverageTime(ec.stats.AverageHitTime, time.Since(startTime), ec.stats.Hits)

	return entry.Result, entry.Error, true
}

// Put stocke un résultat d'évaluation dans le cache
func (ec *EvaluationCache) Put(key *EvaluationKey, result bool, err error, computeCost time.Duration) {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	cacheKey := ec.generateCacheKey(key)

	// Éviction si nécessaire
	if len(ec.cache) >= ec.config.MaxSize {
		ec.evictLRU()
	}

	// Déterminer la TTL appropriée
	ttl := ec.config.DefaultTTL
	if specificTTL, exists := ec.config.TypeSpecificTTL[key.ConditionType]; exists {
		ttl = specificTTL
	}

	// Calculer le score de confiance
	confidence := ec.calculateConfidence(key, computeCost)

	// Créer l'entrée
	entry := &EvalCacheEntry{
		Result:       result,
		Error:        err,
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
		AccessCount:  1,
		TTL:          ttl,
		Confidence:   confidence,
		ComputeCost:  computeCost,
		Tags:         ec.generateTags(key),
	}

	// Stocker dans le cache
	ec.cache[cacheKey] = entry

	// Mettre à jour l'index par type
	if ec.conditionTypeIndex[key.ConditionType] == nil {
		ec.conditionTypeIndex[key.ConditionType] = make([]string, 0)
	}
	ec.conditionTypeIndex[key.ConditionType] = append(ec.conditionTypeIndex[key.ConditionType], cacheKey)

	// Ajouter à la liste LRU
	ec.lruList.AddToFront(cacheKey)
}

// generateCacheKey génère une clé de cache optimisée
func (ec *EvaluationCache) generateCacheKey(key *EvaluationKey) string {
	// Créer une clé structurée
	baseKey := fmt.Sprintf("%s:%s:%s:%s:%v:%s",
		key.ConditionType,
		key.FactType,
		key.FieldName,
		key.Operator,
		key.Value,
		key.FactID)

	// Compresser la clé si elle est trop longue
	if ec.config.EnableKeyCompression && len(baseKey) > ec.config.MaxKeyLength {
		hash := sha256.Sum256([]byte(baseKey))
		compressedKey := fmt.Sprintf("compressed:%x", hash)
		ec.stats.Compressions++
		return compressedKey
	}

	return baseKey
}

// calculateConfidence calcule un score de confiance pour l'entrée
func (ec *EvaluationCache) calculateConfidence(key *EvaluationKey, computeCost time.Duration) float64 {
	confidence := 1.0

	// Réduire la confiance pour des calculs très rapides (possiblement triviaux)
	if computeCost < time.Microsecond {
		confidence *= 0.8
	}

	// Augmenter la confiance pour des calculs coûteux
	if computeCost > time.Millisecond {
		confidence *= 1.2
	}

	// Ajuster selon le type de condition
	switch key.ConditionType {
	case "binary_operation":
		confidence *= 0.95 // Généralement stable
	case "logical_expression":
		confidence *= 0.85 // Plus complexe, moins stable
	case "function_call":
		confidence *= 0.75 // Peut dépendre d'un état externe
	}

	// Limiter à [0, 1]
	if confidence > 1.0 {
		confidence = 1.0
	}
	if confidence < 0.0 {
		confidence = 0.0
	}

	return confidence
}

// generateTags génère des tags pour l'entrée
func (ec *EvaluationCache) generateTags(key *EvaluationKey) []string {
	tags := []string{
		"type:" + key.ConditionType,
		"fact:" + key.FactType,
		"field:" + key.FieldName,
		"op:" + key.Operator,
	}

	// Ajouter des tags spécifiques à la valeur
	switch v := key.Value.(type) {
	case string:
		if len(v) < 10 {
			tags = append(tags, "short_string")
		} else {
			tags = append(tags, "long_string")
		}
	case int, int32, int64, float32, float64:
		tags = append(tags, "numeric")
	case bool:
		tags = append(tags, "boolean")
	default:
		tags = append(tags, "complex_type")
	}

	return tags
}

// evictLRU évince l'entrée la moins récemment utilisée
func (ec *EvaluationCache) evictLRU() {
	if ec.lruList.size == 0 {
		return
	}

	// Récupérer la clé de l'élément le moins récemment utilisé
	lruKey := ec.lruList.RemoveFromTail()
	if lruKey == "" {
		return
	}

	// Supprimer du cache
	ec.removeEntry(lruKey)
	ec.stats.Evictions++
}

// removeEntry supprime une entrée du cache et des index
func (ec *EvaluationCache) removeEntry(cacheKey string) {
	_, exists := ec.cache[cacheKey]
	if !exists {
		return
	}

	// Supprimer du cache principal
	delete(ec.cache, cacheKey)

	// Supprimer des index par type
	for condType, keys := range ec.conditionTypeIndex {
		for i, key := range keys {
			if key == cacheKey {
				// Supprimer de la slice
				ec.conditionTypeIndex[condType] = append(keys[:i], keys[i+1:]...)
				break
			}
		}
	}

	// Supprimer de la liste LRU
	ec.lruList.Remove(cacheKey)
}

// CleanupExpired supprime les entrées expirées
func (ec *EvaluationCache) CleanupExpired() int {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	now := time.Now()
	expired := make([]string, 0)

	// Identifier les entrées expirées
	for key, entry := range ec.cache {
		if now.Sub(entry.CreatedAt) > entry.TTL {
			expired = append(expired, key)
		}
	}

	// Supprimer les entrées expirées
	for _, key := range expired {
		ec.removeEntry(key)
	}

	ec.stats.Cleanups++
	return len(expired)
}

// PrecomputeFrequent pré-calcule les évaluations fréquentes
func (ec *EvaluationCache) PrecomputeFrequent(evaluator func(*EvaluationKey) (bool, error, time.Duration)) {
	ec.mutex.RLock()
	candidates := make([]*EvaluationKey, 0)

	// Identifier les patterns fréquents
	for _, entry := range ec.cache {
		if entry.AccessCount >= int64(ec.config.PrecomputeThreshold) {
			// Reconstituer la clé (simplifiée pour cet exemple)
			// Dans une implémentation complète, nous stockerions la clé originale
		}
	}
	ec.mutex.RUnlock()

	// Pré-calculer les résultats pour ces patterns
	for _, key := range candidates {
		result, err, computeTime := evaluator(key)
		ec.Put(key, result, err, computeTime)
		ec.stats.PrecomputeHits++
	}
}

// InvalidateByType invalide toutes les entrées d'un type donné
func (ec *EvaluationCache) InvalidateByType(conditionType string) int {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	keys, exists := ec.conditionTypeIndex[conditionType]
	if !exists {
		return 0
	}

	count := len(keys)

	// Supprimer toutes les entrées de ce type
	for _, key := range keys {
		ec.removeEntry(key)
	}

	// Vider l'index pour ce type
	delete(ec.conditionTypeIndex, conditionType)

	return count
}

// InvalidateByFact invalide les entrées associées à un fait spécifique
func (ec *EvaluationCache) InvalidateByFact(factID string) int {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	toRemove := make([]string, 0)

	// Trouver toutes les entrées associées à ce fait
	for key, _ := range ec.cache {
		if ec.containsFactID(key, factID) {
			toRemove = append(toRemove, key)
		}
	}

	// Supprimer les entrées trouvées
	for _, key := range toRemove {
		ec.removeEntry(key)
	}

	return len(toRemove)
}

// containsFactID vérifie si une clé contient un ID de fait spécifique
func (ec *EvaluationCache) containsFactID(cacheKey, factID string) bool {
	// Simplifiée - dans une vraie implémentation, nous analysions la structure de la clé
	return fmt.Sprintf(":%s", factID) == cacheKey[len(cacheKey)-len(factID)-1:]
}

// updateAverageTime met à jour une moyenne de temps de façon efficace
func (ec *EvaluationCache) updateAverageTime(currentAvg time.Duration, newTime time.Duration, count int64) time.Duration {
	if count <= 1 {
		return newTime
	}

	// Moyenne pondérée pour éviter les débordements
	return time.Duration(
		(int64(currentAvg)*(count-1) + int64(newTime)) / count,
	)
}

// startCleanupTimer démarre le timer de nettoyage automatique
func (ec *EvaluationCache) startCleanupTimer() {
	ec.cleanupTimer = time.AfterFunc(ec.config.CleanupInterval, func() {
		ec.CleanupExpired()
		ec.startCleanupTimer() // Redémarrer le timer
	})
}

// GetStats retourne les statistiques du cache
func (ec *EvaluationCache) GetStats() CacheStats {
	ec.mutex.RLock()
	defer ec.mutex.RUnlock()

	return ec.stats
}

// Clear vide complètement le cache
func (ec *EvaluationCache) Clear() {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	ec.cache = make(map[string]*EvalCacheEntry)
	ec.conditionTypeIndex = make(map[string][]string)
	ec.lruList = NewLRUList()
	ec.stats = CacheStats{}
}

// Méthodes pour la liste LRU

// AddToFront ajoute une clé au début de la liste LRU
func (lru *LRUList) AddToFront(key string) {
	node := &LRUNode{key: key}

	node.next = lru.head.next
	node.prev = lru.head
	lru.head.next.prev = node
	lru.head.next = node

	lru.size++
}

// MoveToFront déplace une clé vers le début de la liste
func (lru *LRUList) MoveToFront(key string) {
	// Trouver le nœud
	current := lru.head.next
	for current != lru.tail {
		if current.key == key {
			// Supprimer de la position actuelle
			current.prev.next = current.next
			current.next.prev = current.prev

			// Ajouter au début
			current.next = lru.head.next
			current.prev = lru.head
			lru.head.next.prev = current
			lru.head.next = current

			return
		}
		current = current.next
	}
}

// RemoveFromTail supprime et retourne la clé de la fin de la liste
func (lru *LRUList) RemoveFromTail() string {
	if lru.size == 0 {
		return ""
	}

	lastNode := lru.tail.prev
	if lastNode == lru.head {
		return ""
	}

	key := lastNode.key

	lastNode.prev.next = lru.tail
	lru.tail.prev = lastNode.prev

	lru.size--

	return key
}

// Remove supprime une clé spécifique de la liste
func (lru *LRUList) Remove(key string) {
	current := lru.head.next
	for current != lru.tail {
		if current.key == key {
			current.prev.next = current.next
			current.next.prev = current.prev
			lru.size--
			return
		}
		current = current.next
	}
}
