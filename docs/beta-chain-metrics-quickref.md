# BetaChainMetrics Quick Reference

## Table of Contents
- [Creation](#creation)
- [Recording Metrics](#recording-metrics)
- [Querying Metrics](#querying-metrics)
- [Cache Operations](#cache-operations)
- [Analysis Methods](#analysis-methods)
- [Prometheus Export](#prometheus-export)

---

## Creation

```go
// Create new metrics instance
metrics := NewBetaChainMetrics()

// Access from network
metrics := network.GetBetaChainMetrics()
```

---

## Recording Metrics

### Chain Build
```go
metrics.RecordChainBuild(BetaChainMetricDetail{
    RuleID:          "my_rule",
    ChainLength:     3,
    NodesCreated:    1,
    NodesReused:     2,
    BuildTime:       150 * time.Microsecond,
    Timestamp:       time.Now(),
    HashesGenerated: []string{"hash1", "hash2"},
    JoinsExecuted:   5,
    TotalJoinTime:   25 * time.Millisecond,
})
```

### Join Execution
```go
// leftSize=10, rightSize=20, resultSize=50, duration=5ms
metrics.RecordJoinExecution(10, 20, 50, 5*time.Millisecond)
```

### Timing
```go
metrics.AddHashComputeTime(100 * time.Microsecond)
```

---

## Querying Metrics

### Snapshot (Thread-Safe Copy)
```go
snapshot := metrics.GetSnapshot()
fmt.Printf("Chains: %d, Sharing: %.2f%%\n",
    snapshot.TotalChainsBuilt,
    snapshot.SharingRatio * 100)
```

### Summary (Structured Data)
```go
summary := metrics.GetSummary()

// Access sections
chains := summary["chains"].(map[string]interface{})
joins := summary["joins"].(map[string]interface{})
hashCache := summary["hash_cache"].(map[string]interface{})
```

### Direct Field Access
```go
total := metrics.TotalChainsBuilt
avgLength := metrics.AverageChainLength
sharingRatio := metrics.SharingRatio
```

---

## Cache Operations

### Hash Cache
```go
metrics.RecordHashCacheHit()
metrics.RecordHashCacheMiss()
metrics.UpdateHashCacheSize(1000)
efficiency := metrics.GetHashCacheEfficiency() // 0.0 to 1.0
```

### Join Cache
```go
metrics.RecordJoinCacheHit()
metrics.RecordJoinCacheMiss()
metrics.UpdateJoinCacheSize(500)
metrics.RecordJoinCacheEviction()
efficiency := metrics.GetJoinCacheEfficiency() // 0.0 to 1.0
```

### Connection Cache
```go
metrics.RecordConnectionCacheHit()
metrics.RecordConnectionCacheMiss()
efficiency := metrics.GetConnectionCacheEfficiency() // 0.0 to 1.0
```

### Prefix Cache
```go
metrics.RecordPrefixCacheHit()
metrics.RecordPrefixCacheMiss()
metrics.UpdatePrefixCacheSize(200)
efficiency := metrics.GetPrefixCacheEfficiency() // 0.0 to 1.0
```

---

## Analysis Methods

### Top Chains
```go
// By build time (slowest first)
slowest := metrics.GetTopChainsByBuildTime(5)

// By chain length (longest first)
longest := metrics.GetTopChainsByLength(5)

// By join time (most join-intensive first)
joinHeavy := metrics.GetTopChainsByJoinTime(5)

// Iterate results
for _, chain := range slowest {
    fmt.Printf("%s: %v (len=%d)\n",
        chain.RuleID, chain.BuildTime, chain.ChainLength)
}
```

### Performance Stats
```go
// Join performance
joinStats := metrics.GetJoinPerformanceStats()
fmt.Printf("Total: %d, Avg time: %s, Throughput: %.2f joins/sec\n",
    joinStats["total_joins"],
    joinStats["average_time"],
    joinStats["throughput_joins_per_sec"])

// Cache statistics
cacheStats := metrics.GetCacheStats()
hashCache := cacheStats["hash_cache"].(map[string]interface{})
fmt.Printf("Hash cache: %.2f%% efficiency\n",
    hashCache["efficiency_pct"])
```

### Reset
```go
metrics.Reset() // Clear all metrics
```

---

## Prometheus Export

### Setup Exporter
```go
// Alpha only (legacy)
exporter := NewPrometheusExporter(alphaMetrics, config)

// Alpha + Beta
exporter := NewPrometheusExporterWithBeta(alphaMetrics, betaMetrics, config)

// Register metrics
exporter.RegisterMetrics()

// Start HTTP server
go exporter.ServeHTTP(":9090")

// Or manual update
exporter.UpdateMetrics()

// Auto-update every 10s
stopChan := exporter.StartAutoUpdate(10 * time.Second)
defer close(stopChan)
```

### Metric Names

**Chains:**
- `{prefix}_beta_chains_built_total`
- `{prefix}_beta_chains_length_avg`
- `{prefix}_beta_nodes_created_total`
- `{prefix}_beta_nodes_reused_total`
- `{prefix}_beta_nodes_sharing_ratio`

**Joins:**
- `{prefix}_beta_joins_executed_total`
- `{prefix}_beta_joins_time_seconds_avg`
- `{prefix}_beta_joins_selectivity_avg`
- `{prefix}_beta_joins_result_size_avg`

**Hash Cache:**
- `{prefix}_beta_hash_cache_hits_total`
- `{prefix}_beta_hash_cache_misses_total`
- `{prefix}_beta_hash_cache_size`
- `{prefix}_beta_hash_cache_efficiency`

**Join Cache:**
- `{prefix}_beta_join_cache_hits_total`
- `{prefix}_beta_join_cache_misses_total`
- `{prefix}_beta_join_cache_size`
- `{prefix}_beta_join_cache_evictions_total`
- `{prefix}_beta_join_cache_efficiency`

**Connection Cache:**
- `{prefix}_beta_connection_cache_hits_total`
- `{prefix}_beta_connection_cache_misses_total`
- `{prefix}_beta_connection_cache_efficiency`

**Prefix Cache:**
- `{prefix}_beta_prefix_cache_hits_total`
- `{prefix}_beta_prefix_cache_misses_total`
- `{prefix}_beta_prefix_cache_size`
- `{prefix}_beta_prefix_cache_efficiency`

**Timing:**
- `{prefix}_beta_build_time_seconds_total`
- `{prefix}_beta_build_time_seconds_avg`
- `{prefix}_beta_hash_compute_time_seconds_total`

---

## Common Patterns

### Monitoring Loop
```go
ticker := time.NewTicker(1 * time.Minute)
for range ticker.C {
    snapshot := metrics.GetSnapshot()
    log.Printf("Chains: %d, Sharing: %.2f%%, Hash Cache: %.2f%%",
        snapshot.TotalChainsBuilt,
        snapshot.SharingRatio * 100,
        metrics.GetHashCacheEfficiency() * 100)
}
```

### Health Check
```go
func checkMetricsHealth(m *BetaChainMetrics) {
    snapshot := m.GetSnapshot()
    
    if snapshot.SharingRatio < 0.3 {
        log.Warn("Low sharing ratio: %.2f", snapshot.SharingRatio)
    }
    
    if m.GetHashCacheEfficiency() < 0.5 {
        log.Warn("Poor hash cache efficiency")
    }
    
    if snapshot.AverageJoinSelectivity < 0.01 {
        log.Warn("Very low join selectivity")
    }
}
```

### Performance Report
```go
func generateReport(m *BetaChainMetrics) string {
    summary := m.GetSummary()
    
    return fmt.Sprintf(`
Beta Chain Metrics Report
=========================
Chains: %v
Nodes: %v
Joins: %v
Hash Cache: %v
Join Cache: %v
`,
        summary["chains"],
        summary["nodes"],
        summary["joins"],
        summary["hash_cache"],
        summary["join_cache"])
}
```

---

## Key Metrics to Watch

| Metric | Good Range | Action if Outside |
|--------|-----------|-------------------|
| Sharing Ratio | > 0.5 | Enable sharing, check normalization |
| Hash Cache Efficiency | > 0.7 | Increase cache size |
| Join Cache Efficiency | > 0.6 | Tune cache size, check TTL |
| Join Selectivity | 0.01 - 0.3 | Reorder joins, add indexes |
| Avg Build Time | < 1ms | Optimize conditions, check registry |

---

## Thread Safety

✅ **All methods are thread-safe**
- Uses `sync.RWMutex` internally
- Safe for concurrent reads and writes
- `GetSnapshot()` returns copy (no mutex in copy)

---

## Performance

| Operation | Overhead | Notes |
|-----------|---------|-------|
| RecordChainBuild | ~200ns | Mutex + arithmetic |
| RecordJoinExecution | ~300ns | Includes selectivity calc |
| RecordCacheHit/Miss | ~100ns | Mutex + increment |
| GetSnapshot | ~10μs | Deep copy with N chains |
| GetSummary | ~5μs | Constructs map |
| Cache efficiency | ~50ns | Read + division |

---

## License

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License