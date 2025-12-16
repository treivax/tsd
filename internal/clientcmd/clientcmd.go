// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package clientcmd

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/treivax/tsd/internal/tlsconfig"
	"github.com/treivax/tsd/tsdio"
)

const (
	// DefaultServerURL est l'URL par défaut du serveur
	DefaultServerURL = "https://localhost:8080"

	// DefaultTimeout est le timeout par défaut des requêtes
	DefaultTimeout = 30 * time.Second

	// DefaultDialTimeout est le timeout de connexion TCP
	DefaultDialTimeout = 10 * time.Second

	// DefaultTLSHandshakeTimeout est le timeout du handshake TLS
	DefaultTLSHandshakeTimeout = 10 * time.Second

	// DefaultResponseHeaderTimeout est le timeout de lecture des headers de réponse
	DefaultResponseHeaderTimeout = 10 * time.Second

	// DefaultExpectContinueTimeout est le timeout pour Expect: 100-continue
	DefaultExpectContinueTimeout = 1 * time.Second

	// DefaultCAFile est le fichier CA par défaut
	DefaultCAFile = "./certs/ca.crt"

	// ExampleServerHTTPURL est une URL HTTP d'exemple pour la documentation
	ExampleServerHTTPURL = "http://localhost:8080"

	// Content-Type standards
	ContentTypeJSON        = "application/json"
	ContentTypeJSONCharset = "application/json; charset=utf-8"

	// Status codes HTTP métier (aliases pour clarté du code)
	StatusOK                  = http.StatusOK                   // 200
	StatusBadRequest          = http.StatusBadRequest           // 400
	StatusUnauthorized        = http.StatusUnauthorized         // 401
	StatusUnsupportedMedia    = http.StatusUnsupportedMediaType // 415
	StatusInternalServerError = http.StatusInternalServerError  // 500
	StatusServiceUnavailable  = http.StatusServiceUnavailable   // 503
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
	config      *Config
	httpClient  *http.Client
	tlsConfig   *tls.Config
	retryConfig RetryConfig
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
	// Configurer TLS via le package tlsconfig
	clientTLSConfig := &tlsconfig.ClientConfig{
		CAFile:             config.TLSCAFile,
		InsecureSkipVerify: config.Insecure,
	}

	tlsConf, err := tlsconfig.NewClientTLSConfig(clientTLSConfig)
	if err != nil {
		// En cas d'erreur de configuration CA, utiliser config de base
		// mais logger un warning si en mode sécurisé
		if !config.Insecure {
			fmt.Fprintf(os.Stderr, "⚠️  Erreur configuration TLS: %v\n", err)
			fmt.Fprintf(os.Stderr, "⚠️  Utilisation configuration TLS par défaut\n")
		}
		tlsConf = &tls.Config{
			InsecureSkipVerify: config.Insecure,
		}
	}

	// Configurer le transport avec timeouts granulaires
	transport := &http.Transport{
		TLSClientConfig: tlsConf,
		DialContext: (&net.Dialer{
			Timeout: DefaultDialTimeout,
		}).DialContext,
		TLSHandshakeTimeout:   DefaultTLSHandshakeTimeout,
		ResponseHeaderTimeout: DefaultResponseHeaderTimeout,
		ExpectContinueTimeout: DefaultExpectContinueTimeout,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
	}

	return &Client{
		config:      config,
		tlsConfig:   tlsConf,
		retryConfig: DefaultRetryConfig(),
		httpClient: &http.Client{
			Timeout:   config.Timeout,
			Transport: transport,
		},
	}
}

// SetRetryConfig permet de configurer le comportement de retry.
func (c *Client) SetRetryConfig(config RetryConfig) {
	c.retryConfig = config
}

// logRetryAttempt affiche un log de tentative de retry si en mode verbeux.
func (c *Client) logRetryAttempt(attempt int) {
	if attempt > 0 && c.config.Verbose {
		fmt.Printf("🔄 Tentative %d/%d...\n", attempt+1, c.retryConfig.MaxAttempts)
	}
}

// logRetryBackoff affiche un log d'attente avant retry si en mode verbeux.
func (c *Client) logRetryBackoff(backoff time.Duration, err error, statusCode int) {
	if !c.config.Verbose {
		return
	}
	if err != nil {
		fmt.Printf("⏱️  Erreur réseau transitoire, retry dans %v...\n", backoff)
	} else {
		fmt.Printf("⏱️  Erreur HTTP %d (transitoire), retry dans %v...\n", statusCode, backoff)
	}
}

// isSuccessResponse vérifie si la réponse est un succès.
func isSuccessResponse(err error, resp *http.Response) bool {
	return err == nil && resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 500
}

// waitForBackoff attend le délai de backoff ou retourne si le contexte est annulé.
func waitForBackoff(ctx context.Context, backoff time.Duration) error {
	select {
	case <-time.After(backoff):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// executeRequestWithRetry exécute une requête HTTP avec retry automatique.
func (c *Client) executeRequestWithRetry(req *http.Request) (*http.Response, error) {
	var lastErr error
	var resp *http.Response

	for attempt := 0; attempt < c.retryConfig.MaxAttempts; attempt++ {
		c.logRetryAttempt(attempt)

		resp, lastErr = c.httpClient.Do(req.Clone(req.Context()))

		if isSuccessResponse(lastErr, resp) {
			return resp, nil
		}

		statusCode := 0
		if resp != nil {
			statusCode = resp.StatusCode
			resp.Body.Close()
		}

		if !isRetryableError(lastErr, statusCode, c.retryConfig) {
			if lastErr != nil {
				return nil, lastErr
			}
			return resp, nil
		}

		if attempt == c.retryConfig.MaxAttempts-1 {
			break
		}

		backoff := calculateBackoff(attempt, c.retryConfig)
		c.logRetryBackoff(backoff, lastErr, statusCode)

		if err := waitForBackoff(req.Context(), backoff); err != nil {
			return nil, err
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("échec après %d tentatives: %w", c.retryConfig.MaxAttempts, lastErr)
	}

	return resp, nil
}

// validateResponse vérifie que la réponse HTTP est valide.
func validateResponse(resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("réponse nil")
	}

	// Vérifier Content-Type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		return fmt.Errorf("réponse sans Content-Type (status=%d)", resp.StatusCode)
	}

	// Parser Content-Type
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return fmt.Errorf("Content-Type invalide '%s': %w", contentType, err)
	}

	// Valider que c'est JSON
	if mediaType != ContentTypeJSON {
		return fmt.Errorf("Content-Type inattendu: '%s' (attendu: %s)", mediaType, ContentTypeJSON)
	}

	return nil
}

// parseErrorResponse extrait le message d'erreur du body JSON.
func parseErrorResponse(resp *http.Response) string {
	var errResp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		// Impossible de parser, retourner message générique
		return fmt.Sprintf("Erreur HTTP %d", resp.StatusCode)
	}

	// Essayer d'abord le champ 'error', puis 'message'
	if errResp.Error != "" {
		return errResp.Error
	}
	if errResp.Message != "" {
		return errResp.Message
	}

	return fmt.Sprintf("Erreur HTTP %d", resp.StatusCode)
}

// logExecuteRequest affiche des logs de requête si en mode verbeux.
func (c *Client) logExecuteRequest(url string) {
	if !c.config.Verbose {
		return
	}
	fmt.Printf("📤 Envoi requête à %s...\n", url)
	if c.config.AuthToken != "" {
		fmt.Printf("🔒 Authentification: activée\n")
	}
	if c.config.Insecure {
		fmt.Printf("⚠️  TLS: vérification désactivée (mode insecure)\n")
	}
}

// createExecuteRequest crée une requête HTTP pour l'exécution.
func (c *Client) createExecuteRequest(source, sourceName string) (*http.Request, error) {
	req := tsdio.ExecuteRequest{
		Source:     source,
		SourceName: sourceName,
		Verbose:    c.config.Verbose,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encodage JSON: %w", err)
	}

	url := c.config.ServerURL + "/api/v1/execute"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("création requête: %w", err)
	}

	httpReq.Header.Set("Content-Type", ContentTypeJSON)
	if c.config.AuthToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.config.AuthToken)
	}

	return httpReq, nil
}

// handleExecuteResponse traite la réponse HTTP et retourne la réponse TSD.
func handleExecuteResponse(resp *http.Response) (*tsdio.ExecuteResponse, error) {
	switch resp.StatusCode {
	case StatusOK:
		var response tsdio.ExecuteResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("erreur parsing réponse: %w", err)
		}
		return &response, nil

	case StatusBadRequest:
		msg := parseErrorResponse(resp)
		return nil, fmt.Errorf("requête invalide: %s", msg)

	case StatusUnauthorized:
		msg := parseErrorResponse(resp)
		return nil, fmt.Errorf("non autorisé: %s (token invalide/expiré?)", msg)

	case StatusUnsupportedMedia:
		msg := parseErrorResponse(resp)
		return nil, fmt.Errorf("Content-Type non supporté: %s", msg)

	case StatusInternalServerError:
		msg := parseErrorResponse(resp)
		return nil, fmt.Errorf("erreur serveur: %s", msg)

	case StatusServiceUnavailable:
		msg := parseErrorResponse(resp)
		return nil, fmt.Errorf("serveur indisponible: %s", msg)

	default:
		msg := parseErrorResponse(resp)
		return nil, fmt.Errorf("erreur HTTP %d: %s", resp.StatusCode, msg)
	}
}

// Execute envoie une requête d'exécution au serveur
func (c *Client) Execute(source, sourceName string) (*tsdio.ExecuteResponse, error) {
	httpReq, err := c.createExecuteRequest(source, sourceName)
	if err != nil {
		return nil, err
	}

	c.logExecuteRequest(c.config.ServerURL + "/api/v1/execute")

	resp, err := c.executeRequestWithRetry(httpReq)
	if err != nil {
		return nil, fmt.Errorf("envoi requête: %w", err)
	}
	defer resp.Body.Close()

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	return handleExecuteResponse(resp)
}

// HealthCheck vérifie la santé du serveur
func (c *Client) HealthCheck() (*tsdio.HealthResponse, error) {
	url := c.config.ServerURL + "/health"
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("création requête health: %w", err)
	}

	resp, err := c.executeRequestWithRetry(httpReq)
	if err != nil {
		return nil, fmt.Errorf("requête health: %w", err)
	}
	defer resp.Body.Close()

	// Vérifier le status code
	if resp.StatusCode != StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erreur HTTP %d: %s", resp.StatusCode, string(body))
	}

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
	fmt.Fprintf(w, "  tsd client -server %s program.tsd\n", ExampleServerHTTPURL)
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
