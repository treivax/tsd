# TSD Test Suite

This directory contains the complete test suite for the TSD (Type-Safe Datalog) project, organized using Go's standard testing framework with build tags for test categorization.

## 📁 Structure

```
tests/
├── e2e/                    # End-to-end tests (TSD fixtures)
├── integration/            # Module integration tests
├── performance/            # Performance and load tests
├── fixtures/               # Test data files
│   ├── alpha/             # Alpha node test fixtures (26 files)
│   ├── beta/              # Beta node test fixtures (26 files)
│   └── integration/       # Integration test fixtures (31 files)
└── shared/
    └── testutil/          # Shared test utilities
```

## 🚀 Running Tests

### Quick Start

```bash
# Run all unit tests (fast)
make test-unit

# Run all E2E tests
make test-e2e

# Run all tests (unit + integration + E2E)
make test-all

# Run with coverage
make coverage
```

### Unit Tests

Unit tests are co-located with their respective modules (constraint, rete, etc.) following Go conventions.

```bash
# Run all unit tests
go test -short ./constraint/... ./rete/... ./cmd/...

# Or via Makefile
make test-unit
```

**Characteristics:**
- Fast execution (< 1 second per package typically)
- No external dependencies
- Test individual functions/methods
- Run by default with `go test`

### E2E Tests

End-to-end tests execute complete TSD fixtures through the pipeline.

```bash
# Run all E2E tests
go test -tags=e2e ./tests/e2e/...

# Run specific fixture categories
make test-e2e-alpha        # Alpha node fixtures
make test-e2e-beta         # Beta node fixtures
make test-e2e-integration  # Integration fixtures

# Or via Makefile
make test-e2e
```

**Test Files:**
- `tests/e2e/tsd_fixtures_test.go` - Table-driven tests for all 83 fixtures

**Fixtures Coverage:**
- 26 Alpha node test fixtures (arithmetic, comparison, logic operations)
- 26 Beta node test fixtures (joins, multi-pattern matching)
- 31 Integration fixtures (complex scenarios, full pipeline)

### Integration Tests

Integration tests verify module interactions and pipeline flows.

```bash
# Run integration tests
go test -tags=integration ./tests/integration/...

# Or via Makefile
make test-integration
```

**Test Files:**
- `tests/integration/constraint_rete_test.go` - Constraint-RETE integration
- `tests/integration/pipeline_test.go` - Pipeline flows

**Coverage:**
- Cross-module interactions
- Type system validation
- Constraint evaluation with RETE execution
- Multi-rule and multi-type scenarios

### Performance Tests

Performance tests include load tests and benchmarks (slower, run separately).

```bash
# Run performance tests
go test -tags=performance ./tests/performance/...

# Run specific load tests
go test -tags=performance -run=TestLoad ./tests/performance/...

# Run benchmarks
go test -tags=performance -bench=. ./tests/performance/...

# Or via Makefile
make test-performance     # All performance tests
make test-load           # Load tests with profiling
make bench-performance   # Benchmarks only
```

**Test Files:**
- `tests/performance/load_test.go` - Load tests (100-10000 facts)
- `tests/performance/benchmark_test.go` - Go benchmarks

**Coverage:**
- Scalability testing (100, 1000, 5000, 10000 facts)
- Complex constraint evaluation performance
- Join operation performance
- Memory stress tests

## 🏷️ Build Tags

Tests are organized using Go build tags:

| Tag           | Purpose                              | Speed    |
|---------------|--------------------------------------|----------|
| (none)        | Unit tests                           | Fast     |
| `e2e`         | End-to-end TSD fixture tests         | Medium   |
| `integration` | Module integration tests             | Medium   |
| `performance` | Performance and load tests           | Slow     |

### Running Multiple Tags

```bash
# Run E2E and integration tests together
go test -tags=e2e,integration ./tests/...

# Run all tests including performance
go test -tags=e2e,integration,performance ./tests/...
```

## 🔧 Test Utilities

Shared test utilities are available in `tests/shared/testutil/`:

### Runner Utilities (`runner.go`)

Execute TSD files and capture results:

```go
// Simple execution
result := testutil.ExecuteTSDFile(t, "path/to/file.tsd")
testutil.AssertNoError(t, result)

// With options
result := testutil.ExecuteTSDFileWithOptions(t, "path/to/file.tsd", &testutil.ExecutionOptions{
    ExpectError:     false,
    MinActivations:  1,
    ValidateNetwork: true,
    Timeout:         5 * time.Second,
})
```

### Fixture Discovery (`fixtures.go`)

Discover and load test fixtures:

```go
// Discover all fixtures
fixtures := testutil.DiscoverFixtures()

// Get by category
alphaFixtures := testutil.GetFixturesByCategory("alpha")
betaFixtures := testutil.GetFixturesByCategory("beta")

// Load specific fixture
fixture := testutil.LoadFixture(t, "fixtures/alpha/alpha_abs_positive.tsd")
```

### Assertions (`assertions.go`)

TSD-specific assertions:

```go
// Assert network structure
testutil.AssertNetworkStructure(t, result, expectedTypeNodes, expectedTerminalNodes)

// Assert activations
testutil.AssertActivations(t, result, expectedCount)
testutil.AssertMinActivations(t, result, minCount)

// Assert errors
testutil.AssertNoError(t, result)
testutil.AssertError(t, result)
testutil.AssertErrorContains(t, result, "expected message")

// Assert output
testutil.AssertOutputContains(t, result, "expected text")

// Assert facts
testutil.AssertFactCount(t, result, expectedCount)
```

### Helpers (`helpers.go`)

Utility functions:

```go
// Create temporary TSD file
tempFile := testutil.CreateTempTSDFile(t, tsdContent)
defer testutil.CleanupTempFiles(t, tempFile)

// Get test data path
dataPath := testutil.GetTestDataPath()

// Skip if short mode
testutil.SkipIfShort(t, "long-running test")

// Timeouts
err := testutil.WithTimeout(5*time.Second, func() error {
    // ... test logic
})
```

## ⚡ Parallel Execution

All tests are fully thread-safe and support parallel execution:

```bash
# Run tests in parallel (default: 4 workers)
make test-parallel

# Custom parallelism
TEST_PARALLEL=8 make test-parallel

# Or directly with go test
go test -parallel=8 -tags=e2e,integration ./tests/...
```

**Note:** Tests that use `t.Parallel()` will run concurrently within a package.

**Thread-Safety:** As of 2025-12-04, all race conditions have been resolved. The test utilities use a global mutex to protect `os.Stdout` during output capture, ensuring safe parallel execution. See [`PARALLEL_TEST_FIX.md`](../PARALLEL_TEST_FIX.md) for technical details.

## 📊 Coverage

Generate coverage reports:

```bash
# Complete coverage
make coverage
# → Generates: coverage.html

# Unit tests only
make coverage-unit
# → Generates: coverage-unit.html

# E2E tests only
make coverage-e2e
# → Generates: coverage-e2e.html
```

View coverage in browser:

```bash
go tool cover -html=coverage.out
```

## 🎯 Benchmarking

Run benchmarks for performance analysis:

```bash
# All benchmarks
make bench

# Performance benchmarks only
make bench-performance

# With profiling
make bench-profile

# View CPU profile
make profile-cpu

# View memory profile
make profile-mem
```

Example benchmark output:

```
BenchmarkTSDExecution_Simple-8           1000    1234567 ns/op    123456 B/op    1234 allocs/op
BenchmarkFactProcessing_100Facts-8        100   12345678 ns/op   1234567 B/op   12345 allocs/op
```

## 🔍 Debugging Tests

### Verbose Output

```bash
# Verbose test output
make test-verbose

# Or directly
go test -v -tags=e2e ./tests/e2e/...
```

### Run Specific Tests

```bash
# Run specific test function
go test -tags=e2e -run=TestAlphaFixtures ./tests/e2e/...

# Run specific sub-test
go test -tags=e2e -run=TestAlphaFixtures/alpha_abs_positive ./tests/e2e/...

# Run with pattern matching
go test -tags=integration -run=TestPipeline ./tests/integration/...
```

### Race Detection

```bash
# Run with race detector
make test-race

# Or directly
go test -race -tags=e2e,integration ./...
```

## 🧪 Writing New Tests

### Unit Test Example

```go
package mypackage

import "testing"

func TestMyFunction(t *testing.T) {
    t.Parallel() // Enable parallel execution
    
    result := MyFunction("input")
    
    if result != "expected" {
        t.Errorf("Expected 'expected', got '%s'", result)
    }
}
```

### E2E Test Example

```go
//go:build e2e

package e2e

import (
    "testing"
    "github.com/treivax/tsd/tests/shared/testutil"
)

func TestNewFixture(t *testing.T) {
    t.Parallel()
    
    result := testutil.ExecuteTSDFile(t, "fixtures/new/test.tsd")
    
    testutil.AssertNoError(t, result)
    testutil.AssertNetworkStructure(t, result, 1, 1)
    testutil.AssertMinActivations(t, result, 1)
}
```

### Integration Test Example

```go
//go:build integration

package integration

import (
    "testing"
    "github.com/treivax/tsd/rete"
    "github.com/treivax/tsd/tests/shared/testutil"
)

func TestNewIntegration(t *testing.T) {
    t.Parallel()
    
    pipeline := rete.NewConstraintPipeline()
    storage := rete.NewMemoryStorage()
    
    // ... test logic
}
```

### Performance Test Example

```go
//go:build performance

package performance

import (
    "testing"
    "github.com/treivax/tsd/tests/shared/testutil"
)

func TestLoad_NewScenario(t *testing.T) {
    testutil.SkipIfShort(t, "performance test")
    
    // ... load test logic
}

func BenchmarkNewOperation(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // ... benchmark logic
    }
}
```

## 📝 Test Fixtures

### Adding New Fixtures

1. Place `.tsd` file in appropriate directory:
   - `fixtures/alpha/` - Alpha node tests
   - `fixtures/beta/` - Beta node tests
   - `fixtures/integration/` - Integration tests

2. Fixtures are automatically discovered by `testutil.DiscoverFixtures()`

3. Add to table-driven test in `e2e/tsd_fixtures_test.go` if needed

### Fixture Format

```tsd
# Type definitions
type MyType(field1: string, field2: number)

# Rules
rule my_rule : {m: MyType} / m.field2 > 0 ==> print("matched")

# Facts
MyType(field1:"test", field2:42)
```

## 🎯 CI/CD Integration

Tests are organized for easy CI/CD integration:

```yaml
# GitHub Actions example
- name: Unit Tests
  run: make test-unit

- name: Integration Tests
  run: make test-integration

- name: E2E Tests
  run: make test-e2e

- name: Coverage
  run: make coverage
```

## 📊 Test Metrics

Current test coverage:

- **Total Fixtures:** 83 TSD files
  - Alpha: 26 fixtures
  - Beta: 26 fixtures
  - Integration: 31 fixtures
  
- **Test Files:** 
  - Unit tests: ~125 `*_test.go` files in modules
  - E2E tests: 1 file with 83+ test cases
  - Integration tests: 2 files with 20+ test cases
  - Performance tests: 2 files with 15+ tests/benchmarks

## 🚨 Troubleshooting

### Tests Timing Out

```bash
# Increase timeout
TEST_TIMEOUT=30m make test-all

# Or directly
go test -timeout=30m -tags=e2e ./tests/e2e/...
```

### Parallel Test Issues

> **✅ RÉSOLU (2025-12-04):** Les problèmes de parallélisation ont été corrigés.  
> Les tests peuvent maintenant s'exécuter en parallèle de manière fiable.  
> Voir [`PARALLEL_TEST_FIX.md`](../PARALLEL_TEST_FIX.md) pour les détails.

```bash
# Les tests s'exécutent en parallèle par défaut maintenant
make test-integration

# Pour forcer un parallélisme spécifique si nécessaire
go test -parallel=8 -tags=integration ./tests/integration/...
```

### Fixture Not Found

Ensure you're using the correct path:

```go
// Correct - relative to test data path
testutil.ExecuteTSDFile(t, filepath.Join(testutil.GetTestDataPath(), "fixtures/alpha/test.tsd"))

// Or using the full project path
testutil.ExecuteTSDFile(t, "tests/fixtures/alpha/test.tsd")
```

## 📚 Additional Resources

- [Go Testing Package](https://golang.org/pkg/testing/)
- [Table-Driven Tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Go Build Constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints)

## 🔗 Migration Notes

This test structure replaces the previous `universal-rete-runner` approach:

- **Before:** Custom runner binary discovering and executing `.tsd` files
- **After:** Standard `go test` with organized test suites

**Benefits:**
- Standard Go tooling (coverage, profiling, race detection)
- Better IDE integration
- Full parallel execution support (thread-safe)
- Clear test categorization
- CI/CD friendly
- ~4.4x faster with parallel execution

The old runner is deprecated and will be removed in a future version.

---

For questions or contributions, please refer to the main project README.