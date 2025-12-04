// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

var (
	// Cache for discovered fixtures to avoid repeated filesystem scans
	fixtureCache      = make(map[string][]TSDFixture)
	fixtureCacheMutex sync.RWMutex
)

// DiscoverFixtures finds all .tsd files in a directory recursively
func DiscoverFixtures(t *testing.T, baseDir string) []TSDFixture {
	t.Helper()

	// Check cache first
	fixtureCacheMutex.RLock()
	if cached, ok := fixtureCache[baseDir]; ok {
		fixtureCacheMutex.RUnlock()
		return cached
	}
	fixtureCacheMutex.RUnlock()

	var fixtures []TSDFixture

	// Get the project root (assuming we're in tests/shared/testutil)
	projectRoot := getProjectRoot()
	fullPath := filepath.Join(projectRoot, baseDir)

	// Check if directory exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		if t != nil {
			t.Logf("Warning: Directory %s does not exist", fullPath)
		}
		return fixtures
	}

	// Walk directory tree
	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .tsd files
		if !strings.HasSuffix(path, ".tsd") {
			return nil
		}

		// Determine category from path
		category := categorizeFixture(path)

		// Get base name without extension
		baseName := strings.TrimSuffix(filepath.Base(path), ".tsd")

		// Check if this is an error fixture
		shouldError := isErrorFixture(baseName)

		fixtures = append(fixtures, TSDFixture{
			Name:        baseName,
			Path:        path,
			Category:    category,
			ShouldError: shouldError,
		})

		return nil
	})

	if err != nil && t != nil {
		t.Errorf("Error discovering fixtures in %s: %v", baseDir, err)
	}

	// Cache the results
	fixtureCacheMutex.Lock()
	fixtureCache[baseDir] = fixtures
	fixtureCacheMutex.Unlock()

	return fixtures
}

// DiscoverFixturesWithPattern finds .tsd files matching a glob pattern
func DiscoverFixturesWithPattern(t *testing.T, pattern string) []TSDFixture {
	t.Helper()

	var fixtures []TSDFixture

	projectRoot := getProjectRoot()
	fullPattern := filepath.Join(projectRoot, pattern)

	matches, err := filepath.Glob(fullPattern)
	if err != nil {
		if t != nil {
			t.Errorf("Error matching pattern %s: %v", pattern, err)
		}
		return fixtures
	}

	for _, path := range matches {
		// Skip directories
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			continue
		}

		// Only process .tsd files
		if !strings.HasSuffix(path, ".tsd") {
			continue
		}

		category := categorizeFixture(path)
		baseName := strings.TrimSuffix(filepath.Base(path), ".tsd")
		shouldError := isErrorFixture(baseName)

		fixtures = append(fixtures, TSDFixture{
			Name:        baseName,
			Path:        path,
			Category:    category,
			ShouldError: shouldError,
		})
	}

	return fixtures
}

// GetFixturesByCategory returns fixtures for a specific category
func GetFixturesByCategory(t *testing.T, category string) []TSDFixture {
	t.Helper()

	baseDir := fmt.Sprintf("tests/fixtures/%s", category)
	return DiscoverFixtures(t, baseDir)
}

// GetErrorFixtures returns the set of fixtures that should produce errors
func GetErrorFixtures() map[string]bool {
	return map[string]bool{
		"error_args_test":      true,
		"invalid_no_types":     true,
		"invalid_unknown_type": true,
	}
}

// LoadFixture loads a single fixture by name
func LoadFixture(t *testing.T, name string) *TSDFixture {
	t.Helper()

	// Try to find in common locations
	searchPaths := []string{
		"tests/fixtures/alpha",
		"tests/fixtures/beta",
		"tests/fixtures/integration",
	}

	for _, searchPath := range searchPaths {
		fixtures := DiscoverFixtures(t, searchPath)
		for _, fixture := range fixtures {
			if fixture.Name == name {
				return &fixture
			}
		}
	}

	return nil
}

// FixtureExists checks if a fixture file exists
func FixtureExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetAllFixtures returns all fixtures from all categories
func GetAllFixtures(t *testing.T) []TSDFixture {
	t.Helper()

	var allFixtures []TSDFixture

	categories := []string{"alpha", "beta", "integration"}

	for _, category := range categories {
		fixtures := GetFixturesByCategory(t, category)
		allFixtures = append(allFixtures, fixtures...)
	}

	return allFixtures
}

// categorizeFixture determines the category from the file path
func categorizeFixture(path string) string {
	if strings.Contains(path, "alpha") {
		return "alpha"
	}
	if strings.Contains(path, "beta") {
		return "beta"
	}
	if strings.Contains(path, "integration") {
		return "integration"
	}
	return "unknown"
}

// isErrorFixture checks if a fixture name is in the error list
func isErrorFixture(name string) bool {
	errorFixtures := GetErrorFixtures()
	return errorFixtures[name]
}

// getProjectRoot returns the project root directory
func getProjectRoot() string {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}

	// Walk up until we find go.mod
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root without finding go.mod
			return wd
		}
		dir = parent
	}
}

// ClearFixtureCache clears the fixture cache (useful for testing)
func ClearFixtureCache() {
	fixtureCacheMutex.Lock()
	defer fixtureCacheMutex.Unlock()
	fixtureCache = make(map[string][]TSDFixture)
}

// GetFixtureCount returns the total number of discovered fixtures
func GetFixtureCount(t *testing.T) int {
	t.Helper()
	fixtures := GetAllFixtures(t)
	return len(fixtures)
}

// GetFixturesByPattern returns fixtures whose names match a pattern
func GetFixturesByPattern(t *testing.T, namePattern string) []TSDFixture {
	t.Helper()

	allFixtures := GetAllFixtures(t)
	var matching []TSDFixture

	for _, fixture := range allFixtures {
		if strings.Contains(fixture.Name, namePattern) {
			matching = append(matching, fixture)
		}
	}

	return matching
}

// ValidateFixtureStructure checks that the fixture directory structure is correct
func ValidateFixtureStructure(t *testing.T) error {
	t.Helper()

	projectRoot := getProjectRoot()
	requiredDirs := []string{
		"tests/fixtures/alpha",
		"tests/fixtures/beta",
		"tests/fixtures/integration",
	}

	for _, dir := range requiredDirs {
		fullPath := filepath.Join(projectRoot, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("required directory %s does not exist", dir)
		}
	}

	return nil
}
