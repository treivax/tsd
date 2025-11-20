// Package main provides the CLI for RETE validation testing
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/treivax/tsd/internal/validation"
	"github.com/treivax/tsd/pkg/testing"
)

func main() {
	fmt.Println("=== RETE VALIDATION CLI ===")
	fmt.Println("Validation authentique avec réseau RETE réel\n")

	if len(os.Args) == 3 {
		// Mode test spécifique
		constraintFile := os.Args[1]
		factsFile := os.Args[2]

		fmt.Printf("Test spécifique: %s + %s\n\n", constraintFile, factsFile)

		runner := testing.NewTestRunner("")
		result := runner.RunSingleTest(constraintFile, factsFile)

		displayTestResult(result)
		return
	}

	// Mode batch sur tous les tests
	testDir := "/home/resinsec/dev/tsd/beta_coverage_tests"

	runner := testing.NewTestRunner(testDir)
	results, err := runner.RunAllTests()
	if err != nil {
		fmt.Printf("Erreur exécution tests: %v\n", err)
		os.Exit(1)
	}

	generateSummaryReport(results)
}

// displayTestResult affiche le résultat d'un test unique
func displayTestResult(result validation.RETETestResult) {
	fmt.Printf("\n=== RÉSULTATS VALIDATION RETE ===\n")
	fmt.Printf("📋 Test: %s\n", result.TestName)
	fmt.Printf("⏱️  Durée: %v\n", result.ExecutionTime)

	if result.ValidationError != "" {
		fmt.Printf("❌ ERREUR: %s\n", result.ValidationError)
		return
	}

	fmt.Printf("\n📊 MÉTRIQUES:\n")
	fmt.Printf("  • Tokens attendus (simulation): %d\n", len(result.ExpectedTokens))
	fmt.Printf("  • Tokens observés (RETE réel): %d\n", len(result.ObservedTokens))
	fmt.Printf("  • Correspondances: %d\n", len(result.Matches))
	fmt.Printf("  • Mismatches: %d\n", result.Mismatches)
	fmt.Printf("  • Taux de succès: %.1f%%\n", result.SuccessRate)

	if result.IsValid {
		fmt.Printf("\n✅ TEST VALIDÉ\n")
	} else {
		fmt.Printf("\n❌ TEST INVALIDÉ: %s\n", result.ValidationError)
	}

	// Affichage détaillé des tokens
	if len(result.ObservedTokens) > 0 {
		fmt.Printf("\n🔍 TOKENS OBSERVÉS (RETE):\n")
		for i, token := range result.ObservedTokens {
			fmt.Printf("  %d. Règle: %s | Clé: %s\n", i+1, token.RuleName, token.Key)
			for factType, fact := range token.Facts {
				fmt.Printf("     └── %s: %s (ID: %s)\n", factType, formatFactValues(fact.Values), fact.ID)
			}
		}
	}

	if len(result.ExpectedTokens) > 0 {
		fmt.Printf("\n🎯 TOKENS ATTENDUS (simulation):\n")
		for i, token := range result.ExpectedTokens {
			fmt.Printf("  %d. Règle: %s | Clé: %s\n", i+1, token.RuleName, token.Key)
			for factType, fact := range token.Facts {
				fmt.Printf("     └── %s: %s (ID: %s)\n", factType, formatFactValues(fact.Values), fact.ID)
			}
		}
	}
}

// generateSummaryReport génère un rapport de synthèse
func generateSummaryReport(results []validation.RETETestResult) {
	successCount := 0
	for _, result := range results {
		if result.IsValid {
			successCount++
		}
	}

	fmt.Printf("\n=== RAPPORT DE SYNTHÈSE ===\n")
	fmt.Printf("📊 Tests totaux: %d\n", len(results))
	fmt.Printf("✅ Tests réussis: %d\n", successCount)
	fmt.Printf("❌ Tests échoués: %d\n", len(results)-successCount)
	fmt.Printf("📈 Taux de réussite: %.1f%%\n", float64(successCount)/float64(len(results))*100)
	fmt.Printf("🔥 Méthode: Réseau RETE authentique\n")
	fmt.Printf("📅 Date: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	if len(results)-successCount > 0 {
		fmt.Printf("\n❌ TESTS ÉCHOUÉS:\n")
		for _, result := range results {
			if !result.IsValid {
				fmt.Printf("  • %s: %s\n", result.TestName, result.ValidationError)
			}
		}
	}
}

// formatFactValues formate les valeurs d'un fait
func formatFactValues(values map[string]string) string {
	var parts []string
	for key, value := range values {
		parts = append(parts, fmt.Sprintf("%s:%s", key, value))
	}
	return strings.Join(parts, ", ")
}
