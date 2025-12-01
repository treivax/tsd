// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIncrementalFactsParsing_SingleType tests adding facts incrementally for a single type
func TestIncrementalFactsParsing_SingleType(t *testing.T) {
	ps := NewProgramState()

	// Step 1: Define type
	typeContent := `type Person(id: string, name: string, age:number)`
	err := ps.ParseAndMergeContent(typeContent, "types.tsd")
	if err != nil {
		t.Fatalf("Failed to parse type: %v", err)
	}

	if len(ps.Types) != 1 {
		t.Errorf("Expected 1 type, got %d", len(ps.Types))
	}

	// Step 2: Add first batch of facts
	factsContent1 := `
Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)
`
	err = ps.ParseAndMergeContent(factsContent1, "facts1.tsd")
	if err != nil {
		t.Fatalf("Failed to parse first batch of facts: %v", err)
	}

	if len(ps.Facts) != 2 {
		t.Errorf("After first batch: expected 2 facts, got %d", len(ps.Facts))
	}

	// Step 3: Add second batch of facts
	factsContent2 := `
Person(id: "P003", name: "Charlie", age: 35)
Person(id: "P004", name: "Diana", age: 28)
Person(id: "P005", name: "Eve", age: 42)
`
	err = ps.ParseAndMergeContent(factsContent2, "facts2.tsd")
	if err != nil {
		t.Fatalf("Failed to parse second batch of facts: %v", err)
	}

	if len(ps.Facts) != 5 {
		t.Errorf("After second batch: expected 5 facts, got %d", len(ps.Facts))
	}

	// Step 4: Add third batch of facts
	factsContent3 := `
Person(id: "P006", name: "Frank", age: 50)
`
	err = ps.ParseAndMergeContent(factsContent3, "facts3.tsd")
	if err != nil {
		t.Fatalf("Failed to parse third batch of facts: %v", err)
	}

	if len(ps.Facts) != 6 {
		t.Errorf("After third batch: expected 6 facts, got %d", len(ps.Facts))
	}

	// Verify no errors occurred
	if ps.HasErrors() {
		t.Errorf("Unexpected validation errors: %v", ps.GetErrors())
	}

	// Verify all facts are unique
	idsSeen := make(map[string]bool)
	for _, fact := range ps.Facts {
		for _, field := range fact.Fields {
			if field.Name == "id" {
				id := field.Value.Value.(string)
				if idsSeen[id] {
					t.Errorf("Duplicate fact ID found: %s", id)
				}
				idsSeen[id] = true
			}
		}
	}

	if len(idsSeen) != 6 {
		t.Errorf("Expected 6 unique fact IDs, got %d", len(idsSeen))
	}
}

// TestIncrementalFactsParsing_MultipleTypes tests adding facts for different types incrementally
func TestIncrementalFactsParsing_MultipleTypes(t *testing.T) {
	ps := NewProgramState()

	// Step 1: Define multiple types
	typesContent := `
type Person(id: string, name: string, age:number)
type Product(id: string, name: string, price:number)
type Order(id: string, customer_id: string, product_id: string, quantity:number)
`
	err := ps.ParseAndMergeContent(typesContent, "types.tsd")
	if err != nil {
		t.Fatalf("Failed to parse types: %v", err)
	}

	if len(ps.Types) != 3 {
		t.Errorf("Expected 3 types, got %d", len(ps.Types))
	}

	// Step 2: Add Person facts
	personFacts := `
Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)
`
	err = ps.ParseAndMergeContent(personFacts, "person_facts.tsd")
	if err != nil {
		t.Fatalf("Failed to parse Person facts: %v", err)
	}

	if len(ps.Facts) != 2 {
		t.Errorf("After Person facts: expected 2 facts, got %d", len(ps.Facts))
	}

	// Step 3: Add Product facts
	productFacts := `
Product(id: "PR001", name: "Laptop", price: 999.99)
Product(id: "PR002", name: "Mouse", price: 29.99)
Product(id: "PR003", name: "Keyboard", price: 79.99)
`
	err = ps.ParseAndMergeContent(productFacts, "product_facts.tsd")
	if err != nil {
		t.Fatalf("Failed to parse Product facts: %v", err)
	}

	if len(ps.Facts) != 5 {
		t.Errorf("After Product facts: expected 5 facts, got %d", len(ps.Facts))
	}

	// Step 4: Add Order facts
	orderFacts := `
Order(id: "O001", customer_id: "P001", product_id: "PR001", quantity: 1)
Order(id: "O002", customer_id: "P002", product_id: "PR002", quantity: 2)
`
	err = ps.ParseAndMergeContent(orderFacts, "order_facts.tsd")
	if err != nil {
		t.Fatalf("Failed to parse Order facts: %v", err)
	}

	if len(ps.Facts) != 7 {
		t.Errorf("After Order facts: expected 7 facts, got %d", len(ps.Facts))
	}

	// Step 5: Add more mixed facts
	mixedFacts := `
Person(id: "P003", name: "Charlie", age: 35)
Product(id: "PR004", name: "Monitor", price: 299.99)
Order(id: "O003", customer_id: "P003", product_id: "PR004", quantity: 1)
`
	err = ps.ParseAndMergeContent(mixedFacts, "mixed_facts.tsd")
	if err != nil {
		t.Fatalf("Failed to parse mixed facts: %v", err)
	}

	if len(ps.Facts) != 10 {
		t.Errorf("After mixed facts: expected 10 facts, got %d", len(ps.Facts))
	}

	// Verify no errors
	if ps.HasErrors() {
		t.Errorf("Unexpected validation errors: %v", ps.GetErrors())
	}

	// Verify fact type distribution
	factsByType := make(map[string]int)
	for _, fact := range ps.Facts {
		factsByType[fact.TypeName]++
	}

	if factsByType["Person"] != 3 {
		t.Errorf("Expected 3 Person facts, got %d", factsByType["Person"])
	}
	if factsByType["Product"] != 4 {
		t.Errorf("Expected 4 Product facts, got %d", factsByType["Product"])
	}
	if factsByType["Order"] != 3 {
		t.Errorf("Expected 3 Order facts, got %d", factsByType["Order"])
	}
}

// TestIncrementalFactsParsing_WithFiles tests incremental parsing using actual files
func TestIncrementalFactsParsing_WithFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "incremental_facts_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ps := NewProgramState()

	// Create and parse types file
	typesFile := filepath.Join(tempDir, "types.tsd")
	typesContent := `type Person(id: string, name: string, age:number)`
	err = os.WriteFile(typesFile, []byte(typesContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write types file: %v", err)
	}

	err = ps.ParseAndMerge(typesFile)
	if err != nil {
		t.Fatalf("Failed to parse types file: %v", err)
	}

	// Create and parse first facts file
	facts1File := filepath.Join(tempDir, "facts1.tsd")
	facts1Content := `
Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)
`
	err = os.WriteFile(facts1File, []byte(facts1Content), 0644)
	if err != nil {
		t.Fatalf("Failed to write facts1 file: %v", err)
	}

	err = ps.ParseAndMerge(facts1File)
	if err != nil {
		t.Fatalf("Failed to parse facts1 file: %v", err)
	}

	if len(ps.Facts) != 2 {
		t.Errorf("After facts1: expected 2 facts, got %d", len(ps.Facts))
	}

	// Create and parse second facts file
	facts2File := filepath.Join(tempDir, "facts2.tsd")
	facts2Content := `
Person(id: "P003", name: "Charlie", age: 35)
Person(id: "P004", name: "Diana", age: 28)
Person(id: "P005", name: "Eve", age: 42)
`
	err = os.WriteFile(facts2File, []byte(facts2Content), 0644)
	if err != nil {
		t.Fatalf("Failed to write facts2 file: %v", err)
	}

	err = ps.ParseAndMerge(facts2File)
	if err != nil {
		t.Fatalf("Failed to parse facts2 file: %v", err)
	}

	if len(ps.Facts) != 5 {
		t.Errorf("After facts2: expected 5 facts, got %d", len(ps.Facts))
	}

	// Create and parse third facts file
	facts3File := filepath.Join(tempDir, "facts3.tsd")
	facts3Content := `
Person(id: "P006", name: "Frank", age: 50)
Person(id: "P007", name: "Grace", age: 33)
`
	err = os.WriteFile(facts3File, []byte(facts3Content), 0644)
	if err != nil {
		t.Fatalf("Failed to write facts3 file: %v", err)
	}

	err = ps.ParseAndMerge(facts3File)
	if err != nil {
		t.Fatalf("Failed to parse facts3 file: %v", err)
	}

	if len(ps.Facts) != 7 {
		t.Errorf("After facts3: expected 7 facts, got %d", len(ps.Facts))
	}

	// Verify files parsed
	if len(ps.FilesParsed) != 4 {
		t.Errorf("Expected 4 files parsed, got %d", len(ps.FilesParsed))
	}

	// Verify no errors
	if ps.HasErrors() {
		t.Errorf("Unexpected validation errors: %v", ps.GetErrors())
	}
}

// TestIncrementalFactsParsing_WithInvalidFacts tests that invalid facts don't affect valid ones
func TestIncrementalFactsParsing_WithInvalidFacts(t *testing.T) {
	ps := NewProgramState()

	// Define type
	typeContent := `type Person(id: string, name: string, age:number)`
	err := ps.ParseAndMergeContent(typeContent, "types.tsd")
	if err != nil {
		t.Fatalf("Failed to parse type: %v", err)
	}

	// Add first batch of valid facts
	validFacts1 := `
Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)
`
	err = ps.ParseAndMergeContent(validFacts1, "valid1.tsd")
	if err != nil {
		t.Fatalf("Failed to parse valid facts: %v", err)
	}

	if len(ps.Facts) != 2 {
		t.Errorf("After valid batch 1: expected 2 facts, got %d", len(ps.Facts))
	}

	// Add batch with some invalid facts
	mixedFacts := `
Person(id: "P003", name: "Charlie", age: 35)
Person(id: "P004", name: "Diana", salary: 50000)
Person(id: "P005", name: "Eve", age: "forty-two")
`
	err = ps.ParseAndMergeContent(mixedFacts, "mixed.tsd")
	if err != nil {
		t.Fatalf("ParseAndMergeContent should not fail: %v", err)
	}

	// Should have 3 facts (2 from first batch + 1 valid from second)
	if len(ps.Facts) != 3 {
		t.Errorf("After mixed batch: expected 3 facts, got %d", len(ps.Facts))
	}

	// Should have 2 validation errors
	if !ps.HasErrors() {
		t.Error("Expected validation errors for invalid facts")
	}

	errorCount := ps.GetErrorCount()
	if errorCount != 2 {
		t.Errorf("Expected 2 errors, got %d", errorCount)
	}

	// Add another batch of valid facts
	validFacts2 := `
Person(id: "P006", name: "Frank", age: 50)
Person(id: "P007", name: "Grace", age: 33)
`
	err = ps.ParseAndMergeContent(validFacts2, "valid2.tsd")
	if err != nil {
		t.Fatalf("Failed to parse second valid batch: %v", err)
	}

	// Should now have 5 facts total
	if len(ps.Facts) != 5 {
		t.Errorf("After valid batch 2: expected 5 facts, got %d", len(ps.Facts))
	}

	// Error count should still be 2 (no new errors)
	if ps.GetErrorCount() != 2 {
		t.Errorf("Expected 2 total errors, got %d", ps.GetErrorCount())
	}
}

// TestIncrementalFactsParsing_WithReset tests that reset clears facts before adding new ones
func TestIncrementalFactsParsing_WithReset(t *testing.T) {
	ps := NewProgramState()

	// Define type and add facts
	content1 := `
type Person(id: string, name: string, age:number)

Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)
`
	err := ps.ParseAndMergeContent(content1, "initial.tsd")
	if err != nil {
		t.Fatalf("Failed to parse initial content: %v", err)
	}

	if len(ps.Facts) != 2 {
		t.Errorf("After initial: expected 2 facts, got %d", len(ps.Facts))
	}

	// Add more facts
	content2 := `
Person(id: "P003", name: "Charlie", age: 35)
`
	err = ps.ParseAndMergeContent(content2, "additional.tsd")
	if err != nil {
		t.Fatalf("Failed to parse additional content: %v", err)
	}

	if len(ps.Facts) != 3 {
		t.Errorf("After additional: expected 3 facts, got %d", len(ps.Facts))
	}

	// Reset and add new facts
	content3 := `
reset

type Person(id: string, name: string, age:number)

Person(id: "P100", name: "Xavier", age: 40)
Person(id: "P101", name: "Yolanda", age: 45)
`
	err = ps.ParseAndMergeContent(content3, "reset.tsd")
	if err != nil {
		t.Fatalf("Failed to parse reset content: %v", err)
	}

	// Should only have 2 facts (after reset)
	if len(ps.Facts) != 2 {
		t.Errorf("After reset: expected 2 facts, got %d", len(ps.Facts))
	}

	// Verify the facts are the new ones
	for _, fact := range ps.Facts {
		for _, field := range fact.Fields {
			if field.Name == "id" {
				id := field.Value.Value.(string)
				if id != "P100" && id != "P101" {
					t.Errorf("Found old fact ID after reset: %s", id)
				}
			}
		}
	}
}

// TestIncrementalFactsParsing_LargeScale tests parsing many facts incrementally
func TestIncrementalFactsParsing_LargeScale(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large scale test in short mode")
	}

	ps := NewProgramState()

	// Define type
	typeContent := `type Person(id: string, name: string, age:number)`
	err := ps.ParseAndMergeContent(typeContent, "types.tsd")
	if err != nil {
		t.Fatalf("Failed to parse type: %v", err)
	}

	// Add facts in batches of 10, for 100 total facts
	batchSize := 10
	totalBatches := 10

	for batchNum := 0; batchNum < totalBatches; batchNum++ {
		var factsContent string
		for i := 0; i < batchSize; i++ {
			factNum := batchNum*batchSize + i
			factsContent += "\n" + "Person(id: \"P" + padInt(factNum, 3) + "\", name: \"Person" + padInt(factNum, 3) + "\", age: " + intToString(20+factNum%50) + ")"
		}

		err = ps.ParseAndMergeContent(factsContent, "batch_"+intToString(batchNum)+".tsd")
		if err != nil {
			t.Fatalf("Failed to parse batch %d: %v", batchNum, err)
		}

		expectedCount := (batchNum + 1) * batchSize
		if len(ps.Facts) != expectedCount {
			t.Errorf("After batch %d: expected %d facts, got %d", batchNum, expectedCount, len(ps.Facts))
		}
	}

	// Verify final count
	if len(ps.Facts) != totalBatches*batchSize {
		t.Errorf("Final count: expected %d facts, got %d", totalBatches*batchSize, len(ps.Facts))
	}

	// Verify no errors
	if ps.HasErrors() {
		t.Errorf("Unexpected validation errors: %v", ps.GetErrors())
	}
}

// TestIncrementalFactsParsing_MixedWithRules tests adding facts incrementally while rules are present
func TestIncrementalFactsParsing_MixedWithRules(t *testing.T) {
	ps := NewProgramState()

	// Define type and rules
	typeAndRules := `
type Person(id: string, name: string, age:number)

rule adult_check : {p: Person} / p.age >= 18 ==> adult(p.id)
rule senior_check : {p: Person} / p.age >= 60 ==> senior(p.id)
`
	err := ps.ParseAndMergeContent(typeAndRules, "types_rules.tsd")
	if err != nil {
		t.Fatalf("Failed to parse types and rules: %v", err)
	}

	if len(ps.Rules) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(ps.Rules))
	}

	// Add facts batch 1
	facts1 := `
Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 17)
`
	err = ps.ParseAndMergeContent(facts1, "facts1.tsd")
	if err != nil {
		t.Fatalf("Failed to parse facts1: %v", err)
	}

	if len(ps.Facts) != 2 {
		t.Errorf("After facts1: expected 2 facts, got %d", len(ps.Facts))
	}

	// Rules should still be there
	if len(ps.Rules) != 2 {
		t.Errorf("Rules disappeared: expected 2, got %d", len(ps.Rules))
	}

	// Add facts batch 2
	facts2 := `
Person(id: "P003", name: "Charlie", age: 65)
Person(id: "P004", name: "Diana", age: 25)
`
	err = ps.ParseAndMergeContent(facts2, "facts2.tsd")
	if err != nil {
		t.Fatalf("Failed to parse facts2: %v", err)
	}

	if len(ps.Facts) != 4 {
		t.Errorf("After facts2: expected 4 facts, got %d", len(ps.Facts))
	}

	// Rules should still be there
	if len(ps.Rules) != 2 {
		t.Errorf("Rules disappeared: expected 2, got %d", len(ps.Rules))
	}

	// Types should still be there
	if len(ps.Types) != 1 {
		t.Errorf("Types disappeared: expected 1, got %d", len(ps.Types))
	}

	// Verify no errors
	if ps.HasErrors() {
		t.Errorf("Unexpected validation errors: %v", ps.GetErrors())
	}
}

// Helper functions
func padInt(n int, width int) string {
	s := intToString(n)
	for len(s) < width {
		s = "0" + s
	}
	return s
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + intToString(-n)
	}
	digits := ""
	for n > 0 {
		digits = string(rune('0'+n%10)) + digits
		n /= 10
	}
	return digits
}
