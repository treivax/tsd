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
