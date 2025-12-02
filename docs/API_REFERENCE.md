# API Reference - Pipeline TSD

## Vue d'ensemble

Ce document liste toutes les fonctions publiques disponibles pour utiliser le pipeline TSD avec transactions automatiques obligatoires.

**Version** : 2.0.0 - Transactions Obligatoires  
**Date** : 2025-12-02

---

## 📚 Table des Matières

1. [Fonction Principale](#fonction-principale)
2. [Fonctions avec Métriques](#fonctions-avec-métriques)
3. [Fonctions Avancées](#fonctions-avancées)
4. [Fonctions de Construction](#fonctions-de-construction)
5. [Configuration](#configuration)
6. [Types et Structures](#types-et-structures)

---

## Fonction Principale

### `IngestFile()`

**Fonction recommandée** pour la majorité des cas d'usage.

```go
func (cp *ConstraintPipeline) IngestFile(
    filename string,
    network *ReteNetwork,
    storage Storage,
) (*ReteNetwork, error)
```

**Description** :  
Ingère un fichier de contraintes TSD dans le réseau RETE avec **transaction automatique obligatoire**.

**Caractéristiques** :
- ✅ **Transaction automatique** : Créée, committée ou rollbackée automatiquement
- ✅ **Validation incrémentale** : Utilise le contexte du réseau existant
- ✅ **GC automatique** : Après détection d'une commande `reset`
- ✅ **Propagation de faits** : Les faits existants sont propagés aux nouvelles règles

**Paramètres** :
- `filename` : Chemin vers le fichier `.tsd` à ingérer
- `network` : Réseau RETE existant ou `nil` (un nouveau sera créé)
- `storage` : Interface de stockage des faits

**Retour** :
- `*ReteNetwork` : Le réseau RETE mis à jour
- `error` : Erreur en cas d'échec (rollback automatique effectué)

**Exemple** :
```go
storage := rete.NewMemoryStorage()
pipeline := rete.NewConstraintPipeline()

network, err := pipeline.IngestFile("rules.tsd", nil, storage)
if err != nil {
    // ✅ Rollback automatique déjà effectué
    log.Fatalf("Erreur : %v", err)
}
// ✅ Commit automatique déjà effectué
fmt.Println("Ingestion réussie !")
```

**Complexité** :
- Parsing : O(n) où n = taille du fichier
- Validation : O(m) où m = nombre de types/règles
- Transaction : O(1) pour begin, O(k) pour commit/rollback où k = nombre de commandes

---

## Fonctions avec Métriques

### `IngestFileWithMetrics()`

Pour les cas où vous avez besoin de métriques détaillées sur l'ingestion.

```go
func (cp *ConstraintPipeline) IngestFileWithMetrics(
    filename string,
    network *ReteNetwork,
    storage Storage,
) (*ReteNetwork, *IngestionMetrics, error)
```

**Description** :  
Identique à `IngestFile()` mais retourne également des métriques détaillées.

**Caractéristiques** :
- ✅ **Toutes les fonctionnalités de `IngestFile()`**
- ✅ **Métriques détaillées** : Temps de parsing, validation, construction, etc.
- ✅ **Transaction automatique** : Comme `IngestFile()`

**Retour** :
- `*ReteNetwork` : Le réseau RETE mis à jour
- `*IngestionMetrics` : Métriques détaillées de l'ingestion
- `error` : Erreur en cas d'échec

**Métriques Disponibles** :
```go
type IngestionMetrics struct {
    ParsingDuration      time.Duration
    ValidationDuration   time.Duration
    NetworkBuildDuration time.Duration
    TotalDuration        time.Duration
    FactsPropagated      int
    WasReset             bool
    ValidationSkipped    bool
}
```

**Exemple** :
```go
network, metrics, err := pipeline.IngestFileWithMetrics("rules.tsd", nil, storage)
if err != nil {
    log.Fatalf("Erreur : %v", err)
}

fmt.Printf("Parsing : %v\n", metrics.ParsingDuration)
fmt.Printf("Validation : %v\n", metrics.ValidationDuration)
fmt.Printf("Total : %v\n", metrics.TotalDuration)
```

---

## Fonctions Avancées

### `IngestFileWithAdvancedFeatures()`

Pour les cas nécessitant un contrôle fin de la configuration.

```go
func (cp *ConstraintPipeline) IngestFileWithAdvancedFeatures(
    filename string,
    network *ReteNetwork,
    storage Storage,
    config *AdvancedPipelineConfig,
) (*ReteNetwork, *AdvancedMetrics, error)
```

**Description** :  
Ingestion avec configuration avancée et métriques étendues incluant les transactions.

**Caractéristiques** :
- ✅ **Configuration fine** : Timeout, taille max, auto-commit, etc.
- ✅ **Métriques avancées** : Validation, GC, transactions
- ✅ **Transaction automatique** : Toujours activée (non désactivable)

**Configuration** :
```go
type AdvancedPipelineConfig struct {
    // Transactions (toujours activées)
    TransactionTimeout  time.Duration  // Timeout de la transaction
    MaxTransactionSize  int64          // Taille max de l'empreinte mémoire
    AutoCommit          bool           // Commit automatique
    AutoRollbackOnError bool           // Rollback automatique sur erreur
}
```

**Métriques Avancées** :
```go
type AdvancedMetrics struct {
    // Validation incrémentale
    ValidationWithContextDuration time.Duration
    TypesFoundInContext           int
    ValidationErrors              []string
    
    // Garbage Collection
    GCDuration     time.Duration
    NodesCollected int
    MemoryFreed    int64
    GCPerformed    bool
    
    // Transaction (toujours présente)
    TransactionID        string
    TransactionFootprint int64
    ChangesTracked       int
    RollbackPerformed    bool
    RollbackDuration     time.Duration
    TransactionDuration  time.Duration
}
```

**Exemple** :
```go
config := rete.DefaultAdvancedPipelineConfig()
config.TransactionTimeout = 60 * time.Second
config.MaxTransactionSize = 200 * 1024 * 1024 // 200 MB
config.AutoCommit = true

network, metrics, err := pipeline.IngestFileWithAdvancedFeatures(
    "rules.tsd", nil, storage, config,
)

if err != nil {
    log.Fatalf("Erreur : %v", err)
}

// Afficher les métriques avancées
rete.PrintAdvancedMetrics(metrics)
```

### `IngestFileTransactionalSafe()`

Pour obtenir un accès à la transaction (usage avancé).

```go
func (cp *ConstraintPipeline) IngestFileTransactionalSafe(
    filename string,
    network *ReteNetwork,
    storage Storage,
) (*ReteNetwork, *Transaction, error)
```

**Description** :  
Ingestion avec accès à l'objet transaction pour inspection (sans commit automatique).

**Caractéristiques** :
- ✅ **Transaction accessible** : Retourne l'objet `Transaction`
- ✅ **Pas de commit automatique** : Permet inspection avant commit
- ✅ **Rollback automatique** : En cas d'erreur

**Note** : Cette fonction est pour des cas d'usage très spécifiques. Préférez `IngestFile()` ou `IngestFileWithAdvancedFeatures()`.

**Exemple** :
```go
network, tx, err := pipeline.IngestFileTransactionalSafe("rules.tsd", nil, storage)
if err != nil {
    log.Fatalf("Erreur : %v", err)
}

// Inspecter la transaction
fmt.Printf("Transaction ID : %s\n", tx.ID)
fmt.Printf("Commandes : %d\n", tx.GetCommandCount())

// Commit manuel
tx.Commit()
```

---



## Configuration

### `DefaultAdvancedPipelineConfig()`

Retourne une configuration par défaut pour le pipeline avancé.

```go
func DefaultAdvancedPipelineConfig() *AdvancedPipelineConfig
```

**Valeurs par défaut** :
```go
&AdvancedPipelineConfig{
    TransactionTimeout:  30 * time.Second,
    MaxTransactionSize:  100 * 1024 * 1024, // 100 MB
    AutoCommit:          false,
    AutoRollbackOnError: true,
}
```

### `NewConstraintPipeline()`

Crée une nouvelle instance du pipeline.

```go
func NewConstraintPipeline() *ConstraintPipeline
```

**Exemple** :
```go
pipeline := rete.NewConstraintPipeline()
```

---

## Types et Structures

### `ConstraintPipeline`

Structure principale du pipeline.

```go
type ConstraintPipeline struct {
    // Champs internes (privés)
}
```

### `IngestionMetrics`

Métriques basiques d'ingestion.

```go
type IngestionMetrics struct {
    ParsingDuration      time.Duration
    ValidationDuration   time.Duration
    NetworkBuildDuration time.Duration
    TotalDuration        time.Duration
    FactsPropagated      int
    WasReset             bool
    ValidationSkipped    bool
}
```

### `AdvancedMetrics`

Métriques avancées incluant validation, GC et transactions.

```go
type AdvancedMetrics struct {
    // Validation incrémentale
    ValidationWithContextDuration time.Duration
    TypesFoundInContext           int
    ValidationErrors              []string
    
    // Garbage Collection
    GCDuration     time.Duration
    NodesCollected int
    MemoryFreed    int64
    GCPerformed    bool
    
    // Transaction (toujours présente)
    TransactionID        string
    TransactionFootprint int64
    ChangesTracked       int
    RollbackPerformed    bool
    RollbackDuration     time.Duration
    TransactionDuration  time.Duration
}
```

### `AdvancedPipelineConfig`

Configuration pour le pipeline avancé.

```go
type AdvancedPipelineConfig struct {
    // Transactions (toujours activées)
    TransactionTimeout  time.Duration
    MaxTransactionSize  int64
    AutoCommit          bool
    AutoRollbackOnError bool
}
```

---

## Fonctions Utilitaires

### `PrintAdvancedMetrics()`

Affiche les métriques avancées de manière formatée.

```go
func PrintAdvancedMetrics(metrics *AdvancedMetrics)
```

**Exemple** :
```go
network, metrics, err := pipeline.IngestFileWithAdvancedFeatures(
    "rules.tsd", nil, storage, config,
)

rete.PrintAdvancedMetrics(metrics)
```

**Sortie** :
```
📊 MÉTRIQUES AVANCÉES
═══════════════════════════════════════
🔍 Validation incrémentale
   Durée: 125ms
   Types en contexte: 15

🗑️  Garbage Collection
   Durée: 50ms
   Nœuds collectés: 42

🔒 Transaction
   ID: 550e8400-e29b-41d4-a716-446655440000
   Durée: 250ms
   Empreinte mémoire: 2.34 KB
   Changements trackés: 15
═══════════════════════════════════════
```

### `GetAdvancedMetricsSummary()`

Retourne un résumé textuel des métriques.

```go
func GetAdvancedMetricsSummary(metrics *AdvancedMetrics) string
```

---

## Guide de Sélection

### Quel fonction utiliser ?

| Besoin | Fonction Recommandée |
|--------|---------------------|
| **Cas général** | `IngestFile()` |
| **Besoin de métriques** | `IngestFileWithMetrics()` |
| **Configuration fine** | `IngestFileWithAdvancedFeatures()` |
| **Plusieurs fichiers** | Appels successifs à `IngestFile()` |
| **Accès transaction** | `IngestFileTransactionalSafe()` |

### Matrice de Fonctionnalités

| Fonctionnalité | IngestFile | WithMetrics | WithAdvancedFeatures |
|----------------|-----------|-------------|---------------------|
| Transaction automatique | ✅ | ✅ | ✅ |
| Validation incrémentale | ✅ | ✅ | ✅ |
| GC après reset | ✅ | ✅ | ✅ |
| Métriques basiques | ❌ | ✅ | ✅ |
| Métriques avancées | ❌ | ❌ | ✅ |
| Configuration fine | ❌ | ❌ | ✅ |
| Métriques transaction | ❌ | ❌ | ✅ |

---

## Exemples Complets

### Exemple 1 : Usage Simple

```go
package main

import (
    "fmt"
    "log"
    "github.com/treivax/tsd/rete"
)

func main() {
    storage := rete.NewMemoryStorage()
    pipeline := rete.NewConstraintPipeline()
    
    network, err := pipeline.IngestFile("rules.tsd", nil, storage)
    if err != nil {
        log.Fatalf("Erreur : %v", err)
    }
    
    fmt.Println("Ingestion réussie !")
    fmt.Printf("Types : %d\n", len(network.Types))
    fmt.Printf("Règles : %d\n", len(network.TerminalNodes))
}
```

### Exemple 2 : Avec Métriques

```go
package main

import (
    "fmt"
    "log"
    "github.com/treivax/tsd/rete"
)

func main() {
    storage := rete.NewMemoryStorage()
    pipeline := rete.NewConstraintPipeline()
    
    network, metrics, err := pipeline.IngestFileWithMetrics("rules.tsd", nil, storage)
    if err != nil {
        log.Fatalf("Erreur : %v", err)
    }
    
    fmt.Printf("Parsing : %v\n", metrics.ParsingDuration)
    fmt.Printf("Validation : %v\n", metrics.ValidationDuration)
    fmt.Printf("Total : %v\n", metrics.TotalDuration)
    fmt.Printf("Faits propagés : %d\n", metrics.FactsPropagated)
}
```

### Exemple 3 : Configuration Avancée

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/treivax/tsd/rete"
)

func main() {
    storage := rete.NewMemoryStorage()
    pipeline := rete.NewConstraintPipeline()
    
    // Configuration personnalisée
    config := rete.DefaultAdvancedPipelineConfig()
    config.TransactionTimeout = 60 * time.Second
    config.MaxTransactionSize = 200 * 1024 * 1024
    config.AutoCommit = true
    
    network, metrics, err := pipeline.IngestFileWithAdvancedFeatures(
        "rules.tsd", nil, storage, config,
    )
    
    if err != nil {
        log.Fatalf("Erreur : %v", err)
    }
    
    // Afficher métriques détaillées
    rete.PrintAdvancedMetrics(metrics)
    
    fmt.Printf("\nTransaction ID : %s\n", metrics.TransactionID)
    fmt.Printf("Commandes : %d\n", metrics.ChangesTracked)
}
```

### Exemple 4 : Ingestion Incrémentale

```go
package main

import (
    "fmt"
    "log"
    "github.com/treivax/tsd/rete"
)

func main() {
    storage := rete.NewMemoryStorage()
    pipeline := rete.NewConstraintPipeline()
    
    // Premier fichier (types de base)
    network, err := pipeline.IngestFile("types.tsd", nil, storage)
    if err != nil {
        log.Fatalf("Erreur types : %v", err)
    }
    fmt.Printf("Après types : %d types\n", len(network.Types))
    
    // Deuxième fichier (règles) - ingestion incrémentale
    network, err = pipeline.IngestFile("rules.tsd", network, storage)
    if err != nil {
        log.Fatalf("Erreur règles : %v", err)
    }
    fmt.Printf("Après règles : %d règles\n", len(network.TerminalNodes))
    
    // Troisième fichier (faits)
    network, err = pipeline.IngestFile("facts.tsd", network, storage)
    if err != nil {
        log.Fatalf("Erreur faits : %v", err)
    }
    fmt.Printf("Après faits : %d faits\n", len(network.Storage.GetAllFacts()))
}
```

---

## Notes Importantes

### Transactions

- ✅ **Toujours activées** : Les transactions sont obligatoires et automatiques
- ✅ **Rollback automatique** : En cas d'erreur, rollback systématique
- ✅ **Thread-safe** : Protection par mutex sur tous les accès
- ✅ **Performance** : < 1% d'overhead mémoire

### Validation

- ✅ **Incrémentale** : Utilise le contexte du réseau existant
- ✅ **Automatique** : Toujours effectuée
- ✅ **Optimisée** : Seulement les nouveaux types/règles

### Garbage Collection

- ✅ **Automatique** : Déclenchée après détection d'un `reset`
- ✅ **Complète** : Libère tous les nœuds de l'ancien réseau
- ✅ **Métriques** : Nœuds collectés et mémoire libérée

---

## Fonctions Supprimées

Les fonctions suivantes ont été **SUPPRIMÉES** dans la version 2.0.0 :

### Fonctions de Transaction (supprimées)
- ❌ `IngestFileTransactional()` : Remplacée par `IngestFile()`
- ❌ `IngestFileWithTransaction()` : Remplacée par `IngestFile()`

### Fonctions de Construction (supprimées)
- ❌ `BuildNetworkFromConstraintFile()` : Remplacée par `IngestFile(constraintFile, nil, storage)`
- ❌ `BuildNetworkFromMultipleFiles()` : Remplacée par appels successifs à `IngestFile()`
- ❌ `BuildNetworkFromIterativeParser()` : Remplacée par `IngestFile()`
- ❌ `BuildNetworkFromConstraintFileWithFacts()` : Remplacée par deux appels à `IngestFile()`

**Migration** : Utilisez simplement `IngestFile()` pour tous les cas d'usage.

**Exemple - Plusieurs fichiers** :
```go
// Avant
network, err := pipeline.BuildNetworkFromMultipleFiles([]string{"types.tsd", "rules.tsd"}, storage)

// Après
network, err := pipeline.IngestFile("types.tsd", nil, storage)
if err != nil {
    return err
}
network, err = pipeline.IngestFile("rules.tsd", network, storage)
```

**Exemple - Contraintes + Faits** :
```go
// Avant
network, facts, err := pipeline.BuildNetworkFromConstraintFileWithFacts("rules.tsd", "facts.tsd", storage)

// Après
network, err := pipeline.IngestFile("rules.tsd", nil, storage)
if err != nil {
    return err
}
network, err = pipeline.IngestFile("facts.tsd", network, storage)
facts := storage.GetAllFacts()
```

---

## Références

- [Guide des Transactions Obligatoires](./TRANSACTIONS_MANDATORY.md)
- [Changelog v2.0](./CHANGELOG_TRANSACTIONS_V2.md)
- [Résumé d'Implémentation](./IMPLEMENTATION_SUMMARY.md)
- [Migration Complète](./MIGRATION_COMPLETED.md)

---

**Version** : 2.0.0  
**Dernière mise à jour** : 2025-12-02  
**Statut** : Production Ready