// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"sync"
	"time"
)

// DeltaDetector détecte les changements entre deux versions d'un fait.
//
// Il compare un fait "avant" et "après" modification et génère un FactDelta
// contenant uniquement les champs qui ont changé.
//
// Thread-safety : DeltaDetector est safe pour utilisation concurrent.
type DeltaDetector struct {
	config DetectorConfig

	// Cache de comparaisons (si activé)
	cache      map[string]*detectorCacheEntry
	cacheMutex sync.RWMutex

	// Métriques
	comparisons  int64
	cacheHits    int64
	cacheMisses  int64
	metricsMutex sync.RWMutex
}

// detectorCacheEntry représente une entrée dans le cache de comparaisons.
type detectorCacheEntry struct {
	delta     *FactDelta
	createdAt time.Time
}

// NewDeltaDetector crée un nouveau détecteur avec configuration par défaut.
func NewDeltaDetector() *DeltaDetector {
	return NewDeltaDetectorWithConfig(DefaultDetectorConfig())
}

// NewDeltaDetectorWithConfig crée un détecteur avec une configuration spécifique.
func NewDeltaDetectorWithConfig(config DetectorConfig) *DeltaDetector {
	if err := config.Validate(); err != nil {
		config = DefaultDetectorConfig()
	}

	dd := &DeltaDetector{
		config: config,
	}

	if config.CacheComparisons {
		dd.cache = make(map[string]*detectorCacheEntry)
	}

	return dd
}

// DetectDelta compare deux faits et retourne les changements détectés.
//
// Paramètres :
//   - oldFact : fait avant modification (map[string]interface{})
//   - newFact : fait après modification (map[string]interface{})
//   - factID : identifiant interne du fait (ex: "Product~123")
//   - factType : type du fait (ex: "Product")
//
// Retourne :
//   - *FactDelta contenant les champs modifiés
//   - error si la détection échoue
//
// Le FactDelta retourné peut être vide (IsEmpty() == true) si aucun
// changement n'est détecté.
//
// IMPORTANT: L'appelant doit libérer le FactDelta avec ReleaseFactDelta()
// une fois l'utilisation terminée (sauf si mis en cache ou stocké).
func (dd *DeltaDetector) DetectDelta(
	oldFact, newFact map[string]interface{},
	factID, factType string,
) (*FactDelta, error) {
	dd.incrementComparisons()

	// Vérifier cache si activé
	if dd.config.CacheComparisons {
		cacheKey := dd.buildCacheKey(oldFact, newFact, factID)
		if cached := dd.getFromCache(cacheKey); cached != nil {
			dd.incrementCacheHits()
			return cached, nil
		}
		dd.incrementCacheMisses()
	}

	// Acquérir depuis le pool au lieu de NewFactDelta
	delta := AcquireFactDelta(factID, factType)
	delta.FieldCount = len(newFact)

	// Collecter tous les noms de champs (union des deux faits)
	allFields := make(map[string]bool)
	for field := range oldFact {
		allFields[field] = true
	}
	for field := range newFact {
		allFields[field] = true
	}

	// Comparer chaque champ
	for fieldName := range allFields {
		// Ignorer si configuré
		if dd.config.ShouldIgnoreField(fieldName) {
			continue
		}

		oldValue, oldExists := oldFact[fieldName]
		newValue, newExists := newFact[fieldName]

		// Cas 1 : Champ ajouté
		if !oldExists && newExists {
			delta.AddFieldChange(fieldName, nil, newValue)
			continue
		}

		// Cas 2 : Champ supprimé
		if oldExists && !newExists {
			delta.AddFieldChange(fieldName, oldValue, nil)
			continue
		}

		// Cas 3 : Champ modifié (existe dans les deux)
		if oldExists && newExists {
			if !dd.valuesEqual(oldValue, newValue, 0) {
				delta.AddFieldChange(fieldName, oldValue, newValue)
			}
		}
	}

	// Ajouter au cache si activé et pas vide
	if dd.config.CacheComparisons && !delta.IsEmpty() {
		cacheKey := dd.buildCacheKey(oldFact, newFact, factID)
		dd.addToCache(cacheKey, delta)
	}

	// Note: L'appelant doit appeler ReleaseFactDelta(delta) quand fini
	// SAUF si le delta est mis en cache ou stocké ailleurs

	return delta, nil
}

// DetectDeltaQuick est une version optimisée qui ne crée pas de FactDelta
// si aucun changement n'est détecté.
//
// Cette fonction est plus rapide que DetectDelta pour les cas no-op
// (faits identiques).
//
// Retourne :
//   - *FactDelta (nil si aucun changement)
//   - error si échec
func (dd *DeltaDetector) DetectDeltaQuick(
	oldFact, newFact map[string]interface{},
	factID, factType string,
) (*FactDelta, error) {
	dd.incrementComparisons()

	// Quick check : tailles différentes = forcément des changements
	if len(oldFact) != len(newFact) {
		return dd.DetectDelta(oldFact, newFact, factID, factType)
	}

	// Quick check : tous les champs identiques ?
	hasChanges := false
	for fieldName, oldValue := range oldFact {
		if dd.config.ShouldIgnoreField(fieldName) {
			continue
		}

		newValue, exists := newFact[fieldName]
		if !exists || !dd.valuesEqual(oldValue, newValue, 0) {
			hasChanges = true
			break
		}
	}

	if !hasChanges {
		return nil, nil
	}

	return dd.DetectDelta(oldFact, newFact, factID, factType)
}

// valuesEqual compare deux valeurs en respectant la configuration.
//
// Paramètres :
//   - a, b : valeurs à comparer
//   - depth : profondeur de récursion actuelle (pour structures imbriquées)
//
// Retourne true si les valeurs sont égales selon la configuration.
func (dd *DeltaDetector) valuesEqual(a, b interface{}, depth int) bool {
	// Protection contre récursion infinie
	if depth > dd.config.MaxNestingLevel {
		return true
	}

	// Fast path: utiliser version optimisée pour types simples
	// Cela évite reflect.TypeOf pour les types communs
	if depth == 0 {
		// Au premier niveau, utiliser la version optimisée
		return OptimizedValuesEqual(a, b, dd.config.FloatEpsilon)
	}

	// Cas 1 : nil vs non-nil
	if a == nil || b == nil {
		return a == b
	}

	// Cas 2 : Track type changes
	if dd.config.TrackTypeChanges {
		typeA := inferValueType(a)
		typeB := inferValueType(b)
		if typeA != typeB {
			return false
		}
	}

	// Cas 3 : Comparaison profonde si activée
	if dd.config.EnableDeepComparison {
		// Maps imbriquées
		if mapA, okA := a.(map[string]interface{}); okA {
			if mapB, okB := b.(map[string]interface{}); okB {
				return dd.mapsEqual(mapA, mapB, depth+1)
			}
		}

		// Slices imbriquées
		if sliceA, okA := a.([]interface{}); okA {
			if sliceB, okB := b.([]interface{}); okB {
				return dd.slicesEqual(sliceA, sliceB, depth+1)
			}
		}
	}

	return OptimizedValuesEqual(a, b, dd.config.FloatEpsilon)
}

// mapsEqual compare deux maps récursivement.
func (dd *DeltaDetector) mapsEqual(a, b map[string]interface{}, depth int) bool {
	if len(a) != len(b) {
		return false
	}

	for key, valA := range a {
		valB, exists := b[key]
		if !exists {
			return false
		}
		if !dd.valuesEqual(valA, valB, depth) {
			return false
		}
	}

	return true
}

// slicesEqual compare deux slices récursivement.
func (dd *DeltaDetector) slicesEqual(a, b []interface{}, depth int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !dd.valuesEqual(a[i], b[i], depth) {
			return false
		}
	}

	return true
}

// buildCacheKey génère une clé de cache pour une paire de faits.
//
// Note : Cette implémentation simple utilise le factID + timestamp.
// Une version optimisée pourrait hacher les valeurs des faits.
func (dd *DeltaDetector) buildCacheKey(oldFact, newFact map[string]interface{}, factID string) string {
	return factID
}

// getFromCache récupère un delta depuis le cache.
func (dd *DeltaDetector) getFromCache(key string) *FactDelta {
	dd.cacheMutex.RLock()
	defer dd.cacheMutex.RUnlock()

	entry, exists := dd.cache[key]
	if !exists {
		return nil
	}

	// Vérifier expiration
	if time.Since(entry.createdAt) > dd.config.CacheTTL {
		go dd.removeFromCache(key)
		return nil
	}

	return entry.delta
}

// addToCache ajoute un delta au cache.
func (dd *DeltaDetector) addToCache(key string, delta *FactDelta) {
	dd.cacheMutex.Lock()
	defer dd.cacheMutex.Unlock()

	dd.cache[key] = &detectorCacheEntry{
		delta:     delta,
		createdAt: time.Now(),
	}
}

// removeFromCache supprime une entrée du cache.
func (dd *DeltaDetector) removeFromCache(key string) {
	dd.cacheMutex.Lock()
	defer dd.cacheMutex.Unlock()

	delete(dd.cache, key)
}

// ClearCache vide complètement le cache.
func (dd *DeltaDetector) ClearCache() {
	dd.cacheMutex.Lock()
	defer dd.cacheMutex.Unlock()

	dd.cache = make(map[string]*detectorCacheEntry)
}

// DetectorMetrics contient les métriques du détecteur.
type DetectorMetrics struct {
	Comparisons int64
	CacheHits   int64
	CacheMisses int64
	CacheSize   int
	HitRate     float64
}

// GetMetrics retourne les métriques du détecteur.
func (dd *DeltaDetector) GetMetrics() DetectorMetrics {
	dd.metricsMutex.RLock()
	defer dd.metricsMutex.RUnlock()

	dd.cacheMutex.RLock()
	cacheSize := len(dd.cache)
	dd.cacheMutex.RUnlock()

	hitRate := 0.0
	totalCacheAccess := dd.cacheHits + dd.cacheMisses
	if totalCacheAccess > 0 {
		hitRate = float64(dd.cacheHits) / float64(totalCacheAccess)
	}

	return DetectorMetrics{
		Comparisons: dd.comparisons,
		CacheHits:   dd.cacheHits,
		CacheMisses: dd.cacheMisses,
		CacheSize:   cacheSize,
		HitRate:     hitRate,
	}
}

// ResetMetrics réinitialise les métriques à zéro.
func (dd *DeltaDetector) ResetMetrics() {
	dd.metricsMutex.Lock()
	defer dd.metricsMutex.Unlock()

	dd.comparisons = 0
	dd.cacheHits = 0
	dd.cacheMisses = 0
}

// incrementComparisons incrémente le compteur de comparaisons.
func (dd *DeltaDetector) incrementComparisons() {
	dd.metricsMutex.Lock()
	dd.comparisons++
	dd.metricsMutex.Unlock()
}

// incrementCacheHits incrémente le compteur de cache hits.
func (dd *DeltaDetector) incrementCacheHits() {
	dd.metricsMutex.Lock()
	dd.cacheHits++
	dd.metricsMutex.Unlock()
}

// incrementCacheMisses incrémente le compteur de cache misses.
func (dd *DeltaDetector) incrementCacheMisses() {
	dd.metricsMutex.Lock()
	dd.cacheMisses++
	dd.metricsMutex.Unlock()
}

// GetConfig retourne une copie de la configuration actuelle.
func (dd *DeltaDetector) GetConfig() DetectorConfig {
	return dd.config.Clone()
}
