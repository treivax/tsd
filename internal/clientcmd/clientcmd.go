// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package clientcmd

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/treivax/tsd/tsdio"
)

const (
	// DefaultServerURL est l'URL par défaut du serveur
	DefaultServerURL = "https://localhost:8080"

	// DefaultTimeout est le timeout par défaut des requêtes
	DefaultTimeout = 30 * time.Second

	// DefaultCAFile est le fichier CA par défaut
	DefaultCAFile = "./certs/ca.crt"
)

// Config contient la configuration du client
type Config struct {
	ServerURL  string
	File       string
	Text       string
	UseStdin   bool
	Verbose    bool
	Format     string
	Timeout    time.Duration
	ShowHelp   bool
	ShowHealth bool
	AuthToken  string
	AuthType   string
	TLSCAFile  string
	Insecure   bool
}

// Client représente le client HTTP TSD
type Client struct {
	config     *Config
	httpClient *http.Client
	tlsConfig  *tls.Config
}

// Run exécute le client avec les arguments donnés et retourne un code de sortie
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	config, err := parseFlags(args)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur: %v\n", err)
		return 1
	}

	if config.ShowHelp {
		printHelp(stdout)
		return 0
	}

	if config.ShowHealth {
		return runHealthCheck(config, stdout, stderr)
	}

	if err := validateConfig(config); err != nil {
		fmt.Fprintf(stderr, "Erreur: %v\n\n", err)
		printHelp(stderr)
		return 1
	}

	// Lire le code source TSD
	source, sourceName, err := readSource(config, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur lecture source: %v\n", err)
		return 1
	}

	// Créer le client
	client := NewClient(config)

	// Envoyer la requête
	response, err := client.Execute(source, sourceName)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur communication serveur: %v\n", err)
		return 1
	}

	// Afficher les résultats
	if err := printResults(config, response, stdout, stderr); err != nil {
		fmt.Fprintf(stderr, "Erreur affichage résultats: %v\n", err)
		return 1
	}

	if !response.Success {
		return 1
	}

	return 0
}

// parseFlags parse les arguments de ligne de commande
func parseFlags(args []string) (*Config, error) {
	config := &Config{}
	flagSet := flag.NewFlagSet("tsd-client", flag.ContinueOnError)

	flagSet.StringVar(&config.ServerURL, "server", DefaultServerURL, "URL du serveur TSD")
	flagSet.StringVar(&config.File, "file", "", "Fichier TSD (.tsd)")
	flagSet.StringVar(&config.Text, "text", "", "Code TSD directement")
	flagSet.BoolVar(&config.UseStdin, "stdin", false, "Lire depuis stdin")
	flagSet.BoolVar(&config.Verbose, "v", false, "Mode verbeux")
	flagSet.StringVar(&config.Format, "format", "text", "Format de sortie (text, json)")
	flagSet.DurationVar(&config.Timeout, "timeout", DefaultTimeout, "Timeout des requêtes")
	flagSet.BoolVar(&config.ShowHelp, "help", false, "Afficher l'aide")
	flagSet.BoolVar(&config.ShowHealth, "health", false, "Vérifier la santé du serveur")
	flagSet.StringVar(&config.AuthToken, "token", "", "Token d'authentification (clé API ou JWT)")
	flagSet.StringVar(&config.AuthType, "auth-type", "", "Type d'authentification: key ou jwt (optionnel)")

	// TLS
	defaultCAPath := DefaultCAFile
	flagSet.StringVar(&config.TLSCAFile, "tls-ca", defaultCAPath, "Chemin vers le certificat CA pour vérifier le serveur")
	flagSet.BoolVar(&config.Insecure, "insecure", false, "Désactiver la vérification TLS (développement uniquement)")

	if err := flagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			config.ShowHelp = true
			return config, nil
		}
		return nil, err
	}

	// Gérer l'argument positionnel comme fichier
	if config.File == "" && len(flagSet.Args()) > 0 {
		config.File = flagSet.Args()[0]
	}

	// Récupérer le token depuis la variable d'environnement si non fourni
	if config.AuthToken == "" {
		config.AuthToken = os.Getenv("TSD_AUTH_TOKEN")
	}

	// Variables d'environnement pour TLS
	if envCA := os.Getenv("TSD_TLS_CA"); envCA != "" {
		config.TLSCAFile = envCA
	}
	if os.Getenv("TSD_CLIENT_INSECURE") == "true" {
		config.Insecure = true
	}

	return config, nil
}

// validateConfig valide la configuration
func validateConfig(config *Config) error {
	sourcesCount := 0
	if config.File != "" {
		sourcesCount++
	}
	if config.Text != "" {
		sourcesCount++
	}
	if config.UseStdin {
		sourcesCount++
	}

	if sourcesCount == 0 {
		return fmt.Errorf("aucune source spécifiée (-file, -text ou -stdin)")
	}

	if sourcesCount > 1 {
		return fmt.Errorf("une seule source autorisée (-file, -text ou -stdin)")
	}

	if config.Format != "text" && config.Format != "json" {
		return fmt.Errorf("format invalide: %s (doit être 'text' ou 'json')", config.Format)
	}

	return nil
}

// readSource lit le code source TSD
func readSource(config *Config, stdin io.Reader) (string, string, error) {
	if config.UseStdin {
		content, err := io.ReadAll(stdin)
		if err != nil {
			return "", "", fmt.Errorf("lecture stdin: %w", err)
		}
		return string(content), "<stdin>", nil
	}

	if config.Text != "" {
		return config.Text, "<text>", nil
	}

	if config.File != "" {
		if _, err := os.Stat(config.File); os.IsNotExist(err) {
			return "", "", fmt.Errorf("fichier non trouvé: %s", config.File)
		}

		content, err := os.ReadFile(config.File)
		if err != nil {
			return "", "", fmt.Errorf("lecture fichier: %w", err)
		}
		return string(content), config.File, nil
	}

	return "", "", fmt.Errorf("aucune source spécifiée")
}

// NewClient crée un nouveau client TSD
func NewClient(config *Config) *Client {
	// Configurer TLS
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Si mode insecure, désactiver la vérification
	if config.Insecure {
		tlsConfig.InsecureSkipVerify = true
	} else {
		// Charger le CA si fourni et fichier existe
		if config.TLSCAFile != "" {
			if _, err := os.Stat(config.TLSCAFile); err == nil {
				caCert, err := os.ReadFile(config.TLSCAFile)
				if err == nil {
					caCertPool := x509.NewCertPool()
					if caCertPool.AppendCertsFromPEM(caCert) {
						tlsConfig.RootCAs = caCertPool
					}
				}
			}
		}
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return &Client{
		config:    config,
		tlsConfig: tlsConfig,
		httpClient: &http.Client{
			Timeout:   config.Timeout,
			Transport: transport,
		},
	}
}

// Execute envoie une requête d'exécution au serveur
func (c *Client) Execute(source, sourceName string) (*tsdio.ExecuteResponse, error) {
	// Créer la requête
	req := tsdio.ExecuteRequest{
		Source:     source,
		SourceName: sourceName,
		Verbose:    c.config.Verbose,
	}

	// Encoder en JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encodage JSON: %w", err)
	}

	// Créer la requête HTTP
	url := c.config.ServerURL + "/api/v1/execute"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("création requête: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Ajouter le token d'authentification si fourni
	if c.config.AuthToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.config.AuthToken)
	}

	// Envoyer la requête
	if c.config.Verbose {
		fmt.Printf("📤 Envoi requête à %s...\n", url)
		if c.config.AuthToken != "" {
			fmt.Printf("🔒 Authentification: activée\n")
		}
		if c.config.Insecure {
			fmt.Printf("⚠️  TLS: vérification désactivée (mode insecure)\n")
		}
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("envoi requête: %w", err)
	}
	defer resp.Body.Close()

	// Lire la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("lecture réponse: %w", err)
	}

	// Décoder la réponse
	var response tsdio.ExecuteResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("décodage JSON: %w (body: %s)", err, string(body))
	}

	return &response, nil
}

// HealthCheck vérifie la santé du serveur
func (c *Client) HealthCheck() (*tsdio.HealthResponse, error) {
	url := c.config.ServerURL + "/health"
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("requête health: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("lecture réponse: %w", err)
	}

	var health tsdio.HealthResponse
	if err := json.Unmarshal(body, &health); err != nil {
		return nil, fmt.Errorf("décodage JSON: %w", err)
	}

	return &health, nil
}

// runHealthCheck exécute un health check
func runHealthCheck(config *Config, stdout, stderr io.Writer) int {
	client := NewClient(config)

	health, err := client.HealthCheck()
	if err != nil {
		fmt.Fprintf(stderr, "❌ Erreur health check: %v\n", err)
		return 1
	}

	if config.Format == "json" {
		data, _ := json.MarshalIndent(health, "", "  ")
		fmt.Fprintln(stdout, string(data))
	} else {
		fmt.Fprintf(stdout, "✅ Serveur TSD: %s\n", health.Status)
		fmt.Fprintf(stdout, "📊 Version: %s\n", health.Version)
		fmt.Fprintf(stdout, "⏱️  Uptime: %ds\n", health.UptimeSeconds)
		fmt.Fprintf(stdout, "🕐 Timestamp: %s\n", health.Timestamp.Format(time.RFC3339))
	}

	return 0
}

// printResults affiche les résultats
func printResults(config *Config, response *tsdio.ExecuteResponse, stdout, stderr io.Writer) error {
	if config.Format == "json" {
		return printJSON(response, stdout)
	}

	return printText(config, response, stdout, stderr)
}

// printJSON affiche les résultats en JSON
func printJSON(response *tsdio.ExecuteResponse, stdout io.Writer) error {
	data, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return fmt.Errorf("encodage JSON: %w", err)
	}

	fmt.Fprintln(stdout, string(data))
	return nil
}

// printText affiche les résultats en texte
func printText(config *Config, response *tsdio.ExecuteResponse, stdout, stderr io.Writer) error {
	if !response.Success {
		fmt.Fprintf(stderr, "\n❌ ERREUR D'EXÉCUTION\n")
		fmt.Fprintf(stderr, "===================\n")
		fmt.Fprintf(stderr, "Type: %s\n", response.ErrorType)
		fmt.Fprintf(stderr, "Message: %s\n", response.Error)
		fmt.Fprintf(stderr, "Temps: %dms\n", response.ExecutionTimeMs)
		return nil
	}

	fmt.Fprintf(stdout, "\n✅ EXÉCUTION RÉUSSIE\n")
	fmt.Fprintf(stdout, "===================\n")
	fmt.Fprintf(stdout, "Temps d'exécution: %dms\n", response.ExecutionTimeMs)
	fmt.Fprintf(stdout, "Faits injectés: %d\n", response.Results.FactsCount)
	fmt.Fprintf(stdout, "Activations: %d\n\n", response.Results.ActivationsCount)

	if response.Results.ActivationsCount == 0 {
		fmt.Fprintf(stdout, "ℹ️  Aucune action déclenchée\n")
		return nil
	}

	fmt.Fprintf(stdout, "🎯 ACTIONS DÉCLENCHÉES\n")
	fmt.Fprintf(stdout, "======================\n\n")

	for i, activation := range response.Results.Activations {
		fmt.Fprintf(stdout, "%d. Action: %s\n", i+1, activation.ActionName)

		// Afficher les arguments
		if len(activation.Arguments) > 0 {
			fmt.Fprintf(stdout, "   Arguments:\n")
			for _, arg := range activation.Arguments {
				fmt.Fprintf(stdout, "     [%d] %v (%s)\n", arg.Position, arg.Value, arg.Type)
			}
		}

		// Afficher les faits déclencheurs
		if config.Verbose && len(activation.TriggeringFacts) > 0 {
			fmt.Fprintf(stdout, "   Faits déclencheurs:\n")
			for j, fact := range activation.TriggeringFacts {
				fmt.Fprintf(stdout, "     [%d] %s (id: %s)\n", j, fact.Type, fact.ID)
				if len(fact.Fields) > 0 {
					for key, value := range fact.Fields {
						fmt.Fprintf(stdout, "         %s: %v\n", key, value)
					}
				}
			}
		}

		fmt.Fprintln(stdout)
	}

	return nil
}

// printHelp affiche l'aide
func printHelp(w io.Writer) {
	fmt.Fprintln(w, "TSD Client - Client pour le serveur TSD")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "USAGE:")
	fmt.Fprintln(w, "  tsd client <file.tsd> [options]")
	fmt.Fprintln(w, "  tsd client -file <file.tsd> [options]")
	fmt.Fprintln(w, "  tsd client -text \"<tsd code>\" [options]")
	fmt.Fprintln(w, "  tsd client -stdin [options]")
	fmt.Fprintln(w, "  echo \"<tsd code>\" | tsd client -stdin")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "OPTIONS:")
	fmt.Fprintln(w, "  -server <url>       URL du serveur TSD (défaut: https://localhost:8080)")
	fmt.Fprintln(w, "  -file <file>        Fichier TSD (.tsd)")
	fmt.Fprintln(w, "  -text <text>        Code TSD directement")
	fmt.Fprintln(w, "  -stdin              Lire depuis l'entrée standard")
	fmt.Fprintln(w, "  -v                  Mode verbeux (affiche plus de détails)")
	fmt.Fprintln(w, "  -format <format>    Format de sortie: text ou json (défaut: text)")
	fmt.Fprintln(w, "  -timeout <duration> Timeout des requêtes (défaut: 30s)")
	fmt.Fprintln(w, "  -health             Vérifier la santé du serveur")
	fmt.Fprintln(w, "  -token <token>      Token d'authentification (clé API ou JWT)")
	fmt.Fprintln(w, "  -auth-type <type>   Type d'authentification: key ou jwt (optionnel)")
	fmt.Fprintln(w, "  -tls-ca <file>      Certificat CA pour vérifier le serveur (défaut: ./certs/ca.crt)")
	fmt.Fprintln(w, "  -insecure           Désactiver la vérification TLS (⚠️  développement uniquement)")
	fmt.Fprintln(w, "  -help               Afficher cette aide")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "AUTHENTIFICATION:")
	fmt.Fprintln(w, "  Le token peut être fourni via -token ou la variable d'environnement TSD_AUTH_TOKEN")
	fmt.Fprintln(w, "  export TSD_AUTH_TOKEN=\"votre-token-ici\"")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "TLS/HTTPS:")
	fmt.Fprintln(w, "  Par défaut, le client utilise HTTPS et vérifie le certificat du serveur.")
	fmt.Fprintln(w, "  Pour les certificats auto-signés (développement):")
	fmt.Fprintln(w, "    - Option 1: Utiliser le CA généré: -tls-ca ./certs/ca.crt")
	fmt.Fprintln(w, "    - Option 2: Désactiver la vérification: -insecure (⚠️  non sécurisé)")
	fmt.Fprintln(w, "  Variables d'environnement: TSD_TLS_CA, TSD_CLIENT_INSECURE")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "EXEMPLES:")
	fmt.Fprintln(w, "  # Vérifier la santé du serveur (HTTPS par défaut)")
	fmt.Fprintln(w, "  tsd client -health")
	fmt.Fprintln(w, "  tsd client -health -server https://tsd.example.com:8080")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Exécuter un fichier TSD")
	fmt.Fprintln(w, "  tsd client program.tsd")
	fmt.Fprintln(w, "  tsd client -file program.tsd -v")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Exécuter du code TSD directement")
	fmt.Fprintln(w, "  tsd client -text 'type Person : <id: string, name: string>'")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Lire depuis stdin")
	fmt.Fprintln(w, "  echo 'type Person : <id: string>' | tsd client -stdin")
	fmt.Fprintln(w, "  cat program.tsd | tsd client -stdin -v")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Utiliser un serveur distant (HTTPS)")
	fmt.Fprintln(w, "  tsd client -server https://tsd.example.com:8080 program.tsd")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Utiliser un serveur avec certificat auto-signé (développement)")
	fmt.Fprintln(w, "  tsd client program.tsd -tls-ca ./certs/ca.crt")
	fmt.Fprintln(w, "  tsd client program.tsd -insecure")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Utiliser un serveur en HTTP non sécurisé")
	fmt.Fprintln(w, "  tsd client -server http://localhost:8080 program.tsd")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Format JSON pour intégration")
	fmt.Fprintln(w, "  tsd client program.tsd -format json")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Avec authentification par clé API")
	fmt.Fprintln(w, "  tsd client program.tsd -token \"votre-cle-api\"")
	fmt.Fprintln(w, "  export TSD_AUTH_TOKEN=\"votre-cle-api\"")
	fmt.Fprintln(w, "  tsd client program.tsd")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "  # Avec authentification JWT")
	fmt.Fprintln(w, "  tsd client program.tsd -token \"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...\"")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "NOTES:")
	fmt.Fprintln(w, "  - Le serveur utilise HTTPS par défaut (port 8080)")
	fmt.Fprintln(w, "  - Les certificats auto-signés nécessitent -tls-ca ou -insecure")
	fmt.Fprintln(w, "  - En production, utilisez toujours des certificats valides (Let's Encrypt, etc.)")
}
