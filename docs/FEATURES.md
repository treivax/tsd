# Fonctionnalités TSD

Guide complet des fonctionnalités du moteur de règles TSD basé sur l'algorithme RETE.

## Table des matières

- [Syntaxe et Langage](#syntaxe-et-langage)
- [Moteur RETE](#moteur-rete)
- [Actions et Exécution](#actions-et-exécution)
- [Agrégations](#agrégations)
- [Optimisations](#optimisations)
- [Strong Mode - Cohérence](#strong-mode---cohérence)
- [Monitoring et Métriques](#monitoring-et-métriques)
- [Ingestion Incrémentale](#ingestion-incrémentale)

---

## Syntaxe et Langage

### Format de Fichier Unifié (.tsd)

TSD utilise un format de fichier unique pour définir types, faits, règles et actions.

```tsd
# Déclaration de types
type Person : <name: string, age: number, city: string>

# Faits
Person(name="Alice", age=30, city="Paris")
Person(name="Bob", age=25, city="Lyon")

# Règles
rule "adults_in_paris" {
    when {
        p: Person(age >= 18, city == "Paris")
    }
    then {
        print("Adult found: " + p.name)
    }
}
```

### Opérateurs Supportés

**Comparaisons** : `==`, `!=`, `>`, `<`, `>=`, `<=`

**Logique** : `AND`, `OR`, `NOT()`

**Patterns** : `CONTAINS`, `LIKE`, `MATCHES`, `IN`

**Arithmétique** : `+`, `-`, `*`, `/`, `%`

**Fonctions** : `LENGTH()`, `ABS()`, `UPPER()`, `LOWER()`, `ROUND()`, `FLOOR()`, `CEIL()`

### Expressions Complexes

```tsd
rule "complex_conditions" {
    when {
        p: Person(
            age >= 18,
            NOT(city == "Paris" OR city == "Lyon"),
            LENGTH(name) > 3
        )
    }
    then {
        print("Match: " + p.name)
    }
}
```

### Identifiants de Règles

Chaque règle peut avoir un identifiant unique pour faciliter le suivi :

```tsd
rule "check_adult" id="R001" {
    when { p: Person(age >= 18) }
    then { print("Adult: " + p.name) }
}
```

**Avantages** :
- Traçabilité dans les logs
- Débogage facilité
- Metrics par règle

---

## Moteur RETE

### Architecture du Réseau

Le moteur utilise l'algorithme RETE optimisé avec :

- **Alpha Nodes** : Filtrage des faits individuels
- **Beta Nodes** : Jointures entre faits
- **Type Nodes** : Organisation par type
- **Terminal Nodes** : Activation des règles

### Alpha Chains - Optimisation

Les alpha chains réduisent la duplication en partageant les nœuds de filtrage communs.

**Exemple** :
```tsd
# Ces deux règles partagent le même alpha node pour age >= 18
rule "adults" {
    when { p: Person(age >= 18) }
    then { print("Adult") }
}

rule "adults_paris" {
    when { p: Person(age >= 18, city == "Paris") }
    then { print("Adult in Paris") }
}
```

**Gains** : Jusqu'à 60% de réduction du nombre de nœuds

### Beta Sharing - Partage de Jointures

Les jointures identiques sont partagées entre règles :

```tsd
rule "R1" {
    when {
        p: Person(age > 18)
        c: Company(name == p.employer)
    }
    then { print("R1 activated") }
}

rule "R2" {
    when {
        p: Person(age > 18)
        c: Company(name == p.employer)
    }
    then { print("R2 activated") }
}
# La jointure Person-Company est partagée
```

**Gains** : Jusqu'à 40% de réduction du temps d'exécution

### Node Lifecycle Management

Gestion automatique du cycle de vie des nœuds :

- **Création** : Construction à la demande
- **Activation** : Propagation des tokens
- **Désactivation** : Retrait des règles
- **Nettoyage** : Garbage collection

---

## Actions et Exécution

### Actions Prédéfinies

**print** : Affichage de messages
```tsd
then { print("Message: " + value) }
```

**assert** : Ajout de nouveaux faits
```tsd
then { assert(Adult(name=p.name, age=p.age)) }
```

**retract** : Retrait de faits
```tsd
then { retract(p) }
```

### Actions Multiples

Une règle peut exécuter plusieurs actions :

```tsd
rule "process_person" {
    when { p: Person(age >= 18) }
    then {
        print("Found adult: " + p.name)
        assert(Adult(name=p.name, age=p.age))
        retract(p)
    }
}
```

### Expressions Arithmétiques dans les Actions

```tsd
action CalculateTotal : <total: number>

rule "calculate" {
    when {
        o1: Order(amount > 0)
        o2: Order(amount > 0)
    }
    then {
        assert(CalculateTotal(total=o1.amount + o2.amount))
    }
}
```

**Support complet** :
- Variables arithmétiques
- Opérations complexes : `(a + b) * c / d`
- Fonctions : `ABS(x)`, `ROUND(y)`
- Expressions imbriquées

---

## Agrégations

### Agrégations Simples

```tsd
rule "count_adults" {
    when {
        accumulate(
            p: Person(age >= 18),
            $count: count(p)
        )
    }
    then {
        print("Total adults: " + $count)
    }
}
```

**Fonctions d'agrégation** : `count()`, `sum()`, `avg()`, `min()`, `max()`

### Agrégations Multi-Sources

Agrégation sur plusieurs types de faits :

```tsd
rule "total_revenue" {
    when {
        accumulate(
            o: Order(status == "paid"),
            s: Sale(status == "completed"),
            $total: sum(o.amount) + sum(s.amount)
        )
    }
    then {
        print("Total revenue: " + $total)
    }
}
```

### Agrégations avec Seuils

```tsd
rule "alert_high_count" {
    when {
        accumulate(
            e: Event(severity == "critical"),
            $count: count(e)
        )
        eval($count > 10)
    }
    then {
        print("ALERT: Too many critical events!")
    }
}
```

### Agrégations avec Jointures

Combiner agrégations et jointures :

```tsd
rule "department_stats" {
    when {
        d: Department()
        accumulate(
            e: Employee(department == d.name),
            $count: count(e),
            $avgSalary: avg(e.salary)
        )
    }
    then {
        print("Dept " + d.name + ": " + $count + " employees, avg salary: " + $avgSalary)
    }
}
```

---

## Optimisations

### Optimisations par Défaut

Activées automatiquement :

- **Type Node Sharing** : Partage des nœuds de type
- **Alpha Sharing** : Partage des filtres alpha
- **Beta Sharing** : Partage des jointures
- **Normalization Cache** : Cache des expressions normalisées
- **Arithmetic Result Cache** : Cache des résultats arithmétiques

### Optimisations Avancées

#### LRU Cache pour Alpha Sharing

Cache LRU pour les résultats des alpha nodes :

```go
config := rete.NewNetworkConfig()
config.AlphaCacheSize = 1000
config.AlphaCacheTTL = 5 * time.Minute

network := rete.NewReteNetworkWithConfig(storage, config)
```

**Gains** : 30-50% de réduction du temps de filtrage

#### Passthrough Nodes

Optimisation pour les règles sans conditions :

```tsd
rule "always_trigger" {
    when { p: Person() }  # Pas de conditions
    then { print("Person: " + p.name) }
}
```

Le moteur crée un passthrough node au lieu d'un alpha node complet.

### Performance Monitoring

Suivi des performances en temps réel :

```go
metrics := network.GetPerformanceMetrics()
fmt.Printf("Alpha nodes: %d\n", metrics.AlphaNodeCount)
fmt.Printf("Beta nodes: %d\n", metrics.BetaNodeCount)
fmt.Printf("Sharing rate: %.2f%%\n", metrics.SharingRate)
fmt.Printf("Avg rule time: %v\n", metrics.AvgRuleExecutionTime)
```

---

## Strong Mode - Cohérence

Le Strong Mode garantit la cohérence des données entre le réseau RETE et le storage.

### Activation

```go
storage := rete.NewMemoryStorage()
network := rete.NewReteNetwork(storage)

// Activer le Strong Mode avec configuration par défaut
network.EnableStrongMode()
```

### Configuration

```go
config := rete.StrongModeConfig{
    MaxRetries:        5,
    InitialBackoff:    10 * time.Millisecond,
    MaxBackoff:        500 * time.Millisecond,
    Timeout:           30 * time.Second,
    EnableMetrics:     true,
}

network.EnableStrongModeWithConfig(config)
```

### Garanties

- **Read-after-Write** : Les faits sont immédiatement lisibles après insertion
- **Atomicité** : Toutes les opérations sont atomiques
- **Durabilité** : Les données sont persistées avant confirmation
- **Cohérence** : Le réseau et le storage sont toujours synchronisés

### Configurations Optimisées

#### PostgreSQL / MySQL
```go
config := rete.StrongModeConfig{
    MaxRetries:     3,
    InitialBackoff: 5 * time.Millisecond,
    MaxBackoff:     100 * time.Millisecond,
    Timeout:        10 * time.Second,
}
```

#### Redis
```go
config := rete.StrongModeConfig{
    MaxRetries:     2,
    InitialBackoff: 2 * time.Millisecond,
    MaxBackoff:     50 * time.Millisecond,
    Timeout:        5 * time.Second,
}
```

#### Cassandra / DynamoDB
```go
config := rete.StrongModeConfig{
    MaxRetries:     10,
    InitialBackoff: 50 * time.Millisecond,
    MaxBackoff:     2 * time.Second,
    Timeout:        60 * time.Second,
}
```

### Métriques

```go
metrics := network.GetStrongModeMetrics()
fmt.Printf("Total operations: %d\n", metrics.TotalOps)
fmt.Printf("Retries: %d\n", metrics.TotalRetries)
fmt.Printf("Failures: %d\n", metrics.Failures)
fmt.Printf("Avg latency: %v\n", metrics.AvgLatency)
fmt.Printf("Success rate: %.2f%%\n", metrics.SuccessRate)
```

---

## Monitoring et Métriques

### Intégration Prometheus

Export automatique des métriques vers Prometheus :

```go
exporter := rete.NewPrometheusExporter(network)
http.Handle("/metrics", promhttp.Handler())
go http.ListenAndServe(":9090", nil)
```

**Métriques exportées** :
- `rete_rules_total` : Nombre total de règles
- `rete_facts_total` : Nombre total de faits
- `rete_activations_total` : Nombre d'activations
- `rete_rule_execution_duration_seconds` : Durée d'exécution par règle
- `rete_alpha_nodes_total` : Nombre de nœuds alpha
- `rete_beta_nodes_total` : Nombre de nœuds beta
- `rete_sharing_rate` : Taux de partage des nœuds

### Métriques de Chaînes

Suivi détaillé des alpha et beta chains :

```go
chainMetrics := network.GetChainMetrics()
fmt.Printf("Alpha chains: %d\n", chainMetrics.AlphaChainCount)
fmt.Printf("Beta chains: %d\n", chainMetrics.BetaChainCount)
fmt.Printf("Nodes shared: %d\n", chainMetrics.SharedNodeCount)
fmt.Printf("Memory saved: %d bytes\n", chainMetrics.MemorySaved)
```

### Profiling

Profiling automatisé des performances :

```go
profiler := rete.NewProfiler(network)
profiler.Start()

// Exécution des règles...

report := profiler.Stop()
fmt.Printf("Top 10 slowest rules:\n")
for _, rule := range report.SlowestRules {
    fmt.Printf("  %s: %v\n", rule.Name, rule.Duration)
}
```

---

## Ingestion Incrémentale

### Principe

Charger des fichiers TSD sans recréer tout le réseau :

```go
pipeline := rete.NewConstraintPipeline()
storage := rete.NewMemoryStorage()

// Premier chargement
network, err := pipeline.IngestFile("types.tsd", nil, storage)

// Chargement incrémental
network, err = pipeline.IngestFile("rules.tsd", network, storage)
network, err = pipeline.IngestFile("facts.tsd", network, storage)
```

### Validation Incrémentale

Validation avec contexte du réseau existant :

```go
network, err := pipeline.IngestFileWithIncrementalValidation(
    "new_rules.tsd",
    network,
    storage,
)
// Détecte les références à des types non définis
```

### Garbage Collection

Nettoyage automatique après reset :

```go
network, err := pipeline.IngestFileWithGC(
    "reset_and_reload.tsd",
    network,
    storage,
)
// Libère automatiquement la mémoire de l'ancien réseau
```

### Transactions

Ingestion transactionnelle avec rollback automatique :

```go
network, err := pipeline.IngestFile("critical_data.tsd", network, storage)
if err != nil {
    // Rollback automatique déjà effectué
    log.Printf("Erreur, état restauré : %v", err)
} else {
    // Commit automatique déjà effectué
    log.Println("Ingestion réussie")
}
```

### Configuration Avancée

```go
config := rete.DefaultAdvancedPipelineConfig()
config.EnableIncrementalValidation = true
config.EnableGCAfterReset = true
config.EnableTransactions = true
config.AutoCommit = true
config.AutoRollbackOnError = true

network, metrics, err := pipeline.IngestFileWithAdvancedFeatures(
    "data.tsd",
    network,
    storage,
    config,
)
```

---

## Commandes du Réseau

### Reset

Réinitialiser complètement le réseau :

```tsd
reset
```

Ou via l'API :

```go
network.Reset()
```

### Remove Fact

Retirer un fait spécifique :

```tsd
remove Person(name="Alice")
```

### Remove Rule

Retirer une règle par son identifiant :

```tsd
remove rule "check_adult"
```

Ou par son ID :

```tsd
remove rule id="R001"
```

---

## Nested OR Support

Support des expressions OR imbriquées :

```tsd
rule "complex_or" {
    when {
        p: Person(
            (age > 18 AND city == "Paris") OR
            (age > 21 AND city == "Lyon") OR
            (status == "VIP")
        )
    }
    then {
        print("Match: " + p.name)
    }
}
```

Le moteur normalise automatiquement ces expressions en forme normale disjonctive (DNF).

---

## Liens Utiles

- [Tutorial](TUTORIAL.md) - Guide pas à pas
- [API Reference](API_REFERENCE.md) - Documentation complète de l'API
- [Optimizations Guide](../rete/OPTIMIZATIONS_README.md) - Guide des optimisations
- [Strong Mode Tuning](STRONG_MODE_TUNING_GUIDE.md) - Réglage du Strong Mode
- [Examples](EXAMPLES.md) - Exemples d'utilisation

---

**Version** : 1.0  
**Dernière mise à jour** : Janvier 2025