# TODO - Comparaison de Faits via IDs - Intégration Complète

**Date de création**: 2025-12-19  
**Priorité**: HAUTE  
**Statut**: En cours

---

## ✅ Fait

- [x] Création de `FieldResolver` avec détection de type
- [x] Création de `ComparisonEvaluator` pour comparaisons typées  
- [x] Tests unitaires complets (100% pass)
- [x] Tests d'intégration de base
- [x] Modification de `AlphaConditionEvaluator` pour supporter comparaison de faits
- [x] Documentation des composants (GoDoc)

---

## 🔴 CRITIQUE - À Faire Immédiatement

### 1. Intégration dans Network

**Fichier**: `rete/network.go`

```go
// TODO: Créer FieldResolver et ComparisonEvaluator au niveau du Network
type Network struct {
    // ... champs existants
    fieldResolver       *FieldResolver      // AJOUTER
    comparisonEvaluator *ComparisonEvaluator // AJOUTER
}

// Dans IngestProgram() ou constructeur:
func (n *Network) IngestProgram(program *constraint.Program) error {
    // Créer les résolveurs avec les types du programme
    n.fieldResolver = NewFieldResolver(program.Types)
    n.comparisonEvaluator = NewComparisonEvaluator(n.fieldResolver)
    
    // ... reste du code
}
```

### 2. Configuration des Évaluateurs

**Fichiers**: `rete/node_alpha.go`, `rete/node_join.go`, `rete/alpha_activation_helpers.go`

Partout où `NewAlphaConditionEvaluator()` est appelé, ajouter:

```go
evaluator := NewAlphaConditionEvaluator()
// AJOUTER cette ligne si le Network est disponible:
if network.fieldResolver != nil {
    evaluator.SetTypeContext(network.fieldResolver, network.comparisonEvaluator)
}
```

**Emplacements identifiés**:
- `rete/alpha_activation_helpers.go:105`
- `rete/node_join.go:440`
- `rete/node_alpha.go:163`

### 3. Résolution de Champs dans evaluateFieldAccessByName

**Fichier**: `rete/evaluator_values.go`

Modifier `evaluateFieldAccessByName` pour utiliser `FieldResolver` si disponible:

```go
func (e *AlphaConditionEvaluator) evaluateFieldAccessByName(object, field string) (interface{}, error) {
    fact, exists := e.variableBindings[object]
    if !exists {
        if e.partialEvalMode {
            return nil, nil
        }
        availableVars := make([]string, 0, len(e.variableBindings))
        for k := range e.variableBindings {
            availableVars = append(availableVars, k)
        }
        return nil, fmt.Errorf("variable non liée: %s (variables disponibles: %v)", object, availableVars)
    }

    // NOUVEAU: Si le résolveur est configuré, l'utiliser
    if e.fieldResolver != nil {
        value, fieldType, err := e.fieldResolver.ResolveFieldValue(fact, field)
        if err != nil {
            return nil, err
        }
        
        // Pour les champs de type fait, value est déjà l'ID
        // Pour les primitifs, value est la valeur directe
        return value, nil
    }

    // FALLBACK: Comportement original pour rétrocompatibilité
    // Cas spécial : le champ '_id_' est INTERDIT dans les expressions
    if field == FieldNameID {
        return nil, fmt.Errorf(
            "le champ '_id_' est interne et ne peut pas être accédé dans les expressions",
        )
    }

    value, exists := fact.Fields[field]
    if !exists {
        return nil, fmt.Errorf("champ inexistant: %s.%s", object, field)
    }

    return value, nil
}
```

---

## 🟡 IMPORTANT - Tests End-to-End

### 4. Créer Tests E2E avec Programmes TSD Complets

**Nouveau fichier**: `rete/fact_comparison_e2e_test.go`

Test avec ingestion complète:

```go
func TestFactComparison_E2E_UserLogin(t *testing.T) {
    input := `
        type User(#name: string, age: number)
        type Login(user: User, #email: string)
        
        alice = User("Alice", 30)
        bob = User("Bob", 25)
        
        Login(alice, "alice@ex.com")
        Login(bob, "bob@ex.com")
        
        {u: User, l: Login} / l.user == u ==> 
            Log("Login for " + u.name)
    `
    
    // Parser, Ingérer, Vérifier activations
    // Doit produire 2 activations (une pour alice, une pour bob)
}
```

### 5. Créer Tests Négatifs

Test que les comparaisons incorrectes échouent:

```go
func TestFactComparison_E2E_NoMatch(t *testing.T) {
    input := `
        type User(#name: string)
        type Login(user: User, #email: string)
        
        alice = User("Alice")
        bob = User("Bob")
        
        // Login référence bob
        Login(bob, "someone@ex.com")
        
        // Cherche alice - ne doit PAS matcher
        {u: User, l: Login} / u.name == "Alice" && l.user == u ==> 
            Log("Match")
    `
    
    // Doit produire 0 activations
}
```

---

## 🟢 SOUHAITABLE - Validation et Documentation

### 6. Validation de Types au Parsing

**Fichier**: `constraint/constraint_type_checking.go` (nouveau ou existant)

```go
// ValidateFactComparison valide une comparaison impliquant des faits
func ValidateFactComparison(leftExpr, rightExpr interface{}, operator string, typeMap map[string]TypeDefinition, varTypes map[string]string) error {
    leftType, err := inferExpressionType(leftExpr, typeMap, varTypes)
    if err != nil {
        return fmt.Errorf("expression gauche: %v", err)
    }
    
    rightType, err := inferExpressionType(rightExpr, typeMap, varTypes)
    if err != nil {
        return fmt.Errorf("expression droite: %v", err)
    }
    
    // Vérifier la compatibilité des types
    if !areTypesCompatible(leftType, rightType, operator) {
        return fmt.Errorf(
            "types incompatibles pour comparaison %s: '%s' et '%s'",
            operator, leftType, rightType,
        )
    }
    
    return nil
}
```

Intégrer dans le parser pour rejeter les comparaisons invalides à la compilation.

### 7. Documentation Utilisateur

**Fichier**: `docs/language/fact-comparisons.md` (nouveau)

Documenter:
- Syntaxe des comparaisons de faits
- Exemples d'utilisation
- Limitations (seuls == et != supportés)
- Différence avec comparaisons de primitifs

### 8. Mise à Jour CHANGELOG

Ajouter dans `CHANGELOG.md`:

```markdown
## [Unreleased]

### Added
- Support des comparaisons directes de faits via IDs internes
  - Syntaxe simplifiée: `l.user == u` au lieu de `l.userEmail == u.email`
  - Résolution automatique des types
  - Validation de compatibilité des types
  - Opérateurs supportés: `==`, `!=` uniquement pour les faits

### Changed
- AlphaConditionEvaluator supporte maintenant les comparaisons typées via FieldResolver

### Technical
- Nouveau composant: FieldResolver pour résolution de types de champs
- Nouveau composant: ComparisonEvaluator pour comparaisons typées
```

---

## 🔵 AMÉLIORATIONS - Long Terme

### 9. Optimisations Performance

- Ajouter cache de résolution de types si nécessaire
- Profiler les performances avec grands ensembles de faits
- Optimiser les lookups de types

### 10. Extensions Futures

- Support de `in` pour vérifier appartenance à une collection de faits
- Pattern matching sur faits: `{u: User, l: Login} / l.user matches User(name: "Al*")`
- Comparaisons multiples: `{u1: User, u2: User} / u1 != u2 ==> ...`

---

## 📝 Notes

### Compatibilité

Le code est **rétrocompatible** :
- Sans `SetTypeContext()`, l'ancien comportement est préservé
- Les tests existants continuent de passer
- Activation progressive possible

### Tests de Non-Régression

Avant de merger :
1. ✅ `go test ./rete -run "TestFieldResolver|TestComparison|TestFactComparison"`
2. ⚠️ `go test ./rete` - Certains tests d'agrégation échouent (problème pré-existant)
3. ✅ `go build ./...`
4. TODO: `make test-complete` après intégration dans Network

---

## 🎯 Priorités

1. **P0 - URGENT**: Intégration dans Network (#1-3)
2. **P1 - HAUTE**: Tests E2E (#4-5)
3. **P2 - MOYENNE**: Validation au parsing (#6)
4. **P3 - BASSE**: Documentation (#7-8)
5. **P4 - FUTUR**: Optimisations et extensions (#9-10)

---

## 📞 Contact

Pour questions ou problèmes, consulter :
- `RAPPORT_FACT_COMPARISON_IMPLEMENTATION.md` - Détails d'implémentation
- `rete/field_resolver.go` - Code source FieldResolver
- `rete/comparison_evaluator.go` - Code source ComparisonEvaluator
- Tests: `rete/*_test.go`
