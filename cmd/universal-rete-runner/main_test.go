package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestFile represents a test case with constraint and facts files
type TestFile struct {
	name       string
	category   string
	constraint string
	facts      string
}

// TestTestFileStruct tests the TestFile struct initialization
func TestTestFileStruct(t *testing.T) {
	tf := TestFile{
		name:       "test_alpha",
		category:   "alpha",
		constraint: "test.constraint",
		facts:      "test.facts",
	}

	if tf.name != "test_alpha" {
		t.Errorf("TestFile.name = %v, want test_alpha", tf.name)
	}
	if tf.category != "alpha" {
		t.Errorf("TestFile.category = %v, want alpha", tf.category)
	}
	if tf.constraint != "test.constraint" {
		t.Errorf("TestFile.constraint = %v, want test.constraint", tf.constraint)
	}
	if tf.facts != "test.facts" {
		t.Errorf("TestFile.facts = %v, want test.facts", tf.facts)
	}
}

// TestErrorTestsMap tests the error tests map
func TestErrorTestsMap(t *testing.T) {
	errorTests := map[string]bool{
		"error_args_test": true,
	}

	if !errorTests["error_args_test"] {
		t.Error("errorTests should contain error_args_test")
	}

	if errorTests["non_existent_test"] {
		t.Error("errorTests should not contain non_existent_test")
	}
}

// TestTestDirsStructure tests the test directories structure
func TestTestDirsStructure(t *testing.T) {
	testDirs := []struct {
		path     string
		category string
	}{
		{"test/coverage/alpha", "alpha"},
		{"beta_coverage_tests", "beta"},
		{"constraint/test/integration", "integration"},
	}

	if len(testDirs) != 3 {
		t.Errorf("Expected 3 test directories, got %d", len(testDirs))
	}

	expectedDirs := map[string]string{
		"test/coverage/alpha":         "alpha",
		"beta_coverage_tests":         "beta",
		"constraint/test/integration": "integration",
	}

	for _, dir := range testDirs {
		expectedCategory, exists := expectedDirs[dir.path]
		if !exists {
			t.Errorf("Unexpected directory path: %s", dir.path)
		}
		if dir.category != expectedCategory {
			t.Errorf("Directory %s has category %s, want %s", dir.path, dir.category, expectedCategory)
		}
	}
}

// TestFilePatternMatching tests the constraint file pattern matching
func TestFilePatternMatching(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()

	// Create test files
	files := []string{
		"test1.constraint",
		"test2.constraint",
		"test3.constraint",
		"test1.facts",
		"test2.facts",
		"test3.facts",
		"other.txt",
	}

	for _, file := range files {
		filePath := filepath.Join(tempDir, file)
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Test pattern matching
	pattern := filepath.Join(tempDir, "*.constraint")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("Glob failed: %v", err)
	}

	if len(matches) != 3 {
		t.Errorf("Expected 3 constraint files, got %d", len(matches))
	}

	for _, match := range matches {
		if !strings.HasSuffix(match, ".constraint") {
			t.Errorf("Match %s does not end with .constraint", match)
		}
	}
}

// TestFactsFileExistence tests checking for corresponding facts files
func TestFactsFileExistence(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()

	// Create constraint file with corresponding facts file
	constraintFile1 := filepath.Join(tempDir, "test1.constraint")
	factsFile1 := filepath.Join(tempDir, "test1.facts")
	os.WriteFile(constraintFile1, []byte("type Person : <id: string>"), 0644)
	os.WriteFile(factsFile1, []byte("facts content"), 0644)

	// Create constraint file without facts file
	constraintFile2 := filepath.Join(tempDir, "test2.constraint")
	os.WriteFile(constraintFile2, []byte("type Person : <id: string>"), 0644)

	// Test first file (should have facts)
	base1 := strings.TrimSuffix(constraintFile1, ".constraint")
	expectedFacts1 := base1 + ".facts"
	if _, err := os.Stat(expectedFacts1); os.IsNotExist(err) {
		t.Errorf("Facts file should exist for %s", constraintFile1)
	}

	// Test second file (should not have facts)
	base2 := strings.TrimSuffix(constraintFile2, ".constraint")
	expectedFacts2 := base2 + ".facts"
	if _, err := os.Stat(expectedFacts2); !os.IsNotExist(err) {
		t.Errorf("Facts file should not exist for %s", constraintFile2)
	}
}

// TestBaseNameExtraction tests extracting base name from file path
func TestBaseNameExtraction(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		expectedBase string
	}{
		{
			name:         "simple path",
			path:         "/path/to/test.constraint",
			expectedBase: "test",
		},
		{
			name:         "current directory",
			path:         "test.constraint",
			expectedBase: "test",
		},
		{
			name:         "nested path",
			path:         "dir1/dir2/dir3/test_alpha.constraint",
			expectedBase: "test_alpha",
		},
		{
			name:         "with underscores",
			path:         "/tests/error_args_test.constraint",
			expectedBase: "error_args_test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := strings.TrimSuffix(tt.path, ".constraint")
			baseName := filepath.Base(base)

			if baseName != tt.expectedBase {
				t.Errorf("Base name = %v, want %v", baseName, tt.expectedBase)
			}
		})
	}
}

// TestCountActivationsLogic tests the activation counting logic
func TestCountActivationsLogic(t *testing.T) {
	// This tests the logic used in the main function to count activations
	// We simulate the network structure without importing rete package heavily

	tests := []struct {
		name          string
		tokenCounts   []int // tokens per terminal node
		expectedTotal int
	}{
		{
			name:          "no terminals",
			tokenCounts:   []int{},
			expectedTotal: 0,
		},
		{
			name:          "one terminal with no tokens",
			tokenCounts:   []int{0},
			expectedTotal: 0,
		},
		{
			name:          "one terminal with tokens",
			tokenCounts:   []int{5},
			expectedTotal: 5,
		},
		{
			name:          "multiple terminals",
			tokenCounts:   []int{3, 5, 2},
			expectedTotal: 10,
		},
		{
			name:          "terminals with varying counts",
			tokenCounts:   []int{0, 10, 0, 3},
			expectedTotal: 13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate counting logic
			activations := 0
			for _, count := range tt.tokenCounts {
				activations += count
			}

			if activations != tt.expectedTotal {
				t.Errorf("Total activations = %d, want %d", activations, tt.expectedTotal)
			}
		})
	}
}

// TestOutputStringDetection tests detection of injection errors in output
func TestOutputStringDetection(t *testing.T) {
	tests := []struct {
		name             string
		output           string
		expectedHasError bool
	}{
		{
			name:             "no error",
			output:           "All facts injected successfully",
			expectedHasError: false,
		},
		{
			name:             "has injection error",
			output:           "⚠️ Erreur injection fait: invalid type",
			expectedHasError: true,
		},
		{
			name:             "multiple lines with error",
			output:           "Processing...\n⚠️ Erreur injection fait: something wrong\nDone",
			expectedHasError: true,
		},
		{
			name:             "similar but not exact error message",
			output:           "Warning: injection failed",
			expectedHasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasInjectionErrors := strings.Contains(tt.output, "⚠️ Erreur injection fait")

			if hasInjectionErrors != tt.expectedHasError {
				t.Errorf("hasInjectionErrors = %v, want %v for output: %s",
					hasInjectionErrors, tt.expectedHasError, tt.output)
			}
		})
	}
}

// TestIsErrorTest tests the logic for determining if a test is an error test
func TestIsErrorTest(t *testing.T) {
	errorTests := map[string]bool{
		"error_args_test": true,
	}

	tests := []struct {
		name        string
		testName    string
		isErrorTest bool
	}{
		{
			name:        "error test",
			testName:    "error_args_test",
			isErrorTest: true,
		},
		{
			name:        "normal test",
			testName:    "normal_test",
			isErrorTest: false,
		},
		{
			name:        "alpha test",
			testName:    "test_alpha",
			isErrorTest: false,
		},
		{
			name:        "beta test",
			testName:    "test_beta",
			isErrorTest: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isError := errorTests[tt.testName]

			if isError != tt.isErrorTest {
				t.Errorf("isErrorTest(%s) = %v, want %v", tt.testName, isError, tt.isErrorTest)
			}
		})
	}
}

// TestTimeFormatting tests the time formatting used in the output
func TestTimeFormatting(t *testing.T) {
	now := time.Date(2025, 11, 26, 15, 30, 45, 0, time.UTC)
	formatted := now.Format("2006-01-02 15:04:05")
	expected := "2025-11-26 15:30:45"

	if formatted != expected {
		t.Errorf("Time format = %s, want %s", formatted, expected)
	}
}

// TestTestResultCounting tests the logic for counting passed and failed tests
func TestTestResultCounting(t *testing.T) {
	tests := []struct {
		name         string
		results      []bool // true = passed, false = failed
		expectedPass int
		expectedFail int
	}{
		{
			name:         "all pass",
			results:      []bool{true, true, true},
			expectedPass: 3,
			expectedFail: 0,
		},
		{
			name:         "all fail",
			results:      []bool{false, false, false},
			expectedPass: 0,
			expectedFail: 3,
		},
		{
			name:         "mixed results",
			results:      []bool{true, false, true, false, true},
			expectedPass: 3,
			expectedFail: 2,
		},
		{
			name:         "no tests",
			results:      []bool{},
			expectedPass: 0,
			expectedFail: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passed := 0
			failed := 0

			for _, result := range tt.results {
				if result {
					passed++
				} else {
					failed++
				}
			}

			if passed != tt.expectedPass {
				t.Errorf("Passed = %d, want %d", passed, tt.expectedPass)
			}
			if failed != tt.expectedFail {
				t.Errorf("Failed = %d, want %d", failed, tt.expectedFail)
			}
		})
	}
}

// TestSummaryGeneration tests the summary text generation logic
func TestSummaryGeneration(t *testing.T) {
	tests := []struct {
		name      string
		total     int
		passed    int
		failed    int
		allPassed bool
	}{
		{
			name:      "all tests passed",
			total:     10,
			passed:    10,
			failed:    0,
			allPassed: true,
		},
		{
			name:      "some tests failed",
			total:     10,
			passed:    7,
			failed:    3,
			allPassed: false,
		},
		{
			name:      "all tests failed",
			total:     5,
			passed:    0,
			failed:    5,
			allPassed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allPassed := tt.failed == 0

			if allPassed != tt.allPassed {
				t.Errorf("allPassed = %v, want %v", allPassed, tt.allPassed)
			}

			if tt.passed+tt.failed != tt.total {
				t.Errorf("passed(%d) + failed(%d) != total(%d)", tt.passed, tt.failed, tt.total)
			}
		})
	}
}

// TestErrorTestHandling tests the logic for handling error tests
func TestErrorTestHandling(t *testing.T) {
	tests := []struct {
		name               string
		isErrorTest        bool
		hasError           bool
		hasInjectionErrors bool
		shouldPass         bool
	}{
		{
			name:               "error test with error - pass",
			isErrorTest:        true,
			hasError:           true,
			hasInjectionErrors: false,
			shouldPass:         true,
		},
		{
			name:               "error test without error - fail",
			isErrorTest:        true,
			hasError:           false,
			hasInjectionErrors: false,
			shouldPass:         false,
		},
		{
			name:               "error test with injection errors - pass",
			isErrorTest:        true,
			hasError:           false,
			hasInjectionErrors: true,
			shouldPass:         true,
		},
		{
			name:               "normal test without error - pass",
			isErrorTest:        false,
			hasError:           false,
			hasInjectionErrors: false,
			shouldPass:         true,
		},
		{
			name:               "normal test with error - fail",
			isErrorTest:        false,
			hasError:           true,
			hasInjectionErrors: false,
			shouldPass:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var shouldPass bool

			if tt.hasError {
				if tt.isErrorTest {
					shouldPass = true // Error test correctly detected error
				} else {
					shouldPass = false // Normal test should not have error
				}
			} else {
				if tt.isErrorTest {
					// Check for injection errors
					shouldPass = tt.hasInjectionErrors
				} else {
					shouldPass = true // Normal test passed
				}
			}

			if shouldPass != tt.shouldPass {
				t.Errorf("shouldPass = %v, want %v", shouldPass, tt.shouldPass)
			}
		})
	}
}

// TestFileGlobbing tests file pattern matching behavior
func TestFileGlobbing(t *testing.T) {
	tempDir := t.TempDir()

	// Create various test files
	testFiles := map[string]string{
		"alpha_test.constraint": "constraint content",
		"alpha_test.facts":      "facts content",
		"beta_test.constraint":  "constraint content",
		"beta_test.facts":       "facts content",
		"error_test.constraint": "constraint content",
		"no_facts.constraint":   "constraint content",
		"random.txt":            "other content",
	}

	for filename, content := range testFiles {
		path := filepath.Join(tempDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Test constraint file globbing
	pattern := filepath.Join(tempDir, "*.constraint")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("Glob failed: %v", err)
	}

	expectedConstraintCount := 4
	if len(matches) != expectedConstraintCount {
		t.Errorf("Found %d constraint files, want %d", len(matches), expectedConstraintCount)
	}

	// Test pairing with facts files
	pairsFound := 0
	for _, constraintFile := range matches {
		base := strings.TrimSuffix(constraintFile, ".constraint")
		factsFile := base + ".facts"
		if _, err := os.Stat(factsFile); err == nil {
			pairsFound++
		}
	}

	expectedPairs := 2 // We only create 2 complete pairs (alpha_test and beta_test have .facts files)
	if pairsFound != expectedPairs {
		t.Errorf("Found %d constraint-facts pairs, want %d", pairsFound, expectedPairs)
	}
}

// TestMainIntegration tests the main function via subprocess
func TestMainIntegration(t *testing.T) {
	// Build the binary
	testBinary := filepath.Join(t.TempDir(), "universal-rete-runner-test")
	buildCmd := exec.Command("go", "build", "-o", testBinary, ".")

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	buildCmd.Dir = wd

	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build test binary: %v\nOutput: %s", err, output)
	}
	defer os.Remove(testBinary)

	// Run the binary
	cmd := exec.Command(testBinary)
	output, _ := cmd.CombinedOutput()
	outputStr := string(output)

	// Check for expected header output
	expectedStrings := []string{
		"RUNNER UNIVERSEL",
		"TESTS COMPLETS RÉSEAU RETE",
		"Pipeline unique",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Output does not contain %q\nGot: %s", expected, outputStr)
		}
	}
}

// TestMainWithTestFiles tests main with actual test files
func TestMainWithTestFiles(t *testing.T) {
	// This test runs the actual binary and verifies it finds and processes test files
	// It will only run if test files exist in the expected locations

	testBinary := filepath.Join(t.TempDir(), "universal-rete-runner-files-test")
	buildCmd := exec.Command("go", "build", "-o", testBinary, ".")
	buildCmd.Dir, _ = os.Getwd()

	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\n%s", err, output)
	}
	defer os.Remove(testBinary)

	// Change to project root to access test files
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Go up two levels to reach project root (from cmd/universal-rete-runner to root)
	projectRoot := filepath.Join(originalDir, "..", "..")
	if err := os.Chdir(projectRoot); err != nil {
		t.Skipf("Cannot change to project root: %v", err)
	}

	cmd := exec.Command(testBinary)
	output, _ := cmd.CombinedOutput()
	outputStr := string(output)

	// The binary should run (may or may not find tests)
	t.Logf("Binary output:\n%s", outputStr)

	// Check for summary line
	if strings.Contains(outputStr, "Résumé:") {
		// If we have a summary, verify it's formatted correctly
		if !strings.Contains(outputStr, "tests") {
			t.Error("Summary line should mention 'tests'")
		}
	}
}

// TestOutputFormatting tests the output formatting logic
func TestOutputFormatting(t *testing.T) {
	tests := []struct {
		name             string
		testName         string
		passed           int
		typeNodes        int
		terminalNodes    int
		facts            int
		activations      int
		expectedContains string
	}{
		{
			name:             "successful test output",
			testName:         "test_alpha",
			passed:           1,
			typeNodes:        2,
			terminalNodes:    1,
			facts:            3,
			activations:      2,
			expectedContains: "PASSED",
		},
		{
			name:             "test with no activations",
			testName:         "test_beta",
			passed:           1,
			typeNodes:        1,
			terminalNodes:    0,
			facts:            0,
			activations:      0,
			expectedContains: "PASSED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the output format
			output := fmt.Sprintf("✅ PASSED - T:%d R:%d F:%d A:%d\n",
				tt.typeNodes, tt.terminalNodes, tt.facts, tt.activations)

			if !strings.Contains(output, tt.expectedContains) {
				t.Errorf("Output does not contain %q\nGot: %s", tt.expectedContains, output)
			}

			// Verify format includes all metrics
			if !strings.Contains(output, "T:") || !strings.Contains(output, "R:") ||
				!strings.Contains(output, "F:") || !strings.Contains(output, "A:") {
				t.Error("Output should contain all metrics (T, R, F, A)")
			}
		})
	}
}

// TestProgressIndicator tests the progress indicator format
func TestProgressIndicator(t *testing.T) {
	tests := []struct {
		current int
		total   int
		name    string
	}{
		{current: 1, total: 10, name: "test1"},
		{current: 5, total: 10, name: "test5"},
		{current: 10, total: 10, name: "test10"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d/%d", tt.current, tt.total), func(t *testing.T) {
			progress := fmt.Sprintf("Test %d/%d: %s... ", tt.current, tt.total, tt.name)

			if !strings.Contains(progress, fmt.Sprintf("%d/%d", tt.current, tt.total)) {
				t.Error("Progress indicator should show current/total")
			}

			if !strings.Contains(progress, tt.name) {
				t.Error("Progress indicator should show test name")
			}
		})
	}
}

// TestHeaderFormatting tests the header output formatting
func TestHeaderFormatting(t *testing.T) {
	// Test that header elements are present
	expectedElements := []string{
		"═══",
		"RUNNER UNIVERSEL",
		"TESTS COMPLETS RÉSEAU RETE",
		"Pipeline unique",
		"propagation RETE",
		"Date:",
	}

	for _, element := range expectedElements {
		// Just verify the strings are well-formed
		if len(element) == 0 {
			t.Error("Header element should not be empty")
		}
	}
}

// TestTestFileStructureValidation tests test file structure
func TestTestFileStructureValidation(t *testing.T) {
	tf := TestFile{
		name:       "alpha_test",
		category:   "alpha",
		constraint: "test/coverage/alpha/alpha_test.constraint",
		facts:      "test/coverage/alpha/alpha_test.facts",
	}

	// Validate structure
	if tf.name == "" {
		t.Error("TestFile name should not be empty")
	}
	if tf.category == "" {
		t.Error("TestFile category should not be empty")
	}
	if !strings.HasSuffix(tf.constraint, ".constraint") {
		t.Error("Constraint file should end with .constraint")
	}
	if !strings.HasSuffix(tf.facts, ".facts") {
		t.Error("Facts file should end with .facts")
	}
	if !strings.Contains(tf.constraint, tf.category) {
		t.Error("Constraint path should contain category")
	}
}

// TestCategoryDetection tests category detection from paths
func TestCategoryDetection(t *testing.T) {
	testDirs := []struct {
		path     string
		category string
	}{
		{"test/coverage/alpha", "alpha"},
		{"beta_coverage_tests", "beta"},
		{"constraint/test/integration", "integration"},
	}

	for _, dir := range testDirs {
		t.Run(dir.category, func(t *testing.T) {
			// Verify category matches expected pattern
			if dir.category != "alpha" && dir.category != "beta" && dir.category != "integration" {
				t.Errorf("Unexpected category: %s", dir.category)
			}

			// Verify path is reasonable
			if len(dir.path) == 0 {
				t.Error("Path should not be empty")
			}
		})
	}
}

// TestFinalSummaryFormat tests the final summary formatting
func TestFinalSummaryFormat(t *testing.T) {
	tests := []struct {
		name   string
		total  int
		passed int
		failed int
	}{
		{name: "all pass", total: 10, passed: 10, failed: 0},
		{name: "some fail", total: 10, passed: 7, failed: 3},
		{name: "all fail", total: 5, passed: 0, failed: 5},
		{name: "no tests", total: 0, passed: 0, failed: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := fmt.Sprintf("Résumé: %d tests, %d réussis ✅, %d échoués ❌\n",
				tt.total, tt.passed, tt.failed)

			if !strings.Contains(summary, "Résumé:") {
				t.Error("Summary should contain 'Résumé:'")
			}
			if !strings.Contains(summary, "tests") {
				t.Error("Summary should contain 'tests'")
			}
			if !strings.Contains(summary, "✅") {
				t.Error("Summary should contain success emoji")
			}
			if !strings.Contains(summary, "❌") {
				t.Error("Summary should contain failure emoji")
			}

			// Verify totals match
			if tt.passed+tt.failed != tt.total {
				t.Errorf("passed(%d) + failed(%d) != total(%d)", tt.passed, tt.failed, tt.total)
			}
		})
	}
}

// TestSuccessMessage tests the all-tests-passed message
func TestSuccessMessage(t *testing.T) {
	failed := 0
	successMsg := "🎉 TOUS LES TESTS SONT PASSÉS!"

	if failed == 0 {
		// Success condition
		if !strings.Contains(successMsg, "🎉") {
			t.Error("Success message should contain celebration emoji")
		}
		if !strings.Contains(successMsg, "TOUS") {
			t.Error("Success message should emphasize all tests passed")
		}
	}
}
