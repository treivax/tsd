// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/treivax/tsd/auth"
)

const (
	// Version de l'outil
	Version = "1.0.0"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "generate-key":
		exitCode := generateKey(os.Args[2:])
		os.Exit(exitCode)

	case "generate-jwt":
		exitCode := generateJWT(os.Args[2:])
		os.Exit(exitCode)

	case "validate":
		exitCode := validateToken(os.Args[2:])
		os.Exit(exitCode)

	case "help", "-h", "--help":
		printHelp()
		os.Exit(0)

	case "version", "-v", "--version":
		fmt.Printf("tsd-auth version %s\n", Version)
		os.Exit(0)

	default:
		fmt.Fprintf(os.Stderr, "Commande inconnue: %s\n\n", command)
		printHelp()
		os.Exit(1)
	}
}

// generateKey génère une nouvelle clé API
func generateKey(args []string) int {
	fs := flag.NewFlagSet("generate-key", flag.ExitOnError)
	count := fs.Int("count", 1, "Nombre de clés à générer")
	format := fs.String("format", "text", "Format de sortie (text, json)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur: %v\n", err)
		return 1
	}

	keys := make([]string, *count)
	for i := 0; i < *count; i++ {
		key, err := auth.GenerateAuthKey()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur génération clé: %v\n", err)
			return 1
		}
		keys[i] = key
	}

	if *format == "json" {
		output := map[string]interface{}{
			"keys":  keys,
			"count": len(keys),
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println("🔑 Clé(s) API générée(s):")
		fmt.Println("========================")
		for i, key := range keys {
			if *count > 1 {
				fmt.Printf("\nClé %d:\n", i+1)
			}
			fmt.Println(key)
		}
		fmt.Println("\n⚠️  IMPORTANT: Conservez ces clés en lieu sûr!")
		fmt.Println("   Elles ne pourront pas être récupérées ultérieurement.")
	}

	return 0
}

// generateJWT génère un nouveau JWT
func generateJWT(args []string) int {
	fs := flag.NewFlagSet("generate-jwt", flag.ExitOnError)
	secret := fs.String("secret", "", "Secret JWT (requis)")
	username := fs.String("username", "", "Nom d'utilisateur (requis)")
	roles := fs.String("roles", "", "Rôles séparés par des virgules (optionnel)")
	expiration := fs.Duration("expiration", 24*time.Hour, "Durée de validité (ex: 24h, 30m)")
	issuer := fs.String("issuer", "tsd-server", "Émetteur du JWT")
	format := fs.String("format", "text", "Format de sortie (text, json)")
	interactive := fs.Bool("i", false, "Mode interactif (demande le secret)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur: %v\n", err)
		return 1
	}

	// Mode interactif pour le secret
	if *interactive && *secret == "" {
		fmt.Print("Secret JWT: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lecture: %v\n", err)
			return 1
		}
		*secret = strings.TrimSpace(input)
	}

	// Validation
	if *secret == "" {
		fmt.Fprintf(os.Stderr, "Erreur: le secret JWT est requis (utilisez -secret ou -i)\n")
		return 1
	}

	if *username == "" {
		fmt.Fprintf(os.Stderr, "Erreur: le nom d'utilisateur est requis (-username)\n")
		return 1
	}

	// Parser les rôles
	var rolesList []string
	if *roles != "" {
		rolesList = strings.Split(*roles, ",")
		for i, role := range rolesList {
			rolesList[i] = strings.TrimSpace(role)
		}
	}

	// Créer le manager d'authentification
	config := &auth.Config{
		Type:          auth.AuthTypeJWT,
		JWTSecret:     *secret,
		JWTExpiration: *expiration,
		JWTIssuer:     *issuer,
	}

	manager, err := auth.NewManager(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur initialisation: %v\n", err)
		return 1
	}

	// Générer le JWT
	token, err := manager.GenerateJWT(*username, rolesList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur génération JWT: %v\n", err)
		return 1
	}

	if *format == "json" {
		output := map[string]interface{}{
			"token":      token,
			"username":   *username,
			"roles":      rolesList,
			"expiration": expiration.String(),
			"issuer":     *issuer,
			"expires_at": time.Now().Add(*expiration).Format(time.RFC3339),
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println("🎫 JWT généré avec succès:")
		fmt.Println("==========================")
		fmt.Printf("Token: %s\n\n", token)
		fmt.Printf("Utilisateur: %s\n", *username)
		if len(rolesList) > 0 {
			fmt.Printf("Rôles: %s\n", strings.Join(rolesList, ", "))
		}
		fmt.Printf("Expire dans: %s\n", expiration.String())
		fmt.Printf("Expire le: %s\n", time.Now().Add(*expiration).Format(time.RFC3339))
		fmt.Printf("Émetteur: %s\n", *issuer)
		fmt.Println("\n⚠️  IMPORTANT: Conservez ce token en lieu sûr!")
	}

	return 0
}

// validateToken valide un token
func validateToken(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	token := fs.String("token", "", "Token à valider (requis)")
	authType := fs.String("type", "", "Type d'auth: key ou jwt (requis)")
	secret := fs.String("secret", "", "Secret JWT (requis si type=jwt)")
	keys := fs.String("keys", "", "Clés API valides séparées par des virgules (requis si type=key)")
	interactive := fs.Bool("i", false, "Mode interactif")
	format := fs.String("format", "text", "Format de sortie (text, json)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur: %v\n", err)
		return 1
	}

	// Mode interactif
	if *interactive {
		reader := bufio.NewReader(os.Stdin)

		if *token == "" {
			fmt.Print("Token: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erreur lecture: %v\n", err)
				return 1
			}
			*token = strings.TrimSpace(input)
		}

		if *authType == "" {
			fmt.Print("Type d'authentification (key/jwt): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erreur lecture: %v\n", err)
				return 1
			}
			*authType = strings.TrimSpace(input)
		}

		if *authType == "jwt" && *secret == "" {
			fmt.Print("Secret JWT: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erreur lecture: %v\n", err)
				return 1
			}
			*secret = strings.TrimSpace(input)
		}

		if *authType == "key" && *keys == "" {
			fmt.Print("Clés API (séparées par des virgules): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erreur lecture: %v\n", err)
				return 1
			}
			*keys = strings.TrimSpace(input)
		}
	}

	// Validation des arguments
	if *token == "" {
		fmt.Fprintf(os.Stderr, "Erreur: le token est requis (-token)\n")
		return 1
	}

	if *authType == "" {
		fmt.Fprintf(os.Stderr, "Erreur: le type d'authentification est requis (-type key|jwt)\n")
		return 1
	}

	// Créer la config selon le type
	var config *auth.Config
	switch *authType {
	case "key":
		if *keys == "" {
			fmt.Fprintf(os.Stderr, "Erreur: les clés API sont requises pour type=key (-keys)\n")
			return 1
		}
		keysList := strings.Split(*keys, ",")
		for i, key := range keysList {
			keysList[i] = strings.TrimSpace(key)
		}
		config = &auth.Config{
			Type:     auth.AuthTypeKey,
			AuthKeys: keysList,
		}

	case "jwt":
		if *secret == "" {
			fmt.Fprintf(os.Stderr, "Erreur: le secret JWT est requis pour type=jwt (-secret)\n")
			return 1
		}
		config = &auth.Config{
			Type:      auth.AuthTypeJWT,
			JWTSecret: *secret,
		}

	default:
		fmt.Fprintf(os.Stderr, "Erreur: type invalide '%s' (doit être 'key' ou 'jwt')\n", *authType)
		return 1
	}

	// Créer le manager
	manager, err := auth.NewManager(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur initialisation: %v\n", err)
		return 1
	}

	// Valider le token
	info, err := manager.GetTokenInfo(*token)

	if *format == "json" {
		output := map[string]interface{}{
			"valid": err == nil && info.Valid,
			"type":  *authType,
		}
		if err != nil {
			output["error"] = err.Error()
		}
		if info != nil {
			output["username"] = info.Username
			output["roles"] = info.Roles
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
	} else {
		if err == nil && info.Valid {
			fmt.Println("✅ Token valide")
			fmt.Printf("Type: %s\n", *authType)
			if *authType == "jwt" && info.Username != "" {
				fmt.Printf("Utilisateur: %s\n", info.Username)
				if len(info.Roles) > 0 {
					fmt.Printf("Rôles: %s\n", strings.Join(info.Roles, ", "))
				}
			}
			return 0
		} else {
			fmt.Println("❌ Token invalide")
			if err != nil {
				fmt.Printf("Erreur: %v\n", err)
			}
			return 1
		}
	}

	return 0
}

// printHelp affiche l'aide
func printHelp() {
	fmt.Println("TSD Auth - Outil de gestion d'authentification pour TSD")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("  tsd-auth <commande> [options]")
	fmt.Println("")
	fmt.Println("COMMANDES:")
	fmt.Println("  generate-key    Générer une ou plusieurs clés API")
	fmt.Println("  generate-jwt    Générer un JWT")
	fmt.Println("  validate        Valider un token (clé API ou JWT)")
	fmt.Println("  help            Afficher cette aide")
	fmt.Println("  version         Afficher la version")
	fmt.Println("")
	fmt.Println("EXEMPLES:")
	fmt.Println("")
	fmt.Println("  # Générer une clé API")
	fmt.Println("  tsd-auth generate-key")
	fmt.Println("")
	fmt.Println("  # Générer plusieurs clés API")
	fmt.Println("  tsd-auth generate-key -count 3")
	fmt.Println("")
	fmt.Println("  # Générer un JWT")
	fmt.Println("  tsd-auth generate-jwt -secret \"mon-secret-super-securise-de-32-chars\" -username alice")
	fmt.Println("")
	fmt.Println("  # Générer un JWT avec rôles et expiration personnalisée")
	fmt.Println("  tsd-auth generate-jwt -secret \"mon-secret\" -username bob -roles \"admin,user\" -expiration 48h")
	fmt.Println("")
	fmt.Println("  # Générer un JWT en mode interactif (ne pas exposer le secret)")
	fmt.Println("  tsd-auth generate-jwt -i -username alice")
	fmt.Println("")
	fmt.Println("  # Valider une clé API")
	fmt.Println("  tsd-auth validate -type key -token \"ma-cle-api\" -keys \"cle1,cle2,cle3\"")
	fmt.Println("")
	fmt.Println("  # Valider un JWT")
	fmt.Println("  tsd-auth validate -type jwt -token \"eyJhbG...\" -secret \"mon-secret\"")
	fmt.Println("")
	fmt.Println("  # Valider en mode interactif")
	fmt.Println("  tsd-auth validate -i")
	fmt.Println("")
	fmt.Println("OPTIONS COMMUNES:")
	fmt.Println("  -format text|json   Format de sortie (défaut: text)")
	fmt.Println("  -i                  Mode interactif")
	fmt.Println("")
	fmt.Println("Pour plus de détails sur une commande:")
	fmt.Println("  tsd-auth <commande> -h")
}
