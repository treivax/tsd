# Feature Implementation: Complete Test Restructuring with Go Test Integration

## Objective

Restructure all tests in the TSD project to adopt a unified `go test` approach, eliminating the custom `universal-rete-runner` and centralizing test organization. Create a clean separation between unit tests (in modules) and E2E/integration tests (in centralized `tests/` directory).

## Context

The TSD project currently has:
- 125+ `*_test.go` files scattered across modules (`rete/`, `constraint/`, `cmd/`)
- 83 `.tsd` fixture files in 3 different locations
- A custom `universal-rete-runner` binary that discovers and executes `.tsd` files
- Mixed test types (unit, integration, E2E) without clear separation

This restructuring will:
- Centralize E2E and integration tests in a `tests/` directory at project root
- Keep unit tests alongside source code (Go convention)
- Convert all test execution to standard `go test` with build tags
- Remove the custom runner and consolidate logic into reusable test helpers
- Provide clear Makefile commands for different test types

## Architecture Overview

```
tsd/
├── constraint/                     # Module - keep unit tests here
│   ├── *.go                       
│   └── *_test.go                  # Unit tests stay
│
├── rete/                          # Module - keep unit tests here
│   ├── *.go
│   └── *_test.go                  # Unit tests stay
│
├── cmd/                           # Commands - keep unit tests here
│   └── tsd/
│       ├── main.go
│       └── main_test.go           # CLI unit tests stay
│
├── tests/                         # NEW: Centralized E2E/integration
│   ├── README.md                  # Test documentation
│   │
│   ├── fixtures/                  # Centralized .tsd files
│   │   ├── alpha/                 # 26 alpha coverage tests
│   │   ├── beta/                  # 26 beta coverage tests
│   │   └── integration/           # 31 integration tests
│   │
│   ├── e2e/                       # End-to-end tests
│   │   ├── tsd_fixtures_test.go   # Table-driven .tsd tests
│   │   ├── cli_test.go            # CLI binary tests
│   │   └── scenarios_test.go      # Complete workflows
│   │
│   ├── integration/               # Integration tests
│   │   ├── constraint_rete_test.go
│   │   ├── pipeline_test.go
│   │   └── incremental_test.go
│   │
│   ├── performance/               # Performance/load tests
│   │   ├── load_test.go
│   │   ├── stress_test.go
│   │   └── benchmark_test.go
│   │
│   └── shared/                    # Shared test utilities
│       └── testutil/
│           ├── runner.go          # TSD file execution helper
│           ├── fixtures.go        # Fixture discovery
│           ├── assertions.go      # Custom assertions
│           └── helpers.go         # Common test helpers
│
└── Makefile                       # Unified test commands
```

## Implementation Tasks

### Phase 1: Create Test Infrastructure

#### Task 1.1: Create Test Directory Structure
**Action:** Create new directories

Create:
- `tests/`
- `tests/fixtures/alpha/`
- `tests/fixtures/beta/`
- `tests/fixtures/integration/`
- `tests/e2e/`
- `tests/integration/`
- `tests/performance/`
- `tests/shared/testutil/`

**Acceptance Criteria:**
- All directories created
- `.gitkeep` files added where needed
- Structure matches architecture diagram

#### Task 1.2: Implement Test Utilities - Runner
**File:** `tests/shared/testutil/runner.go`

Implement core test execution helper:

```go
package testutil

import (
    "testing"
    "github.com/treivax/tsd/rete"
)

// TSDFixture represents a .tsd test file
type TSDFixture struct {
    Name        string
    Path        string
    Category    string
    ShouldError bool
}

// TSDResult contains execution results
type TSDResult struct {
    TypeNodes     int
    TerminalNodes int
    Facts         int
    Activations   int
    Error         error
    Output        string
}

// ExecuteTSDFile executes a .tsd file and returns results
func ExecuteTSDFile(t *testing.T, path string) *TSDResult

// ExecuteTSDFileWithOptions executes with custom configuration
func ExecuteTSDFileWithOptions(t *testing.T, path string, opts *ExecutionOptions) *TSDResult

// ExecutionOptions configures test execution
type ExecutionOptions struct {
    ExpectError      bool
    MinActivations   int
    MaxActivations   int
    ValidateNetwork  bool
    CaptureOutput    bool
    Timeout          time.Duration
}

// AssertTSDResult validates a TSD execution result
func AssertTSDResult(t *testing.T, result *TSDResult, expected *TSDResult)

// RunTSDFile is a convenience wrapper
func RunTSDFile(t *testing.T, path string) {
    t.Helper()
    result := ExecuteTSDFile(t, path)
    if result.Error != nil {
        t.Fatalf("Failed to execute %s: %v", path, result.Error)
    }
}
```

**Implementation Details:**
- Reuse logic from `universal-rete-runner/main.go`
- Use `rete.NewConstraintPipeline()` and `IngestFile()`
- Capture stdout/stderr for output validation
- Thread-safe execution for parallel tests
- Proper cleanup of resources

**Acceptance Criteria:**
- Executes .tsd files correctly
- Returns detailed results
- Handles errors gracefully
- Works with `t.Parallel()`
- No resource leaks

#### Task 1.3: Implement Test Utilities - Fixture Discovery
**File:** `tests/shared/testutil/fixtures.go`

Implement fixture discovery and management:

```go
package testutil

import (
    "path/filepath"
    "testing"
)

// DiscoverFixtures finds all .tsd files in a directory
func DiscoverFixtures(t *testing.T, baseDir string) []TSDFixture

// DiscoverFixturesWithPattern finds .tsd files matching pattern
func DiscoverFixturesWithPattern(t *testing.T, pattern string) []TSDFixture

// GetFixturesByCategory returns fixtures for a specific category
func GetFixturesByCategory(t *testing.T, category string) []TSDFixture

// GetErrorFixtures returns fixtures that should produce errors
func GetErrorFixtures() map[string]bool

// LoadFixture loads a single fixture by name
func LoadFixture(t *testing.T, name string) *TSDFixture

// FixtureExists checks if a fixture exists
func FixtureExists(path string) bool
```

**Implementation Details:**
- Use `filepath.Glob()` for discovery
- Cache discovered fixtures for performance
- Support recursive directory scanning
- Handle missing directories gracefully
- Error fixtures: `error_args_test`, `invalid_no_types`, `invalid_unknown_type`

**Acceptance Criteria:**
- Discovers all .tsd files correctly
- Categorizes fixtures properly
- Handles missing directories
- Performance: < 100ms for 100 fixtures

#### Task 1.4: Implement Test Utilities - Assertions
**File:** `tests/shared/testutil/assertions.go`

Implement TSD-specific assertions:

```go
package testutil

import "testing"

// AssertNetworkStructure validates RETE network structure
func AssertNetworkStructure(t *testing.T, result *TSDResult, expectedTypeNodes, expectedTerminalNodes int)

// AssertActivations validates activation count
func AssertActivations(t *testing.T, result *TSDResult, expected int)

// AssertMinActivations validates minimum activations
func AssertMinActivations(t *testing.T, result *TSDResult, min int)

// AssertNoError validates successful execution
func AssertNoError(t *testing.T, result *TSDResult)

// AssertError validates expected error
func AssertError(t *testing.T, result *TSDResult)

// AssertErrorContains validates error message
func AssertErrorContains(t *testing.T, result *TSDResult, substring string)

// AssertOutputContains validates captured output
func AssertOutputContains(t *testing.T, result *TSDResult, substring string)

// AssertFactCount validates fact count in storage
func AssertFactCount(t *testing.T, result *TSDResult, expected int)
```

**Acceptance Criteria:**
- Clear error messages with context
- Helper methods use `t.Helper()`
- Integrates with Go's testing package
- Provides detailed failure information

#### Task 1.5: Implement Test Utilities - Helpers
**File:** `tests/shared/testutil/helpers.go`

Implement common test helpers:

```go
package testutil

import "testing"

// WithTimeout runs a test function with timeout
func WithTimeout(t *testing.T, duration time.Duration, fn func())

// CreateTempTSDFile creates a temporary .tsd file for testing
func CreateTempTSDFile(t *testing.T, content string) string

// CleanupTempFiles removes temporary test files
func CleanupTempFiles(t *testing.T, paths ...string)

// CaptureStdout captures stdout during function execution
func CaptureStdout(fn func()) string

// CompareResults compares two TSD execution results
func CompareResults(t *testing.T, result1, result2 *TSDResult)

// SkipIfShort skips test if -short flag is set
func SkipIfShort(t *testing.T, reason string)

// GetTestDataPath returns path to test data directory
func GetTestDataPath() string
```

**Acceptance Criteria:**
- Proper cleanup with `t.Cleanup()`
- Thread-safe implementations
- Cross-platform compatibility
- Clear documentation

### Phase 2: Migrate Fixtures

#### Task 2.1: Move Alpha Coverage Tests
**Action:** Move and reorganize files

Move from:
- `test/coverage/alpha/*.tsd` (26 files)

To:
- `tests/fixtures/alpha/*.tsd`

**Script:**
```bash
mkdir -p tests/fixtures/alpha
mv test/coverage/alpha/*.tsd tests/fixtures/alpha/
```

**Acceptance Criteria:**
- All 26 files moved
- No files lost
- File permissions preserved

#### Task 2.2: Move Beta Coverage Tests
**Action:** Move and reorganize files

Move from:
- `beta_coverage_tests/*.tsd` (26 files)

To:
- `tests/fixtures/beta/*.tsd`

**Script:**
```bash
mkdir -p tests/fixtures/beta
mv beta_coverage_tests/*.tsd tests/fixtures/beta/
```

**Acceptance Criteria:**
- All 26 files moved
- Original directory can be removed
- Fixtures accessible from new location

#### Task 2.3: Move Integration Test Fixtures
**Action:** Move and reorganize files

Move from:
- `constraint/test/integration/*.tsd` (31 files)

To:
- `tests/fixtures/integration/*.tsd`

**Script:**
```bash
mkdir -p tests/fixtures/integration
mv constraint/test/integration/*.tsd tests/fixtures/integration/
```

**Acceptance Criteria:**
- All 31 files moved
- Total: 83 .tsd files in tests/fixtures/
- Old test directories cleaned up

### Phase 3: Create E2E Tests

#### Task 3.1: Implement Table-Driven TSD Tests
**File:** `tests/e2e/tsd_fixtures_test.go`

Create comprehensive table-driven tests:

```go
//go:build e2e

package e2e

import (
    "testing"
    "github.com/treivax/tsd/tests/shared/testutil"
)

func TestAlphaFixtures(t *testing.T) {
    fixtures := testutil.GetFixturesByCategory(t, "alpha")
    
    for _, fixture := range fixtures {
        t.Run(fixture.Name, func(t *testing.T) {
            t.Parallel()
            
            result := testutil.ExecuteTSDFile(t, fixture.Path)
            
            testutil.AssertNoError(t, result)
            testutil.AssertNetworkStructure(t, result, 1, 1) // At least 1 type node and 1 terminal
        })
    }
}

func TestBetaFixtures(t *testing.T) {
    fixtures := testutil.GetFixturesByCategory(t, "beta")
    
    for _, fixture := range fixtures {
        t.Run(fixture.Name, func(t *testing.T) {
            t.Parallel()
            
            result := testutil.ExecuteTSDFile(t, fixture.Path)
            
            testutil.AssertNoError(t, result)
        })
    }
}

func TestIntegrationFixtures(t *testing.T) {
    fixtures := testutil.GetFixturesByCategory(t, "integration")
    errorFixtures := testutil.GetErrorFixtures()
    
    for _, fixture := range fixtures {
        t.Run(fixture.Name, func(t *testing.T) {
            t.Parallel()
            
            result := testutil.ExecuteTSDFile(t, fixture.Path)
            
            if errorFixtures[fixture.Name] {
                testutil.AssertError(t, result)
            } else {
                testutil.AssertNoError(t, result)
            }
        })
    }
}

// Test all fixtures in one run (for coverage)
func TestAllFixtures(t *testing.T) {
    testutil.SkipIfShort(t, "comprehensive test skipped in short mode")
    
    fixtures := testutil.DiscoverFixtures(t, "fixtures")
    
    passed := 0
    failed := 0
    
    for _, fixture := range fixtures {
        result := testutil.ExecuteTSDFile(t, fixture.Path)
        
        if result.Error != nil {
            failed++
            t.Logf("❌ %s: %v", fixture.Name, result.Error)
        } else {
            passed++
            t.Logf("✅ %s: T:%d R:%d F:%d A:%d",
                fixture.Name,
                result.TypeNodes,
                result.TerminalNodes,
                result.Facts,
                result.Activations)
        }
    }
    
    t.Logf("Summary: %d passed, %d failed out of %d total", passed, failed, len(fixtures))
    
    if failed > 0 {
        t.Errorf("%d fixtures failed", failed)
    }
}
```

**Acceptance Criteria:**
- All 83 fixtures tested
- Tests can run in parallel
- Clear test output with statistics
- Errors reported with context

#### Task 3.2: Implement CLI Tests
**File:** `tests/e2e/cli_test.go`

Test the TSD CLI binary:

```go
//go:build e2e

package e2e

import (
    "os/exec"
    "testing"
    "path/filepath"
)

func TestCLI_BasicExecution(t *testing.T) {
    fixture := filepath.Join(testutil.GetTestDataPath(), "fixtures/alpha/alpha_abs_positive.tsd")
    
    cmd := exec.Command("../../bin/tsd", "-file", fixture)
    output, err := cmd.CombinedOutput()
    
    if err != nil {
        t.Fatalf("CLI execution failed: %v\nOutput: %s", err, output)
    }
    
    if !strings.Contains(string(output), "success") {
        t.Errorf("Expected success message in output: %s", output)
    }
}

func TestCLI_InvalidFile(t *testing.T) {
    cmd := exec.Command("../../bin/tsd", "-file", "nonexistent.tsd")
    output, err := cmd.CombinedOutput()
    
    if err == nil {
        t.Error("Expected error for nonexistent file")
    }
    
    if !strings.Contains(string(output), "error") && !strings.Contains(string(output), "not found") {
        t.Errorf("Expected error message in output: %s", output)
    }
}

func TestCLI_Help(t *testing.T) {
    cmd := exec.Command("../../bin/tsd", "-help")
    output, err := cmd.CombinedOutput()
    
    if err != nil {
        t.Fatalf("Help command failed: %v", err)
    }
    
    if !strings.Contains(string(output), "Usage") {
        t.Error("Expected usage information in help output")
    }
}

func TestCLI_Version(t *testing.T) {
    cmd := exec.Command("../../bin/tsd", "-version")
    output, err := cmd.CombinedOutput()
    
    if err != nil {
        t.Fatalf("Version command failed: %v", err)
    }
    
    if len(output) == 0 {
        t.Error("Expected version information")
    }
}
```

**Acceptance Criteria:**
- CLI binary tested end-to-end
- Exit codes validated
- Output validated
- Error cases handled

#### Task 3.3: Implement Scenario Tests
**File:** `tests/e2e/scenarios_test.go`

Test complete user workflows:

```go
//go:build e2e

package e2e

import (
    "testing"
    "github.com/treivax/tsd/rete"
    "github.com/treivax/tsd/tests/shared/testutil"
)

func TestScenario_CompleteWorkflow(t *testing.T) {
    testutil.SkipIfShort(t, "scenario tests skipped in short mode")
    
    // Step 1: Load rules
    pipeline := rete.NewConstraintPipeline()
    storage := rete.NewMemoryStorage()
    
    fixture := filepath.Join(testutil.GetTestDataPath(), "fixtures/integration/actions.tsd")
    network, err := pipeline.IngestFile(fixture, nil, storage)
    
    if err != nil {
        t.Fatalf("Failed to load rules: %v", err)
    }
    
    // Step 2: Verify network structure
    if len(network.TypeNodes) == 0 {
        t.Error("Expected type nodes to be created")
    }
    
    if len(network.TerminalNodes) == 0 {
        t.Error("Expected terminal nodes to be created")
    }
    
    // Step 3: Verify facts loaded
    facts := storage.GetAllFacts()
    if len(facts) == 0 {
        t.Error("Expected facts to be loaded")
    }
    
    // Step 4: Verify activations
    activations := 0
    for _, terminal := range network.TerminalNodes {
        if terminal.Memory != nil && terminal.Memory.Tokens != nil {
            activations += len(terminal.Memory.Tokens)
        }
    }
    
    if activations == 0 {
        t.Log("Warning: No activations found (may be expected for some tests)")
    }
    
    t.Logf("Workflow completed: %d types, %d terminals, %d facts, %d activations",
        len(network.TypeNodes), len(network.TerminalNodes), len(facts), activations)
}

func TestScenario_IncrementalFacts(t *testing.T) {
    testutil.SkipIfShort(t, "scenario tests skipped in short mode")
    
    // Test adding facts incrementally
    pipeline := rete.NewConstraintPipeline()
    storage := rete.NewMemoryStorage()
    
    // Load rules only
    ruleFile := testutil.CreateTempTSDFile(t, `
type Person(name: string, age: number)
action person_found(name: string)
rule r1 : {p: Person} / p.age > 18 ==> person_found(p.name)
`)
    
    network, err := pipeline.IngestFile(ruleFile, nil, storage)
    if err != nil {
        t.Fatalf("Failed to load rules: %v", err)
    }
    
    // Add facts one by one
    facts := []string{
        `Person(name:"Alice", age:25)`,
        `Person(name:"Bob", age:17)`,
        `Person(name:"Charlie", age:30)`,
    }
    
    for _, factStr := range facts {
        // In real scenario, would add facts incrementally
        // For now, just validate the approach
        t.Logf("Would add fact: %s", factStr)
    }
    
    t.Log("Incremental fact scenario completed")
}
```

**Acceptance Criteria:**
- Real-world workflows tested
- Multi-step scenarios validated
- Edge cases covered
- Clear scenario documentation

### Phase 4: Create Integration Tests

#### Task 4.1: Implement Module Integration Tests
**File:** `tests/integration/constraint_rete_test.go`

Test interaction between constraint and rete modules:

```go
//go:build integration

package integration

import (
    "testing"
    "github.com/treivax/tsd/constraint"
    "github.com/treivax/tsd/rete"
)

func TestConstraintReteIntegration(t *testing.T) {
    // Test that constraint validation works with rete execution
    
    pipeline := rete.NewConstraintPipeline()
    storage := rete.NewMemoryStorage()
    
    // Valid constraint
    validRule := `
type Person(name: string, age: number)
rule r1 : {p: Person} / p.age > 18 ==> print("adult")
Person(name:"Alice", age:25)
`
    
    tempFile := testutil.CreateTempTSDFile(t, validRule)
    defer testutil.CleanupTempFiles(t, tempFile)
    
    network, err := pipeline.IngestFile(tempFile, nil, storage)
    
    if err != nil {
        t.Fatalf("Expected valid rule to succeed: %v", err)
    }
    
    if network == nil {
        t.Fatal("Expected network to be created")
    }
}

func TestInvalidConstraintRejection(t *testing.T) {
    pipeline := rete.NewConstraintPipeline()
    storage := rete.NewMemoryStorage()
    
    // Invalid constraint - unknown type
    invalidRule := `
rule r1 : {p: UnknownType} / p.age > 18 ==> print("test")
`
    
    tempFile := testutil.CreateTempTSDFile(t, invalidRule)
    defer testutil.CleanupTempFiles(t, tempFile)
    
    _, err := pipeline.IngestFile(tempFile, nil, storage)
    
    if err == nil {
        t.Error("Expected error for invalid constraint")
    }
}
```

**Acceptance Criteria:**
- Tests cross-module interactions
- Validates integration points
- Tests both success and error paths
- Clear integration boundaries

#### Task 4.2: Implement Pipeline Integration Tests
**File:** `tests/integration/pipeline_test.go`

Test complete pipeline flows:

```go
//go:build integration

package integration

import (
    "testing"
    "github.com/treivax/tsd/rete"
    "github.com/treivax/tsd/tests/shared/testutil"
)

func TestPipeline_CompleteFlow(t *testing.T) {
    result := testutil.ExecuteTSDFile(t, 
        filepath.Join(testutil.GetTestDataPath(), "fixtures/integration/actions.tsd"))
    
    testutil.AssertNoError(t, result)
    testutil.AssertNetworkStructure(t, result, 1, 1)
    
    if result.Facts == 0 {
        t.Error("Expected facts to be loaded")
    }
}

func TestPipeline_MultipleRules(t *testing.T) {
    rule := `
type Person(name: string, age: number)
rule r1 : {p: Person} / p.age > 18 ==> print("adult")
rule r2 : {p: Person} / p.age <= 18 ==> print("minor")
Person(name:"Alice", age:25)
Person(name:"Bob", age:15)
`
    
    tempFile := testutil.CreateTempTSDFile(t, rule)
    defer testutil.CleanupTempFiles(t, tempFile)
    
    result := testutil.ExecuteTSDFile(t, tempFile)
    
    testutil.AssertNoError(t, result)
    
    if result.TerminalNodes < 2 {
        t.Errorf("Expected at least 2 terminal nodes, got %d", result.TerminalNodes)
    }
}
```

**Acceptance Criteria:**
- Pipeline tested end-to-end
- Multiple rules handled
- Facts and activations validated

### Phase 5: Create Performance Tests

#### Task 5.1: Implement Load Tests
**File:** `tests/performance/load_test.go`

Test system under load:

```go
//go:build performance

package performance

import (
    "testing"
    "github.com/treivax/tsd/tests/shared/testutil"
)

func TestLoad_100Facts(t *testing.T) {
    testutil.SkipIfShort(t, "performance tests skipped in short mode")
    
    // Generate rule with 100 facts
    rule := generateRuleWithFacts(100)
    
    tempFile := testutil.CreateTempTSDFile(t, rule)
    defer testutil.CleanupTempFiles(t, tempFile)
    
    result := testutil.ExecuteTSDFile(t, tempFile)
    
    testutil.AssertNoError(t, result)
    testutil.AssertFactCount(t, result, 100)
}

func TestLoad_1000Facts(t *testing.T) {
    testutil.SkipIfShort(t, "performance tests skipped in short mode")
    
    rule := generateRuleWithFacts(1000)
    
    tempFile := testutil.CreateTempTSDFile(t, rule)
    defer testutil.CleanupTempFiles(t, tempFile)
    
    result := testutil.ExecuteTSDFileWithOptions(t, tempFile, &testutil.ExecutionOptions{
        Timeout: 30 * time.Second,
    })
    
    testutil.AssertNoError(t, result)
    testutil.AssertFactCount(t, result, 1000)
}

func generateRuleWithFacts(count int) string {
    // Generate TSD with specified number of facts
}
```

**Acceptance Criteria:**
- Tests with varying load levels
- Performance baselines established
- Timeout handling

#### Task 5.2: Implement Benchmarks
**File:** `tests/performance/benchmark_test.go`

Create Go benchmarks:

```go
//go:build performance

package performance

import (
    "testing"
    "github.com/treivax/tsd/rete"
)

func BenchmarkTSDExecution_Simple(b *testing.B) {
    fixture := filepath.Join(testutil.GetTestDataPath(), "fixtures/alpha/alpha_abs_positive.tsd")
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        pipeline := rete.NewConstraintPipeline()
        storage := rete.NewMemoryStorage()
        _, _ = pipeline.IngestFile(fixture, nil, storage)
    }
}

func BenchmarkTSDExecution_Complex(b *testing.B) {
    fixture := filepath.Join(testutil.GetTestDataPath(), "fixtures/integration/alpha_exhaustive_coverage.tsd")
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        pipeline := rete.NewConstraintPipeline()
        storage := rete.NewMemoryStorage()
        _, _ = pipeline.IngestFile(fixture, nil, storage)
    }
}

func BenchmarkParallel(b *testing.B) {
    fixture := filepath.Join(testutil.GetTestDataPath(), "fixtures/alpha/alpha_abs_positive.tsd")
    
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            pipeline := rete.NewConstraintPipeline()
            storage := rete.NewMemoryStorage()
            _, _ = pipeline.IngestFile(fixture, nil, storage)
        }
    })
}
```

**Acceptance Criteria:**
- Benchmarks for different fixture types
- Memory allocation tracking
- Parallel execution benchmarks

### Phase 6: Update Documentation

#### Task 6.1: Create Tests README
**File:** `tests/README.md`

Document test organization and usage:

```markdown
# TSD Test Suite

Comprehensive test suite for the TSD rule engine using standard `go test`.

## Structure

- `fixtures/` - Test data files (.tsd)
  - `alpha/` - Alpha node coverage tests (26 tests)
  - `beta/` - Beta node coverage tests (26 tests)
  - `integration/` - Integration test fixtures (31 tests)
  
- `e2e/` - End-to-end tests
- `integration/` - Module integration tests
- `performance/` - Load and benchmark tests
- `shared/testutil/` - Shared test utilities

## Running Tests

### Unit Tests (in modules)
```bash
make test-unit
```

### E2E Tests
```bash
make test-e2e
```

### Integration Tests
```bash
make test-integration
```

### All Tests
```bash
make test-all
```

### With Coverage
```bash
make coverage
```

## Build Tags

- `e2e` - End-to-end tests
- `integration` - Integration tests
- `performance` - Performance/load tests

## Parallel Execution

Tests are designed to run in parallel:
```bash
go test -parallel=8 -tags=e2e ./tests/e2e/...
```
```

**Acceptance Criteria:**
- Clear structure documentation
- Usage examples
- Build tag explanations
- Contributing guidelines

#### Task 6.2: Update Project README
**File:** `README.md`

Update main README with new test information:

```markdown
## Testing

The project uses standard Go testing with organized test suites:

- **Unit Tests**: `go test ./...` or `make test-unit`
- **E2E Tests**: `go test -tags=e2e ./tests/e2e/...` or `make test-e2e`
- **Integration Tests**: `go test -tags=integration ./tests/integration/...`
- **Performance Tests**: `go test -tags=performance ./tests/performance/...`

See [tests/README.md](tests/README.md) for detailed information.

### Quick Start

```bash
# Run all unit tests
make test-unit

# Run all tests including E2E
make test-all

# Run with coverage
make coverage

# Run specific test category
make test-e2e
```
```

**Acceptance Criteria:**
- README updated with test commands
- Link to detailed test documentation
- Quick start examples

### Phase 7: Create Makefile Commands

#### Task 7.1: Implement Makefile
**File:** `Makefile` (update existing)

Add/update test targets:

```makefile
.PHONY: test test-unit test-integration test-e2e test-performance test-all coverage bench clean-test

# Test variables
TEST_TIMEOUT ?= 10m
TEST_PARALLEL ?= 4

# Unit tests (fast, no build tags)
test-unit:
	@echo "Running unit tests..."
	go test -v -short -timeout=$(TEST_TIMEOUT) ./constraint/... ./rete/... ./cmd/...

# Integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v -tags=integration -timeout=$(TEST_TIMEOUT) ./tests/integration/...

# E2E tests
test-e2e:
	@echo "Running E2E tests..."
	go test -v -tags=e2e -timeout=$(TEST_TIMEOUT) ./tests/e2e/...

# E2E tests by category
test-e2e-alpha:
	@echo "Running alpha fixture tests..."
	go test -v -tags=e2e -run=TestAlphaFixtures ./tests/e2e/...

test-e2e-beta:
	@echo "Running beta fixture tests..."
	go test -v -tags=e2e -run=TestBetaFixtures ./tests/e2e/...

test-e2e-integration:
	@echo "Running integration fixture tests..."
	go test -v -tags=e2e -run=TestIntegrationFixtures ./tests/e2e/...

# Performance tests
test-performance:
	@echo "Running performance tests..."
	go test -v -tags=performance -timeout=1h ./tests/performance/...

# Load tests with profiling
test-load:
	@echo "Running load tests with profiling..."
	go test -v -tags=performance -run=TestLoad -cpuprofile=cpu.prof -memprofile=mem.prof ./tests/performance/...

# All tests
test-all: test-unit test-integration test-e2e
	@echo "All tests completed!"

# All tests with race detection
test-race:
	@echo "Running tests with race detector..."
	go test -race -tags=e2e,integration ./...

# Parallel execution
test-parallel:
	@echo "Running tests in parallel..."
	go test -v -tags=e2e,integration -parallel=$(TEST_PARALLEL) ./tests/...

# Coverage
coverage:
	@echo "Generating coverage report..."
	go test -tags=e2e,integration -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Coverage by package
coverage-unit:
	go test -short -coverprofile=coverage-unit.out ./constraint/... ./rete/...
	go tool cover -html=coverage-unit.out -o coverage-unit.html

coverage-e2e:
	go test -tags=e2e -coverprofile=coverage-e2e.out ./tests/e2e/...
	go tool cover -html=coverage-e2e.out -o coverage-e2e.html

# Benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem -run=^$$ ./...

bench-performance:
	@echo "Running performance benchmarks..."
	go test -tags=performance -bench=. -benchmem -run=^$$ ./tests/performance/...

# Verbose benchmarks with profiling
bench-profile:
	go test -bench=. -benchmem -cpuprofile=bench-cpu.prof -memprofile=bench-mem.prof ./...

# View profiles
profile-cpu:
	go tool pprof -http=:8080 cpu.prof

profile-mem:
	go tool pprof -http=:8080 mem.prof

# Clean test artifacts
clean-test:
	rm -f coverage*.out coverage*.html
	rm -f *.prof
	rm -f *.test

# Test with verbose output
test-verbose:
	go test -v -tags=e2e,integration ./...

# Quick smoke test
test-smoke:
	@echo "Running smoke tests..."
	go test -short -run=TestAlphaFixtures ./tests/e2e/... 2>&1 | head -20

# Help
help-test:
	@echo "Test targets:"
	@echo "  test-unit          - Run unit tests (fast)"
	@echo "  test-integration   - Run integration tests"
	@echo "  test-e2e          - Run E2E tests"
	@echo "  test-performance  - Run performance tests"
	@echo "  test-all          - Run all tests"
	@echo "  test-race         - Run with race detector"
	@echo "  coverage          - Generate coverage report"
	@echo "  bench             - Run benchmarks"
	@echo ""
	@echo "Environment variables:"
	@echo "  TEST_TIMEOUT      - Test timeout (default: 10m)"
	@echo "  TEST_PARALLEL     - Parallel test count (default: 4)"
```

**Acceptance Criteria:**
- All test types accessible via make
- Clear target names
- Configurable timeouts and parallelism
- Help documentation

### Phase 8: Remove Universal Runner

#### Task 8.1: Delete Universal Runner
**Action:** Remove files and references

Delete:
- `cmd/universal-rete-runner/main.go`
- `cmd/universal-rete-runner/main_test.go`
- `cmd/universal-rete-runner/` directory

Update:
- Remove build targets from Makefile
- Remove from `.gitignore` if present
- Remove from documentation

**Acceptance Criteria:**
- All runner files removed
- No broken references
- Build still works
- No dead code

#### Task 8.2: Clean Up Old Test Directories
**Action:** Remove old test locations

Delete (after migration confirmed):
- `test/coverage/alpha/` (empty after move)
- `test/coverage/beta/` (empty after move)
- `beta_coverage_tests/` (empty after move)
- `constraint/test/integration/` (empty after move)

Keep:
- `test/testutil/` (may have utilities)
- `test/` directory structure for documentation

**Acceptance Criteria:**
- Old fixture locations removed
- No duplicate files
- Git history preserved

### Phase 9: Validation and CI

#### Task 9.1: Validate Test Migration
**Action:** Comprehensive validation

Run:
```bash
# Validate all fixtures found
go test -tags=e2e -run=TestAllFixtures ./tests/e2e/... -v

# Validate count
# Should report: 83 fixtures (26+26+31)

# Run all test types
make test-all

# Check coverage
make coverage
```

Validate:
- All 83 fixtures execute successfully
- Test count matches original runner output
- No regressions in pass/fail status
- Coverage metrics reasonable (>70%)

**Acceptance Criteria:**
- All fixtures migrated and passing
- Test results match pre-migration
- No test gaps
- Coverage maintained or improved

#### Task 9.2: Update CI Configuration
**File:** `.github/workflows/test.yml` (if exists)

Update CI to use new test structure:

```yaml
name: Tests

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run unit tests
        run: make test-unit

  integration-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run integration tests
        run: make test-integration

  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run E2E tests
        run: make test-e2e

  coverage:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Generate coverage
        run: make coverage
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

**Acceptance Criteria:**
- CI runs all test types
- Tests run in parallel where possible
- Coverage tracked
- Fast feedback (< 10 minutes)

## Success Criteria

### Functional
- [ ] All 83 .tsd fixtures migrated to `tests/fixtures/`
- [ ] All fixtures pass via `go test -tags=e2e`
- [ ] Test utilities implement complete TSD execution
- [ ] Unit tests remain in modules
- [ ] E2E/integration tests centralized in `tests/`
- [ ] Universal runner removed completely
- [ ] All tests accessible via make commands

### Non-Functional
- [ ] Test execution time < 5 minutes for all tests
- [ ] Parallel execution works correctly
- [ ] No race conditions (`go test -race` passes)
- [ ] Coverage > 70% overall
- [ ] Clear test organization and documentation

### Documentation
- [ ] `tests/README.md` complete and accurate
- [ ] Main README updated with test info
- [ ] Makefile documented
- [ ] Build tags explained
- [ ] Contributing guide updated

### Quality
- [ ] All tests pass on first run
- [ ] No flaky tests
- [ ] Clear error messages
- [ ] Proper cleanup (no temp file leaks)
- [ ] CI integration working

## Migration Checklist

1. ✅ Create new directory structure
2. ✅ Implement test utilities (runner, fixtures, assertions)
3. ✅ Move alpha fixtures (26 files)
4. ✅ Move beta fixtures (26 files)
5. ✅ Move integration fixtures (31 files)
6. ✅ Create E2E tests (table-driven)
7. ✅ Create integration tests
8. ✅ Create performance tests
9. ✅ Update Makefile
10. ✅ Update documentation
11. ✅ Delete universal runner
12. ✅ Clean up old directories
13. ✅ Validate all tests pass
14. ✅ Update CI configuration

## Testing Strategy

### Test Execution Order
1. Unit tests (fastest, always run)
2. Integration tests (medium speed)
3. E2E tests (comprehensive)
4. Performance tests (slowest, optional)

### Test Isolation
- Each test gets fresh pipeline and storage
- Temp files cleaned up with `t.Cleanup()`
- No shared state between tests
- Parallel execution safe

### Error Handling
- Clear error messages with context
- Failed tests show fixture path
- Captured output included in failures
- Assertions use `t.Helper()` for correct line numbers

## Rollout Plan

1. **Implement utilities** (Phase 1)
2. **Migrate fixtures** (Phase 2) - keep runner as backup
3. **Create new tests** (Phases 3-5)
4. **Validate** (Phase 9.1) - compare runner vs go test results
5. **Remove runner** (Phase 8) - only after validation
6. **Update CI** (Phase 9.2)
7. **Document** (Phase 6)

## Risk Mitigation

- **Risk**: Test migration breaks existing functionality
  - **Mitigation**: Keep runner until full validation complete
  
- **Risk**: Fixtures lost during move
  - **Mitigation**: Use `mv` not `cp`, verify counts
  
- **Risk**: Performance regression
  - **Mitigation**: Benchmark before/after, parallelize tests
  
- **Risk**: Flaky tests in CI
  - **Mitigation**: Proper isolation, cleanup, timeouts

## Out of Scope

- Rewriting existing unit tests in modules
- Changing test assertions library (keep testify)
- Performance optimization of fixtures themselves
- Adding new test fixtures (can be done after migration)

## References

- Current runner: `cmd/universal-rete-runner/main.go`
- Fixture locations: `test/coverage/`, `beta_coverage_tests/`, `constraint/test/integration/`
- Go testing docs: https://pkg.go.dev/testing
- Table-driven tests: https://go.dev/wiki/TableDrivenTests

## Notes

- Follow Go testing conventions and idioms
- Use `t.Parallel()` where safe for performance
- Maintain backward compatibility during migration
- Keep test output clear and actionable
- Document any deviations from this plan