# Incremental Facts Parsing

## Overview

TSD supports **incremental parsing** of facts, allowing you to load data progressively from multiple files or content sources. Facts accumulate in the `ProgramState` across multiple parsing operations, enabling flexible data management strategies.

## Key Features

✅ **Cumulative**: Facts from each parsing operation are added to existing facts  
✅ **Type-safe**: All facts are validated against type definitions  
✅ **Non-blocking**: Invalid facts are rejected with warnings, valid facts are kept  
✅ **File tracking**: Each parsed file is recorded in `FilesParsed`  
✅ **Reset support**: Clear all state with `reset` instruction  
✅ **Mixed content**: Can parse facts alongside types and rules  

## How It Works

### Basic Incremental Flow

```tsd
# File 1: types.tsd
type Person : <id: string, name: string, age: number>

# File 2: facts1.tsd
Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)

# File 3: facts2.tsd
Person(id: "P003", name: "Charlie", age: 35)
Person(id: "P004", name: "Diana", age: 28)
```

**Result**: All 4 facts are accumulated in `ProgramState.Facts`

### Programmatic Usage

```go
ps := constraint.NewProgramState()

// Parse types first
ps.ParseAndMerge("types.tsd")
// State: 1 type, 0 facts

// Parse first batch of facts
ps.ParseAndMerge("facts1.tsd")
// State: 1 type, 2 facts

// Parse second batch of facts
ps.ParseAndMerge("facts2.tsd")
// State: 1 type, 4 facts

// Parse third batch of facts
ps.ParseAndMerge("facts3.tsd")
// State: 1 type, 7 facts
```

## Use Cases

### 1. Batch Data Loading

Load large datasets in smaller batches to manage memory:

```go
ps := constraint.NewProgramState()
ps.ParseAndMerge("types.tsd")

// Load 1000 facts at a time
for i := 0; i < 10; i++ {
    filename := fmt.Sprintf("facts_batch_%d.tsd", i)
    ps.ParseAndMerge(filename)
}
// Total: 10,000 facts loaded incrementally
```

### 2. Multiple Data Sources

Combine facts from different sources:

```go
ps := constraint.NewProgramState()
ps.ParseAndMerge("types.tsd")

// Load from different sources
ps.ParseAndMerge("customers.tsd")      // Customer facts
ps.ParseAndMerge("products.tsd")       // Product facts
ps.ParseAndMerge("orders.tsd")         // Order facts
ps.ParseAndMerge("transactions.tsd")   // Transaction facts
```

### 3. Incremental Updates

Add new facts to existing system:

```go
ps := constraint.NewProgramState()
ps.ParseAndMerge("initial_data.tsd")  // Load base data

// Later... add new data
ps.ParseAndMerge("updates_2025_01.tsd")
ps.ParseAndMerge("updates_2025_02.tsd")
```

### 4. Mixed Content Files

Parse files containing types, rules, and facts together:

```go
ps := constraint.NewProgramState()

// File with types and initial facts
ps.ParseAndMerge("schema_and_seed.tsd")

// File with additional facts
ps.ParseAndMerge("additional_data.tsd")

// File with rules (facts remain intact)
ps.ParseAndMerge("business_rules.tsd")
```

## Validation During Incremental Parsing

Each batch of facts is validated independently:

```tsd
# Batch 1: All valid
Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)
# Result: 2 facts added

# Batch 2: Mixed valid/invalid
Person(id: "P003", name: "Charlie", age: 35)      # ✅ Valid - added
Person(id: "P004", name: "Diana", salary: 50000)  # ❌ Invalid field - rejected
Person(id: "P005", name: "Eve", age: "forty")     # ❌ Wrong type - rejected
# Result: 1 fact added (P003), 2 errors recorded

# Batch 3: All valid again
Person(id: "P006", name: "Frank", age: 50)
# Result: 1 fact added
# Total: 4 facts, 2 errors
```

**Key Point**: Invalid facts in one batch don't affect valid facts in any batch.

## Reset Behavior

The `reset` instruction clears all accumulated state:

```tsd
# File 1: Initial data
type Person : <id: string, name: string, age: number>
Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)
# State: 2 facts

# File 2: Reset and new data
reset

type Person : <id: string, name: string, age: number>
Person(id: "P100", name: "Xavier", age: 40)
# State: 1 fact (previous 2 cleared)
```

**What reset clears:**
- All types
- All rules (and rule IDs)
- All facts
- All errors
- Files parsed list

## Multiple Type Facts

You can incrementally add facts for different types:

```tsd
# types.tsd
type Person : <id: string, name: string, age: number>
type Product : <id: string, name: string, price: number>
type Order : <id: string, customer_id: string, product_id: string>

# people.tsd
Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)

# products.tsd
Product(id: "PR001", name: "Laptop", price: 999.99)
Product(id: "PR002", name: "Mouse", price: 29.99)

# orders.tsd
Order(id: "O001", customer_id: "P001", product_id: "PR001")
Order(id: "O002", customer_id: "P002", product_id: "PR002")
```

**Result**: All facts coexist in the same `ProgramState`:
- 2 Person facts
- 2 Product facts
- 2 Order facts

## Facts and Rules Together

Facts can be added before, after, or alongside rules:

```go
ps := constraint.NewProgramState()

// Scenario 1: Rules first, facts later
ps.ParseAndMerge("types_and_rules.tsd")  // Types + Rules
ps.ParseAndMerge("facts.tsd")             // Facts

// Scenario 2: Facts first, rules later
ps.ParseAndMerge("types_and_facts.tsd")  // Types + Facts
ps.ParseAndMerge("rules.tsd")             // Rules

// Scenario 3: Mixed
ps.ParseAndMerge("schema.tsd")           // Types
ps.ParseAndMerge("initial_facts.tsd")    // Facts
ps.ParseAndMerge("rules.tsd")            // Rules
ps.ParseAndMerge("more_facts.tsd")       // More facts
```

**Rules are preserved** when adding facts incrementally.

## Error Handling

Errors accumulate across batches:

```go
ps := constraint.NewProgramState()
ps.ParseAndMerge("types.tsd")

ps.ParseAndMerge("batch1.tsd")  // 2 valid, 1 invalid
// ps.GetErrorCount() == 1

ps.ParseAndMerge("batch2.tsd")  // 3 valid, 0 invalid
// ps.GetErrorCount() == 1 (unchanged)

ps.ParseAndMerge("batch3.tsd")  // 1 valid, 2 invalid
// ps.GetErrorCount() == 3 (accumulated)

// Check all errors
for _, err := range ps.GetErrors() {
    fmt.Printf("%s: %s\n", err.File, err.Message)
}

// Clear errors if needed
ps.ClearErrors()
```

## Performance Considerations

### Memory Usage

Facts accumulate in memory. For very large datasets:

```go
// Option 1: Process in batches and extract results
for batchFile := range batchFiles {
    ps := constraint.NewProgramState()
    ps.ParseAndMerge("types.tsd")
    ps.ParseAndMerge(batchFile)
    
    // Process facts
    processFacts(ps.Facts)
    
    // ps goes out of scope, memory released
}

// Option 2: Use reset periodically
ps := constraint.NewProgramState()
for batchFile := range batchFiles {
    ps.ParseAndMerge(batchFile)
    
    if len(ps.Facts) > 10000 {
        // Process and reset
        processFacts(ps.Facts)
        ps.ParseAndMergeContent("reset", "reset.tsd")
    }
}
```

### Parsing Performance

Incremental parsing is efficient:
- Each file parsed once
- Validation per fact is O(1) for type lookup
- Field validation is O(n) where n = number of fields
- No re-parsing of previous files

## Testing

The test suite includes comprehensive incremental parsing tests:

### Test Coverage

1. **Single Type Incremental** (`TestIncrementalFactsParsing_SingleType`)
   - Add facts in 3 batches for one type
   - Verify accumulation and uniqueness

2. **Multiple Types** (`TestIncrementalFactsParsing_MultipleTypes`)
   - Add facts for Person, Product, Order types
   - Verify type distribution

3. **File-Based** (`TestIncrementalFactsParsing_WithFiles`)
   - Use actual files on disk
   - Verify file tracking

4. **With Invalid Facts** (`TestIncrementalFactsParsing_WithInvalidFacts`)
   - Mix valid and invalid facts
   - Verify non-blocking behavior

5. **With Reset** (`TestIncrementalFactsParsing_WithReset`)
   - Test reset clears previous facts
   - Verify fresh start

6. **Large Scale** (`TestIncrementalFactsParsing_LargeScale`)
   - 100 facts in 10 batches
   - Performance verification

7. **Mixed With Rules** (`TestIncrementalFactsParsing_MixedWithRules`)
   - Add facts while rules exist
   - Verify rules preserved

### Running Tests

```bash
# All incremental tests
go test -v -run TestIncrementalFactsParsing

# Specific test
go test -v -run TestIncrementalFactsParsing_SingleType

# Include large scale test
go test -v -run TestIncrementalFactsParsing_LargeScale
```

## Best Practices

### 1. Define Types First

Always parse type definitions before facts:

```go
// ✅ GOOD
ps.ParseAndMerge("types.tsd")
ps.ParseAndMerge("facts.tsd")

// ❌ BAD - will cause validation errors
ps.ParseAndMerge("facts.tsd")
ps.ParseAndMerge("types.tsd")
```

### 2. Check Errors After Each Batch

Monitor validation errors incrementally:

```go
ps := constraint.NewProgramState()
ps.ParseAndMerge("types.tsd")

for _, batchFile := range batchFiles {
    beforeCount := ps.GetErrorCount()
    ps.ParseAndMerge(batchFile)
    afterCount := ps.GetErrorCount()
    
    if afterCount > beforeCount {
        fmt.Printf("⚠️  %d errors in %s\n", 
            afterCount-beforeCount, batchFile)
    }
}
```

### 3. Use Meaningful File Names

Organize files logically:

```
data/
  ├── 00_types.tsd           # Types first
  ├── 01_initial_facts.tsd   # Base data
  ├── 02_supplemental.tsd    # Additional data
  └── 03_latest_updates.tsd  # Recent additions
```

### 4. Document Data Sources

Add comments in files:

```tsd
# Customer data from CRM export - 2025-01-15
Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)

# Product catalog - 2025-01-20
Product(id: "PR001", name: "Laptop", price: 999.99)
```

## API Reference

### ParseAndMerge

```go
func (ps *ProgramState) ParseAndMerge(filename string) error
```

Parses a file and merges results into the program state.

**Behavior:**
- Types are merged (compatible types only)
- Rules are appended (with ID uniqueness check)
- Facts are appended (with validation)
- Errors are accumulated

### ParseAndMergeContent

```go
func (ps *ProgramState) ParseAndMergeContent(content, filename string) error
```

Parses string content and merges results.

**Use case**: Dynamic content generation or testing.

### Facts Field

```go
type ProgramState struct {
    Facts []*Fact
    // ...
}
```

Direct access to accumulated facts.

### ToProgram

```go
func (ps *ProgramState) ToProgram() *Program
```

Converts accumulated state to `Program` structure.

## Troubleshooting

### Facts Not Accumulating

**Problem**: Facts count stays at 0 or doesn't increase.

**Possible causes**:
1. Type not defined → Check `ps.HasErrors()`
2. Invalid facts → Check `ps.GetErrors()`
3. Reset instruction → Remove `reset` from files

### Memory Issues

**Problem**: Out of memory with large datasets.

**Solutions**:
1. Process in smaller batches
2. Use reset periodically
3. Stream process facts instead of accumulating

### Duplicate Facts

**Problem**: Same fact appears multiple times.

**Note**: TSD does **not** deduplicate facts automatically. If you parse the same file twice, facts will be duplicated.

**Solution**: Track parsed files yourself or use fact IDs to deduplicate in your application logic.

## Related Documentation

- [Type Validation](./TYPE_VALIDATION.md)
- [Reset Instruction](../../../rete/docs/RESET.md)
- [Program State API](./PROGRAM_STATE_API.md)

## Summary

Incremental facts parsing in TSD provides:

✅ Flexible data loading strategies  
✅ Type-safe validation at each step  
✅ Non-blocking error handling  
✅ Support for multiple types and sources  
✅ Efficient memory and performance  
✅ Reset capability for fresh starts  

Use it to build scalable data processing pipelines while maintaining type safety and system stability.