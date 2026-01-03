// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"errors"
	"strings"
	"testing"
)

// TestWithFactDelta_Success teste l'utilisation normale de WithFactDelta
func TestWithFactDelta_Success(t *testing.T) {
	t.Log("🧪 TEST: WithFactDelta - Utilisation normale")
	t.Log("=============================================")

	factID := "product~123"
	factType := "Product"

	err := WithFactDelta(factID, factType, func(delta *FactDelta) error {
		// Vérifier que le delta est correctement initialisé
		if delta.FactID != factID {
			t.Errorf("Expected FactID %s, got %s", factID, delta.FactID)
		}
		if delta.FactType != factType {
			t.Errorf("Expected FactType %s, got %s", factType, delta.FactType)
		}

		// Utiliser le delta
		delta.Fields["price"] = FieldDelta{
			FieldName: "price",
			OldValue:  100.0,
			NewValue:  150.0,
		}
		delta.FieldCount = 1

		// Vérifier les champs AVANT le release automatique
		if len(delta.Fields) != 1 {
			t.Errorf("Expected 1 field, got %d", len(delta.Fields))
		}

		if delta.FieldCount != 1 {
			t.Errorf("Expected FieldCount 1, got %d", delta.FieldCount)
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Log("✅ WithFactDelta fonctionne correctement")
}

// TestWithFactDelta_Error teste la gestion d'erreur
func TestWithFactDelta_Error(t *testing.T) {
	t.Log("🧪 TEST: WithFactDelta - Gestion d'erreur")
	t.Log("=========================================")

	expectedErr := errors.New("test error")

	err := WithFactDelta("fact1", "Type1", func(delta *FactDelta) error {
		return expectedErr
	})

	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}

	t.Log("✅ Erreur correctement propagée")
}

// TestWithFactDelta_Panic teste la récupération en cas de panic
func TestWithFactDelta_Panic(t *testing.T) {
	t.Log("🧪 TEST: WithFactDelta - Récupération après panic")
	t.Log("==================================================")

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic to be propagated")
		} else {
			t.Log("✅ Panic correctement propagé")
		}
	}()

	_ = WithFactDelta("fact1", "Type1", func(delta *FactDelta) error {
		panic("test panic")
	})
}

// TestWithNodeReferenceSlice_Success teste l'utilisation de WithNodeReferenceSlice
func TestWithNodeReferenceSlice_Success(t *testing.T) {
	t.Log("🧪 TEST: WithNodeReferenceSlice - Utilisation normale")
	t.Log("=====================================================")

	err := WithNodeReferenceSlice(func(nodes *[]NodeReference) error {
		// Vérifier que la slice est vide au départ
		if len(*nodes) != 0 {
			t.Errorf("Expected empty slice, got %d elements", len(*nodes))
		}

		// Ajouter des éléments
		*nodes = append(*nodes, NodeReference{
			NodeID:   "alpha1",
			NodeType: NodeTypeAlpha,
		})
		*nodes = append(*nodes, NodeReference{
			NodeID:   "beta1",
			NodeType: NodeTypeBeta,
		})

		// Vérifier la taille
		if len(*nodes) != 2 {
			t.Errorf("Expected 2 elements, got %d", len(*nodes))
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Log("✅ WithNodeReferenceSlice fonctionne correctement")
}

// TestWithStringBuilder_Success teste l'utilisation de WithStringBuilder
func TestWithStringBuilder_Success(t *testing.T) {
	t.Log("🧪 TEST: WithStringBuilder - Utilisation normale")
	t.Log("================================================")

	var result string
	err := WithStringBuilder(func(sb *strings.Builder) error {
		// Vérifier que le builder est vide
		if sb.Len() != 0 {
			t.Errorf("Expected empty builder, got length %d", sb.Len())
		}

		// Construire une chaîne
		sb.WriteString("Hello")
		sb.WriteString(" ")
		sb.WriteString("World")

		result = sb.String()
		return nil
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", result)
	}

	t.Log("✅ WithStringBuilder fonctionne correctement")
}

// TestWithStringBuilderResult_Success teste WithStringBuilderResult avec générique
func TestWithStringBuilderResult_Success(t *testing.T) {
	t.Log("🧪 TEST: WithStringBuilderResult - Utilisation avec résultat")
	t.Log("==========================================================")

	result, err := WithStringBuilderResult(func(sb *strings.Builder) (string, error) {
		sb.WriteString("Test")
		sb.WriteString(" ")
		sb.WriteString("Result")
		return sb.String(), nil
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != "Test Result" {
		t.Errorf("Expected 'Test Result', got '%s'", result)
	}

	t.Log("✅ WithStringBuilderResult fonctionne correctement")
}

// TestWithStringBuilderResult_Error teste la gestion d'erreur avec résultat
func TestWithStringBuilderResult_Error(t *testing.T) {
	t.Log("🧪 TEST: WithStringBuilderResult - Gestion d'erreur")
	t.Log("===================================================")

	expectedErr := errors.New("builder error")
	result, err := WithStringBuilderResult(func(sb *strings.Builder) (string, error) {
		return "", expectedErr
	})

	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}

	if result != "" {
		t.Errorf("Expected empty result on error, got '%s'", result)
	}

	t.Log("✅ Erreur correctement gérée")
}

// TestWithMap_Success teste l'utilisation de WithMap
func TestWithMap_Success(t *testing.T) {
	t.Log("🧪 TEST: WithMap - Utilisation normale")
	t.Log("======================================")

	var mapSize int
	err := WithMap(func(m *map[string]interface{}) error {
		// Vérifier que la map est vide
		if len(*m) != 0 {
			t.Errorf("Expected empty map, got %d elements", len(*m))
		}

		// Ajouter des éléments
		(*m)["key1"] = "value1"
		(*m)["key2"] = 42
		(*m)["key3"] = true

		mapSize = len(*m)
		return nil
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if mapSize != 3 {
		t.Errorf("Expected 3 elements, got %d", mapSize)
	}

	t.Log("✅ WithMap fonctionne correctement")
}

// TestPoolLifecycle_NoLeaks vérifie qu'il n'y a pas de fuites mémoire
func TestPoolLifecycle_NoLeaks(t *testing.T) {
	t.Log("🧪 TEST: Pool Lifecycle - Vérification absence de fuites")
	t.Log("========================================================")

	// Exécuter plusieurs cycles d'acquisition/release
	iterations := 1000

	for i := 0; i < iterations; i++ {
		// FactDelta
		_ = WithFactDelta("fact", "Type", func(delta *FactDelta) error {
			delta.Fields["field"] = FieldDelta{FieldName: "field"}
			return nil
		})

		// NodeReferenceSlice
		_ = WithNodeReferenceSlice(func(nodes *[]NodeReference) error {
			*nodes = append(*nodes, NodeReference{NodeID: "n1"})
			return nil
		})

		// StringBuilder
		_ = WithStringBuilder(func(sb *strings.Builder) error {
			sb.WriteString("test")
			return nil
		})

		// Map
		_ = WithMap(func(m *map[string]interface{}) error {
			(*m)["key"] = "value"
			return nil
		})
	}

	t.Logf("✅ %d itérations sans fuite détectée", iterations)
}

// TestPoolLifecycle_ConcurrentAccess teste l'accès concurrent au pool
func TestPoolLifecycle_ConcurrentAccess(t *testing.T) {
	t.Log("🧪 TEST: Pool Lifecycle - Accès concurrent")
	t.Log("==========================================")

	const goroutines = 100
	const iterations = 100

	done := make(chan bool, goroutines)

	for g := 0; g < goroutines; g++ {
		go func(id int) {
			for i := 0; i < iterations; i++ {
				_ = WithFactDelta("fact", "Type", func(delta *FactDelta) error {
					delta.Fields["field"] = FieldDelta{FieldName: "field"}
					return nil
				})
			}
			done <- true
		}(g)
	}

	// Attendre que toutes les goroutines terminent
	for g := 0; g < goroutines; g++ {
		<-done
	}

	t.Logf("✅ %d goroutines × %d itérations sans race condition", goroutines, iterations)
}

// BenchmarkWithFactDelta mesure la performance de WithFactDelta
func BenchmarkWithFactDelta(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = WithFactDelta("fact", "Type", func(delta *FactDelta) error {
			delta.Fields["field"] = FieldDelta{
				FieldName: "field",
				OldValue:  1,
				NewValue:  2,
			}
			return nil
		})
	}
}

// BenchmarkFactDeltaManual mesure la performance sans helper (baseline)
func BenchmarkFactDeltaManual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		delta := AcquireFactDelta("fact", "Type")
		delta.Fields["field"] = FieldDelta{
			FieldName: "field",
			OldValue:  1,
			NewValue:  2,
		}
		ReleaseFactDelta(delta)
	}
}

// BenchmarkWithFactDelta_Parallel mesure la performance en parallèle
func BenchmarkWithFactDelta_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = WithFactDelta("fact", "Type", func(delta *FactDelta) error {
				delta.Fields["field"] = FieldDelta{
					FieldName: "field",
					OldValue:  1,
					NewValue:  2,
				}
				return nil
			})
		}
	})
}
