// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strings"
)

// ============================================================================
// Binary Join Rule Creation (2 variables)
// ============================================================================

// createBinaryJoinRule creates a simple binary join rule (2 variables)
func (jrb *JoinRuleBuilder) createBinaryJoinRule(
	network *ReteNetwork,
	ruleID string,
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	terminalNode *TerminalNode,
) error {
	leftVars := []string{variableNames[0]}
	rightVars := []string{variableNames[1]}

	// Create the variable -> type mapping
	varTypes := BuildVarTypesMap(variableNames, variableTypes)

	// STEP 1: Split conditions into alpha (unary) and beta (binary)
	splitter := NewConditionSplitter()
	alphaConditions, betaConditions, err := splitter.SplitConditions(condition)
	if err != nil {
		return fmt.Errorf("error splitting conditions: %w", err)
	}

	fmt.Printf("   🔍 Condition analysis: %d alpha, %d beta\n", len(alphaConditions), len(betaConditions))

	// STEP 2: Create AlphaNodes for alpha conditions with optional decomposition
	alphaNodesByVariable := make(map[string][]*AlphaNode)

	for _, alphaCond := range alphaConditions {
		varName := alphaCond.Variable
		varType := ""

		// Find the variable type
		for j, vn := range variableNames {
			if vn == varName {
				varType = variableTypes[j]
				break
			}
		}

		if varType == "" {
			fmt.Printf("   ⚠️ Could not find type for alpha variable %s, skipping alpha filter\n", varName)
			continue
		}

		// Get TypeNode
		typeNode, exists := network.TypeNodes[varType]
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
		chainBuilder := NewAlphaChainBuilder(network, jrb.utils.storage)
		alphaChain, err := chainBuilder.BuildDecomposedChain(
			decomposedSteps,
			varName,
			typeNode,
			ruleID,
		)
		if err != nil {
			return fmt.Errorf("error building decomposed chain: %w", err)
		}

		// Store the final node of the chain
		if alphaChain.FinalNode != nil {
			alphaNodesByVariable[varName] = append(alphaNodesByVariable[varName], alphaChain.FinalNode)
		}
	}

	// STEP 3: Reconstruct beta-only condition for JoinNode
	var joinCondition map[string]interface{}
	if len(betaConditions) > 0 {
		joinCondition = splitter.ReconstructBetaCondition(betaConditions)
	} else {
		// No beta conditions - use original condition for passthrough
		joinCondition = condition
	}

	// STEP 3b: Build composite condition including alpha conditions for proper sharing
	// The JoinNode hash must include alpha conditions to prevent incorrect sharing
	// between rules with same beta but different alpha conditions
	compositeCondition := map[string]interface{}{
		"beta": joinCondition,
	}

	// Include alpha conditions in the composite to ensure proper hash differentiation
	if len(alphaConditions) > 0 {
		alphaCondMap := make(map[string]interface{})
		for _, alphaCond := range alphaConditions {
			// Use variable name as key and condition as value
			varKey := alphaCond.Variable
			alphaCondMap[varKey] = alphaCond.Condition
		}
		compositeCondition["alpha"] = alphaCondMap
	}

	// STEP 4: Create JoinNode with composite condition (beta + alpha) using BetaSharingRegistry
	var joinNode *JoinNode
	var wasShared bool

	allVars := []string{variableNames[0], variableNames[1]}
	node, hash, shared, createErr := network.BetaSharingRegistry.GetOrCreateJoinNode(
		compositeCondition,
		leftVars,
		rightVars,
		allVars,
		varTypes,
		jrb.utils.storage,
	)
	if createErr != nil {
		return fmt.Errorf("failed to create JoinNode: %w", createErr)
	}

	joinNode = node
	wasShared = shared
	if shared {
		fmt.Printf("   ♻️  Reused shared JoinNode %s (hash: %s)\n", joinNode.ID, hash)
	} else {
		fmt.Printf("   ✨ Created new shared JoinNode %s (hash: %s)\n", joinNode.ID, hash)
	}

	// Register with lifecycle manager
	if network.LifecycleManager != nil {
		network.LifecycleManager.RegisterNode(hash, "JoinNode")
	}

	// STEP 4b: Connect terminal node properly based on sharing status
	if wasShared {
		// JoinNode is shared - use RuleRouterNode to avoid token duplication
		router := NewRuleRouterNode(ruleID, joinNode.ID, jrb.utils.storage)
		router.SetTerminalNode(terminalNode)
		joinNode.AddChild(router)
		fmt.Printf("   🔀 Created RuleRouterNode %s for shared JoinNode %s\n", router.ID, joinNode.ID)
	} else {
		// JoinNode is new - connect terminal directly
		joinNode.AddChild(terminalNode)
	}

	// Store the JoinNode in the network's BetaNodes
	network.BetaNodes[joinNode.ID] = joinNode

	// Also store with legacy key format for test compatibility
	legacyKey := fmt.Sprintf("%s_join", ruleID)
	network.BetaNodes[legacyKey] = joinNode

	// STEP 5: Connect the network correctly
	// For each variable, connect: TypeNode -> [AlphaFilter] -> PassthroughAlpha -> JoinNode
	// IMPORTANT: Skip this step if JoinNode was shared - inputs are already connected
	if wasShared {
		fmt.Printf("   ⏭️  Skipping input reconnection for shared JoinNode %s (already connected)\n", joinNode.ID)
	}

	if !wasShared {
		for i, varName := range variableNames {
			varType := variableTypes[i]
			if varType == "" {
				fmt.Printf("   ⚠️ Type vide pour variable %s\n", varName)
				continue
			}

			side := NodeSideRight
			if i == 0 {
				side = NodeSideLeft
			}

			// Check if we have alpha filters for this variable
			if alphaNodes, hasAlphaFilters := alphaNodesByVariable[varName]; hasAlphaFilters {
				// Connect: AlphaFilter(s) -> PassthroughAlpha -> JoinNode
				// Create passthrough alpha node
				passthroughAlpha := jrb.utils.GetOrCreatePassthroughAlphaNode(network, ruleID, varType, varName, side)

				// Chain each AlphaFilter to the passthrough
				for _, alphaNode := range alphaNodes {
					alphaNode.AddChild(passthroughAlpha)
				}

				// Connect passthrough to JoinNode
				if side == NodeSideLeft {
					passthroughAlpha.AddChild(joinNode)
				} else {
					passthroughAlpha.AddChild(joinNode)
				}

				fmt.Printf("   🔗 Chained %s -> %d AlphaFilter(s)(%s) -> Passthrough -> JoinNode (%s)\n",
					varType, len(alphaNodes), varName, side)
			} else {
				// No alpha filter for this variable, connect directly
				jrb.utils.ConnectTypeNodeToBetaNode(network, ruleID, varName, varType, joinNode, side)
			}
		}
	}

	if wasShared {
		fmt.Printf("   ✅ JoinNode %s réutilisé pour jointure %s\n", joinNode.ID, strings.Join(variableNames, " ⋈ "))
	} else {
		fmt.Printf("   ✅ JoinNode %s créé pour jointure %s\n", joinNode.ID, strings.Join(variableNames, " ⋈ "))
	}

	if len(alphaConditions) > 0 {
		fmt.Printf("   ✅ Alpha/Beta separation: %d alpha filters created\n", len(alphaConditions))
	}

	return nil
}
