// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// TestJoin_WithPKSimpleIDs teste les jointures avec IDs basés sur PK simple
func TestJoin_WithPKSimpleIDs(t *testing.T) {
	t.Log("🧪 TEST: Join - IDs basés sur PK simple")
	t.Log("========================================")

	personFact := &Fact{
		ID:   "Person~Alice",
		Type: "Person",
		Fields: map[string]interface{}{
			"nom": "Alice",
			"age": 30,
		},
	}

	membershipFact := &Fact{
		ID:   "Membership~Alice_Club1",
		Type: "Membership",
		Fields: map[string]interface{}{
			"person_id": "Person~Alice",
			"club":      "Club1",
		},
	}

	// Créer un JoinNode simple
	joinNode := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join_person_membership",
			Children: []Node{},
			Memory: &WorkingMemory{
				NodeID: "join_person_membership",
				Facts:  make(map[string]*Fact),
				Tokens: make(map[string]*Token),
			},
		},
		LeftVariables:  []string{"person"},
		RightVariables: []string{"membership"},
		AllVariables:   []string{"person", "membership"},
		VariableTypes: map[string]string{
			"person":     "Person",
			"membership": "Membership",
		},
		LeftMemory: &WorkingMemory{
			NodeID: "join_person_membership_left",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		RightMemory: &WorkingMemory{
			NodeID: "join_person_membership_right",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		ResultMemory: &WorkingMemory{
			NodeID: "join_person_membership_result",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		JoinConditions: nil, // Pas de condition pour ce test simple
	}

	var capturedToken *Token
	mockTerminal := &mockTerminalNode{
		BaseNode: BaseNode{
			ID: "terminal",
			Memory: &WorkingMemory{
				NodeID: "terminal",
				Facts:  make(map[string]*Fact),
				Tokens: make(map[string]*Token),
			},
		},
		onActivateLeft: func(token *Token) error {
			capturedToken = token
			return nil
		},
	}
	joinNode.AddChild(mockTerminal)

	// Activer left avec person token
	personToken := NewTokenWithFact(personFact, "person", "type_node_person")
	err := joinNode.ActivateLeft(personToken)
	if err != nil {
		t.Fatalf("❌ ActivateLeft erreur: %v", err)
	}

	// Activer right avec membership fact
	err = joinNode.ActivateRight(membershipFact)
	if err != nil {
		t.Fatalf("❌ ActivateRight erreur: %v", err)
	}

	// Vérifier qu'un token a été propagé
	if capturedToken == nil {
		t.Fatal("❌ Aucun token propagé au nœud terminal")
	}

	// Vérifier les bindings
	if !capturedToken.HasBinding("person") {
		t.Errorf("❌ Variable 'person' manquante dans le token")
	}

	if !capturedToken.HasBinding("membership") {
		t.Errorf("❌ Variable 'membership' manquante dans le token")
	}

	// Vérifier que les IDs sont corrects
	personBinding := capturedToken.GetBinding("person")
	if personBinding == nil || personBinding.ID != "Person~Alice" {
		t.Errorf("❌ ID person attendu 'Person~Alice', reçu '%v'", personBinding)
	}

	membershipBinding := capturedToken.GetBinding("membership")
	if membershipBinding == nil || membershipBinding.ID != "Membership~Alice_Club1" {
		t.Errorf("❌ ID membership attendu 'Membership~Alice_Club1', reçu '%v'", membershipBinding)
	}

	t.Log("✅ Test réussi: Join avec IDs PK simple")
}

// TestJoin_WithPKCompositeIDs teste les jointures avec IDs basés sur PK composite
func TestJoin_WithPKCompositeIDs(t *testing.T) {
	t.Log("🧪 TEST: Join - IDs basés sur PK composite")
	t.Log("==========================================")

	personFact := &Fact{
		ID:   "Person~Alice_Dupont",
		Type: "Person",
		Fields: map[string]interface{}{
			"prenom": "Alice",
			"nom":    "Dupont",
			"age":    30,
		},
	}

	contactFact := &Fact{
		ID:   "Contact~Alice_Dupont",
		Type: "Contact",
		Fields: map[string]interface{}{
			"person_id": "Person~Alice_Dupont",
			"email":     "alice.dupont@example.com",
		},
	}

	joinNode := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join_person_contact",
			Children: []Node{},
			Memory: &WorkingMemory{
				NodeID: "join_person_contact",
				Facts:  make(map[string]*Fact),
				Tokens: make(map[string]*Token),
			},
		},
		LeftVariables:  []string{"person"},
		RightVariables: []string{"contact"},
		AllVariables:   []string{"person", "contact"},
		VariableTypes: map[string]string{
			"person":  "Person",
			"contact": "Contact",
		},
		LeftMemory: &WorkingMemory{
			NodeID: "join_person_contact_left",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		RightMemory: &WorkingMemory{
			NodeID: "join_person_contact_right",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		ResultMemory: &WorkingMemory{
			NodeID: "join_person_contact_result",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		JoinConditions: nil,
	}

	var capturedToken *Token
	mockTerminal := &mockTerminalNode{
		BaseNode: BaseNode{
			ID: "terminal",
			Memory: &WorkingMemory{
				NodeID: "terminal",
				Facts:  make(map[string]*Fact),
				Tokens: make(map[string]*Token),
			},
		},
		onActivateLeft: func(token *Token) error {
			capturedToken = token
			return nil
		},
	}
	joinNode.AddChild(mockTerminal)

	personToken := NewTokenWithFact(personFact, "person", "type_node_person")
	err := joinNode.ActivateLeft(personToken)
	if err != nil {
		t.Fatalf("❌ ActivateLeft erreur: %v", err)
	}

	err = joinNode.ActivateRight(contactFact)
	if err != nil {
		t.Fatalf("❌ ActivateRight erreur: %v", err)
	}

	if capturedToken == nil {
		t.Fatal("❌ Aucun token propagé")
	}

	if capturedToken.Bindings.Len() != 2 {
		t.Errorf("❌ Attendu 2 bindings, reçu %d", capturedToken.Bindings.Len())
	}

	personBinding := capturedToken.GetBinding("person")
	if personBinding == nil || personBinding.ID != "Person~Alice_Dupont" {
		t.Errorf("❌ ID person incorrect")
	}

	contactBinding := capturedToken.GetBinding("contact")
	if contactBinding == nil || contactBinding.ID != "Contact~Alice_Dupont" {
		t.Errorf("❌ ID contact incorrect")
	}

	t.Log("✅ Test réussi: Join avec IDs PK composite")
}

// TestJoin_WithHashIDs teste les jointures avec IDs basés sur hash
func TestJoin_WithHashIDs(t *testing.T) {
	t.Log("🧪 TEST: Join - IDs basés sur hash")
	t.Log("===================================")

	eventFact := &Fact{
		ID:   "Event~a1b2c3d4e5f6g7h8",
		Type: "Event",
		Fields: map[string]interface{}{
			"timestamp": 1234567890,
			"message":   "User logged in",
		},
	}

	logFact := &Fact{
		ID:   "Log~x9y8z7w6v5u4t3s2",
		Type: "Log",
		Fields: map[string]interface{}{
			"event_id": "Event~a1b2c3d4e5f6g7h8",
			"details":  "Login successful",
		},
	}

	joinNode := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join_event_log",
			Children: []Node{},
			Memory: &WorkingMemory{
				NodeID: "join_event_log",
				Facts:  make(map[string]*Fact),
				Tokens: make(map[string]*Token),
			},
		},
		LeftVariables:  []string{"event"},
		RightVariables: []string{"log"},
		AllVariables:   []string{"event", "log"},
		VariableTypes: map[string]string{
			"event": "Event",
			"log":   "Log",
		},
		LeftMemory: &WorkingMemory{
			NodeID: "join_event_log_left",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		RightMemory: &WorkingMemory{
			NodeID: "join_event_log_right",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		ResultMemory: &WorkingMemory{
			NodeID: "join_event_log_result",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		JoinConditions: nil,
	}

	var capturedToken *Token
	mockTerminal := &mockTerminalNode{
		BaseNode: BaseNode{
			ID: "terminal",
			Memory: &WorkingMemory{
				NodeID: "terminal",
				Facts:  make(map[string]*Fact),
				Tokens: make(map[string]*Token),
			},
		},
		onActivateLeft: func(token *Token) error {
			capturedToken = token
			return nil
		},
	}
	joinNode.AddChild(mockTerminal)

	eventToken := NewTokenWithFact(eventFact, "event", "type_node_event")
	err := joinNode.ActivateLeft(eventToken)
	if err != nil {
		t.Fatalf("❌ ActivateLeft erreur: %v", err)
	}

	err = joinNode.ActivateRight(logFact)
	if err != nil {
		t.Fatalf("❌ ActivateRight erreur: %v", err)
	}

	if capturedToken == nil {
		t.Fatal("❌ Aucun token propagé")
	}

	eventBinding := capturedToken.GetBinding("event")
	if eventBinding == nil || eventBinding.ID != "Event~a1b2c3d4e5f6g7h8" {
		t.Errorf("❌ ID event incorrect")
	}

	logBinding := capturedToken.GetBinding("log")
	if logBinding == nil || logBinding.ID != "Log~x9y8z7w6v5u4t3s2" {
		t.Errorf("❌ ID log incorrect")
	}

	t.Log("✅ Test réussi: Join avec IDs hash")
}

// TestJoin_WithMixedIDFormats teste les jointures avec formats d'IDs mixtes
func TestJoin_WithMixedIDFormats(t *testing.T) {
	t.Log("🧪 TEST: Join - Formats d'IDs mixtes")
	t.Log("=====================================")

	// Fait avec PK simple
	userFact := &Fact{
		ID:   "User~Alice",
		Type: "User",
		Fields: map[string]interface{}{
			"name": "Alice",
		},
	}

	// Fait avec hash
	sessionFact := &Fact{
		ID:   "Session~f1e2d3c4b5a69788",
		Type: "Session",
		Fields: map[string]interface{}{
			"user_id":   "User~Alice",
			"timestamp": 1234567890,
		},
	}

	joinNode := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join_user_session",
			Children: []Node{},
			Memory: &WorkingMemory{
				NodeID: "join_user_session",
				Facts:  make(map[string]*Fact),
				Tokens: make(map[string]*Token),
			},
		},
		LeftVariables:  []string{"user"},
		RightVariables: []string{"session"},
		AllVariables:   []string{"user", "session"},
		VariableTypes: map[string]string{
			"user":    "User",
			"session": "Session",
		},
		LeftMemory: &WorkingMemory{
			NodeID: "join_user_session_left",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		RightMemory: &WorkingMemory{
			NodeID: "join_user_session_right",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		ResultMemory: &WorkingMemory{
			NodeID: "join_user_session_result",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		JoinConditions: nil,
	}

	var capturedToken *Token
	mockTerminal := &mockTerminalNode{
		BaseNode: BaseNode{
			ID: "terminal",
			Memory: &WorkingMemory{
				NodeID: "terminal",
				Facts:  make(map[string]*Fact),
				Tokens: make(map[string]*Token),
			},
		},
		onActivateLeft: func(token *Token) error {
			capturedToken = token
			return nil
		},
	}
	joinNode.AddChild(mockTerminal)

	userToken := NewTokenWithFact(userFact, "user", "type_node_user")
	err := joinNode.ActivateLeft(userToken)
	if err != nil {
		t.Fatalf("❌ ActivateLeft erreur: %v", err)
	}

	err = joinNode.ActivateRight(sessionFact)
	if err != nil {
		t.Fatalf("❌ ActivateRight erreur: %v", err)
	}

	if capturedToken == nil {
		t.Fatal("❌ Aucun token propagé")
	}

	// Vérifier formats d'IDs différents
	userBinding := capturedToken.GetBinding("user")
	if userBinding == nil || userBinding.ID != "User~Alice" {
		t.Errorf("❌ ID user (PK simple) incorrect: %v", userBinding)
	}

	sessionBinding := capturedToken.GetBinding("session")
	if sessionBinding == nil || sessionBinding.ID != "Session~f1e2d3c4b5a69788" {
		t.Errorf("❌ ID session (hash) incorrect: %v", sessionBinding)
	}

	t.Log("✅ Test réussi: Join avec formats d'IDs mixtes")
}

// TestJoin_NoMatch_DifferentIDs teste qu'une jointure échoue avec des IDs incompatibles
func TestJoin_NoMatch_DifferentIDs(t *testing.T) {
	t.Log("🧪 TEST: Join - Pas de match avec IDs différents")
	t.Log("=================================================")

	personFact := &Fact{
		ID:   "Person~Alice",
		Type: "Person",
		Fields: map[string]interface{}{
			"nom": "Alice",
		},
	}

	// Membership pour Bob, pas Alice
	membershipFact := &Fact{
		ID:   "Membership~Bob_Club1",
		Type: "Membership",
		Fields: map[string]interface{}{
			"person_id": "Person~Bob",
			"club":      "Club1",
		},
	}

	joinNode := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join_person_membership",
			Children: []Node{},
			Memory: &WorkingMemory{
				NodeID: "join_person_membership",
				Facts:  make(map[string]*Fact),
				Tokens: make(map[string]*Token),
			},
		},
		LeftVariables:  []string{"person"},
		RightVariables: []string{"membership"},
		AllVariables:   []string{"person", "membership"},
		VariableTypes: map[string]string{
			"person":     "Person",
			"membership": "Membership",
		},
		LeftMemory: &WorkingMemory{
			NodeID: "join_person_membership_left",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		RightMemory: &WorkingMemory{
			NodeID: "join_person_membership_right",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		ResultMemory: &WorkingMemory{
			NodeID: "join_person_membership_result",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		// Condition de join: person.id == membership.person_id
		JoinConditions: []JoinCondition{
			{
				LeftVar:    "person",
				LeftField:  "id",
				RightVar:   "membership",
				RightField: "person_id",
				Operator:   "==",
			},
		},
	}

	var capturedToken *Token
	mockTerminal := &mockTerminalNode{
		BaseNode: BaseNode{
			ID: "terminal",
			Memory: &WorkingMemory{
				NodeID: "terminal",
				Facts:  make(map[string]*Fact),
				Tokens: make(map[string]*Token),
			},
		},
		onActivateLeft: func(token *Token) error {
			capturedToken = token
			return nil
		},
	}
	joinNode.AddChild(mockTerminal)

	personToken := NewTokenWithFact(personFact, "person", "type_node_person")
	err := joinNode.ActivateLeft(personToken)
	if err != nil {
		t.Fatalf("❌ ActivateLeft erreur: %v", err)
	}

	err = joinNode.ActivateRight(membershipFact)
	if err != nil {
		t.Fatalf("❌ ActivateRight erreur: %v", err)
	}

	// Vérifier qu'aucun token n'a été propagé (pas de match)
	if capturedToken != nil {
		t.Errorf("❌ Aucun token ne devrait être propagé car les IDs ne matchent pas")
	}

	t.Log("✅ Test réussi: Pas de match avec IDs différents")
}

// TestJoin_CascadeWithGeneratedIDs teste une cascade de joins avec IDs générés
func TestJoin_CascadeWithGeneratedIDs(t *testing.T) {
	t.Log("🧪 TEST: Join - Cascade avec IDs générés")
	t.Log("=========================================")

	userFact := &Fact{
		ID:   "User~Alice",
		Type: "User",
		Fields: map[string]interface{}{
			"name": "Alice",
		},
	}

	orderFact := &Fact{
		ID:   "Order~Order123",
		Type: "Order",
		Fields: map[string]interface{}{
			"user_id":    "User~Alice",
			"product_id": "Product~Laptop_15inch",
		},
	}

	productFact := &Fact{
		ID:   "Product~Laptop_15inch",
		Type: "Product",
		Fields: map[string]interface{}{
			"name": "Laptop",
			"size": "15inch",
		},
	}

	// Premier join: User + Order
	joinNode1 := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join1_user_order",
			Children: []Node{},
			Memory: &WorkingMemory{
				NodeID: "join1_user_order",
				Facts:  make(map[string]*Fact),
				Tokens: make(map[string]*Token),
			},
		},
		LeftVariables:  []string{"user"},
		RightVariables: []string{"order"},
		AllVariables:   []string{"user", "order"},
		VariableTypes: map[string]string{
			"user":  "User",
			"order": "Order",
		},
		LeftMemory: &WorkingMemory{
			NodeID: "join1_user_order_left",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		RightMemory: &WorkingMemory{
			NodeID: "join1_user_order_right",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		ResultMemory: &WorkingMemory{
			NodeID: "join1_user_order_result",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		JoinConditions: nil,
	}

	// Deuxième join: (User+Order) + Product
	joinNode2 := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join2_order_product",
			Children: []Node{},
			Memory: &WorkingMemory{
				NodeID: "join2_order_product",
				Facts:  make(map[string]*Fact),
				Tokens: make(map[string]*Token),
			},
		},
		LeftVariables:  []string{"user", "order"},
		RightVariables: []string{"product"},
		AllVariables:   []string{"user", "order", "product"},
		VariableTypes: map[string]string{
			"user":    "User",
			"order":   "Order",
			"product": "Product",
		},
		LeftMemory: &WorkingMemory{
			NodeID: "join2_order_product_left",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		RightMemory: &WorkingMemory{
			NodeID: "join2_order_product_right",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		ResultMemory: &WorkingMemory{
			NodeID: "join2_order_product_result",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		},
		JoinConditions: nil,
	}

	var finalToken *Token
	mockTerminal := &mockTerminalNode{
		BaseNode: BaseNode{
			ID: "terminal",
			Memory: &WorkingMemory{
				NodeID: "terminal",
				Facts:  make(map[string]*Fact),
				Tokens: make(map[string]*Token),
			},
		},
		onActivateLeft: func(token *Token) error {
			finalToken = token
			return nil
		},
	}

	// Connecter les nœuds
	joinNode1.AddChild(joinNode2)
	joinNode2.AddChild(mockTerminal)

	// Activer le premier join
	userToken := NewTokenWithFact(userFact, "user", "type_node_user")
	err := joinNode1.ActivateLeft(userToken)
	if err != nil {
		t.Fatalf("❌ Join1 ActivateLeft erreur: %v", err)
	}

	err = joinNode1.ActivateRight(orderFact)
	if err != nil {
		t.Fatalf("❌ Join1 ActivateRight erreur: %v", err)
	}

	// Activer le deuxième join
	err = joinNode2.ActivateRight(productFact)
	if err != nil {
		t.Fatalf("❌ Join2 ActivateRight erreur: %v", err)
	}

	// Vérifier le token final
	if finalToken == nil {
		t.Fatal("❌ Aucun token final propagé")
	}

	if finalToken.Bindings.Len() != 3 {
		t.Errorf("❌ Attendu 3 bindings, reçu %d", finalToken.Bindings.Len())
	}

	// Vérifier tous les bindings et leurs IDs
	userBinding := finalToken.GetBinding("user")
	if userBinding == nil || userBinding.ID != "User~Alice" {
		t.Errorf("❌ Binding user incorrect")
	}

	orderBinding := finalToken.GetBinding("order")
	if orderBinding == nil || orderBinding.ID != "Order~Order123" {
		t.Errorf("❌ Binding order incorrect")
	}

	productBinding := finalToken.GetBinding("product")
	if productBinding == nil || productBinding.ID != "Product~Laptop_15inch" {
		t.Errorf("❌ Binding product incorrect")
	}

	t.Log("✅ Test réussi: Cascade avec IDs générés")
}
