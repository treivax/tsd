# Beta Sharing Quick Reference

**TSD RETE Engine - Beta Node Sharing System**

---

## Quick Start

### Enable Beta Sharing

```go
// Option 1: Use High Performance preset
config := rete.HighPerformanceConfig()
network := rete.NewReteNetworkWithConfig(storage, config)

// Option 2: Enable in default config
config := rete.DefaultChainPerformanceConfig()
config.BetaSharingEnabled = true
network := rete.NewReteNetworkWithConfig(storage, config)
```

### Check Status

```go
// Check if enabled
if network.BetaSharingRegistry != nil {
    fmt.Println("Beta sharing is enabled")
}

// Get sharing statistics
stats := network.GetBetaSharingStats()
if stats != nil {
    fmt.Printf("Sharing ratio: %.1f%%\n", stats.SharingRatio * 100)
    fmt.Printf("Nodes shared: %d\n", stats.TotalSharedNodes)
}
```

---

## Configuration

### Basic Settings

```go
config.BetaSharingEnabled = true          // Enable/disable sharing
config.BetaHashCacheMaxSize = 10000       // Hash cache size
```

### Presets

| Preset | Sharing | Cache Size | Use Case |
|--------|---------|------------|----------|
| `DefaultChainPerformanceConfig()` | ❌ Off | 10K | Safe default |
| `HighPerformanceConfig()` | ✅ On | 100K | Production |
| `LowMemoryConfig()` | ❌ Off | 1K | Constrained |

---

## Metrics

### Sharing Statistics

```go
stats := network.GetBetaSharingStats()
```

**Fields:**
- `TotalSharedNodes` - Unique shared nodes
- `TotalRequests` - Total GetOrCreateJoinNode calls
- `SharedReuses` - Number of reuses
- `UniqueCreations` - Number of new nodes
- `SharingRatio` - Reuse percentage (0.0-1.0)
- `HashCacheHitRate` - Cache efficiency

### Build Metrics

```go
metrics := network.GetBetaBuildMetrics()
```

**Fields:**
- `TotalJoinNodesRequested` - Total build requests
- `SharedJoinNodesReused` - Nodes reused
- `UniqueJoinNodesCreated` - Nodes created new

### Reset Metrics

```go
network.ResetChainMetrics()  // Resets all metrics
```

---

## Usage Patterns

### Automatic (Recommended)

```go
// Rules automatically use sharing when enabled
// No code changes required!
network := rete.NewReteNetworkWithConfig(storage, config)
// ... build rules as usual
```

### Manual (Advanced)

```go
if network.BetaSharingRegistry != nil {
    node, hash, wasShared, err := network.BetaSharingRegistry.GetOrCreateJoinNode(
        condition,
        leftVars,
        rightVars,
        allVars,
        varTypes,
        storage,
    )
    // Use node...
}
```

---

## Log Messages

### Symbols

- ✨ **New shared node** created
- ♻️ **Existing node** reused
- 🔗 **Connection** established
- 📊 **Refcount** updated
- ⚡ **Optimization** applied

### Example Output

```
✨ Created new shared JoinNode join_abc123 (hash: join_abc123)
♻️ Reused shared JoinNode join_abc123 (hash: join_abc123)
📊 Node join_abc123 now used by 3 rule(s)
```

---

## Troubleshooting

### Problem: Sharing not working

**Check:**
```go
// 1. Is sharing enabled?
fmt.Println(network.Config.BetaSharingEnabled)

// 2. Is registry initialized?
fmt.Println(network.BetaSharingRegistry != nil)

// 3. Check metrics
stats := network.GetBetaSharingStats()
fmt.Printf("Sharing ratio: %.1f%%\n", stats.SharingRatio * 100)
```

### Problem: Low sharing ratio

**Causes:**
- Rules have unique patterns
- Different variable types
- Different condition structures

**Check:**
```go
// List all unique patterns
hashes := network.BetaSharingRegistry.ListSharedJoinNodes()
fmt.Printf("Unique patterns: %d\n", len(hashes))
```

### Problem: High memory usage

**Solutions:**
```go
// Reduce cache size
config.BetaHashCacheMaxSize = 1000

// Or disable sharing
config.BetaSharingEnabled = false
```

---

## Performance Tips

### Optimal Configuration

```go
// Production
config := rete.HighPerformanceConfig()
config.BetaSharingEnabled = true
config.BetaHashCacheMaxSize = 100000

// Development
config := rete.DefaultChainPerformanceConfig()
config.BetaSharingEnabled = true
config.MetricsEnabled = true
config.MetricsDetailedChains = true

// Memory-constrained
config := rete.LowMemoryConfig()
config.BetaSharingEnabled = false
```

### Expected Impact

- **Memory:** 30-50% reduction (with 50%+ identical patterns)
- **Performance:** 20-40% improvement (due to reduced memory pressure)
- **Overhead:** ~1-2 µs per node creation

---

## API Reference

### ReteNetwork Methods

```go
// Get sharing statistics
network.GetBetaSharingStats() *BetaSharingStats

// Get build metrics
network.GetBetaBuildMetrics() *BetaBuildMetrics

// Reset metrics
network.ResetChainMetrics()

// Get configuration
network.GetConfig() *ChainPerformanceConfig
```

### BetaSharingRegistry Methods

```go
// Get or create shared node
GetOrCreateJoinNode(condition, leftVars, rightVars, allVars, varTypes, storage) (*JoinNode, string, bool, error)

// List shared nodes
ListSharedJoinNodes() []string

// Get node details
GetSharedJoinNodeDetails(hash) (*JoinNodeDetails, error)

// Get statistics
GetSharingStats() *BetaSharingStats
```

---

## Code Examples

### Example 1: Enable and Monitor

```go
config := rete.HighPerformanceConfig()
network := rete.NewReteNetworkWithConfig(storage, config)

// Build rules...
// (rules automatically use sharing)

// Check results
stats := network.GetBetaSharingStats()
fmt.Printf("Shared %d nodes, %.1f%% reuse\n", 
    stats.TotalSharedNodes,
    stats.SharingRatio * 100)
```

### Example 2: Conditional Enable

```go
config := rete.DefaultChainPerformanceConfig()

if os.Getenv("ENABLE_BETA_SHARING") == "true" {
    config.BetaSharingEnabled = true
    log.Println("Beta sharing enabled")
}

network := rete.NewReteNetworkWithConfig(storage, config)
```

### Example 3: Development Debug

```go
config := rete.DefaultChainPerformanceConfig()
config.BetaSharingEnabled = true
config.MetricsEnabled = true
config.MetricsDetailedChains = true

network := rete.NewReteNetworkWithConfig(storage, config)

// After building rules
stats := network.GetBetaSharingStats()
fmt.Printf("Debug Info:\n")
fmt.Printf("  Total requests: %d\n", stats.TotalRequests)
fmt.Printf("  Reuses: %d\n", stats.SharedReuses)
fmt.Printf("  New creations: %d\n", stats.UniqueCreations)
fmt.Printf("  Sharing ratio: %.1f%%\n", stats.SharingRatio * 100)
fmt.Printf("  Cache hit rate: %.1f%%\n", stats.HashCacheHitRate * 100)
```

---

## Backward Compatibility

### ✅ Guaranteed Compatible

- Default behavior unchanged (`BetaSharingEnabled = false`)
- All existing tests pass
- Nil-safe accessors (return `nil` when disabled)
- Legacy builder paths preserved
- No breaking API changes

### Migration Path

1. **Phase 1:** Keep disabled, test in dev
2. **Phase 2:** Enable for specific workloads
3. **Phase 3:** Enable by default (future)

---

## Additional Resources

- **Migration Guide:** `BETA_SHARING_MIGRATION.md`
- **Technical Details:** `BETA_SHARING_INTEGRATION_SUMMARY.md`
- **Implementation:** `beta_chain_builder.go`, `beta_sharing.go`
- **Tests:** `beta_sharing_integration_test.go`

---

**Version:** 1.0  
**License:** MIT  
**Project:** TSD RETE Engine