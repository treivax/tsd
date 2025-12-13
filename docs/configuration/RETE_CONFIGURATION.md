# Configuration du Réseau RETE

Ce document décrit toutes les configurations possibles du réseau RETE dans TSD, leurs cas d'usage, et les optimisations disponibles.

## Table des matières

1. [Configuration de Performance des Chaînes](#configuration-de-performance-des-chaînes)
2. [Configuration des Transactions](#configuration-des-transactions)
3. [Configuration du Beta Sharing](#configuration-du-beta-sharing)
4. [Configuration du Cache Arithmétique](#configuration-du-cache-arithmétique)
5. [Configuration du Logger](#configuration-du-logger)
6. [Profils de Configuration Prédéfinis](#profils-de-configuration-prédéfinis)
7. [Exemples d'Utilisation](#exemples-dutilisation)

---

## Configuration de Performance des Chaînes

La structure `ChainPerformanceConfig` contrôle les optimisations de performance pour les chaînes alpha et beta.

### Structure Complète

```go
type ChainPerformanceConfig struct {
    // ═══════════════════════════════════════════════════════════
    // Cache de Hash (Chaînes Alpha)
    // ═══════════════════════════════════════════════════════════
    HashCacheEnabled  bool                `json:"hash_cache_enabled"`
    HashCacheMaxSize  int                 `json:"hash_cache_max_size"`
    HashCacheEviction CacheEvictionPolicy `json:"hash_cache_eviction"`
    HashCacheTTL      time.Duration       `json:"hash_cache_ttl,omitempty"`

    // ═══════════════════════════════════════════════════════════
    // Cache de Connexion (Relations entre Nœuds)
    // ═══════════════════════════════════════════════════════════
    ConnectionCacheEnabled  bool                `json:"connection_cache_enabled"`
    ConnectionCacheMaxSize  int                 `json:"connection_cache_max_size"`
    ConnectionCacheEviction CacheEvictionPolicy `json:"connection_cache_eviction"`
    ConnectionCacheTTL      time.Duration       `json:"connection_cache_ttl,omitempty"`

    // ═══════════════════════════════════════════════════════════
    // Métriques et Monitoring
    // ═══════════════════════════════════════════════════════════
    MetricsEnabled         bool `json:"metrics_enabled"`
    MetricsDetailedChains  bool `json:"metrics_detailed_chains"`
    MetricsMaxChainDetails int  `json:"metrics_max_chain_details"`

    // ═══════════════════════════════════════════════════════════
    // Optimisations de Performance
    // ═══════════════════════════════════════════════════════════
    ParallelHashComputation bool `json:"parallel_hash_computation"`

    // ═══════════════════════════════════════════════════════════
    // Prometheus (Exportation de Métriques)
    // ═══════════════════════════════════════════════════════════
    PrometheusEnabled bool   `json:"prometheus_enabled"`
    PrometheusPrefix  string `json:"prometheus_prefix"`

    // ═══════════════════════════════════════════════════════════
    // Beta Cache (Opérations de Jointure)
    // ═══════════════════════════════════════════════════════════
    BetaCacheEnabled           bool                `json:"beta_cache_enabled"`
    BetaHashCacheMaxSize       int                 `json:"beta_hash_cache_max_size"`
    BetaHashCacheEviction      CacheEvictionPolicy `json:"beta_hash_cache_eviction"`
    BetaHashCacheTTL           time.Duration       `json:"beta_hash_cache_ttl,omitempty"`
    BetaJoinResultCacheEnabled bool                `json:"beta_join_result_cache_enabled"`
    BetaJoinResultCacheMaxSize int                 `json:"beta_join_result_cache_max_size"`
    BetaJoinResultCacheTTL     time.Duration       `json:"beta_join_result_cache_ttl"`
}
```

### Politiques d'Éviction de Cache

```go
type CacheEvictionPolicy string

const (
    EvictionPolicyNone CacheEvictionPolicy = "none"  // Pas d'éviction automatique
    EvictionPolicyLRU  CacheEvictionPolicy = "lru"   // Least Recently Used
    EvictionPolicyLFU  CacheEvictionPolicy = "lfu"   // Least Frequently Used
)
```

### Paramètres Détaillés

#### Cache de Hash (`HashCache*`)

- **Objectif**: Accélère le calcul des hash pour les chaînes alpha
- **HashCacheEnabled**: Active/désactive le cache (défaut: `true`)
- **HashCacheMaxSize**: Nombre max d'entrées (défaut: 10 000)
- **HashCacheEviction**: Politique d'éviction (défaut: `LRU`)
- **HashCacheTTL**: Durée de vie des entrées (0 = illimité)

**Cas d'usage**:
- Activer pour des bases de règles avec beaucoup de conditions similaires
- Augmenter la taille pour des bases de règles très grandes (>1000 règles)

#### Cache de Connexion (`ConnectionCache*`)

- **Objectif**: Mémorise les relations parent-enfant entre nœuds
- **ConnectionCacheEnabled**: Active/désactive (défaut: `true`)
- **ConnectionCacheMaxSize**: Nombre max de connexions (défaut: 50 000)
- **ConnectionCacheEviction**: Politique (défaut: `LRU`)
- **ConnectionCacheTTL**: TTL des connexions

**Cas d'usage**:
- Essentiel pour les réseaux avec beaucoup de nœuds partagés
- Augmenter pour des bases de règles avec >500 règles

#### Métriques (`Metrics*`)

- **MetricsEnabled**: Active la collecte de métriques (défaut: `true`)
- **MetricsDetailedChains**: Stocke les détails de chaque chaîne (défaut: `true`)
- **MetricsMaxChainDetails**: Limite de chaînes trackées (défaut: 1 000)

**Impact**:
- Désactiver `MetricsDetailedChains` réduit l'utilisation mémoire de ~20%
- Utile pour le debugging et l'optimisation

#### Performance Avancée

- **ParallelHashComputation**: Calcul parallèle des hash (défaut: `false`)
  - ⚠️ Expérimental - peut augmenter la performance de 15-30%
  - Recommandé uniquement pour >1000 règles

#### Prometheus

- **PrometheusEnabled**: Exporte métriques vers Prometheus (défaut: `false`)
- **PrometheusPrefix**: Préfixe pour les métriques (défaut: `"tsd_rete"`)

**Métriques exportées**:
- Nombre de nœuds partagés
- Taux de réutilisation des nœuds
- Performance des caches
- Latence des opérations

#### Beta Cache (`BetaCache*`, `BetaJoinResultCache*`)

- **BetaCacheEnabled**: Active le cache pour les jointures (défaut: `true`)
- **BetaHashCacheMaxSize**: Taille du cache de hash de jointure (défaut: 10 000)
- **BetaJoinResultCacheEnabled**: Cache les résultats de jointure (défaut: `true`)
- **BetaJoinResultCacheMaxSize**: Nombre de résultats cachés (défaut: 5 000)
- **BetaJoinResultCacheTTL**: TTL des résultats (défaut: 1 minute)

**Impact**:
- Réduit de 30-50% le temps de jointure pour des patterns répétitifs
- Le cache de résultats est particulièrement efficace pour les règles avec agrégations

---

## Configuration des Transactions

La structure `TransactionOptions` contrôle le comportement des transactions.

### Structure Complète

```go
type TransactionOptions struct {
    // Timeout maximum pour la soumission d'un batch de faits
    SubmissionTimeout time.Duration

    // Délai entre les tentatives de vérification
    VerifyRetryDelay time.Duration

    // Nombre max de tentatives de vérification
    MaxVerifyRetries int

    // Vérifier tous les faits au commit
    VerifyOnCommit bool
}
```

### Valeurs par Défaut

```go
DefaultTransactionOptions() = &TransactionOptions{
    SubmissionTimeout: 30 * time.Second,
    VerifyRetryDelay:  50 * time.Millisecond,
    MaxVerifyRetries:  10,
    VerifyOnCommit:    true,
}
```

### Garanties de Cohérence

TSD fournit une cohérence forte (**Strong Consistency**) pour toutes les opérations:

- ✅ **Vérification synchrone** de tous les faits
- ✅ **Mécanisme de retry** avec backoff exponentiel
- ✅ **Cohérence read-after-write** garantie
- ✅ **Transactions atomiques** (tout ou rien)
- ✅ **Pas de perte de données** en cas d'échec du stockage

### Caractéristiques de Performance

**Stockage en mémoire**:
- ~10 000-50 000 faits/sec (single-node)
- ~1-10ms de latence moyenne par transaction
- Scalabilité linéaire avec la mémoire disponible

**Future: Réplication réseau**:
- Consensus via algorithme Raft (planifié)
- Tolérance aux pannes Byzantine (en cours de recherche)

### Configuration Réseau

```go
type NetworkCoherenceConfig struct {
    DefaultOptions TransactionOptions
    EnableMetrics  bool  // Collecte de métriques (défaut: true)
}
```

---

## Configuration du Beta Sharing

Le Beta Sharing élimine les JoinNodes dupliqués, réduisant la mémoire de 30-50% et améliorant la performance de 20-40%.

### Structure de Configuration

```go
type BetaSharingConfig struct {
    // Active/désactive le partage de JoinNodes
    Enabled bool

    // Taille du cache LRU pour les hash
    HashCacheSize int

    // Nombre max de nœuds partagés (0 = illimité)
    MaxSharedNodes int

    // Active la collecte de métriques
    EnableMetrics bool

    // Normalise l'ordre des variables pour améliorer le partage
    NormalizeOrder bool

    // Active la normalisation avancée (expérimental)
    EnableAdvancedNormalization bool
}
```

### Configuration par Défaut

```go
BetaSharingConfig{
    Enabled:                     true,
    HashCacheSize:               10000,
    MaxSharedNodes:              10000,
    EnableMetrics:               true,
    NormalizeOrder:              true,
    EnableAdvancedNormalization: false,
}
```

### Comportement

1. **Hash Computation**: Calcule un hash canonique pour chaque signature de jointure
2. **Lookup**: Cherche un nœud existant avec le même hash
3. **Reuse or Create**: Réutilise le nœud ou en crée un nouveau
4. **Reference Counting**: Incrémente/décrémente le compteur de références
5. **Cleanup**: Supprime les nœuds non référencés

### Métriques Disponibles

```go
type BetaSharingStats struct {
    TotalCreated     int64  // Nœuds créés
    TotalReused      int64  // Nœuds réutilisés
    CurrentShared    int    // Nœuds actuellement partagés
    SharingRatio     float64 // Ratio de partage (0-100%)
    MemorySavedBytes int64  // Mémoire économisée estimée
}
```

---

## Configuration du Cache Arithmétique

Le cache arithmétique stocke les résultats intermédiaires des expressions arithmétiques.

### Structure de Configuration

```go
type ArithmeticCacheConfig struct {
    Enabled     bool
    MaxSize     int
    Eviction    CacheEvictionPolicy
    TTL         time.Duration
    EnableStats bool
}
```

### Configuration par Défaut

```go
DefaultCacheConfig() = &ArithmeticCacheConfig{
    Enabled:     true,
    MaxSize:     5000,
    Eviction:    EvictionPolicyLRU,
    TTL:         5 * time.Minute,
    EnableStats: true,
}
```

### Cas d'Usage

- **Actions arithmétiques complexes**: `salary * 1.1 + bonus - tax`
- **Expressions répétitives**: Calculs appliqués à plusieurs faits
- **Agrégations**: SUM, AVG, COUNT avec expressions

**Impact**:
- Réduit de 40-60% le temps de calcul pour des expressions complexes
- Particulièrement efficace pour les règles avec beaucoup d'actions

---

## Configuration du Logger

Le système de logging est configurable pour adapter la verbosité.

### Niveaux de Log

```go
type LogLevel int

const (
    LogLevelSilent LogLevel = iota  // Aucun log
    LogLevelError                    // Erreurs uniquement
    LogLevelWarn                     // Warnings + erreurs
    LogLevelInfo                     // Info + warnings + erreurs (défaut)
    LogLevelDebug                    // Tous les logs incluant debug
)
```

### Configuration

```go
// Créer un logger personnalisé
logger := rete.NewLogger(rete.LogLevelDebug, os.Stdout)
logger.SetTimestamps(true)
logger.SetPrefix("[MyApp]")

// Configurer pour le réseau
network.SetLogger(logger)

// Ou utiliser le logger global
rete.SetGlobalLogLevel(rete.LogLevelWarn)
```

### Options

- **Timestamps**: Active/désactive les timestamps
- **Prefix**: Préfixe personnalisé pour identifier la source
- **Output**: Destination des logs (`os.Stdout`, fichier, etc.)
- **Context**: Créer des sous-loggers avec contexte

---

## Profils de Configuration Prédéfinis

TSD fournit 4 profils optimisés pour différents cas d'usage.

### 1. Configuration par Défaut (`DefaultChainPerformanceConfig()`)

**Cas d'usage**: Usage général, équilibre performance/mémoire

```go
config := rete.DefaultChainPerformanceConfig()

// Caractéristiques:
// - Hash Cache: 10k entrées, LRU
// - Connection Cache: 50k entrées, LRU
// - Métriques détaillées: activées
// - Beta Cache: 10k + 5k résultats
// - Prometheus: désactivé
```

**Performance**:
- Mémoire: ~60-80 MB pour 500 règles
- Latence: <5ms par opération
- Throughput: ~20 000 faits/sec

### 2. Configuration Haute Performance (`HighPerformanceConfig()`)

**Cas d'usage**: Systèmes avec beaucoup de mémoire, performance maximale

```go
config := rete.HighPerformanceConfig()

// Caractéristiques:
// - Hash Cache: 100k entrées
// - Connection Cache: 200k entrées
// - Métriques détaillées: désactivées (économie mémoire)
// - Calcul parallèle: activé
// - Beta Cache: 100k + 50k résultats
// - Prometheus: activé
```

**Performance**:
- Mémoire: ~500-800 MB pour 500 règles
- Latence: <2ms par opération
- Throughput: ~50 000 faits/sec

**⚠️ Recommandations**:
- Minimum 8 GB RAM
- Multi-core CPU (4+ cores)
- >1000 règles pour bénéficier pleinement

### 3. Configuration Mémoire Réduite (`LowMemoryConfig()`)

**Cas d'usage**: Systèmes embarqués, containers avec limite de RAM

```go
config := rete.LowMemoryConfig()

// Caractéristiques:
// - Hash Cache: 1k entrées, LRU, TTL 5min
// - Connection Cache: 5k entrées, TTL 5min
// - Métriques détaillées: désactivées
// - Beta Cache: 1k entrées seulement
// - Join Result Cache: désactivé
// - Prometheus: désactivé
```

**Performance**:
- Mémoire: ~10-20 MB pour 500 règles
- Latence: ~10-15ms par opération
- Throughput: ~5 000 faits/sec

**Compromis**:
- Performance réduite de ~50-60%
- Empreinte mémoire réduite de ~80%

### 4. Configuration Sans Caches (`DisabledCachesConfig()`)

**Cas d'usage**: Tests, debugging, validation fonctionnelle

```go
config := rete.DisabledCachesConfig()

// Caractéristiques:
// - Tous les caches: désactivés
// - Métriques détaillées: activées
// - Pas d'optimisations
```

**⚠️ Ne PAS utiliser en production**

---

## Exemples d'Utilisation

### Exemple 1: Configuration de Base

```go
package main

import (
    "github.com/yourusername/tsd/rete"
    "github.com/yourusername/tsd/storage"
)

func main() {
    // Créer un stockage en mémoire
    store := storage.NewMemoryStorage()

    // Créer le réseau avec config par défaut
    network := rete.NewReteNetwork(store)

    // Ajouter des règles...
}
```

### Exemple 2: Configuration Personnalisée

```go
func main() {
    store := storage.NewMemoryStorage()

    // Créer une config personnalisée
    config := rete.DefaultChainPerformanceConfig()
    config.HashCacheMaxSize = 50000           // Cache plus grand
    config.MetricsDetailedChains = false      // Économie mémoire
    config.PrometheusEnabled = true           // Monitoring
    config.BetaJoinResultCacheMaxSize = 10000 // Plus de résultats cachés

    // Valider la config
    if err := config.Validate(); err != nil {
        log.Fatal(err)
    }

    // Créer le réseau avec cette config
    network := rete.NewReteNetworkWithConfig(store, config)
}
```

### Exemple 3: Configuration pour Production

```go
func createProductionNetwork() *rete.ReteNetwork {
    store := storage.NewMemoryStorage()

    // Config haute performance
    config := rete.HighPerformanceConfig()

    // Logger configuré
    logger := rete.NewLogger(rete.LogLevelInfo, os.Stdout)
    logger.SetPrefix("[PROD]")

    // Créer le réseau
    network := rete.NewReteNetworkWithConfig(store, config)
    network.SetLogger(logger)

    // Options de transaction personnalisées
    txOpts := rete.DefaultTransactionOptions()
    txOpts.SubmissionTimeout = 60 * time.Second // Timeout plus long
    txOpts.MaxVerifyRetries = 20                 // Plus de retries

    return network
}
```

### Exemple 4: Configuration pour Tests

```go
func createTestNetwork() *rete.ReteNetwork {
    store := storage.NewMemoryStorage()

    // Config sans caches pour tests déterministes
    config := rete.DisabledCachesConfig()

    // Logger silencieux pour tests
    logger := rete.NewLogger(rete.LogLevelError, os.Stdout)

    network := rete.NewReteNetworkWithConfig(store, config)
    network.SetLogger(logger)

    return network
}
```

### Exemple 5: Monitoring avec Prometheus

```go
func setupMonitoring() *rete.ReteNetwork {
    store := storage.NewMemoryStorage()

    config := rete.HighPerformanceConfig()
    config.PrometheusEnabled = true
    config.PrometheusPrefix = "myapp_rete"

    network := rete.NewReteNetworkWithConfig(store, config)

    // Métriques accessibles sur :9090/metrics
    // Exemples de métriques:
    // - myapp_rete_alpha_nodes_total
    // - myapp_rete_beta_sharing_ratio
    // - myapp_rete_cache_hit_rate
    // - myapp_rete_transaction_latency_ms

    return network
}
```

### Exemple 6: Ajustement Dynamique

```go
func adjustConfigBasedOnLoad(network *rete.ReteNetwork, numRules int) {
    config := network.GetConfig()

    if numRules > 1000 {
        // Beaucoup de règles -> augmenter les caches
        config.HashCacheMaxSize = 100000
        config.ConnectionCacheMaxSize = 200000
        config.ParallelHashComputation = true
    } else if numRules < 100 {
        // Peu de règles -> réduire l'empreinte
        config.HashCacheMaxSize = 5000
        config.ConnectionCacheMaxSize = 10000
        config.MetricsDetailedChains = false
    }

    // Valider les changements
    if err := config.Validate(); err != nil {
        log.Printf("Config invalide: %v", err)
    }
}
```

---

## Métriques et Monitoring

### Métriques de Performance des Chaînes

```go
metrics := network.GetChainMetrics()

// Statistiques disponibles:
// - TotalAlphaChains: nombre de chaînes alpha créées
// - SharedAlphaNodes: nombre de nœuds alpha partagés
// - UniqueAlphaNodes: nombre de nœuds alpha uniques
// - AvgChainLength: longueur moyenne des chaînes
// - MaxChainLength: longueur max d'une chaîne
// - CacheHitRate: taux de hit du cache (0-100%)
```

### Métriques de Beta Sharing

```go
betaStats := network.GetBetaSharingStats()

// Statistiques disponibles:
// - TotalCreated: nœuds créés
// - TotalReused: nœuds réutilisés
// - CurrentShared: nœuds actuellement partagés
// - SharingRatio: ratio de partage (%)
// - MemorySavedBytes: mémoire économisée estimée
```

### Statistiques Réseau

```go
stats := network.GetNetworkStats()

// Informations disponibles:
// - type_nodes: nombre de TypeNodes
// - alpha_nodes: nombre d'AlphaNodes
// - beta_nodes: nombre de JoinNodes
// - terminal_nodes: nombre de TerminalNodes
// - lifecycle_* : stats du gestionnaire de cycle de vie
// - sharing_* : stats du partage alpha
```

---

## Estimation de l'Utilisation Mémoire

```go
config := network.GetConfig()
estimatedBytes := config.EstimateMemoryUsage()

// Formule:
// - Hash Cache: ~500 bytes/entrée
// - Connection Cache: ~100 bytes/entrée
// - Chain Details: ~200 bytes/entrée
// - Beta Hash Cache: ~500 bytes/entrée
// - Beta Join Results: ~1000 bytes/entrée

fmt.Printf("Mémoire estimée pour les caches: %d MB\n", 
    estimatedBytes / (1024 * 1024))
```

---

## Bonnes Pratiques

### 1. Démarrage

✅ **Commencer avec la config par défaut**
```go
network := rete.NewReteNetwork(store)
```

✅ **Monitorer les métriques initialement**
```go
config := rete.DefaultChainPerformanceConfig()
config.PrometheusEnabled = true
```

### 2. Optimisation Progressive

✅ **Identifier les goulots d'étranglement**
- Utiliser les métriques Prometheus
- Analyser les stats de sharing
- Profiler avec pprof

✅ **Ajuster progressivement**
- Augmenter les caches si cache hit rate < 70%
- Activer le calcul parallèle si CPU < 50%
- Réduire les caches si mémoire > 80%

### 3. Production

✅ **Configuration robuste**
```go
config := rete.HighPerformanceConfig()
config.PrometheusEnabled = true

txOpts := rete.DefaultTransactionOptions()
txOpts.SubmissionTimeout = 60 * time.Second
txOpts.MaxVerifyRetries = 20
```

✅ **Logging approprié**
```go
logger := rete.NewLogger(rete.LogLevelWarn, logFile)
network.SetLogger(logger)
```

### 4. Tests

✅ **Isolation et déterminisme**
```go
config := rete.DisabledCachesConfig()
logger := rete.NewLogger(rete.LogLevelError, os.Stdout)
```

### 5. Debugging

✅ **Verbosité maximale**
```go
logger := rete.NewLogger(rete.LogLevelDebug, os.Stdout)
config.MetricsDetailedChains = true
```

---

## Troubleshooting

### Problème: Performance dégradée

**Symptômes**: Latence >50ms, throughput <5000 faits/sec

**Solutions**:
1. Vérifier cache hit rate (devrait être >70%)
2. Augmenter `HashCacheMaxSize` et `ConnectionCacheMaxSize`
3. Activer `ParallelHashComputation` si >1000 règles
4. Passer à `HighPerformanceConfig()`

### Problème: Utilisation mémoire excessive

**Symptômes**: Mémoire >1 GB pour <500 règles

**Solutions**:
1. Désactiver `MetricsDetailedChains`
2. Réduire les tailles de cache
3. Activer TTL sur les caches
4. Passer à `LowMemoryConfig()`

### Problème: Transactions qui échouent

**Symptômes**: Erreurs de timeout, verify failures

**Solutions**:
1. Augmenter `SubmissionTimeout`
2. Augmenter `MaxVerifyRetries`
3. Réduire `VerifyRetryDelay` si latence faible
4. Vérifier les logs avec `LogLevelDebug`

### Problème: Beta sharing ne fonctionne pas

**Symptômes**: `SharingRatio` = 0%, trop de JoinNodes créés

**Solutions**:
1. Vérifier que `BetaCacheEnabled = true`
2. Activer `NormalizeOrder = true`
3. Augmenter `HashCacheSize`
4. Vérifier les logs pour les erreurs de hash computation

---

## Références

- [Architecture RETE](./ARCHITECTURE.md)
- [Guide de Performance](./PERFORMANCE.md)
- [API Reference](./API.md)
- [Scripts de Review](../scripts/review-rete/)

---

**Dernière mise à jour**: 2025-01-XX  
**Version**: 1.0.0