# Guide d'Utilisation des Optimisations - Delta Package

Ce guide explique comment utiliser les optimisations implémentées dans le package `rete/delta`.

---

## 📊 Vue d'Ensemble

Les optimisations sont organisées en 4 catégories:

1. **Object Pooling**: Réutilisation d'objets pour réduire allocations
2. **Cache LRU**: Cache haute performance avec éviction automatique
3. **Fast Paths**: Optimisations de comparaisons
4. **Batch Processing**: Traitement groupé de nœuds

---

## 🔧 1. Object Pooling

### Pools Disponibles

#### FactDelta Pool

**Quand utiliser**: Création fréquente de deltas (> 1000/sec)

```go
// Acquérir un FactDelta depuis le pool
delta := AcquireFactDelta(factID, factType)
defer ReleaseFactDelta(delta) // IMPORTANT: Toujours release

// Utiliser normalement
delta.AddFieldChange("price", 100.0, 150.0)

// Le defer s'occupe du release automatiquement
```

**⚠️ Important**: 
- TOUJOURS appeler `ReleaseFactDelta()` après usage
- SAUF si le delta est mis en cache ou stocké ailleurs
- Utiliser `defer` pour éviter les oublis

#### NodeReference Slice Pool

**Quand utiliser**: Manipulation fréquente de listes de nœuds

```go
// Acquérir une slice depuis le pool
slice := AcquireNodeReferenceSlice()
defer ReleaseNodeReferenceSlice(slice)

// Utiliser
*slice = append(*slice, NodeReference{NodeID: "node1"})

// Le defer s'occupe du release
```

#### StringBuilder Pool

**Quand utiliser**: Construction fréquente de strings

```go
// Acquérir un builder
sb := AcquireStringBuilder()
defer ReleaseStringBuilder(sb)

// Utiliser
sb.WriteString("prefix_")
sb.WriteString(nodeID)
result := sb.String()
```

#### Map Pool

**Quand utiliser**: Besoins temporaires de maps

```go
// Acquérir une map
m := AcquireMap()
defer ReleaseMap(m)

// Utiliser
(*m)["key"] = "value"
```

### Gains de Performance

```
Sans Pool:  586 ns/op   832 B/op   3 allocs/op
Avec Pool:  354 ns/op     0 B/op   0 allocs/op
Gain:       -40%         -100%     -100%
```

---

## 💾 2. Cache LRU Optimisé

### Création

```go
import "time"

// Créer un cache
// Args: maxSize (nombre d'entrées), TTL (durée de vie)
cache := NewOptimizedCache(1000, 5*time.Minute)
```

### Utilisation

```go
// Put
delta := NewFactDelta("Product~123", "Product")
cache.Put("cache_key", delta)

// Get
if result, found := cache.Get("cache_key"); found {
    // Utiliser result
    fmt.Println(result.FactID)
} else {
    // Cache miss
}

// Obtenir statistiques
stats := cache.GetStats()
fmt.Printf("Hit rate: %.2f%%\n", stats.HitRate*100)
fmt.Printf("Size: %d/%d\n", stats.Size, maxSize)
fmt.Printf("Evictions: %d\n", stats.Evictions)
```

### Configuration Recommandée

| Use Case | Max Size | TTL | Raison |
|----------|----------|-----|--------|
| **Dev/Test** | 100 | 1 min | Petit footprint, rotation rapide |
| **Production Light** | 1000 | 5 min | Équilibre hit rate / mémoire |
| **Production Heavy** | 10000 | 10 min | Haute fréquence, gros volume |
| **Real-time** | 5000 | 30 sec | Données fraîches prioritaires |

### Monitoring

```go
// Périodiquement (ex: toutes les minutes)
stats := cache.GetStats()

// Alerting sur anomalies
if stats.HitRate < 0.5 {
    log.Warn("Cache hit rate too low: %.2f", stats.HitRate)
}

if stats.Evictions > 1000 {
    log.Warn("Too many evictions: %d (consider increasing size)", stats.Evictions)
}
```

### Performances

```
Get:    ~40 ns/op     0 allocs
Put:    ~60 ns/op     0 allocs
Evict:  ~50 ns/op     0 allocs
```

---

## ⚡ 3. Fast Paths (Comparaisons Optimisées)

### OptimizedValuesEqual

**Quand utiliser**: Comparaisons haute fréquence de valeurs

```go
// Standard (utilise reflect pour tous les types)
equal := ValuesEqual(a, b, epsilon)

// Optimisé (fast path pour types simples)
equal := OptimizedValuesEqual(a, b, epsilon)
```

**Gains**: 
- Types simples (int, string, bool): ~2ns (équivalent)
- Évite `reflect.TypeOf` pour 95% des cas
- Fallback automatique à `ValuesEqual` pour types complexes

### CopyFactFast

**Quand utiliser**: Copie rapide de faits

```go
original := map[string]interface{}{
    "field1": "value1",
    "field2": 42,
}

// Copie optimisée (pré-allocation)
copy := CopyFactFast(original)
```

### FastHashString

**Quand utiliser**: Hash non-cryptographique rapide (cache keys, etc.)

```go
// FNV-1a hash
hash := FastHashString("my_cache_key_123")
```

---

## 📦 4. Batch Processing

### BatchNodeReferences

**Quand utiliser**: Traitement ordonné de nœuds par type

```go
// Créer un batch
batch := NewBatchNodeReferences(100)

// Ajouter des nœuds
for _, node := range allNodes {
    batch.Add(node)
}

// Traiter dans l'ordre: Alpha → Beta → Terminal
err := batch.ProcessInOrder(func(ref NodeReference) error {
    // Traiter le nœud
    return propagateToNode(ref)
})
```

**Avantages**:
- Ordre optimal pour propagation RETE
- Meilleure localité cache
- Groupage logique par type

---

## 🎯 Patterns d'Utilisation Recommandés

### Pattern 1: Détection Delta avec Pool

```go
func DetectAndPropagate(oldFact, newFact map[string]interface{}) error {
    // Acquérir depuis pool
    delta := AcquireFactDelta(factID, factType)
    defer ReleaseFactDelta(delta)
    
    // Détecter changements
    detector.PopulateDelta(delta, oldFact, newFact)
    
    // Propager
    return propagator.Propagate(delta)
}
```

### Pattern 2: Cache avec Fallback

```go
func GetOrComputeDelta(key string, oldFact, newFact map[string]interface{}) (*FactDelta, error) {
    // Vérifier cache
    if cached, found := cache.Get(key); found {
        return cached, nil
    }
    
    // Cache miss - calculer
    delta, err := detector.DetectDelta(oldFact, newFact, factID, factType)
    if err != nil {
        return nil, err
    }
    
    // Mettre en cache (ne PAS release car en cache)
    cache.Put(key, delta)
    
    return delta, nil
}
```

### Pattern 3: Batch Processing Optimisé

```go
func PropagateToNodes(delta *FactDelta, nodes []NodeReference) error {
    // Créer batch
    batch := NewBatchNodeReferences(len(nodes))
    
    // Grouper par type
    for _, node := range nodes {
        batch.Add(node)
    }
    
    // Traiter dans ordre optimal
    return batch.ProcessInOrder(func(ref NodeReference) error {
        return propagateToNode(ref, delta)
    })
}
```

---

## 📊 Benchmarking

### Exécuter Benchmarks

```bash
# Tous les benchmarks
go test ./rete/delta/... -bench=. -benchmem

# Benchmarks pooling
go test ./rete/delta/... -bench=BenchmarkPool -benchmem

# Benchmarks scalabilité
go test ./rete/delta/... -bench=BenchmarkScalability -benchmem

# Benchmarks cache
go test ./rete/delta/... -bench=BenchmarkOptimizedCache -benchmem
```

### Scripts Automatisés

```bash
# Profiling complet (CPU, mémoire, trace)
./scripts/profile_delta.sh

# Benchmarks avec aggregation
./scripts/benchmark_delta.sh
```

---

## 🔍 Profiling

### CPU Profiling

```bash
cd rete/delta
go test -bench=. -cpuprofile=cpu.prof -benchtime=10s
go tool pprof -http=:8080 cpu.prof
```

### Memory Profiling

```bash
go test -bench=. -memprofile=mem.prof -benchtime=10s
go tool pprof -http=:8081 mem.prof
```

### Trace Analysis

```bash
go test -bench=BenchmarkDeltaDetector -trace=trace.out
go tool trace trace.out
```

---

## ⚠️ Pièges à Éviter

### 1. Oublier ReleaseFactDelta

```go
// ❌ MAUVAIS - Fuite mémoire potentielle
delta := AcquireFactDelta(id, typ)
// ... utiliser delta
// Oubli de release!

// ✅ BON - Toujours release
delta := AcquireFactDelta(id, typ)
defer ReleaseFactDelta(delta)
```

### 2. Release avec Cache

```go
// ❌ MAUVAIS - Delta en cache puis release
delta := AcquireFactDelta(id, typ)
cache.Put(key, delta)
ReleaseFactDelta(delta) // delta en cache sera invalidé!

// ✅ BON - Ne pas release si mis en cache
delta := AcquireFactDelta(id, typ)
cache.Put(key, delta)
// Pas de release - le cache gère le lifecycle
```

### 3. Pool sur Objets Larges

```go
// ❌ MAUVAIS - Objet trop large pour pool
slice := AcquireNodeReferenceSlice()
for i := 0; i < 10000; i++ { // Trop large!
    *slice = append(*slice, node)
}
ReleaseNodeReferenceSlice(slice) // Sera rejeté car trop large

// ✅ BON - Utiliser make directement pour gros objets
largeSlice := make([]NodeReference, 0, 10000)
```

---

## 📈 Métriques de Succès

### Objectifs de Performance

| Métrique | Objectif | Actuel | Status |
|----------|----------|--------|--------|
| DetectDelta (1 champ) | < 500 ns | ~350 ns | ✅ |
| Cache hit latency | < 50 ns | ~40 ns | ✅ |
| Pool overhead | < 30 ns | ~25 ns | ✅ |
| Allocations avec pool | 0 | 0 | ✅ |
| Cache hit rate | > 80% | Variable | 📊 |

### Monitoring Production

```go
// Exemple metrics exposées
type DeltaMetrics struct {
    PoolAcquisitions  int64
    PoolReleases      int64
    CacheHits         int64
    CacheMisses       int64
    CacheEvictions    int64
    DetectorCalls     int64
    AvgDeltaSize      float64
}
```

---

## 🔗 Ressources

- **Rapport Performance**: `REPORTS/delta_performance_report.md`
- **TODO**: `REPORTS/delta_optimizations_TODO.md`
- **Tests**: `rete/delta/*_test.go`
- **Benchmarks**: `rete/delta/benchmark_*_test.go`

---

**Dernière mise à jour**: 2026-01-02  
**Version**: 1.0  
**Contact**: TSD Team
