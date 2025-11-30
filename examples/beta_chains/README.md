# Beta Chains Examples

> **⚠️ WARNING: DEPRECATED EXAMPLES**  
> These examples use an outdated API and are currently non-functional.  
> They need to be updated to use the current `ChainPerformanceConfig` and TSD constraint file format.  
> For working examples, see the test files in `rete/*_test.go`.

This directory contains **executable Go examples** demonstrating the Beta Sharing feature (JoinNode sharing) in TSD's RETE engine.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Examples Overview](#examples-overview)
3. [Running Examples](#running-examples)
4. [Understanding the Output](#understanding-the-output)
5. [Configuration Options](#configuration-options)
6. [Benchmarking](#benchmarking)

---

## Quick Start

```bash
# Run the main interactive example
cd examples/beta_chains
go run main.go

# Run a specific scenario
go run main.go -scenario simple
go run main.go -scenario complex
go run main.go -scenario advanced

# Run with Beta Sharing disabled (comparison)
go run main.go -no-sharing

# Run benchmarks
go test -bench=. -benchmem
```

---

## Examples Overview

### Main Example (`main.go`)

Interactive CLI that lets you:
- Choose between different scenarios
- Compare performance with/without Beta Sharing
- View detailed metrics in real-time
- Export results to JSON/CSV

**Features:**
- ✅ Simple, Complex, and Advanced scenarios
- ✅ Side-by-side comparison
- ✅ Live metrics display
- ✅ ASCII visualizations of chains
- ✅ Performance analysis

### Scenarios

#### 1. Simple Scenario (`scenarios/simple.go`)

**Description:** 2-3 join operations with obvious sharing opportunities.

**Rules:**
- 5 rules sharing the same `Person ⋈ Order` join
- Simple conditions (`amount > 100`, `age > 18`)

**Expected Results:**
- Sharing ratio: ~60-70%
- Memory saved: ~40KB
- Build time: 60-80% faster

**Use Case:** E-commerce recommendations, basic validation

---

#### 2. Complex Scenario (`scenarios/complex.go`)

**Description:** 5+ join operations with optimization opportunities.

**Rules:**
- 10 rules with cascading joins (Person → Order → Item → Review)
- Mixed sharing (some shared, some unique)
- Automatic join order optimization

**Expected Results:**
- Sharing ratio: ~35-45%
- Memory saved: ~120KB
- Build time: 40-50% faster
- Join order optimized for selectivity

**Use Case:** E-commerce with product reviews, order analysis

---

#### 3. Advanced Scenario (`scenarios/advanced.go`)

**Description:** Real-world monitoring system with complex patterns.

**Rules:**
- 20+ rules for anomaly detection
- Diamond patterns (multiple paths converge)
- Prefix reuse optimization
- High cache efficiency

**Expected Results:**
- Sharing ratio: ~40-50%
- Memory saved: ~300KB
- Build time: 50-60% faster
- Cache efficiency: >80%

**Use Case:** Real-time monitoring, alerting systems, fraud detection

---

## Running Examples

### Basic Usage

```bash
# Interactive mode (default)
go run main.go

# Output:
# 🚀 Beta Chains Interactive Example
# ================================
# 
# Select a scenario:
# 1. Simple (2-3 joins, high sharing)
# 2. Complex (5+ joins, mixed sharing)
# 3. Advanced (real-world monitoring)
# 4. Compare all scenarios
# 
# Choice [1-4]: _
```

### Command-line Options

```bash
# Run specific scenario
go run main.go -scenario simple
go run main.go -scenario complex
go run main.go -scenario advanced

# Disable Beta Sharing (comparison baseline)
go run main.go -scenario simple -no-sharing

# Enable verbose logging
go run main.go -scenario complex -verbose

# Export results
go run main.go -scenario advanced -export results.json
go run main.go -scenario advanced -export results.csv

# Run all scenarios and compare
go run main.go -compare
```

### Examples with Configuration

```bash
# High performance config (large caches)
go run main.go -config high-performance

# Memory optimized config (small caches)
go run main.go -config memory-optimized

# Custom cache size
go run main.go -cache-size 5000
```

---

## Understanding the Output

### Example Output

```
🚀 Beta Chains Example: Simple Scenario
========================================

📊 Configuration:
  Beta Sharing:     ✅ Enabled
  Join Cache Size:  1000
  Hash Cache Size:  500
  Metrics:          ✅ Enabled

🏗️  Building network with 5 rules...
  [1/5] rule_high_spender     ✅ (78µs)
  [2/5] rule_vip_customer     ✅ (45µs) - Reused 1 node
  [3/5] rule_discount         ✅ (42µs) - Reused 1 node
  [4/5] rule_loyalty_points   ✅ (41µs) - Reused 1 node
  [5/5] rule_send_email       ✅ (40µs) - Reused 1 node

✅ Network built in 246µs

📈 Beta Sharing Metrics:
┌─────────────────────────────────────────────┐
│ Nodes                                       │
├─────────────────────────────────────────────┤
│ Created:              1                     │
│ Reused:               4                     │
│ Sharing Ratio:        80.0%                 │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│ Performance                                 │
├─────────────────────────────────────────────┤
│ Build Time:           246µs                 │
│ Avg Time/Rule:        49µs                  │
│ Memory Saved:         ~32KB                 │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│ Cache                                       │
├─────────────────────────────────────────────┤
│ Join Cache Hits:      0                     │
│ Join Cache Misses:    0                     │
│ Cache Efficiency:     N/A (no runtime yet)  │
└─────────────────────────────────────────────┘

🔍 Chain Visualization:

TypeNode(Person)    TypeNode(Order)
       └─────────────────┴─────────┐
                              JoinNode(beta_shared)
                              [p.id == o.personId AND o.amount > 100]
                              RefCount: 5 ⭐
                              ├────┬────┬────┬────┬────┐
                         T1   T2   T3   T4   T5
                         (high_spender, vip_customer, ...)

💡 Analysis:
  ✅ Excellent sharing (80%)
  ✅ 4 JoinNodes saved
  ✅ 4x faster build time vs no sharing
  ✅ Significant memory savings
```

### Comparison Output

```
📊 Comparison: WITH vs WITHOUT Beta Sharing
============================================

                    WITH Sharing    WITHOUT Sharing    GAIN
──────────────────────────────────────────────────────────
Nodes Created            1                5           80%
Build Time             246µs            982µs          75%
Memory Used            8KB              40KB           80%
Sharing Ratio          80.0%            0.0%           N/A

💰 Total Savings:
  - 4 JoinNodes saved
  - 736µs faster (75% improvement)
  - 32KB memory saved (80% reduction)
```

---

## Configuration Options

### Predefined Configurations

```go
// Default (balanced)
config := rete.DefaultBetaChainConfig()

// High performance (large caches)
config := rete.HighPerformanceBetaChainConfig()

// Memory optimized (small caches)
config := rete.MemoryOptimizedBetaChainConfig()
```

### Custom Configuration

```go
config := rete.BetaChainConfig{
    EnableBetaSharing:   true,
    EnableMetrics:       true,
    JoinCacheSize:       2000,
    HashCacheSize:       1000,
    EnableOptimization:  true,
    EnablePrefixSharing: true,
}
```

### Configuration via CLI

```bash
# Use predefined config
go run main.go -config high-performance

# Custom cache sizes
go run main.go -join-cache 5000 -hash-cache 2000

# Disable optimizations
go run main.go -no-optimization -no-prefix-sharing
```

---

## Benchmarking

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. -benchmem

# Run specific benchmark
go test -bench=BenchmarkSimple -benchmem
go test -bench=BenchmarkComplex -benchmem
go test -bench=BenchmarkAdvanced -benchmem

# Compare with/without sharing
go test -bench=BenchmarkSimpleWithSharing -benchmem
go test -bench=BenchmarkSimpleWithoutSharing -benchmem

# Generate CPU profile
go test -bench=. -cpuprofile=cpu.prof
go tool pprof cpu.prof

# Generate memory profile
go test -bench=. -memprofile=mem.prof
go tool pprof mem.prof
```

### Benchmark Results (example)

```
BenchmarkSimpleWithSharing-8        10000    112456 ns/op    8192 B/op    5 allocs/op
BenchmarkSimpleWithoutSharing-8      3000    478923 ns/op   40960 B/op   25 allocs/op

BenchmarkComplexWithSharing-8        5000    245678 ns/op   24576 B/op   15 allocs/op
BenchmarkComplexWithoutSharing-8     1000    891234 ns/op  122880 B/op   75 allocs/op

BenchmarkAdvancedWithSharing-8       2000    567890 ns/op   65536 B/op   40 allocs/op
BenchmarkAdvancedWithoutSharing-8     500   1987654 ns/op  327680 B/op  200 allocs/op
```

**Analysis:**
- Simple: 4.3x faster, 5x less memory
- Complex: 3.6x faster, 5x less memory
- Advanced: 3.5x faster, 5x less memory

---

## Metrics and Analysis

### Available Metrics

```go
metrics := network.GetBetaChainMetrics()
snapshot := metrics.GetSnapshot()

// Nodes
snapshot.TotalNodesCreated    // Number of unique JoinNodes
snapshot.TotalNodesReused     // Number of reuses
snapshot.SharingRatio         // Reused / (Created + Reused)

// Cache
snapshot.JoinCacheHits        // Cache hits
snapshot.JoinCacheMisses      // Cache misses
snapshot.JoinCacheSize        // Current cache size

// Performance
snapshot.TotalChainsBuilt     // Number of rules
snapshot.AverageBuildTime     // Avg build time per rule
```

### Exporting Metrics

```bash
# JSON format
go run main.go -scenario advanced -export results.json

# CSV format
go run main.go -scenario advanced -export results.csv

# Both
go run main.go -compare -export comparison.json
```

**JSON Output Example:**
```json
{
  "scenario": "advanced",
  "config": {
    "enable_beta_sharing": true,
    "join_cache_size": 1000,
    "hash_cache_size": 500
  },
  "metrics": {
    "nodes_created": 12,
    "nodes_reused": 8,
    "sharing_ratio": 0.4,
    "build_time_us": 1234,
    "memory_saved_kb": 64,
    "cache_efficiency": 0.85
  },
  "comparison": {
    "without_sharing": {
      "nodes_created": 20,
      "build_time_us": 2987
    },
    "gains": {
      "nodes_saved": 8,
      "time_saved_pct": 58.7,
      "memory_saved_pct": 60.0
    }
  }
}
```

---

## Troubleshooting

### Low Sharing Ratio

**Problem:** Sharing ratio < 10%

**Causes:**
- Rules have very different conditions
- Different variable types
- No common join patterns

**Solution:**
```bash
# Check with verbose logging
go run main.go -scenario simple -verbose

# Analyze the rules
# Look for common patterns that should share
```

### Cache Inefficiency

**Problem:** Cache efficiency < 50%

**Causes:**
- Cache too small for workload
- Unique evaluations (no repetition)

**Solution:**
```bash
# Increase cache size
go run main.go -scenario complex -cache-size 5000

# Monitor efficiency
# If still low, cache might not be beneficial for this workload
```

### High Memory Usage

**Problem:** Memory usage higher than expected

**Causes:**
- Cache too large
- Too many rules loaded

**Solution:**
```bash
# Use memory-optimized config
go run main.go -config memory-optimized

# Or disable metrics
go run main.go -no-metrics
```

---

## Next Steps

1. **Experiment with scenarios:** Try all three scenarios to understand different patterns
2. **Run benchmarks:** Compare performance with your own rules
3. **Customize:** Modify scenarios to match your use cases
4. **Integrate:** Use the learnings in your production code

## Resources

- [BETA_CHAINS_EXAMPLES.md](../../rete/BETA_CHAINS_EXAMPLES.md) - Detailed examples
- [BETA_CHAINS_MIGRATION.md](../../rete/BETA_CHAINS_MIGRATION.md) - Migration guide
- [BETA_CHAINS_USER_GUIDE.md](../../rete/BETA_CHAINS_USER_GUIDE.md) - User guide
- [BETA_CHAINS_TECHNICAL_GUIDE.md](../../rete/BETA_CHAINS_TECHNICAL_GUIDE.md) - Technical details

---

## License

Copyright (c) 2025 TSD Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS