// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package servercmd

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestValidateContentType_ValidJSON(t *testing.T) {
	t.Log("🧪 TEST VALIDATION CONTENT-TYPE - JSON VALIDE")
	t.Log("==============================================")

	server := &Server{
		config: &Config{},
		logger: log.New(io.Discard, "", 0),
	}

	handler := server.validateContentType(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name        string
		contentType string
		wantStatus  int
	}{
		{"application/json", "application/json", http.StatusOK},
		{"avec charset utf-8", "application/json; charset=utf-8", http.StatusOK},
		{"avec charset UTF-8", "application/json; charset=UTF-8", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte("{}")))
			req.Header.Set("Content-Type", tt.contentType)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("❌ Status = %d, attendu %d", rec.Code, tt.wantStatus)
			} else {
				t.Logf("✅ Content-Type accepté: %s", tt.contentType)
			}
		})
	}
}

func TestValidateContentType_Invalid(t *testing.T) {
	t.Log("🧪 TEST VALIDATION CONTENT-TYPE - INVALIDE")
	t.Log("==========================================")

	server := &Server{
		config: &Config{},
		logger: log.New(io.Discard, "", 0),
	}

	handler := server.validateContentType(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name        string
		contentType string
		wantStatus  int
	}{
		{"text/plain", "text/plain", StatusUnsupportedMedia},
		{"application/xml", "application/xml", StatusUnsupportedMedia},
		{"vide", "", StatusUnsupportedMedia},
		{"application/x-www-form-urlencoded", "application/x-www-form-urlencoded", StatusUnsupportedMedia},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte("{}")))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("❌ Status = %d, attendu %d", rec.Code, tt.wantStatus)
			} else {
				t.Logf("✅ Content-Type rejeté: %s", tt.contentType)
			}
		})
	}
}

func TestValidateContentType_GetMethod(t *testing.T) {
	t.Log("🧪 TEST VALIDATION CONTENT-TYPE - GET (EXEMPT)")
	t.Log("==============================================")

	server := &Server{
		config: &Config{},
		logger: log.New(io.Discard, "", 0),
	}

	handler := server.validateContentType(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// GET sans Content-Type devrait passer
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("❌ Status = %d, attendu 200", rec.Code)
	} else {
		t.Logf("✅ GET sans Content-Type accepté")
	}
}

func TestSendErrorResponse(t *testing.T) {
	t.Log("🧪 TEST ENVOI RÉPONSE D'ERREUR")
	t.Log("==============================")

	server := &Server{
		config: &Config{},
		logger: log.New(io.Discard, "", 0),
	}

	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{"bad request", StatusBadRequest, "Requête invalide"},
		{"unauthorized", StatusUnauthorized, "Non autorisé"},
		{"unsupported media", StatusUnsupportedMedia, "Content-Type non supporté"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			server.sendErrorResponse(rec, tt.statusCode, tt.message, time.Now())

			// Vérifier status
			if rec.Code != tt.statusCode {
				t.Errorf("❌ Status = %d, attendu %d", rec.Code, tt.statusCode)
			}

			// Vérifier Content-Type
			ct := rec.Header().Get("Content-Type")
			if ct != ContentTypeJSON {
				t.Errorf("❌ Content-Type = %s, attendu %s", ct, ContentTypeJSON)
			}

			t.Logf("✅ Réponse d'erreur envoyée correctement (status=%d)", tt.statusCode)
		})
	}
}

func TestWriteJSON_ContentType(t *testing.T) {
	t.Log("🧪 TEST ENVOI RÉPONSE JSON - CONTENT-TYPE")
	t.Log("=========================================")

	server := &Server{
		config: &Config{},
		logger: log.New(io.Discard, "", 0),
	}

	tests := []struct {
		name       string
		statusCode int
		data       interface{}
	}{
		{"succès avec data", StatusOK, map[string]string{"result": "ok"}},
		{"erreur avec data", StatusBadRequest, map[string]string{"error": "bad"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			server.writeJSON(rec, tt.data, tt.statusCode)

			// Vérifier status
			if rec.Code != tt.statusCode {
				t.Errorf("❌ Status = %d, attendu %d", rec.Code, tt.statusCode)
			}

			// Vérifier Content-Type
			ct := rec.Header().Get("Content-Type")
			if ct != ContentTypeJSON {
				t.Errorf("❌ Content-Type = %s, attendu %s", ct, ContentTypeJSON)
			}

			t.Logf("✅ Réponse JSON envoyée correctement (status=%d)", tt.statusCode)
		})
	}
}
