// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
	"testing"
)

func TestOptimizedValuesEqual(t *testing.T) {
	t.Log("🧪 TEST OPTIMIZED VALUES EQUAL")
	t.Log("==============================")

	tests := []struct {
		name     string
		a        interface{}
		b        interface{}
		epsilon  float64
		expected bool
	}{
		{"int_equal", 42, 42, 0, true},
		{"int_different", 42, 43, 0, false},
		{"string_equal", "hello", "hello", 0, true},
		{"string_different", "hello", "world", 0, false},
		{"bool_equal", true, true, 0, true},
		{"bool_different", true, false, 0, false},
		{"float64_equal", 3.14, 3.14, 1e-9, true},
		{"float64_close", 3.14, 3.140000001, 1e-6, true},
		{"float64_different", 3.14, 3.15, 1e-9, false},
		{"nil_both", nil, nil, 0, true},
		{"nil_one", nil, 42, 0, false},
		{"different_types", 42, "42", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := OptimizedValuesEqual(tt.a, tt.b, tt.epsilon)
			if result != tt.expected {
				t.Errorf("❌ Expected %v, got %v for %v == %v",
					tt.expected, result, tt.a, tt.b)
			}
		})
	}

	t.Log("✅ OptimizedValuesEqual works correctly")
}

func TestBatchNodeReferences(t *testing.T) {
	t.Log("🧪 TEST BATCH NODE REFERENCES")
	t.Log("=============================")

	batch := NewBatchNodeReferences(12)

	// Ajouter des nœuds de types différents
	batch.Add(NodeReference{NodeID: "alpha1", NodeType: NodeTypeAlpha})
	batch.Add(NodeReference{NodeID: "beta1", NodeType: NodeTypeBeta})
	batch.Add(NodeReference{NodeID: "terminal1", NodeType: NodeTypeTerminal})
	batch.Add(NodeReference{NodeID: "alpha2", NodeType: NodeTypeAlpha})
	batch.Add(NodeReference{NodeID: "beta2", NodeType: NodeTypeBeta})

	if batch.Count() != 5 {
		t.Errorf("❌ Expected count 5, got %d", batch.Count())
	}

	// Traiter dans l'ordre
	var processed []string
	err := batch.ProcessInOrder(func(ref NodeReference) error {
		processed = append(processed, ref.NodeID)
		return nil
	})

	if err != nil {
		t.Errorf("❌ ProcessInOrder failed: %v", err)
	}

	// Vérifier l'ordre: alpha → beta → terminal
	expected := []string{"alpha1", "alpha2", "beta1", "beta2", "terminal1"}
	if len(processed) != len(expected) {
		t.Errorf("❌ Expected %d processed, got %d", len(expected), len(processed))
	}

	for i, id := range expected {
		if processed[i] != id {
			t.Errorf("❌ Expected %s at position %d, got %s", id, i, processed[i])
		}
	}

	t.Log("✅ BatchNodeReferences works correctly")
}

func TestFastHashString(t *testing.T) {
	t.Log("🧪 TEST FAST HASH STRING")
	t.Log("========================")

	// Test que les hash sont cohérents
	s1 := "test"
	h1 := FastHashString(s1)
	h2 := FastHashString(s1)

	if h1 != h2 {
		t.Error("❌ Hash should be consistent for same string")
	}

	// Test que des strings différentes donnent des hash différents
	s2 := "different"
	h3 := FastHashString(s2)

	if h1 == h3 {
		t.Error("❌ Different strings should have different hashes (collision unlikely)")
	}

	// Test string vide
	h4 := FastHashString("")
	if h4 == 0 {
		t.Error("❌ Empty string should have non-zero hash")
	}

	t.Log("✅ FastHashString works correctly")
}

func TestCopyFactFast(t *testing.T) {
	t.Log("🧪 TEST COPY FACT FAST")
	t.Log("======================")

	original := map[string]interface{}{
		"field1": "value1",
		"field2": 42,
		"field3": true,
	}

	copied := CopyFactFast(original)

	// Vérifier égalité
	if len(copied) != len(original) {
		t.Errorf("❌ Expected len %d, got %d", len(original), len(copied))
	}

	for k, v := range original {
		cv, exists := copied[k]
		if !exists {
			t.Errorf("❌ Key %s missing in copy", k)
		}
		if cv != v {
			t.Errorf("❌ Value mismatch for key %s", k)
		}
	}

	// Vérifier indépendance (shallow copy)
	copied["new"] = "added"
	if _, exists := original["new"]; exists {
		t.Error("❌ Copy should be independent")
	}

	t.Log("✅ CopyFactFast works correctly")
}

// Benchmarks comparatifs

func BenchmarkValuesEqual_Optimized_vs_Standard(b *testing.B) {
	testCases := []struct {
		name string
		a    interface{}
		b    interface{}
	}{
		{"int", 42, 42},
		{"string", "hello", "hello"},
		{"float64", 3.14, 3.14},
		{"bool", true, true},
	}

	for _, tc := range testCases {
		b.Run(tc.name+"_optimized", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				OptimizedValuesEqual(tc.a, tc.b, DefaultFloatEpsilon)
			}
		})

		b.Run(tc.name+"_standard", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ValuesEqual(tc.a, tc.b, DefaultFloatEpsilon)
			}
		})
	}
}

func BenchmarkCopyFact(b *testing.B) {
	fact := make(map[string]interface{})
	for i := 0; i < 50; i++ {
		fact[fmt.Sprintf("field%d", i)] = i
	}

	b.Run("CopyFactFast", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = CopyFactFast(fact)
		}
	})

	b.Run("ManualCopy", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			copy := make(map[string]interface{})
			for k, v := range fact {
				copy[k] = v
			}
		}
	})
}

func BenchmarkBatchNodeReferences(b *testing.B) {
	// Préparer des nœuds
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

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		batch := NewBatchNodeReferences(100)
		for _, node := range nodes {
			batch.Add(node)
		}
		_ = batch.ProcessInOrder(func(ref NodeReference) error {
			return nil
		})
	}
}

func BenchmarkFastHashString(b *testing.B) {
	s := "This is a test string for hashing performance"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = FastHashString(s)
	}
}
