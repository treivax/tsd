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
	"regexp"
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
		pattern := filepath.Join(dir.path, "*.tsd")
		matches, _ := filepath.Glob(pattern)

		for _, tsdFile := range matches {
			baseName := filepath.Base(strings.TrimSuffix(tsdFile, ".tsd"))
			allTestFiles = append(allTestFiles, TestFile{
				Name:       baseName,
				Category:   dir.category,
				Constraint: tsdFile,
				Facts:      tsdFile, // .tsd files contain both constraints and facts
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

	// For alpha/beta coverage tests, inject missing action definitions
	var modifiedContent string
	var useModified bool
	if testFile.Category == "alpha" || testFile.Category == "beta" {
		content, err := os.ReadFile(testFile.Constraint)
		if err == nil {
			modifiedContent = InjectMissingActions(string(content))
			if modifiedContent != string(content) {
				useModified = true
			}
		}
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

	var network *rete.ReteNetwork
	var err error

	// Ingest constraint file (with modified content if needed)
	if useModified {
		// Create temporary file with modified content
		tmpFile, tmpErr := os.CreateTemp("", "test-*.tsd")
		if tmpErr == nil {
			tmpFile.WriteString(modifiedContent)
			tmpFile.Close()
			network, err = pipeline.IngestFile(tmpFile.Name(), nil, storage)
			os.Remove(tmpFile.Name())
		} else {
			err = tmpErr
		}
	} else {
		network, err = pipeline.IngestFile(testFile.Constraint, nil, storage)
	}

	if err != nil {
		// Restore stdout
		w.Close()
		os.Stdout = oldStdout

		// Read captured output
		output := <-outputChan
		result.Output = output
		result.Error = fmt.Errorf("failed to ingest constraint file: %w", err)
		result.Passed = expectError
		return result
	}

	// Ingest facts file only if different from constraint file
	// (TSD files contain both constraints and facts, avoid duplicate ingestion)
	if testFile.Facts != testFile.Constraint && !useModified {
		network, err = pipeline.IngestFile(testFile.Facts, network, storage)
		if err != nil {
			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			output := <-outputChan
			result.Output = output
			result.Error = fmt.Errorf("failed to ingest facts file: %w", err)
			result.Passed = expectError
			return result
		}
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

	// Test passes if no error occurred and network was built successfully
	// Activations count doesn't matter - some tests (especially with NOT) may have 0 activations
	result.Passed = err == nil && network != nil && !hasInjectionErrors

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

// InjectMissingActions analyzes TSD content and injects action definitions for any actions
// used in rules that are not already defined. This allows coverage tests to run without
// manually defining every action.
func InjectMissingActions(content string) string {
	// Find all action calls in rules with their full argument lists
	actionCallRegex := regexp.MustCompile(`==>\s*([a-zA-Z_][a-zA-Z0-9_]*)\(([^)]*)\)`)
	matches := actionCallRegex.FindAllStringSubmatch(content, -1)

	if len(matches) == 0 {
		return content
	}

	// Find already defined actions
	actionDefRegex := regexp.MustCompile(`(?m)^action\s+([a-zA-Z_][a-zA-Z0-9_]*)\(`)
	definedActions := make(map[string]bool)
	definedMatches := actionDefRegex.FindAllStringSubmatch(content, -1)
	for _, match := range definedMatches {
		if len(match) > 1 {
			definedActions[match[1]] = true
		}
	}

	// Collect missing actions with their argument counts
	missingActions := make(map[string]int)
	for _, match := range matches {
		if len(match) > 2 {
			actionName := match[1]
			argsStr := strings.TrimSpace(match[2])

			if !definedActions[actionName] {
				// Count arguments (comma-separated)
				argCount := 0
				if argsStr != "" {
					argCount = strings.Count(argsStr, ",") + 1
				}

				// Store the maximum argument count for this action
				if existingCount, exists := missingActions[actionName]; !exists || argCount > existingCount {
					missingActions[actionName] = argCount
				}
			}
		}
	}

	if len(missingActions) == 0 {
		return content
	}

	// Generate action definitions - insert after type definitions, before rules
	var injectedActions strings.Builder
	injectedActions.WriteString("\n// Auto-generated action definitions for testing\n")
	for actionName, argCount := range missingActions {
		// Generate action with correct number of string arguments
		var args []string
		for i := 0; i < argCount; i++ {
			args = append(args, fmt.Sprintf("arg%d: string", i+1))
		}
		argsList := strings.Join(args, ", ")
		injectedActions.WriteString(fmt.Sprintf("action %s(%s)\n", actionName, argsList))
	}
	injectedActions.WriteString("\n")

	// Find position to inject (after last type definition, before first rule)
	ruleRegex := regexp.MustCompile(`(?m)^rule\s+`)
	ruleLoc := ruleRegex.FindStringIndex(content)

	if ruleLoc != nil {
		// Insert before first rule
		insertPos := ruleLoc[0]
		return content[:insertPos] + injectedActions.String() + content[insertPos:]
	}

	// If no rule found, append at end (shouldn't happen in valid tests)
	return content + injectedActions.String()
}
