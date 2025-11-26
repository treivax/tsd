# 🔄 Migrer / Mettre à Niveau (Migrate)

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Tu veux migrer le projet vers une nouvelle version de Go, mettre à jour des dépendances externes, adapter le code suite à un changement d'API externe, ou effectuer toute autre forme de migration technique.

## Objectif

Migrer le projet vers une nouvelle version/dépendance tout en préservant la fonctionnalité, en minimisant les risques, et en documentant les changements nécessaires.

## Types de Migrations

### 1. **Migration Go Version**
- Mettre à jour vers nouvelle version de Go
- Adapter aux changements de langage
- Bénéficier des nouvelles features

### 2. **Migration Dépendances**
- Mettre à jour bibliothèques externes
- Adapter aux breaking changes
- Résoudre vulnérabilités

### 3. **Migration API**
- Adapter à nouveau format d'API externe
- Changer de provider/service
- Modifier protocole de communication

### 4. **Migration Architecture**
- Refonte structure de projet
- Changement de patterns
- Réorganisation modules

### 5. **Migration Données**
- Nouveau format de fichiers .constraint
- Nouveau schéma .facts
- Migration de configurations

## Instructions

### PHASE 1 : PRÉPARATION (Avant Migration)

#### 1.1 Identifier la Migration

**Définir clairement** :

```markdown
## Détails de la Migration

**Type** : Migration Go 1.19 → 1.21

**Motivation** :
- Bénéficier des nouvelles features (generics améliorés)
- Support officiel jusqu'en 2025
- Corrections de sécurité

**Scope** :
- Fichier `go.mod` (version Go)
- Code utilisant features dépréciées
- Tests affectés
- CI/CD pipelines

**Risques** :
- Breaking changes potentiels
- Incompatibilités dépendances
- Tests qui échouent
- Performance différente

**Timeline** :
- Durée estimée : 2-3 jours
- Sprint actuel ou prochain
```

#### 1.2 Analyser l'État Actuel

**Baseline avant migration** :

```bash
# Version actuelle
go version
cat go.mod | grep "^go "

# Dépendances actuelles
go list -m all

# État des tests
make test
make rete-unified

# Métriques de performance
go test -bench=. -benchmem ./rete > baseline_bench.txt

# Analyse statique
go vet ./...
golangci-lint run ./... > baseline_lint.txt
```

**Documenter l'état actuel** :
```markdown
## État Avant Migration

**Version Go** : 1.19.5
**Tests** : ✅ 234/234 passent
**Runner universel** : ✅ 58/58 passent
**Warnings** : 3 (documentés)
**Performance** : Baseline sauvegardée

**Dépendances** :
- github.com/pigeon v1.0.0
- github.com/stretchr/testify v1.8.0
- [... liste complète]
```

#### 1.3 Rechercher les Breaking Changes

**Documentation officielle** :

```bash
# Go release notes
# https://go.dev/doc/go1.21

# Chercher breaking changes
grep -i "breaking\|deprecated\|removed" RELEASE_NOTES.md
```

**Breaking changes identifiés** :
```markdown
## Breaking Changes Go 1.19 → 1.21

### Langage
- ✅ Generics : Améliorations (pas de breaking change)
- ⚠️  `any` devient alias officiel de `interface{}`
- ⚠️  Certaines fonctions `unsafe` modifiées

### Bibliothèque Standard
- ⚠️  `os.Readdir` déprécié → utiliser `os.ReadDir`
- ⚠️  `ioutil` package déprécié → utiliser `os`, `io`
- ✅ Pas d'impact sur notre code

### Tooling
- ✅ `go test` : nouvelles options (pas de breaking)
- ✅ `go build` : optimisations (pas de breaking)

### Notre Code
- 📁 `rete/utils.go` : Utilise `ioutil.ReadFile` → À migrer
- 📁 `test/helpers.go` : Utilise `ioutil.WriteFile` → À migrer
```

#### 1.4 Créer un Plan de Migration

**Stratégie** :

```markdown
## Plan de Migration

### Phase 1 : Environnement (1h)
1. Installer Go 1.21
2. Mettre à jour CI/CD (.github/workflows)
3. Mettre à jour Dockerfile (si applicable)
4. Mettre à jour documentation

### Phase 2 : Code (4h)
1. Mettre à jour `go.mod` (version Go)
2. Exécuter `go mod tidy`
3. Remplacer code déprécié :
   - `ioutil.ReadFile` → `os.ReadFile`
   - `ioutil.WriteFile` → `os.WriteFile`
4. Adapter code si breaking changes
5. Mettre à jour imports

### Phase 3 : Dépendances (2h)
1. Mettre à jour dépendances compatibles
2. Tester compatibilité
3. Résoudre conflits

### Phase 4 : Tests (2h)
1. Exécuter tests unitaires
2. Exécuter tests intégration
3. Exécuter runner universel
4. Comparer benchmarks
5. Corriger régressions

### Phase 5 : Validation (1h)
1. Analyse statique (go vet, golangci-lint)
2. Tests de régression complets
3. Validation performance
4. Review code

### Phase 6 : Documentation (1h)
1. Mettre à jour README.md
2. Mettre à jour CHANGELOG.md
3. Documenter changements
4. Communiquer à l'équipe

**Total estimé** : 11 heures (2 jours)
```

### PHASE 2 : ENVIRONNEMENT (Setup)

#### 2.1 Installer Nouvelle Version

**Installation Go** :

```bash
# Télécharger Go 1.21
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz

# Installer
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Vérifier
go version
# go version go1.21.0 linux/amd64 ✅
```

**Sauvegarder ancienne version** :
```bash
# Au cas où besoin de rollback
which go
# /usr/local/go/bin/go

# Faire backup
sudo cp -r /usr/local/go /usr/local/go-1.19.5-backup
```

#### 2.2 Mettre à Jour CI/CD

**GitHub Actions** :

```yaml
# .github/workflows/test.yml

name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'  # ← Mis à jour
      
      - name: Run tests
        run: make test
```

#### 2.3 Mettre à Jour Docker (si applicable)

```dockerfile
# Dockerfile

FROM golang:1.21-alpine  # ← Mis à jour

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o rete-runner ./cmd/rete-runner

CMD ["./rete-runner"]
```

### PHASE 3 : CODE (Adaptation)

#### 3.1 Mettre à Jour go.mod

```bash
# Mettre à jour version Go dans go.mod
# Avant :
# go 1.19

# Après :
go mod edit -go=1.21

# Nettoyer et réorganiser
go mod tidy

# Vérifier
cat go.mod
```

**Fichier `go.mod` après** :
```go
module github.com/user/tsd

go 1.21  // ← Mis à jour

require (
    github.com/pigeon v1.0.0
    github.com/stretchr/testify v1.8.4  // ← Peut être mis à jour
    // ...
)
```

#### 3.2 Remplacer Code Déprécié

**ioutil → os/io** :

```go
// AVANT (déprécié)
import "io/ioutil"

func readConstraint(path string) ([]byte, error) {
    return ioutil.ReadFile(path)
}

func writeOutput(path string, data []byte) error {
    return ioutil.WriteFile(path, data, 0644)
}
```

```go
// APRÈS (Go 1.21)
import "os"

func readConstraint(path string) ([]byte, error) {
    return os.ReadFile(path)
}

func writeOutput(path string, data []byte) error {
    return os.WriteFile(path, data, 0644)
}
```

**Script de migration automatique** :

```bash
# Remplacer ioutil.ReadFile
find . -name "*.go" -type f -exec sed -i 's/ioutil\.ReadFile/os.ReadFile/g' {} \;

# Remplacer ioutil.WriteFile
find . -name "*.go" -type f -exec sed -i 's/ioutil\.WriteFile/os.WriteFile/g' {} \;

# Mettre à jour imports
goimports -w .

# Vérifier que ça compile
go build ./...
```

#### 3.3 Adapter aux Nouvelles Features (Optionnel)

**Utiliser nouvelles features Go 1.21** :

```go
// AVANT : Utilisation de interface{}
func processValue(v interface{}) error {
    // ...
}

// APRÈS : Utilisation de any (alias officiel)
func processValue(v any) error {
    // ...
}
```

```go
// AVANT : Pas de min/max built-in
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// APRÈS : Utiliser built-in min/max (Go 1.21+)
result := min(a, b)  // Built-in function
```

### PHASE 4 : DÉPENDANCES (Mises à Jour)

#### 4.1 Lister Dépendances à Mettre à Jour

```bash
# Voir dépendances obsolètes
go list -u -m all

# Exemple output :
# github.com/stretchr/testify v1.8.0 [v1.8.4]
# github.com/pigeon v1.0.0 [v1.1.0]
```

#### 4.2 Mettre à Jour Progressivement

**Une par une pour isoler problèmes** :

```bash
# Testify (tests)
go get github.com/stretchr/testify@v1.8.4
go mod tidy
make test  # Vérifier que tests passent

# Pigeon (parser)
go get github.com/pigeon@v1.1.0
go mod tidy
make test  # Vérifier

# Toutes en une fois (si confiant)
go get -u ./...
go mod tidy
```

#### 4.3 Résoudre Conflits de Dépendances

```bash
# Si conflits
go mod why github.com/problematic/package

# Forcer version spécifique
go get github.com/package@v1.2.3

# Vérifier arbre de dépendances
go mod graph | grep package
```

### PHASE 5 : TESTS (Validation)

#### 5.1 Tests Unitaires

```bash
# Tous les tests
go test -v ./...

# Tests par package
go test -v ./rete
go test -v ./constraint
go test -v ./test

# Avec couverture
go test -cover ./...

# Détaillé si échecs
go test -v -failfast ./...
```

**Si tests échouent** :
```markdown
## Tests Échoués Après Migration

### Test: TestJoinNodePropagation
**Erreur** : `panic: interface conversion: interface {} is nil, not *Token`

**Cause** : Changement de comportement des type assertions en Go 1.21

**Fix** :
```go
// AVANT
token := value.(*Token)

// APRÈS
token, ok := value.(*Token)
if !ok {
    return fmt.Errorf("expected *Token, got %T", value)
}
```
```

#### 5.2 Tests d'Intégration

```bash
# Runner universel
make rete-unified

# Tests d'intégration
make test-integration

# Tests end-to-end
make test-e2e
```

#### 5.3 Benchmarks et Performance

```bash
# Benchmarks après migration
go test -bench=. -benchmem ./rete > after_bench.txt

# Comparer avec baseline
benchcmp baseline_bench.txt after_bench.txt

# Exemple output :
# benchmark                      old ns/op     new ns/op     delta
# BenchmarkAlphaNode-8           1000          950          -5.00%
# BenchmarkJoinNode-8            5000          5100         +2.00%
# 
# Acceptable si < ±10%
```

**Si dégradation significative** :
```markdown
## Dégradation Performance Détectée

**Benchmark** : BenchmarkJoinNode
**Avant** : 5000 ns/op
**Après** : 6500 ns/op
**Delta** : +30% ⚠️

**Investigation** :
- Profiling CPU avec `pprof`
- Identifier hotspots
- Vérifier si optimisations Go 1.21 applicables
- Considérer rollback si critique
```

#### 5.4 Tests de Régression

**Checklist** :
```markdown
## Validation Complète Post-Migration

### Tests Automatisés
- [ ] Tests unitaires : ✅ 234/234
- [ ] Tests intégration : ✅ 45/45
- [ ] Runner universel : ✅ 58/58
- [ ] Benchmarks : ✅ Performance acceptable (±5%)

### Tests Manuels
- [ ] Exécuter exemples docs/examples/
- [ ] Tester cas d'usage réels
- [ ] Vérifier output des logs
- [ ] Valider erreurs cohérentes

### Analyse Statique
- [ ] go vet : ✅ Aucun warning
- [ ] golangci-lint : ✅ Aucune erreur nouvelle
- [ ] go fmt : ✅ Code formaté
- [ ] goimports : ✅ Imports propres

### Documentation
- [ ] README.md : ✅ Version Go mise à jour
- [ ] CHANGELOG.md : ✅ Entrée ajoutée
- [ ] CONTRIBUTING.md : ✅ Prérequis mis à jour
- [ ] CI/CD : ✅ Pipelines fonctionnels
```

### PHASE 6 : DOCUMENTATION (Finalisation)

#### 6.1 Mettre à Jour Documentation

**README.md** :

```markdown
## Prérequis

- **Go 1.21+** (précédemment 1.19+)  ← Mis à jour
- Make
- Git

## Installation

```bash
# Vérifier version Go
go version  # Doit être >= 1.21

# Cloner et installer
git clone https://github.com/user/tsd.git
cd tsd
make install
```
```

**CHANGELOG.md** :

```markdown
## [Unreleased]

### Changed
- Migration vers Go 1.21 (depuis Go 1.19)
- Remplacement de `ioutil` par `os` et `io` (dépréciations Go 1.21)
- Mise à jour dépendances :
  - `testify` : v1.8.0 → v1.8.4
  - `pigeon` : v1.0.0 → v1.1.0

### Technical
- CI/CD : Mise à jour vers Go 1.21 dans GitHub Actions
- Dockerfile : Mise à jour image de base golang:1.21
- Performance : Légère amélioration (+5%) grâce aux optimisations Go 1.21

### Migration Guide
Pour les développeurs, voir [MIGRATION_GO_1.21.md](docs/MIGRATION_GO_1.21.md)
```

#### 6.2 Créer Guide de Migration (si applicable)

**docs/MIGRATION_GO_1.21.md** :

```markdown
# Guide de Migration : Go 1.19 → Go 1.21

## Pour les Développeurs

### Mettre à Jour Votre Environnement

1. **Installer Go 1.21** :
   ```bash
   # Linux/macOS
   wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
   
   # Vérifier
   go version
   ```

2. **Mettre à Jour Dépendances** :
   ```bash
   cd tsd/
   go mod download
   go mod tidy
   ```

3. **Vérifier** :
   ```bash
   make test
   make rete-unified
   ```

### Changements à Connaître

#### Code Déprécié Remplacé

- `ioutil.ReadFile` → `os.ReadFile`
- `ioutil.WriteFile` → `os.WriteFile`
- `interface{}` → `any` (recommandé)

#### Nouvelles Features Utilisables

- Built-in `min(a, b)` et `max(a, b)`
- Améliorations des generics
- `clear()` pour maps/slices

### Problèmes Connus

Aucun problème identifié avec notre codebase.

### Support

Si problèmes après migration :
1. Vérifier version Go : `go version`
2. Nettoyer cache : `go clean -modcache`
3. Re-télécharger : `go mod download`
4. Ouvrir issue si persiste
```

#### 6.3 Communiquer à l'Équipe

**Message d'annonce** :

```markdown
## 📢 Migration Go 1.21 Complétée

Bonjour l'équipe,

La migration vers **Go 1.21** est maintenant terminée et mergée dans `main`.

### Ce qui change pour vous

**Si vous développez sur TSD** :
1. Installer Go 1.21 : [Guide](docs/MIGRATION_GO_1.21.md)
2. Faire `go mod download` dans votre repo local
3. Re-runner vos tests

**Si vous utilisez TSD** :
- Aucun changement dans l'API
- Légère amélioration de performance (+5%)
- Toutes les fonctionnalités identiques

### Bénéfices

- ✅ Support officiel jusqu'en 2025
- ✅ Nouvelles features du langage disponibles
- ✅ Corrections de sécurité
- ✅ Optimisations compilateur (+5% performance)

### Questions / Problèmes

Voir [MIGRATION_GO_1.21.md](docs/MIGRATION_GO_1.21.md) ou me contacter.

Merci !
```

### PHASE 7 : ROLLBACK (Si Nécessaire)

#### 7.1 Plan de Rollback

**Préparer avant migration** :

```bash
# Créer branche de sauvegarde
git checkout -b backup/pre-go-1.21-migration
git push origin backup/pre-go-1.21-migration

# Tag de sauvegarde
git tag pre-migration-go-1.21
git push origin pre-migration-go-1.21

# Faire migration sur branche séparée
git checkout main
git checkout -b feature/migrate-go-1.21
```

#### 7.2 Procédure de Rollback

**Si problèmes critiques** :

```bash
# Option 1 : Revert du commit de migration
git revert <commit-hash-migration>
git push origin main

# Option 2 : Reset au tag pré-migration
git reset --hard pre-migration-go-1.21
git push origin main --force  # ⚠️ Coordonner avec équipe

# Option 3 : Restaurer depuis branche backup
git checkout main
git reset --hard backup/pre-go-1.21-migration
git push origin main --force  # ⚠️ Coordonner avec équipe

# Restaurer environnement local
sudo rm -rf /usr/local/go
sudo mv /usr/local/go-1.19.5-backup /usr/local/go
go version  # Vérifier
```

**Communiquer** :
```markdown
## ⚠️ Rollback Migration Go 1.21

Suite à [problème critique détecté], nous avons effectué un rollback
vers Go 1.19.

**Action requise** :
1. Pull la dernière version de `main`
2. Réinstaller Go 1.19 si nécessaire
3. `go mod download`

Le problème est en cours d'investigation. Nous retenterons la migration
une fois résolu.
```

## Critères de Succès

### ✅ Migration Complétée

- [ ] Nouvelle version installée et fonctionnelle
- [ ] `go.mod` mis à jour avec nouvelle version
- [ ] Code adapté (dépréciations remplacées)
- [ ] Dépendances mises à jour
- [ ] Tous les tests passent (unitaires + intégration + runner)

### ✅ Qualité Maintenue

- [ ] Aucune régression fonctionnelle
- [ ] Performance acceptable (±10% max)
- [ ] Analyse statique sans nouveaux warnings
- [ ] Code formaté et propre
- [ ] Pas de breaking changes pour utilisateurs

### ✅ Documentation

- [ ] README.md mis à jour (prérequis)
- [ ] CHANGELOG.md avec entrée migration
- [ ] Guide de migration créé (si nécessaire)
- [ ] CI/CD mis à jour
- [ ] Équipe informée

### ✅ Validation

- [ ] Tests en environnement local : ✅
- [ ] Tests en CI/CD : ✅
- [ ] Tests en staging/preprod : ✅
- [ ] Revue de code : ✅
- [ ] Approbation équipe : ✅

## Format de Réponse

```markdown
# 🔄 MIGRATION : Go 1.19 → Go 1.21

## 📋 Résumé

**Type** : Migration version Go
**De** : Go 1.19.5
**Vers** : Go 1.21.0
**Durée** : 10 heures (2 jours)
**Statut** : ✅ Complétée avec succès

## 🎯 Motivation

- Support officiel Go 1.19 se termine en 2024
- Bénéficier des nouvelles optimisations compilateur
- Accéder aux nouvelles features du langage
- Corrections de sécurité

## 📝 Changements Effectués

### Code
- ✅ `go.mod` : `go 1.19` → `go 1.21`
- ✅ Remplacement `ioutil.ReadFile` → `os.ReadFile` (12 occurrences)
- ✅ Remplacement `ioutil.WriteFile` → `os.WriteFile` (8 occurrences)
- ✅ Mise à jour imports (auto avec `goimports`)

### Dépendances
- ✅ `github.com/stretchr/testify` : v1.8.0 → v1.8.4
- ✅ `github.com/pigeon` : v1.0.0 → v1.1.0
- ✅ Aucun conflit de dépendances

### Infrastructure
- ✅ GitHub Actions : `go-version: '1.21'`
- ✅ Dockerfile : `FROM golang:1.21-alpine`
- ✅ Makefile : Documentation mise à jour

### Documentation
- ✅ README.md : Prérequis Go 1.21+
- ✅ CHANGELOG.md : Entrée migration ajoutée
- ✅ docs/MIGRATION_GO_1.21.md : Guide créé
- ✅ CONTRIBUTING.md : Prérequis mis à jour

## ✅ Validation

### Tests
```bash
$ make test
✅ Tests unitaires : 234/234 PASS

$ make test-integration
✅ Tests intégration : 45/45 PASS

$ make rete-unified
✅ Runner universel : 58/58 PASS
```

### Performance
```
Benchmark Comparaison (baseline vs after)

BenchmarkAlphaNode-8          1000 ns/op → 950 ns/op    (-5%)  ✅
BenchmarkJoinNode-8           5000 ns/op → 5100 ns/op   (+2%)  ✅
BenchmarkPropagation-8       10000 ns/op → 9500 ns/op   (-5%)  ✅

Performance globale : +5% improvement ✅
```

### Analyse Statique
```bash
$ go vet ./...
✅ Aucun warning

$ golangci-lint run ./...
✅ Aucune erreur nouvelle (3 warnings existants conservés)

$ go fmt ./...
✅ Code formaté

$ goimports -w .
✅ Imports organisés
```

## 📊 Statistiques

**Fichiers modifiés** : 25
- Code : 20 fichiers (.go)
- Config : 3 fichiers (go.mod, .github/workflows/, Dockerfile)
- Docs : 2 fichiers (README.md, CHANGELOG.md)

**Lignes modifiées** :
- Ajoutées : 150
- Supprimées : 130
- Net : +20 lignes

**Commits** :
1. `chore(deps): update Go to 1.21`
2. `refactor: replace deprecated ioutil with os/io`
3. `chore(deps): update external dependencies`
4. `ci: update CI/CD to Go 1.21`
5. `docs: update documentation for Go 1.21 migration`

## 🐛 Problèmes Rencontrés

### Problème 1 : Test flaky après migration
**Description** : TestConcurrentPropagation échouait sporadiquement

**Cause** : Changement de comportement scheduler Go 1.21

**Solution** : Ajout d'un `runtime.Gosched()` pour garantir ordonnancement

**Status** : ✅ Résolu

### Problème 2 : Warning golangci-lint
**Description** : Nouveau warning "SA1019: ioutil.ReadAll is deprecated"

**Cause** : Un appel oublié dans test helper

**Solution** : Remplacé par `io.ReadAll`

**Status** : ✅ Résolu

## 💡 Leçons Apprises

**Ce qui a bien marché** :
- ✅ Migration incrémentale (une phase à la fois)
- ✅ Tests après chaque changement
- ✅ Branche dédiée + backup tag
- ✅ Guide de migration pour l'équipe

**Ce qui pourrait être amélioré** :
- ⚠️  Script automatique de remplacement (sed/awk)
- ⚠️  Tests de performance en CI (pas juste local)
- ⚠️  Communication plus tôt avec l'équipe

## 🔗 Ressources

**Documentation** :
- [Go 1.21 Release Notes](https://go.dev/doc/go1.21)
- [Guide Migration Interne](docs/MIGRATION_GO_1.21.md)

**Pull Requests** :
- [#123: Migrate to Go 1.21](https://github.com/user/tsd/pull/123)

**Commits** :
- abc123: chore(deps): update Go to 1.21
- def456: refactor: replace deprecated ioutil
- ghi789: docs: update migration guide

## 📅 Timeline

- **Préparation** : 2h (recherche, plan)
- **Environnement** : 1h (installation, CI/CD)
- **Code** : 4h (adaptation, dépendances)
- **Tests** : 2h (validation, benchmarks)
- **Documentation** : 1h (README, CHANGELOG, guide)

**Total** : 10h (2 jours)

## ✅ Prêt pour Production

- [x] Tous les tests passent
- [x] Performance validée
- [x] Documentation complète
- [x] Équipe informée
- [x] CI/CD fonctionnel
- [x] Guide de migration disponible
- [x] Rollback plan préparé
- [x] Code review approuvée
- [x] Mergé dans main
```

## Exemple d'Utilisation

```
Je veux migrer le projet TSD de Go 1.19 vers Go 1.21 pour bénéficier
des dernières optimisations et features du langage. Je veux m'assurer
que tous les tests passent et qu'il n'y a pas de régression.

Utilise le prompt "migrate" pour m'accompagner dans cette migration.
```

## Checklist de Migration

### Avant Migration

- [ ] Baseline établie (tests, benchmarks, warnings)
- [ ] Breaking changes identifiés et documentés
- [ ] Plan de migration créé avec timeline
- [ ] Backup créé (branche + tag)
- [ ] Équipe informée de la migration à venir
- [ ] Rollback plan préparé

### Pendant Migration

- [ ] Environnement mis à jour (Go, CI/CD)
- [ ] `go.mod` mis à jour
- [ ] Code déprécié remplacé
- [ ] Dépendances mises à jour
- [ ] Tests exécutés après chaque étape
- [ ] Commits atomiques avec messages clairs

### Après Migration

- [ ] Tous tests passent (unitaires + intégration + runner)
- [ ] Benchmarks comparés (acceptable ±10%)
- [ ] Analyse statique sans nouveaux warnings
- [ ] Documentation mise à jour
- [ ] Guide de migration créé
- [ ] Équipe informée et guidée
- [ ] Code review et approbation
- [ ] Mergé dans main

## Commandes Utiles

```bash
# Version actuelle
go version

# Mettre à jour go.mod
go mod edit -go=1.21
go mod tidy

# Mettre à jour dépendances
go get -u ./...
go mod tidy

# Rechercher code déprécié
grep -r "ioutil\." --include="*.go" .

# Remplacer automatiquement
find . -name "*.go" -exec sed -i 's/ioutil\.ReadFile/os.ReadFile/g' {} \;
goimports -w .

# Tests
make test
make test-integration
make rete-unified

# Benchmarks
go test -bench=. -benchmem ./rete > bench.txt
benchcmp before.txt after.txt

# Analyse
go vet ./...
golangci-lint run ./...

# Build
go build ./...

# Nettoyer cache (si problèmes)
go clean -modcache
go mod download
```

## Bonnes Pratiques

### Migration

- **Incrémentale** : Une étape à la fois, tester à chaque étape
- **Backup** : Toujours avoir moyen de rollback rapidement
- **Documentation** : Documenter chaque changement et sa raison
- **Communication** : Tenir l'équipe informée avant, pendant, après
- **Validation** : Tests rigoureux à chaque étape

### Tests

- **Complets** : Unitaires + intégration + end-to-end
- **Performance** : Benchmarks avant/après, comparer
- **Régression** : Suite complète de tests de non-régression
- **Environnements** : Tester en local + CI + staging

### Code

- **Atomique** : Commits petits et focalisés
- **Réversible** : Chaque commit peut être reverté
- **Documenté** : Messages de commit explicites
- **Reviewé** : Code review avant merge

## Anti-Patterns à Éviter

### ❌ Big Bang Migration
```
❌ Tout migrer d'un coup sans tester
✅ Migration incrémentale avec validation à chaque étape
```

### ❌ Pas de Backup
```
❌ Migrer directement sur main sans backup
✅ Créer branche backup + tag avant migration
```

### ❌ Ignorer Breaking Changes
```
❌ "Ça va probablement marcher"
✅ Rechercher et documenter tous les breaking changes
```

### ❌ Migration Sans Tests
```
❌ Migrer et espérer que ça marche
✅ Tests rigoureux après chaque modification
```

### ❌ Pas de Communication
```
❌ Migrer en silence, surprendre l'équipe
✅ Communiquer avant, pendant, et après
```

### ❌ Pas de Guide
```
❌ "Lisez les release notes de Go"
✅ Créer guide spécifique au projet
```

## Outils Recommandés

### Migration Automatique
- `gofmt` - Formatage code
- `goimports` - Organisation imports
- `sed` / `awk` - Remplacement automatique
- Scripts custom - Migrations spécifiques

### Validation
- `go test` - Tests unitaires
- `go test -bench` - Benchmarks
- `go vet` - Analyse statique
- `golangci-lint` - Linter complet
- `go test -race` - Race detector

### Versioning
- `go mod edit` - Édition go.mod
- `go mod tidy` - Nettoyage dépendances
- `go list -u -m all` - Dépendances obsolètes

## Ressources

- [Makefile](../../Makefile) - Commandes disponibles
- [Go Release History](https://go.dev/doc/devel/release) - Versions Go
- [Go Release Policy](https://go.dev/doc/devel/release#policy) - Support timeline
- [Migration Guides](https://go.dev/doc/) - Guides officiels

---

**Version** : 1.0  
**Dernière mise à jour** : Novembre 2025  
**Mainteneur** : Équipe TSD