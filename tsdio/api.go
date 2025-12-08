// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package tsdio

import (
	"time"
)

// ExecuteRequest représente une requête d'exécution de programme TSD
type ExecuteRequest struct {
	// Source contient le code TSD à exécuter (types, facts, rules)
	Source string `json:"source"`

	// SourceName est le nom de la source (pour messages d'erreur)
	// Optionnel, par défaut "<request>"
	SourceName string `json:"source_name,omitempty"`

	// Verbose active le mode verbeux pour plus de détails
	Verbose bool `json:"verbose,omitempty"`
}

// ExecuteResponse représente la réponse d'exécution
type ExecuteResponse struct {
	// Success indique si l'exécution a réussi
	Success bool `json:"success"`

	// Error contient le message d'erreur si Success == false
	Error string `json:"error,omitempty"`

	// ErrorType précise le type d'erreur (parsing, validation, execution)
	ErrorType string `json:"error_type,omitempty"`

	// Results contient les résultats d'exécution si Success == true
	Results *ExecutionResults `json:"results,omitempty"`

	// ExecutionTime est la durée d'exécution en millisecondes
	ExecutionTimeMs int64 `json:"execution_time_ms"`
}

// ExecutionResults contient les détails des résultats d'exécution
type ExecutionResults struct {
	// FactsCount est le nombre de faits injectés dans le réseau
	FactsCount int `json:"facts_count"`

	// ActivationsCount est le nombre total d'activations (actions déclenchées)
	ActivationsCount int `json:"activations_count"`

	// Activations contient les détails de chaque activation
	Activations []Activation `json:"activations"`
}

// Activation représente une action déclenchée avec ses détails
type Activation struct {
	// ActionName est le nom de l'action déclenchée
	ActionName string `json:"action_name"`

	// Arguments contient les arguments évalués de l'action
	Arguments []ArgumentValue `json:"arguments"`

	// TriggeringFacts contient les faits qui ont déclenché cette action
	TriggeringFacts []Fact `json:"triggering_facts"`

	// BindingsCount est le nombre de bindings dans le token
	BindingsCount int `json:"bindings_count"`
}

// ArgumentValue représente un argument évalué d'une action
type ArgumentValue struct {
	// Position est la position de l'argument (0, 1, 2, ...)
	Position int `json:"position"`

	// Value est la valeur de l'argument (peut être string, number, bool, etc.)
	Value interface{} `json:"value"`

	// Type est le type de la valeur (pour affichage)
	Type string `json:"type"`
}

// Fact représente un fait dans le système
type Fact struct {
	// ID est l'identifiant unique du fait
	ID string `json:"id"`

	// Type est le type du fait (ex: "Person", "Order")
	Type string `json:"type"`

	// Fields contient les champs du fait (renommé de Attributes pour cohérence avec domain.Fact)
	Fields map[string]interface{} `json:"fields"`
}

// HealthResponse représente la réponse du endpoint health
type HealthResponse struct {
	// Status est le statut du serveur ("ok", "error")
	Status string `json:"status"`

	// Version est la version du serveur TSD
	Version string `json:"version"`

	// Uptime est le temps depuis le démarrage en secondes
	UptimeSeconds int64 `json:"uptime_seconds"`

	// Timestamp est le timestamp de la réponse
	Timestamp time.Time `json:"timestamp"`
}

// VersionResponse représente la réponse du endpoint version
type VersionResponse struct {
	// Version est la version du serveur TSD
	Version string `json:"version"`

	// BuildTime est la date/heure de build
	BuildTime string `json:"build_time,omitempty"`

	// GitCommit est le hash du commit git
	GitCommit string `json:"git_commit,omitempty"`

	// GoVersion est la version de Go utilisée
	GoVersion string `json:"go_version"`
}

// ErrorTypes constants pour les types d'erreurs
const (
	ErrorTypeParsingError    = "parsing_error"
	ErrorTypeValidationError = "validation_error"
	ErrorTypeExecutionError  = "execution_error"
	ErrorTypeServerError     = "server_error"
)

// NewExecuteRequest crée une nouvelle requête d'exécution
func NewExecuteRequest(source string) *ExecuteRequest {
	return &ExecuteRequest{
		Source:     source,
		SourceName: "<request>",
		Verbose:    false,
	}
}

// NewSuccessResponse crée une réponse de succès
func NewSuccessResponse(results *ExecutionResults, executionTimeMs int64) *ExecuteResponse {
	return &ExecuteResponse{
		Success:         true,
		Results:         results,
		ExecutionTimeMs: executionTimeMs,
	}
}

// NewErrorResponse crée une réponse d'erreur
func NewErrorResponse(errorType, errorMsg string, executionTimeMs int64) *ExecuteResponse {
	return &ExecuteResponse{
		Success:         false,
		Error:           errorMsg,
		ErrorType:       errorType,
		ExecutionTimeMs: executionTimeMs,
	}
}
