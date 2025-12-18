// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import (
	"math/rand"
	"time"
)

// selectByTimestamp sélectionne le xuple selon un comparateur temporel.
//
// Cette fonction helper permet de factoriser la logique de sélection basée
// sur le timestamp pour FIFO et LIFO.
//
// Paramètres:
//   - xuples: liste de xuples parmi lesquels sélectionner
//   - older: si true, retourne le plus ancien; sinon le plus récent
//
// Retourne:
//   - *Xuple: xuple sélectionné, ou nil si la liste est vide
func selectByTimestamp(xuples []*Xuple, older bool) *Xuple {
	if len(xuples) == 0 {
		return nil
	}

	selected := xuples[0]
	for _, xuple := range xuples[1:] {
		if older && xuple.CreatedAt.Before(selected.CreatedAt) {
			selected = xuple
		} else if !older && xuple.CreatedAt.After(selected.CreatedAt) {
			selected = xuple
		}
	}

	return selected
}

// RandomSelectionPolicy sélectionne aléatoirement un xuple.
//
// Utilise un générateur de nombres aléatoires indépendant pour assurer
// la thread-safety et permettre des distributions différentes par instance.
//
// Thread-Safety:
//   - Chaque instance a son propre générateur aléatoire
//   - Non thread-safe en interne : à utiliser avec lock externe (XupleSpace)
type RandomSelectionPolicy struct {
	rng *rand.Rand
}

// NewRandomSelectionPolicy crée une nouvelle politique de sélection aléatoire.
func NewRandomSelectionPolicy() *RandomSelectionPolicy {
	return &RandomSelectionPolicy{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Select sélectionne un xuple aléatoirement.
func (p *RandomSelectionPolicy) Select(xuples []*Xuple) *Xuple {
	if len(xuples) == 0 {
		return nil
	}
	return xuples[p.rng.Intn(len(xuples))]
}

// Name retourne le nom de la politique.
func (p *RandomSelectionPolicy) Name() string {
	return "random"
}

// FIFOSelectionPolicy sélectionne le premier entré (plus ancien).
type FIFOSelectionPolicy struct{}

// NewFIFOSelectionPolicy crée une nouvelle politique FIFO.
func NewFIFOSelectionPolicy() *FIFOSelectionPolicy {
	return &FIFOSelectionPolicy{}
}

// Select sélectionne le xuple le plus ancien.
func (p *FIFOSelectionPolicy) Select(xuples []*Xuple) *Xuple {
	return selectByTimestamp(xuples, true)
}

// Name retourne le nom de la politique.
func (p *FIFOSelectionPolicy) Name() string {
	return "fifo"
}

// LIFOSelectionPolicy sélectionne le dernier entré (plus récent).
type LIFOSelectionPolicy struct{}

// NewLIFOSelectionPolicy crée une nouvelle politique LIFO.
func NewLIFOSelectionPolicy() *LIFOSelectionPolicy {
	return &LIFOSelectionPolicy{}
}

// Select sélectionne le xuple le plus récent.
func (p *LIFOSelectionPolicy) Select(xuples []*Xuple) *Xuple {
	return selectByTimestamp(xuples, false)
}

// Name retourne le nom de la politique.
func (p *LIFOSelectionPolicy) Name() string {
	return "lifo"
}
