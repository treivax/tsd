// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"strings"
	"sync"
	"testing"
)

// TestComputeHashDirect tests the computeHashDirect function with various scenarios
func TestComputeHashDirect(t *testing.T) {
	tests := []struct {
		name        string
		config      BetaSharingConfig
		signature   *JoinNodeSignature
		expectError bool
		description string
	}{
		{
			name: "basic signature with metrics enabled",
			config: BetaSharingConfig{
				Enabled:                     true,
				EnableMetrics:               true,
				NormalizeOrder:              true,
				HashCacheSize:               100,
				MaxSharedNodes:              10,
				EnableAdvancedNormalization: false,
			},
			signature: &JoinNodeSignature{
				LeftVars:  []string{"x", "y"},
				RightVars: []string{"z"},
				AllVars:   []string{"x", "y", "z"},
				Condition: map[string]interface{}{
					"operator": "equals",
					"field":    "value",
				},
				VarTypes: map[string]string{
					"x": "Person",
					"y": "Person",
					"z": "Order",
				},
			},
			expectError: false,
			description: "should compute hash with metrics",
		},
		{
			name: "signature with metrics disabled",
			config: BetaSharingConfig{
				Enabled:                     true,
				EnableMetrics:               false,
				NormalizeOrder:              false,
				HashCacheSize:               50,
				MaxSharedNodes:              5,
				EnableAdvancedNormalization: false,
			},
			signature: &JoinNodeSignature{
				LeftVars:  []string{"a"},
				RightVars: []string{"b"},
				AllVars:   []string{"a", "b"},
				Condition: map[string]interface{}{
					"type": "comparison",
				},
				VarTypes: map[string]string{
					"a": "TypeA",
					"b": "TypeB",
				},
			},
			expectError: false,
			description: "should compute hash without metrics",
		},
		{
			name: "empty signature",
			config: BetaSharingConfig{
				Enabled:        true,
				NormalizeOrder: true,
				HashCacheSize:  10,
			},
			signature: &JoinNodeSignature{
				LeftVars:  []string{},
				RightVars: []string{},
				AllVars:   []string{},
				Condition: map[string]interface{}{},
				VarTypes:  map[string]string{},
			},
			expectError: false,
			description: "should handle empty signature",
		},
		{
			name: "complex nested condition",
			config: BetaSharingConfig{
				Enabled:                     true,
				NormalizeOrder:              true,
				EnableAdvancedNormalization: true,
				HashCacheSize:               200,
			},
			signature: &JoinNodeSignature{
				LeftVars:  []string{"x", "y", "z"},
				RightVars: []string{"a", "b", "c"},
				AllVars:   []string{"x", "y", "z", "a", "b", "c"},
				Condition: map[string]interface{}{
					"type": "and",
					"conditions": []interface{}{
						map[string]interface{}{
							"operator": "equals",
							"left":     "field1",
							"right":    "field2",
						},
						map[string]interface{}{
							"operator": "greater_than",
							"left":     "field3",
							"right":    100,
						},
					},
				},
				VarTypes: map[string]string{
					"x": "Type1",
					"y": "Type2",
					"z": "Type3",
					"a": "Type4",
					"b": "Type5",
					"c": "Type6",
				},
			},
			expectError: false,
			description: "should handle complex nested conditions",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lifecycle := NewLifecycleManager()
			registry := NewBetaSharingRegistry(tt.config, lifecycle)
			hash, err := registry.computeHashDirect(tt.signature)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if hash == "" {
				t.Error("expected non-empty hash")
			}
			// Hash should start with "join_" prefix
			if !strings.HasPrefix(hash, "join_") {
				t.Errorf("expected hash to start with 'join_', got: %s", hash)
			}
			// Note: metrics are internal to the registry and not directly accessible
			// They are tested indirectly through GetSharingStats()
			// Computing the same signature again should produce the same hash
			hash2, err := registry.computeHashDirect(tt.signature)
			if err != nil {
				t.Errorf("unexpected error on second hash: %v", err)
			}
			if hash != hash2 {
				t.Errorf("expected consistent hash, got %s then %s", hash, hash2)
			}
		})
	}
}

// TestNormalizeSignatureFallback tests the normalizeSignatureFallback function
func TestNormalizeSignatureFallback(t *testing.T) {
	tests := []struct {
		name           string
		config         BetaSharingConfig
		signature      *JoinNodeSignature
		expectSorted   bool
		description    string
		verifyVarOrder func(t *testing.T, canonical *CanonicalJoinSignature, original *JoinNodeSignature)
	}{
		{
			name: "normalize with order enabled",
			config: BetaSharingConfig{
				NormalizeOrder: true,
			},
			signature: &JoinNodeSignature{
				LeftVars:  []string{"z", "a", "m"},
				RightVars: []string{"x", "b"},
				AllVars:   []string{"z", "a", "m", "x", "b"},
				Condition: map[string]interface{}{"op": "test"},
				VarTypes: map[string]string{
					"z": "Type1",
					"a": "Type2",
					"m": "Type3",
					"x": "Type4",
					"b": "Type5",
				},
			},
			expectSorted: true,
			description:  "should sort variables when NormalizeOrder is true",
			verifyVarOrder: func(t *testing.T, canonical *CanonicalJoinSignature, original *JoinNodeSignature) {
				// Verify LeftVars are sorted
				expected := []string{"a", "m", "z"}
				if !equalStringSlices(canonical.LeftVars, expected) {
					t.Errorf("expected LeftVars %v, got %v", expected, canonical.LeftVars)
				}
				// Verify RightVars are sorted
				expectedRight := []string{"b", "x"}
				if !equalStringSlices(canonical.RightVars, expectedRight) {
					t.Errorf("expected RightVars %v, got %v", expectedRight, canonical.RightVars)
				}
				// Verify AllVars are sorted
				expectedAll := []string{"a", "b", "m", "x", "z"}
				if !equalStringSlices(canonical.AllVars, expectedAll) {
					t.Errorf("expected AllVars %v, got %v", expectedAll, canonical.AllVars)
				}
			},
		},
		{
			name: "normalize with order disabled",
			config: BetaSharingConfig{
				NormalizeOrder: false,
			},
			signature: &JoinNodeSignature{
				LeftVars:  []string{"z", "a", "m"},
				RightVars: []string{"x", "b"},
				AllVars:   []string{"z", "a", "m", "x", "b"},
				Condition: map[string]interface{}{"op": "test"},
				VarTypes: map[string]string{
					"z": "Type1",
					"a": "Type2",
				},
			},
			expectSorted: false,
			description:  "should preserve order when NormalizeOrder is false",
			verifyVarOrder: func(t *testing.T, canonical *CanonicalJoinSignature, original *JoinNodeSignature) {
				// Variables should remain in original order
				if !equalStringSlices(canonical.LeftVars, original.LeftVars) {
					t.Errorf("expected LeftVars to preserve order: %v, got %v", original.LeftVars, canonical.LeftVars)
				}
				if !equalStringSlices(canonical.RightVars, original.RightVars) {
					t.Errorf("expected RightVars to preserve order: %v, got %v", original.RightVars, canonical.RightVars)
				}
			},
		},
		{
			name: "empty signature",
			config: BetaSharingConfig{
				NormalizeOrder: true,
			},
			signature: &JoinNodeSignature{
				LeftVars:  []string{},
				RightVars: []string{},
				AllVars:   []string{},
				Condition: map[string]interface{}{},
				VarTypes:  map[string]string{},
			},
			expectSorted: true,
			description:  "should handle empty signature",
			verifyVarOrder: func(t *testing.T, canonical *CanonicalJoinSignature, original *JoinNodeSignature) {
				if len(canonical.LeftVars) != 0 {
					t.Error("expected empty LeftVars")
				}
				if len(canonical.VarTypes) != 0 {
					t.Error("expected empty VarTypes")
				}
			},
		},
		{
			name: "single variable",
			config: BetaSharingConfig{
				NormalizeOrder: true,
			},
			signature: &JoinNodeSignature{
				LeftVars:  []string{"x"},
				RightVars: []string{},
				AllVars:   []string{"x"},
				Condition: map[string]interface{}{"simple": true},
				VarTypes:  map[string]string{"x": "TypeX"},
			},
			expectSorted: true,
			description:  "should handle single variable correctly",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lifecycle := NewLifecycleManager()
			registry := NewBetaSharingRegistry(tt.config, lifecycle)
			canonical, err := registry.normalizeSignatureFallback(tt.signature)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			// Verify version
			if canonical.Version != "1.0" {
				t.Errorf("expected version 1.0, got %s", canonical.Version)
			}
			// Verify condition is preserved
			if canonical.Condition == nil {
				t.Error("expected condition to be preserved")
			}
			// Verify VarTypes are always sorted
			if len(canonical.VarTypes) > 1 {
				for i := 1; i < len(canonical.VarTypes); i++ {
					if canonical.VarTypes[i-1].VarName > canonical.VarTypes[i].VarName {
						t.Error("VarTypes should always be sorted")
					}
				}
			}
			// Run custom verification if provided
			if tt.verifyVarOrder != nil {
				tt.verifyVarOrder(t, canonical, tt.signature)
			}
		})
	}
}

// TestCanonicalJoinSignatureToJSON tests the ToJSON method
func TestCanonicalJoinSignatureToJSON(t *testing.T) {
	tests := []struct {
		name        string
		canonical   *CanonicalJoinSignature
		expectError bool
		contains    []string
		description string
	}{
		{
			name: "simple signature",
			canonical: &CanonicalJoinSignature{
				Version:   "1.0",
				LeftVars:  []string{"x"},
				RightVars: []string{"y"},
				AllVars:   []string{"x", "y"},
				Condition: map[string]interface{}{"op": "equals"},
				VarTypes: []VariableTypeMapping{
					{VarName: "x", TypeName: "Person"},
					{VarName: "y", TypeName: "Order"},
				},
			},
			expectError: false,
			contains:    []string{"1.0", "Person", "Order", "equals"},
			description: "should serialize simple signature",
		},
		{
			name: "complex nested structure",
			canonical: &CanonicalJoinSignature{
				Version:   "1.0",
				LeftVars:  []string{"a", "b", "c"},
				RightVars: []string{"d", "e"},
				AllVars:   []string{"a", "b", "c", "d", "e"},
				Condition: map[string]interface{}{
					"type": "and",
					"conditions": []interface{}{
						map[string]interface{}{"op": "eq", "field": "name"},
						map[string]interface{}{"op": "gt", "field": "age", "value": 18},
					},
				},
				VarTypes: []VariableTypeMapping{
					{VarName: "a", TypeName: "TypeA"},
					{VarName: "b", TypeName: "TypeB"},
					{VarName: "c", TypeName: "TypeC"},
					{VarName: "d", TypeName: "TypeD"},
					{VarName: "e", TypeName: "TypeE"},
				},
			},
			expectError: false,
			contains:    []string{"TypeA", "TypeE", "and", "name", "age"},
			description: "should serialize complex structure",
		},
		{
			name: "empty signature",
			canonical: &CanonicalJoinSignature{
				Version:   "1.0",
				LeftVars:  []string{},
				RightVars: []string{},
				AllVars:   []string{},
				Condition: map[string]interface{}{},
				VarTypes:  []VariableTypeMapping{},
			},
			expectError: false,
			contains:    []string{"1.0"},
			description: "should serialize empty signature",
		},
		{
			name: "nil condition",
			canonical: &CanonicalJoinSignature{
				Version:   "1.0",
				LeftVars:  []string{"x"},
				RightVars: []string{"y"},
				AllVars:   []string{"x", "y"},
				Condition: nil,
				VarTypes:  []VariableTypeMapping{},
			},
			expectError: false,
			contains:    []string{"1.0", "\"x\"", "\"y\""},
			description: "should handle nil condition",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonStr, err := tt.canonical.ToJSON()
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if jsonStr == "" {
				t.Error("expected non-empty JSON string")
			}
			// Verify JSON contains expected substrings
			for _, substr := range tt.contains {
				if !strings.Contains(jsonStr, substr) {
					t.Errorf("expected JSON to contain %q, got: %s", substr, jsonStr)
				}
			}
			// Verify it's valid JSON (can be parsed back)
			// This is implicitly tested by json.Marshal not returning error
		})
	}
}

// TestRecordCacheHitAndMiss tests the RecordCacheHit and RecordCacheMiss functions
func TestRecordCacheHitAndMiss(t *testing.T) {
	t.Run("record cache hits", func(t *testing.T) {
		metrics := &BetaBuildMetrics{}
		// Record multiple hits
		for i := 0; i < 10; i++ {
			RecordCacheHit(metrics)
		}
		if metrics.HashCacheHits != 10 {
			t.Errorf("expected 10 cache hits, got %d", metrics.HashCacheHits)
		}
	})
	t.Run("record cache misses", func(t *testing.T) {
		metrics := &BetaBuildMetrics{}
		// Record multiple misses
		for i := 0; i < 5; i++ {
			RecordCacheMiss(metrics)
		}
		if metrics.HashCacheMisses != 5 {
			t.Errorf("expected 5 cache misses, got %d", metrics.HashCacheMisses)
		}
	})
	t.Run("concurrent cache hit recording", func(t *testing.T) {
		metrics := &BetaBuildMetrics{}
		var wg sync.WaitGroup
		// Concurrently record hits
		goroutines := 100
		hitsPerGoroutine := 10
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < hitsPerGoroutine; j++ {
					RecordCacheHit(metrics)
				}
			}()
		}
		wg.Wait()
		expected := int64(goroutines * hitsPerGoroutine)
		if metrics.HashCacheHits != expected {
			t.Errorf("expected %d cache hits, got %d", expected, metrics.HashCacheHits)
		}
	})
	t.Run("concurrent cache miss recording", func(t *testing.T) {
		metrics := &BetaBuildMetrics{}
		var wg sync.WaitGroup
		goroutines := 50
		missesPerGoroutine := 20
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < missesPerGoroutine; j++ {
					RecordCacheMiss(metrics)
				}
			}()
		}
		wg.Wait()
		expected := int64(goroutines * missesPerGoroutine)
		if metrics.HashCacheMisses != expected {
			t.Errorf("expected %d cache misses, got %d", expected, metrics.HashCacheMisses)
		}
	})
	t.Run("mixed hits and misses", func(t *testing.T) {
		metrics := &BetaBuildMetrics{}
		RecordCacheHit(metrics)
		RecordCacheMiss(metrics)
		RecordCacheHit(metrics)
		RecordCacheHit(metrics)
		RecordCacheMiss(metrics)
		if metrics.HashCacheHits != 3 {
			t.Errorf("expected 3 hits, got %d", metrics.HashCacheHits)
		}
		if metrics.HashCacheMisses != 2 {
			t.Errorf("expected 2 misses, got %d", metrics.HashCacheMisses)
		}
	})
}

// TestNewNormalizationContext tests the NewNormalizationContext function
func TestNewNormalizationContext(t *testing.T) {
	tests := []struct {
		name   string
		config BetaSharingConfig
	}{
		{
			name: "all features enabled",
			config: BetaSharingConfig{
				NormalizeOrder:              true,
				EnableAdvancedNormalization: true,
			},
		},
		{
			name: "normalization disabled",
			config: BetaSharingConfig{
				NormalizeOrder:              false,
				EnableAdvancedNormalization: false,
			},
		},
		{
			name: "only order normalization",
			config: BetaSharingConfig{
				NormalizeOrder:              true,
				EnableAdvancedNormalization: false,
			},
		},
		{
			name: "only advanced normalization",
			config: BetaSharingConfig{
				NormalizeOrder:              false,
				EnableAdvancedNormalization: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewNormalizationContext(tt.config)
			if ctx == nil {
				t.Fatal("expected non-nil context")
			}
			// Verify config fields are copied correctly
			if ctx.NormalizeOrder != tt.config.NormalizeOrder {
				t.Errorf("expected NormalizeOrder %v, got %v", tt.config.NormalizeOrder, ctx.NormalizeOrder)
			}
			if ctx.EnableAdvancedNormalization != tt.config.EnableAdvancedNormalization {
				t.Errorf("expected EnableAdvancedNormalization %v, got %v", tt.config.EnableAdvancedNormalization, ctx.EnableAdvancedNormalization)
			}
			// Verify seenNodes map is initialized
			if ctx.seenNodes == nil {
				t.Error("expected seenNodes map to be initialized")
			}
			// Verify seenNodes starts empty
			if len(ctx.seenNodes) != 0 {
				t.Errorf("expected empty seenNodes, got %d entries", len(ctx.seenNodes))
			}
		})
	}
}

// TestNormalizationContextSeenNodes tests the seenNodes tracking in NormalizationContext
func TestNormalizationContextSeenNodes(t *testing.T) {
	config := BetaSharingConfig{
		NormalizeOrder:              true,
		EnableAdvancedNormalization: true,
	}
	ctx := NewNormalizationContext(config)
	// Test tracking multiple nodes
	node1 := "node1"
	node2 := "node2"
	node3 := 123
	ctx.seenNodes[node1] = true
	ctx.seenNodes[node2] = true
	ctx.seenNodes[node3] = true
	if len(ctx.seenNodes) != 3 {
		t.Errorf("expected 3 seen nodes, got %d", len(ctx.seenNodes))
	}
	if !ctx.seenNodes[node1] {
		t.Error("node1 should be marked as seen")
	}
	if !ctx.seenNodes[node2] {
		t.Error("node2 should be marked as seen")
	}
	if !ctx.seenNodes[node3] {
		t.Error("node3 should be marked as seen")
	}
}

// TestComputeHashDirectConsistency verifies hash consistency across multiple calls
func TestComputeHashDirectConsistency(t *testing.T) {
	config := BetaSharingConfig{
		Enabled:        true,
		NormalizeOrder: true,
		EnableMetrics:  true,
		HashCacheSize:  100,
	}
	lifecycle := NewLifecycleManager()
	registry := NewBetaSharingRegistry(config, lifecycle)
	signature := &JoinNodeSignature{
		LeftVars:  []string{"x", "y"},
		RightVars: []string{"z"},
		AllVars:   []string{"x", "y", "z"},
		Condition: map[string]interface{}{
			"operator": "equals",
			"field":    "test",
		},
		VarTypes: map[string]string{
			"x": "TypeX",
			"y": "TypeY",
			"z": "TypeZ",
		},
	}
	// Compute hash multiple times
	hashes := make([]string, 5)
	for i := 0; i < 5; i++ {
		hash, err := registry.computeHashDirect(signature)
		if err != nil {
			t.Fatalf("unexpected error on iteration %d: %v", i, err)
		}
		hashes[i] = hash
	}
	// All hashes should be identical
	firstHash := hashes[0]
	for i := 1; i < len(hashes); i++ {
		if hashes[i] != firstHash {
			t.Errorf("hash inconsistency: iteration 0 produced %s, iteration %d produced %s", firstHash, i, hashes[i])
		}
	}
}

// TestNormalizeSignatureFallbackDifferentOrders verifies that normalization produces
// the same result for signatures with variables in different orders
func TestNormalizeSignatureFallbackDifferentOrders(t *testing.T) {
	config := BetaSharingConfig{
		NormalizeOrder: true,
	}
	lifecycle := NewLifecycleManager()
	registry := NewBetaSharingRegistry(config, lifecycle)
	sig1 := &JoinNodeSignature{
		LeftVars:  []string{"a", "b", "c"},
		RightVars: []string{"x", "y"},
		AllVars:   []string{"a", "b", "c", "x", "y"},
		Condition: map[string]interface{}{"test": true},
		VarTypes: map[string]string{
			"a": "TypeA",
			"b": "TypeB",
			"c": "TypeC",
			"x": "TypeX",
			"y": "TypeY",
		},
	}
	sig2 := &JoinNodeSignature{
		LeftVars:  []string{"c", "a", "b"},
		RightVars: []string{"y", "x"},
		AllVars:   []string{"c", "a", "b", "y", "x"},
		Condition: map[string]interface{}{"test": true},
		VarTypes: map[string]string{
			"a": "TypeA",
			"b": "TypeB",
			"c": "TypeC",
			"x": "TypeX",
			"y": "TypeY",
		},
	}
	canonical1, err := registry.normalizeSignatureFallback(sig1)
	if err != nil {
		t.Fatalf("error normalizing sig1: %v", err)
	}
	canonical2, err := registry.normalizeSignatureFallback(sig2)
	if err != nil {
		t.Fatalf("error normalizing sig2: %v", err)
	}
	// Normalized versions should have identical variable lists
	if !equalStringSlices(canonical1.LeftVars, canonical2.LeftVars) {
		t.Errorf("LeftVars should be identical after normalization: %v vs %v", canonical1.LeftVars, canonical2.LeftVars)
	}
	if !equalStringSlices(canonical1.RightVars, canonical2.RightVars) {
		t.Errorf("RightVars should be identical after normalization: %v vs %v", canonical1.RightVars, canonical2.RightVars)
	}
	if !equalStringSlices(canonical1.AllVars, canonical2.AllVars) {
		t.Errorf("AllVars should be identical after normalization: %v vs %v", canonical1.AllVars, canonical2.AllVars)
	}
	// VarTypes should also match
	if len(canonical1.VarTypes) != len(canonical2.VarTypes) {
		t.Errorf("VarTypes length mismatch: %d vs %d", len(canonical1.VarTypes), len(canonical2.VarTypes))
	}
}
