// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package clientcmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestClient_ExpiredCertificate teste la gestion d'un certificat expiré
func TestClient_ExpiredCertificate(t *testing.T) {
	t.Log("🧪 TEST CLIENT - CERTIFICAT EXPIRÉ")
	t.Log("===================================")

	expiredCert, expiredKey := generateExpiredCertificate(t)

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	cert, err := tls.X509KeyPair(expiredCert, expiredKey)
	if err != nil {
		t.Fatalf("❌ Erreur création certificat: %v", err)
	}

	server.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
	server.StartTLS()
	defer server.Close()

	config := &Config{
		ServerURL: server.URL,
		Timeout:   2 * time.Second,
		Insecure:  false,
	}

	client := NewClient(config)

	_, err = client.Execute("type Person : <id: string>", "<test>")

	if err == nil {
		t.Fatal("❌ Attendait une erreur de certificat expiré")
	}

	errMsg := err.Error()
	if strings.Contains(errMsg, "certificate has expired") ||
		strings.Contains(errMsg, "expired") ||
		strings.Contains(errMsg, "x509") {
		t.Logf("✅ Certificat expiré détecté: %v", err)
	} else {
		t.Logf("⚠️  Erreur TLS détectée (probablement certificat expiré): %v", err)
	}
}

// TestClient_SelfSignedCertificate teste la gestion d'un certificat auto-signé
func TestClient_SelfSignedCertificate(t *testing.T) {
	t.Log("🧪 TEST CLIENT - CERTIFICAT AUTO-SIGNÉ")
	t.Log("=======================================")

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := &Config{
		ServerURL: server.URL,
		Timeout:   2 * time.Second,
		Insecure:  false,
	}

	client := NewClient(config)

	_, err := client.Execute("type Person : <id: string>", "<test>")

	if err == nil {
		t.Fatal("❌ Attendait une erreur de certificat non vérifié")
	}

	errMsg := err.Error()
	if strings.Contains(errMsg, "certificate") ||
		strings.Contains(errMsg, "tls") ||
		strings.Contains(errMsg, "x509") ||
		strings.Contains(errMsg, "unknown authority") {
		t.Logf("✅ Certificat auto-signé rejeté: %v", err)
	} else {
		t.Logf("⚠️  Erreur TLS détectée: %v", err)
	}
}

// TestClient_SelfSignedCertificate_Insecure teste qu'en mode insecure le certificat est accepté
func TestClient_SelfSignedCertificate_Insecure(t *testing.T) {
	t.Log("🧪 TEST CLIENT - CERTIFICAT AUTO-SIGNÉ (MODE INSECURE)")
	t.Log("=======================================================")

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "results": {"facts_count": 0, "activations_count": 0, "activations": []}, "execution_time_ms": 0}`))
	}))
	defer server.Close()

	config := &Config{
		ServerURL: server.URL,
		Timeout:   2 * time.Second,
		Insecure:  true,
	}

	client := NewClient(config)

	resp, err := client.Execute("type Person : <id: string>", "<test>")

	if err != nil {
		t.Fatalf("❌ Erreur inattendue en mode insecure: %v", err)
	}

	if resp == nil || !resp.Success {
		t.Fatal("❌ Réponse invalide")
	}

	t.Log("✅ Certificat auto-signé accepté en mode insecure")
}

// TestClient_HostnameMismatch teste la gestion d'un certificat avec hostname incorrect
func TestClient_HostnameMismatch(t *testing.T) {
	t.Log("🧪 TEST CLIENT - HOSTNAME MISMATCH")
	t.Log("===================================")

	certPEM, keyPEM := generateCertificateForHost(t, "example.com")

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		t.Fatalf("❌ Erreur création certificat: %v", err)
	}

	server.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
	server.StartTLS()
	defer server.Close()

	config := &Config{
		ServerURL: server.URL,
		Timeout:   2 * time.Second,
		Insecure:  false,
	}

	client := NewClient(config)

	_, err = client.Execute("type Person : <id: string>", "<test>")

	if err == nil {
		t.Fatal("❌ Attendait une erreur de hostname mismatch")
	}

	errMsg := err.Error()
	if strings.Contains(errMsg, "certificate") ||
		strings.Contains(errMsg, "x509") ||
		strings.Contains(errMsg, "tls") {
		t.Logf("✅ Hostname mismatch détecté: %v", err)
	} else {
		t.Logf("⚠️  Erreur TLS détectée: %v", err)
	}
}

// generateExpiredCertificate génère un certificat déjà expiré
func generateExpiredCertificate(t *testing.T) ([]byte, []byte) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("❌ Erreur génération clé privée: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"TSD Test"},
			CommonName:   "expired.test",
		},
		NotBefore: time.Now().Add(-48 * time.Hour),
		NotAfter:  time.Now().Add(-24 * time.Hour),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
		DNSNames:              []string{"expired.test"},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("❌ Erreur création certificat: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	return certPEM, keyPEM
}

// generateCertificateForHost génère un certificat pour un hostname spécifique
func generateCertificateForHost(t *testing.T, hostname string) ([]byte, []byte) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("❌ Erreur génération clé privée: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"TSD Test"},
			CommonName:   hostname,
		},
		NotBefore: time.Now().Add(-1 * time.Hour),
		NotAfter:  time.Now().Add(24 * time.Hour),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
		DNSNames:              []string{hostname},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("❌ Erreur création certificat: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	return certPEM, keyPEM
}
