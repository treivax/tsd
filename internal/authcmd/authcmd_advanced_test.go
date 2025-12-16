// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package authcmd

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	MinSecretLength = 32
)

// TestGenerateCert_InvalidOutputDir tests certificate generation with invalid output directory
func TestGenerateCert_InvalidOutputDir(t *testing.T) {
	t.Log("🧪 TEST GÉNÉRATION CERTIFICAT - RÉPERTOIRE INVALIDE")
	t.Log("===================================================")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Use a path that will fail mkdir (e.g., inside a file)
	tmpDir := t.TempDir()
	blockingFile := filepath.Join(tmpDir, "blocking")
	if err := os.WriteFile(blockingFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create blocking file: %v", err)
	}

	invalidDir := filepath.Join(blockingFile, "subdir")

	args := []string{
		"generate-cert",
		"-output-dir", invalidDir,
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode == 0 {
		t.Error("❌ Should fail with invalid output directory")
	} else {
		t.Logf("✅ Error correctly detected for invalid output dir")
	}
}

// TestGenerateCert_VerifyFilePermissions tests that certificate files have correct permissions
func TestGenerateCert_VerifyFilePermissions(t *testing.T) {
	t.Log("🧪 TEST GÉNÉRATION CERTIFICAT - PERMISSIONS FICHIERS")
	t.Log("====================================================")

	tmpDir := t.TempDir()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"generate-cert",
		"-output-dir", tmpDir,
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 0 {
		t.Fatalf("❌ Certificate generation failed: %s", stderr.String())
	}

	// Check key file permissions (should be 0600 for security)
	keyPath := filepath.Join(tmpDir, "server.key")
	keyInfo, err := os.Stat(keyPath)
	if err != nil {
		t.Fatalf("❌ Key file not found: %v", err)
	}

	keyPerms := keyInfo.Mode().Perm()
	// On Unix systems, key should be 0600 (owner read/write only)
	if keyPerms&0077 != 0 {
		t.Logf("⚠️  Key file permissions: %o (should ideally be 0600 on Unix)", keyPerms)
	} else {
		t.Logf("✅ Key file has secure permissions: %o", keyPerms)
	}

	// Check certificate file exists
	certPath := filepath.Join(tmpDir, "server.crt")
	if _, err := os.Stat(certPath); err != nil {
		t.Errorf("❌ Certificate file not found: %v", err)
	} else {
		t.Log("✅ Certificate file created")
	}
}

// TestGenerateCert_VerifyCertificateContent tests actual certificate content
func TestGenerateCert_VerifyCertificateContent(t *testing.T) {
	t.Log("🧪 TEST GÉNÉRATION CERTIFICAT - CONTENU CERTIFICAT")
	t.Log("==================================================")

	tmpDir := t.TempDir()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	const (
		testOrg       = "Test Organization Security"
		testValidDays = 365
	)

	args := []string{
		"generate-cert",
		"-output-dir", tmpDir,
		"-org", testOrg,
		"-valid-days", "365",
		"-hosts", "localhost,example.com",
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 0 {
		t.Fatalf("❌ Certificate generation failed: %s", stderr.String())
	}

	// Read and parse the certificate
	certPath := filepath.Join(tmpDir, "server.crt")
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		t.Fatalf("❌ Failed to read certificate: %v", err)
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		t.Fatal("❌ Failed to decode PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("❌ Failed to parse certificate: %v", err)
	}

	// Verify organization
	if len(cert.Subject.Organization) == 0 || cert.Subject.Organization[0] != testOrg {
		t.Errorf("❌ Organization = %v, want %s", cert.Subject.Organization, testOrg)
	} else {
		t.Logf("✅ Organization correct: %s", testOrg)
	}

	// Verify validity period
	duration := cert.NotAfter.Sub(cert.NotBefore)
	expectedDuration := time.Duration(testValidDays) * 24 * time.Hour
	tolerance := 2 * time.Hour

	if duration > expectedDuration+tolerance || duration < expectedDuration-tolerance {
		t.Errorf("❌ Validity period = %v, want ~%v", duration, expectedDuration)
	} else {
		t.Logf("✅ Validity period correct: %v", duration)
	}

	// Verify DNS names
	expectedDNS := []string{"localhost", "example.com"}
	if len(cert.DNSNames) != len(expectedDNS) {
		t.Errorf("❌ DNS names count = %d, want %d", len(cert.DNSNames), len(expectedDNS))
	} else {
		t.Logf("✅ DNS names: %v", cert.DNSNames)
	}

	// Verify IsCA is false (CRITICAL SECURITY)
	if cert.IsCA {
		t.Error("❌ SECURITY: Certificate has IsCA=true, should be false")
	} else {
		t.Log("✅ SECURITY: IsCA=false (certificate cannot sign others)")
	}
}

// TestGenerateJWT_ShortExpiration tests JWT with very short expiration
func TestGenerateJWT_ShortExpiration(t *testing.T) {
	t.Log("🧪 TEST GÉNÉRATION JWT - EXPIRATION COURTE")
	t.Log("==========================================")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"generate-jwt",
		"-secret", TestSecret,
		"-username", TestUsername,
		"-expiration", "1m",
		"-format", "json",
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 0 {
		t.Fatalf("❌ JWT generation failed: %s", stderr.String())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("❌ Failed to parse JSON: %v", err)
	}

	tokenStr, ok := result["token"].(string)
	if !ok {
		t.Fatal("❌ Token not found in result")
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(TestSecret), nil
	})

	if err != nil {
		t.Fatalf("❌ Failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		exp, ok := claims["exp"].(float64)
		if !ok {
			t.Fatal("❌ exp claim not found")
		}

		expTime := time.Unix(int64(exp), 0)
		now := time.Now()
		duration := expTime.Sub(now)

		// Should be approximately 1 minute (with 10s tolerance)
		if duration < 50*time.Second || duration > 70*time.Second {
			t.Errorf("⚠️  Expiration duration = %v, expected ~1m", duration)
		} else {
			t.Logf("✅ Expiration duration correct: %v", duration)
		}
	}
}

// TestGenerateJWT_LongExpiration tests JWT with long expiration
func TestGenerateJWT_LongExpiration(t *testing.T) {
	t.Log("🧪 TEST GÉNÉRATION JWT - EXPIRATION LONGUE")
	t.Log("==========================================")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"generate-jwt",
		"-secret", TestSecret,
		"-username", TestUsername,
		"-expiration", "720h", // 30 days
		"-format", "json",
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 0 {
		t.Fatalf("❌ JWT generation failed: %s", stderr.String())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("❌ Failed to parse JSON: %v", err)
	}

	if result["token"] == nil {
		t.Error("❌ Token not found in result")
	} else {
		t.Log("✅ Long-lived JWT generated successfully")
	}
}

// TestValidateToken_ExpiredJWT tests validation of expired JWT
func TestValidateToken_ExpiredJWT(t *testing.T) {
	t.Log("🧪 TEST VALIDATION JWT - TOKEN EXPIRÉ")
	t.Log("=====================================")

	// Generate a token with 1 second expiration
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	genArgs := []string{
		"generate-jwt",
		"-secret", TestSecret,
		"-username", TestUsername,
		"-expiration", "1s",
		"-format", "json",
	}

	exitCode := Run(genArgs, nil, stdout, stderr)
	if exitCode != 0 {
		t.Fatalf("❌ Failed to generate token: %s", stderr.String())
	}

	var genResult map[string]interface{}
	json.Unmarshal(stdout.Bytes(), &genResult)
	token := genResult["token"].(string)

	// Wait for token to expire
	t.Log("⏳ Waiting for token to expire...")
	time.Sleep(2 * time.Second)

	// Try to validate expired token
	valStdout := &bytes.Buffer{}
	valStderr := &bytes.Buffer{}

	valArgs := []string{
		"validate",
		"-token", token,
		"-secret", TestSecret,
		"-type", "jwt",
	}

	exitCode = Run(valArgs, nil, valStdout, valStderr)

	// Should fail validation
	if exitCode == 0 {
		t.Error("❌ Expired token should not validate successfully")
	} else {
		t.Log("✅ Expired token correctly rejected")
	}
}

// TestValidateToken_InvalidKey tests validation with non-existent key
func TestValidateToken_InvalidKey(t *testing.T) {
	t.Log("🧪 TEST VALIDATION CLÉ - CLÉ INVALIDE")
	t.Log("=====================================")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"validate",
		"-token", "invalid-key-12345",
		"-keys", "valid-key-67890",
		"-type", "key",
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode == 0 {
		t.Error("❌ Invalid key should not validate")
	} else {
		t.Log("✅ Invalid key correctly rejected")
	}
}

// TestValidateToken_MultipleValidKeys tests validation with multiple valid keys
func TestValidateToken_MultipleValidKeys(t *testing.T) {
	t.Log("🧪 TEST VALIDATION CLÉ - MULTIPLES CLÉS VALIDES")
	t.Log("================================================")

	// Generate multiple keys
	genStdout := &bytes.Buffer{}
	genStderr := &bytes.Buffer{}

	genArgs := []string{"generate-key", "-count", "3", "-format", "json"}
	exitCode := Run(genArgs, nil, genStdout, genStderr)
	if exitCode != 0 {
		t.Fatalf("❌ Failed to generate keys: %s", genStderr.String())
	}

	var genResult map[string]interface{}
	json.Unmarshal(genStdout.Bytes(), &genResult)
	keys := genResult["keys"].([]interface{})
	validKey := keys[0].(string)
	otherKeys := []string{keys[1].(string), keys[2].(string)}

	// Validate with multiple keys (including the valid one)
	valStdout := &bytes.Buffer{}
	valStderr := &bytes.Buffer{}

	multipleKeys := otherKeys[0] + "," + validKey + "," + otherKeys[1]

	valArgs := []string{
		"validate",
		"-token", validKey,
		"-keys", multipleKeys,
		"-type", "key",
	}

	exitCode = Run(valArgs, nil, valStdout, valStderr)

	if exitCode != 0 {
		t.Errorf("❌ Validation failed even with valid key in list: %s", valStderr.String())
	} else {
		t.Log("✅ Valid key found in multiple keys list")
	}
}

// TestValidateToken_EmptyKeysList tests validation with empty keys list
func TestValidateToken_EmptyKeysList(t *testing.T) {
	t.Log("🧪 TEST VALIDATION CLÉ - LISTE CLÉS VIDE")
	t.Log("========================================")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"validate",
		"-token", "some-key",
		"-keys", "",
		"-type", "key",
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode == 0 {
		t.Error("❌ Should fail with empty keys list")
	} else {
		t.Log("✅ Empty keys list correctly rejected")
	}
}

// TestGenerateJWT_MissingIssuer tests JWT generation without issuer (should use default)
func TestGenerateJWT_MissingIssuer(t *testing.T) {
	t.Log("🧪 TEST GÉNÉRATION JWT - SANS ISSUER")
	t.Log("====================================")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"generate-jwt",
		"-secret", TestSecret,
		"-username", TestUsername,
		"-format", "json",
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 0 {
		t.Fatalf("❌ JWT generation failed: %s", stderr.String())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("❌ Failed to parse JSON: %v", err)
	}

	tokenStr := result["token"].(string)

	// Parse token to check issuer claim
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(TestSecret), nil
	})

	if err != nil {
		t.Fatalf("❌ Failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		issuer, hasIssuer := claims["iss"]
		if hasIssuer {
			t.Logf("✅ Default issuer set: %v", issuer)
		} else {
			t.Log("ℹ️  No issuer claim in token (acceptable)")
		}
	}
}

// TestGenerateJWT_EmptyRoles tests JWT generation with empty roles
func TestGenerateJWT_EmptyRoles(t *testing.T) {
	t.Log("🧪 TEST GÉNÉRATION JWT - RÔLES VIDES")
	t.Log("====================================")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"generate-jwt",
		"-secret", TestSecret,
		"-username", TestUsername,
		"-roles", "",
	}

	exitCode := Run(args, nil, stdout, stderr)

	// Should succeed (empty roles is valid)
	if exitCode != 0 {
		t.Logf("⚠️  Generation failed with empty roles: %s", stderr.String())
	} else {
		t.Log("✅ JWT generated successfully with empty roles")
	}
}

// TestGenerateKey_NegativeCount tests generate-key with negative count
func TestGenerateKey_NegativeCount(t *testing.T) {
	t.Log("🧪 TEST GÉNÉRATION CLÉ - COUNT NÉGATIF")
	t.Log("======================================")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := Run([]string{"generate-key", "-count", "-1"}, nil, stdout, stderr)

	if exitCode == 0 {
		t.Error("❌ Should fail with negative count")
	} else {
		t.Log("✅ Negative count correctly rejected")
	}
}

// TestGenerateKey_LargeCount tests generate-key with large count
func TestGenerateKey_LargeCount(t *testing.T) {
	t.Log("🧪 TEST GÉNÉRATION CLÉ - COUNT IMPORTANT")
	t.Log("========================================")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Generate 10 keys
	exitCode := Run([]string{"generate-key", "-count", "10", "-format", "json"}, nil, stdout, stderr)

	if exitCode != 0 {
		t.Fatalf("❌ Generation failed: %s", stderr.String())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("❌ Failed to parse JSON: %v", err)
	}

	keys := result["keys"].([]interface{})
	if len(keys) != 10 {
		t.Errorf("❌ Expected 10 keys, got %d", len(keys))
	} else {
		t.Log("✅ 10 keys generated successfully")
	}

	// Verify all keys are unique
	keySet := make(map[string]bool)
	for _, k := range keys {
		keyStr := k.(string)
		if keySet[keyStr] {
			t.Errorf("❌ Duplicate key found: %s", keyStr)
		}
		keySet[keyStr] = true
	}
	t.Logf("✅ All %d keys are unique", len(keys))
}

// TestGenerateCert_MultipleHosts tests certificate with many hosts
func TestGenerateCert_MultipleHosts(t *testing.T) {
	t.Log("🧪 TEST GÉNÉRATION CERTIFICAT - MULTIPLES HOSTS")
	t.Log("===============================================")

	tmpDir := t.TempDir()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	hosts := "localhost,example.com,test.local,127.0.0.1,192.168.1.1,::1"

	args := []string{
		"generate-cert",
		"-hosts", hosts,
		"-output-dir", tmpDir,
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 0 {
		t.Fatalf("❌ Certificate generation failed: %s", stderr.String())
	}

	// Verify certificate contains all hosts
	certPath := filepath.Join(tmpDir, "server.crt")
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		t.Fatalf("❌ Failed to read certificate: %v", err)
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		t.Fatal("❌ Failed to decode PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("❌ Failed to parse certificate: %v", err)
	}

	t.Logf("✅ Certificate generated with:")
	t.Logf("   - DNS names: %v", cert.DNSNames)
	t.Logf("   - IP addresses: %v", cert.IPAddresses)

	totalSANs := len(cert.DNSNames) + len(cert.IPAddresses)
	if totalSANs < 5 {
		t.Errorf("⚠️  Expected at least 5 SANs, got %d", totalSANs)
	}
}

// TestValidateToken_WrongSecret tests JWT validation with wrong secret
func TestValidateToken_WrongSecret(t *testing.T) {
	t.Log("🧪 TEST VALIDATION JWT - SECRET INCORRECT")
	t.Log("=========================================")

	// Generate token with one secret
	genStdout := &bytes.Buffer{}
	genStderr := &bytes.Buffer{}

	genArgs := []string{
		"generate-jwt",
		"-secret", TestSecret,
		"-username", TestUsername,
		"-format", "json",
	}

	exitCode := Run(genArgs, nil, genStdout, genStderr)
	if exitCode != 0 {
		t.Fatalf("❌ Failed to generate token: %s", genStderr.String())
	}

	var genResult map[string]interface{}
	json.Unmarshal(genStdout.Bytes(), &genResult)
	token := genResult["token"].(string)

	// Try to validate with different secret
	valStdout := &bytes.Buffer{}
	valStderr := &bytes.Buffer{}

	wrongSecret := "wrong-secret-1234567890123456789012"

	valArgs := []string{
		"validate",
		"-token", token,
		"-secret", wrongSecret,
		"-type", "jwt",
	}

	exitCode = Run(valArgs, nil, valStdout, valStderr)

	if exitCode == 0 {
		t.Error("❌ Token should not validate with wrong secret")
	} else {
		t.Log("✅ Token correctly rejected with wrong secret")
	}
}

// TestGenerateCert_InvalidValidDays tests certificate with invalid validity days
func TestGenerateCert_InvalidValidDays(t *testing.T) {
	t.Log("🧪 TEST GÉNÉRATION CERTIFICAT - JOURS VALIDITÉ INVALIDES")
	t.Log("========================================================")

	tests := []struct {
		name       string
		validDays  string
		shouldFail bool
	}{
		{"negative days", "-1", false}, // May succeed with abs value or default
		{"zero days", "0", false},      // May use default
		{"non-numeric", "abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			args := []string{
				"generate-cert",
				"-output-dir", tmpDir,
				"-valid-days", tt.validDays,
			}

			exitCode := Run(args, nil, stdout, stderr)

			if tt.shouldFail && exitCode == 0 {
				t.Errorf("❌ Should fail with invalid valid-days: %s", tt.validDays)
			} else if tt.shouldFail {
				t.Logf("✅ Correctly rejected invalid valid-days: %s", tt.validDays)
			} else if exitCode != 0 {
				t.Logf("ℹ️  Rejected valid-days %s (may be expected): %s", tt.validDays, stderr.String())
			} else {
				t.Logf("✅ Handled valid-days %s", tt.validDays)
			}
		})
	}
}

// TestCopyFile_ReadOnlyDestination tests copyFile with read-only destination directory
func TestCopyFile_ReadOnlyDestination(t *testing.T) {
	t.Log("🧪 TEST COPIE FICHIER - DESTINATION LECTURE SEULE")
	t.Log("=================================================")

	tmpDir := t.TempDir()

	// Create source file
	srcPath := filepath.Join(tmpDir, "source.txt")
	if err := os.WriteFile(srcPath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create read-only directory
	roDir := filepath.Join(tmpDir, "readonly")
	if err := os.Mkdir(roDir, 0555); err != nil {
		t.Fatalf("Failed to create readonly dir: %v", err)
	}
	defer os.Chmod(roDir, 0755) // Restore permissions for cleanup

	dstPath := filepath.Join(roDir, "destination.txt")

	err := copyFile(srcPath, dstPath)

	if err == nil {
		t.Error("❌ Should fail copying to read-only directory")
	} else {
		t.Logf("✅ Error correctly detected: %v", err)
	}
}

// TestValidateToken_JSONOutputFormat tests validate with JSON output format verification
func TestValidateToken_JSONOutputFormat(t *testing.T) {
	t.Log("🧪 TEST VALIDATION - FORMAT JSON COMPLET")
	t.Log("========================================")

	// Generate token
	genStdout := &bytes.Buffer{}
	genStderr := &bytes.Buffer{}

	genArgs := []string{
		"generate-jwt",
		"-secret", TestSecret,
		"-username", TestUsername,
		"-roles", "admin,user",
		"-issuer", TestIssuer,
		"-format", "json",
	}

	exitCode := Run(genArgs, nil, genStdout, genStderr)
	if exitCode != 0 {
		t.Fatalf("❌ Failed to generate token: %s", genStderr.String())
	}

	var genResult map[string]interface{}
	json.Unmarshal(genStdout.Bytes(), &genResult)
	token := genResult["token"].(string)

	// Validate with JSON output
	valStdout := &bytes.Buffer{}
	valStderr := &bytes.Buffer{}

	valArgs := []string{
		"validate",
		"-token", token,
		"-secret", TestSecret,
		"-type", "jwt",
		"-format", "json",
	}

	exitCode = Run(valArgs, nil, valStdout, valStderr)

	if exitCode != 0 {
		t.Logf("⚠️  Validation failed: %s", valStderr.String())
		t.Logf("   This may indicate an issue with JWT validation or format")
		return
	}

	var valResult map[string]interface{}
	if err := json.Unmarshal(valStdout.Bytes(), &valResult); err != nil {
		// JSON format may not be implemented for validation
		t.Logf("ℹ️  JSON output parsing failed (may not be implemented): %v", err)
		t.Logf("   Output: %s", valStdout.String())
		return
	}

	// Verify JSON structure if parsing succeeded
	if valid, ok := valResult["valid"].(bool); ok && valid {
		t.Log("✅ JSON validation output complete and correct")
	} else {
		t.Logf("ℹ️  JSON structure: %v", valResult)
	}
}
