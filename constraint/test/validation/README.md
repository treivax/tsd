# Type Validation Test Files

This directory contains test files for validating the type validation system in TSD.

## Files

### `valid_items.tsd`
Contains only valid types, facts, and rules. Used to verify that the parser and validator correctly accept well-formed TSD code.

**Contents:**
- 3 type definitions (Person, Product, Order)
- 8 valid facts
- 5 valid rules with proper field accesses

**Expected behavior:** All items should be parsed and added to the program state without errors.

### `invalid_items.tsd`
Contains a mix of valid and invalid types, facts, and rules. Used to verify that the validator correctly detects and rejects malformed items.

**Contents:**
- 2 valid type definitions (Person, Product)
- 2 valid facts
- 4 invalid facts (undefined type, undefined fields, wrong field types)
- 2 valid rules
- 5 invalid rules (undefined types, invalid field accesses in constraints and actions)

**Expected behavior:**
- Valid items are parsed and added
- Invalid items trigger non-blocking errors and are rejected
- Each error is recorded with file, type, and descriptive message
- Processing continues for all items

## Test Coverage

These files test the following validation scenarios:

### Facts Validation
1. ✅ Valid facts with all required fields
2. ❌ Facts referencing undefined types
3. ❌ Facts with undefined fields
4. ❌ Facts with wrong field types (e.g., string instead of number)

### Rules Validation
1. ✅ Valid rules with proper variable types and field accesses
2. ❌ Rules referencing undefined types in variables
3. ❌ Rules with invalid field accesses in constraints
4. ❌ Rules with invalid field accesses in actions
5. ❌ Rules with multiple variables and complex constraints

## Usage

These files are used by `type_validation_files_test.go` which tests:

```go
// Valid file - should have no errors
ps.ParseAndMerge("test/validation/valid_items.tsd")
assert(ps.HasErrors() == false)

// Invalid file - should have errors but not crash
ps.ParseAndMerge("test/validation/invalid_items.tsd")
assert(ps.HasErrors() == true)
assert(len(ps.Facts) == 2) // Only valid facts added
assert(len(ps.Rules) == 2) // Only valid rules added
```

## Expected Error Messages

When parsing `invalid_items.tsd`, the following error patterns should appear:

- `"undefined type"` - For facts/rules referencing non-existent types
- `"undefined field"` - For facts with fields not in type definition
- `"not found in type"` - For rules accessing non-existent fields
- `"expected number"` - For type mismatches in fact values

## Adding New Test Cases

To add new validation test cases:

1. Add the type definition if needed
2. Add the invalid item (fact or rule)
3. Add a comment explaining what makes it invalid
4. Update the test expectations in `type_validation_files_test.go`

Example:
```tsd
# INVALID FACT: Missing required field 'age'
Person(id: "P999", name: "Test")
```

## Non-Blocking Error Philosophy

All validation errors are **non-blocking**:
- Invalid items are rejected with a warning (⚠️)
- Errors are recorded in `ProgramState.Errors`
- Valid items continue to be processed
- The system remains in a consistent state

This allows incremental development and debugging without breaking the entire program.