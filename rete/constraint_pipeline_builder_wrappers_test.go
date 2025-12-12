// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
// TestConstraintPipeline_createAlphaRule tests the createAlphaRule wrapper
func TestConstraintPipeline_createAlphaRule(t *testing.T) {
	tests := []struct {
		name          string
		variables     []map[string]interface{}
		variableNames []string
		variableTypes []string
		condition     map[string]interface{}
		action        *Action
		expectError   bool
		description   string
	}{
		{
			name: "simple alpha rule with single variable",
			variables: []map[string]interface{}{
				{
					"name": "p",
					"type": "Person",
				},
			},
			variableNames: []string{"p"},
			variableTypes: []string{"Person"},
			condition: map[string]interface{}{
				"type": "comparison",
				"op":   "==",
				"left": map[string]interface{}{
					"type": "fieldAccess",
					"var":  "p",
					"path": []interface{}{"age"},
				},
				"right": map[string]interface{}{
					"type":  "literal",
					"value": 30.0,
				},
			},
			action: &Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "testAction",
					Args: []interface{}{"Found person"},
				},
			},
			expectError: false,
			description: "Simple alpha rule should succeed",
		},
		{
			name: "alpha rule with complex condition",
			variables: []map[string]interface{}{
				{
					"name": "order",
					"type": "Order",
				},
			},
			variableNames: []string{"order"},
			variableTypes: []string{"Order"},
			condition: map[string]interface{}{
				"type": "and",
				"conditions": []interface{}{
					map[string]interface{}{
						"type": "comparison",
						"op":   ">",
						"left": map[string]interface{}{
							"type": "fieldAccess",
							"var":  "order",
							"path": []interface{}{"amount"},
						},
						"right": map[string]interface{}{
							"type":  "literal",
							"value": 100.0,
						},
					},
					map[string]interface{}{
						"type": "comparison",
						"op":   "==",
						"left": map[string]interface{}{
							"type": "fieldAccess",
							"var":  "order",
							"path": []interface{}{"status"},
						},
						"right": map[string]interface{}{
							"type":  "literal",
							"value": "active",
						},
					},
				},
			},
			action: &Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "logHighValueOrder",
					Args: []interface{}{"High value active order"},
				},
			},
			expectError: false,
			description: "Alpha rule with complex AND condition should succeed",
		},
		{
			name:          "alpha rule with empty variables",
			variables:     []map[string]interface{}{},
			variableNames: []string{},
			variableTypes: []string{},
			condition: map[string]interface{}{
				"type": "comparison",
				"op":   "==",
				"left": map[string]interface{}{
					"type":  "literal",
					"value": true,
				},
				"right": map[string]interface{}{
					"type":  "literal",
					"value": true,
				},
			},
			action: &Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "alwaysTrue",
					Args: []interface{}{"Always true"},
				},
			},
			expectError: true,
			description: "Alpha rule with no variables should fail (no type node)",
		},
		{
			name: "alpha rule with nil action",
			variables: []map[string]interface{}{
				{
					"name": "item",
					"type": "Item",
				},
			},
			variableNames: []string{"item"},
			variableTypes: []string{"Item"},
			condition: map[string]interface{}{
				"type": "comparison",
				"op":   "!=",
				"left": map[string]interface{}{
					"type": "fieldAccess",
					"var":  "item",
					"path": []interface{}{"status"},
				},
				"right": map[string]interface{}{
					"type":  "literal",
					"value": nil,
				},
			},
			action:      nil,
			expectError: false,
			description: "Alpha rule with nil action should succeed (action is optional)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			storage := NewMemoryStorage()
			cp := NewConstraintPipeline()
			network := NewReteNetwork(storage)
			// Create type nodes for the test
			for _, varType := range tt.variableTypes {
				if varType != "" {
					typeDef := TypeDefinition{
						Type:   "typeDefinition",
						Name:   varType,
						Fields: []Field{},
					}
					typeNode := NewTypeNode(varType, typeDef, storage)
					network.TypeNodes[varType] = typeNode
				}
			}
			// Execute
			err := cp.createAlphaRule(
				network,
				"test_rule",
				tt.variables,
				tt.variableNames,
				tt.variableTypes,
				tt.condition,
				tt.action,
				storage,
			)
			// Assert
			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}
// TestConstraintPipeline_createJoinRule tests the createJoinRule wrapper
func TestConstraintPipeline_createJoinRule(t *testing.T) {
	tests := []struct {
		name          string
		variables     []map[string]interface{}
		variableNames []string
		variableTypes []string
		condition     map[string]interface{}
		action        *Action
		expectError   bool
		description   string
	}{
		{
			name: "simple join rule with two variables",
			variables: []map[string]interface{}{
				{
					"name": "p",
					"type": "Person",
				},
				{
					"name": "o",
					"type": "Order",
				},
			},
			variableNames: []string{"p", "o"},
			variableTypes: []string{"Person", "Order"},
			condition: map[string]interface{}{
				"type": "comparison",
				"op":   "==",
				"left": map[string]interface{}{
					"type": "fieldAccess",
					"var":  "p",
					"path": []interface{}{"id"},
				},
				"right": map[string]interface{}{
					"type": "fieldAccess",
					"var":  "o",
					"path": []interface{}{"personId"},
				},
			},
			action: &Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "matchedPersonOrder",
					Args: []interface{}{"Person matched with order"},
				},
			},
			expectError: false,
			description: "Simple join rule should succeed",
		},
		{
			name: "join rule with three variables",
			variables: []map[string]interface{}{
				{"name": "p", "type": "Person"},
				{"name": "o", "type": "Order"},
				{"name": "i", "type": "Item"},
			},
			variableNames: []string{"p", "o", "i"},
			variableTypes: []string{"Person", "Order", "Item"},
			condition: map[string]interface{}{
				"type": "and",
				"conditions": []interface{}{
					map[string]interface{}{
						"type": "comparison",
						"op":   "==",
						"left": map[string]interface{}{
							"type": "fieldAccess",
							"var":  "p",
							"path": []interface{}{"id"},
						},
						"right": map[string]interface{}{
							"type": "fieldAccess",
							"var":  "o",
							"path": []interface{}{"personId"},
						},
					},
					map[string]interface{}{
						"type": "comparison",
						"op":   "==",
						"left": map[string]interface{}{
							"type": "fieldAccess",
							"var":  "o",
							"path": []interface{}{"itemId"},
						},
						"right": map[string]interface{}{
							"type": "fieldAccess",
							"var":  "i",
							"path": []interface{}{"id"},
						},
					},
				},
			},
			action: &Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "threeWayJoin",
					Args: []interface{}{"Three-way join matched"},
				},
			},
			expectError: false,
			description: "Join rule with three variables should succeed",
		},
		{
			name: "join rule with complex condition",
			variables: []map[string]interface{}{
				{"name": "customer", "type": "Customer"},
				{"name": "account", "type": "Account"},
			},
			variableNames: []string{"customer", "account"},
			variableTypes: []string{"Customer", "Account"},
			condition: map[string]interface{}{
				"type": "and",
				"conditions": []interface{}{
					map[string]interface{}{
						"type": "comparison",
						"op":   "==",
						"left": map[string]interface{}{
							"type": "fieldAccess",
							"var":  "customer",
							"path": []interface{}{"id"},
						},
						"right": map[string]interface{}{
							"type": "fieldAccess",
							"var":  "account",
							"path": []interface{}{"customerId"},
						},
					},
					map[string]interface{}{
						"type": "comparison",
						"op":   ">",
						"left": map[string]interface{}{
							"type": "fieldAccess",
							"var":  "account",
							"path": []interface{}{"balance"},
						},
						"right": map[string]interface{}{
							"type":  "literal",
							"value": 1000.0,
						},
					},
				},
			},
			action: &Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "highValueAccount",
					Args: []interface{}{"High value account"},
				},
			},
			expectError: false,
			description: "Join rule with mixed join and filter conditions should succeed",
		},
		{
			name: "join rule with nil condition",
			variables: []map[string]interface{}{
				{"name": "a", "type": "TypeA"},
				{"name": "b", "type": "TypeB"},
			},
			variableNames: []string{"a", "b"},
			variableTypes: []string{"TypeA", "TypeB"},
			condition:     nil,
			action: &Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "unconditionedJoin",
					Args: []interface{}{"Cartesian join"},
				},
			},
			expectError: false,
			description: "Join rule with nil condition should succeed (cartesian product)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			storage := NewMemoryStorage()
			cp := NewConstraintPipeline()
			network := NewReteNetwork(storage)
			// Create type nodes for the test
			for _, varType := range tt.variableTypes {
				if varType != "" {
					typeDef := TypeDefinition{
						Type:   "typeDefinition",
						Name:   varType,
						Fields: []Field{},
					}
					typeNode := NewTypeNode(varType, typeDef, storage)
					network.TypeNodes[varType] = typeNode
				}
			}
			// Execute
			err := cp.createJoinRule(
				network,
				"test_join_rule",
				tt.variables,
				tt.variableNames,
				tt.variableTypes,
				tt.condition,
				tt.action,
				storage,
			)
			// Assert
			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}
// TestConstraintPipeline_createExistsRule tests the createExistsRule wrapper
func TestConstraintPipeline_createExistsRule(t *testing.T) {
	tests := []struct {
		name        string
		exprMap     map[string]interface{}
		condition   map[string]interface{}
		action      *Action
		expectError bool
		description string
	}{
		{
			name: "simple exists rule",
			exprMap: map[string]interface{}{
				"type": "expression",
				"set": map[string]interface{}{
					"variables": []interface{}{
						map[string]interface{}{
							"name":     "p",
							"dataType": "Person",
						},
					},
				},
				"constraints": map[string]interface{}{
					"variable": map[string]interface{}{
						"name":     "o",
						"dataType": "Order",
					},
					"condition": map[string]interface{}{
						"type": "comparison",
						"op":   "==",
						"left": map[string]interface{}{
							"type": "fieldAccess",
							"var":  "o",
							"path": []interface{}{"personId"},
						},
						"right": map[string]interface{}{
							"type": "fieldAccess",
							"var":  "p",
							"path": []interface{}{"id"},
						},
					},
				},
			},
			condition: map[string]interface{}{
				"type": "comparison",
				"op":   ">",
				"left": map[string]interface{}{
					"type": "fieldAccess",
					"var":  "p",
					"path": []interface{}{"age"},
				},
				"right": map[string]interface{}{
					"type":  "literal",
					"value": 18.0,
				},
			},
			action: &Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "personWithOrder",
					Args: []interface{}{"Adult person with orders"},
				},
			},
			expectError: false,
			description: "Simple exists rule should succeed",
		},
		{
			name: "exists rule with complex nested conditions",
			exprMap: map[string]interface{}{
				"type": "expression",
				"set": map[string]interface{}{
					"variables": []interface{}{
						map[string]interface{}{
							"name":     "customer",
							"dataType": "Customer",
						},
					},
				},
				"constraints": map[string]interface{}{
					"variable": map[string]interface{}{
						"name":     "transaction",
						"dataType": "Transaction",
					},
					"conditions": []interface{}{
						map[string]interface{}{
							"type": "comparison",
							"op":   "==",
							"left": map[string]interface{}{
								"type": "fieldAccess",
								"var":  "transaction",
								"path": []interface{}{"customerId"},
							},
							"right": map[string]interface{}{
								"type": "fieldAccess",
								"var":  "customer",
								"path": []interface{}{"id"},
							},
						},
						map[string]interface{}{
							"type": "comparison",
							"op":   ">",
							"left": map[string]interface{}{
								"type": "fieldAccess",
								"var":  "transaction",
								"path": []interface{}{"amount"},
							},
							"right": map[string]interface{}{
								"type":  "literal",
								"value": 500.0,
							},
						},
					},
				},
			},
			condition: map[string]interface{}{
				"type": "comparison",
				"op":   "==",
				"left": map[string]interface{}{
					"type": "fieldAccess",
					"var":  "customer",
					"path": []interface{}{"status"},
				},
				"right": map[string]interface{}{
					"type":  "literal",
					"value": "premium",
				},
			},
			action: &Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "premiumWithLargeTransaction",
					Args: []interface{}{"Premium customer with large transaction"},
				},
			},
			expectError: false,
			description: "Exists rule with complex nested conditions should succeed",
		},
		{
			name: "exists rule with minimal data",
			exprMap: map[string]interface{}{
				"type": "expression",
				"set": map[string]interface{}{
					"variables": []interface{}{
						map[string]interface{}{
							"name":     "item",
							"dataType": "Item",
						},
					},
				},
				"constraints": map[string]interface{}{
					"variable": map[string]interface{}{
						"name":     "review",
						"dataType": "Review",
					},
					"conditions": []interface{}{},
				},
			},
			condition: nil,
			action: &Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "itemWithReview",
					Args: []interface{}{"Item has at least one review"},
				},
			},
			expectError: false,
			description: "Exists rule with minimal conditions should succeed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			storage := NewMemoryStorage()
			cp := NewConstraintPipeline()
			network := NewReteNetwork(storage)
			// Extract types from exprMap and create type nodes
			if set, ok := tt.exprMap["set"].(map[string]interface{}); ok {
				if variables, ok := set["variables"].([]interface{}); ok {
					for _, v := range variables {
						if varMap, ok := v.(map[string]interface{}); ok {
							if varType, ok := varMap["dataType"].(string); ok && varType != "" {
								typeDef := TypeDefinition{
									Type:   "typeDefinition",
									Name:   varType,
									Fields: []Field{},
								}
								typeNode := NewTypeNode(varType, typeDef, storage)
								network.TypeNodes[varType] = typeNode
							}
						}
					}
				}
			}
			// Also create type node for the exists variable
			if constraints, ok := tt.exprMap["constraints"].(map[string]interface{}); ok {
				if variable, ok := constraints["variable"].(map[string]interface{}); ok {
					if existsType, ok := variable["dataType"].(string); ok && existsType != "" {
						typeDef := TypeDefinition{
							Type:   "typeDefinition",
							Name:   existsType,
							Fields: []Field{},
						}
						typeNode := NewTypeNode(existsType, typeDef, storage)
						network.TypeNodes[existsType] = typeNode
					}
				}
			}
			// Execute
			err := cp.createExistsRule(
				network,
				"test_exists_rule",
				tt.exprMap,
				tt.condition,
				tt.action,
				storage,
			)
			// Assert
			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}
// TestConstraintPipeline_extractExistsVariables tests the extractExistsVariables wrapper
func TestConstraintPipeline_extractExistsVariables(t *testing.T) {
	tests := []struct {
		name               string
		exprMap            map[string]interface{}
		expectedMainVar    string
		expectedExistsVar  string
		expectedMainType   string
		expectedExistsType string
		expectError        bool
		description        string
	}{
		{
			name: "standard exists expression",
			exprMap: map[string]interface{}{
				"type": "expression",
				"set": map[string]interface{}{
					"variables": []interface{}{
						map[string]interface{}{
							"name":     "person",
							"dataType": "Person",
						},
					},
				},
				"constraints": map[string]interface{}{
					"variable": map[string]interface{}{
						"name":     "order",
						"dataType": "Order",
					},
				},
			},
			expectedMainVar:    "person",
			expectedExistsVar:  "order",
			expectedMainType:   "Person",
			expectedExistsType: "Order",
			expectError:        false,
			description:        "Standard exists expression should extract all variables correctly",
		},
		{
			name: "exists with multiple main variables (uses first)",
			exprMap: map[string]interface{}{
				"type": "expression",
				"set": map[string]interface{}{
					"variables": []interface{}{
						map[string]interface{}{
							"name":     "customer",
							"dataType": "Customer",
						},
						map[string]interface{}{
							"name":     "account",
							"dataType": "Account",
						},
					},
				},
				"constraints": map[string]interface{}{
					"variable": map[string]interface{}{
						"name":     "transaction",
						"dataType": "Transaction",
					},
				},
			},
			expectedMainVar:    "customer",
			expectedExistsVar:  "transaction",
			expectedMainType:   "Customer",
			expectedExistsType: "Transaction",
			expectError:        false,
			description:        "When multiple variables exist, should use the first one",
		},
		{
			name: "exists without main type specified",
			exprMap: map[string]interface{}{
				"type": "expression",
				"set": map[string]interface{}{
					"variables": []interface{}{
						map[string]interface{}{
							"name": "item",
						},
					},
				},
				"constraints": map[string]interface{}{
					"variable": map[string]interface{}{
						"name":     "review",
						"dataType": "Review",
					},
				},
			},
			expectedMainVar:    "item",
			expectedExistsVar:  "review",
			expectedMainType:   "",
			expectedExistsType: "Review",
			expectError:        false,
			description:        "Missing main type should still succeed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			cp := NewConstraintPipeline()
			// Execute
			mainVar, existsVar, mainType, existsType, err := cp.extractExistsVariables(tt.exprMap)
			// Assert
			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				require.NoError(t, err, tt.description)
				assert.Equal(t, tt.expectedMainVar, mainVar, "Main variable name should match")
				assert.Equal(t, tt.expectedExistsVar, existsVar, "Exists variable name should match")
				assert.Equal(t, tt.expectedMainType, mainType, "Main variable type should match")
				assert.Equal(t, tt.expectedExistsType, existsType, "Exists variable type should match")
			}
		})
	}
}
// TestConstraintPipeline_extractExistsConditions tests the extractExistsConditions wrapper
func TestConstraintPipeline_extractExistsConditions(t *testing.T) {
	tests := []struct {
		name          string
		exprMap       map[string]interface{}
		expectedCount int
		expectError   bool
		description   string
	}{
		{
			name: "exists with single condition",
			exprMap: map[string]interface{}{
				"type": "expression",
				"constraints": map[string]interface{}{
					"condition": map[string]interface{}{
						"type": "comparison",
						"op":   "==",
						"left": map[string]interface{}{
							"type": "fieldAccess",
							"var":  "order",
							"path": []interface{}{"status"},
						},
						"right": map[string]interface{}{
							"type":  "literal",
							"value": "active",
						},
					},
				},
			},
			expectedCount: 1,
			expectError:   false,
			description:   "Single condition should be extracted",
		},
		{
			name: "exists with multiple conditions",
			exprMap: map[string]interface{}{
				"type": "expression",
				"constraints": map[string]interface{}{
					"conditions": []interface{}{
						map[string]interface{}{
							"type": "comparison",
							"op":   ">",
							"left": map[string]interface{}{
								"type": "fieldAccess",
								"var":  "transaction",
								"path": []interface{}{"amount"},
							},
							"right": map[string]interface{}{
								"type":  "literal",
								"value": 100.0,
							},
						},
						map[string]interface{}{
							"type": "comparison",
							"op":   "==",
							"left": map[string]interface{}{
								"type": "fieldAccess",
								"var":  "transaction",
								"path": []interface{}{"type"},
							},
							"right": map[string]interface{}{
								"type":  "literal",
								"value": "credit",
							},
						},
					},
				},
			},
			expectedCount: 2,
			expectError:   false,
			description:   "Multiple conditions should all be extracted",
		},
		{
			name: "exists with no conditions",
			exprMap: map[string]interface{}{
				"type": "expression",
				"constraints": map[string]interface{}{
					"conditions": []interface{}{},
				},
			},
			expectedCount: 0,
			expectError:   false,
			description:   "Empty conditions array should return empty result",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			cp := NewConstraintPipeline()
			// Execute
			conditions, err := cp.extractExistsConditions(tt.exprMap)
			// Assert
			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				require.NoError(t, err, tt.description)
				assert.Equal(t, tt.expectedCount, len(conditions), "Number of extracted conditions should match")
			}
		})
	}
}
// TestConstraintPipeline_connectExistsNodeToTypeNodes tests the connectExistsNodeToTypeNodes wrapper
func TestConstraintPipeline_connectExistsNodeToTypeNodes(t *testing.T) {
	tests := []struct {
		name            string
		mainVarType     string
		existsVarType   string
		createTypeNodes bool
		description     string
	}{
		{
			name:            "connect with both type nodes present",
			mainVarType:     "Person",
			existsVarType:   "Order",
			createTypeNodes: true,
			description:     "Should connect exists node to both type nodes",
		},
		{
			name:            "connect with type nodes missing",
			mainVarType:     "Customer",
			existsVarType:   "Transaction",
			createTypeNodes: false,
			description:     "Should handle missing type nodes gracefully",
		},
		{
			name:            "connect with same type for main and exists",
			mainVarType:     "Node",
			existsVarType:   "Node",
			createTypeNodes: true,
			description:     "Should handle same type for both variables",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			storage := NewMemoryStorage()
			cp := NewConstraintPipeline()
			network := NewReteNetwork(storage)
			if tt.createTypeNodes {
				mainTypeDef := TypeDefinition{
					Type:   "typeDefinition",
					Name:   tt.mainVarType,
					Fields: []Field{},
				}
				network.TypeNodes[tt.mainVarType] = NewTypeNode(tt.mainVarType, mainTypeDef, storage)
				existsTypeDef := TypeDefinition{
					Type:   "typeDefinition",
					Name:   tt.existsVarType,
					Fields: []Field{},
				}
				network.TypeNodes[tt.existsVarType] = NewTypeNode(tt.existsVarType, existsTypeDef, storage)
			}
			existsNode := NewExistsNode("test_exists", nil, "mainVar", "existsVar", map[string]string{}, storage)
			// Execute - should not panic
			require.NotPanics(t, func() {
				cp.connectExistsNodeToTypeNodes(
					network,
					"test_rule",
					existsNode,
					"mainVar",
					tt.mainVarType,
					"existsVar",
					tt.existsVarType,
				)
			}, tt.description)
			// Verify exists node was created
			assert.NotNil(t, existsNode)
		})
	}
}
// TestConstraintPipeline_createAccumulatorRule tests the createAccumulatorRule wrapper
func TestConstraintPipeline_createAccumulatorRule(t *testing.T) {
	tests := []struct {
		name          string
		variables     []map[string]interface{}
		variableNames []string
		variableTypes []string
		aggInfo       *AggregationInfo
		action        *Action
		expectError   bool
		description   string
	}{
		{
			name: "simple count accumulator",
			variables: []map[string]interface{}{
				{
					"name": "order",
					"type": "Order",
				},
			},
			variableNames: []string{"order"},
			variableTypes: []string{"Order"},
			aggInfo: &AggregationInfo{
				Function:    "count",
				AggVariable: "order",
			},
			action: &Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "orderCount",
					Args: []interface{}{"Order count updated"},
				},
			},
			expectError: false,
			description: "Simple count accumulator should succeed",
		},
		{
			name: "sum accumulator with field",
			variables: []map[string]interface{}{
				{
					"name": "transaction",
					"type": "Transaction",
				},
			},
			variableNames: []string{"transaction"},
			variableTypes: []string{"Transaction"},
			aggInfo: &AggregationInfo{
				Function:    "sum",
				AggVariable: "transaction",
				Field:       "amount",
			},
			action: &Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "totalAmount",
					Args: []interface{}{"Total transaction amount"},
				},
			},
			expectError: false,
			description: "Sum accumulator with field should succeed",
		},
		{
			name: "average accumulator",
			variables: []map[string]interface{}{
				{
					"name": "rating",
					"type": "Rating",
				},
			},
			variableNames: []string{"rating"},
			variableTypes: []string{"Rating"},
			aggInfo: &AggregationInfo{
				Function:    "average",
				AggVariable: "rating",
				Field:       "score",
			},
			action: &Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "avgRating",
					Args: []interface{}{"Average rating calculated"},
				},
			},
			expectError: false,
			description: "Average accumulator should succeed",
		},
		{
			name: "accumulator with nil action",
			variables: []map[string]interface{}{
				{
					"name": "item",
					"type": "Item",
				},
			},
			variableNames: []string{"item"},
			variableTypes: []string{"Item"},
			aggInfo: &AggregationInfo{
				Function:    "count",
				AggVariable: "item",
			},
			action:      nil,
			expectError: false,
			description: "Accumulator with nil action should succeed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			storage := NewMemoryStorage()
			cp := NewConstraintPipeline()
			network := NewReteNetwork(storage)
			// Create type nodes
			for _, varType := range tt.variableTypes {
				if varType != "" {
					typeDef := TypeDefinition{
						Type:   "typeDefinition",
						Name:   varType,
						Fields: []Field{},
					}
					typeNode := NewTypeNode(varType, typeDef, storage)
					network.TypeNodes[varType] = typeNode
				}
			}
			// Execute
			err := cp.createAccumulatorRule(
				network,
				"test_accumulator_rule",
				tt.variables,
				tt.variableNames,
				tt.variableTypes,
				tt.aggInfo,
				tt.action,
				storage,
			)
			// Assert
			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}
// TestConstraintPipeline_createMultiSourceAccumulatorRule tests the createMultiSourceAccumulatorRule wrapper
func TestConstraintPipeline_createMultiSourceAccumulatorRule(t *testing.T) {
	tests := []struct {
		name        string
		aggInfo     *AggregationInfo
		action      *Action
		expectError bool
		description string
	}{
		{
			name: "multi-source count aggregation",
			aggInfo: &AggregationInfo{
				MainVariable: "m",
				MainType:     "Main",
				SourcePatterns: []SourcePattern{
					{
						Variable: "order",
						Type:     "Order",
					},
					{
						Variable: "invoice",
						Type:     "Invoice",
					},
				},
				JoinConditions: []JoinCondition{
					{
						LeftVar:    "m",
						LeftField:  "id",
						RightVar:   "order",
						RightField: "mainId",
						Operator:   "==",
					},
					{
						LeftVar:    "m",
						LeftField:  "id",
						RightVar:   "invoice",
						RightField: "mainId",
						Operator:   "==",
					},
				},
				AggregationVars: []AggregationVariable{
					{
						Name:      "count",
						Function:  "COUNT",
						SourceVar: "order",
					},
				},
			},
			action: &Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "multiSourceCount",
					Args: []interface{}{"Multi-source count"},
				},
			},
			expectError: false,
			description: "Multi-source count aggregation should succeed",
		},
		{
			name: "multi-source sum aggregation",
			aggInfo: &AggregationInfo{
				MainVariable: "account",
				MainType:     "Account",
				SourcePatterns: []SourcePattern{
					{
						Variable: "payment",
						Type:     "Payment",
					},
					{
						Variable: "refund",
						Type:     "Refund",
					},
				},
				JoinConditions: []JoinCondition{
					{
						LeftVar:    "account",
						LeftField:  "id",
						RightVar:   "payment",
						RightField: "accountId",
						Operator:   "==",
					},
					{
						LeftVar:    "account",
						LeftField:  "id",
						RightVar:   "refund",
						RightField: "accountId",
						Operator:   "==",
					},
				},
				AggregationVars: []AggregationVariable{
					{
						Name:      "total",
						Function:  "SUM",
						SourceVar: "payment",
						Field:     "amount",
					},
				},
			},
			action: &Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "multiSourceSum",
					Args: []interface{}{"Multi-source sum"},
				},
			},
			expectError: false,
			description: "Multi-source sum aggregation should succeed",
		},
		{
			name: "multi-source with three sources",
			aggInfo: &AggregationInfo{
				MainVariable: "user",
				MainType:     "User",
				SourcePatterns: []SourcePattern{
					{
						Variable: "click",
						Type:     "ClickEvent",
					},
					{
						Variable: "view",
						Type:     "ViewEvent",
					},
					{
						Variable: "purchase",
						Type:     "PurchaseEvent",
					},
				},
				JoinConditions: []JoinCondition{
					{
						LeftVar:    "user",
						LeftField:  "id",
						RightVar:   "click",
						RightField: "userId",
						Operator:   "==",
					},
					{
						LeftVar:    "user",
						LeftField:  "id",
						RightVar:   "view",
						RightField: "userId",
						Operator:   "==",
					},
					{
						LeftVar:    "user",
						LeftField:  "id",
						RightVar:   "purchase",
						RightField: "userId",
						Operator:   "==",
					},
				},
				AggregationVars: []AggregationVariable{
					{
						Name:      "events",
						Function:  "COUNT",
						SourceVar: "click",
					},
				},
			},
			action: &Action{
				Type: "action",
				Job: &JobCall{
					Type: "jobCall",
					Name: "threeSourceAggregation",
					Args: []interface{}{"Three-source aggregation"},
				},
			},
			expectError: false,
			description: "Multi-source with three sources should succeed",
		},
		{
			name: "multi-source with nil action",
			aggInfo: &AggregationInfo{
				MainVariable: "device",
				MainType:     "Device",
				SourcePatterns: []SourcePattern{
					{
						Variable: "sensor1",
						Type:     "Sensor",
					},
					{
						Variable: "sensor2",
						Type:     "Sensor",
					},
				},
				JoinConditions: []JoinCondition{
					{
						LeftVar:    "device",
						LeftField:  "id",
						RightVar:   "sensor1",
						RightField: "deviceId",
						Operator:   "==",
					},
					{
						LeftVar:    "device",
						LeftField:  "id",
						RightVar:   "sensor2",
						RightField: "deviceId",
						Operator:   "==",
					},
				},
				AggregationVars: []AggregationVariable{
					{
						Name:      "avg",
						Function:  "AVG",
						SourceVar: "sensor1",
						Field:     "value",
					},
				},
			},
			action:      nil,
			expectError: false,
			description: "Multi-source with nil action should succeed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			storage := NewMemoryStorage()
			cp := NewConstraintPipeline()
			network := NewReteNetwork(storage)
			// Create type nodes for all sources
			if tt.aggInfo != nil {
				// Create main type node
				if tt.aggInfo.MainType != "" {
					mainTypeDef := TypeDefinition{
						Type:   "typeDefinition",
						Name:   tt.aggInfo.MainType,
						Fields: []Field{},
					}
					mainTypeNode := NewTypeNode(tt.aggInfo.MainType, mainTypeDef, storage)
					network.TypeNodes[tt.aggInfo.MainType] = mainTypeNode
				}
				// Create source type nodes
				if tt.aggInfo.SourcePatterns != nil {
					for _, source := range tt.aggInfo.SourcePatterns {
						if source.Type != "" {
							typeDef := TypeDefinition{
								Type:   "typeDefinition",
								Name:   source.Type,
								Fields: []Field{},
							}
							typeNode := NewTypeNode(source.Type, typeDef, storage)
							network.TypeNodes[source.Type] = typeNode
						}
					}
				}
			}
			// Execute
			err := cp.createMultiSourceAccumulatorRule(
				network,
				"test_multi_source_rule",
				tt.aggInfo,
				tt.action,
				storage,
			)
			// Assert
			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}
// TestConstraintPipeline_createPassthroughAlphaNode tests the createPassthroughAlphaNode wrapper
func TestConstraintPipeline_createPassthroughAlphaNode(t *testing.T) {
	tests := []struct {
		name        string
		ruleID      string
		varName     string
		side        string
		description string
	}{
		{
			name:        "create passthrough for left side",
			ruleID:      "rule_123",
			varName:     "person",
			side:        "left",
			description: "Should create passthrough alpha node for left side",
		},
		{
			name:        "create passthrough for right side",
			ruleID:      "rule_456",
			varName:     "order",
			side:        "right",
			description: "Should create passthrough alpha node for right side",
		},
		{
			name:        "create passthrough with no side specified",
			ruleID:      "rule_789",
			varName:     "item",
			side:        "",
			description: "Should create passthrough alpha node with no side",
		},
		{
			name:        "create passthrough with complex variable name",
			ruleID:      "rule_complex",
			varName:     "customer_account",
			side:        "left",
			description: "Should handle complex variable names",
		},
		{
			name:        "create passthrough with empty rule ID",
			ruleID:      "",
			varName:     "entity",
			side:        "left",
			description: "Should handle empty rule ID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			storage := NewMemoryStorage()
			cp := NewConstraintPipeline()
			// Execute
			alphaNode := cp.createPassthroughAlphaNode(tt.ruleID, tt.varName, tt.side, storage)
			// Assert
			assert.NotNil(t, alphaNode, tt.description)
			assert.NotEmpty(t, alphaNode.ID, "Alpha node should have an ID")
		})
	}
}
// TestConstraintPipeline_connectTypeNodeToBetaNode tests the connectTypeNodeToBetaNode wrapper
func TestConstraintPipeline_connectTypeNodeToBetaNode(t *testing.T) {
	tests := []struct {
		name           string
		varType        string
		side           string
		createTypeNode bool
		description    string
	}{
		{
			name:           "connect to left side with type node present",
			varType:        "Person",
			side:           "left",
			createTypeNode: true,
			description:    "Should connect type node to beta node on left side",
		},
		{
			name:           "connect to right side with type node present",
			varType:        "Order",
			side:           "right",
			createTypeNode: true,
			description:    "Should connect type node to beta node on right side",
		},
		{
			name:           "connect with type node missing",
			varType:        "MissingType",
			side:           "left",
			createTypeNode: false,
			description:    "Should handle missing type node gracefully",
		},
		{
			name:           "connect with no side specified",
			varType:        "Item",
			side:           "",
			createTypeNode: true,
			description:    "Should handle connection with no side specified",
		},
		{
			name:           "connect multiple times to same type",
			varType:        "Customer",
			side:           "left",
			createTypeNode: true,
			description:    "Should handle multiple connections to same type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			storage := NewMemoryStorage()
			cp := NewConstraintPipeline()
			network := NewReteNetwork(storage)
			if tt.createTypeNode {
				typeDef := TypeDefinition{
					Type:   "typeDefinition",
					Name:   tt.varType,
					Fields: []Field{},
				}
				network.TypeNodes[tt.varType] = NewTypeNode(tt.varType, typeDef, storage)
			}
			// Create a beta node
			joinNode := NewJoinNode("test_join", nil, []string{}, []string{}, map[string]string{}, storage)
			// Execute - should not panic
			require.NotPanics(t, func() {
				cp.connectTypeNodeToBetaNode(
					network,
					"test_rule",
					"varName",
					tt.varType,
					joinNode,
					tt.side,
				)
			}, tt.description)
			// Verify join node was created
			assert.NotNil(t, joinNode)
		})
	}
}
// TestConstraintPipeline_getVariableInfo tests the getVariableInfo helper
func TestConstraintPipeline_getVariableInfo(t *testing.T) {
	tests := []struct {
		name            string
		variables       []map[string]interface{}
		variableTypes   []string
		expectedVarName string
		expectedVarType string
		description     string
	}{
		{
			name: "standard variable with name and type",
			variables: []map[string]interface{}{
				{
					"name": "person",
					"type": "Person",
				},
			},
			variableTypes:   []string{"Person"},
			expectedVarName: "person",
			expectedVarType: "Person",
			description:     "Should extract variable name and type correctly",
		},
		{
			name: "variable with only name",
			variables: []map[string]interface{}{
				{
					"name": "order",
				},
			},
			variableTypes:   []string{"Order"},
			expectedVarName: "order",
			expectedVarType: "Order",
			description:     "Should get type from variableTypes when not in variable map",
		},
		{
			name:            "empty variables array",
			variables:       []map[string]interface{}{},
			variableTypes:   []string{},
			expectedVarName: "p",
			expectedVarType: "",
			description:     "Should return default 'p' for empty variables",
		},
		{
			name:            "nil variables",
			variables:       nil,
			variableTypes:   nil,
			expectedVarName: "p",
			expectedVarType: "",
			description:     "Should return default 'p' for nil variables",
		},
		{
			name: "variable without name field",
			variables: []map[string]interface{}{
				{
					"type": "Customer",
				},
			},
			variableTypes:   []string{"Customer"},
			expectedVarName: "p",
			expectedVarType: "Customer",
			description:     "Should use default 'p' when name is not present",
		},
		{
			name: "multiple variables uses first",
			variables: []map[string]interface{}{
				{
					"name": "first",
					"type": "First",
				},
				{
					"name": "second",
					"type": "Second",
				},
			},
			variableTypes:   []string{"First", "Second"},
			expectedVarName: "first",
			expectedVarType: "First",
			description:     "Should use the first variable when multiple are present",
		},
		{
			name: "variable with non-string name",
			variables: []map[string]interface{}{
				{
					"name": 123,
					"type": "Item",
				},
			},
			variableTypes:   []string{"Item"},
			expectedVarName: "p",
			expectedVarType: "Item",
			description:     "Should use default 'p' when name is not a string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			cp := NewConstraintPipeline()
			// Execute
			varName, varType := cp.getVariableInfo(tt.variables, tt.variableTypes)
			// Assert
			assert.Equal(t, tt.expectedVarName, varName, "Variable name should match - "+tt.description)
			assert.Equal(t, tt.expectedVarType, varType, "Variable type should match - "+tt.description)
		})
	}
}