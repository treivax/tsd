// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import "time"

const (
	// DefaultRetentionDuration durée par défaut si aucune n'est spécifiée
	DefaultRetentionDuration = 1 * time.Hour
)

// UnlimitedRetentionPolicy conserve les xuples indéfiniment.
//
// Comportement :
//   - Pas d'expiration basée sur le temps
//   - Les xuples complètement consommés (XupleStateConsumed) peuvent être nettoyés
//   - Les xuples expirés (XupleStateExpired) peuvent être nettoyés
//
// Note: Même avec cette politique, Cleanup() retirera les xuples consommés/expirés
// pour éviter les fuites mémoire.
type UnlimitedRetentionPolicy struct{}

// NewUnlimitedRetentionPolicy crée une nouvelle politique illimitée.
func NewUnlimitedRetentionPolicy() *UnlimitedRetentionPolicy {
	return &UnlimitedRetentionPolicy{}
}

// ComputeExpiration retourne zero time (pas d'expiration basée sur le temps).
func (p *UnlimitedRetentionPolicy) ComputeExpiration(createdAt time.Time) time.Time {
	return time.Time{} // Zero time = pas d'expiration
}

// ShouldRetain retourne true pour les xuples disponibles.
// Les xuples consommés ou expirés peuvent être nettoyés.
func (p *UnlimitedRetentionPolicy) ShouldRetain(xuple *Xuple) bool {
	// Nettoyer uniquement les xuples complètement consommés ou expirés
	return xuple.Metadata.State == XupleStateAvailable
}

// Name retourne le nom de la politique.
func (p *UnlimitedRetentionPolicy) Name() string {
	return "unlimited"
}

// DurationRetentionPolicy expire les xuples après une durée.
type DurationRetentionPolicy struct {
	Duration time.Duration
}

// NewDurationRetentionPolicy crée une nouvelle politique basée sur la durée.
func NewDurationRetentionPolicy(duration time.Duration) *DurationRetentionPolicy {
	if duration <= 0 {
		duration = DefaultRetentionDuration
	}
	return &DurationRetentionPolicy{
		Duration: duration,
	}
}

// ComputeExpiration calcule la date d'expiration.
func (p *DurationRetentionPolicy) ComputeExpiration(createdAt time.Time) time.Time {
	return createdAt.Add(p.Duration)
}

// ShouldRetain vérifie si le xuple doit être conservé.
func (p *DurationRetentionPolicy) ShouldRetain(xuple *Xuple) bool {
	if xuple.Metadata.ExpiresAt.IsZero() {
		return true
	}
	return time.Now().Before(xuple.Metadata.ExpiresAt)
}

// Name retourne le nom de la politique.
func (p *DurationRetentionPolicy) Name() string {
	return "duration"
}
