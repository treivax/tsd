# Démarrage rapide - Optimisations avancées

## Installation

Les optimisations avancées sont intégrées au package principal. Aucune installation supplémentaire nécessaire.

```bash
# Vérifier que tout compile
go build ./rete/...
```

## Validation rapide

```bash
# Exécuter le script de validation
./validate_advanced_features.sh
```

## Utilisation basique

### 1. Validation incrémentale

```go
package main

import (
    "github.com/treivax/tsd/rete"
    "log"
)

func main() {
    pipeline := rete.NewConstraintPipeline()
    storage := rete.NewMemoryStorage()
    
    // Charger avec validation
    network, err := pipeline.IngestFileWithIncrementalValidation("file.tsd", nil, storage)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 2. Avec GC

```go
network, err := pipeline.IngestFileWithGC("reset_file.tsd", network, storage)
```

### 3. Avec transactions

```go
network, err := pipeline.IngestFile("file.tsd", network, storage)
// ✅ Transaction automatique avec auto-commit/auto-rollback
```

### 4. Configuration complète

```go
config := rete.DefaultAdvancedPipelineConfig()
network, metrics, err := pipeline.IngestFileWithAdvancedFeatures(
    "file.tsd",
    network,
    storage,
    config,
)

if err != nil {
    log.Printf("Error: %v", err)
    if metrics.RollbackPerformed {
        log.Printf("Rollback: %v", metrics.RollbackDuration)
    }
}

rete.PrintAdvancedMetrics(metrics)
```

## Exemple complet

```bash
# Compiler et exécuter l'exemple
go run examples/advanced_features_example.go
```

## Tests

```bash
# Tests d'intégration
cd test/integration/incremental/
go test -v -run Advanced

# Test spécifique
go test -v -run TestIncrementalValidation
```

## Documentation

- **Guide complet** : [docs/ADVANCED_FEATURES_README.md](docs/ADVANCED_FEATURES_README.md)
- **Synthèse** : [ADVANCED_FEATURES_SUMMARY.md](ADVANCED_FEATURES_SUMMARY.md)
- **Spécifications** : [docs/ADVANCED_OPTIMIZATIONS.md](docs/ADVANCED_OPTIMIZATIONS.md)

## Commandes utiles

```bash
# Validation complète
./validate_advanced_features.sh

# Compilation
go build ./rete/...

# Tests
go test ./test/integration/incremental/

# Exemple
go run examples/advanced_features_example.go

# Benchmark (optionnel)
go test -bench=. ./rete/...
```

## Support

En cas de problème :
1. Consulter la documentation ci-dessus
2. Vérifier les tests pour des exemples
3. Exécuter le script de validation

**Status** : Production Ready ✅
