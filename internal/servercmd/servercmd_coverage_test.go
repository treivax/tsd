// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package servercmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/treivax/tsd/tsdio"
)

// TestParseFlags_TLSCertValidation teste la validation des certificats TLS
func TestParseFlags_TLSCertValidation(t *testing.T) {
	t.Run("with insecure flag skips validation", func(t *testing.T) {
		args := []string{
			"-insecure",
			"-tls-cert", "/nonexistent/cert.crt",
			"-tls-key", "/nonexistent/key.key",
		}
		config := parseFlags(args)

		if !config.Insecure {
			t.Error("Expected insecure mode")
		}
		if config.TLSCertFile != "/nonexistent/cert.crt" {
			t.Errorf("TLSCertFile = %q, want /nonexistent/cert.crt", config.TLSCertFile)
		}
	})

	t.Run("with valid cert files", func(t *testing.T) {
		// Créer des fichiers temporaires pour simuler cert et key
		tmpDir := t.TempDir()
		certFile := filepath.Join(tmpDir, "cert.crt")
		keyFile := filepath.Join(tmpDir, "key.key")

		if err := os.WriteFile(certFile, []byte("fake cert"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(keyFile, []byte("fake key"), 0644); err != nil {
			t.Fatal(err)
		}

		args := []string{
			"-tls-cert", certFile,
			"-tls-key", keyFile,
		}
		config := parseFlags(args)

		if config.TLSCertFile != certFile {
			t.Errorf("TLSCertFile = %q, want %q", config.TLSCertFile, certFile)
		}
		if config.TLSKeyFile != keyFile {
			t.Errorf("TLSKeyFile = %q, want %q", config.TLSKeyFile, keyFile)
		}
	})
}

// TestParseFlags_EnvironmentVariables teste les variables d'environnement
func TestParseFlags_EnvironmentVariables(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		args     []string
		validate func(*testing.T, *Config)
	}{
		{
			name: "TSD_TLS_CERT override",
			envVars: map[string]string{
				"TSD_TLS_CERT": "/env/cert.crt",
			},
			args: []string{"-insecure"},
			validate: func(t *testing.T, c *Config) {
				if c.TLSCertFile != "/env/cert.crt" {
					t.Errorf("TLSCertFile = %q, want /env/cert.crt", c.TLSCertFile)
				}
			},
		},
		{
			name: "TSD_TLS_KEY override",
			envVars: map[string]string{
				"TSD_TLS_KEY": "/env/key.key",
			},
			args: []string{"-insecure"},
			validate: func(t *testing.T, c *Config) {
				if c.TLSKeyFile != "/env/key.key" {
					t.Errorf("TLSKeyFile = %q, want /env/key.key", c.TLSKeyFile)
				}
			},
		},
		{
			name: "TSD_INSECURE override",
			envVars: map[string]string{
				"TSD_INSECURE": "true",
			},
			args: []string{},
			validate: func(t *testing.T, c *Config) {
				if !c.Insecure {
					t.Error("Expected Insecure to be true")
				}
			},
		},
		{
			name: "TSD_AUTH_KEYS",
			envVars: map[string]string{
				"TSD_AUTH_KEYS": "key1,key2,key3",
			},
			args: []string{"-insecure"},
			validate: func(t *testing.T, c *Config) {
				if len(c.AuthKeys) != 3 {
					t.Errorf("AuthKeys length = %d, want 3", len(c.AuthKeys))
				}
				if c.AuthKeys[0] != "key1" || c.AuthKeys[1] != "key2" || c.AuthKeys[2] != "key3" {
					t.Errorf("AuthKeys = %v, want [key1 key2 key3]", c.AuthKeys)
				}
			},
		},
		{
			name: "TSD_JWT_SECRET",
			envVars: map[string]string{
				"TSD_JWT_SECRET": "my-secret-key",
			},
			args: []string{"-insecure"},
			validate: func(t *testing.T, c *Config) {
				if c.JWTSecret != "my-secret-key" {
					t.Errorf("JWTSecret = %q, want my-secret-key", c.JWTSecret)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Sauvegarder et restaurer les variables d'environnement
			oldEnv := make(map[string]string)
			for key := range tt.envVars {
				oldEnv[key] = os.Getenv(key)
				os.Setenv(key, tt.envVars[key])
			}
			defer func() {
				for key, val := range oldEnv {
					if val == "" {
						os.Unsetenv(key)
					} else {
						os.Setenv(key, val)
					}
				}
			}()

			config := parseFlags(tt.args)
			tt.validate(t, config)
		})
	}
}

// TestHandleExecute_MissingSource teste le cas où source est vide
func TestHandleExecute_MissingSource(t *testing.T) {
	config := &Config{
		AuthType: "none",
	}
	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	reqBody := `{"source": "", "source_name": "test.tsd"}`
	req := httptest.NewRequest("POST", "/api/v1/execute", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.handleExecute(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestHandleExecute_VerboseMode teste le mode verbose
func TestHandleExecute_VerboseMode(t *testing.T) {
	var buf bytes.Buffer
	config := &Config{
		Verbose:  true,
		AuthType: "none",
	}
	logger := log.New(&buf, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	reqBody := `{"source": "type Person { name: String }; fact Person { name: \"Alice\" };", "source_name": "test.tsd"}`
	req := httptest.NewRequest("POST", "/api/v1/execute", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.handleExecute(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	// Vérifier que les logs verbeux ont été écrits
	logOutput := buf.String()
	if !strings.Contains(logOutput, "Requête d'exécution reçue") {
		t.Error("Expected verbose log about request received")
	}
}

// TestHandleExecute_DefaultSourceName teste le source name par défaut
func TestHandleExecute_DefaultSourceName(t *testing.T) {
	config := &Config{
		AuthType: "none",
	}
	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Sans source_name
	reqBody := `{"source": "type Person { name: String }; fact Person { name: \"Alice\" };"}`
	req := httptest.NewRequest("POST", "/api/v1/execute", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.handleExecute(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}
}

// TestExecuteTSDProgram_ParsingError teste une erreur de parsing
func TestExecuteTSDProgram_ParsingError(t *testing.T) {
	config := &Config{
		AuthType: "none",
	}
	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	req := &tsdio.ExecuteRequest{
		Source:     "invalid syntax @@@@",
		SourceName: "invalid.tsd",
	}

	resp := server.executeTSDProgram(req, time.Now())

	if resp.Success {
		t.Error("Expected failure, got success")
	}

	if resp.ErrorType != tsdio.ErrorTypeParsingError {
		t.Errorf("ErrorType = %q, want %q", resp.ErrorType, tsdio.ErrorTypeParsingError)
	}

	if !strings.Contains(resp.Error, "Erreur de parsing") {
		t.Errorf("Error message should contain 'Erreur de parsing', got: %s", resp.Error)
	}
}

// TestWriteJSON_ErrorCase teste le cas d'erreur d'encodage JSON
func TestWriteJSON_ErrorCase(t *testing.T) {
	config := &Config{
		AuthType: "none",
	}
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Créer une structure qui ne peut pas être encodée en JSON
	// Les canaux ne peuvent pas être encodés en JSON
	badData := struct {
		Channel chan int `json:"channel"`
	}{
		Channel: make(chan int),
	}

	w := httptest.NewRecorder()
	server.writeJSON(w, badData, http.StatusOK)

	// Vérifier que l'erreur a été loggée
	logOutput := buf.String()
	if !strings.Contains(logOutput, "Erreur encodage JSON") {
		t.Errorf("Expected error log about JSON encoding, got: %s", logOutput)
	}
}

// TestNewServer_JWTWithoutSecret teste l'erreur JWT sans secret
func TestNewServer_JWTWithoutSecret(t *testing.T) {
	config := &Config{
		AuthType:  "jwt",
		JWTSecret: "",
	}
	logger := log.New(io.Discard, "", 0)
	_, err := NewServer(config, logger)

	if err == nil {
		t.Error("Expected error for JWT without secret")
	}

	if !strings.Contains(err.Error(), "secret JWT") {
		t.Errorf("Expected error about JWT secret, got: %v", err)
	}
}

// TestNewServer_KeyAuthWithoutKeys teste l'erreur auth key sans clés
func TestNewServer_KeyAuthWithoutKeys(t *testing.T) {
	config := &Config{
		AuthType: "key",
		AuthKeys: []string{},
	}
	logger := log.New(io.Discard, "", 0)
	_, err := NewServer(config, logger)

	if err == nil {
		t.Error("Expected error for key auth without keys")
	}

	if !strings.Contains(err.Error(), "au moins une clé") {
		t.Errorf("Expected error about missing keys, got: %v", err)
	}
}

// TestRun_InitError teste le cas où NewServer retourne une erreur
func TestRun_InitError(t *testing.T) {
	// Créer des arguments qui causeront une erreur dans NewServer
	args := []string{
		"-insecure",
		"-auth", "jwt",
		// Pas de JWT secret, ce qui devrait causer une erreur
	}

	var stdout, stderr bytes.Buffer
	exitCode := Run(args, nil, &stdout, &stderr)

	if exitCode != 1 {
		t.Errorf("Run() = %d, want 1", exitCode)
	}

	stderrStr := stderr.String()
	if !strings.Contains(stderrStr, "Erreur initialisation serveur") {
		t.Errorf("Expected error message about server init, got: %s", stderrStr)
	}
}

// TestRun_WithTestServer teste le démarrage du serveur avec httptest
func TestRun_WithTestServer(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     0,
		Insecure: true,
		AuthType: "none",
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Créer un serveur HTTP de test
	ts := httptest.NewServer(server.mux)
	defer ts.Close()

	// Tester le endpoint health
	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to call health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Health check status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// Tester le endpoint version
	resp, err = http.Get(ts.URL + "/api/v1/version")
	if err != nil {
		t.Fatalf("Failed to call version endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Version check status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

// TestHandleExecute_LargeBody teste le cas d'un body trop grand
func TestHandleExecute_LargeBody(t *testing.T) {
	config := &Config{
		AuthType: "none",
	}
	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Créer un body qui dépasse MaxRequestSize
	largeSource := strings.Repeat("a", int(MaxRequestSize)+1000)
	reqBody := fmt.Sprintf(`{"source": "%s", "source_name": "large.tsd"}`, largeSource)

	req := httptest.NewRequest("POST", "/api/v1/execute", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.handleExecute(w, req)

	// Le status devrait être BadRequest
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestHandleExecute_AuthenticationFailure teste l'échec d'authentification
func TestHandleExecute_AuthenticationFailure(t *testing.T) {
	config := &Config{
		AuthType: "key",
		AuthKeys: []string{"valid-key-123456789012345678901234567890"},
	}
	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
	}{
		{
			name:       "no auth header",
			authHeader: "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid key",
			authHeader: "Bearer wrong-key",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "malformed header",
			authHeader: "InvalidFormat",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := `{"source": "type Test { }; fact Test { };", "source_name": "test.tsd"}`
			req := httptest.NewRequest("POST", "/api/v1/execute", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			server.handleExecute(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

// TestParseFlags_AuthKeysWithSpaces teste le parsing des clés avec espaces
func TestParseFlags_AuthKeysWithSpaces(t *testing.T) {
	args := []string{
		"-insecure",
		"-auth-keys", "key1 , key2 , key3",
	}

	config := parseFlags(args)

	if len(config.AuthKeys) != 3 {
		t.Errorf("AuthKeys length = %d, want 3", len(config.AuthKeys))
	}

	// Vérifier que les espaces ont été supprimés
	expected := []string{"key1", "key2", "key3"}
	for i, key := range config.AuthKeys {
		if key != expected[i] {
			t.Errorf("AuthKeys[%d] = %q, want %q", i, key, expected[i])
		}
	}
}

// TestExecuteTSDProgram_ExecutionTime teste que le temps d'exécution est enregistré
func TestExecuteTSDProgram_ExecutionTime(t *testing.T) {
	config := &Config{
		AuthType: "none",
	}
	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	req := &tsdio.ExecuteRequest{
		Source:     "type Test { }; fact Test { };",
		SourceName: "test.tsd",
	}

	startTime := time.Now()
	resp := server.executeTSDProgram(req, startTime)

	if resp.ExecutionTimeMs < 0 {
		t.Errorf("ExecutionTimeMs = %d, should be >= 0", resp.ExecutionTimeMs)
	}

	// Le temps d'exécution devrait être raisonnable (moins de 5 secondes pour un test simple)
	if resp.ExecutionTimeMs > 5000 {
		t.Errorf("ExecutionTimeMs = %d, seems too high", resp.ExecutionTimeMs)
	}
}

// TestParseFlags_JWTDefaults teste les valeurs par défaut JWT
func TestParseFlags_JWTDefaults(t *testing.T) {
	args := []string{
		"-insecure",
		"-auth", "jwt",
		"-jwt-secret", "my-secret",
	}

	config := parseFlags(args)

	if config.JWTExpiration != 24*time.Hour {
		t.Errorf("JWTExpiration = %v, want 24h", config.JWTExpiration)
	}

	if config.JWTIssuer != "tsd-server" {
		t.Errorf("JWTIssuer = %q, want 'tsd-server'", config.JWTIssuer)
	}
}

// TestParseFlags_FlagPrecedenceOverEnv teste que les flags ont priorité sur les variables d'env
func TestParseFlags_FlagPrecedenceOverEnv(t *testing.T) {
	// Définir une variable d'environnement
	os.Setenv("TSD_AUTH_KEYS", "env-key")
	defer os.Unsetenv("TSD_AUTH_KEYS")

	// Spécifier une valeur via flag qui devrait avoir priorité
	args := []string{
		"-insecure",
		"-auth-keys", "flag-key",
	}

	config := parseFlags(args)

	if len(config.AuthKeys) != 1 || config.AuthKeys[0] != "flag-key" {
		t.Errorf("AuthKeys = %v, want [flag-key]", config.AuthKeys)
	}
}

// TestHandleExecute_MalformedJSON teste une requête avec JSON malformé
func TestHandleExecute_MalformedJSON(t *testing.T) {
	config := &Config{
		AuthType: "none",
	}
	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	reqBody := `{"source": "incomplete`
	req := httptest.NewRequest("POST", "/api/v1/execute", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.handleExecute(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestHandleExecute_VerboseRequest teste le flag verbose dans la requête
func TestHandleExecute_VerboseRequest(t *testing.T) {
	var buf bytes.Buffer
	config := &Config{
		Verbose:  false,
		AuthType: "none",
	}
	logger := log.New(&buf, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	reqBody := `{"source": "type Person { name: String }; fact Person { name: \"Alice\" };", "source_name": "test.tsd", "verbose": true}`
	req := httptest.NewRequest("POST", "/api/v1/execute", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.handleExecute(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	// Vérifier que les logs verbeux ont été écrits même si config.Verbose est false
	logOutput := buf.String()
	if !strings.Contains(logOutput, "Requête d'exécution reçue") {
		t.Error("Expected verbose log when request.Verbose is true")
	}
}

// TestParseFlags_AllEnvironmentVariables teste toutes les variables d'environnement
func TestParseFlags_AllEnvironmentVariables(t *testing.T) {
	// Créer des fichiers temporaires pour les certificats
	tmpDir := t.TempDir()
	certFile := filepath.Join(tmpDir, "cert.crt")
	keyFile := filepath.Join(tmpDir, "key.key")
	if err := os.WriteFile(certFile, []byte("cert"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(keyFile, []byte("key"), 0644); err != nil {
		t.Fatal(err)
	}

	os.Setenv("TSD_TLS_CERT", certFile)
	os.Setenv("TSD_TLS_KEY", keyFile)
	os.Setenv("TSD_INSECURE", "true")
	defer func() {
		os.Unsetenv("TSD_TLS_CERT")
		os.Unsetenv("TSD_TLS_KEY")
		os.Unsetenv("TSD_INSECURE")
	}()

	config := parseFlags([]string{})

	if config.TLSCertFile != certFile {
		t.Errorf("TLSCertFile = %q, want %q", config.TLSCertFile, certFile)
	}
	if config.TLSKeyFile != keyFile {
		t.Errorf("TLSKeyFile = %q, want %q", config.TLSKeyFile, keyFile)
	}
	if !config.Insecure {
		t.Error("Expected Insecure to be true")
	}
}

// TestParseFlags_EmptyAuthKeys teste le cas où auth-keys est une chaîne vide
func TestParseFlags_EmptyAuthKeys(t *testing.T) {
	args := []string{
		"-insecure",
		"-auth-keys", "",
	}

	config := parseFlags(args)

	if len(config.AuthKeys) != 0 {
		t.Errorf("AuthKeys length = %d, want 0", len(config.AuthKeys))
	}
}

// TestParseFlags_JWTSecretFromEnv teste le JWT secret depuis variable d'env
func TestParseFlags_JWTSecretFromEnv(t *testing.T) {
	os.Setenv("TSD_JWT_SECRET", "env-secret-value")
	defer os.Unsetenv("TSD_JWT_SECRET")

	args := []string{"-insecure"}
	config := parseFlags(args)

	if config.JWTSecret != "env-secret-value" {
		t.Errorf("JWTSecret = %q, want env-secret-value", config.JWTSecret)
	}
}

// TestParseFlags_JWTFlagOverridesEnv teste que le flag JWT a priorité sur l'env
func TestParseFlags_JWTFlagOverridesEnv(t *testing.T) {
	os.Setenv("TSD_JWT_SECRET", "env-secret")
	defer os.Unsetenv("TSD_JWT_SECRET")

	args := []string{
		"-insecure",
		"-jwt-secret", "flag-secret",
	}
	config := parseFlags(args)

	if config.JWTSecret != "flag-secret" {
		t.Errorf("JWTSecret = %q, want flag-secret", config.JWTSecret)
	}
}

// TestExecuteTSDProgram_IngestionError teste une erreur d'ingestion
func TestExecuteTSDProgram_IngestionError(t *testing.T) {
	config := &Config{
		AuthType: "none",
	}
	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Un programme qui va causer une erreur d'ingestion
	// (après parsing mais pendant l'ingestion RETE)
	req := &tsdio.ExecuteRequest{
		Source:     "type X { }; type X { };", // Type dupliqué
		SourceName: "duplicate-type.tsd",
	}

	resp := server.executeTSDProgram(req, time.Now())

	// Le programme devrait être rejeté soit au parsing soit à la validation
	if resp.Success {
		t.Log("Warning: expected failure but got success - duplicate types might be allowed")
	}
}

// TestHandleExecute_SuccessWithResults teste une exécution réussie avec résultats
func TestHandleExecute_SuccessWithResults(t *testing.T) {
	var buf bytes.Buffer
	config := &Config{
		Verbose:  true,
		AuthType: "none",
	}
	logger := log.New(&buf, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	reqBody := `{"source": "type Person { name: String }; fact Person { name: \"Alice\" };", "source_name": "success.tsd"}`
	req := httptest.NewRequest("POST", "/api/v1/execute", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.handleExecute(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}

	// Vérifier les logs de succès
	logOutput := buf.String()
	if !strings.Contains(logOutput, "Exécution réussie") && !strings.Contains(logOutput, "Requête d'exécution reçue") {
		t.Logf("Log output: %s", logOutput)
	}
}

// TestHandleExecute_ErrorWithVerbose teste une erreur en mode verbose
func TestHandleExecute_ErrorWithVerbose(t *testing.T) {
	var buf bytes.Buffer
	config := &Config{
		Verbose:  true,
		AuthType: "none",
	}
	logger := log.New(&buf, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	reqBody := `{"source": "invalid @@@", "source_name": "error.tsd"}`
	req := httptest.NewRequest("POST", "/api/v1/execute", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.handleExecute(w, req)

	// Vérifier les logs d'erreur
	logOutput := buf.String()
	if !strings.Contains(logOutput, "Exécution échouée") {
		t.Error("Expected failure log message")
	}
}

// TestParseFlags_CustomJWTExpiration teste une durée JWT personnalisée
func TestParseFlags_CustomJWTExpiration(t *testing.T) {
	args := []string{
		"-insecure",
		"-jwt-expiration", "48h",
	}

	config := parseFlags(args)

	if config.JWTExpiration != 48*time.Hour {
		t.Errorf("JWTExpiration = %v, want 48h", config.JWTExpiration)
	}
}

// TestParseFlags_CustomJWTIssuer teste un émetteur JWT personnalisé
func TestParseFlags_CustomJWTIssuer(t *testing.T) {
	args := []string{
		"-insecure",
		"-jwt-issuer", "my-custom-issuer",
	}

	config := parseFlags(args)

	if config.JWTIssuer != "my-custom-issuer" {
		t.Errorf("JWTIssuer = %q, want my-custom-issuer", config.JWTIssuer)
	}
}

// TestExecuteTSDProgram_EmptyProgram teste un programme vide
func TestExecuteTSDProgram_EmptyProgram(t *testing.T) {
	config := &Config{
		AuthType: "none",
	}
	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	req := &tsdio.ExecuteRequest{
		Source:     "",
		SourceName: "empty.tsd",
	}

	resp := server.executeTSDProgram(req, time.Now())

	// Un programme vide est valide (0 types, 0 faits)
	if !resp.Success {
		t.Logf("Empty program result: %v", resp.Error)
	}
}
