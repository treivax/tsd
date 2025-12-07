// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package nodes

import (
	"fmt"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// Aggregate function type constants
const (
	AggregateFunctionSum   = "SUM"
	AggregateFunctionCount = "COUNT"
	AggregateFunctionAvg   = "AVG"
	AggregateFunctionMin   = "MIN"
	AggregateFunctionMax   = "MAX"
)

// AggregationCalculator provides functions for computing aggregations over facts
type AggregationCalculator struct{}

// NewAggregationCalculator creates a new aggregation calculator
func NewAggregationCalculator() *AggregationCalculator {
	return &AggregationCalculator{}
}

// ComputeSum calculates the sum of a numeric field across facts
func (ac *AggregationCalculator) ComputeSum(facts []*domain.Fact, field string) (float64, error) {
	var sum float64
	count := 0

	for _, fact := range facts {
		if value, exists := fact.Fields[field]; exists {
			switch v := value.(type) {
			case int:
				sum += float64(v)
				count++
			case int64:
				sum += float64(v)
				count++
			case float64:
				sum += v
				count++
			case float32:
				sum += float64(v)
				count++
			default:
				// Ignore non-numeric values
				continue
			}
		}
	}

	return sum, nil
}

// ComputeAverage calculates the average of a numeric field across facts
func (ac *AggregationCalculator) ComputeAverage(facts []*domain.Fact, field string) (float64, error) {
	sum, err := ac.ComputeSum(facts, field)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, fact := range facts {
		if _, exists := fact.Fields[field]; exists {
			count++
		}
	}

	if count == 0 {
		return 0, nil
	}

	return sum / float64(count), nil
}

// ComputeCount counts the number of facts
func (ac *AggregationCalculator) ComputeCount(facts []*domain.Fact) int {
	return len(facts)
}

// ComputeMin calculates the minimum value for a field across facts
func (ac *AggregationCalculator) ComputeMin(facts []*domain.Fact, field string) (interface{}, error) {
	return ac.computeMinMax(facts, field, true)
}

// ComputeMax calculates the maximum value for a field across facts
func (ac *AggregationCalculator) ComputeMax(facts []*domain.Fact, field string) (interface{}, error) {
	return ac.computeMinMax(facts, field, false)
}

// computeMinMax calculates either minimum or maximum value for a field across facts
// isMin parameter: true for minimum, false for maximum
func (ac *AggregationCalculator) computeMinMax(facts []*domain.Fact, field string, isMin bool) (interface{}, error) {
	var resultFloat float64
	var resultString string
	var resultOther interface{}
	foundNumeric := false
	foundString := false
	foundOther := false

	for _, fact := range facts {
		if value, exists := fact.Fields[field]; exists {
			switch v := value.(type) {
			case int:
				floatVal := float64(v)
				if !foundNumeric || ac.shouldUpdateNumeric(floatVal, resultFloat, isMin) {
					resultFloat = floatVal
					foundNumeric = true
				}
			case int64:
				floatVal := float64(v)
				if !foundNumeric || ac.shouldUpdateNumeric(floatVal, resultFloat, isMin) {
					resultFloat = floatVal
					foundNumeric = true
				}
			case float32:
				floatVal := float64(v)
				if !foundNumeric || ac.shouldUpdateNumeric(floatVal, resultFloat, isMin) {
					resultFloat = floatVal
					foundNumeric = true
				}
			case float64:
				if !foundNumeric || ac.shouldUpdateNumeric(v, resultFloat, isMin) {
					resultFloat = v
					foundNumeric = true
				}
			case string:
				if !foundString || ac.shouldUpdateString(v, resultString, isMin) {
					resultString = v
					foundString = true
				}
			default:
				if !foundOther {
					resultOther = v
					foundOther = true
				}
			}
		}
	}

	// Return the most appropriate type
	if foundNumeric {
		return resultFloat, nil
	}
	if foundString {
		return resultString, nil
	}
	if foundOther {
		return resultOther, nil
	}

	return nil, fmt.Errorf("no values found for field %s", field)
}

// shouldUpdateNumeric determines if a new numeric value should replace the current one
func (ac *AggregationCalculator) shouldUpdateNumeric(newVal, currentVal float64, isMin bool) bool {
	if isMin {
		return newVal < currentVal
	}
	return newVal > currentVal
}

// shouldUpdateString determines if a new string value should replace the current one
func (ac *AggregationCalculator) shouldUpdateString(newVal, currentVal string, isMin bool) bool {
	if isMin {
		return newVal < currentVal
	}
	return newVal > currentVal
}
