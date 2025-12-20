// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"time"
)

// ============================================================================
// Beta Chain Building Orchestration
// ============================================================================

// betaChainBuildContext holds the state for beta chain building orchestration
type betaChainBuildContext struct {
	builder  *BetaChainBuilder
	patterns []JoinPattern
	ruleID   string

	// Timing and metrics
	startTime           time.Time
	nodesCreated        int
	nodesReused         int
	hashesGenerated     []string
	optimizationApplied bool
	prefixReused        bool

	// Chain state
	chain             *BetaChain
	optimizedPatterns []JoinPattern
	currentParent     Node
	startPatternIndex int
}

// newBetaChainBuildContext creates a new context for beta chain building
func newBetaChainBuildContext(
	bcb *BetaChainBuilder,
	patterns []JoinPattern,
	ruleID string,
) *betaChainBuildContext {
	return &betaChainBuildContext{
		builder:         bcb,
		patterns:        patterns,
		ruleID:          ruleID,
		hashesGenerated: make([]string, 0, len(patterns)),
		chain: &BetaChain{
			Nodes:  make([]*JoinNode, 0, len(patterns)),
			Hashes: make([]string, 0, len(patterns)),
			RuleID: ruleID,
		},
	}
}

// validateInputs validates the input patterns and network state
func (ctx *betaChainBuildContext) validateInputs() error {
	if len(ctx.patterns) == 0 {
		return fmt.Errorf("impossible de construire une chaîne sans patterns")
	}

	return nil
}

// initializeMetrics initializes timing and metrics tracking
func (ctx *betaChainBuildContext) initializeMetrics() {
	ctx.startTime = time.Now()
	ctx.nodesCreated = 0
	ctx.nodesReused = 0
}

// estimateAndOptimizePatterns estimates selectivity and optimizes join order
func (ctx *betaChainBuildContext) estimateAndOptimizePatterns() {
	// Estimate selectivity
	ctx.builder.estimateSelectivity(ctx.patterns)

	// Optimize join order if enabled
	ctx.optimizedPatterns = ctx.patterns
	if ctx.builder.enableOptimization && len(ctx.patterns) > 1 {
		ctx.optimizedPatterns = ctx.builder.optimizeJoinOrder(ctx.patterns)
		if !ctx.builder.patternsEqual(ctx.patterns, ctx.optimizedPatterns) {
			ctx.optimizationApplied = true
			fmt.Printf("⚡ [BetaChainBuilder] Optimisation de l'ordre appliquée (%d patterns réordonnés) pour règle %s\n",
				len(ctx.patterns), ctx.ruleID)
		}
	}
}

// tryReusePrefixChain attempts to reuse an existing chain prefix
func (ctx *betaChainBuildContext) tryReusePrefixChain() {
	ctx.currentParent = nil
	ctx.startPatternIndex = 0

	if !ctx.builder.enablePrefixSharing || len(ctx.optimizedPatterns) <= 1 {
		return
	}

	prefixNode, prefixLen := ctx.builder.findReusablePrefix(ctx.optimizedPatterns, ctx.ruleID)
	if prefixNode != nil && prefixLen > 0 {
		ctx.prefixReused = true
		ctx.currentParent = prefixNode
		ctx.startPatternIndex = prefixLen
		ctx.nodesReused += prefixLen
		fmt.Printf("♻️  [BetaChainBuilder] Préfixe de chaîne réutilisé (%d nœuds) pour règle %s\n",
			prefixLen, ctx.ruleID)
	}
}

// createOrReuseJoinNode creates a new JoinNode or reuses an existing one
func (ctx *betaChainBuildContext) createOrReuseJoinNode(
	pattern JoinPattern,
	patternIndex int,
) (*JoinNode, string, bool, error) {
	var joinNode *JoinNode
	var hash string
	var reused bool
	var err error

	joinNode, hash, reused, err = ctx.builder.betaSharingRegistry.GetOrCreateJoinNode(
		pattern.Condition,
		pattern.LeftVars,
		pattern.RightVars,
		pattern.AllVars,
		pattern.VarTypes,
		ctx.builder.storage,
		patternIndex, // cascadeLevel: distinguishes joins at different cascade depths
	)
	if err != nil {
		return nil, "", false, fmt.Errorf("erreur lors de la création/récupération du JoinNode %d: %w", patternIndex, err)
	}

	return joinNode, hash, reused, nil
}

// registerJoinNodeWithManagers registers the join node with lifecycle and sharing managers
func (ctx *betaChainBuildContext) registerJoinNodeWithManagers(joinNode *JoinNode, hash string) {
	// Register the node if not already registered (for new nodes)
	if _, exists := ctx.builder.network.LifecycleManager.GetNodeLifecycle(hash); !exists {
		ctx.builder.network.LifecycleManager.RegisterNode(hash, "join")
	}
	// Add this rule's reference to the join node
	ctx.builder.network.LifecycleManager.AddRuleToNode(hash, ctx.ruleID, ctx.ruleID)

	// Register rule with beta sharing registry for join node tracking
	if err := ctx.builder.betaSharingRegistry.RegisterRuleForJoinNode(hash, ctx.ruleID); err != nil {
		fmt.Printf("⚠️  [BetaChainBuilder] Warning: failed to register rule %s for join node %s: %v\n",
			ctx.ruleID, hash, err)
	}
}

// handleReusedNode handles the case where a join node is reused
func (ctx *betaChainBuildContext) handleReusedNode(joinNode *JoinNode, patternIndex int) {
	ctx.nodesReused++
	fmt.Printf("♻️  [BetaChainBuilder] Réutilisation du JoinNode %s pour la règle %s (pattern %d/%d)\n",
		joinNode.ID, ctx.ruleID, patternIndex+1, len(ctx.optimizedPatterns))

	// Check connection if we have a parent
	if ctx.currentParent != nil && !ctx.builder.isAlreadyConnectedCached(ctx.currentParent, joinNode) {
		ctx.currentParent.AddChild(joinNode)
		fmt.Printf("🔗 [BetaChainBuilder] Connexion du nœud réutilisé %s au parent %s\n",
			joinNode.ID, ctx.currentParent.GetID())
	} else if ctx.currentParent != nil {
		fmt.Printf("✓  [BetaChainBuilder] Nœud %s déjà connecté au parent %s\n",
			joinNode.ID, ctx.currentParent.GetID())
	}
}

// handleNewNode handles the case where a new join node is created
func (ctx *betaChainBuildContext) handleNewNode(joinNode *JoinNode, patternIndex int) {
	ctx.nodesCreated++

	// Add node to network
	ctx.builder.network.BetaNodes[joinNode.ID] = joinNode

	// Connect to parent if we have one
	if ctx.currentParent != nil {
		ctx.currentParent.AddChild(joinNode)
		ctx.builder.updateConnectionCache(ctx.currentParent.GetID(), joinNode.ID, true)
	}

	fmt.Printf("🆕 [BetaChainBuilder] Nouveau JoinNode %s créé pour la règle %s (pattern %d/%d)\n",
		joinNode.ID, ctx.ruleID, patternIndex+1, len(ctx.optimizedPatterns))
	if ctx.currentParent != nil {
		fmt.Printf("🔗 [BetaChainBuilder] Connexion du nœud %s au parent %s\n",
			joinNode.ID, ctx.currentParent.GetID())
	}
}

// registerNodeLifecycle registers the node with lifecycle manager and logs usage
func (ctx *betaChainBuildContext) registerNodeLifecycle(joinNode *JoinNode, reused bool) {
	lifecycle := ctx.builder.network.LifecycleManager.RegisterNode(joinNode.ID, "join")
	lifecycle.AddRuleReference(ctx.ruleID, "") // RuleName can be added later if needed

	if reused {
		fmt.Printf("📊 [BetaChainBuilder] Nœud %s maintenant utilisé par %d règle(s)\n",
			joinNode.ID, lifecycle.GetRefCount())
	}
}

// updatePrefixCacheIfNeeded updates the prefix cache for future reuse
func (ctx *betaChainBuildContext) updatePrefixCacheIfNeeded(joinNode *JoinNode, patternIndex int) {
	if ctx.builder.enablePrefixSharing && patternIndex < len(ctx.optimizedPatterns)-1 {
		prefixKey := ctx.builder.computePrefixKey(ctx.optimizedPatterns[0:patternIndex+1], ctx.ruleID)
		ctx.builder.updatePrefixCache(prefixKey, joinNode)
	}
}

// buildChainFromPatterns builds the chain pattern by pattern
func (ctx *betaChainBuildContext) buildChainFromPatterns() error {
	for i := ctx.startPatternIndex; i < len(ctx.optimizedPatterns); i++ {
		pattern := ctx.optimizedPatterns[i]

		// Create or reuse join node
		joinNode, hash, reused, err := ctx.createOrReuseJoinNode(pattern, i)
		if err != nil {
			return err
		}

		// Register with managers
		ctx.registerJoinNodeWithManagers(joinNode, hash)

		// Add to chain
		ctx.chain.Nodes = append(ctx.chain.Nodes, joinNode)
		ctx.chain.Hashes = append(ctx.chain.Hashes, hash)
		ctx.hashesGenerated = append(ctx.hashesGenerated, hash)

		// Handle reused vs new node
		if reused {
			ctx.handleReusedNode(joinNode, i)
		} else {
			ctx.handleNewNode(joinNode, i)
		}

		// Register lifecycle
		ctx.registerNodeLifecycle(joinNode, reused)

		// Update prefix cache
		ctx.updatePrefixCacheIfNeeded(joinNode, i)

		// Current node becomes parent for next node
		ctx.currentParent = joinNode
	}

	return nil
}

// finalizeBetaChain finalizes the chain and sets the final node
func (ctx *betaChainBuildContext) finalizeBetaChain() {
	if len(ctx.chain.Nodes) > 0 {
		ctx.chain.FinalNode = ctx.chain.Nodes[len(ctx.chain.Nodes)-1]
	}

	buildTime := time.Since(ctx.startTime)
	fmt.Printf("✅ [BetaChainBuilder] Chaîne beta complète construite pour la règle %s: %d nœud(s) (créés: %d, réutilisés: %d) en %v\n",
		ctx.ruleID, len(ctx.chain.Nodes), ctx.nodesCreated, ctx.nodesReused, buildTime)
}

// recordMetrics records the build metrics
func (ctx *betaChainBuildContext) recordMetrics() {
	if ctx.builder.metrics == nil {
		return
	}

	buildTime := time.Since(ctx.startTime)
	detail := BetaChainMetricDetail{
		RuleID:          ctx.ruleID,
		ChainLength:     len(ctx.chain.Nodes),
		NodesCreated:    ctx.nodesCreated,
		NodesReused:     ctx.nodesReused,
		BuildTime:       buildTime,
		Timestamp:       time.Now(),
		HashesGenerated: ctx.hashesGenerated,
		JoinsExecuted:   0, // Will be updated during runtime
		TotalJoinTime:   0,
	}
	ctx.builder.metrics.RecordChainBuild(detail)
}

// BuildChainOrchestrated orchestrates the building of a beta chain
// using the extract method pattern to separate concerns
func (bcb *BetaChainBuilder) BuildChainOrchestrated(
	patterns []JoinPattern,
	ruleID string,
) (*BetaChain, error) {
	ctx := newBetaChainBuildContext(bcb, patterns, ruleID)

	// Step 1: Validate inputs
	if err := ctx.validateInputs(); err != nil {
		return nil, err
	}

	// Step 2: Initialize metrics
	ctx.initializeMetrics()

	// Step 3: Estimate selectivity and optimize pattern order
	ctx.estimateAndOptimizePatterns()

	// Step 4: Try to reuse existing chain prefix
	ctx.tryReusePrefixChain()

	// Step 5: Build chain from patterns
	if err := ctx.buildChainFromPatterns(); err != nil {
		return nil, err
	}

	// Step 6: Finalize chain
	ctx.finalizeBetaChain()

	// Step 7: Record metrics
	ctx.recordMetrics()

	return ctx.chain, nil
}
