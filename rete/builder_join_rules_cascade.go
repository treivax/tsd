// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strings"
)

const (
	// DefaultJoinSelectivity represents the default selectivity estimate for join operations
	// when no statistical information is available. A value of 0.5 means we assume
	// that half of the potential join combinations will match the join conditions.
	// This is a conservative middle-ground estimate used for join optimization.
	DefaultJoinSelectivity = 0.5
)

// ============================================================================
// Cascade Join Rule Creation (3+ variables) and Helper Functions
// ============================================================================

// createCascadeJoinRule creates a cascade of join nodes for multi-variable rules (3+ variables)
func (jrb *JoinRuleBuilder) createCascadeJoinRule(
	network *ReteNetwork,
	ruleID string,
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	terminalNode *TerminalNode,
) error {
	fmt.Printf("   📍 Règle multi-variables détectée (%d variables): %v\n", len(variableNames), variableNames)

	// Use BetaChainBuilder (always available)
	return jrb.createCascadeJoinRuleWithBuilder(network, ruleID, variableNames, variableTypes, condition, terminalNode)
}

// createCascadeJoinRuleWithBuilder creates a cascade using BetaChainBuilder with sharing support
func (jrb *JoinRuleBuilder) createCascadeJoinRuleWithBuilder(
	network *ReteNetwork,
	ruleID string,
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	terminalNode *TerminalNode,
) error {
	fmt.Printf("   🔧 Construction avec BetaChainBuilder (sharing enabled)\n")

	// STEP 1: Split conditions into alpha (unary) and beta (binary)
	splitter := NewConditionSplitter()
	alphaConditions, betaConditions, err := splitter.SplitConditions(condition)
	if err != nil {
		return fmt.Errorf("error splitting conditions: %w", err)
	}

	fmt.Printf("   🔍 Condition analysis: %d alpha, %d beta\n", len(alphaConditions), len(betaConditions))

	// STEP 2: Create AlphaNodes for alpha conditions (filtering before join)
	alphaNodesByVariable := make(map[string][]*AlphaNode)

	for i, alphaCond := range alphaConditions {
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

		// Create AlphaNode with filtering condition
		alphaNodeID := fmt.Sprintf("%s_alpha_%s_%d", ruleID, varName, i)
		alphaNode := NewAlphaNode(alphaNodeID, alphaCond.Condition, varName, jrb.utils.storage)

		// Register in network
		network.AlphaNodes[alphaNodeID] = alphaNode

		// Store in map
		alphaNodesByVariable[varName] = append(alphaNodesByVariable[varName], alphaNode)

		fmt.Printf("   ✨ Created AlphaNode filter: %s for variable %s\n", alphaNodeID, varName)

		// Connect TypeNode -> AlphaNode
		typeNode, exists := network.TypeNodes[varType]
		if !exists {
			return fmt.Errorf("TypeNode %s not found for alpha filter", varType)
		}
		typeNode.AddChild(alphaNode)
		fmt.Printf("   🔗 Connected %s -> AlphaFilter(%s)\n", varType, varName)
	}

	// STEP 3: Reconstruct beta-only condition for JoinNodes
	var joinCondition map[string]interface{}
	if len(betaConditions) > 0 {
		joinCondition = splitter.ReconstructBetaCondition(betaConditions)
	} else {
		joinCondition = condition
	}

	// STEP 4: Build join patterns for the chain with beta-only condition
	patterns := jrb.buildJoinPatterns(variableNames, variableTypes, joinCondition)

	// STEP 5: Build the chain using BetaChainBuilder
	chain, err := jrb.buildChainWithBuilder(network, ruleID, patterns)
	if err != nil {
		return err
	}

	// STEP 6: Connect the chain to the network with alpha nodes
	return jrb.connectChainToNetworkWithAlpha(network, ruleID, chain, variableNames, variableTypes, terminalNode, alphaNodesByVariable)
}

// buildJoinPatterns creates join patterns for the beta chain.
//
// This function constructs a sequence of JoinPattern objects that define how
// variables should be progressively joined in a cascade. For N variables, it
// creates N-1 patterns where each pattern adds one new variable to the accumulated set.
//
// Algorithm:
//  1. Pattern 1: Join first two variables (v[0] ⋈ v[1])
//  2. Pattern i: Join accumulated variables with next variable (v[0..i] ⋈ v[i+1])
//
// Example for variables [u, o, p]:
//
//	Pattern 1: LeftVars=[u], RightVars=[o], AllVars=[u,o]
//	Pattern 2: LeftVars=[u,o], RightVars=[p], AllVars=[u,o,p]
//
// Parameters:
//   - variableNames: ordered list of variable names (e.g., ["u", "o", "p"])
//   - variableTypes: corresponding types for each variable (e.g., ["User", "Order", "Product"])
//   - condition: join condition that applies to all patterns (will be filtered per pattern if needed)
//
// Returns:
//   - []JoinPattern: list of N-1 patterns defining the cascade join sequence
//
// Key properties:
//   - Each pattern's AllVars incrementally grows: [v0,v1], [v0,v1,v2], [v0,v1,v2,v3], ...
//   - Each pattern's LeftVars contains all previously joined variables
//   - Each pattern's RightVars contains exactly one new variable
//   - VarTypes maps ALL variables to their types (not just those in current pattern)
//   - Selectivity is set to DefaultJoinSelectivity for optimization purposes
//
// Thread-safety: This function is read-only and thread-safe
func (jrb *JoinRuleBuilder) buildJoinPatterns(
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
) []JoinPattern {
	// Create the variable -> type mapping
	varTypes := BuildVarTypesMap(variableNames, variableTypes)

	patterns := make([]JoinPattern, 0, len(variableNames)-1)

	// Pattern 1: First two variables
	patterns = append(patterns, JoinPattern{
		LeftVars:    []string{variableNames[0]},
		RightVars:   []string{variableNames[1]},
		AllVars:     []string{variableNames[0], variableNames[1]},
		VarTypes:    varTypes,
		Condition:   condition,
		Selectivity: DefaultJoinSelectivity,
	})

	// Patterns 2+: Each subsequent variable joins with accumulated results
	for i := 2; i < len(variableNames); i++ {
		accumulatedVars := make([]string, i)
		copy(accumulatedVars, variableNames[0:i])

		allVars := make([]string, i+1)
		copy(allVars, variableNames[0:i+1])

		patterns = append(patterns, JoinPattern{
			LeftVars:    accumulatedVars,
			RightVars:   []string{variableNames[i]},
			AllVars:     allVars,
			VarTypes:    varTypes,
			Condition:   condition,
			Selectivity: DefaultJoinSelectivity,
		})
	}

	return patterns
}

// buildChainWithBuilder builds the beta chain using the BetaChainBuilder
func (jrb *JoinRuleBuilder) buildChainWithBuilder(
	network *ReteNetwork,
	ruleID string,
	patterns []JoinPattern,
) (*BetaChain, error) {
	chain, err := network.BetaChainBuilder.BuildChain(patterns, ruleID)
	if err != nil {
		return nil, fmt.Errorf("failed to build beta chain: %w", err)
	}

	// Add all nodes to network's BetaNodes map
	for _, node := range chain.Nodes {
		network.BetaNodes[node.ID] = node
	}

	return chain, nil
}

// connectChainToNetworkWithAlpha connects the beta chain with alpha node integration
func (jrb *JoinRuleBuilder) connectChainToNetworkWithAlpha(
	network *ReteNetwork,
	ruleID string,
	chain *BetaChain,
	variableNames []string,
	variableTypes []string,
	terminalNode *TerminalNode,
	alphaNodesByVariable map[string][]*AlphaNode,
) error {
	// Connect type nodes to the first join node (for first two variables)
	if len(chain.Nodes) > 0 {
		firstJoin := chain.Nodes[0]
		for i := 0; i < 2 && i < len(variableNames); i++ {
			varName := variableNames[i]
			varType := variableTypes[i]
			side := NodeSideRight
			if i == 0 {
				side = NodeSideLeft
			}

			// Check if we have alpha filters for this variable
			if alphaNodes, hasAlphaFilters := alphaNodesByVariable[varName]; hasAlphaFilters {
				// Connect: AlphaFilter(s) -> PassthroughAlpha -> JoinNode
				passthroughAlpha := jrb.utils.GetOrCreatePassthroughAlphaNode(network, ruleID, varType, varName, side)

				// Chain each AlphaFilter to the passthrough
				for _, alphaNode := range alphaNodes {
					alphaNode.AddChild(passthroughAlpha)
				}

				// Connect passthrough to JoinNode
				passthroughAlpha.AddChild(firstJoin)

				fmt.Printf("   🔗 Chained %s -> %d AlphaFilter(s)(%s) -> Passthrough -> JoinNode (%s)\n",
					varType, len(alphaNodes), varName, side)
			} else {
				// No alpha filter, connect directly
				jrb.utils.ConnectTypeNodeToBetaNode(network, ruleID, varName, varType, firstJoin, side)
			}
		}
		fmt.Printf("   ✓ Connected first two variables to initial JoinNode\n")
	}

	// Connect subsequent variables to their respective join nodes
	for i := 2; i < len(variableNames) && i-1 < len(chain.Nodes); i++ {
		joinNode := chain.Nodes[i-1]
		varName := variableNames[i]
		varType := variableTypes[i]

		// Check if we have alpha filters for this variable
		if alphaNodes, hasAlphaFilters := alphaNodesByVariable[varName]; hasAlphaFilters {
			// Connect: AlphaFilter(s) -> PassthroughAlpha -> JoinNode
			passthroughAlpha := jrb.utils.GetOrCreatePassthroughAlphaNode(network, ruleID, varType, varName, NodeSideRight)

			// Chain each AlphaFilter to the passthrough
			for _, alphaNode := range alphaNodes {
				alphaNode.AddChild(passthroughAlpha)
			}

			// Connect passthrough to JoinNode
			passthroughAlpha.AddChild(joinNode)

			fmt.Printf("   🔗 Chained %s -> %d AlphaFilter(s)(%s) -> Passthrough -> JoinNode (right)\n",
				varType, len(alphaNodes), varName)
		} else {
			// No alpha filter, connect directly
			jrb.utils.ConnectTypeNodeToBetaNode(network, ruleID, varName, varType, joinNode, NodeSideRight)
		}
		fmt.Printf("   ✓ Connected variable %s to cascade level %d\n", varName, i)
	}

	// Connect the final node to the terminal
	if chain.FinalNode != nil {
		chain.FinalNode.AddChild(terminalNode)

		alphaCount := 0
		for _, nodes := range alphaNodesByVariable {
			alphaCount += len(nodes)
		}

		if alphaCount > 0 {
			fmt.Printf("   ✅ Alpha/Beta separation: %d alpha filters created\n", alphaCount)
		}

		fmt.Printf("   ✅ Chain complete: %d JoinNodes, %d shared\n",
			len(chain.Nodes),
			network.BetaChainBuilder.GetMetrics().SharedJoinNodesReused)
	}

	fmt.Printf("   ✅ Architecture en cascade complète avec partage: %s\n", strings.Join(variableNames, " ⋈ "))

	return nil
}
