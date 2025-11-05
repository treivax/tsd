package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/treivax/tsd/constraint"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go <input-file>")
		fmt.Println("")
		fmt.Println("Exemple:")
		fmt.Println("  go run main.go ../tests/test_input.txt")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	input, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Erreur lecture fichier: %v", err)
	}

	// Parse l'input (la fonction Parse sera générée par pigeon)
	result, err := constraint.ParseConstraint("", input)
	if err != nil {
		log.Fatalf("Erreur parsing: %v", err)
	}

	// Validation du programme
	err = constraint.ValidateConstraintProgram(result)
	if err != nil {
		log.Fatalf("Erreur de validation: %v", err)
	}

	// Convertir en JSON pour affichage
	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Erreur JSON: %v", err)
	}

	fmt.Println(string(jsonOutput))
}
