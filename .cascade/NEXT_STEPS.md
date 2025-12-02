# 🚀 Next Steps - Arithmetic Decomposition Implementation

## Current Status

✅ **Phase 1: Core Infrastructure - COMPLETED**

**What's Ready**:
- `EvaluationContext` - Thread-safe intermediate results storage
- `AlphaNode.ActivateWithContext()` - Context-aware activation
- `ConditionEvaluator` - Resolves `tempResult` references
- 23 tests passing, 100% success rate

**What's Next**: Phase 2 - Integration with Decomposer

---

## 🎯 Phase 2: Integration (Next Sprint)

### Task 2.1: Update ArithmeticExpressionDecomposer

**File**: `rete/arithmetic_decomposer.go`

**Objective**: Modify the decomposer to generate `tempResult` references instead of embedded expressions.

**Current Behavior** (Monolithic):
```go
// Returns full nested expression
return map[string]interface{}{
    "type": "binaryOp",
    "operator": "-",
    "left": map[string]interface{}{
        "type": "binaryOp",
        "operator": "*",
        "left": fieldAccessQte,
        "right": numberLiteral23,
    },
    "right": numberLiteral10,
}
```

**Desired Behavior** (Decomposed):
```go
// Step 1: qte * 23 → temp_1
steps = append(steps, SimpleCondition{
    Condition: multiplyCondition,
    ResultName: "temp_1",
})

// Step 2: temp_1 - 10 → temp_2 (references temp_1)
steps = append(steps, SimpleCondition{
    Condition: map[string]interface{}{
        "type": "binaryOp",
        "operator": "-",
        "left": map[string]interface{}{
            "type": "tempResult",
            "step_name": "temp_1",  // ← Reference to previous result
        },
        "right": numberLiteral10,
    },
    ResultName: "temp_2",
})

return map[string]interface{}{
    "type": "tempResult",
    "step_name": "temp_2",  // ← Final result reference
}
```

**Implementation Steps**:

1. Modify `decomposeBinaryOp()`:
```go
func (aed *ArithmeticExpressionDecomposer) decomposeBinaryOp(
    expr map[string]interface{},
    steps *[]SimpleCondition,
    stepCounter *int,
) map[string]interface{} {
    // Recursively decompose left and right
    leftRef := aed.decomposeExpression(expr["left"], steps, stepCounter)
    rightRef := aed.decomposeExpression(expr["right"], steps, stepCounter)
    
    // Create result name for this step
    *stepCounter++
    resultName := fmt.Sprintf("temp_%d", *stepCounter)
    
    // Create atomic condition
    atomicCondition := map[string]interface{}{
        "type": "binaryOp",
        "operator": expr["operator"],
        "left": leftRef,   // Can be tempResult or direct value
        "right": rightRef,
    }
    
    // Store step with result name
    *steps = append(*steps, SimpleCondition{
        Condition: atomicCondition,
        ResultName: resultName,
        Dependencies: extractDepsFrom(leftRef, rightRef),
    })
    
    // Return reference to this result
    return map[string]interface{}{
        "type": "tempResult",
        "step_name": resultName,
    }
}
```

2. Add `extractDepsFrom()` helper:
```go
func extractDepsFrom(left, right interface{}) []string {
    deps := []string{}
    
    if leftMap, ok := left.(map[string]interface{}); ok {
        if leftMap["type"] == "tempResult" {
            if name, ok := leftMap["step_name"].(string); ok {
                deps = append(deps, name)
            }
        }
    }
    
    // Same for right...
    
    return deps
}
```

3. Update `SimpleCondition`:
```go
type SimpleCondition struct {
    Condition    interface{}
    ResultName   string    // NEW: name of result produced
    Dependencies []string  // NEW: required temp results
}
```

**Testing**:
```go
func TestArithmeticDecomposer_GeneratesTempResults(t *testing.T) {
    decomposer := NewArithmeticExpressionDecomposer()
    
    expr := buildExpression("(c.qte * 23 - 10) > 0")
    steps := decomposer.Decompose(expr)
    
    // Verify step 1 produces temp_1
    assert.Equal(t, "temp_1", steps[0].ResultName)
    
    // Verify step 2 references temp_1
    step2Cond := steps[1].Condition.(map[string]interface{})
    leftRef := step2Cond["left"].(map[string]interface{})
    assert.Equal(t, "tempResult", leftRef["type"])
    assert.Equal(t, "temp_1", leftRef["step_name"])
    
    // Verify dependencies
    assert.Contains(t, steps[1].Dependencies, "temp_1")
}
```

---

### Task 2.2: Enhance AlphaChainBuilder

**File**: `rete/alpha_chain_builder.go`

**Objective**: Set decomposition metadata on AlphaNodes during chain building.

**Modifications**:

1. Update `BuildChain()`:
```go
func (acb *AlphaChainBuilder) BuildChain(
    conditions []SimpleCondition,
    varName string,
) (*AlphaNode, error) {
    var root *AlphaNode
    var prev *AlphaNode
    
    for _, cond := range conditions {
        nodeID := acb.generateNodeID(cond)
        
        // Get or create node with sharing
        node := acb.sharingRegistry.GetOrCreateAlphaNode(
            cond.Condition,
            func() *AlphaNode {
                alphaNode := NewAlphaNode(nodeID, cond.Condition, varName, acb.storage)
                
                // SET DECOMPOSITION METADATA
                alphaNode.ResultName = cond.ResultName     // ← NEW
                alphaNode.IsAtomic = true                   // ← NEW
                alphaNode.Dependencies = cond.Dependencies  // ← NEW
                
                return alphaNode
            },
        )
        
        // Link in chain...
    }
    
    return root, nil
}
```

**Testing**:
```go
func TestAlphaChainBuilder_SetsDecompositionMetadata(t *testing.T) {
    builder := NewAlphaChainBuilder(storage, registry)
    
    conditions := []SimpleCondition{
        {
            Condition: multiplyCondition,
            ResultName: "temp_1",
            Dependencies: []string{},
        },
        {
            Condition: subtractCondition,
            ResultName: "temp_2",
            Dependencies: []string{"temp_1"},
        },
    }
    
    chain, _ := builder.BuildChain(conditions, "c")
    
    // Verify first node
    assert.Equal(t, "temp_1", chain.ResultName)
    assert.True(t, chain.IsAtomic)
    assert.Empty(t, chain.Dependencies)
    
    // Verify second node
    secondNode := chain.Children[0].(*AlphaNode)
    assert.Equal(t, "temp_2", secondNode.ResultName)
    assert.True(t, secondNode.IsAtomic)
    assert.Contains(t, secondNode.Dependencies, "temp_1")
}
```

---

### Task 2.3: Update JoinRuleBuilder Integration

**File**: `rete/builder_join_rules.go`

**Objective**: Use decomposition and create EvaluationContext during fact activation.

**Modifications**:

1. Add configuration:
```go
type DecompositionConfig struct {
    Enabled       bool
    MinComplexity int  // Min operations to trigger decomposition
}

// In JoinRuleBuilder
type JoinRuleBuilder struct {
    // ... existing fields ...
    decompositionConfig DecompositionConfig
}
```

2. Modify `createBinaryJoinRule()`:
```go
func (b *JoinRuleBuilder) createBinaryJoinRule(...) error {
    // ... existing alpha condition splitting ...
    
    // For each alpha condition
    for varName, alphaConds := range alphaCondsByVar {
        // Decide: decompose or monolithic?
        useDecomposition := b.shouldDecompose(alphaConds)
        
        var alphaChain *AlphaNode
        
        if useDecomposition {
            // Decompose into steps
            decomposer := NewArithmeticExpressionDecomposer()
            steps := decomposer.Decompose(alphaConds)
            
            // Build chain with decomposition metadata
            chainBuilder := NewAlphaChainBuilder(b.storage, b.sharingRegistry)
            alphaChain, _ = chainBuilder.BuildChain(steps, varName)
        } else {
            // Monolithic approach (existing code)
            alphaChain = NewAlphaNode(...)
        }
        
        // Link to TypeNode
        typeNode.AddChild(alphaChain)
    }
    
    return nil
}
```

3. Modify TypeNode activation to use context:
```go
// In TypeNode.ActivateRight()
func (tn *TypeNode) ActivateRight(fact *Fact) error {
    // ... existing code ...
    
    for _, child := range tn.GetChildren() {
        if alphaNode, ok := child.(*AlphaNode); ok {
            // If child uses decomposition, create context
            if alphaNode.IsAtomic || len(alphaNode.Dependencies) > 0 {
                ctx := NewEvaluationContext(fact)
                err = alphaNode.ActivateWithContext(fact, ctx)
            } else {
                // Standard activation
                err = alphaNode.ActivateRight(fact)
            }
        }
        // ... error handling ...
    }
    
    return nil
}
```

**Testing**:
```go
func TestJoinRuleBuilder_WithDecomposition(t *testing.T) {
    builder := NewJoinRuleBuilder(storage)
    builder.decompositionConfig.Enabled = true
    
    rule := &Rule{
        Conditions: []string{
            "c.type == 'commande'",
            "(c.qte * 23 - 10) > 0",  // ← Will be decomposed
        },
        Action: "approve_order",
    }
    
    err := builder.Build(rule)
    assert.NoError(t, err)
    
    // Submit test fact
    fact := &Fact{
        Type: "commande",
        Fields: map[string]interface{}{
            "qte": 10,
        },
    }
    
    network.SubmitFact(fact)
    
    // Verify token produced (decomposition worked)
    tokens := getTokensInTerminal("approve_order")
    assert.Len(t, tokens, 1)
}
```

---

## 📅 Estimated Timeline

| Task | Effort | Duration |
|------|--------|----------|
| 2.1 - ArithmeticExpressionDecomposer | Medium | 4-6 hours |
| 2.2 - AlphaChainBuilder | Small | 2-3 hours |
| 2.3 - JoinRuleBuilder Integration | Medium | 4-6 hours |
| **Total Phase 2** | | **1-2 days** |

---

## 🧪 Testing Strategy for Phase 2

### Unit Tests
- [ ] Decomposer generates tempResult references
- [ ] AlphaChainBuilder sets metadata correctly
- [ ] Dependencies extracted properly

### Integration Tests
- [ ] End-to-end: complex expression → decomposed chain → correct result
- [ ] Sharing: common sub-expressions shared across rules
- [ ] Mixed: decomposed and monolithic chains coexist

### E2E Tests
- [ ] Extend `TestArithmeticExpressionsE2E` with decomposition mode
- [ ] Compare results: decomposed vs monolithic (must be identical)
- [ ] Performance: measure overhead

---

## 📖 Reference Documents

1. **Spec**: `rete/ARITHMETIC_DECOMPOSITION_SPEC.md`
2. **Prompt**: `.cascade/add-feature-arithmetic-decomposition.md`
3. **Progress**: `rete/ARITHMETIC_DECOMPOSITION_IMPLEMENTATION_PROGRESS.md`
4. **Summary**: `rete/ARITHMETIC_DECOMPOSITION_SUMMARY.md`
5. **Phase 1 Report**: `.cascade/PHASE1_COMPLETION_REPORT.md`

---

## 🚦 Ready to Start?

**Prerequisites**: ✅ All met
- Phase 1 complete and tested
- Documentation up to date
- No blocking issues

**Next Command**:
```bash
# Read the arithmetic decomposer
cat rete/arithmetic_decomposer.go

# Start implementing Task 2.1
# Follow the implementation guide above
```

---

**Good luck with Phase 2!** 🚀

*Last updated: 2025-01-XX*
