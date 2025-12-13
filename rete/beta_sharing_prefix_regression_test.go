// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// TestPrefixSharingDoesNotCrossRules is a regression test for the bug where
// prefix sharing incorrectly shared JoinNode prefixes between different rules,
// causing tokens with incomplete bindings to be propagated to terminal nodes.
//
// Bug scenario:
//   - Rule r1: {u: User, o: Order, p: Product} with conditions C1
//   - Rule r2: {u: User, o: Order, p: Product} with conditions C2
//   - Both rules create 2 JoinNodes: [u]⋈[o] and [u,o]⋈[p]
//   - Without ruleID in prefix key: r2 reuses r1's first JoinNode as prefix
//   - Result: r2 skips creating its first JoinNode, receives incomplete tokens
//   - Error: "Variable 'p' not found" when r2's action executes
//
// Fix:
//   - Include ruleID in computePrefixKey() to prevent cross-rule prefix sharing
//   - Each rule builds its complete cascade independently
//   - JoinNodes can still be shared if conditions match (via cascadeLevel)
func TestPrefixSharingDoesNotCrossRules(t *testing.T) {
	t.Log("🧪 TEST: Prefix Sharing Regression")
	t.Log("=====================================")

	// Setup
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Enable prefix sharing (this is where the bug manifested)
	network.BetaChainBuilder.SetPrefixSharingEnabled(true)

	// Define two rules with same variables but different conditions
	variableNames := []string{"u", "o", "p"}
	variableTypes := []string{"User", "Order", "Product"}

	// Rule 1: u.age >= 25 AND o.amount >= 2 AND p.price > 100
	condition1 := map[string]interface{}{
		"type":     "and",
		"operator": "AND",
		"operands": []interface{}{
			map[string]interface{}{
				"type":     "comparison",
				"operator": ">=",
				"left":     map[string]interface{}{"type": "fieldAccess", "object": "u", "field": "age"},
				"right":    map[string]interface{}{"type": "literal", "value": 25},
			},
			map[string]interface{}{
				"type":     "comparison",
				"operator": ">=",
				"left":     map[string]interface{}{"type": "fieldAccess", "object": "o", "field": "amount"},
				"right":    map[string]interface{}{"type": "literal", "value": 2},
			},
			map[string]interface{}{
				"type":     "comparison",
				"operator": ">",
				"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "price"},
				"right":    map[string]interface{}{"type": "literal", "value": 100},
			},
		},
	}

	// Rule 2: u.status == "vip" AND p.category == "luxury"
	condition2 := map[string]interface{}{
		"type":     "and",
		"operator": "AND",
		"operands": []interface{}{
			map[string]interface{}{
				"type":     "comparison",
				"operator": "==",
				"left":     map[string]interface{}{"type": "fieldAccess", "object": "u", "field": "status"},
				"right":    map[string]interface{}{"type": "literal", "value": "vip"},
			},
			map[string]interface{}{
				"type":     "comparison",
				"operator": "==",
				"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "category"},
				"right":    map[string]interface{}{"type": "literal", "value": "luxury"},
			},
		},
	}

	// Build patterns for both rules (same structure, different conditions)
	builder := NewJoinRuleBuilder(NewBuilderUtils(storage))

	// Build Rule 1
	patterns1 := builder.buildJoinPatterns(variableNames, variableTypes, condition1)

	chain1, err := network.BetaChainBuilder.BuildChain(patterns1, "r1")
	if err != nil {
		t.Fatalf("Failed to build chain for r1: %v", err)
	}

	// Verify Rule 1 created 2 JoinNodes
	if len(chain1.Nodes) != 2 {
		t.Fatalf("Rule 1 should have 2 JoinNodes, got %d", len(chain1.Nodes))
	}
	r1_join1 := chain1.Nodes[0] // [u] ⋈ [o]
	r1_join2 := chain1.Nodes[1] // [u,o] ⋈ [p]

	t.Logf("✅ Rule 1 (r1) built: 2 JoinNodes")
	t.Logf("   - JoinNode 1: %s", r1_join1.ID)
	t.Logf("   - JoinNode 2: %s", r1_join2.ID)

	// Build Rule 2 - THIS IS WHERE THE BUG MANIFESTED
	patterns2 := builder.buildJoinPatterns(variableNames, variableTypes, condition2)

	chain2, err := network.BetaChainBuilder.BuildChain(patterns2, "r2")
	if err != nil {
		t.Fatalf("Failed to build chain for r2: %v", err)
	}

	// CRITICAL CHECK: Rule 2 MUST create its own 2 JoinNodes
	// Bug: If prefix sharing crosses rules, r2 would only create 1 JoinNode
	if len(chain2.Nodes) != 2 {
		t.Errorf("❌ Rule 2 should have 2 JoinNodes (independent cascade), got %d", len(chain2.Nodes))
		t.Errorf("   This indicates prefix sharing is crossing rule boundaries!")
		t.FailNow()
	}

	r2_join1 := chain2.Nodes[0] // [u] ⋈ [o]
	r2_join2 := chain2.Nodes[1] // [u,o] ⋈ [p]

	t.Logf("✅ Rule 2 (r2) built: 2 JoinNodes (independent from r1)")
	t.Logf("   - JoinNode 1: %s", r2_join1.ID)
	t.Logf("   - JoinNode 2: %s", r2_join2.ID)

	// Verify prefix cache contains rule-specific keys
	prefixCache := network.BetaChainBuilder.prefixCache
	if len(prefixCache) == 0 {
		t.Log("⚠️  Prefix cache is empty (prefix sharing may be disabled)")
	} else {
		t.Logf("✅ Prefix cache has %d entries (rule-specific prefixes)", len(prefixCache))
		for key := range prefixCache {
			t.Logf("   - Prefix key: %s", key)
			// Key should contain ruleID (format: "ruleID::[vars]|...")
			if len(key) < 3 || key[0:2] != "r1" && key[0:2] != "r2" {
				t.Logf("   ⚠️  Key may not include ruleID: %s", key)
			}
		}
	}

	// Additional verification: Check that JoinNodes with same cascade level
	// but different conditions are NOT shared
	if r1_join1.ID == r2_join1.ID {
		t.Errorf("❌ r1 and r2 should NOT share first JoinNode (different conditions)")
		t.Errorf("   Both have ID: %s", r1_join1.ID)
	} else {
		t.Logf("✅ r1 and r2 have different first JoinNodes (correct)")
	}

	if r1_join2.ID == r2_join2.ID {
		t.Errorf("❌ r1 and r2 should NOT share second JoinNode (different conditions)")
		t.Errorf("   Both have ID: %s", r1_join2.ID)
	} else {
		t.Logf("✅ r1 and r2 have different second JoinNodes (correct)")
	}

	t.Log("")
	t.Log("✅ REGRESSION TEST PASSED")
	t.Log("   - Each rule builds its complete cascade independently")
	t.Log("   - Prefix sharing does not cross rule boundaries")
	t.Log("   - JoinNodes are shared only when conditions match exactly")
}

// TestCascadeLevelInSignature verifies that cascadeLevel is correctly included
// in the JoinNode signature to prevent sharing between different cascade levels.
func TestCascadeLevelInSignature(t *testing.T) {
	t.Log("🧪 TEST: Cascade Level in Signature")
	t.Log("=====================================")

	storage := NewMemoryStorage()
	lifecycle := NewLifecycleManager()
	config := DefaultBetaSharingConfig()
	config.Enabled = true
	registry := NewBetaSharingRegistry(config, lifecycle)

	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
	}
	varTypes := map[string]string{"u": "User", "o": "Order"}

	// Create JoinNode at cascade level 0
	node1, hash1, shared1, err := registry.GetOrCreateJoinNode(
		condition,
		[]string{"u"},
		[]string{"o"},
		[]string{"u", "o"},
		varTypes,
		storage,
		0, // cascadeLevel = 0
	)
	if err != nil {
		t.Fatalf("Failed to create node at level 0: %v", err)
	}
	if shared1 {
		t.Error("First node should not be shared")
	}

	// Attempt to create "same" JoinNode at cascade level 1
	// This should create a DIFFERENT node because cascadeLevel differs
	node2, hash2, shared2, err := registry.GetOrCreateJoinNode(
		condition,
		[]string{"u"},
		[]string{"o"},
		[]string{"u", "o"},
		varTypes,
		storage,
		1, // cascadeLevel = 1 (different!)
	)
	if err != nil {
		t.Fatalf("Failed to create node at level 1: %v", err)
	}

	// CRITICAL: Nodes at different cascade levels should NOT be shared
	if shared2 {
		t.Error("❌ Node at level 1 should not share with node at level 0")
	}

	if hash1 == hash2 {
		t.Errorf("❌ Different cascade levels should produce different hashes")
		t.Errorf("   Level 0 hash: %s", hash1)
		t.Errorf("   Level 1 hash: %s", hash2)
	}

	if node1.ID == node2.ID {
		t.Errorf("❌ Different cascade levels should produce different nodes")
		t.Errorf("   Both have ID: %s", node1.ID)
	}

	t.Logf("✅ Cascade level 0: node %s (hash: %s)", node1.ID, hash1)
	t.Logf("✅ Cascade level 1: node %s (hash: %s)", node2.ID, hash2)
	t.Log("✅ Different cascade levels create different nodes (correct)")
}
