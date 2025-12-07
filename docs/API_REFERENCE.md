# API Reference - Pipeline TSD

## Vue d'ensemble

Ce document liste toutes les fonctions publiques disponibles pour utiliser le pipeline TSD avec transactions automatiques obligatoires.

**Version** : 2.0.0 - Transactions Obligatoires  
**Date** : 2025-12-02

---

## 📚 Table des Matières

1. [Fonction Principale](#fonction-principale)
2. [Fonctions avec Métriques](#fonctions-avec-métriques)
3. [Fonctions de Construction](#fonctions-de-construction)
4. [Configuration](#configuration)
5. [Types et Structures](#types-et-structures)

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
    FactSubmissionDuration time.Duration
    TotalDuration        time.Duration
    
    TypesAdded       int
    RulesAdded       int
    FactsSubmitted   int
    TokensGenerated  int
    ActivationsProduced int
    
    ResetDetected    bool
    NodesCollected   int
    GCDuration       time.Duration
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

## Fonctions de Construction

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

### `IngestionMetrics`

Structure contenant les métriques d'une ingestion avec `IngestFileWithMetrics()`.

```go
type IngestionMetrics struct {
    ParsingDuration        time.Duration
    ValidationDuration     time.Duration
    NetworkBuildDuration   time.Duration
    FactSubmissionDuration time.Duration
    TotalDuration          time.Duration
    
    TypesAdded          int
    RulesAdded          int
    FactsSubmitted      int
    TokensGenerated     int
    ActivationsProduced int
    
    ResetDetected  bool
    NodesCollected int
    GCDuration     time.Duration
}
```

---

## Guide de Sélection

### Quelle fonction utiliser ?

| Besoin | Fonction Recommandée |
|--------|---------------------|
| **Cas général** | `IngestFile()` |
| **Besoin de métriques** | `IngestFileWithMetrics()` |
| **Plusieurs fichiers** | Appels successifs à `IngestFile()` |

### Fonctionnalités

Les deux fonctions offrent les mêmes fonctionnalités :
- ✅ **Transaction automatique** : Créée, committée ou rollbackée automatiquement
- ✅ **Validation incrémentale** : Utilise le contexte du réseau existant
- ✅ **GC après reset** : Nettoyage automatique après commande `reset`
- ✅ **Propagation de faits** : Les faits existants sont propagés aux nouvelles règles

La seule différence : `IngestFileWithMetrics()` retourne des métriques détaillées.

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
    fmt.Printf("Construction réseau : %v\n", metrics.NetworkBuildDuration)
    fmt.Printf("Total : %v\n", metrics.TotalDuration)
    fmt.Printf("Types ajoutés : %d\n", metrics.TypesAdded)
    fmt.Printf("Règles ajoutées : %d\n", metrics.RulesAdded)
    fmt.Printf("Faits soumis : %d\n", metrics.FactsSubmitted)
    
    if metrics.ResetDetected {
        fmt.Printf("Reset détecté - GC effectué : %d nœuds en %v\n", 
            metrics.NodesCollected, metrics.GCDuration)
    }
}
```

### Exemple 3 : Ingestion Incrémentale Multi-Fichiers

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
    
    // Charger les types
    network, err := pipeline.IngestFile("types.tsd", nil, storage)
    if err != nil {
        log.Fatalf("Erreur types : %v", err)
    }
    fmt.Printf("Types chargés : %d\n", len(network.Types))
    
    // Charger les règles (validation incrémentale)
    network, err = pipeline.IngestFile("rules.tsd", network, storage)
    if err != nil {
        log.Fatalf("Erreur règles : %v", err)
    }
    fmt.Printf("Règles chargées : %d\n", len(network.TerminalNodes))
    
    // Charger les faits
    network, err = pipeline.IngestFile("facts.tsd", network, storage)
    if err != nil {
        log.Fatalf("Erreur faits : %v", err)
    }
    
    fmt.Println("Ingestion multi-fichiers réussie !")
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