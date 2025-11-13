package main

import (
	"fmt"
	"log"

	"github.com/treivax/tsd/rete"
)

func main() {
	fmt.Println("🧪 TEST ARGUMENTS VARIABLES DANS ACTIONS")
	fmt.Println("========================================")

	// Créer le pipeline
	pipeline := rete.NewConstraintPipeline()

	// Créer le storage en mémoire
	storage := rete.NewMemoryStorage()

	// Construire le réseau avec fichier de contraintes et faits
	network, facts, err := pipeline.BuildNetworkFromConstraintFileWithFacts(
		"constraint/test/integration/variable_action_test.constraint",
		"constraint/test/integration/variable_action_test.facts",
		storage,
	)
	if err != nil {
		log.Fatalf("❌ Erreur construction réseau: %v", err)
	}

	fmt.Printf("✅ Réseau construit avec %d faits injectés\n", len(facts))

	// Récupérer l'état du réseau
	state, _ := network.GetNetworkState()
	totalTokens := 0
	for _, memory := range state {
		totalTokens += len(memory.Tokens)
	}

	// Afficher les nœuds terminaux
	fmt.Printf("\n🎯 Nœuds terminaux: %d\n", len(network.TerminalNodes))
	for id, terminal := range network.TerminalNodes {
		if terminal.Action != nil {
			fmt.Printf("   - %s: %s(%v)\n", id, terminal.Action.Job.Name, terminal.Action.Job.Args)
		}
	}

	fmt.Printf("\n📊 État final du réseau: %d tokens\n", totalTokens)
	fmt.Println("🎯 TEST TERMINÉ!")
}
