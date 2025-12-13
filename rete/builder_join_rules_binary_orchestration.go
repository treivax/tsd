// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strings"
)

// ============================================================================
// Binary Join Rule Creation Orchestration
// ============================================================================

// binaryJoinRuleContext holds the state for binary join rule creation orchestration
type binaryJoinRuleContext struct {
	builder       *JoinRuleBuilder
	network       *ReteNetwork
	ruleID        string
	variableNames []string
	variableTypes []string
	condition     map[string]interface{}
	terminalNode  *TerminalNode

	// Processed state
	leftVars           []string
	rightVars          []string
	varTypes           map[string]string
	alphaConditions    []SplitCondition
	betaConditions     []SplitCondition
	alphaNodesByVar    map[string][]*AlphaNode
	joinCondition      map[string]interface{}
	compositeCondition map[string]interface{}
	joinNode           *JoinNode
	joinNodeHash       string
	wasShared          bool
}

// newBinaryJoinRuleContext creates a new context for binary join rule creation
func newBinaryJoinRuleContext(
	jrb *JoinRuleBuilder,
	network *ReteNetwork,
	ruleID string,
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	terminalNode *TerminalNode,
) *binaryJoinRuleContext {
	return &binaryJoinRuleContext{
		builder:       jrb,
		network:       network,
		ruleID:        ruleID,
		variableNames: variableNames,
		variableTypes: variableTypes,
		condition:     condition,
		terminalNode:  terminalNode,
	}
}

// setupVariables sets up left/right variables and type mapping
func (ctx *binaryJoinRuleContext) setupVariables() {
	ctx.leftVars = []string{ctx.variableNames[0]}
	ctx.rightVars = []string{ctx.variableNames[1]}
	ctx.varTypes = BuildVarTypesMap(ctx.variableNames, ctx.variableTypes)
}

// splitConditions splits conditions into alpha (unary) and beta (binary)
func (ctx *binaryJoinRuleContext) splitConditions() error {
	splitter := NewConditionSplitter()
	alphaConditions, betaConditions, err := splitter.SplitConditions(ctx.condition)
	if err != nil {
		return fmt.Errorf("error splitting conditions: %w", err)
	}

	ctx.alphaConditions = alphaConditions
	ctx.betaConditions = betaConditions

	fmt.Printf("   🔍 Condition analysis: %d alpha, %d beta\n", len(alphaConditions), len(betaConditions))

	return nil
}

// createAlphaNodesWithDecomposition creates AlphaNodes for alpha conditions with decomposition
func (ctx *binaryJoinRuleContext) createAlphaNodesWithDecomposition() error {
	ctx.alphaNodesByVar = make(map[string][]*AlphaNode)

	for _, alphaCond := range ctx.alphaConditions {
		varName := alphaCond.Variable
		varType := ""

		// Find the variable type
		for j, vn := range ctx.variableNames {
			if vn == varName {
				varType = ctx.variableTypes[j]
				break
			}
		}

		if varType == "" {
			fmt.Printf("   ⚠️ Could not find type for alpha variable %s, skipping alpha filter\n", varName)
			continue
		}

		// Get TypeNode
		typeNode, exists := ctx.network.TypeNodes[varType]
		if !exists {
			return fmt.Errorf("TypeNode %s not found for alpha filter", varType)
		}

		// ALWAYS DECOMPOSE: Systematic decomposition for all alpha conditions
		fmt.Printf("   🔬 Decomposing alpha condition for %s\n", varName)

		decomposer := NewArithmeticExpressionDecomposer()
		decomposedSteps, err := decomposer.DecomposeToDecomposedConditions(alphaCond.Condition)
		if err != nil {
			return fmt.Errorf("error decomposing alpha condition: %w", err)
		}

		fmt.Printf("   ✨ Decomposed into %d atomic steps for %s\n", len(decomposedSteps), varName)

		// Build decomposed chain
		chainBuilder := NewAlphaChainBuilder(ctx.network, ctx.builder.utils.storage)
		alphaChain, err := chainBuilder.BuildDecomposedChain(
			decomposedSteps,
			varName,
			typeNode,
			ctx.ruleID,
		)
		if err != nil {
			return fmt.Errorf("error building decomposed chain: %w", err)
		}

		// Store the final node of the chain
		if alphaChain.FinalNode != nil {
			ctx.alphaNodesByVar[varName] = append(ctx.alphaNodesByVar[varName], alphaChain.FinalNode)
		}
	}

	return nil
}

// buildCompositeCondition builds the composite condition including alpha and beta parts
func (ctx *binaryJoinRuleContext) buildCompositeCondition() {
	splitter := NewConditionSplitter()

	// Reconstruct beta-only condition for JoinNode
	if len(ctx.betaConditions) > 0 {
		ctx.joinCondition = splitter.ReconstructBetaCondition(ctx.betaConditions)
	} else {
		// No beta conditions - use original condition for passthrough
		ctx.joinCondition = ctx.condition
	}

	// Build composite condition including alpha conditions for proper sharing
	// The JoinNode hash must include alpha conditions to prevent incorrect sharing
	// between rules with same beta but different alpha conditions
	ctx.compositeCondition = map[string]interface{}{
		"beta": ctx.joinCondition,
	}

	// Include alpha conditions in the composite to ensure proper hash differentiation
	if len(ctx.alphaConditions) > 0 {
		alphaCondMap := make(map[string]interface{})
		for _, alphaCond := range ctx.alphaConditions {
			// Use variable name as key and condition as value
			varKey := alphaCond.Variable
			alphaCondMap[varKey] = alphaCond.Condition
		}
		ctx.compositeCondition["alpha"] = alphaCondMap
	}
}

// createOrReuseJoinNode creates a new JoinNode or reuses a shared one
func (ctx *binaryJoinRuleContext) createOrReuseJoinNode() error {
	allVars := []string{ctx.variableNames[0], ctx.variableNames[1]}
	node, hash, shared, createErr := ctx.network.BetaSharingRegistry.GetOrCreateJoinNode(
		ctx.compositeCondition,
		ctx.leftVars,
		ctx.rightVars,
		allVars,
		ctx.varTypes,
		ctx.builder.utils.storage,
		0, // cascadeLevel: binary joins are always level 0 (first join)
	)
	if createErr != nil {
		return fmt.Errorf("failed to create JoinNode: %w", createErr)
	}

	ctx.joinNode = node
	ctx.joinNodeHash = hash
	ctx.wasShared = shared

	if shared {
		fmt.Printf("   ♻️  Reused shared JoinNode %s (hash: %s)\n", ctx.joinNode.ID, hash)
	} else {
		fmt.Printf("   ✨ Created new shared JoinNode %s (hash: %s)\n", ctx.joinNode.ID, hash)
	}

	// Register with lifecycle manager
	if ctx.network.LifecycleManager != nil {
		ctx.network.LifecycleManager.RegisterNode(hash, "JoinNode")
	}

	return nil
}

// connectTerminalNode connects the terminal node properly based on sharing status
func (ctx *binaryJoinRuleContext) connectTerminalNode() {
	if ctx.wasShared {
		// JoinNode is shared - use RuleRouterNode to avoid token duplication
		router := NewRuleRouterNode(ctx.ruleID, ctx.joinNode.ID, ctx.builder.utils.storage)
		router.SetTerminalNode(ctx.terminalNode)
		ctx.joinNode.AddChild(router)
		fmt.Printf("   🔀 Created RuleRouterNode %s for shared JoinNode %s\n", router.ID, ctx.joinNode.ID)
	} else {
		// JoinNode is new - connect terminal directly
		ctx.joinNode.AddChild(ctx.terminalNode)
	}
}

// storeJoinNodeInNetwork stores the JoinNode in the network's BetaNodes
func (ctx *binaryJoinRuleContext) storeJoinNodeInNetwork() {
	ctx.network.BetaNodes[ctx.joinNode.ID] = ctx.joinNode

	// Also store with legacy key format for test compatibility
	legacyKey := fmt.Sprintf("%s_join", ctx.ruleID)
	ctx.network.BetaNodes[legacyKey] = ctx.joinNode
}

// connectNetworkInputs connects the network inputs (TypeNode -> Alpha -> JoinNode)
func (ctx *binaryJoinRuleContext) connectNetworkInputs() {
	// IMPORTANT: Skip this step if JoinNode was shared - inputs are already connected
	if ctx.wasShared {
		fmt.Printf("   ⏭️  Skipping input reconnection for shared JoinNode %s (already connected)\n", ctx.joinNode.ID)
		return
	}

	for i, varName := range ctx.variableNames {
		varType := ctx.variableTypes[i]
		if varType == "" {
			fmt.Printf("   ⚠️ Type vide pour variable %s\n", varName)
			continue
		}

		side := NodeSideRight
		if i == 0 {
			side = NodeSideLeft
		}

		// Check if we have alpha filters for this variable
		if alphaNodes, hasAlphaFilters := ctx.alphaNodesByVar[varName]; hasAlphaFilters {
			// Connect: AlphaFilter(s) -> PassthroughAlpha -> JoinNode
			// Create passthrough alpha node
			passthroughAlpha := ctx.builder.utils.GetOrCreatePassthroughAlphaNode(ctx.network, ctx.ruleID, varType, varName, side)

			// Chain each AlphaFilter to the passthrough
			for _, alphaNode := range alphaNodes {
				alphaNode.AddChild(passthroughAlpha)
			}

			// Connect passthrough to JoinNode
			if side == NodeSideLeft {
				passthroughAlpha.AddChild(ctx.joinNode)
			} else {
				passthroughAlpha.AddChild(ctx.joinNode)
			}

			fmt.Printf("   🔗 Chained %s -> %d AlphaFilter(s)(%s) -> Passthrough -> JoinNode (%s)\n",
				varType, len(alphaNodes), varName, side)
		} else {
			// No alpha filter for this variable, connect directly
			ctx.builder.utils.ConnectTypeNodeToBetaNode(ctx.network, ctx.ruleID, varName, varType, ctx.joinNode, side)
		}
	}
}

// logCompletion logs the completion message
func (ctx *binaryJoinRuleContext) logCompletion() {
	if ctx.wasShared {
		fmt.Printf("   ✅ JoinNode %s réutilisé pour jointure %s\n", ctx.joinNode.ID, strings.Join(ctx.variableNames, " ⋈ "))
	} else {
		fmt.Printf("   ✅ JoinNode %s créé pour jointure %s\n", ctx.joinNode.ID, strings.Join(ctx.variableNames, " ⋈ "))
	}

	if len(ctx.alphaConditions) > 0 {
		fmt.Printf("   ✅ Alpha/Beta separation: %d alpha filters created\n", len(ctx.alphaConditions))
	}
}

// createBinaryJoinRuleOrchestrated orchestrates the creation of a binary join rule
// using the extract method pattern to separate concerns
func (jrb *JoinRuleBuilder) createBinaryJoinRuleOrchestrated(
	network *ReteNetwork,
	ruleID string,
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	terminalNode *TerminalNode,
) error {
	ctx := newBinaryJoinRuleContext(jrb, network, ruleID, variableNames, variableTypes, condition, terminalNode)

	// Step 1: Setup variables and type mapping
	ctx.setupVariables()

	// Step 2: Split conditions into alpha and beta
	if err := ctx.splitConditions(); err != nil {
		return err
	}

	// Step 3: Create AlphaNodes for alpha conditions with decomposition
	if err := ctx.createAlphaNodesWithDecomposition(); err != nil {
		return err
	}

	// Step 4: Build composite condition (beta + alpha for proper sharing)
	ctx.buildCompositeCondition()

	// Step 5: Create or reuse JoinNode using BetaSharingRegistry
	if err := ctx.createOrReuseJoinNode(); err != nil {
		return err
	}

	// Step 6: Connect terminal node (direct or via RuleRouterNode)
	ctx.connectTerminalNode()

	// Step 7: Store JoinNode in network
	ctx.storeJoinNodeInNetwork()

	// Step 8: Connect network inputs (TypeNode -> Alpha -> JoinNode)
	ctx.connectNetworkInputs()

	// Step 9: Log completion
	ctx.logCompletion()

	return nil
}
