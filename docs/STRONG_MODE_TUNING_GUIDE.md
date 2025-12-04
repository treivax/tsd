# Strong Mode - Performance Tuning Guide

**Version**: 1.0  
**Last Updated**: 2025-12-04  
**Target Audience**: DevOps, Performance Engineers, Application Developers

---

## Table of Contents

1. [Overview](#overview)
2. [Quick Start](#quick-start)
3. [Understanding the Parameters](#understanding-the-parameters)
4. [Performance Profiling](#performance-profiling)
5. [Tuning Strategies](#tuning-strategies)
6. [Common Scenarios](#common-scenarios)
7. [Monitoring & Alerting](#monitoring--alerting)
8. [Troubleshooting](#troubleshooting)
9. [Best Practices](#best-practices)

---

## Overview

Strong mode provides the strictest consistency guarantees in TSD, ensuring that all reads reflect the most recent writes. While the default configuration works well for most scenarios, tuning the parameters can significantly improve performance for your specific workload.

### Key Concepts

- **Verification**: The process of confirming that a fact has been persisted to storage
- **Retry**: Re-attempting verification when initial attempts fail
- **Timeout**: Maximum time to wait for an operation to complete
- **Backoff**: Delay between retry attempts (exponential)

### Performance Impact

Well-tuned parameters can:
- ✅ Reduce transaction latency by 30-50%
- ✅ Increase throughput by 2-3x
- ✅ Minimize resource consumption (CPU, network)
- ✅ Improve user experience

---

## Quick Start

### Step 1: Use Default Configuration (Baseline)

Start with the defaults and collect metrics:

```go
network := rete.NewReteNetwork(storage, logger)

// Use default configuration
tx := network.BeginTransaction()
defer tx.Rollback()

// ... add facts ...

err := tx.Commit()
```

**Default Values:**
- `SubmissionTimeout`: 30 seconds
- `VerifyRetryDelay`: 50ms
- `MaxVerifyRetries`: 10
- `VerifyOnCommit`: true

### Step 2: Monitor Performance

Collect metrics for at least 100 transactions:

```go
perfMetrics := rete.NewStrongModePerformanceMetrics()

// For each transaction:
coherenceMetrics := tx.GetCoherenceMetrics()
perfMetrics.RecordTransaction(duration, factCount, success, coherenceMetrics)

// After collecting data:
fmt.Println(perfMetrics.GetReport())
```

### Step 3: Apply Recommendations

The performance metrics system will automatically suggest tuning adjustments:

```go
report := perfMetrics.GetReport()
// Check "Recommendations" section for specific guidance
```

### Step 4: Tune & Re-measure

Apply changes incrementally and measure impact:

```go
opts := rete.DefaultTransactionOptions()
opts.VerifyRetryDelay = 20 * time.Millisecond  // Reduce from 50ms

tx := network.BeginTransactionWithOptions(opts)
```

---

## Understanding the Parameters

### 1. SubmissionTimeout

**What it does:** Maximum time to wait for the entire batch of facts to be submitted to storage.

**Default:** 30 seconds

**When to tune:**
- ✅ Increase if you see timeout errors
- ✅ Decrease for faster failure detection in dev/test environments
- ⚠️ Never set below 5 seconds for production

**Impact:**
- **Too low**: Premature timeouts, transaction failures
- **Too high**: Slow failure detection, resource waste

**Calculation:**
```
Recommended = (avg_verification_time × max_facts_per_tx × 2) + buffer
```

Example:
- Avg verification: 10ms
- Max facts per tx: 100
- Calculation: (10ms × 100 × 2) + 5s = 7 seconds

### 2. VerifyRetryDelay

**What it does:** Initial delay between verification retry attempts. Doubles on each retry (exponential backoff).

**Default:** 50ms

**When to tune:**
- ✅ Decrease if storage is fast and responsive (< 10ms latency)
- ✅ Increase if storage is slow or variable (> 100ms latency)
- ⚠️ Consider storage propagation delay

**Impact:**
- **Too low**: Excessive retry attempts, CPU waste
- **Too high**: Slow recovery, poor throughput

**Retry Schedule (examples):**

| VerifyRetryDelay | Attempt 1 | Attempt 2 | Attempt 3 | Attempt 4 | Total |
|------------------|-----------|-----------|-----------|-----------|-------|
| 20ms             | 20ms      | 40ms      | 80ms      | 160ms     | 300ms |
| 50ms (default)   | 50ms      | 100ms     | 200ms     | 400ms     | 750ms |
| 100ms            | 100ms     | 200ms     | 400ms     | 800ms     | 1.5s  |

**Calculation:**
```
Recommended = storage_latency_p95 / 2
```

Example:
- Storage P95 latency: 40ms
- Calculation: 40ms / 2 = 20ms

### 3. MaxVerifyRetries

**What it does:** Maximum number of verification attempts before giving up.

**Default:** 10

**When to tune:**
- ✅ Increase for eventually consistent storage backends
- ✅ Decrease if storage is strongly consistent (e.g., PostgreSQL)
- ⚠️ Monitor actual retry counts first

**Impact:**
- **Too low**: Premature failures, data loss risk
- **Too high**: Slow failure detection, resource waste

**Total Wait Time:**
```
max_wait = VerifyRetryDelay × (2^MaxVerifyRetries - 1)
```

Examples:
- 50ms × (2^10 - 1) = ~51 seconds (default)
- 50ms × (2^5 - 1) = ~1.5 seconds (fast storage)
- 50ms × (2^15 - 1) = ~27 minutes (eventual consistency)

**Calculation:**
```
Recommended = ceil(log2(desired_max_wait / VerifyRetryDelay))
```

Example:
- Desired max wait: 5 seconds
- VerifyRetryDelay: 50ms
- Calculation: ceil(log2(5000ms / 50ms)) = ceil(log2(100)) = 7 retries

### 4. VerifyOnCommit

**What it does:** Re-verify all facts in the transaction on commit.

**Default:** true

**When to tune:**
- ✅ Set to `false` for high-throughput scenarios where facts are already verified inline
- ⚠️ Only disable if you're confident in your storage consistency

**Impact:**
- **true** (default): Maximum safety, slower commits
- **false**: Faster commits, small risk of undetected failures

**Use Cases:**
- **Keep true**: Financial transactions, critical data, regulatory compliance
- **Consider false**: Analytics pipelines, event streams, high-volume inserts

---

## Performance Profiling

### Collecting Baseline Metrics

```go
package main

import (
    "fmt"
    "time"
    "your-project/rete"
)

func main() {
    network := rete.NewReteNetwork(storage, logger)
    perfMetrics := rete.NewStrongModePerformanceMetrics()
    
    // Run your workload
    for i := 0; i < 1000; i++ {
        start := time.Now()
        
        tx := network.BeginTransaction()
        
        // Add facts
        for j := 0; j < 10; j++ {
            tx.AddFact(fmt.Sprintf("fact-%d-%d", i, j), map[string]interface{}{
                "value": j,
            })
        }
        
        err := tx.Commit()
        duration := time.Since(start)
        
        // Record metrics
        coherenceMetrics := tx.GetCoherenceMetrics()
        perfMetrics.RecordTransaction(duration, 10, err == nil, coherenceMetrics)
    }
    
    // Generate report
    fmt.Println(perfMetrics.GetReport())
    
    // Check health
    if !perfMetrics.IsHealthy {
        fmt.Println("\n⚠️  System needs tuning!")
        for _, rec := range perfMetrics.Recommendations {
            fmt.Println("   ", rec)
        }
    }
}
```

### Key Metrics to Monitor

1. **AvgRetriesPerFact**: Should be < 0.5 for optimal performance
2. **TimeoutRate**: Should be < 1%
3. **AvgVerificationTime**: Baseline for tuning VerifyRetryDelay
4. **SuccessRate**: Should be > 99%
5. **HealthScore**: Should be > 80

---

## Tuning Strategies

### Strategy 1: Fast Storage (< 10ms latency)

**Symptoms:**
- Low retry count (< 0.1 per fact)
- Fast verification (< 10ms average)
- High success rate (> 99.9%)

**Optimization:**
```go
opts := rete.DefaultTransactionOptions()
opts.VerifyRetryDelay = 10 * time.Millisecond  // Reduce from 50ms
opts.MaxVerifyRetries = 5                       // Reduce from 10
opts.SubmissionTimeout = 10 * time.Second       // Reduce from 30s
```

**Expected Improvement:**
- ✅ 40-50% reduction in transaction latency
- ✅ 2-3x increase in throughput
- ⚠️ Slight increase in retry rate (acceptable)

### Strategy 2: Slow Storage (> 100ms latency)

**Symptoms:**
- High retry count (> 1.0 per fact)
- Slow verification (> 100ms average)
- Occasional timeouts

**Optimization:**
```go
opts := rete.DefaultTransactionOptions()
opts.VerifyRetryDelay = 100 * time.Millisecond  // Increase from 50ms
opts.MaxVerifyRetries = 15                       // Increase from 10
opts.SubmissionTimeout = 60 * time.Second        // Increase from 30s
```

**Expected Improvement:**
- ✅ Reduced retry waste
- ✅ Fewer timeouts
- ⚠️ Slower overall, but more reliable

### Strategy 3: Variable Latency Storage

**Symptoms:**
- High variance in verification time
- Occasional timeout spikes
- Inconsistent retry counts

**Optimization:**
```go
opts := rete.DefaultTransactionOptions()
opts.VerifyRetryDelay = 75 * time.Millisecond   // Moderate delay
opts.MaxVerifyRetries = 12                       // More attempts
opts.SubmissionTimeout = 45 * time.Second        // Conservative timeout
```

**Expected Improvement:**
- ✅ Better handling of latency spikes
- ✅ Fewer unnecessary retries during fast periods
- ⚠️ Slight performance trade-off for reliability

### Strategy 4: High Throughput, Best Effort

**Symptoms:**
- Need maximum throughput
- Can tolerate rare failures
- Analytics or non-critical data

**Optimization:**
```go
opts := rete.DefaultTransactionOptions()
opts.VerifyRetryDelay = 20 * time.Millisecond
opts.MaxVerifyRetries = 3                    // Minimal retries
opts.SubmissionTimeout = 5 * time.Second     // Fast timeout
opts.VerifyOnCommit = false                  // Skip commit verification
```

**Expected Improvement:**
- ✅ Maximum throughput (3-5x)
- ✅ Minimum latency
- ⚠️ Slight increase in failure rate (0.1-1%)

---

## Common Scenarios

### Scenario 1: PostgreSQL/MySQL (Synchronous Replication)

**Characteristics:**
- Strong consistency
- Low latency (< 10ms)
- Reliable

**Recommended Config:**
```go
opts := &rete.TransactionOptions{
    SubmissionTimeout: 10 * time.Second,
    VerifyRetryDelay:  10 * time.Millisecond,
    MaxVerifyRetries:  5,
    VerifyOnCommit:    true,
}
```

**Why:**
- Storage confirms writes immediately
- Few retries needed
- Fast feedback loop

### Scenario 2: Cassandra/DynamoDB (Eventually Consistent)

**Characteristics:**
- Eventual consistency
- Variable latency (10-200ms)
- Network-dependent

**Recommended Config:**
```go
opts := &rete.TransactionOptions{
    SubmissionTimeout: 45 * time.Second,
    VerifyRetryDelay:  100 * time.Millisecond,
    MaxVerifyRetries:  12,
    VerifyOnCommit:    true,
}
```

**Why:**
- Need more time for replication
- More retries handle propagation delay
- Conservative timeout prevents false failures

### Scenario 3: Redis (In-Memory)

**Characteristics:**
- Very fast (< 1ms)
- Synchronous
- High throughput

**Recommended Config:**
```go
opts := &rete.TransactionOptions{
    SubmissionTimeout: 5 * time.Second,
    VerifyRetryDelay:  5 * time.Millisecond,
    MaxVerifyRetries:  3,
    VerifyOnCommit:    false,  // Optional: Redis is synchronous
}
```

**Why:**
- Extremely fast verification
- Minimal retries
- Can skip commit verification for max performance

### Scenario 4: Cloud Storage (S3, GCS)

**Characteristics:**
- High latency (100-500ms)
- Eventually consistent
- Network-dependent

**Recommended Config:**
```go
opts := &rete.TransactionOptions{
    SubmissionTimeout: 60 * time.Second,
    VerifyRetryDelay:  200 * time.Millisecond,
    MaxVerifyRetries:  15,
    VerifyOnCommit:    true,
}
```

**Why:**
- Network round-trips are expensive
- Need aggressive retry strategy
- Long timeout for worst-case scenarios

### Scenario 5: Development/Testing

**Characteristics:**
- Fast iteration
- Mock storage or in-memory
- Want quick failures

**Recommended Config:**
```go
opts := &rete.TransactionOptions{
    SubmissionTimeout: 2 * time.Second,   // Fail fast
    VerifyRetryDelay:  5 * time.Millisecond,
    MaxVerifyRetries:  2,                  // Minimal retries
    VerifyOnCommit:    false,              // Speed over safety
}
```

**Why:**
- Quick feedback during development
- Don't waste time on retries
- Failures indicate real bugs

---

## Monitoring & Alerting

### Essential Metrics

Track these metrics in your monitoring system (Prometheus, Datadog, etc.):

```yaml
# Transaction Success Rate
- metric: strong_mode_transaction_success_rate
  alert_threshold: < 99%
  severity: warning

# Average Retries Per Fact
- metric: strong_mode_avg_retries_per_fact
  alert_threshold: > 1.0
  severity: warning
  alert_threshold: > 2.0
  severity: critical

# Timeout Rate
- metric: strong_mode_timeout_rate
  alert_threshold: > 1%
  severity: warning
  alert_threshold: > 5%
  severity: critical

# Health Score
- metric: strong_mode_health_score
  alert_threshold: < 80
  severity: warning
  alert_threshold: < 60
  severity: critical

# Average Transaction Duration
- metric: strong_mode_avg_transaction_duration
  alert_threshold: > 5s
  severity: warning
  baseline: use_historical_average
```

### Grafana Dashboard Example

```json
{
  "dashboard": {
    "title": "TSD Strong Mode Performance",
    "panels": [
      {
        "title": "Transaction Success Rate",
        "targets": [
          {
            "expr": "rate(strong_mode_successful_transactions[5m]) / rate(strong_mode_total_transactions[5m]) * 100"
          }
        ]
      },
      {
        "title": "Average Retries Per Fact",
        "targets": [
          {
            "expr": "strong_mode_total_retries / strong_mode_total_facts_persisted"
          }
        ]
      },
      {
        "title": "Verification Latency (P50, P95, P99)",
        "targets": [
          {
            "expr": "histogram_quantile(0.50, strong_mode_verification_time)"
          },
          {
            "expr": "histogram_quantile(0.95, strong_mode_verification_time)"
          },
          {
            "expr": "histogram_quantile(0.99, strong_mode_verification_time)"
          }
        ]
      }
    ]
  }
}
```

### Alert Response Playbook

#### Alert: High Timeout Rate

**Symptoms:**
- Timeout rate > 5%
- Health score drops

**Investigation:**
1. Check storage latency metrics
2. Verify network connectivity
3. Check storage backend health

**Resolution:**
1. Increase `SubmissionTimeout` by 50%
2. Increase `MaxVerifyRetries` by 2-3
3. Monitor for improvement
4. If persists, investigate storage infrastructure

#### Alert: High Retry Rate

**Symptoms:**
- Avg retries per fact > 2.0
- Slow transaction completion

**Investigation:**
1. Check storage propagation delay
2. Verify `VerifyRetryDelay` is appropriate
3. Check for storage replication lag

**Resolution:**
1. Increase `VerifyRetryDelay` by 50%
2. Reduce `MaxVerifyRetries` slightly
3. Investigate storage replication issues

---

## Troubleshooting

### Problem: High Latency

**Symptoms:**
```
Average transaction time: 5+ seconds
Health Score: B or C
Many timeouts
```

**Diagnostic Steps:**
```go
// 1. Check metrics breakdown
report := perfMetrics.GetReport()
fmt.Println(report)

// 2. Look at verification time
if perfMetrics.AvgVerificationTime > 100*time.Millisecond {
    fmt.Println("Storage is slow")
}

// 3. Check retry pattern
if perfMetrics.AvgRetriesPerFact > 1.5 {
    fmt.Println("Too many retries")
}
```

**Solutions:**
1. Increase `VerifyRetryDelay` to reduce wasted retries
2. Optimize storage queries (add indexes, etc.)
3. Consider batching facts if possible
4. Check network latency to storage

### Problem: Frequent Timeouts

**Symptoms:**
```
Timeout rate: > 5%
Many failed transactions
```

**Diagnostic Steps:**
```go
// Check if timeouts correlate with large transactions
if perfMetrics.AvgFactsPerTransaction > 100 {
    fmt.Println("Large transactions may need longer timeout")
}

// Check timeout timing
if perfMetrics.AvgTimeToTimeout < perfMetrics.CurrentConfig.SubmissionTimeout {
    fmt.Println("Hitting timeout prematurely")
}
```

**Solutions:**
1. Increase `SubmissionTimeout` proportional to transaction size
2. Break large transactions into smaller chunks
3. Investigate storage performance issues

### Problem: Excessive Retries

**Symptoms:**
```
Avg retries per fact: > 2.0
High CPU usage
Slow throughput
```

**Diagnostic Steps:**
```go
// Check retry delay
if perfMetrics.CurrentConfig.VerifyRetryDelay < 10*time.Millisecond {
    fmt.Println("Retry delay too aggressive")
}

// Check storage latency
if perfMetrics.AvgVerificationTime > 50*time.Millisecond {
    fmt.Println("Storage slower than retry delay")
}
```

**Solutions:**
1. Increase `VerifyRetryDelay` to match storage latency
2. Add jitter to avoid thundering herd
3. Investigate storage write amplification

---

## Best Practices

### 1. Start Conservative, Optimize Gradually

```go
// Phase 1: Baseline with defaults (1-2 weeks)
opts := rete.DefaultTransactionOptions()

// Phase 2: Collect metrics and tune conservatively
opts.VerifyRetryDelay = measuredLatency * 1.5  // 50% buffer

// Phase 3: Optimize based on data
if metrics.HealthScore > 95 && metrics.AvgRetriesPerFact < 0.2 {
    opts.MaxVerifyRetries = 5  // Reduce if safe
}
```

### 2. Use Different Configs for Different Workloads

```go
// Critical financial transactions
criticalOpts := &rete.TransactionOptions{
    SubmissionTimeout: 60 * time.Second,
    VerifyRetryDelay:  100 * time.Millisecond,
    MaxVerifyRetries:  15,
    VerifyOnCommit:    true,
}

// Analytics batch processing
batchOpts := &rete.TransactionOptions{
    SubmissionTimeout: 10 * time.Second,
    VerifyRetryDelay:  20 * time.Millisecond,
    MaxVerifyRetries:  5,
    VerifyOnCommit:    false,
}

// Use appropriate config per use case
tx := network.BeginTransactionWithOptions(criticalOpts)
```

### 3. Monitor Continuously

```go
// Set up periodic health checks
go func() {
    ticker := time.NewTicker(5 * time.Minute)
    for range ticker.C {
        if !perfMetrics.IsHealthy {
            logger.Warn("Strong mode health degraded",
                "score", perfMetrics.HealthScore,
                "grade", perfMetrics.PerformanceGrade)
            
            // Auto-tune or alert
            for _, rec := range perfMetrics.Recommendations {
                logger.Info("Recommendation", "msg", rec)
            }
        }
    }
}()
```

### 4. Document Your Configuration

```go
// Document WHY you chose these values
var ProductionConfig = &rete.TransactionOptions{
    // Tuned for PostgreSQL with avg 15ms write latency
    // Measured over 100K transactions, 2025-12-04
    SubmissionTimeout: 15 * time.Second,
    
    // Half of measured P95 latency (30ms)
    VerifyRetryDelay:  15 * time.Millisecond,
    
    // Covers up to 450ms worst case (15ms * 2^5)
    MaxVerifyRetries:  5,
    
    // Required for compliance
    VerifyOnCommit:    true,
}
```

### 5. Test Configuration Changes

```go
func TestConfigurationPerformance(t *testing.T) {
    configs := []*rete.TransactionOptions{
        rete.DefaultTransactionOptions(),
        optimizedConfig,
        aggressiveConfig,
    }
    
    for _, config := range configs {
        metrics := benchmarkConfig(config, 1000)
        t.Logf("Config: %+v\nResults: %s", config, metrics.GetSummary())
        
        require.True(t, metrics.IsHealthy, "Config should be healthy")
        require.Greater(t, metrics.HealthScore, 80.0, "Score too low")
    }
}
```

---

## Summary

### Quick Reference Table

| Storage Type | SubmissionTimeout | VerifyRetryDelay | MaxVerifyRetries | VerifyOnCommit |
|--------------|-------------------|------------------|------------------|----------------|
| PostgreSQL   | 10s               | 10ms             | 5                | true           |
| MySQL        | 10s               | 10ms             | 5                | true           |
| Redis        | 5s                | 5ms              | 3                | false          |
| Cassandra    | 45s               | 100ms            | 12               | true           |
| DynamoDB     | 45s               | 100ms            | 12               | true           |
| S3/GCS       | 60s               | 200ms            | 15               | true           |
| In-Memory    | 2s                | 5ms              | 2                | false          |

### Decision Tree

```
Start with defaults
    ↓
Collect 100+ transaction metrics
    ↓
Check Health Score
    ↓
    ├─ > 90 (Grade A)
    │   └─ Optimize for speed (reduce retries/delays)
    │
    ├─ 80-90 (Grade B)
    │   └─ Keep current config, monitor
    │
    ├─ 70-80 (Grade C)
    │   └─ Investigate specific issues
    │
    └─ < 70 (Grade D/F)
        └─ Increase timeouts/retries, investigate storage
```

---

## Additional Resources

- **API Documentation**: See `rete/coherence_mode.go`
- **Design Document**: `docs/PHASE4_COHERENCE_STRONG_MODE.md`
- **Completion Report**: `docs/PHASE4_STRONG_MODE_COMPLETION.md`
- **Metrics API**: `rete/strong_mode_performance.go`

---

**Questions?** Check the recommendations in your performance report or contact the TSD team.