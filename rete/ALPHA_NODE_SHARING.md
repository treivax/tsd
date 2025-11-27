# AlphaNode Sharing Feature

## Overview

The AlphaNode sharing feature is an optimization in the RETE network that allows multiple rules with identical alpha conditions to share the same AlphaNode. This reduces memory usage, improves evaluation performance, and aligns with classical RETE algorithm optimizations.

## Problem Statement

**Before**: Each rule created its own AlphaNode, even when multiple rules had identical conditions.

```
Rule 1: {p: Person} / p.age > 18 ==> print("Adult")
Rule 2: {p: Person} / p.age > 18 ==> print("Can vote")

Previous behavior:
  TypeNode(Person)
    ├── AlphaNode(rule_1_alpha: p.age > 18)
    │   └── TerminalNode(rule_1_terminal)
    └── AlphaNode(rule_2_alpha: p.age > 18)  ← Duplicate!
        └── TerminalNode(rule_2_terminal)
```

**After**: Rules with identical conditions share a single AlphaNode.

```
Current behavior:
  TypeNode(Person)
    └── AlphaNode(alpha_024a66ab: p.age > 18)  ← Shared!
        ├── TerminalNode(rule_1_terminal)
        └── TerminalNode(rule_2_terminal)
```

## Benefits

1. **Memory Efficiency**: One AlphaNode instead of N duplicate nodes
2. **Performance**: Conditions evaluated once, results propagated to all dependent rules
3. **Scalability**: Better performance with large rulesets having overlapping conditions
4. **Maintainability**: Cleaner network structure

## Architecture

### Key Components

#### 1. `AlphaSharingRegistry`
- **Purpose**: Manages shared AlphaNodes across the network
- **Storage**: Maps condition hashes to AlphaNode instances
- **Thread-safe**: Uses `sync.RWMutex` for concurrent access

#### 2. Condition Hashing (`ConditionHash`)
- **Algorithm**: SHA-256 based hash of normalized condition + variable name
- **Normalization**: Ensures consistent hashing for equivalent conditions
- **Uniqueness**: Hash includes both condition structure and variable name

```go
// Two conditions produce the same hash only if:
// 1. All condition attributes are identical (type, attribute, operator, value)
// 2. Variable names are identical (p vs q would have different hashes)
hash1 := ConditionHash(condition, "p")
hash2 := ConditionHash(condition, "p")  // Same hash
hash3 := ConditionHash(condition, "q")  // Different hash
```

#### 3. Lifecycle Integration
- AlphaNodes are tracked by `LifecycleManager`
- Reference counting: Multiple rules can reference the same AlphaNode
- Automatic cleanup: AlphaNode deleted only when last rule is removed

### Data Flow

```
Rule Creation:
1. Parse rule and extract condition
2. Compute condition hash
3. Check AlphaSharingRegistry:
   - If exists: Reuse existing AlphaNode
   - If new: Create AlphaNode, register in registry
4. Create rule-specific TerminalNode
5. Connect TerminalNode as child of AlphaNode
6. Register rule reference in LifecycleManager

Fact Propagation:
1. Fact enters TypeNode
2. TypeNode propagates to shared AlphaNode
3. AlphaNode evaluates condition once
4. If passed: Propagate to ALL child TerminalNodes
5. Each TerminalNode fires its rule-specific action

Rule Removal:
1. Remove rule reference from LifecycleManager
2. Delete rule's TerminalNode
3. If AlphaNode has no more references:
   - Remove from network.AlphaNodes
   - Remove from AlphaSharingRegistry
4. Otherwise: Keep AlphaNode for remaining rules
```

## API Reference

### AlphaSharingRegistry Methods

```go
// Create a new registry
registry := NewAlphaSharingRegistry()

// Get or create an AlphaNode (core sharing logic)
node, hash, wasShared, err := registry.GetOrCreateAlphaNode(
    condition,      // Condition to evaluate
    variableName,   // Variable name (e.g., "p")
    storage,        // Storage instance
)

// Get an existing AlphaNode by hash
node, exists := registry.GetAlphaNode(hash)

// Remove an AlphaNode (called when no rules use it)
err := registry.RemoveAlphaNode(hash)

// Get sharing statistics
stats := registry.GetStats()
// Returns:
// - total_shared_alpha_nodes: Number of unique AlphaNodes
// - total_rule_references: Total references across all nodes
// - average_sharing_ratio: Average rules per AlphaNode

// List all shared AlphaNode hashes
hashes := registry.ListSharedAlphaNodes()

// Get detailed information about a specific AlphaNode
details := registry.GetSharedAlphaNodeDetails(hash)

// Reset the registry (clear all nodes)
registry.Reset()
```

### Network Integration

```go
// AlphaSharingRegistry is automatically created with ReteNetwork
network := NewReteNetwork(storage)
// network.AlphaSharingManager is initialized

// Get network statistics (includes sharing stats)
stats := network.GetNetworkStats()
// Returns keys like:
// - sharing_total_shared_alpha_nodes
// - sharing_total_rule_references
// - sharing_average_sharing_ratio
```

## Examples

### Example 1: Two Rules, Same Condition

```constraint
type Person : <id: string, age: number>

rule adult_check : {p: Person} / p.age > 18 ==> print("Adult")
rule voting_check : {p: Person} / p.age > 18 ==> print("Can vote")
```

**Result**:
- 1 shared AlphaNode (`alpha_024a66ab...`)
- 2 TerminalNodes (one per rule)
- Sharing ratio: 2.0

**Behavior**:
```
Submit: Person{id: "p1", age: 25}
  ↓
TypeNode(Person)
  ↓
AlphaNode (p.age > 18) ← Evaluated ONCE
  ↓
  ├→ TerminalNode(adult_check)  → fires print("Adult")
  └→ TerminalNode(voting_check) → fires print("Can vote")
```

### Example 2: Mixed Conditions

```constraint
type Person : <id: string, age: number>

rule adult_check : {p: Person} / p.age > 18 ==> print("Adult")
rule voting_check : {p: Person} / p.age > 18 ==> print("Can vote")
rule drinking_check : {p: Person} / p.age > 21 ==> print("Can drink")
```

**Result**:
- 2 AlphaNodes (one for `>18`, one for `>21`)
- 3 TerminalNodes
- Sharing ratio: 1.5 (3 rules / 2 nodes)

### Example 3: Different Variables (No Sharing)

```constraint
type Person : <id: string, age: number>

rule check_p : {p: Person} / p.age > 18 ==> print("P is adult")
rule check_q : {q: Person} / q.age > 18 ==> print("Q is adult")
```

**Result**:
- 2 AlphaNodes (variable names differ: `p` vs `q`)
- Conditions look identical but hash differently due to variable names

**Rationale**: Variable names are part of the semantic identity because they define the binding context for the condition evaluation.

## Performance Considerations

### When Sharing Helps Most

1. **Large rulesets** with many rules checking common conditions
2. **High fact throughput** where each fact evaluation is expensive
3. **Overlapping rule patterns** (e.g., age checks, threshold checks)

### Overhead

- **Negligible**: Hash computation is fast (SHA-256 on small JSON)
- **Registry lookup**: O(1) map access
- **Memory**: Minimal overhead (one map entry per unique condition)

### Benchmarking

For 100 rules with 50 unique conditions:
- **Before**: 100 AlphaNodes created
- **After**: 50 AlphaNodes created (50% reduction)
- **Evaluation**: 50% fewer condition evaluations per fact

## Testing

### Unit Tests (`alpha_sharing_feature_test.go`)

- `TestConditionHash`: Verifies hash consistency and uniqueness
- `TestAlphaSharingRegistry_GetOrCreate`: Tests get-or-create logic
- `TestAlphaSharingRegistry_RemoveAlphaNode`: Tests cleanup
- `TestAlphaSharingRegistry_Stats`: Tests statistics
- `TestAlphaSharingRegistry_ConcurrentAccess`: Tests thread safety

### Integration Tests (`alpha_sharing_integration_test.go`)

- `TestAlphaSharingIntegration_TwoRulesSameCondition`: Basic sharing
- `TestAlphaSharingIntegration_ThreeRulesMixedConditions`: Mixed scenarios
- `TestAlphaSharingIntegration_FactPropagation`: Fact propagation correctness
- `TestAlphaSharingIntegration_RuleRemoval`: Lifecycle management
- `TestAlphaSharingIntegration_NetworkReset`: Reset behavior

### Running Tests

```bash
cd tsd/rete

# Run all alpha sharing tests
go test -v -run "TestAlphaSharing"

# Run specific test
go test -v -run "TestAlphaSharingIntegration_TwoRulesSameCondition"

# Run with coverage
go test -cover -run "TestAlphaSharing"
```

## Implementation Details

### Hash-Based Node IDs

Shared AlphaNodes use hash-based IDs instead of rule-based IDs:

```go
// Old approach (no sharing):
alphaNode := NewAlphaNode(ruleID+"_alpha", condition, ...)
// ID: "rule_0_alpha", "rule_1_alpha" (always unique)

// New approach (with sharing):
hash, _ := ConditionHash(condition, variableName)
alphaNode := NewAlphaNode(hash, condition, ...)
// ID: "alpha_024a66abf888c763" (same for identical conditions)
```

### Condition Normalization

Ensures consistent hashing:

```go
// These normalize to the same representation:
cond1 := map[string]interface{}{"type": "constraint", "attr": "age", "op": ">", "val": 18}
cond2 := map[string]interface{}{"attr": "age", "type": "constraint", "val": 18, "op": ">"}
// → Same hash (order doesn't matter after normalization)
```

### Reference Counting

```go
// AlphaNode lifecycle:
Rule A added: RefCount = 1
Rule B added (same condition): RefCount = 2
Rule A removed: RefCount = 1 (AlphaNode kept)
Rule B removed: RefCount = 0 (AlphaNode deleted)
```

## Migration Notes

### Backward Compatibility

- **Fully backward compatible**: Existing code works without changes
- **Automatic**: Sharing happens transparently during rule creation
- **No API changes**: `NewReteNetwork()` and rule building work as before

### Observing Sharing

Enable debug output or check network stats:

```go
stats := network.GetNetworkStats()
fmt.Printf("Shared AlphaNodes: %d\n", stats["sharing_total_shared_alpha_nodes"])
fmt.Printf("Rule references: %d\n", stats["sharing_total_rule_references"])
fmt.Printf("Sharing ratio: %.2f\n", stats["sharing_average_sharing_ratio"])
```

### Log Messages

When building rules, you'll see:

```
✨ Nouveau AlphaNode partageable créé: alpha_024a66ab (hash: alpha_024a66ab)
♻️  AlphaNode partagé réutilisé: alpha_024a66ab (hash: alpha_024a66ab)
✓ Règle rule_1 attachée à l'AlphaNode partagé alpha_024a66ab via terminal rule_1_terminal
```

## Future Enhancements

### Potential Improvements

1. **Beta Node Sharing**: Extend sharing to join/beta nodes
2. **Condition Subsumption**: Share nodes when one condition subsumes another
3. **Metrics/Monitoring**: Detailed analytics on sharing effectiveness
4. **Persistence**: Save/restore shared node registry
5. **Query API**: Find which rules share which conditions

### Limitations

- **Current scope**: Only simple alpha conditions (not join/beta nodes yet)
- **Hash collisions**: Extremely unlikely with SHA-256, but theoretically possible
- **Variable names**: Different variable names prevent sharing (by design)

## Troubleshooting

### AlphaNode Not Shared When Expected

**Check**:
1. Are variable names identical? (`p` vs `q` won't share)
2. Are conditions exactly identical? (operator, attribute, value)
3. Check logs for "Nouveau AlphaNode" vs "AlphaNode partagé réutilisé"

### Memory Not Freed After Rule Removal

**Check**:
1. Verify `RemoveRule()` was called
2. Check if other rules still reference the AlphaNode
3. Inspect lifecycle stats: `network.LifecycleManager.GetStats()`

### Debugging Tips

```go
// List all shared AlphaNodes
hashes := network.AlphaSharingManager.ListSharedAlphaNodes()
for _, hash := range hashes {
    details := network.AlphaSharingManager.GetSharedAlphaNodeDetails(hash)
    fmt.Printf("AlphaNode %s: %d child(ren)\n", hash, details["child_count"])
}

// Check specific AlphaNode
node, exists := network.AlphaNodes["alpha_024a66ab..."]
if exists {
    fmt.Printf("Children: %d\n", len(node.GetChildren()))
}
```

## Related Documentation

- `NODE_LIFECYCLE_FEATURE.md`: Lifecycle management and reference counting
- `TYPENODE_SHARING_REPORT.md`: TypeNode sharing behavior
- `ALPHA_NODE_SHARING_REPORT.md`: Original investigation and design decisions

## References

- Classic RETE Algorithm: Forgy, C. L. (1982). "Rete: A Fast Algorithm for the Many Pattern/Many Object Pattern Match Problem"
- Node sharing is a standard optimization in production rule systems
- Similar implementations: Drools, Jess, CLIPS

---

**Last Updated**: 2025-01-XX  
**Version**: 1.0  
**Status**: Implemented and tested