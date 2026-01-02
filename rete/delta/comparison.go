// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"reflect"
)

// ValuesEqual compare deux valeurs en profondeur avec optimisations.
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
// Optimisations :
// - Fast path pour types primitifs (évite reflection)
// - Short-circuit pour cas communs
// - Comparaison inline des floats (sans math.Abs)
func ValuesEqual(a, b interface{}, epsilon float64) bool {
	// Fast path 1: nil handling first
	if a == nil || b == nil {
		return a == b
	}

	// Fast path 2: types simples
	if result, handled := compareSimpleTypes(a, b); handled {
		return result
	}

	// Fast path 3: types numériques avec epsilon
	if result, handled := compareNumericTypes(a, b, epsilon); handled {
		return result
	}

	// Slow path: reflection pour types complexes
	return compareComplexTypes(a, b)
}

// compareSimpleTypes compare les types simples non-numériques (string, bool).
// Retourne (résultat, true) si le type est géré, (false, false) sinon.
func compareSimpleTypes(a, b interface{}) (bool, bool) {
	switch va := a.(type) {
	case string:
		vb, ok := b.(string)
		return ok && va == vb, true

	case bool:
		vb, ok := b.(bool)
		return ok && va == vb, true
	}

	return false, false
}

// compareNumericTypes compare les types numériques (int, uint, float).
// Retourne (résultat, true) si le type est géré, (false, false) sinon.
func compareNumericTypes(a, b interface{}, epsilon float64) (bool, bool) {
	// Essayer comparaison signed integers
	if result, handled := compareSignedIntegers(a, b); handled {
		return result, true
	}

	// Essayer comparaison unsigned integers
	if result, handled := compareUnsignedIntegers(a, b); handled {
		return result, true
	}

	// Essayer comparaison floats
	if result, handled := compareFloats(a, b, epsilon); handled {
		return result, true
	}

	return false, false
}

// compareSignedIntegers compare les entiers signés (int, int8, int16, int32, int64).
func compareSignedIntegers(a, b interface{}) (bool, bool) {
	switch va := a.(type) {
	case int:
		vb, ok := b.(int)
		return ok && va == vb, true

	case int64:
		vb, ok := b.(int64)
		return ok && va == vb, true

	case int32:
		vb, ok := b.(int32)
		return ok && va == vb, true

	case int16:
		vb, ok := b.(int16)
		return ok && va == vb, true

	case int8:
		vb, ok := b.(int8)
		return ok && va == vb, true
	}

	return false, false
}

// compareUnsignedIntegers compare les entiers non-signés (uint, uint8, uint16, uint32, uint64).
func compareUnsignedIntegers(a, b interface{}) (bool, bool) {
	switch va := a.(type) {
	case uint:
		vb, ok := b.(uint)
		return ok && va == vb, true

	case uint64:
		vb, ok := b.(uint64)
		return ok && va == vb, true

	case uint32:
		vb, ok := b.(uint32)
		return ok && va == vb, true

	case uint16:
		vb, ok := b.(uint16)
		return ok && va == vb, true

	case uint8:
		vb, ok := b.(uint8)
		return ok && va == vb, true
	}

	return false, false
}

// compareFloats compare les nombres à virgule flottante (float32, float64) avec epsilon.
func compareFloats(a, b interface{}, epsilon float64) (bool, bool) {
	switch va := a.(type) {
	case float64:
		vb, ok := b.(float64)
		if !ok {
			return false, true
		}
		return compareFloat64WithEpsilon(va, vb, epsilon), true

	case float32:
		vb, ok := b.(float32)
		if !ok {
			return false, true
		}
		return compareFloat32WithEpsilon(va, vb, epsilon), true
	}

	return false, false
}

// compareFloat64WithEpsilon compare deux float64 avec epsilon (optimisé sans math.Abs).
func compareFloat64WithEpsilon(a, b, epsilon float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= epsilon
}

// compareFloat32WithEpsilon compare deux float32 avec epsilon (optimisé sans math.Abs).
func compareFloat32WithEpsilon(a, b float32, epsilon float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return float64(diff) <= epsilon
}

// compareComplexTypes compare les types complexes via reflection.
func compareComplexTypes(a, b interface{}) bool {
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
