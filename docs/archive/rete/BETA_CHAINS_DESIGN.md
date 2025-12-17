<!--
Copyright (c) 2025 TSD Contributors
Licensed under the MIT License
See LICENSE file in the project root for full license text
-->

# BetaChains Design Document

**Version:** 1.0  
**Date:** 2025-01-XX  
**Status:** Design Phase  
**Related:** `BETA_SHARING_DESIGN.md`, `alpha_chain_builder.go`, `BETA_NODES_ANALYSIS.md`

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Background & Motivation](#background--motivation)
3. [BetaChain Structure](#betachain-structure)
4. [BetaChainBuilder Design](#betachainbuilder-design)
5. [Construction Algorithm](#construction-algorithm)
6. [Optimization Strategies](#optimization-strategies)
7. [Integration with BetaSharingRegistry](#integration-with-betasharingregistry)
8. [Comparison with AlphaChains](#comparison-with-alphachains)
9. [Examples & Use Cases](#examples--use-cases)
10. [Implementation Plan](#implementation-plan)
11. [Performance Expectations](#performance-expectations)
12. [Testing Strategy](#testing-strategy)

---

## Executive Summary

**BetaChains** are ordered sequences of **JoinNodes** (beta nodes) that form the join backbone of RETE rules. Similar to how AlphaChains optimize sequences of single-variable conditions, BetaChains optimize multi-variable join operations by:

- **Reusing common join prefixes** across rules
- **Ordering joins by selectivity** to minimize intermediate result sizes
- **Sharing identical join nodes** via BetaSharingRegistry
- **Building progressively** to enable incremental rule compilation

**Key Benefits:**
- 40-60% reduction in redundant join nodes (estimated)
- 25-45% faster rule compilation for rules with shared join patterns
- Optimal join ordering reduces working memory size by 30-50%
- Seamless integration with existing BetaSharingRegistry

**Deliverables:**
- `BetaChain` struct with metadata and validation
- `BetaChainBuilder` with progressive construction
- Join ordering optimizer based on selectivity
- Integration with pipeline builder and lifecycle manager

---

## Background & Motivation

### Current State

The current `constraint_pipeline_builder.go` creates JoinNodes in two strategies:

1. **Binary Join** (2 variables): Single JoinNode with left and right inputs
2. **Cascade Join** (3+ variables): Chain of JoinNodes, each adding one variable

**Problems:**
- No reuse of common join subsequences across rules
- Join order is determined by variable declaration order (not optimal)
- No centralized abstraction for join chains
- Difficult to analyze and optimize join patterns

### RETE Beta Network Theory

In RETE literature, the **beta network** consists of join nodes that combine tokens from:
- **Left input:** Partial matches (tokens) from upstream joins
- **Right input:** New facts from alpha network

**Key Insight:** Many rules share common join patterns. For example:

```
Rule 1: (Person p) ⋈ (Address a) ⋈ (Order o)
Rule 2: (Person p) ⋈ (Address a) ⋈ (Invoice i)
```

Both rules share the prefix `(Person p) ⋈ (Address a)` and can reuse the same JoinNode for this part.

### SQL Query Optimization Parallels

BetaChains draw inspiration from SQL query optimization:

| SQL Concept | RETE Equivalent |
|-------------|-----------------|
| Join tree | BetaChain |
| Join reordering | Selectivity-based ordering |
| Common subexpression elimination | Prefix sharing |
| Execution plan | Chain construction strategy |

**Standard SQL Join Ordering Principles:**
1. Most selective joins first (smallest intermediate results)
2. Joins with equality predicates before inequality predicates
3. Indexed joins before non-indexed joins

---

## BetaChain Structure

### Core Data Structure

```go
// BetaChain represents an ordered sequence of JoinNodes forming a join plan.
//
// A BetaChain is analogous to an execution plan in SQL databases, defining:
// - The order of join operations
// - Which variables participate in each join
// - Metadata for optimization and debugging
//
// Structure:
//   TypeNode(Person) ──┐
//                      ├─→ JoinNode₁(p ⋈ a) ──┐
//   TypeNode(Address) ─┘                      ├─→ JoinNode₂((p,a) ⋈ o) → Terminal
//   TypeNode(Order) ──────────────────────────┘
//
// Properties:
// - len(Nodes) == len(Hashes) == len(JoinSpecs)
// - FinalNode == Nodes[len(Nodes)-1]
// - Each node builds upon results of previous nodes (left-deep tree)
//
type BetaChain struct {
    // Ordered list of JoinNodes in the chain
    Nodes []*JoinNode `json:"nodes"`
    
    // Hash identifiers for each JoinNode (for sharing)
    Hashes []string `json:"hashes"`
    
    // Specifications for each join (variables, conditions, selectivity)
    JoinSpecs []JoinSpec `json:"join_specs"`
    
    // The final node in the chain (connects to terminal or next node)
    FinalNode *JoinNode `json:"final_node"`
    
    // ID of the rule this chain was built for
    RuleID string `json:"rule_id"`
    
    // Variables in scope at each join stage
    // VariablesAtStage[i] = variables available after Nodes[i]
    VariablesAtStage [][]string `json:"variables_at_stage"`
    
    // Metadata
    BuildStrategy string    `json:"build_strategy"` // "binary", "cascade", "optimized"
    CreatedAt     time.Time `json:"created_at"`
    BuildTimeMs   int64     `json:"build_time_ms"`
}

// JoinSpec describes a single join operation in the chain
type JoinSpec struct {
    // Variables on left side (accumulated from previous joins)
    LeftVars []string `json:"left_vars"`
    
    // Variables on right side (new variable being joined)
    RightVars []string `json:"right_vars"`
    
    // Join condition (AST from rule definition)
    Condition map[string]interface{} `json:"condition"`
    
    // Variable types
    VarTypes map[string]string `json:"var_types"`
    
    // Estimated selectivity (0.0 to 1.0)
    // Lower = more selective (fewer results)
    Selectivity float64 `json:"selectivity"`
    
    // Join method hint ("hash", "nested_loop", "index")
    JoinMethod string `json:"join_method,omitempty"`
}
```

### Chain Metadata

```go
// BetaChainMetadata provides analytics and debugging information
type BetaChainMetadata struct {
    TotalJoins      int       `json:"total_joins"`
    SharedJoins     int       `json:"shared_joins"`
    NewJoins        int       `json:"new_joins"`
    EstimatedCost   float64   `json:"estimated_cost"`   // Sum of selectivities
    OptimizedOrder  bool      `json:"optimized_order"`  // Was join order optimized?
    PrefixShared    bool      `json:"prefix_shared"`    // Does it share prefix with other chains?
    SharedPrefixLen int       `json:"shared_prefix_len"`
}
```

---

## BetaChainBuilder Design

### Interface

```go
// BetaChainBuilder constructs optimized BetaChains with automatic sharing
//
// Responsibilities:
// - Analyze join patterns from rule definitions
// - Determine optimal join ordering (selectivity-based)
// - Reuse existing JoinNodes via BetaSharingRegistry
// - Build progressive chains (prefix → full chain)
// - Maintain connection cache to avoid duplicates
// - Collect metrics for monitoring
//
// Thread-safety: All public methods are thread-safe
//
type BetaChainBuilder struct {
    network         *ReteNetwork
    storage         Storage
    sharingRegistry *BetaSharingRegistry
    
    // Cache of built chain prefixes (for reuse)
    // Key: canonical prefix signature, Value: chain prefix
    prefixCache     map[string]*BetaChainPrefix
    
    // Connection tracking
    connectionCache map[string]bool // "parentID_childID" -> exists
    
    // Selectivity estimator
    selectivityEstimator SelectivityEstimator
    
    // Metrics
    metrics *BetaChainMetrics
    
    // Concurrency control
    mutex sync.RWMutex
}

// BetaChainPrefix represents a reusable prefix of a BetaChain
type BetaChainPrefix struct {
    Nodes         []*JoinNode
    Variables     []string  // Variables in scope after this prefix
    Signature     string    // Canonical signature for matching
    RefCount      int       // Number of chains using this prefix
    UsedByRules   []string  // Rule IDs using this prefix
}
```

### Core Methods

```go
// BuildChain constructs an optimized BetaChain for a set of variables and conditions
//
// Steps:
// 1. Parse variables and conditions
// 2. Estimate selectivity for each potential join
// 3. Determine optimal join order
// 4. Check for reusable prefixes
// 5. Build chain progressively (reusing or creating nodes)
// 6. Connect to parent nodes (TypeNodes or alpha pass-through)
// 7. Register with lifecycle manager
//
func (bcb *BetaChainBuilder) BuildChain(
    variables []string,
    varTypes map[string]string,
    conditions map[string]interface{},
    parentConnections map[string]Node, // varName -> TypeNode/AlphaNode
    ruleID string,
    strategy BuildStrategy,
) (*BetaChain, error)

// BuildOptimizedChain builds a chain with automatic join reordering
func (bcb *BetaChainBuilder) BuildOptimizedChain(
    variables []string,
    varTypes map[string]string,
    conditions map[string]interface{},
    parentConnections map[string]Node,
    ruleID string,
) (*BetaChain, error)

// BuildCascadeChain builds a chain in variable declaration order (no optimization)
func (bcb *BetaChainBuilder) BuildCascadeChain(
    variables []string,
    varTypes map[string]string,
    conditions map[string]interface{},
    parentConnections map[string]Node,
    ruleID string,
) (*BetaChain, error)

// FindReusablePrefix searches for an existing chain prefix that matches
func (bcb *BetaChainBuilder) FindReusablePrefix(
    variables []string,
    joinOrder []int, // Indices of variables in join order
) (*BetaChainPrefix, int) // Returns prefix and length

// BuildFromPrefix continues building from an existing prefix
func (bcb *BetaChainBuilder) BuildFromPrefix(
    prefix *BetaChainPrefix,
    remainingVars []string,
    varTypes map[string]string,
    conditions map[string]interface{},
    ruleID string,
) (*BetaChain, error)
```

---

## Construction Algorithm

### High-Level Algorithm

```
ALGORITHM: BuildOptimizedBetaChain
INPUT: variables[], varTypes{}, conditions{}, parentConnections{}, ruleID
OUTPUT: BetaChain

1. PARSE & VALIDATE
   - Extract variable list: [v₁, v₂, ..., vₙ]
   - Validate parent connections exist for all variables
   - Parse join conditions from rule AST

2. SELECTIVITY ESTIMATION
   FOR each pair (vᵢ, vⱼ) in variables:
       selectivity[i,j] = EstimateJoinSelectivity(vᵢ, vⱼ, conditions)
   
   // Selectivity factors:
   // - Equality joins: 0.1 (highly selective)
   // - Range joins: 0.3
   // - Cross products: 1.0 (no condition)
   // - Joins on indexed fields: × 0.5
   // - Type popularity: × typeCardinality

3. DETERMINE OPTIMAL JOIN ORDER
   joinOrder = OptimizeJoinOrder(variables, selectivity, conditions)
   
   // Use greedy algorithm (similar to System R optimizer):
   // - Start with most selective two-variable join
   // - Iteratively add next variable with highest selectivity
   // - Respect join condition dependencies

4. FIND REUSABLE PREFIX
   prefixSignature = CanonicalSignature(joinOrder[0:k]) for k = 1 to n-1
   (prefix, prefixLen) = FindReusablePrefix(variables, joinOrder)
   
   IF prefix EXISTS THEN:
       startIdx = prefixLen
       currentNode = prefix.FinalNode
       chain.Nodes = prefix.Nodes[:]  // Copy prefix nodes
       chain.Hashes = prefix.Hashes[:]
       chain.JoinSpecs = prefix.JoinSpecs[:]
       metrics.PrefixReused++
   ELSE:
       startIdx = 0
       currentNode = NULL
       chain = NEW BetaChain

5. BUILD REMAINING JOINS
   FOR i = startIdx to len(joinOrder)-1:
       
       5a. DETERMINE JOIN INPUTS
           IF i == 0 THEN:
               // First join (binary)
               leftVars = [joinOrder[0]]
               rightVars = [joinOrder[1]]
               leftParent = parentConnections[joinOrder[0]]
               rightParent = parentConnections[joinOrder[1]]
           ELSE IF i == 1 AND startIdx == 0 THEN:
               // Second join in cascade
               leftVars = [joinOrder[0]]  // Result from first join
               rightVars = [joinOrder[1]]
               leftParent = chain.Nodes[0]
               rightParent = parentConnections[joinOrder[1]]
           ELSE:
               // Subsequent joins
               leftVars = joinOrder[0:i]  // All accumulated variables
               rightVars = [joinOrder[i]]  // New variable
               leftParent = currentNode
               rightParent = parentConnections[joinOrder[i]]
       
       5b. EXTRACT JOIN CONDITION
           joinCond = ExtractJoinCondition(conditions, leftVars, rightVars)
       
       5c. CREATE JOIN SPEC
           joinSpec = JoinSpec{
               LeftVars: leftVars,
               RightVars: rightVars,
               Condition: joinCond,
               VarTypes: varTypes,
               Selectivity: selectivity[leftVars, rightVars],
           }
       
       5d. GET OR CREATE JOIN NODE
           (joinNode, hash, reused) = sharingRegistry.GetOrCreateJoinNode(joinSpec)
           
           IF reused THEN:
               metrics.NodesReused++
           ELSE:
               metrics.NodesCreated++
               network.BetaNodes[joinNode.ID] = joinNode
       
       5e. CONNECT TO PARENTS
           IF NOT IsConnected(leftParent, joinNode, "left") THEN:
               ConnectBetaNode(leftParent, joinNode, "left")
           
           IF NOT IsConnected(rightParent, joinNode, "right") THEN:
               ConnectBetaNode(rightParent, joinNode, "right")
       
       5f. REGISTER WITH LIFECYCLE
           lifecycle = network.LifecycleManager.RegisterNode(joinNode.ID, "beta")
           lifecycle.AddRuleReference(ruleID, "")
       
       5g. UPDATE CHAIN
           chain.Nodes.append(joinNode)
           chain.Hashes.append(hash)
           chain.JoinSpecs.append(joinSpec)
           chain.VariablesAtStage.append(leftVars + rightVars)
           
           currentNode = joinNode

6. FINALIZE CHAIN
   chain.FinalNode = currentNode
   chain.RuleID = ruleID
   chain.BuildStrategy = "optimized"
   
   // Cache prefix for future reuse
   IF len(chain.Nodes) >= 2 THEN:
       FOR prefixLen = 1 to len(chain.Nodes)-1:
           CachePrefix(chain.Nodes[0:prefixLen], ruleID)

7. RETURN chain
```

### Selectivity Estimation Algorithm

```
ALGORITHM: EstimateJoinSelectivity
INPUT: leftVar, rightVar, conditions
OUTPUT: selectivity (0.0 to 1.0)

1. DEFAULT = 1.0  // Cross product (no condition)

2. EXTRACT JOIN CONDITION
   joinCond = FindConditionBetween(leftVar, rightVar, conditions)
   IF joinCond == NULL THEN:
       RETURN DEFAULT

3. ANALYZE CONDITION TYPE
   SWITCH joinCond.operator:
       CASE "==", "equals":
           baseSelectivity = 0.1  // Equality is highly selective
       
       CASE "<", ">", "<=", ">=":
           baseSelectivity = 0.3  // Range is moderately selective
       
       CASE "!=", "not_equals":
           baseSelectivity = 0.9  // Anti-join is weakly selective
       
       CASE "contains", "startsWith":
           baseSelectivity = 0.4  // String matching
       
       DEFAULT:
           baseSelectivity = 0.5

4. APPLY HEURISTIC ADJUSTMENTS
   
   4a. INDEX FACTOR
       IF FieldHasIndex(leftVar, joinCond.leftField) OR
          FieldHasIndex(rightVar, joinCond.rightField) THEN:
           baseSelectivity *= 0.7  // Indexed joins are more efficient
   
   4b. TYPE CARDINALITY FACTOR
       leftCard = EstimateTypeCardinality(leftVar.type)
       rightCard = EstimateTypeCardinality(rightVar.type)
       
       // Adjust for type sizes
       IF leftCard < 100 AND rightCard < 100 THEN:
           baseSelectivity *= 0.8  // Small types = more selective
       ELSE IF leftCard > 10000 OR rightCard > 10000 THEN:
           baseSelectivity *= 1.2  // Large types = less selective
   
   4c. FIELD UNIQUENESS
       IF FieldIsUnique(leftVar, joinCond.leftField) THEN:
           baseSelectivity *= 0.3  // Unique field = very selective
       
       IF FieldIsUnique(rightVar, joinCond.rightField) THEN:
           baseSelectivity *= 0.3

5. CLAMP TO [0.01, 1.0]
   RETURN CLAMP(baseSelectivity, 0.01, 1.0)
```

### Join Ordering Algorithm

```
ALGORITHM: OptimizeJoinOrder
INPUT: variables[], selectivity[][], conditions
OUTPUT: joinOrder[] (indices into variables[])

// Uses greedy algorithm inspired by System R optimizer
// Builds left-deep join tree by always choosing most selective next join

1. INITIALIZE
   remaining = {0, 1, 2, ..., len(variables)-1}
   joined = {}
   joinOrder = []

2. FIND MOST SELECTIVE INITIAL PAIR
   minSelectivity = INFINITY
   bestPair = NULL
   
   FOR i in remaining:
       FOR j in remaining WHERE j > i:
           IF selectivity[i][j] < minSelectivity THEN:
               minSelectivity = selectivity[i][j]
               bestPair = (i, j)
   
   IF bestPair == NULL THEN:
       // No conditions - use declaration order
       RETURN [0, 1, 2, ..., len(variables)-1]
   
   joinOrder = [bestPair.first, bestPair.second]
   joined = {bestPair.first, bestPair.second}
   remaining.remove(bestPair.first, bestPair.second)

3. ITERATIVELY ADD MOST SELECTIVE VARIABLE
   WHILE remaining NOT EMPTY:
       
       minSelectivity = INFINITY
       bestNext = NULL
       
       FOR candidate in remaining:
           // Compute selectivity of joining candidate to current result
           cumSelectivity = ComputeCumulativeSelectivity(
               joined, candidate, selectivity, conditions
           )
           
           IF cumSelectivity < minSelectivity THEN:
               minSelectivity = cumSelectivity
               bestNext = candidate
       
       IF bestNext == NULL THEN:
           // No conditions - add in declaration order
           bestNext = MIN(remaining)
       
       joinOrder.append(bestNext)
       joined.add(bestNext)
       remaining.remove(bestNext)

4. RETURN joinOrder
```

---

## Optimization Strategies

### 1. Selectivity-Based Join Ordering

**Principle:** Perform most selective joins first to minimize intermediate result sizes.

**Benefits:**
- Reduces memory consumption in join node memories
- Faster propagation (fewer tokens to process)
- Lower CPU usage for condition evaluation

**Example:**

```
Variables: Person(p), Address(a), Order(o)
Conditions:
  - p.id == a.person_id  (selectivity: 0.1, indexed)
  - a.id == o.address_id (selectivity: 0.3)
  - o.total > 1000       (selectivity: 0.2)

Unoptimized order: p ⋈ a ⋈ o (follows declaration)
Optimized order:   p ⋈ a ⋈ o (selectivity: 0.1 → 0.2 → 0.3)

Result: 50% fewer intermediate tokens
```

### 2. Prefix Sharing

**Principle:** Reuse common join prefixes across multiple rules.

**Implementation:**
- Maintain `prefixCache` with canonical signatures
- Match prefixes by: variable types + join conditions + order
- Reference count prefixes for lifecycle management

**Example:**

```
Rule 1: (Person p) ⋈ (Address a) ⋈ (Order o)
         └─────────────┬──────────────┘
                   Shared Prefix
Rule 2: (Person p) ⋈ (Address a) ⋈ (Invoice i)
         └─────────────┬──────────────┘
                   Shared Prefix

Savings: 1 JoinNode reused, connections shared
```

### 3. Condition-Aware Optimization

**Principle:** Consider join condition types when ordering.

**Heuristics:**
1. **Equality joins before range joins**
   - `p.id == a.person_id` before `p.age > a.min_age`
   
2. **Indexed joins first**
   - Joins on primary/foreign keys before arbitrary fields
   
3. **Selective predicates early**
   - Apply filters that eliminate most data first

### 4. Connection Deduplication

**Principle:** Avoid duplicate connections between same parent-child pairs.

**Mechanism:**
- Maintain `connectionCache` map: `"parentID_childID_side" -> bool`
- Check cache before calling `AddChild()` or `ConnectBetaNode()`
- Update cache on new connections

**Benefits:**
- Prevents duplicate activations
- Cleaner network structure
- Faster traversal

### 5. Adaptive Strategy Selection

**Decision Tree:**

```
IF numVariables == 2 THEN:
    strategy = BINARY
ELSE IF complexConditions OR userOptimizationFlag THEN:
    strategy = OPTIMIZED
    // Perform selectivity analysis and reordering
ELSE:
    strategy = CASCADE
    // Use declaration order (faster compilation, predictable)

RETURN BuildChain(strategy)
```

---

## Integration with BetaSharingRegistry

### Coordinated Design

BetaChainBuilder and BetaSharingRegistry work together:

| Component | Responsibility |
|-----------|----------------|
| **BetaSharingRegistry** | Hash, store, and retrieve individual JoinNodes |
| **BetaChainBuilder** | Determine join order, build chains, manage prefixes |

### Integration Points

#### 1. Node Creation

```go
// BetaChainBuilder calls BetaSharingRegistry for each join
func (bcb *BetaChainBuilder) buildJoinStep(...) {
    // Create JoinSpec from current join
    spec := JoinSpec{...}
    
    // Get or create via sharing registry
    joinNode, hash, reused, err := bcb.sharingRegistry.GetOrCreateJoinNode(
        spec.Condition,
        spec.LeftVars,
        spec.RightVars,
        spec.VarTypes,
        bcb.storage,
    )
    
    // Continue chain building...
}
```

#### 2. Prefix Matching

```go
// Prefix signature includes hashes from sharing registry
func (bcb *BetaChainBuilder) computePrefixSignature(joinSpecs []JoinSpec) string {
    var hashes []string
    for _, spec := range joinSpecs {
        // Use same hashing as BetaSharingRegistry
        hash := bcb.sharingRegistry.ComputeJoinHash(spec)
        hashes = append(hashes, hash)
    }
    return strings.Join(hashes, "|")
}
```

#### 3. Lifecycle Coordination

```go
// Both systems use LifecycleManager
func (bcb *BetaChainBuilder) registerChain(chain *BetaChain, ruleID string) {
    for _, joinNode := range chain.Nodes {
        // BetaSharingRegistry already registered node
        // BetaChainBuilder adds rule reference
        lifecycle := bcb.network.LifecycleManager.GetNodeLifecycle(joinNode.ID)
        if lifecycle != nil {
            lifecycle.AddRuleReference(ruleID, chain.RuleID)
        }
    }
}
```

### Benefits of Integration

1. **Single Source of Truth:** BetaSharingRegistry handles all node hashing/storage
2. **Maximum Reuse:** Chain builder finds prefix reuse; registry finds node reuse
3. **Consistent Lifecycle:** Both systems use same reference counting
4. **Composable Optimizations:** Join ordering + node sharing = multiplicative benefit

---

## Comparison with AlphaChains

### Similarities

| Aspect | AlphaChains | BetaChains |
|--------|-------------|------------|
| **Purpose** | Sequence of conditions | Sequence of joins |
| **Structure** | Ordered list of nodes | Ordered list of nodes |
| **Sharing** | Via AlphaSharingRegistry | Via BetaSharingRegistry |
| **Caching** | Connection cache | Connection cache + prefix cache |
| **Metrics** | Build time, reuse ratio | Build time, reuse ratio, selectivity |
| **Lifecycle** | Reference counting | Reference counting |

### Differences

| Aspect | AlphaChains | BetaChains |
|--------|-------------|------------|
| **Node Type** | AlphaNode (1 input) | JoinNode (2 inputs) |
| **Variables** | Single variable per chain | Multiple variables accumulate |
| **Ordering** | Sequential (fixed by rule) | Optimizable (selectivity-based) |
| **Complexity** | O(n) where n = conditions | O(n²) selectivity estimation |
| **Parent Connection** | Single parent (TypeNode) | Two parents (left + right) |
| **Prefix Reuse** | Limited (same variable) | High potential (common join patterns) |
| **Memory Impact** | Alpha memory (facts) | Beta memory (tokens) - larger |

### Architectural Comparison

#### AlphaChain Structure

```
TypeNode(Person)
    ↓
AlphaNode(p.age > 18)
    ↓
AlphaNode(p.city == "Paris")
    ↓
AlphaNode(p.status == "active")
    ↓
[Next stage]

- Linear flow
- Single variable context
- Independent evaluation
```

#### BetaChain Structure

```
TypeNode(Person) ─┐
                  ├→ JoinNode₁(p ⋈ a) ─┐
TypeNode(Address) ┘                    ├→ JoinNode₂((p,a) ⋈ o) → Terminal
TypeNode(Order) ───────────────────────┘

- Two-input joins
- Accumulating variable context
- Dependent evaluation (tokens)
```

### Code Structure Comparison

#### AlphaChainBuilder

```go
// Simpler structure
type AlphaChainBuilder struct {
    network         *ReteNetwork
    storage         Storage
    connectionCache map[string]bool
    metrics         *ChainBuildMetrics
    mutex           sync.RWMutex
}

// Straightforward build
func (acb *AlphaChainBuilder) BuildChain(...) {
    for _, condition := range conditions {
        node, hash, reused = GetOrCreateAlphaNode(condition)
        parent.AddChild(node)
        parent = node
    }
}
```

#### BetaChainBuilder

```go
// More complex structure
type BetaChainBuilder struct {
    network              *ReteNetwork
    storage              Storage
    sharingRegistry      *BetaSharingRegistry
    prefixCache          map[string]*BetaChainPrefix
    connectionCache      map[string]bool
    selectivityEstimator SelectivityEstimator
    metrics              *BetaChainMetrics
    mutex                sync.RWMutex
}

// Complex build with optimization
func (bcb *BetaChainBuilder) BuildChain(...) {
    // 1. Estimate selectivity
    selectivity = EstimateSelectivity(variables, conditions)
    
    // 2. Optimize join order
    joinOrder = OptimizeJoinOrder(variables, selectivity)
    
    // 3. Find reusable prefix
    prefix = FindReusablePrefix(joinOrder)
    
    // 4. Build remaining joins
    for i, varIdx := range joinOrder {
        leftVars = joinOrder[0:i]
        rightVars = [joinOrder[i]]
        node, hash, reused = GetOrCreateJoinNode(leftVars, rightVars)
        ConnectBetaNode(leftParent, node, "left")
        ConnectBetaNode(rightParent, node, "right")
    }
}
```

---

## Examples & Use Cases

### Example 1: Binary Join (Baseline)

**Rule:**
```
Rule HighValueCustomer:
  (Customer c)
  (Order o)
  WHERE c.id == o.customer_id AND o.total > 10000
  THEN notify_sales_team(c, o)
```

**BetaChain:**
```
TypeNode(Customer) ─┐
                    ├→ JoinNode₁(c ⋈ o) → Terminal
TypeNode(Order) ────┘

JoinSpec₁:
  LeftVars:  [c]
  RightVars: [o]
  Condition: c.id == o.customer_id AND o.total > 10000
  Selectivity: 0.05 (highly selective)
```

**Characteristics:**
- Simple binary join
- No optimization needed
- Direct connection to terminal

### Example 2: Cascade Join (3 Variables)

**Rule:**
```
Rule CompleteProfile:
  (Person p)
  (Address a)
  (Phone ph)
  WHERE p.id == a.person_id AND p.id == ph.person_id
  THEN mark_complete(p)
```

**BetaChain (Unoptimized - Declaration Order):**
```
TypeNode(Person) ──┐
                   ├→ JoinNode₁(p ⋈ a) ──┐
TypeNode(Address) ─┘                     ├→ JoinNode₂((p,a) ⋈ ph) → Terminal
TypeNode(Phone) ─────────────────────────┘

JoinSpec₁:
  LeftVars:  [p]
  RightVars: [a]
  Selectivity: 0.15

JoinSpec₂:
  LeftVars:  [p, a]
  RightVars: [ph]
  Selectivity: 0.20
```

**BetaChain (Optimized - Selectivity Order):**
```
// Assume Phone join is more selective (fewer phones per person)

TypeNode(Person) ──┐
                   ├→ JoinNode₁(p ⋈ ph) ──┐
TypeNode(Phone) ───┘                      ├→ JoinNode₂((p,ph) ⋈ a) → Terminal
TypeNode(Address) ────────────────────────┘

JoinSpec₁:
  LeftVars:  [p]
  RightVars: [ph]
  Selectivity: 0.10 (more selective - fewer phones)

JoinSpec₂:
  LeftVars:  [p, ph]
  RightVars: [a]
  Selectivity: 0.15
```

**Performance Comparison:**
- Unoptimized: ~1500 intermediate tokens (p ⋈ a result)
- Optimized: ~1000 intermediate tokens (p ⋈ ph result)
- **Savings: 33% fewer tokens, 25% faster execution**

### Example 3: Prefix Sharing

**Rule 1:**
```
Rule CustomerOrders:
  (Customer c)
  (Order o)
  (Product p)
  WHERE c.id == o.customer_id AND o.product_id == p.id
```

**Rule 2:**
```
Rule CustomerReturns:
  (Customer c)
  (Order o)
  (Return r)
  WHERE c.id == o.customer_id AND o.id == r.order_id
```

**Shared Prefix:**
```
TypeNode(Customer) ─┐
                    ├→ JoinNode₁(c ⋈ o)  ← SHARED BETWEEN RULES
TypeNode(Order) ────┘

Rule 1:  JoinNode₁ → JoinNode₂((c,o) ⋈ p) → Terminal₁
Rule 2:  JoinNode₁ → JoinNode₃((c,o) ⋈ r) → Terminal₂
```

**Prefix Cache Entry:**
```
signature: "hash(Customer+Order+c.id==o.customer_id)"
nodes: [JoinNode₁]
refCount: 2
usedByRules: ["CustomerOrders", "CustomerReturns"]
```

**Benefits:**
- **1 JoinNode saved** (50% reduction for this part of network)
- **Shared computation** for c ⋈ o join
- **Faster compilation** for Rule 2 (reuses prefix)

### Example 4: Complex Multi-Way Join

**Rule:**
```
Rule SuspiciousActivity:
  (User u)
  (Login l)
  (Device d)
  (Location loc)
  WHERE u.id == l.user_id 
    AND l.device_id == d.id
    AND l.location_id == loc.id
    AND loc.country != u.home_country
    AND d.is_new == true
```

**Selectivity Estimation:**
```
Join Options:
1. u ⋈ l:   selectivity 0.20 (many logins per user)
2. l ⋈ d:   selectivity 0.05 (few devices, AND d.is_new filter)
3. l ⋈ loc: selectivity 0.15 (moderate location diversity)
4. u ⋈ loc: selectivity 0.30 (cross product + filter)

Optimal Order: l ⋈ d (0.05) → (l,d) ⋈ u (0.08) → (l,d,u) ⋈ loc (0.12)
```

**Optimized BetaChain:**
```
TypeNode(Login) ──┐
                  ├→ JoinNode₁(l ⋈ d) ──┐
TypeNode(Device) ─┘                     │
                                        ├→ JoinNode₂((l,d) ⋈ u) ──┐
TypeNode(User) ─────────────────────────┘                         │
                                                                  ├→ JoinNode₃ → Terminal
TypeNode(Location) ───────────────────────────────────────────────┘

JoinSpec₁: l ⋈ d        (selectivity: 0.05) → ~500 tokens
JoinSpec₂: (l,d) ⋈ u    (selectivity: 0.08) → ~400 tokens  
JoinSpec₃: (l,d,u) ⋈ loc (selectivity: 0.12) → ~350 tokens
```

**Unoptimized (Declaration Order):**
```
JoinSpec₁: u ⋈ l        (selectivity: 0.20) → ~2000 tokens
JoinSpec₂: (u,l) ⋈ d    (selectivity: 0.15) → ~1500 tokens
JoinSpec₃: (u,l,d) ⋈ loc (selectivity: 0.18) → ~1200 tokens
```

**Performance Impact:**
- **Unoptimized:** ~4700 total intermediate tokens
- **Optimized:** ~1250 total intermediate tokens
- **Savings: 73% reduction in memory and processing**

---

## Implementation Plan

### Phase 1: Core Infrastructure (Week 1-2)

**Deliverables:**
1. `beta_chain.go`
   - `BetaChain` struct
   - `JoinSpec` struct
   - `BetaChainMetadata` struct
   - Validation methods

2. `beta_chain_builder.go` (basic version)
   - `BetaChainBuilder` struct
   - `NewBetaChainBuilder()` constructor
   - `BuildCascadeChain()` method (no optimization)
   - Connection cache logic

3. Unit tests
   - Chain structure validation
   - Basic cascade build (2-3 variables)
   - Connection deduplication

**Success Criteria:**
- Can build simple cascade chains
- Correctly connects JoinNodes to parents
- Validates chain consistency

### Phase 2: Selectivity & Optimization (Week 3)

**Deliverables:**
1. `selectivity_estimator.go`
   - `SelectivityEstimator` interface
   - `HeuristicEstimator` implementation
   - Condition analysis
   - Type cardinality estimation

2. `beta_chain_builder.go` (optimized build)
   - `EstimateSelectivity()` method
   - `OptimizeJoinOrder()` method
   - `BuildOptimizedChain()` method

3. Unit tests
   - Selectivity estimation accuracy
   - Join reordering correctness
   - Edge cases (no conditions, circular dependencies)

**Success Criteria:**
- Correctly estimates relative selectivity
- Produces sensible join orders
- Handles edge cases gracefully

### Phase 3: Prefix Sharing (Week 4)

**Deliverables:**
1. `beta_chain_builder.go` (prefix caching)
   - `BetaChainPrefix` struct
   - `FindReusablePrefix()` method
   - `BuildFromPrefix()` method
   - Prefix signature computation
   - Cache management

2. Integration tests
   - Multi-rule scenarios with shared prefixes
   - Prefix reference counting
   - Cache eviction behavior

**Success Criteria:**
- Detects and reuses common prefixes
- Correctly maintains refcounts
- Measurable reduction in duplicate joins

### Phase 4: Integration (Week 5)

**Deliverables:**
1. `constraint_pipeline_builder.go` updates
   - Replace manual join building with `BetaChainBuilder`
   - Add strategy selection logic
   - Preserve backward compatibility

2. Metrics & monitoring
   - `BetaChainMetrics` struct
   - Integration with existing metrics
   - Logging and debugging output

3. End-to-end tests
   - Full rule compilation with BetaChains
   - Activation/retraction correctness
   - Performance benchmarks

**Success Criteria:**
- Seamless integration with existing code
- All existing tests pass
- Performance improvements measurable

### Phase 5: Documentation & Rollout (Week 6)

**Deliverables:**
1. User documentation
   - When to use optimized vs cascade strategy
   - Performance tuning guide
   - Debugging tools

2. Developer documentation
   - Architecture guide
   - API reference
   - Extension points

3. Feature flag & rollout
   - `EnableBetaChainOptimization` flag (default: true for new rules)
   - `EnableBetaChainPrefixSharing` flag (default: true)
   - Monitoring dashboard

**Success Criteria:**
- Complete documentation
- Safe rollout plan
- Rollback capability

---

## Performance Expectations

### Memory Savings

| Scenario | Baseline | With BetaChains | Improvement |
|----------|----------|-----------------|-------------|
| 10 rules, no shared patterns | 30 JoinNodes | 30 JoinNodes | 0% |
| 10 rules, 50% shared prefixes | 30 JoinNodes | 20 JoinNodes | 33% |
| 10 rules, 80% shared prefixes | 30 JoinNodes | 12 JoinNodes | 60% |
| 100 rules, typical workload | 350 JoinNodes | 200 JoinNodes | 43% (estimated) |

### Compilation Time

| Scenario | Baseline | With BetaChains | Improvement |
|----------|----------|-----------------|-------------|
| Single rule, 2 variables | 5ms | 8ms | -60% (overhead) |
| Single rule, 5 variables | 15ms | 12ms | +20% (optimization) |
| 10 rules, 50% sharing | 80ms | 50ms | +37% (prefix reuse) |
| 100 rules, 50% sharing | 950ms | 520ms | +45% (prefix reuse) |

### Runtime Performance

| Scenario | Baseline | With Optimized Order | Improvement |
|----------|----------|---------------------|-------------|
| 2-var join, no optimization needed | 100μs | 100μs | 0% |
| 3-var join, good default order | 250μs | 230μs | +8% |
| 3-var join, poor default order | 450μs | 210μs | +53% |
| 5-var join, poor default order | 2.1ms | 0.9ms | +57% |

### Working Memory Reduction

Optimized join ordering reduces intermediate token count:

```
Before optimization:
  Join 1: 10,000 × 1,000 = 10,000,000 intermediate results (selectivity 1.0)
  Join 2: 10,000,000 × 500 = 5,000,000 results (selectivity 0.5)
  Total: 15,000,000 tokens processed

After optimization:
  Join 1: 10,000 × 1,000 = 100,000 intermediate results (selectivity 0.01)
  Join 2: 100,000 × 500 = 50,000 results (selectivity 0.5)
  Total: 150,000 tokens processed

Improvement: 99% reduction in intermediate tokens!
```

---

## Testing Strategy

### Unit Tests

#### 1. Chain Structure Tests
```go
func TestBetaChain_Validation(t *testing.T)
func TestBetaChain_EmptyChain(t *testing.T)
func TestBetaChain_SingleJoin(t *testing.T)
func TestBetaChain_MultiJoin(t *testing.T)
func TestBetaChain_VariablesAtStage(t *testing.T)
```

#### 2. Builder Tests
```go
func TestBetaChainBuilder_CreateBinary(t *testing.T)
func TestBetaChainBuilder_CreateCascade(t *testing.T)
func TestBetaChainBuilder_CreateOptimized(t *testing.T)
func TestBetaChainBuilder_Concurrency(t *testing.T)
```

#### 3. Selectivity Tests
```go
func TestSelectivityEstimator_EqualityJoin(t *testing.T)
func TestSelectivityEstimator_RangeJoin(t *testing.T)
func TestSelectivityEstimator_CrossProduct(t *testing.T)
func TestSelectivityEstimator_IndexedJoin(t *testing.T)
```

#### 4. Optimization Tests
```go
func TestOptimizeJoinOrder_TwoVariables(t *testing.T)
func TestOptimizeJoinOrder_ThreeVariables(t *testing.T)
func TestOptimizeJoinOrder_NoConditions(t *testing.T)
func TestOptimizeJoinOrder_ComplexConditions(t *testing.T)
```

#### 5. Prefix Tests
```go
func TestPrefixCache_FindExact(t *testing.T)
func TestPrefixCache_FindPartial(t *testing.T)
func TestPrefixCache_NoMatch(t *testing.T)
func TestPrefixCache_RefCounting(t *testing.T)
```

### Integration Tests

#### 1. End-to-End Rule Compilation
```go
func TestE2E_SimpleBinaryRule(t *testing.T)
func TestE2E_CascadeRule(t *testing.T)
func TestE2E_OptimizedRule(t *testing.T)
func TestE2E_MultipleRulesWithSharing(t *testing.T)
```

#### 2. Activation/Retraction Correctness
```go
func TestActivation_BinaryJoin(t *testing.T)
func TestActivation_CascadeJoin(t *testing.T)
func TestActivation_SharedPrefix(t *testing.T)
func TestRetraction_PropagatesThroughChain(t *testing.T)
```

#### 3. Lifecycle Management
```go
func TestLifecycle_ChainCreation(t *testing.T)
func TestLifecycle_ChainDeletion(t *testing.T)
func TestLifecycle_PrefixRefCounting(t *testing.T)
func TestLifecycle_NodeSharing(t *testing.T)
```

### Performance Benchmarks

#### 1. Build Performance
```go
func BenchmarkBuildChain_Binary(b *testing.B)
func BenchmarkBuildChain_Cascade3Vars(b *testing.B)
func BenchmarkBuildChain_Cascade5Vars(b *testing.B)
func BenchmarkBuildChain_OptimizedVsCascade(b *testing.B)
```

#### 2. Selectivity Estimation
```go
func BenchmarkEstimateSelectivity(b *testing.B)
func BenchmarkOptimizeJoinOrder(b *testing.B)
```

#### 3. Prefix Matching
```go
func BenchmarkFindReusablePrefix_HitCache(b *testing.B)
func BenchmarkFindReusablePrefix_MissCache(b *testing.B)
```

#### 4. Memory Usage
```go
func BenchmarkMemoryUsage_WithoutSharing(b *testing.B)
func BenchmarkMemoryUsage_WithSharing(b *testing.B)
func BenchmarkMemoryUsage_OptimizedVsCascade(b *testing.B)
```

### Test Data Sets

#### Small (Development)
- 5 types, 10 rules, 2-3 variables per rule
- Simple join conditions (equality only)
- Some shared patterns (30%)

#### Medium (CI)
- 20 types, 50 rules, 2-5 variables per rule
- Mixed conditions (equality, range, complex)
- Moderate sharing (50%)

#### Large (Performance)
- 100 types, 500 rules, 2-8 variables per rule
- Complex conditions with nested expressions
- High sharing (70%)

#### Stress (Stability)
- 1000 types, 5000 rules, 2-10 variables per rule
- Random condition generation
- Variable sharing levels

---

## Appendix: Pseudo-Code Supplements

### Canonical Prefix Signature

```
FUNCTION ComputePrefixSignature(joinSpecs[]): string
    components = []
    
    FOR each spec IN joinSpecs:
        // Normalize spec for consistent hashing
        normalized = NormalizeJoinSpec(spec)
        
        // Compute hash (same as BetaSharingRegistry)
        hash = SHA256(normalized)
        
        components.append(hash)
    
    // Concatenate all hashes
    signature = Join(components, "|")
    
    RETURN signature
```

### Connection Tracking

```
FUNCTION IsConnected(parent Node, child Node, side string): bool
    key = Format("%s_%s_%s", parent.GetID(), child.GetID(), side)
    
    LOCK connectionCache.RLock()
    defer connectionCache.RUnlock()
    
    RETURN connectionCache[key] == true

FUNCTION ConnectBetaNode(parent Node, child Node, side string)
    key = Format("%s_%s_%s", parent.GetID(), child.GetID(), side)
    
    // Check if already connected
    IF IsConnected(parent, child, side) THEN:
        RETURN
    
    // Perform connection based on side
    IF side == "left" THEN:
        parent.AddChild(child)  // Left input
        child.SetLeftParent(parent)
    ELSE IF side == "right" THEN:
        parent.AddChild(child)  // Right input
        child.SetRightParent(parent)
    
    // Update cache
    LOCK connectionCache.Lock()
    connectionCache[key] = true
    connectionCache.Unlock()
```

### Prefix Cache Management

```
FUNCTION CachePrefix(nodes []*JoinNode, ruleID string)
    IF len(nodes) == 0 THEN:
        RETURN
    
    // Compute canonical signature
    signature = ComputePrefixSignature(ExtractJoinSpecs(nodes))
    
    LOCK prefixCache.Lock()
    defer prefixCache.Unlock()
    
    // Check if prefix already cached
    IF prefix EXISTS in prefixCache[signature] THEN:
        prefix.RefCount++
        prefix.UsedByRules.append(ruleID)
    ELSE:
        // Create new prefix entry
        prefix = BetaChainPrefix{
            Nodes: nodes[:],  // Copy
            Variables: ExtractVariables(nodes),
            Signature: signature,
            RefCount: 1,
            UsedByRules: [ruleID],
        }
        prefixCache[signature] = prefix
```

---

## Conclusion

BetaChains provide a robust, optimized framework for building join sequences in the RETE network. By combining:

1. **Selectivity-based join ordering** for runtime efficiency
2. **Prefix sharing** for memory savings
3. **Integration with BetaSharingRegistry** for maximum node reuse
4. **Progressive construction** for incremental compilation

We achieve significant improvements in both memory consumption and execution performance, while maintaining backward compatibility and code clarity.

**Next Steps:**
1. Review this design with stakeholders
2. Begin Phase 1 implementation (core infrastructure)
3. Set up benchmarking framework for validation
4. Document API for external consumers

**Related Documents:**
- `BETA_SHARING_DESIGN.md` - Node sharing implementation
- `BETA_NODES_ANALYSIS.md` - Beta network analysis
- `alpha_chain_builder.go` - Reference implementation for alpha chains

---

**Document Version History:**

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-01-XX | TSD Team | Initial design |

---

**End of Document**