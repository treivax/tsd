// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package e2e

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"testing"
	"time"
)

const (
	// Timeouts
	ServerStartTimeout = 10 * time.Second
	ServerReadyWait    = 500 * time.Millisecond
	RequestTimeout     = 15 * time.Second
	ShutdownTimeout    = 5 * time.Second

	// Configuration JWT
	JWTDurationMinutes = 15
	JWTSecret          = "test-jwt-secret-for-e2e-tests-minimum-32-characters-long"

	// Configuration TLS
	CertValidityDays = 1
	CertOrganization = "TSD E2E Tests"

	// Valeurs attendues pour le programme de test
	ExpectedFactsCount       = 2
	ExpectedActivationsCount = 1
	ExpectedAdultActionName  = "adult"
	ExpectedAdultArgument    = "p1"
)

const (
	// SimpleTSDProgram est un programme TSD simple pour les tests E2E
	SimpleTSDProgram = `type Person(id: string, age: number)

action adult(id: string)

rule adult_check : {p: Person} / p.age >= 18 ==> adult(p.id)

Person(id:p1, age:25)
Person(id:p2, age:15)`

	// InvalidTSDProgram est un programme invalide pour tester la gestion d'erreurs
	InvalidTSDProgram = `rule broken_rule : {x: Unknown} / x > 0 ==> action()`
)

// testContext contient tout le contexte nécessaire pour un test E2E
type testContext struct {
	t          *testing.T
	tempDir    string
	serverURL  string
	serverPort int
	apiKey     string
	jwtToken   string
	serverCtx  context.Context
	cancelFunc context.CancelFunc
	caCertPEM  []byte
	certPath   string
	keyPath    string
}

// createHTTPClient crée un client HTTP configuré pour les tests
func (ctx *testContext) createHTTPClient() *http.Client {
	ctx.t.Helper()

	// Configuration TLS si HTTPS
	var tlsConfig *tls.Config
	if ctx.caCertPEM != nil {
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(ctx.caCertPEM) {
			ctx.t.Fatal("❌ Erreur ajout CA au pool")
		}
		tlsConfig = &tls.Config{
			RootCAs: caCertPool,
		}
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: RequestTimeout,
	}
}

// cleanup nettoie les ressources du test
func (ctx *testContext) cleanup() {
	ctx.t.Helper()

	if ctx.cancelFunc != nil {
		ctx.t.Log("🛑 Arrêt du serveur...")
		ctx.cancelFunc()
		time.Sleep(ShutdownTimeout)
	}

	ctx.t.Log("✅ Cleanup terminé")
}
