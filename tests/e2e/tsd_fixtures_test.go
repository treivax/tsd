//go:build e2e

// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package e2e

import (
	"testing"
	"time"

	"github.com/treivax/tsd/tests/shared/testutil"
)

const (
	// ExpectedTotalFixtures nombre total de fixtures attendu
	ExpectedTotalFixtures = 83 // 26 alpha + 26 beta + 31 integration

	// ExpectedAlphaFixtures nombre de fixtures alpha minimum
	ExpectedAlphaFixtures = 26

	// ExpectedBetaFixtures nombre de fixtures beta minimum
	ExpectedBetaFixtures = 26

	// ExpectedIntegrationFixtures nombre de fixtures integration minimum
	ExpectedIntegrationFixtures = 31

	// DefaultTestTimeout timeout par défaut pour les tests
	DefaultTestTimeout = 30 * time.Second
)

// TestAlphaFixtures tests all alpha coverage fixtures
func TestAlphaFixtures(t *testing.T) {
	fixtures := testutil.GetFixturesByCategory(t, "alpha")

	if len(fixtures) == 0 {
		t.Skip("No alpha fixtures found")
	}

	t.Logf("Running %d alpha fixture tests", len(fixtures))

	for _, fixture := range fixtures {
		fixture := fixture // Capture for parallel execution
		t.Run(fixture.Name, func(t *testing.T) {
			t.Parallel()

			result := testutil.ExecuteTSDFile(t, fixture.Path)

			if fixture.ShouldError {
				testutil.AssertError(t, result)
			} else {
				testutil.AssertNoError(t, result)
				testutil.AssertMinNetworkStructure(t, result, 1, 1)

				t.Logf("✅ %s: T:%d R:%d F:%d A:%d",
					fixture.Name,
					result.TypeNodes,
					result.TerminalNodes,
					result.Facts,
					result.Activations)
			}
		})
	}
}

// TestBetaFixtures tests all beta coverage fixtures
func TestBetaFixtures(t *testing.T) {
	fixtures := testutil.GetFixturesByCategory(t, "beta")

	if len(fixtures) == 0 {
		t.Skip("No beta fixtures found")
	}

	t.Logf("Running %d beta fixture tests", len(fixtures))

	for _, fixture := range fixtures {
		fixture := fixture // Capture for parallel execution
		t.Run(fixture.Name, func(t *testing.T) {
			t.Parallel()

			var result *testutil.TSDResult
			if fixture.ShouldError {
				result = testutil.ExecuteTSDFileWithOptions(t, fixture.Path, &testutil.ExecutionOptions{
					ExpectError:     true,
					ValidateNetwork: false,
					CaptureOutput:   true,
					Timeout:         30 * time.Second,
				})
				testutil.AssertError(t, result)
			} else {
				result = testutil.ExecuteTSDFile(t, fixture.Path)
				testutil.AssertNoError(t, result)

				// reset_example.tsd is a documentation file with no rules
				// It only demonstrates the reset command
				if fixture.Name != "reset_example" {
					testutil.AssertMinNetworkStructure(t, result, 1, 1)
				} else {
					// For reset_example, just verify types exist
					if result.TypeNodes < 1 {
						t.Errorf("Expected at least 1 type node, got %d", result.TypeNodes)
					}
					t.Logf("📖 %s: Documentation file (no rules expected)", fixture.Name)
				}

				t.Logf("✅ %s: T:%d R:%d F:%d A:%d",
					fixture.Name,
					result.TypeNodes,
					result.TerminalNodes,
					result.Facts,
					result.Activations)
			}
		})
	}
}

// TestIntegrationFixtures tests all integration fixtures
func TestIntegrationFixtures(t *testing.T) {
	fixtures := testutil.GetFixturesByCategory(t, "integration")
	errorFixtures := testutil.GetErrorFixtures()

	if len(fixtures) == 0 {
		t.Skip("No integration fixtures found")
	}

	t.Logf("Running %d integration fixture tests", len(fixtures))

	for _, fixture := range fixtures {
		fixture := fixture // Capture for parallel execution
		t.Run(fixture.Name, func(t *testing.T) {
			t.Parallel()

			shouldError := errorFixtures[fixture.Name] || fixture.ShouldError

			var result *testutil.TSDResult
			if shouldError {
				result = testutil.ExecuteTSDFileWithOptions(t, fixture.Path, &testutil.ExecutionOptions{
					ExpectError:     true,
					ValidateNetwork: false,
					CaptureOutput:   true,
					Timeout:         DefaultTestTimeout,
				})
				testutil.AssertError(t, result)
				t.Logf("✅ %s: Error detected as expected", fixture.Name)
			} else {
				result = testutil.ExecuteTSDFile(t, fixture.Path)
				testutil.AssertNoError(t, result)

				// Integration tests may have varying network structures
				if result.TypeNodes > 0 && result.TerminalNodes > 0 {
					t.Logf("✅ %s: T:%d R:%d F:%d A:%d",
						fixture.Name,
						result.TypeNodes,
						result.TerminalNodes,
						result.Facts,
						result.Activations)
				} else {
					t.Logf("⚠️  %s: Network structure minimal or empty (may be expected)",
						fixture.Name)
				}
			}
		})
	}
}

// TestAllFixtures runs all fixtures in sequence and reports summary
func TestAllFixtures(t *testing.T) {
	testutil.SkipIfShort(t, "comprehensive fixture test skipped in short mode")

	allFixtures := testutil.GetAllFixtures(t)

	if len(allFixtures) == 0 {
		t.Fatal("No fixtures found - check fixture directory structure")
	}

	t.Logf("Running comprehensive test on %d fixtures", len(allFixtures))

	passed := 0
	failed := 0
	errorExpected := 0

	errorFixtures := testutil.GetErrorFixtures()

	for _, fixture := range allFixtures {
		shouldError := errorFixtures[fixture.Name] || fixture.ShouldError

		var result *testutil.TSDResult
		if shouldError {
			result = testutil.ExecuteTSDFileWithOptions(t, fixture.Path, &testutil.ExecutionOptions{
				ExpectError:     true,
				ValidateNetwork: false,
				CaptureOutput:   true,
				Timeout:         DefaultTestTimeout,
			})
		} else {
			result = testutil.ExecuteTSDFile(t, fixture.Path)
		}

		if shouldError {
			if result.Error != nil {
				errorExpected++
				t.Logf("✅ %s: Error detected as expected", fixture.Name)
			} else {
				failed++
				t.Logf("❌ %s: Expected error but got none", fixture.Name)
			}
		} else {
			if result.Error != nil {
				failed++
				t.Logf("❌ %s: %v", fixture.Name, result.Error)
			} else {
				passed++
				t.Logf("✅ %s: T:%d R:%d F:%d A:%d (%.2fms)",
					fixture.Name,
					result.TypeNodes,
					result.TerminalNodes,
					result.Facts,
					result.Activations,
					float64(result.Duration.Microseconds())/1000.0)
			}
		}
	}

	t.Logf("\n" + "═══════════════════════════════════════")
	t.Logf("Summary: %d total fixtures", len(allFixtures))
	t.Logf("  ✅ Passed: %d", passed)
	t.Logf("  ✅ Error expected: %d", errorExpected)
	t.Logf("  ❌ Failed: %d", failed)
	t.Logf("═══════════════════════════════════════")

	if failed > 0 {
		t.Errorf("%d fixtures failed execution", failed)
	}

	// Validate expected fixture count
	if len(allFixtures) != ExpectedTotalFixtures {
		t.Errorf("Expected %d fixtures, found %d", ExpectedTotalFixtures, len(allFixtures))
	}
}

// TestFixtureStructure validates the fixture directory structure
func TestFixtureStructure(t *testing.T) {
	if err := testutil.ValidateFixtureStructure(t); err != nil {
		t.Fatalf("Fixture structure validation failed: %v", err)
	}

	// Count fixtures per category
	alphaCount := len(testutil.GetFixturesByCategory(t, "alpha"))
	betaCount := len(testutil.GetFixturesByCategory(t, "beta"))
	integrationCount := len(testutil.GetFixturesByCategory(t, "integration"))

	t.Logf("Fixture counts:")
	t.Logf("  Alpha: %d", alphaCount)
	t.Logf("  Beta: %d", betaCount)
	t.Logf("  Integration: %d", integrationCount)
	t.Logf("  Total: %d", alphaCount+betaCount+integrationCount)

	// Validate minimum counts
	if alphaCount == 0 {
		t.Error("No alpha fixtures found")
	}
	if betaCount == 0 {
		t.Error("No beta fixtures found")
	}
	if integrationCount == 0 {
		t.Error("No integration fixtures found")
	}
}

// TestFixtureNaming validates fixture naming conventions
func TestFixtureNaming(t *testing.T) {
	allFixtures := testutil.GetAllFixtures(t)

	for _, fixture := range allFixtures {
		// Check for valid name
		if fixture.Name == "" {
			t.Errorf("Fixture has empty name: %s", fixture.Path)
		}

		// Check for valid path
		if fixture.Path == "" {
			t.Errorf("Fixture %s has empty path", fixture.Name)
		}

		// Check for valid category
		if fixture.Category == "" || fixture.Category == "unknown" {
			t.Errorf("Fixture %s has invalid category: %s", fixture.Name, fixture.Category)
		}
	}
}

// TestErrorFixtures specifically tests fixtures that should error
func TestErrorFixtures(t *testing.T) {
	errorFixtures := testutil.GetErrorFixtures()
	allFixtures := testutil.GetAllFixtures(t)

	for fixtureName := range errorFixtures {
		// Find the fixture
		var found *testutil.TSDFixture
		for _, fixture := range allFixtures {
			if fixture.Name == fixtureName {
				found = &fixture
				break
			}
		}

		if found == nil {
			t.Logf("Warning: Error fixture %s not found (may not be migrated yet)", fixtureName)
			continue
		}

		t.Run(fixtureName, func(t *testing.T) {
			result := testutil.ExecuteTSDFileWithOptions(t, found.Path, &testutil.ExecutionOptions{
				ExpectError:     true,
				ValidateNetwork: false,
				CaptureOutput:   true,
				Timeout:         DefaultTestTimeout,
			})
			testutil.AssertError(t, result)
			t.Logf("✅ %s: Error correctly detected", fixtureName)
		})
	}
}

// TestFixtureCache validates fixture caching works correctly
func TestFixtureCache(t *testing.T) {
	// Clear cache first
	testutil.ClearFixtureCache()

	// First discovery
	fixtures1 := testutil.GetAllFixtures(t)

	// Second discovery (should use cache)
	fixtures2 := testutil.GetAllFixtures(t)

	if len(fixtures1) != len(fixtures2) {
		t.Errorf("Cache inconsistency: first=%d, second=%d", len(fixtures1), len(fixtures2))
	}

	// Validate same fixtures
	for i, f1 := range fixtures1 {
		f2 := fixtures2[i]
		if f1.Name != f2.Name || f1.Path != f2.Path {
			t.Errorf("Fixture mismatch at index %d: %s vs %s", i, f1.Name, f2.Name)
		}
	}
}
