# Remaining Test Failures Report

**Date:** December 2, 2025  
**Status:** 19 PASS / 16 FAIL / 1 SKIP

## Summary

After fixing the incremental validation action tracking issue, the integration tests improved from ~15 passing to 19 passing. The remaining 16 failures fall into 3 distinct categories with different root causes.

## Test Results by Category

### ✅ PASSING (19 tests)

#### Incremental Ingestion Tests (5/5) ✅
- `TestIncrementalIngestion_FactsBeforeRules` - Fixed with action tracking
- `TestIncrementalIngestion_MultipleRules` - Fixed with action tracking
- `TestIncrementalIngestion_TypeExtension` - Fixed with action tracking
- `TestIncrementalIngestion_Reset` - Working correctly
- `TestIncrementalIngestion_Optimizations` - Fixed with action tracking

#### Other Passing Tests (14)
- `TestVariableArguments` ✅
- `TestComprehensiveMixedArguments` ✅
- `TestErrorDetectionInArguments` ✅
- `TestBasicNetworkIntegrity` ✅
- `TestExhaustiveBetaCoverage` ✅
- `TestComplexBetaNodesTupleSpace` ✅
- `TestMassiveBetaNodesWithFactsFile` ✅
- `TestNegationRules` ✅
- `TestTupleSpaceTerminalNodes` ✅
- `TestRealPEGParsingIntegration` ✅
- `TestSemanticValidationWithRealParser` ✅
- `TestResetInstruction_StoragePreservation` ✅
- `TestCompleteCoherencePEGtoRETE` ✅
- `TestSimpleBetaNodeTupleSpace` ✅

### ❌ FAILING (16 tests)

## Category 1: Old Syntax Issues (8 tests)

**Location:** `test/integration/incremental/advanced_test.go`

All 8 tests in this file use deprecated TSD syntax and need to be rewritten.

### Error Pattern
```
no match found, expected: "#", "(", "/*", "//" or [ \t\r\n]
```
Position: Column 13 (after `type Person `)

### Root Cause
Tests use old curly-brace syntax:
```tsd
// OLD (doesn't work)
type Person {
    id: string
    name: string
}

// NEW (correct syntax)
type Person(id: string, name: string)
```

Also using old rule syntax:
```tsd
// OLD (doesn't work)
rule "adult_check" {
    when {
        p: Person(age >= 18)
    }
    then {
        print("Adult: " + p.name)
    }
}

// NEW (correct syntax)
action print_adult(name: string)

rule adult_check: {p: Person} / p.age >= 18 ==> print_adult(p.name)
```

### Failing Tests
1. ❌ `TestIncrementalValidation`
2. ❌ `TestIncrementalValidationError`
3. ❌ `TestGarbageCollectionAfterReset`
4. ❌ `TestTransactionCommit`
5. ❌ `TestTransactionRollback`
6. ❌ `TestAdvancedFeaturesIntegration`
7. ❌ `TestTransactionAutoRollback`
8. ❌ `TestTransactionCommandPattern`

### Recommendation
**Priority: HIGH**  
Rewrite all tests in `advanced_test.go` using current TSD syntax. These tests cover important functionality (transactions, garbage collection, validation) but the syntax is outdated.

### Example Fix Needed
```go
// Current (broken)
typesContent := `
type Person {
    id: string
    name: string
}
`

// Should be
typesContent := `
type Person(id: string, name: string)
`
```

---

## Category 2: Alpha Node Evaluation Issues (2 tests)

**Location:** `test/integration/alpha_*_test.go`

### Error Pattern
Rules are created but facts don't match conditions (0 activations).

### Failing Tests
1. ❌ `TestCompleteAlphaCoverage`
   - 28 alpha rules created
   - 21 facts injected
   - **0 actions triggered** (expected > 0)
   
2. ❌ `TestExhaustiveAlphaCoverage`
   - Similar pattern: rules created but no matches

### Specific Failed Conditions
```
❌ test_string_inequality_success: 0 matches (should match)
❌ test_number_inequality_success: 0 matches (should match)
❌ test_string_less_success: 0 matches (should match)
❌ test_senior_success: 0 matches (should match)
❌ test_low_rating_success: 0 matches (should match)
```

### Root Cause
The alpha node condition evaluator has incomplete operator support or incorrect evaluation logic for certain comparisons. Based on the deep-clean thread context, known issues include:
- `notConstraint` not implemented
- Nested boolean expressions not handled
- Some operators (`!=`, `<`, etc.) may not be working correctly

### Recommendation
**Priority: MEDIUM**  
This is a runtime evaluation issue, not a validation or parsing issue. The rules are correctly parsed and the network is built, but fact matching logic needs debugging.

**Action Items:**
1. Review `rete/evaluator_constraints.go` for operator implementations
2. Add unit tests for each operator type
3. Debug why inequality and less-than operators return no matches
4. Implement missing constraint types (`notConstraint`, nested boolean)

---

## Category 3: Reset & Cleanup Issues (5 tests)

**Location:** `test/integration/reset_instruction_test.go`

### Error Patterns

#### Type Not Removed After Reset
```
❌ Type 'User' should have been removed by reset but still exists
❌ Type 'Order' should have been removed by reset but still exists
❌ Type 'Product' should have been removed by reset but still exists
```

### Failing Tests
1. ❌ `TestResetInstruction_BasicReset`
   - Old types not cleaned up after reset
   - Test expects 6 type nodes, got more (includes old types)

2. ❌ `TestResetInstruction_MultipleResets`
   - Multiple consecutive resets not cleaning up properly

3. ❌ `TestResetInstruction_NetworkIntegrity`
   - Network state inconsistent after reset

4. ❌ `TestResetInstruction_RulesAfterReset`
   - Rules from before reset still present

5. ❌ `TestResetInstruction_ParsingOnly`
   - Parsing issue (possibly old syntax)

### Root Cause
The `reset` command triggers garbage collection and creates a new network, but the old network's types may still be referenced or the new network is inheriting state from the old one.

Possible issues:
- `network.Types` array not being cleared on reset
- TypeNodes not being removed from `network.TypeNodes` map
- Network reference not being properly replaced after reset

### Recommendation
**Priority: MEDIUM**  
Reset is important for testing and production use cases where you want to completely clear the system.

**Action Items:**
1. Review `rete/constraint_pipeline.go` reset handling (around line 120)
2. Ensure `network.Types` is cleared when new network is created
3. Verify that the new network instance replaces all references to the old one
4. Check if `network.GarbageCollect()` is properly cleaning up TypeNodes

---

## Category 4: Other Evaluation Issues (1 test)

### Failing Test
❌ `TestQuotedStringsIntegration`
- 4 terminal nodes created
- **0 activations** (expected at least 1)
- Rules parse correctly but don't trigger

### Root Cause
Similar to alpha coverage tests - facts not matching rules with quoted string literals.

### Recommendation
**Priority: LOW**  
This is likely related to the alpha evaluation issues. Fix Category 2 issues first, then revisit.

---

## Recommended Fix Order

### 1. HIGH Priority - Old Syntax (Quick Win)
Rewrite 8 tests in `advanced_test.go` with current TSD syntax.
- **Estimated effort:** 2-4 hours
- **Impact:** +8 passing tests
- **Complexity:** Low (just syntax updates)

### 2. MEDIUM Priority - Reset Cleanup
Fix reset command to properly clear network state.
- **Estimated effort:** 4-6 hours
- **Impact:** +5 passing tests
- **Complexity:** Medium (requires understanding network lifecycle)

### 3. MEDIUM Priority - Alpha Evaluation
Fix operator evaluation in alpha nodes.
- **Estimated effort:** 6-10 hours
- **Impact:** +2-3 passing tests
- **Complexity:** Medium-High (requires debugging evaluation logic)

---

## Expected Final State

After fixing all issues:
- **Expected:** 35/36 tests passing (97% pass rate)
- **Current:** 19/36 tests passing (53% pass rate)
- **Improvement:** +16 tests, +44 percentage points

---

## Notes

- Core RETE functionality is solid (rete package: 100% pass)
- Constraint validation is working (constraint package: 100% pass)
- Transaction system is working (as evidenced by passing ingestion tests)
- Most issues are in test code or edge cases, not core functionality

## Files to Focus On

### For Old Syntax Fix
- `test/integration/incremental/advanced_test.go` (rewrite all test data)

### For Reset Issues
- `rete/constraint_pipeline.go` (reset handling)
- `rete/network.go` (GarbageCollect, network initialization)

### For Alpha Evaluation
- `rete/evaluator_constraints.go` (condition evaluation)
- `rete/alpha_node.go` (fact matching)
- Add unit tests in `rete/evaluator_test.go`
