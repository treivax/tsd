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
	fact3 := &Fact{ID: "fact3", Type: "Product", Fields: map[string]interface{}{"sku": "ABC"}}
	wm.Facts["fact1"] = fact1
	wm.Facts["fact2"] = fact2
	wm.Facts["fact3"] = fact3

	// Create tokens with bindings
	bindings1 := NewBindingChain().Add("p", fact1)
	token1 := &Token{ID: "token1", Bindings: bindings1}
	wm.Tokens["token1"] = token1

	bindings2 := NewBindingChain().Add("o", fact2)
	token2 := &Token{ID: "token2", Bindings: bindings2}
	wm.Tokens["token2"] = token2

	bindings3 := NewBindingChain().Add("prod", fact3)
	token3 := &Token{ID: "token3", Bindings: bindings3}
	wm.Tokens["token3"] = token3

	// Test 1: Get facts by specific variables
	facts := wm.GetFactsByVariable([]string{"p", "o"})
	assert.Len(t, facts, 2, "Should return facts for variables p and o")
	assert.Contains(t, facts, fact1, "Should contain fact1 (variable p)")
	assert.Contains(t, facts, fact2, "Should contain fact2 (variable o)")

	// Test 2: Get facts by single variable
	factsP := wm.GetFactsByVariable([]string{"p"})
	assert.Len(t, factsP, 1, "Should return only fact for variable p")
	assert.Contains(t, factsP, fact1, "Should contain fact1")

	// Test 3: Get facts by non-existent variable
	factsNone := wm.GetFactsByVariable([]string{"nonexistent"})
	assert.Len(t, factsNone, 0, "Should return empty slice for non-existent variable")

	// Test 4: Get all facts (empty variable list)
	allFacts := wm.GetFactsByVariable([]string{})
	assert.Len(t, allFacts, 3, "Should return all facts when no variables specified")
	assert.Contains(t, allFacts, fact1, "Should contain fact1")
	assert.Contains(t, allFacts, fact2, "Should contain fact2")
	assert.Contains(t, allFacts, fact3, "Should contain fact3")
}

// TestWorkingMemory_GetTokensByVariable tests getting tokens by variable
func TestWorkingMemory_GetTokensByVariable(t *testing.T) {
	t.Log("🧪 TEST: GetTokensByVariable - Filtrage par variable")
	t.Log("=======================================================")

	wm := &WorkingMemory{
		NodeID: "test_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	// Créer des faits
	userFact := &Fact{ID: "u1", Type: "User", Fields: map[string]interface{}{"name": "Alice"}}
	orderFact := &Fact{ID: "o1", Type: "Order", Fields: map[string]interface{}{"id": 100}}
	productFact := &Fact{ID: "p1", Type: "Product", Fields: map[string]interface{}{"name": "Book"}}

	// Créer tokens avec différents bindings
	token1 := &Token{
		ID:       "token1",
		NodeID:   "test_node",
		Facts:    []*Fact{userFact},
		Bindings: NewBindingChainWith("user", userFact),
	}
	token2 := &Token{
		ID:       "token2",
		NodeID:   "test_node",
		Facts:    []*Fact{orderFact},
		Bindings: NewBindingChainWith("order", orderFact),
	}
	token3 := &Token{
		ID:     "token3",
		NodeID: "test_node",
		Facts:  []*Fact{userFact, orderFact},
		Bindings: NewBindingChainWith("user", userFact).
			Add("order", orderFact),
	}
	token4 := &Token{
		ID:       "token4",
		NodeID:   "test_node",
		Facts:    []*Fact{productFact},
		Bindings: NewBindingChainWith("product", productFact),
	}

	// Ajouter tokens à la mémoire
	wm.AddToken(token1)
	wm.AddToken(token2)
	wm.AddToken(token3)
	wm.AddToken(token4)

	// Test 1: Filtrer par une seule variable "user"
	t.Run("Filter by single variable 'user'", func(t *testing.T) {
		tokens := wm.GetTokensByVariable([]string{"user"})
		require.Len(t, tokens, 2, "Should return 2 tokens with 'user' binding")
		// Vérifier que les bons tokens sont retournés
		tokenIDs := make(map[string]bool)
		for _, token := range tokens {
			tokenIDs[token.ID] = true
		}
		assert.True(t, tokenIDs["token1"], "token1 should be included")
		assert.True(t, tokenIDs["token3"], "token3 should be included")
		t.Log("✅ Filtrage par 'user' correct")
	})

	// Test 2: Filtrer par une seule variable "order"
	t.Run("Filter by single variable 'order'", func(t *testing.T) {
		tokens := wm.GetTokensByVariable([]string{"order"})
		require.Len(t, tokens, 2, "Should return 2 tokens with 'order' binding")
		tokenIDs := make(map[string]bool)
		for _, token := range tokens {
			tokenIDs[token.ID] = true
		}
		assert.True(t, tokenIDs["token2"], "token2 should be included")
		assert.True(t, tokenIDs["token3"], "token3 should be included")
		t.Log("✅ Filtrage par 'order' correct")
	})

	// Test 3: Filtrer par multiples variables
	t.Run("Filter by multiple variables", func(t *testing.T) {
		tokens := wm.GetTokensByVariable([]string{"user", "product"})
		require.Len(t, tokens, 3, "Should return 3 tokens with 'user' OR 'product'")
		tokenIDs := make(map[string]bool)
		for _, token := range tokens {
			tokenIDs[token.ID] = true
		}
		assert.True(t, tokenIDs["token1"], "token1 should be included (has user)")
		assert.True(t, tokenIDs["token3"], "token3 should be included (has user)")
		assert.True(t, tokenIDs["token4"], "token4 should be included (has product)")
		assert.False(t, tokenIDs["token2"], "token2 should NOT be included (has only order)")
		t.Log("✅ Filtrage par multiples variables correct")
	})

	// Test 4: Variable inexistante
	t.Run("Filter by non-existent variable", func(t *testing.T) {
		tokens := wm.GetTokensByVariable([]string{"nonexistent"})
		assert.Empty(t, tokens, "Should return empty slice for non-existent variable")
		t.Log("✅ Filtrage par variable inexistante retourne vide")
	})

	// Test 5: Slice vide retourne tous les tokens
	t.Run("Empty variables slice returns all tokens", func(t *testing.T) {
		tokens := wm.GetTokensByVariable([]string{})
		assert.Len(t, tokens, 4, "Should return all 4 tokens")
		t.Log("✅ Slice vide retourne tous les tokens")
	})

	// Test 6: Nil slice retourne tous les tokens
	t.Run("Nil variables slice returns all tokens", func(t *testing.T) {
		tokens := wm.GetTokensByVariable(nil)
		assert.Len(t, tokens, 4, "Should return all 4 tokens")
		t.Log("✅ Slice nil retourne tous les tokens")
	})

	// Test 7: Token sans bindings
	t.Run("Token without bindings", func(t *testing.T) {
		tokenNoBindings := &Token{
			ID:       "token_no_bindings",
			NodeID:   "test_node",
			Facts:    []*Fact{userFact},
			Bindings: nil, // Pas de bindings
		}
		wm.AddToken(tokenNoBindings)

		tokens := wm.GetTokensByVariable([]string{"user"})
		// Ne devrait pas inclure le token sans bindings
		for _, token := range tokens {
			assert.NotEqual(t, "token_no_bindings", token.ID, "Token without bindings should not be included")
		}
		t.Log("✅ Token sans bindings correctement exclu")
	})

	t.Log("✅ Test réussi: GetTokensByVariable fonctionne correctement")
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
		Bindings:     NewBindingChain().Add("p", fact1).Add("o", fact2),
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
	assert.Equal(t, 2, clone.Bindings.Len(), "Should have 2 bindings")
	// Verify facts are copied (shallow copy)
	assert.Equal(t, fact1, clone.Facts[0], "First fact should match")
	assert.Equal(t, fact2, clone.Facts[1], "Second fact should match")
	// Verify bindings are copied (immutable, same reference is OK)
	assert.Equal(t, fact1, clone.Bindings.Get("p"), "Binding 'p' should match")
	assert.Equal(t, fact2, clone.Bindings.Get("o"), "Binding 'o' should match")
}

// TestToken_CloneEmpty tests cloning an empty token
func TestToken_CloneEmpty(t *testing.T) {
	original := &Token{
		ID:       "token1",
		Facts:    []*Fact{},
		NodeID:   "test_node",
		Bindings: NewBindingChain(),
	}
	clone := original.Clone()
	require.NotNil(t, clone, "Clone should not be nil")
	assert.Equal(t, original.ID, clone.ID)
	assert.Equal(t, original.NodeID, clone.NodeID)
	assert.NotNil(t, clone.Facts, "Facts should be initialized")
	assert.Equal(t, 0, clone.Bindings.Len(), "Bindings should be empty")
	assert.Empty(t, clone.Facts, "Facts should be empty")
}

// TestToken_CloneIndependence tests that cloned token is independent
func TestToken_CloneIndependence(t *testing.T) {
	fact := &Fact{ID: "fact1", Type: "Person", Fields: map[string]interface{}{"name": "Alice"}}
	original := &Token{
		ID:       "token1",
		Facts:    []*Fact{fact},
		NodeID:   "test_node",
		Bindings: NewBindingChainWith("p", fact),
	}
	clone := original.Clone()
	// Note: BindingChain is immutable, so we create a NEW chain
	// This test needs to be updated to reflect the immutability
	newFact := &Fact{ID: "fact2", Type: "Order"}
	clone.Facts = append(clone.Facts, newFact)
	clone.Bindings = clone.Bindings.Add("o", newFact)

	// Verify original is unchanged
	assert.Len(t, original.Facts, 1, "Original should still have 1 fact")
	assert.Equal(t, 1, original.Bindings.Len(), "Original should still have 1 binding")
	assert.Nil(t, original.Bindings.Get("o"), "Original should not have 'o' binding")

	// Verify clone has new values
	assert.Len(t, clone.Facts, 2, "Clone should have 2 facts")
	assert.Equal(t, 2, clone.Bindings.Len(), "Clone should have 2 bindings")
	assert.NotNil(t, clone.Bindings.Get("o"), "Clone should have 'o' binding")
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

// TestFact_IDHandling teste la manipulation des IDs de faits
func TestFact_IDHandling(t *testing.T) {
	t.Log("🧪 TEST: Fact ID Handling - Manipulation des IDs")
	t.Log("==================================================")

	tests := []struct {
		name     string
		fact     *Fact
		wantID   string
		wantType string
	}{
		{
			name: "fait avec PK simple",
			fact: &Fact{
				ID:   "Person~Alice",
				Type: "Person",
				Fields: map[string]interface{}{
					"nom": "Alice",
					"age": 30,
				},
			},
			wantID:   "Person~Alice",
			wantType: "Person",
		},
		{
			name: "fait avec PK composite",
			fact: &Fact{
				ID:   "Person~Alice_Dupont",
				Type: "Person",
				Fields: map[string]interface{}{
					"prenom": "Alice",
					"nom":    "Dupont",
					"age":    30,
				},
			},
			wantID:   "Person~Alice_Dupont",
			wantType: "Person",
		},
		{
			name: "fait avec hash",
			fact: &Fact{
				ID:   "Event~a1b2c3d4e5f6g7h8",
				Type: "Event",
				Fields: map[string]interface{}{
					"timestamp": 1234567890,
					"message":   "test",
				},
			},
			wantID:   "Event~a1b2c3d4e5f6g7h8",
			wantType: "Event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fact.ID != tt.wantID {
				t.Errorf("❌ ID = %v, attendu %v", tt.fact.ID, tt.wantID)
			}
			if tt.fact.Type != tt.wantType {
				t.Errorf("❌ Type = %v, attendu %v", tt.fact.Type, tt.wantType)
			}
			t.Log("✅ Test réussi")
		})
	}

	t.Log("✅ Test complet: Fact ID Handling")
}
