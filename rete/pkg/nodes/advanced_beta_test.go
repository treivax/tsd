// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package nodes

import (
	"testing"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// ===== NotNodeImpl Tests =====

func TestNewNotNode(t *testing.T) {
	logger := newMockLogger()
	node := NewNotNode("not1", logger)

	if node == nil {
		t.Fatal("NewNotNode should not return nil")
	}
	if node.ID() != "not1" {
		t.Errorf("Expected ID 'not1', got '%s'", node.ID())
	}
	if node.Type() != "NotNode" {
		t.Errorf("Expected Type 'NotNode', got '%s'", node.Type())
	}
	if node.negationCondition != nil {
		t.Error("Initial negation condition should be nil")
	}
}

func TestNotNode_SetGetNegationCondition(t *testing.T) {
	logger := newMockLogger()
	node := NewNotNode("not1", logger)

	// Initially nil
	if node.GetNegationCondition() != nil {
		t.Error("Initial condition should be nil")
	}

	// Set simple condition
	condition := map[string]interface{}{
		"type":     "simple",
		"field":    "age",
		"operator": "==",
		"value":    0,
	}

	node.SetNegationCondition(condition)

	retrieved := node.GetNegationCondition()
	if retrieved == nil {
		t.Fatal("Condition should not be nil after setting")
	}

	condMap, ok := retrieved.(map[string]interface{})
	if !ok {
		t.Fatal("Retrieved condition should be map[string]interface{}")
	}
	if condMap["type"] != "simple" {
		t.Errorf("Expected type 'simple', got '%v'", condMap["type"])
	}
}

func TestNotNode_ProcessNegation_NilCondition(t *testing.T) {
	logger := newMockLogger()
	node := NewNotNode("not1", logger)

	fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": 25})
	token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

	// With nil condition, should return false
	result := node.ProcessNegation(token, fact)
	if result {
		t.Error("ProcessNegation should return false when condition is nil")
	}
}

func TestNotNode_ProcessNegation_SimpleCondition(t *testing.T) {
	logger := newMockLogger()
	node := NewNotNode("not1", logger)

	// Set simple condition (always true)
	condition := map[string]interface{}{
		"type": "simple",
	}
	node.SetNegationCondition(condition)

	fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": 25})
	token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

	// Simple condition is true, negation should be false
	result := node.ProcessNegation(token, fact)
	if result {
		t.Error("ProcessNegation should return false (negation of true)")
	}
}

func TestNotNode_ProcessNegation_BinaryCondition(t *testing.T) {
	logger := newMockLogger()
	node := NewNotNode("not1", logger)

	tests := []struct {
		name      string
		condition map[string]interface{}
		factAge   int
		wantNeg   bool // true if negation succeeds (condition is false)
	}{
		{
			name: "age == 0, fact has age 0 - negation fails",
			condition: map[string]interface{}{
				"type": "binaryOperation",
				"left": map[string]interface{}{
					"field": "age",
				},
				"operator": "==",
				"right": map[string]interface{}{
					"value": 0,
				},
			},
			factAge: 0,
			wantNeg: false, // condition true, negation false
		},
		{
			name: "age == 0, fact has age 25 - negation succeeds",
			condition: map[string]interface{}{
				"type": "binaryOperation",
				"left": map[string]interface{}{
					"field": "age",
				},
				"operator": "==",
				"right": map[string]interface{}{
					"value": 0,
				},
			},
			factAge: 25,
			wantNeg: true, // condition false, negation true
		},
		{
			name: "age != 0, fact has age 25 - negation fails",
			condition: map[string]interface{}{
				"type": "binary_op",
				"left": map[string]interface{}{
					"field": "age",
				},
				"op": "!=",
				"right": map[string]interface{}{
					"value": 0,
				},
			},
			factAge: 25,
			wantNeg: false, // condition true, negation false
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node.SetNegationCondition(tt.condition)

			fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": tt.factAge})
			token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

			result := node.ProcessNegation(token, fact)
			if result != tt.wantNeg {
				t.Errorf("ProcessNegation() = %v, want %v", result, tt.wantNeg)
			}
		})
	}
}

func TestNotNode_ProcessLeftToken_NoRightFacts(t *testing.T) {
	logger := newMockLogger()
	node := NewNotNode("not1", logger)

	fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": 25})
	token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

	// Add a child to verify propagation
	childLogger := newMockLogger()
	child := NewBaseBetaNode("child1", "BetaNode", childLogger)
	node.AddChild(child)

	err := node.ProcessLeftToken(token)
	if err != nil {
		t.Errorf("ProcessLeftToken should not return error: %v", err)
	}

	// Token should be stored in left memory
	leftTokens := node.GetLeftMemory()
	if len(leftTokens) != 1 {
		t.Errorf("Expected 1 token in left memory, got %d", len(leftTokens))
	}

	// With no right facts, token should be propagated (verify via child's left memory)
	childLeftTokens := child.GetLeftMemory()
	if len(childLeftTokens) != 1 {
		t.Errorf("Expected token to be propagated to child, got %d tokens", len(childLeftTokens))
	}
}

func TestNotNode_ProcessLeftToken_WithRightFacts(t *testing.T) {
	logger := newMockLogger()
	node := NewNotNode("not1", logger)

	// Set condition: age == 0
	condition := map[string]interface{}{
		"type": "binaryOperation",
		"left": map[string]interface{}{
			"field": "age",
		},
		"operator": "==",
		"right": map[string]interface{}{
			"value": 0,
		},
	}
	node.SetNegationCondition(condition)

	// Add a child
	childLogger := newMockLogger()
	child := NewBaseBetaNode("child1", "BetaNode", childLogger)
	node.AddChild(child)

	// Add right fact with age != 0 (condition false, negation true)
	rightFact := domain.NewFact("f_right", "Person", map[string]interface{}{"age": 25})
	node.ProcessRightFact(rightFact)

	// Process left token with age != 0
	leftFact := domain.NewFact("f_left", "Person", map[string]interface{}{"age": 30})
	token := domain.NewToken("t1", "node1", []*domain.Fact{leftFact})

	err := node.ProcessLeftToken(token)
	if err != nil {
		t.Errorf("ProcessLeftToken should not return error: %v", err)
	}

	// Token should be propagated (negation succeeds)
	childLeftTokens := child.GetLeftMemory()
	if len(childLeftTokens) != 1 {
		t.Errorf("Expected token to be propagated, got %d tokens", len(childLeftTokens))
	}
}

func TestNotNode_ProcessRightFact(t *testing.T) {
	logger := newMockLogger()
	node := NewNotNode("not1", logger)

	fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": 25})

	err := node.ProcessRightFact(fact)
	if err != nil {
		t.Errorf("ProcessRightFact should not return error: %v", err)
	}

	// Fact should be stored in right memory
	rightFacts := node.GetRightMemory()
	if len(rightFacts) != 1 {
		t.Errorf("Expected 1 fact in right memory, got %d", len(rightFacts))
	}
}

func TestNotNode_ProcessRightFact_WithLeftTokens(t *testing.T) {
	logger := newMockLogger()
	node := NewNotNode("not1", logger)

	// Set condition: age == 0
	condition := map[string]interface{}{
		"type": "binaryOperation",
		"left": map[string]interface{}{
			"field": "age",
		},
		"operator": "==",
		"right": map[string]interface{}{
			"value": 0,
		},
	}
	node.SetNegationCondition(condition)

	// Add left token first
	leftFact := domain.NewFact("f_left", "Person", map[string]interface{}{"age": 0})
	token := domain.NewToken("t1", "node1", []*domain.Fact{leftFact})
	node.betaMemory.StoreToken(token)

	// Add right fact with age == 0 (condition true, negation false)
	rightFact := domain.NewFact("f_right", "Person", map[string]interface{}{"age": 0})
	err := node.ProcessRightFact(rightFact)
	if err != nil {
		t.Errorf("ProcessRightFact should not return error: %v", err)
	}

	// Debug log should be called (negation evaluation happened)
	if logger.DebugCallCount() == 0 {
		t.Error("Expected debug logs to be called")
	}
}


func TestNotNode_ConcurrentAccess(t *testing.T) {
	logger := newMockLogger()
	node := NewNotNode("not1", logger)

	condition := map[string]interface{}{
		"type": "simple",
	}

	// Concurrent set/get of negation condition
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			node.SetNegationCondition(condition)
			node.GetNegationCondition()
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify final state
	if node.GetNegationCondition() == nil {
		t.Error("Condition should be set after concurrent access")
	}
}

// ===== ExistsNodeImpl Tests =====

func TestNewExistsNode(t *testing.T) {
	logger := newMockLogger()
	node := NewExistsNode("exists1", logger)

	if node == nil {
		t.Fatal("NewExistsNode should not return nil")
	}
	if node.ID() != "exists1" {
		t.Errorf("Expected ID 'exists1', got '%s'", node.ID())
	}
	if node.Type() != "ExistsNode" {
		t.Errorf("Expected Type 'ExistsNode', got '%s'", node.Type())
	}
}

func TestExistsNode_SetGetExistenceCondition(t *testing.T) {
	logger := newMockLogger()
	node := NewExistsNode("exists1", logger)

	variable := domain.TypedVariable{
		Name:     "p",
		DataType: "Person",
	}
	condition := map[string]interface{}{
		"type":  "exists",
		"field": "age",
	}

	// Initially, should return empty
	retVar, retCond := node.GetExistenceCondition()
	if retVar.Name != "" {
		t.Error("Initial variable should be empty")
	}
	if retCond != nil {
		t.Error("Initial condition should be nil")
	}

	// Set condition
	node.SetExistenceCondition(variable, condition)

	// Retrieve and verify
	retVar, retCond = node.GetExistenceCondition()
	if retVar.Name != "p" {
		t.Errorf("Expected variable name 'p', got '%s'", retVar.Name)
	}
	if retVar.DataType != "Person" {
		t.Errorf("Expected DataType 'Person', got '%s'", retVar.DataType)
	}
	if retCond == nil {
		t.Fatal("Condition should not be nil after setting")
	}
}

func TestExistsNode_CheckExistence_NilCondition(t *testing.T) {
	logger := newMockLogger()
	node := NewExistsNode("exists1", logger)

	fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": 25})
	token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

	result := node.CheckExistence(token)
	if result {
		t.Error("CheckExistence should return false when condition is nil")
	}
}

func TestExistsNode_CheckExistence_NoMatchingFacts(t *testing.T) {
	logger := newMockLogger()
	node := NewExistsNode("exists1", logger)

	variable := domain.TypedVariable{
		Name:     "p",
		DataType: "Person",
	}
	condition := map[string]interface{}{"type": "exists"}
	node.SetExistenceCondition(variable, condition)

	fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": 25})
	token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

	// No facts in right memory
	result := node.CheckExistence(token)
	if result {
		t.Error("CheckExistence should return false when no facts in right memory")
	}
}

func TestExistsNode_CheckExistence_WithMatchingFacts(t *testing.T) {
	logger := newMockLogger()
	node := NewExistsNode("exists1", logger)

	variable := domain.TypedVariable{
		Name:     "p",
		DataType: "Person",
	}
	condition := map[string]interface{}{"type": "exists"}
	node.SetExistenceCondition(variable, condition)

	// Add matching fact to right memory
	rightFact := domain.NewFact("f_right", "Person", map[string]interface{}{"age": 30})
	node.betaMemory.StoreFact(rightFact)

	fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": 25})
	token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

	result := node.CheckExistence(token)
	if !result {
		t.Error("CheckExistence should return true when matching fact exists")
	}
}

func TestExistsNode_CheckExistence_WrongType(t *testing.T) {
	logger := newMockLogger()
	node := NewExistsNode("exists1", logger)

	variable := domain.TypedVariable{
		Name:     "p",
		DataType: "Person",
	}
	condition := map[string]interface{}{"type": "exists"}
	node.SetExistenceCondition(variable, condition)

	// Add fact with wrong type
	rightFact := domain.NewFact("f_right", "Product", map[string]interface{}{"price": 100})
	node.betaMemory.StoreFact(rightFact)

	fact := domain.NewFact("f1", "Person", map[string]interface{}{"age": 25})
	token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

	result := node.CheckExistence(token)
	if result {
		t.Error("CheckExistence should return false when fact type doesn't match")
	}
}

func TestExistsNode_ProcessLeftToken_ExistenceSatisfied(t *testing.T) {
	logger := newMockLogger()
	node := NewExistsNode("exists1", logger)

	variable := domain.TypedVariable{
		Name:     "p",
		DataType: "Person",
	}
	condition := map[string]interface{}{"type": "exists"}
	node.SetExistenceCondition(variable, condition)

	// Add child
	childLogger := newMockLogger()
	child := NewBaseBetaNode("child1", "BetaNode", childLogger)
	node.AddChild(child)

	// Add matching fact to right memory
	rightFact := domain.NewFact("f_right", "Person", map[string]interface{}{"age": 30})
	node.betaMemory.StoreFact(rightFact)

	// Process left token
	fact := domain.NewFact("f1", "Order", map[string]interface{}{"id": 1})
	token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

	err := node.ProcessLeftToken(token)
	if err != nil {
		t.Errorf("ProcessLeftToken should not return error: %v", err)
	}

	// Token should be propagated
	childLeftTokens := child.GetLeftMemory()
	if len(childLeftTokens) != 1 {
		t.Errorf("Expected token to be propagated, got %d tokens", len(childLeftTokens))
	}

	// Token should be in left memory
	leftTokens := node.GetLeftMemory()
	if len(leftTokens) != 1 {
		t.Errorf("Expected 1 token in left memory, got %d", len(leftTokens))
	}
}

func TestExistsNode_ProcessLeftToken_ExistenceNotSatisfied(t *testing.T) {
	logger := newMockLogger()
	node := NewExistsNode("exists1", logger)

	variable := domain.TypedVariable{
		Name:     "p",
		DataType: "Person",
	}
	condition := map[string]interface{}{"type": "exists"}
	node.SetExistenceCondition(variable, condition)

	// Add child
	childLogger := newMockLogger()
	child := NewBaseBetaNode("child1", "BetaNode", childLogger)
	node.AddChild(child)

	// No matching facts in right memory

	// Process left token
	fact := domain.NewFact("f1", "Order", map[string]interface{}{"id": 1})
	token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

	err := node.ProcessLeftToken(token)
	if err != nil {
		t.Errorf("ProcessLeftToken should not return error: %v", err)
	}

	// Token should NOT be propagated
	childLeftTokens := child.GetLeftMemory()
	if len(childLeftTokens) != 0 {
		t.Errorf("Expected no propagation, got %d tokens", len(childLeftTokens))
	}

	// Token should still be in left memory
	leftTokens := node.GetLeftMemory()
	if len(leftTokens) != 1 {
		t.Errorf("Expected 1 token in left memory, got %d", len(leftTokens))
	}
}

func TestExistsNode_ProcessRightFact(t *testing.T) {
	logger := newMockLogger()
	node := NewExistsNode("exists1", logger)

	variable := domain.TypedVariable{
		Name:     "p",
		DataType: "Person",
	}
	condition := map[string]interface{}{"type": "exists"}
	node.SetExistenceCondition(variable, condition)

	// Add child
	childLogger := newMockLogger()
	child := NewBaseBetaNode("child1", "BetaNode", childLogger)
	node.AddChild(child)

	// Add left token first (existence not yet satisfied)
	leftFact := domain.NewFact("f_left", "Order", map[string]interface{}{"id": 1})
	token := domain.NewToken("t1", "node1", []*domain.Fact{leftFact})
	node.betaMemory.StoreToken(token)

	// Process right fact (now existence is satisfied)
	rightFact := domain.NewFact("f_right", "Person", map[string]interface{}{"age": 25})
	err := node.ProcessRightFact(rightFact)
	if err != nil {
		t.Errorf("ProcessRightFact should not return error: %v", err)
	}

	// Fact should be in right memory
	rightFacts := node.GetRightMemory()
	if len(rightFacts) != 1 {
		t.Errorf("Expected 1 fact in right memory, got %d", len(rightFacts))
	}

	// Token should now be propagated
	childLeftTokens := child.GetLeftMemory()
	if len(childLeftTokens) != 1 {
		t.Errorf("Expected token to be propagated after right fact, got %d tokens", len(childLeftTokens))
	}
}

func TestExistsNode_ConcurrentAccess(t *testing.T) {
	logger := newMockLogger()
	node := NewExistsNode("exists1", logger)

	variable := domain.TypedVariable{Name: "p", DataType: "Person"}
	condition := map[string]interface{}{"type": "exists"}

	// Concurrent set/get
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			node.SetExistenceCondition(variable, condition)
			node.GetExistenceCondition()
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify final state
	retVar, retCond := node.GetExistenceCondition()
	if retVar.Name != "p" {
		t.Error("Variable should be set after concurrent access")
	}
	if retCond == nil {
		t.Error("Condition should be set after concurrent access")
	}
}

// ===== AccumulateNodeImpl Tests =====

func TestNewAccumulateNode(t *testing.T) {
	logger := newMockLogger()
	accumulator := domain.AccumulateFunction{
		FunctionType: AggregateFunctionSum,
		Field:        "amount",
	}

	node := NewAccumulateNode("acc1", accumulator, logger)

	if node == nil {
		t.Fatal("NewAccumulateNode should not return nil")
	}
	if node.ID() != "acc1" {
		t.Errorf("Expected ID 'acc1', got '%s'", node.ID())
	}
	if node.Type() != "AccumulateNode" {
		t.Errorf("Expected Type 'AccumulateNode', got '%s'", node.Type())
	}
	if node.accumulatedValues == nil {
		t.Error("accumulatedValues should be initialized")
	}
}

func TestAccumulateNode_SetGetAccumulator(t *testing.T) {
	logger := newMockLogger()
	accumulator := domain.AccumulateFunction{
		FunctionType: AggregateFunctionSum,
		Field:        "amount",
	}

	node := NewAccumulateNode("acc1", accumulator, logger)

	// Get initial
	retrieved := node.GetAccumulator()
	if retrieved.FunctionType != AggregateFunctionSum {
		t.Errorf("Expected FunctionType Sum, got '%s'", retrieved.FunctionType)
	}

	// Set new accumulator
	newAcc := domain.AccumulateFunction{
		FunctionType: AggregateFunctionCount,
	}
	node.SetAccumulator(newAcc)

	// Verify change
	retrieved = node.GetAccumulator()
	if retrieved.FunctionType != AggregateFunctionCount {
		t.Errorf("Expected FunctionType Count, got '%s'", retrieved.FunctionType)
	}
}

func TestAccumulateNode_ComputeAggregate_Sum(t *testing.T) {
	logger := newMockLogger()
	accumulator := domain.AccumulateFunction{
		FunctionType: AggregateFunctionSum,
		Field:        "amount",
	}

	node := NewAccumulateNode("acc1", accumulator, logger)

	facts := []*domain.Fact{
		domain.NewFact("f1", "Order", map[string]interface{}{"amount": 100}),
		domain.NewFact("f2", "Order", map[string]interface{}{"amount": 200}),
		domain.NewFact("f3", "Order", map[string]interface{}{"amount": 50}),
	}

	token := domain.NewToken("t1", "node1", []*domain.Fact{})
	result, err := node.ComputeAggregate(token, facts)
	if err != nil {
		t.Errorf("ComputeAggregate should not return error: %v", err)
	}

	// Sum should be 350 (100 + 200 + 50)
	if result.(float64) != 350.0 {
		t.Errorf("Expected sum 350, got %v", result)
	}
}

func TestAccumulateNode_ComputeAggregate_Count(t *testing.T) {
	logger := newMockLogger()
	accumulator := domain.AccumulateFunction{
		FunctionType: AggregateFunctionCount,
	}

	node := NewAccumulateNode("acc1", accumulator, logger)

	facts := []*domain.Fact{
		domain.NewFact("f1", "Order", nil),
		domain.NewFact("f2", "Order", nil),
		domain.NewFact("f3", "Order", nil),
		domain.NewFact("f4", "Order", nil),
	}

	token := domain.NewToken("t1", "node1", []*domain.Fact{})
	result, err := node.ComputeAggregate(token, facts)
	if err != nil {
		t.Errorf("ComputeAggregate should not return error: %v", err)
	}

	if result != 4 {
		t.Errorf("Expected count 4, got %v", result)
	}
}

func TestAccumulateNode_ComputeAggregate_Avg(t *testing.T) {
	logger := newMockLogger()
	accumulator := domain.AccumulateFunction{
		FunctionType: AggregateFunctionAvg,
		Field:        "score",
	}

	node := NewAccumulateNode("acc1", accumulator, logger)

	facts := []*domain.Fact{
		domain.NewFact("f1", "Test", map[string]interface{}{"score": 80}),
		domain.NewFact("f2", "Test", map[string]interface{}{"score": 90}),
		domain.NewFact("f3", "Test", map[string]interface{}{"score": 70}),
	}

	token := domain.NewToken("t1", "node1", []*domain.Fact{})
	result, err := node.ComputeAggregate(token, facts)
	if err != nil {
		t.Errorf("ComputeAggregate should not return error: %v", err)
	}

	// Average should be 80 (240 / 3)
	if result.(float64) != 80.0 {
		t.Errorf("Expected average 80, got %v", result)
	}
}

func TestAccumulateNode_ComputeAggregate_Min(t *testing.T) {
	logger := newMockLogger()
	accumulator := domain.AccumulateFunction{
		FunctionType: AggregateFunctionMin,
		Field:        "price",
	}

	node := NewAccumulateNode("acc1", accumulator, logger)

	facts := []*domain.Fact{
		domain.NewFact("f1", "Product", map[string]interface{}{"price": 100}),
		domain.NewFact("f2", "Product", map[string]interface{}{"price": 50}),
		domain.NewFact("f3", "Product", map[string]interface{}{"price": 200}),
	}

	token := domain.NewToken("t1", "node1", []*domain.Fact{})
	result, err := node.ComputeAggregate(token, facts)
	if err != nil {
		t.Errorf("ComputeAggregate should not return error: %v", err)
	}

	if result.(float64) != 50.0 {
		t.Errorf("Expected min 50, got %v", result)
	}
}

func TestAccumulateNode_ComputeAggregate_Max(t *testing.T) {
	logger := newMockLogger()
	accumulator := domain.AccumulateFunction{
		FunctionType: AggregateFunctionMax,
		Field:        "price",
	}

	node := NewAccumulateNode("acc1", accumulator, logger)

	facts := []*domain.Fact{
		domain.NewFact("f1", "Product", map[string]interface{}{"price": 100}),
		domain.NewFact("f2", "Product", map[string]interface{}{"price": 50}),
		domain.NewFact("f3", "Product", map[string]interface{}{"price": 200}),
	}

	token := domain.NewToken("t1", "node1", []*domain.Fact{})
	result, err := node.ComputeAggregate(token, facts)
	if err != nil {
		t.Errorf("ComputeAggregate should not return error: %v", err)
	}

	if result.(float64) != 200.0 {
		t.Errorf("Expected max 200, got %v", result)
	}
}

func TestAccumulateNode_ComputeAggregate_NoFunction(t *testing.T) {
	logger := newMockLogger()
	accumulator := domain.AccumulateFunction{
		FunctionType: "",
	}

	node := NewAccumulateNode("acc1", accumulator, logger)

	facts := []*domain.Fact{
		domain.NewFact("f1", "Test", map[string]interface{}{"value": 1}),
	}

	token := domain.NewToken("t1", "node1", []*domain.Fact{})
	_, err := node.ComputeAggregate(token, facts)
	if err == nil {
		t.Error("ComputeAggregate should return error when no function defined")
	}
}

func TestAccumulateNode_ComputeAggregate_UnsupportedFunction(t *testing.T) {
	logger := newMockLogger()
	accumulator := domain.AccumulateFunction{
		FunctionType: "INVALID",
	}

	node := NewAccumulateNode("acc1", accumulator, logger)

	facts := []*domain.Fact{
		domain.NewFact("f1", "Test", map[string]interface{}{"value": 1}),
	}

	token := domain.NewToken("t1", "node1", []*domain.Fact{})
	_, err := node.ComputeAggregate(token, facts)
	if err == nil {
		t.Error("ComputeAggregate should return error for unsupported function")
	}
}

func TestAccumulateNode_ProcessLeftToken(t *testing.T) {
	logger := newMockLogger()
	accumulator := domain.AccumulateFunction{
		FunctionType: AggregateFunctionSum,
		Field:        "amount",
	}

	node := NewAccumulateNode("acc1", accumulator, logger)

	// Add child
	childLogger := newMockLogger()
	child := NewBaseBetaNode("child1", "BetaNode", childLogger)
	node.AddChild(child)

	// Add right facts
	node.betaMemory.StoreFact(domain.NewFact("f1", "Order", map[string]interface{}{"amount": 100}))
	node.betaMemory.StoreFact(domain.NewFact("f2", "Order", map[string]interface{}{"amount": 200}))

	// Process left token
	leftFact := domain.NewFact("f_left", "Customer", map[string]interface{}{"id": 1})
	token := domain.NewToken("t1", "node1", []*domain.Fact{leftFact})

	err := node.ProcessLeftToken(token)
	if err != nil {
		t.Errorf("ProcessLeftToken should not return error: %v", err)
	}

	// Token should be in left memory
	leftTokens := node.GetLeftMemory()
	if len(leftTokens) != 1 {
		t.Errorf("Expected 1 token in left memory, got %d", len(leftTokens))
	}

	// Accumulated value should be stored
	node.mu.RLock()
	storedValue := node.accumulatedValues["t1"]
	node.mu.RUnlock()

	if storedValue.(float64) != 300.0 {
		t.Errorf("Expected accumulated value 300, got %v", storedValue)
	}

	// Token should be propagated
	childLeftTokens := child.GetLeftMemory()
	if len(childLeftTokens) != 1 {
		t.Errorf("Expected token to be propagated, got %d tokens", len(childLeftTokens))
	}
}

func TestAccumulateNode_ProcessRightFact(t *testing.T) {
	logger := newMockLogger()
	accumulator := domain.AccumulateFunction{
		FunctionType: AggregateFunctionCount,
	}

	node := NewAccumulateNode("acc1", accumulator, logger)

	// Add child
	childLogger := newMockLogger()
	child := NewBaseBetaNode("child1", "BetaNode", childLogger)
	node.AddChild(child)

	// Add left token first
	leftFact := domain.NewFact("f_left", "Customer", map[string]interface{}{"id": 1})
	token := domain.NewToken("t1", "node1", []*domain.Fact{leftFact})
	node.betaMemory.StoreToken(token)

	// Initially count is 0
	node.mu.Lock()
	node.accumulatedValues["t1"] = 0
	node.mu.Unlock()

	// Process right fact
	rightFact := domain.NewFact("f_right", "Order", map[string]interface{}{"amount": 100})
	err := node.ProcessRightFact(rightFact)
	if err != nil {
		t.Errorf("ProcessRightFact should not return error: %v", err)
	}

	// Fact should be in right memory
	rightFacts := node.GetRightMemory()
	if len(rightFacts) != 1 {
		t.Errorf("Expected 1 fact in right memory, got %d", len(rightFacts))
	}

	// Accumulated value should be updated (count = 1)
	node.mu.RLock()
	storedValue := node.accumulatedValues["t1"]
	node.mu.RUnlock()

	if storedValue != 1 {
		t.Errorf("Expected accumulated value 1, got %v", storedValue)
	}

	// Token should be propagated with updated aggregate
	childLeftTokens := child.GetLeftMemory()
	if len(childLeftTokens) != 1 {
		t.Errorf("Expected token to be propagated, got %d tokens", len(childLeftTokens))
	}
}

func TestAccumulateNode_ConcurrentAccess(t *testing.T) {
	logger := newMockLogger()
	accumulator := domain.AccumulateFunction{
		FunctionType: AggregateFunctionSum,
		Field:        "value",
	}

	node := NewAccumulateNode("acc1", accumulator, logger)

	// Concurrent set/get of accumulator
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			newAcc := domain.AccumulateFunction{
				FunctionType: AggregateFunctionCount,
			}
			node.SetAccumulator(newAcc)
			node.GetAccumulator()
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify final state
	finalAcc := node.GetAccumulator()
	if finalAcc.FunctionType == "" {
		t.Error("Accumulator should be set after concurrent access")
	}
}

func TestAccumulateNode_ComputeAggregate_EmptyFacts(t *testing.T) {
	logger := newMockLogger()
	accumulator := domain.AccumulateFunction{
		FunctionType: AggregateFunctionSum,
		Field:        "amount",
	}

	node := NewAccumulateNode("acc1", accumulator, logger)

	token := domain.NewToken("t1", "node1", []*domain.Fact{})
	result, err := node.ComputeAggregate(token, []*domain.Fact{})
	if err != nil {
		t.Errorf("ComputeAggregate should not return error for empty facts: %v", err)
	}

	// Sum of empty should be 0
	if result.(float64) != 0.0 {
		t.Errorf("Expected sum 0 for empty facts, got %v", result)
	}
}

func TestAccumulateNode_ComputeAggregate_MixedTypes(t *testing.T) {
	logger := newMockLogger()
	accumulator := domain.AccumulateFunction{
		FunctionType: AggregateFunctionSum,
		Field:        "value",
	}

	node := NewAccumulateNode("acc1", accumulator, logger)

	// Mix of int and float
	facts := []*domain.Fact{
		domain.NewFact("f1", "Test", map[string]interface{}{"value": 10}),
		domain.NewFact("f2", "Test", map[string]interface{}{"value": 20.5}),
		domain.NewFact("f3", "Test", map[string]interface{}{"value": 15}),
	}

	token := domain.NewToken("t1", "node1", []*domain.Fact{})
	result, err := node.ComputeAggregate(token, facts)
	if err != nil {
		t.Errorf("ComputeAggregate should handle mixed numeric types: %v", err)
	}

	// Should handle mixed types gracefully
	if result == nil {
		t.Error("Result should not be nil for mixed numeric types")
	}
}
