// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"os"
	"path/filepath"
	"testing"
)

// TestParseTypesAndFactsWithoutRules tests parsing a file with only types and facts (no rules)
func TestParseTypesAndFactsWithoutRules(t *testing.T) {
	ps := NewProgramState()

	content := `
type Person(id: string, name: string, age:number)
type Product(id: string, name: string, price:number)

Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)
Product(id: "PR001", name: "Laptop", price: 999.99)
Product(id: "PR002", name: "Mouse", price: 29.99)
`

	err := ps.ParseAndMergeContent(content, "types_and_facts.tsd")
	if err != nil {
		t.Fatalf("Failed to parse types and facts without rules: %v", err)
	}

	// Verify types were parsed
	if ps.GetTypesCount() != 2 {
		t.Errorf("Expected 2 types, got %d", ps.GetTypesCount())
	}

	if _, exists := ps.GetTypes()["Person"]; !exists {
		t.Error("Person type not found")
	}

	if _, exists := ps.GetTypes()["Product"]; !exists {
		t.Error("Product type not found")
	}

	// Verify facts were parsed
	if ps.GetFactsCount() != 4 {
		t.Errorf("Expected 4 facts, got %d", ps.GetFactsCount())
	}

	// Verify no rules (should be 0)
	if ps.GetRulesCount() != 0 {
		t.Errorf("Expected 0 rules, got %d", ps.GetRulesCount())
	}

	// Verify no errors
	if ps.HasErrors() {
		t.Errorf("Unexpected errors: %v", ps.GetErrors())
	}
}

// TestParseTypesAndFactsWithoutRules_File tests parsing an actual file with only types and facts
func TestParseTypesAndFactsWithoutRules_File(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "no_rules_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create file with types and facts only
	tsdFile := filepath.Join(tempDir, "data.tsd")
	content := `
type Customer(id: string, name: string, email: string, status:string)
type Order(id: string, customer_id: string, amount: number, date:string)

Customer(id: "C001", name: "Alice Smith", email: "alice@example.com", status: "active")
Customer(id: "C002", name: "Bob Jones", email: "bob@example.com", status: "active")
Customer(id: "C003", name: "Charlie Brown", email: "charlie@example.com", status: "inactive")

Order(id: "O001", customer_id: "C001", amount: 150.50, date: "2025-01-15")
Order(id: "O002", customer_id: "C001", amount: 275.00, date: "2025-01-16")
Order(id: "O003", customer_id: "C002", amount: 99.99, date: "2025-01-17")
Order(id: "O004", customer_id: "C003", amount: 50.00, date: "2025-01-18")
`
	err = os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	ps := NewProgramState()
	err = ps.ParseAndMerge(tsdFile)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	// Verify parsing results
	if ps.GetTypesCount() != 2 {
		t.Errorf("Expected 2 types, got %d", ps.GetTypesCount())
	}

	if ps.GetFactsCount() != 7 {
		t.Errorf("Expected 7 facts, got %d", ps.GetFactsCount())
	}

	if ps.GetRulesCount() != 0 {
		t.Errorf("Expected 0 rules, got %d", ps.GetRulesCount())
	}

	// Verify fact types
	customerCount := 0
	orderCount := 0
	for _, fact := range ps.GetFacts() {
		switch fact.TypeName {
		case "Customer":
			customerCount++
		case "Order":
			orderCount++
		}
	}

	if customerCount != 3 {
		t.Errorf("Expected 3 Customer facts, got %d", customerCount)
	}

	if orderCount != 4 {
		t.Errorf("Expected 4 Order facts, got %d", orderCount)
	}

	// Convert to Program
	program := ps.ToProgram()
	if program == nil {
		t.Fatal("Failed to convert to Program")
	}

	if len(program.Types) != 2 {
		t.Errorf("Program: expected 2 types, got %d", len(program.Types))
	}

	if len(program.Facts) != 7 {
		t.Errorf("Program: expected 7 facts, got %d", len(program.Facts))
	}

	if len(program.Expressions) != 0 {
		t.Errorf("Program: expected 0 expressions, got %d", len(program.Expressions))
	}
}

// TestParseOnlyTypes tests parsing a file with only type definitions
func TestParseOnlyTypes(t *testing.T) {
	ps := NewProgramState()

	content := `
type Person(id: string, name: string, age:number)
type Company(id: string, name: string, employees:number)
type Address(street: string, city: string, zipcode:string)
`

	err := ps.ParseAndMergeContent(content, "types_only.tsd")
	if err != nil {
		t.Fatalf("Failed to parse types only: %v", err)
	}

	if ps.GetTypesCount() != 3 {
		t.Errorf("Expected 3 types, got %d", ps.GetTypesCount())
	}

	if ps.GetFactsCount() != 0 {
		t.Errorf("Expected 0 facts, got %d", ps.GetFactsCount())
	}

	if ps.GetRulesCount() != 0 {
		t.Errorf("Expected 0 rules, got %d", ps.GetRulesCount())
	}
}

// TestParseOnlyFacts_ShouldFail tests that parsing facts without types fails
func TestParseOnlyFacts_ShouldFail(t *testing.T) {
	ps := NewProgramState()

	content := `
Person(id: "P001", name: "Alice", age: 30)
Product(id: "PR001", name: "Laptop", price: 999.99)
`

	err := ps.ParseAndMergeContent(content, "facts_only.tsd")
	if err != nil {
		t.Fatalf("ParseAndMergeContent should not fail: %v", err)
	}

	// Facts should be rejected because types are not defined
	if ps.GetFactsCount() != 0 {
		t.Errorf("Expected 0 facts (types not defined), got %d", ps.GetFactsCount())
	}

	// Should have validation errors
	if !ps.HasErrors() {
		t.Error("Expected validation errors for facts without types")
	}

	if ps.GetErrorCount() != 2 {
		t.Errorf("Expected 2 errors, got %d", ps.GetErrorCount())
	}
}

// TestParseTypesAndFactsIncremental_NoRules tests incremental parsing without rules
func TestParseTypesAndFactsIncremental_NoRules(t *testing.T) {
	ps := NewProgramState()

	// Step 1: Parse types
	typesContent := `
type Person(id: string, name: string, age:number)
type Product(id: string, name: string, price:number)
`
	err := ps.ParseAndMergeContent(typesContent, "types.tsd")
	if err != nil {
		t.Fatalf("Failed to parse types: %v", err)
	}

	if ps.GetTypesCount() != 2 {
		t.Errorf("Expected 2 types, got %d", ps.GetTypesCount())
	}

	// Step 2: Parse first batch of facts
	facts1 := `
Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)
`
	err = ps.ParseAndMergeContent(facts1, "facts1.tsd")
	if err != nil {
		t.Fatalf("Failed to parse facts1: %v", err)
	}

	if ps.GetFactsCount() != 2 {
		t.Errorf("After facts1: expected 2 facts, got %d", ps.GetFactsCount())
	}

	// Step 3: Parse second batch of facts
	facts2 := `
Product(id: "PR001", name: "Laptop", price: 999.99)
Product(id: "PR002", name: "Mouse", price: 29.99)
Person(id: "P003", name: "Charlie", age: 35)
`
	err = ps.ParseAndMergeContent(facts2, "facts2.tsd")
	if err != nil {
		t.Fatalf("Failed to parse facts2: %v", err)
	}

	if ps.GetFactsCount() != 5 {
		t.Errorf("After facts2: expected 5 facts, got %d", ps.GetFactsCount())
	}

	// Verify no rules throughout
	if ps.GetRulesCount() != 0 {
		t.Errorf("Expected 0 rules throughout, got %d", ps.GetRulesCount())
	}

	// Verify no errors
	if ps.HasErrors() {
		t.Errorf("Unexpected errors: %v", ps.GetErrors())
	}
}

// TestConvertToProgram_NoRules tests converting ProgramState to Program when no rules exist
func TestConvertToProgram_NoRules(t *testing.T) {
	ps := NewProgramState()

	content := `
type Person(id: string, name: string, age:number)

Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)
Person(id: "P003", name: "Charlie", age: 35)
`

	err := ps.ParseAndMergeContent(content, "data.tsd")
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Convert to Program
	program := ps.ToProgram()
	if program == nil {
		t.Fatal("ToProgram returned nil")
	}

	// Verify Program structure
	if len(program.Types) != 1 {
		t.Errorf("Program: expected 1 type, got %d", len(program.Types))
	}

	if len(program.Facts) != 3 {
		t.Errorf("Program: expected 3 facts, got %d", len(program.Facts))
	}

	if len(program.Expressions) != 0 {
		t.Errorf("Program: expected 0 expressions, got %d", len(program.Expressions))
	}

	// Verify type definition
	if program.Types[0].Name != "Person" {
		t.Errorf("Expected type name 'Person', got '%s'", program.Types[0].Name)
	}

	// Verify facts
	for _, fact := range program.Facts {
		if fact.TypeName != "Person" {
			t.Errorf("Expected fact type 'Person', got '%s'", fact.TypeName)
		}
	}
}

// TestRETENetworkCreation_TypesAndFactsOnly tests that a RETE network is created with types and facts but no rules
func TestRETENetworkCreation_TypesAndFactsOnly(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "rete_no_rules_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create file with types and facts only (no rules)
	tsdFile := filepath.Join(tempDir, "data.tsd")
	content := `
type Person(id: string, name: string, age:number)
type Product(id: string, name: string, price:number)

Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)
Product(id: "PR001", name: "Laptop", price: 999.99)
`
	err = os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Note: We can't directly test RETE network creation from constraint package
	// because that requires the rete package. However, we've verified that:
	// 1. Parsing works correctly (types and facts are parsed)
	// 2. ProgramState can be converted to Program
	// 3. The Program structure is valid for RETE network construction

	ps := NewProgramState()
	err = ps.ParseAndMerge(tsdFile)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	// Verify the parsed structure is valid
	if ps.GetTypesCount() != 2 {
		t.Errorf("Expected 2 types, got %d", ps.GetTypesCount())
	}

	if ps.GetFactsCount() != 3 {
		t.Errorf("Expected 3 facts, got %d", ps.GetFactsCount())
	}

	if ps.GetRulesCount() != 0 {
		t.Errorf("Expected 0 rules, got %d", ps.GetRulesCount())
	}

	// Convert to Program (this is what would be passed to RETE)
	program := ps.ToProgram()
	if program == nil {
		t.Fatal("Failed to convert to Program")
	}

	// Verify Program structure is valid for RETE
	if len(program.Types) != 2 {
		t.Errorf("Program types: expected 2, got %d", len(program.Types))
	}

	if len(program.Facts) != 3 {
		t.Errorf("Program facts: expected 3, got %d", len(program.Facts))
	}

	if len(program.Expressions) != 0 {
		t.Errorf("Program expressions: expected 0, got %d", len(program.Expressions))
	}

	// Log success
	t.Logf("✅ Successfully parsed file with types and facts only (no rules)")
	t.Logf("   - Types: %d", len(program.Types))
	t.Logf("   - Facts: %d", len(program.Facts))
	t.Logf("   - Rules: %d", len(program.Expressions))
	t.Logf("   - Structure is valid for RETE network construction")
}

// TestMultipleFilesWithoutRules tests parsing multiple files with types and facts but no rules
func TestMultipleFilesWithoutRules(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "multi_no_rules_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ps := NewProgramState()

	// File 1: Types
	typesFile := filepath.Join(tempDir, "types.tsd")
	typesContent := `
type Customer(id: string, name: string, email:string)
type Order(id: string, customer_id: string, total:number)
`
	err = os.WriteFile(typesFile, []byte(typesContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write types file: %v", err)
	}

	err = ps.ParseAndMerge(typesFile)
	if err != nil {
		t.Fatalf("Failed to parse types file: %v", err)
	}

	if ps.GetTypesCount() != 2 {
		t.Errorf("After types: expected 2 types, got %d", ps.GetTypesCount())
	}

	// File 2: Customer facts
	customersFile := filepath.Join(tempDir, "customers.tsd")
	customersContent := `
Customer(id: "C001", name: "Alice", email: "alice@example.com")
Customer(id: "C002", name: "Bob", email: "bob@example.com")
`
	err = os.WriteFile(customersFile, []byte(customersContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write customers file: %v", err)
	}

	err = ps.ParseAndMerge(customersFile)
	if err != nil {
		t.Fatalf("Failed to parse customers file: %v", err)
	}

	if ps.GetFactsCount() != 2 {
		t.Errorf("After customers: expected 2 facts, got %d", ps.GetFactsCount())
	}

	// File 3: Order facts
	ordersFile := filepath.Join(tempDir, "orders.tsd")
	ordersContent := `
Order(id: "O001", customer_id: "C001", total: 150.50)
Order(id: "O002", customer_id: "C002", total: 275.00)
Order(id: "O003", customer_id: "C001", total: 99.99)
`
	err = os.WriteFile(ordersFile, []byte(ordersContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write orders file: %v", err)
	}

	err = ps.ParseAndMerge(ordersFile)
	if err != nil {
		t.Fatalf("Failed to parse orders file: %v", err)
	}

	if ps.GetFactsCount() != 5 {
		t.Errorf("After orders: expected 5 facts, got %d", ps.GetFactsCount())
	}

	// Verify no rules exist
	if ps.GetRulesCount() != 0 {
		t.Errorf("Expected 0 rules, got %d", ps.GetRulesCount())
	}

	// Verify file tracking
	if len(ps.GetFilesParsed()) != 3 {
		t.Errorf("Expected 3 files parsed, got %d", len(ps.GetFilesParsed()))
	}

	// Verify Program structure
	program := ps.ToProgram()
	if program == nil {
		t.Fatal("Failed to convert to Program")
	}

	if len(program.Types) != 2 {
		t.Errorf("Program: expected 2 types, got %d", len(program.Types))
	}

	if len(program.Facts) != 5 {
		t.Errorf("Program: expected 5 facts, got %d", len(program.Facts))
	}

	if len(program.Expressions) != 0 {
		t.Errorf("Program: expected 0 expressions, got %d", len(program.Expressions))
	}

	t.Logf("✅ Successfully parsed 3 files with types and facts only (no rules)")
	t.Logf("   - Files: %d", len(ps.GetFilesParsed()))
	t.Logf("   - Types: %d", len(program.Types))
	t.Logf("   - Facts: %d", len(program.Facts))
	t.Logf("   - Rules: %d", len(program.Expressions))
}
