# Remove Rule Command - Documentation

**Feature:** Dynamic Rule Removal  
**Status:** ✅ Production Ready  
**Added:** January 2025  
**License:** MIT

---

## 📋 Overview

The `remove rule <ID>` command allows dynamic removal of rules from the RETE network at runtime. This feature enables incremental rule management without requiring a complete network rebuild.

### Key Features

- ✅ **Grammar Support** - Native PEG grammar support for `remove rule <ruleID>`
- ✅ **Pipeline Integration** - Full integration with constraint pipeline
- ✅ **Lifecycle Management** - Safe removal with reference counting
- ✅ **Node Cleanup** - Automatic cleanup of unused nodes
- ✅ **Alpha Rules** - Full support for alpha chain rules
- ⚠️ **Beta Rules** - Partial support (join rules - lifecycle integration pending)

---

## 🎯 Syntax

```tsd
remove rule <ruleID>
```

### Examples

```tsd
# Remove a single rule
remove rule adult_check

# Remove multiple rules
remove rule rule1
remove rule rule2
remove rule rule3
```

---

## 📖 Usage

### Basic Example

```tsd
# Define types and rules
type Person : <id: string, name: string, age: number>

rule adult_check : {p: Person} / p.age >= 18 ==> adult(p.id)
rule senior_check : {p: Person} / p.age >= 65 ==> senior(p.id)
rule minor_check : {p: Person} / p.age < 18 ==> minor(p.id)

# Add facts
Person(id:p1, name:Alice, age:30)
Person(id:p2, name:Bob, age:70)
Person(id:p3, name:Charlie, age:15)

# Remove a rule
remove rule senior_check
```

**Result:**
- `adult_check` and `minor_check` remain active
- `senior_check` is removed from the network
- All nodes associated only with `senior_check` are cleaned up
- Shared nodes (if any) are preserved

### Incremental Updates

```tsd
# Initial rules
rule r1 : {p: Person} / p.age > 10 ==> action1(p.id)
rule r2 : {p: Person} / p.age > 20 ==> action2(p.id)
rule r3 : {p: Person} / p.age > 30 ==> action3(p.id)

Person(id:p1, age:55)

# Later, remove some rules
remove rule r1
remove rule r3

# Only r2 remains active
```

---

## 🏗️ Architecture

### Grammar Definition

File: `constraint/grammar/constraint.peg`

```peg
RemoveRule <- "remove" _ "rule" _ ruleID:IdentName {
    return map[string]interface{}{
        "type": "ruleRemoval",
        "ruleID": ruleID,
    }, nil
}
```

### Pipeline Processing

File: `rete/constraint_pipeline.go`

The pipeline processes rule removals in the following order:

1. **Parse** - Extract `ruleRemoval` entries from AST
2. **Build Network** - Construct RETE network from rules
3. **Process Removals** - Call `network.RemoveRule(ruleID)` for each removal
4. **Cleanup** - Automatic node cleanup via LifecycleManager

```go
func (cp *ConstraintPipeline) processRuleRemovals(
    network *ReteNetwork, 
    resultMap map[string]interface{},
) error {
    ruleRemovalsData := resultMap["ruleRemovals"]
    ruleRemovals := ruleRemovalsData.([]interface{})
    
    for _, removalData := range ruleRemovals {
        ruleID := removalMap["ruleID"].(string)
        err := network.RemoveRule(ruleID)
        // Handle errors...
    }
    return nil
}
```

### Network Removal

File: `rete/network.go`

```go
func (rn *ReteNetwork) RemoveRule(ruleID string) error {
    // 1. Get nodes associated with this rule
    nodeIDs := rn.LifecycleManager.GetNodesForRule(ruleID)
    
    // 2. Determine removal strategy (alpha vs beta)
    if hasJoinNodes {
        return rn.removeJoinRule(ruleID, nodeIDs)  // Beta removal
    } else {
        return rn.removeAlphaChain(ruleID)         // Alpha removal
    }
    
    // 3. Remove nodes with reference counting
    // 4. Disconnect from parent nodes
    // 5. Clean up network collections
}
```

### Lifecycle Management

The LifecycleManager tracks:
- **Node → Rules mapping** - Which rules use which nodes
- **Rule → Nodes mapping** - Which nodes belong to which rules
- **Reference counting** - How many rules share each node

**Key Operations:**

```go
// Add tracking when creating nodes
lifecycleManager.RegisterNode(nodeID, nodeType, ruleID)

// Check references before removal
shouldDelete := lifecycleManager.RemoveRuleFromNode(nodeID, ruleID)

// Get all nodes for a rule
nodeIDs := lifecycleManager.GetNodesForRule(ruleID)
```

---

## 🔧 Implementation Details

### Node Removal Order

For alpha chains (simple rules):

1. **Terminal Node** - Remove rule's terminal node first
2. **Alpha Nodes** - Remove alpha nodes in reverse chain order
3. **Type Node** - Never removed (shared across all rules of that type)

```
TypeNode (Person)
    ↓
AlphaNode (age > 18)     ← Remove if no other rules use it
    ↓
AlphaNode (name == X)    ← Remove if no other rules use it
    ↓
TerminalNode (rule_X)    ← Always removed with rule
```

### Reference Counting

Each node tracks how many rules reference it:

```go
type NodeLifecycle struct {
    NodeID      string
    NodeType    string
    Rules       map[string]bool  // Set of rule IDs
    RefCount    int              // len(Rules)
    CreatedAt   time.Time
}

// Only delete when RefCount reaches 0
func (lm *LifecycleManager) RemoveRuleFromNode(nodeID, ruleID string) (shouldDelete bool, err error) {
    lifecycle := lm.nodes[nodeID]
    delete(lifecycle.Rules, ruleID)
    lifecycle.RefCount = len(lifecycle.Rules)
    
    return lifecycle.RefCount == 0, nil
}
```

### Node Sharing

Nodes are shared when rules have identical patterns:

```tsd
# These two rules share the AlphaNode (age > 18)
rule r1 : {p: Person} / p.age > 18 ==> action1(p.id)
rule r2 : {p: Person} / p.age > 18 AND p.name == "Alice" ==> action2(p.id)

# Removing r1 keeps the shared node because r2 still needs it
remove rule r1

# Removing r2 deletes the shared node (RefCount = 0)
remove rule r2
```

---

## ✅ Validation & Testing

### Test Coverage

File: `rete/remove_rule_incremental_test.go`

**4 comprehensive tests:**

1. **TestRemoveRuleIncremental_FullPipeline** ✅
   - End-to-end pipeline test
   - Multiple incremental removals
   - Network structure validation

2. **TestRemoveRuleIncremental_WithJoins** ⚠️
   - Join rule removal (skipped - pending beta lifecycle)
   - Beta node cleanup

3. **TestRemoveRuleIncremental_MultipleRemovals** ✅
   - Batch removal testing
   - 5 rules → remove 4 → verify 1 remains

4. **TestRemoveRuleIncremental_ParseOnly** ✅
   - Grammar parsing validation
   - Various rule ID formats

### Running Tests

```bash
# Run all remove rule tests
go test -v ./rete -run TestRemoveRuleIncremental

# Run specific test
go test -v ./rete -run TestRemoveRuleIncremental_FullPipeline

# Run with coverage
go test -coverprofile=coverage.out ./rete -run TestRemoveRuleIncremental
go tool cover -html=coverage.out
```

### Example Test Output

```
🧪 TEST REMOVE RULE - PIPELINE COMPLET INCRÉMENTAL
====================================================

📝 ÉTAPE 1: Création du fichier initial avec 3 règles
✅ Fichier créé avec 3 règles et 3 faits

🔧 ÉTAPE 2: Construction du réseau RETE initial
✅ Réseau construit avec 3 règles

🗑️  ÉTAPE 3: Ajout de la commande 'remove rule senior_check'
✅ Commande de suppression ajoutée au fichier

🔄 ÉTAPE 4: Reconstruction du réseau avec suppression
✅ Règles restantes: 2
✅ La règle 'senior_check' a été correctement supprimée
✅ Les règles 'adult_check' et 'minor_check' sont préservées

✅ TEST COMPLET - Pipeline incrémental validé avec succès!
```

---

## 🚀 Performance Considerations

### Removal Complexity

| Operation | Complexity | Notes |
|-----------|------------|-------|
| Parse removal command | O(1) | Simple AST extraction |
| Find nodes for rule | O(1) | HashMap lookup |
| Remove alpha chain | O(n) | n = chain depth |
| Reference count update | O(1) | HashMap operations |
| Node disconnection | O(m) | m = parent children count |

### Memory Impact

- **Immediate:** Terminal node and rule metadata removed
- **Deferred:** Alpha nodes removed when RefCount = 0
- **Preserved:** Type nodes (never removed)
- **Savings:** Proportional to uniqueness of rule's patterns

### Best Practices

1. **Batch Removals** - Remove multiple rules in one file update
2. **Order Independence** - Removals can be in any order
3. **Network Rebuild** - For many removals, consider full rebuild
4. **Monitoring** - Track node counts before/after removal

---

## ⚠️ Limitations & Future Work

### Current Limitations

1. **Beta Rule Removal** ⚠️
   - Join rules (multi-pattern) removal is incomplete
   - Join node lifecycle tracking needs enhancement
   - See: `beta_backward_compatibility_test.go` (skipped tests)

2. **No Rollback**
   - Rule removal is permanent within a session
   - No undo mechanism (rebuild network to restore)

3. **Network Reconstruction**
   - Each pipeline run rebuilds the entire network
   - Removals are applied after build, not incrementally during

### Future Enhancements

- [ ] **Live Removal** - Remove rules without network rebuild
- [ ] **Join Rule Removal** - Complete beta node lifecycle
- [ ] **Removal Callbacks** - Notify on successful removal
- [ ] **Dry-Run Mode** - Preview what would be removed
- [ ] **Removal Metrics** - Track removal statistics
- [ ] **Partial Removal** - Remove specific patterns within rules

---

## 🔍 Debugging

### Enable Verbose Logging

```go
// During network construction
network, err := pipeline.BuildNetworkFromConstraintFile(file, storage)

// Check removal execution
fmt.Printf("Removing rule: %s\n", ruleID)
err := network.RemoveRule(ruleID)
if err != nil {
    fmt.Printf("Removal failed: %v\n", err)
}

// Verify network state
stats := network.GetNetworkStats()
fmt.Printf("Terminal nodes: %d\n", stats["terminal_nodes"])
fmt.Printf("Alpha nodes: %d\n", stats["alpha_nodes"])
```

### Common Issues

**Issue:** Rule not found
```
Error: règle xyz non trouvée ou aucun nœud associé
```
**Solution:** Verify rule ID matches exactly (case-sensitive)

**Issue:** Node still referenced
```
Error: impossible de supprimer le nœud: encore N référence(s)
```
**Solution:** Other rules are sharing this node (expected behavior)

**Issue:** Join rule removal fails
```
Error: Beta rule removal not fully implemented
```
**Solution:** This is a known limitation - see Future Work

---

## 📚 Related Documentation

- **Lifecycle Management:** `RULE_REMOVAL_WITH_JOINS_FEATURE.md`
- **Alpha Sharing:** `docs/ALPHA_CHAINS_README.md`
- **Beta Sharing:** `BETA_SHARING_SYSTEM.md`
- **Grammar:** `constraint/grammar/constraint.peg`
- **Pipeline:** `rete/constraint_pipeline.go`

---

## 📝 License

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License

---

## ✅ Checklist

Implementation checklist:

- [x] Grammar support for `remove rule <ID>`
- [x] Parser integration
- [x] Pipeline processing
- [x] Alpha rule removal (full support)
- [x] Node lifecycle management
- [x] Reference counting
- [x] Comprehensive tests (4 tests)
- [x] Documentation
- [ ] Beta rule removal (partial - pending)
- [ ] Live removal (future)
- [ ] Removal callbacks (future)

---

**Status:** ✅ PRODUCTION READY (for alpha rules)  
**Version:** 1.0.0  
**Last Updated:** January 2025