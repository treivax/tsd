// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
	"testing"
)

const (
	// Constantes pour tests
	TestFieldCount    = 100
	TestLargeStringMB = 1
	TestConcurrentOps = 100
	TestStressRuns    = 1000
	TestBenchmarkSize = 10000
)

// TestFixtures contient des données de test réutilisables.
type TestFixtures struct {
	SimpleFact  map[string]interface{}
	ComplexFact map[string]interface{}
	LargeFact   map[string]interface{}
	NestedFact  map[string]interface{}
	SliceFact   map[string]interface{}
	UnicodeFact map[string]interface{}
}

// NewTestFixtures crée un ensemble de fixtures de test.
func NewTestFixtures() *TestFixtures {
	return &TestFixtures{
		SimpleFact: map[string]interface{}{
			"id":     "123",
			"name":   "Test Product",
			"price":  100.0,
			"status": "active",
		},
		ComplexFact: map[string]interface{}{
			"id":       "456",
			"name":     "Complex Product",
			"price":    200.0,
			"quantity": 10,
			"category": "Electronics",
			"tags":     []interface{}{"new", "featured"},
			"metadata": map[string]interface{}{
				"brand": "TestBrand",
				"model": "XYZ",
			},
		},
		LargeFact: generateLargeFact(TestFieldCount),
		NestedFact: map[string]interface{}{
			"id": "789",
			"address": map[string]interface{}{
				"street": "123 Main St",
				"city":   "Paris",
				"country": map[string]interface{}{
					"name": "France",
					"code": "FR",
				},
			},
		},
		SliceFact: map[string]interface{}{
			"id":    "slice_1",
			"items": []interface{}{1, 2, 3, 4, 5},
			"tags":  []interface{}{"tag1", "tag2"},
		},
		UnicodeFact: map[string]interface{}{
			"id":    "unicode_1",
			"名前":    "古い値",
			"emoji": "😀",
			"text":  "Voilà un café à Zürich",
		},
	}
}

// generateLargeFact génère un fait avec N champs.
func generateLargeFact(fieldCount int) map[string]interface{} {
	fact := make(map[string]interface{})
	fact["id"] = "large_fact_1"
	for i := 0; i < fieldCount; i++ {
		fieldName := fmt.Sprintf("field_%03d", i)
		fact[fieldName] = i
	}
	return fact
}

// GenerateDeepNestedFact crée structure imbriquée sur N niveaux.
func GenerateDeepNestedFact(depth int, leafValue string) map[string]interface{} {
	if depth == 0 {
		return map[string]interface{}{"value": leafValue}
	}
	return map[string]interface{}{
		"level":  depth,
		"nested": GenerateDeepNestedFact(depth-1, leafValue),
	}
}

// ModifyField crée copie d'un fait avec un champ modifié.
func ModifyField(fact map[string]interface{}, field string, value interface{}) map[string]interface{} {
	modified := make(map[string]interface{})
	for k, v := range fact {
		modified[k] = v
	}
	modified[field] = value
	return modified
}

// ModifyFields crée copie d'un fait avec plusieurs champs modifiés.
func ModifyFields(fact map[string]interface{}, changes map[string]interface{}) map[string]interface{} {
	modified := make(map[string]interface{})
	for k, v := range fact {
		modified[k] = v
	}
	for k, v := range changes {
		modified[k] = v
	}
	return modified
}

// ModifyAllFields modifie tous les champs d'un fait.
func ModifyAllFields(fact map[string]interface{}) map[string]interface{} {
	modified := make(map[string]interface{})
	for k := range fact {
		modified[k] = "modified"
	}
	return modified
}

// AssertDeltaEquals vérifie qu'un delta correspond aux attentes.
func AssertDeltaEquals(t *testing.T, delta *FactDelta, expectedFields map[string]struct {
	oldVal interface{}
	newVal interface{}
}) {
	t.Helper()

	if len(delta.Fields) != len(expectedFields) {
		t.Errorf("❌ Expected %d changed fields, got %d", len(expectedFields), len(delta.Fields))
		return
	}

	for fieldName, expected := range expectedFields {
		fieldDelta, exists := delta.Fields[fieldName]
		if !exists {
			t.Errorf("❌ Expected field '%s' in delta, not found", fieldName)
			continue
		}

		if fieldDelta.OldValue != expected.oldVal {
			t.Errorf("❌ Field '%s': expected old value %v, got %v", fieldName, expected.oldVal, fieldDelta.OldValue)
		}

		if fieldDelta.NewValue != expected.newVal {
			t.Errorf("❌ Field '%s': expected new value %v, got %v", fieldName, expected.newVal, fieldDelta.NewValue)
		}
	}
}

// AssertNodesContain vérifie qu'une liste de nœuds contient les IDs attendus.
func AssertNodesContain(t *testing.T, nodes []NodeReference, expectedIDs []string) {
	t.Helper()

	actualIDs := make(map[string]bool)
	for _, node := range nodes {
		actualIDs[node.NodeID] = true
	}

	for _, expectedID := range expectedIDs {
		if !actualIDs[expectedID] {
			t.Errorf("❌ Expected node '%s' in result, not found", expectedID)
		}
	}
}

// AssertContainsAll vérifie qu'une slice contient toutes les valeurs attendues.
func AssertContainsAll(t *testing.T, actual []string, expected []string) {
	t.Helper()

	actualMap := make(map[string]bool)
	for _, s := range actual {
		actualMap[s] = true
	}

	for _, exp := range expected {
		if !actualMap[exp] {
			t.Errorf("❌ Expected '%s' in result, not found. Got: %v", exp, actual)
		}
	}
}

// BenchmarkHelper facilite l'écriture de benchmarks.
type BenchmarkHelper struct {
	Detector   *DeltaDetector
	Index      *DependencyIndex
	Propagator *DeltaPropagator
}

// NewBenchmarkHelper crée un helper pour benchmarks.
func NewBenchmarkHelper() *BenchmarkHelper {
	index := NewDependencyIndex()
	detector := NewDeltaDetector()

	propagator, _ := NewDeltaPropagatorBuilder().
		WithIndex(index).
		WithDetector(detector).
		Build()

	return &BenchmarkHelper{
		Detector:   detector,
		Index:      index,
		Propagator: propagator,
	}
}

// GenerateLargeString génère une chaîne de N mégaoctets.
func GenerateLargeString(sizeMB int) string {
	size := sizeMB * 1024 * 1024
	bytes := make([]byte, size)
	for i := range bytes {
		bytes[i] = 'a' + byte(i%26)
	}
	return string(bytes)
}

// CreateMultipleNodes crée plusieurs nœuds dans l'index pour tests.
func CreateMultipleNodes(index *DependencyIndex, nodeCount int, factType string, fields []string) {
	for i := 0; i < nodeCount; i++ {
		nodeID := fmt.Sprintf("node_%d", i)
		index.AddAlphaNode(nodeID, factType, fields)
	}
}
