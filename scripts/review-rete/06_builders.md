# 🔍 Revue RETE - Prompt 06: Builders et Construction

> **📋 Standards** : Ce prompt respecte les règles de [.github/prompts/common.md](../../.github/prompts/common.md) et [.github/prompts/review.md](../../.github/prompts/review.md)

**Priorité:** Haute  
**Durée estimée:** 3-4 heures  
**Fichiers concernés:** ~12 fichiers (~3,000 lignes)  
**Date:** 2024-12-15

---

## 📋 Vue d'ensemble

Le module Builders est responsable de :
- La construction du réseau RETE à partir de règles
- La coordination entre builders Alpha, Beta, Exists, Accumulator
- L'orchestration de la construction du réseau complet
- La validation des règles avant construction
- La gestion du contexte de construction
- Le registry des composants construits

Cette revue se concentre sur la qualité, la cohérence et la maintenabilité de cette couche d'orchestration critique.

---

## ⚠️ Rappels Critiques

Avant de commencer, consulter obligatoirement :
- [⚠️ Standards Code Go](../../.github/prompts/common.md#standards-de-code-go) - Conventions, complexité, qualité
- [🎨 Conventions Nommage](../../.github/prompts/common.md#conventions-de-nommage) - Standards projet
- [📋 Checklist Commit](../../.github/prompts/common.md#checklist-avant-commit) - Validation
- [🔍 Revue Code](../../.github/prompts/review.md) - Process et techniques

---

## 🎯 Objectifs de cette revue

### 1. Valider séparation des responsabilités (SRP)
- ✅ Chaque builder a une responsabilité unique et claire
- ✅ Pas de chevauchement entre builders
- ✅ Orchestration coordonne sans faire le travail
- ✅ Registry gère uniquement l'enregistrement

### 2. Optimiser l'orchestration des builders
- ✅ Flux de construction clair et compréhensible
- ✅ Gestion d'erreurs cohérente
- ✅ Contexte minimal et explicite
- ✅ Pas de logique métier dans l'orchestration

### 3. Réduire la complexité (plusieurs >15)
- ✅ Identifier toutes les fonctions >15
- ✅ Décomposer en sous-fonctions cohérentes
- ✅ Améliorer testabilité
- ✅ Target <12 pour fonctions critiques

### 4. Éliminer duplication entre builders
- ✅ Identifier patterns communs
- ✅ Extraire dans fonctions partagées
- ✅ Éviter copier-coller
- ✅ DRY strict

### 5. Améliorer testabilité
- ✅ Chaque builder testable indépendamment
- ✅ Mocks/stubs si nécessaire pour isolation
- ✅ Tests d'intégration pour orchestration
- ✅ Couverture >85%

### 6. Garantir encapsulation et généricité
- ✅ Minimiser exports publics (privé par défaut)
- ✅ Éliminer tout hardcoding
- ✅ Rendre le code générique et réutilisable

---

## 📂 Périmètre des fichiers

```
rete/builder_rules.go                              # Builder règles principal
rete/builder_alpha_rules.go                        # Builder chaînes Alpha
rete/builder_join_rules.go                         # Builder jointures
rete/builder_join_rules_binary_orchestration.go    # Orchestration binaire
rete/builder_exists_rules.go                       # Builder EXISTS patterns
rete/builder_accumulator_rules.go                  # Builder ACCUMULATE patterns
rete/builder_types.go                              # Types builders
rete/builder_utils.go                              # Utilitaires builders
rete/builder_orchestration.go                      # Orchestration globale
rete/builder_registry.go                           # Registry composants
rete/builder_context.go                            # Contexte construction
rete/builder_validation.go                         # Validation règles
```

---

## ✅ Checklist détaillée

### 🏗️ Architecture et Design (Principes SOLID)

- [ ] **Single Responsibility Principle**
  - Chaque builder : une responsabilité unique
  - AlphaBuilder → Alpha seulement
  - JoinBuilder → Jointures seulement
  - ExistsBuilder → EXISTS seulement
  - AccumulatorBuilder → ACCUMULATE seulement
  - Orchestrator → coordination seulement
  - Pas de "God Builders"

- [ ] **Open/Closed Principle**
  - Extensible sans modifier existant
  - Nouveaux types de patterns ajoutables
  - Interfaces pour abstraction

- [ ] **Liskov Substitution Principle**
  - Toutes implémentations respectent contrats
  - Pas de comportements surprenants

- [ ] **Interface Segregation Principle**
  - Interfaces petites et focalisées
  - Pas d'interface monolithique `Builder`
  - Clients dépendent du minimum

- [ ] **Dependency Inversion Principle**
  - Dépendances sur interfaces
  - Injection de dépendances
  - Pas de dépendances hardcodées

- [ ] **Séparation des préoccupations**
  - Construction ≠ Validation
  - Construction ≠ Optimisation
  - Construction ≠ Enregistrement
  - Chaque fichier = une préoccupation

### 🔒 Encapsulation et Visibilité

- [ ] **Variables et fonctions privées par défaut**
  - Tous symboles privés sauf nécessité absolue
  - Exports publics justifiés et documentés
  - Implémentation interne cachée

- [ ] **Minimiser exports publics**
  - Seules interfaces de builders exportées
  - Types internes privés
  - Helpers privés
  - Context interne si possible

- [ ] **Contrats d'interface respectés**
  - API publique stable
  - Breaking changes documentés

### 🚫 Anti-Hardcoding (CRITIQUE)

- [ ] **Aucune valeur hardcodée**
  - Pas de magic numbers
  - Pas de magic strings
  - Pas de limites hardcodées (nombre de patterns, profondeur, etc.)
  - Pas de noms de types hardcodés

- [ ] **Constantes nommées et explicites**
  ```go
  // ❌ MAUVAIS
  if len(patterns) > 20 { return errTooMany }
  
  // ✅ BON
  const MaxPatternsPerRule = 20
  if len(patterns) > MaxPatternsPerRule { return errTooMany }
  ```

- [ ] **Code générique et paramétrable**
  - Builders paramétrés par type
  - Pas de code spécifique à une règle
  - Configuration via options/structures

### 🧪 Tests Fonctionnels RÉELS (CRITIQUE)

- [ ] **Pas de simulation/mocks**
  - Tests construisent vraiment le réseau
  - Vérification de la structure construite
  - Pas de suppositions
  - SAUF: mocks pour isolation de tests unitaires (acceptable)

- [ ] **Tests déterministes et isolés**
  - Chaque test indépendant
  - Setup/teardown propre
  - Reproductibles

- [ ] **Couverture > 85%**
  - Cas nominaux
  - Cas limites
  - Cas d'erreur
  - Edge cases

- [ ] **Tests par builder**
  - Tests unitaires AlphaBuilder
  - Tests unitaires JoinBuilder
  - Tests unitaires ExistsBuilder
  - Tests unitaires AccumulatorBuilder
  - Tests intégration Orchestration

### 📋 Qualité du Code

- [ ] **Complexité cyclomatique < 15**
  - Toutes fonctions <15 (idéalement <10)
  - Identifier toutes >15
  - Décomposer systématiquement
  - Extract Function pattern

- [ ] **Fonctions < 50 lignes**
  - Sauf justification documentée
  - Décomposer fonctions longues
  - Une fonction = une étape claire

- [ ] **Imbrication < 4 niveaux**
  - Pas de deep nesting
  - Early return
  - Extract Function

- [ ] **Pas de duplication (DRY)**
  - Patterns communs extraits
  - Éviter copier-coller entre builders
  - Helpers partagés pour logique commune
  - Constantes pour valeurs répétées

- [ ] **Noms explicites et idiomatiques**
  - Variables: camelCase descriptif
  - Fonctions: MixedCaps, verbes (buildAlpha, validateRule)
  - Types: MixedCaps, noms (AlphaBuilder, BuildContext)
  - Constantes: MixedCaps ou UPPER_CASE
  - Pas d'abréviations: `bldr` → `builder`, `ctx` → `context` (sauf contexte Go standard)

- [ ] **Code auto-documenté**
  - Code lisible comme du texte
  - Logique claire
  - Commentaires si algorithme complexe

### 🔐 Sécurité et Robustesse

- [ ] **Validation des entrées**
  - Règles nulles/vides gérées
  - Patterns validés
  - Types validés
  - Pas de panic

- [ ] **Gestion d'erreurs robuste**
  - Erreurs propagées avec contexte
  - Messages informatifs (quelle règle, quel pattern, pourquoi)
  - Pas de suppression silencieuse
  - Return early on error
  - Wrap errors avec contexte

- [ ] **Thread-safety si nécessaire**
  - Registry thread-safe si accès concurrent
  - Context immutable ou protégé
  - Tests race detector
  - Documentation des garanties

- [ ] **Ressources libérées proprement**
  - Pas de fuites mémoire
  - Defer pour cleanup
  - Context pour annulation si long

### 📚 Documentation

- [ ] **En-tête copyright présent**
  ```go
  // Copyright (c) 2025 TSD Contributors
  // Licensed under the MIT License
  // See LICENSE file in the project root for full license text
  ```

- [ ] **GoDoc pour tous exports**
  - Builders documentés
  - Interfaces documentées
  - Types exportés documentés
  - Exemples si API complexe

- [ ] **Commentaires inline si complexe**
  - Algorithmes construction expliqués
  - Justification choix d'implémentation
  - Références à patterns connus

- [ ] **Pas de commentaires obsolètes**
  - Supprimer code commenté
  - MAJ après changements
  - Pas de redondance

### ⚡ Performance

- [ ] **Construction efficace**
  - Pas de reconstructions inutiles
  - Réutilisation de composants (sharing)
  - Allocations minimisées

- [ ] **Validation précoce**
  - Valider avant construire
  - Fail fast sur erreurs
  - Éviter travail inutile

- [ ] **Registry efficace**
  - Lookups rapides (maps)
  - Pas de scans linéaires
  - Overhead minimal

- [ ] **Context léger**
  - État minimal
  - Pas de copies inutiles
  - Passage par référence si gros

### 🎨 Builders (Spécifique)

- [ ] **Séparation claire des builders**
  - AlphaBuilder → chaînes Alpha uniquement
  - JoinBuilder → nœuds Join uniquement
  - ExistsBuilder → conditions EXISTS uniquement
  - AccumulatorBuilder → accumulations uniquement
  - Pas de chevauchement de responsabilités

- [ ] **Orchestration simple**
  - Coordination sans logique métier
  - Délégation aux builders spécialisés
  - Gestion d'erreurs cohérente
  - Flux compréhensible

- [ ] **Validation cohérente**
  - Validation centralisée ou déléguée clairement
  - Messages d'erreur uniformes
  - Niveaux de validation (syntaxe, sémantique, etc.)

- [ ] **Registry bien défini**
  - Enregistrement uniquement
  - Lookups simples
  - Thread-safe si nécessaire
  - Pas de logique métier

- [ ] **Context approprié**
  - État minimal nécessaire
  - Immutable si possible
  - Scope clair (par règle, par réseau, global)
  - Documentation de la durée de vie

---

## 🚫 Anti-Patterns à Détecter et Éliminer

- [ ] **God Object** - Builder qui fait tout
  - Chercher builders >500 lignes
  - Diviser responsabilités
  - Séparer par type de construction

- [ ] **Long Method** - Fonctions >50-100 lignes
  - Extract Function
  - Décomposer en étapes

- [ ] **Long Parameter List** - >5 paramètres
  - Utiliser BuildContext/Options
  - Grouper paramètres liés

- [ ] **Magic Numbers/Strings** - Hardcoding
  - Extract Constant
  - Constantes nommées

- [ ] **Duplicate Code** - Copier-coller entre builders
  - Extract Function
  - Helpers partagés
  - Composition

- [ ] **Dead Code** - Code inutilisé
  - Supprimer

- [ ] **Deep Nesting** - >4 niveaux
  - Early return
  - Extract Function

- [ ] **Feature Envy** - Builder accède trop à un autre
  - Déplacer logique
  - Encapsulation

- [ ] **Shotgun Surgery** - Changement éparpillé
  - Centraliser logique
  - Composition

---

## 🔧 Commandes de validation

### Tests

```bash
# Tests builders
go test -v ./rete -run "TestBuilder"
go test -v ./rete -run "TestAlphaBuilder"
go test -v ./rete -run "TestJoinBuilder"
go test -v ./rete -run "TestExistsBuilder"
go test -v ./rete -run "TestAccumulatorBuilder"

# Tests orchestration
go test -v ./rete -run "TestOrchestration"
go test -v ./rete -run "TestBuildRule"

# Tests validation
go test -v ./rete -run "TestValidation"

# Tous tests avec couverture
go test -coverprofile=coverage_builders.out ./rete -run "TestBuilder|TestOrchestration|TestValidation"
go tool cover -func=coverage_builders.out
go tool cover -html=coverage_builders.out -o coverage_builders.html

# Race detector
go test -race ./rete -run "TestBuilder|TestRegistry"
```

### Performance

```bash
# Benchmarks construction
go test -bench=BenchmarkBuild -benchmem ./rete

# Benchmarks builders spécifiques
go test -bench=BenchmarkAlphaBuilder -benchmem ./rete
go test -bench=BenchmarkJoinBuilder -benchmem ./rete

# Profiling
go test -bench=BenchmarkBuild -cpuprofile=cpu_build.prof ./rete
go tool pprof -http=:8080 cpu_build.prof
```

### Qualité

```bash
# Complexité (identifier toutes >15)
gocyclo -over 15 rete/builder*.go
gocyclo -top 20 rete/builder*.go

# Vérifications statiques
go vet ./rete/builder*.go
staticcheck ./rete/builder*.go
errcheck ./rete/builder*.go
gosec ./rete/builder*.go

# Formatage
gofmt -l rete/builder*.go
go fmt ./rete/builder*.go
goimports -w rete/builder*.go

# Linting
golangci-lint run ./rete/builder*.go

# Validation complète
make validate
```

### Vérification Copyright

```bash
for file in rete/builder*.go; do
    if ! head -1 "$file" | grep -q "Copyright"; then
        echo "⚠️  COPYRIGHT MANQUANT: $file"
    fi
done
```

---

## 🔄 Processus de revue recommandé

### Phase 1: Analyse initiale (45-60 min)

1. **Mesurer baseline**
   ```bash
   mkdir -p REPORTS/review-rete
   
   # Complexité
   gocyclo -over 10 rete/builder*.go > REPORTS/review-rete/builders_complexity_before.txt
   echo "=== TOP COMPLEXITÉ ==="
   gocyclo -top 30 rete/builder*.go
   
   # Couverture
   go test -coverprofile=REPORTS/review-rete/builders_coverage_before.out ./rete -run "TestBuilder" 2>/dev/null
   go tool cover -func=REPORTS/review-rete/builders_coverage_before.out > REPORTS/review-rete/builders_coverage_before.txt
   
   # Benchmarks
   go test -bench=BenchmarkBuild -benchmem ./rete > REPORTS/review-rete/builders_benchmarks_before.txt 2>&1
   ```

2. **Lire fichiers dans ordre logique**
   - `builder_types.go` (types de base)
   - `builder_context.go` (contexte)
   - `builder_validation.go` (validation)
   - `builder_utils.go` (utilitaires)
   - `builder_alpha_rules.go` (builder Alpha)
   - `builder_join_rules.go` (builder Join)
   - `builder_join_rules_binary_orchestration.go` (orchestration Join)
   - `builder_exists_rules.go` (builder EXISTS)
   - `builder_accumulator_rules.go` (builder ACCUMULATE)
   - `builder_rules.go` (builder principal)
   - `builder_orchestration.go` (orchestration globale)
   - `builder_registry.go` (registry)

3. **Pour chaque fichier, vérifier**
   - [ ] Copyright présent?
   - [ ] Exports minimaux?
   - [ ] Aucun hardcoding?
   - [ ] Code générique?
   - [ ] Complexité <15?
   - [ ] Noms explicites?
   - [ ] Tests présents?
   - [ ] GoDoc complet?
   - [ ] Duplication avec autres builders?
   - [ ] Anti-patterns?

4. **Analyser duplication**
   ```bash
   # Chercher patterns communs
   dupl -threshold 15 rete/builder*.go
   
   # Ou manuellement comparer
   diff -u rete/builder_alpha_rules.go rete/builder_join_rules.go | less
   ```

### Phase 2: Identification des problèmes (45-60 min)

**Créer liste priorisée dans** `REPORTS/review-rete/06_builders_issues.md`:

```markdown
# Problèmes Identifiés - Builders

## P0 - BLOQUANT

### 1. [Si bugs détectés]
- **Fichier:** builder_X.go:XXX
- **Type:** Bug construction / Règle invalide acceptée
- **Impact:** Réseau incorrect
- **Solution:** ...

## P1 - IMPORTANT

### 1. Complexité >15 dans builder_orchestration.go
- **Fichier:** builder_orchestration.go:XXX
- **Fonction:** `buildCompleteNetwork`
- **Complexité:** 22
- **Impact:** Maintenance difficile
- **Solution:** Extract Function - décomposer

### 2. Duplication entre AlphaBuilder et JoinBuilder
- **Fichiers:** builder_alpha_rules.go:XXX, builder_join_rules.go:YYY
- **Type:** Code dupliqué (validation, error handling)
- **Impact:** Maintenance double
- **Solution:** Extract Function partagée

### 3. Hardcoding limites
- **Fichiers:** Multiples
- **Type:** Magic numbers (max patterns, profondeur, etc.)
- **Impact:** Pas configurable
- **Solution:** Extract Constant

## P2 - SOUHAITABLE
...
```

**Problèmes à chercher:**

**P0:**
- Bugs de construction
- Validation incorrecte (règles invalides acceptées)
- Race conditions (registry)
- Panic dans builders

**P1:**
- Complexité >15
- Duplication entre builders (>15 lignes similaires)
- Hardcoding limites/seuils
- Exports non justifiés
- Couverture <70%
- Missing copyright

**P2:**
- Complexité 10-15
- Optimisations mineures
- Refactoring clarté

### Phase 3: Corrections (90-120 min)

#### 3.1 Fixer P0 (bloquants)

**Exemple: Validation incorrecte**

```go
// AVANT - règle invalide acceptée
func validateRule(rule *Rule) error {
    if rule.Name == "" {
        return errors.New("empty name")
    }
    // ❌ Ne valide pas patterns!
    return nil
}

// APRÈS - validation complète
func validateRule(rule *Rule) error {
    if rule == nil {
        return errors.New("nil rule")
    }
    if rule.Name == "" {
        return errors.New("rule must have a name")
    }
    if len(rule.Patterns) == 0 {
        return fmt.Errorf("rule %s has no patterns", rule.Name)
    }
    if len(rule.Patterns) > MaxPatternsPerRule {
        return fmt.Errorf("rule %s has too many patterns: %d > %d", 
            rule.Name, len(rule.Patterns), MaxPatternsPerRule)
    }
    for i, pattern := range rule.Patterns {
        if err := validatePattern(pattern); err != nil {
            return fmt.Errorf("rule %s pattern %d: %w", rule.Name, i, err)
        }
    }
    return nil
}
```

**Commit:**
```bash
git commit -m "[Review-06/Builders] fix(P0): validation complète des règles

- Valide règle non nulle
- Valide au moins un pattern
- Valide nombre max patterns
- Valide chaque pattern
- Messages d'erreur avec contexte
- Tests edge cases ajoutés

Resolves: P0-builders-validation-incomplete
Refs: scripts/review-rete/06_builders.md"
```

#### 3.2 Décomposer fonctions complexes (P1)

**Identifier fonctions >15:**
```bash
gocyclo -over 15 rete/builder*.go
```

**Pattern de décomposition:**

```go
// AVANT - complexité 22
func buildCompleteNetwork(rules []*Rule) (*Network, error) {
    // 40 lignes validation
    // 50 lignes construction alpha
    // 60 lignes construction beta
    // 30 lignes connexions
    // 20 lignes finalisation
}

// APRÈS - décomposer

func buildCompleteNetwork(rules []*Rule) (*Network, error) {
    // Orchestration - complexité ~8
    if err := validateAllRules(rules); err != nil {
        return nil, err
    }
    
    network := newNetwork()
    
    if err := buildAlphaNetwork(network, rules); err != nil {
        return nil, fmt.Errorf("build alpha: %w", err)
    }
    
    if err := buildBetaNetwork(network, rules); err != nil {
        return nil, fmt.Errorf("build beta: %w", err)
    }
    
    if err := connectComponents(network); err != nil {
        return nil, fmt.Errorf("connect: %w", err)
    }
    
    finalizeNetwork(network)
    
    return network, nil
}

// Sous-fonctions <12 complexité
func validateAllRules(rules []*Rule) error { /* 9 */ }
func buildAlphaNetwork(network *Network, rules []*Rule) error { /* 11 */ }
func buildBetaNetwork(network *Network, rules []*Rule) error { /* 12 */ }
func connectComponents(network *Network) error { /* 7 */ }
func finalizeNetwork(network *Network) { /* 5 */ }
```

**Commit:**
```bash
git commit -m "[Review-06/Builders] refactor(P1): décompose buildCompleteNetwork (22→8)

- Extrait validateAllRules() (9)
- Extrait buildAlphaNetwork() (11)
- Extrait buildBetaNetwork() (12)
- Extrait connectComponents() (7)
- Extrait finalizeNetwork() (5)
- Orchestration: 8
- Tests unitaires ajoutés

Resolves: P1-builders-complexity-22
Refs: scripts/review-rete/06_builders.md"
```

#### 3.3 Éliminer duplication (P1)

**Identifier duplication:**
```bash
dupl -threshold 15 rete/builder*.go
```

**Extraire code commun:**

```go
// AVANT - duplication entre builders

// builder_alpha_rules.go
func (b *AlphaBuilder) Build(pattern Pattern) error {
    if pattern.Type == "" {
        return errors.New("empty type")
    }
    if len(pattern.Conditions) == 0 {
        return errors.New("no conditions")
    }
    // ... construction alpha
}

// builder_join_rules.go
func (b *JoinBuilder) Build(pattern Pattern) error {
    if pattern.Type == "" {  // ❌ Duplication
        return errors.New("empty type")
    }
    if len(pattern.Conditions) == 0 {  // ❌ Duplication
        return errors.New("no conditions")
    }
    // ... construction join
}

// APRÈS - extraction

// builder_utils.go (nouveau ou existant)
func validatePattern(pattern Pattern) error {
    if pattern.Type == "" {
        return errors.New("pattern must have a type")
    }
    if len(pattern.Conditions) == 0 {
        return errors.New("pattern must have at least one condition")
    }
    return nil
}

// builder_alpha_rules.go
func (b *AlphaBuilder) Build(pattern Pattern) error {
    if err := validatePattern(pattern); err != nil {
        return fmt.Errorf("alpha build: %w", err)
    }
    // ... construction alpha
}

// builder_join_rules.go
func (b *JoinBuilder) Build(pattern Pattern) error {
    if err := validatePattern(pattern); err != nil {
        return fmt.Errorf("join build: %w", err)
    }
    // ... construction join
}
```

**Commit:**
```bash
git commit -m "[Review-06/Builders] refactor(P1): élimine duplication validation pattern

- Extrait validatePattern() dans builder_utils.go
- Utilisé par AlphaBuilder et JoinBuilder
- Réduit duplication de 30 lignes
- Messages d'erreur uniformes
- Tests unitaires pour validatePattern

Resolves: P1-builders-duplication-validation
Refs: scripts/review-rete/06_builders.md"
```

#### 3.4 Éliminer hardcoding (P1)

```go
// AVANT
if len(patterns) > 50 { return errTooMany }
if depth > 10 { return errTooDeep }
maxRetries := 3

// APRÈS
const (
    MaxPatternsPerRule = 50
    MaxBuildDepth      = 10
    BuildMaxRetries    = 3
)

if len(patterns) > MaxPatternsPerRule { 
    return fmt.Errorf("too many patterns: %d > %d", len(patterns), MaxPatternsPerRule)
}
if depth > MaxBuildDepth { return errTooDeep }
maxRetries := BuildMaxRetries
```

### Phase 4: Validation finale (15-30 min)

```bash
#!/bin/bash
echo "=== VALIDATION FINALE BUILDERS ==="

# 1. Tests
echo "🧪 Tests..."
go test -v ./rete -run "TestBuilder"
TESTS=$?

# 2. Race detector
echo "🏁 Race detector..."
go test -race ./rete -run "TestBuilder|TestRegistry"
RACE=$?

# 3. Complexité
echo "📊 Complexité..."
COMPLEX=$(gocyclo -over 15 rete/builder*.go | wc -l)

# 4. Couverture
echo "📈 Couverture..."
go test -coverprofile=builders_final.out ./rete -run "TestBuilder" 2>/dev/null
COVERAGE=$(go tool cover -func=builders_final.out | tail -1 | awk '{print $3}' | sed 's/%//')

# 5. Copyright
echo "©️  Copyright..."
MISSING_COPYRIGHT=0
for file in rete/builder*.go; do
    if ! head -1 "$file" | grep -q "Copyright"; then
        MISSING_COPYRIGHT=$((MISSING_COPYRIGHT + 1))
        echo "  ⚠️  $file"
    fi
done

# 6. Validation
echo "✅ Validation..."
make validate
VALIDATE=$?

# Résumé
echo ""
echo "=== RÉSULTATS ==="
[ $TESTS -eq 0 ] && echo "✅ Tests: PASS" || echo "❌ Tests: FAIL"
[ $RACE -eq 0 ] && echo "✅ Race: PASS" || echo "❌ Race: FAIL"
[ $COMPLEX -eq 0 ] && echo "✅ Complexité: OK" || echo "❌ Complexité: $COMPLEX >15"
[ $(echo "$COVERAGE >= 85" | bc -l) -eq 1 ] && echo "✅ Couverture: $COVERAGE%" || echo "❌ Couverture: $COVERAGE%"
[ $MISSING_COPYRIGHT -eq 0 ] && echo "✅ Copyright: OK" || echo "❌ Copyright: $MISSING_COPYRIGHT manquants"
[ $VALIDATE -eq 0 ] && echo "✅ Validation: PASS" || echo "❌ Validation: FAIL"

# Verdict
if [ $TESTS -eq 0 ] && [ $RACE -eq 0 ] && [ $COMPLEX -eq 0 ] && [ $(echo "$COVERAGE >= 85" | bc -l) -eq 1 ] && [ $MISSING_COPYRIGHT -eq 0 ] && [ $VALIDATE -eq 0 ]; then
    echo ""
    echo "🎉 VALIDATION RÉUSSIE - Prêt pour Prompt 07!"
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

**Créer:** `REPORTS/review-rete/06_builders_report.md`

**Structure obligatoire:**

```markdown
# 🔍 Revue de Code : Builders et Construction

**Date:** 2024-12-XX  
**Réviseur:** [Nom]  
**Durée:** Xh Ym

---

## 📊 Vue d'Ensemble

- **Fichiers analysés:** 12
- **Lignes de code:** ~3,000
- **Complexité max avant:** XX
- **Complexité max après:** <15
- **Couverture avant:** X%
- **Couverture après:** Y%

---

## ✅ Points Forts

- Séparation claire des builders (Alpha/Join/Exists/Accumulator)
- Orchestration identifiée
- ...

---

## ❌ Problèmes Identifiés et Corrigés

### P0 - BLOQUANT

#### 1. [Si applicable]
...

### P1 - IMPORTANT

#### 1. Complexité XX dans builder_orchestration.go
- **Fonction:** buildCompleteNetwork
- **Avant:** 22
- **Après:** 8
- **Décomposition:** 5 sous-fonctions
- **Commit:** abc1234

#### 2. Duplication validation entre builders
- **Lignes dupliquées:** ~30
- **Solution:** Extrait validatePattern()
- **Commit:** def5678

#### 3. Hardcoding limites
- **Constantes créées:** 8
- **Commit:** ghi9012

---

## 🔧 Changements Apportés

### Refactoring

1. **Décomposition orchestration**
   - 1 fonction 200 lignes → 6 fonctions <40 lignes
   - Complexité 22 → max 12
   - Tests: 5 nouveaux tests unitaires

2. **Élimination duplication**
   - 3 blocs dupliqués extraits
   - Helpers partagés créés
   - 90 lignes dupliquées → 30 lignes partagées

3. **Constantes nommées**
   - 15 magic numbers → constantes
   - 7 magic strings → constantes

---

## 📈 Métriques Avant/Après

| Métrique | Avant | Après | Amélioration |
|----------|-------|-------|--------------|
| Complexité max | 22 | 12 | ✅ -45% |
| Fonctions >15 | 5 | 0 | ✅ 100% |
| Couverture | 72% | 88% | ✅ +16% |
| Duplication | ~90 lignes | 0 | ✅ 100% |
| Magic numbers | 15 | 0 | ✅ 100% |

---

## 💡 Recommandations Futures

### Court terme
1. Benchmarks construction avec règles réelles
2. Profiling mémoire sur grandes règles
3. Documentation architecture builders

### Moyen terme
1. Builder pour nouveaux types de patterns
2. Validation sémantique plus poussée
3. Optimisation ordre de construction

---

## 🏁 Verdict

✅ **APPROUVÉ**

Séparation responsabilités claire, duplication éliminée, standards respectés.
Prêt pour Prompt 07 (Actions).

---

**Prochaines étapes:**
1. Merge commits
2. Lancer Prompt 07
3. Documenter patterns builders
```

### 2. Tests ajoutés

```bash
git diff --name-only | grep "_test.go" > REPORTS/builders_tests_added.txt
```

### 3. Commits atomiques

**Format strict:**
```
[Review-06/Builders] <type>(scope): <description courte>

- Détail 1
- Détail 2
- Resolves: <issue>

Refs: scripts/review-rete/06_builders.md
```

---

## 📊 Métriques de succès

| Métrique | Valeur Actuelle | Cible | Critique |
|----------|----------------|-------|----------|
| Complexité max | À mesurer | <15 | ⚠️ Oui |
| Fonctions >15 | À mesurer | 0 | ⚠️ Oui |
| Couverture tests | À mesurer | >85% | Oui |
| Duplication | À mesurer | 0 | Oui |
| Exports publics | À mesurer | Minimal | Oui |
| Magic numbers | À mesurer | 0 | Oui |
| Copyright | À mesurer | 100% | Oui |
| Race detector | À mesurer | Clean | Oui |

---

## 🎓 Ressources et références

### Standards Projet
- [common.md](../../.github/prompts/common.md)
- [review.md](../../.github/prompts/review.md)
- [Makefile](../../Makefile)

### Design Patterns
- Builder Pattern
- Factory Pattern
- Strategy Pattern (pour builders)
- Registry Pattern

### Refactoring
- [Extract Function](https://refactoring.guru/extract-method)
- [Extract Class](https://refactoring.guru/extract-class)

---

## ✅ Checklist finale avant Prompt 07

**Validation technique:**
- [ ] Tous tests builders passent
- [ ] Race detector clean
- [ ] Aucune fonction >15
- [ ] Couverture >85%
- [ ] `make validate` passe

**Qualité code:**
- [ ] Aucun hardcoding
- [ ] Code générique
- [ ] Exports minimaux
- [ ] Constantes nommées
- [ ] Noms explicites
- [ ] Fonctions <50 lignes
- [ ] Pas de duplication
- [ ] SRP respecté

**Tests:**
- [ ] Tests unitaires par builder
- [ ] Tests intégration orchestration
- [ ] Tests déterministes

**Documentation:**
- [ ] Copyright 100%
- [ ] GoDoc complet
- [ ] Rapport créé

---

## 🚀 Script d'analyse rapide

```bash
#!/bin/bash
# scripts/review-rete/analyze_builders.sh

set -e
echo "=== ANALYSE BUILDERS ==="
echo ""

mkdir -p REPORTS/review-rete

# Baseline
echo "📊 Mesure baseline..."
gocyclo -over 10 rete/builder*.go > REPORTS/review-rete/builders_complexity_before.txt
go test -coverprofile=REPORTS/review-rete/builders_coverage_before.out ./rete -run "TestBuilder" 2>/dev/null
go tool cover -func=REPORTS/review-rete/builders_coverage_before.out > REPORTS/review-rete/builders_coverage_before.txt

echo "✅ Baseline sauvegardée"
echo ""

# Complexité
echo "📈 TOP COMPLEXITÉ (>15)..."
gocyclo -top 30 rete/builder*.go | head -20
echo ""

# Duplication
echo "🔍 DUPLICATION..."
dupl -threshold 15 rete/builder*.go || echo "  (dupl non installé ou aucune duplication)"
echo ""

# Copyright
echo "©️  COPYRIGHT..."
MISSING=0
for file in rete/builder*.go; do
    if ! head -1 "$file" | grep -q "Copyright"; then
        echo "  ❌ $file"
        MISSING=$((MISSING + 1))
    fi
done
[ $MISSING -eq 0 ] && echo "  ✓ OK"

echo ""
echo "=== Analyse terminée ==="
echo "Créer REPORTS/review-rete/06_builders_issues.md"
```

**Lancer:**
```bash
chmod +x scripts/review-rete/analyze_builders.sh
./scripts/review-rete/analyze_builders.sh
```

---

**Prêt à commencer?** 🚀

Bonne revue! Respecter scrupuleusement les standards common.md et review.md.