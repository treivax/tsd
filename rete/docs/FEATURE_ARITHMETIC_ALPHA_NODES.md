# Feature: Arithmetic Expressions Support in AlphaNodes

**Feature ID**: RETE-002  
**Status**: In Development  
**Priority**: High  
**Created**: 2025-01-XX  
**Related**: BUG_RETE001 (Alpha/Beta Separation)

## Overview

This feature extends `AlphaConditionEvaluator` to support complex arithmetic expressions in alpha conditions, enabling more conditions to be evaluated early in AlphaNodes rather than being deferred to JoinNodes.

## Problem Statement

Currently, arithmetic expressions like `c.qte * 23 - 10 > 0` remain in JoinNodes even when they only reference a single variable. This is because:

1. `ConditionSplitter.isSimpleAlphaCondition()` only accepts direct field comparisons
2. `AlphaConditionEvaluator` already supports arithmetic operations but the splitter doesn't recognize them

### Current Behavior

```
Rule: {c: Commande} / c.qte * 23 - 10 > 0 ==> ...

Network Structure:
  TypeNode(Commande) 
    → PassthroughNode 
      → JoinNode [evaluates: c.qte * 23 - 10 > 0]
```

### Desired Behavior

```
Rule: {c: Commande} / c.qte * 23 - 10 > 0 ==> ...

Network Structure:
  TypeNode(Commande)
    → AlphaNode [evaluates: c.qte * 23 - 10 > 0]
      → PassthroughNode
        → TerminalNode
```

## Benefits

1. **Early Filtering**: Arithmetic conditions evaluated before joins
2. **Reduced Join Evaluations**: Fewer facts reach JoinNodes
3. **Better Performance**: Less computational work in beta network
4. **Architectural Correctness**: Respects RETE alpha/beta distinction

### Performance Impact (Estimated)

For rules with single-variable arithmetic conditions:
- **Before**: All facts propagate to JoinNode, evaluated there
- **After**: Only facts passing arithmetic test propagate
- **Expected reduction**: 30-70% fewer beta evaluations (depends on selectivity)

## Technical Design

### 1. Extend `isSimpleAlphaCondition()`

**File**: `rete/condition_splitter.go`

**Current Logic**:
```go
// Only accepts: c.field > literal
if leftType != "fieldAccess" {
    return false // Rejects arithmetic
}
```

**New Logic**:
```go
// Accept both:
// - c.field > literal (simple)
// - c.field * 23 - 10 > literal (arithmetic)

// Check if all variables in the expression belong to single binding
vars := cs.extractVariables(condition)
return len(vars) == 1  // Single variable = alpha condition
```

### 2. Validation Strategy

Before marking a condition as "alpha extractable", validate that `AlphaConditionEvaluator` can handle it:

```go
func (cs *ConditionSplitter) canAlphaEvaluatorHandle(condition map[string]interface{}) bool {
    // Create a test evaluator
    testEval := NewAlphaConditionEvaluator()
    testEval.partialEvalMode = true
    
    // Try to evaluate with nil bindings
    // If it returns an error about unsupported operations -> false
    // If it returns error about missing bindings -> true (normal)
    _, err := testEval.evaluateExpression(condition)
    
    return err == nil || isBindingError(err)
}
```

### 3. Test Coverage

**New Test File**: `rete/arithmetic_alpha_extraction_test.go`

Test cases:
- Simple arithmetic: `c.qte * 23 > 100`
- Nested arithmetic: `(c.qte * 23 - 10) / 2 > 50`
- Mixed operations: `c.qte * c.price - 10 > 0` (alpha - same variable)
- Cross-variable: `c.qte * p.price > 100` (beta - different variables)
- Edge cases: division by zero handling, type mismatches

### 4. Integration Points

**Modified Files**:
1. `rete/condition_splitter.go`
   - Update `isSimpleAlphaCondition()` → `isAlphaExtractable()`
   - Add validation logic

2. `rete/condition_splitter_test.go`
   - Add arithmetic test cases
   - Verify single-variable arithmetic → alpha
   - Verify multi-variable arithmetic → beta

3. `rete/builder_join_rules.go`
   - No changes needed (uses splitter output)

**New Files**:
1. `rete/arithmetic_alpha_extraction_test.go`
   - E2E tests with actual network building
   - Verify AlphaNodes created for arithmetic conditions
   - Performance benchmarks

## Implementation Plan

### Phase 1: Core Functionality (Priority: High)

**Tasks**:
- [ ] Refactor `isSimpleAlphaCondition()` to `isAlphaExtractable()`
- [ ] Implement single-variable detection for arithmetic
- [ ] Add unit tests in `condition_splitter_test.go`
- [ ] Create `arithmetic_alpha_extraction_test.go` with E2E tests

**Validation**:
```bash
go test ./rete -run TestArithmeticAlphaExtraction -v
go test ./rete -run TestConditionSplitter.*Arithmetic -v
```

### Phase 2: Advanced Validation (Priority: Medium)

**Tasks**:
- [ ] Implement `canAlphaEvaluatorHandle()` validation
- [ ] Add error classification (binding vs. unsupported)
- [ ] Test edge cases (division by zero, type errors)

### Phase 3: Optimization (Priority: Low)

**Tasks**:
- [ ] Cache validation results in splitter
- [ ] Add metrics for alpha extraction rate
- [ ] Benchmark performance improvements

## Testing Strategy

### Unit Tests

```go
func TestConditionSplitter_ArithmeticExpressions(t *testing.T) {
    tests := []struct {
        name      string
        condition map[string]interface{}
        wantAlpha bool
        wantBeta  bool
    }{
        {
            name: "simple arithmetic - single variable",
            condition: map[string]interface{}{
                "type": "comparison",
                "left": map[string]interface{}{
                    "type": "binaryOp",
                    "left": map[string]interface{}{
                        "type": "fieldAccess",
                        "object": "c",
                        "field": "qte",
                    },
                    "operator": "*",
                    "right": map[string]interface{}{
                        "type": "number",
                        "value": 23,
                    },
                },
                "operator": ">",
                "right": map[string]interface{}{
                    "type": "number",
                    "value": 100,
                },
            },
            wantAlpha: true,
            wantBeta:  false,
        },
        {
            name: "arithmetic - two variables (join)",
            condition: map[string]interface{}{
                // c.qte * p.price > 100
            },
            wantAlpha: false,
            wantBeta:  true,
        },
    }
    
    // ... test implementation
}
```

### Integration Tests

```go
func TestArithmeticAlphaExtraction_E2E(t *testing.T) {
    tsdContent := `
    type Commande(id: string, qte: number)
    action log(msg: string)
    
    rule high_quantity : {c: Commande} / c.qte * 23 - 10 > 100 
        ==> log("High quantity order")
    `
    
    network := buildNetwork(t, tsdContent)
    
    // Verify network structure
    stats := network.GetNetworkStats()
    alphaNodes := stats["alpha_nodes"].(int)
    
    if alphaNodes < 1 {
        t.Error("Expected AlphaNode for arithmetic condition")
    }
    
    // Verify filtering behavior
    // Submit facts: some pass, some fail the arithmetic condition
    // Only passing facts should reach terminal node
}
```

### Performance Tests

```go
func BenchmarkArithmeticInAlpha(b *testing.B) {
    // Compare: arithmetic in alpha vs. arithmetic in beta
    // Measure: evaluation count, memory allocations, time
}
```

## Migration & Compatibility

### Backward Compatibility

✅ **Fully Compatible**  
- No breaking changes to public API
- Existing networks continue to work
- Improvement is transparent to users

### Migration Path

No migration needed. Simply upgrade and rebuild rules:

```bash
# Rules are re-compiled with new optimization
go build
```

## Known Limitations & Future Work

### Phase 1 Limitations

1. **Cross-variable arithmetic** remains in beta (correct behavior)
   - Example: `c.qte * p.price` → JoinNode (requires both facts)

2. **Function calls in arithmetic** not yet optimized
   - Example: `len(c.items) * 2 > 10` → may stay in beta
   - TODO: Enhance function evaluation

3. **Complex nested expressions** may hit edge cases
   - Extensive testing needed

### Future Enhancements

1. **Partial Evaluation**: If `c.qte * p.price > 100`, extract `c.qte > 100/p.price` as alpha (symbolic)
2. **Constant Folding**: `c.qte * 23 - 10` → precompute `23, -10` coefficients
3. **Common Subexpression Elimination**: Share arithmetic nodes between rules

## Success Metrics

### Functional Metrics

- ✅ All single-variable arithmetic → AlphaNodes
- ✅ All multi-variable arithmetic → JoinNodes
- ✅ All existing tests pass
- ✅ New tests cover arithmetic cases

### Performance Metrics

Measure on `arithmetic_e2e.tsd`:
- **Alpha extraction rate**: % of arithmetic conditions extracted
- **Join evaluations**: Reduction vs. baseline
- **Network build time**: Should not significantly increase
- **Fact processing time**: Should decrease

### Quality Metrics

- ✅ No new compiler warnings
- ✅ 100% test coverage for new code
- ✅ Documentation updated
- ✅ Code review approved

## Documentation Updates

**Files to Update**:
1. `rete/docs/BUG_RETE001_ROOT_CAUSE_ANALYSIS.md`
   - Remove "Arithmetic Expressions" from limitations

2. `rete/docs/CHANGELOG_BUG_RETE001.md`
   - Add Phase 2: Arithmetic support completed

3. `README.md` (if exists in rete/)
   - Add note about arithmetic optimization

4. This document → `rete/docs/FEATURE_ARITHMETIC_ALPHA_NODES.md`

## References

- **Related**: BUG_RETE001 (Alpha/Beta Separation)
- **Code**: `rete/condition_splitter.go`, `rete/evaluator_operators.go`
- **Tests**: `rete/arithmetic_edge_cases_test.go`, `rete/evaluator_complex_expressions_test.go`
- **Theory**: Forgy, C. (1982). "Rete: A Fast Algorithm for the Many Pattern/Many Object Pattern Match Problem"

## Approval & Review

**Reviewers**:
- [ ] Architecture Review
- [ ] Performance Review
- [ ] Security Review (N/A - no security impact)

**Sign-off**:
- [ ] Lead Engineer
- [ ] QA Team

---

**Next Steps**:
1. Review this design document
2. Implement Phase 1 (Core Functionality)
3. Run full test suite
4. Measure performance improvements
5. Update documentation
6. Close related GitHub issues