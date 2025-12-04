# Phase 4 - Strong Coherence Mode Implementation

**Feature**: Strong Consistency Mode (Simplified)  
**Status**: 🚧 Implementation Phase  
**Target Completion**: 2025-12-04  
**Estimated Effort**: 4-6 hours

---

## Executive Summary

This document outlines the simplified implementation of **Strong Consistency Mode** for the TSD RETE engine. After review, the decision was made to implement only Strong mode initially, removing the complexity of multiple coherence modes (Relaxed/Eventual) and focusing on delivering robust, predictable consistency guarantees.

### Strong Consistency Overview

**Strong Consistency** provides the strictest guarantees:
- All reads reflect the most recent write immediately
- Synchronous verification of all facts
- Retry mechanism with exponential backoff
- Read-after-write consistency enforced
- Suitable for ALL production use cases

---

## Motivation

### Current State (Phase 3)
- ✅ Thread-safe transaction model implemented
- ✅ Read-after-write guarantees enforced via `waitForFactPersistence()`
- ✅ Synchronous fact verification with retries
- ⚠️ **Verification logic embedded in various places**
- ⚠️ **No centralized configuration**
- ⚠️ **No metrics collection**

### Goals for Phase 4
1. **Formalize Strong Mode**: Extract and centralize existing verification logic
2. **Configuration**: Provide tunable parameters (timeout, retries, delays)
3. **Metrics**: Track verification attempts, successes, failures
4. **API Consistency**: Clean transaction API with options
5. **Documentation**: Clear guidelines for users

---

## Design Principles

1. **Simplicity First**: One mode, well-implemented, is better than three half-baked modes
2. **Backward Compatibility**: Existing code continues to work unchanged
3. **Explicit Configuration**: Users can tune parameters when needed
4. **Observable Behavior**: Metrics expose what's happening
5. **Safe Defaults**: Conservative timeouts and retries

---

## Strong Mode Guarantees

### What is Guaranteed

✅ **Read-After-Write Consistency**
- If you add a fact in a transaction, it is immediately readable after commit
- No stale reads within the same application

✅ **Synchronous Verification**
- Every `AddFact` is verified in storage before the transaction continues
- Storage operations are blocking

✅ **Retry Mechanism**
- Automatic retries with exponential backoff
- Configurable: max retries (default: 10), initial delay (default: 50ms)

✅ **Atomic Transactions**
- All facts in a transaction are persisted, or none are (on error)
- Rollback on any verification failure

✅ **No Data Loss**
- Any storage failure causes transaction failure
- Errors are propagated to the caller

---

## Technical Implementation

### Configuration Structure

```go
// TransactionOptions configures transaction behavior
type TransactionOptions struct {
    // SubmissionTimeout: maximum time for fact batch submission
    // Default: 30 seconds
    SubmissionTimeout time.Duration
    
    // VerifyRetryDelay: initial delay between verification attempts
    // Doubled on each retry (exponential backoff)
    // Default: 50ms
    VerifyRetryDelay time.Duration
    
    // MaxVerifyRetries: maximum number of verification attempts
    // Default: 10
    MaxVerifyRetries int
    
    // VerifyOnCommit: whether to re-verify all facts on commit
    // Default: true
    VerifyOnCommit bool
}

// DefaultTransactionOptions returns safe defaults
func DefaultTransactionOptions() *TransactionOptions {
    return &TransactionOptions{
        SubmissionTimeout: 30 * time.Second,
        VerifyRetryDelay:  50 * time.Millisecond,
        MaxVerifyRetries:  10,
        VerifyOnCommit:    true,
    }
}
```

### Verification Algorithm

```go
// Pseudo-code for Strong mode verification
func verifyFactStrong(fact *Fact, opts *TransactionOptions) error {
    retryDelay := opts.VerifyRetryDelay
    
    for attempt := 0; attempt < opts.MaxVerifyRetries; attempt++ {
        // Attempt to retrieve fact from storage
        retrieved := storage.GetFact(fact.InternalID)
        
        if retrieved != nil && retrieved.ID == fact.ID {
            logger.Debug("✓ Fact verified: %s (attempt %d)", fact.ID, attempt+1)
            metrics.RecordVerificationSuccess()
            return nil
        }
        
        // Not found yet, retry with backoff
        if attempt < opts.MaxVerifyRetries-1 {
            logger.Debug("Fact %s not found, retry in %v", fact.ID, retryDelay)
            time.Sleep(retryDelay)
            retryDelay *= 2 // Exponential backoff
        }
        
        metrics.RecordVerificationAttempt()
    }
    
    // Failed after all retries
    metrics.RecordVerificationFailure()
    return fmt.Errorf("fact %s not verified after %d attempts", 
        fact.ID, opts.MaxVerifyRetries)
}
```

### Transaction Lifecycle with Strong Mode

```
1. BEGIN TRANSACTION
   ↓
2. RecordAndExecute(AddFactCommand)
   ↓
3. Execute: storage.AddFact(fact)
   ↓
4. Verify: verifyFactStrong(fact, opts)
   ↓ (retry up to 10 times with exponential backoff)
5. Fact verified ✓
   ↓
6. Next command or...
   ↓
7. COMMIT
   ↓
8. Optional: Re-verify all facts (if VerifyOnCommit=true)
   ↓
9. Return SUCCESS
```

On any error at steps 3-8: **ROLLBACK** entire transaction.

---

## API Design

### Basic Usage (Backward Compatible)

```go
// Existing code continues to work
tx := network.BeginTransaction()
tx.RecordAndExecute(AddFactCommand{...})
err := tx.Commit()
```

### With Custom Options

```go
// Custom timeout for high-latency storage
opts := &TransactionOptions{
    SubmissionTimeout: 60 * time.Second,  // Longer timeout
    VerifyRetryDelay:  100 * time.Millisecond,
    MaxVerifyRetries:  20,  // More retries
    VerifyOnCommit:    true,
}

tx := network.BeginTransactionWithOptions(opts)
tx.RecordAndExecute(AddFactCommand{...})
err := tx.Commit()
```

### Global Configuration

```go
// Configure default options for the network
config := rete.NetworkCoherenceConfig{
    DefaultOptions: rete.TransactionOptions{
        SubmissionTimeout: 45 * time.Second,
        VerifyRetryDelay:  75 * time.Millisecond,
        MaxVerifyRetries:  15,
        VerifyOnCommit:    true,
    },
    EnableMetrics: true,
}

network := rete.NewReteNetworkWithConfig(storage, config)
```

---

## Metrics and Observability

### Metrics Collected

```go
type CoherenceMetrics struct {
    // Transaction tracking
    TotalTransactions int64
    TransactionCount  int64
    
    // Verification tracking
    VerificationAttempts  int64  // Total attempts across all facts
    VerificationSuccesses int64  // Successful verifications
    VerificationFailures  int64  // Failed verifications
    
    // Performance tracking
    TotalTransactionDuration time.Duration
    
    // Computed metrics
    AverageTransactionDuration() time.Duration
    VerificationSuccessRate() float64  // Percentage
}
```

### Usage Example

```go
// Access metrics
metrics := network.GetCoherenceMetrics()

fmt.Printf("Total transactions: %d\n", metrics.TotalTransactions)
fmt.Printf("Verification success rate: %.2f%%\n", metrics.VerificationSuccessRate())
fmt.Printf("Average transaction duration: %v\n", metrics.AverageTransactionDuration())
```

### Logging

Strong mode produces detailed logs at appropriate levels:

```go
// DEBUG level
logger.Debug("Beginning transaction with Strong mode")
logger.Debug("Verifying fact: %s (attempt %d/%d)", factID, attempt, maxRetries)
logger.Debug("✓ Fact verified: %s", factID)

// WARN level
logger.Warn("Fact verification slow: %s (attempt %d, elapsed %v)", 
    factID, attempt, elapsed)

// ERROR level
logger.Error("✗ Fact verification failed: %s after %d attempts", 
    factID, maxRetries)
```

---

## Integration with Existing Code

### Transaction Structure Update

```go
type Transaction struct {
    ID       string
    Network  *ReteNetwork
    Commands []*Command
    Options  *TransactionOptions  // NEW: configuration
    Metrics  *CoherenceMetrics    // NEW: per-transaction metrics
    Logger   *Logger
    mu       sync.Mutex
}
```

### Network Extensions

```go
type ReteNetwork struct {
    // ... existing fields ...
    
    CoherenceConfig  NetworkCoherenceConfig  // NEW: global config
    CoherenceMetrics *CoherenceMetrics       // NEW: global metrics
}

// NEW: Begin transaction with custom options
func (n *ReteNetwork) BeginTransactionWithOptions(opts *TransactionOptions) *Transaction {
    if opts == nil {
        opts = &n.CoherenceConfig.DefaultOptions
    }
    
    return &Transaction{
        ID:      generateTransactionID(),
        Network: n,
        Options: opts,
        Metrics: &CoherenceMetrics{},
        Logger:  n.Logger,
    }
}

// Existing method uses default options
func (n *ReteNetwork) BeginTransaction() *Transaction {
    return n.BeginTransactionWithOptions(nil)
}
```

---

## Testing Strategy

### Unit Tests

1. **TransactionOptions validation**
   - Valid options accepted
   - Invalid options (negative values) rejected

2. **Verification logic**
   - Successful verification on first attempt
   - Successful verification after retries
   - Failure after max retries exceeded
   - Exponential backoff verified

3. **Metrics collection**
   - Counters incremented correctly
   - Success rate calculated correctly
   - Average duration computed correctly

### Integration Tests

1. **Transaction with Strong mode**
   - Facts verified synchronously
   - Transaction commits successfully
   - Facts readable immediately after commit

2. **Custom options**
   - Custom timeout respected
   - Custom retry count respected
   - Custom delay respected

3. **Error scenarios**
   - Storage failure causes transaction failure
   - Verification timeout causes transaction failure
   - Rollback cleans up partial state

### Performance Tests

1. **Baseline performance**
   - Measure throughput with default options
   - Measure latency percentiles (p50, p95, p99)

2. **Tuning parameters**
   - Impact of retry count on success rate
   - Impact of retry delay on performance
   - Optimal configuration for different storage backends

---

## Performance Characteristics

### Expected Performance

| Metric | Value | Notes |
|--------|-------|-------|
| **Throughput** | 100-1,000 facts/sec | Depends on storage latency |
| **Latency (p50)** | 10-50ms | Fast storage, first-attempt verification |
| **Latency (p95)** | 50-200ms | Includes retries, exponential backoff |
| **Latency (p99)** | 200-500ms | Multiple retries, slow storage |
| **Verification overhead** | 30-40% | Compared to no verification |
| **Success rate** | >99% | With default retry settings |

### Tuning Guidelines

**For Fast Storage (e.g., in-memory, local SSD)**
```go
opts := &TransactionOptions{
    VerifyRetryDelay: 10 * time.Millisecond,  // Shorter delay
    MaxVerifyRetries: 5,                       // Fewer retries needed
}
```

**For Slow Storage (e.g., network, remote DB)**
```go
opts := &TransactionOptions{
    SubmissionTimeout: 60 * time.Second,      // Longer timeout
    VerifyRetryDelay: 100 * time.Millisecond, // Longer initial delay
    MaxVerifyRetries: 20,                      // More retries
}
```

**For High-Volume Scenarios**
```go
opts := &TransactionOptions{
    VerifyOnCommit: false,  // Skip re-verification at commit
}
// Note: Less safe, but 20-30% faster
```

---

## Migration Path

### Phase 4.1: Core Implementation (This Phase) ✅
- [x] Create `coherence_mode.go` with types and configuration
- [ ] Update `Transaction` to use `TransactionOptions`
- [ ] Extract verification logic into `verifyFactStrong()`
- [ ] Implement `BeginTransactionWithOptions()`
- [ ] Add metrics collection
- [ ] Unit tests for verification logic

### Phase 4.2: Integration (Next)
- [ ] Update existing transaction code to use Strong mode
- [ ] Add logging at appropriate levels
- [ ] Integration tests
- [ ] Performance benchmarks

### Phase 4.3: Documentation (Final)
- [ ] User guide with examples
- [ ] Tuning guide for different scenarios
- [ ] API reference documentation
- [ ] Migration guide for existing users

---

## Examples

### Example 1: Basic Transaction (Default Strong Mode)

```go
// Create network
storage := rete.NewMemoryStorage()
network := rete.NewReteNetwork(storage)

// Begin transaction (uses default Strong mode options)
tx := network.BeginTransaction()

// Add facts
tx.RecordAndExecute(rete.AddFactCommand{
    Fact: &rete.Fact{
        ID:   "user_123",
        Type: "User",
        Fields: map[string]interface{}{
            "id":   "user_123",
            "name": "Alice",
        },
    },
})

// Commit (all facts verified)
if err := tx.Commit(); err != nil {
    log.Fatalf("Transaction failed: %v", err)
}

// Facts are guaranteed to be readable now
fact := storage.GetFact("user_123")
fmt.Printf("User: %v\n", fact.Fields["name"])
```

### Example 2: Custom Options for Slow Storage

```go
// Configure for network storage with higher latency
opts := &rete.TransactionOptions{
    SubmissionTimeout: 60 * time.Second,
    VerifyRetryDelay:  100 * time.Millisecond,
    MaxVerifyRetries:  20,
    VerifyOnCommit:    true,
}

tx := network.BeginTransactionWithOptions(opts)

// Bulk insert with slower verification
for i := 0; i < 1000; i++ {
    tx.RecordAndExecute(rete.AddFactCommand{
        Fact: &rete.Fact{
            ID:   fmt.Sprintf("item_%d", i),
            Type: "Item",
            Fields: map[string]interface{}{
                "id":    fmt.Sprintf("item_%d", i),
                "index": i,
            },
        },
    })
}

if err := tx.Commit(); err != nil {
    log.Fatalf("Bulk insert failed: %v", err)
}
```

### Example 3: Global Configuration

```go
// Set default options for all transactions
config := rete.NetworkCoherenceConfig{
    DefaultOptions: rete.TransactionOptions{
        SubmissionTimeout: 45 * time.Second,
        VerifyRetryDelay:  75 * time.Millisecond,
        MaxVerifyRetries:  15,
        VerifyOnCommit:    true,
    },
    EnableMetrics: true,
}

// Create network with custom config
network := rete.NewReteNetworkWithConfig(storage, config)

// All transactions use these defaults
tx := network.BeginTransaction()
// ... uses SubmissionTimeout=45s, MaxVerifyRetries=15, etc.
```

### Example 4: Monitoring and Metrics

```go
// Get metrics after some transactions
metrics := network.GetCoherenceMetrics()

fmt.Printf("=== Coherence Metrics ===\n")
fmt.Printf("Total transactions: %d\n", metrics.TotalTransactions)
fmt.Printf("Verification attempts: %d\n", metrics.VerificationAttempts)
fmt.Printf("Verification successes: %d\n", metrics.VerificationSuccesses)
fmt.Printf("Verification failures: %d\n", metrics.VerificationFailures)
fmt.Printf("Success rate: %.2f%%\n", metrics.VerificationSuccessRate())
fmt.Printf("Average transaction duration: %v\n", metrics.AverageTransactionDuration())
```

---

## Error Handling

### Verification Failure

```go
tx := network.BeginTransaction()
tx.RecordAndExecute(AddFactCommand{...})

err := tx.Commit()
if err != nil {
    // Check for verification errors
    if strings.Contains(err.Error(), "not verified") {
        log.Printf("Verification failed - possible storage issue")
        // Retry or alert
    }
}
```

### Timeout

```go
// Set aggressive timeout for testing
opts := &TransactionOptions{
    SubmissionTimeout: 100 * time.Millisecond,
    MaxVerifyRetries:  2,
}

tx := network.BeginTransactionWithOptions(opts)
// ... add facts ...

err := tx.Commit()
if err != nil {
    if strings.Contains(err.Error(), "timeout") {
        log.Printf("Transaction timed out - storage too slow")
    }
}
```

---

## Success Criteria

### Functional Requirements
- ✅ Strong mode configuration types defined
- ⏳ Transaction uses TransactionOptions
- ⏳ Verification logic extracted and centralized
- ⏳ Metrics collected and exposed
- ⏳ All tests pass with -race flag

### Performance Requirements
- ⏳ No regression vs current implementation
- ⏳ Throughput: 100-1,000 facts/sec
- ⏳ Success rate: >99% with default settings
- ⏳ p95 latency: <200ms

### Quality Requirements
- ⏳ Unit test coverage >90%
- ⏳ Integration tests for all scenarios
- ⏳ Documentation complete
- ⏳ Examples for common use cases

---

## Future Enhancements (Out of Scope)

If needed in the future, these could be added:

1. **Async Verification** (post-commit background verification)
2. **Relaxed Mode** (bounded staleness)
3. **Batch Verification** (verify multiple facts in one operation)
4. **Adaptive Retry** (adjust retry strategy based on observed latency)

For now, Strong mode provides everything needed for production use.

---

## References

- Phase 3 Completion Report: `docs/PHASE3_COMPLETION.md`
- Logging Guide: `docs/LOGGING_GUIDE.md`
- Transaction API: `rete/transaction.go`
- Storage Interface: `rete/storage.go`

---

**Status**: 🚧 Ready for Implementation  
**Next Step**: Integrate TransactionOptions into Transaction struct and implement verification logic