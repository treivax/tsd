# Aggregation Calculation Bug Fix

## Bug Description

**Issue:** AccumulatorNode aggregation calculations were returning 0.00 even when aggregating over valid numeric values.

**Symptom:** 
- Test showed: `Valeur agrégée = 0.00 pour D1` 
- Expected: `Valeur agrégée = 55000.00` (average of 50000 and 60000)

**Root Cause:** The `calculateAggregateForFacts()` method only handled `float64` type values, but facts were being submitted with `int` type values.

## Code Location

**File:** `tsd/rete/node_accumulate.go`

**Affected Functions:**
- `calculateAggregateForFacts()` - Main aggregation calculation logic
- All aggregation branches: SUM, AVG, MIN, MAX

## The Problem

```go
// OLD CODE (BROKEN)
case "SUM":
    sum := 0.0
    for _, f := range facts {
        if val, ok := f.Fields[an.Field]; ok {
            if numVal, ok := val.(float64); ok {  // ❌ Only accepts float64
                sum += numVal
            }
        }
    }
    return sum, nil
```

When facts were submitted with integer values:
```go
Fields: map[string]interface{}{
    "salary": 50000,  // This is an int, not float64
}
```

The type assertion `val.(float64)` would fail, and the value would be silently skipped, resulting in aggregated value of 0.

## The Fix

Added a helper function `toFloat64()` that handles multiple numeric types:

```go
// NEW CODE (FIXED)
func (an *AccumulatorNode) toFloat64(val interface{}) float64 {
    switch v := val.(type) {
    case float64:
        return v
    case float32:
        return float64(v)
    case int:
        return float64(v)
    case int32:
        return float64(v)
    case int64:
        return float64(v)
    case uint:
        return float64(v)
    case uint32:
        return float64(v)
    case uint64:
        return float64(v)
    default:
        return 0
    }
}
```

Updated all aggregation calculations to use the helper:

```go
case "SUM":
    sum := 0.0
    for _, f := range facts {
        if val, ok := f.Fields[an.Field]; ok {
            numVal := an.toFloat64(val)  // ✅ Handles int, float64, etc.
            sum += numVal
        }
    }
    return sum, nil
```

## Changes Made

### Modified Functions in `node_accumulate.go`

1. **SUM aggregation** - Now uses `toFloat64(val)`
2. **AVG aggregation** - Now uses `toFloat64(val)`
3. **MIN aggregation** - Now uses `toFloat64(val)`
4. **MAX aggregation** - Now uses `toFloat64(val)`
5. **Added `toFloat64()` helper** - New function to convert all numeric types

### Files Modified

- `tsd/rete/node_accumulate.go` - Fixed aggregation calculations
- `tsd/rete/aggregation_calculation_test.go` - New comprehensive test suite (7 tests)

## Test Results

### Before Fix
```
📊 ACCUMULATOR[dept_avg_salary_accum]: Valeur agrégée = 0.00 pour D1
❌ Wrong result - should be 55000.00
```

### After Fix
```
📊 ACCUMULATOR[dept_avg_salary_accum]: Valeur agrégée = 50000.00 pour D1  (1 employee)
📊 ACCUMULATOR[dept_avg_salary_accum]: Valeur agrégée = 55000.00 pour D1  (2 employees)
📊 ACCUMULATOR[dept_avg_salary_accum]: Valeur agrégée = 60000.17 pour D1  (3 employees, mixed int/float)
✅ Correct results!
```

## New Test Coverage

Created `aggregation_calculation_test.go` with comprehensive tests:

1. ✅ **TestAggregationCalculation_AVG** - Average with int and float64 values
2. ✅ **TestAggregationCalculation_SUM** - Sum aggregation
3. ✅ **TestAggregationCalculation_COUNT** - Count aggregation
4. ✅ **TestAggregationCalculation_MIN** - Minimum value aggregation
5. ✅ **TestAggregationCalculation_MAX** - Maximum value aggregation
6. ✅ **TestAggregationCalculation_MultipleAggregates** - Multiple aggregations in one rule
7. ✅ **TestAggregationCalculation_EmptySet** - Aggregation with no matching facts

All tests pass with correct calculations.

## Example Usage

```tsd
type Department : <id: string, name: string>
type Employee : <id: string, deptId: string, salary: number>

rule dept_avg_salary : {d: Department, avg_sal: AVG(e.salary)} / {e: Employee} / e.deptId == d.id ==> print("Avg salary")
```

**Fact Submission:**
```go
emp1 := Fact{
    ID:   "E1",
    Type: "Employee",
    Fields: map[string]interface{}{
        "salary": 50000,  // int - now works correctly!
    },
}
```

**Result:**
- Average correctly calculated as 55000.00 for two employees with salaries 50000 and 60000

## Impact

✅ **Fixed:** All aggregation functions (AVG, SUM, COUNT, MIN, MAX) now work correctly with integer and floating-point values

✅ **Backward Compatible:** No breaking changes - the fix only adds support for more numeric types

✅ **Test Coverage:** Added 7 comprehensive tests covering all aggregation scenarios

✅ **Zero Regressions:** All existing tests continue to pass

## Verification

```bash
# Run aggregation calculation tests
cd rete && go test -v -run TestAggregationCalculation

# Run all RETE tests
cd rete && go test

# Result: All tests pass ✅
```

## Related Issues

- Fixes aggregation calculation bug in new aggregation-with-join syntax feature
- Resolves issue where aggregated values showed as 0.00 in test output
- Completes the aggregation-with-join feature implementation (parser + RETE integration + calculation)

---

**Status:** ✅ FIXED AND TESTED  
**Date:** January 2025  
**Files Modified:** 2  
**Tests Added:** 7  
**Tests Passing:** All (100%)