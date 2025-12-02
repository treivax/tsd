# Métriques de Décomposition Arithmétique

Ce document décrit le système de métriques pour la décomposition et l'évaluation des expressions arithmétiques dans le réseau RETE.

## Vue d'ensemble

Le système de métriques collecte des statistiques détaillées sur :
- Les performances d'évaluation des expressions arithmétiques
- L'efficacité du cache de résultats intermédiaires
- La structure et la complexité des chaînes décomposées
- La détection de dépendances circulaires
- Les activations de règles arithmétiques

## Architecture

### Composants principaux

```
ArithmeticDecompositionMetrics
├── RuleArithmeticMetrics (par règle)
│   ├── Compteurs d'activations
│   ├── Compteurs d'évaluations
│   ├── Statistiques de temps
│   ├── Histogrammes
│   ├── Statistiques de cache
│   └── Informations de dépendances
└── GlobalArithmeticMetrics (agrégées)
    ├── Totaux globaux
    ├── Moyennes
    ├── Percentiles
    └── Statistiques de cache global
```

## Configuration

### Configuration par défaut

```go
config := DefaultMetricsConfig()
// Enabled: true
// CollectHistograms: true
// CollectPercentiles: true
// HistogramBuckets: [1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000] µs
// MaxRulesToTrack: 1000
// RetentionDuration: 24h
// AggregationInterval: 1min
```

### Configuration personnalisée

```go
config := MetricsConfig{
    Enabled:             true,
    CollectHistograms:   true,
    CollectPercentiles:  true,
    HistogramBuckets:    []int{10, 50, 100, 500, 1000, 5000},
    MaxRulesToTrack:     500,
    RetentionDuration:   12 * time.Hour,
    AggregationInterval: 30 * time.Second,
}
metrics := NewArithmeticDecompositionMetrics(config)
```

## Utilisation

### Création et initialisation

```go
// Créer les métriques
metrics := NewArithmeticDecompositionMetrics(DefaultMetricsConfig())

// Désactiver si nécessaire
config := DefaultMetricsConfig()
config.Enabled = false
metrics := NewArithmeticDecompositionMetrics(config)
```

### Enregistrement des métriques

#### 1. Structure de chaîne

Enregistrer la structure d'une chaîne arithmétique décomposée :

```go
ruleID := "price_calculation"
chainLength := 5
atomicSteps := 3
comparisonSteps := 2
intermediateResults := []string{"subtotal", "tax", "total"}
dependencies := map[string][]string{
    "subtotal": {},
    "tax":      {"subtotal"},
    "total":    {"subtotal", "tax"},
}

metrics.RecordChainStructure(
    ruleID,
    chainLength,
    atomicSteps,
    comparisonSteps,
    intermediateResults,
    dependencies,
)
```

#### 2. Évaluations

Enregistrer une évaluation d'expression arithmétique :

```go
start := time.Now()
// ... évaluer l'expression ...
duration := time.Since(start)

success := true
metrics.RecordEvaluation(ruleID, success, duration)
```

#### 3. Activations

Enregistrer une activation de règle :

```go
start := time.Now()
// ... activer la règle ...
duration := time.Since(start)

success := true
metrics.RecordActivation(ruleID, success, duration)
```

#### 4. Cache

Enregistrer les interactions avec le cache :

```go
// Cache hit
result := cache.Get("intermediate_result")
if result != nil {
    metrics.RecordCacheHit(ruleID)
} else {
    metrics.RecordCacheMiss(ruleID)
    // ... calculer et stocker ...
}
```

Mettre à jour les statistiques globales du cache :

```go
stats := cache.Stats()
metrics.UpdateCacheStatistics(
    stats.Size,
    stats.Evictions,
    stats.EstimatedMemoryUsage,
)
```

#### 5. Détection de cycles

Enregistrer la détection de dépendances circulaires :

```go
detector := NewCircularDependencyDetector()
// ... ajouter des dépendances ...

result := detector.Validate()
metrics.RecordGraphValidation(result.MaxDepth, result.HasCycles)

if result.HasCycles {
    for _, cycle := range result.Cycles {
        metrics.RecordCircularDependency(ruleID, cycle)
    }
}
```

### Consultation des métriques

#### Métriques d'une règle spécifique

```go
rule := metrics.GetRuleMetrics(ruleID)
if rule != nil {
    fmt.Printf("Evaluations: %d\n", rule.TotalEvaluations)
    fmt.Printf("Average time: %v\n", rule.AvgEvaluationTime)
    fmt.Printf("Cache hit rate: %.2f%%\n", rule.CacheHitRate*100)
    fmt.Printf("Chain length: %d\n", rule.ChainLength)
}
```

#### Métriques globales

```go
global := metrics.GetGlobalMetrics()
fmt.Printf("Total rules: %d\n", global.TotalRulesWithArithmetic)
fmt.Printf("Total evaluations: %d\n", global.TotalEvaluations)
fmt.Printf("Average evaluation time: %v\n", global.AverageEvaluationTime)
fmt.Printf("P95 time: %v\n", global.EvaluationTimeP95)
fmt.Printf("P99 time: %v\n", global.EvaluationTimeP99)
fmt.Printf("Cache global hit rate: %.2f%%\n", global.CacheGlobalHitRate*100)
```

#### Toutes les métriques

```go
allMetrics := metrics.GetAllRuleMetrics()
for ruleID, rule := range allMetrics {
    fmt.Printf("Rule %s: %d evaluations\n", ruleID, rule.TotalEvaluations)
}
```

### Analyses et classements

#### Règles les plus évaluées

```go
topRules := metrics.GetTopRulesByEvaluations(10)
for i, rule := range topRules {
    fmt.Printf("%d. %s: %d evaluations\n", 
        i+1, rule.RuleID, rule.TotalEvaluations)
}
```

#### Règles les plus coûteuses (temps total)

```go
topByDuration := metrics.GetTopRulesByDuration(10)
for i, rule := range topByDuration {
    fmt.Printf("%d. %s: %v total\n", 
        i+1, rule.RuleID, rule.TotalEvaluationTime)
}
```

#### Règles les plus lentes (temps moyen)

```go
slowest := metrics.GetSlowestRules(10)
for i, rule := range slowest {
    fmt.Printf("%d. %s: %v avg\n", 
        i+1, rule.RuleID, rule.AvgEvaluationTime)
}
```

### Résumé

Obtenir un résumé formaté de toutes les métriques :

```go
summary := metrics.GetSummary()

// Structure du résumé :
// {
//   "rules": {
//     "total_with_arithmetic": 10,
//     "total_decomposed_chains": 15,
//     "tracked_rules": 10
//   },
//   "nodes": {
//     "total_atomic": 45,
//     "total_comparison": 20,
//     "average_chain_length": 3.5
//   },
//   "evaluations": {
//     "total": 1000,
//     "total_time": "150ms",
//     "average_time": "150µs",
//     "min_time": "10µs",
//     "max_time": "5ms"
//   },
//   "cache": {
//     "hits": 750,
//     "misses": 250,
//     "hit_rate": 0.75,
//     "size": 100,
//     "evictions": 5,
//     "memory_usage_bytes": 10240
//   },
//   "validation": {
//     "total_validations": 15,
//     "cycles_detected": 0,
//     "max_graph_depth": 5
//   }
// }
```

## Métriques collectées

### Par règle (RuleArithmeticMetrics)

| Métrique | Type | Description |
|----------|------|-------------|
| `RuleID` | string | Identifiant unique de la règle |
| `RuleName` | string | Nom de la règle |
| `TotalActivations` | int64 | Nombre total d'activations |
| `SuccessfulActivations` | int64 | Activations réussies |
| `FailedActivations` | int64 | Activations échouées |
| `TotalEvaluations` | int64 | Nombre total d'évaluations |
| `SuccessfulEvaluations` | int64 | Évaluations réussies |
| `FailedEvaluations` | int64 | Évaluations échouées |
| `ChainLength` | int | Longueur de la chaîne décomposée |
| `AtomicStepsCount` | int | Nombre d'étapes atomiques |
| `ComparisonStepsCount` | int | Nombre de comparaisons |
| `IntermediateResults` | []string | Liste des résultats intermédiaires |
| `Dependencies` | map[string][]string | Graphe de dépendances |
| `TotalEvaluationTime` | time.Duration | Temps total d'évaluation |
| `MinEvaluationTime` | time.Duration | Temps minimum |
| `MaxEvaluationTime` | time.Duration | Temps maximum |
| `AvgEvaluationTime` | time.Duration | Temps moyen |
| `EvaluationTimeHistogram` | map[int]int64 | Histogramme des temps (buckets → count) |
| `CacheHits` | int64 | Nombre de cache hits |
| `CacheMisses` | int64 | Nombre de cache misses |
| `CacheHitRate` | float64 | Taux de cache hit (0.0-1.0) |
| `CacheEnabled` | bool | Cache activé pour cette règle |
| `MaxDependencyDepth` | int | Profondeur maximale des dépendances |
| `HasCircularDeps` | bool | Présence de dépendances circulaires |
| `IsolatedNodes` | []string | Nœuds isolés dans le graphe |
| `FirstSeen` | time.Time | Première occurrence |
| `LastSeen` | time.Time | Dernière occurrence |
| `Metadata` | map[string]interface{} | Métadonnées supplémentaires |

### Globales (GlobalArithmeticMetrics)

| Métrique | Type | Description |
|----------|------|-------------|
| `TotalRulesWithArithmetic` | int | Nombre de règles avec arithmétique |
| `TotalDecomposedChains` | int | Nombre de chaînes décomposées |
| `TotalAtomicNodes` | int | Nombre total de nœuds atomiques |
| `TotalComparisonNodes` | int | Nombre total de nœuds de comparaison |
| `AverageChainLength` | float64 | Longueur moyenne des chaînes |
| `AverageAtomicStepsPerChain` | float64 | Étapes atomiques moyennes par chaîne |
| `AverageDependencyDepth` | float64 | Profondeur moyenne des dépendances |
| `SharedNodesRatio` | float64 | Ratio de nœuds partagés |
| `CacheGlobalHitRate` | float64 | Taux global de cache hit |
| `TotalActivations` | int64 | Activations totales |
| `TotalEvaluations` | int64 | Évaluations totales |
| `TotalCacheHits` | int64 | Cache hits totaux |
| `TotalCacheMisses` | int64 | Cache misses totaux |
| `TotalCircularDepsDetected` | int64 | Dépendances circulaires détectées |
| `TotalEvaluationTime` | time.Duration | Temps total d'évaluation |
| `AverageEvaluationTime` | time.Duration | Temps moyen d'évaluation |
| `MinEvaluationTime` | time.Duration | Temps minimum global |
| `MaxEvaluationTime` | time.Duration | Temps maximum global |
| `EvaluationTimeP50` | time.Duration | Percentile 50 (médiane) |
| `EvaluationTimeP95` | time.Duration | Percentile 95 |
| `EvaluationTimeP99` | time.Duration | Percentile 99 |
| `CacheSize` | int | Taille actuelle du cache |
| `CacheEvictions` | int64 | Évictions du cache |
| `CacheMemoryUsage` | int64 | Utilisation mémoire du cache (bytes) |
| `GraphValidations` | int64 | Validations de graphe effectuées |
| `CyclesDetected` | int64 | Cycles détectés |
| `MaxGraphDepth` | int | Profondeur maximale de graphe |

## Histogrammes

Les histogrammes collectent la distribution des temps d'évaluation dans des buckets prédéfinis.

### Buckets par défaut (en microsecondes)

```
1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000
```

### Interprétation

```go
rule := metrics.GetRuleMetrics(ruleID)
for bucket, count := range rule.EvaluationTimeHistogram {
    fmt.Printf("≤ %dµs: %d evaluations\n", bucket, count)
}

// Exemple de sortie :
// ≤ 10µs: 150 evaluations
// ≤ 50µs: 500 evaluations
// ≤ 100µs: 300 evaluations
// ≤ 500µs: 50 evaluations
```

## Percentiles

Les percentiles donnent une vue de la distribution des performances :

- **P50 (médiane)** : 50% des évaluations sont plus rapides
- **P95** : 95% des évaluations sont plus rapides
- **P99** : 99% des évaluations sont plus rapides

```go
global := metrics.GetGlobalMetrics()
fmt.Printf("P50: %v\n", global.EvaluationTimeP50)
fmt.Printf("P95: %v\n", global.EvaluationTimeP95)
fmt.Printf("P99: %v\n", global.EvaluationTimeP99)
```

## Intégration avec le cache

Les métriques s'intègrent naturellement avec le cache de résultats arithmétiques :

```go
cache := NewArithmeticResultCache(CacheConfig{
    MaxSize: 100,
    TTL:     5 * time.Minute,
    Enabled: true,
})

metrics := NewArithmeticDecompositionMetrics(DefaultMetricsConfig())

// Lors de l'évaluation
result := cache.Get("intermediate_result")
if result != nil {
    metrics.RecordCacheHit(ruleID)
    // Utiliser le résultat du cache
} else {
    metrics.RecordCacheMiss(ruleID)
    // Calculer et stocker
    result = computeResult()
    cache.Set("intermediate_result", result)
}

// Périodiquement, mettre à jour les stats du cache
stats := cache.Stats()
metrics.UpdateCacheStatistics(stats.Size, stats.Evictions, stats.EstimatedMemoryUsage)
```

## Intégration avec le détecteur de cycles

```go
detector := NewCircularDependencyDetector()
metrics := NewArithmeticDecompositionMetrics(DefaultMetricsConfig())

// Construire le graphe de dépendances
for node, deps := range dependencies {
    for _, dep := range deps {
        detector.AddDependency(node, dep)
    }
}

// Valider
result := detector.Validate()
metrics.RecordGraphValidation(result.MaxDepth, result.HasCycles)

if result.HasCycles {
    for _, cycle := range result.Cycles {
        metrics.RecordCircularDependency(ruleID, cycle)
    }
}
```

## Export Prometheus

Les métriques sont conçues pour être facilement exportées vers Prometheus. Structure recommandée :

### Compteurs

```prometheus
# Activations
rete_arithmetic_activations_total{rule_id="rule1",status="success"} 100
rete_arithmetic_activations_total{rule_id="rule1",status="failed"} 5

# Évaluations
rete_arithmetic_evaluations_total{rule_id="rule1",status="success"} 500
rete_arithmetic_evaluations_total{rule_id="rule1",status="failed"} 10

# Cache
rete_arithmetic_cache_hits_total{rule_id="rule1"} 400
rete_arithmetic_cache_misses_total{rule_id="rule1"} 100
```

### Histogrammes

```prometheus
rete_arithmetic_evaluation_duration_microseconds_bucket{rule_id="rule1",le="10"} 50
rete_arithmetic_evaluation_duration_microseconds_bucket{rule_id="rule1",le="50"} 200
rete_arithmetic_evaluation_duration_microseconds_bucket{rule_id="rule1",le="100"} 400
rete_arithmetic_evaluation_duration_microseconds_bucket{rule_id="rule1",le="+Inf"} 500
```

### Jauges

```prometheus
rete_arithmetic_cache_size 85
rete_arithmetic_cache_memory_bytes 10485760
rete_arithmetic_chain_length{rule_id="rule1"} 5
rete_arithmetic_dependency_depth{rule_id="rule1"} 3
```

### Résumés

```prometheus
rete_arithmetic_evaluation_duration_microseconds{quantile="0.5"} 45
rete_arithmetic_evaluation_duration_microseconds{quantile="0.95"} 180
rete_arithmetic_evaluation_duration_microseconds{quantile="0.99"} 450
```

## Exemples de requêtes Prometheus

### Taux de cache hit global

```promql
sum(rate(rete_arithmetic_cache_hits_total[5m])) / 
(sum(rate(rete_arithmetic_cache_hits_total[5m])) + 
 sum(rate(rete_arithmetic_cache_misses_total[5m])))
```

### Règles les plus lentes (P95)

```promql
topk(10, 
  histogram_quantile(0.95, 
    rate(rete_arithmetic_evaluation_duration_microseconds_bucket[5m])
  )
)
```

### Règles les plus sollicitées

```promql
topk(10, sum(rate(rete_arithmetic_evaluations_total[5m])) by (rule_id))
```

### Détection d'anomalies (P99 > 1ms)

```promql
histogram_quantile(0.99, 
  rate(rete_arithmetic_evaluation_duration_microseconds_bucket[5m])
) > 1000
```

## Bonnes pratiques

### 1. Activation sélective

N'activez les métriques que lorsque nécessaire :

```go
config := DefaultMetricsConfig()
config.Enabled = os.Getenv("ENABLE_METRICS") == "true"
```

### 2. Limitation du nombre de règles

Configurez `MaxRulesToTrack` en fonction de vos besoins :

```go
config.MaxRulesToTrack = 100 // Pour les environnements contraints
```

### 3. Histogrammes adaptatifs

Ajustez les buckets selon vos cas d'usage :

```go
// Pour des évaluations très rapides
config.HistogramBuckets = []int{1, 2, 5, 10, 20, 50, 100}

// Pour des évaluations plus lentes
config.HistogramBuckets = []int{100, 500, 1000, 5000, 10000, 50000}
```

### 4. Nettoyage périodique

Réinitialisez les métriques périodiquement si nécessaire :

```go
// Tous les jours à minuit
ticker := time.NewTicker(24 * time.Hour)
go func() {
    for range ticker.C {
        metrics.Reset()
    }
}()
```

### 5. Export asynchrone

Exportez les métriques dans une goroutine séparée :

```go
go func() {
    ticker := time.NewTicker(1 * time.Minute)
    for range ticker.C {
        summary := metrics.GetSummary()
        exportToPrometheus(summary)
    }
}()
```

## Performance

### Overhead

- **Métriques désactivées** : aucun overhead (early return)
- **Métriques activées** : ~1-2µs par enregistrement
- **Avec histogrammes** : +0.5µs supplémentaire
- **Avec percentiles** : calcul uniquement à la lecture

### Optimisations

- Verrouillage optimisé avec `RWMutex`
- Copies à la demande pour éviter les modifications externes
- Calculs paresseux (percentiles)
- Éviction LRU des règles les moins récentes

### Benchmarks

```
BenchmarkRecordEvaluation-8         1000000    1.2 µs/op
BenchmarkGetRuleMetrics-8          5000000    0.3 µs/op
BenchmarkConcurrentRecording-8      500000    2.5 µs/op
```

## Thread-safety

Toutes les opérations sont thread-safe grâce à l'utilisation de `sync.RWMutex` :

```go
// Sûr à appeler depuis plusieurs goroutines
go metrics.RecordEvaluation("rule1", true, 100*time.Microsecond)
go metrics.RecordEvaluation("rule2", true, 200*time.Microsecond)
go metrics.RecordCacheHit("rule1")
```

## Tests

Exécuter les tests :

```bash
# Tests unitaires
go test -v -run TestArithmeticDecompositionMetrics

# Tests d'intégration
go test -v -run TestMetricsIntegration

# Tests de concurrence
go test -v -run TestConcurrentMetrics

# Benchmarks
go test -bench=. -benchmem

# Couverture
go test -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Troubleshooting

### Métriques non collectées

Vérifiez que les métriques sont activées :

```go
if !metrics.config.Enabled {
    log.Println("Metrics are disabled")
}
```

### Limite de règles atteinte

Augmentez `MaxRulesToTrack` ou vérifiez les évictions :

```go
allMetrics := metrics.GetAllRuleMetrics()
if len(allMetrics) >= config.MaxRulesToTrack {
    log.Printf("Metric limit reached: %d rules", len(allMetrics))
}
```

### Performances dégradées

Désactivez les histogrammes ou les percentiles :

```go
config.CollectHistograms = false
config.CollectPercentiles = false
```

### Mémoire élevée

Réduisez la rétention ou le nombre de règles :

```go
config.MaxRulesToTrack = 100
config.RetentionDuration = 1 * time.Hour
```

## Voir aussi

- [Cache de résultats arithmétiques](ARITHMETIC_CACHE.md)
- [Détecteur de dépendances circulaires](CIRCULAR_DEPENDENCY_DETECTOR.md)
- [Décomposition arithmétique](ARITHMETIC_DECOMPOSITION.md)
- [Export Prometheus](PROMETHEUS_EXPORT.md)