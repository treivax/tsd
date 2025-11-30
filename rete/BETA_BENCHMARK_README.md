# Beta Sharing Benchmarks - Quick Start Guide

**License:** MIT  
**Version:** 1.0  
**Date:** 2025-01-XX

## Overview

This document provides a quick start guide for running and interpreting Beta Sharing System performance benchmarks.

## Quick Start

### Run All Beta Benchmarks

```bash
cd rete
go test -bench=BenchmarkBeta -benchmem -benchtime=1s
```

### Run Specific Benchmark Categories

```bash
# Chain construction benchmarks
go test -bench=BenchmarkBetaChainBuild -benchmem -benchtime=1s

# Join cache benchmarks
go test -bench=BenchmarkJoinCache -benchmem -benchtime=1s

# Hash computation benchmarks
go test -bench=BenchmarkHashCompute -benchmem -benchtime=1s

# High load scenarios
go test -bench=BenchmarkBetaChainBuild_HighLoad -benchmem -benchtime=1s

# Prefix sharing benchmarks
go test -bench=BenchmarkPrefixSharing -benchmem -benchtime=1s
```

## Understanding Benchmark Output

### Example Output

```
BenchmarkBetaChainBuild_WithSharing-16
    39134 ops    15808 ns/op    10.00 nodes_created    391330 nodes_reused    100.0 sharing_%    5668 B/op    105 allocs/op

BenchmarkBetaChainBuild_WithoutSharing-16
    22384 ops    28719 ns/op    22393 nodes_created    6566 B/op    121 allocs/op
```

### Metrics Explained

| Metric | Description | Good Value |
|--------|-------------|------------|
| `ops` | Operations completed in benchmark time | Higher is better |
| `ns/op` | Nanoseconds per operation | Lower is better |
| `sharing_%` | Percentage of nodes reused | 60%+ is good, 80%+ is excellent |
| `nodes_created` | New nodes created | Lower with sharing |
| `nodes_reused` | Existing nodes reused | Higher is better |
| `B/op` | Bytes allocated per operation | Lower is better |
| `allocs/op` | Allocations per operation | Lower is better |

### Interpreting Results

**With Sharing:**
- ✅ 45% faster (28719ns → 15808ns)
- ✅ 13.7% fewer allocations (121 → 105)
- ✅ 100% sharing ratio (perfect reuse)

## Benchmark Categories

### 1. Chain Construction Performance

**Purpose:** Compare building performance with/without sharing

```bash
go test -bench="BenchmarkBetaChainBuild_(WithSharing|WithoutSharing)" -benchmem -benchtime=1s
```

**What to look for:**
- Speed improvement with sharing (expect 15-40% faster)
- Memory reduction (expect 30-50% fewer allocations)
- Sharing ratio (expect 60-90% for similar patterns)

### 2. Similar vs Mixed Patterns

**Purpose:** Understand how pattern similarity affects sharing

```bash
go test -bench="BenchmarkBetaChainBuild_(Similar|Mixed)Patterns" -benchmem -benchtime=1s
```

**What to look for:**
- Similar patterns: 70%+ sharing ratio
- Mixed patterns: 30-50% sharing ratio
- Performance scales consistently

### 3. Join Cache Performance

**Purpose:** Measure join result caching efficiency

```bash
go test -bench=BenchmarkJoinCache -benchmem -benchtime=1s
```

**What to look for:**
- Cache hits: ~125ns (very fast)
- Cache misses: ~234ns (acceptable)
- Hit rate in mixed workload: 60-70%

### 4. Hash Computation

**Purpose:** Measure hashing overhead and cache effectiveness

```bash
go test -bench=BenchmarkHashCompute -benchmem -benchtime=1s
```

**What to look for:**
- Simple conditions: ~1.2μs
- Complex conditions: ~2.3μs
- With cache: ~289ns (4x faster)
- Cache hit rate: 90%+ in steady state

### 5. Scalability Tests

**Purpose:** Verify performance at different scales

```bash
# 10 rules
go test -bench="BenchmarkBetaChainBuild_(Similar|Mixed)Patterns_10Rules" -benchmem

# 100 rules
go test -bench="BenchmarkBetaChainBuild_(Similar|Mixed)Patterns_100Rules" -benchmem
```

**What to look for:**
- Linear scaling (no exponential degradation)
- Consistent sharing ratios across scales
- Memory growth proportional to rule count

### 6. Complex Rules

**Purpose:** Test with complex multi-join rules

```bash
go test -bench=BenchmarkBetaChainBuild_ComplexRules -benchmem -benchtime=1s
```

**What to look for:**
- Acceptable performance (< 500μs per chain)
- Sharing ratio remains strong (50%+)
- Hash cache effectiveness (85%+)

### 7. High Load Scenarios

**Purpose:** Stress test with large deployments

```bash
go test -bench=BenchmarkBetaChainBuild_HighLoad -benchmem -benchtime=1s
```

**What to look for:**
- Performance remains stable
- Higher sharing ratios with more rules
- No memory leaks or degradation

### 8. Prefix Sharing

**Purpose:** Measure prefix sharing optimization impact

```bash
go test -bench=BenchmarkPrefixSharing -benchmem -benchtime=1s
```

**What to look for:**
- 15+ percentage point increase in sharing ratio
- Small performance overhead (< 5%)
- Good prefix cache hit rate (60%+)

## Advanced Usage

### Generate CPU Profile

```bash
go test -bench=BenchmarkBetaChainBuild_WithSharing \
    -benchtime=5s \
    -cpuprofile=cpu.prof \
    -memprofile=mem.prof

# View profile
go tool pprof -http=:8080 cpu.prof
```

### Compare Before/After Changes

```bash
# Baseline
go test -bench=BenchmarkBeta -benchmem > baseline.txt

# After changes
go test -bench=BenchmarkBeta -benchmem > optimized.txt

# Compare (requires benchstat: go install golang.org/x/perf/cmd/benchstat@latest)
benchstat baseline.txt optimized.txt
```

### Run Specific Pattern Count

```bash
# The benchmark functions accept pattern count internally
# To customize, edit the benchmark or use subtests
go test -bench=BenchmarkBetaChainBuild_SimilarPatterns_10Rules -benchmem
```

### Benchmark with Race Detector

```bash
# Check for concurrency issues
go test -race -bench=BenchmarkBeta -benchtime=100ms
```

## Performance Targets

### Baseline Targets

| Scenario | Target | Metric |
|----------|--------|--------|
| Simple chain build (10 patterns) | < 20μs | Time per operation |
| Complex chain build (15+ patterns) | < 500μs | Time per operation |
| Sharing ratio (similar patterns) | > 70% | Reuse percentage |
| Sharing ratio (mixed patterns) | > 30% | Reuse percentage |
| Hash cache hit rate | > 90% | Steady state |
| Join cache hit rate | > 65% | 70% read workload |
| Memory reduction | > 30% | vs without sharing |

### Production Targets

| Deployment Size | Build Time Budget | Sharing Ratio Target |
|-----------------|-------------------|---------------------|
| Small (< 50 rules) | < 1ms per chain | 60%+ |
| Medium (50-500 rules) | < 2ms per chain | 70%+ |
| Large (> 500 rules) | < 5ms per chain | 75%+ |

## Troubleshooting

### Low Sharing Ratio

**Symptoms:** Sharing ratio < 50% with similar patterns

**Possible causes:**
- Patterns are more diverse than expected
- Variable names differ (normalization may help)
- Condition structure varies

**Solutions:**
1. Enable advanced normalization: `config.EnableAdvancedNormalization = true`
2. Review pattern definitions for consistency
3. Check if prefix sharing is enabled

### High Memory Usage

**Symptoms:** Allocations/op higher than expected

**Possible causes:**
- Hash cache too small (frequent evictions)
- Join cache oversized for workload
- Memory leaks in custom conditions

**Solutions:**
1. Increase hash cache size: `config.HashCacheSize = 2000`
2. Run with `-memprofile` to identify allocations
3. Check for goroutine leaks

### Slow Hash Computation

**Symptoms:** Hash operations take > 5μs

**Possible causes:**
- Cache disabled or too small
- Complex nested conditions
- JSON serialization overhead

**Solutions:**
1. Enable hash cache: `config.HashCacheSize = 1000`
2. Simplify condition structures if possible
3. Use simpler operators where equivalent

## Benchmark Development

### Adding New Benchmarks

```go
func BenchmarkBetaChainBuild_YourScenario(b *testing.B) {
    // Setup
    storage := NewMemoryStorage()
    network := NewReteNetwork(storage)
    
    config := BetaSharingConfig{
        Enabled:       true,
        HashCacheSize: 1000,
    }
    betaRegistry := NewBetaSharingRegistry(config, network.LifecycleManager)
    builder := NewBetaChainBuilderWithRegistry(network, storage, betaRegistry)
    
    patterns := createYourPatterns()
    
    b.ResetTimer()
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        ruleID := fmt.Sprintf("rule_%d", i)
        _, err := builder.BuildChain(patterns, ruleID)
        if err != nil {
            b.Fatalf("BuildChain failed: %v", err)
        }
    }
    
    b.StopTimer()
    reportBenchmarkMetrics(b, builder)
}
```

### Benchmark Best Practices

1. **Reset timer after setup:** Use `b.ResetTimer()` after initialization
2. **Report allocations:** Always use `b.ReportAllocs()`
3. **Stop timer for cleanup:** Use `b.StopTimer()` before teardown
4. **Report custom metrics:** Use `b.ReportMetric()` for domain-specific data
5. **Use realistic data:** Pattern complexity should match production

## Related Documentation

- [Beta Performance Report](docs/BETA_PERFORMANCE_REPORT.md) - Detailed analysis and recommendations
- [Beta Sharing README](BETA_SHARING_README.md) - System overview and configuration
- [Beta Chain Builder](BETA_CHAIN_BUILDER_README.md) - Builder API documentation
- [Beta Chain Metrics](docs/BETA_CHAIN_METRICS_README.md) - Metrics API reference

## FAQ

### Q: Why are there debug logs during benchmarks?

A: The system uses structured logging. You can filter output:
```bash
go test -bench=Benchmark 2>&1 | grep "^Benchmark"
```

### Q: How long should I run benchmarks?

A: For quick checks: `-benchtime=100ms`  
For stable results: `-benchtime=1s` or `-benchtime=5s`

### Q: What's a good sharing ratio?

A: 
- Similar patterns: 70-90%
- Somewhat similar: 40-70%
- Diverse patterns: 20-40%

### Q: Should I always enable sharing?

A: Yes for production with 10+ rules. For < 10 simple rules, the overhead may outweigh benefits.

### Q: How do I benchmark my actual rules?

A: Create a benchmark file that loads your rules:
```go
func BenchmarkProductionRules(b *testing.B) {
    rules := loadProductionRules() // Your loader
    // ... setup and benchmark
}
```

## License

This benchmark suite is part of TSD and is licensed under the MIT License.

---

**Last Updated:** 2025-01-XX  
**Maintainer:** TSD Contributors