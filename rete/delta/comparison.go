// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"math"
	"reflect"
)

// ValuesEqual compare deux valeurs en profondeur.
//
// Cette fonction gère les cas spéciaux :
//   - Floats : comparaison avec epsilon (tolérance)
//   - Maps/slices : comparaison récursive
//   - nil : nil == nil → true
//
// Paramètres :
//   - a, b : valeurs à comparer
//   - epsilon : tolérance pour floats (utiliser DefaultFloatEpsilon)
//
// Retourne true si les valeurs sont égales.
func ValuesEqual(a, b interface{}, epsilon float64) bool {
	// Cas 1: nil vs non-nil
	if a == nil || b == nil {
		return a == b
	}

	// Cas 2: Types différents → inégaux
	typeA := reflect.TypeOf(a)
	typeB := reflect.TypeOf(b)
	if typeA != typeB {
		return false
	}

	// Cas 3: Floats avec tolérance epsilon
	switch va := a.(type) {
	case float64:
		vb, ok := b.(float64)
		if !ok {
			return false
		}
		return math.Abs(va-vb) <= epsilon

	case float32:
		vb, ok := b.(float32)
		if !ok {
			return false
		}
		return math.Abs(float64(va-vb)) <= epsilon
	}

	// Cas 4: Types comparables directement (int, string, bool, etc.)
	// Vérifier si le type est comparable
	if typeA.Comparable() {
		return a == b
	}

	// Cas 5: Comparaison profonde (maps, slices, structs non comparables)
	return reflect.DeepEqual(a, b)
}

// FactsEqual compare deux faits (représentés comme maps) en profondeur.
//
// Paramètres :
//   - fact1, fact2 : faits à comparer (map[string]interface{})
//
// Retourne true si les faits sont identiques (mêmes champs, mêmes valeurs).
func FactsEqual(fact1, fact2 map[string]interface{}) bool {
	if len(fact1) != len(fact2) {
		return false
	}

	for key, val1 := range fact1 {
		val2, exists := fact2[key]
		if !exists {
			return false
		}
		if !ValuesEqual(val1, val2, DefaultFloatEpsilon) {
			return false
		}
	}

	return true
}
