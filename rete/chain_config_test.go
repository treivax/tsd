// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
	"time"
)

// TestDefaultChainPerformanceConfig teste la configuration par défaut
func TestDefaultChainPerformanceConfig(t *testing.T) {
	config := DefaultChainPerformanceConfig()

	if config == nil {
		t.Fatal("DefaultChainPerformanceConfig retourne nil")
	}

	// Vérifier les valeurs par défaut
	if !config.HashCacheEnabled {
		t.Error("HashCacheEnabled devrait être true par défaut")
	}
	if config.HashCacheMaxSize != 10000 {
		t.Errorf("HashCacheMaxSize devrait être 10000, obtenu %d", config.HashCacheMaxSize)
	}
	if config.HashCacheEviction != EvictionPolicyLRU {
		t.Errorf("HashCacheEviction devrait être LRU, obtenu %s", config.HashCacheEviction)
	}

	if !config.ConnectionCacheEnabled {
		t.Error("ConnectionCacheEnabled devrait être true par défaut")
	}
	if config.ConnectionCacheMaxSize != 50000 {
		t.Errorf("ConnectionCacheMaxSize devrait être 50000, obtenu %d", config.ConnectionCacheMaxSize)
	}

	if !config.MetricsEnabled {
		t.Error("MetricsEnabled devrait être true par défaut")
	}
	if !config.MetricsDetailedChains {
		t.Error("MetricsDetailedChains devrait être true par défaut")
	}
	if config.MetricsMaxChainDetails != 1000 {
		t.Errorf("MetricsMaxChainDetails devrait être 1000, obtenu %d", config.MetricsMaxChainDetails)
	}

	if config.PrometheusEnabled {
		t.Error("PrometheusEnabled devrait être false par défaut")
	}
	if config.PrometheusPrefix != "tsd_rete" {
		t.Errorf("PrometheusPrefix devrait être 'tsd_rete', obtenu %s", config.PrometheusPrefix)
	}
}

// TestHighPerformanceConfig teste la configuration haute performance
func TestHighPerformanceConfig(t *testing.T) {
	config := HighPerformanceConfig()

	if config.HashCacheMaxSize != 100000 {
		t.Errorf("HashCacheMaxSize devrait être 100000, obtenu %d", config.HashCacheMaxSize)
	}
	if config.ConnectionCacheMaxSize != 200000 {
		t.Errorf("ConnectionCacheMaxSize devrait être 200000, obtenu %d", config.ConnectionCacheMaxSize)
	}
	if config.MetricsDetailedChains {
		t.Error("MetricsDetailedChains devrait être false en haute performance")
	}
	if !config.ParallelHashComputation {
		t.Error("ParallelHashComputation devrait être true en haute performance")
	}
	if !config.PrometheusEnabled {
		t.Error("PrometheusEnabled devrait être true en haute performance")
	}
}

// TestLowMemoryConfig teste la configuration mémoire réduite
func TestLowMemoryConfig(t *testing.T) {
	config := LowMemoryConfig()

	if config.HashCacheMaxSize != 1000 {
		t.Errorf("HashCacheMaxSize devrait être 1000, obtenu %d", config.HashCacheMaxSize)
	}
	if config.ConnectionCacheMaxSize != 5000 {
		t.Errorf("ConnectionCacheMaxSize devrait être 5000, obtenu %d", config.ConnectionCacheMaxSize)
	}
	if config.HashCacheTTL != 5*time.Minute {
		t.Errorf("HashCacheTTL devrait être 5min, obtenu %v", config.HashCacheTTL)
	}
	if config.MetricsDetailedChains {
		t.Error("MetricsDetailedChains devrait être false en mémoire réduite")
	}
}

// TestDisabledCachesConfig teste la configuration sans caches
func TestDisabledCachesConfig(t *testing.T) {
	config := DisabledCachesConfig()

	if config.HashCacheEnabled {
		t.Error("HashCacheEnabled devrait être false")
	}
	if config.ConnectionCacheEnabled {
		t.Error("ConnectionCacheEnabled devrait être false")
	}
	if config.HashCacheMaxSize != 0 {
		t.Errorf("HashCacheMaxSize devrait être 0, obtenu %d", config.HashCacheMaxSize)
	}
}

// TestChainPerformanceConfig_Validate teste la validation
func TestChainPerformanceConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    *ChainPerformanceConfig
		shouldErr bool
	}{
		{
			name:      "configuration par défaut valide",
			config:    DefaultChainPerformanceConfig(),
			shouldErr: false,
		},
		{
			name: "taille hash cache négative",
			config: &ChainPerformanceConfig{
				HashCacheEnabled: true,
				HashCacheMaxSize: -1,
			},
			shouldErr: true,
		},
		{
			name: "taille hash cache zéro avec cache activé",
			config: &ChainPerformanceConfig{
				HashCacheEnabled: true,
				HashCacheMaxSize: 0,
			},
			shouldErr: true,
		},
		{
			name: "taille hash cache trop grande",
			config: &ChainPerformanceConfig{
				HashCacheEnabled: true,
				HashCacheMaxSize: 2000000,
			},
			shouldErr: true,
		},
		{
			name: "politique d'éviction invalide",
			config: &ChainPerformanceConfig{
				HashCacheEnabled:  true,
				HashCacheMaxSize:  1000,
				HashCacheEviction: "invalid",
			},
			shouldErr: true,
		},
		{
			name: "TTL négatif",
			config: &ChainPerformanceConfig{
				HashCacheEnabled: true,
				HashCacheMaxSize: 1000,
				HashCacheTTL:     -1 * time.Second,
			},
			shouldErr: true,
		},
		{
			name: "prometheus activé sans préfixe",
			config: &ChainPerformanceConfig{
				HashCacheEnabled:  true,
				HashCacheMaxSize:  1000,
				PrometheusEnabled: true,
				PrometheusPrefix:  "",
			},
			shouldErr: true,
		},
		{
			name: "métriques détaillées avec limite négative",
			config: &ChainPerformanceConfig{
				HashCacheEnabled:       true,
				HashCacheMaxSize:       1000,
				MetricsEnabled:         true,
				MetricsDetailedChains:  true,
				MetricsMaxChainDetails: -1,
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.shouldErr && err == nil {
				t.Error("Attendu une erreur mais aucune erreur retournée")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Erreur inattendue: %v", err)
			}
		})
	}
}

// TestChainPerformanceConfig_Clone teste le clonage
func TestChainPerformanceConfig_Clone(t *testing.T) {
	original := DefaultChainPerformanceConfig()
	original.HashCacheMaxSize = 12345

	clone := original.Clone()

	if clone == nil {
		t.Fatal("Clone retourne nil")
	}

	// Vérifier que c'est une copie
	if clone == original {
		t.Error("Clone retourne le même pointeur")
	}

	// Vérifier que les valeurs sont identiques
	if clone.HashCacheMaxSize != original.HashCacheMaxSize {
		t.Error("Les valeurs du clone ne correspondent pas")
	}

	// Modifier le clone ne doit pas affecter l'original
	clone.HashCacheMaxSize = 99999
	if original.HashCacheMaxSize == 99999 {
		t.Error("La modification du clone affecte l'original")
	}
}

// TestChainPerformanceConfig_GetCacheInfo teste GetCacheInfo
func TestChainPerformanceConfig_GetCacheInfo(t *testing.T) {
	config := DefaultChainPerformanceConfig()
	info := config.GetCacheInfo()

	if info == nil {
		t.Fatal("GetCacheInfo retourne nil")
	}

	// Vérifier la structure
	if _, ok := info["hash_cache"]; !ok {
		t.Error("hash_cache manquant dans l'info")
	}
	if _, ok := info["connection_cache"]; !ok {
		t.Error("connection_cache manquant dans l'info")
	}
	if _, ok := info["metrics"]; !ok {
		t.Error("metrics manquant dans l'info")
	}
	if _, ok := info["performance"]; !ok {
		t.Error("performance manquant dans l'info")
	}
	if _, ok := info["prometheus"]; !ok {
		t.Error("prometheus manquant dans l'info")
	}
}

// TestChainPerformanceConfig_EstimateMemoryUsage teste l'estimation mémoire
func TestChainPerformanceConfig_EstimateMemoryUsage(t *testing.T) {
	tests := []struct {
		name     string
		config   *ChainPerformanceConfig
		minBytes int64
		maxBytes int64
	}{
		{
			name:     "configuration par défaut",
			config:   DefaultChainPerformanceConfig(),
			minBytes: 5000000,  // ~5MB minimum
			maxBytes: 11000000, // ~11MB maximum
		},
		{
			name:     "configuration haute performance",
			config:   HighPerformanceConfig(),
			minBytes: 50000000,  // ~50MB minimum
			maxBytes: 100000000, // ~100MB maximum
		},
		{
			name:     "configuration mémoire réduite",
			config:   LowMemoryConfig(),
			minBytes: 500000,  // ~500KB minimum
			maxBytes: 1000000, // ~1MB maximum
		},
		{
			name:     "caches désactivés",
			config:   DisabledCachesConfig(),
			minBytes: 0,
			maxBytes: 300000, // ~300KB pour les métriques
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usage := tt.config.EstimateMemoryUsage()
			if usage < tt.minBytes {
				t.Errorf("Utilisation mémoire trop faible: %d bytes (min: %d)", usage, tt.minBytes)
			}
			if usage > tt.maxBytes {
				t.Errorf("Utilisation mémoire trop élevée: %d bytes (max: %d)", usage, tt.maxBytes)
			}
			t.Logf("Utilisation mémoire estimée: %d bytes (~%.2f MB)", usage, float64(usage)/1024/1024)
		})
	}
}

// TestChainPerformanceConfig_String teste la représentation textuelle
func TestChainPerformanceConfig_String(t *testing.T) {
	config := DefaultChainPerformanceConfig()
	str := config.String()

	if str == "" {
		t.Error("String() retourne une chaîne vide")
	}

	// Vérifier que certains éléments clés sont présents
	if len(str) < 50 {
		t.Errorf("String() trop court: %s", str)
	}

	t.Logf("String representation: %s", str)
}

// TestCacheEvictionPolicy teste les constantes de politique d'éviction
func TestCacheEvictionPolicy(t *testing.T) {
	policies := []CacheEvictionPolicy{
		EvictionPolicyNone,
		EvictionPolicyLRU,
		EvictionPolicyLFU,
	}

	for _, policy := range policies {
		if string(policy) == "" {
			t.Errorf("Politique d'éviction vide: %v", policy)
		}
	}

	// Vérifier que les valeurs sont différentes
	if EvictionPolicyNone == EvictionPolicyLRU {
		t.Error("None et LRU ont la même valeur")
	}
	if EvictionPolicyLRU == EvictionPolicyLFU {
		t.Error("LRU et LFU ont la même valeur")
	}
}

// TestChainPerformanceConfig_NilClone teste le clonage d'une config nil
func TestChainPerformanceConfig_NilClone(t *testing.T) {
	var config *ChainPerformanceConfig
	clone := config.Clone()

	if clone != nil {
		t.Error("Clone d'un pointeur nil devrait retourner nil")
	}
}
