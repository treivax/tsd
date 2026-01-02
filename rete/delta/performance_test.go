// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
	"testing"
)

// BenchmarkDeltaDetector_Various mesure la performance de détection selon le contexte
func BenchmarkDeltaDetector_Various(b *testing.B) {
	fixtures := NewTestFixtures()

	benchmarks := []struct {
		name    string
		oldFact map[string]interface{}
		newFact map[string]interface{}
	}{
		{
			"NoChange_Simple",
			fixtures.SimpleFact,
			fixtures.SimpleFact,
		},
		{
			"SingleChange_Simple",
			fixtures.SimpleFact,
			ModifyField(fixtures.SimpleFact, "price", 200.0),
		},
		{
			"MultiChange_Complex",
			fixtures.ComplexFact,
			ModifyFields(fixtures.ComplexFact, map[string]interface{}{
				"price":    300.0,
				"quantity": 20,
			}),
		},
		{
			"AllChange_Large",
			fixtures.LargeFact,
			ModifyAllFields(fixtures.LargeFact),
		},
		{
			"NestedChange",
			fixtures.NestedFact,
			ModifyField(fixtures.NestedFact, "address", map[string]interface{}{
				"street": "456 Other St",
				"city":   "Lyon",
			}),
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			detector := NewDeltaDetector()
			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_, _ = detector.DetectDelta(bm.oldFact, bm.newFact, "Test~1", "Test")
			}
		})
	}
}

// BenchmarkDeltaDetector_WithCache mesure performance avec cache activé
func BenchmarkDeltaDetector_WithCache(b *testing.B) {
	config := DefaultDetectorConfig()
	config.CacheComparisons = true
	detector := NewDeltaDetectorWithConfig(config)

	oldFact := map[string]interface{}{"price": 100.0, "status": "active"}
	newFact := map[string]interface{}{"price": 150.0, "status": "active"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	}
}

// BenchmarkDependencyIndex_Scaling mesure la scalabilité de l'index
func BenchmarkDependencyIndex_Scaling(b *testing.B) {
	sizes := []int{10, 100, 1000, TestBenchmarkSize}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			idx := NewDependencyIndex()

			// Setup index
			for i := 0; i < size; i++ {
				idx.AddAlphaNode(fmt.Sprintf("node%d", i), "Product", []string{"price", "status"})
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = idx.GetAffectedNodes("Product", "price")
			}
		})
	}
}

// BenchmarkDependencyIndex_AddNode mesure performance d'ajout de nœuds
func BenchmarkDependencyIndex_AddNode(b *testing.B) {
	idx := NewDependencyIndex()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		nodeID := fmt.Sprintf("node_%d", i)
		idx.AddAlphaNode(nodeID, "Product", []string{"price", "status", "quantity"})
	}
}

// BenchmarkPropagation_EndToEnd mesure la performance end-to-end
func BenchmarkPropagation_EndToEnd(b *testing.B) {
	helper := NewBenchmarkHelper()
	fixtures := NewTestFixtures()

	// Setup index avec nœuds
	for i := 0; i < 50; i++ {
		helper.Index.AddAlphaNode(fmt.Sprintf("alpha%d", i), "Product", []string{"price"})
	}

	oldFact := fixtures.SimpleFact
	newFact := ModifyField(fixtures.SimpleFact, "price", 200.0)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = helper.Propagator.PropagateUpdate(oldFact, newFact, "Product~123", "Product")
	}
}

// BenchmarkPropagation_DifferentScales mesure performance selon taille des changements
func BenchmarkPropagation_DifferentScales(b *testing.B) {
	scales := []struct {
		name      string
		nodeCount int
		changes   int
	}{
		{"Small_10nodes_1change", 10, 1},
		{"Medium_50nodes_5changes", 50, 5},
		{"Large_100nodes_10changes", 100, 10},
	}

	for _, scale := range scales {
		b.Run(scale.name, func(b *testing.B) {
			helper := NewBenchmarkHelper()

			// Setup nodes
			for i := 0; i < scale.nodeCount; i++ {
				helper.Index.AddAlphaNode(
					fmt.Sprintf("node_%d", i),
					"Product",
					[]string{"price", "status", "quantity"},
				)
			}

			// Créer fact avec changements
			oldFact := generateLargeFact(20)
			newFact := make(map[string]interface{})
			for k, v := range oldFact {
				newFact[k] = v
			}
			// Modifier N champs
			for i := 0; i < scale.changes; i++ {
				fieldName := fmt.Sprintf("field_%03d", i)
				newFact[fieldName] = i + 1000
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = helper.Propagator.PropagateUpdate(oldFact, newFact, "Test~1", "Test")
			}
		})
	}
}

// BenchmarkComparison_ValuesEqual mesure performance de comparaison de valeurs
func BenchmarkComparison_ValuesEqual(b *testing.B) {
	tests := []struct {
		name string
		a    interface{}
		b    interface{}
	}{
		{"int_equal", 42, 42},
		{"int_different", 42, 43},
		{"string_equal", "test string", "test string"},
		{"string_different", "test string", "other string"},
		{"float_equal", 100.5, 100.5},
		{"float_different", 100.5, 100.6},
		{
			"map_equal",
			map[string]interface{}{"a": 1, "b": 2},
			map[string]interface{}{"a": 1, "b": 2},
		},
		{
			"map_different",
			map[string]interface{}{"a": 1, "b": 2},
			map[string]interface{}{"a": 1, "b": 3},
		},
		{
			"slice_equal",
			[]interface{}{1, 2, 3},
			[]interface{}{1, 2, 3},
		},
		{
			"slice_different",
			[]interface{}{1, 2, 3},
			[]interface{}{1, 2, 4},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = ValuesEqual(tt.a, tt.b, DefaultFloatEpsilon)
			}
		})
	}
}

// BenchmarkDeltaPropagator_Mode mesure performance selon mode de propagation
func BenchmarkDeltaPropagator_Mode(b *testing.B) {
	modes := []PropagationMode{
		PropagationModeAuto,
		PropagationModeDelta,
		PropagationModeClassic,
	}

	for _, mode := range modes {
		b.Run(mode.String(), func(b *testing.B) {
			helper := NewBenchmarkHelper()

			config := DefaultPropagationConfig()
			config.DefaultMode = mode
			_ = helper.Propagator.UpdateConfig(config)

			// Setup nodes
			for i := 0; i < 50; i++ {
				helper.Index.AddAlphaNode(fmt.Sprintf("node_%d", i), "Product", []string{"price"})
			}

			oldFact := map[string]interface{}{"id": "123", "price": 100.0}
			newFact := map[string]interface{}{"id": "123", "price": 150.0}

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = helper.Propagator.PropagateUpdate(oldFact, newFact, "Product~123", "Product")
			}
		})
	}
}

// BenchmarkMemoryAllocation mesure allocations mémoire
func BenchmarkMemoryAllocation(b *testing.B) {
	b.Run("DeltaDetector_creation", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = NewDeltaDetector()
		}
	})

	b.Run("DependencyIndex_creation", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = NewDependencyIndex()
		}
	})

	b.Run("FactDelta_creation", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			delta := NewFactDelta("Test~1", "Test")
			delta.AddFieldChange("field1", 1, 2)
			delta.AddFieldChange("field2", "a", "b")
		}
	})
}
