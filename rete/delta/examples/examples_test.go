// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package examples

import (
	"testing"
)

// TestExample1_BasicUsage validates the basic usage example.
func TestExample1_BasicUsage(t *testing.T) {
	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Example1_BasicUsage panicked: %v", r)
		}
	}()

	Example1_BasicUsage()
}

// TestExample2_DependencyIndex validates the dependency index example.
func TestExample2_DependencyIndex(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Example2_DependencyIndex panicked: %v", r)
		}
	}()

	Example2_DependencyIndex()
}

// TestExample3_ConfiguredDetector validates the configured detector example.
func TestExample3_ConfiguredDetector(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Example3_ConfiguredDetector panicked: %v", r)
		}
	}()

	Example3_ConfiguredDetector()
}

// TestExample4_FullIntegration validates the full integration example.
func TestExample4_FullIntegration(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Example4_FullIntegration panicked: %v", r)
		}
	}()

	Example4_FullIntegration()
}

// TestExample5_ConcurrentUpdates validates the concurrent updates example.
func TestExample5_ConcurrentUpdates(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Example5_ConcurrentUpdates panicked: %v", r)
		}
	}()

	Example5_ConcurrentUpdates()
}

// TestExample6_EcommerceScenario validates the e-commerce scenario example.
func TestExample6_EcommerceScenario(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Example6_EcommerceScenario panicked: %v", r)
		}
	}()

	Example6_EcommerceScenario()
}

// TestExample7_InventoryManagement validates the inventory management example.
func TestExample7_InventoryManagement(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Example7_InventoryManagement panicked: %v", r)
		}
	}()

	Example7_InventoryManagement()
}

// TestIntegratedUpdater validates the IntegratedUpdater implementation.
func TestIntegratedUpdater(t *testing.T) {
	updater := NewIntegratedUpdater()

	// Test registration
	updater.RegisterNode("test_node", "TestType", []string{"field1"})
	if updater.GetNodeCount() != 1 {
		t.Errorf("Expected 1 node, got %d", updater.GetNodeCount())
	}

	// Test update with no changes
	oldFact := map[string]interface{}{
		"id":     "t1",
		"field1": "value1",
		"field2": "value2",
		"field3": "value3",
		"field4": "value4",
	}
	newFact := map[string]interface{}{
		"id":     "t1",
		"field1": "value1",
		"field2": "value2",
		"field3": "value3",
		"field4": "value4",
	}

	result := updater.UpdateFact(oldFact, newFact, "TestType~t1", "TestType")
	if result.Strategy != "noop" {
		t.Errorf("Expected noop strategy for identical facts, got %s", result.Strategy)
	}

	// Test delta propagation (1 field changed out of 5 = 20% < 30% threshold)
	newFact2 := map[string]interface{}{
		"id":     "t1",
		"field1": "value_changed",
		"field2": "value2",
		"field3": "value3",
		"field4": "value4",
	}
	result2 := updater.UpdateFact(oldFact, newFact2, "TestType~t1", "TestType")
	if result2.Strategy != "delta" {
		t.Errorf("Expected delta strategy (1/5 fields = 20%%), got %s", result2.Strategy)
	}
	if result2.NodesEvaluated != 1 {
		t.Errorf("Expected 1 node evaluated, got %d", result2.NodesEvaluated)
	}
}

// TestEcommerceSystem validates the EcommerceSystem implementation.
func TestEcommerceSystem(t *testing.T) {
	sys := NewEcommerceSystem()

	// Verify initial state
	product, exists := sys.products["laptop-123"]
	if !exists {
		t.Fatal("Expected laptop-123 to exist")
	}
	if product.Price != 999.99 {
		t.Errorf("Expected initial price 999.99, got %.2f", product.Price)
	}

	// Test flash sale
	sys.ApplyFlashSale("laptop-123", 0.10)
	if product.Price >= 999.99 {
		t.Errorf("Expected price to decrease after flash sale")
	}

	// Test statistics
	stats := sys
	if stats.operations == 0 {
		t.Error("Expected at least one operation")
	}
}

// TestUpdateResult validates the UpdateResult type.
func TestUpdateResult(t *testing.T) {
	tests := []struct {
		name     string
		result   UpdateResult
		expected string
	}{
		{
			name: "delta strategy",
			result: UpdateResult{
				Strategy:       "delta",
				FieldsChanged:  2,
				NodesEvaluated: 3,
			},
			expected: "Delta propagation (2 fields changed, 3 nodes affected)",
		},
		{
			name: "classic strategy",
			result: UpdateResult{
				Strategy:      "classic",
				FieldsChanged: 5,
			},
			expected: "Classic fallback (5 fields changed, threshold exceeded)",
		},
		{
			name: "noop strategy",
			result: UpdateResult{
				Strategy: "noop",
			},
			expected: "No-op (no changes detected)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := tt.result.Summary()
			if summary != tt.expected {
				t.Errorf("Expected summary %q, got %q", tt.expected, summary)
			}
		})
	}
}

// TestUpdateStatistics validates the UpdateStatistics type.
func TestUpdateStatistics(t *testing.T) {
	stats := UpdateStatistics{
		TotalNodesEvaluated: 30,
		TotalNodesAvoided:   70,
	}

	savings := stats.SavingsPercent()
	expected := 70.0

	if savings != expected {
		t.Errorf("Expected savings %.1f%%, got %.1f%%", expected, savings)
	}

	// Test zero case
	emptyStats := UpdateStatistics{}
	if emptyStats.SavingsPercent() != 0.0 {
		t.Error("Expected 0% savings for empty stats")
	}
}

// TestProductToMap validates the Product.ToMap() method.
func TestProductToMap(t *testing.T) {
	product := &Product{
		ID:          "p1",
		Name:        "Test Product",
		Price:       99.99,
		Stock:       10,
		Category:    "test",
		Status:      "active",
		Brand:       "TestBrand",
		Rating:      4.5,
		ReviewCount: 25,
	}

	m := product.ToMap()

	// Verify all fields are present
	expectedFields := []string{"id", "name", "price", "stock", "category", "status", "brand", "rating", "review_count"}
	for _, field := range expectedFields {
		if _, exists := m[field]; !exists {
			t.Errorf("Expected field %s in product map", field)
		}
	}

	// Verify values
	if m["id"] != "p1" {
		t.Errorf("Expected id 'p1', got %v", m["id"])
	}
	if m["price"] != 99.99 {
		t.Errorf("Expected price 99.99, got %v", m["price"])
	}
}

// BenchmarkIntegratedUpdater benchmarks the integrated updater.
func BenchmarkIntegratedUpdater(b *testing.B) {
	updater := NewIntegratedUpdater()
	updater.RegisterNode("node1", "Product", []string{"price"})
	updater.RegisterNode("node2", "Product", []string{"stock"})

	oldFact := map[string]interface{}{
		"id":    "p1",
		"price": 100.0,
		"stock": 50,
	}
	newFact := map[string]interface{}{
		"id":    "p1",
		"price": 120.0,
		"stock": 50,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		updater.UpdateFact(oldFact, newFact, "Product~p1", "Product")
	}
}

// BenchmarkEcommerceSystem benchmarks the e-commerce system.
func BenchmarkEcommerceSystem(b *testing.B) {
	sys := NewEcommerceSystem()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sys.ApplyFlashSale("laptop-123", 0.10)
	}
}
