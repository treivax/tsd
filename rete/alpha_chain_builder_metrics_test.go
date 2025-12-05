// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// TestAlphaChainBuilder_ClearConnectionCache tests clearing the connection cache
func TestAlphaChainBuilder_ClearConnectionCache(t *testing.T) {
	t.Run("clear empty cache", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		builder := NewAlphaChainBuilder(network, storage)

		size := builder.GetConnectionCacheSize()
		if size != 0 {
			t.Errorf("Expected empty cache, got size %d", size)
		}

		builder.ClearConnectionCache()

		size = builder.GetConnectionCacheSize()
		if size != 0 {
			t.Errorf("Expected cache to remain empty after clear, got size %d", size)
		}
	})

	t.Run("clear cache multiple times", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		builder := NewAlphaChainBuilder(network, storage)

		builder.ClearConnectionCache()
		builder.ClearConnectionCache()
		builder.ClearConnectionCache()

		size := builder.GetConnectionCacheSize()
		if size != 0 {
			t.Errorf("Expected cache size 0 after multiple clears, got %d", size)
		}
	})

	t.Run("clear cache with entries", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		builder := NewAlphaChainBuilder(network, storage)

		// Manually populate cache to test clearing
		builder.connectionCache["parent1->child1"] = true
		builder.connectionCache["parent1->child2"] = true
		builder.connectionCache["parent2->child1"] = true

		size := builder.GetConnectionCacheSize()
		if size != 3 {
			t.Errorf("Expected cache size 3, got %d", size)
		}

		builder.ClearConnectionCache()

		size = builder.GetConnectionCacheSize()
		if size != 0 {
			t.Errorf("Expected cache size 0 after clear, got %d", size)
		}
	})
}

// TestAlphaChainBuilder_GetConnectionCacheSize tests getting cache size
func TestAlphaChainBuilder_GetConnectionCacheSize(t *testing.T) {
	t.Run("size of new builder", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		builder := NewAlphaChainBuilder(network, storage)

		size := builder.GetConnectionCacheSize()
		if size != 0 {
			t.Errorf("Expected cache size 0 for new builder, got %d", size)
		}
	})

	t.Run("size reflects cache entries", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		builder := NewAlphaChainBuilder(network, storage)

		initialSize := builder.GetConnectionCacheSize()
		if initialSize != 0 {
			t.Errorf("Expected initial size 0, got %d", initialSize)
		}

		// Add entries manually to test size tracking
		builder.connectionCache["entry1"] = true
		size1 := builder.GetConnectionCacheSize()
		if size1 != 1 {
			t.Errorf("Expected size 1 after adding entry, got %d", size1)
		}

		builder.connectionCache["entry2"] = true
		builder.connectionCache["entry3"] = true
		size2 := builder.GetConnectionCacheSize()
		if size2 != 3 {
			t.Errorf("Expected size 3 after adding more entries, got %d", size2)
		}
	})

	t.Run("size after clear is zero", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		builder := NewAlphaChainBuilder(network, storage)

		// Add entries
		builder.connectionCache["test1"] = true
		builder.connectionCache["test2"] = true

		if builder.GetConnectionCacheSize() != 2 {
			t.Error("Expected size 2 before clear")
		}

		builder.ClearConnectionCache()

		size := builder.GetConnectionCacheSize()
		if size != 0 {
			t.Errorf("Expected cache size 0 after clear, got %d", size)
		}
	})
}

// TestAlphaChainBuilder_GetMetrics tests retrieving metrics
func TestAlphaChainBuilder_GetMetrics(t *testing.T) {
	t.Run("metrics for new builder", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		builder := NewAlphaChainBuilder(network, storage)

		metrics := builder.GetMetrics()
		if metrics == nil {
			t.Fatal("GetMetrics() returned nil")
		}

		if metrics.TotalChainsBuilt != 0 {
			t.Errorf("Expected TotalChainsBuilt = 0, got %d", metrics.TotalChainsBuilt)
		}

		if metrics.TotalNodesCreated != 0 {
			t.Errorf("Expected TotalNodesCreated = 0, got %d", metrics.TotalNodesCreated)
		}

		if metrics.TotalNodesReused != 0 {
			t.Errorf("Expected TotalNodesReused = 0, got %d", metrics.TotalNodesReused)
		}
	})

	t.Run("metrics persist after cache clear", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		builder := NewAlphaChainBuilder(network, storage)

		// Manually update metrics
		builder.metrics.TotalChainsBuilt = 5
		builder.metrics.TotalNodesCreated = 20

		metricsBeforeClear := builder.GetMetrics()
		chainsBeforeClear := metricsBeforeClear.TotalChainsBuilt
		nodesBeforeClear := metricsBeforeClear.TotalNodesCreated

		builder.ClearConnectionCache()

		metricsAfterClear := builder.GetMetrics()
		chainsAfterClear := metricsAfterClear.TotalChainsBuilt
		nodesAfterClear := metricsAfterClear.TotalNodesCreated

		if chainsAfterClear != chainsBeforeClear {
			t.Errorf("Expected metrics to persist after cache clear, TotalChainsBuilt changed from %d to %d", chainsBeforeClear, chainsAfterClear)
		}

		if nodesAfterClear != nodesBeforeClear {
			t.Errorf("Expected metrics to persist after cache clear, TotalNodesCreated changed from %d to %d", nodesBeforeClear, nodesAfterClear)
		}
	})

	t.Run("metrics are not nil", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		builder := NewAlphaChainBuilder(network, storage)

		// Get metrics multiple times
		metrics1 := builder.GetMetrics()
		metrics2 := builder.GetMetrics()
		metrics3 := builder.GetMetrics()

		if metrics1 == nil {
			t.Error("First GetMetrics() returned nil")
		}
		if metrics2 == nil {
			t.Error("Second GetMetrics() returned nil")
		}
		if metrics3 == nil {
			t.Error("Third GetMetrics() returned nil")
		}

		// Should return the same metrics object
		if metrics1 != metrics2 {
			t.Error("Expected GetMetrics() to return same object")
		}
		if metrics2 != metrics3 {
			t.Error("Expected GetMetrics() to return same object")
		}
	})

	t.Run("metrics with custom initial values", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		customMetrics := NewChainBuildMetrics()
		customMetrics.TotalChainsBuilt = 10
		customMetrics.TotalNodesCreated = 50
		customMetrics.TotalNodesReused = 25

		builder := NewAlphaChainBuilderWithMetrics(network, storage, customMetrics)

		metrics := builder.GetMetrics()
		if metrics == nil {
			t.Fatal("GetMetrics() returned nil")
		}

		if metrics.TotalChainsBuilt != 10 {
			t.Errorf("Expected TotalChainsBuilt = 10, got %d", metrics.TotalChainsBuilt)
		}

		if metrics.TotalNodesCreated != 50 {
			t.Errorf("Expected TotalNodesCreated = 50, got %d", metrics.TotalNodesCreated)
		}

		if metrics.TotalNodesReused != 25 {
			t.Errorf("Expected TotalNodesReused = 25, got %d", metrics.TotalNodesReused)
		}
	})

	t.Run("metrics are independent per builder", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)

		builder1 := NewAlphaChainBuilder(network, storage)
		builder2 := NewAlphaChainBuilder(network, storage)

		// Modify metrics in builder1
		builder1.metrics.TotalChainsBuilt = 5

		metrics1 := builder1.GetMetrics()
		metrics2 := builder2.GetMetrics()

		if metrics1.TotalChainsBuilt != 5 {
			t.Errorf("Expected builder1 metrics TotalChainsBuilt = 5, got %d", metrics1.TotalChainsBuilt)
		}

		if metrics2.TotalChainsBuilt != 0 {
			t.Errorf("Expected builder2 metrics TotalChainsBuilt = 0, got %d", metrics2.TotalChainsBuilt)
		}
	})
}

// TestAlphaChainBuilder_CacheAndMetricsIntegration tests cache and metrics together
func TestAlphaChainBuilder_CacheAndMetricsIntegration(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	builder := NewAlphaChainBuilder(network, storage)

	// Initial state
	initialCacheSize := builder.GetConnectionCacheSize()
	initialMetrics := builder.GetMetrics()

	if initialCacheSize != 0 {
		t.Errorf("Expected initial cache size 0, got %d", initialCacheSize)
	}
	if initialMetrics.TotalChainsBuilt != 0 {
		t.Errorf("Expected initial chains built 0, got %d", initialMetrics.TotalChainsBuilt)
	}

	// Simulate some operations
	builder.connectionCache["op1"] = true
	builder.connectionCache["op2"] = true
	builder.metrics.TotalChainsBuilt = 2
	builder.metrics.TotalNodesCreated = 10

	// Check cache and metrics
	cacheSize := builder.GetConnectionCacheSize()
	metrics := builder.GetMetrics()

	if cacheSize != 2 {
		t.Errorf("Expected cache size 2, got %d", cacheSize)
	}
	if metrics.TotalChainsBuilt != 2 {
		t.Errorf("Expected 2 chains built, got %d", metrics.TotalChainsBuilt)
	}

	// Clear cache
	builder.ClearConnectionCache()

	// Verify cache is cleared but metrics persist
	cacheSizeAfterClear := builder.GetConnectionCacheSize()
	if cacheSizeAfterClear != 0 {
		t.Errorf("Expected cache size 0 after clear, got %d", cacheSizeAfterClear)
	}

	metricsAfterClear := builder.GetMetrics()
	if metricsAfterClear.TotalChainsBuilt != 2 {
		t.Errorf("Expected metrics to persist after cache clear, got %d chains", metricsAfterClear.TotalChainsBuilt)
	}

	// Continue operations after clear
	builder.connectionCache["op3"] = true
	builder.metrics.TotalChainsBuilt = 3

	finalMetrics := builder.GetMetrics()
	if finalMetrics.TotalChainsBuilt != 3 {
		t.Errorf("Expected 3 total chains built, got %d", finalMetrics.TotalChainsBuilt)
	}

	finalCacheSize := builder.GetConnectionCacheSize()
	if finalCacheSize != 1 {
		t.Errorf("Expected cache size 1 after new operation, got %d", finalCacheSize)
	}
}
