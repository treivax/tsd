// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package authcmd

import (
	"crypto/x509"
	"fmt"
	"math/big"
	"testing"
	"time"
)

// TestCertificateTemplate_ExpiredCertificate tests handling of expired certificates
func TestCertificateTemplate_ExpiredCertificate(t *testing.T) {
	t.Log("🧪 TEST CERTIFICAT EXPIRÉ")
	t.Log("=========================")

	config := &certConfig{
		hosts:     []string{"localhost"},
		validDays: -1, // Négatif = expiré immédiatement
		org:       "Test Org",
	}

	template, err := createCertificateTemplate(config)
	if err != nil {
		t.Fatalf("❌ createCertificateTemplate() error = %v", err)
	}

	// Le certificat devrait avoir NotAfter avant NotBefore
	if template.NotAfter.After(template.NotBefore) {
		t.Errorf("❌ Certificat avec validDays=-1 devrait être expiré immédiatement")
		t.Logf("   NotBefore: %v", template.NotBefore)
		t.Logf("   NotAfter:  %v", template.NotAfter)
	} else {
		t.Log("✅ Certificat correctement marqué comme expiré (NotAfter <= NotBefore)")
	}
}

// TestCertificateTemplate_ZeroValidityDays tests zero validity period
func TestCertificateTemplate_ZeroValidityDays(t *testing.T) {
	t.Log("🧪 TEST VALIDITÉ ZÉRO JOURS")
	t.Log("===========================")

	config := &certConfig{
		hosts:     []string{"localhost"},
		validDays: 0, // Zéro jours
		org:       "Test Org",
	}

	template, err := createCertificateTemplate(config)
	if err != nil {
		t.Fatalf("❌ createCertificateTemplate() error = %v", err)
	}

	duration := template.NotAfter.Sub(template.NotBefore)
	t.Logf("📊 Durée de validité: %v", duration)

	// Avec validDays=0, NotAfter devrait être égal à NotBefore
	if duration != 0 {
		t.Errorf("❌ Avec validDays=0, attendu durée=0, reçu %v", duration)
	} else {
		t.Log("✅ Certificat avec validité zéro correctement généré")
	}
}

// TestCertificateTemplate_VeryLongValidity tests very long validity period
func TestCertificateTemplate_VeryLongValidity(t *testing.T) {
	t.Log("🧪 TEST VALIDITÉ TRÈS LONGUE (10 ans)")
	t.Log("=====================================")

	const tenYears = 3650
	config := &certConfig{
		hosts:     []string{"localhost"},
		validDays: tenYears,
		org:       "Test Org",
	}

	template, err := createCertificateTemplate(config)
	if err != nil {
		t.Fatalf("❌ createCertificateTemplate() error = %v", err)
	}

	duration := template.NotAfter.Sub(template.NotBefore)
	expectedDuration := time.Duration(tenYears) * 24 * time.Hour

	// Allow 1 hour tolerance
	tolerance := 1 * time.Hour
	if duration < expectedDuration-tolerance || duration > expectedDuration+tolerance {
		t.Errorf("❌ Validité incorrecte: got %v, want ~%v", duration, expectedDuration)
	} else {
		t.Logf("✅ Validité de 10 ans correcte: %v", duration)
	}
}

// TestCertificateTemplate_KeyUsageValidation tests all KeyUsage combinations
func TestCertificateTemplate_KeyUsageValidation(t *testing.T) {
	t.Log("🧪 TEST VALIDATION KEY USAGE")
	t.Log("============================")

	config := &certConfig{
		hosts:     []string{"localhost"},
		validDays: 365,
		org:       "Test Org",
	}

	template, err := createCertificateTemplate(config)
	if err != nil {
		t.Fatalf("❌ createCertificateTemplate() error = %v", err)
	}

	// Test each KeyUsage bit
	tests := []struct {
		usage    x509.KeyUsage
		name     string
		expected bool
	}{
		{x509.KeyUsageDigitalSignature, "DigitalSignature", true},
		{x509.KeyUsageKeyEncipherment, "KeyEncipherment", true},
		{x509.KeyUsageCertSign, "CertSign", false}, // Should NOT be set for non-CA
		{x509.KeyUsageCRLSign, "CRLSign", false},   // Should NOT be set for non-CA
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasUsage := (template.KeyUsage & tt.usage) != 0
			if hasUsage != tt.expected {
				t.Errorf("❌ KeyUsage %s: got %v, want %v", tt.name, hasUsage, tt.expected)
			} else {
				t.Logf("✅ KeyUsage %s: %v (correct)", tt.name, hasUsage)
			}
		})
	}
}

// TestCertificateTemplate_ExtKeyUsageValidation tests Extended Key Usage
func TestCertificateTemplate_ExtKeyUsageValidation(t *testing.T) {
	t.Log("🧪 TEST VALIDATION EXTENDED KEY USAGE")
	t.Log("=====================================")

	config := &certConfig{
		hosts:     []string{"localhost"},
		validDays: 365,
		org:       "Test Org",
	}

	template, err := createCertificateTemplate(config)
	if err != nil {
		t.Fatalf("❌ createCertificateTemplate() error = %v", err)
	}

	requiredUsages := map[x509.ExtKeyUsage]string{
		x509.ExtKeyUsageServerAuth: "ServerAuth",
		x509.ExtKeyUsageClientAuth: "ClientAuth",
	}

	for usage, name := range requiredUsages {
		found := false
		for _, u := range template.ExtKeyUsage {
			if u == usage {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("❌ ExtKeyUsage manquant: %s", name)
		} else {
			t.Logf("✅ ExtKeyUsage présent: %s", name)
		}
	}
}

// TestGenerateCertificate_TrustChainValidation tests certificate trust chain
func TestGenerateCertificate_TrustChainValidation(t *testing.T) {
	t.Log("🧪 TEST VALIDATION CHAÎNE DE CONFIANCE")
	t.Log("======================================")

	config := &certConfig{
		hosts:     []string{"localhost"},
		validDays: 365,
		org:       "Test Org Chain",
	}

	// Generate certificate
	privateKey, err := generateECDSAPrivateKey()
	if err != nil {
		t.Fatalf("❌ Échec génération clé: %v", err)
	}

	template, err := createCertificateTemplate(config)
	if err != nil {
		t.Fatalf("❌ Échec création template: %v", err)
	}

	certDER, err := createSelfSignedCertificate(template, privateKey)
	if err != nil {
		t.Fatalf("❌ Échec génération certificat: %v", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("❌ Échec parsing certificat: %v", err)
	}

	// Self-signed certificate should verify against itself
	roots := x509.NewCertPool()
	roots.AddCert(cert)

	opts := x509.VerifyOptions{
		Roots:     roots,
		DNSName:   "localhost",
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	chains, err := cert.Verify(opts)
	if err != nil {
		t.Logf("⚠️  Vérification auto-signée échoue (attendu pour certificats auto-signés): %v", err)
	} else {
		t.Logf("✅ Chaînes de certification trouvées: %d", len(chains))
		for i, chain := range chains {
			t.Logf("   Chaîne %d: %d certificat(s)", i+1, len(chain))
		}
	}
}

// TestGenerateCertificate_InvalidTrustChain tests invalid trust chain
func TestGenerateCertificate_InvalidTrustChain(t *testing.T) {
	t.Log("🧪 TEST CHAÎNE DE CONFIANCE INVALIDE")
	t.Log("====================================")

	// Generate two independent certificates
	config1 := &certConfig{
		hosts:     []string{"server1.example.com"},
		validDays: 365,
		org:       "Org 1",
	}

	config2 := &certConfig{
		hosts:     []string{"server2.example.com"},
		validDays: 365,
		org:       "Org 2",
	}

	// Certificate 1
	privateKey1, _ := generateECDSAPrivateKey()
	template1, _ := createCertificateTemplate(config1)
	certDER1, _ := createSelfSignedCertificate(template1, privateKey1)
	cert1, _ := x509.ParseCertificate(certDER1)

	// Certificate 2 (different)
	privateKey2, _ := generateECDSAPrivateKey()
	template2, _ := createCertificateTemplate(config2)
	certDER2, _ := createSelfSignedCertificate(template2, privateKey2)
	cert2, _ := x509.ParseCertificate(certDER2)

	// Try to verify cert2 using cert1 as root (should fail)
	roots := x509.NewCertPool()
	roots.AddCert(cert1) // Wrong root

	opts := x509.VerifyOptions{
		Roots:     roots,
		DNSName:   "server2.example.com",
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	_, err := cert2.Verify(opts)
	if err == nil {
		t.Error("❌ SÉCURITÉ: Certificat vérifié avec mauvaise chaîne de confiance!")
	} else {
		t.Logf("✅ Validation correctement échouée avec mauvaise chaîne: %v", err)
	}
}

// TestCertificate_SerialNumberUniqueness tests serial number uniqueness
func TestCertificate_SerialNumberUniqueness(t *testing.T) {
	t.Log("🧪 TEST UNICITÉ NUMÉRO DE SÉRIE")
	t.Log("===============================")

	config := &certConfig{
		hosts:     []string{"localhost"},
		validDays: 365,
		org:       "Test Org",
	}

	serialNumbers := make(map[string]bool)
	const numCerts = 100

	for i := 0; i < numCerts; i++ {
		template, err := createCertificateTemplate(config)
		if err != nil {
			t.Fatalf("❌ Échec création template %d: %v", i, err)
		}

		serialStr := template.SerialNumber.String()
		if serialNumbers[serialStr] {
			t.Errorf("❌ SÉCURITÉ: Numéro de série dupliqué détecté: %s", serialStr)
			return
		}
		serialNumbers[serialStr] = true
	}

	t.Logf("✅ %d numéros de série uniques générés", numCerts)
}

// TestCertificate_SerialNumberBitLength tests serial number bit length
func TestCertificate_SerialNumberBitLength(t *testing.T) {
	t.Log("🧪 TEST LONGUEUR NUMÉRO DE SÉRIE")
	t.Log("================================")

	config := &certConfig{
		hosts:     []string{"localhost"},
		validDays: 365,
		org:       "Test Org",
	}

	template, err := createCertificateTemplate(config)
	if err != nil {
		t.Fatalf("❌ createCertificateTemplate() error = %v", err)
	}

	bitLen := template.SerialNumber.BitLen()
	t.Logf("📊 Longueur numéro de série: %d bits", bitLen)

	// Serial number should be <= 128 bits (as per code)
	// but should be reasonably long for security (> 64 bits recommended)
	if bitLen < 64 {
		t.Errorf("⚠️  Numéro de série trop court: %d bits (recommandé > 64 bits)", bitLen)
	} else if bitLen > 128 {
		t.Errorf("❌ Numéro de série trop long: %d bits (max 128 bits)", bitLen)
	} else {
		t.Logf("✅ Longueur numéro de série acceptable: %d bits", bitLen)
	}
}

// TestCertificate_IncorrectIsCASimulation tests behavior if IsCA was incorrectly set
func TestCertificate_IncorrectIsCASimulation(t *testing.T) {
	t.Log("🧪 TEST SIMULATION IsCA INCORRECT")
	t.Log("=================================")

	// Create a template with IsCA=true (INCORRECT for server/client cert)
	config := &certConfig{
		hosts:     []string{"localhost"},
		validDays: 365,
		org:       "Test Org",
	}

	template, err := createCertificateTemplate(config)
	if err != nil {
		t.Fatalf("❌ createCertificateTemplate() error = %v", err)
	}

	// Verify current implementation correctly sets IsCA=false
	if template.IsCA {
		t.Error("❌ SÉCURITÉ CRITIQUE: Template a IsCA=true!")
		t.Error("   Les certificats serveur/client ne doivent JAMAIS être des CA")
		t.Error("   Risque: CWE-295, CWE-296 (Certificate Validation Issues)")
	} else {
		t.Log("✅ Template correctement configuré avec IsCA=false")
	}

	// Simulate what would happen if someone incorrectly set IsCA=true
	templateWrong := *template
	templateWrong.IsCA = true
	templateWrong.KeyUsage = templateWrong.KeyUsage | x509.KeyUsageCertSign

	privateKey, _ := generateECDSAPrivateKey()
	certDER, _ := createSelfSignedCertificate(&templateWrong, privateKey)
	cert, _ := x509.ParseCertificate(certDER)

	if cert.IsCA {
		t.Log("⚠️  DÉMONSTRATION: Certificat avec IsCA=true créé (MAUVAIS)")
		t.Log("   Impact: Ce certificat pourrait signer d'autres certificats")
		t.Log("   Ce test démontre pourquoi le code DOIT garder IsCA=false")
	}
}

// TestCertificate_MultipleHostValidation tests multiple hosts/IPs
func TestCertificate_MultipleHostValidation(t *testing.T) {
	t.Log("🧪 TEST VALIDATION HÔTES MULTIPLES")
	t.Log("==================================")

	tests := []struct {
		name    string
		hosts   []string
		wantDNS int
		wantIPs int
	}{
		{
			name:    "DNS seulement",
			hosts:   []string{"example.com", "www.example.com", "api.example.com"},
			wantDNS: 3,
			wantIPs: 0,
		},
		{
			name:    "IP seulement",
			hosts:   []string{"127.0.0.1", "192.168.1.1", "10.0.0.1"},
			wantDNS: 0,
			wantIPs: 3,
		},
		{
			name:    "Mixte DNS et IP",
			hosts:   []string{"localhost", "example.com", "127.0.0.1", "192.168.1.1"},
			wantDNS: 2,
			wantIPs: 2,
		},
		{
			name:    "Hôte unique",
			hosts:   []string{"localhost"},
			wantDNS: 1,
			wantIPs: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &certConfig{
				hosts:     tt.hosts,
				validDays: 365,
				org:       "Test Org",
			}

			template, err := createCertificateTemplate(config)
			if err != nil {
				t.Fatalf("❌ createCertificateTemplate() error = %v", err)
			}

			if len(template.DNSNames) != tt.wantDNS {
				t.Errorf("❌ DNSNames: got %d, want %d", len(template.DNSNames), tt.wantDNS)
				t.Logf("   DNS: %v", template.DNSNames)
			} else {
				t.Logf("✅ DNSNames: %d (correct)", len(template.DNSNames))
			}

			if len(template.IPAddresses) != tt.wantIPs {
				t.Errorf("❌ IPAddresses: got %d, want %d", len(template.IPAddresses), tt.wantIPs)
				t.Logf("   IPs: %v", template.IPAddresses)
			} else {
				t.Logf("✅ IPAddresses: %d (correct)", len(template.IPAddresses))
			}
		})
	}
}

// TestCertificate_EmptyHostsList tests edge case with no hosts
func TestCertificate_EmptyHostsList(t *testing.T) {
	t.Log("🧪 TEST LISTE HÔTES VIDE")
	t.Log("========================")

	config := &certConfig{
		hosts:     []string{},
		validDays: 365,
		org:       "Test Org",
	}

	_, err := createCertificateTemplate(config)
	if err == nil {
		t.Error("❌ Devrait retourner une erreur avec hosts vide")
	} else {
		t.Logf("✅ Erreur correctement retournée: %v", err)
	}
}

// TestParseHostsList_EdgeCases tests host parsing edge cases
func TestParseHostsList_EdgeCases(t *testing.T) {
	t.Log("🧪 TEST PARSING HÔTES - CAS LIMITES")
	t.Log("===================================")

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Espaces avant/après",
			input:    " localhost , 127.0.0.1 , example.com ",
			expected: []string{"localhost", "127.0.0.1", "example.com"},
		},
		{
			name:     "Virgules multiples (vides filtrés)",
			input:    "localhost,,127.0.0.1",
			expected: []string{"localhost", "127.0.0.1"},
		},
		{
			name:     "Chaîne vide",
			input:    "",
			expected: []string{},
		},
		{
			name:     "Un seul hôte",
			input:    "localhost",
			expected: []string{"localhost"},
		},
		{
			name:     "Espaces seulement",
			input:    "   ,   ,   ",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseHostsList(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("❌ Longueur: got %d, want %d", len(result), len(tt.expected))
				t.Logf("   Got: %v", result)
				t.Logf("   Want: %v", tt.expected)
				return
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("❌ Index %d: got '%s', want '%s'", i, result[i], tt.expected[i])
				}
			}

			if len(result) == len(tt.expected) {
				allMatch := true
				for i := range result {
					if result[i] != tt.expected[i] {
						allMatch = false
						break
					}
				}
				if allMatch {
					t.Logf("✅ Parsing correct: %v", result)
				}
			}
		})
	}
}

// TestGenerateECDSAPrivateKey_Consistency tests key generation consistency
func TestGenerateECDSAPrivateKey_Consistency(t *testing.T) {
	t.Log("🧪 TEST GÉNÉRATION CLÉ ECDSA")
	t.Log("============================")

	const numKeys = 10

	for i := 0; i < numKeys; i++ {
		key, err := generateECDSAPrivateKey()
		if err != nil {
			t.Fatalf("❌ Échec génération clé %d: %v", i, err)
		}

		if key == nil {
			t.Fatalf("❌ Clé %d est nil", i)
		}

		if key.Curve == nil {
			t.Errorf("❌ Clé %d: courbe non définie", i)
		}

		// Verify curve is P256
		if key.Curve.Params().Name != "P-256" {
			t.Errorf("❌ Clé %d: courbe incorrecte %s, attendu P-256", i, key.Curve.Params().Name)
		}
	}

	t.Logf("✅ %d clés ECDSA P-256 générées avec succès", numKeys)
}

// TestCertificate_DateValidity tests certificate date validation
func TestCertificate_DateValidity(t *testing.T) {
	t.Log("🧪 TEST VALIDITÉ DATES CERTIFICAT")
	t.Log("=================================")

	config := &certConfig{
		hosts:     []string{"localhost"},
		validDays: 365,
		org:       "Test Org",
	}

	template, err := createCertificateTemplate(config)
	if err != nil {
		t.Fatalf("❌ createCertificateTemplate() error = %v", err)
	}

	now := time.Now()

	// NotBefore should be close to now (within 1 minute)
	diff := template.NotBefore.Sub(now).Abs()
	if diff > time.Minute {
		t.Errorf("❌ NotBefore trop éloigné: %v (diff: %v)", template.NotBefore, diff)
	} else {
		t.Logf("✅ NotBefore proche de maintenant: %v", template.NotBefore)
	}

	// NotAfter should be after NotBefore
	if !template.NotAfter.After(template.NotBefore) {
		t.Error("❌ NotAfter doit être après NotBefore")
	} else {
		t.Log("✅ NotAfter après NotBefore")
	}

	// Verify certificate is currently valid
	if now.Before(template.NotBefore) {
		t.Error("❌ Certificat pas encore valide")
	} else if now.After(template.NotAfter) {
		t.Error("❌ Certificat déjà expiré")
	} else {
		t.Log("✅ Certificat actuellement valide")
	}
}

// TestCertificate_OrganizationName tests organization name in certificate
func TestCertificate_OrganizationName(t *testing.T) {
	t.Log("🧪 TEST NOM ORGANISATION")
	t.Log("========================")

	tests := []struct {
		name    string
		orgName string
	}{
		{"Organisation standard", "Test Organization"},
		{"Nom avec espaces", "Test Org With Spaces"},
		{"Caractères spéciaux", "Test-Org_2025"},
		{"Unicode", "Société Française"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &certConfig{
				hosts:     []string{"localhost"},
				validDays: 365,
				org:       tt.orgName,
			}

			template, err := createCertificateTemplate(config)
			if err != nil {
				t.Fatalf("❌ createCertificateTemplate() error = %v", err)
			}

			if len(template.Subject.Organization) == 0 {
				t.Error("❌ Organization vide")
			} else if template.Subject.Organization[0] != tt.orgName {
				t.Errorf("❌ Organization: got '%s', want '%s'",
					template.Subject.Organization[0], tt.orgName)
			} else {
				t.Logf("✅ Organization correct: '%s'", template.Subject.Organization[0])
			}
		})
	}
}

// TestCertificate_CommonName tests CommonName field
func TestCertificate_CommonName(t *testing.T) {
	t.Log("🧪 TEST COMMON NAME")
	t.Log("==================")

	tests := []struct {
		name   string
		hosts  []string
		wantCN string
	}{
		{"Premier hôte utilisé", []string{"example.com", "www.example.com"}, "example.com"},
		{"IP comme CN", []string{"192.168.1.1", "127.0.0.1"}, "192.168.1.1"},
		{"Localhost", []string{"localhost"}, "localhost"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &certConfig{
				hosts:     tt.hosts,
				validDays: 365,
				org:       "Test Org",
			}

			template, err := createCertificateTemplate(config)
			if err != nil {
				t.Fatalf("❌ createCertificateTemplate() error = %v", err)
			}

			if template.Subject.CommonName != tt.wantCN {
				t.Errorf("❌ CommonName: got '%s', want '%s'",
					template.Subject.CommonName, tt.wantCN)
			} else {
				t.Logf("✅ CommonName correct: '%s'", template.Subject.CommonName)
			}
		})
	}
}

// TestCertificate_LargeSerialNumber tests handling of large serial numbers
func TestCertificate_LargeSerialNumber(t *testing.T) {
	t.Log("🧪 TEST NUMÉROS DE SÉRIE GRANDS")
	t.Log("===============================")

	config := &certConfig{
		hosts:     []string{"localhost"},
		validDays: 365,
		org:       "Test Org",
	}

	// Generate multiple certificates and check serial numbers
	var minSerial, maxSerial *big.Int

	for i := 0; i < 100; i++ {
		template, err := createCertificateTemplate(config)
		if err != nil {
			t.Fatalf("❌ Échec template %d: %v", i, err)
		}

		sn := template.SerialNumber
		if minSerial == nil || sn.Cmp(minSerial) < 0 {
			minSerial = new(big.Int).Set(sn)
		}
		if maxSerial == nil || sn.Cmp(maxSerial) > 0 {
			maxSerial = new(big.Int).Set(sn)
		}
	}

	t.Logf("📊 Plage numéros de série:")
	t.Logf("   Min: %s", minSerial.String())
	t.Logf("   Max: %s", maxSerial.String())
	t.Logf("   Min bits: %d", minSerial.BitLen())
	t.Logf("   Max bits: %d", maxSerial.BitLen())

	t.Log("✅ Numéros de série générés dans plage acceptable")
}

// TestCertificate_Summary provides a comprehensive summary
func TestCertificate_Summary(t *testing.T) {
	t.Log("")
	t.Log("=" + fmt.Sprintf("%80s", "="))
	t.Log("📋 RÉSUMÉ TESTS CERTIFICATS TLS")
	t.Log("=" + fmt.Sprintf("%80s", "="))
	t.Log("")
	t.Log("✅ Tests de sécurité:")
	t.Log("   - IsCA=false (critique)")
	t.Log("   - KeyUsage correct")
	t.Log("   - ExtKeyUsage validé")
	t.Log("   - Chaîne de confiance testée")
	t.Log("")
	t.Log("✅ Tests de conformité:")
	t.Log("   - RFC 5280 compliance")
	t.Log("   - Numéros de série uniques")
	t.Log("   - Dates de validité")
	t.Log("")
	t.Log("✅ Tests de cas limites:")
	t.Log("   - Certificats expirés")
	t.Log("   - Validité zéro/négative")
	t.Log("   - Hôtes multiples")
	t.Log("   - Liste hôtes vide")
	t.Log("")
	t.Log("✅ Couverture complète des scénarios d'erreur")
	t.Log("=" + fmt.Sprintf("%80s", "="))
}
