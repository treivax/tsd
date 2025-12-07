// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// TestParseAggregationExpression tests the parseAggregationExpression function
func TestParseAggregationExpression(t *testing.T) {
	cp := &ConstraintPipeline{}

	tests := []struct {
		name        string
		exprMap     map[string]interface{}
		expectError bool
	}{
		{
			name: "valid expression with patterns",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{
						"variables": []interface{}{
							map[string]interface{}{
								"type": "aggregationVariable",
								"name": "avg_sal",
							},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name:        "missing patterns field",
			exprMap:     map[string]interface{}{},
			expectError: true,
		},
		{
			name: "empty patterns list",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{},
			},
			expectError: true,
		},
		{
			name: "patterns not a list",
			exprMap: map[string]interface{}{
				"patterns": "invalid",
			},
			expectError: true,
		},
		{
			name: "missing variables in pattern",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, varsList, err := cp.parseAggregationExpression(tt.exprMap)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if varsList == nil {
					t.Errorf("expected variables list but got nil")
				}
			}
		})
	}
}

// TestExtractAggregationFunction tests the extractAggregationFunction function
func TestExtractAggregationFunction(t *testing.T) {
	cp := &ConstraintPipeline{}

	tests := []struct {
		name         string
		varMap       map[string]interface{}
		expectedFunc string
		expectError  bool
	}{
		{
			name: "direct function field",
			varMap: map[string]interface{}{
				"function": "AVG",
			},
			expectedFunc: "AVG",
			expectError:  false,
		},
		{
			name: "nested function in value",
			varMap: map[string]interface{}{
				"value": map[string]interface{}{
					"type":     "functionCall",
					"function": "SUM",
				},
			},
			expectedFunc: "SUM",
			expectError:  false,
		},
		{
			name: "nested aggregationCall type",
			varMap: map[string]interface{}{
				"value": map[string]interface{}{
					"type":     "aggregationCall",
					"function": "COUNT",
				},
			},
			expectedFunc: "COUNT",
			expectError:  false,
		},
		{
			name:        "no function found",
			varMap:      map[string]interface{}{},
			expectError: true,
		},
		{
			name: "invalid nested structure",
			varMap: map[string]interface{}{
				"value": map[string]interface{}{
					"type": "other",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			function, err := cp.extractAggregationFunction(tt.varMap)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if function != tt.expectedFunc {
					t.Errorf("expected function %q but got %q", tt.expectedFunc, function)
				}
			}
		})
	}
}

// TestExtractAggregationField tests the extractAggregationField function
func TestExtractAggregationField(t *testing.T) {
	cp := &ConstraintPipeline{}

	tests := []struct {
		name        string
		varMap      map[string]interface{}
		expectedVar string
		expectedFld string
		expectError bool
	}{
		{
			name: "direct field format",
			varMap: map[string]interface{}{
				"field": map[string]interface{}{
					"object": "e",
					"field":  "salary",
				},
			},
			expectedVar: "e",
			expectedFld: "salary",
			expectError: false,
		},
		{
			name: "nested in value arguments",
			varMap: map[string]interface{}{
				"value": map[string]interface{}{
					"arguments": []interface{}{
						map[string]interface{}{
							"type":   "fieldAccess",
							"object": "emp",
							"field":  "age",
						},
					},
				},
			},
			expectedVar: "emp",
			expectedFld: "age",
			expectError: false,
		},
		{
			name:        "no field information",
			varMap:      map[string]interface{}{},
			expectError: true,
		},
		{
			name: "empty arguments list",
			varMap: map[string]interface{}{
				"value": map[string]interface{}{
					"arguments": []interface{}{},
				},
			},
			expectError: true,
		},
		{
			name: "wrong argument type",
			varMap: map[string]interface{}{
				"value": map[string]interface{}{
					"arguments": []interface{}{
						map[string]interface{}{
							"type": "literal",
						},
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggVar, field, err := cp.extractAggregationField(tt.varMap)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if aggVar != tt.expectedVar {
					t.Errorf("expected variable %q but got %q", tt.expectedVar, aggVar)
				}
				if field != tt.expectedFld {
					t.Errorf("expected field %q but got %q", tt.expectedFld, field)
				}
			}
		})
	}
}

// TestExtractSourceType tests the extractSourceType function
func TestExtractSourceType(t *testing.T) {
	cp := &ConstraintPipeline{}

	tests := []struct {
		name         string
		exprMap      map[string]interface{}
		expectedType string
		expectError  bool
	}{
		{
			name: "valid source type",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{},
					map[string]interface{}{
						"variables": []interface{}{
							map[string]interface{}{
								"dataType": "Employee",
							},
						},
					},
				},
			},
			expectedType: "Employee",
			expectError:  false,
		},
		{
			name: "missing second pattern",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{},
				},
			},
			expectError: true,
		},
		{
			name: "missing variables in second pattern",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{},
					map[string]interface{}{},
				},
			},
			expectError: true,
		},
		{
			name: "missing dataType",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{},
					map[string]interface{}{
						"variables": []interface{}{
							map[string]interface{}{},
						},
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggType, err := cp.extractSourceType(tt.exprMap)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if aggType != tt.expectedType {
					t.Errorf("expected type %q but got %q", tt.expectedType, aggType)
				}
			}
		})
	}
}

// TestExtractJoinFields tests the extractJoinFields function
func TestExtractJoinFields(t *testing.T) {
	cp := &ConstraintPipeline{}

	tests := []struct {
		name            string
		joinConditions  map[string]interface{}
		expectedJoinFld string
		expectedMainFld string
	}{
		{
			name: "valid comparison with field access",
			joinConditions: map[string]interface{}{
				"type": "comparison",
				"left": map[string]interface{}{
					"type":  "fieldAccess",
					"field": "deptId",
				},
				"right": map[string]interface{}{
					"type":  "fieldAccess",
					"field": "id",
				},
			},
			expectedJoinFld: "deptId",
			expectedMainFld: "id",
		},
		{
			name: "not a comparison type",
			joinConditions: map[string]interface{}{
				"type": "logical",
			},
			expectedJoinFld: "",
			expectedMainFld: "",
		},
		{
			name: "missing field access",
			joinConditions: map[string]interface{}{
				"type": "comparison",
				"left": map[string]interface{}{
					"type": "literal",
				},
			},
			expectedJoinFld: "",
			expectedMainFld: "",
		},
		{
			name:            "empty conditions",
			joinConditions:  map[string]interface{}{},
			expectedJoinFld: "",
			expectedMainFld: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			joinField, mainField := cp.extractJoinFields(tt.joinConditions)
			if joinField != tt.expectedJoinFld {
				t.Errorf("expected join field %q but got %q", tt.expectedJoinFld, joinField)
			}
			if mainField != tt.expectedMainFld {
				t.Errorf("expected main field %q but got %q", tt.expectedMainFld, mainField)
			}
		})
	}
}

// TestExtractThresholdConditions tests the extractThresholdConditions function
func TestExtractThresholdConditions(t *testing.T) {
	cp := &ConstraintPipeline{}

	tests := []struct {
		name              string
		thresholdConds    []map[string]interface{}
		expectedOperator  string
		expectedThreshold float64
	}{
		{
			name: "valid threshold condition",
			thresholdConds: []map[string]interface{}{
				{
					"operator": ">",
					"right": map[string]interface{}{
						"value": 50000.0,
					},
				},
			},
			expectedOperator:  ">",
			expectedThreshold: 50000.0,
		},
		{
			name: "multiple conditions - uses first",
			thresholdConds: []map[string]interface{}{
				{
					"operator": ">=",
					"right": map[string]interface{}{
						"value": 1000.0,
					},
				},
				{
					"operator": "<",
					"right": map[string]interface{}{
						"value": 5000.0,
					},
				},
			},
			expectedOperator:  ">=",
			expectedThreshold: 1000.0,
		},
		{
			name:              "empty conditions - use defaults",
			thresholdConds:    []map[string]interface{}{},
			expectedOperator:  DefaultThresholdOperator,
			expectedThreshold: DefaultThresholdValue,
		},
		{
			name: "missing operator - use default",
			thresholdConds: []map[string]interface{}{
				{
					"right": map[string]interface{}{
						"value": 100.0,
					},
				},
			},
			expectedOperator:  DefaultThresholdOperator,
			expectedThreshold: 100.0,
		},
		{
			name: "missing value - use default",
			thresholdConds: []map[string]interface{}{
				{
					"operator": "<=",
				},
			},
			expectedOperator:  "<=",
			expectedThreshold: DefaultThresholdValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operator, threshold := cp.extractThresholdConditions(tt.thresholdConds)
			if operator != tt.expectedOperator {
				t.Errorf("expected operator %q but got %q", tt.expectedOperator, operator)
			}
			if threshold != tt.expectedThreshold {
				t.Errorf("expected threshold %f but got %f", tt.expectedThreshold, threshold)
			}
		})
	}
}

// TestExtractAggregationInfoFromVariables_Integration tests the refactored orchestrator function
func TestExtractAggregationInfoFromVariables_Integration(t *testing.T) {
	cp := &ConstraintPipeline{}

	tests := []struct {
		name        string
		exprMap     map[string]interface{}
		expectError bool
		validate    func(*testing.T, *AggregationInfo)
	}{
		{
			name: "complete aggregation info",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{
						"variables": []interface{}{
							map[string]interface{}{
								"type":     "aggregationVariable",
								"name":     "avg_sal",
								"function": "AVG",
								"field": map[string]interface{}{
									"object": "e",
									"field":  "salary",
								},
							},
						},
					},
					map[string]interface{}{
						"variables": []interface{}{
							map[string]interface{}{
								"dataType": "Employee",
							},
						},
					},
				},
				"constraints": map[string]interface{}{
					"type": "comparison",
					"left": map[string]interface{}{
						"type":  "fieldAccess",
						"field": "deptId",
					},
					"right": map[string]interface{}{
						"type":  "fieldAccess",
						"field": "id",
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, info *AggregationInfo) {
				if info.Function != "AVG" {
					t.Errorf("expected function AVG but got %s", info.Function)
				}
				if info.AggVariable != "e" {
					t.Errorf("expected variable e but got %s", info.AggVariable)
				}
				if info.Field != "salary" {
					t.Errorf("expected field salary but got %s", info.Field)
				}
				if info.AggType != "Employee" {
					t.Errorf("expected type Employee but got %s", info.AggType)
				}
			},
		},
		{
			name: "no constraints - use defaults",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{
						"variables": []interface{}{
							map[string]interface{}{
								"type":     "aggregationVariable",
								"function": "SUM",
								"field": map[string]interface{}{
									"object": "x",
									"field":  "amount",
								},
							},
						},
					},
					map[string]interface{}{
						"variables": []interface{}{
							map[string]interface{}{
								"dataType": "Transaction",
							},
						},
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, info *AggregationInfo) {
				if info.Operator != DefaultThresholdOperator {
					t.Errorf("expected default operator %s but got %s", DefaultThresholdOperator, info.Operator)
				}
				if info.Threshold != DefaultThresholdValue {
					t.Errorf("expected default threshold %f but got %f", DefaultThresholdValue, info.Threshold)
				}
			},
		},
		{
			name: "missing patterns",
			exprMap: map[string]interface{}{
				"constraints": map[string]interface{}{},
			},
			expectError: true,
		},
		{
			name: "no aggregation variable",
			exprMap: map[string]interface{}{
				"patterns": []interface{}{
					map[string]interface{}{
						"variables": []interface{}{
							map[string]interface{}{
								"type": "regular",
								"name": "d",
							},
						},
					},
					map[string]interface{}{},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := cp.extractAggregationInfoFromVariables(tt.exprMap)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if info == nil {
					t.Fatalf("expected aggregation info but got nil")
				}
				if tt.validate != nil {
					tt.validate(t, info)
				}
			}
		})
	}
}
