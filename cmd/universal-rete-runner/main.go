package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/treivax/tsd/rete"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("🧪 RUNNER UNIVERSEL - TESTS COMPLETS RÉSEAU RETE")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("Pipeline unique avec propagation RETE complète")
	fmt.Printf("Date: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()

	// Trouver tous les fichiers de test
	testDirs := []struct {
		path     string
		category string
	}{
		{"test/coverage/alpha", "alpha"},
		{"beta_coverage_tests", "beta"},
		{"constraint/test/integration", "integration"},
	}

	type TestFile struct {
		name       string
		category   string
		constraint string
		facts      string
	}

	var allTestFiles []TestFile
	for _, dir := range testDirs {
		pattern := filepath.Join(dir.path, "*.constraint")
		matches, _ := filepath.Glob(pattern)

		for _, constraintFile := range matches {
			base := strings.TrimSuffix(constraintFile, ".constraint")
			factsFile := base + ".facts"

			if _, err := os.Stat(factsFile); os.IsNotExist(err) {
				continue
			}

			baseName := filepath.Base(base)
			allTestFiles = append(allTestFiles, TestFile{
				name:       baseName,
				category:   dir.category,
				constraint: constraintFile,
				facts:      factsFile,
			})
		}
	}

	fmt.Printf("�� Trouvé %d tests au total\n\n", len(allTestFiles))

	// Tests qui doivent échouer (tests de détection d'erreurs)
	errorTests := map[string]bool{
		"error_args_test": true,
	}

	// Exécuter tous les tests
	passed := 0
	failed := 0
	for i, testFile := range allTestFiles {
		fmt.Printf("Test %d/%d: %s... ", i+1, len(allTestFiles), testFile.name)

		pipeline := rete.NewConstraintPipeline()
		storage := rete.NewMemoryStorage()

		// Capturer stdout pour détecter les erreurs d'injection
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Canal pour lire la sortie en temps réel
		outputChan := make(chan string)
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outputChan <- buf.String()
		}()

		network, facts, err := pipeline.BuildNetworkFromConstraintFileWithFacts(
			testFile.constraint,
			testFile.facts,
			storage,
		)

		// Restaurer stdout
		w.Close()
		os.Stdout = oldStdout

		// Lire la sortie capturée
		output := <-outputChan

		// Afficher la sortie capturée
		fmt.Print(output)

		// Détecter si des erreurs d'injection ont eu lieu
		hasInjectionErrors := strings.Contains(output, "⚠️ Erreur injection fait")

		// Si c'est un test d'erreur, l'échec est un succès
		isErrorTest := errorTests[testFile.name]

		if err != nil {
			if isErrorTest {
				fmt.Printf("✅ PASSED (error detected as expected)\n")
				passed++
			} else {
				fmt.Printf("❌ FAILED\n")
				failed++
			}
			continue
		}

		// Pour les tests d'erreur, vérifier si des erreurs d'injection ont été détectées
		if isErrorTest {
			if hasInjectionErrors {
				fmt.Printf("✅ PASSED (injection errors detected as expected)\n")
				passed++
			} else {
				fmt.Printf("❌ FAILED (error should have been detected)\n")
				failed++
			}
			continue
		}

		// Compter les activations
		activations := 0
		for _, terminal := range network.TerminalNodes {
			if terminal.Memory != nil && terminal.Memory.Tokens != nil {
				activations += len(terminal.Memory.Tokens)
			}
		}

		fmt.Printf("✅ PASSED - T:%d R:%d F:%d A:%d\n",
			len(network.TypeNodes), len(network.TerminalNodes), len(facts), activations)
		passed++
	}

	fmt.Println()
	fmt.Printf("Résumé: %d tests, %d réussis ✅, %d échoués ❌\n", len(allTestFiles), passed, failed)
	if failed == 0 {
		fmt.Println("🎉 TOUS LES TESTS SONT PASSÉS!")
	}
}
