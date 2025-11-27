# Guide de Démarrage Rapide - Optimisations de Performance

## 🚀 Démarrage en 5 Minutes

### 1. Construction Basique avec Métriques

```go
package main

import (
    "fmt"
    "github.com/treivax/tsd/rete"
)

func main() {
    // Créer le réseau RETE (métriques incluses automatiquement)
    storage := rete.NewMemoryStorage()
    network := rete.NewReteNetwork(storage)
    
    // Construire des règles
    for i := 0; i < 100; i++ {
        conditions := []rete.SimpleCondition{
            {
                Type:     "binaryOperation",
                Left:     map[string]interface{}{"type": "variable", "name": "age"},
                Operator: ">",
                Right:    map[string]interface{}{"type": "literal", "value": float64(18)},
            },
        }
        
        builder := rete.NewAlphaChainBuilderWithMetrics(
            network, 
            storage, 
            network.ChainMetrics,
        )
        
        ruleID := fmt.Sprintf("rule_%d", i)
        _, err := builder.BuildChain(conditions, "person", network.RootNode, ruleID)
        if err != nil {
            panic(err)
        }
    }
    
    // Afficher les métriques
    PrintMetrics(network)
}

func PrintMetrics(network *rete.ReteNetwork) {
    metrics := network.GetChainMetrics()
    summary := metrics.GetSummary()
    
    chains := summary["chains"].(map[string]interface{})
    nodes := summary["nodes"].(map[string]interface{})
    hashCache := summary["hash_cache"].(map[string]interface{})
    
    fmt.Printf("📊 Métriques de Performance\n")
    fmt.Printf("============================\n")
    fmt.Printf("Chaînes construites:    %d\n", chains["total_built"])
    fmt.Printf("Nœuds créés:            %d\n", nodes["total_created"])
    fmt.Printf("Nœuds réutilisés:       %d\n", nodes["total_reused"])
    fmt.Printf("Ratio de partage:       %.2f%%\n", nodes["reuse_rate_pct"])
    fmt.Printf("Efficacité cache hash:  %.2f%%\n", hashCache["efficiency_pct"])
    fmt.Printf("Temps moyen:            %s\n", chains["average_build_time"])
}
```

### 2. Sortie Attendue

```
📊 Métriques de Performance
============================
Chaînes construites:    100
Nœuds créés:            11
Nœuds réutilisés:       189
Ratio de partage:       94.50%
Efficacité cache hash:  94.50%
Temps moyen:            26.612µs
```

## 📈 Cas d'Usage Courants

### Analyser les Chaînes les Plus Lentes

```go
metrics := network.GetChainMetrics()
topSlow := metrics.GetTopChainsByBuildTime(5)

fmt.Println("🐌 Top 5 des chaînes les plus lentes:")
for i, chain := range topSlow {
    fmt.Printf("%d. %s - %v (longueur: %d)\n", 
        i+1, chain.RuleID, chain.BuildTime, chain.ChainLength)
}
```

### Identifier les Chaînes les Plus Longues

```go
metrics := network.GetChainMetrics()
topLong := metrics.GetTopChainsByLength(5)

fmt.Println("📏 Top 5 des chaînes les plus longues:")
for i, chain := range topLong {
    fmt.Printf("%d. %s - %d nœuds\n", 
        i+1, chain.RuleID, chain.ChainLength)
}
```

### Monitoring Continu

```go
import "time"

// Exporter les métriques toutes les minutes
ticker := time.NewTicker(1 * time.Minute)
go func() {
    for range ticker.C {
        metrics := network.GetChainMetrics()
        snapshot := metrics.GetSnapshot()
        
        fmt.Printf("[%s] Chaînes: %d, Partage: %.2f%%, Cache: %.2f%%\n",
            time.Now().Format("15:04:05"),
            snapshot.TotalChainsBuilt,
            snapshot.SharingRatio * 100,
            metrics.GetHashCacheEfficiency() * 100,
        )
    }
}()
```

## 🎯 Métriques Clés à Surveiller

| Métrique | Bon | Moyen | Mauvais | Action |
|----------|-----|-------|---------|--------|
| **Ratio de partage** | >70% | 30-70% | <30% | Vérifier normalisation |
| **Efficacité cache hash** | >80% | 50-80% | <50% | Beaucoup de conditions uniques |
| **Temps moyen** | <50µs | 50-200µs | >200µs | Investiguer la complexité |
| **Taille cache hash** | <10k | 10k-50k | >50k | Considérer nettoyage |

## 🧪 Tests et Benchmarks

### Exécuter les Tests de Performance

```bash
# Tests avec 100 règles
go test -v ./rete -run TestPerformance_LargeRuleset_100Rules

# Tests avec 1000 règles (plus long)
go test -v ./rete -run TestPerformance_LargeRuleset_1000Rules

# Tous les tests de métriques
go test -v ./rete -run TestMetrics_
```

### Exécuter les Benchmarks

```bash
# Benchmarks de construction de chaînes
go test -bench=BenchmarkChainBuild -benchmem ./rete

# Benchmark cache de hash
go test -bench=BenchmarkHashCompute -benchmem ./rete

# Tous les benchmarks
go test -bench=. -benchmem ./rete
```

## 🔧 Configuration Avancée

### Partager les Métriques entre Composants

```go
// Créer des métriques partagées
metrics := rete.NewChainBuildMetrics()

// Les utiliser dans le réseau
network.ChainMetrics = metrics

// Les utiliser dans le registry
registry := rete.NewAlphaSharingRegistryWithMetrics(metrics)
network.AlphaSharingManager = registry

// Les utiliser dans le builder
builder := rete.NewAlphaChainBuilderWithMetrics(network, storage, metrics)
```

### Réinitialiser les Métriques

```go
// Réinitialiser toutes les métriques
network.ResetChainMetrics()

// Vider uniquement le cache de hash
network.AlphaSharingManager.ClearHashCache()

// Vider uniquement le cache de connexion (dans le builder)
builder.ClearConnectionCache()
```

### Exporter les Détails Complets

```go
metrics := network.GetChainMetrics()
snapshot := metrics.GetSnapshot()

// Accéder aux détails de chaque chaîne
for _, detail := range snapshot.ChainDetails {
    fmt.Printf("Règle: %s\n", detail.RuleID)
    fmt.Printf("  Longueur: %d\n", detail.ChainLength)
    fmt.Printf("  Créés: %d, Réutilisés: %d\n", 
        detail.NodesCreated, detail.NodesReused)
    fmt.Printf("  Temps: %v\n", detail.BuildTime)
    fmt.Printf("  Hash: %v\n", detail.HashesGenerated)
}
```

## ⚠️ Pièges à Éviter

### ❌ Ne Pas Faire

```go
// Créer un nouveau builder à chaque fois (perd le cache!)
for i := range rules {
    builder := rete.NewAlphaChainBuilder(network, storage)
    builder.BuildChain(...)
}
```

### ✅ Faire

```go
// Réutiliser le même builder
builder := rete.NewAlphaChainBuilderWithMetrics(network, storage, network.ChainMetrics)
for i := range rules {
    builder.BuildChain(...)
}
```

## 📚 Documentation Complète

- **Guide Détaillé**: `docs/CHAIN_PERFORMANCE_OPTIMIZATION.md`
- **Code Source**: `rete/chain_metrics.go`
- **Tests**: `rete/chain_performance_test.go`
- **Benchmarks**: `rete/chain_metrics_test.go`

## 🤝 Support

Pour des questions ou des problèmes:
1. Consulter la documentation détaillée
2. Vérifier les tests d'exemple
3. Ouvrir une issue sur GitHub

## 📄 License

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License