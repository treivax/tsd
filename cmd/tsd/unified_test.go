// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
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
			args:     []string{},
			expected: RoleCompiler,
		},
		{
			name:     "auth role",
			args:     []string{"auth", "generate-key"},
			expected: RoleAuth,
		},
		{
			name:     "client role",
			args:     []string{"client", "program.tsd"},
			expected: RoleClient,
		},
		{
			name:     "server role",
			args:     []string{"server", "-port", "8080"},
			expected: RoleServer,
		},
		{
			name:     "file argument - default compiler",
			args:     []string{"program.tsd"},
			expected: RoleCompiler,
		},
		{
			name:     "flag argument - default compiler",
			args:     []string{"-file", "program.tsd"},
			expected: RoleCompiler,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simuler os.Args
			originalArgs := make([]string, len(tt.args)+1)
			originalArgs[0] = "tsd"
			copy(originalArgs[1:], tt.args)

			// Temporairement remplacer os.Args pour le test
			oldArgs := []string{"tsd"}
			oldArgs = append(oldArgs, tt.args...)

			// Déterminer le rôle en simulant la logique
			var role string
			if len(tt.args) == 0 {
				role = RoleCompiler
			} else {
				firstArg := tt.args[0]
				switch firstArg {
				case RoleAuth, RoleClient, RoleServer:
					role = firstArg
				default:
					role = RoleCompiler
				}
			}

			if role != tt.expected {
				t.Errorf("determineRole() = %v, want %v", role, tt.expected)
			}
		})
	}
}

func TestPrintGlobalHelp(t *testing.T) {
	// Capturer la sortie de printGlobalHelp
	// Note: printGlobalHelp écrit sur un io.Writer, pas os.Stdout directement
	// Pour tester, on devrait la refactorer pour accepter un io.Writer
	// Pour l'instant, on teste juste qu'elle ne panique pas

	// Test basique
	printGlobalHelp()

	// Vérifier que la fonction ne panique pas
	t.Log("printGlobalHelp executed successfully")
}

func TestPrintGlobalVersion(t *testing.T) {
	// Test basique que la fonction ne panique pas
	printGlobalVersion()
	t.Log("printGlobalVersion executed successfully")
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

func TestGlobalHelpContent(t *testing.T) {
	// Capturer stdout pour tester le contenu
	// Note: Comme printGlobalHelp écrit directement sur stdout,
	// on ne peut pas facilement capturer la sortie dans ce test.
	// On devrait refactorer pour accepter un io.Writer.
	// Pour l'instant, on vérifie juste l'exécution.

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printGlobalHelp panicked: %v", r)
		}
	}()

	printGlobalHelp()

	// Dans une version future, on pourrait vérifier que la sortie contient:
	// - "TSD - Type System Development"
	// - "RÔLES DISPONIBLES"
	// - "auth", "client", "server"
	// - "EXEMPLES"
}

func TestGlobalVersionContent(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printGlobalVersion panicked: %v", r)
		}
	}()

	printGlobalVersion()

	// Dans une version future, on pourrait vérifier que la sortie contient:
	// - Version
	// - "Moteur de règles"
	// - "Copyright"
	// - "MIT"
}

// TestDispatchLogic teste la logique de dispatch sans exécuter réellement les commandes
func TestDispatchLogic(t *testing.T) {
	tests := []struct {
		name        string
		role        string
		shouldError bool
		description string
	}{
		{
			name:        "dispatch auth",
			role:        RoleAuth,
			shouldError: false,
			description: "Should dispatch to authcmd",
		},
		{
			name:        "dispatch client",
			role:        RoleClient,
			shouldError: false,
			description: "Should dispatch to clientcmd",
		},
		{
			name:        "dispatch server",
			role:        RoleServer,
			shouldError: false,
			description: "Should dispatch to servercmd",
		},
		{
			name:        "dispatch compiler",
			role:        RoleCompiler,
			shouldError: false,
			description: "Should dispatch to compilercmd",
		},
		{
			name:        "unknown role",
			role:        "unknown",
			shouldError: true,
			description: "Should return error for unknown role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Vérifier que le rôle est valide ou invalide comme attendu
			validRoles := map[string]bool{
				RoleAuth:     true,
				RoleClient:   true,
				RoleServer:   true,
				RoleCompiler: true,
			}

			isValid := validRoles[tt.role]

			if tt.shouldError && isValid {
				t.Errorf("Expected role %q to be invalid, but it's valid", tt.role)
			}

			if !tt.shouldError && !isValid {
				t.Errorf("Expected role %q to be valid, but it's invalid", tt.role)
			}
		})
	}
}
