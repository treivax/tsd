# BetaChainMetrics Documentation

## Overview

`BetaChainMetrics` is a comprehensive metrics collection system for Beta node (JoinNode) chain construction in the RETE network. It provides detailed insights into chain building performance, node sharing efficiency, join execution statistics, and cache performance.

## Features

- **Chain Construction Metrics**: Track chains built, nodes created/reused, and sharing ratios
- **Join Performance Metrics**: Monitor join execution time, selectivity, and result sizes
- **Multi-Level Cache Metrics**: Hash cache, join cache, connection cache, and prefix cache statistics
- **Thread-Safe Operations**: All operations protected by `sync.RWMutex`
- **Prometheus Integration**: Export metrics to Prometheus for monitoring
- **Detailed Chain Analysis**: Per-chain metrics with sorting and filtering capabilities

## Architecture

```
BetaChainMetrics
├── Chain Metrics (overall statistics)
├── Node Sharing Metrics (created vs reused)
├── Join Performance Metrics (execution stats)
├── Hash Cache Metrics (computation caching)
├── Join Cache Metrics (result caching)
├── Connection Cache Metrics (topology caching)
└── Prefix Cache Metrics (chain reuse)
```

## API Reference

### Constructor

```go
func NewBetaChainMetrics() *BetaChainMetrics
```

Creates a new metrics instance with all counters initialized to zero.

### Recording Methods

#### RecordChainBuild

```go
func (m *BetaChainMetrics) RecordChainBuild(detail BetaChainMetricDetail)
```

Records metrics for a completed chain build. Automatically updates:
- Total chains built
- Nodes created/reused counts
- Average chain length
- Sharing ratio
- Build time statistics

**Example:**

```go
detail := BetaChainMetricDetail{
    RuleID:          "customer_order_join",
    ChainLength:     3,
    NodesCreated:    1,
    NodesReused:     2,
    BuildTime:       150 * time.Microsecond,
    Timestamp:       time.Now(),
    HashesGenerated: []string{"join_abc123", "join_def456"},
    JoinsExecuted:   10,
    TotalJoinTime:   50 * time.Millisecond,
}
metrics.RecordChainBuild(detail)
```

#### RecordJoinExecution

```go
func (m *BetaChainMetrics) RecordJoinExecution(leftSize, rightSize, resultSize int, duration time.Duration)
```

Records a single join operation with automatic selectivity calculation.

**Selectivity Formula:** `resultSize / (leftSize * rightSize)`

**Example:**

```go
// 10 left tokens × 20 right facts = 200 possible results
// Actual: 50 results → selectivity = 0.25 (25%)
metrics.RecordJoinExecution(10, 20, 50, 5*time.Millisecond)
```

#### Cache Metrics

```go
// Hash cache (condition hashing)
func (m *BetaChainMetrics) RecordHashCacheHit()
func (m *BetaChainMetrics) RecordHashCacheMiss()
func (m *BetaChainMetrics) UpdateHashCacheSize(size int)

// Join cache (result caching)
func (m *BetaChainMetrics) RecordJoinCacheHit()
func (m *BetaChainMetrics) RecordJoinCacheMiss()
func (m *BetaChainMetrics) UpdateJoinCacheSize(size int)
func (m *BetaChainMetrics) RecordJoinCacheEviction()

// Connection cache (topology caching)
func (m *BetaChainMetrics) RecordConnectionCacheHit()
func (m *BetaChainMetrics) RecordConnectionCacheMiss()

// Prefix cache (chain reuse)
func (m *BetaChainMetrics) RecordPrefixCacheHit()
func (m *BetaChainMetrics) RecordPrefixCacheMiss()
func (m *BetaChainMetrics) UpdatePrefixCacheSize(size int)
```

### Query Methods

#### GetSnapshot

```go
func (m *BetaChainMetrics) GetSnapshot() BetaChainMetrics
```

Returns a thread-safe copy of all metrics. **Note:** The returned copy does not include the mutex.

**Example:**

```go
snapshot := metrics.GetSnapshot()
fmt.Printf("Total chains: %d\n", snapshot.TotalChainsBuilt)
fmt.Printf("Sharing ratio: %.2f%%\n", snapshot.SharingRatio * 100)
```

#### GetSummary

```go
func (m *BetaChainMetrics) GetSummary() map[string]interface{}
```

Returns a structured summary with nested sections for chains, nodes, joins, and caches.

**Example:**

```go
summary := metrics.GetSummary()
chains := summary["chains"].(map[string]interface{})
fmt.Printf("Chains built: %d\n", chains["total_built"])
fmt.Printf("Avg build time: %s\n", chains["average_build_time"])

joins := summary["joins"].(map[string]interface{})
fmt.Printf("Total joins: %d\n", joins["total_executed"])
fmt.Printf("Avg selectivity: %.4f\n", joins["average_selectivity"])
```

#### Efficiency Methods

```go
func (m *BetaChainMetrics) GetHashCacheEfficiency() float64
func (m *BetaChainMetrics) GetJoinCacheEfficiency() float64
func (m *BetaChainMetrics) GetConnectionCacheEfficiency() float64
func (m *BetaChainMetrics) GetPrefixCacheEfficiency() float64
```

Return cache hit rates as values between 0.0 and 1.0.

**Example:**

```go
hashEff := metrics.GetHashCacheEfficiency()
if hashEff < 0.5 {
    log.Printf("Warning: Hash cache efficiency low: %.2f%%", hashEff*100)
}
```

#### Top Chains Analysis

```go
func (m *BetaChainMetrics) GetTopChainsByBuildTime(n int) []BetaChainMetricDetail
func (m *BetaChainMetrics) GetTopChainsByLength(n int) []BetaChainMetricDetail
func (m *BetaChainMetrics) GetTopChainsByJoinTime(n int) []BetaChainMetricDetail
```

Return the top N chains sorted by the specified metric (descending).

**Example:**

```go
// Find the 5 slowest chains to build
slowest := metrics.GetTopChainsByBuildTime(5)
for i, chain := range slowest {
    fmt.Printf("%d. Rule %s: %v (length: %d)\n",
        i+1, chain.RuleID, chain.BuildTime, chain.ChainLength)
}
```

#### Specialized Statistics

```go
func (m *BetaChainMetrics) GetJoinPerformanceStats() map[string]interface{}
func (m *BetaChainMetrics) GetCacheStats() map[string]interface{}
```

**GetJoinPerformanceStats** returns:
- Total joins executed
- Total/average time
- Average selectivity and result size
- Throughput (joins/sec)

**GetCacheStats** returns efficiency for all four cache types.

### Maintenance Methods

#### Reset

```go
func (m *BetaChainMetrics) Reset()
```

Resets all metrics to zero. Useful for testing or starting a new measurement period.

**Example:**

```go
// Start fresh measurement
metrics.Reset()

// Build some chains...
builder.BuildChain(patterns, "rule1")
builder.BuildChain(patterns, "rule2")

// Analyze this period's metrics
summary := metrics.GetSummary()
```

## Metric Fields

### Core Metrics

| Field | Type | Description |
|-------|------|-------------|
| `TotalChainsBuilt` | `int` | Number of beta chains built |
| `TotalNodesCreated` | `int` | Unique JoinNodes created |
| `TotalNodesReused` | `int` | JoinNodes reused from registry |
| `AverageChainLength` | `float64` | Average nodes per chain |
| `SharingRatio` | `float64` | Ratio of reused to total nodes (0.0-1.0) |

### Sharing Metrics (Compatibility)

| Field | Type | Description |
|-------|------|-------------|
| `TotalJoinNodesRequested` | `int64` | Total node creation/lookup requests |
| `SharedJoinNodesReused` | `int64` | Requests satisfied by existing nodes |
| `UniqueJoinNodesCreated` | `int64` | New nodes created |

### Join Execution Metrics

| Field | Type | Description |
|-------|------|-------------|
| `TotalJoinsExecuted` | `int` | Total join operations performed |
| `AverageJoinTime` | `time.Duration` | Average time per join |
| `TotalJoinTime` | `time.Duration` | Cumulative join time |
| `AverageJoinSelectivity` | `float64` | Average result/input ratio (0.0-1.0) |
| `AverageResultSize` | `float64` | Average number of results per join |

### Cache Metrics

#### Hash Cache (Condition Normalization)

| Field | Type | Description |
|-------|------|-------------|
| `HashCacheHits` | `int` | Cached hash lookups |
| `HashCacheMisses` | `int` | Hash computations performed |
| `HashCacheSize` | `int` | Current cache entries |

#### Join Cache (Result Memoization)

| Field | Type | Description |
|-------|------|-------------|
| `JoinCacheHits` | `int` | Cached result reuses |
| `JoinCacheMisses` | `int` | Results computed |
| `JoinCacheSize` | `int` | Current cached results |
| `JoinCacheEvictions` | `int` | Cache entries evicted |

#### Connection Cache (Topology)

| Field | Type | Description |
|-------|------|-------------|
| `ConnectionCacheHits` | `int` | Topology lookups satisfied |
| `ConnectionCacheMisses` | `int` | New connections checked |

#### Prefix Cache (Chain Reuse)

| Field | Type | Description |
|-------|------|-------------|
| `PrefixCacheHits` | `int` | Reusable prefixes found |
| `PrefixCacheMisses` | `int` | Prefix lookups failed |
| `PrefixCacheSize` | `int` | Current prefix entries |

### Timing Metrics

| Field | Type | Description |
|-------|------|-------------|
| `TotalBuildTime` | `time.Duration` | Cumulative chain build time |
| `AverageBuildTime` | `time.Duration` | Average build time per chain |
| `TotalHashComputeTime` | `time.Duration` | Time spent computing hashes |

## Integration with ReteNetwork

### Accessing Metrics

```go
// Get metrics from network
network := NewReteNetwork(storage)
metrics := network.GetBetaChainMetrics()

if metrics != nil {
    summary := metrics.GetSummary()
    // Process summary...
}
```

### Reset All Metrics

```go
// Reset both alpha and beta metrics
network.ResetChainMetrics()
```

## Prometheus Integration

### Setup

```go
alphaMetrics := network.GetChainMetrics()
betaMetrics := network.GetBetaChainMetrics()
config := network.GetConfig()

// Create exporter with both alpha and beta metrics
exporter := NewPrometheusExporterWithBeta(alphaMetrics, betaMetrics, config)
exporter.RegisterMetrics()

// Start HTTP server
go exporter.ServeHTTP(":9090")
```

### Exported Metrics

All beta metrics are prefixed with `{prefix}_beta_` where prefix is configurable (default: `rete`).

#### Chain Metrics

- `beta_chains_built_total` (counter): Total chains built
- `beta_chains_length_avg` (gauge): Average chain length
- `beta_nodes_created_total` (counter): Total unique nodes
- `beta_nodes_reused_total` (counter): Total reused nodes
- `beta_nodes_sharing_ratio` (gauge): Sharing efficiency (0.0-1.0)

#### Join Metrics

- `beta_joins_executed_total` (counter): Total joins performed
- `beta_joins_time_seconds_avg` (gauge): Average join time
- `beta_joins_selectivity_avg` (gauge): Average selectivity (0.0-1.0)
- `beta_joins_result_size_avg` (gauge): Average result size

#### Cache Metrics

All cache types expose:
- `_hits_total` (counter): Cache hits
- `_misses_total` (counter): Cache misses
- `_efficiency` (gauge): Hit rate (0.0-1.0)
- `_size` (gauge): Current entries (where applicable)

Additional for join cache:
- `beta_join_cache_evictions_total` (counter): Evictions

#### Timing Metrics

- `beta_build_time_seconds_total` (counter): Cumulative build time
- `beta_build_time_seconds_avg` (gauge): Average build time
- `beta_hash_compute_time_seconds_total` (counter): Total hash computation time

### Example Prometheus Queries

```promql
# Sharing efficiency over time
rete_beta_nodes_sharing_ratio

# Join throughput (joins per second)
rate(rete_beta_joins_executed_total[5m])

# Average join latency
rete_beta_joins_time_seconds_avg

# Cache hit rates
rete_beta_hash_cache_efficiency
rete_beta_join_cache_efficiency

# Chain build rate
rate(rete_beta_chains_built_total[1m])

# Join selectivity trend
avg_over_time(rete_beta_joins_selectivity_avg[10m])
```

## Best Practices

### 1. Monitor Sharing Ratio

```go
metrics := network.GetBetaChainMetrics()
sharingRatio := metrics.GetSnapshot().SharingRatio

if sharingRatio < 0.3 {
    log.Warn("Low node sharing ratio: %.2f - consider enabling sharing", sharingRatio)
} else if sharingRatio > 0.7 {
    log.Info("Excellent node sharing: %.2f", sharingRatio)
}
```

### 2. Track Cache Performance

```go
hashEff := metrics.GetHashCacheEfficiency()
joinEff := metrics.GetJoinCacheEfficiency()

if hashEff < 0.5 {
    log.Warn("Hash cache inefficient - consider increasing cache size")
}

if joinEff > 0.8 {
    log.Info("Join cache performing well: %.2f hit rate", joinEff)
}
```

### 3. Identify Slow Chains

```go
slowest := metrics.GetTopChainsByBuildTime(10)
for _, chain := range slowest {
    if chain.BuildTime > 100*time.Millisecond {
        log.Warn("Slow chain build: %s took %v", chain.RuleID, chain.BuildTime)
    }
}
```

### 4. Analyze Join Selectivity

```go
stats := metrics.GetJoinPerformanceStats()
selectivity := stats["average_selectivity"].(float64)

if selectivity < 0.01 {
    log.Warn("Very low join selectivity: %.4f - consider join reordering", selectivity)
}
```

### 5. Periodic Reports

```go
ticker := time.NewTicker(5 * time.Minute)
go func() {
    for range ticker.C {
        snapshot := metrics.GetSnapshot()
        log.Info("Beta Metrics Report:")
        log.Info("  Chains: %d built, avg length: %.2f",
            snapshot.TotalChainsBuilt, snapshot.AverageChainLength)
        log.Info("  Sharing: %.2f%% (%d reused / %d total)",
            snapshot.SharingRatio*100, snapshot.TotalNodesReused,
            snapshot.TotalNodesCreated+snapshot.TotalNodesReused)
        log.Info("  Joins: %d executed, avg time: %v",
            snapshot.TotalJoinsExecuted, snapshot.AverageJoinTime)
        log.Info("  Hash cache: %.2f%% efficiency",
            metrics.GetHashCacheEfficiency()*100)
    }
}()
```

## Performance Considerations

### Thread Safety

All public methods are thread-safe using `sync.RWMutex`:
- Read operations use `RLock()`
- Write operations use `Lock()`

### Memory Usage

- `ChainDetails` slice grows with each chain built
- Consider periodic resets for long-running systems
- Snapshot copies all data → use sparingly in hot paths

### Metrics Overhead

Typical overhead per operation:
- `RecordChainBuild`: ~200ns (mutex + arithmetic)
- `RecordJoinExecution`: ~300ns (includes selectivity calculation)
- `RecordCacheHit/Miss`: ~100ns (mutex + increment)
- `GetSnapshot`: ~10μs (deep copy with N chain details)

## Troubleshooting

### Metrics Always Zero

```go
// Check if metrics are being recorded
metrics := network.GetBetaChainMetrics()
if metrics == nil {
    log.Error("BetaChainBuilder not initialized")
    // Enable beta sharing in config
    config.BetaSharingEnabled = true
}
```

### Sharing Ratio Zero Despite Reuse

Ensure `RecordChainBuild` is called with correct `NodesCreated` and `NodesReused` counts:

```go
detail := BetaChainMetricDetail{
    NodesCreated: uniqueNodes,  // Only new nodes
    NodesReused:  sharedNodes,  // Only reused nodes
    // ...
}
metrics.RecordChainBuild(detail)
```

### Prometheus Metrics Not Updating

```go
// Metrics must be explicitly updated
exporter.UpdateMetrics()

// Or enable auto-update
stopChan := exporter.StartAutoUpdate(10 * time.Second)
defer close(stopChan)
```

## License

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License