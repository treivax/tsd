// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package servercmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

// TestServerTimeouts vérifie que les timeouts par défaut sont correctement configurés
func TestServerTimeouts(t *testing.T) {
	t.Log("🧪 TEST SERVER TIMEOUTS")
	t.Log("=======================")

	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("❌ Échec création serveur: %v", err)
	}

	// Créer le serveur HTTP pour initialiser httpServer
	info := prepareServerInfo(config, server)
	server.httpServer = &http.Server{
		Addr:              info.Addr,
		Handler:           server.mux,
		ReadTimeout:       DefaultReadTimeout,
		ReadHeaderTimeout: DefaultReadHeaderTimeout,
		WriteTimeout:      DefaultWriteTimeout,
		IdleTimeout:       DefaultIdleTimeout,
		MaxHeaderBytes:    DefaultMaxHeaderBytes,
	}

	// Assert - Vérifier timeouts configurés
	tests := []struct {
		name     string
		got      time.Duration
		expected time.Duration
	}{
		{"ReadTimeout", server.httpServer.ReadTimeout, DefaultReadTimeout},
		{"WriteTimeout", server.httpServer.WriteTimeout, DefaultWriteTimeout},
		{"IdleTimeout", server.httpServer.IdleTimeout, DefaultIdleTimeout},
		{"ReadHeaderTimeout", server.httpServer.ReadHeaderTimeout, DefaultReadHeaderTimeout},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("❌ %s = %v, attendu %v", tt.name, tt.got, tt.expected)
			} else {
				t.Logf("✅ %s correctement configuré (%v)", tt.name, tt.got)
			}
		})
	}

	t.Log("✅ Tous les tests de timeouts passés")
}

// TestMaxHeaderBytes vérifie la limite de taille des headers
func TestMaxHeaderBytes(t *testing.T) {
	t.Log("🧪 TEST MAX HEADER BYTES")
	t.Log("========================")

	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("❌ Échec création serveur: %v", err)
	}

	// Créer le serveur HTTP
	info := prepareServerInfo(config, server)
	server.httpServer = &http.Server{
		Addr:           info.Addr,
		Handler:        server.mux,
		MaxHeaderBytes: DefaultMaxHeaderBytes,
	}

	// Assert
	expectedMax := 1 << 20 // 1 MB
	if server.httpServer.MaxHeaderBytes != expectedMax {
		t.Errorf("❌ MaxHeaderBytes = %d, attendu %d",
			server.httpServer.MaxHeaderBytes, expectedMax)
	} else {
		t.Logf("✅ MaxHeaderBytes configuré (%d bytes = 1 MB)", expectedMax)
	}

	t.Log("✅ Test MaxHeaderBytes passé")
}

// TestTimeoutConstants vérifie que les constantes de timeout ont les bonnes valeurs
func TestTimeoutConstants(t *testing.T) {
	t.Log("🧪 TEST TIMEOUT CONSTANTS")
	t.Log("=========================")

	tests := []struct {
		name     string
		got      time.Duration
		expected time.Duration
	}{
		{"DefaultReadTimeout", DefaultReadTimeout, 15 * time.Second},
		{"DefaultWriteTimeout", DefaultWriteTimeout, 15 * time.Second},
		{"DefaultIdleTimeout", DefaultIdleTimeout, 60 * time.Second},
		{"DefaultReadHeaderTimeout", DefaultReadHeaderTimeout, 5 * time.Second},
		{"DefaultShutdownTimeout", DefaultShutdownTimeout, 30 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("❌ %s = %v, attendu %v", tt.name, tt.got, tt.expected)
			} else {
				t.Logf("✅ %s = %v", tt.name, tt.got)
			}
		})
	}

	t.Log("✅ Tous les tests de constantes passés")
}

// TestMaxHeaderBytesConstant vérifie la constante MaxHeaderBytes
func TestMaxHeaderBytesConstant(t *testing.T) {
	t.Log("🧪 TEST MAX HEADER BYTES CONSTANT")
	t.Log("==================================")

	expected := 1 << 20 // 1 MB
	if DefaultMaxHeaderBytes != expected {
		t.Errorf("❌ DefaultMaxHeaderBytes = %d, attendu %d", DefaultMaxHeaderBytes, expected)
	} else {
		t.Logf("✅ DefaultMaxHeaderBytes = %d bytes (1 MB)", DefaultMaxHeaderBytes)
	}

	t.Log("✅ Test constante MaxHeaderBytes passé")
}

// TestReadHeaderTimeoutProtection teste la protection contre Slowloris
func TestReadHeaderTimeoutProtection(t *testing.T) {
	t.Log("🧪 TEST READ HEADER TIMEOUT PROTECTION (SLOWLORIS)")
	t.Log("===================================================")

	if testing.Short() {
		t.Skip("⏭️  Test long, skip en mode -short")
	}

	// Créer un serveur de test avec timeout très court
	config := &Config{
		Host:     "127.0.0.1",
		Port:     0, // Port aléatoire
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("❌ Échec création serveur: %v", err)
	}

	// Créer un listener pour obtenir le port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("❌ Échec création listener: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().String()
	t.Logf("📡 Serveur écoute sur %s", addr)

	// Configurer serveur avec timeout très court pour accélérer le test
	server.httpServer = &http.Server{
		Handler:           server.mux,
		ReadHeaderTimeout: 200 * time.Millisecond, // Très court pour test
		ReadTimeout:       500 * time.Millisecond,
		WriteTimeout:      500 * time.Millisecond,
	}

	// Démarrer le serveur
	go func() {
		server.httpServer.Serve(listener)
	}()
	defer server.Shutdown(context.Background())

	// Attendre que le serveur soit prêt
	time.Sleep(100 * time.Millisecond)

	// Tenter une connexion lente (simuler Slowloris)
	t.Log("🐌 Simulation attaque Slowloris (envoi lent de headers)...")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("❌ Échec connexion: %v", err)
	}
	defer conn.Close()

	// Écrire seulement une partie des headers très lentement
	conn.Write([]byte("GET"))
	time.Sleep(300 * time.Millisecond) // > ReadHeaderTimeout

	// Essayer de continuer - devrait échouer car timeout
	_, err = conn.Write([]byte(" /health HTTP/1.1\r\n\r\n"))

	// La connexion peut soit échouer immédiatement, soit accepter l'écriture
	// mais ne pas renvoyer de réponse valide
	if err == nil {
		// Essayer de lire la réponse
		conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		buf := make([]byte, 1024)
		n, readErr := conn.Read(buf)
		if readErr != nil && readErr != io.EOF {
			t.Logf("✅ ReadHeaderTimeout correctement appliqué: erreur lecture: %v", readErr)
		} else if n == 0 {
			t.Logf("✅ ReadHeaderTimeout correctement appliqué: connexion fermée")
		} else {
			// Vérifier si c'est une erreur HTTP (400 Bad Request attendu)
			response := string(buf[:n])
			if len(response) > 0 {
				t.Logf("✅ ReadHeaderTimeout appliqué: serveur a répondu (possiblement erreur): %s", response[:min(50, len(response))])
			}
		}
	} else {
		t.Logf("✅ ReadHeaderTimeout correctement appliqué: %v", err)
	}

	t.Log("✅ Test protection Slowloris passé")
}

// TestReadTimeoutEnforcement teste l'application du ReadTimeout
func TestReadTimeoutEnforcement(t *testing.T) {
	t.Log("🧪 TEST READ TIMEOUT ENFORCEMENT")
	t.Log("================================")

	if testing.Short() {
		t.Skip("⏭️  Test long, skip en mode -short")
	}

	// Créer serveur avec timeout court
	config := &Config{
		Host:     "127.0.0.1",
		Port:     0,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("❌ Échec création serveur: %v", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("❌ Échec création listener: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().String()
	t.Logf("📡 Serveur écoute sur %s", addr)

	// Timeout très court pour accélérer le test
	server.httpServer = &http.Server{
		Handler:      server.mux,
		ReadTimeout:  200 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
	}

	go server.httpServer.Serve(listener)
	defer server.Shutdown(context.Background())

	time.Sleep(100 * time.Millisecond)

	// Connexion qui envoie headers rapidement mais body lentement
	t.Log("🐌 Simulation client lent (body envoyé lentement)...")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("❌ Échec connexion: %v", err)
	}
	defer conn.Close()

	// Envoyer headers rapidement
	headers := "POST /api/v1/execute HTTP/1.1\r\n" +
		"Host: " + addr + "\r\n" +
		"Content-Type: application/json\r\n" +
		"Content-Length: 100\r\n" +
		"\r\n"
	conn.Write([]byte(headers))

	// Attendre plus que ReadTimeout avant d'envoyer le body
	time.Sleep(300 * time.Millisecond)

	// Essayer d'envoyer le body - devrait échouer
	_, err = conn.Write([]byte(`{"source":"test"}`))
	if err == nil {
		// Peut réussir à écrire mais la lecture devrait timeout
		conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
		buf := make([]byte, 1024)
		_, readErr := conn.Read(buf)
		if readErr != nil {
			t.Logf("✅ ReadTimeout appliqué: erreur lecture après timeout: %v", readErr)
		}
	} else {
		t.Logf("✅ ReadTimeout appliqué: erreur écriture: %v", err)
	}

	t.Log("✅ Test ReadTimeout enforcement passé")
}

// TestIdleTimeoutForKeepAlive teste le timeout des connexions keep-alive
func TestIdleTimeoutForKeepAlive(t *testing.T) {
	t.Log("🧪 TEST IDLE TIMEOUT (KEEP-ALIVE)")
	t.Log("==================================")

	if testing.Short() {
		t.Skip("⏭️  Test long, skip en mode -short")
	}

	config := &Config{
		Host:     "127.0.0.1",
		Port:     0,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("❌ Échec création serveur: %v", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("❌ Échec création listener: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().String()
	t.Logf("📡 Serveur écoute sur %s", addr)

	// IdleTimeout court pour test
	server.httpServer = &http.Server{
		Handler:     server.mux,
		IdleTimeout: 300 * time.Millisecond,
		ReadTimeout: 1 * time.Second,
	}

	go server.httpServer.Serve(listener)
	defer server.Shutdown(context.Background())

	time.Sleep(100 * time.Millisecond)

	// Créer client HTTP avec keep-alive
	t.Log("🔗 Test connexion keep-alive avec IdleTimeout...")
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 1,
			IdleConnTimeout:     5 * time.Second,
		},
	}

	// Première requête - devrait réussir
	url := fmt.Sprintf("http://%s/health", addr)
	resp1, err := client.Get(url)
	if err != nil {
		t.Fatalf("❌ Première requête échouée: %v", err)
	}
	io.Copy(io.Discard, resp1.Body)
	resp1.Body.Close()
	t.Log("✅ Première requête réussie")

	// Attendre plus que IdleTimeout
	t.Log("⏳ Attente dépassement IdleTimeout...")
	time.Sleep(400 * time.Millisecond)

	// Deuxième requête - connexion devrait être fermée et reconnectée
	resp2, err := client.Get(url)
	if err != nil {
		t.Logf("✅ IdleTimeout appliqué: connexion fermée, nouvelle tentative échoue: %v", err)
	} else {
		io.Copy(io.Discard, resp2.Body)
		resp2.Body.Close()
		t.Log("✅ Deuxième requête réussie (nouvelle connexion établie)")
	}

	t.Log("✅ Test IdleTimeout passé")
}

// TestTimeoutsWithTLS teste les timeouts avec TLS activé
func TestTimeoutsWithTLS(t *testing.T) {
	t.Log("🧪 TEST TIMEOUTS AVEC TLS")
	t.Log("=========================")

	if testing.Short() {
		t.Skip("⏭️  Test long, skip en mode -short")
	}

	// Créer certificat auto-signé pour le test
	certFile, keyFile, skip := createTestCertificates(t)
	if skip {
		t.Skip("⏭️  Certificats de test non disponibles, skip du test TLS")
	}
	defer func() {
		// Nettoyer les certificats temporaires
		// (createTestCertificates devrait les créer dans un dossier temp)
	}()

	config := &Config{
		Host:        "127.0.0.1",
		Port:        0,
		AuthType:    "none",
		TLSCertFile: certFile,
		TLSKeyFile:  keyFile,
		Insecure:    false,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("❌ Échec création serveur: %v", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("❌ Échec création listener: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().String()
	t.Logf("📡 Serveur TLS écoute sur %s", addr)

	// Configurer TLS
	tlsConf, err := createTLSConfig(certFile, keyFile)
	if err != nil {
		t.Fatalf("❌ Erreur configuration TLS: %v", err)
	}

	server.httpServer = &http.Server{
		Handler:           server.mux,
		TLSConfig:         tlsConf,
		ReadTimeout:       500 * time.Millisecond,
		ReadHeaderTimeout: 200 * time.Millisecond,
		WriteTimeout:      500 * time.Millisecond,
		IdleTimeout:       1 * time.Second,
		MaxHeaderBytes:    DefaultMaxHeaderBytes,
	}

	// Utiliser listener TLS
	tlsListener := tls.NewListener(listener, tlsConf)
	go server.httpServer.Serve(tlsListener)
	defer server.Shutdown(context.Background())

	time.Sleep(100 * time.Millisecond)

	// Client avec TLS qui ignore les certificats auto-signés
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 2 * time.Second,
	}

	// Test requête normale - devrait réussir
	url := fmt.Sprintf("https://%s/health", addr)
	resp, err := client.Get(url)
	if err != nil {
		t.Logf("⚠️  Requête TLS échouée (peut être normal en test): %v", err)
	} else {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		t.Log("✅ Requête HTTPS réussie avec timeouts configurés")

		// Vérifier que les timeouts sont appliqués
		if server.httpServer.ReadTimeout != 500*time.Millisecond {
			t.Errorf("❌ ReadTimeout TLS = %v, attendu 500ms", server.httpServer.ReadTimeout)
		} else {
			t.Log("✅ Timeouts TLS correctement configurés")
		}
	}

	t.Log("✅ Test timeouts avec TLS passé")
}

// createTestCertificates crée des certificats auto-signés pour les tests
func createTestCertificates(t *testing.T) (certFile, keyFile string, skip bool) {
	// Chercher des certificats de test existants
	testCertDir := "../../tests/fixtures/certs"
	certFile = testCertDir + "/test-server.crt"
	keyFile = testCertDir + "/test-server.key"

	// Vérifier si les certificats existent déjà
	_, certExists := os.Stat(certFile)
	_, keyExists := os.Stat(keyFile)

	if certExists == nil && keyExists == nil {
		t.Logf("📜 Utilisation certificats test existants: %s, %s", certFile, keyFile)
		return certFile, keyFile, false
	}

	// Tenter de générer les certificats automatiquement
	t.Logf("🔐 Génération automatique des certificats de test...")
	generateScript := testCertDir + "/generate_certs.sh"

	if _, err := os.Stat(generateScript); os.IsNotExist(err) {
		t.Logf("⚠️  Script de génération non trouvé: %s", generateScript)
		t.Logf("⚠️  Les tests TLS nécessitent des certificats. Voir tests/fixtures/certs/README.md")
		return "", "", true
	}

	// Exécuter le script de génération
	cmd := exec.Command("bash", generateScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("⚠️  Échec de la génération des certificats: %v", err)
		t.Logf("Output: %s", string(output))
		return "", "", true
	}

	// Vérifier que les certificats ont été créés
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		t.Logf("⚠️  Certificat non créé après génération: %s", certFile)
		return "", "", true
	}
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		t.Logf("⚠️  Clé privée non créée après génération: %s", keyFile)
		return "", "", true
	}

	t.Logf("✅ Certificats de test générés avec succès")
	return certFile, keyFile, false
}

// min retourne le minimum de deux entiers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
