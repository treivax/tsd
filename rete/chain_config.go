// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"time"
)

// CacheEvictionPolicy définit la politique d'éviction du cache
type CacheEvictionPolicy string

const (
	// EvictionPolicyNone pas d'éviction automatique
	EvictionPolicyNone CacheEvictionPolicy = "none"
	// EvictionPolicyLRU éviction Least Recently Used
	EvictionPolicyLRU CacheEvictionPolicy = "lru"
	// EvictionPolicyLFU éviction Least Frequently Used
	EvictionPolicyLFU CacheEvictionPolicy = "lfu"
)

// ChainPerformanceConfig configure les optimisations de performance pour les chaînes alpha
type ChainPerformanceConfig struct {
	// Cache de Hash
	HashCacheEnabled  bool                `json:"hash_cache_enabled"`
	HashCacheMaxSize  int                 `json:"hash_cache_max_size"`
	HashCacheEviction CacheEvictionPolicy `json:"hash_cache_eviction"`
	HashCacheTTL      time.Duration       `json:"hash_cache_ttl,omitempty"` // 0 = pas d'expiration

	// Cache de Connexion
	ConnectionCacheEnabled  bool                `json:"connection_cache_enabled"`
	ConnectionCacheMaxSize  int                 `json:"connection_cache_max_size"`
	ConnectionCacheEviction CacheEvictionPolicy `json:"connection_cache_eviction"`
	ConnectionCacheTTL      time.Duration       `json:"connection_cache_ttl,omitempty"`

	// Métriques
	MetricsEnabled         bool `json:"metrics_enabled"`
	MetricsDetailedChains  bool `json:"metrics_detailed_chains"`   // Stocker les détails de chaque chaîne
	MetricsMaxChainDetails int  `json:"metrics_max_chain_details"` // Limite de détails stockés

	// Performance
	ParallelHashComputation bool `json:"parallel_hash_computation"` // Calcul parallèle des hash (expérimental)

	// Monitoring
	PrometheusEnabled bool   `json:"prometheus_enabled"`
	PrometheusPrefix  string `json:"prometheus_prefix"` // Préfixe pour les métriques Prometheus

	// Beta Cache - Cache pour les opérations de jointure
	BetaCacheEnabled           bool                `json:"beta_cache_enabled"`
	BetaHashCacheMaxSize       int                 `json:"beta_hash_cache_max_size"`
	BetaHashCacheEviction      CacheEvictionPolicy `json:"beta_hash_cache_eviction"`
	BetaHashCacheTTL           time.Duration       `json:"beta_hash_cache_ttl,omitempty"`
	BetaJoinResultCacheEnabled bool                `json:"beta_join_result_cache_enabled"`
	BetaJoinResultCacheMaxSize int                 `json:"beta_join_result_cache_max_size"`
	BetaJoinResultCacheTTL     time.Duration       `json:"beta_join_result_cache_ttl"`
}

// DefaultChainPerformanceConfig retourne la configuration par défaut
func DefaultChainPerformanceConfig() *ChainPerformanceConfig {
	return &ChainPerformanceConfig{
		// Cache de Hash - activé avec LRU
		HashCacheEnabled:  true,
		HashCacheMaxSize:  10000, // 10k entrées max
		HashCacheEviction: EvictionPolicyLRU,
		HashCacheTTL:      0, // Pas d'expiration

		// Cache de Connexion - activé avec limite
		ConnectionCacheEnabled:  true,
		ConnectionCacheMaxSize:  50000, // 50k entrées max
		ConnectionCacheEviction: EvictionPolicyLRU,
		ConnectionCacheTTL:      0,

		// Métriques - activées
		MetricsEnabled:         true,
		MetricsDetailedChains:  true,
		MetricsMaxChainDetails: 1000, // Garder les 1000 dernières chaînes

		// Performance
		ParallelHashComputation: false, // Désactivé par défaut

		// Prometheus - désactivé par défaut
		PrometheusEnabled: false,
		PrometheusPrefix:  "tsd_rete",

		// Beta Cache - activé avec tailles raisonnables
		BetaCacheEnabled:           true,
		BetaHashCacheMaxSize:       10000, // 10k entrées pour les hash de jointure
		BetaHashCacheEviction:      EvictionPolicyLRU,
		BetaHashCacheTTL:           0, // Pas d'expiration
		BetaJoinResultCacheEnabled: true,
		BetaJoinResultCacheMaxSize: 5000,        // 5k résultats de jointure cachés
		BetaJoinResultCacheTTL:     time.Minute, // Expiration après 1 minute
	}
}

// HighPerformanceConfig retourne une configuration optimisée pour la performance
func HighPerformanceConfig() *ChainPerformanceConfig {
	return &ChainPerformanceConfig{
		HashCacheEnabled:  true,
		HashCacheMaxSize:  100000, // 100k entrées
		HashCacheEviction: EvictionPolicyLRU,
		HashCacheTTL:      0,

		ConnectionCacheEnabled:  true,
		ConnectionCacheMaxSize:  200000, // 200k entrées
		ConnectionCacheEviction: EvictionPolicyLRU,
		ConnectionCacheTTL:      0,

		MetricsEnabled:         true,
		MetricsDetailedChains:  false, // Désactiver pour économiser mémoire
		MetricsMaxChainDetails: 0,

		ParallelHashComputation: true, // Activé

		PrometheusEnabled: true,
		PrometheusPrefix:  "tsd_rete",

		// Beta Cache - configuration haute performance
		BetaCacheEnabled:           true,
		BetaHashCacheMaxSize:       100000, // 100k entrées
		BetaHashCacheEviction:      EvictionPolicyLRU,
		BetaHashCacheTTL:           0,
		BetaJoinResultCacheEnabled: true,
		BetaJoinResultCacheMaxSize: 50000,           // 50k résultats cachés
		BetaJoinResultCacheTTL:     5 * time.Minute, // TTL plus long

	}
}

// LowMemoryConfig retourne une configuration optimisée pour la mémoire
func LowMemoryConfig() *ChainPerformanceConfig {
	return &ChainPerformanceConfig{
		HashCacheEnabled:  true,
		HashCacheMaxSize:  1000, // Seulement 1k entrées
		HashCacheEviction: EvictionPolicyLRU,
		HashCacheTTL:      5 * time.Minute, // Expiration après 5min

		ConnectionCacheEnabled:  true,
		ConnectionCacheMaxSize:  5000,
		ConnectionCacheEviction: EvictionPolicyLRU,
		ConnectionCacheTTL:      5 * time.Minute,

		MetricsEnabled:         true,
		MetricsDetailedChains:  false,
		MetricsMaxChainDetails: 0,

		ParallelHashComputation: false,

		PrometheusEnabled: false,
		PrometheusPrefix:  "tsd_rete",

		// Beta Cache - configuration light
		BetaCacheEnabled:           true,
		BetaHashCacheMaxSize:       1000, // 1k entrées seulement
		BetaHashCacheEviction:      EvictionPolicyLRU,
		BetaHashCacheTTL:           0,
		BetaJoinResultCacheEnabled: false, // Désactivé en mode léger
		BetaJoinResultCacheMaxSize: 0,
		BetaJoinResultCacheTTL:     0,
	}
}

// DisabledCachesConfig retourne une configuration sans caches (pour tests/debug)
func DisabledCachesConfig() *ChainPerformanceConfig {
	return &ChainPerformanceConfig{
		HashCacheEnabled:  false,
		HashCacheMaxSize:  0,
		HashCacheEviction: EvictionPolicyNone,
		HashCacheTTL:      0,

		ConnectionCacheEnabled:  false,
		ConnectionCacheMaxSize:  0,
		ConnectionCacheEviction: EvictionPolicyNone,
		ConnectionCacheTTL:      0,

		MetricsEnabled:         true,
		MetricsDetailedChains:  true,
		MetricsMaxChainDetails: 1000,

		ParallelHashComputation: false,

		PrometheusEnabled: false,
		PrometheusPrefix:  "tsd_rete",
	}
}

// Validate vérifie que la configuration est valide
func (c *ChainPerformanceConfig) Validate() error {
	// Valider les tailles de cache
	if c.HashCacheEnabled {
		if c.HashCacheMaxSize < 0 {
			return fmt.Errorf("hash_cache_max_size ne peut pas être négatif: %d", c.HashCacheMaxSize)
		}
		if c.HashCacheMaxSize == 0 {
			return fmt.Errorf("hash_cache_max_size doit être > 0 quand le cache est activé")
		}
		if c.HashCacheMaxSize > 1000000 {
			return fmt.Errorf("hash_cache_max_size trop grand: %d (max: 1000000)", c.HashCacheMaxSize)
		}
	}

	if c.ConnectionCacheEnabled {
		if c.ConnectionCacheMaxSize < 0 {
			return fmt.Errorf("connection_cache_max_size ne peut pas être négatif: %d", c.ConnectionCacheMaxSize)
		}
		if c.ConnectionCacheMaxSize == 0 {
			return fmt.Errorf("connection_cache_max_size doit être > 0 quand le cache est activé")
		}
		if c.ConnectionCacheMaxSize > 10000000 {
			return fmt.Errorf("connection_cache_max_size trop grand: %d (max: 10000000)", c.ConnectionCacheMaxSize)
		}
	}

	// Valider les politiques d'éviction
	validPolicies := map[CacheEvictionPolicy]bool{
		EvictionPolicyNone: true,
		EvictionPolicyLRU:  true,
		EvictionPolicyLFU:  true,
	}

	if !validPolicies[c.HashCacheEviction] {
		return fmt.Errorf("politique d'éviction hash invalide: %s", c.HashCacheEviction)
	}

	if !validPolicies[c.ConnectionCacheEviction] {
		return fmt.Errorf("politique d'éviction connexion invalide: %s", c.ConnectionCacheEviction)
	}

	// Valider les TTL
	if c.HashCacheTTL < 0 {
		return fmt.Errorf("hash_cache_ttl ne peut pas être négatif: %v", c.HashCacheTTL)
	}

	if c.ConnectionCacheTTL < 0 {
		return fmt.Errorf("connection_cache_ttl ne peut pas être négatif: %v", c.ConnectionCacheTTL)
	}

	// Valider les métriques
	if c.MetricsEnabled && c.MetricsDetailedChains {
		if c.MetricsMaxChainDetails < 0 {
			return fmt.Errorf("metrics_max_chain_details ne peut pas être négatif: %d", c.MetricsMaxChainDetails)
		}
		if c.MetricsMaxChainDetails > 100000 {
			return fmt.Errorf("metrics_max_chain_details trop grand: %d (max: 100000)", c.MetricsMaxChainDetails)
		}
	}

	// Valider le préfixe Prometheus
	if c.PrometheusEnabled && c.PrometheusPrefix == "" {
		return fmt.Errorf("prometheus_prefix ne peut pas être vide quand Prometheus est activé")
	}

	// Valider Beta Cache
	if c.BetaCacheEnabled {
		if c.BetaHashCacheMaxSize <= 0 {
			return fmt.Errorf("BetaHashCacheMaxSize doit être > 0 quand BetaCacheEnabled=true")
		}
		if c.BetaJoinResultCacheEnabled && c.BetaJoinResultCacheMaxSize <= 0 {
			return fmt.Errorf("BetaJoinResultCacheMaxSize doit être > 0 quand BetaJoinResultCacheEnabled=true")
		}
	}

	return nil
}

// Clone crée une copie profonde de la configuration
func (c *ChainPerformanceConfig) Clone() *ChainPerformanceConfig {
	if c == nil {
		return nil
	}

	clone := *c
	return &clone
}

// GetCacheInfo retourne des informations sur la configuration des caches
func (c *ChainPerformanceConfig) GetCacheInfo() map[string]interface{} {
	return map[string]interface{}{
		"hash_cache": map[string]interface{}{
			"enabled":  c.HashCacheEnabled,
			"max_size": c.HashCacheMaxSize,
			"eviction": string(c.HashCacheEviction),
			"ttl":      c.HashCacheTTL.String(),
		},
		"connection_cache": map[string]interface{}{
			"enabled":  c.ConnectionCacheEnabled,
			"max_size": c.ConnectionCacheMaxSize,
			"eviction": string(c.ConnectionCacheEviction),
			"ttl":      c.ConnectionCacheTTL.String(),
		},
		"metrics": map[string]interface{}{
			"enabled":           c.MetricsEnabled,
			"detailed_chains":   c.MetricsDetailedChains,
			"max_chain_details": c.MetricsMaxChainDetails,
		},
		"performance": map[string]interface{}{
			"parallel_hash": c.ParallelHashComputation,
		},
		"prometheus": map[string]interface{}{
			"enabled": c.PrometheusEnabled,
			"prefix":  c.PrometheusPrefix,
		},
	}
}

// EstimateMemoryUsage estime l'utilisation mémoire de la configuration (en bytes)
func (c *ChainPerformanceConfig) EstimateMemoryUsage() int64 {
	var total int64

	// Cache de hash: ~500 bytes par entrée (JSON condition + hash)
	if c.HashCacheEnabled {
		total += int64(c.HashCacheMaxSize) * 500
	}

	// Cache de connexion: ~100 bytes par entrée (2 IDs + bool)
	if c.ConnectionCacheEnabled {
		total += int64(c.ConnectionCacheMaxSize) * 100
	}

	// Détails de chaînes: ~200 bytes par entrée
	if c.MetricsEnabled && c.MetricsDetailedChains {
		total += int64(c.MetricsMaxChainDetails) * 200
	}

	return total
}

// String retourne une représentation textuelle de la configuration
func (c *ChainPerformanceConfig) String() string {
	return fmt.Sprintf("ChainPerformanceConfig{HashCache:%v(%d,%s), ConnCache:%v(%d,%s), Metrics:%v, Prometheus:%v}",
		c.HashCacheEnabled, c.HashCacheMaxSize, c.HashCacheEviction,
		c.ConnectionCacheEnabled, c.ConnectionCacheMaxSize, c.ConnectionCacheEviction,
		c.MetricsEnabled, c.PrometheusEnabled)
}
