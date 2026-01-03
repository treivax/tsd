// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package examples

import (
	"fmt"
	"log"

	"github.com/treivax/tsd/rete/delta"
)

// Example1_BasicUsage demonstrates the most basic delta propagation workflow.
//
// This example shows:
// 1. Creating a delta detector
// 2. Detecting changes between two fact versions
// 3. Understanding what fields changed
func Example1_BasicUsage() {
	fmt.Println("=== Example 1: Basic Delta Detection ===")


	// Step 1: Create a delta detector with default configuration
	detector := delta.NewDeltaDetector()

	// Step 2: Define old and new versions of a fact
	oldProduct := map[string]interface{}{
		"id":       "p123",
		"name":     "Widget",
		"price":    100.0,
		"stock":    50,
		"category": "electronics",
	}

	newProduct := map[string]interface{}{
		"id":       "p123",
		"name":     "Widget",
		"price":    120.0, // Changed!
		"stock":    45,    // Changed!
		"category": "electronics",
	}

	// Step 3: Detect changes
	factDelta, err := detector.DetectDelta(
		oldProduct,
		newProduct,
		"Product~p123", // Fact ID
		"Product",      // Fact Type
	)

	if err != nil {
		log.Fatalf("Error detecting delta: %v", err)
	}

	// Step 4: Examine the results
	if factDelta.IsEmpty() {
		fmt.Println("No changes detected")
		return
	}

	fmt.Printf("Fact ID: %s\n", factDelta.FactID)
	fmt.Printf("Fact Type: %s\n", factDelta.FactType)
	fmt.Printf("Number of changes: %d\n\n", len(factDelta.Fields))

	// Step 5: Inspect individual field changes
	for fieldName, fieldChange := range factDelta.Fields {
		fmt.Printf("Field '%s' changed:\n", fieldName)
		fmt.Printf("  Old Value: %v\n", fieldChange.OldValue)
		fmt.Printf("  New Value: %v\n", fieldChange.NewValue)
		fmt.Printf("  Change Type: %s\n\n", fieldChange.ChangeType)
	}

	// Step 6: Check if specific field changed
	if priceChange, exists := factDelta.Fields["price"]; exists {
		fmt.Printf("Price updated: %.2f → %.2f\n",
			priceChange.OldValue.(float64),
			priceChange.NewValue.(float64))
	}

	// Output:
	// === Example 1: Basic Delta Detection ===
	//
	// Fact ID: Product~p123
	// Fact Type: Product
	// Number of changes: 2
	//
	// Field 'price' changed:
	//   Old Value: 100
	//   New Value: 120
	//   Change Type: modified
	//
	// Field 'stock' changed:
	//   Old Value: 50
	//   New Value: 45
	//   Change Type: modified
	//
	// Price updated: 100.00 → 120.00
}

// Example2_DependencyIndex demonstrates how to use the dependency index
// to find which nodes are affected by field changes.
func Example2_DependencyIndex() {
	fmt.Println("\n=== Example 2: Dependency Index ===")


	// Step 1: Create an index
	index := delta.NewDependencyIndex()

	// Step 2: Register nodes and their field dependencies
	// Alpha node checking product price
	index.AddAlphaNode("alpha_price_check", "Product", []string{"price"})

	// Alpha node checking product stock
	index.AddAlphaNode("alpha_stock_check", "Product", []string{"stock"})

	// Beta node joining Product and Order
	index.AddBetaNode("beta_order_product", "Product", []string{"price", "category"})

	// Terminal node for discount rule
	index.AddTerminalNode("term_discount_rule", "Product", []string{"price"})

	// Step 3: Query which nodes care about a specific field
	fmt.Println("Nodes affected by Product.price:")
	priceNodes := index.GetAffectedNodes("Product", "price")
	for _, node := range priceNodes {
		fmt.Printf("  - %s (%s)\n", node.NodeID, node.NodeType)
	}

	fmt.Println("\nNodes affected by Product.stock:")
	stockNodes := index.GetAffectedNodes("Product", "stock")
	for _, node := range stockNodes {
		fmt.Printf("  - %s (%s)\n", node.NodeID, node.NodeType)
	}

	// Step 4: Create a fact delta
	detector := delta.NewDeltaDetector()
	oldProduct := map[string]interface{}{
		"id":    "p123",
		"price": 100.0,
		"stock": 50,
	}
	newProduct := map[string]interface{}{
		"id":    "p123",
		"price": 120.0, // Only price changed
		"stock": 50,
	}

	factDelta, _ := detector.DetectDelta(oldProduct, newProduct, "Product~p123", "Product")

	// Step 5: Find nodes affected by this specific delta
	fmt.Println("\nNodes affected by this delta (price changed):")
	affectedNodes := index.GetAffectedNodesForDelta(factDelta)
	for _, node := range affectedNodes {
		fmt.Printf("  - %s (%s) uses fields: %v\n",
			node.NodeID, node.NodeType, node.Fields)
	}

	// Step 6: Get index statistics
	stats := index.GetStats()
	fmt.Printf("\nIndex Statistics:\n")
	fmt.Printf("  Total nodes: %d\n", stats.NodeCount)
	fmt.Printf("  Unique fields: %d\n", stats.FieldCount)
	fmt.Printf("  Fact types: %v\n", stats.FactTypes)

	// Output:
	// === Example 2: Dependency Index ===
	//
	// Nodes affected by Product.price:
	//   - alpha_price_check (alpha)
	//   - beta_order_product (beta)
	//   - term_discount_rule (terminal)
	//
	// Nodes affected by Product.stock:
	//   - alpha_stock_check (alpha)
	//
	// Nodes affected by this delta (price changed):
	//   - alpha_price_check (alpha) uses fields: [price]
	//   - beta_order_product (beta) uses fields: [price category]
	//   - term_discount_rule (terminal) uses fields: [price]
	//
	// Index Statistics:
	//   Total nodes: 4
	//   Unique fields: 3
	//   Fact types: [Product]
}

// Example3_ConfiguredDetector shows how to customize the delta detector
// for different use cases.
func Example3_ConfiguredDetector() {
	fmt.Println("\n=== Example 3: Custom Detector Configuration ===")


	// Use case: Financial application with price sensitivity
	config := delta.DetectorConfig{
		// Prices can differ by up to 0.01 (1 cent)
		FloatEpsilon: 0.01,

		// Ignore internal/technical fields
		IgnoreInternalFields: true,
		IgnoredFields:        []string{"updated_at", "version", "_sync_status"},

		// Track when types change (e.g., int → string)
		TrackTypeChanges: true,

		// Enable deep comparison for nested objects
		EnableDeepComparison: false, // Disabled for performance

		// Don't cache comparisons (streaming updates)
		CacheComparisons: false,
	}

	detector := delta.NewDeltaDetectorWithConfig(config)

	// Example: Float epsilon in action
	oldPrice := map[string]interface{}{
		"id":    "p1",
		"price": 99.995,
	}
	newPrice := map[string]interface{}{
		"id":    "p1",
		"price": 100.001, // Differs by 0.006 < epsilon
	}

	factDelta, _ := detector.DetectDelta(oldPrice, newPrice, "Product~p1", "Product")

	if factDelta.IsEmpty() {
		fmt.Println("No significant price change detected (within epsilon)")
	} else {
		fmt.Printf("Price change detected: %v\n", factDelta.Fields["price"])
	}

	// Example: Ignored fields
	oldData := map[string]interface{}{
		"id":         "p2",
		"price":      100.0,
		"updated_at": "2025-01-01T10:00:00Z",
	}
	newData := map[string]interface{}{
		"id":         "p2",
		"price":      100.0,
		"updated_at": "2025-01-02T15:30:00Z", // Changed but ignored
	}

	factDelta2, _ := detector.DetectDelta(oldData, newData, "Product~p2", "Product")

	if factDelta2.IsEmpty() {
		fmt.Println("No changes (updated_at was ignored)")
	}

	// Example: Type changes
	oldValue := map[string]interface{}{
		"id":    "p3",
		"stock": 42, // int
	}
	newValue := map[string]interface{}{
		"id":    "p3",
		"stock": "42", // string
	}

	factDelta3, _ := detector.DetectDelta(oldValue, newValue, "Product~p3", "Product")

	if !factDelta3.IsEmpty() {
		if change, exists := factDelta3.Fields["stock"]; exists {
			if change.ChangeType == delta.ChangeTypeModified {
				fmt.Printf("Type change detected: %T → %T\n",
					change.OldValue, change.NewValue)
			}
		}
	}

	// Output:
	// === Example 3: Custom Detector Configuration ===
	//
	// No significant price change detected (within epsilon)
	// No changes (updated_at was ignored)
	// Type change detected: int → string
}

// RunAllBasicExamples executes all basic examples in sequence.
func RunAllBasicExamples() {
	Example1_BasicUsage()
	Example2_DependencyIndex()
	Example3_ConfiguredDetector()
}
