# Phase 4 - Coherence Modes Design Document

**Feature**: Selectable Coherence Modes (Strong / Relaxed / Eventual)  
**Status**: 🚧 Design Phase  
**Target Completion**: 2025-12-04  
**Estimated Effort**: 8-12 hours

---

## Executive Summary

This document outlines the design for implementing selectable coherence modes in the TSD RETE engine, allowing users to choose between different consistency guarantees based on their application requirements.

### Coherence Modes Overview

1. **Strong Consistency** (Default)
   - Strictest guarantees
   - All reads reflect the most recent write
   - Maximum correctness, potential performance cost

2. **Relaxed Consistency**
   - Balanced approach
   - Eventually consistent with bounded staleness
   - Better performance, acceptable for most use cases

3. **Eventual Consistency**
   - Weakest guarantees, highest performance
   - Reads may see stale data temporarily
   - Suitable for high-throughput scenarios where temporary inconsistency is acceptable

---

## Motivation

### Current State (Phase 3)
- ✅ Thread-safe transaction model implemented
- ✅ Read-after-write guarantees enforced
- ✅ Synchronous fact verification with retries
- ⚠️ **All operations use "strong" consistency by default**
- ⚠️ **No flexibility for performance-critical scenarios**

### Use Cases Requiring Different Modes

#### Strong Consistency
- Financial transactions
- Critical business rules
- Compliance and audit systems
- Any scenario where data correctness is paramount

#### Relaxed Consistency
- Real-time analytics dashboards
- Recommendation engines
- User activity tracking
- Moderate-throughput event processing

#### Eventual Consistency
- Logging and telemetry
- Social media feeds
- High-volume sensor data
- Scenarios where occasional stale reads are acceptable

---

## Design Goals

1. **Backward Compatibility**: Existing code continues to work with strong consistency
2. **Explicit Configuration**: Mode selection is clear and intentional
3. **Transaction-Level Control**: Different transactions can use different modes
4. **Observable Behavior**: Mode choice affects observable guarantees
5. **Performance Optimization**: Each mode trades consistency for performance appropriately
6. **Safe Defaults**: Strong consistency remains the default

---

## API Design

### 1. CoherenceMode Type

```go
// CoherenceMode defines the consistency guarantee level for transactions
type CoherenceMode int

const (
    // CoherenceModeStrong provides the strictest consistency guarantees.
    // All reads reflect the most recent writes. Maximum correctness.
    // Default mode for backward compatibility.
    CoherenceModeStrong CoherenceMode = iota
    
    // CoherenceModeRelaxed provides balanced consistency.
    // Reads may lag behind writes by a bounded time (configurable).
    // Suitable for most real-time applications.
    CoherenceModeRelaxed
    
    // CoherenceModeEventual provides eventual consistency.
    // Reads may return stale data temporarily.
    // Highest performance, suitable for high-throughput scenarios.
    CoherenceModeEventual
)

// String returns the human-readable name of the coherence mode
func (m CoherenceMode) String() string {
    switch m {
    case CoherenceModeStrong:
        return "Strong"
    case CoherenceModeRelaxed:
        return "Relaxed"
    case CoherenceModeEventual:
        return "Eventual"
    default:
        return "Unknown"
    }
}
```

### 2. Transaction API Extension

```go
// Transaction options (new)
type TransactionOptions struct {
    CoherenceMode      CoherenceMode
    MaxStaleness       time.Duration // For Relaxed mode
    SkipVerification   bool          // For Eventual mode
    VerifyOnCommit     bool          // Verify all facts on commit (Strong/Relaxed)
}

// BeginTransactionWithMode starts a new transaction with specified coherence mode
func (n *ReteNetwork) BeginTransactionWithMode(mode CoherenceMode, opts *TransactionOptions) *Transaction {
    if opts == nil {
        opts = DefaultOptionsForMode(mode)
    }
    
    tx := &Transaction{
        ID:            generateTransactionID(),
        Network:       n,
        Commands:      make([]*Command, 0),
        Mode:          mode,
        Options:       opts,
        Logger:        n.Logger,
    }
    
    return tx
}

// DefaultOptionsForMode returns sensible defaults for each mode
func DefaultOptionsForMode(mode CoherenceMode) *TransactionOptions {
    switch mode {
    case CoherenceModeStrong:
        return &TransactionOptions{
            CoherenceMode:    CoherenceModeStrong,
            MaxStaleness:     0, // No staleness allowed
            SkipVerification: false,
            VerifyOnCommit:   true,
        }
    case CoherenceModeRelaxed:
        return &TransactionOptions{
            CoherenceMode:    CoherenceModeRelaxed,
            MaxStaleness:     100 * time.Millisecond, // 100ms staleness acceptable
            SkipVerification: false,
            VerifyOnCommit:   true,
        }
    case CoherenceModeEventual:
        return &TransactionOptions{
            CoherenceMode:    CoherenceModeEventual,
            MaxStaleness:     time.Hour, // Effectively unlimited
            SkipVerification: true,      // Skip verification entirely
            VerifyOnCommit:   false,     // Trust that it will eventually be consistent
        }
    default:
        return DefaultOptionsForMode(CoherenceModeStrong)
    }
}
```

### 3. Backward-Compatible API

```go
// BeginTransaction continues to use Strong mode (backward compatible)
func (n *ReteNetwork) BeginTransaction() *Transaction {
    return n.BeginTransactionWithMode(CoherenceModeStrong, nil)
}

// BeginRelaxedTransaction is a convenience method
func (n *ReteNetwork) BeginRelaxedTransaction() *Transaction {
    return n.BeginTransactionWithMode(CoherenceModeRelaxed, nil)
}

// BeginEventualTransaction is a convenience method
func (n *ReteNetwork) BeginEventualTransaction() *Transaction {
    return n.BeginTransactionWithMode(CoherenceModeEventual, nil)
}
```

### 4. Transaction Structure Update

```go
type Transaction struct {
    ID            string
    Network       *ReteNetwork
    Commands      []*Command
    Mode          CoherenceMode           // NEW
    Options       *TransactionOptions     // NEW
    Logger        *Logger
    mu            sync.Mutex
}
```

---

## Implementation Details

### Strong Consistency Behavior

**Guarantees**:
- Every `AddFact` is verified in storage before returning
- Retry mechanism with exponential backoff
- Transactions are atomic: all facts or none
- Read-after-write consistency enforced

**Implementation**:
- Use existing `waitForFactPersistence()` logic
- `VerifyRetryDelay`: 50ms (configurable)
- `MaxVerifyRetries`: 10 (configurable)
- `SubmissionTimeout`: 30s per fact batch

**Performance**:
- Slowest mode
- Best correctness guarantees
- Suitable for: <1000 facts/sec

---

### Relaxed Consistency Behavior

**Guarantees**:
- Facts are verified with a bounded staleness window
- Reads may lag behind writes by up to `MaxStaleness` (default: 100ms)
- Eventual verification: facts not found immediately may be retried less aggressively
- Transactions still atomic on commit

**Implementation**:
```go
func (tx *Transaction) addFactRelaxed(fact *Fact) error {
    // 1. Add to storage immediately
    if err := tx.Network.Storage.AddFact(fact); err != nil {
        return err
    }
    
    // 2. Quick verification (1-2 retries max)
    verified := tx.quickVerify(fact, 2, 20*time.Millisecond)
    
    // 3. If not verified immediately, schedule async verification
    if !verified {
        tx.scheduleAsyncVerification(fact)
    }
    
    // 4. Return immediately (don't block on full verification)
    return nil
}

func (tx *Transaction) quickVerify(fact *Fact, maxRetries int, delay time.Duration) bool {
    for i := 0; i < maxRetries; i++ {
        if retrieved := tx.Network.Storage.GetFact(fact.InternalID); retrieved != nil {
            return true
        }
        if i < maxRetries-1 {
            time.Sleep(delay)
        }
    }
    return false
}
```

**Performance**:
- Medium speed (2-5x faster than Strong)
- Acceptable for most use cases
- Suitable for: 1000-10,000 facts/sec

---

### Eventual Consistency Behavior

**Guarantees**:
- Facts are added to storage with no verification
- No retry mechanism
- Reads may return stale data or miss recent writes temporarily
- Best-effort consistency: trust that storage will eventually reflect writes

**Implementation**:
```go
func (tx *Transaction) addFactEventual(fact *Fact) error {
    // 1. Add to storage immediately
    if err := tx.Network.Storage.AddFact(fact); err != nil {
        return err
    }
    
    // 2. No verification - trust it's there
    tx.Logger.Debug("Fact added (eventual mode, no verification): %s", fact.ID)
    
    // 3. Return immediately
    return nil
}
```

**Performance**:
- Fastest mode (10-100x faster than Strong)
- No blocking on verification
- Suitable for: >10,000 facts/sec

---

## Transaction Lifecycle

### Strong Mode Transaction Flow

```
BeginTransaction() [Strong]
    ↓
RecordAndExecute(AddFactCommand) → Verify fact immediately (blocking)
    ↓ (retry up to 10 times with backoff)
Fact verified ✓
    ↓
RecordAndExecute(AddFactCommand) → Verify next fact
    ↓
...
    ↓
Commit() → Verify all facts one final time
    ↓
Return success
```

### Relaxed Mode Transaction Flow

```
BeginTransactionWithMode(Relaxed)
    ↓
RecordAndExecute(AddFactCommand) → Quick verify (1-2 retries, non-blocking)
    ↓
Fact likely verified ✓ (or scheduled for async verification)
    ↓
RecordAndExecute(AddFactCommand) → Quick verify
    ↓
...
    ↓
Commit() → Verify all facts (with staleness tolerance)
    ↓
Return success (even if some facts still propagating)
```

### Eventual Mode Transaction Flow

```
BeginTransactionWithMode(Eventual)
    ↓
RecordAndExecute(AddFactCommand) → Add to storage, no verification
    ↓
RecordAndExecute(AddFactCommand) → Add to storage, no verification
    ↓
...
    ↓
Commit() → No verification, immediate return
    ↓
Return success (facts may still be propagating)
```

---

## Configuration

### Global Defaults

```go
// Default coherence mode for the network
type NetworkConfig struct {
    DefaultCoherenceMode   CoherenceMode
    StrongModeConfig       StrongModeConfig
    RelaxedModeConfig      RelaxedModeConfig
    EventualModeConfig     EventualModeConfig
}

type StrongModeConfig struct {
    SubmissionTimeout  time.Duration // 30s
    VerifyRetryDelay   time.Duration // 50ms
    MaxVerifyRetries   int           // 10
}

type RelaxedModeConfig struct {
    MaxStaleness       time.Duration // 100ms
    QuickVerifyRetries int           // 2
    QuickVerifyDelay   time.Duration // 20ms
}

type EventualModeConfig struct {
    SkipVerification   bool // true
    AsyncVerification  bool // false (for future: background verification)
}
```

### Per-Transaction Override

```go
// User can override per transaction
tx := network.BeginTransactionWithMode(CoherenceModeRelaxed, &TransactionOptions{
    MaxStaleness: 200 * time.Millisecond, // Custom staleness bound
})
```

---

## Metrics and Observability

### New Metrics

```go
// Coherence mode metrics
type CoherenceMetrics struct {
    // Per-mode counters
    StrongTransactions   int64
    RelaxedTransactions  int64
    EventualTransactions int64
    
    // Verification metrics
    VerificationAttempts int64
    VerificationFailures int64
    QuickVerificationHits int64
    QuickVerificationMisses int64
    
    // Staleness tracking (Relaxed mode)
    StalenessObserved    []time.Duration // Histogram
    MaxStalenessObserved time.Duration
    
    // Performance
    TransactionDuration  map[CoherenceMode]time.Duration
}
```

### Logging

```go
// Transaction start
logger.Info("Starting transaction [mode=%s, opts=%+v]", tx.Mode, tx.Options)

// Fact verification
logger.Debug("Verifying fact [mode=%s, id=%s, attempt=%d]", tx.Mode, fact.ID, attempt)

// Transaction commit
logger.Info("Committing transaction [mode=%s, facts=%d, duration=%s]", 
    tx.Mode, len(tx.Commands), duration)
```

---

## Error Handling

### Strong Mode Errors

- **Verification timeout**: Return error, rollback transaction
- **Storage error**: Return error, rollback transaction
- **Max retries exceeded**: Return error, rollback transaction

### Relaxed Mode Errors

- **Quick verification failed**: Log warning, schedule async verification, continue
- **Storage error on commit**: Return error, rollback transaction
- **Staleness exceeded**: Log warning, continue (configurable to error)

### Eventual Mode Errors

- **Storage error**: Return error immediately (no retry)
- **Network error**: Log error, continue (best-effort)

---

## Testing Strategy

### Unit Tests

1. **Mode selection**
   - Test default mode (Strong)
   - Test explicit mode selection
   - Test mode-specific options

2. **Verification behavior**
   - Strong: All facts verified synchronously
   - Relaxed: Quick verification + optional async
   - Eventual: No verification

3. **Transaction lifecycle**
   - Begin, execute commands, commit for each mode
   - Rollback behavior for each mode

4. **Configuration**
   - Default configs applied correctly
   - Custom configs override defaults
   - Invalid configs rejected

### Integration Tests

1. **Concurrent transactions with different modes**
2. **Mode switching between transactions**
3. **Performance comparison across modes**
4. **Staleness bounds in Relaxed mode**
5. **Eventual consistency convergence**

### Performance Tests

1. **Throughput comparison**: Strong vs Relaxed vs Eventual
2. **Latency percentiles** (p50, p95, p99) per mode
3. **Scalability**: How does each mode scale with fact volume?

---

## Migration Path

### Phase 1: API Introduction (Week 1)
- Add `CoherenceMode` enum
- Add `TransactionOptions` struct
- Implement `BeginTransactionWithMode()`
- **All transactions still use Strong by default**

### Phase 2: Strong Mode Formalization (Week 1)
- Extract existing logic into `addFactStrong()`
- Add comprehensive tests for Strong mode
- Document Strong mode guarantees

### Phase 3: Relaxed Mode Implementation (Week 2)
- Implement `addFactRelaxed()` with quick verification
- Add staleness tracking
- Add tests for Relaxed mode

### Phase 4: Eventual Mode Implementation (Week 2)
- Implement `addFactEventual()` with no verification
- Add performance tests
- Add tests for Eventual mode

### Phase 5: Documentation and Examples (Week 3)
- Update user documentation
- Add examples for each mode
- Add migration guide

---

## Examples

### Example 1: Strong Consistency (Default)

```go
// Financial transaction - requires strict consistency
tx := network.BeginTransaction() // Strong mode (default)

// Add account debit
tx.RecordAndExecute(AddFactCommand{
    Fact: &Fact{ID: "debit_123", Type: "Transaction", ...},
})

// Add account credit
tx.RecordAndExecute(AddFactCommand{
    Fact: &Fact{ID: "credit_456", Type: "Transaction", ...},
})

// Commit - both facts verified before return
if err := tx.Commit(); err != nil {
    log.Fatalf("Transaction failed: %v", err)
}
// ✅ Both facts are guaranteed to be in storage and readable
```

### Example 2: Relaxed Consistency

```go
// Real-time analytics dashboard - 100ms staleness acceptable
tx := network.BeginRelaxedTransaction()

for _, event := range userEvents {
    tx.RecordAndExecute(AddFactCommand{
        Fact: &Fact{ID: event.ID, Type: "UserEvent", ...},
    })
}

if err := tx.Commit(); err != nil {
    log.Fatalf("Failed to log events: %v", err)
}
// ✅ Facts will be readable within ~100ms
```

### Example 3: Eventual Consistency

```go
// High-volume sensor data - eventual consistency sufficient
tx := network.BeginEventualTransaction()

for _, reading := range sensorReadings {
    tx.RecordAndExecute(AddFactCommand{
        Fact: &Fact{ID: reading.ID, Type: "SensorReading", ...},
    })
}

if err := tx.Commit(); err != nil {
    log.Printf("Warning: Some readings may have failed: %v", err)
}
// ✅ Maximum throughput, readings will eventually be consistent
```

### Example 4: Custom Staleness Bound

```go
// Custom relaxed mode with 500ms staleness tolerance
opts := &TransactionOptions{
    CoherenceMode: CoherenceModeRelaxed,
    MaxStaleness:  500 * time.Millisecond,
}

tx := network.BeginTransactionWithMode(CoherenceModeRelaxed, opts)
// ... use transaction
```

---

## Performance Targets

| Mode | Throughput Target | Latency (p95) | Verification Overhead |
|------|-------------------|---------------|----------------------|
| **Strong** | 100-1,000 facts/sec | 50-200ms | ~30-40% |
| **Relaxed** | 1,000-10,000 facts/sec | 5-20ms | ~5-10% |
| **Eventual** | >10,000 facts/sec | <5ms | 0% |

---

## Risks and Mitigations

### Risk 1: User Misunderstanding
**Risk**: Users choose Eventual mode without understanding implications  
**Mitigation**: 
- Clear documentation with warnings
- Strong mode as default
- Require explicit opt-in for Eventual

### Risk 2: Eventual Mode Data Loss
**Risk**: Storage failure in Eventual mode may lose data  
**Mitigation**:
- Log all storage errors even in Eventual mode
- Add optional async verification for critical eventual data
- Provide metrics on verification failures

### Risk 3: Performance Regression
**Risk**: New abstraction slows down existing Strong mode  
**Mitigation**:
- Benchmark Strong mode before and after
- Keep existing code path for Strong mode initially
- Optimize after correctness is verified

---

## Open Questions

1. **Async verification for Relaxed mode**: Should we implement background verification workers?
   - **Decision**: Phase 1 - No. Keep it simple. Evaluate in Phase 2 if needed.

2. **Per-fact mode override**: Should individual facts in a transaction have different modes?
   - **Decision**: No. Transaction-level mode only for simplicity.

3. **Dynamic mode switching**: Should we allow mode change mid-transaction?
   - **Decision**: No. Mode is set at transaction begin, immutable.

4. **Metrics integration**: How to integrate with Prometheus exporter?
   - **Decision**: Add coherence mode labels to existing transaction metrics.

---

## Success Criteria

✅ **Functional**:
- All three modes implemented and tested
- Backward compatibility maintained
- No regressions in existing tests

✅ **Performance**:
- Relaxed mode is 2-5x faster than Strong
- Eventual mode is 10-100x faster than Strong
- Strong mode performance unchanged

✅ **Quality**:
- >90% test coverage for new code
- All tests pass with `-race` flag
- Documentation complete and clear

✅ **Usability**:
- API is intuitive and hard to misuse
- Examples cover common scenarios
- Migration path is clear

---

**Status**: 🚀 Ready for Implementation  
**Next Step**: Begin Phase 1 - API Introduction