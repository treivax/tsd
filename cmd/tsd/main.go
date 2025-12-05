// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"os"

	"github.com/treivax/tsd/internal/authcmd"
	"github.com/treivax/tsd/internal/clientcmd"
	"github.com/treivax/tsd/internal/compilercmd"
	"github.com/treivax/tsd/internal/servercmd"
)

const (
	// Version du binaire unifié TSD
	Version = "1.0.0"

	// Noms des rôles disponibles
	RoleAuth     = "auth"
	RoleClient   = "client"
	RoleServer   = "server"
	RoleCompiler = "" // Rôle par défaut (compilateur)
)

func main() {
	// Gérer --version et --help au niveau global avant de dispatcher
	if len(os.Args) > 1 {
		firstArg := os.Args[1]

		// Aide globale
		if firstArg == "--help" || firstArg == "-h" {
			printGlobalHelp()
			os.Exit(0)
		}

		// Version globale
		if firstArg == "--version" || firstArg == "-v" {
			printGlobalVersion()
			os.Exit(0)
		}
	}

	// Déterminer le rôle selon le premier argument
	role := determineRole()

	// Dispatcher selon le rôle
	exitCode := dispatch(role)
	os.Exit(exitCode)
}

// determineRole détermine le rôle à exécuter selon les arguments
func determineRole() string {
	if len(os.Args) < 2 {
		// Pas d'arguments: comportement par défaut (compilateur)
		return RoleCompiler
	}

	firstArg := os.Args[1]

	// Vérifier si le premier argument est un rôle connu
	switch firstArg {
	case RoleAuth, RoleClient, RoleServer:
		return firstArg
	default:
		// Pas un rôle connu: comportement par défaut (compilateur)
		return RoleCompiler
	}
}

// dispatch exécute le rôle approprié avec les arguments restants
func dispatch(role string) int {
	switch role {
	case RoleAuth:
		// Exécuter la commande auth avec les arguments restants (sans le premier "auth")
		return authcmd.Run(os.Args[2:], os.Stdin, os.Stdout, os.Stderr)

	case RoleClient:
		// Exécuter la commande client
		return clientcmd.Run(os.Args[2:], os.Stdin, os.Stdout, os.Stderr)

	case RoleServer:
		// Exécuter la commande server
		return servercmd.Run(os.Args[2:], os.Stdin, os.Stdout, os.Stderr)

	case RoleCompiler:
		// Exécuter le compilateur/runner avec tous les arguments
		return compilercmd.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)

	default:
		fmt.Fprintf(os.Stderr, "Erreur: rôle inconnu '%s'\n", role)
		return 1
	}
}

// printGlobalHelp affiche l'aide globale du binaire unifié
func printGlobalHelp() {
	fmt.Println("TSD - Type System Development")
	fmt.Println("Moteur de règles basé sur l'algorithme RETE avec système d'authentification")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("  tsd [role] [options]")
	fmt.Println("")
	fmt.Println("RÔLES DISPONIBLES:")
	fmt.Println("  (aucun)         Compilateur/Runner TSD (comportement par défaut)")
	fmt.Println("  auth            Gestion de l'authentification (clés API, JWT)")
	fmt.Println("  client          Client HTTP pour communiquer avec tsd-server")
	fmt.Println("  server          Serveur HTTP TSD")
	fmt.Println("")
	fmt.Println("OPTIONS GLOBALES:")
	fmt.Println("  --help, -h      Afficher cette aide")
	fmt.Println("  --version, -v   Afficher la version")
	fmt.Println("")
	fmt.Println("EXEMPLES:")
	fmt.Println("")
	fmt.Println("  # Compiler/exécuter un programme TSD (comportement par défaut)")
	fmt.Println("  tsd program.tsd")
	fmt.Println("  tsd -file program.tsd -v")
	fmt.Println("")
	fmt.Println("  # Gestion de l'authentification")
	fmt.Println("  tsd auth generate-key")
	fmt.Println("  tsd auth generate-jwt -secret \"mon-secret\" -username alice")
	fmt.Println("  tsd auth validate -type jwt -token \"eyJhbG...\" -secret \"mon-secret\"")
	fmt.Println("")
	fmt.Println("  # Client HTTP")
	fmt.Println("  tsd client program.tsd -server http://localhost:8080")
	fmt.Println("  tsd client -health")
	fmt.Println("")
	fmt.Println("  # Serveur HTTP")
	fmt.Println("  tsd server -port 8080")
	fmt.Println("  tsd server -auth jwt -jwt-secret \"mon-secret\"")
	fmt.Println("")
	fmt.Println("AIDE SPÉCIFIQUE À UN RÔLE:")
	fmt.Println("  tsd auth --help")
	fmt.Println("  tsd client --help")
	fmt.Println("  tsd server --help")
	fmt.Println("  tsd --help          (aide du compilateur)")
	fmt.Println("")
	fmt.Println("NOTES:")
	fmt.Println("  Le binaire unique 'tsd' remplace les anciens binaires séparés.")
	fmt.Println("")
	fmt.Println("DOCUMENTATION:")
	fmt.Println("  Consultez README.md et docs/ pour plus d'informations.")
}

// printGlobalVersion affiche la version globale du binaire unifié
func printGlobalVersion() {
	fmt.Printf("TSD (Type System Development) v%s\n", Version)
	fmt.Println("Moteur de règles basé sur l'algorithme RETE")
	fmt.Println("")
	fmt.Println("Composants:")
	fmt.Println("  - Compilateur/Runner TSD")
	fmt.Println("  - Gestion d'authentification (Auth Key + JWT)")
	fmt.Println("  - Client HTTP")
	fmt.Println("  - Serveur HTTP")
	fmt.Println("")
	fmt.Println("Copyright (c) 2025 TSD Contributors")
	fmt.Println("Licence: MIT")
}
