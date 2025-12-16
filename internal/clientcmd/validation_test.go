// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package clientcmd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateResponse_ValidContentType(t *testing.T) {
	t.Log("🧪 TEST VALIDATION RÉPONSE - CONTENT-TYPE VALIDE")
	t.Log("=================================================")

	tests := []struct {
		name        string
		contentType string
		wantErr     bool
	}{
		{"application/json", "application/json", false},
		{"avec charset utf-8", "application/json; charset=utf-8", false},
		{"avec charset UTF-8", "application/json; charset=UTF-8", false},
		{"text/plain", "text/plain", true},
		{"application/xml", "application/xml", true},
		{"vide", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
			}
			if tt.contentType != "" {
				resp.Header.Set("Content-Type", tt.contentType)
			}

			err := validateResponse(resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("❌ Erreur = %v, wantErr %v", err, tt.wantErr)
			} else if tt.wantErr {
				t.Logf("✅ Content-Type invalide rejeté: %s", tt.contentType)
			} else {
				t.Logf("✅ Content-Type valide accepté: %s", tt.contentType)
			}
		})
	}
}

func TestExecute_StatusCodes(t *testing.T) {
	t.Log("🧪 TEST EXECUTE - GESTION STATUS CODES")
	t.Log("=======================================")

	tests := []struct {
		name         string
		statusCode   int
		responseBody string
		wantErr      bool
		errContains  string
	}{
		{
			name:         "200 OK",
			statusCode:   200,
			responseBody: `{"success": true, "execution_time_ms": 10}`,
			wantErr:      false,
		},
		{
			name:         "400 Bad Request",
			statusCode:   400,
			responseBody: `{"error": "program vide"}`,
			wantErr:      true,
			errContains:  "requête invalide",
		},
		{
			name:         "401 Unauthorized",
			statusCode:   401,
			responseBody: `{"error": "token expiré"}`,
			wantErr:      true,
			errContains:  "non autorisé",
		},
		{
			name:         "415 Unsupported Media Type",
			statusCode:   415,
			responseBody: `{"error": "Content-Type invalide"}`,
			wantErr:      true,
			errContains:  "Content-Type non supporté",
		},
		{
			name:         "500 Internal Server Error",
			statusCode:   500,
			responseBody: `{"error": "erreur compilation"}`,
			wantErr:      true,
			errContains:  "erreur serveur",
		},
		{
			name:         "503 Service Unavailable",
			statusCode:   503,
			responseBody: `{"error": "serveur surchargé"}`,
			wantErr:      true,
			errContains:  "serveur indisponible",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock serveur
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			config := &Config{
				ServerURL: server.URL,
				Timeout:   DefaultTimeout,
			}
			client := NewClient(config)

			// Désactiver retry pour les tests (pour éviter les retries sur 500/503)
			client.SetRetryConfig(RetryConfig{
				MaxAttempts:          1,
				BaseDelay:            0,
				MaxDelay:             0,
				Jitter:               0,
				RetryableStatusCodes: []int{},
			})

			_, err := client.Execute("test program", "<test>")

			if (err != nil) != tt.wantErr {
				t.Errorf("❌ Erreur = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("❌ Erreur ne contient pas '%s': %v", tt.errContains, err)
				} else {
					t.Logf("✅ Status %d géré correctement: %v", tt.statusCode, err)
				}
			} else if !tt.wantErr {
				t.Logf("✅ Succès pour status %d", tt.statusCode)
			}
		})
	}
}

func TestExecute_InvalidContentType(t *testing.T) {
	t.Log("🧪 TEST EXECUTE - CONTENT-TYPE INVALIDE")
	t.Log("========================================")

	// Mock serveur qui retourne text/plain
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("This is plain text"))
	}))
	defer server.Close()

	config := &Config{
		ServerURL: server.URL,
		Timeout:   DefaultTimeout,
	}
	client := NewClient(config)

	// Désactiver retry
	client.SetRetryConfig(RetryConfig{
		MaxAttempts:          1,
		BaseDelay:            0,
		MaxDelay:             0,
		Jitter:               0,
		RetryableStatusCodes: []int{},
	})

	_, err := client.Execute("test program", "<test>")

	if err == nil {
		t.Errorf("❌ Attendu une erreur pour Content-Type invalide")
		return
	}

	if !strings.Contains(err.Error(), "Content-Type inattendu") {
		t.Errorf("❌ Message d'erreur incorrect: %v", err)
		return
	}

	t.Logf("✅ Content-Type invalide rejeté correctement: %v", err)
}

func TestParseErrorResponse(t *testing.T) {
	t.Log("🧪 TEST PARSE ERROR RESPONSE")
	t.Log("============================")

	tests := []struct {
		name        string
		body        string
		contentType string
		want        string
	}{
		{
			name:        "erreur avec champ error",
			body:        `{"error": "message d'erreur"}`,
			contentType: "application/json",
			want:        "message d'erreur",
		},
		{
			name:        "erreur avec champ message",
			body:        `{"message": "autre message"}`,
			contentType: "application/json",
			want:        "autre message",
		},
		{
			name:        "JSON invalide",
			body:        "not json",
			contentType: "application/json",
			want:        "Erreur HTTP 400",
		},
		{
			name:        "JSON vide",
			body:        "{}",
			contentType: "application/json",
			want:        "Erreur HTTP 400",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tt.contentType)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(tt.body))
			}))
			defer server.Close()

			resp, err := http.Get(server.URL)
			if err != nil {
				t.Fatalf("❌ Erreur HTTP Get: %v", err)
			}
			defer resp.Body.Close()

			got := parseErrorResponse(resp)

			if got != tt.want {
				t.Errorf("❌ parseErrorResponse() = %v, want %v", got, tt.want)
			} else {
				t.Logf("✅ Message d'erreur parsé: %s", got)
			}
		})
	}
}
