# 🗂️ Package Delta - Système de Propagation Optimale

## 📋 Vue d'ensemble

Le package `delta` implémente le **système de propagation optimale** pour le réseau RETE de TSD. Il comprend :

1. **Détection de changements** (DeltaDetector) - Identifier précisément les champs modifiés
2. **Indexation des dépendances** (DependencyIndex) - Déterminer quels nœuds sont affectés
3. **Modèle de données** (FieldDelta, FactDelta) - Représenter les changements

Ce système permet de propager les modifications de faits uniquement vers les nœuds RETE réellement impactés, évitant ainsi de ré-évaluer l'ensemble du réseau.

## 🎯 Objectifs

Le système de propagation delta vise à :

1. **Détecter précisément** les changements entre deux versions d'un fait
2. **Identifier rapidement** les nœuds RETE affectés par ces changements
3. **Propager efficacement** uniquement vers les nœuds concernés
4. **Éviter** la ré-évaluation complète du réseau (retract + insert classique)

**Gain attendu** : Jusqu'à 70% de réduction du temps de propagation pour des faits avec peu de changements.

## 🏗️ Architecture

### Composants principaux

#### 1. **DeltaDetector** (`delta_detector.go`) - Nouveau dans Prompt 04

Détecteur de changements entre deux versions d'un fait :

```go
detector := delta.NewDeltaDetector()

// Détecter les changements
oldFact := map[string]interface{}{
    "id": "p123",
    "price": 100.0,
    "status": "active",
}

newFact := map[string]interface{}{
    "id": "p123",
    "price": 150.0,  // Modifié
    "status": "active",
}

factDelta, err := detector.DetectDelta(
    oldFact, 
    newFact, 
    "Product~p123", 
    "Product",
)

// factDelta contient uniquement le champ "price"
```

**Fonctionnalités** :
- Détection complète (`DetectDelta`) ou optimisée (`DetectDeltaQuick`)
- Configuration flexible (epsilon floats, champs ignorés, comparaison profonde)
- Cache optionnel avec TTL pour comparaisons répétées
- Métriques de performance intégrées
- Thread-safe (sync.RWMutex)

#### 2. **DetectorConfig** (`detector_config.go`) - Nouveau dans Prompt 04

Configuration du détecteur pour différents cas d'usage :

```go
config := delta.DetectorConfig{
    FloatEpsilon:         0.01,  // Tolérance 1% pour floats
    IgnoreInternalFields: true,  // Ignorer champs "_*"
    IgnoredFields:        []string{"timestamp", "updated_at"},
    TrackTypeChanges:     true,  // Détecter 42 (int) → "42" (string)
    EnableDeepComparison: true,  // Comparaison récursive maps/slices
    MaxNestingLevel:      10,    // Protection stack overflow
    CacheComparisons:     true,  // Activer cache
    CacheTTL:             5 * time.Minute,
}

detector := delta.NewDeltaDetectorWithConfig(config)
```

#### 3. **DependencyIndex** (`dependency_index.go`)

Structure centrale qui maintient un index inversé :
- **Structure** : `factType → field → [nodeIDs]`
- **Exemple** : `Product → price → [alpha1, alpha2, beta3, terminal5]`
- **Thread-safety** : Toutes les opérations sont protégées par mutex

```go
idx := delta.NewDependencyIndex()

// Ajouter des nœuds
idx.AddAlphaNode("alpha1", "Product", []string{"price", "status"})
idx.AddBetaNode("beta1", "Order", []string{"customer_id"})
idx.AddTerminalNode("term1", "Product", []string{"price"})

// Requête : qui est affecté par Product.price ?
affected := idx.GetAffectedNodes("Product", "price")
// Retourne : [alpha1, term1]

// Requête avec un FactDelta
delta := delta.NewFactDelta("Product~123", "Product")
delta.AddFieldChange("price", 100.0, 150.0)
affected := idx.GetAffectedNodesForDelta(delta)
```

#### 4. **FieldExtractor** (`field_extractor.go`)

Extracteurs de champs depuis les structures AST :

- **AlphaConditionExtractor** : Extrait les champs depuis les conditions alpha
- **BetaConditionExtractor** : Extrait les champs depuis les conditions de jointure
- **ActionFieldExtractor** : Extrait les champs depuis les actions des nœuds terminaux

```go
// Exemple : extraction depuis une condition alpha
condition := map[string]interface{}{
    "type": "comparison",
    "left": map[string]interface{}{
        "type": "fieldAccess",
        "field": "price",
    },
    "right": 100,
}

fields, err := delta.ExtractFieldsFromAlphaCondition(condition)
// Retourne : ["price"]
```

#### 5. **IndexBuilder** (`index_builder.go`)

Orchestrateur de construction d'index depuis un réseau RETE complet :

```go
builder := delta.NewIndexBuilder()
builder.EnableDiagnostics() // Optionnel

idx := delta.NewDependencyIndex()

// Construire depuis différents types de nœuds
err := builder.BuildFromAlphaNode(idx, "alpha1", "Product", condition)
err := builder.BuildFromBetaNode(idx, "beta1", "Order", joinCondition)
err := builder.BuildFromTerminalNode(idx, "term1", "Product", actions)

// Consulter les diagnostics
diag := builder.GetDiagnostics()
fmt.Printf("Nœuds traités: %d, Champs extraits: %d\n", 
    diag.NodesProcessed, diag.FieldsExtracted)
```

#### 6. **IndexMetrics** (`index_metrics.go`)

Métriques de performance de l'index :

```go
metrics := delta.NewIndexMetrics()

// Enregistrer des opérations
metrics.RecordLookup(duration, nodesFound)
metrics.RecordNodeAdd()

// Consulter les statistiques
avgTime := metrics.GetAverageLookupTime()
avgNodes := metrics.GetAverageNodesPerLookup()

// Créer un snapshot
snapshot := metrics.Snapshot()
```

## 📊 Types de Nœuds Indexés

| Type de Nœud | Description | Champs Indexés |
|--------------|-------------|----------------|
| **Alpha** | Tests sur un seul fait | Champs testés dans les conditions |
| **Beta** | Jointures entre faits | Champs utilisés dans les comparaisons |
| **Terminal** | Actions de règles | Champs modifiés/lus dans les actions |

## 🚀 Utilisation

### Cas d'usage 1 : Construction d'index depuis un réseau RETE

```go
// Créer l'index et le builder
idx := delta.NewDependencyIndex()
builder := delta.NewIndexBuilder()
builder.EnableDiagnostics()

// Pour chaque nœud du réseau RETE, appeler le builder approprié
for _, alphaNode := range reteNetwork.GetAlphaNodes() {
    err := builder.BuildFromAlphaNode(
        idx,
        alphaNode.ID,
        alphaNode.FactType,
        alphaNode.Condition,
    )
    if err != nil {
        log.Printf("Erreur indexation alpha: %v", err)
    }
}

// Statistiques finales
stats := idx.GetStats()
fmt.Printf("Index construit : %d nœuds, %d champs\n", 
    stats.NodeCount, stats.FieldCount)
```

### Cas d'usage 2 : Propagation delta complète (Prompt 04)

```go
// 1. Initialiser le détecteur et l'index
detector := delta.NewDeltaDetector()
idx := delta.NewDependencyIndex()
// ... (construction de l'index depuis le réseau RETE)

// 2. Capture du fait avant modification
oldFact := getCurrentFact("Product~123")

// 3. Application de la modification
newFact := applyModification(oldFact, updates)

// 4. Détection des changements
factDelta, err := detector.DetectDelta(
    oldFact,
    newFact,
    "Product~123",
    "Product",
)

if err != nil {
    return err
}

// 5. Décision de propagation basée sur le ratio
if factDelta.ChangeRatio() < 0.3 {
    // < 30% de changements → propagation delta optimisée
    affectedNodes := idx.GetAffectedNodesForDelta(factDelta)
    
    // 6. Propager uniquement vers nœuds affectés
    for _, nodeRef := range affectedNodes {
        propagateDeltaToNode(nodeRef, factDelta)
    }
    
    fmt.Printf("✅ Delta propagation: %d nœuds affectés sur %d\n", 
        len(affectedNodes), totalNodes)
} else {
    // >= 30% de changements → retract + insert classique plus efficace
    retractFact(oldFact)
    insertFact(newFact)
    
    fmt.Printf("⚠️  Classique (trop de changements): ratio=%.2f\n", 
        factDelta.ChangeRatio())
}
```

### Cas d'usage 3 : Optimisation avec DetectDeltaQuick

```go
detector := delta.NewDeltaDetector()

// Pour des comparaisons fréquentes sans changements
factDelta, err := detector.DetectDeltaQuick(
    oldFact,
    newFact,
    factID,
    factType,
)

if factDelta == nil {
    // Aucun changement → pas de propagation nécessaire
    return nil
}

// Changements détectés → continuer la propagation
processDelta(factDelta)
```

### Cas d'usage 4 : Métriques et monitoring

```go
detector := delta.NewDeltaDetector()

// Faire des détections
for _, modification := range modifications {
    detector.DetectDelta(
        modification.Old,
        modification.New,
        modification.ID,
        modification.Type,
    )
}

// Obtenir les métriques
metrics := detector.GetMetrics()
fmt.Printf("Comparaisons effectuées: %d\n", metrics.Comparisons)
fmt.Printf("Cache hits: %d (%.2f%%)\n", 
    metrics.CacheHits, 
    metrics.HitRate * 100)
fmt.Printf("Taille du cache: %d entrées\n", metrics.CacheSize)

// Reset pour nouveau cycle
detector.ResetMetrics()
```

### Cas d'usage 5 : Propagation delta (legacy)

```go
// Un fait a été modifié
factDelta := delta.NewFactDelta("Product~123", "Product")
factDelta.AddFieldChange("price", 100.0, 150.0)
factDelta.AddFieldChange("status", "active", "inactive")

// Trouver les nœuds affectés
affectedNodes := idx.GetAffectedNodesForDelta(factDelta)

// Propager uniquement vers ces nœuds
for _, nodeRef := range affectedNodes {
    switch nodeRef.NodeType {
    case "alpha":
        // Ré-évaluer nœud alpha
    case "beta":
        // Ré-évaluer nœud beta
    case "terminal":
        // Ré-évaluer nœud terminal
    }
}
```

### Cas d'usage 6 : Analyse de dépendances

```go
// Analyser les dépendances d'un champ
affectedByPrice := idx.GetAffectedNodes("Product", "price")

fmt.Printf("Le champ Product.price affecte %d nœuds:\n", len(affectedByPrice))
for _, node := range affectedByPrice {
    fmt.Printf("  - %s\n", node.String())
    // Exemple: "alpha[alpha_price_check](Product)"
}
```

## 📈 Performances

### DeltaDetector (Prompt 04)

D'après les benchmarks :

| Opération | Temps | Allocations | Performance |
|-----------|-------|-------------|-------------|
| DetectDeltaQuick (no-op) | 140.6 ns/op | 0 B | ⭐⭐⭐⭐⭐ |
| DetectDelta (no-op) | 454.9 ns/op | 128 B | ⭐⭐⭐⭐ |
| DetectDelta (single change) | 645.5 ns/op | 832 B | ⭐⭐⭐⭐ |
| DetectDelta (multiple changes) | 664.0 ns/op | 832 B | ⭐⭐⭐⭐ |
| DetectDelta (large 50 fields) | 5,530 ns/op | 4,040 B | ⭐⭐⭐⭐ |
| DetectDelta (with cache) | 40.27 ns/op | 0 B | ⭐⭐⭐⭐⭐ |
| DetectDelta (deep nested) | 605.3 ns/op | 832 B | ⭐⭐⭐⭐ |

**Points clés** :
- DetectDeltaQuick **3× plus rapide** que DetectDelta pour faits identiques
- Cache activé offre un gain de **16×** (40ns vs 645ns)
- 0 allocation pour DetectDeltaQuick si aucun changement
- Performance linéaire : ~110ns par champ

### DependencyIndex (Prompt 03)

D'après les benchmarks :

| Opération | Temps | Allocations |
|-----------|-------|-------------|
| AddAlphaNode | ~161 µs | 565 B/op |
| GetAffectedNodes (100 nœuds) | ~8.3 µs | 14.9 KB/op |
| GetAffectedNodesForDelta | ~4.7 µs | 7.3 KB/op |
| ExtractFieldsFromAlphaCondition | ~393 ns | 32 B/op |

### Estimation mémoire

Pour un index typique :
- **Par nœud** : ~150 bytes (nodeID + NodeReference + overhead)
- **Par entrée de champ** : ~74 bytes (index entry)
- **Exemple** : 100 nœuds avec 3 champs chacun ≈ 37 KB

## ✅ Tests et Validation

### Couverture de tests

```bash
go test ./rete/delta/... -cover
# coverage: 85.4% of statements
```

**Par composant** :
- DeltaDetector : 100% des cas critiques couverts
- DetectorConfig : Validation complète
- DependencyIndex : 83.8% de couverture
- FieldExtractor : Tests complets
- IndexBuilder : Diagnostics validés

### Exécuter les tests

```bash
# Tests unitaires
go test ./rete/delta/...

# Tests avec verbose
go test -v ./rete/delta/...

# Tests d'intégration
go test -v ./rete/delta/... -run TestIndexation

# Benchmarks
go test -bench=. -benchmem ./rete/delta/...
```

## 🔄 Workflow d'intégration

1. **Construction de l'index** : Lors de la compilation des règles
2. **Mise à jour incrémentale** : Lors de l'ajout/suppression de règles
3. **Interrogation** : Lors de la propagation delta
4. **Maintenance** : Clear() si reconstruction complète nécessaire

## 📝 Implémentation et Limitations

### Composants implémentés (Prompt 01-04)

- ✅ Modèle de données (FieldDelta, FactDelta, ChangeType, ValueType)
- ✅ Comparaison de valeurs (ValuesEqual, FactsEqual)
- ✅ Détecteur de changements (DeltaDetector, DetectorConfig)
- ✅ Index de dépendances (DependencyIndex)
- ✅ Extraction de champs (FieldExtractor)
- ✅ Construction d'index (IndexBuilder)
- ✅ Métriques de performance (IndexMetrics, DetectorMetrics)
- ✅ Tests complets (85.4% couverture)
- ✅ Benchmarks de performance

### Intégration future (Prompt 05)

- [ ] Intégration dans Network.ModifyFact()
- [ ] Génération de commandes delta pour propagation
- [ ] Seuil de ratio configurable
- [ ] Métriques de gain end-to-end
- [ ] Tests d'intégration complets

### Limitations actuelles

- **BuildFromNetwork()** : Retourne un index vide (stub pour conception)
- **Champs imbriqués** : Pas de support des nested fields dans l'extraction
- **Validation de cohérence** : Pas de vérification entre types et champs
- **Cache DeltaDetector** : Clé basée uniquement sur factID (hashing possible)

## 📚 Documentation complémentaire

### Rapports de Validation

- **Prompt 02** : `IMPLEMENTATION_REPORT_PROMPT02.md` - Modèle de données delta
- **Prompt 03** : `IMPLEMENTATION_REPORT_PROMPT03.md` - Indexation des dépendances
- **Prompt 04** : `VALIDATION_REPORT_PROMPT04.md` - Détection delta
- **Prompt 04** : `EXECUTION_SUMMARY_PROMPT04.md` - Résumé d'exécution

### Documentation Technique

- **Conception globale** : `REPORTS/conception_delta_architecture.md`
- **Métadonnées nœuds** : `REPORTS/metadata_noeuds.md`
- **Mapping AST** : `REPORTS/ast_conditions_mapping.md`

### Spécifications

- **Prompt 02** : `scripts/propagation_optimale/02_modelisation_delta.md`
- **Prompt 03** : `scripts/propagation_optimale/03_indexation_dependances.md`
- **Prompt 04** : `scripts/propagation_optimale/04_detection_delta.md`

## 🤝 Contribution

Pour modifier le système d'indexation :

1. Respecter les standards de `common.md`
2. Maintenir la couverture > 80%
3. Ajouter des tests pour tout nouveau comportement
4. Documenter les fonctions exportées (GoDoc)
5. Mettre à jour ce README si changement d'API

## 📄 Licence

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License
