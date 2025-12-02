# Demo: Alpha/Beta Separation in Action

This document demonstrates the alpha/beta separation feature with concrete examples showing the network structure and behavior before and after the integration.

---

## Demo 1: Simple Join with Alpha Filter

### Rule Definition
```tsd
type Person(id: string, age: number)
type Order(id: string, personId: string, amount: number)

action notify(personId: string, orderId: string)

rule large_orders : {p: Person, o: Order} / 
    p.id == o.personId AND o.amount > 100 
    ==> notify(p.id, o.id)
```

### Network Structure: BEFORE Integration

```
RootNode
  │
  ├─> TypeNode(Person)
  │     └─> PassthroughAlpha
  │           └─> JoinNode [p.id == o.personId AND o.amount > 100]
  │                 └─> TerminalNode
  │
  └─> TypeNode(Order)
        └─> PassthroughAlpha
              └─> JoinNode↑
```

**Behavior:**
- ALL Person facts → PassthroughAlpha → JoinNode
- ALL Order facts → PassthroughAlpha → JoinNode
- JoinNode evaluates BOTH conditions for every pair

**Example with 1000 Orders (100 with amount > 100):**
- Facts reaching JoinNode: 1000 Orders
- Pairs evaluated: 1000 pairs
- Successful matches: 100 pairs

### Network Structure: AFTER Integration

```
RootNode
  │
  ├─> TypeNode(Person)
  │     └─> PassthroughAlpha[per-rule]
  │           └─> JoinNode [p.id == o.personId]
  │                 └─> TerminalNode
  │
  └─> TypeNode(Order)
        └─> AlphaNode [o.amount > 100]  ← NEW!
              └─> PassthroughAlpha[per-rule]
                    └─> JoinNode↑
```

**Behavior:**
- ALL Person facts → PassthroughAlpha → JoinNode
- ALL Order facts → AlphaNode → **FILTERED** → PassthroughAlpha → JoinNode
- AlphaNode filters to only Orders with amount > 100
- JoinNode evaluates ONLY the join condition (p.id == o.personId)

**Example with 1000 Orders (100 with amount > 100):**
- Facts reaching AlphaNode: 1000 Orders
- Facts passing AlphaNode: 100 Orders (90% filtered!)
- Pairs evaluated in JoinNode: 100 pairs
- Successful matches: 100 pairs

**Result:** 90% reduction in join evaluations! ⚡

---

## Demo 2: Multiple Alpha Filters on Same Variable

### Rule Definition
```tsd
type Product(id: string, price: number, stock: number)
type Order(id: string, productId: string, quantity: number)

action process(orderId: string, productId: string)

rule valid_orders : {p: Product, o: Order} / 
    p.id == o.productId AND 
    p.stock >= o.quantity AND 
    p.price > 0 AND 
    o.quantity > 0
    ==> process(o.id, p.id)
```

### Network Structure: AFTER Integration

```
RootNode
  │
  ├─> TypeNode(Product)
  │     └─> AlphaNode [p.price > 0]  ← Alpha filter 1
  │           └─> PassthroughAlpha
  │                 └─> JoinNode [p.id == o.productId AND p.stock >= o.quantity]
  │                       └─> TerminalNode
  │
  └─> TypeNode(Order)
        └─> AlphaNode [o.quantity > 0]  ← Alpha filter 2
              └─> PassthroughAlpha
                    └─> JoinNode↑
```

**Key Points:**
- `p.price > 0` - Alpha condition (single variable) → AlphaNode
- `o.quantity > 0` - Alpha condition (single variable) → AlphaNode
- `p.id == o.productId` - Beta condition (two variables) → JoinNode
- `p.stock >= o.quantity` - Beta condition (two variables) → JoinNode

**Performance:**
- Products with price ≤ 0: Filtered at TypeNode level
- Orders with quantity ≤ 0: Filtered at TypeNode level
- Only valid facts reach the join
- Significant reduction in join space!

---

## Demo 3: Three-Way Join with Cascade

### Rule Definition
```tsd
type Customer(id: string, tier: string)
type Order(id: string, customerId: string, amount: number)
type Shipment(id: string, orderId: string, status: string)

action notify(customerId: string, orderId: string, shipmentId: string)

rule vip_ready_shipments : {c: Customer, o: Order, s: Shipment} / 
    c.id == o.customerId AND 
    o.id == s.orderId AND 
    c.tier == "VIP" AND 
    s.status == "READY"
    ==> notify(c.id, o.id, s.id)
```

### Network Structure: AFTER Integration (Cascade)

```
RootNode
  │
  ├─> TypeNode(Customer)
  │     └─> AlphaNode [c.tier == "VIP"]
  │           └─> PassthroughAlpha
  │                 └─> JoinNode₁ [c.id == o.customerId]
  │                       │
  ├─> TypeNode(Order)     │
  │     └─> PassthroughAlpha
  │           └─> JoinNode₁↑
  │                 └─> JoinNode₂ [o.id == s.orderId]
  │                       │
  └─> TypeNode(Shipment)  │
        └─> AlphaNode [s.status == "READY"]
              └─> PassthroughAlpha
                    └─> JoinNode₂↑
                          └─> TerminalNode
```

**Cascade Behavior:**
1. **Level 1:** Customer (filtered for VIP) ⋈ Order
2. **Level 2:** (Customer, Order) ⋈ Shipment (filtered for READY)

**Example with 10,000 Customers (1% VIP), 100,000 Shipments (10% READY):**

Without Alpha Filters:
- Customer × Order: 10,000 × all_orders pairs
- Then × Shipment: result × 100,000
- **Massive join space!**

With Alpha Filters:
- AlphaNode filters Customers to 100 VIPs
- Customer × Order: 100 × relevant_orders pairs
- AlphaNode filters Shipments to 10,000 READY
- (Customer, Order) × Shipment: result × 10,000
- **99% reduction in initial join space!** 🚀

---

## Demo 4: Per-Rule Passthrough Isolation

### Rule Definitions
```tsd
type Person(id: string, age: number)
type Order(id: string, personId: string, amount: number)

action notify_large(personId: string)
action notify_very_large(personId: string)

rule large_orders : {p: Person, o: Order} / 
    p.id == o.personId AND o.amount > 100 
    ==> notify_large(p.id)

rule very_large_orders : {p: Person, o: Order} / 
    p.id == o.personId AND o.amount > 500 
    ==> notify_very_large(p.id)
```

### Why Per-Rule Passthroughs Matter

**Without per-rule isolation (BUGGY):**
```
TypeNode(Order)
  └─> AlphaNode [o.amount > 100]  ← Created for first rule
        └─> SHARED PassthroughAlpha
              ├─> JoinNode[rule1]
              └─> JoinNode[rule2]  ← WRONG! Gets filtered facts from rule1's filter!
```

**Problem:** rule2 would receive Orders filtered by rule1's condition (amount > 100)
- rule2 expects Orders with amount > 500
- But gets Orders with amount > 100
- Result: **INCORRECT ACTIVATIONS**

**With per-rule isolation (CORRECT):**
```
TypeNode(Order)
  ├─> AlphaNode [o.amount > 100]
  │     └─> PassthroughAlpha[rule1]
  │           └─> JoinNode[rule1]
  │
  └─> AlphaNode [o.amount > 500]
        └─> PassthroughAlpha[rule2]
              └─> JoinNode[rule2]
```

**Result:** Each rule gets exactly the facts it needs! ✅

**Test Case:**
- Person(id: p1)
- Order(id: o1, personId: p1, amount: 150)
- Order(id: o2, personId: p1, amount: 600)

**Expected Activations:**
- `large_orders`: 2 activations (o1 and o2)
- `very_large_orders`: 1 activation (o2 only)

✅ **With per-rule passthroughs:** CORRECT
❌ **With shared passthroughs:** INCORRECT (both rules would see same filtered set)

---

## Demo 5: Arithmetic in Alpha Filters

### Rule Definition
```tsd
type Product(id: string, price: number, quantity: number)
type Cart(id: string, productId: string, items: number)

action checkout(cartId: string, productId: string)

rule high_value_carts : {p: Product, c: Cart} / 
    p.id == c.productId AND 
    (p.price * c.items) > 1000
    ==> checkout(c.id, p.id)
```

### Network Structure: AFTER Integration

```
RootNode
  │
  ├─> TypeNode(Product)
  │     └─> PassthroughAlpha
  │           └─> JoinNode [p.id == c.productId AND (p.price * c.items) > 1000]
  │                 └─> TerminalNode
  │
  └─> TypeNode(Cart)
        └─> PassthroughAlpha
              └─> JoinNode↑
```

**Note:** `(p.price * c.items) > 1000` is a **beta condition** because it uses variables from both Product (`p.price`) and Cart (`c.items`).

**Cannot be alpha-extracted** because it requires both facts to be present.

**However, if the condition was `p.price * 10 > 1000` (single variable):**
```
TypeNode(Product)
  └─> AlphaNode [p.price * 10 > 1000]  ← Arithmetic alpha filter!
        └─> PassthroughAlpha
              └─> JoinNode [p.id == c.productId]
```

---

## Demo 6: Real-World Performance Example

### Scenario: E-commerce Order Processing

**System:** 
- 1,000,000 products
- 100,000 orders per day
- 50,000 customers (1% are VIP)

**Rule:**
```tsd
rule vip_large_orders : {c: Customer, o: Order, p: Product} / 
    c.id == o.customerId AND 
    o.productId == p.id AND 
    c.tier == "VIP" AND 
    o.amount > 500 AND 
    p.inStock == true
    ==> prioritize_order(o.id)
```

### Without Alpha/Beta Separation

**Network:**
- All Customers → JoinNode₁
- All Orders → JoinNode₁
- JoinNode₁ output → JoinNode₂
- All Products → JoinNode₂

**Join Space:**
- JoinNode₁: 50,000 × 100,000 = 5 billion pairs
- Then filter for VIP and amount > 500
- JoinNode₂: (reduced) × 1,000,000 products

**Result:** Extremely slow, potentially unusable

### With Alpha/Beta Separation

**Network:**
- Customers → AlphaNode[tier=="VIP"] → 500 VIP customers
- Orders → AlphaNode[amount>500] → 10,000 qualifying orders
- Products → AlphaNode[inStock==true] → 800,000 in-stock products

**Join Space:**
- JoinNode₁: 500 × 10,000 = 5 million pairs (99.9% reduction!)
- JoinNode₂: (reduced) × 800,000 products (but much smaller input)

**Performance Improvement:**
- **1000x fewer join evaluations at first level**
- **Queries that took minutes now take seconds**
- **System can handle real-time processing** ✅

---

## Verification: How to Check Alpha/Beta Separation

### 1. Network Stats
```go
stats := network.GetNetworkStats()
alphaNodes := stats["alpha_nodes"].(int)
betaNodes := stats["beta_nodes"].(int)

// With proper separation, alphaNodes > 0 for filtered rules
```

### 2. Visual Inspection
Walk the network from TypeNodes:
```
TypeNode → AlphaNode? → Passthrough → JoinNode
            ↑
         Should exist for filtered variables
```

### 3. Test Activations
Submit test facts and verify:
- Facts filtered at alpha level don't reach JoinNodes
- Only correctly filtered facts produce activations
- No duplicate or incorrect activations

### 4. Performance Monitoring
```go
// Count join evaluations
joinEvaluations := 0
for _, betaNode := range network.BetaNodes {
    joinEvaluations += betaNode.GetEvaluationCount()
}

// Should be much lower with alpha filtering
```

---

## Summary

### Key Takeaways

1. **Alpha conditions** (single variable) → **AlphaNodes** → Early filtering
2. **Beta conditions** (multiple variables) → **JoinNodes** → Join evaluation
3. **Per-rule passthroughs** → Correct isolation between rules
4. **Performance gains** → Up to 99% reduction in join evaluations
5. **Correct semantics** → Follows classical RETE architecture

### When Alpha/Beta Separation Helps Most

✅ **Large fact bases** - More facts = bigger win from early filtering  
✅ **High selectivity filters** - Filters that eliminate 90%+ of facts  
✅ **Multi-way joins** - Cascade of joins with intermediate filtering  
✅ **Real-time systems** - Need fast response times  

### Rules of Thumb

- If a condition uses **one variable** → Alpha (filtered early)
- If a condition uses **multiple variables** → Beta (join condition)
- More alpha filters = Less work for JoinNodes = Faster execution

---

**Status:** ✅ Fully implemented and tested  
**Documentation:** See [Implementation Details](./IMPLEMENTATION_ALPHA_BETA_INTEGRATION.md)  
**Date:** 2025-12-02