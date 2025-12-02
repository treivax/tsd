# Feature Implementation: Complete Arithmetic Decomposition for Alpha Expressions

## Objective

Implement the complete architecture for decomposing complex arithmetic alpha expressions into chains of atomic AlphaNodes with intermediate result propagation, as specified in `rete/ARITHMETIC_DECOMPOSITION_SPEC.md`.

## Context

The TSD RETE network currently handles arithmetic alpha expressions (e.g., `(c.qte * 23 - 10 + c.remise * 43) > 0`) as monolithic conditions evaluated in a single AlphaNode. While this works correctly and is optimized with AlphaSharingRegistry, it prevents fine-grained sharing of common sub-expressions across rules.

This feature implements a decomposition system that breaks complex arithmetic expressions into atomic steps with intermediate results that can be shared and cached.

## Reference Specification

Follow the detailed specification in: `rete/ARITHMETIC_DECOMPOSITION_SPEC.md`

Key sections:
- Architecture Required (lines 38-516)
- Complete Execution Flow (lines 517-597)
- Required Tests (lines 598-641)
- Implementation Plan (lines 710-740)

## Implementation Tasks

### Phase 1: Core Infrastructure

#### Task 1.1: Create EvaluationContext
**File:** `rete/evaluation_context.go`

Implement the `EvaluationContext` type that stores intermediate results during alpha chain evaluation:

```go
type EvaluationContext struct {
    OriginalFact        *Fact
    IntermediateResults map[string]interface{}
    EvaluationPath      []string
    Timestamp           time.Time
    Metadata            map[string]interface{}
}
```

Include methods:
- `NewEvaluationContext(fact *Fact) *EvaluationContext`
- `SetIntermediateResult(key string, value interface{})`
- `GetIntermediateResult(key string) (interface{}, bool)`
- `Clone() *EvaluationContext`

**Acceptance Criteria:**
- Context stores and retrieves intermediate results by key
- Clone creates deep copy of all maps
- EvaluationPath tracks execution order
- Thread-safe for concurrent access (add sync.RWMutex if needed)

#### Task 1.2: Extend AlphaNode Structure
**File:** `rete/alpha_node.go`

Add new fields to AlphaNode:
```go
type AlphaNode struct {
    // ... existing fields ...
    
    // NEW: Decomposition support
    ResultName   string      // Name of intermediate result produced (e.g., "temp_1")
    IsAtomic     bool        // true if atomic operation (single step)
    Dependencies []string    // Required intermediate results
}
```

**Acceptance Criteria:**
- Backward compatible with existing AlphaNode usage
- Fields properly initialized in constructors
- Serialization/deserialization handles new fields

#### Task 1.3: Implement ActivateWithContext
**File:** `rete/alpha_node.go`

Add context-aware activation method:

```go
func (an *AlphaNode) ActivateWithContext(fact *Fact, context *EvaluationContext) error
```

Implementation must:
1. Check all dependencies are satisfied in context
2. Evaluate condition using ConditionEvaluator
3. Store result in context if ResultName is set
4. Propagate to children with enriched context
5. Handle both AlphaNode and non-AlphaNode children

**Acceptance Criteria:**
- Dependencies validated before evaluation
- Errors returned for missing dependencies
- Context passed correctly to child nodes
- Backward compatible with standard Activate() method

#### Task 1.4: Create ConditionEvaluator
**File:** `rete/condition_evaluator.go`

Implement condition evaluator with context support:

```go
type ConditionEvaluator struct {
    storage Storage
}

func (ce *ConditionEvaluator) EvaluateWithContext(
    condition interface{},
    fact *Fact,
    context *EvaluationContext,
) (interface{}, error)
```

Must handle:
- `binaryOp` / `binaryOperation` - arithmetic operations
- `comparison` - comparison operations
- `fieldAccess` - field value extraction
- `number` / `numberLiteral` - literal values
- `tempResult` - intermediate result references ⚠️ KEY FEATURE

Include helper methods:
- `resolveTempResult()` - resolve intermediate result from context
- `evaluateBinaryOp()` - evaluate arithmetic operations
- `evaluateComparison()` - evaluate comparisons
- `applyOperator()` - apply arithmetic operators
- `toFloat64()` - type conversion

**Acceptance Criteria:**
- All condition types handled correctly
- TempResult references resolved from context
- Error handling for missing intermediate results
- Supports nested expressions
- Type conversions work for int, int64, float64

### Phase 2: Integration with Decomposer

#### Task 2.1: Update ArithmeticExpressionDecomposer
**File:** `rete/arithmetic_decomposer.go`

Modify `decomposeBinaryOp()` to generate tempResult references:

```go
func (aed *ArithmeticExpressionDecomposer) decomposeBinaryOp(
    expr map[string]interface{},
    steps *[]SimpleCondition,
    stepCounter *int,
) map[string]interface{}
```

Changes:
- Create unique result names: `fmt.Sprintf("temp_%d", *stepCounter)`
- Store result name in SimpleCondition
- Return tempResult reference map instead of original expression
- Handle nested expressions recursively

**Acceptance Criteria:**
- Each atomic step has unique result name
- TempResult references include correct step_name
- Dependency chain is correct
- Original expression behavior preserved

#### Task 2.2: Enhance AlphaChainBuilder
**File:** `rete/alpha_chain_builder.go`

Modify `BuildChain()` to set decomposition metadata:

```go
func (acb *AlphaChainBuilder) BuildChain(
    conditions []SimpleCondition,
    varName string,
) (*AlphaNode, error)
```

Changes:
- Set `ResultName` on each AlphaNode from SimpleCondition
- Set `IsAtomic = true` for decomposed nodes
- Extract and set `Dependencies` from condition
- Link nodes in correct dependency order

Add helper:
```go
func (acb *AlphaChainBuilder) extractDependencies(condition interface{}) []string
```

**Acceptance Criteria:**
- Each node in chain has correct ResultName
- Dependencies extracted from tempResult references
- Nodes linked in execution order
- Sharing works with decomposed chains

#### Task 2.3: Update JoinRuleBuilder Integration
**File:** `rete/builder_join_rules.go`

Modify `createBinaryJoinRule()` to use decomposition with context propagation:

```go
func (b *JoinRuleBuilder) createBinaryJoinRule(...) error
```

Changes:
- Enable arithmetic decomposition by default
- Use `ActivateWithContext()` for decomposed chains
- Create EvaluationContext for each fact activation
- Handle both decomposed and monolithic alpha nodes

Add configuration:
```go
type DecompositionConfig struct {
    Enabled           bool
    MinComplexity     int  // Minimum operations to trigger decomposition
    EnableSharing     bool
    EnableCaching     bool
}
```

**Acceptance Criteria:**
- Decomposition enabled for complex expressions
- Simple expressions use monolithic approach
- Context created and propagated correctly
- Backward compatible with existing rules

### Phase 3: Testing

#### Task 3.1: Unit Tests for EvaluationContext
**File:** `rete/evaluation_context_test.go`

Tests:
- `TestEvaluationContext_SetGet` - basic set/get operations
- `TestEvaluationContext_Clone` - deep copy verification
- `TestEvaluationContext_EvaluationPath` - path tracking
- `TestEvaluationContext_Concurrent` - thread safety

#### Task 3.2: Unit Tests for ConditionEvaluator
**File:** `rete/condition_evaluator_test.go`

Tests:
- `TestConditionEvaluator_BinaryOp` - arithmetic operations
- `TestConditionEvaluator_Comparison` - comparison operations
- `TestConditionEvaluator_TempResult` - intermediate result resolution
- `TestConditionEvaluator_NestedExpression` - complex nested cases
- `TestConditionEvaluator_MissingDependency` - error handling

#### Task 3.3: Integration Tests
**File:** `rete/arithmetic_decomposition_integration_test.go`

Tests:
- `TestArithmeticDecomposition_SimpleExpression` - basic decomposition
  - Expression: `c.qte * 23 > 100`
  - Expected: 2 steps (multiply, compare)
  
- `TestArithmeticDecomposition_ComplexExpression` - complex decomposition
  - Expression: `(c.qte * 23 - 10 + c.remise * 43) > 0`
  - Expected: 5 steps as per spec
  
- `TestAlphaChain_EvaluateWithContext` - context propagation
  - Create chain with 3 nodes
  - Verify intermediate results stored
  - Verify final result correct
  
- `TestChainSharing_DecomposedExpressions` - sharing verification
  - Two rules with common sub-expression `c.qte * 23`
  - Verify sub-expression shared (single AlphaNode)
  - Verify both rules produce correct results

#### Task 3.4: End-to-End Test Enhancement
**File:** `rete/action_arithmetic_e2e_test.go`

Enhance existing `TestArithmeticExpressionsE2E`:
- Add test case with decomposition enabled
- Print decomposition statistics:
  - Number of atomic steps per rule
  - Shared intermediate nodes
  - Intermediate result cache hits
- Compare results with monolithic mode
- Verify token counts identical

**Acceptance Criteria:**
- All tests pass
- Decomposed and monolithic modes produce same results
- Sharing detected and measured
- Performance benchmarks included

### Phase 4: Observability & Documentation

#### Task 4.1: Add Metrics Collection
**File:** `rete/decomposition_metrics.go`

Implement:
```go
type DecompositionMetrics struct {
    TotalExpressions         int64
    ExpressionsDecomposed    int64
    AverageStepsPerChain     float64
    IntermediateResultsCount int64
    SharedStepsRatio         float64
    EvaluationTimePerStep    time.Duration
}

func (dm *DecompositionMetrics) Record(...)
func (dm *DecompositionMetrics) Report() string
```

**Acceptance Criteria:**
- Metrics collected during network construction
- Metrics collected during evaluation
- Thread-safe counters
- Human-readable report format

#### Task 4.2: Add Debug Logging
**File:** Throughout implementation

Add structured logging:
- Context creation: log fact ID and timestamp
- Step evaluation: log step name, dependencies, result
- Sharing detection: log when nodes are reused
- Performance: log evaluation time per step

Use levels:
- DEBUG: detailed step-by-step execution
- INFO: decomposition summary
- WARN: missing dependencies, fallback to monolithic
- ERROR: evaluation failures

#### Task 4.3: Update Documentation
**Files:**
- `rete/ARITHMETIC_DECOMPOSITION_SPEC.md` - update status to "Implemented"
- `README.md` - add section on arithmetic decomposition
- `docs/architecture.md` - document new components
- Add inline godoc comments to all new types and functions

**Acceptance Criteria:**
- All public APIs documented
- Examples included in documentation
- Architecture diagrams updated
- Migration guide for existing code

### Phase 5: Performance & Safety

#### Task 5.1: Add Feature Flag
**File:** `rete/config.go`

Add configuration:
```go
type ReteConfig struct {
    // ... existing fields ...
    
    ArithmeticDecomposition DecompositionConfig
}
```

Environment variable: `TSD_ENABLE_ARITHMETIC_DECOMPOSITION=true|false`

**Acceptance Criteria:**
- Feature can be toggled at runtime
- Default: enabled for new deployments
- Graceful fallback to monolithic if disabled

#### Task 5.2: Add Safety Checks
**File:** Throughout implementation

Add validation:
- Circular dependency detection in chains
- Maximum chain depth limit (default: 20)
- Maximum intermediate results per context (default: 100)
- Timeout for long evaluations

**Acceptance Criteria:**
- Circular dependencies detected and rejected
- Deep recursion prevented
- Memory bounded
- Timeouts configurable

#### Task 5.3: Performance Benchmarks
**File:** `rete/arithmetic_decomposition_bench_test.go`

Benchmarks:
- `BenchmarkDecomposition_vs_Monolithic` - compare modes
- `BenchmarkContextCreation` - context overhead
- `BenchmarkIntermediateResultLookup` - map access cost
- `BenchmarkChainEvaluation_Depth5` - chain of 5 nodes
- `BenchmarkChainEvaluation_Depth10` - chain of 10 nodes
- `BenchmarkSharing_10Rules` - sharing efficiency

**Acceptance Criteria:**
- Decomposition overhead < 20% for simple expressions
- Sharing provides net benefit at 3+ rules
- Memory usage reasonable (< 1KB per context)
- Benchmarks run in CI

## Success Criteria

### Functional
- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] E2E test produces same results in both modes
- [ ] Complex expression `(c.qte * 23 - 10 + c.remise * 43) > 0` decomposes into 5 steps
- [ ] Common sub-expressions are shared across rules
- [ ] Intermediate results are correctly propagated

### Non-Functional
- [ ] Performance overhead < 20% for simple expressions
- [ ] Memory usage < 1KB per evaluation context
- [ ] Thread-safe implementation
- [ ] No race conditions (verified with `go test -race`)
- [ ] Backward compatible with existing code

### Documentation
- [ ] All public APIs documented with godoc
- [ ] Architecture documentation updated
- [ ] Migration guide provided
- [ ] Metrics and observability documented

### Quality
- [ ] Code coverage > 85% for new code
- [ ] No linter warnings
- [ ] All diagnostics resolved
- [ ] Performance benchmarks included

## Testing Strategy

1. **Unit Tests**: Test each component in isolation
2. **Integration Tests**: Test component interactions
3. **E2E Tests**: Test complete rule evaluation with decomposition
4. **Comparison Tests**: Verify decomposed == monolithic results
5. **Performance Tests**: Benchmark overhead and sharing benefits
6. **Concurrency Tests**: Verify thread safety with `go test -race`

## Rollout Plan

1. **Phase 1**: Implement with feature flag OFF by default
2. **Phase 2**: Enable in test environments, collect metrics
3. **Phase 3**: Enable for specific rules (complexity threshold)
4. **Phase 4**: Enable by default for all rules
5. **Phase 5**: Remove monolithic fallback (if metrics good)

## Risk Mitigation

- **Risk**: Performance regression
  - **Mitigation**: Feature flag, benchmarks, gradual rollout
  
- **Risk**: Bugs in new evaluation logic
  - **Mitigation**: Extensive tests, comparison with monolithic mode
  
- **Risk**: Increased complexity
  - **Mitigation**: Clear documentation, simple APIs, good tests
  
- **Risk**: Memory leaks
  - **Mitigation**: Context lifecycle management, bounded maps, tests

## Out of Scope

- Advanced optimizations (CSE - Common Subexpression Elimination)
- Persistent caching of intermediate results across facts
- Parallel evaluation of independent sub-expressions
- Dynamic rewriting of chains based on runtime statistics

These can be added in future iterations after measuring real-world usage.

## References

- Primary spec: `rete/ARITHMETIC_DECOMPOSITION_SPEC.md`
- Existing decomposer: `rete/arithmetic_decomposer.go`
- Existing chain builder: `rete/alpha_chain_builder.go`
- E2E test: `rete/action_arithmetic_e2e_test.go`
- Join rule builder: `rete/builder_join_rules.go`

## Notes

- Follow Go idioms and existing code style
- Use existing test patterns and utilities
- Reuse `AlphaSharingRegistry` for node sharing
- Keep backward compatibility with existing AlphaNode usage
- Add metrics and logging for observability
- Document all design decisions in code comments