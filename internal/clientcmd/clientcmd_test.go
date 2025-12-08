// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package clientcmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/treivax/tsd/tsdio"
)

// Test constants
const (
	TestTimeout = 5 * time.Second
)

// TestParseFlags_Help vérifie le parsing du flag help
func TestParseFlags_Help(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantHelp bool
	}{
		{
			name:     "help flag",
			args:     []string{"-help"},
			wantHelp: true,
		},
		{
			name:     "no help flag",
			args:     []string{"-file", "test.tsd"},
			wantHelp: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := parseFlags(tt.args)
			if err != nil {
				t.Fatalf("parseFlags() error = %v", err)
			}

			if config.ShowHelp != tt.wantHelp {
				t.Errorf("ShowHelp = %v, want %v", config.ShowHelp, tt.wantHelp)
			}
		})
	}
}

// TestParseFlags_Sources vérifie le parsing des différentes sources
func TestParseFlags_Sources(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantFile  string
		wantText  string
		wantStdin bool
	}{
		{
			name:     "file flag",
			args:     []string{"-file", "test.tsd"},
			wantFile: "test.tsd",
		},
		{
			name:     "positional file",
			args:     []string{"program.tsd"},
			wantFile: "program.tsd",
		},
		{
			name:     "text flag",
			args:     []string{"-text", "type Person"},
			wantText: "type Person",
		},
		{
			name:      "stdin flag",
			args:      []string{"-stdin"},
			wantStdin: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := parseFlags(tt.args)
			if err != nil {
				t.Fatalf("parseFlags() error = %v", err)
			}

			if config.File != tt.wantFile {
				t.Errorf("File = %q, want %q", config.File, tt.wantFile)
			}
			if config.Text != tt.wantText {
				t.Errorf("Text = %q, want %q", config.Text, tt.wantText)
			}
			if config.UseStdin != tt.wantStdin {
				t.Errorf("UseStdin = %v, want %v", config.UseStdin, tt.wantStdin)
			}
		})
	}
}

// TestParseFlags_Options vérifie le parsing des options
func TestParseFlags_Options(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantServer  string
		wantVerbose bool
		wantFormat  string
		wantTimeout time.Duration
		wantHealth  bool
	}{
		{
			name:        "default options",
			args:        []string{"-file", "test.tsd"},
			wantServer:  DefaultServerURL,
			wantVerbose: false,
			wantFormat:  "text",
			wantTimeout: DefaultTimeout,
			wantHealth:  false,
		},
		{
			name:        "custom server",
			args:        []string{"-file", "test.tsd", "-server", "https://example.com:9000"},
			wantServer:  "https://example.com:9000",
			wantVerbose: false,
			wantFormat:  "text",
		},
		{
			name:        "verbose mode",
			args:        []string{"-file", "test.tsd", "-v"},
			wantVerbose: true,
		},
		{
			name:       "json format",
			args:       []string{"-file", "test.tsd", "-format", "json"},
			wantFormat: "json",
		},
		{
			name:        "custom timeout",
			args:        []string{"-file", "test.tsd", "-timeout", "60s"},
			wantTimeout: 60 * time.Second,
		},
		{
			name:       "health check",
			args:       []string{"-health"},
			wantHealth: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := parseFlags(tt.args)
			if err != nil {
				t.Fatalf("parseFlags() error = %v", err)
			}

			if tt.wantServer != "" && config.ServerURL != tt.wantServer {
				t.Errorf("ServerURL = %q, want %q", config.ServerURL, tt.wantServer)
			}
			if config.Verbose != tt.wantVerbose {
				t.Errorf("Verbose = %v, want %v", config.Verbose, tt.wantVerbose)
			}
			if tt.wantFormat != "" && config.Format != tt.wantFormat {
				t.Errorf("Format = %q, want %q", config.Format, tt.wantFormat)
			}
			if tt.wantTimeout != 0 && config.Timeout != tt.wantTimeout {
				t.Errorf("Timeout = %v, want %v", config.Timeout, tt.wantTimeout)
			}
			if config.ShowHealth != tt.wantHealth {
				t.Errorf("ShowHealth = %v, want %v", config.ShowHealth, tt.wantHealth)
			}
		})
	}
}

// TestParseFlags_TLS vérifie le parsing des options TLS
func TestParseFlags_TLS(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantCAFile   string
		wantInsecure bool
	}{
		{
			name:         "default TLS",
			args:         []string{"-file", "test.tsd"},
			wantCAFile:   DefaultCAFile,
			wantInsecure: false,
		},
		{
			name:       "custom CA",
			args:       []string{"-file", "test.tsd", "-tls-ca", "/custom/ca.crt"},
			wantCAFile: "/custom/ca.crt",
		},
		{
			name:         "insecure mode",
			args:         []string{"-file", "test.tsd", "-insecure"},
			wantCAFile:   DefaultCAFile, // Le CA file garde sa valeur par défaut
			wantInsecure: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := parseFlags(tt.args)
			if err != nil {
				t.Fatalf("parseFlags() error = %v", err)
			}

			if config.TLSCAFile != tt.wantCAFile {
				t.Errorf("TLSCAFile = %q, want %q", config.TLSCAFile, tt.wantCAFile)
			}
			if config.Insecure != tt.wantInsecure {
				t.Errorf("Insecure = %v, want %v", config.Insecure, tt.wantInsecure)
			}
		})
	}
}

// TestParseFlags_Auth vérifie le parsing des options d'authentification
func TestParseFlags_Auth(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		envToken     string
		wantToken    string
		wantAuthType string
	}{
		{
			name:      "token flag",
			args:      []string{"-file", "test.tsd", "-token", "mytoken123"},
			wantToken: "mytoken123",
		},
		{
			name:         "token with auth type",
			args:         []string{"-file", "test.tsd", "-token", "jwt123", "-auth-type", "jwt"},
			wantToken:    "jwt123",
			wantAuthType: "jwt",
		},
		{
			name:      "token from env",
			args:      []string{"-file", "test.tsd"},
			envToken:  "envtoken456",
			wantToken: "envtoken456",
		},
		{
			name:      "flag overrides env",
			args:      []string{"-file", "test.tsd", "-token", "flagtoken"},
			envToken:  "envtoken",
			wantToken: "flagtoken",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			oldToken := os.Getenv("TSD_AUTH_TOKEN")
			defer os.Setenv("TSD_AUTH_TOKEN", oldToken)

			if tt.envToken != "" {
				os.Setenv("TSD_AUTH_TOKEN", tt.envToken)
			} else {
				os.Unsetenv("TSD_AUTH_TOKEN")
			}

			config, err := parseFlags(tt.args)
			if err != nil {
				t.Fatalf("parseFlags() error = %v", err)
			}

			if config.AuthToken != tt.wantToken {
				t.Errorf("AuthToken = %q, want %q", config.AuthToken, tt.wantToken)
			}
			if tt.wantAuthType != "" && config.AuthType != tt.wantAuthType {
				t.Errorf("AuthType = %q, want %q", config.AuthType, tt.wantAuthType)
			}
		})
	}
}

// TestValidateConfig_NoSource vérifie la validation sans source
func TestValidateConfig_NoSource(t *testing.T) {
	config := &Config{}

	err := validateConfig(config)
	if err == nil {
		t.Fatal("validateConfig() expected error for no source, got nil")
	}

	if !strings.Contains(err.Error(), "aucune source") {
		t.Errorf("error message = %q, want to contain 'aucune source'", err.Error())
	}
}

// TestValidateConfig_MultipleSources vérifie la validation avec plusieurs sources
func TestValidateConfig_MultipleSources(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "file and text",
			config: &Config{
				File: "test.tsd",
				Text: "type Person",
			},
		},
		{
			name: "file and stdin",
			config: &Config{
				File:     "test.tsd",
				UseStdin: true,
			},
		},
		{
			name: "text and stdin",
			config: &Config{
				Text:     "type Person",
				UseStdin: true,
			},
		},
		{
			name: "all three",
			config: &Config{
				File:     "test.tsd",
				Text:     "type Person",
				UseStdin: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if err == nil {
				t.Fatal("validateConfig() expected error for multiple sources, got nil")
			}

			if !strings.Contains(err.Error(), "une seule source") {
				t.Errorf("error message = %q, want to contain 'une seule source'", err.Error())
			}
		})
	}
}

// TestValidateConfig_InvalidFormat vérifie la validation d'un format invalide
func TestValidateConfig_InvalidFormat(t *testing.T) {
	config := &Config{
		File:   "test.tsd",
		Format: "xml",
	}

	err := validateConfig(config)
	if err == nil {
		t.Fatal("validateConfig() expected error for invalid format, got nil")
	}

	if !strings.Contains(err.Error(), "format invalide") {
		t.Errorf("error message = %q, want to contain 'format invalide'", err.Error())
	}
}

// TestValidateConfig_Valid vérifie la validation avec une config valide
func TestValidateConfig_Valid(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "file source",
			config: &Config{
				File:   "test.tsd",
				Format: "text",
			},
		},
		{
			name: "text source",
			config: &Config{
				Text:   "type Person",
				Format: "json",
			},
		},
		{
			name: "stdin source",
			config: &Config{
				UseStdin: true,
				Format:   "text",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if err != nil {
				t.Errorf("validateConfig() unexpected error = %v", err)
			}
		})
	}
}

// TestReadSource_Stdin vérifie la lecture depuis stdin
func TestReadSource_Stdin(t *testing.T) {
	config := &Config{UseStdin: true}
	stdin := strings.NewReader("type Person : <id: string>")

	source, sourceName, err := readSource(config, stdin)
	if err != nil {
		t.Fatalf("readSource() error = %v", err)
	}

	if source != "type Person : <id: string>" {
		t.Errorf("source = %q, want %q", source, "type Person : <id: string>")
	}
	if sourceName != "<stdin>" {
		t.Errorf("sourceName = %q, want %q", sourceName, "<stdin>")
	}
}

// TestReadSource_Text vérifie la lecture depuis du texte direct
func TestReadSource_Text(t *testing.T) {
	config := &Config{Text: "type Order : <id: string>"}

	source, sourceName, err := readSource(config, nil)
	if err != nil {
		t.Fatalf("readSource() error = %v", err)
	}

	if source != "type Order : <id: string>" {
		t.Errorf("source = %q, want %q", source, "type Order : <id: string>")
	}
	if sourceName != "<text>" {
		t.Errorf("sourceName = %q, want %q", sourceName, "<text>")
	}
}

// TestReadSource_File vérifie la lecture depuis un fichier
func TestReadSource_File(t *testing.T) {
	// Créer un fichier temporaire
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.tsd")
	content := "type Product : <id: string, price: number>"

	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	config := &Config{File: tmpFile}

	source, sourceName, err := readSource(config, nil)
	if err != nil {
		t.Fatalf("readSource() error = %v", err)
	}

	if source != content {
		t.Errorf("source = %q, want %q", source, content)
	}
	if sourceName != tmpFile {
		t.Errorf("sourceName = %q, want %q", sourceName, tmpFile)
	}
}

// TestReadSource_FileNotFound vérifie l'erreur quand le fichier n'existe pas
func TestReadSource_FileNotFound(t *testing.T) {
	config := &Config{File: "/nonexistent/file.tsd"}

	_, _, err := readSource(config, nil)
	if err == nil {
		t.Fatal("readSource() expected error for nonexistent file, got nil")
	}

	if !strings.Contains(err.Error(), "fichier non trouvé") {
		t.Errorf("error message = %q, want to contain 'fichier non trouvé'", err.Error())
	}
}

// TestNewClient_Insecure vérifie la création d'un client en mode insecure
func TestNewClient_Insecure(t *testing.T) {
	config := &Config{
		ServerURL: "https://localhost:8080",
		Insecure:  true,
		Timeout:   TestTimeout,
	}

	client := NewClient(config)
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}

	if client.config != config {
		t.Error("client config not set correctly")
	}

	if client.tlsConfig == nil {
		t.Fatal("client tlsConfig is nil")
	}

	if !client.tlsConfig.InsecureSkipVerify {
		t.Error("InsecureSkipVerify should be true")
	}
}

// TestNewClient_WithCA vérifie la création d'un client avec CA
func TestNewClient_WithCA(t *testing.T) {
	// Créer un fichier CA temporaire
	tmpDir := t.TempDir()
	caFile := filepath.Join(tmpDir, "ca.crt")

	// Contenu CA de test (cert PEM bidon mais valide syntaxiquement)
	caContent := `-----BEGIN CERTIFICATE-----
MIIBkTCB+wIJAKHHCgVZU6N9MA0GCSqGSIb3DQEBCwUAMBExDzANBgNVBAMMBnRl
c3RjYTAeFw0yNTAxMDEwMDAwMDBaFw0yNjAxMDEwMDAwMDBaMBExDzANBgNVBAMM
BnRlc3RjYTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAwXvLRIyF8TJ9bCQ1
-----END CERTIFICATE-----`

	err := os.WriteFile(caFile, []byte(caContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create CA file: %v", err)
	}

	config := &Config{
		ServerURL: "https://localhost:8080",
		TLSCAFile: caFile,
		Timeout:   TestTimeout,
	}

	client := NewClient(config)
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}

	if client.tlsConfig.InsecureSkipVerify {
		t.Error("InsecureSkipVerify should be false")
	}
}

// TestClient_Execute vérifie l'exécution d'une requête
func TestClient_Execute(t *testing.T) {
	// Créer un serveur HTTP de test
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Vérifier la méthode
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Vérifier le Content-Type
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", ct)
		}

		// Vérifier l'authorization header
		if auth := r.Header.Get("Authorization"); auth != "Bearer testtoken123" {
			t.Errorf("Expected Authorization Bearer testtoken123, got %s", auth)
		}

		// Lire la requête
		var req tsdio.ExecuteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Vérifier les champs
		if req.Source != "type Person" {
			t.Errorf("Expected source 'type Person', got %s", req.Source)
		}
		if req.SourceName != "test.tsd" {
			t.Errorf("Expected sourceName 'test.tsd', got %s", req.SourceName)
		}

		// Envoyer une réponse
		response := tsdio.ExecuteResponse{
			Success:         true,
			ExecutionTimeMs: 42,
			Results: &tsdio.ExecutionResults{
				FactsCount:       1,
				ActivationsCount: 0,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &Config{
		ServerURL: server.URL,
		AuthToken: "testtoken123",
		Insecure:  true,
		Timeout:   TestTimeout,
	}

	client := NewClient(config)
	client.httpClient = server.Client()

	response, err := client.Execute("type Person", "test.tsd")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !response.Success {
		t.Error("Expected Success = true")
	}
	if response.ExecutionTimeMs != 42 {
		t.Errorf("ExecutionTimeMs = %d, want 42", response.ExecutionTimeMs)
	}
	if response.Results.FactsCount != 1 {
		t.Errorf("FactsCount = %d, want 1", response.Results.FactsCount)
	}
}

// TestClient_HealthCheck vérifie le health check
func TestClient_HealthCheck(t *testing.T) {
	// Créer un serveur HTTP de test
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Errorf("Expected path /health, got %s", r.URL.Path)
		}

		response := tsdio.HealthResponse{
			Status:        "healthy",
			Version:       "1.0.0",
			UptimeSeconds: 123,
			Timestamp:     time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &Config{
		ServerURL: server.URL,
		Insecure:  true,
		Timeout:   TestTimeout,
	}

	client := NewClient(config)
	client.httpClient = server.Client()

	health, err := client.HealthCheck()
	if err != nil {
		t.Fatalf("HealthCheck() error = %v", err)
	}

	if health.Status != "healthy" {
		t.Errorf("Status = %q, want 'healthy'", health.Status)
	}
	if health.Version != "1.0.0" {
		t.Errorf("Version = %q, want '1.0.0'", health.Version)
	}
	if health.UptimeSeconds != 123 {
		t.Errorf("UptimeSeconds = %d, want 123", health.UptimeSeconds)
	}
}

// TestPrintResults_JSON vérifie l'affichage en format JSON
func TestPrintResults_JSON(t *testing.T) {
	config := &Config{Format: "json"}

	response := &tsdio.ExecuteResponse{
		Success:         true,
		ExecutionTimeMs: 50,
		Results: &tsdio.ExecutionResults{
			FactsCount:       2,
			ActivationsCount: 1,
		},
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := printResults(config, response, &stdout, &stderr)
	if err != nil {
		t.Fatalf("printResults() error = %v", err)
	}

	// Vérifier que c'est du JSON valide
	var decoded tsdio.ExecuteResponse
	if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	if decoded.Success != response.Success {
		t.Errorf("Success = %v, want %v", decoded.Success, response.Success)
	}
}

// TestPrintResults_Text_Success vérifie l'affichage en format texte (succès)
func TestPrintResults_Text_Success(t *testing.T) {
	config := &Config{Format: "text", Verbose: false}

	response := &tsdio.ExecuteResponse{
		Success:         true,
		ExecutionTimeMs: 100,
		Results: &tsdio.ExecutionResults{
			FactsCount:       3,
			ActivationsCount: 0,
		},
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := printResults(config, response, &stdout, &stderr)
	if err != nil {
		t.Fatalf("printResults() error = %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "EXÉCUTION RÉUSSIE") {
		t.Error("Output should contain 'EXÉCUTION RÉUSSIE'")
	}
	if !strings.Contains(output, "100ms") {
		t.Error("Output should contain execution time")
	}
	if !strings.Contains(output, "Faits injectés: 3") {
		t.Error("Output should contain facts count")
	}
}

// TestPrintResults_Text_Error vérifie l'affichage en format texte (erreur)
func TestPrintResults_Text_Error(t *testing.T) {
	config := &Config{Format: "text"}

	response := &tsdio.ExecuteResponse{
		Success:         false,
		Error:           "syntax error",
		ErrorType:       "ParseError",
		ExecutionTimeMs: 5,
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := printResults(config, response, &stdout, &stderr)
	if err != nil {
		t.Fatalf("printResults() error = %v", err)
	}

	output := stderr.String()
	if !strings.Contains(output, "ERREUR D'EXÉCUTION") {
		t.Error("Output should contain 'ERREUR D'EXÉCUTION'")
	}
	if !strings.Contains(output, "syntax error") {
		t.Error("Output should contain error message")
	}
	if !strings.Contains(output, "ParseError") {
		t.Error("Output should contain error type")
	}
}

// TestPrintResults_Text_WithActivations vérifie l'affichage avec des activations
func TestPrintResults_Text_WithActivations(t *testing.T) {
	config := &Config{Format: "text", Verbose: true}

	response := &tsdio.ExecuteResponse{
		Success:         true,
		ExecutionTimeMs: 75,
		Results: &tsdio.ExecutionResults{
			FactsCount:       2,
			ActivationsCount: 1,
			Activations: []tsdio.Activation{
				{
					ActionName: "greet",
					Arguments: []tsdio.ArgumentValue{
						{Position: 0, Value: "Alice", Type: "string"},
					},
					TriggeringFacts: []tsdio.Fact{
						{
							ID:   "f1",
							Type: "Person",
							Fields: map[string]interface{}{
								"name": "Alice",
							},
						},
					},
				},
			},
		},
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := printResults(config, response, &stdout, &stderr)
	if err != nil {
		t.Fatalf("printResults() error = %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "ACTIONS DÉCLENCHÉES") {
		t.Error("Output should contain 'ACTIONS DÉCLENCHÉES'")
	}
	if !strings.Contains(output, "greet") {
		t.Error("Output should contain action name")
	}
	if !strings.Contains(output, "Alice") {
		t.Error("Output should contain argument value")
	}
	if !strings.Contains(output, "Person") {
		t.Error("Output should contain fact type")
	}
}

// TestRun_Help vérifie l'exécution avec l'aide
func TestRun_Help(t *testing.T) {
	args := []string{"-help"}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(args, nil, &stdout, &stderr)
	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "TSD Client") {
		t.Error("Help output should contain 'TSD Client'")
	}
	if !strings.Contains(output, "USAGE:") {
		t.Error("Help output should contain 'USAGE:'")
	}
}

// TestRun_ValidationError vérifie l'exécution avec une erreur de validation
func TestRun_ValidationError(t *testing.T) {
	args := []string{} // Pas de source
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(args, nil, &stdout, &stderr)
	if exitCode != 1 {
		t.Errorf("Run() exitCode = %d, want 1", exitCode)
	}

	output := stderr.String()
	if !strings.Contains(output, "aucune source") {
		t.Errorf("Error output should contain 'aucune source', got: %s", output)
	}
}

// TestRun_FileNotFound vérifie l'exécution avec un fichier non trouvé
func TestRun_FileNotFound(t *testing.T) {
	args := []string{"-file", "/nonexistent/file.tsd"}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(args, nil, &stdout, &stderr)
	if exitCode != 1 {
		t.Errorf("Run() exitCode = %d, want 1", exitCode)
	}

	output := stderr.String()
	if !strings.Contains(output, "fichier non trouvé") {
		t.Errorf("Error output should contain 'fichier non trouvé', got: %s", output)
	}
}

// TestRun_HealthCheck vérifie l'exécution du health check
func TestRun_HealthCheck(t *testing.T) {
	// Créer un serveur HTTP de test
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := tsdio.HealthResponse{
			Status:        "healthy",
			Version:       "1.0.0",
			UptimeSeconds: 100,
			Timestamp:     time.Now(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	args := []string{"-health", "-server", server.URL, "-insecure"}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(args, nil, &stdout, &stderr)
	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "healthy") {
		t.Error("Output should contain 'healthy'")
	}
}

// TestPrintHelp vérifie l'affichage de l'aide
func TestPrintHelp(t *testing.T) {
	var stdout bytes.Buffer
	printHelp(&stdout)

	output := stdout.String()

	expectedStrings := []string{
		"TSD Client",
		"USAGE:",
		"OPTIONS:",
		"EXEMPLES:",
		"-help",
		"-file",
		"-text",
		"-stdin",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Help output should contain %q", expected)
		}
	}
}
