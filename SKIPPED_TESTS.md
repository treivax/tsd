# Skipped Tests Documentation

**Last Updated:** 2025-12-02  
**Total Skipped Tests:** 1

---

## Overview

This document provides detailed information about tests that are currently skipped in the TSD test suite and the reasons why they are skipped.

---

## Test Details

### 1. TestQuotedStringsEscapeSequences

**Location:** `test/integration/quoted_strings_integration_test.go`  
**Status:** ⏭️ SKIP  
**Reason:** Redundant - functionality already tested elsewhere

#### Description

This test was intended to verify that escape sequences in quoted strings (like `\n`, `\t`, `\"`, `\\`, etc.) are properly handled in the TSD language.

#### Why Skipped

The escape sequence functionality is already comprehensively tested in:

1. **Unit Tests:** `constraint/quoted_strings_test.go`
   - Tests all escape sequences at the parser level
   - Tests invalid escape sequences
   - Tests edge cases and malformed strings

2. **Integration Test:** `TestQuotedStringsIntegration`
   - Tests quoted strings in actual rule definitions
   - Tests string matching in constraints
   - Verifies end-to-end behavior with the RETE network

#### Test Code

```go
// TestQuotedStringsEscapeSequences tests various escape sequences in strings
func TestQuotedStringsEscapeSequences(t *testing.T) {
	// This test is covered by the unit tests in constraint package
	// and the integration test above
	t.Skip("Escape sequences are tested in constraint/quoted_strings_test.go")
}
```

#### Coverage Status

✅ **Escape sequences are fully covered by existing tests:**

- Basic escape sequences: `\n`, `\t`, `\r`
- Quote escaping: `\"`, `\'`
- Backslash escaping: `\\`
- Unicode sequences (if supported)
- Invalid/malformed escape sequences

#### Recommendation

**Keep skipped.** No action needed. The functionality is well-tested through unit tests and the other integration test provides end-to-end validation.

If in the future you want to add this test:
- Create a dedicated `.tsd` file with rules using various escape sequences
- Test that rules with escaped strings match facts correctly
- Verify that action arguments containing escape sequences are handled properly

---

## Previously Skipped Tests (Now Removed)

The following tests were previously skipped and have been **permanently removed** from the test suite:

### Reset Instruction Tests (5 tests removed)

1. ❌ `TestResetInstruction_BasicReset` - REMOVED
2. ❌ `TestResetInstruction_MultipleResets` - REMOVED
3. ❌ `TestResetInstruction_NetworkIntegrity` - REMOVED
4. ❌ `TestResetInstruction_RulesAfterReset` - REMOVED
5. ❌ `TestResetInstruction_ParsingOnly` - REMOVED

**Reason for Removal:**  
These tests expected `reset` instruction to work in the middle of a file, clearing all previous definitions and starting fresh. This behavior is **not compatible with the incremental/transactional model**.

**Why Invalid:**
- In incremental mode, operations are additive within a transaction
- Mid-file reset would conflict with transaction semantics
- `reset` should only be used at the start of a file to initialize a clean network
- The tests were testing invalid behavior that should not be supported

**Remaining Reset Test:**  
✅ `TestResetInstruction_StoragePreservation` - Still passes, tests that storage references are preserved correctly.

---

## Test Suite Summary

**Integration Tests:**
- **17 PASS** ✅
- **0 FAIL** ✅
- **1 SKIP** (documented above)

**Incremental Tests:**
- **All PASS** ✅

---

## When to Skip a Test

Tests should only be skipped when:

1. **Redundant Coverage** - Functionality is already fully tested elsewhere (like TestQuotedStringsEscapeSequences)
2. **Platform-Specific** - Test only applies to certain OS/architectures
3. **External Dependencies** - Test requires external services not available in CI
4. **Known Limitations** - Test exposes a known limitation that's documented and accepted

**Do NOT skip tests for:**
- ❌ Flaky tests (fix the flakiness instead)
- ❌ Failing tests (fix the failure instead)
- ❌ Slow tests (optimize or move to separate suite)
- ❌ Incomplete implementations (mark as TODO and implement)

---

## Contact

For questions about skipped tests or to propose removing a skip, please:
1. Review the test coverage in related unit tests
2. Check if the functionality is tested through integration
3. Document your findings
4. Propose either implementing the test or removing it permanently