// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package authcmd

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/treivax/tsd/auth"
)

const (
	// Version de l'outil auth
	Version = "1.0.0"
)

// Run exécute la commande auth avec les arguments donnés
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printHelp(stdout)
		return 1
	}

	command := args[0]

	switch command {
	case "generate-key":
		return generateKey(args[1:], stdout, stderr)

	case "generate-jwt":
		return generateJWT(args[1:], stdin, stdout, stderr)

	case "validate":
		return validateToken(args[1:], stdin, stdout, stderr)

	case "help", "-h", "--help":
		printHelp(stdout)
		return 0

	case "version", "-v", "--version":
		fmt.Fprintf(stdout, "tsd auth version %s\n", Version)
		return 0

	default:
		fmt.Fprintf(stderr, "Commande inconnue: %s\n\n", command)
		printHelp(stderr)
		return 1
	}
}

// generateKey génère une nouvelle clé API
func generateKey(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("generate-key", flag.ContinueOnError)
	fs.SetOutput(stderr)
	count := fs.Int("count", 1, "Nombre de clés à générer")
	format := fs.String("format", "text", "Format de sortie (text, json)")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	keys := make([]string, *count)
	for i := 0; i < *count; i++ {
		key, err := auth.GenerateAuthKey()
		if err != nil {
			fmt.Fprintf(stderr, "Erreur génération clé: %v\n", err)
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
		fmt.Fprintln(stdout, string(data))
	} else {
		fmt.Fprintln(stdout, "🔑 Clé(s) API générée(s):")
		fmt.Fprintln(stdout, "========================")
		for i, key := range keys {
			if *count > 1 {
				fmt.Fprintf(stdout, "\nClé %d:\n", i+1)
			}
			fmt.Fprintln(stdout, key)
		}
		fmt.Fprintln(stdout, "\n⚠️  IMPORTANT: Conservez ces clés en lieu sûr!")
		fmt.Fprintln(stdout, "   Elles ne pourront pas être récupérées ultérieurement.")
	}

	return 0
}

// generateJWT génère un nouveau JWT
func generateJWT(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("generate-jwt", flag.ContinueOnError)
	fs.SetOutput(stderr)
	secret := fs.String("secret", "", "Secret JWT (requis)")
	username := fs.String("username", "", "Nom d'utilisateur (requis)")
	roles := fs.String("roles", "", "Rôles séparés par des virgules (optionnel)")
	expiration := fs.Duration("expiration", 24*time.Hour, "Durée de validité (ex: 24h, 30m)")
	issuer := fs.String("issuer", "tsd-server", "Émetteur du JWT")
	format := fs.String("format", "text", "Format de sortie (text, json)")
	interactive := fs.Bool("i", false, "Mode interactif (demande le secret)")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	// Mode interactif pour le secret
	if *interactive && *secret == "" {
		fmt.Fprint(stdout, "Secret JWT: ")
		reader := bufio.NewReader(stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(stderr, "Erreur lecture: %v\n", err)
			return 1
		}
		*secret = strings.TrimSpace(input)
	}

	// Validation
	if *secret == "" {
		fmt.Fprintln(stderr, "Erreur: le secret JWT est requis (utilisez -secret ou -i)")
		return 1
	}

	if *username == "" {
		fmt.Fprintln(stderr, "Erreur: le nom d'utilisateur est requis (-username)")
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
		fmt.Fprintf(stderr, "Erreur initialisation: %v\n", err)
		return 1
	}

	// Générer le JWT
	token, err := manager.GenerateJWT(*username, rolesList)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur génération JWT: %v\n", err)
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
		fmt.Fprintln(stdout, string(data))
	} else {
		fmt.Fprintln(stdout, "🎫 JWT généré avec succès:")
		fmt.Fprintln(stdout, "==========================")
		fmt.Fprintf(stdout, "Token: %s\n\n", token)
		fmt.Fprintf(stdout, "Utilisateur: %s\n", *username)
		if len(rolesList) > 0 {
			fmt.Fprintf(stdout, "Rôles: %s\n", strings.Join(rolesList, ", "))
		}
		fmt.Fprintf(stdout, "Expire dans: %s\n", expiration.String())
		fmt.Fprintf(stdout, "Expire le: %s\n", time.Now().Add(*expiration).Format(time.RFC3339))
		fmt.Fprintf(stdout, "Émetteur: %s\n", *issuer)
		fmt.Fprintln(stdout, "\n⚠️  IMPORTANT: Conservez ce token en lieu sûr!")
	}

	return 0
}

// validateToken valide un token
func validateToken(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	token := fs.String("token", "", "Token à valider (requis)")
	authType := fs.String("type", "", "Type d'auth: key ou jwt (requis)")
	secret := fs.String("secret", "", "Secret JWT (requis si type=jwt)")
	keys := fs.String("keys", "", "Clés API valides séparées par des virgules (requis si type=key)")
	interactive := fs.Bool("i", false, "Mode interactif")
	format := fs.String("format", "text", "Format de sortie (text, json)")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	// Mode interactif
	if *interactive {
		reader := bufio.NewReader(stdin)

		if *token == "" {
			fmt.Fprint(stdout, "Token: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(stderr, "Erreur lecture: %v\n", err)
				return 1
			}
			*token = strings.TrimSpace(input)
		}

		if *authType == "" {
			fmt.Fprint(stdout, "Type d'authentification (key/jwt): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(stderr, "Erreur lecture: %v\n", err)
				return 1
			}
			*authType = strings.TrimSpace(input)
		}

		if *authType == "jwt" && *secret == "" {
			fmt.Fprint(stdout, "Secret JWT: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(stderr, "Erreur lecture: %v\n", err)
				return 1
			}
			*secret = strings.TrimSpace(input)
		}

		if *authType == "key" && *keys == "" {
			fmt.Fprint(stdout, "Clés API (séparées par des virgules): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(stderr, "Erreur lecture: %v\n", err)
				return 1
			}
			*keys = strings.TrimSpace(input)
		}
	}

	// Validation des arguments
	if *token == "" {
		fmt.Fprintln(stderr, "Erreur: le token est requis (-token)")
		return 1
	}

	if *authType == "" {
		fmt.Fprintln(stderr, "Erreur: le type d'authentification est requis (-type key|jwt)")
		return 1
	}

	// Créer la config selon le type
	var config *auth.Config
	switch *authType {
	case "key":
		if *keys == "" {
			fmt.Fprintln(stderr, "Erreur: les clés API sont requises pour type=key (-keys)")
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
			fmt.Fprintln(stderr, "Erreur: le secret JWT est requis pour type=jwt (-secret)")
			return 1
		}
		config = &auth.Config{
			Type:      auth.AuthTypeJWT,
			JWTSecret: *secret,
		}

	default:
		fmt.Fprintf(stderr, "Erreur: type invalide '%s' (doit être 'key' ou 'jwt')\n", *authType)
		return 1
	}

	// Créer le manager
	manager, err := auth.NewManager(config)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur initialisation: %v\n", err)
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
		fmt.Fprintln(stdout, string(data))
	} else {
		if err == nil && info.Valid {
			fmt.Fprintln(stdout, "✅ Token valide")
			fmt.Fprintf(stdout, "Type: %s\n", *authType)
			if *authType == "jwt" && info.Username != "" {
				fmt.Fprintf(stdout, "Utilisateur: %s\n", info.Username)
				if len(info.Roles) > 0 {
					fmt.Fprintf(stdout, "Rôles: %s\n", strings.Join(info.Roles, ", "))
				}
			}
			return 0
		}
		fmt.Fprintln(stdout, "❌ Token invalide")
		if err != nil {
			fmt.Fprintf(stdout, "Erreur: %v\n", err)
		}
		return 1
	}

	return 0
}

// printHelp affiche l'aide
func printHelp(w io.Writer) {
	fmt.Fprintln(w, "TSD Auth - Outil de gestion d'authentification pour TSD")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "USAGE:")
	fmt.Fprintln(w, "  tsd auth <commande> [options]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "COMMANDES:")
	fmt.Fprintln(w, "  generate-key    Générer une ou plusieurs clés API")
	fmt.Fprintln(w, "  generate-jwt    Générer un JWT")
	fmt.Fprintln(w, "  validate        Valider un token (clé API ou JWT)")
	fmt.Fprintln(w, "  help            Afficher cette aide")
	fmt.Fprintln(w, "  version         Afficher la version")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "EXEMPLES:")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Générer une clé API")
	fmt.Fprintln(w, "  tsd auth generate-key")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Générer plusieurs clés API")
	fmt.Fprintln(w, "  tsd auth generate-key -count 3")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Générer un JWT")
	fmt.Fprintln(w, "  tsd auth generate-jwt -secret \"mon-secret-super-securise-de-32-chars\" -username alice")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Générer un JWT avec rôles et expiration personnalisée")
	fmt.Fprintln(w, "  tsd auth generate-jwt -secret \"mon-secret\" -username bob -roles \"admin,user\" -expiration 48h")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Générer un JWT en mode interactif (ne pas exposer le secret)")
	fmt.Fprintln(w, "  tsd auth generate-jwt -i -username alice")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Valider une clé API")
	fmt.Fprintln(w, "  tsd auth validate -type key -token \"ma-cle-api\" -keys \"cle1,cle2,cle3\"")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Valider un JWT")
	fmt.Fprintln(w, "  tsd auth validate -type jwt -token \"eyJhbG...\" -secret \"mon-secret\"")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Valider en mode interactif")
	fmt.Fprintln(w, "  tsd auth validate -i")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "OPTIONS COMMUNES:")
	fmt.Fprintln(w, "  -format text|json   Format de sortie (défaut: text)")
	fmt.Fprintln(w, "  -i                  Mode interactif")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Pour plus de détails sur une commande:")
	fmt.Fprintln(w, "  tsd auth <commande> -h")
}
