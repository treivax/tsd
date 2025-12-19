// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package servercmd

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/treivax/tsd/auth"
	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/internal/tlsconfig"
	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/rete/actions"
	"github.com/treivax/tsd/tsdio"
	"github.com/treivax/tsd/xuples"
)

const (
	// DefaultPort est le port par défaut du serveur
	DefaultPort = 8080

	// DefaultHost est l'hôte par défaut du serveur
	DefaultHost = "0.0.0.0"

	// Version est la version du serveur
	Version = "1.0.0"

	// MaxRequestSize est la taille maximale d'une requête (10MB)
	MaxRequestSize = 10 * 1024 * 1024

	// DefaultCertDir est le répertoire par défaut des certificats
	DefaultCertDir = "./certs"

	// DefaultCertFile est le fichier de certificat par défaut
	DefaultCertFile = "server.crt"

	// DefaultKeyFile est le fichier de clé privée par défaut
	DefaultKeyFile = "server.key"

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

	// DefaultReadTimeout limite le temps pour lire la requête complète.
	// Protège contre slow client attacks.
	DefaultReadTimeout = 15 * time.Second

	// DefaultWriteTimeout limite le temps pour écrire la réponse.
	// Empêche blocage sur clients lents.
	DefaultWriteTimeout = 15 * time.Second

	// DefaultIdleTimeout limite le temps d'inactivité sur connexions keep-alive.
	// Libère ressources des connexions zombies.
	DefaultIdleTimeout = 60 * time.Second

	// DefaultReadHeaderTimeout limite le temps pour lire les headers HTTP.
	// Protection spécifique contre Slowloris.
	DefaultReadHeaderTimeout = 5 * time.Second

	// DefaultMaxHeaderBytes limite la taille maximale des headers HTTP (1 MB).
	// Protège contre headers excessifs et attaques par volume.
	DefaultMaxHeaderBytes = 1 << 20

	// DefaultShutdownTimeout est le timeout pour le graceful shutdown (30 secondes)
	DefaultShutdownTimeout = 30 * time.Second

	// Headers de sécurité HTTP recommandés pour API TSD

	// HeaderStrictTransportSecurity force HTTPS pour 1 an avec subdomains
	HeaderStrictTransportSecurity = "Strict-Transport-Security"
	// ValueHSTS est la valeur HSTS avec max-age d'un an et includeSubDomains
	ValueHSTS = "max-age=31536000; includeSubDomains"

	// HeaderXContentTypeOptions empêche MIME sniffing
	HeaderXContentTypeOptions = "X-Content-Type-Options"
	// ValueNoSniff désactive le MIME sniffing
	ValueNoSniff = "nosniff"

	// HeaderXFrameOptions empêche clickjacking
	HeaderXFrameOptions = "X-Frame-Options"
	// ValueDeny bloque totalement l'affichage en iframe
	ValueDeny = "DENY"

	// HeaderContentSecurityPolicy définit la politique de sécurité du contenu
	HeaderContentSecurityPolicy = "Content-Security-Policy"
	// ValueCSP est une politique stricte pour API (pas de contenu HTML/JS)
	ValueCSP = "default-src 'none'; frame-ancestors 'none'"

	// HeaderXXSSProtection active la protection XSS (legacy browsers)
	HeaderXXSSProtection = "X-XSS-Protection"
	// ValueXSSBlock active le blocage XSS
	ValueXSSBlock = "1; mode=block"

	// HeaderReferrerPolicy contrôle les informations de referrer
	HeaderReferrerPolicy = "Referrer-Policy"
	// ValueNoReferrer n'envoie aucune information de referrer
	ValueNoReferrer = "no-referrer"

	// HeaderServer masque la version du serveur
	HeaderServer = "Server"
	// ValueServerName est le nom générique du serveur
	ValueServerName = "TSD"
)

var (
	// startTime est l'heure de démarrage du serveur
	startTime = time.Now()
)

// Config contient la configuration du serveur
type Config struct {
	Host          string
	Port          int
	Verbose       bool
	AuthType      string
	AuthKeys      []string
	JWTSecret     string
	JWTExpiration time.Duration
	JWTIssuer     string
	TLSCertFile   string
	TLSKeyFile    string
	Insecure      bool
}

// Server représente le serveur HTTP TSD
type Server struct {
	config      *Config
	logger      *log.Logger
	mux         *http.ServeMux
	authManager *auth.Manager
	httpServer  *http.Server
}

// Run démarre le serveur TSD avec les arguments donnés et retourne un code de sortie
// ServerInfo contient les informations de configuration du serveur pour affichage
type ServerInfo struct {
	Addr        string
	Protocol    string
	Version     string
	TLSEnabled  bool
	TLSCertFile string
	TLSKeyFile  string
	AuthEnabled bool
	AuthType    string
	Endpoints   []string
}

// prepareServerInfo prépare les informations du serveur (logique testable)
func prepareServerInfo(config *Config, server *Server) *ServerInfo {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	// Déterminer le protocole
	protocol := "https"
	if config.Insecure {
		protocol = "http"
	}

	info := &ServerInfo{
		Addr:        addr,
		Protocol:    protocol,
		Version:     Version,
		TLSEnabled:  !config.Insecure,
		TLSCertFile: config.TLSCertFile,
		TLSKeyFile:  config.TLSKeyFile,
		AuthEnabled: server.authManager.IsEnabled(),
	}

	if info.AuthEnabled {
		info.AuthType = server.authManager.GetAuthType()
	}

	info.Endpoints = []string{
		fmt.Sprintf("POST %s://%s/api/v1/execute - Exécuter un programme TSD", protocol, addr),
		fmt.Sprintf("GET  %s://%s/health - Health check", protocol, addr),
		fmt.Sprintf("GET  %s://%s/api/v1/version - Version info", protocol, addr),
	}

	return info
}

// logServerInfo affiche les informations du serveur (logique testable)
func logServerInfo(logger *log.Logger, info *ServerInfo) {
	logger.Printf("🚀 Démarrage du serveur TSD sur %s://%s", info.Protocol, info.Addr)
	logger.Printf("📊 Version: %s", info.Version)

	// Afficher le statut TLS
	if info.TLSEnabled {
		logger.Printf("🔒 TLS: activé")
		logger.Printf("   Certificat: %s", info.TLSCertFile)
		logger.Printf("   Clé: %s", info.TLSKeyFile)
	} else {
		logger.Printf("⚠️  TLS: désactivé (mode HTTP non sécurisé)")
		logger.Printf("⚠️  AVERTISSEMENT: Ne pas utiliser en production!")
	}

	// Afficher le statut d'authentification
	if info.AuthEnabled {
		logger.Printf("🔒 Authentification: activée (%s)", info.AuthType)
	} else {
		logger.Printf("⚠️  Authentification: désactivée (mode développement)")
	}

	logger.Printf("🔗 Endpoints disponibles:")
	for _, endpoint := range info.Endpoints {
		logger.Printf("   %s", endpoint)
	}
}

// createTLSConfig crée la configuration TLS (logique testable)
func createTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	serverConfig := tlsconfig.DefaultServerConfig(certFile, keyFile)
	return tlsconfig.NewServerTLSConfig(serverConfig)
}

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	config := parseFlags(args)

	logger := log.New(stdout, "[TSD-SERVER] ", log.LstdFlags)

	server, initErr := NewServer(config, logger)
	if initErr != nil {
		fmt.Fprintf(stderr, "❌ Erreur initialisation serveur: %v\n", initErr)
		return 1
	}

	// Préparer et afficher les informations du serveur
	info := prepareServerInfo(config, server)
	logServerInfo(logger, info)

	// Créer le serveur HTTP avec timeouts de sécurité et l'attacher à la struct Server
	server.httpServer = &http.Server{
		Addr:    info.Addr,
		Handler: server.mux,

		// Timeouts de sécurité
		ReadTimeout:       DefaultReadTimeout,
		ReadHeaderTimeout: DefaultReadHeaderTimeout,
		WriteTimeout:      DefaultWriteTimeout,
		IdleTimeout:       DefaultIdleTimeout,

		// Limites supplémentaires
		MaxHeaderBytes: DefaultMaxHeaderBytes,
	}

	// Si TLS activé, configurer TLS
	if !config.Insecure {
		tlsConf, tlsErr := createTLSConfig(config.TLSCertFile, config.TLSKeyFile)
		if tlsErr != nil {
			fmt.Fprintf(stderr, "❌ Erreur configuration TLS: %v\n", tlsErr)
			return 1
		}
		server.httpServer.TLSConfig = tlsConf
	}

	// Canal pour capturer les erreurs du serveur
	serverErrors := make(chan error, 1)

	// Démarrer le serveur dans une goroutine
	go func() {
		var err error
		if config.Insecure {
			logger.Printf("⚠️  Démarrage en mode HTTP non sécurisé")
			err = server.httpServer.ListenAndServe()
		} else {
			logger.Printf("🔒 Démarrage en mode HTTPS sécurisé")
			err = server.httpServer.ListenAndServeTLS(config.TLSCertFile, config.TLSKeyFile)
		}
		if err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Canal pour capturer les signaux OS
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Attendre signal ou erreur
	select {
	case err := <-serverErrors:
		fmt.Fprintf(stderr, "❌ Erreur démarrage serveur: %v\n", err)
		return 1
	case sig := <-sigChan:
		logger.Printf("📡 Signal %v reçu, arrêt gracieux du serveur...", sig)
	}

	// Graceful shutdown avec timeout
	ctx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Fprintf(stderr, "❌ Erreur lors du shutdown: %v\n", err)
		return 1
	}

	return 0
}

// parseFlags parse les arguments de ligne de commande
func parseFlags(args []string) *Config {
	config := &Config{}

	fs := flag.NewFlagSet("tsd-server", flag.ContinueOnError)
	fs.StringVar(&config.Host, "host", DefaultHost, "Hôte du serveur")
	fs.IntVar(&config.Port, "port", DefaultPort, "Port du serveur")
	fs.BoolVar(&config.Verbose, "v", false, "Mode verbeux")

	// TLS
	defaultCertPath := filepath.Join(DefaultCertDir, DefaultCertFile)
	defaultKeyPath := filepath.Join(DefaultCertDir, DefaultKeyFile)
	fs.StringVar(&config.TLSCertFile, "tls-cert", defaultCertPath, "Chemin vers le certificat TLS")
	fs.StringVar(&config.TLSKeyFile, "tls-key", defaultKeyPath, "Chemin vers la clé privée TLS")
	fs.BoolVar(&config.Insecure, "insecure", false, "Désactiver TLS (mode HTTP non sécurisé)")

	// Authentification
	fs.StringVar(&config.AuthType, "auth", "none", "Type d'authentification: none, key, jwt")
	authKeysStr := fs.String("auth-keys", "", "Clés API (séparées par des virgules)")
	fs.StringVar(&config.JWTSecret, "jwt-secret", "", "Secret pour JWT")
	fs.DurationVar(&config.JWTExpiration, "jwt-expiration", 24*time.Hour, "Durée de validité JWT")
	fs.StringVar(&config.JWTIssuer, "jwt-issuer", "tsd-server", "Émetteur JWT")

	fs.Parse(args)

	// Variables d'environnement pour TLS
	if envCert := os.Getenv("TSD_TLS_CERT"); envCert != "" {
		config.TLSCertFile = envCert
	}
	if envKey := os.Getenv("TSD_TLS_KEY"); envKey != "" {
		config.TLSKeyFile = envKey
	}
	if os.Getenv("TSD_INSECURE") == "true" {
		config.Insecure = true
	}

	// Vérifier que les certificats existent si TLS est activé
	if !config.Insecure {
		if _, err := os.Stat(config.TLSCertFile); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "❌ Certificat TLS non trouvé: %s\n", config.TLSCertFile)
			fmt.Fprintf(os.Stderr, "\n💡 Solutions:\n")
			fmt.Fprintf(os.Stderr, "   1. Générer des certificats: tsd auth generate-cert\n")
			fmt.Fprintf(os.Stderr, "   2. Spécifier un certificat: --tls-cert /path/to/cert.crt\n")
			fmt.Fprintf(os.Stderr, "   3. Démarrer en mode non sécurisé: --insecure (déconseillé en production)\n")
			os.Exit(1)
		}
		if _, err := os.Stat(config.TLSKeyFile); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "❌ Clé privée TLS non trouvée: %s\n", config.TLSKeyFile)
			fmt.Fprintf(os.Stderr, "\n💡 Solutions:\n")
			fmt.Fprintf(os.Stderr, "   1. Générer des certificats: tsd auth generate-cert\n")
			fmt.Fprintf(os.Stderr, "   2. Spécifier une clé: --tls-key /path/to/key.key\n")
			fmt.Fprintf(os.Stderr, "   3. Démarrer en mode non sécurisé: --insecure (déconseillé en production)\n")
			os.Exit(1)
		}
	}

	// Parser les clés API depuis la variable d'environnement ou le flag
	if *authKeysStr == "" {
		*authKeysStr = os.Getenv("TSD_AUTH_KEYS")
	}
	if *authKeysStr != "" {
		config.AuthKeys = strings.Split(*authKeysStr, ",")
		for i, key := range config.AuthKeys {
			config.AuthKeys[i] = strings.TrimSpace(key)
		}
	}

	// Récupérer le secret JWT depuis la variable d'environnement si non fourni
	if config.JWTSecret == "" {
		config.JWTSecret = os.Getenv("TSD_JWT_SECRET")
	}

	return config
}

// NewServer crée un nouveau serveur TSD
func NewServer(config *Config, logger *log.Logger) (*Server, error) {
	// Créer le gestionnaire d'authentification
	authConfig := &auth.Config{
		Type:          config.AuthType,
		AuthKeys:      config.AuthKeys,
		JWTSecret:     config.JWTSecret,
		JWTExpiration: config.JWTExpiration,
		JWTIssuer:     config.JWTIssuer,
	}

	authManager, err := auth.NewManager(authConfig)
	if err != nil {
		return nil, fmt.Errorf("erreur initialisation authentification: %w", err)
	}

	s := &Server{
		config:      config,
		logger:      logger,
		mux:         http.NewServeMux(),
		authManager: authManager,
	}

	// Enregistrer les routes
	s.registerRoutes()

	return s, nil
}

// registerRoutes enregistre les routes HTTP
func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/api/v1/execute", s.withSecurityHeaders(s.validateContentType(s.handleExecute)))
	s.mux.HandleFunc("/health", s.withSecurityHeaders(s.handleHealth))
	s.mux.HandleFunc("/api/v1/version", s.withSecurityHeaders(s.handleVersion))
}

// Shutdown effectue un arrêt gracieux du serveur.
// Les nouvelles connexions sont refusées et les requêtes en cours sont
// terminées dans la limite du timeout spécifié via le contexte.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	s.logger.Printf("🛑 Arrêt gracieux du serveur démarré...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Printf("❌ Erreur lors du shutdown: %v", err)
		return fmt.Errorf("erreur shutdown serveur: %w", err)
	}

	s.logger.Printf("✅ Serveur arrêté proprement")
	return nil
}

// withSecurityHeaders ajoute les headers de sécurité HTTP à toutes les réponses.
// Ces headers protègent contre XSS, clickjacking, MIME sniffing et downgrade attacks.
func (s *Server) withSecurityHeaders(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// HSTS - Force HTTPS (1 an)
		w.Header().Set(HeaderStrictTransportSecurity, ValueHSTS)

		// Empêche MIME sniffing
		w.Header().Set(HeaderXContentTypeOptions, ValueNoSniff)

		// Empêche clickjacking
		w.Header().Set(HeaderXFrameOptions, ValueDeny)

		// CSP stricte (API uniquement, pas de contenu HTML/JS)
		w.Header().Set(HeaderContentSecurityPolicy, ValueCSP)

		// XSS Protection (legacy browsers)
		w.Header().Set(HeaderXXSSProtection, ValueXSSBlock)

		// Pas de referrer
		w.Header().Set(HeaderReferrerPolicy, ValueNoReferrer)

		// Masque version serveur
		w.Header().Set(HeaderServer, ValueServerName)

		handler(w, r)
	}
}

// validateContentType vérifie que le Content-Type est application/json pour les requêtes POST.
// Retourne 415 Unsupported Media Type si invalide.
func (s *Server) validateContentType(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Exemption pour GET (pas de body)
		if r.Method == http.MethodGet {
			handler(w, r)
			return
		}

		// Extraire Content-Type (ignorer charset)
		contentType := r.Header.Get("Content-Type")
		if contentType == "" {
			s.sendErrorResponse(w, StatusUnsupportedMedia, "Content-Type header requis", time.Now())
			return
		}

		// Parser pour ignorer charset
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			s.sendErrorResponse(w, StatusUnsupportedMedia, "Content-Type invalide", time.Now())
			return
		}

		// Valider que c'est application/json
		if mediaType != ContentTypeJSON {
			msg := fmt.Sprintf("Content-Type '%s' non supporté. Attendu: %s", mediaType, ContentTypeJSON)
			s.sendErrorResponse(w, StatusUnsupportedMedia, msg, time.Now())
			return
		}

		handler(w, r)
	}
}

// sendErrorResponse envoie une réponse d'erreur JSON standardisée.
func (s *Server) sendErrorResponse(w http.ResponseWriter, statusCode int, message string, startTime time.Time) {
	executionTimeMs := time.Since(startTime).Milliseconds()
	response := tsdio.NewErrorResponse(tsdio.ErrorTypeServerError, message, executionTimeMs)
	s.writeJSON(w, response, statusCode)
}

// handleExecute gère les requêtes d'exécution de programmes TSD
func (s *Server) handleExecute(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Vérifier la méthode HTTP
	if r.Method != http.MethodPost {
		s.sendErrorResponse(w, http.StatusMethodNotAllowed, "Méthode non autorisée", startTime)
		return
	}

	// Authentification
	if err := s.authenticate(r); err != nil {
		s.sendErrorResponse(w, StatusUnauthorized, "Authentification échouée: "+err.Error(), startTime)
		return
	}

	// Limiter la taille de la requête
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestSize)

	// Décoder la requête JSON avec validation stricte
	var req tsdio.ExecuteRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		s.sendErrorResponse(w, StatusBadRequest, fmt.Sprintf("JSON invalide: %v", err), startTime)
		return
	}

	// Valider la requête
	if req.Source == "" {
		s.sendErrorResponse(w, StatusBadRequest, "Le champ 'source' est requis", startTime)
		return
	}

	if req.SourceName == "" {
		req.SourceName = "<request>"
	}

	if s.config.Verbose || req.Verbose {
		s.logger.Printf("📥 Requête d'exécution reçue: source=%s, length=%d", req.SourceName, len(req.Source))
	}

	// Exécuter le programme TSD
	response := s.executeTSDProgram(&req, startTime)

	// Écrire la réponse
	s.writeJSON(w, response, StatusOK)

	if s.config.Verbose || req.Verbose {
		if response.Success {
			s.logger.Printf("✅ Exécution réussie: %d activations en %dms",
				response.Results.ActivationsCount, response.ExecutionTimeMs)
		} else {
			s.logger.Printf("❌ Exécution échouée: %s (%s) en %dms",
				response.ErrorType, response.Error, response.ExecutionTimeMs)
		}
	}
}

// executeTSDProgram exécute un programme TSD et retourne la réponse
func (s *Server) executeTSDProgram(req *tsdio.ExecuteRequest, startTime time.Time) *tsdio.ExecuteResponse {
	// Parser le programme TSD
	resultRaw, err := constraint.ParseConstraint(req.SourceName, []byte(req.Source))
	if err != nil {
		executionTimeMs := time.Since(startTime).Milliseconds()
		return tsdio.NewErrorResponse(tsdio.ErrorTypeParsingError, fmt.Sprintf("Erreur de parsing: %v", err), executionTimeMs)
	}

	// Valider le programme
	if err := constraint.ValidateConstraintProgram(resultRaw); err != nil {
		executionTimeMs := time.Since(startTime).Milliseconds()
		return tsdio.NewErrorResponse(tsdio.ErrorTypeValidationError, fmt.Sprintf("Erreur de validation: %v", err), executionTimeMs)
	}

	// Convertir le résultat en Program pour accéder aux XupleSpaces
	result, err := constraint.ConvertResultToProgram(resultRaw)
	if err != nil {
		executionTimeMs := time.Since(startTime).Milliseconds()
		return tsdio.NewErrorResponse(tsdio.ErrorTypeParsingError, fmt.Sprintf("Erreur conversion: %v", err), executionTimeMs)
	}

	// Créer le XupleManager et instancier les xuple-spaces déclarés
	xupleManager := xuples.NewXupleManager()
	if err := instantiateXupleSpaces(xupleManager, result.XupleSpaces); err != nil {
		executionTimeMs := time.Since(startTime).Milliseconds()
		return tsdio.NewErrorResponse(tsdio.ErrorTypeValidationError, fmt.Sprintf("Erreur création xuple-spaces: %v", err), executionTimeMs)
	}

	// Créer le pipeline RETE
	pipeline := rete.NewConstraintPipeline()
	storage := rete.NewMemoryStorage()
	
	// Créer le réseau RETE et configurer le XupleHandler AVANT l'ingestion
	network := rete.NewReteNetwork(storage)
	network.SetXupleManager(xupleManager)
	network.SetXupleHandler(func(xuplespace string, fact *rete.Fact, triggeringFacts []*rete.Fact) error {
		return xupleManager.CreateXuple(xuplespace, fact, triggeringFacts)
	})
	
	// Créer un collecteur d'exécutions (observer pattern) et le configurer AVANT l'ingestion
	statsCollector := NewExecutionStatsCollector()
	network.SetActionObserver(statsCollector)

	// Créer un fichier temporaire pour le source
	tmpFile, err := os.CreateTemp("", "tsd-*.tsd")
	if err != nil {
		executionTimeMs := time.Since(startTime).Milliseconds()
		return tsdio.NewErrorResponse(tsdio.ErrorTypeServerError, fmt.Sprintf("Erreur création fichier temporaire: %v", err), executionTimeMs)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Écrire le source dans le fichier temporaire
	if _, err := tmpFile.Write([]byte(req.Source)); err != nil {
		executionTimeMs := time.Since(startTime).Milliseconds()
		return tsdio.NewErrorResponse(tsdio.ErrorTypeServerError, fmt.Sprintf("Erreur écriture fichier temporaire: %v", err), executionTimeMs)
	}
	tmpFile.Close()

	// Ingérer le fichier avec le réseau pré-configuré
	network, _, err = pipeline.IngestFile(tmpFile.Name(), network, storage)
	if err != nil {
		executionTimeMs := time.Since(startTime).Milliseconds()
		return tsdio.NewErrorResponse(tsdio.ErrorTypeExecutionError, fmt.Sprintf("Erreur ingestion: %v", err), executionTimeMs)
	}

	// Configurer le BuiltinActionExecutor avec le XupleManager
	builtinExecutor := actions.NewBuiltinActionExecutor(network, xupleManager, os.Stdout, s.logger)
	_ = builtinExecutor // Éviter warning unused - sera utilisé dans une version future pour d'autres actions

	// Les activations sont maintenant capturées automatiquement via l'observer
	// pendant l'exécution des règles (observer configuré avant l'ingestion)

	// Collecter les résultats
	facts := storage.GetAllFacts()
	activations := statsCollector.GetActivations()

	executionTimeMs := time.Since(startTime).Milliseconds()

	results := &tsdio.ExecutionResults{
		FactsCount:       len(facts),
		ActivationsCount: len(activations),
		Activations:      activations,
	}

	return tsdio.NewSuccessResponse(results, executionTimeMs)
}

// instantiateXupleSpaces crée les xuple-spaces déclarés dans le programme.
//
// Cette fonction convertit les déclarations AST en configurations xuples
// et les instancie dans le XupleManager.
//
// Paramètres:
//   - xupleManager: gestionnaire de xuple-spaces
//   - declarations: déclarations de xuple-spaces depuis le parser
//
// Retourne:
//   - error: erreur si la création échoue
func instantiateXupleSpaces(xupleManager xuples.XupleManager, declarations []constraint.XupleSpaceDeclaration) error {
	for _, decl := range declarations {
		// Construire la configuration depuis la déclaration AST
		config, err := buildXupleSpaceConfig(decl)
		if err != nil {
			return fmt.Errorf("failed to build config for xuple-space '%s': %w", decl.Name, err)
		}

		// Créer le xuple-space
		if err := xupleManager.CreateXupleSpace(decl.Name, config); err != nil {
			return fmt.Errorf("failed to create xuple-space '%s': %w", decl.Name, err)
		}

		log.Printf("✅ Created xuple-space '%s' with policies: selection=%s, consumption=%s, retention=%s",
			decl.Name,
			config.SelectionPolicy.Name(),
			config.ConsumptionPolicy.Name(),
			config.RetentionPolicy.Name())
	}

	return nil
}

// buildXupleSpaceConfig construit une configuration xuples depuis une déclaration AST.
//
// Convertit les politiques définies dans le langage TSD en politiques concrètes
// du module xuples.
//
// Paramètres:
//   - decl: déclaration xuple-space depuis le parser
//
// Retourne:
//   - xuples.XupleSpaceConfig: configuration construite
//   - error: erreur si la conversion échoue
func buildXupleSpaceConfig(decl constraint.XupleSpaceDeclaration) (xuples.XupleSpaceConfig, error) {
	// Construire la politique de sélection
	selectionPolicy, err := buildSelectionPolicy(decl.SelectionPolicy)
	if err != nil {
		return xuples.XupleSpaceConfig{}, err
	}

	// Construire la politique de consommation
	consumptionPolicy, err := buildConsumptionPolicy(decl.ConsumptionPolicy)
	if err != nil {
		return xuples.XupleSpaceConfig{}, err
	}

	// Construire la politique de rétention
	retentionPolicy, err := buildRetentionPolicy(decl.RetentionPolicy)
	if err != nil {
		return xuples.XupleSpaceConfig{}, err
	}

	return xuples.XupleSpaceConfig{
		Name:              decl.Name,
		SelectionPolicy:   selectionPolicy,
		ConsumptionPolicy: consumptionPolicy,
		RetentionPolicy:   retentionPolicy,
	}, nil
}

// buildSelectionPolicy construit une politique de sélection.
func buildSelectionPolicy(policyName string) (xuples.SelectionPolicy, error) {
	switch policyName {
	case "random":
		return xuples.NewRandomSelectionPolicy(), nil
	case "fifo":
		return xuples.NewFIFOSelectionPolicy(), nil
	case "lifo":
		return xuples.NewLIFOSelectionPolicy(), nil
	default:
		return nil, fmt.Errorf("unknown selection policy: %s", policyName)
	}
}

// buildConsumptionPolicy construit une politique de consommation.
func buildConsumptionPolicy(conf constraint.XupleConsumptionPolicyConf) (xuples.ConsumptionPolicy, error) {
	switch conf.Type {
	case "once":
		return xuples.NewOnceConsumptionPolicy(), nil
	case "per-agent":
		return xuples.NewPerAgentConsumptionPolicy(), nil
	case "limited":
		if conf.Limit <= 0 {
			return nil, fmt.Errorf("limited consumption policy requires limit > 0, got %d", conf.Limit)
		}
		return xuples.NewLimitedConsumptionPolicy(conf.Limit), nil
	default:
		return nil, fmt.Errorf("unknown consumption policy: %s", conf.Type)
	}
}

// buildRetentionPolicy construit une politique de rétention.
func buildRetentionPolicy(conf constraint.XupleRetentionPolicyConf) (xuples.RetentionPolicy, error) {
	switch conf.Type {
	case "unlimited":
		return xuples.NewUnlimitedRetentionPolicy(), nil
	case "duration":
		if conf.Duration <= 0 {
			return nil, fmt.Errorf("duration retention policy requires duration > 0, got %d", conf.Duration)
		}
		return xuples.NewDurationRetentionPolicy(time.Duration(conf.Duration) * time.Second), nil
	default:
		return nil, fmt.Errorf("unknown retention policy: %s", conf.Type)
	}
}

// collectActivations collecte toutes les activations du réseau
// DEPRECATED: Utiliser ExecutionStatsCollector avec observer pattern à la place.
// Cette méthode est conservée temporairement pour compatibilité avec les tests existants.
//
// TODO(xuples): Supprimer cette méthode une fois tous les tests migrés vers observer pattern.
func (s *Server) collectActivations(network *rete.ReteNetwork) []tsdio.Activation {
	if network == nil {
		return []tsdio.Activation{}
	}

	activations := []tsdio.Activation{}

	for _, terminal := range network.TerminalNodes {
		// Utiliser les statistiques d'exécution au lieu de Memory.Tokens
		count := terminal.GetExecutionCount()
		if count == 0 {
			continue
		}

		actionName := "unknown"
		if terminal.Action != nil {
			jobs := terminal.Action.GetJobs()
			if len(jobs) > 0 {
				actionName = jobs[0].Name
			}
		}

		// Pour compatibilité, créer une activation par exécution
		// Note: Les détails du token ne sont plus disponibles via cette méthode
		lastResult := terminal.GetLastExecutionResult()
		if lastResult != nil {
			activation := tsdio.Activation{
				ActionName:      actionName,
				Arguments:       formatArguments(lastResult.Arguments),
				TriggeringFacts: extractFacts(lastResult.Context.Token),
				BindingsCount:   len(lastResult.Context.Token.Facts),
			}
			activations = append(activations, activation)
		}
	}

	return activations
}

// extractArguments extrait les arguments d'une activation
func (s *Server) extractArguments(terminal *rete.TerminalNode, token *rete.Token) []tsdio.ArgumentValue {
	args := []tsdio.ArgumentValue{}

	if terminal.Action == nil || terminal.Action.Job == nil {
		return args
	}

	// Note: Les arguments ne peuvent pas être facilement évalués ici car
	// l'évaluateur d'arguments n'est pas exporté. On retourne les expressions
	// brutes converties en string.
	for i, argExpr := range terminal.Action.Job.Args {
		// Convertir l'expression en string
		value := fmt.Sprintf("%v", argExpr)

		argValue := tsdio.ArgumentValue{
			Position: i,
			Value:    value,
			Type:     "expression",
		}
		args = append(args, argValue)
	}

	return args
}

// getValueType retourne le type d'une valeur
func (s *Server) getValueType(value interface{}) string {
	if value == nil {
		return "nil"
	}

	switch value.(type) {
	case string:
		return "string"
	case int, int8, int16, int32, int64:
		return "int"
	case uint, uint8, uint16, uint32, uint64:
		return "uint"
	case float32, float64:
		return "float"
	case bool:
		return "bool"
	default:
		return "unknown"
	}
}

// extractFacts extrait les faits d'un token
func (s *Server) extractFacts(token *rete.Token) []tsdio.Fact {
	facts := []tsdio.Fact{}

	for _, fact := range token.Facts {
		if fact == nil {
			continue
		}

		f := tsdio.Fact{
			ID:     fact.ID,
			Type:   fact.Type,
			Fields: s.extractAttributes(fact),
		}
		facts = append(facts, f)
	}

	return facts
}

// extractAttributes extrait les attributs d'un fait
func (s *Server) extractAttributes(fact *rete.Fact) map[string]interface{} {
	attrs := make(map[string]interface{})

	if fact.Fields != nil {
		for key, value := range fact.Fields {
			attrs[key] = value
		}
	}

	return attrs
}

// authenticate vérifie l'authentification de la requête
func (s *Server) authenticate(r *http.Request) error {
	if !s.authManager.IsEnabled() {
		return nil
	}

	// Extraire le token du header Authorization
	authHeader := r.Header.Get("Authorization")
	token := auth.ExtractTokenFromHeader(authHeader)

	// Valider le token
	if err := s.authManager.ValidateToken(token); err != nil {
		return err
	}

	return nil
}

// handleHealth gère les requêtes de health check
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendErrorResponse(w, http.StatusMethodNotAllowed, "Méthode non autorisée", time.Now())
		return
	}

	uptime := time.Since(startTime).Seconds()

	response := tsdio.HealthResponse{
		Status:        "ok",
		Version:       Version,
		UptimeSeconds: int64(uptime),
		Timestamp:     time.Now(),
	}

	s.writeJSON(w, response, StatusOK)
}

// handleVersion gère les requêtes de version
func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendErrorResponse(w, http.StatusMethodNotAllowed, "Méthode non autorisée", time.Now())
		return
	}

	response := tsdio.VersionResponse{
		Version:   Version,
		GoVersion: runtime.Version(),
	}

	s.writeJSON(w, response, StatusOK)
}

// writeJSON écrit une réponse JSON
func (s *Server) writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Printf("❌ Erreur encodage JSON: %v", err)
	}
}
