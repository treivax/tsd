// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// extractJoinConditionsRecursive recursively extracts all join conditions from constraint tree
// This function traverses the constraint AST and identifies comparison operations that
// represent join conditions (as opposed to threshold conditions on aggregation variables).
func (cp *ConstraintPipeline) extractJoinConditionsRecursive(constraints map[string]interface{}, aggVarNames map[string]bool, joinConditions *[]JoinCondition) {
	constraintType, _ := constraints["type"].(string)

	if constraintType == "comparison" {
		// Check if this is a join condition (not a threshold)
		if !cp.isThresholdCondition(constraints, aggVarNames) {
			// Extract join condition details
			joinCond := JoinCondition{}

			if operator, ok := constraints["operator"].(string); ok {
				joinCond.Operator = operator
			}

			if leftData, ok := constraints["left"].(map[string]interface{}); ok {
				if leftData["type"] == "fieldAccess" {
					if obj, ok := leftData["object"].(string); ok {
						joinCond.LeftVar = obj
					}
					if field, ok := leftData["field"].(string); ok {
						joinCond.LeftField = field
					}
				}
			}

			if rightData, ok := constraints["right"].(map[string]interface{}); ok {
				if rightData["type"] == "fieldAccess" {
					if obj, ok := rightData["object"].(string); ok {
						joinCond.RightVar = obj
					}
					if field, ok := rightData["field"].(string); ok {
						joinCond.RightField = field
					}
				}
			}

			*joinConditions = append(*joinConditions, joinCond)
		}
	} else if constraintType == "logicalExpr" {
		// Recursively process left side
		if leftData, ok := constraints["left"].(map[string]interface{}); ok {
			cp.extractJoinConditionsRecursive(leftData, aggVarNames, joinConditions)
		}

		// Recursively process operations
		if ops, ok := constraints["operations"].([]interface{}); ok {
			for _, opInterface := range ops {
				if opMap, ok := opInterface.(map[string]interface{}); ok {
					if rightData, ok := opMap["right"].(map[string]interface{}); ok {
						cp.extractJoinConditionsRecursive(rightData, aggVarNames, joinConditions)
					}
				}
			}
		}
	}
}

// separateAggregationConstraints separates join conditions from threshold conditions
// Returns (joinConditions, thresholdConditions)
//
// Join conditions are comparisons between fields of different variables (e.g., p.id == e.id)
// Threshold conditions are comparisons involving aggregation variables (e.g., totalSales > 1000)
func (cp *ConstraintPipeline) separateAggregationConstraints(constraints map[string]interface{}, aggVarNames map[string]bool) (map[string]interface{}, []map[string]interface{}) {
	var joinConditions map[string]interface{}
	var thresholdConditions []map[string]interface{}

	constraintType, _ := constraints["type"].(string)

	if constraintType == "comparison" {
		// Single comparison - check if it's a threshold or join
		if cp.isThresholdCondition(constraints, aggVarNames) {
			thresholdConditions = append(thresholdConditions, constraints)
		} else {
			joinConditions = constraints
		}
	} else if constraintType == "logicalExpr" {
		// Logical expression (AND/OR) - separate conditions
		leftData, _ := constraints["left"]

		// Handle operations field - it might be []interface{} or []map[string]interface{}
		var operations []interface{}
		if ops, ok := constraints["operations"].([]interface{}); ok {
			operations = ops
		} else if ops, ok := constraints["operations"].([]map[string]interface{}); ok {
			// Convert to []interface{}
			for _, op := range ops {
				operations = append(operations, op)
			}
		}

		// Check left condition
		if leftMap, ok := leftData.(map[string]interface{}); ok {
			if cp.isThresholdCondition(leftMap, aggVarNames) {
				thresholdConditions = append(thresholdConditions, leftMap)
			} else {
				joinConditions = leftMap
			}
		}

		// Check operations
		for _, opInterface := range operations {
			if opMap, ok := opInterface.(map[string]interface{}); ok {
				if rightData, ok := opMap["right"].(map[string]interface{}); ok {
					if cp.isThresholdCondition(rightData, aggVarNames) {
						thresholdConditions = append(thresholdConditions, rightData)
					} else if joinConditions == nil {
						joinConditions = rightData
					}
				}
			}
		}
	}

	return joinConditions, thresholdConditions
}

// isThresholdCondition checks if a comparison references an aggregation variable
// A threshold condition involves at least one aggregation variable (e.g., totalSales > 1000)
// whereas a join condition involves only regular variables (e.g., p.id == e.id)
func (cp *ConstraintPipeline) isThresholdCondition(condition map[string]interface{}, aggVarNames map[string]bool) bool {
	condType, _ := condition["type"].(string)
	if condType != "comparison" {
		return false
	}

	// Check if left side is an aggregation variable
	if leftData, ok := condition["left"].(map[string]interface{}); ok {
		if leftData["type"] == "variable" {
			if varName, ok := leftData["name"].(string); ok {
				return aggVarNames[varName]
			}
		}
	}

	// Check if right side is an aggregation variable
	if rightData, ok := condition["right"].(map[string]interface{}); ok {
		if rightData["type"] == "variable" {
			if varName, ok := rightData["name"].(string); ok {
				return aggVarNames[varName]
			}
		}
	}

	return false
}
