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
	"os"
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

	case "generate-cert":
		return generateCert(args[1:], stdout, stderr)

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

	if *count < 0 {
		fmt.Fprintf(stderr, "❌ Erreur: le nombre de clés ne peut pas être négatif\n")
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
	// Étape 1: Parser les arguments de ligne de commande
	config, _, err := parseValidationFlags(args, stderr)
	if err != nil {
		return 1
	}

	// Étape 2: Mode interactif - lire les inputs manquants
	if config.Interactive {
		if err := readInteractiveInput(config, stdin, stdout, stderr); err != nil {
			fmt.Fprintf(stderr, "Erreur: %v\n", err)
			return 1
		}
	}

	// Étape 3: Valider que tous les paramètres requis sont présents
	if err := validateConfigParameters(config); err != nil {
		fmt.Fprintf(stderr, "Erreur: %v\n", err)
		return 1
	}

	// Étape 4: Créer la configuration auth
	authConfig, err := createAuthConfig(config)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur: %v\n", err)
		return 1
	}

	// Étape 5: Créer le manager d'authentification
	manager, err := auth.NewManager(authConfig)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur initialisation: %v\n", err)
		return 1
	}

	// Étape 6: Valider le token
	result := validateTokenWithManager(manager, config.Token)

	// Étape 7: Formater et afficher le résultat
	output := formatValidationOutput(result, config)
	fmt.Fprint(stdout, output)

	// Retourner le code de sortie approprié
	if result.Valid {
		return 0
	}
	return 1
}

// generateCert génère des certificats TLS auto-signés
func generateCert(args []string, stdout, stderr io.Writer) int {
	// 1. Parser la configuration
	config, err := parseCertFlags(args, stderr)
	if err != nil {
		return 1
	}

	// 2. Créer le répertoire de sortie
	if err := os.MkdirAll(config.outputDir, 0755); err != nil {
		fmt.Fprintf(stderr, "❌ Erreur création répertoire: %v\n", err)
		return 1
	}

	// 3. Générer la clé privée
	privateKey, err := generateECDSAPrivateKey()
	if err != nil {
		fmt.Fprintf(stderr, "❌ Erreur génération clé privée: %v\n", err)
		return 1
	}

	// 4. Créer le template du certificat
	template, err := createCertificateTemplate(config)
	if err != nil {
		fmt.Fprintf(stderr, "❌ Erreur création template: %v\n", err)
		return 1
	}

	// 5. Créer le certificat auto-signé
	certDER, err := createSelfSignedCertificate(template, privateKey)
	if err != nil {
		fmt.Fprintf(stderr, "❌ Erreur création certificat: %v\n", err)
		return 1
	}

	// 6. Écrire les fichiers
	result, err := writeCertificateFiles(config, certDER, privateKey)
	if err != nil {
		fmt.Fprintf(stderr, "❌ Erreur écriture fichiers: %v\n", err)
		return 1
	}

	// 7. Afficher la sortie
	if config.format == "json" {
		formatCertOutputJSON(result, config, stdout)
	} else {
		formatCertOutputText(result, config, stdout)
	}

	return 0
}

// copyFile copie un fichier
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

// printHelp affiche l'aide
func printHelp(w io.Writer) {
	fmt.Fprintln(w, "TSD Auth - Outil de gestion d'authentification pour TSD")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "USAGE:")
	fmt.Fprintln(w, "  tsd auth <commande> [options]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "COMMANDES:")
	fmt.Fprintln(w, "  generate-key     Générer une ou plusieurs clés API")
	fmt.Fprintln(w, "  generate-jwt     Générer un JWT")
	fmt.Fprintln(w, "  generate-cert    Générer des certificats TLS auto-signés")
	fmt.Fprintln(w, "  validate         Valider un token (clé API ou JWT)")
	fmt.Fprintln(w, "  help             Afficher cette aide")
	fmt.Fprintln(w, "  version          Afficher la version")
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
	fmt.Fprintln(w, "  # Générer des certificats TLS pour développement")
	fmt.Fprintln(w, "  tsd auth generate-cert")
	fmt.Fprintln(w, "  tsd auth generate-cert -output-dir ./my-certs -hosts \"localhost,127.0.0.1,192.168.1.100\"")
	fmt.Fprintln(w, "  tsd auth generate-cert -valid-days 730 -org \"My Company\"")
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
