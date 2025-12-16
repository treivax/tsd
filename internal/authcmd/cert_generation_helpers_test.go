// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package authcmd

import (
	"crypto/x509"
	"testing"
	"time"
)

// TestCreateCertificateTemplate_IsNotCA tests that generated certificate is not a CA
func TestCreateCertificateTemplate_IsNotCA(t *testing.T) {
	config := &certConfig{
		hosts:     []string{"localhost", "127.0.0.1"},
		validDays: 365,
		org:       "Test Organization",
	}

	template, err := createCertificateTemplate(config)
	if err != nil {
		t.Fatalf("❌ createCertificateTemplate() error = %v", err)
	}

	if template.IsCA {
		t.Error("❌ SECURITY: Certificate template has IsCA=true, should be false for server/client cert")
	} else {
		t.Log("✅ Certificate correctly marked as IsCA=false (not a CA)")
	}

	// Verify BasicConstraintsValid is still set
	if !template.BasicConstraintsValid {
		t.Error("❌ BasicConstraintsValid should be true")
	}
}

// TestCreateCertificateTemplate_ValidityPeriod tests certificate validity period
func TestCreateCertificateTemplate_ValidityPeriod(t *testing.T) {
	tests := []struct {
		name      string
		validDays int
		wantDays  int
	}{
		{"1 year", 365, 365},
		{"2 years", 730, 730},
		{"90 days", 90, 90},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &certConfig{
				hosts:     []string{"localhost"},
				validDays: tt.validDays,
				org:       "Test Org",
			}

			template, err := createCertificateTemplate(config)
			if err != nil {
				t.Fatalf("createCertificateTemplate() error = %v", err)
			}

			duration := template.NotAfter.Sub(template.NotBefore)
			expectedDuration := time.Duration(tt.wantDays) * 24 * time.Hour

			// Allow 1 hour tolerance for test execution time
			tolerance := 1 * time.Hour
			if duration > expectedDuration+tolerance || duration < expectedDuration-tolerance {
				t.Errorf("❌ Validity period incorrect: got %v, want ~%v", duration, expectedDuration)
			} else {
				t.Logf("✅ Validity period correct: %v (expected ~%v)", duration, expectedDuration)
			}
		})
	}
}

// TestCreateCertificateTemplate_DNSAndIPAddresses tests SANs configuration
func TestCreateCertificateTemplate_DNSAndIPAddresses(t *testing.T) {
	config := &certConfig{
		hosts:     []string{"localhost", "example.com", "127.0.0.1", "192.168.1.1"},
		validDays: 365,
		org:       "Test Org",
	}

	template, err := createCertificateTemplate(config)
	if err != nil {
		t.Fatalf("createCertificateTemplate() error = %v", err)
	}

	// Should have 2 DNS names (localhost, example.com)
	if len(template.DNSNames) != 2 {
		t.Errorf("❌ Expected 2 DNS names, got %d", len(template.DNSNames))
	} else {
		t.Logf("✅ DNS names: %v", template.DNSNames)
	}

	// Should have 2 IP addresses (127.0.0.1, 192.168.1.1)
	if len(template.IPAddresses) != 2 {
		t.Errorf("❌ Expected 2 IP addresses, got %d", len(template.IPAddresses))
	} else {
		t.Logf("✅ IP addresses: %v", template.IPAddresses)
	}
}

// TestCreateCertificateTemplate_KeyUsage tests that KeyUsage is correctly set
func TestCreateCertificateTemplate_KeyUsage(t *testing.T) {
	config := &certConfig{
		hosts:     []string{"localhost"},
		validDays: 365,
		org:       "Test Org",
	}

	template, err := createCertificateTemplate(config)
	if err != nil {
		t.Fatalf("createCertificateTemplate() error = %v", err)
	}

	// Verify key usage includes required values
	if template.KeyUsage == 0 {
		t.Error("❌ KeyUsage is not set")
	}

	// Verify ExtKeyUsage includes ServerAuth and ClientAuth
	hasServerAuth := false
	hasClientAuth := false
	for _, usage := range template.ExtKeyUsage {
		if usage == 1 { // x509.ExtKeyUsageServerAuth
			hasServerAuth = true
		}
		if usage == 2 { // x509.ExtKeyUsageClientAuth
			hasClientAuth = true
		}
	}

	if !hasServerAuth {
		t.Error("❌ ExtKeyUsage missing ServerAuth")
	} else {
		t.Log("✅ ExtKeyUsage includes ServerAuth")
	}

	if !hasClientAuth {
		t.Error("❌ ExtKeyUsage missing ClientAuth")
	} else {
		t.Log("✅ ExtKeyUsage includes ClientAuth")
	}
}

// TestCreateCertificateTemplate_RFC5280Compliance tests RFC 5280 compliance
func TestCreateCertificateTemplate_RFC5280Compliance(t *testing.T) {
	t.Log("🧪 TEST CONFORMITÉ RFC 5280")
	t.Log("===========================")

	config := &certConfig{
		hosts:     []string{"localhost"},
		validDays: 365,
		org:       "Test Org",
	}

	template, err := createCertificateTemplate(config)
	if err != nil {
		t.Fatalf("❌ createCertificateTemplate() error = %v", err)
	}

	// RFC 5280 Section 4.2.1.9: Basic Constraints
	// For end-entity certificates (non-CA), cA MUST be FALSE
	if template.IsCA {
		t.Error("❌ RFC 5280 VIOLATION: IsCA must be false for end-entity certificates")
		t.Error("   Section 4.2.1.9: The cA boolean indicates whether the certified public key")
		t.Error("   may be used to verify certificate signatures. For end-entity certs, cA=FALSE")
	} else {
		t.Log("✅ RFC 5280 compliant: IsCA=false for end-entity certificate")
	}

	// BasicConstraintsValid must be true to indicate the extension is present
	if !template.BasicConstraintsValid {
		t.Error("❌ RFC 5280: BasicConstraintsValid should be true")
	} else {
		t.Log("✅ BasicConstraintsValid=true (extension présente)")
	}

	// RFC 5280 Section 4.2.1.3: Key Usage
	// For server/client certs: digitalSignature, keyEncipherment
	expectedKeyUsage := x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
	if (template.KeyUsage & expectedKeyUsage) == 0 {
		t.Errorf("❌ RFC 5280: KeyUsage should include DigitalSignature and KeyEncipherment")
	} else {
		t.Log("✅ RFC 5280 compliant: KeyUsage includes required values")
	}

	// RFC 5280 Section 4.2.1.12: Extended Key Usage
	// For TLS server: id-kp-serverAuth (1.3.6.1.5.5.7.3.1)
	// For TLS client: id-kp-clientAuth (1.3.6.1.5.5.7.3.2)
	if len(template.ExtKeyUsage) == 0 {
		t.Error("❌ RFC 5280: ExtKeyUsage should be present for TLS certificates")
	} else {
		t.Logf("✅ RFC 5280 compliant: ExtKeyUsage present with %d values", len(template.ExtKeyUsage))
	}

	// Serial number should be positive and unique
	if template.SerialNumber == nil || template.SerialNumber.Sign() <= 0 {
		t.Error("❌ RFC 5280: Serial number must be positive")
	} else {
		t.Logf("✅ Serial number: %s", template.SerialNumber.String())
	}

	t.Log("")
	t.Log("📋 Résumé conformité RFC 5280:")
	t.Logf("   IsCA: %v (attendu: false)", template.IsCA)
	t.Logf("   BasicConstraintsValid: %v", template.BasicConstraintsValid)
	t.Logf("   KeyUsage: %d", template.KeyUsage)
	t.Logf("   ExtKeyUsage: %v", template.ExtKeyUsage)
}

// TestGeneratedCertificate_SecurityProperties tests actual generated certificate
func TestGeneratedCertificate_SecurityProperties(t *testing.T) {
	t.Log("🧪 TEST PROPRIÉTÉS SÉCURITÉ CERTIFICAT GÉNÉRÉ")
	t.Log("==============================================")

	config := &certConfig{
		hosts:     []string{"localhost"},
		validDays: 365,
		org:       "Test Org Security",
	}

	// Generate private key
	privateKey, err := generateECDSAPrivateKey()
	if err != nil {
		t.Fatalf("❌ Échec génération clé: %v", err)
	}

	// Create template
	template, err := createCertificateTemplate(config)
	if err != nil {
		t.Fatalf("❌ Échec création template: %v", err)
	}

	// Generate certificate
	certDER, err := createSelfSignedCertificate(template, privateKey)
	if err != nil {
		t.Fatalf("❌ Échec génération certificat: %v", err)
	}

	// Parse the generated certificate using x509
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("❌ Échec parsing certificat: %v", err)
	}

	// Verify IsCA is false (CRITICAL SECURITY CHECK)
	if cert.IsCA {
		t.Error("❌ SÉCURITÉ CRITIQUE: Le certificat généré a IsCA=true")
		t.Error("   IMPACT: Certificat pourrait être utilisé pour signer d'autres certificats")
		t.Error("   RISQUE: CWE-295 (Improper Certificate Validation)")
		t.Error("   SOLUTION: IsCA doit être false pour certificats serveur/client")
	} else {
		t.Log("✅ SÉCURITÉ: IsCA=false (certificat ne peut pas signer d'autres certificats)")
	}

	// Verify BasicConstraints
	if !cert.BasicConstraintsValid {
		t.Error("⚠️ BasicConstraintsValid devrait être true")
	} else {
		t.Log("✅ BasicConstraintsValid correctement défini")
	}

	// Verify the certificate cannot be used as CA
	if cert.IsCA {
		// Check if KeyCertSign is set (it should NOT be for non-CA)
		if (cert.KeyUsage & x509.KeyUsageCertSign) != 0 {
			t.Error("❌ SÉCURITÉ: KeyUsageCertSign ne devrait PAS être défini pour certificat non-CA")
		}
	} else {
		t.Log("✅ Certificat ne peut pas signer d'autres certificats (usage CA interdit)")
	}

	// Verify signature algorithm
	if cert.SignatureAlgorithm == 0 {
		t.Error("❌ Algorithme de signature non défini")
	} else {
		t.Logf("✅ Algorithme signature: %v", cert.SignatureAlgorithm)
	}

	// Verify public key algorithm
	if cert.PublicKeyAlgorithm == 0 {
		t.Error("❌ Algorithme de clé publique non défini")
	} else {
		t.Logf("✅ Algorithme clé publique: %v", cert.PublicKeyAlgorithm)
	}

	t.Log("")
	t.Log("📊 Propriétés de sécurité vérifiées:")
	t.Logf("   ✅ IsCA: %v (DOIT être false)", cert.IsCA)
	t.Logf("   ✅ BasicConstraintsValid: %v", cert.BasicConstraintsValid)
	t.Logf("   ✅ KeyUsage: %d", cert.KeyUsage)
	t.Logf("   ✅ ExtKeyUsage: %v", cert.ExtKeyUsage)
	t.Logf("   ✅ SignatureAlgorithm: %v", cert.SignatureAlgorithm)
	t.Logf("   ✅ PublicKeyAlgorithm: %v", cert.PublicKeyAlgorithm)
}

// TestCertificate_CannotSignOtherCerts tests that generated cert cannot sign others
func TestCertificate_CannotSignOtherCerts(t *testing.T) {
	t.Log("🧪 TEST CERTIFICAT NE PEUT PAS SIGNER D'AUTRES CERTIFICATS")
	t.Log("===========================================================")

	config := &certConfig{
		hosts:     []string{"localhost"},
		validDays: 365,
		org:       "Test Org",
	}

	// Generate first certificate (should NOT be usable as CA)
	privateKey1, err := generateECDSAPrivateKey()
	if err != nil {
		t.Fatalf("❌ Échec génération clé 1: %v", err)
	}
	template1, err := createCertificateTemplate(config)
	if err != nil {
		t.Fatalf("❌ Échec création template 1: %v", err)
	}
	certDER1, err := createSelfSignedCertificate(template1, privateKey1)
	if err != nil {
		t.Fatalf("❌ Échec génération certificat 1: %v", err)
	}
	cert1, err := x509.ParseCertificate(certDER1)
	if err != nil {
		t.Fatalf("❌ Échec parsing certificat 1: %v", err)
	}

	// Verify this certificate is NOT a CA
	if cert1.IsCA {
		t.Fatal("❌ Le certificat test est marqué comme CA, test invalide")
	}

	t.Logf("✅ Certificat 1 correctement marqué IsCA=false")
	t.Log("   Ce certificat NE DOIT PAS pouvoir signer d'autres certificats")
	t.Log("   (Validation par inspection - OpenSSL rejetterait toute tentative de signature)")

	// Verify KeyUsage does NOT include CertSign
	if (cert1.KeyUsage & x509.KeyUsageCertSign) != 0 {
		t.Error("❌ SÉCURITÉ: KeyUsageCertSign présent sur certificat non-CA")
		t.Error("   Un certificat serveur/client ne devrait jamais avoir KeyUsageCertSign")
	} else {
		t.Log("✅ KeyUsageCertSign absent (correct pour certificat non-CA)")
	}
}
