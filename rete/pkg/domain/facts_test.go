package domain

import (
	"testing"
	"time"
)

// ===== Fact Tests =====

func TestNewFact(t *testing.T) {
	fields := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	fact := NewFact("fact1", "Person", fields)

	if fact == nil {
		t.Fatal("NewFact should not return nil")
	}
	if fact.ID != "fact1" {
		t.Errorf("Expected ID 'fact1', got '%s'", fact.ID)
	}
	if fact.Type != "Person" {
		t.Errorf("Expected Type 'Person', got '%s'", fact.Type)
	}
	if len(fact.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(fact.Fields))
	}
	if fact.Timestamp.IsZero() {
		t.Error("Timestamp should be set")
	}
	if time.Since(fact.Timestamp) > time.Second {
		t.Error("Timestamp should be recent")
	}
}

func TestFact_String(t *testing.T) {
	fact := NewFact("fact1", "Person", nil)

	result := fact.String()
	expected := "fact1:Person"

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFact_GetField(t *testing.T) {
	fields := map[string]interface{}{
		"name":   "Alice",
		"age":    25,
		"salary": 50000.0,
		"active": true,
	}
	fact := NewFact("fact1", "Person", fields)

	tests := []struct {
		name       string
		fieldName  string
		wantValue  interface{}
		wantExists bool
	}{
		{"existing string field", "name", "Alice", true},
		{"existing int field", "age", 25, true},
		{"existing float field", "salary", 50000.0, true},
		{"existing bool field", "active", true, true},
		{"non-existing field", "unknown", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, exists := fact.GetField(tt.fieldName)
			if exists != tt.wantExists {
				t.Errorf("GetField(%s) exists = %v, want %v", tt.fieldName, exists, tt.wantExists)
			}
			if tt.wantExists && value != tt.wantValue {
				t.Errorf("GetField(%s) value = %v, want %v", tt.fieldName, value, tt.wantValue)
			}
		})
	}
}

func TestFact_EmptyFields(t *testing.T) {
	fact := NewFact("fact1", "Empty", nil)

	// NewFact with nil fields will result in fact.Fields being nil
	// This is acceptable behavior

	_, exists := fact.GetField("any")
	if exists {
		t.Error("GetField should return false for non-existing field")
	}
}

func TestFact_FieldTypes(t *testing.T) {
	fields := map[string]interface{}{
		"string": "value",
		"int":    42,
		"float":  3.14,
		"bool":   true,
		"nil":    nil,
		"slice":  []int{1, 2, 3},
		"map":    map[string]string{"key": "value"},
	}

	fact := NewFact("fact1", "Mixed", fields)

	for fieldName, expectedValue := range fields {
		value, exists := fact.GetField(fieldName)
		if !exists {
			t.Errorf("Field %s should exist", fieldName)
		}
		// Pour slice et map, on vérifie juste qu'ils ne sont pas nil
		if expectedValue == nil {
			if value != nil {
				t.Errorf("Field %s should be nil, got %v", fieldName, value)
			}
		}
	}
}

// ===== Token Tests =====

func TestNewToken(t *testing.T) {
	fact1 := NewFact("f1", "Type1", map[string]interface{}{"x": 1})
	fact2 := NewFact("f2", "Type2", map[string]interface{}{"y": 2})
	facts := []*Fact{fact1, fact2}

	token := NewToken("token1", "node1", facts)

	if token == nil {
		t.Fatal("NewToken should not return nil")
	}
	if token.ID != "token1" {
		t.Errorf("Expected ID 'token1', got '%s'", token.ID)
	}
	if token.NodeID != "node1" {
		t.Errorf("Expected NodeID 'node1', got '%s'", token.NodeID)
	}
	if len(token.Facts) != 2 {
		t.Errorf("Expected 2 facts, got %d", len(token.Facts))
	}
	if token.Parent != nil {
		t.Error("Parent should be nil by default")
	}
}

func TestToken_EmptyFacts(t *testing.T) {
	token := NewToken("token1", "node1", []*Fact{})

	if token == nil {
		t.Fatal("NewToken should not return nil")
	}
	if len(token.Facts) != 0 {
		t.Errorf("Expected 0 facts, got %d", len(token.Facts))
	}
}

func TestToken_WithParent(t *testing.T) {
	fact1 := NewFact("f1", "Type1", nil)
	fact2 := NewFact("f2", "Type2", nil)

	parentToken := NewToken("parent", "node1", []*Fact{fact1})
	childToken := NewToken("child", "node2", []*Fact{fact1, fact2})
	childToken.Parent = parentToken

	if childToken.Parent == nil {
		t.Error("Parent should be set")
	}
	if childToken.Parent.ID != "parent" {
		t.Errorf("Expected parent ID 'parent', got '%s'", childToken.Parent.ID)
	}
}

// ===== WorkingMemory Tests =====

func TestNewWorkingMemory(t *testing.T) {
	wm := NewWorkingMemory("node1")

	if wm == nil {
		t.Fatal("NewWorkingMemory should not return nil")
	}
	if wm.NodeID != "node1" {
		t.Errorf("Expected NodeID 'node1', got '%s'", wm.NodeID)
	}
	if wm.Facts == nil {
		t.Error("Facts map should be initialized")
	}
	if wm.Tokens == nil {
		t.Error("Tokens map should be initialized")
	}
	if len(wm.Facts) != 0 {
		t.Errorf("Expected 0 facts, got %d", len(wm.Facts))
	}
	if len(wm.Tokens) != 0 {
		t.Errorf("Expected 0 tokens, got %d", len(wm.Tokens))
	}
}

func TestWorkingMemory_AddFact(t *testing.T) {
	wm := NewWorkingMemory("node1")
	fact1 := NewFact("f1", "Type1", map[string]interface{}{"x": 1})
	fact2 := NewFact("f2", "Type2", map[string]interface{}{"y": 2})

	wm.AddFact(fact1)
	if len(wm.Facts) != 1 {
		t.Errorf("Expected 1 fact, got %d", len(wm.Facts))
	}

	wm.AddFact(fact2)
	if len(wm.Facts) != 2 {
		t.Errorf("Expected 2 facts, got %d", len(wm.Facts))
	}

	// Vérifier que les faits sont présents
	if _, exists := wm.Facts["f1"]; !exists {
		t.Error("Fact f1 should exist")
	}
	if _, exists := wm.Facts["f2"]; !exists {
		t.Error("Fact f2 should exist")
	}
}

func TestWorkingMemory_AddFact_NilMap(t *testing.T) {
	wm := &WorkingMemory{NodeID: "node1"} // Facts map non initialisé
	fact := NewFact("f1", "Type1", nil)

	wm.AddFact(fact)

	if wm.Facts == nil {
		t.Error("Facts map should be initialized after AddFact")
	}
	if len(wm.Facts) != 1 {
		t.Errorf("Expected 1 fact, got %d", len(wm.Facts))
	}
}

func TestWorkingMemory_RemoveFact(t *testing.T) {
	wm := NewWorkingMemory("node1")
	fact1 := NewFact("f1", "Type1", nil)
	fact2 := NewFact("f2", "Type2", nil)

	wm.AddFact(fact1)
	wm.AddFact(fact2)

	if len(wm.Facts) != 2 {
		t.Fatalf("Expected 2 facts before removal, got %d", len(wm.Facts))
	}

	wm.RemoveFact("f1")

	if len(wm.Facts) != 1 {
		t.Errorf("Expected 1 fact after removal, got %d", len(wm.Facts))
	}
	if _, exists := wm.Facts["f1"]; exists {
		t.Error("Fact f1 should not exist after removal")
	}
	if _, exists := wm.Facts["f2"]; !exists {
		t.Error("Fact f2 should still exist")
	}
}

func TestWorkingMemory_RemoveFact_NonExisting(t *testing.T) {
	wm := NewWorkingMemory("node1")
	fact := NewFact("f1", "Type1", nil)
	wm.AddFact(fact)

	// Supprimer un fait inexistant ne devrait pas causer d'erreur
	wm.RemoveFact("nonexistent")

	if len(wm.Facts) != 1 {
		t.Errorf("Expected 1 fact, got %d", len(wm.Facts))
	}
}

func TestWorkingMemory_GetFacts(t *testing.T) {
	wm := NewWorkingMemory("node1")

	// Empty memory
	facts := wm.GetFacts()
	if len(facts) != 0 {
		t.Errorf("Expected 0 facts, got %d", len(facts))
	}

	// Add facts
	fact1 := NewFact("f1", "Type1", nil)
	fact2 := NewFact("f2", "Type2", nil)
	fact3 := NewFact("f3", "Type3", nil)

	wm.AddFact(fact1)
	wm.AddFact(fact2)
	wm.AddFact(fact3)

	facts = wm.GetFacts()
	if len(facts) != 3 {
		t.Errorf("Expected 3 facts, got %d", len(facts))
	}

	// Vérifier que tous les faits sont présents
	factIDs := make(map[string]bool)
	for _, f := range facts {
		factIDs[f.ID] = true
	}
	if !factIDs["f1"] || !factIDs["f2"] || !factIDs["f3"] {
		t.Error("Not all facts are returned")
	}
}

func TestWorkingMemory_AddToken(t *testing.T) {
	wm := NewWorkingMemory("node1")
	token1 := NewToken("t1", "node1", []*Fact{})
	token2 := NewToken("t2", "node1", []*Fact{})

	wm.AddToken(token1)
	if len(wm.Tokens) != 1 {
		t.Errorf("Expected 1 token, got %d", len(wm.Tokens))
	}

	wm.AddToken(token2)
	if len(wm.Tokens) != 2 {
		t.Errorf("Expected 2 tokens, got %d", len(wm.Tokens))
	}

	if _, exists := wm.Tokens["t1"]; !exists {
		t.Error("Token t1 should exist")
	}
	if _, exists := wm.Tokens["t2"]; !exists {
		t.Error("Token t2 should exist")
	}
}

func TestWorkingMemory_AddToken_NilMap(t *testing.T) {
	wm := &WorkingMemory{NodeID: "node1"} // Tokens map non initialisé
	token := NewToken("t1", "node1", []*Fact{})

	wm.AddToken(token)

	if wm.Tokens == nil {
		t.Error("Tokens map should be initialized after AddToken")
	}
	if len(wm.Tokens) != 1 {
		t.Errorf("Expected 1 token, got %d", len(wm.Tokens))
	}
}

func TestWorkingMemory_RemoveToken(t *testing.T) {
	wm := NewWorkingMemory("node1")
	token1 := NewToken("t1", "node1", []*Fact{})
	token2 := NewToken("t2", "node1", []*Fact{})

	wm.AddToken(token1)
	wm.AddToken(token2)

	if len(wm.Tokens) != 2 {
		t.Fatalf("Expected 2 tokens before removal, got %d", len(wm.Tokens))
	}

	wm.RemoveToken("t1")

	if len(wm.Tokens) != 1 {
		t.Errorf("Expected 1 token after removal, got %d", len(wm.Tokens))
	}
	if _, exists := wm.Tokens["t1"]; exists {
		t.Error("Token t1 should not exist after removal")
	}
	if _, exists := wm.Tokens["t2"]; !exists {
		t.Error("Token t2 should still exist")
	}
}

func TestWorkingMemory_GetTokens(t *testing.T) {
	wm := NewWorkingMemory("node1")

	// Empty memory
	tokens := wm.GetTokens()
	if len(tokens) != 0 {
		t.Errorf("Expected 0 tokens, got %d", len(tokens))
	}

	// Add tokens
	token1 := NewToken("t1", "node1", []*Fact{})
	token2 := NewToken("t2", "node1", []*Fact{})
	token3 := NewToken("t3", "node1", []*Fact{})

	wm.AddToken(token1)
	wm.AddToken(token2)
	wm.AddToken(token3)

	tokens = wm.GetTokens()
	if len(tokens) != 3 {
		t.Errorf("Expected 3 tokens, got %d", len(tokens))
	}

	// Vérifier que tous les tokens sont présents
	tokenIDs := make(map[string]bool)
	for _, tk := range tokens {
		tokenIDs[tk.ID] = true
	}
	if !tokenIDs["t1"] || !tokenIDs["t2"] || !tokenIDs["t3"] {
		t.Error("Not all tokens are returned")
	}
}

// ===== BasicJoinCondition Tests =====

func TestNewBasicJoinCondition(t *testing.T) {
	condition := NewBasicJoinCondition("leftField", "rightField", "==")

	if condition == nil {
		t.Fatal("NewBasicJoinCondition should not return nil")
	}
	if condition.LeftField != "leftField" {
		t.Errorf("Expected LeftField 'leftField', got '%s'", condition.LeftField)
	}
	if condition.RightField != "rightField" {
		t.Errorf("Expected RightField 'rightField', got '%s'", condition.RightField)
	}
	if condition.Operator != "==" {
		t.Errorf("Expected Operator '==', got '%s'", condition.Operator)
	}
}

func TestBasicJoinCondition_GetMethods(t *testing.T) {
	condition := NewBasicJoinCondition("left", "right", "!=")

	if condition.GetLeftField() != "left" {
		t.Errorf("Expected 'left', got '%s'", condition.GetLeftField())
	}
	if condition.GetRightField() != "right" {
		t.Errorf("Expected 'right', got '%s'", condition.GetRightField())
	}
	if condition.GetOperator() != "!=" {
		t.Errorf("Expected '!=', got '%s'", condition.GetOperator())
	}
}

func TestBasicJoinCondition_Evaluate_EmptyToken(t *testing.T) {
	condition := NewBasicJoinCondition("age", "age", "==")
	token := NewToken("t1", "node1", []*Fact{})
	fact := NewFact("f1", "Person", map[string]interface{}{"age": 30})

	result := condition.Evaluate(token, fact)
	if result {
		t.Error("Evaluate should return false for empty token")
	}
}

func TestBasicJoinCondition_Evaluate_Equality(t *testing.T) {
	tests := []struct {
		name       string
		leftValue  interface{}
		rightValue interface{}
		operator   string
		want       bool
	}{
		{"equal integers", 42, 42, "==", true},
		{"equal strings", "test", "test", "==", true},
		{"equal floats", 3.14, 3.14, "==", true},
		{"equal bools", true, true, "==", true},
		{"not equal integers", 42, 43, "==", false},
		{"not equal strings", "test", "other", "==", false},
		{"alternative equality operator", 42, 42, "=", true},
		{"not equal operator", 42, 42, "!=", false},
		{"not equal operator different values", 42, 43, "!=", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := NewBasicJoinCondition("field", "field", tt.operator)

			leftFact := NewFact("f1", "Type1", map[string]interface{}{"field": tt.leftValue})
			token := NewToken("t1", "node1", []*Fact{leftFact})
			rightFact := NewFact("f2", "Type2", map[string]interface{}{"field": tt.rightValue})

			result := condition.Evaluate(token, rightFact)
			if result != tt.want {
				t.Errorf("Evaluate() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestBasicJoinCondition_Evaluate_Comparison(t *testing.T) {
	tests := []struct {
		name       string
		leftValue  interface{}
		rightValue interface{}
		operator   string
		want       bool
	}{
		// Integer comparisons
		{"int less than true", 10, 20, "<", true},
		{"int less than false", 20, 10, "<", false},
		{"int less than equal", 10, 10, "<", false},
		{"int less or equal true", 10, 20, "<=", true},
		{"int less or equal equal", 10, 10, "<=", true},
		{"int less or equal false", 20, 10, "<=", false},
		{"int greater than true", 20, 10, ">", true},
		{"int greater than false", 10, 20, ">", false},
		{"int greater or equal true", 20, 10, ">=", true},
		{"int greater or equal equal", 10, 10, ">=", true},
		{"int greater or equal false", 10, 20, ">=", false},

		// Float comparisons
		{"float less than", 1.5, 2.5, "<", true},
		{"float greater than", 3.14, 2.71, ">", true},
		{"float equal", 3.14, 3.14, "<=", true},

		// String comparisons
		{"string less than", "apple", "banana", "<", true},
		{"string greater than", "zebra", "apple", ">", true},
		{"string equal", "test", "test", "<=", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := NewBasicJoinCondition("field", "field", tt.operator)

			leftFact := NewFact("f1", "Type1", map[string]interface{}{"field": tt.leftValue})
			token := NewToken("t1", "node1", []*Fact{leftFact})
			rightFact := NewFact("f2", "Type2", map[string]interface{}{"field": tt.rightValue})

			result := condition.Evaluate(token, rightFact)
			if result != tt.want {
				t.Errorf("Evaluate() = %v, want %v (left=%v, right=%v, op=%s)",
					result, tt.want, tt.leftValue, tt.rightValue, tt.operator)
			}
		})
	}
}

func TestBasicJoinCondition_Evaluate_MissingFields(t *testing.T) {
	condition := NewBasicJoinCondition("leftField", "rightField", "==")

	// Left field missing
	leftFact := NewFact("f1", "Type1", map[string]interface{}{"other": 1})
	token := NewToken("t1", "node1", []*Fact{leftFact})
	rightFact := NewFact("f2", "Type2", map[string]interface{}{"rightField": 1})

	result := condition.Evaluate(token, rightFact)
	if result {
		t.Error("Evaluate should return false when left field is missing")
	}

	// Right field missing
	leftFact = NewFact("f1", "Type1", map[string]interface{}{"leftField": 1})
	token = NewToken("t1", "node1", []*Fact{leftFact})
	rightFact = NewFact("f2", "Type2", map[string]interface{}{"other": 1})

	result = condition.Evaluate(token, rightFact)
	if result {
		t.Error("Evaluate should return false when right field is missing")
	}
}

func TestBasicJoinCondition_Evaluate_InvalidOperator(t *testing.T) {
	condition := NewBasicJoinCondition("field", "field", "INVALID")

	leftFact := NewFact("f1", "Type1", map[string]interface{}{"field": 42})
	token := NewToken("t1", "node1", []*Fact{leftFact})
	rightFact := NewFact("f2", "Type2", map[string]interface{}{"field": 42})

	result := condition.Evaluate(token, rightFact)
	if result {
		t.Error("Evaluate should return false for invalid operator")
	}
}

func TestBasicJoinCondition_Evaluate_TypeMismatch(t *testing.T) {
	condition := NewBasicJoinCondition("field", "field", "<")

	// Comparing int with string
	leftFact := NewFact("f1", "Type1", map[string]interface{}{"field": 42})
	token := NewToken("t1", "node1", []*Fact{leftFact})
	rightFact := NewFact("f2", "Type2", map[string]interface{}{"field": "string"})

	result := condition.Evaluate(token, rightFact)
	// Should handle gracefully and return false
	if result {
		t.Error("Evaluate should return false for type mismatch in comparison")
	}
}

func TestBasicJoinCondition_Evaluate_MultipleFactsInToken(t *testing.T) {
	condition := NewBasicJoinCondition("age", "age", "==")

	// Token with multiple facts - should use the last one
	fact1 := NewFact("f1", "Person", map[string]interface{}{"age": 25})
	fact2 := NewFact("f2", "Person", map[string]interface{}{"age": 30})
	token := NewToken("t1", "node1", []*Fact{fact1, fact2})

	rightFact := NewFact("f3", "Person", map[string]interface{}{"age": 30})

	result := condition.Evaluate(token, rightFact)
	if !result {
		t.Error("Evaluate should use the last fact in token and return true")
	}
}

// ===== Error Tests =====

func TestValidationError(t *testing.T) {
	err := NewValidationError("age", -5, "must be positive")

	if err == nil {
		t.Fatal("NewValidationError should not return nil")
	}
	if err.Field != "age" {
		t.Errorf("Expected Field 'age', got '%s'", err.Field)
	}
	if err.Value != -5 {
		t.Errorf("Expected Value -5, got %v", err.Value)
	}
	if err.Message != "must be positive" {
		t.Errorf("Expected Message 'must be positive', got '%s'", err.Message)
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("Error() should return non-empty string")
	}
	// Vérifier que le message contient les informations clés
	expectedParts := []string{"validation error", "age", "must be positive", "-5"}
	for _, part := range expectedParts {
		found := false
		for i := 0; i <= len(errMsg)-len(part); i++ {
			if errMsg[i:i+len(part)] == part {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Error message should contain '%s', got '%s'", part, errMsg)
		}
	}
}

func TestNodeError(t *testing.T) {
	cause := NewValidationError("field", "value", "invalid")
	err := NewNodeError("node1", "JoinNode", cause)

	if err == nil {
		t.Fatal("NewNodeError should not return nil")
	}
	if err.NodeID != "node1" {
		t.Errorf("Expected NodeID 'node1', got '%s'", err.NodeID)
	}
	if err.NodeType != "JoinNode" {
		t.Errorf("Expected NodeType 'JoinNode', got '%s'", err.NodeType)
	}
	if err.Cause == nil {
		t.Error("Cause should not be nil")
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("Error() should return non-empty string")
	}

	// Test Unwrap
	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Error("Unwrap() should return the original cause")
	}
}

func TestPredefinedErrors(t *testing.T) {
	errors := []error{
		ErrFactNotFound,
		ErrInvalidFactType,
		ErrInvalidFieldType,
		ErrNodeNotFound,
		ErrStorageError,
		ErrValidationFailed,
	}

	for _, err := range errors {
		if err == nil {
			t.Error("Predefined error should not be nil")
		}
		if err.Error() == "" {
			t.Error("Predefined error should have non-empty message")
		}
	}
}
