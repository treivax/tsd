// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package clientcmd

import (
	"errors"
	"net"
	"net/url"
	"syscall"
	"testing"
	"time"
)

func TestIsRetryableError_NetworkErrors(t *testing.T) {
	t.Log("🧪 TEST ERREURS RÉSEAU RETRYABLE")
	t.Log("=================================")

	config := DefaultRetryConfig()

	tests := []struct {
		name       string
		err        error
		statusCode int
		want       bool
	}{
		{
			name:       "timeout réseau",
			err:        &net.DNSError{IsTimeout: true},
			statusCode: 0,
			want:       true,
		},
		{
			name:       "connexion refusée",
			err:        syscall.ECONNREFUSED,
			statusCode: 0,
			want:       true,
		},
		{
			name:       "connexion reset",
			err:        syscall.ECONNRESET,
			statusCode: 0,
			want:       true,
		},
		{
			name:       "timeout syscall",
			err:        syscall.ETIMEDOUT,
			statusCode: 0,
			want:       true,
		},
		{
			name:       "network unreachable",
			err:        syscall.ENETUNREACH,
			statusCode: 0,
			want:       true,
		},
		{
			name:       "url error",
			err:        &url.Error{Op: "Get", URL: "http://example.com", Err: errors.New("test")},
			statusCode: 0,
			want:       true,
		},
		{
			name:       "500 Internal Server Error",
			err:        nil,
			statusCode: 500,
			want:       true,
		},
		{
			name:       "502 Bad Gateway",
			err:        nil,
			statusCode: 502,
			want:       true,
		},
		{
			name:       "503 Service Unavailable",
			err:        nil,
			statusCode: 503,
			want:       true,
		},
		{
			name:       "504 Gateway Timeout",
			err:        nil,
			statusCode: 504,
			want:       true,
		},
		{
			name:       "400 Bad Request",
			err:        nil,
			statusCode: 400,
			want:       false,
		},
		{
			name:       "401 Unauthorized",
			err:        nil,
			statusCode: 401,
			want:       false,
		},
		{
			name:       "404 Not Found",
			err:        nil,
			statusCode: 404,
			want:       false,
		},
		{
			name:       "200 OK",
			err:        nil,
			statusCode: 200,
			want:       false,
		},
		{
			name:       "pas d'erreur ni status code",
			err:        nil,
			statusCode: 0,
			want:       false,
		},
		{
			name:       "erreur générique",
			err:        errors.New("generic error"),
			statusCode: 0,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRetryableError(tt.err, tt.statusCode, config)
			if got != tt.want {
				t.Errorf("❌ isRetryableError() = %v, attendu %v", got, tt.want)
			} else {
				t.Logf("✅ Correct: retryable=%v", got)
			}
		})
	}
}

func TestCalculateBackoff(t *testing.T) {
	t.Log("🧪 TEST CALCUL BACKOFF")
	t.Log("======================")

	config := RetryConfig{
		BaseDelay: 1 * time.Second,
		MaxDelay:  10 * time.Second,
		Jitter:    0, // Pas de jitter pour tests déterministes
	}

	tests := []struct {
		name    string
		attempt int
		want    time.Duration
	}{
		{"premier retry (2^0)", 0, 1 * time.Second},
		{"deuxième retry (2^1)", 1, 2 * time.Second},
		{"troisième retry (2^2)", 2, 4 * time.Second},
		{"quatrième retry (2^3)", 3, 8 * time.Second},
		{"plafonné au max", 4, 10 * time.Second}, // 16s mais max=10s
		{"plafonné au max", 10, 10 * time.Second},
		{"attempt négatif", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateBackoff(tt.attempt, config)
			if got != tt.want {
				t.Errorf("❌ calculateBackoff(%d) = %v, attendu %v", tt.attempt, got, tt.want)
			} else {
				t.Logf("✅ Backoff correct: %v", got)
			}
		})
	}
}

func TestCalculateBackoff_WithJitter(t *testing.T) {
	t.Log("🧪 TEST BACKOFF AVEC JITTER")
	t.Log("============================")

	config := RetryConfig{
		BaseDelay: 1 * time.Second,
		MaxDelay:  10 * time.Second,
		Jitter:    0.2, // ±20%
	}

	// Tester que le jitter reste dans les bornes
	baseDelay := 2 * time.Second                           // Pour attempt=1
	minExpected := time.Duration(float64(baseDelay) * 0.8) // -20%
	maxExpected := time.Duration(float64(baseDelay) * 1.2) // +20%

	// Exécuter plusieurs fois pour tester la distribution
	for i := 0; i < 100; i++ {
		got := calculateBackoff(1, config)

		if got < minExpected || got > maxExpected {
			t.Errorf("❌ Backoff hors bornes: %v (attendu entre %v et %v)", got, minExpected, maxExpected)
		}
	}

	t.Logf("✅ Jitter reste dans les bornes ±20%%")
}

func TestDefaultRetryConfig(t *testing.T) {
	t.Log("🧪 TEST CONFIGURATION PAR DÉFAUT")
	t.Log("=================================")

	config := DefaultRetryConfig()

	if config.MaxAttempts != DefaultMaxAttempts {
		t.Errorf("❌ MaxAttempts = %d, attendu %d", config.MaxAttempts, DefaultMaxAttempts)
	}

	if config.BaseDelay != DefaultBaseDelay {
		t.Errorf("❌ BaseDelay = %v, attendu %v", config.BaseDelay, DefaultBaseDelay)
	}

	if config.MaxDelay != DefaultMaxDelay {
		t.Errorf("❌ MaxDelay = %v, attendu %v", config.MaxDelay, DefaultMaxDelay)
	}

	if config.Jitter != DefaultJitter {
		t.Errorf("❌ Jitter = %f, attendu %f", config.Jitter, DefaultJitter)
	}

	if len(config.RetryableStatusCodes) != 4 {
		t.Errorf("❌ RetryableStatusCodes count = %d, attendu 4", len(config.RetryableStatusCodes))
	}

	expectedCodes := []int{500, 502, 503, 504}
	for i, code := range expectedCodes {
		if config.RetryableStatusCodes[i] != code {
			t.Errorf("❌ RetryableStatusCodes[%d] = %d, attendu %d", i, config.RetryableStatusCodes[i], code)
		}
	}

	t.Logf("✅ Configuration par défaut correcte")
}

func TestIsNetworkError(t *testing.T) {
	t.Log("🧪 TEST DÉTECTION ERREURS RÉSEAU")
	t.Log("=================================")

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "timeout error",
			err:  &net.DNSError{IsTimeout: true},
			want: true,
		},
		{
			name: "url error",
			err:  &url.Error{Op: "Get", URL: "http://test.com", Err: errors.New("test")},
			want: true,
		},
		{
			name: "ECONNREFUSED",
			err:  syscall.ECONNREFUSED,
			want: true,
		},
		{
			name: "ECONNRESET",
			err:  syscall.ECONNRESET,
			want: true,
		},
		{
			name: "ETIMEDOUT",
			err:  syscall.ETIMEDOUT,
			want: true,
		},
		{
			name: "ENETUNREACH",
			err:  syscall.ENETUNREACH,
			want: true,
		},
		{
			name: "generic error",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isNetworkError(tt.err)
			if got != tt.want {
				t.Errorf("❌ isNetworkError() = %v, attendu %v", got, tt.want)
			} else {
				t.Logf("✅ Détection correcte: %v", got)
			}
		})
	}
}
