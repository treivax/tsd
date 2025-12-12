// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package config

import (
	"os"
	"testing"

	"github.com/treivax/tsd/constraint/pkg/domain"
)

func TestLoadFromEnv(t *testing.T) {
	tests := []struct {
		name      string
		envVars   map[string]string
		checkFn   func(*testing.T, *ConfigManager)
		wantError bool
	}{
		{
			name: "load max expressions from env",
			envVars: map[string]string{
				EnvMaxExpressions: "5000",
			},
			checkFn: func(t *testing.T, cm *ConfigManager) {
				if cm.config.Parser.MaxExpressions != 5000 {
					t.Errorf("MaxExpressions = %d, want 5000", cm.config.Parser.MaxExpressions)
				}
			},
			wantError: false,
		},
		{
			name: "load max depth from env",
			envVars: map[string]string{
				EnvMaxDepth: "30",
			},
			checkFn: func(t *testing.T, cm *ConfigManager) {
				if cm.config.Validator.MaxDepth != 30 {
					t.Errorf("MaxDepth = %d, want 30", cm.config.Validator.MaxDepth)
				}
			},
			wantError: false,
		},
		{
			name: "load debug from env",
			envVars: map[string]string{
				EnvDebug: "true",
			},
			checkFn: func(t *testing.T, cm *ConfigManager) {
				if !cm.config.Debug {
					t.Error("Debug should be true")
				}
			},
			wantError: false,
		},
		{
			name: "load strict mode from env",
			envVars: map[string]string{
				EnvStrictMode: "false",
			},
			checkFn: func(t *testing.T, cm *ConfigManager) {
				if cm.config.Validator.StrictMode {
					t.Error("StrictMode should be false")
				}
			},
			wantError: false,
		},
		{
			name: "load log level from env",
			envVars: map[string]string{
				EnvLogLevel: "debug",
			},
			checkFn: func(t *testing.T, cm *ConfigManager) {
				if cm.config.Logger.Level != "debug" {
					t.Errorf("LogLevel = %s, want debug", cm.config.Logger.Level)
				}
			},
			wantError: false,
		},
		{
			name: "load log format from env",
			envVars: map[string]string{
				EnvLogFormat: "text",
			},
			checkFn: func(t *testing.T, cm *ConfigManager) {
				if cm.config.Logger.Format != "text" {
					t.Errorf("LogFormat = %s, want text", cm.config.Logger.Format)
				}
			},
			wantError: false,
		},
		{
			name: "load log output from env",
			envVars: map[string]string{
				EnvLogOutput: "stderr",
			},
			checkFn: func(t *testing.T, cm *ConfigManager) {
				if cm.config.Logger.Output != "stderr" {
					t.Errorf("LogOutput = %s, want stderr", cm.config.Logger.Output)
				}
			},
			wantError: false,
		},
		{
			name: "invalid max expressions",
			envVars: map[string]string{
				EnvMaxExpressions: "invalid",
			},
			wantError: true,
		},
		{
			name: "invalid max depth",
			envVars: map[string]string{
				EnvMaxDepth: "not-a-number",
			},
			wantError: true,
		},
		{
			name: "invalid debug value",
			envVars: map[string]string{
				EnvDebug: "maybe",
			},
			wantError: true,
		},
		{
			name: "multiple env vars",
			envVars: map[string]string{
				EnvMaxExpressions: "2000",
				EnvDebug:          "true",
				EnvLogLevel:       "warn",
			},
			checkFn: func(t *testing.T, cm *ConfigManager) {
				if cm.config.Parser.MaxExpressions != 2000 {
					t.Errorf("MaxExpressions = %d, want 2000", cm.config.Parser.MaxExpressions)
				}
				if !cm.config.Debug {
					t.Error("Debug should be true")
				}
				if cm.config.Logger.Level != "warn" {
					t.Errorf("LogLevel = %s, want warn", cm.config.Logger.Level)
				}
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer func() {
				// Clean up
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			cm := NewConfigManager("")
			err := cm.LoadFromEnv()

			if tt.wantError {
				if err == nil {
					t.Error("LoadFromEnv() error = nil, want error")
				}
			} else {
				if err != nil {
					t.Errorf("LoadFromEnv() error = %v, want nil", err)
				}
				if tt.checkFn != nil {
					tt.checkFn(t, cm)
				}
			}
		})
	}
}

func TestGetConfigFilePath(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		defaultPath string
		want        string
	}{
		{
			name:        "no env - use default",
			envValue:    "",
			defaultPath: "/default/config.json",
			want:        "/default/config.json",
		},
		{
			name:        "env set - use env",
			envValue:    "/custom/config.json",
			defaultPath: "/default/config.json",
			want:        "/custom/config.json",
		},
		{
			name:        "empty default",
			envValue:    "",
			defaultPath: "",
			want:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(EnvConfigFile, tt.envValue)
				defer os.Unsetenv(EnvConfigFile)
			} else {
				os.Unsetenv(EnvConfigFile)
			}

			got := GetConfigFilePath(tt.defaultPath)
			if got != tt.want {
				t.Errorf("GetConfigFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeConfig(t *testing.T) {
	tests := []struct {
		name     string
		base     *Config
		source   *Config
		checkFn  func(*testing.T, *Config)
	}{
		{
			name:   "merge nil source",
			base:   DefaultConfig(),
			source: nil,
			checkFn: func(t *testing.T, c *Config) {
				// Should remain unchanged
				if c.Parser.MaxExpressions != DefaultMaxExpressions {
					t.Error("Config should not change when merging nil")
				}
			},
		},
		{
			name: "merge max expressions",
			base: DefaultConfig(),
			source: &Config{
				Parser: domain.ParserConfig{
					MaxExpressions: 5000,
				},
			},
			checkFn: func(t *testing.T, c *Config) {
				if c.Parser.MaxExpressions != 5000 {
					t.Errorf("MaxExpressions = %d, want 5000", c.Parser.MaxExpressions)
				}
			},
		},
		{
			name: "merge allowed operators",
			base: DefaultConfig(),
			source: &Config{
				Validator: domain.ValidatorConfig{
					AllowedOperators: []string{"==", "!="},
				},
			},
			checkFn: func(t *testing.T, c *Config) {
				if len(c.Validator.AllowedOperators) != 2 {
					t.Errorf("AllowedOperators length = %d, want 2", len(c.Validator.AllowedOperators))
				}
			},
		},
		{
			name: "merge log level",
			base: DefaultConfig(),
			source: &Config{
				Logger: domain.LoggerConfig{
					Level: "debug",
				},
			},
			checkFn: func(t *testing.T, c *Config) {
				if c.Logger.Level != "debug" {
					t.Errorf("LogLevel = %s, want debug", c.Logger.Level)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := NewConfigManager("")
			cm.config = tt.base
			cm.MergeConfig(tt.source)
			if tt.checkFn != nil {
				tt.checkFn(t, cm.config)
			}
		})
	}
}

func TestCloneDeepCopy(t *testing.T) {
	// Verify Clone creates a deep copy
	original := NewConfigManager("")
	original.config.Validator.AllowedOperators = []string{"==", "!="}

	cloned := original.Clone()

	// Modify cloned config's slice
	cloned.config.Validator.AllowedOperators[0] = ">"

	// Original should not be affected
	if original.config.Validator.AllowedOperators[0] != "==" {
		t.Error("Clone should be a deep copy - original was modified")
	}
}
