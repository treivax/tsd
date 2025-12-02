# Transactions RETE - Guide d'Utilisation

**Version** : 2.0.0 (Command Pattern)  
**Date** : 2025-12-02  
**Status** : Production-ready ✅

---

## 🎯 Qu'est-ce qu'une transaction ?

Une transaction permet de grouper plusieurs opérations sur le réseau RETE et de les **valider** (commit) ou **annuler** (rollback) atomiquement.

### Cas d'usage

- ✅ **Ingestion de fichiers** : Valider seulement si tout le fichier est correct
- ✅ **Modifications batch** : Annuler si une opération échoue
- ✅ **Tests** : Restaurer l'état initial après chaque test
- ✅ **Validation** : Tester des changements sans les appliquer

---

## 🚀 Démarrage rapide

### Installation

```bash
go get github.com/treivax/tsd/rete
```

### Exemple minimal

```go
package main

import (
    "fmt"
    "github.com/treivax/tsd/rete"
)

func main() {
    // Créer un réseau
    storage := rete.NewMemoryStorage()
    network := rete.NewReteNetwork(storage)
    
    // Commencer une transaction
    tx := network.BeginTransaction()
    network.SetTransaction(tx)
    
    // Effectuer des opérations
    err := network.SubmitFact(&rete.Fact{
        ID:   "user1",
        Type: "User",
        Fields: map[string]interface{}{
            "name": "Alice",
            "age":  30,
        },
    })
    
    if err != nil {
        // En cas d'erreur : annuler
        tx.Rollback()
        network.SetTransaction(nil)
        fmt.Printf("Erreur: %v\n", err)
        return
    }
    
    // Tout s'est bien passé : valider
    tx.Commit()
    network.SetTransaction(nil)
    
    fmt.Println("✅ Transaction réussie!")
}
```

---

## 📖 Guide d'utilisation

### 1. Créer une transaction

```go
tx := network.BeginTransaction()
network.SetTransaction(tx)
```

**Complexité** : O(1) - Instantané  
**Mémoire** : ~432 bytes

### 2. Effectuer des opérations

```go
// Toutes les opérations sont automatiquement enregistrées
network.SubmitFact(fact1)
network.SubmitFact(fact2)
network.RemoveFact(factID)
```

**Overhead** : ~200 bytes par opération

### 3A. Valider (Commit)

```go
err := tx.Commit()
if err != nil {
    // Gérer l'erreur
}
network.SetTransaction(nil)
```

Les modifications sont **conservées**.

### 3B. Annuler (Rollback)

```go
err := tx.Rollback()
if err != nil {
    // Gérer l'erreur
}
network.SetTransaction(nil)
```

Les modifications sont **annulées**, le réseau revient à l'état initial.

---

## 💡 Patterns recommandés

### Pattern 1 : Gestion d'erreur simple

```go
tx := network.BeginTransaction()
network.SetTransaction(tx)

err := doSomeWork(network)
if err != nil {
    tx.Rollback()
    network.SetTransaction(nil)
    return err
}

tx.Commit()
network.SetTransaction(nil)
return nil
```

### Pattern 2 : Defer pour cleanup

```go
tx := network.BeginTransaction()
network.SetTransaction(tx)
defer network.SetTransaction(nil)

// Si une erreur se produit, rollback
committed := false
defer func() {
    if !committed {
        tx.Rollback()
    }
}()

// Effectuer les opérations
if err := doWork(network); err != nil {
    return err
}

// Tout s'est bien passé
tx.Commit()
committed = true
return nil
```

### Pattern 3 : Batch processing

```go
tx := network.BeginTransaction()
network.SetTransaction(tx)
defer network.SetTransaction(nil)

for i, fact := range facts {
    if err := network.SubmitFact(fact); err != nil {
        tx.Rollback()
        return fmt.Errorf("fact %d failed: %w", i, err)
    }
}

return tx.Commit()
```

---

## 📊 Performance

### Overhead mémoire

| Taille réseau | Transaction | Overhead |
|--------------|-------------|----------|
| 1,000 faits  | ~2 KB       | < 1%     |
| 10,000 faits | ~2 KB       | < 0.1%   |
| 100,000 faits| ~2 KB       | < 0.01%  |

### Temps d'exécution

| Opération | Temps | Complexité |
|-----------|-------|------------|
| BeginTransaction | ~270 ns | O(1) |
| SubmitFact (transactionnel) | ~1-5 µs | O(1) |
| Rollback (10 ops) | ~50 µs | O(k)* |
| Rollback (100 ops) | ~500 µs | O(k)* |
| Commit | ~1 µs | O(1) |

\* k = nombre d'opérations

---

## ⚙️ API Reference

### Transaction

#### Création

```go
tx := network.BeginTransaction()
```

Crée une nouvelle transaction. Instantané, O(1).

#### Activation

```go
network.SetTransaction(tx)
```

Active la transaction sur le réseau. Toutes les opérations suivantes seront enregistrées.

#### Désactivation

```go
network.SetTransaction(nil)
```

Désactive la transaction. Les opérations suivantes seront exécutées en mode normal.

#### Commit

```go
err := tx.Commit() error
```

Valide la transaction. Les modifications sont conservées.

**Retourne** : Erreur si transaction non active ou déjà committée.

#### Rollback

```go
err := tx.Rollback() error
```

Annule la transaction. Les modifications sont annulées en rejouant les commandes à l'envers.

**Retourne** : Erreur si transaction non active ou déjà committée.

#### Métriques

```go
// Nombre de commandes enregistrées
count := tx.GetCommandCount()

// Durée depuis le début de la transaction
duration := tx.GetDuration()

// Empreinte mémoire estimée (bytes)
footprint := tx.GetMemoryFootprint()

// Représentation textuelle
str := tx.String()
```

---

## 🧪 Tests

### Tester avec des transactions

```go
func TestMyFeature(t *testing.T) {
    storage := rete.NewMemoryStorage()
    network := rete.NewReteNetwork(storage)
    
    // Sauvegarder l'état initial
    initialCount := len(storage.GetAllFacts())
    
    // Commencer une transaction
    tx := network.BeginTransaction()
    network.SetTransaction(tx)
    
    // Tester la feature
    // ...
    
    // Rollback pour restaurer l'état
    tx.Rollback()
    network.SetTransaction(nil)
    
    // Vérifier que l'état est restauré
    finalCount := len(storage.GetAllFacts())
    if finalCount != initialCount {
        t.Errorf("State not restored: %d != %d", finalCount, initialCount)
    }
}
```

---

## ⚠️ Limitations et contraintes

### Thread-safety

- ✅ Une transaction par réseau à la fois
- ⚠️ Pas de transactions concurrentes sur le même réseau
- ⚠️ Les opérations sur le réseau ne sont pas thread-safe

**Recommandation** : Utiliser un mutex si accès concurrent nécessaire.

### États invalides

Ces opérations retournent une erreur :

- ❌ Commit deux fois
- ❌ Rollback après Commit
- ❌ Rollback deux fois
- ❌ Opérations sur transaction inactive

### Rollback partiel

Si une commande échoue pendant `Rollback()` :
- Le rollback s'arrête
- Le réseau peut être dans un état inconsistant
- **Recommandation** : Implémenter Undo() de façon infaillible

---

## 🔧 Extension

### Ajouter une nouvelle commande

Pour supporter de nouvelles opérations transactionnelles :

1. **Créer la commande**

```go
type MyCommand struct {
    // Champs nécessaires pour Execute et Undo
}

func (c *MyCommand) Execute() error {
    // Effectuer l'opération
}

func (c *MyCommand) Undo() error {
    // Annuler l'opération
}

func (c *MyCommand) String() string {
    return "MyCommand(...)"
}
```

2. **Intégrer dans le réseau**

```go
func (network *ReteNetwork) MyOperation(...) error {
    tx := network.GetTransaction()
    if tx != nil && tx.IsActive {
        cmd := NewMyCommand(...)
        return tx.RecordAndExecute(cmd)
    }
    
    // Mode normal
    return network.myOperationInternal(...)
}
```

---

## 📚 Documentation complète

- **Architecture** : `docs/TRANSACTION_ARCHITECTURE.md`
- **Résumé optimisation** : `docs/TRANSACTION_OPTIMIZATION_SUMMARY.md`
- **Checklist** : `docs/TRANSACTION_IMPLEMENTATION_CHECKLIST.md`
- **Exemple complet** : `examples/transaction_example.go`

---

## 🐛 Dépannage

### "transaction not active"

**Cause** : Tentative d'utiliser une transaction après commit/rollback.

**Solution** : Vérifier que `tx.IsActive == true` avant les opérations.

### "transaction already committed"

**Cause** : Tentative de commit/rollback d'une transaction déjà terminée.

**Solution** : Ne pas réutiliser une transaction après commit.

### Overhead mémoire élevé

**Cause** : Trop d'opérations dans une seule transaction.

**Solution** : 
- Diviser en plusieurs transactions plus petites
- Commit intermédiaires si possible

### Rollback lent

**Cause** : Beaucoup d'opérations à annuler (Rollback = O(k)).

**Solution** :
- Normal si k est grand
- Optimiser en réduisant le nombre d'opérations par transaction

---

## ❓ FAQ

### Puis-je imbriquer des transactions ?

Non, l'implémentation actuelle ne supporte pas les transactions imbriquées.

### Puis-je faire un Redo après Undo ?

Non, après un Rollback, les commandes sont libérées.

### Les transactions sont-elles persistées ?

Non, les transactions sont en mémoire uniquement. Pour la persistence, voir les évolutions futures (WAL).

### Quelle est la différence avec l'ancienne implémentation ?

| Aspect | V1 (Snapshot) | V2 (Command) |
|--------|---------------|--------------|
| Overhead mémoire | ~100% | < 1% |
| BeginTransaction | O(N) | O(1) |
| Scalabilité | < 100k faits | Millions de faits |

---

## 🚀 Évolutions futures

### Court terme
- [ ] AddRuleCommand, RemoveRuleCommand
- [ ] ModifyFactCommand

### Moyen terme
- [ ] Savepoints (rollback partiel)
- [ ] Transactions imbriquées
- [ ] WAL (Write-Ahead Log)

### Long terme
- [ ] Optimistic Concurrency Control
- [ ] Two-Phase Commit
- [ ] Event Sourcing

---

## 📝 Changelog

### v2.0.0 (2025-12-02) - Command Pattern

**Breaking Changes** :
- ❌ Suppression de NetworkSnapshot
- ❌ Suppression de RecordChange(), GetChanges()
- ❌ Suppression de GetSnapshotSize()

**Nouveautés** :
- ✅ Command Pattern avec rejeu inversé
- ✅ BeginTransaction O(1) au lieu de O(N)
- ✅ Overhead mémoire < 1% au lieu de ~100%
- ✅ Scalable jusqu'à millions de faits

**Performance** :
- 🚀 ~99% plus rapide
- 🚀 ~99% moins de mémoire
- 🚀 Production-ready

---

## 📧 Support

- **Issues** : GitHub Issues
- **Documentation** : `docs/TRANSACTION_ARCHITECTURE.md`
- **Exemples** : `examples/transaction_example.go`

---

**Auteurs** : TSD Contributors  
**Licence** : MIT  
**Status** : Production-ready ✅