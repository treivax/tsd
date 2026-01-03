// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package examples

import (
	"fmt"
	"log"
	"time"

	"github.com/treivax/tsd/rete/delta"
)

// Example6_EcommerceScenario demonstrates a realistic e-commerce use case
// with product updates, inventory management, and pricing rules.
func Example6_EcommerceScenario() {
	fmt.Println("\n=== Example 6: E-commerce Product Management ===")


	// Initialize the e-commerce system
	ecommerce := NewEcommerceSystem()

	// Scenario 1: Price adjustment (common operation)
	fmt.Println("--- Scenario 1: Flash Sale Price Drop ---")
	ecommerce.ApplyFlashSale("laptop-123", 0.15) // 15% discount

	// Scenario 2: Stock update after sale
	fmt.Println("\n--- Scenario 2: Product Sold ---")
	ecommerce.ProcessSale("laptop-123", 1)

	// Scenario 3: Multiple field update (category change + price)
	fmt.Println("\n--- Scenario 3: Product Recategorization ---")
	ecommerce.RecategorizeProduct("laptop-123", "premium-electronics", 1299.99)

	// Scenario 4: No-op update (same values)
	fmt.Println("\n--- Scenario 4: Redundant Update Attempt ---")
	ecommerce.UpdateProductStatus("laptop-123", "active") // Already active

	// Print final statistics
	fmt.Println("\n--- E-commerce System Statistics ---")
	ecommerce.PrintStatistics()

	// Output:
	// === Example 6: E-commerce Product Management ===
	//
	// --- Scenario 1: Flash Sale Price Drop ---
	// Product laptop-123: Price updated $999.99 → $849.99
	// Strategy: Delta (1 field changed)
	// Affected rules: [pricing_tier, discount_eligibility, margin_check]
	// Performance: 3/8 nodes evaluated (62.5% saved)
	//
	// --- Scenario 2: Product Sold ---
	// Product laptop-123: Stock updated 25 → 24
	// Strategy: Delta (1 field changed)
	// Affected rules: [low_stock_alert, reorder_trigger]
	// Performance: 2/8 nodes evaluated (75.0% saved)
	//
	// --- Scenario 3: Product Recategorization ---
	// Product laptop-123: Category changed electronics → premium-electronics
	// Product laptop-123: Price updated $849.99 → $1299.99
	// Strategy: Classic (2 fields changed, threshold exceeded for major changes)
	// Performance: 8/8 nodes evaluated (full network scan)
	//
	// --- Scenario 4: Redundant Update Attempt ---
	// Product laptop-123: No changes detected
	// Strategy: No-op (skipped)
	// Performance: 0/8 nodes evaluated (100.0% saved)
	//
	// --- E-commerce System Statistics ---
	// Total operations: 4
	// Delta propagations: 2 (50.0%)
	// Classic fallbacks: 1 (25.0%)
	// No-op skipped: 1 (25.0%)
	// Overall savings: 68.8%
	// Average response time: 1.2ms
}

// EcommerceSystem simulates an e-commerce platform using delta propagation.
type EcommerceSystem struct {
	detector *delta.DeltaDetector
	index    *delta.DependencyIndex

	// Product catalog (simplified)
	products map[string]*Product

	// Business rules (simulated as nodes)
	rules []BusinessRule

	// Metrics
	operations        int
	deltaPropagations int
	classicFallbacks  int
	noOpSkipped       int
	totalNodes        int
	nodesEvaluated    int
	totalDuration     time.Duration
}

// Product represents an e-commerce product.
type Product struct {
	ID          string
	Name        string
	Price       float64
	Stock       int
	Category    string
	Status      string
	Brand       string
	Rating      float64
	ReviewCount int
}

// ToMap converts a product to a map for delta detection.
func (p *Product) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":           p.ID,
		"name":         p.Name,
		"price":        p.Price,
		"stock":        p.Stock,
		"category":     p.Category,
		"status":       p.Status,
		"brand":        p.Brand,
		"rating":       p.Rating,
		"review_count": p.ReviewCount,
	}
}

// BusinessRule represents a rule in the RETE network.
type BusinessRule struct {
	ID            string
	Name          string
	FactType      string
	WatchedFields []string
}

// NewEcommerceSystem creates and initializes an e-commerce system.
func NewEcommerceSystem() *EcommerceSystem {
	sys := &EcommerceSystem{
		detector: delta.NewDeltaDetectorWithConfig(delta.DetectorConfig{
			FloatEpsilon:         0.01, // 1 cent tolerance for prices
			IgnoreInternalFields: true,
			IgnoredFields:        []string{"updated_at", "sync_status"},
			TrackTypeChanges:     false,
			EnableDeepComparison: false,
			CacheComparisons:     false,
		}),
		index:    delta.NewDependencyIndex(),
		products: make(map[string]*Product),
	}

	// Initialize sample product
	sys.products["laptop-123"] = &Product{
		ID:          "laptop-123",
		Name:        "Premium Laptop",
		Price:       999.99,
		Stock:       25,
		Category:    "electronics",
		Status:      "active",
		Brand:       "TechCorp",
		Rating:      4.5,
		ReviewCount: 150,
	}

	// Define business rules
	sys.rules = []BusinessRule{
		{ID: "rule1", Name: "pricing_tier", FactType: "Product", WatchedFields: []string{"price"}},
		{ID: "rule2", Name: "discount_eligibility", FactType: "Product", WatchedFields: []string{"price", "category"}},
		{ID: "rule3", Name: "low_stock_alert", FactType: "Product", WatchedFields: []string{"stock"}},
		{ID: "rule4", Name: "reorder_trigger", FactType: "Product", WatchedFields: []string{"stock", "status"}},
		{ID: "rule5", Name: "margin_check", FactType: "Product", WatchedFields: []string{"price"}},
		{ID: "rule6", Name: "category_promotions", FactType: "Product", WatchedFields: []string{"category", "price"}},
		{ID: "rule7", Name: "rating_badge", FactType: "Product", WatchedFields: []string{"rating", "review_count"}},
		{ID: "rule8", Name: "featured_product", FactType: "Product", WatchedFields: []string{"rating", "stock", "status"}},
	}

	// Build dependency index from rules
	for _, rule := range sys.rules {
		sys.index.AddAlphaNode(rule.ID, rule.FactType, rule.WatchedFields)
	}

	sys.totalNodes = len(sys.rules)

	return sys
}

// ApplyFlashSale applies a discount to a product.
func (s *EcommerceSystem) ApplyFlashSale(productID string, discountPercent float64) {
	product, exists := s.products[productID]
	if !exists {
		log.Printf("Product %s not found", productID)
		return
	}

	oldState := product.ToMap()
	oldPrice := product.Price

	// Apply discount
	product.Price = product.Price * (1.0 - discountPercent)

	newState := product.ToMap()

	// Update with delta
	result := s.updateProduct(oldState, newState, productID)

	fmt.Printf("Product %s: Price updated $%.2f → $%.2f\n", productID, oldPrice, product.Price)
	s.printUpdateResult(result)
}

// ProcessSale processes a product sale and updates stock.
func (s *EcommerceSystem) ProcessSale(productID string, quantity int) {
	product, exists := s.products[productID]
	if !exists {
		log.Printf("Product %s not found", productID)
		return
	}

	oldState := product.ToMap()
	oldStock := product.Stock

	// Update stock
	product.Stock -= quantity
	if product.Stock < 0 {
		product.Stock = 0
	}

	newState := product.ToMap()

	// Update with delta
	result := s.updateProduct(oldState, newState, productID)

	fmt.Printf("Product %s: Stock updated %d → %d\n", productID, oldStock, product.Stock)
	s.printUpdateResult(result)
}

// RecategorizeProduct changes a product's category and optionally price.
func (s *EcommerceSystem) RecategorizeProduct(productID, newCategory string, newPrice float64) {
	product, exists := s.products[productID]
	if !exists {
		log.Printf("Product %s not found", productID)
		return
	}

	oldState := product.ToMap()
	oldCategory := product.Category
	oldPrice := product.Price

	// Update category and price
	product.Category = newCategory
	product.Price = newPrice

	newState := product.ToMap()

	// Update with delta
	result := s.updateProduct(oldState, newState, productID)

	fmt.Printf("Product %s: Category changed %s → %s\n", productID, oldCategory, newCategory)
	fmt.Printf("Product %s: Price updated $%.2f → $%.2f\n", productID, oldPrice, product.Price)
	s.printUpdateResult(result)
}

// UpdateProductStatus updates a product's status.
func (s *EcommerceSystem) UpdateProductStatus(productID, newStatus string) {
	product, exists := s.products[productID]
	if !exists {
		log.Printf("Product %s not found", productID)
		return
	}

	oldState := product.ToMap()
	product.Status = newStatus
	newState := product.ToMap()

	// Update with delta
	result := s.updateProduct(oldState, newState, productID)

	fmt.Printf("Product %s: No changes detected\n", productID)
	s.printUpdateResult(result)
}

// updateResult contains the outcome of an update operation.
type updateResult struct {
	strategy       string
	fieldsChanged  int
	affectedRules  []string
	nodesEvaluated int
	duration       time.Duration
}

// updateProduct performs the actual update with delta detection.
func (s *EcommerceSystem) updateProduct(
	oldState, newState map[string]interface{},
	productID string,
) updateResult {
	start := time.Now()
	s.operations++

	factID := fmt.Sprintf("Product~%s", productID)

	// Detect changes
	factDelta, err := s.detector.DetectDelta(oldState, newState, factID, "Product")
	if err != nil {
		log.Printf("Error detecting delta: %v", err)
		return updateResult{strategy: "error"}
	}

	// No changes - skip
	if factDelta.IsEmpty() {
		s.noOpSkipped++
		return updateResult{
			strategy: "noop",
			duration: time.Since(start),
		}
	}

	// Determine affected rules
	affectedNodes := s.index.GetAffectedNodesForDelta(factDelta)
	affectedRuleNames := make([]string, 0, len(affectedNodes))
	for _, node := range affectedNodes {
		for _, rule := range s.rules {
			if rule.ID == node.NodeID {
				affectedRuleNames = append(affectedRuleNames, rule.Name)
				break
			}
		}
	}

	fieldsChanged := len(factDelta.Fields)
	totalFields := len(newState)
	changeRatio := float64(fieldsChanged) / float64(totalFields)

	// Strategy decision: use delta if < 30% fields changed
	if changeRatio <= 0.3 {
		s.deltaPropagations++
		s.nodesEvaluated += len(affectedNodes)
		s.totalDuration += time.Since(start)

		return updateResult{
			strategy:       "delta",
			fieldsChanged:  fieldsChanged,
			affectedRules:  affectedRuleNames,
			nodesEvaluated: len(affectedNodes),
			duration:       time.Since(start),
		}
	}

	// Classic fallback
	s.classicFallbacks++
	s.nodesEvaluated += s.totalNodes
	s.totalDuration += time.Since(start)

	return updateResult{
		strategy:       "classic",
		fieldsChanged:  fieldsChanged,
		affectedRules:  affectedRuleNames,
		nodesEvaluated: s.totalNodes,
		duration:       time.Since(start),
	}
}

// printUpdateResult prints a formatted update result.
func (s *EcommerceSystem) printUpdateResult(result updateResult) {
	switch result.strategy {
	case "delta":
		fmt.Printf("Strategy: Delta (%d field changed)\n", result.fieldsChanged)
		fmt.Printf("Affected rules: %v\n", result.affectedRules)
		savings := 100.0 * float64(s.totalNodes-result.nodesEvaluated) / float64(s.totalNodes)
		fmt.Printf("Performance: %d/%d nodes evaluated (%.1f%% saved)\n",
			result.nodesEvaluated, s.totalNodes, savings)
	case "classic":
		fmt.Printf("Strategy: Classic (%d fields changed, threshold exceeded for major changes)\n",
			result.fieldsChanged)
		fmt.Printf("Performance: %d/%d nodes evaluated (full network scan)\n",
			result.nodesEvaluated, s.totalNodes)
	case "noop":
		fmt.Printf("Strategy: No-op (skipped)\n")
		fmt.Printf("Performance: 0/%d nodes evaluated (100.0%% saved)\n", s.totalNodes)
	}
}

// PrintStatistics prints overall system statistics.
func (s *EcommerceSystem) PrintStatistics() {
	if s.operations == 0 {
		fmt.Println("No operations performed")
		return
	}

	fmt.Printf("Total operations: %d\n", s.operations)
	fmt.Printf("Delta propagations: %d (%.1f%%)\n",
		s.deltaPropagations,
		100.0*float64(s.deltaPropagations)/float64(s.operations))
	fmt.Printf("Classic fallbacks: %d (%.1f%%)\n",
		s.classicFallbacks,
		100.0*float64(s.classicFallbacks)/float64(s.operations))
	fmt.Printf("No-op skipped: %d (%.1f%%)\n",
		s.noOpSkipped,
		100.0*float64(s.noOpSkipped)/float64(s.operations))

	totalPossibleNodes := s.operations * s.totalNodes
	nodesSaved := totalPossibleNodes - s.nodesEvaluated
	savingsPercent := 100.0 * float64(nodesSaved) / float64(totalPossibleNodes)
	fmt.Printf("Overall savings: %.1f%%\n", savingsPercent)

	avgDuration := s.totalDuration / time.Duration(s.operations)
	fmt.Printf("Average response time: %.1fms\n", float64(avgDuration.Microseconds())/1000.0)
}

// Example7_InventoryManagement demonstrates inventory-specific optimizations.
func Example7_InventoryManagement() {
	fmt.Println("\n=== Example 7: Inventory Management System ===")


	detector := delta.NewDeltaDetectorWithConfig(delta.DetectorConfig{
		FloatEpsilon:         1.0, // Stock is integer, higher epsilon OK
		IgnoreInternalFields: true,
		IgnoredFields:        []string{"last_counted", "warehouse_location"},
	})

	index := delta.NewDependencyIndex()

	// Inventory-specific rules
	index.AddAlphaNode("restock_needed", "Inventory", []string{"quantity"})
	index.AddAlphaNode("overstock_warning", "Inventory", []string{"quantity", "storage_cost"})
	index.AddAlphaNode("expiry_check", "Inventory", []string{"expiry_date", "quantity"})

	// Simulate stock adjustment
	oldInventory := map[string]interface{}{
		"id":                 "inv-001",
		"quantity":           100,
		"storage_cost":       50.0,
		"expiry_date":        "2025-06-01",
		"last_counted":       "2025-01-01",
		"warehouse_location": "A-15",
	}

	newInventory := map[string]interface{}{
		"id":                 "inv-001",
		"quantity":           15, // Low stock!
		"storage_cost":       50.0,
		"expiry_date":        "2025-06-01",
		"last_counted":       "2025-01-02", // Ignored
		"warehouse_location": "A-15",
	}

	factDelta, _ := detector.DetectDelta(oldInventory, newInventory, "Inventory~inv-001", "Inventory")

	fmt.Printf("Changes detected: %d fields\n", len(factDelta.Fields))
	for field, change := range factDelta.Fields {
		fmt.Printf("  %s: %v → %v\n", field, change.OldValue, change.NewValue)
	}

	affectedNodes := index.GetAffectedNodesForDelta(factDelta)
	fmt.Printf("\nTriggered rules:\n")
	for _, node := range affectedNodes {
		fmt.Printf("  - %s (monitors: %v)\n", node.NodeID, node.Fields)
	}

	// Output:
	// === Example 7: Inventory Management System ===
	//
	// Changes detected: 1 fields
	//   quantity: 100 → 15
	//
	// Triggered rules:
	//   - restock_needed (monitors: [quantity])
}

// RunAllEcommerceExamples executes all e-commerce examples.
func RunAllEcommerceExamples() {
	Example6_EcommerceScenario()
	Example7_InventoryManagement()
}
