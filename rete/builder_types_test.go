// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

func TestNewTypeBuilder(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)

	tb := NewTypeBuilder(utils)

	if tb == nil {
		t.Fatal("NewTypeBuilder returned nil")
	}

	if tb.utils != utils {
		t.Error("TypeBuilder.utils not set correctly")
	}
}

func TestTypeBuilder_CreateTypeDefinition(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	tb := NewTypeBuilder(utils)

	tests := []struct {
		name      string
		typeName  string
		typeMap   map[string]interface{}
		wantName  string
		wantType  string
		numFields int
	}{
		{
			name:     "simple type with no fields",
			typeName: "Person",
			typeMap: map[string]interface{}{
				"name": "Person",
			},
			wantName:  "Person",
			wantType:  "type",
			numFields: 0,
		},
		{
			name:     "type with single field",
			typeName: "Person",
			typeMap: map[string]interface{}{
				"name": "Person",
				"fields": []interface{}{
					map[string]interface{}{
						"name": "age",
						"type": "number",
					},
				},
			},
			wantName:  "Person",
			wantType:  "type",
			numFields: 1,
		},
		{
			name:     "type with multiple fields",
			typeName: "Employee",
			typeMap: map[string]interface{}{
				"name": "Employee",
				"fields": []interface{}{
					map[string]interface{}{
						"name": "id",
						"type": "number",
					},
					map[string]interface{}{
						"name": "name",
						"type": "string",
					},
					map[string]interface{}{
						"name": "salary",
						"type": "number",
					},
				},
			},
			wantName:  "Employee",
			wantType:  "type",
			numFields: 3,
		},
		{
			name:     "type with invalid field format",
			typeName: "BadType",
			typeMap: map[string]interface{}{
				"name": "BadType",
				"fields": []interface{}{
					"not a map", // Invalid field
				},
			},
			wantName:  "BadType",
			wantType:  "type",
			numFields: 0, // Invalid field should be skipped
		},
		{
			name:     "type with missing field name",
			typeName: "PartialType",
			typeMap: map[string]interface{}{
				"name": "PartialType",
				"fields": []interface{}{
					map[string]interface{}{
						"type": "string", // Missing name
					},
				},
			},
			wantName:  "PartialType",
			wantType:  "type",
			numFields: 0, // Field without name should be skipped
		},
		{
			name:     "type with missing field type",
			typeName: "PartialType2",
			typeMap: map[string]interface{}{
				"name": "PartialType2",
				"fields": []interface{}{
					map[string]interface{}{
						"name": "field1", // Missing type
					},
				},
			},
			wantName:  "PartialType2",
			wantType:  "type",
			numFields: 0, // Field without type should be skipped
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeDef := tb.CreateTypeDefinition(tt.typeName, tt.typeMap)

			if typeDef.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", typeDef.Name, tt.wantName)
			}

			if typeDef.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", typeDef.Type, tt.wantType)
			}

			if len(typeDef.Fields) != tt.numFields {
				t.Errorf("Number of fields = %d, want %d", len(typeDef.Fields), tt.numFields)
			}

			// Verify field details for multi-field test
			if tt.numFields == 3 && tt.typeName == "Employee" {
				expectedFields := []struct {
					name string
					typ  string
				}{
					{"id", "number"},
					{"name", "string"},
					{"salary", "number"},
				}

				for i, expected := range expectedFields {
					if typeDef.Fields[i].Name != expected.name {
						t.Errorf("Field[%d].Name = %q, want %q", i, typeDef.Fields[i].Name, expected.name)
					}
					if typeDef.Fields[i].Type != expected.typ {
						t.Errorf("Field[%d].Type = %q, want %q", i, typeDef.Fields[i].Type, expected.typ)
					}
				}
			}
		})
	}
}

func TestTypeBuilder_CreateTypeNodes(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	tb := NewTypeBuilder(utils)

	t.Run("create single type node", func(t *testing.T) {
		network := NewReteNetwork(storage)

		types := []interface{}{
			map[string]interface{}{
				"name": "Person",
				"fields": []interface{}{
					map[string]interface{}{
						"name": "age",
						"type": "number",
					},
				},
			},
		}

		err := tb.CreateTypeNodes(network, types, storage)
		if err != nil {
			t.Fatalf("CreateTypeNodes failed: %v", err)
		}

		// Verify TypeNode was created
		typeNode, exists := network.TypeNodes["Person"]
		if !exists {
			t.Fatal("TypeNode 'Person' not created")
		}

		if typeNode.TypeName != "Person" {
			t.Errorf("TypeNode.TypeName = %q, want 'Person'", typeNode.TypeName)
		}

		// Verify TypeNode is connected to RootNode
		foundChild := false
		for _, child := range network.RootNode.Children {
			if child == typeNode {
				foundChild = true
				break
			}
		}
		if !foundChild {
			t.Error("TypeNode not connected to RootNode")
		}
	})

	t.Run("create multiple type nodes", func(t *testing.T) {
		network := NewReteNetwork(storage)

		types := []interface{}{
			map[string]interface{}{
				"name": "Person",
				"fields": []interface{}{
					map[string]interface{}{
						"name": "age",
						"type": "number",
					},
				},
			},
			map[string]interface{}{
				"name": "Employee",
				"fields": []interface{}{
					map[string]interface{}{
						"name": "salary",
						"type": "number",
					},
				},
			},
			map[string]interface{}{
				"name": "Department",
			},
		}

		err := tb.CreateTypeNodes(network, types, storage)
		if err != nil {
			t.Fatalf("CreateTypeNodes failed: %v", err)
		}

		// Verify all TypeNodes were created
		expectedTypes := []string{"Person", "Employee", "Department"}
		for _, typeName := range expectedTypes {
			if _, exists := network.TypeNodes[typeName]; !exists {
				t.Errorf("TypeNode '%s' not created", typeName)
			}
		}

		if len(network.TypeNodes) != 3 {
			t.Errorf("Expected 3 TypeNodes, got %d", len(network.TypeNodes))
		}

		// Verify all are connected to RootNode
		if len(network.RootNode.Children) != 3 {
			t.Errorf("Expected 3 children on RootNode, got %d", len(network.RootNode.Children))
		}
	})

	t.Run("error on invalid type format", func(t *testing.T) {
		network := NewReteNetwork(storage)

		types := []interface{}{
			"not a map", // Invalid format
		}

		err := tb.CreateTypeNodes(network, types, storage)
		if err == nil {
			t.Error("Expected error for invalid type format, got nil")
		}
	})

	t.Run("error on missing type name", func(t *testing.T) {
		network := NewReteNetwork(storage)

		types := []interface{}{
			map[string]interface{}{
				// Missing "name" field
				"fields": []interface{}{},
			},
		}

		err := tb.CreateTypeNodes(network, types, storage)
		if err == nil {
			t.Error("Expected error for missing type name, got nil")
		}
	})

	t.Run("error on invalid type name format", func(t *testing.T) {
		network := NewReteNetwork(storage)

		types := []interface{}{
			map[string]interface{}{
				"name": 123, // Not a string
			},
		}

		err := tb.CreateTypeNodes(network, types, storage)
		if err == nil {
			t.Error("Expected error for invalid type name format, got nil")
		}
	})

	t.Run("with lifecycle manager", func(t *testing.T) {
		network := NewReteNetwork(storage)
		network.LifecycleManager = NewLifecycleManager()

		types := []interface{}{
			map[string]interface{}{
				"name": "TestType",
			},
		}

		err := tb.CreateTypeNodes(network, types, storage)
		if err != nil {
			t.Fatalf("CreateTypeNodes failed: %v", err)
		}

		// Verify TypeNode was registered with LifecycleManager
		typeNode := network.TypeNodes["TestType"]
		nodeType := network.LifecycleManager.GetNodeType(typeNode.GetID())
		if nodeType != "type" {
			t.Errorf("Node type in LifecycleManager = %q, want 'type'", nodeType)
		}
	})
}

func TestTypeBuilder_Integration(t *testing.T) {
	// Integration test: Create types and verify the complete graph structure
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	tb := NewTypeBuilder(utils)
	network := NewReteNetwork(storage)

	types := []interface{}{
		map[string]interface{}{
			"name": "Person",
			"fields": []interface{}{
				map[string]interface{}{"name": "name", "type": "string"},
				map[string]interface{}{"name": "age", "type": "number"},
			},
		},
		map[string]interface{}{
			"name": "Employee",
			"fields": []interface{}{
				map[string]interface{}{"name": "id", "type": "number"},
				map[string]interface{}{"name": "department", "type": "string"},
				map[string]interface{}{"name": "salary", "type": "number"},
			},
		},
	}

	err := tb.CreateTypeNodes(network, types, storage)
	if err != nil {
		t.Fatalf("CreateTypeNodes failed: %v", err)
	}

	// Verify network structure
	if len(network.TypeNodes) != 2 {
		t.Errorf("Expected 2 type nodes, got %d", len(network.TypeNodes))
	}

	// Verify Person type
	personNode := network.TypeNodes["Person"]
	if personNode == nil {
		t.Fatal("Person TypeNode is nil")
	}
	if len(personNode.TypeDef.Fields) != 2 {
		t.Errorf("Person should have 2 fields, got %d", len(personNode.TypeDef.Fields))
	}

	// Verify Employee type
	employeeNode := network.TypeNodes["Employee"]
	if employeeNode == nil {
		t.Fatal("Employee TypeNode is nil")
	}
	if len(employeeNode.TypeDef.Fields) != 3 {
		t.Errorf("Employee should have 3 fields, got %d", len(employeeNode.TypeDef.Fields))
	}

	// Verify both are children of RootNode
	if len(network.RootNode.Children) != 2 {
		t.Errorf("RootNode should have 2 children, got %d", len(network.RootNode.Children))
	}

	// Verify propagation path exists
	for _, child := range network.RootNode.Children {
		typeNode, ok := child.(*TypeNode)
		if !ok {
			t.Error("RootNode child is not a TypeNode")
			continue
		}

		if typeNode.TypeName != "Person" && typeNode.TypeName != "Employee" {
			t.Errorf("Unexpected TypeNode: %s", typeNode.TypeName)
		}
	}
}
