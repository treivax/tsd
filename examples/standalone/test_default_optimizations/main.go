// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Test rapide pour vérifier que les optimisations sont activées par défaut
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/treivax/tsd/rete"
)

func main() {
	fmt.Println("=== Test des optimisations par défaut ===")

	tmpDir, _ := os.MkdirTemp("", "test_opt_")
	defer os.RemoveAll(tmpDir)

	storage := rete.NewMemoryStorage()
	pipeline := rete.NewConstraintPipeline()

	// Test 1: Validation incrémentale
	fmt.Println("1️⃣  Test validation incrémentale")

	typesFile := filepath.Join(tmpDir, "types.tsd")
	os.WriteFile(typesFile, []byte(`type Person {
    id: string
    name: string
    age: number
}`), 0644)

	network, _, err := pipeline.IngestFile(typesFile, nil, storage)
	if err != nil {
		fmt.Printf("   ❌ Erreur: %v\n", err)
		return
	}
	fmt.Println("   ✅ Types chargés")

	// Règle avec type existant (doit passer avec validation incrémentale)
	rulesFile := filepath.Join(tmpDir, "rules.tsd")
	os.WriteFile(rulesFile, []byte(`rule check_person {
    when {
        p: Person(name == "John")
    }
    then {
        print("Found John")
    }
}`), 0644)

	network, _, err = pipeline.IngestFile(rulesFile, network, storage)
	if err != nil {
		fmt.Printf("   ❌ Erreur: %v\n", err)
		return
	}
	fmt.Println("   ✅ Validation incrémentale OK")

	// Règle avec type non existant (doit échouer)
	invalidFile := filepath.Join(tmpDir, "invalid.tsd")
	os.WriteFile(invalidFile, []byte(`rule check_company {
    when {
        c: Company(employees > 10)
    }
    then {
        print("Company")
    }
}`), 0644)

	_, _, err = pipeline.IngestFile(invalidFile, network, storage)
	if err != nil {
		fmt.Println("   ✅ Erreur détectée (comme attendu)")
	} else {
		fmt.Println("   ❌ Erreur NON détectée (problème)")
		return
	}

	// Test 2: GC après reset
	fmt.Println("\n2️⃣  Test GC après reset")

	nodesBefore := len(network.TypeNodes) + len(network.AlphaNodes) + len(network.TerminalNodes)
	fmt.Printf("   Nœuds avant reset: %d\n", nodesBefore)

	resetFile := filepath.Join(tmpDir, "reset.tsd")
	os.WriteFile(resetFile, []byte(`reset

type Vehicle {
    id: string
    brand: string
}`), 0644)

	network, _, err = pipeline.IngestFile(resetFile, network, storage)
	if err != nil {
		fmt.Printf("   ❌ Erreur: %v\n", err)
		return
	}

	nodesAfter := len(network.TypeNodes) + len(network.AlphaNodes) + len(network.TerminalNodes)
	fmt.Printf("   Nœuds après reset: %d\n", nodesAfter)

	if len(network.Types) == 1 && network.Types[0].Name == "Vehicle" {
		fmt.Println("   ✅ GC effectué correctement")
	} else {
		fmt.Println("   ❌ GC non effectué (problème)")
		return
	}

	fmt.Println("\n=== ✅ TOUS LES TESTS PASSENT ===")
	fmt.Println("Les optimisations sont bien activées par défaut !")
}
