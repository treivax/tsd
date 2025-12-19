// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package api

import "testing"

func TestDefaultConfig(t *testing.T) {
	t.Log("🧪 TEST DEFAULT CONFIG")

	config := DefaultConfig()
	if config == nil {
		t.Fatal("❌ Config ne devrait pas être nil")
	}

	if config.LogLevel != LogLevelInfo {
		t.Errorf("❌ LogLevel par défaut attendu: Info, reçu: %v", config.LogLevel)
	}

	if !config.EnableMetrics {
		t.Error("❌ EnableMetrics devrait être true par défaut")
	}

	if config.XupleSpaceDefaults == nil {
		t.Fatal("❌ XupleSpaceDefaults ne devrait pas être nil")
	}

	if config.XupleSpaceDefaults.Selection != SelectionFIFO {
		t.Errorf("❌ Selection par défaut attendue: FIFO, reçue: %s",
			config.XupleSpaceDefaults.Selection)
	}

	t.Log("✅ Configuration par défaut correcte")
}

func TestConfigValidate_Valid(t *testing.T) {
	t.Log("🧪 TEST CONFIG VALIDATE VALID")

	config := DefaultConfig()
	err := config.Validate()
	if err != nil {
		t.Errorf("❌ Configuration par défaut devrait être valide, erreur: %v", err)
	}

	t.Log("✅ Validation config valide réussie")
}

func TestConfigValidate_InvalidTransactionTimeout(t *testing.T) {
	t.Log("🧪 TEST CONFIG VALIDATE INVALID TRANSACTION TIMEOUT")

	config := DefaultConfig()
	config.TransactionTimeout = -1

	err := config.Validate()
	if err == nil {
		t.Fatal("❌ Devrait retourner une erreur pour TransactionTimeout négatif")
	}

	if configErr, ok := err.(*ConfigError); ok {
		if configErr.Field != "TransactionTimeout" {
			t.Errorf("❌ Field attendu: TransactionTimeout, reçu: %s", configErr.Field)
		}
	} else {
		t.Errorf("❌ Devrait être une *ConfigError, reçu: %T", err)
	}

	t.Log("✅ Validation TransactionTimeout négatif détectée")
}

func TestConfigValidate_InvalidMaxFactsInMemory(t *testing.T) {
	t.Log("🧪 TEST CONFIG VALIDATE INVALID MAX FACTS")

	config := DefaultConfig()
	config.MaxFactsInMemory = -100

	err := config.Validate()
	if err == nil {
		t.Fatal("❌ Devrait retourner une erreur pour MaxFactsInMemory négatif")
	}

	t.Log("✅ Validation MaxFactsInMemory négatif détectée")
}

func TestConfigValidate_InvalidSelectionPolicy(t *testing.T) {
	t.Log("🧪 TEST CONFIG VALIDATE INVALID SELECTION POLICY")

	config := DefaultConfig()
	config.XupleSpaceDefaults.Selection = "invalid"

	err := config.Validate()
	if err == nil {
		t.Fatal("❌ Devrait retourner une erreur pour politique de sélection invalide")
	}

	t.Log("✅ Validation politique sélection invalide détectée")
}

func TestConfigValidate_InvalidConsumptionPolicy(t *testing.T) {
	t.Log("🧪 TEST CONFIG VALIDATE INVALID CONSUMPTION POLICY")

	config := DefaultConfig()
	config.XupleSpaceDefaults.Consumption = "invalid"

	err := config.Validate()
	if err == nil {
		t.Fatal("❌ Devrait retourner une erreur pour politique de consommation invalide")
	}

	t.Log("✅ Validation politique consommation invalide détectée")
}

func TestConfigValidate_InvalidRetentionPolicy(t *testing.T) {
	t.Log("🧪 TEST CONFIG VALIDATE INVALID RETENTION POLICY")

	config := DefaultConfig()
	config.XupleSpaceDefaults.Retention = "invalid"

	err := config.Validate()
	if err == nil {
		t.Fatal("❌ Devrait retourner une erreur pour politique de rétention invalide")
	}

	t.Log("✅ Validation politique rétention invalide détectée")
}

func TestConfigValidate_EmptyPolicies(t *testing.T) {
	t.Log("🧪 TEST CONFIG VALIDATE EMPTY POLICIES")

	config := &Config{
		LogLevel:      LogLevelInfo,
		EnableMetrics: true,
		XupleSpaceDefaults: &XupleSpaceDefaults{
			Selection:   "",
			Consumption: "",
			Retention:   "",
		},
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("❌ Devrait accepter les politiques vides (défauts), erreur: %v", err)
	}

	if config.XupleSpaceDefaults.Selection != SelectionFIFO {
		t.Errorf("❌ Politique vide devrait devenir FIFO, reçu: %s",
			config.XupleSpaceDefaults.Selection)
	}

	t.Log("✅ Validation politiques vides (défauts appliqués) réussie")
}

func TestConfigValidate_RetentionDurationRequired(t *testing.T) {
	t.Log("🧪 TEST CONFIG VALIDATE RETENTION DURATION REQUIRED")

	config := DefaultConfig()
	config.XupleSpaceDefaults.Retention = RetentionDuration
	config.XupleSpaceDefaults.RetentionDuration = 0

	err := config.Validate()
	if err == nil {
		t.Fatal("❌ Devrait exiger RetentionDuration > 0 quand Retention = duration")
	}

	t.Log("✅ Validation RetentionDuration requis détectée")
}

func TestConfigValidate_NegativeMaxSize(t *testing.T) {
	t.Log("🧪 TEST CONFIG VALIDATE NEGATIVE MAX SIZE")

	config := DefaultConfig()
	config.XupleSpaceDefaults.MaxSize = -50

	err := config.Validate()
	if err == nil {
		t.Fatal("❌ Devrait retourner une erreur pour MaxSize négatif")
	}

	t.Log("✅ Validation MaxSize négatif détectée")
}

func TestSelectionPolicyConstants(t *testing.T) {
	t.Log("🧪 TEST SELECTION POLICY CONSTANTS")

	tests := []struct {
		policy   SelectionPolicy
		expected string
	}{
		{SelectionFIFO, "fifo"},
		{SelectionLIFO, "lifo"},
		{SelectionRandom, "random"},
	}

	for _, tt := range tests {
		if string(tt.policy) != tt.expected {
			t.Errorf("❌ Constante %s attendue: %s, reçue: %s",
				tt.expected, tt.expected, tt.policy)
		}
	}

	t.Log("✅ Constantes politique sélection correctes")
}

func TestConsumptionPolicyConstants(t *testing.T) {
	t.Log("🧪 TEST CONSUMPTION POLICY CONSTANTS")

	tests := []struct {
		policy   ConsumptionPolicy
		expected string
	}{
		{ConsumptionOnce, "once"},
		{ConsumptionPerAgent, "per-agent"},
	}

	for _, tt := range tests {
		if string(tt.policy) != tt.expected {
			t.Errorf("❌ Constante %s attendue: %s, reçue: %s",
				tt.expected, tt.expected, tt.policy)
		}
	}

	t.Log("✅ Constantes politique consommation correctes")
}

func TestRetentionPolicyConstants(t *testing.T) {
	t.Log("🧪 TEST RETENTION POLICY CONSTANTS")

	tests := []struct {
		policy   RetentionPolicy
		expected string
	}{
		{RetentionUnlimited, "unlimited"},
		{RetentionDuration, "duration"},
	}

	for _, tt := range tests {
		if string(tt.policy) != tt.expected {
			t.Errorf("❌ Constante %s attendue: %s, reçue: %s",
				tt.expected, tt.expected, tt.policy)
		}
	}

	t.Log("✅ Constantes politique rétention correctes")
}

func TestLogLevelConstants(t *testing.T) {
	t.Log("🧪 TEST LOG LEVEL CONSTANTS")

	levels := []LogLevel{
		LogLevelSilent,
		LogLevelError,
		LogLevelWarn,
		LogLevelInfo,
		LogLevelDebug,
	}

	for i, level := range levels {
		if int(level) != i {
			t.Errorf("❌ LogLevel %d devrait avoir valeur %d, reçu: %d",
				i, i, int(level))
		}
	}

	t.Log("✅ Constantes niveau log correctes")
}
