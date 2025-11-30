# Multi-Source Aggregation - Quick Reference

## Basic Syntax

```
rule rule_name : 
  {main: MainType, agg1: FUNC1(var1.field1), agg2: FUNC2(var2.field2)} 
  / {var1: Type1} 
  / {var2: Type2} 
  / join_conditions 
  ==> action
```

## Supported Functions

| Function | Description | Example |
|----------|-------------|---------|
| `AVG(var.field)` | Average of values | `avg_sal: AVG(e.salary)` |
| `SUM(var.field)` | Sum of values | `total_sal: SUM(e.salary)` |
| `COUNT(var.field)` | Count of records | `emp_count: COUNT(e.id)` |
| `MIN(var.field)` | Minimum value | `min_sal: MIN(e.salary)` |
| `MAX(var.field)` | Maximum value | `max_sal: MAX(e.salary)` |

## Quick Examples

### Two Sources - Basic
```
type Department : <id: string, name: string>
type Employee : <id: string, deptId: string, salary: number>
type Performance : <id: string, employeeId: string, score: number>

rule dept_metrics : 
  {d: Department, avg_sal: AVG(e.salary), avg_score: AVG(p.score)} 
  / {e: Employee} 
  / {p: Performance} 
  / e.deptId == d.id AND p.employeeId == e.id 
  ==> print("Dept", d.name, "metrics computed")
```

### With Thresholds
```
rule high_performers : 
  {d: Department, avg_sal: AVG(e.salary), avg_score: AVG(p.score)} 
  / {e: Employee} 
  / {p: Performance} 
  / e.deptId == d.id AND p.employeeId == e.id 
    AND avg_sal > 50000 AND avg_score > 85 
  ==> print("High performing:", d.name)
```

### Three Sources
```
rule comprehensive_stats : 
  {d: Department, 
   emp_count: COUNT(e.id),
   avg_score: AVG(p.score), 
   total_hours: SUM(t.hours)} 
  / {e: Employee} 
  / {p: Performance} 
  / {t: Training} 
  / e.deptId == d.id 
    AND p.employeeId == e.id 
    AND t.employeeId == e.id 
  ==> print("Stats for", d.name)
```

### Multiple Functions on Same Source
```
rule salary_range : 
  {d: Department, 
   min_sal: MIN(e.salary), 
   max_sal: MAX(e.salary), 
   avg_sal: AVG(e.salary)} 
  / {e: Employee} 
  / e.deptId == d.id 
  ==> print("Range:", min_sal, "-", max_sal)
```

## Threshold Operators

| Operator | Meaning |
|----------|---------|
| `>` | Greater than |
| `>=` | Greater than or equal |
| `<` | Less than |
| `<=` | Less than or equal |
| `==` | Equal to |
| `!=` | Not equal to |

## Pattern Structure

### First Pattern Block
```
{main_var: MainType, agg_var1: FUNC1(...), agg_var2: FUNC2(...)}
```
- Main variable (mandatory)
- One or more aggregation variables

### Source Pattern Blocks
```
/ {source1: Type1}
/ {source2: Type2}
/ {source3: Type3}
```
- One block per source type
- Each declares a fact type to join

### Join Conditions
```
/ source1.field == main_var.field 
  AND source2.field == source1.other_field
  AND agg_var1 > threshold1
  AND agg_var2 < threshold2
```
- Join conditions link sources together
- Threshold conditions filter on aggregated values
- Use `AND` or `OR` to combine conditions

## Common Patterns

### Department Analytics
```
rule dept_overview :
  {d: Department,
   emp_count: COUNT(e.id),
   avg_salary: AVG(e.salary),
   total_budget: SUM(e.salary)}
  / {e: Employee}
  / e.deptId == d.id
  ==> update_dashboard(d.id)
```

### Performance Monitoring
```
rule service_health :
  {svc: Service,
   avg_latency: AVG(r.latency),
   error_count: COUNT(e.id)}
  / {r: Request}
  / {e: Error}
  / r.serviceId == svc.id 
    AND e.serviceId == svc.id
    AND avg_latency < 100
  ==> mark_healthy(svc.id)
```

### Budget Analysis
```
rule over_budget :
  {dept: Department,
   total_spend: SUM(e.salary)}
  / {e: Employee}
  / e.deptId == dept.id
    AND total_spend > dept.budget
  ==> alert("Over budget", dept.name)
```

## How It Works

1. **Facts Submitted** → TypeNodes
2. **Join Chain** → Combines facts from all sources
3. **Accumulator** → Computes all aggregations
4. **Thresholds** → Filters results
5. **Action** → Fires only if all thresholds satisfied

## Tips & Best Practices

### ✅ Do
- Use descriptive aggregation variable names
- Group related aggregations in one rule
- Add thresholds to filter unwanted results
- Document complex join conditions

### ❌ Don't
- Mix unrelated aggregations in one rule
- Create circular join dependencies
- Forget to link all sources with join conditions
- Use aggregation variable names that conflict with source variables

## Debugging

### Check Join Conditions
Ensure all sources are connected:
```
✅ Good: e.deptId == d.id AND p.employeeId == e.id
❌ Bad:  e.deptId == d.id (p is orphaned)
```

### Verify Field Names
Make sure fields exist in fact types:
```
✅ Good: AVG(e.salary) where Employee has salary field
❌ Bad:  AVG(e.pay) where Employee has no pay field
```

### Test Without Thresholds First
Start simple, add thresholds after:
```
Step 1: {d: Dept, avg: AVG(e.sal)} / {e: Emp} / e.deptId == d.id
Step 2: Add threshold: ... AND avg > 50000
```

## Performance Considerations

- **Small datasets (<1000 facts):** No optimization needed
- **Medium datasets (1000-10000):** Consider selective join order
- **Large datasets (>10000):** Monitor memory usage, consider windowing

## Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| "Variable not found" | Misspelled variable name | Check spelling |
| "Field not found" | Field doesn't exist | Verify type definition |
| "No activation" | Thresholds too strict | Relax conditions |
| "Too many activations" | Missing threshold | Add filtering conditions |

## Migration from Single-Source

**Before (Single-Source):**
```
rule r1 : {d: Dept, avg: AVG(e.sal)} / {e: Emp} / e.deptId == d.id
```

**After (Multi-Source - just add more!):**
```
rule r1 : 
  {d: Dept, avg_sal: AVG(e.sal), avg_score: AVG(p.score)} 
  / {e: Emp} 
  / {p: Perf} 
  / e.deptId == d.id AND p.empId == e.id
```

## See Also

- [Full Documentation](./multi-source-aggregation.md)
- [Implementation Details](./MULTI_SOURCE_AGGREGATION_SUMMARY.md)
- [Examples](../examples/multi_source_aggregation.tsd)
- [Changelog](./CHANGELOG_MULTI_SOURCE_AGG.md)

## Questions?

Common questions answered in full documentation:
- How are ties handled in MIN/MAX?
- What happens with NULL values?
- Can I nest aggregations?
- How does retraction work?
- Performance tuning tips