// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewDeltaPropagatorBuilder(t *testing.T) {
	builder := NewDeltaPropagatorBuilder()

	if builder == nil {
		t.Fatal("Builder should not be nil")
	}

	if builder.config.DefaultMode != PropagationModeAuto {
		t.Error("Expected default config")
	}
}

func TestDeltaPropagatorBuilder_Build_RequiresIndex(t *testing.T) {
	builder := NewDeltaPropagatorBuilder()

	_, err := builder.Build()

	if err == nil {
		t.Error("Expected error when index is missing")
	}
}

func TestDeltaPropagatorBuilder_Build_ValidConfig(t *testing.T) {
	builder := NewDeltaPropagatorBuilder().
		WithIndex(NewDependencyIndex()).
		WithDetector(NewDeltaDetector()).
		WithStrategy(&SequentialStrategy{})

	propagator, err := builder.Build()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if propagator == nil {
		t.Fatal("Propagator should not be nil")
	}
}

func TestDeltaPropagatorBuilder_Build_InvalidConfig(t *testing.T) {
	invalidConfig := DefaultPropagationConfig()
	invalidConfig.DeltaThreshold = 1.5

	builder := NewDeltaPropagatorBuilder().
		WithIndex(NewDependencyIndex()).
		WithConfig(invalidConfig)

	_, err := builder.Build()

	if err == nil {
		t.Error("Expected error with invalid config")
	}
}

func TestDeltaPropagatorBuilder_Build_DefaultDetector(t *testing.T) {
	builder := NewDeltaPropagatorBuilder().
		WithIndex(NewDependencyIndex())

	propagator, err := builder.Build()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if propagator.detector == nil {
		t.Error("Expected default detector to be created")
	}
}

func TestDeltaPropagatorBuilder_Build_DefaultStrategy(t *testing.T) {
	builder := NewDeltaPropagatorBuilder().
		WithIndex(NewDependencyIndex())

	propagator, err := builder.Build()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if propagator.strategy == nil {
		t.Error("Expected default strategy to be created")
	}

	if propagator.strategy.GetName() != "Sequential" {
		t.Errorf("Expected Sequential strategy, got %s", propagator.strategy.GetName())
	}
}

func TestDeltaPropagator_PropagateUpdate_FeatureDisabled(t *testing.T) {
	config := DefaultPropagationConfig()
	config.EnableDeltaPropagation = false

	propagator, _ := NewDeltaPropagatorBuilder().
		WithIndex(NewDependencyIndex()).
		WithConfig(config).
		Build()

	oldFact := map[string]interface{}{"price": 100}
	newFact := map[string]interface{}{"price": 120}

	err := propagator.PropagateUpdate(oldFact, newFact, "Product~1", "Product")

	// Classic propagation callback not configured
	if err == nil {
		t.Error("Expected error for classic propagation callback not configured")
	}

	expectedMsg := "classic propagation callback not configured"
	if err != nil && !contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error message to contain '%s', got: %v", expectedMsg, err)
	}
}

func TestDeltaPropagator_PropagateUpdate_NoChanges(t *testing.T) {
	config := DefaultPropagationConfig()
	config.MinFieldsForDelta = 1

	propagator, _ := NewDeltaPropagatorBuilder().
		WithIndex(NewDependencyIndex()).
		WithConfig(config).
		Build()

	sameFact := map[string]interface{}{"price": 100, "name": "A", "stock": 10}

	err := propagator.PropagateUpdate(sameFact, sameFact, "Product~1", "Product")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestDeltaPropagator_PropagateUpdate_WithCallback(t *testing.T) {
	index := NewDependencyIndex()
	index.AddAlphaNode("alpha1", "Product", []string{"price"})

	callbackCalled := false
	var receivedDelta *FactDelta

	callback := func(nodeID string, delta *FactDelta) error {
		callbackCalled = true
		receivedDelta = delta
		return nil
	}

	config := DefaultPropagationConfig()
	config.MinFieldsForDelta = 2

	propagator, _ := NewDeltaPropagatorBuilder().
		WithIndex(index).
		WithPropagateCallback(callback).
		WithConfig(config).
		Build()

	oldFact := map[string]interface{}{"price": 100, "name": "Product A", "stock": 10}
	newFact := map[string]interface{}{"price": 120, "name": "Product A", "stock": 10}

	err := propagator.PropagateUpdate(oldFact, newFact, "Product~1", "Product")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !callbackCalled {
		t.Error("Callback should have been called")
	}

	if receivedDelta == nil {
		t.Fatal("Delta should not be nil")
	}

	if receivedDelta.FactID != "Product~1" {
		t.Errorf("Expected FactID Product~1, got %s", receivedDelta.FactID)
	}
}

func TestDeltaPropagator_PropagateUpdate_CallbackError(t *testing.T) {
	index := NewDependencyIndex()
	index.AddAlphaNode("alpha1", "Product", []string{"price"})

	expectedError := errors.New("propagation failed")

	callback := func(nodeID string, delta *FactDelta) error {
		return expectedError
	}

	config := DefaultPropagationConfig()
	config.RetryOnError = false
	config.MinFieldsForDelta = 1

	propagator, _ := NewDeltaPropagatorBuilder().
		WithIndex(index).
		WithPropagateCallback(callback).
		WithConfig(config).
		Build()

	oldFact := map[string]interface{}{"price": 100, "name": "A", "stock": 10}
	newFact := map[string]interface{}{"price": 120, "name": "A", "stock": 10}

	err := propagator.PropagateUpdate(oldFact, newFact, "Product~1", "Product")

	if err == nil {
		t.Error("Expected error from callback")
	}
}

func TestDeltaPropagator_PropagateUpdate_MetricsRecorded(t *testing.T) {
	index := NewDependencyIndex()
	index.AddAlphaNode("alpha1", "Product", []string{"price"})

	callback := func(nodeID string, delta *FactDelta) error {
		return nil
	}

	config := DefaultPropagationConfig()
	config.MinFieldsForDelta = 1

	propagator, _ := NewDeltaPropagatorBuilder().
		WithIndex(index).
		WithPropagateCallback(callback).
		WithConfig(config).
		Build()

	oldFact := map[string]interface{}{"price": 100, "name": "A", "stock": 10}
	newFact := map[string]interface{}{"price": 120, "name": "A", "stock": 10}

	err := propagator.PropagateUpdate(oldFact, newFact, "Product~1", "Product")
	if err != nil {
		t.Fatalf("❌ PropagateUpdate failed: %v", err)
	}

	metrics := propagator.GetMetrics()

	if metrics.TotalPropagations == 0 {
		t.Error("Expected at least one propagation recorded")
	}

	if metrics.DeltaPropagations == 0 {
		t.Error("Expected delta propagation recorded")
	}
}

func TestDeltaPropagator_PropagateUpdateWithContext_Timeout(t *testing.T) {
	index := NewDependencyIndex()
	index.AddAlphaNode("alpha1", "Product", []string{"price"})

	callback := func(nodeID string, delta *FactDelta) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	propagator, _ := NewDeltaPropagatorBuilder().
		WithIndex(index).
		WithPropagateCallback(callback).
		Build()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	oldFact := map[string]interface{}{"price": 100}
	newFact := map[string]interface{}{"price": 120}

	err := propagator.PropagateUpdateWithContext(ctx, oldFact, newFact, "Product~1", "Product")

	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestDeltaPropagator_GetConfig(t *testing.T) {
	config := DefaultPropagationConfig()
	config.DeltaThreshold = 0.3

	propagator, _ := NewDeltaPropagatorBuilder().
		WithIndex(NewDependencyIndex()).
		WithConfig(config).
		Build()

	retrievedConfig := propagator.GetConfig()

	if retrievedConfig.DeltaThreshold != 0.3 {
		t.Errorf("Expected threshold 0.3, got %v", retrievedConfig.DeltaThreshold)
	}

	retrievedConfig.DeltaThreshold = 0.5

	retrievedConfig2 := propagator.GetConfig()
	if retrievedConfig2.DeltaThreshold != 0.3 {
		t.Error("Config mutation should not affect original")
	}
}

func TestDeltaPropagator_UpdateConfig(t *testing.T) {
	propagator, _ := NewDeltaPropagatorBuilder().
		WithIndex(NewDependencyIndex()).
		Build()

	newConfig := DefaultPropagationConfig()
	newConfig.DeltaThreshold = 0.7

	err := propagator.UpdateConfig(newConfig)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	retrievedConfig := propagator.GetConfig()
	if retrievedConfig.DeltaThreshold != 0.7 {
		t.Errorf("Expected threshold 0.7, got %v", retrievedConfig.DeltaThreshold)
	}
}

func TestDeltaPropagator_UpdateConfig_Invalid(t *testing.T) {
	propagator, _ := NewDeltaPropagatorBuilder().
		WithIndex(NewDependencyIndex()).
		Build()

	invalidConfig := DefaultPropagationConfig()
	invalidConfig.DeltaThreshold = 2.0

	err := propagator.UpdateConfig(invalidConfig)

	if err == nil {
		t.Error("Expected error with invalid config")
	}
}

func TestDeltaPropagator_ResetMetrics(t *testing.T) {
	index := NewDependencyIndex()
	index.AddAlphaNode("alpha1", "Product", []string{"price"})

	deltaCallback := func(nodeID string, delta *FactDelta) error {
		return nil
	}

	classicCallback := func(factID, factType string, oldFact, newFact map[string]interface{}) error {
		return nil
	}

	propagator, _ := NewDeltaPropagatorBuilder().
		WithIndex(index).
		WithPropagateCallback(deltaCallback).
		WithClassicPropagationCallback(classicCallback).
		Build()

	oldFact := map[string]interface{}{"price": 100}
	newFact := map[string]interface{}{"price": 120}

	err := propagator.PropagateUpdate(oldFact, newFact, "Product~1", "Product")
	if err != nil {
		t.Fatalf("❌ PropagateUpdate failed: %v", err)
	}

	metrics := propagator.GetMetrics()
	if metrics.TotalPropagations == 0 {
		t.Fatal("Expected propagations recorded")
	}

	propagator.ResetMetrics()

	metrics = propagator.GetMetrics()
	if metrics.TotalPropagations != 0 {
		t.Error("Expected metrics to be reset")
	}
}

func TestDeltaPropagator_ConcurrentPropagations(t *testing.T) {
	index := NewDependencyIndex()
	index.AddAlphaNode("alpha1", "Product", []string{"price"})

	callCount := 0
	var mutex sync.Mutex

	callback := func(nodeID string, delta *FactDelta) error {
		mutex.Lock()
		callCount++
		mutex.Unlock()
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	config := DefaultPropagationConfig()
	config.MaxConcurrentPropagations = 5
	config.MinFieldsForDelta = 1

	propagator, _ := NewDeltaPropagatorBuilder().
		WithIndex(index).
		WithPropagateCallback(callback).
		WithConfig(config).
		Build()

	oldFact := map[string]interface{}{"price": 100, "name": "A", "stock": 10}
	newFact := map[string]interface{}{"price": 120, "name": "A", "stock": 10}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = propagator.PropagateUpdate(oldFact, newFact, "Product~1", "Product")
		}()
	}

	wg.Wait()

	mutex.Lock()
	defer mutex.Unlock()
	if callCount != 10 {
		t.Errorf("Expected 10 callback calls, got %d", callCount)
	}
}

// TestDeltaPropagator_ClassicPropagation teste la propagation classique via callback.
func TestDeltaPropagator_ClassicPropagation(t *testing.T) {
	t.Log("🧪 TEST DELTA PROPAGATOR - CLASSIC PROPAGATION")
	t.Log("==============================================")

	index := NewDependencyIndex()

	var receivedFactID, receivedFactType string
	var receivedOldFact, receivedNewFact map[string]interface{}
	callbackCalled := false

	classicCallback := func(factID, factType string, oldFact, newFact map[string]interface{}) error {
		callbackCalled = true
		receivedFactID = factID
		receivedFactType = factType
		receivedOldFact = oldFact
		receivedNewFact = newFact
		return nil
	}

	config := DefaultPropagationConfig()
	config.EnableDeltaPropagation = false

	propagator, err := NewDeltaPropagatorBuilder().
		WithIndex(index).
		WithConfig(config).
		WithClassicPropagationCallback(classicCallback).
		Build()

	if err != nil {
		t.Fatalf("❌ Failed to build propagator: %v", err)
	}

	oldFact := map[string]interface{}{"price": 100, "name": "Product A"}
	newFact := map[string]interface{}{"price": 120, "name": "Product B"}

	err = propagator.PropagateUpdate(oldFact, newFact, "Product~1", "Product")

	if err != nil {
		t.Fatalf("❌ PropagateUpdate failed: %v", err)
	}

	if !callbackCalled {
		t.Fatal("❌ Classic callback was not called")
	}

	if receivedFactID != "Product~1" {
		t.Errorf("❌ Expected factID 'Product~1', got '%s'", receivedFactID)
	}

	if receivedFactType != "Product" {
		t.Errorf("❌ Expected factType 'Product', got '%s'", receivedFactType)
	}

	if receivedOldFact["price"] != 100 {
		t.Errorf("❌ Expected old price 100, got %v", receivedOldFact["price"])
	}

	if receivedNewFact["price"] != 120 {
		t.Errorf("❌ Expected new price 120, got %v", receivedNewFact["price"])
	}

	metrics := propagator.GetMetrics()
	if metrics.ClassicPropagations != 1 {
		t.Errorf("❌ Expected 1 classic propagation, got %d", metrics.ClassicPropagations)
	}

	t.Log("✅ Classic propagation callback executed successfully")
}

// TestDeltaPropagator_ClassicPropagationError teste la gestion d'erreur du callback classique.
func TestDeltaPropagator_ClassicPropagationError(t *testing.T) {
	t.Log("🧪 TEST DELTA PROPAGATOR - CLASSIC PROPAGATION ERROR")
	t.Log("====================================================")

	index := NewDependencyIndex()

	expectedErr := fmt.Errorf("simulated classic propagation error")
	classicCallback := func(factID, factType string, oldFact, newFact map[string]interface{}) error {
		return expectedErr
	}

	config := DefaultPropagationConfig()
	config.EnableDeltaPropagation = false

	propagator, _ := NewDeltaPropagatorBuilder().
		WithIndex(index).
		WithConfig(config).
		WithClassicPropagationCallback(classicCallback).
		Build()

	oldFact := map[string]interface{}{"price": 100}
	newFact := map[string]interface{}{"price": 120}

	err := propagator.PropagateUpdate(oldFact, newFact, "Product~1", "Product")

	if err == nil {
		t.Fatal("❌ Expected error from classic propagation")
	}

	if err != expectedErr {
		t.Errorf("❌ Expected specific error, got: %v", err)
	}

	t.Log("✅ Classic propagation error handled correctly")
}

// TestDeltaPropagator_FallbackToClassic teste le fallback vers propagation classique.
func TestDeltaPropagator_FallbackToClassic(t *testing.T) {
	t.Log("🧪 TEST DELTA PROPAGATOR - FALLBACK TO CLASSIC")
	t.Log("================================================")

	index := NewDependencyIndex()
	index.AddAlphaNode("alpha1", "Product", []string{"price"})

	deltaCallCount := 0
	classicCallCount := 0

	deltaCallback := func(nodeID string, delta *FactDelta) error {
		deltaCallCount++
		return nil
	}

	classicCallback := func(factID, factType string, oldFact, newFact map[string]interface{}) error {
		classicCallCount++
		return nil
	}

	config := DefaultPropagationConfig()
	config.DeltaThreshold = 0.3  // 30% threshold
	config.MinFieldsForDelta = 1 // Accept delta for facts with 1+ fields
	config.EnableDeltaPropagation = true

	propagator, _ := NewDeltaPropagatorBuilder().
		WithIndex(index).
		WithConfig(config).
		WithPropagateCallback(deltaCallback).
		WithClassicPropagationCallback(classicCallback).
		Build()

	// Cas 1: Petit changement (1/5 = 20% < 30%) -> propagation delta
	oldFact1 := map[string]interface{}{"price": 100, "name": "A", "stock": 10, "category": "X", "active": true}
	newFact1 := map[string]interface{}{"price": 120, "name": "A", "stock": 10, "category": "X", "active": true}

	err := propagator.PropagateUpdate(oldFact1, newFact1, "Product~1", "Product")
	if err != nil {
		t.Fatalf("❌ Propagation failed: %v", err)
	}

	if deltaCallCount != 1 {
		t.Errorf("❌ Expected 1 delta call, got %d", deltaCallCount)
	}

	if classicCallCount != 0 {
		t.Errorf("❌ Expected 0 classic call, got %d", classicCallCount)
	}

	// Cas 2: Grand changement (2/5 = 40% > 30%) -> fallback classique
	oldFact2 := map[string]interface{}{"price": 100, "name": "A", "stock": 10, "category": "X", "active": true}
	newFact2 := map[string]interface{}{"price": 200, "name": "B", "stock": 10, "category": "X", "active": true}

	err = propagator.PropagateUpdate(oldFact2, newFact2, "Product~2", "Product")
	if err != nil {
		t.Fatalf("❌ Propagation failed: %v", err)
	}

	if deltaCallCount != 1 { // Toujours 1 (pas augmenté)
		t.Errorf("❌ Expected 1 delta call, got %d", deltaCallCount)
	}

	if classicCallCount != 1 { // Devrait être 1 maintenant
		t.Errorf("❌ Expected 1 classic call, got %d", classicCallCount)
	}

	metrics := propagator.GetMetrics()
	if metrics.DeltaPropagations != 1 {
		t.Errorf("❌ Expected 1 delta propagation, got %d", metrics.DeltaPropagations)
	}

	if metrics.ClassicPropagations != 1 {
		t.Errorf("❌ Expected 1 classic propagation, got %d", metrics.ClassicPropagations)
	}

	t.Log("✅ Fallback to classic propagation works correctly")
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestDeltaPropagator_recordFallbackReason teste la fonction privée recordFallbackReason
func TestDeltaPropagator_recordFallbackReason(t *testing.T) {
	t.Run("fallback reason: fields", func(t *testing.T) {
		config := DefaultPropagationConfig()
		config.MinFieldsForDelta = 5 // Nécessite au moins 5 champs

		propagator, _ := NewDeltaPropagatorBuilder().
			WithIndex(NewDependencyIndex()).
			WithDetector(NewDeltaDetector()).
			WithConfig(config).
			Build()

		// Delta avec seulement 3 champs (< MinFieldsForDelta)
		delta := NewFactDelta("fact1", "Product")
		delta.FieldCount = 3 // < 5
		delta.AddFieldChange("price", 100, 120)
		delta.AddFieldChange("stock", 10, 5)
		delta.AddFieldChange("name", "A", "B")

		// Appeler recordFallbackReason
		propagator.recordFallbackReason(delta, 2)

		// Vérifier que la métrique "fields" a été enregistrée
		snapshot := propagator.metrics.GetSnapshot()
		if snapshot.FallbacksDueToFields != 1 {
			t.Errorf("❌ Expected FallbacksDueToFields=1, got %d", snapshot.FallbacksDueToFields)
		}

		t.Log("✅ Fallback reason 'fields' enregistré correctement")
	})

	t.Run("fallback reason: ratio", func(t *testing.T) {
		config := DefaultPropagationConfig()
		config.MinFieldsForDelta = 2
		config.DeltaThreshold = 0.3 // Maximum 30% de changements

		propagator, _ := NewDeltaPropagatorBuilder().
			WithIndex(NewDependencyIndex()).
			WithDetector(NewDeltaDetector()).
			WithConfig(config).
			Build()

		// Delta avec 4/5 champs changés = 80% > 30%
		delta := NewFactDelta("fact1", "Product")
		delta.FieldCount = 5
		delta.AddFieldChange("field1", 1, 2)
		delta.AddFieldChange("field2", 3, 4)
		delta.AddFieldChange("field3", 5, 6)
		delta.AddFieldChange("field4", 7, 8)

		propagator.recordFallbackReason(delta, 2)

		snapshot := propagator.metrics.GetSnapshot()
		if snapshot.FallbacksDueToRatio != 1 {
			t.Errorf("❌ Expected FallbacksDueToRatio=1, got %d", snapshot.FallbacksDueToRatio)
		}

		t.Log("✅ Fallback reason 'ratio' enregistré correctement")
	})

	t.Run("fallback reason: nodes", func(t *testing.T) {
		config := DefaultPropagationConfig()
		config.MinFieldsForDelta = 2
		config.DeltaThreshold = 0.8         // 80% autorisé
		config.MaxAffectedNodesForDelta = 5 // Maximum 5 nœuds affectés

		propagator, _ := NewDeltaPropagatorBuilder().
			WithIndex(NewDependencyIndex()).
			WithDetector(NewDeltaDetector()).
			WithConfig(config).
			Build()

		// Delta avec peu de changements mais beaucoup de nœuds affectés
		delta := NewFactDelta("fact1", "Product")
		delta.FieldCount = 10
		delta.AddFieldChange("price", 100, 120)

		// 10 nœuds affectés > MaxAffectedNodesForDelta (5)
		propagator.recordFallbackReason(delta, 10)

		snapshot := propagator.metrics.GetSnapshot()
		if snapshot.FallbacksDueToNodes != 1 {
			t.Errorf("❌ Expected FallbacksDueToNodes=1, got %d", snapshot.FallbacksDueToNodes)
		}

		t.Log("✅ Fallback reason 'nodes' enregistré correctement")
	})

	t.Run("fallback reason: pk (primary key)", func(t *testing.T) {
		config := DefaultPropagationConfig()
		config.MinFieldsForDelta = 2
		config.DeltaThreshold = 0.8
		config.MaxAffectedNodesForDelta = 100
		config.PrimaryKeyFields = []string{"id", "uuid"} // Champs PK

		propagator, _ := NewDeltaPropagatorBuilder().
			WithIndex(NewDependencyIndex()).
			WithDetector(NewDeltaDetector()).
			WithConfig(config).
			Build()

		// Delta avec changement de clé primaire
		delta := NewFactDelta("fact1", "Product")
		delta.FieldCount = 10
		delta.AddFieldChange("id", "123", "456") // PK change!
		delta.AddFieldChange("name", "A", "B")

		propagator.recordFallbackReason(delta, 2)

		snapshot := propagator.metrics.GetSnapshot()
		if snapshot.FallbacksDueToPK != 1 {
			t.Errorf("❌ Expected FallbacksDueToPK=1, got %d", snapshot.FallbacksDueToPK)
		}

		t.Log("✅ Fallback reason 'pk' enregistré correctement")
	})

	t.Run("no fallback - delta should be used", func(t *testing.T) {
		config := DefaultPropagationConfig()
		config.MinFieldsForDelta = 2
		config.DeltaThreshold = 0.5 // 50% autorisé
		config.MaxAffectedNodesForDelta = 100

		propagator, _ := NewDeltaPropagatorBuilder().
			WithIndex(NewDependencyIndex()).
			WithDetector(NewDeltaDetector()).
			WithConfig(config).
			Build()

		// Delta qui devrait passer tous les checks
		delta := NewFactDelta("fact1", "Product")
		delta.FieldCount = 10
		delta.AddFieldChange("price", 100, 120)
		delta.AddFieldChange("stock", 10, 5)

		initialSnapshot := propagator.metrics.GetSnapshot()
		initialTotal := initialSnapshot.FallbacksDueToFields +
			initialSnapshot.FallbacksDueToRatio +
			initialSnapshot.FallbacksDueToNodes +
			initialSnapshot.FallbacksDueToPK +
			initialSnapshot.FallbacksDueToError

		// Appeler recordFallbackReason avec des conditions valides
		propagator.recordFallbackReason(delta, 3)

		// Aucun fallback ne devrait être enregistré
		newSnapshot := propagator.metrics.GetSnapshot()
		newTotal := newSnapshot.FallbacksDueToFields +
			newSnapshot.FallbacksDueToRatio +
			newSnapshot.FallbacksDueToNodes +
			newSnapshot.FallbacksDueToPK +
			newSnapshot.FallbacksDueToError

		if newTotal != initialTotal {
			t.Errorf("❌ No fallback should be recorded, initial=%d, new=%d", initialTotal, newTotal)
		}

		t.Log("✅ Aucun fallback enregistré pour delta valide")
	})

	t.Run("multiple fallback reasons", func(t *testing.T) {
		config := DefaultPropagationConfig()
		config.MinFieldsForDelta = 5

		propagator, _ := NewDeltaPropagatorBuilder().
			WithIndex(NewDependencyIndex()).
			WithDetector(NewDeltaDetector()).
			WithConfig(config).
			Build()

		// Premier delta: fallback sur fields
		delta1 := NewFactDelta("fact1", "Product")
		delta1.FieldCount = 3
		delta1.AddFieldChange("a", 1, 2)
		propagator.recordFallbackReason(delta1, 2)

		// Deuxième delta: fallback sur fields aussi
		delta2 := NewFactDelta("fact2", "Product")
		delta2.FieldCount = 2
		delta2.AddFieldChange("b", 1, 2)
		propagator.recordFallbackReason(delta2, 1)

		snapshot := propagator.metrics.GetSnapshot()
		if snapshot.FallbacksDueToFields != 2 {
			t.Errorf("❌ Expected FallbacksDueToFields=2, got %d", snapshot.FallbacksDueToFields)
		}

		t.Log("✅ Multiples fallbacks enregistrés correctement")
	})
}
