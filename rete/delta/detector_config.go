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
		return newInvalidConfigError("FloatEpsilon", "must be >= 0")
	}

	if dc.MaxNestingLevel < 1 {
		return newInvalidConfigError("MaxNestingLevel", "must be >= 1")
	}

	if dc.CacheTTL < 0 {
		return newInvalidConfigError("CacheTTL", "must be >= 0")
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
