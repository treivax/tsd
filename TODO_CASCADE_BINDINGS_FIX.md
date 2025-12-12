# TODO: Fix Cascade Join Bindings Loss Issue

## Problem Description

**Status**: CRITICAL BUG
**Affects**: 3+ variable cascade joins
**Tests Failing**: 
- `beta_join_complex.tsd`
- `join_multi_variable_complex.tsd` 
- `beta_exhaustive_coverage.tsd`

### Symptoms

When a rule has 3 or more variables (e.g., `{u: User, o: Order, p: Product}`), the cascade join architecture creates multiple JoinNodes:
- Join1: u ⋈ o
- Join2: (u,o) ⋈ p

The problem is that when the final token reaches the TerminalNode, it only contains bindings for 2 variables instead of all 3.

**Example error**:
```
erreur évaluation argument 1: ❌ Erreur d'exécution d'action:
   Variable 'p' non trouvée dans le contexte
   Variables disponibles: [u o]
```

Expected: `[u, o, p]`
Actual: `[u, o]`

### Root Cause Analysis

#### What Works ✅
1. **BindingChain.Merge()**: Tested independently and works correctly
   - `chain1 [u]`.Merge(`chain2 [o]`) = `[u, o]` ✅
   - `chain1 [u, o]`.Merge(`chain2 [p]`) = `[u, o, p]` ✅

2. **Token creation for right facts**: `createTokenForRightFact()` correctly creates bindings

3. **Pattern building**: `buildJoinPatterns()` correctly creates patterns with AllVars

#### What Doesn't Work ❌
The joined token from the second cascade join loses bindings when propagated to the terminal node.

#### Debugging Challenges
- `fmt.Printf()` statements in the code don't appear in test output (possibly suppressed by test framework or buffered)
- Test cache issues required `go clean -testcache` to see changes
- The error occurs during token propagation, making it hard to trace

### Investigation Steps Completed

1. ✅ Verified BindingChain.Merge() works correctly in isolation
2. ✅ Checked NewJoinNode() correctly computes AllVariables = leftVars + rightVars  
3. ✅ Verified GetOrCreateJoinNode() passes allVars but doesn't use it (not a bug - NewJoinNode recomputes it)
4. ✅ Checked passthrough alpha nodes create tokens with bindings correctly
5. ✅ Verified PropagateToChildren() passes tokens as-is (doesn't modify them)
6. ⚠️  Unable to get debug output from performJoinWithTokens() during tests

### Next Steps

#### Priority 1: Add Comprehensive Logging

Create a diagnostic mode that logs to a file instead of stdout:

```go
// In node_join.go performJoinWithTokens()
if len(jn.AllVariables) >= 3 {
    logFile, _ := os.OpenFile("/tmp/cascade_join_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    defer logFile.Close()
    fmt.Fprintf(logFile, "\n=== JOIN %s ===\n", jn.ID)
    fmt.Fprintf(logFile, "AllVariables: %v\n", jn.AllVariables)
    fmt.Fprintf(logFile, "Token1 bindings: %v\n", token1.GetVariables())
    fmt.Fprintf(logFile, "Token2 bindings: %v\n", token2.GetVariables())
    if newBindings != nil {
        fmt.Fprintf(logFile, "Merged bindings: %v\n", newBindings.Variables())
    }
}
```

#### Priority 2: Verify Token Propagation

Check each step of the propagation chain:

1. **Join1 (u ⋈ o)**:
   - Input left: token with [u]
   - Input right: fact o → creates token with [o]
   - Output: joined token with [u, o] ← VERIFY THIS

2. **Join2 ((u,o) ⋈ p)**:
   - Input left: token with [u, o] ← VERIFY THIS ARRIVES CORRECTLY
   - Input right: fact p → creates token with [p]
   - Output: joined token with [u, o, p] ← VERIFY THIS

3. **TerminalNode**:
   - Input: token should have [u, o, p] ← CURRENTLY HAS ONLY [u, o]

#### Priority 3: Potential Fixes

Based on analysis, check these potential issues:

1. **Check if performJoinWithTokens() is being called**:
   - Maybe evaluateJoinConditions() is failing and returning nil?
   - Add logging before the nil check

2. **Check if a different code path is being used**:
   - Maybe cascade joins use a different join method?
   - Search for other places where tokens are created

3. **Check if the token is being cloned/copied somewhere**:
   - Maybe the bindings are being lost during a copy operation?
   - Check Token.Clone() or any similar methods

4. **Check connection logic**:
   - Verify Join1 output actually goes to Join2.ActivateLeft()
   - Verify Join2 output goes to TerminalNode.ActivateLeft()

### Test Case for Verification

Create a minimal test that can be run independently:

```go
func TestCascadeBindingsPreservation(t *testing.T) {
    // Create network
    storage := NewMemoryStorage()
    network := NewReteNetwork(storage)
    
    // Define types
    network.AddType("User", map[string]string{"id": "string"})
    network.AddType("Order", map[string]string{"id": "string", "user_id": "string"})
    network.AddType("Product", map[string]string{"id": "string", "name": "string"})
    
    // Create cascade join manually
    varTypes := map[string]string{"u": "User", "o": "Order", "p": "Product"}
    
    // Join1: u ⋈ o
    join1 := NewJoinNode("join1", condition1, []string{"u"}, []string{"o"}, varTypes, storage)
    
    // Join2: (u,o) ⋈ p  
    join2 := NewJoinNode("join2", condition2, []string{"u", "o"}, []string{"p"}, varTypes, storage)
    
    // Connect: join1 → join2
    join1.AddChild(join2)
    
    // Terminal
    terminal := NewTerminalNode("term", action, storage)
    join2.AddChild(terminal)
    
    // Submit facts and verify
    userFact := &Fact{ID: "U1", Type: "User", Fields: map[string]interface{}{"id": "U1"}}
    orderFact := &Fact{ID: "O1", Type: "Order", Fields: map[string]interface{}{"id": "O1", "user_id": "U1"}}
    productFact := &Fact{ID: "P1", Type: "Product", Fields: map[string]interface{}{"id": "P1", "name": "Laptop"}}
    
    // ... submit facts and check terminal token has all 3 bindings
}
```

### Files to Review

- `rete/node_join.go` - performJoinWithTokens(), ActivateLeft(), ActivateRight()
- `rete/binding_chain.go` - Merge() implementation
- `rete/builder_join_rules_cascade.go` - Cascade join construction
- `rete/beta_chain_builder_orchestration.go` - Chain building logic
- `rete/node_alpha.go` - Passthrough alpha nodes
- `rete/node_base.go` - PropagateToChildren()

### Related Documentation

- `docs/architecture/BINDINGS_DESIGN.md` - Design of binding system
- `docs/architecture/BINDINGS_ANALYSIS.md` - Analysis of binding requirements
- `scripts/multi-jointures/10_validation_e2e.md` - Validation requirements

---

**Created**: 2025-12-12
**Priority**: CRITICAL  
**Assignee**: To be determined
**Estimated Effort**: 4-8 hours for proper debugging and fix
