// Package main provides the CLI for RETE validation testing
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/treivax/tsd/internal/validation"
)

// TestResult represents the outcome of a single test execution in the RETE validation process.
// It captures essential metrics and status information for each test run, including:
// - TestName: the identifier of the executed test
// - ExecutionTime: the duration taken to complete the test
// - TokensGenerated: the number of tokens produced during test execution
// - ValidationError: any error message encountered during validation (empty if successful)
// - Success: a boolean flag indicating whether the test passed or failed
type TestResult struct {
	TestName        string
	ExecutionTime   time.Duration
	TokensGenerated int
	ValidationError string
	Success         bool
}

func main() {
	fmt.Println("=== RETE VALIDATION CLI ===")
	fmt.Println("Validation authentique avec réseau RETE réel\n")

	if len(os.Args) == 3 {
		// Mode test spécifique
		constraintFile := os.Args[1]
		factsFile := os.Args[2]

		fmt.Printf("Test spécifique: %s + %s\n\n", constraintFile, factsFile)
		result := runSingleTest(constraintFile, factsFile)
		displayTestResult(result)
		return
	}

	// Mode batch sur tous les tests
	testDir := "/home/resinsec/dev/tsd/beta_coverage_tests"
	results := runAllTests(testDir)
	generateSummaryReport(results)
}

// runSingleTest exécute un test unique
func runSingleTest(constraintFile, factsFile string) TestResult {
	start := time.Now()

	// Créer le réseau RETE
	network := validation.NewRETEValidationNetwork()

	// Charger les contraintes
	err := network.ParseConstraintFile(constraintFile)
	if err != nil {
		return TestResult{
			TestName:        filepath.Base(constraintFile),
			ExecutionTime:   time.Since(start),
			ValidationError: fmt.Sprintf("Erreur parsing contraintes: %v", err),
			Success:         false,
		}
	}

	// Charger les faits
	err = network.LoadFactsFile(factsFile)
	if err != nil {
		return TestResult{
			TestName:        filepath.Base(constraintFile),
			ExecutionTime:   time.Since(start),
			ValidationError: fmt.Sprintf("Erreur chargement faits: %v", err),
			Success:         false,
		}
	}

	// Obtenir les résultats
	tokensCount, _ := network.GetValidationResults()

	return TestResult{
		TestName:        filepath.Base(constraintFile),
		ExecutionTime:   time.Since(start),
		TokensGenerated: tokensCount,
		Success:         tokensCount > 0,
	}
}

// runAllTests exécute tous les tests dans un répertoire
func runAllTests(testDir string) []TestResult {
	var results []TestResult

	// Parcourir les fichiers .constraint
	files, err := filepath.Glob(filepath.Join(testDir, "*.constraint"))
	if err != nil {
		fmt.Printf("Erreur listing fichiers: %v\n", err)
		return results
	}

	for _, constraintFile := range files {
		baseName := strings.TrimSuffix(filepath.Base(constraintFile), ".constraint")
		factsFile := filepath.Join(testDir, baseName+".facts")

		if _, err := os.Stat(factsFile); os.IsNotExist(err) {
			continue // Skip if no corresponding .facts file
		}

		result := runSingleTest(constraintFile, factsFile)
		results = append(results, result)

		status := "❌ ÉCHEC"
		if result.Success {
			status = "✅ SUCCÈS"
		}
		fmt.Printf("%s %s (%v)\n", status, baseName, result.ExecutionTime)
	}

	return results
}

// displayTestResult affiche le résultat d'un test unique
func displayTestResult(result TestResult) {
	fmt.Printf("\n=== RÉSULTATS VALIDATION RETE ===\n")
	fmt.Printf("📋 Test: %s\n", result.TestName)
	fmt.Printf("⏱️  Durée: %v\n", result.ExecutionTime)

	if result.ValidationError != "" {
		fmt.Printf("❌ ERREUR: %s\n", result.ValidationError)
		return
	}

	fmt.Printf("\n📊 MÉTRIQUES:\n")
	fmt.Printf("  • Tokens générés: %d\n", result.TokensGenerated)

	if result.Success {
		fmt.Printf("✅ TEST RÉUSSI\n")
	} else {
		fmt.Printf("❌ TEST ÉCHOUÉ\n")
	}
}

// generateSummaryReport génère un rapport de synthèse
func generateSummaryReport(results []TestResult) {
	successCount := 0
	totalTime := time.Duration(0)

	fmt.Printf("\n=== RAPPORT DE SYNTHÈSE ===\n")
	fmt.Printf("Tests exécutés: %d\n", len(results))

	for _, result := range results {
		if result.Success {
			successCount++
		}
		totalTime += result.ExecutionTime
	}

	fmt.Printf("Tests réussis: %d\n", successCount)
	fmt.Printf("Tests échoués: %d\n", len(results)-successCount)
	fmt.Printf("Taux de succès: %.1f%%\n", float64(successCount)/float64(len(results))*100)
	fmt.Printf("Durée totale: %v\n", totalTime)
}
