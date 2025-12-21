// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"log"

	"github.com/treivax/tsd/api"
)

func main() {
	fmt.Println("=== TSD Multi-File Example: Basic Schema/Data Separation ===")
	fmt.Println()

	// Create a new pipeline
	pipeline := api.NewPipeline()

	// Step 1: Load the schema file (types)
	fmt.Println("📁 Loading schema.tsd...")
	schemaResult, err := pipeline.IngestFile("schema.tsd")
	if err != nil {
		log.Fatalf("❌ Error loading schema: %v", err)
	}
	fmt.Printf("✅ Schema loaded: %d types defined\n", schemaResult.TypeCount())
	fmt.Println()

	// Step 2: Load the data file (facts)
	// The types from schema.tsd are automatically available
	fmt.Println("📁 Loading data.tsd...")
	dataResult, err := pipeline.IngestFile("data.tsd")
	if err != nil {
		log.Fatalf("❌ Error loading data: %v", err)
	}
	fmt.Printf("✅ Data loaded: %d facts submitted\n", dataResult.FactCount())
	fmt.Println()

	// Display summary
	fmt.Println("📊 Summary:")
	fmt.Printf("   - Total types: %d\n", dataResult.TypeCount())
	fmt.Printf("   - Total facts: %d\n", dataResult.FactCount())
	fmt.Println()

	// List the types
	fmt.Println("📋 Types in network:")
	// Access types through network
	network := dataResult.Network()
	for _, typeDef := range network.Types {
		fmt.Printf("   - %s (%d fields)\n", typeDef.Name, len(typeDef.Fields))
	}
	fmt.Println()

	fmt.Println("✅ Multi-file loading completed successfully!")
	fmt.Println()
	fmt.Println("Key takeaway:")
	fmt.Println("  - Types defined in schema.tsd were automatically available in data.tsd")
	fmt.Println("  - The RETE network maintains context between file loads")
	fmt.Println("  - This enables modular organization of TSD programs")
}
