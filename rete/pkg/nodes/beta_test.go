package nodes

import (
	"fmt"
	"sync"
	"testing"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// Tests pour NewBetaMemory
func TestNewBetaMemory(t *testing.T) {
	memory := NewBetaMemory()

	if memory == nil {
		t.Fatal("NewBetaMemory returned nil")
	}

	tokens, facts := memory.Size()
	if tokens != 0 {
		t.Errorf("Expected 0 tokens, got %d", tokens)
	}
	if facts != 0 {
		t.Errorf("Expected 0 facts, got %d", facts)
	}
}

// Tests pour StoreToken
func TestBetaMemory_StoreToken(t *testing.T) {
	memory := NewBetaMemory()

	token := &domain.Token{
		ID: "token-1",
		Facts: []*domain.Fact{
			{ID: "fact-1", Type: "Customer"},
		},
	}

	memory.StoreToken(token)

	tokens, _ := memory.Size()
	if tokens != 1 {
		t.Errorf("Expected 1 token, got %d", tokens)
	}

	retrieved := memory.GetTokens()
	if len(retrieved) != 1 {
		t.Fatalf("Expected 1 token in GetTokens, got %d", len(retrieved))
	}
	if retrieved[0].ID != "token-1" {
		t.Errorf("Expected token ID 'token-1', got '%s'", retrieved[0].ID)
	}
}

// Tests pour StoreToken - multiple tokens
func TestBetaMemory_StoreToken_Multiple(t *testing.T) {
	memory := NewBetaMemory()

	for i := 0; i < 5; i++ {
		token := &domain.Token{
			ID: fmt.Sprintf("token-%d", i),
			Facts: []*domain.Fact{
				{ID: fmt.Sprintf("fact-%d", i), Type: "Customer"},
			},
		}
		memory.StoreToken(token)
	}

	tokens, _ := memory.Size()
	if tokens != 5 {
		t.Errorf("Expected 5 tokens, got %d", tokens)
	}

	retrieved := memory.GetTokens()
	if len(retrieved) != 5 {
		t.Errorf("Expected 5 tokens in GetTokens, got %d", len(retrieved))
	}
}

// Tests pour StoreToken - overwrite existing
func TestBetaMemory_StoreToken_Overwrite(t *testing.T) {
	memory := NewBetaMemory()

	token1 := &domain.Token{
		ID: "token-1",
		Facts: []*domain.Fact{
			{ID: "fact-1", Type: "Customer"},
		},
	}
	memory.StoreToken(token1)

	// Stocker un autre token avec le même ID
	token2 := &domain.Token{
		ID: "token-1",
		Facts: []*domain.Fact{
			{ID: "fact-2", Type: "Order"},
		},
	}
	memory.StoreToken(token2)

	tokens, _ := memory.Size()
	if tokens != 1 {
		t.Errorf("Expected 1 token (overwritten), got %d", tokens)
	}

	retrieved := memory.GetTokens()
	if len(retrieved) != 1 {
		t.Fatalf("Expected 1 token in GetTokens, got %d", len(retrieved))
	}
	// Vérifier que c'est le nouveau token
	if len(retrieved[0].Facts) != 1 || retrieved[0].Facts[0].ID != "fact-2" {
		t.Error("Expected token to be overwritten with new data")
	}
}

// Tests pour RemoveToken
func TestBetaMemory_RemoveToken(t *testing.T) {
	memory := NewBetaMemory()

	token := &domain.Token{
		ID: "token-1",
		Facts: []*domain.Fact{
			{ID: "fact-1", Type: "Customer"},
		},
	}
	memory.StoreToken(token)

	// Retirer le token
	removed := memory.RemoveToken("token-1")
	if !removed {
		t.Error("Expected RemoveToken to return true")
	}

	tokens, _ := memory.Size()
	if tokens != 0 {
		t.Errorf("Expected 0 tokens after removal, got %d", tokens)
	}
}

// Tests pour RemoveToken - non-existent
func TestBetaMemory_RemoveToken_NonExistent(t *testing.T) {
	memory := NewBetaMemory()

	removed := memory.RemoveToken("non-existent")
	if removed {
		t.Error("Expected RemoveToken to return false for non-existent token")
	}
}

// Tests pour GetTokens
func TestBetaMemory_GetTokens(t *testing.T) {
	t.Run("empty memory", func(t *testing.T) {
		memory := NewBetaMemory()

		tokens := memory.GetTokens()
		if tokens == nil {
			t.Fatal("GetTokens returned nil")
		}
		if len(tokens) != 0 {
			t.Errorf("Expected 0 tokens, got %d", len(tokens))
		}
	})

	t.Run("with tokens", func(t *testing.T) {
		memory := NewBetaMemory()

		for i := 0; i < 3; i++ {
			token := &domain.Token{
				ID: fmt.Sprintf("token-%d", i),
				Facts: []*domain.Fact{
					{ID: fmt.Sprintf("fact-%d", i), Type: "Customer"},
				},
			}
			memory.StoreToken(token)
		}

		tokens := memory.GetTokens()
		if len(tokens) != 3 {
			t.Errorf("Expected 3 tokens, got %d", len(tokens))
		}

		// Vérifier que tous les tokens sont présents
		tokenIDs := make(map[string]bool)
		for _, token := range tokens {
			tokenIDs[token.ID] = true
		}

		for i := 0; i < 3; i++ {
			expectedID := fmt.Sprintf("token-%d", i)
			if !tokenIDs[expectedID] {
				t.Errorf("Expected token '%s' in results", expectedID)
			}
		}
	})
}

// Tests pour StoreFact
func TestBetaMemory_StoreFact(t *testing.T) {
	memory := NewBetaMemory()

	fact := &domain.Fact{
		ID:   "fact-1",
		Type: "Customer",
		Fields: map[string]interface{}{
			"name": "John",
			"age":  30,
		},
	}

	memory.StoreFact(fact)

	_, facts := memory.Size()
	if facts != 1 {
		t.Errorf("Expected 1 fact, got %d", facts)
	}

	retrieved := memory.GetFacts()
	if len(retrieved) != 1 {
		t.Fatalf("Expected 1 fact in GetFacts, got %d", len(retrieved))
	}
	if retrieved[0].ID != "fact-1" {
		t.Errorf("Expected fact ID 'fact-1', got '%s'", retrieved[0].ID)
	}
}

// Tests pour StoreFact - multiple facts
func TestBetaMemory_StoreFact_Multiple(t *testing.T) {
	memory := NewBetaMemory()

	for i := 0; i < 5; i++ {
		fact := &domain.Fact{
			ID:   fmt.Sprintf("fact-%d", i),
			Type: "Customer",
			Fields: map[string]interface{}{
				"id": i,
			},
		}
		memory.StoreFact(fact)
	}

	_, facts := memory.Size()
	if facts != 5 {
		t.Errorf("Expected 5 facts, got %d", facts)
	}
}

// Tests pour RemoveFact
func TestBetaMemory_RemoveFact(t *testing.T) {
	memory := NewBetaMemory()

	fact := &domain.Fact{
		ID:   "fact-1",
		Type: "Customer",
	}
	memory.StoreFact(fact)

	removed := memory.RemoveFact("fact-1")
	if !removed {
		t.Error("Expected RemoveFact to return true")
	}

	_, facts := memory.Size()
	if facts != 0 {
		t.Errorf("Expected 0 facts after removal, got %d", facts)
	}
}

// Tests pour RemoveFact - non-existent
func TestBetaMemory_RemoveFact_NonExistent(t *testing.T) {
	memory := NewBetaMemory()

	removed := memory.RemoveFact("non-existent")
	if removed {
		t.Error("Expected RemoveFact to return false for non-existent fact")
	}
}

// Tests pour GetFacts
func TestBetaMemory_GetFacts(t *testing.T) {
	t.Run("empty memory", func(t *testing.T) {
		memory := NewBetaMemory()

		facts := memory.GetFacts()
		if facts == nil {
			t.Fatal("GetFacts returned nil")
		}
		if len(facts) != 0 {
			t.Errorf("Expected 0 facts, got %d", len(facts))
		}
	})

	t.Run("with facts", func(t *testing.T) {
		memory := NewBetaMemory()

		for i := 0; i < 4; i++ {
			fact := &domain.Fact{
				ID:   fmt.Sprintf("fact-%d", i),
				Type: "Customer",
			}
			memory.StoreFact(fact)
		}

		facts := memory.GetFacts()
		if len(facts) != 4 {
			t.Errorf("Expected 4 facts, got %d", len(facts))
		}
	})
}

// Tests pour Clear
func TestBetaMemory_Clear(t *testing.T) {
	memory := NewBetaMemory()

	// Ajouter des tokens et des faits
	for i := 0; i < 3; i++ {
		token := &domain.Token{
			ID: fmt.Sprintf("token-%d", i),
		}
		memory.StoreToken(token)

		fact := &domain.Fact{
			ID:   fmt.Sprintf("fact-%d", i),
			Type: "Customer",
		}
		memory.StoreFact(fact)
	}

	tokens, facts := memory.Size()
	if tokens != 3 || facts != 3 {
		t.Errorf("Expected 3 tokens and 3 facts before clear, got %d tokens and %d facts", tokens, facts)
	}

	// Clear
	memory.Clear()

	tokens, facts = memory.Size()
	if tokens != 0 {
		t.Errorf("Expected 0 tokens after clear, got %d", tokens)
	}
	if facts != 0 {
		t.Errorf("Expected 0 facts after clear, got %d", facts)
	}

	// Vérifier que GetTokens et GetFacts retournent des slices vides
	retrievedTokens := memory.GetTokens()
	if len(retrievedTokens) != 0 {
		t.Errorf("Expected 0 tokens in GetTokens after clear, got %d", len(retrievedTokens))
	}

	retrievedFacts := memory.GetFacts()
	if len(retrievedFacts) != 0 {
		t.Errorf("Expected 0 facts in GetFacts after clear, got %d", len(retrievedFacts))
	}
}

// Tests pour Size
func TestBetaMemory_Size(t *testing.T) {
	memory := NewBetaMemory()

	// Vide
	tokens, facts := memory.Size()
	if tokens != 0 || facts != 0 {
		t.Errorf("Expected 0/0, got %d/%d", tokens, facts)
	}

	// Ajouter des tokens
	for i := 0; i < 3; i++ {
		token := &domain.Token{ID: fmt.Sprintf("token-%d", i)}
		memory.StoreToken(token)
	}

	tokens, facts = memory.Size()
	if tokens != 3 || facts != 0 {
		t.Errorf("Expected 3/0, got %d/%d", tokens, facts)
	}

	// Ajouter des faits
	for i := 0; i < 2; i++ {
		fact := &domain.Fact{ID: fmt.Sprintf("fact-%d", i), Type: "Customer"}
		memory.StoreFact(fact)
	}

	tokens, facts = memory.Size()
	if tokens != 3 || facts != 2 {
		t.Errorf("Expected 3/2, got %d/%d", tokens, facts)
	}

	// Retirer un token
	memory.RemoveToken("token-1")

	tokens, facts = memory.Size()
	if tokens != 2 || facts != 2 {
		t.Errorf("Expected 2/2, got %d/%d", tokens, facts)
	}

	// Retirer un fait
	memory.RemoveFact("fact-0")

	tokens, facts = memory.Size()
	if tokens != 2 || facts != 1 {
		t.Errorf("Expected 2/1, got %d/%d", tokens, facts)
	}
}

// Tests de concurrence pour BetaMemory
func TestBetaMemory_Concurrency_StoreToken(t *testing.T) {
	memory := NewBetaMemory()

	const goroutines = 10
	const tokensPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(index int) {
			defer wg.Done()
			for j := 0; j < tokensPerGoroutine; j++ {
				token := &domain.Token{
					ID: fmt.Sprintf("token-%d-%d", index, j),
				}
				memory.StoreToken(token)
			}
		}(i)
	}

	wg.Wait()

	tokens, _ := memory.Size()
	expected := goroutines * tokensPerGoroutine
	if tokens != expected {
		t.Errorf("Expected %d tokens, got %d", expected, tokens)
	}
}

func TestBetaMemory_Concurrency_StoreFact(t *testing.T) {
	memory := NewBetaMemory()

	const goroutines = 10
	const factsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(index int) {
			defer wg.Done()
			for j := 0; j < factsPerGoroutine; j++ {
				fact := &domain.Fact{
					ID:   fmt.Sprintf("fact-%d-%d", index, j),
					Type: "Customer",
				}
				memory.StoreFact(fact)
			}
		}(i)
	}

	wg.Wait()

	_, facts := memory.Size()
	expected := goroutines * factsPerGoroutine
	if facts != expected {
		t.Errorf("Expected %d facts, got %d", expected, facts)
	}
}

func TestBetaMemory_Concurrency_Mixed(t *testing.T) {
	memory := NewBetaMemory()

	const operations = 1000

	var wg sync.WaitGroup
	wg.Add(4)

	// Writer de tokens
	go func() {
		defer wg.Done()
		for i := 0; i < operations; i++ {
			token := &domain.Token{ID: fmt.Sprintf("token-%d", i)}
			memory.StoreToken(token)
		}
	}()

	// Writer de faits
	go func() {
		defer wg.Done()
		for i := 0; i < operations; i++ {
			fact := &domain.Fact{ID: fmt.Sprintf("fact-%d", i), Type: "Customer"}
			memory.StoreFact(fact)
		}
	}()

	// Reader de tokens
	go func() {
		defer wg.Done()
		for i := 0; i < operations; i++ {
			_ = memory.GetTokens()
		}
	}()

	// Reader de faits
	go func() {
		defer wg.Done()
		for i := 0; i < operations; i++ {
			_ = memory.GetFacts()
		}
	}()

	wg.Wait()

	tokens, facts := memory.Size()
	if tokens != operations {
		t.Errorf("Expected %d tokens, got %d", operations, tokens)
	}
	if facts != operations {
		t.Errorf("Expected %d facts, got %d", operations, facts)
	}
}

// Tests pour NewBaseBetaNode
func TestNewBaseBetaNode(t *testing.T) {
	logger := newMockLogger()
	node := NewBaseBetaNode("beta-node-1", "join", logger)

	if node == nil {
		t.Fatal("NewBaseBetaNode returned nil")
	}

	if node.ID() != "beta-node-1" {
		t.Errorf("Expected ID 'beta-node-1', got '%s'", node.ID())
	}

	if node.Type() != "join" {
		t.Errorf("Expected Type 'join', got '%s'", node.Type())
	}

	if node.betaMemory == nil {
		t.Error("Expected non-nil betaMemory")
	}
}

// Tests cas limites
func TestBetaMemory_EdgeCases(t *testing.T) {
	t.Run("nil token", func(t *testing.T) {
		memory := NewBetaMemory()

		// Le stockage d'un token nil devrait paniquer
		defer func() {
			if r := recover(); r == nil {
				t.Error("StoreToken should panic with nil token")
			}
		}()

		memory.StoreToken(nil)
	})

	t.Run("nil fact", func(t *testing.T) {
		memory := NewBetaMemory()

		// Le stockage d'un fait nil devrait paniquer
		defer func() {
			if r := recover(); r == nil {
				t.Error("StoreFact should panic with nil fact")
			}
		}()

		memory.StoreFact(nil)
	})

	t.Run("empty token ID", func(t *testing.T) {
		memory := NewBetaMemory()

		token := &domain.Token{ID: ""}
		memory.StoreToken(token)

		retrieved := memory.GetTokens()
		if len(retrieved) != 1 {
			t.Errorf("Expected 1 token with empty ID, got %d", len(retrieved))
		}
	})

	t.Run("clear empty memory", func(t *testing.T) {
		memory := NewBetaMemory()

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Clear panicked on empty memory: %v", r)
			}
		}()

		memory.Clear()

		tokens, facts := memory.Size()
		if tokens != 0 || facts != 0 {
			t.Errorf("Expected 0/0 after clearing empty memory, got %d/%d", tokens, facts)
		}
	})
}

// Benchmarks
func BenchmarkBetaMemory_StoreToken(b *testing.B) {
	memory := NewBetaMemory()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		token := &domain.Token{ID: fmt.Sprintf("token-%d", i)}
		memory.StoreToken(token)
	}
}

func BenchmarkBetaMemory_GetTokens(b *testing.B) {
	memory := NewBetaMemory()

	// Pré-remplir avec des tokens
	for i := 0; i < 1000; i++ {
		token := &domain.Token{ID: fmt.Sprintf("token-%d", i)}
		memory.StoreToken(token)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = memory.GetTokens()
	}
}

func BenchmarkBetaMemory_StoreFact(b *testing.B) {
	memory := NewBetaMemory()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fact := &domain.Fact{ID: fmt.Sprintf("fact-%d", i), Type: "Customer"}
		memory.StoreFact(fact)
	}
}

func BenchmarkBetaMemory_Size(b *testing.B) {
	memory := NewBetaMemory()

	// Pré-remplir
	for i := 0; i < 100; i++ {
		token := &domain.Token{ID: fmt.Sprintf("token-%d", i)}
		memory.StoreToken(token)
		fact := &domain.Fact{ID: fmt.Sprintf("fact-%d", i), Type: "Customer"}
		memory.StoreFact(fact)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = memory.Size()
	}
}
