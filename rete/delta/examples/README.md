# 📚 Delta Propagation Examples

This directory contains practical, runnable examples demonstrating how to use the RETE delta propagation system in real-world scenarios.

## 📋 Overview

The examples are organized from basic to advanced, covering:

1. **Basic Usage** - Core delta detection and indexing concepts
2. **Full Integration** - Complete workflow with network simulation
3. **E-commerce Scenarios** - Real-world business use cases

## 🚀 Quick Start

### Running All Examples

```bash
# From project root
cd rete/delta/examples
go test -v
```

### Running Specific Examples

```go
package main

import "github.com/treivax/tsd/rete/delta/examples"

func main() {
    // Run all basic examples
    examples.RunAllBasicExamples()
    
    // Run all integration examples
    examples.RunAllIntegrationExamples()
    
    // Run all e-commerce examples
    examples.RunAllEcommerceExamples()
}
```

## 📖 Examples Guide

### 01_basic_usage.go

**What it demonstrates:**
- Creating and configuring a delta detector
- Detecting changes between fact versions
- Understanding field-level changes
- Using the dependency index
- Customizing detector configuration

**Examples included:**
- `Example1_BasicUsage()` - Basic delta detection
- `Example2_DependencyIndex()` - Index queries and node dependencies
- `Example3_ConfiguredDetector()` - Custom configurations for different use cases

**Key concepts:**
```go
// Detect what changed
detector := delta.NewDeltaDetector()
factDelta, _ := detector.DetectDelta(oldFact, newFact, factID, factType)

// Find affected nodes
index := delta.NewDependencyIndex()
affectedNodes := index.GetAffectedNodesForDelta(factDelta)
```

**When to study:**
- ✅ You're new to delta propagation
- ✅ You want to understand the core concepts
- ✅ You need to configure the detector for your use case

---

### 02_full_integration.go

**What it demonstrates:**
- Complete integration pattern combining all components
- Automatic strategy selection (delta vs classic)
- Thread-safe concurrent updates
- Performance metrics and monitoring
- Fallback strategies

**Examples included:**
- `Example4_FullIntegration()` - Full workflow with decision logic
- `Example5_ConcurrentUpdates()` - Thread-safe concurrent processing

**Key components:**
```go
// IntegratedUpdater - Production-ready wrapper
updater := NewIntegratedUpdater()
updater.RegisterNode(nodeID, factType, fields)

// Automatic strategy selection
result := updater.UpdateFact(oldFact, newFact, factID, factType)
// Returns: delta, classic, or noop based on change profile
```

**When to study:**
- ✅ You're ready to integrate delta into your application
- ✅ You need a production-ready pattern
- ✅ You want to handle concurrent updates
- ✅ You need metrics and monitoring

---

### 03_ecommerce_scenario.go

**What it demonstrates:**
- Real-world e-commerce use cases
- Domain-specific optimizations
- Business rule modeling
- Performance analysis in context

**Examples included:**
- `Example6_EcommerceScenario()` - Complete e-commerce system
  - Flash sales (price updates)
  - Inventory management (stock updates)
  - Product recategorization
  - No-op detection
- `Example7_InventoryManagement()` - Inventory-specific patterns

**Scenarios covered:**

| Scenario | Fields Changed | Strategy | Savings |
|----------|----------------|----------|---------|
| Flash Sale (price drop) | 1 | Delta | ~65% |
| Stock update (sale) | 1 | Delta | ~75% |
| Recategorization | 2+ | Classic | 0% |
| Redundant update | 0 | No-op | 100% |

**When to study:**
- ✅ You're building an e-commerce platform
- ✅ You want to see realistic performance gains
- ✅ You need to model business rules
- ✅ You want to optimize for your domain

---

## 🎯 Learning Path

### Beginner (30 minutes)

1. Read [QUICK_START.md](../QUICK_START.md)
2. Run `Example1_BasicUsage()`
3. Experiment with `Example2_DependencyIndex()`
4. Try different configurations in `Example3_ConfiguredDetector()`

### Intermediate (1 hour)

1. Study `IntegratedUpdater` in `02_full_integration.go`
2. Run `Example4_FullIntegration()` and analyze output
3. Review the decision logic (delta vs classic)
4. Understand metrics and statistics

### Advanced (2 hours)

1. Read [MIGRATION.md](../MIGRATION.md)
2. Study `EcommerceSystem` implementation
3. Run `Example6_EcommerceScenario()` with profiling
4. Adapt patterns to your domain
5. Run benchmarks and compare performance

---

## 💡 Common Patterns

### Pattern 1: Wrapper with Automatic Fallback

```go
type FactUpdater struct {
    detector *delta.DeltaDetector
    index    *delta.DependencyIndex
    threshold float64
}

func (u *FactUpdater) UpdateFact(old, new map[string]interface{}) {
    factDelta, _ := u.detector.DetectDelta(old, new, id, typ)
    
    if factDelta.IsEmpty() {
        return // No changes
    }
    
    changeRatio := float64(len(factDelta.Changes)) / float64(len(new))
    if changeRatio <= u.threshold {
        // Delta propagation
        nodes := u.index.GetAffectedNodesForDelta(factDelta)
        propagateToNodes(nodes, factDelta)
    } else {
        // Classic fallback
        retractAndAssert(old, new)
    }
}
```

**Use when:** You want automatic, smart strategy selection.

---

### Pattern 2: Domain-Specific Detector

```go
// Financial domain
financialDetector := delta.NewDeltaDetectorWithConfig(delta.DetectorConfig{
    FloatEpsilon: 0.01, // 1 cent precision
    IgnoredFields: []string{"updated_at", "audit_log"},
})

// IoT sensors
sensorDetector := delta.NewDeltaDetectorWithConfig(delta.DetectorConfig{
    FloatEpsilon: 0.1, // Sensor tolerance
    IgnoreInternalFields: true,
    EnableDeepComparison: false, // Performance
})
```

**Use when:** Your domain has specific tolerance or field requirements.

---

### Pattern 3: Metrics-Driven Optimization

```go
type MetricsCollector struct {
    deltaPropagations int64
    classicFallbacks  int64
    nodesEvaluated    int64
    nodesAvoided      int64
}

func (m *MetricsCollector) RecordUpdate(strategy string, evaluated, total int) {
    switch strategy {
    case "delta":
        atomic.AddInt64(&m.deltaPropagations, 1)
        atomic.AddInt64(&m.nodesEvaluated, int64(evaluated))
        atomic.AddInt64(&m.nodesAvoided, int64(total-evaluated))
    case "classic":
        atomic.AddInt64(&m.classicFallbacks, 1)
        atomic.AddInt64(&m.nodesEvaluated, int64(total))
    }
}

func (m *MetricsCollector) SavingsPercent() float64 {
    total := m.nodesEvaluated + m.nodesAvoided
    return 100.0 * float64(m.nodesAvoided) / float64(total)
}
```

**Use when:** You need to monitor and optimize performance in production.

---

## 🧪 Testing Your Integration

All examples include tests demonstrating how to validate your integration:

```bash
# Run all example tests
go test -v

# Run with coverage
go test -cover

# Run benchmarks
go test -bench=. -benchmem

# Run specific example
go test -v -run TestExample1_BasicUsage
```

---

## 📊 Performance Expectations

Based on benchmarks from `e2e_business_test.go`:

| Network Size | Update Type | Strategy | Speedup | Nodes Saved |
|--------------|-------------|----------|---------|-------------|
| 100 nodes | 1 field | Delta | 3.4x | 80% |
| 100 nodes | 5 fields | Classic | 1.0x | 0% |
| 50 nodes | 1 field | Delta | 2.1x | 60% |
| 20 nodes | 1 field | Classic | 0.9x | N/A (overhead) |

**Key takeaways:**
- Delta shines with >50 nodes and <30% field changes
- Small networks may see overhead
- Classic fallback ensures correctness

---

## 🔧 Customization Guide

### Custom Node Types

```go
type CustomNode struct {
    ID       string
    Type     string // "validation", "transformation", "action"
    FactType string
    Fields   []string
}

// Register with index
index.AddAlphaNode(node.ID, node.FactType, node.Fields)
```

### Custom Change Detection

```go
// Override comparison for specific fields
config := delta.DetectorConfig{
    FloatEpsilon: 0.01,
    IgnoredFields: []string{"metadata"},
    // Add custom validators here
}

detector := delta.NewDeltaDetectorWithConfig(config)
```

---

## 🆘 Troubleshooting

### Issue: "No nodes affected by delta"

**Cause:** Index not populated or field names don't match.

**Solution:**
```go
// Verify index was built
stats := index.GetStats()
fmt.Printf("Index has %d nodes\n", stats.NodeCount)

// Check field names match exactly
affectedNodes := index.GetAffectedNodes("Product", "price") // case-sensitive
```

---

### Issue: "Too many classic fallbacks"

**Cause:** Threshold too low or many fields changing.

**Solution:**
```go
// Increase threshold
updater.deltaThreshold = 0.5 // Allow up to 50% changes

// Or ignore irrelevant fields
detector := delta.NewDeltaDetectorWithConfig(delta.DetectorConfig{
    IgnoredFields: []string{"updated_at", "version"},
})
```

---

### Issue: "Performance worse than classic"

**Cause:** Network too small or overhead not amortized.

**Solution:**
```go
// Use adaptive strategy
if networkSize < 50 {
    useClassic()
} else {
    useDelta()
}
```

---

## 📚 Additional Resources

- [MIGRATION.md](../MIGRATION.md) - Complete migration guide
- [README.md](../README.md) - Architecture and design
- [QUICK_START.md](../QUICK_START.md) - Getting started
- [OPTIMIZATION_GUIDE.md](../OPTIMIZATION_GUIDE.md) - Performance tuning

---

## 🤝 Contributing

Have a new example or pattern? Contributions welcome!

1. Follow the existing structure (Example[N]_[Name])
2. Include tests in `examples_test.go`
3. Add documentation comments
4. Update this README

---

## 📝 License

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License

---

**Last updated:** 2025-01-02  
**Examples version:** 1.0.0