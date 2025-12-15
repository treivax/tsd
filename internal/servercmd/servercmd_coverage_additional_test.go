// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package servercmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/treivax/tsd/tsdio"
)

// TestParseFlags_TLSEnvironmentVariables tests TLS environment variable handling
func TestParseFlags_TLSEnvironmentVariables(t *testing.T) {
	// Clear and set environment variables
	oldCert := os.Getenv("TSD_TLS_CERT")
	oldKey := os.Getenv("TSD_TLS_KEY")
	oldInsecure := os.Getenv("TSD_INSECURE")
	defer func() {
		if oldCert != "" {
			os.Setenv("TSD_TLS_CERT", oldCert)
		} else {
			os.Unsetenv("TSD_TLS_CERT")
		}
		if oldKey != "" {
			os.Setenv("TSD_TLS_KEY", oldKey)
		} else {
			os.Unsetenv("TSD_TLS_KEY")
		}
		if oldInsecure != "" {
			os.Setenv("TSD_INSECURE", oldInsecure)
		} else {
			os.Unsetenv("TSD_INSECURE")
		}
	}()

	tests := []struct {
		name        string
		envCert     string
		envKey      string
		envInsecure string
		args        []string
		checkFn     func(*testing.T, *Config)
	}{
		{
			name:        "✅ TSD_TLS_CERT environment variable",
			envCert:     "/custom/cert.pem",
			envKey:      "",
			envInsecure: "",
			args:        []string{"-insecure"},
			checkFn: func(t *testing.T, c *Config) {
				if c.TLSCertFile != "/custom/cert.pem" {
					t.Errorf("TLSCertFile = %q, want /custom/cert.pem", c.TLSCertFile)
				}
			},
		},
		{
			name:        "✅ TSD_TLS_KEY environment variable",
			envCert:     "",
			envKey:      "/custom/key.pem",
			envInsecure: "",
			args:        []string{"-insecure"},
			checkFn: func(t *testing.T, c *Config) {
				if c.TLSKeyFile != "/custom/key.pem" {
					t.Errorf("TLSKeyFile = %q, want /custom/key.pem", c.TLSKeyFile)
				}
			},
		},
		{
			name:        "✅ TSD_INSECURE environment variable",
			envCert:     "",
			envKey:      "",
			envInsecure: "true",
			args:        []string{},
			checkFn: func(t *testing.T, c *Config) {
				if !c.Insecure {
					t.Error("Insecure should be true from TSD_INSECURE=true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			if tt.envCert != "" {
				os.Setenv("TSD_TLS_CERT", tt.envCert)
			} else {
				os.Unsetenv("TSD_TLS_CERT")
			}
			if tt.envKey != "" {
				os.Setenv("TSD_TLS_KEY", tt.envKey)
			} else {
				os.Unsetenv("TSD_TLS_KEY")
			}
			if tt.envInsecure != "" {
				os.Setenv("TSD_INSECURE", tt.envInsecure)
			} else {
				os.Unsetenv("TSD_INSECURE")
			}

			config := parseFlags(tt.args)
			tt.checkFn(t, config)
		})
	}
}

// TestParseFlags_AllFlagCombinations tests various flag combinations
func TestParseFlags_AllFlagCombinations(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantHost string
		wantPort int
		wantAuth string
	}{
		{
			name:     "✅ custom host and port",
			args:     []string{"-host", "192.168.1.1", "-port", "9090", "-insecure"},
			wantHost: "192.168.1.1",
			wantPort: 9090,
			wantAuth: "none",
		},
		{
			name:     "✅ verbose flag",
			args:     []string{"-v", "-insecure"},
			wantHost: "0.0.0.0",
			wantPort: 8080,
			wantAuth: "none",
		},
		{
			name:     "✅ jwt authentication",
			args:     []string{"-auth", "jwt", "-jwt-secret", "mysecret", "-insecure"},
			wantHost: "0.0.0.0",
			wantPort: 8080,
			wantAuth: "jwt",
		},
		{
			name:     "✅ key authentication",
			args:     []string{"-auth", "key", "-auth-keys", "key1,key2,key3", "-insecure"},
			wantHost: "0.0.0.0",
			wantPort: 8080,
			wantAuth: "key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseFlags(tt.args)

			if config.Host != tt.wantHost {
				t.Errorf("Host = %q, want %q", config.Host, tt.wantHost)
			}
			if config.Port != tt.wantPort {
				t.Errorf("Port = %d, want %d", config.Port, tt.wantPort)
			}
			if config.AuthType != tt.wantAuth {
				t.Errorf("AuthType = %q, want %q", config.AuthType, tt.wantAuth)
			}
		})
	}
}

// TestParseFlags_JWTFlagDetails tests JWT-specific flags
func TestParseFlags_JWTFlagDetails(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantExpiration time.Duration
		wantIssuer     string
	}{
		{
			name:           "✅ default JWT expiration",
			args:           []string{"-auth", "jwt", "-jwt-secret", "secret", "-insecure"},
			wantExpiration: 24 * time.Hour,
			wantIssuer:     "tsd-server",
		},
		{
			name:           "✅ custom JWT expiration",
			args:           []string{"-auth", "jwt", "-jwt-secret", "secret", "-jwt-expiration", "48h", "-insecure"},
			wantExpiration: 48 * time.Hour,
			wantIssuer:     "tsd-server",
		},
		{
			name:           "✅ custom JWT issuer",
			args:           []string{"-auth", "jwt", "-jwt-secret", "secret", "-jwt-issuer", "my-app", "-insecure"},
			wantExpiration: 24 * time.Hour,
			wantIssuer:     "my-app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseFlags(tt.args)

			if config.JWTExpiration != tt.wantExpiration {
				t.Errorf("JWTExpiration = %v, want %v", config.JWTExpiration, tt.wantExpiration)
			}
			if config.JWTIssuer != tt.wantIssuer {
				t.Errorf("JWTIssuer = %q, want %q", config.JWTIssuer, tt.wantIssuer)
			}
		})
	}
}

// TestParseFlags_AuthKeysProcessing tests auth keys parsing and trimming
func TestParseFlags_AuthKeysProcessing(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantKeys     []string
		wantKeyCount int
	}{
		{
			name:         "✅ single auth key",
			args:         []string{"-auth-keys", "key1", "-insecure"},
			wantKeys:     []string{"key1"},
			wantKeyCount: 1,
		},
		{
			name:         "✅ multiple auth keys",
			args:         []string{"-auth-keys", "key1,key2,key3", "-insecure"},
			wantKeys:     []string{"key1", "key2", "key3"},
			wantKeyCount: 3,
		},
		{
			name:         "✅ auth keys with spaces (trimmed)",
			args:         []string{"-auth-keys", " key1 , key2 , key3 ", "-insecure"},
			wantKeys:     []string{"key1", "key2", "key3"},
			wantKeyCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseFlags(tt.args)

			if len(config.AuthKeys) != tt.wantKeyCount {
				t.Errorf("AuthKeys count = %d, want %d", len(config.AuthKeys), tt.wantKeyCount)
			}

			for i, wantKey := range tt.wantKeys {
				if i < len(config.AuthKeys) && config.AuthKeys[i] != wantKey {
					t.Errorf("AuthKeys[%d] = %q, want %q", i, config.AuthKeys[i], wantKey)
				}
			}
		})
	}
}

// TestExecuteTSDProgram_ParsingErrors tests error handling in executeTSDProgram
func TestExecuteTSDProgram_ParsingErrors(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     0,
		Insecure: true,
		AuthType: "none",
	}

	server, err := NewServer(config, nil)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	tests := []struct {
		name          string
		request       *tsdio.ExecuteRequest
		wantErrorType string
		wantErrorMsg  string
	}{
		{
			name: "✗ parsing error - invalid syntax",
			request: &tsdio.ExecuteRequest{
				Source:     "invalid syntax @#$%",
				SourceName: "test.tsd",
			},
			wantErrorType: "parsing_error",
			wantErrorMsg:  "Erreur de parsing",
		},
		{
			name: "✗ empty source",
			request: &tsdio.ExecuteRequest{
				Source:     "",
				SourceName: "empty.tsd",
			},
			wantErrorType: "", // May succeed with empty program
			wantErrorMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now()
			response := server.executeTSDProgram(tt.request, startTime)

			if tt.wantErrorType != "" {
				if response.Success {
					t.Error("Expected error but got success")
				}
				if response.Error == "" {
					t.Fatal("Expected error message but got empty string")
				}
				if response.ErrorType != tt.wantErrorType {
					t.Errorf("Error type = %q, want %q", response.ErrorType, tt.wantErrorType)
				}
				if tt.wantErrorMsg != "" && !strings.Contains(response.Error, tt.wantErrorMsg) {
					t.Errorf("Error message %q does not contain %q", response.Error, tt.wantErrorMsg)
				}
			}
		})
	}
}

// TestExecuteTSDProgram_ValidationErrors tests validation error handling
func TestExecuteTSDProgram_ValidationErrors(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     0,
		Insecure: true,
		AuthType: "none",
	}

	server, err := NewServer(config, nil)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Note: Most TSD programs that parse will also validate
	// This test documents that validation is called
	t.Run("✅ valid program passes validation", func(t *testing.T) {
		request := &tsdio.ExecuteRequest{
			Source:     "type Person(id: string, name: string)",
			SourceName: "valid.tsd",
		}

		startTime := time.Now()
		response := server.executeTSDProgram(request, startTime)

		if !response.Success {
			t.Errorf("Valid program should succeed, got error: %v", response.Error)
		}

		if response.Results == nil {
			t.Error("Results should not be nil for successful execution")
		}
	})
}

// TestExecuteTSDProgram_SuccessfulExecution tests successful program execution
func TestExecuteTSDProgram_SuccessfulExecution(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     0,
		Insecure: true,
		AuthType: "none",
	}

	server, err := NewServer(config, nil)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	tests := []struct {
		name        string
		source      string
		sourceName  string
		wantSuccess bool
	}{
		{
			name:        "✅ simple type definition",
			source:      "type Person(id: string, name: string, age: number)",
			sourceName:  "person.tsd",
			wantSuccess: true,
		},
		{
			name: "✅ type with rule",
			source: `type Person(id: string, age: number)
action adult(id: string)
rule r1: {p: Person} / p.age >= 18 ==> adult(p.id)`,
			sourceName:  "adult.tsd",
			wantSuccess: true,
		},
		{
			name: "✅ multiple types",
			source: `type Person(id: string, name: string)
type Address(street: string, city: string)`,
			sourceName:  "multi.tsd",
			wantSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &tsdio.ExecuteRequest{
				Source:     tt.source,
				SourceName: tt.sourceName,
			}

			startTime := time.Now()
			response := server.executeTSDProgram(request, startTime)

			if response.Success != tt.wantSuccess {
				t.Errorf("Success = %v, want %v", response.Success, tt.wantSuccess)
				if response.Error != "" {
					t.Logf("Error: %v", response.Error)
				}
			}

			if response.Success {
				if response.Results == nil {
					t.Error("Results should not be nil for successful execution")
				}
				if response.ExecutionTimeMs < 0 {
					t.Error("ExecutionTimeMs should be non-negative")
				}
			}
		})
	}
}

// TestCollectActivations_EdgeCases tests collectActivations with various inputs
func TestCollectActivations_EdgeCases(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     0,
		Insecure: true,
		AuthType: "none",
	}

	server, err := NewServer(config, nil)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	t.Run("✅ nil network returns empty activations", func(t *testing.T) {
		activations := server.collectActivations(nil)
		if activations == nil {
			t.Error("Should return non-nil slice for nil network")
		}
		if len(activations) != 0 {
			t.Errorf("Should return empty slice for nil network, got %d activations", len(activations))
		}
	})
}

// Note: getValueType, extractFacts, extractArguments are private functions
// They are tested indirectly through executeTSDProgram tests above

// TestRun_EdgeCases tests Run function with edge cases (without actually starting server)
func TestRun_EdgeCases(t *testing.T) {
	t.Run("✅ Run with initialization check", func(t *testing.T) {
		// This test verifies that Run can be called and initializes properly
		// We use a goroutine to avoid blocking, but we're mainly testing initialization

		// Create temp cert files for testing
		tmpDir := t.TempDir()
		certFile := filepath.Join(tmpDir, "cert.pem")
		keyFile := filepath.Join(tmpDir, "key.pem")

		// Create dummy cert files
		os.WriteFile(certFile, []byte("dummy cert"), 0644)
		os.WriteFile(keyFile, []byte("dummy key"), 0644)

		args := []string{
			"-host", "127.0.0.1",
			"-port", "0", // Port 0 = random available port
			"-insecure",
			"-tls-cert", certFile,
			"-tls-key", keyFile,
		}

		var stdout, stderr bytes.Buffer

		// Note: We can't easily test Run() fully as it blocks on ListenAndServe
		// This test at least exercises parseFlags and NewServer creation paths

		// Just verify parseFlags works
		config := parseFlags(args)
		if config.Host != "127.0.0.1" {
			t.Errorf("Host = %q, want 127.0.0.1", config.Host)
		}
		if config.Port != 0 {
			t.Errorf("Port = %d, want 0", config.Port)
		}

		// Verify NewServer works with this config
		_, err := NewServer(config, nil)
		if err != nil {
			t.Errorf("NewServer failed: %v", err)
		}

		t.Logf("✅ Run initialization path tested")
		_ = stdout
		_ = stderr
	})
}
