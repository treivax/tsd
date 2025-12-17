<!--
Copyright (c) 2025 TSD Contributors
Licensed under the MIT License
See LICENSE file in the project root for full license text
-->

# BetaChains Examples & Visual Guide

**Version:** 1.0  
**Date:** 2025-01-XX  
**Related:** `BETA_CHAINS_DESIGN.md`, `BETA_SHARING_EXAMPLES.md`

---

## Table of Contents

1. [Visual Diagrams](#visual-diagrams)
2. [Basic Examples](#basic-examples)
3. [Optimization Examples](#optimization-examples)
4. [Prefix Sharing Examples](#prefix-sharing-examples)
5. [Real-World Use Cases](#real-world-use-cases)
6. [Performance Comparisons](#performance-comparisons)
7. [Anti-Patterns](#anti-patterns)

---

## Visual Diagrams

### Basic Chain Structure

#### Binary Join (2 Variables)

```
┌────────────────┐
│ TypeNode       │
│ (Customer)     │
└───────┬────────┘
        │ Facts: Customer objects
        │
        ▼
    ┌───────────────────────────┐
    │  PassThroughAlphaNode     │
    │  (c: Customer)            │
    └───────┬───────────────────┘
            │ Left Input
            ▼
        ┌─────────────────────────────┐
        │                             │◄──── Right Input
        │  JoinNode                   │
        │  c ⋈ o                      │      ┌────────────────┐
        │                             │      │ TypeNode       │
        │  Condition:                 │      │ (Order)        │
        │  c.id == o.customer_id      │      └───────┬────────┘
        │                             │              │
        │  LeftMemory: [c₁, c₂, ...]  │              │
        │  RightMemory: [o₁, o₂, ...] │              ▼
        │  ResultMemory: [(c₁,o₁)...] │      ┌───────────────┐
        └──────────┬──────────────────┘      │ PassThrough   │
                   │                          │ (o: Order)    │
                   │ Tokens: (c, o) pairs     └───────┬───────┘
                   ▼                                  │
           ┌────────────────┐                         │
           │ TerminalNode   │◄────────────────────────┘
           │ (action)       │
           └────────────────┘
```

#### Cascade Join (3 Variables)

```
TypeNode(Person) ──────────┐
                           │ Left Input
                           ▼
                    ┌──────────────┐
                    │  JoinNode₁   │◄──── Right Input
                    │  p ⋈ a       │              │
                    │              │      TypeNode(Address)
                    └──────┬───────┘
                           │ Tokens: (p, a)
                           │ Left Input for next join
                           ▼
                    ┌──────────────┐
                    │  JoinNode₂   │◄──── Right Input
                    │  (p,a) ⋈ ph  │              │
                    │              │      TypeNode(Phone)
                    └──────┬───────┘
                           │ Tokens: (p, a, ph)
                           ▼
                    ┌──────────────┐
                    │ TerminalNode │
                    └──────────────┘
```

### Optimized vs Unoptimized

#### Scenario: 4-Variable Join

**Rule:**
```
(User u), (Session s), (Page p), (Event e)
WHERE u.id == s.user_id 
  AND s.id == e.session_id
  AND p.id == e.page_id
```

**Unoptimized (Declaration Order):**

```
           ┌──► JoinNode₁ ◄──┐
           │    u ⋈ s        │
    TypeNode(User)     TypeNode(Session)
           │                 
           └─► Result: 10,000 tokens
                     │
                     ▼
           ┌──► JoinNode₂ ◄──┐
           │  (u,s) ⋈ p      │
         (u,s)         TypeNode(Page)
           │
           └─► Result: 8,000 tokens
                     │
                     ▼
           ┌──► JoinNode₃ ◄──┐
           │ (u,s,p) ⋈ e     │
        (u,s,p)        TypeNode(Event)
           │
           └─► Result: 500 tokens

Total Intermediate Tokens: 18,500
```

**Optimized (Selectivity Order):**

```
Selectivity Analysis:
  s ⋈ e: 0.05 (session_id is highly selective)
  u ⋈ s: 0.15 (users have multiple sessions)
  p ⋈ e: 0.20 (pages have many events)

Optimal Order: s ⋈ e → (s,e) ⋈ u → (s,e,u) ⋈ p
```

```
           ┌──► JoinNode₁ ◄──┐
           │    s ⋈ e        │
    TypeNode(Session)  TypeNode(Event)
           │                 
           └─► Result: 500 tokens  (highly selective!)
                     │
                     ▼
           ┌──► JoinNode₂ ◄──┐
           │  (s,e) ⋈ u      │
         (s,e)         TypeNode(User)
           │
           └─► Result: 450 tokens
                     │
                     ▼
           ┌──► JoinNode₃ ◄──┐
           │ (s,e,u) ⋈ p     │
        (s,e,u)        TypeNode(Page)
           │
           └─► Result: 400 tokens

Total Intermediate Tokens: 1,350

IMPROVEMENT: 93% reduction! (18,500 → 1,350)
```

### Prefix Sharing Visualization

#### Two Rules with Common Prefix

**Rule 1:** Customer Order Analysis
**Rule 2:** Customer Return Analysis

```
                    ┌─────────────────┐
                    │ TypeNode        │
                    │ (Customer)      │
                    └────────┬────────┘
                             │
                             ▼
                    ┌────────────────┐
                    │ PassThrough    │
                    │ (c: Customer)  │
                    └────────┬───────┘
                             │ Left
                             ▼
      ┌──────────────────────────────────┐
      │  JoinNode_SHARED                 │◄── Right
      │  c ⋈ o                           │         │
      │  hash: join_a7f3b2c8             │    ┌────┴──────┐
      │  Condition: c.id == o.customer_id│    │ TypeNode  │
      │  RefCount: 2                     │    │ (Order)   │
      │  UsedBy: [Rule1, Rule2]          │    └───────────┘
      └──────┬───────────────┬───────────┘
             │               │
             │ Rule 1        │ Rule 2
             ▼               ▼
    ┌────────────────┐  ┌────────────────┐
    │  JoinNode₂     │  │  JoinNode₃     │
    │  (c,o) ⋈ p     │  │  (c,o) ⋈ r     │
    │  (Product)     │  │  (Return)      │
    └────────┬───────┘  └────────┬───────┘
             │                   │
             ▼                   ▼
      ┌───────────┐       ┌───────────┐
      │ Terminal₁ │       │ Terminal₂ │
      └───────────┘       └───────────┘

PREFIX CACHE:
  Signature: "hash(c⋈o|Customer+Order|c.id==o.customer_id)"
  Nodes: [JoinNode_SHARED]
  RefCount: 2
  Savings: 1 JoinNode, 1 set of connections
```

---

## Basic Examples

### Example 1: Simple Binary Join

**Scenario:** Find high-value customers

**Rule Definition:**
```go
Rule: "HighValueCustomers"
Variables:
  c: Customer
  o: Order
Condition:
  c.id == o.customer_id AND o.total > 10000
Action:
  notify_sales_team(c, o)
```

**BetaChain Structure:**
```go
BetaChain{
  Nodes: [
    JoinNode{
      ID: "HighValueCustomers_join",
      LeftVars: ["c"],
      RightVars: ["o"],
      Condition: {
        type: "and",
        conditions: [
          {operator: "==", left: "c.id", right: "o.customer_id"},
          {operator: ">", left: "o.total", right: 10000}
        ]
      }
    }
  ],
  Hashes: ["join_b4c7e9f1"],
  JoinSpecs: [
    JoinSpec{
      LeftVars: ["c"],
      RightVars: ["o"],
      VarTypes: {"c": "Customer", "o": "Order"},
      Selectivity: 0.05,  // Very selective (5% of orders)
      JoinMethod: "hash"
    }
  ],
  FinalNode: JoinNode_ref,
  RuleID: "HighValueCustomers",
  BuildStrategy: "binary"
}
```

**Network Diagram:**
```
Customer Facts → PassThrough(c) ──┐
                                  ├→ JoinNode → Terminal
Order Facts → PassThrough(o) ─────┘

Memories:
  LeftMemory: [Customer{id:1, name:"Alice"}, ...]
  RightMemory: [Order{id:101, customer_id:1, total:15000}, ...]
  ResultMemory: [(Customer{1}, Order{101}), ...]
```

### Example 2: Three-Variable Cascade

**Scenario:** Complete user profile check

**Rule Definition:**
```go
Rule: "CompleteProfile"
Variables:
  u: User
  e: Email
  p: Phone
Condition:
  u.id == e.user_id AND u.id == p.user_id
Action:
  mark_profile_complete(u)
```

**BetaChain Structure:**
```go
BetaChain{
  Nodes: [
    JoinNode{
      ID: "CompleteProfile_join_0_1",
      LeftVars: ["u"],
      RightVars: ["e"],
      Condition: {operator: "==", left: "u.id", right: "e.user_id"}
    },
    JoinNode{
      ID: "CompleteProfile_join_2",
      LeftVars: ["u", "e"],
      RightVars: ["p"],
      Condition: {operator: "==", left: "u.id", right: "p.user_id"}
    }
  ],
  Hashes: ["join_c8d9a3e2", "join_f1b5d8a7"],
  JoinSpecs: [
    JoinSpec{
      LeftVars: ["u"],
      RightVars: ["e"],
      Selectivity: 0.15,
      JoinMethod: "hash"
    },
    JoinSpec{
      LeftVars: ["u", "e"],
      RightVars: ["p"],
      Selectivity: 0.20,
      JoinMethod: "nested_loop"
    }
  ],
  VariablesAtStage: [
    ["u", "e"],      // After first join
    ["u", "e", "p"]  // After second join
  ],
  FinalNode: JoinNode2_ref,
  RuleID: "CompleteProfile",
  BuildStrategy: "cascade"
}
```

**Step-by-Step Execution:**

```
Step 1: User fact arrives
  User{id:1, name:"Alice"} → PassThrough(u) → LeftMemory of JoinNode₁

Step 2: Email fact arrives
  Email{id:101, user_id:1, address:"alice@example.com"} 
    → PassThrough(e) → RightMemory of JoinNode₁
    → Join evaluation: u.id(1) == e.user_id(1) ✓
    → Token created: (User{1}, Email{101})
    → Propagate to LeftMemory of JoinNode₂

Step 3: Phone fact arrives
  Phone{id:201, user_id:1, number:"555-1234"}
    → PassThrough(p) → RightMemory of JoinNode₂
    → Join evaluation: u.id(1) == p.user_id(1) ✓
    → Token created: (User{1}, Email{101}, Phone{201})
    → Propagate to Terminal
    → Action executed: mark_profile_complete(User{1})
```

---

## Optimization Examples

### Example 3: Suboptimal Declaration Order

**Scenario:** User activity tracking with poor variable ordering

**Rule:**
```go
Rule: "SuspiciousLogin"
Variables:
  u: User        // 100,000 users
  l: Login       // 10,000,000 logins
  d: Device      // 50,000 devices
  loc: Location  // 1,000 locations
Condition:
  u.id == l.user_id
  AND l.device_id == d.id
  AND l.location_id == loc.id
  AND d.is_new == true
  AND loc.country != u.home_country
```

**Unoptimized Build (Declaration Order):**

```
Join Order: u → l → d → loc

Step 1: u ⋈ l
  Selectivity: 0.20 (users have many logins)
  Left: 100,000 users
  Right: 10,000,000 logins
  Result: 2,000,000 intermediate tokens
  
Step 2: (u,l) ⋈ d
  Selectivity: 0.30 (many logins per device)
  Left: 2,000,000 (u,l) pairs
  Right: 50,000 devices
  Result: 600,000 intermediate tokens
  
Step 3: (u,l,d) ⋈ loc
  Selectivity: 0.25 (locations filtered by country mismatch)
  Left: 600,000 (u,l,d) triples
  Right: 1,000 locations
  Result: 150,000 final tokens

Total Intermediate Processing: 2,750,000 tokens
Peak Memory: ~220 MB (assuming 80 bytes/token)
Build Time: 850ms
```

**Optimized Build (Selectivity Order):**

```
Selectivity Analysis:
  l ⋈ d (with d.is_new filter): 0.03  ← Most selective!
  (l,d) ⋈ u: 0.10
  (l,d,u) ⋈ loc (with country filter): 0.15

Optimal Join Order: l → d → u → loc
```

```
Step 1: l ⋈ d (with d.is_new == true)
  Selectivity: 0.03 (only new devices)
  Left: 10,000,000 logins
  Right: 50,000 devices (1,500 are new)
  Result: 300,000 intermediate tokens
  
Step 2: (l,d) ⋈ u
  Selectivity: 0.10
  Left: 300,000 (l,d) pairs
  Right: 100,000 users
  Result: 30,000 intermediate tokens
  
Step 3: (l,d,u) ⋈ loc (with country filter)
  Selectivity: 0.15
  Left: 30,000 (l,d,u) triples
  Right: 1,000 locations
  Result: 4,500 final tokens

Total Intermediate Processing: 334,500 tokens
Peak Memory: ~27 MB
Build Time: 180ms

IMPROVEMENTS:
  - 88% reduction in intermediate tokens (2.75M → 335K)
  - 88% reduction in memory (220MB → 27MB)
  - 79% faster compilation (850ms → 180ms)
```

**Visual Comparison:**

```
Unoptimized (Declaration Order):
████████████████████████████████████████████████ 2.75M tokens

Optimized (Selectivity Order):
██████ 335K tokens

Memory Usage:
Unoptimized: ████████████████████ 220MB
Optimized:   ███ 27MB
```

### Example 4: Complex Multi-Way Join

**Scenario:** E-commerce order fulfillment tracking

**Rule:**
```go
Rule: "OrderReadyToShip"
Variables:
  o: Order
  p: Payment
  i: Inventory
  s: Shipping
  c: Customer
Condition:
  o.id == p.order_id AND p.status == "completed"
  AND o.product_id == i.product_id AND i.quantity > 0
  AND o.id == s.order_id AND s.status == "pending"
  AND o.customer_id == c.id AND c.address_verified == true
```

**Selectivity Matrix:**

```
         o      p      i      s      c
    ┌──────┬──────┬──────┬──────┬──────┐
  o │  -   │ 0.15 │ 0.30 │ 0.20 │ 0.25 │
  p │ 0.15 │  -   │  -   │  -   │  -   │
  i │ 0.30 │  -   │  -   │  -   │  -   │
  s │ 0.20 │  -   │  -   │  -   │  -   │
  c │ 0.25 │  -   │  -   │  -   │  -   │
    └──────┴──────┴──────┴──────┴──────┘

Most selective joins:
  1. o ⋈ p: 0.15 (completed payments)
  2. (o,p) ⋈ s: 0.20 (pending shipping)
  3. (o,p,s) ⋈ c: 0.25 (verified addresses)
  4. (o,p,s,c) ⋈ i: 0.30 (in-stock items)
```

**Optimized BetaChain:**

```
                   o ────────┐
                             ├─► JoinNode₁ (o ⋈ p) ────┐
                   p ────────┘  selectivity: 0.15      │
                                result: 1,500 tokens   │
                                                        ├─► JoinNode₂ ((o,p) ⋈ s) ──┐
                   s ──────────────────────────────────┘   selectivity: 0.20       │
                                                            result: 300 tokens       │
                                                                                     ├─► JoinNode₃
                   c ───────────────────────────────────────────────────────────────┘  ((o,p,s) ⋈ c)
                                                                                        selectivity: 0.25
                                                                                        result: 75 tokens
                                                                                            │
                   i ───────────────────────────────────────────────────────────────────────┘
                                                            JoinNode₄ ((o,p,s,c) ⋈ i)
                                                            selectivity: 0.30
                                                            result: 22 tokens → Terminal
```

---

## Prefix Sharing Examples

### Example 5: Customer Analysis Rules

**Scenario:** Multiple rules analyze customer-order relationships

**Rule 1: High-Value Customers**
```go
Variables: (Customer c), (Order o), (Product p)
Condition: c.id == o.customer_id AND o.product_id == p.id AND p.price > 1000
```

**Rule 2: Frequent Customers**
```go
Variables: (Customer c), (Order o), (Shipping s)
Condition: c.id == o.customer_id AND o.id == s.order_id AND s.delivery_time < 24
```

**Rule 3: Customer Feedback**
```go
Variables: (Customer c), (Order o), (Review r)
Condition: c.id == o.customer_id AND o.id == r.order_id AND r.rating > 4
```

**Shared Prefix Detection:**

```
All three rules share: (Customer c) ⋈ (Order o) with c.id == o.customer_id

Prefix Signature: "join_c4d8f2a1|Customer+Order|c.id==o.customer_id"
```

**Network Structure:**

```
TypeNode(Customer) ──┐
                     ├─► JoinNode_SHARED (c ⋈ o) ───┬─► JoinNode₂ (Rule 1) → Terminal₁
TypeNode(Order) ─────┘      hash: join_c4d8f2a1    │    (c,o) ⋈ p
                            RefCount: 3              │
                                                     ├─► JoinNode₃ (Rule 2) → Terminal₂
                                                     │    (c,o) ⋈ s
                                                     │
                                                     └─► JoinNode₄ (Rule 3) → Terminal₃
                                                          (c,o) ⋈ r
```

**Prefix Cache Entry:**

```go
prefixCache["join_c4d8f2a1|Customer+Order|c.id==o.customer_id"] = &BetaChainPrefix{
  Nodes: [JoinNode_SHARED],
  Variables: ["c", "o"],
  Signature: "join_c4d8f2a1|Customer+Order|c.id==o.customer_id",
  RefCount: 3,
  UsedByRules: ["HighValueCustomers", "FrequentCustomers", "CustomerFeedback"],
}
```

**Savings:**
- **Nodes saved:** 2 JoinNodes (3 rules share 1 node instead of 3 separate nodes)
- **Memory saved:** 2 × (LeftMemory + RightMemory + ResultMemory) = ~40-60% reduction
- **Computation saved:** Join evaluation happens once, results shared

### Example 6: Progressive Prefix Sharing

**Scenario:** Increasingly specific rules

**Rule 1:** `(Person p) ⋈ (Address a)`
**Rule 2:** `(Person p) ⋈ (Address a) ⋈ (Phone ph)`
**Rule 3:** `(Person p) ⋈ (Address a) ⋈ (Phone ph) ⋈ (Email e)`

**Network with Progressive Sharing:**

```
Level 0: TypeNodes
┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐
│ Person   │  │ Address  │  │ Phone    │  │ Email    │
└────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘
     │             │              │              │
Level 1: First Join (Shared by ALL rules)
     └──────┬──────┘              │              │
            ▼                     │              │
    ┌────────────────┐            │              │
    │ JoinNode₁      │────────────┴──────────────┼──► Terminal₁ (Rule 1)
    │ p ⋈ a          │                           │
    │ RefCount: 3    │                           │
    └───────┬────────┘                           │
            │                                    │
Level 2: Second Join (Shared by Rule 2 & 3)     │
            └──────┬─────────────────────────────┘
                   ▼
           ┌────────────────┐
           │ JoinNode₂      │──────────────────────► Terminal₂ (Rule 2)
           │ (p,a) ⋈ ph     │
           │ RefCount: 2    │
           └───────┬────────┘
                   │
Level 3: Third Join (Only Rule 3)               
                   └──────┬──────────────────────┘
                          ▼
                  ┌────────────────┐
                  │ JoinNode₃      │──────────────► Terminal₃ (Rule 3)
                  │ (p,a,ph) ⋈ e   │
                  │ RefCount: 1    │
                  └────────────────┘
```

**Prefix Cache:**

```go
// Prefix 1: p ⋈ a
prefixCache["prefix_1"] = {
  Nodes: [JoinNode₁],
  Variables: ["p", "a"],
  RefCount: 3,
  UsedByRules: ["Rule1", "Rule2", "Rule3"]
}

// Prefix 2: p ⋈ a ⋈ ph
prefixCache["prefix_2"] = {
  Nodes: [JoinNode₁, JoinNode₂],
  Variables: ["p", "a", "ph"],
  RefCount: 2,
  UsedByRules: ["Rule2", "Rule3"]
}

// No prefix for Rule 3 alone (RefCount would be 1)
```

**Compilation Timeline:**

```
Time = 0ms: Compile Rule 1
  - Create JoinNode₁ (p ⋈ a)
  - Cache prefix_1
  - Total: 1 node created

Time = 50ms: Compile Rule 2
  - Find prefix_1 in cache ✓
  - Reuse JoinNode₁
  - Create JoinNode₂ (p,a) ⋈ ph)
  - Cache prefix_2
  - Total: 1 node created (1 reused)

Time = 95ms: Compile Rule 3
  - Find prefix_2 in cache ✓
  - Reuse JoinNode₁ and JoinNode₂
  - Create JoinNode₃ ((p,a,ph) ⋈ e)
  - Total: 1 node created (2 reused)

Summary: 3 nodes created, 3 nodes reused, 50% sharing ratio
```

---

## Real-World Use Cases

### Use Case 1: Fraud Detection System

**Context:** Detect potentially fraudulent transactions by correlating user, transaction, device, and location data.

**Rules:**

```go
// Rule 1: Suspicious large transaction
Rule "LargeTransactionNewDevice":
  (User u), (Transaction t), (Device d)
  WHERE u.id == t.user_id
    AND t.device_id == d.id
    AND t.amount > u.average_transaction * 5
    AND d.first_seen > NOW() - 7 days

// Rule 2: Unusual location
Rule "UnusualLocation":
  (User u), (Transaction t), (Location loc)
  WHERE u.id == t.user_id
    AND t.location_id == loc.id
    AND Distance(loc, u.home_location) > 1000 miles
    AND t.amount > 500

// Rule 3: Rapid succession transactions
Rule "RapidTransactions":
  (User u), (Transaction t1), (Transaction t2)
  WHERE u.id == t1.user_id
    AND u.id == t2.user_id
    AND t2.timestamp - t1.timestamp < 5 minutes
    AND t1.location_id != t2.location_id
```

**BetaChain Optimization Impact:**

```
Before Optimization:
  - 9 JoinNodes total (3 per rule)
  - No sharing
  - Average rule execution: 450ms
  - Memory: 180MB

After Optimization:
  - 6 JoinNodes total (3 shared u ⋈ t joins)
  - 33% node reduction
  - Average rule execution: 320ms (29% faster)
  - Memory: 125MB (31% reduction)
  - Prefix cache hit rate: 67%
```

### Use Case 2: Supply Chain Management

**Context:** Track order fulfillment across multiple entities.

**Rules:**

```go
// Rule 1: Order ready for pickup
Rule "ReadyForPickup":
  (Order o), (Inventory i), (Warehouse w), (Carrier c)
  WHERE o.product_id == i.product_id
    AND i.warehouse_id == w.id
    AND o.carrier_id == c.id
    AND i.quantity >= o.quantity
    AND c.availability == true

// Rule 2: Delayed shipment
Rule "DelayedShipment":
  (Order o), (Inventory i), (Warehouse w), (Weather wt)
  WHERE o.product_id == i.product_id
    AND i.warehouse_id == w.id
    AND w.location_id == wt.location_id
    AND wt.severity > 3
    AND o.priority == "urgent"

// Rule 3: Inventory reorder
Rule "ReorderNeeded":
  (Order o), (Inventory i), (Warehouse w), (Supplier s)
  WHERE o.product_id == i.product_id
    AND i.warehouse_id == w.id
    AND i.supplier_id == s.id
    AND i.quantity < i.reorder_point
```

**Shared Prefix:**

All three rules share: `(Order o) ⋈ (Inventory i) ⋈ (Warehouse w)`

**Optimized Network:**

```
TypeNode(Order) ─────┐
                     ├─► JoinNode₁ (o ⋈ i) ──┐
TypeNode(Inventory) ─┘    Selectivity: 0.12   │
                                               ├─► JoinNode₂ ((o,i) ⋈ w) ─┬─► Rule 1 joins
                     TypeNode(Warehouse) ──────┘   Selectivity: 0.08      │   
                                                    SHARED PREFIX          ├─► Rule 2 joins
                                                    RefCount: 3            │
                                                                          └─► Rule 3 joins
```

**Performance Metrics:**

```
Compilation Time:
  Rule 1: 85ms (builds prefix)
  Rule 2: 42ms (reuses prefix, 51% faster)
  Rule 3: 38ms (reuses prefix, 55% faster)

Runtime Performance:
  Without sharing: 3 separate (o⋈i⋈w) computations = 270ms total
  With sharing: 1 (o⋈i⋈w) computation = 90ms total
  Improvement: 67% faster

Memory Usage:
  Without sharing: 3 × (LeftMem + RightMem + ResultMem) = 135MB
  With sharing: 1 × memories + 3 × downstream = 87MB
  Improvement: 36% reduction
```

### Use Case 3: Healthcare Patient Monitoring

**Context:** Monitor patient vital signs, medications, and alerts.

**Rules:**

```go
// Rule 1: Critical condition
Rule "CriticalCondition":
  (Patient p), (VitalSign v), (Doctor d)
  WHERE p.id == v.patient_id
    AND p.doctor_id == d.id
    AND (v.heart_rate > 120 OR v.blood_pressure > 180)
    AND d.on_call == true

// Rule 2: Medication interaction
Rule "MedicationAlert":
  (Patient p), (Prescription prx1), (Prescription prx2)
  WHERE p.id == prx1.patient_id
    AND p.id == prx2.patient_id
    AND prx1.id != prx2.id
    AND HasInteraction(prx1.medication, prx2.medication)

// Rule 3: Allergy warning
Rule "AllergyWarning":
  (Patient p), (Prescription prx), (Allergy a)
  WHERE p.id == prx.patient_id
    AND p.id == a.patient_id
    AND prx.medication == a.substance
```

**Selectivity-Optimized Chains:**

```
Rule 1 (CriticalCondition):
  Join Order: v → p → d (start with critical vitals)
  Selectivity: 0.02 → 0.15 → 0.30
  Rationale: Only 2% of vitals are critical, filter first

Rule 2 (MedicationAlert):
  Join Order: prx1 ⋈ prx2 → p (self-join on prescriptions)
  Selectivity: 0.05 → 0.20
  Rationale: Few prescription pairs interact, check first

Rule 3 (AllergyWarning):
  Join Order: a → p → prx (start with allergy matches)
  Selectivity: 0.03 → 0.10 → 0.25
  Rationale: Exact substance match is highly selective
```

---

## Performance Comparisons

### Benchmark Results

#### Scenario A: 10 Rules, Low Sharing (20%)

```
Metric                    | Without BetaChains | With BetaChains | Improvement
--------------------------|-------------------|-----------------|------------
Total JoinNodes           | 45                | 38              | 16%
Compilation Time          | 420ms             | 380ms           | 10%
Memory Usage              | 92MB              | 81MB            | 12%
Average Rule Exec Time    | 185μs             | 172μs           | 7%
```

#### Scenario B: 50 Rules, Medium Sharing (50%)

```
Metric                    | Without BetaChains | With BetaChains | Improvement
--------------------------|-------------------|-----------------|------------
Total JoinNodes           | 215               | 128             | 40%
Compilation Time          | 2.8s              | 1.6s            | 43%
Memory Usage              | 580MB             | 365MB           | 37%
Average Rule Exec Time    | 425μs             | 280μs           | 34%
Prefix Cache Hit Rate     | N/A               | 58%             | N/A
```

#### Scenario C: 100 Rules, High Sharing (70%)

```
Metric                    | Without BetaChains | With BetaChains | Improvement
--------------------------|-------------------|-----------------|------------
Total JoinNodes           | 520               | 210             | 60%
Compilation Time          | 7.2s              | 2.9s            | 60%
Memory Usage              | 1.4GB             | 620MB           | 56%
Average Rule Exec Time    | 850μs             | 420μs           | 51%
Prefix Cache Hit Rate     | N/A               | 73%             | N/A
Avg Prefix Length         | N/A               | 2.3 nodes       | N/A
```

### Memory Breakdown

**Before BetaChains (50 rules):**
```
JoinNode structures:      215 nodes × 1.2KB = 258MB
Left memories:            215 × 500KB = 107.5MB
Right memories:           215 × 500KB = 107.5MB
Result memories:          215 × 500KB = 107.5MB
Total:                    580MB
```

**After BetaChains (50 rules):**
```
JoinNode structures:      128 nodes × 1.2KB = 154MB  (40% reduction)
Shared left memories:     128 × 500KB = 64MB
Shared right memories:    128 × 500KB = 64MB
Shared result memories:   128 × 500KB = 64MB
Prefix cache overhead:    ~3MB
BetaChainBuilder cache:   ~16MB
Total:                    365MB (37% reduction)
```

---

## Anti-Patterns

### Anti-Pattern 1: Over-Optimization

**Problem:** Spending more time optimizing than saved in execution

**Example:**
```go
// Simple 2-variable rule - no optimization needed
Rule "SimpleJoin":
  (A a), (B b)
  WHERE a.id == b.a_id

// Don't do this:
strategy := OPTIMIZED  // Wastes time on selectivity estimation
// Overhead: +5ms for optimization, saves: 0ms (already optimal)

// Do this:
strategy := BINARY  // Direct binary join
// Fast compilation, same runtime performance
```

**When to avoid optimization:**
- 2-variable joins (already optimal)
- Rules that execute rarely (compile time > execution time savings)
- Prototype/development phase (use CASCADE for predictability)

### Anti-Pattern 2: Ignoring Selectivity Hints

**Problem:** Not providing selectivity hints for complex conditions

**Example:**
```go
// Without hints (poor estimation):
Rule "ComplexFilter":
  (User u), (Event e)
  WHERE u.id == e.user_id
    AND ComplexFunction(u, e) == true  // Unknown selectivity!

// Builder estimates: 0.5 (default)
// Actual selectivity: 0.01 (very selective)
// Result: Suboptimal join order

// With hints (better):
Rule "ComplexFilter":
  (User u), (Event e)
  WHERE u.id == e.user_id
    AND ComplexFunction(u, e) == true
  SELECTIVITY_HINT: 0.01  // Inform optimizer

// Result: Correct join order
```

### Anti-Pattern 3: Excessive Prefix Caching

**Problem:** Caching prefixes that are rarely reused

**Example:**
```go
// 100 rules, all unique join patterns
// Prefix cache size: 100 entries
// Prefix cache hit rate: 0%
// Memory waste: ~50MB for unused cache

// Solution: Set cache eviction policy
BetaChainBuilder{
  maxPrefixCacheSize: 50,
  evictionPolicy: LRU,
  minRefCountForCache: 2,  // Only cache if reused
}
```

### Anti-Pattern 4: Premature Prefix Sharing

**Problem:** Forcing prefix sharing when join semantics differ

**Example:**
```go
// These look similar but have different semantics!

Rule1: (A a), (B b) WHERE a.id == b.a_id AND a.status == "active"
                                              ^^^^^^^^^^^^^^^^
                                              Alpha condition!

Rule2: (A a), (B b) WHERE a.id == b.a_id AND b.status == "active"
                                              ^^^^^^^^^^^^^^^^
                                              Different condition!

// Don't share prefix: join conditions are semantically different
// Even though base join (a.id == b.a_id) is the same,
// the filtering happens at different stages
```

### Anti-Pattern 5: Circular Dependencies

**Problem:** Join conditions that create cycles

**Example:**
```go
Rule "Circular":
  (A a), (B b), (C c)
  WHERE a.id == b.a_id
    AND b.id == c.b_id
    AND c.id == a.c_id  // Creates cycle!

// Optimizer may fail or create suboptimal plan
// Solution: Detect cycles and use cascade strategy
```

---

## Summary

### Key Takeaways

1. **BetaChains optimize join sequences** through selectivity-based ordering and prefix sharing
2. **Binary joins (2 vars)** need no optimization; use `BINARY` strategy
3. **Cascade joins (3+ vars)** benefit from optimization; use `OPTIMIZED` strategy when conditions are complex
4. **Prefix sharing** provides 30-60% memory savings when rules share join patterns
5. **Selectivity estimation** is crucial for optimal join ordering; provide hints for complex conditions

### When to Use Each Strategy

| Strategy | Use When | Benefits | Overhead |
|----------|----------|----------|----------|
| BINARY | 2 variables | Simple, fast | Minimal |
| CASCADE | 3+ vars, simple conditions | Predictable, fast compile | Low |
| OPTIMIZED | 3+ vars, complex conditions | Best runtime performance | Medium compile time |

### Expected Improvements

- **Memory:** 30-60% reduction with high prefix sharing
- **Compilation:** 20-50% faster for rules with shared prefixes  
- **Runtime:** 25-57% faster for optimized multi-way joins
- **Overall:** 40-60% fewer JoinNodes in typical workloads

---

**End of Examples Document**