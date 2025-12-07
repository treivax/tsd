// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package config

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	// Test Storage config
	if config.Storage.Type != "memory" {
		t.Errorf("Storage.Type = %s, want memory", config.Storage.Type)
	}
	if config.Storage.Timeout != 30*time.Second {
		t.Errorf("Storage.Timeout = %v, want 30s", config.Storage.Timeout)
	}

	// Test Network config
	if config.Network.MaxNodes != 1000 {
		t.Errorf("Network.MaxNodes = %d, want 1000", config.Network.MaxNodes)
	}
	if config.Network.MaxFactsPerNode != 10000 {
		t.Errorf("Network.MaxFactsPerNode = %d, want 10000", config.Network.MaxFactsPerNode)
	}
	if config.Network.GCInterval != 5*time.Minute {
		t.Errorf("Network.GCInterval = %v, want 5m", config.Network.GCInterval)
	}

	// Test Logger config
	if config.Logger.Level != "info" {
		t.Errorf("Logger.Level = %s, want info", config.Logger.Level)
	}
	if config.Logger.Format != "text" {
		t.Errorf("Logger.Format = %s, want text", config.Logger.Format)
	}
	if config.Logger.Output != "stdout" {
		t.Errorf("Logger.Output = %s, want stdout", config.Logger.Output)
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid default config",
			config:    DefaultConfig(),
			wantError: false,
		},
		{
			name: "valid memory storage",
			config: &Config{
				Storage: StorageConfig{Type: "memory"},
				Logger:  LoggerConfig{Level: "info"},
			},
			wantError: false,
		},

		{
			name: "invalid storage type",
			config: &Config{
				Storage: StorageConfig{Type: "invalid"},
				Logger:  LoggerConfig{Level: "info"},
			},
			wantError: true,
			errorMsg:  "storage.type",
		},
		{
			name: "empty storage type",
			config: &Config{
				Storage: StorageConfig{Type: ""},
				Logger:  LoggerConfig{Level: "info"},
			},
			wantError: true,
			errorMsg:  "storage.type",
		},
		{
			name: "valid logger level - debug",
			config: &Config{
				Storage: StorageConfig{Type: "memory"},
				Logger:  LoggerConfig{Level: "debug"},
			},
			wantError: false,
		},
		{
			name: "valid logger level - warn",
			config: &Config{
				Storage: StorageConfig{Type: "memory"},
				Logger:  LoggerConfig{Level: "warn"},
			},
			wantError: false,
		},
		{
			name: "valid logger level - error",
			config: &Config{
				Storage: StorageConfig{Type: "memory"},
				Logger:  LoggerConfig{Level: "error"},
			},
			wantError: false,
		},
		{
			name: "invalid logger level",
			config: &Config{
				Storage: StorageConfig{Type: "memory"},
				Logger:  LoggerConfig{Level: "trace"},
			},
			wantError: true,
			errorMsg:  "logger.level",
		},
		{
			name: "empty logger level",
			config: &Config{
				Storage: StorageConfig{Type: "memory"},
				Logger:  LoggerConfig{Level: ""},
			},
			wantError: true,
			errorMsg:  "logger.level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantError {
				if err == nil {
					t.Errorf("Validate() error = nil, want error containing %q", tt.errorMsg)
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Validate() error = %v, want error containing %q", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() error = %v, want nil", err)
				}
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ValidationError
		expected string
	}{
		{
			name: "storage type error",
			err: &ValidationError{
				Field:   "storage.type",
				Value:   "invalid",
				Message: "must be 'memory' (in-memory storage only)",
			},
			expected: "config validation error on field 'storage.type': must be 'memory' (in-memory storage only)",
		},
		{
			name: "logger level error",
			err: &ValidationError{
				Field:   "logger.level",
				Value:   "trace",
				Message: "must be one of: debug, info, warn, error",
			},
			expected: "config validation error on field 'logger.level': must be one of: debug, info, warn, error",
		},
		{
			name: "empty message",
			err: &ValidationError{
				Field:   "some.field",
				Value:   123,
				Message: "",
			},
			expected: "config validation error on field 'some.field': ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("Error() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestStorageConfig(t *testing.T) {
	storage := StorageConfig{
		Type:    "memory",
		Timeout: 10 * time.Second,
	}

	if storage.Type != "memory" {
		t.Errorf("Type = %s, want memory", storage.Type)
	}
	if storage.Timeout != 10*time.Second {
		t.Errorf("Timeout = %v, want 10s", storage.Timeout)
	}
}

func TestNetworkConfig(t *testing.T) {
	network := NetworkConfig{
		MaxNodes:        5000,
		MaxFactsPerNode: 20000,
		GCInterval:      10 * time.Minute,
	}

	if network.MaxNodes != 5000 {
		t.Errorf("MaxNodes = %d, want 5000", network.MaxNodes)
	}
	if network.MaxFactsPerNode != 20000 {
		t.Errorf("MaxFactsPerNode = %d, want 20000", network.MaxFactsPerNode)
	}
	if network.GCInterval != 10*time.Minute {
		t.Errorf("GCInterval = %v, want 10m", network.GCInterval)
	}
}

func TestLoggerConfig(t *testing.T) {
	logger := LoggerConfig{
		Level:  "debug",
		Format: "json",
		Output: "/var/log/rete.log",
	}

	if logger.Level != "debug" {
		t.Errorf("Level = %s, want debug", logger.Level)
	}
	if logger.Format != "json" {
		t.Errorf("Format = %s, want json", logger.Format)
	}
	if logger.Output != "/var/log/rete.log" {
		t.Errorf("Output = %s, want /var/log/rete.log", logger.Output)
	}
}

func TestConfig_JSONMarshaling(t *testing.T) {
	original := DefaultConfig()
	original.Storage.Type = "memory"
	original.Network.MaxNodes = 2000
	original.Logger.Level = "debug"

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	// Unmarshal back
	var unmarshaled Config
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify values are preserved
	if unmarshaled.Storage.Type != original.Storage.Type {
		t.Errorf("Storage.Type not preserved: got %s, want %s", unmarshaled.Storage.Type, original.Storage.Type)
	}
	if unmarshaled.Network.MaxNodes != original.Network.MaxNodes {
		t.Errorf("Network.MaxNodes not preserved: got %d, want %d", unmarshaled.Network.MaxNodes, original.Network.MaxNodes)
	}
	if unmarshaled.Logger.Level != original.Logger.Level {
		t.Errorf("Logger.Level not preserved: got %s, want %s", unmarshaled.Logger.Level, original.Logger.Level)
	}
}

func TestConfig_MultipleValidationErrors(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		checkFunc func(*testing.T, error)
	}{
		{
			name: "invalid storage type first",
			config: &Config{
				Storage: StorageConfig{Type: "invalid"},
				Logger:  LoggerConfig{Level: "invalid"},
			},
			checkFunc: func(t *testing.T, err error) {
				if err == nil {
					t.Fatal("Expected error, got nil")
				}
				// Should fail on storage.type first
				if !strings.Contains(err.Error(), "storage.type") {
					t.Errorf("Expected error about storage.type, got: %v", err)
				}
			},
		},
		{
			name: "invalid logger level when storage valid",
			config: &Config{
				Storage: StorageConfig{Type: "memory"},
				Logger:  LoggerConfig{Level: "invalid"},
			},
			checkFunc: func(t *testing.T, err error) {
				if err == nil {
					t.Fatal("Expected error, got nil")
				}
				if !strings.Contains(err.Error(), "logger.level") {
					t.Errorf("Expected error about logger.level, got: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			tt.checkFunc(t, err)
		})
	}
}

func TestConfig_EdgeCases(t *testing.T) {
	t.Run("zero values", func(t *testing.T) {
		config := &Config{}
		err := config.Validate()
		if err == nil {
			t.Error("Empty config should fail validation")
		}
	})

	t.Run("nil config panic check", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for nil config")
			}
		}()
		var config *Config
		_ = config.Validate()
	})

	t.Run("very large timeout", func(t *testing.T) {
		config := DefaultConfig()
		config.Storage.Timeout = 24 * time.Hour
		err := config.Validate()
		if err != nil {
			t.Errorf("Large timeout should be valid, got error: %v", err)
		}
	})

	t.Run("zero timeout", func(t *testing.T) {
		config := DefaultConfig()
		config.Storage.Timeout = 0
		err := config.Validate()
		if err != nil {
			t.Errorf("Zero timeout should be valid, got error: %v", err)
		}
	})

	t.Run("very large max nodes", func(t *testing.T) {
		config := DefaultConfig()
		config.Network.MaxNodes = 1000000
		err := config.Validate()
		if err != nil {
			t.Errorf("Large MaxNodes should be valid, got error: %v", err)
		}
	})

	t.Run("zero max nodes", func(t *testing.T) {
		config := DefaultConfig()
		config.Network.MaxNodes = 0
		err := config.Validate()
		if err != nil {
			t.Errorf("Zero MaxNodes should be valid, got error: %v", err)
		}
	})
}

func TestValidationError_Fields(t *testing.T) {
	err := &ValidationError{
		Field:   "test.field",
		Value:   "test value",
		Message: "test message",
	}

	if err.Field != "test.field" {
		t.Errorf("Field = %s, want test.field", err.Field)
	}
	if err.Value != "test value" {
		t.Errorf("Value = %v, want test value", err.Value)
	}
	if err.Message != "test message" {
		t.Errorf("Message = %s, want test message", err.Message)
	}
}

func TestConfig_AllStorageTypes(t *testing.T) {
	types := []struct {
		name  string
		value string
		valid bool
	}{
		{"memory", "memory", true},
		{"invalid type", "invalid", false},
		{"empty", "", false},
		{"uppercase MEMORY", "MEMORY", false},
		{"mixed case Memory", "Memory", false},
	}

	for _, tt := range types {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.Storage.Type = tt.value
			err := config.Validate()

			if tt.valid && err != nil {
				t.Errorf("Storage type %q should be valid, got error: %v", tt.value, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("Storage type %q should be invalid, got no error", tt.value)
			}
		})
	}
}

func TestConfig_AllLoggerLevels(t *testing.T) {
	levels := []struct {
		name  string
		value string
		valid bool
	}{
		{"debug", "debug", true},
		{"info", "info", true},
		{"warn", "warn", true},
		{"error", "error", true},
		{"trace", "trace", false},
		{"fatal", "fatal", false},
		{"empty", "", false},
		{"uppercase DEBUG", "DEBUG", false},
		{"mixed case Info", "Info", false},
	}

	for _, tt := range levels {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.Logger.Level = tt.value
			err := config.Validate()

			if tt.valid && err != nil {
				t.Errorf("Logger level %q should be valid, got error: %v", tt.value, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("Logger level %q should be invalid, got no error", tt.value)
			}
		})
	}
}
