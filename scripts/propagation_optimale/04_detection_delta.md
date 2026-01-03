# 🔎 Prompt 04 - Détection de Delta

> **📋 Standards** : Ce prompt respecte les règles de `.github/prompts/common.md` et `.github/prompts/develop.md`

## 🎯 Objectif

Implémenter le système de détection de changements entre deux versions d'un fait : `DeltaDetector` qui compare un fait avant et après modification et génère un `FactDelta` précis.

Cette détection est cruciale pour la propagation optimale : elle identifie exactement quels champs ont changé, permettant ainsi de ne propager que vers les nœuds concernés.

**⚠️ IMPORTANT** : Ce prompt génère du code. Respecter strictement les standards de `common.md`.

---

## 📋 Prérequis

Avant de commencer ce prompt :

- [x] **Prompt 01 validé** : Conception disponible
- [x] **Prompt 02 validé** : Modèle de données delta implémenté
- [x] **Prompt 03 validé** : Indexation des dépendances implémentée
- [x] **Tests passent** : `go test ./rete/delta/... -v` (100% success)
- [x] **Documents de référence** :
  - `REPORTS/conception_delta_architecture.md`
  - `rete/delta/field_delta.go` - Structures FieldDelta, FactDelta
  - `rete/delta/comparison.go` - Comparaison de valeurs

---

## 📂 Fichiers à Créer

Ajouter au package `rete/delta` :

```
rete/delta/
├── delta_detector.go           # Détecteur de changements
├── delta_detector_test.go      # Tests unitaires
├── detector_config.go          # Configuration du détecteur
├── detector_config_test.go     # Tests configuration
└── detector_benchmark_test.go  # Benchmarks performance
```

---

## 🔧 Tâche 1 : Configuration du Détecteur

### Fichier : `rete/delta/detector_config.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"time"
)

// DetectorConfig contient la configuration du DeltaDetector.
//
// Cette configuration permet d'ajuster le comportement de la détection
// selon les besoins de performance et de précision.
type DetectorConfig struct {
	// FloatEpsilon est la tolérance pour la comparaison de floats.
	// Valeurs différant de moins que epsilon sont considérées égales.
	// Default: DefaultFloatEpsilon (1e-9)
	FloatEpsilon float64
	
	// IgnoreInternalFields indique si les champs internes (préfixe "_")
	// doivent être ignorés lors de la détection.
	// Default: true
	IgnoreInternalFields bool
	
	// IgnoredFields est une liste de noms de champs à ignorer
	// lors de la détection (ex: timestamps auto-générés).
	// Default: []
	IgnoredFields []string
	
	// TrackTypeChanges indique si un changement de type de valeur
	// doit être détecté (ex: 42 (int) → "42" (string)).
	// Default: true
	TrackTypeChanges bool
	
	// EnableDeepComparison active la comparaison profonde pour
	// les structures imbriquées (maps, slices).
	// Default: true
	EnableDeepComparison bool
	
	// MaxNestingLevel est la profondeur maximale pour la comparaison
	// récursive des structures imbriquées (protection stack overflow).
	// Default: 10
	MaxNestingLevel int
	
	// CacheComparisons active le cache des résultats de comparaison
	// pour optimiser les comparaisons répétées.
	// Default: false (car overhead mémoire)
	CacheComparisons bool
	
	// CacheTTL est la durée de vie des entrées du cache.
	// Default: 1 minute
	CacheTTL time.Duration
}

// DefaultDetectorConfig retourne une configuration par défaut.
func DefaultDetectorConfig() DetectorConfig {
	return DetectorConfig{
		FloatEpsilon:         DefaultFloatEpsilon,
		IgnoreInternalFields: true,
		IgnoredFields:        []string{},
		TrackTypeChanges:     true,
		EnableDeepComparison: true,
		MaxNestingLevel:      10,
		CacheComparisons:     false,
		CacheTTL:             1 * time.Minute,
	}
}

// Validate vérifie que la configuration est valide.
//
// Retourne une erreur si des paramètres sont incohérents.
func (dc *DetectorConfig) Validate() error {
	if dc.FloatEpsilon < 0 {
		return ErrInvalidConfig("FloatEpsilon must be >= 0")
	}
	
	if dc.MaxNestingLevel < 1 {
		return ErrInvalidConfig("MaxNestingLevel must be >= 1")
	}
	
	if dc.CacheTTL < 0 {
		return ErrInvalidConfig("CacheTTL must be >= 0")
	}
	
	return nil
}

// ShouldIgnoreField retourne true si un champ doit être ignoré.
//
// Un champ est ignoré si :
//   - Il est dans IgnoredFields
//   - Il commence par "_" et IgnoreInternalFields est true
func (dc *DetectorConfig) ShouldIgnoreField(fieldName string) bool {
	// Vérifier liste explicite
	for _, ignored := range dc.IgnoredFields {
		if fieldName == ignored {
			return true
		}
	}
	
	// Vérifier champs internes
	if dc.IgnoreInternalFields && len(fieldName) > 0 && fieldName[0] == '_' {
		return true
	}
	
	return false
}

// Clone crée une copie de la configuration.
func (dc *DetectorConfig) Clone() DetectorConfig {
	ignoredFields := make([]string, len(dc.IgnoredFields))
	copy(ignoredFields, dc.IgnoredFields)
	
	return DetectorConfig{
		FloatEpsilon:         dc.FloatEpsilon,
		IgnoreInternalFields: dc.IgnoreInternalFields,
		IgnoredFields:        ignoredFields,
		TrackTypeChanges:     dc.TrackTypeChanges,
		EnableDeepComparison: dc.EnableDeepComparison,
		MaxNestingLevel:      dc.MaxNestingLevel,
		CacheComparisons:     dc.CacheComparisons,
		CacheTTL:             dc.CacheTTL,
	}
}

// ErrInvalidConfig représente une erreur de configuration invalide.
type ErrInvalidConfig string

func (e ErrInvalidConfig) Error() string {
	return "invalid detector config: " + string(e)
}
```

### Tests : `rete/delta/detector_config_test.go`

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
	"time"
)

func TestDefaultDetectorConfig(t *testing.T) {
	config := DefaultDetectorConfig()
	
	if config.FloatEpsilon != DefaultFloatEpsilon {
		t.Errorf("Expected FloatEpsilon = %v, got %v", DefaultFloatEpsilon, config.FloatEpsilon)
	}
	
	if !config.IgnoreInternalFields {
		t.Error("Expected IgnoreInternalFields = true")
	}
	
	if !config.TrackTypeChanges {
		t.Error("Expected TrackTypeChanges = true")
	}
	
	if !config.EnableDeepComparison {
		t.Error("Expected EnableDeepComparison = true")
	}
	
	if config.MaxNestingLevel != 10 {
		t.Errorf("Expected MaxNestingLevel = 10, got %d", config.MaxNestingLevel)
	}
	
	if config.CacheComparisons {
		t.Error("Expected CacheComparisons = false by default")
	}
}

func TestDetectorConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    DetectorConfig
		wantError bool
	}{
		{
			name:      "valid default config",
			config:    DefaultDetectorConfig(),
			wantError: false,
		},
		{
			name: "negative epsilon",
			config: DetectorConfig{
				FloatEpsilon:    -0.1,
				MaxNestingLevel: 10,
				CacheTTL:        time.Minute,
			},
			wantError: true,
		},
		{
			name: "zero nesting level",
			config: DetectorConfig{
				FloatEpsilon:    DefaultFloatEpsilon,
				MaxNestingLevel: 0,
				CacheTTL:        time.Minute,
			},
			wantError: true,
		},
		{
			name: "negative cache TTL",
			config: DetectorConfig{
				FloatEpsilon:    DefaultFloatEpsilon,
				MaxNestingLevel: 10,
				CacheTTL:        -time.Second,
			},
			wantError: true,
		},
		{
			name: "edge case: zero epsilon",
			config: DetectorConfig{
				FloatEpsilon:    0,
				MaxNestingLevel: 1,
				CacheTTL:        0,
			},
			wantError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestDetectorConfig_ShouldIgnoreField(t *testing.T) {
	tests := []struct {
		name      string
		config    DetectorConfig
		fieldName string
		want      bool
	}{
		{
			name:      "normal field",
			config:    DefaultDetectorConfig(),
			fieldName: "price",
			want:      false,
		},
		{
			name:      "internal field with underscore",
			config:    DefaultDetectorConfig(),
			fieldName: "_internal",
			want:      true,
		},
		{
			name: "internal field but not ignored",
			config: DetectorConfig{
				IgnoreInternalFields: false,
			},
			fieldName: "_internal",
			want:      false,
		},
		{
			name: "explicitly ignored field",
			config: DetectorConfig{
				IgnoredFields: []string{"timestamp", "updated_at"},
			},
			fieldName: "timestamp",
			want:      true,
		},
		{
			name: "not in ignored list",
			config: DetectorConfig{
				IgnoredFields: []string{"timestamp"},
			},
			fieldName: "price",
			want:      false,
		},
		{
			name: "multiple ignored fields",
			config: DetectorConfig{
				IgnoredFields: []string{"field1", "field2", "field3"},
			},
			fieldName: "field2",
			want:      true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.ShouldIgnoreField(tt.fieldName)
			if got != tt.want {
				t.Errorf("ShouldIgnoreField(%s) = %v, want %v", tt.fieldName, got, tt.want)
			}
		})
	}
}

func TestDetectorConfig_Clone(t *testing.T) {
	original := DetectorConfig{
		FloatEpsilon:         0.001,
		IgnoreInternalFields: true,
		IgnoredFields:        []string{"field1", "field2"},
		TrackTypeChanges:     true,
		EnableDeepComparison: true,
		MaxNestingLevel:      5,
		CacheComparisons:     true,
		CacheTTL:             2 * time.Minute,
	}
	
	cloned := original.Clone()
	
	// Vérifier égalité des valeurs
	if cloned.FloatEpsilon != original.FloatEpsilon {
		t.Error("FloatEpsilon not cloned correctly")
	}
	if cloned.IgnoreInternalFields != original.IgnoreInternalFields {
		t.Error("IgnoreInternalFields not cloned correctly")
	}
	if cloned.MaxNestingLevel != original.MaxNestingLevel {
		t.Error("MaxNestingLevel not cloned correctly")
	}
	
	// Vérifier que les slices sont des copies indépendantes
	if len(cloned.IgnoredFields) != len(original.IgnoredFields) {
		t.Error("IgnoredFields length mismatch")
	}
	
	// Modifier le clone ne doit pas affecter l'original
	cloned.IgnoredFields[0] = "modified"
	if original.IgnoredFields[0] == "modified" {
		t.Error("Clone is not independent (slice mutation affected original)")
	}
	
	cloned.FloatEpsilon = 999
	if original.FloatEpsilon == 999 {
		t.Error("Clone is not independent (field mutation affected original)")
	}
}

func TestErrInvalidConfig_Error(t *testing.T) {
	err := ErrInvalidConfig("test message")
	expected := "invalid detector config: test message"
	
	if err.Error() != expected {
		t.Errorf("Error() = %s, want %s", err.Error(), expected)
	}
}
```

---

## 🔧 Tâche 2 : Implémentation du DeltaDetector

### Fichier : `rete/delta/delta_detector.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
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
	cache      map[string]*cacheEntry
	cacheMutex sync.RWMutex
	
	// Métriques
	comparisons     int64
	cacheHits       int64
	cacheMisses     int64
	metricsMutex    sync.RWMutex
}

// cacheEntry représente une entrée dans le cache de comparaisons.
type cacheEntry struct {
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
		// Fallback sur config par défaut si invalide
		config = DefaultDetectorConfig()
	}
	
	dd := &DeltaDetector{
		config: config,
	}
	
	if config.CacheComparisons {
		dd.cache = make(map[string]*cacheEntry)
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
	
	// Créer le FactDelta
	delta := NewFactDelta(factID, factType)
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
	
	// Ajouter au cache si activé
	if dd.config.CacheComparisons && !delta.IsEmpty() {
		cacheKey := dd.buildCacheKey(oldFact, newFact, factID)
		dd.addToCache(cacheKey, delta)
	}
	
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
		// Aucun changement détecté
		return nil, nil
	}
	
	// Des changements existent, faire détection complète
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
		// Au-delà de la limite, considérer égaux par défaut
		return true
	}
	
	// Cas 1 : Valeurs identiques (référence ou valeur simple)
	if a == b {
		return true
	}
	
	// Cas 2 : nil vs non-nil
	if a == nil || b == nil {
		return a == b
	}
	
	// Cas 3 : Track type changes
	if dd.config.TrackTypeChanges {
		typeA := inferValueType(a)
		typeB := inferValueType(b)
		if typeA != typeB {
			return false
		}
	}
	
	// Cas 4 : Comparaison profonde si activée
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
	
	// Cas 5 : Comparaison standard avec epsilon pour floats
	return ValuesEqual(a, b, dd.config.FloatEpsilon)
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
	// Implémentation simple : factID suffit car on suppose
	// qu'un fait donné ne change pas fréquemment
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
		// Expiré, supprimer (en async pour ne pas bloquer)
		go dd.removeFromCache(key)
		return nil
	}
	
	return entry.delta
}

// addToCache ajoute un delta au cache.
func (dd *DeltaDetector) addToCache(key string, delta *FactDelta) {
	dd.cacheMutex.Lock()
	defer dd.cacheMutex.Unlock()
	
	dd.cache[key] = &cacheEntry{
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
	
	dd.cache = make(map[string]*cacheEntry)
}

// GetMetrics retourne les métriques du détecteur.
type DetectorMetrics struct {
	Comparisons int64
	CacheHits   int64
	CacheMisses int64
	CacheSize   int
	HitRate     float64
}

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
```

### Tests : `rete/delta/delta_detector_test.go`

**Contenu** (partie 1/3) :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
	"time"
)

func TestNewDeltaDetector(t *testing.T) {
	detector := NewDeltaDetector()
	
	if detector == nil {
		t.Fatal("NewDeltaDetector returned nil")
	}
	
	config := detector.GetConfig()
	if config.FloatEpsilon != DefaultFloatEpsilon {
		t.Error("Default config not applied")
	}
}

func TestNewDeltaDetectorWithConfig(t *testing.T) {
	customConfig := DetectorConfig{
		FloatEpsilon:         0.001,
		IgnoreInternalFields: false,
		MaxNestingLevel:      5,
	}
	
	detector := NewDeltaDetectorWithConfig(customConfig)
	config := detector.GetConfig()
	
	if config.FloatEpsilon != 0.001 {
		t.Errorf("Custom FloatEpsilon not set: got %v", config.FloatEpsilon)
	}
	
	if config.IgnoreInternalFields {
		t.Error("Custom IgnoreInternalFields not set")
	}
}

func TestDeltaDetector_DetectDelta_NoChanges(t *testing.T) {
	detector := NewDeltaDetector()
	
	fact := map[string]interface{}{
		"id":     "123",
		"price":  100.0,
		"status": "active",
	}
	
	delta, err := detector.DetectDelta(fact, fact, "Product~123", "Product")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if !delta.IsEmpty() {
		t.Errorf("Expected empty delta for identical facts, got %d changes", len(delta.Fields))
	}
}

func TestDeltaDetector_DetectDelta_SingleFieldChange(t *testing.T) {
	detector := NewDeltaDetector()
	
	oldFact := map[string]interface{}{
		"id":     "123",
		"price":  100.0,
		"status": "active",
	}
	
	newFact := map[string]interface{}{
		"id":     "123",
		"price":  150.0,
		"status": "active",
	}
	
	delta, err := detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if delta.IsEmpty() {
		t.Fatal("Expected changes to be detected")
	}
	
	if len(delta.Fields) != 1 {
		t.Errorf("Expected 1 changed field, got %d", len(delta.Fields))
	}
	
	priceChange, exists := delta.Fields["price"]
	if !exists {
		t.Fatal("Expected 'price' field to be in delta")
	}
	
	if priceChange.OldValue != 100.0 {
		t.Errorf("Expected old price = 100.0, got %v", priceChange.OldValue)
	}
	
	if priceChange.NewValue != 150.0 {
		t.Errorf("Expected new price = 150.0, got %v", priceChange.NewValue)
	}
	
	if priceChange.ChangeType != ChangeTypeModified {
		t.Errorf("Expected ChangeTypeModified, got %v", priceChange.ChangeType)
	}
}

func TestDeltaDetector_DetectDelta_MultipleFieldChanges(t *testing.T) {
	detector := NewDeltaDetector()
	
	oldFact := map[string]interface{}{
		"id":       "123",
		"price":    100.0,
		"status":   "active",
		"quantity": 10,
	}
	
	newFact := map[string]interface{}{
		"id":       "123",
		"price":    150.0,
		"status":   "inactive",
		"quantity": 10,
	}
	
	delta, err := detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if len(delta.Fields) != 2 {
		t.Errorf("Expected 2 changed fields, got %d", len(delta.Fields))
	}
	
	// Vérifier que price et status ont changé
	if _, exists := delta.Fields["price"]; !exists {
		t.Error("Expected 'price' in delta")
	}
	
	if _, exists := delta.Fields["status"]; !exists {
		t.Error("Expected 'status' in delta")
	}
	
	// Vérifier que quantity n'a PAS changé
	if _, exists := delta.Fields["quantity"]; exists {
		t.Error("Did not expect 'quantity' in delta")
	}
}

func TestDeltaDetector_DetectDelta_FieldAdded(t *testing.T) {
	detector := NewDeltaDetector()
	
	oldFact := map[string]interface{}{
		"id":    "123",
		"price": 100.0,
	}
	
	newFact := map[string]interface{}{
		"id":       "123",
		"price":    100.0,
		"category": "Electronics",
	}
	
	delta, err := detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if len(delta.Fields) != 1 {
		t.Errorf("Expected 1 change (added field), got %d", len(delta.Fields))
	}
	
	categoryChange, exists := delta.Fields["category"]
	if !exists {
		t.Fatal("Expected 'category' in delta")
	}
	
	if categoryChange.ChangeType != ChangeTypeAdded {
		t.Errorf("Expected ChangeTypeAdded, got %v", categoryChange.ChangeType)
	}
	
	if categoryChange.OldValue != nil {
		t.Errorf("Expected old value = nil for added field, got %v", categoryChange.OldValue)
	}
	
	if categoryChange.NewValue != "Electronics" {
		t.Errorf("Expected new value = Electronics, got %v", categoryChange.NewValue)
	}
}

func TestDeltaDetector_DetectDelta_FieldRemoved(t *testing.T) {
	detector := NewDeltaDetector()
	
	oldFact := map[string]interface{}{
		"id":       "123",
		"price":    100.0,
		"category": "Electronics",
	}
	
	newFact := map[string]interface{}{
		"id":    "123",
		"price": 100.0,
	}
	
	delta, err := detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if len(delta.Fields) != 1 {
		t.Errorf("Expected 1 change (removed field), got %d", len(delta.Fields))
	}
	
	categoryChange, exists := delta.Fields["category"]
	if !exists {
		t.Fatal("Expected 'category' in delta")
	}
	
	if categoryChange.ChangeType != ChangeTypeRemoved {
		t.Errorf("Expected ChangeTypeRemoved, got %v", categoryChange.ChangeType)
	}
	
	if categoryChange.OldValue != "Electronics" {
		t.Errorf("Expected old value = Electronics, got %v", categoryChange.OldValue)
	}
	
	if categoryChange.NewValue != nil {
		t.Errorf("Expected new value = nil for removed field, got %v", categoryChange.NewValue)
	}
}

func TestDeltaDetector_DetectDelta_IgnoreInternalFields(t *testing.T) {
	config := DefaultDetectorConfig()
	config.IgnoreInternalFields = true
	detector := NewDeltaDetectorWithConfig(config)
	
	oldFact := map[string]interface{}{
		"id":        "123",
		"price":     100.0,
		"_internal": "old_value",
	}
	
	newFact := map[string]interface{}{
		"id":        "123",
		"price":     100.0,
		"_internal": "new_value",
	}
	
	delta, err := detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if !delta.IsEmpty() {
		t.Errorf("Expected no changes (internal field ignored), got %d changes", len(delta.Fields))
	}
}

func TestDeltaDetector_DetectDelta_IgnoredFieldsList(t *testing.T) {
	config := DefaultDetectorConfig()
	config.IgnoredFields = []string{"timestamp", "updated_at"}
	detector := NewDeltaDetectorWithConfig(config)
	
	oldFact := map[string]interface{}{
		"id":         "123",
		"price":      100.0,
		"timestamp":  "2024-01-01T00:00:00Z",
		"updated_at": "2024-01-01T00:00:00Z",
	}
	
	newFact := map[string]interface{}{
		"id":         "123",
		"price":      100.0,
		"timestamp":  "2024-01-02T00:00:00Z",
		"updated_at": "2024-01-02T00:00:00Z",
	}
	
	delta, err := detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if !delta.IsEmpty() {
		t.Errorf("Expected no changes (ignored fields), got %d changes", len(delta.Fields))
	}
}

func TestDeltaDetector_DetectDelta_FloatEpsilon(t *testing.T) {
	config := DefaultDetectorConfig()
	config.FloatEpsilon = 0.01 // Tolérance de 1%
	detector := NewDeltaDetectorWithConfig(config)
	
	oldFact := map[string]interface{}{
		"price": 100.0,
	}
	
	newFact := map[string]interface{}{
		"price": 100.005, // Différence < epsilon
	}
	
	delta, err := detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if !delta.IsEmpty() {
		t.Error("Expected no change (within epsilon tolerance)")
	}
}

func TestDeltaDetector_DetectDelta_FloatOutsideEpsilon(t *testing.T) {
	config := DefaultDetectorConfig()
	config.FloatEpsilon = 0.01
	detector := NewDeltaDetectorWithConfig(config)
	
	oldFact := map[string]interface{}{
		"price": 100.0,
	}
	
	newFact := map[string]interface{}{
		"price": 100.5, // Différence > epsilon
	}
	
	delta, err := detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if delta.IsEmpty() {
		t.Error("Expected change detected (outside epsilon)")
	}
}

func TestDeltaDetector_DetectDelta_TypeChange(t *testing.T) {
	config := DefaultDetectorConfig()
	config.TrackTypeChanges = true
	detector := NewDeltaDetectorWithConfig(config)
	
	oldFact := map[string]interface{}{
		"value": 42,
	}
	
	newFact := map[string]interface{}{
		"value": "42", // int → string
	}
	
	delta, err := detector.DetectDelta(oldFact, newFact, "Test~1", "Test")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if delta.IsEmpty() {
		t.Error("Expected change detected (type change)")
	}
}

func TestDeltaDetector_DetectDelta_NestedMaps(t *testing.T) {
	config := DefaultDetectorConfig()
	config.EnableDeepComparison = true
	detector := NewDeltaDetectorWithConfig(config)
	
	oldFact := map[string]interface{}{
		"id": "123",
		"address": map[string]interface{}{
			"city":  "Paris",
			"zip":   "75001",
		},
	}
	
	newFact := map[string]interface{}{
		"id": "123",
		"address": map[string]interface{}{
			"city":  "Lyon",
			"zip":   "75001",
		},
	}
	
	delta, err := detector.DetectDelta(oldFact, newFact, "Customer~123", "Customer")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if delta.IsEmpty() {
		t.Error("Expected change in nested map")
	}
	
	// Vérifier que 'address' a changé
	if _, exists := delta.Fields["address"]; !exists {
		t.Error("Expected 'address' field in delta")
	}
}

func TestDeltaDetector_DetectDeltaQuick_NoChanges(t *testing.T) {
	detector := NewDeltaDetector()
	
	fact := map[string]interface{}{
		"id":    "123",
		"price": 100.0,
	}
	
	delta, err := detector.DetectDeltaQuick(fact, fact, "Product~123", "Product")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if delta != nil {
		t.Error("Expected nil delta for no changes (quick detect)")
	}
}

func TestDeltaDetector_DetectDeltaQuick_WithChanges(t *testing.T) {
	detector := NewDeltaDetector()
	
	oldFact := map[string]interface{}{
		"id":    "123",
		"price": 100.0,
	}
	
	newFact := map[string]interface{}{
		"id":    "123",
		"price": 150.0,
	}
	
	delta, err := detector.DetectDeltaQuick(oldFact, newFact, "Product~123", "Product")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if delta == nil {
		t.Fatal("Expected delta to be returned")
	}
	
	if delta.IsEmpty() {
		t.Error("Expected changes in delta")
	}
}

func TestDeltaDetector_ChangeRatio(t *testing.T) {
	detector := NewDeltaDetector()
	
	oldFact := map[string]interface{}{
		"field1":  "value1",
		"field2":  "value2",
		"field3":  "value3",
		"field4":  "value4",
		"field5":  "value5",
		"field6":  "value6",
		"field7":  "value7",
		"field8":  "value8",
		"field9":  "value9",
		"field10": "value10",
	}
	
	// Modifier 2 champs sur 10 = 20%
	newFact := make(map[string]interface{})
	for k, v := range oldFact {
		newFact[k] = v
	}
	newFact["field1"] = "modified1"
	newFact["field2"] = "modified2"
	
	delta, err := detector.DetectDelta(oldFact, newFact, "Test~1", "Test")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	ratio := delta.ChangeRatio()
	expectedRatio := 0.2 // 2/10
	
	if ratio != expectedRatio {
		t.Errorf("Expected change ratio = %v, got %v", expectedRatio, ratio)
	}
}
```

**Contenu** (partie 2/3 - Tests cache et métriques) :

```go
func TestDeltaDetector_CacheEnabled(t *testing.T) {
	config := DefaultDetectorConfig()
	config.CacheComparisons = true
	config.CacheTTL = 1 * time.Minute
	detector := NewDeltaDetectorWithConfig(config)
	
	oldFact := map[string]interface{}{"price": 100.0}
	newFact := map[string]interface{}{"price": 150.0}
	
	// Première détection
	delta1, err := detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Deuxième détection (devrait venir du cache)
	delta2, err := detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Vérifier que les deltas sont identiques
	if len(delta1.Fields) != len(delta2.Fields) {
		t.Error("Cache returned different delta")
	}
	
	// Vérifier métriques cache
	metrics := detector.GetMetrics()
	if metrics.CacheHits == 0 {
		t.Error("Expected cache hit, got 0")
	}
}

func TestDeltaDetector_CacheExpiration(t *testing.T) {
	config := DefaultDetectorConfig()
	config.CacheComparisons = true
	config.CacheTTL = 10 * time.Millisecond
	detector := NewDeltaDetectorWithConfig(config)
	
	oldFact := map[string]interface{}{"price": 100.0}
	newFact := map[string]interface{}{"price": 150.0}
	
	// Première détection
	_, err := detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Attendre expiration
	time.Sleep(20 * time.Millisecond)
	
	// Deuxième détection (cache expiré)
	detector.ResetMetrics() // Reset pour isoler le test
	_, err = detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	metrics := detector.GetMetrics()
	if metrics.CacheHits > 0 {
		t.Error("Expected cache miss after expiration")
	}
}

func TestDeltaDetector_ClearCache(t *testing.T) {
	config := DefaultDetectorConfig()
	config.CacheComparisons = true
	detector := NewDeltaDetectorWithConfig(config)
	
	oldFact := map[string]interface{}{"price": 100.0}
	newFact := map[string]interface{}{"price": 150.0}
	
	// Ajouter au cache
	_, err := detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	metrics := detector.GetMetrics()
	if metrics.CacheSize == 0 {
		t.Fatal("Expected cache to have entries")
	}
	
	// Clear cache
	detector.ClearCache()
	
	metrics = detector.GetMetrics()
	if metrics.CacheSize != 0 {
		t.Errorf("Expected cache size = 0 after clear, got %d", metrics.CacheSize)
	}
}

func TestDeltaDetector_GetMetrics(t *testing.T) {
	detector := NewDeltaDetector()
	
	oldFact := map[string]interface{}{"price": 100.0}
	newFact := map[string]interface{}{"price": 150.0}
	
	// Faire plusieurs détections
	for i := 0; i < 5; i++ {
		_, _ = detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	}
	
	metrics := detector.GetMetrics()
	
	if metrics.Comparisons != 5 {
		t.Errorf("Expected 5 comparisons, got %d", metrics.Comparisons)
	}
}

func TestDeltaDetector_ResetMetrics(t *testing.T) {
	detector := NewDeltaDetector()
	
	oldFact := map[string]interface{}{"price": 100.0}
	newFact := map[string]interface{}{"price": 150.0}
	
	// Faire quelques détections
	_, _ = detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	_, _ = detector.DetectDelta(oldFact, newFact, "Product~124", "Product")
	
	metrics := detector.GetMetrics()
	if metrics.Comparisons == 0 {
		t.Fatal("Expected non-zero comparisons before reset")
	}
	
	// Reset
	detector.ResetMetrics()
	
	metrics = detector.GetMetrics()
	if metrics.Comparisons != 0 {
		t.Errorf("Expected 0 comparisons after reset, got %d", metrics.Comparisons)
	}
}

func TestDeltaDetector_ConcurrentDetection(t *testing.T) {
	detector := NewDeltaDetector()
	
	oldFact := map[string]interface{}{"price": 100.0}
	newFact := map[string]interface{}{"price": 150.0}
	
	// Lancer plusieurs goroutines
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				_, err := detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
				if err != nil {
					t.Errorf("Goroutine %d: unexpected error: %v", id, err)
				}
			}
			done <- true
		}(i)
	}
	
	// Attendre fin
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Vérifier métriques cohérentes
	metrics := detector.GetMetrics()
	if metrics.Comparisons != 1000 {
		t.Errorf("Expected 1000 comparisons, got %d", metrics.Comparisons)
	}
}
```

### Benchmarks : `rete/delta/detector_benchmark_test.go`

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
)

func BenchmarkDeltaDetector_DetectDelta_NoChanges(b *testing.B) {
	detector := NewDeltaDetector()
	
	fact := map[string]interface{}{
		"id":       "123",
		"price":    100.0,
		"status":   "active",
		"quantity": 10,
		"category": "Electronics",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.DetectDelta(fact, fact, "Product~123", "Product")
	}
}

func BenchmarkDeltaDetector_DetectDelta_SingleChange(b *testing.B) {
	detector := NewDeltaDetector()
	
	oldFact := map[string]interface{}{
		"id":       "123",
		"price":    100.0,
		"status":   "active",
		"quantity": 10,
		"category": "Electronics",
	}
	
	newFact := map[string]interface{}{
		"id":       "123",
		"price":    150.0,
		"status":   "active",
		"quantity": 10,
		"category": "Electronics",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	}
}

func BenchmarkDeltaDetector_DetectDelta_MultipleChanges(b *testing.B) {
	detector := NewDeltaDetector()
	
	oldFact := map[string]interface{}{
		"id":       "123",
		"price":    100.0,
		"status":   "active",
		"quantity": 10,
		"category": "Electronics",
	}
	
	newFact := map[string]interface{}{
		"id":       "123",
		"price":    150.0,
		"status":   "inactive",
		"quantity": 5,
		"category": "Electronics",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	}
}

func BenchmarkDeltaDetector_DetectDeltaQuick_NoChanges(b *testing.B) {
	detector := NewDeltaDetector()
	
	fact := map[string]interface{}{
		"id":       "123",
		"price":    100.0,
		"status":   "active",
		"quantity": 10,
		"category": "Electronics",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.DetectDeltaQuick(fact, fact, "Product~123", "Product")
	}
}

func BenchmarkDeltaDetector_DetectDelta_LargeFact(b *testing.B) {
	detector := NewDeltaDetector()
	
	// Fait avec 50 champs
	oldFact := make(map[string]interface{})
	newFact := make(map[string]interface{})
	
	for i := 0; i < 50; i++ {
		fieldName := "field" + string(rune('0'+i))
		oldFact[fieldName] = i
		newFact[fieldName] = i
	}
	
	// Modifier 5 champs
	for i := 0; i < 5; i++ {
		fieldName := "field" + string(rune('0'+i))
		newFact[fieldName] = i + 1000
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.DetectDelta(oldFact, newFact, "Large~1", "Large")
	}
}

func BenchmarkDeltaDetector_DetectDelta_WithCache(b *testing.B) {
	config := DefaultDetectorConfig()
	config.CacheComparisons = true
	detector := NewDeltaDetectorWithConfig(config)
	
	oldFact := map[string]interface{}{
		"id":    "123",
		"price": 100.0,
	}
	
	newFact := map[string]interface{}{
		"id":    "123",
		"price": 150.0,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	}
}

func BenchmarkDeltaDetector_DetectDelta_DeepNested(b *testing.B) {
	config := DefaultDetectorConfig()
	config.EnableDeepComparison = true
	config.MaxNestingLevel = 5
	detector := NewDeltaDetectorWithConfig(config)
	
	oldFact := map[string]interface{}{
		"id": "123",
		"nested": map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"value": "old",
				},
			},
		},
	}
	
	newFact := map[string]interface{}{
		"id": "123",
		"nested": map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"value": "new",
				},
			},
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.DetectDelta(oldFact, newFact, "Nested~1", "Nested")
	}
}
```

---

## ✅ Validation

Après implémentation, exécuter :

```bash
# 1. Formattage
go fmt ./rete/delta/...
goimports -w ./rete/delta/

# 2. Validation statique
go vet ./rete/delta/...
staticcheck ./rete/delta/...

# 3. Tests unitaires
go test ./rete/delta/... -v
go test ./rete/delta/... -cover

# 4. Benchmarks
go test ./rete/delta/... -bench=. -benchmem

# 5. Race detector
go test ./rete/delta/... -race

# 6. Validation complète
make validate
```

**Critères de succès** :
- [ ] Tous les tests passent (100%)
- [ ] Couverture > 90%
- [ ] Aucune erreur `go vet`, `staticcheck`
- [ ] Benchmarks montrent performance acceptable
- [ ] Aucune race condition détectée
- [ ] DetectDeltaQuick plus rapide que DetectDelta pour cas no-op

---

## 📊 Livrables

À la fin de ce prompt :

1. **Code** :
   - ✅ `rete/delta/detector_config.go` - Configuration détecteur
   - ✅ `rete/delta/delta_detector.go` - Implémentation détecteur
   
2. **Tests** :
   - ✅ `rete/delta/detector_config_test.go`
   - ✅ `rete/delta/delta_detector_test.go`
   - ✅ `rete/delta/detector_benchmark_test.go`

3. **Validation** :
   - ✅ Rapport de couverture (> 90%)
   - ✅ Résultats benchmarks
   - ✅ Rapport race detector (0 races)

---

## 🚀 Commit

Une fois validé :

```bash
git add rete/delta/
git commit -m "feat(rete): [Prompt 04] Implémentation détection delta

- DeltaDetector avec configuration flexible
- Détection précise des changements de champs
- Support comparaison profonde (nested maps/slices)
- Gestion champs ignorés et internes
- Cache optionnel avec TTL
- Optimisation DetectDeltaQuick pour cas no-op
- Métriques de performance
- Thread-safe (sync.RWMutex)
- Tests unitaires complets (> 90%