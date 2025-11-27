# Exemple d'Intégration du Cache LRU

Ce programme démontre l'utilisation du cache LRU (Least Recently Used) intégré dans le système de partage des AlphaNodes du réseau RETE.

## 🎯 Objectif

Illustrer les fonctionnalités du cache LRU :
- Configuration avec différents presets (default, high-performance, low-memory)
- Éviction automatique des entrées les moins récemment utilisées
- Expiration TTL (Time To Live)
- Statistiques détaillées de performance
- Nettoyage manuel des entrées expirées

## 🚀 Exécution

```bash
cd tsd
go run examples/lru_cache/main.go
```

## 📋 Ce que l'exemple démontre

### 1. Configuration par défaut
- Cache LRU activé automatiquement
- Capacité : 10,000 entrées
- Pas d'expiration TTL

### 2. Configuration haute performance
- Capacité : 100,000 entrées
- Estimation mémoire : ~67 MB
- Métriques détaillées désactivées pour économiser mémoire

### 3. Configuration basse mémoire
- Capacité : 1,000 entrées
- TTL : 5 minutes
- Estimation mémoire : ~0.95 MB

### 4. Simulation d'utilisation
- 100 conditions traitées
- 10 valeurs uniques répétées
- Hit rate attendu : ~90%

### 5. Statistiques du cache
Affiche :
- Type de cache (LRU)
- Taille actuelle vs capacité
- Nombre de hits/misses
- Taux de performance (hit rate, eviction rate, fill rate)

### 6. Métriques du réseau
Intégration avec les métriques globales du réseau RETE

### 7. Démonstration de l'éviction LRU
- Cache limité à 5 entrées
- Ajout de 10 conditions
- Vérification que les 5 premières sont évincées

### 8. Démonstration du TTL
- TTL configuré à 500ms
- Vérification de l'expiration automatique
- Cache hit puis cache miss après expiration

### 9. Nettoyage des entrées expirées
- Ajout de plusieurs entrées
- Attente de l'expiration
- Nettoyage manuel avec `CleanExpiredHashCache()`

### 10. Configuration personnalisée
- Création d'une configuration sur mesure
- Validation de la configuration
- Estimation de l'utilisation mémoire

### 11. Comparaison des configurations
Tableau comparatif :
- Configuration par défaut
- Haute performance
- Basse mémoire

## 📊 Sortie Attendue

```
🔧 Exemple d'intégration du Cache LRU dans Alpha Sharing
=========================================================

1️⃣  Création d'un réseau avec configuration par défaut
   ✓ Cache LRU activé: true
   ✓ Capacité: 10000 entrées
   ✓ Politique d'éviction: lru
   ✓ TTL: 0s

...

5️⃣  Statistiques du cache LRU
   Type de cache: lru
   Taille actuelle: 10 entrées
   Capacité: 10000 entrées
   Cache hits: 90
   Cache misses: 10
   Évictions: 0
   Hit rate: 90.00%
   Éviction rate: 0.00%
   Fill rate: 0.10%

...

7️⃣  Démonstration de l'éviction LRU
   Capacité du cache: 5 entrées
   Ajout de 10 conditions...
   ✓ Taille finale du cache: 5 entrées (limité par capacité)
   ✓ Évictions: 5 (10 - 5 = 5 évictions attendues)

...
```

## 🎓 Concepts Démontrés

### Cache LRU (Least Recently Used)
- Conserve les N entrées les plus récemment utilisées
- Évince automatiquement les entrées les moins récentes
- Optimal pour les patterns d'accès avec localité temporelle

### TTL (Time To Live)
- Expiration automatique après une durée configurée
- Utile pour les environnements contraints en mémoire
- Nettoyage manuel ou automatique à chaque accès

### Statistiques de Performance
- **Hit rate** : Pourcentage de cache hits (90% dans l'exemple)
- **Eviction rate** : Taux d'évictions par rapport aux insertions
- **Fill rate** : Pourcentage de remplissage du cache

## 🔧 Utilisation dans Votre Code

### Création d'un réseau avec LRU (configuration par défaut)
```go
storage := rete.NewMemoryStorage()
network := rete.NewReteNetwork(storage)
// Le cache LRU est automatiquement activé
```

### Création avec configuration personnalisée
```go
config := rete.HighPerformanceConfig()
network := rete.NewReteNetworkWithConfig(storage, config)
```

### Accès aux statistiques
```go
stats := network.AlphaSharingManager.GetHashCacheStats()
fmt.Printf("Hit rate: %.2f%%\n", stats["hit_rate"].(float64) * 100)
```

### Configuration personnalisée
```go
config := rete.DefaultChainPerformanceConfig()
config.HashCacheMaxSize = 25000
config.HashCacheTTL = 10 * time.Minute
config.MetricsEnabled = true

if err := config.Validate(); err != nil {
    log.Fatal(err)
}

network := rete.NewReteNetworkWithConfig(storage, config)
```

## 📚 Documentation Complémentaire

- [`rete/LRU_INTEGRATION_SUMMARY.md`](../../rete/LRU_INTEGRATION_SUMMARY.md) - Documentation complète
- [`rete/CHANGELOG_LRU_INTEGRATION.md`](../../rete/CHANGELOG_LRU_INTEGRATION.md) - Changelog détaillé
- [`rete/PERFORMANCE_QUICKSTART.md`](../../rete/PERFORMANCE_QUICKSTART.md) - Guide de performance
- [`docs/PROMETHEUS_INTEGRATION.md`](../../docs/PROMETHEUS_INTEGRATION.md) - Intégration Prometheus

## ⚡ Points Clés

1. **Automatique** : Le cache LRU est activé par défaut, aucune configuration requise
2. **Configurable** : Trois presets disponibles + configuration personnalisée
3. **Performant** : Hit rate typique de 90% sur patterns courants
4. **Contrôlé** : Limite stricte de mémoire avec éviction LRU
5. **Observable** : Statistiques détaillées pour le monitoring
6. **Thread-safe** : Utilisable en environnement concurrent

## 🐛 Dépannage

### Le cache n'a aucun hit
- Vérifiez que `HashCacheEnabled` est `true`
- Vérifiez que vous utilisez des conditions identiques
- Les conditions doivent être strictement identiques (même structure JSON)

### Trop d'évictions
- Augmentez `HashCacheMaxSize`
- Utilisez `HighPerformanceConfig()`
- Vérifiez les patterns d'accès (localité temporelle faible ?)

### Utilisation mémoire élevée
- Réduisez `HashCacheMaxSize`
- Utilisez `LowMemoryConfig()`
- Activez le TTL pour expiration automatique

## 📞 Support

Pour plus d'informations, consultez :
- La documentation du projet TSD
- Les tests d'intégration dans `rete/alpha_sharing_lru_integration_test.go`
- Le code source dans `rete/alpha_sharing.go`

---

*Exemple créé le : 2025-01-27*
*Version TSD : avec cache LRU intégré*