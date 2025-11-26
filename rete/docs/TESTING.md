# RETE Evaluator Testing Documentation

## Overview

This document describes the comprehensive test suite added to ensure the correctness of the RETE evaluator refactor and multi-variable join cascade implementation.

## Test Suite Structure

### 1. Join Node Cascade Tests (`node_join_cascade_test.go`)

These integration tests verify the correctness of cascading join operations for multi-variable rules.

#### `TestJoinNodeCascade_TwoVariablesIntegration`
**Purpose**: Validates that 2-variable joins work correctly through the constraint pipeline.

**Test Scenario**:
- Creates a constraint with User and Order types
- Submits User fact first → expects 0 terminal tokens
- Submits matching Order fact → expects terminal token created
- Submits non-matching Order fact → expects filtering works correctly

**Key Validations**:
- Join condition evaluation (u.id == o.user_id)
- Incremental token propagation
- Condition-based filtering

#### `TestJoinNodeCascade_ThreeVariablesIntegration`
**Purpose**: Validates that 3-variable cascade joins work correctly (User ⋈ Order ⋈ Product).

**Test Scenario**:
- Uses the existing `test/incremental_propagation.constraint` file
- Submits facts incrementally: User → Order → Product
- Verifies terminal tokens only created when all three facts present
- Tests conditions: u.id == o.user_id AND o.product_id == p.id AND u.age >= 18

**Key Validations**:
- Multi-level join cascade propagation
- Intermediate join state management
- Full N-tuple assembly only when complete

#### `TestJoinNodeCascade_OrderIndependence`
**Purpose**: Ensures join results are independent of fact submission order.

**Test Scenarios**:
- Submits facts in different orders: User→Order, Order→User
- Verifies same terminal token count regardless of order

**Key Validations**:
- Left/Right memory management in JoinNodes
- Bidirectional join evaluation
- Memory consistency across submission patterns

#### `TestJoinNodeCascade_MultipleMatchingFacts`
**Purpose**: Validates cartesian product behavior with multiple facts per type.

**Test Scenario**:
- Submits 2 Users (U1, U2)
- Submits 3 Orders: 2 for U1, 1 for U2
- Expects 3 terminal tokens: (U1,O1), (U1,O2), (U2,O3)

**Key Validations**:
- Correct cartesian product generation
- Join condition filtering across multiple candidates
- Memory scaling with fact volume

#### `TestJoinNodeCascade_Retraction`
**Purpose**: Tests that fact retraction correctly removes dependent join tokens.

**Test Scenario**:
- Submits User and Order, creating a joined token
- Retracts User fact
- Verifies all dependent tokens removed

**Key Validations**:
- Retraction propagation through cascade
- Memory cleanup in left/right/result memories
- Terminal token removal

---

### 2. Partial Evaluation Tests (`evaluator_partial_eval_test.go`)

These tests validate the partial evaluation mode of `AlphaConditionEvaluator`, which allows condition evaluation when not all variables are bound (critical for cascade joins).

#### `TestPartialEval_UnboundVariables`
**Purpose**: Verifies that partial eval mode tolerates references to unbound variables.

**Test Scenario**:
- Binds only variable 'u'
- Evaluates condition u.age >= 18 (should succeed)
- Evaluates condition o.user_id == u.id where 'o' is unbound (should tolerate)

**Key Validations**:
- Evaluable conditions succeed even in partial mode
- Unevaluable conditions fail gracefully without errors

#### `TestPartialEval_LogicalExpressions`
**Purpose**: Tests logical operators (AND/OR) in partial eval mode.

**Test Scenario**:
- Tests: u.age >= 18 AND u.age <= 65
- Tests short-circuit behavior when conditions fail

**Key Validations**:
- AND operator evaluation
- Condition chaining
- Short-circuit optimization

#### `TestPartialEval_MixedBoundUnbound`
**Purpose**: Validates conditions mixing bound and unbound variables.

**Test Scenario**:
- Binds 'u', leaves 'o' and 'p' unbound
- Tests: u.age >= 18 AND o.user_id == u.id AND p.price > 100

**Key Validations**:
- Partial evaluation of evaluable sub-expressions
- Graceful handling of unevaluable sub-expressions

#### `TestPartialEval_ComparisonOperators`
**Purpose**: Comprehensive test of all comparison operators in partial mode.

**Test Cases**:
- Equality (==), Inequality (!=)
- Greater Than (>), Greater Or Equal (>=)
- Less Than (<), Less Or Equal (<=)
- Both true and false cases

**Key Validations**:
- All operators work correctly
- Type coercion for numeric comparisons

#### `TestPartialEval_StringComparisons`
**Purpose**: Tests string field comparisons.

**Test Scenario**:
- Tests: u.status == "active"
- Tests: u.status != "inactive"

**Key Validations**:
- String equality/inequality
- Case-sensitive comparison

#### `TestPartialEval_NestedFieldAccess`
**Purpose**: Validates multiple field accesses in complex conditions.

**Test Scenario**:
- Tests: u.age > 18 AND u.name == "Bob"

**Key Validations**:
- Multiple field dereferences
- Field access caching/optimization

#### `TestPartialEval_NormalModeComparison`
**Purpose**: Documents behavioral difference between normal and partial eval modes.

**Test Scenario**:
- Same condition evaluated in both modes
- Unbound variable 'o' referenced

**Key Validations**:
- Normal mode: may error on unbound variable
- Partial mode: tolerates unbound variable

#### `TestPartialEval_ArithmeticExpressions`
**Purpose**: Tests arithmetic operations in partial mode.

**Test Scenario**:
- Tests: o.quantity * o.price > 50

**Key Validations**:
- Arithmetic operator evaluation
- Numeric type handling

#### `TestPartialEval_EdgeCases`
**Purpose**: Tests error conditions and edge cases.

**Test Cases**:
- Missing field access
- Nil condition
- Empty condition map

**Key Validations**:
- Graceful error handling
- No crashes on invalid input

---

## Test Execution

### Run All RETE Tests
```bash
go test ./rete -v
```

### Run Join Cascade Tests Only
```bash
go test ./rete -run "TestJoinNodeCascade" -v
```

### Run Partial Eval Tests Only
```bash
go test ./rete -run "TestPartialEval" -v
```

### Run Integration Tests
```bash
go test ./test/integration -v
```

### Run All Tests
```bash
go test ./... -v
```

---

## Test Coverage Areas

### ✅ Covered
1. **Multi-variable joins (2 and 3 variables)**
   - Cascade architecture validation
   - Incremental propagation
   - Order independence

2. **Partial evaluation mode**
   - Unbound variable tolerance
   - Logical operators (AND/OR)
   - All comparison operators
   - Arithmetic expressions
   - String operations

3. **Fact retraction**
   - Cascading retraction through joins
   - Memory cleanup

4. **Cartesian products**
   - Multiple facts per type
   - Join filtering

5. **Integration with constraint pipeline**
   - End-to-end validation
   - Parser integration
   - Semantic validation

### 📝 Future Test Improvements

1. **Performance Tests**
   - Large fact volumes (1000+ facts)
   - Deep cascade chains (5+ variables)
   - Memory usage profiling

2. **Concurrent Operations**
   - Parallel fact submission
   - Race condition detection
   - Lock contention measurement

3. **Specialized Join Patterns**
   - Self-joins
   - Multi-way joins (N>3)
   - Negative joins with NOT

4. **Error Recovery**
   - Invalid fact types
   - Malformed conditions
   - Memory exhaustion scenarios

---

## Debugging Failed Tests

### Common Issues

#### Issue: "Expected N terminal tokens, got M"
**Diagnosis**: Join condition may not be evaluating correctly.

**Debug Steps**:
1. Check join condition syntax
2. Verify variable types match fact types
3. Enable debug logging in JoinNode
4. Inspect left/right memory contents

#### Issue: "Variable non liée" errors
**Diagnosis**: Partial eval mode not enabled or condition references truly missing variable.

**Debug Steps**:
1. Verify `SetPartialEvalMode(true)` is called
2. Check that all referenced variables are in scope
3. Verify cascade join architecture

#### Issue: "Terminal tokens not removed after retraction"
**Diagnosis**: Retraction not propagating through network.

**Debug Steps**:
1. Verify `GetInternalID()` returns correct format
2. Check retraction calls propagate to children
3. Inspect memory cleanup in all node types

---

## Test Maintenance

### Adding New Tests

1. **Create descriptive test function name**: `TestFeature_Scenario`
2. **Add structured logging**:
   ```go
   t.Log("🧪 TEST: Description")
   t.Log("====================")
   ```
3. **Use helper functions**: `createTempConstraintFile`, `countAllTerminalTokens`
4. **Verify positive and negative cases**
5. **Add clear success message**: `t.Log("\n🎊 TEST PASSED: Summary")`

### Updating Tests After Refactors

When modifying RETE core functionality:

1. Run full test suite before changes
2. Document expected behavior changes
3. Update affected tests
4. Verify integration tests still pass
5. Add new tests for new features
6. Update this documentation

---

## Related Documentation

- [Network Architecture](./NETWORK_ARCHITECTURE.md) - RETE network structure
- [Evaluator Design](./EVALUATOR_DESIGN.md) - Condition evaluation internals
- [Pipeline Flow](./PIPELINE_FLOW.md) - Constraint file processing

---

## Test Statistics

- **Total Test Files**: 3
- **Total Test Functions**: 20+
- **Coverage Focus**: Join operations, partial evaluation, integration
- **Execution Time**: ~350ms (all tests)
- **Success Rate**: 100% (as of last refactor completion)

---

## Continuous Integration

### Pre-commit Checks
```bash
# Run all tests
make test

# Run with race detection
go test -race ./...

# Check test coverage
go test -cover ./...
```

### CI Pipeline (Recommended)
1. Run `go test ./...` on every commit
2. Require 100% test pass rate for merge
3. Run integration tests on main branch
4. Generate coverage reports
5. Run performance benchmarks weekly

---

**Document Version**: 1.0  
**Last Updated**: 2024  
**Authors**: Engineering Team  
**Status**: Active