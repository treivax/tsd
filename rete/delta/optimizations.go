// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"math"
	"reflect"
)

// OptimizedValuesEqual est une version optimisée de ValuesEqual.
//
// Optimisations :
// - Short-circuit pour types simples
// - Évite reflect.DeepEqual quand possible
// - Inline hints pour compiler
func OptimizedValuesEqual(a, b interface{}, epsilon float64) bool {
	// Fast path : nil handling first
	if a == nil || b == nil {
		return a == b
	}

	// Fast path : types simples (inlined par compiler)
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
		// Comparaison float optimisée
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

	// Slow path : comparaison générique
	return ValuesEqual(a, b, epsilon)
}

// FastHashString calcule un hash rapide d'une string (non-cryptographique).
//
// Utilise l'algorithme FNV-1a pour rapidité et bonne distribution.
func FastHashString(s string) uint64 {
	// FNV-1a hash
	const (
		offset64 = 14695981039346656037
		prime64  = 1099511628211
	)

	hash := uint64(offset64)
	for i := 0; i < len(s); i++ {
		hash ^= uint64(s[i])
		hash *= prime64
	}
	return hash
}

// FastHashBytes calcule un hash rapide de bytes.
func FastHashBytes(b []byte) uint64 {
	const (
		offset64 = 14695981039346656037
		prime64  = 1099511628211
	)

	hash := uint64(offset64)
	for i := 0; i < len(b); i++ {
		hash ^= uint64(b[i])
		hash *= prime64
	}
	return hash
}

// PreallocatedMap crée une map pré-allouée avec capacité.
func PreallocatedMap(size int) map[string]interface{} {
	return make(map[string]interface{}, size)
}

// CopyFactFast copie rapidement un fait (optimisé).
func CopyFactFast(fact map[string]interface{}) map[string]interface{} {
	if len(fact) == 0 {
		return make(map[string]interface{})
	}

	copy := make(map[string]interface{}, len(fact))
	for k, v := range fact {
		copy[k] = v
	}
	return copy
}

// BatchNodeReferences groupe les références de nœuds pour traitement batch.
type BatchNodeReferences struct {
	alpha    []NodeReference
	beta     []NodeReference
	terminal []NodeReference
}

// NewBatchNodeReferences crée un batch pré-alloué.
func NewBatchNodeReferences(expectedSize int) *BatchNodeReferences {
	thirdSize := expectedSize / 3
	if thirdSize < 4 {
		thirdSize = 4
	}

	return &BatchNodeReferences{
		alpha:    make([]NodeReference, 0, thirdSize),
		beta:     make([]NodeReference, 0, thirdSize),
		terminal: make([]NodeReference, 0, thirdSize),
	}
}

// Add ajoute une référence au batch approprié.
func (b *BatchNodeReferences) Add(ref NodeReference) {
	switch ref.NodeType {
	case NodeTypeAlpha:
		b.alpha = append(b.alpha, ref)
	case NodeTypeBeta:
		b.beta = append(b.beta, ref)
	case NodeTypeTerminal:
		b.terminal = append(b.terminal, ref)
	}
}

// ProcessInOrder traite les nœuds dans l'ordre optimal (Alpha → Beta → Terminal).
func (b *BatchNodeReferences) ProcessInOrder(processor func(NodeReference) error) error {
	// Alpha d'abord
	for i := range b.alpha {
		if err := processor(b.alpha[i]); err != nil {
			return err
		}
	}

	// Beta ensuite
	for i := range b.beta {
		if err := processor(b.beta[i]); err != nil {
			return err
		}
	}

	// Terminal en dernier
	for i := range b.terminal {
		if err := processor(b.terminal[i]); err != nil {
			return err
		}
	}

	return nil
}

// Count retourne le nombre total de nœuds.
func (b *BatchNodeReferences) Count() int {
	return len(b.alpha) + len(b.beta) + len(b.terminal)
}

// quickTypeCheck effectue une vérification rapide de type sans reflection.
//
// Retourne true si les types sont différents (optimisation early-exit).
func quickTypeCheck(a, b interface{}) bool {
	// Utiliser type switch pour éviter reflect
	switch a.(type) {
	case string:
		_, ok := b.(string)
		return !ok
	case int, int8, int16, int32, int64:
		switch b.(type) {
		case int, int8, int16, int32, int64:
			return false
		default:
			return true
		}
	case uint, uint8, uint16, uint32, uint64:
		switch b.(type) {
		case uint, uint8, uint16, uint32, uint64:
			return false
		default:
			return true
		}
	case float32, float64:
		switch b.(type) {
		case float32, float64:
			return false
		default:
			return true
		}
	case bool:
		_, ok := b.(bool)
		return !ok
	case nil:
		return b != nil
	}

	// Fallback: utiliser reflect uniquement si nécessaire
	return reflect.TypeOf(a) != reflect.TypeOf(b)
}

// InlinedAbsFloat64 calcule la valeur absolue (inlined).
func inlinedAbsFloat64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// InlinedAbsFloat32 calcule la valeur absolue (inlined).
func inlinedAbsFloat32(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

// CompareFloatsFast compare deux floats avec epsilon (optimisé inline).
func compareFloatsFast(a, b, epsilon float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= epsilon
}

// FloatEqualityCheck optimisé sans math.Abs.
func floatEqualityCheck(a, b float64, epsilon float64) bool {
	// Cas spéciaux
	if a == b {
		return true
	}

	// NaN check
	if a != a || b != b {
		return false
	}

	// Infinity check
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return a == b
	}

	// Différence absolue
	diff := a - b
	if diff < 0 {
		diff = -diff
	}

	return diff <= epsilon
}
