package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/treivax/tsd/internal/validation"
)

const (
	// Couleurs pour l'affichage
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

// TestSuite représente un ensemble de tests
type TestSuite struct {
	Name        string
	Description string
	Directory   string
	Tests       []TestPair
}

// TestPair représente un couple constraint/facts
type TestPair struct {
	Name           string
	ConstraintFile string
	FactsFile      string
	Category       string
}

// UnifiedTestRunner gère l'exécution de tous les tests RETE
type UnifiedTestRunner struct {
	projectRoot string
	suites      []TestSuite
}

func NewUnifiedTestRunner(projectRoot string) *UnifiedTestRunner {
	return &UnifiedTestRunner{
		projectRoot: projectRoot,
		suites:      make([]TestSuite, 0),
	}
}

// executeTestWithCompleteRete exécute un test avec la nouvelle implémentation RETE complète
func (u *UnifiedTestRunner) executeTestWithCompleteRete(constraintFile, factsFile string) validation.RETETestResult {
	testPath := strings.TrimSuffix(constraintFile, ".constraint")

	// Utiliser la fonction de validation RETE complète
	result, err := validation.ValidateRETEWithFile(testPath, 5*time.Minute)
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

	return *result
}

func (u *UnifiedTestRunner) DiscoverTests() error {
	fmt.Printf("%s🔍 DÉCOUVERTE AUTOMATIQUE DES TESTS RETE%s\n", ColorBlue, ColorReset)
	fmt.Printf("================================================\n\n")

	// Répertoires à scanner
	testDirs := []struct {
		path        string
		name        string
		description string
	}{
		{
			path:        filepath.Join(u.projectRoot, "beta_coverage_tests"),
			name:        "Beta Tests",
			description: "Tests de couverture beta - Validation complète RETE",
		},
		{
			path:        filepath.Join(u.projectRoot, "test", "coverage", "alpha"),
			name:        "Alpha Tests",
			description: "Tests de couverture alpha - Nœuds alpha du réseau RETE",
		},
		{
			path:        filepath.Join(u.projectRoot, "constraint", "test", "integration"),
			name:        "Integration Tests",
			description: "Tests d'intégration - Validation fonctionnelle complète",
		},
	}

	totalTests := 0
	for _, dir := range testDirs {
		if _, err := os.Stat(dir.path); os.IsNotExist(err) {
			fmt.Printf("%s⚠️  Répertoire %s non trouvé : %s%s\n", ColorYellow, dir.name, dir.path, ColorReset)
			continue
		}

		suite := TestSuite{
			Name:        dir.name,
			Description: dir.description,
			Directory:   dir.path,
			Tests:       make([]TestPair, 0),
		}

		// Scanner les fichiers .constraint
		constraintFiles, err := filepath.Glob(filepath.Join(dir.path, "*.constraint"))
		if err != nil {
			return fmt.Errorf("erreur scan %s: %v", dir.path, err)
		}

		for _, constraintFile := range constraintFiles {
			baseName := strings.TrimSuffix(filepath.Base(constraintFile), ".constraint")
			factsFile := filepath.Join(dir.path, baseName+".facts")

			// Vérifier si le fichier .facts existe
			if _, err := os.Stat(factsFile); err == nil {
				category := u.categorizeTest(baseName)
				suite.Tests = append(suite.Tests, TestPair{
					Name:           baseName,
					ConstraintFile: constraintFile,
					FactsFile:      factsFile,
					Category:       category,
				})
			} else {
				fmt.Printf("%s⚠️  Test %s - fichier .facts manquant%s\n", ColorYellow, baseName, ColorReset)
			}
		}

		// Trier les tests par nom
		sort.Slice(suite.Tests, func(i, j int) bool {
			return suite.Tests[i].Name < suite.Tests[j].Name
		})

		if len(suite.Tests) > 0 {
			u.suites = append(u.suites, suite)
			totalTests += len(suite.Tests)
			fmt.Printf("%s✅ Suite %s%s : %d tests découverts\n", ColorGreen, dir.name, ColorReset, len(suite.Tests))
		}
	}

	fmt.Printf("\n%s📊 RÉSUMÉ DÉCOUVERTE%s\n", ColorCyan, ColorReset)
	fmt.Printf("==================\n")
	fmt.Printf("• Suites de tests : %d\n", len(u.suites))
	fmt.Printf("• Total tests : %d\n", totalTests)
	fmt.Printf("\n")

	return nil
}

func (u *UnifiedTestRunner) categorizeTest(testName string) string {
	name := strings.ToLower(testName)

	switch {
	case strings.Contains(name, "join"):
		return "Join Operations"
	case strings.Contains(name, "exists"):
		return "Exists Conditions"
	case strings.Contains(name, "not"):
		return "Negation Logic"
	case strings.Contains(name, "alpha"):
		return "Alpha Nodes"
	case strings.Contains(name, "beta"):
		return "Beta Nodes"
	case strings.Contains(name, "comparison"):
		return "Comparison Ops"
	case strings.Contains(name, "arithmetic"):
		return "Arithmetic Ops"
	case strings.Contains(name, "contains") || strings.Contains(name, "like") || strings.Contains(name, "matches"):
		return "String Ops"
	case strings.Contains(name, "integration") || strings.Contains(name, "comprehensive"):
		return "Integration"
	default:
		return "General"
	}
}

func (u *UnifiedTestRunner) RunAllTests() ([]validation.RETETestResult, error) {
	fmt.Printf("%s🚀 EXÉCUTION COMPLÈTE DE TOUS LES TESTS RETE%s\n", ColorBlue+ColorBold, ColorReset)
	fmt.Printf("==============================================\n\n")

	var allResults []validation.RETETestResult
	totalTests := 0
	for _, suite := range u.suites {
		totalTests += len(suite.Tests)
	}

	currentTest := 0
	startTime := time.Now()

	for _, suite := range u.suites {
		fmt.Printf("%s📁 SUITE: %s%s\n", ColorCyan+ColorBold, suite.Name, ColorReset)
		fmt.Printf("   %s\n", suite.Description)
		fmt.Printf("   Tests: %d\n\n", len(suite.Tests))

		for _, test := range suite.Tests {
			currentTest++
			fmt.Printf("%s[%d/%d]%s %s🎯 %s%s (%s%s%s) ... ",
				ColorWhite, currentTest, totalTests, ColorReset,
				ColorCyan, test.Name, ColorReset,
				ColorYellow, test.Category, ColorReset)

			// Utiliser la nouvelle implémentation RETE complète
			result := u.executeTestWithCompleteRete(test.ConstraintFile, test.FactsFile)

			// Ajouter les métadonnées de catégorie
			result.TestName = test.Name
			allResults = append(allResults, result)

			if result.Success {
				fmt.Printf("%s✅ RÉUSSI%s (%s%d tokens%s)\n",
					ColorGreen, ColorReset,
					ColorWhite, len(result.ObservedTokens), ColorReset)
			} else {
				fmt.Printf("%s❌ ÉCHEC%s\n", ColorRed, ColorReset)
				if result.ValidationNote != "" {
					fmt.Printf("      %sErreur: %s%s\n", ColorRed, result.ValidationNote, ColorReset)
				}
			}
		}
		fmt.Println()
	}

	totalDuration := time.Since(startTime)

	fmt.Printf("%s📊 EXÉCUTION TERMINÉE%s\n", ColorBlue+ColorBold, ColorReset)
	fmt.Printf("====================\n")
	fmt.Printf("• Total tests: %d\n", totalTests)
	fmt.Printf("• Durée totale: %v\n", totalDuration.Round(time.Millisecond))
	fmt.Printf("\n")

	return allResults, nil
}

func (u *UnifiedTestRunner) GenerateDetailedReport(results []validation.RETETestResult) string {
	report := strings.Builder{}
	now := time.Now()

	// En-tête du rapport
	report.WriteString(fmt.Sprintf("# RAPPORT DÉTAILLÉ VALIDATION RETE UNIFIÉE\n\n"))
	report.WriteString(fmt.Sprintf("**Date:** %s  \n", now.Format("2006-01-02 15:04:05")))
	report.WriteString(fmt.Sprintf("**Système:** Runner unifié tous tests (Alpha + Beta + Intégration)  \n"))
	report.WriteString(fmt.Sprintf("**Méthode:** Tokens RÉELLEMENT extraits du réseau RETE  \n\n"))

	// Statistiques globales
	passed := 0
	failed := 0
	totalDuration := time.Duration(0)
	categoriesStats := make(map[string]struct{ passed, total int })

	for _, result := range results {
		if result.Success {
			passed++
		} else {
			failed++
		}
		// totalDuration += result.ExecutionTime // Pas de temps d'exécution individuel avec RETE complet

		// Catégoriser le test
		category := u.categorizeTest(result.TestName)
		stats := categoriesStats[category]
		stats.total++
		if result.Success {
			stats.passed++
		}
		categoriesStats[category] = stats
	}

	total := len(results)
	successRate := float64(passed) / float64(total) * 100

	report.WriteString("## 📊 RÉSUMÉ EXÉCUTIF\n\n")
	report.WriteString(fmt.Sprintf("- **Tests totaux:** %d\n", total))
	report.WriteString(fmt.Sprintf("- **Tests réussis:** %d\n", passed))
	report.WriteString(fmt.Sprintf("- **Tests échoués:** %d\n", failed))
	report.WriteString(fmt.Sprintf("- **Taux de succès:** %.1f%%\n", successRate))
	report.WriteString(fmt.Sprintf("- **Durée totale:** %v\n", totalDuration.Round(time.Millisecond)))
	report.WriteString("\n")

	// Statistiques par catégorie
	report.WriteString("## 📈 STATISTIQUES PAR CATÉGORIE\n\n")
	for category, stats := range categoriesStats {
		rate := float64(stats.passed) / float64(stats.total) * 100
		status := "✅"
		if rate < 100 {
			status = "⚠️"
		}
		report.WriteString(fmt.Sprintf("- **%s** %s: %d/%d (%.1f%%)\n", category, status, stats.passed, stats.total, rate))
	}
	report.WriteString("\n")

	// Détails par test
	currentSuite := ""
	for _, result := range results {
		// Déterminer la suite basée sur le nom du test
		suite := u.determineSuiteFromTestName(result.TestName)
		if suite != currentSuite {
			report.WriteString(fmt.Sprintf("## 🔥 SUITE: %s\n\n", suite))
			currentSuite = suite
		}

		// En-tête du test
		status := "✅ RÉUSSI"
		statusEmoji := "✅"
		if !result.Success {
			status = "❌ ÉCHEC"
			statusEmoji = "❌"
		}

		report.WriteString(fmt.Sprintf("### %s Test: %s\n\n", statusEmoji, result.TestName))
		report.WriteString(fmt.Sprintf("**Catégorie:** %s  \n", u.categorizeTest(result.TestName)))
		report.WriteString(fmt.Sprintf("**Statut:** %s  \n", status))
		report.WriteString(fmt.Sprintf("**Tokens:** %d observés  \n", len(result.ObservedTokens)))
		report.WriteString(fmt.Sprintf("**Note:** %s  \n\n", result.ValidationNote))

		if result.ValidationNote != "" && !result.Success {
			report.WriteString(fmt.Sprintf("**❌ Erreur:** %s\n\n", result.ValidationNote))
		}

		// Règles du test
		report.WriteString("#### 📋 Règles du Test\n\n")
		for i, rule := range result.Rules {
			report.WriteString(fmt.Sprintf("%d. `%s`\n", i+1, rule))
		}
		report.WriteString("\n")

		// Faits soumis
		report.WriteString("#### 📊 Faits Soumis au Réseau RETE\n\n")
		for i, fact := range result.Facts {
			report.WriteString(fmt.Sprintf("%d. `%s`\n", i+1, fact))
		}
		report.WriteString("\n")

		// Tokens observés (RETE réel)
		report.WriteString("#### 🔥 Tokens Observés (Extraits du Réseau RETE)\n\n")
		if len(result.ObservedTokens) > 0 {
			for i, token := range result.ObservedTokens {
				report.WriteString(fmt.Sprintf("**Token %d:**\n", i+1))
				report.WriteString(fmt.Sprintf("- **Règle:** %s\n", token.RuleName))
				report.WriteString(fmt.Sprintf("- **Clé:** `%s`\n", token.Key))
				report.WriteString("- **Faits composant le token:**\n")
				j := 1
				for _, fact := range token.Facts {
					report.WriteString(fmt.Sprintf("  %d. %s: %s (ID: %s)\n", j, fact.Type, fact.Values, fact.ID))
					j++
				}
				report.WriteString("\n")
			}
		} else {
			report.WriteString("*Aucun token observé*\n\n")
		}

		report.WriteString("---\n\n")
	}

	// Conclusion
	report.WriteString("## 🎯 CONCLUSION\n\n")
	if failed == 0 {
		report.WriteString("✅ **VALIDATION COMPLÈTE RÉUSSIE**\n\n")
		report.WriteString("Tous les tests du réseau RETE (Alpha, Beta, et Intégration) ont été validés avec succès. ")
		report.WriteString("Le système d'extraction de tokens depuis le réseau RETE fonctionne parfaitement.\n\n")
	} else {
		report.WriteString("⚠️ **VALIDATION PARTIELLE**\n\n")
		report.WriteString(fmt.Sprintf("%d test(s) nécessitent une attention particulière. ", failed))
		report.WriteString("Consultez les détails ci-dessus pour les erreurs spécifiques.\n\n")
	}

	report.WriteString("**Méthode de validation:** Tokens réellement extraits du réseau RETE vs simulation\n")
	report.WriteString("**Architecture:** Réseau RETE authentique avec nœuds alpha et beta\n")

	return report.String()
}

func (u *UnifiedTestRunner) determineSuiteFromTestName(testName string) string {
	name := strings.ToLower(testName)
	switch {
	case strings.Contains(name, "alpha_"):
		return "Alpha Tests"
	case strings.HasPrefix(name, "join_") || strings.HasPrefix(name, "exists_") || strings.HasPrefix(name, "not_") || strings.Contains(name, "complex"):
		return "Beta Tests"
	default:
		return "Integration Tests"
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%s=== RUNNER RETE UNIFIÉ ====%s\n", ColorBlue+ColorBold, ColorReset)
		fmt.Printf("Exécution de tous les tests RETE (Alpha + Beta + Intégration)\n\n")
		fmt.Printf("Usage:\n")
		fmt.Printf("  %s <project_root>          - Exécuter tous les tests\n", os.Args[0])
		fmt.Printf("  %s <project_root> report   - Générer seulement le rapport\n", os.Args[0])
		fmt.Printf("\nExemple:\n")
		fmt.Printf("  %s /home/resinsec/dev/tsd\n", os.Args[0])
		os.Exit(1)
	}

	projectRoot := os.Args[1]
	reportOnly := len(os.Args) > 2 && os.Args[2] == "report"

	runner := NewUnifiedTestRunner(projectRoot)

	// Découverte des tests
	if err := runner.DiscoverTests(); err != nil {
		fmt.Printf("%sErreur découverte tests: %v%s\n", ColorRed, err, ColorReset)
		os.Exit(1)
	}

	if len(runner.suites) == 0 {
		fmt.Printf("%sAucune suite de tests trouvée%s\n", ColorRed, ColorReset)
		os.Exit(1)
	}

	var results []validation.RETETestResult
	var err error

	if !reportOnly {
		// Exécution des tests
		results, err = runner.RunAllTests()
		if err != nil {
			fmt.Printf("%sErreur exécution tests: %v%s\n", ColorRed, err, ColorReset)
			os.Exit(1)
		}
	}

	// Génération du rapport
	if len(results) > 0 {
		report := runner.GenerateDetailedReport(results)

		// Sauvegarde du rapport (nom fixe pour éviter la prolifération)
		reportFile := filepath.Join(projectRoot, "RAPPORT_RETE_UNIFIE.md")
		if err := os.WriteFile(reportFile, []byte(report), 0644); err != nil {
			fmt.Printf("%sErreur sauvegarde rapport: %v%s\n", ColorRed, err, ColorReset)
		} else {
			fmt.Printf("%s📄 Rapport détaillé généré: %s%s\n", ColorCyan, reportFile, ColorReset)
		}

		// Statistiques finales
		passed := 0
		for _, result := range results {
			if result.Success {
				passed++
			}
		}

		fmt.Printf("\n%s🎉 VALIDATION TERMINÉE%s\n", ColorGreen+ColorBold, ColorReset)
		fmt.Printf("===================\n")
		fmt.Printf("• Tests réussis: %d/%d\n", passed, len(results))
		fmt.Printf("• Taux de succès: %.1f%%\n", float64(passed)/float64(len(results))*100)

		if passed == len(results) {
			fmt.Printf("%s✅ TOUS LES TESTS RETE ONT RÉUSSI !%s\n", ColorGreen, ColorReset)
			os.Exit(0)
		} else {
			fmt.Printf("%s⚠️  %d test(s) nécessitent une attention%s\n", ColorYellow, len(results)-passed, ColorReset)
			os.Exit(1)
		}
	}
}
