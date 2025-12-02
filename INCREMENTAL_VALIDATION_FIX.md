# Incremental Validation Fix - Action Definitions Support

## Summary

Fixed a critical issue in incremental validation where action definitions were not being tracked or merged when validating TSD files incrementally. This caused all rules with actions to fail validation when loaded in separate files from their action definitions.

## Problem

When ingesting TSD files incrementally:
1. First file: types and facts → ✅ Works
2. Second file: actions and rules → ❌ Failed with "action manquante dans l'expression"

The validator was merging type definitions from previous files but not action definitions, causing semantic validation to fail even when actions were properly defined.

## Root Causes

### 1. Network Not Tracking Actions
The `ReteNetwork` struct had no field to store action definitions for incremental validation context.

### 2. Actions Not Converted to RETE Format
The `ConvertToReteProgram` function in `constraint/api.go` was only converting types and expressions, but not actions.

### 3. Actions Not Extracted from AST
No code was extracting and storing action definitions when processing ingested files.

### 4. Actions Not Merged During Validation
The `IncrementalValidator` was only merging types, not actions, when creating the context for validation.

### 5. Expressions Not Fully Converted
When converting the merged program back to AST for validation, the `programToAST` method was only including `type` and `set` fields, omitting the crucial `action` field.

## Changes Made

### 1. Added Action Storage to ReteNetwork (`rete/network.go`)
```go
type ReteNetwork struct {
    // ... existing fields ...
    Actions []ActionDefinition `json:"actions"` // Action definitions for incremental validation
    // ... rest ...
}
```

### 2. Added ActionDefinition Types (`rete/structures.go`)
```go
type Parameter struct {
    Name string `json:"name"`
    Type string `json:"type"`
}

type ActionDefinition struct {
    Type       string      `json:"type"`
    Name       string      `json:"name"`
    Parameters []Parameter `json:"parameters"`
}
```

### 3. Updated ConvertToReteProgram (`constraint/api.go`)
Now includes actions in the conversion:
```go
reteProgram := map[string]interface{}{
    "types":       typesInterface,
    "actions":     actionsInterface,  // NEW
    "expressions": expressionsInterface,
}
```

### 4. Added Action Extraction (`rete/constraint_pipeline_parser.go`)
New function `extractAndStoreActions` that:
- Extracts action definitions from parsed AST
- Stores them in `network.Actions`
- Handles action redefinitions

### 5. Updated Incremental Validator (`rete/incremental_validation.go`)
- Added `extractExistingActions()` to extract actions from network context
- Updated `mergePrograms()` to merge both types and actions
- Fixed `programToAST()` to include all expression fields (action, constraints, patterns, ruleId)

### 6. Integrated into Pipeline (`rete/constraint_pipeline.go`)
Added call to `extractAndStoreActions` after type processing.

## Test Results

### Before Fix
```
--- FAIL: TestIncrementalIngestion_FactsBeforeRules
--- FAIL: TestIncrementalIngestion_MultipleRules
--- FAIL: TestIncrementalIngestion_TypeExtension
--- FAIL: TestIncrementalIngestion_Optimizations
```

### After Fix
```
--- PASS: TestIncrementalIngestion_FactsBeforeRules (0.00s)
--- PASS: TestIncrementalIngestion_MultipleRules (0.00s)
--- PASS: TestIncrementalIngestion_TypeExtension (0.00s)
--- PASS: TestIncrementalIngestion_Optimizations (0.00s)
--- PASS: TestIncrementalIngestion_Reset (0.00s)
```

### Overall Integration Test Status
- **19 tests PASS** (up from ~15 before)
- **16 tests FAIL** (down from ~19 before)
- **1 test SKIP**

## Remaining Issues (Not Related to This Fix)

### 1. Advanced Tests Using Old Syntax
Tests in `test/integration/incremental/advanced_test.go` use deprecated syntax:
- Old: `type Person { id: string }`
- New: `type Person(id: string)`
These need to be rewritten.

### 2. Reset Tests
Some reset instruction tests fail because types are not being fully cleaned up after reset command.

### 3. Alpha Coverage Tests
Rules are not matching facts due to alpha node evaluation issues (unrelated to action validation).

### 4. Quoted Strings Test
No rule activations occurring (evaluation issue, not validation issue).

## Files Modified

1. `rete/network.go` - Added Actions field
2. `rete/structures.go` - Added ActionDefinition and Parameter types
3. `rete/incremental_validation.go` - Added action merging logic
4. `rete/constraint_pipeline.go` - Integrated action extraction
5. `rete/constraint_pipeline_parser.go` - Added extractAndStoreActions function
6. `constraint/api.go` - Updated ConvertToReteProgram to include actions

## Impact

- ✅ Incremental ingestion now works correctly with actions defined separately from rules
- ✅ Action definitions are properly tracked across file ingestions
- ✅ Validation context includes both types and actions from previous files
- ✅ No breaking changes to existing APIs
- ✅ All existing rete and constraint tests still pass

## Example Usage

```tsd
# File 1: types_and_facts.tsd
type Person(id: string, name: string, age: number)

Person(id: P001, name: Alice, age: 25)
Person(id: P002, name: Bob, age: 30)
```

```tsd
# File 2: actions_and_rules.tsd
action adult(id: string)

rule r1: {p: Person} / p.age >= 18 ==> adult(p.id)
```

```go
// Incremental ingestion now works!
network, _ := pipeline.IngestFile("types_and_facts.tsd", nil, storage)
network, _ = pipeline.IngestFile("actions_and_rules.tsd", network, storage)
```

## Date
December 2, 2025