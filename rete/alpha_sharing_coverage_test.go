// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// TestConditionHashCached_Coverage tests the cached hash function
func TestConditionHashCached_Coverage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		enableCache   bool
		condition     interface{}
		variableName  string
		expectError   bool
		errorContains string
	}{
		{
			name:        "cache disabled - simple comparison",
			enableCache: false,
			condition: map[string]interface{}{
				"type":     "comparison",
				"left":     "a",
				"operator": ">",
				"right":    "b",
			},
			variableName: "var1",
			expectError:  false,
		},
		{
			name:        "cache enabled - simple comparison",
			enableCache: true,
			condition: map[string]interface{}{
				"type":     "comparison",
				"left":     "a",
				"operator": ">",
				"right":    "b",
			},
			variableName: "var1",
			expectError:  false,
		},
		{
			name:        "cache enabled - binary operation",
			enableCache: true,
			condition: map[string]interface{}{
				"type":     "binaryOperation",
				"left":     "x",
				"operator": "+",
				"right":    "y",
			},
			variableName: "var2",
			expectError:  false,
		},
		{
			name:        "cache enabled - field access",
			enableCache: true,
			condition: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "age",
			},
			variableName: "p",
			expectError:  false,
		},
		{
			name:        "cache enabled - same condition twice (cache hit)",
			enableCache: true,
			condition: map[string]interface{}{
				"type":     "comparison",
				"left":     "a",
				"operator": "==",
				"right":    "b",
			},
			variableName: "var1",
			expectError:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &ChainPerformanceConfig{
				HashCacheEnabled: tt.enableCache,
				HashCacheMaxSize: 10,
			}
			metrics := NewChainBuildMetrics()
			registry := NewAlphaSharingRegistryWithConfig(config, metrics)
			// Call once
			hash1, err := registry.ConditionHashCached(tt.condition, tt.variableName)
			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, hash1)
			// Call again with same condition to test cache
			if tt.enableCache {
				hash2, err := registry.ConditionHashCached(tt.condition, tt.variableName)
				require.NoError(t, err)
				require.Equal(t, hash1, hash2, "hashes should be identical")
			}
		})
	}
}

// TestAlphaSharingRegistry_ClearHashCache tests cache clearing
func TestAlphaSharingRegistry_ClearHashCache(t *testing.T) {
	t.Parallel()
	config := &ChainPerformanceConfig{
		HashCacheEnabled: true,
		HashCacheMaxSize: 10,
	}
	metrics := NewChainBuildMetrics()
	registry := NewAlphaSharingRegistryWithConfig(config, metrics)
	// Add some entries to cache
	condition := map[string]interface{}{
		"type":     "comparison",
		"left":     "a",
		"operator": ">",
		"right":    "b",
	}
	_, err := registry.ConditionHashCached(condition, "var1")
	require.NoError(t, err)
	// Clear cache
	registry.ClearHashCache()
	// Try to get stats after clear
	stats := registry.GetHashCacheStats()
	require.NotNil(t, stats)
}

// TestAlphaSharingRegistry_GetHashCacheStats tests cache statistics
func TestAlphaSharingRegistry_GetHashCacheStats(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		enableCache bool
	}{
		{
			name:        "cache enabled",
			enableCache: true,
		},
		{
			name:        "cache disabled",
			enableCache: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &ChainPerformanceConfig{
				HashCacheEnabled: tt.enableCache,
				HashCacheMaxSize: 10,
			}
			metrics := NewChainBuildMetrics()
			registry := NewAlphaSharingRegistryWithConfig(config, metrics)
			stats := registry.GetHashCacheStats()
			// Stats are always returned, even when cache is disabled (returns simple map)
			require.NotNil(t, stats)
		})
	}
}

// TestNewAlphaSharingRegistryWithConfig_Coverage tests registry creation with various configs
func TestNewAlphaSharingRegistryWithConfig_Coverage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		config *ChainPerformanceConfig
	}{
		{
			name:   "nil config - uses defaults",
			config: nil,
		},
		{
			name: "cache enabled with size",
			config: &ChainPerformanceConfig{
				HashCacheEnabled: true,
				HashCacheMaxSize: 100,
			},
		},
		{
			name: "cache disabled",
			config: &ChainPerformanceConfig{
				HashCacheEnabled: false,
			},
		},
		{
			name: "small cache size",
			config: &ChainPerformanceConfig{
				HashCacheEnabled: true,
				HashCacheMaxSize: 1,
			},
		},
		{
			name: "sharing disabled",
			config: &ChainPerformanceConfig{
				HashCacheEnabled: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := NewChainBuildMetrics()
			registry := NewAlphaSharingRegistryWithConfig(tt.config, metrics)
			require.NotNil(t, registry)
			// Verify registry is functional
			condition := map[string]interface{}{
				"type":     "comparison",
				"left":     "a",
				"operator": ">",
				"right":    "b",
			}
			hash, err := registry.ConditionHashCached(condition, "var1")
			require.NoError(t, err)
			require.NotEmpty(t, hash)
		})
	}
}

// TestAlphaSharingRegistry_ResetCoverage tests registry reset functionality
func TestAlphaSharingRegistry_ResetCoverage(t *testing.T) {
	t.Parallel()
	config := &ChainPerformanceConfig{
		HashCacheEnabled: true,
		HashCacheMaxSize: 10,
	}
	metrics := NewChainBuildMetrics()
	registry := NewAlphaSharingRegistryWithConfig(config, metrics)
	// Add some entries to cache
	condition := map[string]interface{}{
		"type":     "comparison",
		"left":     "x",
		"operator": ">",
		"right":    "y",
	}
	_, err := registry.ConditionHashCached(condition, "p")
	require.NoError(t, err)
	// Reset the registry
	registry.Reset()
	// Verify we can still use the registry after reset
	stats := registry.GetHashCacheStats()
	require.NotNil(t, stats)
	// Verify we can still hash conditions after reset
	_, err = registry.ConditionHashCached(condition, "p")
	require.NoError(t, err)
}
