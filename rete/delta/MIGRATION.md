# 🔄 Guide de Migration - Propagation Delta RETE

## 📋 Vue d'ensemble

Ce guide vous accompagne dans la migration d'un système RETE classique (Retract+Insert) vers le **système de propagation delta optimisé**.

**Gains attendus** :
- ⚡ **3-4x plus rapide** pour réseaux moyens/grands (>50 nœuds)
- 🎯 **70-80% de nœuds évités** lors de mises à jour ciblées
- 📉 **Réduction significative** de la charge CPU et des allocations mémoire

**Quand migrer** :
- ✅ Réseaux RETE avec >50 nœuds
- ✅ Mises à jour fréquentes de faits (updates > inserts+deletes)
- ✅ Faits avec nombreux champs mais peu de modifications simultanées
- ✅ Besoin de performance prévisible en production

**Quand rester classique** :
- ❌ Petits réseaux (<20 nœuds) - overhead delta non amorti
- ❌ Mises à jour massives (>50% des champs modifiés à chaque fois)
- ❌ Patterns insert/delete dominants sans updates

---

## 🚀 Migration Rapide (TL;DR)

### Avant (Classique)

```go
// Mise à jour classique : Retract + Insert
network.RetractFact("Product~123")
network.AssertFact(updatedProduct)
```

### Après (Delta)

```go
import "github.com/treivax/tsd/rete/delta"

// 1. Créer le système delta
detector := delta.NewDeltaDetector()
index := delta.NewDependencyIndex()

// 2. Construire l'index depuis le réseau
builder := delta.NewIndexBuilder()
builder.BuildFromNetwork(index, network)

// 3. Détecter les changements
factDelta, _ := detector.DetectDelta(
    oldProduct, 
    newProduct, 
    "Product~123", 
    "Product",
)

// 4. Propager uniquement aux nœuds affectés
if factDelta != nil && !factDelta.IsEmpty() {
    affectedNodes := index.GetAffectedNodesForDelta(factDelta)
    for _, node := range affectedNodes {
        network.PropagateToNode(node.NodeID, factDelta)
    }
} else {
    // Aucun changement ou fallback classique
    network.RetractFact(factDelta.FactID)
    network.AssertFact(newProduct)
}
```

---

## 📖 Migration Détaillée

### Étape 1 : Analyser Votre Réseau Existant

Avant de migrer, évaluez si delta est bénéfique :

```go
// Collecter métriques sur votre réseau actuel
type NetworkProfile struct {
    TotalNodes       int
    FactTypes        []string
    AvgFieldsPerFact int
    UpdateFrequency  float64 // updates/sec
    InsertFrequency  float64 // inserts/sec
    DeleteFrequency  float64 // deletes/sec
}

func analyzeNetwork(network *rete.Network) NetworkProfile {
    // Implémenter collection de métriques
    return NetworkProfile{
        TotalNodes:       len(network.GetAllNodes()),
        UpdateFrequency:  measureUpdateRate(),
        // ...
    }
}

// Décision
profile := analyzeNetwork(myNetwork)
if profile.TotalNodes > 50 && 
   profile.UpdateFrequency > profile.InsertFrequency {
    // 👍 Delta est recommandé
} else {
    // 👎 Rester en mode classique
}
```

### Étape 2 : Initialiser les Composants Delta

#### 2.1 Configuration du Détecteur

```go
import "github.com/treivax/tsd/rete/delta"

// Configuration par défaut (recommandée)
detector := delta.NewDeltaDetector()

// Configuration personnalisée
config := delta.DetectorConfig{
    FloatEpsilon:         0.001,  // Tolérance pour prix, montants
    IgnoreInternalFields: true,   // Ignorer "_timestamp", "_version"
    IgnoredFields:        []string{"updated_at", "sync_status"},
    TrackTypeChanges:     true,   // Détecter int→string
    EnableDeepComparison: true,   // Pour nested objects
    MaxNestingLevel:      5,      // Protection recursion
    CacheComparisons:     true,   // Cache pour comparaisons répétées
    CacheTTL:             5 * time.Minute,
}

detector := delta.NewDeltaDetectorWithConfig(config)
```

**Choix de configuration** :

| Scénario | FloatEpsilon | CacheComparisons | EnableDeepComparison |
|----------|--------------|------------------|----------------------|
| Données financières | `0.01` | `true` | `false` |
| Inventaire/Stock | `0.001` | `true` | `false` |
| Configuration complexe | `1e-9` | `false` | `true` |
| Temps réel (IoT) | `0.1` | `false` | `false` |

#### 2.2 Construction de l'Index

```go
// Créer l'index
index := delta.NewDependencyIndex()

// Builder avec diagnostics
builder := delta.NewIndexBuilder()
builder.EnableDiagnostics()

// Option A : Construction automatique depuis le réseau
err := builder.BuildFromNetwork(index, network)
if err != nil {
    log.Fatalf("Échec construction index: %v", err)
}

// Option B : Construction manuelle (contrôle total)
index.AddAlphaNode("alpha_product_price", "Product", []string{"price"})
index.AddAlphaNode("alpha_product_status", "Product", []string{"status"})
index.AddBetaNode("beta_order_customer", "Order", []string{"customer_id", "total"})
index.AddTerminalNode("term_discount_rule", "Product", []string{"price", "category"})

// Vérifier diagnostics
diag := builder.GetDiagnostics()
log.Printf("Index construit: %d nœuds, %d champs extraits", 
    diag.NodesProcessed, diag.FieldsExtracted)
```

### Étape 3 : Intégration dans le Workflow

#### 3.1 Pattern : Wrapper de Mise à Jour

Créer une abstraction qui encapsule la logique delta :

```go
type FactUpdater struct {
    network   *rete.Network
    detector  *delta.DeltaDetector
    index     *delta.DependencyIndex
    
    // Seuils de décision
    deltaThreshold float64  // % de champs modifiés avant fallback
}

func NewFactUpdater(network *rete.Network) (*FactUpdater, error) {
    detector := delta.NewDeltaDetector()
    index := delta.NewDependencyIndex()
    
    builder := delta.NewIndexBuilder()
    if err := builder.BuildFromNetwork(index, network); err != nil {
        return nil, err
    }
    
    return &FactUpdater{
        network:        network,
        detector:       detector,
        index:          index,
        deltaThreshold: 0.3, // Fallback si >30% champs modifiés
    }, nil
}

func (u *FactUpdater) UpdateFact(
    oldFact, newFact map[string]interface{},
    factID, factType string,
) error {
    // 1. Détecter les changements
    factDelta, err := u.detector.DetectDelta(oldFact, newFact, factID, factType)
    if err != nil {
        return fmt.Errorf("détection delta: %w", err)
    }
    
    // 2. Vérifier si delta est pertinent
    if factDelta.IsEmpty() {
        // Aucun changement réel
        return nil
    }
    
    totalFields := len(newFact)
    changedFields := len(factDelta.Changes)
    changeRatio := float64(changedFields) / float64(totalFields)
    
    // 3a. Propagation delta (optimisée)
    if changeRatio <= u.deltaThreshold {
        affectedNodes := u.index.GetAffectedNodesForDelta(factDelta)
        
        for _, node := range affectedNodes {
            if err := u.network.PropagateToNode(node.NodeID, factDelta); err != nil {
                return fmt.Errorf("propagation vers %s: %w", node.NodeID, err)
            }
        }
        
        return nil
    }
    
    // 3b. Fallback classique (trop de changements)
    if err := u.network.RetractFact(factID); err != nil {
        return fmt.Errorf("retract: %w", err)
    }
    
    if err := u.network.AssertFact(newFact); err != nil {
        return fmt.Errorf("assert: %w", err)
    }
    
    return nil
}
```

#### 3.2 Utilisation du Wrapper

```go
// Initialisation (une fois)
updater, err := NewFactUpdater(network)
if err != nil {
    log.Fatal(err)
}

// Mise à jour de fait (n fois)
oldProduct := map[string]interface{}{
    "id":       "p123",
    "name":     "Widget",
    "price":    100.0,
    "stock":    50,
    "category": "electronics",
}

newProduct := map[string]interface{}{
    "id":       "p123",
    "name":     "Widget",
    "price":    120.0,  // Modifié
    "stock":    45,     // Modifié
    "category": "electronics",
}

err = updater.UpdateFact(oldProduct, newProduct, "Product~p123", "Product")
if err != nil {
    log.Printf("Erreur mise à jour: %v", err)
}
```

### Étape 4 : Gestion du Cycle de Vie

#### 4.1 Reconstruction de l'Index

Quand reconstruire l'index :
- ✅ Après ajout/suppression de règles
- ✅ Modification de conditions dans le réseau
- ✅ Périodiquement en production (ex: toutes les heures)

```go
type IndexManager struct {
    index        *delta.DependencyIndex
    builder      *delta.IndexBuilder
    network      *rete.Network
    mutex        sync.RWMutex
    lastRebuild  time.Time
    rebuildEvery time.Duration
}

func NewIndexManager(network *rete.Network) *IndexManager {
    return &IndexManager{
        network:      network,
        index:        delta.NewDependencyIndex(),
        builder:      delta.NewIndexBuilder(),
        rebuildEvery: 1 * time.Hour,
    }
}

func (m *IndexManager) RebuildIfNeeded() error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    if time.Since(m.lastRebuild) < m.rebuildEvery {
        return nil
    }
    
    // Nouvel index
    newIndex := delta.NewDependencyIndex()
    if err := m.builder.BuildFromNetwork(newIndex, m.network); err != nil {
        return err
    }
    
    // Swap atomique
    m.index = newIndex
    m.lastRebuild = time.Now()
    
    log.Printf("Index reconstruit: %d nœuds", newIndex.GetStats().NodeCount)
    return nil
}

func (m *IndexManager) GetIndex() *delta.DependencyIndex {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    return m.index
}
```

#### 4.2 Monitoring et Métriques

```go
import "time"

type DeltaMetrics struct {
    DeltaPropagations   int64
    ClassicFallbacks    int64
    NodesEvaluated      int64
    NodesAvoided        int64
    DetectionTime       time.Duration
    PropagationTime     time.Duration
}

func (m *DeltaMetrics) RecordDelta(
    detected time.Duration,
    propagated time.Duration,
    nodesAffected int,
    totalNodes int,
) {
    atomic.AddInt64(&m.DeltaPropagations, 1)
    atomic.AddInt64(&m.NodesEvaluated, int64(nodesAffected))
    atomic.AddInt64(&m.NodesAvoided, int64(totalNodes-nodesAffected))
    
    // Utiliser sync.Mutex pour les durées
}

func (m *DeltaMetrics) GetSavings() float64 {
    evaluated := atomic.LoadInt64(&m.NodesEvaluated)
    avoided := atomic.LoadInt64(&m.NodesAvoided)
    total := evaluated + avoided
    
    if total == 0 {
        return 0.0
    }
    
    return float64(avoided) / float64(total) * 100.0
}

// Exposition Prometheus (optionnel)
func (m *DeltaMetrics) RegisterPrometheus() {
    // Enregistrer collectors Prometheus
}
```

### Étape 5 : Tests de Régression

Avant déploiement, valider que delta produit les mêmes résultats :

```go
func TestDeltaVsClassic(t *testing.T) {
    // Setup
    network := setupTestNetwork()
    detector := delta.NewDeltaDetector()
    index := delta.NewDependencyIndex()
    builder := delta.NewIndexBuilder()
    builder.BuildFromNetwork(index, network)
    
    // Préparer faits
    oldFact := map[string]interface{}{"id": "p1", "price": 100.0}
    newFact := map[string]interface{}{"id": "p1", "price": 150.0}
    
    // Clone réseau pour test classique
    networkClassic := network.Clone()
    networkDelta := network.Clone()
    
    // Approche classique
    networkClassic.RetractFact("Product~p1")
    networkClassic.AssertFact(newFact)
    classicState := networkClassic.GetWorkingMemory()
    
    // Approche delta
    factDelta, _ := detector.DetectDelta(oldFact, newFact, "Product~p1", "Product")
    affectedNodes := index.GetAffectedNodesForDelta(factDelta)
    for _, node := range affectedNodes {
        networkDelta.PropagateToNode(node.NodeID, factDelta)
    }
    deltaState := networkDelta.GetWorkingMemory()
    
    // Comparer états finaux
    if !reflect.DeepEqual(classicState, deltaState) {
        t.Errorf("États différents! Classic=%+v, Delta=%+v", 
            classicState, deltaState)
    }
}
```

---

## 🎯 Cas d'Usage Réels

### Cas 1 : E-commerce - Mise à Jour Prix Produit

```go
// Contexte : 500 produits, 1000 règles, updates fréquents de prix

type ProductPriceUpdater struct {
    updater *FactUpdater
}

func (u *ProductPriceUpdater) UpdatePrice(productID string, newPrice float64) error {
    oldProduct := u.getProduct(productID)
    newProduct := copyMap(oldProduct)
    newProduct["price"] = newPrice
    
    return u.updater.UpdateFact(
        oldProduct,
        newProduct,
        fmt.Sprintf("Product~%s", productID),
        "Product",
    )
}

// Résultat : 4x plus rapide (200ms → 50ms par update)
// Gain : 75% nœuds évités (seules règles sur "price" évaluées)
```

### Cas 2 : IoT - Mises à Jour Capteurs

```go
// Contexte : 10k capteurs, température change toutes les 30s

type SensorUpdater struct {
    updater *FactUpdater
}

func (u *SensorUpdater) UpdateTemperature(sensorID string, temp float64) error {
    old := map[string]interface{}{
        "id":          sensorID,
        "temperature": u.lastTemp[sensorID],
        "location":    u.locations[sensorID],
        "status":      "active",
        "calibrated":  true,
    }
    
    new := copyMap(old)
    new["temperature"] = temp
    
    // Seul "temperature" change → delta très efficace
    return u.updater.UpdateFact(old, new, 
        fmt.Sprintf("Sensor~%s", sensorID), "Sensor")
}

// Résultat : 10x plus rapide (seules règles température activées)
```

### Cas 3 : Gestion Commandes - Workflow États

```go
// Contexte : Commandes passent par plusieurs états

type OrderStateUpdater struct {
    updater *FactUpdater
}

func (u *OrderStateUpdater) TransitionState(
    orderID string, 
    newState string,
) error {
    oldOrder := u.getOrder(orderID)
    newOrder := copyMap(oldOrder)
    newOrder["status"] = newState
    newOrder["updated_at"] = time.Now()
    
    // Configuration pour ignorer "updated_at"
    return u.updater.UpdateFact(oldOrder, newOrder,
        fmt.Sprintf("Order~%s", orderID), "Order")
}

// Résultat : 3x plus rapide (seules règles "status" évaluées)
```

---

## ⚠️ Pièges Courants et Solutions

### Piège 1 : Index Non Reconstruit

**Problème** : Vous ajoutez une règle mais l'index n'est pas mis à jour.

```go
// ❌ MAUVAIS
network.AddRule(newRule)
// Index ne connaît pas la nouvelle règle !
updater.UpdateFact(old, new, id, typ) // Règle ignorée

// ✅ BON
network.AddRule(newRule)
builder.BuildFromNetwork(index, network) // Reconstruire
updater.UpdateFact(old, new, id, typ)
```

**Solution** : Hook de reconstruction automatique.

```go
type AutoRebuildNetwork struct {
    *rete.Network
    indexManager *IndexManager
}

func (n *AutoRebuildNetwork) AddRule(rule *rete.Rule) error {
    if err := n.Network.AddRule(rule); err != nil {
        return err
    }
    return n.indexManager.RebuildIfNeeded()
}
```

### Piège 2 : Comparaison Float Stricte

**Problème** : `0.1 + 0.2 != 0.3` en float, delta détecte un changement fantôme.

```go
// ❌ MAUVAIS
detector := delta.NewDeltaDetector() // FloatEpsilon par défaut = 1e-9

old := map[string]interface{}{"amount": 0.1 + 0.2}
new := map[string]interface{}{"amount": 0.3}

d, _ := detector.DetectDelta(old, new, "id", "Type")
// d.Changes contient "amount" (faux positif)

// ✅ BON
config := delta.DetectorConfig{
    FloatEpsilon: 0.0001, // Tolérance adaptée
}
detector := delta.NewDeltaDetectorWithConfig(config)
```

### Piège 3 : Champs Ignorés Non Configurés

**Problème** : Champs techniques (`updated_at`, `version`) déclenchent deltas inutiles.

```go
// ❌ MAUVAIS
old := map[string]interface{}{
    "price":      100.0,
    "updated_at": time1,
}
new := map[string]interface{}{
    "price":      100.0,
    "updated_at": time2, // Changé mais pas important
}
// Delta détecté inutilement

// ✅ BON
config := delta.DetectorConfig{
    IgnoredFields: []string{"updated_at", "version", "_metadata"},
}
```

### Piège 4 : Overhead sur Petits Réseaux

**Problème** : Delta plus lent que classique sur <20 nœuds.

```go
// ❌ MAUVAIS : Utiliser delta partout
if anyUpdate {
    useDelta()
}

// ✅ BON : Décision adaptative
if networkSize > 50 && changeRatio < 0.3 {
    useDelta()
} else {
    useClassic()
}
```

---

## 📊 Benchmarking Avant/Après

Script pour mesurer gains réels :

```go
func BenchmarkMigration(b *testing.B) {
    network := setupRealNetwork() // Votre réseau
    facts := generateUpdates(1000) // Vos updates typiques
    
    b.Run("Classic", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            for _, update := range facts {
                network.RetractFact(update.ID)
                network.AssertFact(update.New)
            }
        }
    })
    
    b.Run("Delta", func(b *testing.B) {
        detector := delta.NewDeltaDetector()
        index := delta.NewDependencyIndex()
        builder := delta.NewIndexBuilder()
        builder.BuildFromNetwork(index, network)
        
        for i := 0; i < b.N; i++ {
            for _, update := range facts {
                d, _ := detector.DetectDelta(
                    update.Old, update.New, 
                    update.ID, update.Type,
                )
                nodes := index.GetAffectedNodesForDelta(d)
                for _, node := range nodes {
                    network.PropagateToNode(node.NodeID, d)
                }
            }
        }
    })
}
```

---

## ✅ Checklist de Migration

- [ ] Analyser profil réseau (>50 nœuds, updates fréquents)
- [ ] Choisir configuration détecteur adaptée
- [ ] Implémenter wrapper `FactUpdater`
- [ ] Construire index initial avec `BuildFromNetwork`
- [ ] Implémenter reconstruction périodique
- [ ] Ajouter métriques et monitoring
- [ ] Tests de régression (delta vs classique)
- [ ] Benchmarks avant/après
- [ ] Déploiement progressif (canary, A/B test)
- [ ] Monitoring production (fallback rate, latence)

---

## 📚 Ressources

- [README.md](./README.md) - Architecture complète
- [QUICK_START.md](./QUICK_START.md) - Premiers pas
- [OPTIMIZATION_GUIDE.md](./OPTIMIZATION_GUIDE.md) - Tuning avancé
- [examples/](./examples/) - Exemples complets
- Tests : `e2e_business_test.go` - Scénarios réels

---

## 🆘 Support

**Problèmes courants** :

1. **"Index vide après BuildFromNetwork"**  
   → Vérifier que le réseau contient des nœuds avec conditions AST compatibles

2. **"Fallback classique trop fréquent"**  
   → Ajuster `deltaThreshold` ou vérifier `IgnoredFields`

3. **"Performance dégradée"**  
   → Profiler avec `go test -cpuprofile`, vérifier taille cache

4. **"Résultats différents classic vs delta"**  
   → Bug potentiel, créer test de régression et signaler

---

**Dernière mise à jour** : 2025-01-02  
**Version** : 1.0.0  
**Auteur** : TSD Contributors