// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
	"time"
)

// TestIndexBuilder_BuildFromNetwork_NilNetwork tests nil network handling
func TestIndexBuilder_BuildFromNetwork_NilNetwork(t *testing.T) {
	builder := NewIndexBuilder()
	idx, err := builder.BuildFromNetwork(nil)

	// Nil network should return empty index without error
	if err != nil {
		t.Errorf("❌ BuildFromNetwork(nil) should not error, got: %v", err)
	}

	if idx == nil {
		t.Errorf("❌ BuildFromNetwork(nil) should return non-nil index")
	}
}

// TestCacheOptimized_RemoveTail tests the cache LRU tail removal
func TestCacheOptimized_RemoveTail(t *testing.T) {
	cache := NewOptimizedCache(3, 5*time.Minute)

	// Fill cache to capacity
	cache.Put("key1", &FactDelta{FactID: "fact1", FactType: "Type1"})
	cache.Put("key2", &FactDelta{FactID: "fact2", FactType: "Type2"})
	cache.Put("key3", &FactDelta{FactID: "fact3", FactType: "Type3"})

	// Get stats before eviction
	statsBefore := cache.GetStats()
	if statsBefore.Size != 3 {
		t.Errorf("❌ Cache size before = %d, want 3", statsBefore.Size)
	}

	// Add one more - should evict the tail (least recently used)
	cache.Put("key4", &FactDelta{FactID: "fact4", FactType: "Type4"})

	// Verify eviction happened
	statsAfter := cache.GetStats()
	if statsAfter.Size != 3 {
		t.Errorf("❌ Cache size after = %d, want 3", statsAfter.Size)
	}

	if statsAfter.Evictions != 1 {
		t.Errorf("❌ Evictions = %d, want 1", statsAfter.Evictions)
	}

	// Verify oldest entry (key1) was removed
	result, found := cache.Get("key1")
	if found || result != nil {
		t.Errorf("❌ key1 should have been evicted but is still present")
	}

	// Verify newest entry is present
	result, found = cache.Get("key4")
	if !found || result == nil {
		t.Errorf("❌ key4 should be present")
	}
}

// TestCacheOptimized_LRUOrdering tests that LRU ordering is maintained
func TestCacheOptimized_LRUOrdering(t *testing.T) {
	cache := NewOptimizedCache(2, 5*time.Minute)

	// Add two entries
	cache.Put("key1", &FactDelta{FactID: "fact1", FactType: "Type1"})
	cache.Put("key2", &FactDelta{FactID: "fact2", FactType: "Type2"})

	// Access key1 to make it more recently used
	cache.Get("key1")

	// Add key3 - should evict key2 (least recently used)
	cache.Put("key3", &FactDelta{FactID: "fact3", FactType: "Type3"})

	// Verify key2 was evicted
	_, found := cache.Get("key2")
	if found {
		t.Errorf("❌ key2 should have been evicted")
	}

	// Verify key1 and key3 are present
	_, found = cache.Get("key1")
	if !found {
		t.Errorf("❌ key1 should still be present")
	}
	_, found = cache.Get("key3")
	if !found {
		t.Errorf("❌ key3 should be present")
	}
}

// TestDependencyIndex_GetAffectedNodesForDelta_ComplexScenario tests complex delta scenarios
func TestDependencyIndex_GetAffectedNodesForDelta_ComplexScenario(t *testing.T) {
	idx := NewDependencyIndex()

	// Add multiple nodes watching different fields
	idx.AddAlphaNode("price_check", "Product", []string{"price"})
	idx.AddAlphaNode("stock_check", "Product", []string{"stock"})
	idx.AddAlphaNode("price_and_status", "Product", []string{"price", "status"})
	idx.AddBetaNode("order_product", "Product", []string{"price", "category"})
	idx.AddTerminalNode("discount_rule", "Product", []string{"price"})

	t.Run("single field change", func(t *testing.T) {
		delta := NewFactDelta("Product~123", "Product")
		delta.AddFieldChange("price", 100.0, 150.0)

		nodes := idx.GetAffectedNodesForDelta(delta)

		// Should return all nodes that watch "price"
		expectedCount := 4 // price_check, price_and_status, order_product, discount_rule
		if len(nodes) != expectedCount {
			t.Errorf("❌ Got %d affected nodes, want %d", len(nodes), expectedCount)
		}
	})

	t.Run("multiple field changes", func(t *testing.T) {
		delta := NewFactDelta("Product~123", "Product")
		delta.AddFieldChange("price", 100.0, 150.0)
		delta.AddFieldChange("stock", 50, 45)

		nodes := idx.GetAffectedNodesForDelta(delta)

		// Should return all nodes that watch either "price" or "stock"
		expectedCount := 5 // all nodes
		if len(nodes) != expectedCount {
			t.Errorf("❌ Got %d affected nodes, want %d", len(nodes), expectedCount)
		}
	})

	t.Run("field with no watchers", func(t *testing.T) {
		delta := NewFactDelta("Product~123", "Product")
		delta.AddFieldChange("description", "old", "new")

		nodes := idx.GetAffectedNodesForDelta(delta)

		// Should return no nodes
		if len(nodes) != 0 {
			t.Errorf("❌ Got %d affected nodes, want 0", len(nodes))
		}
	})
}

// TestDetectorConfig_ValidateEdgeCases tests additional config validation edge cases
func TestDetectorConfig_ValidateEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		config  DetectorConfig
		wantErr bool
	}{
		{
			name:    "default config is valid",
			config:  DefaultDetectorConfig(),
			wantErr: false,
		},
		{
			name: "valid custom config",
			config: DetectorConfig{
				FloatEpsilon:         0.01,
				IgnoreInternalFields: true,
				TrackTypeChanges:     true,
				EnableDeepComparison: false,
				MaxNestingLevel:      10,
				CacheComparisons:     false,
			},
			wantErr: false,
		},
		{
			name: "negative float epsilon",
			config: DetectorConfig{
				FloatEpsilon:         -0.01,
				IgnoreInternalFields: true,
				MaxNestingLevel:      10,
			},
			wantErr: true,
		},
		{
			name: "zero nesting level",
			config: DetectorConfig{
				FloatEpsilon:    DefaultFloatEpsilon,
				MaxNestingLevel: 0,
			},
			wantErr: true,
		},
		{
			name: "negative nesting level",
			config: DetectorConfig{
				FloatEpsilon:    DefaultFloatEpsilon,
				MaxNestingLevel: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("❌ Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestPropagationConfig_ValidateEdgeCases tests additional propagation config validation
func TestPropagationConfig_ValidateEdgeCases(t *testing.T) {
	t.Run("default config is valid", func(t *testing.T) {
		config := DefaultPropagationConfig()
		err := config.Validate()
		if err != nil {
			t.Errorf("❌ Default config should be valid, got error: %v", err)
		}
	})

	t.Run("config with delta enabled", func(t *testing.T) {
		config := DefaultPropagationConfig()
		config.EnableDeltaPropagation = true
		config.DefaultMode = PropagationModeDelta

		err := config.Validate()
		if err != nil {
			t.Errorf("❌ Config with delta enabled should be valid, got error: %v", err)
		}
	})

	t.Run("config with delta disabled", func(t *testing.T) {
		config := DefaultPropagationConfig()
		config.EnableDeltaPropagation = false
		config.DefaultMode = PropagationModeClassic

		err := config.Validate()
		if err != nil {
			t.Errorf("❌ Config with delta disabled should be valid, got error: %v", err)
		}
	})
}
