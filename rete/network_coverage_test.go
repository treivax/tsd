// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"
	"time"
)

// Tests for network_validator.go
func TestValidateNetwork_ValidStructure(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	err := network.ValidateNetwork()
	if err != nil {
		t.Errorf("ValidateNetwork() failed on empty network: %v", err)
	}
}
func TestValidateNetwork_NilRootNode(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.RootNode = nil
	err := network.ValidateNetwork()
	if err == nil {
		t.Error("Expected error for nil RootNode")
	}
}
func TestValidateNetwork_NilStorage(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.Storage = nil
	err := network.ValidateNetwork()
	if err == nil {
		t.Error("Expected error for nil Storage")
	}
}
func TestValidateNetwork_NilLifecycleManager(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.LifecycleManager = nil
	err := network.ValidateNetwork()
	if err == nil {
		t.Error("Expected error for nil LifecycleManager")
	}
}
func TestValidateRule_NonexistentRule(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	err := network.ValidateRule("nonexistent_rule")
	if err == nil {
		t.Error("Expected error for nonexistent rule")
	}
}
func TestValidateFactIntegrity_FactNotFound(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	err := network.ValidateFactIntegrity("nonexistent_fact")
	if err == nil {
		t.Error("Expected error for nonexistent fact")
	}
}
func TestValidateMemoryConsistency_EmptyNetwork(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	err := network.ValidateMemoryConsistency()
	if err != nil {
		t.Errorf("ValidateMemoryConsistency() failed: %v", err)
	}
}

// Tests for network_builder.go
func TestNewReteNetwork_Initialization(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	if network == nil {
		t.Fatal("NewReteNetwork() returned nil")
	}
	if network.RootNode == nil {
		t.Error("RootNode not initialized")
	}
	if network.Storage != storage {
		t.Error("Storage not set correctly")
	}
	if network.TypeNodes == nil {
		t.Error("TypeNodes map not initialized")
	}
	if network.AlphaNodes == nil {
		t.Error("AlphaNodes map not initialized")
	}
	if network.BetaNodes == nil {
		t.Error("BetaNodes map not initialized")
	}
	if network.TerminalNodes == nil {
		t.Error("TerminalNodes map not initialized")
	}
	if network.LifecycleManager == nil {
		t.Error("LifecycleManager not initialized")
	}
	if network.AlphaSharingManager == nil {
		t.Error("AlphaSharingManager not initialized")
	}
	if network.BetaSharingRegistry == nil {
		t.Error("BetaSharingRegistry not initialized")
	}
	if network.ArithmeticResultCache == nil {
		t.Error("ArithmeticResultCache not initialized")
	}
	if network.Config == nil {
		t.Error("Config not initialized")
	}
	if network.GetLogger() == nil {
		t.Error("Logger not initialized")
	}
}
func TestNewReteNetworkWithConfig_NilConfig(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetworkWithConfig(storage, nil)
	if network == nil {
		t.Fatal("NewReteNetworkWithConfig() returned nil")
	}
	if network.Config == nil {
		t.Error("Config should be initialized with default when nil passed")
	}
}
func TestNewReteNetworkWithConfig_CustomConfig(t *testing.T) {
	storage := NewMemoryStorage()
	config := &ChainPerformanceConfig{
		HashCacheMaxSize:     500,
		BetaHashCacheMaxSize: 500,
	}
	network := NewReteNetworkWithConfig(storage, config)
	if network == nil {
		t.Fatal("NewReteNetworkWithConfig() returned nil")
	}
	if network.Config.BetaHashCacheMaxSize != 500 {
		t.Errorf("BetaHashCacheMaxSize = %d, want 500", network.Config.BetaHashCacheMaxSize)
	}
}

// Tests for network_manager.go
func TestSubmitFact_Basic(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	fact := &Fact{
		ID:   "test_fact",
		Type: "Person",
		Fields: map[string]interface{}{
			"age": 25,
		},
	}
	err := network.SubmitFact(fact)
	if err != nil {
		t.Errorf("SubmitFact() failed: %v", err)
	}
	storedFact := storage.GetFact(fact.GetInternalID())
	if storedFact == nil {
		t.Error("Fact not found in storage")
	}
}
func TestRemoveFact_Basic(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	fact := &Fact{
		ID:   "test_fact",
		Type: "Person",
		Fields: map[string]interface{}{
			"age": 25,
		},
	}
	storage.AddFact(fact)
	err := network.RemoveFact(fact.GetInternalID())
	if err != nil {
		t.Errorf("RemoveFact() failed: %v", err)
	}
	storedFact := storage.GetFact(fact.GetInternalID())
	if storedFact != nil {
		t.Error("Fact should be removed from storage")
	}
}
func TestRetractFact_FactNotFound(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	err := network.RetractFact("nonexistent_fact")
	if err == nil {
		t.Error("Expected error for nonexistent fact")
	}
}
func TestReset_ClearsNetwork(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Add some data
	network.Types = append(network.Types, TypeDefinition{Name: "Test"})
	network.Reset()
	if len(network.TypeNodes) != 0 {
		t.Error("TypeNodes should be empty after reset")
	}
	if len(network.Types) != 0 {
		t.Error("Types should be empty after reset")
	}
}
func TestClearMemory_Basic(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.ClearMemory()
	// Should not panic
}
func TestGarbageCollect_Basic(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.GarbageCollect()
	if len(network.TypeNodes) != 0 {
		t.Error("TypeNodes should be empty after GC")
	}
}
func TestSubmitFactsFromGrammar_Empty(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	err := network.SubmitFactsFromGrammar([]map[string]interface{}{})
	if err != nil {
		t.Errorf("SubmitFactsFromGrammar() with empty list failed: %v", err)
	}
}
func TestSubmitFactsFromGrammar_Single(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.SubmissionTimeout = 5 * time.Second
	facts := []map[string]interface{}{
		{
			"id":       "fact1",
			"reteType": "Person",
			"age":      25,
		},
	}
	err := network.SubmitFactsFromGrammar(facts)
	if err != nil {
		t.Errorf("SubmitFactsFromGrammar() failed: %v", err)
	}
	storedFact := storage.GetFact("Person~fact1")
	if storedFact == nil {
		t.Error("Fact not found in storage")
	}
}
func TestRepropagateExistingFact_TypeNotFound(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	fact := &Fact{
		ID:   "test_fact",
		Type: "UnknownType",
		Fields: map[string]interface{}{
			"age": 25,
		},
	}
	err := network.RepropagateExistingFact(fact)
	if err == nil {
		t.Error("Expected error for unknown type")
	}
}

// Tests for network.go (getters/setters)
func TestGetChainMetrics(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	metrics := network.GetChainMetrics()
	if metrics == nil {
		t.Error("GetChainMetrics() returned nil")
	}
}
func TestGetBetaSharingStats(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	stats := network.GetBetaSharingStats()
	if stats == nil {
		t.Error("GetBetaSharingStats() returned nil")
	}
}
func TestGetConfig(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	config := network.GetConfig()
	if config == nil {
		t.Error("GetConfig() returned nil")
	}
}
func TestSetLogger(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	logger := NewLogger(LogLevelDebug, nil)
	network.SetLogger(logger)
	if network.GetLogger() != logger {
		t.Error("Logger not set correctly")
	}
}
func TestResetChainMetrics(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.ResetChainMetrics()
	metrics := network.GetChainMetrics()
	if metrics == nil {
		t.Error("Metrics should still exist after reset")
	}
}
func TestGetTypeDefinition(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.Types = append(network.Types, TypeDefinition{Name: "Person"})
	typeDef := network.GetTypeDefinition("Person")
	if typeDef == nil {
		t.Error("GetTypeDefinition() returned nil")
		return
	}
	if typeDef.Name != "Person" {
		t.Errorf("Type name = %s, want Person", typeDef.Name)
	}
	nonExistent := network.GetTypeDefinition("NonExistent")
	if nonExistent != nil {
		t.Error("Should return nil for non-existent type")
	}
}
func TestGetNetworkStats(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	stats := network.GetNetworkStats()
	if stats == nil {
		t.Error("GetNetworkStats() returned nil")
		return
	}
	if _, ok := stats["type_nodes"]; !ok {
		t.Error("Stats should include type_nodes")
	}
	if _, ok := stats["alpha_nodes"]; !ok {
		t.Error("Stats should include alpha_nodes")
	}
}
func TestGetRuleInfo(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	info := network.GetRuleInfo("test_rule")
	if info == nil {
		t.Error("GetRuleInfo() returned nil")
		return
	}
	if info.RuleID != "test_rule" {
		t.Errorf("Rule ID = %s, want test_rule", info.RuleID)
	}
}
func TestGetSetTransaction(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	tx := network.GetTransaction()
	if tx != nil {
		t.Error("Transaction should be nil initially")
	}
	// Test setting nil transaction (should be safe)
	network.SetTransaction(nil)
	tx = network.GetTransaction()
	if tx != nil {
		t.Error("Transaction should still be nil")
	}
}
