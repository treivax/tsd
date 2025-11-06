package main

import (
	"fmt"
	"time"

	"github.com/treivax/tsd/rete"
)

func main() {
	fmt.Println("🚀 BENCHMARK DE PERFORMANCE RETE (SANS PERSISTANCE)")
	fmt.Println("==================================================")

	// Démonstration des conditions Alpha
	demonstrateAlphaConditions()

	fmt.Println("\n==================================================")
	fmt.Println("🚀 BENCHMARK DE PERFORMANCE PRINCIPAL")
	fmt.Println("==================================================")

	// Créer un programme simple
	program := &rete.Program{
		Types: []rete.TypeDefinition{
			{
				Type: "typeDefinition",
				Name: "Event",
				Fields: []rete.Field{
					{Name: "id", Type: "number"},
					{Name: "priority", Type: "number"},
					{Name: "active", Type: "bool"},
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
							Name:     "evt",
							DataType: "Event",
						},
					},
				},
				Constraints: map[string]interface{}{
					"type":     "binaryOperation",
					"operator": ">",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "evt",
						"field":  "priority",
					},
					"right": map[string]interface{}{
						"type":  "numberLiteral",
						"value": 5.0,
					},
				},
				Action: &rete.Action{
					Type: "action",
					Job: rete.JobCall{
						Type: "jobCall",
						Name: "process",
						Args: []string{"evt"},
					},
				},
			},
		},
	}

	// Créer le réseau (storage en mémoire seulement)
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)

	err := network.LoadFromAST(program)
	if err != nil {
		panic(err)
	}

	fmt.Printf("✅ Réseau RETE initialisé\n")

	// Test de performance avec différentes charges
	testCases := []struct {
		name      string
		numFacts  int
		batchSize int
	}{
		{"Petite charge", 1000, 100},
		{"Charge moyenne", 10000, 1000},
		{"Grande charge", 50000, 5000},
	}

	for _, tc := range testCases {
		fmt.Printf("\n🔥 Test: %s (%d faits par batch de %d)\n", tc.name, tc.numFacts, tc.batchSize)

		start := time.Now()

		for i := 0; i < tc.numFacts; i++ {
			fact := &rete.Fact{
				ID:   fmt.Sprintf("event_%d", i),
				Type: "Event",
				Fields: map[string]interface{}{
					"id":       float64(i),
					"priority": float64(i % 10),
					"active":   i%2 == 0,
				},
				Timestamp: time.Now(),
			}

			err := network.SubmitFact(fact)
			if err != nil {
				panic(err)
			}

			// Afficher progression pour grandes charges
			if tc.numFacts > 1000 && i%tc.batchSize == 0 && i > 0 {
				elapsed := time.Since(start)
				rate := float64(i) / elapsed.Seconds()
				fmt.Printf("   Progress: %d/%d faits (%.0f faits/sec)\n", i, tc.numFacts, rate)
			}
		}

		elapsed := time.Since(start)
		rate := float64(tc.numFacts) / elapsed.Seconds()

		fmt.Printf("   ✅ Terminé en %v\n", elapsed)
		fmt.Printf("   📊 Performance: %.0f faits/seconde\n", rate)
		fmt.Printf("   📊 Temps par fait: %.2f µs\n", float64(elapsed.Nanoseconds())/float64(tc.numFacts)/1000)

		// Statistiques du réseau
		state, _ := network.GetNetworkState()
		totalFacts := 0
		totalTokens := 0
		for _, memory := range state {
			totalFacts += len(memory.Facts)
			totalTokens += len(memory.Tokens)
		}
		fmt.Printf("   💾 État final: %d faits, %d tokens dans le réseau\n", totalFacts, totalTokens)
	}

	fmt.Println("\n🎯 BENCHMARK TERMINÉ!")
	fmt.Println("Performance optimisée sans persistance etcd")
}

func demonstrateAlphaConditions() {
	fmt.Println("🔬 DÉMONSTRATION DES CONDITIONS ALPHA")
	fmt.Println("=====================================")

	storage := rete.NewMemoryStorage()
	builder := rete.NewAlphaConditionBuilder()

	// Créer quelques nœuds Alpha avec différentes conditions
	conditions := map[string]interface{}{
		"Priorité élevée": builder.And(
			builder.FieldEquals("evt", "active", true),
			builder.FieldGreaterOrEqual("evt", "priority", 8),
		),
		"Priorité moyenne": builder.FieldRange("evt", "priority", 4, 7),
		"Événements critiques": builder.AndMultiple(
			builder.FieldEquals("evt", "active", true),
			builder.FieldGreaterThan("evt", "score", 90.0),
		),
	}

	alphaNodes := make(map[string]*rete.AlphaNode)
	for name, condition := range conditions {
		nodeId := fmt.Sprintf("alpha_%s", name)
		alphaNodes[name] = rete.NewAlphaNode(nodeId, condition, "evt", storage)
	}

	// Créer des faits de test
	testFacts := []*rete.Fact{
		{
			ID:   "event_1",
			Type: "Event",
			Fields: map[string]interface{}{
				"id":       1,
				"priority": 9,
				"active":   true,
				"score":    85.0,
			},
		},
		{
			ID:   "event_2",
			Type: "Event",
			Fields: map[string]interface{}{
				"id":       2,
				"priority": 5,
				"active":   true,
				"score":    75.0,
			},
		},
		{
			ID:   "event_3",
			Type: "Event",
			Fields: map[string]interface{}{
				"id":       3,
				"priority": 7,
				"active":   true,
				"score":    95.5,
			},
		},
	}

	fmt.Println("\n📊 Test des conditions sur des faits d'exemple:")

	for i, fact := range testFacts {
		fmt.Printf("\n🔹 Fait %d: id=%v, priority=%v, active=%v, score=%v\n",
			i+1, fact.Fields["id"], fact.Fields["priority"],
			fact.Fields["active"], fact.Fields["score"])

		for name, node := range alphaNodes {
			memoryBefore := len(node.GetMemory().Facts)
			err := node.ActivateRight(fact)
			if err != nil {
				fmt.Printf("   ❌ %s: ERREUR (%v)\n", name, err)
				continue
			}
			memoryAfter := len(node.GetMemory().Facts)

			if memoryAfter > memoryBefore {
				fmt.Printf("   ✅ %s: MATCH\n", name)
			} else {
				fmt.Printf("   ❌ %s: NO MATCH\n", name)
			}
		}
	}

	fmt.Println("\n📈 Résumé des correspondances:")
	for name, node := range alphaNodes {
		fmt.Printf("   🎯 %s: %d faits correspondants\n", name, len(node.GetMemory().Facts))
	}
}
