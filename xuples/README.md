# Package Xuples - Documentation

## 📋 Vue d'Ensemble

Le package `xuples` implémente le système de xuple-space pour TSD, permettant de publier et gérer les activations de règles RETE de manière découplée et configurable.

## 🎯 Concepts

### Xuple
Une **activation de règle** disponible dans un xuple-space, contenant :
- L'action RETE déclenchée
- Le token avec tous les bindings
- Les faits déclencheurs
- Un statut (pending, consumed, expired, archived)
- Des métadonnées de tracking

### XupleSpace
Un **espace nommé** gérant des xuples avec des politiques configurables :
- **Sélection** : comment choisir parmi les xuples disponibles (FIFO, LIFO, random)
- **Consommation** : règles de consommation (once, unlimited, limited)
- **Rétention** : durée de vie des xuples (unlimited, duration)

### XupleManager
Gestionnaire global permettant de créer et gérer plusieurs xuple-spaces.

## 🚀 Utilisation

### Exemple Basique

```go
package main

import (
    "fmt"
    "github.com/treivax/tsd/xuples"
    "github.com/treivax/tsd/rete"
)

func main() {
    // Créer un gestionnaire
    manager := xuples.NewXupleManager()
    
    // Créer un xuple-space avec politiques par défaut
    space, err := manager.CreateSpace("alerts", nil)
    if err != nil {
        panic(err)
    }
    
    // Ajouter un xuple (normalement fait par RETE)
    action := &rete.Action{/* ... */}
    token := &rete.Token{/* ... */}
    facts := []*rete.Fact{/* ... */}
    
    xuple, err := space.Add(action, token, facts, manager)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Xuple créé: %s\n", xuple.ID)
    
    // Consommer un xuple
    consumed, err := space.Consume("agent-1", nil)
    if err != nil {
        panic(err)
    }
    
    if consumed != nil {
        fmt.Printf("Xuple consommé par agent-1: %s\n", consumed.ID)
    }
}
```

### Exemple avec Politiques Personnalisées

```go
config := xuples.PolicyConfig{
    Selection: xuples.SelectionPolicyConfig{
        Type: "fifo", // Premier arrivé, premier servi
    },
    Consumption: xuples.ConsumptionPolicyConfig{
        Type:            "limited",
        MaxConsumptions: 3, // Jusqu'à 3 consommations
    },
    Retention: xuples.RetentionPolicyConfig{
        Type:     "duration",
        Duration: 5 * time.Minute, // Expire après 5 minutes
    },
}

space, err := manager.CreateSpace("limited-alerts", &config)
```

### Filtrage

```go
// Consommer seulement les xuples pending
consumed, err := space.Consume("agent-1", xuples.FilterByStatus(xuples.StatusPending))

// Consommer seulement les xuples d'une action spécifique
consumed, err := space.Consume("agent-1", xuples.FilterByAction("send_notification"))

// Filtre personnalisé
consumed, err := space.Consume("agent-1", func(x *xuples.Xuple) bool {
    return x.Status == xuples.StatusPending && x.ConsumptionCount == 0
})
```

### Statistiques

```go
stats := space.GetStats()
fmt.Printf("Total créés: %d\n", stats.TotalCreated)
fmt.Printf("Total consommés: %d\n", stats.TotalConsumed)
fmt.Printf("Pending actuels: %d\n", stats.CurrentPending)
```

### Nettoyage

```go
// Nettoyer les xuples expirés
cleaned := space.Cleanup()
fmt.Printf("%d xuples nettoyés\n", cleaned)

// Vider complètement le space
space.Clear()
```

## 📋 Politiques Disponibles

### Sélection

| Type | Description |
|------|-------------|
| `fifo` | Premier arrivé, premier servi (par CreatedAt) |
| `lifo` | Dernier arrivé, premier servi (par CreatedAt) |
| `random` | Sélection aléatoire |

### Consommation

| Type | Description | Paramètres |
|------|-------------|------------|
| `once` | Une seule consommation | - |
| `unlimited` | Consommations illimitées | - |
| `limited` | Nombre limité de consommations | `MaxConsumptions` |

### Rétention

| Type | Description | Paramètres |
|------|-------------|------------|
| `unlimited` | Pas d'expiration | - |
| `duration` | Expire après une durée | `Duration` |

## 🔒 Thread-Safety

- Toutes les opérations sont thread-safe via `sync.RWMutex`
- Génération d'IDs via compteur atomique
- Les xuples sont immuables après création (sauf Status et ConsumedBy)

## 🧪 Tests

```bash
# Exécuter les tests
go test ./xuples/...

# Avec couverture
go test ./xuples/... -cover

# Avec verbosité
go test ./xuples/... -v
```

## 📚 Documentation

- [Structures de données](../docs/xuples/design/01-data-structures.md)
- [Code Review](../REPORTS/code-review-refactoring-xuples-2025-12-17.md)
- GoDoc : `go doc github.com/treivax/tsd/xuples`

## 🎯 Intégration RETE (Future)

Le module xuples est conçu pour s'intégrer avec RETE via un publisher :

```go
// Dans rete/node_terminal.go (exemple futur)
if network.XuplePublisher != nil {
    network.XuplePublisher.Publish(tn.Action, token, token.Facts)
}
```

Cela permettra :
- Découplage complet entre RETE et xuples
- Publication asynchrone possible
- Agents externes récupérant les activations

## ⚠️ TODO

- [ ] Implémenter `TupleSpacePublisher` interface
- [ ] Intégrer avec `TerminalNode`
- [ ] Ajouter configuration enable/disable
- [ ] Support pour sérialisation JSON/YAML
- [ ] API REST pour agents externes (futur)
- [ ] Métriques Prometheus (futur)

## 📝 Licence

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License
