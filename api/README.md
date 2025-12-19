# Package API - TSD

Le package `api` fournit une interface simplifiée pour utiliser le moteur de règles TSD.

## 🎯 Objectif

Ce package est le **point d'entrée recommandé** pour toutes les applications utilisant TSD. Il encapsule la complexité des packages `rete`, `xuples`, et `constraint`, et fournit une API simple et intuitive.

## ✨ Caractéristiques

- **API Simple** : Une seule ligne pour créer un pipeline
- **Configuration Automatique** : Gestion automatique des xuple-spaces et des actions
- **Ingestion Incrémentale** : Support de l'ajout progressif de règles et de faits
- **Thread-Safe** : Utilisation sécurisée en environnement concurrent
- **Métriques Intégrées** : Collecte automatique des statistiques de performance
- **Gestion d'Erreurs** : Erreurs détaillées avec position dans le fichier source

## 📦 Installation

```go
import "github.com/treivax/tsd/api"
```

## 🚀 Utilisation Rapide

### Exemple Basique

```go
package main

import (
    "fmt"
    "log"
    "github.com/treivax/tsd/api"
)

func main() {
    // Créer un pipeline
    pipeline := api.NewPipeline()

    // Ingérer un programme TSD
    result, err := pipeline.IngestFile("program.tsd")
    if err != nil {
        log.Fatal(err)
    }

    // Afficher les résultats
    fmt.Printf("Types définis: %d\n", result.TypeCount())
    fmt.Printf("Règles actives: %d\n", result.RuleCount())
    fmt.Printf("Faits dans le réseau: %d\n", result.FactCount())
}
```

### Ingestion Incrémentale

```go
pipeline := api.NewPipeline()

// Charger les types
_, err := pipeline.IngestFile("types.tsd")
if err != nil {
    log.Fatal(err)
}

// Ajouter des règles
_, err = pipeline.IngestFile("rules.tsd")
if err != nil {
    log.Fatal(err)
}

// Soumettre des faits
result, err := pipeline.IngestFile("facts.tsd")
if err != nil {
    log.Fatal(err)
}
```

### Configuration Personnalisée

```go
config := &api.Config{
    LogLevel:          api.LogLevelDebug,
    EnableMetrics:     true,
    MaxFactsInMemory:  100000,
    XupleSpaceDefaults: &api.XupleSpaceDefaults{
        Selection:   api.SelectionFIFO,
        Consumption: api.ConsumptionOnce,
        Retention:   api.RetentionUnlimited,
    },
}

pipeline := api.NewPipelineWithConfig(config)
result, err := pipeline.IngestFile("program.tsd")
```

### Accès aux Xuples

```go
result, _ := pipeline.IngestFile("monitoring.tsd")

// Récupérer tous les xuples d'un xuple-space
alerts, err := result.GetXuples("critical_alerts")
if err != nil {
    log.Fatal(err)
}

for _, xuple := range alerts {
    fmt.Printf("Alert: %v\n", xuple.Fact.Fields)
}

// Consommer un xuple (retrieve)
xuple, err := result.Retrieve("critical_alerts", "agent1")
if err == nil {
    fmt.Printf("Consumed: %v\n", xuple.Fact.Fields)
}
```

## 📊 Métriques

Les métriques d'ingestion sont collectées automatiquement :

```go
result, _ := pipeline.IngestFile("program.tsd")
metrics := result.Metrics()

fmt.Printf("Temps de parsing: %v\n", metrics.ParseDuration)
fmt.Printf("Temps de construction réseau: %v\n", metrics.BuildDuration)
fmt.Printf("Nombre de propagations: %d\n", metrics.PropagationCount)
```

## 🔧 Configuration

### Niveaux de Logs

- `LogLevelSilent` : Aucun log
- `LogLevelError` : Erreurs uniquement
- `LogLevelWarn` : Erreurs et avertissements
- `LogLevelInfo` : Informations, erreurs et avertissements (défaut)
- `LogLevelDebug` : Tous les logs y compris debug

### Politiques Xuple-Spaces

**Sélection** :
- `SelectionFIFO` : Premier arrivé, premier servi (défaut)
- `SelectionLIFO` : Dernier arrivé, premier servi
- `SelectionRandom` : Sélection aléatoire

**Consommation** :
- `ConsumptionOnce` : Chaque xuple peut être consommé une seule fois (défaut)
- `ConsumptionPerAgent` : Chaque agent peut consommer chaque xuple une fois

**Rétention** :
- `RetentionUnlimited` : Conservation illimitée (défaut)
- `RetentionDuration` : Conservation pendant une durée limitée

## ⚠️ Gestion d'Erreurs

Les erreurs sont détaillées et typées :

```go
_, err := pipeline.IngestFile("invalid.tsd")
if err != nil {
    switch e := err.(type) {
    case *api.ParseError:
        fmt.Printf("Erreur de parsing ligne %d, colonne %d: %s\n",
            e.Line, e.Column, e.Message)
    case *api.XupleSpaceError:
        fmt.Printf("Erreur xuple-space '%s': %s\n",
            e.SpaceName, e.Message)
    case *api.ConfigError:
        fmt.Printf("Erreur de configuration '%s': %s\n",
            e.Field, e.Message)
    default:
        fmt.Printf("Erreur: %v\n", err)
    }
}
```

## 🧵 Thread Safety

Le pipeline est thread-safe. Plusieurs goroutines peuvent appeler `IngestFile` en parallèle :

```go
pipeline := api.NewPipeline()

var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        filename := fmt.Sprintf("file%d.tsd", id)
        pipeline.IngestFile(filename)
    }(i)
}
wg.Wait()
```

⚠️ **Note** : L'ordre d'exécution des règles peut varier en parallèle. Pour un contrôle strict, utilisez un seul goroutine.

## 🏗️ Architecture

```
api (High-Level API)
├── Pipeline         → Point d'entrée principal
├── Result           → Résultats d'ingestion
├── Config           → Configuration
└── Errors           → Gestion d'erreurs typées
    ↓
rete (RETE Engine)
├── ReteNetwork      → Réseau RETE
├── ConstraintPipeline → Pipeline d'ingestion
└── Actions          → Exécuteurs d'actions
    ↓
xuples (Xuple Management)
├── XupleManager     → Gestionnaire de xuple-spaces
├── XupleSpace       → Espace de stockage de xuples
└── Policies         → Politiques de sélection/consommation
    ↓
constraint (Parser)
└── ParseConstraint  → Parser PEG pour fichiers TSD
```

## 📖 Documentation

- [Documentation complète du package](https://pkg.go.dev/github.com/treivax/tsd/api)
- [Guide TSD principal](../README.md)
- [Documentation RETE](../rete/README.md)
- [Documentation Xuples](../xuples/README.md)

## 🧪 Tests

Exécuter les tests :

```bash
cd api
go test -v
```

Avec couverture :

```bash
go test -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 📝 Exemples

Voir le fichier [examples_test.go](examples_test.go) pour des exemples complets et testables.

## 🤝 Contribution

Ce package suit les standards définis dans [CONTRIBUTING.md](../CONTRIBUTING.md).

## 📄 Licence

Copyright (c) 2025 TSD Contributors - Licence MIT

Voir le fichier [LICENSE](../LICENSE) pour plus de détails.
