# Contributing to TSD

Thank you for your interest in contributing to TSD! This guide will help you get started with development, testing, and submitting contributions.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Coding Standards](#coding-standards)
- [Testing Requirements](#testing-requirements)
- [Making Changes](#making-changes)
- [Submitting Pull Requests](#submitting-pull-requests)
- [Review Process](#review-process)
- [Performance Guidelines](#performance-guidelines)
- [Documentation](#documentation)
- [Getting Help](#getting-help)

---

## Code of Conduct

We are committed to providing a welcoming and inclusive environment. Please:

- Be respectful and constructive in all interactions
- Welcome newcomers and help them get started
- Focus on what is best for the community and project
- Show empathy towards other community members

Report unacceptable behavior to the project maintainers.

---

## Getting Started

### Prerequisites

Before contributing, ensure you have:

- **Go 1.21+** installed ([download](https://golang.org/dl/))
- **Git** for version control
- **Make** for build automation
- A **GitHub account** for pull requests

### Quick Start

```bash
# Clone the repository
git clone https://github.com/treivax/tsd.git
cd tsd

# Install dependencies
make deps

# Install development tools
make deps-dev

# Build the project
make build

# Run tests
make test-all

# Validate your setup
make validate
```

If all commands succeed, you're ready to contribute!

---

## Development Setup

### 1. Fork and Clone

```bash
# Fork the repository on GitHub, then clone your fork
git clone https://github.com/YOUR_USERNAME/tsd.git
cd tsd

# Add upstream remote
git remote add upstream https://github.com/treivax/tsd.git

# Verify remotes
git remote -v
```

### 2. Install Development Tools

```bash
# Install all development tools
make deps-dev
```

This installs:
- `goimports` - Import formatting
- `golangci-lint` - Static analysis
- Other Go tools as needed

### 3. Configure Your Editor

We use `.editorconfig` for consistent formatting. Install the EditorConfig plugin for your editor:

- **VSCode**: EditorConfig for VS Code
- **IntelliJ/GoLand**: Built-in support
- **Vim**: editorconfig-vim
- **Emacs**: editorconfig-emacs

### 4. Verify Your Setup

```bash
# Run quick validation
make quick-check

# Run full validation
make validate
```

---

## Project Structure

```
tsd/
├── cmd/                    # Command-line applications
│   ├── tsd/               # Main TSD binary (unified CLI)
│   └── */                 # Other utilities
├── rete/                  # RETE algorithm implementation
│   ├── network.go         # Core RETE network
│   ├── node_*.go          # Node implementations
│   ├── evaluator_*.go     # Expression evaluators
│   └── action_*.go        # Action execution
├── constraint/            # Language parser and AST
│   ├── parser.go          # PEG parser
│   └── ast.go             # Abstract syntax tree
├── auth/                  # Authentication system
│   ├── manager.go         # Auth manager
│   └── middleware.go      # HTTP middleware
├── tsdio/                 # I/O utilities
├── internal/              # Private packages
├── tests/                 # Test suites
│   ├── e2e/              # End-to-end tests
│   ├── integration/      # Integration tests
│   └── performance/      # Performance tests
├── examples/              # Example TSD programs
├── docs/                  # Documentation
└── scripts/               # Build and utility scripts
```

### Key Packages

- **`rete`**: Core rule engine (RETE algorithm)
- **`constraint`**: Language parsing and AST
- **`auth`**: Authentication and authorization
- **`cmd/tsd`**: Main CLI application
- **`tests`**: Comprehensive test suites

---

## Coding Standards

### Go Style Guide

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html).

### Naming Conventions

**Packages**:
```go
// Use lowercase, single-word names
package rete
package constraint
```

**Functions/Methods**:
```go
// Exported: PascalCase
func EvaluateCondition(cond interface{}) bool

// Unexported: camelCase
func evaluateFieldAccess(field string) interface{}
```

**Variables**:
```go
// Descriptive names, avoid abbreviations
network := NewReteNetwork()
alphaNode := &AlphaNode{}

// Short names OK in tight scopes
for i, v := range items {
    process(v)
}
```

**Constants**:
```go
// Exported: PascalCase
const DefaultTimeout = 30 * time.Second

// Unexported: camelCase
const maxRetries = 10
```

### Code Formatting

**Always run before committing**:
```bash
make format
```

**Rules**:
- Use `gofmt` (tabs, not spaces)
- One blank line between functions
- Group imports: stdlib, external, internal
- Add comments to exported symbols

**Example**:
```go
package rete

import (
    "fmt"
    "sync"

    "github.com/pkg/errors"

    "github.com/treivax/tsd/constraint"
)

// AlphaNode tests conditions on a single fact.
// It is part of the alpha network in the RETE algorithm.
type AlphaNode struct {
    BaseNode
    Condition    interface{}
    VariableName string
}

// EvaluateCondition tests if the fact satisfies the condition.
func (a *AlphaNode) EvaluateCondition(fact *Fact) (bool, error) {
    if fact == nil {
        return false, errors.New("fact cannot be nil")
    }
    // Implementation...
    return true, nil
}
```

### Error Handling

**Always handle errors explicitly**:
```go
// Good
result, err := operation()
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Bad
result, _ := operation()  // Never ignore errors
```

**Error Messages**:
- Start with lowercase (except proper nouns)
- No trailing punctuation
- Provide context

```go
// Good
return fmt.Errorf("failed to evaluate condition for variable %s: %w", varName, err)

// Bad
return fmt.Errorf("Error: %s.", err.Error())
```

### Concurrency

**Use mutexes for shared state**:
```go
type AlphaNode struct {
    // ...
    mutex sync.RWMutex
}

func (a *AlphaNode) AddFact(fact *Fact) {
    a.mutex.Lock()
    defer a.mutex.Unlock()
    // Modify state
}

func (a *AlphaNode) GetFacts() []*Fact {
    a.mutex.RLock()
    defer a.mutex.RUnlock()
    // Read state
}
```

**Avoid goroutine leaks**:
```go
// Always provide cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go func() {
    select {
    case <-ctx.Done():
        return
    case result := <-resultChan:
        process(result)
    }
}()
```

---

## Testing Requirements

### Test Coverage

**Minimum requirements**:
- New code: **80%+ coverage**
- Bug fixes: Add regression test
- Refactoring: Maintain existing coverage

**Check coverage**:
```bash
make coverage
# Opens coverage.html in browser
```

### Test Types

#### 1. Unit Tests

Test individual functions/methods in isolation.

**Location**: Same package as code (`*_test.go`)

**Example**:
```go
func TestAlphaNode_EvaluateCondition(t *testing.T) {
    tests := []struct {
        name      string
        condition interface{}
        fact      *Fact
        want      bool
        wantErr   bool
    }{
        {
            name: "simple equality",
            condition: map[string]interface{}{
                "field": "status",
                "op":    "==",
                "value": "active",
            },
            fact: &Fact{
                Type: "Account",
                Fields: map[string]interface{}{
                    "status": "active",
                },
            },
            want:    true,
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            node := &AlphaNode{Condition: tt.condition}
            got, err := node.EvaluateCondition(tt.fact)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("EvaluateCondition() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("EvaluateCondition() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Run unit tests**:
```bash
make test-unit
```

#### 2. Integration Tests

Test interactions between components.

**Location**: `tests/integration/`

**Tags**: `// +build integration`

**Run integration tests**:
```bash
make test-integration
```

#### 3. E2E Tests

Test complete TSD programs (`.tsd` files).

**Location**: `tests/e2e/`

**Fixtures**: `tests/e2e/fixtures/`

**Run E2E tests**:
```bash
make test-e2e              # All E2E tests
make test-e2e-alpha        # Alpha network tests
make test-e2e-beta         # Beta network tests
make test-e2e-integration  # Full integration tests
```

#### 4. Performance Tests

Benchmark critical paths.

**Example**:
```go
func BenchmarkAlphaNode_Evaluate(b *testing.B) {
    node := &AlphaNode{/* setup */}
    fact := &Fact{/* setup */}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        node.EvaluateCondition(fact)
    }
}
```

**Run benchmarks**:
```bash
make bench              # All benchmarks
make bench-performance  # Performance-specific
make bench-profile      # With CPU/memory profiling
```

### Test Best Practices

**1. Table-Driven Tests**:
```go
tests := []struct {
    name string
    // inputs
    // expected outputs
}{
    {name: "case 1", /* ... */},
    {name: "case 2", /* ... */},
}
```

**2. Clear Test Names**:
```go
// Good
func TestAlphaNode_EvaluateCondition_WithNilFact_ReturnsError(t *testing.T)

// Bad
func TestAlpha(t *testing.T)
```

**3. Test Edge Cases**:
- Nil inputs
- Empty collections
- Boundary values
- Concurrent access

**4. Use Subtests**:
```go
t.Run("nil fact", func(t *testing.T) {
    // Test nil handling
})
```

**5. Clean Up Resources**:
```go
t.Cleanup(func() {
    // Cleanup code
})
```

---

## Making Changes

### Branching Strategy

**Branch naming**:
```bash
# Features
feature/add-aggregation-support
feature/improve-beta-sharing

# Bug fixes
fix/alpha-node-memory-leak
fix/casting-error-messages

# Documentation
docs/architecture-guide
docs/update-readme

# Refactoring
refactor/simplify-evaluator
refactor/extract-interface
```

### Workflow

**1. Create a branch**:
```bash
git checkout -b feature/my-feature
```

**2. Make your changes**:
```bash
# Edit files
vim rete/new_feature.go

# Format code
make format

# Run tests
make test-all

# Check for issues
make lint
```

**3. Commit your changes**:
```bash
git add .
git commit -m "feat: add support for complex aggregations

- Implement MultiSourceAccumulatorNode
- Add tests for multi-source patterns
- Update documentation

Fixes #123"
```

**Commit Message Format**:
```
<type>: <subject>

<body>

<footer>
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `refactor`: Code refactoring
- `test`: Test additions/changes
- `perf`: Performance improvements
- `chore`: Build/tooling changes

**4. Push and create PR**:
```bash
git push origin feature/my-feature
```

Then create a pull request on GitHub.

---

## Submitting Pull Requests

### Before Submitting

**Checklist**:
- [ ] Code follows style guidelines
- [ ] All tests pass (`make test-all`)
- [ ] Added tests for new code
- [ ] Updated documentation
- [ ] Ran `make validate`
- [ ] No linting errors (`make lint`)
- [ ] Rebased on latest `main`

**Validation**:
```bash
# Run complete validation
make validate

# Should output:
# ✅ Formatage
# ✅ Analyse statique
# ✅ Compilation
# ✅ Tests unitaires
# ✅ Tests d'intégration
# ✅ Tests E2E
```

### PR Guidelines

**Title**: Clear and descriptive
```
Good: "Add support for EXISTS patterns in beta network"
Bad:  "Update code"
```

**Description**: Include:
```markdown
## What
Brief description of changes

## Why
Motivation and context

## How
Implementation approach

## Testing
How to test the changes

## Checklist
- [ ] Tests added
- [ ] Docs updated
- [ ] Breaking changes documented
```

**Size**: Keep PRs focused and reasonably sized
- Small PRs: < 300 lines (preferred)
- Medium PRs: 300-1000 lines
- Large PRs: Discuss with maintainers first

---

## Review Process

### What Reviewers Look For

1. **Correctness**: Does it work as intended?
2. **Tests**: Adequate test coverage?
3. **Style**: Follows coding standards?
4. **Performance**: No obvious bottlenecks?
5. **Documentation**: Clear and complete?
6. **Breaking Changes**: Are they necessary and documented?

### Responding to Feedback

**Be responsive**:
- Address comments promptly
- Ask questions if unclear
- Push updates as new commits

**Example response**:
```markdown
> Can this be simplified?

Good point! I've extracted the logic into a helper function.
Committed in abc123.
```

### Approval and Merge

**Requirements**:
- At least 1 approval from maintainer
- All CI checks passing
- No unresolved comments
- Up-to-date with `main`

**Merge strategy**:
- Squash small feature PRs
- Merge commits for large features
- Maintainers will handle the merge

---

## Performance Guidelines

### Profiling

**CPU profiling**:
```bash
make test-load
make profile-cpu
```

**Memory profiling**:
```bash
make bench-profile
make profile-mem
```

### Performance Checklist

**Algorithm complexity**:
- [ ] Time complexity documented
- [ ] Space complexity acceptable
- [ ] Compared to alternatives

**Memory management**:
- [ ] No unnecessary allocations in hot paths
- [ ] Objects pooled where appropriate
- [ ] No memory leaks

**Concurrency**:
- [ ] Lock contention minimized
- [ ] No unnecessary goroutines
- [ ] Proper synchronization

### Benchmarking

**Always benchmark performance changes**:
```go
func BenchmarkBefore(b *testing.B) {
    // Baseline
}

func BenchmarkAfter(b *testing.B) {
    // Your optimization
}
```

**Compare results**:
```bash
go test -bench=. -benchmem ./... > old.txt
# Make changes
go test -bench=. -benchmem ./... > new.txt
benchcmp old.txt new.txt
```

---

## Documentation

### Code Documentation

**Document all exported symbols**:
```go
// AlphaNode tests conditions on a single fact.
// It is part of the alpha network in RETE algorithm.
//
// The node maintains working memory of facts that satisfy
// its condition and propagates matching facts to children.
type AlphaNode struct {
    // ...
}

// EvaluateCondition tests if the fact satisfies the node's condition.
// It returns true if the fact matches, false otherwise.
//
// Returns an error if the condition is malformed or evaluation fails.
func (a *AlphaNode) EvaluateCondition(fact *Fact) (bool, error) {
    // ...
}
```

### User Documentation

Update relevant documentation in `docs/`:

- **New features**: Add to `USER_GUIDE.md`
- **API changes**: Update `API_REFERENCE.md`
- **Breaking changes**: Update `CHANGELOG.md`
- **Examples**: Add to `examples/`

### Documentation Standards

**Format**:
- Use Markdown
- Clear headings
- Code examples with output
- Cross-reference related docs

**Example**:
```markdown
## Type Casting

TSD supports explicit type casting using the `cast` operator.

### Syntax

\```
cast(expression as type)
\```

### Supported Casts

| From   | To     | Example                    |
|--------|--------|----------------------------|
| number | string | `cast(123 as string)`      |
| string | number | `cast("123" as number)`    |
| bool   | string | `cast(true as string)`     |

### Examples

\```tsd
rule convert_amount
when
    order: Order(
        amount_str: string_amount,
        cast(string_amount as number) > 1000
    )
then
    Print("High value order: " + string_amount);
end
\```

See also: [Operators](USER_GUIDE.md#operators)
```

---

## Getting Help

### Documentation

- **[User Guide](USER_GUIDE.md)**: Complete feature reference
- **[Architecture](ARCHITECTURE.md)**: Technical design
- **[Tutorial](TUTORIAL.md)**: Step-by-step learning

### Communication

- **Issues**: [GitHub Issues](https://github.com/treivax/tsd/issues)
- **Discussions**: [GitHub Discussions](https://github.com/treivax/tsd/discussions)
- **Email**: maintainers@tsd.dev (for security issues)

### Common Questions

**Q: Where should I start?**
- Look for issues labeled `good first issue`
- Fix documentation typos/improvements
- Add test coverage
- Improve error messages

**Q: How do I run a specific test?**
```bash
go test -v -run TestAlphaNode ./rete/
```

**Q: My PR is failing CI. What do I do?**
```bash
# Run the same checks locally
make ci
```

**Q: How do I update my branch?**
```bash
git fetch upstream
git rebase upstream/main
```

---

## License

By contributing to TSD, you agree that your contributions will be licensed under the MIT License.

See [LICENSE](../LICENSE) for details.

---

## Thank You!

Every contribution, no matter how small, makes TSD better. We appreciate your time and effort!

**Contributors are recognized in**:
- [README.md](../README.md#contributors)
- [GitHub Contributors](https://github.com/treivax/tsd/graphs/contributors)

---

*Last updated: 2024-12-07*