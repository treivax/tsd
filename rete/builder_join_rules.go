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
	utils *BuilderUtils
}

// NewJoinRuleBuilder creates a new JoinRuleBuilder instance
func NewJoinRuleBuilder(utils *BuilderUtils) *JoinRuleBuilder {
	return &JoinRuleBuilder{
		utils: utils,
	}
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

	var joinNode *JoinNode
	var wasShared bool

	// Try to use BetaSharingRegistry if available and enabled
	if network.BetaSharingRegistry != nil && network.Config != nil && network.Config.BetaSharingEnabled {
		allVars := []string{variableNames[0], variableNames[1]}
		node, hash, shared, err := network.BetaSharingRegistry.GetOrCreateJoinNode(
			condition,
			leftVars,
			rightVars,
			allVars,
			varTypes,
			jrb.utils.storage,
		)
		if err != nil {
			// Fallback to direct creation on error
			fmt.Printf("   ⚠️ Beta sharing failed: %v, falling back to direct creation\n", err)
			joinNode = NewJoinNode(ruleID+"_join", condition, leftVars, rightVars, varTypes, jrb.utils.storage)
		} else {
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
		}
	} else {
		// Legacy mode: direct creation
		joinNode = NewJoinNode(ruleID+"_join", condition, leftVars, rightVars, varTypes, jrb.utils.storage)
	}

	joinNode.AddChild(terminalNode)

	// Store the JoinNode in the network's BetaNodes
	network.BetaNodes[joinNode.ID] = joinNode

	// Connect the TypeNodes via pass-through AlphaNodes
	for i, varName := range variableNames {
		varType := variableTypes[i]
		if varType != "" {
			side := NodeSideRight
			if i == 0 {
				side = NodeSideLeft
			}
			jrb.utils.ConnectTypeNodeToBetaNode(network, ruleID, varName, varType, joinNode, side)
		} else {
			fmt.Printf("   ⚠️ Type vide pour variable %s\n", varName)
		}
	}

	if wasShared {
		fmt.Printf("   ✅ JoinNode %s réutilisé pour jointure %s\n", joinNode.ID, strings.Join(variableNames, " ⋈ "))
	} else {
		fmt.Printf("   ✅ JoinNode %s créé pour jointure %s\n", joinNode.ID, strings.Join(variableNames, " ⋈ "))
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

	// Try to use BetaChainBuilder if available and enabled
	if network.BetaChainBuilder != nil && network.Config != nil && network.Config.BetaSharingEnabled {
		return jrb.createCascadeJoinRuleWithBuilder(network, ruleID, variableNames, variableTypes, condition, terminalNode)
	}

	// Fallback to legacy cascade implementation
	return jrb.createCascadeJoinRuleLegacy(network, ruleID, variableNames, variableTypes, condition, terminalNode)
}

// createCascadeJoinRuleLegacy creates a cascade without BetaChainBuilder
func (jrb *JoinRuleBuilder) createCascadeJoinRuleLegacy(
	network *ReteNetwork,
	ruleID string,
	variableNames []string,
	variableTypes []string,
	condition map[string]interface{},
	terminalNode *TerminalNode,
) error {
	fmt.Printf("   🔧 Construction d'architecture en cascade de JoinNodes (legacy mode)\n")

	// Step 1: Create the first JoinNode for the first 2 variables
	leftVars := []string{variableNames[0]}
	rightVars := []string{variableNames[1]}
	currentVarTypes := map[string]string{
		variableNames[0]: variableTypes[0],
		variableNames[1]: variableTypes[1],
	}

	currentJoinNode := NewJoinNode(
		fmt.Sprintf("%s_join_%d_%d", ruleID, 0, 1),
		condition,
		leftVars,
		rightVars,
		currentVarTypes,
		jrb.utils.storage,
	)
	network.BetaNodes[currentJoinNode.ID] = currentJoinNode

	// Connect the first 2 variables to the first JoinNode
	for i := 0; i < 2; i++ {
		varName := variableNames[i]
		varType := variableTypes[i]
		side := NodeSideRight
		if i == 0 {
			side = NodeSideLeft
		}
		jrb.utils.ConnectTypeNodeToBetaNode(network, ruleID, varName, varType, currentJoinNode, side)
		fmt.Printf("   ✓ Cascade level 1 connection\n")
	}

	fmt.Printf("   ✅ JoinNode cascade level 1: %s ⋈ %s\n", variableNames[0], variableNames[1])

	// Step 2+: Join each subsequent variable to the previous result
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

		// Create the next JoinNode
		nextJoinNode := NewJoinNode(
			fmt.Sprintf("%s_join_%d", ruleID, i),
			condition,
			accumulatedVars,
			[]string{nextVarName},
			accumulatedVarTypes,
			jrb.utils.storage,
		)
		network.BetaNodes[nextJoinNode.ID] = nextJoinNode

		// Connect the previous JoinNode to the new JoinNode
		currentJoinNode.AddChild(nextJoinNode)

		// Connect the new variable to the JoinNode
		jrb.utils.ConnectTypeNodeToBetaNode(network, ruleID, nextVarName, nextVarType, nextJoinNode, NodeSideRight)
		fmt.Printf("   ✓ Cascade level %d connection\n", i)

		fmt.Printf("   ✅ JoinNode cascade level %d: (%s) ⋈ %s\n", i, strings.Join(accumulatedVars, " ⋈ "), nextVarName)

		currentJoinNode = nextJoinNode
	}

	// Connect the last JoinNode to the terminal
	currentJoinNode.AddChild(terminalNode)
	fmt.Printf("   ✅ Architecture en cascade complète: %s\n", strings.Join(variableNames, " ⋈ "))

	return nil
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

	// Build join patterns for the chain
	patterns := jrb.buildJoinPatterns(variableNames, variableTypes, condition)

	// Build the chain using BetaChainBuilder
	chain, err := jrb.buildChainWithBuilder(network, ruleID, patterns)
	if err != nil {
		return err
	}

	// Connect the chain to the network
	return jrb.connectChainToNetwork(network, ruleID, chain, variableNames, variableTypes, terminalNode)
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
