# 🚀 Quick Start - Propagation Delta

Guide rapide pour utiliser le système de propagation sélective.

---

## 📋 Concepts Clés

**Propagation Delta** : Propager uniquement vers les nœuds affectés par les champs modifiés.

**Propagation Classique** : Retract+Insert (tous les nœuds ré-évalués).

**Mode Auto** : Choisit automatiquement le mode optimal selon des heuristiques.

---

## 🔧 Construction du Propagateur

### Minimal (avec défauts)

```go
import "github.com/treivax/tsd/rete/delta"

// Index de dépendances obligatoire
index := delta.NewDependencyIndex()

// Construction
propagator, err := delta.NewDeltaPropagatorBuilder().
    WithIndex(index).
    Build()

if err != nil {
    log.Fatal(err)
}
```

### Complet (configuration personnalisée)

```go
// Configuration
config := delta.DefaultPropagationConfig()
config.DeltaThreshold = 0.3           // 30% de changements max
config.MinFieldsForDelta = 5          // 5 champs minimum
config.MaxAffectedNodesForDelta = 50  // 50 nœuds max
config.EnableMetrics = true           // Activer métriques

// Détecteur de delta personnalisé
detectorConfig := delta.DefaultDetectorConfig()
detectorConfig.FloatEpsilon = 1e-6
detector := delta.NewDeltaDetectorWithConfig(detectorConfig)

// Stratégie de propagation
strategy := &delta.OptimizedStrategy{}

// Callback de propagation vers RETE
callback := func(nodeID string, factDelta *delta.FactDelta) error {
    // Logique de propagation vers le nœud RETE
    return reteNetwork.PropagateToNode(nodeID, factDelta)
}

// Construction complète
propagator, err := delta.NewDeltaPropagatorBuilder().
    WithIndex(index).
    WithDetector(detector).
    WithStrategy(strategy).
    WithConfig(config).
    WithPropagateCallback(callback).
    Build()
```

---

## 🎯 Utilisation

### Propager une Mise à Jour

```go
oldFact := map[string]interface{}{
    "id":    "p123",
    "name":  "Product A",
    "price": 100.0,
    "stock": 50,
}

newFact := map[string]interface{}{
    "id":    "p123",
    "name":  "Product A",
    "price": 120.0,  // Changement
    "stock": 50,
}

err := propagator.PropagateUpdate(
    oldFact, newFact,
    "Product~p123", "Product",
)

if err != nil {
    log.Printf("Propagation failed: %v", err)
}
```

### Avec Contexte et Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := propagator.PropagateUpdateWithContext(
    ctx,
    oldFact, newFact,
    "Product~p123", "Product",
)
```

---

## 📊 Métriques

### Consulter les Métriques

```go
metrics := propagator.GetMetrics()

fmt.Printf("Total propagations: %d\n", metrics.TotalPropagations)
fmt.Printf("Delta propagations: %d\n", metrics.DeltaPropagations)
fmt.Printf("Classic propagations: %d\n", metrics.ClassicPropagations)
fmt.Printf("Efficiency ratio: %.2f%%\n", metrics.GetEfficiencyRatio()*100)
fmt.Printf("Delta usage: %.2f%%\n", metrics.GetDeltaUsageRatio()*100)
fmt.Printf("Avg time: %v\n", metrics.AvgPropagationTime)
fmt.Printf("Nodes skipped: %d\n", metrics.NodesSkippedByDelta)
```

### Réinitialiser les Métriques

```go
propagator.ResetMetrics()
```

---

## ⚙️ Configuration

### Modifier la Configuration

```go
newConfig := propagator.GetConfig()
newConfig.DeltaThreshold = 0.5

err := propagator.UpdateConfig(newConfig)
if err != nil {
    log.Printf("Invalid config: %v", err)
}
```

### Modes de Propagation

```go
config := delta.DefaultPropagationConfig()

// Mode 1 : Toujours delta
config.DefaultMode = delta.PropagationModeDelta

// Mode 2 : Toujours classique
config.DefaultMode = delta.PropagationModeClassic

// Mode 3 : Automatique (recommandé)
config.DefaultMode = delta.PropagationModeAuto  // Défaut
```

---

## 🎨 Stratégies de Propagation

### Sequential (simple et prévisible)

```go
strategy := &delta.SequentialStrategy{}
// Ordre : alpha → beta → terminal
```

### Topological (respect dépendances)

```go
strategy := delta.NewTopologicalStrategy()
strategy.SetNodeDepth("alpha1", 1)
strategy.SetNodeDepth("beta1", 2)
strategy.SetNodeDepth("terminal1", 3)
// Ordre : profondeur croissante
```

### Optimized (hybride, recommandé)

```go
strategy := &delta.OptimizedStrategy{}
// Combine : type + factType + optimisations
```

---

## 🔍 Debugging

### Activer Logging Détaillé

```go
config := delta.DefaultPropagationConfig()
config.LogPropagationDetails = true
```

### Consulter Fallbacks

```go
metrics := propagator.GetMetrics()
fmt.Printf("Fallbacks due to ratio: %d\n", metrics.FallbacksDueToRatio)
fmt.Printf("Fallbacks due to nodes: %d\n", metrics.FallbacksDueToNodes)
fmt.Printf("Fallbacks due to PK: %d\n", metrics.FallbacksDueToPK)
fmt.Printf("Fallbacks due to error: %d\n", metrics.FallbacksDueToError)
```

---

## 🚨 Gestion d'Erreurs

### Retry sur Erreur

```go
config := delta.DefaultPropagationConfig()
config.RetryOnError = true  // Défaut : true

// Si propagation delta échoue → retry en mode classique
```

### Timeout

```go
config := delta.DefaultPropagationConfig()
config.PropagationTimeout = 30 * time.Second  // Défaut : 30s
```

---

## 📈 Optimisation

### Ajuster les Seuils

```go
config := delta.DefaultPropagationConfig()

// Favoriser delta (relâcher contraintes)
config.MinFieldsForDelta = 1          // Au lieu de 3
config.DeltaThreshold = 0.7           // Au lieu de 0.5
config.MaxAffectedNodesForDelta = 200 // Au lieu de 100

// Favoriser classic (renforcer contraintes)
config.MinFieldsForDelta = 10
config.DeltaThreshold = 0.2
config.MaxAffectedNodesForDelta = 20
```

### Concurrence

```go
config := delta.DefaultPropagationConfig()
config.MaxConcurrentPropagations = 10  // Défaut : 10

// Plus = meilleure parallélisation, plus de charge
// Moins = moins de charge, sérialisation
```

---

## 🎯 Cas d'Usage

### Mise à Jour Faible Nombre de Champs

```go
// 1 champ modifié sur 10 → Delta optimal
oldFact := map[string]interface{}{"field1": 1, ..., "field10": 10}
newFact := map[string]interface{}{"field1": 2, ..., "field10": 10}

// Propagation delta : uniquement nœuds dépendant de field1
```

### Mise à Jour Massive

```go
// 9 champs modifiés sur 10 → Classic optimal
oldFact := map[string]interface{}{"field1": 1, ..., "field10": 10}
newFact := map[string]interface{}{"field1": 99, ..., "field9": 99, "field10": 10}

// Propagation classique : Retract+Insert plus efficace
```

### Changement de Clé Primaire

```go
config := delta.DefaultPropagationConfig()
config.PrimaryKeyFields = []string{"id"}
config.AllowPrimaryKeyChange = false  // Défaut : false

// Changement de PK → Forcer mode classique
oldFact := map[string]interface{}{"id": "123", "name": "A"}
newFact := map[string]interface{}{"id": "456", "name": "A"}

// Auto fallback classique car PK changé
```

---

## ✅ Checklist Intégration

- [ ] DependencyIndex construit et peuplé
- [ ] Callback de propagation implémenté
- [ ] Configuration validée (`config.Validate()`)
- [ ] Stratégie choisie
- [ ] Métriques activées si besoin
- [ ] Tests d'intégration avec réseau RETE
- [ ] Gestion erreurs et fallback testés
- [ ] Performance monitorée

---

**Astuce** : Commencer avec `PropagationModeAuto` et ajuster selon les métriques.

**Documentation** : Voir `IMPLEMENTATION_REPORT_PROMPT05.md` pour détails complets.
