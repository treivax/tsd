# AlphaNode Sharing - Quick Reference Guide

## What is AlphaNode Sharing?

Multiple rules with identical conditions now **share a single AlphaNode**, reducing memory usage and improving performance.

## Quick Example

```constraint
type Person : <id: string, age: number>

rule adult_check : {p: Person} / p.age > 18 ==> print("Adult")
rule voting_check : {p: Person} / p.age > 18 ==> print("Can vote")
```

**Before**: 2 AlphaNodes (one per rule)  
**After**: 1 shared AlphaNode + 2 TerminalNodes

```
TypeNode(Person)
  └── AlphaNode(alpha_024a66ab: p.age > 18)  ← Shared!
      ├── TerminalNode(adult_check)
      └── TerminalNode(voting_check)
```

---

## Key Concepts

### Condition Hash
- Each condition gets a unique SHA-256 hash
- Hash includes: condition structure + variable name
- Identical conditions = identical hash = same AlphaNode

### Sharing Criteria
Two conditions share an AlphaNode if:
- ✅ Same attribute (`age` vs `age`)
- ✅ Same operator (`>` vs `>`)
- ✅ Same value (`18` vs `18`)
- ✅ Same variable name (`p` vs `p`)

### No Sharing When:
- ❌ Different variable names (`p` vs `q`)
- ❌ Different operators (`>` vs `>=`)
- ❌ Different values (`18` vs `21`)
- ❌ Different attributes (`age` vs `salary`)

---

## Usage

### Automatic (No Code Changes Required)

```go
// Just use the network as before
storage := NewMemoryStorage()
pipeline := NewConstraintPipeline()
network, err := pipeline.BuildNetworkFromConstraintFile("rules.constraint", storage)
// Sharing happens automatically!
```

### Check Sharing Statistics

```go
stats := network.GetNetworkStats()
fmt.Printf("Shared AlphaNodes: %d\n", 
    stats["sharing_total_shared_alpha_nodes"])
fmt.Printf("Rule references: %d\n", 
    stats["sharing_total_rule_references"])
fmt.Printf("Sharing ratio: %.2f\n", 
    stats["sharing_average_sharing_ratio"])
```

### Inspect Shared Nodes

```go
// List all shared AlphaNodes
hashes := network.AlphaSharingManager.ListSharedAlphaNodes()

// Get details for a specific node
for _, hash := range hashes {
    details := network.AlphaSharingManager.GetSharedAlphaNodeDetails(hash)
    fmt.Printf("AlphaNode %s: %d rules\n", 
        details["node_id"], 
        details["child_count"])
}
```

---

## Log Messages

### When building rules:

**New AlphaNode created:**
```
✨ Nouveau AlphaNode partageable créé: alpha_024a66ab (hash: alpha_024a66ab)
```

**Existing AlphaNode reused:**
```
♻️  AlphaNode partagé réutilisé: alpha_024a66ab (hash: alpha_024a66ab)
✓ Règle rule_1 attachée à l'AlphaNode partagé alpha_024a66ab via terminal rule_1_terminal
```

---

## Common Patterns

### Pattern 1: Age Thresholds
```constraint
rule adult : {p: Person} / p.age > 18 ==> action1()
rule voter : {p: Person} / p.age > 18 ==> action2()
rule senior : {p: Person} / p.age > 65 ==> action3()
```
**Result**: 2 AlphaNodes (one for `>18`, one for `>65`)

### Pattern 2: Salary Ranges
```constraint
rule high_earner1 : {p: Person} / p.salary > 100000 ==> action1()
rule high_earner2 : {p: Person} / p.salary > 100000 ==> action2()
rule mid_earner : {p: Person} / p.salary > 50000 ==> action3()
```
**Result**: 2 AlphaNodes (one for `>100000`, one for `>50000`)

### Pattern 3: Different Variables (No Sharing)
```constraint
rule check_p : {p: Person} / p.age > 18 ==> action1()
rule check_q : {q: Person} / q.age > 18 ==> action2()
```
**Result**: 2 AlphaNodes (variables `p` and `q` are different)

---

## Performance Benefits

| Scenario | Before | After | Improvement |
|----------|--------|-------|-------------|
| 100 rules, 50 unique conditions | 100 AlphaNodes | 50 AlphaNodes | 50% reduction |
| Evaluations per fact | 100 | 50 | 2x faster |
| Memory usage | High | Low | 50% less |

---

## Testing

### Run tests:
```bash
cd tsd/rete

# All alpha sharing tests
go test -v -run "TestAlphaSharing"

# Specific test
go test -v -run "TestAlphaSharingIntegration_TwoRulesSameCondition"

# All RETE tests
go test -v
```

---

## API Reference (Quick)

### AlphaSharingRegistry Methods

```go
// Get or create an AlphaNode
node, hash, wasShared, err := registry.GetOrCreateAlphaNode(
    condition, variableName, storage)

// Get statistics
stats := registry.GetStats()

// List all shared nodes
hashes := registry.ListSharedAlphaNodes()

// Get node details
details := registry.GetSharedAlphaNodeDetails(hash)

// Reset registry
registry.Reset()
```

### Network Integration

```go
// Access the registry
network.AlphaSharingManager

// Get network stats (includes sharing stats)
stats := network.GetNetworkStats()
// Keys: sharing_total_shared_alpha_nodes, 
//       sharing_total_rule_references,
//       sharing_average_sharing_ratio
```

---

## Troubleshooting

### AlphaNode not shared when expected?

1. **Check variable names**: `p` vs `q` won't share
2. **Check conditions exactly**: Even spaces in values matter
3. **Enable debug logs**: Look for "Nouveau" vs "réutilisé"

### Debug commands:
```go
// List all AlphaNodes
for id, node := range network.AlphaNodes {
    fmt.Printf("%s: %d children\n", id, len(node.GetChildren()))
}

// Check lifecycle
lifecycle, _ := network.LifecycleManager.GetNodeLifecycle("alpha_024a66ab...")
fmt.Printf("References: %d\n", lifecycle.GetRefCount())
```

---

## Important Notes

✅ **Fully backward compatible** - existing code works without changes  
✅ **Thread-safe** - safe for concurrent access  
✅ **Automatic cleanup** - nodes deleted when no rules reference them  
✅ **Zero overhead** - hash lookup is O(1)  

❌ **Current limitation**: Only simple alpha conditions (not join/beta yet)  
❌ **Variable names matter**: Different names = different AlphaNodes  

---

## Examples with Output

### Example 1: Basic Sharing
```constraint
type Person : <id: string, age: number>
rule r1 : {p: Person} / p.age > 18 ==> print("A")
rule r2 : {p: Person} / p.age > 18 ==> print("B")
```

**Statistics:**
```
alpha_nodes: 1
terminal_nodes: 2
sharing_total_shared_alpha_nodes: 1
sharing_total_rule_references: 2
sharing_average_sharing_ratio: 2.0
```

### Example 2: Partial Sharing
```constraint
type Person : <id: string, age: number>
rule r1 : {p: Person} / p.age > 18 ==> print("A")
rule r2 : {p: Person} / p.age > 18 ==> print("B")
rule r3 : {p: Person} / p.age > 21 ==> print("C")
```

**Statistics:**
```
alpha_nodes: 2
terminal_nodes: 3
sharing_total_shared_alpha_nodes: 2
sharing_total_rule_references: 3
sharing_average_sharing_ratio: 1.5
```

---

## Related Documentation

- **Full Guide**: `ALPHA_NODE_SHARING.md`
- **Implementation**: `ALPHA_SHARING_IMPLEMENTATION_SUMMARY.md`
- **Lifecycle**: `NODE_LIFECYCLE_FEATURE.md`
- **TypeNode Sharing**: `TYPENODE_SHARING_REPORT.md`

---

**Version**: 1.0  
**Status**: Production Ready ✅  
**Last Updated**: January 2025