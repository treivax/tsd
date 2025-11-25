package constraint

import (
	"os"
	"testing"
)

func TestParseRemoveFact(t *testing.T) {
	// Test du parsing d'une commande remove simple
	input := `
type Person : <id:string, name:string>

Person(id:P1, name:Alice)
remove Person P1
`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	parsed, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected result to be a map")
	}

	// Vérifier les types
	types := parsed["types"].([]interface{})
	if len(types) != 1 {
		t.Errorf("Expected 1 type, got %d", len(types))
	}

	// Vérifier les faits
	facts := parsed["facts"].([]interface{})
	if len(facts) != 1 {
		t.Errorf("Expected 1 fact, got %d", len(facts))
	}

	// Vérifier les rétractions
	retractions := parsed["retractions"].([]interface{})
	if len(retractions) != 1 {
		t.Fatalf("Expected 1 retraction, got %d", len(retractions))
	}

	retraction := retractions[0].(map[string]interface{})
	if retraction["type"] != "retraction" {
		t.Errorf("Expected type 'retraction', got '%s'", retraction["type"])
	}
	if retraction["typeName"] != "Person" {
		t.Errorf("Expected typeName 'Person', got '%s'", retraction["typeName"])
	}
	if retraction["factID"] != "P1" {
		t.Errorf("Expected factID 'P1', got '%s'", retraction["factID"])
	}
}

func TestParseMultipleRemoveFacts(t *testing.T) {
	// Test du parsing de plusieurs commandes remove
	input := `
type Person : <id:string, name:string>
type Order : <id:string, customer_id:string>

Person(id:P1, name:Alice)
Person(id:P2, name:Bob)
Order(id:O1, customer_id:P1)
Order(id:O2, customer_id:P2)

remove Person P1
remove Order O2
remove Person P2
`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	parsed := result.(map[string]interface{})

	// Vérifier les rétractions
	retractions := parsed["retractions"].([]interface{})
	if len(retractions) != 3 {
		t.Fatalf("Expected 3 retractions, got %d", len(retractions))
	}

	// Vérifier la première rétraction
	retraction1 := retractions[0].(map[string]interface{})
	if retraction1["typeName"] != "Person" || retraction1["factID"] != "P1" {
		t.Errorf("First retraction incorrect: %v", retraction1)
	}

	// Vérifier la deuxième rétraction
	retraction2 := retractions[1].(map[string]interface{})
	if retraction2["typeName"] != "Order" || retraction2["factID"] != "O2" {
		t.Errorf("Second retraction incorrect: %v", retraction2)
	}

	// Vérifier la troisième rétraction
	retraction3 := retractions[2].(map[string]interface{})
	if retraction3["typeName"] != "Person" || retraction3["factID"] != "P2" {
		t.Errorf("Third retraction incorrect: %v", retraction3)
	}
}

func TestParseRemoveFactFromFile(t *testing.T) {
	// Test du parsing du fichier de test
	content, err := os.ReadFile("test/remove_fact_test.constraint")
	if err != nil {
		t.Skipf("Could not read test file: %v", err)
		return
	}

	result, err := Parse("test/remove_fact_test.constraint", content)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	parsed := result.(map[string]interface{})

	// Vérifier les types
	types := parsed["types"].([]interface{})
	if len(types) != 2 {
		t.Errorf("Expected 2 types, got %d", len(types))
	}

	// Vérifier les faits (4 ajouts + 1 après rétractions = 5)
	facts := parsed["facts"].([]interface{})
	if len(facts) != 5 {
		t.Errorf("Expected 5 facts, got %d", len(facts))
	}

	// Vérifier les rétractions (2)
	retractions := parsed["retractions"].([]interface{})
	if len(retractions) != 2 {
		t.Errorf("Expected 2 retractions, got %d", len(retractions))
	}

	// Vérifier le contenu des rétractions
	retraction1 := retractions[0].(map[string]interface{})
	if retraction1["typeName"] != "Person" || retraction1["factID"] != "P1" {
		t.Errorf("First retraction incorrect: typeName=%s, factID=%s",
			retraction1["typeName"], retraction1["factID"])
	}

	retraction2 := retractions[1].(map[string]interface{})
	if retraction2["typeName"] != "Order" || retraction2["factID"] != "O2" {
		t.Errorf("Second retraction incorrect: typeName=%s, factID=%s",
			retraction2["typeName"], retraction2["factID"])
	}
}

func TestRemoveFactWithComplexID(t *testing.T) {
	// Test avec des IDs contenant des caractères spéciaux
	input := `
type Product : <id:string, name:string>

Product(id:PROD-123-ABC, name:Widget)
remove Product PROD-123-ABC
`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	parsed := result.(map[string]interface{})
	retractions := parsed["retractions"].([]interface{})

	if len(retractions) != 1 {
		t.Fatalf("Expected 1 retraction, got %d", len(retractions))
	}

	retraction := retractions[0].(map[string]interface{})
	if retraction["factID"] != "PROD-123-ABC" {
		t.Errorf("Expected factID 'PROD-123-ABC', got '%s'", retraction["factID"])
	}
}
