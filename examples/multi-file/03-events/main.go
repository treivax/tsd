// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/treivax/tsd/api"
)

func main() {
	fmt.Println("=== TSD Multi-File Example: Event Management System ===")
	fmt.Println()

	// Create a new pipeline
	pipeline := api.NewPipeline()

	// Define the files to load in order
	files := []string{
		"schemas/events.tsd",
		"rules/events.tsd",
		"data/conference-2025.tsd",
	}

	// Variable to store the last result
	var lastResult *api.Result

	// Load each file incrementally
	for i, file := range files {
		fmt.Printf("📁 [%d/%d] Loading %s...\n", i+1, len(files), filepath.Base(file))

		result, err := pipeline.IngestFile(file)
		if err != nil {
			log.Fatalf("❌ Error loading %s: %v", file, err)
		}

		// Store the last result
		lastResult = result

		// Display what was loaded
		typeCount := result.TypeCount()
		ruleCount := result.RuleCount()
		factCount := result.FactCount()

		fmt.Printf("   ✅ Loaded: %d types, %d rules, %d facts\n", typeCount, ruleCount, factCount)
	}
	fmt.Println()

	// Display final summary
	fmt.Println("📊 Final Network State:")
	fmt.Println()

	// Get xuple-spaces created by rules
	if lastResult != nil {
		spaceNames := lastResult.XupleSpaceNames()
		if len(spaceNames) > 0 {
			fmt.Printf("📦 Xuple-Spaces (%d):\n", len(spaceNames))
			for _, name := range spaceNames {
				xuples, err := lastResult.GetXuples(name)
				if err != nil {
					fmt.Printf("   - %s: error accessing xuples: %v\n", name, err)
					continue
				}
				fmt.Printf("   - %s: %d xuples\n", name, len(xuples))

				// Show first 3 xuples as examples
				if len(xuples) > 0 {
					displayCount := 3
					if len(xuples) < displayCount {
						displayCount = len(xuples)
					}
					for i := 0; i < displayCount; i++ {
						fmt.Printf("     • %v\n", xuples[i])
					}
					if len(xuples) > displayCount {
						fmt.Printf("     ... and %d more\n", len(xuples)-displayCount)
					}
				}
			}
			fmt.Println()
		}
	}

	// Display metrics
	fmt.Println("📈 Key Insights:")
	fmt.Println("   ✓ Schema loaded from separate file (schemas/events.tsd)")
	fmt.Println("   ✓ Business rules loaded with xuple-space declarations")
	fmt.Println("   ✓ Data loaded referencing types from schema file")
	fmt.Println("   ✓ Rules executed automatically on loaded facts")
	fmt.Println("   ✓ Xuples generated and stored in declared spaces")
	fmt.Println()

	// Demonstrate the power of multi-file organization
	fmt.Println("🎯 Benefits Demonstrated:")
	fmt.Println("   • Modular organization: schemas, rules, and data separated")
	fmt.Println("   • Incremental loading: each file builds on previous context")
	fmt.Println("   • Type safety: facts validated against schema definitions")
	fmt.Println("   • Automatic execution: rules fire as facts are loaded")
	fmt.Println("   • Maintainability: easy to update schemas, rules, or data independently")
	fmt.Println()

	fmt.Println("✅ Multi-file event management system loaded successfully!")
	fmt.Println()
	fmt.Println("💡 Next steps:")
	fmt.Println("   - Modify data/conference-2025.tsd to add more events")
	fmt.Println("   - Add new rules in rules/events.tsd")
	fmt.Println("   - Extend the schema in schemas/events.tsd")
	fmt.Println("   - All changes are isolated and independently testable!")
}
