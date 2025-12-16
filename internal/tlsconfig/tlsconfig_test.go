// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package tlsconfig

import (
	"crypto/tls"
	"os"
	"path/filepath"
	"testing"
)

// TestDefaultCipherSuites tests the DefaultCipherSuites function
func TestDefaultCipherSuites(t *testing.T) {
	t.Log("🧪 TEST: DefaultCipherSuites")
	t.Log("==================================")

	suites := DefaultCipherSuites()

	if len(suites) == 0 {
		t.Fatal("❌ DefaultCipherSuites returned empty slice")
	}

	expectedSuites := []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	}

	if len(suites) != len(expectedSuites) {
		t.Errorf("❌ Expected %d cipher suites, got %d", len(expectedSuites), len(suites))
	}

	for i, expected := range expectedSuites {
		if i < len(suites) && suites[i] != expected {
			t.Errorf("❌ Suite[%d] = %d, want %d", i, suites[i], expected)
		}
	}

	t.Log("✅ DefaultCipherSuites test passed")
}

// TestNewServerTLSConfig tests the NewServerTLSConfig function
func TestNewServerTLSConfig(t *testing.T) {
	t.Log("🧪 TEST: NewServerTLSConfig")
	t.Log("==================================")

	tmpDir := t.TempDir()
	certFile := filepath.Join(tmpDir, "server.crt")
	keyFile := filepath.Join(tmpDir, "server.key")

	// Create dummy cert and key files
	if err := os.WriteFile(certFile, []byte("dummy cert"), 0644); err != nil {
		t.Fatalf("❌ Failed to create cert: %v", err)
	}
	if err := os.WriteFile(keyFile, []byte("dummy key"), 0644); err != nil {
		t.Fatalf("❌ Failed to create key: %v", err)
	}

	tests := []struct {
		name      string
		config    *ServerConfig
		wantErr   bool
		checkFunc func(*testing.T, *tls.Config)
	}{
		{
			name:    "✅ Nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "✅ Missing CertFile",
			config: &ServerConfig{
				KeyFile: keyFile,
			},
			wantErr: true,
		},
		{
			name: "✅ Missing KeyFile",
			config: &ServerConfig{
				CertFile: certFile,
			},
			wantErr: true,
		},
		{
			name: "✅ Non-existent CertFile",
			config: &ServerConfig{
				CertFile: "/nonexistent/cert.crt",
				KeyFile:  keyFile,
			},
			wantErr: true,
		},
		{
			name: "✅ Non-existent KeyFile",
			config: &ServerConfig{
				CertFile: certFile,
				KeyFile:  "/nonexistent/key.key",
			},
			wantErr: true,
		},
		{
			name: "✅ Valid config with defaults",
			config: &ServerConfig{
				CertFile: certFile,
				KeyFile:  keyFile,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, cfg *tls.Config) {
				if cfg.MinVersion != MinTLSVersion {
					t.Errorf("❌ MinVersion = %d, want %d", cfg.MinVersion, MinTLSVersion)
				}
				if len(cfg.CipherSuites) == 0 {
					t.Error("❌ CipherSuites should not be empty")
				}
				// PreferServerCipherSuites reflects what's set in config (false if not specified)
				// This is expected behavior - caller should use DefaultServerConfig for opinionated defaults
			},
		},
		{
			name: "✅ Valid config with custom values",
			config: &ServerConfig{
				CertFile:                 certFile,
				KeyFile:                  keyFile,
				MinVersion:               tls.VersionTLS13,
				CipherSuites:             []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384},
				PreferServerCipherSuites: false,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, cfg *tls.Config) {
				if cfg.MinVersion != tls.VersionTLS13 {
					t.Errorf("❌ MinVersion = %d, want %d", cfg.MinVersion, tls.VersionTLS13)
				}
				if len(cfg.CipherSuites) != 1 {
					t.Errorf("❌ CipherSuites length = %d, want 1", len(cfg.CipherSuites))
				}
				if cfg.PreferServerCipherSuites {
					t.Error("❌ PreferServerCipherSuites should be false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewServerTLSConfig(tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("❌ NewServerTLSConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkFunc != nil {
				tt.checkFunc(t, cfg)
			}

			t.Log("✅ Test passed")
		})
	}
}

// TestNewClientTLSConfig tests the NewClientTLSConfig function
func TestNewClientTLSConfig(t *testing.T) {
	t.Log("🧪 TEST: NewClientTLSConfig")
	t.Log("==================================")

	tmpDir := t.TempDir()
	caFile := filepath.Join(tmpDir, "ca.crt")

	// Create dummy CA file
	caPEM := []byte(`-----BEGIN CERTIFICATE-----
MIIBhTCCASugAwIBAgIQIRi6zePL6mKjOipn+dNuaTAKBggqhkjOPQQDAjASMRAw
DgYDVQQKEwdBY21lIENvMB4XDTE3MTAyMDE5NDMwNloXDTE4MTAyMDE5NDMwNlow
EjEQMA4GA1UEChMHQWNtZSBDbzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABD0d
7VNhbWvZLWPuj/RtHFjvtJBEwOkhbN/BnnE8rnZR8+sbwnc/KhCk3FhnpHZnQz7B
5aETbbIgmuvewdjvSBSjYzBhMA4GA1UdDwEB/wQEAwICpDATBgNVHSUEDDAKBggr
BgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1UdEQQiMCCCDmxvY2FsaG9zdDo1
NDUzgg4xMjcuMC4wLjE6NTQ1MzAKBggqhkjOPQQDAgNIADBFAiEA2zpJEPQyz6/l
Wf86aX6PepsntZv2GYlA5UpabfT2EZICICpJ5h/iI+i341gBmLiAFQOyTDT+/wQc
6MF9+Yw1Yy0t
-----END CERTIFICATE-----`)

	if err := os.WriteFile(caFile, caPEM, 0644); err != nil {
		t.Fatalf("❌ Failed to create CA file: %v", err)
	}

	tests := []struct {
		name      string
		config    *ClientConfig
		wantErr   bool
		checkFunc func(*testing.T, *tls.Config)
	}{
		{
			name:    "✅ Nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "✅ Insecure mode",
			config: &ClientConfig{
				InsecureSkipVerify: true,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, cfg *tls.Config) {
				if !cfg.InsecureSkipVerify {
					t.Error("❌ InsecureSkipVerify should be true")
				}
				if cfg.MinVersion != MinTLSVersion {
					t.Errorf("❌ MinVersion = %d, want %d", cfg.MinVersion, MinTLSVersion)
				}
			},
		},
		{
			name: "✅ Valid CA file",
			config: &ClientConfig{
				CAFile:             caFile,
				InsecureSkipVerify: false,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, cfg *tls.Config) {
				if cfg.InsecureSkipVerify {
					t.Error("❌ InsecureSkipVerify should be false")
				}
				if cfg.RootCAs == nil {
					t.Error("❌ RootCAs should not be nil")
				}
			},
		},
		{
			name: "✅ Non-existent CA file (should not error, just skip)",
			config: &ClientConfig{
				CAFile:             "/nonexistent/ca.crt",
				InsecureSkipVerify: false,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, cfg *tls.Config) {
				if cfg.InsecureSkipVerify {
					t.Error("❌ InsecureSkipVerify should be false")
				}
				// RootCAs should be nil if CA file doesn't exist
				if cfg.RootCAs != nil {
					t.Error("❌ RootCAs should be nil when CA file doesn't exist")
				}
			},
		},
		{
			name: "✅ Custom MinVersion",
			config: &ClientConfig{
				MinVersion:         tls.VersionTLS13,
				InsecureSkipVerify: false,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, cfg *tls.Config) {
				if cfg.MinVersion != tls.VersionTLS13 {
					t.Errorf("❌ MinVersion = %d, want %d", cfg.MinVersion, tls.VersionTLS13)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewClientTLSConfig(tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("❌ NewClientTLSConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkFunc != nil {
				tt.checkFunc(t, cfg)
			}

			t.Log("✅ Test passed")
		})
	}
}

// TestDefaultServerConfig tests the DefaultServerConfig function
func TestDefaultServerConfig(t *testing.T) {
	t.Log("🧪 TEST: DefaultServerConfig")
	t.Log("==================================")

	certFile := "/path/to/cert.crt"
	keyFile := "/path/to/key.key"

	config := DefaultServerConfig(certFile, keyFile)

	if config == nil {
		t.Fatal("❌ DefaultServerConfig returned nil")
	}

	if config.CertFile != certFile {
		t.Errorf("❌ CertFile = %s, want %s", config.CertFile, certFile)
	}

	if config.KeyFile != keyFile {
		t.Errorf("❌ KeyFile = %s, want %s", config.KeyFile, keyFile)
	}

	if config.MinVersion != MinTLSVersion {
		t.Errorf("❌ MinVersion = %d, want %d", config.MinVersion, MinTLSVersion)
	}

	if len(config.CipherSuites) == 0 {
		t.Error("❌ CipherSuites should not be empty")
	}

	if !config.PreferServerCipherSuites {
		t.Error("❌ PreferServerCipherSuites should be true")
	}

	t.Log("✅ DefaultServerConfig test passed")
}

// TestDefaultClientConfig tests the DefaultClientConfig function
func TestDefaultClientConfig(t *testing.T) {
	t.Log("🧪 TEST: DefaultClientConfig")
	t.Log("==================================")

	config := DefaultClientConfig()

	if config == nil {
		t.Fatal("❌ DefaultClientConfig returned nil")
	}

	if config.MinVersion != MinTLSVersion {
		t.Errorf("❌ MinVersion = %d, want %d", config.MinVersion, MinTLSVersion)
	}

	if config.InsecureSkipVerify {
		t.Error("❌ InsecureSkipVerify should be false by default")
	}

	t.Log("✅ DefaultClientConfig test passed")
}
