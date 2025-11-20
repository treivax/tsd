package testing

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/treivax/tsd/internal/validation"
)

// TestRunner gère l'exécution des tests de validation RETE
type TestRunner struct {
	TestDirectory   string
	OutputDir       string
	Timeout         time.Duration
	Verbose         bool
	CreateReports   bool
	IncludePatterns []string
	ExcludePatterns []string
	results         []validation.RETETestResult
}

// NewTestRunner crée une nouvelle instance de TestRunner
func NewTestRunner(testDir string) *TestRunner {
	return &TestRunner{
		TestDirectory:   testDir,
		OutputDir:       "test_reports",
		Timeout:         5 * time.Minute,
		Verbose:         false,
		CreateReports:   true,
		IncludePatterns: []string{},
		ExcludePatterns: []string{},
		results:         []validation.RETETestResult{},
	}
}

// SetVerbose active ou désactive le mode verbose
func (tr *TestRunner) SetVerbose(verbose bool) {
	tr.Verbose = verbose
}

// SetTimeout définit le timeout pour les tests
func (tr *TestRunner) SetTimeout(timeout time.Duration) {
	tr.Timeout = timeout
}

// RunAllTests exécute tous les tests trouvés dans le répertoire
func (tr *TestRunner) RunAllTests() error {
	fmt.Printf("🚀 Démarrage des tests RETE dans: %s\n", tr.TestDirectory)

	// Découvrir tous les tests
	testPairs, err := tr.discoverTests()
	if err != nil {
		return fmt.Errorf("erreur lors de la découverte des tests: %v", err)
	}

	fmt.Printf("📂 Tests découverts: %d paires\n", len(testPairs))

	// Exécuter tous les tests
	for i, pair := range testPairs {
		fmt.Printf("\n📝 Test %d/%d: %s\n", i+1, len(testPairs), pair.Name)

		result := tr.executeSpecificTest(pair.ConstraintFile, pair.FactsFile)
		tr.results = append(tr.results, result)

		tr.displayTestResult(result)
	}

	// Afficher un résumé final
	tr.displayFinalSummary()

	return nil
}

// TestPair représente une paire de fichiers de test
type TestPair struct {
	Name           string
	ConstraintFile string
	FactsFile      string
}

// discoverTests découvre tous les tests dans le répertoire
func (tr *TestRunner) discoverTests() ([]TestPair, error) {
	var testPairs []TestPair

	err := filepath.Walk(tr.TestDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Chercher les fichiers .constraint
		if strings.HasSuffix(path, ".constraint") {
			baseName := strings.TrimSuffix(path, ".constraint")
			factsFile := baseName + ".facts"

			// Vérifier que le fichier .facts existe
			if _, err := os.Stat(factsFile); err == nil {
				testName := filepath.Base(baseName)

				// Appliquer les filtres d'inclusion/exclusion
				if tr.shouldIncludeTest(testName) {
					testPairs = append(testPairs, TestPair{
						Name:           testName,
						ConstraintFile: path,
						FactsFile:      factsFile,
					})
				}
			}
		}

		return nil
	})

	return testPairs, err
}

// shouldIncludeTest vérifie si un test doit être inclus selon les patterns
func (tr *TestRunner) shouldIncludeTest(testName string) bool {
	// Vérifier les patterns d'exclusion
	for _, exclude := range tr.ExcludePatterns {
		if strings.Contains(testName, exclude) {
			return false
		}
	}

	// Si pas de patterns d'inclusion, inclure par défaut
	if len(tr.IncludePatterns) == 0 {
		return true
	}

	// Vérifier les patterns d'inclusion
	for _, include := range tr.IncludePatterns {
		if strings.Contains(testName, include) {
			return true
		}
	}

	return false
}

// executeSpecificTest exécute un test spécifique avec des fichiers donnés
func (tr *TestRunner) executeSpecificTest(constraintFile, factsFile string) validation.RETETestResult {
	startTime := time.Now()

	testPath := strings.TrimSuffix(constraintFile, ".constraint")

	if tr.Verbose {
		fmt.Printf("📂 Fichiers: %s.{constraint,facts}\n", testPath)
	}

	// Utiliser la nouvelle fonction de validation RETE complète
	result, err := validation.ValidateRETEWithFile(testPath, tr.Timeout)
	if err != nil {
		return validation.RETETestResult{
			TestName:       filepath.Base(constraintFile),
			Rules:          []string{},
			Facts:          []string{},
			ObservedTokens: []validation.RETETokenInfo{},
			Success:        false,
			ValidationNote: fmt.Sprintf("Erreur: %v", err),
		}
	}

	if tr.Verbose {
		fmt.Printf("⏱️  Temps d'exécution: %v\n", time.Since(startTime))
	}

	return *result
}

// displayTestResult affiche les résultats d'un test
func (tr *TestRunner) displayTestResult(result validation.RETETestResult) {
	fmt.Printf("📊 Résultats test: %s\n", result.TestName)
	fmt.Printf("  Tokens observés RETE: %d\n", len(result.ObservedTokens))

	if result.Success {
		fmt.Printf("  ✅ VALIDÉ: %s\n", result.ValidationNote)
	} else {
		fmt.Printf("  ❌ INVALIDÉ: %s\n", result.ValidationNote)
	}
}

// displayFinalSummary affiche un résumé final de tous les tests
func (tr *TestRunner) displayFinalSummary() {
	fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Printf("📈 RÉSUMÉ FINAL - VALIDATION RETE COMPLÈTE\n")
	fmt.Printf(strings.Repeat("=", 80) + "\n")

	totalTests := len(tr.results)
	successfulTests := 0
	totalTokens := 0

	for _, result := range tr.results {
		if result.Success {
			successfulTests++
		}
		totalTokens += len(result.ObservedTokens)
	}

	fmt.Printf("Tests exécutés: %d\n", totalTests)
	fmt.Printf("Tests réussis: %d\n", successfulTests)
	fmt.Printf("Taux de réussite: %.1f%%\n", float64(successfulTests)/float64(totalTests)*100)
	fmt.Printf("Tokens RETE générés: %d\n", totalTokens)

	if successfulTests == totalTests {
		fmt.Printf("\n🎉 TOUS LES TESTS RETE ONT RÉUSSI !\n")
	} else {
		fmt.Printf("\n⚠️  %d tests ont échoué\n", totalTests-successfulTests)
	}

	fmt.Printf(strings.Repeat("=", 80) + "\n")
}

// GetResults retourne les résultats de tous les tests exécutés
func (tr *TestRunner) GetResults() []validation.RETETestResult {
	return tr.results
}

// SetIncludePatterns définit les patterns de test à inclure
func (tr *TestRunner) SetIncludePatterns(patterns []string) {
	tr.IncludePatterns = patterns
}

// SetExcludePatterns définit les patterns de test à exclure
func (tr *TestRunner) SetExcludePatterns(patterns []string) {
	tr.ExcludePatterns = patterns
}
