# 🤝 Contributing to TSD

Merci de votre intérêt pour contribuer au projet TSD ! Ce guide vous aidera à soumettre des contributions de qualité.

---

## 📋 Table des Matières

- [Code de Conduite](#code-de-conduite)
- [Comment Contribuer](#comment-contribuer)
- [Setup Environnement](#setup-environnement)
- [Workflow de Contribution](#workflow-de-contribution)
- [Standards de Code](#standards-de-code)
- [Standards de Tests](#standards-de-tests)
- [Standards de Documentation](#standards-de-documentation)
- [Process de Review](#process-de-review)
- [Licence et Copyright](#licence-et-copyright)
- [Ressources](#ressources)

---

## 🤝 Code de Conduite

### Nos Standards

En participant à ce projet, vous vous engagez à :

- ✅ Être respectueux et inclusif
- ✅ Accepter les critiques constructives
- ✅ Collaborer de bonne foi
- ✅ Respecter les opinions divergentes
- ❌ Ne pas harceler ou discriminer
- ❌ Ne pas publier d'informations privées d'autrui

### Reporting

Si vous observez un comportement inapproprié, contactez les mainteneurs du projet.

---

## 🔒 Reporting de Vulnérabilités de Sécurité

**⚠️ Important : Ne reportez JAMAIS de vulnérabilités de sécurité via des issues publiques GitHub.**

Si vous découvrez une vulnérabilité de sécurité dans TSD :

1. **NE PAS** créer d'issue publique
2. **Consultez** notre [Security Policy](SECURITY.md)
3. **Utilisez** GitHub Security Advisory (recommandé)
4. **Ou contactez** directement les mainteneurs de manière privée

Notre [Security Policy](SECURITY.md) détaille :
- Comment reporter de manière responsable
- Nos délais de réponse
- Le processus de gestion des vulnérabilités
- La politique de divulgation coordonnée

**Merci de protéger les utilisateurs de TSD en suivant cette procédure.**

---

## 🎯 Comment Contribuer

### Types de Contributions Acceptées

- 🐛 **Bug Reports** - Signaler un problème
- ✨ **Features** - Proposer une amélioration
- 📝 **Documentation** - Améliorer la doc
- 🧪 **Tests** - Ajouter/améliorer tests
- 🔧 **Code** - Corriger bugs, implémenter features

### Avant de Commencer

1. **Chercher des issues existantes** : Éviter les doublons
2. **Good First Issue** : Issues taguées pour débutants
3. **Help Wanted** : Issues où de l'aide est bienvenue

---

## 📋 Utilisation des Templates

### Créer une Issue

Utilisez les templates appropriés :

- **🐛 Bug Report** : Pour signaler un bug
- **✨ Feature Request** : Pour proposer une amélioration
- **❓ Question** : Pour poser une question

Les templates vous guideront pour fournir toutes les informations nécessaires.

#### Proposer une Feature

1. Créer une issue avec le template "✨ Feature Request"
2. Décrire clairement le besoin et la solution proposée
3. Attendre validation avant d'implémenter
4. Discuter de l'approche avec les mainteneurs

#### Reporter un Bug

1. Créer une issue avec le template "🐛 Bug Report"
2. Fournir :
   - Description du problème
   - Étapes de reproduction
   - Comportement attendu vs observé
   - Version de Go, OS, etc.
   - Logs/erreurs si disponibles

### Créer une Pull Request

Le template de PR inclut une checklist complète. Assurez-vous de :

1. Remplir toutes les sections
2. Cocher tous les items de la checklist
3. Lier l'issue correspondante
4. Fournir des instructions de test

---

## 🛠️ Setup Environnement

### Prerequisites

- **Go** : Version 1.24+ (voir `go.mod`)
- **Make** : Pour utiliser le Makefile
- **Git** : Pour le versioning
- **Outils Go** (installés automatiquement) :
  - `goimports`
  - `staticcheck`
  - `errcheck`
  - `golangci-lint`
  - `gosec`
  - `govulncheck`

### Initial Setup

```bash
# 1. Fork le repo sur GitHub
# 2. Clone votre fork
git clone https://github.com/VOTRE_USERNAME/tsd.git
cd tsd

# 3. Ajouter le repo upstream
git remote add upstream https://github.com/treivax/tsd.git

# 4. Installer les dépendances
go mod download

# 5. Installer les outils de développement
make deps-dev

# 6. Installer les hooks pre-commit (optionnel mais recommandé)
pip install pre-commit  # Si pas déjà installé
make install-hooks

# 7. Build le projet
make build

# 8. Vérifier l'installation
make validate
```

### Verify Installation

```bash
# Quick validation (format + lint + build)
make quick-check

# Full validation (includes all tests)
make validate
```

### Commandes Utiles

```bash
# Build
make build              # Compiler le binaire TSD

# Tests
make test-unit         # Tests unitaires (rapides)
make test-fixtures     # Tests des fixtures partagées
make test-e2e          # Tests E2E (fixtures TSD)
make test-integration  # Tests d'intégration
make test-performance  # Tests de performance
make test-all          # Tous les tests standards
make test-complete     # TOUS les tests (complet, recommandé avant commit)
make coverage-prod     # Rapport de couverture (code production uniquement)

# Validation
make format            # Formater le code
make lint              # Analyse statique
make security-scan     # Scan de sécurité complet
make validate          # Validation complète (format+lint+build+test-complete)
make quick-check       # Validation rapide sans tests

# Nettoyage
make clean             # Nettoyer les artefacts

# Aide
make help              # Liste complète des commandes
```

---

## 🔄 Workflow de Contribution

### 1. Créer une Branche

```bash
# Mettre à jour main
git checkout main
git pull upstream main

# Créer une branche de feature
git checkout -b feature/ma-feature

# OU pour un bugfix
git checkout -b fix/mon-bugfix
```

**Convention de nommage** :
- `feature/description` - Nouvelle fonctionnalité
- `fix/description` - Correction de bug
- `docs/description` - Documentation
- `test/description` - Ajout/amélioration tests
- `refactor/description` - Refactoring
- `perf/description` - Amélioration performance
- `chore/description` - Maintenance (deps, build, etc.)

### 2. Développer

```bash
# Faire vos modifications
# ...

# Formater le code
make format

# Vérifier en continu
make validate
```

### 3. Committer

**Convention de commit** (Conventional Commits) :

```
<type>(<scope>): <description>

[body optionnel]

[footer optionnel]
```

**Types** :
- `feat` : Nouvelle fonctionnalité
- `fix` : Correction de bug
- `docs` : Documentation uniquement
- `test` : Ajout/modification tests
- `refactor` : Refactoring (pas de changement fonctionnel)
- `perf` : Amélioration performance
- `chore` : Maintenance (deps, build, etc.)
- `ci` : CI/CD changes

**Exemples** :

```bash
git commit -m "feat(rete): ajouter support jointures N-variables"
git commit -m "fix(constraint): corriger parsing règles imbriquées"
git commit -m "docs(readme): ajouter section installation"
git commit -m "test(rete): améliorer couverture nodes join"
git commit -m "perf(rete): optimiser propagation tokens avec cache"
```

### 4. Pousser et Créer PR

```bash
# Pousser la branche
git push origin feature/ma-feature

# Aller sur GitHub et créer une Pull Request
# Utiliser le template de PR fourni
```

### 5. Process de Review

1. **CI/CD** : Attendre que les checks passent (GitHub Actions)
2. **Review** : Répondre aux commentaires des reviewers
3. **Modifications** : Pousser les corrections
4. **Approval** : Obtenir l'approbation d'un mainteneur
5. **Merge** : Le mainteneur merge la PR

---

## ✅ Standards de Code

### Documentation Complète

Consulter **obligatoirement** :
- [`.github/prompts/common.md`](.github/prompts/common.md) - Standards du projet ⭐
- [`.github/prompts/develop.md`](.github/prompts/develop.md) - Guide développement
- [`.github/prompts/review.md`](.github/prompts/review.md) - Guide de revue

### Règles Strictes

#### ❌ INTERDICTIONS ABSOLUES

1. **Aucun hardcoding** :
   ```go
   // ❌ INTERDIT
   timeout := 30 * time.Second
   if userID == "special-user-123" { ... }
   
   // ✅ AUTORISÉ
   const DefaultTimeout = 30 * time.Second
   timeout := config.Timeout
   ```

2. **Pas de code spécifique** :
   - Code doit être générique et réutilisable
   - Utiliser paramètres, interfaces, configuration

3. **Pas de modification du code généré** :
   - Ne jamais modifier `parser.go` ou autre code généré

#### ✅ OBLIGATOIRE

1. **En-tête copyright** (tous les fichiers `.go`) :
   ```go
   // Copyright (c) 2025 TSD Contributors
   // Licensed under the MIT License
   // See LICENSE file in the project root for full license text
   
   package monpackage
   ```

2. **Formatage** :
   ```bash
   go fmt ./...
   goimports -w .
   ```

3. **Linting** :
   ```bash
   go vet ./...
   staticcheck ./...
   errcheck ./...
   ```

4. **Visibilité** :
   - **Tout privé par défaut** (non exporté)
   - N'exporter que ce qui fait partie de l'API publique

5. **Gestion d'erreurs** :
   - Jamais ignorer les erreurs
   - Messages clairs et contextuels
   - Wrap avec `fmt.Errorf("context: %w", err)`

### Code Quality Standards

| Element | Convention | Exemple |
|---------|------------|---------|
| Packages | lowercase, singulier | `rete`, `constraint` |
| Files | lowercase_underscore | `node_join.go` |
| Tests | *_test.go | `node_join_test.go` |
| Constants | MixedCaps/UPPER | `MaxNodes`, `DefaultTimeout` |
| Variables | camelCase | `nodeCount`, `resultToken` |
| Functions | MixedCaps | `ProcessToken`, `EvaluateCondition` |
| Types | MixedCaps | `AlphaNode`, `TokenMemory` |
| Interfaces | MixedCaps + "er" | `Evaluator`, `Processor` |

### Qualité du Code

- **Complexité**: Cyclomatique < 15
- **Longueur fonctions**: < 50 lignes (sauf justification)
- **Imbrication**: < 4 niveaux
- **DRY**: Don't Repeat Yourself
- **Single Responsibility**: Une fonction, une responsabilité

### Principes Architecturaux

- **SOLID** : Single Responsibility, Open/Closed, etc.
- **Dependency Injection** : Pas de dépendances globales
- **Composition over Inheritance** : Interfaces et embedding
- **Interfaces** : Petites, focalisées, cohésives
- **Découplage** : Couplage faible, cohésion forte

### Formatting

```bash
# Format code (automatic)
make format

# Check formatting
gofmt -l .

# Import organization
goimports -w .
```

### Linting

```bash
# Run all linters
make lint

# Individual checks
go vet ./...
golangci-lint run
staticcheck ./...
```

---

## 🧪 Standards de Tests

### Obligatoire

- ✅ **Couverture > 80%** sur nouveau code (MANDATORY)
- ✅ **Tests réels** : Extraction résultats réels, PAS de mocks (sauf si explicitement demandé)
- ✅ **Tests déterministes** : Pas de flaky tests
- ✅ **Tests isolés** : Indépendants, pas de dépendances entre tests
- ✅ **Table-driven tests** quand applicable
- ✅ **Messages clairs** avec émojis (✅ ❌ ⚠️)

### Structure de Test

```go
func TestFeature(t *testing.T) {
    t.Log("🧪 TEST FEATURE")
    t.Log("================")
    
    tests := []struct {
        name     string
        input    interface{}
        expected interface{}
        wantErr  bool
    }{
        {"cas nominal", validInput, expectedOutput, false},
        {"cas erreur", invalidInput, nil, true},
        {"cas limite", edgeInput, edgeOutput, false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            // Act
            result, err := Feature(tt.input)
            
            // Assert
            if (err != nil) != tt.wantErr {
                t.Errorf("❌ Erreur = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("❌ Attendu %v, reçu %v", tt.expected, result)
            }
            t.Log("✅ Test réussi")
        })
    }
}
```

### Running Tests

```bash
# Tests unitaires (rapides)
make test-unit

# Tests fixtures partagées
make test-fixtures

# Tests E2E (fixtures TSD)
make test-e2e

# Tests d'intégration
make test-integration

# Tests de performance
make test-performance

# Tous les tests standards
make test-all

# TOUS les tests (validation complète)
make test-complete

# Rapport de couverture (code production uniquement)
make coverage-prod
```

### Test Categories

- **Unit tests** (`*_test.go` dans les modules) - Tester fonctions individuelles
- **Fixtures tests** (`tests/fixtures/`) - Tests des fixtures partagées
- **E2E tests** (`tests/e2e/`) - Tests end-to-end avec fichiers `.tsd`
- **Integration tests** (`tests/integration/`) - Tests d'interaction entre modules
- **Performance tests** (`tests/performance/`) - Tests de charge et benchmarks

---

## 📚 Standards de Documentation

### GoDoc

Documenter **toutes** les fonctions/types exportés :

```go
// ExecuteProgram compile et exécute un programme TSD.
// 
// Le programme est d'abord parsé, puis compilé en réseau RETE,
// et enfin exécuté avec les données fournies.
//
// Paramètres :
//   - program : Code source TSD à exécuter
//   - data : Données d'entrée (peut être nil)
//
// Retourne :
//   - Résultat de l'exécution
//   - Erreur si parsing, compilation ou exécution échoue
func ExecuteProgram(program string, data map[string]interface{}) (*Result, error) {
    // ...
}
```

### README

- Mettre à jour si changement d'API
- Ajouter exemples d'utilisation
- Documenter nouvelles commandes CLI

### CHANGELOG.md

Ajouter entrée dans la section `[Unreleased]` :

```markdown
## [Unreleased]

### Added
- Support jointures N-variables avec contraintes multiples

### Fixed
- Correction parsing règles avec négation imbriquée

### Changed
- Amélioration messages d'erreur compilation
```

---

## 🔍 Process de Review

### Checklist Avant PR

Vérifier **TOUJOURS** avant de soumettre :

- [ ] ✅ **Copyright header** présent dans tous les nouveaux fichiers `.go`
- [ ] ✅ **Aucun hardcoding** (valeurs, chemins, configs)
- [ ] ✅ **Code générique** avec paramètres/interfaces
- [ ] ✅ **Constantes nommées** pour toutes les valeurs
- [ ] ✅ **Formatage** : `make format` appliqué
- [ ] ✅ **Linting** : `make lint` passe sans erreur
- [ ] ✅ **Sécurité** : `make security-scan` passe
- [ ] ✅ **Tests** : `make test-complete` passe
- [ ] ✅ **Couverture** : ≥ 80% maintenue (vérifier avec `make coverage-prod`)
- [ ] ✅ **Documentation** : GoDoc + README mis à jour si nécessaire
- [ ] ✅ **CHANGELOG.md** : Entrée ajoutée sous `[Unreleased]`
- [ ] ✅ **Branch à jour** avec main

### En Review

Les reviewers vérifient :

1. **Conformité standards** (common.md)
2. **Qualité code** (lisibilité, simplicité)
3. **Tests** (couverture, pertinence)
4. **Documentation** (clarté, complétude)
5. **Pas de régression** (CI passe)

### Répondre aux Commentaires

- ✅ Être constructif et ouvert
- ✅ Expliquer vos choix si nécessaire
- ✅ Implémenter les suggestions raisonnables
- ✅ Demander clarification si besoin
- ❌ Ne pas prendre personnellement
- ❌ Ne pas argumenter sans raison valable

### Merge

Une PR est mergée quand :

1. ✅ Tous les checks CI passent
2. ✅ Au moins 1 approbation d'un mainteneur
3. ✅ Tous les commentaires résolus
4. ✅ Branch à jour avec main

---

## 📜 Licence et Copyright

### Licence du Projet

TSD est sous **MIT License**. Toute contribution sera sous cette même licence.

### En-tête Copyright

**Obligatoire** dans tous les nouveaux fichiers `.go` :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
```

### Vérification Avant Commit

```bash
# Vérifier que tous les fichiers ont le copyright
for file in $(find . -name "*.go" -type f ! -path "./.git/*"); do
    if ! head -1 "$file" | grep -q "Copyright\|Code generated"; then
        echo "⚠️  EN-TÊTE MANQUANT: $file"
    fi
done
```

### Dépendances Externes

Avant d'ajouter une dépendance, vérifier :

- ✅ Licence compatible (MIT, BSD, Apache 2.0)
- ❌ Éviter GPL, AGPL (copyleft)
- 📝 Documenter dans `THIRD_PARTY_LICENSES.md`

---

## 📚 Ressources

### Documentation Projet

- [common.md](.github/prompts/common.md) - Standards du projet ⭐
- [develop.md](.github/prompts/develop.md) - Guide développement
- [review.md](.github/prompts/review.md) - Guide de revue
- [test.md](.github/prompts/test.md) - Guide tests
- [Documentation technique](docs/) - Architecture et guides

### Outils

- [Makefile](Makefile) - Toutes les commandes
- [GitHub Actions](.github/workflows/) - CI/CD
- [Issues](../../issues) - Bugs et features

### Ressources Go

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

### Communication

- **Issues** : Pour bugs, features, questions
- **Discussions** : Pour discussions générales
- **Pull Requests** : Pour proposer du code

---

## 🎉 Merci !

Merci de contribuer à TSD ! Chaque contribution, petite ou grande, est appréciée.

**Questions ?** N'hésitez pas à ouvrir une issue ou une discussion.

**Débutant ?** Cherchez les issues taguées `good first issue`.

---

**Happy Coding! 🚀**
