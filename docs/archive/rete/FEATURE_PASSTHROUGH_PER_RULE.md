# Feature: Passthrough Per-Rule with Alpha Filters

**Feature ID**: RETE-003  
**Status**: In Development  
**Priority**: High  
**Created**: 2025-01-XX  
**Related**: BUG_RETE001 (Alpha/Beta Separation), RETE-002 (Arithmetic Alpha)

## Overview

This feature fixes passthrough node sharing when multiple rules use the same variable types but have different alpha filters. Currently, a shared passthrough can propagate facts incorrectly when one rule's alpha filter should have blocked them.

## Problem Statement

### Current Behavior (Incorrect)

```tsd
type Order(id: string, amount: number)

rule large_orders : {o: Order} / o.amount > 100 ==> ...
rule very_large_orders : {o: Order} / o.amount > 500 ==> ...
```

**Current Network**:
```
TypeNode(Order)
  → AlphaNode[amount > 100]
  → AlphaNode[amount > 500]
  → PassthroughNode[SHARED]  ← PROBLEM: shared between both rules
    → TerminalNode[large_orders]
    → TerminalNode[very_large_orders]
```

**Problem**: Facts passing the first alpha filter (amount > 100) propagate through the shared passthrough to BOTH terminal nodes, even if they don't pass the second alpha filter (amount > 500).

### Example Failure

```
Submit: Order(id: "O1", amount: 150)
  ✓ Passes: amount > 100 → goes to large_orders (CORRECT)
  ✗ ALSO goes to very_large_orders (WRONG - should be blocked by amount > 500)
```

### Root Cause

From `rete/builder_join_rules.go`:

```go
// Current: Single passthrough per type
passthroughKey := fmt.Sprintf("passthrough_%s", typeName)
if existingNode, exists := network.PassthroughRegistry[passthroughKey]; exists {
    return existingNode // PROBLEM: reuses same node for different filters
}
```

## Solution Design

### Strategy 1: Passthrough Per-Rule (Recommended)

Each rule gets its own passthrough nodes, even if they share the same types.

**Pros**:
- Simple to implement
- Guarantees correctness
- Clear ownership: each rule controls its own path

**Cons**:
- More passthrough nodes (memory overhead)
- Less sharing

### Strategy 2: Passthrough Per-Alpha-Filter-Chain

Share passthrough only when the complete alpha filter chain is identical.

**Pros**:
- Maximum sharing while maintaining correctness
- Optimal memory usage

**Cons**:
- Complex: requires alpha chain hashing
- Registry management becomes more complex

### Strategy 3: Hybrid - Passthrough with Filter Context

Store filter context in passthrough node, check on propagation.

**Pros**:
- Maintains sharing
- Dynamic filtering

**Cons**:
- Runtime overhead (filter check on each propagation)
- Breaks RETE purity (passthrough should be dumb)

## Recommended: Strategy 1 (Passthrough Per-Rule)

### Implementation

**File**: `rete/builder_join_rules.go`

#### Current Code:
```go
func (b *JoinRuleBuilder) getOrCreatePassthrough(
    network *Network,
    typeName string,
    sourceNode Node,
) (*PassthroughNode, error) {
    passthroughKey := fmt.Sprintf("passthrough_%s", typeName)
    // ...
}
```

#### New Code:
```go
func (b *JoinRuleBuilder) getOrCreatePassthrough(
    network *Network,
    ruleName string,
    varName string,
    typeName string,
    sourceNode Node,
) (*PassthroughNode, error) {
    // Include rule name in key to prevent incorrect sharing
    passthroughKey := fmt.Sprintf("passthrough_%s_%s_%s", ruleName, varName, typeName)
    
    if existingNode, exists := network.PassthroughRegistry[passthroughKey]; exists {
        return existingNode, nil
    }
    
    passthrough := NewPassthroughNode(passthroughKey, typeName)
    passthrough.SetParent(sourceNode)
    
    network.PassthroughRegistry[passthroughKey] = passthrough
    
    if network.LifecycleManager != nil {
        network.LifecycleManager.RegisterNode(
            passthrough,
            passthroughKey,
            []string{ruleName},
        )
    }
    
    return passthrough, nil
}
```

### Updated Call Sites

All calls to `getOrCreatePassthrough` need to pass `ruleName`:

```go
// Before:
passthrough, err := b.getOrCreatePassthrough(network, typeName, lastNode)

// After:
passthrough, err := b.getOrCreatePassthrough(network, ruleName, varName, typeName, lastNode)
```

### Network Structure (After Fix)

```
TypeNode(Order)
  → AlphaNode[amount > 100]
      → PassthroughNode[large_orders_o_Order]
          → TerminalNode[large_orders]
  → AlphaNode[amount > 500]
      → PassthroughNode[very_large_orders_o_Order]
          → TerminalNode[very_large_orders]
```

**Result**: Each rule has its own path, filters work correctly.

## Testing Strategy

### Unit Tests

```go
func TestPassthroughPerRule_DifferentAlphaFilters(t *testing.T) {
    // Two rules: same types, different alpha filters
    content := `
    type Order(id: string, amount: number)
    action log(msg: string)
    
    rule large : {o: Order} / o.amount > 100 ==> log("Large")
    rule xlarge : {o: Order} / o.amount > 500 ==> log("XLarge")
    `
    
    network := buildNetwork(t, content)
    
    // Verify: 2 different passthrough nodes
    passthroughCount := 0
    for key := range network.PassthroughRegistry {
        if strings.Contains(key, "passthrough_") {
            passthroughCount++
        }
    }
    
    if passthroughCount != 2 {
        t.Errorf("Expected 2 passthrough nodes, got %d", passthroughCount)
    }
    
    // Test filtering behavior
    order150 := &Fact{ID: "O1", Type: "Order", Fields: map[string]interface{}{
        "id": "o1", "amount": 150.0,
    }}
    
    order600 := &Fact{ID: "O2", Type: "Order", Fields: map[string]interface{}{
        "id": "o2", "amount": 600.0,
    }}
    
    network.SubmitFact(order150)
    network.SubmitFact(order600)
    
    // order150 should activate only "large"
    largeTokens := network.TerminalNodes["large_terminal"].GetMemory().Tokens
    if len(largeTokens) != 2 {
        t.Errorf("Expected 2 activations for large (both orders), got %d", len(largeTokens))
    }
    
    // order600 should activate only "xlarge" (not order150)
    xlargeTokens := network.TerminalNodes["xlarge_terminal"].GetMemory().Tokens
    if len(xlargeTokens) != 1 {
        t.Errorf("Expected 1 activation for xlarge (only order600), got %d", len(xlargeTokens))
    }
}
```

### Integration Test

Re-enable the skipped test:

```go
func TestBetaBackwardCompatibility_JoinNodeSharing(t *testing.T) {
    // REMOVE: t.Skip("TODO: Fix passthrough sharing...")
    
    // Test should now pass with per-rule passthroughs
    // ...
}
```

## Migration & Compatibility

### Backward Compatibility

⚠️ **Breaking Change for Network Structure**

- More passthrough nodes will be created
- Memory usage may increase (small - typically < 1KB per passthrough)
- Network statistics will show higher node counts

### Performance Impact

**Memory**:
- Before: 1 passthrough per type
- After: 1 passthrough per rule per variable
- Impact: Negligible (passthrough nodes are lightweight)

**Runtime**:
- No performance degradation
- May actually improve (fewer incorrect propagations)

### Migration Path

No user action required - transparent optimization:

```bash
# Simply rebuild rules
go build && ./your-app
```

## Future Optimizations

### Phase 2: Smart Sharing (Optional)

Re-enable sharing when alpha filter chains are identical:

```go
// Hash the complete alpha chain
alphaChainHash := computeAlphaChainHash(alphaNodes)
passthroughKey := fmt.Sprintf("passthrough_%s_%s", typeName, alphaChainHash)

// Share passthrough if chains match
if existingNode, exists := network.PassthroughRegistry[passthroughKey]; exists {
    return existingNode
}
```

**When to implement**: After metrics show significant passthrough duplication (> 1000 nodes).

## Success Metrics

### Functional

- ✅ All alpha filters correctly block facts
- ✅ No cross-contamination between rules
- ✅ Test `TestBetaBackwardCompatibility_JoinNodeSharing` passes

### Performance

- ✅ Network build time: < 5% increase
- ✅ Fact propagation time: no degradation
- ✅ Memory overhead: < 10% increase for typical rulesets

### Quality

- ✅ 100% test coverage for new code
- ✅ All existing tests pass
- ✅ Documentation updated

## Implementation Checklist

- [ ] Update `getOrCreatePassthrough` signature (add `ruleName`, `varName`)
- [ ] Update all call sites in `builder_join_rules.go`
- [ ] Update call sites in `builder_alpha_rules.go` (if any)
- [ ] Add passthrough key format documentation
- [ ] Write unit test: `TestPassthroughPerRule_DifferentAlphaFilters`
- [ ] Write integration test: verify correct filtering
- [ ] Re-enable `TestBetaBackwardCompatibility_JoinNodeSharing`
- [ ] Run full test suite
- [ ] Update documentation
- [ ] Measure memory/performance impact

## Files to Modify

1. `rete/builder_join_rules.go`
   - `getOrCreatePassthrough()` signature and implementation
   - All call sites

2. `rete/builder_join_rules_test.go` (if exists)
   - Update tests with new signature

3. `rete/beta_backward_compatibility_test.go`
   - Remove `t.Skip()` from `TestBetaBackwardCompatibility_JoinNodeSharing`

4. `rete/passthrough_per_rule_test.go` (new)
   - Tests for per-rule passthrough behavior

## References

- **Related**: BUG_RETE001 (Alpha/Beta Separation)
- **Code**: `rete/builder_join_rules.go`, `rete/node_passthrough.go`
- **Theory**: RETE requires each rule to maintain independent fact propagation paths

## Approval

**Reviewers**:
- [ ] Architecture Review
- [ ] Performance Review

**Sign-off**:
- [ ] Lead Engineer
- [ ] QA Team

---

**Status**: Ready for implementation after RETE-002 (Arithmetic Alpha) is complete.