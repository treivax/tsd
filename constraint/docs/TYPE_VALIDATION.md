# Type Validation in TSD

## Overview

TSD provides comprehensive type validation for both rules and facts to ensure that all references to types and fields are valid and conform to their definitions. This validation is **non-blocking**, meaning that invalid items produce warnings and are rejected, but do not prevent the processing of valid items.

## Features

### 1. Type Definition Validation

When a rule or fact references a type, the system validates that:
- The type has been previously defined
- The type definition is accessible in the current program state

### 2. Field Access Validation

When a rule or fact accesses a field of a type, the system validates that:
- The field exists in the type definition
- The field name matches exactly (case-sensitive)
- All field accesses in constraints and actions are valid

### 3. Field Type Validation

When a fact assigns a value to a field, the system validates that:
- The value type matches the field type definition
- Type conversions are appropriate (e.g., `identifier` can be treated as `string`)

## Validation Behavior

### Non-Blocking Errors

When validation fails:
1. A warning message is printed to stderr with the `⚠️` icon
2. The error is recorded in `ProgramState.Errors`
3. The invalid item is **rejected** (not added to the program state)
4. Processing continues with the next item

This approach ensures that:
- Valid items are always processed
- Invalid items don't pollute the system
- Debugging information is preserved
- The system remains stable

### Error Information

Each validation error contains:
- **File**: The source file where the error occurred
- **Type**: Either "rule" or "fact"
- **Message**: A descriptive error message
- **Line**: Line number (when available)

## Examples

### Example 1: Rule with Undefined Type

```tsd
type Person : <id: string, name: string, age: number>

# This rule will be REJECTED with a warning
rule r1 : {u: UnknownType} / u.name == "test" ==> log(u.id)
```

**Output:**
```
⚠️  Skipping invalid rule in example.tsd: variable u references undefined type UnknownType in example.tsd
```

**Result:** The rule is not added to the program state.

### Example 2: Rule with Invalid Field Access

```tsd
type Person : <id: string, name: string, age: number>

# This rule will be REJECTED with a warning (Person has no 'salary' field)
rule r2 : {p: Person} / p.salary > 1000 ==> high_earner(p.id)
```

**Output:**
```
⚠️  Skipping invalid rule in example.tsd: constraint validation failed: field p.salary not found in type Person
```

**Result:** The rule is not added to the program state.

### Example 3: Valid Rule

```tsd
type Person : <id: string, name: string, age: number>

# This rule is VALID and will be accepted
rule r3 : {p: Person} / p.age > 18 ==> adult(p.id, p.name)
```

**Output:** (no warning)

**Result:** The rule is added to the program state.

### Example 4: Fact with Undefined Type

```tsd
type Person : <id: string, name: string, age: number>

# This fact will be REJECTED with a warning
UnknownType(id: "U001", value: "test")
```

**Output:**
```
⚠️  Skipping invalid fact in example.tsd: fact references undefined type UnknownType in example.tsd
```

**Result:** The fact is not added to the program state.

### Example 5: Fact with Invalid Field

```tsd
type Person : <id: string, name: string, age: number>

# This fact will be REJECTED with a warning (Person has no 'salary' field)
Person(id: "P001", name: "Alice", salary: 50000)
```

**Output:**
```
⚠️  Skipping invalid fact in example.tsd: fact contains undefined field salary for type Person in example.tsd
```

**Result:** The fact is not added to the program state.

### Example 6: Fact with Wrong Field Type

```tsd
type Person : <id: string, name: string, age: number>

# This fact will be REJECTED with a warning (age should be number, not string)
Person(id: "P001", name: "Alice", age: "twenty-five")
```

**Output:**
```
⚠️  Skipping invalid fact in example.tsd: field age validation failed in example.tsd: expected number, got string
```

**Result:** The fact is not added to the program state.

### Example 7: Valid Fact

```tsd
type Person : <id: string, name: string, age: number>

# This fact is VALID and will be accepted
Person(id: "P001", name: "Alice", age: 25)
```

**Output:** (no warning)

**Result:** The fact is added to the program state.

## Validation in Actions

Field accesses in rule actions are also validated:

```tsd
type Person : <id: string, name: string, age: number>

# VALID: both p.id and p.name exist
rule r1 : {p: Person} / p.age > 18 ==> process(p.id, p.name)

# INVALID: p.address doesn't exist
rule r2 : {p: Person} / p.age > 18 ==> process(p.id, p.address)
```

**Output:**
```
⚠️  Skipping invalid rule in example.tsd: action validation failed: field p.address not found in type Person
```

## Mixed Valid/Invalid Items

The system processes valid items even when some items are invalid:

```tsd
type Person : <id: string, name: string, age: number>
type Product : <id: string, name: string, price: number>

# Valid facts
Person(id: "P001", name: "Alice", age: 25)
Product(id: "PR001", name: "Widget", price: 19.99)

# Invalid fact (unknown field)
Person(id: "P002", name: "Bob", invalidField: "test")

# Valid rule
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)

# Invalid rule (unknown type)
rule r2 : {x: UnknownType} / x.value > 0 ==> process(x.id)

# Valid rule
rule r3 : {pr: Product} / pr.price < 100 ==> affordable(pr.id)
```

**Result:**
- 2 types defined
- 2 facts added (P001, PR001)
- 2 rules added (r1, r3)
- 2 errors recorded (invalid fact P002, invalid rule r2)

## Checking for Validation Errors

### Programmatically

```go
ps := constraint.NewProgramState()
err := ps.ParseAndMergeContent(content, "example.tsd")

// Check if there are validation errors
if ps.HasErrors() {
    fmt.Printf("Found %d validation error(s)\n", ps.GetErrorCount())
    
    for _, validationErr := range ps.GetErrors() {
        fmt.Printf("  - %s in %s: %s\n", 
            validationErr.Type, 
            validationErr.File, 
            validationErr.Message)
    }
}
```

### Command Line

When using the TSD CLI, validation warnings are automatically printed to stderr:

```bash
$ tsd -file example.tsd
⚠️  Skipping invalid rule in example.tsd: variable u references undefined type UnknownType in example.tsd
```

## Best Practices

### 1. Define Types First

Always define types before using them in rules or facts:

```tsd
# ✓ GOOD: Type defined first
type Person : <id: string, name: string, age: number>
Person(id: "P001", name: "Alice", age: 25)

# ✗ BAD: Fact before type definition
Person(id: "P001", name: "Alice", age: 25)
type Person : <id: string, name: string, age: number>
```

### 2. Use Consistent Field Names

Field names are case-sensitive. Be consistent:

```tsd
type Person : <id: string, name: string, age: number>

# ✓ GOOD: Exact match
Person(id: "P001", name: "Alice", age: 25)

# ✗ BAD: Case mismatch
Person(ID: "P001", Name: "Alice", Age: 25)
```

### 3. Split Files Logically

Organize your code with types in one file and facts/rules in others:

```bash
# types.tsd - Type definitions
type Person : <id: string, name: string, age: number>
type Product : <id: string, name: string, price: number>

# facts.tsd - Data facts
Person(id: "P001", name: "Alice", age: 25)
Product(id: "PR001", name: "Widget", price: 19.99)

# rules.tsd - Business rules
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
rule r2 : {pr: Product} / pr.price < 100 ==> affordable(pr.id)
```

### 4. Monitor Validation Errors

Don't ignore validation warnings. They indicate potential bugs:

```bash
# Check for validation errors in CI/CD
if tsd -file *.tsd 2>&1 | grep -q "⚠️"; then
    echo "Validation errors found!"
    exit 1
fi
```

## Implementation Details

### Validation Flow

1. **Parse file** → Convert to AST
2. **Process resets** → Clear state if reset instruction present
3. **Merge types** → Add/validate type definitions
4. **Merge rules** → Validate and add rules
5. **Merge facts** → Validate and add facts

### Validation Methods

- `validateRule()` - Validates rule variables, constraints, and actions
- `validateFact()` - Validates fact type and field assignments
- `validateFieldAccesses()` - Recursively scans AST for field accesses
- `validateFactValue()` - Validates field value types

### Type System

Supported field types:
- `string` - Text values
- `number` - Numeric values (integer or float)
- `bool` - Boolean values (true/false)
- `identifier` - Special string type for IDs (treated as string)

## See Also

- [Rule Identifier Syntax](./RULE_IDENTIFIERS.md)
- [Type Definition Syntax](./GRAMMAR_COMPLETE.md)
- [Reset Instruction](./RESET.md)