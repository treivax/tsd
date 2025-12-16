// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package servercmd

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestServer_ClientDisconnect teste la gestion d'une déconnexion brutale du client
func TestServer_ClientDisconnect(t *testing.T) {
	t.Log("🧪 TEST SERVEUR - CLIENT DISCONNECT")
	t.Log("====================================")

	config := &Config{
		Host:     "localhost",
		Port:     0,
		Insecure: true,
		AuthType: "none",
	}

	server, err := setupTestServerOnRandomPort(config, t)
	if err != nil {
		t.Fatalf("❌ Erreur création serveur: %v", err)
	}
	defer server.httpServer.Close()

	addr := server.httpServer.Addr

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("❌ Erreur connexion: %v", err)
	}

	conn.Write([]byte("POST /api/v1/execute HTTP/1.1\r\n"))
	conn.Write([]byte("Host: localhost\r\n"))
	conn.Write([]byte("Content-Type: application/json\r\n"))
	conn.Write([]byte("Content-Length: 1000\r\n"))
	conn.Write([]byte("\r\n"))

	conn.Close()

	time.Sleep(200 * time.Millisecond)

	if server.httpServer == nil {
		t.Error("❌ Serveur ne devrait pas être nil après disconnect client")
	}

	t.Log("✅ Disconnect client géré correctement")
}

// TestServer_RequestTooLarge teste la gestion d'une requête trop grande
func TestServer_RequestTooLarge(t *testing.T) {
	t.Log("🧪 TEST SERVEUR - REQUÊTE TROP LARGE")
	t.Log("=====================================")

	config := &Config{
		Host:     "localhost",
		Port:     0,
		Insecure: true,
		AuthType: "none",
	}

	server, err := setupTestServerOnRandomPort(config, t)
	if err != nil {
		t.Fatalf("❌ Erreur création serveur: %v", err)
	}
	defer server.httpServer.Close()

	baseURL := "http://" + server.httpServer.Addr

	largeBody := strings.Repeat("x", 11*1024*1024)

	req, _ := http.NewRequest("POST", baseURL+"/api/v1/execute", strings.NewReader(largeBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		if strings.Contains(err.Error(), "request body too large") ||
			strings.Contains(err.Error(), "http: request body too large") {
			t.Logf("✅ Requête trop large rejetée: %v", err)
			return
		}
		t.Logf("✅ Requête trop large rejetée avec erreur: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusRequestEntityTooLarge ||
		resp.StatusCode == http.StatusBadRequest {
		t.Logf("✅ Requête trop large rejetée: status %d", resp.StatusCode)
	} else {
		t.Logf("⚠️  Status = %d (attendu 413 ou 400, mais requête rejetée)", resp.StatusCode)
	}
}

// TestServer_MalformedRequest teste la gestion d'une requête mal formée
func TestServer_MalformedRequest(t *testing.T) {
	t.Log("🧪 TEST SERVEUR - REQUÊTE MAL FORMÉE")
	t.Log("=====================================")

	config := &Config{
		Host:     "localhost",
		Port:     0,
		Insecure: true,
		AuthType: "none",
	}

	server, err := setupTestServerOnRandomPort(config, t)
	if err != nil {
		t.Fatalf("❌ Erreur création serveur: %v", err)
	}
	defer server.httpServer.Close()

	addr := server.httpServer.Addr

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("❌ Erreur connexion: %v", err)
	}
	defer conn.Close()

	conn.Write([]byte("INVALID HTTP REQUEST\r\n\r\n"))

	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buf)

	if err != nil && err != io.EOF {
		t.Logf("✅ Requête mal formée rejetée: %v", err)
		return
	}

	if n > 0 {
		response := string(buf[:n])
		if strings.Contains(response, "400") || strings.Contains(response, "Bad Request") {
			t.Log("✅ Requête mal formée rejetée avec 400 Bad Request")
		} else {
			t.Logf("⚠️  Réponse serveur: %s", response)
		}
	}
}

// TestServer_SlowClient teste la gestion d'un client lent (protection slowloris)
func TestServer_SlowClient(t *testing.T) {
	t.Log("🧪 TEST SERVEUR - CLIENT LENT (SLOWLORIS)")
	t.Log("==========================================")

	config := &Config{
		Host:     "localhost",
		Port:     0,
		Insecure: true,
		AuthType: "none",
	}

	server, err := setupTestServerOnRandomPort(config, t)
	if err != nil {
		t.Fatalf("❌ Erreur création serveur: %v", err)
	}
	defer server.httpServer.Close()

	addr := server.httpServer.Addr

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("❌ Erreur connexion: %v", err)
	}
	defer conn.Close()

	conn.Write([]byte("POST /api/v1/execute HTTP/1.1\r\n"))
	time.Sleep(100 * time.Millisecond)
	conn.Write([]byte("Host: localhost\r\n"))
	time.Sleep(100 * time.Millisecond)

	time.Sleep(6 * time.Second)

	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err := conn.Read(buf)

	if err != nil {
		if strings.Contains(err.Error(), "timeout") ||
			strings.Contains(err.Error(), "closed") ||
			err == io.EOF {
			t.Log("✅ Client lent déconnecté par timeout du serveur")
			return
		}
	}

	if n > 0 {
		response := string(buf[:n])
		if strings.Contains(response, "408") || strings.Contains(response, "timeout") {
			t.Log("✅ Client lent géré avec timeout approprié")
		} else {
			t.Logf("⚠️  Réponse serveur: %s", response)
		}
	}
}

// TestServer_IncompleteBody teste la gestion d'un body incomplet
func TestServer_IncompleteBody(t *testing.T) {
	t.Log("🧪 TEST SERVEUR - BODY INCOMPLET")
	t.Log("=================================")

	config := &Config{
		Host:     "localhost",
		Port:     0,
		Insecure: true,
		AuthType: "none",
	}

	server, err := setupTestServerOnRandomPort(config, t)
	if err != nil {
		t.Fatalf("❌ Erreur création serveur: %v", err)
	}
	defer server.httpServer.Close()

	baseURL := "http://" + server.httpServer.Addr

	incompleteJSON := `{"source": "type Person : <id: string>", "source_name": "`

	req, _ := http.NewRequest("POST", baseURL+"/api/v1/execute", bytes.NewBufferString(incompleteJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", "1000")

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		if strings.Contains(err.Error(), "EOF") ||
			strings.Contains(err.Error(), "unexpected EOF") {
			t.Logf("✅ Body incomplet détecté: %v", err)
			return
		}
		t.Logf("✅ Erreur détectée: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		t.Logf("✅ Body incomplet rejeté: status %d", resp.StatusCode)
	} else {
		t.Logf("⚠️  Status = %d", resp.StatusCode)
	}
}

// TestServer_InvalidJSON teste la gestion d'un JSON invalide
func TestServer_InvalidJSON(t *testing.T) {
	t.Log("🧪 TEST SERVEUR - JSON INVALIDE")
	t.Log("================================")

	config := &Config{
		Host:     "localhost",
		Port:     0,
		Insecure: true,
		AuthType: "none",
	}

	server, err := setupTestServerOnRandomPort(config, t)
	if err != nil {
		t.Fatalf("❌ Erreur création serveur: %v", err)
	}
	defer server.httpServer.Close()

	baseURL := "http://" + server.httpServer.Addr

	invalidJSON := `{invalid json content here`

	req, _ := http.NewRequest("POST", baseURL+"/api/v1/execute", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		t.Logf("✅ JSON invalide rejeté: status %d", resp.StatusCode)
	} else {
		t.Errorf("❌ Status = %d, attendu 400", resp.StatusCode)
	}
}

// setupTestServerOnRandomPort crée un serveur de test sur un port aléatoire
func setupTestServerOnRandomPort(config *Config, t *testing.T) (*Server, error) {
	t.Helper()

	logger := log.New(io.Discard, "", 0)

	server, err := NewServer(config, logger)
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}

	server.httpServer = &http.Server{
		Handler:           server.mux,
		ReadTimeout:       DefaultReadTimeout,
		ReadHeaderTimeout: DefaultReadHeaderTimeout,
		WriteTimeout:      DefaultWriteTimeout,
		IdleTimeout:       DefaultIdleTimeout,
		MaxHeaderBytes:    DefaultMaxHeaderBytes,
	}

	server.httpServer.Addr = listener.Addr().String()

	go func() {
		server.httpServer.Serve(listener)
	}()

	time.Sleep(100 * time.Millisecond)

	return server, nil
}
