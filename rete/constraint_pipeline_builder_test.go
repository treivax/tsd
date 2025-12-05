// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConstraintPipeline_buildNetwork tests the main network building function
func TestConstraintPipeline_buildNetwork(t *testing.T) {
	tests := []struct {
		name        string
		types       []interface{}
		expressions []interface{}
		expectError bool
		description string
	}{
		{
			name:        "empty network",
			types:       []interface{}{},
			expressions: []interface{}{},
			expectError: false,
			description: "Building empty network should succeed",
		},
		{
			name: "network with single type",
			types: []interface{}{
				map[string]interface{}{
					"type": "typeDefinition",
					"name": "Person",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "name",
							"type": "string",
						},
						map[string]interface{}{
							"name": "age",
							"type": "number",
						},
					},
				},
			},
			expressions: []interface{}{},
			expectError: false,
			description: "Network with single type should succeed",
		},
		{
			name: "network with multiple types",
			types: []interface{}{
				map[string]interface{}{
					"type": "typeDefinition",
					"name": "Person",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "id",
							"type": "string",
						},
					},
				},
				map[string]interface{}{
					"type": "typeDefinition",
					"name": "Order",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "id",
							"type": "string",
						},
						map[string]interface{}{
							"name": "amount",
							"type": "number",
						},
					},
				},
			},
			expressions: []interface{}{},
			expectError: false,
			description: "Network with multiple types should succeed",
		},
		{
			name: "network with single expression",
			types: []interface{}{
				map[string]interface{}{
					"type": "typeDefinition",
					"name": "Person",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "name",
							"type": "string",
						},
					},
				},
			},
			expressions: []interface{}{
				map[string]interface{}{
					"type":   "expression",
					"ruleId": "rule1",
					"set": map[string]interface{}{
						"type": "set",
						"variables": []interface{}{
							map[string]interface{}{
								"type":     "typedVariable",
								"name":     "p",
								"dataType": "Person",
							},
						},
					},
					"constraints": map[string]interface{}{
						"type": "constraint",
					},
					"action": map[string]interface{}{
						"type": "action",
						"job": map[string]interface{}{
							"type": "jobCall",
							"name": "print",
							"args": []interface{}{"test"},
						},
					},
				},
			},
			expectError: false,
			description: "Network with single expression should succeed",
		},
		{
			name: "complete network with types and expressions",
			types: []interface{}{
				map[string]interface{}{
					"type": "typeDefinition",
					"name": "Person",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "name",
							"type": "string",
						},
					},
				},
			},
			expressions: []interface{}{
				map[string]interface{}{
					"type":   "expression",
					"ruleId": "rule1",
					"set": map[string]interface{}{
						"type": "set",
						"variables": []interface{}{
							map[string]interface{}{
								"type":     "typedVariable",
								"name":     "p",
								"dataType": "Person",
							},
						},
					},
					"constraints": map[string]interface{}{
						"type": "constraint",
					},
					"action": map[string]interface{}{
						"type": "action",
						"job": map[string]interface{}{
							"type": "jobCall",
							"name": "notify",
							"args": []interface{}{"hello"},
						},
					},
				},
			},
			expectError: false,
			description: "Complete network with types and expressions should succeed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMemoryStorage()
			cp := NewConstraintPipeline()

			network, err := cp.buildNetwork(storage, tt.types, tt.expressions)

			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Nil(t, network)
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, network, "Network should not be nil")
				assert.NotNil(t, network.RootNode, "RootNode should not be nil")
			}
		})
	}
}

// TestConstraintPipeline_createTypeDefinition tests type definition creation
func TestConstraintPipeline_createTypeDefinition(t *testing.T) {
	tests := []struct {
		name         string
		typeName     string
		typeMap      map[string]interface{}
		expectedName string
		numFields    int
		description  string
	}{
		{
			name:     "simple type with one field",
			typeName: "Person",
			typeMap: map[string]interface{}{
				"type": "typeDefinition",
				"name": "Person",
				"fields": []interface{}{
					map[string]interface{}{
						"name": "name",
						"type": "string",
					},
				},
			},
			expectedName: "Person",
			numFields:    1,
			description:  "Simple type with one field should be created",
		},
		{
			name:     "type with multiple fields",
			typeName: "Person",
			typeMap: map[string]interface{}{
				"type": "typeDefinition",
				"name": "Person",
				"fields": []interface{}{
					map[string]interface{}{
						"name": "id",
						"type": "string",
					},
					map[string]interface{}{
						"name": "name",
						"type": "string",
					},
					map[string]interface{}{
						"name": "age",
						"type": "number",
					},
					map[string]interface{}{
						"name": "active",
						"type": "bool",
					},
				},
			},
			expectedName: "Person",
			numFields:    4,
			description:  "Type with multiple fields should be created",
		},
		{
			name:     "type with no fields",
			typeName: "Empty",
			typeMap: map[string]interface{}{
				"type":   "typeDefinition",
				"name":   "Empty",
				"fields": []interface{}{},
			},
			expectedName: "Empty",
			numFields:    0,
			description:  "Type with no fields should be created",
		},
		{
			name:     "type with different field types",
			typeName: "Mixed",
			typeMap: map[string]interface{}{
				"type": "typeDefinition",
				"name": "Mixed",
				"fields": []interface{}{
					map[string]interface{}{
						"name": "stringField",
						"type": "string",
					},
					map[string]interface{}{
						"name": "numberField",
						"type": "number",
					},
					map[string]interface{}{
						"name": "boolField",
						"type": "bool",
					},
				},
			},
			expectedName: "Mixed",
			numFields:    3,
			description:  "Type with different field types should be created",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := NewConstraintPipeline()

			typeDef := cp.createTypeDefinition(tt.typeName, tt.typeMap)

			assert.Equal(t, tt.expectedName, typeDef.Name, tt.description)
			assert.Len(t, typeDef.Fields, tt.numFields, "Number of fields should match")
		})
	}
}

// TestConstraintPipeline_createTypeNodes tests type node creation
func TestConstraintPipeline_createTypeNodes(t *testing.T) {
	tests := []struct {
		name        string
		types       []interface{}
		expectError bool
		description string
	}{
		{
			name:        "empty type list",
			types:       []interface{}{},
			expectError: false,
			description: "Empty type list should succeed",
		},
		{
			name: "single type",
			types: []interface{}{
				map[string]interface{}{
					"type": "typeDefinition",
					"name": "Person",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "name",
							"type": "string",
						},
					},
				},
			},
			expectError: false,
			description: "Single type should be created successfully",
		},
		{
			name: "multiple types",
			types: []interface{}{
				map[string]interface{}{
					"type": "typeDefinition",
					"name": "Person",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "name",
							"type": "string",
						},
					},
				},
				map[string]interface{}{
					"type": "typeDefinition",
					"name": "Order",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "id",
							"type": "string",
						},
					},
				},
			},
			expectError: false,
			description: "Multiple types should be created successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMemoryStorage()
			cp := NewConstraintPipeline()
			network := NewReteNetwork(storage)

			err := cp.createTypeNodes(network, tt.types, storage)

			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}

// TestConstraintPipeline_createRuleNodes tests rule node creation
func TestConstraintPipeline_createRuleNodes(t *testing.T) {
	tests := []struct {
		name        string
		expressions []interface{}
		expectError bool
		description string
	}{
		{
			name:        "empty expression list",
			expressions: []interface{}{},
			expectError: false,
			description: "Empty expression list should succeed",
		},
		{
			name: "single expression",
			expressions: []interface{}{
				map[string]interface{}{
					"type":   "expression",
					"ruleId": "rule1",
					"set": map[string]interface{}{
						"type": "set",
						"variables": []interface{}{
							map[string]interface{}{
								"type":     "typedVariable",
								"name":     "p",
								"dataType": "Person",
							},
						},
					},
					"constraints": map[string]interface{}{
						"type": "constraint",
					},
					"action": map[string]interface{}{
						"type": "action",
						"job": map[string]interface{}{
							"type": "jobCall",
							"name": "print",
							"args": []interface{}{"test"},
						},
					},
				},
			},
			expectError: true, // Will fail because Person type doesn't exist
			description: "Single expression without type should fail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMemoryStorage()
			cp := NewConstraintPipeline()
			network := NewReteNetwork(storage)

			// Create Person type if test expects no error
			if !tt.expectError {
				types := []interface{}{
					map[string]interface{}{
						"type": "typeDefinition",
						"name": "Person",
						"fields": []interface{}{
							map[string]interface{}{
								"name": "name",
								"type": "string",
							},
						},
					},
				}
				_ = cp.createTypeNodes(network, types, storage)
			}

			err := cp.createRuleNodes(network, tt.expressions, storage)

			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}

// TestConstraintPipeline_createSingleRule tests single rule creation
func TestConstraintPipeline_createSingleRule(t *testing.T) {
	tests := []struct {
		name        string
		ruleID      string
		exprMap     map[string]interface{}
		expectError bool
		description string
	}{
		{
			name:   "simple rule",
			ruleID: "rule1",
			exprMap: map[string]interface{}{
				"type":   "expression",
				"ruleId": "rule1",
				"set": map[string]interface{}{
					"type": "set",
					"variables": []interface{}{
						map[string]interface{}{
							"type":     "typedVariable",
							"name":     "p",
							"dataType": "Person",
						},
					},
				},
				"constraints": map[string]interface{}{
					"type": "constraint",
				},
				"action": map[string]interface{}{
					"type": "action",
					"job": map[string]interface{}{
						"type": "jobCall",
						"name": "notify",
						"args": []interface{}{"hello"},
					},
				},
			},
			expectError: false,
			description: "Simple rule should be created successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMemoryStorage()
			cp := NewConstraintPipeline()
			network := NewReteNetwork(storage)

			// Create Person type first
			types := []interface{}{
				map[string]interface{}{
					"type": "typeDefinition",
					"name": "Person",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "name",
							"type": "string",
						},
					},
				},
			}
			_ = cp.createTypeNodes(network, types, storage)

			err := cp.createSingleRule(network, tt.ruleID, tt.exprMap, storage)

			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}

// TestConstraintPipeline_isMultiSourceAggregation tests multi-source aggregation detection
func TestConstraintPipeline_isMultiSourceAggregation(t *testing.T) {
	tests := []struct {
		name        string
		exprMap     map[string]interface{}
		expected    bool
		description string
	}{
		{
			name: "single pattern - not multi-source",
			exprMap: map[string]interface{}{
				"set": map[string]interface{}{
					"variables": []interface{}{
						map[string]interface{}{
							"type": "typedVariable",
						},
					},
				},
			},
			expected:    false,
			description: "Single pattern should not be multi-source",
		},
		{
			name: "multiple patterns (3+) - multi-source",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{
						"variables": []interface{}{},
					},
					map[string]interface{}{
						"variables": []interface{}{},
					},
					map[string]interface{}{
						"variables": []interface{}{},
					},
				},
			},
			expected:    true,
			description: "More than 2 patterns should be multi-source",
		},
		{
			name: "two patterns - not multi-source",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{
						"variables": []interface{}{},
					},
					map[string]interface{}{
						"variables": []interface{}{},
					},
				},
			},
			expected:    false,
			description: "Two patterns alone are not multi-source",
		},
		{
			name:        "no patterns - not multi-source",
			exprMap:     map[string]interface{}{},
			expected:    false,
			description: "No patterns should not be multi-source",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := NewConstraintPipeline()

			result := cp.isMultiSourceAggregation(tt.exprMap)

			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

// TestConstraintPipeline_Integration tests complete pipeline scenarios
func TestConstraintPipeline_Integration(t *testing.T) {
	t.Run("complete workflow with types and rules", func(t *testing.T) {
		storage := NewMemoryStorage()
		cp := NewConstraintPipeline()

		// Define types
		types := []interface{}{
			map[string]interface{}{
				"type": "typeDefinition",
				"name": "Person",
				"fields": []interface{}{
					map[string]interface{}{
						"name": "id",
						"type": "string",
					},
					map[string]interface{}{
						"name": "name",
						"type": "string",
					},
					map[string]interface{}{
						"name": "age",
						"type": "number",
					},
				},
			},
			map[string]interface{}{
				"type": "typeDefinition",
				"name": "Order",
				"fields": []interface{}{
					map[string]interface{}{
						"name": "id",
						"type": "string",
					},
					map[string]interface{}{
						"name": "person_id",
						"type": "string",
					},
					map[string]interface{}{
						"name": "amount",
						"type": "number",
					},
				},
			},
		}

		// Define expressions/rules
		expressions := []interface{}{
			map[string]interface{}{
				"type":   "expression",
				"ruleId": "adult_rule",
				"set": map[string]interface{}{
					"type": "set",
					"variables": []interface{}{
						map[string]interface{}{
							"type":     "typedVariable",
							"name":     "p",
							"dataType": "Person",
						},
					},
				},
				"constraints": map[string]interface{}{
					"type": "constraint",
				},
				"action": map[string]interface{}{
					"type": "action",
					"job": map[string]interface{}{
						"type": "jobCall",
						"name": "processAdult",
						"args": []interface{}{"p"},
					},
				},
			},
		}

		// Build network
		network, err := cp.buildNetwork(storage, types, expressions)
		require.NoError(t, err, "Building network should succeed")
		require.NotNil(t, network, "Network should not be nil")
		assert.NotNil(t, network.RootNode, "RootNode should exist")
	})

	t.Run("network with empty inputs", func(t *testing.T) {
		storage := NewMemoryStorage()
		cp := NewConstraintPipeline()

		network, err := cp.buildNetwork(storage, []interface{}{}, []interface{}{})
		require.NoError(t, err, "Building empty network should succeed")
		require.NotNil(t, network, "Network should not be nil")
		assert.NotNil(t, network.RootNode, "RootNode should exist even in empty network")
	})

	t.Run("network with only types, no rules", func(t *testing.T) {
		storage := NewMemoryStorage()
		cp := NewConstraintPipeline()

		types := []interface{}{
			map[string]interface{}{
				"type": "typeDefinition",
				"name": "Person",
				"fields": []interface{}{
					map[string]interface{}{
						"name": "name",
						"type": "string",
					},
				},
			},
		}

		network, err := cp.buildNetwork(storage, types, []interface{}{})
		require.NoError(t, err, "Building network with only types should succeed")
		require.NotNil(t, network, "Network should not be nil")
	})
}
