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

**Fonction UNIQUE** pour tous les cas d'usage. C'est la seule fonction d'ingestion.

```go
func (cp *ConstraintPipeline) IngestFile(
    filename string,
    network *ReteNetwork,
    storage Storage,
) (*ReteNetwork, *IngestionMetrics, error)
```

**Description** :  
Ingère un fichier de contraintes TSD dans le réseau RETE avec **transaction automatique obligatoire**.  
**Les métriques sont toujours collectées** (coût négligeable < 0.1%).

**Caractéristiques** :
- ✅ **Transaction automatique** : Créée, committée ou rollbackée automatiquement
- ✅ **Validation incrémentale** : Utilise le contexte du réseau existant
- ✅ **GC automatique** : Après détection d'une commande `reset`
- ✅ **Propagation de faits** : Les faits existants sont propagés aux nouvelles règles
- ✅ **Métriques incluses** : Toujours retournées sans impact sur les performances

**Paramètres** :
- `filename` : Chemin vers le fichier `.tsd` à ingérer
- `network` : Réseau RETE existant ou `nil` (un nouveau sera créé)
- `storage` : Interface de stockage des faits

**Retour** :
- `*ReteNetwork` : Le réseau RETE mis à jour
- `*IngestionMetrics` : Métriques détaillées de l'ingestion (toujours collectées)
- `error` : Erreur en cas d'échec (rollback automatique effectué)

**Exemple** :
```go
storage := rete.NewMemoryStorage()
pipeline := rete.NewConstraintPipeline()

network, metrics, err := pipeline.IngestFile("rules.tsd", nil, storage)
if err != nil {
    // ✅ Rollback automatique déjà effectué
    log.Fatalf("Erreur : %v", err)
}
// ✅ Commit automatique déjà effectué
fmt.Printf("Ingestion réussie en %v\n", metrics.TotalDuration)
fmt.Printf("Types: %d, Règles: %d, Faits: %d\n", 
    metrics.TypesAdded, metrics.RulesAdded, metrics.FactsSubmitted)
```

**Complexité** :
- Parsing : O(n) où n = taille du fichier
- Validation : O(m) où m = nombre de types/règles
- Transaction : O(1) pour begin, O(k) pour commit/rollback où k = nombre de commandes
- **Collecte de métriques** : O(1) - coût négligeable

---

## Métriques d'Ingestion

### Structure `IngestionMetrics`

Les métriques sont **toujours retournées** par `IngestFile()` sans impact sur les performances (< 0.1%).

**Métriques disponibles** :
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
network, metrics, err := pipeline.IngestFile("rules.tsd", nil, storage)
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

Structure contenant les métriques d'une ingestion avec `IngestFile()`.

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

**Il n'y a qu'UNE SEULE fonction** : `IngestFile()`

| Besoin | Solution |
|--------|----------|
| **Tous les cas** | `IngestFile()` - retourne toujours les métriques |
| **Plusieurs fichiers** | Appels successifs à `IngestFile()` |
| **Ignorer les métriques** | Utiliser `_` : `network, _, err := pipeline.IngestFile(...)` |
| **Utiliser les métriques** | Capturer : `network, metrics, err := pipeline.IngestFile(...)` |

### Fonctionnalités (toutes incluses)

- ✅ **Transaction automatique** : Créée, committée ou rollbackée automatiquement
- ✅ **Validation incrémentale** : Utilise le contexte du réseau existant
- ✅ **GC après reset** : Nettoyage automatique après commande `reset`
- ✅ **Propagation de faits** : Les faits existants sont propagés aux nouvelles règles
- ✅ **Métriques incluses** : Toujours collectées, coût négligeable (< 0.1%)

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
    
    // Les métriques sont toujours retournées
    network, metrics, err := pipeline.IngestFile("rules.tsd", nil, storage)
    if err != nil {
        log.Fatalf("Erreur : %v", err)
    }
    
    fmt.Println("Ingestion réussie !")
    fmt.Printf("Types : %d\n", len(network.Types))
    fmt.Printf("Règles : %d\n", len(network.TerminalNodes))
    fmt.Printf("Durée totale : %v\n", metrics.TotalDuration)
}
```

### Exemple 2 : Utiliser les Métriques pour le Monitoring

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
    
    **Exemple d'utilisation des métriques** :
    ```go
    network, metrics, err := pipeline.IngestFile("rules.tsd", nil, storage)
    if err != nil {
        log.Fatalf("Erreur : %v", err)
    }

    // Afficher les métriques de performance
    fmt.Printf("⏱️  Durées:\n")
    fmt.Printf("  Parsing : %v\n", metrics.ParsingDuration)
    fmt.Printf("  Validation : %v\n", metrics.ValidationDuration)
    fmt.Printf("  Création types : %v\n", metrics.TypeCreationDuration)
    fmt.Printf("  Création règles : %v\n", metrics.RuleCreationDuration)
    fmt.Printf("  Total : %v\n", metrics.TotalDuration)

    fmt.Printf("\n📊 Compteurs:\n")
    fmt.Printf("  Types ajoutés : %d\n", metrics.TypesAdded)
    fmt.Printf("  Règles ajoutées : %d\n", metrics.RulesAdded)
    fmt.Printf("  Faits soumis : %d\n", metrics.FactsSubmitted)
    fmt.Printf("  Faits propagés : %d\n", metrics.FactsPropagated)

    if metrics.WasReset {
        fmt.Printf("\n🔄 Reset détecté - Ancien réseau nettoyé\n")
    }

    // Identifier les goulots d'étranglement
    fmt.Printf("\n🎯 Goulot : %s\n", metrics.GetBottleneck())
    ```

    **Méthodes utiles** :
    ```go
    // Affichage formaté complet
    fmt.Println(metrics.String())

    // Résumé court
    fmt.Println(metrics.Summary())

    // Vérifier l'efficacité
    if metrics.IsEfficient() {
        fmt.Println("✅ Ingestion efficace")
    }

    // Identifier le goulot d'étranglement
    bottleneck := metrics.GetBottleneck()
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
    network, metrics1, err := pipeline.IngestFile("types.tsd", nil, storage)
    if err != nil {
        log.Fatalf("Erreur types : %v", err)
    }
    fmt.Printf("Types chargés : %d (en %v)\n", 
        len(network.Types), metrics1.TotalDuration)
    
    // Charger les règles (validation incrémentale)
    network, metrics2, err := pipeline.IngestFile("rules.tsd", network, storage)
    if err != nil {
        log.Fatalf("Erreur règles : %v", err)
    }
    fmt.Printf("Règles chargées : %d (en %v)\n", 
        len(network.TerminalNodes), metrics2.TotalDuration)
    
    // Charger les faits
    network, metrics3, err := pipeline.IngestFile("facts.tsd", network, storage)
    if err != nil {
        log.Fatalf("Erreur faits : %v", err)
    }
    
    totalTime := metrics1.TotalDuration + metrics2.TotalDuration + metrics3.TotalDuration
    fmt.Printf("\n✅ Ingestion multi-fichiers réussie en %v\n", totalTime)
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

### Fonctions Supprimées (v2.0.0+)

**Fonctions d'ingestion multiples (supprimées)** :
- ❌ `IngestFileWithMetrics()` : Fusionnée dans `IngestFile()` qui retourne toujours les métriques
- ❌ `IngestFileWithAdvancedFeatures()` : Remplacée par `IngestFile()`
- ❌ `IngestFileTransactionalSafe()` : Remplacée par `IngestFile()`
- ❌ `IngestFileTransactional()` : Remplacée par `IngestFile()`
- ❌ `IngestFileWithTransaction()` : Remplacée par `IngestFile()`

### Fonctions de Construction (supprimées)
- ❌ `BuildNetworkFromConstraintFile()` : Remplacée par `IngestFile(constraintFile, nil, storage)`
- ❌ `BuildNetworkFromMultipleFiles()` : Remplacée par appels successifs à `IngestFile()`
- ❌ `BuildNetworkFromIterativeParser()` : Remplacée par `IngestFile()`
- ❌ `BuildNetworkFromConstraintFileWithFacts()` : Remplacée par deux appels à `IngestFile()`

**Migration** : Utilisez simplement `IngestFile()` qui retourne toujours `(network, metrics, error)`.

**Exemple - Plusieurs fichiers** :
```go
// Avant
network, err := pipeline.BuildNetworkFromMultipleFiles([]string{"types.tsd", "rules.tsd"}, storage)

// Après (noter le retour de metrics)
network, _, err := pipeline.IngestFile("types.tsd", nil, storage)
if err != nil {
    return err
}
network, _, err = pipeline.IngestFile("rules.tsd", network, storage)
```

**Exemple - Contraintes + Faits** :
```go
// Avant
network, facts, err := pipeline.BuildNetworkFromConstraintFileWithFacts("rules.tsd", "facts.tsd", storage)

// Après (noter le retour de metrics)
network, _, err := pipeline.IngestFile("rules.tsd", nil, storage)
if err != nil {
    return err
}
network, _, err = pipeline.IngestFile("facts.tsd", network, storage)
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