// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
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
//
// Note: Pour de meilleures performances, utiliser OptimizedValuesEqual
// qui évite reflect.TypeOf pour les types communs.
func ValuesEqual(a, b interface{}, epsilon float64) bool {
	// Essayer la version optimisée d'abord (fast path)
	// Elle retourne true/false rapidement pour les types simples
	// et délègue à cette fonction pour les cas complexes
	return optimizedValuesEqualInternal(a, b, epsilon)
}

// optimizedValuesEqualInternal implémentation interne avec fast paths.
func optimizedValuesEqualInternal(a, b interface{}, epsilon float64) bool {
	// Fast path 1: nil handling first
	if a == nil || b == nil {
		return a == b
	}

	// Fast path 2: types simples via type switch (évite reflect)
	switch va := a.(type) {
	case string:
		vb, ok := b.(string)
		return ok && va == vb

	case int:
		vb, ok := b.(int)
		return ok && va == vb

	case int64:
		vb, ok := b.(int64)
		return ok && va == vb

	case int32:
		vb, ok := b.(int32)
		return ok && va == vb

	case int16:
		vb, ok := b.(int16)
		return ok && va == vb

	case int8:
		vb, ok := b.(int8)
		return ok && va == vb

	case uint:
		vb, ok := b.(uint)
		return ok && va == vb

	case uint64:
		vb, ok := b.(uint64)
		return ok && va == vb

	case uint32:
		vb, ok := b.(uint32)
		return ok && va == vb

	case uint16:
		vb, ok := b.(uint16)
		return ok && va == vb

	case uint8:
		vb, ok := b.(uint8)
		return ok && va == vb

	case bool:
		vb, ok := b.(bool)
		return ok && va == vb

	case float64:
		vb, ok := b.(float64)
		if !ok {
			return false
		}
		// Optimisé sans math.Abs
		diff := va - vb
		if diff < 0 {
			diff = -diff
		}
		return diff <= epsilon

	case float32:
		vb, ok := b.(float32)
		if !ok {
			return false
		}
		diff := va - vb
		if diff < 0 {
			diff = -diff
		}
		return float64(diff) <= epsilon
	}

	// Slow path: reflection pour types complexes
	typeA := reflect.TypeOf(a)
	typeB := reflect.TypeOf(b)
	if typeA != typeB {
		return false
	}

	// Pour types complexes (maps, slices, etc.), utiliser DeepEqual
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
