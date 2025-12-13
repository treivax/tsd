// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"
)

func TestNewRuleBuilder(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	// Create a mock pipeline (nil for basic test)
	rb := NewRuleBuilder(utils, nil)
	if rb == nil {
		t.Fatal("NewRuleBuilder returned nil")
	}
	if rb.utils != utils {
		t.Error("RuleBuilder.utils not set correctly")
	}
	if rb.alphaBuilder == nil {
		t.Error("RuleBuilder.alphaBuilder not initialized")
	}
	if rb.joinBuilder == nil {
		t.Error("RuleBuilder.joinBuilder not initialized")
	}
	if rb.existsBuilder == nil {
		t.Error("RuleBuilder.existsBuilder not initialized")
	}
	if rb.accumulatorBuilder == nil {
		t.Error("RuleBuilder.accumulatorBuilder not initialized")
	}
}
func TestRuleBuilder_CreateRuleNodes_InvalidExpression(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	utils := NewBuilderUtils(storage)
	rb := NewRuleBuilder(utils, nil)
	t.Run("error on invalid expression format", func(t *testing.T) {
		expressions := []interface{}{
			"not a map", // Invalid format
		}
		err := rb.CreateRuleNodes(network, expressions)
		if err == nil {
			t.Error("Expected error for invalid expression format, got nil")
		}
	})
}
func TestRuleBuilder_CreateRuleNodes_EmptyList(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	utils := NewBuilderUtils(storage)
	rb := NewRuleBuilder(utils, nil)
	t.Run("empty expression list", func(t *testing.T) {
		expressions := []interface{}{}
		err := rb.CreateRuleNodes(network, expressions)
		if err != nil {
			t.Errorf("Empty expression list should not error, got: %v", err)
		}
	})
}
func TestRuleBuilder_CreateSingleRule_NoPipeline(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	utils := NewBuilderUtils(storage)
	rb := NewRuleBuilder(utils, nil)
	t.Run("error when pipeline is nil", func(t *testing.T) {
		exprMap := map[string]interface{}{
			"type": "rule",
		}
		err := rb.CreateSingleRule(network, "test_rule", exprMap)
		if err == nil {
			t.Error("Expected error when pipeline is nil, got nil")
		}
		if err != nil && err.Error() != "pipeline is nil - cannot create rule" {
			t.Errorf("Expected 'pipeline is nil' error message, got: %v", err)
		}
	})
}
func TestRuleBuilder_CreateRuleByType(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	utils := NewBuilderUtils(storage)
	rb := NewRuleBuilder(utils, nil)
	action := &Action{
		Type: "test",
	}
	condition := map[string]interface{}{
		"type": ConditionTypeSimple,
	}
	variables := []map[string]interface{}{}
	variableNames := []string{"x"}
	variableTypes := []string{"Person"}
	t.Run("unknown rule type returns error", func(t *testing.T) {
		err := rb.createRuleByType(
			network,
			"test_rule",
			"unknown_type",
			map[string]interface{}{},
			condition,
			action,
			variables,
			variableNames,
			variableTypes,
			nil,
			false,
		)
		if err == nil {
			t.Error("Expected error for unknown rule type, got nil")
		}
	})
}
func TestRuleBuilder_Delegation(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	rb := NewRuleBuilder(utils, nil)
	t.Run("alphaBuilder is properly initialized", func(t *testing.T) {
		if rb.alphaBuilder == nil {
			t.Error("alphaBuilder should be initialized")
		}
		if rb.alphaBuilder.utils != utils {
			t.Error("alphaBuilder should share utils")
		}
	})
	t.Run("joinBuilder is properly initialized", func(t *testing.T) {
		if rb.joinBuilder == nil {
			t.Error("joinBuilder should be initialized")
		}
		if rb.joinBuilder.utils != utils {
			t.Error("joinBuilder should share utils")
		}
	})
	t.Run("existsBuilder is properly initialized", func(t *testing.T) {
		if rb.existsBuilder == nil {
			t.Error("existsBuilder should be initialized")
		}
		if rb.existsBuilder.utils != utils {
			t.Error("existsBuilder should share utils")
		}
	})
	t.Run("accumulatorBuilder is properly initialized", func(t *testing.T) {
		if rb.accumulatorBuilder == nil {
			t.Error("accumulatorBuilder should be initialized")
		}
		if rb.accumulatorBuilder.utils != utils {
			t.Error("accumulatorBuilder should share utils")
		}
	})
}
func TestRuleBuilder_PipelineReference(t *testing.T) {
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)

	// Create a mock pipeline that implements PipelineHelper interface
	mockPipeline := &mockPipelineHelper{}

	t.Run("pipeline reference is stored", func(t *testing.T) {
		rb := NewRuleBuilder(utils, mockPipeline)
		if rb.pipeline == nil {
			t.Error("pipeline should be stored")
		}
		// The pipeline is now typed as PipelineHelper, no need for type assertion
		if rb.pipeline != mockPipeline {
			t.Error("pipeline reference not stored correctly")
		}
	})
	t.Run("nil pipeline is acceptable", func(t *testing.T) {
		rb := NewRuleBuilder(utils, nil)
		if rb == nil {
			t.Error("RuleBuilder should be created even with nil pipeline")
		}
	})
}

// mockPipelineHelper is a minimal mock implementation of PipelineHelper for testing
type mockPipelineHelper struct{}

func (m *mockPipelineHelper) extractActionFromExpression(exprMap map[string]interface{}, ruleID string) (*Action, error) {
	return nil, nil
}

func (m *mockPipelineHelper) detectAggregation(constraintsData interface{}) bool {
	return false
}

func (m *mockPipelineHelper) extractAggregationInfo(constraintsData interface{}) (*AggregationInfo, error) {
	return nil, nil
}

func (m *mockPipelineHelper) extractAggregationInfoFromVariables(exprMap map[string]interface{}) (*AggregationInfo, error) {
	return nil, nil
}

func (m *mockPipelineHelper) extractMultiSourceAggregationInfo(exprMap map[string]interface{}) (*AggregationInfo, error) {
	return nil, nil
}

func (m *mockPipelineHelper) hasAggregationVariables(exprMap map[string]interface{}) bool {
	return false
}

func (m *mockPipelineHelper) buildConditionFromConstraints(constraintsData interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (m *mockPipelineHelper) extractVariablesFromExpression(exprMap map[string]interface{}) ([]map[string]interface{}, []string, []string) {
	return nil, nil, nil
}

func (m *mockPipelineHelper) determineRuleType(exprMap map[string]interface{}, varCount int, hasAggregation bool) string {
	return "alpha"
}

func (m *mockPipelineHelper) logRuleCreation(ruleType string, ruleID string, variableNames []string) {
}
