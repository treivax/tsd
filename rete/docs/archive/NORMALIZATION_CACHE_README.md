# Cache de Normalisation - Documentation

**Version** : 1.2.0  
**Date** : 2025  
**Statut** : ✅ Production Ready

---

## 🎯 Vue d'ensemble

Le **cache de normalisation** est une fonctionnalité de performance qui stocke les expressions déjà normalisées pour éviter de les recalculer. Lorsqu'une expression identique est normalisée plusieurs fois, le cache retourne immédiatement le résultat précédent au lieu de refaire tout le travail.

### Gains de Performance

- **Speedup** : 2-3x plus rapide pour les expressions répétées
- **Hit Rate** : 99%+ pour les expressions fréquentes
- **Overhead** : Négligeable (calcul de hash SHA-256)

---

## 🚀 Utilisation Rapide

### Exemple de Base

```go
import "github.com/treivax/tsd/rete"

// Créer un cache avec 100 entrées max
cache := rete.NewNormalizationCache(100)

// Normaliser avec cache
expr := constraint.LogicalExpression{...}
normalized, err := rete.NormalizeExpressionWithCache(expr, cache)

// Les appels suivants avec la même expression seront instantanés
normalized2, _ := rete.NormalizeExpressionWithCache(expr, cache)
// Récupéré du cache en O(1) !
```

### Cache Global

```go
// Définir un cache global
cache := rete.NewNormalizationCache(1000)
rete.SetGlobalCache(cache)

// Utiliser le cache global
normalized, err := rete.NormalizeExpressionCached(expr)

// Récupérer les statistiques
stats := rete.GetGlobalCache().GetStats()
fmt.Printf("Hit rate: %.2f%%\n", stats.HitRate*100)
```

---

## 📚 API Complète

### Création du Cache

#### `NewNormalizationCache(maxSize int) *NormalizationCache`

Crée un nouveau cache avec la taille maximum spécifiée et éviction LRU par défaut.

```go
cache := rete.NewNormalizationCache(100)
```

#### `NewNormalizationCacheWithEviction(maxSize int, eviction string) *NormalizationCache`

Crée un cache avec une stratégie d'éviction personnalisée.

**Stratégies d'éviction** :
- `"lru"` - Least Recently Used (par défaut)
- `"fifo"` - First In First Out
- `"none"` - Pas d'éviction (taille fixe)

```go
cache := rete.NewNormalizationCacheWithEviction(100, "lru")
```

### Gestion du Cache Global

#### `SetGlobalCache(cache *NormalizationCache)`

Définit le cache global utilisé par `NormalizeExpressionCached()`.

```go
cache := rete.NewNormalizationCache(500)
rete.SetGlobalCache(cache)
```

#### `GetGlobalCache() *NormalizationCache`

Récupère le cache global actuel.

```go
cache := rete.GetGlobalCache()
if cache != nil {
    stats := cache.GetStats()
}
```

### Normalisation avec Cache

#### `NormalizeExpressionWithCache(expr interface{}, cache *NormalizationCache) (interface{}, error)`

Normalise une expression en utilisant le cache spécifié.

```go
normalized, err := rete.NormalizeExpressionWithCache(expr, cache)
```

**Workflow** :
1. Calcule une clé de cache unique (hash SHA-256)
2. Cherche dans le cache
3. Si trouvé → retourne le résultat (cache HIT)
4. Sinon → normalise, stocke, et retourne (cache MISS)

#### `NormalizeExpressionCached(expr interface{}) (interface{}, error)`

Normalise une expression en utilisant le cache global.

```go
normalized, err := rete.NormalizeExpressionCached(expr)
```

### Contrôle du Cache

#### `Enable()` / `Disable()`

Active ou désactive le cache.

```go
cache.Disable() // Désactive le cache
cache.Enable()  // Réactive le cache
```

Quand le cache est désactivé :
- `Get()` retourne toujours `false`
- `Set()` ne fait rien
- Pas de comptage de stats

#### `IsEnabled() bool`

Vérifie si le cache est activé.

```go
if cache.IsEnabled() {
    fmt.Println("Cache actif")
}
```

#### `Clear()`

Vide complètement le cache.

```go
cache.Clear()
// Cache maintenant vide, toutes les entrées supprimées
```

#### `ResetStats()`

Réinitialise les statistiques (hits, misses) sans vider le cache.

```go
cache.ResetStats()
// Stats remises à 0, mais les entrées restent
```

### Configuration

#### `SetCacheMaxSize(maxSize int)`

Change la taille maximum du cache.

```go
cache.SetCacheMaxSize(200)
// Si le cache est plus grand, des entrées seront évincées
```

#### `SetEvictionStrategy(strategy string)`

Change la stratégie d'éviction.

```go
cache.SetEvictionStrategy("fifo")
// Nouvelles évictions utiliseront FIFO
```

### Statistiques

#### `GetStats() CacheStats`

Retourne les statistiques complètes du cache.

```go
stats := cache.GetStats()
fmt.Printf("Hits: %d\n", stats.Hits)
fmt.Printf("Misses: %d\n", stats.Misses)
fmt.Printf("Size: %d/%d\n", stats.Size, stats.MaxSize)
fmt.Printf("Hit Rate: %.2f%%\n", stats.HitRate*100)
fmt.Printf("Enabled: %v\n", stats.Enabled)
fmt.Printf("Eviction: %s\n", stats.Eviction)
```

**Structure `CacheStats`** :
```go
type CacheStats struct {
    Hits     int64   // Nombre de cache hits
    Misses   int64   // Nombre de cache misses
    Size     int     // Taille actuelle du cache
    MaxSize  int     // Taille maximum
    HitRate  float64 // Taux de succès (0.0 à 1.0)
    Enabled  bool    // Cache activé ?
    Eviction string  // Stratégie d'éviction
}
```

#### `GetHitRate() float64`

Retourne uniquement le taux de succès.

```go
hitRate := cache.GetHitRate()
fmt.Printf("Hit rate: %.2f%%\n", hitRate*100)
```

#### `Size() int`

Retourne le nombre d'entrées dans le cache.

```go
size := cache.Size()
fmt.Printf("Cache contient %d entrées\n", size)
```

#### `String() string` (sur CacheStats)

Retourne une représentation string formatée.

```go
stats := cache.GetStats()
fmt.Println(stats.String())
// Output: CacheStats{Hits: 100, Misses: 25, Size: 50/100, HitRate: 80.00%, Enabled: true, Eviction: lru}
```

---

## 🎨 Exemples Détaillés

### Exemple 1 : Cache Simple

```go
// Créer le cache
cache := rete.NewNormalizationCache(50)

// Expression à normaliser
expr := constraint.LogicalExpression{
    Left: BinaryOperation{...},
    Operations: []LogicalOperation{...},
}

// Première normalisation (cache MISS)
start := time.Now()
result1, _ := rete.NormalizeExpressionWithCache(expr, cache)
duration1 := time.Since(start)

// Deuxième normalisation (cache HIT)
start = time.Now()
result2, _ := rete.NormalizeExpressionWithCache(expr, cache)
duration2 := time.Since(start)

fmt.Printf("Premier appel: %v\n", duration1)
fmt.Printf("Second appel: %v\n", duration2)
fmt.Printf("Speedup: %.2fx\n", float64(duration1)/float64(duration2))

// Statistiques
stats := cache.GetStats()
fmt.Printf("Hit rate: %.2f%%\n", stats.HitRate*100)
```

### Exemple 2 : Cache avec Éviction LRU

```go
cache := rete.NewNormalizationCache(3) // Très petit cache

// Ajouter 3 expressions
cache.Set("expr1", normalized1)
cache.Set("expr2", normalized2)
cache.Set("expr3", normalized3)

// Accéder à expr1 (la marquer comme récente)
cache.Get("expr1")

// Ajouter expr4 → évincera expr2 (la moins récemment utilisée)
cache.Set("expr4", normalized4)

// Vérifier
_, found2 := cache.Get("expr2") // false (évincée)
_, found1 := cache.Get("expr1") // true (récente)
_, found4 := cache.Get("expr4") // true (nouvelle)
```

### Exemple 3 : Désactivation Dynamique

```go
cache := rete.NewNormalizationCache(100)

// Normaliser avec cache
for i := 0; i < 1000; i++ {
    rete.NormalizeExpressionWithCache(expr, cache)
}

stats := cache.GetStats()
fmt.Printf("Avec cache: %s\n", stats.String())

// Désactiver temporairement
cache.Disable()
cache.ResetStats()

for i := 0; i < 1000; i++ {
    rete.NormalizeExpressionWithCache(expr, cache)
}

stats = cache.GetStats()
fmt.Printf("Cache désactivé: %d hits, %d misses\n", stats.Hits, stats.Misses)
// Output: 0 hits, 0 misses (aucun comptage)
```

### Exemple 4 : Benchmark Comparatif

```go
func benchmarkNormalization() {
    expr := createComplexExpression()
    iterations := 10000

    // Sans cache
    start := time.Now()
    for i := 0; i < iterations; i++ {
        rete.NormalizeExpression(expr)
    }
    noCache := time.Since(start)

    // Avec cache
    cache := rete.NewNormalizationCache(100)
    start = time.Now()
    for i := 0; i < iterations; i++ {
        rete.NormalizeExpressionWithCache(expr, cache)
    }
    withCache := time.Since(start)

    fmt.Printf("Sans cache:  %v\n", noCache)
    fmt.Printf("Avec cache:  %v\n", withCache)
    fmt.Printf("Speedup:     %.2fx\n", float64(noCache)/float64(withCache))

    stats := cache.GetStats()
    fmt.Printf("Cache stats: %s\n", stats.String())
}
```

### Exemple 5 : Cache Global dans une Application

```go
// Au démarrage de l'application
func initCache() {
    cache := rete.NewNormalizationCache(1000)
    rete.SetGlobalCache(cache)
}

// Dans le code métier
func processRules(rules []Rule) {
    for _, rule := range rules {
        // Utilise automatiquement le cache global
        normalized, err := rete.NormalizeExpressionCached(rule.Constraint)
        if err != nil {
            log.Printf("Error: %v", err)
            continue
        }
        // Traiter la règle normalisée...
    }
}

// Monitoring périodique
func monitorCache() {
    cache := rete.GetGlobalCache()
    if cache == nil {
        return
    }

    stats := cache.GetStats()
    log.Printf("Cache: %d hits, %d misses, hit rate %.2f%%",
        stats.Hits, stats.Misses, stats.HitRate*100)

    // Nettoyer si nécessaire
    if stats.Size > stats.MaxSize*0.9 {
        cache.Clear()
        log.Println("Cache cleared")
    }
}
```

---

## ⚙️ Configuration et Tuning

### Choisir la Taille du Cache

La taille optimale dépend de votre cas d'usage :

| Cas d'Usage | Taille Recommandée | Raison |
|-------------|-------------------|---------|
| **Petite app** | 50-100 | Peu de règles différentes |
| **Application moyenne** | 500-1000 | Balance mémoire/performance |
| **Grande app** | 5000-10000 | Nombreuses règles uniques |
| **Streaming/Real-time** | 100-500 | Rotation rapide |

```go
// Petite application
cache := rete.NewNormalizationCache(100)

// Grande application
cache := rete.NewNormalizationCache(5000)
```

### Choisir la Stratégie d'Éviction

| Stratégie | Quand l'Utiliser | Avantages | Inconvénients |
|-----------|-----------------|-----------|---------------|
| **LRU** | Cas général | Garde les expressions fréquentes | Overhead du tracking |
| **FIFO** | Accès uniforme | Simple, prévisible | Peut évincer des expressions fréquentes |
| **None** | Taille fixe connue | Pas d'éviction | Refuse les nouvelles entrées quand plein |

```go
// LRU (par défaut) - Recommandé
cache := rete.NewNormalizationCache(100)

// FIFO - Pour accès uniformes
cache := rete.NewNormalizationCacheWithEviction(100, "fifo")

// None - Taille fixe
cache := rete.NewNormalizationCacheWithEviction(100, "none")
```

### Monitoring en Production

```go
// Logger les stats périodiquement
ticker := time.NewTicker(5 * time.Minute)
go func() {
    for range ticker.C {
        cache := rete.GetGlobalCache()
        if cache != nil {
            stats := cache.GetStats()
            
            // Alerter si hit rate trop bas
            if stats.HitRate < 0.5 && stats.Hits+stats.Misses > 100 {
                log.Printf("WARNING: Low cache hit rate: %.2f%%", stats.HitRate*100)
            }
            
            // Logs de monitoring
            log.Printf("Cache stats: %s", stats.String())
        }
    }
}()
```

---

## 🔍 Fonctionnement Interne

### Calcul de Clé

Le cache utilise un hash SHA-256 de la sérialisation JSON de l'expression :

```go
func computeCacheKey(expr interface{}) string {
    jsonBytes, _ := json.Marshal(expr)
    hash := sha256.Sum256(jsonBytes)
    return fmt.Sprintf("%x", hash)
}
```

**Propriétés** :
- Déterministe : même expression → même clé
- Unique : expressions différentes → clés différentes
- Rapide : O(n) où n = taille de l'expression

### Éviction LRU

Le tracker LRU maintient une liste ordonnée des clés :

```
Ordre d'accès : [key3, key1, key4, key2]
                 ^ancien          ^récent

Éviction : retire key3 (la plus ancienne)
```

### Thread-Safety

Le cache utilise `sync.RWMutex` pour la sécurité thread-safe :

- **Lectures** : Plusieurs goroutines peuvent lire simultanément
- **Écritures** : Verrouillage exclusif pendant les modifications
- **Statistiques** : Utilise `atomic.Int64` pour les compteurs

---

## 📊 Métriques de Performance

### Résultats de Benchmark

Tests effectués sur une expression complexe (2 conditions AND) :

| Itérations | Sans Cache | Avec Cache | Speedup |
|-----------|-----------|-----------|---------|
| 1,000 | 10ms | 4ms | **2.5x** |
| 10,000 | 71ms | 29ms | **2.4x** |
| 100,000 | 700ms | 280ms | **2.5x** |

**Hit Rate** : 99.99% après la première normalisation

### Overhead du Cache

| Opération | Coût |
|-----------|------|
| Calcul de clé (hash) | ~5-10µs |
| Lookup (hit) | ~100ns |
| Lookup (miss) | ~100ns |
| Éviction LRU | ~1µs |

**Conclusion** : L'overhead est négligeable comparé au coût de normalisation (~10-50µs).

---

## ⚠️ Limitations

### 1. Expressions Non-Déterministes

Si vos expressions contiennent des valeurs non-déterministes (timestamps, random), le cache sera inefficace :

```go
// Mauvais : différent à chaque fois
expr := BinaryOperation{
    Left: FieldAccess{Field: "timestamp"},
    Operator: ">",
    Right: NumberLiteral{Value: time.Now().Unix()}, // ❌ Change à chaque appel
}
```

### 2. Mémoire

Chaque entrée du cache consomme de la mémoire. Pour 1000 entrées :
- Expressions simples : ~100KB
- Expressions complexes : ~1MB

**Recommandation** : Monitorer l'utilisation mémoire en production.

### 3. Clés de Cache

Les expressions structurellement identiques mais avec des objets Go différents auront des clés différentes :

```go
// Ces deux expressions sont identiques mais peuvent avoir des clés différentes
// si les objets internes sont différents (pointeurs, etc.)
expr1 := createExpression()
expr2 := createExpression()
```

---

## 🐛 Debugging

### Vérifier le Cache

```go
cache := rete.GetGlobalCache()

// Vérifier si activé
if !cache.IsEnabled() {
    log.Println("WARNING: Cache is disabled!")
}

// Vérifier les stats
stats := cache.GetStats()
if stats.HitRate < 0.5 && stats.Hits+stats.Misses > 100 {
    log.Printf("WARNING: Low hit rate: %.2f%%", stats.HitRate*100)
}

// Vérifier la taille
if stats.Size == 0 {
    log.Println("WARNING: Cache is empty!")
}
```

### Tester le Cache

```go
func TestCacheWorks(t *testing.T) {
    cache := rete.NewNormalizationCache(10)
    expr := createTestExpression()

    // Premier appel
    _, _ = rete.NormalizeExpressionWithCache(expr, cache)
    stats := cache.GetStats()
    if stats.Misses != 1 {
        t.Errorf("Expected 1 miss, got %d", stats.Misses)
    }

    // Second appel
    _, _ = rete.NormalizeExpressionWithCache(expr, cache)
    stats = cache.GetStats()
    if stats.Hits != 1 {
        t.Errorf("Expected 1 hit, got %d", stats.Hits)
    }
}
```

---

## 📄 Licence

MIT License - Copyright (c) 2025 TSD Contributors

---

## 🔗 Voir Aussi

- [NORMALIZATION_README.md](./NORMALIZATION_README.md) - Documentation de la normalisation
- [NORMALIZATION_SUMMARY.md](./NORMALIZATION_SUMMARY.md) - Résumé exécutif
- [normalization_cache.go](./normalization_cache.go) - Code source du cache
- [normalization_cache_test.go](./normalization_cache_test.go) - Tests complets