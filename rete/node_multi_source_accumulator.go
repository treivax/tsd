// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License

package rete

import (
	"fmt"
	"sync"
)

// MultiSourceAccumulatorNode computes multiple aggregations over multiple joined sources
// This node sits at the end of a join chain and performs aggregation calculations
// on the combined facts, evaluating thresholds before firing actions.
type MultiSourceAccumulatorNode struct {
	BaseNode
	MainVariable    string                `json:"main_variable"`    // Main variable (e.g., "d")
	MainType        string                `json:"main_type"`        // Main type (e.g., "Department")
	AggregationVars []AggregationVariable `json:"aggregation_vars"` // Multiple aggregations to compute
	SourcePatterns  []SourcePattern       `json:"source_patterns"`  // Source patterns being aggregated
	Storage         Storage               `json:"-"`
	mutex           sync.RWMutex

	// Memory structures
	MainFacts      map[string]*Fact              `json:"-"` // Main facts indexed by ID
	CombinedTokens map[string]map[string]*Token  `json:"-"` // Tokens grouped by main fact ID
	AggregateCache map[string]map[string]float64 `json:"-"` // Cached aggregation results: mainFactID -> aggVarName -> value
}

// NewMultiSourceAccumulatorNode creates a new multi-source accumulator node
func NewMultiSourceAccumulatorNode(
	id string,
	mainVar string,
	mainType string,
	aggVars []AggregationVariable,
	sources []SourcePattern,
	storage Storage,
) *MultiSourceAccumulatorNode {
	return &MultiSourceAccumulatorNode{
		BaseNode: BaseNode{
			ID:       id,
			Type:     "multi_source_accumulator",
			Children: make([]Node, 0),
			Memory:   &WorkingMemory{Tokens: make(map[string]*Token), Facts: make(map[string]*Fact)},
		},
		MainVariable:    mainVar,
		MainType:        mainType,
		AggregationVars: aggVars,
		SourcePatterns:  sources,
		Storage:         storage,
		MainFacts:       make(map[string]*Fact),
		CombinedTokens:  make(map[string]map[string]*Token),
		AggregateCache:  make(map[string]map[string]float64),
	}
}

// Activate processes incoming tokens from the join chain
func (msn *MultiSourceAccumulatorNode) Activate(fact *Fact, token *Token) error {
	msn.mutex.Lock()
	defer msn.mutex.Unlock()

	if token == nil {
		return fmt.Errorf("multi-source accumulator requires token")
	}

	// Extract main fact from token bindings
	mainFact := token.GetBinding(msn.MainVariable)
	if mainFact == nil {
		return fmt.Errorf("main variable %s not found in token bindings", msn.MainVariable)
	}

	mainFactID := mainFact.ID

	// Store main fact
	msn.MainFacts[mainFactID] = mainFact

	// Store token indexed by main fact
	if msn.CombinedTokens[mainFactID] == nil {
		msn.CombinedTokens[mainFactID] = make(map[string]*Token)
	}
	msn.CombinedTokens[mainFactID][token.ID] = token

	fmt.Printf("📊 MULTI_ACCUMULATOR[%s]: Received token %s for main fact %s\n",
		msn.ID, token.ID, mainFactID)

	// Recompute aggregations for this main fact
	return msn.processMainFact(mainFact)
}

// processMainFact computes all aggregations for a main fact and evaluates thresholds
func (msn *MultiSourceAccumulatorNode) processMainFact(mainFact *Fact) error {
	mainFactID := mainFact.ID

	// Get all tokens for this main fact
	tokens := msn.CombinedTokens[mainFactID]
	if len(tokens) == 0 {
		fmt.Printf("📊 MULTI_ACCUMULATOR[%s]: No tokens for main fact %s\n", msn.ID, mainFactID)
		return nil
	}

	// Initialize aggregate cache for this main fact
	if msn.AggregateCache[mainFactID] == nil {
		msn.AggregateCache[mainFactID] = make(map[string]float64)
	}

	// Compute each aggregation variable
	allThresholdsSatisfied := true
	for _, aggVar := range msn.AggregationVars {
		value, err := msn.computeAggregation(mainFactID, aggVar, tokens)
		if err != nil {
			return fmt.Errorf("error computing aggregation %s: %w", aggVar.Name, err)
		}

		// Cache the result
		msn.AggregateCache[mainFactID][aggVar.Name] = value

		fmt.Printf("📊 MULTI_ACCUMULATOR[%s]: %s = %.2f for main fact %s\n",
			msn.ID, aggVar.Name, value, mainFactID)

		// Check threshold if specified
		if aggVar.Operator != "" && aggVar.Operator != ">=" || aggVar.Threshold != 0 {
			if !msn.evaluateThreshold(value, aggVar.Operator, aggVar.Threshold) {
				fmt.Printf("❌ MULTI_ACCUMULATOR[%s]: Threshold not satisfied: %s (%.2f %s %.2f)\n",
					msn.ID, aggVar.Name, value, aggVar.Operator, aggVar.Threshold)
				allThresholdsSatisfied = false
			} else {
				fmt.Printf("✅ MULTI_ACCUMULATOR[%s]: Threshold satisfied: %s (%.2f %s %.2f)\n",
					msn.ID, aggVar.Name, value, aggVar.Operator, aggVar.Threshold)
			}
		}
	}

	// Only fire if all thresholds are satisfied
	if allThresholdsSatisfied {
		fmt.Printf("✅ MULTI_ACCUMULATOR[%s]: All thresholds satisfied for main fact %s\n",
			msn.ID, mainFactID)

		// Create a token with the main fact and aggregation results
		// We'll use the first token as a template and add aggregate values
		var firstToken *Token
		for _, t := range tokens {
			firstToken = t
			break
		}

		if firstToken != nil {
			newToken := &Token{
				ID:       fmt.Sprintf("msaccum_%s", mainFactID),
				Facts:    firstToken.Facts,
				Bindings: firstToken.Bindings,
			}

			// Store token in memory
			msn.Memory.AddToken(newToken)

			// Propagate to children
			for _, child := range msn.Children {
				if err := child.ActivateLeft(newToken); err != nil {
					fmt.Printf("⚠️  Error activating child: %v\n", err)
				}
			}
		}
	} else {
		fmt.Printf("❌ MULTI_ACCUMULATOR[%s]: Not all thresholds satisfied for main fact %s\n",
			msn.ID, mainFactID)
	}

	return nil
}

// computeAggregation computes a single aggregation value
func (msn *MultiSourceAccumulatorNode) computeAggregation(
	mainFactID string,
	aggVar AggregationVariable,
	tokens map[string]*Token,
) (float64, error) {
	// Collect all relevant facts for this aggregation
	var values []float64
	seenFacts := make(map[string]bool) // Deduplicate facts

	for _, token := range tokens {
		// Get the source fact for this aggregation
		sourceFact := token.GetBinding(aggVar.SourceVar)
		if sourceFact == nil {
			continue
		}

		// Skip if we've already seen this fact
		if seenFacts[sourceFact.ID] {
			continue
		}
		seenFacts[sourceFact.ID] = true

		// For COUNT, just count the facts
		if aggVar.Function == "COUNT" {
			values = append(values, 1.0)
			continue
		}

		// Get the field value
		fieldValue, exists := sourceFact.Fields[aggVar.Field]
		if !exists {
			fmt.Printf("⚠️  MULTI_ACCUMULATOR[%s]: Field %s not found in fact %s\n",
				msn.ID, aggVar.Field, sourceFact.ID)
			continue
		}

		// Convert to float64
		numValue := msn.toFloat64(fieldValue)
		values = append(values, numValue)
	}

	// Compute the aggregation based on function
	return msn.calculateAggregate(aggVar.Function, values)
}

// calculateAggregate performs the actual aggregation calculation
func (msn *MultiSourceAccumulatorNode) calculateAggregate(function string, values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, nil
	}

	switch function {
	case "COUNT":
		return float64(len(values)), nil

	case "SUM":
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		return sum, nil

	case "AVG":
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		return sum / float64(len(values)), nil

	case "MIN":
		min := values[0]
		for _, v := range values {
			if v < min {
				min = v
			}
		}
		return min, nil

	case "MAX":
		max := values[0]
		for _, v := range values {
			if v > max {
				max = v
			}
		}
		return max, nil

	default:
		return 0, fmt.Errorf("unsupported aggregation function: %s", function)
	}
}

// evaluateThreshold checks if a value satisfies a threshold condition
func (msn *MultiSourceAccumulatorNode) evaluateThreshold(value float64, operator string, threshold float64) bool {
	switch operator {
	case ">":
		return value > threshold
	case ">=":
		return value >= threshold
	case "<":
		return value < threshold
	case "<=":
		return value <= threshold
	case "==":
		return value == threshold
	case "!=":
		return value != threshold
	default:
		// No operator or ">=" with threshold 0 (default) - always satisfied
		if operator == "" || (operator == ">=" && threshold == 0) {
			return true
		}
		return false
	}
}

// toFloat64 converts various numeric types to float64
func (msn *MultiSourceAccumulatorNode) toFloat64(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	default:
		fmt.Printf("⚠️  Cannot convert %v (type %T) to float64, returning 0\n", val, val)
		return 0
	}
}

// ActivateLeft handles left-side activation (from join chain)
func (msn *MultiSourceAccumulatorNode) ActivateLeft(token *Token) error {
	if len(token.Facts) > 0 {
		return msn.Activate(token.Facts[0], token)
	}
	return msn.Activate(nil, token)
}

// ActivateRight is not used for multi-source accumulator (only receives from left)
func (msn *MultiSourceAccumulatorNode) ActivateRight(fact *Fact) error {
	return fmt.Errorf("multi-source accumulator does not support right activation")
}

// ActivateRetract handles fact retraction
func (msn *MultiSourceAccumulatorNode) ActivateRetract(factID string) error {
	msn.mutex.Lock()
	defer msn.mutex.Unlock()

	// Check if this is a main fact being retracted
	if mainFact, exists := msn.MainFacts[factID]; exists {
		fmt.Printf("🔄 MULTI_ACCUMULATOR[%s]: Retracting main fact %s\n", msn.ID, factID)
		msn.ClearMainFact(mainFact.ID)

		// Propagate retraction to children
		for _, child := range msn.Children {
			if err := child.ActivateRetract(factID); err != nil {
				fmt.Printf("⚠️  Error propagating retraction to child: %v\n", err)
			}
		}
		return nil
	}

	// Otherwise, scan tokens to find if this fact is part of any combined tokens
	for mainFactID, tokens := range msn.CombinedTokens {
		for tokenID, token := range tokens {
			// Check if this fact is in the token
			factInToken := false
			for _, f := range token.Facts {
				if f.ID == factID {
					factInToken = true
					break
				}
			}

			if factInToken {
				// Remove this token
				delete(msn.CombinedTokens[mainFactID], tokenID)
				fmt.Printf("🔄 MULTI_ACCUMULATOR[%s]: Removed token %s due to fact %s retraction\n",
					msn.ID, tokenID, factID)
			}
		}

		// If no tokens left for this main fact, clear its cache
		if len(msn.CombinedTokens[mainFactID]) == 0 {
			delete(msn.CombinedTokens, mainFactID)
			delete(msn.AggregateCache, mainFactID)
		} else {
			// Recompute aggregations for this main fact
			if mainFact, exists := msn.MainFacts[mainFactID]; exists {
				if err := msn.processMainFact(mainFact); err != nil {
					fmt.Printf("⚠️  Error recomputing aggregations after retraction: %v\n", err)
				}
			}
		}
	}

	return nil
}

// GetAggregateValue retrieves a cached aggregation value
func (msn *MultiSourceAccumulatorNode) GetAggregateValue(mainFactID string, aggVarName string) (float64, bool) {
	msn.mutex.RLock()
	defer msn.mutex.RUnlock()

	if msn.AggregateCache[mainFactID] == nil {
		return 0, false
	}

	value, exists := msn.AggregateCache[mainFactID][aggVarName]
	return value, exists
}

// ClearMainFact removes all data associated with a main fact (internal - no locking)
func (msn *MultiSourceAccumulatorNode) ClearMainFact(mainFactID string) {
	delete(msn.MainFacts, mainFactID)
	delete(msn.CombinedTokens, mainFactID)
	delete(msn.AggregateCache, mainFactID)

	fmt.Printf("🧹 MULTI_ACCUMULATOR[%s]: Cleared data for main fact %s\n", msn.ID, mainFactID)
}

// GetStats returns statistics about the accumulator state
func (msn *MultiSourceAccumulatorNode) GetStats() map[string]interface{} {
	msn.mutex.RLock()
	defer msn.mutex.RUnlock()

	totalTokens := 0
	for _, tokens := range msn.CombinedTokens {
		totalTokens += len(tokens)
	}

	return map[string]interface{}{
		"node_id":           msn.ID,
		"node_type":         msn.Type,
		"main_variable":     msn.MainVariable,
		"main_facts_count":  len(msn.MainFacts),
		"total_tokens":      totalTokens,
		"aggregation_count": len(msn.AggregationVars),
		"source_count":      len(msn.SourcePatterns),
	}
}
