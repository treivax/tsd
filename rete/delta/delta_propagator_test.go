// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"context"
	"errors"
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

	// Classic propagation not yet implemented
	if err == nil {
		t.Error("Expected error for classic propagation not implemented")
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

	callback := func(nodeID string, delta *FactDelta) error {
		return nil
	}

	propagator, _ := NewDeltaPropagatorBuilder().
		WithIndex(index).
		WithPropagateCallback(callback).
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
