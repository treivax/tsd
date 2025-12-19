// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
	"time"
)

// Constantes pour la conversion de durées
const (
	SecondsPerMinute = 60
	SecondsPerHour   = 3600
	SecondsPerDay    = 86400
	SecondsPerWeek   = 604800
)

// Constantes pour les politiques de xuple-space
const (
	// Politiques de sélection
	SelectionFIFO   = "fifo"
	SelectionLIFO   = "lifo"
	SelectionRandom = "random"

	// Politiques de consommation
	ConsumptionOnce     = "once"
	ConsumptionPerAgent = "per-agent"
	ConsumptionLimited  = "limited"

	// Politiques de rétention
	RetentionUnlimited = "unlimited"
	RetentionDuration  = "duration"
)

// Valeurs par défaut pour xuple-spaces
const (
	DefaultConsumptionType   = ConsumptionOnce
	DefaultRetentionType     = RetentionUnlimited
	DefaultMaxSize           = 0 // 0 = illimité
	DefaultRetentionDuration = 0 // 0 = illimité
)

// ParseDuration parse une durée au format Go (24h, 7d, etc.) et retourne les secondes.
// Supporte les unités: s (secondes), m (minutes), h (heures), d (jours), w (semaines).
func ParseDuration(s string) (int, error) {
	if s == "" {
		return 0, nil
	}

	// Utiliser le parser Go standard pour les formats standards
	dur, err := time.ParseDuration(s)
	if err == nil {
		return int(dur.Seconds()), nil
	}

	// Sinon, essayer notre parser personnalisé pour les jours et semaines
	if len(s) < 2 {
		return 0, err
	}

	unit := s[len(s)-1]
	valueStr := s[:len(s)-1]

	var value int
	if _, scanErr := fmt.Sscanf(valueStr, "%d", &value); scanErr != nil {
		return 0, fmt.Errorf("invalid duration format '%s': %w", s, scanErr)
	}

	if value <= 0 {
		return 0, fmt.Errorf("duration must be positive, got %d", value)
	}

	switch unit {
	case 'd':
		return value * SecondsPerDay, nil
	case 'w':
		return value * SecondsPerWeek, nil
	default:
		return 0, fmt.Errorf("invalid duration format '%s', expected formats: 24h, 7d, 1w, etc.", s)
	}
}
