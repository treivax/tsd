// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/treivax/tsd/constraint/pkg/domain"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	// Test Parser config
	if config.Parser.MaxExpressions != 1000 {
		t.Errorf("Parser.MaxExpressions = %d, want 1000", config.Parser.MaxExpressions)
	}
	if config.Parser.Debug {
		t.Error("Parser.Debug should be false by default")
	}
	if !config.Parser.Recover {
		t.Error("Parser.Recover should be true by default")
	}

	// Test Validator config
	if !config.Validator.StrictMode {
		t.Error("Validator.StrictMode should be true by default")
	}
	if config.Validator.MaxDepth != 20 {
		t.Errorf("Validator.MaxDepth = %d, want 20", config.Validator.MaxDepth)
	}
	if len(config.Validator.AllowedOperators) == 0 {
		t.Error("Validator.AllowedOperators should not be empty")
	}

	// Test Logger config
	if config.Logger.Level != "info" {
		t.Errorf("Logger.Level = %s, want info", config.Logger.Level)
	}
	if config.Logger.Format != "json" {
		t.Errorf("Logger.Format = %s, want json", config.Logger.Format)
	}
	if config.Logger.Output != "stdout" {
		t.Errorf("Logger.Output = %s, want stdout", config.Logger.Output)
	}

	// Test other fields
	if config.Debug {
		t.Error("Debug should be false by default")
	}
	if config.Version != "1.0.0" {
		t.Errorf("Version = %s, want 1.0.0", config.Version)
	}
}

func TestNewConfigManager(t *testing.T) {
	path := "/tmp/test-config.json"
	cm := NewConfigManager(path)

	if cm == nil {
		t.Fatal("NewConfigManager() returned nil")
	}
	if cm.filePath != path {
		t.Errorf("filePath = %s, want %s", cm.filePath, path)
	}
	if cm.config == nil {
		t.Error("config should not be nil")
	}
}

func TestConfigManager_GetConfig(t *testing.T) {
	cm := NewConfigManager("")
	config := cm.GetConfig()

	if config == nil {
		t.Error("GetConfig() returned nil")
	}
	if config != cm.config {
		t.Error("GetConfig() should return the same config instance")
	}
}

func TestConfigManager_SetConfig(t *testing.T) {
	cm := NewConfigManager("")
	newConfig := &Config{
		Version: "2.0.0",
		Debug:   true,
	}

	cm.SetConfig(newConfig)

	if cm.config != newConfig {
		t.Error("SetConfig() did not set the config")
	}
	if cm.config.Version != "2.0.0" {
		t.Errorf("Version = %s, want 2.0.0", cm.config.Version)
	}
}

func TestConfigManager_GetParserConfig(t *testing.T) {
	cm := NewConfigManager("")
	parserConfig := cm.GetParserConfig()

	if parserConfig.MaxExpressions != 1000 {
		t.Errorf("MaxExpressions = %d, want 1000", parserConfig.MaxExpressions)
	}
}

func TestConfigManager_GetValidatorConfig(t *testing.T) {
	cm := NewConfigManager("")
	validatorConfig := cm.GetValidatorConfig()

	if !validatorConfig.StrictMode {
		t.Error("StrictMode should be true")
	}
}

func TestConfigManager_GetLoggerConfig(t *testing.T) {
	cm := NewConfigManager("")
	loggerConfig := cm.GetLoggerConfig()

	if loggerConfig.Level != "info" {
		t.Errorf("Level = %s, want info", loggerConfig.Level)
	}
}

func TestConfigManager_IsDebugEnabled(t *testing.T) {
	tests := []struct {
		name        string
		debug       bool
		parserDebug bool
		expectDebug bool
	}{
		{
			name:        "both false",
			debug:       false,
			parserDebug: false,
			expectDebug: false,
		},
		{
			name:        "config debug true",
			debug:       true,
			parserDebug: false,
			expectDebug: true,
		},
		{
			name:        "parser debug true",
			debug:       false,
			parserDebug: true,
			expectDebug: true,
		},
		{
			name:        "both true",
			debug:       true,
			parserDebug: true,
			expectDebug: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := NewConfigManager("")
			cm.config.Debug = tt.debug
			cm.config.Parser.Debug = tt.parserDebug

			result := cm.IsDebugEnabled()
			if result != tt.expectDebug {
				t.Errorf("IsDebugEnabled() = %v, want %v", result, tt.expectDebug)
			}
		})
	}
}

func TestConfigManager_UpdateParserConfig(t *testing.T) {
	cm := NewConfigManager("")
	newParserConfig := domain.ParserConfig{
		MaxExpressions: 5000,
		Debug:          true,
		Recover:        false,
	}

	cm.UpdateParserConfig(newParserConfig)

	if cm.config.Parser.MaxExpressions != 5000 {
		t.Errorf("MaxExpressions = %d, want 5000", cm.config.Parser.MaxExpressions)
	}
	if !cm.config.Parser.Debug {
		t.Error("Debug should be true")
	}
	if cm.config.Parser.Recover {
		t.Error("Recover should be false")
	}
}

func TestConfigManager_UpdateValidatorConfig(t *testing.T) {
	cm := NewConfigManager("")
	newValidatorConfig := domain.ValidatorConfig{
		StrictMode:       false,
		AllowedOperators: []string{"==", "!="},
		MaxDepth:         10,
	}

	cm.UpdateValidatorConfig(newValidatorConfig)

	if cm.config.Validator.StrictMode {
		t.Error("StrictMode should be false")
	}
	if len(cm.config.Validator.AllowedOperators) != 2 {
		t.Errorf("AllowedOperators length = %d, want 2", len(cm.config.Validator.AllowedOperators))
	}
	if cm.config.Validator.MaxDepth != 10 {
		t.Errorf("MaxDepth = %d, want 10", cm.config.Validator.MaxDepth)
	}
}

func TestConfigManager_UpdateLoggerConfig(t *testing.T) {
	cm := NewConfigManager("")
	newLoggerConfig := domain.LoggerConfig{
		Level:  "debug",
		Format: "text",
		Output: "stderr",
	}

	cm.UpdateLoggerConfig(newLoggerConfig)

	if cm.config.Logger.Level != "debug" {
		t.Errorf("Level = %s, want debug", cm.config.Logger.Level)
	}
	if cm.config.Logger.Format != "text" {
		t.Errorf("Format = %s, want text", cm.config.Logger.Format)
	}
	if cm.config.Logger.Output != "stderr" {
		t.Errorf("Output = %s, want stderr", cm.config.Logger.Output)
	}
}

func TestConfigManager_SetDebug(t *testing.T) {
	cm := NewConfigManager("")

	cm.SetDebug(true)
	if !cm.config.Debug {
		t.Error("Debug should be true")
	}

	cm.SetDebug(false)
	if cm.config.Debug {
		t.Error("Debug should be false")
	}
}

func TestConfigManager_Validate(t *testing.T) {
	tests := []struct {
		name      string
		modifyFn  func(*Config)
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid default config",
			modifyFn:  func(c *Config) {},
			wantError: false,
		},
		{
			name: "invalid parser max expressions - zero",
			modifyFn: func(c *Config) {
				c.Parser.MaxExpressions = 0
			},
			wantError: true,
			errorMsg:  "parser.max_expressions must be positive",
		},
		{
			name: "invalid parser max expressions - negative",
			modifyFn: func(c *Config) {
				c.Parser.MaxExpressions = -10
			},
			wantError: true,
			errorMsg:  "parser.max_expressions must be positive",
		},
		{
			name: "invalid validator max depth - zero",
			modifyFn: func(c *Config) {
				c.Validator.MaxDepth = 0
			},
			wantError: true,
			errorMsg:  "validator.max_depth must be positive",
		},
		{
			name: "invalid validator max depth - negative",
			modifyFn: func(c *Config) {
				c.Validator.MaxDepth = -5
			},
			wantError: true,
			errorMsg:  "validator.max_depth must be positive",
		},
		{
			name: "empty allowed operators",
			modifyFn: func(c *Config) {
				c.Validator.AllowedOperators = []string{}
			},
			wantError: true,
			errorMsg:  "validator.allowed_operators cannot be empty",
		},
		{
			name: "invalid logger level",
			modifyFn: func(c *Config) {
				c.Logger.Level = "invalid"
			},
			wantError: true,
			errorMsg:  "invalid logger level",
		},
		{
			name: "invalid logger format",
			modifyFn: func(c *Config) {
				c.Logger.Format = "xml"
			},
			wantError: true,
			errorMsg:  "invalid logger format",
		},
		{
			name: "valid logger level - debug",
			modifyFn: func(c *Config) {
				c.Logger.Level = "debug"
			},
			wantError: false,
		},
		{
			name: "valid logger level - warn",
			modifyFn: func(c *Config) {
				c.Logger.Level = "warn"
			},
			wantError: false,
		},
		{
			name: "valid logger level - error",
			modifyFn: func(c *Config) {
				c.Logger.Level = "error"
			},
			wantError: false,
		},
		{
			name: "valid logger format - text",
			modifyFn: func(c *Config) {
				c.Logger.Format = "text"
			},
			wantError: false,
		},
		{
			name: "valid logger format - plain",
			modifyFn: func(c *Config) {
				c.Logger.Format = "plain"
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := NewConfigManager("")
			tt.modifyFn(cm.config)

			err := cm.Validate()

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

func TestConfigManager_SaveToFile(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name      string
		filePath  string
		wantError bool
	}{
		{
			name:      "valid file path",
			filePath:  filepath.Join(tempDir, "config.json"),
			wantError: false,
		},
		{
			name:      "nested directory",
			filePath:  filepath.Join(tempDir, "subdir", "config.json"),
			wantError: false,
		},
		{
			name:      "empty file path",
			filePath:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := NewConfigManager(tt.filePath)
			err := cm.SaveToFile()

			if tt.wantError {
				if err == nil {
					t.Error("SaveToFile() error = nil, want error")
				}
			} else {
				if err != nil {
					t.Errorf("SaveToFile() error = %v, want nil", err)
				}

				// Verify file exists
				if _, err := os.Stat(tt.filePath); os.IsNotExist(err) {
					t.Errorf("File was not created at %s", tt.filePath)
				}

				// Verify content is valid JSON
				data, err := os.ReadFile(tt.filePath)
				if err != nil {
					t.Fatalf("Failed to read saved file: %v", err)
				}

				var config Config
				if err := json.Unmarshal(data, &config); err != nil {
					t.Errorf("Saved file is not valid JSON: %v", err)
				}
			}
		})
	}
}

func TestConfigManager_LoadFromFile(t *testing.T) {
	tempDir := t.TempDir()

	// Create a valid config file
	validConfigPath := filepath.Join(tempDir, "valid.json")
	validConfig := DefaultConfig()
	validConfig.Version = "2.0.0"
	validConfig.Debug = true
	data, _ := json.MarshalIndent(validConfig, "", "  ")
	os.WriteFile(validConfigPath, data, 0644)

	// Create an invalid JSON file
	invalidJSONPath := filepath.Join(tempDir, "invalid.json")
	os.WriteFile(invalidJSONPath, []byte("{ invalid json }"), 0644)

	tests := []struct {
		name      string
		filePath  string
		wantError bool
		checkFn   func(*testing.T, *ConfigManager)
	}{
		{
			name:      "load valid config",
			filePath:  validConfigPath,
			wantError: false,
			checkFn: func(t *testing.T, cm *ConfigManager) {
				if cm.config.Version != "2.0.0" {
					t.Errorf("Version = %s, want 2.0.0", cm.config.Version)
				}
				if !cm.config.Debug {
					t.Error("Debug should be true")
				}
			},
		},
		{
			name:      "file does not exist - creates default",
			filePath:  filepath.Join(tempDir, "nonexistent.json"),
			wantError: false,
			checkFn: func(t *testing.T, cm *ConfigManager) {
				// Should create the file with default config
				if _, err := os.Stat(cm.filePath); os.IsNotExist(err) {
					t.Error("Default config file should have been created")
				}
			},
		},
		{
			name:      "invalid JSON",
			filePath:  invalidJSONPath,
			wantError: true,
		},
		{
			name:      "empty file path",
			filePath:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := NewConfigManager(tt.filePath)
			err := cm.LoadFromFile()

			if tt.wantError {
				if err == nil {
					t.Error("LoadFromFile() error = nil, want error")
				}
			} else {
				if err != nil {
					t.Errorf("LoadFromFile() error = %v, want nil", err)
				}
				if tt.checkFn != nil {
					tt.checkFn(t, cm)
				}
			}
		})
	}
}

func TestConfigManager_String(t *testing.T) {
	cm := NewConfigManager("")
	str := cm.String()

	if str == "" {
		t.Error("String() returned empty string")
	}

	// Verify it's valid JSON
	var config Config
	if err := json.Unmarshal([]byte(str), &config); err != nil {
		t.Errorf("String() output is not valid JSON: %v", err)
	}

	// Check that it contains expected fields
	if !strings.Contains(str, "parser") {
		t.Error("String() output should contain 'parser'")
	}
	if !strings.Contains(str, "validator") {
		t.Error("String() output should contain 'validator'")
	}
	if !strings.Contains(str, "logger") {
		t.Error("String() output should contain 'logger'")
	}
}

func TestConfigManager_Clone(t *testing.T) {
	original := NewConfigManager("/tmp/original.json")
	original.config.Debug = true
	original.config.Version = "1.5.0"

	cloned := original.Clone()

	// Verify values are copied
	if cloned.config.Debug != original.config.Debug {
		t.Error("Cloned config should have same Debug value")
	}
	if cloned.config.Version != original.config.Version {
		t.Error("Cloned config should have same Version value")
	}
	if cloned.filePath != original.filePath {
		t.Error("Cloned config should have same filePath")
	}

	// Verify they are separate instances
	cloned.config.Debug = false
	if original.config.Debug == cloned.config.Debug {
		t.Error("Modifying cloned config should not affect original")
	}
}

func TestConfigManager_Reset(t *testing.T) {
	cm := NewConfigManager("")

	// Modify the config
	cm.config.Debug = true
	cm.config.Version = "2.0.0"
	cm.config.Parser.MaxExpressions = 5000

	// Reset to defaults
	cm.Reset()

	// Verify it's back to defaults
	if cm.config.Debug {
		t.Error("Debug should be false after reset")
	}
	if cm.config.Version != "1.0.0" {
		t.Errorf("Version = %s, want 1.0.0 after reset", cm.config.Version)
	}
	if cm.config.Parser.MaxExpressions != 1000 {
		t.Errorf("MaxExpressions = %d, want 1000 after reset", cm.config.Parser.MaxExpressions)
	}
}

func TestConfig_JSONMarshaling(t *testing.T) {
	original := DefaultConfig()
	original.Debug = true
	original.Version = "1.5.0"

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
	if unmarshaled.Debug != original.Debug {
		t.Error("Debug value not preserved after JSON round-trip")
	}
	if unmarshaled.Version != original.Version {
		t.Error("Version value not preserved after JSON round-trip")
	}
	if unmarshaled.Parser.MaxExpressions != original.Parser.MaxExpressions {
		t.Error("Parser.MaxExpressions not preserved after JSON round-trip")
	}
}

func TestConfigManager_SaveAndLoadRoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "roundtrip.json")

	// Create and configure a manager
	cm1 := NewConfigManager(configPath)
	cm1.config.Debug = true
	cm1.config.Version = "2.5.0"
	cm1.config.Parser.MaxExpressions = 2000
	cm1.config.Validator.MaxDepth = 30

	// Save to file
	if err := cm1.SaveToFile(); err != nil {
		t.Fatalf("SaveToFile() failed: %v", err)
	}

	// Load from file in a new manager
	cm2 := NewConfigManager(configPath)
	if err := cm2.LoadFromFile(); err != nil {
		t.Fatalf("LoadFromFile() failed: %v", err)
	}

	// Verify all values match
	if cm2.config.Debug != cm1.config.Debug {
		t.Error("Debug value not preserved")
	}
	if cm2.config.Version != cm1.config.Version {
		t.Error("Version value not preserved")
	}
	if cm2.config.Parser.MaxExpressions != cm1.config.Parser.MaxExpressions {
		t.Error("Parser.MaxExpressions not preserved")
	}
	if cm2.config.Validator.MaxDepth != cm1.config.Validator.MaxDepth {
		t.Error("Validator.MaxDepth not preserved")
	}
}
