# RETE Test Suite - Quick Start Guide

## 🚀 Quick Start

### Run All Tests
```bash
cd /home/resinsec/dev/tsd
go test ./...
```

### Run RETE Tests Only
```bash
go test ./rete -v
```

### Run Specific Test Categories
```bash
# Join cascade tests
go test ./rete -run "TestJoinNodeCascade" -v

# Partial evaluator tests
go test ./rete -run "TestPartialEval" -v

# Integration tests
go test ./test/integration -v
```

---

## 📋 Test Overview

### New Tests Added (Post-Refactor)

#### 1. Join Cascade Tests (`node_join_cascade_test.go`)
Tests the multi-variable join cascade architecture:
- ✅ 2-variable joins (User ⋈ Order)
- ✅ 3-variable joins (User ⋈ Order ⋈ Product)
- ✅ Order independence (facts can be submitted in any order)
- ✅ Multiple matching facts (cartesian products)
- ✅ Fact retraction through cascades

#### 2. Partial Evaluator Tests (`evaluator_partial_eval_test.go`)
Tests the partial evaluation mode for intermediate join stages:
- ✅ Unbound variable tolerance
- ✅ Logical operators (AND/OR)
- ✅ All comparison operators (==, !=, <, >, <=, >=)
- ✅ String comparisons
- ✅ Arithmetic expressions
- ✅ Edge cases and error handling

---

## 🎯 What's Being Tested

### Core Functionality
1. **Multi-variable joins** - Ensuring N-variable rules work via binary join cascades
2. **Incremental propagation** - Facts propagate through cascade levels correctly
3. **Partial evaluation** - Conditions evaluate even when not all variables bound
4. **Fact retraction** - Removing facts cleans up dependent tokens
5. **Join filtering** - Only matching facts create terminal tokens

### Integration
- End-to-end constraint pipeline
- Parser → Validator → Network builder → Fact submission
- Action extraction and tuple-space population

---

## 📊 Test Statistics

| Metric | Value |
|--------|-------|
| Total Test Files | 3+ (new) |
| Total Test Functions | 14+ (new) |
| Test Execution Time | ~350ms |
| Pass Rate | 100% ✅ |

---

## 🐛 Debugging Failed Tests

### Common Issues

#### "Expected N terminal tokens, got M"
**Cause**: Join condition not evaluating correctly.

**Fix**:
1. Check constraint file syntax
2. Verify variable bindings match fact types
3. Enable debug logging in JoinNode

#### "Variable non liée" (unbound variable)
**Cause**: Partial eval mode not enabled or cascade architecture issue.

**Fix**:
1. Verify `SetPartialEvalMode(true)` is called
2. Check join cascade construction
3. Inspect left/right memory contents

#### "Tokens not removed after retraction"
**Cause**: Retraction not propagating correctly.

**Fix**:
1. Verify `GetInternalID()` format
2. Check propagation to all child nodes
3. Inspect memory cleanup in all three memories (left/right/result)

---

## 📁 Test Files Location

```
tsd/
├── rete/
│   ├── node_join_cascade_test.go      # Join cascade integration tests
│   ├── evaluator_partial_eval_test.go # Partial evaluator unit tests
│   ├── rete_test.go                    # Core RETE tests
│   ├── docs/
│   │   └── TESTING.md                  # Comprehensive test documentation
│   └── test/
│       ├── incremental_propagation.constraint
│       └── incremental_propagation.facts
├── test/
│   └── integration/                    # Integration test suite
└── TESTING_IMPROVEMENTS_SUMMARY.md     # Complete summary of testing work
```

---

## 🔍 Test Details

### TestJoinNodeCascade_TwoVariablesIntegration
Tests 2-variable join (User ⋈ Order):
- Submit User → 0 terminal tokens
- Submit matching Order → 1 terminal token created
- Submit non-matching Order → still 1 token (filtering works)

### TestJoinNodeCascade_ThreeVariablesIntegration
Tests 3-variable cascade (User ⋈ Order ⋈ Product):
- Submit User → 0 tokens
- Submit Order → 0 tokens (missing Product)
- Submit Product → 1 token (complete cascade)

### TestPartialEval_UnboundVariables
Tests partial evaluation with unbound variables:
- Bound variable conditions evaluate successfully
- Unbound variable conditions tolerated gracefully

---

## 🛠️ Advanced Usage

### Run with Race Detection
```bash
go test -race ./rete
```

### Generate Coverage Report
```bash
go test -cover ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run Specific Test
```bash
go test ./rete -run TestJoinNodeCascade_ThreeVariables -v
```

### Run with Short Mode (Skip Long Tests)
```bash
go test -short ./...
```

---

## ✅ Test Checklist

Before committing code:
- [ ] Run `go test ./...` - all tests pass
- [ ] Run `go test -race ./rete` - no race conditions
- [ ] Run integration tests - pipeline works end-to-end
- [ ] Check coverage if modifying core logic
- [ ] Update tests if adding new features
- [ ] Update documentation if changing behavior

---

## 📚 Further Reading

- **[docs/TESTING.md](docs/TESTING.md)** - Comprehensive test documentation
- **[TESTING_IMPROVEMENTS_SUMMARY.md](../TESTING_IMPROVEMENTS_SUMMARY.md)** - Complete work summary
- **[rete_test.go](rete_test.go)** - Core RETE functionality tests

---

## 🤝 Contributing

When adding new tests:
1. Use descriptive test names: `TestFeature_Scenario`
2. Add structured logging with emojis: `t.Log("🧪 TEST: Description")`
3. Use helper functions to reduce duplication
4. Test both positive and negative cases
5. Add clear success/failure messages
6. Update this README if adding new test categories

---

## 📞 Support

If tests fail unexpectedly:
1. Check recent code changes
2. Review test logs for specific errors
3. Consult [docs/TESTING.md](docs/TESTING.md) debugging section
4. Run individual failing test with `-v` flag
5. Enable debug logging in relevant modules

---

**Last Updated**: 2024  
**Status**: ✅ All Tests Passing  
**Maintainer**: Engineering Team