// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"log"
	"os"

	"github.com/treivax/tsd/rete/delta"
)

// NewReteNetwork crée un nouveau réseau RETE avec la configuration par défaut
func NewReteNetwork(storage Storage) *ReteNetwork {
	return NewReteNetworkWithConfig(storage, DefaultChainPerformanceConfig())
}

// NewReteNetworkWithConfig crée un nouveau réseau RETE avec une configuration personnalisée
func NewReteNetworkWithConfig(storage Storage, config *ChainPerformanceConfig) *ReteNetwork {
	if config == nil {
		config = DefaultChainPerformanceConfig()
	}

	rootNode := NewRootNode(storage)
	metrics := NewChainBuildMetrics()
	lifecycleManager := NewLifecycleManager()

	// Initialize Beta sharing (always enabled)
	betaSharingConfig := BetaSharingConfig{
		Enabled:                     true,
		HashCacheSize:               config.BetaHashCacheMaxSize,
		MaxSharedNodes:              10000, // Default limit
		EnableMetrics:               true,
		NormalizeOrder:              true,
		EnableAdvancedNormalization: false,
	}
	betaSharingRegistry := NewBetaSharingRegistry(betaSharingConfig, lifecycleManager)

	// Initialize arithmetic result cache with default config
	arithmeticCacheConfig := DefaultCacheConfig()
	arithmeticCache := NewArithmeticResultCache(arithmeticCacheConfig)

	network := &ReteNetwork{
		RootNode:              rootNode,
		TypeNodes:             make(map[string]*TypeNode),
		AlphaNodes:            make(map[string]*AlphaNode),
		BetaNodes:             make(map[string]interface{}),
		TerminalNodes:         make(map[string]*TerminalNode),
		Storage:               storage,
		Types:                 make([]TypeDefinition, 0),
		BetaBuilder:           nil, // Deprecated field, kept for backward compatibility
		LifecycleManager:      lifecycleManager,
		AlphaSharingManager:   NewAlphaSharingRegistryWithConfig(config, metrics),
		PassthroughRegistry:   make(map[string]*AlphaNode),
		BetaSharingRegistry:   betaSharingRegistry,
		BetaChainBuilder:      nil, // Will be initialized below
		ChainMetrics:          metrics,
		Config:                config,
		ArithmeticResultCache: arithmeticCache,
		logger:                NewLogger(LogLevelInfo, os.Stdout), // Logger par défaut niveau Info

		// Phase 2: Initialiser les paramètres de synchronisation
		SubmissionTimeout: DefaultSubmissionTimeout,
		VerifyRetryDelay:  DefaultVerifyRetryDelay,
		MaxVerifyRetries:  DefaultMaxVerifyRetries,
	}

	// Initialize action executor
	network.ActionExecutor = NewActionExecutor(network, log.Default())

	// Initialize BetaChainBuilder (always enabled)
	betaChainBuilder := NewBetaChainBuilderWithComponents(
		network,
		storage,
		betaSharingRegistry,
		lifecycleManager,
	)
	betaChainBuilder.SetOptimizationEnabled(true)
	betaChainBuilder.SetPrefixSharingEnabled(true)
	network.BetaChainBuilder = betaChainBuilder

	// Initialize delta propagation components
	deltaIndex := delta.NewDependencyIndex()
	deltaDetector := delta.NewDeltaDetector()
	deltaPropagator, err := delta.NewDeltaPropagatorBuilder().
		WithIndex(deltaIndex).
		WithDetector(deltaDetector).
		Build()
	if err != nil {
		// Log warning but continue - delta is optional optimization
		network.logger.Warn("Failed to initialize delta propagator: %v", err)
	} else {
		// Create callbacks for delta to interact with network
		callbacks := newReteNetworkCallbacks(network)

		// Create integration helper with all required components
		integrationHelper := delta.NewIntegrationHelper(deltaPropagator, deltaIndex, callbacks)
		integrationHelper.SetNetwork(network)

		network.DeltaPropagator = deltaPropagator
		network.DependencyIndex = deltaIndex
		network.EnableDeltaPropagation = true
		network.IntegrationHelper = integrationHelper
	}

	return network
}
