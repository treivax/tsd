package rete

import (
	"fmt"
	"testing"
	"time"
)

func TestNewIndexedFactStorage(t *testing.T) {
	config := IndexConfig{
		IndexedFields:        []string{"name", "age"},
		MaxCacheSize:         1000,
		CacheTTL:             time.Minute,
		EnableCompositeIndex: true,
		AutoIndexThreshold:   100,
	}

	storage := NewIndexedFactStorage(config)

	if storage == nil {
		t.Fatal("NewIndexedFactStorage returned nil")
	}

	if storage.factsByID == nil {
		t.Error("factsByID not initialized")
	}

	if storage.factsByType == nil {
		t.Error("factsByType not initialized")
	}

	if storage.factsByField == nil {
		t.Error("factsByField not initialized")
	}

	if storage.compositeIndex == nil {
		t.Error("compositeIndex not initialized")
	}

	if storage.accessStats == nil {
		t.Error("accessStats not initialized")
	}

	if len(storage.config.IndexedFields) != 2 {
		t.Errorf("Expected 2 indexed fields, got %d", len(storage.config.IndexedFields))
	}
}

func TestStoreFact(t *testing.T) {
	config := IndexConfig{
		IndexedFields:        []string{"name"},
		EnableCompositeIndex: true,
	}
	storage := NewIndexedFactStorage(config)

	fact := &Fact{
		ID:   "fact1",
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Alice",
			"age":  30,
		},
	}

	err := storage.StoreFact(fact)
	if err != nil {
		t.Fatalf("StoreFact failed: %v", err)
	}

	// Verify fact is stored in main index
	storedFact, exists := storage.factsByID[fact.ID]
	if !exists {
		t.Error("Fact not found in factsByID")
	}
	if storedFact.ID != fact.ID {
		t.Errorf("Expected fact ID %s, got %s", fact.ID, storedFact.ID)
	}

	// Verify fact is indexed by type
	typeMap, exists := storage.factsByType[fact.Type]
	if !exists {
		t.Error("Fact type not indexed")
	}
	if _, exists := typeMap[fact.ID]; !exists {
		t.Error("Fact not found in type index")
	}

	// Verify fact is indexed by fields
	if fieldIndex, exists := storage.factsByField["name"]; exists {
		if valueMap, exists := fieldIndex["Alice"]; exists {
			if _, exists := valueMap[fact.ID]; !exists {
				t.Error("Fact not found in name field index")
			}
		} else {
			t.Error("Name value not indexed")
		}
	} else {
		t.Error("Name field not indexed")
	}
}

func TestGetFactByID(t *testing.T) {
	storage := NewIndexedFactStorage(IndexConfig{})

	fact := &Fact{
		ID:     "fact1",
		Type:   "Person",
		Fields: map[string]interface{}{"name": "Bob"},
	}

	storage.StoreFact(fact)

	// Test existing fact
	retrieved, exists := storage.GetFactByID("fact1")
	if !exists {
		t.Error("Fact should exist")
	}
	if retrieved.ID != "fact1" {
		t.Errorf("Expected fact1, got %s", retrieved.ID)
	}

	// Test non-existing fact
	_, exists = storage.GetFactByID("nonexistent")
	if exists {
		t.Error("Non-existent fact should not be found")
	}

	// Verify access is recorded
	stats := storage.GetAccessStats()
	if stats["id:fact1"] == 0 {
		t.Error("Access to fact1 was not recorded")
	}
}

func TestGetFactsByType(t *testing.T) {
	storage := NewIndexedFactStorage(IndexConfig{})

	fact1 := &Fact{ID: "f1", Type: "Person", Fields: map[string]interface{}{"name": "Alice"}}
	fact2 := &Fact{ID: "f2", Type: "Person", Fields: map[string]interface{}{"name": "Bob"}}
	fact3 := &Fact{ID: "f3", Type: "Car", Fields: map[string]interface{}{"model": "Tesla"}}

	storage.StoreFact(fact1)
	storage.StoreFact(fact2)
	storage.StoreFact(fact3)

	// Test getting Person facts
	persons := storage.GetFactsByType("Person")
	if len(persons) != 2 {
		t.Errorf("Expected 2 Person facts, got %d", len(persons))
	}

	// Test getting Car facts
	cars := storage.GetFactsByType("Car")
	if len(cars) != 1 {
		t.Errorf("Expected 1 Car fact, got %d", len(cars))
	}

	// Test non-existing type
	robots := storage.GetFactsByType("Robot")
	if len(robots) != 0 {
		t.Errorf("Expected 0 Robot facts, got %d", len(robots))
	}

	// Verify access is recorded
	stats := storage.GetAccessStats()
	if stats["type:Person"] == 0 {
		t.Error("Access to Person type was not recorded")
	}
}

func TestGetFactsByField(t *testing.T) {
	storage := NewIndexedFactStorage(IndexConfig{})

	fact1 := &Fact{ID: "f1", Type: "Person", Fields: map[string]interface{}{"age": 30}}
	fact2 := &Fact{ID: "f2", Type: "Person", Fields: map[string]interface{}{"age": 30}}
	fact3 := &Fact{ID: "f3", Type: "Person", Fields: map[string]interface{}{"age": 25}}

	storage.StoreFact(fact1)
	storage.StoreFact(fact2)
	storage.StoreFact(fact3)

	// Test getting facts with age 30
	age30Facts := storage.GetFactsByField("age", 30)
	if len(age30Facts) != 2 {
		t.Errorf("Expected 2 facts with age 30, got %d", len(age30Facts))
	}

	// Test getting facts with age 25
	age25Facts := storage.GetFactsByField("age", 25)
	if len(age25Facts) != 1 {
		t.Errorf("Expected 1 fact with age 25, got %d", len(age25Facts))
	}

	// Test non-existing field
	nameField := storage.GetFactsByField("name", "Alice")
	if len(nameField) != 0 {
		t.Errorf("Expected 0 facts with name Alice, got %d", len(nameField))
	}

	// Test non-existing value
	age99Facts := storage.GetFactsByField("age", 99)
	if len(age99Facts) != 0 {
		t.Errorf("Expected 0 facts with age 99, got %d", len(age99Facts))
	}
}

func TestCompositeIndex(t *testing.T) {
	config := IndexConfig{
		EnableCompositeIndex: true,
	}
	storage := NewIndexedFactStorage(config)

	fact := &Fact{
		ID:   "f1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   123,
			"name": "Alice",
		},
	}

	storage.StoreFact(fact)

	// Verify composite index was created
	compositeKey := "id_name:123_Alice"
	facts := storage.GetFactsByCompositeKey(compositeKey)
	if len(facts) != 1 {
		t.Errorf("Expected 1 fact in composite index, got %d", len(facts))
	}
	if facts[0].ID != "f1" {
		t.Errorf("Expected fact f1, got %s", facts[0].ID)
	}

	// Test non-existing composite key
	facts = storage.GetFactsByCompositeKey("id_name:456_Bob")
	if len(facts) != 0 {
		t.Errorf("Expected 0 facts for non-existing key, got %d", len(facts))
	}
}

func TestRemoveFact(t *testing.T) {
	config := IndexConfig{
		EnableCompositeIndex: true,
	}
	storage := NewIndexedFactStorage(config)

	fact := &Fact{
		ID:   "f1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   123,
			"name": "Alice",
			"age":  30,
		},
	}

	storage.StoreFact(fact)

	// Verify fact exists
	if storage.Size() != 1 {
		t.Errorf("Expected size 1, got %d", storage.Size())
	}

	// Remove fact
	removed := storage.RemoveFact("f1")
	if !removed {
		t.Error("RemoveFact should return true")
	}

	// Verify fact is removed from main index
	_, exists := storage.factsByID["f1"]
	if exists {
		t.Error("Fact should be removed from factsByID")
	}

	// Verify fact is removed from type index
	if typeMap, exists := storage.factsByType["Person"]; exists {
		if _, exists := typeMap["f1"]; exists {
			t.Error("Fact should be removed from type index")
		}
	}

	// Verify fact is removed from field indexes
	if fieldIndex, exists := storage.factsByField["name"]; exists {
		if valueMap, exists := fieldIndex["Alice"]; exists {
			if _, exists := valueMap["f1"]; exists {
				t.Error("Fact should be removed from field index")
			}
		}
	}

	// Verify size is updated
	if storage.Size() != 0 {
		t.Errorf("Expected size 0 after removal, got %d", storage.Size())
	}

	// Try to remove non-existing fact
	removed = storage.RemoveFact("nonexistent")
	if removed {
		t.Error("RemoveFact should return false for non-existing fact")
	}
}

func TestRemoveFactWithCompositeIndex(t *testing.T) {
	config := IndexConfig{
		EnableCompositeIndex: true,
	}
	storage := NewIndexedFactStorage(config)

	fact := &Fact{
		ID:   "f1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   123,
			"name": "Alice",
		},
	}

	storage.StoreFact(fact)

	// Verify composite index exists
	compositeKey := "id_name:123_Alice"
	facts := storage.GetFactsByCompositeKey(compositeKey)
	if len(facts) != 1 {
		t.Fatalf("Expected 1 fact in composite index before removal, got %d", len(facts))
	}

	// Remove fact
	storage.RemoveFact("f1")

	// Verify composite index is cleaned up
	facts = storage.GetFactsByCompositeKey(compositeKey)
	if len(facts) != 0 {
		t.Errorf("Expected 0 facts in composite index after removal, got %d", len(facts))
	}
}

func TestGetAccessStats(t *testing.T) {
	storage := NewIndexedFactStorage(IndexConfig{})

	fact := &Fact{ID: "f1", Type: "Person", Fields: map[string]interface{}{"name": "Alice"}}
	storage.StoreFact(fact)

	// Generate some accesses
	storage.GetFactByID("f1")
	storage.GetFactByID("f1")
	storage.GetFactsByType("Person")
	storage.GetFactsByField("name", "Alice")

	stats := storage.GetAccessStats()

	if stats["id:f1"] != 2 {
		t.Errorf("Expected 2 accesses to id:f1, got %d", stats["id:f1"])
	}

	if stats["type:Person"] != 1 {
		t.Errorf("Expected 1 access to type:Person, got %d", stats["type:Person"])
	}

	if stats["field:name:Alice"] != 1 {
		t.Errorf("Expected 1 access to field:name:Alice, got %d", stats["field:name:Alice"])
	}
}

func TestOptimizeIndexes(t *testing.T) {
	config := IndexConfig{
		AutoIndexThreshold: 5,
	}
	storage := NewIndexedFactStorage(config)

	fact := &Fact{ID: "f1", Type: "Person", Fields: map[string]interface{}{"name": "Alice"}}
	storage.StoreFact(fact)

	// Generate accesses beyond threshold
	for i := 0; i < 10; i++ {
		storage.GetFactByID("f1")
	}

	// Call optimize (currently doesn't do much, but tests the method)
	storage.OptimizeIndexes()

	// Verify stats still exist
	stats := storage.GetAccessStats()
	if stats["id:f1"] != 10 {
		t.Errorf("Expected 10 accesses after optimization, got %d", stats["id:f1"])
	}
}

func TestClear(t *testing.T) {
	storage := NewIndexedFactStorage(IndexConfig{})

	// Add multiple facts
	for i := 0; i < 5; i++ {
		fact := &Fact{
			ID:     fmt.Sprintf("fact_%d", i),
			Type:   "Person",
			Fields: map[string]interface{}{"index": i},
		}
		storage.StoreFact(fact)
	}

	if storage.Size() != 5 {
		t.Errorf("Expected size 5 before clear, got %d", storage.Size())
	}

	// Clear storage
	storage.Clear()

	// Verify everything is cleared
	if storage.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", storage.Size())
	}

	if len(storage.factsByType) != 0 {
		t.Errorf("Expected empty factsByType, got %d entries", len(storage.factsByType))
	}

	if len(storage.factsByField) != 0 {
		t.Errorf("Expected empty factsByField, got %d entries", len(storage.factsByField))
	}

	if len(storage.compositeIndex) != 0 {
		t.Errorf("Expected empty compositeIndex, got %d entries", len(storage.compositeIndex))
	}

	if len(storage.accessStats) != 0 {
		t.Errorf("Expected empty accessStats, got %d entries", len(storage.accessStats))
	}
}

func TestSize(t *testing.T) {
	storage := NewIndexedFactStorage(IndexConfig{})

	if storage.Size() != 0 {
		t.Errorf("Expected size 0 for empty storage, got %d", storage.Size())
	}

	// Add facts
	for i := 0; i < 10; i++ {
		fact := &Fact{
			ID:     fmt.Sprintf("fact_%d", i),
			Type:   "Test",
			Fields: map[string]interface{}{},
		}
		storage.StoreFact(fact)
	}

	if storage.Size() != 10 {
		t.Errorf("Expected size 10, got %d", storage.Size())
	}

	// Remove some facts
	storage.RemoveFact("fact_0")
	storage.RemoveFact("fact_1")

	if storage.Size() != 8 {
		t.Errorf("Expected size 8 after removals, got %d", storage.Size())
	}
}

func TestConcurrentAccess(t *testing.T) {
	storage := NewIndexedFactStorage(IndexConfig{})

	// Test concurrent writes (StoreFact has proper locking)
	numFacts := 50
	numGoroutines := 5
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer func() { done <- true }()
			for j := 0; j < numFacts/numGoroutines; j++ {
				fact := &Fact{
					ID:     fmt.Sprintf("fact_%d_%d", goroutineID, j),
					Type:   "Test",
					Fields: map[string]interface{}{"index": j, "goroutine": goroutineID},
				}
				storage.StoreFact(fact)
			}
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify storage integrity
	expectedSize := numFacts
	if storage.Size() != expectedSize {
		t.Errorf("Expected size %d after concurrent writes, got %d", expectedSize, storage.Size())
	}

	// Verify we can retrieve facts by type
	facts := storage.GetFactsByType("Test")
	if len(facts) != expectedSize {
		t.Errorf("Expected %d facts by type, got %d", expectedSize, len(facts))
	}
}

func TestMultipleFieldIndexing(t *testing.T) {
	config := IndexConfig{
		IndexedFields: []string{"name", "age", "city"},
	}
	storage := NewIndexedFactStorage(config)

	fact := &Fact{
		ID:   "f1",
		Type: "Person",
		Fields: map[string]interface{}{
			"name": "Alice",
			"age":  30,
			"city": "Paris",
		},
	}

	storage.StoreFact(fact)

	// Verify all fields are indexed
	nameFacts := storage.GetFactsByField("name", "Alice")
	if len(nameFacts) != 1 {
		t.Errorf("Expected 1 fact indexed by name, got %d", len(nameFacts))
	}

	ageFacts := storage.GetFactsByField("age", 30)
	if len(ageFacts) != 1 {
		t.Errorf("Expected 1 fact indexed by age, got %d", len(ageFacts))
	}

	cityFacts := storage.GetFactsByField("city", "Paris")
	if len(cityFacts) != 1 {
		t.Errorf("Expected 1 fact indexed by city, got %d", len(cityFacts))
	}
}
