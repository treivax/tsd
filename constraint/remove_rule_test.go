// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"testing"
)

// TestParseRemoveFactNewSyntax vérifie que la nouvelle syntaxe "remove fact" fonctionne
func TestParseRemoveFactNewSyntax(t *testing.T) {
	t.Log("🧪 TEST PARSING REMOVE FACT - NOUVELLE SYNTAXE")
	t.Log("===============================================")

	input := `
type Person(id: string, name:string)

Person(id: "P1", name: "Alice")

remove fact Person P1
`

	result, err := ParseConstraint("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur parsing: %v", err)
	}

	parsed := result.(map[string]interface{})

	// Vérifier les faits
	facts := parsed["facts"].([]interface{})
	if len(facts) != 1 {
		t.Fatalf("Expected 1 fact, got %d", len(facts))
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

	t.Log("✅ Test réussi - 'remove fact' parse correctement")
}

// TestParseRemoveRule vérifie que la nouvelle commande "remove rule" fonctionne
func TestParseRemoveRule(t *testing.T) {
	t.Log("🧪 TEST PARSING REMOVE RULE")
	t.Log("===========================")

	input := `
type Person(id: string, name: string, age:number)

rule adult_check : {p: Person} / p.age >= 18 ==> notify(p.id)

remove rule adult_check
`

	result, err := ParseConstraint("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur parsing: %v", err)
	}

	parsed := result.(map[string]interface{})

	// Vérifier les règles
	expressions := parsed["expressions"].([]interface{})
	if len(expressions) != 1 {
		t.Fatalf("Expected 1 expression, got %d", len(expressions))
	}

	// Vérifier les suppressions de règles
	ruleRemovals := parsed["ruleRemovals"].([]interface{})
	if len(ruleRemovals) != 1 {
		t.Fatalf("Expected 1 rule removal, got %d", len(ruleRemovals))
	}

	removal := ruleRemovals[0].(map[string]interface{})
	if removal["type"] != "ruleRemoval" {
		t.Errorf("Expected type 'ruleRemoval', got '%s'", removal["type"])
	}
	if removal["ruleID"] != "adult_check" {
		t.Errorf("Expected ruleID 'adult_check', got '%s'", removal["ruleID"])
	}

	t.Log("✅ Test réussi - 'remove rule' parse correctement")
}

// TestParseMultipleRemoveCommandsMixed vérifie plusieurs commandes de suppression
func TestParseMultipleRemoveCommandsMixed(t *testing.T) {
	t.Log("🧪 TEST PARSING MULTIPLE REMOVE COMMANDS (MIXED)")
	t.Log("=================================================")

	input := `
type Person(id: string, name: string, age:number)
type Order(id: string, customer:string)

rule adult_check : {p: Person} / p.age >= 18 ==> notify(p.id)
rule order_check : {o: Order} / o.customer == "VIP" ==> process(o.id)

Person(id: "P1", name: "Alice", age: 25)
Order(id: "O1", customer: "VIP")

remove fact Person P1
remove fact Order O1
remove rule adult_check
remove rule order_check
`

	result, err := ParseConstraint("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur parsing: %v", err)
	}

	parsed := result.(map[string]interface{})

	// Vérifier les types
	types := parsed["types"].([]interface{})
	if len(types) != 2 {
		t.Errorf("Expected 2 types, got %d", len(types))
	}

	// Vérifier les règles
	expressions := parsed["expressions"].([]interface{})
	if len(expressions) != 2 {
		t.Errorf("Expected 2 expressions, got %d", len(expressions))
	}

	// Vérifier les faits
	facts := parsed["facts"].([]interface{})
	if len(facts) != 2 {
		t.Errorf("Expected 2 facts, got %d", len(facts))
	}

	// Vérifier les rétractions de faits
	retractions := parsed["retractions"].([]interface{})
	if len(retractions) != 2 {
		t.Fatalf("Expected 2 retractions, got %d", len(retractions))
	}

	retraction1 := retractions[0].(map[string]interface{})
	if retraction1["typeName"] != "Person" || retraction1["factID"] != "P1" {
		t.Errorf("First retraction incorrect: %v", retraction1)
	}

	retraction2 := retractions[1].(map[string]interface{})
	if retraction2["typeName"] != "Order" || retraction2["factID"] != "O1" {
		t.Errorf("Second retraction incorrect: %v", retraction2)
	}

	// Vérifier les suppressions de règles
	ruleRemovals := parsed["ruleRemovals"].([]interface{})
	if len(ruleRemovals) != 2 {
		t.Fatalf("Expected 2 rule removals, got %d", len(ruleRemovals))
	}

	removal1 := ruleRemovals[0].(map[string]interface{})
	if removal1["ruleID"] != "adult_check" {
		t.Errorf("First rule removal incorrect: expected 'adult_check', got '%s'", removal1["ruleID"])
	}

	removal2 := ruleRemovals[1].(map[string]interface{})
	if removal2["ruleID"] != "order_check" {
		t.Errorf("Second rule removal incorrect: expected 'order_check', got '%s'", removal2["ruleID"])
	}

	t.Log("✅ Test réussi - Multiples commandes de suppression parsent correctement")
}

// TestParseRemoveRuleWithComplexID vérifie les identifiants complexes
func TestParseRemoveRuleWithComplexID(t *testing.T) {
	t.Log("🧪 TEST PARSING REMOVE RULE WITH COMPLEX ID")
	t.Log("===========================================")

	input := `
type Person(id: string, name:string)

rule my_complex_rule_123 : {p: Person} / p.name == "Test" ==> action(p.id)

remove rule my_complex_rule_123
`

	result, err := ParseConstraint("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur parsing: %v", err)
	}

	parsed := result.(map[string]interface{})
	ruleRemovals := parsed["ruleRemovals"].([]interface{})

	if len(ruleRemovals) != 1 {
		t.Fatalf("Expected 1 rule removal, got %d", len(ruleRemovals))
	}

	removal := ruleRemovals[0].(map[string]interface{})
	if removal["ruleID"] != "my_complex_rule_123" {
		t.Errorf("Expected ruleID 'my_complex_rule_123', got '%s'", removal["ruleID"])
	}

	t.Log("✅ Test réussi - Identifiant complexe de règle parse correctement")
}

// TestParseRemoveRuleFromFile vérifie le parsing depuis un fichier
func TestParseRemoveRuleFromFile(t *testing.T) {
	t.Log("🧪 TEST PARSING REMOVE RULE FROM FILE")
	t.Log("=====================================")

	// Créer un fichier temporaire
	content := `
type Person(id: string, name: string, age:number)

rule rule1 : {p: Person} / p.age > 18 ==> action1(p.id)
rule rule2 : {p: Person} / p.age < 65 ==> action2(p.id)

remove rule rule1
remove rule rule2
`

	result, err := ParseConstraint("test.constraint", []byte(content))
	if err != nil {
		t.Fatalf("❌ Erreur parsing fichier: %v", err)
	}

	parsed := result.(map[string]interface{})

	// Vérifier les suppressions de règles (2)
	ruleRemovals := parsed["ruleRemovals"].([]interface{})
	if len(ruleRemovals) != 2 {
		t.Errorf("Expected 2 rule removals, got %d", len(ruleRemovals))
	}

	// Vérifier le contenu des suppressions
	removal1 := ruleRemovals[0].(map[string]interface{})
	if removal1["ruleID"] != "rule1" {
		t.Errorf("First removal incorrect: ruleID=%s", removal1["ruleID"])
	}

	removal2 := ruleRemovals[1].(map[string]interface{})
	if removal2["ruleID"] != "rule2" {
		t.Errorf("Second removal incorrect: ruleID=%s", removal2["ruleID"])
	}

	t.Log("✅ Test réussi - Parsing depuis fichier fonctionne")
}

// TestOldRemoveSyntaxShouldFail vérifie que l'ancienne syntaxe échoue maintenant
func TestOldRemoveSyntaxShouldFail(t *testing.T) {
	t.Log("🧪 TEST OLD REMOVE SYNTAX SHOULD FAIL")
	t.Log("======================================")

	// L'ancienne syntaxe "remove TypeName ID" ne devrait plus fonctionner
	input := `
type Person(id: string, name:string)

Person(id: "P1", name: "Alice")

remove Person P1
`

	_, err := ParseConstraint("test", []byte(input))
	if err == nil {
		t.Errorf("⚠️ L'ancienne syntaxe 'remove TypeName ID' devrait échouer maintenant")
	} else {
		t.Logf("✅ Comme attendu, l'ancienne syntaxe échoue: %v", err)
	}
}
