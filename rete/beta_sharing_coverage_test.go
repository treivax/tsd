// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"
)

// TestBetaSharingRegistry_AddRuleToJoinNode tests adding rules to join nodes
func TestBetaSharingRegistry_AddRuleToJoinNode(t *testing.T) {
	t.Run("add rule when sharing enabled", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		err := registry.AddRuleToJoinNode("node1", "rule1")
		if err != nil {
			t.Errorf("AddRuleToJoinNode() unexpected error: %v", err)
		}
		// Verify rule was added
		rules := registry.GetJoinNodeRules("node1")
		if len(rules) != 1 {
			t.Errorf("Expected 1 rule, got %d", len(rules))
		}
		if rules[0] != "rule1" {
			t.Errorf("Expected rule1, got %s", rules[0])
		}
	})
	t.Run("add multiple rules to same node", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		err := registry.AddRuleToJoinNode("node1", "rule1")
		if err != nil {
			t.Errorf("AddRuleToJoinNode() error: %v", err)
		}
		err = registry.AddRuleToJoinNode("node1", "rule2")
		if err != nil {
			t.Errorf("AddRuleToJoinNode() error: %v", err)
		}
		err = registry.AddRuleToJoinNode("node1", "rule3")
		if err != nil {
			t.Errorf("AddRuleToJoinNode() error: %v", err)
		}
		rules := registry.GetJoinNodeRules("node1")
		if len(rules) != 3 {
			t.Errorf("Expected 3 rules, got %d", len(rules))
		}
	})
	t.Run("add rule when sharing disabled", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = false
		registry := NewBetaSharingRegistry(config, nil)
		err := registry.AddRuleToJoinNode("node1", "rule1")
		if err != nil {
			t.Errorf("AddRuleToJoinNode() unexpected error: %v", err)
		}
		// Should be no-op when disabled
		rules := registry.GetJoinNodeRules("node1")
		if len(rules) != 0 {
			t.Errorf("Expected 0 rules when sharing disabled, got %d", len(rules))
		}
	})
	t.Run("add same rule twice to same node", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		err := registry.AddRuleToJoinNode("node1", "rule1")
		if err != nil {
			t.Errorf("AddRuleToJoinNode() error: %v", err)
		}
		err = registry.AddRuleToJoinNode("node1", "rule1")
		if err != nil {
			t.Errorf("AddRuleToJoinNode() error: %v", err)
		}
		// Should still have only 1 rule (map deduplicates)
		rules := registry.GetJoinNodeRules("node1")
		if len(rules) != 1 {
			t.Errorf("Expected 1 rule (deduplicated), got %d", len(rules))
		}
	})
}

// TestBetaSharingRegistry_RemoveRuleFromJoinNode tests removing rules from join nodes
func TestBetaSharingRegistry_RemoveRuleFromJoinNode(t *testing.T) {
	t.Run("remove rule from node with multiple rules", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		registry.AddRuleToJoinNode("node1", "rule1")
		registry.AddRuleToJoinNode("node1", "rule2")
		registry.AddRuleToJoinNode("node1", "rule3")
		canDelete, err := registry.RemoveRuleFromJoinNode("node1", "rule2")
		if err != nil {
			t.Errorf("RemoveRuleFromJoinNode() error: %v", err)
		}
		if canDelete {
			t.Error("Node should not be deletable, still has rules")
		}
		rules := registry.GetJoinNodeRules("node1")
		if len(rules) != 2 {
			t.Errorf("Expected 2 rules remaining, got %d", len(rules))
		}
	})
	t.Run("remove last rule from node", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		registry.AddRuleToJoinNode("node1", "rule1")
		canDelete, err := registry.RemoveRuleFromJoinNode("node1", "rule1")
		if err != nil {
			t.Errorf("RemoveRuleFromJoinNode() error: %v", err)
		}
		if !canDelete {
			t.Error("Node should be deletable when last rule is removed")
		}
		rules := registry.GetJoinNodeRules("node1")
		if len(rules) != 0 {
			t.Errorf("Expected 0 rules, got %d", len(rules))
		}
	})
	t.Run("remove rule from non-existent node", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		canDelete, err := registry.RemoveRuleFromJoinNode("nonexistent", "rule1")
		if err == nil {
			t.Error("Expected error when removing from non-existent node")
		}
		if canDelete {
			t.Error("Should not report node as deletable when it doesn't exist")
		}
	})
	t.Run("remove non-existent rule from node", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		registry.AddRuleToJoinNode("node1", "rule1")
		_, err := registry.RemoveRuleFromJoinNode("node1", "rule2")
		if err != nil {
			t.Errorf("RemoveRuleFromJoinNode() error: %v", err)
		}
		// Node should still have rule1
		rules := registry.GetJoinNodeRules("node1")
		if len(rules) != 1 {
			t.Errorf("Expected 1 rule remaining, got %d", len(rules))
		}
	})
	t.Run("remove rule when sharing disabled", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = false
		registry := NewBetaSharingRegistry(config, nil)
		canDelete, err := registry.RemoveRuleFromJoinNode("node1", "rule1")
		if err != nil {
			t.Errorf("RemoveRuleFromJoinNode() unexpected error: %v", err)
		}
		if canDelete {
			t.Error("Should return false when sharing disabled")
		}
	})
}

// TestBetaSharingRegistry_GetJoinNodeRules tests retrieving rules for a node
func TestBetaSharingRegistry_GetJoinNodeRules(t *testing.T) {
	t.Run("get rules for node with rules", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		registry.AddRuleToJoinNode("node1", "rule1")
		registry.AddRuleToJoinNode("node1", "rule2")
		registry.AddRuleToJoinNode("node1", "rule3")
		rules := registry.GetJoinNodeRules("node1")
		if len(rules) != 3 {
			t.Errorf("Expected 3 rules, got %d", len(rules))
		}
		// Verify all rules are present
		ruleMap := make(map[string]bool)
		for _, rule := range rules {
			ruleMap[rule] = true
		}
		if !ruleMap["rule1"] || !ruleMap["rule2"] || !ruleMap["rule3"] {
			t.Error("Not all expected rules were returned")
		}
	})
	t.Run("get rules for non-existent node", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		rules := registry.GetJoinNodeRules("nonexistent")
		if len(rules) != 0 {
			t.Errorf("Expected empty slice for non-existent node, got %d rules", len(rules))
		}
	})
	t.Run("get rules for node with no rules", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		// Add then remove all rules
		registry.AddRuleToJoinNode("node1", "rule1")
		registry.RemoveRuleFromJoinNode("node1", "rule1")
		rules := registry.GetJoinNodeRules("node1")
		if len(rules) != 0 {
			t.Errorf("Expected empty slice after removing all rules, got %d rules", len(rules))
		}
	})
}

// TestBetaSharingRegistry_GetJoinNodeRefCount tests reference counting
func TestBetaSharingRegistry_GetJoinNodeRefCount(t *testing.T) {
	t.Run("refcount for node with multiple rules", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		registry.AddRuleToJoinNode("node1", "rule1")
		registry.AddRuleToJoinNode("node1", "rule2")
		registry.AddRuleToJoinNode("node1", "rule3")
		count := registry.GetJoinNodeRefCount("node1")
		if count != 3 {
			t.Errorf("Expected refcount 3, got %d", count)
		}
	})
	t.Run("refcount for non-existent node", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		count := registry.GetJoinNodeRefCount("nonexistent")
		if count != 0 {
			t.Errorf("Expected refcount 0 for non-existent node, got %d", count)
		}
	})
	t.Run("refcount after removing rules", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		registry.AddRuleToJoinNode("node1", "rule1")
		registry.AddRuleToJoinNode("node1", "rule2")
		registry.AddRuleToJoinNode("node1", "rule3")
		registry.RemoveRuleFromJoinNode("node1", "rule2")
		count := registry.GetJoinNodeRefCount("node1")
		if count != 2 {
			t.Errorf("Expected refcount 2 after removal, got %d", count)
		}
	})
	t.Run("refcount zero after removing all rules", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		registry.AddRuleToJoinNode("node1", "rule1")
		registry.RemoveRuleFromJoinNode("node1", "rule1")
		count := registry.GetJoinNodeRefCount("node1")
		if count != 0 {
			t.Errorf("Expected refcount 0 after removing all rules, got %d", count)
		}
	})
}

// TestBetaSharingRegistry_UnregisterJoinNode tests unregistering nodes
func TestBetaSharingRegistry_UnregisterJoinNode(t *testing.T) {
	t.Run("unregister node with no rules", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		registry.AddRuleToJoinNode("node1", "rule1")
		registry.RemoveRuleFromJoinNode("node1", "rule1")
		err := registry.UnregisterJoinNode("node1")
		if err != nil {
			t.Errorf("UnregisterJoinNode() error: %v", err)
		}
		// Verify node is gone
		count := registry.GetJoinNodeRefCount("node1")
		if count != 0 {
			t.Errorf("Expected node to be unregistered, got refcount %d", count)
		}
	})
	t.Run("unregister non-existent node", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		err := registry.UnregisterJoinNode("nonexistent")
		if err != nil {
			t.Errorf("UnregisterJoinNode() should not error for non-existent node: %v", err)
		}
	})
	t.Run("unregister when sharing disabled", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = false
		registry := NewBetaSharingRegistry(config, nil)
		err := registry.UnregisterJoinNode("node1")
		if err != nil {
			t.Errorf("UnregisterJoinNode() error: %v", err)
		}
	})
	t.Run("unregister node then try to get rules", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		registry.AddRuleToJoinNode("node1", "rule1")
		registry.RemoveRuleFromJoinNode("node1", "rule1")
		registry.UnregisterJoinNode("node1")
		rules := registry.GetJoinNodeRules("node1")
		if len(rules) != 0 {
			t.Errorf("Expected no rules after unregister, got %d", len(rules))
		}
	})
}

// TestBetaSharingRegistry_ReleaseJoinNodeByID tests releasing nodes by ID
func TestBetaSharingRegistry_ReleaseJoinNodeByID(t *testing.T) {
	t.Run("release node with no rules", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		storage := NewMemoryStorage()
		condition := map[string]interface{}{
			"type":     "comparison",
			"operator": "==",
			"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "id"},
			"right":    map[string]interface{}{"type": "fieldAccess", "object": "o", "field": "id"},
		}
		node, _, _, err := registry.GetOrCreateJoinNode(
			condition,
			[]string{"p"},
			[]string{"o"},
			[]string{"p", "o"},
			map[string]string{"p": "Person", "o": "Order"},
			storage,
		)
		if err != nil {
			t.Fatalf("GetOrCreateJoinNode() error: %v", err)
		}
		nodeID := node.ID
		// Release should succeed when no rules reference it
		released, err := registry.ReleaseJoinNodeByID(nodeID)
		if err != nil {
			t.Errorf("ReleaseJoinNodeByID() error: %v", err)
		}
		if !released {
			t.Error("Expected node to be released")
		}
	})
	t.Run("release node with active rules", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		storage := NewMemoryStorage()
		condition := map[string]interface{}{
			"type":     "comparison",
			"operator": "==",
			"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "id"},
			"right":    map[string]interface{}{"type": "fieldAccess", "object": "o", "field": "id"},
		}
		node, _, _, err := registry.GetOrCreateJoinNode(
			condition,
			[]string{"p"},
			[]string{"o"},
			[]string{"p", "o"},
			map[string]string{"p": "Person", "o": "Order"},
			storage,
		)
		if err != nil {
			t.Fatalf("GetOrCreateJoinNode() error: %v", err)
		}
		nodeID := node.ID
		registry.AddRuleToJoinNode(nodeID, "rule1")
		// Release should fail when rules still reference it
		released, err := registry.ReleaseJoinNodeByID(nodeID)
		if err == nil {
			t.Error("Expected error when releasing node with active rules")
		}
		if released {
			t.Error("Node should not be released when rules are active")
		}
	})
	t.Run("release non-existent node", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = true
		registry := NewBetaSharingRegistry(config, nil)
		released, err := registry.ReleaseJoinNodeByID("nonexistent")
		if err != nil {
			t.Errorf("ReleaseJoinNodeByID() should not error for non-existent node: %v", err)
		}
		if released {
			t.Error("Should not report non-existent node as released")
		}
	})
	t.Run("release when sharing disabled", func(t *testing.T) {
		config := DefaultBetaSharingConfig()
		config.Enabled = false
		registry := NewBetaSharingRegistry(config, nil)
		released, err := registry.ReleaseJoinNodeByID("node1")
		if err != nil {
			t.Errorf("ReleaseJoinNodeByID() error: %v", err)
		}
		if released {
			t.Error("Should return false when sharing disabled")
		}
	})
}

// TestBetaSharingRegistry_RuleLifecycle tests complete rule lifecycle
func TestBetaSharingRegistry_RuleLifecycle(t *testing.T) {
	config := DefaultBetaSharingConfig()
	config.Enabled = true
	registry := NewBetaSharingRegistry(config, nil)
	// Add multiple rules to a node
	registry.AddRuleToJoinNode("node1", "rule1")
	registry.AddRuleToJoinNode("node1", "rule2")
	registry.AddRuleToJoinNode("node1", "rule3")
	// Verify refcount
	if count := registry.GetJoinNodeRefCount("node1"); count != 3 {
		t.Errorf("Expected refcount 3, got %d", count)
	}
	// Remove one rule
	canDelete, err := registry.RemoveRuleFromJoinNode("node1", "rule2")
	if err != nil {
		t.Fatalf("RemoveRuleFromJoinNode() error: %v", err)
	}
	if canDelete {
		t.Error("Node should not be deletable yet")
	}
	// Verify refcount decreased
	if count := registry.GetJoinNodeRefCount("node1"); count != 2 {
		t.Errorf("Expected refcount 2 after removal, got %d", count)
	}
	// Remove another rule
	canDelete, err = registry.RemoveRuleFromJoinNode("node1", "rule1")
	if err != nil {
		t.Fatalf("RemoveRuleFromJoinNode() error: %v", err)
	}
	if canDelete {
		t.Error("Node should not be deletable yet")
	}
	// Remove last rule
	canDelete, err = registry.RemoveRuleFromJoinNode("node1", "rule3")
	if err != nil {
		t.Fatalf("RemoveRuleFromJoinNode() error: %v", err)
	}
	if !canDelete {
		t.Error("Node should be deletable after removing last rule")
	}
	// Verify node has no rules
	if count := registry.GetJoinNodeRefCount("node1"); count != 0 {
		t.Errorf("Expected refcount 0 after removing all rules, got %d", count)
	}
	// Unregister the node
	err = registry.UnregisterJoinNode("node1")
	if err != nil {
		t.Errorf("UnregisterJoinNode() error: %v", err)
	}
}

// TestBetaSharingRegistry_ConcurrentRuleManagement tests thread safety
func TestBetaSharingRegistry_ConcurrentRuleManagement(t *testing.T) {
	config := DefaultBetaSharingConfig()
	config.Enabled = true
	registry := NewBetaSharingRegistry(config, nil)
	done := make(chan bool)
	numGoroutines := 10
	rulesPerGoroutine := 5
	// Concurrently add rules
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < rulesPerGoroutine; j++ {
				nodeID := "node1"
				ruleID := string(rune('A' + id*rulesPerGoroutine + j))
				registry.AddRuleToJoinNode(nodeID, ruleID)
			}
			done <- true
		}(i)
	}
	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	// Verify all rules were added
	count := registry.GetJoinNodeRefCount("node1")
	if count != numGoroutines*rulesPerGoroutine {
		t.Errorf("Expected %d rules, got %d", numGoroutines*rulesPerGoroutine, count)
	}
	// Concurrently query rules
	for i := 0; i < numGoroutines; i++ {
		go func() {
			rules := registry.GetJoinNodeRules("node1")
			if len(rules) == 0 {
				t.Error("Expected rules but got empty slice")
			}
			done <- true
		}()
	}
	// Wait for all queries
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}
