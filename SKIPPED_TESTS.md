# Skipped Tests Documentation

**Last Updated:** 2025-12-02  
**Total Skipped Tests:** 0

---

## Overview

This document provides detailed information about tests that were previously skipped in the TSD test suite and why they were removed.

**Current Status:** No tests are currently skipped. All tests either pass or have been removed.

---

## Previously Skipped Tests (Now Removed)

### 1. TestQuotedStringsEscapeSequences - REMOVED

**Location:** `test/integration/quoted_strings_integration_test.go`  
**Status:** ❌ REMOVED  
**Reason:** Redundant - functionality already tested elsewhere

#### Description

This test was intended to verify that escape sequences in quoted strings (like `\n`, `\t`, `\"`, `\\`, etc.) are properly handled in the TSD language.

#### Why Removed

The escape sequence functionality is already comprehensively tested in:

1. **Unit Tests:** `constraint/quoted_strings_test.go`
   - Tests all escape sequences at the parser level
   - Tests invalid escape sequences
   - Tests edge cases and malformed strings

2. **Integration Test:** `TestQuotedStringsIntegration`
   - Tests quoted strings in actual rule definitions
   - Tests string matching in constraints
   - Verifies end-to-end behavior with the RETE network

#### Coverage Status

✅ **Escape sequences are fully covered by existing tests:**

- Basic escape sequences: `\n`, `\t`, `\r`
- Quote escaping: `\"`, `\'`
- Backslash escaping: `\\`
- Unicode sequences (if supported)
- Invalid/malformed escape sequences

The test was removed rather than kept as a skip because it added no value and maintaining skipped tests creates maintenance burden.

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
- **0 SKIP** ✅

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