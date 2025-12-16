// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package clientcmd

import (
	"errors"
	"math"
	"math/rand"
	"net"
	"net/url"
	"syscall"
	"time"
)

// Configuration retry par défaut
const (
	DefaultMaxAttempts = 3
	DefaultBaseDelay   = 1 * time.Second
	DefaultMaxDelay    = 10 * time.Second
	DefaultJitter      = 0.2
)

// DefaultRetryableStatusCodes définit les codes HTTP qui justifient un retry
var DefaultRetryableStatusCodes = []int{
	500, // Internal Server Error
	502, // Bad Gateway
	503, // Service Unavailable
	504, // Gateway Timeout
}

// RetryConfig définit la configuration de retry pour le client HTTP.
type RetryConfig struct {
	// MaxAttempts est le nombre total de tentatives (initiale + retries).
	// Valeur par défaut : 3
	MaxAttempts int

	// BaseDelay est le délai initial avant le premier retry.
	// Valeur par défaut : 1 seconde
	BaseDelay time.Duration

	// MaxDelay est le délai maximum entre deux tentatives.
	// Évite les délais exponentiels trop longs.
	// Valeur par défaut : 10 secondes
	MaxDelay time.Duration

	// Jitter est le pourcentage de variation aléatoire (0.0-1.0).
	// Exemple : 0.2 = ±20%
	// Valeur par défaut : 0.2 (20%)
	Jitter float64

	// RetryableStatusCodes liste les codes HTTP à retry.
	// Valeur par défaut : [500, 502, 503, 504]
	RetryableStatusCodes []int
}

// DefaultRetryConfig retourne la configuration par défaut.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:          DefaultMaxAttempts,
		BaseDelay:            DefaultBaseDelay,
		MaxDelay:             DefaultMaxDelay,
		Jitter:               DefaultJitter,
		RetryableStatusCodes: DefaultRetryableStatusCodes,
	}
}

// isRetryableError détermine si une erreur justifie un retry.
// Retourne true pour les erreurs réseau transitoires et certains codes HTTP.
func isRetryableError(err error, statusCode int, config RetryConfig) bool {
	if err == nil && statusCode == 0 {
		return false
	}

	// Vérifier si c'est une erreur réseau transitoire
	if err != nil && isNetworkError(err) {
		return true
	}

	// Vérifier si le code HTTP est dans la liste retryable
	for _, code := range config.RetryableStatusCodes {
		if statusCode == code {
			return true
		}
	}

	return false
}

// isNetworkError détecte les erreurs réseau transitoires.
func isNetworkError(err error) bool {
	if err == nil {
		return false
	}

	// Timeout réseau
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	// URL errors (connexion refusée, DNS, etc.)
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return true
	}

	// Erreurs syscall spécifiques
	if errors.Is(err, syscall.ECONNREFUSED) ||
		errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, syscall.ETIMEDOUT) ||
		errors.Is(err, syscall.ENETUNREACH) {
		return true
	}

	return false
}

// calculateBackoff calcule le délai avant le prochain retry avec backoff exponentiel.
// Formule : min(maxDelay, baseDelay * 2^attempt) ± jitter
func calculateBackoff(attempt int, config RetryConfig) time.Duration {
	if attempt < 0 {
		return 0
	}

	// Backoff exponentiel
	exponential := float64(config.BaseDelay) * math.Pow(2, float64(attempt))
	delay := time.Duration(exponential)

	// Plafonner au max
	if delay > config.MaxDelay {
		delay = config.MaxDelay
	}

	// Ajouter jitter aléatoire (évite thundering herd)
	if config.Jitter > 0 {
		jitterRange := float64(delay) * config.Jitter
		jitterDelta := (rand.Float64()*2 - 1) * jitterRange // -jitter à +jitter
		delay = delay + time.Duration(jitterDelta)

		// S'assurer que le délai reste positif
		if delay < 0 {
			delay = 0
		}
	}

	return delay
}
