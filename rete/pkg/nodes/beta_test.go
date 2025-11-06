package nodes

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// mockLogger pour les tests
type mockLogger struct{}

func (m *mockLogger) Debug(msg string, fields map[string]interface{})            {}
func (m *mockLogger) Info(msg string, fields map[string]interface{})             {}
func (m *mockLogger) Warn(msg string, fields map[string]interface{})             {}
func (m *mockLogger) Error(msg string, err error, fields map[string]interface{}) {}

// Helper pour créer des faits de test
func createTestFact(id, factType string, fields map[string]interface{}) *domain.Fact {
	return domain.NewFact(id, factType, fields)
}

// Helper pour créer des tokens de test
func createTestToken(id, nodeID string, facts []*domain.Fact) *domain.Token {
	return domain.NewToken(id, nodeID, facts)
}

// Helper pour extraire les IDs des faits
func getFactIDs(facts []*domain.Fact) []string {
	ids := make([]string, len(facts))
	for i, fact := range facts {
		ids[i] = fact.ID
	}
	return ids
}

func TestBetaMemoryImpl(t *testing.T) {
	t.Run("NewBetaMemory", func(t *testing.T) {
		memory := NewBetaMemory()
		if memory == nil {
			t.Fatal("NewBetaMemory should not return nil")
		}

		tokens, facts := memory.Size()
		if tokens != 0 || facts != 0 {
			t.Errorf("New memory should be empty, got %d tokens and %d facts", tokens, facts)
		}
	})

	t.Run("StoreToken", func(t *testing.T) {
		memory := NewBetaMemory()
		fact := createTestFact("f1", "Person", map[string]interface{}{"name": "John"})
		token := createTestToken("t1", "node1", []*domain.Fact{fact})

		memory.StoreToken(token)

		tokens := memory.GetTokens()
		if len(tokens) != 1 {
			t.Errorf("Expected 1 token, got %d", len(tokens))
		}
		if tokens[0].ID != "t1" {
			t.Errorf("Expected token ID 't1', got '%s'", tokens[0].ID)
		}
	})

	t.Run("RemoveToken", func(t *testing.T) {
		memory := NewBetaMemory()
		fact := createTestFact("f1", "Person", map[string]interface{}{"name": "John"})
		token := createTestToken("t1", "node1", []*domain.Fact{fact})

		memory.StoreToken(token)

		removed := memory.RemoveToken("t1")
		if !removed {
			t.Error("Token should have been removed")
		}

		tokens := memory.GetTokens()
		if len(tokens) != 0 {
			t.Errorf("Expected 0 tokens after removal, got %d", len(tokens))
		}
	})

	t.Run("StoreFact", func(t *testing.T) {
		memory := NewBetaMemory()
		fact := createTestFact("f1", "Person", map[string]interface{}{"name": "John"})

		memory.StoreFact(fact)

		facts := memory.GetFacts()
		if len(facts) != 1 {
			t.Errorf("Expected 1 fact, got %d", len(facts))
		}
		if facts[0].ID != "f1" {
			t.Errorf("Expected fact ID 'f1', got '%s'", facts[0].ID)
		}
	})

	t.Run("Clear", func(t *testing.T) {
		memory := NewBetaMemory()
		fact := createTestFact("f1", "Person", map[string]interface{}{"name": "John"})
		token := createTestToken("t1", "node1", []*domain.Fact{fact})

		memory.StoreToken(token)
		memory.StoreFact(fact)

		memory.Clear()

		tokens, facts := memory.Size()
		if tokens != 0 || facts != 0 {
			t.Errorf("Memory should be empty after clear, got %d tokens and %d facts", tokens, facts)
		}
	})
}

func TestBaseBetaNode(t *testing.T) {
	logger := &mockLogger{}

	t.Run("NewBaseBetaNode", func(t *testing.T) {
		node := NewBaseBetaNode("beta1", "BetaNode", logger)
		if node == nil {
			t.Fatal("NewBaseBetaNode should not return nil")
		}

		if node.ID() != "beta1" {
			t.Errorf("Expected ID 'beta1', got '%s'", node.ID())
		}
		if node.Type() != "BetaNode" {
			t.Errorf("Expected type 'BetaNode', got '%s'", node.Type())
		}
	})

	t.Run("ProcessLeftToken", func(t *testing.T) {
		node := NewBaseBetaNode("beta1", "BetaNode", logger)
		fact := createTestFact("f1", "Person", map[string]interface{}{"name": "John"})
		token := createTestToken("t1", "node1", []*domain.Fact{fact})

		err := node.ProcessLeftToken(token)
		if err != nil {
			t.Errorf("ProcessLeftToken should not return error: %v", err)
		}

		leftMemory := node.GetLeftMemory()
		if len(leftMemory) != 1 {
			t.Errorf("Expected 1 token in left memory, got %d", len(leftMemory))
		}
	})

	t.Run("ProcessRightFact", func(t *testing.T) {
		node := NewBaseBetaNode("beta1", "BetaNode", logger)
		fact := createTestFact("f1", "Person", map[string]interface{}{"name": "John"})

		err := node.ProcessRightFact(fact)
		if err != nil {
			t.Errorf("ProcessRightFact should not return error: %v", err)
		}

		rightMemory := node.GetRightMemory()
		if len(rightMemory) != 1 {
			t.Errorf("Expected 1 fact in right memory, got %d", len(rightMemory))
		}
	})

	t.Run("JoinOperation", func(t *testing.T) {
		node := NewBaseBetaNode("beta1", "BetaNode", logger)

		// Créer un token avec un fait
		leftFact := createTestFact("f1", "Person", map[string]interface{}{"name": "John"})
		token := createTestToken("t1", "node1", []*domain.Fact{leftFact})

		// Créer un fait pour le côté droit
		rightFact := createTestFact("f2", "Address", map[string]interface{}{"city": "Paris"})

		// Mock child pour capturer les résultats
		capturedTokens := []*domain.Token{}
		mockChild := &mockBetaNode{
			processLeftToken: func(t *domain.Token) error {
				capturedTokens = append(capturedTokens, t)
				return nil
			},
		}
		node.AddChild(mockChild)

		// Traiter le token gauche puis le fait droit
		err := node.ProcessLeftToken(token)
		if err != nil {
			t.Errorf("ProcessLeftToken failed: %v", err)
		}

		err = node.ProcessRightFact(rightFact)
		if err != nil {
			t.Errorf("ProcessRightFact failed: %v", err)
		}

		// Vérifier qu'une jointure a été produite
		if len(capturedTokens) != 1 {
			t.Errorf("Expected 1 joined token, got %d", len(capturedTokens))
		}

		if len(capturedTokens) > 0 {
			joinedToken := capturedTokens[0]
			if len(joinedToken.Facts) != 2 {
				t.Errorf("Expected joined token to have 2 facts, got %d", len(joinedToken.Facts))
			}
		}
	})

	t.Run("ClearMemory", func(t *testing.T) {
		node := NewBaseBetaNode("beta1", "BetaNode", logger)

		fact := createTestFact("f1", "Person", map[string]interface{}{"name": "John"})
		token := createTestToken("t1", "node1", []*domain.Fact{fact})

		node.ProcessLeftToken(token)
		node.ProcessRightFact(fact)

		node.ClearMemory()

		if len(node.GetLeftMemory()) != 0 {
			t.Error("Left memory should be empty after clear")
		}
		if len(node.GetRightMemory()) != 0 {
			t.Error("Right memory should be empty after clear")
		}
	})
}

func TestJoinNodeImpl(t *testing.T) {
	logger := &mockLogger{}

	t.Run("NewJoinNode", func(t *testing.T) {
		node := NewJoinNode("join1", logger)
		if node == nil {
			t.Fatal("NewJoinNode should not return nil")
		}

		if node.ID() != "join1" {
			t.Errorf("Expected ID 'join1', got '%s'", node.ID())
		}
		if node.Type() != "JoinNode" {
			t.Errorf("Expected type 'JoinNode', got '%s'", node.Type())
		}
	})

	t.Run("SetAndGetJoinConditions", func(t *testing.T) {
		node := NewJoinNode("join1", logger)

		condition := domain.NewBasicJoinCondition("name", "name", "==")
		conditions := []domain.JoinCondition{condition}

		node.SetJoinConditions(conditions)

		retrievedConditions := node.GetJoinConditions()
		if len(retrievedConditions) != 1 {
			t.Errorf("Expected 1 condition, got %d", len(retrievedConditions))
		}
	})

	t.Run("EvaluateJoin_NoConditions", func(t *testing.T) {
		node := NewJoinNode("join1", logger)

		leftFact := createTestFact("f1", "Person", map[string]interface{}{"name": "John"})
		token := createTestToken("t1", "node1", []*domain.Fact{leftFact})
		rightFact := createTestFact("f2", "Person", map[string]interface{}{"name": "Jane"})

		result := node.EvaluateJoin(token, rightFact)
		if !result {
			t.Error("Join should succeed when no conditions are set")
		}
	})

	t.Run("EvaluateJoin_WithConditions_Success", func(t *testing.T) {
		node := NewJoinNode("join1", logger)

		condition := domain.NewBasicJoinCondition("name", "name", "==")
		node.SetJoinConditions([]domain.JoinCondition{condition})

		leftFact := createTestFact("f1", "Person", map[string]interface{}{"name": "John"})
		token := createTestToken("t1", "node1", []*domain.Fact{leftFact})
		rightFact := createTestFact("f2", "Person", map[string]interface{}{"name": "John"})

		result := node.EvaluateJoin(token, rightFact)
		if !result {
			t.Error("Join should succeed when conditions are met")
		}
	})

	t.Run("EvaluateJoin_WithConditions_Failure", func(t *testing.T) {
		node := NewJoinNode("join1", logger)

		condition := domain.NewBasicJoinCondition("name", "name", "==")
		node.SetJoinConditions([]domain.JoinCondition{condition})

		leftFact := createTestFact("f1", "Person", map[string]interface{}{"name": "John"})
		token := createTestToken("t1", "node1", []*domain.Fact{leftFact})
		rightFact := createTestFact("f2", "Person", map[string]interface{}{"name": "Jane"})

		result := node.EvaluateJoin(token, rightFact)
		if result {
			t.Error("Join should fail when conditions are not met")
		}
	})

	t.Run("JoinWithConditions_Integration", func(t *testing.T) {
		node := NewJoinNode("join1", logger)

		// Condition: les âges doivent être égaux
		condition := domain.NewBasicJoinCondition("age", "age", "==")
		node.SetJoinConditions([]domain.JoinCondition{condition})

		// Mock child pour capturer les résultats
		capturedTokens := []*domain.Token{}
		mockChild := &mockBetaNode{
			processLeftToken: func(t *domain.Token) error {
				capturedTokens = append(capturedTokens, t)
				return nil
			},
		}
		node.AddChild(mockChild)

		// Token avec une personne de 25 ans
		leftFact := createTestFact("f1", "Person", map[string]interface{}{"name": "John", "age": 25})
		token := createTestToken("t1", "node1", []*domain.Fact{leftFact})

		// Fait avec une personne de 30 ans (ne devrait pas matcher)
		rightFact1 := createTestFact("f2", "Person", map[string]interface{}{"name": "Jane", "age": 30})

		// Fait avec une personne de 25 ans (devrait matcher)
		rightFact2 := createTestFact("f3", "Person", map[string]interface{}{"name": "Bob", "age": 25})

		// Traitement - ordre différent pour tester correctement
		node.ProcessRightFact(rightFact1) // Stocker d'abord le fait qui ne match pas
		node.ProcessRightFact(rightFact2) // Stocker le fait qui match
		node.ProcessLeftToken(token)      // Maintenant traiter le token - devrait créer une seule jointure

		// Vérification
		if len(capturedTokens) != 1 {
			t.Errorf("Expected exactly 1 successful join, got %d", len(capturedTokens))
			for i, token := range capturedTokens {
				t.Logf("Token %d: ID=%s, FactIDs=%v", i, token.ID, getFactIDs(token.Facts))
			}
		}

		if len(capturedTokens) > 0 {
			joinedToken := capturedTokens[0]
			if len(joinedToken.Facts) != 2 {
				t.Errorf("Expected joined token to have 2 facts, got %d", len(joinedToken.Facts))
			}

			// Vérifier que c'est bien le bon fait qui a été joint (celui avec age=25)
			found := false
			for _, fact := range joinedToken.Facts {
				if fact.ID == "f3" {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected f3 to be in joined token, got facts: %v", getFactIDs(joinedToken.Facts))
			}
		}
	})
}

func TestBasicJoinCondition(t *testing.T) {
	t.Run("NewBasicJoinCondition", func(t *testing.T) {
		condition := domain.NewBasicJoinCondition("left", "right", "==")

		if condition.GetLeftField() != "left" {
			t.Errorf("Expected left field 'left', got '%s'", condition.GetLeftField())
		}
		if condition.GetRightField() != "right" {
			t.Errorf("Expected right field 'right', got '%s'", condition.GetRightField())
		}
		if condition.GetOperator() != "==" {
			t.Errorf("Expected operator '==', got '%s'", condition.GetOperator())
		}
	})

	t.Run("Evaluate_Equality", func(t *testing.T) {
		condition := domain.NewBasicJoinCondition("name", "name", "==")

		leftFact := createTestFact("f1", "Person", map[string]interface{}{"name": "John"})
		token := createTestToken("t1", "node1", []*domain.Fact{leftFact})
		rightFact := createTestFact("f2", "Person", map[string]interface{}{"name": "John"})

		result := condition.Evaluate(token, rightFact)
		if !result {
			t.Error("Equality condition should return true for equal values")
		}
	})

	t.Run("Evaluate_Inequality", func(t *testing.T) {
		condition := domain.NewBasicJoinCondition("age", "age", "!=")

		leftFact := createTestFact("f1", "Person", map[string]interface{}{"age": 25})
		token := createTestToken("t1", "node1", []*domain.Fact{leftFact})
		rightFact := createTestFact("f2", "Person", map[string]interface{}{"age": 30})

		result := condition.Evaluate(token, rightFact)
		if !result {
			t.Error("Inequality condition should return true for different values")
		}
	})

	t.Run("Evaluate_LessThan", func(t *testing.T) {
		condition := domain.NewBasicJoinCondition("age", "age", "<")

		leftFact := createTestFact("f1", "Person", map[string]interface{}{"age": 25})
		token := createTestToken("t1", "node1", []*domain.Fact{leftFact})
		rightFact := createTestFact("f2", "Person", map[string]interface{}{"age": 30})

		result := condition.Evaluate(token, rightFact)
		if !result {
			t.Error("Less than condition should return true when left < right")
		}
	})

	t.Run("Evaluate_MissingField", func(t *testing.T) {
		condition := domain.NewBasicJoinCondition("nonexistent", "name", "==")

		leftFact := createTestFact("f1", "Person", map[string]interface{}{"name": "John"})
		token := createTestToken("t1", "node1", []*domain.Fact{leftFact})
		rightFact := createTestFact("f2", "Person", map[string]interface{}{"name": "John"})

		result := condition.Evaluate(token, rightFact)
		if result {
			t.Error("Condition should return false when field doesn't exist")
		}
	})
}

func TestEdgeCases(t *testing.T) {
	logger := &mockLogger{}

	t.Run("ProcessFact_DelegatesToProcessRightFact", func(t *testing.T) {
		node := NewBaseBetaNode("beta1", "BetaNode", logger)
		fact := createTestFact("f1", "Person", map[string]interface{}{"name": "John"})

		err := node.ProcessFact(fact)
		if err != nil {
			t.Errorf("ProcessFact should delegate to ProcessRightFact without error: %v", err)
		}

		rightMemory := node.GetRightMemory()
		if len(rightMemory) != 1 {
			t.Errorf("ProcessFact should store fact in right memory via delegation")
		}
	})

	t.Run("JoinCondition_AllOperators", func(t *testing.T) {
		operators := []struct {
			op       string
			left     interface{}
			right    interface{}
			expected bool
		}{
			{"==", 25, 25, true},
			{"=", 25, 25, true},
			{"!=", 25, 30, true},
			{"!=", 25, 25, false},
			{"<", 25, 30, true},
			{"<", 30, 25, false},
			{"<=", 25, 25, true},
			{"<=", 25, 30, true},
			{"<=", 30, 25, false},
			{">", 30, 25, true},
			{">", 25, 30, false},
			{">=", 25, 25, true},
			{">=", 30, 25, true},
			{">=", 25, 30, false},
			{"unknown", 25, 25, false},
		}

		for _, op := range operators {
			condition := domain.NewBasicJoinCondition("age", "age", op.op)
			leftFact := createTestFact("f1", "Person", map[string]interface{}{"age": op.left})
			token := createTestToken("t1", "node1", []*domain.Fact{leftFact})
			rightFact := createTestFact("f2", "Person", map[string]interface{}{"age": op.right})

			result := condition.Evaluate(token, rightFact)
			if result != op.expected {
				t.Errorf("Operator %s with %v %s %v: expected %v, got %v",
					op.op, op.left, op.op, op.right, op.expected, result)
			}
		}
	})

	t.Run("JoinCondition_TypeHandling", func(t *testing.T) {
		// Test with strings
		condition := domain.NewBasicJoinCondition("name", "name", "<")
		leftFact := createTestFact("f1", "Person", map[string]interface{}{"name": "Alice"})
		token := createTestToken("t1", "node1", []*domain.Fact{leftFact})
		rightFact := createTestFact("f2", "Person", map[string]interface{}{"name": "Bob"})

		result := condition.Evaluate(token, rightFact)
		if !result {
			t.Error("String comparison 'Alice' < 'Bob' should be true")
		}

		// Test with floats
		condition = domain.NewBasicJoinCondition("score", "score", ">")
		leftFact = createTestFact("f1", "Result", map[string]interface{}{"score": 85.5})
		token = createTestToken("t1", "node1", []*domain.Fact{leftFact})
		rightFact = createTestFact("f2", "Result", map[string]interface{}{"score": 80.0})

		result = condition.Evaluate(token, rightFact)
		if !result {
			t.Error("Float comparison 85.5 > 80.0 should be true")
		}
	})

	t.Run("BetaMemory_Size", func(t *testing.T) {
		memory := NewBetaMemory()

		// Test initial size
		tokens, facts := memory.Size()
		if tokens != 0 || facts != 0 {
			t.Errorf("Empty memory should report 0,0 size, got %d,%d", tokens, facts)
		}

		// Add items and check size
		fact := createTestFact("f1", "Person", map[string]interface{}{"name": "John"})
		token := createTestToken("t1", "node1", []*domain.Fact{fact})

		memory.StoreToken(token)
		memory.StoreFact(fact)

		tokens, facts = memory.Size()
		if tokens != 1 || facts != 1 {
			t.Errorf("Memory with 1 token and 1 fact should report 1,1 size, got %d,%d", tokens, facts)
		}
	})

	t.Run("BetaMemory_RemoveNonExistent", func(t *testing.T) {
		memory := NewBetaMemory()

		// Try to remove non-existent token and fact
		tokenRemoved := memory.RemoveToken("nonexistent")
		factRemoved := memory.RemoveFact("nonexistent")

		if tokenRemoved {
			t.Error("Removing non-existent token should return false")
		}
		if factRemoved {
			t.Error("Removing non-existent fact should return false")
		}
	})

	t.Run("JoinCondition_EmptyToken", func(t *testing.T) {
		condition := domain.NewBasicJoinCondition("name", "name", "==")
		emptyToken := createTestToken("t1", "node1", []*domain.Fact{}) // Empty facts slice
		rightFact := createTestFact("f2", "Person", map[string]interface{}{"name": "John"})

		result := condition.Evaluate(emptyToken, rightFact)
		if result {
			t.Error("Condition evaluation with empty token should return false")
		}
	})

	t.Run("JoinNode_MultipleConditions", func(t *testing.T) {
		node := NewJoinNode("join1", logger)

		// Multiple conditions (AND logic)
		condition1 := domain.NewBasicJoinCondition("age", "age", "==")
		condition2 := domain.NewBasicJoinCondition("city", "city", "==")
		node.SetJoinConditions([]domain.JoinCondition{condition1, condition2})

		capturedTokens := []*domain.Token{}
		mockChild := &mockBetaNode{
			processLeftToken: func(t *domain.Token) error {
				capturedTokens = append(capturedTokens, t)
				return nil
			},
		}
		node.AddChild(mockChild)

		// Token that matches on age but not city
		leftFact := createTestFact("f1", "Person", map[string]interface{}{"age": 25, "city": "Paris"})
		token := createTestToken("t1", "node1", []*domain.Fact{leftFact})

		rightFact1 := createTestFact("f2", "Person", map[string]interface{}{"age": 25, "city": "London"}) // age matches, city doesn't
		rightFact2 := createTestFact("f3", "Person", map[string]interface{}{"age": 25, "city": "Paris"})  // both match

		node.ProcessRightFact(rightFact1)
		node.ProcessRightFact(rightFact2)
		node.ProcessLeftToken(token)

		// Only one join should succeed (with rightFact2)
		if len(capturedTokens) != 1 {
			t.Errorf("Expected 1 join with multiple conditions, got %d", len(capturedTokens))
		}
	})

	t.Run("PropagateToNonBetaChildren", func(t *testing.T) {
		node := NewBaseBetaNode("beta1", "BetaNode", logger)

		// Mock non-beta child
		capturedFacts := []*domain.Fact{}
		mockChild := &mockNonBetaNode{
			processFact: func(f *domain.Fact) error {
				capturedFacts = append(capturedFacts, f)
				return nil
			},
		}
		node.AddChild(mockChild)

		// Process a token with multiple facts
		fact1 := createTestFact("f1", "Person", map[string]interface{}{"name": "John"})
		fact2 := createTestFact("f2", "Address", map[string]interface{}{"city": "Paris"})
		token := createTestToken("t1", "node1", []*domain.Fact{fact1, fact2})

		// Trigger join (we'll use a right fact to cause propagation)
		rightFact := createTestFact("f3", "Company", map[string]interface{}{"name": "ACME"})

		node.ProcessLeftToken(token)
		node.ProcessRightFact(rightFact)

		// Should have propagated the last fact of the joined token to non-beta child
		if len(capturedFacts) != 1 {
			t.Errorf("Expected 1 fact propagated to non-beta child, got %d", len(capturedFacts))
		}
		if len(capturedFacts) > 0 && capturedFacts[0].ID != "f3" {
			t.Errorf("Expected last fact 'f3' to be propagated, got '%s'", capturedFacts[0].ID)
		}
	})
}

func TestConcurrency(t *testing.T) {
	logger := &mockLogger{}

	t.Run("BetaMemory_ConcurrentAccess", func(t *testing.T) {
		memory := NewBetaMemory()

		const numGoroutines = 10
		const operationsPerGoroutine = 100

		var wg sync.WaitGroup

		// Lanceur pour écriture concurrent de tokens
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < operationsPerGoroutine; j++ {
					fact := createTestFact(fmt.Sprintf("f%d_%d", id, j), "Test", map[string]interface{}{"id": j})
					token := createTestToken(fmt.Sprintf("t%d_%d", id, j), "node", []*domain.Fact{fact})
					memory.StoreToken(token)
				}
			}(i)
		}

		// Lanceur pour écriture concurrent de facts
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < operationsPerGoroutine; j++ {
					fact := createTestFact(fmt.Sprintf("cf%d_%d", id, j), "Test", map[string]interface{}{"id": j})
					memory.StoreFact(fact)
				}
			}(i)
		}

		// Lanceur pour lecture concurrent
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < operationsPerGoroutine; j++ {
					memory.GetTokens()
					memory.GetFacts()
					memory.Size()
				}
			}()
		}

		wg.Wait()

		// Vérifier l'état final
		tokens, facts := memory.Size()
		expectedTokens := numGoroutines * operationsPerGoroutine
		expectedFacts := numGoroutines * operationsPerGoroutine

		if tokens != expectedTokens {
			t.Errorf("Expected %d tokens, got %d", expectedTokens, tokens)
		}
		if facts != expectedFacts {
			t.Errorf("Expected %d facts, got %d", expectedFacts, facts)
		}
	})

	t.Run("JoinNode_ConcurrentProcessing", func(t *testing.T) {
		node := NewJoinNode("join1", logger)

		const numGoroutines = 5
		var wg sync.WaitGroup

		capturedTokens := make([]*domain.Token, 0)
		var capturedMutex sync.Mutex

		mockChild := &mockBetaNode{
			processLeftToken: func(t *domain.Token) error {
				capturedMutex.Lock()
				capturedTokens = append(capturedTokens, t)
				capturedMutex.Unlock()
				return nil
			},
		}
		node.AddChild(mockChild)

		// Traiter des tokens en parallèle
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				fact := createTestFact(fmt.Sprintf("f%d", id), "Person", map[string]interface{}{"id": id})
				token := createTestToken(fmt.Sprintf("t%d", id), "node", []*domain.Fact{fact})
				node.ProcessLeftToken(token)
			}(i)
		}

		// Traiter des faits en parallèle
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				fact := createTestFact(fmt.Sprintf("rf%d", id), "Person", map[string]interface{}{"id": id})
				node.ProcessRightFact(fact)
			}(i)
		}

		wg.Wait()

		// Donner un peu de temps pour que tous les traitements soient terminés
		time.Sleep(100 * time.Millisecond)

		capturedMutex.Lock()
		joinCount := len(capturedTokens)
		capturedMutex.Unlock()

		// Chaque token devrait être joint avec chaque fait
		expectedJoins := numGoroutines * numGoroutines
		if joinCount != expectedJoins {
			t.Errorf("Expected %d joins, got %d", expectedJoins, joinCount)
		}
	})
}

// Mock BetaNode pour les tests
type mockBetaNode struct {
	id               string
	processLeftToken func(*domain.Token) error
	processRightFact func(*domain.Fact) error
}

func (m *mockBetaNode) ID() string                                    { return m.id }
func (m *mockBetaNode) Type() string                                  { return "MockBetaNode" }
func (m *mockBetaNode) ProcessFact(fact *domain.Fact) error           { return m.processRightFact(fact) }
func (m *mockBetaNode) GetMemory() *domain.WorkingMemory              { return nil }
func (m *mockBetaNode) AddChild(domain.Node)                          {}
func (m *mockBetaNode) GetChildren() []domain.Node                    { return nil }
func (m *mockBetaNode) ProcessLeftToken(token *domain.Token) error    { return m.processLeftToken(token) }
func (m *mockBetaNode) ProcessRightFact(fact *domain.Fact) error      { return m.processRightFact(fact) }
func (m *mockBetaNode) GetLeftMemory() []*domain.Token                { return nil }
func (m *mockBetaNode) GetRightMemory() []*domain.Fact                { return nil }
func (m *mockBetaNode) ClearMemory()                                  {}
func (m *mockBetaNode) SetJoinConditions([]domain.JoinCondition)      {}
func (m *mockBetaNode) GetJoinConditions() []domain.JoinCondition     { return nil }
func (m *mockBetaNode) EvaluateJoin(*domain.Token, *domain.Fact) bool { return true }

// Mock NonBetaNode pour les tests
type mockNonBetaNode struct {
	id          string
	processFact func(*domain.Fact) error
}

func (m *mockNonBetaNode) ID() string                          { return m.id }
func (m *mockNonBetaNode) Type() string                        { return "MockNonBetaNode" }
func (m *mockNonBetaNode) ProcessFact(fact *domain.Fact) error { return m.processFact(fact) }
