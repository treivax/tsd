// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/treivax/tsd/constraint"
)

// TestEnrichProgramWithNetworkTypes tests the enrichProgramWithNetworkTypes function
// that merges types from the network into the program for incremental validation
func TestEnrichProgramWithNetworkTypes(t *testing.T) {
	t.Run("merge network types into program with no types", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)

		// Add some types to the network
		network.Types = []TypeDefinition{
			{
				Type: "typeDefinition",
				Name: "Person",
				Fields: []Field{
					{Name: "id", Type: "string", IsPrimaryKey: true},
					{Name: "name", Type: "string"},
				},
			},
			{
				Type: "typeDefinition",
				Name: "Company",
				Fields: []Field{
					{Name: "id", Type: "string", IsPrimaryKey: true},
					{Name: "name", Type: "string"},
				},
			},
		}

		// Create a program with no types (just facts)
		program := &constraint.Program{
			Types: []constraint.TypeDefinition{},
			Facts: []constraint.Fact{
				{TypeName: "Person"},
			},
		}

		pipeline := NewConstraintPipeline()
		enriched := pipeline.enrichProgramWithNetworkTypes(program, network)

		// Should have 2 types from network
		if len(enriched.Types) != 2 {
			t.Errorf("Expected 2 types, got %d", len(enriched.Types))
		}

		// Verify Person type
		personFound := false
		companyFound := false
		for _, typeDef := range enriched.Types {
			if typeDef.Name == "Person" {
				personFound = true
				if len(typeDef.Fields) != 2 {
					t.Errorf("Person type should have 2 fields, got %d", len(typeDef.Fields))
				}
				// Check primary key is preserved
				pkFound := false
				for _, field := range typeDef.Fields {
					if field.Name == "id" && field.IsPrimaryKey {
						pkFound = true
					}
				}
				if !pkFound {
					t.Error("Primary key 'id' not found in Person type")
				}
			}
			if typeDef.Name == "Company" {
				companyFound = true
			}
		}

		if !personFound {
			t.Error("Person type not found in enriched program")
		}
		if !companyFound {
			t.Error("Company type not found in enriched program")
		}
	})

	t.Run("avoid duplicates when program already has some types", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)

		// Add types to the network
		network.Types = []TypeDefinition{
			{
				Type: "typeDefinition",
				Name: "Person",
				Fields: []Field{
					{Name: "id", Type: "string", IsPrimaryKey: true},
					{Name: "name", Type: "string"},
				},
			},
			{
				Type: "typeDefinition",
				Name: "Company",
				Fields: []Field{
					{Name: "id", Type: "string", IsPrimaryKey: true},
					{Name: "name", Type: "string"},
				},
			},
		}

		// Create a program that already has Person type (avoid duplication)
		program := &constraint.Program{
			Types: []constraint.TypeDefinition{
				{
					Type: "typeDefinition",
					Name: "Person",
					Fields: []constraint.Field{
						{Name: "id", Type: "string", IsPrimaryKey: true},
						{Name: "name", Type: "string"},
					},
				},
			},
		}

		pipeline := NewConstraintPipeline()
		enriched := pipeline.enrichProgramWithNetworkTypes(program, network)

		// Should have 2 types total (Person from program, Company from network)
		if len(enriched.Types) != 2 {
			t.Errorf("Expected 2 types (no duplicates), got %d", len(enriched.Types))
		}

		// Count Person types - should be exactly 1
		personCount := 0
		for _, typeDef := range enriched.Types {
			if typeDef.Name == "Person" {
				personCount++
			}
		}
		if personCount != 1 {
			t.Errorf("Expected exactly 1 Person type, got %d", personCount)
		}
	})

	t.Run("preserve original program when no network types", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		// No types in network

		program := &constraint.Program{
			Types: []constraint.TypeDefinition{
				{
					Type: "typeDefinition",
					Name: "Person",
					Fields: []constraint.Field{
						{Name: "id", Type: "string"},
					},
				},
			},
		}

		pipeline := NewConstraintPipeline()
		enriched := pipeline.enrichProgramWithNetworkTypes(program, network)

		// Should have only the program's own type
		if len(enriched.Types) != 1 {
			t.Errorf("Expected 1 type, got %d", len(enriched.Types))
		}
		if enriched.Types[0].Name != "Person" {
			t.Errorf("Expected Person type, got %s", enriched.Types[0].Name)
		}
	})
}

// TestIncrementalValidationWithMultipleFiles tests the complete workflow
// of loading types and facts from separate files
func TestIncrementalValidationWithMultipleFiles(t *testing.T) {
	t.Run("types in file 1, facts in file 2", func(t *testing.T) {
		tempDir := t.TempDir()

		// File 1: Types with primary keys
		typesFile := filepath.Join(tempDir, "types.tsd")
		typesContent := `
type Person(#id: string, name: string, age: number)
type Company(#id: string, name: string, employees: number)
`
		err := os.WriteFile(typesFile, []byte(typesContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write types file: %v", err)
		}

		// File 2: Facts referencing the types
		factsFile := filepath.Join(tempDir, "facts.tsd")
		factsContent := `
Person(id: "P001", name: "Alice", age: 30)
Person(id: "P002", name: "Bob", age: 25)
Company(id: "C001", name: "TechCorp", employees: 250)
`
		err = os.WriteFile(factsFile, []byte(factsContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write facts file: %v", err)
		}

		pipeline := NewConstraintPipeline()
		storage := NewMemoryStorage()

		// Ingest types first
		network, _, err := pipeline.IngestFile(typesFile, nil, storage)
		if err != nil {
			t.Fatalf("Failed to ingest types file: %v", err)
		}

		if network == nil {
			t.Fatal("Network is nil after ingesting types")
		}
		if len(network.TypeNodes) != 2 {
			t.Errorf("Expected 2 TypeNodes, got %d", len(network.TypeNodes))
		}

		// Verify types have primary keys
		personNode, exists := network.TypeNodes["Person"]
		if !exists {
			t.Fatal("Person TypeNode not found")
		}
		pkFound := false
		for _, field := range personNode.TypeDefinition.Fields {
			if field.Name == "id" && field.IsPrimaryKey {
				pkFound = true
			}
		}
		if !pkFound {
			t.Error("Primary key 'id' not found in Person TypeNode")
		}

		// Ingest facts (should succeed because types are enriched from network)
		network, _, err = pipeline.IngestFile(factsFile, network, storage)
		if err != nil {
			t.Fatalf("Failed to ingest facts file: %v", err)
		}

		// Verify facts were submitted
		// The facts should be in storage
		allFacts := storage.GetAllFacts()
		if len(allFacts) != 3 {
			t.Errorf("Expected 3 facts in storage, got %d", len(allFacts))
		}
	})

	t.Run("types split across multiple files", func(t *testing.T) {
		tempDir := t.TempDir()

		// File 1: Person type
		file1 := filepath.Join(tempDir, "person.tsd")
		content1 := `type Person(#id: string, name: string)`
		err := os.WriteFile(file1, []byte(content1), 0644)
		if err != nil {
			t.Fatalf("Failed to write file1: %v", err)
		}

		// File 2: Company type
		file2 := filepath.Join(tempDir, "company.tsd")
		content2 := `type Company(#id: string, name: string)`
		err = os.WriteFile(file2, []byte(content2), 0644)
		if err != nil {
			t.Fatalf("Failed to write file2: %v", err)
		}

		// File 3: Facts using both types
		file3 := filepath.Join(tempDir, "data.tsd")
		content3 := `
Person(id: "P1", name: "Alice")
Company(id: "C1", name: "TechCorp")
`
		err = os.WriteFile(file3, []byte(content3), 0644)
		if err != nil {
			t.Fatalf("Failed to write file3: %v", err)
		}

		pipeline := NewConstraintPipeline()
		storage := NewMemoryStorage()

		// Ingest files sequentially
		var network *ReteNetwork
		files := []string{file1, file2, file3}
		for i, file := range files {
			network, _, err = pipeline.IngestFile(file, network, storage)
			if err != nil {
				t.Fatalf("Failed to ingest file %d (%s): %v", i+1, file, err)
			}
		}

		if network == nil {
			t.Fatal("Network is nil")
		}

		// Should have both types
		if len(network.TypeNodes) != 2 {
			t.Errorf("Expected 2 TypeNodes, got %d", len(network.TypeNodes))
		}

		// Should have both facts
		allFacts := storage.GetAllFacts()
		if len(allFacts) != 2 {
			t.Errorf("Expected 2 facts, got %d", len(allFacts))
		}
	})

	t.Run("complex types with multiple fields reference from previous file", func(t *testing.T) {
		tempDir := t.TempDir()

		// File 1: Type definition with many fields
		typesFile := filepath.Join(tempDir, "schema.tsd")
		typesContent := `type User(#id: string, username: string, email: string, age: number, active: bool)`
		err := os.WriteFile(typesFile, []byte(typesContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write types file: %v", err)
		}

		// File 2: Regular facts (not assignments)
		factsFile := filepath.Join(tempDir, "users.tsd")
		factsContent := `
User(id: "U1", username: "alice", email: "alice@example.com", age: 30, active: true)
User(id: "U2", username: "bob", email: "bob@example.com", age: 25, active: false)
User(id: "U3", username: "charlie", email: "charlie@example.com", age: 35, active: true)
`
		err = os.WriteFile(factsFile, []byte(factsContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write facts file: %v", err)
		}

		pipeline := NewConstraintPipeline()
		storage := NewMemoryStorage()

		// Ingest types
		network, _, err := pipeline.IngestFile(typesFile, nil, storage)
		if err != nil {
			t.Fatalf("Failed to ingest types: %v", err)
		}

		// Ingest facts (should work with enriched types from previous file)
		network, _, err = pipeline.IngestFile(factsFile, network, storage)
		if err != nil {
			t.Fatalf("Failed to ingest facts: %v", err)
		}

		// Should have the facts
		allFacts := storage.GetAllFacts()
		if len(allFacts) != 3 {
			t.Errorf("Expected 3 facts, got %d", len(allFacts))
		}
	})
}
