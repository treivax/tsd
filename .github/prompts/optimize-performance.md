# ⚡ Optimiser les Performances (Optimize Performance)

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Tu veux optimiser les performances d'une partie du code qui est trop lente, consomme trop de mémoire, ou présente des goulots d'étranglement.

## Objectif

Optimiser les performances de manière méthodique en :
- Mesurant les performances actuelles (baseline)
- Identifiant les goulots d'étranglement réels
- Optimisant de façon ciblée et mesurée
- Validant les gains de performance
- Préservant la sémantique et la qualité du code

## 📄 RÈGLES DE LICENCE ET COPYRIGHT - OBLIGATOIRE

### 🔒 Vérification de Compatibilité de Licence

**SI l'optimisation nécessite du code externe ou une nouvelle bibliothèque** :

1. **Vérifier la licence** :
   - ✅ Licences permissives acceptées : MIT, BSD, Apache-2.0, ISC
   - ⚠️ Licences à éviter : GPL, AGPL, LGPL (copyleft)
   - ❌ Code sans licence = NE PAS UTILISER
   - ❌ Code propriétaire = NE PAS UTILISER

2. **Documenter l'origine** :
   - Si code inspiré/adapté : ajouter commentaire avec source
   - Si bibliothèque tierce : mettre à jour `go.mod` et `THIRD_PARTY_LICENSES.md`
   - Si algorithme connu : citer la référence académique

### 📝 En-tête de Copyright OBLIGATOIRE

**SI création de nouveaux fichiers durant l'optimisation** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package [nom_du_package]
```

**VÉRIFICATION** :
- ✅ Tous les nouveaux fichiers .go ont l'en-tête de copyright
- ✅ Les fichiers existants conservent leur en-tête
- ✅ Aucun code externe non vérifié n'est introduit

### ⚠️ INTERDICTIONS STRICTES

- ❌ **Ne JAMAIS copier du code** sans vérifier la licence
- ❌ **Ne JAMAIS utiliser de code GPL/AGPL** (incompatible avec MIT)
- ❌ **Ne JAMAIS omettre les en-têtes de copyright** dans les nouveaux fichiers
- ✅ **TOUJOURS écrire du code original** lors d'optimisations

## ⚠️ RÈGLES STRICTES - OPTIMISATION

### 🚫 INTERDICTIONS ABSOLUES

1. **CODE GOLANG** :
   - ❌ AUCUN HARDCODING introduit
   - ❌ AUCUNE optimisation prématurée
   - ❌ AUCUNE optimisation sans mesure
   - ❌ AUCUN sacrifice de lisibilité sans gain prouvé
   - ✅ Code générique avec paramètres/interfaces
   - ✅ Constantes nommées pour toutes les valeurs
   - ✅ Respect strict Effective Go

2. **TESTS RETE** :
   - ❌ AUCUNE simulation de résultats
   - ❌ AUCUN test qui ne valide pas la sémantique
   - ✅ Extraction depuis réseau RETE réel uniquement
   - ✅ Validation sémantique après optimisation
   - ✅ Benchmarks avec données réelles

3. **MÉTHODOLOGIE** :
   - ❌ Pas d'optimisation sans profiling
   - ❌ Pas de commit sans benchmarks
   - ✅ Mesures avant/après obligatoires
   - ✅ Validation sémantique obligatoire
   - ✅ Documentation des gains

## Instructions

### PHASE 1 : MESURE (Baseline)

#### 1.1 Identifier le Problème de Performance

**Collecte d'informations** :

```
Performance Issue ID : [Numéro si applicable]
Titre : [Description courte]
Type : CPU / Mémoire / I/O / Latence / Débit
Sévérité : Critique / Majeure / Mineure

Description :
[Description du problème de performance]

Comportement actuel :
- Temps d'exécution : X ms/s
- Consommation mémoire : Y MB
- Débit : Z ops/s

Comportement cible :
- Temps d'exécution : < X ms/s
- Consommation mémoire : < Y MB
- Débit : > Z ops/s

Contexte :
- Charge typique : ...
- Cas d'usage : ...
- Environnement : ...
```

#### 1.2 Établir la Baseline (Mesure Initiale)

**Benchmarks Go** :

```go
// rete/performance_test.go
func BenchmarkCurrentImplementation(b *testing.B) {
    // Setup
    network := buildLargeNetwork()
    facts := generateTestFacts(1000)
    
    b.ResetTimer()
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        for _, fact := range facts {
            network.SubmitFact(fact)
        }
    }
}

func BenchmarkSpecificFunction(b *testing.B) {
    bindings := setupBindings()
    condition := setupCondition()
    
    b.ResetTimer()
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        evaluateCondition(bindings, condition)
    }
}
```

**Exécuter benchmarks** :

```bash
# Baseline initiale
go test -bench=. -benchmem ./rete > baseline.txt
go test -bench=. -benchmem -benchtime=10s ./rete >> baseline.txt

# Résultats typiques :
# BenchmarkCurrentImplementation-8    1000    1234567 ns/op    123456 B/op    1234 allocs/op
```

#### 1.3 Profiling Détaillé

**CPU Profiling** :

```bash
# Générer profil CPU
go test -bench=BenchmarkCurrentImplementation -cpuprofile=cpu.prof ./rete

# Analyser avec pprof
go tool pprof cpu.prof
# (pprof) top10
# (pprof) list functionName
# (pprof) web  # Génère graphique si graphviz installé

# Ou interactif
go tool pprof -http=:8080 cpu.prof
```

**Memory Profiling** :

```bash
# Générer profil mémoire
go test -bench=BenchmarkCurrentImplementation -memprofile=mem.prof ./rete

# Analyser
go tool pprof mem.prof
# (pprof) top10
# (pprof) list functionName

# Vérifier allocations
go tool pprof -alloc_space mem.prof
go tool pprof -alloc_objects mem.prof
```

**Trace Analysis** :

```bash
# Générer trace
go test -bench=BenchmarkCurrentImplementation -trace=trace.out ./rete

# Visualiser
go tool trace trace.out
# Ouvre navigateur avec visualisation détaillée
```

#### 1.4 Identifier les Bottlenecks

**Analyse des résultats** :

```
CPU Hotspots (top consommateurs) :
1. functionA : 45% du temps CPU
2. functionB : 30% du temps CPU
3. functionC : 15% du temps CPU

Allocations Mémoire :
1. []byte allocations : 50 MB
2. map[string]*Fact : 30 MB
3. Token structs : 20 MB

Goulots d'étranglement identifiés :
- Boucles imbriquées dans evaluateCondition
- Allocations répétées dans createToken
- Copy de slices dans propagateToChildren
```

**Priorisation** :

1. **Impact élevé / Effort faible** → Priorité 1
2. **Impact élevé / Effort élevé** → Priorité 2
3. **Impact faible / Effort faible** → Priorité 3
4. **Impact faible / Effort élevé** → À éviter

### PHASE 2 : OPTIMISATION (Action)

#### 2.1 Stratégies d'Optimisation

**Optimisations Algorithmiques** :

```go
// ❌ AVANT - O(n²)
func findMatch(items []Item, target string) *Item {
    for _, item := range items {
        for _, prop := range item.Properties {
            if prop == target {
                return &item
            }
        }
    }
    return nil
}

// ✅ APRÈS - O(1) avec index
type ItemIndex struct {
    byProperty map[string]*Item
}

func (idx *ItemIndex) findMatch(target string) *Item {
    return idx.byProperty[target]
}
```

**Optimisations Mémoire** :

```go
// ❌ AVANT - Allocations répétées
func process(items []string) []Result {
    results := []Result{}  // Croissance dynamique
    for _, item := range items {
        results = append(results, processItem(item))
    }
    return results
}

// ✅ APRÈS - Pré-allocation
func process(items []string) []Result {
    results := make([]Result, 0, len(items))  // Capacité connue
    for _, item := range items {
        results = append(results, processItem(item))
    }
    return results
}
```

**Sync.Pool pour Réutilisation** :

```go
// Pool pour réutiliser objets
var tokenPool = sync.Pool{
    New: func() interface{} {
        return &Token{
            Bindings: make(map[string]*Fact, 8),
            Facts:    make([]*Fact, 0, 8),
        }
    },
}

// Obtenir token du pool
func getToken() *Token {
    return tokenPool.Get().(*Token)
}

// Retourner au pool
func putToken(t *Token) {
    // Reset
    for k := range t.Bindings {
        delete(t.Bindings, k)
    }
    t.Facts = t.Facts[:0]
    tokenPool.Put(t)
}
```

**Éviter Copies Inutiles** :

```go
// ❌ AVANT - Copy slice
func propagate(tokens []*Token) {
    for _, token := range tokens {
        // Copy complète du slice
        newTokens := make([]*Token, len(tokens))
        copy(newTokens, tokens)
        process(newTokens)
    }
}

// ✅ APRÈS - Utiliser directement
func propagate(tokens []*Token) {
    for i := range tokens {
        // Utiliser référence directe
        process(tokens[i])
    }
}
```

**Batch Processing** :

```go
// ❌ AVANT - Traitement un par un
func submitFacts(facts []*Fact) {
    for _, fact := range facts {
        network.SubmitFact(fact)
        network.Propagate()  // Propagation à chaque fois
    }
}

// ✅ APRÈS - Batch
func submitFacts(facts []*Fact) {
    for _, fact := range facts {
        network.SubmitFact(fact)
    }
    network.PropagateAll()  // Propagation groupée
}
```

#### 2.2 Implémenter les Optimisations

**Processus** :

1. **Une optimisation à la fois** :
   ```bash
   git checkout -b perf/optimize-function-name
   ```

2. **Mesurer avant** :
   ```bash
   go test -bench=BenchmarkFunction -benchmem ./rete > before.txt
   ```

3. **Optimiser** :
   ```go
   // ⚠️ VÉRIFIER : Aucun hardcoding introduit
   // ⚠️ VÉRIFIER : Code générique maintenu
   // ⚠️ VÉRIFIER : Sémantique préservée
   
   // Implémentation optimisée
   ```

4. **Mesurer après** :
   ```bash
   go test -bench=BenchmarkFunction -benchmem ./rete > after.txt
   ```

5. **Comparer** :
   ```bash
   # Utiliser benchstat pour comparaison statistique
   go install golang.org/x/perf/cmd/benchstat@latest
   benchstat before.txt after.txt
   ```

#### 2.3 Valider la Sémantique

**Tests de non-régression** :

```go
func TestOptimization_SemanticsPreserved(t *testing.T) {
    t.Log("🔍 VALIDATION SÉMANTIQUE POST-OPTIMISATION")
    t.Log("============================================")
    
    // Même setup que baseline
    network := buildNetwork()
    facts := generateTestFacts(100)
    
    // Soumettre faits
    for _, fact := range facts {
        network.SubmitFact(fact)
    }
    
    // ✅ OBLIGATOIRE : Extraction depuis réseau RETE réel
    actualTokens := 0
    for _, terminal := range network.TerminalNodes {
        actualTokens += len(terminal.Memory.GetTokens())
    }
    
    // ❌ INTERDIT : expectedTokens := 50 (hardcodé)
    
    // Vérifier que résultat sémantique est identique
    if actualTokens == 0 {
        t.Error("❌ Optimisation a changé la sémantique")
    }
    
    // Vérifier contenu des tokens
    for _, terminal := range network.TerminalNodes {
        for _, token := range terminal.Memory.GetTokens() {
            if len(token.Bindings) == 0 {
                t.Error("❌ Tokens invalides après optimisation")
            }
        }
    }
    
    t.Logf("✅ Sémantique préservée : %d tokens", actualTokens)
}
```

**Tests complets** :

```bash
# Tous les tests doivent passer
go test ./...
go test -race ./...
make test-integration
make rete-unified  # 58/58 ✅
```

### PHASE 3 : VALIDATION (Mesure des Gains)

#### 3.1 Benchmarks Comparatifs

**Analyse détaillée** :

```bash
# Benchmarks avec statistiques
go test -bench=. -benchmem -count=10 ./rete > optimized.txt

# Comparaison statistique
benchstat baseline.txt optimized.txt

# Résultat attendu :
# name                        old time/op    new time/op    delta
# CurrentImplementation-8       1.23ms ± 2%    0.45ms ± 1%  -63.41%  (p=0.000 n=10+10)
#
# name                        old alloc/op   new alloc/op   delta
# CurrentImplementation-8       123kB ± 0%      45kB ± 0%  -63.41%  (p=0.000 n=10+10)
#
# name                        old allocs/op  new allocs/op  delta
# CurrentImplementation-8       1.23k ± 0%     0.45k ± 0%  -63.41%  (p=0.000 n=10+10)
```

#### 3.2 Profiling Post-Optimisation

**Vérifier les gains** :

```bash
# Nouveau profil CPU
go test -bench=BenchmarkOptimized -cpuprofile=cpu_opt.prof ./rete
go tool pprof -top cpu_opt.prof

# Comparer avec baseline
go tool pprof -base cpu.prof cpu_opt.prof

# Nouveau profil mémoire
go test -bench=BenchmarkOptimized -memprofile=mem_opt.prof ./rete
go tool pprof -top mem_opt.prof
```

#### 3.3 Tests de Charge

**Tests avec données réelles** :

```go
func BenchmarkRealWorldScenario(b *testing.B) {
    // Scénario réaliste : 10k faits, 50 règles
    network := buildProductionNetwork()
    facts := loadRealWorldFacts(10000)
    
    b.ResetTimer()
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        for _, fact := range facts {
            network.SubmitFact(fact)
        }
    }
}
```

**Tests de scalabilité** :

```bash
# Différentes tailles
go test -bench=. -benchmem ./rete -args -size=100
go test -bench=. -benchmem ./rete -args -size=1000
go test -bench=. -benchmem ./rete -args -size=10000
go test -bench=. -benchmem ./rete -args -size=100000
```

## Critères de Succès

### ✅ Performance Améliorée

- [ ] Baseline mesurée et documentée
- [ ] Bottlenecks identifiés par profiling
- [ ] Optimisations implémentées **sans hardcoding**
- [ ] Gains mesurés avec benchstat
- [ ] Performance cible atteinte
- [ ] Scalabilité validée

### ✅ Sémantique Préservée

- [ ] **Tous les tests passent** (unit, race, integration)
- [ ] **Tests RETE avec extraction réseau réel**
- [ ] Runner universel : 58/58 ✅
- [ ] Résultats identiques à baseline
- [ ] Aucune régression fonctionnelle

### ✅ Qualité Maintenue

- [ ] **Aucun hardcoding** introduit
- [ ] **Code générique** maintenu
- [ ] go vet et golangci-lint sans erreur
- [ ] Code reste lisible et maintenable
- [ ] Documentation à jour

### ✅ Gains Documentés

- [ ] Benchmarks avant/après
- [ ] Profils CPU/mémoire
- [ ] Rapport d'optimisation
- [ ] CHANGELOG mis à jour

## Format de Réponse

```
=== OPTIMISATION DE PERFORMANCE ===

📋 IDENTIFICATION

Performance Issue : [Titre]
Type : CPU / Mémoire / Latence / Débit
Sévérité : [Critique/Majeure/Mineure]

Fonction/Module : [Localisation]

📊 BASELINE (AVANT)

Benchmarks :
  • Temps d'exécution : 1.23 ms/op
  • Allocations mémoire : 123 kB/op
  • Nombre d'allocations : 1234 allocs/op
  • Débit : 813 ops/s

Profiling CPU :
  • functionA : 45% CPU
  • functionB : 30% CPU
  • functionC : 15% CPU

Profiling Mémoire :
  • []byte allocations : 50 MB
  • map allocations : 30 MB
  • struct allocations : 20 MB

🎯 BOTTLENECKS IDENTIFIÉS

1. Boucles imbriquées O(n²) dans evaluateCondition
   Impact : 45% du temps CPU
   Priorité : Haute

2. Allocations répétées dans createToken
   Impact : 40% de la mémoire
   Priorité : Haute

3. Copies inutiles dans propagateToChildren
   Impact : 15% du temps CPU
   Priorité : Moyenne

⚡ OPTIMISATIONS APPLIQUÉES

1. Algorithme O(n²) → O(1) avec index
   Fichier : rete/node_join.go
   Lignes : 265-280
   ⚠️ **VÉRIFIÉ** : Aucun hardcoding
   ⚠️ **VÉRIFIÉ** : Code générique

2. Utilisation sync.Pool pour tokens
   Fichier : rete/token.go
   Lignes : 45-60
   ⚠️ **VÉRIFIÉ** : Pas de hardcoding

3. Élimination copies avec références directes
   Fichier : rete/propagation.go
   Lignes : 120-135

📊 RÉSULTATS (APRÈS)

Benchmarks :
  • Temps d'exécution : 0.45 ms/op  (-63.4%)
  • Allocations mémoire : 45 kB/op   (-63.4%)
  • Nombre d'allocations : 450 allocs/op (-63.4%)
  • Débit : 2222 ops/s  (+173%)

Comparaison statistique (benchstat) :
name                    old time/op    new time/op    delta
Implementation-8          1.23ms ± 2%    0.45ms ± 1%  -63.41%

name                    old alloc/op   new alloc/op   delta
Implementation-8           123kB ± 0%      45kB ± 0%  -63.41%

✅ VALIDATION SÉMANTIQUE

Tests :
  ✅ go test ./... : PASS
  ✅ go test -race ./... : PASS
  ✅ make test-integration : PASS
  ✅ make rete-unified : 58/58 ✅
  ⚠️ **VÉRIFIÉ** : Extraction réseau RETE réel

Régression :
  ✅ Résultats identiques à baseline
  ✅ Aucune régression fonctionnelle
  ✅ Sémantique 100% préservée

Qualité :
  ✅ go vet : 0 erreur
  ✅ golangci-lint : 0 erreur
  ✅ Code lisible et maintenable

📈 GAINS MESURÉS

Performance :
  • Temps CPU : -63.4% (1.23ms → 0.45ms)
  • Mémoire : -63.4% (123kB → 45kB)
  • Allocations : -63.4% (1234 → 450)
  • Débit : +173% (813 → 2222 ops/s)

Scalabilité :
  • 100 faits : 45ms (avant: 123ms)
  • 1000 faits : 450ms (avant: 1230ms)
  • 10000 faits : 4.5s (avant: 12.3s)

🎯 VERDICT : OPTIMISATION RÉUSSIE ✅

Objectif : < 1ms/op → Atteint (0.45ms)
Gains : 63.4% temps CPU, 63.4% mémoire
Sémantique : Préservée ✅
Qualité : Maintenue ✅
```

## Exemple d'Utilisation

```
La propagation RETE est trop lente avec plus de 1000 faits.
Profiling montre que evaluateCondition prend 45% du temps CPU.

Utilise le prompt "optimize-performance" pour :
1. Établir baseline avec benchmarks
2. Profiler CPU et mémoire
3. Identifier bottlenecks
4. Optimiser de façon ciblée
5. Valider gains et sémantique
```

## Checklist d'Optimisation

### Avant de Commencer
- [ ] Problème de performance identifié
- [ ] Baseline mesurée (benchmarks)
- [ ] Profiling effectué (CPU, mémoire)
- [ ] Bottlenecks identifiés
- [ ] Objectif de performance défini

### Pendant l'Optimisation
- [ ] Optimisation ciblée sur bottleneck
- [ ] Une optimisation à la fois
- [ ] Benchmarks avant/après chaque changement
- [ ] **Aucun hardcoding** introduit
- [ ] **Code générique** maintenu
- [ ] Sémantique préservée

### Après l'Optimisation
- [ ] **Gains mesurés** avec benchstat ✅
- [ ] **Tous les tests passent** ✅
- [ ] **Tests RETE extraction réseau réel** ✅
- [ ] **Aucune régression** ✅
- [ ] go vet et golangci-lint sans erreur ✅
- [ ] Code reste lisible ✅
- [ ] Documentation mise à jour ✅
- [ ] CHANGELOG mis à jour ✅

## Commandes Utiles

```bash
# Benchmarks
go test -bench=. -benchmem ./rete
go test -bench=. -benchmem -benchtime=10s ./rete
go test -bench=BenchmarkSpecific -count=10 -benchmem ./rete

# Profiling
go test -bench=. -cpuprofile=cpu.prof ./rete
go test -bench=. -memprofile=mem.prof ./rete
go test -bench=. -trace=trace.out ./rete

# Analyse
go tool pprof cpu.prof
go tool pprof -http=:8080 cpu.prof
go tool trace trace.out

# Comparaison
benchstat before.txt after.txt

# Validation
go test ./...
go test -race ./...
make rete-unified
```

## Bonnes Pratiques

1. **Mesurer d'abord** : Pas d'optimisation sans profiling
2. **Une chose à la fois** : Isoler chaque optimisation
3. **Comparer** : Toujours benchstat avant/après
4. **Valider sémantique** : Tests RETE avec extraction réelle
5. **Documenter gains** : Rapport avec chiffres précis
6. **Préserver lisibilité** : Pas d'optimisation obscure sans gain significatif
7. **Respecter règles** : Pas de hardcoding, code générique

## Patterns d'Optimisation Go

### Allocations
```go
// Pré-allocation de slices
s := make([]T, 0, capacity)

// Réutilisation avec sync.Pool
var pool = sync.Pool{New: func() interface{} { return new(T) }}

// Éviter conversions string <-> []byte
// Utiliser []byte directement
```

### CPU
```go
// Cache pour calculs répétés
var cache = make(map[Key]Value)

// Batch processing au lieu de un par un
// Traiter par lots de 100-1000

// Éviter réflexion (reflection) dans hot path
// Utiliser type assertions ou code gen
```

### Algorithmes
```go
// O(n²) → O(n log n) avec tri
// O(n) → O(1) avec map/index
// Éviter boucles imbriquées
```

## Outils Recommandés

- **benchstat** : Comparaison statistique de benchmarks
- **pprof** : Profiling CPU et mémoire
- **trace** : Analyse de traces d'exécution
- **golangci-lint** : Détection issues de performance
- **staticcheck** : Analyse statique

## Ressources

- [Go Performance Tips](https://github.com/dgryski/go-perfbook)
- [High Performance Go Workshop](https://dave.cheney.net/high-performance-go-workshop/gopherchina-2019.html)
- [Profiling Go Programs](https://go.dev/blog/pprof)
- [Go Benchmarks](https://pkg.go.dev/testing#hdr-Benchmarks)

---

**Version** : 1.0  
**Dernière mise à jour** : Novembre 2025  
**Type** : Optimisation de performance avec mesures rigoureuses