// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// RemovalStrategy defines the interface for rule removal strategies
type RemovalStrategy interface {
	// RemoveRule removes a rule using the specific strategy
	// Returns the number of nodes deleted and any error
	RemoveRule(ruleID string, nodeIDs []string) (int, error)

	// CanHandle returns true if this strategy can handle the given rule
	CanHandle(ruleID string, nodeIDs []string) bool

	// Name returns the name of this strategy for logging purposes
	Name() string
}

// NodeRemover handles the physical removal of nodes from the network
type NodeRemover interface {
	// RemoveNodeFromNetwork removes a node from the appropriate network map
	RemoveNodeFromNetwork(nodeID string) error

	// RemoveJoinNodeFromNetwork removes a join node and its dependents
	RemoveJoinNodeFromNetwork(nodeID string) error

	// RemoveNodeWithCheck removes a node only if RefCount == 0
	RemoveNodeWithCheck(nodeID, ruleID string) error
}

// NodeConnector handles parent-child relationships between nodes
type NodeConnector interface {
	// RemoveChildFromNode removes a child node from a parent node
	RemoveChildFromNode(parent Node, child Node)

	// DisconnectChild removes a child from a node's children list
	DisconnectChild(parent Node, child Node)
}

// ChainAnalyzer provides utilities for analyzing alpha node chains
type ChainAnalyzer interface {
	// IsPartOfChain returns true if the node is part of an alpha chain
	IsPartOfChain(nodeID string) bool

	// GetChainParent returns the parent node in an alpha chain
	GetChainParent(alphaNode *AlphaNode) Node

	// OrderAlphaNodesReverse orders alpha nodes from terminal to type node
	OrderAlphaNodesReverse(alphaNodeIDs []string) []string
}

// NodeClassifier classifies nodes by type for strategy selection
type NodeClassifier interface {
	// IsJoinNode returns true if the node is a join node
	IsJoinNode(nodeID string) bool

	// ClassifyNodes separates nodes by type
	ClassifyNodes(nodeIDs []string) *NodeClassification
}

// NodeClassification holds nodes separated by type
type NodeClassification struct {
	TerminalNodes []string
	JoinNodes     []string
	AlphaNodes    []string
	TypeNodes     []string
	OtherNodes    []string
}

// BaseRemovalContext provides common context for all removal strategies
type BaseRemovalContext struct {
	Network          *ReteNetwork
	LifecycleManager *LifecycleManager
	Logger           Logger
}

// StrategySelector selects the appropriate removal strategy for a rule
type StrategySelector interface {
	// SelectStrategy chooses the best strategy for removing the given rule
	SelectStrategy(ruleID string, nodeIDs []string) RemovalStrategy
}

// DefaultStrategySelector implements strategy selection based on node analysis
type DefaultStrategySelector struct {
	simpleStrategy     RemovalStrategy
	alphaChainStrategy RemovalStrategy
	joinStrategy       RemovalStrategy
	network            *ReteNetwork
}

// NewDefaultStrategySelector creates a new strategy selector
func NewDefaultStrategySelector(
	network *ReteNetwork,
	simpleStrategy RemovalStrategy,
	alphaChainStrategy RemovalStrategy,
	joinStrategy RemovalStrategy,
) *DefaultStrategySelector {
	return &DefaultStrategySelector{
		simpleStrategy:     simpleStrategy,
		alphaChainStrategy: alphaChainStrategy,
		joinStrategy:       joinStrategy,
		network:            network,
	}
}

// SelectStrategy chooses the appropriate strategy based on rule analysis
func (s *DefaultStrategySelector) SelectStrategy(ruleID string, nodeIDs []string) RemovalStrategy {
	// Check for join nodes first (highest priority)
	if s.joinStrategy.CanHandle(ruleID, nodeIDs) {
		return s.joinStrategy
	}

	// Check for alpha chains
	if s.alphaChainStrategy.CanHandle(ruleID, nodeIDs) {
		return s.alphaChainStrategy
	}

	// Default to simple strategy
	return s.simpleStrategy
}
