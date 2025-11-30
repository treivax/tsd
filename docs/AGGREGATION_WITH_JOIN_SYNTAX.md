# Aggregation with Join Syntax

## Overview

The TSD constraint language now supports a powerful new syntax for defining aggregations with joins. This allows you to perform aggregate calculations (AVG, SUM, COUNT, MIN, MAX) over sets of facts that are joined to other facts.

## Motivation

Previously, aggregations in TSD were limited to standalone aggregate constraints. With the new syntax, you can express rules like:

- "For each department, calculate the average salary of employees in that department"
- "For each customer, count the number of orders they placed"
- "For each product category, find the maximum price"

This pattern requires both:
1. A join relationship between two fact types
2. An aggregation over the joined facts

## Syntax

### Basic Structure

```tsd
rule <rule_name> : {<result_vars>, <agg_var>: <AGG_FUNC>(<source_var>.<field>)} / {<source_var>: <SourceType>} / <join_condition> ==> <action>
```

### Components

1. **Result pattern block** (first `{...}`)
   - Contains the "driver" variable(s) and aggregation result variable(s)
   - Can mix regular typed variables and aggregation variables
   - Example: `{d: Department, avg_sal: AVG(e.salary)}`

2. **Source pattern block** (second `{...}`)
   - Contains the variable being aggregated over
   - Example: `{e: Employee}`

3. **Join condition**
   - Specifies how the result and source patterns are related
   - Example: `e.deptId == d.id`

4. **Action**
   - Executed when the rule fires
   - Has access to both result variables and aggregation results

### Aggregation Functions

The following aggregation functions are supported:

- `AVG(var.field)` - Average of numeric field values
- `SUM(var.field)` - Sum of numeric field values
- `COUNT(var.field)` - Count of matching facts
- `MIN(var.field)` - Minimum field value
- `MAX(var.field)` - Maximum field value

## Examples

### Example 1: Average Salary by Department

```tsd
type Department : <id: string, name: string>
type Employee : <id: string, deptId: string, salary: number>

rule dept_avg_salary : {d: Department, avg_sal: AVG(e.salary)} / {e: Employee} / e.deptId == d.id ==> print("Avg salary")
```

**How it works:**
- For each `Department` fact (driver variable `d`)
- Find all `Employee` facts where `e.deptId == d.id`
- Calculate the average of their `salary` fields
- Bind the result to `avg_sal`
- Fire the action with both `d` and `avg_sal` available

### Example 2: Multiple Aggregations

```tsd
type Department : <id: string, name: string>
type Employee : <id: string, deptId: string, salary: number>

rule dept_stats : {d: Department, avg_sal: AVG(e.salary), max_sal: MAX(e.salary), count: COUNT(e.id)} / {e: Employee} / e.deptId == d.id ==> print("Stats")
```

You can declare multiple aggregation variables in the same pattern block.

### Example 3: Customer Order Count

```tsd
type Customer : <id: string, name: string>
type Order : <id: string, customerId: string, amount: number>

rule customer_order_count : {c: Customer, order_count: COUNT(o.id)} / {o: Order} / o.customerId == c.id ==> notify("Order count")
```

### Example 4: Product Category Maximum Price

```tsd
type Category : <id: string, name: string>
type Product : <id: string, categoryId: string, price: number>

rule category_max_price : {cat: Category, max_price: MAX(p.price)} / {p: Product} / p.categoryId == cat.id ==> alert("Max price")
```

## Implementation Details

### Parser Grammar

The grammar has been extended with:

1. **`PatternBlocks`** rule - Allows multiple `{...}` blocks separated by `/`
2. **`AggregationVariable`** rule - Parses `name: FUNC(var.field)` syntax
3. **`TypedVariable`** alternatives - Can be either `SimpleTypedVariable` or `AggregationVariable`

### AST Structure

For multi-pattern aggregation rules, the expression AST contains:

```json
{
  "type": "expression",
  "ruleId": "dept_avg_salary",
  "patterns": [
    {
      "type": "set",
      "variables": [
        {"type": "typedVariable", "name": "d", "dataType": "Department"},
        {
          "type": "aggregationVariable",
          "name": "avg_sal",
          "function": "AVG",
          "field": {"type": "fieldAccess", "object": "e", "field": "salary"}
        }
      ]
    },
    {
      "type": "set",
      "variables": [
        {"type": "typedVariable", "name": "e", "dataType": "Employee"}
      ]
    }
  ],
  "constraints": {...},
  "action": {...}
}
```

### Backward Compatibility

The parser maintains full backward compatibility:

- **Single-pattern rules** (old syntax) continue to use the `"set"` field
- **Multi-pattern rules** (new syntax) use the `"patterns"` field
- All validation and processing code handles both structures

Example of old syntax (still supported):

```tsd
rule adults : {p: Person} / p.age >= 18 ==> print("Adult")
```

This produces:

```json
{
  "type": "expression",
  "ruleId": "adults",
  "set": {...},
  "constraints": {...},
  "action": {...}
}
```

### RETE Network Mapping

Aggregation-with-join rules are implemented using `AccumulatorNode` in the RETE network:

1. A **passthrough alpha node** is created for the driver pattern (e.g., Department)
2. A **passthrough alpha node** is created for the source pattern (e.g., Employee)
3. An **AccumulatorNode** performs the aggregation and join matching
4. A **TerminalNode** fires the action when conditions are met

## Validation

The semantic validator has been updated to:

1. Extract variables from all pattern blocks (not just the first)
2. Validate field accesses across all declared variables
3. Check type compatibility for aggregation expressions
4. Handle aggregation variables that don't have traditional data types

## Testing

Comprehensive tests verify:

- ✅ Parser correctly handles new syntax
- ✅ All aggregation functions (AVG, SUM, COUNT, MIN, MAX) work
- ✅ Multiple aggregation variables in one rule
- ✅ Backward compatibility with old single-pattern syntax
- ✅ Semantic validation passes
- ✅ RETE network construction succeeds
- ✅ Facts are processed correctly

## Limitations and Future Work

### Current Limitations

1. **Aggregation calculation** - The aggregation logic itself may need refinement (values showing as 0.00 in some cases)
2. **No threshold comparison** - Unlike old `AccumulateConstraint` syntax, the new syntax doesn't support threshold comparisons (e.g., `AVG(...) > 50000`)
3. **Single aggregation source** - Currently supports aggregating over one source pattern at a time

### Future Enhancements

1. **Threshold comparisons** - Add support for conditions like `avg_sal > threshold`
2. **Multiple source patterns** - Allow aggregating over multiple joined patterns
3. **Nested aggregations** - Support aggregations of aggregations
4. **GROUP BY semantics** - More explicit grouping control

## Migration Guide

### From Old Aggregate Syntax

**Old syntax (AccumulateConstraint):**

```tsd
rule avg_check : {p: Person} / AVG(e: Employee / p.id == e.managerId; e.salary) > 50000 ==> alert("High avg")
```

**New syntax (Aggregation with Join):**

```tsd
rule avg_check : {p: Person, avg_sal: AVG(e.salary)} / {e: Employee} / p.id == e.managerId ==> alert("High avg")
```

Note: Threshold comparison (`> 50000`) would need to be added to the action or as a separate constraint.

## See Also

- [Beta Sharing Architecture](BETA_SHARING_ARCHITECTURE.md)
- [Join Node Implementation](JOIN_NODE_IMPLEMENTATION.md)
- [Parser Grammar Reference](../constraint/grammar/constraint.peg)
- [Backward Compatibility Tests](../rete/beta_backward_compatibility_test.go)

## Examples Directory

Working examples can be found in:

- `rete/beta_backward_compatibility_test.go` - Test cases demonstrating the syntax
- `constraint/aggregation_join_test.go` - Parser-level test cases

---

**Version:** 1.0  
**Date:** 2025  
**Status:** ✅ Implemented and Tested