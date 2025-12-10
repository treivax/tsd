# 🔧 Maintenance - Prompt Universel

> **📋 Standards** : Ce prompt respecte les règles de [common.md](./common.md)

## 🎯 Objectif

Maintenir le projet TSD : migration, nettoyage, vérification licence, statistiques, ou optimisation.

---

## ⚠️ Rappels Critiques

Avant de commencer, consulter [common.md](./common.md) :
- [🔒 Licence et Copyright](./common.md#licence-et-copyright) - Vérifications obligatoires
- [🏗️ Architecture](./common.md#architecture-et-organisation) - Principes projet
- [🔧 Outils](./common.md#outils-et-commandes) - Commandes validation

---

## 📋 Instructions

### 1. Définir l'Action

**Précise** :
- **Type** : [ ] Migration  [ ] Nettoyage  [ ] Licence  [ ] Stats  [ ] Optimisation
- **Portée** : Projet entier, module(s) spécifique(s)
- **Objectif** : Améliorer quoi exactement ?

### 2. Migration

#### Cas d'Usage
- Migration version Go
- Migration dépendances
- Migration API (breaking changes)
- Migration structure projet

#### Process Migration

1. **Planifier**
   - Documenter changements nécessaires
   - Identifier impacts
   - Prévoir rollback si échec

2. **Préparer**
   ```bash
   # Sauvegarder état actuel
   git checkout -b migration-backup
   
   # Créer branche migration
   git checkout -b migrate-to-X
   ```

3. **Migrer progressivement**
   - Un module à la fois
   - Tests passent après chaque étape
   - Commits atomiques

4. **Valider**
   ```bash
   # Tests complets
   go test ./...
   
   # Validation complète
   make validate
   
   # Vérifier dépendances
   go mod tidy
   go mod verify
   ```

5. **Documenter**
   - CHANGELOG.md mis à jour
   - README si changements API
   - Guide migration si breaking changes

#### Migration Version Go

```bash
# Mettre à jour go.mod
go mod edit -go=1.21

# Mettre à jour dépendances
go get -u ./...
go mod tidy

# Tester
go test ./...
make validate
```

### 3. Nettoyage (Deep Clean)

#### Checklist Nettoyage

**Code Mort** :
```bash
# Trouver code non utilisé
go run golang.org/x/tools/cmd/deadcode@latest ./...

# Supprimer imports non utilisés
goimports -w .
```

**Fichiers Temporaires** :
```bash
# Nettoyer builds
go clean -cache -testcache -modcache

# Supprimer fichiers générés
rm -f *.prof *.out *.test
```

**Documentation Obsolète** :
- [ ] README à jour avec code actuel
- [ ] Docs obsolètes supprimées
- [ ] Liens cassés corrigés
- [ ] Exemples fonctionnels

**Tests Obsolètes** :
- [ ] Tests pour code supprimé → supprimer
- [ ] Tests commentés → supprimer ou corriger
- [ ] Fixtures non utilisées → supprimer

**Dépendances** :
```bash
# Nettoyer dépendances non utilisées
go mod tidy

# Vérifier vulnérabilités
go list -m all | nancy sleuth
```

**Refactoring** :
- [ ] Duplication éliminée (DRY)
- [ ] Fonctions trop longues décomposées
- [ ] Complexité réduite
- [ ] Nommage amélioré

### 4. Vérification Licence

#### En-têtes Copyright

```bash
# Vérifier tous les fichiers .go
for file in $(find . -name "*.go" -type f ! -path "./.git/*" ! -path "*/vendor/*"); do
    if ! head -1 "$file" | grep -q "Copyright\|Code generated"; then
        echo "⚠️  EN-TÊTE MANQUANT: $file"
    fi
done
```

#### Ajouter en-tête manquant

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package monpackage
```

#### Vérifier Dépendances

```bash
# Lister toutes les dépendances
go list -m all

# Vérifier licences dépendances
go-licenses report ./... --template=licenses.tpl
```

**Licences acceptées** (voir [common.md](./common.md#licence-et-copyright)) :
- ✅ MIT, BSD, Apache-2.0, ISC
- ⚠️ Éviter GPL, AGPL, LGPL
- ❌ Code sans licence, propriétaire

#### Documentation Licence

Si dépendance tierce utilisée :
1. Ajouter à `go.mod`
2. Documenter dans `THIRD_PARTY_LICENSES.md`
3. Vérifier compatibilité MIT

### 5. Statistiques Code

#### Métriques de Base

```bash
# Lignes de code
find . -name "*.go" -not -path "*/vendor/*" | xargs wc -l | tail -1

# Nombre de fichiers
find . -name "*.go" -not -path "*/vendor/*" | wc -l

# Nombre de packages
go list ./... | wc -l
```

#### Complexité

```bash
# Complexité cyclomatique
gocyclo -over 15 .

# Top 10 fonctions complexes
gocyclo -top 10 .
```

#### Couverture Tests

```bash
# Couverture globale
go test -cover ./...

# Rapport détaillé
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Visualisation HTML
go tool cover -html=coverage.out
```

#### Dépendances

```bash
# Graphe dépendances
go mod graph

# Pourquoi dépendance X ?
go mod why github.com/some/package

# Packages obsolètes
go list -u -m all
```

#### Rapport Complet

```bash
# Générer rapport stats
cat > REPORTS/stats-$(date +%Y%m%d).md << EOF
# Statistiques Code - $(date +%Y-%m-%d)

## Métriques Globales
- Lignes de code: $(find . -name "*.go" -not -path "*/vendor/*" | xargs wc -l | tail -1 | awk '{print $1}')
- Fichiers Go: $(find . -name "*.go" -not -path "*/vendor/*" | wc -l)
- Packages: $(go list ./... | wc -l)

## Complexité
$(gocyclo -top 10 .)

## Couverture Tests
$(go test -cover ./... 2>&1)
EOF
```

### 6. Optimisation Performance

#### Avant d'Optimiser

⚠️ **Règle d'Or** : Ne jamais optimiser sans mesure !

1. **Identifier le problème**
   - Profiling montre goulot d'étranglement réel
   - Impact utilisateur significatif
   - Baseline mesurée

2. **Mesurer avant**
   ```bash
   # Benchmark avant optimisation
   go test -bench=. -benchmem > bench-before.txt
   ```

#### Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
# Commandes dans pprof: top, list, web

# Memory profiling
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# Allocation tracking
go test -benchmem -bench=.

# Trace complet
go test -trace=trace.out
go tool trace trace.out
```

#### Optimisations Courantes

**Allocations mémoire** :
```go
// ❌ Allocations répétées
func process(items []Item) {
    for _, item := range items {
        result := make([]byte, 1024)  // Allocation à chaque itération
        // ...
    }
}

// ✅ Réutilisation
func process(items []Item) {
    result := make([]byte, 1024)  // Une seule allocation
    for _, item := range items {
        result = result[:0]  // Réinitialiser sans allouer
        // ...
    }
}
```

**String concatenation** :
```go
// ❌ Lent avec beaucoup de concaténations
var result string
for _, s := range strings {
    result += s  // Allocation à chaque fois
}

// ✅ Utiliser strings.Builder
var builder strings.Builder
for _, s := range strings {
    builder.WriteString(s)
}
result := builder.String()
```

**Slices pré-allocation** :
```go
// ❌ Réallocations multiples
var results []Result
for _, item := range items {
    results = append(results, process(item))
}

// ✅ Pré-allouer avec capacité connue
results := make([]Result, 0, len(items))
for _, item := range items {
    results = append(results, process(item))
}
```

#### Valider Optimisation

```bash
# Benchmark après optimisation
go test -bench=. -benchmem > bench-after.txt

# Comparer
benchcmp bench-before.txt bench-after.txt

# Vérifier non-régression
go test ./...
make validate
```

---

## ✅ Checklist Maintenance

**Migration** :
- [ ] Changements planifiés et documentés
- [ ] Migration incrémentale (étapes atomiques)
- [ ] Tests passent après chaque étape
- [ ] `go.mod` mis à jour
- [ ] CHANGELOG.md à jour

**Nettoyage** :
- [ ] Code mort supprimé
- [ ] Documentation obsolète supprimée
- [ ] Tests inutiles supprimés
- [ ] Dépendances nettoyées (`go mod tidy`)
- [ ] Duplication éliminée

**Licence** :
- [ ] En-têtes copyright présents (tous les .go)
- [ ] Dépendances vérifiées (licences compatibles)
- [ ] `THIRD_PARTY_LICENSES.md` à jour
- [ ] Pas de code GPL/AGPL/propriétaire

**Stats** :
- [ ] Métriques collectées
- [ ] Complexité < 15
- [ ] Couverture > 80%
- [ ] Rapport généré dans REPORTS/

**Optimisation** :
- [ ] Profiling effectué (goulot identifié)
- [ ] Benchmark avant mesuré
- [ ] Optimisation appliquée
- [ ] Benchmark après mesuré
- [ ] Amélioration > 20% (sinon pas worth it)
- [ ] Tests passent (comportement identique)

---

## 🎯 Principes

1. **Mesurer** : Données objectives, pas intuitions
2. **Incrémental** : Petits changements validés
3. **Non-régression** : Tests passent toujours
4. **Documentation** : Changements documentés
5. **Prudence** : Backup avant changements majeurs

---

## 🚫 Anti-Patterns

- ❌ Migration big bang (tout d'un coup)
- ❌ Optimisation prématurée sans mesure
- ❌ Supprimer code sans vérifier utilisation
- ❌ Négliger documentation lors migration
- ❌ Ignorer licences dépendances
- ❌ Optimiser ce qui n'est pas le goulot
- ❌ Pas de backup avant changements majeurs

---

## 📊 Commandes Utiles

```bash
# Maintenance générale
go clean -cache -testcache -modcache    # Nettoyer caches
go mod tidy                             # Nettoyer dépendances
go mod verify                           # Vérifier intégrité
goimports -w .                          # Nettoyer imports
go fmt ./...                            # Formater code

# Analyse
gocyclo -over 15 .                      # Complexité
go test -cover ./...                    # Couverture
go list -u -m all                       # Dépendances obsolètes
staticcheck ./...                       # Analyse statique

# Profiling
go test -cpuprofile=cpu.prof -bench=.   # CPU
go test -memprofile=mem.prof -bench=.   # Mémoire
go tool pprof cpu.prof                  # Analyser

# Validation
make validate                           # Validation complète
```

---

## 📚 Ressources

- [common.md](./common.md) - Standards projet
- [Go Modules](https://go.dev/ref/mod) - Gestion dépendances
- [pprof](https://github.com/google/pprof) - Profiling
- [Optimization Guide](https://go.dev/doc/diagnostics) - Guide Go
- [Makefile](../../Makefile) - Commandes projet

---

**Workflow** : Planifier → Mesurer → Exécuter → Valider → Documenter → Rapport