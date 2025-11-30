# Multi-Source Aggregation Performance Guide

**Copyright (c) 2025 TSD Contributors**  
**Licensed under the MIT License**

## Table of Contents

1. [Overview](#overview)
2. [Quick Start](#quick-start)
3. [Benchmark Suite](#benchmark-suite)
4. [Profiling Tools](#profiling-tools)
5. [Performance Characteristics](#performance-characteristics)
6. [Optimization Strategies](#optimization-strategies)
7. [Monitoring & Metrics](#monitoring--metrics)
8. [Troubleshooting](#troubleshooting)

---

## Overview

This guide provides comprehensive information on profiling, benchmarking, and optimizing multi-source aggregation performance in the RETE engine. Multi-source aggregations involve joining multiple fact types and computing aggregate functions (AVG, SUM, COUNT, MIN, MAX) across the joined results.

### Performance Factors

Key factors affecting multi-source aggregation performance:

- **Number of main facts**: Primary entities being aggregated (e.g., Departments)
- **Number of source facts**: Facts being joined and aggregated (e.g., Employees, Performance records)
- **Join fanout**: Average number of source facts per main fact
- **Number of aggregation variables**: Count of aggregate functions being computed
- **Number of sources**: Count of distinct fact types being joined
- **Threshold complexity**: Number and type of threshold conditions

---

## Quick Start

### Running All Benchmarks

```bash
# Run the complete benchmark suite
cd rete
go test -bench=BenchmarkMultiSourceAggregation -benchmem -benchtime=5s

# Run with profiling
./scripts/profile_multi_source.sh
```

### Running Specific Benchmarks

```bash
# Small scale 2-source aggregation
go test -bench=BenchmarkMultiSourceAggregation_TwoSources_SmallScale -benchmem

# Large scale with profiling
go test -bench=BenchmarkMultiSourceAggregation_TwoSources_LargeScale \
    -cpuprofile=cpu.prof \
    -memprofile=mem.prof \
    -benchmem
```

### Quick Profile Analysis

```bash
# View CPU profile interactively
go tool pprof -http=:8080 profiles/cpu_two_sources_large.prof

# Show top CPU consumers
go tool pprof -top profiles/cpu_two_sources_large.prof

# Show top memory allocators
go tool pprof -top profiles/mem_two_sources_large.prof
```

---

## Benchmark Suite

### Benchmark Categories

#### 1. Scale Benchmarks

Test performance across different data sizes:

| Benchmark | Main Facts | Source Facts (each) | Purpose |
|-----------|-----------|---------------------|---------|
| `TwoSources_SmallScale` | 10 | 50 | Baseline performance |
| `TwoSources_MediumScale` | 100 | 500 | Typical workload |
| `TwoSources_LargeScale` | 1,000 | 5,000 | High-volume scenario |
| `ThreeSources_SmallScale` | 10 | 50 | 3-way join baseline |
| `ThreeSources_MediumScale` | 100 | 500 | 3-way join typical |
| `ThreeSources_LargeScale` | 500 | 2,500 | 3-way join high-volume |

#### 2. Fanout Benchmarks

Test join cardinality impact:

- **HighFanout**: 10 main facts, 1,000 source facts each (100:1 ratio)
- **LowFanout**: 100 main facts, 150 source facts each (1.5:1 ratio)

#### 3. Aggregation Function Benchmarks

- **ManyAggregates**: Tests with 5+ aggregation functions (AVG, SUM, COUNT, MIN, MAX)

#### 4. Operational Benchmarks

- **WithThresholds**: Performance impact of threshold evaluation
- **Retraction**: Fact retraction and recomputation overhead
- **IncrementalUpdate**: Single-fact addition performance

#### 5. Memory Benchmarks

- **Memory_SmallScale**: Memory usage for small datasets
- **Memory_LargeScale**: Memory usage for large datasets

### Interpreting Results

Example benchmark output:

```
BenchmarkMultiSourceAggregation_TwoSources_MediumScale-8
    1000    1234567 ns/op    456789 B/op    1234 allocs/op
    facts/sec: 890.5
    total_facts: 1100
    main_facts: 100
    agg_vars: 2
    activations: 100
```

**Key Metrics:**
- **ns/op**: Nanoseconds per operation (lower is better)
- **B/op**: Bytes allocated per operation (lower is better)
- **allocs/op**: Number of allocations per operation (lower is better)
- **facts/sec**: Throughput in facts processed per second (higher is better)
- **activations**: Number of rule activations (validates correctness)

---

## Profiling Tools

### CPU Profiling

#### Generate CPU Profile

```bash
go test -bench=BenchmarkMultiSourceAggregation_Profile \
    -cpuprofile=cpu.prof \
    -benchtime=30s
```

#### Analyze CPU Profile

```bash
# Interactive web UI (recommended)
go tool pprof -http=:8080 cpu.prof

# Command-line analysis
go tool pprof cpu.prof
(pprof) top 20          # Show top 20 CPU consumers
(pprof) list MultiSourceAccumulatorNode.Activate  # Show source
(pprof) web            # Generate graph (requires graphviz)
```

#### Key Areas to Examine

1. **Activation path**: `MultiSourceAccumulatorNode.Activate()`
2. **Aggregation computation**: `computeAggregation()`
3. **Join processing**: `JoinNode.Activate()`
4. **Token creation**: Token allocation and copying
5. **Field extraction**: Getting values from fact fields

### Memory Profiling

#### Generate Memory Profile

```bash
go test -bench=BenchmarkMultiSourceAggregation_Memory_LargeScale \
    -memprofile=mem.prof \
    -benchtime=10s
```

#### Analyze Memory Profile

```bash
# Interactive web UI
go tool pprof -http=:8080 mem.prof

# Show allocations
go tool pprof -alloc_space mem.prof
(pprof) top 20

# Show in-use memory
go tool pprof -inuse_space mem.prof
(pprof) top 20
```

#### Key Memory Hotspots

1. **Token allocation**: `CombinedTokens` map growth
2. **Cache storage**: `AggregateCache` size
3. **Fact storage**: `MainFacts` and binding maps
4. **String allocations**: Fact IDs, token IDs

### Trace Profiling

For detailed execution flow analysis:

```bash
# Generate trace
go test -bench=BenchmarkMultiSourceAggregation_TwoSources_MediumScale \
    -trace=trace.out \
    -benchtime=5s

# View trace
go tool trace trace.out
```

Trace view shows:
- Goroutine execution timeline
- GC pauses
- Network/syscall blocking
- Lock contention

---

## Performance Characteristics

### Time Complexity

| Operation | Complexity | Notes |
|-----------|-----------|-------|
| Fact submission | O(J × F) | J = join depth, F = fanout |
| Aggregation computation | O(T × A) | T = tokens, A = aggregation functions |
| Threshold evaluation | O(A) | Per aggregation variable |
| Fact retraction | O(T × A) | Must recompute affected aggregations |

### Space Complexity

| Structure | Complexity | Notes |
|-----------|-----------|-------|
| MainFacts | O(M) | M = main facts |
| CombinedTokens | O(M × T) | T = avg tokens per main fact |
| AggregateCache | O(M × A) | A = aggregation variables |
| Join chain memory | O(J × F × M) | J = join depth, F = fanout |

### Scalability Limits

Based on benchmark data:

| Scale | Main Facts | Source Facts | Performance | Memory |
|-------|-----------|--------------|-------------|--------|
| Small | < 100 | < 500 | Excellent | < 10 MB |
| Medium | 100-1,000 | 500-5,000 | Good | 10-100 MB |
| Large | 1,000-10,000 | 5,000-50,000 | Acceptable | 100-500 MB |
| Very Large | > 10,000 | > 50,000 | Needs optimization | > 500 MB |

---

## Optimization Strategies

### 1. Join Order Optimization

**Problem**: Join order significantly impacts performance.

**Solution**: Join more selective sources first.

```go
// Before: Join least selective first
// Employee (5000 facts) → Performance (5000 facts) → Training (5000 facts)

// After: Join most selective first
// Training (500 facts) → Employee (5000 facts) → Performance (5000 facts)
```

**Implementation**: Add selectivity hints or statistics to pattern ordering.

### 2. Early Threshold Evaluation

**Problem**: Computing all aggregations when thresholds will fail wastes CPU.

**Solution**: Evaluate thresholds incrementally.

```go
// Pseudocode for early threshold evaluation
for _, aggVar := range msn.AggregationVars {
    value := computeAggregation(...)
    
    // Check threshold immediately
    if aggVar.Threshold != 0 {
        if !evaluateThreshold(value, aggVar.Operator, aggVar.Threshold) {
            // Short-circuit: don't compute remaining aggregations
            return nil
        }
    }
}
```

**Expected Improvement**: 20-40% for threshold-heavy workloads.

### 3. Incremental Aggregation

**Problem**: Recomputing entire aggregations on each token is expensive.

**Solution**: Maintain running aggregates and update incrementally.

```go
type IncrementalAggregate struct {
    Count float64
    Sum   float64
    Min   float64
    Max   float64
}

// Update incrementally
func (ia *IncrementalAggregate) Add(value float64) {
    ia.Count++
    ia.Sum += value
    if value < ia.Min { ia.Min = value }
    if value > ia.Max { ia.Max = value }
}

func (ia *IncrementalAggregate) Remove(value float64) {
    ia.Count--
    ia.Sum -= value
    // MIN/MAX require full recompute or sorted set
}
```

**Expected Improvement**: 50-70% for incremental updates.

### 4. Token Pooling

**Problem**: Token allocation creates GC pressure.

**Solution**: Reuse token objects via sync.Pool.

```go
var tokenPool = sync.Pool{
    New: func() interface{} {
        return &Token{
            Bindings: make(map[string]*Fact, 4),
        }
    },
}

func acquireToken() *Token {
    return tokenPool.Get().(*Token)
}

func releaseToken(t *Token) {
    // Clear for reuse
    for k := range t.Bindings {
        delete(t.Bindings, k)
    }
    t.Facts = t.Facts[:0]
    tokenPool.Put(t)
}
```

**Expected Improvement**: 10-20% reduction in allocations.

### 5. Cache-Aware Data Structures

**Problem**: Map iterations have poor cache locality.

**Solution**: Use arrays/slices where possible.

```go
// Before: map[string]*Token
CombinedTokens map[string]map[string]*Token

// After: Use slice for better cache locality
type TokenBucket struct {
    MainFactID string
    Tokens     []*Token
}
CombinedTokens []TokenBucket
```

**Expected Improvement**: 5-15% for large datasets.

### 6. Parallel Aggregation

**Problem**: Aggregating multiple main facts is embarrassingly parallel.

**Solution**: Process independent main facts concurrently.

```go
var wg sync.WaitGroup
results := make(chan *AggregationResult, len(mainFacts))

for _, mainFact := range mainFacts {
    wg.Add(1)
    go func(mf *Fact) {
        defer wg.Done()
        result := computeAggregations(mf)
        results <- result
    }(mainFact)
}

wg.Wait()
close(results)
```

**Expected Improvement**: Near-linear scaling with CPU cores (for sufficient main facts).

### 7. Memory Management

**Problem**: Large datasets cause memory pressure and GC pauses.

**Solutions**:

#### a) Eviction Policies

```go
type LRUCache struct {
    maxSize int
    items   map[string]*cacheEntry
    lruList *list.List
}

// Evict least-recently-used when cache is full
```

#### b) Streaming Aggregation

For very large datasets, don't keep all tokens in memory:

```go
// Process in batches
const batchSize = 1000
for batch := range streamTokens(batchSize) {
    processAggregations(batch)
    // Release memory after each batch
}
```

#### c) Fact Compression

Store only essential fields in memory:

```go
type CompressedFact struct {
    ID     string
    Type   string
    Fields []byte // Serialized, compressed
}
```

### 8. Index-Based Deduplication

**Problem**: Checking for duplicate source facts is O(n).

**Solution**: Use bitmap or hash set for O(1) lookups.

```go
type FactSet struct {
    ids map[string]struct{} // Set of fact IDs
}

func (fs *FactSet) Contains(id string) bool {
    _, exists := fs.ids[id]
    return exists
}

func (fs *FactSet) Add(id string) {
    fs.ids[id] = struct{}{}
}
```

---

## Monitoring & Metrics

### Real-Time Metrics

Key metrics to track in production:

```go
type AggregationMetrics struct {
    TotalActivations     int64
    AggregationsComputed int64
    TokensProcessed      int64
    ThresholdsFailed     int64
    AvgComputeTimeMs     float64
    PeakMemoryMB         float64
}
```

### Prometheus Integration

Export metrics for monitoring:

```go
var (
    aggregationsComputed = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "rete_aggregations_computed_total",
            Help: "Total aggregations computed",
        },
    )
    
    aggregationDuration = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "rete_aggregation_duration_seconds",
            Help:    "Time spent computing aggregations",
            Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
        },
    )
)
```

### Logging for Performance

Use structured logging with timing:

```go
start := time.Now()
err := msn.processMainFact(mainFact)
duration := time.Since(start)

if duration > threshold {
    log.Warn("Slow aggregation",
        "main_fact_id", mainFact.ID,
        "duration_ms", duration.Milliseconds(),
        "token_count", len(tokens),
    )
}
```

---

## Troubleshooting

### Issue: Slow Aggregation Performance

**Symptoms:**
- High CPU usage
- Long benchmark times
- Low facts/sec throughput

**Diagnosis:**
1. Run CPU profile: `go tool pprof -http=:8080 cpu.prof`
2. Check for hot paths in `computeAggregation()` or `Activate()`
3. Look for excessive token copying or field access

**Solutions:**
- Implement incremental aggregation (see Optimization #3)
- Add early threshold evaluation (see Optimization #2)
- Optimize join order (see Optimization #1)

### Issue: High Memory Usage

**Symptoms:**
- Large heap allocations
- Frequent GC pauses
- Out-of-memory errors

**Diagnosis:**
1. Run memory profile: `go tool pprof -alloc_space mem.prof`
2. Check `CombinedTokens` and `AggregateCache` sizes
3. Look for leaked references

**Solutions:**
- Implement LRU eviction (see Optimization #7a)
- Use token pooling (see Optimization #4)
- Add fact compression (see Optimization #7c)

### Issue: Poor Scaling with More Sources

**Symptoms:**
- 3-source benchmarks much slower than 2-source
- Non-linear performance degradation

**Diagnosis:**
1. Measure join chain performance separately
2. Check join fanout ratios
3. Profile join node activation

**Solutions:**
- Optimize join order
- Add join result caching
- Consider bloom filters for join pruning

### Issue: Threshold Evaluation Overhead

**Symptoms:**
- WithThresholds benchmark slower than without
- Many failed activations

**Diagnosis:**
1. Count threshold failures vs. successes
2. Profile threshold evaluation code
3. Check for redundant computations

**Solutions:**
- Implement early threshold evaluation
- Use approximate filtering (e.g., skip if partial aggregate already fails)
- Cache threshold results when inputs haven't changed

### Issue: Retraction Performance

**Symptoms:**
- Retraction benchmark much slower than insertion
- Performance degrades over time

**Diagnosis:**
1. Profile retraction codepath
2. Check recomputation frequency
3. Measure cache invalidation overhead

**Solutions:**
- Implement incremental updates instead of full recompute
- Track dependencies more precisely
- Use lazy recomputation (compute on next access)

---

## Advanced Topics

### Custom Aggregation Functions

Extend the system with custom aggregators:

```go
type CustomAggregator interface {
    Init()
    Add(value float64)
    Remove(value float64)
    Result() float64
}

// Example: Median aggregator
type MedianAggregator struct {
    values []float64
}

func (m *MedianAggregator) Result() float64 {
    sort.Float64s(m.values)
    n := len(m.values)
    if n%2 == 0 {
        return (m.values[n/2-1] + m.values[n/2]) / 2
    }
    return m.values[n/2]
}
```

### Distributed Aggregation

For very large datasets, distribute computation:

```go
// Map phase: Compute partial aggregates per node
partials := distributeCompute(mainFacts, func(facts []*Fact) PartialAggregate {
    return computePartialAggregate(facts)
})

// Reduce phase: Combine partials
final := combinePartialAggregates(partials)
```

### Approximate Aggregation

Trade accuracy for speed with sampling:

```go
// Reservoir sampling for approximate COUNT/AVG
type ApproximateAggregator struct {
    reservoir []float64
    size      int
    count     int64
}

func (a *ApproximateAggregator) Add(value float64) {
    a.count++
    if len(a.reservoir) < a.size {
        a.reservoir = append(a.reservoir, value)
    } else {
        // Random replacement
        i := rand.Int63n(a.count)
        if i < int64(a.size) {
            a.reservoir[i] = value
        }
    }
}
```

---

## References

- [Go Profiling Documentation](https://go.dev/blog/pprof)
- [Benchmark Writing Best Practices](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)
- [RETE Algorithm Overview](../README.md)
- [Multi-Source Aggregation Quick Reference](AGGREGATION_JOIN_FEATURE_SUMMARY.md)

---

## Contributing

Found a performance issue or have an optimization idea?

1. Run benchmarks to quantify the improvement
2. Profile before and after
3. Document the change in this guide
4. Submit a PR with benchmarks and profile data

---

**Last Updated:** 2025-01-XX  
**Maintainers:** TSD Contributors