// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"testing"
)

// MockNode est un nœud de test qui capture les tokens propagés
type mockTerminalNode struct {
	BaseNode
	onActivateLeft func(*Token) error
}

func (m *mockTerminalNode) ActivateLeft(token *Token) error {
	if m.onActivateLeft != nil {
		return m.onActivateLeft(token)
	}
	m.Memory.AddToken(token)
	return nil
}

func (m *mockTerminalNode) ActivateRight(fact *Fact) error {
	return nil
}

func (m *mockTerminalNode) ActivateRetract(factID string) error {
	return nil
}

// TestJoinCascade_2Variables_UserOrder teste la cascade de jointure à 2 variables (régression)
func TestJoinCascade_2Variables_UserOrder(t *testing.T) {
	t.Log("🧪 TEST Cascade 2 variables - User-Order (régression)")
	t.Log("======================================================")

	userFact := &Fact{
		ID:     "u1",
		Type:   "User",
		Fields: map[string]interface{}{"id": 1, "name": "Alice"},
	}

	orderFact := &Fact{
		ID:     "o1",
		Type:   "Order",
		Fields: map[string]interface{}{"id": 100, "user_id": 1},
	}

	joinNode := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join_user_order",
			Children: []Node{},
			Memory:   &WorkingMemory{NodeID: "join_user_order", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		},
		LeftVariables:  []string{"user"},
		RightVariables: []string{"order"},
		AllVariables:   []string{"user", "order"},
		VariableTypes: map[string]string{
			"user":  "User",
			"order": "Order",
		},
		LeftMemory:     &WorkingMemory{NodeID: "join_user_order_left", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		RightMemory:    &WorkingMemory{NodeID: "join_user_order_right", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		ResultMemory:   &WorkingMemory{NodeID: "join_user_order_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		JoinConditions: nil,
	}

	var finalToken *Token
	mockTerminal := &mockTerminalNode{
		BaseNode: BaseNode{
			ID:     "terminal",
			Memory: &WorkingMemory{NodeID: "terminal", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		},
		onActivateLeft: func(token *Token) error {
			finalToken = token
			return nil
		},
	}
	joinNode.AddChild(mockTerminal)

	userToken := NewTokenWithFact(userFact, "user", "type_node_user")
	err := joinNode.ActivateLeft(userToken)
	if err != nil {
		t.Fatalf("❌ ActivateLeft erreur: %v", err)
	}

	err = joinNode.ActivateRight(orderFact)
	if err != nil {
		t.Fatalf("❌ ActivateRight erreur: %v", err)
	}

	if finalToken == nil {
		t.Fatal("❌ Aucun token propagé au terminal")
	}

	if finalToken.Bindings.Len() != 2 {
		t.Errorf("❌ Attendu 2 bindings, got %d", finalToken.Bindings.Len())
	}

	if !finalToken.HasBinding("user") {
		t.Errorf("❌ Variable 'user' manquante")
	}

	if !finalToken.HasBinding("order") {
		t.Errorf("❌ Variable 'order' manquante")
	}

	if finalToken.GetBinding("user") != userFact {
		t.Errorf("❌ Binding 'user' incorrect")
	}

	if finalToken.GetBinding("order") != orderFact {
		t.Errorf("❌ Binding 'order' incorrect")
	}

	t.Log("✅ Cascade 2 variables fonctionne (régression OK)")
}

// TestJoinCascade_3Variables_UserOrderProduct teste la cascade de jointure à 3 variables
func TestJoinCascade_3Variables_UserOrderProduct(t *testing.T) {
	t.Log("🧪 TEST Cascade 3 variables - User-Order-Product")
	t.Log("=================================================")

	userFact := &Fact{
		ID:     "u1",
		Type:   "User",
		Fields: map[string]interface{}{"id": 1, "name": "Alice"},
	}

	orderFact := &Fact{
		ID:     "o1",
		Type:   "Order",
		Fields: map[string]interface{}{"id": 100, "user_id": 1, "product_id": 200},
	}

	productFact := &Fact{
		ID:     "p1",
		Type:   "Product",
		Fields: map[string]interface{}{"id": 200, "name": "Laptop"},
	}

	joinNode1 := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join1_user_order",
			Children: []Node{},
			Memory:   &WorkingMemory{NodeID: "join1_user_order", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		},
		LeftVariables:  []string{"user"},
		RightVariables: []string{"order"},
		AllVariables:   []string{"user", "order"},
		VariableTypes: map[string]string{
			"user":    "User",
			"order":   "Order",
			"product": "Product",
		},
		LeftMemory:     &WorkingMemory{NodeID: "join1_left", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		RightMemory:    &WorkingMemory{NodeID: "join1_right", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		ResultMemory:   &WorkingMemory{NodeID: "join1_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		JoinConditions: nil,
	}

	joinNode2 := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join2_add_product",
			Children: []Node{},
			Memory:   &WorkingMemory{NodeID: "join2_add_product", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		},
		LeftVariables:  []string{"user", "order"},
		RightVariables: []string{"product"},
		AllVariables:   []string{"user", "order", "product"},
		VariableTypes: map[string]string{
			"user":    "User",
			"order":   "Order",
			"product": "Product",
		},
		LeftMemory:     &WorkingMemory{NodeID: "join2_left", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		RightMemory:    &WorkingMemory{NodeID: "join2_right", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		ResultMemory:   &WorkingMemory{NodeID: "join2_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		JoinConditions: nil,
	}

	joinNode1.AddChild(joinNode2)

	var finalToken *Token
	mockTerminal := &mockTerminalNode{
		BaseNode: BaseNode{
			ID:     "terminal",
			Memory: &WorkingMemory{NodeID: "terminal", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		},
		onActivateLeft: func(token *Token) error {
			finalToken = token
			return nil
		},
	}
	joinNode2.AddChild(mockTerminal)

	userToken := NewTokenWithFact(userFact, "user", "type_user")
	err := joinNode1.ActivateLeft(userToken)
	if err != nil {
		t.Fatalf("❌ JoinNode1 ActivateLeft erreur: %v", err)
	}

	err = joinNode1.ActivateRight(orderFact)
	if err != nil {
		t.Fatalf("❌ JoinNode1 ActivateRight erreur: %v", err)
	}

	err = joinNode2.ActivateRight(productFact)
	if err != nil {
		t.Fatalf("❌ JoinNode2 ActivateRight erreur: %v", err)
	}

	if finalToken == nil {
		t.Fatal("❌ Aucun token propagé au terminal")
	}

	if finalToken.Bindings.Len() != 3 {
		t.Errorf("❌ CRITIQUE: Attendu 3 bindings, got %d", finalToken.Bindings.Len())
		t.Errorf("   Variables présentes: %v", finalToken.GetVariables())
	}

	expectedVars := []string{"user", "order", "product"}
	for _, v := range expectedVars {
		if !finalToken.HasBinding(v) {
			t.Errorf("❌ Variable '%s' manquante dans le token final", v)
		}
	}

	if finalToken.GetBinding("user") != userFact {
		t.Errorf("❌ Binding 'user' incorrect")
	}
	if finalToken.GetBinding("order") != orderFact {
		t.Errorf("❌ Binding 'order' incorrect")
	}
	if finalToken.GetBinding("product") != productFact {
		t.Errorf("❌ Binding 'product' incorrect")
	}

	t.Log("✅ Cascade 3 variables fonctionne - TOUS les bindings préservés")
}

// TestJoinCascade_3Variables_DifferentOrders teste différents ordres de soumission
func TestJoinCascade_3Variables_DifferentOrders(t *testing.T) {
	t.Log("🧪 TEST Cascade 3 variables - Différents ordres de soumission")
	t.Log("==============================================================")

	orders := []struct {
		name  string
		order []string
	}{
		{"User→Order→Product", []string{"user", "order", "product"}},
		{"User→Product→Order", []string{"user", "product", "order"}},
		{"Order→User→Product", []string{"order", "user", "product"}},
		{"Product→User→Order", []string{"product", "user", "order"}},
	}

	for _, tc := range orders {
		t.Run(tc.name, func(t *testing.T) {
			facts := map[string]*Fact{
				"user":    {ID: "u1", Type: "User", Fields: map[string]interface{}{"id": 1}},
				"order":   {ID: "o1", Type: "Order", Fields: map[string]interface{}{"id": 100}},
				"product": {ID: "p1", Type: "Product", Fields: map[string]interface{}{"id": 200}},
			}

			joinNode1, joinNode2, mockTerminal := setupCascade3Variables()

			var finalToken *Token
			mockTerminal.onActivateLeft = func(token *Token) error {
				finalToken = token
				return nil
			}

			for _, factName := range tc.order {
				fact := facts[factName]
				if factName == "user" {
					token := NewTokenWithFact(fact, "user", "type_user")
					joinNode1.ActivateLeft(token)
				} else if factName == "order" {
					joinNode1.ActivateRight(fact)
				} else if factName == "product" {
					joinNode2.ActivateRight(fact)
				}
			}

			if finalToken == nil {
				t.Errorf("❌ Aucun token final pour ordre %v", tc.order)
				return
			}

			if finalToken.Bindings.Len() != 3 {
				t.Errorf("❌ Ordre %v: attendu 3 bindings, got %d", tc.order, finalToken.Bindings.Len())
			}

			t.Logf("✅ Ordre %v: %d bindings", tc.name, finalToken.Bindings.Len())
		})
	}

	t.Log("✅ Résultats cohérents quel que soit l'ordre de soumission")
}

// TestJoinCascade_NVariables teste la scalabilité jusqu'à N=10 variables
func TestJoinCascade_NVariables(t *testing.T) {
	t.Log("🧪 TEST Cascade N variables - Scalabilité")
	t.Log("==========================================")

	for n := 2; n <= 10; n++ {
		t.Run(fmt.Sprintf("n=%d_variables", n), func(t *testing.T) {
			facts := make([]*Fact, n)
			varNames := make([]string, n)

			for i := 0; i < n; i++ {
				varNames[i] = fmt.Sprintf("var%d", i)
				facts[i] = &Fact{
					ID:     fmt.Sprintf("f%d", i),
					Type:   fmt.Sprintf("Type%d", i),
					Fields: map[string]interface{}{"id": i},
				}
			}

			joinNodes := buildCascade(n, varNames)

			var finalToken *Token
			mockTerminal := &mockTerminalNode{
				BaseNode: BaseNode{
					ID:     "terminal",
					Memory: &WorkingMemory{NodeID: "terminal", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
				},
				onActivateLeft: func(token *Token) error {
					finalToken = token
					return nil
				},
			}
			lastJoinNode := joinNodes[len(joinNodes)-1]
			lastJoinNode.AddChild(mockTerminal)

			for i, fact := range facts {
				if i == 0 {
					token := NewTokenWithFact(fact, varNames[i], "type_node")
					joinNodes[0].ActivateLeft(token)
				} else if i == 1 {
					joinNodes[0].ActivateRight(fact)
				} else {
					joinNodes[i-1].ActivateRight(fact)
				}
			}

			if finalToken == nil {
				t.Fatalf("❌ N=%d: Aucun token final", n)
			}

			if finalToken.Bindings.Len() != n {
				t.Errorf("❌ N=%d: Attendu %d bindings, got %d", n, n, finalToken.Bindings.Len())
				t.Errorf("   Variables présentes: %v", finalToken.GetVariables())
			}

			for _, varName := range varNames {
				if !finalToken.HasBinding(varName) {
					t.Errorf("❌ N=%d: Variable '%s' manquante", n, varName)
				}
			}

			t.Logf("✅ N=%d variables: Tous les bindings préservés", n)
		})
	}

	t.Log("✅ Scalabilité validée jusqu'à N=10 variables")
}

// setupCascade3Variables crée une cascade de 2 JoinNodes pour 3 variables
func setupCascade3Variables() (*JoinNode, *JoinNode, *mockTerminalNode) {
	joinNode1 := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join1",
			Children: []Node{},
			Memory:   &WorkingMemory{NodeID: "join1", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		},
		LeftVariables:  []string{"user"},
		RightVariables: []string{"order"},
		AllVariables:   []string{"user", "order"},
		VariableTypes: map[string]string{
			"user":    "User",
			"order":   "Order",
			"product": "Product",
		},
		LeftMemory:     &WorkingMemory{NodeID: "join1_left", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		RightMemory:    &WorkingMemory{NodeID: "join1_right", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		ResultMemory:   &WorkingMemory{NodeID: "join1_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		JoinConditions: nil,
	}

	joinNode2 := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join2",
			Children: []Node{},
			Memory:   &WorkingMemory{NodeID: "join2", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		},
		LeftVariables:  []string{"user", "order"},
		RightVariables: []string{"product"},
		AllVariables:   []string{"user", "order", "product"},
		VariableTypes: map[string]string{
			"user":    "User",
			"order":   "Order",
			"product": "Product",
		},
		LeftMemory:     &WorkingMemory{NodeID: "join2_left", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		RightMemory:    &WorkingMemory{NodeID: "join2_right", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		ResultMemory:   &WorkingMemory{NodeID: "join2_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		JoinConditions: nil,
	}

	joinNode1.AddChild(joinNode2)

	mockTerminal := &mockTerminalNode{
		BaseNode: BaseNode{
			ID:     "terminal",
			Memory: &WorkingMemory{NodeID: "terminal", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		},
	}
	joinNode2.AddChild(mockTerminal)

	return joinNode1, joinNode2, mockTerminal
}

// buildCascade construit une cascade de (n-1) JoinNodes pour n variables
func buildCascade(n int, varNames []string) []*JoinNode {
	if n < 2 {
		return []*JoinNode{}
	}

	joinNodes := make([]*JoinNode, n-1)
	varTypes := make(map[string]string)
	for i, v := range varNames {
		varTypes[v] = fmt.Sprintf("Type%d", i)
	}

	joinNodes[0] = &JoinNode{
		BaseNode: BaseNode{
			ID:       fmt.Sprintf("join_%d", 0),
			Children: []Node{},
			Memory:   &WorkingMemory{NodeID: fmt.Sprintf("join_%d", 0), Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		},
		LeftVariables:  []string{varNames[0]},
		RightVariables: []string{varNames[1]},
		AllVariables:   []string{varNames[0], varNames[1]},
		VariableTypes:  varTypes,
		LeftMemory:     &WorkingMemory{NodeID: fmt.Sprintf("join_%d_left", 0), Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		RightMemory:    &WorkingMemory{NodeID: fmt.Sprintf("join_%d_right", 0), Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		ResultMemory:   &WorkingMemory{NodeID: fmt.Sprintf("join_%d_result", 0), Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
	}

	for i := 2; i < n; i++ {
		leftVars := make([]string, i)
		copy(leftVars, varNames[0:i])

		allVars := make([]string, i+1)
		copy(allVars, varNames[0:i+1])

		joinNodes[i-1] = &JoinNode{
			BaseNode: BaseNode{
				ID:       fmt.Sprintf("join_%d", i-1),
				Children: []Node{},
				Memory:   &WorkingMemory{NodeID: fmt.Sprintf("join_%d", i-1), Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			},
			LeftVariables:  leftVars,
			RightVariables: []string{varNames[i]},
			AllVariables:   allVars,
			VariableTypes:  varTypes,
			LeftMemory:     &WorkingMemory{NodeID: fmt.Sprintf("join_%d_left", i-1), Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			RightMemory:    &WorkingMemory{NodeID: fmt.Sprintf("join_%d_right", i-1), Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			ResultMemory:   &WorkingMemory{NodeID: fmt.Sprintf("join_%d_result", i-1), Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		}

		joinNodes[i-2].AddChild(joinNodes[i-1])
	}

	return joinNodes
}
