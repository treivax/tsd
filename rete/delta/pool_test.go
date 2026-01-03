// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
)

func TestPool_FactDelta(t *testing.T) {
	t.Log("🧪 TEST POOL - FactDelta")
	t.Log("========================")

	// Acquérir un FactDelta
	delta := AcquireFactDelta("Test~1", "Test")

	if delta == nil {
		t.Fatal("❌ AcquireFactDelta returned nil")
	}

	if delta.FactID != "Test~1" {
		t.Errorf("❌ Expected FactID 'Test~1', got '%s'", delta.FactID)
	}

	if delta.FactType != "Test" {
		t.Errorf("❌ Expected FactType 'Test', got '%s'", delta.FactType)
	}

	if delta.Fields == nil {
		t.Error("❌ Fields map should be initialized")
	}

	// Ajouter des données
	delta.AddFieldChange("field1", "old", "new")

	// Relâcher
	ReleaseFactDelta(delta)

	// Acquérir à nouveau (devrait être réutilisé)
	delta2 := AcquireFactDelta("Test~2", "Test2")

	if delta2 == nil {
		t.Fatal("❌ Second AcquireFactDelta returned nil")
	}

	// Devrait être nettoyé
	if len(delta2.Fields) != 0 {
		t.Errorf("❌ Fields should be empty after reuse, got %d", len(delta2.Fields))
	}

	if delta2.FactID != "Test~2" {
		t.Errorf("❌ FactID not reset properly: got '%s'", delta2.FactID)
	}

	ReleaseFactDelta(delta2)

	t.Log("✅ Pool FactDelta works correctly")
}

func TestPool_NodeReferenceSlice(t *testing.T) {
	t.Log("🧪 TEST POOL - NodeReference Slice")
	t.Log("==================================")

	// Acquérir une slice
	slice := AcquireNodeReferenceSlice()

	if slice == nil {
		t.Fatal("❌ AcquireNodeReferenceSlice returned nil")
	}

	if len(*slice) != 0 {
		t.Errorf("❌ Expected empty slice, got len=%d", len(*slice))
	}

	// Ajouter des éléments
	*slice = append(*slice, NodeReference{NodeID: "node1"})
	*slice = append(*slice, NodeReference{NodeID: "node2"})

	if len(*slice) != 2 {
		t.Errorf("❌ Expected len=2, got %d", len(*slice))
	}

	// Relâcher
	ReleaseNodeReferenceSlice(slice)

	// Acquérir à nouveau
	slice2 := AcquireNodeReferenceSlice()

	if slice2 == nil {
		t.Fatal("❌ Second acquire returned nil")
	}

	// Devrait être vide
	if len(*slice2) != 0 {
		t.Errorf("❌ Slice should be empty after reuse, got len=%d", len(*slice2))
	}

	ReleaseNodeReferenceSlice(slice2)

	t.Log("✅ Pool NodeReferenceSlice works correctly")
}

func TestPool_StringBuilder(t *testing.T) {
	t.Log("🧪 TEST POOL - StringBuilder")
	t.Log("============================")

	// Acquérir un builder
	sb := AcquireStringBuilder()

	if sb == nil {
		t.Fatal("❌ AcquireStringBuilder returned nil")
	}

	// Utiliser
	sb.WriteString("test")
	result := sb.String()

	if result != "test" {
		t.Errorf("❌ Expected 'test', got '%s'", result)
	}

	// Relâcher
	ReleaseStringBuilder(sb)

	// Acquérir à nouveau
	sb2 := AcquireStringBuilder()

	if sb2 == nil {
		t.Fatal("❌ Second acquire returned nil")
	}

	// Devrait être vide
	if sb2.Len() != 0 {
		t.Errorf("❌ Builder should be empty after reuse, got len=%d", sb2.Len())
	}

	ReleaseStringBuilder(sb2)

	t.Log("✅ Pool StringBuilder works correctly")
}

func TestPool_Map(t *testing.T) {
	t.Log("🧪 TEST POOL - Map")
	t.Log("==================")

	// Acquérir une map
	m := AcquireMap()

	if m == nil {
		t.Fatal("❌ AcquireMap returned nil")
	}

	if len(*m) != 0 {
		t.Errorf("❌ Expected empty map, got len=%d", len(*m))
	}

	// Ajouter des éléments
	(*m)["key1"] = "value1"
	(*m)["key2"] = "value2"

	if len(*m) != 2 {
		t.Errorf("❌ Expected len=2, got %d", len(*m))
	}

	// Relâcher
	ReleaseMap(m)

	// Acquérir à nouveau
	m2 := AcquireMap()

	if m2 == nil {
		t.Fatal("❌ Second acquire returned nil")
	}

	// Devrait être vide
	if len(*m2) != 0 {
		t.Errorf("❌ Map should be empty after reuse, got len=%d", len(*m2))
	}

	ReleaseMap(m2)

	t.Log("✅ Pool Map works correctly")
}

func BenchmarkPool_FactDelta(b *testing.B) {
	b.Run("WithPool", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			delta := AcquireFactDelta("Test~1", "Test")
			delta.AddFieldChange("field1", "old", "new")
			ReleaseFactDelta(delta)
		}
	})

	b.Run("WithoutPool", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			delta := NewFactDelta("Test~1", "Test")
			delta.AddFieldChange("field1", "old", "new")
			// Pas de release - garbage collected
		}
	})
}

func BenchmarkPool_NodeReferenceSlice(b *testing.B) {
	b.Run("WithPool", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			slice := AcquireNodeReferenceSlice()
			*slice = append(*slice, NodeReference{NodeID: "node1"})
			*slice = append(*slice, NodeReference{NodeID: "node2"})
			ReleaseNodeReferenceSlice(slice)
		}
	})

	b.Run("WithoutPool", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			slice := make([]NodeReference, 0, 16)
			slice = append(slice, NodeReference{NodeID: "node1"})
			slice = append(slice, NodeReference{NodeID: "node2"})
			// Utilisation intentionnelle pour éviter warning SA4010
			_ = slice
			// Pas de release - garbage collected
		}
	})
}
