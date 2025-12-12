// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"
)

// TestSimpleRuleRemovalStrategy_CanHandle tests strategy selection for simple rules
func TestSimpleRuleRemovalStrategy_CanHandle(t *testing.T) {
	network := createTestNetwork(t)
	strategy := NewSimpleRuleRemovalStrategy(network)
	tests := []struct {
		name     string
		setup    func() []string
		expected bool
	}{
		{
			name: "simple rule without chains or joins",
			setup: func() []string {
				// Create a simple alpha node
				alphaNode := &AlphaNode{
					BaseNode: BaseNode{
						ID:   "alpha_simple",
						Type: "alpha",
					},
				}
				network.AlphaNodes["alpha_simple"] = alphaNode
				// Register in lifecycle
				network.LifecycleManager.RegisterNode("alpha_simple", "alpha")
				network.LifecycleManager.AddRuleToNode("alpha_simple", "test_rule", "Test Rule")
				return []string{"alpha_simple"}
			},
			expected: true,
		},
		{
			name: "rule with chain - cannot handle",
			setup: func() []string {
				// Create a chain of alpha nodes
				alpha1 := &AlphaNode{
					BaseNode: BaseNode{
						ID:   "alpha_chain1",
						Type: "alpha",
					},
				}
				alpha2 := &AlphaNode{
					BaseNode: BaseNode{
						ID:       "alpha_chain2",
						Type:     "alpha",
						Children: []Node{alpha1},
					},
				}
				network.AlphaNodes["alpha_chain1"] = alpha1
				network.AlphaNodes["alpha_chain2"] = alpha2
				network.LifecycleManager.RegisterNode("alpha_chain1", "alpha")
				network.LifecycleManager.RegisterNode("alpha_chain2", "alpha")
				network.LifecycleManager.AddRuleToNode("alpha_chain1", "test_rule", "Test Rule")
				network.LifecycleManager.AddRuleToNode("alpha_chain2", "test_rule", "Test Rule")
				return []string{"alpha_chain1", "alpha_chain2"}
			},
			expected: false,
		},
		{
			name: "rule with join node - cannot handle",
			setup: func() []string {
				// Create a join node
				joinNode := &JoinNode{
					BaseNode: BaseNode{
						ID:   "join_test",
						Type: "join",
					},
				}
				network.BetaNodes["join_test"] = joinNode
				network.LifecycleManager.RegisterNode("join_test", "join")
				network.LifecycleManager.AddRuleToNode("join_test", "test_rule", "Test Rule")
				return []string{"join_test"}
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean network for each test
			network.AlphaNodes = make(map[string]*AlphaNode)
			network.BetaNodes = make(map[string]interface{})
			network.LifecycleManager = NewLifecycleManager()
			nodeIDs := tt.setup()
			result := strategy.CanHandle("test_rule", nodeIDs)
			if result != tt.expected {
				t.Errorf("CanHandle() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestAlphaChainRemovalStrategy_CanHandle tests strategy selection for alpha chains
func TestAlphaChainRemovalStrategy_CanHandle(t *testing.T) {
	network := createTestNetwork(t)
	strategy := NewAlphaChainRemovalStrategy(network)
	tests := []struct {
		name     string
		setup    func() []string
		expected bool
	}{
		{
			name: "rule with alpha chain - can handle",
			setup: func() []string {
				// Create a chain of alpha nodes
				alpha1 := &AlphaNode{
					BaseNode: BaseNode{
						ID:   "alpha_chain1",
						Type: "alpha",
					},
				}
				alpha2 := &AlphaNode{
					BaseNode: BaseNode{
						ID:       "alpha_chain2",
						Type:     "alpha",
						Children: []Node{alpha1},
					},
				}
				network.AlphaNodes["alpha_chain1"] = alpha1
				network.AlphaNodes["alpha_chain2"] = alpha2
				network.LifecycleManager.RegisterNode("alpha_chain1", "alpha")
				network.LifecycleManager.RegisterNode("alpha_chain2", "alpha")
				network.LifecycleManager.AddRuleToNode("alpha_chain1", "test_rule", "Test Rule")
				network.LifecycleManager.AddRuleToNode("alpha_chain2", "test_rule", "Test Rule")
				return []string{"alpha_chain1", "alpha_chain2"}
			},
			expected: true,
		},
		{
			name: "rule with join node - cannot handle",
			setup: func() []string {
				// Create chain with join node
				joinNode := &JoinNode{
					BaseNode: BaseNode{
						ID:   "join_test",
						Type: "join",
					},
				}
				network.BetaNodes["join_test"] = joinNode
				network.LifecycleManager.RegisterNode("join_test", "join")
				network.LifecycleManager.AddRuleToNode("join_test", "test_rule", "Test Rule")
				return []string{"join_test"}
			},
			expected: false,
		},
		{
			name: "simple rule without chain - cannot handle",
			setup: func() []string {
				alphaNode := &AlphaNode{
					BaseNode: BaseNode{
						ID:   "alpha_simple",
						Type: "alpha",
					},
				}
				network.AlphaNodes["alpha_simple"] = alphaNode
				network.LifecycleManager.RegisterNode("alpha_simple", "alpha")
				network.LifecycleManager.AddRuleToNode("alpha_simple", "test_rule", "Test Rule")
				return []string{"alpha_simple"}
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean network for each test
			network.AlphaNodes = make(map[string]*AlphaNode)
			network.BetaNodes = make(map[string]interface{})
			network.LifecycleManager = NewLifecycleManager()
			nodeIDs := tt.setup()
			result := strategy.CanHandle("test_rule", nodeIDs)
			if result != tt.expected {
				t.Errorf("CanHandle() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestJoinRuleRemovalStrategy_CanHandle tests strategy selection for join rules
func TestJoinRuleRemovalStrategy_CanHandle(t *testing.T) {
	network := createTestNetwork(t)
	strategy := NewJoinRuleRemovalStrategy(network)
	tests := []struct {
		name     string
		setup    func() []string
		expected bool
	}{
		{
			name: "rule with join node - can handle",
			setup: func() []string {
				joinNode := &JoinNode{
					BaseNode: BaseNode{
						ID:   "join_test",
						Type: "join",
					},
				}
				network.BetaNodes["join_test"] = joinNode
				network.LifecycleManager.RegisterNode("join_test", "join")
				network.LifecycleManager.AddRuleToNode("join_test", "test_rule", "Test Rule")
				return []string{"join_test"}
			},
			expected: true,
		},
		{
			name: "simple rule without join - cannot handle",
			setup: func() []string {
				alphaNode := &AlphaNode{
					BaseNode: BaseNode{
						ID:   "alpha_simple",
						Type: "alpha",
					},
				}
				network.AlphaNodes["alpha_simple"] = alphaNode
				network.LifecycleManager.RegisterNode("alpha_simple", "alpha")
				network.LifecycleManager.AddRuleToNode("alpha_simple", "test_rule", "Test Rule")
				return []string{"alpha_simple"}
			},
			expected: false,
		},
		{
			name: "mixed nodes with join - can handle",
			setup: func() []string {
				alphaNode := &AlphaNode{
					BaseNode: BaseNode{
						ID:   "alpha_mixed",
						Type: "alpha",
					},
				}
				joinNode := &JoinNode{
					BaseNode: BaseNode{
						ID:   "join_mixed",
						Type: "join",
					},
				}
				network.AlphaNodes["alpha_mixed"] = alphaNode
				network.BetaNodes["join_mixed"] = joinNode
				network.LifecycleManager.RegisterNode("alpha_mixed", "alpha")
				network.LifecycleManager.RegisterNode("join_mixed", "join")
				network.LifecycleManager.AddRuleToNode("alpha_mixed", "test_rule", "Test Rule")
				network.LifecycleManager.AddRuleToNode("join_mixed", "test_rule", "Test Rule")
				return []string{"alpha_mixed", "join_mixed"}
			},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean network for each test
			network.AlphaNodes = make(map[string]*AlphaNode)
			network.BetaNodes = make(map[string]interface{})
			network.LifecycleManager = NewLifecycleManager()
			nodeIDs := tt.setup()
			result := strategy.CanHandle("test_rule", nodeIDs)
			if result != tt.expected {
				t.Errorf("CanHandle() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestDefaultStrategySelector_SelectStrategy tests the strategy selector
func TestDefaultStrategySelector_SelectStrategy(t *testing.T) {
	network := createTestNetwork(t)
	simpleStrategy := NewSimpleRuleRemovalStrategy(network)
	alphaChainStrategy := NewAlphaChainRemovalStrategy(network)
	joinStrategy := NewJoinRuleRemovalStrategy(network)
	selector := NewDefaultStrategySelector(network, simpleStrategy, alphaChainStrategy, joinStrategy)
	tests := []struct {
		name             string
		setup            func() []string
		expectedStrategy string
	}{
		{
			name: "selects join strategy for join nodes",
			setup: func() []string {
				joinNode := &JoinNode{
					BaseNode: BaseNode{
						ID:   "join_test",
						Type: "join",
					},
				}
				network.BetaNodes["join_test"] = joinNode
				network.LifecycleManager.RegisterNode("join_test", "join")
				network.LifecycleManager.AddRuleToNode("join_test", "test_rule", "Test Rule")
				return []string{"join_test"}
			},
			expectedStrategy: "JoinRuleRemoval",
		},
		{
			name: "selects alpha chain strategy for chains",
			setup: func() []string {
				alpha1 := &AlphaNode{
					BaseNode: BaseNode{
						ID:   "alpha_chain1",
						Type: "alpha",
					},
				}
				alpha2 := &AlphaNode{
					BaseNode: BaseNode{
						ID:       "alpha_chain2",
						Type:     "alpha",
						Children: []Node{alpha1},
					},
				}
				network.AlphaNodes["alpha_chain1"] = alpha1
				network.AlphaNodes["alpha_chain2"] = alpha2
				network.LifecycleManager.RegisterNode("alpha_chain1", "alpha")
				network.LifecycleManager.RegisterNode("alpha_chain2", "alpha")
				network.LifecycleManager.AddRuleToNode("alpha_chain1", "test_rule", "Test Rule")
				network.LifecycleManager.AddRuleToNode("alpha_chain2", "test_rule", "Test Rule")
				return []string{"alpha_chain1", "alpha_chain2"}
			},
			expectedStrategy: "AlphaChainRemoval",
		},
		{
			name: "selects simple strategy for simple rules",
			setup: func() []string {
				alphaNode := &AlphaNode{
					BaseNode: BaseNode{
						ID:   "alpha_simple",
						Type: "alpha",
					},
				}
				network.AlphaNodes["alpha_simple"] = alphaNode
				network.LifecycleManager.RegisterNode("alpha_simple", "alpha")
				network.LifecycleManager.AddRuleToNode("alpha_simple", "test_rule", "Test Rule")
				return []string{"alpha_simple"}
			},
			expectedStrategy: "SimpleRuleRemoval",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean network for each test
			network.AlphaNodes = make(map[string]*AlphaNode)
			network.BetaNodes = make(map[string]interface{})
			network.LifecycleManager = NewLifecycleManager()
			nodeIDs := tt.setup()
			strategy := selector.SelectStrategy("test_rule", nodeIDs)
			if strategy.Name() != tt.expectedStrategy {
				t.Errorf("SelectStrategy() returned %v, expected %v", strategy.Name(), tt.expectedStrategy)
			}
		})
	}
}

// TestStrategyNames verifies that each strategy has a unique name
func TestStrategyNames(t *testing.T) {
	network := createTestNetwork(t)
	strategies := []RemovalStrategy{
		NewSimpleRuleRemovalStrategy(network),
		NewAlphaChainRemovalStrategy(network),
		NewJoinRuleRemovalStrategy(network),
	}
	names := make(map[string]bool)
	for _, strategy := range strategies {
		name := strategy.Name()
		if names[name] {
			t.Errorf("Duplicate strategy name: %s", name)
		}
		names[name] = true
		if name == "" {
			t.Error("Strategy name cannot be empty")
		}
	}
	// Verify expected names
	expectedNames := map[string]bool{
		"SimpleRuleRemoval": true,
		"AlphaChainRemoval": true,
		"JoinRuleRemoval":   true,
	}
	for name := range names {
		if !expectedNames[name] {
			t.Errorf("Unexpected strategy name: %s", name)
		}
	}
}

// createTestNetwork creates a minimal RETE network for testing
func createTestNetwork(t *testing.T) *ReteNetwork {
	t.Helper()
	network := &ReteNetwork{
		RootNode:            &RootNode{BaseNode: BaseNode{ID: "root", Type: "root"}},
		TypeNodes:           make(map[string]*TypeNode),
		AlphaNodes:          make(map[string]*AlphaNode),
		BetaNodes:           make(map[string]interface{}),
		TerminalNodes:       make(map[string]*TerminalNode),
		LifecycleManager:    NewLifecycleManager(),
		AlphaSharingManager: NewAlphaSharingRegistry(),
		logger:              NewLogger(LogLevelInfo, nil),
	}
	return network
}
