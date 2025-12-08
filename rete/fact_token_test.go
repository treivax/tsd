// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWorkingMemory_RemoveToken tests removing a token from working memory
func TestWorkingMemory_RemoveToken(t *testing.T) {
	wm := &WorkingMemory{
		NodeID: "test_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	// Add some tokens
	token1 := &Token{ID: "token1", NodeID: "test_node"}
	token2 := &Token{ID: "token2", NodeID: "test_node"}
	wm.Tokens["token1"] = token1
	wm.Tokens["token2"] = token2

	assert.Len(t, wm.Tokens, 2, "Should have 2 tokens initially")

	// Remove one token
	wm.RemoveToken("token1")

	assert.Len(t, wm.Tokens, 1, "Should have 1 token after removal")
	assert.Nil(t, wm.Tokens["token1"], "token1 should be removed")
	assert.NotNil(t, wm.Tokens["token2"], "token2 should still exist")

	// Remove non-existent token (should not panic)
	wm.RemoveToken("non_existent")
	assert.Len(t, wm.Tokens, 1, "Should still have 1 token")
}

// TestWorkingMemory_GetFactsByVariable tests getting facts by variable
func TestWorkingMemory_GetFactsByVariable(t *testing.T) {
	wm := &WorkingMemory{
		NodeID: "test_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	// Add some facts
	fact1 := &Fact{ID: "fact1", Type: "Person", Fields: map[string]interface{}{"name": "Alice"}}
	fact2 := &Fact{ID: "fact2", Type: "Order", Fields: map[string]interface{}{"id": "123"}}
	wm.Facts["fact1"] = fact1
	wm.Facts["fact2"] = fact2

	// Get facts by variable (current implementation returns all facts)
	facts := wm.GetFactsByVariable([]string{"p", "o"})

	assert.Len(t, facts, 2, "Should return all facts")
	assert.Contains(t, facts, fact1, "Should contain fact1")
	assert.Contains(t, facts, fact2, "Should contain fact2")
}

// TestWorkingMemory_GetTokensByVariable tests getting tokens by variable
func TestWorkingMemory_GetTokensByVariable(t *testing.T) {
	wm := &WorkingMemory{
		NodeID: "test_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	// Add some tokens
	token1 := &Token{ID: "token1", NodeID: "test_node"}
	token2 := &Token{ID: "token2", NodeID: "test_node"}
	wm.Tokens["token1"] = token1
	wm.Tokens["token2"] = token2

	// Get tokens by variable (current implementation returns all tokens)
	tokens := wm.GetTokensByVariable([]string{"p", "o"})

	assert.Len(t, tokens, 2, "Should return all tokens")
	assert.Contains(t, tokens, token1, "Should contain token1")
	assert.Contains(t, tokens, token2, "Should contain token2")
}

// TestFact_Clone tests cloning a fact
func TestFact_Clone(t *testing.T) {
	original := &Fact{
		ID:     "fact1",
		Type:   "Person",
		Fields: map[string]interface{}{"name": "Alice", "age": 30},
	}

	// Clone the fact
	clone := original.Clone()

	// Verify clone is equal but independent
	require.NotNil(t, clone, "Clone should not be nil")
	assert.Equal(t, original.ID, clone.ID, "ID should match")
	assert.Equal(t, original.Type, clone.Type, "Type should match")

	assert.Equal(t, original.Fields["name"], clone.Fields["name"], "Fields should match")
	assert.Equal(t, original.Fields["age"], clone.Fields["age"], "Fields should match")

	// Modify clone and verify original is unchanged
	clone.Fields["name"] = "Bob"
	clone.Fields["age"] = 40

	assert.Equal(t, "Alice", original.Fields["name"], "Original should be unchanged")
	assert.Equal(t, 30, original.Fields["age"], "Original should be unchanged")
	assert.Equal(t, "Bob", clone.Fields["name"], "Clone should have new value")
	assert.Equal(t, 40, clone.Fields["age"], "Clone should have new value")
}

// TestFact_CloneWithEmptyFields tests cloning a fact with empty fields
func TestFact_CloneWithEmptyFields(t *testing.T) {
	original := &Fact{
		ID:     "fact1",
		Type:   "Empty",
		Fields: map[string]interface{}{},
	}

	clone := original.Clone()

	require.NotNil(t, clone, "Clone should not be nil")
	assert.Equal(t, original.ID, clone.ID)
	assert.Equal(t, original.Type, clone.Type)
	assert.NotNil(t, clone.Fields, "Fields should be initialized")
	assert.Empty(t, clone.Fields, "Fields should be empty")
}

// TestWorkingMemory_Clone tests cloning working memory
func TestWorkingMemory_Clone(t *testing.T) {
	original := &WorkingMemory{
		NodeID: "test_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	// Add facts
	fact1 := &Fact{ID: "fact1", Type: "Person", Fields: map[string]interface{}{"name": "Alice"}}
	fact2 := &Fact{ID: "fact2", Type: "Order", Fields: map[string]interface{}{"id": "123"}}
	original.Facts["fact1"] = fact1
	original.Facts["fact2"] = fact2

	// Add tokens
	token1 := &Token{ID: "token1", NodeID: "test_node", Facts: []*Fact{fact1}}
	token2 := &Token{ID: "token2", NodeID: "test_node", Facts: []*Fact{fact2}}
	original.Tokens["token1"] = token1
	original.Tokens["token2"] = token2

	// Clone working memory
	clone := original.Clone()

	// Verify clone is equal but independent
	require.NotNil(t, clone, "Clone should not be nil")
	assert.Equal(t, original.NodeID, clone.NodeID, "NodeID should match")
	assert.Len(t, clone.Facts, 2, "Should have 2 facts")
	assert.Len(t, clone.Tokens, 2, "Should have 2 tokens")

	// Verify facts are cloned
	clonedFact1 := clone.Facts["fact1"]
	require.NotNil(t, clonedFact1, "Cloned fact should exist")
	assert.Equal(t, fact1.ID, clonedFact1.ID, "Fact ID should match")
	assert.Equal(t, fact1.Type, clonedFact1.Type, "Fact type should match")

	// Modify clone and verify original is unchanged
	clonedFact1.Fields["name"] = "Bob"
	assert.Equal(t, "Alice", fact1.Fields["name"], "Original fact should be unchanged")
	assert.Equal(t, "Bob", clonedFact1.Fields["name"], "Cloned fact should have new value")

	// Verify tokens are cloned
	clonedToken1 := clone.Tokens["token1"]
	require.NotNil(t, clonedToken1, "Cloned token should exist")
	assert.Equal(t, token1.ID, clonedToken1.ID, "Token ID should match")
}

// TestWorkingMemory_CloneEmpty tests cloning empty working memory
func TestWorkingMemory_CloneEmpty(t *testing.T) {
	original := &WorkingMemory{
		NodeID: "empty_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	clone := original.Clone()

	require.NotNil(t, clone, "Clone should not be nil")
	assert.Equal(t, original.NodeID, clone.NodeID)
	assert.NotNil(t, clone.Facts, "Facts map should be initialized")
	assert.NotNil(t, clone.Tokens, "Tokens map should be initialized")
	assert.Empty(t, clone.Facts, "Facts should be empty")
	assert.Empty(t, clone.Tokens, "Tokens should be empty")
}

// TestToken_Clone tests cloning a token
func TestToken_Clone(t *testing.T) {
	fact1 := &Fact{ID: "fact1", Type: "Person", Fields: map[string]interface{}{"name": "Alice"}}
	fact2 := &Fact{ID: "fact2", Type: "Order", Fields: map[string]interface{}{"id": "123"}}

	original := &Token{
		ID:           "token1",
		Facts:        []*Fact{fact1, fact2},
		NodeID:       "test_node",
		Bindings:     map[string]*Fact{"p": fact1, "o": fact2},
		IsJoinResult: true,
	}

	// Clone the token
	clone := original.Clone()

	// Verify clone is equal
	require.NotNil(t, clone, "Clone should not be nil")
	assert.Equal(t, original.ID, clone.ID, "ID should match")
	assert.Equal(t, original.NodeID, clone.NodeID, "NodeID should match")
	assert.Equal(t, original.IsJoinResult, clone.IsJoinResult, "IsJoinResult should match")
	assert.Len(t, clone.Facts, 2, "Should have 2 facts")
	assert.Len(t, clone.Bindings, 2, "Should have 2 bindings")

	// Verify facts are copied (shallow copy)
	assert.Equal(t, fact1, clone.Facts[0], "First fact should match")
	assert.Equal(t, fact2, clone.Facts[1], "Second fact should match")

	// Verify bindings are copied
	assert.Equal(t, fact1, clone.Bindings["p"], "Binding 'p' should match")
	assert.Equal(t, fact2, clone.Bindings["o"], "Binding 'o' should match")
}

// TestToken_CloneEmpty tests cloning an empty token
func TestToken_CloneEmpty(t *testing.T) {
	original := &Token{
		ID:       "token1",
		Facts:    []*Fact{},
		NodeID:   "test_node",
		Bindings: map[string]*Fact{},
	}

	clone := original.Clone()

	require.NotNil(t, clone, "Clone should not be nil")
	assert.Equal(t, original.ID, clone.ID)
	assert.Equal(t, original.NodeID, clone.NodeID)
	assert.NotNil(t, clone.Facts, "Facts should be initialized")
	assert.NotNil(t, clone.Bindings, "Bindings should be initialized")
	assert.Empty(t, clone.Facts, "Facts should be empty")
	assert.Empty(t, clone.Bindings, "Bindings should be empty")
}

// TestToken_CloneIndependence tests that cloned token is independent
func TestToken_CloneIndependence(t *testing.T) {
	fact := &Fact{ID: "fact1", Type: "Person", Fields: map[string]interface{}{"name": "Alice"}}

	original := &Token{
		ID:       "token1",
		Facts:    []*Fact{fact},
		NodeID:   "test_node",
		Bindings: map[string]*Fact{"p": fact},
	}

	clone := original.Clone()

	// Modify clone's slices and maps
	newFact := &Fact{ID: "fact2", Type: "Order"}
	clone.Facts = append(clone.Facts, newFact)
	clone.Bindings["o"] = newFact

	// Verify original is unchanged
	assert.Len(t, original.Facts, 1, "Original should still have 1 fact")
	assert.Len(t, original.Bindings, 1, "Original should still have 1 binding")
	assert.Nil(t, original.Bindings["o"], "Original should not have 'o' binding")

	// Verify clone has new values
	assert.Len(t, clone.Facts, 2, "Clone should have 2 facts")
	assert.Len(t, clone.Bindings, 2, "Clone should have 2 bindings")
	assert.NotNil(t, clone.Bindings["o"], "Clone should have 'o' binding")
}

// TestWorkingMemory_ComplexClone tests cloning with complex nested structures
func TestWorkingMemory_ComplexClone(t *testing.T) {
	wm := &WorkingMemory{
		NodeID: "complex_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	// Add facts with complex field values
	fact1 := &Fact{
		ID:   "fact1",
		Type: "Person",
		Fields: map[string]interface{}{
			"name":    "Alice",
			"age":     30,
			"active":  true,
			"nested":  map[string]interface{}{"city": "Paris", "country": "France"},
			"tags":    []string{"premium", "verified"},
			"balance": 1250.50,
		},
	}
	wm.Facts["fact1"] = fact1

	// Clone
	clone := wm.Clone()

	// Verify clone
	require.NotNil(t, clone)
	clonedFact := clone.Facts["fact1"]
	require.NotNil(t, clonedFact)

	assert.Equal(t, "Alice", clonedFact.Fields["name"])
	assert.Equal(t, 30, clonedFact.Fields["age"])
	assert.Equal(t, true, clonedFact.Fields["active"])
	assert.Equal(t, 1250.50, clonedFact.Fields["balance"])

	// Note: nested maps and slices are shallow copied in the current implementation
	// This test documents the current behavior
	nestedMap := clonedFact.Fields["nested"].(map[string]interface{})
	assert.Equal(t, "Paris", nestedMap["city"])
}
