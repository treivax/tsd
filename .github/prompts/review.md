# 🔍 Revue et Qualité - Prompt Universel

> **📋 Standards** : Ce prompt respecte les règles de [common.md](./common.md)

## 🎯 Objectif

Analyser et améliorer la qualité du code : revue de code, refactoring, ou optimisation.

---

## ⚠️ Rappels Critiques

Avant de commencer, consulter [common.md](./common.md) :
- [⚠️ Standards Code Go](./common.md#standards-de-code-go) - Conventions, complexité, qualité
- [🎨 Conventions Nommage](./common.md#conventions-de-nommage) - Standards projet
- [📋 Checklist Commit](./common.md#checklist-avant-commit) - Validation

---

## 📋 Instructions

### 1. Définir l'Action

**Précise** :
- **Type** : [ ] Revue code  [ ] Refactoring  [ ] Optimisation  [ ] Audit qualité
- **Portée** : Fichier(s), module(s), fonction(s) concerné(s)
- **Objectif** : Améliorer quoi ? (lisibilité, performance, maintenabilité)
- **Contraintes** : Ne pas changer le comportement (sauf si optimisation)

### 2. Revue de Code

#### Points de Vérification

**Architecture et Design** :
- [ ] Respect principes SOLID
- [ ] Séparation des responsabilités claire
- [ ] Pas de couplage fort
- [ ] Interfaces appropriées
- [ ] Composition over inheritance

**Qualité du Code** :
- [ ] Noms explicites (variables, fonctions, types)
- [ ] Fonctions < 50 lignes (sauf justification)
- [ ] Complexité cyclomatique < 15
- [ ] Pas de duplication (DRY)
- [ ] Code auto-documenté

**Conventions Go** :
- [ ] `go fmt` appliqué
- [ ] `goimports` utilisé
- [ ] Conventions nommage respectées (voir [common.md](./common.md#conventions-de-nommage))
- [ ] Erreurs gérées explicitement
- [ ] Pas de panic (sauf cas critique)

**Encapsulation** :
- [ ] Variables/fonctions privées par défaut
- [ ] Exports publics minimaux et justifiés
- [ ] Contrats d'interface respectés
- [ ] Pas d'exposition interne inutile

**Standards Projet** :
- [ ] En-tête copyright présent
- [ ] Aucun hardcoding (valeurs, chemins, configs)
- [ ] Code générique avec paramètres
- [ ] Constantes nommées pour valeurs

**Tests** :
- [ ] Tests présents (couverture > 80%)
- [ ] Tests déterministes
- [ ] Tests isolés
- [ ] Messages d'erreur clairs

**Documentation** :
- [ ] GoDoc pour exports
- [ ] Commentaires inline si complexe
- [ ] Exemples d'utilisation
- [ ] README module à jour

**Performance** :
- [ ] Complexité algorithmique acceptable
- [ ] Pas de boucles inutiles
- [ ] Pas de calculs redondants
- [ ] Ressources libérées proprement

**Sécurité** :
- [ ] Validation des entrées
- [ ] Gestion des erreurs robuste
- [ ] Pas d'injection possible
- [ ] Gestion cas nil/vides

### 3. Refactoring

#### Objectifs Refactoring

1. **Améliorer lisibilité** sans changer comportement
2. **Réduire complexité** en décomposant
3. **Éliminer duplication** (DRY)
4. **Améliorer maintenabilité** par meilleure structure

#### Process Refactoring

1. **Analyser l'existant**
   - Identifier les code smells
   - Repérer la duplication
   - Mesurer la complexité

2. **Planifier les étapes**
   - Refactoring incrémental
   - Chaque étape validée par tests
   - Commits atomiques

3. **Exécuter**
   - Une technique à la fois
   - Tests passent après chaque étape
   - Valider avec `make validate`

4. **Valider**
   - Comportement identique
   - Tests passent
   - Pas de régression performance

#### Techniques Refactoring

**Extract Function** :
```go
// Avant - fonction longue
func processOrder(order Order) error {
    // 50 lignes de validation
    // 30 lignes de traitement
    // 20 lignes de notification
}

// Après - décomposé
func processOrder(order Order) error {
    if err := validateOrder(order); err != nil {
        return err
    }
    if err := executeOrder(order); err != nil {
        return err
    }
    return notifyOrderProcessed(order)
}

func validateOrder(order Order) error { /* ... */ }
func executeOrder(order Order) error { /* ... */ }
func notifyOrderProcessed(order Order) error { /* ... */ }
```

**Extract Constant** :
```go
// Avant - magic numbers
func isValid(age int) bool {
    return age >= 18 && age <= 120
}

// Après - constantes nommées
const (
    MinAge = 18
    MaxAge = 120
)

func isValid(age int) bool {
    return age >= MinAge && age <= MaxAge
}
```

**Simplify Conditional** :
```go
// Avant - condition complexe
if status == "active" && user != nil && user.HasPermission("read") && !expired {
    // ...
}

// Après - fonction explicite
func canAccess(status string, user *User, expired bool) bool {
    return status == "active" &&
           user != nil &&
           user.HasPermission("read") &&
           !expired
}

if canAccess(status, user, expired) {
    // ...
}
```

**Rename** :
```go
// Avant - noms peu clairs
func proc(d []byte) []byte { /* ... */ }

// Après - noms explicites
func processData(data []byte) []byte { /* ... */ }
```

### 4. Optimisation

#### Quand Optimiser ?

- ✅ Mesure prouve un problème (profiling)
- ✅ Goulot d'étranglement identifié
- ✅ Impact utilisateur significatif
- ❌ Optimisation prématurée
- ❌ Micro-optimisations sans mesure

#### Process Optimisation

1. **Mesurer** - Benchmark avant
2. **Identifier** - Profiler pour trouver goulot
3. **Optimiser** - Une chose à la fois
4. **Mesurer** - Benchmark après
5. **Valider** - Tests passent, comportement identique

```bash
# Profiling
go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=.
go tool pprof cpu.prof

# Benchmarking
go test -bench=BenchmarkFunction -benchmem
```

---

## ✅ Checklist Revue

- [ ] Architecture respecte SOLID
- [ ] Code suit conventions Go
- [ ] Encapsulation respectée (privé par défaut)
- [ ] Aucun hardcoding
- [ ] Code générique et réutilisable
- [ ] Constantes nommées
- [ ] Noms explicites
- [ ] Complexité < 15
- [ ] Fonctions < 50 lignes
- [ ] Pas de duplication
- [ ] Tests présents (> 80%)
- [ ] GoDoc complet
- [ ] `go vet` + `staticcheck` OK
- [ ] Gestion erreurs robuste
- [ ] Performance acceptable

---

## 🎯 Principes

1. **Comportement préservé** - Refactoring ne change pas le comportement
2. **Incrémental** - Petites étapes validées par tests
3. **Mesurable** - Métriques avant/après
4. **Simple** - La solution la plus simple
5. **Testable** - Tests valident chaque étape

---

## 🚫 Anti-Patterns

- ❌ God Object (classe qui fait tout)
- ❌ Long Method (> 100 lignes)
- ❌ Long Parameter List (> 5 params)
- ❌ Duplicate Code
- ❌ Dead Code
- ❌ Magic Numbers/Strings
- ❌ Deep Nesting (> 4 niveaux)
- ❌ Shotgun Surgery (changement éparpillé)
- ❌ Feature Envy (méthode dans mauvaise classe)
- ❌ Primitive Obsession (types primitifs partout)

---

## 📊 Métriques Qualité

```bash
# Complexité cyclomatique
gocyclo -over 15 .

# Duplication
dupl -threshold 15 .

# Vérifications statiques
go vet ./...
staticcheck ./...
errcheck ./...
gosec ./...

# Linting
golangci-lint run

# Couverture tests
go test -cover ./...

# Validation complète
make validate
```

---

## 📝 Format Réponse Revue

```markdown
## 🔍 Revue de Code : [Module/Fichier]

### 📊 Vue d'Ensemble
- Lignes de code : X
- Complexité : Faible/Moyenne/Élevée
- Couverture tests : X%

### ✅ Points Forts
- Point fort 1
- Point fort 2

### ⚠️ Points d'Attention
- Point attention 1 (ligne X)
- Point attention 2 (ligne Y)

### ❌ Problèmes Identifiés
- Problème 1 (critique/majeur/mineur)
- Problème 2 (critique/majeur/mineur)

### 💡 Recommandations
1. Recommandation 1
2. Recommandation 2

### 📈 Métriques
- Avant : [métriques]
- Après : [métriques si refactoring]

### 🏁 Verdict
✅ Approuvé / ⚠️ Approuvé avec réserves / ❌ Changements requis
```

---

## 📚 Ressources

- [common.md](./common.md) - Standards projet
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review](https://github.com/golang/go/wiki/CodeReviewComments)
- [Refactoring Guru](https://refactoring.guru/)
- [Makefile](../../Makefile) - Commandes validation

---

**Workflow** : Analyser → Identifier → Planifier → Refactorer → Valider → Documenter