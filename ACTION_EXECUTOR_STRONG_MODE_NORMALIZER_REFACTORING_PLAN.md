# 🔄 REFACTORING PLAN: action_executor.go, strong_mode_performance.go & nested_or_normalizer.go

## 📋 Executive Summary

**Objective**: Refactor three large files in the rete package to improve maintainability, readability, and testability while preserving all existing behavior.

**Files to refactor**:
1. `rete/action_executor.go` (619 lines, 35 functions)
2. `rete/strong_mode_performance.go` (596 lines, 68 symbols)
3. `rete/nested_or_normalizer.go` (623 lines, 34 functions)

**Total**: 1,838 lines across 3 files

**Approach**: Incremental, behavior-preserving refactoring with validation after each step.

---

## 🎯 Goals

### ✅ Primary Objectives

1. **Reduce file size and complexity**
   - Split large files into focused modules
   - Each module should have a single, clear responsibility
   - Target: < 300 lines per file

2. **Improve code organization**
   - Group related functions logically
   - Separate concerns (execution, validation, evaluation, metrics, normalization)
   - Make navigation and maintenance easier

3. **Preserve behavior 100%**
   - No changes to public APIs
   - All existing tests must pass without modification
   - No performance regressions

4. **Maintain backward compatibility**
   - Keep all public function signatures unchanged
   - Preserve all existing constants and types
   - Support both old and new code patterns

---

## 📊 Current State Analysis

### File 1: rete/action_executor.go (619 lines)

**Structure**:
- ActionExecutor type (lines 14-19): Main executor struct
- Constructor & configuration (lines 22-61): NewActionExecutor, RegisterDefaultActions, SetLogging
- Execution logic (lines 64-123): ExecuteAction, executeJob
- Argument evaluation (lines 126-205): evaluateArgument (80 lines - very long!)
- Fact operations (lines 208-303): evaluateFactCreation, evaluateFactModification
- Arithmetic evaluation (lines 306-413): evaluateArithmetic, evaluateBinaryOperation, evaluateComparison
- Field validation (lines 430-498): validateFactFields, validateFieldValue, validateFieldType
- Helpers (lines 501-589): logAction, formatArgument, formatArgs, toNumber, generateFactID
- ExecutionContext type (lines 592-619): Context management

**Identified Issues**:
- ❌ Too many responsibilities in one file (execution, evaluation, validation, formatting)
- ❌ Very long functions: evaluateArgument (80 lines), evaluateFactModification (52 lines)
- ❌ Mixed concerns: execution + validation + evaluation + formatting
- ❌ Helper functions mixed with core logic
- ❌ No clear separation between public API and private helpers

**Code Smells**:
- Long method: evaluateArgument (80 lines with nested switch)
- Feature envy: Many functions access ctx fields
- Data clumps: Arguments passed around repeatedly
- Magic strings: "type", "variable", "fieldAccess", etc.

### File 2: rete/strong_mode_performance.go (596 lines)

**Structure**:
- StrongModePerformanceMetrics type (lines 14-71): Large struct with 40+ fields
- ConfigChange type (lines 74-82): Configuration tracking
- Constructor (lines 85-97): NewStrongModePerformanceMetrics
- Recording methods (lines 100-208): RecordTransaction, RecordCommit, RecordConfigChange
- Analysis methods (lines 211-326): updateHealthIndicators, generateRecommendations, getTopRollbackReasons
- Report generation (lines 358-454): GetReport (97 lines - very long!)
- Calculation helpers (lines 457-504): Various rate calculations
- Summary methods (lines 523-596): GetSummary, Clone

**Identified Issues**:
- ❌ God object: StrongModePerformanceMetrics has too many responsibilities
- ❌ Very long methods: GetReport (97 lines), RecordTransaction (70 lines)
- ❌ Mixed concerns: metrics collection + analysis + reporting + health checking
- ❌ Large struct with 40+ fields makes it hard to understand
- ❌ Duplication in rate calculation methods

**Code Smells**:
- Large class: 40+ fields in main struct
- Long method: GetReport (97 lines)
- Duplication: Multiple similar rate calculation functions
- Primitive obsession: Many float64 fields for rates/percentages
- Feature envy: Methods accessing many fields of the same struct

### File 3: rete/nested_or_normalizer.go (623 lines)

**Structure**:
- NestedORComplexity constants (lines 16-28): Complexity levels
- NestedORAnalysis type (lines 32-40): Analysis result
- Analysis functions (lines 43-196): AnalyzeNestedOR, analyzeLogicalExpressionNesting, analyzeMapExpressionNesting
- Flattening functions (lines 199-392): FlattenNestedOR, flattenLogicalExpression, flattenMapExpression
- DNF transformation (lines 396-578): TransformToDNF, transformLogicalExpressionToDNF, normalizeDNFTerms
- Main normalization (lines 589-623): NormalizeNestedOR

**Identified Issues**:
- ❌ Three major responsibilities: analysis, flattening, DNF transformation
- ❌ Long functions: normalizeDNFTerms (51 lines), generateDNFTerms (41 lines)
- ❌ Duplication: Similar patterns for LogicalExpression vs Map handling
- ❌ Complex nested logic in transformation functions
- ❌ No clear separation of concerns

**Code Smells**:
- Long methods: normalizeDNFTerms (51 lines)
- Duplication: Parallel functions for two formats (LogicalExpression vs Map)
- Complex conditionals: Deep nesting in transformation logic
- Magic numbers: Hardcoded complexity thresholds

---

## 🗂️ Proposed File Structure

### For rete/action_executor.go → Split into 6 files:

1. **rete/action_executor.go** (REFACTORED)
   - ActionExecutor type
   - NewActionExecutor, RegisterDefaultActions, SetLogging
   - ExecuteAction, executeJob (orchestration only)
   - ~120 lines

2. **rete/action_executor_evaluation.go** (NEW)
   - evaluateArgument (refactored into smaller functions)
   - evaluateArithmetic
   - evaluateBinaryOperation
   - evaluateArithmeticOperation
   - evaluateComparison
   - areEqual
   - toNumber
   - ~180 lines

3. **rete/action_executor_facts.go** (NEW)
   - evaluateFactCreation
   - evaluateFactModification
   - generateFactID
   - factCounter
   - ~100 lines

4. **rete/action_executor_validation.go** (NEW)
   - validateFactFields
   - validateFieldValue
   - validateFieldType
   - ~80 lines

5. **rete/action_executor_context.go** (NEW)
   - ExecutionContext type
   - NewExecutionContext
   - GetVariable
   - ~40 lines

6. **rete/action_executor_helpers.go** (NEW)
   - logAction
   - formatArgument
   - formatArgs
   - ~80 lines

### For rete/strong_mode_performance.go → Split into 6 files:

1. **rete/strong_mode_performance.go** (REFACTORED)
   - StrongModePerformanceMetrics type (simplified)
   - NewStrongModePerformanceMetrics
   - RecordTransaction
   - RecordCommit
   - RecordConfigChange
   - ~150 lines

2. **rete/strong_mode_performance_types.go** (NEW)
   - ConfigChange type
   - Helper types (reasonCount, etc.)
   - ~40 lines

3. **rete/strong_mode_performance_health.go** (NEW)
   - updateHealthIndicators
   - getHealthStatus
   - Health-related calculations
   - ~100 lines

4. **rete/strong_mode_performance_analysis.go** (NEW)
   - generateRecommendations
   - getTopRollbackReasons
   - Analysis helper functions
   - ~100 lines

5. **rete/strong_mode_performance_calculations.go** (NEW)
   - All rate calculation methods (getSuccessRate, etc.)
   - Aggregation helpers
   - ~100 lines

6. **rete/strong_mode_performance_reporting.go** (NEW)
   - GetReport
   - GetSummary
   - formatRecommendations
   - Clone
   - ~120 lines

### For rete/nested_or_normalizer.go → Split into 5 files:

1. **rete/nested_or_normalizer.go** (REFACTORED)
   - NestedORComplexity constants
   - NestedORAnalysis type
   - NormalizeNestedOR (main entry point)
   - ~80 lines

2. **rete/nested_or_normalizer_analysis.go** (NEW)
   - AnalyzeNestedOR
   - analyzeLogicalExpressionNesting
   - analyzeMapExpressionNesting
   - ~160 lines

3. **rete/nested_or_normalizer_flattening.go** (NEW)
   - FlattenNestedOR
   - flattenLogicalExpression
   - flattenMapExpression
   - collectORTermsRecursive
   - collectORTermsFromMap
   - hasOnlyOR, hasOnlyORInMap
   - ~200 lines

4. **rete/nested_or_normalizer_dnf.go** (NEW)
   - TransformToDNF
   - transformLogicalExpressionToDNF
   - transformMapToDNF
   - extractANDGroups
   - generateDNFTerms
   - ~140 lines

5. **rete/nested_or_normalizer_helpers.go** (NEW)
   - normalizeDNFTerms
   - Canonicalization helpers
   - termWithCanonical type
   - ~80 lines

---

## 📝 Detailed Refactoring Steps

### Phase 1: rete/action_executor.go

#### Step 1: Extract ExecutionContext ✅
**Action**: Create `rete/action_executor_context.go`
- Move ExecutionContext type
- Move NewExecutionContext
- Move GetVariable
- Add MIT license header
- Run tests

**Validation**:
```bash
go test ./rete/...
go vet ./rete/...
```

#### Step 2: Extract Helpers ✅
**Action**: Create `rete/action_executor_helpers.go`
- Move logAction
- Move formatArgument
- Move formatArgs
- Run tests

#### Step 3: Extract Fact Operations ✅
**Action**: Create `rete/action_executor_facts.go`
- Move evaluateFactCreation
- Move evaluateFactModification
- Move generateFactID
- Move factCounter variable
- Run tests

#### Step 4: Extract Validation ✅
**Action**: Create `rete/action_executor_validation.go`
- Move validateFactFields
- Move validateFieldValue
- Move validateFieldType
- Run tests

#### Step 5: Extract Evaluation Logic ✅
**Action**: Create `rete/action_executor_evaluation.go`
- Move evaluateArgument
- Move evaluateArithmetic
- Move evaluateBinaryOperation
- Move evaluateArithmeticOperation
- Move evaluateComparison
- Move areEqual
- Move toNumber
- Run tests

#### Step 6: Clean Up Main File ✅
**Action**: Refactor `rete/action_executor.go`
- Keep only ActionExecutor type
- Keep constructor and public API
- Keep orchestration methods (ExecuteAction, executeJob)
- Ensure clear delegation to helper modules
- Run tests

### Phase 2: rete/strong_mode_performance.go

#### Step 7: Extract Types ✅
**Action**: Create `rete/strong_mode_performance_types.go`
- Move ConfigChange type
- Move reasonCount type
- Move other helper types
- Add MIT license header
- Run tests

#### Step 8: Extract Calculations ✅
**Action**: Create `rete/strong_mode_performance_calculations.go`
- Move all getXxxRate methods
- Move calculation helper functions
- Run tests

#### Step 9: Extract Health Indicators ✅
**Action**: Create `rete/strong_mode_performance_health.go`
- Move updateHealthIndicators
- Move getHealthStatus
- Move health-related logic
- Run tests

#### Step 10: Extract Analysis ✅
**Action**: Create `rete/strong_mode_performance_analysis.go`
- Move generateRecommendations
- Move getTopRollbackReasons
- Move analysis helper functions
- Run tests

#### Step 11: Extract Reporting ✅
**Action**: Create `rete/strong_mode_performance_reporting.go`
- Move GetReport (refactor if needed)
- Move GetSummary
- Move formatRecommendations
- Move Clone
- Run tests

#### Step 12: Clean Up Main File ✅
**Action**: Refactor `rete/strong_mode_performance.go`
- Keep only main type and constructor
- Keep recording methods
- Keep public API
- Ensure delegation to helper modules
- Run tests

### Phase 3: rete/nested_or_normalizer.go

#### Step 13: Extract Analysis ✅
**Action**: Create `rete/nested_or_normalizer_analysis.go`
- Move AnalyzeNestedOR
- Move analyzeLogicalExpressionNesting
- Move analyzeMapExpressionNesting
- Add MIT license header
- Run tests

#### Step 14: Extract Flattening ✅
**Action**: Create `rete/nested_or_normalizer_flattening.go`
- Move FlattenNestedOR
- Move flattenLogicalExpression
- Move flattenMapExpression
- Move collectORTermsRecursive
- Move collectORTermsFromMap
- Move hasOnlyOR, hasOnlyORInMap
- Run tests

#### Step 15: Extract DNF Transformation ✅
**Action**: Create `rete/nested_or_normalizer_dnf.go`
- Move TransformToDNF
- Move transformLogicalExpressionToDNF
- Move transformMapToDNF
- Move extractANDGroups
- Move generateDNFTerms
- Run tests

#### Step 16: Extract Helpers ✅
**Action**: Create `rete/nested_or_normalizer_helpers.go`
- Move normalizeDNFTerms
- Move termWithCanonical type
- Move canonicalization helpers
- Run tests

#### Step 17: Clean Up Main File ✅
**Action**: Refactor `rete/nested_or_normalizer.go`
- Keep only constants and main type
- Keep NormalizeNestedOR (main entry point)
- Ensure delegation to helper modules
- Add clear documentation
- Run tests

---

## ✅ Validation Strategy

### After Each Step

1. **Compile Check**:
   ```bash
   go build ./rete/...
   ```

2. **Unit Tests**:
   ```bash
   go test ./rete/...
   ```

3. **Static Analysis**:
   ```bash
   go vet ./rete/...
   golangci-lint run ./rete/...
   ```

4. **Integration Tests**:
   ```bash
   make test
   make rete-unified
   ```

### Final Validation

1. **Full Test Suite**:
   ```bash
   make test-all
   ```

2. **Benchmarks** (ensure no performance regression):
   ```bash
   go test -bench=. -benchmem ./rete/...
   ```

3. **Coverage Check**:
   ```bash
   go test -cover ./rete/...
   ```

---

## 📦 Commit Strategy

Each phase will have its own commit:

**Phase 1 Commit**:
```
refactor(rete): split action_executor.go into focused modules

- Extract ExecutionContext to action_executor_context.go
- Extract helpers to action_executor_helpers.go
- Extract fact operations to action_executor_facts.go
- Extract validation to action_executor_validation.go
- Extract evaluation logic to action_executor_evaluation.go
- Refactor main file to focus on orchestration

Behavior preserved: All tests passing
Tests: ✅ go test ./rete/...
Lint: ✅ go vet ./rete/...
```

**Phase 2 Commit**:
```
refactor(rete): split strong_mode_performance.go into focused modules

- Extract types to strong_mode_performance_types.go
- Extract calculations to strong_mode_performance_calculations.go
- Extract health indicators to strong_mode_performance_health.go
- Extract analysis to strong_mode_performance_analysis.go
- Extract reporting to strong_mode_performance_reporting.go
- Refactor main file to focus on metrics collection

Behavior preserved: All tests passing
Tests: ✅ go test ./rete/...
Lint: ✅ go vet ./rete/...
```

**Phase 3 Commit**:
```
refactor(rete): split nested_or_normalizer.go into focused modules

- Extract analysis to nested_or_normalizer_analysis.go
- Extract flattening to nested_or_normalizer_flattening.go
- Extract DNF transformation to nested_or_normalizer_dnf.go
- Extract helpers to nested_or_normalizer_helpers.go
- Refactor main file to focus on main entry point

Behavior preserved: All tests passing
Tests: ✅ go test ./rete/...
Lint: ✅ go vet ./rete/...
```

---

## 🎓 Expected Outcomes

### Before Refactoring

**Total**:
- Files: 3
- Lines: 1,838
- Functions: ~137
- Average lines/file: 613
- Max lines/file: 623
- Responsibilities: Mixed (5-7 per file)

### After Refactoring

**Total**:
- Files: 17 (from 3)
- Lines: ~1,900 (accounting for headers)
- Functions: ~137 (same)
- Average lines/file: ~112
- Max lines/file: ~200
- Responsibilities: 1-2 per file

### Quality Improvements

1. ✅ **Readability**: Clear separation of concerns
2. ✅ **Maintainability**: Small, focused files
3. ✅ **Testability**: Isolated functions easier to test
4. ✅ **Navigation**: File names indicate contents
5. ✅ **Onboarding**: Modules can be understood independently
6. ✅ **Reusability**: Extracted helpers can be reused
7. ✅ **Behavior**: 100% preserved

---

## 🚀 Execution Timeline

**Estimated time**: 4-5 hours

- Phase 1 (action_executor.go): 1.5-2 hours
  - Steps 1-6: 15-20 minutes each

- Phase 2 (strong_mode_performance.go): 1.5-2 hours
  - Steps 7-12: 15-20 minutes each

- Phase 3 (nested_or_normalizer.go): 1.5-2 hours
  - Steps 13-17: 15-20 minutes each

---

## 🔒 Risk Mitigation

1. **Incremental approach**: Each step is small and reversible
2. **Test after each step**: Catch issues immediately
3. **Preserve all tests**: No test modifications needed
4. **Clear commits**: Easy to identify and revert problematic changes
5. **No API changes**: External code unaffected

---

## ✅ Success Criteria

- [ ] All files < 300 lines
- [ ] Each file has single responsibility
- [ ] All tests passing
- [ ] No lint errors
- [ ] No performance regression
- [ ] All public APIs unchanged
- [ ] MIT license headers on all new files
- [ ] Clear package documentation
- [ ] No hardcoded values introduced
- [ ] Code remains generic and reusable

---

**Ready to begin execution!**