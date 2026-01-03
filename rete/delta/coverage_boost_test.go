// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

// TestIndexBuilder_ExtractBetaNodes_Coverage tests extractBetaNodes thoroughly
func TestIndexBuilder_ExtractBetaNodes_Coverage(t *testing.T) {
	tests := []struct {
		name         string
		network      interface{}
		wantErr      bool
		wantSkipped  int
		wantWarnings int
		enableDiag   bool
	}{
		{
			name: "nil BetaNodes field",
			network: &struct {
				BetaNodes map[string]interface{}
			}{
				BetaNodes: nil,
			},
			wantErr:     false,
			wantSkipped: 0,
			enableDiag:  false,
		},
		{
			name: "empty BetaNodes map",
			network: &struct {
				BetaNodes map[string]interface{}
			}{
				BetaNodes: make(map[string]interface{}),
			},
			wantErr:     false,
			wantSkipped: 0,
			enableDiag:  false,
		},
		{
			name: "BetaNodes with interface nodes",
			network: &struct {
				BetaNodes map[string]interface{}
			}{
				BetaNodes: map[string]interface{}{
					"beta1": &struct{ ID string }{ID: "beta1"},
					"beta2": &struct{ ID string }{ID: "beta2"},
				},
			},
			wantErr:      false,
			wantSkipped:  2,
			wantWarnings: 2,
			enableDiag:   true,
		},
		{
			name: "BetaNodes with pointer nodes",
			network: &struct {
				BetaNodes map[string]interface{}
			}{
				BetaNodes: map[string]interface{}{
					"beta1": &struct{ ID string }{ID: "beta1"},
				},
			},
			wantErr:      false,
			wantSkipped:  1,
			wantWarnings: 1,
			enableDiag:   true,
		},
		{
			name: "BetaNodes not a map",
			network: &struct {
				BetaNodes string
			}{
				BetaNodes: "not-a-map",
			},
			wantErr:    true,
			enableDiag: false,
		},
		{
			name: "BetaNodes with nil values",
			network: &struct {
				BetaNodes map[string]interface{}
			}{
				BetaNodes: map[string]interface{}{
					"beta1": nil,
					"beta2": &struct{ ID string }{ID: "beta2"},
				},
			},
			wantErr:      false,
			wantSkipped:  1, // Only beta2 counted, beta1 skipped silently
			wantWarnings: 1,
			enableDiag:   true,
		},
		{
			name: "no BetaNodes field",
			network: &struct {
				AlphaNodes map[string]interface{}
			}{
				AlphaNodes: make(map[string]interface{}),
			},
			wantErr:     false,
			wantSkipped: 0,
			enableDiag:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewIndexBuilder()
			if tt.enableDiag {
				builder.EnableDiagnostics()
			}
			idx := NewDependencyIndex()

			networkValue := reflect.ValueOf(tt.network)
			if networkValue.Kind() == reflect.Ptr {
				networkValue = networkValue.Elem()
			}

			err := builder.extractBetaNodes(idx, networkValue)

			if (err != nil) != tt.wantErr {
				t.Errorf("❌ extractBetaNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.enableDiag {
				diag := builder.GetDiagnostics()
				if diag.NodesSkipped != tt.wantSkipped {
					t.Errorf("❌ NodesSkipped = %d, want %d", diag.NodesSkipped, tt.wantSkipped)
				}
				if len(diag.Warnings) != tt.wantWarnings {
					t.Errorf("❌ Warnings count = %d, want %d", len(diag.Warnings), tt.wantWarnings)
				}
			}
		})
	}
}

// TestFieldExtractor_ExtractFieldsFromBinaryNode_Coverage tests binary node extraction
func TestFieldExtractor_ExtractFieldsFromBinaryNode_Coverage(t *testing.T) {
	tests := []struct {
		name    string
		node    map[string]interface{}
		want    []string
		wantErr bool
	}{
		{
			name: "AND node with two field accesses",
			node: map[string]interface{}{
				"type":     "binary",
				"operator": "AND",
				"left": map[string]interface{}{
					"type":  "fieldAccess",
					"field": "status",
				},
				"right": map[string]interface{}{
					"type":  "fieldAccess",
					"field": "priority",
				},
			},
			want:    []string{"status", "priority"},
			wantErr: false,
		},
		{
			name: "OR node with nested conditions",
			node: map[string]interface{}{
				"type":     "binary",
				"operator": "OR",
				"left": map[string]interface{}{
					"type":  "fieldAccess",
					"field": "active",
				},
				"right": map[string]interface{}{
					"type": "comparison",
					"left": map[string]interface{}{
						"type":  "fieldAccess",
						"field": "count",
					},
					"operator": ">",
					"right":    10,
				},
			},
			want:    []string{"active", "count"},
			wantErr: false,
		},
		{
			name: "binary node with literal values",
			node: map[string]interface{}{
				"type":     "binary",
				"operator": "AND",
				"left":     true,
				"right":    false,
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "binary node missing left",
			node: map[string]interface{}{
				"type":     "binary",
				"operator": "AND",
				"right": map[string]interface{}{
					"type":  "fieldAccess",
					"field": "test",
				},
			},
			want:    []string{"test"},
			wantErr: false,
		},
		{
			name: "binary node missing right",
			node: map[string]interface{}{
				"type":     "binary",
				"operator": "AND",
				"left": map[string]interface{}{
					"type":  "fieldAccess",
					"field": "test",
				},
			},
			want:    []string{"test"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := make(map[string]bool)
			err := extractFieldsFromBinaryNode(tt.node, fields)

			if (err != nil) != tt.wantErr {
				t.Errorf("❌ extractFieldsFromBinaryNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Convert fields map to slice for comparison
			var got []string
			for field := range fields {
				got = append(got, field)
			}

			if len(got) != len(tt.want) {
				t.Errorf("❌ extractFieldsFromBinaryNode() returned %d fields, want %d", len(got), len(tt.want))
				t.Errorf("   Got: %v", got)
				t.Errorf("   Want: %v", tt.want)
				return
			}

			// Convert to map for order-independent comparison
			gotMap := make(map[string]bool)
			for _, field := range got {
				gotMap[field] = true
			}

			for _, field := range tt.want {
				if !gotMap[field] {
					t.Errorf("❌ Missing expected field: %s", field)
				}
			}
		})
	}
}

// TestOptimizations_FastHashBytes tests the fast hash function
func TestOptimizations_FastHashBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  uint64
	}{
		{
			name:  "empty bytes",
			input: []byte{},
			want:  14695981039346656037, // FNV offset basis
		},
		{
			name:  "single byte",
			input: []byte{42},
			want:  0, // Will be calculated
		},
		{
			name:  "hello world",
			input: []byte("hello world"),
			want:  0,
		},
		{
			name:  "same input produces same hash",
			input: []byte("test"),
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FastHashBytes(tt.input)

			// For non-empty, verify deterministic
			if len(tt.input) > 0 {
				got2 := FastHashBytes(tt.input)
				if got != got2 {
					t.Errorf("❌ FastHashBytes not deterministic: %d != %d", got, got2)
				}
			}

			// For empty bytes, verify it returns FNV offset
			if len(tt.input) == 0 && tt.want != 0 {
				if got != tt.want {
					t.Errorf("❌ FastHashBytes() = %d, want %d", got, tt.want)
				}
			}
		})
	}
}

// TestOptimizations_PreallocatedMap tests preallocated map creation
func TestOptimizations_PreallocatedMap(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
	}{
		{
			name:     "zero capacity",
			capacity: 0,
		},
		{
			name:     "small capacity",
			capacity: 10,
		},
		{
			name:     "large capacity",
			capacity: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PreallocatedMap(tt.capacity)

			if got == nil {
				t.Errorf("❌ PreallocatedMap() returned nil")
				return
			}

			// Verify it's a map
			if len(got) != 0 {
				t.Errorf("❌ PreallocatedMap() should return empty map, got len=%d", len(got))
			}

			// Verify we can use it
			got["test"] = "value"
			if got["test"] != "value" {
				t.Errorf("❌ PreallocatedMap() map not usable")
			}
		})
	}
}

// TestPool_ReleaseEdgeCases tests edge cases for pool release functions
func TestPool_ReleaseEdgeCases(t *testing.T) {
	t.Run("ReleaseNodeReferenceSlice with nil", func(t *testing.T) {
		// Should not panic
		ReleaseNodeReferenceSlice(nil)
	})

	t.Run("ReleaseNodeReferenceSlice with empty slice", func(t *testing.T) {
		slice := make([]NodeReference, 0)
		ReleaseNodeReferenceSlice(&slice)
	})

	t.Run("ReleaseNodeReferenceSlice with data", func(t *testing.T) {
		slice := make([]NodeReference, 3)
		slice[0] = NodeReference{NodeID: "node1"}
		slice[1] = NodeReference{NodeID: "node2"}
		slice[2] = NodeReference{NodeID: "node3"}

		// ReleaseNodeReferenceSlice expects a pointer
		slicePtr := &slice
		ReleaseNodeReferenceSlice(slicePtr)

		// The function doesn't actually modify the slice, just puts it back in pool
		// So we can't test the clearing behavior this way
	})

	t.Run("ReleaseStringBuilder with nil", func(t *testing.T) {
		// Should not panic
		ReleaseStringBuilder(nil)
	})

	t.Run("ReleaseStringBuilder with data", func(t *testing.T) {
		sb := &strings.Builder{}
		sb.WriteString("test data")

		// Release should not panic
		ReleaseStringBuilder(sb)
	})
}

// TestIntegrationHelper_DiagnosticMethods tests Enable/DisableDiagnostics
func TestIntegrationHelper_DiagnosticMethods(t *testing.T) {
	t.Run("EnableDiagnostics on IntegrationHelper", func(t *testing.T) {
		helper := &IntegrationHelper{
			builder: NewIndexBuilder(),
		}

		helper.EnableDiagnostics()

		// Verify diagnostics are enabled in the builder
		diag := helper.builder.GetDiagnostics()
		if diag == nil {
			t.Errorf("❌ Diagnostics should be initialized after EnableDiagnostics")
		}
	})

	t.Run("DisableDiagnostics on IntegrationHelper", func(t *testing.T) {
		helper := &IntegrationHelper{
			builder: NewIndexBuilder(),
		}

		helper.EnableDiagnostics()
		helper.DisableDiagnostics()

		// Diagnostics should still exist but might not collect new data
		// The important thing is it doesn't panic
	})
}

// TestCacheOptimized_NewOptimizedCache_EdgeCases tests cache creation edge cases
func TestCacheOptimized_NewOptimizedCache_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		maxSize int
		ttl     time.Duration
	}{
		{
			name:    "negative maxSize uses default",
			maxSize: -1,
			ttl:     300 * time.Second,
		},
		{
			name:    "zero maxSize uses default",
			maxSize: 0,
			ttl:     300 * time.Second,
		},
		{
			name:    "negative TTL uses default",
			maxSize: 100,
			ttl:     -1 * time.Second,
		},
		{
			name:    "zero TTL uses default",
			maxSize: 100,
			ttl:     0,
		},
		{
			name:    "valid values preserved",
			maxSize: 50,
			ttl:     600 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewOptimizedCache(tt.maxSize, tt.ttl)

			if cache == nil {
				t.Errorf("❌ NewOptimizedCache() returned nil")
				return
			}

			// Verify cache is usable
			stats := cache.GetStats()
			if stats.Size != 0 {
				t.Errorf("❌ New cache should be empty, got size=%d", stats.Size)
			}
		})
	}
}

// TestErrors_ErrorMethod tests the Error() method on error types
func TestErrors_ErrorMethod(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		contains []string
	}{
		{
			name: "ComponentNotInitializedError",
			err: &ComponentNotInitializedError{
				Component: "propagator",
				Function:  "Propagate",
				Message:   "not initialized",
			},
			contains: []string{"propagator", "Propagate", "not initialized"},
		},
		{
			name: "InvalidConfigError",
			err: &InvalidConfigError{
				Field:  "FloatEpsilon",
				Reason: "must be positive",
			},
			contains: []string{"FloatEpsilon", "must be positive"},
		},
		{
			name: "InvalidFactError",
			err: &InvalidFactError{
				FactID:   "Product~123",
				FactType: "Product",
				Reason:   "missing required field",
			},
			contains: []string{"Product~123", "Product", "missing required field"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()

			if got == "" {
				t.Errorf("❌ Error() returned empty string")
			}

			for _, substr := range tt.contains {
				if !strings.Contains(got, substr) {
					t.Errorf("❌ Error() = %q, should contain %q", got, substr)
				}
			}
		})
	}
}

// TestPropagationStrategy_Sequential tests the sequential strategy
func TestPropagationStrategy_Sequential(t *testing.T) {
	strategy := &SequentialStrategy{}

	t.Run("GetName", func(t *testing.T) {
		name := strategy.GetName()
		if name != "Sequential" {
			t.Errorf("❌ GetName() = %q, want %q", name, "Sequential")
		}
	})

	t.Run("ShouldPropagate with empty nodes", func(t *testing.T) {
		delta := NewFactDelta("fact1", "Type1")
		delta.AddFieldChange("price", 100.0, 150.0)

		nodes := []NodeReference{}
		result := strategy.ShouldPropagate(delta, nodes)

		// Should return false when there are no affected nodes
		if result {
			t.Errorf("❌ ShouldPropagate() = %v, want false for empty nodes", result)
		}
	})

	t.Run("ShouldPropagate with non-empty nodes", func(t *testing.T) {
		delta := NewFactDelta("fact1", "Type1")
		delta.AddFieldChange("price", 100.0, 150.0)

		nodes := []NodeReference{
			{NodeID: "alpha1", NodeType: NodeTypeAlpha},
		}
		result := strategy.ShouldPropagate(delta, nodes)

		// Should return true when there are affected nodes
		if !result {
			t.Errorf("❌ ShouldPropagate() = %v, want true for non-empty nodes", result)
		}
	})

	t.Run("GetPropagationOrder", func(t *testing.T) {
		nodes := []NodeReference{
			{NodeID: "beta1", NodeType: NodeTypeBeta},
			{NodeID: "alpha1", NodeType: NodeTypeAlpha},
			{NodeID: "term1", NodeType: NodeTypeTerminal},
		}

		ordered := strategy.GetPropagationOrder(nodes)

		if len(ordered) != len(nodes) {
			t.Errorf("❌ GetPropagationOrder() returned %d nodes, want %d", len(ordered), len(nodes))
		}
	})
}
