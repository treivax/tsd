// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package testutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// GetTestDataPath returns the path to the test data directory
func GetTestDataPath() string {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "tests"
	}

	// Navigate up to find project root (contains go.mod)
	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, "tests")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root without finding go.mod
			return "tests"
		}
		dir = parent
	}
}

// AssertNoError fails the test if result contains an error
func AssertNoError(t *testing.T, result *TSDResult) {
	t.Helper()

	if result == nil {
		t.Fatal("Result is nil")
	}

	if result.Error != nil {
		t.Fatalf("Expected no error, got: %v", result.Error)
	}
}

// AssertNetworkStructure verifies network structure meets minimum requirements
func AssertNetworkStructure(t *testing.T, result *TSDResult, minTypeNodes, minTerminalNodes int) {
	t.Helper()

	if result == nil {
		t.Fatal("Result is nil")
	}

	if result.TypeNodes < minTypeNodes {
		t.Errorf("Expected at least %d type nodes, got %d", minTypeNodes, result.TypeNodes)
	}

	if result.TerminalNodes < minTerminalNodes {
		t.Errorf("Expected at least %d terminal nodes, got %d", minTerminalNodes, result.TerminalNodes)
	}
}

// AssertFactCount verifies the number of facts matches expected count
func AssertFactCount(t *testing.T, result *TSDResult, expectedCount int) {
	t.Helper()

	if result == nil {
		t.Fatal("Result is nil")
	}

	if result.Facts != expectedCount {
		t.Errorf("Expected %d facts, got %d", expectedCount, result.Facts)
	}
}

// AssertMinActivations verifies at least the minimum number of activations occurred
func AssertMinActivations(t *testing.T, result *TSDResult, minActivations int) {
	t.Helper()

	if result == nil {
		t.Fatal("Result is nil")
	}

	if result.Activations < minActivations {
		t.Errorf("Expected at least %d activations, got %d", minActivations, result.Activations)
	}
}

// AssertMinNetworkStructure is an alias for AssertNetworkStructure
func AssertMinNetworkStructure(t *testing.T, result *TSDResult, minTypeNodes, minTerminalNodes int) {
	AssertNetworkStructure(t, result, minTypeNodes, minTerminalNodes)
}

// AssertError fails the test if result does not contain an error
func AssertError(t *testing.T, result *TSDResult) {
	t.Helper()

	if result == nil {
		t.Fatal("Result is nil")
	}

	if result.Error == nil {
		t.Fatal("Expected an error, got nil")
	}
}

// SkipIfShort skips test if -short flag is set
func SkipIfShort(t *testing.T, reason string) {
	t.Helper()

	if testing.Short() {
		if reason != "" {
			t.Skipf("Skipping test in short mode: %s", reason)
		} else {
			t.Skip("Skipping test in short mode")
		}
	}
}

// GetFixturesByCategory returns fixtures filtered by category
func GetFixturesByCategory(t *testing.T, category string) []TSDFixture {
	t.Helper()

	testDataPath := GetTestDataPath()
	fixturesPath := filepath.Join(testDataPath, "fixtures", category)

	var fixtures []TSDFixture

	entries, err := os.ReadDir(fixturesPath)
	if err != nil {
		// Directory doesn't exist or can't be read, return empty
		return fixtures
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tsd") {
			continue
		}

		fixture := TSDFixture{
			Name:        strings.TrimSuffix(entry.Name(), ".tsd"),
			Path:        filepath.Join(fixturesPath, entry.Name()),
			Category:    category,
			ShouldError: strings.Contains(entry.Name(), "error") || strings.Contains(entry.Name(), "invalid"),
		}
		fixtures = append(fixtures, fixture)
	}

	return fixtures
}

// GetErrorFixtures returns fixtures that are expected to produce errors
func GetErrorFixtures() map[string]bool {
	testDataPath := GetTestDataPath()
	fixturesPath := filepath.Join(testDataPath, "fixtures", "errors")

	errorMap := make(map[string]bool)

	entries, err := os.ReadDir(fixturesPath)
	if err != nil {
		return errorMap
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tsd") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".tsd")
		errorMap[name] = true
	}

	return errorMap
}

// GetAllFixtures returns all available fixtures
func GetAllFixtures(t *testing.T) []TSDFixture {
	t.Helper()

	var allFixtures []TSDFixture

	categories := []string{"alpha", "beta", "integration", "errors"}
	for _, category := range categories {
		fixtures := GetFixturesByCategory(t, category)
		allFixtures = append(allFixtures, fixtures...)
	}

	return allFixtures
}

// ClearFixtureCache clears any cached fixture data
func ClearFixtureCache() {
	// No-op for now
}

// ValidateFixtureStructure validates fixture file structure
func ValidateFixtureStructure(t *testing.T) error {
	t.Helper()

	// Validate that fixtures directory exists and has expected structure
	testDataPath := GetTestDataPath()
	fixturesPath := filepath.Join(testDataPath, "fixtures")

	if _, err := os.Stat(fixturesPath); err != nil {
		return err
	}

	// Check that expected categories exist
	categories := []string{"alpha", "beta", "integration", "errors"}
	for _, category := range categories {
		categoryPath := filepath.Join(fixturesPath, category)
		if _, err := os.Stat(categoryPath); err != nil {
			// Category doesn't exist, which is okay - some might not be present yet
			continue
		}
	}

	return nil
}

// CreateTempTSDFile creates a temporary .tsd file for testing
func CreateTempTSDFile(t *testing.T, content string) string {
	t.Helper()

	// Create temp file with .tsd extension
	tmpFile, err := os.CreateTemp("", "test-*.tsd")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	path := tmpFile.Name()

	// Write content
	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		os.Remove(path)
		t.Fatalf("Failed to write temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(path)
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Register cleanup
	t.Cleanup(func() {
		os.Remove(path)
	})

	return path
}

// CreateTempTSDFileWithName creates a temp .tsd file with specific name
func CreateTempTSDFileWithName(t *testing.T, name, content string) string {
	t.Helper()

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "tsd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create file with specific name
	path := filepath.Join(tmpDir, name)
	if !strings.HasSuffix(path, ".tsd") {
		path += ".tsd"
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to write temp file: %v", err)
	}

	// Register cleanup
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	return path
}

// CleanupTempFiles removes temporary test files
func CleanupTempFiles(t *testing.T, paths ...string) {
	t.Helper()

	for _, path := range paths {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			t.Logf("Warning: Failed to cleanup %s: %v", path, err)
		}
	}
}
