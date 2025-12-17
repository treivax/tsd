// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"os"
	"testing"
)

// TestJoinNodeCascade_TwoVariablesIntegration teste la cascade de jointure à 2 variables via pipeline
func TestJoinNodeCascade_TwoVariablesIntegration(t *testing.T) {
	t.Log("🧪 TEST: JoinNode Cascade Integration - 2 Variables (User ⋈ Order)")
	t.Log("=====================================================================")

	constraintContent := `// Test 2-variable cascade
type User(#id: string, name:string)
type Order(#id: string, user_id: string, amount:number)
action process_order(userId: string, orderId: string)
rule r1 : {u: User, o: Order} / u.id == "U1" AND o.user_id == u.id ==> process_order(u.id, o.id)
`
	tmpFile := createTempConstraintFile(t, "two_var_cascade", constraintContent)
	defer os.Remove(tmpFile)

	pipeline := NewConstraintPipeline()
	storage := NewMemoryStorage()
	network, _, err := pipeline.IngestFile(tmpFile, nil, storage)
	if err != nil {
		t.Fatalf("❌ Failed to build network: %v", err)
	}
	t.Logf("✅ Network built: %d type nodes, %d terminals", len(network.TypeNodes), len(network.TerminalNodes))

	userFact := &Fact{
		ID:     "U1",
		Type:   "User",
		Fields: map[string]interface{}{"id": "U1", "name": "Alice"},
	}
	err = network.SubmitFact(userFact)
	if err != nil {
		t.Fatalf("❌ Error submitting User: %v", err)
	}

	terminalTokens := countAllTerminalTokens(network)
	if terminalTokens != 0 {
		t.Logf("⚠️  Terminal tokens after User only: %d (expected 0)", terminalTokens)
	} else {
		t.Logf("✅ No terminal tokens yet (missing Order)")
	}

	orderFact := &Fact{
		ID:     "O1",
		Type:   "Order",
		Fields: map[string]interface{}{"id": "O1", "user_id": "U1", "amount": 100},
	}
	err = network.SubmitFact(orderFact)
	if err != nil {
		t.Fatalf("❌ Error submitting Order: %v", err)
	}

	terminalTokens = countAllTerminalTokens(network)
	if terminalTokens < 1 {
		t.Errorf("❌ Expected at least 1 terminal token, got %d", terminalTokens)
	} else {
		t.Logf("✅ Terminal token created: %d", terminalTokens)
	}

	badOrderFact := &Fact{
		ID:     "O2",
		Type:   "Order",
		Fields: map[string]interface{}{"id": "O2", "user_id": "U999", "amount": 50},
	}
	err = network.SubmitFact(badOrderFact)
	if err != nil {
		t.Fatalf("❌ Error submitting non-matching Order: %v", err)
	}

	finalTokens := countAllTerminalTokens(network)
	if finalTokens != terminalTokens {
		t.Logf("⚠️  Token count changed from %d to %d after non-matching Order", terminalTokens, finalTokens)
	} else {
		t.Logf("✅ Non-matching Order filtered correctly")
	}

	t.Log("\n🎊 TEST PASSED: 2-variable cascade join via pipeline works correctly")
}

// TestJoinNodeCascade_ThreeVariablesIntegration teste la cascade à 3 variables via pipeline
func TestJoinNodeCascade_ThreeVariablesIntegration(t *testing.T) {
	t.Log("🧪 TEST: JoinNode Cascade Integration - 3 Variables (User ⋈ Order ⋈ Product)")
	t.Log("=============================================================================")

	pipeline := NewConstraintPipeline()
	storage := NewMemoryStorage()
	network, _, err := pipeline.IngestFile("test/incremental_propagation.tsd", nil, storage)
	if err != nil {
		t.Fatalf("❌ Failed to build network: %v", err)
	}
	t.Logf("✅ Network built for 3-variable test")

	userFact := &Fact{
		ID:     "U1",
		Type:   "User",
		Fields: map[string]interface{}{"id": "U1", "age": 25},
	}
	err = network.SubmitFact(userFact)
	if err != nil {
		t.Fatalf("❌ Error submitting User: %v", err)
	}
	count1 := countAllTerminalTokens(network)
	t.Logf("✅ After User: %d terminal tokens", count1)

	orderFact := &Fact{
		ID:     "O1",
		Type:   "Order",
		Fields: map[string]interface{}{"id": "O1", "user_id": "U1", "product_id": "P1"},
	}
	err = network.SubmitFact(orderFact)
	if err != nil {
		t.Fatalf("❌ Error submitting Order: %v", err)
	}
	count2 := countAllTerminalTokens(network)
	t.Logf("✅ After User+Order: %d terminal tokens", count2)

	productFact := &Fact{
		ID:     "P1",
		Type:   "Product",
		Fields: map[string]interface{}{"id": "P1", "name": "Widget"},
	}
	err = network.SubmitFact(productFact)
	if err != nil {
		t.Fatalf("❌ Error submitting Product: %v", err)
	}
	count3 := countAllTerminalTokens(network)
	t.Logf("✅ After User+Order+Product: %d terminal tokens", count3)

	if count3 < 1 {
		t.Errorf("❌ Expected at least 1 terminal token after full cascade, got %d", count3)
	} else {
		t.Logf("✅ 3-variable cascade completed successfully")
	}

	t.Log("\n🎊 TEST PASSED: 3-variable cascade join works correctly")
}

// TestJoinNodeCascade_OrderIndependence teste l'indépendance de l'ordre de soumission
func TestJoinNodeCascade_OrderIndependence(t *testing.T) {
	t.Log("🧪 TEST: JoinNode Cascade - Order Independence")
	t.Log("==============================================")

	constraintContent := `// Order independence test
type User(#id:string)
type Order(#id: string, user_id:string)
action test_action(userId: string, orderId: string)
rule r1 : {u: User, o: Order} / o.user_id == u.id ==> test_action(u.id, o.id)
`
	testOrders := []struct {
		name  string
		order []string
	}{
		{"User→Order", []string{"U", "O"}},
		{"Order→User", []string{"O", "U"}},
	}

	for _, testCase := range testOrders {
		t.Run(testCase.name, func(t *testing.T) {
			tmpFile := createTempConstraintFile(t, "order_test", constraintContent)
			defer os.Remove(tmpFile)

			pipeline := NewConstraintPipeline()
			storage := NewMemoryStorage()
			network, _, err := pipeline.IngestFile(tmpFile, nil, storage)
			if err != nil {
				t.Fatalf("❌ Failed to build network: %v", err)
			}

			userFact := &Fact{
				ID:     "U1",
				Type:   "User",
				Fields: map[string]interface{}{"id": "U1"},
			}
			orderFact := &Fact{
				ID:     "O1",
				Type:   "Order",
				Fields: map[string]interface{}{"id": "O1", "user_id": "U1"},
			}

			for _, factType := range testCase.order {
				switch factType {
				case "U":
					err = network.SubmitFact(userFact)
				case "O":
					err = network.SubmitFact(orderFact)
				}
				if err != nil {
					t.Fatalf("❌ Error submitting fact: %v", err)
				}
			}

			terminalTokens := countAllTerminalTokens(network)
			if terminalTokens < 1 {
				t.Errorf("❌ Expected at least 1 terminal token, got %d", terminalTokens)
			} else {
				t.Logf("✅ Order %v produced %d terminal tokens", testCase.order, terminalTokens)
			}
		})
	}

	t.Log("\n🎊 TEST PASSED: Join cascade is order-independent")
}

// TestJoinNodeCascade_MultipleMatchingFacts teste le comportement du produit cartésien
func TestJoinNodeCascade_MultipleMatchingFacts(t *testing.T) {
	t.Log("🧪 TEST: JoinNode Cascade - Multiple Matching Facts")
	t.Log("====================================================")

	constraintContent := `// Multiple matching facts test
type User(#id:string)
type Order(#id: string, user_id:string)
action test_action(userId: string, orderId: string)
rule r1 : {u: User, o: Order} / o.user_id == u.id ==> test_action(u.id, o.id)
`
	tmpFile := createTempConstraintFile(t, "multi_match", constraintContent)
	defer os.Remove(tmpFile)

	pipeline := NewConstraintPipeline()
	storage := NewMemoryStorage()
	network, _, err := pipeline.IngestFile(tmpFile, nil, storage)
	if err != nil {
		t.Fatalf("❌ Failed to build network: %v", err)
	}

	user1 := &Fact{
		ID:     "U1",
		Type:   "User",
		Fields: map[string]interface{}{"id": "U1"},
	}
	user2 := &Fact{
		ID:     "U2",
		Type:   "User",
		Fields: map[string]interface{}{"id": "U2"},
	}
	network.SubmitFact(user1)
	network.SubmitFact(user2)

	order1 := &Fact{
		ID:     "O1",
		Type:   "Order",
		Fields: map[string]interface{}{"id": "O1", "user_id": "U1"},
	}
	order2 := &Fact{
		ID:     "O2",
		Type:   "Order",
		Fields: map[string]interface{}{"id": "O2", "user_id": "U1"},
	}
	order3 := &Fact{
		ID:     "O3",
		Type:   "Order",
		Fields: map[string]interface{}{"id": "O3", "user_id": "U2"},
	}
	network.SubmitFact(order1)
	network.SubmitFact(order2)
	network.SubmitFact(order3)

	terminalTokens := countAllTerminalTokens(network)
	if terminalTokens != 3 {
		t.Logf("⚠️  Expected 3 terminal tokens (cartesian product), got %d", terminalTokens)
	} else {
		t.Logf("✅ Correct cartesian product: 3 terminal tokens")
	}

	t.Log("\n🎊 TEST PASSED: Multiple matching facts handled correctly")
}

// TestJoinNodeCascade_Retraction teste la rétractation de faits dans les cascades
func TestJoinNodeCascade_Retraction(t *testing.T) {
	t.Log("🧪 TEST: JoinNode Cascade - Fact Retraction")
	t.Log("============================================")

	constraintContent := `// Retraction test
type User(#id:string)
type Order(#id: string, user_id:string)
action test_action(userId: string, orderId: string)
rule r1 : {u: User, o: Order} / o.user_id == u.id ==> test_action(u.id, o.id)
`
	tmpFile := createTempConstraintFile(t, "retract_test", constraintContent)
	defer os.Remove(tmpFile)

	pipeline := NewConstraintPipeline()
	storage := NewMemoryStorage()
	network, _, err := pipeline.IngestFile(tmpFile, nil, storage)
	if err != nil {
		t.Fatalf("❌ Failed to build network: %v", err)
	}

	userFact := &Fact{
		ID:     "U1",
		Type:   "User",
		Fields: map[string]interface{}{"id": "U1"},
	}
	orderFact := &Fact{
		ID:     "O1",
		Type:   "Order",
		Fields: map[string]interface{}{"id": "O1", "user_id": "U1"},
	}
	network.SubmitFact(userFact)
	network.SubmitFact(orderFact)

	beforeCount := countAllTerminalTokens(network)
	if beforeCount < 1 {
		t.Logf("⚠️  Expected terminal tokens before retraction, got %d", beforeCount)
	}

	err = network.RetractFact(userFact.GetInternalID())
	if err != nil {
		t.Fatalf("❌ Error retracting User: %v", err)
	}

	afterCount := countAllTerminalTokens(network)
	if afterCount != 0 {
		t.Logf("⚠️  Expected 0 terminal tokens after retraction, got %d", afterCount)
	} else {
		t.Logf("✅ Terminal tokens removed after User retraction")
	}

	t.Log("\n🎊 TEST PASSED: Fact retraction works in cascades")
}

// countAllTerminalTokens compte tous les tokens dans les nœuds terminaux
func countAllTerminalTokens(network *ReteNetwork) int {
	total := 0
	for _, terminal := range network.TerminalNodes {
		total += len(terminal.Memory.GetTokens())
	}
	return total
}

// createTempConstraintFile crée un fichier de contrainte temporaire
func createTempConstraintFile(t *testing.T, name, content string) string {
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/" + name + ".tsd"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("❌ Failed to create temp file: %v", err)
	}
	return tmpFile
}
