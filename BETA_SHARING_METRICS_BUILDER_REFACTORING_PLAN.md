# Refactoring Plan: Beta Sharing, Metrics, and Join Rule Builder

## Overview
This document outlines the refactoring plan for three large files in the `rete` package:
- `rete/beta_sharing.go` (~729 lines)
- `rete/arithmetic_decomposition_metrics.go` (~713 lines)
- `rete/builder_join_rules.go` (~759 lines)

## Objectives
1. **Improve Maintainability**: Split large files into focused, cohesive modules
2. **Preserve Behavior**: Ensure all existing tests pass without modification
3. **Maintain Public APIs**: No changes to exported types, functions, or signatures
4. **Enhance Readability**: Group related functionality for easier navigation
5. **Follow Standards**: Add MIT license headers to all new files

## File 1: beta_sharing.go Refactoring

### Current Structure (~729 lines)
- BetaSharingRegistryImpl methods (lifecycle, node management)
- Hash computation and normalization
- DefaultJoinNodeNormalizer implementation
- DefaultJoinNodeHasher implementation
- Standalone helper functions

### Target Structure (4 files)

#### 1. `beta_sharing.go` - Core Registry (~250 lines)
**Responsibility**: Core registry implementation and node lifecycle
**Contents**:
- `GetOrCreateJoinNode` - Main entry point for node creation/retrieval
- `RegisterJoinNode` - Explicit node registration
- `AddRuleToJoinNode` - Associate rule with node
- `RemoveRuleFromJoinNode` - Remove rule reference
- `RegisterRuleForJoinNode` - Register rule usage
- `UnregisterJoinNode` - Complete node removal
- `GetJoinNodeRules` - Get rules for a node
- `GetJoinNodeRefCount` - Get reference count
- `Shutdown` - Cleanup and resource release
- `ClearCache` - Clear hash cache

#### 2. `beta_sharing_hash.go` - Hash & Normalization (~250 lines)
**Responsibility**: Hash computation, normalization, and canonical representation
**Contents**:
- `computeHashDirect` - Direct hash computation fallback
- `normalizeSignatureFallback` - Basic normalization when no normalizer configured
- `defaultJoinNodeNormalizer` type and implementation
  - `NewDefaultJoinNodeNormalizer`
  - `NormalizeSignature`
  - `NormalizeCondition`
  - `normalizeConditionMap`
- `defaultJoinNodeHasher` type and implementation
  - `NewDefaultJoinNodeHasher`
  - `ComputeHash`
  - `ComputeHashCached`

#### 3. `beta_sharing_stats.go` - Statistics & Introspection (~150 lines)
**Responsibility**: Statistics gathering, metrics, and node introspection
**Contents**:
- `ReleaseJoinNode` - Decrement refcount and remove if unused
- `ReleaseJoinNodeByID` - Release by node ID
- `GetSharingStats` - Current sharing metrics
- `ListSharedJoinNodes` - List all shared node hashes
- `GetSharedJoinNodeDetails` - Detailed node information

#### 4. `beta_sharing_helpers.go` - Standalone Helpers (~80 lines)
**Responsibility**: Backward-compatible standalone helper functions
**Contents**:
- `NormalizeJoinCondition` - Standalone normalization function
- `ComputeJoinHash` - Standalone hash computation function

## File 2: arithmetic_decomposition_metrics.go Refactoring

### Current Structure (~713 lines)
- Type definitions (4 types: ArithmeticDecompositionMetrics, RuleArithmeticMetrics, GlobalArithmeticMetrics, MetricsConfig)
- Constructor and configuration
- Recording functions (6 functions)
- Query/retrieval functions (7 functions)
- Private helper functions (9 functions)

### Target Structure (3 files)

#### 1. `arithmetic_decomposition_metrics.go` - Core Types & Recording (~350 lines)
**Responsibility**: Type definitions, constructor, and metric recording
**Contents**:
- Type definitions:
  - `ArithmeticDecompositionMetrics`
  - `RuleArithmeticMetrics`
  - `GlobalArithmeticMetrics`
  - `MetricsConfig`
- Constructor:
  - `DefaultMetricsConfig`
  - `NewArithmeticDecompositionMetrics`
- Recording functions:
  - `RecordActivation`
  - `RecordEvaluation`
  - `RecordCacheHit`
  - `RecordCacheMiss`
  - `RecordChainStructure`
  - `RecordCircularDependency`
  - `RecordGraphValidation`
  - `UpdateCacheStatistics`
  - `Reset`

#### 2. `arithmetic_decomposition_metrics_query.go` - Query & Retrieval (~200 lines)
**Responsibility**: Metric queries, retrieval, and aggregation
**Contents**:
- `GetRuleMetrics` - Get metrics for specific rule
- `GetGlobalMetrics` - Get global aggregated metrics
- `GetAllRuleMetrics` - Get all rule metrics
- `GetTopRulesByEvaluations` - Top N by evaluation count
- `GetTopRulesByDuration` - Top N by total duration
- `GetSlowestRules` - Top N by average duration
- `GetSummary` - Formatted summary

#### 3. `arithmetic_decomposition_metrics_helpers.go` - Private Helpers (~200 lines)
**Responsibility**: Internal calculations and helper functions
**Contents**:
- `getOrCreateRuleMetrics` - Get or create rule metrics
- `getHistogramBucket` - Find histogram bucket
- `updateCacheHitRate` - Update per-rule cache hit rate
- `updateGlobalCacheHitRate` - Update global cache hit rate
- `calculateMaxDepth` - Calculate dependency graph depth
- `recalculateGlobalAverages` - Recalculate global averages
- `calculatePercentiles` - Calculate time percentiles
- `evictOldestRule` - LRU eviction
- `copyRuleMetrics` - Deep copy for thread safety

## File 3: builder_join_rules.go Refactoring

### Current Structure (~759 lines)
- JoinRuleBuilder type and configuration
- CreateJoinRule entry point
- Binary join creation (large function ~213 lines)
- Cascade join creation (large function ~80 lines + helper ~80 lines)
- Helper functions (4 functions for chain building and connection)

### Target Structure (3 files)

#### 1. `builder_join_rules.go` - Core Type & Entry Point (~100 lines)
**Responsibility**: Type definition, constructor, configuration, and entry point
**Contents**:
- `JoinRuleBuilder` type definition
- `NewJoinRuleBuilder` - Constructor
- `SetDecompositionEnabled` - Configuration
- `SetDecompositionComplexity` - Configuration
- `CreateJoinRule` - Main entry point (delegates to binary or cascade)
- `createCascadeJoinRule` - Delegation wrapper

#### 2. `builder_join_rules_binary.go` - Binary Join (~350 lines)
**Responsibility**: Binary join rule creation (2 variables)
**Contents**:
- `createBinaryJoinRule` - Complete binary join implementation
  - Condition splitting (alpha/beta)
  - AlphaNode creation with decomposition
  - JoinNode creation with sharing
  - Network connection with alpha filters

#### 3. `builder_join_rules_cascade.go` - Cascade Join & Helpers (~350 lines)
**Responsibility**: Cascade join creation (3+ variables) and helper functions
**Contents**:
- `createCascadeJoinRuleWithBuilder` - Cascade implementation with BetaChainBuilder
  - Alpha/beta condition splitting
  - AlphaNode creation
  - Chain pattern building
  - Network connection
- Helper functions:
  - `buildJoinPatterns` - Create join patterns for beta chain
  - `buildChainWithBuilder` - Build chain using BetaChainBuilder
  - `connectChainToNetwork` - Connect chain to network (basic)
  - `connectChainToNetworkWithAlpha` - Connect chain with alpha integration

## Validation Strategy

### Phase 1: Pre-Refactoring Validation
1. Run all tests: `go test ./rete/...`
2. Run static analysis: `go vet ./rete/...`
3. Record test coverage baseline
4. Document current behavior

### Phase 2: Incremental Refactoring
For each file group:
1. Create new files with extracted code
2. Update imports and references
3. Run tests after each extraction
4. Verify no behavior changes

### Phase 3: Post-Refactoring Validation
1. Run all tests: `go test ./rete/...`
2. Verify all tests pass unchanged
3. Run static analysis: `go vet ./rete/...`
4. Check test coverage (should be same or better)
5. Build entire project: `go build ./...`

## Test Coverage

### Existing Tests to Verify
**Beta Sharing**:
- `beta_sharing_coverage_test.go` - Component-level tests
- `beta_sharing_integration_test.go` - Integration tests
- Tests in other files that use BetaSharingRegistry

**Arithmetic Decomposition Metrics**:
- `arithmetic_decomposition_metrics_test.go` - Unit tests
- `arithmetic_decomposition_integration_test.go` - Integration tests
- `arithmetic_node_sharing_validation_test.go` - Sharing validation

**Join Rule Builder**:
- `alpha_filters_diagnostic_test.go` - Alpha filter tests
- `arithmetic_decomposition_integration_test.go` - Decomposition with joins
- Various integration tests using join rules

## Success Criteria

### Must Have
- [ ] All existing tests pass without modification
- [ ] No changes to public APIs (exported types/functions)
- [ ] All new files have MIT license headers
- [ ] `go vet` passes with no warnings
- [ ] `go build ./...` succeeds
- [ ] No behavioral changes detected

### Should Have
- [ ] Code is more readable and maintainable
- [ ] Related functionality is grouped logically
- [ ] File sizes are reasonable (<400 lines per file)
- [ ] Documentation is clear and accurate

### Nice to Have
- [ ] Improved code organization enables future enhancements
- [ ] Easier to add new features to focused modules
- [ ] Better separation of concerns

## Execution Order

1. **beta_sharing.go** (First - foundational component)
   - Create `beta_sharing_hash.go`
   - Create `beta_sharing_stats.go`
   - Create `beta_sharing_helpers.go`
   - Update `beta_sharing.go`
   - Test and validate

2. **arithmetic_decomposition_metrics.go** (Second - independent component)
   - Create `arithmetic_decomposition_metrics_query.go`
   - Create `arithmetic_decomposition_metrics_helpers.go`
   - Update `arithmetic_decomposition_metrics.go`
   - Test and validate

3. **builder_join_rules.go** (Third - depends on both above)
   - Create `builder_join_rules_binary.go`
   - Create `builder_join_rules_cascade.go`
   - Update `builder_join_rules.go`
   - Test and validate

## Risk Mitigation

### Potential Risks
1. **Breaking internal dependencies**: Careful import management
2. **Race conditions**: Preserve all mutex usage patterns
3. **Test failures**: Incremental approach with validation
4. **Performance regression**: No algorithmic changes

### Mitigation Strategies
1. Extract one file at a time
2. Run tests after each extraction
3. Use git branches for safe experimentation
4. Keep original files until validation complete
5. Document any edge cases discovered

## Rollback Plan

If issues are discovered:
1. Each refactoring is a separate commit
2. Can revert individual file extractions
3. Original files preserved until full validation
4. Git history maintains complete audit trail

## Timeline

- **Phase 1** (beta_sharing): 30-45 minutes
- **Phase 2** (arithmetic_decomposition_metrics): 30-45 minutes
- **Phase 3** (builder_join_rules): 30-45 minutes
- **Total**: ~2-2.5 hours including testing and validation

## Notes

- Follow existing code style and conventions
- Preserve all comments and documentation
- Maintain consistent error handling patterns
- Keep thread-safety guarantees intact
- No optimization or feature changes during refactoring