# 🎯 RÈGLES DE GÉNÉRATION DE CODE - CONVENTIONS TSD GO

**Date :** 13 novembre 2025  
**Scope :** Toute génération de code dans le projet TSD
**Status :** **RÈGLES SYSTÉMATIQUES ACTIVES**

## 📋 CONVENTIONS SYSTÉMATIQUES À RESPECTER

Basé sur l'analyse du projet TSD (conformité 87%) et les standards Go officiels.

### 🏷️ TYPES ET STRUCTURES - PascalCase

```go
// ✅ TOUJOURS générer ainsi
type Program struct { ... }
type TypeDefinition struct { ... }
type ConstraintValidator struct { ... }
type ReteNetwork struct { ... }

// ❌ JAMAIS générer ainsi
type program struct { ... }          // Incorrect
type typeDefinition struct { ... }   // Incorrect
```

**Règle :** Tous les types exportés en **PascalCase**, privés en **camelCase**.

### 🔧 FONCTIONS ET MÉTHODES

```go
// ✅ Fonctions exportées - PascalCase
func (v *ConstraintValidator) ValidateProgram(ctx context.Context) error
func (cp *ConstraintPipeline) BuildNetwork(types []TypeDefinition) (*ReteNetwork, error)

// ✅ Fonctions privées - camelCase  
func createTypeNodes(defs []TypeDefinition) map[string]*TypeNode
func validateNetwork(network *ReteNetwork) error

// ✅ Récepteurs - Abréviations courtes et cohérentes
func (cp *ConstraintPipeline) ...   // cp pour ConstraintPipeline
func (v *Validator) ...             // v pour Validator
func (n *Network) ...               // n pour Network
```

**Règles :**
- Fonctions exportées : **PascalCase**
- Fonctions privées : **camelCase**
- Récepteurs : **abréviations courtes cohérentes**

### 🔀 VARIABLES ET CONSTANTES

```go
// ✅ Variables - camelCase
var statePool = &sync.Pool{...}
var networkCache map[string]*ReteNetwork
ctx := context.Background()

// ✅ Constantes - camelCase
const choiceNoMatch = -1
const defaultTimeout = 30 * time.Second

// ✅ Constantes exportées - PascalCase si nécessaire
const MaxCacheSize = 1000
```

**Règle :** Variables et constantes en **camelCase**, exportées en **PascalCase**.

### 📂 STRUCTURE DE FICHIERS

```go
// ✅ Noms de fichiers - snake_case (convention TSD)
constraint_parser.go        // Parser pour contraintes
rete_network.go            // Réseau RETE  
type_converter.go          // Convertisseur de types
constraint_validator.go    // Validateur
test_utils.go             // Utilitaires de test

// ✅ Répertoires - snake_case
pkg/domain/
internal/config/
test/integration/
rete/pkg/nodes/
```

**Règle :** Fichiers et répertoires en **snake_case** pour cohérence avec TSD.

### 🧪 TESTS

```go
// ✅ Fonctions de test
func TestConstraintValidator_ValidateProgram(t *testing.T) { ... }
func TestReteNetwork_BuildFromConstraints(t *testing.T) { ... }

// ✅ Sous-tests  
t.Run("ValidConstraint", func(t *testing.T) { ... })
t.Run("InvalidSyntax", func(t *testing.T) { ... })

// ✅ Fichiers de test
constraint_validator_test.go
rete_network_test.go
integration_test.go
```

**Règles :**
- Tests : **TestType_Method** pattern
- Sous-tests : **PascalCase descriptif**
- Fichiers : **component_test.go**

### 🏗️ PACKAGES ET INTERFACES

```go
// ✅ Packages courts et descriptifs
package domain
package validator  
package nodes

// ✅ Interfaces - nommage simple et clair
type Storage interface {
    Store(fact Fact) error
    Retrieve(id string) (Fact, error)
}

type ConstraintValidator interface {
    ValidateProgram(prog *Program) error
    ValidateExpression(expr Expression) error
}
```

**Règles :**
- Packages : **noms courts, snake_case directories**
- Interfaces : **PascalCase, descriptif**

## 🎯 PATTERNS SPÉCIFIQUES TSD

### 🔍 Context et Error Handling

```go
// ✅ Toujours générer avec context en premier paramètre
func (v *ConstraintValidator) ValidateProgram(ctx context.Context, prog *Program) error

// ✅ Erreurs wrappées avec contexte
return fmt.Errorf("failed to validate constraint %s: %w", constraint.ID, err)

// ✅ Validation des paramètres
if prog == nil {
    return errors.New("program cannot be nil")
}
```

### 🏷️ Types de Domaine

```go
// ✅ Types du domaine TSD
type Constraint struct {
    ID       string
    Type     ConstraintType
    Expr     Expression
}

type ReteNode interface {
    Process(ctx context.Context, fact Fact) error
    GetID() string
}

type AlphaNode struct {
    nodeBase
    condition AlphaCondition
}
```

### 🔧 Builder Pattern (fréquent dans TSD)

```go
// ✅ Builders pour structures complexes
type NetworkBuilder struct {
    types map[string]*TypeDefinition
    rules []Rule
}

func NewNetworkBuilder() *NetworkBuilder {
    return &NetworkBuilder{
        types: make(map[string]*TypeDefinition),
        rules: make([]Rule, 0),
    }
}

func (b *NetworkBuilder) AddType(def *TypeDefinition) *NetworkBuilder {
    b.types[def.Name] = def
    return b
}

func (b *NetworkBuilder) Build() (*ReteNetwork, error) {
    // ...
}
```

## ✅ CHECKLIST GÉNÉRATION DE CODE

Avant de générer du code Go pour TSD :

- [ ] Types exportés en **PascalCase**
- [ ] Fonctions exportées en **PascalCase**, privées en **camelCase**
- [ ] Variables en **camelCase**
- [ ] Fichiers en **snake_case** 
- [ ] Récepteurs courts et cohérents
- [ ] Context en premier paramètre si applicable
- [ ] Error handling avec wrapping
- [ ] Validation des paramètres nil
- [ ] Tests avec pattern **TestType_Method**
- [ ] Documentation godoc sur exports

## 🔄 ÉVOLUTION CONTINUE

Ces règles seront :
- **Appliquées systématiquement** à tout nouveau code généré
- **Mises à jour** si les conventions TSD évoluent
- **Vérifiées** par les pre-commit hooks existants

---

**🎯 ENGAGEMENT :** Tout code généré respectera ces conventions pour maintenir la cohérence et qualité du projet TSD.

*Règles actives depuis le 13 novembre 2025*