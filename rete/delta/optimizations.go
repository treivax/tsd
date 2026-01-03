// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

// OptimizedValuesEqual est un alias vers ValuesEqual pour rétrocompatibilité.
//
// Note: Cette fonction est maintenant identique à ValuesEqual car celle-ci
// intègre déjà toutes les optimisations (fast paths, inline, etc.).
//
// Utiliser ValuesEqual directement est préférable.
//
// Deprecated: Utiliser ValuesEqual qui intègre déjà les optimisations.
func OptimizedValuesEqual(a, b interface{}, epsilon float64) bool {
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
