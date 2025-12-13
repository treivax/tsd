# Configuration Globale TSD

Guide complet de configuration du système TSD (Type System Development).

## Table des Matières

- [Vue d'Ensemble](#vue-densemble)
- [Architecture de Configuration](#architecture-de-configuration)
- [Configurations par Composant](#configurations-par-composant)
- [Profils de Déploiement](#profils-de-déploiement)
- [Variables d'Environnement](#variables-denvironnement)
- [Fichiers de Configuration](#fichiers-de-configuration)
- [Exemples Pratiques](#exemples-pratiques)

---

## Vue d'Ensemble

TSD offre un système de configuration modulaire permettant d'adapter le moteur de règles à différents contextes :

- **Développement** : Debugging détaillé, caches désactivés
- **Tests** : Configuration déterministe, métriques activées
- **Production** : Performance optimale, monitoring complet
- **Embarqué** : Empreinte mémoire minimale

### Philosophie de Configuration

```
Configuration = Defaults + Environment + Explicit
```

1. **Defaults** : Valeurs optimales pour usage général
2. **Environment** : Variables d'environnement (12-factor app)
3. **Explicit** : Configuration programmatique ou fichiers

---

## Architecture de Configuration

### Hiérarchie des Composants

```
TSD System
├── Network (Réseau RETE)
│   ├── Performance (Caches, optimisations)
│   ├── Coherence (Transactions, cohérence)
│   └── Sharing (Alpha/Beta sharing)
├── Storage (Stockage)
│   ├── Type (memory, future: disk, distributed)
│   └── Timeouts
├── Constraint (Analyseur de contraintes)
│   ├── Parser
│   ├── Validator
│   └── Logger
├── Server (Serveur HTTP/HTTPS)
│   ├── Network (Host, Port, TLS)
│   └── Authentication (None, Key, JWT)
└── Client (Client CLI/API)
    ├── Connection (URL, Timeout)
    └── Output (Format, Verbose)
```

---

## Configurations par Composant

### 1. Réseau RETE

> **Documentation détaillée** : [RETE_CONFIGURATION.md](RETE_CONFIGURATION.md)

#### ChainPerformanceConfig

Configure les caches et optimisations du réseau RETE.

```go
type ChainPerformanceConfig struct {
    // Caches Alpha
    HashCacheEnabled       bool
    HashCacheMaxSize       int
    HashCacheEviction      CacheEvictionPolicy
    HashCacheTTL           time.Duration
    
    ConnectionCacheEnabled bool
    ConnectionCacheMaxSize int
    ConnectionCacheTTL     time.Duration
    
    // Métriques
    MetricsEnabled         bool
    MetricsDetailedChains  bool
    MetricsMaxChainDetails int
    
    // Performance
    ParallelHashComputation bool
    
    // Prometheus
    PrometheusEnabled bool
    PrometheusPrefix  string
    
    // Beta Cache
    BetaCacheEnabled           bool
    BetaHashCacheMaxSize       int
    BetaJoinResultCacheEnabled bool
    BetaJoinResultCacheMaxSize int
    BetaJoinResultCacheTTL     time.Duration
}
```

**Profils prédéfinis** :
- `DefaultChainPerformanceConfig()` - Usage général
- `HighPerformanceConfig()` - Performance maximale
- `LowMemoryConfig()` - Empreinte mémoire réduite
- `DisabledCachesConfig()` - Tests/debugging

**Exemple** :
```go
config := rete.HighPerformanceConfig()
config.PrometheusEnabled = true
network := rete.NewReteNetworkWithConfig(store, config)
```

---

#### TransactionOptions

Configure le comportement des transactions.

```go
type TransactionOptions struct {
    SubmissionTimeout time.Duration  // Timeout soumission
    VerifyRetryDelay  time.Duration  // Délai entre retries
    MaxVerifyRetries  int            // Nombre max de retries
    VerifyOnCommit    bool           // Vérification au commit
}
```

**Valeurs par défaut** :
```go
DefaultTransactionOptions() = &TransactionOptions{
    SubmissionTimeout: 30 * time.Second,
    VerifyRetryDelay:  50 * time.Millisecond,
    MaxVerifyRetries:  10,
    VerifyOnCommit:    true,
}
```

**Garanties** :
- ✅ Cohérence forte (Strong Consistency)
- ✅ Read-after-write garantie
- ✅ Transactions atomiques
- ✅ Pas de perte de données

---

#### BetaSharingConfig

Configure le partage des JoinNodes (économie 30-50% mémoire).

```go
type BetaSharingConfig struct {
    Enabled                     bool
    HashCacheSize               int
    MaxSharedNodes              int
    EnableMetrics               bool
    NormalizeOrder              bool
    EnableAdvancedNormalization bool
}
```

**Exemple** :
```go
config := BetaSharingConfig{
    Enabled:        true,
    HashCacheSize:  10000,
    MaxSharedNodes: 10000,
    EnableMetrics:  true,
    NormalizeOrder: true,
}
```

---

#### Logger

Configure le système de logging.

```go
type LogLevel int

const (
    LogLevelSilent LogLevel = iota
    LogLevelError
    LogLevelWarn
    LogLevelInfo
    LogLevelDebug
)
```

**Exemple** :
```go
logger := rete.NewLogger(rete.LogLevelInfo, os.Stdout)
logger.SetPrefix("[PROD]")
logger.SetTimestamps(true)
network.SetLogger(logger)
```

---

### 2. Storage (Stockage)

#### MemoryStorage

Configuration du stockage en mémoire (actuel).

```go
type StorageConfig struct {
    Type    string        // "memory" (seul type supporté)
    Timeout time.Duration // Timeout opérations
}
```

**Défaut** :
```go
StorageConfig{
    Type:    "memory",
    Timeout: 5 * time.Second,
}
```

**Future** : Disk storage, distributed storage via Raft.

---

### 3. Constraint (Analyseur)

#### ParserConfig

Configure l'analyseur de contraintes.

```go
type ParserConfig struct {
    MaxExpressions int  // Nombre max d'expressions
    Debug          bool // Mode debug
    Recover        bool // Récupération sur erreur
}
```

**Défaut** :
```go
ParserConfig{
    MaxExpressions: 10000,
    Debug:          false,
    Recover:        true,
}
```

---

#### ValidatorConfig

Configure le validateur de contraintes.

```go
type ValidatorConfig struct {
    StrictMode       bool     // Mode strict
    AllowedOperators []string // Opérateurs autorisés
    MaxDepth         int      // Profondeur max expressions
}
```

**Défaut** :
```go
ValidatorConfig{
    StrictMode:       false,
    AllowedOperators: []string{"==", "!=", "<", ">", "<=", ">="},
    MaxDepth:         100,
}
```

---

### 4. Server (Serveur HTTP/HTTPS)

#### ServerConfig

Configure le serveur HTTP/HTTPS.

```go
type Config struct {
    // Network
    Host string // Adresse d'écoute
    Port int    // Port d'écoute
    
    // TLS/HTTPS
    TLSCertFile string // Certificat TLS
    TLSKeyFile  string // Clé privée TLS
    
    // Authentication
    AuthType      string        // "none", "key", "jwt"
    AuthKeys      []string      // Clés API (mode "key")
    JWTSecret     string        // Secret JWT (mode "jwt")
    JWTExpiration time.Duration // Expiration JWT
    JWTIssuer     string        // Issuer JWT
    
    // Behavior
    Verbose bool // Logs détaillés
}
```

**Exemple HTTP** :
```go
config := &Config{
    Host:     "0.0.0.0",
    Port:     8080,
    AuthType: "key",
    AuthKeys: []string{"my-secret-key"},
}
```

**Exemple HTTPS** :
```go
config := &Config{
    Host:        "0.0.0.0",
    Port:        8443,
    TLSCertFile: "/etc/tsd/certs/server.crt",
    TLSKeyFile:  "/etc/tsd/certs/server.key",
    AuthType:    "jwt",
    JWTSecret:   "your-jwt-secret",
}
```

---

### 5. Client (Client CLI)

#### ClientConfig

Configure le client CLI/API.

```go
type Config struct {
    ServerURL string        // URL du serveur
    Timeout   time.Duration // Timeout requêtes
    AuthToken string        // Token d'authentification
    
    // Input
    File     string // Fichier d'entrée
    Text     string // Texte direct
    UseStdin bool   // Lire depuis stdin
    
    // Output
    Format  string // json, text, yaml
    Verbose bool   // Logs détaillés
}
```

**Exemple** :
```go
config := &Config{
    ServerURL: "https://localhost:8443",
    Timeout:   30 * time.Second,
    AuthToken: "Bearer your-jwt-token",
    Format:    "json",
    Verbose:   true,
}
```

---

### 6. Authentication (Authentification)

#### AuthConfig

Configure l'authentification.

```go
type Config struct {
    Type          string        // "none", "key", "jwt"
    AuthKeys      []string      // Clés API
    JWTSecret     string        // Secret JWT
    JWTExpiration time.Duration // Durée validité JWT
    JWTIssuer     string        // Issuer JWT
}
```

**Types d'authentification** :

1. **None** : Aucune authentification
   ```go
   Config{Type: "none"}
   ```

2. **Key** : Clés API statiques
   ```go
   Config{
       Type:     "key",
       AuthKeys: []string{"key1", "key2", "key3"},
   }
   ```

3. **JWT** : JSON Web Tokens
   ```go
   Config{
       Type:          "jwt",
       JWTSecret:     "your-256-bit-secret",
       JWTExpiration: 24 * time.Hour,
       JWTIssuer:     "tsd-server",
   }
   ```

---

## Profils de Déploiement

### Développement

```go
// RETE Network
reteConfig := rete.DefaultChainPerformanceConfig()
reteConfig.MetricsEnabled = true
reteConfig.MetricsDetailedChains = true

// Logger
logger := rete.NewLogger(rete.LogLevelDebug, os.Stdout)

// Server
serverConfig := &ServerConfig{
    Host:     "127.0.0.1",
    Port:     8080,
    AuthType: "none",
    Verbose:  true,
}
```

**Caractéristiques** :
- Logging détaillé
- Métriques complètes
- Pas d'authentification
- Local uniquement

---

### Test

```go
// RETE Network - Déterministe
reteConfig := rete.DisabledCachesConfig()

// Logger - Silencieux
logger := rete.NewLogger(rete.LogLevelError, os.Stdout)

// Transactions - Strictes
txOpts := rete.DefaultTransactionOptions()
txOpts.MaxVerifyRetries = 3
txOpts.VerifyOnCommit = true
```

**Caractéristiques** :
- Caches désactivés (déterminisme)
- Logging minimal
- Vérifications strictes

---

### Production

```go
// RETE Network - Haute Performance
reteConfig := rete.HighPerformanceConfig()
reteConfig.PrometheusEnabled = true
reteConfig.PrometheusPrefix = "prod_tsd"

// Logger - Optimisé
logger := rete.NewLogger(rete.LogLevelWarn, logFile)
logger.SetPrefix("[PROD]")

// Server - HTTPS + JWT
serverConfig := &ServerConfig{
    Host:          "0.0.0.0",
    Port:          8443,
    TLSCertFile:   "/etc/tsd/certs/server.crt",
    TLSKeyFile:    "/etc/tsd/certs/server.key",
    AuthType:      "jwt",
    JWTSecret:     os.Getenv("TSD_JWT_SECRET"),
    JWTExpiration: 1 * time.Hour,
    Verbose:       false,
}

// Transactions - Robustes
txOpts := rete.DefaultTransactionOptions()
txOpts.SubmissionTimeout = 60 * time.Second
txOpts.MaxVerifyRetries = 20
```

**Caractéristiques** :
- Performance maximale
- HTTPS + JWT
- Monitoring Prometheus
- Logging production
- Transactions robustes

---

### Embarqué / IoT

```go
// RETE Network - Mémoire minimale
reteConfig := rete.LowMemoryConfig()
reteConfig.PrometheusEnabled = false
reteConfig.MetricsDetailedChains = false

// Logger - Minimal
logger := rete.NewLogger(rete.LogLevelError, os.Stdout)

// Pas de serveur - Utilisation programmatique uniquement
```

**Caractéristiques** :
- Empreinte mémoire ~10-20 MB
- Caches réduits
- Pas de métriques détaillées
- Logging minimal

---

## Variables d'Environnement

TSD supporte la configuration via variables d'environnement (12-factor app).

### Serveur

```bash
# Network
TSD_HOST=0.0.0.0
TSD_PORT=8443

# TLS
TSD_TLS_CERT=/etc/tsd/certs/server.crt
TSD_TLS_KEY=/etc/tsd/certs/server.key

# Authentication
TSD_AUTH_TYPE=jwt
TSD_JWT_SECRET=your-256-bit-secret
TSD_JWT_EXPIRATION=1h
TSD_JWT_ISSUER=tsd-production

# Behavior
TSD_VERBOSE=false
TSD_LOG_LEVEL=warn
```

### Client

```bash
# Connection
TSD_SERVER_URL=https://api.example.com:8443
TSD_TIMEOUT=30s
TSD_AUTH_TOKEN=your-jwt-token

# Output
TSD_FORMAT=json
TSD_VERBOSE=true
```

### RETE Network

```bash
# Performance
TSD_RETE_HASH_CACHE_SIZE=100000
TSD_RETE_CONNECTION_CACHE_SIZE=200000
TSD_RETE_PARALLEL_HASH=true

# Prometheus
TSD_PROMETHEUS_ENABLED=true
TSD_PROMETHEUS_PREFIX=prod_tsd

# Beta Sharing
TSD_BETA_SHARING_ENABLED=true
TSD_BETA_CACHE_SIZE=10000
```

---

## Fichiers de Configuration

### Format JSON

**Emplacement** : `config.json`, `/etc/tsd/config.json`, `~/.tsd/config.json`

```json
{
  "server": {
    "host": "0.0.0.0",
    "port": 8443,
    "tls_cert": "/etc/tsd/certs/server.crt",
    "tls_key": "/etc/tsd/certs/server.key"
  },
  "auth": {
    "type": "jwt",
    "jwt_secret": "your-secret",
    "jwt_expiration": "1h",
    "jwt_issuer": "tsd-prod"
  },
  "rete": {
    "hash_cache_size": 100000,
    "connection_cache_size": 200000,
    "metrics_enabled": true,
    "prometheus_enabled": true
  },
  "storage": {
    "type": "memory",
    "timeout": "5s"
  },
  "logger": {
    "level": "warn",
    "format": "json",
    "output": "/var/log/tsd/server.log"
  }
}
```

### Format YAML

```yaml
server:
  host: 0.0.0.0
  port: 8443
  tls:
    cert: /etc/tsd/certs/server.crt
    key: /etc/tsd/certs/server.key

auth:
  type: jwt
  jwt:
    secret: your-secret
    expiration: 1h
    issuer: tsd-prod

rete:
  performance:
    hash_cache_size: 100000
    connection_cache_size: 200000
    parallel_hash: true
  metrics:
    enabled: true
    detailed_chains: false
  prometheus:
    enabled: true
    prefix: prod_tsd

logger:
  level: warn
  format: json
  output: /var/log/tsd/server.log
```

---

## Exemples Pratiques

### Exemple 1 : Serveur de Développement

```go
package main

import (
    "github.com/yourusername/tsd/rete"
    "github.com/yourusername/tsd/internal/servercmd"
    "github.com/yourusername/tsd/storage"
)

func main() {
    // Storage
    store := storage.NewMemoryStorage()
    
    // RETE Network
    reteConfig := rete.DefaultChainPerformanceConfig()
    network := rete.NewReteNetworkWithConfig(store, reteConfig)
    
    // Logger
    logger := rete.NewLogger(rete.LogLevelDebug, os.Stdout)
    network.SetLogger(logger)
    
    // Server
    serverConfig := &servercmd.Config{
        Host:     "127.0.0.1",
        Port:     8080,
        AuthType: "none",
        Verbose:  true,
    }
    
    // Démarrer
    servercmd.Run(serverConfig, network)
}
```

---

### Exemple 2 : Production avec Docker

**Dockerfile** :
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o tsd-server ./cmd/tsd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/tsd-server .
COPY --from=builder /app/config/prod.json /etc/tsd/config.json
EXPOSE 8443
CMD ["./tsd-server", "server", "--config", "/etc/tsd/config.json"]
```

**docker-compose.yml** :
```yaml
version: '3.8'

services:
  tsd-server:
    build: .
    ports:
      - "8443:8443"
    environment:
      - TSD_JWT_SECRET=${TSD_JWT_SECRET}
      - TSD_LOG_LEVEL=warn
      - TSD_PROMETHEUS_ENABLED=true
    volumes:
      - ./certs:/etc/tsd/certs:ro
      - ./logs:/var/log/tsd
    restart: unless-stopped
```

**.env** :
```bash
TSD_JWT_SECRET=your-production-secret-256-bits
```

---

### Exemple 3 : Client avec Configuration

```go
package main

import (
    "github.com/yourusername/tsd/internal/clientcmd"
)

func main() {
    config := &clientcmd.Config{
        ServerURL: "https://api.example.com:8443",
        AuthToken: os.Getenv("TSD_AUTH_TOKEN"),
        File:      "rules.tsd",
        Format:    "json",
        Timeout:   30 * time.Second,
        Verbose:   true,
    }
    
    clientcmd.Run(config)
}
```

---

### Exemple 4 : Configuration Complète Programmatique

```go
package main

import (
    "github.com/yourusername/tsd/rete"
    "github.com/yourusername/tsd/storage"
    "github.com/yourusername/tsd/auth"
)

func createProductionSystem() *rete.ReteNetwork {
    // Storage
    store := storage.NewMemoryStorage()
    
    // RETE Performance Config
    perfConfig := rete.HighPerformanceConfig()
    perfConfig.PrometheusEnabled = true
    perfConfig.PrometheusPrefix = "prod_tsd"
    perfConfig.HashCacheMaxSize = 100000
    perfConfig.ConnectionCacheMaxSize = 200000
    
    // Create Network
    network := rete.NewReteNetworkWithConfig(store, perfConfig)
    
    // Logger
    logFile, _ := os.OpenFile("/var/log/tsd/server.log", 
        os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    logger := rete.NewLogger(rete.LogLevelWarn, logFile)
    logger.SetPrefix("[PROD]")
    logger.SetTimestamps(true)
    network.SetLogger(logger)
    
    // Transaction Options
    txOpts := rete.DefaultTransactionOptions()
    txOpts.SubmissionTimeout = 60 * time.Second
    txOpts.MaxVerifyRetries = 20
    
    // Beta Sharing
    betaConfig := rete.BetaSharingConfig{
        Enabled:        true,
        HashCacheSize:  50000,
        MaxSharedNodes: 50000,
        EnableMetrics:  true,
        NormalizeOrder: true,
    }
    // (Beta sharing est activé par défaut dans le network)
    
    return network
}
```

---

## Validation de Configuration

### Validation Programmatique

```go
// Valider RETE config
config := rete.HighPerformanceConfig()
if err := config.Validate(); err != nil {
    log.Fatalf("Configuration invalide: %v", err)
}

// Estimer mémoire
estimatedMB := config.EstimateMemoryUsage() / (1024 * 1024)
log.Printf("Mémoire estimée: %d MB", estimatedMB)
```

### Vérifications Recommandées

```go
// 1. Vérifier cohérence caches
if config.HashCacheEnabled && config.HashCacheMaxSize == 0 {
    log.Fatal("HashCache activé mais taille = 0")
}

// 2. Vérifier limites raisonnables
if config.HashCacheMaxSize > 1000000 {
    log.Warn("HashCache très grand (>1M), risque mémoire")
}

// 3. Vérifier Prometheus prefix
if config.PrometheusEnabled && config.PrometheusPrefix == "" {
    log.Fatal("PrometheusPrefix requis si activé")
}
```

---

## Monitoring et Métriques

### Exposition Prometheus

Lorsque `PrometheusEnabled = true`, métriques exposées sur `/metrics` :

```
# RETE Network
tsd_rete_alpha_nodes_total
tsd_rete_beta_nodes_total
tsd_rete_terminal_nodes_total

# Beta Sharing
tsd_rete_beta_sharing_ratio
tsd_rete_beta_nodes_shared
tsd_rete_beta_nodes_unique

# Caches
tsd_rete_hash_cache_hit_rate
tsd_rete_hash_cache_size
tsd_rete_connection_cache_hit_rate

# Performance
tsd_rete_transaction_latency_ms
tsd_rete_facts_processed_total
tsd_rete_rules_executed_total
```

### Grafana Dashboard

Exemple de requêtes :

```promql
# Taux de hit du cache
rate(tsd_rete_hash_cache_hits[5m]) / 
rate(tsd_rete_hash_cache_requests[5m])

# Latence transactions (p95)
histogram_quantile(0.95, 
  rate(tsd_rete_transaction_latency_ms_bucket[5m]))

# Ratio de partage Beta
tsd_rete_beta_sharing_ratio
```

---

## Troubleshooting

### Performance Dégradée

**Symptôme** : Latence >50ms

**Solutions** :
1. Vérifier cache hit rate (>70%)
2. Augmenter caches si hit rate faible
3. Activer `ParallelHashComputation`
4. Passer à `HighPerformanceConfig()`

### Mémoire Excessive

**Symptôme** : >1 GB pour <500 règles

**Solutions** :
1. Désactiver `MetricsDetailedChains`
2. Réduire tailles de cache
3. Activer TTL sur caches
4. Passer à `LowMemoryConfig()`

### Erreurs d'Authentification

**Symptôme** : HTTP 401/403

**Solutions** :
1. Vérifier `AuthType` correspond
2. Vérifier token/key valide
3. Vérifier `JWTSecret` identique client/serveur
4. Vérifier expiration JWT

---

## Références

- [RETE Configuration](RETE_CONFIGURATION.md) - Configuration détaillée RETE
- [Architecture](../ARCHITECTURE.md) - Architecture système
- [API Reference](../api/PUBLIC_API.md) - Documentation API
- [User Guide](../USER_GUIDE.md) - Guide utilisateur
- [Logging Guide](../LOGGING_GUIDE.md) - Configuration logging

---

**Version** : 1.0.0  
**Dernière mise à jour** : 2025-01-XX  
**Mainteneur** : TSD Contributors