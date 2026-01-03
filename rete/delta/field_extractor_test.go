// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"sort"
	"testing"
)

func TestAlphaConditionExtractor_SimpleFieldAccess(t *testing.T) {
	condition := map[string]interface{}{
		"type":  "fieldAccess",
		"field": "price",
	}

	extractor := &AlphaConditionExtractor{}
	fields, err := extractor.ExtractFields(condition)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}

	if fields[0] != "price" {
		t.Errorf("Expected 'price', got '%s'", fields[0])
	}
}

func TestAlphaConditionExtractor_Comparison(t *testing.T) {
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "price",
		},
		"right": map[string]interface{}{
			"type":  "number",
			"value": 100,
		},
	}

	extractor := &AlphaConditionExtractor{}
	fields, err := extractor.ExtractFields(condition)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}

	if fields[0] != "price" {
		t.Errorf("Expected 'price', got '%s'", fields[0])
	}
}

func TestAlphaConditionExtractor_BinaryOp(t *testing.T) {
	condition := map[string]interface{}{
		"type":     "binaryOp",
		"operator": "&&",
		"left": map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":  "fieldAccess",
				"field": "price",
			},
			"right": 100,
		},
		"right": map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":  "fieldAccess",
				"field": "status",
			},
			"right": "active",
		},
	}

	extractor := &AlphaConditionExtractor{}
	fields, err := extractor.ExtractFields(condition)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(fields))
	}

	sort.Strings(fields)

	if fields[0] != "price" || fields[1] != "status" {
		t.Errorf("Expected ['price', 'status'], got %v", fields)
	}
}

func TestBetaConditionExtractor_JoinCondition(t *testing.T) {
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
		"left": map[string]interface{}{
			"type":     "fieldAccess",
			"variable": "order",
			"field":    "customer_id",
		},
		"right": map[string]interface{}{
			"type":     "fieldAccess",
			"variable": "customer",
			"field":    "id",
		},
	}

	extractor := &BetaConditionExtractor{}
	fields, err := extractor.ExtractFields(condition)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(fields))
	}

	fieldMap := make(map[string]bool)
	for _, field := range fields {
		fieldMap[field] = true
	}

	if !fieldMap["customer_id"] || !fieldMap["id"] {
		t.Errorf("Expected ['customer_id', 'id'], got %v", fields)
	}
}

func TestActionFieldExtractor_UpdateAction(t *testing.T) {
	action := map[string]interface{}{
		"type":     "updateWithModifications",
		"variable": "product",
		"modifications": map[string]interface{}{
			"price":  150,
			"status": "updated",
		},
	}

	extractor := &ActionFieldExtractor{}
	fields, err := extractor.ExtractFields(action)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(fields))
	}

	fieldMap := make(map[string]bool)
	for _, field := range fields {
		fieldMap[field] = true
	}

	if !fieldMap["price"] || !fieldMap["status"] {
		t.Errorf("Expected ['price', 'status'], got %v", fields)
	}
}

func TestActionFieldExtractor_InsertAction(t *testing.T) {
	action := map[string]interface{}{
		"type":     "factCreation",
		"factType": "Alert",
		"fields": map[string]interface{}{
			"severity": "high",
			"message":  "test",
		},
	}

	extractor := &ActionFieldExtractor{}
	fields, err := extractor.ExtractFields(action)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(fields))
	}

	fieldMap := make(map[string]bool)
	for _, field := range fields {
		fieldMap[field] = true
	}

	if !fieldMap["severity"] || !fieldMap["message"] {
		t.Errorf("Expected ['severity', 'message'], got %v", fields)
	}
}

func TestExtractFieldsRecursive_EmptyNode(t *testing.T) {
	fields := make(map[string]bool)
	err := extractFieldsRecursive(nil, fields)

	if err != nil {
		t.Errorf("Expected no error for nil node, got %v", err)
	}

	if len(fields) != 0 {
		t.Error("Expected no fields for nil node")
	}
}

func TestExtractFieldsRecursive_NestedStructure(t *testing.T) {
	node := map[string]interface{}{
		"type": "binaryOp",
		"left": map[string]interface{}{
			"type": "binaryOp",
			"left": map[string]interface{}{
				"type":  "fieldAccess",
				"field": "field1",
			},
			"right": map[string]interface{}{
				"type":  "fieldAccess",
				"field": "field2",
			},
		},
		"right": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "field3",
		},
	}

	fields := make(map[string]bool)
	err := extractFieldsRecursive(node, fields)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(fields) != 3 {
		t.Fatalf("Expected 3 fields, got %d", len(fields))
	}

	for i := 1; i <= 3; i++ {
		fieldName := "field" + string(rune('0'+i))
		if !fields[fieldName] {
			t.Errorf("Expected field '%s' to be extracted", fieldName)
		}
	}
}

func TestExtractFieldsRecursive_SliceInput(t *testing.T) {
	nodes := []interface{}{
		map[string]interface{}{
			"type":  "fieldAccess",
			"field": "field1",
		},
		map[string]interface{}{
			"type":  "fieldAccess",
			"field": "field2",
		},
	}

	fields := make(map[string]bool)
	err := extractFieldsRecursive(nodes, fields)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(fields))
	}
}

func TestExtractFieldsFromAlphaCondition_HelperFunction(t *testing.T) {
	condition := map[string]interface{}{
		"type":  "fieldAccess",
		"field": "testField",
	}

	fields, err := ExtractFieldsFromAlphaCondition(condition)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(fields) != 1 || fields[0] != "testField" {
		t.Errorf("Expected ['testField'], got %v", fields)
	}
}

func TestExtractFieldsFromBetaCondition_HelperFunction(t *testing.T) {
	condition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "leftField",
		},
		"right": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "rightField",
		},
	}

	fields, err := ExtractFieldsFromBetaCondition(condition)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(fields))
	}
}

func TestExtractFieldsFromAction_HelperFunction(t *testing.T) {
	action := map[string]interface{}{
		"type": "updateWithModifications",
		"modifications": map[string]interface{}{
			"field1": "value1",
		},
	}

	fields, err := ExtractFieldsFromAction(action)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(fields) != 1 || fields[0] != "field1" {
		t.Errorf("Expected ['field1'], got %v", fields)
	}
}

func BenchmarkExtractFieldsFromAlphaCondition(b *testing.B) {
	condition := map[string]interface{}{
		"type": "binaryOp",
		"left": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "price",
		},
		"right": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "status",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ExtractFieldsFromAlphaCondition(condition)
	}
}
