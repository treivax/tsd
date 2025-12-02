// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strings"
)

// JoinRuleBuilder handles the creation of join rules
type JoinRuleBuilder struct {
	utils                   *BuilderUtils
	enableDecomposition     bool // Enable arithmetic decomposition
	decompositionComplexity int  // Minimum complexity to trigger decomposition
}

// NewJoinRuleBuilder creates a new JoinRuleBuilder instance
func NewJoinRuleBuilder(utils *BuilderUtils) *JoinRuleBuilder {
	return &JoinRuleBuilder{
		utils:                   utils,
		enableDecomposition:     true, // Always enable decomposition
		decompositionComplexity: 1,    // Decompose all arithmetic expressions (even single operations)
	}
}

// SetDecompositionEnabled enables or disables arithmetic decomposition
func (jrb *JoinRuleBuilder) SetDecompositionEnabled(enabled bool) {
	jrb.enableDecomposition = enabled
}

// SetDecompositionComplexity sets the minimum complexity threshold for decomposition
func (jrb *JoinRuleBuilder) SetDecompositionComplexity(complexity int) {
	jrb.decompositionComplexity = complexity
}

// CreateJoinRule creates a join rule with JoinNode
func (jrb *JoinRuleBuilder) CreateJoinRule(
	network *ReteNetwork,
	ruleID string,
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	action *Action,
) error {
	// Create the terminal node for this rule
	terminalNode := jrb.utils.CreateTerminalNode(network, ruleID, action)

	// Delegate to the appropriate function based on the number of variables
	if len(variableNames) > 2 {
		return jrb.createCascadeJoinRule(network, ruleID, variableNames, variableTypes, condition, terminalNode)
	}

	return jrb.createBinaryJoinRule(network, ruleID, variableNames, variableTypes, condition, terminalNode)
}

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

// createCascadeJoinRuleLegacy is deprecated and removed - BetaChainBuilder is now always used
// This function is kept as a comment for historical reference
/*
func (jrb *JoinRuleBuilder) createCascadeJoinRuleLegacy(
	network *ReteNetwork,
	ruleID string,
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	terminalNode *TerminalNode,
) error {
	fmt.Printf("   🔧 Construction d'architecture en cascade de JoinNodes (legacy mode)\n")

	// STEP 1: Split conditions into alpha (unary) and beta (binary)
	splitter := NewConditionSplitter()
	alphaConditions, betaConditions, err := splitter.SplitConditions(condition)
	if err != nil {
		return fmt.Errorf("error splitting conditions: %w", err)
	}

	fmt.Printf("   🔍 Condition analysis: %d alpha, %d beta\n", len(alphaConditions), len(betaConditions))

	// STEP 2: Create AlphaNodes for alpha conditions (filtering before join)
	// Track alpha nodes by variable to connect them properly in the cascade
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

		// Store in map - multiple alpha nodes can exist for same variable
		alphaNodesByVariable[varName] = append(alphaNodesByVariable[varName], alphaNode)

		fmt.Printf("   ✨ Created AlphaNode filter: %s for variable %s\n", alphaNodeID, varName)

		// Connect TypeNode -> AlphaNode (filtering)
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
		// No beta conditions - use original condition for compatibility
		joinCondition = condition
	}

	// STEP 4: Create the first JoinNode for the first 2 variables
	leftVars := []string{variableNames[0]}
	rightVars := []string{variableNames[1]}
	currentVarTypes := map[string]string{
		variableNames[0]: variableTypes[0],
		variableNames[1]: variableTypes[1],
	}

	currentJoinNode := NewJoinNode(
		fmt.Sprintf("%s_join_%d_%d", ruleID, 0, 1),
		joinCondition,
		leftVars,
		rightVars,
		currentVarTypes,
		jrb.utils.storage,
	)
	network.BetaNodes[currentJoinNode.ID] = currentJoinNode

	// Connect the first 2 variables to the first JoinNode via AlphaNodes if they exist
	for i := 0; i < 2; i++ {
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
			passthroughAlpha.AddChild(currentJoinNode)

			fmt.Printf("   🔗 Chained %s -> %d AlphaFilter(s)(%s) -> Passthrough -> JoinNode (%s)\n",
				varType, len(alphaNodes), varName, side)
		} else {
			// No alpha filter for this variable, connect directly
			jrb.utils.ConnectTypeNodeToBetaNode(network, ruleID, varName, varType, currentJoinNode, side)
		}
		fmt.Printf("   ✓ Cascade level 1 connection\n")
	}

	fmt.Printf("   ✅ JoinNode cascade level 1: %s ⋈ %s\n", variableNames[0], variableNames[1])

	// STEP 5: Join each subsequent variable to the previous result
	for i := 2; i < len(variableNames); i++ {
		nextVarName := variableNames[i]
		nextVarType := variableTypes[i]

		// Accumulated variables so far
		accumulatedVars := variableNames[0:i]
		accumulatedVarTypes := make(map[string]string)
		for j := 0; j < i; j++ {
			accumulatedVarTypes[variableNames[j]] = variableTypes[j]
		}
		accumulatedVarTypes[nextVarName] = nextVarType

		// Create the next JoinNode with beta-only condition
		nextJoinNode := NewJoinNode(
			fmt.Sprintf("%s_join_%d", ruleID, i),
			joinCondition,
			accumulatedVars,
			[]string{nextVarName},
			accumulatedVarTypes,
			jrb.utils.storage,
		)
		network.BetaNodes[nextJoinNode.ID] = nextJoinNode

		// Connect the previous JoinNode to the new JoinNode
		currentJoinNode.AddChild(nextJoinNode)

		// Connect the new variable to the JoinNode via AlphaNodes if they exist
		if alphaNodes, hasAlphaFilters := alphaNodesByVariable[nextVarName]; hasAlphaFilters {
			// Connect: AlphaFilter(s) -> PassthroughAlpha -> JoinNode
			passthroughAlpha := jrb.utils.GetOrCreatePassthroughAlphaNode(network, ruleID, nextVarType, nextVarName, NodeSideRight)

			// Chain each AlphaFilter to the passthrough
			for _, alphaNode := range alphaNodes {
				alphaNode.AddChild(passthroughAlpha)
			}

			// Connect passthrough to JoinNode
			passthroughAlpha.AddChild(nextJoinNode)

			fmt.Printf("   🔗 Chained %s -> %d AlphaFilter(s)(%s) -> Passthrough -> JoinNode (right)\n",
				nextVarType, len(alphaNodes), nextVarName)
		} else {
			// No alpha filter, connect directly
			jrb.utils.ConnectTypeNodeToBetaNode(network, ruleID, nextVarName, nextVarType, nextJoinNode, NodeSideRight)
		}
		fmt.Printf("   ✓ Cascade level %d connection\n", i)

		fmt.Printf("   ✅ JoinNode cascade level %d: (%s) ⋈ %s\n", i, strings.Join(accumulatedVars, " ⋈ "), nextVarName)

		currentJoinNode = nextJoinNode
	}

	// Connect the last JoinNode to the terminal
	currentJoinNode.AddChild(terminalNode)

	if len(alphaConditions) > 0 {
		fmt.Printf("   ✅ Alpha/Beta separation: %d alpha filters created\n", len(alphaConditions))
	}
	fmt.Printf("   ✅ Architecture en cascade complète: %s\n", strings.Join(variableNames, " ⋈ "))

	return nil
}
*/

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

// buildJoinPatterns creates join patterns for the beta chain
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
		Selectivity: 0.5, // Default selectivity
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
			Selectivity: 0.5, // Default selectivity
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

// connectChainToNetwork connects the beta chain to the network's type nodes and terminal
func (jrb *JoinRuleBuilder) connectChainToNetwork(
	network *ReteNetwork,
	ruleID string,
	chain *BetaChain,
	variableNames []string,
	variableTypes []string,
	terminalNode *TerminalNode,
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
			jrb.utils.ConnectTypeNodeToBetaNode(network, ruleID, varName, varType, firstJoin, side)
		}
		fmt.Printf("   ✓ Connected first two variables to initial JoinNode\n")
	}

	// Connect subsequent variables to their respective join nodes
	for i := 2; i < len(variableNames) && i-1 < len(chain.Nodes); i++ {
		joinNode := chain.Nodes[i-1]
		varName := variableNames[i]
		varType := variableTypes[i]
		jrb.utils.ConnectTypeNodeToBetaNode(network, ruleID, varName, varType, joinNode, NodeSideRight)
		fmt.Printf("   ✓ Connected variable %s to cascade level %d\n", varName, i)
	}

	// Connect the final node to the terminal
	if chain.FinalNode != nil {
		chain.FinalNode.AddChild(terminalNode)
		fmt.Printf("   ✅ Chain complete: %d JoinNodes, %d shared\n",
			len(chain.Nodes),
			network.BetaChainBuilder.GetMetrics().SharedJoinNodesReused)
	}

	fmt.Printf("   ✅ Architecture en cascade complète avec partage: %s\n", strings.Join(variableNames, " ⋈ "))

	return nil
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
