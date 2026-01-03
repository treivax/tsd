// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
	"testing"
)

func BenchmarkDeltaDetector_DetectDelta_NoChanges(b *testing.B) {
	detector := NewDeltaDetector()

	fact := map[string]interface{}{
		"id":       "123",
		"price":    100.0,
		"status":   "active",
		"quantity": 10,
		"category": "Electronics",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.DetectDelta(fact, fact, "Product~123", "Product")
	}
}

func BenchmarkDeltaDetector_DetectDelta_SingleChange(b *testing.B) {
	detector := NewDeltaDetector()

	oldFact := map[string]interface{}{
		"id":       "123",
		"price":    100.0,
		"status":   "active",
		"quantity": 10,
		"category": "Electronics",
	}

	newFact := map[string]interface{}{
		"id":       "123",
		"price":    150.0,
		"status":   "active",
		"quantity": 10,
		"category": "Electronics",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	}
}

func BenchmarkDeltaDetector_DetectDelta_MultipleChanges(b *testing.B) {
	detector := NewDeltaDetector()

	oldFact := map[string]interface{}{
		"id":       "123",
		"price":    100.0,
		"status":   "active",
		"quantity": 10,
		"category": "Electronics",
	}

	newFact := map[string]interface{}{
		"id":       "123",
		"price":    150.0,
		"status":   "inactive",
		"quantity": 5,
		"category": "Electronics",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	}
}

func BenchmarkDeltaDetector_DetectDeltaQuick_NoChanges(b *testing.B) {
	detector := NewDeltaDetector()

	fact := map[string]interface{}{
		"id":       "123",
		"price":    100.0,
		"status":   "active",
		"quantity": 10,
		"category": "Electronics",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.DetectDeltaQuick(fact, fact, "Product~123", "Product")
	}
}

func BenchmarkDeltaDetector_DetectDelta_LargeFact(b *testing.B) {
	detector := NewDeltaDetector()

	oldFact := make(map[string]interface{})
	newFact := make(map[string]interface{})

	for i := 0; i < 50; i++ {
		fieldName := fmt.Sprintf("field%d", i)
		oldFact[fieldName] = i
		newFact[fieldName] = i
	}

	for i := 0; i < 5; i++ {
		fieldName := fmt.Sprintf("field%d", i)
		newFact[fieldName] = i + 1000
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.DetectDelta(oldFact, newFact, "Large~1", "Large")
	}
}

func BenchmarkDeltaDetector_DetectDelta_WithCache(b *testing.B) {
	config := DefaultDetectorConfig()
	config.CacheComparisons = true
	detector := NewDeltaDetectorWithConfig(config)

	oldFact := map[string]interface{}{
		"id":    "123",
		"price": 100.0,
	}

	newFact := map[string]interface{}{
		"id":    "123",
		"price": 150.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.DetectDelta(oldFact, newFact, "Product~123", "Product")
	}
}

func BenchmarkDeltaDetector_DetectDelta_DeepNested(b *testing.B) {
	config := DefaultDetectorConfig()
	config.EnableDeepComparison = true
	config.MaxNestingLevel = 5
	detector := NewDeltaDetectorWithConfig(config)

	oldFact := map[string]interface{}{
		"id": "123",
		"nested": map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"value": "old",
				},
			},
		},
	}

	newFact := map[string]interface{}{
		"id": "123",
		"nested": map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"value": "new",
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.DetectDelta(oldFact, newFact, "Nested~1", "Nested")
	}
}
