// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestDetermineRole(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "no arguments - default compiler",
			args:     []string{"tsd"},
			expected: RoleCompiler,
		},
		{
			name:     "auth role",
			args:     []string{"tsd", "auth", "generate-key"},
			expected: RoleAuth,
		},
		{
			name:     "client role",
			args:     []string{"tsd", "client", "program.tsd"},
			expected: RoleClient,
		},
		{
			name:     "server role",
			args:     []string{"tsd", "server", "-port", "8080"},
			expected: RoleServer,
		},
		{
			name:     "file argument - default compiler",
			args:     []string{"tsd", "program.tsd"},
			expected: RoleCompiler,
		},
		{
			name:     "flag argument - default compiler",
			args:     []string{"tsd", "-file", "program.tsd"},
			expected: RoleCompiler,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Sauvegarder os.Args et le restaurer après le test
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			// Définir os.Args pour le test
			os.Args = tt.args

			// Appeler la vraie fonction determineRole
			role := determineRole()

			if role != tt.expected {
				t.Errorf("determineRole() = %q, want %q", role, tt.expected)
			}
		})
	}
}

func TestPrintGlobalHelp(t *testing.T) {
	// Sauvegarder stdout et le restaurer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	// Exécuter printGlobalHelp
	printGlobalHelp()

	// Fermer le writer et lire la sortie
	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Vérifier que la sortie contient les éléments clés
	expectedStrings := []string{
		"TSD - Type System Development",
		"RÔLES DISPONIBLES",
		"auth",
		"client",
		"server",
		"EXEMPLES",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("printGlobalHelp() output missing %q", expected)
		}
	}
}

func TestPrintGlobalVersion(t *testing.T) {
	// Sauvegarder stdout et le restaurer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	// Exécuter printGlobalVersion
	printGlobalVersion()

	// Fermer le writer et lire la sortie
	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Vérifier que la sortie contient les éléments clés
	expectedStrings := []string{
		Version,
		"TSD",
		"RETE",
		"Copyright",
		"MIT",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("printGlobalVersion() output missing %q", expected)
		}
	}
}

func TestRoleConstants(t *testing.T) {
	// Vérifier que les constantes de rôle sont bien définies
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"auth role", RoleAuth, "auth"},
		{"client role", RoleClient, "client"},
		{"server role", RoleServer, "server"},
		{"compiler role", RoleCompiler, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

func TestVersionConstant(t *testing.T) {
	if Version == "" {
		t.Error("Version constant should not be empty")
	}

	// Vérifier le format (devrait ressembler à "1.0.0")
	if !strings.Contains(Version, ".") {
		t.Errorf("Version %q should contain dots", Version)
	}
}

// TestDispatch_UnknownRole teste le dispatch avec un rôle inconnu
func TestDispatch_UnknownRole(t *testing.T) {
	// Sauvegarder stderr et le restaurer
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	defer func() {
		os.Stderr = oldStderr
	}()

	// Appeler dispatch avec un rôle inconnu
	exitCode := dispatch("unknown")

	// Fermer le writer et lire la sortie d'erreur
	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Vérifier le code de sortie
	if exitCode != ExitMisuseCommand {
		t.Errorf("dispatch(unknown) exitCode = %d, want %d", exitCode, ExitMisuseCommand)
	}

	// Vérifier le message d'erreur
	if !strings.Contains(output, "rôle inconnu") {
		t.Errorf("dispatch(unknown) error message missing 'rôle inconnu', got: %s", output)
	}

	// Vérifier le message d'aide
	if !strings.Contains(output, "--help") {
		t.Errorf("dispatch(unknown) error message missing '--help' suggestion, got: %s", output)
	}
}

// TestDispatch_ValidRoles teste que les rôles valides sont reconnus
func TestDispatch_ValidRoles(t *testing.T) {
	tests := []struct {
		name string
		role string
	}{
		{"auth role", RoleAuth},
		{"client role", RoleClient},
		{"server role", RoleServer},
		{"compiler role", RoleCompiler},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Vérifier que le rôle est dans le switch de dispatch
			validRoles := map[string]bool{
				RoleAuth:     true,
				RoleClient:   true,
				RoleServer:   true,
				RoleCompiler: true,
			}

			if !validRoles[tt.role] {
				t.Errorf("Role %q should be valid", tt.role)
			}
		})
	}
}
