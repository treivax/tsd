# Multi-Source Aggregation Support

## Overview

The TSD rule engine now supports **multi-source aggregations** - the ability to aggregate data from multiple joined fact types in a single rule. This feature extends the existing single-source aggregation capability to handle complex analytical queries across related entities.

## Syntax

### Basic Multi-Source Aggregation

```
rule rule_name : {main: MainType, agg1: AGG_FUNC1(var1.field1), agg2: AGG_FUNC2(var2.field2)} 
  / {var1: Type1} 
  / {var2: Type2} 
  / join_conditions 
  ==> action
```

### Example: Two Sources

```
type Department : <id: string, name: string>
type Employee : <id: string, deptId: string, salary: number>
type Performance : <id: string, employeeId: string, score: number>

rule dept_combined_metric : 
  {d: Department, avg_sal: AVG(e.salary), avg_score: AVG(p.score)} 
  / {e: Employee} 
  / {p: Performance} 
  / e.deptId == d.id AND p.employeeId == e.id 
  ==> print("Combined metrics")
```

This rule:
- Joins `Department`, `Employee`, and `Performance` facts
- Computes average salary from Employee records
- Computes average performance score from Performance records
- Fires the action when facts match the join conditions

### Example: Three Sources

```
type Department : <id: string>
type Employee : <id: string, deptId: string, salary: number>
type Performance : <id: string, employeeId: string, score: number>
type Training : <id: string, employeeId: string, hours: number>

rule dept_full_metrics : 
  {d: Department, 
   avg_sal: AVG(e.salary), 
   avg_score: AVG(p.score), 
   total_hours: SUM(t.hours)} 
  / {e: Employee} 
  / {p: Performance} 
  / {t: Training} 
  / e.deptId == d.id AND p.employeeId == e.id AND t.employeeId == e.id 
  ==> print("Full metrics")
```

## Supported Aggregation Functions

All standard aggregation functions are supported for each aggregation variable:

- `AVG(var.field)` - Average of numeric values
- `SUM(var.field)` - Sum of numeric values
- `COUNT(var.field)` - Count of records
- `MIN(var.field)` - Minimum value
- `MAX(var.field)` - Maximum value

## Pattern Structure

### First Pattern Block (Main + Aggregation Variables)

```
{main_var: MainType, agg_var1: AGG_FUNC1(source1.field1), agg_var2: AGG_FUNC2(source2.field2), ...}
```

- **main_var**: The primary entity around which aggregation occurs
- **agg_var1, agg_var2, ...**: Aggregation variables that compute metrics from joined sources

### Source Pattern Blocks

```
/ {source1: Type1}
/ {source2: Type2}
/ {source3: Type3}
...
```

Each source pattern block declares a fact type that will be joined and potentially aggregated.

### Join Conditions

```
/ source1.join_field == main_var.field AND source2.join_field == source1.other_field AND ...
```

Join conditions specify how the sources relate to each other using logical AND operations.

## Threshold Support

You can add threshold conditions on aggregation results:

```
rule high_performing_dept : 
  {d: Department, avg_sal: AVG(e.salary), avg_score: AVG(p.score)} 
  / {e: Employee} 
  / {p: Performance} 
  / e.deptId == d.id AND p.employeeId == e.id 
    AND avg_sal > 50000 AND avg_score > 80 
  ==> print("High performing department")
```

The rule only fires when:
1. The join conditions are satisfied
2. The average salary exceeds 50,000
3. The average performance score exceeds 80

### Supported Threshold Operators

- `>` - Greater than
- `>=` - Greater than or equal
- `<` - Less than
- `<=` - Less than or equal
- `==` - Equal to
- `!=` - Not equal to

## Implementation Details

### Parser Support

The constraint parser recognizes:
- Multiple aggregation variables in the first pattern block
- Multiple source pattern blocks (via `/` separator)
- Complex join conditions with `AND`/`OR` operators
- Mixed join conditions and threshold conditions

### RETE Network Construction

For multi-source aggregations, the pipeline:

1. **Detects** multi-source pattern (multiple aggregation vars or 3+ pattern blocks)
2. **Extracts** all aggregation variables with their functions and source fields
3. **Identifies** all source patterns and their types
4. **Parses** join conditions between patterns
5. **Builds** a chain of JoinNodes to connect all sources
6. **Creates** terminal nodes to fire actions when conditions are met

### Current Implementation Status

**✅ Completed:**
- Parser grammar supports multi-source syntax
- AST extraction for multiple aggregation variables
- AST extraction for multiple source patterns
- Join condition extraction and separation
- Multi-join chain construction in RETE network
- Basic join-based activation

**🚧 In Progress:**
- Multi-source AccumulatorNode implementation
- Aggregation computation across multiple joined sources
- Threshold evaluation on computed aggregates

**📋 Planned:**
- Performance optimization for large multi-join scenarios
- Incremental aggregation updates
- Memory management for join results

## Use Cases

### Business Analytics

Track comprehensive metrics across organizational structures:

```
rule executive_dashboard :
  {exec: Executive,
   dept_count: COUNT(d.id),
   avg_team_size: AVG(t.member_count),
   total_revenue: SUM(p.revenue)}
  / {d: Department}
  / {t: Team}
  / {p: Project}
  / d.executiveId == exec.id 
    AND t.departmentId == d.id 
    AND p.teamId == t.id
  ==> update_dashboard(exec.id)
```

### Performance Monitoring

Combine operational metrics from multiple sources:

```
rule service_health :
  {svc: Service,
   avg_latency: AVG(r.latency),
   error_rate: AVG(e.rate),
   throughput: SUM(r.count)}
  / {r: Request}
  / {e: Error}
  / r.serviceId == svc.id 
    AND e.serviceId == svc.id
    AND avg_latency < 100 
    AND error_rate < 0.01
  ==> alert_success(svc.id)
```

### Complex Relationship Analysis

Analyze patterns across interconnected entities:

```
rule project_risk_assessment :
  {proj: Project,
   avg_task_delay: AVG(t.delay_days),
   team_turnover: COUNT(tc.employee_id),
   budget_variance: SUM(e.amount - e.budgeted)}
  / {t: Task}
  / {tc: TeamChange}
  / {e: Expense}
  / t.projectId == proj.id
    AND tc.projectId == proj.id
    AND e.projectId == proj.id
    AND avg_task_delay > 5
  ==> escalate_risk(proj.id)
```

## Testing

The feature includes comprehensive test coverage:

### Parser Tests (`constraint/multi_source_aggregation_test.go`)
- ✅ Two-source aggregation syntax
- ✅ Three-source aggregation syntax
- ✅ Multiple thresholds
- ✅ Different aggregation functions (AVG, SUM, COUNT, MIN, MAX)
- ✅ Error cases and validation
- ✅ Complex join conditions

### RETE Tests (`rete/multi_source_aggregation_test.go`)
- ✅ Two-source join chain construction
- ✅ Three-source join chain construction
- ⚠️ Threshold evaluation (partial - needs aggregation logic)
- ✅ Different function types

## Migration Guide

### From Single-Source to Multi-Source

**Before (Single Source):**
```
rule dept_avg : 
  {d: Department, avg_sal: AVG(e.salary)} 
  / {e: Employee} 
  / e.deptId == d.id 
  ==> print("Average")
```

**After (Multi-Source):**
```
rule dept_combined : 
  {d: Department, avg_sal: AVG(e.salary), avg_score: AVG(p.score)} 
  / {e: Employee} 
  / {p: Performance} 
  / e.deptId == d.id AND p.employeeId == e.id 
  ==> print("Combined")
```

The syntax is backward-compatible - single-source rules continue to work unchanged.

## Performance Considerations

1. **Join Order**: The system joins sources in the order specified. Consider placing selective sources earlier.

2. **Fact Volume**: Multi-source joins can produce many intermediate results. Monitor memory usage with high fact counts.

3. **Index Usage**: Ensure join fields are efficiently accessible in fact structures.

4. **Incremental Updates**: Currently, join chains rebuild on fact changes. Future optimizations will support incremental updates.

## Limitations

1. **Aggregation Implementation**: Current version creates join chains but full aggregation computation is in progress.

2. **Threshold Evaluation**: Threshold checking on aggregated values needs completion of aggregation logic.

3. **Memory Management**: Large multi-join results currently stored in memory without eviction.

4. **Cycle Detection**: No validation against cyclic join patterns (e.g., A→B→C→A).

## Future Enhancements

- **Window Functions**: Time-based or count-based aggregation windows
- **Nested Aggregations**: Aggregations over aggregation results
- **Dynamic Thresholds**: Compare aggregation results against fact field values
- **Streaming Aggregation**: Real-time aggregation updates with fact stream processing
- **Parallel Execution**: Concurrent aggregation computation across partitions

## References

- [Single-Source Aggregation](./aggregation-with-joins.md)
- [RETE Algorithm](./rete-architecture.md)
- [Threshold Support](./aggregation-thresholds.md)
- [Parser Grammar](../constraint/grammar/constraint.peg)