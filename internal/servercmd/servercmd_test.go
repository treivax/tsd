// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package servercmd

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/tsdio"
)

// Test constants
const (
	TestTimeout = 5 * time.Second
)

// TestParseFlags_Defaults vérifie les valeurs par défaut
func TestParseFlags_Defaults(t *testing.T) {
	args := []string{"-insecure"}
	config := parseFlags(args)

	if config.Host != DefaultHost {
		t.Errorf("Host = %q, want %q", config.Host, DefaultHost)
	}
	if config.Port != DefaultPort {
		t.Errorf("Port = %d, want %d", config.Port, DefaultPort)
	}
	if config.Verbose {
		t.Error("Verbose should be false by default")
	}
	if config.AuthType != "none" {
		t.Errorf("AuthType = %q, want 'none'", config.AuthType)
	}
}

// TestParseFlags_CustomValues vérifie les valeurs personnalisées
func TestParseFlags_CustomValues(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantHost    string
		wantPort    int
		wantVerbose bool
		wantAuth    string
	}{
		{
			name:     "custom host",
			args:     []string{"-host", "127.0.0.1", "-insecure"},
			wantHost: "127.0.0.1",
			wantPort: DefaultPort,
		},
		{
			name:     "custom port",
			args:     []string{"-port", "9000", "-insecure"},
			wantHost: DefaultHost,
			wantPort: 9000,
		},
		{
			name:        "verbose mode",
			args:        []string{"-v", "-insecure"},
			wantVerbose: true,
		},
		{
			name:     "key auth",
			args:     []string{"-auth", "key", "-insecure"},
			wantAuth: "key",
		},
		{
			name:     "jwt auth",
			args:     []string{"-auth", "jwt", "-insecure"},
			wantAuth: "jwt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseFlags(tt.args)

			if tt.wantHost != "" && config.Host != tt.wantHost {
				t.Errorf("Host = %q, want %q", config.Host, tt.wantHost)
			}
			if tt.wantPort != 0 && config.Port != tt.wantPort {
				t.Errorf("Port = %d, want %d", config.Port, tt.wantPort)
			}
			if config.Verbose != tt.wantVerbose {
				t.Errorf("Verbose = %v, want %v", config.Verbose, tt.wantVerbose)
			}
			if tt.wantAuth != "" && config.AuthType != tt.wantAuth {
				t.Errorf("AuthType = %q, want %q", config.AuthType, tt.wantAuth)
			}
		})
	}
}

// TestParseFlags_TLS vérifie les options TLS
func TestParseFlags_TLS(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantInsecure bool
		wantCert     string
		wantKey      string
	}{
		{
			name:         "insecure mode",
			args:         []string{"-insecure"},
			wantInsecure: true,
		},
		{
			name:         "custom cert",
			args:         []string{"-tls-cert", "/custom/cert.crt", "-insecure"},
			wantCert:     "/custom/cert.crt",
			wantInsecure: true,
		},
		{
			name:         "custom key",
			args:         []string{"-tls-key", "/custom/key.key", "-insecure"},
			wantKey:      "/custom/key.key",
			wantInsecure: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseFlags(tt.args)

			if config.Insecure != tt.wantInsecure {
				t.Errorf("Insecure = %v, want %v", config.Insecure, tt.wantInsecure)
			}
			if tt.wantCert != "" && config.TLSCertFile != tt.wantCert {
				t.Errorf("TLSCertFile = %q, want %q", config.TLSCertFile, tt.wantCert)
			}
			if tt.wantKey != "" && config.TLSKeyFile != tt.wantKey {
				t.Errorf("TLSKeyFile = %q, want %q", config.TLSKeyFile, tt.wantKey)
			}
		})
	}
}

// TestParseFlags_JWT vérifie les options JWT
func TestParseFlags_JWT(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantSecret     string
		wantExpiration time.Duration
		wantIssuer     string
	}{
		{
			name:       "jwt secret",
			args:       []string{"-jwt-secret", "mysecret123", "-insecure"},
			wantSecret: "mysecret123",
		},
		{
			name:           "jwt expiration",
			args:           []string{"-jwt-expiration", "48h", "-insecure"},
			wantExpiration: 48 * time.Hour,
		},
		{
			name:       "jwt issuer",
			args:       []string{"-jwt-issuer", "my-app", "-insecure"},
			wantIssuer: "my-app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := parseFlags(tt.args)

			if tt.wantSecret != "" && config.JWTSecret != tt.wantSecret {
				t.Errorf("JWTSecret = %q, want %q", config.JWTSecret, tt.wantSecret)
			}
			if tt.wantExpiration != 0 && config.JWTExpiration != tt.wantExpiration {
				t.Errorf("JWTExpiration = %v, want %v", config.JWTExpiration, tt.wantExpiration)
			}
			if tt.wantIssuer != "" && config.JWTIssuer != tt.wantIssuer {
				t.Errorf("JWTIssuer = %q, want %q", config.JWTIssuer, tt.wantIssuer)
			}
		})
	}
}

// TestParseFlags_AuthKeys vérifie le parsing des clés d'authentification
func TestParseFlags_AuthKeys(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		envValue  string
		wantKeys  []string
		wantCount int
	}{
		{
			name:      "single key from flag",
			args:      []string{"-auth-keys", "key1", "-insecure"},
			wantKeys:  []string{"key1"},
			wantCount: 1,
		},
		{
			name:      "multiple keys from flag",
			args:      []string{"-auth-keys", "key1,key2,key3", "-insecure"},
			wantKeys:  []string{"key1", "key2", "key3"},
			wantCount: 3,
		},
		{
			name:      "keys from env",
			args:      []string{"-insecure"},
			envValue:  "envkey1,envkey2",
			wantKeys:  []string{"envkey1", "envkey2"},
			wantCount: 2,
		},
		{
			name:      "keys with spaces",
			args:      []string{"-auth-keys", "key1 , key2 , key3", "-insecure"},
			wantKeys:  []string{"key1", "key2", "key3"},
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			oldEnv := os.Getenv("TSD_AUTH_KEYS")
			defer os.Setenv("TSD_AUTH_KEYS", oldEnv)

			if tt.envValue != "" {
				os.Setenv("TSD_AUTH_KEYS", tt.envValue)
			} else {
				os.Unsetenv("TSD_AUTH_KEYS")
			}

			config := parseFlags(tt.args)

			if len(config.AuthKeys) != tt.wantCount {
				t.Errorf("len(AuthKeys) = %d, want %d", len(config.AuthKeys), tt.wantCount)
			}

			for i, want := range tt.wantKeys {
				if i >= len(config.AuthKeys) {
					t.Errorf("Missing key at index %d: want %q", i, want)
					continue
				}
				if config.AuthKeys[i] != want {
					t.Errorf("AuthKeys[%d] = %q, want %q", i, config.AuthKeys[i], want)
				}
			}
		})
	}
}

// TestNewServer_NoAuth vérifie la création d'un serveur sans authentification
func TestNewServer_NoAuth(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)

	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	if server == nil {
		t.Fatal("NewServer() returned nil")
	}

	if server.config != config {
		t.Error("server config not set correctly")
	}

	if server.authManager == nil {
		t.Fatal("authManager is nil")
	}

	if server.authManager.IsEnabled() {
		t.Error("authManager should not be enabled")
	}
}

// TestNewServer_WithKeyAuth vérifie la création d'un serveur avec authentification par clé
func TestNewServer_WithKeyAuth(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "key",
		AuthKeys: []string{"testkey123456789012345678901234567890"},
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)

	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	if !server.authManager.IsEnabled() {
		t.Error("authManager should be enabled")
	}

	if server.authManager.GetAuthType() != "key" {
		t.Errorf("authType = %q, want 'key'", server.authManager.GetAuthType())
	}
}

// TestNewServer_WithJWTAuth vérifie la création d'un serveur avec authentification JWT
func TestNewServer_WithJWTAuth(t *testing.T) {
	config := &Config{
		Host:          "localhost",
		Port:          8080,
		AuthType:      "jwt",
		JWTSecret:     "secret123456789012345678901234567890",
		JWTExpiration: 24 * time.Hour,
		JWTIssuer:     "test-server",
		Insecure:      true,
	}

	logger := log.New(io.Discard, "", 0)

	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	if !server.authManager.IsEnabled() {
		t.Error("authManager should be enabled")
	}

	if server.authManager.GetAuthType() != "jwt" {
		t.Errorf("authType = %q, want 'jwt'", server.authManager.GetAuthType())
	}
}

// TestHandleHealth vérifie le handler /health
func TestHandleHealth(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	server.handleHealth(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	var health tsdio.HealthResponse
	if err := json.Unmarshal(body, &health); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if health.Status != "ok" {
		t.Errorf("Status = %q, want 'ok'", health.Status)
	}

	if health.Version != Version {
		t.Errorf("Version = %q, want %q", health.Version, Version)
	}

	if health.UptimeSeconds < 0 {
		t.Errorf("UptimeSeconds = %d, should be >= 0", health.UptimeSeconds)
	}
}

// TestHandleVersion vérifie le handler /api/v1/version
func TestHandleVersion(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/version", nil)
	w := httptest.NewRecorder()

	server.handleVersion(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	var version map[string]interface{}
	if err := json.Unmarshal(body, &version); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if v, ok := version["version"].(string); !ok || v != Version {
		t.Errorf("version = %v, want %q", version["version"], Version)
	}
}

// TestHandleExecute_NoAuth vérifie l'exécution sans authentification
func TestHandleExecute_NoAuth(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Programme TSD simple avec commentaire
	program := `// Type definition
type Person(id: string, name: string)
`

	reqBody := tsdio.ExecuteRequest{
		Source:     program,
		SourceName: "test.tsd",
		Verbose:    false,
	}

	jsonData, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/execute", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.handleExecute(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	var execResp tsdio.ExecuteResponse
	if err := json.Unmarshal(body, &execResp); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if !execResp.Success {
		t.Errorf("Success = false, error: %s", execResp.Error)
	}
}

// TestHandleExecute_WithAuth vérifie l'exécution avec authentification
func TestHandleExecute_WithAuth(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "key",
		AuthKeys: []string{"validkey123456789012345678901234567890"},
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	program := `// Person type
type Person(id: string)
`

	reqBody := tsdio.ExecuteRequest{
		Source:     program,
		SourceName: "test.tsd",
	}

	tests := []struct {
		name           string
		token          string
		wantStatusCode int
		wantSuccess    bool
	}{
		{
			name:           "valid token",
			token:          "validkey123456789012345678901234567890",
			wantStatusCode: http.StatusOK,
			wantSuccess:    true,
		},
		{
			name:           "invalid token",
			token:          "invalidkey",
			wantStatusCode: http.StatusUnauthorized,
			wantSuccess:    false,
		},
		{
			name:           "no token",
			token:          "",
			wantStatusCode: http.StatusUnauthorized,
			wantSuccess:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("POST", "/api/v1/execute", bytes.NewReader(jsonData))
			req.Header.Set("Content-Type", "application/json")

			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			w := httptest.NewRecorder()
			server.handleExecute(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatusCode {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatusCode)
			}

			if tt.wantSuccess {
				body, _ := io.ReadAll(resp.Body)
				var execResp tsdio.ExecuteResponse
				if err := json.Unmarshal(body, &execResp); err != nil {
					t.Fatalf("Failed to decode JSON: %v", err)
				}

				if !execResp.Success {
					t.Errorf("Success = false, error: %s", execResp.Error)
				}
			}
		})
	}
}

// TestHandleExecute_InvalidJSON vérifie le comportement avec un JSON invalide
func TestHandleExecute_InvalidJSON(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest("POST", "/api/v1/execute", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.handleExecute(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}

// TestHandleExecute_ParseError vérifie le comportement avec une erreur de parsing
func TestHandleExecute_ParseError(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Programme TSD avec erreur de syntaxe
	program := `type Person(invalid syntax here)`

	reqBody := tsdio.ExecuteRequest{
		Source:     program,
		SourceName: "test.tsd",
	}

	jsonData, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/execute", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.handleExecute(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	var execResp tsdio.ExecuteResponse
	if err := json.Unmarshal(body, &execResp); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if execResp.Success {
		t.Error("Success should be false for invalid TSD program")
	}

	if execResp.Error == "" {
		t.Error("Error message should not be empty")
	}
}

// TestAuthenticate vérifie la fonction d'authentification
func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name      string
		authType  string
		authKeys  []string
		token     string
		wantAllow bool
	}{
		{
			name:      "no auth required",
			authType:  "none",
			wantAllow: true,
		},
		{
			name:      "valid key",
			authType:  "key",
			authKeys:  []string{"key123456789012345678901234567890AB", "key223456789012345678901234567890CD"},
			token:     "key123456789012345678901234567890AB",
			wantAllow: true,
		},
		{
			name:      "invalid key",
			authType:  "key",
			authKeys:  []string{"key123456789012345678901234567890AB", "key223456789012345678901234567890CD"},
			token:     "wrongkey",
			wantAllow: false,
		},
		{
			name:      "no token with auth required",
			authType:  "key",
			authKeys:  []string{"key123456789012345678901234567890AB"},
			token:     "",
			wantAllow: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Host:     "localhost",
				Port:     8080,
				AuthType: tt.authType,
				AuthKeys: tt.authKeys,
				Insecure: true,
			}

			logger := log.New(io.Discard, "", 0)
			server, err := NewServer(config, logger)
			if err != nil {
				t.Fatalf("NewServer() error = %v", err)
			}

			req := httptest.NewRequest("GET", "/test", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			err = server.authenticate(req)
			allowed := (err == nil)
			if allowed != tt.wantAllow {
				t.Errorf("authenticate() allowed = %v, want %v (err: %v)", allowed, tt.wantAllow, err)
			}
		})
	}
}

// TestGetValueType vérifie la détection du type de valeur
func TestGetValueType(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	tests := []struct {
		name  string
		value interface{}
		want  string
	}{
		{
			name:  "string",
			value: "hello",
			want:  "string",
		},
		{
			name:  "int",
			value: 42,
			want:  "int",
		},
		{
			name:  "float64",
			value: 3.14,
			want:  "float",
		},
		{
			name:  "bool",
			value: true,
			want:  "bool",
		},
		{
			name:  "map",
			value: map[string]interface{}{"key": "value"},
			want:  "unknown",
		},
		{
			name:  "slice",
			value: []interface{}{1, 2, 3},
			want:  "unknown",
		},
		{
			name:  "nil",
			value: nil,
			want:  "nil",
		},
		{
			name:  "unknown",
			value: struct{}{},
			want:  "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := server.getValueType(tt.value)
			if got != tt.want {
				t.Errorf("getValueType(%v) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}

// TestWriteJSON vérifie l'écriture JSON
func TestWriteJSON(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"message": "test",
		"value":   42,
	}

	server.writeJSON(w, data, http.StatusOK)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want 'application/json'", ct)
	}

	body, _ := io.ReadAll(resp.Body)
	var decoded map[string]interface{}
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if decoded["message"] != "test" {
		t.Errorf("message = %v, want 'test'", decoded["message"])
	}
}

// TestWriteError vérifie l'écriture d'erreur
func TestWriteError(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	w := httptest.NewRecorder()
	server.writeError(w, tsdio.ErrorTypeServerError, "test error", http.StatusBadRequest, time.Now())

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}

	body, _ := io.ReadAll(resp.Body)
	var decoded map[string]interface{}
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if decoded["error"] != "test error" {
		t.Errorf("error = %v, want 'test error'", decoded["error"])
	}
}

// TestRegisterRoutes vérifie l'enregistrement des routes
func TestRegisterRoutes(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Vérifier que les routes sont enregistrées
	tests := []struct {
		path   string
		method string
	}{
		{"/health", "GET"},
		{"/api/v1/version", "GET"},
		{"/api/v1/execute", "POST"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			server.mux.ServeHTTP(w, req)

			// Vérifier que la route existe (pas de 404)
			if w.Code == http.StatusNotFound {
				t.Errorf("Route %s %s not found", tt.method, tt.path)
			}
		})
	}
}

// TestEnvironmentVariables vérifie la prise en compte des variables d'environnement
func TestEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name     string
		envKey   string
		envValue string
		checkFn  func(*Config) bool
	}{
		{
			name:     "TSD_TLS_CERT",
			envKey:   "TSD_TLS_CERT",
			envValue: "/custom/cert.crt",
			checkFn: func(c *Config) bool {
				return c.TLSCertFile == "/custom/cert.crt"
			},
		},
		{
			name:     "TSD_TLS_KEY",
			envKey:   "TSD_TLS_KEY",
			envValue: "/custom/key.key",
			checkFn: func(c *Config) bool {
				return c.TLSKeyFile == "/custom/key.key"
			},
		},
		{
			name:     "TSD_INSECURE",
			envKey:   "TSD_INSECURE",
			envValue: "true",
			checkFn: func(c *Config) bool {
				return c.Insecure
			},
		},
		{
			name:     "TSD_JWT_SECRET",
			envKey:   "TSD_JWT_SECRET",
			envValue: "envsecret",
			checkFn: func(c *Config) bool {
				return c.JWTSecret == "envsecret"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldValue := os.Getenv(tt.envKey)
			defer os.Setenv(tt.envKey, oldValue)

			os.Setenv(tt.envKey, tt.envValue)

			config := parseFlags([]string{"-insecure"})

			if !tt.checkFn(config) {
				t.Errorf("Environment variable %s not applied correctly", tt.envKey)
			}
		})
	}
}

// TestExecuteTSDProgram_Simple vérifie l'exécution d'un programme simple
func TestExecuteTSDProgram_Simple(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Programme simple qui définit un type avec commentaire
	program := `// Person type
type Person(id: string, name: string)
`

	req := &tsdio.ExecuteRequest{
		Source:     program,
		SourceName: "test.tsd",
	}

	result := server.executeTSDProgram(req, time.Now())
	if result == nil {
		t.Fatal("result is nil")
	}

	if !result.Success {
		t.Errorf("executeTSDProgram() failed: %s", result.Error)
	}
}

// TestExtractAttributes vérifie l'extraction des attributs
func TestExtractAttributes(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Créer un fait RETE avec des champs
	fact := &rete.Fact{
		ID:   "f1",
		Type: "Person",
		Fields: map[string]interface{}{
			"name":   "Alice",
			"age":    30,
			"active": true,
		},
	}

	result := server.extractAttributes(fact)

	if len(result) != 3 {
		t.Errorf("len(result) = %d, want 3", len(result))
	}

	if result["name"] != "Alice" {
		t.Errorf("name = %v, want 'Alice'", result["name"])
	}
	if result["age"] != 30 {
		t.Errorf("age = %v, want 30", result["age"])
	}
	if result["active"] != true {
		t.Errorf("active = %v, want true", result["active"])
	}
}

// TestCollectActivations_WithRealNetwork vérifie la collecte des activations avec un vrai réseau RETE
func TestCollectActivations_WithRealNetwork(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Programme TSD avec une règle et des faits pour générer des activations
	program := `// Test program with rule and facts
type Person(id: string, name: string, age: number)

action notify_adult(name: string, age: number)

rule adult_check : {p: Person} / p.age >= 18 ==> notify_adult(p.name, p.age)

Person(id:p1, name:Alice, age:25)
Person(id:p2, name:Bob, age:30)
`

	req := &tsdio.ExecuteRequest{
		Source:     program,
		SourceName: "test_activations.tsd",
		Verbose:    false,
	}

	// Exécuter le programme pour créer le réseau RETE
	response := server.executeTSDProgram(req, time.Now())
	if !response.Success {
		t.Fatalf("executeTSDProgram() failed: %s", response.Error)
	}

	// Vérifier que nous avons des résultats
	if response.Results == nil {
		t.Fatal("Results is nil")
	}

	// Vérifier que nous avons des activations (2 faits correspondant à la règle)
	if response.Results.ActivationsCount == 0 {
		t.Error("Expected activations, got 0")
	}

	if len(response.Results.Activations) != response.Results.ActivationsCount {
		t.Errorf("len(Activations) = %d, want %d", len(response.Results.Activations), response.Results.ActivationsCount)
	}

	// Vérifier que les activations contiennent les bonnes informations
	for i, activation := range response.Results.Activations {
		if activation.ActionName != "notify_adult" {
			t.Errorf("Activation[%d].ActionName = %q, want 'notify_adult'", i, activation.ActionName)
		}

		if len(activation.TriggeringFacts) == 0 {
			t.Errorf("Activation[%d] has no triggering facts", i)
		}

		// Vérifier que les faits déclencheurs sont des Person
		for _, fact := range activation.TriggeringFacts {
			if fact.Type != "Person" {
				t.Errorf("Fact type = %q, want 'Person'", fact.Type)
			}
		}
	}
}

// TestCollectActivations vérifie la collecte des activations
func TestCollectActivations(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Test avec un réseau nil
	activations := server.collectActivations(nil)
	if len(activations) != 0 {
		t.Errorf("len(activations) = %d, want 0 for nil network", len(activations))
	}

	// Test avec un réseau vide
	network := &rete.ReteNetwork{}
	activations = server.collectActivations(network)
	if len(activations) != 0 {
		t.Errorf("len(activations) = %d, want 0 for empty network", len(activations))
	}
}

// TestExtractFacts vérifie l'extraction des faits
func TestExtractFacts(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Créer un token avec des faits
	token := &rete.Token{
		Facts: []*rete.Fact{
			{
				ID:   "f1",
				Type: "Person",
				Fields: map[string]interface{}{
					"name": "Alice",
					"age":  30,
				},
			},
			{
				ID:   "f2",
				Type: "Order",
				Fields: map[string]interface{}{
					"id":     "o123",
					"amount": 100.50,
				},
			},
		},
	}

	facts := server.extractFacts(token)

	if len(facts) != 2 {
		t.Fatalf("len(facts) = %d, want 2", len(facts))
	}

	if facts[0].ID != "f1" {
		t.Errorf("facts[0].ID = %q, want 'f1'", facts[0].ID)
	}
	if facts[0].Type != "Person" {
		t.Errorf("facts[0].Type = %q, want 'Person'", facts[0].Type)
	}
	if facts[0].Attributes["name"] != "Alice" {
		t.Errorf("facts[0].Attributes[name] = %v, want 'Alice'", facts[0].Attributes["name"])
	}

	if facts[1].ID != "f2" {
		t.Errorf("facts[1].ID = %q, want 'f2'", facts[1].ID)
	}
}

// TestExtractArguments vérifie l'extraction des arguments
func TestExtractArguments(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Créer un terminal node avec une action
	terminal := &rete.TerminalNode{
		Action: &rete.Action{
			Job: &rete.JobCall{
				Name: "test_action",
				Args: []interface{}{
					"string_value",
					42,
					3.14,
					true,
					nil,
				},
			},
		},
	}

	// Créer un token
	token := &rete.Token{}

	args := server.extractArguments(terminal, token)

	if len(args) != 5 {
		t.Fatalf("len(args) = %d, want 5", len(args))
	}

	// Vérifier que les arguments sont extraits (ils sont convertis en string "expression")
	for i, arg := range args {
		if arg.Position != i {
			t.Errorf("args[%d].Position = %d, want %d", i, arg.Position, i)
		}
		if arg.Type != "expression" {
			t.Errorf("args[%d].Type = %q, want 'expression'", i, arg.Type)
		}
	}
}

// TestHandleHealth_InvalidMethod vérifie le rejet des méthodes non-GET
func TestHandleHealth_InvalidMethod(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest("POST", "/health", nil)
	w := httptest.NewRecorder()

	server.handleHealth(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	}
}

// TestHandleVersion_InvalidMethod vérifie le rejet des méthodes non-GET
func TestHandleVersion_InvalidMethod(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest("POST", "/api/v1/version", nil)
	w := httptest.NewRecorder()

	server.handleVersion(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	}
}

// TestHandleExecute_MethodNotAllowed vérifie le rejet des méthodes non-POST
func TestHandleExecute_MethodNotAllowed(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/execute", nil)
	w := httptest.NewRecorder()

	server.handleExecute(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	}
}

// TestHandleExecute_TooLarge vérifie le rejet des requêtes trop grandes
func TestHandleExecute_TooLarge(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Créer une requête avec un body trop grand
	largeData := make([]byte, MaxRequestSize+1)
	req := httptest.NewRequest("POST", "/api/v1/execute", bytes.NewReader(largeData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.handleExecute(w, req)

	resp := w.Result()
	// MaxBytesReader renvoie une erreur qui est capturée lors du décodage JSON,
	// ce qui produit un 400 Bad Request au lieu de 413
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want %d (for oversized request)", resp.StatusCode, http.StatusBadRequest)
	}
}

// TestExecuteTSDProgram_WithFacts vérifie l'exécution avec des faits
func TestExecuteTSDProgram_WithFacts(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Programme avec des faits
	program := `// Test with facts
type Order(id: string, amount: number, status: string)

Order(id:o1, amount:100.50, status:pending)
Order(id:o2, amount:250.00, status:completed)
Order(id:o3, amount:75.25, status:pending)
`

	req := &tsdio.ExecuteRequest{
		Source:     program,
		SourceName: "test_with_facts.tsd",
	}

	result := server.executeTSDProgram(req, time.Now())
	if !result.Success {
		t.Fatalf("executeTSDProgram() failed: %s", result.Error)
	}

	// Vérifier que nous avons des résultats
	if result.Results == nil {
		t.Fatal("Results is nil")
	}

	// Vérifier que les faits sont injectés
	if result.Results.FactsCount != 3 {
		t.Errorf("FactsCount = %d, want 3", result.Results.FactsCount)
	}
}

// TestExecuteTSDProgram_WithRule vérifie l'exécution avec une règle
func TestExecuteTSDProgram_WithRule(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Programme avec une règle qui devrait se déclencher
	program := `// Test with rule and facts
type Product(id: string, price: number, quantity: number)

action alert_expensive(id: string, price: number)

rule check_expensive : {p: Product} / p.price > 100 ==> alert_expensive(p.id, p.price)

Product(id:p1, price:150.00, quantity:5)
Product(id:p2, price:50.00, quantity:10)
Product(id:p3, price:200.00, quantity:0)
`

	req := &tsdio.ExecuteRequest{
		Source:     program,
		SourceName: "test_with_rule.tsd",
	}

	result := server.executeTSDProgram(req, time.Now())
	if !result.Success {
		t.Fatalf("executeTSDProgram() failed: %s", result.Error)
	}

	// Vérifier que nous avons des résultats
	if result.Results == nil {
		t.Fatal("Results is nil")
	}

	// Vérifier que les faits sont injectés
	if result.Results.FactsCount != 3 {
		t.Errorf("FactsCount = %d, want 3", result.Results.FactsCount)
	}

	// Vérifier qu'il y a des activations (p1 et p3 devraient correspondre car prix > 100)
	if result.Results.ActivationsCount < 1 {
		t.Error("Expected at least 1 activation for expensive product")
	}

	// Vérifier que l'action est correcte
	if len(result.Results.Activations) > 0 {
		if result.Results.Activations[0].ActionName != "alert_expensive" {
			t.Errorf("ActionName = %q, want 'alert_expensive'", result.Results.Activations[0].ActionName)
		}
	}
}

// TestHandleExecute_ComplexProgram vérifie l'exécution d'un programme complexe
func TestHandleExecute_ComplexProgram(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Programme complexe avec plusieurs types, règles et faits
	program := `// Complex program
type Customer(id: string, name: string, points: number)
type Order(customerId: string, amount: number, discount: number)

action apply_vip_discount(customerId: string, amount: number)

rule vip_discount : {c: Customer, o: Order} / c.points > 100 AND o.customerId == c.id AND o.discount == 0 ==> apply_vip_discount(o.customerId, o.amount)

Customer(id:c1, name:Alice, points:500)
Customer(id:c2, name:Bob, points:50)
Order(customerId:c1, amount:500.00, discount:0)
Order(customerId:c2, amount:300.00, discount:0)
`

	reqBody := tsdio.ExecuteRequest{
		Source:     program,
		SourceName: "complex_test.tsd",
		Verbose:    false,
	}

	jsonData, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/execute", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.handleExecute(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, _ := io.ReadAll(resp.Body)
	var execResp tsdio.ExecuteResponse
	if err := json.Unmarshal(body, &execResp); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if !execResp.Success {
		t.Errorf("Success = false, error: %s", execResp.Error)
	}

	// Vérifier que nous avons des résultats
	if execResp.Results == nil {
		t.Fatal("Results is nil")
	}

	// Devrait avoir 4 faits
	if execResp.Results.FactsCount != 4 {
		t.Errorf("FactsCount = %d, want 4", execResp.Results.FactsCount)
	}

	// Devrait avoir au moins une activation (customer avec points > 100 et order)
	if execResp.Results.ActivationsCount == 0 {
		t.Error("Expected at least 1 activation for VIP discount rule")
	}
}
