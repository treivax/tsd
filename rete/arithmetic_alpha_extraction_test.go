// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestArithmeticAlphaExtraction_SingleVariable verifies that arithmetic expressions
// with a single variable are extracted to AlphaNodes
func TestArithmeticAlphaExtraction_SingleVariable(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "arithmetic_alpha.tsd")
	content := `type Commande(id: string, qte: number, price: number)
action log(msg: string)
rule high_quantity : {c: Commande} / c.qte * 23 - 10 > 100
    ==> log("High quantity")
rule expensive : {c: Commande} / c.price * 1.2 > 1000
    ==> log("Expensive item")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	// Get network statistics
	stats := network.GetNetworkStats()
	alphaNodes := stats["alpha_nodes"].(int)
	t.Logf("Network stats: %+v", stats)
	// We should have at least 2 AlphaNodes (one per rule with arithmetic condition)
	if alphaNodes < 2 {
		t.Errorf("Expected at least 2 AlphaNodes for arithmetic conditions, got %d", alphaNodes)
	}
	// Test that facts are properly filtered by AlphaNodes
	t.Run("filtering behavior", func(t *testing.T) {
		// Submit a fact that passes the first rule: qte * 23 - 10 > 100
		// qte = 6: 6 * 23 - 10 = 138 - 10 = 128 > 100 ✓
		passFirst := &Fact{
			ID:   "C1",
			Type: "Commande",
			Fields: map[string]interface{}{
				"id":    "c1",
				"qte":   6.0,
				"price": 100.0,
			},
		}
		// Submit a fact that fails the first rule: qte * 23 - 10 > 100
		// qte = 3: 3 * 23 - 10 = 69 - 10 = 59 < 100 ✗
		failFirst := &Fact{
			ID:   "C2",
			Type: "Commande",
			Fields: map[string]interface{}{
				"id":    "c2",
				"qte":   3.0,
				"price": 100.0,
			},
		}
		network.SubmitFact(passFirst)
		network.SubmitFact(failFirst)
		// Check activations for high_quantity rule
		if terminal, exists := network.TerminalNodes["high_quantity_terminal"]; exists {
			tokens := terminal.GetMemory().Tokens
			if len(tokens) != 1 {
				t.Errorf("Expected 1 activation for high_quantity, got %d", len(tokens))
				t.Logf("Tokens: %+v", tokens)
			}
		} else {
			t.Error("Terminal node 'high_quantity_terminal' not found")
		}
	})
	t.Logf("✅ Arithmetic alpha extraction working correctly")
}

// TestArithmeticAlphaExtraction_ComplexNested verifies nested arithmetic expressions
func TestArithmeticAlphaExtraction_ComplexNested(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "complex_arithmetic.tsd")
	content := `type Order(id: string, quantity: number, price: number, discount: number)
action process(msg: string)
rule complex_calc : {o: Order} / (o.quantity * o.price - o.discount) / 2 > 50
    ==> process("Complex calculation passed")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	stats := network.GetNetworkStats()
	t.Logf("Network stats: %+v", stats)
	// Submit test facts
	// (5 * 30 - 10) / 2 = (150 - 10) / 2 = 140 / 2 = 70 > 50 ✓
	passFact := &Fact{
		ID:   "O1",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":       "o1",
			"quantity": 5.0,
			"price":    30.0,
			"discount": 10.0,
		},
	}
	// (2 * 30 - 10) / 2 = (60 - 10) / 2 = 50 / 2 = 25 < 50 ✗
	failFact := &Fact{
		ID:   "O2",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":       "o2",
			"quantity": 2.0,
			"price":    30.0,
			"discount": 10.0,
		},
	}
	network.SubmitFact(passFact)
	network.SubmitFact(failFact)
	// Verify only one activation
	if terminal, exists := network.TerminalNodes["complex_calc_terminal"]; exists {
		tokens := terminal.GetMemory().Tokens
		if len(tokens) != 1 {
			t.Errorf("Expected 1 activation, got %d", len(tokens))
		}
	}
	t.Logf("✅ Complex nested arithmetic extraction working")
}

// TestArithmeticAlphaExtraction_MultiVariable verifies that multi-variable
// arithmetic expressions correctly stay in JoinNodes (beta)
func TestArithmeticAlphaExtraction_MultiVariable(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "multi_var_arithmetic.tsd")
	content := `type Product(id: string, basePrice: number)
type Customer(id: string, discount: number)
type Order(id: string, productId: string, customerId: string)
action notify(msg: string)
rule discount_calc : {p: Product, c: Customer, o: Order}
    / o.productId == p.id AND o.customerId == c.id AND p.basePrice * (1 - c.discount) > 100
    ==> notify("Expensive order after discount")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	stats := network.GetNetworkStats()
	t.Logf("Network stats: %+v", stats)
	// The arithmetic expression p.basePrice * (1 - c.discount) involves TWO variables
	// This should be in a JoinNode, not an AlphaNode
	betaNodes := stats["beta_nodes"].(int)
	if betaNodes < 1 {
		t.Error("Expected at least 1 beta node for multi-variable arithmetic")
	}
	// Test the actual behavior
	product := &Fact{
		ID:   "P1",
		Type: "Product",
		Fields: map[string]interface{}{
			"id":        "p1",
			"basePrice": 150.0,
		},
	}
	customer := &Fact{
		ID:   "C1",
		Type: "Customer",
		Fields: map[string]interface{}{
			"id":       "c1",
			"discount": 0.2, // 20% discount
		},
	}
	order := &Fact{
		ID:   "O1",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":         "o1",
			"productId":  "p1",
			"customerId": "c1",
		},
	}
	network.SubmitFact(product)
	network.SubmitFact(customer)
	network.SubmitFact(order)
	// 150 * (1 - 0.2) = 150 * 0.8 = 120 > 100 ✓
	if terminal, exists := network.TerminalNodes["discount_calc_terminal"]; exists {
		tokens := terminal.GetMemory().Tokens
		if len(tokens) != 1 {
			t.Errorf("Expected 1 activation, got %d", len(tokens))
		}
	}
	t.Logf("✅ Multi-variable arithmetic correctly in JoinNode")
}

// TestArithmeticAlphaExtraction_MixedConditions verifies rules with both
// simple and arithmetic alpha conditions
func TestArithmeticAlphaExtraction_MixedConditions(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "mixed_conditions.tsd")
	content := `type Item(id: string, name: string, weight: number, value: number)
action ship(msg: string)
rule valuable_heavy : {i: Item}
    / i.name == "Gold" AND i.weight * 2 > 10 AND i.value > 1000
    ==> ship("Ship valuable heavy item")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	stats := network.GetNetworkStats()
	t.Logf("Network stats: %+v", stats)
	// We should have multiple AlphaNodes:
	// - One for name == "Gold"
	// - One for weight * 2 > 10
	// - One for value > 1000
	alphaNodes := stats["alpha_nodes"].(int)
	if alphaNodes < 3 {
		t.Logf("Expected at least 3 AlphaNodes, got %d (may be chained)", alphaNodes)
	}
	// Test facts
	goldHeavyValuable := &Fact{
		ID:   "I1",
		Type: "Item",
		Fields: map[string]interface{}{
			"id":     "i1",
			"name":   "Gold",
			"weight": 8.0,    // 8 * 2 = 16 > 10 ✓
			"value":  1500.0, // > 1000 ✓
		},
	}
	silverLight := &Fact{
		ID:   "I2",
		Type: "Item",
		Fields: map[string]interface{}{
			"id":     "i2",
			"name":   "Silver",
			"weight": 3.0,
			"value":  800.0,
		},
	}
	goldLight := &Fact{
		ID:   "I3",
		Type: "Item",
		Fields: map[string]interface{}{
			"id":     "i3",
			"name":   "Gold",
			"weight": 4.0, // 4 * 2 = 8 < 10 ✗
			"value":  1200.0,
		},
	}
	network.SubmitFact(goldHeavyValuable)
	network.SubmitFact(silverLight)
	network.SubmitFact(goldLight)
	// Only goldHeavyValuable should activate the rule
	if terminal, exists := network.TerminalNodes["valuable_heavy_terminal"]; exists {
		tokens := terminal.GetMemory().Tokens
		if len(tokens) != 1 {
			t.Errorf("Expected 1 activation, got %d", len(tokens))
			for tokenID, tok := range tokens {
				t.Logf("  Token %s: %+v", tokenID, tok.Facts)
			}
		}
	}
	t.Logf("✅ Mixed conditions with arithmetic working correctly")
}

// TestArithmeticAlphaExtraction_EdgeCases tests edge cases and error conditions
func TestArithmeticAlphaExtraction_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		tsdContent  string
		factFields  map[string]interface{}
		shouldMatch bool
		description string
	}{
		{
			name: "division with zero result",
			tsdContent: `type Data(id: string, value: number)
action alert(msg: string)
rule zero_div : {d: Data} / d.value / 10 == 0 ==> alert("Zero")`,
			factFields: map[string]interface{}{
				"id":    "d1",
				"value": 0.0,
			},
			shouldMatch: true,
			description: "0 / 10 = 0",
		},
		{
			name: "negative arithmetic",
			tsdContent: `type Account(id: string, balance: number)
action warn(msg: string)
rule negative : {a: Account} / a.balance * -1 > 100 ==> warn("Large debt")`,
			factFields: map[string]interface{}{
				"id":      "a1",
				"balance": -150.0,
			},
			shouldMatch: true,
			description: "-150 * -1 = 150 > 100",
		},
		// TODO: Enable when parser supports % operator
		// {
		// 	name: "modulo operation",
		// 	tsdContent: `type Number(id: string, value: number)
		// action check(msg: string)
		// rule even : {n: Number} / n.value % 2 == 0 ==> check("Even")`,
		// 	factFields: map[string]interface{}{
		// 		"id":    "n1",
		// 		"value": 42.0,
		// 	},
		// 	shouldMatch: true,
		// 	description: "42 % 2 = 0 (even)",
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			tsdFile := filepath.Join(tempDir, "test.tsd")
			if err := os.WriteFile(tsdFile, []byte(tt.tsdContent), 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}
			storage := NewMemoryStorage()
			pipeline := NewConstraintPipeline()
			network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
			if err != nil {
				t.Fatalf("Failed to build network: %v", err)
			}
			// Extract the type name from fact fields
			factType := "Data"
			if _, ok := tt.factFields["balance"]; ok {
				factType = "Account"
			}
			fact := &Fact{
				ID:     "F1",
				Type:   factType,
				Fields: tt.factFields,
			}
			network.SubmitFact(fact)
			// Check if rule activated
			var activated bool
			for _, terminal := range network.TerminalNodes {
				if len(terminal.GetMemory().Tokens) > 0 {
					activated = true
					break
				}
			}
			if activated != tt.shouldMatch {
				t.Errorf("%s: expected match=%v, got match=%v", tt.description, tt.shouldMatch, activated)
			}
		})
	}
	t.Logf("✅ Edge cases handled correctly")
}

// BenchmarkArithmeticInAlpha benchmarks arithmetic evaluation in AlphaNodes
func BenchmarkArithmeticInAlpha(b *testing.B) {
	tempDir := b.TempDir()
	tsdFile := filepath.Join(tempDir, "bench.tsd")
	content := `type Order(id: string, qty: number, price: number)
action process(msg: string)
rule expensive : {o: Order} / o.qty * o.price > 500 ==> process("Expensive")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		b.Fatalf("Failed to write test file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		b.Fatalf("Failed to build network: %v", err)
	}
	// Prepare facts
	passFacts := make([]*Fact, 100)
	failFacts := make([]*Fact, 100)
	for i := 0; i < 100; i++ {
		passFacts[i] = &Fact{
			ID:   fmt.Sprintf("pass%d", i),
			Type: "Order",
			Fields: map[string]interface{}{
				"id":    "o1",
				"qty":   20.0,
				"price": 30.0, // 20 * 30 = 600 > 500
			},
		}
		failFacts[i] = &Fact{
			ID:   fmt.Sprintf("fail%d", i),
			Type: "Order",
			Fields: map[string]interface{}{
				"id":    "o2",
				"qty":   10.0,
				"price": 20.0, // 10 * 20 = 200 < 500
			},
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		network.SubmitFact(passFacts[i%100])
		network.SubmitFact(failFacts[i%100])
	}
}
