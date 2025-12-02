# Transactions Obligatoires dans TSD

## Vue d'ensemble

À partir de cette version, **les transactions sont OBLIGATOIRES et AUTOMATIQUES** dans tout le pipeline d'ingestion TSD. Cette décision architecturale garantit la cohérence des données et simplifie l'utilisation de l'API.

## 🔒 Changements Majeurs

### 1. Transactions Automatiques

Toutes les fonctions d'ingestion utilisent maintenant **automatiquement** les transactions :

- `IngestFile()` : Transaction automatique avec auto-commit/auto-rollback
- `IngestFileWithMetrics()` : Transaction automatique avec collecte de métriques
- `IngestFileWithAdvancedFeatures()` : Transaction automatique avec configuration avancée

### 2. Gestion Automatique du Cycle de Vie

```
┌─────────────┐
│   Parsing   │
└──────┬──────┘
       │
       ▼
┌─────────────────┐
│ BeginTransaction│ ◄── AUTOMATIQUE
└──────┬──────────┘
       │
       ▼
┌─────────────┐
│ Validation  │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Ingestion  │
└──────┬──────┘
       │
   ┌───┴───┐
   │       │
   ▼       ▼
┌──────┐ ┌──────────┐
│Commit│ │ Rollback │ ◄── AUTOMATIQUE selon le résultat
└──────┘ └──────────┘
```

### 3. Comportement en Cas d'Erreur

**Avant** (transactions optionnelles) :
```go
network, err := pipeline.IngestFile(filename, network, storage)
if err != nil {
    // ❌ État du réseau possiblement inconsistant
    // ❌ Pas de moyen de revenir en arrière
}
```

**Maintenant** (transactions obligatoires) :
```go
network, err := pipeline.IngestFile(filename, network, storage)
if err != nil {
    // ✅ Rollback automatique effectué
    // ✅ Le réseau est dans son état initial
    // ✅ Aucune corruption de données possible
}
```

## 📋 API Simplifiée

### Fonction Principale : `IngestFile`

```go
// AVANT : Plusieurs variantes, transactions optionnelles
network, err := pipeline.IngestFile(filename, network, storage)                    // Sans transaction
// Ces fonctions ont été SUPPRIMÉES :
// network, err := pipeline.IngestFileWithTransaction(filename, network, storage)
// err := pipeline.IngestFileTransactional(filename, network, storage, tx)

// MAINTENANT : Une seule fonction, transactions automatiques
network, err := pipeline.IngestFile(filename, network, storage)
// ✅ Transaction créée automatiquement
// ✅ Commit automatique si succès
// ✅ Rollback automatique si erreur
```

### Fonction avec Métriques : `IngestFileWithMetrics`

```go
network, metrics, err := pipeline.IngestFileWithMetrics(filename, network, storage)
// ✅ Transaction automatique + métriques détaillées
```

### Fonction Avancée : `IngestFileWithAdvancedFeatures`

```go
config := rete.DefaultAdvancedPipelineConfig()
config.AutoCommit = true
config.AutoRollbackOnError = true

network, metrics, err := pipeline.IngestFileWithAdvancedFeatures(filename, network, storage, config)
// ✅ Transaction automatique + validation incrémentale + GC + métriques avancées
```

## 🔧 Configuration des Transactions

### `AdvancedPipelineConfig`

Les transactions ne peuvent plus être désactivées. Les options suivantes restent configurables :

```go
type AdvancedPipelineConfig struct {
    // Transactions (toujours activées - non configurable)
    TransactionTimeout  time.Duration  // Timeout de la transaction (défaut: 30s)
    MaxTransactionSize  int64          // Taille max de l'empreinte mémoire (défaut: 100 MB)
    AutoCommit          bool           // Commit automatique (défaut: false)
    AutoRollbackOnError bool           // Rollback automatique sur erreur (défaut: true)
}
```

### Différence `AutoCommit` true vs false

**`AutoCommit = false`** (défaut) :
```go
network, metrics, err := pipeline.IngestFileWithAdvancedFeatures(filename, network, storage, config)
if err != nil {
    // ✅ Rollback déjà effectué automatiquement
} else {
    // ⚠️  Transaction ACTIVE mais pas committée
    // L'utilisateur doit commit manuellement si nécessaire
    tx := network.GetTransaction()
    tx.Commit()
}
```

**`AutoCommit = true`** :
```go
network, metrics, err := pipeline.IngestFileWithAdvancedFeatures(filename, network, storage, config)
if err != nil {
    // ✅ Rollback déjà effectué automatiquement
} else {
    // ✅ Commit déjà effectué automatiquement
    // Rien à faire, tout est terminé
}
```

## 🔄 Migration depuis l'Ancienne API

### Cas 1 : Utilisation simple sans transaction

**Avant** :
```go
network, err := pipeline.IngestFile(filename, network, storage)
```

**Après** :
```go
// ✅ Aucun changement nécessaire !
network, err := pipeline.IngestFile(filename, network, storage)
// La transaction est maintenant automatique
```

### Cas 2 : Utilisation avec transaction manuelle

**Avant** :
```go
// ❌ Cette approche n'est plus possible (fonction supprimée)
// Les transactions sont maintenant automatiques dans IngestFile()
network, err := pipeline.IngestFile(filename, network, storage)
// ✅ Transaction gérée automatiquement
```

**Après** :
```go
// ✅ Simplification drastique
network, err := pipeline.IngestFile(filename, network, storage)
// Transaction automatique avec commit/rollback automatique
```

### Cas 3 : Gestion fine des transactions

**Avant** :
```go
config := DefaultAdvancedPipelineConfig()
config.EnableTransactions = true  // Optionnel
config.AutoCommit = true
```

**Après** :
```go
config := DefaultAdvancedPipelineConfig()
// EnableTransactions a été supprimé (toujours activé)
config.AutoCommit = true
```

## 📊 Métriques de Transaction

Les métriques incluent maintenant toujours les informations de transaction :

```go
type AdvancedMetrics struct {
    // Transaction (toujours présente)
    TransactionID        string        // ID unique de la transaction
    TransactionFootprint int64         // Empreinte mémoire de la transaction
    ChangesTracked       int           // Nombre de commandes enregistrées
    RollbackPerformed    bool          // True si rollback effectué
    RollbackDuration     time.Duration // Durée du rollback
    TransactionDuration  time.Duration // Durée totale de la transaction
}
```

**Affichage des métriques** :
```go
rete.PrintAdvancedMetrics(metrics)
```

Sortie :
```
📊 MÉTRIQUES AVANCÉES
═══════════════════════════════════════
🔒 Transaction
   ID: 550e8400-e29b-41d4-a716-446655440000
   Durée: 125ms
   Empreinte mémoire: 2.34 KB
   Changements trackés: 15
═══════════════════════════════════════
```

## 🎯 Avantages des Transactions Obligatoires

### 1. **Cohérence Garantie**
- Aucun état intermédiaire corrompu possible
- Rollback automatique en cas d'erreur
- Isolation complète des modifications

### 2. **API Simplifiée**
- Une seule fonction `IngestFile` pour tous les cas
- Pas besoin de gérer manuellement les transactions
- Moins de code boilerplate

### 3. **Sécurité Renforcée**
- Impossible d'oublier un rollback
- Protection contre les erreurs de parsing/validation
- Gestion automatique des erreurs

### 4. **Performance**
- Empreinte mémoire minimale (< 1% du réseau)
- BeginTransaction en O(1) constant
- Rollback en O(k) où k = nombre de commandes (pas O(N) du réseau)

### 5. **Observabilité**
- Métriques de transaction toujours disponibles
- Traçabilité complète des modifications
- ID unique de transaction pour le debugging

## ⚠️ Breaking Changes

### Changements dans `AdvancedPipelineConfig`

| Avant | Après | Migration |
|-------|-------|-----------|
| `EnableTransactions bool` | ❌ Supprimé | Retirer cette ligne (toujours activé) |
| `SnapshotSize int64` (métriques) | `TransactionFootprint int64` | Renommer dans votre code |
| `TransactionUsed bool` (métriques) | ❌ Supprimé | Toujours true maintenant |

### Fonctions Supprimées

| Fonction | Statut | Alternative |
|----------|--------|-------------|
| `IngestFileTransactional()` | ❌ SUPPRIMÉ | Utiliser `IngestFile()` |
| `IngestFileWithTransaction()` | ❌ SUPPRIMÉ | Utiliser `IngestFile()` |

Ces fonctions ont été complètement supprimées. Utilisez `IngestFile()` qui gère automatiquement les transactions.

## 🧪 Tests

Tous les tests de transactions passent avec succès :

```bash
go test ./rete -run Transaction -v
```

Résultat :
```
✅ TestTransaction_CommitAppliesChanges
✅ TestTransaction_RollbackRevertsAllChanges
✅ TestTransaction_MultipleOperations
✅ TestTransaction_MemoryScalability
✅ TestTransaction_TimeScalability
✅ TestTransaction_RollbackScalability
✅ TestTransaction_LargeNumberOfOperations
✅ TestTransaction_CommitMemoryRelease
... (31 tests au total)
```

## 📚 Exemples Complets

### Exemple 1 : Ingestion Simple

```go
package main

import (
    "fmt"
    "github.com/treivax/tsd/rete"
)

func main() {
    storage := rete.NewMemoryStorage()
    pipeline := rete.NewConstraintPipeline()
    
    // Ingestion avec transaction automatique
    network, err := pipeline.IngestFile("rules.tsd", nil, storage)
    if err != nil {
        // ✅ Rollback déjà effectué
        fmt.Printf("Erreur: %v\n", err)
        return
    }
    
    // ✅ Commit automatique déjà effectué
    fmt.Println("Ingestion réussie!")
}
```

### Exemple 2 : Ingestion avec Métriques

```go
package main

import (
    "fmt"
    "github.com/treivax/tsd/rete"
)

func main() {
    storage := rete.NewMemoryStorage()
    pipeline := rete.NewConstraintPipeline()
    
    // Ingestion avec métriques et transaction automatique
    network, metrics, err := pipeline.IngestFileWithMetrics("rules.tsd", nil, storage)
    if err != nil {
        fmt.Printf("Erreur: %v\n", err)
        return
    }
    
    // Afficher les métriques (incluent la transaction)
    fmt.Printf("Parsing: %v\n", metrics.ParsingDuration)
    fmt.Printf("Validation: %v\n", metrics.ValidationDuration)
    fmt.Printf("Total: %v\n", metrics.TotalDuration)
}
```

### Exemple 3 : Configuration Avancée

```go
package main

import (
    "fmt"
    "time"
    "github.com/treivax/tsd/rete"
)

func main() {
    storage := rete.NewMemoryStorage()
    pipeline := rete.NewConstraintPipeline()
    
    // Configuration personnalisée
    config := rete.DefaultAdvancedPipelineConfig()
    config.TransactionTimeout = 60 * time.Second
    config.MaxTransactionSize = 200 * 1024 * 1024 // 200 MB
    config.AutoCommit = true
    config.AutoRollbackOnError = true
    
    network, metrics, err := pipeline.IngestFileWithAdvancedFeatures(
        "rules.tsd", nil, storage, config,
    )
    
    if err != nil {
        fmt.Printf("Erreur: %v\n", err)
        return
    }
    
    // Afficher les métriques avancées
    rete.PrintAdvancedMetrics(metrics)
}
```

## 🔍 Debugging

### Vérifier l'État de la Transaction

```go
// Récupérer la transaction courante du réseau
tx := network.GetTransaction()

if tx != nil {
    fmt.Printf("Transaction ID: %s\n", tx.ID)
    fmt.Printf("Active: %t\n", tx.IsActive)
    fmt.Printf("Committed: %t\n", tx.IsCommitted)
    fmt.Printf("Rolled back: %t\n", tx.IsRolledBack)
    fmt.Printf("Commands: %d\n", tx.GetCommandCount())
}
```

### Inspecter les Commandes

```go
tx := network.GetTransaction()
commands := tx.GetCommands()

for i, cmd := range commands {
    fmt.Printf("Command %d: %s\n", i, cmd.String())
}
```

## 📖 Références

- [Architecture des Transactions](./TRANSACTION_ARCHITECTURE.md)
- [Guide d'Utilisation des Transactions](./TRANSACTION_README.md)
- [Tests de Scalabilité](../rete/transaction_scalability_test.go)
- [Benchmarks de Performance](../rete/transaction_benchmark_test.go)

## 🤝 Support

Pour toute question ou problème lié aux transactions obligatoires, veuillez :

1. Consulter cette documentation
2. Exécuter les tests : `go test ./rete -run Transaction -v`
3. Vérifier les logs de transaction dans la sortie console
4. Ouvrir une issue sur le projet avec les logs complets

---

**Dernière mise à jour** : 2025-12-02  
**Version** : TSD v2.0 - Transactions Obligatoires