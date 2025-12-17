# Arithmetic Expressions in TSD Rules

This document provides comprehensive documentation for using arithmetic expressions in TSD (Type System Description) rules and actions.

## Table of Contents

1. [Overview](#overview)
2. [TSD File Syntax](#tsd-file-syntax)
3. [Arithmetic Operators](#arithmetic-operators)
4. [Using Arithmetic in Rule Conditions](#using-arithmetic-in-rule-conditions)
5. [Using Arithmetic in Actions](#using-arithmetic-in-actions)
6. [Type Conversions](#type-conversions)
7. [Edge Cases and Limitations](#edge-cases-and-limitations)
8. [Operator Encoding](#operator-encoding)
9. [Examples](#examples)
10. [Best Practices](#best-practices)

---

## Overview

The TSD RETE engine supports arithmetic expressions in both:
- **Rule premises (conditions)**: Used to filter and match facts
- **Rule consequents (actions)**: Used to compute values when rules fire

Arithmetic expressions can include:
- Field accesses from bound facts (`p.price`, `c.quantity`)
- Numeric literals (`10`, `3.14`, `-5`)
- Nested expressions with operator precedence
- Multiple operators in complex calculations

---

## TSD File Syntax

### File Structure

A `.tsd` file contains:

```
// Comments start with //

// 1. Type definitions
type TypeName(field1: type1, field2: type2, ...)

// 2. Action definitions
action action_name(param1: type1, param2: type2, ...)

// 3. Rule definitions
rule rule_name : {var1: Type1, var2: Type2} /
    condition_expression
    ==> action_call(arg1, arg2, ...)

// 4. Fact instances (optional)
TypeName(field1: value1, field2: value2, ...)
```

### Type Definitions

```tsd
type Product(id: string, price: number, weight: number, stock: number)
type Order(id: string, product_id: string, quantity: number, discount: number)
type Customer(id: string, tier: string, credit_limit: number)
```

**Supported field types:**
- `string`: Text values
- `number`: Integer or floating-point numbers
- `boolean`: True/false values

### Action Definitions

```tsd
action process_order(
    order_id: string,
    total_price: number,
    discounted_price: number,
    tax_amount: number
)

action log_event(message: string, severity: number)
```

### Rule Definitions

```tsd
rule calculate_invoice : {p: Product, o: Order} /
    o.product_id == p.id AND o.quantity > 0
    ==> process_order(
        o.id,
        p.price * o.quantity,
        (p.price * o.quantity) * (1 - o.discount / 100),
        (p.price * o.quantity) * 0.08
    )
```

### Fact Instances

Facts can be defined directly in `.tsd` files:

```tsd
Product(id: "PROD001", price: 99.99, weight: 2.5, stock: 100)
Order(id: "ORD001", product_id: "PROD001", quantity: 5, discount: 10)
Customer(id: "CUST001", tier: "gold", credit_limit: 5000)
```

**Important notes:**
- No `fact` keyword is required (unlike some rule languages)
- The parser may emit facts with a `reteType` key instead of `type`
- Fact IDs are automatically copied into the `Fields` map for join conditions

---

## Arithmetic Operators

### Supported Operators

| Operator | Description | Example | Result |
|----------|-------------|---------|--------|
| `+` | Addition | `5 + 3` | `8` |
| `-` | Subtraction | `10 - 4` | `6` |
| `*` | Multiplication | `6 * 7` | `42` |
| `/` | Division | `20 / 4` | `5` |
| `%` | Modulo (remainder) | `17 % 5` | `2` |

### Operator Precedence

Arithmetic expressions follow standard mathematical precedence:

1. **Highest**: `*`, `/`, `%` (multiplication, division, modulo)
2. **Lowest**: `+`, `-` (addition, subtraction)

Use parentheses to override precedence:

```tsd
// Without parentheses: 2 + 3 * 4 = 2 + (3 * 4) = 14
// With parentheses: (2 + 3) * 4 = 5 * 4 = 20
```

**Note**: The parser handles operator precedence during parsing. The evaluator processes the parsed tree structure.

---

## Using Arithmetic in Rule Conditions

Arithmetic expressions can appear in rule conditions (premises) to filter facts.

### Basic Comparison with Arithmetic

```tsd
rule expensive_orders : {p: Product, o: Order} /
    o.product_id == p.id AND (p.price * o.quantity) > 1000
    ==> log_event("High value order", 2)
```

### Complex Condition Example

```tsd
rule bulk_discount_eligible : {p: Product, o: Order, c: Customer} /
    o.product_id == p.id AND
    o.customer_id == c.id AND
    (p.price * o.quantity * 0.9) < c.credit_limit AND
    o.quantity >= 10
    ==> apply_bulk_discount(o.id, p.price * o.quantity * 0.1)
```

### Nested Arithmetic in Conditions

```tsd
rule shipping_cost_check : {p: Product, o: Order} /
    o.product_id == p.id AND
    ((p.weight * o.quantity) + 10) * 2.5 > 100
    ==> flag_for_review(o.id)
```

---

## Using Arithmetic in Actions

Action arguments can contain arithmetic expressions that are evaluated when the rule fires.

### Simple Calculation

```tsd
action invoice(order_id: string, total: number, tax: number)

rule generate_invoice : {p: Product, o: Order} /
    o.product_id == p.id
    ==> invoice(
        o.id,
        p.price * o.quantity,
        (p.price * o.quantity) * 0.08
    )
```

### Complex Nested Expressions

```tsd
action billing(
    order_id: string,
    subtotal: number,
    discount_amount: number,
    tax_amount: number,
    grand_total: number,
    shipping: number
)

rule calculate_billing : {p: Product, o: Order} /
    o.product_id == p.id AND o.quantity > 0
    ==> billing(
        o.id,
        p.price * o.quantity,
        (p.price * o.quantity) * (o.discount / 100),
        (p.price * o.quantity) * (1 - o.discount / 100) * 0.08,
        (p.price * o.quantity) * (1 - o.discount / 100) * 1.08,
        (p.weight * o.quantity) + 10
    )
```

### Combining Variables and Literals

```tsd
rule special_pricing : {p: Product, o: Order} /
    o.product_id == p.id
    ==> calculate_price(
        o.id,
        p.price * o.quantity * 0.9,           // 10% discount
        (p.price * o.quantity * 0.9) + 5.99,  // Add fixed shipping
        p.price * 3 * 2 + 1                   // Complex calculation
    )
```

---

## Type Conversions

### Automatic Conversions

The evaluator automatically handles numeric type conversions:

```tsd
// Integer + Float -> Float
5 + 2.5 = 7.5

// Float * Integer -> Float
2.5 * 3 = 7.5

// Integer / Integer -> Float (if result is fractional)
10 / 4 = 2.5
```

### Supported Numeric Types

The engine supports:
- `int`, `int32`, `int64`
- `float32`, `float64`
- `uint`, `uint32`, `uint64`

All are converted to `float64` for arithmetic operations.

---

## Edge Cases and Limitations

### Division by Zero

Division or modulo by zero results in an error:

```tsd
// This will cause an error if denominator is 0
p.price / denominator

// This will cause an error if modulus is 0
p.price % modulus
```

**Recommendation**: Ensure denominators are non-zero in your conditions:

```tsd
rule safe_division : {p: Product, o: Order} /
    o.product_id == p.id AND o.quantity > 0
    ==> calculate(p.price / o.quantity)  // Safe: quantity > 0
```

### Invalid Type Combinations

Arithmetic operations require numeric operands. These will fail:

```tsd
// ❌ Invalid: String + Number
o.id + 5

// ❌ Invalid: Boolean * Number
o.is_active * p.price

// ✅ Valid: Number + Number
p.price + o.quantity
```

### Floating-Point Precision

Be aware of floating-point precision limitations:

```tsd
// May have small rounding errors
0.1 + 0.2 = 0.30000000000000004

// For currency, consider storing cents as integers:
// price_in_cents / 100
```

### Nil/Null Values

Arithmetic operations on `nil` values will fail:

```tsd
// If p.price is nil, this will error
p.price * o.quantity
```

**Recommendation**: Use conditions to filter out nil values:

```tsd
rule safe_calc : {p: Product, o: Order} /
    o.product_id == p.id AND
    p.price != nil AND
    o.quantity != nil
    ==> calculate(p.price * o.quantity)
```

---

## Operator Encoding

### Base64 Encoding Issue

The parser sometimes encodes operators as Base64 strings:

| Operator | Base64 Encoding |
|----------|-----------------|
| `+` | `Kw==` |
| `-` | `LQ==` |
| `*` | `Kg==` |
| `/` | `Lw==` |
| `%` | `JQ==` |
| `==` | `PT0=` |
| `!=` | `IT0=` |
| `<` | `PA==` |
| `<=` | `PD0=` |
| `>` | `Pg==` |
| `>=` | `Pj0=` |

### Automatic Decoding

The RETE engine automatically decodes Base64-encoded operators using the centralized `operator_utils.go` utility:

```go
// Automatic decoding happens in:
// - evaluator_values.go
// - evaluator_expressions.go
// - action_executor.go
```

**Note**: This is transparent to users writing `.tsd` files. The encoding happens during parsing, and decoding happens during evaluation.

---

## Examples

### Example 1: E-commerce Order Processing

```tsd
type Product(id: string, price: number, weight: number)
type Order(id: string, product_id: string, qty: number, discount: number)

action process_order(
    order_id: string,
    gross_total: number,
    discount_amount: number,
    net_total: number,
    shipping_cost: number
)

rule calculate_order : {p: Product, o: Order} /
    o.product_id == p.id AND o.qty > 0
    ==> process_order(
        o.id,
        p.price * o.qty,
        (p.price * o.qty) * (o.discount / 100),
        (p.price * o.qty) * (1 - o.discount / 100),
        (p.weight * o.qty) + 10
    )

// Facts
Product(id: "PROD001", price: 100, weight: 2)
Order(id: "ORD001", product_id: "PROD001", qty: 5, discount: 10)
```

**Result when rule fires:**
- `gross_total`: 100 * 5 = 500
- `discount_amount`: 500 * 0.1 = 50
- `net_total`: 500 * 0.9 = 450
- `shipping_cost`: (2 * 5) + 10 = 20

### Example 2: Tiered Pricing

```tsd
type Product(id: string, base_price: number)
type Order(id: string, product_id: string, quantity: number)

action apply_pricing(order_id: string, unit_price: number, total: number)

// Small orders: full price
rule small_order_pricing : {p: Product, o: Order} /
    o.product_id == p.id AND o.quantity < 10
    ==> apply_pricing(o.id, p.base_price, p.base_price * o.quantity)

// Medium orders: 10% discount
rule medium_order_pricing : {p: Product, o: Order} /
    o.product_id == p.id AND o.quantity >= 10 AND o.quantity < 50
    ==> apply_pricing(
        o.id,
        p.base_price * 0.9,
        (p.base_price * 0.9) * o.quantity
    )

// Large orders: 20% discount
rule large_order_pricing : {p: Product, o: Order} /
    o.product_id == p.id AND o.quantity >= 50
    ==> apply_pricing(
        o.id,
        p.base_price * 0.8,
        (p.base_price * 0.8) * o.quantity
    )

Product(id: "P1", base_price: 100)
Order(id: "O1", product_id: "P1", quantity: 5)    // Small order
Order(id: "O2", product_id: "P1", quantity: 25)   // Medium order
Order(id: "O3", product_id: "P1", quantity: 100)  // Large order
```

### Example 3: Complex Calculation with Multiple Operations

```tsd
type Material(id: string, cost_per_kg: number, density: number)
type Production(id: string, material_id: string, volume: number, hours: number)

action calculate_cost(
    production_id: string,
    material_cost: number,
    labor_cost: number,
    overhead: number,
    total_cost: number
)

rule production_costing : {m: Material, pr: Production} /
    pr.material_id == m.id AND pr.volume > 0
    ==> calculate_cost(
        pr.id,
        // Material cost = volume * density * cost per kg
        pr.volume * m.density * m.cost_per_kg,
        // Labor cost = hours * rate (50/hour) with overtime multiplier
        pr.hours * 50 * (1 + (pr.hours - 40) * 0.5 / 40),
        // Overhead = 20% of material + labor
        (pr.volume * m.density * m.cost_per_kg + pr.hours * 50) * 0.2,
        // Total
        (pr.volume * m.density * m.cost_per_kg) +
        (pr.hours * 50 * (1 + (pr.hours - 40) * 0.5 / 40)) +
        ((pr.volume * m.density * m.cost_per_kg + pr.hours * 50) * 0.2)
    )

Material(id: "MAT001", cost_per_kg: 15, density: 1.2)
Production(id: "PROD001", material_id: "MAT001", volume: 10, hours: 45)
```

---

## Best Practices

### 1. Use Parentheses for Clarity

Even when not strictly necessary, parentheses improve readability:

```tsd
// Less clear
action calc(result: number)
rule r : {x: X} / true ==> calc(x.a + x.b * x.c - x.d / x.e)

// More clear
rule r : {x: X} / true ==> calc((x.a + (x.b * x.c)) - (x.d / x.e))
```

### 2. Validate Denominators

Always ensure division/modulo denominators are non-zero:

```tsd
rule safe_average : {p: Product, o: Order} /
    o.product_id == p.id AND o.quantity > 0
    ==> calculate_avg(p.price / o.quantity)
```

### 3. Consider Precision Requirements

For financial calculations, be aware of floating-point limitations:

```tsd
// For exact currency: store as cents and convert
// price_cents = 9999  (means $99.99)
// total_cents = price_cents * quantity
// total_dollars = total_cents / 100
```

### 4. Document Complex Expressions

Use comments to explain non-obvious calculations:

```tsd
rule calculate_shipping : {p: Product, o: Order} /
    o.product_id == p.id
    ==> shipping_cost(
        o.id,
        // Base shipping: $10 + weight-based charge ($2.50/kg)
        ((p.weight * o.quantity) + 10) * 2.5
    )
```

### 5. Test Edge Cases

Always test with:
- Zero values
- Negative numbers
- Very large numbers
- Very small numbers (for floating-point)
- Boundary conditions

### 6. Avoid Overly Complex Expressions

If an expression becomes too complex, consider:
- Breaking it into multiple rules
- Pre-computing values in upstream systems
- Creating intermediate facts

```tsd
// Instead of one massive expression in a rule:
// total = ((a + b) * (c - d) / (e + f)) * ((g * h) - (i / j)) + ...

// Consider breaking into steps:
rule step1 : ... ==> intermediate1(...)
rule step2 : ... ==> intermediate2(...)
rule final : ... ==> result(...)
```

---

## Additional Resources

- **Action System Documentation**: See `ACTIONS_README.md` for action execution details
- **RETE Algorithm**: See main `README.md` for RETE network architecture
- **Testing**: See `action_arithmetic_e2e_test.go` for comprehensive examples
- **Operator Utilities**: See `operator_utils.go` for operator handling internals

---

## Troubleshooting

### Error: "division by zero"

**Cause**: Dividing or taking modulo by zero.

**Solution**: Add condition to ensure denominator is non-zero:
```tsd
rule safe : {x: X} / x.denominator != 0 ==> calc(x.numerator / x.denominator)
```

### Error: "operator not supported"

**Cause**: Unknown or unsupported operator.

**Solution**: Verify you're using supported operators: `+`, `-`, `*`, `/`, `%`

### Error: "arithmetic operation requires numbers"

**Cause**: Attempting arithmetic on non-numeric types.

**Solution**: Ensure operands are numeric fields or literals:
```tsd
// ❌ Wrong: o.id * 5  (id is string)
// ✅ Right: o.quantity * 5  (quantity is number)
```

### Unexpected Results

**Cause**: Operator precedence or type conversion issues.

**Solution**:
- Add explicit parentheses
- Check data types of fields
- Add debug logging to inspect intermediate values

---

## Version History

- **v1.0** (2025-01): Initial arithmetic expression support
  - Basic operators in conditions and actions
  - Base64 operator decoding
  - Type conversion handling
  
- **v1.1** (2025-01): Centralized operator utilities
  - Created `operator_utils.go` for consistent operator handling
  - Added comprehensive edge-case tests
  - Improved error messages

---

**Copyright © 2025 TSD Contributors**  
Licensed under the MIT License