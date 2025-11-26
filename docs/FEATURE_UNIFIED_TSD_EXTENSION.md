# Feature: Unified .tsd File Extension

## Overview
This feature unifies the previously separate `.constraint` and `.facts` file extensions into a single `.tsd` extension. A TSD file can contain any combination of type definitions, facts, and rules, providing a more flexible and unified approach to defining TSD programs.

## Motivation

### Current State
- **`.constraint` files**: Contain type definitions and rules
- **`.facts` files**: Contain fact assertions
- **Limitation**: Artificial separation requiring multiple files for a complete program

### Problems
1. **Fragmentation**: A single logical program is split across multiple files
2. **Confusion**: Users must understand which elements go in which file type
3. **Verbosity**: Test files and examples require paired files (.constraint + .facts)
4. **Inconsistency**: The language supports all elements in one file, but the tooling enforces separation

### Solution
- **Single extension**: `.tsd` for all TSD code files
- **Flexible content**: A .tsd file can contain types, facts, rules, or any combination
- **Backward compatibility**: Migration path for existing .constraint and .facts files

## Syntax

### File Extension
```
filename.tsd
```

### Content Examples

#### Types Only
```tsd
type Person : <id: string, name: string, age: int>
type Product : <id: string, price: float>
```

#### Facts Only
```tsd
Person("p1", "Alice", 30)
Person("p2", "Bob", 25)
Product("prod1", 99.99)
```

#### Rules Only
```tsd
rule check_adult : {p: Person} / p.age >= 18 ==> adult(p.id)
rule expensive : {prod: Product} / prod.price > 100 ==> expensive_item(prod.id)
```

#### Complete Program (Mixed)
```tsd
type Person : <id: string, name: string, age: int>

Person("p1", "Alice", 30)
Person("p2", "Bob", 17)

rule check_adult : {p: Person} / p.age >= 18 ==> adult(p.id)
rule check_minor : {p: Person} / p.age < 18 ==> minor(p.id)
```

## Implementation Details

### File Processing
1. **Parser**: Already supports mixed content (types, facts, rules)
2. **Extension change**: Rename all `.constraint` and `.facts` files to `.tsd`
3. **File merging**: When paired files exist (same basename, different extensions):
   - Read content from `.constraint` file
   - Append content from `.facts` file
   - Save combined content as `.tsd` file
   - Delete original paired files

### Migration Strategy

#### Automated Migration Script
- Script: `scripts/migrate_to_tsd.sh`
- Actions:
  1. Find all `.constraint` and `.facts` files
  2. For each unique basename:
     - If both `.constraint` and `.facts` exist: merge them
     - If only one exists: rename it
  3. Create `.tsd` file with combined/renamed content
  4. Remove original files

#### Manual Migration
```bash
# Single file rename
mv file.constraint file.tsd

# Paired files merge
cat file.constraint file.facts > file.tsd
```

### CLI Changes

#### Current Flags
```bash
tsd -constraint rules.constraint -facts data.facts
```

#### New Unified Flag
```bash
tsd -file program.tsd
# or simply
tsd program.tsd
```

#### Backward Compatibility (Optional)
The `-constraint` and `-facts` flags could remain as aliases to `-file` for a transition period, with deprecation warnings.

### Code Changes Required

1. **CLI flags** (`cmd/tsd/main.go`):
   - Add `-file` flag (or make positional argument)
   - Deprecate `-constraint` and `-facts` (with warnings)
   - Update help text and examples

2. **Documentation**:
   - Update all documentation to use `.tsd` extension
   - Update examples in README
   - Update CHANGELOG

3. **Test files**:
   - Rename all `.constraint` files to `.tsd`
   - Merge paired test files (same basename)
   - Update test code to reference `.tsd` files

4. **File search patterns**:
   - Update any glob patterns or file filters
   - Update `.gitignore` if applicable
   - Update editor configurations

## Migration Checklist

- [ ] Create migration script (`scripts/migrate_to_tsd.sh`)
- [ ] Migrate all test files in `beta_coverage_tests/`
- [ ] Migrate all test files in `test/coverage/`
- [ ] Migrate all test files in `constraint/test/`
- [ ] Migrate all test files in `rete/test/`
- [ ] Update CLI argument parsing
- [ ] Update help text and examples
- [ ] Update test code references
- [ ] Update documentation (README, examples)
- [ ] Update CHANGELOG
- [ ] Run full test suite to verify migration

## Benefits

1. **Simplicity**: One file extension to remember
2. **Flexibility**: Mix types, facts, and rules as needed
3. **Cohesion**: Related code stays together
4. **Clarity**: File extension matches project name (TSD)
5. **Reduced boilerplate**: Fewer files to manage

## Examples

### Before (Paired Files)

**person.constraint:**
```
type Person : <id: string, name: string, age: int>
rule adult : {p: Person} / p.age >= 18 ==> is_adult(p.id)
```

**person.facts:**
```
Person("p1", "Alice", 30)
Person("p2", "Bob", 17)
```

### After (Single File)

**person.tsd:**
```
type Person : <id: string, name: string, age: int>

Person("p1", "Alice", 30)
Person("p2", "Bob", 17)

rule adult : {p: Person} / p.age >= 18 ==> is_adult(p.id)
```

## Testing Strategy

1. **Unit tests**: Verify file parsing with `.tsd` extension
2. **Integration tests**: Run full pipeline with merged `.tsd` files
3. **Regression tests**: Ensure all existing tests pass after migration
4. **Edge cases**:
   - Files with no facts (types + rules only)
   - Files with no rules (types + facts only)
   - Empty files
   - Unicode content
   - Large merged files

## Documentation Updates

- [ ] `README.md`: Update examples and usage
- [ ] `docs/rule_identifiers.md`: Update file extension references
- [ ] `docs/rule_id_uniqueness.md`: Update examples
- [ ] CLI help text
- [ ] Error messages referencing file types
- [ ] Any tutorial or guide documents

## Version Impact

This is a **breaking change** that should be released as a new major version (e.g., v3.0.0), as it changes:
- File extensions (user-facing)
- CLI flags (potentially, if deprecated)
- File structure expectations

However, the language syntax and semantics remain unchanged—only the packaging changes.

## Timeline

1. **Phase 1**: Create migration script and documentation (this document)
2. **Phase 2**: Migrate all test files
3. **Phase 3**: Update CLI and code references
4. **Phase 4**: Run full test suite and fix any issues
5. **Phase 5**: Update documentation
6. **Phase 6**: Commit and release

## Success Criteria

- [ ] All test files migrated to `.tsd` extension
- [ ] All paired files successfully merged
- [ ] All tests passing
- [ ] No references to `.constraint` or `.facts` in code (except deprecation warnings)
- [ ] Documentation fully updated
- [ ] CLI accepts `.tsd` files
- [ ] Users can create a complete program in a single `.tsd` file