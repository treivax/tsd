# API Publique TSD

Documentation de l'API publique du moteur de règles TSD.

## Table des Matières

- [Vue d'Ensemble](#vue-densemble)
- [API Programmatique Go](#api-programmatique-go)
- [API HTTP/REST](#api-httprest)
- [Exemples d'Utilisation](#exemples-dutilisation)

---

## Vue d'Ensemble

TSD expose deux types d'API :

1. **API Programmatique Go** : Utilisez TSD comme bibliothèque dans votre application Go
2. **API HTTP/REST** : Utilisez TSD comme service via HTTP/HTTPS

---

## API Programmatique Go

### Package Principal : `rete`

#### Création d'un Réseau RETE

```go
import (
    "github.com/yourusername/tsd/rete"
    "github.com/yourusername/tsd/storage"
)

// Créer avec configuration par défaut
store := storage.NewMemoryStorage()
network := rete.NewReteNetwork(store)

// Créer avec configuration personnalisée
config := rete.HighPerformanceConfig()
network := rete.NewReteNetworkWithConfig(store, config)
```

#### Configurations Disponibles

```go
// Profils prédéfinis
config := rete.DefaultChainPerformanceConfig()  // Usage général
config := rete.HighPerformanceConfig()          // Performance max
config := rete.LowMemoryConfig()                // Mémoire minimale
config := rete.DisabledCachesConfig()           // Tests/debug

// Validation
if err := config.Validate(); err != nil {
    log.Fatal(err)
}

// Estimation mémoire
bytes := config.EstimateMemoryUsage()
```

#### Configuration du Logger

```go
// Créer un logger
logger := rete.NewLogger(rete.LogLevelInfo, os.Stdout)

// Configurer
logger.SetLevel(rete.LogLevelDebug)
logger.SetTimestamps(true)
logger.SetPrefix("[MyApp]")

// Attacher au réseau
network.SetLogger(logger)

// Logger avec contexte
contextLogger := logger.WithContext("RuleEngine")
```

**Niveaux de log** :
- `LogLevelSilent` - Aucun log
- `LogLevelError` - Erreurs uniquement
- `LogLevelWarn` - Warnings + erreurs
- `LogLevelInfo` - Info + warnings + erreurs (défaut)
- `LogLevelDebug` - Tous les logs

#### Métriques et Statistiques

```go
// Métriques de performance des chaînes
metrics := network.GetChainMetrics()
fmt.Printf("Chaînes alpha créées: %d\n", metrics.TotalAlphaChains)
fmt.Printf("Nœuds alpha partagés: %d\n", metrics.SharedAlphaNodes)

// Statistiques de partage Beta
betaStats := network.GetBetaSharingStats()
fmt.Printf("Ratio de partage: %.2f%%\n", betaStats.SharingRatio)
fmt.Printf("Mémoire économisée: %d MB\n", betaStats.MemorySavedBytes/(1024*1024))

// Statistiques réseau
stats := network.GetNetworkStats()
fmt.Printf("Type nodes: %d\n", stats["type_nodes"])
fmt.Printf("Alpha nodes: %d\n", stats["alpha_nodes"])
fmt.Printf("Beta nodes: %d\n", stats["beta_nodes"])

// Réinitialiser métriques
network.ResetChainMetrics()
```

#### Transactions

```go
// Options de transaction par défaut
txOpts := rete.DefaultTransactionOptions()

// Configuration personnalisée
txOpts := &rete.TransactionOptions{
    SubmissionTimeout: 60 * time.Second,
    VerifyRetryDelay:  100 * time.Millisecond,
    MaxVerifyRetries:  20,
    VerifyOnCommit:    true,
}

// Validation
if err := txOpts.Validate(); err != nil {
    log.Fatal(err)
}

// Cloner options
clone := txOpts.Clone()
```

#### Génération d'IDs de Faits

```go
// Génération thread-safe d'ID unique
factID := network.GenerateFactID("Person")
// Résultat: "Person_1", "Person_2", etc.
```

---

### Package : `storage`

#### Memory Storage

```go
import "github.com/yourusername/tsd/storage"

// Créer un stockage en mémoire
store := storage.NewMemoryStorage()

// Opérations (implémentation de l'interface Storage)
// - AddFact(fact *Fact) error
// - GetFact(id string) (*Fact, error)
// - RemoveFact(id string) error
// - GetFactsByType(typeName string) ([]*Fact, error)
// - GetAllFacts() ([]*Fact, error)
```

---

### Package : `constraint`

#### Parser

```go
import "github.com/yourusername/tsd/constraint"

// Configuration parser
config := constraint.ParserConfig{
    MaxExpressions: 10000,
    Debug:          false,
    Recover:        true,
}

// Parser de contraintes
parser := constraint.NewParser(config)

// Parser fichier .tsd
result, err := parser.ParseFile("rules.tsd")
if err != nil {
    log.Fatal(err)
}

// Accéder aux résultats
for _, typeDef := range result.Types {
    fmt.Printf("Type: %s\n", typeDef.Name)
}

for _, rule := range result.Rules {
    fmt.Printf("Rule: %s\n", rule.ID)
}
```

#### Validator

```go
// Configuration validateur
validatorConfig := constraint.ValidatorConfig{
    StrictMode:       false,
    AllowedOperators: []string{"==", "!=", "<", ">"},
    MaxDepth:         100,
}

validator := constraint.NewValidator(validatorConfig)

// Valider règles
errors := validator.Validate(result)
if len(errors) > 0 {
    for _, err := range errors {
        fmt.Printf("Erreur: %v\n", err)
    }
}
```

---

### Package : `auth`

#### Configuration Authentification

```go
import "github.com/yourusername/tsd/auth"

// Aucune authentification
config := &auth.Config{Type: "none"}

// Authentification par clé API
config := &auth.Config{
    Type:     "key",
    AuthKeys: []string{"key1", "key2"},
}

// Authentification JWT
config := &auth.Config{
    Type:          "jwt",
    JWTSecret:     "your-secret",
    JWTExpiration: 24 * time.Hour,
    JWTIssuer:     "tsd-server",
}

// Créer authenticator
authenticator := auth.NewAuthenticator(config)

// Valider token/key
valid, err := authenticator.Validate(token)
```

---

## API HTTP/REST

### Endpoints Disponibles

#### POST /compile

Compile et exécute un fichier TSD.

**Request** :
```http
POST /compile HTTP/1.1
Content-Type: text/plain
Authorization: Bearer <token>

# Contenu du fichier .tsd
type Person
Person("Alice", 30)
{p: Person} / p.age > 25 ==> print("Adult: ", p.name)
```

**Response** :
```json
{
  "status": "success",
  "results": [
    {"type": "print", "message": "Adult: Alice"}
  ],
  "execution_time_ms": 12
}
```

#### GET /health

Vérifier l'état du serveur.

**Request** :
```http
GET /health HTTP/1.1
```

**Response** :
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime_seconds": 3600
}
```

#### GET /metrics

Métriques Prometheus (si activé).

**Request** :
```http
GET /metrics HTTP/1.1
```

**Response** :
```
# HELP tsd_rete_alpha_nodes_total Total number of alpha nodes
# TYPE tsd_rete_alpha_nodes_total gauge
tsd_rete_alpha_nodes_total 150

# HELP tsd_rete_beta_sharing_ratio Beta node sharing ratio
# TYPE tsd_rete_beta_sharing_ratio gauge
tsd_rete_beta_sharing_ratio 0.65
```

### Authentification

#### API Key

```http
POST /compile HTTP/1.1
Authorization: ApiKey your-api-key-here
Content-Type: text/plain

[contenu TSD]
```

#### JWT Token

```http
POST /compile HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: text/plain

[contenu TSD]
```

### Codes de Statut HTTP

| Code | Signification |
|------|---------------|
| 200  | Succès |
| 400  | Requête invalide |
| 401  | Non authentifié |
| 403  | Accès refusé |
| 500  | Erreur serveur |

---

## Exemples d'Utilisation

### Exemple 1 : Application Go Basique

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yourusername/tsd/constraint"
    "github.com/yourusername/tsd/rete"
    "github.com/yourusername/tsd/storage"
)

func main() {
    // 1. Créer le stockage
    store := storage.NewMemoryStorage()
    
    // 2. Créer le réseau RETE
    config := rete.DefaultChainPerformanceConfig()
    network := rete.NewReteNetworkWithConfig(store, config)
    
    // 3. Configurer le logger
    logger := rete.NewLogger(rete.LogLevelInfo, os.Stdout)
    network.SetLogger(logger)
    
    // 4. Parser le fichier TSD
    parser := constraint.NewParser(constraint.ParserConfig{})
    result, err := parser.ParseFile("rules.tsd")
    if err != nil {
        log.Fatal(err)
    }
    
    // 5. Charger dans le réseau
    // (code de chargement ici)
    
    // 6. Exécuter
    // (code d'exécution ici)
    
    // 7. Récupérer métriques
    metrics := network.GetChainMetrics()
    fmt.Printf("Performance: %+v\n", metrics)
}
```

### Exemple 2 : Serveur HTTP avec Configuration Complète

```go
package main

import (
    "log"
    "os"
    "time"
    
    "github.com/yourusername/tsd/auth"
    "github.com/yourusername/tsd/internal/servercmd"
    "github.com/yourusername/tsd/rete"
    "github.com/yourusername/tsd/storage"
)

func main() {
    // Storage
    store := storage.NewMemoryStorage()
    
    // RETE avec monitoring
    reteConfig := rete.HighPerformanceConfig()
    reteConfig.PrometheusEnabled = true
    reteConfig.PrometheusPrefix = "prod_tsd"
    network := rete.NewReteNetworkWithConfig(store, reteConfig)
    
    // Logger production
    logFile, _ := os.OpenFile("tsd.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    logger := rete.NewLogger(rete.LogLevelWarn, logFile)
    network.SetLogger(logger)
    
    // Serveur HTTPS avec JWT
    serverConfig := &servercmd.Config{
        Host:          "0.0.0.0",
        Port:          8443,
        TLSCertFile:   "/etc/tsd/server.crt",
        TLSKeyFile:    "/etc/tsd/server.key",
        AuthType:      "jwt",
        JWTSecret:     os.Getenv("TSD_JWT_SECRET"),
        JWTExpiration: 1 * time.Hour,
        Verbose:       false,
    }
    
    // Démarrer serveur
    log.Println("Starting TSD server on :8443")
    servercmd.Run(serverConfig, network)
}
```

### Exemple 3 : Client HTTP

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

func main() {
    // Lire fichier TSD
    content, _ := ioutil.ReadFile("rules.tsd")
    
    // Créer requête
    req, _ := http.NewRequest("POST", "https://localhost:8443/compile", bytes.NewReader(content))
    req.Header.Set("Content-Type", "text/plain")
    req.Header.Set("Authorization", "Bearer "+os.Getenv("TSD_TOKEN"))
    
    // Envoyer requête
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    
    // Lire réponse
    body, _ := ioutil.ReadAll(resp.Body)
    
    // Parser JSON
    var result map[string]interface{}
    json.Unmarshal(body, &result)
    
    fmt.Printf("Status: %s\n", result["status"])
    fmt.Printf("Results: %v\n", result["results"])
}
```

### Exemple 4 : Utilisation avec Métriques

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/yourusername/tsd/rete"
    "github.com/yourusername/tsd/storage"
)

func main() {
    // Créer réseau avec métriques
    config := rete.DefaultChainPerformanceConfig()
    config.MetricsEnabled = true
    config.MetricsDetailedChains = true
    config.PrometheusEnabled = true
    
    store := storage.NewMemoryStorage()
    network := rete.NewReteNetworkWithConfig(store, config)
    
    // Charger et exécuter règles
    // ...
    
    // Afficher métriques
    displayMetrics(network)
    
    // Monitoring continu
    ticker := time.NewTicker(10 * time.Second)
    for range ticker.C {
        displayMetrics(network)
    }
}

func displayMetrics(network *rete.ReteNetwork) {
    // Métriques chaînes
    metrics := network.GetChainMetrics()
    fmt.Printf("\n=== Chain Metrics ===\n")
    fmt.Printf("Total chains: %d\n", metrics.TotalAlphaChains)
    fmt.Printf("Shared nodes: %d\n", metrics.SharedAlphaNodes)
    
    // Beta sharing
    betaStats := network.GetBetaSharingStats()
    if betaStats != nil {
        fmt.Printf("\n=== Beta Sharing ===\n")
        fmt.Printf("Sharing ratio: %.2f%%\n", betaStats.SharingRatio)
        fmt.Printf("Memory saved: %d KB\n", betaStats.MemorySavedBytes/1024)
    }
    
    // Stats réseau
    stats := network.GetNetworkStats()
    fmt.Printf("\n=== Network Stats ===\n")
    for k, v := range stats {
        fmt.Printf("%s: %v\n", k, v)
    }
}
```

---

## Interfaces Publiques

### Interface Storage

```go
type Storage interface {
    AddFact(fact *Fact) error
    GetFact(id string) (*Fact, error)
    RemoveFact(id string) error
    GetFactsByType(typeName string) ([]*Fact, error)
    GetAllFacts() ([]*Fact, error)
}
```

### Interface Node

```go
type Node interface {
    GetID() string
    GetType() NodeType
    Activate(token *Token) error
}
```

---

## Types Publics Principaux

### ReteNetwork

```go
type ReteNetwork struct {
    // Configuration
    Config *ChainPerformanceConfig
    
    // Méthodes publiques
    GetChainMetrics() *ChainBuildMetrics
    GetBetaSharingStats() *BetaSharingStats
    GetNetworkStats() map[string]interface{}
    SetLogger(logger *Logger)
    GenerateFactID(typeName string) string
    ResetChainMetrics()
}
```

### ChainPerformanceConfig

Voir [RETE Configuration](../configuration/RETE_CONFIGURATION.md) pour les détails complets.

### TransactionOptions

```go
type TransactionOptions struct {
    SubmissionTimeout time.Duration
    VerifyRetryDelay  time.Duration
    MaxVerifyRetries  int
    VerifyOnCommit    bool
}
```

### Logger

```go
type Logger struct {
    // Méthodes publiques
    Debug(format string, args ...interface{})
    Info(format string, args ...interface{})
    Warn(format string, args ...interface{})
    Error(format string, args ...interface{})
    
    SetLevel(level LogLevel)
    SetTimestamps(enabled bool)
    SetPrefix(prefix string)
    WithContext(context string) *Logger
}
```

---

## Bonnes Pratiques

### 1. Gestion des Ressources

```go
// ✅ Bon : Fermer les ressources
store := storage.NewMemoryStorage()
defer store.Close() // si disponible

// ✅ Bon : Logger vers fichier avec rotation
logFile, err := os.OpenFile("tsd.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
    log.Fatal(err)
}
defer logFile.Close()
```

### 2. Configuration

```go
// ✅ Bon : Valider configuration
config := rete.HighPerformanceConfig()
if err := config.Validate(); err != nil {
    log.Fatalf("Invalid config: %v", err)
}

// ✅ Bon : Estimer impact mémoire
bytes := config.EstimateMemoryUsage()
if bytes > 1024*1024*1024 { // > 1GB
    log.Warn("Config uses >1GB memory")
}
```

### 3. Logging

```go
// ✅ Bon : Adapter niveau selon environnement
level := rete.LogLevelWarn
if os.Getenv("ENV") == "development" {
    level = rete.LogLevelDebug
}
logger := rete.NewLogger(level, os.Stdout)
```

### 4. Métriques

```go
// ✅ Bon : Monitoring régulier
ticker := time.NewTicker(1 * time.Minute)
go func() {
    for range ticker.C {
        metrics := network.GetChainMetrics()
        // Envoyer métriques à système monitoring
    }
}()
```

---

## Migration et Compatibilité

### Versions

- **v1.x** : API stable actuelle
- **v2.x** : Futures améliorations (rétrocompatibilité garantie)

### Déprécations

Aucune dépréciation actuellement.

---

## Références

- [Configuration Globale](../configuration/README.md)
- [RETE Configuration](../configuration/RETE_CONFIGURATION.md)
- [User Guide](../USER_GUIDE.md)
- [Architecture](../ARCHITECTURE.md)

---

**Version** : 1.0.0  
**Dernière mise à jour** : 2025-01-XX  
**Mainteneur** : TSD Contributors