# Multi-Source Aggregation Examples

**Copyright (c) 2025 TSD Contributors**  
**Licensed under the MIT License**

This directory contains comprehensive real-world examples demonstrating multi-source aggregation capabilities in the TSD RETE engine. These examples showcase how to combine data from multiple fact types, compute aggregate functions, and apply business logic using threshold conditions.

---

## Table of Contents

1. [Overview](#overview)
2. [Examples Index](#examples-index)
3. [Quick Start](#quick-start)
4. [Common Patterns](#common-patterns)
5. [Syntax Reference](#syntax-reference)
6. [Running the Examples](#running-the-examples)
7. [Creating Custom Examples](#creating-custom-examples)

---

## Overview

Multi-source aggregations enable you to:

- **Join multiple fact types** on related fields (e.g., Customer → Order → Payment)
- **Compute aggregations** across joined data (AVG, SUM, COUNT, MIN, MAX)
- **Apply thresholds** to trigger rules only when conditions are met
- **Generate insights** from complex relationships in your data

### When to Use Multi-Source Aggregations

✅ **Perfect for:**
- Customer analytics (LTV, churn prediction, segmentation)
- Supply chain monitoring (supplier performance, inventory optimization)
- IoT device management (predictive maintenance, anomaly detection)
- Financial analysis (revenue reconciliation, risk scoring)
- E-commerce analytics (order patterns, product performance)

❌ **Not ideal for:**
- Simple single-fact pattern matching
- High-frequency event processing requiring microsecond latency
- Cases where denormalized data would be simpler

---

## Examples Index

### 1. E-Commerce Analytics
**File:** [`ecommerce_analytics.tsd`](ecommerce_analytics.tsd)

**Business Domain:** Online retail platform  
**Fact Types:** Customer, Order, OrderItem, Payment, Product, Review  
**Use Cases:**
- Customer lifetime value (LTV) calculation
- VIP customer identification
- Abandoned cart detection
- Product performance analysis
- Revenue reconciliation
- Inventory turnover tracking

**Key Rules:**
- `identify_vip_customers` - High-value customer detection
- `detect_abandoned_carts` - Cart abandonment monitoring
- `low_stock_alert` - Inventory alerts with sales velocity
- `revenue_reconciliation` - Payment matching

**Complexity:** ⭐⭐⭐ Intermediate  
**Sources:** Up to 4 concurrent joins  
**Aggregations:** Up to 7 per rule

---

### 2. Supply Chain Monitoring
**File:** [`supply_chain_monitoring.tsd`](supply_chain_monitoring.tsd)

**Business Domain:** Manufacturing and logistics  
**Fact Types:** Supplier, Shipment, ShipmentItem, Part, QualityInspection, SupplierInvoice, Warehouse, WarehouseInventory  
**Use Cases:**
- Supplier reliability scoring
- Quality control monitoring
- Shipment tracking and delays
- Inventory optimization
- Invoice reconciliation
- Warehouse capacity management

**Key Rules:**
- `supplier_reliability_score` - Performance evaluation
- `critical_stock_shortage` - Inventory alerts
- `quality_failure_patterns` - Defect detection
- `supply_chain_risk_score` - Risk assessment

**Complexity:** ⭐⭐⭐⭐ Advanced  
**Sources:** Up to 5 concurrent joins  
**Aggregations:** Up to 9 per rule

---

### 3. IoT Sensor Monitoring
**File:** [`iot_sensor_monitoring.tsd`](iot_sensor_monitoring.tsd)

**Business Domain:** Industrial IoT and smart facilities  
**Fact Types:** Device, SensorReading, Alert, MaintenanceRecord, DeviceFailure, EnergyConsumption, ProductionMetric  
**Use Cases:**
- Device health monitoring
- Predictive maintenance
- Anomaly detection
- Energy efficiency tracking
- Production quality monitoring
- Alert storm detection

**Key Rules:**
- `device_health_score` - Overall health calculation
- `predictive_maintenance_trigger` - Failure prediction
- `alert_storm_detection` - Cascading failure detection
- `device_total_cost_ownership` - TCO analysis

**Complexity:** ⭐⭐⭐⭐⭐ Expert  
**Sources:** Up to 6 concurrent joins  
**Aggregations:** Up to 9 per rule

---

## Quick Start

### 1. View an Example

```bash
cat ecommerce_analytics.tsd
```

### 2. Test with Sample Data

Create a test file with facts:

```tsd
// Include the example rules
// Add your fact definitions here

fact Customer {
  id: "customer123",
  name: "John Doe",
  email: "john@example.com",
  tier: "gold",
  joinDate: "2024-01-15"
}

fact Order {
  id: "order456",
  customerId: "customer123",
  orderDate: "2024-06-20",
  status: "delivered",
  totalAmount: 299.99
}

fact Payment {
  id: "payment789",
  orderId: "order456",
  amount: 299.99,
  method: "credit_card",
  status: "completed",
  processedDate: "2024-06-20"
}
```

### 3. Run with TSD Engine

```bash
# Build and run
tsd run ecommerce_test.tsd

# Or use Go directly
go run cmd/tsd/main.go ecommerce_test.tsd
```

---

## Common Patterns

### Pattern 1: Customer Analytics

**Scenario:** Calculate customer lifetime value

```tsd
rule customer_lifetime_value :
  {c: Customer,
   order_count: COUNT(o.id),
   total_revenue: SUM(p.amount),
   avg_order_value: AVG(o.totalAmount)}
  / {o: Order}
  / {p: Payment}
  / o.customerId == c.id
    AND p.orderId == o.id
    AND p.status == "completed"
  ==> print("Customer LTV:", c.name, "Revenue:", total_revenue)
```

**Key Points:**
- Main entity: Customer
- Source 1: Orders linked by `customerId`
- Source 2: Payments linked by `orderId`
- Three aggregations: COUNT, SUM, AVG

---

### Pattern 2: Quality Monitoring

**Scenario:** Detect quality issues from inspections

```tsd
rule quality_failure_detection :
  {p: Part,
   total_inspections: COUNT(qi.id),
   avg_defect_rate: AVG(qi.defectRate),
   total_failed: SUM(qi.failed)}
  / {si: ShipmentItem}
  / {qi: QualityInspection}
  / si.partId == p.id
    AND qi.shipmentId == si.shipmentId
    AND avg_defect_rate > 10.0
  ==> print("Quality Alert:", p.name, "Defect Rate:", avg_defect_rate)
```

**Key Points:**
- Threshold condition: `avg_defect_rate > 10.0`
- Rule only fires when threshold is met
- Combines shipment and inspection data

---

### Pattern 3: Predictive Analytics

**Scenario:** Predict device maintenance needs

```tsd
rule predictive_maintenance :
  {d: Device,
   failure_count: COUNT(df.id),
   total_downtime: SUM(df.downtime),
   avg_sensor_value: AVG(sr.value),
   alert_count: COUNT(a.id)}
  / {df: DeviceFailure}
  / {sr: SensorReading}
  / {a: Alert}
  / df.deviceId == d.id
    AND sr.deviceId == d.id
    AND a.deviceId == d.id
    AND failure_count >= 2
    AND alert_count >= 5
  ==> print("Maintenance Required:", d.name)
```

**Key Points:**
- Three source types joined to main device
- Multiple threshold conditions (AND)
- Predictive based on historical patterns

---

## Syntax Reference

### Basic Structure

```tsd
rule rule_name :
  {main_var: MainType, agg1: AGG_FUNC(source1.field), agg2: AGG_FUNC(source2.field)}
  / {source1: SourceType1}
  / {source2: SourceType2}
  / join_conditions AND threshold_conditions
  ==> action
```

### Aggregation Functions

| Function | Purpose | Example |
|----------|---------|---------|
| `AVG()` | Average value | `avg_salary: AVG(e.salary)` |
| `SUM()` | Total sum | `total_revenue: SUM(o.amount)` |
| `COUNT()` | Count items | `order_count: COUNT(o.id)` |
| `MIN()` | Minimum value | `min_price: MIN(p.price)` |
| `MAX()` | Maximum value | `max_score: MAX(s.score)` |

### Join Conditions

```tsd
// Simple equality
e.deptId == d.id

// Multiple joins
e.deptId == d.id AND p.employeeId == e.id

// Three-way join
e.deptId == d.id AND p.employeeId == e.id AND t.employeeId == e.id
```

### Threshold Conditions

```tsd
// Single threshold
avg_salary > 60000

// Multiple thresholds
avg_salary > 60000 AND order_count >= 10

// Comparison operators
==  // Equal
>   // Greater than
<   // Less than
>=  // Greater than or equal
<=  // Less than or equal
!=  // Not equal
```

---

## Running the Examples

### Method 1: Direct Execution

```bash
# Navigate to examples directory
cd tsd/examples/multi_source_aggregations

# Run an example
tsd run ecommerce_analytics.tsd

# With verbose output
tsd run -v supply_chain_monitoring.tsd
```

### Method 2: Integration Testing

```bash
# Run Go tests that use examples
cd tsd/rete
go test -v -run TestMultiSourceAggregation
```

### Method 3: Interactive Exploration

```go
package main

import (
    "github.com/tsd/rete"
)

func main() {
    storage := rete.NewMemoryStorage()
    pipeline := rete.NewConstraintPipeline()
    
    // Load example rules
    network, err := pipeline.BuildNetworkFromConstraintFile(
        "examples/multi_source_aggregations/ecommerce_analytics.tsd",
        storage,
    )
    
    // Submit facts and observe activations
    network.SubmitFact(customerFact)
    network.SubmitFact(orderFact)
    // ...
}
```

---

## Creating Custom Examples

### Step 1: Define Your Domain

Identify:
- Main entity (e.g., Customer, Device, Supplier)
- Related entities (Orders, Sensors, Shipments)
- Relationships between entities
- Metrics you want to compute

### Step 2: Create Type Definitions

```tsd
type MainEntity : <
  id: string,
  name: string,
  // ... other fields
>

type RelatedEntity1 : <
  id: string,
  mainEntityId: string,  // Foreign key
  metric: number
>
```

### Step 3: Write Aggregation Rules

```tsd
rule my_analytics_rule :
  {m: MainEntity,
   metric1: AVG(r1.metric),
   metric2: SUM(r2.value)}
  / {r1: RelatedEntity1}
  / {r2: RelatedEntity2}
  / r1.mainEntityId == m.id
    AND r2.relatedEntity1Id == r1.id
    AND metric1 > threshold
  ==> print("Alert:", m.name, "Metric:", metric1)
```

### Step 4: Add Documentation

Include:
- Business context
- Use case descriptions
- Sample data scenarios
- Expected behavior
- Performance considerations

---

## Best Practices

### 1. **Start Simple**
Begin with 2-source joins before moving to 3+ sources.

### 2. **Use Meaningful Names**
```tsd
// Good
avg_customer_lifetime_value: AVG(o.totalAmount)

// Avoid
x: AVG(o.totalAmount)
```

### 3. **Set Appropriate Thresholds**
Base thresholds on business requirements and historical data:
```tsd
// Business-driven threshold
AND avg_order_value > 100  // Company's AOV target
```

### 4. **Document Business Logic**
```tsd
// Rule: Identify at-risk customers
// Triggers when: orders < 3 AND failed_payments > 0
// Action: Flag for retention campaign
rule at_risk_customers : ...
```

### 5. **Consider Performance**
- Fewer sources = faster execution
- Selective joins reduce memory usage
- Set thresholds to filter early

---

## Performance Considerations

### Small Scale (< 1,000 main facts)
- Any number of aggregation functions
- 2-4 source joins
- Real-time processing suitable

### Medium Scale (1,000 - 10,000 main facts)
- Optimize join order (most selective first)
- Consider caching strategies
- Monitor memory usage

### Large Scale (> 10,000 main facts)
- Batch processing recommended
- Use windowing for time-series data
- Consider distributed processing

See [`MULTI_SOURCE_PERFORMANCE_GUIDE.md`](../../rete/MULTI_SOURCE_PERFORMANCE_GUIDE.md) for optimization details.

---

## Troubleshooting

### Issue: Rule Not Firing

**Check:**
1. Are all required facts present?
2. Do join conditions match fact fields?
3. Are threshold conditions realistic?
4. Is data type compatible (string vs number)?

**Debug:**
```tsd
// Add rule without thresholds first
rule debug_rule :
  {m: MainEntity, count: COUNT(r.id)}
  / {r: RelatedEntity}
  / r.mainEntityId == m.id
  ==> print("Debug:", m.id, "Count:", count)
```

### Issue: Too Many Activations

**Solutions:**
- Add more restrictive thresholds
- Refine join conditions
- Use EXISTS patterns for existence checks
- Filter at fact submission level

### Issue: Poor Performance

**Investigate:**
- Run benchmarks (see performance guide)
- Profile with pprof
- Check join fanout ratios
- Optimize join order

---

## Additional Resources

- **Feature Documentation:** [`AGGREGATION_JOIN_FEATURE_SUMMARY.md`](../../AGGREGATION_JOIN_FEATURE_SUMMARY.md)
- **Performance Guide:** [`MULTI_SOURCE_PERFORMANCE_GUIDE.md`](../../rete/MULTI_SOURCE_PERFORMANCE_GUIDE.md)
- **Quick Reference:** See individual example files for inline documentation
- **Test Cases:** [`rete/multi_source_aggregation_test.go`](../../rete/multi_source_aggregation_test.go)

---

## Contributing

Have a great real-world example? Contributions welcome!

1. Create a new `.tsd` file in this directory
2. Include comprehensive comments and documentation
3. Add sample data scenarios
4. Document expected behavior
5. Submit a pull request

**Example Template:**
```tsd
// [Title] with Multi-Source Aggregations
// ========================================
// Business Context: [Describe the scenario]
//
// Fact Types: [List types used]
// Use Cases: [List key use cases]

// Type definitions
type MainType : <...>
type SourceType1 : <...>

// Rules with documentation
rule descriptive_name :
  {m: MainType, agg: AGG(s.field)}
  / {s: SourceType1}
  / join_conditions
  ==> action

// Usage notes at the end
```

---

## License

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License

See [LICENSE](../../LICENSE) file in the project root for full license text.

---

**Last Updated:** 2025-01-XX  
**Maintainers:** TSD Contributors  
**Status:** Production Ready ✅