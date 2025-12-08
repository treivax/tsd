// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

func TestNewAlphaNode(t *testing.T) {
	storage := NewMemoryStorage()
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
		"left":     "field1",
		"right":    "value1",
	}

	node := NewAlphaNode("alpha1", condition, "p", storage)

	if node == nil {
		t.Fatal("NewAlphaNode returned nil")
	}

	if node.ID != "alpha1" {
		t.Errorf("Expected ID 'alpha1', got '%s'", node.ID)
	}

	if node.Type != "alpha" {
		t.Errorf("Expected Type 'alpha', got '%s'", node.Type)
	}

	if node.VariableName != "p" {
		t.Errorf("Expected VariableName 'p', got '%s'", node.VariableName)
	}

	if node.Condition == nil {
		t.Error("Condition should not be nil")
	}

	if node.Memory == nil {
		t.Error("Memory should not be nil")
	}

	if node.Children == nil {
		t.Error("Children should not be nil")
	}
}

func TestAlphaNodeActivateLeft(t *testing.T) {
	storage := NewMemoryStorage()
	node := NewAlphaNode("alpha1", nil, "p", storage)

	token := &Token{
		ID:    "t1",
		Facts: []*Fact{},
	}

	// ActivateLeft now silently ignores tokens (used during retroactive propagation)
	err := node.ActivateLeft(token)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestAlphaNodeActivateRetract(t *testing.T) {
	storage := NewMemoryStorage()
	node := NewAlphaNode("alpha1", nil, "p", storage)

	// Add a fact to memory
	fact := &Fact{
		ID:     "f1",
		Type:   "Person",
		Fields: map[string]interface{}{"name": "Alice"},
	}

	node.Memory.Facts[fact.ID] = fact

	// Verify fact exists
	_, exists := node.Memory.Facts[fact.ID]
	if !exists {
		t.Fatal("Fact should exist in memory before retraction")
	}

	// Retract the fact
	err := node.ActivateRetract(fact.ID)
	if err != nil {
		t.Errorf("ActivateRetract failed: %v", err)
	}

	// Verify fact was removed
	_, exists = node.Memory.Facts[fact.ID]
	if exists {
		t.Error("Fact should have been removed from memory")
	}
}

func TestAlphaNodeActivateRetractNonExistent(t *testing.T) {
	storage := NewMemoryStorage()
	node := NewAlphaNode("alpha1", nil, "p", storage)

	// Try to retract a fact that doesn't exist
	err := node.ActivateRetract("nonexistent")
	if err != nil {
		t.Errorf("ActivateRetract should not error for non-existent fact: %v", err)
	}
}

func TestAlphaNodePassthroughLeft(t *testing.T) {
	storage := NewMemoryStorage()
	condition := map[string]interface{}{
		"type": "passthrough",
		"side": "left",
	}

	node := NewAlphaNode("alpha1", condition, "p", storage)

	// Add a child node to receive propagation
	childCalled := false
	mockChild := &MockNode{
		activateLeft: func(token *Token) error {
			childCalled = true
			if token == nil {
				t.Error("Token should not be nil")
				return nil
			}
			if len(token.Facts) != 1 {
				t.Errorf("Expected 1 fact in token, got %d", len(token.Facts))
			}
			if token.Bindings["p"] == nil {
				t.Error("Expected binding for variable 'p'")
			}
			return nil
		},
	}
	node.AddChild(mockChild)

	fact := &Fact{
		ID:     "f1",
		Type:   "Person",
		Fields: map[string]interface{}{"name": "Alice"},
	}

	err := node.ActivateRight(fact)
	if err != nil {
		t.Errorf("ActivateRight failed: %v", err)
	}

	if !childCalled {
		t.Error("Child node should have been called")
	}
}

func TestAlphaNodePassthroughRight(t *testing.T) {
	storage := NewMemoryStorage()
	condition := map[string]interface{}{
		"type": "passthrough",
		"side": "right",
	}

	node := NewAlphaNode("alpha1", condition, "p", storage)

	// Add a child node to receive propagation
	childCalled := false
	mockChild := &MockNode{
		activateRight: func(fact *Fact) error {
			childCalled = true
			if fact == nil {
				t.Error("Fact should not be nil")
			}
			return nil
		},
	}
	node.AddChild(mockChild)

	fact := &Fact{
		ID:     "f1",
		Type:   "Person",
		Fields: map[string]interface{}{"name": "Bob"},
	}

	err := node.ActivateRight(fact)
	if err != nil {
		t.Errorf("ActivateRight failed: %v", err)
	}

	if !childCalled {
		t.Error("Child node should have been called with ActivateRight")
	}
}

func TestAlphaNodePassthroughDefault(t *testing.T) {
	storage := NewMemoryStorage()
	condition := map[string]interface{}{
		"type": "passthrough",
		// No side specified, should default to right
	}

	node := NewAlphaNode("alpha1", condition, "p", storage)

	childCalled := false
	mockChild := &MockNode{
		activateRight: func(fact *Fact) error {
			childCalled = true
			return nil
		},
	}
	node.AddChild(mockChild)

	fact := &Fact{
		ID:     "f1",
		Type:   "Person",
		Fields: map[string]interface{}{"name": "Charlie"},
	}

	err := node.ActivateRight(fact)
	if err != nil {
		t.Errorf("ActivateRight failed: %v", err)
	}

	if !childCalled {
		t.Error("Child node should have been called")
	}
}

func TestAlphaNodeRetractWithChildren(t *testing.T) {
	storage := NewMemoryStorage()
	node := NewAlphaNode("alpha1", nil, "p", storage)

	// Add child node
	childRetractCalled := false
	mockChild := &MockNode{
		activateRetract: func(factID string) error {
			childRetractCalled = true
			if factID != "f1" {
				t.Errorf("Expected factID 'f1', got '%s'", factID)
			}
			return nil
		},
	}
	node.AddChild(mockChild)

	// Add fact
	fact := &Fact{
		ID:     "f1",
		Type:   "Person",
		Fields: map[string]interface{}{"name": "Dave"},
	}
	node.Memory.Facts[fact.ID] = fact

	// Retract
	err := node.ActivateRetract(fact.ID)
	if err != nil {
		t.Errorf("ActivateRetract failed: %v", err)
	}

	if !childRetractCalled {
		t.Error("Child node's ActivateRetract should have been called")
	}
}

func TestAlphaNodeMemoryIsolation(t *testing.T) {
	storage := NewMemoryStorage()
	node1 := NewAlphaNode("alpha1", nil, "p", storage)
	node2 := NewAlphaNode("alpha2", nil, "q", storage)

	fact1 := &Fact{ID: "f1", Type: "Person", Fields: map[string]interface{}{"name": "Alice"}}
	fact2 := &Fact{ID: "f2", Type: "Person", Fields: map[string]interface{}{"name": "Bob"}}

	node1.Memory.Facts[fact1.ID] = fact1
	node2.Memory.Facts[fact2.ID] = fact2

	// Verify isolation
	_, exists1 := node1.Memory.Facts["f1"]
	if !exists1 {
		t.Error("Node1 should have fact f1")
	}

	_, exists2 := node1.Memory.Facts["f2"]
	if exists2 {
		t.Error("Node1 should not have fact f2")
	}

	_, exists3 := node2.Memory.Facts["f2"]
	if !exists3 {
		t.Error("Node2 should have fact f2")
	}

	_, exists4 := node2.Memory.Facts["f1"]
	if exists4 {
		t.Error("Node2 should not have fact f1")
	}
}

// MockNode for testing
type MockNode struct {
	BaseNode
	activateLeft    func(*Token) error
	activateRight   func(*Fact) error
	activateRetract func(string) error
}

func (m *MockNode) ActivateLeft(token *Token) error {
	if m.activateLeft != nil {
		return m.activateLeft(token)
	}
	return nil
}

func (m *MockNode) ActivateRight(fact *Fact) error {
	if m.activateRight != nil {
		return m.activateRight(fact)
	}
	return nil
}

func (m *MockNode) ActivateRetract(factID string) error {
	if m.activateRetract != nil {
		return m.activateRetract(factID)
	}
	return nil
}
