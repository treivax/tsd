// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import "fmt"

// ============================================================================
// Multi-Source Aggregation Extraction Orchestration
// ============================================================================

// aggregationExtractionContext holds the state for aggregation info extraction
type aggregationExtractionContext struct {
	pipeline *ConstraintPipeline
	exprMap  map[string]interface{}

	// Extracted data
	aggInfo      *AggregationInfo
	patternsList []interface{}
	firstPattern map[string]interface{}
}

// newAggregationExtractionContext creates a new context for aggregation extraction
func newAggregationExtractionContext(
	cp *ConstraintPipeline,
	exprMap map[string]interface{},
) *aggregationExtractionContext {
	return &aggregationExtractionContext{
		pipeline: cp,
		exprMap:  exprMap,
		aggInfo: &AggregationInfo{
			AggregationVars: []AggregationVariable{},
			SourcePatterns:  []SourcePattern{},
			JoinConditions:  []JoinCondition{},
		},
	}
}

// validateAndExtractPatterns validates the patterns field and extracts the patterns list
func (ctx *aggregationExtractionContext) validateAndExtractPatterns() error {
	// Check for multi-pattern syntax
	patternsData, hasPatterns := ctx.exprMap["patterns"]
	if !hasPatterns {
		return fmt.Errorf("no patterns field found")
	}

	patternsList, ok := patternsData.([]interface{})
	if !ok || len(patternsList) < 2 {
		return fmt.Errorf("expected at least 2 pattern blocks for aggregation")
	}

	ctx.patternsList = patternsList
	return nil
}

// extractFirstPattern extracts the first pattern block containing main and aggregation variables
func (ctx *aggregationExtractionContext) extractFirstPattern() error {
	firstPattern, ok := ctx.patternsList[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("first pattern is not a map")
	}

	varsData, hasVars := firstPattern["variables"]
	if !hasVars {
		return fmt.Errorf("no variables in first pattern")
	}

	varsList, ok := varsData.([]interface{})
	if !ok {
		return fmt.Errorf("variables is not a list")
	}

	ctx.firstPattern = firstPattern

	// Extract variables from first pattern
	return ctx.extractVariablesFromFirstPattern(varsList)
}

// extractVariablesFromFirstPattern extracts main variable and aggregation variables
func (ctx *aggregationExtractionContext) extractVariablesFromFirstPattern(varsList []interface{}) error {
	for _, varInterface := range varsList {
		varMap, ok := varInterface.(map[string]interface{})
		if !ok {
			continue
		}

		varType, _ := varMap["type"].(string)
		if varType == "aggregationVariable" {
			// This is an aggregation variable
			if err := ctx.extractAggregationVariable(varMap); err != nil {
				return err
			}
		} else {
			// This is the main variable
			ctx.extractMainVariable(varMap)
		}
	}

	return nil
}

// extractAggregationVariable extracts a single aggregation variable
func (ctx *aggregationExtractionContext) extractAggregationVariable(varMap map[string]interface{}) error {
	aggVar := AggregationVariable{}

	if name, ok := varMap["name"].(string); ok {
		aggVar.Name = name
	}

	if function, ok := varMap["function"].(string); ok {
		aggVar.Function = function
	}

	// Extract field being aggregated
	if fieldData, ok := varMap["field"].(map[string]interface{}); ok {
		if fieldObj, ok := fieldData["object"].(string); ok {
			aggVar.SourceVar = fieldObj
		}
		if fieldName, ok := fieldData["field"].(string); ok {
			aggVar.Field = fieldName
		}
	}

	ctx.aggInfo.AggregationVars = append(ctx.aggInfo.AggregationVars, aggVar)
	return nil
}

// extractMainVariable extracts the main variable from the pattern
func (ctx *aggregationExtractionContext) extractMainVariable(varMap map[string]interface{}) {
	if name, ok := varMap["name"].(string); ok {
		ctx.aggInfo.MainVariable = name
	}
	if dataType, ok := varMap["dataType"].(string); ok {
		ctx.aggInfo.MainType = dataType
	}
}

// extractSourcePatterns extracts source patterns from remaining pattern blocks
func (ctx *aggregationExtractionContext) extractSourcePatterns() error {
	for i := 1; i < len(ctx.patternsList); i++ {
		pattern, ok := ctx.patternsList[i].(map[string]interface{})
		if !ok {
			continue
		}

		if varsData, hasVars := pattern["variables"]; hasVars {
			if varsList, ok := varsData.([]interface{}); ok && len(varsList) > 0 {
				if varMap, ok := varsList[0].(map[string]interface{}); ok {
					sourcePattern := SourcePattern{}

					if varName, ok := varMap["name"].(string); ok {
						sourcePattern.Variable = varName
					}
					if varType, ok := varMap["dataType"].(string); ok {
						sourcePattern.Type = varType
					}

					ctx.aggInfo.SourcePatterns = append(ctx.aggInfo.SourcePatterns, sourcePattern)

					// Also update legacy fields if this is the first/primary aggregation source
					if i == 1 && len(ctx.aggInfo.AggregationVars) > 0 {
						ctx.aggInfo.AggVariable = sourcePattern.Variable
						ctx.aggInfo.AggType = sourcePattern.Type
						ctx.aggInfo.Function = ctx.aggInfo.AggregationVars[0].Function
						ctx.aggInfo.Field = ctx.aggInfo.AggregationVars[0].Field
					}
				}
			}
		}
	}

	return nil
}

// extractConstraintsAndConditions extracts join conditions and thresholds from constraints
func (ctx *aggregationExtractionContext) extractConstraintsAndConditions() error {
	constraintsData, hasConstraints := ctx.exprMap["constraints"]
	if !hasConstraints {
		// No constraints - set default threshold
		ctx.setDefaultThresholds()
		return nil
	}

	constraintMap, ok := constraintsData.(map[string]interface{})
	if !ok {
		// Invalid constraints - set default threshold
		ctx.setDefaultThresholds()
		return nil
	}

	ctx.aggInfo.JoinCondition = constraintMap

	// Get list of aggregation variable names
	aggVarNames := ctx.pipeline.getAggregationVariableNames(ctx.exprMap)

	// Separate join conditions and threshold conditions
	joinConditionsMap, thresholdConditions := ctx.pipeline.separateAggregationConstraints(constraintMap, aggVarNames)

	// Extract all join conditions
	ctx.pipeline.extractJoinConditionsRecursive(constraintMap, aggVarNames, &ctx.aggInfo.JoinConditions)

	// Extract join fields from the first join condition (for backward compatibility)
	ctx.extractJoinFieldsForBackwardCompatibility(joinConditionsMap)

	// Extract thresholds and apply to aggregation variables
	ctx.extractAndApplyThresholds(thresholdConditions)

	// Set default threshold for first aggregation variable (backward compatibility)
	ctx.setDefaultThresholdForFirstVariable()

	return nil
}

// extractJoinFieldsForBackwardCompatibility extracts join fields for legacy support
func (ctx *aggregationExtractionContext) extractJoinFieldsForBackwardCompatibility(joinConditionsMap map[string]interface{}) {
	if joinConditionsMap == nil || joinConditionsMap["type"] != "comparison" {
		return
	}

	// Extract left side
	if leftData, ok := joinConditionsMap["left"].(map[string]interface{}); ok {
		if leftData["type"] == "fieldAccess" {
			if leftObj, ok := leftData["object"].(string); ok {
				if field, ok := leftData["field"].(string); ok {
					// Determine if this is the main or agg side
					if leftObj == ctx.aggInfo.MainVariable {
						ctx.aggInfo.MainField = field
					} else {
						ctx.aggInfo.JoinField = field
					}
				}
			}
		}
	}

	// Extract right side
	if rightData, ok := joinConditionsMap["right"].(map[string]interface{}); ok {
		if rightData["type"] == "fieldAccess" {
			if rightObj, ok := rightData["object"].(string); ok {
				if field, ok := rightData["field"].(string); ok {
					if rightObj == ctx.aggInfo.MainVariable {
						ctx.aggInfo.MainField = field
					} else {
						ctx.aggInfo.JoinField = field
					}
				}
			}
		}
	}
}

// extractAndApplyThresholds extracts thresholds and applies them to aggregation variables
func (ctx *aggregationExtractionContext) extractAndApplyThresholds(thresholdConditions []map[string]interface{}) {
	for _, threshold := range thresholdConditions {
		if leftData, ok := threshold["left"].(map[string]interface{}); ok {
			if leftData["type"] == "variable" {
				if aggVarName, ok := leftData["name"].(string); ok {
					// Find the matching aggregation variable
					for i := range ctx.aggInfo.AggregationVars {
						if ctx.aggInfo.AggregationVars[i].Name == aggVarName {
							if operator, ok := threshold["operator"].(string); ok {
								ctx.aggInfo.AggregationVars[i].Operator = operator
							}
							if rightData, ok := threshold["right"].(map[string]interface{}); ok {
								if value, ok := rightData["value"].(float64); ok {
									ctx.aggInfo.AggregationVars[i].Threshold = value
								}
							}
							break
						}
					}
				}
			}
		}
	}
}

// setDefaultThresholdForFirstVariable sets default threshold for first aggregation variable
func (ctx *aggregationExtractionContext) setDefaultThresholdForFirstVariable() {
	if len(ctx.aggInfo.AggregationVars) == 0 {
		return
	}

	if ctx.aggInfo.AggregationVars[0].Operator != "" {
		ctx.aggInfo.Operator = ctx.aggInfo.AggregationVars[0].Operator
		ctx.aggInfo.Threshold = ctx.aggInfo.AggregationVars[0].Threshold
	} else {
		ctx.aggInfo.Operator = ">="
		ctx.aggInfo.Threshold = 0
		ctx.aggInfo.AggregationVars[0].Operator = ">="
		ctx.aggInfo.AggregationVars[0].Threshold = 0
	}
}

// setDefaultThresholds sets default thresholds when no constraints are present
func (ctx *aggregationExtractionContext) setDefaultThresholds() {
	ctx.aggInfo.Operator = ">="
	ctx.aggInfo.Threshold = 0
	for i := range ctx.aggInfo.AggregationVars {
		ctx.aggInfo.AggregationVars[i].Operator = ">="
		ctx.aggInfo.AggregationVars[i].Threshold = 0
	}
}

// extractMultiSourceAggregationInfoOrchestrated orchestrates the extraction of multi-source aggregation info
// using the extract method pattern to separate concerns
func (cp *ConstraintPipeline) extractMultiSourceAggregationInfoOrchestrated(
	exprMap map[string]interface{},
) (*AggregationInfo, error) {
	ctx := newAggregationExtractionContext(cp, exprMap)

	// Step 1: Validate and extract patterns list
	if err := ctx.validateAndExtractPatterns(); err != nil {
		return nil, err
	}

	// Step 2: Extract first pattern (main + aggregation variables)
	if err := ctx.extractFirstPattern(); err != nil {
		return nil, err
	}

	// Step 3: Extract source patterns from remaining blocks
	if err := ctx.extractSourcePatterns(); err != nil {
		return nil, err
	}

	// Step 4: Extract constraints, join conditions, and thresholds
	if err := ctx.extractConstraintsAndConditions(); err != nil {
		return nil, err
	}

	return ctx.aggInfo, nil
}
