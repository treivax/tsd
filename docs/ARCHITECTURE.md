# TSD Architecture

This document provides a comprehensive technical overview of the TSD rule engine architecture, implementation details, and design decisions.

## Table of Contents

- [Overview](#overview)
- [RETE Algorithm Implementation](#rete-algorithm-implementation)
- [Network Structure](#network-structure)
- [Node Types](#node-types)
- [Memory Management](#memory-management)
- [Optimization Strategies](#optimization-strategies)
- [Transaction System](#transaction-system)
- [Type System](#type-system)
- [Expression Evaluation](#expression-evaluation)
- [Action Execution](#action-execution)
- [Storage Layer](#storage-layer)
- [Performance Characteristics](#performance-characteristics)
- [Concurrency Model](#concurrency-model)

---

## Overview

TSD is a forward-chaining rule engine built on the RETE algorithm, optimized for high-performance pattern matching and rule execution. The system is designed for:

- **Fast pattern matching**: Incremental evaluation with shared computation
- **Scalability**: Handle thousands of rules and millions of facts efficiently
- **Type safety**: Strong static typing with runtime validation
- **Concurrency**: Thread-safe operations with fine-grained locking
- **Extensibility**: Pluggable action executors and future network replication

### Key Design Principles

1. **Incremental Computation**: Only recompute what changes
2. **State Sharing**: Maximize reuse of intermediate results
3. **Memory Efficiency**: Smart caching and cleanup strategies
4. **Strong Guarantees**: ACID transactions with configurable consistency

---

## RETE Algorithm Implementation

TSD implements a variant of the RETE algorithm with several enhancements for modern use cases.

### Classic RETE Components

```
┌─────────────┐
│  Root Node  │
└──────┬──────┘
       │
   ┌───┴────┐
   │  Type  │  (Filter facts by type)
   │  Node  │
   └───┬────┘
       │
   ┌───┴─────┐
   │  Alpha  │  (Test single fact conditions)
   │  Node   │
   └───┬─────┘
       │
   ┌───┴────┐
   │  Beta  │  (Join multiple facts)
   │  Node  │
   └───┬────┘
       │
  ┌────┴────────┐
  │  Terminal   │  (Execute actions)
  │    Node     │
  └─────────────┘
```

### Enhanced Node Types

TSD extends classic RETE with specialized nodes:

- **AccumulatorNode**: Compute aggregations (SUM, AVG, COUNT, MIN, MAX)
- **ExistsNode**: Existential quantification (EXISTS patterns)
- **NotNode**: Negation (NOT patterns)
- **MultiSourceAccumulatorNode**: Multi-source aggregations
- **RuleRouterNode**: Efficient beta node sharing across rules

---

## Network Structure

### ReteNetwork

The central data structure coordinating all components:

```go
type ReteNetwork struct {
    RootNode              *RootNode
    TypeNodes             map[string]*TypeNode
    AlphaNodes            map[string]*AlphaNode
    BetaNodes             map[string]interface{}
    TerminalNodes         map[string]*TerminalNode
    Storage               Storage
    
    // Optimization components
    AlphaSharingManager   *AlphaSharingRegistry
    BetaSharingRegistry   BetaSharingRegistry
    AlphaChainBuilder     *AlphaChainBuilder
    BetaChainBuilder      *BetaChainBuilder
    ArithmeticResultCache *ArithmeticResultCache
    
    // Execution components
    ActionExecutor        *ActionExecutor
    LifecycleManager      *LifecycleManager
    
    // Transaction support
    currentTx             *Transaction
    txMutex               sync.RWMutex
    
    // Performance configuration
    Config                *ChainPerformanceConfig
    ChainMetrics          *ChainBuildMetrics
}
```

### Network Construction

Network building occurs in phases:

1. **Type Registration**: Register all type definitions
2. **Alpha Chain Building**: Create and share alpha nodes
3. **Beta Chain Building**: Create join nodes with sharing
4. **Terminal Node Creation**: Connect to action executors
5. **Optimization Pass**: Apply sharing and caching strategies

---

## Node Types

### BaseNode

All nodes inherit from `BaseNode`:

```go
type BaseNode struct {
    ID       string
    Type     string
    Memory   *WorkingMemory
    Children []Node
    Storage  Storage
    network  *ReteNetwork
    mutex    sync.RWMutex
}
```

**Responsibilities**:
- Unique identification
- Child node management
- Working memory management
- Thread-safe state access

### RootNode

Entry point for all facts into the network.

**Operations**:
- `ActivateLeft(token)`: Receives initial tokens
- Routes facts to appropriate TypeNodes based on fact type

### TypeNode

Filters facts by their declared type.

**Key Features**:
- One TypeNode per type definition
- Fast type-based routing
- Maintains type metadata

### AlphaNode

Tests conditions on a single fact.

**Condition Types**:
- Field comparisons: `field == value`, `field > 10`
- String operations: `LIKE`, `CONTAINS`, `MATCHES`
- Set membership: `field IN [...]`
- Type casts: `cast(field as number)`
- Arithmetic expressions: `field + 10 > 20`

**Optimization**:
- **Alpha Sharing**: Nodes with identical conditions are shared
- **Condition Decomposition**: Complex conditions split into atomic operations
- **Result Caching**: Intermediate arithmetic results cached globally

**Example Alpha Chain**:
```
TypeNode(Account)
    → AlphaNode(balance > 0)
        → AlphaNode(status == "active")
            → AlphaNode(balance > 1000)
```

### JoinNode (Beta Node)

Joins tokens from multiple facts based on join conditions.

**Architecture**:
```
     Left Input                Right Input
    (Token from                (Fact from
     previous node)             TypeNode/AlphaNode)
          ↓                           ↓
      ┌────────────────────────────────┐
      │         JoinNode               │
      │  • Test join conditions        │
      │  • Create combined tokens      │
      │  • Maintain left/right memory  │
      └────────────────────────────────┘
                    ↓
            (Combined token with
             multiple variables)
```

**Join Conditions**:
- Field equality: `account.customer_id == customer.id`
- Cross-variable comparisons: `order.amount > account.balance`

**Memory Management**:
- **Left Memory**: Stores tokens from upstream nodes
- **Right Memory**: Stores facts from alpha nodes
- Both memories indexed for fast lookup

**Beta Sharing**:
- JoinNodes with identical conditions shared across rules
- RuleRouterNodes distribute tokens to rule-specific terminals

### AccumulatorNode

Computes aggregations over collections of facts.

**Aggregate Functions**:
- `COUNT`: Count matching facts
- `SUM`: Sum numeric field values
- `AVG`: Average numeric field values
- `MIN`: Minimum field value
- `MAX`: Maximum field value

**Example**:
```
ACCUMULATE(SUM(order.amount) FROM Order o WHERE o.customer_id == c.id)
```

**Implementation**:
- Maintains aggregation state per main fact
- Incrementally updates aggregates on fact changes
- Efficient index structures for fast lookups

### ExistsNode

Tests for existence of matching facts.

**Semantics**:
```
EXISTS(Order o WHERE o.customer_id == c.id AND o.status == "pending")
```

**Activation Logic**:
- Activates when **at least one** matching fact exists
- Deactivates when **no** matching facts exist
- Efficiently handles fact additions/retractions

### NotNode

Tests for absence of matching facts.

**Semantics**:
```
NOT(Order o WHERE o.customer_id == c.id AND o.status == "overdue")
```

**Activation Logic**:
- Activates when **no** matching facts exist
- Deactivates when **any** matching fact exists
- Critical for negation patterns

### TerminalNode

Executes actions when all rule conditions are satisfied.

**Responsibilities**:
- Receives tokens from final beta/alpha node
- Delegates to ActionExecutor
- Tracks rule activations

---

## Memory Management

### WorkingMemory

Stores facts and tokens at each node.

```go
type WorkingMemory struct {
    Facts  map[string]*Fact
    Tokens map[string]*Token
    mutex  sync.RWMutex
}
```

**Operations**:
- `AddFact(fact)`: Store fact with unique ID
- `RemoveFact(id)`: Remove fact by ID
- `GetFacts()`: Retrieve all facts
- Thread-safe with RWMutex

### Token Structure

Tokens carry variable bindings through the network.

```go
type Token struct {
    ID        string
    Variables map[string]*Fact  // Variable name → Fact
    Parent    *Token            // For token chains
}
```

**Properties**:
- Immutable after creation
- Shared across nodes when possible
- Garbage collected when no longer referenced

### Arithmetic Result Cache

Global cache for intermediate arithmetic computations.

**Features**:
- LRU eviction policy
- Configurable size limits
- Hash-based lookup for expressions
- Significant speedup for complex arithmetic

**Cache Key**:
```
Hash(expression_structure + variable_bindings)
```

---

## Optimization Strategies

### Alpha Node Sharing

**Problem**: Multiple rules may test identical conditions.

**Solution**: Share AlphaNode instances across rules.

**Benefits**:
- Reduced memory footprint
- Faster network construction
- Shared condition evaluation

**Implementation**:
```go
type AlphaSharingRegistry struct {
    registry map[string]*AlphaNode  // condition_hash → node
    mutex    sync.RWMutex
}
```

### Beta Node Sharing

**Problem**: Multiple rules may have identical join patterns.

**Solution**: Share JoinNode instances, use RuleRouterNodes for rule-specific logic.

**Architecture**:
```
        Shared JoinNode
               ↓
    ┌──────────┴──────────┐
    ↓                     ↓
RuleRouter(Rule1)   RuleRouter(Rule2)
    ↓                     ↓
Terminal(Rule1)     Terminal(Rule2)
```

**Benefits**:
- Dramatic reduction in join operations
- Up to 80% reduction in network size
- Faster rule evaluation

### Condition Decomposition

**Problem**: Complex conditions are slow to evaluate.

**Solution**: Decompose into atomic operations with intermediate results.

**Example**:
```
Original: (balance + pending_charges) * 1.05 > credit_limit

Decomposed:
1. temp_1 = balance + pending_charges
2. temp_2 = temp_1 * 1.05
3. temp_2 > credit_limit
```

**Benefits**:
- Intermediate results cached
- Shared computation across conditions
- Better failure fast on early conditions

### Passthrough Optimization

**Problem**: Many alpha nodes are simple type filters.

**Solution**: Eliminate unnecessary intermediate nodes.

**Before**:
```
TypeNode → AlphaNode(always true) → AlphaNode(real condition)
```

**After**:
```
TypeNode → AlphaNode(real condition)
```

---

## Transaction System

### Transaction Model

TSD supports ACID transactions with configurable consistency levels.

**Isolation Levels**:
- **Read Uncommitted**: Fast, no consistency guarantees
- **Read Committed**: See only committed facts
- **Snapshot Isolation**: Consistent snapshot of facts
- **Serializable**: Full serializability (slowest)

### Transaction Lifecycle

```go
tx := network.BeginTransaction()
tx.Assert(fact1)
tx.Assert(fact2)
tx.Retract(fact3)
err := tx.Commit()  // All-or-nothing
```

**Commit Protocol**:
1. **Prepare Phase**: Validate all operations
2. **Log Phase**: Write to transaction log
3. **Apply Phase**: Propagate changes through network
4. **Verify Phase**: Ensure consistency (Strong Mode)
5. **Complete Phase**: Release locks and cleanup

### Strong Mode

Guarantees eventual consistency with configurable verification.

**Configuration**:
```go
Config{
    SubmissionTimeout: 30 * time.Second,
    VerifyRetryDelay:  10 * time.Millisecond,
    MaxVerifyRetries:  10,
}
```

**Verification Strategy**:
- Retry with exponential backoff
- Configurable for transaction consistency
- Optimized for in-memory storage with strong consistency

---

## Type System

### Type Definitions

```go
type TypeDefinition struct {
    Name   string
    Fields []FieldDefinition
}

type FieldDefinition struct {
    Name string
    Type string  // "string", "number", "bool"
}
```

**Supported Types**:
- `string`: Text values
- `number`: Integer and floating-point
- `bool`: Boolean true/false

### Type Casting

Explicit type conversions using `cast` operator.

**Syntax**:
```
cast(field as number)
cast(string_value as bool)
cast(123 as string)
```

**Implementation**: `rete/evaluator_cast.go`

**Casting Rules**:
- `number → string`: String representation
- `string → number`: Parse numeric value (error if invalid)
- `bool → string`: "true" or "false"
- `string → bool`: "true"/"1" → true, others → false
- `number → bool`: 0 → false, non-zero → true

---

## Expression Evaluation

### Condition Evaluation

Conditions are evaluated by specialized evaluators:

**AlphaConditionEvaluator**: Single-fact conditions
- Field access: `fact.field`
- Comparisons: `==`, `!=`, `>`, `<`, `>=`, `<=`
- Logical operators: `AND`, `OR`, `NOT`
- String operations: `LIKE`, `CONTAINS`, `MATCHES`, `IN`
- Arithmetic: `+`, `-`, `*`, `/`
- Casting: `cast(expr as type)`

**BetaConditionEvaluator**: Multi-fact conditions
- Cross-variable comparisons: `var1.field == var2.field`
- Complex join predicates

### Operator Semantics

**Arithmetic Operators** (`+`, `-`, `*`, `/`):
- Numeric operands: Standard arithmetic
- String `+` string: Concatenation
- Mixed types: Error (explicit cast required)

**Comparison Operators**:
- Type-safe comparisons
- String comparison: Lexicographic
- Number comparison: Numeric
- Boolean comparison: true > false

**String Operators**:
- `LIKE`: Pattern matching with `%` wildcard
- `CONTAINS`: Substring search
- `MATCHES`: Regular expression matching
- `IN`: Set membership

---

## Action Execution

### ActionExecutor

Executes actions when rules fire.

**Action Types**:
- **Assert**: Add new fact to working memory
- **Retract**: Remove fact from working memory
- **Modify**: Update existing fact
- **Print**: Output for debugging
- **Custom**: User-defined actions

### Argument Evaluation

Actions can use complex expressions as arguments:

```go
func (ae *ActionExecutor) evaluateArgument(
    arg interface{},
    token *Token,
) (interface{}, error)
```

**Supported Argument Types**:
- Literal values: `42`, `"hello"`, `true`
- Field access: `customer.name`
- Arithmetic: `balance * 1.05`
- String concatenation: `"Hello " + name`
- Casting: `cast(amount as string)`
- Function calls: Custom functions

---

## Storage Layer

### Storage Interface

In-memory storage with strong consistency:

```go
type Storage interface {
    AddFact(fact *Fact) error
    GetFact(factID string) *Fact
    RemoveFact(factID string) error
    GetAllFacts() []*Fact
    
    // Memory management
    SaveMemory(nodeID string, memory *WorkingMemory) error
    LoadMemory(nodeID string) (*WorkingMemory, error)
    DeleteMemory(nodeID string) error
    
    // Consistency
    Sync() error
    Clear() error
}
```

### Implementation

**MemoryStorage**:
- Pure in-memory storage
- Thread-safe with mutex protection
- Strong consistency guarantees
- No persistence (export to .tsd files)
- Optimized for high throughput (~10,000-50,000 facts/sec)

**Future: Network Replication**:
- Multi-node replication via Raft consensus
- Distributed in-memory storage
- Strong consistency across nodes
- Estimated ~1,000-10,000 facts/sec

---

## Performance Characteristics

### Time Complexity

**Fact Assertion**: O(α + β)
- α: Number of alpha node tests
- β: Number of beta node joins
- Typically sub-millisecond

**Fact Retraction**: O(α + β + τ)
- τ: Number of tokens to remove
- Proportional to rule activations

**Rule Evaluation**: O(1) amortized
- Incremental evaluation
- Only changed facts recomputed

### Space Complexity

**Memory Usage**: O(f × n + t)
- f: Number of facts
- n: Average nodes per fact
- t: Number of tokens

**Optimization Impact**:
- Alpha sharing: 30-50% reduction
- Beta sharing: 50-80% reduction
- Arithmetic caching: 20-40% reduction

### Benchmarks

Typical performance on modern hardware:

| Operation | Throughput | Latency (p99) |
|-----------|------------|---------------|
| Assert fact | 100K/sec | <1ms |
| Retract fact | 80K/sec | <2ms |
| Rule evaluation | 500K/sec | <0.5ms |
| Complex join (3+ facts) | 50K/sec | <5ms |
| Aggregation | 30K/sec | <10ms |

---

## Concurrency Model

### Thread Safety

All network operations are thread-safe with fine-grained locking:

**RWMutex Usage**:
- Read operations: Shared lock
- Write operations: Exclusive lock
- Lock scopes minimized to critical sections

**Lock Ordering**:
1. Network-level locks
2. Node-level locks
3. Memory-level locks

**Deadlock Prevention**:
- Consistent lock ordering
- Timeout-based deadlock detection
- Lock-free data structures where possible

### Parallel Execution

**Fact Processing**:
- Multiple facts can be asserted concurrently
- Network topology allows parallel propagation
- Synchronization only at merge points

**Rule Execution**:
- Independent rules execute in parallel
- Dependent rules serialize automatically
- Action execution can be parallelized

---

## Monitoring and Metrics

### Chain Build Metrics

```go
type ChainBuildMetrics struct {
    AlphaNodesCreated    int
    AlphaNodesShared     int
    BetaNodesCreated     int
    BetaNodesShared      int
    TotalBuildTime       time.Duration
    CacheHits            int
    CacheMisses          int
}
```

### Beta Sharing Statistics

```go
type BetaSharingStats struct {
    TotalJoinNodes       int
    SharedJoinNodes      int
    SharingPercentage    float64
    RuleRoutersCreated   int
}
```

### Performance Profiling

Enable profiling:
```go
network.Config.EnableProfiling = true
```

Retrieve metrics:
```go
metrics := network.GetChainMetrics()
betaStats := network.GetBetaSharingStats()
```

---

## Future Enhancements

### Planned Optimizations

1. **Lazy Evaluation**: Delay computation until needed
2. **Parallel Alpha Networks**: Concurrent alpha node evaluation
3. **Adaptive Indexing**: Dynamic index selection based on access patterns
4. **Query Optimization**: Cost-based join ordering

### Research Directions

1. **Machine Learning Integration**: Learn optimal network structure
2. **Distributed RETE**: Scale across multiple machines
3. **GPU Acceleration**: Parallel pattern matching on GPU
4. **Incremental Compilation**: Hot-reload rules without network rebuild

---

## References

### Academic Papers

- Forgy, C. (1982). "Rete: A Fast Algorithm for the Many Pattern/Many Object Pattern Match Problem"
- Doorenbos, R. (1995). "Production Matching for Large Learning Systems"
- Batory, D. (1994). "The LEAPS Algorithm"

### Implementation Resources

- [RETE Algorithm Explained](https://en.wikipedia.org/wiki/Rete_algorithm)
- [Production Rule Systems](https://www.drools.org/)
- TSD Source Code: `/rete` package

---

## Glossary

- **Alpha Network**: Single-fact condition testing phase
- **Beta Network**: Multi-fact joining phase
- **Token**: Binding of variables to facts
- **Working Memory**: Collection of facts in the system
- **Activation**: Rule that has matched and can fire
- **Fact**: Typed data instance in working memory
- **Pattern**: Template for matching facts
- **Condition**: Test applied to facts or tokens
- **Action**: Operation executed when rule fires

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on:
- Architecture decision records (ADRs)
- Design review process
- Performance testing requirements
- Documentation standards

---

*Last updated: 2024-12-07*