# Beta Sharing Design (JoinNode Sharing Registry)

**Version:** 1.0  
**Status:** Design Draft  
**Author:** TSD RETE Team  
**Date:** 2024

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Background & Motivation](#background--motivation)
3. [Architecture Overview](#architecture-overview)
4. [BetaSharingRegistry Design](#betasharingregistry-design)
5. [Sharing Criteria & Compatibility](#sharing-criteria--compatibility)
6. [Normalization & Hashing](#normalization--hashing)
7. [Public API](#public-api)
8. [Integration with Builder & Lifecycle](#integration-with-builder--lifecycle)
9. [Sequence Diagrams](#sequence-diagrams)
10. [Usage Examples](#usage-examples)
11. [Performance Considerations](#performance-considerations)
12. [Testing Strategy](#testing-strategy)
13. [Migration & Rollout](#migration--rollout)
14. [Future Enhancements](#future-enhancements)

---

## Executive Summary

This document describes the design of the **BetaSharingRegistry**, a system for sharing JoinNodes (Beta nodes) in the TSD RETE engine. Similar to the existing AlphaSharingRegistry, the BetaSharingRegistry will:

- **Eliminate duplicate JoinNodes** by identifying and reusing nodes with identical join patterns
- **Reduce memory consumption** by 30-50% in rule bases with overlapping join conditions
- **Improve runtime performance** by 20-40% through reduced node activation overhead
- **Maintain backward compatibility** with existing rule compilation and execution

The system follows these principles:
- **Performance over memory**: Caching and fast lookups prioritized
- **Thread-safe**: All operations protected by sync.RWMutex
- **Inspired by AlphaSharingRegistry**: Proven patterns reused
- **Incremental adoption**: Feature-flagged for safe rollout

---

## Background & Motivation

### Current State

The RETE network currently implements **AlphaNode sharing** through the AlphaSharingRegistry. However, **JoinNodes (Beta nodes) are never shared**, even when multiple rules contain identical join patterns.

For example, these two rules:
```/dev/null/example.tsd#L1-10
rule "HighValueOrder" {
    when {
        order: Order(order.value > 1000)
        customer: Customer(customer.id == order.customerId)
    }
    then { ... }
}

rule "PremiumCustomerOrder" {
    when {
        order: Order(order.value > 500)
        customer: Customer(customer.id == order.customerId)  // Same join!
    }
    then { ... }
}
```

Both rules contain the join condition `customer.id == order.customerId`, but the current implementation creates **two separate JoinNodes**, each with its own:
- Left memory (for order tokens)
- Right memory (for customer tokens)
- Result memory (for joined tokens)
- Join evaluation logic

### Problem Statement

**Without sharing:**
- Memory usage scales linearly with rule count, even for duplicate join patterns
- Join evaluation is performed redundantly across identical JoinNodes
- Network size grows unnecessarily, impacting compilation and debugging
- Common patterns (FK joins, temporal joins) are repeated extensively

### Expected Benefits

**With Beta Sharing:**
- **Memory reduction**: 30-50% in typical rule bases (based on Alpha sharing results)
- **Performance improvement**: 20-40% reduction in join evaluation overhead
- **Smaller network graphs**: Easier visualization and debugging
- **Faster compilation**: Fewer nodes to create and connect

### Design Goals

1. **Correctness**: Shared nodes must produce identical results to independent nodes
2. **Performance**: Hashing and lookup must be fast (sub-millisecond)
3. **Safety**: Thread-safe operations with minimal lock contention
4. **Compatibility**: Existing rules continue to work without changes
5. **Observability**: Metrics and debugging support for sharing behavior

---

## Architecture Overview

### High-Level Components

```/dev/null/architecture.txt#L1-25
┌─────────────────────────────────────────────────────────────┐
│                    RETE Network                              │
│                                                              │
│  ┌────────────┐         ┌──────────────────────┐           │
│  │  TypeNode  │────────>│  AlphaNode (shared)  │           │
│  └────────────┘         └──────────┬───────────┘           │
│                                    │                         │
│                                    v                         │
│                         ┌─────────────────────┐             │
│                         │   JoinNode (shared) │<────┐       │
│                         │  ┌──────────────┐   │     │       │
│  ┌────────────┐         │  │ Left Memory  │   │     │       │
│  │  TypeNode  │────────>│  │ Right Memory │   │     │       │
│  └────────────┘         │  │Result Memory │   │     │       │
│                         │  └──────────────┘   │     │       │
│                         └─────────┬───────────┘     │       │
│                                   │                 │       │
│                                   v                 │       │
│                         ┌─────────────────────┐    │       │
│                         │  BetaSharingRegistry │────┘       │
│                         │  ┌──────────────┐   │            │
│                         │  │ Hash -> Node │   │            │
│                         │  │  LRU Cache   │   │            │
│                         │  │   Metrics    │   │            │
│                         │  └──────────────┘   │            │
│                         └─────────────────────┘            │
└─────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

| Component | Responsibility |
|-----------|---------------|
| **BetaSharingRegistry** | Central registry for shared JoinNodes; hash-based lookup; lifecycle management |
| **JoinNode** | Unchanged; performs joins between left and right inputs |
| **ConstraintPipelineBuilder** | Modified to call GetOrCreateJoinNode instead of NewJoinNode |
| **LifecycleManager** | Tracks reference counts; triggers cleanup when refcount reaches zero |
| **BetaBuildMetrics** | Collects statistics on sharing effectiveness |

---

## BetaSharingRegistry Design

### Core Data Structures

```/dev/null/structures.go#L1-35
type BetaSharingRegistry struct {
    // Shared JoinNodes indexed by canonical hash
    sharedJoinNodes map[string]*JoinNode
    
    // LRU cache for condition hash computation
    // Key: condition AST structure (JSON)
    // Value: computed hash string
    hashCache *LRUCache
    
    // Lifecycle manager for reference counting
    lifecycle *LifecycleManager
    
    // Metrics collection
    metrics *BetaBuildMetrics
    
    // Thread safety
    mutex sync.RWMutex
    
    // Configuration
    config BetaSharingConfig
}

type BetaSharingConfig struct {
    Enabled           bool
    HashCacheSize     int  // Default: 1000
    MaxSharedNodes    int  // Default: 10000
    EnableMetrics     bool
    NormalizeOrder    bool // Canonicalize variable order
}

type BetaBuildMetrics struct {
    TotalJoinNodesRequested   int64
    SharedJoinNodesReused     int64
    UniqueJoinNodesCreated    int64
    HashCacheHits             int64
    HashCacheMisses           int64
    AverageHashTimeNs         int64
}
```

### Storage Strategy

The registry maintains:

1. **Primary Map**: `map[string]*JoinNode`
   - Key: Canonical hash of join signature
   - Value: Shared JoinNode instance

2. **Hash Cache**: LRU cache for expensive hash computations
   - Caches normalized condition AST → hash mappings
   - Configurable size (default: 1000 entries)
   - Invalidated when normalization rules change

3. **Reverse Index** (optional, for debugging):
   - `map[*JoinNode]string`: Node → Hash lookup
   - Enables fast "which hash does this node belong to?" queries

### Lifecycle Management

```/dev/null/lifecycle.txt#L1-20
JoinNode Lifecycle with Sharing:

1. Rule Compilation
   ├─> GetOrCreateJoinNode(signature)
   ├─> Hash signature → check registry
   └─> If exists:
       ├─> Increment refcount (via LifecycleManager)
       └─> Return existing node
       If not exists:
       ├─> Create new JoinNode
       ├─> Register in registry
       ├─> Initialize refcount = 1
       └─> Return new node

2. Rule Execution
   └─> Shared JoinNode operates normally
       (all referencing rules share same memories)

3. Rule Removal
   ├─> Decrement refcount
   └─> If refcount == 0:
       ├─> Remove from registry
       ├─> Clear memories
       └─> Disconnect from network
```

**Key Points:**
- Reference counting managed by LifecycleManager (existing component)
- JoinNode is removed only when no rules reference it
- Memories are cleared on removal (left, right, result)
- Thread-safe increment/decrement with mutex protection

---

## Sharing Criteria & Compatibility

### When Can Two JoinNodes Be Shared?

Two JoinNodes can share the same instance if and only if:

#### 1. **Join Conditions Are Semantically Identical**

After normalization, the join condition ASTs must be equivalent:

```/dev/null/example.go#L1-15
// These CAN be shared (after normalization):
Condition A: customer.id == order.customerId
Condition B: order.customerId == customer.id  // Commutative equality

// These CANNOT be shared:
Condition C: customer.id == order.customerId
Condition D: customer.name == order.customerName  // Different fields
```

#### 2. **Variable Sets Are Compatible**

The left and right variable sets must match:

```/dev/null/variables.txt#L1-10
// Shareable:
JoinNode1: LeftVars=[order], RightVars=[customer]
JoinNode2: LeftVars=[order], RightVars=[customer]

// NOT shareable:
JoinNode3: LeftVars=[order], RightVars=[customer]
JoinNode4: LeftVars=[order, product], RightVars=[customer]
  // Different left variable count
```

#### 3. **Variable Types Are Consistent**

Type mappings for variables must be compatible:

```/dev/null/types.go#L1-12
// Shareable:
Node1: {order: "Order", customer: "Customer"}
Node2: {order: "Order", customer: "Customer"}

// NOT shareable:
Node3: {order: "Order", customer: "Customer"}
Node4: {order: "PurchaseOrder", customer: "Customer"}
  // Different type for 'order' variable
```

#### 4. **Join Condition Operators Match**

Operators must have identical semantics:

```/dev/null/operators.txt#L1-8
// Shareable:
a.x == b.y
a.x == b.y

// NOT shareable:
a.x == b.y
a.x != b.y  // Different operator
```

### Compatibility Test Algorithm

```/dev/null/algorithm.go#L1-30
func CanShareJoinNodes(node1, node2 *JoinNodeSignature) bool {
    // 1. Check variable set equality
    if !equalVarSets(node1.LeftVars, node2.LeftVars) {
        return false
    }
    if !equalVarSets(node1.RightVars, node2.RightVars) {
        return false
    }
    
    // 2. Check variable type compatibility
    if !compatibleVarTypes(node1.VarTypes, node2.VarTypes) {
        return false
    }
    
    // 3. Normalize and compare join conditions
    norm1 := NormalizeJoinCondition(node1.Condition)
    norm2 := NormalizeJoinCondition(node2.Condition)
    
    if !deepEqual(norm1, norm2) {
        return false
    }
    
    // 4. Check operator semantics
    if !compatibleOperators(norm1, norm2) {
        return false
    }
    
    return true
}
```

### Edge Cases

| Case | Shareable? | Reason |
|------|-----------|--------|
| Different variable names in condition | No | Variable bindings differ |
| Commutative operator order (a==b vs b==a) | Yes | Normalized to canonical form |
| Different evaluation order in multi-condition join | Yes | If conditions are logically equivalent |
| Same condition, different AllVariables | Depends | If AllVariables is superset, may share |
| Cascaded joins with partial overlap | Partial | Shared prefix, unique suffixes |

---

## Normalization & Hashing

### Normalization Process

Normalization converts join signatures into a **canonical form** that is identical for semantically equivalent joins.

#### Step 1: Unwrap AST Wrappers

Remove wrapper nodes that don't affect semantics:

```/dev/null/unwrap.go#L1-15
// Before:
{
  "type": "constraint",
  "expression": {
    "type": "binary",
    "op": "==",
    "left": {"type": "field", "var": "customer", "field": "id"},
    "right": {"type": "field", "var": "order", "field": "customerId"}
  }
}

// After:
{
  "type": "binary",
  "op": "==",
  "left": {"type": "field", "var": "customer", "field": "id"},
  "right": {"type": "field", "var": "order", "field": "customerId"}
}
```

#### Step 2: Canonicalize Commutative Operations

For commutative operators (==, !=, &&, ||), order operands canonically:

```/dev/null/commutative.go#L1-10
// Before:
order.customerId == customer.id

// After (canonical order: lexicographic by variable name):
customer.id == order.customerId

// Rule: left operand variable name <= right operand variable name
```

#### Step 3: Sort Variable Lists

Variables in LeftVars, RightVars, AllVariables are sorted alphabetically:

```/dev/null/sorting.go#L1-8
// Before:
LeftVars: [order, product, customer]

// After:
LeftVars: [customer, order, product]

// Ensures consistent ordering regardless of declaration order
```

#### Step 4: Normalize Type Mappings

Type maps are converted to sorted key-value pairs:

```/dev/null/types.json#L1-10
// Before:
{
  "order": "Order",
  "customer": "Customer",
  "product": "Product"
}

// After (sorted by key):
[["customer", "Customer"], ["order", "Order"], ["product", "Product"]]
```

#### Step 5: Create Canonical Structure

Combine all normalized components into a single structure:

```/dev/null/canonical.go#L1-20
type CanonicalJoinSignature struct {
    // Sorted variable lists
    LeftVars  []string
    RightVars []string
    AllVars   []string
    
    // Sorted type mappings
    VarTypes []VariableTypeMapping
    
    // Normalized condition AST
    Condition ConditionNode
    
    // Metadata (not part of hash)
    Version string  // "1.0" - allows future evolution
}

type VariableTypeMapping struct {
    VarName  string
    TypeName string
}
```

### Hashing Algorithm

```/dev/null/hashing.txt#L1-25
Hash Computation:

1. Serialize canonical structure to JSON
   ├─> Use deterministic ordering (sorted keys)
   ├─> No whitespace or formatting variations
   └─> Consistent encoding (UTF-8)

2. Compute SHA-256 hash
   └─> crypto/sha256.Sum256(jsonBytes)

3. Encode hash to hex string
   └─> hex.EncodeToString(hash[:8])  // First 8 bytes

4. Prepend prefix for readability
   └─> "join_" + hexString
   
Example Output:
   join_a3f2c1d4e5b6f7a8

Collision Probability:
   - 64-bit hash (8 bytes)
   - Birthday paradox: ~1% collision at 5M nodes
   - Acceptable for typical networks (< 100K nodes)
   - Full hash stored internally for verification
```

### Hash Caching

```/dev/null/cache.go#L1-20
type LRUCache struct {
    capacity int
    cache    map[string]string  // normalized JSON -> hash
    order    *list.List
    lookup   map[string]*list.Element
    mutex    sync.RWMutex
}

func (r *BetaSharingRegistry) JoinNodeHashCached(sig *JoinNodeSignature) (string, error) {
    // 1. Normalize signature
    canonical := NormalizeJoinSignature(sig)
    
    // 2. Serialize to JSON (deterministic)
    jsonKey := canonical.ToJSON()
    
    // 3. Check cache
    if hash, ok := r.hashCache.Get(jsonKey); ok {
        atomic.AddInt64(&r.metrics.HashCacheHits, 1)
        return hash, nil
    }
    
    // 4. Compute and cache
    hash := ComputeHash(canonical)
    r.hashCache.Put(jsonKey, hash)
    atomic.AddInt64(&r.metrics.HashCacheMisses, 1)
    return hash, nil
}
```

**Cache Invalidation:**
- Cache entries never expire (deterministic hashing)
- Cache is cleared only on:
  - Manual reset (testing/debugging)
  - Normalization algorithm version change
  - LRU eviction when cache is full

---

## Public API

### Core API Functions

```/dev/null/api.go#L1-80
// GetOrCreateJoinNode returns a shared JoinNode for the given signature.
// If a compatible node exists, it is returned and refcount is incremented.
// Otherwise, a new node is created and registered.
//
// Thread-safe: Multiple goroutines can call concurrently.
func (r *BetaSharingRegistry) GetOrCreateJoinNode(
    condition ConditionNode,
    leftVars []string,
    rightVars []string,
    allVars []string,
    varTypes map[string]string,
    storage Storage,
) (*JoinNode, string, bool, error) {
    
    // Build signature
    sig := &JoinNodeSignature{
        Condition:  condition,
        LeftVars:   leftVars,
        RightVars:  rightVars,
        AllVars:    allVars,
        VarTypes:   varTypes,
    }
    
    // Compute hash (cached)
    hash, err := r.JoinNodeHashCached(sig)
    if err != nil {
        return nil, "", false, err
    }
    
    // Fast path: read lock for lookup
    r.mutex.RLock()
    if existing, ok := r.sharedJoinNodes[hash]; ok {
        r.mutex.RUnlock()
        
        // Increment refcount
        r.lifecycle.AddRuleReference(existing.ID())
        
        atomic.AddInt64(&r.metrics.SharedJoinNodesReused, 1)
        return existing, hash, true, nil
    }
    r.mutex.RUnlock()
    
    // Slow path: write lock for creation
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    // Double-check: another goroutine may have created it
    if existing, ok := r.sharedJoinNodes[hash]; ok {
        r.lifecycle.AddRuleReference(existing.ID())
        atomic.AddInt64(&r.metrics.SharedJoinNodesReused, 1)
        return existing, hash, true, nil
    }
    
    // Create new JoinNode
    node := NewJoinNode(condition, leftVars, rightVars, allVars, varTypes, storage)
    
    // Register in registry
    r.sharedJoinNodes[hash] = node
    r.lifecycle.RegisterNode(node, hash)
    
    atomic.AddInt64(&r.metrics.UniqueJoinNodesCreated, 1)
    return node, hash, false, nil
}

// RegisterJoinNode explicitly registers an existing JoinNode.
// Used when migration from non-shared to shared nodes.
func (r *BetaSharingRegistry) RegisterJoinNode(
    node *JoinNode,
    hash string,
) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    if existing, ok := r.sharedJoinNodes[hash]; ok {
        if existing != node {
            return fmt.Errorf("hash collision: different node already registered")
        }
        // Already registered, no-op
        return nil
    }
    
    r.sharedJoinNodes[hash] = node
    r.lifecycle.RegisterNode(node, hash)
    return nil
}

// ReleaseJoinNode decrements refcount and removes node if unused.
func (r *BetaSharingRegistry) ReleaseJoinNode(hash string) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    node, ok := r.sharedJoinNodes[hash]
    if !ok {
        return fmt.Errorf("join node not found: %s", hash)
    }
    
    refcount := r.lifecycle.RemoveRuleReference(node.ID())
    
    if refcount == 0 {
        // No more references, clean up
        delete(r.sharedJoinNodes, hash)
        
        // Clear memories
        node.ClearMemories()
        
        // Disconnect from network (handled by LifecycleManager)
        r.lifecycle.UnregisterNode(node.ID())
    }
    
    return nil
}

// GetSharingStats returns current sharing metrics.
func (r *BetaSharingRegistry) GetSharingStats() *BetaSharingStats {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    return &BetaSharingStats{
        TotalSharedNodes:        len(r.sharedJoinNodes),
        TotalRequests:           atomic.LoadInt64(&r.metrics.TotalJoinNodesRequested),
        SharedReuses:            atomic.LoadInt64(&r.metrics.SharedJoinNodesReused),
        UniqueCreations:         atomic.LoadInt64(&r.metrics.UniqueJoinNodesCreated),
        SharingRatio:            r.calculateSharingRatio(),
        HashCacheHitRate:        r.calculateCacheHitRate(),
        AverageHashTimeMs:       float64(atomic.LoadInt64(&r.metrics.AverageHashTimeNs)) / 1e6,
    }
}

// ListSharedJoinNodes returns all shared join node hashes.
func (r *BetaSharingRegistry) ListSharedJoinNodes() []string {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    hashes := make([]string, 0, len(r.sharedJoinNodes))
    for hash := range r.sharedJoinNodes {
        hashes = append(hashes, hash)
    }
    sort.Strings(hashes)
    return hashes
}

// GetSharedJoinNodeDetails returns detailed info about a shared node.
func (r *BetaSharingRegistry) GetSharedJoinNodeDetails(hash string) (*JoinNodeDetails, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    node, ok := r.sharedJoinNodes[hash]
    if !ok {
        return nil, fmt.Errorf("join node not found: %s", hash)
    }
    
    refcount := r.lifecycle.GetReferenceCount(node.ID())
    
    return &JoinNodeDetails{
        Hash:            hash,
        NodeID:          node.ID(),
        ReferenceCount:  refcount,
        LeftVars:        node.LeftVariables,
        RightVars:       node.RightVariables,
        AllVars:         node.AllVariables,
        VarTypes:        node.VariableTypes,
        LeftMemorySize:  len(node.LeftMemory),
        RightMemorySize: len(node.RightMemory),
        ResultMemorySize: len(node.ResultMemory),
    }, nil
}
```

### Statistics & Monitoring

```/dev/null/stats.go#L1-25
type BetaSharingStats struct {
    // Current state
    TotalSharedNodes   int
    
    // Cumulative counters
    TotalRequests      int64
    SharedReuses       int64
    UniqueCreations    int64
    
    // Computed metrics
    SharingRatio       float64  // SharedReuses / TotalRequests
    HashCacheHitRate   float64  // CacheHits / (CacheHits + CacheMisses)
    AverageHashTimeMs  float64
}

type JoinNodeDetails struct {
    Hash              string
    NodeID            string
    ReferenceCount    int
    LeftVars          []string
    RightVars         []string
    AllVars           []string
    VarTypes          map[string]string
    LeftMemorySize    int
    RightMemorySize   int
    ResultMemorySize  int
}
```

---

## Integration with Builder & Lifecycle

### Builder Integration

**Before (without sharing):**

```/dev/null/before.go#L1-15
func (b *ConstraintPipelineBuilder) createBinaryJoinRule(
    var1, var2 *VariableTypeMapping,
    condition ConditionNode,
) (*JoinNode, error) {
    
    // Always create new JoinNode
    joinNode := NewJoinNode(
        condition,
        []string{var1.Name},
        []string{var2.Name},
        []string{var1.Name, var2.Name},
        map[string]string{var1.Name: var1.Type, var2.Name: var2.Type},
        b.storage,
    )
    
    // Connect to network...
    return joinNode, nil
}
```

**After (with sharing):**

```/dev/null/after.go#L1-30
func (b *ConstraintPipelineBuilder) createBinaryJoinRule(
    var1, var2 *VariableTypeMapping,
    condition ConditionNode,
) (*JoinNode, error) {
    
    // Use BetaSharingRegistry
    joinNode, hash, wasShared, err := b.betaRegistry.GetOrCreateJoinNode(
        condition,
        []string{var1.Name},
        []string{var2.Name},
        []string{var1.Name, var2.Name},
        map[string]string{var1.Name: var1.Type, var2.Name: var2.Type},
        b.storage,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to get/create join node: %w", err)
    }
    
    // Log sharing info
    if wasShared {
        b.logger.Debugf("Reused shared JoinNode %s for %s+%s", hash, var1.Name, var2.Name)
    } else {
        b.logger.Debugf("Created new JoinNode %s for %s+%s", hash, var1.Name, var2.Name)
    }
    
    // Connect to network (avoid duplicate connections if shared)...
    b.connectToNetwork(joinNode, var1, var2, wasShared)
    
    return joinNode, nil
}
```

### Connection Handling

**Challenge:** When a JoinNode is shared, we must avoid creating duplicate pass-through AlphaNodes.

**Solution:**

```/dev/null/connections.go#L1-35
func (b *ConstraintPipelineBuilder) connectToNetwork(
    joinNode *JoinNode,
    leftVar, rightVar *VariableTypeMapping,
    wasShared bool,
) error {
    
    // Left input: TypeNode -> PassThrough AlphaNode -> JoinNode
    leftTypeNode := b.network.GetTypeNode(leftVar.Type)
    
    if !wasShared {
        // New node: create connection
        leftAlpha := NewPassThroughAlphaNode(leftVar.Name)
        leftTypeNode.AddChild(leftAlpha)
        leftAlpha.AddChild(joinNode)
        b.network.RegisterAlphaNode(leftAlpha)
    } else {
        // Shared node: check if connection already exists
        if !b.connectionExists(leftTypeNode, joinNode) {
            leftAlpha := NewPassThroughAlphaNode(leftVar.Name)
            leftTypeNode.AddChild(leftAlpha)
            leftAlpha.AddChild(joinNode)
            b.network.RegisterAlphaNode(leftAlpha)
        }
        // Else: connection exists, no-op
    }
    
    // Right input: similar logic
    rightTypeNode := b.network.GetTypeNode(rightVar.Type)
    if !wasShared || !b.connectionExists(rightTypeNode, joinNode) {
        rightAlpha := NewPassThroughAlphaNode(rightVar.Name)
        rightTypeNode.AddChild(rightAlpha)
        rightAlpha.AddChild(joinNode)
        b.network.RegisterAlphaNode(rightAlpha)
    }
    
    return nil
}
```

### Lifecycle Integration

```/dev/null/lifecycle_integration.go#L1-40
// Rule Addition
func (n *ReteNetwork) AddRule(rule *Rule) error {
    // ... parse and build rule ...
    
    // Builder uses BetaSharingRegistry.GetOrCreateJoinNode()
    // which internally calls:
    //   lifecycle.RegisterNode(joinNode, hash)
    //   lifecycle.AddRuleReference(joinNode.ID())
    
    // ... compile and register ...
    return nil
}

// Rule Removal
func (n *ReteNetwork) RemoveRule(ruleName string) error {
    rule, ok := n.rules[ruleName]
    if !ok {
        return fmt.Errorf("rule not found: %s", ruleName)
    }
    
    // For each JoinNode in the rule's network path:
    for _, joinHash := range rule.JoinNodeHashes {
        // Decrement refcount and potentially remove
        err := n.betaRegistry.ReleaseJoinNode(joinHash)
        if err != nil {
            return fmt.Errorf("failed to release join node %s: %w", joinHash, err)
        }
    }
    
    // Remove rule from registry
    delete(n.rules, ruleName)
    
    return nil
}

// Network Shutdown
func (n *ReteNetwork) Shutdown() error {
    // Clear all shared nodes (refcounts will be 0)
    for _, hash := range n.betaRegistry.ListSharedJoinNodes() {
        n.betaRegistry.ReleaseJoinNode(hash)
    }
    
    // ... other cleanup ...
    return nil
}
```

---

## Sequence Diagrams

### Sequence 1: GetOrCreateJoinNode - Cache Hit (Shared Node)

```/dev/null/sequence1.txt#L1-25
Builder                Registry              HashCache        Lifecycle        JoinNode
  |                       |                      |                |                |
  |--GetOrCreateJoinNode->|                      |                |                |
  |                       |                      |                |                |
  |                       |--JoinNodeHashCached->|                |                |
  |                       |                      |                |                |
  |                       |<--return hash--------|                |                |
  |                       |    (cache HIT)       |                |                |
  |                       |                      |                |                |
  |                       |--RLock()             |                |                |
  |                       |--lookup(hash)        |                |                |
  |                       |  [FOUND]             |                |                |
  |                       |--RUnlock()           |                |                |
  |                       |                      |                |                |
  |                       |--AddRuleReference--->|                |                |
  |                       |                      |                |                |
  |                       |                      |   increment refcount            |
  |                       |                      |                |                |
  |<--return (node, hash, true)------------------|                |                |
  |                       |                      |                |                |
  |--use shared node------|----------------------|----------------|--------------->|
  |                       |                      |                |                |
```

### Sequence 2: GetOrCreateJoinNode - Cache Miss (New Node)

```/dev/null/sequence2.txt#L1-30
Builder                Registry              HashCache        Lifecycle        JoinNode
  |                       |                      |                |                |
  |--GetOrCreateJoinNode->|                      |                |                |
  |                       |                      |                |                |
  |                       |--JoinNodeHashCached->|                |                |
  |                       |                      |                |                |
  |                       |  [normalize sig]     |                |                |
  |                       |  [compute SHA-256]   |                |                |
  |                       |  [cache result]      |                |                |
  |                       |                      |                |                |
  |                       |<--return hash--------|                |                |
  |                       |    (cache MISS)      |                |                |
  |                       |                      |                |                |
  |                       |--RLock()             |                |                |
  |                       |--lookup(hash)        |                |                |
  |                       |  [NOT FOUND]         |                |                |
  |                       |--RUnlock()           |                |                |
  |                       |                      |                |                |
  |                       |--Lock()              |                |                |
  |                       |--double-check        |                |                |
  |                       |  [still not found]   |                |                |
  |                       |                      |                |                |
  |                       |--NewJoinNode()-------|----------------|--------------->|
  |                       |                      |                |                |
  |                       |<--return new node----|----------------|----------------|
  |                       |                      |                |                |
  |                       |--store(hash, node)   |                |                |
  |                       |--RegisterNode------->|                |                |
  |                       |     (refcount=1)     |                |                |
  |                       |--Unlock()            |                |                |
  |                       |                      |                |                |
  |<--return (node, hash, false)-----------------|                |                |
  |                       |                      |                |                |
```

### Sequence 3: ReleaseJoinNode - Refcount Zero (Cleanup)

```/dev/null/sequence3.txt#L1-25
Network                Registry              Lifecycle        JoinNode
  |                       |                      |                |
  |--RemoveRule()         |                      |                |
  |  [iterate join nodes] |                      |                |
  |                       |                      |                |
  |--ReleaseJoinNode(hash)|                      |                |
  |                       |                      |                |
  |                       |--Lock()              |                |
  |                       |--lookup(hash)        |                |
  |                       |  [FOUND]             |                |
  |                       |                      |                |
  |                       |--RemoveRuleReference>|                |
  |                       |                      |                |
  |                       |<--return refcount=0--|                |
  |                       |                      |                |
  |                       |--delete(hash)        |                |
  |                       |                      |                |
  |                       |--ClearMemories()-----|--------------->|
  |                       |                      |                |
  |                       |--UnregisterNode----->|                |
  |                       |                      |                |
  |                       |--Unlock()            |                |
  |                       |                      |                |
  |<--return nil----------|                      |                |
  |                       |                      |                |
```

---

## Usage Examples

### Example 1: Simple Binary Join Sharing

**Input Rules:**

```/dev/null/rules1.tsd#L1-20
rule "HighValueCustomerOrder" {
    when {
        customer: Customer(customer.tier == "GOLD")
        order: Order(order.customerId == customer.id && order.value > 1000)
    }
    then {
        log("High value order from gold customer")
    }
}

rule "RecentCustomerOrder" {
    when {
        customer: Customer(customer.signupDate > "2024-01-01")
        order: Order(order.customerId == customer.id && order.value > 500)
    }
    then {
        log("Recent customer placed order")
    }
}
```

**Network Structure (Without Sharing):**

```/dev/null/without_sharing.txt#L1-15
TypeNode[Customer]                    TypeNode[Order]
      |                                     |
      v                                     v
  AlphaNode[tier=="GOLD"]         AlphaNode[value>1000]
      |                                     |
      +---------------JoinNode1-------------+
                         |
                         v
                   TerminalNode[Rule1]

  AlphaNode[signup>"2024"]        AlphaNode[value>500]
      |                                     |
      +---------------JoinNode2-------------+  // DUPLICATE!
                         |
                         v
                   TerminalNode[Rule2]
```

**Network Structure (With Sharing):**

```/dev/null/with_sharing.txt#L1-15
TypeNode[Customer]                    TypeNode[Order]
      |                                     |
      +-------+                       +-----+
      |       |                       |     |
      v       v                       v     v
Alpha[tier] Alpha[signup]      Alpha[val>1000] Alpha[val>500]
      |       |                       |           |
      |       |                       v           v
      |       |                  JoinNode (SHARED!)
      |       |                   /         \
      |       +------------------+           +---> TerminalNode[Rule2]
      |                                      |
      +--------------------------------------+
                                             |
                                             v
                                       TerminalNode[Rule1]
```

**Result:**
- 1 JoinNode instead of 2 (50% reduction)
- Both rules share the same join memories
- Hash: `join_a3f2c1d4` (based on `order.customerId == customer.id`)

### Example 2: Cascade Join with Partial Sharing

**Input Rule:**

```/dev/null/rules2.tsd#L1-10
rule "CustomerOrderProduct" {
    when {
        customer: Customer(customer.tier == "GOLD")
        order: Order(order.customerId == customer.id)
        product: Product(product.id == order.productId && product.inStock == true)
    }
    then {
        ship(order, product, customer)
    }
}
```

**Cascade Structure:**

```/dev/null/cascade.txt#L1-20
Step 1: Join Customer + Order
   JoinNode1[customer.id == order.customerId]
   Hash: join_abc123
   
Step 2: Join (Customer+Order) + Product
   JoinNode2[(customer,order).productId == product.id]
   Hash: join_def456

If another rule has the same cascade:
   rule "TrackGoldCustomerProducts" {
       when {
           customer: Customer(customer.tier == "GOLD")
           order: Order(order.customerId == customer.id)
           product: Product(product.id == order.productId)
       }
       ...
   }

Then:
   - JoinNode1 is SHARED (both rules use it)
   - JoinNode2 may be DIFFERENT (due to different product filter)
```

### Example 3: Non-Shareable Joins (Different Conditions)

```/dev/null/non_shareable.tsd#L1-20
rule "Rule1" {
    when {
        customer: Customer(customer.id > 1000)
        order: Order(order.customerId == customer.id)  // Join on .id
    }
    then { ... }
}

rule "Rule2" {
    when {
        customer: Customer(customer.tier == "GOLD")
        order: Order(order.customerEmail == customer.email)  // Join on .email
    }
    then { ... }
}

// Result: Two DIFFERENT JoinNodes
//   Join1: customer.id == order.customerId
//   Join2: customer.email == order.customerEmail
```

### Example 4: Metrics & Debugging

```/dev/null/metrics_example.go#L1-40
// Query sharing statistics
stats := network.BetaRegistry.GetSharingStats()

fmt.Printf("Beta Sharing Statistics:\n")
fmt.Printf("  Total shared nodes:    %d\n", stats.TotalSharedNodes)
fmt.Printf("  Total requests:        %d\n", stats.TotalRequests)
fmt.Printf("  Shared reuses:         %d\n", stats.SharedReuses)
fmt.Printf("  Unique creations:      %d\n", stats.UniqueCreations)
fmt.Printf("  Sharing ratio:         %.2f%%\n", stats.SharingRatio*100)
fmt.Printf("  Hash cache hit rate:   %.2f%%\n", stats.HashCacheHitRate*100)
fmt.Printf("  Avg hash time:         %.3f ms\n", stats.AverageHashTimeMs)

// Output example:
// Beta Sharing Statistics:
//   Total shared nodes:    15
//   Total requests:        42
//   Shared reuses:         27
//   Unique creations:      15
//   Sharing ratio:         64.29%
//   Hash cache hit rate:   85.71%
//   Avg hash time:         0.123 ms

// Inspect specific shared node
details, _ := network.BetaRegistry.GetSharedJoinNodeDetails("join_a3f2c1d4")

fmt.Printf("Join Node Details:\n")
fmt.Printf("  Hash:              %s\n", details.Hash)
fmt.Printf("  Node ID:           %s\n", details.NodeID)
fmt.Printf("  Reference count:   %d\n", details.ReferenceCount)
fmt.Printf("  Left vars:         %v\n", details.LeftVars)
fmt.Printf("  Right vars:        %v\n", details.RightVars)
fmt.Printf("  Left memory size:  %d tokens\n", details.LeftMemorySize)
fmt.Printf("  Right memory size: %d tokens\n", details.RightMemorySize)

// Output example:
// Join Node Details:
//   Hash:              join_a3f2c1d4
//   Node ID:           node_12345
//   Reference count:   3
//   Left vars:         [customer]
//   Right vars:        [order]
//   Left memory size:  25 tokens
//   Right memory size: 120 tokens
```

---

## Performance Considerations

### Memory Impact

**Before Sharing:**
```/dev/null/memory_before.txt#L1-8
100 rules with FK joins:
  - 100 JoinNodes
  - Each with 3 memory maps (left, right, result)
  - Average 1KB per memory map
  Total: 100 × 3 × 1KB = 300KB
```

**After Sharing (50% sharing ratio):**
```/dev/null/memory_after.txt#L1-10
100 rules with FK joins:
  - 50 unique JoinNodes (50 shared)
  - Each with 3 memory maps
  - Average 1KB per memory map
  Total: 50 × 3 × 1KB = 150KB
  
Savings: 150KB (50% reduction)
```

### Runtime Performance

**Join Evaluation:**
- Shared nodes perform joins once per activation
- Multiple rules benefit from same evaluation
- Expected speedup: 20-40% for rules with shared joins

**Hash Computation:**
- Normalized once per signature
- Cached in LRU (typical hit rate: 80-90%)
- Sub-millisecond overhead even on cache miss

**Lock Contention:**
- Read-heavy workload (lookups >> modifications)
- RWMutex allows concurrent reads
- Write locks only during node creation/removal
- Expected contention: minimal (<1% overhead)

### Scalability

| Network Size | Shared Nodes | Hash Cache Size | Memory Usage | Lookup Time |
|--------------|--------------|-----------------|--------------|-------------|
| 100 rules    | ~50          | 1000            | ~150KB       | <0.1ms      |
| 1,000 rules  | ~400         | 2000            | ~1.2MB       | <0.2ms      |
| 10,000 rules | ~3,000       | 5000            | ~9MB         | <0.5ms      |

**Key Observations:**
- Sublinear growth in shared nodes (rules share patterns)
- Hash cache prevents redundant computation
- Lookup time scales logarithmically with network size

### Optimization Opportunities

1. **Bloom Filters** for fast negative lookups (future enhancement)
2. **Shard Hash Maps** to reduce lock contention (if needed)
3. **Lazy Memory Allocation** for rarely-activated nodes
4. **Compressed Hash Storage** using integer IDs instead of strings

---

## Testing Strategy

### Unit Tests

```/dev/null/unit_tests.txt#L1-20
Package: rete/beta_sharing_test.go

TestGetOrCreateJoinNode_NewNode
  - Create first instance of a join signature
  - Verify: wasShared=false, refcount=1

TestGetOrCreateJoinNode_SharedNode
  - Create same signature twice
  - Verify: wasShared=true, same node returned, refcount=2

TestJoinNodeHashCached_Normalization
  - Test commutative operator normalization
  - Test variable ordering normalization
  - Verify: equivalent signatures produce same hash

TestJoinNodeHashCached_CacheHit
  - Hash same signature twice
  - Verify: second call hits cache

TestReleaseJoinNode_RefcountDecrement
  - Create shared node (refcount=2)
  - Release once
  - Verify: node still exists, refcount=1

TestReleaseJoinNode_Cleanup
  - Create node (refcount=1)
  - Release
  - Verify: node removed from registry

TestCompatibility_DifferentVariables
  - Signatures with different variables
  - Verify: different hashes

TestCompatibility_DifferentTypes
  - Same variables, different types
  - Verify: different hashes

TestConcurrency_ParallelGetOrCreate
  - 100 goroutines request same signature
  - Verify: only one node created, refcount=100

TestMetrics_SharingRatio
  - Create mix of shared/unique nodes
  - Verify: metrics match expectations
```

### Integration Tests

```/dev/null/integration_tests.txt#L1-25
Package: rete/integration_test.go

TestBetaSharing_TwoRulesSharedJoin
  - Add two rules with identical join
  - Verify: 1 JoinNode created
  - Assert facts through network
  - Verify: both rules activate correctly

TestBetaSharing_CascadePartialSharing
  - Add two rules with shared first join, different second join
  - Verify: 1st join shared, 2nd joins unique
  - Assert facts, verify correctness

TestBetaSharing_RuleRemoval
  - Add 3 rules sharing same join
  - Remove 1 rule
  - Verify: JoinNode still exists, refcount=2
  - Remove 2nd rule
  - Verify: JoinNode still exists, refcount=1
  - Remove 3rd rule
  - Verify: JoinNode removed

TestBetaSharing_MemoryBehavior
  - Add facts to shared JoinNode
  - Verify: both rules see same memory contents
  - Retract fact
  - Verify: both rules retract correctly

TestBetaSharing_DisabledFlag
  - Disable beta sharing (feature flag)
  - Add two rules with identical join
  - Verify: 2 JoinNodes created (no sharing)
```

### Performance Tests

```/dev/null/perf_tests.txt#L1-20
Package: rete/beta_sharing_bench_test.go

BenchmarkGetOrCreateJoinNode_New
  - Measure time to create unique JoinNodes
  - Baseline for comparison

BenchmarkGetOrCreateJoinNode_Shared
  - Measure time to retrieve shared JoinNodes
  - Should be faster than creation

BenchmarkHashComputation_CacheMiss
  - Measure normalization + hashing time
  - Target: <1ms per signature

BenchmarkHashComputation_CacheHit
  - Measure cached hash retrieval
  - Target: <0.1ms per lookup

BenchmarkConcurrentAccess_ReadHeavy
  - 90% reads (GetOrCreate existing), 10% writes (new nodes)
  - Measure throughput under concurrency

BenchmarkMemoryUsage_WithSharing
  - Compare memory usage: 1000 rules with/without sharing
  - Verify 30-50% reduction
```

---

## Migration & Rollout

### Feature Flag

```/dev/null/feature_flag.go#L1-15
type NetworkConfig struct {
    // ... existing config ...
    
    // Beta sharing configuration
    EnableBetaSharing      bool   `yaml:"enable_beta_sharing"`
    BetaHashCacheSize      int    `yaml:"beta_hash_cache_size"`
    BetaMaxSharedNodes     int    `yaml:"beta_max_shared_nodes"`
    BetaSharingMetrics     bool   `yaml:"beta_sharing_metrics"`
}

// Default values
var DefaultConfig = NetworkConfig{
    EnableBetaSharing:      false,  // Disabled by default
    BetaHashCacheSize:      1000,
    BetaMaxSharedNodes:     10000,
    BetaSharingMetrics:     true,
}
```

### Rollout Phases

**Phase 1: Internal Testing (Week 1-2)**
- Enable in development environment
- Run existing test suite
- Verify correctness and performance
- Fix any issues discovered

**Phase 2: Canary Deployment (Week 3)**
- Enable for 10% of production workloads
- Monitor metrics:
  - Sharing ratio
  - Memory usage
  - Latency (p50, p95, p99)
  - Error rates
- Rollback if issues detected

**Phase 3: Gradual Rollout (Week 4-5)**
- Increase to 50% of workloads
- Continue monitoring
- Gather feedback from users
- Tune configuration if needed

**Phase 4: Full Deployment (Week 6)**
- Enable for 100% of workloads
- Change default to `EnableBetaSharing: true`
- Document in user guide

**Phase 5: Optimization (Week 7+)**
- Analyze sharing patterns in production
- Identify opportunities for improved normalization
- Tune cache sizes based on real usage

### Backward Compatibility

**Guarantees:**
- Existing rules compile and execute identically
- Serialized network formats unchanged
- API surface additions only (no breaking changes)
- Can disable sharing via feature flag

**Migration Path:**
- No code changes required for existing users
- Opt-in via configuration
- Gradual adoption supported

---

## Future Enhancements

### 1. Advanced Normalization

**Goal:** Share more nodes by recognizing additional equivalences

**Examples:**
```/dev/null/advanced_norm.txt#L1-10
// Recognize associativity:
(a && b) && c  ≡  a && (b && c)

// Recognize distributivity:
a && (b || c)  ≡  (a && b) || (a && c)

// Recognize De Morgan's laws:
!(a && b)  ≡  !a || !b
```

**Implementation:**
- AST canonicalization pass
- Pattern matching and rewriting
- Verification testing required

### 2. Partial Sharing

**Goal:** Share JoinNodes even when variable sets partially overlap

**Example:**
```/dev/null/partial_sharing.txt#L1-15
// Rule 1: needs {customer, order}
JoinNode1: customer.id == order.customerId
  LeftVars: [customer]
  RightVars: [order]

// Rule 2: needs {customer, order, product}
JoinNode2: customer.id == order.customerId
  LeftVars: [customer]
  RightVars: [order]

// Share JoinNode, but Rule 2 extends with additional variables
```

**Challenges:**
- AllVariables field differs
- Must ensure extended variables don't conflict
- May require node wrapping or adapters

### 3. Adaptive Cache Sizing

**Goal:** Dynamically adjust hash cache size based on hit rate

**Strategy:**
- Monitor cache hit rate over time windows
- Increase size if hit rate < 80%
- Decrease size if hit rate > 95% (wasted memory)
- Use exponential moving average

### 4. Distributed Sharing

**Goal:** Share JoinNodes across multiple RETE instances

**Use Case:** Microservices architecture with shared rule base

**Approach:**
- Centralized sharing registry (Redis/etcd)
- Consensus protocol for node creation
- Local cache for fast lookups
- Fallback to local-only sharing

### 5. Sharing Visualization

**Goal:** Visualize which rules share which nodes

**Features:**
- Graph UI showing rule → JoinNode edges
- Heat map of most-shared nodes
- Dependency analysis (which rules would break if node X removed)
- Optimization suggestions

### 6. Machine Learning Optimization

**Goal:** Learn optimal normalization strategies from production data

**Approach:**
- Collect join condition patterns
- Train model to predict sharing likelihood
- Suggest rule refactoring for better sharing
- Auto-generate shared condition libraries

---

## Appendix A: Hash Collision Handling

**Probability:**
With 64-bit hashes and N nodes:
- Birthday paradox: P(collision) ≈ N²/(2×2⁶⁴)
- At 10,000 nodes: P ≈ 0.0000027%
- Negligible risk

**Detection:**
```/dev/null/collision_detection.go#L1-15
func (r *BetaSharingRegistry) detectCollision(hash string, sig *JoinNodeSignature) bool {
    node := r.sharedJoinNodes[hash]
    if node == nil {
        return false
    }
    
    // Compare full canonical structures
    existingSig := node.GetSignature()
    if !deepEqual(sig, existingSig) {
        // COLLISION DETECTED!
        return true
    }
    return false
}
```

**Resolution:**
- Log error with full details
- Fall back to creating new node with hash suffix
- Alert monitoring system
- Investigate normalization bug

---

## Appendix B: Debugging Tools

### 1. Hash Inspector

```/dev/null/hash_inspector.sh#L1-10
# CLI tool to inspect join node hashes
$ rete-debug beta-hash inspect join_a3f2c1d4

Join Node: join_a3f2c1d4
  Reference Count: 3
  Left Vars: [customer]
  Right Vars: [order]
  Condition: customer.id == order.customerId
  Rules Using:
    - HighValueCustomerOrder
    - RecentCustomerOrder
    - LoyaltyRewardCalculation
```

### 2. Sharing Analyzer

```/dev/null/sharing_analyzer.sh#L1-15
# Analyze sharing patterns in rule base
$ rete-debug beta-sharing analyze

Beta Sharing Analysis:
  Total JoinNodes: 50
  Shared: 32 (64%)
  Unique: 18 (36%)
  
  Most Shared Nodes:
    1. join_a3f2c1d4: 12 rules (customer-order FK join)
    2. join_7f8e9d0a: 8 rules (order-product FK join)
    3. join_b1c2d3e4: 5 rules (customer-address join)
  
  Sharing Opportunities:
    - 4 rules have similar joins differing only in field order
      → Recommend enabling advanced normalization
```

### 3. Memory Profiler

```/dev/null/memory_profiler.sh#L1-10
# Compare memory usage with/without sharing
$ rete-debug beta-sharing memory-profile

Memory Usage:
  With Sharing:    2.4 MB (50 nodes)
  Without Sharing: 4.8 MB (100 nodes)
  Savings:         2.4 MB (50%)
  
  Per-Node Breakdown:
    Avg Left Memory:   120 tokens × 48 bytes = 5.8 KB
    Avg Right Memory:  200 tokens × 48 bytes = 9.6 KB
    Avg Result Memory: 50 tokens × 96 bytes = 4.8 KB
```

---

## Appendix C: Configuration Examples

### Minimal Configuration

```/dev/null/config_minimal.yaml#L1-5
network:
  enable_beta_sharing: true
  
# Uses all default values:
#   beta_hash_cache_size: 1000
#   beta_max_shared_nodes: 10000
```

### Production Configuration

```/dev/null/config_production.yaml#L1-15
network:
  enable_beta_sharing: true
  beta_hash_cache_size: 5000      # Larger cache for high-throughput
  beta_max_shared_nodes: 50000    # Support large rule bases
  beta_sharing_metrics: true      # Enable detailed metrics
  
  # Lifecycle integration
  enable_lifecycle_manager: true
  node_gc_interval: "5m"          # Clean up unused nodes every 5 min
  
  # Logging
  log_level: "info"
  log_beta_sharing_events: true   # Log sharing decisions
```

### Development Configuration

```/dev/null/config_development.yaml#L1-10
network:
  enable_beta_sharing: true
  beta_hash_cache_size: 100       # Small cache for testing
  beta_max_shared_nodes: 1000
  beta_sharing_metrics: true
  
  # Debug settings
  log_level: "debug"
  log_beta_sharing_events: true
  log_hash_computations: true     # Verbose hashing logs
```

---

## Conclusion

The BetaSharingRegistry provides a robust, performant solution for sharing JoinNodes in the TSD RETE engine. By following proven patterns from AlphaSharingRegistry and adapting them to the unique challenges of Beta nodes (two inputs, complex join conditions), we achieve:

- **Significant memory savings** (30-50%)
- **Performance improvements** (20-40%)
- **Backward compatibility**
- **Production-ready observability**

The design is modular, testable, and extensible, providing a solid foundation for future optimizations.

---

**Next Steps:**
1. Review and approve this design
2. Implement `rete/beta_sharing.go` (core registry)
3. Implement `rete/beta_normalization.go` (hashing and normalization)
4. Update `rete/constraint_pipeline_builder.go` (integration)
5. Add comprehensive tests
6. Begin rollout with feature flag

**Questions or Feedback:**
Please direct to the TSD RETE team or open an issue in the repository.

---

**Document History:**
- v1.0 (2024): Initial design