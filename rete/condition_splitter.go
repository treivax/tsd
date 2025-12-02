// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

// ConditionType indicates whether a condition is Alpha or Beta
type ConditionType int

const (
	// ConditionTypeAlpha - test on a single variable (unary filter)
	ConditionTypeAlpha ConditionType = iota
	// ConditionTypeBeta - test between multiple variables (binary join)
	ConditionTypeBeta
)

// SplitCondition represents a decomposed condition with metadata
type SplitCondition struct {
	Type      ConditionType          // Alpha or Beta
	Condition map[string]interface{} // The condition itself
	Variables []string               // Variables used in this condition
	Variable  string                 // Primary variable (for alpha conditions)
}

// ConditionSplitter decomposes complex conditions into alpha and beta conditions
type ConditionSplitter struct{}

// NewConditionSplitter creates a new ConditionSplitter instance
func NewConditionSplitter() *ConditionSplitter {
	return &ConditionSplitter{}
}

// SplitConditions decomposes a condition (potentially with AND operators) into alpha and beta conditions
// Alpha conditions: tests on a single variable (should be evaluated in AlphaNodes)
// Beta conditions: tests between multiple variables (should be evaluated in JoinNodes)
//
// NOTE: Currently only simple alpha conditions (direct field comparisons) are extracted.
// Complex arithmetic expressions are kept in beta conditions until AlphaConditionEvaluator
// is enhanced to support them.
func (cs *ConditionSplitter) SplitConditions(
	condition map[string]interface{},
) (alphaConditions []SplitCondition, betaConditions []SplitCondition, err error) {

	if condition == nil {
		return nil, nil, nil
	}

	// Check if this is a constraint wrapper (from TSD parser)
	if constraint, ok := condition["constraint"].(map[string]interface{}); ok {
		condition = constraint
	}

	// Check if this is a type wrapper
	if condType, ok := condition["type"].(string); ok && condType == "constraint" {
		if innerConstraint, ok := condition["constraint"].(map[string]interface{}); ok {
			condition = innerConstraint
		}
	}

	// Handle logicalExpr (AND operations)
	if exprType, ok := condition["type"].(string); ok && exprType == "logicalExpr" {
		return cs.splitLogicalExpression(condition)
	}

	// Single condition - classify it
	vars := cs.extractVariables(condition)
	splitCond := SplitCondition{
		Condition: condition,
		Variables: vars,
	}

	if len(vars) == 1 {
		// Alpha condition (single variable)
		// Check if this is a SIMPLE condition that AlphaNode can evaluate
		if cs.isSimpleAlphaCondition(condition) {
			splitCond.Type = ConditionTypeAlpha
			splitCond.Variable = vars[0]
			alphaConditions = append(alphaConditions, splitCond)
		} else {
			// Complex arithmetic expression - keep in beta for now
			// TODO: Enhance AlphaConditionEvaluator to handle arithmetic
			splitCond.Type = ConditionTypeBeta
			betaConditions = append(betaConditions, splitCond)
		}
	} else if len(vars) > 1 {
		// Beta condition (multiple variables)
		splitCond.Type = ConditionTypeBeta
		betaConditions = append(betaConditions, splitCond)
	}

	return alphaConditions, betaConditions, nil
}

// splitLogicalExpression splits a logicalExpr (AND operations) into individual conditions
func (cs *ConditionSplitter) splitLogicalExpression(
	logicalExpr map[string]interface{},
) (alphaConditions []SplitCondition, betaConditions []SplitCondition, err error) {

	// Extract the left condition
	left, hasLeft := logicalExpr["left"].(map[string]interface{})
	if hasLeft {
		alphas, betas, err := cs.SplitConditions(left)
		if err != nil {
			return nil, nil, fmt.Errorf("error splitting left condition: %w", err)
		}
		alphaConditions = append(alphaConditions, alphas...)
		betaConditions = append(betaConditions, betas...)
	}

	// Extract the operations (right-hand side AND operations)
	// Try to extract operations - handle different possible types
	var operations []interface{}
	hasOps := false

	if opsRaw, exists := logicalExpr["operations"]; exists && opsRaw != nil {
		// Try []interface{} first
		if opsSlice, ok := opsRaw.([]interface{}); ok {
			operations = opsSlice
			hasOps = true
		} else if opsSlice, ok := opsRaw.([]map[string]interface{}); ok {
			// Convert []map[string]interface{} to []interface{}
			operations = make([]interface{}, len(opsSlice))
			for i, op := range opsSlice {
				operations[i] = op
			}
			hasOps = true
		}
	}

	if hasOps {
		for _, op := range operations {
			opMap, ok := op.(map[string]interface{})
			if !ok {
				continue
			}

			// Each operation should have an "op" field (AND) and a "right" field
			right, hasRight := opMap["right"].(map[string]interface{})
			if !hasRight {
				continue
			}

			alphas, betas, err := cs.SplitConditions(right)
			if err != nil {
				return nil, nil, fmt.Errorf("error splitting operation: %w", err)
			}
			alphaConditions = append(alphaConditions, alphas...)
			betaConditions = append(betaConditions, betas...)
		}
	}

	return alphaConditions, betaConditions, nil
}

// ClassifyCondition determines if a condition is alpha or beta
func (cs *ConditionSplitter) ClassifyCondition(condition map[string]interface{}) ConditionType {
	vars := cs.extractVariables(condition)
	if len(vars) == 1 {
		return ConditionTypeAlpha
	}
	return ConditionTypeBeta
}

// IsAlphaCondition returns true if the condition is an alpha condition (single variable)
func (cs *ConditionSplitter) IsAlphaCondition(condition map[string]interface{}) bool {
	return cs.ClassifyCondition(condition) == ConditionTypeAlpha
}

// IsBetaCondition returns true if the condition is a beta condition (multiple variables)
func (cs *ConditionSplitter) IsBetaCondition(condition map[string]interface{}) bool {
	return cs.ClassifyCondition(condition) == ConditionTypeBeta
}

// ExtractVariables extracts all variable names from a condition
func (cs *ConditionSplitter) ExtractVariables(condition map[string]interface{}) []string {
	return cs.extractVariables(condition)
}

// extractVariables recursively extracts variable names from a condition AST
func (cs *ConditionSplitter) extractVariables(node interface{}) []string {
	vars := make(map[string]bool)
	cs.extractVariablesRecursive(node, vars)

	// Convert map to slice
	result := make([]string, 0, len(vars))
	for v := range vars {
		result = append(result, v)
	}
	return result
}

// extractVariablesRecursive traverses the AST and collects variable names
func (cs *ConditionSplitter) extractVariablesRecursive(node interface{}, vars map[string]bool) {
	switch v := node.(type) {
	case map[string]interface{}:
		// Check for fieldAccess type (variable.field)
		if nodeType, ok := v["type"].(string); ok {
			if nodeType == "fieldAccess" {
				// Extract the object (variable name)
				if obj, ok := v["object"].(string); ok {
					vars[obj] = true
				}
			}
		}

		// Also check for "variable" field directly (alternative format)
		if variable, ok := v["variable"].(string); ok {
			vars[variable] = true
		}

		// Check for "object" field (another format)
		if object, ok := v["object"].(string); ok {
			vars[object] = true
		}

		// Recursively process all map values
		for _, value := range v {
			cs.extractVariablesRecursive(value, vars)
		}

	case []interface{}:
		// Recursively process all array elements
		for _, elem := range v {
			cs.extractVariablesRecursive(elem, vars)
		}
	}
}

// GetPrimaryVariable returns the primary variable for an alpha condition
// Returns empty string if not an alpha condition
func (cs *ConditionSplitter) GetPrimaryVariable(condition map[string]interface{}) string {
	vars := cs.extractVariables(condition)
	if len(vars) == 1 {
		return vars[0]
	}
	return ""
}

// ReconstructBetaCondition reconstructs a beta-only condition from split conditions
// This is useful when all alpha conditions have been extracted and we need to create
// a condition map containing only the beta conditions
func (cs *ConditionSplitter) ReconstructBetaCondition(betaConditions []SplitCondition) map[string]interface{} {
	if len(betaConditions) == 0 {
		return nil
	}

	if len(betaConditions) == 1 {
		return betaConditions[0].Condition
	}

	// Multiple beta conditions - reconstruct as logicalExpr with AND
	logicalExpr := map[string]interface{}{
		"type": "logicalExpr",
		"left": betaConditions[0].Condition,
	}

	operations := make([]interface{}, 0, len(betaConditions)-1)
	for i := 1; i < len(betaConditions); i++ {
		operation := map[string]interface{}{
			"op":    "AND",
			"right": betaConditions[i].Condition,
		}
		operations = append(operations, operation)
	}

	logicalExpr["operations"] = operations

	return map[string]interface{}{
		"type":       "constraint",
		"constraint": logicalExpr,
	}
}

// HasMixedConditions checks if a condition contains both alpha and beta conditions
func (cs *ConditionSplitter) HasMixedConditions(condition map[string]interface{}) bool {
	alphas, betas, err := cs.SplitConditions(condition)
	if err != nil {
		return false
	}
	return len(alphas) > 0 && len(betas) > 0
}

// CountConditions returns the total number of atomic conditions in a complex condition
func (cs *ConditionSplitter) CountConditions(condition map[string]interface{}) int {
	alphas, betas, err := cs.SplitConditions(condition)
	if err != nil {
		return 0
	}
	return len(alphas) + len(betas)
}

// isSimpleAlphaCondition checks if an alpha condition can be evaluated by AlphaNode
// Returns true for:
//   - Direct field comparisons (e.g., c.qte > 5)
//   - Arithmetic expressions with single variable (e.g., c.qte * 23 - 10 > 0)
//
// Returns false for:
//   - Multi-variable expressions (e.g., c.qte * p.price > 0) - requires join
//   - Complex expressions that AlphaConditionEvaluator cannot handle
func (cs *ConditionSplitter) isSimpleAlphaCondition(condition map[string]interface{}) bool {
	// Check for comparison type
	condType, hasType := condition["type"].(string)
	if !hasType || condType != "comparison" {
		return false
	}

	// Extract all variables from the entire condition (left and right sides)
	vars := cs.extractVariables(condition)

	// If more than one variable, this is a join condition (beta)
	if len(vars) > 1 {
		return false
	}

	// If exactly one variable, this is an alpha condition
	// The AlphaConditionEvaluator can handle:
	// - Simple field access: c.qte > 5
	// - Arithmetic operations: c.qte * 23 - 10 > 0
	// - Nested expressions: (c.qte * 23 - 10) / 2 > 50
	if len(vars) == 1 {
		return true
	}

	// Zero variables means constant comparison (e.g., 5 > 3)
	// This should be evaluated at compile time, but for now treat as alpha
	return false
}
