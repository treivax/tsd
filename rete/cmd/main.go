package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/treivax/tsd/rete"
)

func main() {
	fmt.Println("🔥 DÉMONSTRATION DU RÉSEAU RETE (VERSION SIMPLE)")
	fmt.Println("===============================================")

	// 1. Créer un programme RETE manuellement (sans parser constraint pour l'instant)
	fmt.Println("\n📋 ÉTAPE 1: Création du programme RETE")
	
	program := &rete.Program{
		Types: []rete.TypeDefinition{
			{
				Type: "typeDefinition",
				Name: "Personne",
				Fields: []rete.Field{
					{Name: "nom", Type: "string"},
					{Name: "age", Type: "number"},
					{Name: "adulte", Type: "bool"},
				},
			},
		},
		Expressions: []rete.Expression{
			{
				Type: "expression",
				Set: rete.Set{
					Type: "set",
					Variables: []rete.TypedVariable{
						{
							Type:     "typedVariable",
							Name:     "client",
							DataType: "Personne",
						},
					},
				},
				Constraints: map[string]interface{}{
					"type": "logicalExpr",
					"conditions": []string{"client.age = 25", "client.adulte = true"},
				},
				Action: &rete.Action{
					Type: "action",
					Job: rete.JobCall{
						Type: "jobCall",
						Name: "action",
						Args: []string{"client"},
					},
				},
			},
		},
	}

	fmt.Printf("✅ Programme créé avec %d type(s) et %d expression(s)\n", 
		len(program.Types), len(program.Expressions))

	// 2. Créer le storage
	fmt.Println("\n💾 ÉTAPE 2: Initialisation du storage")
	storage := rete.NewMemoryStorage()
	fmt.Printf("✅ Storage en mémoire initialisé\n")

	// 3. Créer et construire le réseau RETE
	fmt.Println("\n🏗️  ÉTAPE 3: Construction du réseau RETE")
	network := rete.NewReteNetwork(storage)
	
	err := network.LoadFromAST(program)
	if err != nil {
		log.Fatalf("Erreur construction réseau: %v", err)
	}

	// Afficher la structure
	network.PrintNetworkStructure()

	// 4. Créer et soumettre des faits de test
	fmt.Println("\n🔥 ÉTAPE 4: Soumission de faits")
	
	// Fait qui devrait déclencher l'action (age=25 ET adulte=true)
	fact1 := &rete.Fact{
		ID:        "personne_1",
		Type:      "Personne",
		Fields: map[string]interface{}{
			"nom":    "Alice",
			"age":    25.0,
			"adulte": true,
		},
		Timestamp: time.Now(),
	}

	fmt.Printf("Soumission du fait 1: %s\n", fact1.String())
	err = network.SubmitFact(fact1)
	if err != nil {
		log.Fatalf("Erreur soumission fait 1: %v", err)
	}

	// Fait qui ne devrait PAS déclencher l'action (age=17 ET adulte=false)
	fact2 := &rete.Fact{
		ID:        "personne_2", 
		Type:      "Personne",
		Fields: map[string]interface{}{
			"nom":    "Bob",
			"age":    17.0,
			"adulte": false,
		},
		Timestamp: time.Now(),
	}

	fmt.Printf("Soumission du fait 2: %s\n", fact2.String())
	err = network.SubmitFact(fact2)
	if err != nil {
		log.Fatalf("Erreur soumission fait 2: %v", err)
	}

	// Fait d'un autre type (devrait être ignoré)
	fact3 := &rete.Fact{
		ID:        "autre_1",
		Type:      "AutreType",
		Fields: map[string]interface{}{
			"nom": "Test",
		},
		Timestamp: time.Now(),
	}

	fmt.Printf("Soumission du fait 3 (autre type): %s\n", fact3.String())
	err = network.SubmitFact(fact3)
	if err != nil {
		log.Fatalf("Erreur soumission fait 3: %v", err)
	}

	// 5. Afficher l'état du réseau
	fmt.Println("\n📊 ÉTAPE 5: État final du réseau")
	state, err := network.GetNetworkState()
	if err != nil {
		log.Fatalf("Erreur récupération état: %v", err)
	}

	for nodeID, memory := range state {
		fmt.Printf("\n🎯 Nœud: %s\n", nodeID)
		fmt.Printf("   Faits en mémoire: %d\n", len(memory.Facts))
		fmt.Printf("   Tokens en mémoire: %d\n", len(memory.Tokens))
		
		if len(memory.Facts) > 0 {
			fmt.Printf("   Détail des faits:\n")
			for _, fact := range memory.Facts {
				factJSON, _ := json.MarshalIndent(fact, "     ", "  ")
				fmt.Printf("     %s\n", factJSON)
			}
		}

		if len(memory.Tokens) > 0 {
			fmt.Printf("   Détail des tokens:\n")
			for _, token := range memory.Tokens {
				tokenJSON, _ := json.MarshalIndent(token, "     ", "  ")
				fmt.Printf("     %s\n", tokenJSON)
			}
		}
	}

	fmt.Println("\n🎉 DÉMONSTRATION TERMINÉE!")
	fmt.Println("Le réseau RETE a traité les faits et déclenché les actions appropriées.")
	fmt.Println("========================")
}