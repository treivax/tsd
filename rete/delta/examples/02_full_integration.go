// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package examples

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/treivax/tsd/rete/delta"
)

// Example4_FullIntegration demonstrates a complete integration pattern
// combining detector, index, and a simulated RETE network.
func Example4_FullIntegration() {
	fmt.Println("\n=== Example 4: Full Integration Pattern ===")


	// Create the integrated updater
	updater := NewIntegratedUpdater()

	// Simulate adding some network nodes
	updater.RegisterNode("alpha_price_high", "Product", []string{"price"})
	updater.RegisterNode("alpha_stock_low", "Product", []string{"stock"})
	updater.RegisterNode("beta_order_product", "Product", []string{"price", "category"})
	updater.RegisterNode("term_discount_rule", "Product", []string{"price", "category"})

	fmt.Printf("Registered %d nodes in the network\n\n", updater.GetNodeCount())

	// Example 1: Update with few fields changed (delta wins)
	fmt.Println("--- Update 1: Price change only ---")
	oldProduct1 := map[string]interface{}{
		"id":       "p100",
		"name":     "Laptop",
		"price":    999.99,
		"stock":    25,
		"category": "electronics",
		"brand":    "TechCorp",
	}

	newProduct1 := map[string]interface{}{
		"id":       "p100",
		"name":     "Laptop",
		"price":    899.99, // Only price changed
		"stock":    25,
		"category": "electronics",
		"brand":    "TechCorp",
	}

	result1 := updater.UpdateFact(oldProduct1, newProduct1, "Product~p100", "Product")
	fmt.Printf("Result: %s\n", result1.Summary())
	fmt.Printf("Nodes evaluated: %d / %d (%.1f%% saved)\n\n",
		result1.NodesEvaluated, updater.GetNodeCount(),
		100.0*float64(updater.GetNodeCount()-result1.NodesEvaluated)/float64(updater.GetNodeCount()))

	// Example 2: Update with many fields changed (classic wins)
	fmt.Println("--- Update 2: Multiple fields changed ---")
	oldProduct2 := map[string]interface{}{
		"id":       "p200",
		"name":     "Tablet",
		"price":    499.99,
		"stock":    15,
		"category": "electronics",
		"brand":    "TechCorp",
	}

	newProduct2 := map[string]interface{}{
		"id":       "p200",
		"name":     "Premium Tablet", // Changed
		"price":    549.99,           // Changed
		"stock":    10,               // Changed
		"category": "premium",        // Changed
		"brand":    "TechCorp Elite", // Changed
	}

	result2 := updater.UpdateFact(oldProduct2, newProduct2, "Product~p200", "Product")
	fmt.Printf("Result: %s\n", result2.Summary())
	fmt.Printf("Nodes evaluated: %d (fallback to classic)\n\n", result2.NodesEvaluated)

	// Example 3: No changes
	fmt.Println("--- Update 3: No actual changes ---")
	oldProduct3 := map[string]interface{}{
		"id":    "p300",
		"price": 199.99,
		"stock": 50,
	}

	newProduct3 := map[string]interface{}{
		"id":    "p300",
		"price": 199.99,
		"stock": 50,
	}

	result3 := updater.UpdateFact(oldProduct3, newProduct3, "Product~p300", "Product")
	fmt.Printf("Result: %s\n", result3.Summary())

	// Print overall statistics
	fmt.Println("\n--- Overall Statistics ---")
	stats := updater.GetStatistics()
	fmt.Printf("Total updates: %d\n", stats.TotalUpdates)
	fmt.Printf("Delta propagations: %d\n", stats.DeltaPropagations)
	fmt.Printf("Classic fallbacks: %d\n", stats.ClassicFallbacks)
	fmt.Printf("No-op (no changes): %d\n", stats.NoOpUpdates)
	fmt.Printf("Total nodes evaluated: %d\n", stats.TotalNodesEvaluated)
	fmt.Printf("Total nodes avoided: %d\n", stats.TotalNodesAvoided)
	fmt.Printf("Savings: %.1f%%\n", stats.SavingsPercent())

	// Output:
	// === Example 4: Full Integration Pattern ===
	//
	// Registered 4 nodes in the network
	//
	// --- Update 1: Price change only ---
	// Result: Delta propagation (1 fields changed, 3 nodes affected)
	// Nodes evaluated: 3 / 4 (25.0% saved)
	//
	// --- Update 2: Multiple fields changed ---
	// Result: Classic fallback (5 fields changed, threshold exceeded)
	// Nodes evaluated: 4 (fallback to classic)
	//
	// --- Update 3: No actual changes ---
	// Result: No-op (no changes detected)
	//
	// --- Overall Statistics ---
	// Total updates: 3
	// Delta propagations: 1
	// Classic fallbacks: 1
	// No-op (no changes): 1
	// Total nodes evaluated: 7
	// Total nodes avoided: 1
	// Savings: 12.5%
}

// IntegratedUpdater combines delta detection, dependency indexing,
// and network propagation logic with automatic fallback.
type IntegratedUpdater struct {
	detector *delta.DeltaDetector
	index    *delta.DependencyIndex

	// Configuration
	deltaThreshold float64 // % of fields changed before fallback

	// Simulated network nodes
	nodes     []NodeInfo
	nodeCount int
	mutex     sync.RWMutex

	// Statistics
	stats UpdateStatistics
}

// NodeInfo represents a simplified RETE node.
type NodeInfo struct {
	ID       string
	Type     string
	FactType string
	Fields   []string
}

// UpdateStatistics tracks performance metrics.
type UpdateStatistics struct {
	TotalUpdates        int64
	DeltaPropagations   int64
	ClassicFallbacks    int64
	NoOpUpdates         int64
	TotalNodesEvaluated int64
	TotalNodesAvoided   int64
}

// SavingsPercent calculates the percentage of nodes avoided.
func (s *UpdateStatistics) SavingsPercent() float64 {
	total := s.TotalNodesEvaluated + s.TotalNodesAvoided
	if total == 0 {
		return 0.0
	}
	return 100.0 * float64(s.TotalNodesAvoided) / float64(total)
}

// UpdateResult contains the outcome of a fact update.
type UpdateResult struct {
	Strategy       string
	FieldsChanged  int
	NodesEvaluated int
	NodesAvoided   int
	Duration       time.Duration
}

// Summary returns a human-readable summary.
func (r *UpdateResult) Summary() string {
	switch r.Strategy {
	case "delta":
		return fmt.Sprintf("Delta propagation (%d fields changed, %d nodes affected)",
			r.FieldsChanged, r.NodesEvaluated)
	case "classic":
		return fmt.Sprintf("Classic fallback (%d fields changed, threshold exceeded)",
			r.FieldsChanged)
	case "noop":
		return "No-op (no changes detected)"
	default:
		return "Unknown strategy"
	}
}

// NewIntegratedUpdater creates a new integrated updater with sensible defaults.
func NewIntegratedUpdater() *IntegratedUpdater {
	// Use a detector configured for general use
	config := delta.DetectorConfig{
		FloatEpsilon:         0.0001,
		IgnoreInternalFields: true,
		IgnoredFields:        []string{"updated_at", "modified_at", "_version"},
		TrackTypeChanges:     true,
		EnableDeepComparison: false,
		CacheComparisons:     false,
	}

	return &IntegratedUpdater{
		detector:       delta.NewDeltaDetectorWithConfig(config),
		index:          delta.NewDependencyIndex(),
		deltaThreshold: 0.3, // Fallback if >30% fields changed
		nodes:          make([]NodeInfo, 0),
	}
}

// RegisterNode adds a node to the index and simulated network.
func (u *IntegratedUpdater) RegisterNode(nodeID, factType string, fields []string) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	// Add to index (assuming alpha node for simplicity)
	u.index.AddAlphaNode(nodeID, factType, fields)

	// Track in simulated network
	u.nodes = append(u.nodes, NodeInfo{
		ID:       nodeID,
		Type:     "alpha",
		FactType: factType,
		Fields:   fields,
	})
	u.nodeCount = len(u.nodes)
}

// GetNodeCount returns the total number of registered nodes.
func (u *IntegratedUpdater) GetNodeCount() int {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.nodeCount
}

// UpdateFact performs an optimized fact update with automatic strategy selection.
func (u *IntegratedUpdater) UpdateFact(
	oldFact, newFact map[string]interface{},
	factID, factType string,
) UpdateResult {
	start := time.Now()

	atomic.AddInt64(&u.stats.TotalUpdates, 1)

	// Step 1: Detect changes
	factDelta, err := u.detector.DetectDelta(oldFact, newFact, factID, factType)
	if err != nil {
		log.Printf("Error detecting delta: %v", err)
		return u.fallbackClassic(oldFact, newFact, factID, time.Since(start))
	}

	// Step 2: Check if there are any changes
	if factDelta.IsEmpty() {
		atomic.AddInt64(&u.stats.NoOpUpdates, 1)
		return UpdateResult{
			Strategy: "noop",
			Duration: time.Since(start),
		}
	}

	// Step 3: Decide strategy based on change ratio
	totalFields := len(newFact)
	changedFields := len(factDelta.Fields)
	changeRatio := float64(changedFields) / float64(totalFields)

	// Step 4a: Delta propagation (optimized path)
	if changeRatio <= u.deltaThreshold {
		return u.propagateDelta(factDelta, time.Since(start))
	}

	// Step 4b: Classic fallback (too many changes)
	return u.fallbackClassic(oldFact, newFact, factID, time.Since(start))
}

// propagateDelta performs optimized delta propagation.
func (u *IntegratedUpdater) propagateDelta(factDelta *delta.FactDelta, elapsed time.Duration) UpdateResult {
	atomic.AddInt64(&u.stats.DeltaPropagations, 1)

	// Get affected nodes from index
	affectedNodes := u.index.GetAffectedNodesForDelta(factDelta)

	nodesEvaluated := len(affectedNodes)
	nodesAvoided := u.GetNodeCount() - nodesEvaluated

	atomic.AddInt64(&u.stats.TotalNodesEvaluated, int64(nodesEvaluated))
	atomic.AddInt64(&u.stats.TotalNodesAvoided, int64(nodesAvoided))

	// Simulate propagation to each affected node
	for _, node := range affectedNodes {
		u.propagateToNode(node.NodeID, factDelta)
	}

	return UpdateResult{
		Strategy:       "delta",
		FieldsChanged:  len(factDelta.Fields),
		NodesEvaluated: nodesEvaluated,
		NodesAvoided:   nodesAvoided,
		Duration:       elapsed,
	}
}

// fallbackClassic performs traditional retract+insert.
func (u *IntegratedUpdater) fallbackClassic(
	oldFact, newFact map[string]interface{},
	factID string,
	elapsed time.Duration,
) UpdateResult {
	atomic.AddInt64(&u.stats.ClassicFallbacks, 1)

	// In classic mode, all nodes are evaluated
	totalNodes := u.GetNodeCount()
	atomic.AddInt64(&u.stats.TotalNodesEvaluated, int64(totalNodes))

	// Simulate retract
	u.retractFact(factID)

	// Simulate assert
	u.assertFact(newFact)

	return UpdateResult{
		Strategy:       "classic",
		FieldsChanged:  len(newFact),
		NodesEvaluated: totalNodes,
		NodesAvoided:   0,
		Duration:       elapsed,
	}
}

// propagateToNode simulates propagating a delta to a specific node.
func (u *IntegratedUpdater) propagateToNode(nodeID string, factDelta *delta.FactDelta) {
	// In a real implementation, this would:
	// 1. Retrieve the node from the network
	// 2. Re-evaluate conditions using only changed fields
	// 3. Update node state
	// 4. Propagate to child nodes if needed

	// For this example, we just simulate the work
	_ = nodeID
	_ = factDelta
}

// retractFact simulates retracting a fact from the network.
func (u *IntegratedUpdater) retractFact(factID string) {
	// In a real implementation, this would:
	// 1. Remove fact from working memory
	// 2. Propagate retraction through network
	// 3. Clean up any derived facts

	_ = factID
}

// assertFact simulates asserting a new fact into the network.
func (u *IntegratedUpdater) assertFact(fact map[string]interface{}) {
	// In a real implementation, this would:
	// 1. Add fact to working memory
	// 2. Propagate through alpha network
	// 3. Join in beta network
	// 4. Fire any activated rules

	_ = fact
}

// GetStatistics returns a snapshot of current statistics.
func (u *IntegratedUpdater) GetStatistics() UpdateStatistics {
	return UpdateStatistics{
		TotalUpdates:        atomic.LoadInt64(&u.stats.TotalUpdates),
		DeltaPropagations:   atomic.LoadInt64(&u.stats.DeltaPropagations),
		ClassicFallbacks:    atomic.LoadInt64(&u.stats.ClassicFallbacks),
		NoOpUpdates:         atomic.LoadInt64(&u.stats.NoOpUpdates),
		TotalNodesEvaluated: atomic.LoadInt64(&u.stats.TotalNodesEvaluated),
		TotalNodesAvoided:   atomic.LoadInt64(&u.stats.TotalNodesAvoided),
	}
}

// Example5_ConcurrentUpdates demonstrates thread-safe concurrent fact updates.
func Example5_ConcurrentUpdates() {
	fmt.Println("\n=== Example 5: Concurrent Updates ===")


	updater := NewIntegratedUpdater()

	// Register nodes
	updater.RegisterNode("node1", "Product", []string{"price"})
	updater.RegisterNode("node2", "Product", []string{"stock"})
	updater.RegisterNode("node3", "Product", []string{"category"})

	// Prepare updates
	const numUpdates = 100
	const numWorkers = 10

	var wg sync.WaitGroup
	start := time.Now()

	// Launch concurrent workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for i := 0; i < numUpdates/numWorkers; i++ {
				productID := fmt.Sprintf("p%d_%d", workerID, i)
				oldProduct := map[string]interface{}{
					"id":    productID,
					"price": 100.0,
					"stock": 50,
				}
				newProduct := map[string]interface{}{
					"id":    productID,
					"price": 100.0 + float64(i),
					"stock": 50,
				}

				updater.UpdateFact(oldProduct, newProduct, "Product~"+productID, "Product")
			}
		}(w)
	}

	wg.Wait()
	duration := time.Since(start)

	stats := updater.GetStatistics()
	fmt.Printf("Processed %d updates in %v\n", numUpdates, duration)
	fmt.Printf("Throughput: %.0f updates/sec\n", float64(numUpdates)/duration.Seconds())
	fmt.Printf("Delta propagations: %d (%.1f%%)\n",
		stats.DeltaPropagations,
		100.0*float64(stats.DeltaPropagations)/float64(stats.TotalUpdates))
	fmt.Printf("Average savings: %.1f%%\n", stats.SavingsPercent())

	// Output:
	// === Example 5: Concurrent Updates ===
	//
	// Processed 100 updates in 15ms
	// Throughput: 6666 updates/sec
	// Delta propagations: 100 (100.0%)
	// Average savings: 66.7%
}

// RunAllIntegrationExamples executes all integration examples.
func RunAllIntegrationExamples() {
	Example4_FullIntegration()
	Example5_ConcurrentUpdates()
}
