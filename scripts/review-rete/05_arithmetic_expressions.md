# 🔍 Revue RETE - Prompt 05: Expressions Arithmétiques

> **📋 Standards** : Ce prompt respecte les règles de [.github/prompts/common.md](../../.github/prompts/common.md) et [.github/prompts/review.md](../../.github/prompts/review.md)

**Priorité:** Moyenne  
**Durée estimée:** 2-3 heures  
**Fichiers concernés:** ~8 fichiers (~2,800 lignes)  
**Date:** 2024-12-15

---

## 📋 Vue d'ensemble

Le module d'expressions arithmétiques est responsable de :
- L'évaluation des expressions arithmétiques dans les conditions
- La décomposition d'expressions complexes en sous-expressions
- La normalisation d'expressions logiques (OR/AND imbriqués)
- Le cache des résultats d'évaluation pour performance
- L'analyse et la simplification d'expressions

Cette revue se concentre sur la qualité, la performance et la maintenabilité de ce module critique.

---

## ⚠️ Rappels Critiques

Avant de commencer, consulter obligatoirement :
- [⚠️ Standards Code Go](../../.github/prompts/common.md#standards-de-code-go) - Conventions, complexité, qualité
- [🎨 Conventions Nommage](../../.github/prompts/common.md#conventions-de-nommage) - Standards projet
- [📋 Checklist Commit](../../.github/prompts/common.md#checklist-avant-commit) - Validation
- [🔍 Revue Code](../../.github/prompts/review.md) - Process et techniques

---

## 🎯 Objectifs de cette revue

### 1. Optimiser le cache de résultats
- ✅ Analyser l'efficacité du cache (hit ratio)
- ✅ Vérifier la stabilité des clés de cache
- ✅ Mesurer l'overhead du cache
- ✅ Optimiser si hit ratio <70%

### 2. Valider la décomposition d'expressions
- ✅ Vérifier que la décomposition préserve la sémantique
- ✅ Tester les cas complexes (nested, multi-opérateurs)
- ✅ S'assurer qu'aucune sous-expression n'est perdue
- ✅ Valider l'ordre d'évaluation

### 3. Simplifier l'analyzer (complexité 28)
- ✅ Identifier la fonction `analyzeLogicalExpressionMap` (complexité ~28)
- ✅ Décomposer en sous-fonctions cohérentes (<15 chacune)
- ✅ Améliorer testabilité

### 4. Améliorer performance d'évaluation
- ✅ Identifier les bottlenecks (profiling)
- ✅ Optimiser les allocations
- ✅ Court-circuiter quand possible
- ✅ Benchmarks avant/après

### 5. Valider normalisation OR/AND
- ✅ Vérifier que les formes normalisées sont équivalentes
- ✅ Tester les équivalences logiques (De Morgan, etc.)
- ✅ S'assurer de la stabilité de la normalisation

### 6. Garantir encapsulation et généricité
- ✅ Minimiser exports publics (privé par défaut)
- ✅ Éliminer tout hardcoding
- ✅ Rendre le code générique et réutilisable

---

## 📂 Périmètre des fichiers

```
rete/arithmetic_expression_decomposer.go           # Décomposition expressions
rete/arithmetic_result_cache.go                    # Cache résultats
rete/arithmetic_decomposition_metrics.go           # Métriques décomposition
rete/arithmetic_decomposition_metrics_helpers.go   # Helpers métriques
rete/expression_analyzer.go                        # ⚠️ COMPLEXITÉ 28!
rete/nested_or_normalizer.go                       # Normalisation OR/AND
rete/nested_or_normalizer_analysis.go              # Analyse normalisation
rete/arithmetic_evaluator.go                       # Évaluation expressions
```

---

## ✅ Checklist détaillée

### 🏗️ Architecture et Design (Principes SOLID)

- [ ] **Single Responsibility Principle**
  - Decomposer, Cache, Evaluator, Normalizer séparés
  - Chaque fichier = une responsabilité unique
  - Pas de "God Objects"

- [ ] **Open/Closed Principle**
  - Extensible sans modifier code existant
  - Nouveaux opérateurs ajoutables facilement
  - Interfaces pour abstraction

- [ ] **Liskov Substitution Principle**
  - Toutes implémentations respectent contrats
  - Pas de comportements surprenants

- [ ] **Interface Segregation Principle**
  - Interfaces petites et focalisées
  - Clients ne dépendent que du nécessaire

- [ ] **Dependency Inversion Principle**
  - Dépendances sur interfaces
  - Injection de dépendances
  - Pas de dépendances globales

### 🔒 Encapsulation et Visibilité

- [ ] **Variables et fonctions privées par défaut**
  - Tous symboles privés sauf nécessité absolue
  - Exports publics justifiés et documentés
  - Implémentation interne cachée

- [ ] **Minimiser exports publics**
  - Seules interfaces/types du contrat public exportés
  - Helpers/utilitaires privés
  - Structures internes privées

- [ ] **Contrats d'interface respectés**
  - API publique stable
  - Breaking changes documentés

### 🚫 Anti-Hardcoding (CRITIQUE)

- [ ] **Aucune valeur hardcodée**
  - Pas de magic numbers
  - Pas de magic strings
  - Pas de limites hardcodées (profondeur, taille, etc.)
  - Pas de timeouts hardcodés

- [ ] **Constantes nommées et explicites**
  ```go
  // ❌ MAUVAIS
  if depth > 10 { return errTooDeep }
  if len(expr) > 100 { return errTooLong }
  
  // ✅ BON
  const (
      MaxExpressionDepth = 10
      MaxExpressionLength = 100
  )
  if depth > MaxExpressionDepth { return errTooDeep }
  if len(expr) > MaxExpressionLength { return errTooLong }
  ```

- [ ] **Code générique et paramétrable**
  - Paramètres de fonction pour valeurs variables
  - Interfaces pour opérateurs/évaluateurs
  - Configuration via structures
  - Pas de code spécifique à un opérateur

### 🧪 Tests Fonctionnels RÉELS (CRITIQUE)

- [ ] **Pas de simulation/mocks**
  - Tests évaluent vraiment les expressions
  - Résultats réels comparés aux attendus
  - Pas de suppositions sur les résultats

- [ ] **Tests déterministes et isolés**
  - Chaque test indépendant
  - Pas de dépendances entre tests
  - Setup/teardown propre
  - Reproductibles

- [ ] **Couverture > 80%**
  - Cas nominaux
  - Cas limites (division par zéro, overflow, etc.)
  - Cas d'erreur
  - Edge cases (expressions vides, très longues, etc.)

- [ ] **Tests équivalences mathématiques**
  - Commutativité: `a + b = b + a`
  - Associativité: `(a+b)+c = a+(b+c)`
  - Distributivité: `a*(b+c) = a*b + a*c`
  - Identités: `a + 0 = a`, `a * 1 = a`
  - De Morgan: `!(a && b) = !a || !b`

### 📋 Qualité du Code

- [ ] **Complexité cyclomatique < 15**
  - ⚠️ CRITIQUE: Décomposer `analyzeLogicalExpressionMap` (complexité 28)
  - Toutes autres fonctions <15 (idéalement <10)
  - Extract Function pattern

- [ ] **Fonctions < 50 lignes**
  - Sauf justification documentée
  - Décomposer fonctions longues
  - Une fonction = une responsabilité

- [ ] **Imbrication < 4 niveaux**
  - Pas de deep nesting
  - Early return
  - Extract Function

- [ ] **Pas de duplication (DRY)**
  - Code partagé extrait
  - Composition/interfaces
  - Constantes pour valeurs répétées

- [ ] **Noms explicites et idiomatiques**
  - Variables: camelCase descriptif
  - Fonctions: MixedCaps, verbes
  - Types: MixedCaps, noms
  - Constantes: MixedCaps ou UPPER_CASE
  - Éviter abréviations: `expr` → `expression`, `eval` → `evaluate`

- [ ] **Code auto-documenté**
  - Code lisible comme du texte
  - Logique claire
  - Commentaires seulement si algorithme complexe

### 🔐 Sécurité et Robustesse

- [ ] **Validation des entrées**
  - Expressions nulles/vides gérées
  - Profondeur excessive détectée (stack overflow)
  - Types validés
  - Pas de panic

- [ ] **Gestion d'erreurs robuste**
  - Erreurs propagées avec contexte
  - Messages informatifs (quelle expression, où, pourquoi)
  - Pas de suppression silencieuse
  - Return early on error

- [ ] **Protection contre overflow/underflow**
  - Validation des opérations arithmétiques
  - Détection overflow avant calcul si possible
  - Gestion division par zéro

- [ ] **Thread-safety si nécessaire**
  - Cache thread-safe si accès concurrent
  - Synchronisation correcte (RWMutex pour cache)
  - Tests race detector
  - Pas de race conditions

- [ ] **Ressources libérées proprement**
  - Pas de fuites mémoire
  - Defer pour cleanup
  - Éviction cache si limite mémoire

### 📚 Documentation

- [ ] **En-tête copyright présent**
  ```go
  // Copyright (c) 2025 TSD Contributors
  // Licensed under the MIT License
  // See LICENSE file in the project root for full license text
  ```

- [ ] **GoDoc pour tous exports**
  - Fonctions exportées documentées
  - Types exportés documentés
  - Exemples pour API complexes

- [ ] **Commentaires inline si complexe**
  - Algorithmes mathématiques expliqués
  - Normalisation logique documentée
  - Références si formules connues

- [ ] **Pas de commentaires obsolètes**
  - Supprimer code commenté
  - MAJ après changements
  - Pas de redondance

### ⚡ Performance

- [ ] **Cache efficace**
  - Hit ratio >70% mesuré et documenté
  - Clés de cache stables et déterministes
  - Overhead cache <1% temps total
  - Éviction intelligente si limite mémoire

- [ ] **Évaluation optimisée**
  - Court-circuit pour AND/OR (`&&`, `||`)
  - Pas de calculs redondants
  - Mémoïzation si sous-expressions répétées

- [ ] **Allocations minimisées**
  - Réutilisation d'objets si possible
  - Pas de copies inutiles
  - Slices pré-alloués si taille connue

- [ ] **Décomposition efficace**
  - Pas de décomposition inutile (expressions simples)
  - Partage de sous-expressions communes
  - Overhead décomposition mesuré

### 🎨 Expressions Arithmétiques (Spécifique)

- [ ] **Opérateurs supportés clairement documentés**
  - Liste complète : `+`, `-`, `*`, `/`, `%` (si supporté), etc.
  - Précédence correcte
  - Associativité correcte

- [ ] **Décomposition préserve sémantique**
  - Tests exhaustifs de non-régression
  - Équivalence mathématique validée
  - Ordre d'évaluation respecté

- [ ] **Normalisation logique correcte**
  - Formes normales équivalentes
  - De Morgan appliqué correctement
  - Double négation éliminée
  - Tests avec tables de vérité

- [ ] **Gestion types numériques**
  - int, int64, float64 supportés
  - Conversions explicites et documentées
  - Pas de perte de précision non documentée

---

## 🚫 Anti-Patterns à Détecter et Éliminer

- [ ] **God Object** - Analyzer fait tout
  - ⚠️ Chercher `expression_analyzer.go`
  - Diviser responsabilités
  - Extract Function

- [ ] **Long Method** - Fonctions >50-100 lignes
  - ⚠️ La fonction à complexité 28 est probablement longue
  - Extract Function
  - Décomposer en étapes

- [ ] **Long Parameter List** - >5 paramètres
  - Utiliser structure d'options
  - Grouper paramètres liés

- [ ] **Magic Numbers/Strings** - Hardcoding
  - Extract Constant
  - Constantes nommées

- [ ] **Duplicate Code** - Répétition
  - Extract Function
  - Composition

- [ ] **Dead Code** - Code inutilisé
  - Supprimer

- [ ] **Deep Nesting** - >4 niveaux
  - Early return
  - Extract Function

- [ ] **Type Checking Instead of Polymorphism**
  - Switch sur types → interfaces
  - Polymorphisme

---

## 🔧 Commandes de validation

### Tests

```bash
# Tests expressions arithmétiques
go test -v ./rete -run "TestArithmetic"
go test -v ./rete -run "TestExpression"

# Tests décomposition
go test -v ./rete -run "TestDecompos"

# Tests normalisation
go test -v ./rete -run "TestNormali"
go test -v ./rete -run "TestNestedOr"

# Tests cache
go test -v ./rete -run "TestCache"

# Tous tests avec couverture
go test -coverprofile=coverage_arith.out ./rete -run "TestArithmetic|TestExpression|TestDecompos|TestNormali|TestCache"
go tool cover -func=coverage_arith.out
go tool cover -html=coverage_arith.out -o coverage_arith.html

# Race detector
go test -race ./rete -run "TestArithmetic|TestCache"
```

### Performance

```bash
# Benchmarks expressions
go test -bench=BenchmarkArithmetic -benchmem ./rete
go test -bench=BenchmarkExpression -benchmem ./rete

# Benchmarks cache (mesurer hit ratio)
go test -bench=BenchmarkCache -benchmem ./rete

# Benchmarks décomposition
go test -bench=BenchmarkDecompos -benchmem ./rete

# Profiling CPU
go test -bench=BenchmarkArithmetic -cpuprofile=cpu_arith.prof ./rete
go tool pprof -http=:8080 cpu_arith.prof

# Profiling mémoire
go test -bench=BenchmarkArithmetic -memprofile=mem_arith.prof ./rete
go tool pprof -http=:8080 mem_arith.prof
```

### Qualité

```bash
# Complexité (CRITIQUE: trouver la fonction à 28)
gocyclo -over 15 rete/arithmetic*.go rete/expression*.go rete/nested*.go
gocyclo -top 20 rete/arithmetic*.go rete/expression*.go rete/nested*.go

# Vérifications statiques (obligatoires)
go vet ./rete/arithmetic*.go ./rete/expression*.go ./rete/nested*.go
staticcheck ./rete/arithmetic*.go ./rete/expression*.go ./rete/nested*.go
errcheck ./rete/arithmetic*.go ./rete/expression*.go ./rete/nested*.go
gosec ./rete/arithmetic*.go ./rete/expression*.go ./rete/nested*.go

# Formatage (obligatoire avant commit)
gofmt -l rete/arithmetic*.go rete/expression*.go rete/nested*.go
go fmt ./rete/arithmetic*.go ./rete/expression*.go ./rete/nested*.go
goimports -w rete/arithmetic*.go rete/expression*.go rete/nested*.go

# Linting complet
golangci-lint run ./rete/arithmetic*.go ./rete/expression*.go ./rete/nested*.go

# Validation complète
make validate
```

### Vérification Copyright

```bash
# Vérifier en-têtes
for file in rete/arithmetic*.go rete/expression*.go rete/nested*.go; do
    if ! head -1 "$file" | grep -q "Copyright"; then
        echo "⚠️  COPYRIGHT MANQUANT: $file"
    fi
done
```

---

## 🔄 Processus de revue recommandé

### Phase 1: Analyse initiale (30-45 min)

1. **Mesurer baseline actuelle**
   ```bash
   mkdir -p REPORTS/review-rete
   
   # Complexité (TROUVER la fonction à 28!)
   gocyclo -over 10 rete/arithmetic*.go rete/expression*.go rete/nested*.go > REPORTS/review-rete/arith_complexity_before.txt
   echo "=== TOP COMPLEXITÉ ==="
   gocyclo -top 20 rete/arithmetic*.go rete/expression*.go rete/nested*.go
   
   # Couverture
   go test -coverprofile=REPORTS/review-rete/arith_coverage_before.out ./rete -run "TestArithmetic|TestExpression" 2>/dev/null
   go tool cover -func=REPORTS/review-rete/arith_coverage_before.out > REPORTS/review-rete/arith_coverage_before.txt
   
   # Benchmarks
   go test -bench=BenchmarkArithmetic -benchmem ./rete > REPORTS/review-rete/arith_benchmarks_before.txt 2>&1
   ```

2. **Lire fichiers dans ordre logique**
   - `arithmetic_evaluator.go` (fondation - évaluation)
   - `arithmetic_expression_decomposer.go` (décomposition)
   - `arithmetic_result_cache.go` (cache)
   - `arithmetic_decomposition_metrics.go` (métriques)
   - `arithmetic_decomposition_metrics_helpers.go` (helpers)
   - `expression_analyzer.go` (⚠️ COMPLEXITÉ 28!)
   - `nested_or_normalizer.go` (normalisation)
   - `nested_or_normalizer_analysis.go` (analyse)

3. **Pour chaque fichier, vérifier**
   - [ ] En-tête copyright présent?
   - [ ] Exports minimaux?
   - [ ] Aucun hardcoding?
   - [ ] Code générique?
   - [ ] Complexité <15? (⚠️ identifier la 28)
   - [ ] Noms explicites (pas `expr`, `eval` abrégés)?
   - [ ] Tests présents?
   - [ ] GoDoc complet?
   - [ ] Anti-patterns?

4. **Analyser cache**
   ```bash
   # Chercher métriques de cache
   grep -n "hit" rete/arithmetic_result_cache.go
   grep -n "miss" rete/arithmetic_result_cache.go
   grep -n "CacheHit\|CacheMiss" rete/arithmetic*.go
   ```

### Phase 2: Identification des problèmes (30-45 min)

**Créer liste priorisée dans** `REPORTS/review-rete/05_arith_issues.md`:

```markdown
# Problèmes Identifiés - Expressions Arithmétiques

## P0 - BLOQUANT

### 1. [Si bugs détectés]
- **Fichier:** arithmetic_evaluator.go:XXX
- **Type:** Bug évaluation / Division par zéro non gérée
- **Impact:** Crash runtime
- **Solution:** ...

## P1 - IMPORTANT

### 1. Complexité 28 dans expression_analyzer.go
- **Fichier:** expression_analyzer.go:XXX
- **Fonction:** `analyzeLogicalExpressionMap`
- **Type:** Complexité excessive (28)
- **Impact:** Impossible à maintenir
- **Solution:** Extract Function - décomposer en 4-5 sous-fonctions

### 2. Cache hit ratio faible
- **Fichier:** arithmetic_result_cache.go
- **Métriques:** Hit ratio <70% (si mesuré)
- **Impact:** Performance sous-optimale
- **Solution:** Améliorer clés de cache, stabilité

### 3. [Hardcoding détecté]
- **Fichiers:** Multiples
- **Type:** Magic numbers (limites, profondeurs, etc.)
- **Impact:** Pas configurable, rigide
- **Solution:** Extract Constant

## P2 - SOUHAITABLE
...
```

**Problèmes à chercher:**

**P0 - Bloquant:**
- Division par zéro non gérée
- Overflow non détecté
- Panic dans évaluation
- Race conditions (cache)
- Bug normalisation (équivalence cassée)

**P1 - Important:**
- **Complexité 28 (PRIORITÉ)**
- Autres complexités 15-28
- Cache inefficace (hit ratio <70%)
- Hardcoding limites/profondeurs
- Exports non justifiés
- Couverture <70%
- Missing copyright

**P2 - Souhaitable:**
- Complexité 10-15
- Optimisations mineures
- Refactoring clarté

### Phase 3: Corrections (60-90 min)

#### 3.1 Fixer P0 (bloquants)

**Exemple: Division par zéro**

```go
// AVANT
func divide(a, b float64) float64 {
    return a / b  // ❌ Panic si b == 0
}

// APRÈS
const DivisionByZeroError = "division by zero"

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf(DivisionByZeroError)
    }
    return a / b, nil
}
```

**Tests:**
```go
func TestDivide_ZeroDivisor(t *testing.T) {
    _, err := divide(10, 0)
    require.Error(t, err)
    assert.Contains(t, err.Error(), "division by zero")
}
```

**Commit:**
```bash
git commit -m "[Review-05/Arith] fix(P0): gère division par zéro

- Ajoute validation diviseur != 0
- Retourne erreur explicite
- Évite panic runtime
- Tests edge case ajoutés

Resolves: P0-arith-division-zero
Refs: scripts/review-rete/05_arithmetic_expressions.md"
```

#### 3.2 Décomposer fonction à complexité 28 (P1 PRIORITÉ)

**Identifier la fonction:**
```bash
gocyclo -over 25 rete/expression*.go
# Probablement analyzeLogicalExpressionMap dans expression_analyzer.go
```

**Pattern de décomposition:**

```go
// AVANT - complexité 28, ~120 lignes
func analyzeLogicalExpressionMap(expr LogicalExpr) (*AnalysisResult, error) {
    // 30 lignes parsing
    // 40 lignes validation
    // 30 lignes transformation
    // 20 lignes génération résultat
}

// APRÈS - décomposer

func analyzeLogicalExpressionMap(expr LogicalExpr) (*AnalysisResult, error) {
    // Orchestration - complexité ~8
    parsed := parseLogicalExpression(expr)
    
    if err := validateLogicalStructure(parsed); err != nil {
        return nil, err
    }
    
    normalized := normalizeLogicalForm(parsed)
    simplified := simplifyLogicalExpression(normalized)
    
    return buildAnalysisResult(simplified), nil
}

// Chaque sous-fonction <12 complexité
func parseLogicalExpression(expr LogicalExpr) *ParsedExpr {
    // Complexité ~10
}

func validateLogicalStructure(parsed *ParsedExpr) error {
    // Complexité ~8
}

func normalizeLogicalForm(parsed *ParsedExpr) *NormalizedExpr {
    // Complexité ~12
}

func simplifyLogicalExpression(norm *NormalizedExpr) *SimplifiedExpr {
    // Complexité ~9
}

func buildAnalysisResult(simplified *SimplifiedExpr) *AnalysisResult {
    // Complexité ~6
}
```

**Tests:**
```go
func TestParseLogicalExpression(t *testing.T) { /* ... */ }
func TestValidateLogicalStructure(t *testing.T) { /* ... */ }
func TestNormalizeLogicalForm(t *testing.T) { /* ... */ }
func TestSimplifyLogicalExpression(t *testing.T) { /* ... */ }
func TestBuildAnalysisResult(t *testing.T) { /* ... */ }
```

**Commit:**
```bash
git commit -m "[Review-05/Arith] refactor(P1): décompose analyzeLogicalExpressionMap (28→8)

- Extrait parseLogicalExpression() (complexité 10)
- Extrait validateLogicalStructure() (complexité 8)
- Extrait normalizeLogicalForm() (complexité 12)
- Extrait simplifyLogicalExpression() (complexité 9)
- Extrait buildAnalysisResult() (complexité 6)
- Orchestration: complexité 8
- Tests unitaires pour chaque fonction

Resolves: P1-arith-complexity-28
Refs: scripts/review-rete/05_arithmetic_expressions.md"
```

#### 3.3 Optimiser cache (P1)

**Analyser hit ratio:**

```go
// Ajouter instrumentation si absente
type CacheMetrics struct {
    Hits   int64
    Misses int64
    mu     sync.RWMutex
}

func (m *CacheMetrics) HitRatio() float64 {
    m.mu.RLock()
    defer m.mu.RUnlock()
    total := m.Hits + m.Misses
    if total == 0 {
        return 0
    }
    return float64(m.Hits) / float64(total)
}
```

**Tests cache:**
```go
func TestCache_HitRatio(t *testing.T) {
    cache := NewArithmeticCache()
    
    // Première évaluation - miss
    result1, _ := cache.Evaluate("x + y", bindings1)
    
    // Même expression, mêmes bindings - hit
    result2, _ := cache.Evaluate("x + y", bindings1)
    
    assert.Equal(t, result1, result2)
    assert.GreaterOrEqual(t, cache.Metrics().HitRatio(), 0.5)
}
```

**Améliorer stabilité clés:**
```go
// AVANT - non déterministe si map
func cacheKey(expr string, bindings map[string]interface{}) string {
    return fmt.Sprintf("%s:%v", expr, bindings)  // ❌ map order
}

// APRÈS - déterministe
func cacheKey(expr string, bindings map[string]interface{}) string {
    keys := make([]string, 0, len(bindings))
    for k := range bindings {
        keys = append(keys, k)
    }
    sort.Strings(keys)  // ✅ ordre stable
    
    var buf strings.Builder
    buf.WriteString(expr)
    buf.WriteString(":")
    for _, k := range keys {
        fmt.Fprintf(&buf, "%s=%v;", k, bindings[k])
    }
    return buf.String()
}
```

#### 3.4 Éliminer hardcoding (P1)

```go
// AVANT
if depth > 10 { return errTooDeep }
if len(expr) > 100 { return errTooLong }
timeout := 30 * time.Second

// APRÈS
const (
    MaxExpressionDepth     = 10
    MaxExpressionLength    = 100
    ExpressionEvalTimeout  = 30 * time.Second
)

if depth > MaxExpressionDepth { return errTooDeep }
if len(expr) > MaxExpressionLength { return errTooLong }
timeout := ExpressionEvalTimeout
```

**Commits atomiques pour chaque fix.**

### Phase 4: Validation finale (15-30 min)

```bash
#!/bin/bash
echo "=== VALIDATION FINALE ARITHMETIC ==="

# 1. Tests
echo "🧪 Tests..."
go test -v ./rete -run "TestArithmetic|TestExpression|TestDecompos|TestNormali"
TESTS=$?

# 2. Race detector
echo "🏁 Race detector..."
go test -race ./rete -run "TestArithmetic|TestCache"
RACE=$?

# 3. Complexité
echo "📊 Complexité..."
COMPLEX=$(gocyclo -over 15 rete/arithmetic*.go rete/expression*.go rete/nested*.go | wc -l)

# 4. Couverture
echo "📈 Couverture..."
go test -coverprofile=arith_final.out ./rete -run "TestArithmetic|TestExpression" 2>/dev/null
COVERAGE=$(go tool cover -func=arith_final.out | tail -1 | awk '{print $3}' | sed 's/%//')

# 5. Copyright
echo "©️  Copyright..."
MISSING_COPYRIGHT=0
for file in rete/arithmetic*.go rete/expression*.go rete/nested*.go; do
    if ! head -1 "$file" | grep -q "Copyright"; then
        MISSING_COPYRIGHT=$((MISSING_COPYRIGHT + 1))
        echo "  ⚠️  $file"
    fi
done

# 6. Validation complète
echo "✅ Validation..."
make validate
VALIDATE=$?

# Résumé
echo ""
echo "=== RÉSULTATS ==="
[ $TESTS -eq 0 ] && echo "✅ Tests: PASS" || echo "❌ Tests: FAIL"
[ $RACE -eq 0 ] && echo "✅ Race: PASS" || echo "❌ Race: FAIL"
[ $COMPLEX -eq 0 ] && echo "✅ Complexité: OK" || echo "❌ Complexité: $COMPLEX >15"
[ $(echo "$COVERAGE >= 80" | bc -l) -eq 1 ] && echo "✅ Couverture: $COVERAGE%" || echo "❌ Couverture: $COVERAGE%"
[ $MISSING_COPYRIGHT -eq 0 ] && echo "✅ Copyright: OK" || echo "❌ Copyright: $MISSING_COPYRIGHT manquants"
[ $VALIDATE -eq 0 ] && echo "✅ Validation: PASS" || echo "❌ Validation: FAIL"

# Verdict
if [ $TESTS -eq 0 ] && [ $RACE -eq 0 ] && [ $COMPLEX -eq 0 ] && [ $(echo "$COVERAGE >= 80" | bc -l) -eq 1 ] && [ $MISSING_COPYRIGHT -eq 0 ] && [ $VALIDATE -eq 0 ]; then
    echo ""
    echo "🎉 VALIDATION RÉUSSIE - Prêt pour Prompt 06!"
    exit 0
else
    echo ""
    echo "❌ VALIDATION ÉCHOUÉE"
    exit 1
fi
```

---

## 📝 Livrables attendus

### 1. Rapport d'analyse

**Créer:** `REPORTS/review-rete/05_arithmetic_expressions_report.md`

**Structure obligatoire:**

```markdown
# 🔍 Revue de Code : Expressions Arithmétiques

**Date:** 2024-12-XX  
**Réviseur:** [Nom]  
**Durée:** Xh Ym

---

## 📊 Vue d'Ensemble

- **Fichiers analysés:** 8
- **Lignes de code:** ~2,800
- **Complexité avant:** Max 28
- **Complexité après:** Max <15
- **Couverture avant:** X%
- **Couverture après:** Y%
- **Cache hit ratio:** X% → Y%

---

## ✅ Points Forts

- Séparation decomposer/cache/evaluator/normalizer
- Cache implémenté (à optimiser)
- ...

---

## ❌ Problèmes Identifiés et Corrigés

### P0 - BLOQUANT

#### 1. Division par zéro non gérée
- **Fichier:** arithmetic_evaluator.go:XXX
- **Type:** Bug critique
- **Impact:** Panic runtime
- **Solution:** Validation + erreur
- **Commit:** abc1234
- **Status:** ✅ Résolu

### P1 - IMPORTANT

#### 1. Complexité 28 dans expression_analyzer.go
- **Fonction:** analyzeLogicalExpressionMap
- **Solution:** Décomposé en 5 fonctions (<12 chacune)
- **Commit:** def5678
- **Status:** ✅ Résolu

#### 2. Cache hit ratio 45% (sous-optimal)
- **Problème:** Clés non stables (map order)
- **Solution:** Tri des clés, ordre déterministe
- **Hit ratio après:** 78%
- **Commit:** ghi9012
- **Status:** ✅ Résolu

---

## 🔧 Changements Apportés

### Refactoring

1. **Décomposition analyzeLogicalExpressionMap**
   - 1 fonction 120 lignes → 5 fonctions ~25 lignes
   - Complexité 28 → max 12
   - Tests unitaires: 5

2. **Optimisation cache**
   - Clés stabilisées (tri)
   - Hit ratio: 45% → 78%
   - Thread-safety confirmée (RWMutex)

3. **Élimination hardcoding**
   - 18 magic numbers → constantes
   - 6 magic strings → constantes

---

## 📈 Métriques Avant/Après

| Métrique | Avant | Après | Amélioration |
|----------|-------|-------|--------------|
| Complexité max | 28 | 12 | ✅ -57% |
| Fonctions >15 | 3 | 0 | ✅ 100% |
| Couverture | 74% | 86% | ✅ +12% |
| Cache hit ratio | 45% | 78% | ✅ +73% |
| Magic numbers | 18 | 0 | ✅ 100% |

---

## 💡 Recommandations Futures

### Court terme
1. Monitorer cache hit ratio en production
2. Benchmarks sur expressions réelles
3. Documenter opérateurs supportés

### Moyen terme
1. Support opérateur modulo (%)
2. Évaluation lazy pour performance
3. Constant folding (simplification compile-time)

---

## 🏁 Verdict

✅ **APPROUVÉ**

Complexité 28 éliminée, cache optimisé, standards respectés.
Prêt pour Prompt 06 (Builders).

---

**Prochaines étapes:**
1. Merge commits
2. Lancer Prompt 06
3. Monitorer cache metrics production
```

### 2. Tests ajoutés/améliorés

```bash
git diff --name-only | grep "_test.go" > REPORTS/arith_tests_added.txt
diff <(cat REPORTS/review-rete/arith_coverage_before.txt) \
     <(cat REPORTS/review-rete/arith_coverage_after.txt) >> REPORTS/arith_tests_added.txt
```

### 3. Commits atomiques

**Format strict:**
```
[Review-05/Arith] <type>(scope): <description courte>

- Détail 1
- Détail 2
- Resolves: <issue>

Refs: scripts/review-rete/05_arithmetic_expressions.md
```

**Exemples:**

```
[Review-05/Arith] fix(P0): gère division par zéro dans evaluator

- Validation diviseur != 0
- Erreur explicite retournée
- Évite panic runtime
- Tests edge case ajoutés

Resolves: P0-arith-division-zero
Refs: scripts/review-rete/05_arithmetic_expressions.md
```

```
[Review-05/Arith] refactor(P1): décompose analyzeLogicalExpressionMap (28→12)

- Extrait 5 sous-fonctions (complexités 6-12)
- Améliore testabilité
- Tests unitaires ajoutés
- Orchestration: complexité 8

Resolves: P1-arith-complexity-28
Refs: scripts/review-rete/05_arithmetic_expressions.md
```

```
[Review-05/Arith] perf(P1): optimise cache expressions (hit ratio 45%→78%)

- Stabilise clés de cache (tri des bindings)
- Élimine non-déterminisme map order
- Hit ratio: 45% → 78% (+73%)
- Benchmarks confirment amélioration

Resolves: P1-arith-cache-inefficient
Refs: scripts/review-rete/05_arithmetic_expressions.md
```

---

## 📊 Métriques de succès

| Métrique | Valeur Actuelle | Cible | Critique |
|----------|----------------|-------|----------|
| Complexité max | À mesurer | <15 | ⚠️ Oui |
| Fonction à 28 | Identifier | 0 | ⚠️ OUI! |
| Fonctions >15 | À mesurer | 0 | Oui |
| Couverture tests | À mesurer | >80% | Oui |
| Cache hit ratio | À mesurer | >70% | ⚠️ Oui |
| Exports publics | À mesurer | Minimal | Oui |
| Magic numbers | À mesurer | 0 | Oui |
| Copyright headers | À mesurer | 100% | Oui |
| Race detector | À mesurer | Clean | Oui |

---

## 🎓 Ressources et références

### Standards Projet
- [common.md](../../.github/prompts/common.md)
- [review.md](../../.github/prompts/review.md)
- [Makefile](../../Makefile)

### Mathématiques & Logique
- Équivalences logiques (De Morgan, distributivité)
- Formes normales (DNF, CNF)
- Simplification expressions booléennes

### Performance
- Caching strategies
- Memoization
- Lazy evaluation

### Refactoring
- [Extract Function](https://refactoring.guru/extract-method)
- [Replace Magic Number](https://refactoring.guru/replace-magic-number-with-symbolic-constant)

---

## ✅ Checklist finale avant Prompt 06

**Validation technique:**
- [ ] Tous tests expressions passent
- [ ] Race detector clean (cache)
- [ ] Aucune fonction >15
- [ ] Complexité 28 ÉLIMINÉE
- [ ] Couverture >80%
- [ ] Cache hit ratio >70%
- [ ] `make validate` passe

**Qualité code:**
- [ ] Aucun hardcoding
- [ ] Code générique
- [ ] Exports minimaux
- [ ] Constantes nommées
- [ ] Noms explicites (pas abréviations)
- [ ] Fonctions <50 lignes
- [ ] Imbrication <4 niveaux
- [ ] Pas de duplication

**Robustesse:**
- [ ] Division par zéro gérée
- [ ] Overflow détecté
- [ ] Profondeur excessive gérée
- [ ] Validation entrées

**Tests:**
- [ ] Tests réels
- [ ] Tests déterministes
- [ ] Tests équivalences mathématiques
- [ ] Tests normalisation logique

**Documentation:**
- [ ] Copyright 100%
- [ ] GoDoc complet
- [ ] Opérateurs documentés
- [ ] Rapport créé

**Commande validation finale:** (voir script Phase 4 ci-dessus)

---

## 🚀 Script d'analyse rapide

```bash
#!/bin/bash
# scripts/review-rete/analyze_arith.sh

set -e
echo "=== ANALYSE EXPRESSIONS ARITHMÉTIQUES ==="
echo ""

mkdir -p REPORTS/review-rete

# Baseline
echo "📊 Mesure baseline..."
gocyclo -over 10 rete/arithmetic*.go rete/expression*.go rete/nested*.go > REPORTS/review-rete/arith_complexity_before.txt
go test -coverprofile=REPORTS/review-rete/arith_coverage_before.out ./rete -run "TestArithmetic|TestExpression" 2>/dev/null
go tool cover -func=REPORTS/review-rete/arith_coverage_before.out > REPORTS/review-rete/arith_coverage_before.txt

echo "✅ Baseline sauvegardée"
echo ""

# CRITIQUE: Trouver fonction à 28
echo "🚨 RECHERCHE COMPLEXITÉ 28..."
gocyclo -top 20 rete/arithmetic*.go rete/expression*.go rete/nested*.go | head -15
echo ""

# Checks
echo "🔍 Vérifications..."
echo "  go vet:"
go vet ./rete/arithmetic*.go ./rete/expression*.go ./rete/nested*.go 2>&1 | grep -v "exit status" || echo "    ✓ OK"

echo "  copyright:"
MISSING=0
for file in rete/arithmetic*.go rete/expression*.go rete/nested*.go; do
    if ! head -1 "$file" | grep -q "Copyright"; then
        echo "    ❌ $file"
        MISSING=$((MISSING + 1))
    fi
done
[ $MISSING -eq 0 ] && echo "    ✓ OK"

echo ""
echo "=== Analyse terminée ==="
echo "Créer REPORTS/review-rete/05_arith_issues.md"
```

**Lancer:**
```bash
chmod +x scripts/review-rete/analyze_arith.sh
./scripts/review-rete/analyze_arith.sh
```

---

**Prêt à commencer?** 🚀

Bonne revue! Respecter scrupuleusement les standards common.md et review.md.