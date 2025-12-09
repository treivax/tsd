# Module RETE - Moteur de Règles

Moteur d'inférence basé sur l'algorithme RETE avec optimisations avancées et persistance.

## 🎯 Vue d'ensemble

Le module RETE implémente un réseau d'inférence qui :

- **Construit automatiquement** un réseau de nœuds à partir de règles TSD
- **Exécute efficacement** les règles sur les faits
- **Optimise le partage** de nœuds entre règles (alpha chains, beta sharing)
- **Persiste l'état** dans un storage configurable
- **Fournit des métriques** détaillées de performance

## 🏗️ Architecture

```
Programme TSD → AST → Réseau RETE → Exécution
                         ↓
                    Storage (persistance)
```

### Types de Nœuds

| Nœud | Description |
|------|-------------|
| **RootNode** | Point d'entrée pour tous les faits |
| **TypeNode** | Filtre et valide les faits par type |
| **AlphaNode** | Évalue les conditions sur faits individuels |
| **BetaNode** | Gère les jointures multi-faits |
| **JoinNode** | Effectue les jointures conditionnelles |
| **AccumulateNode** | Agrégations (count, sum, avg, etc.) |
| **TerminalNode** | Déclenche les actions |

### Optimisations

- ✅ **Alpha Chains** : Partage des nœuds alpha (40-60% réduction)
- ✅ **Beta Sharing** : Partage des jointures (30-50% réduction)
- ✅ **Node Lifecycle** : Gestion optimisée du cycle de vie
- ✅ **Caches LRU** : Cache des résultats de filtrage
- ✅ **Arithmetic Cache** : Cache des calculs arithmétiques
- ✅ **Normalization** : Détection d'équivalences

## 🚀 Utilisation Rapide

```go
package main

import (
    "github.com/yourusername/tsd/rete"
    "github.com/yourusername/tsd/constraint"
)

func main() {
    // 1. Créer le storage
    storage := rete.NewMemoryStorage()
    
    // 2. Créer le réseau
    network := rete.NewReteNetwork(storage)
    
    // 3. Charger les règles
    pipeline := constraint.NewConstraintPipeline()
    network, err := pipeline.IngestFile("rules.tsd", network, storage)
    if err != nil {
        panic(err)
    }
    
    // 4. Asserter des faits
    fact := map[string]interface{}{
        "type": "Person",
        "name": "Alice",
        "age":  25,
    }
    network.Assert(fact)
    
    // 5. Consulter les métriques
    metrics := network.GetMetrics()
    fmt.Printf("Rules: %d, Facts: %d, Activations: %d\n",
        metrics.RuleCount, metrics.FactCount, metrics.ActivationCount)
}
```

## 📖 Documentation

### Guides Utilisateur

- [**Actions**](docs/ACTIONS.md) - Guide des actions (print, assert, retract)
- [**Alpha Chains**](docs/ALPHA_CHAINS.md) - Partage des nœuds alpha
- [**Beta Sharing**](docs/BETA_SHARING.md) - Partage des jointures
- [**Beta Chains**](docs/BETA_CHAINS.md) - Chaînes de jointures optimisées
- [**Nested OR**](docs/NESTED_OR.md) - Expressions OR imbriquées
- [**Arithmetic**](docs/ARITHMETIC.md) - Expressions arithmétiques
- [**Multi-Source Aggregation**](docs/MULTI_SOURCE_AGGREGATION.md) - Agrégations avancées

### Guides Techniques

- [**Network Architecture**](docs/NETWORK_ARCHITECTURE.md) - Architecture modulaire du réseau RETE
- [**Advanced Nodes**](docs/ADVANCED_NODES_IMPLEMENTATION.md) - Implémentation des nœuds avancés
- [**Advanced Nodes Usage**](docs/ADVANCED_NODES_USAGE_GUIDE.md) - Guide d'utilisation des nœuds avancés
- [**Node Lifecycle**](docs/NODE_LIFECYCLE.md) - Cycle de vie des nœuds
- [**Normalization**](docs/NORMALIZATION.md) - Normalisation des expressions
- [**Optimizations**](docs/OPTIMIZATIONS.md) - Guide des optimisations
- [**Testing**](docs/TESTING.md) - Guide des tests
- [**Tuple Space**](docs/TUPLE_SPACE_IMPLEMENTATION.md) - Implémentation de l'espace de tuples

### Guides de Conception

- [**Alpha/Beta Chains Comparison**](docs/ALPHA_BETA_CHAINS_COMPARISON.md) - Comparaison alpha vs beta chains
- [**Beta Chains Design**](docs/BETA_CHAINS_DESIGN.md) - Conception des beta chains
- [**Beta Sharing Design**](docs/BETA_SHARING_DESIGN.md) - Conception du beta sharing
- [**Beta Nodes Architecture**](docs/BETA_NODES_ARCHITECTURE_DIAGRAMS.md) - Diagrammes d'architecture
- [**Beta Nodes Guide**](docs/BETA_NODES_GUIDE.md) - Guide des nœuds beta
- [**Expression Analyzer**](docs/EXPRESSION_ANALYZER_README.md) - Analyseur d'expressions
- [**Feature: Arithmetic Alpha Nodes**](docs/FEATURE_ARITHMETIC_ALPHA_NODES.md) - Alpha nodes arithmétiques
- [**Feature: Passthrough Per Rule**](docs/FEATURE_PASSTHROUGH_PER_RULE.md) - Optimisation passthrough

### Exemples

- [**Beta Chains Examples**](docs/BETA_CHAINS_EXAMPLES.md) - Exemples de beta chains
- [**Beta Sharing Examples**](docs/BETA_SHARING_EXAMPLES.md) - Exemples de beta sharing

## 🔧 Configuration

### Configuration par Défaut

```go
network := rete.NewReteNetwork(storage)
// Toutes les optimisations activées par défaut
```

### Configuration Personnalisée

```go
config := rete.NetworkConfig{
    // Alpha Chains
    EnableAlphaChains:   true,
    AlphaChainMinLength: 2,
    AlphaChainMaxDepth:  10,
    
    // Beta Sharing
    EnableBetaSharing:   true,
    BetaSharingStrategy: "aggressive",
    
    // Caches
    AlphaCacheEnabled:   true,
    AlphaCacheSize:      10000,
    AlphaCacheTTL:       5 * time.Minute,
    
    // Node Lifecycle
    EnableLazyNodeCreation: true,
    EnableNodeGC:           true,
}

network := rete.NewReteNetworkWithConfig(storage, &config)
```

## 📊 Métriques et Monitoring

```go
// Métriques générales
metrics := network.GetMetrics()
fmt.Printf("Rules: %d\n", metrics.RuleCount)
fmt.Printf("Facts: %d\n", metrics.FactCount)
fmt.Printf("Alpha nodes: %d\n", metrics.AlphaNodeCount)
fmt.Printf("Beta nodes: %d\n", metrics.BetaNodeCount)

// Métriques alpha chains
alphaMetrics := network.GetAlphaChainMetrics()
fmt.Printf("Alpha sharing rate: %.2f%%\n", alphaMetrics.SharingRate)

// Métriques beta sharing
betaMetrics := network.GetBetaSharingMetrics()
fmt.Printf("Beta sharing rate: %.2f%%\n", betaMetrics.SharingRate)

// Export Prometheus
exporter := rete.NewPrometheusExporter(network)
http.Handle("/metrics", promhttp.Handler())
```

## 🧪 Tests

```bash
# Tous les tests
go test ./rete/...

# Tests spécifiques
go test ./rete/alpha_chain_integration_test.go
go test ./rete/beta_sharing_integration_test.go

# Tests avec coverage
go test ./rete/... -cover

# Benchmarks
go test ./rete/... -bench=. -benchmem
```

## 🔗 Liens Utiles

### Documentation Générale
- [README Principal](../README.md)
- [Documentation Complète](../docs/README.md)
- [Tutorial](../docs/TUTORIAL.md)
- [Features](../docs/FEATURES.md)
- [Optimizations Guide](../docs/OPTIMIZATIONS.md)

### Package Constraint
- [Constraint Parser](../constraint/README.md)
- [Grammar Guide](../constraint/docs/GRAMMAR_COMPLETE.md)

## 📈 Performance

### Benchmarks Typiques

| Réseau | Sans Optimisations | Avec Optimisations | Amélioration |
|--------|--------------------|--------------------|--------------|
| Petit (10 règles) | 45ms | 15ms | **3.0x** |
| Moyen (100 règles) | 780ms | 165ms | **4.7x** |
| Grand (1000 règles) | 12.5s | 2.1s | **6.0x** |

### Mémoire

| Réseau | Sans Optimisations | Avec Optimisations | Réduction |
|--------|--------------------|--------------------|-----------| 
| Petit | 12 MB | 5 MB | **58%** |
| Moyen | 145 MB | 52 MB | **64%** |
| Grand | 2.8 GB | 0.9 GB | **68%** |

## 🤝 Contribution

Pour contribuer au module RETE :

1. Lire [Development Guidelines](../docs/development_guidelines.md)
2. Créer une branche pour votre feature
3. Ajouter des tests
4. Soumettre une Pull Request

## 📄 License

Voir [LICENSE](../LICENSE) à la racine du projet.

---

**Version** : 1.0  
**Dernière mise à jour** : Janvier 2025