// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package authcmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// Test constants
const (
	TestSecret   = "test-secret-12345678901234567890123"
	TestUsername = "testuser"
	TestRole     = "admin"
	TestIssuer   = "test-issuer"
)

// TestRun_NoArgs tests Run with no arguments
func TestRun_NoArgs(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := Run([]string{}, nil, stdout, stderr)

	if exitCode != 1 {
		t.Errorf("Run() exitCode = %d, want 1", exitCode)
	}

	if !strings.Contains(stdout.String(), "USAGE") {
		t.Errorf("output should contain USAGE")
	}
}

// TestRun_Help tests Run with help command
func TestRun_Help(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "help", args: []string{"help"}},
		{name: "-h", args: []string{"-h"}},
		{name: "--help", args: []string{"--help"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			exitCode := Run(tt.args, nil, stdout, stderr)

			if exitCode != 0 {
				t.Errorf("Run() exitCode = %d, want 0", exitCode)
			}

			if !strings.Contains(stdout.String(), "USAGE") {
				t.Errorf("help output should contain USAGE")
			}
		})
	}
}

// TestRun_Version tests Run with version command
func TestRun_Version(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "version", args: []string{"version"}},
		{name: "-v", args: []string{"-v"}},
		{name: "--version", args: []string{"--version"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			exitCode := Run(tt.args, nil, stdout, stderr)

			if exitCode != 0 {
				t.Errorf("Run() exitCode = %d, want 0", exitCode)
			}

			output := stdout.String()
			if !strings.Contains(output, "version") || !strings.Contains(output, Version) {
				t.Errorf("version output should contain version info, got: %s", output)
			}
		})
	}
}

// TestRun_UnknownCommand tests Run with unknown command
func TestRun_UnknownCommand(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := Run([]string{"unknown-command"}, nil, stdout, stderr)

	if exitCode != 1 {
		t.Errorf("Run() exitCode = %d, want 1", exitCode)
	}

	if !strings.Contains(stderr.String(), "Commande inconnue") {
		t.Errorf("error should mention unknown command")
	}
}

// TestGenerateKey_Single tests generating a single key
func TestGenerateKey_Single(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := Run([]string{"generate-key"}, nil, stdout, stderr)

	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Clé(s) API générée(s)") {
		t.Errorf("output should contain key generation message")
	}

	if !strings.Contains(output, "IMPORTANT") {
		t.Errorf("output should contain warning message")
	}
}

// TestGenerateKey_Multiple tests generating multiple keys
func TestGenerateKey_Multiple(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := Run([]string{"generate-key", "-count", "3"}, nil, stdout, stderr)

	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Clé 1") {
		t.Errorf("output should contain key numbering")
	}
	if !strings.Contains(output, "Clé 3") {
		t.Errorf("output should contain all keys")
	}
}

// TestGenerateKey_JSONFormat tests generating keys in JSON format
func TestGenerateKey_JSONFormat(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := Run([]string{"generate-key", "-count", "2", "-format", "json"}, nil, stdout, stderr)

	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("output should be valid JSON, got error: %v", err)
		return
	}

	keys, ok := result["keys"].([]interface{})
	if !ok {
		t.Errorf("result should contain keys array")
		return
	}

	if len(keys) != 2 {
		t.Errorf("keys length = %d, want 2", len(keys))
	}

	count, ok := result["count"].(float64)
	if !ok || int(count) != 2 {
		t.Errorf("count should be 2")
	}
}

// TestGenerateKey_InvalidFlag tests generate-key with invalid flag
func TestGenerateKey_InvalidFlag(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := Run([]string{"generate-key", "-invalid"}, nil, stdout, stderr)

	if exitCode != 1 {
		t.Errorf("Run() exitCode = %d, want 1", exitCode)
	}
}

// TestGenerateJWT_Success tests generating a JWT successfully
func TestGenerateJWT_Success(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"generate-jwt",
		"-secret", TestSecret,
		"-username", TestUsername,
		"-roles", TestRole,
		"-issuer", TestIssuer,
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "JWT généré") {
		t.Errorf("output should contain JWT generation message")
	}
}

// TestGenerateJWT_MissingSecret tests generating JWT without secret
func TestGenerateJWT_MissingSecret(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"generate-jwt",
		"-username", TestUsername,
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 1 {
		t.Errorf("Run() exitCode = %d, want 1", exitCode)
	}

	if !strings.Contains(stderr.String(), "secret") {
		t.Errorf("error should mention missing secret")
	}
}

// TestGenerateJWT_MissingUsername tests generating JWT without username
func TestGenerateJWT_MissingUsername(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"generate-jwt",
		"-secret", TestSecret,
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 1 {
		t.Errorf("Run() exitCode = %d, want 1", exitCode)
	}

	if !strings.Contains(stderr.String(), "username") {
		t.Errorf("error should mention missing username")
	}
}

// TestGenerateJWT_JSONFormat tests generating JWT in JSON format
func TestGenerateJWT_JSONFormat(t *testing.T) {
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
		t.Errorf("Run() exitCode = %d, want 0, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("output should be valid JSON, got error: %v", err)
		return
	}

	if _, ok := result["token"]; !ok {
		t.Errorf("result should contain token field")
	}

	if _, ok := result["username"]; !ok {
		t.Errorf("result should contain username field")
	}
}

// TestValidateToken_JWT tests validating a JWT token
func TestValidateToken_JWT(t *testing.T) {
	// First generate a token
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
		t.Fatalf("Failed to generate token: %s", genStderr.String())
	}

	var genResult map[string]interface{}
	if err := json.Unmarshal(genStdout.Bytes(), &genResult); err != nil {
		t.Fatalf("Failed to parse generated token: %v", err)
	}

	token := genResult["token"].(string)

	// Now validate it
	valStdout := &bytes.Buffer{}
	valStderr := &bytes.Buffer{}

	valArgs := []string{
		"validate",
		"-token", token,
		"-secret", TestSecret,
		"-type", "jwt",
	}

	exitCode = Run(valArgs, nil, valStdout, valStderr)

	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0, stderr: %s", exitCode, valStderr.String())
	}

	output := valStdout.String()
	if !strings.Contains(output, "valide") {
		t.Errorf("output should indicate token is valid")
	}
}

// TestValidateToken_InvalidJWT tests validating an invalid JWT
func TestValidateToken_InvalidJWT(t *testing.T) {
	valStdout := &bytes.Buffer{}
	valStderr := &bytes.Buffer{}

	valArgs := []string{
		"validate",
		"-token", "invalid.jwt.token",
		"-secret", TestSecret,
		"-type", "jwt",
	}

	exitCode := Run(valArgs, nil, valStdout, valStderr)

	if exitCode != 1 {
		t.Errorf("Run() exitCode = %d, want 1", exitCode)
	}

	// Check that there was an error in either stdout or stderr
	if valStderr.Len() == 0 && valStdout.Len() == 0 {
		t.Errorf("should have error output in stdout or stderr, both empty")
	}
}

// TestValidateToken_Key tests validating an API key
func TestValidateToken_Key(t *testing.T) {
	// Generate a key first
	genStdout := &bytes.Buffer{}
	genStderr := &bytes.Buffer{}

	genArgs := []string{"generate-key", "-format", "json"}
	exitCode := Run(genArgs, nil, genStdout, genStderr)
	if exitCode != 0 {
		t.Fatalf("Failed to generate key: %s", genStderr.String())
	}

	var genResult map[string]interface{}
	if err := json.Unmarshal(genStdout.Bytes(), &genResult); err != nil {
		t.Fatalf("Failed to parse generated key: %v", err)
	}

	keys := genResult["keys"].([]interface{})
	key := keys[0].(string)

	// Now validate it
	valStdout := &bytes.Buffer{}
	valStderr := &bytes.Buffer{}

	valArgs := []string{
		"validate",
		"-token", key,
		"-keys", key,
		"-type", "key",
	}

	exitCode = Run(valArgs, nil, valStdout, valStderr)

	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0, stderr: %s", exitCode, valStderr.String())
	}

	output := valStdout.String()
	if !strings.Contains(output, "valide") {
		t.Errorf("output should indicate token is valid")
	}
}

// TestValidateToken_MissingType tests validate without type flag
func TestValidateToken_MissingType(t *testing.T) {
	valStdout := &bytes.Buffer{}
	valStderr := &bytes.Buffer{}

	valArgs := []string{"validate", "-token", "some-token"}

	exitCode := Run(valArgs, nil, valStdout, valStderr)

	if exitCode != 1 {
		t.Errorf("Run() exitCode = %d, want 1", exitCode)
	}

	if !strings.Contains(valStderr.String(), "type") {
		t.Errorf("error should mention missing type")
	}
}

// TestValidateToken_JSONFormat tests validate with JSON output
func TestValidateToken_JSONFormat(t *testing.T) {
	// Generate a token
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
		t.Fatalf("Failed to generate token: %s", genStderr.String())
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
		t.Errorf("Run() exitCode = %d, want 0, stderr: %s", exitCode, valStderr.String())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(valStdout.Bytes(), &result); err != nil {
		t.Errorf("output should be valid JSON, got error: %v", err)
		return
	}

	if valid, ok := result["valid"].(bool); !ok || !valid {
		t.Errorf("result should indicate token is valid")
	}
}

// TestGenerateCert_Success tests generating a certificate
func TestGenerateCert_Success(t *testing.T) {
	tmpDir := t.TempDir()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"generate-cert",
		"-hosts", "localhost",
		"-output-dir", tmpDir,
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Certificat") {
		t.Errorf("output should contain certificate generation message")
	}
}

// TestGenerateCert_DefaultHosts tests generate-cert with default hosts
func TestGenerateCert_DefaultHosts(t *testing.T) {
	tmpDir := t.TempDir()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{"generate-cert", "-output-dir", tmpDir}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Certificat") {
		t.Errorf("output should contain certificate generation message")
	}
}

// TestPrintHelp tests printHelp function
func TestPrintHelp(t *testing.T) {
	buf := &bytes.Buffer{}
	printHelp(buf)

	output := buf.String()
	requiredSections := []string{
		"USAGE",
		"COMMANDES",
		"generate-key",
		"generate-jwt",
		"validate",
		"generate-cert",
	}

	for _, section := range requiredSections {
		if !strings.Contains(output, section) {
			t.Errorf("help output should contain %q", section)
		}
	}
}

// TestVersionConstant tests Version constant
func TestVersionConstant(t *testing.T) {
	if Version == "" {
		t.Errorf("Version constant should not be empty")
	}
}

// TestGenerateKey_InvalidCount tests generate-key with invalid count
func TestGenerateKey_InvalidCount(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := Run([]string{"generate-key", "-count", "invalid"}, nil, stdout, stderr)

	if exitCode != 1 {
		t.Errorf("Run() exitCode = %d, want 1", exitCode)
	}
}

// TestGenerateJWT_MultipleRoles tests generating JWT with multiple roles
func TestGenerateJWT_MultipleRoles(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"generate-jwt",
		"-secret", TestSecret,
		"-username", TestUsername,
		"-roles", "admin,user,viewer",
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "JWT généré") {
		t.Errorf("output should contain JWT generation message")
	}
}

// TestValidateToken_EmptyToken tests validating empty token
func TestValidateToken_EmptyToken(t *testing.T) {
	valStdout := &bytes.Buffer{}
	valStderr := &bytes.Buffer{}

	valArgs := []string{
		"validate",
		"-type", "jwt",
		"-secret", TestSecret,
	}

	exitCode := Run(valArgs, nil, valStdout, valStderr)

	if exitCode != 1 {
		t.Errorf("Run() exitCode = %d, want 1", exitCode)
	}

	if !strings.Contains(valStderr.String(), "token") {
		t.Errorf("error should mention missing token")
	}
}

// TestGenerateJWT_WithExpiration tests JWT generation with custom expiration
func TestGenerateJWT_WithExpiration(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"generate-jwt",
		"-secret", TestSecret,
		"-username", TestUsername,
		"-expiration", "2h",
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "JWT généré") {
		t.Errorf("output should contain JWT generation message")
	}
}

// TestGenerateJWT_InteractiveMode tests interactive mode for JWT generation
func TestGenerateJWT_InteractiveMode(t *testing.T) {
	t.Skip("Interactive mode is difficult to test properly - requires real stdin interaction")
}

// TestValidateToken_InteractiveMode tests validate in interactive mode
func TestValidateToken_InteractiveMode(t *testing.T) {
	// First generate a token
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
		t.Fatalf("Failed to generate token: %s", genStderr.String())
	}

	var genResult map[string]interface{}
	json.Unmarshal(genStdout.Bytes(), &genResult)
	token := genResult["token"].(string)

	// Now validate in interactive mode
	input := token + "\njwt\n" + TestSecret + "\n"
	stdin := strings.NewReader(input)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"validate",
		"-i",
	}

	exitCode = Run(args, stdin, stdout, stderr)

	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0, stderr: %s", exitCode, stderr.String())
	}
}

// TestValidateToken_KeyInteractive tests key validation in interactive mode
func TestValidateToken_KeyInteractive(t *testing.T) {
	// Generate a key
	genStdout := &bytes.Buffer{}
	genStderr := &bytes.Buffer{}

	genArgs := []string{"generate-key", "-format", "json"}
	exitCode := Run(genArgs, nil, genStdout, genStderr)
	if exitCode != 0 {
		t.Fatalf("Failed to generate key: %s", genStderr.String())
	}

	var genResult map[string]interface{}
	json.Unmarshal(genStdout.Bytes(), &genResult)
	keys := genResult["keys"].([]interface{})
	key := keys[0].(string)

	// Validate in interactive mode
	input := key + "\nkey\n" + key + "\n"
	stdin := strings.NewReader(input)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"validate",
		"-i",
	}

	exitCode = Run(args, stdin, stdout, stderr)

	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0, stderr: %s", exitCode, stderr.String())
	}
}

// TestGenerateKey_ZeroCount tests generate-key with count=0
func TestGenerateKey_ZeroCount(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := Run([]string{"generate-key", "-count", "0"}, nil, stdout, stderr)

	// Should succeed but generate no keys
	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0", exitCode)
	}
}

// TestGenerateCert_CustomValidDays tests certificate with custom validity
func TestGenerateCert_CustomValidDays(t *testing.T) {
	tmpDir := t.TempDir()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"generate-cert",
		"-hosts", "localhost,example.com",
		"-output-dir", tmpDir,
		"-valid-days", "730",
		"-org", "Test Organization",
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Certificat") {
		t.Errorf("output should contain certificate generation message")
	}
}

// TestGenerateCert_JSONFormat tests certificate generation with JSON output
func TestGenerateCert_JSONFormat(t *testing.T) {
	tmpDir := t.TempDir()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"generate-cert",
		"-hosts", "localhost",
		"-output-dir", tmpDir,
		"-format", "json",
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 0 {
		t.Errorf("Run() exitCode = %d, want 0, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("output should be valid JSON, got error: %v", err)
		return
	}

	// Check for any certificate-related fields
	if len(result) == 0 {
		t.Errorf("result should not be empty")
	}
}

// TestValidateToken_MissingSecret tests JWT validation without secret
func TestValidateToken_MissingSecret(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"validate",
		"-token", "some-token",
		"-type", "jwt",
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 1 {
		t.Errorf("Run() exitCode = %d, want 1", exitCode)
	}

	if !strings.Contains(stderr.String(), "secret") {
		t.Errorf("error should mention missing secret")
	}
}

// TestValidateToken_MissingKeys tests key validation without keys
func TestValidateToken_MissingKeys(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"validate",
		"-token", "some-token",
		"-type", "key",
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 1 {
		t.Errorf("Run() exitCode = %d, want 1", exitCode)
	}

	if !strings.Contains(stderr.String(), "clés") && !strings.Contains(stderr.String(), "keys") {
		t.Errorf("error should mention missing keys")
	}
}

// TestGenerateJWT_InvalidExpiration tests JWT with invalid expiration format
func TestGenerateJWT_InvalidExpiration(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"generate-jwt",
		"-secret", TestSecret,
		"-username", TestUsername,
		"-expiration", "invalid",
	}

	exitCode := Run(args, nil, stdout, stderr)

	if exitCode != 1 {
		t.Errorf("Run() exitCode = %d, want 1", exitCode)
	}
}

// TestGenerateCert_InvalidHosts tests certificate with empty hosts
func TestGenerateCert_InvalidHosts(t *testing.T) {
	tmpDir := t.TempDir()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"generate-cert",
		"-hosts", "",
		"-output-dir", tmpDir,
	}

	exitCode := Run(args, nil, stdout, stderr)

	// Should handle empty hosts gracefully
	if exitCode == 0 {
		t.Logf("Empty hosts handled, exitCode = 0")
	}
}
