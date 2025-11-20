package main

import (
	"fmt"

	"github.com/treivax/tsd/internal/validation"
)

func main() {
	fmt.Println("=== DIAGNOSTIC RETE ===")

	// Test simple
	rules := []string{
		"{u: User} / u.age > 18 ==> adult(u.id)",
	}

	facts := []string{
		"User(id:\"U1\", age:25)",
	}

	// Parser les règles
	types, reteRules := validation.ParseRulesForRETENetwork(rules)
	fmt.Printf("Types extraits: %v\n", types)
	fmt.Printf("Règles parsées: %v\n", reteRules)

	// Convertir les faits
	factData := validation.ConvertFactsToRETEData(facts)
	fmt.Printf("Faits convertis: %v\n", factData)

	// Créer le réseau
	network := validation.NewMiniRETENetwork()

	// Ajouter les types
	for _, typeName := range types {
		network.AddTypeNode(typeName)
	}

	// Ajouter les règles
	for _, rule := range reteRules {
		network.AddRuleNode(rule)
	}

	// Insérer les faits
	for _, fact := range factData {
		network.InsertFact(fact)
	}

	// Extraire les tokens
	tokens := network.ExtractTerminalTokens()
	fmt.Printf("Tokens extraits: %v\n", tokens)
}
