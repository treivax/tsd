# 🏁 Guide du Race Detector pour TSD

## 📋 Résumé Exécutif

Le flag `-race` de Go détecte les **race conditions** (accès concurrents non synchronisés aux mêmes variables). C'est **OBLIGATOIRE** pour tous les prompts qui génèrent ou exécutent des tests.

## 🚫 Règle Absolue

```
TOUJOURS exécuter `go test -race` lors de la validation de tests.
NE JAMAIS skip cette étape, même si plus lente (~10x).
```

## ⚠️ Pourquoi C'est Critique

### Les Race Conditions Causent :
- ❌ **Bugs intermittents** : Apparaissent aléatoirement, impossibles à reproduire
- ❌ **Corruption de données** : Lectures/écritures simultanées → données invalides
- ❌ **Crashes production** : Panics inexpliqués sous charge
- ❌ **Tests flaky** : Passent parfois, échouent parfois
- ❌ **Métriques incorrectes** : Compteurs corrompus

### Race Conditions Sont :
- 🎲 **Timing-dependent** : Dépendent du scheduling des goroutines
- 🙈 **Invisibles sans `-race`** : Tests normaux ne les détectent PAS
- 🐛 **Difficiles à debug** : Symptômes apparaissent loin de la cause
- 💣 **Silencieuses** : Pas d'erreur visible, juste des résultats incorrects

## ✅ Quand Utiliser `-race`

| Situation | Commande | Obligatoire ? |
|-----------|----------|---------------|
| Dev local rapide | `go test ./pkg` | ❌ Non |
| Avant commit | `go test ./...` | ⚠️ Recommandé |
| **Avant pull request** | `make test-race` | ✅ **OUI** |
| **CI/CD** | `go test -race ./...` | ✅ **OUI** |
| **Deep-clean** | `go test -race ./...` | ✅ **OUI** |
| **Ajout/modification tests** | `go test -race ./...` | ✅ **OUI** |
| **Debug test flaky** | `go test -race -count=100` | ✅ **OUI** |
| **Avant release** | `make test-race` | ✅ **OUI** |

## 🔧 Commandes

### Commande Standard
```bash
# Tests avec race detector
go test -race ./...
```

### Via Makefile (Recommandé)
```bash
# Target dédié
make test-race
```

### Test Spécifique
```bash
# Un seul test
go test -race -run TestNomDuTest ./rete

# Un seul package
go test -race ./rete
```

### Avec Tags
```bash
# Avec tags e2e et integration
go test -race -tags=e2e,integration ./...
```

### Tests Répétés (Détecter Flaky)
```bash
# Répéter 10 fois
go test -race -count=10 ./...

# Répéter jusqu'à échec
go test -race -count=100 ./...
```

## 📊 Interpréter les Résultats

### ✅ Succès (Aucune Race)
```
ok      github.com/treivax/tsd/rete     7.402s
ok      github.com/treivax/tsd/constraint 2.855s
```
**Action** : Continuer ✅

### ❌ Échec (Race Détectée)
```
==================
WARNING: DATA RACE
Read at 0x00c000012345 by goroutine 21:
  github.com/treivax/tsd/rete.NewPipeline()
      /path/to/file.go:28 +0x184

Previous write at 0x00c000012345 by goroutine 9:
  github.com/treivax/tsd/test/util.captureOutput()
      /path/to/test.go:174 +0x64
==================
--- FAIL: TestFeature (0.24s)
    testing.go:1490: race detected during execution of test
FAIL
```

**Action** : 
1. 🛑 **STOP** - Ne pas continuer
2. 📝 Analyser la race (fichiers, lignes, goroutines)
3. 🔧 Fixer le problème de synchronisation
4. ✅ Relancer `go test -race` jusqu'à succès

## 🐛 Types de Race Conditions Communes

### 1. Compteur Non-Protégé
```go
// ❌ MAUVAIS - Race condition
var counter int
go func() { counter++ }()
go func() { counter++ }()

// ✅ BON - Avec mutex
var mu sync.Mutex
var counter int
go func() { mu.Lock(); counter++; mu.Unlock() }()
go func() { mu.Lock(); counter++; mu.Unlock() }()

// ✅ BON - Avec atomic
var counter int64
go func() { atomic.AddInt64(&counter, 1) }()
go func() { atomic.AddInt64(&counter, 1) }()
```

### 2. Map Non-Protégée
```go
// ❌ MAUVAIS - Race condition
m := make(map[string]int)
go func() { m["key"] = 1 }()
go func() { _ = m["key"] }()

// ✅ BON - Avec mutex
var mu sync.RWMutex
m := make(map[string]int)
go func() { mu.Lock(); m["key"] = 1; mu.Unlock() }()
go func() { mu.RLock(); _ = m["key"]; mu.RUnlock() }()

// ✅ BON - Avec sync.Map
var m sync.Map
go func() { m.Store("key", 1) }()
go func() { m.Load("key") }()
```

### 3. Variable Globale
```go
// ❌ MAUVAIS - Race condition
var config Config
go func() { config = loadConfig() }()
go func() { use(config) }()

// ✅ BON - Avec sync.Once
var config Config
var once sync.Once
go func() { once.Do(func() { config = loadConfig() }) }()
go func() { once.Do(func() { config = loadConfig() }); use(config) }()
```

### 4. Slice Non-Protégée
```go
// ❌ MAUVAIS - Race condition
var items []int
go func() { items = append(items, 1) }()
go func() { items = append(items, 2) }()

// ✅ BON - Avec channel
ch := make(chan int, 10)
go func() { ch <- 1 }()
go func() { ch <- 2 }()
go func() {
    items := []int{}
    for i := range ch {
        items = append(items, i)
    }
}()
```

## 🎯 Checklist de Validation

### Lors de l'Ajout de Tests
- [ ] Tests écrits
- [ ] Tests passent : `go test ./...`
- [ ] 🏁 **Race detector : `go test -race ./...`** (OBLIGATOIRE)
- [ ] Aucune race condition détectée
- [ ] Si race détectée → fixée

### Lors du Debug de Tests
- [ ] Problème identifié
- [ ] Correction implémentée
- [ ] Test corrigé passe : `go test -run TestName`
- [ ] 🏁 **Race detector : `go test -race -run TestName`** (OBLIGATOIRE)
- [ ] Tous tests : `make test`
- [ ] 🏁 **Race detector global : `make test-race`** (OBLIGATOIRE)
- [ ] Aucune race condition

### Lors du Deep-Clean
- [ ] Code nettoyé
- [ ] Tests passent : `go test ./...`
- [ ] 🏁 **Race detector : `go test -race ./...`** (OBLIGATOIRE)
- [ ] go vet : `go vet ./...`
- [ ] staticcheck : `staticcheck ./...`
- [ ] Build : `make build`

## ⏱️ Performance

### Coût du Race Detector
- **Temps d'exécution** : ~10x plus lent
- **Mémoire** : 5-10x plus d'utilisation
- **Build size** : Plus gros (instrumentation)

### Pourquoi C'est Acceptable
```
30 secondes de test avec -race
VS
Des heures/jours de debug en production

→ Toujours en valoir la peine !
```

## 🚀 Workflow Recommandé

### Development Rapide
```bash
# Tests rapides pendant le dev
go test ./pkg

# Tests complets avant commit
go test ./...
```

### Avant Pull Request
```bash
# Validation complète (OBLIGATOIRE)
make test          # Tests normaux
make test-race     # 🏁 Race detector (OBLIGATOIRE)
make build         # Build
```

### CI/CD Pipeline
```yaml
# .github/workflows/test.yml
- name: Test with race detector
  run: go test -race -tags=e2e,integration ./...
```

## 📚 Références

### Documentation Go
- [Go Race Detector](https://go.dev/doc/articles/race_detector)
- [Go Blog: Data Races](https://go.dev/blog/race-detector)
- [Effective Go: Concurrency](https://go.dev/doc/effective_go#concurrency)

### Projet TSD
- `docs/INSTALLATION.md` - Instructions d'installation
- `rete/docs/TESTING.md` - Tests RETE
- `tests/README.md` - Infrastructure de test
- `Makefile` - Target `test-race`

### Prompts Concernés
- `.github/prompts/add-test.md` - Ajout de tests
- `.github/prompts/debug-test.md` - Debug de tests
- `.github/prompts/run-tests.md` - Exécution de tests
- `.github/prompts/deep-clean.md` - Validation complète

## 💡 Cas d'Usage TSD

### Pourquoi TSD Est à Risque

TSD est un **moteur RETE** avec :
- ✅ Concurrence (goroutines multiples)
- ✅ État partagé (network nodes, tokens)
- ✅ Caches (beta join cache, LRU cache)
- ✅ Métriques (compteurs partagés)

→ **Risque ÉLEVÉ de race conditions !**

### Zones Critiques
1. **ReteNetwork** : Accès concurrent aux nodes
2. **BetaJoinCache** : Cache partagé entre goroutines
3. **Metrics** : Compteurs incrémentés en parallèle
4. **MemoryStorage** : Stockage des tokens
5. **Logger** : Écriture logs concurrent

## ⚠️ Ne JAMAIS

### ❌ Skip le Race Detector
```bash
# ❌ MAUVAIS - Validation incomplète
go test ./...
# "Ça passe, c'est bon !"

# ✅ BON - Validation complète
go test ./...
go test -race ./...  # OBLIGATOIRE
```

### ❌ Ignorer les Warnings
```bash
# ❌ MAUVAIS - Ignorer la race
$ go test -race ./...
WARNING: DATA RACE
# "C'est juste dans les tests, pas grave"

# ✅ BON - Fixer la race
$ go test -race ./...
WARNING: DATA RACE
# → Analyser et FIXER avant de continuer
```

### ❌ Assumer que Tests Normaux Suffisent
```bash
# ❌ MAUVAIS - Assumption dangereuse
$ go test ./...
ok  # "Pas de bug !"

# ✅ BON - Toujours vérifier races
$ go test ./...
ok
$ go test -race ./...  # Peut révéler des races !
WARNING: DATA RACE
```

## ✅ Toujours

### ✅ Exécuter Avant Validation
```bash
# Checklist complète
go test ./...           # Tests normaux
go test -race ./...     # 🏁 Race detector (OBLIGATOIRE)
go test -cover ./...    # Couverture
go vet ./...           # Analyse statique
staticcheck ./...      # Linter avancé
```

### ✅ Fixer les Races Immédiatement
```
Race détectée → STOP → Analyser → Fixer → Re-tester
```

### ✅ Documenter les Fixes
```go
// Fixed race condition: access to counter was not synchronized
// Now using sync.Mutex to protect concurrent access
var mu sync.Mutex
var counter int
```

## 🎯 Résumé

### Commande à Retenir
```bash
go test -race ./...
```

### Règle à Retenir
```
TOUJOURS exécuter -race lors de la validation de tests.
JAMAIS skip cette étape.
```

### Raison à Retenir
```
Race conditions = bugs invisibles qui causent
crashes production, corruption données, tests flaky.

Seul moyen de les détecter : go test -race
```

---

**Date de création** : 2025-12-08  
**Version** : 1.0  
**Statut** : Obligatoire pour tous les prompts de test

---

## 📞 Support

Si vous détectez une race condition et avez besoin d'aide :

1. **Copier le rapport complet** de `go test -race`
2. **Noter les fichiers/lignes** concernés
3. **Identifier les goroutines** impliquées
4. **Consulter ce guide** pour patterns de fix
5. **Demander review** si solution incertaine

**Ne jamais ignorer une race condition détectée.**