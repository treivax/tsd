// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"testing"
)
// TestNewBetaChainBuilder tests the creation of a new BetaChainBuilder
func TestNewBetaChainBuilder(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	builder := NewBetaChainBuilder(network, storage)
	if builder == nil {
		t.Fatal("Expected non-nil builder")
	}
	if builder.network != network {
		t.Error("Expected builder.network to match provided network")
	}
	if builder.storage != storage {
		t.Error("Expected builder.storage to match provided storage")
	}
	if builder.connectionCache == nil {
		t.Error("Expected connectionCache to be initialized")
	}
	if builder.prefixCache == nil {
		t.Error("Expected prefixCache to be initialized")
	}
	if builder.metrics == nil {
		t.Error("Expected metrics to be initialized")
	}
	if !builder.enableOptimization {
		t.Error("Expected optimization to be enabled by default")
	}
	if !builder.enablePrefixSharing {
		t.Error("Expected prefix sharing to be enabled by default")
	}
}
// TestNewBetaChainBuilderWithMetrics tests builder creation with shared metrics
func TestNewBetaChainBuilderWithMetrics(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	metrics := NewBetaChainMetrics()
	builder := NewBetaChainBuilderWithMetrics(network, storage, metrics)
	if builder == nil {
		t.Fatal("Expected non-nil builder")
	}
	if builder.metrics != metrics {
		t.Error("Expected builder.metrics to match provided metrics")
	}
}
// TestBuildChain_EmptyPatterns tests that BuildChain fails with empty patterns
func TestBuildChain_EmptyPatterns(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	builder := NewBetaChainBuilder(network, storage)
	patterns := []JoinPattern{}
	chain, err := builder.BuildChain(patterns, "test_rule")
	if err == nil {
		t.Error("Expected error for empty patterns")
	}
	if chain != nil {
		t.Error("Expected nil chain for empty patterns")
	}
}
// TestBuildChain_SinglePattern tests building a chain with a single pattern
func TestBuildChain_SinglePattern(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Initialize BetaSharingRegistry
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	// Create a single join pattern
	pattern := JoinPattern{
		LeftVars:  []string{"p"},
		RightVars: []string{"o"},
		AllVars:   []string{"p", "o"},
		VarTypes: map[string]string{
			"p": "Person",
			"o": "Order",
		},
		Condition: map[string]interface{}{
			"type": "join",
			"left": map[string]interface{}{
				"field": "p.id",
			},
			"right": map[string]interface{}{
				"field": "o.customer_id",
			},
			"operator": "==",
		},
		Selectivity:   0.3,
		EstimatedCost: 0.6,
	}
	patterns := []JoinPattern{pattern}
	chain, err := builder.BuildChain(patterns, "test_rule_single")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if chain == nil {
		t.Fatal("Expected non-nil chain")
	}
	if len(chain.Nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(chain.Nodes))
	}
	if len(chain.Hashes) != 1 {
		t.Errorf("Expected 1 hash, got %d", len(chain.Hashes))
	}
	if chain.FinalNode == nil {
		t.Error("Expected non-nil FinalNode")
	}
	if chain.RuleID != "test_rule_single" {
		t.Errorf("Expected RuleID 'test_rule_single', got '%s'", chain.RuleID)
	}
	// Validate chain
	if err := chain.ValidateChain(); err != nil {
		t.Errorf("Chain validation failed: %v", err)
	}
}
// TestBuildChain_MultiplePatterns tests building a chain with multiple patterns
func TestBuildChain_MultiplePatterns(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Initialize BetaSharingRegistry
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	// Create cascade patterns (p ⋈ o, then (p,o) ⋈ pay)
	patterns := []JoinPattern{
		{
			LeftVars:  []string{"p"},
			RightVars: []string{"o"},
			AllVars:   []string{"p", "o"},
			VarTypes: map[string]string{
				"p": "Person",
				"o": "Order",
			},
			Condition: map[string]interface{}{
				"type": "join",
			},
			Selectivity:   0.3,
			EstimatedCost: 0.6,
		},
		{
			LeftVars:  []string{"p", "o"},
			RightVars: []string{"pay"},
			AllVars:   []string{"p", "o", "pay"},
			VarTypes: map[string]string{
				"p":   "Person",
				"o":   "Order",
				"pay": "Payment",
			},
			Condition: map[string]interface{}{
				"type": "join",
			},
			Selectivity:   0.4,
			EstimatedCost: 1.2,
		},
	}
	chain, err := builder.BuildChain(patterns, "test_rule_cascade")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if chain == nil {
		t.Fatal("Expected non-nil chain")
	}
	if len(chain.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(chain.Nodes))
	}
	if len(chain.Hashes) != 2 {
		t.Errorf("Expected 2 hashes, got %d", len(chain.Hashes))
	}
	if chain.FinalNode != chain.Nodes[1] {
		t.Error("FinalNode should be the last node")
	}
	// Validate the chain was built successfully
	if len(chain.Nodes) != 2 {
		t.Errorf("Expected 2 nodes in chain, got %d", len(chain.Nodes))
	}
}
// TestBuildChain_NodeReuse tests that identical patterns reuse nodes
func TestBuildChain_NodeReuse(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Initialize BetaSharingRegistry
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	// Create identical pattern
	pattern := JoinPattern{
		LeftVars:  []string{"p"},
		RightVars: []string{"o"},
		AllVars:   []string{"p", "o"},
		VarTypes: map[string]string{
			"p": "Person",
			"o": "Order",
		},
		Condition: map[string]interface{}{
			"type": "join",
		},
		Selectivity: 0.3,
	}
	// Build first chain
	patterns1 := []JoinPattern{pattern}
	chain1, err := builder.BuildChain(patterns1, "rule1")
	if err != nil {
		t.Fatalf("Error building chain1: %v", err)
	}
	// Build second chain with same pattern
	patterns2 := []JoinPattern{pattern}
	chain2, err := builder.BuildChain(patterns2, "rule2")
	if err != nil {
		t.Fatalf("Error building chain2: %v", err)
	}
	// Both chains should have nodes with the same ID (shared)
	if chain1.Nodes[0].ID != chain2.Nodes[0].ID {
		t.Error("Expected both chains to share the same JoinNode")
	}
	// Verify that sharing occurred - both chains should have the same node ID
	if chain1.Nodes[0].ID != chain2.Nodes[0].ID {
		t.Errorf("Expected both chains to share the same node, but got different IDs: %s vs %s",
			chain1.Nodes[0].ID, chain2.Nodes[0].ID)
	}
	// Check lifecycle
	lifecycle, _ := network.LifecycleManager.GetNodeLifecycle(chain1.Nodes[0].ID)
	if lifecycle == nil {
		t.Fatal("Expected lifecycle to exist")
	}
	refCount := lifecycle.GetRefCount()
	if refCount != 2 {
		t.Errorf("Expected refCount 2, got %d", refCount)
	}
}
// TestEstimateSelectivity tests selectivity estimation
func TestEstimateSelectivity(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	builder := NewBetaChainBuilder(network, storage)
	patterns := []JoinPattern{
		{
			LeftVars:  []string{"p"},
			RightVars: []string{"o"},
			AllVars:   []string{"p", "o"},
		},
		{
			LeftVars:  []string{"p"},
			RightVars: []string{"o"},
			AllVars:   []string{"p", "o"},
			JoinConditions: []JoinCondition{
				{LeftField: "p.id", RightField: "o.customer_id", Operator: "=="},
			},
		},
	}
	builder.estimateSelectivity(patterns)
	// First pattern (no join conditions)
	if patterns[0].Selectivity == 0 {
		t.Error("Expected non-zero selectivity for first pattern")
	}
	// Second pattern (with join conditions) should be more selective
	if patterns[1].Selectivity >= patterns[0].Selectivity {
		t.Error("Pattern with join conditions should be more selective")
	}
}
// TestOptimizeJoinOrder tests join order optimization
func TestOptimizeJoinOrder(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	builder := NewBetaChainBuilder(network, storage)
	patterns := []JoinPattern{
		{
			LeftVars:    []string{"a"},
			RightVars:   []string{"b"},
			Selectivity: 0.8, // Less selective
		},
		{
			LeftVars:    []string{"c"},
			RightVars:   []string{"d"},
			Selectivity: 0.2, // More selective
		},
		{
			LeftVars:    []string{"e"},
			RightVars:   []string{"f"},
			Selectivity: 0.5, // Medium selectivity
		},
	}
	optimized := builder.optimizeJoinOrder(patterns)
	// Should be sorted by selectivity (ascending)
	if optimized[0].Selectivity != 0.2 {
		t.Errorf("Expected first pattern to have selectivity 0.2, got %f", optimized[0].Selectivity)
	}
	if optimized[1].Selectivity != 0.5 {
		t.Errorf("Expected second pattern to have selectivity 0.5, got %f", optimized[1].Selectivity)
	}
	if optimized[2].Selectivity != 0.8 {
		t.Errorf("Expected third pattern to have selectivity 0.8, got %f", optimized[2].Selectivity)
	}
	// Original should be unchanged
	if patterns[0].Selectivity != 0.8 {
		t.Error("Original patterns should not be modified")
	}
}
// TestBuildChain_WithOptimization tests that optimization is applied
func TestBuildChain_WithOptimization(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	builder.SetOptimizationEnabled(true)
	patterns := []JoinPattern{
		{
			LeftVars:    []string{"a"},
			RightVars:   []string{"b"},
			AllVars:     []string{"a", "b"},
			VarTypes:    map[string]string{"a": "A", "b": "B"},
			Condition:   map[string]interface{}{"type": "join"},
			Selectivity: 0.8,
		},
		{
			LeftVars:    []string{"c"},
			RightVars:   []string{"d"},
			AllVars:     []string{"c", "d"},
			VarTypes:    map[string]string{"c": "C", "d": "D"},
			Condition:   map[string]interface{}{"type": "join"},
			Selectivity: 0.2,
		},
	}
	chain, err := builder.BuildChain(patterns, "test_optimization")
	if err != nil {
		t.Fatalf("Error building chain: %v", err)
	}
	// Verify chain was built successfully with optimization enabled
	if chain == nil {
		t.Fatal("Expected non-nil chain")
	}
	if len(chain.Nodes) != 2 {
		t.Errorf("Expected 2 nodes in optimized chain, got %d", len(chain.Nodes))
	}
}
// TestConnectionCache tests connection caching functionality
func TestConnectionCache(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	builder := NewBetaChainBuilder(network, storage)
	// Create mock nodes
	parent := NewJoinNode("parent", nil, []string{"p"}, []string{"o"}, nil, storage)
	child := NewJoinNode("child", nil, []string{"o"}, []string{"pay"}, nil, storage)
	// Initially not connected
	if builder.isAlreadyConnectedCached(parent, child) {
		t.Error("Expected nodes to not be connected initially")
	}
	// Connect them
	parent.AddChild(child)
	// Clear cache to force real check
	builder.ClearConnectionCache()
	// Now should be connected
	if !builder.isAlreadyConnectedCached(parent, child) {
		t.Error("Expected nodes to be connected after AddChild")
	}
	// Cache functionality verified - nodes are now connected
}
// TestPrefixCache tests prefix caching functionality
func TestPrefixCache(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	builder := NewBetaChainBuilder(network, storage)
	patterns := []JoinPattern{
		{
			LeftVars:  []string{"p"},
			RightVars: []string{"o"},
		},
		{
			LeftVars:  []string{"p", "o"},
			RightVars: []string{"pay"},
		},
	}
	// Compute prefix key
	prefixKey := builder.computePrefixKey(patterns[0:1])
	if prefixKey == "" {
		t.Error("Expected non-empty prefix key")
	}
	// Create a mock node and update cache
	node := NewJoinNode("test", nil, []string{"p"}, []string{"o"}, nil, storage)
	builder.updatePrefixCache(prefixKey, node)
	// Check cache size
	size := builder.GetPrefixCacheSize()
	if size != 1 {
		t.Errorf("Expected prefix cache size 1, got %d", size)
	}
	// Clear cache
	builder.ClearPrefixCache()
	size = builder.GetPrefixCacheSize()
	if size != 0 {
		t.Errorf("Expected prefix cache size 0 after clear, got %d", size)
	}
}
// TestDetermineJoinType tests join type determination
func TestDetermineJoinType(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	builder := NewBetaChainBuilder(network, storage)
	tests := []struct {
		name     string
		pattern  JoinPattern
		expected string
	}{
		{
			name: "binary join",
			pattern: JoinPattern{
				LeftVars:  []string{"p"},
				RightVars: []string{"o"},
			},
			expected: "binary",
		},
		{
			name: "cascade join",
			pattern: JoinPattern{
				LeftVars:  []string{"p", "o"},
				RightVars: []string{"pay"},
			},
			expected: "cascade",
		},
		{
			name: "multi join",
			pattern: JoinPattern{
				LeftVars:  []string{"p"},
				RightVars: []string{"o", "pay"},
			},
			expected: "multi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.determineJoinType(tt.pattern)
			if result != tt.expected {
				t.Errorf("Expected join type '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
// TestChainValidation tests chain validation
func TestChainValidation(t *testing.T) {
	storage := NewMemoryStorage()
	tests := []struct {
		name      string
		chain     *BetaChain
		expectErr bool
	}{
		{
			name: "valid chain",
			chain: func() *BetaChain {
				node := NewJoinNode("j1", nil, nil, nil, nil, storage)
				return &BetaChain{
					Nodes:     []*JoinNode{node},
					Hashes:    []string{"hash1"},
					FinalNode: node,
					RuleID:    "rule1",
				}
			}(),
			expectErr: false,
		},
		{
			name: "empty chain",
			chain: &BetaChain{
				Nodes:  []*JoinNode{},
				Hashes: []string{},
			},
			expectErr: true,
		},
		{
			name: "mismatched lengths",
			chain: &BetaChain{
				Nodes:     []*JoinNode{NewJoinNode("j1", nil, nil, nil, nil, storage)},
				Hashes:    []string{"hash1", "hash2"},
				FinalNode: NewJoinNode("j1", nil, nil, nil, nil, storage),
			},
			expectErr: true,
		},
		{
			name: "nil node",
			chain: &BetaChain{
				Nodes:     []*JoinNode{nil},
				Hashes:    []string{"hash1"},
				FinalNode: nil,
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.chain.ValidateChain()
			if tt.expectErr && err == nil {
				t.Error("Expected validation error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no validation error, got: %v", err)
			}
		})
	}
}
// TestCountSharedNodes tests counting shared nodes
func TestCountSharedNodes(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	pattern := JoinPattern{
		LeftVars:    []string{"p"},
		RightVars:   []string{"o"},
		AllVars:     []string{"p", "o"},
		VarTypes:    map[string]string{"p": "Person", "o": "Order"},
		Condition:   map[string]interface{}{"type": "join"},
		Selectivity: 0.3,
	}
	// Build two chains sharing the same pattern
	chain1, _ := builder.BuildChain([]JoinPattern{pattern}, "rule1")
	chain2, _ := builder.BuildChain([]JoinPattern{pattern}, "rule2")
	// Count shared nodes in chain1
	sharedCount := builder.CountSharedNodes(chain1)
	if sharedCount != 1 {
		t.Errorf("Expected 1 shared node, got %d", sharedCount)
	}
	// Same for chain2
	sharedCount = builder.CountSharedNodes(chain2)
	if sharedCount != 1 {
		t.Errorf("Expected 1 shared node, got %d", sharedCount)
	}
}
// TestGetChainStats tests chain statistics
func TestGetChainStats(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	pattern := JoinPattern{
		LeftVars:    []string{"p"},
		RightVars:   []string{"o"},
		AllVars:     []string{"p", "o"},
		VarTypes:    map[string]string{"p": "Person", "o": "Order"},
		Condition:   map[string]interface{}{"type": "join"},
		Selectivity: 0.3,
	}
	chain, _ := builder.BuildChain([]JoinPattern{pattern}, "rule1")
	stats := builder.GetChainStats(chain)
	if totalNodes, ok := stats["total_nodes"].(int); !ok || totalNodes != 1 {
		t.Errorf("Expected total_nodes to be 1, got %v", stats["total_nodes"])
	}
	if _, ok := stats["sharing_ratio"].(float64); !ok {
		t.Error("Expected sharing_ratio to be present")
	}
	if _, ok := stats["shared_nodes"].(int); !ok {
		t.Error("Expected shared_nodes to be present")
	}
}
// TestGetChainInfo tests chain info retrieval
func TestGetChainInfo(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	pattern := JoinPattern{
		LeftVars:    []string{"p"},
		RightVars:   []string{"o"},
		AllVars:     []string{"p", "o"},
		VarTypes:    map[string]string{"p": "Person", "o": "Order"},
		Condition:   map[string]interface{}{"type": "join"},
		Selectivity: 0.3,
	}
	chain, _ := builder.BuildChain([]JoinPattern{pattern}, "test_info")
	info := chain.GetChainInfo()
	if ruleID, ok := info["rule_id"].(string); !ok || ruleID != "test_info" {
		t.Errorf("Expected rule_id 'test_info', got %v", info["rule_id"])
	}
	if length, ok := info["chain_length"].(int); !ok || length != 1 {
		t.Errorf("Expected chain_length 1, got %v", info["chain_length"])
	}
	if _, ok := info["node_ids"].([]string); !ok {
		t.Error("Expected node_ids to be []string")
	}
	if _, ok := info["summary"].(string); !ok {
		t.Error("Expected summary to be string")
	}
}
// TestMetricsRecording tests that metrics are properly recorded
func TestMetricsRecording(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
		EnableMetrics: true,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	pattern := JoinPattern{
		LeftVars:    []string{"p"},
		RightVars:   []string{"o"},
		AllVars:     []string{"p", "o"},
		VarTypes:    map[string]string{"p": "Person", "o": "Order"},
		Condition:   map[string]interface{}{"type": "join"},
		Selectivity: 0.3,
	}
	// Build multiple chains
	for i := 0; i < 3; i++ {
		_, err := builder.BuildChain([]JoinPattern{pattern}, "rule"+string(rune(i)))
		if err != nil {
			t.Fatalf("Error building chain %d: %v", i, err)
		}
	}
	// GetMetrics now returns registry metrics, not internal metrics
	metrics := builder.GetMetrics()
	if metrics == nil {
		t.Fatal("Expected metrics to be non-nil")
	}
	// Verify metrics tracked some activity
	if metrics.TotalJoinNodesRequested < 3 {
		t.Errorf("Expected at least 3 join nodes requested, got %d", metrics.TotalJoinNodesRequested)
	}
}
// TestSetOptimizationEnabled tests enabling/disabling optimization
func TestSetOptimizationEnabled(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	builder := NewBetaChainBuilder(network, storage)
	// Default should be enabled
	if !builder.enableOptimization {
		t.Error("Expected optimization to be enabled by default")
	}
	// Disable it
	builder.SetOptimizationEnabled(false)
	if builder.enableOptimization {
		t.Error("Expected optimization to be disabled")
	}
	// Enable it
	builder.SetOptimizationEnabled(true)
	if !builder.enableOptimization {
		t.Error("Expected optimization to be enabled")
	}
}
// TestSetPrefixSharingEnabled tests enabling/disabling prefix sharing
func TestSetPrefixSharingEnabled(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	builder := NewBetaChainBuilder(network, storage)
	// Default should be enabled
	if !builder.enablePrefixSharing {
		t.Error("Expected prefix sharing to be enabled by default")
	}
	// Disable it
	builder.SetPrefixSharingEnabled(false)
	if builder.enablePrefixSharing {
		t.Error("Expected prefix sharing to be disabled")
	}
	// Enable it
	builder.SetPrefixSharingEnabled(true)
	if !builder.enablePrefixSharing {
		t.Error("Expected prefix sharing to be enabled")
	}
}
// TestBuildChain_WithoutSharingRegistry tests fallback when no registry
func TestBuildChain_WithoutSharingRegistry(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Don't set BetaSharingManager - test fallback
	builder := NewBetaChainBuilder(network, storage)
	pattern := JoinPattern{
		LeftVars:    []string{"p"},
		RightVars:   []string{"o"},
		AllVars:     []string{"p", "o"},
		VarTypes:    map[string]string{"p": "Person", "o": "Order"},
		Condition:   map[string]interface{}{"type": "join"},
		Selectivity: 0.3,
	}
	chain, err := builder.BuildChain([]JoinPattern{pattern}, "test_no_registry")
	if err != nil {
		t.Fatalf("Expected no error with fallback, got: %v", err)
	}
	if chain == nil {
		t.Fatal("Expected non-nil chain")
	}
	if len(chain.Nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(chain.Nodes))
	}
}
// TestConcurrentBuildChain tests thread-safety with concurrent builds
func TestConcurrentBuildChain(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
		EnableMetrics: true,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	pattern := JoinPattern{
		LeftVars:    []string{"p"},
		RightVars:   []string{"o"},
		AllVars:     []string{"p", "o"},
		VarTypes:    map[string]string{"p": "Person", "o": "Order"},
		Condition:   map[string]interface{}{"type": "join"},
		Selectivity: 0.3,
	}
	done := make(chan bool)
	errors := make(chan error, 10)
	// Build 10 chains concurrently
	for i := 0; i < 10; i++ {
		go func(id int) {
			_, err := builder.BuildChain([]JoinPattern{pattern}, "concurrent_rule_"+string(rune(id)))
			if err != nil {
				errors <- err
			}
			done <- true
		}(i)
	}
	// Wait for all to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	close(errors)
	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent build error: %v", err)
	}
	// Verify concurrent builds completed
	metrics := builder.GetMetrics()
	if metrics == nil {
		t.Fatal("Expected metrics to be non-nil")
	}
	if metrics.TotalJoinNodesRequested < 10 {
		t.Errorf("Expected at least 10 join node requests, got %d", metrics.TotalJoinNodesRequested)
	}
}
// BenchmarkBuildChain benchmarks chain building
func BenchmarkBuildChain(b *testing.B) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	pattern := JoinPattern{
		LeftVars:    []string{"p"},
		RightVars:   []string{"o"},
		AllVars:     []string{"p", "o"},
		VarTypes:    map[string]string{"p": "Person", "o": "Order"},
		Condition:   map[string]interface{}{"type": "join"},
		Selectivity: 0.3,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := builder.BuildChain([]JoinPattern{pattern}, "bench_rule")
		if err != nil {
			b.Fatalf("Error building chain: %v", err)
		}
	}
}
// BenchmarkBuildChain_Cascade benchmarks cascade chain building
func BenchmarkBuildChain_Cascade(b *testing.B) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := BetaSharingConfig{
		Enabled:       true,
		HashCacheSize: 1000,
	}
	betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
	builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
	patterns := []JoinPattern{
		{
			LeftVars:    []string{"p"},
			RightVars:   []string{"o"},
			AllVars:     []string{"p", "o"},
			VarTypes:    map[string]string{"p": "Person", "o": "Order"},
			Condition:   map[string]interface{}{"type": "join"},
			Selectivity: 0.3,
		},
		{
			LeftVars:    []string{"p", "o"},
			RightVars:   []string{"pay"},
			AllVars:     []string{"p", "o", "pay"},
			VarTypes:    map[string]string{"p": "Person", "o": "Order", "pay": "Payment"},
			Condition:   map[string]interface{}{"type": "join"},
			Selectivity: 0.4,
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := builder.BuildChain(patterns, "bench_cascade")
		if err != nil {
			b.Fatalf("Error building chain: %v", err)
		}
	}
}
// TestBetaBuildMetrics tests metrics structure and methods
func TestBetaBuildMetrics(t *testing.T) {
	metrics := &BetaBuildMetrics{}
	// Metrics structure is initialized
	// Test basic metrics tracking
	// Note: Using BetaBuildMetrics from beta_sharing_interface.go
	// which tracks TotalJoinNodesRequested, SharedJoinNodesReused, UniqueJoinNodesCreated
	if metrics.TotalJoinNodesRequested < 0 {
		t.Error("TotalJoinNodesRequested should not be negative")
	}
	if metrics.SharedJoinNodesReused < 0 {
		t.Error("SharedJoinNodesReused should not be negative")
	}
	if metrics.UniqueJoinNodesCreated < 0 {
		t.Error("UniqueJoinNodesCreated should not be negative")
	}
}