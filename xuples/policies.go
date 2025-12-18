// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import "time"

// SelectionPolicy définit comment sélectionner un xuple parmi plusieurs disponibles.
type SelectionPolicy interface {
	// Select sélectionne un xuple parmi une liste de xuples disponibles.
	// Retourne nil si aucun xuple n'est disponible.
	Select(xuples []*Xuple) *Xuple

	// Name retourne le nom de la politique.
	Name() string
}

// ConsumptionPolicy définit comment les xuples peuvent être consommés.
type ConsumptionPolicy interface {
	// CanConsume vérifie si un agent peut consommer un xuple.
	CanConsume(xuple *Xuple, agentID string) bool

	// OnConsumed est appelé après qu'un xuple ait été consommé.
	// Retourne true si le xuple doit être marqué comme complètement consommé.
	OnConsumed(xuple *Xuple, agentID string) bool

	// Name retourne le nom de la politique.
	Name() string
}

// RetentionPolicy définit la durée de vie des xuples.
type RetentionPolicy interface {
	// ComputeExpiration calcule la date d'expiration pour un nouveau xuple.
	// Retourne zero time si pas d'expiration.
	ComputeExpiration(createdAt time.Time) time.Time

	// ShouldRetain vérifie si un xuple doit être conservé.
	ShouldRetain(xuple *Xuple) bool

	// Name retourne le nom de la politique.
	Name() string
}

// PolicyType représente le type de politique.
type PolicyType int

const (
	// PolicyTypeSelection politique de sélection
	PolicyTypeSelection PolicyType = iota

	// PolicyTypeConsumption politique de consommation
	PolicyTypeConsumption

	// PolicyTypeRetention politique de rétention
	PolicyTypeRetention
)

// String retourne la représentation textuelle du type.
func (p PolicyType) String() string {
	switch p {
	case PolicyTypeSelection:
		return "selection"
	case PolicyTypeConsumption:
		return "consumption"
	case PolicyTypeRetention:
		return "retention"
	default:
		return "unknown"
	}
}
