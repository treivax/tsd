# Intégration Prometheus pour RETE TSD

## Vue d'ensemble

Ce guide explique comment intégrer les métriques RETE avec Prometheus pour le monitoring en production.

## Table des Matières

1. [Installation](#installation)
2. [Configuration](#configuration)
3. [Métriques Disponibles](#métriques-disponibles)
4. [Exemples d'Utilisation](#exemples-dutilisation)
5. [Dashboard Grafana](#dashboard-grafana)
6. [Alerting](#alerting)
7. [Bonnes Pratiques](#bonnes-pratiques)

---

## Installation

### Prérequis

- Prometheus installé et configuré
- Application TSD utilisant le package `rete`

### Activation dans TSD

```go
import "github.com/treivax/tsd/rete"

// Configuration avec Prometheus activé
config := rete.DefaultChainPerformanceConfig()
config.PrometheusEnabled = true
config.PrometheusPrefix = "tsd_rete"

// Ou utiliser la configuration haute performance (Prometheus activé par défaut)
config := rete.HighPerformanceConfig()
```

---

## Configuration

### Configuration Basique

```go
storage := rete.NewMemoryStorage()
network := rete.NewReteNetwork(storage)

// Créer l'exporteur Prometheus
config := rete.DefaultChainPerformanceConfig()
config.PrometheusEnabled = true

exporter := rete.NewPrometheusExporter(network.ChainMetrics, config)
exporter.RegisterMetrics()

// Démarrer le serveur HTTP pour /metrics
go func() {
    if err := exporter.ServeHTTP(":9090"); err != nil {
        log.Fatal(err)
    }
}()
```

### Configuration Avancée

```go
// Créer l'exporteur avec mise à jour automatique
exporter := rete.NewPrometheusExporter(network.ChainMetrics, config)
exporter.RegisterMetrics()

// Mise à jour automatique toutes les 10 secondes
stopAutoUpdate := exporter.StartAutoUpdate(10 * time.Second)
defer close(stopAutoUpdate)

// Exposer les métriques sur un endpoint personnalisé
http.Handle("/custom/metrics", exporter.Handler())
http.ListenAndServe(":8080", nil)
```

### Configuration Prometheus

Ajouter dans `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'tsd_rete'
    scrape_interval: 15s
    static_configs:
      - targets: ['localhost:9090']
        labels:
          application: 'tsd'
          environment: 'production'
```

---

## Métriques Disponibles

### Métriques de Chaînes

#### `tsd_rete_chains_built_total`
- **Type**: Counter
- **Description**: Nombre total de chaînes alpha construites
- **Utilisation**: Suivre le volume de règles traitées

```promql
# Taux de construction de chaînes par seconde
rate(tsd_rete_chains_built_total[5m])
```

#### `tsd_rete_chains_length_avg`
- **Type**: Gauge
- **Description**: Longueur moyenne des chaînes alpha
- **Utilisation**: Identifier la complexité des règles

```promql
# Longueur moyenne des chaînes
tsd_rete_chains_length_avg
```

### Métriques de Nœuds

#### `tsd_rete_nodes_created_total`
- **Type**: Counter
- **Description**: Nombre total de nœuds alpha créés

#### `tsd_rete_nodes_reused_total`
- **Type**: Counter
- **Description**: Nombre total de nœuds alpha réutilisés

#### `tsd_rete_nodes_sharing_ratio`
- **Type**: Gauge
- **Description**: Ratio de partage de nœuds (0.0 à 1.0)
- **Utilisation**: Mesurer l'efficacité du partage

```promql
# Ratio de partage en pourcentage
tsd_rete_nodes_sharing_ratio * 100
```

### Métriques de Cache de Hash

#### `tsd_rete_hash_cache_hits_total`
- **Type**: Counter
- **Description**: Nombre total de hits du cache de hash

#### `tsd_rete_hash_cache_misses_total`
- **Type**: Counter
- **Description**: Nombre total de misses du cache de hash

#### `tsd_rete_hash_cache_size`
- **Type**: Gauge
- **Description**: Taille actuelle du cache de hash

#### `tsd_rete_hash_cache_efficiency`
- **Type**: Gauge
- **Description**: Efficacité du cache de hash (0.0 à 1.0)

```promql
# Taux de hit du cache (%)
tsd_rete_hash_cache_efficiency * 100

# Taux de hit calculé manuellement
rate(tsd_rete_hash_cache_hits_total[5m]) / 
(rate(tsd_rete_hash_cache_hits_total[5m]) + rate(tsd_rete_hash_cache_misses_total[5m]))
```

### Métriques de Cache de Connexion

#### `tsd_rete_connection_cache_hits_total`
- **Type**: Counter
- **Description**: Nombre total de hits du cache de connexion

#### `tsd_rete_connection_cache_misses_total`
- **Type**: Counter
- **Description**: Nombre total de misses du cache de connexion

#### `tsd_rete_connection_cache_efficiency`
- **Type**: Gauge
- **Description**: Efficacité du cache de connexion (0.0 à 1.0)

### Métriques de Temps

#### `tsd_rete_build_time_seconds_total`
- **Type**: Counter
- **Description**: Temps total passé à construire des chaînes (en secondes)

#### `tsd_rete_build_time_seconds_avg`
- **Type**: Gauge
- **Description**: Temps moyen de construction d'une chaîne (en secondes)

```promql
# Temps moyen de construction en millisecondes
tsd_rete_build_time_seconds_avg * 1000
```

#### `tsd_rete_hash_compute_time_seconds_total`
- **Type**: Counter
- **Description**: Temps total de calcul des hash (en secondes)

---

## Exemples d'Utilisation

### Exemple Complet

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
    
    "github.com/treivax/tsd/rete"
)

func main() {
    // Configuration
    storage := rete.NewMemoryStorage()
    network := rete.NewReteNetwork(storage)
    
    config := rete.DefaultChainPerformanceConfig()
    config.PrometheusEnabled = true
    config.PrometheusPrefix = "myapp_rete"
    
    // Créer et configurer l'exporteur
    exporter := rete.NewPrometheusExporter(network.ChainMetrics, config)
    exporter.RegisterMetrics()
    
    // Mise à jour automatique des métriques
    stopUpdate := exporter.StartAutoUpdate(5 * time.Second)
    defer close(stopUpdate)
    
    // Endpoints HTTP
    http.Handle("/metrics", exporter.Handler())
    http.HandleFunc("/health", healthHandler)
    
    // Démarrer le serveur
    go func() {
        log.Printf("Métriques Prometheus disponibles sur :9090/metrics")
        if err := http.ListenAndServe(":9090", nil); err != nil {
            log.Fatal(err)
        }
    }()
    
    // Construire des règles (exemple)
    buildRules(network, storage)
    
    // Garder l'application en cours d'exécution
    select {}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "OK")
}

func buildRules(network *rete.ReteNetwork, storage rete.Storage) {
    builder := rete.NewAlphaChainBuilderWithMetrics(
        network, 
        storage, 
        network.ChainMetrics,
    )
    
    for i := 0; i < 100; i++ {
        conditions := []rete.SimpleCondition{
            {
                Type:     "binaryOperation",
                Left:     map[string]interface{}{"type": "variable", "name": "age"},
                Operator: ">",
                Right:    map[string]interface{}{"type": "literal", "value": float64(18)},
            },
        }
        
        ruleID := fmt.Sprintf("rule_%d", i)
        _, err := builder.BuildChain(conditions, "person", network.RootNode, ruleID)
        if err != nil {
            log.Printf("Erreur construction règle %s: %v", ruleID, err)
        }
    }
}
```

### Vérifier les Métriques

```bash
# Vérifier que le endpoint fonctionne
curl http://localhost:9090/metrics

# Filtrer les métriques RETE
curl http://localhost:9090/metrics | grep tsd_rete

# Vérifier une métrique spécifique
curl -s http://localhost:9090/metrics | grep tsd_rete_chains_built_total
```

---

## Dashboard Grafana

### Import du Dashboard

Créer un fichier `tsd-rete-dashboard.json`:

```json
{
  "dashboard": {
    "title": "TSD RETE Performance",
    "panels": [
      {
        "title": "Chaînes Construites",
        "targets": [
          {
            "expr": "rate(tsd_rete_chains_built_total[5m])"
          }
        ],
        "type": "graph"
      },
      {
        "title": "Ratio de Partage (%)",
        "targets": [
          {
            "expr": "tsd_rete_nodes_sharing_ratio * 100"
          }
        ],
        "type": "stat"
      },
      {
        "title": "Efficacité Cache Hash (%)",
        "targets": [
          {
            "expr": "tsd_rete_hash_cache_efficiency * 100"
          }
        ],
        "type": "gauge"
      },
      {
        "title": "Temps Moyen Construction (ms)",
        "targets": [
          {
            "expr": "tsd_rete_build_time_seconds_avg * 1000"
          }
        ],
        "type": "graph"
      }
    ]
  }
}
```

### Requêtes PromQL Utiles

```promql
# Taux de construction de chaînes par minute
rate(tsd_rete_chains_built_total[1m]) * 60

# Nombre de nœuds créés vs réutilisés (dernière heure)
sum(increase(tsd_rete_nodes_created_total[1h]))
sum(increase(tsd_rete_nodes_reused_total[1h]))

# Efficacité globale du cache
(tsd_rete_hash_cache_efficiency + tsd_rete_connection_cache_efficiency) / 2

# Latence p95 (si histogrammes disponibles)
histogram_quantile(0.95, rate(tsd_rete_build_time_seconds[5m]))

# Tendance du ratio de partage sur 24h
avg_over_time(tsd_rete_nodes_sharing_ratio[24h])
```

---

## Alerting

### Règles d'Alerte Prometheus

Créer un fichier `rete-alerts.yml`:

```yaml
groups:
  - name: tsd_rete_alerts
    interval: 30s
    rules:
      # Alerte si le ratio de partage chute
      - alert: LowNodeSharingRatio
        expr: tsd_rete_nodes_sharing_ratio < 0.3
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Faible ratio de partage de nœuds"
          description: "Le ratio de partage est à {{ $value | humanizePercentage }}"
      
      # Alerte si le cache de hash est inefficace
      - alert: LowHashCacheEfficiency
        expr: tsd_rete_hash_cache_efficiency < 0.5
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "Efficacité du cache de hash faible"
          description: "L'efficacité du cache de hash est à {{ $value | humanizePercentage }}"
      
      # Alerte si le temps de construction augmente
      - alert: HighBuildTime
        expr: tsd_rete_build_time_seconds_avg > 0.001
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Temps de construction élevé"
          description: "Le temps moyen de construction est {{ $value | humanizeDuration }}"
      
      # Alerte si le cache atteint sa capacité max
      - alert: CacheNearCapacity
        expr: |
          (tsd_rete_hash_cache_size / 10000) > 0.9
        for: 15m
        labels:
          severity: info
        annotations:
          summary: "Cache de hash proche de la capacité"
          description: "Le cache est rempli à {{ $value | humanizePercentage }}"
```

### Configuration AlertManager

```yaml
route:
  group_by: ['alertname', 'severity']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'team-ops'
  routes:
    - match:
        severity: critical
      receiver: 'pagerduty'

receivers:
  - name: 'team-ops'
    email_configs:
      - to: 'ops@example.com'
  
  - name: 'pagerduty'
    pagerduty_configs:
      - service_key: '<your-key>'
```

---

## Bonnes Pratiques

### 1. Préfixe des Métriques

Toujours utiliser un préfixe cohérent:

```go
config.PrometheusPrefix = "myapp_rete"  // ✅ Bon
config.PrometheusPrefix = "rete"        // ❌ Trop générique
```

### 2. Labels

Ajouter des labels pour segmenter les métriques:

```go
// Exemple d'ajout de labels personnalisés (à implémenter)
labels := map[string]string{
    "environment": "production",
    "region":      "us-east-1",
    "version":     "1.2.3",
}
```

### 3. Intervalles de Scraping

- **Production**: 15-30 secondes
- **Développement**: 5-10 secondes
- **Tests**: 1-5 secondes

### 4. Rétention des Données

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

storage:
  tsdb:
    retention.time: 15d  # Conserver 15 jours de données
    retention.size: 50GB # Ou limite de taille
```

### 5. Cardinalité

⚠️ **Attention**: Ne pas créer de métriques avec trop de labels uniques

```go
// ❌ Mauvais (cardinalité élevée)
metric_by_rule_id{rule_id="rule_1"}
metric_by_rule_id{rule_id="rule_2"}
// ... 10000 règles = 10000 séries temporelles

// ✅ Bon (cardinalité contrôlée)
total_rules 10000
```

### 6. Monitoring du Monitoring

Surveiller Prometheus lui-même:

```promql
# Utilisation mémoire de Prometheus
process_resident_memory_bytes

# Durée de scraping
scrape_duration_seconds

# Nombre de séries temporelles
prometheus_tsdb_symbol_table_size_bytes
```

---

## Dépannage

### Métriques Manquantes

```bash
# Vérifier que le endpoint est accessible
curl http://localhost:9090/metrics

# Vérifier les logs Prometheus
journalctl -u prometheus -f

# Vérifier la configuration Prometheus
promtool check config prometheus.yml
```

### Valeurs Incorrectes

```go
// Forcer une mise à jour des métriques
exporter.UpdateMetrics()

// Vérifier les métriques brutes
snapshot := network.GetChainMetrics().GetSnapshot()
fmt.Printf("%+v\n", snapshot)
```

### Performance

Si Prometheus ralentit:

1. Réduire l'intervalle de scraping
2. Activer la compression
3. Augmenter les ressources
4. Réduire la rétention

```yaml
# Activer la compression
global:
  scrape_timeout: 10s
  
scrape_configs:
  - job_name: 'tsd_rete'
    scrape_interval: 30s  # Augmenter l'intervalle
    metric_relabel_configs:
      - action: drop
        regex: 'unwanted_.*'  # Supprimer les métriques non nécessaires
```

---

## Exemples de Requêtes Avancées

### Analyse de Tendances

```promql
# Tendance du ratio de partage (régression linéaire)
predict_linear(tsd_rete_nodes_sharing_ratio[1h], 3600)

# Détection d'anomalies
abs(tsd_rete_chains_built_total - avg_over_time(tsd_rete_chains_built_total[1h])) > 
  2 * stddev_over_time(tsd_rete_chains_built_total[1h])
```

### Agrégation Multi-Instances

```promql
# Total de chaînes construites sur toutes les instances
sum(tsd_rete_chains_built_total)

# Ratio de partage moyen par région
avg by(region) (tsd_rete_nodes_sharing_ratio)
```

### Calculs Dérivés

```promql
# Taux de création de nœuds par seconde
rate(tsd_rete_nodes_created_total[5m])

# Ratio création/réutilisation
rate(tsd_rete_nodes_created_total[5m]) / 
rate(tsd_rete_nodes_reused_total[5m])
```

---

## Ressources

### Documentation

- [Prometheus Docs](https://prometheus.io/docs/)
- [PromQL Cheatsheet](https://promlabs.com/promql-cheat-sheet/)
- [Grafana Docs](https://grafana.com/docs/)

### Outils

- [PromLens](https://promlens.com/) - Éditeur PromQL
- [Promtool](https://prometheus.io/docs/prometheus/latest/command-line/promtool/) - Validation
- [Grafana Play](https://play.grafana.org/) - Test de dashboards

### Support

- GitHub Issues: [https://github.com/treivax/tsd/issues](https://github.com/treivax/tsd/issues)
- Documentation: `docs/CHAIN_PERFORMANCE_OPTIMIZATION.md`

---

## License

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License