// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

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
	exitCode := Run(os.Stdout, os.Stderr)
	os.Exit(exitCode)
}

// TestResult holds the result of a single test
type TestResult struct {
	Name          string
	Category      string
	Passed        bool
	TypeNodes     int
	TerminalNodes int
	Facts         int
	Activations   int
	Error         error
	Output        string
}

// RunResult holds the overall test run results
type RunResult struct {
	Total   int
	Passed  int
	Failed  int
	Results []TestResult
}

// Run executes the universal RETE runner and returns an exit code
// This function is testable and doesn't call os.Exit
func Run(stdout, stderr io.Writer) int {
	result := RunTests(stdout)

	// Print summary
	fmt.Fprintln(stdout)
	fmt.Fprintf(stdout, "Résumé: %d tests, %d réussis ✅, %d échoués ❌\n", result.Total, result.Passed, result.Failed)
	if result.Failed == 0 {
		fmt.Fprintln(stdout, "🎉 TOUS LES TESTS SONT PASSÉS!")
		return 0
	}

	return 1
}

// RunTests discovers and executes all tests, returning the results
func RunTests(stdout io.Writer) *RunResult {
	PrintHeader(stdout)

	testFiles := DiscoverTests()
	fmt.Fprintf(stdout, "🔍 Trouvé %d tests au total\n\n", len(testFiles))

	result := &RunResult{
		Total:   len(testFiles),
		Results: make([]TestResult, 0, len(testFiles)),
	}

	errorTests := GetErrorTests()

	for i, testFile := range testFiles {
		fmt.Fprintf(stdout, "Test %d/%d: %s... ", i+1, len(testFiles), testFile.Name)

		testResult := ExecuteTest(testFile, errorTests[testFile.Name])
		result.Results = append(result.Results, testResult)

		// Print output captured during test
		if testResult.Output != "" {
			fmt.Fprint(stdout, testResult.Output)
		}

		if testResult.Passed {
			if testResult.Error != nil {
				fmt.Fprintf(stdout, "✅ PASSED (error detected as expected)\n")
			} else if errorTests[testFile.Name] {
				fmt.Fprintf(stdout, "✅ PASSED (injection errors detected as expected)\n")
			} else {
				fmt.Fprintf(stdout, "✅ PASSED - T:%d R:%d F:%d A:%d\n",
					testResult.TypeNodes, testResult.TerminalNodes, testResult.Facts, testResult.Activations)
			}
			result.Passed++
		} else {
			if errorTests[testFile.Name] {
				fmt.Fprintf(stdout, "❌ FAILED (error should have been detected)\n")
			} else {
				fmt.Fprintf(stdout, "❌ FAILED\n")
			}
			result.Failed++
		}
	}

	return result
}

// TestFile represents a test with constraint and facts files
type TestFile struct {
	Name       string
	Category   string
	Constraint string
	Facts      string
}

// DiscoverTests finds all test files in the test directories
func DiscoverTests() []TestFile {
	testDirs := []struct {
		path     string
		category string
	}{
		{"test/coverage/alpha", "alpha"},
		{"beta_coverage_tests", "beta"},
		{"constraint/test/integration", "integration"},
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
				Name:       baseName,
				Category:   dir.category,
				Constraint: constraintFile,
				Facts:      factsFile,
			})
		}
	}

	return allTestFiles
}

// GetErrorTests returns the set of tests that should produce errors
func GetErrorTests() map[string]bool {
	return map[string]bool{
		"error_args_test": true,
	}
}

// ExecuteTest runs a single test and returns the result
func ExecuteTest(testFile TestFile, expectError bool) TestResult {
	result := TestResult{
		Name:     testFile.Name,
		Category: testFile.Category,
	}

	pipeline := rete.NewConstraintPipeline()
	storage := rete.NewMemoryStorage()

	// Capture stdout to detect injection errors
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Channel to read the output in real-time
	outputChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outputChan <- buf.String()
	}()

	// Ingest constraint file
	network, err := pipeline.IngestFile(testFile.Constraint, nil, storage)
	if err != nil {
		// Restore stdout
		w.Close()
		os.Stdout = oldStdout

		// Read captured output
		output := <-outputChan
		result.Output = output
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to ingest constraint file: %v", err))
		return result
	}

	// Ingest facts file
	network, err = pipeline.IngestFile(testFile.Facts, network, storage)
	if err != nil {
		// Restore stdout
		w.Close()
		os.Stdout = oldStdout

		// Read captured output
		output := <-outputChan
		result.Output = output
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to ingest facts file: %v", err))
		return result
	}

	// Collect facts from storage
	facts := storage.GetAllFacts()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	output := <-outputChan
	result.Output = output

	// Detect injection errors in output
	hasInjectionErrors := strings.Contains(output, "⚠️ Erreur injection fait")

	if err != nil {
		result.Error = err
		result.Passed = expectError
		return result
	}

	// For error tests, check if injection errors were detected
	if expectError {
		result.Passed = hasInjectionErrors
		return result
	}

	// Count activations
	activations := 0
	for _, terminal := range network.TerminalNodes {
		if terminal.Memory != nil && terminal.Memory.Tokens != nil {
			activations += len(terminal.Memory.Tokens)
		}
	}

	result.TypeNodes = len(network.TypeNodes)
	result.TerminalNodes = len(network.TerminalNodes)
	result.Facts = len(facts)
	result.Activations = activations
	result.Passed = true

	return result
}

// PrintHeader prints the test runner header
func PrintHeader(w io.Writer) {
	fmt.Fprintln(w, "═══════════════════════════════════════════════════════════")
	fmt.Fprintln(w, "🧪 RUNNER UNIVERSEL - TESTS COMPLETS RÉSEAU RETE")
	fmt.Fprintln(w, "═══════════════════════════════════════════════════════════")
	fmt.Fprintln(w, "Pipeline unique avec propagation RETE complète")
	fmt.Fprintf(w, "Date: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintln(w)
}
