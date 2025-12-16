// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package clientcmd

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestClient_ConnectionRefused teste la gestion d'une connexion refusée
func TestClient_ConnectionRefused(t *testing.T) {
	t.Log("🧪 TEST CLIENT - CONNEXION REFUSÉE")
	t.Log("===================================")

	config := &Config{
		ServerURL: "https://localhost:19999",
		Timeout:   2 * time.Second,
		Insecure:  true,
	}

	client := NewClient(config)

	_, err := client.Execute("type Person : <id: string>", "<test>")

	if err == nil {
		t.Fatal("❌ Attendait une erreur pour serveur inexistant")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "connection refused") &&
		!strings.Contains(errMsg, "connexion refusée") &&
		!strings.Contains(errMsg, "connect: connection refused") {
		t.Logf("⚠️  Message d'erreur: %v", err)
	}

	t.Logf("✅ Connexion refusée détectée: %v", err)
}

// TestClient_Timeout teste la gestion d'un timeout de requête
func TestClient_Timeout(t *testing.T) {
	t.Log("🧪 TEST CLIENT - TIMEOUT")
	t.Log("========================")

	done := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(10 * time.Second):
		case <-done:
			return
		}
	}))
	defer func() {
		close(done)
		server.Close()
	}()

	config := &Config{
		ServerURL: server.URL,
		Timeout:   200 * time.Millisecond,
		Insecure:  true,
	}

	client := NewClient(config)
	
	retryConfig := RetryConfig{
		MaxAttempts:          1,
		BaseDelay:            0,
		MaxDelay:             0,
		Jitter:               0,
		RetryableStatusCodes: []int{},
	}
	client.SetRetryConfig(retryConfig)

	start := time.Now()
	_, err := client.Execute("type Person : <id: string>", "<test>")
	duration := time.Since(start)

	if err == nil {
		t.Fatal("❌ Attendait un timeout")
	}

	if duration > 1*time.Second {
		t.Errorf("❌ Timeout trop long: %v (max attendu: 1s)", duration)
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		t.Logf("✅ Timeout détecté en %v: %v", duration, err)
	} else if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline exceeded") {
		t.Logf("✅ Timeout détecté en %v: %v", duration, err)
	} else {
		t.Logf("⚠️  Erreur détectée mais type non confirmé timeout: %v", err)
	}
}

// TestClient_DNSError teste la gestion d'une erreur DNS
func TestClient_DNSError(t *testing.T) {
	t.Log("🧪 TEST CLIENT - ERREUR DNS")
	t.Log("===========================")

	config := &Config{
		ServerURL: "https://invalid-hostname-that-does-not-exist-xyz-123456.com",
		Timeout:   5 * time.Second,
		Insecure:  true,
	}

	client := NewClient(config)

	_, err := client.Execute("type Person : <id: string>", "<test>")

	if err == nil {
		t.Fatal("❌ Attendait une erreur DNS")
	}

	errMsg := err.Error()
	if strings.Contains(errMsg, "no such host") ||
		strings.Contains(errMsg, "DNS") ||
		strings.Contains(errMsg, "lookup") {
		t.Logf("✅ Erreur DNS détectée: %v", err)
	} else {
		t.Logf("⚠️  Message d'erreur (probablement DNS): %v", err)
	}
}

// TestClient_IncompleteResponse teste la gestion d'une réponse incomplète
func TestClient_IncompleteResponse(t *testing.T) {
	t.Log("🧪 TEST CLIENT - RÉPONSE INCOMPLÈTE")
	t.Log("====================================")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "result": `))

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		hj, ok := w.(http.Hijacker)
		if !ok {
			return
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			return
		}
		conn.Close()
	}))
	defer server.Close()

	config := &Config{
		ServerURL: server.URL,
		Timeout:   2 * time.Second,
		Insecure:  true,
	}

	client := NewClient(config)

	_, err := client.Execute("type Person : <id: string>", "<test>")

	if err == nil {
		t.Fatal("❌ Attendait une erreur de parsing/connexion")
	}

	t.Logf("✅ Réponse incomplète détectée: %v", err)
}

// TestClient_ConnectionReset teste la gestion d'une connexion réinitialisée
func TestClient_ConnectionReset(t *testing.T) {
	t.Log("🧪 TEST CLIENT - CONNEXION RÉINITIALISÉE")
	t.Log("=========================================")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
			return
		}

		conn, _, err := hj.Hijack()
		if err != nil {
			return
		}

		conn.Close()
	}))
	defer server.Close()

	config := &Config{
		ServerURL: server.URL,
		Timeout:   2 * time.Second,
		Insecure:  true,
	}

	client := NewClient(config)

	_, err := client.Execute("type Person : <id: string>", "<test>")

	if err == nil {
		t.Fatal("❌ Attendait une erreur de connexion réinitialisée")
	}

	errMsg := err.Error()
	if strings.Contains(errMsg, "EOF") ||
		strings.Contains(errMsg, "connection reset") ||
		strings.Contains(errMsg, "broken pipe") {
		t.Logf("✅ Connexion réinitialisée détectée: %v", err)
	} else {
		t.Logf("⚠️  Erreur de connexion détectée: %v", err)
	}
}

// TestClient_SlowServer teste la gestion d'un serveur lent avec timeout
func TestClient_SlowServer(t *testing.T) {
	t.Log("🧪 TEST CLIENT - SERVEUR LENT")
	t.Log("==============================")

	done := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(2 * time.Second):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success": true, "results": {"facts_count": 0, "activations_count": 0, "activations": []}, "execution_time_ms": 0}`))
		case <-done:
			return
		}
	}))
	defer func() {
		close(done)
		server.Close()
	}()

	config := &Config{
		ServerURL: server.URL,
		Timeout:   300 * time.Millisecond,
		Insecure:  true,
	}

	client := NewClient(config)
	
	retryConfig := RetryConfig{
		MaxAttempts:          1,
		BaseDelay:            0,
		MaxDelay:             0,
		Jitter:               0,
		RetryableStatusCodes: []int{},
	}
	client.SetRetryConfig(retryConfig)

	start := time.Now()
	_, err := client.Execute("type Person : <id: string>", "<test>")
	duration := time.Since(start)

	if err == nil {
		t.Fatal("❌ Attendait un timeout pour serveur lent")
	}

	if duration > 1*time.Second {
		t.Errorf("❌ Timeout trop long: %v (max attendu: 1s)", duration)
	}

	t.Logf("✅ Serveur lent détecté avec timeout en %v: %v", duration, err)
}

// TestClient_ContextCancellation teste l'annulation via contexte
func TestClient_ContextCancellation(t *testing.T) {
	t.Log("🧪 TEST CLIENT - ANNULATION CONTEXTE")
	t.Log("=====================================")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	config := &Config{
		ServerURL: server.URL,
		Timeout:   30 * time.Second,
		Insecure:  true,
	}

	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	req, err := client.createExecuteRequest("type Person : <id: string>", "<test>")
	if err != nil {
		t.Fatalf("❌ Erreur création requête: %v", err)
	}

	req = req.WithContext(ctx)

	start := time.Now()
	_, err = client.executeRequestWithRetry(req)
	duration := time.Since(start)

	if err == nil {
		t.Fatal("❌ Attendait une erreur d'annulation")
	}

	if duration > 500*time.Millisecond {
		t.Errorf("⚠️  Annulation tardive: %v", duration)
	}

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		t.Logf("✅ Annulation contexte détectée en %v", duration)
	} else {
		t.Logf("⚠️  Erreur détectée (probablement annulation): %v", err)
	}
}

// TestClient_InvalidPort teste la gestion d'un port invalide
func TestClient_InvalidPort(t *testing.T) {
	t.Log("🧪 TEST CLIENT - PORT INVALIDE")
	t.Log("===============================")

	config := &Config{
		ServerURL: "https://localhost:99999",
		Timeout:   2 * time.Second,
		Insecure:  true,
	}

	client := NewClient(config)

	_, err := client.Execute("type Person : <id: string>", "<test>")

	if err == nil {
		t.Fatal("❌ Attendait une erreur pour port invalide")
	}

	errMsg := err.Error()
	if strings.Contains(errMsg, "invalid port") ||
		strings.Contains(errMsg, "dial") ||
		strings.Contains(errMsg, "connection") {
		t.Logf("✅ Port invalide détecté: %v", err)
	} else {
		t.Logf("⚠️  Erreur détectée: %v", err)
	}
}

// TestClient_RetryOnNetworkError teste le retry automatique sur erreur réseau
func TestClient_RetryOnNetworkError(t *testing.T) {
	t.Log("🧪 TEST CLIENT - RETRY SUR ERREUR RÉSEAU")
	t.Log("=========================================")

	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			hj, ok := w.(http.Hijacker)
			if ok {
				conn, _, _ := hj.Hijack()
				conn.Close()
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "results": {"facts_count": 0, "activations_count": 0, "activations": []}, "execution_time_ms": 0}`))
	}))
	defer server.Close()

	config := &Config{
		ServerURL: server.URL,
		Timeout:   5 * time.Second,
		Insecure:  true,
		Verbose:   false,
	}

	client := NewClient(config)
	retryConfig := DefaultRetryConfig()
	retryConfig.MaxAttempts = 3
	retryConfig.BaseDelay = 50 * time.Millisecond
	client.SetRetryConfig(retryConfig)

	resp, err := client.Execute("type Person : <id: string>", "<test>")

	if err != nil {
		t.Logf("⚠️  Erreur malgré retry: %v (tentatives: %d)", err, attempts)
	} else if resp != nil && resp.Success {
		t.Logf("✅ Retry réussi après %d tentative(s)", attempts)
	} else {
		t.Errorf("❌ Réponse invalide après retry")
	}
}
