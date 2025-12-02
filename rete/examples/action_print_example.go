// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"log"

	"github.com/treivax/tsd/rete"
)

func main() {
	fmt.Println("=== Exemple d'utilisation de l'action print ===")

	// Créer un réseau RETE
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)

	// Définir un type Person
	personType := rete.TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []rete.Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
			{Name: "age", Type: "number"},
			{Name: "city", Type: "string"},
		},
	}
	network.Types = append(network.Types, personType)

	fmt.Println("✅ Type Person défini")

	// Créer quelques faits
	facts := []*rete.Fact{
		{
			ID:   "person_1",
			Type: "Person",
			Fields: map[string]interface{}{
				"id":   "1",
				"name": "Alice",
				"age":  25.0,
				"city": "Paris",
			},
		},
		{
			ID:   "person_2",
			Type: "Person",
			Fields: map[string]interface{}{
				"id":   "2",
				"name": "Bob",
				"age":  30.0,
				"city": "Lyon",
			},
		},
		{
			ID:   "person_3",
			Type: "Person",
			Fields: map[string]interface{}{
				"id":   "3",
				"name": "Charlie",
				"age":  35.0,
				"city": "Paris",
			},
		},
	}

	// Exemple 1: Afficher une chaîne littérale
	fmt.Println("--- Exemple 1: Afficher une chaîne littérale ---")
	action1 := &rete.Action{
		Type: "action",
		Jobs: []rete.JobCall{
			{
				Type: "jobCall",
				Name: "print",
				Args: []interface{}{
					map[string]interface{}{
						"type":  "string",
						"value": "=== Bienvenue dans TSD ===",
					},
				},
			},
		},
	}

	token1 := &rete.Token{
		ID:       "token1",
		Facts:    []*rete.Fact{facts[0]},
		Bindings: map[string]*rete.Fact{},
	}

	err := network.ActionExecutor.ExecuteAction(action1, token1)
	if err != nil {
		log.Fatalf("Erreur: %v", err)
	}
	fmt.Println()

	// Exemple 2: Afficher un champ d'un fait
	fmt.Println("--- Exemple 2: Afficher un champ d'un fait ---")
	action2 := &rete.Action{
		Type: "action",
		Jobs: []rete.JobCall{
			{
				Type: "jobCall",
				Name: "print",
				Args: []interface{}{
					map[string]interface{}{
						"type":   "fieldAccess",
						"object": "p",
						"field":  "name",
					},
				},
			},
		},
	}

	token2 := &rete.Token{
		ID:    "token2",
		Facts: []*rete.Fact{facts[0]},
		Bindings: map[string]*rete.Fact{
			"p": facts[0],
		},
	}

	err = network.ActionExecutor.ExecuteAction(action2, token2)
	if err != nil {
		log.Fatalf("Erreur: %v", err)
	}
	fmt.Println()

	// Exemple 3: Afficher plusieurs champs avec des messages
	fmt.Println("--- Exemple 3: Afficher plusieurs informations ---")
	for i, fact := range facts {
		action3 := &rete.Action{
			Type: "action",
			Jobs: []rete.JobCall{
				{
					Type: "jobCall",
					Name: "print",
					Args: []interface{}{
						map[string]interface{}{
							"type":  "string",
							"value": fmt.Sprintf("Personne %d:", i+1),
						},
					},
				},
				{
					Type: "jobCall",
					Name: "print",
					Args: []interface{}{
						map[string]interface{}{
							"type":   "fieldAccess",
							"object": "p",
							"field":  "name",
						},
					},
				},
				{
					Type: "jobCall",
					Name: "print",
					Args: []interface{}{
						map[string]interface{}{
							"type":   "fieldAccess",
							"object": "p",
							"field":  "age",
						},
					},
				},
				{
					Type: "jobCall",
					Name: "print",
					Args: []interface{}{
						map[string]interface{}{
							"type":   "fieldAccess",
							"object": "p",
							"field":  "city",
						},
					},
				},
			},
		}

		token := &rete.Token{
			ID:    fmt.Sprintf("token_%d", i+3),
			Facts: []*rete.Fact{fact},
			Bindings: map[string]*rete.Fact{
				"p": fact,
			},
		}

		err = network.ActionExecutor.ExecuteAction(action3, token)
		if err != nil {
			log.Fatalf("Erreur: %v", err)
		}
		fmt.Println()
	}

	// Exemple 4: Afficher un fait complet
	fmt.Println("--- Exemple 4: Afficher un fait complet ---")
	action4 := &rete.Action{
		Type: "action",
		Jobs: []rete.JobCall{
			{
				Type: "jobCall",
				Name: "print",
				Args: []interface{}{
					map[string]interface{}{
						"type": "variable",
						"name": "p",
					},
				},
			},
		},
	}

	token4 := &rete.Token{
		ID:    "token_complete",
		Facts: []*rete.Fact{facts[1]},
		Bindings: map[string]*rete.Fact{
			"p": facts[1],
		},
	}

	err = network.ActionExecutor.ExecuteAction(action4, token4)
	if err != nil {
		log.Fatalf("Erreur: %v", err)
	}
	fmt.Println()

	// Exemple 5: Action non définie (sera loguée uniquement)
	fmt.Println("--- Exemple 5: Action non définie ---")
	action5 := &rete.Action{
		Type: "action",
		Jobs: []rete.JobCall{
			{
				Type: "jobCall",
				Name: "send_email",
				Args: []interface{}{
					map[string]interface{}{
						"type":   "fieldAccess",
						"object": "p",
						"field":  "name",
					},
				},
			},
		},
	}

	token5 := &rete.Token{
		ID:    "token5",
		Facts: []*rete.Fact{facts[2]},
		Bindings: map[string]*rete.Fact{
			"p": facts[2],
		},
	}

	err = network.ActionExecutor.ExecuteAction(action5, token5)
	if err != nil {
		log.Fatalf("Erreur: %v", err)
	}
	fmt.Println("ℹ️  L'action 'send_email' n'est pas définie, elle est juste loguée")
	fmt.Println()

	// Afficher les actions enregistrées
	fmt.Println("--- Actions enregistrées ---")
	names := network.ActionExecutor.GetRegistry().GetRegisteredNames()
	fmt.Printf("Nombre d'actions: %d\n", len(names))
	for _, name := range names {
		fmt.Printf("  • %s\n", name)
	}
	fmt.Println()

	fmt.Println("=== Exemple terminé avec succès ===")
}
