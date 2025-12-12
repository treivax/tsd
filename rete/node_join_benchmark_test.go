// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"testing"
)

// ============================================================================
// BENCHMARKS DE BINDINGCHAIN (COMPLÉMENTAIRES)
// Note: Les benchmarks basiques Add/Get/Variables/ToMap sont dans binding_chain_test.go
// ============================================================================

// BenchmarkBindingChain_Add_10Variables mesure la performance d'ajout de 10 bindings
func BenchmarkBindingChain_Add_10Variables(b *testing.B) {
	facts := make([]*Fact, 10)
	for i := 0; i < 10; i++ {
		facts[i] = &Fact{
			ID:     fmt.Sprintf("f%d", i),
			Type:   "Type",
			Fields: map[string]interface{}{"id": i},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chain := NewBindingChain()
		for j := 0; j < 10; j++ {
			chain = chain.Add(fmt.Sprintf("var%d", j), facts[j])
		}
	}
}

// BenchmarkBindingChain_Get_SmallChain mesure la performance de Get() sur une petite chaîne (3 bindings)
func BenchmarkBindingChain_Get_SmallChain(b *testing.B) {
	chain := NewBindingChain()
	for i := 0; i < 3; i++ {
		chain = chain.Add(
			fmt.Sprintf("var%d", i),
			&Fact{ID: fmt.Sprintf("f%d", i), Type: "Type", Fields: map[string]interface{}{"id": i}},
		)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chain.Get("var1")
	}
}

// BenchmarkBindingChain_Get_MediumChain mesure la performance de Get() sur une chaîne moyenne (10 bindings)
func BenchmarkBindingChain_Get_MediumChain(b *testing.B) {
	chain := NewBindingChain()
	for i := 0; i < 10; i++ {
		chain = chain.Add(
			fmt.Sprintf("var%d", i),
			&Fact{ID: fmt.Sprintf("f%d", i), Type: "Type", Fields: map[string]interface{}{"id": i}},
		)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chain.Get("var0") // Pire cas : chercher le plus ancien
	}
}

// BenchmarkBindingChain_Merge mesure la performance de fusion de deux chaînes
func BenchmarkBindingChain_Merge(b *testing.B) {
	chain1 := NewBindingChain()
	chain2 := NewBindingChain()

	for i := 0; i < 5; i++ {
		chain1 = chain1.Add(
			fmt.Sprintf("v1_%d", i),
			&Fact{ID: fmt.Sprintf("f1_%d", i), Type: "Type", Fields: map[string]interface{}{"id": i}},
		)
		chain2 = chain2.Add(
			fmt.Sprintf("v2_%d", i),
			&Fact{ID: fmt.Sprintf("f2_%d", i), Type: "Type", Fields: map[string]interface{}{"id": i}},
		)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chain1.Merge(chain2)
	}
}

// BenchmarkBindingChain_Len mesure la performance de Len()
func BenchmarkBindingChain_Len(b *testing.B) {
	chain := NewBindingChain()
	for i := 0; i < 10; i++ {
		chain = chain.Add(
			fmt.Sprintf("var%d", i),
			&Fact{ID: fmt.Sprintf("f%d", i), Type: "Type", Fields: map[string]interface{}{"id": i}},
		)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chain.Len()
	}
}

// ============================================================================
// BENCHMARKS DE JOINNODE
// ============================================================================

// BenchmarkJoinNode_2Variables mesure la performance de jointure avec 2 variables
func BenchmarkJoinNode_2Variables(b *testing.B) {
	userFact := &Fact{
		ID:     "u1",
		Type:   "User",
		Fields: map[string]interface{}{"id": 1},
	}
	orderFact := &Fact{
		ID:     "o1",
		Type:   "Order",
		Fields: map[string]interface{}{"user_id": 1},
	}

	joinNode := NewJoinNode(
		"join_bench",
		nil,
		[]string{"user"},
		[]string{"order"},
		map[string]string{
			"user":  "User",
			"order": "Order",
		},
		nil, // Storage non nécessaire pour le benchmark
	)

	// Ajouter manuellement les conditions de jointure
	joinNode.JoinConditions = []JoinCondition{
		{
			LeftField:  "id",
			RightField: "user_id",
			LeftVar:    "user",
			RightVar:   "order",
			Operator:   "==",
		},
	}

	leftToken := NewTokenWithFact(userFact, "user", "test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		joinNode.LeftMemory = &WorkingMemory{NodeID: "join_bench_left", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode.RightMemory = &WorkingMemory{NodeID: "join_bench_right", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode.ResultMemory = &WorkingMemory{NodeID: "join_bench_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}

		_ = joinNode.ActivateLeft(leftToken)
		_ = joinNode.ActivateRight(orderFact)
	}
}

// BenchmarkJoinNode_3Variables mesure la performance de jointure en cascade avec 3 variables
func BenchmarkJoinNode_3Variables(b *testing.B) {
	userFact := &Fact{
		ID:     "u1",
		Type:   "User",
		Fields: map[string]interface{}{"id": 1},
	}
	orderFact := &Fact{
		ID:     "o1",
		Type:   "Order",
		Fields: map[string]interface{}{"user_id": 1, "id": 100},
	}
	productFact := &Fact{
		ID:     "p1",
		Type:   "Product",
		Fields: map[string]interface{}{"order_id": 100},
	}

	// Premier JoinNode: user + order
	joinNode1 := NewJoinNode(
		"join1",
		nil,
		[]string{"user"},
		[]string{"order"},
		map[string]string{
			"user":    "User",
			"order":   "Order",
			"product": "Product",
		},
		nil,
	)
	joinNode1.JoinConditions = []JoinCondition{
		{
			LeftField:  "id",
			RightField: "user_id",
			LeftVar:    "user",
			RightVar:   "order",
			Operator:   "==",
		},
	}

	// Deuxième JoinNode: (user + order) + product
	joinNode2 := NewJoinNode(
		"join2",
		nil,
		[]string{"user", "order"},
		[]string{"product"},
		map[string]string{
			"user":    "User",
			"order":   "Order",
			"product": "Product",
		},
		nil,
	)
	joinNode2.JoinConditions = []JoinCondition{
		{
			LeftField:  "id",
			RightField: "order_id",
			LeftVar:    "order",
			RightVar:   "product",
			Operator:   "==",
		},
	}

	joinNode1.AddChild(joinNode2)

	userToken := NewTokenWithFact(userFact, "user", "test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		joinNode1.LeftMemory = &WorkingMemory{NodeID: "join1_left", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode1.RightMemory = &WorkingMemory{NodeID: "join1_right", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode1.ResultMemory = &WorkingMemory{NodeID: "join1_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode2.LeftMemory = &WorkingMemory{NodeID: "join2_left", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode2.RightMemory = &WorkingMemory{NodeID: "join2_right", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode2.ResultMemory = &WorkingMemory{NodeID: "join2_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}

		_ = joinNode1.ActivateLeft(userToken)
		_ = joinNode1.ActivateRight(orderFact)
		_ = joinNode2.ActivateRight(productFact)
	}
}

// BenchmarkJoinNode_4Variables mesure la performance de jointure en cascade avec 4 variables
func BenchmarkJoinNode_4Variables(b *testing.B) {
	userFact := &Fact{
		ID:     "u1",
		Type:   "User",
		Fields: map[string]interface{}{"id": 1},
	}
	orderFact := &Fact{
		ID:     "o1",
		Type:   "Order",
		Fields: map[string]interface{}{"user_id": 1, "id": 100},
	}
	productFact := &Fact{
		ID:     "p1",
		Type:   "Product",
		Fields: map[string]interface{}{"order_id": 100, "id": 200},
	}
	paymentFact := &Fact{
		ID:     "pay1",
		Type:   "Payment",
		Fields: map[string]interface{}{"product_id": 200},
	}

	varTypes := map[string]string{
		"user":    "User",
		"order":   "Order",
		"product": "Product",
		"payment": "Payment",
	}

	// JoinNode 1: user + order
	joinNode1 := NewJoinNode("join1", nil, []string{"user"}, []string{"order"}, varTypes, nil)
	joinNode1.JoinConditions = []JoinCondition{
		{LeftField: "id", RightField: "user_id", LeftVar: "user", RightVar: "order", Operator: "=="},
	}

	// JoinNode 2: (user + order) + product
	joinNode2 := NewJoinNode("join2", nil, []string{"user", "order"}, []string{"product"}, varTypes, nil)
	joinNode2.JoinConditions = []JoinCondition{
		{LeftField: "id", RightField: "order_id", LeftVar: "order", RightVar: "product", Operator: "=="},
	}

	// JoinNode 3: (user + order + product) + payment
	joinNode3 := NewJoinNode("join3", nil, []string{"user", "order", "product"}, []string{"payment"}, varTypes, nil)
	joinNode3.JoinConditions = []JoinCondition{
		{LeftField: "id", RightField: "product_id", LeftVar: "product", RightVar: "payment", Operator: "=="},
	}

	joinNode1.AddChild(joinNode2)
	joinNode2.AddChild(joinNode3)

	userToken := NewTokenWithFact(userFact, "user", "test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		joinNode1.LeftMemory = &WorkingMemory{NodeID: "join1_left", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode1.RightMemory = &WorkingMemory{NodeID: "join1_right", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode1.ResultMemory = &WorkingMemory{NodeID: "join1_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode2.LeftMemory = &WorkingMemory{NodeID: "join2_left", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode2.RightMemory = &WorkingMemory{NodeID: "join2_right", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode2.ResultMemory = &WorkingMemory{NodeID: "join2_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode3.LeftMemory = &WorkingMemory{NodeID: "join3_left", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode3.RightMemory = &WorkingMemory{NodeID: "join3_right", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode3.ResultMemory = &WorkingMemory{NodeID: "join3_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}

		_ = joinNode1.ActivateLeft(userToken)
		_ = joinNode1.ActivateRight(orderFact)
		_ = joinNode2.ActivateRight(productFact)
		_ = joinNode3.ActivateRight(paymentFact)
	}
}

// BenchmarkJoinNode_PerformJoinWithTokens mesure la performance de performJoinWithTokens()
func BenchmarkJoinNode_PerformJoinWithTokens(b *testing.B) {
	userFact := &Fact{ID: "u1", Type: "User", Fields: map[string]interface{}{"id": 1}}
	orderFact := &Fact{ID: "o1", Type: "Order", Fields: map[string]interface{}{"user_id": 1}}

	token1 := &Token{
		ID:       "t1",
		Facts:    []*Fact{userFact},
		Bindings: NewBindingChainWith("user", userFact),
		NodeID:   "test",
		Metadata: TokenMetadata{JoinLevel: 0},
	}

	token2 := &Token{
		ID:       "t2",
		Facts:    []*Fact{orderFact},
		Bindings: NewBindingChainWith("order", orderFact),
		NodeID:   "test",
		Metadata: TokenMetadata{JoinLevel: 0},
	}

	joinNode := NewJoinNode(
		"join_bench",
		nil,
		[]string{"user"},
		[]string{"order"},
		map[string]string{
			"user":  "User",
			"order": "Order",
		},
		nil,
	)
	joinNode.JoinConditions = []JoinCondition{
		{
			LeftField:  "id",
			RightField: "user_id",
			LeftVar:    "user",
			RightVar:   "order",
			Operator:   "==",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = joinNode.performJoinWithTokens(token1, token2)
	}
}

// ============================================================================
// BENCHMARKS DE MÉMOIRE
// ============================================================================

// BenchmarkBindingChain_Memory mesure les allocations mémoire de BindingChain
func BenchmarkBindingChain_Memory(b *testing.B) {
	b.ReportAllocs()

	fact := &Fact{ID: "f1", Type: "User", Fields: map[string]interface{}{"id": 1}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chain := NewBindingChain()
		for j := 0; j < 10; j++ {
			chain = chain.Add(fmt.Sprintf("var%d", j), fact)
		}
	}
}

// BenchmarkJoinNode_Memory_3Variables mesure les allocations mémoire pour jointure 3 variables
func BenchmarkJoinNode_Memory_3Variables(b *testing.B) {
	b.ReportAllocs()

	userFact := &Fact{
		ID:     "u1",
		Type:   "User",
		Fields: map[string]interface{}{"id": 1},
	}
	orderFact := &Fact{
		ID:     "o1",
		Type:   "Order",
		Fields: map[string]interface{}{"user_id": 1, "id": 100},
	}
	productFact := &Fact{
		ID:     "p1",
		Type:   "Product",
		Fields: map[string]interface{}{"order_id": 100},
	}

	varTypes := map[string]string{
		"user":    "User",
		"order":   "Order",
		"product": "Product",
	}

	joinNode1 := NewJoinNode("join1", nil, []string{"user"}, []string{"order"}, varTypes, nil)
	joinNode1.JoinConditions = []JoinCondition{
		{LeftField: "id", RightField: "user_id", LeftVar: "user", RightVar: "order", Operator: "=="},
	}

	joinNode2 := NewJoinNode("join2", nil, []string{"user", "order"}, []string{"product"}, varTypes, nil)
	joinNode2.JoinConditions = []JoinCondition{
		{LeftField: "id", RightField: "order_id", LeftVar: "order", RightVar: "product", Operator: "=="},
	}

	joinNode1.AddChild(joinNode2)

	userToken := NewTokenWithFact(userFact, "user", "test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		joinNode1.LeftMemory = &WorkingMemory{NodeID: "join1_left", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode1.RightMemory = &WorkingMemory{NodeID: "join1_right", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode1.ResultMemory = &WorkingMemory{NodeID: "join1_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode2.LeftMemory = &WorkingMemory{NodeID: "join2_left", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode2.RightMemory = &WorkingMemory{NodeID: "join2_right", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}
		joinNode2.ResultMemory = &WorkingMemory{NodeID: "join2_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)}

		_ = joinNode1.ActivateLeft(userToken)
		_ = joinNode1.ActivateRight(orderFact)
		_ = joinNode2.ActivateRight(productFact)
	}
}

// ============================================================================
// BENCHMARKS COMPARATIFS
// ============================================================================

// BenchmarkBindingChain_vs_Map_Get compare BindingChain.Get() avec map[string]*Fact
func BenchmarkBindingChain_vs_Map_Get(b *testing.B) {
	fact := &Fact{ID: "f1", Type: "User", Fields: map[string]interface{}{"id": 1}}

	// Setup BindingChain
	chain := NewBindingChain()
	for i := 0; i < 10; i++ {
		chain = chain.Add(fmt.Sprintf("var%d", i), fact)
	}

	// Setup map
	m := make(map[string]*Fact)
	for i := 0; i < 10; i++ {
		m[fmt.Sprintf("var%d", i)] = fact
	}

	b.Run("BindingChain", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = chain.Get("var0")
		}
	})

	b.Run("Map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = m["var0"]
		}
	})
}
