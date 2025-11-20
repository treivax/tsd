// Package main provides a clean RETE validation runner
package main

import (
	"fmt"
	"os"

	"github.com/treivax/tsd/internal/validation"
	"github.com/treivax/tsd/pkg/testing"
)

func main() {
	fmt.Println("=== RUNNER RETE VALIDATION PROPRE ===")
	fmt.Println("Architecture refactorisée selon les bonnes pratiques Go\n")

	if len(os.Args) == 3 {
		// Mode test spécifique
		constraintFile := os.Args[1]
		factsFile := os.Args[2]

		fmt.Printf("Test spécifique: %s + %s\n\n", constraintFile, factsFile)

		runner := testing.NewTestRunner("")
		result := runner.RunSingleTest(constraintFile, factsFile)

		displayResult(result)
		return
	}

	// Afficher l'aide
	fmt.Println("Usage:")
	fmt.Println("  go run runner_clean.go <constraint_file> <facts_file>")
	fmt.Println("")
	fmt.Println("Exemples:")
	fmt.Println("  go run runner_clean.go /path/to/join_simple.constraint /path/to/join_simple.facts")
	fmt.Println("")
	fmt.Println("Pour exécuter tous les tests:")
	fmt.Println("  go run ../../../cmd/rete-validate/main.go")
}

func displayResult(result validation.RETETestResult) {
	fmt.Printf("=== RÉSULTATS VALIDATION RETE ===\n")
	fmt.Printf("📋 Test: %s\n", result.TestName)
	fmt.Printf("⏱️  Durée: %v\n", result.ExecutionTime)

	if result.ValidationError != "" {
		fmt.Printf("❌ ERREUR: %s\n", result.ValidationError)
		return
	}

	fmt.Printf("\n📊 MÉTRIQUES:\n")
	fmt.Printf("  • Tokens attendus: %d\n", len(result.ExpectedTokens))
	fmt.Printf("  • Tokens observés: %d\n", len(result.ObservedTokens))
	fmt.Printf("  • Correspondances: %d\n", len(result.Matches))
	fmt.Printf("  • Taux de succès: %.1f%%\n", result.SuccessRate)

	if result.IsValid {
		fmt.Printf("\n✅ TEST VALIDÉ\n")
	} else {
		fmt.Printf("\n❌ TEST INVALIDÉ: %s\n", result.ValidationError)
	}
}
