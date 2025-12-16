// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package e2e

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/treivax/tsd/auth"
	"github.com/treivax/tsd/internal/servercmd"
)

// TestClientServerRoundtrip_Complete teste le scénario E2E complet avec HTTPS et JWT
func TestClientServerRoundtrip_Complete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Log("🧪 TEST E2E COMPLET : Client-Server HTTPS + JWT")
	t.Log("===============================================")

	ctx := setupCompleteTestContext(t)
	defer ctx.cleanup()

	t.Run("Scénario_Complet_HTTPS_JWT", func(t *testing.T) {
		testCompleteRoundtrip(ctx)
	})

	t.Run("Authentification_Invalide", func(t *testing.T) {
		testInvalidAuthentication(ctx)
	})

	t.Run("Programme_TSD_Invalide", func(t *testing.T) {
		testInvalidProgram(ctx)
	})

	t.Run("Requete_Sans_Token", func(t *testing.T) {
		testUnauthorizedRequest(ctx)
	})
}

// TestClientServerRoundtrip_HTTP_NoAuth teste le scénario simple HTTP sans auth (legacy)
func TestClientServerRoundtrip_HTTP_NoAuth(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Log("🧪 TEST E2E SIMPLE : Client-Server HTTP sans auth")
	t.Log("=================================================")

	ctx := setupHTTPTestContext(t)
	defer ctx.cleanup()

	testSimpleHTTPRoundtrip(ctx)
}

// setupCompleteTestContext configure un environnement de test E2E complet (HTTPS + JWT)
func setupCompleteTestContext(t *testing.T) *testContext {
	t.Helper()

	ctx := &testContext{t: t}

	// 1. Créer répertoire temporaire
	tempDir := t.TempDir()
	ctx.tempDir = tempDir
	t.Logf("📁 Répertoire temporaire: %s", tempDir)

	// 2. Générer certificats TLS
	t.Log("🔐 Génération certificats TLS...")
	ctx.generateTLSCertificates()

	// 3. Générer clé API
	t.Log("🔑 Génération clé API...")
	ctx.generateAPIKey()

	// 4. Générer JWT
	t.Log("🎫 Génération JWT...")
	ctx.generateJWT()

	// 5. Démarrer serveur HTTPS avec JWT
	t.Log("🚀 Démarrage serveur HTTPS + JWT...")
	ctx.startHTTPSServerWithJWT()

	// 6. Attendre serveur prêt
	t.Log("⏳ Attente serveur prêt...")
	ctx.waitForServerReady()

	t.Log("✅ Setup E2E complet terminé")
	return ctx
}

// setupHTTPTestContext configure un environnement de test E2E simple (HTTP sans auth)
func setupHTTPTestContext(t *testing.T) *testContext {
	t.Helper()

	ctx := &testContext{t: t}

	// Créer répertoire temporaire
	tempDir := t.TempDir()
	ctx.tempDir = tempDir
	t.Logf("📁 Répertoire temporaire: %s", tempDir)

	// Démarrer serveur HTTP simple
	t.Log("🚀 Démarrage serveur HTTP simple...")
	ctx.startHTTPServer()

	// Attendre serveur prêt
	t.Log("⏳ Attente serveur prêt...")
	ctx.waitForServerReady()

	t.Log("✅ Setup HTTP simple terminé")
	return ctx
}

// generateTLSCertificates génère les certificats TLS pour le test
func (ctx *testContext) generateTLSCertificates() {
	ctx.t.Helper()

	// Générer CA auto-signé
	caCertPEM, caKeyPEM, err := generateSelfSignedCA(CertOrganization, CertValidityDays)
	if err != nil {
		ctx.t.Fatalf("❌ Erreur génération CA: %v", err)
	}
	ctx.caCertPEM = caCertPEM

	// Générer certificat serveur
	serverCertPEM, serverKeyPEM, err := generateServerCert(
		caCertPEM, caKeyPEM,
		[]string{"localhost", "127.0.0.1"},
		CertOrganization,
		CertValidityDays,
	)
	if err != nil {
		ctx.t.Fatalf("❌ Erreur génération certificat serveur: %v", err)
	}

	// Sauvegarder certificats dans fichiers
	ctx.certPath = filepath.Join(ctx.tempDir, "server.crt")
	ctx.keyPath = filepath.Join(ctx.tempDir, "server.key")
	caPath := filepath.Join(ctx.tempDir, "ca.crt")

	if err := os.WriteFile(ctx.certPath, serverCertPEM, 0600); err != nil {
		ctx.t.Fatalf("❌ Erreur écriture certificat: %v", err)
	}
	if err := os.WriteFile(ctx.keyPath, serverKeyPEM, 0600); err != nil {
		ctx.t.Fatalf("❌ Erreur écriture clé: %v", err)
	}
	if err := os.WriteFile(caPath, caCertPEM, 0600); err != nil {
		ctx.t.Fatalf("❌ Erreur écriture CA: %v", err)
	}

	ctx.t.Logf("   ✅ Certificats générés: %s", ctx.certPath)
}

// generateAPIKey génère une clé API pour le test
func (ctx *testContext) generateAPIKey() {
	ctx.t.Helper()

	apiKey, err := auth.GenerateAuthKey()
	if err != nil {
		ctx.t.Fatalf("❌ Erreur génération clé API: %v", err)
	}
	ctx.apiKey = apiKey
	ctx.t.Logf("   ✅ Clé API: %s...", apiKey[:16])
}

// generateJWT génère un JWT pour le test
func (ctx *testContext) generateJWT() {
	ctx.t.Helper()

	// Créer un Auth Manager pour générer le JWT
	authConfig := &auth.Config{
		Type:          auth.AuthTypeJWT,
		JWTSecret:     JWTSecret,
		JWTExpiration: time.Duration(JWTDurationMinutes) * time.Minute,
		JWTIssuer:     "tsd-e2e-test",
	}

	manager, err := auth.NewManager(authConfig)
	if err != nil {
		ctx.t.Fatalf("❌ Erreur création Auth Manager: %v", err)
	}

	token, err := manager.GenerateJWT("test-user", []string{"admin"})
	if err != nil {
		ctx.t.Fatalf("❌ Erreur génération JWT: %v", err)
	}

	ctx.jwtToken = token
	ctx.t.Logf("   ✅ JWT: %s...", token[:32])
}

// findAvailablePort trouve un port disponible pour le serveur
func findAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	return port, nil
}

// startHTTPSServerWithJWT démarre un serveur HTTPS avec authentification JWT
func (ctx *testContext) startHTTPSServerWithJWT() {
	ctx.t.Helper()

	port, err := findAvailablePort()
	if err != nil {
		ctx.t.Fatalf("❌ Erreur allocation port: %v", err)
	}
	ctx.serverPort = port
	ctx.serverURL = fmt.Sprintf("https://127.0.0.1:%d", port)

	// Créer contexte avec annulation pour le serveur
	serverCtx, cancel := context.WithCancel(context.Background())
	ctx.serverCtx = serverCtx
	ctx.cancelFunc = cancel

	// Démarrer serveur dans goroutine
	go func() {
		args := []string{
			"-host", "127.0.0.1",
			"-port", fmt.Sprintf("%d", port),
			"-auth", auth.AuthTypeJWT,
			"-jwt-secret", JWTSecret,
			"-jwt-issuer", "tsd-e2e-test",
			"-tls-cert", ctx.certPath,
			"-tls-key", ctx.keyPath,
		}

		var stdout, stderr bytes.Buffer
		servercmd.Run(args, os.Stdin, &stdout, &stderr)
	}()

	ctx.t.Logf("   ✅ Serveur HTTPS démarré sur %s", ctx.serverURL)
}

// startHTTPServer démarre un serveur HTTP simple sans auth
func (ctx *testContext) startHTTPServer() {
	ctx.t.Helper()

	port, err := findAvailablePort()
	if err != nil {
		ctx.t.Fatalf("❌ Erreur allocation port: %v", err)
	}
	ctx.serverPort = port
	ctx.serverURL = fmt.Sprintf("http://127.0.0.1:%d", port)

	// Créer contexte avec annulation pour le serveur
	serverCtx, cancel := context.WithCancel(context.Background())
	ctx.serverCtx = serverCtx
	ctx.cancelFunc = cancel

	// Démarrer serveur dans goroutine
	go func() {
		args := []string{
			"-host", "127.0.0.1",
			"-port", fmt.Sprintf("%d", port),
			"-auth", "none",
			"-insecure",
		}

		var stdout, stderr bytes.Buffer
		servercmd.Run(args, os.Stdin, &stdout, &stderr)
	}()

	ctx.t.Logf("   ✅ Serveur HTTP démarré sur %s", ctx.serverURL)
}

// waitForServerReady attend que le serveur soit prêt
func (ctx *testContext) waitForServerReady() {
	ctx.t.Helper()

	// Créer client HTTP pour health check
	httpClient := ctx.createHTTPClient()

	healthURL := ctx.serverURL + "/health"
	deadline := time.Now().Add(ServerStartTimeout)

	for time.Now().Before(deadline) {
		resp, err := httpClient.Get(healthURL)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				ctx.t.Log("   ✅ Serveur prêt")
				return
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	ctx.t.Fatal("❌ Timeout attente serveur")
}
