# Types and Facts Without Rules - Behavior and Limitations

## Overview

This document clarifies the behavior when parsing TSD files that contain only type definitions and facts, without any rules.

## Current Behavior

### ✅ What Works

#### 1. Parsing Types and Facts
Files containing only types and facts parse successfully:

```tsd
type Person : <id: string, name: string, age: number>
type Product : <id: string, name: string, price: number>

Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)
Product(id: "PR001", name: "Laptop", price: 999.99)
```

**Result**: 
- ✅ Parsing succeeds
- ✅ Types are stored in `ProgramState.Types`
- ✅ Facts are validated against types
- ✅ Valid facts are stored in `ProgramState.Facts`
- ✅ Can be converted to `Program` structure

#### 2. Type Validation
Facts are validated against type definitions even without rules:

```tsd
type Person : <id: string, name: string, age: number>

Person(id: "P001", name: "Alice", age: 30)           # ✅ Valid
Person(id: "P002", name: "Bob", salary: 50000)       # ❌ Invalid field
Person(id: "P003", name: "Charlie", age: "thirty")   # ❌ Wrong type
```

**Result**:
- Valid facts are accepted
- Invalid facts are rejected with warnings
- Errors are recorded in `ProgramState.Errors`

#### 3. Incremental Parsing
You can parse types and facts incrementally across multiple files:

```go
ps := constraint.NewProgramState()

// Parse types
ps.ParseAndMerge("types.tsd")          // ✅ Works

// Parse facts (batch 1)
ps.ParseAndMerge("batch1_facts.tsd")   // ✅ Works

// Parse facts (batch 2)
ps.ParseAndMerge("batch2_facts.tsd")   // ✅ Works

// All facts accumulated
fmt.Printf("Total facts: %d\n", len(ps.Facts))
```

### ❌ What Doesn't Work

#### RETE Network Creation Without Rules

**The RETE network builder requires at least one rule** to create a network.

```go
// This will FAIL
pipeline := rete.NewConstraintPipeline()
storage := rete.NewMemoryStorage()
network, err := pipeline.BuildNetworkFromConstraintFile("types_and_facts_only.tsd", storage)
// err: "aucun nœud terminal dans le réseau" (no terminal node in network)
```

**Why?**

The RETE network validation explicitly checks for terminal nodes (rules):

```go
// From constraint_pipeline_validator.go
func (cp *ConstraintPipeline) validateNetwork(network *ReteNetwork) error {
    if len(network.TerminalNodes) == 0 {
        return fmt.Errorf("aucun nœud terminal dans le réseau")
    }
    return nil
}
```

This is intentional because:
- A RETE network without rules does nothing
- Facts would just be stored without any processing
- There are no actions to execute
- No pattern matching occurs

## Use Cases

### ✅ Valid Use Cases (Without RETE Network)

#### 1. Data Validation
Use TSD to validate data against schemas:

```go
ps := constraint.NewProgramState()

// Define schema
ps.ParseAndMergeContent(`
type Customer : <id: string, name: string, email: string>
`, "schema.tsd")

// Validate data
ps.ParseAndMergeContent(`
Customer(id: "C001", name: "Alice", email: "alice@example.com")
Customer(id: "C002", name: "Bob", email: "invalid")  // Will validate
`, "data.tsd")

// Check for validation errors
if ps.HasErrors() {
    for _, err := range ps.GetErrors() {
        fmt.Printf("Validation error: %s\n", err.Message)
    }
}
```

#### 2. Data Transformation
Parse and convert data to different formats:

```go
ps := constraint.NewProgramState()
ps.ParseAndMerge("data.tsd")

// Convert to Program structure
program := ps.ToProgram()

// Export to JSON
jsonData, _ := json.Marshal(program)
fmt.Println(string(jsonData))
```

#### 3. Schema Management
Manage type definitions separately from rules:

```go
ps := constraint.NewProgramState()

// Load schema
ps.ParseAndMerge("schemas/v1_types.tsd")

// Later, load rules
ps.ParseAndMerge("rules/business_logic.tsd")

// Now build RETE network (works because rules exist)
// ... convert to RETE format and build network
```

#### 4. Data Accumulation
Incrementally accumulate facts from multiple sources:

```go
ps := constraint.NewProgramState()
ps.ParseAndMerge("schema.tsd")

// Load data from different sources
for _, source := range dataSources {
    ps.ParseAndMerge(source)
}

// Process accumulated facts
for _, fact := range ps.Facts {
    // Custom processing logic
    processFact(fact)
}
```

### ❌ Invalid Use Case

#### Trying to Build RETE Network Without Rules

```go
// ❌ THIS WILL FAIL
ps := constraint.NewProgramState()
ps.ParseAndMerge("types_and_facts_only.tsd")

pipeline := rete.NewConstraintPipeline()
storage := rete.NewMemoryStorage()
network, err := pipeline.BuildNetworkFromConstraintFile("types_and_facts_only.tsd", storage)
// Error: "aucun nœud terminal dans le réseau"
```

## Workarounds

### Option 1: Add a Dummy Rule

If you need to build a RETE network but don't have real business logic yet:

```tsd
type Person : <id: string, name: string, age: number>

Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)

# Dummy rule to satisfy network creation
rule dummy : {p: Person} / p.id != "" ==> noop()
```

### Option 2: Use ProgramState Without RETE

Work directly with `ProgramState` for validation and data management:

```go
ps := constraint.NewProgramState()
ps.ParseAndMerge("types_and_facts.tsd")

// Use ProgramState directly
if ps.HasErrors() {
    handleErrors(ps.GetErrors())
}

// Access facts directly
for _, fact := range ps.Facts {
    processFact(fact)
}

// No RETE network needed
```

### Option 3: Separate Schema from Rules

Keep schemas separate and combine with rules when needed:

```go
// schemas/types.tsd
type Person : <id: string, name: string, age: number>

// data/facts.tsd
Person(id: "P001", name: "Alice", age: 30)

// rules/business.tsd
rule adult_check : {p: Person} / p.age >= 18 ==> adult(p.id)

// Build network with all files
files := []string{
    "schemas/types.tsd",
    "rules/business.tsd",
    "data/facts.tsd",
}
network, err := pipeline.BuildNetworkFromMultipleFiles(files, storage)
```

## API Behavior Summary

| Operation | Types Only | Types + Facts | Types + Facts + Rules |
|-----------|-----------|---------------|----------------------|
| `ParseAndMerge()` | ✅ Works | ✅ Works | ✅ Works |
| `ParseAndMergeContent()` | ✅ Works | ✅ Works | ✅ Works |
| Type validation | ✅ Works | ✅ Works | ✅ Works |
| Fact validation | N/A | ✅ Works | ✅ Works |
| `ToProgram()` | ✅ Works | ✅ Works | ✅ Works |
| `BuildNetworkFromConstraintFile()` | ❌ Fails* | ❌ Fails* | ✅ Works |
| `BuildNetworkFromMultipleFiles()` | ❌ Fails* | ❌ Fails* | ✅ Works |

\* Fails with error: "aucun nœud terminal dans le réseau"

## Recommendations

### For Schema-Only Files

If you only have types (schema definitions):

```tsd
# schemas/customer.tsd
type Customer : <id: string, name: string, email: string>
type Order : <id: string, customer_id: string, total: number>
```

**Use**: `ProgramState` for schema management and validation setup.

**Don't**: Try to build RETE network.

### For Data-Only Files

If you have types and facts but no rules:

```tsd
# data.tsd
type Person : <id: string, name: string>

Person(id: "P001", name: "Alice")
Person(id: "P002", name: "Bob")
```

**Use**: 
- `ProgramState` for data validation
- `ToProgram()` for data export
- Custom processing logic

**Don't**: Try to build RETE network.

### For Complete Applications

If you need RETE network capabilities:

```tsd
# app.tsd
type Person : <id: string, name: string, age: number>

Person(id: "P001", name: "Alice", age: 30)

rule adult_check : {p: Person} / p.age >= 18 ==> adult(p.id)
```

**Use**: Full RETE pipeline with network creation.

## Testing

Tests confirm current behavior:

```go
// ✅ This test passes
func TestParseTypesAndFactsWithoutRules(t *testing.T) {
    ps := constraint.NewProgramState()
    ps.ParseAndMergeContent("type Person : <id: string, name: string>", "types.tsd")
    ps.ParseAndMergeContent("Person(id: \"P001\", name: \"Alice\")", "facts.tsd")
    
    // Parsing works
    assert.Equal(t, 1, len(ps.Types))
    assert.Equal(t, 1, len(ps.Facts))
    assert.Equal(t, 0, len(ps.Rules))
}

// ❌ This test fails (expected behavior)
func TestBuildNetworkWithoutRules(t *testing.T) {
    pipeline := rete.NewConstraintPipeline()
    storage := rete.NewMemoryStorage()
    
    _, err := pipeline.BuildNetworkFromConstraintFile("types_and_facts_only.tsd", storage)
    
    // Expected error
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "aucun nœud terminal")
}
```

## Future Considerations

If there's demand for building RETE networks without rules, the validation could be made optional:

```go
// Hypothetical future API
pipeline := rete.NewConstraintPipeline()
pipeline.SetOption("allow_empty_network", true)  // Allow networks without rules
network, err := pipeline.BuildNetwork(...)        // Would succeed
```

However, this would require:
- Clear use cases for empty networks
- Decision on what happens when facts are submitted
- Storage and retrieval mechanisms without rule execution

## Conclusion

**Current Design Philosophy:**

TSD separates concerns:
- **Parsing & Validation**: Works with types and facts alone
- **RETE Network**: Requires rules to provide meaningful behavior

This is intentional and makes sense because:
- A RETE network without rules has no purpose
- Facts need rules to trigger actions
- Validation can happen at parse time without RETE

**When to use what:**
- **Types + Facts only** → Use `ProgramState` for validation
- **Types + Facts + Rules** → Build RETE network for execution

## Related Documentation

- [Type Validation](./TYPE_VALIDATION.md)
- [Incremental Facts Parsing](./INCREMENTAL_FACTS_PARSING.md)
- [RETE Network Architecture](../../../rete/docs/RETE_ARCHITECTURE.md)
- [Program State API](./PROGRAM_STATE_API.md)