// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"time"

	"github.com/treivax/tsd/rete"
)

func main() {
	fmt.Println("🔧 Exemple d'intégration du Cache LRU dans Alpha Sharing")
	fmt.Println("=========================================================")
	fmt.Println()

	// 1. Créer un réseau avec configuration par défaut (LRU activé)
	fmt.Println("1️⃣  Création d'un réseau avec configuration par défaut")
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)

	config := network.GetConfig()
	fmt.Printf("   ✓ Cache LRU activé: %v\n", config.HashCacheEnabled)
	fmt.Printf("   ✓ Capacité: %d entrées\n", config.HashCacheMaxSize)
	fmt.Printf("   ✓ Politique d'éviction: %s\n", config.HashCacheEviction)
	fmt.Printf("   ✓ TTL: %v\n\n", config.HashCacheTTL)

	// 2. Créer un réseau avec configuration haute performance
	fmt.Println("2️⃣  Création d'un réseau haute performance")
	highPerfConfig := rete.HighPerformanceConfig()
	networkHP := rete.NewReteNetworkWithConfig(storage, highPerfConfig)

	fmt.Printf("   ✓ Capacité: %d entrées\n", highPerfConfig.HashCacheMaxSize)
	fmt.Printf("   ✓ Estimation mémoire: %.2f MB\n", float64(highPerfConfig.EstimateMemoryUsage())/(1024*1024))
	fmt.Printf("   ✓ Métriques détaillées: %v\n", highPerfConfig.MetricsDetailedChains)
	fmt.Printf("   ✓ Réseau créé: %s\n\n", networkHP.RootNode.GetID())

	// 3. Créer un réseau avec configuration basse mémoire
	fmt.Println("3️⃣  Création d'un réseau basse mémoire")
	lowMemConfig := rete.LowMemoryConfig()
	networkLM := rete.NewReteNetworkWithConfig(storage, lowMemConfig)

	fmt.Printf("   ✓ Capacité: %d entrées\n", lowMemConfig.HashCacheMaxSize)
	fmt.Printf("   ✓ TTL: %v\n", lowMemConfig.HashCacheTTL)
	fmt.Printf("   ✓ Estimation mémoire: %.2f MB\n", float64(lowMemConfig.EstimateMemoryUsage())/(1024*1024))
	fmt.Printf("   ✓ Réseau créé: %s\n\n", networkLM.RootNode.GetID())

	// 4. Simuler l'utilisation du cache
	fmt.Println("4️⃣  Simulation d'utilisation du cache")
	registry := network.AlphaSharingManager

	// Créer plusieurs conditions similaires
	for i := 0; i < 100; i++ {
		condition := map[string]interface{}{
			"type":     "binaryOperation",
			"operator": "==",
			"left":     "age",
			"right":    18 + (i % 10), // 10 valeurs différentes répétées
		}

		_, err := registry.ConditionHashCached(condition, "p")
		if err != nil {
			fmt.Printf("   ✗ Erreur: %v\n", err)
			return
		}
	}

	fmt.Printf("   ✓ 100 conditions traitées\n")

	// 5. Afficher les statistiques du cache
	fmt.Println("\n5️⃣  Statistiques du cache LRU")
	stats := registry.GetHashCacheStats()

	fmt.Printf("   Type de cache: %v\n", stats["type"])
	fmt.Printf("   Taille actuelle: %v entrées\n", stats["size"])
	fmt.Printf("   Capacité: %v entrées\n", stats["capacity"])
	fmt.Printf("   Cache hits: %v\n", stats["hits"])
	fmt.Printf("   Cache misses: %v\n", stats["misses"])
	fmt.Printf("   Évictions: %v\n", stats["evictions"])
	fmt.Printf("   Hit rate: %.2f%%\n", stats["hit_rate"].(float64)*100)
	fmt.Printf("   Éviction rate: %.2f%%\n", stats["eviction_rate"].(float64)*100)
	fmt.Printf("   Fill rate: %.2f%%\n\n", stats["fill_rate"].(float64)*100)

	// 6. Afficher les métriques du réseau
	fmt.Println("6️⃣  Métriques du réseau")
	metrics := network.GetChainMetrics()
	summary := metrics.GetSummary()

	fmt.Printf("   Cache hits: %v\n", summary["hash_cache_hits"])
	fmt.Printf("   Cache misses: %v\n", summary["hash_cache_misses"])
	fmt.Printf("   Taille cache: %v\n\n", summary["hash_cache_size"])

	// 7. Démonstration du comportement LRU avec petite capacité
	fmt.Println("7️⃣  Démonstration de l'éviction LRU")
	smallConfig := rete.DefaultChainPerformanceConfig()
	smallConfig.HashCacheMaxSize = 5 // Très petite capacité
	smallRegistry := rete.NewAlphaSharingRegistryWithConfig(smallConfig, rete.NewChainBuildMetrics())

	fmt.Printf("   Capacité du cache: %d entrées\n", smallConfig.HashCacheMaxSize)
	fmt.Println("   Ajout de 10 conditions...")

	for i := 0; i < 10; i++ {
		condition := map[string]interface{}{
			"type":     "binaryOperation",
			"operator": "==",
			"left":     fmt.Sprintf("field%d", i),
			"right":    i,
		}

		_, err := smallRegistry.ConditionHashCached(condition, "p")
		if err != nil {
			fmt.Printf("   ✗ Erreur: %v\n", err)
			return
		}
	}

	smallStats := smallRegistry.GetHashCacheStats()
	fmt.Printf("   ✓ Taille finale du cache: %v entrées (limité par capacité)\n", smallStats["size"])
	fmt.Printf("   ✓ Évictions: %v (10 - 5 = 5 évictions attendues)\n\n", smallStats["evictions"])

	// 8. Démonstration du TTL
	fmt.Println("8️⃣  Démonstration du TTL (expiration)")
	ttlConfig := rete.DefaultChainPerformanceConfig()
	ttlConfig.HashCacheTTL = 500 * time.Millisecond
	ttlRegistry := rete.NewAlphaSharingRegistryWithConfig(ttlConfig, rete.NewChainBuildMetrics())

	condition := map[string]interface{}{
		"type":     "binaryOperation",
		"operator": "==",
		"left":     "status",
		"right":    "active",
	}

	fmt.Printf("   TTL configuré: %v\n", ttlConfig.HashCacheTTL)
	fmt.Println("   Ajout d'une condition...")
	_, _ = ttlRegistry.ConditionHashCached(condition, "p")

	fmt.Println("   Vérification immédiate (devrait être un hit)...")
	_, _ = ttlRegistry.ConditionHashCached(condition, "p")

	ttlStats1 := ttlRegistry.GetHashCacheStats()
	fmt.Printf("   ✓ Hits: %v\n", ttlStats1["hits"])

	fmt.Println("   Attente de l'expiration (600ms)...")
	time.Sleep(600 * time.Millisecond)

	fmt.Println("   Vérification après expiration (devrait être un miss)...")
	_, _ = ttlRegistry.ConditionHashCached(condition, "p")

	ttlStats2 := ttlRegistry.GetHashCacheStats()
	fmt.Printf("   ✓ Misses après expiration: %v\n\n", ttlStats2["misses"])

	// 9. Nettoyage des entrées expirées
	fmt.Println("9️⃣  Nettoyage des entrées expirées")
	for i := 0; i < 5; i++ {
		cond := map[string]interface{}{
			"type":     "binaryOperation",
			"operator": "==",
			"left":     fmt.Sprintf("temp%d", i),
			"right":    i,
		}
		_, _ = ttlRegistry.ConditionHashCached(cond, "p")
	}

	sizeBefore := ttlRegistry.GetHashCacheSize()
	fmt.Printf("   Taille avant nettoyage: %d entrées\n", sizeBefore)

	time.Sleep(600 * time.Millisecond)

	cleaned := ttlRegistry.CleanExpiredHashCache()
	fmt.Printf("   ✓ Entrées nettoyées: %d\n", cleaned)
	fmt.Printf("   ✓ Taille après nettoyage: %d entrées\n\n", ttlRegistry.GetHashCacheSize())

	// 10. Configuration personnalisée
	fmt.Println("🔟 Configuration personnalisée")
	customConfig := rete.DefaultChainPerformanceConfig()
	customConfig.HashCacheMaxSize = 25000
	customConfig.HashCacheTTL = 10 * time.Minute
	customConfig.MetricsEnabled = true
	customConfig.MetricsDetailedChains = true

	fmt.Println("   Configuration personnalisée:")
	fmt.Printf("   - Capacité: %d entrées\n", customConfig.HashCacheMaxSize)
	fmt.Printf("   - TTL: %v\n", customConfig.HashCacheTTL)
	fmt.Printf("   - Métriques: %v\n", customConfig.MetricsEnabled)
	fmt.Printf("   - Estimation mémoire: %.2f MB\n", float64(customConfig.EstimateMemoryUsage())/(1024*1024))

	// Validation de la configuration
	if err := customConfig.Validate(); err != nil {
		fmt.Printf("   ✗ Configuration invalide: %v\n", err)
	} else {
		fmt.Println("   ✓ Configuration valide")
		fmt.Println()
	}

	// 11. Comparaison des configurations
	fmt.Println("1️⃣1️⃣  Comparaison des configurations")
	fmt.Println("\n   Configuration           | Capacité  | TTL        | Mémoire (MB)")
	fmt.Println("   ----------------------- | --------- | ---------- | ------------")

	configs := map[string]*rete.ChainPerformanceConfig{
		"Par défaut":        rete.DefaultChainPerformanceConfig(),
		"Haute performance": rete.HighPerformanceConfig(),
		"Basse mémoire":     rete.LowMemoryConfig(),
	}

	for name, cfg := range configs {
		memMB := float64(cfg.EstimateMemoryUsage()) / (1024 * 1024)
		ttlStr := "Aucun"
		if cfg.HashCacheTTL > 0 {
			ttlStr = cfg.HashCacheTTL.String()
		}
		fmt.Printf("   %-23s | %-9d | %-10s | %.2f\n", name, cfg.HashCacheMaxSize, ttlStr, memMB)
	}

	// 12. Conclusion
	fmt.Println("\n" + repeat("=", 60))
	fmt.Println("✅ Démonstration terminée avec succès!")
	fmt.Println("\nPoints clés:")
	fmt.Println("  • Le cache LRU est automatiquement activé avec la configuration par défaut")
	fmt.Println("  • Trois configurations prédéfinies disponibles (default, high-perf, low-mem)")
	fmt.Println("  • Contrôle fin de la mémoire via capacité et TTL")
	fmt.Println("  • Statistiques détaillées pour le monitoring")
	fmt.Println("  • Éviction LRU automatique quand la capacité est atteinte")
	fmt.Println("  • Expiration TTL optionnelle pour les environnements contraints")
	fmt.Println("  • Thread-safe et production-ready")
}

// repeat répète une chaîne n fois
func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
