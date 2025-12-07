// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package authcmd

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"path/filepath"
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
	fs := flag.NewFlagSet("generate-cert", flag.ContinueOnError)
	fs.SetOutput(stderr)
	outputDir := fs.String("output-dir", "./certs", "Répertoire de sortie pour les certificats")
	hosts := fs.String("hosts", "localhost,127.0.0.1", "Hôtes/IPs séparés par des virgules")
	validDays := fs.Int("valid-days", 365, "Durée de validité en jours")
	org := fs.String("org", "TSD Development", "Nom de l'organisation")
	format := fs.String("format", "text", "Format de sortie (text, json)")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	// Créer le répertoire de sortie si nécessaire
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(stderr, "❌ Erreur création répertoire: %v\n", err)
		return 1
	}

	// Parser les hôtes
	hostList := strings.Split(*hosts, ",")
	for i, h := range hostList {
		hostList[i] = strings.TrimSpace(h)
	}

	// Générer la clé privée ECDSA
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Fprintf(stderr, "❌ Erreur génération clé privée: %v\n", err)
		return 1
	}

	// Préparer le template du certificat
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Duration(*validDays) * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		fmt.Fprintf(stderr, "❌ Erreur génération numéro série: %v\n", err)
		return 1
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{*org},
			CommonName:   hostList[0],
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Ajouter les hôtes au certificat
	for _, h := range hostList {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	// Créer le certificat auto-signé
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		fmt.Fprintf(stderr, "❌ Erreur création certificat: %v\n", err)
		return 1
	}

	// Sauvegarder le certificat
	certPath := filepath.Join(*outputDir, "server.crt")
	certOut, err := os.Create(certPath)
	if err != nil {
		fmt.Fprintf(stderr, "❌ Erreur création fichier certificat: %v\n", err)
		return 1
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		certOut.Close()
		fmt.Fprintf(stderr, "❌ Erreur encodage certificat: %v\n", err)
		return 1
	}
	certOut.Close()

	// Sauvegarder la clé privée
	keyPath := filepath.Join(*outputDir, "server.key")
	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Fprintf(stderr, "❌ Erreur création fichier clé: %v\n", err)
		return 1
	}

	privBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		keyOut.Close()
		fmt.Fprintf(stderr, "❌ Erreur marshalling clé: %v\n", err)
		return 1
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}); err != nil {
		keyOut.Close()
		fmt.Fprintf(stderr, "❌ Erreur encodage clé: %v\n", err)
		return 1
	}
	keyOut.Close()

	// Créer aussi une copie du certificat comme CA (pour les clients)
	caPath := filepath.Join(*outputDir, "ca.crt")
	if err := copyFile(certPath, caPath); err != nil {
		fmt.Fprintf(stderr, "⚠️  Avertissement: impossible de créer ca.crt: %v\n", err)
	}

	if *format == "json" {
		output := map[string]interface{}{
			"success":      true,
			"cert_path":    certPath,
			"key_path":     keyPath,
			"ca_path":      caPath,
			"hosts":        hostList,
			"valid_days":   *validDays,
			"not_before":   notBefore.Format(time.RFC3339),
			"not_after":    notAfter.Format(time.RFC3339),
			"organization": *org,
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Fprintln(stdout, string(data))
	} else {
		fmt.Fprintln(stdout, "🔐 Certificats TLS générés avec succès!")
		fmt.Fprintln(stdout, "=====================================")
		fmt.Fprintf(stdout, "\n📁 Répertoire: %s\n\n", *outputDir)
		fmt.Fprintf(stdout, "📄 Fichiers générés:\n")
		fmt.Fprintf(stdout, "   - %s (certificat serveur)\n", certPath)
		fmt.Fprintf(stdout, "   - %s (clé privée serveur)\n", keyPath)
		fmt.Fprintf(stdout, "   - %s (certificat CA pour clients)\n\n", caPath)
		fmt.Fprintf(stdout, "🏷️  Hôtes autorisés: %s\n", strings.Join(hostList, ", "))
		fmt.Fprintf(stdout, "📅 Valide de %s à %s\n", notBefore.Format("2006-01-02"), notAfter.Format("2006-01-02"))
		fmt.Fprintf(stdout, "🏢 Organisation: %s\n\n", *org)
		fmt.Fprintln(stdout, "⚠️  IMPORTANT:")
		fmt.Fprintf(stdout, "   - La clé privée (%s) doit rester SECRÈTE\n", keyPath)
		fmt.Fprintln(stdout, "   - Ne JAMAIS committer les certificats dans Git")
		fmt.Fprintln(stdout, "   - Ces certificats sont auto-signés (pour développement)")
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "📝 Utilisation:")
		fmt.Fprintf(stdout, "   Serveur: tsd server --tls-cert %s --tls-key %s\n", certPath, keyPath)
		fmt.Fprintf(stdout, "   Client:  tsd client --server https://localhost:8080 --tls-ca %s\n", caPath)
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "💡 Pour production, utilisez des certificats signés par une CA reconnue (Let's Encrypt, etc.)")
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
