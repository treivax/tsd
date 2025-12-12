# 📖 Guide de Migration - Validator Package Refactoring

## 🎯 Objectif

Ce guide aide à migrer le code existant vers la nouvelle API du package `constraint/pkg/validator` suite au refactoring de la Session 3.

---

## 📋 Résumé des Changements

### Breaking Changes
1. **NewConstraintValidator** : Signature modifiée (injection de configuration)
2. **ActionValidator** : Struct supprimée, remplacée par fonctions pures

### Nouveaux Exports
- `NewConstraintValidatorWithDefaults()` : Constructeur avec config par défaut
- `ValidateAction()` : Fonction pure de validation d'action
- `ValidateJobCall()` : Fonction pure de validation de job call
- Constantes : `ComparisonOperators`, `LogicalOperators`, `ArithmeticOperators`, `OrderableTypes`, `NumericTypes`

---

## 🔄 Migration Pas à Pas

### Étape 1 : Migration NewConstraintValidator

#### Scénario A : Migration Simple (Recommandé)

**Avant :**
```go
import "github.com/treivax/tsd/constraint/pkg/validator"

registry := validator.NewTypeRegistry()
checker := validator.NewTypeChecker(registry)
v := validator.NewConstraintValidator(registry, checker)
```

**Après :**
```go
import "github.com/treivax/tsd/constraint/pkg/validator"

registry := validator.NewTypeRegistry()
checker := validator.NewTypeChecker(registry)
v := validator.NewConstraintValidatorWithDefaults(registry, checker)
```

**Changement :** Remplacer `NewConstraintValidator` par `NewConstraintValidatorWithDefaults`

#### Scénario B : Configuration Personnalisée

**Si vous voulez personnaliser la configuration :**

```go
import (
    "github.com/treivax/tsd/constraint/pkg/domain"
    "github.com/treivax/tsd/constraint/pkg/validator"
)

registry := validator.NewTypeRegistry()
checker := validator.NewTypeChecker(registry)

// Définir configuration personnalisée
config := domain.ValidatorConfig{
    StrictMode: true,
    AllowedOperators: []string{
        "==", "!=", "<", ">", "<=", ">=",
        "AND", "OR", "NOT",
        "+", "-", "*", "/", "%",
    },
    MaxDepth: 20, // Au lieu de 10 par défaut
}

v := validator.NewConstraintValidator(registry, checker, config)
```

### Étape 2 : Migration ActionValidator

#### Avant :
```go
import "github.com/treivax/tsd/constraint/pkg/validator"

av := validator.NewActionValidator()

// Validation d'action
err := av.ValidateAction(action)
if err != nil {
    return err
}

// Validation de job call
err = av.ValidateJobCall(jobCall)
if err != nil {
    return err
}
```

#### Après :
```go
import "github.com/treivax/tsd/constraint/pkg/validator"

// Plus besoin de créer une instance

// Validation d'action
err := validator.ValidateAction(action)
if err != nil {
    return err
}

// Validation de job call
err = validator.ValidateJobCall(jobCall)
if err != nil {
    return err
}
```

**Changement :** Appeler directement les fonctions, supprimer `NewActionValidator()`

---

## 🔍 Exemples Complets

### Exemple 1 : Validator Simple

**Code Original :**
```go
package myapp

import (
    "github.com/treivax/tsd/constraint/pkg/validator"
    "github.com/treivax/tsd/constraint/pkg/domain"
)

func setupValidator() *validator.ConstraintValidator {
    registry := validator.NewTypeRegistry()
    checker := validator.NewTypeChecker(registry)
    return validator.NewConstraintValidator(registry, checker)
}
```

**Code Migré :**
```go
package myapp

import (
    "github.com/treivax/tsd/constraint/pkg/validator"
    "github.com/treivax/tsd/constraint/pkg/domain"
)

func setupValidator() *validator.ConstraintValidator {
    registry := validator.NewTypeRegistry()
    checker := validator.NewTypeChecker(registry)
    return validator.NewConstraintValidatorWithDefaults(registry, checker)
}
```

### Exemple 2 : Validation d'Actions

**Code Original :**
```go
package myapp

import (
    "github.com/treivax/tsd/constraint/pkg/validator"
    "github.com/treivax/tsd/constraint/pkg/domain"
)

func validateProgramActions(program *domain.Program) error {
    av := validator.NewActionValidator()
    
    for _, expr := range program.Expressions {
        if expr.Action != nil {
            if err := av.ValidateAction(expr.Action); err != nil {
                return err
            }
        }
    }
    
    return nil
}
```

**Code Migré :**
```go
package myapp

import (
    "github.com/treivax/tsd/constraint/pkg/validator"
    "github.com/treivax/tsd/constraint/pkg/domain"
)

func validateProgramActions(program *domain.Program) error {
    // Plus besoin de créer NewActionValidator()
    
    for _, expr := range program.Expressions {
        if expr.Action != nil {
            if err := validator.ValidateAction(expr.Action); err != nil {
                return err
            }
        }
    }
    
    return nil
}
```

### Exemple 3 : Tests

**Code Original :**
```go
func TestMyValidator(t *testing.T) {
    registry := validator.NewTypeRegistry()
    checker := validator.NewTypeChecker(registry)
    v := validator.NewConstraintValidator(registry, checker)
    
    // ... tests
}
```

**Code Migré :**
```go
func TestMyValidator(t *testing.T) {
    registry := validator.NewTypeRegistry()
    checker := validator.NewTypeChecker(registry)
    v := validator.NewConstraintValidatorWithDefaults(registry, checker)
    
    // ... tests
}
```

---

## 🛠️ Script de Migration Automatique

Pour faciliter la migration, voici un script bash :

```bash
#!/bin/bash
# migration_validator.sh

# Remplacer NewConstraintValidator par NewConstraintValidatorWithDefaults
find . -name "*.go" -type f -exec sed -i 's/validator\.NewConstraintValidator(/validator.NewConstraintValidatorWithDefaults(/g' {} \;

# Supprimer les lignes contenant NewActionValidator
find . -name "*.go" -type f -exec sed -i '/av := validator\.NewActionValidator()/d' {} \;
find . -name "*.go" -type f -exec sed -i '/actionValidator := validator\.NewActionValidator()/d' {} \;

# Remplacer av.ValidateAction par validator.ValidateAction
find . -name "*.go" -type f -exec sed -i 's/av\.ValidateAction(/validator.ValidateAction(/g' {} \;

# Remplacer av.ValidateJobCall par validator.ValidateJobCall  
find . -name "*.go" -type f -exec sed -i 's/av\.ValidateJobCall(/validator.ValidateJobCall(/g' {} \;

echo "✅ Migration terminée"
echo "⚠️  Vérifiez manuellement les changements avec: git diff"
```

**Utilisation :**
```bash
chmod +x migration_validator.sh
./migration_validator.sh
git diff  # Vérifier les changements
go test ./...  # Valider que tout fonctionne
```

---

## ✅ Checklist de Migration

Pour chaque fichier à migrer :

- [ ] Identifier les usages de `NewConstraintValidator`
- [ ] Remplacer par `NewConstraintValidatorWithDefaults` ou avec config
- [ ] Identifier les usages de `NewActionValidator()`
- [ ] Supprimer les lignes `av := NewActionValidator()`
- [ ] Remplacer `av.ValidateAction` par `validator.ValidateAction`
- [ ] Remplacer `av.ValidateJobCall` par `validator.ValidateJobCall`
- [ ] Compiler : `go build ./...`
- [ ] Tester : `go test ./...`
- [ ] Vérifier : `go vet ./...`

---

## 🔍 Vérification Post-Migration

### 1. Rechercher les usages restants

```bash
# Rechercher NewConstraintValidator (devrait être remplacé)
grep -rn "NewConstraintValidator(" --include="*.go" . | grep -v "WithDefaults"

# Rechercher NewActionValidator (devrait être supprimé)
grep -rn "NewActionValidator" --include="*.go" .

# Rechercher av.Validate ou actionValidator.Validate
grep -rn "av\.Validate\|actionValidator\.Validate" --include="*.go" .
```

Si ces recherches retournent des résultats, il reste du code à migrer.

### 2. Valider la compilation

```bash
go build ./...
```

### 3. Valider les tests

```bash
go test ./...
```

### 4. Vérifier la couverture

```bash
go test -cover ./constraint/pkg/validator/...
```

Devrait être ≥ 80%

---

## 🆘 Résolution de Problèmes

### Erreur : "not enough arguments in call to validator.NewConstraintValidator"

**Cause :** Appel à `NewConstraintValidator` sans le paramètre config

**Solution :**
```go
// Option 1 (Recommandé)
v := validator.NewConstraintValidatorWithDefaults(registry, checker)

// Option 2
config := domain.ValidatorConfig{ /* ... */ }
v := validator.NewConstraintValidator(registry, checker, config)
```

### Erreur : "undefined: validator.NewActionValidator"

**Cause :** `NewActionValidator` a été supprimé

**Solution :**
```go
// Avant
av := validator.NewActionValidator()
err := av.ValidateAction(action)

// Après
err := validator.ValidateAction(action)
```

### Erreur : "av.ValidateAction undefined (type has no method)"

**Cause :** Variable `av` référence un ActionValidator qui n'existe plus

**Solution :**
```go
// Supprimer la ligne
av := validator.NewActionValidator()

// Remplacer
av.ValidateAction(action)
// Par
validator.ValidateAction(action)
```

---

## 📚 Références

- Rapport de revue : `REPORTS/REVIEW_CONSTRAINT_SESSION_3_PKG_VALIDATOR.md`
- Résumé refactoring : `REPORTS/REFACTORING_SESSION_3_SUMMARY.md`
- TODO technique : `constraint/pkg/validator/TODO_REFACTORING.md`
- Standards projet : `.github/prompts/common.md`

---

## 💬 Support

En cas de problème lors de la migration :

1. Vérifier ce guide
2. Consulter les exemples ci-dessus
3. Vérifier les tests existants dans `constraint/pkg/validator/*_test.go`
4. Consulter le rapport de revue pour comprendre les changements

---

**Date** : 2025-12-11  
**Version** : 1.0  
**Auteur** : GitHub Copilot CLI
