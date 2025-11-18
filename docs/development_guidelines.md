# GUIDE DE DÉVELOPPEMENT GO - PROJET TSD

**Version :** 1.0
**Date :** 13 novembre 2025
**Statut :** Standards obligatoires pour tous les développeurs

## 🎯 CONVENTIONS DE NOMMAGE OBLIGATOIRES

### 📁 **Noms de Fichiers**
```bash
# ✅ CORRECT - snake_case
user_service.go
constraint_parser.go
rete_network.go
test_utils.go

# ❌ INCORRECT - camelCase
userService.go
constraintParser.go
reteNetwork.go
```

### 🏷️ **Types et Structures**
```go
// ✅ CORRECT - PascalCase pour types exportés
type UserService struct {
    config Config
}

type ConstraintValidator interface {
    Validate(input interface{}) error
}

// ✅ CORRECT - camelCase pour types internes
type parserState struct {
    position int
}

// ❌ INCORRECT - snake_case
type user_service struct { } // INTERDIT
type constraint_validator interface { } // INTERDIT
```

### 🔧 **Fonctions et Méthodes**
```go
// ✅ CORRECT - PascalCase pour exportées
func (s *UserService) ValidateUser(user User) error { }
func NewConstraintValidator() *ConstraintValidator { }

// ✅ CORRECT - camelCase pour privées
func (s *userService) parseInput(input string) error { }
func createConnection() *Connection { }

// ❌ INCORRECT - snake_case (sauf tests)
func validate_user(user User) error { } // INTERDIT
func create_connection() *Connection { } // INTERDIT

// ✅ EXCEPTION - Tests peuvent utiliser snake_case
func TestUserService_ValidateUser(t *testing.T) { } // OK dans tests
```

### 🔀 **Variables et Constantes**
```go
// ✅ CORRECT - camelCase
var globalConfig Config
var userCount int
const defaultTimeout = 30

// ✅ CORRECT - Constantes exportées
const MaxRetries = 3
const DefaultBufferSize = 1024

// ❌ INCORRECT - snake_case ou mixed case
var global_config Config // INTERDIT
var Global_Config Config // INTERDIT
const default_timeout = 30 // INTERDIT
```

### 📂 **Packages et Répertoires**
```bash
# ✅ CORRECT - snake_case, descriptifs
pkg/domain/
pkg/validator/
internal/config/
test/integration/
cmd/constraint_parser/

# ❌ INCORRECT - camelCase ou trop génériques
pkg/Domain/ # INTERDIT
internal/Config/ # INTERDIT
test/Integration/ # INTERDIT
cmd/constraintParser/ # INTERDIT
```

## 🏗️ **ARCHITECTURE RECOMMANDÉE**

### Structure des Packages
```
project/
├── cmd/                    # Applications principales
│   └── app_name/
├── internal/               # Code privé au projet
│   ├── config/
│   └── service/
├── pkg/                    # Code réutilisable
│   ├── domain/            # Types métier
│   ├── validator/         # Validation
│   └── storage/           # Persistence
├── test/                  # Tests centralisés
│   ├── integration/
│   ├── unit/
│   └── benchmark/
└── scripts/               # Outils et utilitaires
```

### Nommage des Packages
```go
// ✅ CORRECT - Noms simples et descriptifs
package domain
package validator
package storage

// ❌ INCORRECT - Noms génériques ou répétitifs
package utils          // Trop générique
package domainservice  // Répétitif
package validator_pkg  // Snake case interdit
```

## 🧪 **CONVENTIONS DE TESTS**

### Structure des Tests
```go
// ✅ CORRECT - Organisation claire
func TestUserService_ValidateUser(t *testing.T) {
    t.Run("valid_user_should_pass", func(t *testing.T) {
        // Test avec snake_case OK
    })

    t.Run("invalid_user_should_fail", func(t *testing.T) {
        // Test avec snake_case OK
    })
}

// ✅ CORRECT - Helpers de test
func createTestUser(name string) *User {
    return &User{Name: name}
}

// ✅ CORRECT - Benchmarks
func BenchmarkUserService_ValidateUser(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Benchmark code
    }
}
```

### Organisation des Fichiers de Test
```bash
# ✅ CORRECT - Co-localisation avec suffixe
user_service.go
user_service_test.go

# ✅ CORRECT - Tests d'intégration séparés
test/integration/user_integration_test.go
test/unit/user_unit_test.go
```

## 🔍 **VALIDATION AUTOMATIQUE**

### Pre-commit Hook
Le projet inclut un hook pre-commit qui valide automatiquement :
- ✅ Noms de fichiers en snake_case
- ✅ Types en PascalCase
- ✅ Fonctions correctement nommées
- ✅ Variables en camelCase
- ✅ Compilation sans erreurs
- ✅ Tests rapides

### Scripts de Validation
```bash
# Analyser les conventions du projet
./scripts/analyze_naming.sh

# Valider la conformité globale
./scripts/validate_conventions.sh

# Rapport de validation final
./scripts/final_validation_report.sh
```

## 📋 **CHECKLIST DÉVELOPPEUR**

Avant chaque commit, vérifier :

### ✅ **Code**
- [ ] Noms de fichiers en snake_case
- [ ] Types exportés en PascalCase
- [ ] Fonctions exportées en PascalCase
- [ ] Variables en camelCase
- [ ] Pas de snake_case dans les fonctions (sauf tests)

### ✅ **Tests**
- [ ] Tests ajoutés pour nouveaux features
- [ ] Tests existants passent
- [ ] Noms de tests descriptifs avec snake_case

### ✅ **Documentation**
- [ ] Commentaires godoc pour types/fonctions exportés
- [ ] README mis à jour si nécessaire
- [ ] Exemples d'utilisation si API publique

### ✅ **Build**
- [ ] `go build ./...` passe sans erreurs
- [ ] `go test ./...` passe sans erreurs
- [ ] `go vet ./...` ne rapporte aucun problème
- [ ] `gofmt -d .` ne montre aucune différence

## 🚨 **ERREURS FRÉQUENTES À ÉVITER**

### ❌ **Anti-patterns de Nommage**
```go
// ERREUR: Mélange de conventions
type User_Service struct { } // snake_case interdit pour types
func Process_User() { }      // snake_case interdit pour fonctions
var Global_Config Config    // mixed case interdit

// ERREUR: Noms non descriptifs
func DoStuff() { }          // Nom trop vague
var data interface{}        // Nom trop générique
type Thing struct { }       // Nom non descriptif
```

### ❌ **Anti-patterns de Structure**
```go
// ERREUR: Package mal nommé
package UtilityFunctions    // PascalCase interdit pour packages
package constraint_utils    // snake_case interdit pour packages

// ERREUR: Récepteurs mal nommés
func (constraintValidator *ConstraintValidator) Validate() {} // Trop verbeux
func (cv ConstraintValidator) Validate() {}                  // Devrait être pointeur
func (this *ConstraintValidator) Validate() {}               // "this" interdit en Go
```

## 🎯 **OBJECTIFS DE QUALITÉ**

### Métriques Cibles
- **Conformité nommage :** > 90%
- **Couverture de tests :** > 80%
- **Compilation sans warnings :** 100%
- **Tests passants :** 100%

### Outils Recommandés
```bash
# Linting et formatage
go vet ./...
gofmt -w .
golangci-lint run

# Tests et couverture
go test -v ./...
go test -cover ./...
go test -race ./...
```

## 📚 **RESSOURCES**

### Documentation Officielle
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)
- [Package Names](https://go.dev/blog/package-names)

### Standards du Projet
- Voir `NAMING_CONVENTIONS_FINAL_REPORT.md` pour l'état actuel
- Voir `REFACTORING_PLAN.md` pour l'architecture
- **Voir `CODE_GENERATION_CONVENTIONS.md` pour la génération de code automatique**

---

## 🤖 **GÉNÉRATION DE CODE AUTOMATIQUE**

**RÈGLE CRITIQUE :** Toute génération de code automatique (IA, outils, templates) doit respecter les conventions TSD.

### Application Systématique
- ✅ **Types exportés :** PascalCase uniquement
- ✅ **Fonctions exportées :** PascalCase, privées camelCase
- ✅ **Variables :** camelCase systématique
- ✅ **Fichiers :** snake_case (convention TSD)
- ✅ **Tests :** pattern TestType_Method
- ✅ **Context :** premier paramètre si applicable
- ✅ **Error handling :** avec wrapping et validation

### Référence Complète
**Document obligatoire :** `CODE_GENERATION_CONVENTIONS.md`

---

## 🏁 **CONCLUSION**

**Ce guide est OBLIGATOIRE** pour maintenir la qualité et la cohérence du projet TSD.

Le hook pre-commit aide à détecter automatiquement les violations, mais chaque développeur est responsable de respecter ces conventions.

**En cas de doute :** Suivre les conventions existantes dans le projet et consulter la documentation Go officielle.

*Dernière mise à jour : 13 novembre 2025*
