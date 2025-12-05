// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/treivax/tsd/auth"
	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/tsdio"
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
}

// Server représente le serveur HTTP TSD
type Server struct {
	config      *Config
	logger      *log.Logger
	mux         *http.ServeMux
	authManager *auth.Manager
}

func main() {
	config := parseFlags()

	logger := log.New(os.Stdout, "[TSD-SERVER] ", log.LstdFlags)

	server, err := NewServer(config, logger)
	if err != nil {
		logger.Fatalf("❌ Erreur initialisation serveur: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	logger.Printf("🚀 Démarrage du serveur TSD sur %s", addr)
	logger.Printf("📊 Version: %s", Version)

	// Afficher le statut d'authentification
	if server.authManager.IsEnabled() {
		logger.Printf("🔒 Authentification: activée (%s)", server.authManager.GetAuthType())
	} else {
		logger.Printf("⚠️  Authentification: désactivée (mode développement)")
	}

	logger.Printf("🔗 Endpoints disponibles:")
	logger.Printf("   POST http://%s/api/v1/execute - Exécuter un programme TSD", addr)
	logger.Printf("   GET  http://%s/health - Health check", addr)
	logger.Printf("   GET  http://%s/api/v1/version - Version info", addr)

	if err := http.ListenAndServe(addr, server.mux); err != nil {
		logger.Fatalf("❌ Erreur démarrage serveur: %v", err)
	}
}

// parseFlags parse les arguments de ligne de commande
func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.Host, "host", DefaultHost, "Hôte du serveur")
	flag.IntVar(&config.Port, "port", DefaultPort, "Port du serveur")
	flag.BoolVar(&config.Verbose, "v", false, "Mode verbeux")

	// Authentification
	flag.StringVar(&config.AuthType, "auth", "none", "Type d'authentification: none, key, jwt")
	authKeysStr := flag.String("auth-keys", "", "Clés API (séparées par des virgules)")
	flag.StringVar(&config.JWTSecret, "jwt-secret", "", "Secret pour JWT")
	flag.DurationVar(&config.JWTExpiration, "jwt-expiration", 24*time.Hour, "Durée de validité JWT")
	flag.StringVar(&config.JWTIssuer, "jwt-issuer", "tsd-server", "Émetteur JWT")

	flag.Parse()

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
	s.mux.HandleFunc("/api/v1/execute", s.handleExecute)
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/api/v1/version", s.handleVersion)
}

// handleExecute gère les requêtes d'exécution de programmes TSD
func (s *Server) handleExecute(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Vérifier la méthode HTTP
	if r.Method != http.MethodPost {
		s.writeError(w, tsdio.ErrorTypeServerError, "Méthode non autorisée", http.StatusMethodNotAllowed, startTime)
		return
	}

	// Authentification
	if err := s.authenticate(r); err != nil {
		s.writeError(w, tsdio.ErrorTypeServerError, "Authentification échouée: "+err.Error(), http.StatusUnauthorized, startTime)
		return
	}

	// Limiter la taille de la requête
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestSize)

	// Décoder la requête JSON
	var req tsdio.ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, tsdio.ErrorTypeServerError, fmt.Sprintf("Erreur décodage JSON: %v", err), http.StatusBadRequest, startTime)
		return
	}

	// Valider la requête
	if req.Source == "" {
		s.writeError(w, tsdio.ErrorTypeServerError, "Le champ 'source' est requis", http.StatusBadRequest, startTime)
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
	s.writeJSON(w, response, http.StatusOK)

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
	result, err := constraint.ParseConstraint(req.SourceName, []byte(req.Source))
	if err != nil {
		executionTimeMs := time.Since(startTime).Milliseconds()
		return tsdio.NewErrorResponse(tsdio.ErrorTypeParsingError, fmt.Sprintf("Erreur de parsing: %v", err), executionTimeMs)
	}

	// Valider le programme
	if err := constraint.ValidateConstraintProgram(result); err != nil {
		executionTimeMs := time.Since(startTime).Milliseconds()
		return tsdio.NewErrorResponse(tsdio.ErrorTypeValidationError, fmt.Sprintf("Erreur de validation: %v", err), executionTimeMs)
	}

	// Créer le pipeline RETE
	pipeline := rete.NewConstraintPipeline()
	storage := rete.NewMemoryStorage()

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

	// Ingérer le fichier
	network, err := pipeline.IngestFile(tmpFile.Name(), nil, storage)
	if err != nil {
		executionTimeMs := time.Since(startTime).Milliseconds()
		return tsdio.NewErrorResponse(tsdio.ErrorTypeExecutionError, fmt.Sprintf("Erreur ingestion: %v", err), executionTimeMs)
	}

	// Collecter les résultats
	facts := storage.GetAllFacts()
	activations := s.collectActivations(network)

	executionTimeMs := time.Since(startTime).Milliseconds()

	results := &tsdio.ExecutionResults{
		FactsCount:       len(facts),
		ActivationsCount: len(activations),
		Activations:      activations,
	}

	return tsdio.NewSuccessResponse(results, executionTimeMs)
}

// collectActivations collecte toutes les activations du réseau
func (s *Server) collectActivations(network *rete.ReteNetwork) []tsdio.Activation {
	if network == nil {
		return []tsdio.Activation{}
	}

	activations := []tsdio.Activation{}

	for _, terminal := range network.TerminalNodes {
		if terminal.Memory == nil || terminal.Memory.Tokens == nil {
			continue
		}

		actionName := "unknown"
		if terminal.Action != nil && terminal.Action.Job != nil {
			actionName = terminal.Action.Job.Name
		}

		for _, token := range terminal.Memory.Tokens {
			activation := tsdio.Activation{
				ActionName:      actionName,
				Arguments:       s.extractArguments(terminal, token),
				TriggeringFacts: s.extractFacts(token),
				BindingsCount:   len(token.Facts),
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
			ID:         fact.ID,
			Type:       fact.Type,
			Attributes: s.extractAttributes(fact),
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
		s.writeError(w, tsdio.ErrorTypeServerError, "Méthode non autorisée", http.StatusMethodNotAllowed, time.Now())
		return
	}

	uptime := time.Since(startTime).Seconds()

	response := tsdio.HealthResponse{
		Status:        "ok",
		Version:       Version,
		UptimeSeconds: int64(uptime),
		Timestamp:     time.Now(),
	}

	s.writeJSON(w, response, http.StatusOK)
}

// handleVersion gère les requêtes de version
func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeError(w, tsdio.ErrorTypeServerError, "Méthode non autorisée", http.StatusMethodNotAllowed, time.Now())
		return
	}

	response := tsdio.VersionResponse{
		Version:   Version,
		GoVersion: runtime.Version(),
	}

	s.writeJSON(w, response, http.StatusOK)
}

// writeJSON écrit une réponse JSON
func (s *Server) writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Printf("❌ Erreur encodage JSON: %v", err)
	}
}

// writeError écrit une réponse d'erreur
func (s *Server) writeError(w http.ResponseWriter, errorType, message string, statusCode int, startTime time.Time) {
	executionTimeMs := time.Since(startTime).Milliseconds()
	response := tsdio.NewErrorResponse(errorType, message, executionTimeMs)
	s.writeJSON(w, response, statusCode)
}
