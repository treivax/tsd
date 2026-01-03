// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
	"testing"
)

// BenchmarkDetector_WithPooling compare l'impact du pooling.
func BenchmarkDetector_WithPooling(b *testing.B) {
	detector := NewDeltaDetector()

	oldFact := map[string]interface{}{
		"id":       "123",
		"price":    100.0,
		"status":   "active",
		"quantity": 10,
	}

	newFact := map[string]interface{}{
		"id":       "123",
		"price":    150.0,
		"status":   "active",
		"quantity": 10,
	}

	b.Run("WithPoolRelease", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			delta, _ := detector.DetectDelta(oldFact, newFact, "Test~1", "Test")
			ReleaseFactDelta(delta)
		}
	})

	b.Run("WithoutPoolRelease", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			delta, _ := detector.DetectDelta(oldFact, newFact, "Test~1", "Test")
			_ = delta // Pas de release - garbage collected
		}
	})
}

// BenchmarkOptimizedCache_Performance teste les performances du cache optimisé.
func BenchmarkOptimizedCache_Performance(b *testing.B) {
	scenarios := []struct {
		name     string
		size     int
		hitRate  float64 // Taux de hit souhaité (0.0 - 1.0)
		keyRange int     // Plage de clés (pour contrôler hit rate)
	}{
		{"SmallCache_HighHitRate", 100, 0.9, 50},
		{"SmallCache_LowHitRate", 100, 0.1, 1000},
		{"LargeCache_HighHitRate", 10000, 0.9, 5000},
		{"LargeCache_MediumHitRate", 10000, 0.5, 20000},
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			cache := NewOptimizedCache(scenario.size, DefaultCacheTTL)

			// Préremplir le cache
			delta := NewFactDelta("Test~1", "Test")
			for i := 0; i < scenario.size/2; i++ {
				key := fmt.Sprintf("key%d", i%scenario.keyRange)
				cache.Put(key, delta)
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				key := fmt.Sprintf("key%d", i%scenario.keyRange)
				cache.Get(key)
			}
		})
	}
}

// BenchmarkBatchProcessing compare le traitement batch vs séquentiel.
func BenchmarkBatchProcessing(b *testing.B) {
	// Créer 100 nœuds de types mixtes
	nodes := make([]NodeReference, 100)
	for i := 0; i < 100; i++ {
		nodeType := NodeTypeAlpha
		if i%3 == 1 {
			nodeType = NodeTypeBeta
		} else if i%3 == 2 {
			nodeType = NodeTypeTerminal
		}
		nodes[i] = NodeReference{
			NodeID:   fmt.Sprintf("node%d", i),
			NodeType: nodeType,
		}
	}

	// Fonction de traitement (simule du travail)
	process := func(ref NodeReference) error {
		_ = ref.NodeID
		return nil
	}

	b.Run("Batch", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			batch := NewBatchNodeReferences(100)
			for _, node := range nodes {
				batch.Add(node)
			}
			_ = batch.ProcessInOrder(process)
		}
	})

	b.Run("Sequential", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for _, node := range nodes {
				_ = process(node)
			}
		}
	})
}

// BenchmarkMemoryFootprint mesure l'empreinte mémoire.
func BenchmarkMemoryFootprint(b *testing.B) {
	b.Run("FactDelta_Small", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			delta := NewFactDelta("Test~1", "Test")
			delta.AddFieldChange("field1", "old", "new")
			_ = delta
		}
	})

	b.Run("FactDelta_Large", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			delta := NewFactDelta("Test~1", "Test")
			for j := 0; j < 50; j++ {
				delta.AddFieldChange(fmt.Sprintf("field%d", j), j, j+1)
			}
			_ = delta
		}
	})

	b.Run("Cache_Entries", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			cache := NewOptimizedCache(1000, DefaultCacheTTL)
			delta := NewFactDelta("Test~1", "Test")
			for j := 0; j < 100; j++ {
				cache.Put(fmt.Sprintf("key%d", j), delta)
			}
			_ = cache
		}
	})
}

// BenchmarkScalability teste la scalabilité avec différentes tailles de faits.
func BenchmarkScalability(b *testing.B) {
	sizes := []int{10, 50, 100, 500, 1000}

	detector := NewDeltaDetector()

	for _, size := range sizes {
		b.Run(fmt.Sprintf("FactSize_%d", size), func(b *testing.B) {
			oldFact := make(map[string]interface{}, size)
			newFact := make(map[string]interface{}, size)

			for i := 0; i < size; i++ {
				fieldName := fmt.Sprintf("field%d", i)
				oldFact[fieldName] = i
				newFact[fieldName] = i
			}

			// Modifier 10% des champs
			changeCount := size / 10
			if changeCount == 0 {
				changeCount = 1
			}
			for i := 0; i < changeCount; i++ {
				fieldName := fmt.Sprintf("field%d", i)
				newFact[fieldName] = i + 1000
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				delta, _ := detector.DetectDelta(oldFact, newFact, "Test~1", "Test")
				_ = delta
			}
		})
	}
}

// BenchmarkConcurrentAccess teste l'accès concurrent.
func BenchmarkConcurrentAccess(b *testing.B) {
	cache := NewOptimizedCache(1000, DefaultCacheTTL)
	delta := NewFactDelta("Test~1", "Test")

	// Préremplir
	for i := 0; i < 100; i++ {
		cache.Put(fmt.Sprintf("key%d", i), delta)
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i%100)
			cache.Get(key)
			i++
		}
	})
}
