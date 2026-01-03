// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"strings"
	"sync"
	"time"
)

// FactDeltaPool est un pool d'objets FactDelta pour réduire les allocations.
var FactDeltaPool = sync.Pool{
	New: func() interface{} {
		return &FactDelta{
			Fields: make(map[string]FieldDelta, 8), // Pré-allouer taille typique
		}
	},
}

// AcquireFactDelta obtient un FactDelta depuis le pool.
func AcquireFactDelta(factID, factType string) *FactDelta {
	delta := FactDeltaPool.Get().(*FactDelta)
	delta.FactID = factID
	delta.FactType = factType
	delta.Timestamp = time.Now()
	return delta
}

// ReleaseFactDelta retourne un FactDelta au pool.
func ReleaseFactDelta(delta *FactDelta) {
	if delta == nil {
		return
	}

	// Reset pour réutilisation
	delta.FactID = ""
	delta.FactType = ""
	delta.FieldCount = 0

	// Clear map sans réallouer
	for k := range delta.Fields {
		delete(delta.Fields, k)
	}

	FactDeltaPool.Put(delta)
}

// NodeReferencePool est un pool de slices de NodeReference.
var NodeReferencePool = sync.Pool{
	New: func() interface{} {
		slice := make([]NodeReference, 0, 16)
		return &slice
	},
}

// AcquireNodeReferenceSlice obtient une slice depuis le pool.
func AcquireNodeReferenceSlice() *[]NodeReference {
	slice := NodeReferencePool.Get().(*[]NodeReference)
	*slice = (*slice)[:0] // Reset length, keep capacity
	return slice
}

// ReleaseNodeReferenceSlice retourne une slice au pool.
func ReleaseNodeReferenceSlice(slice *[]NodeReference) {
	if slice == nil {
		return
	}

	if cap(*slice) > maxPooledSliceCapacity {
		// Trop grande, ne pas réutiliser
		return
	}
	NodeReferencePool.Put(slice)
}

// StringBuilderPool est un pool de strings.Builder pour construction efficace.
var StringBuilderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

// AcquireStringBuilder obtient un builder depuis le pool.
func AcquireStringBuilder() *strings.Builder {
	sb := StringBuilderPool.Get().(*strings.Builder)
	sb.Reset()
	return sb
}

// ReleaseStringBuilder retourne un builder au pool.
func ReleaseStringBuilder(sb *strings.Builder) {
	if sb == nil {
		return
	}

	if sb.Cap() > maxPooledBuilderCapacity {
		// Trop grand, ne pas réutiliser
		return
	}
	StringBuilderPool.Put(sb)
}

// MapPool est un pool de maps pour allocations temporaires.
var MapPool = sync.Pool{
	New: func() interface{} {
		m := make(map[string]interface{}, 16)
		return &m
	},
}

// AcquireMap obtient une map depuis le pool.
func AcquireMap() *map[string]interface{} {
	m := MapPool.Get().(*map[string]interface{})
	// Clear la map
	for k := range *m {
		delete(*m, k)
	}
	return m
}

// ReleaseMap retourne une map au pool.
func ReleaseMap(m *map[string]interface{}) {
	if m == nil {
		return
	}

	// Clear avant retour
	for k := range *m {
		delete(*m, k)
	}

	if len(*m) > maxPooledMapSize {
		// Trop grande, ne pas réutiliser
		return
	}

	MapPool.Put(m)
}

// Pool capacity limits
const (
	maxPooledSliceCapacity   = 1024
	maxPooledBuilderCapacity = 4096
	maxPooledMapSize         = 128
)

// WithFactDelta exécute une fonction avec un FactDelta acquis depuis le pool.
//
// Le FactDelta est automatiquement retourné au pool à la fin de l'exécution,
// même en cas d'erreur ou de panic. Ce pattern garantit qu'aucune fuite
// mémoire ne se produira.
//
// Utilisation recommandée pour tout code utilisant le pool :
//
//	err := WithFactDelta(factID, factType, func(delta *FactDelta) error {
//	    // Remplir le delta
//	    delta.Fields["name"] = FieldDelta{...}
//
//	    // Propager le delta
//	    return propagator.Propagate(context.Background(), oldFact, newFact)
//	})
//
// Paramètres :
//   - factID : identifiant du fait
//   - factType : type du fait
//   - fn : fonction à exécuter avec le delta acquis
//
// Retourne :
//   - error : erreur retournée par fn, ou nil
func WithFactDelta(factID, factType string, fn func(*FactDelta) error) error {
	delta := AcquireFactDelta(factID, factType)
	defer ReleaseFactDelta(delta)
	return fn(delta)
}

// WithNodeReferenceSlice exécute une fonction avec une slice acquise depuis le pool.
//
// La slice est automatiquement retournée au pool à la fin de l'exécution.
//
// Utilisation :
//
//	err := WithNodeReferenceSlice(func(nodes *[]NodeReference) error {
//	    *nodes = append(*nodes, NodeReference{NodeID: "n1"})
//	    return processNodes(*nodes)
//	})
//
// Paramètres :
//   - fn : fonction à exécuter avec la slice acquise
//
// Retourne :
//   - error : erreur retournée par fn, ou nil
func WithNodeReferenceSlice(fn func(*[]NodeReference) error) error {
	slice := AcquireNodeReferenceSlice()
	defer ReleaseNodeReferenceSlice(slice)
	return fn(slice)
}

// WithStringBuilder exécute une fonction avec un StringBuilder acquis depuis le pool.
//
// Le builder est automatiquement retourné au pool à la fin de l'exécution.
//
// Utilisation :
//
//	result, err := WithStringBuilderResult(func(sb *strings.Builder) (string, error) {
//	    sb.WriteString("Hello")
//	    sb.WriteString(" World")
//	    return sb.String(), nil
//	})
//
// Paramètres :
//   - fn : fonction à exécuter avec le builder acquis
//
// Retourne :
//   - error : erreur retournée par fn, ou nil
func WithStringBuilder(fn func(*strings.Builder) error) error {
	sb := AcquireStringBuilder()
	defer ReleaseStringBuilder(sb)
	return fn(sb)
}

// WithStringBuilderResult exécute une fonction avec un StringBuilder et retourne un résultat.
//
// Le builder est automatiquement retourné au pool à la fin de l'exécution.
//
// Paramètres :
//   - fn : fonction à exécuter avec le builder acquis
//
// Retourne :
//   - T : résultat retourné par fn
//   - error : erreur retournée par fn, ou nil
func WithStringBuilderResult[T any](fn func(*strings.Builder) (T, error)) (T, error) {
	sb := AcquireStringBuilder()
	defer ReleaseStringBuilder(sb)
	return fn(sb)
}

// WithMap exécute une fonction avec une map acquise depuis le pool.
//
// La map est automatiquement retournée au pool à la fin de l'exécution.
//
// Utilisation :
//
//	err := WithMap(func(m *map[string]interface{}) error {
//	    (*m)["key"] = "value"
//	    return processMap(*m)
//	})
//
// Paramètres :
//   - fn : fonction à exécuter avec la map acquise
//
// Retourne :
//   - error : erreur retournée par fn, ou nil
func WithMap(fn func(*map[string]interface{}) error) error {
	m := AcquireMap()
	defer ReleaseMap(m)
	return fn(m)
}
