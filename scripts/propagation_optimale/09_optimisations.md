# ⚡ Prompt 09 - Optimisations et Profiling

> **📋 Standards** : Ce prompt respecte les règles de `.github/prompts/common.md` et `.github/prompts/develop.md`

## 🎯 Objectif

Analyser les performances du système de propagation delta, identifier les goulots d'étranglement, et implémenter les optimisations nécessaires pour atteindre les objectifs de performance.

**Objectifs de performance** :
- Update partiel (1 champ) : **> 10x plus rapide** que Retract+Insert
- Overhead index : **< 5%** mémoire
- Latency p99 : **< 1ms** pour propagation delta
- Throughput : **> 10000 updates/sec**

**⚠️ IMPORTANT** : Ce prompt optimise du code existant. Respecter strictement les standards de `common.md`.

---

## 📋 Prérequis

Avant de commencer ce prompt :

- [x] **Prompts 01-08 validés** : Système complet + tests passent
- [x] **Tests d'intégration passent** : 100% success
- [x] **Documents de référence** :
  - `REPORTS/conception_delta_architecture.md`
  - Résultats benchmarks des prompts précédents

---

## 📂 Fichiers à Créer/Modifier

```
rete/delta/
├── optimizations.go              # Optimisations (nouveau)
├── pool.go                       # Object pooling (nouveau)
├── cache_optimized.go            # Cache optimisé (nouveau)
└── [fichiers existants]          # Modifications optimisations

scripts/
├── profile_delta.sh              # Script profiling (nouveau)
└── benchmark_delta.sh            # Script benchmarks (nouveau)

REPORTS/
└── delta_performance_report.md  # Rapport performance (nouveau)
```

---

## 🔧 Tâche 1 : Profiling et Analyse

### Script : `scripts/profile_delta.sh`

**Contenu** :

```bash
#!/bin/bash
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License

set -e

echo "🔍 Profiling Delta Propagation System"
echo "======================================"

# 1. CPU Profiling
echo ""
echo "📊 CPU Profiling..."
go test ./rete/delta/... -bench=. -cpuprofile=cpu.prof -benchtime=10s
go tool pprof -http=:8080 cpu.prof &

# 2. Memory Profiling
echo ""
echo "💾 Memory Profiling..."
go test ./rete/delta/... -bench=. -memprofile=mem.prof -benchtime=10s
go tool pprof -http=:8081 mem.prof &

# 3. Allocation Profiling
echo ""
echo "🔢 Allocation Analysis..."
go test ./rete/delta/... -bench=. -benchmem -benchtime=5s > alloc_report.txt

# 4. Trace Analysis
echo ""
echo "🔬 Trace Analysis..."
go test ./rete/delta/... -bench=BenchmarkDeltaDetector -trace=trace.out
go tool trace trace.out &

# 5. Escape Analysis
echo ""
echo "🏃 Escape Analysis..."
go build -gcflags='-m -m' ./rete/delta/... 2> escape_analysis.txt

echo ""
echo "✅ Profiling complete!"
echo "   - CPU profile: http://localhost:8080"
echo "   - Memory profile: http://localhost:8081"
echo "   - Allocations: alloc_report.txt"
echo "   - Escape analysis: escape_analysis.txt"
```

### Script : `scripts/benchmark_delta.sh`

**Contenu** :

```bash
#!/bin/bash
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License

set -e

echo "⚡ Benchmarking Delta Propagation System"
echo "========================================"

ITERATIONS=5
OUTPUT_DIR="benchmark_results"

mkdir -p $OUTPUT_DIR

echo ""
echo "🏃 Running benchmarks ($ITERATIONS iterations)..."

for i in $(seq 1 $ITERATIONS); do
    echo "  Iteration $i/$ITERATIONS..."
    go test ./rete/delta/... \
        -bench=. \
        -benchmem \
        -benchtime=5s \
        -count=1 \
        > "$OUTPUT_DIR/bench_$i.txt"
done

echo ""
echo "📊 Aggregating results..."

# Utiliser benchstat si disponible
if command -v benchstat &> /dev/null; then
    benchstat $OUTPUT_DIR/bench_*.txt > $OUTPUT_DIR/aggregate.txt
    cat $OUTPUT_DIR/aggregate.txt
else
    echo "⚠️  Install benchstat for detailed analysis: go install golang.org/x/perf/cmd/benchstat@latest"
fi

echo ""
echo "✅ Benchmarks complete! Results in: $OUTPUT_DIR/"
```

### Analyse : Identifier les Goulots d'Étranglement

Exécuter les scripts et analyser :

```bash
chmod +x scripts/profile_delta.sh scripts/benchmark_delta.sh
./scripts/profile_delta.sh
./scripts/benchmark_delta.sh
```

**Points à analyser** :

1. **CPU hotspots** : Fonctions consommant le plus de CPU
2. **Allocations** : Fonctions allouant le plus de mémoire
3. **Escape analysis** : Variables s'échappant vers le heap
4. **Lock contention** : Temps passé à attendre des locks

**Documenter dans** : `REPORTS/delta_profiling_analysis.md`

---

## 🔧 Tâche 2 : Object Pooling

### Fichier : `rete/delta/pool.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
    "sync"
)

// FactDeltaPool est un pool d'objets FactDelta pour réduire les allocations.
var FactDeltaPool = sync.Pool{
    New: func() interface{} {
        return &FactDelta{
            Fields: make(map[string]FieldDelta, 8), // Pré-allouer taille typique
        }
    },
}

// AcquireFactDelta obtient un FactDelta depuis le pool.
func AcquireFactDelta(factID, factType string) *FactDelta {
    delta := FactDeltaPool.Get().(*FactDelta)
    delta.FactID = factID
    delta.FactType = factType
    delta.Timestamp = time.Now()
    return delta
}

// ReleaseFactDelta retourne un FactDelta au pool.
func ReleaseFactDelta(delta *FactDelta) {
    // Reset pour réutilisation
    delta.FactID = ""
    delta.FactType = ""
    delta.FieldCount = 0
    
    // Clear map sans réallouer
    for k := range delta.Fields {
        delete(delta.Fields, k)
    }
    
    FactDeltaPool.Put(delta)
}

// NodeReferencePool est un pool de slices de NodeReference.
var NodeReferencePool = sync.Pool{
    New: func() interface{} {
        slice := make([]NodeReference, 0, 16)
        return &slice
    },
}

// AcquireNodeReferenceSlice obtient une slice depuis le pool.
func AcquireNodeReferenceSlice() *[]NodeReference {
    slice := NodeReferencePool.Get().(*[]NodeReference)
    *slice = (*slice)[:0] // Reset length, keep capacity
    return slice
}

// ReleaseNodeReferenceSlice retourne une slice au pool.
func ReleaseNodeReferenceSlice(slice *[]NodeReference) {
    if cap(*slice) > 1024 {
        // Trop grande, ne pas réutiliser
        return
    }
    NodeReferencePool.Put(slice)
}

// StringBuilderPool est un pool de strings.Builder pour construction efficace.
var StringBuilderPool = sync.Pool{
    New: func() interface{} {
        return &strings.Builder{}
    },
}

// AcquireStringBuilder obtient un builder depuis le pool.
func AcquireStringBuilder() *strings.Builder {
    sb := StringBuilderPool.Get().(*strings.Builder)
    sb.Reset()
    return sb
}

// ReleaseStringBuilder retourne un builder au pool.
func ReleaseStringBuilder(sb *strings.Builder) {
    if sb.Cap() > 4096 {
        // Trop grand, ne pas réutiliser
        return
    }
    StringBuilderPool.Put(sb)
}
```

**Modification de `delta_detector.go`** :

```go
// Utiliser le pool dans DetectDelta
func (dd *DeltaDetector) DetectDelta(
    oldFact, newFact map[string]interface{},
    factID, factType string,
) (*FactDelta, error) {
    dd.incrementComparisons()
    
    // Acquérir depuis le pool au lieu de NewFactDelta
    delta := AcquireFactDelta(factID, factType)
    delta.FieldCount = len(newFact)
    
    // ... reste de la logique ...
    
    // Note : L'appelant doit appeler ReleaseFactDelta(delta) quand fini
    
    return delta, nil
}
```

---

## 🔧 Tâche 3 : Optimisations Cache

### Fichier : `rete/delta/cache_optimized.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
    "sync"
    "time"
)

// OptimizedCache est un cache haute performance avec éviction LRU.
type OptimizedCache struct {
    entries    map[string]*cacheEntry
    lruList    *lruList
    maxSize    int
    ttl        time.Duration
    hits       int64
    misses     int64
    evictions  int64
    mutex      sync.RWMutex
}

// lruList est une liste doublement chaînée pour LRU.
type lruList struct {
    head *lruNode
    tail *lruNode
    size int
}

type lruNode struct {
    key   string
    entry *cacheEntry
    prev  *lruNode
    next  *lruNode
}

// cacheEntry représente une entrée de cache.
type cacheEntry struct {
    delta     *FactDelta
    createdAt time.Time
    accessedAt time.Time
    accessCount int64
}

// NewOptimizedCache crée un cache optimisé.
func NewOptimizedCache(maxSize int, ttl time.Duration) *OptimizedCache {
    return &OptimizedCache{
        entries: make(map[string]*cacheEntry, maxSize),
        lruList: &lruList{},
        maxSize: maxSize,
        ttl:     ttl,
    }
}

// Get récupère une valeur depuis le cache.
func (oc *OptimizedCache) Get(key string) (*FactDelta, bool) {
    oc.mutex.Lock()
    defer oc.mutex.Unlock()
    
    entry, exists := oc.entries[key]
    if !exists {
        oc.misses++
        return nil, false
    }
    
    // Vérifier expiration
    if time.Since(entry.createdAt) > oc.ttl {
        oc.remove(key)
        oc.misses++
        return nil, false
    }
    
    // Mettre à jour LRU
    oc.lruList.moveToFront(key)
    entry.accessedAt = time.Now()
    entry.accessCount++
    
    oc.hits++
    return entry.delta, true
}

// Put ajoute une valeur au cache.
func (oc *OptimizedCache) Put(key string, delta *FactDelta) {
    oc.mutex.Lock()
    defer oc.mutex.Unlock()
    
    // Si existe déjà, mettre à jour
    if entry, exists := oc.entries[key]; exists {
        entry.delta = delta
        entry.accessedAt = time.Now()
        oc.lruList.moveToFront(key)
        return
    }
    
    // Si plein, évincer LRU
    if len(oc.entries) >= oc.maxSize {
        oc.evictLRU()
    }
    
    // Ajouter nouvelle entrée
    entry := &cacheEntry{
        delta:      delta,
        createdAt:  time.Now(),
        accessedAt: time.Now(),
        accessCount: 0,
    }
    
    oc.entries[key] = entry
    oc.lruList.addToFront(key, entry)
}

// evictLRU évince l'entrée la moins récemment utilisée.
func (oc *OptimizedCache) evictLRU() {
    if oc.lruList.tail == nil {
        return
    }
    
    key := oc.lruList.tail.key
    oc.remove(key)
    oc.evictions++
}

// remove supprime une entrée du cache.
func (oc *OptimizedCache) remove(key string) {
    delete(oc.entries, key)
    oc.lruList.remove(key)
}

// GetStats retourne les statistiques du cache.
func (oc *OptimizedCache) GetStats() struct {
    Size       int
    Hits       int64
    Misses     int64
    Evictions  int64
    HitRate    float64
} {
    oc.mutex.RLock()
    defer oc.mutex.RUnlock()
    
    total := oc.hits + oc.misses
    hitRate := 0.0
    if total > 0 {
        hitRate = float64(oc.hits) / float64(total)
    }
    
    return struct {
        Size       int
        Hits       int64
        Misses     int64
        Evictions  int64
        HitRate    float64
    }{
        Size:      len(oc.entries),
        Hits:      oc.hits,
        Misses:    oc.misses,
        Evictions: oc.evictions,
        HitRate:   hitRate,
    }
}

// Implémentation lruList

func (l *lruList) addToFront(key string, entry *cacheEntry) {
    node := &lruNode{
        key:   key,
        entry: entry,
    }
    
    if l.head == nil {
        l.head = node
        l.tail = node
    } else {
        node.next = l.head
        l.head.prev = node
        l.head = node
    }
    
    l.size++
}

func (l *lruList) moveToFront(key string) {
    node := l.find(key)
    if node == nil || node == l.head {
        return
    }
    
    // Détacher
    if node.prev != nil {
        node.prev.next = node.next
    }
    if node.next != nil {
        node.next.prev = node.prev
    }
    if node == l.tail {
        l.tail = node.prev
    }
    
    // Mettre en tête
    node.prev = nil
    node.next = l.head
    l.head.prev = node
    l.head = node
}

func (l *lruList) remove(key string) {
    node := l.find(key)
    if node == nil {
        return
    }
    
    if node.prev != nil {
        node.prev.next = node.next
    } else {
        l.head = node.next
    }
    
    if node.next != nil {
        node.next.prev = node.prev
    } else {
        l.tail = node.prev
    }
    
    l.size--
}

func (l *lruList) find(key string) *lruNode {
    current := l.head
    for current != nil {
        if current.key == key {
            return current
        }
        current = current.next
    }
    return nil
}
```

---

## 🔧 Tâche 4 : Optimisations Diverses

### Fichier : `rete/delta/optimizations.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

// OptimizedValuesEqual est une version optimisée de ValuesEqual.
//
// Optimisations :
// - Short-circuit pour types simples
// - Évite reflect.DeepEqual quand possible
// - Inline hints pour compiler
func OptimizedValuesEqual(a, b interface{}, epsilon float64) bool {
    // Fast path : égalité stricte (très fréquent)
    if a == b {
        return true
    }
    
    // Fast path : nil
    if a == nil || b == nil {
        return false
    }
    
    // Fast path : types simples (inlined par compiler)
    switch va := a.(type) {
    case string:
        vb, ok := b.(string)
        return ok && va == vb
    
    case int:
        vb, ok := b.(int)
        return ok && va == vb
    
    case int64:
        vb, ok := b.(int64)
        return ok && va == vb
    
    case bool:
        vb, ok := b.(bool)
        return ok && va == vb
    
    case float64:
        vb, ok := b.(float64)
        if !ok {
            return false
        }
        // Comparaison float optimisée
        diff := va - vb
        if diff < 0 {
            diff = -diff
        }
        return diff <= epsilon
    
    case float32:
        vb, ok := b.(float32)
        if !ok {
            return false
        }
        diff := va - vb
        if diff < 0 {
            diff = -diff
        }
        return float64(diff) <= epsilon
    }
    
    // Slow path : comparaison générique
    return ValuesEqual(a, b, epsilon)
}

// FastHashString calcule un hash rapide d'une string (non-cryptographique).
func FastHashString(s string) uint64 {
    // FNV-1a hash
    const (
        offset64 = 14695981039346656037
        prime64  = 1099511628211
    )
    
    hash := uint64(offset64)
    for i := 0; i < len(s); i++ {
        hash ^= uint64(s[i])
        hash *= prime64
    }
    return hash
}

// PreallocatedMap crée une map pré-allouée avec capacité.
func PreallocatedMap(size int) map[string]interface{} {
    return make(map[string]interface{}, size)
}

// CopyFactFast copie rapidement un fait (optimisé).
func CopyFactFast(fact map[string]interface{}) map[string]interface{} {
    if len(fact) == 0 {
        return make(map[string]interface{})
    }
    
    copy := make(map[string]interface{}, len(fact))
    for k, v := range fact {
        copy[k] = v
    }
    return copy
}

// BatchNodeReferences groupe les références de nœuds pour traitement batch.
type BatchNodeReferences struct {
    alpha    []NodeReference
    beta     []NodeReference
    terminal []NodeReference
}

// NewBatchNodeReferences crée un batch pré-alloué.
func NewBatchNodeReferences(expectedSize int) *BatchNodeReferences {
    return &BatchNodeReferences{
        alpha:    make([]NodeReference, 0, expectedSize/3),
        beta:     make([]NodeReference, 0, expectedSize/3),
        terminal: make([]NodeReference, 0, expectedSize/3),
    }
}

// Add ajoute une référence au batch approprié.
func (b *BatchNodeReferences) Add(ref NodeReference) {
    switch ref.NodeType {
    case "alpha":
        b.alpha = append(b.alpha, ref)
    case "beta":
        b.beta = append(b.beta, ref)
    case "terminal":
        b.terminal = append(b.terminal, ref)
    }
}

// ProcessInOrder traite les nœuds dans l'ordre optimal.
func (b *BatchNodeReferences) ProcessInOrder(processor func(NodeReference) error) error {
    // Alpha d'abord
    for _, ref := range b.alpha {
        if err := processor(ref); err != nil {
            return err
        }
    }
    
    // Beta ensuite
    for _, ref := range b.beta {
        if err := processor(ref); err != nil {
            return err
        }
    }
    
    // Terminal en dernier
    for _, ref := range b.terminal {
        if err := processor(ref); err != nil {
            return err
        }
    }
    
    return nil
}
```

---

## 🔧 Tâche 5 : Benchmarks Comparatifs

### Fichier : `rete/delta/benchmark_comparison_test.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
    "testing"
)

// BenchmarkComparison_DeltaVsClassic compare delta vs classique.
func BenchmarkComparison_DeltaVsClassic(b *testing.B) {
    scenarios := []struct {
        name         string
        fieldCount   int
        changedCount int
    }{
        {"1_field_of_10", 10, 1},
        {"2_fields_of_10", 10, 2},
        {"5_fields_of_10", 10, 5},
        {"1_field_of_100", 100, 1},
        {"10_fields_of_100", 100, 10},
        {"50_fields_of_100", 100, 50},
    }
    
    for _, scenario := range scenarios {
        oldFact := generateFact(scenario.fieldCount)
        newFact := modifyNFields(oldFact, scenario.changedCount)
        
        b.Run(scenario.name+"_delta", func(b *testing.B) {
            detector := NewDeltaDetector()
            propagator := setupPropagator()
            
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                delta, _ := detector.DetectDelta(oldFact, newFact, "Test~1", "Test")
                _ = propagator.PropagateUpdate(oldFact, newFact, "Test~1", "Test")
                ReleaseFactDelta(delta)
            }
        })
        
        b.Run(scenario.name+"_classic", func(b *testing.B) {
            network := setupClassicNetwork()
            
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                _ = network.RetractFact("Test~1", "Test")
                _ = network.InsertFact(newFact, "Test~1", "Test")
            }
        })
    }
}

// BenchmarkMemoryUsage compare l'utilisation mémoire.
func BenchmarkMemoryUsage_DeltaVsClassic(b *testing.B) {
    b.Run("delta", func(b *testing.B) {
        detector := NewDeltaDetector()
        oldFact := generateFact(50)
        newFact := modifyNFields(oldFact, 5)
        
        b.ReportAllocs()
        b.ResetTimer()
        
        for i := 0; i < b.N; i++ {
            delta, _ := detector.DetectDelta(oldFact, newFact, "Test~1", "Test")
            ReleaseFactDelta(delta)
        }
    })
    
    b.Run("classic", func(b *testing.B) {
        network := setupClassicNetwork()
        newFact := generateFact(50)
        
        b.ReportAllocs()
        b.ResetTimer()
        
        for i := 0; i < b.N; i++ {
            _ = network.RetractFact("Test~1", "Test")
            _ = network.InsertFact(newFact, "Test~1", "Test")
        }
    })
}

// Helpers
func generateFact(fieldCount int) map[string]interface{} {
    fact := make(map[string]interface{}, fieldCount)
    for i := 0; i < fieldCount; i++ {
        fact[fmt.Sprintf("field%d", i)] = i
    }
    return fact
}

func modifyNFields(fact map[string]interface{}, n int) map[string]interface{} {
    modified := CopyFactFast(fact)
    count := 0
    for k := range modified {
        if count >= n {
            break
        }
        modified[k] = "modified"
        count++
    }
    return modified
}
```

---

## ✅ Validation

Après optimisations, exécuter :

```bash
# 1. Profiling
./scripts/profile_delta.sh

# 2. Benchmarks comparatifs
./scripts/benchmark_delta.sh

# 3. Vérifier gains
# Comparer avant/après optimisations

# 4. Tests de régression
go test ./rete/delta/... -v
go test ./tests/integration/... -v

# 5. Validation complète
make test
make benchmark
```

**Critères de succès** :
- [ ] Update 1 champ : **> 10x plus rapide** que classique
- [ ] Allocations réduites : **> 50%** moins d'allocations
- [ ] Overhead mémoire index : **< 5%**
- [ ] Latency p99 : **< 1ms**
- [ ] Aucune régression fonctionnelle

---

## 📊 Rapport de Performance

### Fichier : `REPORTS/delta_performance_report.md`

**Contenu** :

```markdown
# Rapport de Performance - Propagation Delta

## Résumé Exécutif

- **Objectif atteint** : ✅ / ❌
- **Gain moyen** : XXx plus rapide
- **Allocations réduites** : XX%
- **Overhead mémoire** : XX%

## Benchmarks Comparatifs

### Update 1 champ sur 10

| Métrique | Classique | Delta | Gain |
|----------|-----------|-------|------|
| Latency | XXX ns | XX ns | XXx |
| Allocations | XX | XX | -XX% |
| Bytes allocated | XXX B | XX B | -XX% |

### Update 5 champs sur 100

| Métrique | Classique | Delta | Gain |
|----------|-----------|-------|------|
| Latency | XXX ns | XX ns | XXx |
| Allocations | XX | XX | -XX% |

## Profiling CPU

Top 5 hotspots :
1. Function A - XX%
2. Function B - XX%
3. ...

## Profiling Mémoire

Top 5 allocations :
1. Type A - XX MB
2. Type B - XX MB
3. ...

## Optimisations Implémentées

1. **Object pooling** : -XX% allocations
2. **Cache LRU optimisé** : +XX% hit rate
3. **Inline hints** : -XX% latency
4. **Batch processing** : -XX% overhead

## Conclusions

- Delta propagation apporte un gain de XXx en moyenne
- Particulièrement efficace pour updates partiels (< 30% champs)
- Overhead acceptable (< 5% mémoire)

## Recommandations

1. Activer delta par défaut pour production
2. Configurer seuil à XX% pour fallback
3. Monitoring : surveiller hit rate cache
```

---

## 🚀 Commit

Une fois validé :

```bash
git add rete/delta/ scripts/ REPORTS/
git commit -m "perf(rete): [Prompt 09] Optimisations propagation delta

- Object pooling (FactDelta, NodeReference) : -50% allocations
- Cache LRU optimisé avec éviction : +30% hit rate
- Optimisations comparaison valeurs : -40% latency
- Batch processing nœuds : -20% overhead
- Scripts profiling et benchmarking
- Rapport performance détaillé
- Gain moyen : 12x plus rapide que classique
- Aucune régression fonctionnelle"
```

---

## 🚦 Prochaine Étape

Passer au **Prompt 10 - Documentation**

---

**Durée estimée** : 2-3 heures  
**Difficulté** : Élevée (expertise performance)  
**Prérequis** : Prompts 01-08 validés  
**Objectif** : > 10x speedup